<p align="center">
  <img src="assets/logo.webp" alt="NestGo Logo" width="120" />
</p>

<h1 align="center">NestGo</h1>

<p align="center">
  A progressive Go framework for building efficient, reliable, and scalable server-side applications.
</p>

<p align="center">
  <a href="https://github.com/Ashishkapoor1469/Nestgo/actions/workflows/ci.yml">
    <img src="https://github.com/Ashishkapoor1469/Nestgo/actions/workflows/ci.yml/badge.svg" alt="CI" />
  </a>
  <a href="https://pkg.go.dev/github.com/Ashishkapoor1469/Nestgo">
    <img src="https://pkg.go.dev/badge/github.com/Ashishkapoor1469/Nestgo.svg" alt="Go Reference" />
  </a>
  <a href="https://github.com/Ashishkapoor1469/Nestgo/releases">
    <img src="https://img.shields.io/github/v/release/Ashishkapoor1469/Nestgo?color=blue" alt="Latest Release" />
  </a>
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go" alt="Go Version" />
  <img src="https://img.shields.io/badge/license-MIT-green" alt="License" />
</p>

---

## What is NestGo?

NestGo is an opinionated, modular backend framework for Go — inspired by NestJS — that brings **enterprise-grade architecture** to Go development without sacrificing performance.

It gives you:
- **Modular architecture** with automatic dependency resolution
- **A powerful CLI** (`nestgo`) that scaffolds, generates, and runs your app
- **Compile-time dependency injection** — no reflection, no magic, full type safety
- **Built-in middleware**, guards, interceptors, and lifecycle hooks
- **Production-ready from day one** — not a toy, not a micro-library

---

## Why NestGo?

| Feature | NestGo | Gin | Fiber |
|---------|--------|-----|-------|
| Modular Architecture | ✅ Built-in | ❌ Manual | ❌ Manual |
| Dependency Injection | ✅ Compile-time | ❌ None | ❌ None |
| CLI Tooling | ✅ Advanced | ❌ None | ❌ None |
| Code Generation | ✅ Full CRUD | ❌ None | ❌ None |
| Architecture Linting | ✅ Built-in | ❌ Manual | ❌ Manual |
| Route Explorer (AST) | ✅ Built-in | ❌ None | ❌ None |
| Migration System | ✅ Built-in | ❌ Manual | ❌ Manual |
| Performance | ✅ Native Go | ✅ Native Go | ✅ Native Go |

---

## Installation

```bash
go install github.com/Ashishkapoor1469/Nestgo/cmd/nestgo@latest
```

Verify installation:

```bash
nestgo version
```

---

## Quick Start

```bash
# 1. Create a new project
nestgo new myapp

# 2. Enter the directory
cd myapp

# 3. Start the development server
nestgo dev
```

Your API is live at **http://localhost:3000/api**

```bash
# Built-in health check — works immediately
curl http://localhost:3000/api/health
# → { "status": "ok" }
```

---

## Generate a Resource

The signature feature of NestGo. One command generates a fully wired CRUD resource:

```bash
nestgo generate resource users
```

This creates:

```
internal/modules/users/
  ├── module.go       # Module definition with DI wiring
  ├── controller.go   # REST controller with 5 CRUD routes
  ├── service.go      # Business logic layer
  ├── dto.go          # Create + Update DTOs
  ├── entity.go       # Data model
  └── users_test.go   # Table-driven service tests

migrations/
  └── 1714000000_create_users.sql   # Auto-generated migration
```

Then register it in `internal/modules/app_module.go`:

```go
Imports: []common.Module{
    &health.HealthModule{},
    &users.UsersModule{},  // ← add this
},
```

