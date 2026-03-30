# Runbook: Authoritative Verify

**Version:** 1.0  
**Trigger:** Manual or scheduled authoritative verification of the repository's assurance state  
**Scope:** Full repo — all capability slices present at time of execution

---

## Purpose

This runbook describes how to perform a full authoritative verification of the `itdlab` repository. It answers the question:

> "At this point in time, what is the actual, evidenced state of quality in this repository?"

Authoritative verify is not a development check. It is the formal procedure used before:
- a promotion decision to the stable repository,
- a gate status report for G4 or G5,
- a milestone sign-off,
- an external audit or review.

---

## Prerequisites

Before starting, confirm all of the following:

- [ ] Latest code is committed and pushed to `origin/main`
- [ ] `git status` shows no uncommitted changes
- [ ] `bin/itdlab` is built from the current HEAD (`make build`)
- [ ] SQLite is at the latest schema version (`make db-init` or verify `schema_version = 1`)
- [ ] No other `itdlab` run is in progress (check `runs/events.jsonl` for an open run)
- [ ] Disk space available: ≥ 500 MB free on the partition hosting `reports/`
- [ ] The operator has reviewed `docs/SKIP_REGISTER.md` for any active skips

---

## Procedure

### Phase 1: Environment snapshot

```bash
# Record binary version
./bin/itdlab version  # or ./bin/itdlab --help | head -1

# Record schema version
sqlite3 db/semantic_index.sqlite "SELECT version, applied_at FROM schema_version;"

# Record event log baseline (line count before this run)
wc -l runs/events.jsonl

# Record git HEAD
git rev-parse HEAD
git log --oneline -1

# Write all of the above to a snapshot file
{
  echo "# Authoritative Verify Snapshot"
  echo "Date: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "Git HEAD: $(git rev-parse HEAD)"
  echo "Git message: $(git log --oneline -1)"
  echo "Binary: $(./bin/itdlab --version 2>/dev/null || echo 'see help')"
  sqlite3 db/semantic_index.sqlite "SELECT 'Schema version: ' || version FROM schema_version;"
  echo "Event log lines before: $(wc -l < runs/events.jsonl)"
} > /tmp/av_snapshot.txt
cat /tmp/av_snapshot.txt
```

### Phase 2: Re-ingest from sources

A clean ingest verifies that the current sources can be fully ingested without error.

```bash
./bin/itdlab ingest run --source sources/
```

Stop if exit code ≠ 0. Record the `run_id` from stdout.

```bash
AV_INGEST_RUN_ID=<run_id from above>
```

### Phase 3: Normalize

```bash
./bin/itdlab normalize apply
```

Stop if exit code ≠ 0. Record the `run_id`.

```bash
AV_NORMALIZE_RUN_ID=<run_id from above>
```

### Phase 4: Check evidence completeness for each run

```bash
./bin/itdlab audit evidence $AV_INGEST_RUN_ID
./bin/itdlab audit evidence $AV_NORMALIZE_RUN_ID
```

Expected: both report `trusted = true` and `evidence.complete = true`.

If either reports `trusted = false`, stop and investigate before proceeding.

### Phase 5: Run all tests

```bash
make test 2>&1 | tee /tmp/av_test_output.txt
```

Record exit code. A non-zero exit code means at least one test failed. Do not proceed to gate evaluation if tests fail.

### Phase 6: Evaluate quality gates

For each gate (G1 through G5, as applicable to completed capability slices):

```bash
# Check gate status — currently manual; see docs/QUALITY_GATES.md
# For each gate, confirm:
# - all required layers were executed (no unregistered skips)
# - all layers produced trusted runs
# - evidence artifacts exist
```

Summarise gate status:

| Gate | Status | Notes |
|------|--------|-------|
| G1 | PASS / degraded / FAIL / not-evaluated | |
| G2 | PASS / degraded / FAIL / not-evaluated | |
| G3 | PASS / degraded / FAIL / not-evaluated | |
| G4 | PASS / degraded / FAIL / not-evaluated | |
| G5 | PASS / degraded / FAIL / not-evaluated | |

### Phase 7: Check active skips

```bash
cat docs/SKIP_REGISTER.md
```

For each active skip:
- Confirm it has a valid Skip ID, owner, and review date.
- Confirm it is reflected in the gate status table above.
- Confirm no Category 3 or Category 4 skip is present without approval record.

### Phase 8: Write authoritative verify report

```bash
REPORT_DIR="reports/av-$(date -u +%Y%m%d-%H%M%S)"
mkdir -p $REPORT_DIR

# Copy snapshot
cp /tmp/av_snapshot.txt $REPORT_DIR/

# Copy test output
cp /tmp/av_test_output.txt $REPORT_DIR/

# Write gate status
cat > $REPORT_DIR/gate_status.md << 'EOF'
# Gate Status

Date: <ISO 8601>
Operator: <name>
Git HEAD: <sha>

## Summary

| Gate | Status | Notes |
|------|--------|-------|
| G1   |        |       |
| G2   |        |       |
| G3   |        |       |
| G4   |        |       |
| G5   |        |       |

## Active skips
(copy from SKIP_REGISTER.md)

## Trusted run IDs cited
(list all run_ids used as evidence)
EOF

echo "Authoritative verify report written to: $REPORT_DIR"
```

---

## Stop Conditions

Stop and record as INCOMPLETE if:

- Any capability-slice run exits non-zero without a registered explanation.
- Any run cited as gate evidence is `trusted = false`.
- `make test` exits non-zero.
- Any `mock-forbidden` layer was found to have been executed with mocks.
- `docs/SKIP_REGISTER.md` contains an unowned or expired skip.

---

## Completion Criteria

The authoritative verify is complete and credible when:

- [ ] All runs in Phase 2–3 exited 0 or 2 (no unhandled errors)
- [ ] All cited runs are `trusted = true`
- [ ] All tests passed or all failures are registered in SKIP_REGISTER.md
- [ ] Gate status table is filled in with supporting run IDs
- [ ] Report exists in `reports/av-<timestamp>/`
- [ ] Report is committed: `git add reports/av-<timestamp>/ && git commit -m "audit: authoritative verify <date>"`

---

## Internal references
- `docs/EXECUTION_ASSURANCE_PROGRAM.md`
- `docs/QUALITY_GATES.md`
- `docs/QUALITY_GATES_POLICY.md`
- `docs/EVIDENCE_MODEL.md`
- `docs/SKIP_REGISTER.md`
- `docs/PLAYBOOKS/PLAYBOOK_REAL_PATH_VERIFICATION.md`

## Review metadata
- Owner: project team
- Status: draft
- Last reviewed: 2026-03-30
