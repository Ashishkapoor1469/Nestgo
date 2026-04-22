package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Ashishkapoor1469/Nestgo/cli/utils"
	"github.com/spf13/cobra"
)

// NewCmd creates the `nestgo new` command.
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new [project-name]",
		Short: "Create a new NestGo project",
		Long:  "Scaffolds a complete NestGo project with best-practice structure.",
		Args:  cobra.ExactArgs(1),
		RunE:  runNew,
	}
	cmd.Flags().StringP("template", "t", "rest", "Project template: rest, microservice, graphql")
	return cmd
}

func runNew(cmd *cobra.Command, args []string) error {
	name := args[0]
	tmpl, _ := cmd.Flags().GetString("template")

	utils.PrintHeader("Creating NestGo Project: " + name)
	utils.PrintInfo("Template: " + tmpl)
	fmt.Println()

	// Create project directory.
	if err := os.MkdirAll(name, 0755); err != nil {
		return fmt.Errorf("failed to create project directory: %w", err)
	}

	// Define project structure.
	dirs := []string{
		"cmd/server",
		"internal/modules",
		"internal/common/dto",
		"internal/common/middleware",
		"internal/common/guards",
		"internal/common/schemas",
		"internal/config",
		"migrations",
		"test",
	}

	// Add template-specific directories.
	switch tmpl {
	case "microservice":
		dirs = append(dirs, "internal/health", "internal/metrics")
	case "graphql":
		dirs = append(dirs, "internal/graphql/schema", "internal/graphql/resolvers")
	}

	spinner := utils.StartSpinner("Creating project structure...")

	for _, dir := range dirs {
		path := filepath.Join(name, dir)
		if err := os.MkdirAll(path, 0755); err != nil {
			spinner.StopWithError("Failed to create directory: " + dir)
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	spinner.StopWithSuccess("Project structure created")

	// Generate files.
	files := map[string]string{
		"go.mod":                         goModTemplate,
		"cmd/server/main.go":             mainTemplate,
		"internal/config/config.go":      configTemplate,
		"internal/modules/app_module.go": appModuleTemplate,
		".env":                           envTemplate,
		".env.example":                   envExampleTemplate,
		".gitignore":                     gitignoreTemplate,
		"Makefile":                       makefileTemplate,
		"README.md":                      readmeTemplate,
	}

	modName := "github.com/" + name

	for path, tmplStr := range files {
		fullPath := filepath.Join(name, path)
		if err := writeTemplate(fullPath, tmplStr, map[string]string{
			"Name":       name,
			"ModuleName": modName,
		}); err != nil {
			return fmt.Errorf("failed to create %s: %w", path, err)
		}
		utils.PrintStep("📄", path)
	}

	// Initialize git.
	fmt.Println()
	gitSpinner := utils.StartSpinner("Initializing git repository...")
	gitCmd := exec.Command("git", "init")
	gitCmd.Dir = name
	_ = gitCmd.Run()
	gitSpinner.StopWithSuccess("Git repository initialized")

	// Run go mod tidy.
	tidySpinner := utils.StartSpinner("Installing dependencies...")
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = name
	tidyCmd.Stdout = os.Stdout
	tidyCmd.Stderr = os.Stderr
	_ = tidyCmd.Run()
	tidySpinner.StopWithSuccess("Dependencies installed")

	fmt.Println()
	utils.PrintSuccess("Project " + name + " created successfully!")
	fmt.Println()
	utils.PrintDim("  Get started:")
	utils.PrintDim("    cd " + name)
	utils.PrintDim("    nestgo dev")
	fmt.Println()
	utils.PrintDim("  Generate a resource:")
	utils.PrintDim("    nestgo generate resource users")
	fmt.Println()

	return nil
}

func writeTemplate(path, tmplStr string, data any) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	funcMap := template.FuncMap{
		"title": strings.Title,
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
	}

	tmpl, err := template.New("file").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("template parse error: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, data)
}

// --- Project Templates ---

var goModTemplate = `module {{.ModuleName}}
go 1.22
require (
	github.com/Ashishkapoor1469/Nestgo v0.3.0
)
`

var mainTemplate = `package main

import (
	"log"

	"github.com/Ashishkapoor1469/Nestgo/core"
	"github.com/Ashishkapoor1469/Nestgo/config"
	"{{.ModuleName}}/internal/modules"
	appconfig "{{.ModuleName}}/internal/config"
)

func main() {
	// Load configuration.
	cfg := config.MustLoad[appconfig.AppConfig](".")

	// Create the application.
	app := core.New(
		core.WithAddress(":"+cfg.Port),
		core.WithGlobalPrefix("/api"),
	)

	// Register root module.
	app.RegisterModule(&modules.AppModule{})

	// Start the server.
	log.Printf("🚀 Starting %s on :%s", "{{.Name}}", cfg.Port)
	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
`

var configTemplate = `package config

// AppConfig holds the application configuration.
type AppConfig struct {
	Port        string ` + "`" + `env:"PORT" default:"3000"` + "`" + `
	Environment string ` + "`" + `env:"APP_ENV" default:"development"` + "`" + `
	DBHost      string ` + "`" + `env:"DB_HOST" default:"localhost"` + "`" + `
	DBPort      int    ` + "`" + `env:"DB_PORT" default:"5432"` + "`" + `
	DBUser      string ` + "`" + `env:"DB_USER" default:"postgres"` + "`" + `
	DBPassword  string ` + "`" + `env:"DB_PASSWORD"` + "`" + `
	DBName      string ` + "`" + `env:"DB_NAME" default:"{{.Name}}"` + "`" + `
	JWTSecret   string ` + "`" + `env:"JWT_SECRET" default:"change-me-in-production"` + "`" + `
}
`

var appModuleTemplate = `package modules

import (
	"github.com/Ashishkapoor1469/Nestgo/di"
	"github.com/Ashishkapoor1469/Nestgo/common"
)

// AppModule is the root module of the application.
type AppModule struct{}

func (m *AppModule) Module() common.ModuleConfig {
	return common.ModuleConfig{
		Name:        "app",
		Imports:     []common.Module{},
		Controllers: []common.Controller{},
		Providers:   []di.Provider{},
	}
}
`

var envTemplate = `# {{.Name}} Environment Configuration
PORT=3000
APP_ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME={{.Name}}

# Auth
JWT_SECRET=change-me-in-production
`

var envExampleTemplate = `# {{.Name}} Environment Configuration
PORT=3000
APP_ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME={{.Name}}

# Auth
JWT_SECRET=
`

var gitignoreTemplate = `# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
/bin/
/dist/

# Test
*.test
*.out
coverage.html

# Environment
.env
.env.local
.env.*.local

# IDE
.idea/
.vscode/
*.swp
*.swo

# OS
.DS_Store
Thumbs.db

# Go
vendor/
`

var makefileTemplate = `.PHONY: dev build test lint run

# Development
dev:
	@go run cmd/server/main.go

# Build
build:
	@go build -ldflags="-s -w" -o bin/server cmd/server/main.go

# Run production
run: build
	@./bin/server

# Test
test:
	@go test -v -race ./...

# Test with coverage
test-cover:
	@go test -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html

# Lint
lint:
	@golangci-lint run ./...

# Tidy
tidy:
	@go mod tidy

# Generate
generate:
	@go generate ./...
`

var readmeTemplate = `# {{.Name}}
 
A NestGo application.
 
[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![NestGo](https://img.shields.io/badge/NestGo-Framework-E34F26)](https://github.com/Ashishkapoor1469/Nestgo)
 
---
 
## 🚀 Quick Start
 
` + "```bash" + `
# Install dependencies
go mod tidy
 
# Start development server
nestgo dev
` + "```" + `
 
**Your API is running at: http://localhost:3000/api**
 
> Note: All routes are prefixed with ` + "`/api`" + ` by default
 
---
 
## 📁 Project Structure
 
` + "```" + `
{{.Name}}/
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
` + "```" + `
 
---
 
## 🌐 API Endpoints
 
**Base URL:** ` + "`http://localhost:3000/api`" + `
 
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | ` + "`/api`" + ` | API welcome |
| GET | ` + "`/api/health`" + ` | Health check |
 
` + "```bash" + `
# Test the API
curl http://localhost:3000/api
` + "```" + `
 
---
 
## 🛠️ Commands
 
### Development
` + "```bash" + `
nestgo dev                          # Start with hot-reload
nestgo dev --port=8080              # Custom port
go run cmd/server/main.go           # Run directly
` + "```" + `
 
### Code Generation
` + "```bash" + `
nestgo generate resource <name>     # Generate CRUD resource
nestgo generate module <name>       # Generate module
nestgo generate controller <name>   # Generate controller
nestgo generate service <name>      # Generate service
` + "```" + `
 
### Database
` + "```bash" + `
nestgo migration:create <name>      # Create migration
nestgo migration:run                # Run migrations
nestgo migration:rollback           # Rollback
` + "```" + `
 
### Build & Test
` + "```bash" + `
make build                          # Build binary
make test                           # Run tests
nestgo doctor                       # Check project health
nestgo graph                        # Visualize dependencies
` + "```" + `
 
---
 
## ⚙️ Configuration
 
Create ` + "`.env`" + ` file:
 
` + "```bash" + `
APP_NAME={{.Name}}
APP_PORT=3000
APP_GLOBAL_PREFIX=/api
 
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_USER=postgres
DATABASE_PASSWORD=password
DATABASE_NAME={{.Name}}
 
JWT_SECRET=your-secret-key
LOG_LEVEL=info
` + "```" + `
 
---
 
## 🐳 Docker
 
` + "```bash" + `
# Build and run
docker-compose up -d
 
# Or manually
docker build -t {{.Name}} .
docker run -p 3000:3000 {{.Name}}
` + "```" + `
 
---
 
## 📚 Resources
 
- [NestGo Documentation](https://github.com/Ashishkapoor1469/Nestgo)
- [Go Documentation](https://golang.org/doc/)
 
---
 
**Built with [NestGo](https://github.com/Ashishkapoor1469/Nestgo)** ⭐
`
