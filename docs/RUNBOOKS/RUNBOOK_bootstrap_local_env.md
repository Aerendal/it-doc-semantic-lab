# Runbook: Bootstrap Local Environment

## Prerequisites

- Go 1.21+
- `make`
- `sqlite3` CLI (for `db-init` and inspection)
- Git

## Steps

### 1. Clone and enter the repository

```sh
cd /path/to/it-doc-semantic-lab
```

### 2. Tidy dependencies

```sh
make tidy
```

### 3. Build the CLI

```sh
make build
```

Binary is placed at `bin/itdlab`.

### 4. Initialise the database

```sh
make db-init
```

Creates `db/semantic_index.sqlite` with schema v1.

### 5. Create runtime directories (if not present)

```sh
mkdir -p runs reports normalized
```

### 6. Verify the build

```sh
./bin/itdlab --help
```

Expected: help text listing `ingest`, `normalize`, `classify`, `relations`, `export`, `audit`.

### 7. Run short tests

```sh
make test-short
```

Expected: all tests pass.

## Stop Conditions

- If `make tidy` fails: check internet connectivity for Go module download
- If `make build` fails: check Go version (`go version` must be 1.21+)
- If `make db-init` fails: check `sqlite3` is installed (`which sqlite3`)
