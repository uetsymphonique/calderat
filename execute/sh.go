package execute

import (
	logger "calderat/utils"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

type Sh struct {
	shortName string
	logger    *logger.Logger
}

// NewSh initializes a new SH executor
func NewSh(log *logger.Logger) *Sh {
	return &Sh{
		shortName: "SH",
		logger:    log,
	}
}

// Execute runs an SH command with a specified timeout
func (se *Sh) Execute(command string, timeout time.Duration) (string, error) {
	if runtime.GOOS != "linux" {
		se.logger.Log(logger.ERROR, "Command execution failed: Unsupported OS")
		return "", fmt.Errorf("%s is only supported on Linux systems", se.shortName)
	}

	se.logger.Log(logger.INFO, "Executing command: %s (Timeout: %v)", command, timeout)

	// Create a context with the specified timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Construct the execution command
	cmd := exec.CommandContext(ctx, "sh", "-c", command)

	// Capture the output
	output, err := cmd.CombinedOutput()

	// Check for context timeout
	if ctx.Err() == context.DeadlineExceeded {
		se.logger.Log(logger.WARN, "Command timed out after %v", timeout)
		return "", fmt.Errorf("command timed out after %v", timeout)
	}

	if err != nil {
		se.logger.Log(logger.ERROR, "Command execution failed: %v\nOutput: %s", err, string(output))
		return "", fmt.Errorf("failed to execute %s command: %v\nOutput: %s", se.shortName, err, string(output))
	}

	se.logger.Log(logger.DEBUG, "Command executed successfully. Output:\n%s", string(output))
	return string(output), nil
}
