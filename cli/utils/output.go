package utils

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// ─── Color Palette ──────────────────────────────────────────────────────────

var (
	ColorSuccess = lipgloss.Color("#22c55e") // green
	ColorError   = lipgloss.Color("#ef4444") // red
	ColorWarning = lipgloss.Color("#eab308") // yellow
	ColorInfo    = lipgloss.Color("#3b82f6") // blue
	ColorAccent  = lipgloss.Color("#06b6d4") // cyan
	ColorDim     = lipgloss.Color("#6b7280") // gray
	ColorWhite   = lipgloss.Color("#f9fafb")
	ColorPurple  = lipgloss.Color("#a855f7")
)

// ─── Styles ─────────────────────────────────────────────────────────────────

var (
	StyleSuccess = lipgloss.NewStyle().Foreground(ColorSuccess).Bold(true)
	StyleError   = lipgloss.NewStyle().Foreground(ColorError).Bold(true)
	StyleWarning = lipgloss.NewStyle().Foreground(ColorWarning).Bold(true)
	StyleInfo    = lipgloss.NewStyle().Foreground(ColorInfo).Bold(true)
	StyleAccent  = lipgloss.NewStyle().Foreground(ColorAccent)
	StyleDim     = lipgloss.NewStyle().Foreground(ColorDim)
	StyleBold    = lipgloss.NewStyle().Bold(true)
	StyleHeader  = lipgloss.NewStyle().
			Foreground(ColorWhite).
			Bold(true).
			Padding(0, 1)
	StyleBanner = lipgloss.NewStyle().
			Foreground(ColorAccent).
			Bold(true)
	StyleTableHeader = lipgloss.NewStyle().
				Foreground(ColorAccent).
				Bold(true).
				Underline(true)
	StyleGroupTitle = lipgloss.NewStyle().
			Foreground(ColorPurple).
			Bold(true).
			MarginTop(1)
)

// ─── Print Functions ────────────────────────────────────────────────────────

// PrintSuccess prints a green success message with checkmark.
func PrintSuccess(msg string) {
	fmt.Println(StyleSuccess.Render("  ✔ " + msg))
}

// PrintError prints a red error message with cross mark.
func PrintError(msg string) {
	fmt.Fprintln(os.Stderr, StyleError.Render("  ✖ "+msg))
}

// PrintWarning prints a yellow warning message.
func PrintWarning(msg string) {
	fmt.Println(StyleWarning.Render("  ⚠ " + msg))
}

// PrintInfo prints a blue informational message.
func PrintInfo(msg string) {
	fmt.Println(StyleInfo.Render("  ℹ " + msg))
}

// PrintStep prints a step with an icon and message.
func PrintStep(icon, msg string) {
	fmt.Printf("  %s %s\n", icon, msg)
}

// PrintDim prints dimmed/muted text.
func PrintDim(msg string) {
	fmt.Println(StyleDim.Render("  " + msg))
}

// ─── Header / Banner ────────────────────────────────────────────────────────

// PrintHeader prints a styled boxed header.
func PrintHeader(title string) {
	width := len(title) + 8
	if width < 44 {
		width = 44
	}

	border := strings.Repeat("═", width-2)
	padding := (width - 2 - len(title)) / 2
	paddedTitle := strings.Repeat(" ", padding) + title + strings.Repeat(" ", width-2-padding-len(title))

	box := fmt.Sprintf(
		"╔%s╗\n║%s║\n╚%s╝",
		border,
		paddedTitle,
		border,
	)

	fmt.Println()
	fmt.Println(StyleBanner.Render(box))
	fmt.Println()
}

// PrintBanner prints the NestGo CLI banner.
func PrintBanner() {
	PrintHeader("🚀 NestGo CLI — Enterprise-grade Go Framework")
}

// ─── Table Rendering ────────────────────────────────────────────────────────

