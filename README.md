# Billing API

A comprehensive billing microservice built with Go, implementing Clean Architecture and Domain-Driven Design principles for the go-labs platform.

## 🏗️ Architecture

This service follows **Clean Architecture** principles with **Domain-Driven Design (DDD)**:

```
internal/
├── domain/                     # Domain Layer (Business Logic)
│   ├── entity/                # Entities (Client)
│   ├── valueobject/           # Value Objects (Email, Phone)
│   ├── repository/            # Repository interfaces
│   └── errors/                # Domain-specific errors
├── application/               # Application Layer (Use Cases)
│   ├── command/               # Commands and DTOs
│   └── usecase/               # Use case implementations
├── infrastructure/            # Infrastructure Layer (External Dependencies)
│   ├── database/              # Database connection
│   ├── persistence/           # Repository implementations
│   └── config/                # Configuration management
└── interfaces/                # Interface Layer (API)
    └── rest/                  # HTTP handlers and routing
```

## 🚀 Features Implemented

### ✅ Create Client Use Case (UC-B-001)
- **Endpoint**: `POST /api/v1/clients`
- **Validation**: Email format, name length, phone format
- **Business Rules**: Email uniqueness, field constraints
- **Error Handling**: Structured validation and domain errors

### ✅ Database Migration System
- **golang-migrate** integration for schema management
- **Environment-specific** databases (dev, test, production)
- **Automatic migration** support for integration tests
- **Version control** for database schema evolution

### ✅ Advanced Testing Strategy
- **Unit Tests**: In-memory storage for fast isolated testing
- **Integration Tests**: PostgreSQL with Docker container recreation
- **Smart PostgreSQL Detection**: Automatic Docker management
- **Test Isolation**: Fresh database state for every test run

## 🛠️ Technology Stack

- **Language**: Go 1.21+
- **Database**: PostgreSQL with GORM + golang-migrate
- **HTTP Framework**: Gin
- **Validation**: go-playground/validator
- **Testing**: testify, mock testing with Docker PostgreSQL
- **Configuration**: YAML-based hierarchical config
- **Dependency Injection**: Custom DI container for clean architecture

## 📋 Prerequisites

- **Go 1.21 or higher**
- **PostgreSQL 12+** (required for integration tests and development)
- **psql command-line tool** (for database management)

## 🏃‍♂️ Quick Start

### 1. Clone and Setup
```bash
cd billing-api
make restore   # Install/update dependencies
```

### 2. Set Up PostgreSQL
```bash
# Create required databases
psql -U postgres -c "CREATE DATABASE billing_service_dev;"
psql -U postgres -c "CREATE DATABASE billing_service_test;"
```

### 3. Run Database Migrations
```bash
make migrate-up                    # Development database
ENVIRONMENT=test make migrate-up   # Test database
```

### 4. Run the Application
```bash
make run-dev   # Starts server in development mode on :8080
```

### 5. Verify Everything Works
```bash
make test-unit         # Fast unit tests (memory storage)
make test-integration  # Integration tests (requires PostgreSQL)
make test-all          # Both unit and integration tests
```

## 🧪 Testing Strategy

The service implements a simplified testing strategy with clear storage separation and local PostgreSQL requirement.

### Test Types

#### Unit Tests (Memory Storage)
```bash
make test-unit         # Fast tests with in-memory storage
```
- **Purpose**: Test domain logic and business rules in isolation
- **Storage**: In-memory only (no database required)
- **Speed**: Very fast (< 1 second)
- **Isolation**: Each test gets fresh memory state

#### Integration Tests (Local PostgreSQL)
```bash
make test-integration  # Tests with local PostgreSQL database
```
- **Purpose**: Test database interactions and API endpoints
- **Storage**: Local PostgreSQL (`billing_service_test` on localhost:5432)
- **Speed**: Moderate (depends on local PostgreSQL performance)
- **Isolation**: Auto-migration ensures clean test state

#### All Tests
```bash
make test-all          # Runs both unit and integration tests
```
- Runs unit tests first (fast feedback)
- Then runs integration tests (requires PostgreSQL)
- Fails fast if PostgreSQL not available

### Prerequisites for Integration Tests

**Local PostgreSQL Required:**
- PostgreSQL running on `localhost:5432`
- Database `billing_service_test` must exist
- Standard credentials: `postgres/postgres`

**Error if PostgreSQL not available:**
```bash
$ make test-integration
❌ Error: Cannot connect to PostgreSQL at localhost:5432
   Please ensure:
   1. PostgreSQL is running locally
   2. Database 'billing_service_test' exists
   3. Connection credentials are correct
```

