package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/taqu/agentcap/internal/bench"
	"github.com/taqu/agentcap/internal/buildparse"
	"github.com/taqu/agentcap/internal/delta"
	"github.com/taqu/agentcap/internal/exec"
	"github.com/taqu/agentcap/internal/gitparse"
	"github.com/taqu/agentcap/internal/integration/claude"
	"github.com/taqu/agentcap/internal/integration/codex"
	"github.com/taqu/agentcap/internal/integration/common"
	"github.com/taqu/agentcap/internal/integration/engine"
	"github.com/taqu/agentcap/internal/integration/protocol"
	"github.com/taqu/agentcap/internal/project"
	"github.com/taqu/agentcap/internal/query"
	"github.com/taqu/agentcap/internal/reduce"
	"github.com/taqu/agentcap/internal/session"
	"github.com/taqu/agentcap/internal/stats"
	"github.com/taqu/agentcap/internal/store"
)

const usage = `acap — command output compression layer for AI coding agents

Usage:
  acap run <command> [args...]
  acap show <id> [--meta] [--lines X:Y] [--match <text>] [--path <path>]
  acap raw  <id> [--stdout|--stderr]
  acap clean [--older-than <duration>]
  acap stats
  acap bench command [--json] -- <command> [args...]
  acap session <start|info|list|history>
  acap history
  acap exec [--protocol=json]
  acap hook <claude|codex>
  acap integrate <claude|codex|status> [--dry-run] [--remove]
  acap doctor

Examples:
  acap run ls -la
  acap run find . -type f
  acap run rg "Workspace" .
  acap run cat src/main.go
  acap show 8f31c2
  acap show 8f31c2 --lines 1:20
  acap show 8f31c2 --match "func "
  acap raw  8f31c2 --stderr
  acap clean --older-than 7d
  acap bench command -- git diff
  acap session start
  acap session info
  acap session list
  acap history
  acap integrate claude
  acap integrate status
  acap doctor

Environment:
  ACAP_DEBUG=1        Print reducer debug info to stderr.
  ACAP_SESSION_ID=xx  Enable session-aware delta compression.
  ACAP_ROOT=<path>    Override project root for .acap/store.db discovery.
  ACAP_BYPASS=1       Skip AgentCap interception entirely.
`

const runUsage = `Usage: acap run <command> [args...]

Executes <command> directly (no shell interpretation), reduces the output
for LLM consumption, and writes it to stdout.  The child's exit code is
preserved.
`

const benchUsage = `Usage: acap bench <subcommand>

Subcommands:
  command   Benchmark a single command execution.

Run "acap bench command --help" for details.
`

const benchCommandUsage = `Usage: acap bench command [--json] -- <command> [args...]

Executes <command> exactly once through the normal AgentCap pipeline
(no shell interpretation) and reports:

  raw        bytes the command wrote to stdout and stderr
  visible    bytes of the AgentCap result presentation an agent receives
  reduction  1 - visible/raw for this single command ("n/a" if raw is 0)
  duration   command execution time (process start to exit)
  processing AgentCap post-execution time, including reduction and storage
  id         stored result ID, inspectable with "acap show" / "acap raw"

Reduction is the reduction in agent-visible bytes for this one command.
It is not a workflow-level or end-to-end agent savings figure.

The measurement is stateless: no session delta is applied, even if
ACAP_SESSION_ID is set. The benchmark exits with the command's exit code.

Flags:
  --json    Print the measurement as JSON.
`

var debugMode = os.Getenv("ACAP_DEBUG") == "1"

func main() {
	args := os.Args[1:]

	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprint(os.Stdout, usage)
		os.Exit(0)
	}

	switch args[0] {
	case "run":
		runCmd(args[1:])
	case "show":
		showCmd(args[1:])
	case "raw":
		rawCmd(args[1:])
	case "clean":
		cleanCmd(args[1:])
	case "stats":
		statsCmd()
	case "bench":
		benchCmd(args[1:])
	case "session":
		sessionCmd(args[1:])
	case "history":
		historyCmd(args[1:])
	case "exec":
		execCmd(args[1:])
	case "hook":
		hookCmd(args[1:])
	case "integrate":
		integrateCmd(args[1:])
	case "doctor":
		doctorCmd(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "acap: unknown command %q\n\n%s", args[0], usage)
		os.Exit(1)
	}
}

