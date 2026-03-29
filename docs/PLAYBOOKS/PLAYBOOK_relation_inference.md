# Playbook: Relation Inference

## Purpose

Defines the strategy for inferring semantic relations between IT documentation artifacts.

---

## Relation Types

| Type | Meaning |
|------|---------|
| `depends_on` | Document A requires Document B to exist and be current |
| `implements` | Document A gives concrete form to a policy defined in B |
| `references` | Document A cites Document B |
| `supersedes` | Document A replaces Document B |
| `complements` | Document A and B together cover a shared concern |
| `derived_from` | Document A is a derivative or specialization of B |

---

## Rule-Based Inference

Relations are inferred by rules stored in the `relation_rules` table. A rule has:
- `from_class`: document class of the source
- `to_class`: document class of the target
- `rel_type`: relation type to infer
- `description`: human-readable explanation

### Example Rule

> A `procedure` document `depends_on` a `policy` document in the same domain.

### Rule Application

1. For each rule, query all document pairs matching `(from_class, to_class)`
2. For each pair, create a `Relation` row with `source = 'rule_engine'`
3. Set `explanation` to the rule's description + matched document names
4. Append event to JSONL log

---

## Explainability Requirement

Every inferred relation must have a non-empty `explanation` field. This is enforced by Layer 18 (Explainability Tests) and Gate 4.

The explanation must be human-readable and include:
- The rule that fired
- The from/to document names
- The confidence score if below 1.0

---

## When to Re-Run Relation Inference

- After new documents are ingested and classified
- After relation rules are modified
- After canonical IDs are updated (re-run normalize first)