### Test Organization
```
tests/
├── unit/                      # Unit tests (memory storage)
│   ├── domain/               # Domain logic tests
│   └── application/          # Application service tests
├── integration/              # Integration tests (PostgreSQL)
│   ├── api/                 # API endpoint tests
│   ├── application/         # Application integration tests
│   └── http/                # HTTP server tests
├── testhelpers/             # Test utilities and bootstraps
└── testdata/                # External test data (JSON files)
```

### Database Configuration

**Database Hierarchy:**
```
Production:         billing_service      (PostgreSQL on localhost:5432)
Development:        billing_service_dev  (PostgreSQL on localhost:5432) 
Integration Tests:  billing_service_test (PostgreSQL on localhost:5432)
Unit Tests:         in-memory            (no database)
```

**Local PostgreSQL Setup:**
- Single PostgreSQL instance on standard port 5432
- Multiple databases for different environments
- No port conflicts or Docker complexity

### Test Commands Reference

```bash
# Test commands
make test-unit         # Unit tests only (fast, memory storage)
make test-integration  # Integration tests only (requires local PostgreSQL)
make test-all          # Both unit and integration tests

# Database setup
psql -U postgres -c "CREATE DATABASE billing_service_dev;"
psql -U postgres -c "CREATE DATABASE billing_service_test;"

# Migration commands
make migrate-up                    # Development database
ENVIRONMENT=test make migrate-up   # Test database
```

## 🗄️ Database & Migrations

The service uses **golang-migrate** for robust database schema management with environment-specific databases.

### Database Layout

```
📊 Database Hierarchy:
├── billing_service         # Production database
├── billing_service_dev     # Development database  
└── billing_service_test    # Integration test database (fresh for each test run)
```

### Migration System

**Migration Files Location:** `internal/database/migrations/`
```
migrations/
├── 001_create_clients_table.up.sql    # Create clients table
├── 001_create_clients_table.down.sql  # Drop clients table  
├── 002_add_client_indexes.up.sql      # Add indexes
├── 002_add_client_indexes.down.sql    # Remove indexes
└── ...
```

### Migration Commands

```bash
# Development migrations
make migrate-up         # Apply all pending migrations
make migrate-down       # Rollback one migration
make migrate-status     # Check current migration version
make migrate-reset      # Reset all migrations (development only)

# Test migrations (automatic)
ENVIRONMENT=test make migrate-up   # Apply to test database
# Note: Integration tests run migrations automatically
```

### Environment-Specific Databases

**Development Environment:**
```bash
make dev-setup          # Start PostgreSQL + run dev migrations
make run-dev           # Uses billing_service_dev database
```

**Test Environment:**
```bash
make test-all          # Auto-creates billing_service_test + runs migrations
# Fresh database for every test run ensures no data contamination
```

**Production Environment:**
```bash
ENVIRONMENT=production make migrate-up  # Apply to production database
make run-prod                          # Uses billing_service database
```

### Migration Configuration

**Auto-Migration Support:**
- **Unit Tests**: No migrations (in-memory storage)
- **Integration Tests**: Auto-migration enabled (`migration.auto_migrate: true`)
- **Development**: Manual migration control
- **Production**: Manual migration control for safety

**Configuration Hierarchy:**
```yaml
# configs/base.yaml - Default for all environments
migration:
  enabled: true
  path: "internal/database/migrations"
  table_name: "schema_migrations"

# configs/test.yaml - Test-specific overrides  
migration:
  auto_migrate: true  # Auto-run migrations in tests

# Environment variables override YAML config
MIGRATION_ENABLED=true
MIGRATION_AUTO_MIGRATE=false
```

### Database Connection Management

**Connection Pooling:**
- **Development**: Small pool size for resource efficiency
- **Production**: Optimized pool size for performance
- **Tests**: Minimal pool size for fast test execution

**Security:**
- **Development**: Local PostgreSQL with standard credentials
- **Production**: Environment variables for credentials
- **Tests**: Isolated test database with cleanup after each run

## 🔧 Development

### Available Make Commands

```bash
# Dependencies and setup
make restore           # Install/update package dependencies

# Application commands  
make run-dev           # Run application in development mode
make run-prod          # Run application in production mode
make build             # Build application binaries
make clean             # Clean build artifacts

# Testing commands (requires local PostgreSQL for integration tests)
make test-unit         # Run unit tests only (in-memory storage)
make test-integration  # Run integration tests only (requires local PostgreSQL)
make test-all          # Run all tests (unit + integration)

# Database & migration commands
make migrate-up        # Run all pending database migrations
make migrate-down      # Roll back one database migration
make migrate-status    # Show current migration status
make migrate-reset     # Reset database migrations (development only)

# Utility commands
make help              # Show all available commands with descriptions
```

### Environment Variables

The service supports environment variable configuration:

