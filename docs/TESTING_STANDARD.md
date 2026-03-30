# Testing Standard

## Filozofia

Testowanie w tym repozytorium jest **oparte na dowodach**: każdy test produkuje lub waliduje artefakt, który może być niezależnie skontrolowany. Testy to nie tylko sprawdzenia regresji — są częścią modelu dowodów przebiegu.

---

## Cele negatywne

Standard ten **nie** zakłada ani nie twierdzi, że:

1. **Wykonanie bez błędów.** Repozytorium zakłada, że błędy będą się zdarzać. Obowiązkiem jest, żeby były widoczne — nie żeby nie występowały.
2. **Zielono = gotowe.** Przechodzący zestaw testów bez kompletnego pakietu dowodów nie jest wiarygodnym wynikiem.
3. **Zielono = możliwe do promocji.** Kod funkcjonalności może być eksperymentalnie zielony bez spełnienia wymagań Gate 4 + dowodów potrzebnych do promocji do stabilnego repozytorium.
4. **Dowolny test = silny dowód.** Test oparty na mockach w warstwie mock-forbidden nie wlicza się do wiarygodności bramy.
5. **Pokrycie = kompletność.** Wysokie pokrycie linii nie zastępuje pokrycia warstw. Warstwa, która nie została wykonana, jest luką — niezależnie od metryk pokrycia linii.

---

## 6 poziomów, 30 warstw

Testy są zorganizowane w 6 poziomów. Każdy poziom adresuje odrębną klasę ryzyka awarii.

| Level | Name | Layers | Primary Risk Addressed |
|-------|------|--------|------------------------|
| A | Contract & Input | 1–5 | Corrupt or missing source data |
| B | Parser & Extraction | 6–10 | Wrong output from parsing |
| C | Normalization & Canonical Model | 11–15 | Identity collisions, type drift |
| D | Relations & Semantics | 16–20 | Incorrect inference, unexplained edges |
| E | Interface & Run | 21–25 | CLI contract breakage, run state corruption |
| F | Operational & Audit | 26–30 | Non-reproducibility, missing evidence, gate failures |

Pełne definicje warstw: zob. [TEST_CATALOG.md](TEST_CATALOG.md).

---

## Podsumowanie polityki mocków

Każda warstwa ma przypisaną politykę mocków. Pełne definicje i tabelę według warstw zob. w `docs/POLICY_MOCKS_AND_REAL_PATHS.md`.

| Policy | Meaning |
|--------|---------|
| **mock-allowed** | Test doubles permitted for all dependencies |
| **mock-restricted** | Core I/O (filesystem, SQLite, event log) must use real implementations; clock/random may be mocked |
| **mock-forbidden** | All primary dependencies must be real; no mocking of filesystem, SQLite, or event log |

Użycie mocków w warstwie `mock-forbidden` unieważnia wkład tej warstwy do bramy. Musi zostać zarejestrowane jako pominięcie (Kategoria 3) zgodnie z `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`.

---

## Podsumowanie polityki pomijania

Każde pominięcie musi być zarejestrowane w `docs/SKIP_REGISTER.md`. Pełną procedurę zob. w `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`.

| Category | Condition | Gate impact |
|----------|-----------|-------------|
| 1 | Infrastructure unavailable, layer not on critical gate path | None |
| 2 | Layer not yet implemented, gate not being evaluated | None |
| 3 | Layer on critical gate path, approved with owner + review date | Gate marked `degraded` — not promotable |
| 4 / unregistered | Forbidden layers, or any skip without registration | Automatic gate failure |

**Zabronione pominięcia (nigdy niedozwolone pod żadnym warunkiem):**
- Evidence pack production (Layer 27)
- Exit code recording
- Event log append
- SQLite run record creation

---

## Bramy jakości

Przebieg **nie może** być promowany dopóki:

