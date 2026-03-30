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

## Pojęcia kontekstu wykonania

### Context / Execution context
Zbiór warunków środowiskowych i deklaracji operatora, które określają, w jakim trybie wykonywany jest przebieg. Kontekst jest walidowany przed pierwszym zapisem stanu. Przebieg, którego kontekst nie przeszedł walidacji, jest przebiegiem **odrzuconym** (`REJECTED`).

### context_id
Unikalny identyfikator kontekstu wykonania przebiegu. Zakodowany w manifeście uruchomienia i w rekordzie przebiegu w SQLite. Umożliwia późniejszą ocenę: „czy ten przebieg był wykonany w dozwolonym kontekście?". Format: `<typ>-<środowisko>-<timestamp>`, np. `authoritative-local-20260330`.

### run_context
Konkretna instancja kontekstu wykonania dla przebiegu. Obejmuje: `context_id`, typ kontekstu (`untrusted_local` | `trusted_ci` | `authoritative_verify`), wynik walidacji oraz ewentualne naruszenia.

### Allowed context
Kontekst, który przeszedł walidację i spełnia minimalne wymagania dla danego polecenia. Tylko przebieg w `allowed context` może produkować artefakty kwalifikujące się jako dowód bramki.

### Rejected context
Kontekst, który nie przeszedł walidacji — brakuje wymaganych warunków (np. baza danych w trybie in-memory, brak wymaganej wersji schematu, aktywny `mock-forbidden` naruszony). Przebieg w `rejected context` musi zakończyć się z kodem 1. Pole `context_validation_result` w manifeście zawiera przyczyny odrzucenia.

### Authoritative run
Przebieg wykonany w kontekście `authoritative_verify` — celowo, przez operatora, z pełną izolacją, rzeczywistymi danymi wejściowymi, bez aktywnych pomijań Kategorii 3/4 oraz z kompletnym pakietem dowodowym. Jest to najwyższy poziom zaufania. Może być cytowany w raporcie bramki G5 i decyzji o promocji.

### Sealed run
Przebieg, który zakończył się pomyślnie, a jego manifest i artefakty zostały zablokowane (oznaczone jako `seal_status: sealed`). Zapieczętowanego przebiegu nie można zmodyfikować ani zaktualizować. Jest to prererekwizyt dla `promotion-eligible run`.

### Promotion-eligible run
Przebieg spełniający wszystkie poniższe warunki:
- `trust_level: authoritative_verify`,
- `evidence.complete = true`,
- `seal_status: sealed`,
- wszystkie wymagane bramki w statusie `PASS`,
- brak aktywnych pomijań Kategorii 3 lub 4 na warstwach ścieżki krytycznej.

Tylko `promotion-eligible run` może być cytowany jako podstawa decyzji promocji do stabilnego repozytorium.

### REJECTED
Status przebiegu lub kontekstu wskazujący, że walidacja wstępna nie powiodła się. Przebieg `REJECTED` nie zapisuje stanu (poza zapisem odmowy na stderr), nie tworzy pakietu dowodowego i nie może być cytowany.

### degraded
Status bramki lub przebiegu wskazujący, że warunki są częściowo spełnione: co najmniej jedno aktywne pomijanie Kategorii 3 lub warstwa `mock-restricted` bez weryfikacji ścieżką rzeczywistą. Bramka `degraded` blokuje promocję, lecz nie blokuje dalszego wytwarzania.

### Mandatory artifact
Artefakt wymagany dla każdego przebiegu zmieniającego stan, niezależnie od polecenia: `stdout.txt`, `db_checksum.txt`, `summary.md`, `run_manifest.json`. Zdefiniowane w `docs/EVIDENCE_MODEL.md`. Brak któregokolwiek z mandatory artifacts czyni przebieg INCOMPLETE.

### Command-specific artifact
Artefakt wymagany tylko dla określonego polecenia CLI (np. `parse_report.json` dla `itdlab ingest run`, `collision_report.md` dla `itdlab normalize apply`). Lista per polecenie w `docs/EVIDENCE_MODEL.md`. Brak command-specific artifact czyni przebieg INCOMPLETE dla danego polecenia.

### Blocking gate
Bramka jakości, której wynik `FAIL` bezpośrednio blokuje dalsze operacje lub promocję. Każda warstwa testowa ma przypisaną blokującą bramkę w `docs/TEST_CATALOG.md` (pole Blocking gate). Przebieg nie może przejść do kolejnego etapu, jeśli jego blokująca bramka jest w statusie `FAIL`.

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
- `docs/EXECUTION_CONTRACT.md`

## Punkty odniesienia autorytetu
- `docs/REFERENCES.md` — RFC 2119 (semantyka MUST / SHOULD)
- `docs/REFERENCES.md` — ISO/IEC/IEEE 29148 (kontrakt wymagań)

## Review metadata
- Owner: project team
- Status: draft
- Last reviewed: 2026-03-30
