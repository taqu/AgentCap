package bench

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/taqu/agentcap/internal/integration/engine"
	"github.com/taqu/agentcap/internal/integration/protocol"
	"github.com/taqu/agentcap/internal/store"
)

// helperArg makes the test binary act as a deterministic target command.
const helperArg = "-bench-helper"

func TestMain(m *testing.M) {
	if len(os.Args) > 1 && os.Args[1] == helperArg {
		os.Exit(runHelper(os.Args[2:]))
	}
	os.Exit(m.Run())
}

// runHelper executes a sequence of actions:
//
//	out <s>      write s to stdout
//	err <s>      write s to stderr
//	touch <path> append one line to path (observable side effect)
//	argv         write each remaining arg on its own line, then stop
//	exit <n>     exit with status n
//	file <path>  copy a file to stdout
func runHelper(args []string) int {
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "out":
			i++
			os.Stdout.WriteString(args[i])
		case "err":
			i++
			os.Stderr.WriteString(args[i])
		case "touch":
			i++
			f, err := os.OpenFile(args[i], os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
			if err != nil {
				return 99
			}
			f.WriteString("x\n")
			f.Close()
		case "file":
			i++
			data, err := os.ReadFile(args[i])
			if err != nil {
				return 98
			}
			os.Stdout.Write(data)
		case "argv":
			for _, a := range args[i+1:] {
				fmt.Fprintf(os.Stdout, "%s\n", a)
			}
			return 0
		case "exit":
			i++
			n, _ := strconv.Atoi(args[i])
			return n
		}
	}
	return 0
}

func TestSessionRepeatedResultAndAggregate(t *testing.T) {
	root := isolate(t)
	id, err := StartSession(root)
	if err != nil {
		t.Fatal(err)
	}
	output := strings.Repeat("a useful repeated output line\n", 100)
	command := helper(t, "out", output)
	first, err := MeasureSessionCommand(context.Background(), command, root, id)
	if err != nil {
		t.Fatal(err)
	}
	second, err := MeasureSessionCommand(context.Background(), command, root, id)
	if err != nil {
		t.Fatal(err)
	}
	if first.Presentation != "full" || second.Presentation != "unchanged" {
		t.Fatalf("presentations = %q, %q", first.Presentation, second.Presentation)
	}
	if second.StatefulVisibleBytes >= second.StatelessVisibleBytes {
		t.Errorf("second stateful/stateless = %d/%d", second.StatefulVisibleBytes, second.StatelessVisibleBytes)
	}
	if first.StatelessVisibleBytes != first.StatefulVisibleBytes {
		t.Errorf("first stateless/stateful = %d/%d", first.StatelessVisibleBytes, first.StatefulVisibleBytes)
	}
	agg, err := LoadSession(root, id)
	if err != nil {
		t.Fatal(err)
	}
	if agg.CommandCount != 2 || agg.FullCount != 1 || agg.UnchangedCount != 1 || agg.DeltaCount != 0 {
		t.Errorf("aggregate counts = %+v", agg)
	}
	if agg.RawBytes != first.RawBytes+second.RawBytes ||
		agg.StatelessBytes != first.StatelessVisibleBytes+second.StatelessVisibleBytes ||
		agg.StatefulBytes != first.StatefulVisibleBytes+second.StatefulVisibleBytes {
		t.Errorf("aggregate byte totals = %+v", agg)
	}
}

