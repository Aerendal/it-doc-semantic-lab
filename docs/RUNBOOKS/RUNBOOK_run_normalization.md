# Runbook: Run Normalization

## Prerequisites

- Ingest completed successfully
- `documents` table has rows with `status = 'ingested'`

## Steps

### 1. Preview normalization

```sh
./bin/itdlab normalize preview
```

Review output — check for unexpected canonical ID changes or collisions.

### 2. Apply normalization

```sh
./bin/itdlab normalize apply
```

### 3. Verify canonical IDs

```sh
sqlite3 db/semantic_index.sqlite "SELECT id, canonical_id, raw_name FROM documents;"
```

### 4. Check for collisions

```sh
sqlite3 db/semantic_index.sqlite \
  "SELECT canonical_id, count(*) as n FROM documents GROUP BY canonical_id HAVING n > 1;"
```

Expected: zero rows. If collisions exist, resolve per PLAYBOOK_normalization.md before proceeding.

### 5. Review normalization report

```sh
cat reports/<run_id>/normalization_report.json
```

## Stop Conditions

- Collision count > 0 → resolve before export (Gate 3)
- Unexpected canonical ID changes → re-check normalization rules
