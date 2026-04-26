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
	cmd.Flags().StringP("template", "t", "rest", "Project template: rest, microservice")
	return cmd
}

func runNew(cmd *cobra.Command, args []string) error {
	name := args[0]
	tmpl, _ := cmd.Flags().GetString("template")

	utils.PrintBanner()
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
		"internal/modules/health",
		"internal/common/dto",
		"internal/common/middleware",
		"internal/common/guards",
		"internal/config",
		"migrations",
		"test",
	}

	if tmpl == "microservice" {
		dirs = append(dirs, "internal/health", "internal/metrics")
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
	modName := "github.com/" + name
	data := map[string]string{
		"Name":       name,
		"ModuleName": modName,
	}

	files := map[string]string{
		"go.mod":                                          goModTemplate,
		"cmd/server/main.go":                             mainTemplate,
		"internal/config/config.go":                      configTemplate,
		"internal/modules/app_module.go":                 appModuleTemplate,
		"internal/modules/health/health_controller.go":   healthControllerTemplate,
		"internal/modules/health/health_module.go":       healthModuleTemplate,
		".env":                                           envTemplate,
		".env.example":                                   envExampleTemplate,
		".gitignore":                                     gitignoreTemplate,
		"nestgo.json":                                    nestgoJSONTemplate,
		"Makefile":                                       makefileTemplate,
		"README.md":                                      projectReadmeTemplate,
	}

	fmt.Println()
	for path, tmplStr := range files {
		fullPath := filepath.Join(name, path)
		if err := writeTemplate(fullPath, tmplStr, data); err != nil {
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
	utils.PrintSuccess("Project '" + name + "' created successfully!")
	fmt.Println()
	fmt.Println("  Next steps:")
	utils.PrintDim("    cd " + name)
	utils.PrintDim("    nestgo dev")
	fmt.Println()
	fmt.Println("  Your app starts with:")
	utils.PrintDim("    GET /api/health  →  { \"status\": \"ok\" }")
	fmt.Println()
	fmt.Println("  Generate a resource:")
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
		"lower": strings.ToLower,
		"upper": strings.ToUpper,
	}

	t, err := template.New("file").Funcs(funcMap).Parse(tmplStr)
	if err != nil {
		return fmt.Errorf("template parse error: %w", err)
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Execute(f, data)
}

// ─── Project Templates ────────────────────────────────────────────────────────

var goModTemplate = `module {{.ModuleName}}

go 1.22

require (
	github.com/Ashishkapoor1469/Nestgo v0.3.0
)
`

var mainTemplate = `package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Ashishkapoor1469/Nestgo/core"
	"github.com/Ashishkapoor1469/Nestgo/config"
	"{{.ModuleName}}/internal/modules"
	appconfig "{{.ModuleName}}/internal/config"
)

func main() {
	// Load configuration from .env
	cfg := config.MustLoad[appconfig.AppConfig](".")

	port := cfg.Port
	if p := os.Getenv("PORT"); p != "" {
		port = p
	}

	// Create the application.
	app := core.New(
		core.WithAddress(":"+port),
		core.WithGlobalPrefix("/api"),
	)

	// Register root module (which imports all feature modules).
	app.RegisterModule(&modules.AppModule{})

	fmt.Printf("\n  🚀 %s running at http://localhost:%s\n", "{{.Name}}", port)
	fmt.Printf("  🩺 Health check: http://localhost:%s/api/health\n\n", port)

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}
`

var configTemplate = `package config

// AppConfig holds the application configuration loaded from environment.
type AppConfig struct {
	Port        string ` + "`" + `env:"PORT"         default:"3000"` + "`" + `
	Environment string ` + "`" + `env:"APP_ENV"      default:"development"` + "`" + `
	DBHost      string ` + "`" + `env:"DB_HOST"      default:"localhost"` + "`" + `
	DBPort      int    ` + "`" + `env:"DB_PORT"      default:"5432"` + "`" + `
	DBUser      string ` + "`" + `env:"DB_USER"      default:"postgres"` + "`" + `
	DBPassword  string ` + "`" + `env:"DB_PASSWORD"` + "`" + `
	DBName      string ` + "`" + `env:"DB_NAME"      default:"{{.Name}}"` + "`" + `
	JWTSecret   string ` + "`" + `env:"JWT_SECRET"   default:"change-me-in-production"` + "`" + `
}
`

var healthControllerTemplate = `package health

import "github.com/Ashishkapoor1469/Nestgo/common"

// HealthController handles health check requests.
type HealthController struct{}

// NewHealthController creates a new HealthController.
func NewHealthController() *HealthController {
	return &HealthController{}
}

// Prefix returns the route prefix.
func (c *HealthController) Prefix() string {
	return "/health"
}

// Routes returns the controller's route definitions.
func (c *HealthController) Routes() []common.Route {
	return []common.Route{
		{
			Method:  "GET",
			Path:    "/",
			Handler: c.Check,
			Summary: "Health check",
		},
	}
}

// Check returns a 200 OK with status "ok".
func (c *HealthController) Check(ctx *common.Context) error {
	return ctx.OK(map[string]string{
		"status": "ok",
	})
}
`

var healthModuleTemplate = `package health

import (
	"github.com/Ashishkapoor1469/Nestgo/common"
	"github.com/Ashishkapoor1469/Nestgo/di"
)

// HealthModule wires up the health check feature.
type HealthModule struct{}

func (m *HealthModule) Module() common.ModuleConfig {
	controller := NewHealthController()
	return common.ModuleConfig{
		Name:        "health",
		Controllers: []common.Controller{controller},
		Providers:   []di.Provider{},
	}
}
`

var appModuleTemplate = `package modules

import (
	"github.com/Ashishkapoor1469/Nestgo/common"
	"github.com/Ashishkapoor1469/Nestgo/di"
	"{{.ModuleName}}/internal/modules/health"
)

// AppModule is the root module — import all feature modules here.
type AppModule struct{}

func (m *AppModule) Module() common.ModuleConfig {
	return common.ModuleConfig{
		Name: "app",
		Imports: []common.Module{
			&health.HealthModule{},
		},
		Controllers: []common.Controller{},
		Providers:   []di.Provider{},
	}
}
`

var nestgoJSONTemplate = `{
  "name": "{{.Name}}",
  "version": "0.1.0",
  "language": "go",
  "entrypoint": "cmd/server/main.go",
  "sourceRoot": "internal",
  "prefix": "/api",
  "compilerOptions": {
    "deleteOutDir": true,
    "outDir": "bin"
  }
}
`

var envTemplate = `# {{.Name}} — Environment Variables
PORT=3000
APP_ENV=development

# Database (configure when ready)
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=
DB_NAME={{.Name}}

# Security
JWT_SECRET=change-me-in-production
`

var envExampleTemplate = `# {{.Name}} — Environment Variables Example
# Copy this to .env and fill in your values
PORT=3000
APP_ENV=development

# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password-here
DB_NAME={{.Name}}

# Security — CHANGE THIS IN PRODUCTION
JWT_SECRET=your-secret-key-here
`

var gitignoreTemplate = `# Binaries
*.exe
*.exe~
*.dll
*.so
*.dylib
/bin/
/dist/

# Test output
*.test
*.out
coverage.html
coverage.out

# Environment (never commit secrets)
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

var makefileTemplate = `.PHONY: dev build test lint run tidy

# Development with hot reload
dev:
	@nestgo dev

# Build production binary
build:
	@go build -ldflags="-s -w" -o bin/server ./cmd/server

# Run production binary
run: build
	@./bin/server

# Run tests
test:
	@go test -v -race -cover ./...

# Lint
lint:
	@golangci-lint run ./...

# Tidy dependencies
tidy:
	@go mod tidy

# Generate a resource
generate:
	@nestgo generate resource $(name)
`

var projectReadmeTemplate = `# {{.Name}}

A production-ready API built with [NestGo](https://github.com/Ashishkapoor1469/Nestgo).

[![Go](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![NestGo](https://img.shields.io/badge/NestGo-v0.3.0-E34F26)](https://github.com/Ashishkapoor1469/Nestgo)

---

## 🚀 Quick Start

` + "```bash" + `
# Start development server with hot reload
nestgo dev
` + "```" + `

Your API is live at: **http://localhost:3000/api**

---

## ✅ Built-in Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | ` + "`/api/health`" + ` | Health check → ` + "`{\"status\":\"ok\"}`" + ` |

---

## 📁 Project Structure

` + "```" + `
{{.Name}}/
├── cmd/server/          # Application entry point
├── internal/
│   ├── config/          # App configuration (env → struct)
│   └── modules/         # Feature modules
│       └── health/      # Built-in health check module
├── migrations/          # Database migrations
├── test/                # Integration tests
├── .env                 # Environment variables (git-ignored)
├── .env.example         # Example env (commit this)
├── nestgo.json          # NestGo project config
├── Makefile             # Convenience scripts
└── go.mod
` + "```" + `

---

## 🛠️ CLI Commands

` + "```bash" + `
# Development
nestgo dev                         # Start with hot reload

# Code generation
nestgo generate resource users     # Full CRUD resource
nestgo generate module payments    # Module only
nestgo generate controller orders  # Controller only

# Database
nestgo migration:create add_users  # Create migration file
nestgo migration:run               # Run pending migrations
nestgo migration:rollback          # Rollback last migration

# Utilities
nestgo doctor                      # Project health check
nestgo routes                      # Show all registered routes
nestgo build                       # Build production binary
` + "```" + `

---

## ⚙️ Configuration

Edit ` + "`.env`" + `:

` + "```bash" + `
PORT=3000
APP_ENV=development
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=your-password
DB_NAME={{.Name}}
JWT_SECRET=your-secret
` + "```" + `

---

## 📚 Resources

- [NestGo Documentation](https://github.com/Ashishkapoor1469/Nestgo)
- [NestGo Examples](https://github.com/Ashishkapoor1469/Nestgo/tree/main/examples)

---

Built with ❤️ using [NestGo](https://github.com/Ashishkapoor1469/Nestgo)
`
