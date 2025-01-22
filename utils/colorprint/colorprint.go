package colorprint

import (
	"os"
	"strings"
)

const (
	BLACK   = "\033[30m"
	RED     = "\033[31m"
	GREEN   = "\033[32m"
	YELLOW  = "\033[33m"
	BLUE    = "\033[34m"
	MAGENTA = "\033[35m"
	CYAN    = "\033[36m"
	WHITE   = "\033[37m"
	RESET   = "\033[0m"
)

func supportsColor() bool {
	term := os.Getenv("TERM")
	return strings.Contains(term, "xterm") || strings.Contains(term, "screen") || strings.Contains(term, "color")
}

func ColorString(s string, color string) string {
	return color + s + RESET
}
