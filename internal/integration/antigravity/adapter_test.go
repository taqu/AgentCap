package antigravity

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	osexec "os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/taqu/agentcap/internal/integration/conformance"
	"github.com/taqu/agentcap/internal/integration/protocol"
	"github.com/taqu/agentcap/internal/project"
	"github.com/taqu/agentcap/internal/store"
)

// hookData mirrors the payload captured from a real Antigravity CLI 1.2.14.
func hookData(command, cwd, conversation string, waitMs any) []byte {
	args := map[string]any{"CommandLine": command, "Cwd": cwd, "toolAction": "Running shell command", "toolSummary": "Run", "FutureArg": "keep"}
	if waitMs != nil {
		args["WaitMsBeforeAsync"] = waitMs
	}
	data, _ := json.Marshal(map[string]any{
		"conversationId": conversation, "stepIdx": 7, "modelName": "gemini-test",
		"workspacePaths": []string{"/elsewhere/a", "/elsewhere/b"}, "transcriptPath": "/t.jsonl",
		"artifactDirectoryPath": "/artifacts", "futureTopLevel": true,
		"toolCall": map[string]any{"name": "run_command", "args": args},
	})
	return data
}

type hookResponse struct {
	Decision  string            `json:"decision"`
	Overwrite map[string]string `json:"overwrite"`
	Other     map[string]any    `json:"-"`
}

func decode(t *testing.T, data []byte) hookResponse {
	t.Helper()
	var r hookResponse
	if err := json.Unmarshal(data, &r); err != nil {
		t.Fatalf("%s: %v", data, err)
	}
	if r.Decision != "ask" {
		t.Fatalf("adapter must never grant or deny: %s", data)
	}
	return r
}

func preparedPayload(t *testing.T, command, cwd, conversation string) string {
	t.Helper()
	t.Setenv(EnvPermissions, permissionsDeclared)
	executable, _ := os.Executable()
	response, reason, err := Prepare(hookData(command, cwd, conversation, 5000), executable)
	if err != nil || reason != "intercepted" {
		t.Fatalf("Prepare: %s %v", reason, err)
	}
	r := decode(t, response)
	if len(r.Overwrite) != 1 {
		t.Fatalf("only CommandLine may be overwritten: %s", response)
	}
	_, payload, ok := strings.Cut(r.Overwrite["CommandLine"], " antigravity-exec ")
	if !ok {
		t.Fatal(r.Overwrite["CommandLine"])
	}
	return strings.Trim(payload, "'")
}

func invoke(t *testing.T) conformance.Invoke {
	return func(ctx context.Context, command, cwd, session string) *protocol.ToolResponse {
		payload := preparedPayload(t, command, cwd, session)
		var stdout, stderr bytes.Buffer
		code := ExecutePayload(ctx, payload, &stdout, &stderr)
		resp := &protocol.ToolResponse{ExitCode: code, Stdout: stdout.String(), Stderr: stderr.String()}
		fields := strings.Fields(resp.Stdout)
		if len(fields) > 1 && fields[0] == "@acap" {
			resp.ResultID = fields[1]
			st, err := store.Open(project.FindRoot(cwd))
			if err != nil {
				t.Error(err)
				return resp
			}
			defer st.Close()
			entry, err := st.Open(resp.ResultID)
			if err != nil {
				t.Error(err)
				return resp
			}
			resp.Presentation = entry.Meta.Presentation
		}
		return resp
	}
}

func TestAntigravityConformance(t *testing.T) {
	if _, err := shellPath(); err != nil {
		t.Skip(err)
	}
	conformance.Run(t, invoke(t))
}

func TestParseHookInput(t *testing.T) {
	h, err := ParseHookInput(hookData("go test ./...", "/w", "conv", 5000))
	if err != nil {
		t.Fatal(err)
	}
	a, err := ParseRunCommandArgs(h.ToolCall.Args)
	if err != nil {
		t.Fatal(err)
	}
	if h.ConversationID != "conv" || h.ModelName != "gemini-test" || *h.StepIdx != 7 || len(h.WorkspacePaths) != 2 ||
		a.CommandLine != "go test ./..." || a.Cwd != "/w" || *a.WaitMsBeforeAsync != 5000 {
		t.Fatalf("%+v %+v", h, a)
	}
	for name, data := range map[string]string{
		"malformed":     `{`,
		"no tool":       `{"conversationId":"c"}`,
		"no tool name":  `{"toolCall":{"args":{}}}`,
		"array payload": `[]`,
	} {
		if _, err := ParseHookInput([]byte(data)); err == nil {
			t.Errorf("%s: expected error", name)
		}
	}
	for _, args := range []string{``, `{}`, `{"CommandLine":"  "}`, `{"CommandLine":7}`} {
		if _, err := ParseRunCommandArgs(json.RawMessage(args)); err == nil {
			t.Errorf("%q: expected error", args)
		}
	}
}

