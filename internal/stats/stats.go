// Package stats provides the Stats type and formatting utilities for acap.
package stats

import "fmt"

// Stats holds cumulative reduction statistics.
type Stats struct {
	Commands          int64 `json:"commands"`
	RawBytes          int64 `json:"raw_bytes"`
	RetBytes          int64 `json:"ret_bytes"`
	ShowCalls         int64 `json:"show_calls,omitempty"`
	ShowRetBytes      int64 `json:"show_ret_bytes,omitempty"`
	RawCalls          int64 `json:"raw_calls,omitempty"`
	RawRetBytes       int64 `json:"raw_ret_bytes,omitempty"`
	UnchangedCount    int64 `json:"unchanged_count,omitempty"`
	DeltaCount        int64 `json:"delta_count,omitempty"`
	FullFallbackCount int64 `json:"full_fallback_count,omitempty"`
	StatelessBytes    int64 `json:"stateless_bytes,omitempty"`
	StatefulBytes     int64 `json:"stateful_bytes,omitempty"`
}

// FormatBytes formats a byte count as a human-readable string (KB/MB/GB).
func FormatBytes(n int64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)
	switch {
	case n >= GB:
		return fmt.Sprintf("%.1fGB", float64(n)/float64(GB))
	case n >= MB:
		return fmt.Sprintf("%.1fMB", float64(n)/float64(MB))
	case n >= KB:
		return fmt.Sprintf("%.1fKB", float64(n)/float64(KB))
	default:
		return fmt.Sprintf("%dB", n)
	}
}
