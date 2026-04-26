package commands

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

// VersionCmd creates the `nestgo version` command.
func VersionCmd(cliVersion string) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show NestGo version information",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println()
			fmt.Println("  NestGo CLI")
			fmt.Println("  ─────────────────────────────────")
			fmt.Printf("  CLI Version:       v%s\n", cliVersion)
			fmt.Printf("  Framework Version: v0.3.0\n")
			fmt.Printf("  Go Version:        %s\n", runtime.Version())
			fmt.Printf("  OS/Arch:           %s/%s\n", runtime.GOOS, runtime.GOARCH)
			fmt.Println("  ─────────────────────────────────")
			fmt.Println()
		},
	}
}