func TestPassThroughNeverRewritesOrGrants(t *testing.T) {
	executable, _ := os.Executable()
	cwd := t.TempDir()
	check := func(t *testing.T, data []byte, wantReason string) {
		t.Helper()
		response, reason, _ := Prepare(data, executable)
		if reason != wantReason || string(response) != `{"decision":"ask"}` {
			t.Fatalf("got %s %s, want pass-through %s", response, reason, wantReason)
		}
	}
	t.Run("permission mode not declared", func(t *testing.T) {
		for _, v := range []string{"", "1", "request-review", "UNRESTRICTED"} {
			t.Setenv(EnvPermissions, v)
			for _, command := range []string{"git status", "rm temporary-file", "echo foo | grep foo", "cd subdir && go test ./...", "FOO=bar ./script.sh"} {
				check(t, hookData(command, cwd, "c", 5000), "permission-mode")
			}
		}
	})
	t.Setenv(EnvPermissions, permissionsDeclared)
	t.Run("other tools", func(t *testing.T) {
		for _, tool := range []string{"view_file", "write_to_file", "list_dir", "browser_click_element"} {
			data, _ := json.Marshal(map[string]any{"toolCall": map[string]any{"name": tool, "args": map[string]any{"CommandLine": "x"}}})
			check(t, data, "unsupported-tool")
		}
	})
	t.Run("invalid", func(t *testing.T) {
		check(t, []byte(`{`), "invalid-input")
		check(t, []byte(`{"toolCall":{"name":"run_command","args":{"Cwd":"/"}}}`), "invalid-input")
	})
	t.Run("recursion", func(t *testing.T) {
		for _, command := range []string{"acap show abc", "acap raw abc --lines 1:3", "acap stats", "acap clean", "/usr/local/bin/acap show x", "cd sub && acap show x"} {
			check(t, hookData(command, cwd, "c", 5000), "recursion")
		}
		t.Setenv(protocol.EnvDepth, "1")
		check(t, hookData("echo ok", cwd, "c", 5000), "recursion")
	})
	for _, env := range []string{protocol.EnvBypass, protocol.EnvBenchmarkMode} {
		t.Run(env, func(t *testing.T) {
			value := "1"
			if env == protocol.EnvBenchmarkMode {
				value = "disabled"
			}
			t.Setenv(env, value)
			check(t, hookData("echo ok", cwd, "c", 5000), "explicit-bypass")
		})
	}
	t.Run("async", func(t *testing.T) {
		for _, wait := range []any{0, 500, 1999} {
			check(t, hookData("go test ./...", cwd, "c", wait), "async")
		}
	})
	t.Run("long-running", func(t *testing.T) {
		for _, command := range []string{"tail -f log.txt", "npm run dev", "go test ./... --watch", "vim main.go", "python -m http.server && serve ."} {
			check(t, hookData(command, cwd, "c", 5000), "long-running-or-interactive")
		}
	})
	t.Run("cwd", func(t *testing.T) {
		check(t, hookData("echo ok", "", "c", 5000), "missing-cwd")
		check(t, hookData("echo ok", "relative/dir", "c", 5000), "missing-cwd")
	})
	t.Run("binary missing", func(t *testing.T) {
		response, reason, err := Prepare(hookData("echo ok", cwd, "c", 5000), filepath.Join(cwd, "missing-acap"))
		if reason != "binary-unavailable" || err == nil || string(response) != `{"decision":"ask"}` {
			t.Fatal(string(response), reason, err)
		}
	})
	t.Run("intercepts ordinary commands when declared", func(t *testing.T) {
		for _, wait := range []any{nil, 2000, 10000} {
			response, reason, err := Prepare(hookData("go test ./...", cwd, "c", wait), executable)
			if reason != "intercepted" || err != nil || !strings.Contains(string(response), "antigravity-exec") {
				t.Fatal(string(response), reason, err)
			}
		}
	})
}