Your endpoints are now live:

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/users` | List all users (paginated) |
| `GET` | `/api/users/:id` | Get user by ID |
| `POST` | `/api/users` | Create a user |
| `PUT` | `/api/users/:id` | Update a user |
| `DELETE` | `/api/users/:id` | Delete a user |

---

## Project Structure

```
myapp/
├── cmd/server/
│   └── main.go               # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go         # Env → struct config
│   └── modules/
│       ├── app_module.go     # Root module (imports all features)
│       └── health/           # Built-in health module
│           ├── health_controller.go
│           └── health_module.go
├── migrations/               # SQL migration files
├── test/                     # Integration tests
├── .env                      # Environment variables (git-ignored)
├── .env.example              # Example env (commit this)
├── nestgo.json               # NestGo project configuration
├── Makefile                  # Convenience commands
└── go.mod
```

---

## CLI Reference

### Project Commands

```bash
nestgo new <name>                  # Scaffold a new project
nestgo dev                         # Start dev server with hot reload
nestgo dev --port 8080             # Override port
nestgo build                       # Build production binary
nestgo build --output bin/api      # Custom output path
```

### Code Generation

```bash
# Full CRUD resource (recommended)
nestgo generate resource <name>

# Individual components
nestgo generate module <name>
nestgo generate controller <name>
nestgo generate service <name>
nestgo generate guard <name>
nestgo generate middleware <name>
nestgo generate interceptor <name>
nestgo generate dto <name>
nestgo generate test <name>

# Auth scaffolding
nestgo generate auth
```

### Database Migrations

```bash
nestgo migration:create <name>     # Create a new migration file
nestgo migration:run               # Run all pending migrations
nestgo migration:rollback          # Rollback the last migration
nestgo migration:status            # Show migration status
```

### Diagnostics

```bash
nestgo doctor                      # Full project health check with fixes
nestgo routes                      # Show all registered routes (AST-based)
nestgo graph                       # Visualize module dependency graph
nestgo lint-arch                   # Enforce clean architecture boundaries
nestgo version                     # Show CLI + framework + Go versions
```

---

## Configuration

NestGo uses environment variables loaded from `.env`:

```bash
# Server
PORT=3000
APP_ENV=development            # development | production

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=myapp

# Security
JWT_SECRET=change-me-in-production
```

The config is automatically mapped to a typed struct at startup:

```go
type AppConfig struct {
    Port      string `env:"PORT"       default:"3000"`
    DBHost    string `env:"DB_HOST"    default:"localhost"`
    JWTSecret string `env:"JWT_SECRET"`
}
```

---

## Migration System

```bash
# Create a new migration
nestgo migration:create add_users_table

# This creates:
# migrations/1714000000_add_users_table.sql

# Run all pending migrations
nestgo migration:run

# Rollback last migration
nestgo migration:rollback
```

---

## Examples

Working examples are in the [`/examples`](./examples) directory:

| Example | Description |
|---------|-------------|
| [`basic-api`](./examples/basic-api) | Health check + users CRUD (in-memory) |
| [`todo-api`](./examples/todo-api) | Full todo app with persistence |

Run an example:

```bash
cd examples/basic-api
go run ./cmd

# GET /api/health  → { "status": "ok" }
# GET /api/users   → list of users
```

---

## Roadmap

- [x] Modular architecture with DI
- [x] CLI with code generation
- [x] Hot reload dev server
- [x] Route explorer (AST)
- [x] Migration system
- [x] Auth scaffolding (JWT + bcrypt)
- [x] Architecture linter
- [ ] WebSocket support (in progress)
- [ ] OpenAPI/Swagger auto-generation
- [ ] Plugin system
- [ ] gRPC transport layer
- [ ] Rate limiting middleware
- [ ] Caching integration

---

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) first.

```bash
# Clone the repo
git clone https://github.com/Ashishkapoor1469/Nestgo.git
cd Nestgo

# Run tests
go test ./...

# Build the CLI
go build -o nestgo ./cmd/nestgo
```

---

## License

MIT © [Ashish Kapoor](https://github.com/Ashishkapoor1469)

---

<p align="center">
  Built with ❤️ for Go developers who want more structure without more complexity.
</p>
