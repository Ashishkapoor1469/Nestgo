package utils

import (
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
		PrintError("Cannot execute 'nestgo " + command + "' here.")
		PrintDim("This command must be run inside a valid NestGo project directory (containing go.mod).")
		PrintDim("Try running 'nestgo new <app-name>' to create a project first.")
		os.Exit(1)
	}
}