func runCmd(args []string) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprint(os.Stdout, runUsage)
		os.Exit(0)
	}

	startTime := time.Now().UTC()

	result, err := exec.Run(context.Background(), args, nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(result.ExitCode)
	}

	execDuration := result.Duration

	t0 := time.Now()
	reducer := reduce.Select(args)
	reduced := reducer.Reduce(result)
	reduceDuration := time.Since(t0)

	// Compute hashes.
	stdoutHash := delta.HashBytes(result.Stdout)
	stderrHash := delta.HashBytes(result.Stderr)

	// Get working directory.
	cwd, _ := os.Getwd()

	// Build metadata.
	reducerName := reducerTypeName(reducer)
	meta := store.Meta{
		Command:     args,
		ExitCode:    result.ExitCode,
		StartedAt:   startTime,
		DurationMs:  result.Duration.Milliseconds(),
		StdoutBytes: int64(len(result.Stdout)),
		StderrBytes: int64(len(result.Stderr)),
		Reducer:     reducerName,
		CreatedAt:   time.Now().UTC(),
		Truncated:   result.Truncated,
		StdoutHash:  stdoutHash,
		StderrHash:  stderrHash,
		WorkDir:     cwd,
	}

	// Find session and baseline.
	var sess *session.Session
	var baselineEntry *store.Entry
	var baselineRec *session.HistoryRecord

	if s, err := session.Current(); err == nil && s != nil {
		sess = s
		meta.SessionID = s.ID
		key := session.NewKey(args, cwd)
		if rec, err := s.LatestBaseline(key); err == nil && rec != nil {
			baselineRec = rec
		}
	}

	// Store result.
	st, storeInitErr := store.New()
	if storeInitErr == nil {
		defer st.Close()
	}
	var entry *store.Entry
	if storeInitErr == nil {
		entry, err = st.Save(meta, result.Stdout, result.Stderr, reduced.Output)
		if err != nil {
			fmt.Fprintf(os.Stderr, "acap: warning: result not stored: %v\n", err)
			entry = nil
		}
	} else {
		fmt.Fprintf(os.Stderr, "acap: warning: result not stored: %v\n", storeInitErr)
	}

	// Load baseline entry if we have a reference.
	if sess != nil && baselineRec != nil && st != nil {
		if be, err := st.Open(baselineRec.ResultID); err == nil {
			baselineEntry = be
		}
		// Missing baseline is OK: fall back to full capsule.
	}

	// Determine presentation.
	output := reduced.Output
	presentation := string(delta.PresentationFull)

	// Save git diff metadata if applicable.
	if entry != nil && storeInitErr == nil {
		if gdr, ok := reducer.(*reduce.GitDiffReducer); ok && gdr.ParsedDiff != nil {
			if err := st.SaveGitDiff(entry.ID, gdr.ParsedDiff.Files); err != nil && debugMode {
				fmt.Fprintf(os.Stderr, "acap: warning: save git diff: %v\n", err)
			}
		}
	}

	// Save build diagnostics / test failures if applicable.
	if entry != nil && st != nil {
		switch rv := reducer.(type) {
		case *reduce.GoBuildReducer:
			if rv.ParsedBuild != nil && len(rv.ParsedBuild.Diagnostics) > 0 {
				if err := st.SaveDiagnostics(entry.ID, rv.ParsedBuild.Diagnostics); err != nil && debugMode {
					fmt.Fprintf(os.Stderr, "acap: warning: save diagnostics: %v\n", err)
				}
			}
		case *reduce.GoTestReducer:
			if rv.ParsedRun != nil && len(rv.ParsedRun.Failures) > 0 {
				if err := st.SaveTestFailures(entry.ID, rv.ParsedRun.Failures); err != nil && debugMode {
					fmt.Fprintf(os.Stderr, "acap: warning: save test failures: %v\n", err)
				}
			}
		case *reduce.GccReducer:
			if rv.ParsedBuild != nil && len(rv.ParsedBuild.Diagnostics) > 0 {
				if err := st.SaveDiagnostics(entry.ID, rv.ParsedBuild.Diagnostics); err != nil && debugMode {
					fmt.Fprintf(os.Stderr, "acap: warning: save diagnostics: %v\n", err)
				}
			}
		case *reduce.CargoBuildReducer:
			if rv.ParsedBuild != nil && len(rv.ParsedBuild.Diagnostics) > 0 {
				if err := st.SaveDiagnostics(entry.ID, rv.ParsedBuild.Diagnostics); err != nil && debugMode {
					fmt.Fprintf(os.Stderr, "acap: warning: save diagnostics: %v\n", err)
				}
			}
		case *reduce.CargoTestReducer:
			if rv.ParsedRun != nil && len(rv.ParsedRun.Failures) > 0 {
				if err := st.SaveTestFailures(entry.ID, rv.ParsedRun.Failures); err != nil && debugMode {
					fmt.Fprintf(os.Stderr, "acap: warning: save test failures: %v\n", err)
				}
			}
		}
	}

	if sess != nil && baselineEntry != nil && entry != nil {
		dr := selectDeltaReducer(reducer, st)
		deltaResult := delta.Compare(
			context.Background(),
			baselineEntry,
			entry,
			reduced.Output,
			dr,
			stdoutHash,
			stderrHash,
		)
		output = deltaResult.Output
		presentation = string(deltaResult.Presentation)

		// Update stored metadata with baseline and presentation info.
		updatedMeta := entry.Meta
		updatedMeta.BaselineID = baselineEntry.Meta.ID
		updatedMeta.Presentation = presentation
		if err := st.UpdateMeta(entry, updatedMeta); err != nil && debugMode {
			fmt.Fprintf(os.Stderr, "acap: warning: update meta: %v\n", err)
		}
	} else if entry != nil {
		updatedMeta := entry.Meta
		updatedMeta.Presentation = "full"
		_ = st.UpdateMeta(entry, updatedMeta)
	}

	// Inject result ID.
	if entry != nil {
		output = injectResultID(output, entry.ID)
	}
	fmt.Print(output)

	// Record in session history.
	if sess != nil && entry != nil {
		key := session.NewKey(args, cwd)
		_ = sess.Record(entry.ID, key, result.ExitCode, stdoutHash, stderrHash, presentation)
	}

	// Persist statistics.
	stateless := reduced.RetBytes
	stateful := len(output)
	if st != nil {
		_ = st.RecordRun(reduced.RawBytes, stateless, stateful, presentation)
	}

	if debugMode {
		var idStr string
		if entry != nil {
			idStr = entry.ID
		} else {
			idStr = "(not stored)"
		}
		var sessionStr string
		if sess != nil {
			sessionStr = sess.ID
		}
		fmt.Fprintf(os.Stderr, "acap: cmd=%s reducer=%s id=%s session=%s presentation=%s raw=%d stateless=%d stateful=%d exec=%s reduce=%s\n",
			args[0], reducerName, idStr, sessionStr, presentation,
			reduced.RawBytes, stateless, stateful,
			execDuration.Round(time.Millisecond),
			reduceDuration.Round(time.Millisecond),
		)
	}

	os.Exit(result.ExitCode)
}

