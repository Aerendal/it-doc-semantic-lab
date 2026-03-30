# ADR-002: Event Log as Audit Backbone

| Field | Value |
|-------|-------|
| Status | Accepted |
| Date | 2026-03-30 |
| Deciders | project team |
| Supersedes | — |
| Superseded by | — |

---

## Context

The system processes documents, infers relations, and normalises names through a series of runs. Each run modifies SQLite state. Without an independent record of what happened, when, and in what order, there is no way to:
- reconstruct the history of how a document reached its current state,
- compare two runs to identify regressions,
- audit a run after the fact without re-running it,
- detect whether SQLite state was modified outside a managed run.

SQLite is the source of truth for current state, but it does not natively track historical state or cross-run provenance. Something else must serve as the audit backbone.

---

## Decision

Every state-changing operation in `itdlab` appends a structured JSON event to an append-only event log file (`runs/events.jsonl`). The event log is:

1. **Append-only** — existing lines are never modified or deleted.
2. **Structured** — each line is a valid JSON object with mandatory fields.
3. **Run-scoped** — every event carries a `run_id` that ties it to a specific invocation.
4. **Independent of SQLite** — the log is written to a separate file; it does not depend on SQLite transactions.

### Mandatory event fields

```json
{
  "ts": "<ISO 8601 UTC>",
  "run_id": "<string>",
  "step": "<string>",
  "entity": "<string>",
  "entity_id": "<string>",
  "action": "<string>",
  "before": "<any | null>",
  "after": "<any | null>"
}
```

Additional fields are permitted. The mandatory set must always be present.

---

## Alternatives Considered

### 1. SQLite-only history tables

**Approach:** Add `_history` shadow tables that capture previous values on UPDATE/DELETE.

**Rejected because:**
- Triggers in SQLite are brittle and hard to test.
- History tables in the same file as current state create a single point of failure.
- Reading run history requires SQL joins across many tables; a flat log is easier to stream and grep.
- History tables do not survive if the database is rebuilt from scratch.

### 2. Write-ahead log (WAL) as audit record

**Approach:** Rely on SQLite's WAL file as the event record.

**Rejected because:**
- WAL files are checkpointed and overwritten; they are not a permanent audit trail.
- WAL format is binary; it is not human-readable without tooling.
- WAL is implementation-specific and not part of the public API.

### 3. Structured logging only (stdout)

**Approach:** Emit structured log lines to stdout; redirect to file by the caller.

**Rejected because:**
- Stdout capture depends on the shell invocation; it can be lost.
- Stdout mixes diagnostic output with audit events; filtering is fragile.
- Stdout is not append-only by design; it can be truncated.

---

## Consequences

### Positive

- Complete history of all state changes is available for any run, including runs that failed or were aborted.
- A run can be reconstructed from the event log without access to the current SQLite state.
- Two runs can be diffed by comparing their event log slices.
- The event log can be replayed to detect if SQLite diverges from the expected state.
- Forensic review after a contract breach (exit code 3) is possible even if SQLite is corrupted.

### Negative / Accepted tradeoffs

- The event log grows indefinitely. There is currently no compaction policy.
- Appending to JSONL on every event is a sequential I/O operation; it is a bottleneck for high-volume runs.
- The event log and SQLite can temporarily diverge between a SQLite write and the subsequent log append. This is accepted; the invariant (I1 in `EXECUTION_CONTRACT.md`) requires SQLite write before log append, not atomicity of both.

### Deferred

- Cross-run replay tooling (`itdlab audit replay`) is not implemented in v1.
- Compaction / archival policy is deferred.
- Event log signing / tamper detection is deferred.

---

## Implementation Notes

- Event log is written by `internal/adapters/jsonl/event_log.go`.
- Mutex-protected append; safe for concurrent goroutines within a single run.
- Each run produces a separate slice of events identifiable by `run_id`.
- `itdlab audit runs` reads the event log to list historical runs.

---

## Internal references
- `docs/EXECUTION_ASSURANCE_PROGRAM.md`
- `docs/EXECUTION_CONTRACT.md`
- `docs/EVIDENCE_MODEL.md`
- `docs/ADR/ADR-001-sqlite-as-source-of-truth.md`

## Review metadata
- Owner: project team
- Status: accepted
- Last reviewed: 2026-03-30
