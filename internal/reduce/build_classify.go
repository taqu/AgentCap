package reduce

import (
	"path/filepath"
	"strings"
)

func classifyBuildCommand(args []string) string {
	if len(args) == 0 {
		return ""
	}
	base := filepath.Base(args[0])
	// Strip .exe on Windows
	base = strings.TrimSuffix(base, ".exe")

	switch base {
	case "go":
		for _, a := range args[1:] {
			if strings.HasPrefix(a, "-") {
				continue
			}
			switch a {
			case "build", "install", "vet":
				return "go-build"
			case "test":
				return "go-test"
			}
			return ""
		}
	case "gcc", "g++":
		return "gcc"
	case "clang", "clang++":
		return "clang"
	case "cargo":
		for _, a := range args[1:] {
			if strings.HasPrefix(a, "-") {
				continue
			}
			switch a {
			case "build":
				return "cargo-build"
			case "check":
				return "cargo-check"
			case "test":
				return "cargo-test"
			}
			return ""
		}
	case "make":
		return "make"
	case "ninja":
		return "ninja"
	}
	return ""
}

// SelectBuild picks the right build/test reducer.
func SelectBuild(args []string) Reducer {
	switch classifyBuildCommand(args) {
	case "go-build":
		return &GoBuildReducer{Args: args}
	case "go-test":
		return &GoTestReducer{Args: args}
	case "gcc", "g++":
		return &GccReducer{Args: args, Tool: "gcc"}
	case "clang", "clang++":
		return &GccReducer{Args: args, Tool: "clang"}
	case "cargo-build", "cargo-check":
		return &CargoBuildReducer{Args: args}
	case "cargo-test":
		return &CargoTestReducer{Args: args}
	case "make":
		return &MakeReducer{Args: args}
	case "ninja":
		return &NinjaReducer{Args: args}
	default:
		return &GenericReducer{}
	}
}