// selectDeltaReducer maps a reduce.Reducer to its corresponding delta.Reducer.
func selectDeltaReducer(r reduce.Reducer, st *store.Store) delta.Reducer {
	switch r.(type) {
	case *reduce.GrepReducer:
		return &delta.GrepDelta{}
	case *reduce.FindReducer:
		return &delta.FindDelta{}
	case *reduce.LsReducer:
		return &delta.LsDelta{}
	case *reduce.GitStatusReducer:
		return &delta.GitStatusDelta{}
	case *reduce.GitDiffReducer:
		return &delta.GitDiffDelta{Store: st}
	case *reduce.GoBuildReducer:
		return &delta.DiagnosticsDelta{Store: st}
	case *reduce.GoTestReducer:
		return &delta.TestResultsDelta{Store: st}
	case *reduce.GccReducer:
		return &delta.DiagnosticsDelta{Store: st}
	case *reduce.CargoBuildReducer:
		return &delta.DiagnosticsDelta{Store: st}
	case *reduce.CargoTestReducer:
		return &delta.TestResultsDelta{Store: st}
	default:
		return nil
	}
}

// injectResultID inserts the result ID into the output header.
// If output starts with "@acap ", the ID is inserted as the second token.
// Otherwise "@acap <id>" is prepended.
func injectResultID(output, id string) string {
	if strings.HasPrefix(output, "@acap ") {
		// Insert id after "@acap ".
		rest := output[len("@acap "):]
		return "@acap " + id + " " + rest
	}
	return "@acap " + id + "\n" + output
}

func showCmd(args []string) {
	fs := flag.NewFlagSet("show", flag.ExitOnError)
	metaFlag := fs.Bool("meta", false, "print metadata")
	linesFlag := fs.String("lines", "", "line range X:Y (1-based inclusive)")
	matchFlag := fs.String("match", "", "search text (case-insensitive)")
	pathFlag := fs.String("path", "", "path substring filter")
	fileFlag := fs.String("file", "", "show file diff (git diff results)")
	hunkFlag := fs.Int("hunk", 0, "show specific hunk (requires --file, 1-based)")
	errorsFlag := fs.Bool("errors", false, "show errors for build result")
	warningsFlag := fs.Bool("warnings", false, "show warnings for build result")
	testFlag := fs.String("test", "", "show specific test failure")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: acap show <id> [--meta] [--lines X:Y] [--match text] [--path path] [--file path] [--hunk N] [--errors] [--warnings] [--test name]")
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}
	remaining := fs.Args()
	if len(remaining) == 0 {
		fs.Usage()
		os.Exit(1)
	}
	id := remaining[0]

	s, err := store.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: show: %v\n", err)
		os.Exit(1)
	}
	defer s.Close()
	entry, err := s.Open(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: show: %v\n", err)
		os.Exit(1)
	}

	var written int

	switch {
	case *fileFlag != "":
		written = showGitFile(s, entry, *fileFlag, *hunkFlag)
	case *errorsFlag:
		written = showBuildErrors(s, entry)
	case *warningsFlag:
		written = showBuildWarnings(s, entry)
	case *testFlag != "":
		written = showTestFailure(s, entry, *testFlag)
	case *metaFlag:
		written = showMeta(entry)
	case *linesFlag != "":
		written = showLines(entry, *linesFlag)
	case *matchFlag != "":
		written = showMatch(entry, *matchFlag)
	case *pathFlag != "":
		written = showPath(entry, *pathFlag)
	default:
		cap, err := entry.Capsule()
		if err != nil {
			fmt.Fprintf(os.Stderr, "acap: show: %v\n", err)
			os.Exit(1)
		}
		fmt.Print(cap)
		written = len(cap)
	}

	_ = s.RecordShow(written)
}

