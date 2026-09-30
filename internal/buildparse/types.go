package buildparse

import "fmt"

type Severity int

const (
	SeverityError Severity = iota
	SeverityWarning
	SeverityNote
	SeverityHelp
	SeverityFatal
)

func (s Severity) String() string {
	switch s {
	case SeverityError:
		return "E"
	case SeverityWarning:
		return "W"
	case SeverityNote:
		return "N"
	case SeverityHelp:
		return "H"
	case SeverityFatal:
		return "F"
	default:
		return "?"
	}
}

type DiagnosticNote struct {
	Severity Severity
	File     string
	Line     int
	Col      int
	Message  string
}

type Diagnostic struct {
	Tool     string
	Severity Severity
	Code     string
	File     string
	Line     int
	Col      int
	Message  string
	Notes    []DiagnosticNote
	RawStart int64
	RawEnd   int64
}

// Key returns a stable identity string for delta comparison.
func (d Diagnostic) Key() string {
	return fmt.Sprintf("%s:%d:%s", d.File, d.Line, d.Message)
}

type BuildResult struct {
	Tool        string
	Pass        bool
	Duration    string
	Diagnostics []Diagnostic
	Truncated   bool
}

type TestFailure struct {
	Package  string
	Suite    string
	Name     string
	File     string
	Line     int
	Output   []string
	Panic    bool
	RawStart int64
	RawEnd   int64
}

func (f *TestFailure) FullName() string {
	if f.Package != "" {
		return f.Package + "::" + f.Name
	}
	return f.Name
}

type TestPackageResult struct {
	Name    string
	Pass    bool
	Elapsed float64
	Passed  int
	Failed  int
	Skipped int
}

type TestRun struct {
	Tool     string
	Pass     bool
	Packages []TestPackageResult
	Failures []TestFailure
	Duration string
}

const MaxDiagnostics = 200
const MaxBuildBytes = 50 * 1024 * 1024
