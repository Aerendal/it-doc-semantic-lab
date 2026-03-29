# Quality Gates

Quality gates are blocking conditions. A run or promotion **cannot proceed** past a gate unless all conditions for that gate are satisfied.

---

## Gate 1: Input Contract (before any processing)

| Condition | Test Layer |
|-----------|-----------|
| All required source files present | Layer 1 |
| All files readable | Layer 2 |
| Zero encoding violations | Layer 3 |
| All files have parseable structure | Layer 4 |
| All files satisfy source schema contract | Layer 5 |

**Blocked if:** Any condition fails → run aborts with `exit 1` and a contract failure report.

---

## Gate 2: Parse Quality (after ingest)

| Condition | Test Layer |
|-----------|-----------|
| All fixture golden tests pass | Layer 8 |
| No parser panics on any input | Layer 9 |
| Parser output is deterministic | Layer 10 |

**Blocked if:** Any condition fails → normalization step does not start.

---

## Gate 3: Normalization Integrity (after normalize)

| Condition | Test Layer |
|-----------|-----------|
| All canonical IDs are unique | Layer 11 |
| Zero unresolved collisions | Layer 12 |
| All aliases resolve cleanly | Layer 13 |
| All domain fields have correct types | Layer 14 |

**Blocked if:** Any condition fails → classify step does not start.

---

## Gate 4: Semantic Consistency (after classify + relations)

| Condition | Test Layer |
|-----------|-----------|
| All inferred relations have non-empty explanation | Layer 18 |
| Zero cycles in depends_on graph | Layer 19 |
| Zero contradictory relations | Layer 17 |

**Blocked if:** Any condition fails → export step does not start.

---

## Gate 5: Run Completeness (before export)

| Condition | Test Layer |
|-----------|-----------|
| Evidence pack is complete | Layer 27 |
| SQLite state is consistent with event log | Layer 25 |
| Run record in SQLite has `status = 'completed'` | Layer 22 |
| CLI exit code was 0 for all prior steps | Layer 21 |

**Blocked if:** Any condition fails → export to reference repository is blocked.

---

## Gate Failure Protocol

1. Log failure to event log (`action: "gate_failed"`, `entity: "gate"`, `entity_id: "<gate_id>"`)
2. Write failure details to `reports/<run_id>/gate_failures.json`
3. Exit with code `2` (gate failure, distinct from error code `1`)
4. Do not modify any downstream state
