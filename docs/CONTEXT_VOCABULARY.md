# Context Vocabulary

Ten dokument definiuje kanoniczne znaczenie terminów używanych we wszystkich dokumentach dotyczących polityki, testowania, dowodów i bramek w tym repozytorium. Gdy termin pojawia się w jakimkolwiek dokumencie w `docs/`, jego znaczenie jest takie jak tu zdefiniowane, chyba że dokument używający explicite je nadpisuje i stwierdza to nadpisanie.

---

## Pojęcia podstawowe

### Run
Pojedyncze wywołanie CLI `itdlab`, które produkuje dającą się prześledzić, opartą na dowodach tranzycję stanu. Przebieg:
- ma unikalny `run_id`,
- zaczyna się rekordem startowym w SQLite,
- kończy się rekordem zakończenia w SQLite,
- przez cały czas dołącza zdarzenia do dziennika zdarzeń JSONL,
- przy zakończeniu produkuje pakiet dowodowy.

Wywołanie procesu, które nie spełnia tych kryteriów, nie jest przebiegiem — jest nieśledzonym wykonaniem.

### Trusted run
Przebieg spełniający wszystkie poniższe warunki:
1. Pakiet dowodowy jest kompletny (zob. `docs/EVIDENCE_MODEL.md`).
2. Wszystkie wymagane warstwy testowe dla polecenia przebiegu zostały wykonane (brak niezarejestrowanych pomijań).
3. Warstwy `mock-forbidden` nie zostały wykonane z mockami.
4. Status bramki jest zapisany jako `PASS` lub `degraded` (nie domniemany ani nieobecny).

Trusted run może być przytaczany jako dowód w ocenie bramki lub decyzjach o promocji.

### Untrusted run
Każdy przebieg, który nie spełnia kryteriów trusted run. Untrusted run:
- nie może być przytaczany jako dowód bramki,
- nie może być używany do uzasadnienia promocji,
- nie może być raportowany jako `PASS` bez zastrzeżeń.

Untrusted runs nie są zabronione. Są dozwolone do eksploracji i debugowania. Rozróżnienie ma znaczenie wyłącznie przy podejmowaniu decyzji dotyczących bramek lub promocji.

### Evidence pack
Zestaw artefaktów produkowanych przez przebieg, umożliwiający niezależny audyt. Zdefiniowany w pełni w `docs/EVIDENCE_MODEL.md`. Pakiet dowodowy jest albo **kompletny**, albo **niekompletny** — bez ocen pośrednich.

### INCOMPLETE run
Przebieg, którego pakiet dowodowy jest niekompletny. Zdefiniowany przez sześć warunków w `docs/EVIDENCE_MODEL.md`. INCOMPLETE run nie może być przytaczany jako dowód dla żadnej bramki.

---

## Pojęcia bramek

### Quality gate
Nazwany, explicite blokujący warunek, który musi być spełniony, zanim praca może przejść do następnego etapu. Bramki są zdefiniowane w `docs/QUALITY_GATES.md` i egzekwowane zgodnie z `docs/QUALITY_GATES_POLICY.md`.

### Gate status
Jedna z trzech wartości:
- `PASS` — wszystkie warunki spełnione, dowody kompletne, brak niezarejestrowanych pomijań na warstwach krytycznych.
- `degraded` — warunki częściowo spełnione; co najmniej jedno pomijanie Kategorii 3 jest aktywne na warstwie ścieżki krytycznej lub warstwa `mock-restricted` nie ma weryfikacji ścieżką rzeczywistą.
- `FAIL` — co najmniej jeden warunek blokujący niespełniony lub obecne jest pomijanie Kategorii 4 / niezarejestrowane.

Bramka ze statusem `degraded` blokuje promocję, lecz nie blokuje dalszego wytwarzania.

### Critical path (gate)
Warstwa testowa jest na ścieżce krytycznej bramki, jeśli jej wynik PASS / FAIL bezpośrednio determinuje wynik bramki. Zdefiniowana per bramka w `docs/TEST_CATALOG.md` (pole Blocking gate).

### Credible green
Wynik bramki `PASS`, który jest niezależnie weryfikowalny z pakietu dowodowego bez udziału autora. Pełna definicja w `docs/TESTING_STANDARD.md`.

---

## Pojęcia testowe

### Layer
Nazwany, ograniczony obszar troski testowej. Istnieje 30 warstw zorganizowanych w 6 poziomów (A–F). Każda warstwa ma zdefiniowany cel, dane wejściowe, artefakt, kryterium PASS, blokującą bramkę, politykę mocków i siłę dowodową. Zdefiniowane w `docs/TEST_CATALOG.md`.

### Level
Grupowanie warstw adresujących wspólną klasę ryzyka awarii (np. Level A = kontrakt i wejście, Level F = operacje i audyt).

