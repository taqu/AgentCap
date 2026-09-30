package workload

// WorkflowMetrics is the Phase 7 workflow-level record attached to coding-agent
// trial results. It is additive to schema version 4: older consumers ignore it.
//
// Metric boundary: every agent-visible byte is counted exactly once, from the
// host agent's own tool-result event stream. Store statistics (raw bytes,
// stateful bytes, processing time) are recorded separately as attribution and
// are never added to the visible totals.
type WorkflowMetrics struct {
	RunID      string `json:"run_id"`
	Batch      string `json:"batch,omitempty"`
	StartedAt  string `json:"started_at"`
	Platform   string `json:"platform"`
	Category   string `json:"category,omitempty"`
	Language   string `json:"language,omitempty"`
	AgentVer   string `json:"agent_version,omitempty"`
	Model      string `json:"model,omitempty"`
	AcapCommit string `json:"agentcap_version,omitempty"`
	AdapterVer string `json:"adapter_version,omitempty"`
	Outcome    string `json:"outcome"` // PASS, FAIL, TIMEOUT, INFRA_ERROR

	// Shell/command tool results as seen by the agent (both modes).
	ShellCommands       int   `json:"shell_commands"`
	RepeatedCommands    int   `json:"repeated_commands"`
	ShellVisibleBytes   int64 `json:"shell_visible_bytes"`
	ShellVisibleTokens  int64 `json:"shell_visible_tokens_est"`
	ShowCalls           int   `json:"show_calls"`
	ShowVisibleBytes    int64 `json:"show_visible_bytes"`
	RawCalls            int   `json:"raw_calls"`
	RawVisibleBytes     int64 `json:"raw_visible_bytes"`
	ImmediateRawCalls   int   `json:"immediate_raw_calls"`
	AcapResults         int   `json:"acap_results"`
	InstructionBytes    int64 `json:"instruction_bytes"`
	TotalAcapCostBytes  int64 `json:"total_context_cost_bytes"`
	TotalAcapCostTokens int64 `json:"total_context_cost_tokens_est"`

	// Native, non-shell tool payloads (file reads, native search...). Not
	// AgentCap-controlled; reported so shell savings are seen in proportion.
	NativeToolCalls int   `json:"native_tool_calls"`
	NativeToolBytes int64 `json:"native_tool_bytes"`

	FamilyBytes map[string]int64 `json:"family_bytes,omitempty"`
	FamilyCount map[string]int   `json:"family_count,omitempty"`

	// Adapter metrics log (.acap/adapter-metrics.jsonl), ON mode only.
	Intercepted      int            `json:"commands_intercepted"`
	Bypassed         int            `json:"commands_bypassed"`
	AdapterFailures  int            `json:"adapter_failures"`
	AdapterLatencyNS int64          `json:"adapter_latency_ns"`
	CoreLatencyNS    int64          `json:"core_latency_ns"`
	BypassReasons    map[string]int `json:"bypass_reasons,omitempty"`

	// Store presentation kinds (ON mode only).
	FullResults      int `json:"full_results"`
	DeltaResults     int `json:"delta_results"`
	UnchangedResults int `json:"unchanged_results"`

	// Host-reported model tokens (measured, not estimated). Zero means unavailable.
	HostInputTokens      int64 `json:"host_input_tokens"`
	HostCacheReadTokens  int64 `json:"host_cache_read_tokens"`
	HostCacheWriteTokens int64 `json:"host_cache_write_tokens"`
	HostOutputTokens     int64 `json:"host_output_tokens"`

	Anomalies []string     `json:"anomalies,omitempty"`
	Trace     []TraceEntry `json:"trace,omitempty"`
}

// TraceEntry is one ordered tool call in the agent workflow.
type TraceEntry struct {
	Index    int    `json:"i"`
	Tool     string `json:"tool"`
	Family   string `json:"family"`
	Command  string `json:"command,omitempty"`
	Bytes    int64  `json:"bytes"`
	IsError  bool   `json:"is_error,omitempty"`
	ResultID string `json:"result_id,omitempty"`
	Kind     string `json:"kind,omitempty"` // full, delta, unchanged, bypass
}