- Wszystkie testy Level A przechodzą (kontrakt wejściowy zwalidowany)
- Wszystkie złote testy Level B są zgodne
- Level C: zero nierozwiązanych kolizji
- Level D: wszystkie relacje mają niepuste wyjaśnienie
- Level E: kody wyjścia CLI zgodne z kontraktem
- Level F: pakiet dowodów jest kompletny

Pełne definicje bram: zob. [QUALITY_GATES.md](QUALITY_GATES.md) i [QUALITY_GATES_POLICY.md](QUALITY_GATES_POLICY.md).

---

## Warunki promocji

Funkcjonalność może być **eksperymentalnie zielona** (lokalne przejście testów, działający prototyp) bez bycia **możliwą do promocji**.

Do promocji do stabilnego repozytorium (`IT-Dokumentacja`) funkcjonalność musi spełniać:
1. Gate 4 (spójność semantyczna) w pełni zaliczony — brak bram `degraded`.
2. Pakiet dowodów kompletny dla przebiegu promocji.
3. Brak pomięć Kategorii 3 lub Kategorii 4 w żadnej warstwie zasilającej bramę promocji.
4. Złote pliki są aktualne i zrecenzowane.
5. `itdlab export repo1` kończy się kodem 0.

Funkcjonalność, która jest eksperymentalnie zielona, ale niemożliwa do promocji, powinna być udokumentowana jako taka w `docs/DEVELOPMENT_PLAN.md`.

---

## Standard przeglądu

Recenzent musi być w stanie niezależnie zweryfikować wynik przebiegu **bez rozmowy z autorem**.

Konkretnie, przegląd jest wiarygodny gdy:
- pakiet dowodów dla przebiegu jest kompletny (`itdlab audit evidence <run_id>` kończy się z kodem 0),
- wpisy w dzienniku zdarzeń dla przebiegu są obecne i możliwe do przetworzenia,
- złote pliki są aktualne i zgodne z wynikiem przebiegu,
- wszelkie pominięcia lub tryby zdegradowane są zarejestrowane w `docs/SKIP_REGISTER.md`,
- status bramy jest jawnie odnotowany (PASS / degraded / FAIL), a nie domniemany.

Przegląd, który nie jest niezależnie reprodukowalny na podstawie pakietu dowodów, nie jest wiarygodnym zatwierdzeniem.

---

## Typy testów

### Testy jednostkowe
- Na poziomie pakietu, w plikach `_test.go` obok kodu źródłowego
- Bez I/O, bez systemu plików
- Muszą być deterministyczne

### Testy z plikami wzorcowymi (fixture tests)
- Używają plików z `internal/testkit/fixtures/`
- Wejście to znany plik → weryfikacja dokładnego wyjścia

### Testy złotego pliku (golden tests)
- Porównują wyjście z plikami w `internal/testkit/golden/`
- Aktualizuj poleceniem `go test ./... -update` gdy złoty plik celowo się zmienia
- Zmiany w złotych plikach muszą być recenzowane w diff PR

### Testy integracyjne
- W `internal/testkit/` lub plikach `_integration_test.go`
- Wymagają prawdziwego SQLite i dziennika zdarzeń
- Pomijane z flagą `-short`

### Testy kontraktowe
- Weryfikują kontrakty flag CLI i kody wyjścia
- Uruchamiane jako część standardowego zestawu testów

---

## Konwencje

- Nazwy plików testów: `<subject>_test.go`
- Nazwy plików wzorcowych: `<subject>_<case>.md` / `.json`
- Nazwy złotych plików: `<test_name>.golden`
- Uruchom `make test-short` aby pominąć testy integracyjne
- Uruchom `make test` dla pełnego zestawu

---

## Wymóg determinizmu

Każda funkcja produkująca wyjście używane w testach **musi** zwracać identyczne wyjście dla identycznego wejścia. Niedeterminizm jest błędem testu.

Funkcje niedeterministyczne (np. generatory UUID, producenci znaczników czasu) muszą być wstrzykiwalne i zastępowane deterministycznymi atrapami w testach. Muszą być oznaczone w kodzie adnotacją: `// non-deterministic: <reason>`.
