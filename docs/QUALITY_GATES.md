# Quality Gates

Bramki jakości są warunkami blokującymi. Przebieg lub promocja **nie może kontynuować** za bramką, dopóki wszystkie jej warunki nie zostaną spełnione.

---

## Gate 1: Input Contract (before any processing)

| Warunek | Test Layer |
|-----------|-----------|
| Wszystkie wymagane pliki źródłowe są obecne | Layer 1 |
| Wszystkie pliki są czytelne | Layer 2 |
| Zero naruszeń kodowania | Layer 3 |
| Wszystkie pliki mają strukturę możliwą do sparsowania | Layer 4 |
| Wszystkie pliki spełniają kontrakt schematu źródłowego | Layer 5 |

**Zablokowane jeśli:** Dowolny warunek nie zostanie spełniony → przebieg jest przerywany z `exit 1` i raportem błędu kontraktu.

---

## Gate 2: Parse Quality (after ingest)

| Warunek | Test Layer |
|-----------|-----------|
| Wszystkie fixture golden tests przechodzą pomyślnie | Layer 8 |
| Brak panik parsera na jakimkolwiek wejściu | Layer 9 |
| Wynik parsera jest deterministyczny | Layer 10 |

**Zablokowane jeśli:** Dowolny warunek nie zostanie spełniony → krok normalizacji nie uruchamia się.

---

## Gate 3: Normalization Integrity (after normalize)

| Warunek | Test Layer |
|-----------|-----------|
| Wszystkie kanoniczne identyfikatory są unikalne | Layer 11 |
| Zero nierozwiązanych kolizji | Layer 12 |
| Wszystkie aliasy są rozwiązywane poprawnie | Layer 13 |
| Wszystkie pola domeny mają poprawne typy | Layer 14 |

**Zablokowane jeśli:** Dowolny warunek nie zostanie spełniony → krok klasyfikacji nie uruchamia się.

---

## Gate 4: Semantic Consistency (after classify + relations)

| Warunek | Test Layer |
|-----------|-----------|
| Wszystkie wywnioskowane relacje mają niepuste wyjaśnienie | Layer 18 |
| Zero cykli w grafie depends_on | Layer 19 |
| Zero sprzecznych relacji | Layer 17 |

**Zablokowane jeśli:** Dowolny warunek nie zostanie spełniony → krok eksportu nie uruchamia się.

---

## Gate 5: Run Completeness (before export)

| Warunek | Test Layer |
|-----------|-----------|
| Paczka dowodów jest kompletna | Layer 27 |
| Stan SQLite jest spójny z logiem zdarzeń | Layer 25 |
| Rekord przebiegu w SQLite ma `status = 'completed'` | Layer 22 |
| Kod wyjścia CLI wynosił 0 dla wszystkich poprzednich kroków | Layer 21 |

**Zablokowane jeśli:** Dowolny warunek nie zostanie spełniony → eksport do repozytorium referencyjnego jest zablokowany.

---

## Protokół awarii bramki

1. Zaloguj błąd do logu zdarzeń (`action: "gate_failed"`, `entity: "gate"`, `entity_id: "<gate_id>"`)
2. Zapisz szczegóły błędu do `reports/<run_id>/gate_failures.json`
3. Zakończ z kodem `2` (błąd bramki, różny od kodu błędu `1`)
4. Nie modyfikuj żadnego stanu poniżej
