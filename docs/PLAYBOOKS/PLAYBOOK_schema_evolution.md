# Playbook: Schema Evolution

## Purpose

Defines the strategy for evolving the SQLite schema without breaking existing data or runs.

---

## Principles

1. Schema changes are **always additive** — no column drops, no renames in place
2. Each schema version has its own DDL file: `db/schema_v<N>.sql`
3. The `schema_version` table tracks applied versions
4. Migration tests (Layer 15) must pass before applying any migration

---

## Adding a New Column

```sql
-- db/schema_v2.sql (example)
ALTER TABLE documents ADD COLUMN confidence REAL NOT NULL DEFAULT 0.0;
INSERT OR IGNORE INTO schema_version (version, applied_at, description)
VALUES (2, datetime('now'), 'add confidence column to documents');
```

Rules:
- New columns must have a `DEFAULT` value
- Never `NOT NULL` without a default on an existing table
- Add an index if the column will be queried

---

## Adding a New Table

Add the new `CREATE TABLE IF NOT EXISTS` statement to the new schema file. New tables do not require data migration.

---

## Migration Process

1. Write `db/schema_v<N>.sql`
2. Write a migration test (Layer 15) using a pre-migration SQLite snapshot
3. Run `make test` — verify migration test passes
4. Apply: `sqlite3 db/semantic_index.sqlite < db/schema_v<N>.sql`
5. Update `internal/adapters/sqlite/schema.go` to include the new DDL

---

## What Never to Do

- Do not `DROP COLUMN` or `DROP TABLE` on existing schema versions
- Do not rename columns (use a new column + deprecated old one)
- Do not change `CHECK` constraints on existing columns in place
