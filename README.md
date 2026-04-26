<p align="center">
  <img src="assets/logo.png" alt="NestGo Logo" width="100%" />
</p>

<h1 align="center">NestGo</h1>

<p align="center">
  A progressive Go framework for building efficient, reliable, and scalable server-side applications.
</p>

<p align="center">
  <a href="https://github.com/Ashishkapoor1469/Nestgo/actions/workflows/ci.yml">
    <img src="https://github.com/Ashishkapoor1469/Nestgo/actions/workflows/ci.yml/badge.svg" alt="CI Status" />
  </a>
  &nbsp;
  <a href="https://pkg.go.dev/github.com/Ashishkapoor1469/Nestgo">
    <img src="https://pkg.go.dev/badge/github.com/Ashishkapoor1469/Nestgo.svg" alt="Go Reference" />
  </a>
  &nbsp;
  <img src="https://img.shields.io/badge/version-v0.4.0-blue" alt="Version" />
  &nbsp;
  <img src="https://img.shields.io/badge/Go-1.22+-00ADD8?logo=go&logoColor=white" alt="Go Version" />
  &nbsp;
  <img src="https://img.shields.io/badge/license-MIT-22c55e" alt="License MIT" />
</p>

<br />

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
|---------|:------:|:---:|:-----:|
| Modular Architecture | ✅ Built-in | ❌ Manual | ❌ Manual |
| Dependency Injection | ✅ Compile-time | ❌ None | ❌ None |
| CLI Tooling | ✅ Advanced | ❌ None | ❌ None |
| Code Generation | ✅ Full CRUD | ❌ None | ❌ None |
| Architecture Linting | ✅ Built-in | ❌ None | ❌ None |
| Route Explorer | ✅ Built-in | ❌ None | ❌ None |
| Migration System | ✅ Built-in | ❌ Manual | ❌ Manual |
| Hot Reload Dev Server | ✅ Built-in | ❌ None | ❌ None |
| Performance | ✅ Native Go | ✅ Native Go | ✅ Native Go |

---

## Installation

```bash
go install github.com/Ashishkapoor1469/Nestgo/cmd/nestgo@latest
```

Verify:

```bash
nestgo version
# CLI Version:       v0.4.0
# Framework Version: v0.4.0
# Go Version:        go1.22.x
```

---

## Quick Start

Three commands to a running API:

```bash
nestgo new myapp
cd myapp
nestgo dev
```

Your API is live at **`http://localhost:3000/api`**

```bash
# Built-in health check — works immediately, no setup needed
curl http://localhost:3000/api/health
# → { "status": "ok" }
```

---

## Generate a Resource

The signature feature of NestGo. One command generates a fully wired CRUD resource:

```bash
nestgo generate resource users
```

Output:

```
  ✦ Scaffolding resource: users

  ├── internal/modules/users/module.go        # Module with DI wiring
  ├── internal/modules/users/controller.go    # 5 REST endpoints
  ├── internal/modules/users/service.go       # Business logic
  ├── internal/modules/users/dto.go           # Create + Update DTOs
  ├── internal/modules/users/entity.go        # Data model
  ├── internal/modules/users/users_test.go    # Table-driven tests
  └── migrations/1714000000_create_users.sql  # SQL migration

  Next steps:
    1. Register in app_module.go
    2. Run migrations: nestgo migration:run
    3. Test: curl http://localhost:3000/api/users
```

Register the module in `internal/modules/app_module.go`:

```go
Imports: []common.Module{
    &health.HealthModule{},
    &users.UsersModule{},  // ← add this line
},
```

Your CRUD endpoints are now live:

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/users` | List all users (paginated) |
| `GET` | `/api/users/{id}` | Get user by ID |
| `POST` | `/api/users` | Create a user |
| `PUT` | `/api/users/{id}` | Update a user |
| `DELETE` | `/api/users/{id}` | Delete a user |

---

## Project Structure

```
myapp/
├── cmd/
│   └── server/
│       └── main.go               # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go             # Env → typed struct config
│   └── modules/
│       ├── app_module.go         # Root module — register features here
│       └── health/               # Built-in health check module
│           ├── health_controller.go
│           └── health_module.go
├── migrations/                   # SQL migration files
├── test/                         # Integration tests
├── .env                          # Environment variables (git-ignored)
├── .env.example                  # Template env file (commit this)
├── nestgo.json                   # NestGo project configuration
├── Makefile                      # make dev / build / test
└── go.mod
```

---

## CLI Reference

### Project

```bash
nestgo new <name>              # Scaffold a new NestGo project
nestgo dev                     # Start dev server with hot reload
nestgo dev --port 8080         # Override port
nestgo build                   # Build production binary (bin/server)
nestgo build --output bin/api  # Custom output path
nestgo version                 # Show version info
```

### Code Generation

```bash
# Full CRUD resource — the fast path
nestgo generate resource <name>

