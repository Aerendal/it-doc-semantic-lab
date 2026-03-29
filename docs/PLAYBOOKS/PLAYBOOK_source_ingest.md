# Playbook: Source Ingest

## Purpose

Defines the strategy for ingesting raw IT documentation sources into the semantic lab.

---

## When to Use This Playbook

- Adding new source documents to `sources/`
- Re-ingesting after source files are updated
- Debugging a failed ingest run

---

## Ingest Strategy

### 1. Source Placement

Place raw Markdown files in the appropriate subdirectory:

```
sources/
  matrices/    — industry document matrices
  metagraph/   — metagraph definitions
  plans/       — project plans and roadmaps
```

Files must be UTF-8 encoded Markdown. No BOM. No binary attachments.

### 2. Pre-Ingest Validation

Before running ingest:
- Run Layer 1–5 tests (contract & input validation)
- Verify file encodings
- Confirm at least one heading per file

### 3. Ingest Run

```
itdlab ingest run --source sources/
```

This:
1. Discovers all `.md` files under the source path
2. Parses each file (headings, sections, metadata)
3. Writes a `Document` row to SQLite per file
4. Appends `ingested` events to `runs/events.jsonl`
5. Writes `source_manifest.json` and `parse_report.json` to `reports/<run_id>/`

### 4. Post-Ingest Verification

```
itdlab ingest inspect <path>   # review a specific file's parse result
```

Check:
- All expected documents appear in SQLite
- No parse errors in `parse_report.json`
- Event log has entries for all ingested files

---

## Failure Modes

| Symptom | Likely Cause | Action |
|---------|-------------|--------|
| File not in source manifest | File not in source path | Verify placement |
| Encoding error | BOM or non-UTF-8 | Convert with `iconv` |
| Zero headings | File has no `#` headings | Check Markdown structure |
| Gate 1 failure | Contract violation | Fix source file, re-run |
