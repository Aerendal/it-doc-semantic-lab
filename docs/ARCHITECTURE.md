# Architecture — IT Documentation Semantic Lab

## Przegląd

To repozytorium jest **silnikiem analizy semantycznej** dla projektu IT-Dokumentacja.  
Działa jako narzędzie CLI w Go, wspierane przez SQLite (źródło prawdy) i dziennik zdarzeń JSONL (ścieżka audytowa).

Stabilne repozytorium referencyjne (`IT-Dokumentacja`) otrzymuje wyłącznie **promowane, stabilne metadane** eksportowane stąd.

---

## Zasady projektowe

| Zasada | Wyraz |
|--------|-------|
| Local-first | Brak zależności sieciowych w czasie wykonania |
| Auditable | Każda zmiana stanu jest logowana do `runs/events.jsonl` |
| Reproducible | Te same dane wejściowe + ten sam przebieg = te same dane wyjściowe |
| Incremental | Wycinki możliwości, nie wielkie warstwy naraz |
| Evidence-driven | Każdy przebieg produkuje weryfikowalną paczkę dowodową |

---

## Stos technologiczny

| Obszar | Wybór | Uzasadnienie |
|--------|-------|-------------|
| Język | Go | Statyczny plik binarny, silna biblioteka standardowa, łatwe testowanie |
| Baza danych | SQLite (`modernc.org/sqlite`) | Lokalna, bez CGO, zero konfiguracji |
| Dziennik audytu | JSONL (tylko dopisywanie) | Czytelny dla człowieka, grepowany, bezpieczny dla dopisywania |
| Framework CLI | Cobra | Drzewo podpoleceń, flagi, generowanie pomocy |
| Narzędzia testowe | stdlib `testing` + golden files | Brak zewnętrznych zależności |

---

## Struktura repozytorium

```
cmd/itdlab/          — CLI entrypoint
internal/
  app/               — application-layer use cases (one package per capability)
    ingest/
    normalize/
    classify/
    relations/
    sections/
    authority/
    export/
    audit/
  domain/            — pure domain types (no I/O)
  ports/             — interfaces (source reader, event log, stores, report writer)
  adapters/
    sqlite/          — SQLite implementation of stores
    jsonl/           — JSONL event log implementation
    filesystem/      — filesystem source reader
    markdown/        — Markdown parser
  cli/               — Cobra command definitions
  testkit/           — test helpers, fixtures, golden files, builders
db/
  schema_v1.sql      — canonical DDL
  semantic_index.sqlite  — runtime database (gitignored)
runs/
  events.jsonl       — append-only event log (gitignored)
sources/             — raw IT documentation inputs
normalized/          — normalized outputs
reports/             — generated reports per run
docs/
  ARCHITECTURE.md    — this file
  PLAYBOOKS/         — strategic how-to guides
  RUNBOOKS/          — step-by-step operational procedures
  ADR/               — architectural decision records
  TESTING_STANDARD.md
  TEST_CATALOG.md
  EVIDENCE_MODEL.md
  QUALITY_GATES.md
```

---

## Wycinki możliwości

Rozwój odbywa się według pionowych wycinków możliwości, nie poziomych warstw.

| Wycinek | Tło | Interfejs | Dowody |
|---------|-----|-----------|--------|
| 1 — Ingest | Parser Markdown, magazyn SQLite, dziennik zdarzeń | `itdlab ingest run`, `itdlab ingest inspect` | Manifest źródłowy, raport parsowania |
| 2 — Normalize | Kanoniczne ID, deduplikacja, wykrywanie kolizji | `itdlab normalize preview`, `itdlab normalize apply` | Raport normalizacji |
| 3 — Sections | Archetypy sekcji, wnioskowanie ról | `itdlab sections show`, `itdlab sections explain` | Mapa sekcji, raport anomalii |
| 4 — Relations | Wnioskowanie oparte na regułach, zależności między dokumentami | `itdlab relations show`, `itdlab relations explain` | Graf relacji, raport kandydatów |
| 5 — Authority | Powiązania regulacyjne | `itdlab authority check` | Raport pokrycia przez autorytet |
| 6 — Export | Promocja stabilnych metadanych | `itdlab export repo1` | Manifest eksportu |

---

## Model stanu

```
raw → ingested → normalized → classified → exported
```

Każde przejście jest:
1. Zapisywane jako mutacja wiersza w SQLite
2. Dopisywane jako zdarzenie do `runs/events.jsonl`

---

## Powiązane decyzje

- [ADR-001: SQLite as source of truth](ADR/ADR-001-sqlite-as-source-of-truth.md)
