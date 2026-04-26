package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Ashishkapoor1469/Nestgo/cli/utils"
	"github.com/spf13/cobra"
)

// DoctorCmd creates the `nestgo doctor` command.
func DoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Run project health checks and diagnose issues",
		Long:  "Checks project structure, configuration, dependencies, and environment. Provides actionable fixes for every issue found.",
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

type check struct {
	label string
	pass  bool
	warn  bool
	fix   string
}

func runDoctor(cmd *cobra.Command, args []string) error {
	utils.EnsureProjectContext("doctor")

	fmt.Println()
	fmt.Println("  🩺 NestGo Doctor — Project Health Check")
	fmt.Println()

	var checks []check

	// ── Checks ────────────────────────────────────────────────────────────────

	// 1. go.mod
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		checks = append(checks, check{
			label: "go.mod exists",
			fix:   "Run: go mod init <your-module-name>",
		})
	} else {
		checks = append(checks, check{label: "go.mod exists", pass: true})
	}

	// 2. nestgo.json
	if _, err := os.Stat("nestgo.json"); os.IsNotExist(err) {
		checks = append(checks, check{
			label: "nestgo.json exists",
			warn:  true,
			fix:   "Create nestgo.json (run: nestgo new <name> to get one automatically)",
		})
	} else {
		checks = append(checks, check{label: "nestgo.json exists", pass: true})
	}

	// 3. .env file
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		checks = append(checks, check{
			label: ".env file exists",
			warn:  true,
			fix:   "Run: cp .env.example .env",
		})
	} else {
		checks = append(checks, check{label: ".env file exists", pass: true})
	}

	// 4. .env.example
	if _, err := os.Stat(".env.example"); os.IsNotExist(err) {
		checks = append(checks, check{
			label: ".env.example exists",
			warn:  true,
			fix:   "Create .env.example so teammates can set up their environment",
		})
	} else {
		checks = append(checks, check{label: ".env.example exists", pass: true})
	}

	// 5. cmd/ directory
	if _, err := os.Stat("cmd"); os.IsNotExist(err) {
		checks = append(checks, check{
			label: "cmd/ directory exists",
			fix:   "Create cmd/server/main.go — this is your entry point",
		})
	} else {
		checks = append(checks, check{label: "cmd/ directory exists", pass: true})
	}

	// 6. internal/ directory
	if _, err := os.Stat("internal"); os.IsNotExist(err) {
		checks = append(checks, check{
			label: "internal/ directory exists",
			fix:   "Create internal/modules/ to hold your feature modules",
		})
	} else {
		checks = append(checks, check{label: "internal/ directory exists", pass: true})
	}

	// 7. internal/modules/
	if _, err := os.Stat(filepath.Join("internal", "modules")); os.IsNotExist(err) {
		checks = append(checks, check{
			label: "internal/modules/ exists",
			warn:  true,
			fix:   "Create internal/modules/ and add your feature modules there",
		})
	} else {
		checks = append(checks, check{label: "internal/modules/ exists", pass: true})
	}

	// 8. migrations/ directory
	if _, err := os.Stat("migrations"); os.IsNotExist(err) {
		checks = append(checks, check{
			label: "migrations/ directory exists",
			warn:  true,
			fix:   "Run: mkdir migrations",
		})
	} else {
		checks = append(checks, check{label: "migrations/ directory exists", pass: true})
	}

	// 9. Entry point exists
	entry := findEntryPoint()
	if entry == "" {
		checks = append(checks, check{
			label: "Entry point (main.go) found",
			fix:   "Create cmd/server/main.go — see the NestGo docs for a starter template",
		})
	} else {
		checks = append(checks, check{label: fmt.Sprintf("Entry point found (%s)", entry), pass: true})
	}

	// 10. Go toolchain
	goCmd := exec.Command("go", "version")
	goOut, err := goCmd.Output()
	if err != nil {
		checks = append(checks, check{
			label: "Go toolchain available",
			fix:   "Install Go from https://golang.org/dl/",
		})
	} else {
		goVer := strings.TrimSpace(strings.TrimPrefix(string(goOut), "go version "))
		checks = append(checks, check{label: fmt.Sprintf("Go toolchain: %s", goVer), pass: true})
	}

	// 11. go mod dependencies
	tidyCmd := exec.Command("go", "mod", "verify")
	if err := tidyCmd.Run(); err != nil {
		checks = append(checks, check{
			label: "Dependencies verified",
			warn:  true,
			fix:   "Run: go mod tidy",
		})
	} else {
		checks = append(checks, check{label: "Dependencies verified", pass: true})
	}

	// 12. Test files
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
		checks = append(checks, check{
			label: "Test files found",
			warn:  true,
			fix:   "Run: nestgo generate test <module-name> to scaffold tests",
		})
	} else {
		checks = append(checks, check{label: fmt.Sprintf("Test files found (%d)", testCount), pass: true})
	}

	// ── Print Results ──────────────────────────────────────────────────────────

	issues := 0
	warnings := 0
	passed := 0

	for _, c := range checks {
		switch {
		case c.pass:
			fmt.Printf("  ✅ %s\n", c.label)
			passed++
		case c.warn:
			fmt.Printf("  ⚠️  %s\n", c.label)
			if c.fix != "" {
				fmt.Printf("      → %s\n", c.fix)
			}
			warnings++
		default:
			fmt.Printf("  ❌ %s\n", c.label)
			if c.fix != "" {
				fmt.Printf("      → %s\n", c.fix)
			}
			issues++
		}
	}

	// Summary
	fmt.Println()
	fmt.Println("  ─────────────────────────────────")
	fmt.Printf("  ✅ Passed:   %d\n", passed)
	fmt.Printf("  ⚠️  Warnings: %d\n", warnings)
	fmt.Printf("  ❌ Issues:   %d\n", issues)
	fmt.Println("  ─────────────────────────────────")

	switch {
	case issues > 0:
		fmt.Println("\n  🔴 Project has critical issues that must be fixed.")
	case warnings > 0:
		fmt.Println("\n  🟡 Project is healthy but has some recommendations.")
	default:
		fmt.Println("\n  🟢 Project is in excellent shape!")
	}
	fmt.Println()

	// Module scan
	fmt.Println("  📦 Module Summary:")
	moduleCount, controllerCount, serviceCount := scanModules()
	fmt.Printf("     Modules:     %d\n", moduleCount)
	fmt.Printf("     Controllers: %d\n", controllerCount)
	fmt.Printf("     Services:    %d\n", serviceCount)
	fmt.Println()

	return nil
}