# Individual pieces
nestgo generate module <name>
nestgo generate controller <name>
nestgo generate service <name>
nestgo generate guard <name>
nestgo generate middleware <name>
nestgo generate interceptor <name>
nestgo generate dto <name>
nestgo generate test <name>

# Auth module (JWT + bcrypt)
nestgo generate auth
```

### Database

```bash
nestgo migration:create <name>   # Create a new .sql migration file
nestgo migration:run             # Run all pending migrations
nestgo migration:rollback        # Rollback the last migration
nestgo migration:status          # Show pending / applied migrations
```

### Diagnostics

```bash
nestgo doctor      # Full health check with actionable fix suggestions
nestgo routes      # List all registered routes (parsed from AST)
nestgo graph       # Visualize module dependency tree
nestgo lint-arch   # Detect clean architecture violations
```

---

## Configuration

NestGo maps environment variables to a typed config struct at startup — no manual parsing:

```go
// internal/config/config.go
type AppConfig struct {
    Port        string `env:"PORT"         default:"3000"`
    Environment string `env:"APP_ENV"      default:"development"`
    DBHost      string `env:"DB_HOST"      default:"localhost"`
    DBPort      int    `env:"DB_PORT"      default:"5432"`
    DBUser      string `env:"DB_USER"      default:"postgres"`
    DBPassword  string `env:"DB_PASSWORD"`
    DBName      string `env:"DB_NAME"`
    JWTSecret   string `env:"JWT_SECRET"   default:"change-me"`
}
```

Your `.env` file:

```env
PORT=3000
APP_ENV=development

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME=myapp

JWT_SECRET=change-me-in-production
```

---

## Migration System

```bash
# 1. Create a migration
nestgo migration:create add_users_table
# → creates migrations/1714000000_add_users_table.sql

# 2. Run pending migrations
nestgo migration:run

# 3. Check status
nestgo migration:status

# 4. Rollback if needed
nestgo migration:rollback
```

---

## Examples

Fully working examples in the [`/examples`](./examples) directory:

| Example | Description | Run |
|---------|-------------|-----|
| [`basic-api`](./examples/basic-api) | Health + users CRUD (in-memory) | `cd examples/basic-api && go run ./cmd` |
| [`todo-api`](./examples/todo-api) | Todo app with persistence | `cd examples/todo-api && go run ./cmd` |

```bash
cd examples/basic-api
go run ./cmd

# GET /api/health  → { "status": "ok", "service": "basic-api" }
# GET /api/users   → { "data": [...], "total": 2 }
# POST /api/users  → created user
```

---

## `nestgo doctor` Output

```
  🩺 NestGo Doctor — Project Health Check

  ✅ go.mod exists
  ✅ nestgo.json exists
  ✅ .env file exists
  ✅ .env.example exists
  ✅ cmd/ directory exists
  ✅ internal/ directory exists
  ✅ internal/modules/ exists
  ✅ migrations/ directory exists
  ✅ Entry point found (cmd/server/main.go)
  ✅ Go toolchain: go1.22.3 linux/amd64
  ✅ Dependencies verified
  ⚠️  Test files found (0)
      → Run: nestgo generate test <module-name>

  ─────────────────────────────────
  ✅ Passed:   11
  ⚠️  Warnings: 1
  ❌ Issues:   0
  ─────────────────────────────────

  🟡 Project is healthy but has some recommendations.
```

---

## Roadmap

- [x] Modular architecture with DI
- [x] CLI with full code generation
- [x] Hot reload dev server
- [x] Route explorer (AST-based)
- [x] SQL migration system
- [x] Auth scaffolding (JWT + bcrypt)
- [x] Architecture linter
- [x] `nestgo doctor` health checks
- [ ] WebSocket support
- [ ] OpenAPI / Swagger auto-generation
- [ ] Rate limiting middleware
- [ ] Caching integration (Redis)
- [ ] gRPC transport layer

---

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) first.

```bash
# Clone
git clone https://github.com/Ashishkapoor1469/Nestgo.git
cd Nestgo

# Run tests
go test ./...

# Build the CLI locally
go build -o bin/nestgo ./cmd/nestgo

# Try it
./bin/nestgo --help
```

---

## License

MIT © [Ashish Kapoor](https://github.com/Ashishkapoor1469)

---

<p align="center">
  Built with ❤️ for Go developers who want structure without complexity.
  <br />
  <a href="https://github.com/Ashishkapoor1469/Nestgo">⭐ Star us on GitHub</a>
</p>
