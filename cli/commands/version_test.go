package commands

import (
	"bytes"
	"strings"
	"testing"
)

func TestVersionCommand(t *testing.T) {
	cmd := VersionCmd("1.0.0-test")

	// Capture output
	var buf bytes.Buffer
	cmd.SetOut(&buf)

	// Execute the command without args
	cmd.SetArgs([]string{})
	err := cmd.Execute()

	if err != nil {
		t.Fatalf("VersionCmd execution failed: %v", err)
	}

	output := buf.String()

	if !strings.Contains(output, "NestGo CLI") {
		t.Errorf("Expected output to contain 'NestGo CLI', got:\n%s", output)
	}
	if !strings.Contains(output, "1.0.0-test") {
		t.Errorf("Expected output to contain '1.0.0-test', got:\n%s", output)
	}
}
