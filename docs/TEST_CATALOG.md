# Test Catalog

30 warstw testowych zorganizowanych w 6 poziomów.

Każdy wpis warstwy zawiera:
- **Goal** — jakie ryzyko awarii adresuje
- **Input** — co jest podawane na wejście testu
- **Output** — co test produkuje
- **Artifact** — plik zapisywany do pakietu dowodów (jeśli dotyczy)
- **PASS criterion** — warunek wiarygodnego zaliczenia
- **Blocking gate** — która brama jakości jest zasilana przez tę warstwę
- **Mock policy** — `allowed` / `restricted` / `forbidden` (zob. `docs/POLICY_MOCKS_AND_REAL_PATHS.md`)
- **Evidence strength** — `low` / `medium` / `high` / `promotion-critical`

---

## Level A — Contract & Input (Layers 1–5)

### Layer 1: File Presence Tests
- **Goal:** Weryfikacja, czy wszystkie wymagane pliki źródłowe istnieją przed jakimkolwiek przetwarzaniem.
- **Input:** Ścieżka do katalogu źródłowego
- **Output:** PASS jeśli wszystkie wymagane pliki są obecne, FAIL z listą brakujących plików
- **Artifact:** `presence_check.json`
- **PASS criterion:** Zero brakujących wymaganych plików
- **Blocking gate:** G1
- **Mock policy:** forbidden
- **Evidence strength:** high

### Layer 2: File Readability Tests
- **Goal:** Weryfikacja, czy wszystkie pliki źródłowe można otworzyć i odczytać.
- **Input:** Ścieżki do plików źródłowych
- **Output:** PASS jeśli wszystkie pliki czytelne, FAIL z błędem dla każdego pliku
- **Artifact:** `readability_report.json`
- **PASS criterion:** Zero błędów odczytu
- **Blocking gate:** G1
- **Mock policy:** forbidden
- **Evidence strength:** high

### Layer 3: Encoding Tests
- **Goal:** Wykrywanie naruszeń UTF-8, znaczników BOM, bajtów zerowych.
- **Input:** Surowe bajty pliku
- **Output:** Status kodowania dla każdego pliku
- **Artifact:** `encoding_report.json`
- **PASS criterion:** Zero naruszeń kodowania w wymaganych plikach
- **Blocking gate:** G1
- **Mock policy:** forbidden
- **Evidence strength:** medium

### Layer 4: Markdown Structure Tests
- **Goal:** Weryfikacja, czy parser może wykryć co najmniej jeden nagłówek, listę lub granicę sekcji.
- **Input:** Plik źródłowy Markdown
- **Output:** Podsumowanie struktury (liczba nagłówków, liczba list)
- **Artifact:** osadzony w `parse_report.json`
- **PASS criterion:** Wykryto co najmniej jeden nagłówek
- **Blocking gate:** G1
- **Mock policy:** forbidden
- **Evidence strength:** medium

### Layer 5: Source Schema Tests
- **Goal:** Weryfikacja, czy pliki źródłowe spełniają minimalny kontrakt wejściowy (pole klasy, niepuste ciało).
- **Input:** Przetworzone źródło
- **Output:** Lista naruszeń kontraktu
- **Artifact:** `source_contract_report.json`
- **PASS criterion:** Zero naruszeń kontraktu
- **Blocking gate:** G1
- **Mock policy:** forbidden
- **Evidence strength:** high

---

## Level B — Parser & Extraction (Layers 6–10)

### Layer 6: Parser Unit Tests
- **Goal:** Każda funkcja parsera poprawnie obsługuje minimalny przypadek.
- **Input:** Literały łańcuchowe inline
- **Output:** Przetworzona struktura
- **Artifact:** brak (wyjście testu jednostkowego)
- **PASS criterion:** Wyjście zgodne z oczekiwaną strukturą
- **Blocking gate:** G2
- **Mock policy:** allowed
- **Evidence strength:** low

