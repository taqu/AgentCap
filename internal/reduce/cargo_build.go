package reduce

import (
	"fmt"
	"strings"

	"github.com/taqu/agentcap/internal/buildparse"
	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/exec"
)

type CargoBuildReducer struct {
	Args        []string
	ParsedBuild *buildparse.BuildResult
}

func (r *CargoBuildReducer) Reduce(result *exec.Result) *ReducedResult {
	raw := append(result.Stdout, result.Stderr...)
	rawBytes := len(result.Stdout) + len(result.Stderr)
	cleaned := clean.StripANSI(raw)

	parsed := buildparse.ParseCargoBuild(cleaned)
	if parsed == nil {
		return (&GenericReducer{}).Reduce(result)
	}
	parsed.Pass = result.ExitCode == 0
	r.ParsedBuild = parsed

	var sb strings.Builder
	errors, warnings := 0, 0
	for _, d := range parsed.Diagnostics {
		switch d.Severity {
		case buildparse.SeverityError, buildparse.SeverityFatal:
			errors++
		case buildparse.SeverityWarning:
			warnings++
		}
	}

	if parsed.Pass {
		sb.WriteString("@acap cargo-build PASS\n")
		if warnings > 0 {
			fmt.Fprintf(&sb, "warnings=%d\n\n", warnings)
			shown := 0
			for _, d := range parsed.Diagnostics {
				if d.Severity != buildparse.SeverityWarning {
					continue
				}
				if shown >= 10 {
					break
				}
				fmt.Fprintf(&sb, "W %s %s:%d %s\n", d.Code, d.File, d.Line, d.Message)
				shown++
			}
		}
	} else {
		fmt.Fprintf(&sb, "@acap cargo-build FAIL\nerrors=%d", errors)
		if warnings > 0 {
			fmt.Fprintf(&sb, " warnings=%d", warnings)
		}
		sb.WriteString("\n\n")
		shown := 0
		for _, d := range parsed.Diagnostics {
			if d.Severity == buildparse.SeverityWarning || d.Severity == buildparse.SeverityNote {
				continue
			}
			if shown >= 50 {
				fmt.Fprintf(&sb, "...omitted=%d\n", errors-shown)
				break
			}
			code := ""
			if d.Code != "" {
				code = "[" + d.Code + "] "
			}
			fmt.Fprintf(&sb, "E%s %s:%d:%d\n  %s\n", code, d.File, d.Line, d.Col, d.Message)
			for _, n := range d.Notes {
				if n.Severity == buildparse.SeverityHelp {
					fmt.Fprintf(&sb, "  help: %s\n", n.Message)
				}
			}
			shown++
		}
	}

	out := sb.String()
	return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
}
