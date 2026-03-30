# Execution Contract

Ten dokument definiuje, co każdy przebieg `itdlab` musi gwarantować — warunki wstępne, które musi spełniać, niezmienniki, które musi utrzymywać, oraz warunki końcowe, które musi zapewniać. Przebieg naruszający którykolwiek niezmiennik lub warunek końcowy **nie jest prawidłowym przebiegiem** i nie może być przytaczany jako dowód.

---

## Cel

Kontrakt wykonania jest wiążącą umową między narzędziem CLI a jego operatorami. Nie opisuje *jak* narzędzie działa — opisuje *co* jest gwarantowane jako prawdziwe dla każdego przebiegu, za który narzędzie przyjmuje odpowiedzialność.

---

## Warunki wstępne

Przed rozpoczęciem przebiegu muszą być spełnione następujące warunki. Jeśli którykolwiek z nich nie jest spełniony, przebieg musi odmówić uruchomienia, zakończyć z kodem 1 i zapisać opisowy błąd na stderr.

| # | Warunek wstępny | Sprawdzenie |
|---|-------------|-------|
| P1 | Plik bazy danych SQLite jest dostępny (istnieje lub może zostać utworzony) | Sprawdzane przy starcie |
| P2 | Wersja schematu SQLite ≥ 1 (lub może zostać zainicjowana) | Sprawdzane przy starcie |
| P3 | Ścieżka do dziennika zdarzeń jest zapisywalna (plik istnieje lub może zostać utworzony) | Sprawdzane przy starcie |
| P4 | Katalog raportów jest zapisywalny | Sprawdzane przy starcie |
| P5 | Możliwe jest wygenerowanie unikalnego `run_id` | Sprawdzane przy starcie |
| P6 | Żaden inny przebieg z tym samym `run_id` nie istnieje w SQLite | Sprawdzane przed pierwszym zapisem |

Przebieg, który nie może spełnić P1–P6, jest **przebiegiem odmówionym**. Odmówione przebiegi nie zapisują do SQLite ani do dziennika zdarzeń (z wyjątkiem przyczyny odmowy zapisanej na stderr).

---

## Niezmienniki

Poniższe warunki muszą pozostawać prawdziwe przez cały czas trwania przebiegu — od pierwszego zapisu do ostatniego. Naruszenie któregokolwiek niezmiennika stanowi **naruszenie kontraktu** i musi skutkować natychmiastowym kontrolowanym przerwaniem.

| # | Niezmiennik | Konsekwencja naruszenia |
|---|-----------|-------------------------|
| I1 | Każda operacja zmieniająca stan jest zapisywana do SQLite przed dołączeniem do dziennika zdarzeń | Przerwanie; oznacz przebieg `status = 'aborted'` |
| I2 | Dziennik zdarzeń jest wyłącznie do dołączania; żadna istniejąca linia nie jest modyfikowana ani usuwana | Przerwanie; nie obcinaj dziennika |
| I3 | `run_id` we wszystkich wpisach dziennika zdarzeń odpowiada bieżącemu przebiegowi | Przerwanie w przypadku wykrycia niezgodności |
| I4 | Ograniczenia klucza obcego w SQLite nigdy nie są omijane | Przerwanie; cofnij transakcję |
| I5 | Żaden częściowy zapis nie pozostawia SQLite w niespójnym stanie | Używaj transakcji dla wszystkich wieloetapowych zapisów |
| I6 | Pole `status` przebiegu w SQLite przechodzi tylko do przodu: `running` → `completed` / `failed` / `aborted` | Przerwanie w przypadku próby cofnięcia przejścia |

---

## Warunki końcowe

Po zakończeniu przebiegu (exit 0 lub exit 1) wszystkie poniższe warunki muszą być prawdziwe. Jeśli którykolwiek warunek końcowy nie może zostać spełniony, przebieg musi zapisać `status = 'failed'` i zakończyć z kodem niezerowym.

| # | Warunek końcowy | Weryfikacja |
|---|--------------|--------------|
| Q1 | Rekord przebiegu istnieje w SQLite z ustawionymi `finished_at` i `exit_code` | `SELECT * FROM runs WHERE run_id = ?` |
| Q2 | Co najmniej jeden wpis dziennika zdarzeń istnieje dla przebiegu | `grep run_id events.jsonl` |
| Q3 | `reports/<run_id>/stdout.txt` istnieje i jest niepusty | Sprawdzenie pliku |
| Q4 | `reports/<run_id>/db_checksum.txt` istnieje i zawiera prawidłową linię SHA-256 | Sprawdzenie pliku + formatu |
| Q5 | `reports/<run_id>/summary.md` istnieje i jest niepusty | Sprawdzenie pliku |
| Q6 | Wszystkie wymagane artefakty specyficzne dla polecenia istnieją (zob. `docs/EVIDENCE_MODEL.md`) | Sprawdzenie per polecenie |
| Q7 | SQLite WAL jest zapisany do punktu kontrolnego przed obliczeniem sumy kontrolnej | `PRAGMA wal_checkpoint(FULL)` przed haszowaniem |