### Layer 7: Fixture Parsing Tests
- **Goal:** Znane pliki wzorcowe produkują znane wyjście.
- **Input:** `internal/testkit/fixtures/`
- **Output:** Przetworzone dokumenty
- **Artifact:** brak (wyjście testu)
- **PASS criterion:** Wyjście zgodne z oczekiwanym JSON
- **Blocking gate:** G2
- **Mock policy:** allowed
- **Evidence strength:** medium

### Layer 8: Golden Extraction Tests
- **Goal:** Pełne wyjście parsera zgodne z złotym plikiem bajt po bajcie.
- **Input:** Plik wzorcowy
- **Output:** Dokument JSON
- **Artifact:** raport różnic przy niezgodności
- **PASS criterion:** Wyjście == golden
- **Blocking gate:** G2
- **Mock policy:** restricted
- **Evidence strength:** promotion-critical

### Layer 9: Partial Corruption Tests
- **Goal:** Zniekształcone wejście nie powoduje paniki; błąd zwracany jest w kontrolowany sposób.
- **Input:** Skrócony lub nieprawidłowy Markdown
- **Output:** Błąd (nie panika)
- **Artifact:** brak
- **PASS criterion:** Brak paniki, błąd jest niepusty i opisowy
- **Blocking gate:** G2
- **Mock policy:** allowed
- **Evidence strength:** medium

### Layer 10: Determinism Tests
- **Goal:** To samo wejście zawsze produkuje to samo wyjście.
- **Input:** Plik wzorcowy, uruchomiony N razy
- **Output:** N identycznych wyników
- **Artifact:** brak
- **PASS criterion:** Wszystkie N wyników identyczne
- **Blocking gate:** G2
- **Mock policy:** restricted
- **Evidence strength:** promotion-critical

---

## Level C — Normalization & Canonical Model (Layers 11–15)

### Layer 11: Canonical ID Tests
- **Goal:** Semantycznie równoważne nazwy produkują ten sam kanoniczny identyfikator.
- **Input:** Wariantowe ciągi nazw (np. "Risk Register", "risk_register", "Risk-Register")
- **Output:** Kanoniczny identyfikator
- **Artifact:** brak
- **PASS criterion:** Wszystkie warianty → ten sam kanoniczny identyfikator
- **Blocking gate:** G3
- **Mock policy:** allowed
- **Evidence strength:** medium

### Layer 12: Collision Detection Tests
- **Goal:** System wykrywa, gdy dwa różne dokumenty produkują ten sam kanoniczny identyfikator.
- **Input:** Dwa dokumenty o różnej treści, ale podobnych nazwach
- **Output:** Wpis w raporcie kolizji
- **Artifact:** `collision_report.json`
- **PASS criterion:** Kolizja wykryta i odnotowana
- **Blocking gate:** G3
- **Mock policy:** allowed
- **Evidence strength:** high

### Layer 13: Alias Resolution Tests
- **Goal:** Aliasy poprawnie wskazują na swój kanoniczny dokument.
- **Input:** Ciągi aliasów zarejestrowane dla dokumentu
- **Output:** Rozwiązany identyfikator dokumentu
- **Artifact:** brak
- **PASS criterion:** Wszystkie aliasy wskazują na poprawny dokument kanoniczny
- **Blocking gate:** G3
- **Mock policy:** allowed
- **Evidence strength:** medium

### Layer 14: Typing Tests
- **Goal:** Wszystkie pola domenowe mają poprawne typy i prawidłowe wartości.
- **Input:** Utrwalony dokument z SQLite
- **Output:** Raport walidacji typów
- **Artifact:** `type_validation_report.json`
- **PASS criterion:** Zero naruszeń typów
- **Blocking gate:** G3
- **Mock policy:** restricted
- **Evidence strength:** medium

### Layer 15: Migration Tests
- **Goal:** Migracje schematu nie niszczą ani nie uszkadzają istniejących danych.
- **Input:** Snapshot SQLite sprzed migracji + skrypt migracji
- **Output:** Stan po migracji
- **Artifact:** raport porównania liczby wierszy
- **PASS criterion:** Wszystkie wiersze nienaruszone, nowe kolumny mają poprawne wartości domyślne
- **Blocking gate:** G3
- **Mock policy:** restricted
- **Evidence strength:** high

