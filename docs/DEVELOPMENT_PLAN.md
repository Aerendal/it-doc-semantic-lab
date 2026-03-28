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

## Promotion policy
Stable components are promoted to IT-Dokumentacja/itdoc/ only when:
1. Experiment has a gold standard
2. Experiment passes on all corpus cases
3. Dry-run/apply/idempotency confirmed
4. Validation passes on IT-Dokumentacja after integration
