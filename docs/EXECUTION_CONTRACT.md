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

## Punkty odniesienia autorytetu
- `docs/REFERENCES.md` — RFC 2119 (semantyka MUST / SHOULD)

## Review metadata
- Owner: project team
- Status: draft
- Last reviewed: 2026-03-30
