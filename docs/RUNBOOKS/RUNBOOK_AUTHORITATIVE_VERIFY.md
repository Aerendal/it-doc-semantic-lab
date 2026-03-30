# Runbook: Authoritative Verify

**Version:** 1.0  
**Trigger:** Ręczna lub zaplanowana autoryzowana weryfikacja stanu pewności repozytorium  
**Scope:** Całe repo — wszystkie wycinki funkcjonalności obecne w chwili wykonania

---

## Cel

Ten runbook opisuje sposób przeprowadzenia pełnej autoryzowanej weryfikacji repozytorium `itdlab`. Odpowiada na pytanie:

> „W tej chwili, jaki jest rzeczywisty, udokumentowany stan jakości tego repozytorium?"

Autoryzowana weryfikacja nie jest sprawdzeniem deweloperskim. Jest to formalna procedura stosowana przed:
- decyzją o promowaniu do stabilnego repozytorium,
- raportem statusu bramek dla G4 lub G5,
- zatwierdzeniem kamienia milowego,
- zewnętrznym audytem lub przeglądem.

---

## Wymagania wstępne

Przed rozpoczęciem potwierdź wszystkie poniższe punkty:

- [ ] Najnowszy kod jest zatwierdzony i wypchnięty do `origin/main`
- [ ] `git status` nie wykazuje niezatwierdzonych zmian
- [ ] `bin/itdlab` jest zbudowane z bieżącego HEAD (`make build`)
- [ ] SQLite jest w najnowszej wersji schematu (`make db-init` lub sprawdź `schema_version = 1`)
- [ ] Żadne inne uruchomienie `itdlab` nie jest w toku (sprawdź `runs/events.jsonl` pod kątem otwartego uruchomienia)
- [ ] Dostępne miejsce na dysku: ≥ 500 MB wolnego miejsca na partycji zawierającej `reports/`
- [ ] Operator zapoznał się z `docs/SKIP_REGISTER.md` w poszukiwaniu aktywnych pominięć

---

## Procedura

### Faza 1: Migawka środowiska

```bash
# Record binary version
./bin/itdlab version  # or ./bin/itdlab --help | head -1

# Record schema version
sqlite3 db/semantic_index.sqlite "SELECT version, applied_at FROM schema_version;"

# Record event log baseline (line count before this run)
wc -l runs/events.jsonl

# Record git HEAD
git rev-parse HEAD
git log --oneline -1

# Write all of the above to a snapshot file
{
  echo "# Authoritative Verify Snapshot"
  echo "Date: $(date -u +%Y-%m-%dT%H:%M:%SZ)"
  echo "Git HEAD: $(git rev-parse HEAD)"
  echo "Git message: $(git log --oneline -1)"
  echo "Binary: $(./bin/itdlab --version 2>/dev/null || echo 'see help')"
  sqlite3 db/semantic_index.sqlite "SELECT 'Schema version: ' || version FROM schema_version;"
  echo "Event log lines before: $(wc -l < runs/events.jsonl)"
} > av_snapshot.txt
cat av_snapshot.txt
```

### Faza 2: Ponowny ingest ze źródeł

Czysty ingest weryfikuje, że bieżące źródła mogą być w pełni przetworzone bez błędów.

```bash
./bin/itdlab ingest run --source sources/
```

Zatrzymaj się, jeśli kod wyjścia ≠ 0. Zapisz `run_id` ze standardowego wyjścia.

```bash
AV_INGEST_RUN_ID=<run_id from above>
```

### Faza 3: Normalizacja

```bash
./bin/itdlab normalize apply
```

Zatrzymaj się, jeśli kod wyjścia ≠ 0. Zapisz `run_id`.

```bash
AV_NORMALIZE_RUN_ID=<run_id from above>
```

### Faza 4: Sprawdzenie kompletności dowodów dla każdego uruchomienia

```bash
./bin/itdlab audit evidence $AV_INGEST_RUN_ID
./bin/itdlab audit evidence $AV_NORMALIZE_RUN_ID
```

Oczekiwany wynik: oba raporty zwracają `trusted = true` oraz `evidence.complete = true`.

Jeśli którekolwiek zwraca `trusted = false`, zatrzymaj się i zbadaj sprawę przed kontynuacją.

### Faza 5: Uruchomienie wszystkich testów

```bash
make test 2>&1 | tee av_test_output.txt
```

