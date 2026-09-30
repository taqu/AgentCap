package common

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"github.com/taqu/agentcap/internal/project"
)

// AdapterMetric is one best-effort adapter record. It never contains command
// text, external session IDs or environment values.
type AdapterMetric struct {
	Agent             string
	Adapter           string
	Version           string
	Session           string // external session ID; stored only as its mapping
	ToolUseID         string
	SubagentID        string
	Reason            string
	Failed            bool
	AdapterLatency    time.Duration
	ProcessingLatency time.Duration
}

// RecordMetric appends one bounded line to <root>/.acap/adapter-metrics.jsonl.
// Reasons "intercepted" and "processing-failure" count as interceptions;
// everything else is a bypass. Failures never affect execution.
func RecordMetric(cwd string, m AdapterMetric) {
	if cwd == "" || !filepath.IsAbs(cwd) {
		return
	}
	root := project.FindRoot(cwd)
	dir := filepath.Join(root, ".acap")
	if os.MkdirAll(dir, 0755) != nil {
		return
	}
	record := struct {
		Agent        string `json:"agent"`
		Adapter      string `json:"adapter"`
		Version      string `json:"adapter_version"`
		Session      string `json:"session_mapping"`
		ToolUseID    string `json:"tool_use_id,omitempty"`
		SubagentID   string `json:"subagent_id,omitempty"`
		Reason       string `json:"reason"`
		Intercepted  int    `json:"commands_intercepted"`
		Bypassed     int    `json:"commands_bypassed"`
		Failures     int    `json:"adapter_failures"`
		AdapterNs    int64  `json:"adapter_latency_ns"`
		ProcessingNs int64  `json:"processing_latency_ns"`
	}{Agent: m.Agent, Adapter: m.Adapter, Version: m.Version, Session: MapSession(m.Agent, m.Session, root), ToolUseID: m.ToolUseID, SubagentID: m.SubagentID, Reason: m.Reason, AdapterNs: m.AdapterLatency.Nanoseconds(), ProcessingNs: m.ProcessingLatency.Nanoseconds()}
	if m.Reason == "intercepted" || m.Reason == "processing-failure" {
		record.Intercepted = 1
	} else {
		record.Bypassed = 1
	}
	if m.Failed {
		record.Failures = 1
	}
	data, _ := json.Marshal(record)
	f, err := os.OpenFile(filepath.Join(dir, "adapter-metrics.jsonl"), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return
	}
	defer f.Close()
	f.Write(append(data, '\n'))
}
