package session

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// HistoryRecord is one entry in a session's JSONL history.
type HistoryRecord struct {
	Seq          int       `json:"seq"`
	ResultID     string    `json:"result_id"`
	CommandKey   string    `json:"command_key"` // hash
	CreatedAt    time.Time `json:"created_at"`
	StdoutHash   string    `json:"stdout_hash"`
	StderrHash   string    `json:"stderr_hash"`
	ExitCode     int       `json:"exit_code"`
	Presentation string    `json:"presentation"` // "full", "unchanged", "delta"
}

// History is an append-only JSONL session history file.
type History struct {
	path string
	mu   sync.Mutex
}

func newHistory(path string) *History {
	return &History{path: path}
}

// Append writes a new record, assigning the next sequence number.
func (h *History) Append(rec HistoryRecord) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	records, _ := h.readAll()
	rec.Seq = len(records) + 1

	f, err := os.OpenFile(h.path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o600)
	if err != nil {
		return fmt.Errorf("session: append: %w", err)
	}
	defer f.Close()

	data, err := json.Marshal(rec)
	if err != nil {
		return fmt.Errorf("session: marshal: %w", err)
	}
	_, err = fmt.Fprintf(f, "%s\n", data)
	return err
}

// LatestBaseline returns the most recent record matching commandKeyHash, or nil.
func (h *History) LatestBaseline(commandKeyHash string) (*HistoryRecord, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	records, err := h.readAll()
	if err != nil {
		return nil, err
	}
	for i := len(records) - 1; i >= 0; i-- {
		if records[i].CommandKey == commandKeyHash {
			cp := records[i]
			return &cp, nil
		}
	}
	return nil, nil
}

// Records returns all records in order.
func (h *History) Records() ([]HistoryRecord, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.readAll()
}

func (h *History) readAll() ([]HistoryRecord, error) {
	f, err := os.Open(h.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("session: read history: %w", err)
	}
	defer f.Close()

	var records []HistoryRecord
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 64*1024)
	for scanner.Scan() {
		var rec HistoryRecord
		if err := json.Unmarshal(scanner.Bytes(), &rec); err != nil {
			continue // skip corrupt entries
		}
		records = append(records, rec)
	}
	return records, nil
}