func showGitFile(st *store.Store, entry *store.Entry, filePath string, hunkN int) int {
	files, err := st.GetGitDiffFiles(entry.ID)
	if err != nil || len(files) == 0 {
		fmt.Fprintf(os.Stderr, "acap: show --file: no git diff data for %s\n", entry.ID)
		os.Exit(1)
	}

	// Find file by exact or suffix match.
	var matched *gitparse.GitDiffFile
	for i := range files {
		f := &files[i]
		p := f.Path()
		if p == filePath || strings.HasSuffix(p, "/"+filePath) || strings.HasSuffix(p, "\\"+filePath) {
			matched = f
			break
		}
	}
	if matched == nil {
		fmt.Fprintf(os.Stderr, "acap: show --file: file %q not found in git diff result %s\n", filePath, entry.ID)
		os.Exit(1)
	}

	rawFile, err := os.Open(entry.StdoutPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: show --file: %v\n", err)
		os.Exit(1)
	}
	defer rawFile.Close()

	if hunkN > 0 {
		// Look up hunk.
		hunks, err := st.GetGitDiffHunks(entry.ID, matched.Index)
		if err != nil {
			fmt.Fprintf(os.Stderr, "acap: show --file --hunk: %v\n", err)
			os.Exit(1)
		}
		if hunkN > len(hunks) {
			fmt.Fprintf(os.Stderr, "acap: show --file --hunk: hunk %d out of range (file has %d hunks)\n", hunkN, len(hunks))
			os.Exit(1)
		}
		h := hunks[hunkN-1]
		size := h.RawEnd - h.RawStart
		if size <= 0 {
			fmt.Fprintf(os.Stderr, "acap: show --file --hunk: empty hunk\n")
			os.Exit(1)
		}
		if _, err := rawFile.Seek(h.RawStart, io.SeekStart); err != nil {
			fmt.Fprintf(os.Stderr, "acap: show --file --hunk: seek: %v\n", err)
			os.Exit(1)
		}
		n, err := io.Copy(os.Stdout, io.LimitReader(rawFile, size))
		if err != nil {
			fmt.Fprintf(os.Stderr, "acap: show --file --hunk: read: %v\n", err)
			os.Exit(1)
		}
		return int(n)
	}

	// Read entire file range.
	size := matched.RawEnd - matched.RawStart
	if size <= 0 {
		fmt.Fprintf(os.Stderr, "acap: show --file: empty file range\n")
		os.Exit(1)
	}
	if _, err := rawFile.Seek(matched.RawStart, io.SeekStart); err != nil {
		fmt.Fprintf(os.Stderr, "acap: show --file: seek: %v\n", err)
		os.Exit(1)
	}
	n, err := io.Copy(os.Stdout, io.LimitReader(rawFile, size))
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: show --file: read: %v\n", err)
		os.Exit(1)
	}
	return int(n)
}

func showBuildErrors(st *store.Store, entry *store.Entry) int {
	diags, err := st.GetDiagnostics(entry.ID)
	if err != nil || len(diags) == 0 {
		fmt.Fprintf(os.Stderr, "acap: show --errors: no diagnostic data for %s\n", entry.ID)
		os.Exit(1)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap %s errors\n\n", entry.ID)
	for _, d := range diags {
		if d.Severity != buildparse.SeverityError && d.Severity != buildparse.SeverityFatal {
			continue
		}
		if d.File == "(linker)" {
			fmt.Fprintf(&sb, "E (linker): %s\n", d.Message)
		} else {
			code := ""
			if d.Code != "" {
				code = "[" + d.Code + "] "
			}
			fmt.Fprintf(&sb, "E%s %s:%d:%d\n  %s\n", code, d.File, d.Line, d.Col, d.Message)
		}
	}
	out := sb.String()
	fmt.Print(out)
	return len(out)
}

func showBuildWarnings(st *store.Store, entry *store.Entry) int {
	diags, err := st.GetDiagnostics(entry.ID)
	if err != nil || len(diags) == 0 {
		fmt.Fprintf(os.Stderr, "acap: show --warnings: no diagnostic data for %s\n", entry.ID)
		os.Exit(1)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap %s warnings\n\n", entry.ID)
	for _, d := range diags {
		if d.Severity != buildparse.SeverityWarning {
			continue
		}
		code := ""
		if d.Code != "" {
			code = "[" + d.Code + "] "
		}
		fmt.Fprintf(&sb, "W%s %s:%d %s\n", code, d.File, d.Line, d.Message)
	}
	out := sb.String()
	fmt.Print(out)
	return len(out)
}

func showTestFailure(st *store.Store, entry *store.Entry, testName string) int {
	failures, err := st.GetTestFailures(entry.ID)
	if err != nil || len(failures) == 0 {
		fmt.Fprintf(os.Stderr, "acap: show --test: no test failure data for %s\n", entry.ID)
		os.Exit(1)
	}
	var matched *buildparse.TestFailure
	for i := range failures {
		f := &failures[i]
		if f.FullName() == testName || f.Name == testName {
			matched = f
			break
		}
	}
	if matched == nil {
		fmt.Fprintf(os.Stderr, "acap: show --test: test %q not found in result %s\n", testName, entry.ID)
		os.Exit(1)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap %s test=%s\n", entry.ID, matched.FullName())
	if matched.File != "" {
		fmt.Fprintf(&sb, "%s:%d\n", matched.File, matched.Line)
	}
	if matched.Panic {
		sb.WriteString("PANIC\n")
	}
	for _, l := range matched.Output {
		fmt.Fprintf(&sb, "%s\n", l)
	}
	out := sb.String()
	fmt.Print(out)
	return len(out)
}

func showMeta(entry *store.Entry) int {
	m := entry.Meta
	cmd := strings.Join(m.Command, " ")
	out := fmt.Sprintf("id=%s\ncommand=%s\nexit=%d\nduration=%dms\nstdout=%dB\nstderr=%dB\nreducer=%s\ncreated=%s\n",
		m.ID, cmd, m.ExitCode, m.DurationMs, m.StdoutBytes, m.StderrBytes, m.Reducer,
		m.CreatedAt.Format(time.RFC3339),
	)
	fmt.Print(out)
	return len(out)
}

func showLines(entry *store.Entry, linesFlag string) int {
	parts := strings.SplitN(linesFlag, ":", 2)
	if len(parts) != 2 {
		fmt.Fprintf(os.Stderr, "acap: show --lines: expected X:Y format, got %q\n", linesFlag)
		os.Exit(1)
	}
	from, err1 := strconv.Atoi(parts[0])
	to, err2 := strconv.Atoi(parts[1])
	if err1 != nil || err2 != nil {
		fmt.Fprintf(os.Stderr, "acap: show --lines: invalid range %q\n", linesFlag)
		os.Exit(1)
	}

	f, err := os.Open(entry.StdoutPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: show --lines: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	lines, err := query.Lines(f, from, to)
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: show --lines: %v\n", err)
		os.Exit(1)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap %s lines=%d:%d\n\n", entry.ID, from, to)
	for _, l := range lines {
		sb.WriteString(l)
		sb.WriteByte('\n')
	}
	out := sb.String()
	fmt.Print(out)
	return len(out)
}

func showMatch(entry *store.Entry, text string) int {
	f, err := os.Open(entry.StdoutPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: show --match: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	res, err := query.Match(f, text, 50)
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: show --match: %v\n", err)
		os.Exit(1)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap %s match=%q hits=%d\n\n", entry.ID, text, res.Total)
	for _, h := range res.Hits {
		fmt.Fprintf(&sb, "%4d: %s\n", h.LineNum, h.Line)
	}
	if res.Omitted > 0 {
		fmt.Fprintf(&sb, "\nomitted=%d\n", res.Omitted)
	}
	out := sb.String()
	fmt.Print(out)
	return len(out)
}

func showPath(entry *store.Entry, path string) int {
	f, err := os.Open(entry.StdoutPath())
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: show --path: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	lines, total, err := query.PathFilter(f, path, 100)
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: show --path: %v\n", err)
		os.Exit(1)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap %s path=%s matches=%d\n\n", entry.ID, path, total)
	for _, l := range lines {
		sb.WriteString(l)
		sb.WriteByte('\n')
	}
	out := sb.String()
	fmt.Print(out)
	return len(out)
}

func rawCmd(args []string) {
	fs := flag.NewFlagSet("raw", flag.ExitOnError)
	stdoutFlag := fs.Bool("stdout", false, "print stdout (default)")
	stderrFlag := fs.Bool("stderr", false, "print stderr")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: acap raw <id> [--stdout|--stderr]")
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}
	remaining := fs.Args()
	if len(remaining) == 0 {
		fs.Usage()
		os.Exit(1)
	}
	id := remaining[0]

	s, err := store.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: raw: %v\n", err)
		os.Exit(1)
	}
	defer s.Close()
	entry, err := s.Open(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: raw: %v\n", err)
		os.Exit(1)
	}

	var path string
	if *stderrFlag {
		path = entry.StderrPath()
	} else if *stdoutFlag {
		path = entry.StdoutPath()
	} else {
		path = entry.StdoutPath() // default: stdout
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: raw: %v\n", err)
		os.Exit(1)
	}
	defer f.Close()

	n, err := io.Copy(os.Stdout, f)
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: raw: copy: %v\n", err)
		os.Exit(1)
	}
	_ = s.RecordRaw(int(n))
}

func cleanCmd(args []string) {
	fs := flag.NewFlagSet("clean", flag.ExitOnError)
	olderThan := fs.String("older-than", "7d", "remove results older than this duration (e.g. 24h, 30d)")
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: acap clean [--older-than <duration>]")
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	dur, err := parseDuration(*olderThan)
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: clean: invalid duration %q: %v\n", *olderThan, err)
		os.Exit(1)
	}

	s, err := store.New()
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: clean: %v\n", err)
		os.Exit(1)
	}
	defer s.Close()

	removed, err := s.Cleanup(dur)
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: clean: %v\n", err)
		os.Exit(1)
	}

	if removed == 0 {
		fmt.Println("nothing to remove")
	} else {
		fmt.Printf("removed %d result(s)\n", removed)
	}
}