func TestShellSemanticsAndTransport(t *testing.T) {
	if _, err := shellPath(); err != nil {
		t.Skip(err)
	}
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".git"), 0755)
	for command, want := range map[string]string{
		`printf '%s\n' "a b"`:                            "a b",
		`FOO="hello world" sh -c 'printf "%s\n" "$FOO"'`: "hello world",
		`printf 'a\nb\n' | grep b`:                       "b",
		`false || echo recovered`:                        "recovered",
		`mkdir -p nested && cd nested && pwd`:            filepath.Join(dir, "nested"),
		`printf 'x\n' > file.txt`:                        "",
		`for f in a b; do echo "item-$f"; done`:          "item-b",
		`printf '%s' '\"; touch injected; #'`:            `\"; touch injected; #`,
	} {
		payload := preparedPayload(t, command, dir, "")
		var p execution
		decoded, _ := base64.RawURLEncoding.DecodeString(payload)
		json.Unmarshal(decoded, &p)
		if *p.Request.ShellCommand != command || p.Request.WorkingDir != dir {
			t.Fatal("shell source or cwd changed")
		}
		var stdout, stderr bytes.Buffer
		if code := ExecutePayload(context.Background(), payload, &stdout, &stderr); code != 0 || !strings.Contains(stdout.String(), want) {
			t.Fatalf("%s: %d %q %q", command, code, stdout.String(), stderr.String())
		}
	}
	if _, err := os.Stat(filepath.Join(dir, "injected")); !os.IsNotExist(err) {
		t.Fatal("wrapper injected command syntax")
	}
	if data, _ := os.ReadFile(filepath.Join(dir, "file.txt")); string(data) != "x\n" {
		t.Fatalf("redirection: %q", data)
	}
}

// The overwritten CommandLine is executed by Antigravity's shell; run it the
// same way to prove the full transport, including quoting of the binary path.
func TestOverwrittenCommandLineExecutesOnce(t *testing.T) {
	sh, err := osexec.LookPath("sh")
	if err != nil {
		t.Skip(err)
	}
	acap := filepath.Join(t.TempDir(), "acap bin", "acap")
	os.MkdirAll(filepath.Dir(acap), 0755)
	build := osexec.Command("go", "build", "-o", acap, "github.com/taqu/agentcap/cmd/acap")
	if out, err := build.CombinedOutput(); err != nil {
		t.Skip("cannot build acap:", err, string(out))
	}
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".git"), 0755)
	os.WriteFile(filepath.Join(dir, "increment-counter.sh"), []byte("#!/bin/sh\nn=$(cat counter 2>/dev/null || echo 0)\necho $((n+1)) > counter\necho counted\nexit ${1:-0}\n"), 0755)
	t.Setenv(EnvPermissions, permissionsDeclared)
	for i, tc := range []struct {
		command string
		code    int
	}{{"./increment-counter.sh", 0}, {"./increment-counter.sh 7", 7}} {
		response, reason, err := Prepare(hookData(tc.command, dir, "conv", 5000), acap)
		if reason != "intercepted" || err != nil {
			t.Fatal(reason, err)
		}
		cmd := osexec.Command(sh, "-c", decode(t, response).Overwrite["CommandLine"])
		cmd.Dir = dir
		out, _ := cmd.CombinedOutput()
		if cmd.ProcessState.ExitCode() != tc.code || !strings.HasPrefix(string(out), "@acap ") {
			t.Fatalf("%s: exit %d output %q", tc.command, cmd.ProcessState.ExitCode(), out)
		}
		if data, _ := os.ReadFile(filepath.Join(dir, "counter")); strings.TrimSpace(string(data)) != string(rune('1'+i)) {
			t.Fatalf("counter after %d runs = %q", i+1, data)
		}
		id := strings.Fields(string(out))[1]
		raw := osexec.Command(acap, "raw", id)
		raw.Dir = dir
		if rawOut, err := raw.Output(); err != nil || !strings.Contains(string(rawOut), "counted") {
			t.Fatalf("acap raw %s: %v %q", id, err, rawOut)
		}
	}
}

func TestGoTestFailureIsCompactAndStateful(t *testing.T) {
	if _, err := osexec.LookPath("go"); err != nil {
		t.Skip(err)
	}
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".git"), 0755)
	os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module fixture\n\ngo 1.21\n"), 0644)
	var body strings.Builder
	body.WriteString("package fixture\n\nimport \"testing\"\n\nfunc TestBroken(t *testing.T) {\n")
	for i := 0; i < 300; i++ {
		body.WriteString("\tt.Log(\"noise line padding padding padding\")\n")
	}
	body.WriteString("\tt.Fatal(\"FIX_ME\")\n}\n\nfunc TestOK(t *testing.T) {}\n")
	os.WriteFile(filepath.Join(dir, "fixture_test.go"), []byte(body.String()), 0644)

	run := func() *protocol.ToolResponse {
		return invoke(t)(context.Background(), "go test ./...", dir, "conversation-A")
	}
	first := run()
	if first.ExitCode == 0 || first.ResultID == "" || first.Presentation != "full" {
		t.Fatalf("first: %+v", first)
	}
	st, _ := store.Open(dir)
	entry, _ := st.Open(first.ResultID)
	raw, _ := os.ReadFile(entry.StdoutPath())
	st.Close()
	if len(first.Stdout) >= len(raw) || !strings.Contains(string(raw), "FIX_ME") || strings.Count(string(raw), "noise line") != 300 {
		t.Fatalf("expected compact presentation of complete raw output: presented=%d raw=%d", len(first.Stdout), len(raw))
	}
	if entry.Meta.Integration == nil || entry.Meta.Integration.Agent != "antigravity" || entry.Meta.Integration.Model != "gemini-test" {
		t.Fatalf("attribution: %+v", entry.Meta.Integration)
	}
	second := run()
	if second.ExitCode == 0 || second.Presentation == "full" {
		t.Fatalf("same conversation should compare with baseline: %+v", second)
	}
	other := invoke(t)(context.Background(), "go test ./...", dir, "conversation-B")
	if other.Presentation != "full" {
		t.Fatalf("conversations must be isolated: %+v", other)
	}
}

