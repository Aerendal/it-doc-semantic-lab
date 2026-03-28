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

## Relation to stable repo

This lab → promotes to → IT-Dokumentacja/itdoc/

Promotion conditions: experiment has gold standard, passes on all corpus cases, idempotency confirmed, IT-Dokumentacja validates green after integration.
