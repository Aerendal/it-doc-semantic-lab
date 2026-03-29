# Evidence Model

Every completed run of `itdlab` must produce a verifiable **evidence pack** — a set of artifacts that allow independent audit of what happened, why, and with what result.

A run without a complete evidence pack is treated as an **INCOMPLETE run** regardless of exit code.

---

## Principles

1. Evidence must be produced by the run itself — not manually assembled after the fact.
2. Evidence must be sufficient for an independent reviewer to reconstruct what happened without talking to the author.
3. Evidence must survive the session — artifacts written only to memory or temporary directories do not count.
4. A green exit code without evidence is not a credible result.

---

## Required Artifacts per Run

Every run of any `itdlab` command must produce **all** of the following:

| # | Artifact | Location | Format | Description |
|---|----------|----------|--------|-------------|
| 1 | Run record | SQLite `runs` table | row | start time, finish time, exit code, status |
| 2 | Event log entries | `runs/events.jsonl` | JSONL | All events appended during this run |
| 3 | Command output | `reports/<run_id>/stdout.txt` | text | Captured CLI stdout |
| 4 | Exit code | Run record in SQLite | integer | 0 = success, non-zero = failure |
| 5 | SQLite checksum | `reports/<run_id>/db_checksum.txt` | text | SHA-256 of DB file after run |
| 6 | Run summary | `reports/<run_id>/summary.md` | Markdown | Human-readable run summary |

Artifacts 1–6 are mandatory for **all** runs.

---

## Command-Specific Required Artifacts

In addition to the universal set, each command requires its own artifacts:

| Command | Required Artifact | Location |
|---------|-------------------|----------|
| `ingest run` | Source manifest | `reports/<run_id>/source_manifest.json` |
| `ingest run` | Parse report | `reports/<run_id>/parse_report.json` |
| `normalize apply` | Normalization report | `reports/<run_id>/normalization_report.json` |
| `normalize apply` | Collision report | `reports/<run_id>/collision_report.json` |
| `relations show` | Relation graph | `reports/<run_id>/relation_graph.json` |
| `authority check` | Authority coverage report | `reports/<run_id>/authority_coverage_report.json` |
| `export repo1` | Export manifest | `reports/<run_id>/export_manifest.json` |
| `export repo1` | Gate pass report | `reports/<run_id>/gate_pass_report.json` |
| `audit evidence` | Evidence manifest | `reports/<run_id>/evidence_manifest.json` |

---

## INCOMPLETE Run Definition

A run is declared **INCOMPLETE** if any of the following is true:

1. The universal artifact set (items 1–6) is not fully present.
2. A command-specific required artifact is missing.
3. Any artifact exists but is empty (zero bytes).
4. The SQLite `runs` table has no row for this `run_id`.
5. The event log has no entries for this `run_id`.
6. The run record has `status = 'running'` and the process is no longer active.

An INCOMPLETE run may not be cited as evidence for gate evaluation.

---

## Evidence Pack Completeness Check

Evidence pack completeness is verified by:
- **Layer 27** (Evidence Pack Tests) in `docs/TEST_CATALOG.md`
- **Gate 5** in `docs/QUALITY_GATES.md`

The `itdlab audit evidence <run_id>` command performs the completeness check and exits non-zero if any artifact is missing.

---

## Event Log Format

Each line in `runs/events.jsonl` is a JSON object. The log is append-only and must never be modified after a line is written.

```json
{
  "ts": "2026-03-29T12:00:00Z",
  "run_id": "run_001",
  "step": "normalize",
  "entity": "document_family",
  "entity_id": "risk_register",
  "action": "canonicalized",
  "before": "Risk Register",
  "after": "risk_register",
  "meta": {}
}
```

**Required fields:** `ts`, `run_id`, `step`, `entity`, `entity_id`, `action`

**Forbidden operations:** delete, modify, truncate, re-order existing lines

---

## SQLite Checksum

The checksum artifact (`db_checksum.txt`) must contain:

```
sha256:<hex>  db/semantic_index.sqlite
```

It is produced immediately after the last write of the run, before process exit. It covers the full database file including WAL-merged state.

---

## Run Summary Format

`reports/<run_id>/summary.md` must contain at minimum:

- Run ID
- Command invoked (with flags)
- Start and finish times
- Exit code
- Number of entities processed (by type)
- Any errors or warnings encountered
- Gate status (pass / fail / not evaluated)

---

## Reproducibility Guarantee

Given the same source files and the same run configuration, the evidence pack (including SQLite state and event log entries) must be functionally identical between runs.

"Functionally identical" means:
- same row counts in all tables,
- same canonical IDs assigned,
- same relations inferred,
- same gate outcomes.

Timestamps and run IDs are exempt from the reproducibility requirement.

Any function that breaks the functional reproducibility guarantee must be explicitly marked in code with a comment: `// non-deterministic: <reason>`.
