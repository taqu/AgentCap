package util

import "strings"

// PadRight pads s with spaces to width w.
func PadRight(s string, w int) string {
	if len(s) >= w {
		return s
	}
	return s + strings.Repeat(" ", w-len(s))
}

// Truncate shortens s to at most n bytes, adding "..." when cut.
func Truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	if n <= 3 {
		return s[:n]
	}
	return s[:n-3] + "..."
}