---

## Level D — Relations & Semantics (Layers 16–20)

### Layer 16: Relation Rule Unit Tests
- **Goal:** Każda reguła relacji poprawnie wyzwala się na minimalnej parze wejściowej.
- **Input:** Dwa uproszczone dokumenty + definicja reguły
- **Output:** Relacja (lub brak)
- **Artifact:** brak
- **PASS criterion:** Wyjście zgodne z oczekiwaną relacją
- **Blocking gate:** G4
- **Mock policy:** allowed
- **Evidence strength:** low

### Layer 17: Relation Consistency Tests
- **Goal:** Nie istnieją sprzeczne relacje.
- **Input:** Pełny zestaw relacji z SQLite
- **Output:** Lista sprzeczności
- **Artifact:** `relation_consistency_report.json`
- **PASS criterion:** Zero sprzeczności
- **Blocking gate:** G4
- **Mock policy:** restricted
- **Evidence strength:** high

### Layer 18: Explainability Tests
- **Goal:** Każda wywnioskowana relacja ma niepuste pole wyjaśnienia.
- **Input:** Wszystkie relacje w SQLite gdzie source = 'rule_engine'
- **Output:** Wiersze z pustym wyjaśnieniem
- **Artifact:** `explainability_report.json`
- **PASS criterion:** Zero wierszy z pustym wyjaśnieniem
- **Blocking gate:** G4
- **Mock policy:** restricted
- **Evidence strength:** promotion-critical

### Layer 19: Graph Acyclicity Tests
- **Goal:** Wykrywanie cykli w skierowanych grafach relacji, gdzie cykle są nieprawidłowe.
- **Input:** Krawędzie relacji depends_on
- **Output:** Lista cykli
- **Artifact:** `acyclicity_report.json`
- **PASS criterion:** Zero cykli w grafie depends_on
- **Blocking gate:** G4
- **Mock policy:** restricted
- **Evidence strength:** high

### Layer 20: Section Influence Tests
- **Goal:** Reguły wnioskowania sekcja-do-sekcji wyzwalają się dla znanych przypadków kanonicznych.
- **Input:** Dokument o znanych strukturach sekcji
- **Output:** Wywnioskowane relacje sekcji
- **Artifact:** brak
- **PASS criterion:** Wyjście zgodne z oczekiwaną mapą wpływu sekcji
- **Blocking gate:** G4
- **Mock policy:** restricted
- **Evidence strength:** medium

---

## Level E — Interface & Run (Layers 21–25)

### Layer 21: CLI Contract Tests
- **Goal:** Wszystkie polecenia i flagi istnieją; kody wyjścia są poprawne.
- **Input:** Wywołanie pliku binarnego CLI
- **Output:** Kod wyjścia, struktura stdout
- **Artifact:** brak
- **PASS criterion:** Kod wyjścia 0 dla prawidłowego wejścia, niezerowy dla nieprawidłowego
- **Blocking gate:** G4, G5
- **Mock policy:** restricted
- **Evidence strength:** medium

### Layer 22: End-to-End Slice Tests
- **Goal:** Pełny wycinek funkcjonalności działa bez błędów i produkuje oczekiwany stan.
- **Input:** Wzorcowy katalog źródłowy
- **Output:** Stan SQLite + zdarzenia JSONL
- **Artifact:** liczba wierszy SQLite, liczba linii w dzienniku zdarzeń
- **PASS criterion:** SQLite zgodny z oczekiwanymi wierszami; zdarzenia są niepuste i prawidłowym JSONL
- **Blocking gate:** G5
- **Mock policy:** forbidden
- **Evidence strength:** promotion-critical

### Layer 23: Resume / Restart Tests
- **Goal:** Przebieg przerwany w połowie może zostać wznowiony bez utraty ani duplikacji danych.
- **Input:** Częściowo ukończony stan przebiegu
- **Output:** Ukończony stan przebiegu
- **Artifact:** raport porównania liczby wierszy
- **PASS criterion:** Stan końcowy identyczny z nieprzerwianym przebiegiem
- **Blocking gate:** G5
- **Mock policy:** forbidden
- **Evidence strength:** high

