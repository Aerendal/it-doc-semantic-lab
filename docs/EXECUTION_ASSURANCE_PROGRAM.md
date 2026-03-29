# Execution Assurance Program

## Purpose

This document defines the execution-assurance program for the experimental repository.

The goal is not only to add new semantic capabilities, but to build a development and verification system that reduces the risk of:
- silently skipped tests,
- unjustified mock-based green runs,
- incomplete execution paths,
- undocumented shortcuts,
- regressions hidden behind weak evidence,
- architecture drift between the lab repository and the stable repository.

This program is designed for an **auditable experimental repository**.  
The repository may remain experimental, but the execution model, evidence model, and test discipline must be explicit and reviewable.

---

## Scope

The program applies to:
- source ingestion,
- normalization,
- document classification,
- section-role mapping,
- relation inference,
- authority alignment,
- export/promotion to the stable repository,
- CLI workflows,
- run manifests,
- evidence packs,
- test design and test execution.

---

## Core engineering principles

1. **No silent skipping**  
   A skipped test or skipped execution step must be explicit, justified, and reportable.

2. **No hidden mocks for critical paths**  
   Critical runtime or data-path behavior must not be validated only through mocks when a real-path verification is required.

3. **Evidence before trust**  
   A green run is not treated as trustworthy unless it produces a manifest, logs, and traceable evidence.

4. **Determinism first**  
   The same inputs should produce the same outputs or a clearly explained delta.

5. **Promotion only after proof**  
   Components are promoted to the stable repository only after satisfying defined technical and evidence gates.

6. **Experimentation is allowed, but not undocumented improvisation**  
   Every important deviation, assumption, or simplification must be documented.

---

# Program structure

The execution-assurance program is divided into **4 phases** and **10 implementation stages**.

## Phase A — foundations
Stages 1–3 define the rules of execution and the evidence contract.

## Phase B — defensive mechanics
Stages 4–6 build mechanisms that prevent weak validation practices.

## Phase C — enforcement and operations
Stages 7–9 connect rules with tooling, documentation, and audit evidence.

## Phase D — rollout and hardening
Stage 10 integrates the program into normal repository work.

---

# 10 implementation stages

## Stage 1 — Gate model and enforcement policy

### Objective
Define what is allowed, forbidden, and conditionally allowed in repository execution and testing.

### Inputs
- current repository structure,
- current testing workflow,
- current experiment plan,
- promotion policy to the stable repository.

### Outputs
- gate model,
- policy for skips,
- policy for mocks/fakes/stubs,
- policy for evidence generation,
- policy for promotion decisions.

### Required artifacts
- `docs/QUALITY_GATES.md`
- `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`
- `docs/POLICY_MOCKS_AND_REAL_PATHS.md`

### Closure criteria
- forbidden and allowed patterns are explicit,
- repository can classify a run as valid / invalid / incomplete,
- critical-path definition exists.

---

## Stage 2 — Test-layer catalog and risk mapping

### Objective
Define the full map of test layers and what each one is supposed to detect.

### Inputs
- semantic pipeline goals,
- known failure modes,
- planned relation model,
- planned ingest / normalization / export flows.

### Outputs
- test-layer catalog,
- risk-to-test mapping,
- justification for each test layer.

### Required artifacts
- `docs/TESTING_STANDARD.md`
- `docs/TEST_CATALOG.md`
- `docs/RISK_TO_TEST_MATRIX.md`

### Closure criteria
- each test layer has a named purpose,
- each critical risk is mapped to at least one test layer,
- the repository can explain why a given test class exists.

---

## Stage 3 — Execution contract and evidence model

### Objective
Define what counts as a real execution run and what evidence must be generated.

### Inputs
- CLI execution model,
- current reports,
- intended audit requirements.

### Outputs
- run manifest schema,
- evidence pack schema,
- execution contract.

### Required artifacts
- `docs/EVIDENCE_MODEL.md`
- `docs/RUN_MANIFEST_SCHEMA.md`
- `docs/EXECUTION_CONTRACT.md`

### Closure criteria
- every meaningful run can generate a manifest,
- manifests include executed steps, skipped steps, configuration, and outcome,
- evidence packs are reproducible and reviewable.

---

## Stage 4 — Anti-skip mechanisms

### Objective
Prevent silent weakening of the test suite through unjustified skips or soft bypasses.

### Inputs
- existing tests,
- marker usage,
- current skip/xfail patterns.

### Outputs
- skip/xfail linting,
- skip whitelist policy,
- delta report for newly introduced skips.

### Required artifacts
- test linter rules,
- skip audit report,
- exception register.

### Closure criteria
- new skips are visible,
- unjustified skips fail gates,
- skip usage can be audited historically.

---

## Stage 5 — Anti-mock mechanisms

### Objective
Prevent fake green runs caused by hidden mocking of critical paths.

### Inputs
- current adapters,
- current test style,
- critical-path definition from Stage 1.

### Outputs
- mock policy by layer,
- real-path vs simulated-path classification,
- detection rules for disallowed patching.

### Required artifacts
- `docs/POLICY_MOCKS_AND_REAL_PATHS.md`
- critical-path registry,
- mock-usage audit report.

### Closure criteria
- critical paths have explicit mock rules,
- disallowed mocking is detectable,
- repository can distinguish real verification from simulated verification.

---

## Stage 6 — Deterministic fixtures and real test harness

### Objective
Create repeatable, audited fixtures and goldens that reduce pressure to “just mock it”.

### Inputs
- source corpora,
- normalized outputs,
- gold cases,
- relation candidates,
- authority references.

