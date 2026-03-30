# Development Plan — it-doc-semantic-lab

## Cel
To repozytorium jest laboratorium badań semantycznych i eksperymentów dla projektu IT-Dokumentacja.

Stabilne repozytorium produkcyjne (IT-Dokumentacja) zawiera możliwości P1.5.
Niniejsze laboratorium buduje warstwę semantyczną P2→P3.

## Cele warstw

### P2 — structural-rule
- [ ] Detektor klasy dokumentu (exp_001)
- [ ] Maper ról sekcji (exp_002)
- [ ] Rejestr klas i ról (schemas/document_classes.yaml)
- [ ] Ewaluator kompletności ról

### P3 — semantic-role
- [ ] Generator raportów luk (exp_003)
- [ ] Planer podziału (exp_004)
- [ ] Skorer kompletności treści

## Pierwsze przypadki korpusu
1. nlp_algorithm_architecture_spec → case_001 (morfologia_polska_algorytmy) ✓
2. project_architecture_concept → case_001 (do dodania)
3. compliance_procedure → case_001 (do dodania)

## Zależność od zapewnienia wykonania

Praca implementacyjna P2→P3 nie jest traktowana jako zakończona jedynie z powodu istnienia kodu.

Każdy wycinek możliwości powinien być oceniany względem:
- bramek jakości,
- odpowiednich warstw testowych,
- wymagań dowodowych,
- zasad rzeczywistej vs zamockowanej ścieżki,
- kryteriów gotowości do promocji.

Kolejność implementacji powinna zatem śledzić nie tylko ukończenie funkcji, ale również:
1. pokrycie warstwami testowymi,
2. kompletność paczki dowodowej,
3. status bramki,
4. wykonalność promocji do stabilnego repozytorium.

Patrz: `docs/EXECUTION_ASSURANCE_PROGRAM.md`, `docs/QUALITY_GATES_POLICY.md`, `docs/TESTING_STANDARD.md`.

## Polityka promocji
Stabilne komponenty są promowane do IT-Dokumentacja/itdoc/ tylko wtedy, gdy:
1. Eksperyment ma wzorzec złoty
2. Eksperyment przechodzi na wszystkich przypadkach korpusu
3. Dry-run/apply/idempotency potwierdzone
4. Walidacja przechodzi w IT-Dokumentacja po integracji
