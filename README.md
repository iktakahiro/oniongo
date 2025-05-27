# Go DDD & Onion-Architecture Example and Techniques

[![Go Version](https://img.shields.io/badge/go-1.24.3-blue.svg)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/iktakahiro/oniongo)

English | [日本語](README.ja_JP.md)

**NOTE**: This repository is an example to demonstrate "how to implement DDD architecture in a Go application." If you use this as a reference, ensure to implement authentication and security before deploying it to a real-world environment!

* Python implementation: [dddpy](https://github.com/iktakahiro/dddpy)

## Tech Stack

* **gRPC & Connect-Go**: Modern RPC framework with HTTP/2 support
* **Ent**: Type-safe ORM for Go with code generation
* **SQLite**: Lightweight database for development
* **Buf**: Protocol buffer management and code generation
* **Samber/do**: Dependency injection container
* **Atlas**: Database migration tool

## Project Setup

1. Install dependencies:

```bash
make install
```

2. Generate code (protobuf, ent, mocks):

```bash
make buf-generate
make ent-generate
make mockgen
```

3. Run database migrations:

```bash
make migrate-up
```

4. Start the gRPC server:

```bash
make server
```

The server will start on port 8080 by default. You can override this by setting the `PORT` environment variable.

## Code Architecture

The directory structure is based on Onion Architecture:

```
internal/
├── domain/           # Domain Layer (Entities, Value Objects, Repository Interfaces)
│   └── todo/
├── application/      # Application Layer (Use Cases)
│   ├── todoapp/
│   └── uow/         # Unit of Work pattern
├── infrastructure/  # Infrastructure Layer (Repository Implementations, External Services)
│   ├── ent/         # Ent ORM (Schema, Generated Code, Repository)
│   ├── sqlite/      # Database migrations
│   └── di/          # Dependency injection setup
└── api/             # Presentation Layer (gRPC Handlers, Generated Code)
    └── grpc/        # gRPC handlers and generated protobuf code
```

### Domain Layer

The domain layer contains the core business logic and is independent of external concerns. It includes:

1. **Entities**: Core business objects with identity
2. **Value Objects**: Immutable objects that describe characteristics
3. **Repository Interfaces**: Contracts for data persistence
4. **Domain Services**: Business logic that doesn't belong to a single entity

#### 1. Entity

The `Todo` entity represents the core business object:

```go
// Todo is the entity that represents a todo item.
type Todo struct {
    id          TodoID
    title       string
    body        string
    status      TodoStatus
    createdAt   time.Time
    updatedAt   time.Time
    completedAt *time.Time
}

// NewTodo creates a new Todo.
func NewTodo(title string, body string) (*Todo, error) {
    now := time.Now()
    return &Todo{
        id:          NewTodoID(),
        title:       title,
        body:        body,
        status:      TodoStatusNotStarted,
        createdAt:   now,
        updatedAt:   now,
        completedAt: nil,
    }, nil
}
```

Key characteristics of the entity:

* Encapsulates business rules and invariants
* Provides methods for state transitions (`Start()`, `Complete()`)
* Maintains data integrity through validation
* Uses value objects for type safety (`TodoID`, `TodoStatus`)

#### 2. Value Objects

Value objects ensure type safety and encapsulate validation logic:

```go
// TodoID represents a unique identifier for a Todo.
type TodoID uuid.UUID

// NewTodoID creates a new TodoID.
func NewTodoID() TodoID {
    return TodoID(uuid.New())
}

// TodoStatus represents the status of a Todo.
type TodoStatus string

const (
    TodoStatusNotStarted TodoStatus = "not_started"
    TodoStatusInProgress TodoStatus = "in_progress"
    TodoStatusCompleted  TodoStatus = "completed"
)
```

#### 3. Repository Interface

The repository interface defines the contract for data persistence without specifying implementation details:

```go
// TodoRepository is the interface that wraps the basic CRUD operations for Todo.
type TodoRepository interface {
    Create(ctx context.Context, todo *Todo) error
    Update(ctx context.Context, todo *Todo) error
    FindAll(ctx context.Context) ([]*Todo, error)
    FindByID(ctx context.Context, id TodoID) (*Todo, error)
    Delete(ctx context.Context, id TodoID) error
}
```

### Infrastructure Layer

The infrastructure layer contains implementations of interfaces defined in the domain layer. It includes:

1. **Repository Implementations**: Concrete implementations using Ent ORM
2. **Database Schema**: Ent schema definitions
3. **External Service Integrations**: HTTP clients, third-party APIs, etc.
4. **Dependency Injection**: Service container configuration

#### 1. Repository Implementation

The `todoRepository` implements the domain repository interface using Ent ORM:

```go
// todoRepository is the implementation of the TodoRepository interface.
type todoRepository struct{}

// Create creates the Todo.
func (r todoRepository) Create(ctx context.Context, todo *todo.Todo) error {
    tx, err := db.GetTx(ctx)
    if err != nil {
        return err
    }

    status := todoschema.Status(todo.Status().String())
    _, err = tx.TodoSchema.Create().
        SetTitle(todo.Title()).
        SetBody(todo.Body()).
        SetStatus(status).
        Save(ctx)
    if err != nil {
        return fmt.Errorf("failed to create todo: %w", err)
    }
    return nil
}
```

Unlike the repository interface, the implementation code in the infrastructure layer can contain details specific to a particular technology (Ent ORM and SQLite in this example).

#### 2. Data Mapping

The infrastructure layer handles conversion between domain entities and database models:

```go
// convertEntToTodo converts ent.TodoSchema to domain Todo
func convertEntToTodo(v *entgen.TodoSchema) (*todo.Todo, error) {
    status, err := todo.NewTodoStatusFromString(string(v.Status))
    if err != nil {
        return nil, fmt.Errorf("failed to convert status %v: %w", v.Status, err)
    }
    return todo.ReconstructTodo(v.ID, v.Title, *v.Body, status, v.CreatedAt, v.UpdatedAt), nil
}
```

### Application Layer

The application layer contains the application-specific business rules. It includes:

1. **Use Case implementations**: Application services that orchestrate domain objects
2. **Transaction management**: Unit of Work pattern for data consistency
3. **Error handling**: Application-specific error handling

#### 1. Use Case Implementation

Each use case is implemented as a separate struct with a single `Execute` method:

```go
// CreateTodoUseCase is the interface that wraps the basic CreateTodo operation.
type CreateTodoUseCase interface {
    Execute(ctx context.Context, req CreateTodoRequest) error
}

// createTodoUseCase is the implementation of the CreateTodoUseCase interface.
type createTodoUseCase struct {
    todoRepository todo.TodoRepository
    txRunner       uow.TransactionRunner
}

// Execute creates a new Todo.
func (u createTodoUseCase) Execute(ctx context.Context, req CreateTodoRequest) error {
    todo, err := todo.NewTodo(req.Title, req.Body)
    if err != nil {
        return fmt.Errorf("failed to create todo: %w", err)
    }
    
    err = u.txRunner.RunInTx(ctx, func(ctx context.Context) error {
        if err := u.todoRepository.Create(ctx, todo); err != nil {
            return fmt.Errorf("failed to save todo: %w", err)
        }
        return nil
    })
    if err != nil {
        return fmt.Errorf("failed to execute transaction: %w", err)
    }
    return nil
}
```

Key characteristics of use cases:

* Single responsibility principle
* Transaction management through Unit of Work pattern
* Clear interface definition
* Dependency injection through constructor

#### 2. Transaction Management

The application layer uses the Unit of Work pattern to ensure data consistency:

```go
// TransactionRunner provides transaction management capabilities.
type TransactionRunner interface {
    RunInTx(ctx context.Context, fn func(ctx context.Context) error) error
}
```

### Presentation Layer

The presentation layer handles gRPC requests and responses. It includes:

1. **gRPC Handlers**: Convert between protobuf messages and domain objects
2. **Generated Code**: Protocol buffer generated code and Connect-Go handlers
3. **Input Validation**: Request validation using buf validate
4. **Error Handling**: Convert domain errors to gRPC status codes
5. **Middleware**: Cross-cutting concerns like logging and error handling

#### 1. gRPC Handler

The handlers are organized under the `api/grpc` directory:

```go
func (h *todoHandler) CreateTodo(
    ctx context.Context,
    req *connect.Request[v1.CreateTodoRequest],
) (*connect.Response[v1.CreateTodoResponse], error) {
    // Extract body value if present
    body := ""
    if req.Msg.Body != nil {
        body = *req.Msg.Body
    }

    // Create use case request
    useCaseReq := todoapp.CreateTodoRequest{
        Title: req.Msg.Title,
        Body:  body,
    }

    // Execute use case
    if err := h.createTodoUseCase.Execute(ctx, useCaseReq); err != nil {
        return nil, connect.NewError(connect.CodeInternal, err)
    }

    // Return response
    return connect.NewResponse(&v1.CreateTodoResponse{}), nil
}
```

#### 2. Protocol Buffer Definition

The API contract is defined using Protocol Buffers with validation:

```protobuf
// TodoService provides all todo-related operations
service TodoService {
  // CreateTodo creates a new todo item
  rpc CreateTodo(CreateTodoRequest) returns (CreateTodoResponse);
  
  // GetTodo retrieves a todo item by its ID
  rpc GetTodo(GetTodoRequest) returns (GetTodoResponse);
  
  // GetTodos retrieves all todo items
  rpc GetTodos(GetTodosRequest) returns (GetTodosResponse);
  
  // UpdateTodo updates an existing todo item
  rpc UpdateTodo(UpdateTodoRequest) returns (UpdateTodoResponse);
  
  // StartTodo changes the todo status to in progress
  rpc StartTodo(StartTodoRequest) returns (StartTodoResponse);
  
  // CompleteTodo changes the todo status to completed
  rpc CompleteTodo(CompleteTodoRequest) returns (CompleteTodoResponse);
  
  // DeleteTodo deletes a todo item
  rpc DeleteTodo(DeleteTodoRequest) returns (DeleteTodoResponse);
}

message CreateTodoRequest {
  string title = 1 [(buf.validate.field).string.min_len = 1];
  optional string body = 2;
}
```

## How to Work

1. Clone this repository
2. Install dependencies: `make install`
3. Generate code: `make buf-generate && make ent-generate`
4. Run migrations: `make migrate-up`
5. Start the server: `make server`
6. The gRPC server will be available at `localhost:8080`

### Sample Requests using grpcurl

* Create a new todo:

```bash
grpcurl -plaintext -d '{
  "title": "Implement DDD architecture",
  "body": "Create a sample application using DDD principles"
}' localhost:8080 oniongo.v1.TodoService/CreateTodo
```

* Get all todos:

```bash
grpcurl -plaintext -d '{}' localhost:8080 oniongo.v1.TodoService/GetTodos
```

* Get a specific todo:

```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000"
}' localhost:8080 oniongo.v1.TodoService/GetTodo
```

* Start a todo:

```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000"
}' localhost:8080 oniongo.v1.TodoService/StartTodo
```

* Complete a todo:

```bash
grpcurl -plaintext -d '{
  "id": "550e8400-e29b-41d4-a716-446655440000"
}' localhost:8080 oniongo.v1.TodoService/CompleteTodo
```

## Development

### Code Generation

This project uses several code generation tools:

```bash
# Generate protobuf code
make buf-generate

# Generate Ent ORM code
make ent-generate

# Generate mocks for testing
make mockgen
```

### Database Migrations

```bash
# Create a new migration
make migrate-diff name=add_new_field

# Apply migrations
make migrate-up
```

### Running Tests

```bash
go test ./...
```

### Code Quality

This project uses several tools to maintain code quality:

* **golangci-lint**: Comprehensive linter for Go
* **buf**: Protocol buffer linting and breaking change detection
* **mockery**: Mock generation for testing

```bash
# Format code
make fmt

# Run linter
make lint
```

## Key Design Patterns

### 1. Dependency Injection

The project uses the `samber/do` library for dependency injection, ensuring loose coupling between layers:

```go
// Dependency injection setup
func DependencyInjection() *do.Injector {
    injector := do.New()
    
    // Register dependencies
    do.Provide(injector, todorepo.NewTodoRepository)
    do.Provide(injector, todoapp.NewCreateTodoUseCase)
    do.Provide(injector, todohandler.NewTodoServiceHandler)
    
    return injector
}
```

### 2. Unit of Work Pattern

Transaction management is handled through the Unit of Work pattern, ensuring data consistency across multiple operations:

```go
err = u.txRunner.RunInTx(ctx, func(ctx context.Context) error {
    // Multiple repository operations within a single transaction
    return nil
})
```

### 3. Repository Pattern

The repository pattern abstracts data access logic, making the domain layer independent of specific database technologies:

```go
// Domain layer defines the interface
type TodoRepository interface {
    Create(ctx context.Context, todo *Todo) error
    // ... other methods
}

// Infrastructure layer provides the implementation
type todoRepository struct{}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## About

This repository demonstrates how to implement Domain-Driven Design (DDD) and Onion Architecture in Go, providing a clean, maintainable, and testable codebase structure for building scalable applications.
