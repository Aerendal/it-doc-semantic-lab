# Runbook: Export Metadata to Repo 1

## Prerequisites

- All quality gates (1–5) have passed
- Evidence pack is complete
- Run has `status = 'completed'` in SQLite
- Reference repository is checked out and on correct branch

## Steps

### 1. Verify gate status

```sh
sqlite3 db/semantic_index.sqlite \
  "SELECT run_id, status FROM runs ORDER BY started_at DESC LIMIT 1;"
```

Expected: `status = 'completed'`

### 2. Run export (dry-run first)

```sh
./bin/itdlab export repo1 --target ../IT-Dokumentacja/ --dry-run
```

Review the list of files to be promoted.

### 3. Apply export

```sh
./bin/itdlab export repo1 --target ../IT-Dokumentacja/
```

### 4. Review export manifest

```sh
cat reports/<run_id>/export_manifest.json
```

### 5. Commit in reference repository

```sh
cd ../IT-Dokumentacja
git add .
git status
git commit -m "chore: promote semantic metadata from lab run <run_id>"
```

## Stop Conditions

- Exit code 2 → quality gate failure; check `gate_failures.json`
- Files already up to date → no changes to promote
- Target directory not found → verify `--target` path
