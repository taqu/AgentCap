// Package project provides project-root discovery for acap.
package project

import (
	"os"
	"path/filepath"
)

const EnvRoot = "ACAP_ROOT"

// FindRoot discovers the AgentCap project root starting from cwd.
// Priority: ACAP_ROOT env > nearest .acap/ ancestor > nearest .git ancestor > cwd
func FindRoot(cwd string) string {
	if root := os.Getenv(EnvRoot); root != "" {
		return filepath.Clean(root)
	}

	// Walk up: .acap/ wins first.
	dir := filepath.Clean(cwd)
	for {
		if fi, err := os.Stat(filepath.Join(dir, ".acap")); err == nil && fi.IsDir() {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// Walk up: .git wins second.
	dir = filepath.Clean(cwd)
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// Fallback: cwd itself.
	return filepath.Clean(cwd)
}
