package reduce

import (
	"fmt"
	"strings"

	"github.com/taqu/agentcap/internal/buildparse"
	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/exec"
)

type GoBuildReducer struct {
	Args        []string
	ParsedBuild *buildparse.BuildResult
}

func (r *GoBuildReducer) Reduce(result *exec.Result) *ReducedResult {
	// go build writes errors to stderr
	raw := append(result.Stdout, result.Stderr...)
	rawBytes := len(result.Stdout) + len(result.Stderr)
	cleaned := clean.StripANSI(raw)

	parsed := buildparse.ParseGoBuild(cleaned)
	if parsed == nil {
		return (&GenericReducer{}).Reduce(result)
	}
	parsed.Pass = result.ExitCode == 0
	r.ParsedBuild = parsed

	var sb strings.Builder
	if parsed.Pass {
		sb.WriteString("@acap go-build PASS\n")
		if parsed.Duration != "" {
			fmt.Fprintf(&sb, "duration=%s\n", parsed.Duration)
		}
		var warnings []buildparse.Diagnostic
		for _, d := range parsed.Diagnostics {
			if d.Severity == buildparse.SeverityWarning {
				warnings = append(warnings, d)
			}
		}
		if len(warnings) > 0 {
			fmt.Fprintf(&sb, "warnings=%d\n\n", len(warnings))
			for i, w := range warnings {
				if i >= 10 {
					fmt.Fprintf(&sb, "...omitted=%d\n", len(warnings)-10)
					break
				}
				fmt.Fprintf(&sb, "W %s:%d %s\n", w.File, w.Line, w.Message)
			}
		}
	} else {
		errors := 0
		for _, d := range parsed.Diagnostics {
			if d.Severity == buildparse.SeverityError || d.Severity == buildparse.SeverityFatal {
				errors++
			}
		}
		warnings := 0
		for _, d := range parsed.Diagnostics {
			if d.Severity == buildparse.SeverityWarning {
				warnings++
			}
		}
		fmt.Fprintf(&sb, "@acap go-build FAIL\nerrors=%d\n", errors)
		if warnings > 0 {
			fmt.Fprintf(&sb, "warnings=%d\n", warnings)
		}
		sb.WriteByte('\n')
		shown := 0
		for _, d := range parsed.Diagnostics {
			if d.Severity == buildparse.SeverityWarning || d.Severity == buildparse.SeverityNote {
				continue
			}
			if shown >= 50 {
				fmt.Fprintf(&sb, "...omitted=%d\n", errors-shown)
				break
			}
			fmt.Fprintf(&sb, "E %s:%d:%d\n  %s\n", d.File, d.Line, d.Col, d.Message)
			shown++
		}
	}

	out := sb.String()
	return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
}
