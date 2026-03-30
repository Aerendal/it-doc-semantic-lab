# Risk-to-Test Matrix

Dokument ten mapuje zidentyfikowane kategorie ryzyka na warstwy testów, które je adresują, obowiązującą politykę mocków oraz wymaganą siłę dowodów dla wiarygodnego wyniku.

---

## Cel

Macierz ryzyka do testów jawnie określa powiązanie między tym, co może pójść nie tak, a tym, które warstwy testów stanowią główną obronę przed każdym trybem awarii. Wspiera ona:
- ocenę bramek (które ryzyka są pokryte przed promocją),
- ocenę wpływu pominięć (które ryzyka stają się niepokryte po pominięciu warstwy),
- priorytetyzację testów (które warstwy budować najpierw dla nowego wycinka funkcjonalności).

---

## Kategorie ryzyka

| ID | Kategoria ryzyka | Opis |
|----|-----------------|------|
| R1 | Uszkodzone lub brakujące dane wejściowe | Pliki źródłowe nieobecne, nieczytelne lub zniekształcone |
| R2 | Awaria parsera | Parser produkuje błędne wyjście, panikuje lub jest niedeterministyczny |
| R3 | Niestabilność tożsamości | Ten sam koncept otrzymuje różne kanoniczne identyfikatory między uruchomieniami |
| R4 | Cicha kolizja nazw | Dwa różne dokumenty zmapowane na ten sam kanoniczny identyfikator bez wykrycia |
| R5 | Dryf schematu | Zmiany modelu dziedzinowego psują istniejące utrwalone dane |
| R6 | Błędne wnioskowanie relacji | Relacje wywnioskowane nieprawidłowo lub bez wyjaśnienia |
| R7 | Naruszenie integralności grafu | Nieodkryte cykle zależności lub sprzeczne relacje |
| R8 | Zerwanie kontraktu CLI | Flagi poleceń, kody wyjścia lub format wyjścia zmienione po cichu |
| R9 | Uszkodzenie stanu uruchomienia | Częściowe zapisy pozostawiają SQLite w niespójnym stanie |
| R10 | Utrata dowodów | Uruchomienie kończy się, ale artefakty są brakujące lub puste |
| R11 | Brak odtwarzalności | Te same dane wejściowe produkują różne wyjście przy powtarzanych uruchomieniach |
| R12 | Ciche pominięcie | Warstwa testów nie jest wykonywana bez rejestracji ani powiadomienia |
| R13 | Promocja bez bramki | Funkcja promowana do stabilnego repozytorium bez zaliczenia wszystkich bramek |
| R14 | Maskowanie rzeczywistych awarii przez mocki | Warstwa `mock-forbidden` wykonana z mockami, ukrywając rzeczywiste awarie |

---

## Mapowanie ryzyk na warstwy

| Ryzyko | Warstwy główne | Warstwy wspierające | Blokująca bramka |
|--------|---------------|---------------------|-----------------|
| R1 — Uszkodzone/brakujące dane wejściowe | 1, 2, 3, 4, 5 | — | G1 |
| R2 — Awaria parsera | 6, 7, 8, 9, 10 | — | G2 |
| R3 — Niestabilność tożsamości | 11, 13 | 10 | G3 |
| R4 — Cicha kolizja nazw | 12 | 11 | G3 |
| R5 — Dryf schematu | 14, 15 | — | G3 |
| R6 — Błędne wnioskowanie relacji | 16, 17, 18 | 20 | G4 |
| R7 — Naruszenie integralności grafu | 17, 19 | — | G4 |
| R8 — Zerwanie kontraktu CLI | 21 | 22 | G4, G5 |
| R9 — Uszkodzenie stanu uruchomienia | 23, 25 | 22 | G5 |
| R10 — Utrata dowodów | 27, 24 | 22 | G5 |
| R11 — Brak odtwarzalności | 26, 10 | 8 | G5 |
| R12 — Ciche pominięcie | (egzekwowanie polityki) | 27 | all |
| R13 — Promocja bez bramki | 30 | 27 | G5 |
| R14 — Maskowanie przez mocki | (egzekwowanie polityki) | 22, 25, 26 | G4, G5 |

---

## Ryzyko a polityka mocków

Ryzyka najbardziej wrażliwe na naruszenia polityki mocków:

| Ryzyko | Dlaczego mocki są tu niebezpieczne |
|--------|-----------------------------------|
| R2 — Awaria parsera | Mockowany parser zawsze zwraca oczekiwane wyjście; rzeczywiste awarie niewidoczne |
| R9 — Uszkodzenie stanu uruchomienia | SQLite w pamięci nie ćwiczy WAL, odzyskiwania po awarii ani błędów zapisu na dysk |
| R10 — Utrata dowodów | Mockowany system plików po cichu pochłania artefakty |
| R11 — Brak odtwarzalności | Testy determinizmu z mockowanym I/O nie wykrywają zmienności na poziomie systemu plików |
| R14 — Maskowanie przez mocki | Z definicji |

Warstwy adresujące R2, R9, R10, R11 są `mock-restricted` lub `mock-forbidden` zgodnie z `docs/POLICY_MOCKS_AND_REAL_PATHS.md`.

---

## Ryzyko a siła dowodów

| Siła dowodów | Ryzyka, które wiarygodnie pokrywa |
|-------------|----------------------------------|
| `low` | R2 (częściowo), R3 (częściowo) — wyłącznie na poziomie jednostkowym |
| `medium` | R1, R2, R3, R4, R5, R6, R7, R8 |
| `high` | R1, R4, R5, R7, R9, R13 |
| `promotion-critical` | R2, R3, R10, R11, R13, R14 |

Do promocji wszystkie warstwy `promotion-critical` muszą być wykonane bez mocków i muszą produkować wyniki PASS z artefaktami.

---

## Pokrycie według wycinka funkcjonalności

| Wycinek | Główne pokrywane ryzyka | Minimalne wymagane warstwy |
|---------|------------------------|---------------------------|
| 1 — Ingest | R1, R2 | Layers 1–10 |
| 2 — Normalize | R3, R4, R5 | Layers 11–15 |
| 3 — Relations | R6, R7 | Layers 16–20 |
| 4 — CLI + Run | R8, R9 | Layers 21–25 |
| 5 — Audit + Export | R10, R11, R12, R13, R14 | Layers 26–30 |

Wycinek funkcjonalności nie jest gotowy do promocji, dopóki wszystkie jego minimalnie wymagane warstwy nie przechodzą z wykonaniem `mock-forbidden` lub `mock-restricted + real-path`.

---

## Wewnętrzne odniesienia
- `docs/TESTING_STANDARD.md`
- `docs/TEST_CATALOG.md`
- `docs/QUALITY_GATES_POLICY.md`
- `docs/POLICY_MOCKS_AND_REAL_PATHS.md`
- `docs/CONTEXT_VOCABULARY.md`

## Metadane przeglądu
- Właściciel: zespół projektowy
- Status: szkic
- Last reviewed: 2026-03-30
