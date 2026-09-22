package buildparse

import (
	"bytes"
	"regexp"
	"strconv"
	"strings"
)

var (
	goTestOkRe       = regexp.MustCompile(`^ok  \t(.+)\t([\d.]+)s`)
	goTestFailRe     = regexp.MustCompile(`^FAIL\t(.+)\t([\d.]+)s`)
	goTestCachedRe   = regexp.MustCompile(`^ok  \t(.+)\t\(cached\)`)
	goTestFailLineRe = regexp.MustCompile(`^--- FAIL: (\S+)`)
	goTestPassLineRe = regexp.MustCompile(`^--- PASS: (\S+)`)
	goTestFileLineRe = regexp.MustCompile(`^\s+(.+_test\.go):(\d+): (.*)`)
)

func ParseGoTest(raw []byte) *TestRun {
	if len(raw) > MaxBuildBytes {
		return nil
	}
	run := &TestRun{Tool: "go-test"}

	lines := bytes.Split(bytes.TrimRight(raw, "\n"), []byte("\n"))

	var currentFailure *TestFailure
	var currentPkg string

	for i := 0; i < len(lines); i++ {
		line := string(lines[i])

		if m := goTestOkRe.FindStringSubmatch(line); m != nil {
			elapsed, _ := strconv.ParseFloat(m[2], 64)
			run.Packages = append(run.Packages, TestPackageResult{Name: m[1], Pass: true, Elapsed: elapsed})
			currentPkg = m[1]
			continue
		}
		if m := goTestCachedRe.FindStringSubmatch(line); m != nil {
			run.Packages = append(run.Packages, TestPackageResult{Name: m[1], Pass: true})
			currentPkg = m[1]
			continue
		}
		if m := goTestFailRe.FindStringSubmatch(line); m != nil {
			elapsed, _ := strconv.ParseFloat(m[2], 64)
			run.Packages = append(run.Packages, TestPackageResult{Name: m[1], Pass: false, Elapsed: elapsed})
			currentPkg = m[1]
			continue
		}
		if m := goTestFailLineRe.FindStringSubmatch(line); m != nil {
			if currentFailure != nil {
				run.Failures = append(run.Failures, *currentFailure)
			}
			currentFailure = &TestFailure{Package: currentPkg, Name: m[1]}
			continue
		}
		if goTestPassLineRe.MatchString(line) {
			if currentFailure != nil {
				run.Failures = append(run.Failures, *currentFailure)
				currentFailure = nil
			}
			continue
		}
		if strings.HasPrefix(line, "panic:") {
			if currentFailure != nil {
				currentFailure.Panic = true
			}
			continue
		}
		// Indented failure output
		if currentFailure != nil && (strings.HasPrefix(line, "\t") || strings.HasPrefix(line, "    ")) {
			trimmed := strings.TrimSpace(line)
			if m := goTestFileLineRe.FindStringSubmatch(line); m != nil {
				if currentFailure.File == "" {
					currentFailure.File = m[1]
					currentFailure.Line, _ = strconv.Atoi(m[2])
				}
			}
			if len(currentFailure.Output) < 30 {
				currentFailure.Output = append(currentFailure.Output, trimmed)
			}
		}
	}
	if currentFailure != nil {
		run.Failures = append(run.Failures, *currentFailure)
	}
	return run
}
