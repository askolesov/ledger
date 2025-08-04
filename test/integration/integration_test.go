package integration

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	binaryName = "ledger"
	binaryPath = "../../build/ledger"
)

// Helper function to run CLI command and capture output
func runCommand(t *testing.T, args ...string) (stdout, stderr string, exitCode int) {
	t.Helper()

	// Get absolute path to binary
	absPath, err := filepath.Abs(binaryPath)
	if err != nil {
		t.Fatalf("Failed to get absolute path to binary: %v", err)
	}

	// Check if binary exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		t.Fatalf("Binary not found at %s. Run 'make build' first.", absPath)
	}

	cmd := exec.Command(absPath, args...)

	// Use CombinedOutput to get both stdout and stderr together
	output, err := cmd.CombinedOutput()

	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			exitCode = exitError.ExitCode()
		} else {
			t.Fatalf("Failed to run command: %v", err)
		}
	} else {
		exitCode = 0
	}

	// For simplicity, put all output in stdout since CombinedOutput mixes them
	stdout = string(output)
	stderr = ""
	return stdout, stderr, exitCode
}

// Helper function to get test data file path
func getTestDataPath(filename string) string {
	return filepath.Join("testdata", filename)
}

// V1 Tests
func TestV1ValidateValid(t *testing.T) {
	stdout, stderr, exitCode := runCommand(t, "v1", "validate", getTestDataPath("v1/valid.yaml"))

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
	}

	if !strings.Contains(stdout, "✓ Data is valid") {
		t.Errorf("Expected success message in stdout, got: %s", stdout)
	}
}

func TestV1ValidateInvalidBalance(t *testing.T) {
	stdout, stderr, exitCode := runCommand(t, "v1", "validate", getTestDataPath("v1/invalid-balance.yaml"))

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code for invalid file, got 0")
	}

	// Should contain validation error message
	output := stdout + stderr
	if !strings.Contains(output, "validation failed") {
		t.Errorf("Expected validation error message, got: %s", output)
	}
}

func TestV1ValidateInvalidStructure(t *testing.T) {
	stdout, stderr, exitCode := runCommand(t, "v1", "validate", getTestDataPath("v1/invalid-structure.yaml"))

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code for invalid file, got 0")
	}

	// Should contain error message (either parsing or validation)
	output := stdout + stderr
	if !strings.Contains(output, "failed") && !strings.Contains(output, "error") {
		t.Errorf("Expected error message, got: %s", output)
	}
}

// V2 Tests
func TestV2ValidateValid(t *testing.T) {
	stdout, stderr, exitCode := runCommand(t, "validate", getTestDataPath("v2/valid.yaml"))

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
	}

	if !strings.Contains(stdout, "✓ Ledger is valid") {
		t.Errorf("Expected success message in stdout, got: %s", stdout)
	}
}

func TestV2ValidateInvalidBalance(t *testing.T) {
	stdout, stderr, exitCode := runCommand(t, "validate", getTestDataPath("v2/invalid-balance.yaml"))

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code for invalid file, got 0")
	}

	// Should contain validation error message
	output := stdout + stderr
	if !strings.Contains(output, "validation failed") {
		t.Errorf("Expected validation error message, got: %s", output)
	}
}

func TestV2ValidateInvalidStructure(t *testing.T) {
	stdout, stderr, exitCode := runCommand(t, "validate", getTestDataPath("v2/invalid-structure.yaml"))

	if exitCode == 0 {
		t.Errorf("Expected non-zero exit code for invalid file, got 0")
	}

	// Should contain error message (either parsing or validation)
	output := stdout + stderr
	if !strings.Contains(output, "failed") && !strings.Contains(output, "error") {
		t.Errorf("Expected error message, got: %s", output)
	}
}

// Migration Test
func TestV1MigrateToV2(t *testing.T) {
	// Create temporary output file
	outputFile := filepath.Join("testdata", "migration", "output.yaml")
	defer func() {
		require.NoError(t, os.Remove(outputFile)) // Clean up after test
	}()

	// Run migration
	stdout, stderr, exitCode := runCommand(t, "v1", "migrate",
		getTestDataPath("migration/v1-source.yaml"),
		outputFile)

	if exitCode != 0 {
		t.Errorf("Expected exit code 0, got %d. Stderr: %s", exitCode, stderr)
	}

	if !strings.Contains(stdout, "✓ Migration completed successfully") {
		t.Errorf("Expected success message in stdout, got: %s", stdout)
	}

	// Verify output file was created
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Errorf("Output file was not created: %s", outputFile)
		return
	}

	// Validate the migrated file using v2 validate command
	stdout2, stderr2, exitCode2 := runCommand(t, "validate", outputFile)

	if exitCode2 != 0 {
		t.Errorf("Migrated file failed v2 validation. Exit code: %d, Stderr: %s", exitCode2, stderr2)
	}

	if !strings.Contains(stdout2, "✓ Ledger is valid") {
		t.Errorf("Expected migrated file to be valid, got: %s", stdout2)
	}
}
