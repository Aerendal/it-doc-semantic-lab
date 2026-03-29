# it-doc-semantic-lab

> Research and experimentation repository for the semantic layer of the IT-Dokumentacja project.

## Role of this repository

This is a **lab repository** — not a stable product.  
The stable reference repo is [IT-Dokumentacja](../IT_Dokumentacja/IT-Dokumentacja).

This repo builds the P2→P3 semantic audit capabilities:
- document class detection
- section role mapping
- class-based completeness validation
- gap reporting
- document split planning

## Current project maturity (parent repo)

**P1.5** — structural + safe autofix (V1), pre-semantic.  
See: `IT-Dokumentacja/docs/assessment/PROJECT_CAPABILITY_ASSESSMENT.md`

## Repository structure

```
corpora/           — corpus cases per document class
  <class_id>/
    case_NNN/
      raw/         — original source document
      normalized/  — extracted markdown
      gold/        — gold standard: class, roles, gaps, split plan
      notes/       — rationale and notes
schemas/           — document class and role contract definitions
experiments/       — experiment directories with results
semantic_audit/    — Python modules (stubs → implementations)
reports/           — experiment results
docs/              — development plan and notes
```

## Corpus cases

| Class | Cases | Status |
|---|---|---|
| `nlp_algorithm_architecture_spec` | case_001 | gold standard ready |
| `project_architecture_concept` | — | planned |
| `compliance_procedure` | — | planned |

## Experiments

| ID | Name | Status |
|---|---|---|
| exp_001 | Document class detector | planned |
| exp_002 | Section role mapper | planned |
| exp_003 | Gap report generator | planned |
| exp_004 | Split planner | planned |

## Execution discipline

This repository is experimental, but it is not treated as a loose sandbox.

Work in this repository is expected to follow:
- explicit quality gates,
- evidence-backed execution,
- auditable test layers,
- documented exceptions for skips, degraded modes, and restricted mocks.

The repository does **not** assume failure-free execution.  
It assumes that important failures and deviations must be visible, reproducible, and reviewable.

See:
- `docs/EXECUTION_ASSURANCE_PROGRAM.md`
- `docs/QUALITY_GATES_POLICY.md`
- `docs/TESTING_STANDARD.md`
- `docs/TEST_CATALOG.md`

## Documentation

- `docs/ARCHITECTURE.md` — overview of the Go CLI architecture and capability slices
- `docs/TESTING_STANDARD.md` — testing philosophy, 6 levels, mandatory rules and conventions
- `docs/TEST_CATALOG.md` — 30-layer operational test catalog with gate, mock policy and evidence strength
- `docs/EVIDENCE_MODEL.md` — required artifacts per run, mandatory files, INCOMPLETE run criteria
- `docs/QUALITY_GATES.md` — blocking gate conditions (itdlab CLI)
- `docs/QUALITY_GATES_POLICY.md` — gate enforcement policy and promotion rules
- `docs/EXECUTION_ASSURANCE_PROGRAM.md` — execution assurance program stages
- `docs/POLICY_SKIPS_AND_EXCEPTIONS.md` — when skips are allowed, exception registration, gate impact
- `docs/POLICY_MOCKS_AND_REAL_PATHS.md` — mock-allowed / mock-restricted / mock-forbidden layers
- `docs/REFERENCES.md` — authoritative sources for standards, regulations and documentation practices
- `docs/PLAYBOOKS/` — strategic how-to guides
- `docs/RUNBOOKS/` — step-by-step operational procedures
- `docs/ADR/` — architectural decision records

## Relation to stable repo

This lab → promotes to → IT-Dokumentacja/itdoc/

Promotion conditions: experiment has gold standard, passes on all corpus cases, idempotency confirmed, IT-Dokumentacja validates green after integration.
