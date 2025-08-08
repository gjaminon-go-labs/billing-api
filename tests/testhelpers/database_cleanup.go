// Database Cleanup Utilities for Integration Tests
//
// This file provides database cleanup functionality for integration tests.
// Provides: Table truncation, data isolation, test setup utilities
// Pattern: Test helper utilities for maintaining clean database state
// Used by: Integration test setup, test isolation, debugging support
package testhelpers

import (
	"fmt"
	"log"

	"gorm.io/gorm"
)

// DatabaseCleaner provides methods for cleaning up test data
type DatabaseCleaner struct {
	db *gorm.DB
}

// NewDatabaseCleaner creates a new database cleaner instance
func NewDatabaseCleaner(db *gorm.DB) *DatabaseCleaner {
	return &DatabaseCleaner{
		db: db,
	}
}

// CleanupTestData removes all data from billing schema tables
// This method deletes data in the correct order to handle foreign key constraints
func (c *DatabaseCleaner) CleanupTestData() error {
	log.Println("🧹 Cleaning up test data...")

	// List of tables in dependency order (child tables first)
	// This ensures foreign key constraints are respected during cleanup
	tablesToClean := []string{
		"storage_records", // No foreign keys, safe to clean first
		"clients",         // No foreign keys, safe to clean
	}

	// Delete data from each table (safer than TRUNCATE for permissions)
	for _, table := range tablesToClean {
		if err := c.deleteFromTable(table); err != nil {
			return fmt.Errorf("failed to clean table %s: %w", table, err)
		}
	}

	log.Println("✅ Test data cleanup completed")
	return nil
}

// deleteFromTable deletes all data from a specific table
// This is safer than TRUNCATE as it doesn't require special permissions
func (c *DatabaseCleaner) deleteFromTable(tableName string) error {
	query := fmt.Sprintf("DELETE FROM billing.%s", tableName)

	result := c.db.Exec(query)
	if result.Error != nil {
		return fmt.Errorf("failed to delete from table %s: %w", tableName, result.Error)
	}

	log.Printf("🗑️  Cleaned table: billing.%s (%d records deleted)", tableName, result.RowsAffected)
	return nil
}

// truncateTable truncates a specific table and resets its sequences
// This method requires higher privileges and is kept for optional use
func (c *DatabaseCleaner) truncateTable(tableName string) error {
	// Use TRUNCATE with RESTART IDENTITY to reset any auto-incrementing sequences
	// CASCADE option handles any remaining foreign key dependencies
	query := fmt.Sprintf("TRUNCATE TABLE billing.%s RESTART IDENTITY CASCADE", tableName)

	if err := c.db.Exec(query).Error; err != nil {
		return fmt.Errorf("failed to truncate table %s: %w", tableName, err)
	}

	log.Printf("🗑️  Truncated table: billing.%s", tableName)
	return nil
}

// VerifyCleanState checks if all test tables are empty
// This is useful for debugging and ensuring cleanup worked correctly
func (c *DatabaseCleaner) VerifyCleanState() error {
	tablesToCheck := []string{"clients", "storage_records"}

	for _, table := range tablesToCheck {
		var count int64
		query := fmt.Sprintf("SELECT COUNT(*) FROM billing.%s", table)

		if err := c.db.Raw(query).Scan(&count).Error; err != nil {
			return fmt.Errorf("failed to check table %s: %w", table, err)
		}

		if count > 0 {
			return fmt.Errorf("table billing.%s is not empty: contains %d records", table, count)
		}
	}

	log.Println("✅ All test tables are clean")
	return nil
}

// GetTableCounts returns the number of records in each test table
// Useful for debugging and understanding test data state
func (c *DatabaseCleaner) GetTableCounts() (map[string]int64, error) {
	tablesToCheck := []string{"clients", "storage_records"}
	counts := make(map[string]int64)

	for _, table := range tablesToCheck {
		var count int64
		query := fmt.Sprintf("SELECT COUNT(*) FROM billing.%s", table)

		if err := c.db.Raw(query).Scan(&count).Error; err != nil {
			return nil, fmt.Errorf("failed to count records in table %s: %w", table, err)
		}

		counts[table] = count
	}

	return counts, nil
}

// CleanupSpecificTable deletes data from only a specific table
// Useful for targeted cleanup in specific test scenarios
func (c *DatabaseCleaner) CleanupSpecificTable(tableName string) error {
	log.Printf("🧹 Cleaning up table: billing.%s", tableName)

	if err := c.deleteFromTable(tableName); err != nil {
		return fmt.Errorf("failed to cleanup table %s: %w", tableName, err)
	}

	return nil
}
