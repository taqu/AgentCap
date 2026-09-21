package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/taqu/agentcap/internal/delta"
	"github.com/taqu/agentcap/internal/exec"
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
  acap session <start|info|list|history>
  acap history

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
  acap session start
  acap session info
  acap session list
  acap history

Environment:
  ACAP_DEBUG=1        Print reducer debug info to stderr.
  ACAP_SESSION_ID=xx  Enable session-aware delta compression.
`

const runUsage = `Usage: acap run <command> [args...]

Executes <command> directly (no shell interpretation), reduces the output
for LLM consumption, and writes it to stdout.  The child's exit code is
preserved.
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
	case "session":
		sessionCmd(args[1:])
	case "history":
		historyCmd(args[1:])
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

	if sess != nil && baselineEntry != nil && entry != nil {
		dr := selectDeltaReducer(reducer)
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
		if err := entry.UpdateMeta(updatedMeta); err != nil && debugMode {
			fmt.Fprintf(os.Stderr, "acap: warning: update meta: %v\n", err)
		}
	} else if entry != nil {
		updatedMeta := entry.Meta
		updatedMeta.Presentation = "full"
		_ = entry.UpdateMeta(updatedMeta)
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
	_ = stats.RecordRun(reduced.RawBytes, stateless, stateful, presentation)

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
func selectDeltaReducer(r reduce.Reducer) delta.Reducer {
	switch r.(type) {
	case *reduce.GrepReducer:
		return &delta.GrepDelta{}
	case *reduce.FindReducer:
		return &delta.FindDelta{}
	case *reduce.LsReducer:
		return &delta.LsDelta{}
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
	fs.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: acap show <id> [--meta] [--lines X:Y] [--match text] [--path path]")
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
	entry, err := s.Open(id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "acap: show: %v\n", err)
		os.Exit(1)
	}

	var written int

	switch {
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

	_ = stats.RecordShow(written)
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
	_ = stats.RecordRaw(int(n))
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
	s, err := stats.Load()
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

func reducerTypeName(r reduce.Reducer) string {
	switch r.(type) {
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
	case *reduce.GenericReducer:
		return "generic"
	default:
		return "unknown"
	}
}
