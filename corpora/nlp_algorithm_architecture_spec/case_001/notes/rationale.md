# Case 001 — morfologia_polska_algorytmy

## Source
Real technical document: architecture spec for Polish morphology algorithm network.

## Why this case
- First diagnostic catalyst for project maturity assessment
- Used to establish P1.5 baseline in controlled capability test
- Reveals gap between structural detection (working) and semantic role mapping (not yet built)

## Document structure
6 top-level chapters:
1. Domain taxonomy (L1-L5 morphological layers)
2. Architecture overview + probabilistic model
3. Module specification (MOD-1..MOD-6)
4. Output contract (rule vector format)
5. Data model (SQLite schema)
6. Summary + next steps

## Gold standard assessment result
- Detected class: nlp_algorithm_architecture_spec ✓
- Roles present: document_goal, domain_taxonomy, architecture_overview, module_specification, output_contract, data_model (all complete)
- Roles partial: implementation_plan, test_strategy
- Roles missing: risk_register, decision_log, acceptance_criteria, mvp_scope, ownership_model

## Project maturity level at time of test
P1.5 — structural + safe autofix, pre-semantic
