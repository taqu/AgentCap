package gitparse

import (
	"strings"
)

// ParseBranch parses git branch output (with or without -a/-r flags).
func ParseBranch(raw []byte) *GitBranchList {
	if len(raw) == 0 {
		return &GitBranchList{}
	}

	text := string(raw)
	lines := strings.Split(text, "\n")

	result := &GitBranchList{}

	for _, line := range lines {
		if line == "" {
			continue
		}

		current := false
		if strings.HasPrefix(line, "* ") {
			current = true
			line = line[2:]
		} else if strings.HasPrefix(line, "  ") {
			line = line[2:]
		} else {
			continue
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Strip " -> ref" annotation (e.g. "origin/HEAD -> origin/main")
		if idx := strings.Index(line, " -> "); idx >= 0 {
			line = line[:idx]
		}

		// Detect remote branches
		isRemote := false
		if strings.HasPrefix(line, "remotes/") {
			isRemote = true
			line = strings.TrimPrefix(line, "remotes/")
		}

		entry := GitBranchEntry{
			Name:    line,
			Current: current,
			Remote:  isRemote,
		}

		if isRemote {
			result.Remote = append(result.Remote, entry)
		} else {
			result.Local = append(result.Local, entry)
			if current {
				result.Current = line
			}
		}
	}

	return result
}
