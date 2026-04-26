package commands

import (
	"bytes"
	"testing"
)

func TestNewCommand_MissingArg(t *testing.T) {
	cmd := NewCmd()

	var buf bytes.Buffer
	cmd.SetErr(&buf)
	cmd.SetOut(&buf)

	// No arguments provided
	cmd.SetArgs([]string{})

	err := cmd.Execute()

	// Ensure we got an error and not a panic
	if err == nil {
		t.Fatal("Expected error when no arguments are provided, got nil")
	}

	// Wait, cobra defaults to showing help on Missing args, or it throws an error.
	// We want to ensure it handles it gracefully.
}
