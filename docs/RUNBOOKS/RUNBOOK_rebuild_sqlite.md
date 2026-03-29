# Runbook: Rebuild SQLite

## When to Use

- Database file is corrupt or missing
- Schema migration failed midway
- Starting fresh after major structural changes

## Prerequisites

- `db/schema_v1.sql` exists
- Event log `runs/events.jsonl` is intact (for replay if needed)

## Steps

### 1. Back up existing database (if recoverable)

```sh
cp db/semantic_index.sqlite db/semantic_index.sqlite.bak
```

### 2. Remove corrupt database

```sh
rm db/semantic_index.sqlite
```

### 3. Reinitialise schema

```sh
make db-init
```

### 4. Re-ingest from sources

```sh
./bin/itdlab ingest run --source sources/
```

### 5. Re-run normalization

```sh
./bin/itdlab normalize apply
```

### 6. Re-run relations

```sh
./bin/itdlab relations show
```

### 7. Verify row counts match backup (if available)

```sh
sqlite3 db/semantic_index.sqlite "SELECT count(*) FROM documents;"
sqlite3 db/semantic_index.sqlite.bak "SELECT count(*) FROM documents;"
```

## Notes

- The event log is the audit trail — it does not need to be rebuilt.
- If the event log is also lost, re-ingest from original source files.
- Do not commit `db/semantic_index.sqlite` to git (it is gitignored).
