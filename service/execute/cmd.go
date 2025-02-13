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

type Cmd struct {
	shortName string
	logger    *logger.Logger
	path      string
}

// NewCmd initializes a new CMD executor
func NewCmd(log *logger.Logger) *Cmd {
	return &Cmd{
		shortName: "cmd",
		logger:    log,
		path:      "cmd",
	}
}

// Execute runs a CMD command with a specified timeout
func (ce *Cmd) Execute(command string, timeout time.Duration) (string, error) {
	if runtime.GOOS != "windows" {
		ce.logger.Log(logger.ERROR, "Command execution failed: Unsupported OS")
		return "", fmt.Errorf("%s is only supported on Windows systems", ce.shortName)
	}

	ce.logger.Log(logger.INFO, "Executing command: %s by %s (Timeout: %v)", command, ce.shortName, timeout)

	// Create a context with the specified timeout
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Construct the execution command
	cmd := exec.CommandContext(ctx, ce.path, "/C", command)

	// Capture the output
	output, err := cmd.CombinedOutput()

	// Check for context timeout
	if ctx.Err() == context.DeadlineExceeded {
		ce.logger.Log(logger.WARN, "Command timed out after %v", timeout)
		return "", fmt.Errorf("command timed out after %v", timeout)
	}

	if err != nil {
		fmt.Println(colorprint.ColorString(fmt.Sprintf("Command execution failed: %v\nOutput: %s", err, string(output)), colorprint.RED))
		return "", fmt.Errorf("failed to execute %s command: %v\nOutput: %s", ce.shortName, err, string(output))
	}

	ce.logger.Log(logger.DEBUG, "Command executed successfully. Output:\n%s", string(output))
	return string(output), nil
}

func (ce *Cmd) ShortName() string {
	return ce.shortName
}

func (ce *Cmd) Path() string {
	return ce.path
}
