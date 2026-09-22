package buildparse

import (
	"os"
	"testing"
)

func readFixture(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("readFixture %s: %v", path, err)
	}
	return data
}

func TestParseGoBuildError(t *testing.T) {
	raw := readFixture(t, "../../testdata/build/go-build-error.txt")
	result := ParseGoBuild(raw)
	if result == nil {
		t.Fatal("ParseGoBuild returned nil")
	}
	errors, warnings := 0, 0
	for _, d := range result.Diagnostics {
		switch d.Severity {
		case SeverityError, SeverityFatal:
			errors++
		case SeverityWarning:
			warnings++
		}
	}
	if errors != 2 {
		t.Errorf("expected 2 errors, got %d (diags: %+v)", errors, result.Diagnostics)
	}
	if warnings != 0 {
		t.Errorf("expected 0 warnings, got %d", warnings)
	}
}

func TestParseGoBuildEmpty(t *testing.T) {
	result := ParseGoBuild([]byte{})
	if result == nil {
		t.Fatal("ParseGoBuild(empty) returned nil")
	}
	if len(result.Diagnostics) != 0 {
		t.Errorf("expected 0 diagnostics, got %d", len(result.Diagnostics))
	}
}

func TestParseGoTestFail(t *testing.T) {
	raw := readFixture(t, "../../testdata/build/go-test-fail.txt")
	run := ParseGoTest(raw)
	if run == nil {
		t.Fatal("ParseGoTest returned nil")
	}
	if len(run.Failures) != 2 {
		t.Errorf("expected 2 failures, got %d (failures: %+v)", len(run.Failures), run.Failures)
	}
	passPkgs := 0
	for _, p := range run.Packages {
		if p.Pass {
			passPkgs++
		}
	}
	if passPkgs < 1 {
		t.Errorf("expected at least 1 pass package, got %d (packages: %+v)", passPkgs, run.Packages)
	}
}

func TestParseGoTestPass(t *testing.T) {
	raw := readFixture(t, "../../testdata/build/go-test-pass.txt")
	run := ParseGoTest(raw)
	if run == nil {
		t.Fatal("ParseGoTest returned nil")
	}
	if len(run.Packages) != 3 {
		t.Errorf("expected 3 packages, got %d (packages: %+v)", len(run.Packages), run.Packages)
	}
	if len(run.Failures) != 0 {
		t.Errorf("expected 0 failures, got %d", len(run.Failures))
	}
}

func TestParseGoTestPanic(t *testing.T) {
	raw := readFixture(t, "../../testdata/build/go-test-panic.txt")
	run := ParseGoTest(raw)
	if run == nil {
		t.Fatal("ParseGoTest returned nil")
	}
	if len(run.Failures) != 1 {
		t.Errorf("expected 1 failure, got %d (failures: %+v)", len(run.Failures), run.Failures)
	}
	if len(run.Failures) > 0 && !run.Failures[0].Panic {
		t.Errorf("expected Panic=true for TestWorkspace failure")
	}
}

func TestParseGCCSingleError(t *testing.T) {
	raw := readFixture(t, "../../testdata/build/gcc-single-error.txt")
	result := ParseGCC(raw)
	if result == nil {
		t.Fatal("ParseGCC returned nil")
	}
	errors := 0
	for _, d := range result.Diagnostics {
		if d.Severity == SeverityError || d.Severity == SeverityFatal {
			errors++
		}
	}
	if errors != 1 {
		t.Errorf("expected 1 error, got %d (diags: %+v)", errors, result.Diagnostics)
	}
	if len(result.Diagnostics) > 0 && len(result.Diagnostics[0].Notes) == 0 {
		t.Errorf("expected note on first diagnostic, got none")
	}
}

func TestParseGCCManyErrors(t *testing.T) {
	raw := readFixture(t, "../../testdata/build/gcc-many-errors.txt")
	result := ParseGCC(raw)
	if result == nil {
		t.Fatal("ParseGCC returned nil")
	}
	errors, warnings := 0, 0
	for _, d := range result.Diagnostics {
		switch d.Severity {
		case SeverityError, SeverityFatal:
			errors++
		case SeverityWarning:
			warnings++
		}
	}
	if errors != 3 {
		t.Errorf("expected 3 errors, got %d (diags: %+v)", errors, result.Diagnostics)
	}
	if warnings != 2 {
		t.Errorf("expected 2 warnings, got %d", warnings)
	}
}

func TestParseClangError(t *testing.T) {
	raw := readFixture(t, "../../testdata/build/clang-error.txt")
	result := ParseClang(raw)
	if result == nil {
		t.Fatal("ParseClang returned nil")
	}
	errors := 0
	for _, d := range result.Diagnostics {
		if d.Severity == SeverityError || d.Severity == SeverityFatal {
			errors++
		}
	}
	if errors != 1 {
		t.Errorf("expected 1 error, got %d (diags: %+v)", errors, result.Diagnostics)
	}
}

func TestParseCargoBuildErrors(t *testing.T) {
	raw := readFixture(t, "../../testdata/build/cargo-check-error.txt")
	result := ParseCargoBuild(raw)
	if result == nil {
		t.Fatal("ParseCargoBuild returned nil")
	}
	if len(result.Diagnostics) < 2 {
		t.Fatalf("expected at least 2 diagnostics, got %d (diags: %+v)", len(result.Diagnostics), result.Diagnostics)
	}
	codes := make(map[string]bool)
	for _, d := range result.Diagnostics {
		if d.Code != "" {
			codes[d.Code] = true
		}
	}
	if !codes["E0382"] {
		t.Errorf("expected E0382 in diagnostics, got codes: %v", codes)
	}
	if !codes["E0277"] {
		t.Errorf("expected E0277 in diagnostics, got codes: %v", codes)
	}
}

func TestParseCargoTestPass(t *testing.T) {
	raw := readFixture(t, "../../testdata/build/cargo-test-pass.txt")
	run := ParseCargoTest(raw)
	if run == nil {
		t.Fatal("ParseCargoTest returned nil")
	}
	if len(run.Failures) != 0 {
		t.Errorf("expected 0 failures, got %d (failures: %+v)", len(run.Failures), run.Failures)
	}
}

func TestParseCargoTestFail(t *testing.T) {
	raw := readFixture(t, "../../testdata/build/cargo-test-fail.txt")
	run := ParseCargoTest(raw)
	if run == nil {
		t.Fatal("ParseCargoTest returned nil")
	}
	if len(run.Failures) != 1 {
		t.Errorf("expected 1 failure, got %d (failures: %+v)", len(run.Failures), run.Failures)
	}
}
