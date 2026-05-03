package commands

import (
	"runtime"

	"github.com/spf13/cobra"
)

// VersionCmd creates the `nestgo version` command.
func VersionCmd(cliVersion string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show NestGo version information",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println()
			cmd.Println("  NestGo CLI")
			cmd.Println("  ─────────────────────────────────")
			cmd.Printf("  CLI Version:       v%s\n", cliVersion)
			cmd.Printf("  Framework Version: v0.5.0\n")
			cmd.Printf("  Go Version:        %s\n", runtime.Version())
			cmd.Printf("  OS/Arch:           %s/%s\n", runtime.GOOS, runtime.GOARCH)
			cmd.Println("  ─────────────────────────────────")
			cmd.Println()
		},
	}
}
