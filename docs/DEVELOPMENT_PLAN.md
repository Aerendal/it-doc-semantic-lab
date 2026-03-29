# Development Plan — it-doc-semantic-lab

## Purpose
This repository is the semantic research and experimentation lab for the IT-Dokumentacja project.

The stable production repo (IT-Dokumentacja) contains P1.5 capabilities.
This lab builds the P2→P3 semantic layer.

## Layer targets

### P2 — structural-rule
- [ ] Document class detector (exp_001)
- [ ] Section role mapper (exp_002)
- [ ] Class-role registry (schemas/document_classes.yaml)
- [ ] Role completeness evaluator

### P3 — semantic-role
- [ ] Gap report generator (exp_003)
- [ ] Split planner (exp_004)
- [ ] Content completeness scorer

## Corpus first cases
1. nlp_algorithm_architecture_spec → case_001 (morfologia_polska_algorytmy) ✓
2. project_architecture_concept → case_001 (to be added)
3. compliance_procedure → case_001 (to be added)

## Execution-assurance dependency

P2→P3 implementation work is not treated as complete only when code exists.

Each capability slice should be evaluated against:
- quality gates,
- applicable test layers,
- evidence requirements,
- real-path vs mocked-path rules,
- promotion-readiness criteria.

Implementation order should therefore track not only feature completion, but also:
1. test-layer coverage,
2. evidence-pack completeness,
3. gate status,
4. promotion feasibility to the stable repository.

See: `docs/EXECUTION_ASSURANCE_PROGRAM.md`, `docs/QUALITY_GATES_POLICY.md`, `docs/TESTING_STANDARD.md`.

## Promotion policy
Stable components are promoted to IT-Dokumentacja/itdoc/ only when:
1. Experiment has a gold standard
2. Experiment passes on all corpus cases
3. Dry-run/apply/idempotency confirmed
4. Validation passes on IT-Dokumentacja after integration
