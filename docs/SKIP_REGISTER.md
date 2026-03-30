# Skip Register

This is the living register of active test-layer skips and approved exceptions.

Maintained per `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`.

A skip entered here does **not** make it disappear — it makes it visible, owned, and reviewable.

---

## Active Skips

| Skip ID | Layer | Category | Owner | Registered | Review Date | Gate Impact | Reason |
|---------|-------|----------|-------|------------|-------------|-------------|--------|
| — | — | — | — | — | — | — | *(no active skips)* |

---

## Skip ID Format

`SKIP-<YYYY>-<NNN>` — e.g. `SKIP-2026-001`

Increment `NNN` sequentially per year. Do not reuse IDs.

---

## How to Register a Skip

1. Add a row to the **Active Skips** table above.
2. Record the Skip ID in the relevant test file as a comment:
   ```go
   // SKIP-2026-001: Layer 28 skipped; no performance baseline yet. Review: 2026-09-01.
   t.Skip("SKIP-2026-001")
   ```
3. Record the skip in any run manifest for runs affected by it (`skips` array in `run_manifest.json`).
4. Add to **Closed Skips** when the skip is resolved or expired.

---

## Mandatory Fields

Per `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`:

| Field | Required for Category |
|-------|-----------------------|
| Skip ID | 1, 2, 3, 4 |
| Layer | 1, 2, 3, 4 |
| Category (1/2/3/4) | 1, 2, 3, 4 |
| Owner | 1, 2, 3, 4 |
| Registered date | 1, 2, 3, 4 |
| Review date | 1, 2 |
| Gate impact | 1, 2, 3, 4 |
| Reason | 1, 2, 3, 4 |

Category 3 (skip on blocked critical path) and Category 4 (permanent skip) require approval before registration.

---

## Closed Skips

| Skip ID | Layer | Closed Date | Resolution |
|---------|-------|-------------|------------|
| — | — | — | *(no closed skips)* |

---

## Internal references
- `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`
- `docs/TESTING_STANDARD.md`
- `docs/QUALITY_GATES_POLICY.md`

## Review metadata
- Owner: project team
- Status: active (living document)
- Last reviewed: 2026-03-30
