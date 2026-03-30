# Run Manifest Schema

Manifest uruchomienia (`reports/<run_id>/run_manifest.json`) jest czytelnym maszynowo podsumowaniem tożsamości, zakresu, wyniku i dowodów uruchomienia. Jest produkowany po zakończeniu uruchomienia i jest wymagany dla wszystkich uruchomień `trusted`.

---

## Cel

Manifest uruchomienia jest jedynym autorytatywnym indeksem evidence pack uruchomienia. Recenzent lub narzędzie automatyczne może odczytać manifest, aby ustalić:
- jakie polecenie zostało uruchomione i z jakimi parametrami,
- jaki był wynik,
- które artefakty zostały wyprodukowane,
- czy evidence pack jest kompletny,
- które bramki jakości zostały ocenione i jaki był ich status.

---

## Schemat

```json
{
  "schema_version": 1,
  "run_id": "<string — unique run identifier>",
  "command": "<string — full CLI command as invoked, e.g. 'itdlab ingest run --source sources/'>",
  "started_at": "<string — ISO 8601 datetime>",
  "finished_at": "<string — ISO 8601 datetime>",
  "exit_code": "<integer — 0 | 1 | 2 | 3>",
  "status": "<string — 'completed' | 'failed' | 'aborted'>",
  "trusted": "<boolean — true if all trusted-run criteria are met>",
  "environment": {
    "binary_version": "<string — itdlab version>",
    "go_version": "<string — Go runtime version>",
    "os": "<string — operating system>",
    "db_path": "<string — path to SQLite file used>",
    "log_path": "<string — path to event log file used>"
  },
  "evidence": {
    "db_checksum": "<string — sha256:<hex>>",
    "event_count": "<integer — number of events appended in this run>",
    "artifacts": [
      {
        "name": "<string — artifact logical name, e.g. 'parse_report'>",
        "path": "<string — relative path from repo root>",
        "size_bytes": "<integer>",
        "sha256": "<string — sha256:<hex>>"
      }
    ],
    "complete": "<boolean — true if all required artifacts are present and non-empty>"
  },
  "gates": [
    {
      "gate_id": "<string — e.g. 'G1'>",
      "status": "<string — 'PASS' | 'degraded' | 'FAIL' | 'not_evaluated'>",
      "evaluated_at": "<string — ISO 8601 datetime>",
      "failures": [
        "<string — description of each failing condition>"
      ]
    }
  ],
  "skips": [
    {
      "skip_id": "<string — from SKIP_REGISTER.md>",
      "layer": "<string — e.g. 'Layer 28 — Performance Budget Tests'>",
      "category": "<integer — 1 | 2 | 3>",
      "gate_impact": "<string — 'none' | 'degraded:G4'>"
    }
  ],
  "entities_processed": {
    "documents": "<integer>",
    "sections": "<integer>",
    "relations": "<integer>",
    "normalizations": "<integer>"
  },
  "errors": [
    {
      "step": "<string>",
      "entity_id": "<string>",
      "message": "<string>"
    }
  ]
}
```

---

## Definicje pól

### Pola najwyższego poziomu

| Pole | Typ | Wymagane | Opis |
|------|-----|----------|------|
| `schema_version` | liczba całkowita | tak | Zawsze `1` dla tej wersji schematu |
| `run_id` | ciąg znaków | tak | Unikalny identyfikator uruchomienia. Musi odpowiadać SQLite `runs.run_id` |
| `command` | ciąg znaków | tak | Pełny ciąg wywołania CLI |
| `started_at` | ciąg znaków | tak | Data i czas UTC w formacie ISO 8601 |
| `finished_at` | ciąg znaków | tak | Data i czas UTC w formacie ISO 8601 |
| `exit_code` | liczba całkowita | tak | 0, 1, 2 lub 3 zgodnie z kontraktem kodu wyjścia |
| `status` | ciąg znaków | tak | `completed`, `failed` lub `aborted` |
| `trusted` | wartość logiczna | tak | `true` tylko jeśli wszystkie kryteria uruchomienia `trusted` są spełnione |

### Obiekt `evidence`

| Pole | Typ | Wymagane | Opis |
|------|-----|----------|------|
| `db_checksum` | ciąg znaków | tak | SHA-256 pliku SQLite na końcu uruchomienia. Format: `sha256:<hex>` |
| `event_count` | liczba całkowita | tak | Liczba linii dziennika zdarzeń dołączonych w tym uruchomieniu |
| `artifacts` | tablica | tak | Jeden wpis na każdy wyprodukowany artefakt |
| `complete` | wartość logiczna | tak | `true` jeśli wszystkie wymagane artefakty są obecne i niepuste |

### Tablica `gates`

Każdy element opisuje jedną ocenę bramki. Bramki nieistotne dla polecenia są rejestrowane ze `status: "not_evaluated"`.

### Tablica `skips`

Każde aktywne pominięcie z `SKIP_REGISTER.md`, które wpłynęło na to uruchomienie. Pusta tablica, jeśli żadne pominięcia nie były aktywne.

### Tablica `errors`

Niekrytyczne błędy napotkane podczas uruchomienia. Uruchomienie może zakończyć się kodem 0 z niepustą tablicą `errors`, jeśli błędy dotyczyły poszczególnych encji i nie wpłynęły na ogólny wynik.

---

## Reguły walidacji

1. `run_id` musi odpowiadać wierszowi w tabeli SQLite `runs` z tym samym `run_id`.
2. `db_checksum` musi odpowiadać rzeczywistej wartości SHA-256 pliku bazy danych w momencie zapisania manifestu.
3. `evidence.complete` musi być `false`, jeśli jakikolwiek artefakt w tablicy `artifacts` ma `size_bytes = 0`.
4. `trusted` musi być `false`, jeśli `evidence.complete = false`.
5. `trusted` musi być `false`, jeśli jakiekolwiek pominięcie w `skips` ma `category = 3` lub `category = 4`.
6. Status bramki `PASS` wymaga `evidence.complete = true`.

---

## Produkcja

Manifest jest produkowany przez `itdlab` na końcu każdego uruchomienia zmieniającego stan, zapisywany do:

```
reports/<run_id>/run_manifest.json
```

Jest też odwoływany przez `itdlab audit evidence <run_id>` do weryfikacji kompletności.

---

## Wewnętrzne odniesienia
- `docs/EXECUTION_CONTRACT.md`
- `docs/EVIDENCE_MODEL.md`
- `docs/CONTEXT_VOCABULARY.md`
- `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`

## Metadane przeglądu
- Właściciel: zespół projektowy
- Status: szkic
- Last reviewed: 2026-03-30
