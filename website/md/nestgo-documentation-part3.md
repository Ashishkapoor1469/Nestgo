# NestGo Framework - Complete Documentation (Part 3)

---

# Security

## Authentication

### JWT Authentication

```go
package auth

import (
    "github.com/golang-jwt/jwt/v5"
    "time"
)

type JWTService struct {
    secretKey []byte
}

func NewJWTService(secret string) *JWTService {
    return &JWTService{secretKey: []byte(secret)}
}

type Claims struct {
    UserID string `json:"userId"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func (s *JWTService) GenerateToken(user *User) (string, error) {
    claims := &Claims{
        UserID: user.ID,
        Email:  user.Email,
        Role:   user.Role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "nestgo-app",
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(s.secretKey)
}

func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return s.secretKey, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        return claims, nil
    }
    
    return nil, errors.New("invalid token")
}
```

### Auth Controller

```go
type AuthController struct {
    authService *AuthService
    jwtService  *JWTService
}

func (c *AuthController) Login(ctx *common.Context) error {
    var dto LoginDTO
    if err := ctx.Bind(&dto); err != nil {
        return common.NewBadRequestError("Invalid request")
    }
    
    user, err := c.authService.ValidateCredentials(dto.Email, dto.Password)
    if err != nil {
        return common.NewUnauthorizedError("Invalid credentials")
    }
    
    token, err := c.jwtService.GenerateToken(user)
    if err != nil {
        return err
    }
    
    return ctx.JSON(http.StatusOK, map[string]interface{}{
        "token": token,
        "user":  user,
    })
}

func (c *AuthController) Register(ctx *common.Context) error {
    var dto RegisterDTO
    if err := ctx.Bind(&dto); err != nil {
        return common.NewBadRequestError("Invalid request")
    }
    
    user, err := c.authService.Register(&dto)
    if err != nil {
        return err
    }
    
    token, err := c.jwtService.GenerateToken(user)
    if err != nil {
        return err
    }
    
    return ctx.JSON(http.StatusCreated, map[string]interface{}{
        "token": token,
        "user":  user,
    })
}

func (c *AuthController) Me(ctx *common.Context) error {
    user := ctx.Get("user").(*User)
    return ctx.JSON(http.StatusOK, user)
}
```

### Password Hashing

```go
package auth

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### OAuth2 Integration

```go
package auth

import (
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

type OAuth2Service struct {
    config *oauth2.Config
}

func NewOAuth2Service(clientID, clientSecret, redirectURL string) *OAuth2Service {
    return &OAuth2Service{
        config: &oauth2.Config{
            ClientID:     clientID,
            ClientSecret: clientSecret,
            RedirectURL:  redirectURL,
            Scopes:       []string{"email", "profile"},
            Endpoint:     google.Endpoint,
        },
    }
}

func (s *OAuth2Service) GetAuthURL(state string) string {
    return s.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *OAuth2Service) ExchangeCode(code string) (*oauth2.Token, error) {
    return s.config.Exchange(context.Background(), code)
}
```

---

## Authorization (RBAC)

### Permission-based Authorization

```go
package rbac

type Permission string

const (
    PermissionReadUser   Permission = "user:read"
    PermissionWriteUser  Permission = "user:write"
    PermissionDeleteUser Permission = "user:delete"
    PermissionReadOrder  Permission = "order:read"
    PermissionWriteOrder Permission = "order:write"
)

type Role struct {
    Name        string
    Permissions []Permission
}

var Roles = map[string]Role{
    "admin": {
        Name: "admin",
        Permissions: []Permission{
            PermissionReadUser,
            PermissionWriteUser,
            PermissionDeleteUser,
            PermissionReadOrder,
            PermissionWriteOrder,
        },
    },
    "user": {
        Name: "user",
        Permissions: []Permission{
            PermissionReadUser,
            PermissionReadOrder,
        },
    },
}

func HasPermission(role string, permission Permission) bool {
    r, exists := Roles[role]
    if !exists {
        return false
    }
    
    for _, p := range r.Permissions {
        if p == permission {
            return true
        }
    }
    return false
}
```

### Permission Guard

