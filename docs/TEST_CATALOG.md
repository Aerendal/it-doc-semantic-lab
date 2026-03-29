# Test Catalog

30 test layers organized in 6 levels.

---

## Level A — Contract & Input (Layers 1–5)

### Layer 1: File Presence Tests
- **Goal:** Verify all required source files exist before any processing.
- **Input:** Source directory path
- **Output:** PASS if all required files present, FAIL with missing file list
- **Artifact:** `presence_check.json`
- **PASS criterion:** Zero missing required files

### Layer 2: File Readability Tests
- **Goal:** Verify all source files can be opened and read.
- **Input:** Source file paths
- **Output:** PASS if all files readable, FAIL with error per file
- **Artifact:** `readability_report.json`
- **PASS criterion:** Zero read errors

### Layer 3: Encoding Tests
- **Goal:** Detect UTF-8 violations, BOM markers, null bytes.
- **Input:** Raw file bytes
- **Output:** Per-file encoding status
- **Artifact:** `encoding_report.json`
- **PASS criterion:** Zero encoding violations in required files

### Layer 4: Markdown Structure Tests
- **Goal:** Verify parser can detect at least one heading, list, or section boundary.
- **Input:** Markdown source file
- **Output:** Structural summary (heading count, list count)
- **Artifact:** embedded in `parse_report.json`
- **PASS criterion:** At least one heading detected

### Layer 5: Source Schema Tests
- **Goal:** Verify source files satisfy minimum input contract (class field, non-empty body).
- **Input:** Parsed source
- **Output:** Contract violations list
- **Artifact:** `source_contract_report.json`
- **PASS criterion:** Zero contract violations

---

## Level B — Parser & Extraction (Layers 6–10)

### Layer 6: Parser Unit Tests
- **Goal:** Each parser function handles a minimal case correctly.
- **Input:** Inline string literals
- **Output:** Parsed struct
- **PASS criterion:** Output matches expected struct

### Layer 7: Fixture Parsing Tests
- **Goal:** Known fixture files produce known output.
- **Input:** `internal/testkit/fixtures/`
- **Output:** Parsed document
- **PASS criterion:** Output matches expected JSON

### Layer 8: Golden Extraction Tests
- **Goal:** Full parser output matches golden file byte-for-byte.
- **Input:** Fixture file
- **Output:** JSON document
- **Artifact:** Updated on `-update` flag
- **PASS criterion:** Output == golden

### Layer 9: Partial Corruption Tests
- **Goal:** Malformed input does not panic; returns error gracefully.
- **Input:** Truncated or invalid Markdown
- **Output:** Error (not panic)
- **PASS criterion:** No panic, error is non-nil and descriptive

### Layer 10: Determinism Tests
- **Goal:** Same input always produces same output.
- **Input:** Fixture file, run N times
- **Output:** N identical results
- **PASS criterion:** All N outputs identical

---

## Level C — Normalization & Canonical Model (Layers 11–15)

### Layer 11: Canonical ID Tests
- **Goal:** Semantically equivalent names produce the same canonical ID.
- **Input:** Variant name strings (e.g. "Risk Register", "risk_register", "Risk-Register")
- **Output:** Canonical ID
- **PASS criterion:** All variants → same canonical ID

### Layer 12: Collision Detection Tests
- **Goal:** System detects when two distinct documents produce the same canonical ID.
- **Input:** Two documents with different content but similar names
- **Output:** Collision report entry
- **PASS criterion:** Collision detected and logged

### Layer 13: Alias Resolution Tests
- **Goal:** Aliases correctly map to their canonical document.
- **Input:** Alias strings registered against a document
- **Output:** Resolved document ID
- **PASS criterion:** All aliases resolve to correct canonical document

### Layer 14: Typing Tests
- **Goal:** All domain fields have correct types and valid values.
- **Input:** Persisted document from SQLite
- **Output:** Type validation report
- **PASS criterion:** Zero type violations

### Layer 15: Migration Tests
- **Goal:** Schema migrations do not destroy or corrupt existing data.
- **Input:** Pre-migration SQLite snapshot + migration script
- **Output:** Post-migration state
- **PASS criterion:** All rows intact, new columns have correct defaults