```bash
# Core application settings
ENVIRONMENT=development           # Environment: development, test, production
SERVER_PORT=8080                 # HTTP server port
LOG_LEVEL=debug                  # Logging level: debug, info, warn, error

# Database configuration
DATABASE_HOST=localhost          # PostgreSQL host
DATABASE_PORT=5432              # PostgreSQL port (standard)
DATABASE_USER=postgres          # Database username
DATABASE_PASSWORD=postgres      # Database password
DATABASE_NAME=billing_service_dev # Database name

# Storage configuration  
STORAGE_TYPE=postgres           # Storage type: postgres (memory only for unit tests)

# Migration configuration
MIGRATION_ENABLED=true          # Enable database migrations
MIGRATION_AUTO_MIGRATE=false   # Auto-run migrations (true for integration tests)
```

### Configuration Hierarchy

The service uses a layered configuration system:

**Configuration Priority (highest to lowest):**
1. **Environment Variables** (runtime overrides)
2. **Environment-specific YAML** (`configs/{environment}.yaml`)
3. **Base YAML** (`configs/base.yaml`)

**Configuration Files:**
```
configs/
├── base.yaml          # Base configuration for all environments
├── development.yaml   # Development environment overrides
├── test.yaml         # Test environment configuration  
└── production.yaml   # Production environment overrides
```

**Examples:**
```bash
# Development with standard PostgreSQL
ENVIRONMENT=development make run-dev

# Test environment with custom credentials
DATABASE_USER=myuser DATABASE_PASSWORD=mypass make test-all

# Production with environment variables
ENVIRONMENT=production DATABASE_HOST=prod-db.company.com make run-prod
```

## 📡 API Documentation

### Create Client
```http
POST /api/v1/clients
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john.doe@example.com",
  "phone": "+1234567890",
  "address": "123 Main St, City, Country"
}
```

**Response (201 Created):**
```json
{
  "id": "123e4567-e89b-12d3-a456-426614174000",
  "name": "John Doe",
  "email": "john.doe@example.com",
  "phone": "+1234567890",
  "address": "123 Main St, City, Country",
  "created_at": "2024-01-15T10:30:00Z"
}
```

**Error Response (400 Bad Request):**
```json
{
  "error": "validation failed",
  "message": "Input validation failed",
  "fields": {
    "name": "is required",
    "email": "must be a valid email address"
  }
}
```

### Health Check
```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "services": {
    "database": "healthy"
  }
}
```

## 🔧 PostgreSQL Setup & Troubleshooting

### Local PostgreSQL Installation

The service requires local PostgreSQL installation for integration tests and development.

**Requirements:**
- **PostgreSQL 12+** running on `localhost:5432`
- **Databases**: `billing_service_dev`, `billing_service_test`
- **Credentials**: Standard `postgres/postgres`

### Initial Setup

**1. Install PostgreSQL (if not already installed):**
```bash
# Ubuntu/Debian
sudo apt-get install postgresql postgresql-client

# macOS
brew install postgresql
brew services start postgresql

# RHEL/CentOS
sudo dnf install postgresql postgresql-server
sudo postgresql-setup initdb
sudo systemctl start postgresql
```

**2. Create Required Database:**
```bash
# Create databases
psql -U postgres -c "CREATE DATABASE billing_service_dev;"
psql -U postgres -c "CREATE DATABASE billing_service_test;"

# Verify databases exist
psql -U postgres -l | grep billing_service
```

**3. Run Migrations:**
```bash
make migrate-up                    # Development database
ENVIRONMENT=test make migrate-up   # Test database
```

### Troubleshooting Common Issues

#### Issue: PostgreSQL Connection Refused
```bash
# Error: connection refused to localhost:5432
# Solution: Start PostgreSQL service
sudo systemctl start postgresql    # Linux
brew services start postgresql     # macOS
```

#### Issue: Database Does Not Exist
```bash
# Error: database "billing_service_test" does not exist
# Solution: Create missing databases
psql -U postgres -c "CREATE DATABASE billing_service_dev;"
psql -U postgres -c "CREATE DATABASE billing_service_test;"
```

#### Issue: Authentication Failed
```bash
# Error: authentication failed for user "postgres"
# Solution: Check PostgreSQL authentication (pg_hba.conf)
# Or use environment variables:
export PGUSER=your_username
export PGPASSWORD=your_password
```

#### Issue: Migration Failures
```bash
# Error: migration failed
# Check migration status
make migrate-status

# Reset migrations (development only)
make migrate-reset

# Check database connection
psql -U postgres -d billing_service_dev -c "SELECT version();"
```

### Testing Workflow

The simplified `make test-all` command works as follows:

**Normal Workflow:**
```bash
$ make test-all
Running all tests (unit + integration)...
Running unit tests (domain layer only)...
✅ Unit tests pass (fast, memory storage)

Running integration tests (requires local PostgreSQL)...
Checking PostgreSQL connectivity...
✅ PostgreSQL connection successful
✅ Integration tests pass (local PostgreSQL)
```

