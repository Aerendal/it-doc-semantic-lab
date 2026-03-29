# Playbook: Export to Repo 1

## Purpose

Defines the strategy for promoting stable, validated semantic metadata from this lab to the reference repository (`IT-Dokumentacja`).

---

## What Is Exported?

Only **stable, gate-validated** metadata is promoted. Never raw or intermediate state.

Exported artifacts:
- Document canonical IDs and classes
- Section role mappings (stable ones only)
- Relation graph (validated edges only)
- Authority coverage report

---

## Pre-Export Requirements

All of the following must be true before export:

1. Gate 1–5 have all passed for the current run
2. All relations have non-empty explanations (Layer 18)
3. Evidence pack is complete (Layer 27)
4. No unresolved canonical ID collisions (Gate 3)

If any condition is not met, `itdlab export repo1` exits with code `2`.

---

## Export Process

```
itdlab export repo1 --target ../IT-Dokumentacja/
```

This:
1. Validates all quality gates
2. Generates export artifacts in `normalized/`
3. Copies stable files to the target repository path
4. Writes `export_manifest.json` listing all promoted files
5. Appends `exported` events to JSONL log

---

## What Is NOT Exported

- Raw source files
- Intermediate normalization state
- Experimental or low-confidence relations (confidence < 0.8)
- SQLite database (stays in lab only)
- Event log (stays in lab only)

---

## Rollback

Export is not destructive to the lab. To undo a promotion to repo 1:
- Use `git revert` in the reference repository
- The lab state is unaffected
