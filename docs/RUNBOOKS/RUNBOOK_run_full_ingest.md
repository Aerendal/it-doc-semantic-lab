# Runbook: Run Full Ingest

## Prerequisites

- `bin/itdlab` built (`make build`)
- `db/semantic_index.sqlite` initialised (`make db-init`)
- Source files present in `sources/`

## Steps

### 1. Verify sources exist

```sh
ls sources/
```

Expected: at least one `.md` file in a subdirectory.

### 2. Run ingest

```sh
./bin/itdlab ingest run --source sources/
```

### 3. Verify event log entries

```sh
grep '"step":"ingest"' runs/events.jsonl | wc -l
```

Expected: count equals number of ingested files.

### 4. Inspect a document

```sh
./bin/itdlab ingest inspect sources/<path-to-file>.md
```

### 5. Check database

```sh
sqlite3 db/semantic_index.sqlite "SELECT count(*) FROM documents;"
```

Expected: non-zero row count.

### 6. Review parse report

```sh
cat reports/<run_id>/parse_report.json
```

## Stop Conditions

- Exit code non-zero → check `reports/<run_id>/stdout.txt` and event log
- Zero rows in `documents` table → no `.md` files found; check source path
- Gate 1 failure → fix source files per `gate_failures.json` and re-run