```go
type PermissionGuard struct {
    requiredPermission Permission
}

func NewPermissionGuard(permission Permission) *PermissionGuard {
    return &PermissionGuard{requiredPermission: permission}
}

func (g *PermissionGuard) CanActivate(ctx *common.Context) (bool, error) {
    user := ctx.Get("user").(*User)
    if user == nil {
        return false, common.NewUnauthorizedError("Not authenticated")
    }
    
    if !HasPermission(user.Role, g.requiredPermission) {
        return false, common.NewForbiddenError("Insufficient permissions")
    }
    
    return true, nil
}

// Usage
{
    Method: "DELETE",
    Path: "/users/{id}",
    Handler: c.Delete,
    Guards: []common.Guard{
        NewPermissionGuard(PermissionDeleteUser),
    },
}
```

---

## CORS

### CORS Configuration

```go
package middleware

import "github.com/rs/cors"

func CORS() common.Middleware {
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"http://localhost:3000", "https://example.com"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Content-Type", "Authorization"},
        ExposedHeaders:   []string{"X-Total-Count"},
        AllowCredentials: true,
        MaxAge:           300,
    })
    
    return func(next common.HandlerFunc) common.HandlerFunc {
        return func(ctx *common.Context) error {
            c.HandlerFunc(ctx.Response, ctx.Request)
            
            if ctx.Request.Method == "OPTIONS" {
                return ctx.NoContent(http.StatusNoContent)
            }
            
            return next(ctx)
        }
    }
}

// Apply globally
app.Use(middleware.CORS())
```

---

## Rate Limiting

### Token Bucket Rate Limiter

```go
package middleware

import (
    "golang.org/x/time/rate"
    "sync"
)

type RateLimiter struct {
    limiters map[string]*rate.Limiter
    mu       sync.RWMutex
    rate     rate.Limit
    burst    int
}

func NewRateLimiter(rps int, burst int) *RateLimiter {
    return &RateLimiter{
        limiters: make(map[string]*rate.Limiter),
        rate:     rate.Limit(rps),
        burst:    burst,
    }
}

func (rl *RateLimiter) getLimiter(key string) *rate.Limiter {
    rl.mu.RLock()
    limiter, exists := rl.limiters[key]
    rl.mu.RUnlock()
    
    if !exists {
        rl.mu.Lock()
        limiter = rate.NewLimiter(rl.rate, rl.burst)
        rl.limiters[key] = limiter
        rl.mu.Unlock()
    }
    
    return limiter
}

func (rl *RateLimiter) Allow(key string) bool {
    return rl.getLimiter(key).Allow()
}

func RateLimit(rps, burst int) common.Middleware {
    limiter := NewRateLimiter(rps, burst)
    
    return func(next common.HandlerFunc) common.HandlerFunc {
        return func(ctx *common.Context) error {
            ip := ctx.Request.RemoteAddr
            
            if !limiter.Allow(ip) {
                return common.NewTooManyRequestsError("Rate limit exceeded")
            }
            
            return next(ctx)
        }
    }
}

// Usage
app.Use(middleware.RateLimit(100, 200)) // 100 requests per second, burst 200
```

---

## Helmet & Security Headers

### Security Headers Middleware

```go
package middleware

func SecurityHeaders() common.Middleware {
    return func(next common.HandlerFunc) common.HandlerFunc {
        return func(ctx *common.Context) error {
            // Prevent XSS attacks
            ctx.SetHeader("X-Content-Type-Options", "nosniff")
            
            // Prevent clickjacking
            ctx.SetHeader("X-Frame-Options", "DENY")
            
            // Enable XSS filter
            ctx.SetHeader("X-XSS-Protection", "1; mode=block")
            
            // HSTS
            ctx.SetHeader("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
            
            // Content Security Policy
            ctx.SetHeader("Content-Security-Policy", "default-src 'self'")
            
            // Referrer Policy
            ctx.SetHeader("Referrer-Policy", "no-referrer")
            
            // Permissions Policy
            ctx.SetHeader("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
            
            return next(ctx)
        }
    }
}

app.Use(middleware.SecurityHeaders())
```

---

# Testing

## Unit Testing

### Service Testing

