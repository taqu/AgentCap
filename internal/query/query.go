// Package query provides streaming line-oriented queries over io.Reader sources.
// All functions work on io.Reader so they never load entire files into memory.
package query

import (
	"bufio"
	"io"
	"strings"
)

// Lines returns the lines at 1-based positions [from, to] inclusive.
// Lines are returned as-is (no newline appended).
// If from > to or from < 1, an empty slice is returned.
// If the reader has fewer lines than to, lines up to EOF are returned.
func Lines(r io.Reader, from, to int) ([]string, error) {
	if from < 1 {
		from = 1
	}
	if to < from {
		return nil, nil
	}

	scanner := bufio.NewScanner(r)
	// Allow long lines (up to 1MB each).
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	var out []string
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		if lineNum < from {
			continue
		}
		if lineNum > to {
			break
		}
		out = append(out, scanner.Text())
	}
	return out, scanner.Err()
}

// Hit represents a single matching line.
type Hit struct {
	LineNum int
	Line    string
}

// MatchResult holds the output of a Match call.
type MatchResult struct {
	Hits    []Hit
	Total   int // total matching lines (may exceed len(Hits))
	Omitted int
}

// Match searches r for lines containing text (case-insensitive).
// Up to maxHits hits are returned. If more exist, Omitted is set.
func Match(r io.Reader, text string, maxHits int) (*MatchResult, error) {
	lower := strings.ToLower(text)
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	res := &MatchResult{}
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), lower) {
			res.Total++
			if len(res.Hits) < maxHits {
				res.Hits = append(res.Hits, Hit{LineNum: lineNum, Line: line})
			}
		}
	}
	res.Omitted = res.Total - len(res.Hits)
	return res, scanner.Err()
}

// PathFilter returns lines from r that contain the given path substring.
// Up to maxLines results are collected. Returns the lines, the total count of
// matching lines, and any error.
func PathFilter(r io.Reader, path string, maxLines int) ([]string, int, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)

	var out []string
	total := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, path) {
			total++
			if len(out) < maxLines {
				out = append(out, line)
			}
		}
	}
	return out, total, scanner.Err()
}
