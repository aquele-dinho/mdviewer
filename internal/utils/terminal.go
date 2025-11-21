package utils

import (
	"os"

	"golang.org/x/term"
)

// GetTerminalWidth returns the current terminal width
func GetTerminalWidth() int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil || width == 0 {
		// Default to 80 columns if we can't detect
		return 80
	}
	return width
}

// IsTerminal checks if stdout is a terminal
func IsTerminal() bool {
	return term.IsTerminal(int(os.Stdout.Fd()))
}

// IsTTY checks if the given file descriptor is a TTY
func IsTTY(fd int) bool {
	return term.IsTerminal(fd)
}
