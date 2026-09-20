package store

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// newTestStore creates a Store backed by a temp directory.
func newTestStore(t *testing.T) *Store {
	t.Helper()
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "results"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(dir, "tmp"), 0o700); err != nil {
		t.Fatal(err)
	}
	return &Store{dir: dir}
}

func sampleMeta() Meta {
	return Meta{
		Command:   []string{"rg", "Workspace", "."},
		ExitCode:  0,
		StartedAt: time.Now().UTC().Truncate(time.Second),
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		Reducer:   "rg",
	}
}

// TestSave verifies all four files are created with correct content.
func TestSave(t *testing.T) {
	s := newTestStore(t)
	stdout := []byte("hello stdout")
	stderr := []byte("hello stderr")
	capsule := "@acap abc123 rg matches=5 files=2\n"

	e, err := s.Save(sampleMeta(), stdout, stderr, capsule)
	if err != nil {
		t.Fatalf("Save: %v", err)
	}
	if len(e.ID) != 6 {
		t.Errorf("expected 6-char id, got %q", e.ID)
	}

	if got, err := os.ReadFile(e.StdoutPath()); err != nil || !bytes.Equal(got, stdout) {
		t.Errorf("stdout file: err=%v content=%q", err, got)
	}
	if got, err := os.ReadFile(e.StderrPath()); err != nil || !bytes.Equal(got, stderr) {
		t.Errorf("stderr file: err=%v content=%q", err, got)
	}
	if got, err := os.ReadFile(e.CapsulePath()); err != nil || string(got) != capsule {
		t.Errorf("capsule file: err=%v content=%q", err, got)
	}
	if got, err := os.ReadFile(filepath.Join(e.Dir, "meta.json")); err != nil || len(got) == 0 {
		t.Errorf("meta.json: err=%v", err)
	}
}

// TestOpen_Exact verifies exact-ID lookup.
func TestOpen_Exact(t *testing.T) {
	s := newTestStore(t)
	e, err := s.Save(sampleMeta(), []byte("out"), []byte("err"), "capsule")
	if err != nil {
		t.Fatal(err)
	}

	got, err := s.Open(e.ID)
	if err != nil {
		t.Fatalf("Open exact: %v", err)
	}
	if got.ID != e.ID {
		t.Errorf("id mismatch: got %s want %s", got.ID, e.ID)
	}
}

// TestOpen_Prefix verifies prefix lookup and ambiguity detection.
func TestOpen_Prefix(t *testing.T) {
	s := newTestStore(t)

	// Inject two known directories directly so we can control their names.
	id1 := "aabbcc"
	id2 := "aabbdd"
	for _, id := range []string{id1, id2} {
		dir := filepath.Join(s.dir, "results", id)
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		m := Meta{ID: id, Command: []string{"echo"}, CreatedAt: time.Now().UTC()}
		data, _ := jsonMarshal(m)
		if err := os.WriteFile(filepath.Join(dir, "meta.json"), data, 0o600); err != nil {
			t.Fatal(err)
		}
		// create required files
		for _, f := range []string{"stdout", "stderr", "capsule"} {
			_ = os.WriteFile(filepath.Join(dir, f), []byte{}, 0o600)
		}
	}

	// Unique prefix "aabbcc" should resolve to id1.
	e, err := s.Open("aabbcc")
	if err != nil {
		t.Fatalf("Open prefix: %v", err)
	}
	if e.ID != id1 {
		t.Errorf("expected %s, got %s", id1, e.ID)
	}

	// Ambiguous prefix "aabb" should fail.
	_, err = s.Open("aabb")
	if err == nil {
		t.Error("expected ambiguity error for prefix 'aabb'")
	}
	if !strings.Contains(err.Error(), "ambiguous") {
		t.Errorf("expected 'ambiguous' in error, got: %v", err)
	}
}

// TestOpen_Missing verifies a clean error for unknown IDs.
func TestOpen_Missing(t *testing.T) {
	s := newTestStore(t)
	_, err := s.Open("zzzzzz")
	if err == nil {
		t.Error("expected error for unknown id")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Errorf("expected 'not found' in error, got: %v", err)
	}
}