func TestSessionChangedResultUsesNormalDeltaDecision(t *testing.T) {
	root := isolate(t)
	id, err := StartSession(root)
	if err != nil {
		t.Fatal(err)
	}
	fixture := filepath.Join(root, "output.txt")
	base := strings.Repeat("same line\n", 100)
	if err := os.WriteFile(fixture, []byte(base+"old line\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	command := helper(t, "file", fixture)
	if _, err := MeasureSessionCommand(context.Background(), command, root, id); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(fixture, []byte(base+"new line\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	changed, err := MeasureSessionCommand(context.Background(), command, root, id)
	if err != nil {
		t.Fatal(err)
	}
	if changed.Presentation != "delta" && changed.Presentation != "full" {
		t.Fatalf("changed presentation = %q", changed.Presentation)
	}
	if changed.Presentation == "delta" && changed.StatefulVisibleBytes >= changed.StatelessVisibleBytes {
		t.Errorf("delta stateful/stateless = %d/%d", changed.StatefulVisibleBytes, changed.StatelessVisibleBytes)
	}
}

func TestSessionCommandIdentityAndIsolation(t *testing.T) {
	root := isolate(t)
	output := strings.Repeat("identical output\n", 80)
	one, err := StartSession(root)
	if err != nil {
		t.Fatal(err)
	}
	firstCommand := helper(t, "out", output)
	secondCommand := helper(t, "noop", "ignored", "out", output)
	if _, err := MeasureSessionCommand(context.Background(), firstCommand, root, one); err != nil {
		t.Fatal(err)
	}
	different, err := MeasureSessionCommand(context.Background(), secondCommand, root, one)
	if err != nil {
		t.Fatal(err)
	}
	if different.Presentation != "full" {
		t.Errorf("different command presentation = %q, want full", different.Presentation)
	}

	two, err := StartSession(root)
	if err != nil {
		t.Fatal(err)
	}
	isolated, err := MeasureSessionCommand(context.Background(), firstCommand, root, two)
	if err != nil {
		t.Fatal(err)
	}
	if isolated.Presentation != "full" {
		t.Errorf("new session presentation = %q, want full", isolated.Presentation)
	}
}

func TestSessionExactlyOnceFailedEmptyAndRecoverable(t *testing.T) {
	root := isolate(t)
	id, err := StartSession(root)
	if err != nil {
		t.Fatal(err)
	}
	counter := filepath.Join(root, "session-counter.txt")
	m, err := MeasureSessionCommand(context.Background(),
		helper(t, "touch", counter, "out", "failure output\n", "exit", "3"), root, id)
	if err != nil {
		t.Fatal(err)
	}
	if countLines(t, counter) != 1 || m.ExitCode != 3 {
		t.Errorf("side effects=%d exit=%d", countLines(t, counter), m.ExitCode)
	}
	empty, err := MeasureSessionCommand(context.Background(), helper(t), root, id)
	if err != nil {
		t.Fatal(err)
	}
	if empty.RawBytes != 0 {
		t.Errorf("empty RawBytes = %d", empty.RawBytes)
	}
	st, err := store.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	entry, err := st.Open(m.ResultID)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(entry.StdoutPath())
	if err != nil || string(raw) != "failure output\n" {
		t.Errorf("stored raw = %q, err=%v", raw, err)
	}
}

func helper(t *testing.T, actions ...string) []string {
	t.Helper()
	exe, err := os.Executable()
	if err != nil {
		t.Fatal(err)
	}
	return append([]string{exe, helperArg}, actions...)
}

// isolate points result storage at a fresh temp root.
func isolate(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	t.Setenv("ACAP_ROOT", root)
	t.Setenv("ACAP_SESSION_ID", "")
	t.Setenv("LOCALAPPDATA", root)
	t.Setenv("XDG_CACHE_HOME", root)
	return root
}

func countLines(t *testing.T, path string) int {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0
		}
		t.Fatal(err)
	}
	return strings.Count(string(data), "\n")
}

func TestMeasureExecutesExactlyOnce(t *testing.T) {
	root := isolate(t)
	counter := filepath.Join(root, "counter.txt")

	m, err := MeasureCommand(context.Background(),
		helper(t, "touch", counter, "out", "hello\n"), root)
	if err != nil {
		t.Fatal(err)
	}
	if n := countLines(t, counter); n != 1 {
		t.Fatalf("target executed %d times, want exactly 1", n)
	}
	if m.RawStdoutBytes != int64(len("hello\n")) {
		t.Errorf("RawStdoutBytes = %d", m.RawStdoutBytes)
	}
}

func TestMeasureInvokesPipelineOnce(t *testing.T) {
	isolate(t)
	calls := 0
	counting := func(ctx context.Context, req *protocol.ToolRequest) (*engine.Outcome, error) {
		calls++
		return engine.Run(ctx, req)
	}
	if _, err := measure(context.Background(), counting, helper(t, "out", "x"), t.TempDir(), ""); err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("pipeline invoked %d times, want 1", calls)
	}
}

func TestMeasureIsStateless(t *testing.T) {
	isolate(t)
	var got *protocol.ToolRequest
	spy := func(ctx context.Context, req *protocol.ToolRequest) (*engine.Outcome, error) {
		got = req
		return engine.Run(ctx, req)
	}
	t.Setenv("ACAP_SESSION_ID", "some-session")
	m, err := measure(context.Background(), spy, helper(t, "out", "x"), t.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	if got.SessionID != "" {
		t.Errorf("SessionID = %q, want empty (stateless)", got.SessionID)
	}
	if m.Presentation != "full" {
		t.Errorf("Presentation = %q, want full", m.Presentation)
	}
}

func TestRawByteCounting(t *testing.T) {
	isolate(t)
	cases := []struct {
		name           string
		actions        []string
		stdout, stderr int64
	}{
		{"stdout only", []string{"out", "abcdef"}, 6, 0},
		{"stderr only", []string{"err", "oops!"}, 0, 5},
		{"both", []string{"out", "12345678", "err", "abc"}, 8, 3},
		{"empty", nil, 0, 0},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			m, err := MeasureCommand(context.Background(), helper(t, tc.actions...), t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			if m.RawStdoutBytes != tc.stdout || m.RawStderrBytes != tc.stderr {
				t.Errorf("stdout/stderr = %d/%d, want %d/%d",
					m.RawStdoutBytes, m.RawStderrBytes, tc.stdout, tc.stderr)
			}
			if m.RawBytes != tc.stdout+tc.stderr {
				t.Errorf("RawBytes = %d, want %d", m.RawBytes, tc.stdout+tc.stderr)
			}
		})
	}
}

// AgentVisibleBytes must equal the rendered presentation: the stored capsule
// plus the injected "@acap <id>" header, and nothing else.
func TestAgentVisibleBytesIsRenderedPresentation(t *testing.T) {
	isolate(t)
	m, err := MeasureCommand(context.Background(),
		helper(t, "out", "line one\nline two\n", "err", "warn\n"), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if m.ResultID == "" {
		t.Fatal("result not stored")
	}
	st, err := store.New()
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	e, err := st.Open(m.ResultID)
	if err != nil {
		t.Fatal(err)
	}
	capsule, _ := e.Capsule()
	// Generic small output: capsule does not start with "@acap ", so the
	// pipeline prepends a header line.
	want := "@acap " + m.ResultID + "\n" + capsule
	if m.AgentVisibleBytes != int64(len(want)) {
		t.Errorf("AgentVisibleBytes = %d, want %d (%q)", m.AgentVisibleBytes, len(want), want)
	}
	if !strings.Contains(capsule, "warn") {
		t.Errorf("stderr not folded into presentation: %q", capsule)
	}
}

func TestReportNotCountedAsVisible(t *testing.T) {
	isolate(t)
	m, err := MeasureCommand(context.Background(), helper(t, "out", "hello\n"), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	before := *m
	var buf bytes.Buffer
	if err := WriteText(&buf, m); err != nil {
		t.Fatal(err)
	}
	if err := WriteJSON(&buf, m); err != nil {
		t.Fatal(err)
	}
	if m.AgentVisibleBytes != before.AgentVisibleBytes {
		t.Error("rendering the report changed AgentVisibleBytes")
	}
	if int64(buf.Len()) == m.AgentVisibleBytes {
		t.Error("report size unexpectedly equals visible bytes")
	}
}

func TestReductionRatio(t *testing.T) {
	m := &Measurement{RawBytes: 1000, AgentVisibleBytes: 250}
	r, ok := m.ReductionRatio()
	if !ok || r != 0.75 {
		t.Errorf("ratio = %v, %v; want 0.75, true", r, ok)
	}

	zero := &Measurement{RawBytes: 0, AgentVisibleBytes: 20}
	if _, ok := zero.ReductionRatio(); ok {
		t.Error("ratio defined for zero raw bytes")
	}
}

func TestZeroOutputCommand(t *testing.T) {
	isolate(t)
	m, err := MeasureCommand(context.Background(), helper(t), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if m.RawBytes != 0 {
		t.Fatalf("RawBytes = %d, want 0", m.RawBytes)
	}
	var text bytes.Buffer
	if err := WriteText(&text, m); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(text.String(), "reduction:  n/a") {
		t.Errorf("zero-output report missing n/a reduction:\n%s", text.String())
	}
	var js bytes.Buffer
	if err := WriteJSON(&js, m); err != nil {
		t.Fatal(err)
	}
	var obj map[string]any
	if err := json.Unmarshal(js.Bytes(), &obj); err != nil {
		t.Fatal(err)
	}
	if v, present := obj["reduction_ratio"]; !present || v != nil {
		t.Errorf("reduction_ratio = %v (present=%v), want null", v, present)
	}
}

func TestSessionReportZeroDenominators(t *testing.T) {
	m := &SessionMeasurement{SessionID: "empty"}
	var text bytes.Buffer
	if err := WriteSessionText(&text, m); err != nil {
		t.Fatal(err)
	}
	if strings.Count(text.String(), "n/a") != 3 {
		t.Errorf("zero-denominator report:\n%s", text.String())
	}
	var js bytes.Buffer
	if err := WriteSessionJSON(&js, m); err != nil {
		t.Fatal(err)
	}
	var obj map[string]any
	if err := json.Unmarshal(js.Bytes(), &obj); err != nil {
		t.Fatal(err)
	}
	for _, key := range []string{"stateless_vs_raw", "stateful_vs_raw", "stateful_vs_stateless"} {
		if obj[key] != nil {
			t.Errorf("%s = %v, want null", key, obj[key])
		}
	}
}

func TestSuccessfulCommandReport(t *testing.T) {
	isolate(t)
	m, err := MeasureCommand(context.Background(), helper(t, "out", "ok\n"), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if m.ExitCode != 0 {
		t.Errorf("ExitCode = %d", m.ExitCode)
	}
	if m.ExecutionDuration <= 0 {
		t.Error("ExecutionDuration not measured")
	}
	if m.ProcessingDuration <= 0 || m.ProcessingDuration < m.ReduceDuration {
		t.Errorf("ProcessingDuration = %s, ReduceDuration = %s", m.ProcessingDuration, m.ReduceDuration)
	}
	var buf bytes.Buffer
	if err := WriteText(&buf, m); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"exit:       0", "total:      3 B", "id:         " + m.ResultID} {
		if !strings.Contains(buf.String(), want) {
			t.Errorf("report missing %q:\n%s", want, buf.String())
		}
	}
}

func TestFailedCommandStillMeasured(t *testing.T) {
	isolate(t)
	m, err := MeasureCommand(context.Background(),
		helper(t, "out", "FAIL\n", "err", "boom\n", "exit", "3"), t.TempDir())
	if err != nil {
		t.Fatalf("non-zero exit treated as benchmark failure: %v", err)
	}
	if m.ExitCode != 3 {
		t.Errorf("ExitCode = %d, want 3", m.ExitCode)
	}
	if m.RawBytes != 10 {
		t.Errorf("RawBytes = %d, want 10", m.RawBytes)
	}
	if m.ResultID == "" {
		t.Error("failed command result not stored")
	}
}

func TestCommandNotFound(t *testing.T) {
	isolate(t)
	_, err := MeasureCommand(context.Background(),
		[]string{"acap-definitely-not-a-command-xyz"}, t.TempDir())
	var ee *ExecError
	if !errors.As(err, &ee) {
		t.Fatalf("err = %v, want *ExecError", err)
	}
	if ee.ExitCode != 127 {
		t.Errorf("ExitCode = %d, want 127", ee.ExitCode)
	}
}

func TestArgvPreserved(t *testing.T) {
	isolate(t)
	args := []string{"hello world", `%s\n`, "*", "a'b\"c", "$HOME", ""}
	m, err := MeasureCommand(context.Background(),
		helper(t, append([]string{"argv"}, args...)...), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	st, err := store.New()
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	e, err := st.Open(m.ResultID)
	if err != nil {
		t.Fatal(err)
	}
	raw, err := os.ReadFile(e.StdoutPath())
	if err != nil {
		t.Fatal(err)
	}
	want := strings.Join(args, "\n") + "\n"
	if string(raw) != want {
		t.Errorf("argv not preserved:\n got %q\nwant %q", raw, want)
	}
}

func TestStoredResultRetrievable(t *testing.T) {
	isolate(t)
	m, err := MeasureCommand(context.Background(),
		helper(t, "out", "stored-out\n", "err", "stored-err\n"), t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	st, err := store.New()
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()
	e, err := st.Open(m.ResultID)
	if err != nil {
		t.Fatal(err)
	}
	if e.Meta.StdoutBytes != m.RawStdoutBytes || e.Meta.StderrBytes != m.RawStderrBytes {
		t.Errorf("stored meta sizes %d/%d != measured %d/%d",
			e.Meta.StdoutBytes, e.Meta.StderrBytes, m.RawStdoutBytes, m.RawStderrBytes)
	}
	out, _ := os.ReadFile(e.StdoutPath())
	errb, _ := os.ReadFile(e.StderrPath())
	if string(out) != "stored-out\n" || string(errb) != "stored-err\n" {
		t.Errorf("stored raw = %q / %q", out, errb)
	}
}