// parseDuration extends time.ParseDuration to support a "d" suffix for days.
func parseDuration(s string) (time.Duration, error) {
	if strings.HasSuffix(s, "d") {
		n, err := strconv.ParseFloat(s[:len(s)-1], 64)
		if err != nil {
			return 0, fmt.Errorf("invalid days value: %w", err)
		}
		return time.Duration(n * float64(24*time.Hour)), nil
	}
	return time.ParseDuration(s)
}

func statsCmd() {
	cwd, _ := os.Getwd()
	root := project.FindRoot(cwd)
	st, err := store.Open(root)
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: stats: %v\n", err)
		os.Exit(1)
	}
	defer st.Close()

	s, err := st.LoadStats()
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: stats: %v\n", err)
		os.Exit(1)
	}

	saved := s.RawBytes - s.RetBytes
	var reduction float64
	if s.RawBytes > 0 {
		reduction = float64(saved) / float64(s.RawBytes) * 100
	}

	fmt.Printf("commands: %d\n", s.Commands)
	fmt.Printf("raw:      %s\n", stats.FormatBytes(s.RawBytes))
	fmt.Printf("returned: %s\n", stats.FormatBytes(s.RetBytes))
	fmt.Printf("saved:    %s\n", stats.FormatBytes(saved))
	fmt.Printf("reduction: %.1f%%\n", reduction)

	if s.ShowCalls > 0 {
		fmt.Printf("show_calls: %d\n", s.ShowCalls)
		fmt.Printf("show_returned: %s\n", stats.FormatBytes(s.ShowRetBytes))
	}
	if s.RawCalls > 0 {
		fmt.Printf("raw_calls: %d\n", s.RawCalls)
		fmt.Printf("raw_returned: %s\n", stats.FormatBytes(s.RawRetBytes))
	}

	// Session stats.
	if s.UnchangedCount+s.DeltaCount+s.FullFallbackCount > 0 {
		fmt.Println()
		fmt.Printf("unchanged: %d\n", s.UnchangedCount)
		fmt.Printf("delta:     %d\n", s.DeltaCount)
		fmt.Printf("full:      %d\n", s.FullFallbackCount)
	}
	if s.StatelessBytes > 0 && s.StatefulBytes > 0 {
		sessionSaved := s.StatelessBytes - s.StatefulBytes
		var sessionReduction float64
		if s.StatelessBytes > 0 {
			sessionReduction = float64(sessionSaved) / float64(s.StatelessBytes) * 100
		}
		fmt.Printf("\nstateless: %s\n", stats.FormatBytes(s.StatelessBytes))
		fmt.Printf("stateful:  %s\n", stats.FormatBytes(s.StatefulBytes))
		fmt.Printf("session_saved: %s (%.1f%%)\n", stats.FormatBytes(sessionSaved), sessionReduction)
	}
}

