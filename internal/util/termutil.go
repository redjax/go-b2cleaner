package util

import (
	"os"

	"github.com/mattn/go-isatty"
	"golang.org/x/term"
)

// DetectTerminalWidth tries to get the terminal width, falling back to a default if necessary.
func DetectTerminalWidth(fallback int) int {
	fd := os.Stdout.Fd()
	if isatty.IsTerminal(fd) {
		w, _, err := term.GetSize(int(fd))
		if err == nil && w >= 80 {
			return w
		}
	}
	return fallback
}

// MaxNameLen calculates the max length for the "Name" column given terminal width and other column widths.
func MaxNameLen(termWidth, typeCol, sizeCol, createdCol, borders int) int {
	maxNameLen := termWidth - (typeCol + sizeCol + createdCol + borders)
	if maxNameLen < 10 {
		maxNameLen = 10
	}
	return maxNameLen
}
