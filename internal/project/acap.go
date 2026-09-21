package project

import (
	"os"
	"path/filepath"
)

// AcapDir returns the .acap directory path under root.
func AcapDir(root string) string {
	return filepath.Join(root, ".acap")
}

// DBPath returns the canonical SQLite database path.
func DBPath(root string) string {
	return filepath.Join(root, ".acap", "store.db")
}

// ObjectsDir returns the objects directory path.
func ObjectsDir(root string) string {
	return filepath.Join(root, ".acap", "objects")
}

// EnsureDirs creates the .acap and objects directories under root.
func EnsureDirs(root string) error {
	if err := os.MkdirAll(filepath.Join(root, ".acap", "objects"), 0o700); err != nil {
		return err
	}
	return nil
}
