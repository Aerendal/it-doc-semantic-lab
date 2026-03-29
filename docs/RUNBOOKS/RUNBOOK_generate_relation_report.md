# Runbook: Generate Relation Report

## Prerequisites

- Normalize completed successfully
- `documents` table has rows with `status = 'normalized'`
- Relation rules exist in `relation_rules` table

## Steps

### 1. Show all relations

```sh
./bin/itdlab relations show
```

### 2. Show relations for a specific document

```sh
./bin/itdlab relations show --doc <document_id>
```

### 3. Explain a specific relation

```sh
./bin/itdlab relations explain --rel <relation_id>
```

### 4. Check relation count in database

```sh
sqlite3 db/semantic_index.sqlite "SELECT type, count(*) FROM relations GROUP BY type;"
```

### 5. Check for relations without explanation

```sh
sqlite3 db/semantic_index.sqlite \
  "SELECT id, from_id, to_id FROM relations WHERE explanation = '' OR explanation IS NULL;"
```

Expected: zero rows. (Gate 4 blocks if any exist.)

### 6. Check for dependency cycles

```sh
# Manual check for small graphs:
sqlite3 db/semantic_index.sqlite \
  "SELECT from_id, to_id FROM relations WHERE type = 'depends_on';"
```

## Stop Conditions

- Relations with empty explanation → fix inference rules; re-run
- Cycles detected in depends_on → review rules and document structure
