package execute

import (
	"calderat/utils/colorprint"
	"calderat/utils/logger"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"time"
)

type Sh struct {
	shortName string
	logger    *logger.Logger
	path      string
}

// NewSh initializes a new SH executor
func NewSh(log *logger.Logger) *Sh {
	return &Sh{
		shortName: "SH",
		logger:    log,
		path:      "sh",
	}
}

// Execute runs an SH command with a specified timeout
func (se *Sh) Execute(command string, timeout time.Duration) (string, error) {
	if runtime.GOOS != "linux" {
		se.logger.Log(logger.ERROR, "Command execution failed: Unsupported OS")
		return "", fmt.Errorf("%s is only supported on Linux systems", se.shortName)
	}

	se.logger.Log(logger.INFO, "Executing command: %s by %s (Timeout: %v)", command, se.shortName, timeout)

	// Create a context with the specified timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Construct the execution command
	cmd := exec.CommandContext(ctx, se.path, "-c", command)

	// Capture the output
	output, err := cmd.CombinedOutput()

	// Check for context timeout
	if ctx.Err() == context.DeadlineExceeded {
		se.logger.Log(logger.WARN, "Command timed out after %v", timeout)
		return "", fmt.Errorf("command timed out after %v", timeout)
	}

	if err != nil {
		fmt.Println(colorprint.ColorString(fmt.Sprintf("Command execution failed: %v\nOutput: %s", err, string(output)), colorprint.RED))
		return "", fmt.Errorf("failed to execute %s command: %v\nOutput: %s", se.shortName, err, string(output))
	}

	se.logger.Log(logger.DEBUG, "Command executed successfully. Output:\n%s", string(output))
	return string(output), nil
}

func (se *Sh) ShortName() string {
	return se.shortName
}

func (se *Sh) Path() string {
	return se.path
}
