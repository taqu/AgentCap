package gitparse

import (
	"strconv"
	"strings"
)

// ParseStatus parses git status output (verbose or short/porcelain).
func ParseStatus(raw []byte) *GitStatus {
	if len(raw) == 0 {
		return &GitStatus{IsClean: true}
	}

	text := string(raw)
	lines := strings.Split(text, "\n")

	gs := &GitStatus{}

	// Detect format: verbose starts with "On branch" or "HEAD detached"
	// or "nothing to commit". Short format lines start with 2-char status
	// or "## " for header.
	isVerbose := false
	for _, line := range lines {
		if strings.HasPrefix(line, "On branch ") ||
			strings.HasPrefix(line, "HEAD detached at") ||
			strings.HasPrefix(line, "nothing to commit") ||
			strings.HasPrefix(line, "No commits yet") ||
			strings.HasPrefix(line, "Changes to be committed") ||
			strings.HasPrefix(line, "Changes not staged") ||
			strings.HasPrefix(line, "Untracked files:") {
			isVerbose = true
			break
		}
	}

	if isVerbose {
		parseVerboseStatus(lines, gs)
	} else {
		parseShortStatus(lines, gs)
	}

	return gs
}

func parseVerboseStatus(lines []string, gs *GitStatus) {
	// Sections: 0=none, 1=staged, 2=unstaged, 3=untracked
	section := 0

	for _, line := range lines {
		// Branch line
		if strings.HasPrefix(line, "On branch ") {
			gs.Branch = strings.TrimPrefix(line, "On branch ")
			gs.Branch = strings.TrimSpace(gs.Branch)
			continue
		}
		if strings.HasPrefix(line, "HEAD detached at ") {
			gs.Detached = true
			gs.HeadHash = strings.TrimSpace(strings.TrimPrefix(line, "HEAD detached at "))
			gs.Branch = "DETACHED"
			continue
		}

		// Ahead/behind
		if strings.Contains(line, "Your branch is ahead of") {
			// "Your branch is ahead of 'origin/main' by 2 commits."
			gs.Upstream = extractQuoted(line)
			gs.Ahead = extractNumber(line, "by ", " commit")
			continue
		}
		if strings.Contains(line, "Your branch is behind") {
			gs.Upstream = extractQuoted(line)
			gs.Behind = extractNumber(line, "by ", " commit")
			continue
		}
		if strings.Contains(line, "Your branch and") && strings.Contains(line, "have diverged") {
			gs.Upstream = extractQuoted(line)
			continue
		}

		// Clean
		if strings.Contains(line, "nothing to commit") {
			gs.IsClean = true
			continue
		}

		// Section headers
		if strings.HasPrefix(line, "Changes to be committed:") {
			section = 1
			continue
		}
		if strings.HasPrefix(line, "Changes not staged for commit:") {
			section = 2
			continue
		}
		if strings.HasPrefix(line, "Untracked files:") {
			section = 3
			continue
		}
		// Other section-like lines that reset context
		if line == "" {
			// blank lines don't reset section
			continue
		}

		// File lines are indented
		trimmed := strings.TrimLeft(line, "\t ")
		if trimmed == line {
			// Not indented - could be a new top-level section or informational line
			// Reset section if it looks like a new header
			if strings.HasSuffix(line, ":") {
				section = 0
			}
			continue
		}

		switch section {
		case 1:
			parseVerboseStagedLine(trimmed, gs)
		case 2:
			parseVerboseUnstagedLine(trimmed, gs)
		case 3:
			// Untracked - skip "use ... to include" hints
			if strings.HasPrefix(trimmed, "(use") || trimmed == "" {
				continue
			}
			gs.Files = append(gs.Files, GitStatusFile{Path: trimmed, XY: "??"})
		}
	}
}

func parseVerboseStagedLine(line string, gs *GitStatus) {
	// Format: "statusword:   path" or "renamed:   old -> new"
	if strings.HasPrefix(line, "(use") {
		return
	}
	colon := strings.IndexByte(line, ':')
	if colon < 0 {
		return
	}
	word := strings.TrimSpace(line[:colon])
	path := strings.TrimSpace(line[colon+1:])
	if path == "" {
		return
	}

	var xy string
	switch word {
	case "modified":
		xy = "M "
	case "new file":
		xy = "A "
	case "deleted":
		xy = "D "
	case "renamed":
		xy = "R "
	case "copied":
		xy = "C "
	case "typechange":
		xy = "T "
	default:
		xy = "M "
	}

	f := GitStatusFile{XY: xy}
	if word == "renamed" || word == "copied" {
		// "old -> new"
		parts := strings.SplitN(path, " -> ", 2)
		if len(parts) == 2 {
			f.OldPath = strings.TrimSpace(parts[0])
			f.Path = strings.TrimSpace(parts[1])
		} else {
			f.Path = path
		}
	} else {
		f.Path = path
	}
	gs.Files = append(gs.Files, f)
}

