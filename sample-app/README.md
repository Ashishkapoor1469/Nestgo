# sample-app
 
A NestGo application.
 
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![NestGo](https://img.shields.io/badge/NestGo-Framework-E34F26)](https://github.com/Ashishkapoor1469/Nestgo)
 
---
 
## 🚀 Quick Start
 
```bash
# Install dependencies
go mod tidy
 
# Start development server
nestgo dev
```

**Your API is running at: http://localhost:8080/api**  
*(Note: Ensure your local PostgreSQL cluster is running before executing requests)*

> 💡 **Pro-Tip**: You can interactively test the entire API lifecycle (Auth + Users CRUD) by opening `rest.http` with the VS Code REST Client!
 
---
 
## 📁 Project Structure
 
```
sample-app/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # Configuration
│   ├── modules/         # Feature modules
│   └── common/          # Shared code (guards, middleware, etc.)
├── migrations/          # Database migrations
├── test/                # Integration tests
├── .env                 # Environment variables
├── Makefile             # Build scripts
├── nestgo.yaml          # NestGo configuration
└── go.mod
```
 
---
 
## 🌐 API Endpoints
 
**Base URL:** `http://localhost:8080/api`
 
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/auth/register` | Register a new user |
| POST | `/api/auth/login` | Login safely and return JWT Token |
| GET | `/api/auth/profile` | View profile via authenticated token |
| GET | `/api/users` | List all users (Protected) |
| GET | `/api/users/{id}` | Read a single user |
| PUT | `/api/users/{id}` | Update user attributes |
| DELETE| `/api/users/{id}` | Delete a user |
 
```bash
# The absolute easiest way to test these is natively with the provided rest.http file!
```
 
---
 
## 🛠️ Commands
 
### Development
```bash
nestgo dev                          # Start with hot-reload
nestgo dev --port=8080              # Custom port
go run cmd/server/main.go           # Run directly
```
 
### Code Generation
```bash
nestgo generate resource <name>     # Generate CRUD resource
nestgo generate module <name>       # Generate module
nestgo generate controller <name>   # Generate controller
nestgo generate service <name>      # Generate service
```
 
### Database
```bash
nestgo migration:create <name>      # Create migration
nestgo migration:run                # Run migrations
nestgo migration:rollback           # Rollback
```
 
### Build & Test
```bash
make build                          # Build binary
make test                           # Run tests
nestgo doctor                       # Check project health
nestgo graph                        # Visualize dependencies
```
 
---
 
## ⚙️ Configuration
 
Create `.env` file:
 
```bash
APP_NAME=sample-app
APP_PORT=3000
APP_GLOBAL_PREFIX=/api
 
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=password
DATABASE_NAME=sample-app
 
JWT_SECRET=your-secret-key
LOG_LEVEL=info
```
 
---
 
## 🐳 Docker
 
```bash
# Build and run
docker-compose up -d
 
# Or manually
docker build -t sample-app .
docker run -p 3000:3000 sample-app
```
 
---
 
## 📚 Resources
 
- [NestGo Documentation](https://github.com/Ashishkapoor1469/Nestgo)
- [Go Documentation](https://golang.org/doc/)
 
---
 
**Built with [NestGo](https://github.com/Ashishkapoor1469/Nestgo)** ⭐
