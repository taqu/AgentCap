package agentbench

import (
	"encoding/base64"
	"encoding/json"
	"strings"
	"testing"
)

func TestParseClaudeStream(t *testing.T) {
	lines := []string{
		`{"type":"system","subtype":"init","claude_code_version":"9.9.9"}`,
		`{"type":"assistant","message":{"content":[{"type":"tool_use","id":"a","name":"Bash","input":{"command":"go test ./..."}}]}}`,
		`{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"a","content":"@acap 1a2b3c go test FAIL\nfoo","is_error":true}]}}`,
		`{"type":"assistant","message":{"content":[{"type":"tool_use","id":"b","name":"Bash","input":{"command":"acap raw 1a2b3c"}}]}}`,
		`{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"b","content":[{"type":"text","text":"raw raw raw"}]}]}}`,
		`{"type":"assistant","message":{"content":[{"type":"tool_use","id":"c","name":"Read","input":{"file_path":"x"}}]}}`,
		`{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"c","content":"12345"}]}}`,
		`{"type":"assistant","message":{"content":[{"type":"tool_use","id":"d","name":"Bash","input":{"command":"go test ./..."}}]}}`,
		`{"type":"user","message":{"content":[{"type":"tool_result","tool_use_id":"d","content":"@acap 4d5e6f unchanged from 1a2b3c\nexit=1\n"}]}}`,
		`{"type":"result","usage":{"input_tokens":7,"cache_read_input_tokens":100,"cache_creation_input_tokens":20,"output_tokens":30}}`,
	}
	r := &AgentRunResult{Stdout: []byte(strings.Join(lines, "\n"))}
	parseClaudeStream(r)
	finishTrace(r)
	if r.AgentVersion != "9.9.9" || r.HostInput != 7 || r.HostCacheRead != 100 || r.HostOutput != 30 {
		t.Fatalf("metadata: %+v", r)
	}
	if r.CommandCount != 3 || len(r.Trace) != 4 {
		t.Fatalf("commands=%d trace=%d", r.CommandCount, len(r.Trace))
	}
	w := buildWorkflow(r.Trace, 0)
	if w.RawCalls != 1 || w.ImmediateRawCalls != 1 || w.RawVisibleBytes != 11 {
		t.Fatalf("raw metrics: %+v", w)
	}
	if w.AcapResults != 2 || w.UnchangedResults != 1 || w.FullResults != 1 || w.RepeatedCommands != 1 {
		t.Fatalf("result kinds: %+v", w)
	}
	if w.NativeToolCalls != 1 || w.NativeToolBytes != 5 {
		t.Fatalf("native: %+v", w)
	}
	if w.ShellVisibleBytes != w.FamilyBytes["test"]+w.FamilyBytes["acap-raw"] {
		t.Fatalf("family bytes do not partition shell bytes: %+v", w.FamilyBytes)
	}
}

func TestParseAntigravityStreamDecodesWrapper(t *testing.T) {
	payload, _ := json.Marshal(map[string]any{"request": map[string]any{"shell_command": "git diff"}})
	wrapped := "acap antigravity-exec '" + base64.RawURLEncoding.EncodeToString(payload) + "'"
	ev := func(cmd, out string) string {
		b, _ := json.Marshal(map[string]any{"event": "step_update", "step_update": map[string]any{
			"state": "DONE", "step_type": "tool", "tool_name": "run_command",
			"tool_info": map[string]any{"parameters": map[string]any{"CommandLine": cmd}, "output": out}}})
		return string(b)
	}
	lines := []string{ev(wrapped, "@acap abcdef git diff\r\n1 file\r\n"),
		`{"event":"result","result":{"usage":{"input_tokens":10,"output_tokens":2,"thinking_tokens":3}}}`}
	r := &AgentRunResult{Stdout: []byte(strings.Join(lines, "\n"))}
	parseAntigravityStream(r)
	if len(r.Trace) != 1 || r.Trace[0].Command != "git diff" || r.Trace[0].Family != "git" || r.Trace[0].ResultID != "abcdef" {
		t.Fatalf("trace: %+v", r.Trace)
	}
	if r.Trace[0].Bytes != int64(len("@acap abcdef git diff\n1 file\n")) || r.HostOutput != 5 {
		t.Fatalf("bytes/tokens: %+v %d", r.Trace[0], r.HostOutput)
	}
}

func TestClassifyFamily(t *testing.T) {
	cases := map[string]string{
		"grep -rn foo .":             "search",
		"git grep foo":               "search",
		"git diff":                   "git",
		"go test ./...":              "test",
		"go build ./...":             "build",
		"make test":                  "test",
		"make":                       "build",
		"g++ -c a.cpp":               "compiler",
		"find . -name x":             "listing",
		"cat a.go":                   "file-read",
		"acap show 1a2b3c --errors":  "acap-show",
		"cd sub && go test ./...":    "test",
		"GOFLAGS=-v go test ./x/...": "test",
		"echo hi":                    "generic",
	}
	for in, want := range cases {
		if got := classifyFamily(in); got != want {
			t.Errorf("classifyFamily(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestProductionAdaptersRejectUnsupportedModes(t *testing.T) {
	for _, a := range []ModeSupporter{&ClaudeAdapter{}, &AntigravityAdapter{}} {
		if a.SupportsMode(ModeStateless) || a.SupportsMode(ModeStateful) || !a.SupportsMode(ModeDisabled) || !a.SupportsMode(ModeIntegrated) {
			t.Fatalf("%T mode support", a)
		}
	}
}

func TestHostQuotaErrorsAreDetected(t *testing.T) {
	agy := &AgentRunResult{Stdout: []byte(`{"event":"result","result":{"status":"ERROR","error":"Individual quota reached.","usage":{}}}`)}
	parseAntigravityStream(agy)
	cl := &AgentRunResult{Stdout: []byte(`{"type":"result","is_error":true,"result":"API Error: 529 overloaded","usage":{}}`)}
	parseClaudeStream(cl)
	ok := &AgentRunResult{Stdout: []byte(`{"type":"result","is_error":true,"result":"I could not fix the test","usage":{}}`)}
	parseClaudeStream(ok)
	if agy.HostError == "" || cl.HostError == "" || ok.HostError != "" {
		t.Fatalf("host errors: agy=%q claude=%q task=%q", agy.HostError, cl.HostError, ok.HostError)
	}
}
