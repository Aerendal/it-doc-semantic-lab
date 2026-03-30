# Evidence Model

Każde ukończone uruchomienie `itdlab` musi produkować weryfikowalny **evidence pack** — zestaw artefaktów umożliwiających niezależny audyt tego, co się stało, dlaczego i z jakim wynikiem.

Uruchomienie bez kompletnego evidence pack jest traktowane jako **INCOMPLETE run** niezależnie od kodu wyjścia.

---

## Zasady

1. Dowody muszą być produkowane przez samo uruchomienie — nie montowane ręcznie po fakcie.
2. Dowody muszą być wystarczające dla niezależnego recenzenta, aby odtworzyć przebieg zdarzeń bez kontaktu z autorem.
3. Dowody muszą przetrwać sesję — artefakty zapisane wyłącznie w pamięci lub katalogach tymczasowych nie są liczone.
4. Zielony kod wyjścia bez dowodów nie jest wiarygodnym wynikiem.

---

## Wymagane artefakty na uruchomienie

Każde uruchomienie dowolnego polecenia `itdlab` musi produkować **wszystkie** poniższe artefakty:

| # | Artefakt | Lokalizacja | Format | Opis |
|---|----------|-------------|--------|------|
| 1 | Rekord uruchomienia | tabela SQLite `runs` | wiersz | czas rozpoczęcia, czas zakończenia, kod wyjścia, status |
| 2 | Wpisy dziennika zdarzeń | `runs/events.jsonl` | JSONL | Wszystkie zdarzenia dołączone podczas tego uruchomienia |
| 3 | Wyjście polecenia | `reports/<run_id>/stdout.txt` | tekst | Przechwycone wyjście standardowe CLI |
| 4 | Kod wyjścia | Rekord uruchomienia w SQLite | liczba całkowita | 0 = sukces, wartość niezerowa = błąd |
| 5 | Suma kontrolna SQLite | `reports/<run_id>/db_checksum.txt` | tekst | SHA-256 pliku bazy danych po uruchomieniu |
| 6 | Podsumowanie uruchomienia | `reports/<run_id>/summary.md` | Markdown | Czytelne dla człowieka podsumowanie uruchomienia |

Artefakty 1–6 są obowiązkowe dla **wszystkich** uruchomień.

---

## Artefakty wymagane dla poszczególnych poleceń

Poza zestawem uniwersalnym każde polecenie wymaga własnych artefaktów:

| Polecenie | Wymagany artefakt | Lokalizacja |
|-----------|-------------------|-------------|
| `ingest run` | Manifest źródłowy | `reports/<run_id>/source_manifest.json` |
| `ingest run` | Raport parsowania | `reports/<run_id>/parse_report.json` |
| `normalize apply` | Raport normalizacji | `reports/<run_id>/normalization_report.json` |
| `normalize apply` | Raport kolizji | `reports/<run_id>/collision_report.json` |
| `relations show` | Graf relacji | `reports/<run_id>/relation_graph.json` |
| `authority check` | Raport pokrycia autorytetu | `reports/<run_id>/authority_coverage_report.json` |
| `export repo1` | Manifest eksportu | `reports/<run_id>/export_manifest.json` |
| `export repo1` | Raport przejścia bramki | `reports/<run_id>/gate_pass_report.json` |
| `audit evidence` | Manifest dowodów | `reports/<run_id>/evidence_manifest.json` |

---

## Definicja uruchomienia INCOMPLETE

Uruchomienie jest uznawane za **INCOMPLETE** jeśli zachodzi którykolwiek z poniższych warunków:

1. Uniwersalny zestaw artefaktów (pozycje 1–6) nie jest w pełni obecny.
2. Brakuje artefaktu wymaganego przez dane polecenie.
3. Jakiś artefakt istnieje, ale jest pusty (zero bajtów).
4. Tabela SQLite `runs` nie zawiera wiersza dla tego `run_id`.
5. Dziennik zdarzeń nie zawiera wpisów dla tego `run_id`.
6. Rekord uruchomienia ma `status = 'running'`, a proces nie jest już aktywny.

Uruchomienie INCOMPLETE nie może być cytowane jako dowód przy ocenie bramki.

---

## Weryfikacja kompletności evidence pack

Kompletność evidence pack jest weryfikowana przez:
- **Layer 27** (Evidence Pack Tests) w `docs/TEST_CATALOG.md`
- **Gate 5** w `docs/QUALITY_GATES.md`

Polecenie `itdlab audit evidence <run_id>` przeprowadza weryfikację kompletności i kończy się kodem niezerowym, jeśli brakuje jakiegoś artefaktu.

---

## Format dziennika zdarzeń

Każda linia w `runs/events.jsonl` jest obiektem JSON. Dziennik jest wyłącznie dołączany i nie może być modyfikowany po zapisaniu wiersza.

```json
{
  "ts": "2026-03-29T12:00:00Z",
  "run_id": "run_001",
  "step": "normalize",
  "entity": "document_family",
  "entity_id": "risk_register",
  "action": "canonicalized",
  "before": "Risk Register",
  "after": "risk_register",
  "meta": {}
}
```

**Wymagane pola:** `ts`, `run_id`, `step`, `entity`, `entity_id`, `action`

**Zabronione operacje:** usuwanie, modyfikacja, obcinanie, zmiana kolejności istniejących linii

---

## Suma kontrolna SQLite

Artefakt sumy kontrolnej (`db_checksum.txt`) musi zawierać:

```
sha256:<hex>  db/semantic_index.sqlite
```

Jest produkowana bezpośrednio po ostatnim zapisie uruchomienia, przed zakończeniem procesu. Obejmuje pełny plik bazy danych, włącznie ze stanem scalonym WAL.

---

## Format podsumowania uruchomienia

`reports/<run_id>/summary.md` musi zawierać co najmniej:

- Run ID
- Wywołane polecenie (z flagami)
- Czasy rozpoczęcia i zakończenia
- Kod wyjścia
- Liczba przetworzonych encji (według typu)
- Napotkane błędy lub ostrzeżenia
- Status bramki (pass / fail / not evaluated)

---

## Gwarancja odtwarzalności

Przy tych samych plikach źródłowych i tej samej konfiguracji uruchomienia, evidence pack (w tym stan SQLite i wpisy dziennika zdarzeń) musi być funkcjonalnie identyczny między uruchomieniami.

„Funkcjonalnie identyczny" oznacza:
- ta sama liczba wierszy we wszystkich tabelach,
- te same przypisane kanoniczne identyfikatory,
- te same wywnioskowane relacje,
- te same wyniki bramek.

Znaczniki czasu i identyfikatory uruchomień są zwolnione z wymogu odtwarzalności.

Każda funkcja, która narusza gwarancję funkcjonalnej odtwarzalności, musi być wyraźnie oznaczona w kodzie komentarzem: `// non-deterministic: <reason>`.