// PrintTable renders a formatted table with headers and rows.
func PrintTable(headers []string, rows [][]string) {
	if len(headers) == 0 {
		return
	}

	// Calculate column widths.
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if i < len(widths) && len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Add padding.
	for i := range widths {
		widths[i] += 2
	}

	// Print header.
	headerLine := "  "
	for i, h := range headers {
		headerLine += fmt.Sprintf("%-*s", widths[i], h)
	}
	fmt.Println(StyleTableHeader.Render(headerLine))

	// Print separator.
	sepLine := "  "
	for _, w := range widths {
		sepLine += strings.Repeat("─", w)
	}
	fmt.Println(StyleDim.Render(sepLine))

	// Print rows.
	for _, row := range rows {
		line := "  "
		for i, cell := range row {
			if i < len(widths) {
				line += fmt.Sprintf("%-*s", widths[i], cell)
			}
		}
		fmt.Println(line)
	}
}

// ─── Grouped Help ───────────────────────────────────────────────────────────

// CommandGroup represents a group of CLI commands for help display.
type CommandGroup struct {
	Title    string
	Commands []CommandEntry
}

// CommandEntry is a single command in a help group.
type CommandEntry struct {
	Name        string
	Description string
}

// PrintGroupedHelp renders grouped command listings.
func PrintGroupedHelp(groups []CommandGroup) {
	for _, group := range groups {
		fmt.Println(StyleGroupTitle.Render("  " + group.Title))
		for _, cmd := range group.Commands {
			name := StyleAccent.Render(fmt.Sprintf("    %-24s", cmd.Name))
			desc := StyleDim.Render(cmd.Description)
			fmt.Printf("%s%s\n", name, desc)
		}
	}
	fmt.Println()
}

// ─── Spinner ────────────────────────────────────────────────────────────────

// Spinner shows a loading spinner with a message.
type Spinner struct {
	message string
	stop    chan struct{}
	done    chan struct{}
	mu      sync.Mutex
	stopped bool
}

// StartSpinner creates and starts a new spinner.
func StartSpinner(message string) *Spinner {
	s := &Spinner{
		message: message,
		stop:    make(chan struct{}),
		done:    make(chan struct{}),
	}
	go s.run()
	return s
}

func (s *Spinner) run() {
	defer close(s.done)

	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	i := 0
	ticker := time.NewTicker(80 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-s.stop:
			// Clear the spinner line.
			fmt.Printf("\r%s\r", strings.Repeat(" ", len(s.message)+10))
			return
		case <-ticker.C:
			frame := StyleAccent.Render(frames[i%len(frames)])
			fmt.Printf("\r  %s %s", frame, s.message)
			i++
		}
	}
}

// Stop stops the spinner and prints a completion message.
func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.stopped {
		return
	}
	s.stopped = true
	close(s.stop)
	<-s.done
}

// StopWithSuccess stops the spinner and prints a success message.
func (s *Spinner) StopWithSuccess(msg string) {
	s.Stop()
	PrintSuccess(msg)
}

// StopWithError stops the spinner and prints an error message.
func (s *Spinner) StopWithError(msg string) {
	s.Stop()
	PrintError(msg)
}

// ─── Progress Bar ───────────────────────────────────────────────────────────

// PrintProgress prints a simple progress indicator.
func PrintProgress(current, total int, label string) {
	width := 30
	filled := (current * width) / total
	bar := strings.Repeat("█", filled) + strings.Repeat("░", width-filled)
	pct := (current * 100) / total
	fmt.Printf("\r  %s [%s] %d%%", label, StyleAccent.Render(bar), pct)
	if current == total {
		fmt.Println()
	}
}

// ─── Confirmation ───────────────────────────────────────────────────────────

// Confirm prompts the user for a yes/no confirmation.
func Confirm(prompt string) bool {
	fmt.Printf("  %s [y/N]: ", prompt)
	var input string
	_, _ = fmt.Scanln(&input)
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}
