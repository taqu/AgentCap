package buildparse

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

var goBuildLineRe = regexp.MustCompile(`^(.+\.go):(\d+):(\d+): (.+)$`)

func ParseGoBuild(raw []byte) *BuildResult {
	if len(raw) > MaxBuildBytes {
		return nil
	}
	result := &BuildResult{Tool: "go-build"}
	lines := bytes.Split(bytes.TrimRight(raw, "\n"), []byte("\n"))

	for _, lb := range lines {
		line := string(lb)
		if strings.HasPrefix(line, "#") {
			continue
		}
		m := goBuildLineRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		lineN, _ := strconv.Atoi(m[2])
		colN, _ := strconv.Atoi(m[3])
		msg := m[4]
		sev := SeverityError
		if strings.HasPrefix(msg, "warning: ") {
			sev = SeverityWarning
			msg = strings.TrimPrefix(msg, "warning: ")
		}
		result.Diagnostics = append(result.Diagnostics, Diagnostic{
			Tool: "go-build", Severity: sev,
			File: m[1], Line: lineN, Col: colN, Message: msg,
		})
		if len(result.Diagnostics) >= MaxDiagnostics {
			result.Truncated = true
			break
		}
	}
	return result
}
