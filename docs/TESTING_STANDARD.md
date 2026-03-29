# Testing Standard

## Philosophy

Testing in this repository is **evidence-driven**: every test produces or validates an artifact that can be independently audited. Tests are not just regression checks — they are part of the run evidence model.

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

## Quality Gates

A run **must not** be promoted until:

- All Level A tests pass (input contract validated)
- All Level B golden tests match
- Level C: zero unresolved collisions
- Level D: all relations have non-empty explanation
- Level E: CLI exit codes conform to contract
- Level F: evidence pack is complete

Full gate definitions: see [QUALITY_GATES.md](QUALITY_GATES.md).

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
