# NestGo Framework - Complete Documentation

---

## Table of Contents

### Getting Started
1. [Introduction](#introduction)
2. [Installation](#installation)
3. [Quick Start](#quick-start)
4. [Project Structure](#project-structure)
5. [Your First Application](#your-first-application)

### Core Concepts
6. [Modules](#modules)
7. [Controllers](#controllers)
8. [Services & Providers](#services--providers)
9. [Dependency Injection](#dependency-injection)
10. [Middleware](#middleware)

### HTTP
11. [Request & Response](#request--response)
12. [Routing](#routing)
13. [Request Validation](#request-validation)
14. [Response Formatting](#response-formatting)
15. [Exception Handling](#exception-handling)

### Advanced Features
16. [Guards](#guards)
17. [Interceptors](#interceptors)
18. [Pipes](#pipes)
19. [Exception Filters](#exception-filters)

### Data Management
20. [Database Integration](#database-integration)
21. [ORMs & Query Builders](#orms--query-builders)
22. [Migrations](#migrations)
23. [Transactions](#transactions)

### Real-time & Events
24. [WebSockets](#websockets)
25. [Event Emitters](#event-emitters)
26. [Message Queues](#message-queues)

### Background Processing
27. [Task Scheduling](#task-scheduling)
28. [Background Jobs](#background-jobs)
29. [Cron Jobs](#cron-jobs)

### Security
30. [Authentication](#authentication)
31. [Authorization (RBAC)](#authorization-rbac)
32. [CORS](#cors)
33. [Rate Limiting](#rate-limiting)
34. [Helmet & Security Headers](#helmet--security-headers)

### Testing
35. [Unit Testing](#unit-testing)
36. [Integration Testing](#integration-testing)
37. [E2E Testing](#e2e-testing)
38. [Mocking Dependencies](#mocking-dependencies)

### CLI Tool
39. [CLI Overview](#cli-overview)
40. [Commands Reference](#commands-reference)
41. [Generators](#generators)
42. [Custom Templates](#custom-templates)

### Deployment
43. [Production Build](#production-build)
44. [Docker](#docker)
45. [Kubernetes](#kubernetes)
46. [Environment Configuration](#environment-configuration)

### Monitoring & Observability
47. [Logging](#logging)
48. [Metrics (Prometheus)](#metrics-prometheus)
49. [Tracing (OpenTelemetry)](#tracing-opentelemetry)
50. [Health Checks](#health-checks)

### Advanced Topics
51. [Custom Decorators](#custom-decorators)
52. [Dynamic Modules](#dynamic-modules)
53. [Circular Dependencies](#circular-dependencies)
54. [Performance Optimization](#performance-optimization)

### API Reference
55. [Core API](#core-api)
56. [HTTP API](#http-api)
57. [DI Container API](#di-container-api)
58. [Common Utilities](#common-utilities)

---

# Getting Started

## Introduction

### What is NestGo?

NestGo is a production-grade backend framework for Go that brings the elegant architectural patterns of NestJS to the Go ecosystem. It provides:

- **Modular Architecture**: Organize code in reusable, testable modules
- **Dependency Injection**: Type-safe, compile-time validated dependencies
- **Enterprise Patterns**: Guards, Interceptors, Exception Filters, and more
- **Developer Experience**: Powerful CLI, hot-reloading, and scaffolding tools
- **Performance**: Native Go speed with zero runtime reflection

### Philosophy

NestGo follows three core principles:

1. **Explicit over Magic**: No runtime reflection or hidden behavior
2. **Type Safety**: Catch errors at compile-time, not runtime
3. **Convention over Configuration**: Opinionated structure for consistency

### When to Use NestGo?

**✅ Use NestGo when:**
- Building REST APIs with complex business logic
- Creating microservices architectures
- Need real-time features (WebSockets)
- Require background job processing
- Building SaaS or enterprise applications
- Team needs consistent code structure

**❌ Consider alternatives when:**
- Building simple static file servers
- Need maximum performance at the cost of structure
- Project is a small single-endpoint service

---

## Installation

### Prerequisites

- **Go**: Version 1.22 or higher
- **Git**: For version control
- **Text Editor**: VS Code, GoLand, or your preferred editor

### Installing the CLI

The NestGo CLI is the primary tool for creating and managing NestGo applications.

```bash
# Install globally
go install github.com/nestgo/nestgo/cmd/nestgo@latest

# Verify installation
nestgo --version
```

### Shell Autocompletion (Optional but Recommended)

Enable command autocompletion for faster development:

**Bash:**
```bash
# Add to ~/.bashrc
echo 'source <(nestgo completion bash)' >> ~/.bashrc
source ~/.bashrc
```

**Zsh:**
```bash
# Add to ~/.zshrc
echo 'source <(nestgo completion zsh)' >> ~/.zshrc
source ~/.zshrc
```

**Fish:**
```bash
nestgo completion fish | source
```

### Updating NestGo

To update to the latest version:

```bash
go install github.com/nestgo/nestgo/cmd/nestgo@latest
```

---

## Quick Start

### Creating Your First Project

```bash
# Create a new NestGo project
nestgo new my-api

# Navigate to project directory
cd my-api

# Start development server with hot-reload
nestgo dev
```

The server will start at `http://localhost:8080`.

### Project Creation Options

```bash
# Specify package manager
nestgo new my-api --package-manager=go-modules

# Skip Git initialization
nestgo new my-api --skip-git

# Use a specific template
nestgo new my-api --template=microservice
```

### Testing Your Application

```bash
# Visit in browser or use curl
curl http://localhost:8080/api/v1/health

# Expected response
{
  "status": "ok",
  "timestamp": "2026-04-18T10:30:00Z"
}
```

---

## Project Structure

### Default Directory Layout

```
my-api/
├── cmd/
│   └── nestgo/          # Application entry point
│       └── main.go
├── internal/
│   ├── app/             # Application module
│   │   ├── app.module.go
│   │   └── app.controller.go
│   ├── users/           # Users feature module
│   │   ├── users.module.go
│   │   ├── users.controller.go
│   │   ├── users.service.go
│   │   └── dto/
│   │       ├── create-user.dto.go
│   │       └── update-user.dto.go
│   └── common/          # Shared utilities
│       ├── guards/
│       ├── interceptors/
│       └── filters/
├── config/              # Configuration files
│   ├── config.go
│   └── config.yaml
├── pkg/                 # Reusable packages
├── test/                # E2E tests
├── .env                 # Environment variables
├── .gitignore
├── go.mod
├── go.sum
├── nestgo.yaml          # NestGo configuration
└── README.md
```

### Understanding the Structure

- **`cmd/`**: Application entry points (for multiple binaries)
- **`internal/`**: Private application code (cannot be imported by other projects)
- **`pkg/`**: Public libraries (can be imported by other projects)
- **`config/`**: Configuration management
- **`test/`**: End-to-end tests
- **`nestgo.yaml`**: Framework configuration

---

## Your First Application

### Step 1: Create a Module

```bash
nestgo generate module products
```

This creates:
```
internal/products/
└── products.module.go
```

**products.module.go:**
```go
package products

import (
    "github.com/nestgo/nestgo/common"
    "github.com/nestgo/nestgo/di"
)

type ProductsModule struct{}

func (m *ProductsModule) Module() common.ModuleConfig {
    return common.ModuleConfig{
        Name:        "products",
        Controllers: []common.Controller{},
        Providers:   []di.Provider{},
        Imports:     []common.Module{},
        Exports:     []di.Provider{},
    }
}
```

### Step 2: Create a Service

```bash
nestgo generate service products
```

**products.service.go:**
```go
package products

import "context"

type Product struct {
    ID    string `json:"id"`
    Name  string `json:"name"`
    Price float64 `json:"price"`
}

type ProductsService struct {
    // Dependencies will be injected here
}

func NewProductsService() *ProductsService {
    return &ProductsService{}
}

func (s *ProductsService) FindAll(ctx context.Context) ([]*Product, error) {
    // Business logic here
    return []*Product{
        {ID: "1", Name: "Product 1", Price: 99.99},
        {ID: "2", Name: "Product 2", Price: 149.99},
    }, nil
}

func (s *ProductsService) FindOne(ctx context.Context, id string) (*Product, error) {
    // Find product by ID
    return &Product{ID: id, Name: "Product 1", Price: 99.99}, nil
}

func (s *ProductsService) Create(ctx context.Context, product *Product) (*Product, error) {
    // Create new product
    return product, nil
}
```

### Step 3: Create a Controller

```bash
nestgo generate controller products
```

**products.controller.go:**
```go
package products

import (
    "github.com/nestgo/nestgo/common"
    "net/http"
)

type ProductsController struct {
    service *ProductsService
}

func NewProductsController(service *ProductsService) *ProductsController {
    return &ProductsController{service: service}
}

func (c *ProductsController) Prefix() string {
    return "/products"
}

func (c *ProductsController) Routes() []common.Route {
    return []common.Route{
        {
            Method:  http.MethodGet,
            Path:    "/",
            Handler: c.FindAll,
        },
        {
            Method:  http.MethodGet,
            Path:    "/{id}",
            Handler: c.FindOne,
        },
        {
            Method:  http.MethodPost,
            Path:    "/",
            Handler: c.Create,
        },
    }
}

func (c *ProductsController) FindAll(ctx *common.Context) error {
    products, err := c.service.FindAll(ctx.Request.Context())
    if err != nil {
        return err
    }
    return ctx.JSON(http.StatusOK, products)
}

func (c *ProductsController) FindOne(ctx *common.Context) error {
    id := ctx.Param("id")
    product, err := c.service.FindOne(ctx.Request.Context(), id)
    if err != nil {
        return err
    }
    return ctx.JSON(http.StatusOK, product)
}

func (c *ProductsController) Create(ctx *common.Context) error {
    var product Product
    if err := ctx.Bind(&product); err != nil {
        return err
    }
    
    created, err := c.service.Create(ctx.Request.Context(), &product)
    if err != nil {
        return err
    }
    
    return ctx.JSON(http.StatusCreated, created)
}
```

### Step 4: Register in Module

Update **products.module.go:**
```go
func (m *ProductsModule) Module() common.ModuleConfig {
    service := NewProductsService()
    controller := NewProductsController(service)

    return common.ModuleConfig{
        Name:        "products",
        Controllers: []common.Controller{controller},
        Providers: []di.Provider{
            {Instance: service},
        },
    }
}
```

### Step 5: Import into App Module

**app.module.go:**
```go
package app

import (
    "github.com/nestgo/nestgo/common"
    "my-api/internal/products"
)

type AppModule struct{}

func (m *AppModule) Module() common.ModuleConfig {
    return common.ModuleConfig{
        Name: "app",
        Imports: []common.Module{
            &products.ProductsModule{},
        },
    }
}
```

### Step 6: Test Your API

```bash
# Start server
nestgo dev

# Test endpoints
curl http://localhost:8080/api/v1/products
curl http://localhost:8080/api/v1/products/1
curl -X POST http://localhost:8080/api/v1/products \
  -H "Content-Type: application/json" \
  -d '{"name":"New Product","price":199.99}'
```

---

# Core Concepts

## Modules

### What are Modules?

Modules are the fundamental building blocks of NestGo applications. They organize code into cohesive, reusable units with clear boundaries.

### Module Structure

```go
type ModuleConfig struct {
    Name        string           // Unique module identifier
    Controllers []Controller     // HTTP controllers
    Providers   []Provider       // Injectable services
    Imports     []Module         // Other modules to import
    Exports     []Provider       // Providers to make available to other modules
}
```

### Creating a Module

**Simple Module:**
```go
package users

type UsersModule struct{}

func (m *UsersModule) Module() common.ModuleConfig {
    return common.ModuleConfig{
        Name: "users",
        Controllers: []common.Controller{
            NewUsersController(),
        },
        Providers: []di.Provider{
            {Instance: NewUsersService()},
        },
    }
}
```

**Module with Dependencies:**
```go
func (m *UsersModule) Module() common.ModuleConfig {
    // Create service instances
    emailService := NewEmailService()
    usersService := NewUsersService(emailService)
    controller := NewUsersController(usersService)

    return common.ModuleConfig{
        Name: "users",
        Controllers: []common.Controller{controller},
        Providers: []di.Provider{
            {Instance: emailService},
            {Instance: usersService},
        },
        Exports: []di.Provider{
            {Instance: usersService}, // Other modules can use this
        },
    }
}
```

### Module Imports

```go
type OrdersModule struct{}

func (m *OrdersModule) Module() common.ModuleConfig {
    return common.ModuleConfig{
        Name: "orders",
        Imports: []common.Module{
            &users.UsersModule{},    // Import users module
            &products.ProductsModule{}, // Import products module
        },
        // Now you can inject UsersService and ProductsService
    }
}
```

### Feature Modules

Organize by feature (recommended):
```
internal/
├── users/
│   ├── users.module.go
│   ├── users.controller.go
│   ├── users.service.go
│   └── users.repository.go
├── products/
│   ├── products.module.go
│   ├── products.controller.go
│   └── products.service.go
└── orders/
    ├── orders.module.go
    ├── orders.controller.go
    └── orders.service.go
```

### Shared Modules

Create reusable modules for cross-cutting concerns:

```go
package database

type DatabaseModule struct{}

func (m *DatabaseModule) Module() common.ModuleConfig {
    client := NewDatabaseClient()
    
    return common.ModuleConfig{
        Name: "database",
        Providers: []di.Provider{
            {Instance: client, Scope: di.Singleton},
        },
        Exports: []di.Provider{
            {Instance: client},
        },
    }
}
```

### Global Modules

Modules available everywhere without explicit import:

```go
func (m *ConfigModule) Module() common.ModuleConfig {
    return common.ModuleConfig{
        Name:   "config",
        Global: true, // Available to all modules
        Providers: []di.Provider{
            {Instance: NewConfigService()},
        },
    }
}
```

---

## Controllers

### What are Controllers?

Controllers handle incoming HTTP requests and return responses. They define routes and delegate business logic to services.

### Controller Interface

Every controller must implement:
```go
type Controller interface {
    Prefix() string              // Base path for all routes
    Routes() []Route             // Route definitions
}
```

### Basic Controller

```go
package users

import (
    "github.com/nestgo/nestgo/common"
    "net/http"
)

type UsersController struct {
    service *UsersService
}

func NewUsersController(service *UsersService) *UsersController {
    return &UsersController{service: service}
}

func (c *UsersController) Prefix() string {
    return "/users"
}

func (c *UsersController) Routes() []common.Route {
    return []common.Route{
        {Method: http.MethodGet, Path: "/", Handler: c.FindAll},
        {Method: http.MethodGet, Path: "/{id}", Handler: c.FindOne},
        {Method: http.MethodPost, Path: "/", Handler: c.Create},
        {Method: http.MethodPut, Path: "/{id}", Handler: c.Update},
        {Method: http.MethodDelete, Path: "/{id}", Handler: c.Delete},
    }
}
```

### Route Handlers

```go
// GET /users
func (c *UsersController) FindAll(ctx *common.Context) error {
    users, err := c.service.FindAll(ctx.Request.Context())
    if err != nil {
        return err
    }
    return ctx.JSON(http.StatusOK, users)
}

// GET /users/{id}
func (c *UsersController) FindOne(ctx *common.Context) error {
    id := ctx.Param("id")
    user, err := c.service.FindOne(ctx.Request.Context(), id)
    if err != nil {
        return err
    }
    return ctx.JSON(http.StatusOK, user)
}

// POST /users
func (c *UsersController) Create(ctx *common.Context) error {
    var dto CreateUserDTO
    if err := ctx.Bind(&dto); err != nil {
        return common.NewBadRequestError("Invalid request body")
    }
    
    user, err := c.service.Create(ctx.Request.Context(), &dto)
    if err != nil {
        return err
    }
    
    return ctx.JSON(http.StatusCreated, user)
}
```

### Request Parameters

**Path Parameters:**
```go
id := ctx.Param("id")          // /users/{id}
slug := ctx.Param("slug")      // /posts/{slug}
```

**Query Parameters:**
```go
page := ctx.Query("page")          // ?page=1
limit := ctx.Query("limit")        // ?limit=10
search := ctx.Query("search")      // ?search=john

// With defaults
page := ctx.QueryDefault("page", "1")
```

**Request Body:**
```go
var dto CreateUserDTO
if err := ctx.Bind(&dto); err != nil {
    return err
}
```

### Response Methods

```go
// JSON response
ctx.JSON(http.StatusOK, data)

// String response
ctx.String(http.StatusOK, "Hello World")

// HTML response
ctx.HTML(http.StatusOK, "<h1>Hello</h1>")

// No content
ctx.NoContent(http.StatusNoContent)

// Redirect
ctx.Redirect(http.StatusFound, "/login")

// File download
ctx.File("/path/to/file.pdf")

// Stream response
ctx.Stream(http.StatusOK, "text/csv", reader)
```

### Controller Middleware

Apply middleware to specific controllers:

```go
func (c *UsersController) Middlewares() []common.Middleware {
    return []common.Middleware{
        middleware.Auth(),
        middleware.RateLimit(100, time.Minute),
    }
}
```

### Grouped Routes

```go
func (c *AdminController) Routes() []common.Route {
    return []common.Route{
        // Public routes
        {Method: "GET", Path: "/health", Handler: c.Health},
        
        // Admin routes (with middleware)
        {
            Method: "GET",
            Path: "/users",
            Handler: c.ListUsers,
            Middlewares: []common.Middleware{middleware.AdminOnly()},
        },
        {
            Method: "DELETE",
            Path: "/users/{id}",
            Handler: c.DeleteUser,
            Middlewares: []common.Middleware{middleware.AdminOnly()},
        },
    }
}
```

---

## Services & Providers

### What are Services?

Services contain business logic and are injected into controllers via dependency injection.

### Creating a Service

```go
package users

type UsersService struct {
    db    *database.Client
    cache *cache.Client
}

// Constructor for DI
func NewUsersService(db *database.Client, cache *cache.Client) *UsersService {
    return &UsersService{
        db:    db,
        cache: cache,
    }
}

func (s *UsersService) FindAll(ctx context.Context) ([]*User, error) {
    // Check cache
    cached, err := s.cache.Get(ctx, "users:all")
    if err == nil {
        return cached.([]*User), nil
    }
    
    // Query database
    users, err := s.db.Query(ctx, "SELECT * FROM users")
    if err != nil {
        return nil, err
    }
    
    // Cache result
    s.cache.Set(ctx, "users:all", users, 5*time.Minute)
    
    return users, nil
}
```

### Service Patterns

**Repository Pattern:**
```go
type UserRepository struct {
    db *database.Client
}

func (r *UserRepository) FindByID(ctx context.Context, id string) (*User, error) {
    // Database query
}

func (r *UserRepository) Create(ctx context.Context, user *User) error {
    // Insert into database
}
```

**Business Logic Service:**
```go
type OrdersService struct {
    repo           *OrderRepository
    productsService *ProductsService
    emailService   *EmailService
}

func (s *OrdersService) PlaceOrder(ctx context.Context, dto *CreateOrderDTO) (*Order, error) {
    // Validate products
    for _, item := range dto.Items {
        product, err := s.productsService.FindOne(ctx, item.ProductID)
        if err != nil {
            return nil, err
        }
        if product.Stock < item.Quantity {
            return nil, errors.New("insufficient stock")
        }
    }
    
    // Create order
    order, err := s.repo.Create(ctx, dto)
    if err != nil {
        return nil, err
    }
    
    // Send confirmation email
    go s.emailService.SendOrderConfirmation(order)
    
    return order, nil
}
```

### Provider Scopes

**Singleton (Default):**
```go
{Instance: NewUsersService(), Scope: di.Singleton}
// Created once, shared across all requests
```

**Request:**
```go
{Instance: NewRequestLogger(), Scope: di.Request}
// New instance per HTTP request
```

### Interface-based Services

```go
// Define interface
type EmailSender interface {
    Send(to, subject, body string) error
}

// Implementation
type SMTPEmailService struct{}

func (s *SMTPEmailService) Send(to, subject, body string) error {
    // Send email via SMTP
}

// Register in module
Providers: []di.Provider{
    {Instance: &SMTPEmailService{}, As: (*EmailSender)(nil)},
}

// Inject interface
type NotificationService struct {
    email EmailSender
}
```

---

## Dependency Injection

### How DI Works in NestGo

NestGo uses **constructor-based dependency injection** without runtime reflection.

### Constructor Injection

```go
type OrdersService struct {
    db           *database.Client
    products     *ProductsService
    email        *EmailService
}

// Dependencies are injected via constructor
func NewOrdersService(
    db *database.Client,
    products *ProductsService,
    email *EmailService,
) *OrdersService {
    return &OrdersService{
        db:       db,
        products: products,
        email:    email,
    }
}
```

### Registering Providers

```go
func (m *OrdersModule) Module() common.ModuleConfig {
    // Create instances with dependencies
    db := NewDatabaseClient()
    productsService := NewProductsService(db)
    emailService := NewEmailService()
    ordersService := NewOrdersService(db, productsService, emailService)
    
    return common.ModuleConfig{
        Name: "orders",
        Providers: []di.Provider{
            {Instance: db},
            {Instance: productsService},
            {Instance: emailService},
            {Instance: ordersService},
        },
    }
}
```

### Dependency Resolution

The DI container resolves dependencies in order:

1. **Topological Sort**: Analyzes dependency graph
2. **Cycle Detection**: Fails at compile-time if circular dependencies exist
3. **Instantiation**: Creates instances in correct order
4. **Injection**: Passes dependencies to constructors

### Optional Dependencies

```go
type CacheService struct {
    redis *redis.Client
}

func NewCacheService(redis *redis.Client) *CacheService {
    // redis can be nil
    return &CacheService{redis: redis}
}

func (s *CacheService) Get(key string) (interface{}, error) {
    if s.redis == nil {
        return nil, errors.New("cache not configured")
    }
    // Use redis
}
```

### Provider Factories

```go
// Factory function
func NewDatabaseClient(config *Config) *database.Client {
    return database.Connect(config.DatabaseURL)
}

// Register with factory
Providers: []di.Provider{
    {
        Factory: func() interface{} {
            config := GetConfig()
            return NewDatabaseClient(config)
        },
    },
}
```

### Async Providers

```go
type AsyncProvider struct {
    Name    string
    Factory func() (interface{}, error)
}

// Example: Database connection
{
    Name: "database",
    Factory: func() (interface{}, error) {
        db, err := sql.Open("postgres", connString)
        if err != nil {
            return nil, err
        }
        
        // Wait for connection
        if err := db.Ping(); err != nil {
            return nil, err
        }
        
        return db, nil
    },
}
```

---

## Middleware

### What is Middleware?

Middleware functions execute before route handlers, allowing you to:
- Authenticate requests
- Log request/response
- Modify request/response
- Handle CORS
- Rate limiting

### Creating Middleware

```go
package middleware

import "github.com/nestgo/nestgo/common"

func Logger() common.Middleware {
    return func(next common.HandlerFunc) common.HandlerFunc {
        return func(ctx *common.Context) error {
            start := time.Now()
            
            // Before request
            log.Printf("Request: %s %s", ctx.Request.Method, ctx.Request.URL.Path)
            
            // Execute handler
            err := next(ctx)
            
            // After request
            duration := time.Since(start)
            log.Printf("Response: %d (took %v)", ctx.Response.Status, duration)
            
            return err
        }
    }
}
```

### Applying Middleware

**Global Middleware:**
```go
func main() {
    app := core.New()
    
    // Apply to all routes
    app.Use(middleware.Logger())
    app.Use(middleware.CORS())
    app.Use(middleware.Recover())
    
    app.Start()
}
```

**Module-level:**
```go
func (m *UsersModule) Middlewares() []common.Middleware {
    return []common.Middleware{
        middleware.Auth(),
        middleware.RateLimit(100, time.Minute),
    }
}
```

**Route-level:**
```go
{
    Method: "POST",
    Path: "/upload",
    Handler: c.Upload,
    Middlewares: []common.Middleware{
        middleware.MaxBodySize(10 * 1024 * 1024), // 10MB
    },
}
```

### Built-in Middleware

**CORS:**
```go
middleware.CORS(middleware.CORSConfig{
    AllowOrigins: []string{"https://example.com"},
    AllowMethods: []string{"GET", "POST", "PUT", "DELETE"},
    AllowHeaders: []string{"Content-Type", "Authorization"},
})
```

**Rate Limiting:**
```go
middleware.RateLimit(
    100,           // 100 requests
    time.Minute,   // per minute
)
```

**Request ID:**
```go
middleware.RequestID()
```

**Timeout:**
```go
middleware.Timeout(30 * time.Second)
```

### Custom Middleware Examples

**Authentication:**
```go
func Auth() common.Middleware {
    return func(next common.HandlerFunc) common.HandlerFunc {
        return func(ctx *common.Context) error {
            token := ctx.Header("Authorization")
            if token == "" {
                return common.NewUnauthorizedError("Missing token")
            }
            
            user, err := validateToken(token)
            if err != nil {
                return common.NewUnauthorizedError("Invalid token")
            }
            
            // Add user to context
            ctx.Set("user", user)
            
            return next(ctx)
        }
    }
}
```

**Request Logging:**
```go
func RequestLogger(logger *slog.Logger) common.Middleware {
    return func(next common.HandlerFunc) common.HandlerFunc {
        return func(ctx *common.Context) error {
            start := time.Now()
            
            logger.Info("Request started",
                "method", ctx.Request.Method,
                "path", ctx.Request.URL.Path,
                "ip", ctx.Request.RemoteAddr,
            )
            
            err := next(ctx)
            
            logger.Info("Request completed",
                "method", ctx.Request.Method,
                "path", ctx.Request.URL.Path,
                "status", ctx.Response.Status,
                "duration", time.Since(start),
            )
            
            return err
        }
    }
}
```

---

# HTTP

## Request & Response

### Context Object

The `Context` provides access to request and response:

```go
type Context struct {
    Request  *http.Request
    Response *ResponseWriter
    Params   map[string]string
    // ... other fields
}
```

### Reading Request Data

**Headers:**
```go
contentType := ctx.Header("Content-Type")
auth := ctx.Header("Authorization")

// Set response header
ctx.SetHeader("X-Custom-Header", "value")
```

**Cookies:**
```go
cookie, err := ctx.Cookie("session_id")
if err != nil {
    // Cookie not found
}

// Set cookie
ctx.SetCookie(&http.Cookie{
    Name:     "session_id",
    Value:    "abc123",
    MaxAge:   3600,
    HttpOnly: true,
    Secure:   true,
    SameSite: http.SameSiteStrictMode,
})
```

**Request Body:**
```go
// JSON
var dto CreateUserDTO
if err := ctx.Bind(&dto); err != nil {
    return common.NewBadRequestError("Invalid JSON")
}

// Form data
name := ctx.FormValue("name")
email := ctx.FormValue("email")

// Multipart form
file, header, err := ctx.FormFile("avatar")
if err != nil {
    return err
}
defer file.Close()
```

**Raw Body:**
```go
body, err := io.ReadAll(ctx.Request.Body)
if err != nil {
    return err
}
```

### Sending Responses

**JSON:**
```go
return ctx.JSON(http.StatusOK, map[string]interface{}{
    "message": "Success",
    "data": users,
})
```

**With Custom Headers:**
```go
ctx.SetHeader("X-Total-Count", "100")
ctx.SetHeader("Cache-Control", "max-age=3600")
return ctx.JSON(http.StatusOK, data)
```

**Paginated Response:**
```go
return ctx.Paginated(users, pagination.Meta{
    Page:       1,
    PerPage:    10,
    TotalPages: 5,
    TotalItems: 50,
})
```

**File Response:**
```go
// Download file
return ctx.File("/path/to/report.pdf")

// Stream file
return ctx.FileAttachment("/path/to/file.zip", "download.zip")
```

**Stream Response:**
```go
return ctx.Stream(http.StatusOK, "text/csv", reader)
```

### Status Codes

```go
// Success
ctx.JSON(http.StatusOK, data)              // 200
ctx.JSON(http.StatusCreated, data)         // 201
ctx.NoContent(http.StatusNoContent)        // 204

// Client Errors
return common.NewBadRequestError("...")    // 400
return common.NewUnauthorizedError("...")  // 401
return common.NewForbiddenError("...")     // 403
return common.NewNotFoundError("...")      // 404

// Server Errors
return common.NewInternalError("...")      // 500
```

---

## Routing

### Route Definition

```go
type Route struct {
    Method      string
    Path        string
    Handler     HandlerFunc
    Middlewares []Middleware
}
```

### Path Patterns

**Static:**
```go
{Method: "GET", Path: "/users", Handler: c.GetUsers}
```

**Parameters:**
```go
{Method: "GET", Path: "/users/{id}", Handler: c.GetUser}
{Method: "GET", Path: "/posts/{year}/{month}/{slug}", Handler: c.GetPost}
```

**Wildcard:**
```go
{Method: "GET", Path: "/files/*", Handler: c.ServeFiles}
```

**Query Strings:**
```go
// Route: /search
// URL: /search?q=golang&page=2
q := ctx.Query("q")
page := ctx.QueryDefault("page", "1")
```

### Route Groups

```go
func (c *APIController) Routes() []common.Route {
    v1Routes := []common.Route{
        {Method: "GET", Path: "/users", Handler: c.GetUsers},
        {Method: "GET", Path: "/posts", Handler: c.GetPosts},
    }
    
    v2Routes := []common.Route{
        {Method: "GET", Path: "/users", Handler: c.GetUsersV2},
        {Method: "GET", Path: "/posts", Handler: c.GetPostsV2},
    }
    
    return append(
        withPrefix("/v1", v1Routes),
        withPrefix("/v2", v2Routes)...,
    )
}
```

### Versioning

**URL Versioning:**
```go
app.RegisterModule(&v1.APIModuleV1{}, core.WithPrefix("/api/v1"))
app.RegisterModule(&v2.APIModuleV2{}, core.WithPrefix("/api/v2"))
```

**Header Versioning:**
```go
func (c *Controller) Routes() []common.Route {
    return []common.Route{
        {
            Method: "GET",
            Path: "/users",
            Handler: c.RouteByVersion,
        },
    }
}

func (c *Controller) RouteByVersion(ctx *common.Context) error {
    version := ctx.Header("API-Version")
    switch version {
    case "2":
        return c.GetUsersV2(ctx)
    default:
        return c.GetUsersV1(ctx)
    }
}
```

---

## Request Validation

### DTOs (Data Transfer Objects)

```go
package dto

type CreateUserDTO struct {
    Email     string `json:"email" validate:"required,email"`
    Password  string `json:"password" validate:"required,min=8"`
    FirstName string `json:"firstName" validate:"required"`
    LastName  string `json:"lastName" validate:"required"`
    Age       int    `json:"age" validate:"gte=18,lte=100"`
}

type UpdateUserDTO struct {
    FirstName *string `json:"firstName,omitempty"`
    LastName  *string `json:"lastName,omitempty"`
    Age       *int    `json:"age,omitempty" validate:"omitempty,gte=18"`
}
```

### Validation Tags

```go
`validate:"required"`              // Must be present
`validate:"email"`                 // Valid email format
`validate:"min=8,max=50"`          // String length
`validate:"gte=18,lte=100"`        // Number range
`validate:"url"`                   // Valid URL
`validate:"uuid"`                  // Valid UUID
`validate:"oneof=admin user"`      // Enum values
`validate:"omitempty"`             // Skip if empty
```

### Using Validation

```go
import "github.com/go-playground/validator/v10"

var validate = validator.New()

func (c *UsersController) Create(ctx *common.Context) error {
    var dto CreateUserDTO
    if err := ctx.Bind(&dto); err != nil {
        return common.NewBadRequestError("Invalid JSON")
    }
    
    // Validate
    if err := validate.Struct(dto); err != nil {
        return common.NewValidationError(err)
    }
    
    user, err := c.service.Create(ctx.Request.Context(), &dto)
    if err != nil {
        return err
    }
    
    return ctx.JSON(http.StatusCreated, user)
}
```

### Custom Validators

```go
func init() {
    validate.RegisterValidation("username", validateUsername)
}

func validateUsername(fl validator.FieldLevel) bool {
    username := fl.Field().String()
    // Only alphanumeric and underscore, 3-20 chars
    match, _ := regexp.MatchString(`^[a-zA-Z0-9_]{3,20}$`, username)
    return match
}

// Usage
type CreateUserDTO struct {
    Username string `json:"username" validate:"required,username"`
}
```

### Validation Middleware

```go
func ValidateDTO(dtoType interface{}) common.Middleware {
    return func(next common.HandlerFunc) common.HandlerFunc {
        return func(ctx *common.Context) error {
            dto := reflect.New(reflect.TypeOf(dtoType)).Interface()
            
            if err := ctx.Bind(dto); err != nil {
                return common.NewBadRequestError("Invalid request body")
            }
            
            if err := validate.Struct(dto); err != nil {
                return common.NewValidationError(err)
            }
            
            ctx.Set("dto", dto)
            return next(ctx)
        }
    }
}

// Usage
{
    Method: "POST",
    Path: "/users",
    Handler: c.Create,
    Middlewares: []common.Middleware{
        ValidateDTO(CreateUserDTO{}),
    },
}
```

---

## Exception Handling

### Built-in Exceptions

```go
// 400 Bad Request
common.NewBadRequestError("Invalid input")

// 401 Unauthorized
common.NewUnauthorizedError("Invalid credentials")

// 403 Forbidden
common.NewForbiddenError("Access denied")

// 404 Not Found
common.NewNotFoundError("User not found")

// 409 Conflict
common.NewConflictError("Email already exists")

// 422 Unprocessable Entity
common.NewValidationError(err)

// 500 Internal Server Error
common.NewInternalError("Database connection failed")
```

### Custom Exceptions

```go
type CustomError struct {
    Code    string
    Message string
    Status  int
}

func (e *CustomError) Error() string {
    return e.Message
}

func NewPaymentError(message string) error {
    return &CustomError{
        Code:    "PAYMENT_FAILED",
        Message: message,
        Status:  402,
    }
}
```

### Exception Filters

```go
package filters

import "github.com/nestgo/nestgo/common"

func GlobalExceptionFilter() common.ExceptionFilter {
    return func(ctx *common.Context, err error) {
        // Log error
        log.Error("Request failed", "error", err, "path", ctx.Request.URL.Path)
        
        // Handle different error types
        switch e := err.(type) {
        case *common.HTTPError:
            ctx.JSON(e.Status, map[string]interface{}{
                "error": e.Message,
                "code":  e.Code,
            })
        case *validator.ValidationErrors:
            ctx.JSON(http.StatusUnprocessableEntity, map[string]interface{}{
                "error": "Validation failed",
                "fields": formatValidationErrors(e),
            })
        default:
            // Hide internal errors from client
            ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
                "error": "Internal server error",
            })
        }
    }
}

// Apply globally
app.UseExceptionFilter(filters.GlobalExceptionFilter())
```

---

This comprehensive documentation continues with all 55 sections covering every aspect of NestGo. Would you like me to continue with the remaining sections?
