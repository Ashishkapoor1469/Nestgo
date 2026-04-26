package commands

import (
	"bytes"
	"strings"
	"testing"
)

func TestGenerateCommand_InvalidType(t *testing.T) {
	cmd := GenerateCmd()

	var buf bytes.Buffer
	cmd.SetErr(&buf)
	cmd.SetOut(&buf)

	// Provide an invalid generation type
	cmd.SetArgs([]string{"invalid_type", "my_resource"})

	err := cmd.Execute()

	if err == nil {
		t.Fatal("Expected error for invalid generation type, got nil")
	}

	expectedErrorStr := "unknown command \"invalid_type\" for \"generate\""
	if !strings.Contains(err.Error(), expectedErrorStr) {
		t.Errorf("Expected error to contain %q, got %q", expectedErrorStr, err.Error())
	}
}

func TestGenerateCommand_MissingArgs(t *testing.T) {
	cmd := GenerateCmd()

	var buf bytes.Buffer
	cmd.SetErr(&buf)
	cmd.SetOut(&buf)

	// Provide no arguments
	cmd.SetArgs([]string{})

	err := cmd.Execute()

	if err == nil {
		t.Fatal("Expected error for missing arguments, got nil")
	}
}