Zapisz kod wyjścia. Niezerowy kod wyjścia oznacza, że co najmniej jeden test nie przeszedł. Nie kontynuuj oceny bramek, jeśli testy nie przeszły.

### Faza 6: Ocena bramek jakości

Dla każdej bramki (G1 do G5, stosownie do ukończonych wycinków funkcjonalności):

```bash
# Check gate status — currently manual; see docs/QUALITY_GATES.md
# For each gate, confirm:
# - all required layers were executed (no unregistered skips)
# - all layers produced trusted runs
# - evidence artifacts exist
```

Podsumowanie statusu bramek:

| Gate | Status | Notes |
|------|--------|-------|
| G1 | PASS / degraded / FAIL / not-evaluated | |
| G2 | PASS / degraded / FAIL / not-evaluated | |
| G3 | PASS / degraded / FAIL / not-evaluated | |
| G4 | PASS / degraded / FAIL / not-evaluated | |
| G5 | PASS / degraded / FAIL / not-evaluated | |

### Faza 7: Sprawdzenie aktywnych pominięć

```bash
cat docs/SKIP_REGISTER.md
```

Dla każdego aktywnego pominięcia:
- Potwierdź, że posiada prawidłowy identyfikator Skip ID, właściciela i datę przeglądu.
- Potwierdź, że jest uwzględnione w powyższej tabeli statusu bramek.
- Potwierdź, że żadne pominięcie Kategorii 3 lub Kategorii 4 nie jest obecne bez zapisu o zatwierdzeniu.

### Faza 8: Sporządzenie raportu autoryzowanej weryfikacji

```bash
REPORT_DIR="reports/av-$(date -u +%Y%m%d-%H%M%S)"
mkdir -p $REPORT_DIR

# Copy snapshot
cp av_snapshot.txt $REPORT_DIR/

# Copy test output
cp av_test_output.txt $REPORT_DIR/

# Write gate status
cat > $REPORT_DIR/gate_status.md << 'EOF'
# Gate Status

Date: <ISO 8601>
Operator: <name>
Git HEAD: <sha>

## Summary

| Gate | Status | Notes |
|------|--------|-------|
| G1   |        |       |
| G2   |        |       |
| G3   |        |       |
| G4   |        |       |
| G5   |        |       |

## Active skips
(copy from SKIP_REGISTER.md)

## Trusted run IDs cited
(list all run_ids used as evidence)
EOF

echo "Authoritative verify report written to: $REPORT_DIR"
```

---

## Warunki zatrzymania

Zatrzymaj się i zapisz jako NIEKOMPLETNE, jeśli:

- Jakiekolwiek uruchomienie wycinka funkcjonalności zakończyło się z niezerowym kodem bez zarejestrowanego wyjaśnienia.
- Jakiekolwiek uruchomienie powołane jako dowód bramki ma `trusted = false`.
- `make test` zakończył się z niezerowym kodem.
- Jakakolwiek warstwa `mock-forbidden` została uruchomiona z mockami.
- `docs/SKIP_REGISTER.md` zawiera pominięcie bez właściciela lub z przeterminowaną datą.

---

## Kryteria ukończenia

Autoryzowana weryfikacja jest kompletna i wiarygodna, gdy:

- [ ] Wszystkie uruchomienia w Fazach 2–3 zakończyły się z kodem 0 lub 2 (brak nieobsłużonych błędów)
- [ ] Wszystkie powołane uruchomienia mają `trusted = true`
- [ ] Wszystkie testy przeszły lub wszystkie niepowodzenia są zarejestrowane w SKIP_REGISTER.md
- [ ] Tabela statusu bramek jest wypełniona z podaniem identyfikatorów uruchomień jako dowodów
- [ ] Raport istnieje w `reports/av-<timestamp>/`
- [ ] Raport jest zatwierdzony: `git add reports/av-<timestamp>/ && git commit -m "audit: authoritative verify <date>"`

---

## Odnośniki wewnętrzne
- `docs/EXECUTION_ASSURANCE_PROGRAM.md`
- `docs/QUALITY_GATES.md`
- `docs/QUALITY_GATES_POLICY.md`
- `docs/EVIDENCE_MODEL.md`
- `docs/SKIP_REGISTER.md`
- `docs/PLAYBOOKS/PLAYBOOK_REAL_PATH_VERIFICATION.md`

## Metadane przeglądu
- Owner: project team
- Status: draft
- Last reviewed: 2026-03-30
