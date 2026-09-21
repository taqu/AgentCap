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
