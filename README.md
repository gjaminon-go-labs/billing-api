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
- **Docker** (essential for testing and local development)
- **PostgreSQL 12+** (optional - Docker provides this automatically)

## 🏃‍♂️ Quick Start

### 1. Clone and Setup
```bash
cd billing-api
make restore   # Install/update dependencies
```

### 2. Development Environment Setup
```bash
make dev-setup # Starts PostgreSQL + runs migrations automatically
```

### 3. Run the Application
```bash
make run-dev   # Starts server in development mode on :8080
```

### 4. Verify Everything Works
```bash
make test-all  # Runs all tests with smart PostgreSQL detection
```

### Alternative: Manual Setup
```bash
# 1. Start PostgreSQL manually
make docker-up        # Start PostgreSQL container (port 5433)

# 2. Run migrations
make migrate-up       # Apply database migrations

# 3. Run application
make run-dev          # Start development server
```

## 🧪 Testing Strategy

The service implements a comprehensive testing strategy with automatic PostgreSQL management and test isolation.

### Test Types

#### Unit Tests (Memory Storage)
```bash
make test-unit         # Fast tests with in-memory storage
```
- **Purpose**: Test domain logic and business rules in isolation
- **Storage**: In-memory (no database required)
- **Speed**: Very fast (< 1 second)
- **Isolation**: Each test gets fresh memory state

#### Integration Tests (PostgreSQL Storage)
```bash
make test-integration  # Tests with real PostgreSQL database
```
- **Purpose**: Test database interactions and API endpoints
- **Storage**: PostgreSQL test database (`billing_service_test`)
- **Speed**: Moderate (with Docker container recreation)
- **Isolation**: Fresh Docker container for every test run

#### All Tests (Smart PostgreSQL Detection)
```bash
make test-all          # Runs both unit and integration tests
```

**Smart Detection Logic:**
1. **Docker container exists** → Recreates for fresh test isolation
2. **Local PostgreSQL running** → Uses existing instance  
3. **No PostgreSQL found** → Creates fresh Docker container

### Test Isolation Benefits

```bash
# First run
make test-all
# 🔄 Docker PostgreSQL container found - recreating for fresh test isolation...
# ✅ All tests pass with clean database state

# Second run  
make test-all
# 🔄 Docker PostgreSQL container found - recreating for fresh test isolation...
# ✅ All tests pass with completely fresh state (no data contamination)
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

### PostgreSQL Configuration for Tests

**Port Configuration:**
- **Docker PostgreSQL**: Port 5433 (avoids conflicts with local PostgreSQL on 5432)
- **Test Database**: `billing_service_test` (separate from development)
- **Custom Port**: Set `DOCKER_DB_PORT=xxxx` for different port

**Database Hierarchy:**
```
Production:         billing_service      (PostgreSQL)
Development:        billing_service_dev  (PostgreSQL) 
Integration Tests:  billing_service_test (PostgreSQL)
Unit Tests:         in-memory            (no database)
```

### Test Commands Reference

```bash
# Basic test commands
make test-unit         # Unit tests only (fast, memory storage)
make test-integration  # Integration tests only (PostgreSQL)
make test-all         # All tests (smart PostgreSQL management)

# PostgreSQL management  
make check-postgres   # Check PostgreSQL status and recreate if needed
make docker-up        # Start PostgreSQL container manually
make docker-down      # Stop PostgreSQL container
make recreate-docker-postgres # Force recreate container

# Test environment setup
make test-setup       # Set up test databases and run migrations
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
make dev-setup         # Complete development setup (PostgreSQL + migrations)

# Application commands  
make run-dev           # Run application in development mode
make run-prod          # Run application in production mode
make build             # Build application binaries
make clean             # Clean build artifacts

# Testing commands
make test-unit         # Run unit tests only (in-memory storage)
make test-integration  # Run integration tests only (PostgreSQL test DB)
make test-all          # Run all tests (smart: fresh Docker or existing local PostgreSQL)
make test-setup        # Set up PostgreSQL test environment

# Database & migration commands
make migrate-up        # Run all pending database migrations
make migrate-down      # Roll back one database migration
make migrate-status    # Show current migration status
make migrate-reset     # Reset database migrations (development only)

# PostgreSQL management
make docker-up         # Start PostgreSQL container (port 5433)
make docker-down       # Stop PostgreSQL container  
make check-postgres    # Check PostgreSQL status, recreate Docker for test isolation
make recreate-docker-postgres # Recreate container for fresh state

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
DATABASE_PORT=5432              # PostgreSQL port (5433 for Docker)
DATABASE_USER=postgres          # Database username
DATABASE_PASSWORD=postgres      # Database password
DATABASE_NAME=billing_service_dev # Database name

# Storage configuration  
STORAGE_TYPE=postgres           # Storage type: postgres (memory only for tests)

