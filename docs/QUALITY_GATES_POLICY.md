# Quality Gates

## Purpose

This document defines the **quality gates** used in the experimental repository.

The goal of the gate model is to ensure that repository work is not treated as complete merely because code exists or because a partial run produced a superficially green result.

A change, run, experiment, or promotion candidate must pass the relevant gate set before it is considered:
- technically credible,
- operationally reproducible,
- auditable,
- eligible for promotion,
- or suitable for external review.

This document is part of the execution-assurance program described in `docs/EXECUTION_ASSURANCE_PROGRAM.md`.

---

## Principles

1. **Evidence over assertion**  
   A claim that a capability works must be backed by logs, manifests, reports, and repeatable execution.

2. **Scope-aware gates**  
   Not every change requires every gate. Gates are applied according to the type and impact of the change.

3. **No hidden softening of standards**  
   Skips, mocks, bypasses, degraded modes, and manual overrides must be explicit and reviewable.

4. **Promotion only after proof**  
   Promotion to the stable repository is allowed only after the promotion gate has passed.

5. **Fail closed where needed**  
   For critical paths, missing evidence or ambiguous status is treated as failure, not success.

---

# Gate levels

The repository uses five gate levels.

## Gate G0 — Repository hygiene gate

### Purpose
Ensure the repository is in a reviewable and non-chaotic state before deeper validation is attempted.

### Applies to
- every branch proposed for review,
- every PR candidate,
- every release candidate,
- every promotion candidate.

### Minimum checks
- working tree is clean for the evaluated state,
- required documentation files exist,
- no forbidden local artifacts are committed,
- no temporary/debug-only files are part of the change,
- file naming and repository structure remain consistent.

### Required evidence
- repository status summary,
- file inventory delta,
- forbidden-artifact check result.

### Failure examples
- stray archives committed,
- temporary reports committed,
- debug helper files left in tracked tree,
- undocumented repo structure changes.

---

## Gate G1 — Static contract gate

### Purpose
Ensure the repository contract is internally coherent before runtime execution.

### Applies to
- schema changes,
- CLI changes,
- configuration changes,
- source model changes,
- relation model changes,
- export model changes.

### Minimum checks
- required files and schemas load successfully,
- YAML/JSON/Markdown contract files parse correctly,
- field names, required fields, and enumerations are valid,
- no forbidden or deprecated field usage is introduced,
- command signatures remain interpretable.

### Required evidence
- schema validation report,
- static lint report,
- contract diff summary.

### Failure examples
- invalid YAML in schema files,
- missing required schema fields,
- mixed typing for the same field,
- broken command definitions,
- undocumented contract drift.

---

## Gate G2 — Execution integrity gate

### Purpose
Ensure that the run actually happened in a technically credible way.

### Applies to
- ingest runs,
- normalization runs,
- relation inference runs,
- export runs,
- experiment result claims.

### Minimum checks
- run manifest is generated,
- executed steps are listed,
- skipped steps are listed explicitly,
- configuration and input set are recorded,
- logs exist and correspond to the manifest,
- result status is explicit (`PASS`, `FAIL`, `WARN`, `INCOMPLETE`).

### Required evidence
- run manifest,
- execution logs,
- configuration snapshot,
- checksums or fingerprints of critical inputs,
- summary report.

### Failure examples
- no manifest,
- ambiguous run status,
- steps executed but not recorded,
- evidence pack missing mandatory files,
- mismatch between manifest and produced outputs.

---

## Gate G3 — Verification credibility gate

### Purpose
Ensure that validation was meaningful and not weakened by silent skips, unjustified mocks, or low-value evidence.

### Applies to
- all experiment slices,
- all semantic feature claims,
- all promotion candidates,
- all externally reviewed runs.

### Minimum checks
- applicable tests executed,
- skip usage is explicit and justified,
- xfail usage is explicit and justified,
- mocking on critical paths complies with policy,
- golden outputs are current and reviewable,
- determinism expectations are met or explained.

### Required evidence
- test summary,
- skip/xfail audit report,
- mock-path vs real-path report,
- determinism report,
- golden verification report.

### Failure examples
- critical-path behavior validated only through hidden mocks,
- new skip added without registration,
- green result with missing golden verification,
- unstable output without explanation,
- tests passing only in degraded mode.

---

## Gate G4 — Promotion gate

### Purpose
Decide whether an experimental capability is eligible for promotion to the stable repository.

### Applies to
- features intended for `IT-Dokumentacja`,
- exports intended to become repository contract,
- metadata fields intended for stable templates,
- reusable mechanisms intended for stable runtime.