func parseVerboseUnstagedLine(line string, gs *GitStatus) {
	if strings.HasPrefix(line, "(use") {
		return
	}
	colon := strings.IndexByte(line, ':')
	if colon < 0 {
		return
	}
	word := strings.TrimSpace(line[:colon])
	path := strings.TrimSpace(line[colon+1:])
	if path == "" {
		return
	}

	var xy string
	switch word {
	case "modified":
		xy = " M"
	case "deleted":
		xy = " D"
	case "renamed":
		xy = " R"
	default:
		xy = " M"
	}

	f := GitStatusFile{XY: xy, Path: path}
	if word == "renamed" {
		parts := strings.SplitN(path, " -> ", 2)
		if len(parts) == 2 {
			f.OldPath = strings.TrimSpace(parts[0])
			f.Path = strings.TrimSpace(parts[1])
		}
	}
	gs.Files = append(gs.Files, f)
}

func parseShortStatus(lines []string, gs *GitStatus) {
	for _, line := range lines {
		if line == "" {
			continue
		}
		// Branch header: "## main...origin/main [ahead 2]"
		if strings.HasPrefix(line, "## ") {
			parseShortBranchHeader(line[3:], gs)
			continue
		}
		// Need at least 3 chars: XY + space
		if len(line) < 3 {
			continue
		}
		xy := line[:2]
		rest := line[3:] // skip XY + space
		f := GitStatusFile{XY: xy}
		// Renames: "old -> new"
		if strings.Contains(rest, " -> ") {
			parts := strings.SplitN(rest, " -> ", 2)
			f.OldPath = strings.TrimSpace(parts[0])
			f.Path = strings.TrimSpace(parts[1])
		} else {
			f.Path = rest
		}
		gs.Files = append(gs.Files, f)
	}
	if len(gs.Files) == 0 && gs.Branch == "" {
		gs.IsClean = true
	}
}

func parseShortBranchHeader(header string, gs *GitStatus) {
	// "main...origin/main [ahead 2, behind 1]"
	// "HEAD (no branch)"
	// "main"

	if strings.HasPrefix(header, "HEAD (no branch)") {
		gs.Detached = true
		gs.Branch = "DETACHED"
		return
	}

	// Split off tracking info
	dotdot := strings.Index(header, "...")
	if dotdot >= 0 {
		gs.Branch = strings.TrimSpace(header[:dotdot])
		rest := header[dotdot+3:]
		// rest: "origin/main [ahead 2]" or "origin/main [behind 3]" etc.
		bracket := strings.IndexByte(rest, '[')
		if bracket >= 0 {
			gs.Upstream = strings.TrimSpace(rest[:bracket])
			info := rest[bracket+1:]
			end := strings.IndexByte(info, ']')
			if end >= 0 {
				info = info[:end]
			}
			for _, part := range strings.Split(info, ",") {
				part = strings.TrimSpace(part)
				if strings.HasPrefix(part, "ahead ") {
					n, _ := strconv.Atoi(strings.TrimPrefix(part, "ahead "))
					gs.Ahead = n
				} else if strings.HasPrefix(part, "behind ") {
					n, _ := strconv.Atoi(strings.TrimPrefix(part, "behind "))
					gs.Behind = n
				}
			}
		} else {
			gs.Upstream = strings.TrimSpace(rest)
		}
	} else {
		gs.Branch = strings.TrimSpace(header)
	}
}

// extractQuoted returns the first single-quoted string in s.
func extractQuoted(s string) string {
	start := strings.IndexByte(s, '\'')
	if start < 0 {
		return ""
	}
	end := strings.IndexByte(s[start+1:], '\'')
	if end < 0 {
		return ""
	}
	return s[start+1 : start+1+end]
}

// extractNumber extracts an integer from s that appears between prefix and suffix.
func extractNumber(s, prefix, suffix string) int {
	idx := strings.Index(s, prefix)
	if idx < 0 {
		return 0
	}
	rest := s[idx+len(prefix):]
	end := strings.Index(rest, suffix)
	if end < 0 {
		// try to read until space
		end = strings.IndexByte(rest, ' ')
		if end < 0 {
			end = len(rest)
		}
	}
	n, _ := strconv.Atoi(strings.TrimSpace(rest[:end]))
	return n
}
