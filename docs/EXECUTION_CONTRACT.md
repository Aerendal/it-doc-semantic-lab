# Execution Contract

This document defines what every `itdlab` run must guarantee — the preconditions it requires, the invariants it must maintain, and the postconditions it must satisfy. A run that violates any invariant or postcondition is **not a valid run** and must not be cited as evidence.

---

## Purpose

The execution contract is the binding agreement between the CLI tool and its operators. It does not describe *how* the tool works — it describes *what* is guaranteed to be true about every run that the tool accepts responsibility for.

---

## Preconditions

Before a run begins, the following must be true. If any precondition fails, the run must refuse to start, exit with code 1, and write a descriptive error to stderr.

| # | Precondition | Check |
|---|-------------|-------|
| P1 | SQLite database file is accessible (exists or creatable) | Checked at startup |
| P2 | SQLite schema version ≥ 1 (or can be initialised) | Checked at startup |
| P3 | Event log path is writable (file exists or can be created) | Checked at startup |
| P4 | Reports directory is writable | Checked at startup |
| P5 | A unique `run_id` can be generated | Checked at startup |
| P6 | No other run with the same `run_id` exists in SQLite | Checked before first write |

A run that cannot satisfy P1–P6 is a **refused run**. Refused runs do not write to SQLite or the event log (except for the refusal reason to stderr).

---

## Invariants

The following must remain true throughout the run — from the first write to the last. Violation of any invariant is a **contract breach** and must result in immediate controlled abort.

| # | Invariant | Consequence of violation |
|---|-----------|-------------------------|
| I1 | Every state-changing operation is written to SQLite before it is appended to the event log | Abort; mark run `status = 'aborted'` |
| I2 | The event log is append-only; no existing line is modified or deleted | Abort; do not truncate log |
| I3 | The `run_id` in all event log entries matches the current run | Abort if mismatch detected |
| I4 | Foreign key constraints in SQLite are never bypassed | Abort; roll back transaction |
| I5 | No partial write leaves SQLite in an inconsistent state | Use transactions for all multi-step writes |
| I6 | The run's `status` field in SQLite transitions only forward: `running` → `completed` / `failed` / `aborted` | Abort if backward transition attempted |

---

## Postconditions

After a run completes (exit 0 or exit 1), all of the following must be true. If any postcondition cannot be satisfied, the run must record `status = 'failed'` and exit non-zero.

| # | Postcondition | Verification |
|---|--------------|--------------|
| Q1 | Run record exists in SQLite with `finished_at` and `exit_code` set | `SELECT * FROM runs WHERE run_id = ?` |
| Q2 | At least one event log entry exists for the run | `grep run_id events.jsonl` |
| Q3 | `reports/<run_id>/stdout.txt` exists and is non-empty | File check |
| Q4 | `reports/<run_id>/db_checksum.txt` exists and contains a valid SHA-256 line | File check + format check |
| Q5 | `reports/<run_id>/summary.md` exists and is non-empty | File check |
| Q6 | All command-specific required artifacts exist (see `docs/EVIDENCE_MODEL.md`) | Per-command check |
| Q7 | SQLite WAL is checkpointed before checksum is computed | `PRAGMA wal_checkpoint(FULL)` before hash |

A run that exits 0 but fails Q1–Q7 is an **INCOMPLETE run** per `docs/EVIDENCE_MODEL.md` and may not be cited as gate evidence.

---

## Contract Breach Protocol

When a contract breach is detected (invariant violation mid-run):

1. Stop all further writes immediately.
2. Attempt to write a final `action: "contract_breach"` event to the event log with the invariant ID and reason.
3. Set run `status = 'aborted'` in SQLite if the DB is still writable.
4. Exit with code 3 (contract breach, distinct from error=1 and gate failure=2).
5. Write the breach reason to stderr.
6. Do **not** clean up partial state — preserve it for forensic review.

---

## Exit Code Contract

| Code | Meaning |
|------|---------|
| 0 | Run completed; all postconditions satisfied |
| 1 | Run failed; at least one postcondition not satisfied; evidence pack may be incomplete |
| 2 | Gate failure; postconditions satisfied but at least one quality gate did not pass |
| 3 | Contract breach; invariant violated mid-run; partial state preserved |

Exit codes are part of the CLI contract. Any code that produces a different exit code for these conditions has a defect.

---

## Scope

This contract applies to:
- All `itdlab` commands that modify SQLite state
- All `itdlab` commands that append to the event log
- All `itdlab` commands that produce evidence artifacts

It does not apply to:
- `itdlab --help` and `itdlab [command] --help`
- `itdlab ingest inspect` (read-only inspection, no state change)

Read-only commands must still write to stdout but are not required to produce evidence packs.

---

## Internal references
- `docs/EXECUTION_ASSURANCE_PROGRAM.md`
- `docs/QUALITY_GATES_POLICY.md`
- `docs/TESTING_STANDARD.md`
- `docs/EVIDENCE_MODEL.md`
- `docs/CONTEXT_VOCABULARY.md`

## Authority anchors
- `docs/REFERENCES.md` — RFC 2119 (MUST / SHOULD semantics)

## Review metadata
- Owner: project team
- Status: draft
- Last reviewed: 2026-03-30
