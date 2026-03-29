# Architecture — IT Documentation Semantic Lab

## Overview

This repository is the **semantic analysis engine** for the IT-Dokumentacja project.  
It operates as a Go CLI tool backed by SQLite (source of truth) and a JSONL event log (audit trail).

The stable reference repository (`IT-Dokumentacja`) receives only **promoted, stable metadata** exported from here.

---

## Design Principles

| Principle | Expression |
|-----------|-----------|
| Local-first | No network dependencies at runtime |
| Auditable | Every state change is logged to `runs/events.jsonl` |
| Reproducible | Same input + same run = same output |
| Incremental | Capability slices, not big-bang layers |
| Evidence-driven | Every run produces a verifiable evidence pack |

---

## Technology Stack

| Concern | Choice | Rationale |
|---------|--------|-----------|
| Language | Go | Static binary, strong stdlib, easy testing |
| Database | SQLite (`modernc.org/sqlite`) | Local, no CGO, zero-config |
| Audit log | JSONL (append-only) | Human-readable, grep-able, append-safe |
| CLI framework | Cobra | Subcommand tree, flags, help generation |
| Test tooling | stdlib `testing` + golden files | No external deps |

---

## Repository Layout

```
cmd/itdlab/          — CLI entrypoint
internal/
  app/               — application-layer use cases (one package per capability)
    ingest/
    normalize/
    classify/
    relations/
    sections/
    authority/
    export/
    audit/
  domain/            — pure domain types (no I/O)
  ports/             — interfaces (source reader, event log, stores, report writer)
  adapters/
    sqlite/          — SQLite implementation of stores
    jsonl/           — JSONL event log implementation
    filesystem/      — filesystem source reader
    markdown/        — Markdown parser
  cli/               — Cobra command definitions
  testkit/           — test helpers, fixtures, golden files, builders
db/
  schema_v1.sql      — canonical DDL
  semantic_index.sqlite  — runtime database (gitignored)
runs/
  events.jsonl       — append-only event log (gitignored)
sources/             — raw IT documentation inputs
normalized/          — normalized outputs
reports/             — generated reports per run
docs/
  ARCHITECTURE.md    — this file
  PLAYBOOKS/         — strategic how-to guides
  RUNBOOKS/          — step-by-step operational procedures
  ADR/               — architectural decision records
  TESTING_STANDARD.md
  TEST_CATALOG.md
  EVIDENCE_MODEL.md
  QUALITY_GATES.md
```

---

## Capability Slices

Development follows vertical capability slices, not horizontal layers.

| Slice | Background | Interface | Evidence |
|-------|-----------|-----------|---------|
| 1 — Ingest | Markdown parser, SQLite store, event log | `itdlab ingest run`, `itdlab ingest inspect` | Source manifest, parse report |
| 2 — Normalize | Canonical IDs, dedup, collision detection | `itdlab normalize preview`, `itdlab normalize apply` | Normalization report |
| 3 — Sections | Section archetypes, role inference | `itdlab sections show`, `itdlab sections explain` | Section map, anomaly report |
| 4 — Relations | Rule-based inference, cross-doc deps | `itdlab relations show`, `itdlab relations explain` | Relation graph, candidate report |
| 5 — Authority | Regulatory linkage | `itdlab authority check` | Authority coverage report |
| 6 — Export | Stable metadata promotion | `itdlab export repo1` | Export manifest |

---

## State Model

```
raw → ingested → normalized → classified → exported
```

Each transition is:
1. Recorded as a row mutation in SQLite
2. Appended as an event to `runs/events.jsonl`

---

## Related Decisions

- [ADR-001: SQLite as source of truth](ADR/ADR-001-sqlite-as-source-of-truth.md)
