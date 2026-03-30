# Risk-to-Test Matrix

This document maps identified risk categories to the test layers that address them, the mock policy that applies, and the evidence strength required for a credible result.

---

## Purpose

The risk-to-test matrix makes explicit the connection between what can go wrong and which test layers are the primary defense against each failure mode. It supports:
- gate evaluation (which risks are covered before promotion),
- skip impact assessment (which risks become uncovered when a layer is skipped),
- test prioritization (which layers to build first for a new capability slice).

---

## Risk Categories

| ID | Risk Category | Description |
|----|--------------|-------------|
| R1 | Corrupt or missing input | Source files absent, unreadable, or malformed |
| R2 | Parser failure | Parser produces wrong output, panics, or is non-deterministic |
| R3 | Identity instability | Same concept gets different canonical IDs across runs |
| R4 | Silent name collision | Two distinct documents mapped to the same canonical ID without detection |
| R5 | Schema drift | Domain model changes break existing persisted data |
| R6 | Incorrect relation inference | Relations inferred incorrectly or without explanation |
| R7 | Graph integrity violation | Dependency cycles or contradictory relations undetected |
| R8 | CLI contract breakage | Command flags, exit codes, or output format changed silently |
| R9 | Run state corruption | Partial writes leave SQLite in inconsistent state |
| R10 | Evidence loss | Run completes but artifacts missing or empty |
| R11 | Non-reproducibility | Same input produces different output on repeated runs |
| R12 | Silent skip | Test layer not executed without registration or notice |
| R13 | Promotion without gate | Feature promoted to stable repo without all gates passing |
| R14 | Mock masking real failure | Mock-forbidden layer executed with mocks, hiding real failures |

---

## Risk-to-Layer Mapping

| Risk | Primary layers | Supporting layers | Blocking gate |
|------|---------------|-------------------|---------------|
| R1 — Corrupt/missing input | 1, 2, 3, 4, 5 | — | G1 |
| R2 — Parser failure | 6, 7, 8, 9, 10 | — | G2 |
| R3 — Identity instability | 11, 13 | 10 | G3 |
| R4 — Silent name collision | 12 | 11 | G3 |
| R5 — Schema drift | 14, 15 | — | G3 |
| R6 — Incorrect relation inference | 16, 17, 18 | 20 | G4 |
| R7 — Graph integrity violation | 17, 19 | — | G4 |
| R8 — CLI contract breakage | 21 | 22 | G4, G5 |
| R9 — Run state corruption | 23, 25 | 22 | G5 |
| R10 — Evidence loss | 27, 24 | 22 | G5 |
| R11 — Non-reproducibility | 26, 10 | 8 | G5 |
| R12 — Silent skip | (policy enforcement) | 27 | all |
| R13 — Promotion without gate | 30 | 27 | G5 |
| R14 — Mock masking real failure | (policy enforcement) | 22, 25, 26 | G4, G5 |

---

## Risk vs. Mock Policy

Risks that are most sensitive to mock policy violations:

| Risk | Why mocks are dangerous here |
|------|------------------------------|
| R2 — Parser failure | A mock parser always returns expected output; real failures invisible |
| R9 — Run state corruption | In-memory SQLite does not exercise WAL, crash recovery, or disk write failures |
| R10 — Evidence loss | Mock filesystem silently swallows artifacts |
| R11 — Non-reproducibility | Determinism tests with mocked I/O do not detect file-system-level variability |
| R14 — Mock masking | By definition |

Layers addressing R2, R9, R10, R11 are `mock-restricted` or `mock-forbidden` per `docs/POLICY_MOCKS_AND_REAL_PATHS.md`.

---

## Risk vs. Evidence Strength

| Evidence strength | Risks it credibly covers |
|------------------|--------------------------|
| `low` | R2 (partial), R3 (partial) — unit-level only |
| `medium` | R1, R2, R3, R4, R5, R6, R7, R8 |
| `high` | R1, R4, R5, R7, R9, R13 |
| `promotion-critical` | R2, R3, R10, R11, R13, R14 |

For promotion, all `promotion-critical` layers must be executed without mocks and must produce PASS results with artifacts.

---

## Coverage by Capability Slice

| Slice | Primary risks covered | Minimum layers required |
|-------|-----------------------|------------------------|
| 1 — Ingest | R1, R2 | Layers 1–10 |
| 2 — Normalize | R3, R4, R5 | Layers 11–15 |
| 3 — Relations | R6, R7 | Layers 16–20 |
| 4 — CLI + Run | R8, R9 | Layers 21–25 |
| 5 — Audit + Export | R10, R11, R12, R13, R14 | Layers 26–30 |

A capability slice is not promotion-ready until all its minimum required layers are passing with `mock-forbidden` or `mock-restricted + real-path` execution.

---

## Internal references
- `docs/TESTING_STANDARD.md`
- `docs/TEST_CATALOG.md`
- `docs/QUALITY_GATES_POLICY.md`
- `docs/POLICY_MOCKS_AND_REAL_PATHS.md`
- `docs/CONTEXT_VOCABULARY.md`

## Review metadata
- Owner: project team
- Status: draft
- Last reviewed: 2026-03-30
