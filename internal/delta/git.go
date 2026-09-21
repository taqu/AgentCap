package delta

import (
	"context"
	"fmt"
	"strings"

	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/gitparse"
	"github.com/taqu/agentcap/internal/store"
)

// GitStatusDelta implements Reducer for git status.
type GitStatusDelta struct{}

func (d *GitStatusDelta) Delta(ctx context.Context, baseline, current *store.Entry) (*Result, error) {
	prevRaw := readFile(baseline.StdoutPath())
	currRaw := readFile(current.StdoutPath())
	if prevRaw == nil || currRaw == nil {
		return nil, nil
	}

	prevStatus := gitparse.ParseStatus(clean.StripANSI(prevRaw))
	currStatus := gitparse.ParseStatus(clean.StripANSI(currRaw))
	if prevStatus == nil || currStatus == nil {
		return nil, nil
	}

	// Build file maps keyed by path.
	prevMap := statusFileMap(prevStatus.Files)
	currMap := statusFileMap(currStatus.Files)

	var added, removed, changed []string
	for path, f := range currMap {
		if prev, ok := prevMap[path]; !ok {
			added = append(added, fmt.Sprintf("%s %s", strings.TrimSpace(f.XY), path))
		} else if prev.XY != f.XY {
			changed = append(changed, fmt.Sprintf("%s %s -> %s", path, strings.TrimSpace(prev.XY), strings.TrimSpace(f.XY)))
		}
	}
	for path, f := range prevMap {
		if _, ok := currMap[path]; !ok {
			removed = append(removed, fmt.Sprintf("%s %s", strings.TrimSpace(f.XY), path))
		}
	}

	if len(added)+len(removed)+len(changed) == 0 {
		return nil, nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap delta from %s git-status\n", baseline.Meta.ID)
	if len(added) > 0 {
		sb.WriteString("\nadded:\n")
		for _, l := range added {
			sb.WriteString(l + "\n")
		}
	}
	if len(removed) > 0 {
		sb.WriteString("\nresolved:\n")
		for _, l := range removed {
			sb.WriteString(l + "\n")
		}
	}
	if len(changed) > 0 {
		sb.WriteString("\nchanged:\n")
		for _, l := range changed {
			sb.WriteString(l + "\n")
		}
	}
	return &Result{Output: sb.String(), Presentation: PresentationDelta}, nil
}

func statusFileMap(files []gitparse.GitStatusFile) map[string]gitparse.GitStatusFile {
	m := make(map[string]gitparse.GitStatusFile, len(files))
	for _, f := range files {
		m[f.Path] = f
	}
	return m
}

// GitDiffDelta implements Reducer for git diff structural comparison.
type GitDiffDelta struct {
	Store *store.Store
}

func (d *GitDiffDelta) Delta(ctx context.Context, baseline, current *store.Entry) (*Result, error) {
	// Load baseline files from store.
	prevFiles, err := d.Store.GetGitDiffFiles(baseline.Meta.ID)
	if err != nil || len(prevFiles) == 0 {
		return nil, nil
	}

	currRaw := readFile(current.StdoutPath())
	if currRaw == nil {
		return nil, nil
	}

	currDiff := gitparse.ParseDiff(clean.StripANSI(currRaw))
	if currDiff == nil {
		return nil, nil
	}

	// Build maps by path.
	prevMap := diffFileMap(prevFiles)
	currMap := diffFileMap(currDiff.Files)

	var changedFiles []gitparse.GitDiffFile
	var newFiles []gitparse.GitDiffFile
	var resolvedPaths []string
	unchangedCount := 0

	for path, cf := range currMap {
		if pf, ok := prevMap[path]; ok {
			if cf.Additions == pf.Additions && cf.Deletions == pf.Deletions && cf.Status == pf.Status {
				unchangedCount++
			} else {
				changedFiles = append(changedFiles, cf)
			}
		} else {
			newFiles = append(newFiles, cf)
		}
	}
	for path := range prevMap {
		if _, ok := currMap[path]; !ok {
			resolvedPaths = append(resolvedPaths, path)
		}
	}

	if len(changedFiles)+len(newFiles)+len(resolvedPaths) == 0 {
		return nil, nil
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap delta from %s git-diff\n", baseline.Meta.ID)

	if len(newFiles) > 0 {
		sb.WriteString("\nnew files:\n")
		for _, f := range newFiles {
			renderDeltaFileLine(&sb, f)
		}
	}
	if len(changedFiles) > 0 {
		sb.WriteString("\nchanged:\n")
		for _, f := range changedFiles {
			renderDeltaFileLine(&sb, f)
		}
	}
	if len(resolvedPaths) > 0 {
		sb.WriteString("\nresolved:\n")
		for _, p := range resolvedPaths {
			sb.WriteString(p + "\n")
		}
	}
	if unchangedCount > 0 {
		fmt.Fprintf(&sb, "\nunchanged_files=%d\n", unchangedCount)
	}

	return &Result{Output: sb.String(), Presentation: PresentationDelta}, nil
}

func diffFileMap(files []gitparse.GitDiffFile) map[string]gitparse.GitDiffFile {
	m := make(map[string]gitparse.GitDiffFile, len(files))
	for _, f := range files {
		m[f.Path()] = f
	}
	return m
}

func renderDeltaFileLine(sb *strings.Builder, f gitparse.GitDiffFile) {
	fmt.Fprintf(sb, "%s %s +%d -%d\n", f.Status, f.Path(), f.Additions, f.Deletions)
}
