# Playbook: Normalization

## Purpose

Defines the strategy for normalizing document names and assigning canonical IDs.

---

## Normalization Rules

### 1. Canonical ID Generation

A canonical ID is derived from the raw document name by:
1. Lowercasing all characters
2. Replacing spaces, hyphens, and slashes with underscores
3. Stripping non-alphanumeric characters (except underscores)
4. Collapsing multiple underscores into one

Examples:
- `"Risk Register"` → `risk_register`
- `"IT Security Policy v2"` → `it_security_policy_v2`
- `"Change Management Process (CMP)"` → `change_management_process_cmp`

### 2. Collision Handling

If two different documents produce the same canonical ID:
1. A collision entry is written to SQLite
2. The run is NOT blocked — normalization continues
3. Gate 3 blocks promotion until all collisions are resolved

Resolution options:
- Rename one of the documents
- Add a distinguishing suffix (e.g., `_v2`, `_legacy`)
- Merge if the documents are truly duplicates

### 3. Alias Registration

A document may have multiple known aliases. Aliases:
- Are stored in the `normalizations` table
- Always resolve to exactly one canonical ID
- Cannot be shared between canonical IDs

---

## Preview Before Apply

Always preview before applying:

```
itdlab normalize preview
```

Review the `normalization_report.json` output. Only apply when all changes are intentional:

```
itdlab normalize apply
```

---

## When to Re-Normalize

- After adding new source documents
- After renaming source files
- After changing normalization rules (requires migration test — Layer 15)
