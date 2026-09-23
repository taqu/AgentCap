// Package common provides shared helpers for AgentCap integration adapters.
package common

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/taqu/agentcap/internal/integration/protocol"
)

// IsBypassed returns true if ACAP_BYPASS=1 in the current environment.
func IsBypassed() bool {
	return os.Getenv(protocol.EnvBypass) == "1"
}

// IsRecursive returns true if ACAP_INTERCEPT_DEPTH > 0 (we are already inside an interception).
func IsRecursive() bool {
	n, _ := strconv.Atoi(os.Getenv(protocol.EnvDepth))
	return n > 0
}

// IsAcapCommand returns true if args[0] appears to be the acap binary itself.
// This prevents recursive interception of "acap show", "acap raw", etc.
func IsAcapCommand(args []string) bool {
	if len(args) == 0 {
		return false
	}
	base := filepath.Base(args[0])
	return base == "acap" || base == "acap.exe"
}

// EnvWithDepth returns a copy of env with ACAP_INTERCEPT_DEPTH incremented by 1.
// env is a []string in "KEY=VALUE" form (os.Environ() format).
func EnvWithDepth(env []string) []string {
	key := protocol.EnvDepth + "="
	current := 0
	result := make([]string, 0, len(env)+1)
	found := false
	for _, e := range env {
		if strings.HasPrefix(e, key) {
			current, _ = strconv.Atoi(e[len(key):])
			result = append(result, key+strconv.Itoa(current+1))
			found = true
		} else {
			result = append(result, e)
		}
	}
	if !found {
		result = append(result, key+"1")
	}
	return result
}
