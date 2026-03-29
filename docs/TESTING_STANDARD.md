# Testing Standard

## Philosophy

Testing in this repository is **evidence-driven**: every test produces or validates an artifact that can be independently audited. Tests are not just regression checks — they are part of the run evidence model.

---

## Non-Goals

This standard does **not** claim or assume:

1. **Failure-free execution.** The repository assumes failures will occur. The obligation is that they must be visible, not that they must not happen.
2. **Green = ready.** A passing test suite without a complete evidence pack is not a credible result.
3. **Green = promotable.** Feature code may be experimentally green without satisfying Gate 4 + evidence requirements for promotion to the stable repository.
4. **Any test = strong evidence.** A mock-based test of a mock-forbidden layer does not count toward gate credibility.
5. **Coverage = completeness.** High line coverage does not replace layer coverage. A layer not executed is a gap, regardless of line coverage metrics.

---

## 6 Levels, 30 Layers

Tests are organized into 6 levels. Each level addresses a distinct class of failure risk.

| Level | Name | Layers | Primary Risk Addressed |
|-------|------|--------|------------------------|
| A | Contract & Input | 1–5 | Corrupt or missing source data |
| B | Parser & Extraction | 6–10 | Wrong output from parsing |
| C | Normalization & Canonical Model | 11–15 | Identity collisions, type drift |
| D | Relations & Semantics | 16–20 | Incorrect inference, unexplained edges |
| E | Interface & Run | 21–25 | CLI contract breakage, run state corruption |
| F | Operational & Audit | 26–30 | Non-reproducibility, missing evidence, gate failures |

Full layer definitions: see [TEST_CATALOG.md](TEST_CATALOG.md).

---

## Mock Policy Summary

Each layer has an assigned mock policy. See `docs/POLICY_MOCKS_AND_REAL_PATHS.md` for full definitions and the layer-by-layer table.

| Policy | Meaning |
|--------|---------|
| **mock-allowed** | Test doubles permitted for all dependencies |
| **mock-restricted** | Core I/O (filesystem, SQLite, event log) must use real implementations; clock/random may be mocked |
| **mock-forbidden** | All primary dependencies must be real; no mocking of filesystem, SQLite, or event log |

Using mocks in a `mock-forbidden` layer invalidates that layer's gate contribution. It must be registered as a skip (Category 3) per `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`.

---

## Skip Policy Summary

Each skip must be registered in `docs/SKIP_REGISTER.md`. See `docs/POLICY_SKIPS_AND_EXCEPTIONS.md` for full procedure.

| Category | Condition | Gate impact |
|----------|-----------|-------------|
| 1 | Infrastructure unavailable, layer not on critical gate path | None |
| 2 | Layer not yet implemented, gate not being evaluated | None |
| 3 | Layer on critical gate path, approved with owner + review date | Gate marked `degraded` — not promotable |
| 4 / unregistered | Forbidden layers, or any skip without registration | Automatic gate failure |

**Forbidden skips (never permitted under any condition):**
- Evidence pack production (Layer 27)
- Exit code recording
- Event log append
- SQLite run record creation

---

## Quality Gates

A run **must not** be promoted until:

- All Level A tests pass (input contract validated)
- All Level B golden tests match
- Level C: zero unresolved collisions
- Level D: all relations have non-empty explanation
- Level E: CLI exit codes conform to contract
- Level F: evidence pack is complete

Full gate definitions: see [QUALITY_GATES.md](QUALITY_GATES.md) and [QUALITY_GATES_POLICY.md](QUALITY_GATES_POLICY.md).

---

## Promotion Expectation

A feature may be **experimentally green** (local test pass, working prototype) without being **promotable**.

For promotion to the stable repository (`IT-Dokumentacja`), a feature must satisfy:
1. Gate 4 (semantic consistency) fully passed — no `degraded` gates.
2. Evidence pack complete for the promotion run.
3. No Category 3 or Category 4 skips on any layer that feeds the promotion gate.
4. Golden files are current and reviewed.
5. `itdlab export repo1` exits with code 0.

A feature that is experimentally green but not promotable should be documented as such in `docs/DEVELOPMENT_PLAN.md`.

---

## Review Standard

A reviewer must be able to independently verify a run result **without talking to the author**.

Specifically, a review is credible when:
- the evidence pack for the run is complete (`itdlab audit evidence <run_id>` exits 0),
- the event log entries for the run are present and parseable,
- the golden files are current and match the run output,
- any skips or degraded modes are registered in `docs/SKIP_REGISTER.md`,
- the gate status is explicitly recorded (PASS / degraded / FAIL), not implied.

A review that is not independently reproducible from the evidence pack is not a credible approval.

---

## Test Types

### Unit tests
- Package-level, in `_test.go` files alongside source
- No I/O, no filesystem
- Must be deterministic

### Fixture tests
- Use files from `internal/testkit/fixtures/`
- Input is a known file → verify exact output

### Golden tests
- Compare output to files in `internal/testkit/golden/`
- Update with `go test ./... -update` when golden intentionally changes
- Golden file changes must be reviewed in PR diff

### Integration tests
- In `internal/testkit/` or `_integration_test.go` files
- Require real SQLite and event log
- Skipped with `-short` flag

### Contract tests
- Verify CLI flag contracts and exit codes
- Run as part of standard test suite

---

## Conventions

- Test file names: `<subject>_test.go`
- Fixture file names: `<subject>_<case>.md` / `.json`
- Golden file names: `<test_name>.golden`
- Run `make test-short` to skip integration tests
- Run `make test` for full suite

---

## Determinism Requirement

Any function that produces output used in tests **must** produce identical output for identical input. Non-determinism is a test failure.

Non-deterministic functions (e.g., UUID generators, timestamp producers) must be injectable and replaced with deterministic fakes in tests. They must be marked in code with: `// non-deterministic: <reason>`.
