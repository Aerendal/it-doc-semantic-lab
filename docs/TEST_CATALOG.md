# Test Catalog

30 test layers organized in 6 levels.

Each layer entry includes:
- **Goal** — what failure risk it addresses
- **Input** — what is fed to the test
- **Output** — what the test produces
- **Artifact** — file written to the evidence pack (if any)
- **PASS criterion** — condition for a credible pass
- **Blocking gate** — which quality gate(s) this layer feeds
- **Mock policy** — `allowed` / `restricted` / `forbidden` (see `docs/POLICY_MOCKS_AND_REAL_PATHS.md`)
- **Evidence strength** — `low` / `medium` / `high` / `promotion-critical`

---

## Level A — Contract & Input (Layers 1–5)

### Layer 1: File Presence Tests
- **Goal:** Verify all required source files exist before any processing.
- **Input:** Source directory path
- **Output:** PASS if all required files present, FAIL with missing file list
- **Artifact:** `presence_check.json`
- **PASS criterion:** Zero missing required files
- **Blocking gate:** G1
- **Mock policy:** forbidden
- **Evidence strength:** high

### Layer 2: File Readability Tests
- **Goal:** Verify all source files can be opened and read.
- **Input:** Source file paths
- **Output:** PASS if all files readable, FAIL with error per file
- **Artifact:** `readability_report.json`
- **PASS criterion:** Zero read errors
- **Blocking gate:** G1
- **Mock policy:** forbidden
- **Evidence strength:** high

### Layer 3: Encoding Tests
- **Goal:** Detect UTF-8 violations, BOM markers, null bytes.
- **Input:** Raw file bytes
- **Output:** Per-file encoding status
- **Artifact:** `encoding_report.json`
- **PASS criterion:** Zero encoding violations in required files
- **Blocking gate:** G1
- **Mock policy:** forbidden
- **Evidence strength:** medium

### Layer 4: Markdown Structure Tests
- **Goal:** Verify parser can detect at least one heading, list, or section boundary.
- **Input:** Markdown source file
- **Output:** Structural summary (heading count, list count)
- **Artifact:** embedded in `parse_report.json`
- **PASS criterion:** At least one heading detected
- **Blocking gate:** G1
- **Mock policy:** forbidden
- **Evidence strength:** medium

### Layer 5: Source Schema Tests
- **Goal:** Verify source files satisfy minimum input contract (class field, non-empty body).
- **Input:** Parsed source
- **Output:** Contract violations list
- **Artifact:** `source_contract_report.json`
- **PASS criterion:** Zero contract violations
- **Blocking gate:** G1
- **Mock policy:** forbidden
- **Evidence strength:** high

---

## Level B — Parser & Extraction (Layers 6–10)

### Layer 6: Parser Unit Tests
- **Goal:** Each parser function handles a minimal case correctly.
- **Input:** Inline string literals
- **Output:** Parsed struct
- **Artifact:** none (unit test output)
- **PASS criterion:** Output matches expected struct
- **Blocking gate:** G2
- **Mock policy:** allowed
- **Evidence strength:** low

### Layer 7: Fixture Parsing Tests
- **Goal:** Known fixture files produce known output.
- **Input:** `internal/testkit/fixtures/`
- **Output:** Parsed document
- **Artifact:** none (test output)
- **PASS criterion:** Output matches expected JSON
- **Blocking gate:** G2
- **Mock policy:** allowed
- **Evidence strength:** medium

### Layer 8: Golden Extraction Tests
- **Goal:** Full parser output matches golden file byte-for-byte.
- **Input:** Fixture file
- **Output:** JSON document
- **Artifact:** diff report if mismatch
- **PASS criterion:** Output == golden
- **Blocking gate:** G2
- **Mock policy:** restricted
- **Evidence strength:** promotion-critical

### Layer 9: Partial Corruption Tests
- **Goal:** Malformed input does not panic; returns error gracefully.
- **Input:** Truncated or invalid Markdown
- **Output:** Error (not panic)
- **Artifact:** none
- **PASS criterion:** No panic, error is non-nil and descriptive
- **Blocking gate:** G2
- **Mock policy:** allowed
- **Evidence strength:** medium