func benchCmd(args []string) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprint(os.Stdout, benchUsage)
		os.Exit(0)
	}
	switch args[0] {
	case "command":
		benchCommandCmd(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "acap: bench: unknown subcommand %q\n\n%s", args[0], benchUsage)
		os.Exit(1)
	}
}

func benchCommandCmd(args []string) {
	fs := flag.NewFlagSet("bench command", flag.ExitOnError)
	jsonFlag := fs.Bool("json", false, "print measurement as JSON")
	fs.Usage = func() { fmt.Fprint(os.Stderr, benchCommandUsage) }
	if len(args) > 0 && (args[0] == "--help" || args[0] == "-h") {
		fmt.Fprint(os.Stdout, benchCommandUsage)
		os.Exit(0)
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}
	target := fs.Args()
	if len(target) == 0 {
		fs.Usage()
		os.Exit(1)
	}

	cwd, _ := os.Getwd()
	m, err := bench.MeasureCommand(context.Background(), target, cwd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: bench: %v\n", err)
		var ee *bench.ExecError
		if errors.As(err, &ee) && ee.ExitCode != 0 {
			os.Exit(ee.ExitCode)
		}
		os.Exit(1)
	}

	if *jsonFlag {
		err = bench.WriteJSON(os.Stdout, m)
	} else {
		err = bench.WriteText(os.Stdout, m)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: bench: %v\n", err)
		os.Exit(1)
	}
	os.Exit(m.ExitCode)
}

func sessionCmd(args []string) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		fmt.Fprintln(os.Stdout, "Usage: acap session <start|info|list|history>")
		os.Exit(0)
	}
	switch args[0] {
	case "start":
		id, err := session.GenerateID()
		if err != nil {
			fmt.Fprintf(os.Stderr, "acap: session start: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("agentcap-session=%s\n", id)
		fmt.Printf("export %s=%s\n", session.EnvKey, id)
	case "info":
		sess, err := session.Current()
		if err != nil || sess == nil {
			fmt.Fprintf(os.Stderr, "acap: no session (set %s)\n", session.EnvKey)
			os.Exit(1)
		}
		records, _ := sess.History().Records()
		fmt.Printf("session=%s\ncommands=%d\n", sess.ID, len(records))
	case "list":
		sessions, err := session.List()
		if err != nil {
			fmt.Fprintf(os.Stderr, "acap: session list: %v\n", err)
			os.Exit(1)
		}
		if len(sessions) == 0 {
			fmt.Println("no sessions")
			return
		}
		for _, si := range sessions {
			age := time.Since(si.Updated).Round(time.Minute)
			fmt.Printf("%s results=%d updated=%s\n", si.ID, si.Count, formatAge(age))
		}
	case "history":
		historyCmd(args[1:])
	default:
		fmt.Fprintf(os.Stderr, "acap: session: unknown subcommand %q\n", args[0])
		os.Exit(1)
	}
}

func historyCmd(args []string) {
	sess, err := session.Current()
	if err != nil || sess == nil {
		fmt.Fprintf(os.Stderr, "acap: no session (set %s)\n", session.EnvKey)
		os.Exit(1)
	}
	records, err := sess.History().Records()
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: history: %v\n", err)
		os.Exit(1)
	}
	if len(records) == 0 {
		fmt.Println("no history")
		return
	}
	for i := len(records) - 1; i >= 0; i-- {
		r := records[i]
		// We don't store command text in history, just the key hash.
		// Show result ID and presentation.
		fmt.Printf("%d %s ? %s\n", r.Seq, r.ResultID, r.Presentation)
	}
}

func formatAge(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}

// execCmd implements "acap exec [--protocol=json]".
// Reads a ToolRequest JSON from stdin and writes a ToolResponse JSON to stdout.
func execCmd(args []string) {
	fs := flag.NewFlagSet("exec", flag.ExitOnError)
	protocolFlag := fs.String("protocol", "json", "protocol format (only json supported)")
	_ = fs.Parse(args)

	if *protocolFlag != "json" {
		writeToolError(fmt.Sprintf("unsupported protocol: %s", *protocolFlag))
		os.Exit(1)
	}

	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		writeToolError(fmt.Sprintf("read stdin: %v", err))
		os.Exit(1)
	}

	var req protocol.ToolRequest
	if err := json.Unmarshal(data, &req); err != nil {
		writeToolError(fmt.Sprintf("parse request: %v", err))
		os.Exit(1)
	}

	resp, err := engine.Execute(context.Background(), &req)
	if err != nil {
		writeToolError(err.Error())
		os.Exit(1)
	}

	out, _ := json.Marshal(resp)
	fmt.Println(string(out))
	os.Exit(resp.ExitCode)
}

