package session

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// CommandKey identifies a command execution for baseline comparison.
type CommandKey struct {
	Args []string
	Cwd  string
}

// Hash returns a stable hex string for this key.
func (k CommandKey) Hash() string {
	h := sha256.Sum256([]byte(fmt.Sprintf("cwd=%s\nargs=%s", k.Cwd, strings.Join(k.Args, "\x00"))))
	return hex.EncodeToString(h[:])
}

// NewKey creates a CommandKey from args and working directory.
func NewKey(args []string, cwd string) CommandKey {
	return CommandKey{Args: args, Cwd: cwd}
}