```go
package users_test

import (
    "context"
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

type MockRepository struct {
    mock.Mock
}

func (m *MockRepository) FindByID(ctx context.Context, id string) (*User, error) {
    args := m.Called(ctx, id)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).(*User), args.Error(1)
}

func TestUsersService_FindOne(t *testing.T) {
    // Arrange
    mockRepo := new(MockRepository)
    service := NewUsersService(mockRepo)
    
    expectedUser := &User{
        ID:    "123",
        Email: "test@example.com",
        Name:  "Test User",
    }
    
    mockRepo.On("FindByID", mock.Anything, "123").Return(expectedUser, nil)
    
    // Act
    user, err := service.FindOne(context.Background(), "123")
    
    // Assert
    assert.NoError(t, err)
    assert.Equal(t, expectedUser, user)
    mockRepo.AssertExpectations(t)
}

func TestUsersService_FindOne_NotFound(t *testing.T) {
    mockRepo := new(MockRepository)
    service := NewUsersService(mockRepo)
    
    mockRepo.On("FindByID", mock.Anything, "999").Return(nil, errors.New("not found"))
    
    user, err := service.FindOne(context.Background(), "999")
    
    assert.Error(t, err)
    assert.Nil(t, user)
}
```

### Testing with Table-Driven Tests

```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "user@example.com", false},
        {"missing @", "userexample.com", true},
        {"missing domain", "user@", true},
        {"empty string", "", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

---

## Integration Testing

### HTTP Integration Tests

```go
package integration_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
)

func TestUsersAPI(t *testing.T) {
    // Setup
    app := setupTestApp()
    
    t.Run("POST /users - Create User", func(t *testing.T) {
        dto := CreateUserDTO{
            Email:    "test@example.com",
            Password: "password123",
            Name:     "Test User",
        }
        
        body, _ := json.Marshal(dto)
        req := httptest.NewRequest("POST", "/api/v1/users", bytes.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        
        resp := httptest.NewRecorder()
        app.ServeHTTP(resp, req)
        
        assert.Equal(t, http.StatusCreated, resp.Code)
        
        var result map[string]interface{}
        json.NewDecoder(resp.Body).Decode(&result)
        assert.NotEmpty(t, result["id"])
    })
    
    t.Run("GET /users/{id} - Get User", func(t *testing.T) {
        req := httptest.NewRequest("GET", "/api/v1/users/123", nil)
        resp := httptest.NewRecorder()
        
        app.ServeHTTP(resp, req)
        
        assert.Equal(t, http.StatusOK, resp.Code)
    })
}

func setupTestApp() *core.Application {
    app := core.New(core.WithGlobalPrefix("/api/v1"))
    app.RegisterModule(&users.UsersModule{})
    return app
}
```

### Database Integration Tests

```go
func TestUsersRepository_Integration(t *testing.T) {
    // Setup test database
    db := setupTestDB(t)
    defer db.Close()
    
    repo := NewUsersRepository(db)
    
    t.Run("Create and Find User", func(t *testing.T) {
        user := &User{
            Email: "test@example.com",
            Name:  "Test User",
        }
        
        err := repo.Create(context.Background(), user)
        assert.NoError(t, err)
        assert.NotEmpty(t, user.ID)
        
        found, err := repo.FindByID(context.Background(), user.ID)
        assert.NoError(t, err)
        assert.Equal(t, user.Email, found.Email)
    })
}

func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("postgres", "postgres://test:test@localhost/testdb")
    if err != nil {
        t.Fatal(err)
    }
    
    // Run migrations
    runMigrations(db)
    
    // Cleanup on test end
    t.Cleanup(func() {
        db.Exec("TRUNCATE TABLE users CASCADE")
    })
    
    return db
}
```

---

## E2E Testing

### End-to-End Test Suite

```go
package e2e_test

import (
    "testing"
    "github.com/gavv/httpexpect/v2"
)

func TestUserWorkflow(t *testing.T) {
    // Start test server
    app := startTestServer(t)
    defer app.Stop()
    
    e := httpexpect.New(t, "http://localhost:8080")
    
    // Register
    registerResp := e.POST("/api/v1/auth/register").
        WithJSON(map[string]string{
            "email":    "test@example.com",
            "password": "password123",
            "name":     "Test User",
        }).
        Expect().
        Status(http.StatusCreated).
        JSON().Object()
    
    token := registerResp.Value("token").String().Raw()
    userID := registerResp.Value("user").Object().Value("id").String().Raw()
    
    // Get profile
    e.GET("/api/v1/users/{id}", userID).
        WithHeader("Authorization", "Bearer "+token).
        Expect().
        Status(http.StatusOK).
        JSON().Object().
        ValueEqual("email", "test@example.com")
    
    // Update profile
    e.PUT("/api/v1/users/{id}", userID).
        WithHeader("Authorization", "Bearer "+token).
        WithJSON(map[string]string{"name": "Updated Name"}).
        Expect().
        Status(http.StatusOK)
    
    // Verify update
    e.GET("/api/v1/users/{id}", userID).
        WithHeader("Authorization", "Bearer "+token).
        Expect().
        Status(http.StatusOK).
        JSON().Object().
        ValueEqual("name", "Updated Name")
}
```

---

## Mocking Dependencies

### Interface-based Mocking

```go
// Define interface
type EmailSender interface {
    Send(to, subject, body string) error
}