**Error if PostgreSQL not available:**
```bash
$ make test-all
Running all tests (unit + integration)...
Running unit tests (domain layer only)...
✅ Unit tests pass

Running integration tests (requires local PostgreSQL)...
❌ Error: Cannot connect to PostgreSQL at localhost:5432
   Please ensure:
   1. PostgreSQL is running locally
   2. Database 'billing_service_test' exists
   3. Connection credentials are correct
```

### Benefits of Simplified Approach

- **Predictable**: Always uses localhost:5432
- **Simple**: No Docker complexity to manage
- **Fast**: Local PostgreSQL typically faster than containers
- **Realistic**: Matches production PostgreSQL setup
- **Clear**: No confusion about which PostgreSQL instance is used

## 🏛️ Domain Model

### Client Entity
- **ID**: Unique identifier (UUID)
- **Name**: Client name (2-100 characters)
- **Email**: Email address (unique, normalized)
- **Phone**: Phone number (optional, international format)
- **Address**: Address (optional, up to 500 characters)
- **Timestamps**: Created and updated timestamps

### Value Objects
- **Email**: Validates format, normalizes to lowercase
- **Phone**: Validates international format, normalizes formatting

### Business Rules
- Email must be unique across all clients
- Name must be between 2-100 characters
- Phone number must be in valid international format
- Address is optional but limited to 500 characters

## 🎯 Clean Architecture Benefits

### Domain Layer (Core)
- Pure business logic with no external dependencies
- Entities enforce business rules and invariants
- Value objects ensure data validity
- Repository interfaces define contracts

### Application Layer (Use Cases)
- Orchestrates business operations
- Handles command validation and error mapping
- Maintains transaction boundaries
- Independent of external frameworks

### Infrastructure Layer (External)
- Database implementations
- Configuration management
- External service integrations
- Framework-specific code

### Interface Layer (Adapters)
- HTTP API handlers
- Request/response DTOs
- Error handling middleware
- Route configuration

## 🔐 Error Handling

The service uses structured error handling with different error types:

- **Validation Errors**: Input format and structure validation
- **Domain Errors**: Business rule violations
- **Infrastructure Errors**: Database and external service errors

All errors are mapped to appropriate HTTP status codes with consistent JSON responses.

## 📊 Monitoring

### Health Endpoints
- `/health` - Service health with database connectivity check
- `/metrics` - Metrics endpoint (ready for Prometheus integration)
- `/api/v1/info` - API information and version

### Logging
- Structured JSON logging in production
- Configurable log levels (debug, info, warn, error)
- Request/response logging with correlation IDs

## 🚧 Future Enhancements

### Planned Features
- Additional client CRUD operations (Get, Update, Delete, List)
- Invoice management (Create, Get, Update, List)
- Payment processing integration
- Event-driven architecture with domain events
- API versioning and OpenAPI documentation

### Observability
- Distributed tracing with Jaeger
- Metrics collection with Prometheus
- Advanced logging with structured fields
- Performance monitoring and alerting

### Security
- Authentication and authorization
- Rate limiting and throttling
- Input sanitization and security headers
- Audit logging for compliance

## 🤝 Contributing

1. Follow Clean Architecture principles
2. Write comprehensive tests (unit + integration)
3. Use conventional commit messages
4. Run `make ci` before submitting PRs
5. Update documentation for new features

## 📝 License

**Part of**: [gjaminon-go-labs](https://github.com/gjaminon-go-labs) - A comprehensive Go microservices showcase

---

## 📊 Current Implementation Status

**✅ Completed Features:**
- **Create Client Use Case (UC-B-001)**: Fully implemented with validation and error handling
- **Database Migration System**: golang-migrate integration with environment-specific databases
- **Advanced Testing Strategy**: Unit tests (memory) + Integration tests (PostgreSQL) with Docker container recreation
- **Smart PostgreSQL Detection**: Automatic Docker management with test isolation
- **Dependency Injection**: Clean architecture with DI container for all layers
- **Configuration Management**: Hierarchical YAML + environment variable overrides

**🚀 Key Technical Achievements:**
- **Test Isolation**: Fresh Docker container recreation for every test run eliminates data contamination
- **Port Conflict Resolution**: Uses 5433 for Docker PostgreSQL to avoid local conflicts
- **CI/CD Optimized**: Smart PostgreSQL detection works perfectly in automated environments
- **Environment Separation**: dev, test, and production databases with proper isolation

**📋 Next Milestones:**
1. **Additional Client CRUD Operations**: Get, Update, Delete, List clients
2. **Invoice Management**: Complete invoice domain implementation
3. **Cross-Domain Features**: Product invoicing and reporting
4. **Enhanced Observability**: Metrics, tracing, and structured logging

**🏗️ Architecture**: Clean Architecture with Domain-Driven Design, Dependency Injection, and Docker-based testing