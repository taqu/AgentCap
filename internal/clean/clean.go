// Package clean provides utilities for cleaning terminal output for LLM consumption.
package clean

import (
	"bytes"
	"strings"
)

// StripANSI removes ANSI escape sequences (color, cursor movement, etc.) from b.
func StripANSI(b []byte) []byte {
	if !bytes.ContainsRune(b, '\x1b') {
		return b
	}
	out := make([]byte, 0, len(b))
	i := 0
	for i < len(b) {
		if b[i] == '\x1b' && i+1 < len(b) && b[i+1] == '[' {
			// CSI sequence: ESC [ ... (letter)
			i += 2
			for i < len(b) && !isCSITerminator(b[i]) {
				i++
			}
			if i < len(b) {
				i++ // consume terminator
			}
			continue
		}
		if b[i] == '\x1b' && i+1 < len(b) {
			// Other escape: ESC + single char (OSC, etc.) — skip both
			i += 2
			continue
		}
		out = append(out, b[i])
		i++
	}
	return out
}

func isCSITerminator(c byte) bool {
	return c >= 0x40 && c <= 0x7E
}

// StripProgress removes carriage-return overwritten progress lines.
// Terminal progress bars often write partial lines followed by \r to
// overwrite them; only the final content on each line survives.
func StripProgress(b []byte) []byte {
	if !bytes.ContainsRune(b, '\r') {
		return b
	}
	lines := bytes.Split(b, []byte{'\n'})
	out := make([]byte, 0, len(b))
	for i, line := range lines {
		// Within each \n-separated segment, keep only the last \r-segment.
		parts := bytes.Split(line, []byte{'\r'})
		kept := parts[len(parts)-1]
		out = append(out, kept...)
		if i < len(lines)-1 {
			out = append(out, '\n')
		}
	}
	return out
}

// Lines splits b into lines, trimming a single trailing empty line.
func Lines(b []byte) []string {
	if len(b) == 0 {
		return nil
	}
	s := string(b)
	lines := strings.Split(s, "\n")
	// Trim trailing empty line that results from a final newline.
	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	return lines
}

// CollapseRepeated collapses runs of more than 5 consecutive identical lines
// into one representative line plus a count annotation.
func CollapseRepeated(lines []string) []string {
	const threshold = 5
	if len(lines) == 0 {
		return lines
	}
	out := make([]string, 0, len(lines))
	i := 0
	for i < len(lines) {
		j := i + 1
		for j < len(lines) && lines[j] == lines[i] {
			j++
		}
		count := j - i
		if count > threshold {
			out = append(out, lines[i])
			out = append(out, strings.Repeat("", 0)+"("+itoa(count-1)+" repeated)")
		} else {
			out = append(out, lines[i:j]...)
		}
		i = j
	}
	return out
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	buf := [20]byte{}
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		pos--
		buf[pos] = '-'
	}
	return string(buf[pos:])
}
