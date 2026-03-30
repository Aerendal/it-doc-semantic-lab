# Policy: Mocks and Real Paths

Dokument definiuje, które warstwy testowe mogą używać zaślepek (mocków), które wymagają weryfikacji na rzeczywistej ścieżce, oraz co oznacza „rzeczywista ścieżka" w tym repozytorium.

---

## Zasada

**Mock ukrywający realne błędy jest gorszy niż brak testu.**

Mocki są dopuszczalne tam, gdzie izolują testowaną jednostkę od nieistotnej infrastruktury. Są ograniczone lub zakazane tam, gdzie maskowałyby błędy poprawności, ukrywały problemy integracyjne lub uniemożliwiały produkcję dowodów na podstawie rzeczywistych danych wejściowych.

---

## Definicje

### Mock-allowed
Warstwa może używać zastępników testowych (stubów, fejków, implementacji in-memory) dla swoich zależności. Wkład warstwy do bramki jest w pełni wiarygodny.

### Mock-restricted
Warstwa może używać mocków jedynie dla infrastruktury wyraźnie nieistotnej (np. zegar, ziarno losowości). Główne ścieżki I/O (system plików, SQLite, dziennik zdarzeń) muszą być realizowane z użyciem rzeczywistych implementacji lub wiernych fejków opartych na rzeczywistym I/O plikowym. Wkład do bramki jest wiarygodny tylko wtedy, gdy przeprowadzono również weryfikację na rzeczywistej ścieżce.

### Mock-forbidden
Warstwa musi korzystać z rzeczywistych implementacji — rzeczywistego systemu plików, rzeczywistego SQLite na dysku, rzeczywistego pliku dziennika zdarzeń, rzeczywistych dokumentów źródłowych. Użycie mocków dla głównych zależności unieważnia wkład warstwy do bramki.

### Weryfikacja na rzeczywistej ścieżce
Przebieg testowy, w którym:
1. Dane wejściowe to rzeczywisty plik z `internal/testkit/fixtures/` lub `sources/`.
2. Zapisy SQLite trafiają do rzeczywistego (tymczasowego) pliku bazy danych na dysku.
3. Dopisywanie do dziennika zdarzeń trafia do rzeczywistego (tymczasowego) pliku JSONL na dysku.
4. Dane wyjściowe są porównywane ze znanym poprawnym plikiem golden lub oczekiwaną strukturą.

---

## Polityka mocków dla warstw

| Poziom | Warstwy | Polityka | Uwagi |
|--------|---------|----------|-------|
| A — Contract & Input | 1–5 | **Mock-forbidden** | Muszą czytać rzeczywiste pliki z `sources/` lub fixtures |
| B — Parser & Extraction | 6–7 | Mock-allowed | Testy jednostkowe i fixture mogą używać ciągów in-memory |
| B — Parser & Extraction | 8–10 | **Mock-restricted** | Testy golden i determinizmu muszą używać rzeczywistych plików fixture |
| C — Normalization | 11–13 | Mock-allowed | Logika kanonicznego ID i aliasów może używać danych wejściowych in-memory |
| C — Normalization | 14–15 | **Mock-restricted** | Testy typów i migracji wymagają rzeczywistego SQLite na dysku |
| D — Relations & Semantics | 16 | Mock-allowed | Testy jednostkowe reguł mogą używać stubów dokumentów in-memory |
| D — Relations & Semantics | 17–20 | **Mock-restricted** | Testy spójności, wyjaśnialności, cykli i wpływu wymagają rzeczywistego SQLite |
| E — Interface & Run | 21 | **Mock-restricted** | Testy kontraktu CLI wymagają wywołania rzeczywistego binarnego; mogą używać tymczasowej bazy |
| E — Interface & Run | 22–25 | **Mock-forbidden** | Testy end-to-end, wznawiania, integralności dziennika zdarzeń i materializacji wymagają pełnego stosu |
| F — Operational & Audit | 26–30 | **Mock-forbidden** | Wszystkie warstwy operacyjne i audytowe wymagają rzeczywistego systemu plików, bazy danych i dziennika zdarzeń |

---

## Wymóg weryfikacji na rzeczywistej ścieżce

Następujące warstwy mają **obowiązkowy przebieg weryfikacji na rzeczywistej ścieżce** poza wszelkimi testami jednostkowymi opartymi na mockach:

- Layer 8 (Golden Extraction Tests)
- Layer 10 (Determinism Tests)
- Layer 22 (End-to-End Slice Tests)
- Layer 24 (Event Log Integrity Tests)
- Layer 25 (SQLite Materialization Tests)
- Layer 26 (Reproducibility Tests)
- Layer 27 (Evidence Pack Tests)

Dla tych warstw zestaw testów uruchamiany wyłącznie na mockach lub stanie in-memory nie spełnia kryterium PASS warstwy.

---

## Polityka plików fixture i golden

- Pliki fixture znajdują się w `internal/testkit/fixtures/`. Są to rzeczywiste pliki Markdown reprezentujące znane dane wejściowe dokumentów.
- Pliki golden znajdują się w `internal/testkit/golden/`. Są to znane poprawne dane wyjściowe dla konkretnych danych wejściowych.
- Ani pliki fixture, ani golden nie mogą być automatycznie generowane przez ten sam testowany kod bez niezależnego kroku przeglądu.
- Pliki golden muszą być aktualizowane celowo poleceniem `go test ./... -update`, przeglądane w diff PR i zatwierdzone przed scaleniem.

---

## Kiedy mock-restricted oznacza „rzeczywiste dla rdzenia, mock dla krawędzi"

Dla warstw mock-restricted dopuszczalny wzorzec to:

```
core path: real filesystem / real SQLite / real event log
edges:     mock clock (use fixed time in tests)
edges:     mock random (use fixed seed)
edges:     mock external HTTP (not applicable in this repo)
```

Ścieżka rdzenia nigdy nie może być zastąpiona mockiem w warstwie mock-restricted.

---

## Wiarygodność bramki a polityka mocków

| Polityka mocków | Wkład do bramki |
|----------------|----------------|
| Mock-allowed | W pełni wiarygodny |
| Mock-restricted — weryfikacja na rzeczywistej ścieżce obecna | W pełni wiarygodny |
| Mock-restricted — weryfikacja na rzeczywistej ścieżce nieobecna | Wkład do bramki: `degraded` |
| Mock-forbidden — użyto mocka | Wkład do bramki: **nieważny** (warstwa nie jest liczona) |

Warstwa oznaczona `mock-forbidden`, wykonana z użyciem mocków, musi być zarejestrowana jako pominięcie (Kategoria 3) zgodnie z `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`.