# Migration configuration
MIGRATION_ENABLED=true          # Enable database migrations
MIGRATION_AUTO_MIGRATE=false   # Auto-run migrations (true for tests)

# Docker configuration
DOCKER_DB_PORT=5433            # PostgreSQL port for Docker container
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
# Development with environment override
ENVIRONMENT=development DATABASE_PORT=5433 make run-dev

# Test environment with custom Docker port
DOCKER_DB_PORT=15432 make test-all

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

## 🐳 Docker Configuration & Troubleshooting

### Docker PostgreSQL Management

The service uses Docker to provide consistent PostgreSQL environments across development and testing.

**Default Configuration:**
- **Container Name**: `billing-postgres`
- **Docker Port**: `5433` (avoids conflicts with local PostgreSQL on 5432)
- **PostgreSQL Version**: `15`
- **Databases Created**: `billing_service_dev`, `billing_service_test`

### Port Configuration

**Avoiding PostgreSQL Conflicts:**
```bash
# Default setup (recommended)
make docker-up         # Uses port 5433 (conflict-free)

# Custom port if needed  
DOCKER_DB_PORT=15432 make docker-up   # Uses port 15432

# Check what's running on PostgreSQL ports
make check-postgres    # Smart detection of existing PostgreSQL
```

**Port Usage:**
- **5432**: Local PostgreSQL (if installed)
- **5433**: Docker PostgreSQL (default)
- **Custom**: Set via `DOCKER_DB_PORT` environment variable

### Troubleshooting Common Issues

#### Issue: Container Name Conflict
```bash
# Error: container name 'billing-postgres' is already in use
# Solution: Recreate the container
make recreate-docker-postgres
```

#### Issue: Port Already in Use
```bash
# Error: port 5433 is already allocated
# Solution 1: Use different port
DOCKER_DB_PORT=15432 make docker-up

# Solution 2: Stop conflicting container
docker ps | grep 5433
docker stop <container-name>
make docker-up
```

#### Issue: PostgreSQL Connection Refused
```bash
# Error: connection refused to localhost:5433
# Check container status
docker ps | grep billing-postgres

# If not running, start it
make docker-up

# If still failing, recreate
make recreate-docker-postgres
```

#### Issue: Database Does Not Exist
```bash
# Error: database "billing_service_test" does not exist
# Solution: Recreate container (auto-creates databases)
make recreate-docker-postgres

# Or create manually
docker exec billing-postgres psql -U postgres -c "CREATE DATABASE billing_service_test;"
```

#### Issue: Migration Failures
```bash
# Error: migration failed
# Check migration status
make migrate-status

# Reset migrations (development only)
make migrate-reset

# Check database connection
docker exec billing-postgres psql -U postgres -l
```

### Smart PostgreSQL Detection Scenarios

The `make test-all` command intelligently handles different PostgreSQL setups:

**Scenario 1: Docker PostgreSQL Running**
```bash
$ make test-all
🔍 Checking PostgreSQL status on port 5433...
🔄 Docker PostgreSQL container found - recreating for fresh test isolation...
🗑️  Removing existing container...
🐘 Creating fresh PostgreSQL container...
✅ Fresh PostgreSQL container ready for testing
🧪 Running all tests (unit + integration)...
```

**Scenario 2: Local PostgreSQL Running**
```bash
$ make test-all
🔍 Checking PostgreSQL status on port 5433...
✅ Local PostgreSQL detected on port 5433 - using existing instance
🧪 Running all tests (unit + integration)...
```

**Scenario 3: No PostgreSQL Running**
```bash
$ make test-all
🔍 Checking PostgreSQL status on port 5433...
🐘 No PostgreSQL found - creating fresh Docker container...
✅ Fresh PostgreSQL container ready for testing
🧪 Running all tests (unit + integration)...
```

### CI/CD Considerations

**GitHub Actions/CI Environment:**
```yaml
# .github/workflows/test.yml example
- name: Run tests with Docker PostgreSQL
  run: |
    # CI environments are clean, so Docker will be created fresh
    make test-all
    # ✅ Optimal: No conflicts, fresh state, isolated tests
```

**Benefits for CI/CD:**
- **No manual setup**: Automatic PostgreSQL management
- **Fresh state**: Each CI run gets clean database
- **Port conflict free**: Uses non-standard port (5433)
- **Resource efficient**: Container recreation instead of data cleanup

### Manual Docker Operations

```bash
# Manual container management
docker run --name billing-postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=billing_service_dev \
  -p 5433:5432 -d postgres:15

# Create test database manually
docker exec billing-postgres psql -U postgres \
  -c "CREATE DATABASE billing_service_test;"

# Check container logs
docker logs billing-postgres

# Connect to PostgreSQL directly
docker exec -it billing-postgres psql -U postgres -d billing_service_dev
```

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