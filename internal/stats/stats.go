// Package stats persists reduction statistics to ~/.cache/acap/stats.json.
package stats

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Stats holds cumulative reduction statistics.
type Stats struct {
	Commands     int64 `json:"commands"`
	RawBytes     int64 `json:"raw_bytes"`
	RetBytes     int64 `json:"ret_bytes"`
	ShowCalls    int64 `json:"show_calls,omitempty"`
	ShowRetBytes int64 `json:"show_ret_bytes,omitempty"`
	RawCalls     int64 `json:"raw_calls,omitempty"`
	RawRetBytes  int64 `json:"raw_ret_bytes,omitempty"`
	// Phase 3 fields
	UnchangedCount    int64 `json:"unchanged_count,omitempty"`
	DeltaCount        int64 `json:"delta_count,omitempty"`
	FullFallbackCount int64 `json:"full_fallback_count,omitempty"`
	StatelessBytes    int64 `json:"stateless_bytes,omitempty"`
	StatefulBytes     int64 `json:"stateful_bytes,omitempty"`
}

// RecordRun records a run with optional session-aware stats.
// presentation is "full", "unchanged", "delta", or "" for non-session runs.
// stateless is what Phase 1/2 would have returned; stateful is what was actually returned.
func RecordRun(raw, stateless, stateful int, presentation string) error {
	path, err := statsPath()
	if err != nil {
		return err
	}
	s, err := loadFrom(path)
	if err != nil {
		s = &Stats{}
	}
	s.Commands++
	s.RawBytes += int64(raw)
	s.RetBytes += int64(stateful)
	s.StatelessBytes += int64(stateless)
	s.StatefulBytes += int64(stateful)
	switch presentation {
	case "unchanged":
		s.UnchangedCount++
	case "delta":
		s.DeltaCount++
	case "full":
		s.FullFallbackCount++
	}
	return saveTo(path, s)
}

// Record is kept for backwards compat; calls RecordRun with empty presentation.
func Record(raw, ret int) error {
	return RecordRun(raw, ret, ret, "")
}

// RecordShow records a show command call and the bytes returned.
func RecordShow(retBytes int) error {
	path, err := statsPath()
	if err != nil {
		return err
	}

	s, err := loadFrom(path)
	if err != nil {
		s = &Stats{}
	}

	s.ShowCalls++
	s.ShowRetBytes += int64(retBytes)

	return saveTo(path, s)
}

// RecordRaw records a raw command call and the bytes returned.
func RecordRaw(retBytes int) error {
	path, err := statsPath()
	if err != nil {
		return err
	}

	s, err := loadFrom(path)
	if err != nil {
		s = &Stats{}
	}

	s.RawCalls++
	s.RawRetBytes += int64(retBytes)

	return saveTo(path, s)
}

// Load returns accumulated statistics from disk. Returns empty stats if the
// file does not exist.
func Load() (*Stats, error) {
	path, err := statsPath()
	if err != nil {
		return nil, err
	}
	return loadFrom(path)
}

// statsPath returns the path to the stats file, creating the directory if needed.
func statsPath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("stats: cannot determine cache dir: %w", err)
	}
	dir := filepath.Join(cacheDir, "acap")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("stats: cannot create cache dir: %w", err)
	}
	return filepath.Join(dir, "stats.json"), nil
}

func loadFrom(path string) (*Stats, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Stats{}, nil
		}
		return nil, fmt.Errorf("stats: read %s: %w", path, err)
	}
	var s Stats
	if err := json.Unmarshal(data, &s); err != nil {
		// Corrupt file — return empty rather than error.
		return &Stats{}, nil
	}
	return &s, nil
}

func saveTo(path string, s *Stats) error {
	data, err := json.Marshal(s)
	if err != nil {
		return fmt.Errorf("stats: marshal: %w", err)
	}
	// Write atomically via temp file.
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return fmt.Errorf("stats: write tmp: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("stats: rename: %w", err)
	}
	return nil
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
