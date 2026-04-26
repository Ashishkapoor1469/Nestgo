package commands

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/Ashishkapoor1469/Nestgo/cli/utils"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

// DevCmd creates the `nestgo dev` command.
func DevCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "dev",
		Short: "Start development server with hot reload",
		Long:  "Watches for file changes and automatically rebuilds and restarts the server.",
		RunE:  runDev,
	}
	cmd.Flags().StringP("port", "p", "", "Override port (default from .env or 3000)")
	return cmd
}

// BuildCmd creates the `nestgo build` command.
func BuildCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "build",
		Short: "Build optimized production binary",
		RunE:  runBuild,
	}
	cmd.Flags().StringP("output", "o", "bin/server", "Output binary path")
	return cmd
}

func runDev(cmd *cobra.Command, args []string) error {
	utils.EnsureProjectContext("dev")

	portOverride, _ := cmd.Flags().GetString("port")

	// ── Startup Banner ──────────────────────────────────────────────────────
	fmt.Println()
	fmt.Println("  ┌─────────────────────────────────────────┐")
	fmt.Println("  │         NestGo Dev Server               │")
	fmt.Println("  │         Hot Reload Enabled              │")
	fmt.Println("  └─────────────────────────────────────────┘")
	fmt.Println()

	// ── Validate project ────────────────────────────────────────────────────
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return fmt.Errorf("no go.mod found — are you in a NestGo project directory?\n  → Run: nestgo new <name>")
	}

	entryPoint := findEntryPoint()
	if entryPoint == "" {
		return fmt.Errorf(
			"could not find main.go entry point\n" +
				"  Searched: cmd/server/main.go, cmd/main.go, main.go\n" +
				"  → Run: nestgo new <name> to scaffold a project",
		)
	}

	fmt.Printf("  ✔ NestGo project detected\n")
	fmt.Printf("  ✔ Entry point: %s\n", entryPoint)

	// ── Port resolution ──────────────────────────────────────────────────────
	port := resolvePort(portOverride)
	if isPortInUse(port) {
		suggested := findFreePort(port)
		fmt.Printf("  ⚠️  Port %s is in use\n", port)
		if suggested != "" {
			fmt.Printf("  → Using port %s instead\n", suggested)
			port = suggested
		} else {
			return fmt.Errorf("port %s is busy and no alternative found", port)
		}
	}
	fmt.Printf("  ✔ Port: %s\n", port)

	// Inject PORT so the app uses our chosen port
	if err := os.Setenv("PORT", port); err != nil {
		return fmt.Errorf("failed to set PORT: %w", err)
	}

	fmt.Printf("  ✔ Watching for file changes...\n\n")

	// ── Start process ────────────────────────────────────────────────────────
	var process *exec.Cmd
	restart := make(chan struct{}, 1)

	process = startDevProcess(entryPoint, port)

	// ── File watcher ─────────────────────────────────────────────────────────
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			fmt.Printf("  ❌ Watcher error: %v\n", err)
			return
		}
		defer func() { _ = watcher.Close() }()

		// Watch all directories recursively, excluding noise dirs
		_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				skip := map[string]bool{
					"vendor": true, "node_modules": true,
					".git": true, "bin": true, "dist": true,
					".next": true, "website": true,
				}
				if skip[filepath.Base(path)] {
					return filepath.SkipDir
				}
				return watcher.Add(path)
			}
			return nil
		})

		var lastEvent time.Time
		debounce := 400 * time.Millisecond

		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) || event.Has(fsnotify.Create) {
					ext := filepath.Ext(event.Name)
					if ext == ".go" || ext == ".env" {
						now := time.Now()
						if now.Sub(lastEvent) > debounce {
							lastEvent = now
							rel, _ := filepath.Rel(".", event.Name)
							fmt.Printf("\n  🔄 Changed: %s\n", rel)
							select {
							case restart <- struct{}{}:
							default:
							}
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Printf("  ⚠️  Watch error: %v\n", err)
			}
		}
	}()

	// ── Restart handler ──────────────────────────────────────────────────────
	go func() {
		for range restart {
			fmt.Println("  ⏳ Rebuilding...")
			if process != nil && process.Process != nil {
				_ = process.Process.Signal(syscall.SIGTERM)
				_ = process.Wait()
			}
			process = startDevProcess(entryPoint, port)
		}
	}()

	// ── Graceful shutdown ────────────────────────────────────────────────────
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n  🛑 Shutting down dev server...")
	if process != nil && process.Process != nil {
		_ = process.Process.Signal(syscall.SIGTERM)
		_ = process.Wait()
	}
	fmt.Println("  ✔ Goodbye!")
	return nil
}

