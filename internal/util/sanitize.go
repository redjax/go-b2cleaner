package util

import (
	"regexp"
	"strings"
)

// SanitizeFileName replaces unsafe characters and ensures a single extension at the end.
func SanitizeFileName(input, extension string) string {
	// Remove extension if already present
	base := strings.TrimSuffix(input, extension)

	// Replace unsafe characters with "_"
	unsafe := regexp.MustCompile(`[:/\\]`)
	safe := unsafe.ReplaceAllString(base, "_")

	// Optionally, replace spaces with underscores or hyphens
	safe = strings.ReplaceAll(safe, " ", "_")

	// Ensure only one extension at the end
	return safe + extension
}