Przebieg, który kończy się z kodem 0, lecz nie spełnia Q1–Q7, jest przebiegiem **INCOMPLETE** zgodnie z `docs/EVIDENCE_MODEL.md` i nie może być przytaczany jako dowód bramki.

---

## Protokół naruszenia kontraktu

Gdy wykryto naruszenie kontraktu (naruszenie niezmiennika w trakcie przebiegu):

1. Natychmiast zatrzymaj wszystkie dalsze zapisy.
2. Spróbuj zapisać końcowe zdarzenie `action: "contract_breach"` do dziennika zdarzeń z identyfikatorem niezmiennika i przyczyną.
3. Ustaw `status = 'aborted'` przebiegu w SQLite, jeśli baza danych jest nadal zapisywalna.
4. Zakończ z kodem 3 (naruszenie kontraktu, odrębne od błędu=1 i niepowodzenia bramki=2).
5. Zapisz przyczynę naruszenia na stderr.
6. **Nie** czyść częściowego stanu — zachowaj go do przeglądu forensycznego.

---

## Kontrakt kodów wyjścia

| Code | Znaczenie |
|------|---------|
| 0 | Przebieg zakończony; wszystkie warunki końcowe spełnione |
| 1 | Przebieg nieudany; co najmniej jeden warunek końcowy niespełniony; pakiet dowodowy może być niekompletny |
| 2 | Niepowodzenie bramki; warunki końcowe spełnione, ale co najmniej jedna bramka jakości nie przeszła |
| 3 | Naruszenie kontraktu; niezmiennik naruszony w trakcie przebiegu; częściowy stan zachowany |

Kody wyjścia są częścią kontraktu CLI. Każdy kod produkujący inny kod wyjścia dla tych warunków ma defekt.

---

## Walidacja kontekstu wykonania

Przed pierwszym zapisem stanu narzędzie waliduje kontekst wykonania. Walidacja kontekstu jest odrębna od warunków wstępnych (P1–P6) — sprawdza *jakość* środowiska, nie tylko jego dostępność.

### Wymagane pola kontekstu

| Pole | Opis | Wartości |
|------|------|---------|
| `context_id` | Unikalny identyfikator kontekstu dla tego przebiegu | Format: `<typ>-<środowisko>-<timestamp>` |
| `context_type` | Klasyfikacja kontekstu | `untrusted_local` \| `trusted_ci` \| `authoritative_verify` |
| `db_mode` | Tryb połączenia SQLite | Musi być `file` — tryb `memory` jest odrzucany |
| `schema_version` | Wersja schematu bazy danych | Musi być ≥ 1 |
| `active_skips` | Lista aktywnych pomijań z `SKIP_REGISTER.md` | Sprawdź, czy brak Kategorii 4 lub przeterminowanych |

### Reguły walidacji

1. Jeśli `db_mode = memory`, kontekst jest **odrzucony**. Przebieg kończy się kodem 1.
2. Jeśli jakiekolwiek aktywne pominięcie jest Kategorii 4 lub przeterminowane, kontekst jest **odrzucony** dla tego polecenia.
3. Jeśli `context_type = authoritative_verify` i jakikolwiek layer `mock-forbidden` ma aktywny mock — kontekst jest **odrzucony**.
4. Wynik walidacji (`allowed` | `rejected`) i przyczyny są zapisywane w `run_manifest.json` jako `context_validation_result`.

---

## Maszyna stanów przebiegu

Każdy przebieg przechodzi przez zdefiniowaną sekwencję stanów. Przejścia do przodu są jedyne legalne.

### Stany

| Stan | Opis |
|------|------|
| `created` | Przebieg zainicjowany; `run_id` wygenerowany; rekord w SQLite utworzony |
| `context_validated` | Walidacja kontekstu zakończona wynikiem `allowed` |
| `running` | Operacje zmieniające stan w toku |
| `completed` | Wszystkie warunki końcowe spełnione; exit code 0 lub 2 |
| `failed` | Co najmniej jeden warunek końcowy niespełniony; exit code 1 |
| `aborted` | Niezmiennik naruszony; przebieg zatrzymany awaryjnie; exit code 3 |
| `sealed` | Manifest i artefakty zapieczętowane; nie można modyfikować |
| `promoted` | Przebieg użyty jako podstawa promocji do stabilnego repozytorium |

### Legalne przejścia

```
created → context_validated
created → failed          (walidacja kontekstu odrzucona)
context_validated → running
running → completed
running → failed
running → aborted
completed → sealed        (wyłącznie dla authoritative_verify)
sealed → promoted
```

### Nielegalne przejścia

Każde przejście niespełnione na powyższej liście jest nielegalne. W szczególności:
- `completed → running` — NIEDOZWOLONE (przebieg nie może być wznowiony)
- `aborted → running` — NIEDOZWOLONE (aborted jest stanem terminalnym)
- `failed → completed` — NIEDOZWOLONE
- `sealed → completed` / `sealed → running` — NIEDOZWOLONE