func runBuild(cmd *cobra.Command, args []string) error {
	utils.EnsureProjectContext("build")
	output, _ := cmd.Flags().GetString("output")

	fmt.Println("\n  🔨 Building NestGo application...")

	entryPoint := findEntryPoint()
	if entryPoint == "" {
		return fmt.Errorf("could not find main.go entry point")
	}

	if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
		return err
	}

	start := time.Now()
	buildCmd := exec.Command("go", "build",
		"-ldflags=-s -w",
		"-o", output,
		"./"+filepath.Dir(entryPoint),
	)
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("build failed: %w", err)
	}

	duration := time.Since(start).Round(time.Millisecond)
	info, _ := os.Stat(output)
	size := "unknown"
	if info != nil {
		size = fmt.Sprintf("%.1f MB", float64(info.Size())/1024/1024)
	}

	fmt.Printf("  ✅ Built in %s\n", duration)
	fmt.Printf("  📦 Binary: %s (%s)\n", output, size)
	fmt.Printf("  🚀 Run: ./%s\n\n", output)
	return nil
}

// findEntryPoint searches common locations for main.go.
func findEntryPoint() string {
	candidates := []string{
		"cmd/server/main.go",
		"cmd/main.go",
		"main.go",
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			return c
		}
	}
	return ""
}

// startDevProcess compiles and runs the app, printing output with a prefix.
func startDevProcess(entryPoint, port string) *exec.Cmd {
	cmd := exec.Command("go", "run", "./"+filepath.Dir(entryPoint))
	cmd.Env = append(os.Environ(), "PORT="+port)

	// Pipe stdout/stderr with a nice prefix
	stdoutPipe, _ := cmd.StdoutPipe()
	stderrPipe, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		fmt.Printf("  ❌ Failed to start: %v\n", err)
		// Give a useful hint for common errors
		if strings.Contains(err.Error(), "exec") {
			fmt.Println("  → Make sure Go is installed and in PATH")
		}
		return nil
	}

	// Stream stdout
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			fmt.Printf("  │ %s\n", scanner.Text())
		}
	}()

	// Stream stderr — detect compile errors
	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "syntax error") || strings.Contains(line, "undefined") {
				fmt.Printf("  ❌ %s\n", line)
			} else if strings.Contains(line, "cannot use") || strings.Contains(line, "does not implement") {
				fmt.Printf("  ❌ %s\n", line)
			} else {
				fmt.Printf("  │ %s\n", line)
			}
		}
	}()

	fmt.Printf("  ✔ Server started (PID: %d) → http://localhost:%s\n", cmd.Process.Pid, port)
	return cmd
}

// resolvePort determines port from override flag or .env file.
func resolvePort(override string) string {
	if override != "" {
		return override
	}
	// Try reading from .env
	if data, err := os.ReadFile(".env"); err == nil {
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "PORT=") {
				if p := strings.TrimPrefix(line, "PORT="); p != "" {
					return p
				}
			}
		}
	}
	return "3000"
}

// isPortInUse checks if a TCP port is already bound.
func isPortInUse(port string) bool {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return true
	}
	_ = ln.Close()
	return false
}

// findFreePort tries incrementing the port until a free one is found.
func findFreePort(startPort string) string {
	base := 3000
	_, _ = fmt.Sscanf(startPort, "%d", &base)
	for i := 1; i <= 10; i++ {
		p := fmt.Sprintf("%d", base+i)
		if !isPortInUse(p) {
			return p
		}
	}
	return ""
}
