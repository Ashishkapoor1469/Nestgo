package commands

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/nestgo/nestgo/cli/utils"
	"github.com/spf13/cobra"
)

// DoctorCmd creates the `nestgo doctor` command.
func DoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Analyze project structure and detect issues",
		Long:  "Runs diagnostics on your NestGo project to detect anti-patterns, missing imports, and architectural issues.",
		RunE:  runDoctor,
	}
}

// GraphCmd creates the `nestgo graph` command.
func GraphCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "graph",
		Short: "Visualize module dependency graph",
		Long:  "Generates a text-based dependency graph of your NestGo modules.",
		RunE:  runGraph,
	}
}

// MigrateCmd creates the `nestgo migrate` command.
func MigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "up",
			Short: "Run all pending migrations",
			RunE: func(cmd *cobra.Command, args []string) error {
				utils.EnsureProjectContext("migrate up")
				fmt.Println("⬆️  Running migrations...")
				fmt.Println("   (Implement by calling your migration runner)")
				return nil
			},
		},
		&cobra.Command{
			Use:   "down",
			Short: "Rollback the last migration",
			RunE: func(cmd *cobra.Command, args []string) error {
				utils.EnsureProjectContext("migrate down")
				fmt.Println("⬇️  Rolling back last migration...")
				return nil
			},
		},
		&cobra.Command{
			Use:   "create [name]",
			Short: "Create a new migration file",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return createMigration(args[0])
			},
		},
	)

	return cmd
}

func runDoctor(cmd *cobra.Command, args []string) error {
	utils.EnsureProjectContext("doctor")
	fmt.Println("\n🩺 NestGo Doctor — Project Health Check")

	issues := 0
	warnings := 0
	passed := 0

	// Check 1: Project structure.
	fmt.Println("  Checking project structure...")
	requiredDirs := []string{"cmd", "internal"}
	for _, dir := range requiredDirs {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			fmt.Printf("    ❌ Missing directory: %s/\n", dir)
			issues++
		} else {
			fmt.Printf("    ✅ Found %s/\n", dir)
			passed++
		}
	}

	// Check 2: Go mod.
	fmt.Println("\n  Checking Go module...")
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		fmt.Println("    ❌ Missing go.mod file")
		issues++
	} else {
		fmt.Println("    ✅ go.mod found")
		passed++
	}

	// Check 3: Environment files.
	fmt.Println("\n  Checking configuration...")
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		fmt.Println("    ⚠️  Missing .env file (recommended)")
		warnings++
	} else {
		fmt.Println("    ✅ .env file found")
		passed++
	}

	if _, err := os.Stat(".env.example"); os.IsNotExist(err) {
		fmt.Println("    ⚠️  Missing .env.example file (recommended for team)")
		warnings++
	} else {
		fmt.Println("    ✅ .env.example found")
		passed++
	}

	// Check 4: Find modules.
	fmt.Println("\n  Scanning modules...")
	moduleCount := 0
	controllerCount := 0
	serviceCount := 0

	_ = filepath.Walk("internal/modules", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}

		src := string(content)
		if strings.Contains(src, "ModuleConfig") {
			moduleCount++
		}
		if strings.Contains(src, "ControllerDefinition") || strings.Contains(src, "Prefix()") {
			controllerCount++
		}
		if strings.Contains(src, "Service struct") {
			serviceCount++
		}

		return nil
	})

	fmt.Printf("    📦 Modules: %d\n", moduleCount)
	fmt.Printf("    🎮 Controllers: %d\n", controllerCount)
	fmt.Printf("    ⚙️  Services: %d\n", serviceCount)

	// Check 5: Look for common anti-patterns.
	fmt.Println("\n  Checking for anti-patterns...")
	antiPatterns := checkAntiPatterns()
	if len(antiPatterns) == 0 {
		fmt.Println("    ✅ No anti-patterns detected")
		passed++
	} else {
		for _, ap := range antiPatterns {
			fmt.Printf("    ⚠️  %s\n", ap)
			warnings++
		}
	}

	// Check 6: Test files.
	fmt.Println("\n  Checking test coverage...")
	testCount := 0
	_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			testCount++
		}
		return nil
	})
	if testCount == 0 {
		fmt.Println("    ⚠️  No test files found")
		warnings++
	} else {
		fmt.Printf("    ✅ Found %d test files\n", testCount)
		passed++
	}

	// Summary.
	fmt.Println("\n  ─────────────────────────────────")
	fmt.Printf("  ✅ Passed:   %d\n", passed)
	fmt.Printf("  ⚠️  Warnings: %d\n", warnings)
	fmt.Printf("  ❌ Issues:   %d\n", issues)
	fmt.Println("  ─────────────────────────────────")

	if issues > 0 {
		fmt.Println("\n  🔴 Project has issues that should be fixed.")
	} else if warnings > 0 {
		fmt.Println("\n  🟡 Project is healthy with some recommendations.")
	} else {
		fmt.Println("\n  🟢 Project is in great shape!")
	}
	fmt.Println()

	return nil
}