---

## Level D — Relations & Semantics (Layers 16–20)

### Layer 16: Relation Rule Unit Tests
- **Goal:** Each relation rule fires correctly on a minimal input pair.
- **Input:** Two document stubs + rule definition
- **Output:** Relation (or none)
- **PASS criterion:** Output matches expected relation

### Layer 17: Relation Consistency Tests
- **Goal:** No contradictory relations exist (e.g., A depends_on B and B depends_on A simultaneously as depends_on).
- **Input:** Full relation set from SQLite
- **Output:** Contradiction list
- **PASS criterion:** Zero contradictions

### Layer 18: Explainability Tests
- **Goal:** Every inferred relation has a non-empty explanation field.
- **Input:** All relations in SQLite where source = 'rule_engine'
- **Output:** Rows with empty explanation
- **PASS criterion:** Zero rows with empty explanation

### Layer 19: Graph Acyclicity Tests
- **Goal:** Detect cycles in directed relation graphs where cycles are invalid.
- **Input:** depends_on relation edges
- **Output:** Cycle list
- **PASS criterion:** Zero cycles in depends_on graph

### Layer 20: Section Influence Tests
- **Goal:** Section-to-section inference rules fire for known canonical cases.
- **Input:** Document with known section structure
- **Output:** Inferred section relations
- **PASS criterion:** Output matches expected section influence map

---

## Level E — Interface & Run (Layers 21–25)

### Layer 21: CLI Contract Tests
- **Goal:** All commands and flags exist; exit codes are correct.
- **Input:** CLI binary invocation
- **Output:** Exit code, stdout structure
- **PASS criterion:** Exit code 0 for valid input, non-zero for invalid

### Layer 22: End-to-End Slice Tests
- **Goal:** Full capability slice runs without error and produces expected state.
- **Input:** Fixture source directory
- **Output:** SQLite state + JSONL events
- **PASS criterion:** SQLite matches expected rows; events are non-empty and valid JSONL

### Layer 23: Resume / Restart Tests
- **Goal:** A run interrupted mid-way can be resumed without data loss or duplication.
- **Input:** Partially completed run state
- **Output:** Completed run state
- **PASS criterion:** Final state identical to uninterrupted run

### Layer 24: Event Log Integrity Tests
- **Goal:** JSONL event log is valid, append-only, and parseable line by line.
- **Input:** `runs/events.jsonl`
- **Output:** Parse result per line
- **PASS criterion:** Zero parse errors; no line modified after write

### Layer 25: SQLite Materialization Tests
- **Goal:** SQLite state is consistent with the event log.
- **Input:** SQLite + events.jsonl for same run
- **Output:** Diff between replayed events and SQLite state
- **PASS criterion:** Zero diffs

---

## Level F — Operational & Audit (Layers 26–30)

### Layer 26: Reproducibility Tests
- **Goal:** Running the same command on the same data twice produces the same SQLite state.
- **Input:** Source data + two identical run invocations
- **Output:** SQLite checksum pair
- **PASS criterion:** Both checksums identical

### Layer 27: Evidence Pack Tests
- **Goal:** Every completed run produces a full evidence pack.
- **Input:** Completed run
- **Output:** Evidence pack manifest
- **PASS criterion:** All required artifacts present (see EVIDENCE_MODEL.md)

### Layer 28: Performance Budget Tests
- **Goal:** Key operations complete within defined time and memory budgets.
- **Input:** Standard fixture corpus
- **Output:** Timing and memory measurements
- **PASS criterion:** Under budget thresholds (defined per command)

### Layer 29: Failure-Mode Tests
- **Goal:** Errors are caught, reported correctly, and do not leave partial state.
- **Input:** Injected failure scenarios (missing file, corrupt DB, write error)
- **Output:** Error message, exit code, state of DB and event log
- **PASS criterion:** Non-zero exit, descriptive error, no partial writes

### Layer 30: Release Gate Tests
- **Goal:** Promotion to the reference repository is blocked unless all quality gates pass.
- **Input:** Full run result
- **Output:** Gate pass/fail report
- **PASS criterion:** All gates pass before export command succeeds
