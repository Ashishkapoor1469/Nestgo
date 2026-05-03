package main

import (
	"fmt"
	"os"

	"github.com/Ashishkapoor1469/Nestgo/cli/commands"
	"github.com/Ashishkapoor1469/Nestgo/cli/config"
	"github.com/Ashishkapoor1469/Nestgo/cli/utils"
	"github.com/spf13/cobra"
)

// Version is injected dynamically at build time using -ldflags="-X main.Version=1.0.0"
var Version = "0.5.0" // Default fallback

func main() {
	// Initialize Global Config (~/.nestgo) behind the scenes
	_, err := config.LoadGlobalConfig()
	if err != nil {
		// Just warn, don't crash
		utils.PrintWarning("Could not load global config: " + err.Error())
	}

	rootCmd := &cobra.Command{
		Use:     "nestgo",
		Short:   "NestGo — Enterprise-grade Go backend CLI",
		Version: Version,
	}

	// Override default help output with grouped, styled display.
	rootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		utils.PrintBanner()

		fmt.Printf("  %s\n\n", utils.StyleDim.Render("Enterprise-grade Go Backend Framework — v"+Version))

		utils.PrintGroupedHelp([]utils.CommandGroup{
			{
				Title: "📦 Project",
				Commands: []utils.CommandEntry{
					{Name: "new <name>", Description: "Create a new NestGo project"},
				},
			},
			{
				Title: "⚙️  Generation",
				Commands: []utils.CommandEntry{
					{Name: "generate resource", Description: "Scaffold complete CRUD module"},
					{Name: "generate module", Description: "Generate a module"},
					{Name: "generate controller", Description: "Generate a controller"},
					{Name: "generate service", Description: "Generate a service"},
					{Name: "generate schema", Description: "Generate a validation schema"},
					{Name: "generate dto", Description: "Generate DTOs (Create/Update)"},
					{Name: "generate test", Description: "Generate test file"},
					{Name: "generate middleware", Description: "Generate middleware"},
					{Name: "generate guard", Description: "Generate auth guard"},
					{Name: "generate interceptor", Description: "Generate interceptor"},
				},
			},
			{
				Title: "🗄️  Database",
				Commands: []utils.CommandEntry{
					{Name: "migration create", Description: "Create new migration file"},
					{Name: "migration run", Description: "Run pending migrations"},
					{Name: "migration rollback", Description: "Rollback last migration"},
					{Name: "migration status", Description: "Show migration status"},
				},
			},
			{
				Title: "🚀 Development",
				Commands: []utils.CommandEntry{
					{Name: "dev", Description: "Start hot-reloading dev server"},
					{Name: "build", Description: "Build optimized production binary"},
					{Name: "routes", Description: "Display all registered REST API routes"},
					{Name: "docs:generate", Description: "Generate OpenAPI/Swagger documentation"},
				},
			},
			{
				Title: "🔍 Analysis",
				Commands: []utils.CommandEntry{
					{Name: "doctor", Description: "Run project health checks"},
					{Name: "graph", Description: "Visualize module dependencies"},
					{Name: "lint-arch", Description: "Check clean architecture violations"},
				},
			},
			{
				Title: "📖 Info / Guides",
				Commands: []utils.CommandEntry{
					{Name: "version", Description: "Show CLI, framework, and Go version"},
					{Name: "metrics", Description: "Learn how to enable the metrics endpoint"},
					{Name: "versioning", Description: "Learn how to use API versioning"},
				},
			},
		})

		fmt.Printf("  %s\n", utils.StyleDim.Render("Use \"nestgo <command> --help\" for more information about a command."))
		fmt.Println()
	})

	// Make sure completion is enabled completely.
	rootCmd.CompletionOptions.DisableDefaultCmd = false

	// Register all subcommands.
	rootCmd.AddCommand(
		commands.NewCmd(),
		commands.GenerateCmd(),
		commands.DevCmd(),
		commands.BuildCmd(),
		commands.DoctorCmd(),
		commands.GraphCmd(),
		commands.MigrationCmd(),
		commands.RoutesCmd(),
		commands.LintArchCmd(),
		commands.DocsCmd(),
		commands.MetricsCmd(),
		commands.VersioningCmd(),
		commands.VersionCmd(Version),

		// Colon-separated aliases for migration commands.
		commands.MigrationCreateAliasCmd(),
		commands.MigrationRunAliasCmd(),
		commands.MigrationRollbackAliasCmd(),
		commands.MigrationStatusAliasCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		utils.PrintError(err.Error())
		os.Exit(1)
	}
}
