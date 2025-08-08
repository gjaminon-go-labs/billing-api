# Billing API

Billing microservice for the go-labs platform handling client management, invoicing, and payment operations.

## 🚀 Quick Start

```bash
# Install dependencies
make restore

# Setup database (required for development)
make migrate-up

# Run development server (port 8080)
make run-dev

# Run all tests
make test-all
```

## 📡 API Endpoints

### Client Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/clients` | Create new client |
| GET | `/api/v1/clients/:id` | Get client by ID |
| PUT | `/api/v1/clients/:id` | Update client |
| DELETE | `/api/v1/clients/:id` | Delete client |
| GET | `/api/v1/clients` | List all clients (coming soon) |

### System
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/health` | Health check with database status |

## 🧪 Testing

```bash
# Unit tests only (fast, in-memory)
make test-unit

# Integration tests only (requires PostgreSQL)
make test-integration

# All tests with coverage
make test-all

# Generate business coverage report
make test-integration-report
```

## 🛠️ Development Commands

```bash
# Database migrations
make migrate-up          # Apply migrations
make migrate-down        # Rollback one migration
make migrate-status      # Check migration status

# Development
make run-dev            # Start with hot reload
make build              # Build binary
make clean              # Clean build artifacts

# Code quality
make lint               # Run linter
make fmt                # Format code
```

## 🏗️ Architecture

This service implements:
- **Clean Architecture** with clear layer separation
- **Domain-Driven Design** principles
- **Repository Pattern** for data access
- **PostgreSQL** with GORM and golang-migrate
- **Comprehensive testing** (unit + integration)

For detailed architecture documentation, see [platform-docs](../../platform-docs/ARCHITECTURE.md).

## 📋 Prerequisites

- Go 1.21+
- PostgreSQL 12+
- Make

## 🔗 Related Documentation

- [Platform Architecture](../../platform-docs/ARCHITECTURE.md) - Complete system design
- [TDD Workflow](../../CLAUDE.md) - Development methodology
- [API Documentation](docs/API.md) - Detailed endpoint specs (coming soon)

## 🤝 Contributing

All contributions must follow the TDD workflow defined in the root `CLAUDE.md`. 
No exceptions - tests first, always!

---

Part of [go-labs](https://github.com/gjaminon-go-labs) - A comprehensive Go microservices platform