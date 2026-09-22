package reduce

import (
	"fmt"
	"strings"

	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/exec"
)

type MakeReducer struct{ Args []string }
type NinjaReducer struct{ Args []string }

func (r *MakeReducer) Reduce(result *exec.Result) *ReducedResult {
	raw := append(result.Stdout, result.Stderr...)
	rawBytes := len(result.Stdout) + len(result.Stderr)
	cleaned := string(clean.StripANSI(raw))

	if len(cleaned) <= smallThresholdBytes {
		return &ReducedResult{Output: cleaned, RawBytes: rawBytes, RetBytes: len(cleaned)}
	}

	lines := strings.Split(cleaned, "\n")
	var kept []string
	for _, l := range lines {
		// Keep error/failure lines and diagnostics
		if isUsefulBuildLine(l) {
			kept = append(kept, l)
		}
	}

	if result.ExitCode == 0 {
		out := "@acap make PASS\n"
		return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap make FAIL\n\n")
	shown := 0
	for _, l := range kept {
		if shown >= 50 {
			break
		}
		sb.WriteString(l + "\n")
		shown++
	}
	out := sb.String()
	return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
}

func (r *NinjaReducer) Reduce(result *exec.Result) *ReducedResult {
	raw := append(result.Stdout, result.Stderr...)
	rawBytes := len(result.Stdout) + len(result.Stderr)
	cleaned := string(clean.StripANSI(raw))

	lines := strings.Split(cleaned, "\n")
	var kept []string
	steps := 0
	for _, l := range lines {
		// Skip [N/M] progress lines
		if len(l) > 0 && l[0] == '[' {
			steps++
			continue
		}
		if isUsefulBuildLine(l) || l == "FAILED" {
			kept = append(kept, l)
		}
	}

	if result.ExitCode == 0 {
		out := fmt.Sprintf("@acap ninja PASS steps=%d\n", steps)
		return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap ninja FAIL\n\n")
	for i, l := range kept {
		if i >= 50 {
			break
		}
		sb.WriteString(l + "\n")
	}
	out := sb.String()
	return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
}

func isUsefulBuildLine(l string) bool {
	if strings.Contains(l, ": error:") || strings.Contains(l, ": fatal error:") {
		return true
	}
	if strings.Contains(l, "Error") || strings.Contains(l, "FAILED") || strings.Contains(l, "FAIL") {
		return true
	}
	if strings.Contains(l, "undefined reference") || strings.Contains(l, "ld: error") {
		return true
	}
	if strings.Contains(l, "make[") && strings.Contains(l, "Error") {
		return true
	}
	return false
}
