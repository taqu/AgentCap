package reduce

import (
	"fmt"
	"strings"

	"github.com/taqu/agentcap/internal/clean"
	"github.com/taqu/agentcap/internal/exec"
	"github.com/taqu/agentcap/internal/gitparse"
)

// GitDiffReducer handles git diff output.
type GitDiffReducer struct {
	Args       []string
	ParsedDiff *gitparse.GitDiff // exposed so main.go can store it
}

func (r *GitDiffReducer) Reduce(result *exec.Result) *ReducedResult {
	raw := result.Stdout
	rawBytes := len(raw)

	if result.ExitCode != 0 {
		out := string(clean.StripANSI(raw)) + buildStderr(result.Stderr)
		return &ReducedResult{
			Output:   out,
			RawBytes: rawBytes,
			RetBytes: rawBytes,
		}
	}

	// Check for --stat flag.
	isStat := hasDiffStatFlag(r.Args)

	if isStat {
		// --stat output is already compact; minimal processing.
		cleaned := string(clean.StripANSI(raw))
		out := "@acap git-diff stat\n" + cleaned + buildStderr(result.Stderr)
		return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
	}

	// Empty diff (no output or just newline).
	if len(strings.TrimSpace(string(raw))) == 0 {
		scope := diffScope(r.Args)
		header := "@acap git-diff"
		if scope != "" {
			header += " " + scope
		}
		out := header + "\nno changes\n"
		return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
	}

	parsed := gitparse.ParseDiff(clean.StripANSI(raw))
	if parsed == nil || len(parsed.Files) == 0 {
		// Fall back to generic.
		return (&GenericReducer{}).Reduce(result)
	}
	r.ParsedDiff = parsed

	scope := diffScope(r.Args)
	totalAdd, totalDel := parsed.Totals()

	var sb strings.Builder
	fmt.Fprintf(&sb, "@acap git-diff")
	if scope != "" {
		fmt.Fprintf(&sb, " %s", scope)
	}
	sb.WriteByte('\n')
	fmt.Fprintf(&sb, "files=%d +%d -%d\n", len(parsed.Files), totalAdd, totalDel)

	// Binary files section.
	var binaryFiles []gitparse.GitDiffFile
	for _, f := range parsed.Files {
		if f.Binary {
			binaryFiles = append(binaryFiles, f)
		}
	}

	sb.WriteByte('\n')
	for _, f := range parsed.Files {
		if f.Binary {
			continue
		}
		renderDiffFileLine(&sb, f)
	}

	if len(binaryFiles) > 0 {
		fmt.Fprintf(&sb, "\nbinary=%d\n", len(binaryFiles))
		for _, f := range binaryFiles {
			fmt.Fprintf(&sb, "%s\n", f.Path())
		}
	}

	sb.WriteString(buildStderr(result.Stderr))
	out := sb.String()
	return &ReducedResult{Output: out, RawBytes: rawBytes, RetBytes: len(out)}
}

func renderDiffFileLine(sb *strings.Builder, f gitparse.GitDiffFile) {
	// Status prefix.
	status := f.Status
	if f.Similarity > 0 {
		status = fmt.Sprintf("R%d", f.Similarity)
	}

	// Path.
	var pathStr string
	if f.OldPath != "" && f.NewPath != "" && f.OldPath != f.NewPath {
		pathStr = f.OldPath + " -> " + f.NewPath
	} else {
		pathStr = f.Path()
	}

	// Mode change?
	if f.OldMode != "" && f.NewMode != "" {
		fmt.Fprintf(sb, "%s %s mode %s -> %s\n", status, pathStr, f.OldMode, f.NewMode)
		return
	}

	// Additions/deletions + hunks.
	hunks := len(f.Hunks)
	switch status {
	case "A":
		fmt.Fprintf(sb, "A %s +%d", pathStr, f.Additions)
	case "D":
		fmt.Fprintf(sb, "D %s -%d", pathStr, f.Deletions)
	default:
		fmt.Fprintf(sb, "%s %s +%d -%d", status, pathStr, f.Additions, f.Deletions)
	}
	if hunks > 0 {
		fmt.Fprintf(sb, " hunks=%d", hunks)
	}
	sb.WriteByte('\n')

	// Show hunk headers.
	for _, h := range f.Hunks {
		if h.Header != "" {
			fmt.Fprintf(sb, "  h%d @@ -%d,%d +%d,%d @@ %s\n", h.Index+1, h.OldStart, h.OldLines, h.NewStart, h.NewLines, h.Header)
		} else {
			fmt.Fprintf(sb, "  h%d @@ -%d,%d +%d,%d @@\n", h.Index+1, h.OldStart, h.OldLines, h.NewStart, h.NewLines)
		}
	}
}

// diffScope returns "staged", "cached", or "" based on args.
func diffScope(args []string) string {
	for _, a := range args {
		if a == "--cached" || a == "--staged" {
			return "staged"
		}
	}
	return ""
}

func hasDiffStatFlag(args []string) bool {
	for _, a := range args {
		if a == "--stat" || a == "--shortstat" || strings.HasPrefix(a, "--stat=") {
			return true
		}
	}
	return false
}
