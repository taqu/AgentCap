package query

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

// makeReader returns a reader with n lines of the form "line N: content".
func makeReader(lines []string) *strings.Reader {
	return strings.NewReader(strings.Join(lines, "\n") + "\n")
}

// sampleLines generates n distinct lines.
func sampleLines(n int) []string {
	out := make([]string, n)
	for i := range out {
		out[i] = fmt.Sprintf("line %d: content goes here", i+1)
	}
	return out
}

// TestLines_Basic verifies basic range extraction.
func TestLines_Basic(t *testing.T) {
	src := []string{"alpha", "beta", "gamma", "delta", "epsilon"}
	r := makeReader(src)
	got, err := Lines(r, 2, 4)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"beta", "gamma", "delta"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("[%d] got %q want %q", i, got[i], want[i])
		}
	}
}

// TestLines_OutOfBounds verifies that requesting beyond EOF returns what exists.
func TestLines_OutOfBounds(t *testing.T) {
	src := []string{"a", "b", "c"}
	r := makeReader(src)
	got, err := Lines(r, 2, 100)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{"b", "c"}
	if len(got) != len(want) {
		t.Fatalf("got %v, want %v", got, want)
	}
}

// TestLines_Empty verifies empty input returns empty slice.
func TestLines_Empty(t *testing.T) {
	got, err := Lines(strings.NewReader(""), 1, 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty, got %v", got)
	}
}

// TestLines_OneBased verifies 1-based numbering (from=1 returns first line).
func TestLines_OneBased(t *testing.T) {
	src := []string{"first", "second"}
	r := makeReader(src)
	got, err := Lines(r, 1, 1)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != "first" {
		t.Errorf("expected [first], got %v", got)
	}
}

// TestLines_InvalidRange verifies that from > to returns nil.
func TestLines_InvalidRange(t *testing.T) {
	src := []string{"a", "b", "c"}
	r := makeReader(src)
	got, err := Lines(r, 5, 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty for invalid range, got %v", got)
	}
}

// TestMatch_Basic verifies basic case-insensitive match.
func TestMatch_Basic(t *testing.T) {
	src := []string{
		"Hello World",
		"goodbye world",
		"Hello again",
		"nothing here",
	}
	r := makeReader(src)
	res, err := Match(r, "hello", 10)
	if err != nil {
		t.Fatal(err)
	}
	if res.Total != 2 {
		t.Errorf("total: got %d want 2", res.Total)
	}
	if len(res.Hits) != 2 {
		t.Errorf("hits: got %d want 2", len(res.Hits))
	}
	if res.Omitted != 0 {
		t.Errorf("omitted: got %d want 0", res.Omitted)
	}
}

// TestMatch_CaseInsensitive verifies uppercase/lowercase both match.
func TestMatch_CaseInsensitive(t *testing.T) {
	src := []string{"UPPER", "lower", "Mixed"}
	r := makeReader(src)
	res, err := Match(r, "UPPER", 10)
	if err != nil {
		t.Fatal(err)
	}
	if res.Total != 1 || res.Hits[0].Line != "UPPER" {
		t.Errorf("unexpected result: %+v", res)
	}
}

// TestMatch_MaxLimit verifies maxHits truncation and Omitted count.
func TestMatch_MaxLimit(t *testing.T) {
	// 10 matching lines, max 3.
	lines := make([]string, 10)
	for i := range lines {
		lines[i] = fmt.Sprintf("match line %d", i)
	}
	r := makeReader(lines)
	res, err := Match(r, "match", 3)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Hits) != 3 {
		t.Errorf("hits: got %d want 3", len(res.Hits))
	}
	if res.Total != 10 {
		t.Errorf("total: got %d want 10", res.Total)
	}
	if res.Omitted != 7 {
		t.Errorf("omitted: got %d want 7", res.Omitted)
	}
}

// TestMatch_Empty verifies no matches on empty input.
func TestMatch_Empty(t *testing.T) {
	res, err := Match(strings.NewReader(""), "anything", 10)
	if err != nil {
		t.Fatal(err)
	}
	if res.Total != 0 || len(res.Hits) != 0 {
		t.Errorf("expected empty result, got %+v", res)
	}
}

// TestMatch_LineNumbers verifies that Hit.LineNum is 1-based and correct.
func TestMatch_LineNumbers(t *testing.T) {
	src := []string{"no", "yes match", "no", "yes here"}
	r := makeReader(src)
	res, err := Match(r, "yes", 10)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Hits) != 2 {
		t.Fatalf("expected 2 hits, got %d", len(res.Hits))
	}
	if res.Hits[0].LineNum != 2 {
		t.Errorf("first hit linenum: got %d want 2", res.Hits[0].LineNum)
	}
	if res.Hits[1].LineNum != 4 {
		t.Errorf("second hit linenum: got %d want 4", res.Hits[1].LineNum)
	}
}

// TestPathFilter_Basic verifies path substring filtering.
func TestPathFilter_Basic(t *testing.T) {
	src := []string{
		"src/foo.go:1:package main",
		"internal/bar.go:5:func Bar()",
		"src/baz.go:10:var x = 1",
		"other/file.go:1:package other",
	}
	r := makeReader(src)
	lines, total, err := PathFilter(r, "src/", 10)
	if err != nil {
		t.Fatal(err)
	}
	if total != 2 {
		t.Errorf("total: got %d want 2", total)
	}
	if len(lines) != 2 {
		t.Errorf("lines: got %d want 2", len(lines))
	}
}

// TestPathFilter_Limit verifies the maxLines limit is respected.
func TestPathFilter_Limit(t *testing.T) {
	src := make([]string, 20)
	for i := range src {
		src[i] = fmt.Sprintf("path/file%d.go content", i)
	}
	r := makeReader(src)
	lines, total, err := PathFilter(r, "path/", 5)
	if err != nil {
		t.Fatal(err)
	}
	if total != 20 {
		t.Errorf("total: got %d want 20", total)
	}
	if len(lines) != 5 {
		t.Errorf("lines: got %d want 5", len(lines))
	}
}

// BenchmarkLines benchmarks Lines on a 10k-line file.
func BenchmarkLines(b *testing.B) {
	data := []byte(strings.Join(sampleLines(10000), "\n") + "\n")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(data)
		if _, err := Lines(r, 100, 200); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkMatch benchmarks Match on a 10k-line file.
func BenchmarkMatch(b *testing.B) {
	data := []byte(strings.Join(sampleLines(10000), "\n") + "\n")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(data)
		if _, err := Match(r, "content", 50); err != nil {
			b.Fatal(err)
		}
	}
}

// BenchmarkPathFilter benchmarks PathFilter on a 10k-line file.
func BenchmarkPathFilter(b *testing.B) {
	src := make([]string, 10000)
	for i := range src {
		if i%3 == 0 {
			src[i] = fmt.Sprintf("src/file%d.go:%d:content", i, i)
		} else {
			src[i] = fmt.Sprintf("other/file%d.go:%d:content", i, i)
		}
	}
	data := []byte(strings.Join(src, "\n") + "\n")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		r := bytes.NewReader(data)
		if _, _, err := PathFilter(r, "src/", 50); err != nil {
			b.Fatal(err)
		}
	}
}
