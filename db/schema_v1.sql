-- SQLite schema v1 for it-doc-semantic-lab
-- Source of truth for all lab state.
-- Event log (runs/events.jsonl) is the complementary append-only audit trail.

PRAGMA journal_mode = WAL;
PRAGMA foreign_keys = ON;

-- ─────────────────────────────────────────────
-- RUNS
-- Tracks each invocation of the CLI.
-- ─────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS runs (
    run_id      TEXT PRIMARY KEY,
    command     TEXT NOT NULL,
    started_at  TEXT NOT NULL,  -- ISO 8601
    finished_at TEXT,
    exit_code   INTEGER,
    status      TEXT NOT NULL DEFAULT 'running'
        CHECK (status IN ('running','completed','failed','aborted'))
);

-- ─────────────────────────────────────────────
-- DOCUMENTS
-- Core entity: one row per ingested document.
-- ─────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS documents (
    id              TEXT PRIMARY KEY,
    canonical_id    TEXT NOT NULL,
    raw_name        TEXT NOT NULL,
    source_path     TEXT NOT NULL,
    class           TEXT NOT NULL DEFAULT '',
    industry        TEXT NOT NULL DEFAULT 'generic',
    phase           TEXT NOT NULL DEFAULT 'unspecified',
    status          TEXT NOT NULL DEFAULT 'raw'
        CHECK (status IN ('raw','ingested','normalized','classified','exported')),
    checksum        TEXT NOT NULL DEFAULT '',
    ingested_at     TEXT NOT NULL,
    updated_at      TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_documents_canonical_id ON documents(canonical_id);
CREATE INDEX IF NOT EXISTS idx_documents_class        ON documents(class);
CREATE INDEX IF NOT EXISTS idx_documents_status       ON documents(status);

-- ─────────────────────────────────────────────
-- NORMALIZATIONS
-- Records every normalization decision for a document.
-- ─────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS normalizations (
    id              TEXT PRIMARY KEY,
    document_id     TEXT NOT NULL REFERENCES documents(id),
    run_id          TEXT NOT NULL REFERENCES runs(run_id),
    field           TEXT NOT NULL,
    before_value    TEXT NOT NULL,
    after_value     TEXT NOT NULL,
    rule_applied    TEXT NOT NULL DEFAULT '',
    normalized_at   TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_normalizations_document_id ON normalizations(document_id);

-- ─────────────────────────────────────────────
-- SECTIONS
-- Parsed sections within a document.
-- ─────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS sections (
    id          TEXT PRIMARY KEY,
    document_id TEXT NOT NULL REFERENCES documents(id),
    heading     TEXT NOT NULL,
    level       INTEGER NOT NULL DEFAULT 1,
    role        TEXT NOT NULL DEFAULT 'unknown',
    confidence  REAL NOT NULL DEFAULT 0.0,
    position    INTEGER NOT NULL DEFAULT 0,
    word_count  INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_sections_document_id ON sections(document_id);
CREATE INDEX IF NOT EXISTS idx_sections_role        ON sections(role);

-- ─────────────────────────────────────────────
-- SECTION ANOMALIES
-- Detected problems with sections.
-- ─────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS section_anomalies (
    id           TEXT PRIMARY KEY,
    section_id   TEXT NOT NULL REFERENCES sections(id),
    document_id  TEXT NOT NULL REFERENCES documents(id),
    anomaly_type TEXT NOT NULL,
    description  TEXT NOT NULL
);

-- ─────────────────────────────────────────────
-- RELATION RULES
-- Rule definitions used to infer relations between documents.
-- ─────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS relation_rules (
    id          TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    from_class  TEXT NOT NULL,
    to_class    TEXT NOT NULL,
    rel_type    TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT ''
);

-- ─────────────────────────────────────────────
-- RELATIONS
-- Directed semantic relationships between documents.
-- ─────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS relations (
    id            TEXT PRIMARY KEY,
    from_id       TEXT NOT NULL REFERENCES documents(id),
    to_id         TEXT NOT NULL REFERENCES documents(id),
    type          TEXT NOT NULL,
    source        TEXT NOT NULL DEFAULT 'rule_engine'
        CHECK (source IN ('rule_engine','explicit','inferred')),
    confidence    REAL NOT NULL DEFAULT 1.0,
    explanation   TEXT NOT NULL DEFAULT '',
    rule_id       TEXT REFERENCES relation_rules(id),
    discovered_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_relations_from_id ON relations(from_id);
CREATE INDEX IF NOT EXISTS idx_relations_to_id   ON relations(to_id);
CREATE INDEX IF NOT EXISTS idx_relations_type    ON relations(type);

-- ─────────────────────────────────────────────
-- AUTHORITY REFS
-- Links documents/sections to regulatory authorities.
-- ─────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS authority_refs (
    id          TEXT PRIMARY KEY,
    document_id TEXT NOT NULL REFERENCES documents(id),
    section_id  TEXT REFERENCES sections(id),
    authority   TEXT NOT NULL,
    clause      TEXT NOT NULL,
    url         TEXT NOT NULL DEFAULT '',
    linked_at   TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_authority_refs_document_id ON authority_refs(document_id);
CREATE INDEX IF NOT EXISTS idx_authority_refs_authority   ON authority_refs(authority);

-- ─────────────────────────────────────────────
-- SCHEMA VERSION
-- Single-row table for migration tracking.
-- ─────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS schema_version (
    version     INTEGER PRIMARY KEY,
    applied_at  TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT ''
);

INSERT OR IGNORE INTO schema_version (version, applied_at, description)
VALUES (1, datetime('now'), 'initial schema');
