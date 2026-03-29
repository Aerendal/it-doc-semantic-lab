# Runbook: Cut Release Candidate

## When to Use

When a stable set of semantic metadata is ready for promotion to the reference repository and a versioned snapshot is needed.

## Prerequisites

- All quality gates passed
- Export to repo 1 completed successfully
- All tests pass (`make test`)

## Steps

### 1. Run full test suite

```sh
make test
```

All tests must pass. Zero failures.

### 2. Verify evidence pack

```sh
./bin/itdlab audit evidence <run_id>
```

All required artifacts must be present.

### 3. Tag the release in this repository

```sh
git tag -a "rc/<date>-<run_id>" -m "Release candidate: <date> run <run_id>"
git push origin "rc/<date>-<run_id>"
```

### 4. Record the release in the export manifest

```sh
cat reports/<run_id>/export_manifest.json
```

Save a copy of the manifest to `reports/rc/<date>/manifest.json`.

### 5. Tag the reference repository

In `IT-Dokumentacja`:

```sh
git tag -a "semantic-rc/<date>" -m "Semantic metadata from lab rc/<date>-<run_id>"
git push origin "semantic-rc/<date>"
```

## Stop Conditions

- Any test failure → fix before cutting RC
- Gate failure → fix before cutting RC
- Evidence pack incomplete → complete it first
