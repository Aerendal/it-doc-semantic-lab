# Skip Register

To jest żywy rejestr aktywnych pominięć warstw testowych i zatwierdzonych wyjątków.

Utrzymywany zgodnie z `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`.

Wpis pominięcia tutaj **nie** sprawia, że znika — sprawia, że jest widoczne, ma właściciela i podlega przeglądowi.

---

## Aktywne pominięcia

| Skip ID | Layer | Category | Owner | Registered | Review Date | Gate Impact | Reason |
|---------|-------|----------|-------|------------|-------------|-------------|--------|
| — | — | — | — | — | — | — | *(no active skips)* |

---

## Format Skip ID

`SKIP-<YYYY>-<NNN>` — np. `SKIP-2026-001`

Zwiększaj `NNN` kolejno w ramach roku. Nie używaj ponownie identyfikatorów.

---

## Jak zarejestrować pominięcie

1. Dodaj wiersz do tabeli **Aktywne pominięcia** powyżej.
2. Zapisz Skip ID w odpowiednim pliku testowym jako komentarz:
   ```go
   // SKIP-2026-001: Layer 28 skipped; no performance baseline yet. Review: 2026-09-01.
   t.Skip("SKIP-2026-001")
   ```
3. Zapisz pominięcie w manifeście przebiegu dla przebiegów nim objętych (tablica `skips` w `run_manifest.json`).
4. Przenieś do **Zamknięte pominięcia**, gdy pominięcie zostanie rozwiązane lub wygaśnie.

---

## Wymagane pola

Zgodnie z `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`:

| Pole | Wymagane dla kategorii |
|------|------------------------|
| Skip ID | 1, 2, 3, 4 |
| Layer | 1, 2, 3, 4 |
| Category (1/2/3/4) | 1, 2, 3, 4 |
| Owner | 1, 2, 3, 4 |
| Registered date | 1, 2, 3, 4 |
| Review date | 1, 2 |
| Gate impact | 1, 2, 3, 4 |
| Reason | 1, 2, 3, 4 |

Kategoria 3 (pominięcie na zablokowanej ścieżce krytycznej) i Kategoria 4 (pominięcie stałe) wymagają zatwierdzenia przed rejestracją.

---

## Zamknięte pominięcia

| Skip ID | Layer | Closed Date | Resolution |
|---------|-------|-------------|------------|
| — | — | — | *(no closed skips)* |

---

## Odniesienia wewnętrzne
- `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`
- `docs/TESTING_STANDARD.md`
- `docs/QUALITY_GATES_POLICY.md`

## Metadane przeglądu
- Właściciel: zespół projektowy
- Status: aktywny (dokument żywy)
- Ostatni przegląd: 2026-03-30
