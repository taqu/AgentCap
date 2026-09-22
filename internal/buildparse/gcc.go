package buildparse

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

var gccDiagRe = regexp.MustCompile(`^(.+):(\d+):(\d+): (fatal error|error|warning|note|help): (.+)$`)
var linkerErrorRe = regexp.MustCompile(`^(.+): undefined reference to`)

func ParseGCC(raw []byte) *BuildResult {
	return parseGCCLike(raw, "gcc")
}

func ParseClang(raw []byte) *BuildResult {
	return parseGCCLike(raw, "clang")
}

func parseGCCLike(raw []byte, tool string) *BuildResult {
	if len(raw) > MaxBuildBytes {
		return nil
	}
	result := &BuildResult{Tool: tool}

	lines := bytes.Split(bytes.TrimRight(raw, "\n"), []byte("\n"))
	var current *Diagnostic

	for _, lb := range lines {
		line := string(lb)

		if m := gccDiagRe.FindStringSubmatch(line); m != nil {
			lineN, _ := strconv.Atoi(m[2])
			colN, _ := strconv.Atoi(m[3])
			sev := parseSeverity(m[4])

			if sev == SeverityNote || sev == SeverityHelp {
				if current != nil {
					current.Notes = append(current.Notes, DiagnosticNote{
						Severity: sev, File: m[1], Line: lineN, Col: colN, Message: m[5],
					})
					continue
				}
			}

			if current != nil {
				result.Diagnostics = append(result.Diagnostics, *current)
			}
			current = &Diagnostic{
				Tool: tool, Severity: sev,
				File: m[1], Line: lineN, Col: colN, Message: m[5],
			}
			if len(result.Diagnostics) >= MaxDiagnostics {
				result.Truncated = true
				current = nil
				break
			}
			continue
		}

		// Linker errors: "undefined reference to 'foo'"
		if strings.Contains(line, "undefined reference to") || strings.Contains(line, "ld: error") {
			if current != nil {
				result.Diagnostics = append(result.Diagnostics, *current)
				current = nil
			}
			result.Diagnostics = append(result.Diagnostics, Diagnostic{
				Tool: tool, Severity: SeverityError,
				File: "(linker)", Message: strings.TrimSpace(line),
			})
		}
		// Skip caret lines (^^^), source snippet lines
	}
	if current != nil {
		result.Diagnostics = append(result.Diagnostics, *current)
	}
	return result
}

func parseSeverity(s string) Severity {
	switch s {
	case "fatal error":
		return SeverityFatal
	case "error":
		return SeverityError
	case "warning":
		return SeverityWarning
	case "note":
		return SeverityNote
	case "help":
		return SeverityHelp
	default:
		return SeverityError
	}
}
