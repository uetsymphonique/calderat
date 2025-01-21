package logger

import (
	"fmt"
	"log"
	"os"
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
			prefix = "[TRACE] "
		case DEBUG:
			prefix = "[DEBUG] "
		case INFO:
			prefix = "[INFO] "
		case WARN:
			prefix = "[WARN] "
		case ERROR:
			prefix = "[ERROR] "
		}
		l.logger.Printf(prefix+format, v...)
	}
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