### Layer 10: Determinism Tests
- **Goal:** Same input always produces same output.
- **Input:** Fixture file, run N times
- **Output:** N identical results
- **Artifact:** none
- **PASS criterion:** All N outputs identical
- **Blocking gate:** G2
- **Mock policy:** restricted
- **Evidence strength:** promotion-critical

---

## Level C — Normalization & Canonical Model (Layers 11–15)

### Layer 11: Canonical ID Tests
- **Goal:** Semantically equivalent names produce the same canonical ID.
- **Input:** Variant name strings (e.g. "Risk Register", "risk_register", "Risk-Register")
- **Output:** Canonical ID
- **Artifact:** none
- **PASS criterion:** All variants → same canonical ID
- **Blocking gate:** G3
- **Mock policy:** allowed
- **Evidence strength:** medium

### Layer 12: Collision Detection Tests
- **Goal:** System detects when two distinct documents produce the same canonical ID.
- **Input:** Two documents with different content but similar names
- **Output:** Collision report entry
- **Artifact:** `collision_report.json`
- **PASS criterion:** Collision detected and logged
- **Blocking gate:** G3
- **Mock policy:** allowed
- **Evidence strength:** high

### Layer 13: Alias Resolution Tests
- **Goal:** Aliases correctly map to their canonical document.
- **Input:** Alias strings registered against a document
- **Output:** Resolved document ID
- **Artifact:** none
- **PASS criterion:** All aliases resolve to correct canonical document
- **Blocking gate:** G3
- **Mock policy:** allowed
- **Evidence strength:** medium

### Layer 14: Typing Tests
- **Goal:** All domain fields have correct types and valid values.
- **Input:** Persisted document from SQLite
- **Output:** Type validation report
- **Artifact:** `type_validation_report.json`
- **PASS criterion:** Zero type violations
- **Blocking gate:** G3
- **Mock policy:** restricted
- **Evidence strength:** medium

### Layer 15: Migration Tests
- **Goal:** Schema migrations do not destroy or corrupt existing data.
- **Input:** Pre-migration SQLite snapshot + migration script
- **Output:** Post-migration state
- **Artifact:** row-count comparison report
- **PASS criterion:** All rows intact, new columns have correct defaults
- **Blocking gate:** G3
- **Mock policy:** restricted
- **Evidence strength:** high

---

## Level D — Relations & Semantics (Layers 16–20)

### Layer 16: Relation Rule Unit Tests
- **Goal:** Each relation rule fires correctly on a minimal input pair.
- **Input:** Two document stubs + rule definition
- **Output:** Relation (or none)
- **Artifact:** none
- **PASS criterion:** Output matches expected relation
- **Blocking gate:** G4
- **Mock policy:** allowed
- **Evidence strength:** low

### Layer 17: Relation Consistency Tests
- **Goal:** No contradictory relations exist.
- **Input:** Full relation set from SQLite
- **Output:** Contradiction list
- **Artifact:** `relation_consistency_report.json`
- **PASS criterion:** Zero contradictions
- **Blocking gate:** G4
- **Mock policy:** restricted
- **Evidence strength:** high

### Layer 18: Explainability Tests
- **Goal:** Every inferred relation has a non-empty explanation field.
- **Input:** All relations in SQLite where source = 'rule_engine'
- **Output:** Rows with empty explanation
- **Artifact:** `explainability_report.json`
- **PASS criterion:** Zero rows with empty explanation
- **Blocking gate:** G4
- **Mock policy:** restricted
- **Evidence strength:** promotion-critical

### Layer 19: Graph Acyclicity Tests
- **Goal:** Detect cycles in directed relation graphs where cycles are invalid.
- **Input:** depends_on relation edges
- **Output:** Cycle list
- **Artifact:** `acyclicity_report.json`
- **PASS criterion:** Zero cycles in depends_on graph
- **Blocking gate:** G4
- **Mock policy:** restricted
- **Evidence strength:** high

### Layer 20: Section Influence Tests
- **Goal:** Section-to-section inference rules fire for known canonical cases.
- **Input:** Document with known section structure
- **Output:** Inferred section relations
- **Artifact:** none
- **PASS criterion:** Output matches expected section influence map
- **Blocking gate:** G4
- **Mock policy:** restricted
- **Evidence strength:** medium

---

## Level E — Interface & Run (Layers 21–25)

