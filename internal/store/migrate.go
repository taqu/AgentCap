package store

import (
	"context"
	"database/sql"
	"fmt"
)

// runMigrations applies sequential migrations from `from` to `to`.
func runMigrations(db *sql.DB, from, to int) error {
	for v := from; v < to; v++ {
		fn, ok := migrations[v]
		if !ok {
			return fmt.Errorf("store: no migration from version %d", v)
		}
		tx, err := db.BeginTx(context.Background(), nil)
		if err != nil {
			return fmt.Errorf("store: migration %d begin: %w", v, err)
		}
		if err := fn(tx); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("store: migration %d: %w", v, err)
		}
		if _, err := tx.ExecContext(context.Background(),
			`UPDATE schema_info SET version=?`, v+1); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("store: migration %d update version: %w", v, err)
		}
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("store: migration %d commit: %w", v, err)
		}
	}
	return nil
}

// migrations maps from-version to the function that upgrades to version+1.
var migrations = map[int]func(*sql.Tx) error{
	1: migrateV1ToV2,
	2: migrateV2ToV3,
	3: migrateV3ToV4,
	4: migrateV4ToV5,
}

func migrateV4ToV5(tx *sql.Tx) error {
	_, err := tx.ExecContext(context.Background(), `ALTER TABLE results ADD COLUMN integration_json TEXT NOT NULL DEFAULT ''`)
	return err
}

func migrateV3ToV4(tx *sql.Tx) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS benchmark_sessions (
    id         TEXT PRIMARY KEY,
    created_at TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS benchmark_measurements (
    session_id             TEXT NOT NULL,
    sequence               INTEGER NOT NULL,
    result_id              TEXT NOT NULL,
    command_json           TEXT NOT NULL DEFAULT '[]',
    exit_code              INTEGER NOT NULL DEFAULT 0,
    raw_stdout_bytes       INTEGER NOT NULL DEFAULT 0,
    raw_stderr_bytes       INTEGER NOT NULL DEFAULT 0,
    raw_bytes              INTEGER NOT NULL DEFAULT 0,
    stateless_bytes        INTEGER NOT NULL DEFAULT 0,
    stateful_bytes         INTEGER NOT NULL DEFAULT 0,
    presentation           TEXT NOT NULL DEFAULT 'full',
    truncated              INTEGER NOT NULL DEFAULT 0,
    execution_duration_ns  INTEGER NOT NULL DEFAULT 0,
    reduce_duration_ns     INTEGER NOT NULL DEFAULT 0,
    processing_duration_ns INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (session_id, sequence),
    UNIQUE (result_id),
    FOREIGN KEY (session_id) REFERENCES benchmark_sessions(id)
);`
	_, err := tx.ExecContext(context.Background(), ddl)
	return err
}

func migrateV2ToV3(tx *sql.Tx) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS diagnostics (
    result_id        TEXT NOT NULL,
    diagnostic_index INTEGER NOT NULL,
    tool             TEXT NOT NULL DEFAULT '',
    severity         TEXT NOT NULL DEFAULT 'E',
    code             TEXT NOT NULL DEFAULT '',
    file             TEXT NOT NULL DEFAULT '',
    line             INTEGER NOT NULL DEFAULT 0,
    column_no        INTEGER NOT NULL DEFAULT 0,
    message          TEXT NOT NULL DEFAULT '',
    raw_start        INTEGER NOT NULL DEFAULT 0,
    raw_end          INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (result_id, diagnostic_index)
);

CREATE TABLE IF NOT EXISTS test_failures (
    result_id     TEXT NOT NULL,
    failure_index INTEGER NOT NULL,
    suite         TEXT NOT NULL DEFAULT '',
    test_name     TEXT NOT NULL DEFAULT '',
    file          TEXT NOT NULL DEFAULT '',
    line          INTEGER NOT NULL DEFAULT 0,
    panic         INTEGER NOT NULL DEFAULT 0,
    raw_start     INTEGER NOT NULL DEFAULT 0,
    raw_end       INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (result_id, failure_index)
);`
	_, err := tx.ExecContext(context.Background(), ddl)
	return err
}

func migrateV1ToV2(tx *sql.Tx) error {
	const ddl = `
CREATE TABLE IF NOT EXISTS git_diff_files (
    result_id   TEXT NOT NULL,
    file_index  INTEGER NOT NULL,
    old_path    TEXT NOT NULL DEFAULT '',
    new_path    TEXT NOT NULL DEFAULT '',
    status      TEXT NOT NULL DEFAULT '',
    additions   INTEGER NOT NULL DEFAULT 0,
    deletions   INTEGER NOT NULL DEFAULT 0,
    binary      INTEGER NOT NULL DEFAULT 0,
    old_mode    TEXT NOT NULL DEFAULT '',
    new_mode    TEXT NOT NULL DEFAULT '',
    similarity  INTEGER NOT NULL DEFAULT 0,
    raw_start   INTEGER NOT NULL DEFAULT 0,
    raw_end     INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (result_id, file_index)
);

CREATE TABLE IF NOT EXISTS git_diff_hunks (
    result_id   TEXT NOT NULL,
    file_index  INTEGER NOT NULL,
    hunk_index  INTEGER NOT NULL,
    old_start   INTEGER NOT NULL DEFAULT 0,
    old_lines   INTEGER NOT NULL DEFAULT 0,
    new_start   INTEGER NOT NULL DEFAULT 0,
    new_lines   INTEGER NOT NULL DEFAULT 0,
    header      TEXT NOT NULL DEFAULT '',
    raw_start   INTEGER NOT NULL DEFAULT 0,
    raw_end     INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (result_id, file_index, hunk_index)
);`
	_, err := tx.ExecContext(context.Background(), ddl)
	return err
}
