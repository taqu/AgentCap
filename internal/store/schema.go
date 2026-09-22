package store

import (
	"context"
	"fmt"
)

const currentSchemaVersion = 3

const schemaV1 = `
CREATE TABLE IF NOT EXISTS schema_info (
	version INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS results (
	id           TEXT PRIMARY KEY,
	created_at   TEXT NOT NULL,
	started_at   TEXT NOT NULL DEFAULT '',
	cwd          TEXT NOT NULL DEFAULT '',
	argv_json    TEXT NOT NULL DEFAULT '[]',
	exit_code    INTEGER NOT NULL DEFAULT 0,
	duration_ms  INTEGER NOT NULL DEFAULT 0,
	truncated    INTEGER NOT NULL DEFAULT 0,
	reducer      TEXT NOT NULL DEFAULT '',
	stdout_object TEXT NOT NULL DEFAULT '',
	stderr_object TEXT NOT NULL DEFAULT '',
	stdout_bytes INTEGER NOT NULL DEFAULT 0,
	stderr_bytes INTEGER NOT NULL DEFAULT 0,
	stdout_hash  TEXT NOT NULL DEFAULT '',
	stderr_hash  TEXT NOT NULL DEFAULT '',
	capsule      TEXT NOT NULL DEFAULT '',
	session_id   TEXT NOT NULL DEFAULT '',
	sequence     INTEGER NOT NULL DEFAULT 0,
	baseline_id  TEXT NOT NULL DEFAULT '',
	presentation TEXT NOT NULL DEFAULT '',
	work_dir     TEXT NOT NULL DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_results_created_at ON results(created_at);

CREATE TABLE IF NOT EXISTS stats (
	key   TEXT PRIMARY KEY,
	value INTEGER NOT NULL DEFAULT 0
);
`

func (s *Store) initSchema() error {
	// Check if schema_info exists and has a version.
	// On ANY error (table not found, no rows, etc.) treat as new database.
	var version int
	err := s.db.QueryRowContext(context.Background(),
		`SELECT version FROM schema_info LIMIT 1`).Scan(&version)

	if err != nil {
		// New database (or first open): create schema.
		if _, err := s.db.ExecContext(context.Background(), schemaV1); err != nil {
			return fmt.Errorf("create schema: %w", err)
		}
		// Insert initial schema version, then run migrations to current.
		var count int
		_ = s.db.QueryRowContext(context.Background(),
			`SELECT COUNT(*) FROM schema_info`).Scan(&count)
		if count == 0 {
			if _, err := s.db.ExecContext(context.Background(),
				`INSERT INTO schema_info(version) VALUES(?)`, 1); err != nil {
				return fmt.Errorf("set schema version: %w", err)
			}
		}
		return runMigrations(s.db, 1, currentSchemaVersion)
	}

	// Check version compatibility.
	if version > currentSchemaVersion {
		return fmt.Errorf("store: schema version %d is newer than supported version %d; please upgrade acap",
			version, currentSchemaVersion)
	}

	// Future: run migrations from version to currentSchemaVersion.
	return runMigrations(s.db, version, currentSchemaVersion)
}
