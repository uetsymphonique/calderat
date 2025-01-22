package execute_test

import (
	"calderat/execute"
	"calderat/utils/logger"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

func TestNewPowerShell(t *testing.T) {
	log, _ := logger.New("DEBUG") // Initialize a test logger
	ps := execute.NewPowerShell(log)

	if ps == nil {
		t.Fatal("Expected NewPowerShell to return a valid instance, got nil")
	}
	if ps.ShortName() != "PowerShell" {
		t.Errorf("Expected shortName to be 'PowerShell', got '%s'", ps.ShortName())
	}
	if ps.Path() != "powershell" {
		t.Errorf("Expected path to be 'powershell', got '%s'", ps.Path())
	}
}

func TestPowerShellExecute(t *testing.T) {
	logFile, err := os.CreateTemp("", "powershell_test.log")
	if err != nil {
		t.Fatalf("Failed to create temporary log file: %v", err)
	}
	defer os.Remove(logFile.Name())

	log, _ := logger.NewWithOutput("DEBUG", logFile) // Initialize logger with custom output
	ps := execute.NewPowerShell(log)

	// Case 1: Unsupported OS
	if runtime.GOOS != "windows" {
		_, err := ps.Execute("Get-Date", 2*time.Second)
		if err == nil || err.Error() != "PowerShell is only supported on Windows systems" {
			t.Error("Expected unsupported OS error, got nil or wrong error")
		}
		return
	}

	// Case 2: Successful Command Execution
	output, err := ps.Execute("Write-Output 'Hello, PowerShell!'", 2*time.Second)
	if err != nil {
		t.Errorf("Expected successful execution, got error: %v", err)
	}
	fmt.Println(output)

	expected := "Hello, PowerShell!"
	fmt.Println(expected)
	if output != expected {
		t.Errorf("Expected output: '%s', got: '%s'", expected, output)
	}

	// Case 3: Timeout Scenario
	_, err = ps.Execute("Start-Sleep -Seconds 5", 2*time.Second)
	if err == nil || !strings.Contains(err.Error(), "command timed out") {
		t.Error("Expected timeout error, got nil or wrong error")
	}

	// Verify Logs
	logContent, err := os.ReadFile(logFile.Name())
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	logText := string(logContent)
	if !strings.Contains(logText, "Executing command: Write-Output 'Hello, PowerShell!'") {
		t.Error("Expected log entry for command execution not found")
	}
}

func TestPowerShellLogger(t *testing.T) {
	log, _ := logger.New("DEBUG")
	ps := execute.NewPowerShell(log)

	// Test logging interactions
	ps.Execute("echo test", 2*time.Second)
}

func TestPowerShellTimeout(t *testing.T) {
	log, _ := logger.New("DEBUG")
	ps := execute.NewPowerShell(log)

	// Test timeout logic
	_, err := ps.Execute("Start-Sleep -Seconds 5", 1*time.Second)
	if err == nil || !strings.Contains(err.Error(), "command timed out") {
		t.Error("Expected timeout error, got nil or wrong error")
	}
}
