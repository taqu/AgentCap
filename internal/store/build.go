package store

import (
	"context"
	"fmt"

	"github.com/taqu/agentcap/internal/buildparse"
)

func (s *Store) SaveDiagnostics(resultID string, diags []buildparse.Diagnostic) error {
	if len(diags) == 0 {
		return nil
	}
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for i, d := range diags {
		_, err := tx.ExecContext(context.Background(), `
			INSERT OR REPLACE INTO diagnostics
			(result_id, diagnostic_index, tool, severity, code, file, line, column_no, message, raw_start, raw_end)
			VALUES (?,?,?,?,?,?,?,?,?,?,?)`,
			resultID, i, d.Tool, d.Severity.String(), d.Code, d.File, d.Line, d.Col, d.Message, d.RawStart, d.RawEnd)
		if err != nil {
			return fmt.Errorf("store: save diagnostic: %w", err)
		}
	}
	return tx.Commit()
}

func (s *Store) GetDiagnostics(resultID string) ([]buildparse.Diagnostic, error) {
	rows, err := s.db.QueryContext(context.Background(), `
		SELECT tool, severity, code, file, line, column_no, message, raw_start, raw_end
		FROM diagnostics WHERE result_id = ? ORDER BY diagnostic_index`, resultID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var diags []buildparse.Diagnostic
	for rows.Next() {
		var d buildparse.Diagnostic
		var sevStr string
		if err := rows.Scan(&d.Tool, &sevStr, &d.Code, &d.File, &d.Line, &d.Col, &d.Message, &d.RawStart, &d.RawEnd); err == nil {
			d.Severity = parseSeverityStr(sevStr)
			diags = append(diags, d)
		}
	}
	return diags, nil
}

func parseSeverityStr(s string) buildparse.Severity {
	switch s {
	case "E":
		return buildparse.SeverityError
	case "W":
		return buildparse.SeverityWarning
	case "N":
		return buildparse.SeverityNote
	case "H":
		return buildparse.SeverityHelp
	case "F":
		return buildparse.SeverityFatal
	default:
		return buildparse.SeverityError
	}
}

func (s *Store) SaveTestFailures(resultID string, failures []buildparse.TestFailure) error {
	if len(failures) == 0 {
		return nil
	}
	tx, err := s.db.BeginTx(context.Background(), nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for i, f := range failures {
		_, err := tx.ExecContext(context.Background(), `
			INSERT OR REPLACE INTO test_failures
			(result_id, failure_index, suite, test_name, file, line, panic, raw_start, raw_end)
			VALUES (?,?,?,?,?,?,?,?,?)`,
			resultID, i, f.Suite, f.Name, f.File, f.Line, boolToInt(f.Panic), f.RawStart, f.RawEnd)
		if err != nil {
			return fmt.Errorf("store: save test failure: %w", err)
		}
	}
	return tx.Commit()
}

func (s *Store) GetTestFailures(resultID string) ([]buildparse.TestFailure, error) {
	rows, err := s.db.QueryContext(context.Background(), `
		SELECT suite, test_name, file, line, panic, raw_start, raw_end
		FROM test_failures WHERE result_id = ? ORDER BY failure_index`, resultID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var failures []buildparse.TestFailure
	for rows.Next() {
		var f buildparse.TestFailure
		var panicInt int
		if err := rows.Scan(&f.Suite, &f.Name, &f.File, &f.Line, &panicInt, &f.RawStart, &f.RawEnd); err == nil {
			f.Panic = panicInt != 0
			failures = append(failures, f)
		}
	}
	return failures, nil
}
