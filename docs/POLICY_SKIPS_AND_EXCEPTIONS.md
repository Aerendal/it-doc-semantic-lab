# Policy: Skips and Exceptions

This document defines when a test layer, quality gate check, or evidence requirement may be skipped, how to register the exception, and how skips affect gate credibility.

---

## Principle

**Skipping is a documented decision, not a convenience.**

A silent skip — one that is not registered, reviewed, and time-bounded — is treated as a gate failure equivalent. It removes the credibility of any gate that depends on the skipped layer.

---

## Skip Categories

### Category 1: Permitted skip — infrastructure not available

A layer may be skipped without gate impact if:
- the infrastructure it depends on is genuinely not available in the current environment (e.g., no network, no external service),
- the skip is registered per the procedure below,
- the skipped layer is **not** on the critical path for the gate being evaluated.

**Example:** Layer 28 (Performance Budget Tests) may be skipped in a low-resource CI environment if the performance budget is not yet formally defined.

### Category 2: Permitted skip — layer not yet implemented

A layer may be skipped without gate impact if:
- it is explicitly listed as `not yet implemented` in `docs/TEST_CATALOG.md`,
- the gate that depends on it is not being evaluated in this run,
- the skip is registered.

### Category 3: Restricted skip — requires explicit approval

A layer may be skipped with gate impact acknowledged if:
- the layer is on the critical path for an active gate,
- the skip is registered with an owner and a review date,
- the gate result is marked `degraded` (not `PASS`) until the layer is reinstated.

**Layers that may not be skipped without approval:**
- Layer 1 (File Presence), Layer 5 (Source Schema) — Gate 1
- Layer 8 (Golden Extraction), Layer 10 (Determinism) — Gate 2
- Layer 18 (Explainability), Layer 19 (Acyclicity) — Gate 4
- Layer 27 (Evidence Pack) — Gate 5

### Category 4: Forbidden skip — never permitted

The following may **never** be skipped under any circumstance:

| Item | Reason |
|------|--------|
| Evidence pack production (Layer 27) | Without evidence, the run cannot be audited |
| Exit code recording | Required for all gate evaluations |
| Event log append | Required for reproducibility and audit |
| SQLite run record creation | Required for run tracking |

Attempting to skip a Category 4 item results in automatic gate failure for all gates in the run.

---

## Exception Registration

Every skip (Categories 1–3) must be registered. An unregistered skip is treated as Category 4.

### Required fields

```
skip_id:       <unique identifier, e.g., SKIP-2026-001>
layer:         <layer number and name, e.g., Layer 28 — Performance Budget Tests>
category:      <1 | 2 | 3>
reason:        <brief factual description — what is missing and why>
owner:         <name or role responsible for this exception>
registered_at: <ISO 8601 date>
review_date:   <ISO 8601 date — when this skip must be reviewed>
gate_impact:   <none | degraded:<gate_id>>
```

### Registration location

Skips are registered in `docs/SKIP_REGISTER.md`. Each skip is a single entry in that file.

---

## Skip Expiry

A skip expires on its `review_date`. On expiry:

1. The owner must either reinstate the layer or renew the skip with a new `review_date`.
2. An expired skip without renewal is automatically treated as Category 4 (forbidden).
3. The `SKIP_REGISTER.md` entry must be updated with the outcome.

---

## Gate Impact of Skips

| Category | Gate impact |
|----------|------------|
| 1 (infra unavailable, non-critical) | None — gate evaluates remaining layers |
| 2 (not yet implemented, non-critical) | None — gate evaluates remaining layers |
| 3 (approved, critical path) | Gate marked `degraded` — not eligible for promotion |
| 4 (forbidden) or unregistered | Automatic gate failure for all active gates |

A `degraded` gate result means:
- the run is not promotable to the stable repository,
- the gate status must be recorded as `degraded` in the run summary,
- the degraded state must be visible in `audit evidence` output.

---

## What a "Credible Green" Means

A gate result of `PASS` is credible only if:
1. All required test layers were executed (not skipped).
2. All registered skips are Category 1 or 2 (non-critical).
3. No Category 4 items were bypassed.
4. The evidence pack is complete.

A result that does not meet these conditions must be reported as `degraded`, not `PASS`.
