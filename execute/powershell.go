package execute

import (
	logger "calderat/utils"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type PowerShell struct {
	shortName string   // Friendly name or alias for the shell
	path      string   // Path to the PowerShell executable
	execArgs  []string // Default arguments for the shell execution
	logger    *logger.Logger
}

// NewPowerShell initializes a new PowerShell instance
func NewPowerShell(log *logger.Logger) *PowerShell {
	return &PowerShell{
		shortName: "PowerShell",
		path:      "powershell", // Default PowerShell executable
		execArgs:  []string{"-ExecutionPolicy", "Bypass", "-Command"},
		logger:    log,
	}
}

// Execute runs a PowerShell command with a specified timeout and logs the process
func (ps *PowerShell) Execute(command string, timeout time.Duration) (string, error) {
	if runtime.GOOS != "windows" {
		ps.logger.Log(logger.ERROR, "Command execution failed: Unsupported OS")
		return "", fmt.Errorf("%s is only supported on Windows systems", ps.shortName)
	}

	ps.logger.Log(logger.INFO, "Executing command: %s (Timeout: %v)", command, timeout)
	ps.logger.Log(logger.TRACE, "Full arguments: %v", ps.execArgs)

	// Create a context with the specified timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Append the command to the default execution arguments
	args := append(ps.execArgs, command)

	// Construct the execution command
	cmd := exec.CommandContext(ctx, ps.path, args...)

	// Capture the output
	output, err := cmd.CombinedOutput()

	// Check for context timeout
	if ctx.Err() == context.DeadlineExceeded {
		ps.logger.Log(logger.ERROR, "Command timed out after %v", timeout)
		return "", fmt.Errorf("command timed out after %v", timeout)
	}

	if err != nil {
		ps.logger.Log(logger.ERROR, "Command execution failed: %v\nOutput: %s", err, string(output))
		return "", fmt.Errorf("failed to execute %s command: %v\nOutput: %s", ps.shortName, err, string(output))
	}

	ps.logger.Log(logger.DEBUG, "Command executed successfully. Output:\n%s", strings.TrimRight(string(output), " \n\r"))
	return strings.TrimRight(string(output), " \n\r"), nil
}

func (p *PowerShell) ShortName() string {
	return p.shortName
}

func (p *PowerShell) Path() string {
	return p.path
}
