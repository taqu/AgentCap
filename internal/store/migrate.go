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
	// version 1->2 will be added when Phase 3 requires schema changes.
}
