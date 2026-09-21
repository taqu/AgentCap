package session

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// EnvKey is the environment variable name for the session ID.
const EnvKey = "ACAP_SESSION_ID"

// Session groups command executions for cross-command comparison.
type Session struct {
	ID      string
	dir     string
	history *History
}

// Current returns the session identified by ACAP_SESSION_ID, or nil if not set.
func Current() (*Session, error) {
	id := os.Getenv(EnvKey)
	if id == "" {
		return nil, nil
	}
	return open(id)
}

// Open opens or creates a session by ID.
func Open(id string) (*Session, error) {
	return open(id)
}

// GenerateID generates a new random 6-hex session ID.
func GenerateID() (string, error) {
	b := make([]byte, 3)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("session: generate id: %w", err)
	}
	return hex.EncodeToString(b), nil
}

func open(id string) (*Session, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("session: cache dir: %w", err)
	}
	dir := filepath.Join(cacheDir, "agentcap", "sessions", id)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, fmt.Errorf("session: create dir: %w", err)
	}
	return &Session{
		ID:      id,
		dir:     dir,
		history: newHistory(filepath.Join(dir, "history.jsonl")),
	}, nil
}

// History returns the underlying history.
func (s *Session) History() *History {
	return s.history
}

// LatestBaseline finds the most recent stored result with the same CommandKey.
func (s *Session) LatestBaseline(key CommandKey) (*HistoryRecord, error) {
	return s.history.LatestBaseline(key.Hash())
}

// Record appends this execution to the session history.
func (s *Session) Record(resultID string, key CommandKey, exitCode int, stdoutHash, stderrHash, presentation string) error {
	return s.history.Append(HistoryRecord{
		ResultID:     resultID,
		CommandKey:   key.Hash(),
		CreatedAt:    time.Now().UTC(),
		StdoutHash:   stdoutHash,
		StderrHash:   stderrHash,
		ExitCode:     exitCode,
		Presentation: presentation,
	})
}

// SessionInfo holds summary info about a session.
type SessionInfo struct {
	ID      string
	Count   int
	Updated time.Time
}

// List returns all existing sessions.
func List() ([]SessionInfo, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("session: cache dir: %w", err)
	}
	sessionsDir := filepath.Join(cacheDir, "agentcap", "sessions")
	entries, err := os.ReadDir(sessionsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var result []SessionInfo
	for _, de := range entries {
		if !de.IsDir() {
			continue
		}
		sess, err := open(de.Name())
		if err != nil {
			continue
		}
		records, _ := sess.history.Records()
		var updated time.Time
		if len(records) > 0 {
			updated = records[len(records)-1].CreatedAt
		}
		result = append(result, SessionInfo{
			ID:      de.Name(),
			Count:   len(records),
			Updated: updated,
		})
	}
	return result, nil
}
