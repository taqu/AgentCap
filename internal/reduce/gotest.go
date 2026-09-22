package reduce

import (
	"fmt"
	"strings"

	"github.com/taqu/agentcap/internal/buildparse"
	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/exec"
)

type GoTestReducer struct {
	Args      []string
	ParsedRun *buildparse.TestRun
}

func (r *GoTestReducer) Reduce(result *exec.Result) *ReducedResult {
	raw := append(result.Stdout, result.Stderr...)
	rawBytes := len(result.Stdout) + len(result.Stderr)
	cleaned := clean.StripANSI(raw)

	parsed := buildparse.ParseGoTest(cleaned)
	if parsed == nil {
		return (&GenericReducer{}).Reduce(result)
	}
	parsed.Pass = result.ExitCode == 0
	r.ParsedRun = parsed

	var sb strings.Builder
	totalPass, totalFail := 0, 0
	for _, p := range parsed.Packages {
		if p.Pass {
			totalPass++
		} else {
			totalFail++
		}
	}

	if parsed.Pass {
		fmt.Fprintf(&sb, "@acap go-test PASS\npackages=%d", len(parsed.Packages))
		if parsed.Duration != "" {
			fmt.Fprintf(&sb, " duration=%s", parsed.Duration)
		}
		sb.WriteByte('\n')
	} else {
		fmt.Fprintf(&sb, "@acap go-test FAIL\npackages=%d pass=%d fail=%d\ntests_failed=%d\n",
			len(parsed.Packages), totalPass, totalFail, len(parsed.Failures))

		// Group by package
		pkgOrder := []string{}
		pkgMap := make(map[string][]buildparse.TestFailure)
		for _, f := range parsed.Failures {
			if _, ok := pkgMap[f.Package]; !ok {
				pkgOrder = append(pkgOrder, f.Package)
			}
			pkgMap[f.Package] = append(pkgMap[f.Package], f)
		}
		for _, pkg := range pkgOrder {
			fmt.Fprintf(&sb, "\n%s\n", pkg)
			for _, f := range pkgMap[pkg] {
				fmt.Fprintf(&sb, "  %s\n", f.Name)
				if f.File != "" && f.Line > 0 {
					fmt.Fprintf(&sb, "    %s:%d\n", f.File, f.Line)
				}
				for i, l := range f.Output {
					if i >= 3 {
						break
					}
					fmt.Fprintf(&sb, "    %s\n", l)
				}
			}
		}
	}

	out := sb.String()
	return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
}