// A model change inside one conversation must not reset the session.
func TestModelNameIsNotSessionIdentity(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".git"), 0755)
	t.Setenv(EnvPermissions, permissionsDeclared)
	executable, _ := os.Executable()
	var presentations []string
	for _, model := range []string{"model-a", "model-b"} {
		var h map[string]any
		json.Unmarshal(hookData("printf same", dir, "conv", 5000), &h)
		h["modelName"] = model
		data, _ := json.Marshal(h)
		response, _, _ := Prepare(data, executable)
		_, payload, _ := strings.Cut(decode(t, response).Overwrite["CommandLine"], " antigravity-exec ")
		var stdout, stderr bytes.Buffer
		ExecutePayload(context.Background(), strings.Trim(payload, "'"), &stdout, &stderr)
		st, _ := store.Open(dir)
		entry, _ := st.Open(strings.Fields(stdout.String())[1])
		presentations = append(presentations, entry.Meta.Presentation)
		st.Close()
	}
	if presentations[1] != "unchanged" {
		t.Fatalf("model change reset session: %v", presentations)
	}
}

func TestLargeOutputNoDeadlockOrTruncation(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".git"), 0755)
	resp := invoke(t)(context.Background(), `i=0; while [ $i -lt 40000 ]; do echo "out line $i"; echo "err line $i" >&2; i=$((i+1)); done; exit 3`, dir, "large")
	if resp.ExitCode != 3 || resp.ResultID == "" {
		t.Fatalf("%d %q", resp.ExitCode, resp.ResultID)
	}
	st, _ := store.Open(dir)
	defer st.Close()
	entry, _ := st.Open(resp.ResultID)
	stdout, _ := os.ReadFile(entry.StdoutPath())
	stderr, _ := os.ReadFile(entry.StderrPath())
	if strings.Count(string(stdout), "\n") != 40000 || strings.Count(string(stderr), "\n") != 40000 {
		t.Fatalf("stdout=%d stderr=%d", len(stdout), len(stderr))
	}
	// Complex shell programs use the existing generic reducer; a classifiable
	// command shows reduction through the same adapter path.
	os.WriteFile(filepath.Join(dir, "big.log"), stdout, 0644)
	cat := invoke(t)(context.Background(), "cat big.log", dir, "large")
	if cat.ExitCode != 0 || cat.ResultID == "" || len(cat.Stdout) >= len(stdout)/10 {
		t.Fatalf("presentation not reduced: %d bytes", len(cat.Stdout))
	}
}

func TestInvalidPayloadNeverExecutes(t *testing.T) {
	for _, payload := range []string{"invalid", base64.RawURLEncoding.EncodeToString([]byte(`{}`))} {
		var stdout, stderr bytes.Buffer
		if code := ExecutePayload(context.Background(), payload, &stdout, &stderr); code != 2 || stdout.Len() != 0 || !strings.Contains(stderr.String(), "not started") {
			t.Fatalf("%d %s", code, stdout.String())
		}
	}
}

func TestMetricsRecorded(t *testing.T) {
	dir := t.TempDir()
	os.MkdirAll(filepath.Join(dir, ".git"), 0755)
	invoke(t)(context.Background(), "printf ok", dir, "secret-conversation-id")
	t.Setenv(EnvPermissions, "")
	executable, _ := os.Executable()
	Prepare(hookData("printf ok", dir, "secret-conversation-id", 5000), executable)
	data, _ := os.ReadFile(filepath.Join(dir, ".acap", "adapter-metrics.jsonl"))
	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) != 2 || strings.Contains(string(data), "secret-conversation-id") || strings.Contains(string(data), "printf") {
		t.Fatalf("%s", data)
	}
	var m map[string]any
	json.Unmarshal([]byte(lines[0]), &m)
	if m["agent"] != "antigravity" || m["adapter_version"] != AdapterVersion || m["commands_intercepted"] != 1.0 || m["tool_use_id"] != "7" {
		t.Fatal(lines[0])
	}
	json.Unmarshal([]byte(lines[1]), &m)
	if m["reason"] != "permission-mode" || m["commands_bypassed"] != 1.0 {
		t.Fatal(lines[1])
	}
}