// Mock implementation
type MockEmailSender struct {
    mock.Mock
}

func (m *MockEmailSender) Send(to, subject, body string) error {
    args := m.Called(to, subject, body)
    return args.Error(0)
}

// Test
func TestNotificationService(t *testing.T) {
    mockEmail := new(MockEmailSender)
    service := NewNotificationService(mockEmail)
    
    mockEmail.On("Send", "user@example.com", "Welcome", mock.Anything).Return(nil)
    
    err := service.SendWelcome("user@example.com")
    
    assert.NoError(t, err)
    mockEmail.AssertExpectations(t)
}
```

---

# CLI Tool

## CLI Overview

The NestGo CLI is a powerful code generation and project management tool.

### Installation

```bash
go install github.com/nestgo/nestgo/cmd/nestgo@latest
```

### Available Commands

```bash
nestgo --help

Available Commands:
  new           Create a new NestGo project
  generate      Generate code (module, controller, service, etc.)
  dev           Start development server with hot-reload
  build         Build production binary
  test          Run tests
  migration     Database migration commands
  doctor        Check project health
  graph         Visualize dependency graph
  completion    Generate shell completion scripts
```

---

## Commands Reference

### new - Create New Project

```bash
# Basic usage
nestgo new my-api

# Options
nestgo new my-api --package-manager=go-modules
nestgo new my-api --skip-git
nestgo new my-api --template=microservice
nestgo new my-api --database=postgres
```

**Generated Structure:**
```
my-api/
├── cmd/nestgo/main.go
├── internal/
│   ├── app/
│   └── config/
├── pkg/
├── test/
├── .env
├── go.mod
└── nestgo.yaml
```

### generate - Code Generation

```bash
# Generate module
nestgo generate module users

# Generate controller
nestgo generate controller users

# Generate service
nestgo generate service users

# Generate full resource (module + controller + service + dto)
nestgo generate resource products

# Generate guard
nestgo generate guard auth

# Generate interceptor
nestgo generate interceptor logging

# Generate middleware
nestgo generate middleware cors
```

**Resource Generation:**
```bash
nestgo generate resource orders
? What transport layer would you like to use? REST API
? Would you like to generate CRUD entry points? Yes

CREATE internal/orders/orders.module.go
CREATE internal/orders/orders.controller.go
CREATE internal/orders/orders.service.go
CREATE internal/orders/orders.repository.go
CREATE internal/orders/dto/create-order.dto.go
CREATE internal/orders/dto/update-order.dto.go
CREATE internal/orders/entities/order.entity.go
CREATE internal/orders/orders_test.go
```

### dev - Development Server

```bash
# Start with hot-reload
nestgo dev

# Custom port
nestgo dev --port=3000

# Verbose logging
nestgo dev --verbose

# Skip build cache
nestgo dev --no-cache
```

### build - Production Build

```bash
# Build binary
nestgo build

# Custom output
nestgo build --output=./bin/app

# Cross-compile
nestgo build --os=linux --arch=amd64

# Optimized build
nestgo build --optimize
```

### test - Run Tests

```bash
# Run all tests
nestgo test

# Run specific package
nestgo test ./internal/users

# With coverage
nestgo test --coverage

# Watch mode
nestgo test --watch

# E2E tests only
nestgo test --e2e
```

### migration - Database Migrations

```bash
# Create migration
nestgo migration:create create_users_table

# Run migrations
nestgo migration:run

# Rollback last migration
nestgo migration:rollback

# Rollback all
nestgo migration:reset

# Migration status
nestgo migration:status
```

### doctor - Health Check

```bash
nestgo doctor

Checking NestGo project health...
✓ Go version: 1.22.1
✓ Dependencies: up to date
✓ Module structure: valid
✓ No circular dependencies detected
✓ Test coverage: 78%
⚠ Warning: Some files missing package comments

Recommendations:
  - Add package documentation to internal/users
  - Consider increasing test coverage to 80%
```

### graph - Dependency Visualization

```bash
# Generate dependency graph
nestgo graph

