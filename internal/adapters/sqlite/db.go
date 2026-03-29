package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

// DB wraps a SQLite connection with migration support.
type DB struct {
	conn *sql.DB
	path string
}

// Open opens (or creates) the SQLite database at path and applies the embedded schema.
func Open(path string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, fmt.Errorf("sqlite: create db dir: %w", err)
	}

	conn, err := sql.Open("sqlite", path+"?_pragma=journal_mode(WAL)&_pragma=foreign_keys(ON)")
	if err != nil {
		return nil, fmt.Errorf("sqlite: open: %w", err)
	}

	db := &DB{conn: conn, path: path}
	if err := db.migrate(context.Background()); err != nil {
		conn.Close()
		return nil, fmt.Errorf("sqlite: migrate: %w", err)
	}

	return db, nil
}

// Close closes the underlying database connection.
func (db *DB) Close() error {
	return db.conn.Close()
}

// Conn returns the underlying *sql.DB for use by store implementations.
func (db *DB) Conn() *sql.DB {
	return db.conn
}

// migrate applies schema_v1.sql if the schema_version table does not yet exist.
func (db *DB) migrate(ctx context.Context) error {
	_, err := db.conn.ExecContext(ctx, schemaV1)
	if err != nil {
		return fmt.Errorf("apply schema_v1: %w", err)
	}
	return nil
}
