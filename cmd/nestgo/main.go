package main

import (
	"fmt"
	"os"

	"github.com/Ashishkapoor1469/Nestgo/cli/commands"
	"github.com/Ashishkapoor1469/Nestgo/cli/config"
	"github.com/spf13/cobra"
)

// Version is injected dynamically at build time using -ldflags="-X main.Version=1.0.0"
var Version = "0.1.0" // Default fallback

func main() {
	// Initialize Global Config (~/.nestgo) behind the scenes
	_, err := config.LoadGlobalConfig()
	if err != nil {
		// Just warn, don't crash
		fmt.Printf("\033[33m⚠️  Warning: Could not load global config: %v\033[0m\n", err)
	}

	rootCmd := &cobra.Command{
		Use:   "nestgo",
		Short: "NestGo — Enterprise-grade Go backend CLI",
		Long: `
    ╔══════════════════════════════════════════╗
    ║           🚀 NestGo CLI Engine           ║
    ║   Enterprise-grade Go Backend Framework  ║
    ╚══════════════════════════════════════════╝

NestGo CLI accelerates project scaffolding, module generations,
and runs hot-reloading development servers seamlessly.

Commands:
  nestgo new <name>         Create a new production project
  nestgo generate resource  Scaffold a complete CRUD module
  nestgo dev                Start hot-reloading development server
  nestgo doctor             Run architecture health checks

Autocomplete:
  nestgo completion bash    Generate bash autocompletion script
  nestgo completion zsh     Generate zsh autocompletion script
`,
		Version: Version,
	}

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
		commands.MigrateCmd(),
	)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
