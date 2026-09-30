package agentbench

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/taqu/agentcap/internal/workload"
)

// TokenEstimateDivisor is the single, consistent estimation method used for
// every configuration: estimated tokens = ceil(bytes / 4). It is labelled
// "_est" everywhere and never mixed with host-measured token counts.
const TokenEstimateDivisor = 4

func estTokens(b int64) int64 { return (b + TokenEstimateDivisor - 1) / TokenEstimateDivisor }

var (
	capsuleHeader = regexp.MustCompile(`(?m)^(?:.*\n)?@acap ([0-9a-f]{6,})(?: (delta|unchanged) from)?`)
	acapRetrieval = regexp.MustCompile(`(?:^|[\s;&|/])acap\s+(show|raw)\s+([0-9a-f]{4,})`)
)

// classifyFamily maps a shell command to a coarse command family (§32).
func classifyFamily(command string) string {
	if m := acapRetrieval.FindStringSubmatch(command); m != nil {
		return "acap-" + m[1]
	}
	fields := strings.Fields(strings.TrimSpace(command))
	for len(fields) > 0 && (strings.Contains(fields[0], "=") || fields[0] == "cd" || fields[0] == "&&") {
		// skip leading env assignments and "cd dir &&"
		if fields[0] == "cd" && len(fields) > 2 {
			fields = fields[2:]
			continue
		}
		fields = fields[1:]
	}
	if len(fields) == 0 {
		return "generic"
	}
	name := filepath.Base(fields[0])
	sub := ""
	if len(fields) > 1 {
		sub = fields[1]
	}
	switch name {
	case "grep", "rg", "ag", "ack", "git-grep":
		return "search"
	case "find", "ls", "tree", "fd":
		return "listing"
	case "git":
		if sub == "grep" {
			return "search"
		}
		return "git"
	case "go":
		switch sub {
		case "test":
			return "test"
		case "build", "vet", "run", "install":
			return "build"
		}
		return "generic"
	case "make", "cmake", "ninja", "cargo":
		if strings.Contains(command, "test") || strings.Contains(command, "check") {
			return "test"
		}
		return "build"
	case "gcc", "g++", "cc", "c++", "clang", "clang++", "rustc":
		return "compiler"
	case "cat", "head", "tail", "sed", "awk", "nl", "wc", "less":
		return "file-read"
	}
	if strings.HasPrefix(name, "./") || strings.Contains(name, "test") {
		return "test"
	}
	return "generic"
}

// buildWorkflow computes all trace-derived Phase 7 metrics. Visible bytes are
// counted once, from the agent-visible tool result text.
func buildWorkflow(trace []workload.TraceEntry, instructionBytes int64) *workload.WorkflowMetrics {
	w := &workload.WorkflowMetrics{FamilyBytes: map[string]int64{}, FamilyCount: map[string]int{}, InstructionBytes: instructionBytes}
	seen := map[string]int{}
	lastResultID := ""
	for i := range trace {
		e := &trace[i]
		e.Index = i + 1
		if e.Tool != "shell" {
			w.NativeToolCalls++
			w.NativeToolBytes += e.Bytes
			continue
		}
		w.ShellCommands++
		w.ShellVisibleBytes += e.Bytes
		w.FamilyBytes[e.Family] += e.Bytes
		w.FamilyCount[e.Family]++
		key := strings.Join(strings.Fields(e.Command), " ")
		if seen[key] > 0 {
			w.RepeatedCommands++
		}
		seen[key]++
		switch e.Family {
		case "acap-show":
			w.ShowCalls++
			w.ShowVisibleBytes += e.Bytes
		case "acap-raw":
			w.RawCalls++
			w.RawVisibleBytes += e.Bytes
			if m := acapRetrieval.FindStringSubmatch(e.Command); m != nil && lastResultID != "" && strings.HasPrefix(lastResultID, m[2]) {
				w.ImmediateRawCalls++
				w.Anomalies = append(w.Anomalies, "immediate raw fallback at trace step "+itoa(e.Index)+" for "+lastResultID)
			}
		}
		if e.ResultID != "" {
			w.AcapResults++
			switch e.Kind {
			case "delta":
				w.DeltaResults++
			case "unchanged":
				w.UnchangedResults++
			default:
				w.FullResults++
			}
		}
		if !strings.HasPrefix(e.Family, "acap-") {
			lastResultID = e.ResultID
		}
	}
	w.ShellVisibleTokens = estTokens(w.ShellVisibleBytes)
	w.TotalAcapCostBytes = w.ShellVisibleBytes + w.InstructionBytes
	w.TotalAcapCostTokens = estTokens(w.TotalAcapCostBytes)
	return w
}

// shellEntry builds one shell trace entry, extracting AgentCap capsule metadata
// from the visible output when present.
func shellEntry(command, output string, isErr bool) workload.TraceEntry {
	e := workload.TraceEntry{Tool: "shell", Family: classifyFamily(command), Command: truncate(command, 300), Bytes: int64(len(output)), IsError: isErr}
	if m := capsuleHeader.FindStringSubmatch(output); m != nil && !strings.HasPrefix(e.Family, "acap-") {
		e.ResultID = m[1]
		e.Kind = "full"
		if m[2] != "" {
			e.Kind = m[2]
		}
	}
	return e
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func itoa(i int) string {
	b, _ := json.Marshal(i)
	return string(b)
}

// applyAdapterMetrics folds <root>/.acap/adapter-metrics.jsonl into w.
func applyAdapterMetrics(root string, w *workload.WorkflowMetrics) {
	f, err := os.Open(filepath.Join(root, ".acap", "adapter-metrics.jsonl"))
	if err != nil {
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		var r struct {
			Version      string `json:"adapter_version"`
			Reason       string `json:"reason"`
			Intercepted  int    `json:"commands_intercepted"`
			Bypassed     int    `json:"commands_bypassed"`
			Failures     int    `json:"adapter_failures"`
			AdapterNs    int64  `json:"adapter_latency_ns"`
			ProcessingNs int64  `json:"processing_latency_ns"`
		}
		if json.Unmarshal(scanner.Bytes(), &r) != nil {
			continue
		}
		w.AdapterVer = r.Version
		w.Bypassed += r.Bypassed
		w.AdapterFailures += r.Failures
		if r.Bypassed > 0 {
			if w.BypassReasons == nil {
				w.BypassReasons = map[string]int{}
			}
			w.BypassReasons[r.Reason]++
		}
		// Only the wrapper records an interception; it carries both latencies.
		if r.Intercepted > 0 {
			w.Intercepted += r.Intercepted
			w.AdapterLatencyNS += r.AdapterNs
			w.CoreLatencyNS += r.ProcessingNs
		}
	}
}
