package main

import (
	"fmt"
	"os"

	"github.com/nestgo/nestgo/cli/commands"
	"github.com/spf13/cobra"
)

var version = "0.1.0"

func main() {
	rootCmd := &cobra.Command{
		Use:   "nestgo",
		Short: "NestGo — Production-grade Go backend framework",
		Long: `
    ╔══════════════════════════════════════════╗
    ║           🚀 NestGo Framework            ║
    ║   Production-grade Go Backend Framework   ║
    ╚══════════════════════════════════════════╝

NestGo is a modular, high-performance backend framework for Go,
inspired by NestJS but redesigned with idiomatic Go patterns.

Get started:
  nestgo new my-app         Create a new project
  nestgo generate resource  Generate a full CRUD resource
  nestgo dev                Start development server
  nestgo doctor             Analyze project health`,
		Version: version,
	}

	// Register all commands.
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