// TestSave_NonZeroExit verifies failed commands are stored correctly.
func TestSave_NonZeroExit(t *testing.T) {
	s := newTestStore(t)
	m := sampleMeta()
	m.ExitCode = 1
	m.Truncated = true

	e, err := s.Save(m, []byte(""), []byte("error: not found"), "")
	if err != nil {
		t.Fatal(err)
	}
	got, err := s.Open(e.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Meta.ExitCode != 1 {
		t.Errorf("exit_code: got %d want 1", got.Meta.ExitCode)
	}
	if !got.Meta.Truncated {
		t.Error("expected truncated=true")
	}
}

// TestCorrupt verifies that an entry with a missing meta.json is skipped in List
// and returns an error from Open.
func TestCorrupt(t *testing.T) {
	s := newTestStore(t)

	// Create an entry dir without meta.json.
	corruptDir := filepath.Join(s.dir, "results", "corrupt")
	if err := os.MkdirAll(corruptDir, 0o700); err != nil {
		t.Fatal(err)
	}

	// List should not error out and should skip the corrupt entry.
	entries, err := s.List()
	if err != nil {
		t.Fatalf("List with corrupt entry: %v", err)
	}
	for _, e := range entries {
		if e.ID == "corrupt" {
			t.Error("corrupt entry should be skipped in List")
		}
	}

	// Open should return an error.
	_, err = s.Open("corrupt")
	if err == nil {
		t.Error("expected error opening corrupt entry")
	}
}

// TestCleanup verifies old results are removed, new ones kept.
func TestCleanup(t *testing.T) {
	s := newTestStore(t)

	// Save an "old" entry by manually adjusting its meta.
	e, err := s.Save(sampleMeta(), []byte("old"), []byte{}, "old capsule")
	if err != nil {
		t.Fatal(err)
	}
	// Patch meta.json with an old timestamp.
	m := e.Meta
	m.CreatedAt = time.Now().Add(-10 * 24 * time.Hour)
	data, _ := jsonMarshal(m)
	_ = os.WriteFile(filepath.Join(e.Dir, "meta.json"), data, 0o600)

	// Save a new entry.
	eNew, err := s.Save(sampleMeta(), []byte("new"), []byte{}, "new capsule")
	if err != nil {
		t.Fatal(err)
	}

	removed, err := s.Cleanup(7 * 24 * time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if removed != 1 {
		t.Errorf("expected 1 removed, got %d", removed)
	}

	// Old entry should be gone.
	if _, err := os.Stat(e.Dir); !os.IsNotExist(err) {
		t.Error("old entry should have been removed")
	}
	// New entry should remain.
	if _, err := os.Stat(eNew.Dir); err != nil {
		t.Errorf("new entry should still exist: %v", err)
	}
}

// TestAtomicCreation verifies that the tmp directory is not visible as a valid
// result during creation.
func TestAtomicCreation(t *testing.T) {
	s := newTestStore(t)

	// List results before any save — should be empty.
	entries, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries before save, got %d", len(entries))
	}

	// Manually place a dir in tmp/ — it should not appear in List or Open.
	tmpDir := filepath.Join(s.dir, "tmp", "phantom")
	if err := os.MkdirAll(tmpDir, 0o700); err != nil {
		t.Fatal(err)
	}
	m := Meta{ID: "phantom", Command: []string{"echo"}, CreatedAt: time.Now().UTC()}
	data, _ := jsonMarshal(m)
	_ = os.WriteFile(filepath.Join(tmpDir, "meta.json"), data, 0o600)

	// List should still return 0 results (tmp/ is separate from results/).
	entries, err = s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Errorf("tmp entry should not appear in List, got %d", len(entries))
	}

	// Open should not find it either.
	if _, err := s.Open("phantom"); err == nil {
		t.Error("tmp entry should not be openable via Open")
	}
}

// TestLargeOutput verifies that a 20MB stdout is saved fully (not capped at 10MB).
func TestLargeOutput(t *testing.T) {
	s := newTestStore(t)

	const size = 20 * 1024 * 1024
	large := bytes.Repeat([]byte("x"), size)

	e, err := s.Save(sampleMeta(), large, []byte{}, "capsule")
	if err != nil {
		t.Fatalf("Save large: %v", err)
	}

	fi, err := os.Stat(e.StdoutPath())
	if err != nil {
		t.Fatal(err)
	}
	if fi.Size() != int64(size) {
		t.Errorf("stdout file size: got %d, want %d", fi.Size(), size)
	}
}

// TestRoundTrip tests Save → Open → read capsule/stdout/stderr.
func TestRoundTrip(t *testing.T) {
	s := newTestStore(t)

	stdout := []byte("round trip stdout\n")
	stderr := []byte("round trip stderr\n")
	capsule := "@acap xxxxxx rg matches=1 files=1\n"

	entry, err := s.Save(sampleMeta(), stdout, stderr, capsule)
	if err != nil {
		t.Fatal(err)
	}

	got, err := s.Open(entry.ID)
	if err != nil {
		t.Fatal(err)
	}

	cap2, err := got.Capsule()
	if err != nil {
		t.Fatal(err)
	}
	if cap2 != capsule {
		t.Errorf("capsule mismatch: got %q want %q", cap2, capsule)
	}

	if data, err := os.ReadFile(got.StdoutPath()); err != nil || !bytes.Equal(data, stdout) {
		t.Errorf("stdout mismatch: %v %q", err, data)
	}
	if data, err := os.ReadFile(got.StderrPath()); err != nil || !bytes.Equal(data, stderr) {
		t.Errorf("stderr mismatch: %v %q", err, data)
	}
}

// BenchmarkSave benchmarks the Save operation.
func BenchmarkSave(b *testing.B) {
	dir := b.TempDir()
	_ = os.MkdirAll(filepath.Join(dir, "results"), 0o700)
	_ = os.MkdirAll(filepath.Join(dir, "tmp"), 0o700)
	s := &Store{dir: dir}

	stdout := bytes.Repeat([]byte("bench line\n"), 1000)
	stderr := []byte{}
	capsule := "@acap bench\n"
	m := sampleMeta()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := s.Save(m, stdout, stderr, capsule); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkOpen benchmarks the Open operation.
func BenchmarkOpen(b *testing.B) {
	dir := b.TempDir()
	_ = os.MkdirAll(filepath.Join(dir, "results"), 0o700)
	_ = os.MkdirAll(filepath.Join(dir, "tmp"), 0o700)
	s := &Store{dir: dir}

	e, err := s.Save(sampleMeta(), []byte("bench"), []byte{}, "bench")
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := s.Open(e.ID); err != nil {
			b.Fatal(err)
		}
	}
}

// jsonMarshal is a test helper to marshal without importing encoding/json directly.
func jsonMarshal(v interface{}) ([]byte, error) {
	// Use the same encoding/json used in production code.
	return marshalJSON(v)
}
