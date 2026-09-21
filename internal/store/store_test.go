package store

import (
	"bytes"
	"os"
	"testing"
	"time"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	root := t.TempDir()
	s, err := Open(root)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
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

	// Verify raw objects are readable.
	if got, err := os.ReadFile(e.StdoutPath()); err != nil || !bytes.Equal(got, stdout) {
		t.Errorf("stdout object: err=%v content=%q", err, got)
	}
	if got, err := os.ReadFile(e.StderrPath()); err != nil || !bytes.Equal(got, stderr) {
		t.Errorf("stderr object: err=%v content=%q", err, got)
	}

	// Verify capsule.
	cap, err := e.Capsule()
	if err != nil || cap != capsule {
		t.Errorf("capsule: err=%v got=%q", err, cap)
	}
}

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

func TestOpen_Prefix(t *testing.T) {
	s := newTestStore(t)
	e, err := s.Save(sampleMeta(), []byte("out"), []byte{}, "cap")
	if err != nil {
		t.Fatal(err)
	}
	// Use the first 4 chars as prefix.
	prefix := e.ID[:4]
	got, err := s.Open(prefix)
	if err != nil {
		t.Fatalf("Open prefix: %v", err)
	}
	if got.ID != e.ID {
		t.Errorf("id mismatch: got %s want %s", got.ID, e.ID)
	}
}

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
		t.Errorf("stdout mismatch: %v", err)
	}
	if data, err := os.ReadFile(got.StderrPath()); err != nil || !bytes.Equal(data, stderr) {
		t.Errorf("stderr mismatch: %v", err)
	}
}

func TestList(t *testing.T) {
	s := newTestStore(t)
	for i := 0; i < 3; i++ {
		if _, err := s.Save(sampleMeta(), []byte("out"), []byte{}, "cap"); err != nil {
			t.Fatal(err)
		}
	}
	entries, err := s.List()
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(entries))
	}
}

func TestDelete(t *testing.T) {
	s := newTestStore(t)
	e, err := s.Save(sampleMeta(), []byte("del"), []byte{}, "cap")
	if err != nil {
		t.Fatal(err)
	}
	if err := s.Delete(e.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Open(e.ID); err == nil {
		t.Error("expected error opening deleted entry")
	}
}

func TestCleanup(t *testing.T) {
	s := newTestStore(t)
	// Save one entry with a backdated CreatedAt.
	m := sampleMeta()
	m.CreatedAt = time.Now().UTC().Add(-10 * 24 * time.Hour)
	if _, err := s.Save(m, []byte("old"), []byte{}, "cap"); err != nil {
		t.Fatal(err)
	}
	// Clean entries older than 7 days.
	n, err := s.Cleanup(7 * 24 * time.Hour)
	if err != nil {
		t.Fatal(err)
	}
	if n != 1 {
		t.Errorf("expected 1 removed, got %d", n)
	}
}

func TestProjectIsolation(t *testing.T) {
	rootA := t.TempDir()
	rootB := t.TempDir()

	storeA, err := Open(rootA)
	if err != nil {
		t.Fatal(err)
	}
	defer storeA.Close()

	storeB, err := Open(rootB)
	if err != nil {
		t.Fatal(err)
	}
	defer storeB.Close()

	e, err := storeA.Save(sampleMeta(), []byte("a"), []byte{}, "cap-a")
	if err != nil {
		t.Fatal(err)
	}

	// storeB must not see storeA's result.
	if _, err := storeB.Open(e.ID); err == nil {
		t.Error("storeB should not find storeA result")
	}
}

func TestStats(t *testing.T) {
	s := newTestStore(t)
	if err := s.RecordRun(1000, 200, 100, "full"); err != nil {
		t.Fatal(err)
	}
	if err := s.RecordRun(2000, 400, 0, "unchanged"); err != nil {
		t.Fatal(err)
	}
	st, err := s.LoadStats()
	if err != nil {
		t.Fatal(err)
	}
	if st.Commands != 2 {
		t.Errorf("commands: got %d want 2", st.Commands)
	}
	if st.RawBytes != 3000 {
		t.Errorf("raw_bytes: got %d want 3000", st.RawBytes)
	}
	if st.UnchangedCount != 1 {
		t.Errorf("unchanged_count: got %d want 1", st.UnchangedCount)
	}
}

func TestSchemaVersion(t *testing.T) {
	root := t.TempDir()
	s, err := Open(root)
	if err != nil {
		t.Fatal(err)
	}

	// Manually bump version beyond supported.
	if _, err := s.db.Exec(`UPDATE schema_info SET version=999`); err != nil {
		t.Fatal(err)
	}
	s.Close()

	// Reopening should fail with schema version error.
	_, err = Open(root)
	if err == nil {
		t.Error("expected error for unsupported schema version")
	}
}

func BenchmarkSave(b *testing.B) {
	root := b.TempDir()
	s, err := Open(root)
	if err != nil {
		b.Fatal(err)
	}
	defer s.Close()

	stdout := bytes.Repeat([]byte("bench line\n"), 1000)
	m := sampleMeta()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := s.Save(m, stdout, []byte{}, "@acap bench\n"); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkOpen(b *testing.B) {
	root := b.TempDir()
	s, err := Open(root)
	if err != nil {
		b.Fatal(err)
	}
	defer s.Close()

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
