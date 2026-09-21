// Package gitparse provides structured parsers for common git command output.
package gitparse

// GitStatusFile represents one file line in git status output.
type GitStatusFile struct {
	Path    string // current path (new path for renames)
	OldPath string // original path for renames/copies
	XY      string // two-character status code (e.g. " M", "M ", "??", "UU")
}

// IsConflicted returns true for merge conflict statuses.
func (f GitStatusFile) IsConflicted() bool {
	if len(f.XY) < 2 {
		return false
	}
	switch f.XY {
	case "DD", "AU", "UD", "UA", "DU", "AA", "UU":
		return true
	}
	return false
}

// IsUntracked returns true for ?? status.
func (f GitStatusFile) IsUntracked() bool { return f.XY == "??" }

// GitStatus is the structured result of git status parsing.
type GitStatus struct {
	Branch   string
	Upstream string
	Ahead    int
	Behind   int
	Detached bool
	HeadHash string
	Files    []GitStatusFile
	IsClean  bool
}

// GitDiffHunk represents a single hunk in a unified diff.
type GitDiffHunk struct {
	Index    int
	OldStart int
	OldLines int
	NewStart int
	NewLines int
	Header   string // function context after @@
	RawStart int64  // byte offset in raw stdout (start of @@ line)
	RawEnd   int64  // exclusive end byte offset
}

// GitDiffFile represents a single changed file in a git diff.
type GitDiffFile struct {
	Index      int
	OldPath    string
	NewPath    string
	Status     string // M, A, D, R, C, T
	Additions  int
	Deletions  int
	Binary     bool
	OldMode    string
	NewMode    string
	Similarity int // rename/copy percentage
	Hunks      []GitDiffHunk
	RawStart   int64 // byte offset of "diff --git" line
	RawEnd     int64 // exclusive end
}

// Path returns NewPath or OldPath if NewPath is empty.
func (f *GitDiffFile) Path() string {
	if f.NewPath != "" {
		return f.NewPath
	}
	return f.OldPath
}

// GitDiff is the structured result of git diff parsing.
type GitDiff struct {
	Files    []GitDiffFile
	StatOnly bool
}

// Totals returns total additions, deletions.
func (d *GitDiff) Totals() (add, del int) {
	for _, f := range d.Files {
		add += f.Additions
		del += f.Deletions
	}
	return
}

// GitCommit represents commit metadata.
type GitCommit struct {
	Hash    string
	Author  string
	Date    string
	Subject string
	Body    string
}

// GitShow is the structured result of git show.
type GitShow struct {
	Commit GitCommit
	Diff   *GitDiff
}

// GitLog is the structured result of git log.
type GitLog struct {
	Commits []GitCommit
	Total   int
}

// GitBranchEntry is a single branch in git branch output.
type GitBranchEntry struct {
	Name    string
	Current bool
	Remote  bool
}

// GitBranchList is the structured result of git branch.
type GitBranchList struct {
	Current string
	Local   []GitBranchEntry
	Remote  []GitBranchEntry
}