func runGraph(cmd *cobra.Command, args []string) error {
	utils.EnsureProjectContext("graph")
	fmt.Println("\n📊 Module Dependency Graph")

	modules := findModules()
	if len(modules) == 0 {
		fmt.Println("  No modules found in internal/modules/")
		return nil
	}

	for _, mod := range modules {
		fmt.Printf("  📦 %s\n", mod.name)
		for _, dep := range mod.imports {
			fmt.Printf("     └── %s\n", dep)
		}
		if len(mod.controllers) > 0 {
			for _, ctrl := range mod.controllers {
				fmt.Printf("     ├── 🎮 %s\n", ctrl)
			}
		}
		if len(mod.services) > 0 {
			for _, svc := range mod.services {
				fmt.Printf("     ├── ⚙️  %s\n", svc)
			}
		}
	}
	fmt.Println()

	return nil
}

func createMigration(name string) error {
	utils.EnsureProjectContext("migrate create " + name)
	dir := "migrations"
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	timestamp := strings.ReplaceAll(strings.ReplaceAll(
		strings.Split(fmt.Sprintf("%v", os.Getenv("timestamp")), ".")[0],
		"-", ""), ":", "")
	if timestamp == "" {
		timestamp = fmt.Sprintf("%d", os.Getpid())
	}

	fileName := fmt.Sprintf("%s_%s.go", timestamp, name)
	filePath := filepath.Join(dir, fileName)

	content := fmt.Sprintf(`package migrations

import "database/sql"

// %s migration
func Up_%s(db *sql.DB) error {
	_, err := db.Exec(`+"`"+`
		-- Add your migration SQL here
	`+"`"+`)
	return err
}

func Down_%s(db *sql.DB) error {
	_, err := db.Exec(`+"`"+`
		-- Add your rollback SQL here
	`+"`"+`)
	return err
}
`, name, name, name)

	if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
		return err
	}

	fmt.Printf("✅ Created migration: %s\n", filePath)
	return nil
}

func checkAntiPatterns() []string {
	var issues []string

	_ = filepath.Walk("internal", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}

		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, path, nil, parser.ImportsOnly)
		if err != nil {
			return nil
		}

		for _, imp := range node.Imports {
			importPath := strings.Trim(imp.Path.Value, `"`)
			if importPath == "fmt" {
				// Check if fmt.Println is used (should use logger instead).
				content, _ := os.ReadFile(path)
				if strings.Contains(string(content), "fmt.Println") {
					issues = append(issues, fmt.Sprintf("%s: Use structured logger instead of fmt.Println", path))
				}
			}
		}

		return nil
	})

	return issues
}

type moduleInfo struct {
	name        string
	imports     []string
	controllers []string
	services    []string
}

func findModules() []moduleInfo {
	var modules []moduleInfo

	modulesDir := filepath.Join("internal", "modules")
	entries, err := os.ReadDir(modulesDir)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		modDir := filepath.Join(modulesDir, entry.Name())
		mod := moduleInfo{name: entry.Name()}

		_ = filepath.Walk(modDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || filepath.Ext(path) != ".go" {
				return nil
			}

			fset := token.NewFileSet()
			node, parseErr := parser.ParseFile(fset, path, nil, 0)
			if parseErr != nil {
				return nil
			}

			for _, decl := range node.Decls {
				if genDecl, ok := decl.(*ast.GenDecl); ok {
					for _, spec := range genDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok {
							name := typeSpec.Name.Name
							if strings.HasSuffix(name, "Controller") {
								mod.controllers = append(mod.controllers, name)
							}
							if strings.HasSuffix(name, "Service") {
								mod.services = append(mod.services, name)
							}
						}
					}
				}
			}

			return nil
		})

		modules = append(modules, mod)
	}

	return modules
}
