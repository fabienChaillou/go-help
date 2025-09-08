// main_test.go
package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestRealMain(t *testing.T) {
	// Redirect stdout to capture output
	oldStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// Call the function we want to test
	exitCode := realMain()

	// Restore stdout
	w.Close()
	os.Stdout = oldStdout

	// Read captured output
	var buf bytes.Buffer
	buf.ReadFrom(r)
	output := buf.String()

	// Check exit code
	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d", exitCode)
	}

	// Check output
	expectedOutput := "Hello, World!"
	if !strings.Contains(output, expectedOutput) {
		t.Errorf("Expected output to contain %q, got %q", expectedOutput, output)
	}
}

func TestMain(m *testing.M) {
	// TestMain is a special function that allows setup and teardown for tests
	// Run all the tests
	exitCode := m.Run()

	// Exit with the same code
	os.Exit(exitCode)
}
