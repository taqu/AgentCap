package reduce

import (
	"fmt"
	"strings"

	"github.com/taqu/agentcap/internal/buildparse"
	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/exec"
)

type GccReducer struct {
	Args        []string
	Tool        string
	ParsedBuild *buildparse.BuildResult
}

func (r *GccReducer) Reduce(result *exec.Result) *ReducedResult {
	raw := append(result.Stdout, result.Stderr...)
	rawBytes := len(result.Stdout) + len(result.Stderr)
	cleaned := clean.StripANSI(raw)

	var parsed *buildparse.BuildResult
	if r.Tool == "clang" {
		parsed = buildparse.ParseClang(cleaned)
	} else {
		parsed = buildparse.ParseGCC(cleaned)
	}
	if parsed == nil {
		return (&GenericReducer{}).Reduce(result)
	}
	parsed.Pass = result.ExitCode == 0
	r.ParsedBuild = parsed

	var sb strings.Builder
	tool := r.Tool
	if tool == "" {
		tool = "gcc"
	}

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
		fmt.Fprintf(&sb, "@acap %s-build PASS\n", tool)
		if warnings > 0 {
			fmt.Fprintf(&sb, "warnings=%d\n\n", warnings)
			shown := 0
			for _, d := range parsed.Diagnostics {
				if d.Severity != buildparse.SeverityWarning {
					continue
				}
				if shown >= 10 {
					fmt.Fprintf(&sb, "...omitted=%d\n", warnings-shown)
					break
				}
				fmt.Fprintf(&sb, "W %s:%d %s\n", d.File, d.Line, d.Message)
				shown++
			}
		}
	} else {
		fmt.Fprintf(&sb, "@acap %s-build FAIL\nerrors=%d", tool, errors)
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
			if d.File == "(linker)" {
				fmt.Fprintf(&sb, "E (linker)\n  %s\n", d.Message)
			} else {
				fmt.Fprintf(&sb, "E %s:%d:%d\n  %s\n", d.File, d.Line, d.Col, d.Message)
			}
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
