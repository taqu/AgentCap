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
	Commands int64 `json:"commands"`
	RawBytes int64 `json:"raw_bytes"`
	RetBytes int64 `json:"ret_bytes"`
}

// Record loads existing stats, adds the new raw/ret bytes, and saves back.
func Record(raw, ret int) error {
	path, err := statsPath()
	if err != nil {
		return err
	}

	s, err := loadFrom(path)
	if err != nil {
		// If we can't load (file not found, corrupt), start fresh.
		s = &Stats{}
	}

	s.Commands++
	s.RawBytes += int64(raw)
	s.RetBytes += int64(ret)

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