Próba nielegalnego przejścia musi skutkować naruszeniem kontraktu (exit code 3).

---

## Reguły zapieczętowania

Zapieczętowanie (`sealed`) jest opsjonalne i dotyczy wyłącznie przebiegów `authoritative_verify`.

### Warunki zapieczętowania

1. Przebieg jest w stanie `completed` z exit code 0.
2. `evidence.complete = true` w manifeście.
3. `context_type = authoritative_verify`.
4. Brak aktywnych pomijań Kategorii 3 lub 4.
5. Wszystkie wymagane bramki mają status `PASS`.

### Skutki zapieczętowania

- Plik `run_manifest.json` jest oznaczany `seal_status: sealed`.
- Suma kontrolna manifestu jest zapisywana w `reports/<run_id>/manifest_checksum.txt`.
- Artefakty w `reports/<run_id>/` nie mogą być nadpisywane po zapieczętowaniu.
- SQLite rejestruje `sealed_at` w rekordzie przebiegu.

Zapieczętowanie jest prererekwizytem dla `promotion-eligible run`.

---

## Przebieg autorytatywny a nieautorytatywny

### Przebieg nieautorytatywny (`untrusted_local` lub `trusted_ci`)

- Może być używany do exploracji, debugowania, weryfikacji CI.
- Może być cytowany jako dowód dla bramek G1–G3, jeśli `trust_level = trusted_ci` i `evidence.complete = true`.
- **Nie może** być podstawą decyzji promocji do stabilnego repozytorium.
- **Nie może** być zapieczętowany.

### Przebieg autorytatywny (`authoritative_verify`)

- Wykonany celowo, z pełną izolacją środowiska, rzeczywistymi danymi wejściowymi.
- Może być cytowany jako dowód dla bramek G1–G5.
- Może być zapieczętowany.
- Jest prererekwizytem `promotion-eligible run`.
- Wymaga explicitnej deklaracji `context_type = authoritative_verify` przy wywołaniu CLI.

### Konsekwencje błędnej klasyfikacji

Deklaracja `authoritative_verify` dla przebiegu, który nie spełnia wymagań tego kontekstu (np. aktywny mock-forbidden z mockiem), jest naruszeniem kontraktu i skutkuje rejection kontekstu + exit code 1. Nie jest silent downgrade do `untrusted_local`.

---

## Warunki wstępne promocji

Przed wykonaniem `itdlab export repo1` (promocja do stabilnego repozytorium) muszą być spełnione wszystkie poniższe warunki:

| # | Warunek | Weryfikacja |
|---|---------|-------------|
| PR1 | Co najmniej jeden `promotion-eligible run` istnieje dla każdego wycinka funkcjonalnego | `itdlab audit runs --promotion-eligible` |
| PR2 | Wszystkie bramki G1–G5 mają status `PASS` | Raport bramek |
| PR3 | Brak aktywnych pomijań Kategorii 3 lub 4 na ścieżkach krytycznych | `docs/SKIP_REGISTER.md` |
| PR4 | Każdy `promotion-eligible run` ma `seal_status: sealed` | Manifest przebiegu |
| PR5 | Żaden z cytowanych przebiegów nie ma `context_validation_result: rejected` | Manifest przebiegu |
| PR6 | `itdlab export repo1 --dry-run` kończy się z exit code 0 | Weryfikacja przed produkcją |

Niespełnienie któregokolwiek z PR1–PR6 musi blokować wykonanie `itdlab export repo1` z exit code 2 (gate failure).

---

## Zakres

Ten kontrakt ma zastosowanie do:
- Wszystkich poleceń `itdlab` modyfikujących stan SQLite
- Wszystkich poleceń `itdlab` dołączających do dziennika zdarzeń
- Wszystkich poleceń `itdlab` produkujących artefakty dowodowe

Nie ma zastosowania do:
- `itdlab --help` i `itdlab [command] --help`
- `itdlab ingest inspect` (inspekcja tylko do odczytu, bez zmiany stanu)

Polecenia tylko do odczytu nadal muszą zapisywać na stdout, lecz nie są zobowiązane do tworzenia pakietów dowodowych.

---

## Odwołania wewnętrzne
- `docs/EXECUTION_ASSURANCE_PROGRAM.md`
- `docs/QUALITY_GATES_POLICY.md`
- `docs/TESTING_STANDARD.md`
- `docs/EVIDENCE_MODEL.md`
- `docs/CONTEXT_VOCABULARY.md`
- `docs/RUN_MANIFEST_SCHEMA.md`

## Punkty odniesienia autorytetu
- `docs/REFERENCES.md` — RFC 2119 (semantyka MUST / SHOULD)
- `docs/REFERENCES.md` — ISO/IEC/IEEE 29148 (kontrakt wymagań i warunki postcondition)

## Review metadata
- Owner: project team
- Status: draft
- Last reviewed: 2026-03-30
