package delta

import (
	"context"
	"fmt"
	"strings"

	"github.com/taqu/agentcap/internal/buildparse"
	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/store"
)

// DiagnosticsDelta compares two build/compiler results.
type DiagnosticsDelta struct {
	Store *store.Store
}

func (d *DiagnosticsDelta) Delta(ctx context.Context, baseline, current *store.Entry) (*Result, error) {
	prevDiags, err := d.Store.GetDiagnostics(baseline.Meta.ID)
	if err != nil || prevDiags == nil {
		return nil, nil
	}

	currRaw := readFile(current.StdoutPath())
	currErr := readFile(current.StderrPath())
	allCurrRaw := append(currRaw, currErr...)

	// Re-parse current raw output based on reducer type
	var currDiags []buildparse.Diagnostic
	reducerName := current.Meta.Reducer
	switch {
	case reducerName == "go-build":
		if p := buildparse.ParseGoBuild(clean.StripANSI(allCurrRaw)); p != nil {
			currDiags = p.Diagnostics
		}
	case reducerName == "gcc" || reducerName == "clang":
		if p := buildparse.ParseGCC(clean.StripANSI(allCurrRaw)); p != nil {
			currDiags = p.Diagnostics
		}
	case reducerName == "cargo-build" || reducerName == "cargo-check":
		if p := buildparse.ParseCargoBuild(clean.StripANSI(allCurrRaw)); p != nil {
			currDiags = p.Diagnostics
		}
	default:
		return nil, nil
	}

	// Build identity sets
	prevKeys := make(map[string]buildparse.Diagnostic, len(prevDiags))
	for _, d := range prevDiags {
		prevKeys[d.Key()] = d
	}
	currKeys := make(map[string]buildparse.Diagnostic, len(currDiags))
	for _, d := range currDiags {
		currKeys[d.Key()] = d
	}

	var newDiags, resolvedDiags []buildparse.Diagnostic
	for k, d := range currKeys {
		if _, ok := prevKeys[k]; !ok {
			newDiags = append(newDiags, d)
		}
	}
	for k, d := range prevKeys {
		if _, ok := currKeys[k]; !ok {
			resolvedDiags = append(resolvedDiags, d)
		}
	}

	// Count transitions
	prevErrors := countErrors(prevDiags)
	currErrors := countErrors(currDiags)

	if len(newDiags) == 0 && len(resolvedDiags) == 0 {
		return nil, nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap delta from %s %s\n", baseline.Meta.ID, reducerName)

	if prevErrors != currErrors {
		fmt.Fprintf(&sb, "errors %d -> %d\n", prevErrors, currErrors)
	}

	if len(newDiags) > 0 {
		fmt.Fprintf(&sb, "\nnew_errors=%d\n", len(newDiags))
		for _, d := range newDiags {
			fmt.Fprintf(&sb, "E %s:%d %s\n", d.File, d.Line, d.Message)
		}
	}
	if len(resolvedDiags) > 0 {
		fmt.Fprintf(&sb, "\nresolved=%d\n", len(resolvedDiags))
		for _, d := range resolvedDiags {
			fmt.Fprintf(&sb, "  %s:%d %s\n", d.File, d.Line, d.Message)
		}
	}
	unchanged := len(currDiags) - len(newDiags)
	if unchanged > 0 {
		fmt.Fprintf(&sb, "\nunchanged=%d\n", unchanged)
	}

	return &Result{Output: sb.String(), Presentation: PresentationDelta}, nil
}

func countErrors(diags []buildparse.Diagnostic) int {
	n := 0
	for _, d := range diags {
		if d.Severity == buildparse.SeverityError || d.Severity == buildparse.SeverityFatal {
			n++
		}
	}
	return n
}

// TestResultsDelta compares two test runs by failure identity.
type TestResultsDelta struct {
	Store *store.Store
}

func (d *TestResultsDelta) Delta(ctx context.Context, baseline, current *store.Entry) (*Result, error) {
	prevFailures, err := d.Store.GetTestFailures(baseline.Meta.ID)
	if err != nil {
		return nil, nil
	}

	currRaw := readFile(current.StdoutPath())
	currErr := readFile(current.StderrPath())
	allCurrRaw := append(currRaw, currErr...)

	var currRun *buildparse.TestRun
	reducerName := current.Meta.Reducer
	switch reducerName {
	case "go-test":
		currRun = buildparse.ParseGoTest(clean.StripANSI(allCurrRaw))
	case "cargo-test":
		currRun = buildparse.ParseCargoTest(clean.StripANSI(allCurrRaw))
	default:
		return nil, nil
	}
	if currRun == nil {
		return nil, nil
	}

	// Build identity sets
	prevSet := make(map[string]bool, len(prevFailures))
	for _, f := range prevFailures {
		prevSet[f.FullName()] = true
	}
	currSet := make(map[string]buildparse.TestFailure, len(currRun.Failures))
	for _, f := range currRun.Failures {
		currSet[f.FullName()] = f
	}

	var newFails []buildparse.TestFailure
	var resolvedNames []string
	for name, f := range currSet {
		if !prevSet[name] {
			newFails = append(newFails, f)
		}
	}
	for name := range prevSet {
		if _, ok := currSet[name]; !ok {
			resolvedNames = append(resolvedNames, name)
		}
	}

	if len(newFails) == 0 && len(resolvedNames) == 0 {
		// Check pass/fail transition
		currPass := current.Meta.ExitCode == 0
		basePass := baseline.Meta.ExitCode == 0
		if currPass && !basePass {
			out := fmt.Sprintf("@acap delta from %s\n%s FAIL -> PASS\nresolved_failures=%d\n",
				baseline.Meta.ID, reducerName, len(prevFailures))
			return &Result{Output: out, Presentation: PresentationDelta}, nil
		}
		return nil, nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap delta from %s %s\nfailed %d -> %d\n",
		baseline.Meta.ID, reducerName, len(prevFailures), len(currRun.Failures))

	if len(newFails) > 0 {
		fmt.Fprintf(&sb, "\nnew:\n")
		for _, f := range newFails {
			fmt.Fprintf(&sb, "  %s\n", f.FullName())
			if f.File != "" && f.Line > 0 {
				fmt.Fprintf(&sb, "    %s:%d\n", f.File, f.Line)
			}
		}
	}
	if len(resolvedNames) > 0 {
		fmt.Fprintf(&sb, "\nresolved:\n")
		for _, name := range resolvedNames {
			fmt.Fprintf(&sb, "  %s\n", name)
		}
	}

	return &Result{Output: sb.String(), Presentation: PresentationDelta}, nil
}
