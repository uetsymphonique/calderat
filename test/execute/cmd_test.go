package execute

import (
	"calderat/service/execute"
	"calderat/utils/logger"
	"runtime"
	"testing"
	"time"
)

// TestNewCmd ensures that the Cmd struct initializes correctly

// TestExecuteSuccess verifies command execution (only runs on Windows)
func TestExecuteSuccess(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Skipping test: Execute function only works on Windows")
	}

	log, _ := logger.New("DEBUG")
	cmdExecutor := execute.NewCmd(log)

	// Simple command to test (verifying it doesn't fail)
	output, err := cmdExecutor.Execute(`dir /s c:\`, 60*time.Second)
	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}

	t.Log(output)
}