### Outputs
- deterministic fixtures,
- golden outputs,
- negative cases,
- corruption cases,
- boundary cases.

### Required artifacts
- `testdata/fixtures/`
- `testdata/golden/`
- `docs/FIXTURE_POLICY.md`

### Closure criteria
- major flows have stable fixtures,
- expected outputs are versioned,
- fixture selection is explicit and reproducible.

---

## Stage 7 — CLI and CI quality gates

### Objective
Connect the assurance model to executable commands and automated gates.

### Inputs
- outputs of Stages 1–6,
- repository CLI design,
- CI workflow strategy.

### Outputs
- verification commands,
- gate command(s),
- CI pipeline integration,
- fail/warn/pass semantics.

### Required artifacts
- `make verify`
- `make audit-run`
- CI workflow definitions
- gate summary report

### Closure criteria
- repository can run gate checks consistently,
- CI reflects real gate status,
- failures are actionable and interpretable.

---

## Stage 8 — Playbooks and runbooks

### Objective
Document how work should be performed, repeated, and audited.

### Inputs
- actual implemented workflows,
- evidence contract,
- gate model.

### Outputs
- playbooks,
- runbooks,
- operator procedures,
- troubleshooting logic.

### Required artifacts
- `docs/PLAYBOOKS/`
- `docs/RUNBOOKS/`
- `docs/TROUBLESHOOTING.md`

### Closure criteria
- major workflows are documented step by step,
- a reviewer can reproduce a run without oral explanations,
- operational ambiguity is reduced.

---

## Stage 9 — Audit layer and evidence packs

### Objective
Ensure that each important run leaves a verifiable trace.

### Inputs
- run manifest model,
- gate execution,
- fixture and report outputs.

### Outputs
- evidence packs,
- run summaries,
- checksums/fingerprints,
- audit reports.

### Required artifacts
- `runs/`
- `reports/`
- `evidence/`
- standardized audit summary files

### Closure criteria
- each important run has a reviewable evidence pack,
- evidence is sufficient to explain PASS/FAIL/WARN,
- execution traces can be compared across runs.

---

## Stage 10 — Rollout and hardening

### Objective
Apply the assurance system to day-to-day work and close the remaining bypasses.

### Inputs
- all prior stages,
- real development usage,
- observed failure/bypass patterns.

### Outputs
- hardened gates,
- reduced bypass surface,
- stable promotion policy,
- updated operating discipline.

### Required artifacts
- rollout checklist,
- hardening backlog,
- post-rollout review report.

### Closure criteria
- the assurance model is used in normal development,
- bypasses are reduced and visible,
- promotion to the stable repository is evidence-driven.

---

# 30 test/control layers

The program uses **30 test/control layers**, grouped into 6 levels.

## Level A — source and input contract (1–5)
1. file presence tests  
2. file readability tests  
3. encoding tests  
4. markdown structure tests  
5. source schema tests

## Level B — parser and extraction (6–10)
6. parser unit tests  
7. fixture parsing tests  
8. golden extraction tests  
9. partial corruption tests  
10. determinism tests

## Level C — normalization and canonical model (11–15)
11. canonical ID tests  
12. collision detection tests  
13. alias resolution tests  
14. typing tests  
15. migration tests

## Level D — relations and semantic logic (16–20)
16. relation rule unit tests  
17. relation consistency tests  
18. explainability tests  
19. cycle / acyclicity tests  
20. section influence tests

## Level E — interfaces and execution flows (21–25)
21. CLI contract tests  
22. end-to-end slice tests  
23. resume / restart tests  
24. event log integrity tests  
25. SQLite materialization tests

## Level F — auditability and release control (26–30)
26. reproducibility tests  
27. evidence pack tests  
28. performance budget tests  
29. failure-mode tests  
30. release gate tests

---

# Expected repository documentation

At minimum, the repository should eventually contain:

- `docs/QUALITY_GATES.md`
- `docs/TESTING_STANDARD.md`
- `docs/TEST_CATALOG.md`
- `docs/EVIDENCE_MODEL.md`
- `docs/EXECUTION_CONTRACT.md`
- `docs/PLAYBOOKS/`
- `docs/RUNBOOKS/`
- `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`
- `docs/POLICY_MOCKS_AND_REAL_PATHS.md`
- `docs/RISK_TO_TEST_MATRIX.md`
- `docs/FIXTURE_POLICY.md`
- `docs/ADR/ADR-001-sqlite-as-source-of-truth.md`

---

# Immediate implementation order

The recommended order of work is:

1. Stage 1 — gate model and enforcement policy
2. Stage 2 — test-layer catalog and risk mapping
3. Stage 3 — execution contract and evidence model
4. Stage 4 — anti-skip mechanisms
5. Stage 5 — anti-mock mechanisms
6. Stage 6 — deterministic fixtures and real harness
7. Stage 7 — CLI and CI gates
8. Stage 8 — playbooks and runbooks
9. Stage 9 — audit layer and evidence packs
10. Stage 10 — rollout and hardening

---

# Relation to the experimental repository

This program is intended for the **experimental repository**, where semantic capabilities are researched and implemented under controlled conditions.

The stable repository remains the promotion target.  
The experimental repository is where:
- new semantic mechanisms are designed,
- the execution discipline is hardened,
- evidence models are developed,
- verification logic is validated before promotion.

---

# Final note

The purpose of this program is not to eliminate experimentation.
The purpose is to ensure that experimentation is:
- explicit,
- reviewable,
- repeatable,
- evidence-backed,
- and resistant to silent weakening of quality standards.
