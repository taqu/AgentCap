package stats

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRecordAndLoad(t *testing.T) {
	// Use a temp dir for stats.
	dir := t.TempDir()
	// Override the stats path for this test.
	origLoad := loadFrom
	origSave := saveTo
	_ = origLoad
	_ = origSave

	path := filepath.Join(dir, "stats.json")

	// Record some entries.
	if err := recordTo(path, 1000, 200); err != nil {
		t.Fatalf("recordTo: %v", err)
	}
	if err := recordTo(path, 2000, 400); err != nil {
		t.Fatalf("recordTo: %v", err)
	}

	s, err := loadFrom(path)
	if err != nil {
		t.Fatalf("loadFrom: %v", err)
	}
	if s.Commands != 2 {
		t.Errorf("commands: got %d, want 2", s.Commands)
	}
	if s.RawBytes != 3000 {
		t.Errorf("raw_bytes: got %d, want 3000", s.RawBytes)
	}
	if s.RetBytes != 600 {
		t.Errorf("ret_bytes: got %d, want 600", s.RetBytes)
	}
}

func TestLoadMissing(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nonexistent.json")
	s, err := loadFrom(path)
	if err != nil {
		t.Fatalf("loadFrom missing: %v", err)
	}
	if s.Commands != 0 {
		t.Errorf("expected 0 commands for missing file")
	}
}

func TestLoadCorrupt(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "stats.json")
	if err := os.WriteFile(path, []byte("not valid json"), 0o644); err != nil {
		t.Fatal(err)
	}
	s, err := loadFrom(path)
	if err != nil {
		t.Fatalf("loadFrom corrupt: %v", err)
	}
	if s.Commands != 0 {
		t.Errorf("expected empty stats for corrupt file")
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		n    int64
		want string
	}{
		{0, "0B"},
		{512, "512B"},
		{1024, "1.0KB"},
		{1536, "1.5KB"},
		{1048576, "1.0MB"},
		{1073741824, "1.0GB"},
	}
	for _, tt := range tests {
		got := FormatBytes(tt.n)
		if got != tt.want {
			t.Errorf("FormatBytes(%d) = %q, want %q", tt.n, got, tt.want)
		}
	}
}

// recordTo is a test helper that calls the internal save logic.
func recordTo(path string, raw, ret int) error {
	s, err := loadFrom(path)
	if err != nil {
		s = &Stats{}
	}
	s.Commands++
	s.RawBytes += int64(raw)
	s.RetBytes += int64(ret)
	return saveTo(path, s)
}
