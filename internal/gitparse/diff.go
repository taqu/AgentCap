package gitparse

import (
	"bytes"
	"strconv"
	"strings"
)

const maxDiffBytes = 50 * 1024 * 1024 // 50MB safety limit

// ParseDiff parses unified git diff output.
// Returns nil if parsing fails or input exceeds limit.
func ParseDiff(raw []byte) *GitDiff {
	if len(raw) > maxDiffBytes {
		return nil
	}
	if len(raw) == 0 {
		return &GitDiff{}
	}

	// Find all "diff --git " positions to get byte ranges for each file.
	prefix := []byte("diff --git ")
	var starts []int
	offset := 0
	for {
		idx := bytes.Index(raw[offset:], prefix)
		if idx < 0 {
			break
		}
		starts = append(starts, offset+idx)
		offset += idx + len(prefix)
	}

	if len(starts) == 0 {
		return &GitDiff{}
	}

	diff := &GitDiff{}
	for i, start := range starts {
		var end int
		if i+1 < len(starts) {
			end = starts[i+1]
		} else {
			end = len(raw)
		}
		f := parseOneFileDiff(raw[start:end], int64(start))
		if f != nil {
			f.Index = i
			diff.Files = append(diff.Files, *f)
		}
	}
	return diff
}

func parseOneFileDiff(data []byte, startOffset int64) *GitDiffFile {
	text := string(data)
	lines := strings.Split(text, "\n")

	f := &GitDiffFile{
		RawStart: startOffset,
		RawEnd:   startOffset + int64(len(data)),
	}

	// Parse "diff --git a/X b/Y"
	if len(lines) > 0 {
		parseDiffGitLine(lines[0], f)
	}

	// Track byte offset for hunk positioning.
	byteOffset := startOffset
	hunkIndex := 0

	var currentHunk *GitDiffHunk
	inHunk := false

	for li, line := range lines {
		lineLen := int64(len(line))
		if li > 0 {
			lineLen++ // account for '\n'
		}

		if li == 0 {
			byteOffset += int64(len(lines[0]))
			if len(lines) > 1 {
				byteOffset++ // '\n'
			}
			continue
		}

		// Metadata lines (before any @@)
		if !inHunk {
			if strings.HasPrefix(line, "similarity index ") {
				// "similarity index 85%"
				numStr := strings.TrimPrefix(line, "similarity index ")
				numStr = strings.TrimSuffix(numStr, "%")
				n, _ := strconv.Atoi(strings.TrimSpace(numStr))
				f.Similarity = n
			} else if strings.HasPrefix(line, "rename from ") {
				f.OldPath = strings.TrimPrefix(line, "rename from ")
				f.OldPath = strings.TrimSpace(f.OldPath)
				f.Status = "R"
			} else if strings.HasPrefix(line, "rename to ") {
				f.NewPath = strings.TrimPrefix(line, "rename to ")
				f.NewPath = strings.TrimSpace(f.NewPath)
				f.Status = "R"
			} else if strings.HasPrefix(line, "copy from ") {
				f.OldPath = strings.TrimPrefix(line, "copy from ")
				f.OldPath = strings.TrimSpace(f.OldPath)
				f.Status = "C"
			} else if strings.HasPrefix(line, "copy to ") {
				f.NewPath = strings.TrimPrefix(line, "copy to ")
				f.NewPath = strings.TrimSpace(f.NewPath)
				f.Status = "C"
			} else if strings.HasPrefix(line, "old mode ") {
				f.OldMode = strings.TrimPrefix(line, "old mode ")
				f.OldMode = strings.TrimSpace(f.OldMode)
			} else if strings.HasPrefix(line, "new mode ") {
				f.NewMode = strings.TrimPrefix(line, "new mode ")
				f.NewMode = strings.TrimSpace(f.NewMode)
			} else if strings.HasPrefix(line, "Binary files") || strings.HasPrefix(line, "GIT binary patch") {
				f.Binary = true
			} else if strings.HasPrefix(line, "--- ") {
				// Use --- line to refine path if not set
				path := strings.TrimPrefix(line, "--- ")
				path = strings.TrimSpace(path)
				if path != "/dev/null" && path != "a/"+f.OldPath {
					// strip a/ prefix
					if strings.HasPrefix(path, "a/") {
						path = path[2:]
					}
					if f.OldPath == "" {
						f.OldPath = path
					}
				}
				if path == "/dev/null" {
					// new file
					if f.Status == "" {
						f.Status = "A"
					}
				}
			} else if strings.HasPrefix(line, "+++ ") {
				path := strings.TrimPrefix(line, "+++ ")
				path = strings.TrimSpace(path)
				if path != "/dev/null" {
					if strings.HasPrefix(path, "b/") {
						path = path[2:]
					}
					if f.NewPath == "" {
						f.NewPath = path
					}
				} else {
					// deleted file
					if f.Status == "" {
						f.Status = "D"
					}
				}
			} else if strings.HasPrefix(line, "@@ ") {
				inHunk = true
				// Fall through to hunk handling below.
				goto handleHunk
			}
			byteOffset += lineLen
			continue
		}

	handleHunk:
		if strings.HasPrefix(line, "@@ ") {
			// End previous hunk.
			if currentHunk != nil {
				currentHunk.RawEnd = byteOffset
				f.Hunks = append(f.Hunks, *currentHunk)
			}
			oldStart, oldLines, newStart, newLines, header := parseHunkHeader(line)
			currentHunk = &GitDiffHunk{
				Index:    hunkIndex,
				OldStart: oldStart,
				OldLines: oldLines,
				NewStart: newStart,
				NewLines: newLines,
				Header:   header,
				RawStart: byteOffset,
			}
			hunkIndex++
		} else if currentHunk != nil {
			if len(line) > 0 {
				switch line[0] {
				case '+':
					if !strings.HasPrefix(line, "+++") {
						f.Additions++
					}
				case '-':
					if !strings.HasPrefix(line, "---") {
						f.Deletions++
					}
				}
			}
		}
		byteOffset += lineLen
	}

	// Close last hunk.
	if currentHunk != nil {
		currentHunk.RawEnd = startOffset + int64(len(data))
		f.Hunks = append(f.Hunks, *currentHunk)
	}

	// Determine status if not set.
	if f.Status == "" {
		if f.OldPath == "" || f.OldPath == "/dev/null" {
			f.Status = "A"
		} else if f.NewPath == "" || f.NewPath == "/dev/null" {
			f.Status = "D"
		} else {
			f.Status = "M"
		}
	}

	// Ensure paths are set.
	if f.NewPath == "" && f.OldPath != "" {
		f.NewPath = f.OldPath
	}
	if f.OldPath == "" && f.NewPath != "" {
		f.OldPath = f.NewPath
	}

	return f
}

