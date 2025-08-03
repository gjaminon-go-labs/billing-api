# DDD/Clean Architecture Layers Explained

## **Domain-Driven Design Layers (Inside-Out)**

### **1. Domain Layer** 🎯 **(Core Business Logic)**
The heart of your application - pure business logic with no external dependencies.

- **Entities** - Objects with identity and lifecycle (e.g., Client with ID)
- **Value Objects** - Immutable objects without identity (e.g., Email, Phone)
- **Aggregates** - Consistency boundaries around entities
- **Domain Services** - Pure business logic that doesn't fit in entities
- **Repository Interfaces** - Contracts for data access (no implementation)

**Example:** `Client` entity with `Email` value object - handles business rules and validation

### **2. Application Layer** ⚙️ **(Use Case Orchestration)**
Coordinates domain objects to fulfill specific business use cases.

- **Application Services** - Orchestrate domain operations for specific use cases
- **Commands/Queries** - Input models for operations (CQRS pattern)
- **DTOs** - Data transfer objects for external communication
- **Use Case Handlers** - Coordinate domain objects to fulfill business scenarios

**Example:** `ClientService.CreateClient()` orchestrates Client creation using domain rules

### **3. Infrastructure Layer** 🔧 **(External Concerns)**
Handles all external dependencies and technical implementation details.

- **Repository Implementations** - Actual data access (database, API calls)
- **External Services** - Email, SMS, payment gateways
- **Configuration** - App settings, connection strings
- **Logging, Caching** - Cross-cutting concerns

**Example:** PostgreSQL repository implementing domain repository interface

### **4. Interface/Presentation Layer** 🌐 **(External Communication)**
Handles communication with external world (users, other systems).

- **HTTP Handlers** - REST API endpoints
- **Controllers** - Route requests to application services
- **Middleware** - Authentication, validation, error handling
- **DTOs/ViewModels** - Request/response models

**Example:** REST API handlers that accept HTTP requests and call application services

## **Current Status in Our Project**

### **✅ Implemented Layers:**

**Domain Layer (Complete):**
```
internal/domain/
├── entity/client.go          # Client aggregate root
└── valueobject/
    ├── email.go             # Email value object
    └── phone.go             # Phone value object
```

**Application Layer (Phase 1 Complete):**
```
internal/application/
└── client_service.go        # ClientService for use case orchestration
```

### **❌ Missing Layers:**

**Infrastructure Layer (Not implemented):**
- Repository implementations
- Database connections
- External service integrations

**Interface Layer (Not implemented):**
- HTTP handlers
- REST API endpoints
- Request/response DTOs

## **Dependencies Flow (Clean Architecture)**

```
Interface Layer → Application Layer → Domain Layer ← Infrastructure Layer
     🌐              ⚙️               🎯                🔧
```

**Key Rules:**
- **Inner layers never depend on outer layers**
- **Domain layer has no external dependencies** 
- **Infrastructure layer implements domain interfaces**
- **Interface layer only knows about application layer**

## **Enterprise Benefits**

### **Testability:**
- Domain: Pure unit tests (no mocks needed)
- Application: Service tests with domain objects
- Infrastructure: Integration tests with real databases
- Interface: API tests with mocked application services

### **Maintainability:**
- Business logic isolated in domain
- Use cases clearly defined in application layer
- Technical concerns separated in infrastructure
- External interfaces decoupled

### **Flexibility:**
- Can change databases without touching domain
- Can add new interfaces (GraphQL, gRPC) easily
- Can modify business rules without breaking infrastructure
- Can test business logic independently

## **Next Steps**

### **Phase 2: HTTP Interface Layer**
Add REST API endpoints that use the application service:
- HTTP handlers for Client operations
- Request/response DTOs
- Route configuration

### **Phase 3: Infrastructure Layer**
Add persistence capabilities:
- Repository interface in domain
- In-memory repository implementation
- Database integration (PostgreSQL)

### **Phase 4: Advanced Patterns**
Enterprise-grade features:
- CQRS (Command Query Responsibility Segregation)
- Event sourcing
- Domain events
- Saga patterns

## **Learning Objectives Achieved**

✅ **Enterprise Architecture:** Proper layered design  
✅ **Domain-Driven Design:** Entities, value objects, aggregates  
✅ **Clean Architecture:** Dependency inversion, separation of concerns  
✅ **Application Services:** Use case orchestration patterns  
✅ **Value Objects:** Immutable domain concepts with behavior  
✅ **Factory Methods:** Controlled object creation with validation  

---

*This document serves as a learning reference for understanding enterprise software architecture patterns implemented in the GoTuto billing service.*