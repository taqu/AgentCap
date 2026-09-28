package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// BenchmarkRecord is one persisted command measurement in a benchmark session.
type BenchmarkRecord struct {
	SessionID string
	Sequence  int
	ResultID  string
	Command   []string
	ExitCode  int

	RawStdoutBytes int64
	RawStderrBytes int64
	RawBytes       int64
	StatelessBytes int64
	StatefulBytes  int64
	Presentation   string
	Truncated      bool

	ExecutionDuration  time.Duration
	ReduceDuration     time.Duration
	ProcessingDuration time.Duration
}

// CreateBenchmarkSession registers an isolated benchmark session boundary.
func (s *Store) CreateBenchmarkSession(id string) error {
	_, err := s.db.ExecContext(context.Background(),
		`INSERT INTO benchmark_sessions(id, created_at) VALUES(?, ?)`,
		id, time.Now().UTC().Format(time.RFC3339Nano))
	if err != nil {
		return fmt.Errorf("store: create benchmark session: %w", err)
	}
	return nil
}

// BenchmarkSessionExists reports whether id was created as a benchmark session.
func (s *Store) BenchmarkSessionExists(id string) (bool, error) {
	var n int
	err := s.db.QueryRowContext(context.Background(),
		`SELECT COUNT(*) FROM benchmark_sessions WHERE id=?`, id).Scan(&n)
	return n != 0, err
}

// SaveBenchmarkRecord appends rec and assigns its session sequence number.
func (s *Store) SaveBenchmarkRecord(rec *BenchmarkRecord) error {
	commandJSON, err := json.Marshal(rec.Command)
	if err != nil {
		return fmt.Errorf("store: marshal benchmark command: %w", err)
	}
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := tx.QueryRowContext(context.Background(),
		`SELECT COALESCE(MAX(sequence), 0) + 1 FROM benchmark_measurements WHERE session_id=?`,
		rec.SessionID).Scan(&rec.Sequence); err != nil {
		return fmt.Errorf("store: benchmark sequence: %w", err)
	}
	_, err = tx.ExecContext(context.Background(), `INSERT INTO benchmark_measurements (
		session_id, sequence, result_id, command_json, exit_code,
		raw_stdout_bytes, raw_stderr_bytes, raw_bytes, stateless_bytes, stateful_bytes,
		presentation, truncated, execution_duration_ns, reduce_duration_ns, processing_duration_ns
	) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`,
		rec.SessionID, rec.Sequence, rec.ResultID, string(commandJSON), rec.ExitCode,
		rec.RawStdoutBytes, rec.RawStderrBytes, rec.RawBytes, rec.StatelessBytes, rec.StatefulBytes,
		rec.Presentation, boolToInt(rec.Truncated), rec.ExecutionDuration.Nanoseconds(),
		rec.ReduceDuration.Nanoseconds(), rec.ProcessingDuration.Nanoseconds())
	if err != nil {
		return fmt.Errorf("store: save benchmark measurement: %w", err)
	}
	return tx.Commit()
}

// LoadBenchmarkRecords returns a benchmark session's measurements in order.
func (s *Store) LoadBenchmarkRecords(sessionID string) ([]BenchmarkRecord, error) {
	exists, err := s.BenchmarkSessionExists(sessionID)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, sql.ErrNoRows
	}
	rows, err := s.db.QueryContext(context.Background(), `SELECT
		session_id, sequence, result_id, command_json, exit_code,
		raw_stdout_bytes, raw_stderr_bytes, raw_bytes, stateless_bytes, stateful_bytes,
		presentation, truncated, execution_duration_ns, reduce_duration_ns, processing_duration_ns
		FROM benchmark_measurements WHERE session_id=? ORDER BY sequence`, sessionID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []BenchmarkRecord
	for rows.Next() {
		var rec BenchmarkRecord
		var commandJSON string
		var truncated int
		var executionNS, reduceNS, processingNS int64
		if err := rows.Scan(&rec.SessionID, &rec.Sequence, &rec.ResultID, &commandJSON, &rec.ExitCode,
			&rec.RawStdoutBytes, &rec.RawStderrBytes, &rec.RawBytes, &rec.StatelessBytes, &rec.StatefulBytes,
			&rec.Presentation, &truncated, &executionNS, &reduceNS, &processingNS); err != nil {
			return nil, err
		}
		_ = json.Unmarshal([]byte(commandJSON), &rec.Command)
		rec.Truncated = truncated != 0
		rec.ExecutionDuration = time.Duration(executionNS)
		rec.ReduceDuration = time.Duration(reduceNS)
		rec.ProcessingDuration = time.Duration(processingNS)
		out = append(out, rec)
	}
	return out, rows.Err()
}
