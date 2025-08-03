# Test Data Isolation - Database Cleanup Strategy

## Overview

This document describes the test data isolation system implemented for the billing-api service. The system ensures that integration tests start with a clean database state, preventing test pollution and non-deterministic results.

## Problem Solved

Without data cleanup, integration tests suffer from:
- **Test Pollution**: Later tests affected by data from earlier tests
- **Non-Deterministic Results**: Tests passing/failing based on existing data
- **Unique Constraint Violations**: Repeated tests failing due to duplicate data
- **Debugging Difficulty**: Hard to isolate test failures

## Implementation Strategy

### Selected Approach: Table-Level DELETE Operations

We chose DELETE operations over other approaches because:
- **Secure**: Works with application user permissions (no DDL required)
- **Fast**: Efficient for test data volumes
- **Safe**: Respects foreign key constraints
- **Reliable**: No special database permissions needed

### Architecture

```
Integration Test
    ↓
NewIntegrationTestServer()
    ↓
DatabaseCleaner.CleanupTestData()
    ↓
DELETE FROM billing.storage_records
DELETE FROM billing.clients
    ↓
Clean Test Environment
```

## Configuration

### Test Cleanup Settings

Integration tests automatically enable cleanup by default:

```go
// In IntegrationTestConfig()
TestCleanupEnabled: true,  // Enable test data cleanup
TestCleanupOnSetup: true,  // Cleanup before each test
```

### Environment Variables

Override cleanup behavior with environment variables:
```bash
TEST_CLEANUP_ENABLED=false  # Disable cleanup (for debugging)
TEST_CLEANUP_ON_SETUP=false # Disable automatic cleanup on setup
```

## Usage Patterns

### Automatic Cleanup (Default)

```go
func TestMyIntegration(t *testing.T) {
    // This automatically triggers cleanup
    server := testhelpers.NewIntegrationTestServer()
    
    // Test runs with clean database state
    // ...
}
```

### Manual Cleanup Control

```go
func TestWithManualCleanup(t *testing.T) {
    // No automatic cleanup
    stack := testhelpers.NewIntegrationTestStackNoCleanup()
    
    // Manual cleanup when needed
    err := stack.DatabaseCleaner.CleanupTestData()
    assert.NoError(t, err)
    
    // Test logic...
}
```

### Debugging Without Cleanup

```go
func TestDebugging(t *testing.T) {
    // For debugging - preserves data between runs
    server := testhelpers.NewIntegrationTestServerNoCleanup()
    
    // Test logic...
    
    // Manually inspect database after test
}
```

## Available Test Helpers

### Automatic Cleanup Helpers

- `NewIntegrationTestServer()` - HTTP server with automatic cleanup
- `NewIntegrationTestStack()` - Full test stack with automatic cleanup
- `NewCleanIntegrationTestServer()` - Explicitly clean server (same as default)
- `NewCleanIntegrationTestStack()` - Explicitly clean stack (same as default)

### No-Cleanup Helpers (for debugging)

- `NewIntegrationTestServerNoCleanup()` - Server without automatic cleanup
- `NewIntegrationTestStackNoCleanup()` - Stack without automatic cleanup

### Manual Cleanup Functions

- `CleanupIntegrationTestData()` - Standalone cleanup function
- `stack.DatabaseCleaner.CleanupTestData()` - Full cleanup
- `stack.DatabaseCleaner.CleanupSpecificTable(tableName)` - Single table
- `stack.DatabaseCleaner.VerifyCleanState()` - Verify all tables empty
- `stack.DatabaseCleaner.GetTableCounts()` - Get record counts for debugging

## Implementation Details

### Database Cleaner

The `DatabaseCleaner` provides low-level cleanup operations:

```go
// Create cleaner
cleaner := NewDatabaseCleaner(db)

// Full cleanup
err := cleaner.CleanupTestData()

// Specific table
err := cleaner.CleanupSpecificTable("clients")

// Verification
err := cleaner.VerifyCleanState()

// Debugging
counts, err := cleaner.GetTableCounts()
```

### Cleanup Order

Tables are cleaned in dependency order to respect foreign key constraints:

```go
tablesToClean := []string{
    "storage_records", // No foreign keys, safe first
    "clients",         // No foreign keys, safe
    // Future tables will be ordered appropriately
}
```

### Error Handling

Cleanup errors cause test panics to ensure test reliability:

```go
if err := cleaner.CleanupTestData(); err != nil {
    panic("Failed to cleanup test data: " + err.Error())
}
```

## Security Model

### Permission Requirements

- **Application User**: DELETE permissions on test tables (already has this)
- **Migration User**: Not needed for cleanup operations
- **No Special Privileges**: Works with standard DML permissions

### Safe Operations

The cleanup system uses:
- DELETE statements (not TRUNCATE)
- No session parameter changes
- No foreign key constraint manipulation
- Standard SQL operations only

## Performance Characteristics

### Timing

```
Clean Database:     ~5ms
Small Dataset:      ~10ms  (< 100 records)
Medium Dataset:     ~50ms  (< 1000 records)
Large Dataset:      ~200ms (1000+ records)
```

### Optimization

- Uses single DELETE statements per table
- Minimal database round trips
- No transaction overhead
- Efficient for test data volumes

## Monitoring and Debugging

### Cleanup Logging

All cleanup operations are logged:

```
🧹 Cleaning up test data...
🗑️  Cleaned table: billing.storage_records (4 records deleted)
🗑️  Cleaned table: billing.clients (0 records deleted)
✅ Test data cleanup completed
```

### Debugging Tools

```go
// Check current state
counts, err := cleaner.GetTableCounts()
fmt.Printf("Current state: %+v\n", counts)

// Verify cleanup
err = cleaner.VerifyCleanState()
if err != nil {
    fmt.Printf("Cleanup verification failed: %v\n", err)
}
```

### Disable Cleanup for Investigation

```go
// Use no-cleanup helpers to preserve data
stack := testhelpers.NewIntegrationTestStackNoCleanup()

// Run tests and investigate database manually
```

## Testing the Cleanup System

The cleanup system includes its own tests:

```bash
# Test cleanup functionality
go test -v ./tests/integration/testhelpers/

# Test isolation between test runs
go test -v ./tests/integration/testhelpers/ -count=3
```

## Migration to New System

### Existing Tests

All existing integration tests automatically benefit from the new cleanup system:
- No code changes required
- Automatic cleanup on every test run
- Improved test reliability

### Opt-out for Debugging

If debugging requires preserved data:

```go
// Replace this:
server := testhelpers.NewIntegrationTestServer()

// With this:
server := testhelpers.NewIntegrationTestServerNoCleanup()
```

## Future Enhancements

### Planned Features

1. **Selective Cleanup**: Clean only specific data types
2. **Test Fixtures**: Pre-populate known test data
3. **Cleanup Metrics**: Track cleanup performance
4. **Parallel Test Safety**: Enhanced isolation for parallel test execution

### Table Addition

When adding new tables to the system:

1. Add table name to `tablesToClean` slice in correct dependency order
2. Update documentation
3. Test cleanup with new table
4. Verify foreign key constraint handling

## Troubleshooting

### Common Issues

**Q: Tests failing with "permission denied"**  
A: Ensure application user has DELETE permissions on all test tables

**Q: Foreign key constraint violations during cleanup**  
A: Check table cleanup order in `tablesToClean` slice

**Q: Cleanup taking too long**  
A: Check for large datasets, consider test data volume limits

**Q: Tests still see data from previous runs**  
A: Verify cleanup is enabled and check for errors in cleanup logs

### Manual Recovery

If cleanup system fails:

```sql
-- Manual cleanup (as migration user)
DELETE FROM billing.storage_records;
DELETE FROM billing.clients;

-- Verify clean state
SELECT 'storage_records' as table_name, COUNT(*) as count FROM billing.storage_records
UNION ALL
SELECT 'clients' as table_name, COUNT(*) as count FROM billing.clients;
```

## Best Practices

### For Test Authors

1. **Use default cleanup**: `NewIntegrationTestServer()` for normal tests
2. **Debug with no-cleanup**: `NewIntegrationTestServerNoCleanup()` when investigating
3. **Verify test isolation**: Tests should pass in any order
4. **Avoid large datasets**: Keep test data volumes reasonable
5. **Use manual cleanup**: When you need precise control over timing

### For Test Data

1. **Unique identifiers**: Use unique emails, names to avoid conflicts
2. **Reasonable volume**: Don't create excessive test data
3. **Clean patterns**: Use consistent naming for test data
4. **Avoid external dependencies**: Keep tests self-contained

---

*This system ensures reliable, isolated integration tests while maintaining performance and simplicity.*