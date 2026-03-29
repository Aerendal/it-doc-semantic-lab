# ADR-001: SQLite as Source of Truth

**Status:** Accepted  
**Deciders:** Project team  
**Date:** 2026-03-29

---

## Context

The semantic lab needs a queryable, persistent store for documents, sections, relations, normalizations, and run history.  
Options considered: flat files (YAML/JSON), PostgreSQL, Neo4j, SQLite.

## Decision

Use **SQLite** as the single source of truth for all lab state.

Complement it with an **append-only JSONL event log** (`runs/events.jsonl`) for audit, reproducibility, and run history — but SQLite is authoritative for current state.

## Rationale

| Criterion | SQLite | Alternatives |
|-----------|--------|-------------|
| Zero-config | ✅ file-based | ❌ PostgreSQL requires server |
| Queryable | ✅ full SQL | ❌ flat JSON not queryable |
| Graph relations | ✅ JOIN-based | Neo4j is overkill at this scale |
| CGO-free | ✅ `modernc.org/sqlite` | `mattn/go-sqlite3` requires CGO |
| Portable binary | ✅ | ❌ server-based DBs |
| Auditability | ✅ WAL mode + event log | |

## Consequences

- All writes go to SQLite first, then the event log.
- The JSONL log can reconstruct SQLite state from scratch (reproducibility guarantee).
- Schema migrations use `schema_version` table. Each version is an additive DDL file.
- Neo4j or other graph stores may be added later as **read-only views** on top of SQLite exports — never as the primary store.

## Rejected Alternatives

- **PostgreSQL** — requires external server, unnecessary for local lab tooling.
- **Neo4j** — appropriate for graph traversal at scale; premature here.
- **Flat YAML/JSON** — not queryable, no relational integrity.
