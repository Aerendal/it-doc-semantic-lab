# Playbook: Real-Path Verification

**Version:** 1.0  
**Scope:** All capability slices — applicable whenever a `mock-restricted` or `mock-forbidden` test layer is being executed

---

## Purpose

This playbook describes *how* to perform real-path verification: what it means, when to do it, what counts as passing, and what to do when it fails.

Real-path verification is not the same as running tests. It is the deliberate confirmation that the code path responsible for a test-covered behaviour was exercised against real, non-mocked dependencies.

---

## When to Apply This Playbook

Apply when any of the following is true:

- You are preparing to claim a `mock-restricted` test layer as PASS.
- You are preparing to claim a `mock-forbidden` test layer as PASS.
- A reviewer has asked for real-path verification evidence.
- You are preparing a promotion-readiness package for G4 or G5.
- A run manifest shows `trusted = false` and the reason is a mock policy violation.

---

## What Real-Path Verification Covers

Per `docs/POLICY_MOCKS_AND_REAL_PATHS.md`, the following layers require real-path verification:

| Layer | Name | Mock policy |
|-------|------|-------------|
| 5 | File system contract tests | mock-restricted |
| 7 | Golden output regression tests | mock-restricted |
| 8 | Determinism tests | mock-restricted |
| 22 | End-to-end slice tests | mock-forbidden |
| 23 | Database integrity tests | mock-forbidden |
| 25 | Run reproducibility tests | mock-forbidden |
| 26 | Evidence completeness tests | mock-forbidden |

For `mock-restricted` layers: at least one full test execution per slice must use real files with no mocking of file I/O or SQLite.

For `mock-forbidden` layers: every execution must use real files and real SQLite. No exceptions.

---

## Step-by-Step Procedure

### Step 1: Identify the target layer and slice

```
Layer: <e.g. Layer 22 — End-to-end slice tests>
Slice: <e.g. Slice 1 — Ingest>
Run ID (if verifying a previous run): <run_id>
```

### Step 2: Confirm no mocks are active for the layer

Check the test files for the target layer:

```bash
grep -r "t.Skip\|mock\|Mock\|stub\|Stub\|FakeDB\|InMemory" \
  <test_file_or_directory>
```

Expected: No mock setup for the components being tested. If mocks are found:
- Check if they are for *external* dependencies only (acceptable for `mock-restricted`).
- If they wrap SQLite, the file system, or the event log → this is a violation for `mock-forbidden` layers.
- Do not proceed until mocks are removed or the layer is reclassified.

### Step 3: Prepare the real-path test environment

```bash
# Create a fresh test directory
mkdir -p /tmp/itdlab-realpath-test
cd /tmp/itdlab-realpath-test

# Copy real source fixtures
cp -r <repo>/sources/  ./sources/

# Initialise a fresh database (no pre-existing state)
<repo>/bin/itdlab --db ./semantic_index.sqlite --log ./events.jsonl ingest run --source ./sources/
```

Confirm:
- `semantic_index.sqlite` is a real file on disk (not `:memory:`).
- `events.jsonl` is a real file on disk.
- Sources are real files, not generated stubs.

### Step 4: Run the target layer's tests against real dependencies

```bash
# From repo root
DB_PATH=/tmp/itdlab-realpath-test/semantic_index.sqlite \
LOG_PATH=/tmp/itdlab-realpath-test/events.jsonl \
go test ./... -run <TestPattern> -v 2>&1 | tee /tmp/itdlab-realpath-test/test_output.txt
```

### Step 5: Verify the evidence artifacts

After the test run:

```bash
# Check that SQLite was actually used (non-empty, not in-memory artefact)
ls -lh /tmp/itdlab-realpath-test/semantic_index.sqlite

# Check that events were appended
wc -l /tmp/itdlab-realpath-test/events.jsonl

# Check that reports were produced
ls -lh reports/<run_id>/
```

All of the following must be true:
- [ ] SQLite file exists and is larger than the empty-schema baseline
- [ ] At least 1 event line in the event log
- [ ] `run_manifest.json` exists and `evidence.complete = true`
- [ ] `run_manifest.json` has `trusted = true`

### Step 6: Record the real-path evidence

In `reports/<run_id>/`:
```
real_path_verification.md  — document: layer, slice, test command, result, timestamp
```

Template:
```md
## Real-Path Verification Record

- Layer: <layer name>
- Slice: <slice name>
- Run ID: <run_id>
- Test command: `<full command>`
- SQLite path: <path>
- Event log path: <path>
- Result: PASS / FAIL
- Verified at: <ISO 8601 datetime>
- Verified by: <owner>
- Notes: <any deviations or observations>
```

---

## Stop Conditions

Stop and do not record a PASS if any of the following is true:

- Any `mock-forbidden` component was wrapped in a mock or stub.
- SQLite was opened as `:memory:`.
- Source files were generated dynamically without review.
- The test environment was shared with a running development server that could modify state.
- `run_manifest.json` shows `trusted = false`.

---

## Recovery if Real-Path Verification Fails

1. Record the failure in the run manifest (`evidence.complete = false` if artifacts are missing, or `trusted = false` if a policy violation was found).
2. Register a skip in `docs/SKIP_REGISTER.md` if the layer cannot be executed for a documented reason.
3. Do not proceed to G4/G5 gate evaluation until the layer is passing or a Category 1/2 skip is registered with owner and review date.

---

## Internal references
- `docs/POLICY_MOCKS_AND_REAL_PATHS.md`
- `docs/TESTING_STANDARD.md`
- `docs/TEST_CATALOG.md`
- `docs/EVIDENCE_MODEL.md`
- `docs/SKIP_REGISTER.md`

## Review metadata
- Owner: project team
- Status: draft
- Last reviewed: 2026-03-30
