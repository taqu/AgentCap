package reduce

import (
	"path/filepath"
	"strings"
)

// isGitExecutable returns true if the executable name is "git" or ends with "/git".
func isGitExecutable(name string) bool {
	base := filepath.Base(name)
	return base == "git" || base == "git.exe"
}

// classifyGitSubcommand returns the git subcommand: "status", "diff", "show", "log", "branch", or "".
func classifyGitSubcommand(args []string) string {
	// args[0] is "git", args[1] would be subcommand (possibly with flags before it)
	for i := 1; i < len(args); i++ {
		if strings.HasPrefix(args[i], "-") {
			continue // skip git-level flags like -C, --git-dir
		}
		switch args[i] {
		case "status":
			return "status"
		case "diff":
			return "diff"
		case "show":
			return "show"
		case "log":
			return "log"
		case "branch":
			return "branch"
		}
		return "" // unknown subcommand
	}
	return ""
}

// SelectGit picks the right git reducer based on subcommand.
func SelectGit(args []string) Reducer {
	switch classifyGitSubcommand(args) {
	case "status":
		return &GitStatusReducer{Args: args}
	case "diff":
		return &GitDiffReducer{Args: args}
	case "show":
		return &GitShowReducer{Args: args}
	case "log":
		return &GitLogReducer{Args: args}
	case "branch":
		return &GitBranchReducer{Args: args}
	default:
		return &GenericReducer{}
	}
}