// parseDiffGitLine parses "diff --git a/X b/Y".
func parseDiffGitLine(line string, f *GitDiffFile) {
	// Remove "diff --git " prefix
	rest := strings.TrimPrefix(line, "diff --git ")
	// Format: "a/path b/path"
	// Split on " b/" - but paths can contain spaces, so we need to be careful.
	// Strategy: find " b/" that balances with the "a/" prefix.
	if strings.HasPrefix(rest, "a/") {
		// Find " b/" - the split point
		idx := strings.Index(rest, " b/")
		if idx >= 0 {
			oldP := rest[2:idx] // strip "a/"
			newP := rest[idx+3:] // strip " b/"
			f.OldPath = oldP
			f.NewPath = newP
		} else {
			// No " b/" found, try splitting on space
			parts := strings.SplitN(rest, " ", 2)
			if len(parts) == 2 {
				f.OldPath = strings.TrimPrefix(parts[0], "a/")
				f.NewPath = strings.TrimPrefix(parts[1], "b/")
			}
		}
	} else {
		// No a/ prefix
		parts := strings.SplitN(rest, " ", 2)
		if len(parts) == 2 {
			f.OldPath = parts[0]
			f.NewPath = parts[1]
		}
	}
}

// parseHunkHeader parses "@@ -10,5 +10,7 @@ func Foo() {".
func parseHunkHeader(line string) (oldStart, oldLines, newStart, newLines int, header string) {
	// Strip leading "@@ "
	line = strings.TrimPrefix(line, "@@ ")
	// Find closing " @@"
	end := strings.Index(line, " @@")
	if end >= 0 {
		header = strings.TrimSpace(line[end+3:])
		line = line[:end]
	}

	// line is now "-10,5 +10,7"
	parts := strings.Fields(line)
	for _, p := range parts {
		if strings.HasPrefix(p, "-") {
			oldStart, oldLines = parseHunkRange(p[1:])
		} else if strings.HasPrefix(p, "+") {
			newStart, newLines = parseHunkRange(p[1:])
		}
	}
	return
}

func parseHunkRange(s string) (start, count int) {
	comma := strings.IndexByte(s, ',')
	if comma < 0 {
		start, _ = strconv.Atoi(s)
		count = 1
		return
	}
	start, _ = strconv.Atoi(s[:comma])
	count, _ = strconv.Atoi(s[comma+1:])
	return
}