### Layer 24: Event Log Integrity Tests
- **Goal:** Dziennik zdarzeń JSONL jest prawidłowy, wyłącznie dołączany i możliwy do przetworzenia linia po linii.
- **Input:** `runs/events.jsonl`
- **Output:** Wynik przetwarzania dla każdej linii
- **Artifact:** `event_log_integrity_report.json`
- **PASS criterion:** Zero błędów przetwarzania; żadna linia nie jest modyfikowana po zapisie
- **Blocking gate:** G5
- **Mock policy:** forbidden
- **Evidence strength:** promotion-critical

### Layer 25: SQLite Materialization Tests
- **Goal:** Stan SQLite jest spójny z dziennikiem zdarzeń.
- **Input:** SQLite + events.jsonl dla tego samego przebiegu
- **Output:** Różnica między odtworzonym stanem zdarzeń a stanem SQLite
- **Artifact:** `materialization_diff_report.json`
- **PASS criterion:** Zero różnic
- **Blocking gate:** G5
- **Mock policy:** forbidden
- **Evidence strength:** promotion-critical

---

## Level F — Operational & Audit (Layers 26–30)

### Layer 26: Reproducibility Tests
- **Goal:** Uruchomienie tego samego polecenia na tych samych danych dwa razy produkuje ten sam stan SQLite.
- **Input:** Dane źródłowe + dwa identyczne wywołania przebiegu
- **Output:** Para sum kontrolnych SQLite
- **Artifact:** `reproducibility_report.json`
- **PASS criterion:** Obie sumy kontrolne identyczne
- **Blocking gate:** G5
- **Mock policy:** forbidden
- **Evidence strength:** promotion-critical

### Layer 27: Evidence Pack Tests
- **Goal:** Każdy ukończony przebieg produkuje kompletny pakiet dowodów.
- **Input:** Ukończony przebieg
- **Output:** Manifest pakietu dowodów
- **Artifact:** `evidence_manifest.json`
- **PASS criterion:** Wszystkie wymagane artefakty obecne (zob. `docs/EVIDENCE_MODEL.md`)
- **Blocking gate:** G5
- **Mock policy:** forbidden
- **Evidence strength:** promotion-critical

### Layer 28: Performance Budget Tests
- **Goal:** Kluczowe operacje kończą się w określonych limitach czasu i pamięci.
- **Input:** Standardowy korpus wzorcowy
- **Output:** Pomiary czasu i pamięci
- **Artifact:** `performance_report.json`
- **PASS criterion:** Poniżej progów budżetu (zdefiniowanych per polecenie w Makefile)
- **Blocking gate:** brak (doradczy)
- **Mock policy:** forbidden
- **Evidence strength:** medium

### Layer 29: Failure-Mode Tests
- **Goal:** Błędy są przechwytywane, poprawnie raportowane i nie pozostawiają częściowego stanu.
- **Input:** Wstrzyknięte scenariusze awarii (brakujący plik, uszkodzona baza danych, błąd zapisu)
- **Output:** Komunikat błędu, kod wyjścia, stan bazy danych i dziennika zdarzeń
- **Artifact:** `failure_mode_report.json`
- **PASS criterion:** Niezerowy kod wyjścia, opisowy błąd, brak częściowych zapisów
- **Blocking gate:** G5
- **Mock policy:** forbidden
- **Evidence strength:** high

### Layer 30: Release Gate Tests
- **Goal:** Promocja do repozytorium referencyjnego jest zablokowana, jeśli nie wszystkie bramy jakości zostały zaliczone.
- **Input:** Pełny wynik przebiegu
- **Output:** Raport zaliczenia/niezaliczenia bram
- **Artifact:** `gate_pass_report.json`
- **PASS criterion:** Wszystkie bramy G0–G4 zaliczone; pakiet dowodów kompletny; kod wyjścia 0
- **Blocking gate:** G5
- **Mock policy:** forbidden
- **Evidence strength:** promotion-critical