# Output to file
nestgo graph --output=graph.png

# Include only modules
nestgo graph --modules-only

# Detect cycles
nestgo graph --check-cycles
```

---

## Generators

### Custom Templates

Create custom templates in `.nestgo/templates/`:

```go
// .nestgo/templates/service.tmpl
package {{ .Package }}

import (
    "context"
    "github.com/nestgo/nestgo/common"
)

type {{ .Name }}Service struct {
    repository *{{ .Name }}Repository
}

func New{{ .Name }}Service(repository *{{ .Name }}Repository) *{{ .Name }}Service {
    return &{{ .Name }}Service{
        repository: repository,
    }
}

func (s *{{ .Name }}Service) FindAll(ctx context.Context) ([]*{{ .Name }}, error) {
    return s.repository.FindAll(ctx)
}
```

**Use custom template:**
```bash
nestgo generate service users --template=custom
```

---

# Deployment

## Production Build

### Build Configuration

```yaml
# nestgo.yaml
build:
  output: ./bin/app
  optimize: true
  ldflags:
    - -s -w
    - -X main.version={{.Version}}
    - -X main.buildTime={{.BuildTime}}
```

### Build Script

```bash
#!/bin/bash
# build.sh

VERSION=$(git describe --tags --always)
BUILD_TIME=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

go build \
  -ldflags="-s -w -X main.version=${VERSION} -X main.buildTime=${BUILD_TIME}" \
  -o ./bin/app \
  ./cmd/nestgo
```

---

## Docker

### Dockerfile

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /app/bin/nestgo \
    ./cmd/nestgo

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/bin/nestgo .

# Copy config files
COPY --from=builder /app/config ./config

EXPOSE 8080

CMD ["./nestgo"]
```

### Docker Compose

```yaml
version: '3.8'

services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:pass@db:5432/mydb
      - REDIS_URL=redis://redis:6379
    depends_on:
      - db
      - redis
    restart: unless-stopped

  db:
    image: postgres:15-alpine
    environment:
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=pass
      - POSTGRES_DB=mydb
    volumes:
      - postgres_data:/var/lib/postgresql/data
    ports:
      - "5432:5432"

  redis:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data

volumes:
  postgres_data:
  redis_data:
```

---

## Kubernetes

### Deployment

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nestgo-app
  labels:
    app: nestgo
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nestgo
  template:
    metadata:
      labels:
        app: nestgo
    spec:
      containers:
      - name: app
        image: your-registry/nestgo-app:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: app-secrets
              key: database-url
        - name: REDIS_URL
          value: redis://redis:6379
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "256Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

### Service

```yaml
# service.yaml
apiVersion: v1
kind: Service
metadata:
  name: nestgo-service
spec:
  selector:
    app: nestgo
  ports:
  - protocol: TCP
    port: 80
    targetPort: 8080
  type: LoadBalancer
```

### ConfigMap

```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: app-config
data:
  APP_ENV: production
  LOG_LEVEL: info
  PORT: "8080"
```

---

## Environment Configuration

### Config Structure

```go
package config

import (
    "github.com/spf13/viper"
)

type Config struct {
    App      AppConfig
    Database DatabaseConfig
    Redis    RedisConfig
    JWT      JWTConfig
}

type AppConfig struct {
    Name        string
    Environment string
    Port        int
    LogLevel    string
}

type DatabaseConfig struct {
    Host     string
    Port     int
    User     string
    Password string
    Database string
    SSLMode  string
}

func Load() (*Config, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./config")
    viper.AddConfigPath(".")
    
    viper.AutomaticEnv()
    
    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }
    
    var config Config
    if err := viper.Unmarshal(&config); err != nil {
        return nil, err
    }
    
    return &config, nil
}
```

### Environment Files

**.env:**
```bash
APP_ENV=development
APP_PORT=8080
DATABASE_URL=postgres://user:pass@localhost:5432/mydb
REDIS_URL=redis://localhost:6379
JWT_SECRET=your-secret-key
```

**config.yaml:**
```yaml
app:
  name: NestGo API
  environment: ${APP_ENV}
  port: ${APP_PORT}
  logLevel: info

database:
  host: localhost
  port: 5432
  user: ${DB_USER}
  password: ${DB_PASSWORD}
  database: ${DB_NAME}
  sslMode: disable

redis:
  url: ${REDIS_URL}

jwt:
  secret: ${JWT_SECRET}
  expiresIn: 24h
```