func writeToolError(msg string) {
	resp := protocol.ToolResponse{
		Protocol: protocol.Version,
		ExitCode: 1,
		Error:    msg,
	}
	out, _ := json.Marshal(resp)
	fmt.Println(string(out))
}

// shellMetaChars contains characters that indicate complex shell expressions.
const shellMetaChars = "|;&$><`\\"

// hasShellMeta returns true if the command string contains shell metacharacters.
func hasShellMeta(cmd string) bool {
	return strings.ContainsAny(cmd, shellMetaChars)
}

// splitWords performs a simple whitespace split of a command string.
// Does not handle quoting or escaping — only used for simple commands.
func splitWords(s string) []string {
	return strings.Fields(s)
}

// hookCmd implements "acap hook <claude|codex>".
func hookCmd(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: acap hook <claude|codex>")
		os.Exit(1)
	}
	adapter := args[0]

	// Check bypass/recursive conditions first.
	if common.IsBypassed() || common.IsRecursive() {
		switch adapter {
		case "claude":
			os.Stdout.Write(claude.MakeAllowResponse())
		case "codex":
			os.Stdout.Write(codex.MakeAllowResponse())
		}
		os.Exit(0)
	}

	switch adapter {
	case "claude":
		hookClaudeCmd()
	case "codex":
		hookCodexCmd()
	default:
		fmt.Fprintf(os.Stderr, "acap: hook: unknown adapter %q\n", adapter)
		os.Exit(1)
	}
}

func hookClaudeCmd() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		os.Stdout.Write(claude.MakeAllowResponse())
		os.Exit(0)
	}

	hookInput, err := claude.ParseHookInput(data)
	if err != nil || hookInput.ToolName != "Bash" {
		os.Stdout.Write(claude.MakeAllowResponse())
		os.Exit(0)
	}

	bashInput, err := claude.ParseBashInput(hookInput.ToolInput)
	if err != nil {
		os.Stdout.Write(claude.MakeAllowResponse())
		os.Exit(0)
	}

	cmd := bashInput.Command

	// Bypass for shell metacharacters — let original Bash tool handle complex expressions.
	if hasShellMeta(cmd) {
		os.Stdout.Write(claude.MakeAllowResponse())
		os.Exit(0)
	}

	cmdArgs := splitWords(cmd)
	if len(cmdArgs) == 0 {
		os.Stdout.Write(claude.MakeAllowResponse())
		os.Exit(0)
	}

	// Bypass for acap commands to prevent recursion.
	if common.IsAcapCommand(cmdArgs) {
		os.Stdout.Write(claude.MakeAllowResponse())
		os.Exit(0)
	}

	cwd, _ := os.Getwd()
	sessionID := common.MapAgentSession(hookInput.SessionID)

	req := &protocol.ToolRequest{
		Protocol:   protocol.Version,
		Command:    cmdArgs,
		WorkingDir: cwd,
		SessionID:  sessionID,
	}

	resp, err := engine.Execute(context.Background(), req)
	if err != nil {
		os.Stdout.Write(claude.MakeAllowResponse())
		os.Exit(0)
	}

	blockResp, err := claude.MakeBlockResponse(resp.Stdout)
	if err != nil {
		os.Stdout.Write(claude.MakeAllowResponse())
		os.Exit(0)
	}

	os.Stdout.Write(blockResp)
	os.Exit(2)
}

func hookCodexCmd() {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		os.Stdout.Write(codex.MakeAllowResponse())
		os.Exit(0)
	}

	hookInput, err := codex.ParseHookInput(data)
	if err != nil || hookInput.Tool != "shell" {
		os.Stdout.Write(codex.MakeAllowResponse())
		os.Exit(0)
	}

	shellInput, err := codex.ParseShellInput(hookInput.ToolInput)
	if err != nil {
		os.Stdout.Write(codex.MakeAllowResponse())
		os.Exit(0)
	}

	// Prefer cmdline if available; fall back to cmd string split.
	var cmdArgs []string
	if len(shellInput.Cmdline) > 0 {
		cmdArgs = shellInput.Cmdline
	} else {
		cmd := shellInput.Cmd
		if hasShellMeta(cmd) {
			os.Stdout.Write(codex.MakeAllowResponse())
			os.Exit(0)
		}
		cmdArgs = splitWords(cmd)
	}

	if len(cmdArgs) == 0 {
		os.Stdout.Write(codex.MakeAllowResponse())
		os.Exit(0)
	}

	if common.IsAcapCommand(cmdArgs) {
		os.Stdout.Write(codex.MakeAllowResponse())
		os.Exit(0)
	}

	cwd, _ := os.Getwd()
	sessionID := common.MapAgentSession(hookInput.SessionID)

	req := &protocol.ToolRequest{
		Protocol:   protocol.Version,
		Command:    cmdArgs,
		WorkingDir: cwd,
		SessionID:  sessionID,
	}

	resp, err := engine.Execute(context.Background(), req)
	if err != nil {
		os.Stdout.Write(codex.MakeAllowResponse())
		os.Exit(0)
	}

	blockResp, err := codex.MakeBlockResponse(resp.Stdout)
	if err != nil {
		os.Stdout.Write(codex.MakeAllowResponse())
		os.Exit(0)
	}

	os.Stdout.Write(blockResp)
	os.Exit(2)
}

