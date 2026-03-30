# Run Manifest Schema

The run manifest (`reports/<run_id>/run_manifest.json`) is the machine-readable summary of a run's identity, scope, outcome, and evidence. It is produced at run completion and is required for all trusted runs.

---

## Purpose

The run manifest is the single authoritative index of a run's evidence pack. A reviewer or automated tool can read the manifest to determine:
- what command was run and with what parameters,
- what the outcome was,
- which artifacts were produced,
- whether the evidence pack is complete,
- which quality gates were evaluated and what their status was.

---

## Schema

```json
{
  "schema_version": 1,
  "run_id": "<string — unique run identifier>",
  "command": "<string — full CLI command as invoked, e.g. 'itdlab ingest run --source sources/'>",
  "started_at": "<string — ISO 8601 datetime>",
  "finished_at": "<string — ISO 8601 datetime>",
  "exit_code": "<integer — 0 | 1 | 2 | 3>",
  "status": "<string — 'completed' | 'failed' | 'aborted'>",
  "trusted": "<boolean — true if all trusted-run criteria are met>",
  "environment": {
    "binary_version": "<string — itdlab version>",
    "go_version": "<string — Go runtime version>",
    "os": "<string — operating system>",
    "db_path": "<string — path to SQLite file used>",
    "log_path": "<string — path to event log file used>"
  },
  "evidence": {
    "db_checksum": "<string — sha256:<hex>>",
    "event_count": "<integer — number of events appended in this run>",
    "artifacts": [
      {
        "name": "<string — artifact logical name, e.g. 'parse_report'>",
        "path": "<string — relative path from repo root>",
        "size_bytes": "<integer>",
        "sha256": "<string — sha256:<hex>>"
      }
    ],
    "complete": "<boolean — true if all required artifacts are present and non-empty>"
  },
  "gates": [
    {
      "gate_id": "<string — e.g. 'G1'>",
      "status": "<string — 'PASS' | 'degraded' | 'FAIL' | 'not_evaluated'>",
      "evaluated_at": "<string — ISO 8601 datetime>",
      "failures": [
        "<string — description of each failing condition>"
      ]
    }
  ],
  "skips": [
    {
      "skip_id": "<string — from SKIP_REGISTER.md>",
      "layer": "<string — e.g. 'Layer 28 — Performance Budget Tests'>",
      "category": "<integer — 1 | 2 | 3>",
      "gate_impact": "<string — 'none' | 'degraded:G4'>"
    }
  ],
  "entities_processed": {
    "documents": "<integer>",
    "sections": "<integer>",
    "relations": "<integer>",
    "normalizations": "<integer>"
  },
  "errors": [
    {
      "step": "<string>",
      "entity_id": "<string>",
      "message": "<string>"
    }
  ]
}
```

---

## Field Definitions

### Top-level fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `schema_version` | integer | yes | Always `1` for this schema version |
| `run_id` | string | yes | Unique run identifier. Must match SQLite `runs.run_id` |
| `command` | string | yes | Full CLI invocation string |
| `started_at` | string | yes | ISO 8601 UTC datetime |
| `finished_at` | string | yes | ISO 8601 UTC datetime |
| `exit_code` | integer | yes | 0, 1, 2, or 3 per exit code contract |
| `status` | string | yes | `completed`, `failed`, or `aborted` |
| `trusted` | boolean | yes | `true` only if all trusted-run criteria are met |

### `evidence` object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `db_checksum` | string | yes | SHA-256 of SQLite file at run end. Format: `sha256:<hex>` |
| `event_count` | integer | yes | Number of event log lines appended in this run |
| `artifacts` | array | yes | One entry per artifact produced |
| `complete` | boolean | yes | `true` if all required artifacts are present and non-empty |

### `gates` array

Each element describes one gate evaluation. Gates not relevant to the command are recorded with `status: "not_evaluated"`.

### `skips` array

Each active skip from `SKIP_REGISTER.md` that affected this run. Empty array if no skips were active.

### `errors` array

Non-fatal errors encountered during the run. A run may exit 0 with non-empty `errors` if errors were per-entity and did not affect the overall outcome.

---

## Validation Rules

1. `run_id` must match a row in SQLite `runs` table with the same `run_id`.
2. `db_checksum` must match the actual SHA-256 of the database file at the time the manifest was written.
3. `evidence.complete` must be `false` if any artifact in the `artifacts` array has `size_bytes = 0`.
4. `trusted` must be `false` if `evidence.complete = false`.
5. `trusted` must be `false` if any skip in `skips` has `category = 3` or `category = 4`.
6. Gate status `PASS` requires `evidence.complete = true`.

---

## Production

The manifest is produced by `itdlab` at the end of every state-changing run, written to:

```
reports/<run_id>/run_manifest.json
```

It is also referenced by `itdlab audit evidence <run_id>` for completeness verification.

---

## Internal references
- `docs/EXECUTION_CONTRACT.md`
- `docs/EVIDENCE_MODEL.md`
- `docs/CONTEXT_VOCABULARY.md`
- `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`

## Review metadata
- Owner: project team
- Status: draft
- Last reviewed: 2026-03-30
