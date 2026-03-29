# Playbook: Authority Alignment

## Purpose

Defines the strategy for linking IT documentation artifacts to regulatory and standards authorities.

---

## What Is an Authority Reference?

An authority reference (`authority_ref`) links a document or section to a specific clause in an external authority:
- Regulatory frameworks (e.g., ISO 27001, GDPR, 21 CFR Part 11)
- Industry standards (e.g., ITIL, COBIT)
- Internal policies (treated as internal authorities)

---

## When to Create Authority Refs

1. A document explicitly cites a regulation or standard in its text
2. A section heading or content maps to a known clause
3. An industry matrix requires specific documents for compliance

---

## Authority Alignment Process

### 1. Define authorities

Authorities are referenced by name + clause, e.g.:
- `authority: "ISO 27001"`, `clause: "A.12.6.1"`
- `authority: "GDPR"`, `clause: "Art. 32"`

### 2. Link documents

```
itdlab authority check
```

This:
1. Scans all ingested documents for authority patterns
2. Creates `authority_ref` rows for detected matches
3. Reports coverage per authority in `authority_coverage_report.json`

### 3. Review coverage

A document set's authority coverage is the ratio of linked clauses to total clauses for that authority. Target: ≥ 80% for any required authority.

---

## Failure Modes

| Symptom | Action |
|---------|--------|
| Low coverage for a required authority | Review documents for missing citations |
| Duplicate authority refs | Check for re-runs without deduplication |
| Authority name mismatch | Standardize authority names before linking |