// integrateCmd implements "acap integrate <claude|codex|status> [--dry-run] [--remove]".
func integrateCmd(args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Usage: acap integrate <claude|codex|status> [--dry-run] [--remove]")
		os.Exit(1)
	}

	sub := args[0]
	rest := args[1:]

	cwd, _ := os.Getwd()
	root := project.FindRoot(cwd)

	switch sub {
	case "status":
		claudeInstalled := claude.IsInstalled(root)
		codexInstalled := codex.IsInstalled(root)
		if claudeInstalled {
			fmt.Println("claude: configured")
		} else {
			fmt.Println("claude: not configured")
		}
		if codexInstalled {
			fmt.Println("codex: configured")
		} else {
			fmt.Println("codex: not configured")
		}

	case "claude":
		fs := flag.NewFlagSet("integrate claude", flag.ExitOnError)
		dryRun := fs.Bool("dry-run", false, "print what would be done")
		remove := fs.Bool("remove", false, "remove hook")
		_ = fs.Parse(rest)

		if *remove {
			if err := claude.Remove(root); err != nil {
				fmt.Fprintf(os.Stderr, "acap: integrate claude --remove: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("claude: hook removed")
		} else {
			if err := claude.Install(root, *dryRun); err != nil {
				fmt.Fprintf(os.Stderr, "acap: integrate claude: %v\n", err)
				os.Exit(1)
			}
			if !*dryRun {
				fmt.Println("claude: hook installed")
			}
		}

	case "codex":
		fs := flag.NewFlagSet("integrate codex", flag.ExitOnError)
		dryRun := fs.Bool("dry-run", false, "print what would be done")
		remove := fs.Bool("remove", false, "remove hook")
		_ = fs.Parse(rest)

		if *remove {
			if err := codex.Remove(root); err != nil {
				fmt.Fprintf(os.Stderr, "acap: integrate codex --remove: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("codex: hook removed")
		} else {
			if err := codex.Install(root, *dryRun); err != nil {
				fmt.Fprintf(os.Stderr, "acap: integrate codex: %v\n", err)
				os.Exit(1)
			}
			if !*dryRun {
				fmt.Println("codex: hook installed")
			}
		}

	default:
		fmt.Fprintf(os.Stderr, "acap: integrate: unknown adapter %q\n", sub)
		os.Exit(1)
	}
}

// doctorCmd implements "acap doctor" — read-only diagnostic check.
func doctorCmd(args []string) {
	_ = args

	// acap binary path.
	exe, err := os.Executable()
	if err != nil {
		fmt.Printf("acap binary: OK (path unknown: %v)\n", err)
	} else {
		fmt.Printf("acap binary: OK (running as %s)\n", exe)
	}

	// Project root.
	cwd, _ := os.Getwd()
	root := project.FindRoot(cwd)
	fmt.Printf("project root: %s\n", root)

	// .acap writable.
	acapDir := project.AcapDir(root)
	if err := checkWritable(acapDir); err != nil {
		fmt.Printf(".acap writable: FAIL (%v)\n", err)
	} else {
		fmt.Printf(".acap writable: OK\n")
	}

	// SQLite store.
	st, err := store.Open(root)
	if err != nil {
		fmt.Printf("SQLite store: FAIL (%v)\n", err)
	} else {
		st.Close()
		fmt.Printf("SQLite store: OK\n")
	}

	// Claude integration.
	if claude.IsInstalled(root) {
		fmt.Printf("claude integration: configured\n")
	} else {
		fmt.Printf("claude integration: not configured\n")
	}

	// Codex integration.
	if codex.IsInstalled(root) {
		fmt.Printf("codex integration: configured\n")
	} else {
		fmt.Printf("codex integration: not configured\n")
	}
}

func checkWritable(dir string) error {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	tmp := dir + "/.acap_write_test"
	if err := os.WriteFile(tmp, []byte("ok"), 0o644); err != nil {
		return err
	}
	return os.Remove(tmp)
}

func reducerTypeName(r reduce.Reducer) string {
	switch rv := r.(type) {
	case *reduce.LsReducer:
		return "ls"
	case *reduce.FindReducer:
		return "find"
	case *reduce.GrepReducer:
		return "grep"
	case *reduce.CatReducer:
		return "cat"
	case *reduce.HeadTailReducer:
		return "headtail"
	case *reduce.TreeReducer:
		return "tree"
	case *reduce.DuReducer:
		return "du"
	case *reduce.WcReducer:
		return "wc"
	case *reduce.GitStatusReducer:
		return "git-status"
	case *reduce.GitDiffReducer:
		return "git-diff"
	case *reduce.GitShowReducer:
		return "git-show"
	case *reduce.GitLogReducer:
		return "git-log"
	case *reduce.GitBranchReducer:
		return "git-branch"
	case *reduce.GoBuildReducer:
		return "go-build"
	case *reduce.GoTestReducer:
		return "go-test"
	case *reduce.GccReducer:
		if rv.Tool != "" {
			return rv.Tool
		}
		return "gcc"
	case *reduce.CargoBuildReducer:
		return "cargo-build"
	case *reduce.CargoTestReducer:
		return "cargo-test"
	case *reduce.MakeReducer:
		return "make"
	case *reduce.NinjaReducer:
		return "ninja"
	case *reduce.GenericReducer:
		return "generic"
	default:
		return "unknown"
	}
}