---

# Monitoring & Observability

## Logging

### Structured Logging with slog

```go
package logger

import (
    "log/slog"
    "os"
)

func New(level string) *slog.Logger {
    var logLevel slog.Level
    switch level {
    case "debug":
        logLevel = slog.LevelDebug
    case "info":
        logLevel = slog.LevelInfo
    case "warn":
        logLevel = slog.LevelWarn
    case "error":
        logLevel = slog.LevelError
    default:
        logLevel = slog.LevelInfo
    }
    
    opts := &slog.HandlerOptions{
        Level: logLevel,
    }
    
    handler := slog.NewJSONHandler(os.Stdout, opts)
    return slog.New(handler)
}

// Usage
logger := logger.New("info")
logger.Info("User created",
    "userId", user.ID,
    "email", user.Email,
)
logger.Error("Database error",
    "error", err,
    "query", query,
)
```

### Request Logging Middleware

```go
func RequestLogger(logger *slog.Logger) common.Middleware {
    return func(next common.HandlerFunc) common.HandlerFunc {
        return func(ctx *common.Context) error {
            start := time.Now()
            
            // Log request
            logger.Info("Request started",
                "method", ctx.Request.Method,
                "path", ctx.Request.URL.Path,
                "ip", ctx.Request.RemoteAddr,
            )
            
            // Process request
            err := next(ctx)
            
            // Log response
            logger.Info("Request completed",
                "method", ctx.Request.Method,
                "path", ctx.Request.URL.Path,
                "status", ctx.Response.Status,
                "duration", time.Since(start).Milliseconds(),
                "error", err,
            )
            
            return err
        }
    }
}
```

---

## Metrics (Prometheus)

### Prometheus Metrics

```go
package metrics

import (
    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
)

var (
    httpRequestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests",
        },
        []string{"method", "endpoint", "status"},
    )
    
    httpRequestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "HTTP request duration in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"method", "endpoint"},
    )
    
    activeConnections = promauto.NewGauge(
        prometheus.GaugeOpts{
            Name: "active_connections",
            Help: "Number of active connections",
        },
    )
)

func MetricsMiddleware() common.Middleware {
    return func(next common.HandlerFunc) common.HandlerFunc {
        return func(ctx *common.Context) error {
            start := time.Now()
            activeConnections.Inc()
            defer activeConnections.Dec()
            
            err := next(ctx)
            
            duration := time.Since(start).Seconds()
            status := ctx.Response.Status
            
            httpRequestsTotal.WithLabelValues(
                ctx.Request.Method,
                ctx.Request.URL.Path,
                fmt.Sprintf("%d", status),
            ).Inc()
            
            httpRequestDuration.WithLabelValues(
                ctx.Request.Method,
                ctx.Request.URL.Path,
            ).Observe(duration)
            
            return err
        }
    }
}
```

### Metrics Endpoint

```go
import "github.com/prometheus/client_golang/prometheus/promhttp"

func (c *MetricsController) Prefix() string {
    return "/metrics"
}

func (c *MetricsController) Routes() []common.Route {
    return []common.Route{
        {
            Method: "GET",
            Path: "/",
            Handler: func(ctx *common.Context) error {
                promhttp.Handler().ServeHTTP(ctx.Response, ctx.Request)
                return nil
            },
        },
    }
}
```

---

## Health Checks

### Health Check Implementation

```go
package health

type HealthCheck struct {
    db    *sql.DB
    redis *redis.Client
}

type HealthStatus struct {
    Status   string            `json:"status"`
    Version  string            `json:"version"`
    Checks   map[string]string `json:"checks"`
}

func (h *HealthCheck) Check(ctx context.Context) *HealthStatus {
    status := &HealthStatus{
        Status:  "ok",
        Version: "1.0.0",
        Checks:  make(map[string]string),
    }
    
    // Check database
    if err := h.db.PingContext(ctx); err != nil {
        status.Checks["database"] = "unhealthy"
        status.Status = "degraded"
    } else {
        status.Checks["database"] = "healthy"
    }
    
    // Check Redis
    if err := h.redis.Ping(ctx).Err(); err != nil {
        status.Checks["redis"] = "unhealthy"
        status.Status = "degraded"
    } else {
        status.Checks["redis"] = "healthy"
    }
    
    return status
}
```

---

This completes the comprehensive NestGo Framework documentation covering all major aspects from getting started to deployment and monitoring!
