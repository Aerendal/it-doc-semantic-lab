# Evidence Model

Every completed run of `itdlab` must produce a verifiable **evidence pack** — a set of artifacts that allow independent audit of what happened, why, and with what result.

---

## Required Artifacts per Run

| Artifact | Location | Description |
|----------|----------|-------------|
| Event log entries | `runs/events.jsonl` | All events appended during this run |
| Run record | SQLite `runs` table | start time, finish time, exit code, status |
| Source manifest | `reports/<run_id>/source_manifest.json` | List of all source files processed |
| Parse report | `reports/<run_id>/parse_report.json` | Per-file parse outcome |
| Command output | `reports/<run_id>/stdout.txt` | Captured CLI stdout |
| Exit code | run record in SQLite | 0 = success, non-zero = failure |
| SQLite checksum | `reports/<run_id>/db_checksum.txt` | SHA-256 of DB file after run |
| Run summary | `reports/<run_id>/summary.md` | Human-readable run summary |

---

## Optional Artifacts (by command)

| Command | Additional Artifact |
|---------|-------------------|
| `normalize apply` | `normalization_report.json` |
| `relations show` | `relation_graph.json` |
| `authority check` | `authority_coverage_report.json` |
| `export repo1` | `export_manifest.json` |

---

## Evidence Pack Completeness Check

A run's evidence pack is **complete** when all required artifacts exist and are non-empty.

Verified by Layer 27 (Evidence Pack Tests) in the test catalog.

---

## Event Log Format

Each line in `runs/events.jsonl` is a JSON object:

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

### Required fields: `ts`, `run_id`, `step`, `entity`, `entity_id`, `action`

---

## Reproducibility Guarantee

Given the same source files and the same run configuration, the evidence pack (including SQLite state) must be byte-for-byte identical between runs.

Any function that breaks this guarantee must be explicitly marked and justified.
