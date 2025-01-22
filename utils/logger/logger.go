package logger

import (
	"calderat/utils/colorprint"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

type LogLevel int

const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
)

type Logger struct {
	level  LogLevel
	logger *log.Logger
}

// New creates a new Logger instance with the specified log level and default output (stdout).
func New(logLevel string) (*Logger, error) {
	level, err := parseLogLevel(logLevel)
	if err != nil {
		return nil, err
	}
	return &Logger{
		level:  level,
		logger: log.New(os.Stdout, "", log.LstdFlags),
	}, nil
}

// NewWithOutput creates a new Logger instance with the specified log level and custom output.
func NewWithOutput(logLevel string, output *os.File) (*Logger, error) {
	level, err := parseLogLevel(logLevel)
	if err != nil {
		return nil, err
	}
	return &Logger{
		level:  level,
		logger: log.New(output, "", log.LstdFlags),
	}, nil
}

// parseLogLevel converts a log level string to the LogLevel type.
func parseLogLevel(level string) (LogLevel, error) {
	switch strings.ToUpper(level) {
	case "TRACE":
		return TRACE, nil
	case "DEBUG":
		return DEBUG, nil
	case "INFO":
		return INFO, nil
	case "WARN":
		return WARN, nil
	case "ERROR":
		return ERROR, nil
	default:
		return INFO, fmt.Errorf("invalid log level: %s", level)
	}
}

// Log logs a message if the message level is greater than or equal to the logger's level.
func (l *Logger) Log(level LogLevel, format string, v ...interface{}) {
	if level >= l.level {
		prefix := ""
		switch level {
		case TRACE:
			prefix = colorprint.ColorString("[TRACE] "+format, colorprint.CYAN)
		case DEBUG:
			prefix = colorprint.ColorString("[DEBUG] "+format, colorprint.GREEN)
		case INFO:
			prefix = colorprint.ColorString("[INFO] "+format, colorprint.CYAN)
		case WARN:
			prefix = colorprint.ColorString("[WARN] "+format, colorprint.YELLOW)
		case ERROR:
			prefix = colorprint.ColorString("[ERROR] "+format, colorprint.RED)
		}
		l.logger.Printf(prefix, v...)

		// Add a stack trace for WARN and ERROR levels
		if level >= WARN {
			stack := getStackTrace()
			l.logger.Printf(colorprint.ColorString("[STACK TRACE]\n%s", colorprint.MAGENTA), stack)
		}
	}
}

func getStackTrace() string {
	var sb strings.Builder
	pc := make([]uintptr, 10) // Limit stack trace depth to 10 frames
	n := runtime.Callers(3, pc)

	frames := runtime.CallersFrames(pc[:n])
	for {
		frame, more := frames.Next()
		// Format: function_name (file:line)
		sb.WriteString(fmt.Sprintf("%s (%s:%d)\n", frame.Function, frame.File, frame.Line))
		if !more {
			break
		}
	}

	return sb.String()
}

// SetLevel dynamically updates the log level.
func (l *Logger) SetLevel(logLevel string) error {
	level, err := parseLogLevel(logLevel)
	if err != nil {
		return err
	}
	l.level = level
	return nil
}