### Layer 21: CLI Contract Tests
- **Goal:** All commands and flags exist; exit codes are correct.
- **Input:** CLI binary invocation
- **Output:** Exit code, stdout structure
- **Artifact:** none
- **PASS criterion:** Exit code 0 for valid input, non-zero for invalid
- **Blocking gate:** G4, G5
- **Mock policy:** restricted
- **Evidence strength:** medium

### Layer 22: End-to-End Slice Tests
- **Goal:** Full capability slice runs without error and produces expected state.
- **Input:** Fixture source directory
- **Output:** SQLite state + JSONL events
- **Artifact:** SQLite row counts, event log line count
- **PASS criterion:** SQLite matches expected rows; events are non-empty and valid JSONL
- **Blocking gate:** G5
- **Mock policy:** forbidden
- **Evidence strength:** promotion-critical

### Layer 23: Resume / Restart Tests
- **Goal:** A run interrupted mid-way can be resumed without data loss or duplication.
- **Input:** Partially completed run state
- **Output:** Completed run state
- **Artifact:** row-count comparison report
- **PASS criterion:** Final state identical to uninterrupted run
- **Blocking gate:** G5
- **Mock policy:** forbidden
- **Evidence strength:** high

### Layer 24: Event Log Integrity Tests
- **Goal:** JSONL event log is valid, append-only, and parseable line by line.
- **Input:** `runs/events.jsonl`
- **Output:** Parse result per line
- **Artifact:** `event_log_integrity_report.json`
- **PASS criterion:** Zero parse errors; no line modified after write
- **Blocking gate:** G5
- **Mock policy:** forbidden
- **Evidence strength:** promotion-critical

### Layer 25: SQLite Materialization Tests
- **Goal:** SQLite state is consistent with the event log.
- **Input:** SQLite + events.jsonl for same run
- **Output:** Diff between replayed events and SQLite state
- **Artifact:** `materialization_diff_report.json`
- **PASS criterion:** Zero diffs
- **Blocking gate:** G5
- **Mock policy:** forbidden
- **Evidence strength:** promotion-critical

---

## Level F — Operational & Audit (Layers 26–30)

### Layer 26: Reproducibility Tests
- **Goal:** Running the same command on the same data twice produces the same SQLite state.
- **Input:** Source data + two identical run invocations
- **Output:** SQLite checksum pair
- **Artifact:** `reproducibility_report.json`
- **PASS criterion:** Both checksums identical
- **Blocking gate:** G5
- **Mock policy:** forbidden
- **Evidence strength:** promotion-critical

### Layer 27: Evidence Pack Tests
- **Goal:** Every completed run produces a full evidence pack.
- **Input:** Completed run
- **Output:** Evidence pack manifest
- **Artifact:** `evidence_manifest.json`
- **PASS criterion:** All required artifacts present (see `docs/EVIDENCE_MODEL.md`)
- **Blocking gate:** G5
- **Mock policy:** forbidden
- **Evidence strength:** promotion-critical

### Layer 28: Performance Budget Tests
- **Goal:** Key operations complete within defined time and memory budgets.
- **Input:** Standard fixture corpus
- **Output:** Timing and memory measurements
- **Artifact:** `performance_report.json`
- **PASS criterion:** Under budget thresholds (defined per command in Makefile)
- **Blocking gate:** none (advisory)
- **Mock policy:** forbidden
- **Evidence strength:** medium

### Layer 29: Failure-Mode Tests
- **Goal:** Errors are caught, reported correctly, and do not leave partial state.
- **Input:** Injected failure scenarios (missing file, corrupt DB, write error)
- **Output:** Error message, exit code, state of DB and event log
- **Artifact:** `failure_mode_report.json`
- **PASS criterion:** Non-zero exit, descriptive error, no partial writes
- **Blocking gate:** G5
- **Mock policy:** forbidden
- **Evidence strength:** high

### Layer 30: Release Gate Tests
- **Goal:** Promotion to the reference repository is blocked unless all quality gates pass.
- **Input:** Full run result
- **Output:** Gate pass/fail report
- **Artifact:** `gate_pass_report.json`
- **PASS criterion:** All gates G0–G4 pass; evidence pack complete; exit code 0
- **Blocking gate:** G5
- **Mock policy:** forbidden
- **Evidence strength:** promotion-critical
