package gitparse

import (
	"bytes"
	"strings"
)

const maxLogCommits = 100

// ParseLog parses git log output (standard verbose or --oneline format).
func ParseLog(raw []byte) *GitLog {
	if len(raw) == 0 {
		return &GitLog{}
	}

	text := string(raw)
	lines := strings.Split(text, "\n")

	gl := &GitLog{}

	// Detect oneline format: lines like "abc1234 subject" (7-char hex + space).
	if isOnelineFormat(lines) {
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			// Could be "* abc1234 subject" (from --graph) - strip graph chars
			line = stripGraphPrefix(line)
			sp := strings.IndexByte(line, ' ')
			if sp < 0 {
				continue
			}
			hash := line[:sp]
			if !isHex(hash) {
				continue
			}
			subject := strings.TrimSpace(line[sp+1:])
			gl.Total++
			if len(gl.Commits) < maxLogCommits {
				gl.Commits = append(gl.Commits, GitCommit{Hash: hash, Subject: subject})
			}
		}
		return gl
	}

	// Verbose format.
	parseVerboseLog(lines, gl)
	return gl
}

func parseVerboseLog(lines []string, gl *GitLog) {
	var current *GitCommit
	inBody := false
	bodyLines := []string{}

	flush := func() {
		if current == nil {
			return
		}
		if len(bodyLines) > 0 {
			current.Body = strings.Join(bodyLines, "\n")
		}
		gl.Total++
		if len(gl.Commits) < maxLogCommits {
			gl.Commits = append(gl.Commits, *current)
		}
		current = nil
		bodyLines = bodyLines[:0]
		inBody = false
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "commit ") {
			flush()
			hash := strings.TrimSpace(strings.TrimPrefix(line, "commit "))
			// strip parenthetical "(HEAD -> main)"
			if sp := strings.IndexByte(hash, ' '); sp >= 0 {
				hash = hash[:sp]
			}
			current = &GitCommit{Hash: hash}
			inBody = false
		} else if current != nil {
			if strings.HasPrefix(line, "Author: ") {
				current.Author = strings.TrimPrefix(line, "Author: ")
				current.Author = strings.TrimSpace(current.Author)
				inBody = false
			} else if strings.HasPrefix(line, "AuthorDate: ") || strings.HasPrefix(line, "Date: ") {
				dateStr := line
				if strings.HasPrefix(line, "AuthorDate: ") {
					dateStr = strings.TrimPrefix(line, "AuthorDate: ")
				} else {
					dateStr = strings.TrimPrefix(line, "Date: ")
				}
				current.Date = strings.TrimSpace(dateStr)
				inBody = false
			} else if strings.HasPrefix(line, "Merge: ") ||
				strings.HasPrefix(line, "Commit: ") ||
				strings.HasPrefix(line, "CommitDate: ") {
				// skip
			} else if line == "" {
				if current.Subject != "" {
					inBody = true
				}
			} else if !inBody {
				// Subject line (first non-blank after headers)
				trimmed := strings.TrimSpace(line)
				if trimmed != "" && current.Subject == "" {
					current.Subject = trimmed
				}
			} else {
				bodyLines = append(bodyLines, strings.TrimSpace(line))
			}
		}
	}
	flush()
}

// ParseShow parses git show output (commit header + optional diff).
func ParseShow(raw []byte) *GitShow {
	if len(raw) == 0 {
		return &GitShow{}
	}

	// Find the diff --git section if present.
	diffPrefix := []byte("\ndiff --git ")
	diffIdx := bytes.Index(raw, diffPrefix)

	var commitBytes, diffBytes []byte
	if diffIdx >= 0 {
		commitBytes = raw[:diffIdx]
		diffBytes = raw[diffIdx+1:] // skip leading \n
	} else {
		commitBytes = raw
	}

	gs := &GitShow{}

	// Parse commit header (verbose log format).
	gl := ParseLog(commitBytes)
	if len(gl.Commits) > 0 {
		gs.Commit = gl.Commits[0]
	}

	if len(diffBytes) > 0 {
		gs.Diff = ParseDiff(diffBytes)
	}

	return gs
}

func isOnelineFormat(lines []string) bool {
	// Check first non-empty line.
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		line = stripGraphPrefix(line)
		if strings.HasPrefix(line, "commit ") {
			return false
		}
		// Check if it looks like "7hexchars subject"
		sp := strings.IndexByte(line, ' ')
		if sp >= 4 && sp <= 40 && isHex(line[:sp]) {
			return true
		}
		return false
	}
	return false
}

func stripGraphPrefix(line string) string {
	// Strip git log --graph decorations like "* ", "| * ", "| |/ " etc.
	for i, c := range line {
		if c == '*' {
			rest := strings.TrimSpace(line[i+1:])
			return rest
		}
		if c != '|' && c != ' ' && c != '/' && c != '\\' {
			return line[i:]
		}
	}
	return line
}

func isHex(s string) bool {
	if len(s) == 0 {
		return false
	}
	for _, c := range s {
		if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f') || (c >= 'A' && c <= 'F')) {
			return false
		}
	}
	return true
}
