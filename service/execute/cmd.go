package execute

import (
	"calderat/utils/colorprint"
	"calderat/utils/logger"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type Cmd struct {
	shortName string
	logger    *logger.Logger
	path      string
}

func parseWindowsCmd(command string) ([]string, error) {
	var args []string
	var current strings.Builder
	inQuotes := false

	for i := 0; i < len(command); i++ {
		c := command[i]
		switch c {
		case '"':
			// Toggle whether we're inside quotes.
			inQuotes = !inQuotes
		case ' ':
			if inQuotes {
				current.WriteByte(c)
			} else {
				if current.Len() > 0 {
					args = append(args, current.String())
					current.Reset()
				}
			}
		default:
			current.WriteByte(c)
		}
	}
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	return args, nil
}

// NewCmd initializes a new CMD executor
func NewCmd(log *logger.Logger) *Cmd {
	return &Cmd{
		shortName: "cmd",
		logger:    log,
		path:      "cmd.exe",
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
	// Use shlex to split the command string into arguments.
	parsedArgs, err := parseWindowsCmd(command)
	if err != nil {
		ce.logger.Log(logger.ERROR, "Error parsing command: %v", err)
		return "", fmt.Errorf("error parsing command: %v", err)
	}

	// Prepend "/C" so that cmd.exe executes the command and then exits.
	args := append([]string{"/C"}, parsedArgs...)

	// Construct the execution command with the parsed arguments.
	cmd := exec.CommandContext(ctx, ce.path, args...)

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
