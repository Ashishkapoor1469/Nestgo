package commands

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/nestgo/nestgo/cli/utils"
	"github.com/spf13/cobra"
)

// DevCmd creates the `nestgo dev` command.
func DevCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "dev",
		Short: "Start development server with hot reload",
		Long:  "Watches for file changes and automatically rebuilds and restarts the server.",
		RunE:  runDev,
	}
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
	fmt.Println("\n🔥 NestGo Dev Server — Hot Reload Mode")

	// Find the main entry point.
	entryPoint := findEntryPoint()
	if entryPoint == "" {
		return fmt.Errorf("could not find main.go entry point (looked in cmd/server/main.go, cmd/main.go, main.go)")
	}

	fmt.Printf("  📂 Entry: %s\n", entryPoint)
	fmt.Println("  👀 Watching for changes...")

	// Start the process.
	var process *exec.Cmd
	restart := make(chan struct{}, 1)

	// Initial start.
	process = startProcess(entryPoint)

	// Watch for file changes.
	go func() {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			fmt.Printf("❌ Watch error: %v\n", err)
			return
		}
		defer watcher.Close()

		// Watch all .go files recursively.
		_ = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				if path == "vendor" || path == "node_modules" || path == ".git" || path == "bin" {
					return filepath.SkipDir
				}
				return watcher.Add(path)
			}
			return nil
		})

		var lastEvent time.Time
		debounce := 500 * time.Millisecond

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
							fmt.Printf("\n  🔄 Change detected: %s\n", event.Name)
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

	// Handle restarts.
	go func() {
		for range restart {
			fmt.Println("  ⏳ Rebuilding...")
			if process != nil && process.Process != nil {
				_ = process.Process.Signal(syscall.SIGTERM)
				_ = process.Wait()
			}
			process = startProcess(entryPoint)
		}
	}()

	// Wait for interrupt.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\n  🛑 Stopping dev server...")
	if process != nil && process.Process != nil {
		_ = process.Process.Signal(syscall.SIGTERM)
		_ = process.Wait()
	}

	return nil
}

func runBuild(cmd *cobra.Command, args []string) error {
	utils.EnsureProjectContext("build")
	output, _ := cmd.Flags().GetString("output")

	fmt.Println("\n🔨 Building NestGo application...")

	entryPoint := findEntryPoint()
	if entryPoint == "" {
		return fmt.Errorf("could not find main.go entry point")
	}

	// Ensure output directory exists.
	dir := filepath.Dir(output)
	if err := os.MkdirAll(dir, 0755); err != nil {
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

	duration := time.Since(start)

	// Get binary size.
	info, _ := os.Stat(output)
	size := "unknown"
	if info != nil {
		sizeMB := float64(info.Size()) / 1024 / 1024
		size = fmt.Sprintf("%.1f MB", sizeMB)
	}

	fmt.Printf("  ✅ Build complete in %s\n", duration)
	fmt.Printf("  📦 Binary: %s (%s)\n", output, size)
	fmt.Printf("  🚀 Run: ./%s\n\n", output)

	return nil
}

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

func startProcess(entryPoint string) *exec.Cmd {
	cmd := exec.Command("go", "run", "./"+filepath.Dir(entryPoint))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		fmt.Printf("  ❌ Failed to start: %v\n", err)
		return nil
	}

	fmt.Printf("  ✅ Server running (PID: %d)\n", cmd.Process.Pid)
	return cmd
}