### Minimum checks
- G0–G3 passed,
- gold standard exists,
- all required corpus cases pass,
- idempotency is confirmed,
- integration dry-run is green,
- stable repository validation remains green after integration.

### Required evidence
- promotion candidate report,
- corpus coverage report,
- idempotency report,
- integration verification report,
- rollback note or rollback feasibility statement.

### Failure examples
- passing only on one happy-path case,
- no gold standard,
- unstable repeat execution,
- stable repository broken after merge attempt,
- unclear migration consequences.

---

# Gate application matrix

## Change types and required minimum gate

| Change type | Minimum gate |
|---|---|
| Documentation-only note, no contract impact | G0 |
| Schema / YAML / contract definition change | G1 |
| Parser / normalizer / relation execution change | G2 |
| Test-affecting or verification-affecting change | G3 |
| Promotion to stable repository | G4 |

## Escalation rule
If a change touches more than one category, the **highest applicable gate wins**.

Example:
- a parser refactor that also changes schema rules and test expectations is **not** G2 only; it must satisfy **G3**.

---

# Gate status model

Each gate outcome must use one of the following statuses:

- `PASS` — gate satisfied,
- `WARN` — gate satisfied with a documented non-blocking issue,
- `FAIL` — gate not satisfied,
- `INCOMPLETE` — evaluation could not be completed due to missing required evidence.

## Interpretation rules

### PASS
Used only when all mandatory checks for the gate are satisfied.

### WARN
Used only when:
- the gate was executed,
- evidence is complete,
- the issue is documented,
- the issue is explicitly classified as non-blocking.

### FAIL
Used when any blocking requirement is violated.

### INCOMPLETE
Used when the evaluation itself is not trustworthy due to missing mandatory evidence.

> `INCOMPLETE` must never be reported as if it were equivalent to `PASS`.

---

# Forbidden patterns

The following patterns are forbidden unless explicitly documented and approved under the exception policy.

## Skip-related
- silent introduction of `skip` or `xfail`,
- broad conditional skip without explicit reason,
- skip used as substitute for fixing a broken path.

## Mock-related
- mocking critical runtime/data paths without explicit policy allowance,
- replacing real verification with fake adapters while reporting the result as real validation,
- monkeypatching behavior in a way that hides actual integration risk.

## Evidence-related
- claiming success without manifest and logs,
- deleting or omitting evidence for a claimed run,
- publishing a PASS summary with missing mandatory artifacts.

## Contract-related
- introducing new required fields without schema update,
- using the same field with conflicting semantics,
- changing export semantics without documenting migration consequences.

---

# Exceptions and deviations

Exceptions are allowed only when all of the following are true:

1. the exception is documented,
2. the reason is specific and technical,
3. the scope is limited,
4. the duration is limited or review-triggered,
5. the resulting risk is described,
6. the repository still reports the degraded status honestly.

## Required exception record fields
- exception ID,
- date,
- owner,
- affected gate,
- affected component,
- reason,
- risk,
- expiry/review date,
- mitigation plan.

Recommended file:
- `docs/EXCEPTION_REGISTER.md`

---

# Required evidence per gate

| Gate | Minimum evidence |
|---|---|
| G0 | repo status summary, forbidden-artifact check |
| G1 | schema validation report, static lint report |
| G2 | run manifest, logs, config snapshot, summary |
| G3 | test summary, skip audit, mock audit, determinism report |
| G4 | promotion report, corpus coverage, idempotency report, stable integration verification |

---

# Release and promotion rule

A component may be promoted from the experimental repository to the stable repository only if:

1. the intended stable scope is explicit,
2. G0–G4 pass,
3. the feature has a gold standard,
4. corpus coverage is sufficient for the claimed scope,
5. integration into the stable repository remains green,
6. rollback is feasible or explicitly analyzed.

If any of these conditions is missing, the feature remains experimental.

---

# Mapping to the execution-assurance program

This document implements **Stage 1 — Gate model and enforcement policy** from `docs/EXECUTION_ASSURANCE_PROGRAM.md`.

Follow-on documents expected after this one:
- `docs/TESTING_STANDARD.md`
- `docs/TEST_CATALOG.md`
- `docs/EVIDENCE_MODEL.md`
- `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`
- `docs/POLICY_MOCKS_AND_REAL_PATHS.md`

---

# Final note

The purpose of quality gates is not to make experimentation impossible.

The purpose is to ensure that experimentation remains:
- explicit,
- disciplined,
- reviewable,
- evidence-backed,
- and resistant to accidental or intentional weakening of technical standards.