### Mock policy
Zasada regulująca, jakie test double mogą być używane w danej warstwie. Trzy wartości:
- `mock-allowed` — test double dozwolone dla wszystkich zależności.
- `mock-restricted` — podstawowe ścieżki I/O muszą używać rzeczywistych implementacji; clock/random może być mockowany.
- `mock-forbidden` — wszystkie podstawowe zależności muszą być rzeczywiste.

Zdefiniowana w pełni w `docs/POLICY_MOCKS_AND_REAL_PATHS.md`.

### Real-path verification
Wykonanie testu, w którym dane wejściowe to rzeczywisty plik z `internal/testkit/fixtures/` lub `sources/`, zapisy SQLite trafiają do rzeczywistego pliku na dysku, a dołączenia do dziennika zdarzeń trafiają do rzeczywistego pliku JSONL. Zdefiniowana w `docs/POLICY_MOCKS_AND_REAL_PATHS.md`.

### Fixture
Znany dobry plik wejściowy przechowywany w `internal/testkit/fixtures/`. Fixtures to rzeczywiste pliki, niezależnie przeglądane i wersjonowane w git. Nie są auto-generowane przez testowany kod.

### Golden file
Znane poprawne oczekiwane dane wyjściowe przechowywane w `internal/testkit/golden/`. Aktualizowane celowo za pomocą `go test ./... -update` i przeglądane w diff PR przed scaleniem.

### Evidence strength
Ocena per warstwa tego, jak mocno wynik zaliczający wspiera wiarygodność bramki. Cztery wartości:
- `low` — poziom jednostkowy; wspiera rozumienie, lecz nie decyzje bramek samodzielnie.
- `medium` — poziom fixture lub integracyjny; przyczynia się do wiarygodności bramki.
- `high` — integracja z rzeczywistym I/O; mocno wspiera decyzje bramek.
- `promotion-critical` — wymagane dla wiarygodnego wyniku nadającego się do promocji; brak blokuje promocję.

---

## Pojęcia pomijań i wyjątków

### Skip
Zarejestrowana decyzja o niewykonaniu warstwy testowej dla konkretnego przebiegu lub okresu. Pomijania są klasyfikowane do czterech kategorii (1–4) zdefiniowanych w `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`.

### Unregistered skip
Każda warstwa testowa niewykonana bez odpowiadającego wpisu w `docs/SKIP_REGISTER.md`. Traktowane jako Kategoria 4 (automatyczne niepowodzenie bramki).

### Forbidden skip
Pomijanie elementu Kategorii 4. Nigdy niedozwolone. Skutkuje automatycznym niepowodzeniem bramki dla wszystkich aktywnych bramek w przebiegu.

### Expired skip
Pomijanie, którego `review_date` minęła bez odnowienia. Traktowane jako Kategoria 4 od dnia po wygaśnięciu.

### Exception
Synonim pomijania Kategorii 3 — pomijanie, które jest zatwierdzone, udokumentowane i ograniczone czasowo, lecz które stawia bramkę, której dotyczy, w statusie `degraded`.

---

## Pojęcia promocji

### Promotion
Akt kopiowania stabilnych, zwalidowanych metadanych z tego repozytorium eksperymentalnego do stabilnego repozytorium referencyjnego (`IT-Dokumentacja`). Wymaga Gate 5 `PASS`, kompletnego pakietu dowodowego i pomyślnego kodu wyjścia 0 polecenia `itdlab export repo1`.

### Promotion-ready
Wynik przebiegu spełniający wszystkie warunki promocji: G0–G4 `PASS`, brak bramek `degraded`, kompletny pakiet dowodowy, brak pomijań Kategorii 3 lub 4 na warstwach ścieżki krytycznej.

### Experimentally green
Funkcja, która przechodzi zestaw testów lokalnie lub w CI, lecz jeszcze nie spełnia kryteriów promocji. Jest to prawidłowy stan pośredni. Musi być zapisana w `docs/DEVELOPMENT_PLAN.md` jako odrębna od promotion-ready.

---

## Terminy statusu dokumentu

### draft
Dokument jest pisany; nie był jeszcze przeglądany; nie może być przytaczany jako normatywny.

### in-review
Dokument jest kompletny i w trakcie przeglądu; może być używany jako robocze odwołanie.

### approved
Dokument został przejrzany i zaakceptowany; normatywny dla tego repozytorium.

---

## Odwołania wewnętrzne
- `docs/EXECUTION_ASSURANCE_PROGRAM.md`
- `docs/QUALITY_GATES_POLICY.md`
- `docs/TESTING_STANDARD.md`
- `docs/EVIDENCE_MODEL.md`
- `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`
- `docs/POLICY_MOCKS_AND_REAL_PATHS.md`

## Review metadata
- Owner: project team
- Status: draft
- Last reviewed: 2026-03-30
