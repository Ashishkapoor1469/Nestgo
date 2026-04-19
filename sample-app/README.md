# 🚀 sample-app

A modern, production-ready NestGo application built with Go 1.22+

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![Framework](https://img.shields.io/badge/Framework-NestGo-E34F26?style=flat)](https://github.com/Ashishkapoor1469/Nestgo)
[![License](https://img.shields.io/badge/License-MIT-success?style=flat)](LICENSE)

---

## 📋 Table of Contents

- [Overview](#overview)
- [Quick Start](#quick-start)
- [Project Structure](#project-structure)
- [Available Endpoints](#available-endpoints)
- [Configuration](#configuration)
- [Development](#development)
- [Commands Reference](#commands-reference)
- [Testing](#testing)
- [Deployment](#deployment)
- [Built With](#built-with)

---

## 🌟 Overview

This is a sample application demonstrating the power and elegance of the **NestGo** framework. It showcases:

- ✅ **Modular Architecture** - Clean separation of concerns with feature modules
- ✅ **Dependency Injection** - Type-safe, compile-time validated dependencies
- ✅ **RESTful API** - Well-structured HTTP endpoints with proper routing
- ✅ **Hot Reload** - Fast development with live code reloading
- ✅ **Database Integration** - Ready for PostgreSQL, MySQL, or MongoDB
- ✅ **Middleware Support** - CORS, logging, authentication, and more
- ✅ **Production Ready** - Dockerized and ready for deployment

---

## 🚀 Quick Start

### Prerequisites

- **Go**: 1.22 or higher ([Download](https://golang.org/dl/))
- **NestGo CLI**: Install globally
  ```bash
  go install github.com/nestgo/nestgo/cmd/nestgo@latest
  ```

### Installation & Setup

```bash
# Clone the repository
git clone https://github.com/yourusername/sample-app.git
cd sample-app

# Install dependencies
go mod tidy

# Copy environment file
cp .env.example .env

# Start development server
nestgo dev
```

### 🎉 Your API is now running!

Open your browser and visit:

**🌐 http://localhost:3000/api**

You should see:
```json
{
  "message": "Welcome to sample-app API",
  "version": "1.0.0",
  "status": "ok"
}
```

> **Note:** The global prefix `/api` is set by default. All routes are prefixed with `/api`

---

## 📁 Project Structure

```
sample-app/
├── cmd/
│   └── server/              # Application entry point
│       └── main.go          # Bootstrap and configuration
│
├── internal/                # Private application code
│   ├── app/                 # Root application module
│   │   ├── app.module.go
│   │   └── app.controller.go
│   │
│   ├── modules/             # Feature modules
│   │   ├── users/           # Users module
│   │   │   ├── users.module.go
│   │   │   ├── users.controller.go
│   │   │   ├── users.service.go
│   │   │   └── dto/
│   │   │       ├── create-user.dto.go
│   │   │       └── update-user.dto.go
│   │   │
│   │   └── products/        # Products module (example)
│   │       ├── products.module.go
│   │       ├── products.controller.go
│   │       └── products.service.go
│   │
│   ├── common/              # Shared utilities
│   │   ├── guards/          # Authentication guards
│   │   ├── interceptors/    # Response interceptors
│   │   ├── middleware/      # HTTP middleware
│   │   └── filters/         # Exception filters
│   │
│   └── config/              # Configuration management
│       ├── config.go
│       └── config.yaml
│
├── migrations/              # Database migrations
│   ├── 001_create_users_table.up.sql
│   └── 001_create_users_table.down.sql
│
├── test/                    # Integration & E2E tests
│   ├── integration/
│   └── e2e/
│
├── pkg/                     # Public reusable packages
├── .env                     # Environment variables (git-ignored)
├── .env.example             # Environment template
├── Dockerfile               # Docker configuration
├── docker-compose.yml       # Docker Compose setup
├── Makefile                 # Build automation scripts
├── go.mod                   # Go module dependencies
├── go.sum                   # Dependency checksums
├── nestgo.yaml              # NestGo configuration
└── README.md               # This file
```

### 📦 Module Organization

Each feature module follows a consistent pattern:

```
users/
├── users.module.go       # Module definition & DI container
├── users.controller.go   # HTTP routes & handlers
├── users.service.go      # Business logic
├── users.repository.go   # Data access layer
├── dto/                  # Data Transfer Objects
│   ├── create-user.dto.go
│   └── update-user.dto.go
└── entities/             # Domain models
    └── user.entity.go
```

---

## 🌐 Available Endpoints

### Base URL
```
http://localhost:3000/api
```

### Core Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api` | API welcome message |
| `GET` | `/api/health` | Health check endpoint |
| `GET` | `/api/health/ready` | Readiness probe |
| `GET` | `/api/health/live` | Liveness probe |

### User Endpoints (Example)

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| `GET` | `/api/users` | Get all users | ❌ |
| `GET` | `/api/users/:id` | Get user by ID | ❌ |
| `POST` | `/api/users` | Create new user | ✅ |
| `PUT` | `/api/users/:id` | Update user | ✅ |
| `DELETE` | `/api/users/:id` | Delete user | ✅ |

### Testing the API

```bash
# Get API info
curl http://localhost:3000/api

# Health check
curl http://localhost:3000/api/health

# Get all users
curl http://localhost:3000/api/users

# Get specific user
curl http://localhost:3000/api/users/123

# Create user
curl -X POST http://localhost:3000/api/users \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","name":"John Doe"}'
```

---

## ⚙️ Configuration

### Environment Variables

Create a `.env` file in the root directory:

```bash
# Application
APP_NAME=sample-app
APP_ENV=development
APP_PORT=3000
APP_GLOBAL_PREFIX=/api

# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=password
DATABASE_NAME=sample_app
DATABASE_SSL_MODE=disable

# Redis (optional)
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT Authentication
JWT_SECRET=your-super-secret-key-change-in-production
JWT_EXPIRES_IN=24h

# CORS
CORS_ORIGINS=http://localhost:3000,http://localhost:8080
CORS_CREDENTIALS=true

# Logging
LOG_LEVEL=info
LOG_FORMAT=json
```

### NestGo Configuration

Edit `nestgo.yaml` for framework-specific settings:

```yaml
app:
  name: sample-app
  port: ${APP_PORT}
  globalPrefix: /api
  
server:
  timeout: 30s
  maxBodySize: 10MB
  
features:
  swagger: true
  metrics: true
  healthCheck: true
```

---

## 🛠️ Development

### Running the Application

```bash
# Development mode with hot-reload
nestgo dev

# Run directly with Go
go run cmd/server/main.go

# Development with verbose logging
nestgo dev --verbose

# Custom port
nestgo dev --port=8080
```

### Code Generation

```bash
# Generate a complete CRUD resource
nestgo generate resource products

# Generate individual components
nestgo generate module orders
nestgo generate controller orders
nestgo generate service orders

# Generate guards, interceptors, middleware
nestgo generate guard auth
nestgo generate interceptor logging
nestgo generate middleware cors
```

### Database Operations

```bash
# Create a new migration
nestgo migration:create add_user_avatar

# Run all pending migrations
nestgo migration:run

# Rollback last migration
nestgo migration:rollback

# Check migration status
nestgo migration:status
```

### Project Health

```bash
# Analyze project health
nestgo doctor

# Visualize module dependencies
nestgo graph

# Check for circular dependencies
nestgo graph --check-cycles
```

---

## 📝 Commands Reference

| Command | Description | Example |
|---------|-------------|---------|
| `nestgo dev` | Start dev server with hot-reload | `nestgo dev --port=8080` |
| `nestgo build` | Build production binary | `nestgo build --optimize` |
| `nestgo generate resource <name>` | Generate CRUD resource | `nestgo generate resource products` |
| `nestgo generate module <name>` | Generate a module | `nestgo generate module auth` |
| `nestgo generate controller <name>` | Generate a controller | `nestgo generate controller users` |
| `nestgo generate service <name>` | Generate a service | `nestgo generate service email` |
| `nestgo migration:create <name>` | Create new migration | `nestgo migration:create add_posts` |
| `nestgo migration:run` | Run migrations | `nestgo migration:run` |
| `nestgo migration:rollback` | Rollback last migration | `nestgo migration:rollback` |
| `nestgo doctor` | Check project health | `nestgo doctor` |
| `nestgo graph` | Visualize dependencies | `nestgo graph --output=graph.png` |
| `make build` | Build binary (via Makefile) | `make build` |
| `make test` | Run all tests | `make test` |
| `make docker-build` | Build Docker image | `make docker-build` |
| `make docker-run` | Run with Docker | `make docker-run` |

---

## 🧪 Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run specific package
go test ./internal/modules/users/...

# Run tests with verbose output
go test -v ./...

# Using Make
make test
make test-coverage
```

### Test Structure

```bash
# Unit tests (alongside code)
internal/modules/users/users.service_test.go

# Integration tests
test/integration/users_test.go

# E2E tests
test/e2e/api_test.go
```

### Example Test

```go
package users_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestUsersService_Create(t *testing.T) {
    service := NewUsersService()
    
    user, err := service.Create(&CreateUserDTO{
        Email: "test@example.com",
        Name:  "Test User",
    })
    
    assert.NoError(t, err)
    assert.NotEmpty(t, user.ID)
    assert.Equal(t, "test@example.com", user.Email)
}
```

---

## 🐳 Deployment

### Docker

```bash
# Build Docker image
docker build -t sample-app:latest .

# Run container
docker run -p 3000:3000 --env-file .env sample-app:latest

# Using Docker Compose
docker-compose up -d

# Stop containers
docker-compose down
```

### Production Build

```bash
# Build optimized binary
make build

# Or manually
go build -ldflags="-s -w" -o bin/app cmd/server/main.go

# Run production binary
./bin/app
```

### Kubernetes

```bash
# Apply Kubernetes manifests
kubectl apply -f k8s/

# Check deployment status
kubectl get pods -l app=sample-app

# View logs
kubectl logs -f deployment/sample-app
```

---

## 📊 Monitoring

### Metrics

Prometheus metrics available at:
```
http://localhost:3000/metrics
```

### Health Checks

- **Liveness**: `GET /api/health/live` - Is the app running?
- **Readiness**: `GET /api/health/ready` - Can the app serve traffic?

### Logging

Structured JSON logs with `slog`:

```json
{
  "time": "2026-04-18T10:30:00Z",
  "level": "INFO",
  "msg": "Request completed",
  "method": "GET",
  "path": "/api/users",
  "status": 200,
  "duration": 45
}
```

---

## 🔐 Security

- ✅ CORS configured and enabled
- ✅ Security headers (Helmet)
- ✅ Rate limiting on public endpoints
- ✅ Input validation with DTOs
- ✅ JWT authentication ready
- ✅ SQL injection prevention
- ✅ XSS protection

---

## 📚 Learn More

### NestGo Documentation
- [Official Repository](https://github.com/Ashishkapoor1469/Nestgo)
- [Getting Started Guide](https://github.com/Ashishkapoor1469/Nestgo#getting-started)
- [CLI Documentation](https://github.com/Ashishkapoor1469/Nestgo#cli-tool)

### Go Resources
- [Effective Go](https://golang.org/doc/effective_go)
- [Go by Example](https://gobyexample.com/)
- [Go Modules](https://blog.golang.org/using-go-modules)

---

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

Built with ❤️ using [NestGo](https://github.com/Ashishkapoor1469/Nestgo)

**NestGo** - A next-generation, production-grade backend framework for Go, inspired by NestJS.

---

## 📞 Support

- 📧 Email: support@sample-app.com
- 💬 Discord: [Join our community](https://discord.gg/sample-app)
- 🐛 Issues: [GitHub Issues](https://github.com/yourusername/sample-app/issues)

---

<div align="center">

**⭐ Star this repo if you find it helpful!**

Made with [NestGo](https://github.com/Ashishkapoor1469/Nestgo) | [Documentation](https://docs.nestgo.dev) | [Examples](https://github.com/Ashishkapoor1469/Nestgo/tree/main/examples)

</div>