func scanModules() (modules, controllers, services int) {
	_ = filepath.Walk("internal/modules", func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() || filepath.Ext(path) != ".go" {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		src := string(content)
		if strings.Contains(src, "ModuleConfig") {
			modules++
		}
		if strings.Contains(src, "Prefix()") {
			controllers++
		}
		if strings.Contains(src, "Service struct") {
			services++
		}
		return nil
	})
	return
}

func runGraph(cmd *cobra.Command, args []string) error {
	utils.EnsureProjectContext("graph")

	fmt.Println()
	fmt.Println("  📊 Module Dependency Graph")
	fmt.Println()

	mods := findModules()
	if len(mods) == 0 {
		fmt.Println("  No modules found in internal/modules/")
		fmt.Println("  Run: nestgo generate module <name>")
		fmt.Println()
		return nil
	}

	for i, mod := range mods {
		connector := "├──"
		if i == len(mods)-1 {
			connector = "└──"
		}
		fmt.Printf("  %s 📦 %s\n", connector, mod.name)
		for _, ctrl := range mod.controllers {
			fmt.Printf("  │    ├── 🎮 %s\n", ctrl)
		}
		for _, svc := range mod.services {
			fmt.Printf("  │    └── ⚙️  %s\n", svc)
		}
	}
	fmt.Println()
	return nil
}

type moduleInfo struct {
	name        string
	controllers []string
	services    []string
}

func findModules() []moduleInfo {
	var modules []moduleInfo

	entries, err := os.ReadDir(filepath.Join("internal", "modules"))
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		mod := moduleInfo{name: entry.Name()}

		modDir := filepath.Join("internal", "modules", entry.Name())
		_ = filepath.Walk(modDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() || filepath.Ext(path) != ".go" {
				return nil
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return nil
			}
			src := string(content)

			base := filepath.Base(path)
			if strings.HasSuffix(base, "controller.go") || strings.Contains(src, "Prefix()") {
				name := strings.TrimSuffix(base, ".go")
				mod.controllers = append(mod.controllers, name)
			}
			if strings.HasSuffix(base, "service.go") {
				name := strings.TrimSuffix(base, ".go")
				mod.services = append(mod.services, name)
			}
			return nil
		})

		modules = append(modules, mod)
	}

	return modules
}


