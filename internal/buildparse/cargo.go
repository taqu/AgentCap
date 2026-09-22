package buildparse

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

var (
	cargoErrorCodeRe = regexp.MustCompile(`^error\[(E\d+)\]: (.+)$`)
	cargoErrorRe     = regexp.MustCompile(`^error: (.+)$`)
	cargoWarningRe   = regexp.MustCompile(`^warning(?:\[(\w+)\])?: (.+)$`)
	cargoLocRe       = regexp.MustCompile(`^\s*--> (.+):(\d+):(\d+)$`)
	cargoHelpRe      = regexp.MustCompile(`^\s*= (help|note): (.+)$`)

	cargoTestRunningRe = regexp.MustCompile(`^running (\d+) tests?`)
	cargoTestResultRe  = regexp.MustCompile(`^test result: (ok|FAILED)\. (\d+) passed; (\d+) failed; (\d+) ignored`)
	cargoTestLineRe    = regexp.MustCompile(`^test (.+) \.\.\. (ok|FAILED|ignored)$`)
	cargoTestSectionRe = regexp.MustCompile(`^---- (.+) stdout ----$`)
)

func ParseCargoBuild(raw []byte) *BuildResult {
	if len(raw) > MaxBuildBytes {
		return nil
	}
	result := &BuildResult{Tool: "cargo-build"}
	lines := bytes.Split(bytes.TrimRight(raw, "\n"), []byte("\n"))

	var current *Diagnostic

	for _, lb := range lines {
		line := string(lb)

		if m := cargoErrorCodeRe.FindStringSubmatch(line); m != nil {
			if current != nil {
				result.Diagnostics = append(result.Diagnostics, *current)
			}
			current = &Diagnostic{Tool: "cargo", Severity: SeverityError, Code: m[1], Message: m[2]}
			if len(result.Diagnostics) >= MaxDiagnostics {
				result.Truncated = true
				break
			}
			continue
		}
		if m := cargoErrorRe.FindStringSubmatch(line); m != nil {
			msg := m[1]
			if strings.HasPrefix(msg, "aborting due to") {
				continue
			}
			if current != nil {
				result.Diagnostics = append(result.Diagnostics, *current)
			}
			current = &Diagnostic{Tool: "cargo", Severity: SeverityError, Message: msg}
			if len(result.Diagnostics) >= MaxDiagnostics {
				result.Truncated = true
				break
			}
			continue
		}
		if m := cargoWarningRe.FindStringSubmatch(line); m != nil {
			if current != nil {
				result.Diagnostics = append(result.Diagnostics, *current)
			}
			code := m[1]
			msg := m[2]
			if strings.HasPrefix(msg, "unused import") || strings.HasPrefix(msg, "generated") {
				current = nil
				continue
			}
			current = &Diagnostic{Tool: "cargo", Severity: SeverityWarning, Code: code, Message: msg}
			if len(result.Diagnostics) >= MaxDiagnostics {
				result.Truncated = true
				break
			}
			continue
		}
		if current != nil {
			if m := cargoLocRe.FindStringSubmatch(line); m != nil {
				if current.File == "" {
					current.File = m[1]
					current.Line, _ = strconv.Atoi(m[2])
					current.Col, _ = strconv.Atoi(m[3])
				}
				continue
			}
			if m := cargoHelpRe.FindStringSubmatch(line); m != nil {
				kind := m[1]
				msg := m[2]
				sev := SeverityNote
				if kind == "help" {
					sev = SeverityHelp
				}
				current.Notes = append(current.Notes, DiagnosticNote{Severity: sev, Message: msg})
			}
		}
	}
	if current != nil {
		result.Diagnostics = append(result.Diagnostics, *current)
	}
	return result
}

func ParseCargoTest(raw []byte) *TestRun {
	if len(raw) > MaxBuildBytes {
		return nil
	}
	run := &TestRun{Tool: "cargo-test"}
	lines := bytes.Split(bytes.TrimRight(raw, "\n"), []byte("\n"))

	var currentPkg TestPackageResult
	var currentFailure *TestFailure
	inFailureSection := false
	var failedNames []string

	for _, lb := range lines {
		line := string(lb)

		if cargoTestRunningRe.MatchString(line) {
			if currentPkg.Name != "" {
				run.Packages = append(run.Packages, currentPkg)
			}
			currentPkg = TestPackageResult{}
			inFailureSection = false
			failedNames = nil
			continue
		}
		if m := cargoTestLineRe.FindStringSubmatch(line); m != nil {
			name := m[1]
			status := m[2]
			switch status {
			case "ok":
				currentPkg.Passed++
			case "ignored":
				currentPkg.Skipped++
			case "FAILED":
				currentPkg.Failed++
				failedNames = append(failedNames, name)
			}
			continue
		}
		if m := cargoTestResultRe.FindStringSubmatch(line); m != nil {
			currentPkg.Pass = m[1] == "ok"
			currentPkg.Passed, _ = strconv.Atoi(m[2])
			currentPkg.Failed, _ = strconv.Atoi(m[3])
			currentPkg.Skipped, _ = strconv.Atoi(m[4])
			run.Packages = append(run.Packages, currentPkg)
			currentPkg = TestPackageResult{}
			continue
		}
		if line == "failures:" {
			inFailureSection = true
			continue
		}
		if inFailureSection {
			if m := cargoTestSectionRe.FindStringSubmatch(line); m != nil {
				if currentFailure != nil {
					run.Failures = append(run.Failures, *currentFailure)
				}
				name := m[1]
				// Extract suite from name (e.g. "foo::bar::test_name" -> suite="foo::bar", name="test_name")
				parts := strings.Split(name, "::")
				testName := parts[len(parts)-1]
				suite := strings.Join(parts[:len(parts)-1], "::")
				currentFailure = &TestFailure{Suite: suite, Name: testName}
				continue
			}
			if currentFailure != nil && line != "" && line != "failures:" {
				if strings.HasPrefix(line, "thread '") && strings.Contains(line, "panicked") {
					currentFailure.Panic = true
				}
				if len(currentFailure.Output) < 30 {
					currentFailure.Output = append(currentFailure.Output, strings.TrimSpace(line))
				}
			}
		}
	}
	if currentFailure != nil {
		run.Failures = append(run.Failures, *currentFailure)
	}
	_ = failedNames
	return run
}
