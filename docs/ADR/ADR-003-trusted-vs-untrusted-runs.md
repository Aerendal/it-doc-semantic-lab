# ADR-003: Trusted vs. Untrusted Runs

| Field | Value |
|-------|-------|
| Status | Accepted |
| Date | 2026-03-30 |
| Deciders | project team |
| Supersedes | — |
| Superseded by | — |

---

## Context

Not all runs are equal. A run executed in a development environment with mocked dependencies, skipped layers, or an incomplete evidence pack should not be treated the same as a run executed on a clean environment with real inputs and full evidence. Treating all green runs as equivalent creates a false sense of quality.

The problem is not that incomplete runs are wrong — they are often useful during development. The problem is when an incomplete run is *cited* as gate evidence or *used* as the basis for a promotion decision.

The system needs a formal distinction between runs that can be cited and runs that cannot.

---

## Decision

Every run is classified as either **trusted** or **untrusted** at the time the run manifest is written.

### A run is trusted if and only if:

1. All preconditions in `docs/EXECUTION_CONTRACT.md` were met.
2. All postconditions in `docs/EXECUTION_CONTRACT.md` were satisfied (`exit_code = 0` or `exit_code = 2`).
3. `evidence.complete = true` in `run_manifest.json`.
4. No `mock-forbidden` layer was executed with mocks.
5. No Category 3 or Category 4 skip was active (per `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`).
6. The environment was not a local development machine with overridden DB or log paths that bypass standard isolation (specifically: `--db` and `--log` must point to the standard paths or explicitly declared alternate paths with a documented reason).

If any of these conditions is false, `trusted = false` in the run manifest. The run is recorded and available for inspection, but it may not be cited as evidence for gate evaluation or promotion.

### Consequences of untrusted classification

- The run may not be referenced in a gate report as a PASS citation.
- The run may not appear in a promotion-readiness checklist as a completed item.
- `itdlab audit evidence <run_id>` will report the run as untrusted.
- The run's test results may still be used informally for development feedback, but this must not be confused with gate evidence.

---

## Alternatives Considered

### 1. All runs are equal; trust is asserted externally

**Approach:** Do not classify runs; let the operator declare which runs they trust.

**Rejected because:**
- This removes objective enforceability. An operator under pressure will declare untrusted runs as trusted.
- The classification must be deterministic and computed by the tool, not declared by the operator.

### 2. Trust is binary: pass/fail only

**Approach:** A run that exits 0 is trusted; any other exit code is not.

**Rejected because:**
- This allows runs with incomplete evidence packs and active mock-forbidden violations to be classified as trusted, as long as the run exits 0.
- Green exit code is necessary but not sufficient for trust.

### 3. Separate "audit mode" flag

**Approach:** Add a `--trusted` flag that enables stricter checking.

**Rejected because:**
- This creates a two-track system where most runs are never audited.
- Developers will not use `--trusted` during normal development; it provides no protection against accidental misuse.
- Trust classification should be automatic, not opt-in.

---

## Consequences

### Positive

- Gate evidence reports can clearly distinguish trusted from untrusted runs.
- Promotion decisions have a clear, objective basis: at least one trusted run per gate.
- Development workflow is not blocked — untrusted runs are still fast, useful, and recorded.
- The classification is computed deterministically; no ambiguity about a run's status.

### Negative / Accepted tradeoffs

- Some false positives: a run on a developer machine with full evidence but non-standard paths will be classified as untrusted. The developer can override by providing the standard paths.
- The trusted/untrusted distinction introduces a classification step in the audit report that must be maintained as new conditions are added.

### Deferred

- Per-organisation trust profiles (e.g., CI always trusted regardless of path) are deferred.
- Cryptographic signing of trusted run manifests is deferred.

---

## Implementation Notes

- `trusted` field is computed and written to `run_manifest.json` by the run finalisation logic.
- `itdlab audit runs` displays trust status for each run.
- `itdlab audit evidence <run_id>` validates all conditions and reports which failed if `trusted = false`.
- Go implementation: trust classification logic should be in `internal/app/audit/trust.go`.

---

## Internal references
- `docs/EXECUTION_CONTRACT.md`
- `docs/EVIDENCE_MODEL.md`
- `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`
- `docs/POLICY_MOCKS_AND_REAL_PATHS.md`
- `docs/CONTEXT_VOCABULARY.md`
- `docs/ADR/ADR-001-sqlite-as-source-of-truth.md`
- `docs/ADR/ADR-002-event-log-as-audit-backbone.md`

## Review metadata
- Owner: project team
- Status: accepted
- Last reviewed: 2026-03-30
