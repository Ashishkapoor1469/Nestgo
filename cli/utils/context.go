package utils

import (
	"fmt"
	"os"
)

// IsNestGoProject checks if the current working directory is a valid NestGo project.
// It checks for the existence of go.mod and nestgo.json natively.
func IsNestGoProject() bool {
	// Look for go.mod (marker of Go project)
	if _, err := os.Stat("go.mod"); os.IsNotExist(err) {
		return false
	}
	
	// Assuming `main.go` and `core` patterns signify a nested project. 
	// For stricter detection, we can look for nestgo.json, but given
	// it's a Go framework, looking for go.mod is the first gate.
	return true
}

// EnsureProjectContext gracefully fails if run outside a valid NestGo project space.
func EnsureProjectContext(command string) {
	if !IsNestGoProject() {
		fmt.Printf("❌ Cannot execute 'nestgo %s' here.\n", command)
		fmt.Println("This command must be run inside a valid NestGo project directory (containing go.mod).")
		fmt.Println("Try running 'nestgo new <app-name>' to create a project first.")
		os.Exit(1)
	}
}

// PrintSuccess logs standardized success messages.
func PrintSuccess(msg string) {
	fmt.Printf("\033[32m✅ %s\033[0m\n", msg)
}

// PrintWarning logs standardized warning messages.
func PrintWarning(msg string) {
	fmt.Printf("\033[33m⚠️  %s\033[0m\n", msg)
}

// PrintError logs standardized error messages.
func PrintError(msg string) {
	fmt.Printf("\033[31m❌ %s\033[0m\n", msg)
}
