package clean

import (
	"testing"
)

func TestStripANSI(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"no escapes", "hello world", "hello world"},
		{"color reset", "hello\x1b[0mworld", "helloworld"},
		{"bold red", "\x1b[1;31merror\x1b[0m: file not found", "error: file not found"},
		{"cursor movement", "\x1b[2Jhello", "hello"},
		{"multiple sequences", "\x1b[32mok\x1b[0m \x1b[31mfail\x1b[0m", "ok fail"},
		{"unicode preserved", "\x1b[32m日本語\x1b[0m", "日本語"},
		{"other escape", "\x1bXsomething", "something"}, // non-CSI
		{"no trailing newline loss", "line1\nline2\n", "line1\nline2\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(StripANSI([]byte(tt.input)))
			if got != tt.want {
				t.Errorf("StripANSI(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestStripProgress(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{"empty", "", ""},
		{"no CR", "line1\nline2\n", "line1\nline2\n"},
		{"simple progress", "10%\r20%\r30%\n", "30%\n"},
		{"multi-line with progress", "start\n10%\r20%\rdone\nend\n", "start\ndone\nend\n"},
		{"only CR", "a\rb\rc\n", "c\n"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := string(StripProgress([]byte(tt.input)))
			if got != tt.want {
				t.Errorf("StripProgress(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestLines(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{"empty", "", nil},
		{"single line no newline", "hello", []string{"hello"}},
		{"single line with newline", "hello\n", []string{"hello"}},
		{"two lines", "a\nb\n", []string{"a", "b"}},
		{"blank lines preserved", "a\n\nb\n", []string{"a", "", "b"}},
		{"unicode", "こんにちは\n世界\n", []string{"こんにちは", "世界"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Lines([]byte(tt.input))
			if len(got) != len(tt.want) {
				t.Errorf("Lines(%q): got %d lines, want %d; got=%v want=%v", tt.input, len(got), len(tt.want), got, tt.want)
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Lines(%q)[%d] = %q, want %q", tt.input, i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestCollapseRepeated(t *testing.T) {
	tests := []struct {
		name      string
		input     []string
		wantLines int
		checkFunc func([]string) bool
	}{
		{
			name:      "empty",
			input:     nil,
			wantLines: 0,
		},
		{
			name:      "no repeats",
			input:     []string{"a", "b", "c"},
			wantLines: 3,
		},
		{
			name:      "exactly 5 repeats — not collapsed",
			input:     []string{"x", "x", "x", "x", "x"},
			wantLines: 5,
		},
		{
			name:      "6 repeats — collapsed",
			input:     []string{"x", "x", "x", "x", "x", "x"},
			wantLines: 2, // "x" + "(5 repeated)"
		},
		{
			name:      "10 repeats",
			input:     []string{"y", "y", "y", "y", "y", "y", "y", "y", "y", "y"},
			wantLines: 2,
		},
		{
			name:      "mixed",
			input:     []string{"a", "b", "b", "b", "b", "b", "b", "c"},
			wantLines: 4, // a, b, "(5 repeated)", c
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CollapseRepeated(tt.input)
			if len(got) != tt.wantLines {
				t.Errorf("CollapseRepeated: got %d lines %v, want %d", len(got), got, tt.wantLines)
			}
		})
	}
}
