# sample-app

A NestGo application.

## Getting Started

```bash
# Install dependencies
go mod tidy

# Start development server
nestgo dev

# Or run directly
go run cmd/server/main.go
```

## Project Structure

```
sample-app/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # Configuration
│   ├── modules/         # Feature modules
│   └── common/          # Shared code
├── migrations/          # Database migrations
├── test/                # Integration tests
├── .env                 # Environment variables
├── Makefile             # Build scripts
└── go.mod
```

## Commands

| Command | Description |
|---|---|
| `nestgo dev` | Start dev server with hot reload |
| `nestgo generate resource <name>` | Generate a CRUD resource |
| `nestgo generate module <name>` | Generate a module |
| `nestgo doctor` | Analyze project health |
| `make build` | Build production binary |
| `make test` | Run tests |

## Built with [NestGo](https://github.com/Ashishkapoor1469/Nestgo)
