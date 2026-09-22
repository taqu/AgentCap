package reduce

import (
	"fmt"
	"strings"

	"github.com/taqu/agentcap/internal/buildparse"
	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/exec"
)

type CargoTestReducer struct {
	Args      []string
	ParsedRun *buildparse.TestRun
}

func (r *CargoTestReducer) Reduce(result *exec.Result) *ReducedResult {
	raw := append(result.Stdout, result.Stderr...)
	rawBytes := len(result.Stdout) + len(result.Stderr)
	cleaned := clean.StripANSI(raw)

	parsed := buildparse.ParseCargoTest(cleaned)
	if parsed == nil {
		return (&GenericReducer{}).Reduce(result)
	}
	parsed.Pass = result.ExitCode == 0
	r.ParsedRun = parsed

	var sb strings.Builder
	totalPass, totalFail, totalIgnore := 0, 0, 0
	for _, p := range parsed.Packages {
		totalPass += p.Passed
		totalFail += p.Failed
		totalIgnore += p.Skipped
	}

	if parsed.Pass {
		fmt.Fprintf(&sb, "@acap cargo-test PASS\npass=%d ignored=%d\n", totalPass, totalIgnore)
	} else {
		fmt.Fprintf(&sb, "@acap cargo-test FAIL\npass=%d fail=%d ignored=%d\n\nfailures:\n",
			totalPass, totalFail, totalIgnore)
		for _, f := range parsed.Failures {
			name := f.Name
			if f.Suite != "" {
				name = f.Suite + "::" + f.Name
			}
			fmt.Fprintf(&sb, "%s\n", name)
			for i, l := range f.Output {
				if i >= 3 {
					break
				}
				fmt.Fprintf(&sb, "  %s\n", l)
			}
		}
	}

	out := sb.String()
	return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
}
