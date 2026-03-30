# Playbook: Real-Path Verification

**Wersja:** 1.0  
**Zakres:** Wszystkie wycinki (capability slices) — stosowany zawsze, gdy wykonywana jest warstwa testów `mock-restricted` lub `mock-forbidden`

---

## Cel

Ten playbook opisuje *jak* przeprowadzać weryfikację rzeczywistej ścieżki (real-path verification): co to oznacza, kiedy ją wykonywać, co uznaje się za zaliczenie, oraz co robić w razie niepowodzenia.

Weryfikacja rzeczywistej ścieżki nie jest tym samym, co uruchamianie testów. Jest to celowe potwierdzenie, że ścieżka kodu odpowiedzialna za zachowanie objęte testami została wykonana na rzeczywistych, niemockowanych zależnościach.

---

## Kiedy stosować ten Playbook

Stosuj, gdy spełniony jest którykolwiek z poniższych warunków:

- Przygotowujesz się do uznania warstwy testów `mock-restricted` za PASS.
- Przygotowujesz się do uznania warstwy testów `mock-forbidden` za PASS.
- Recenzent zażądał dowodów weryfikacji rzeczywistej ścieżki.
- Przygotowujesz pakiet gotowości do promocji dla G4 lub G5.
- Manifest uruchomienia wskazuje `trusted = false`, a przyczyną jest naruszenie polityki mocków.

---

## Zakres weryfikacji rzeczywistej ścieżki

Zgodnie z `docs/POLICY_MOCKS_AND_REAL_PATHS.md`, następujące warstwy wymagają weryfikacji rzeczywistej ścieżki:

| Warstwa | Nazwa | Polityka mocków |
|---------|-------|-----------------|
| 5 | Testy kontraktów systemu plików | mock-restricted |
| 7 | Testy regresji złotych danych wyjściowych | mock-restricted |
| 8 | Testy deterministyczności | mock-restricted |
| 22 | Testy end-to-end wycinka | mock-forbidden |
| 23 | Testy integralności bazy danych | mock-forbidden |
| 25 | Testy odtwarzalności uruchomień | mock-forbidden |
| 26 | Testy kompletności dowodów | mock-forbidden |

Dla warstw `mock-restricted`: przynajmniej jedno pełne wykonanie testów na wycinek musi używać rzeczywistych plików bez mockowania operacji I/O na plikach ani SQLite.

Dla warstw `mock-forbidden`: każde wykonanie musi używać rzeczywistych plików i rzeczywistego SQLite. Bez wyjątków.

---

## Procedura krok po kroku

### Krok 1: Zidentyfikuj docelową warstwę i wycinek

```
Layer: <np. Layer 22 — End-to-end slice tests>
Slice: <np. Slice 1 — Ingest>
Run ID (jeśli weryfikujesz poprzednie uruchomienie): <run_id>
```

### Krok 2: Potwierdź brak aktywnych mocków dla warstwy

Sprawdź pliki testowe dla docelowej warstwy:

```bash
grep -r "t.Skip\|mock\|Mock\|stub\|Stub\|FakeDB\|InMemory" \
  <test_file_or_directory>
```

Oczekiwany wynik: Brak konfiguracji mocków dla testowanych komponentów. Jeśli znajdziesz mocki:
- Sprawdź, czy dotyczą wyłącznie *zewnętrznych* zależności (dopuszczalne dla `mock-restricted`).
- Jeśli owijają SQLite, system plików lub dziennik zdarzeń → jest to naruszenie dla warstw `mock-forbidden`.
- Nie kontynuuj, dopóki mocki nie zostaną usunięte lub warstwa nie zostanie przeklasyfikowana.

### Krok 3: Przygotuj środowisko testów rzeczywistej ścieżki

```bash
# Create a fresh test directory
mkdir -p /tmp/itdlab-realpath-test
cd /tmp/itdlab-realpath-test

# Copy real source fixtures
cp -r <repo>/sources/  ./sources/

# Initialise a fresh database (no pre-existing state)
<repo>/bin/itdlab --db ./semantic_index.sqlite --log ./events.jsonl ingest run --source ./sources/
```

Potwierdź:
- `semantic_index.sqlite` jest rzeczywistym plikiem na dysku (nie `:memory:`).
- `events.jsonl` jest rzeczywistym plikiem na dysku.
- Źródła są rzeczywistymi plikami, nie dynamicznie generowanymi zaślepkami.

### Krok 4: Uruchom testy docelowej warstwy na rzeczywistych zależnościach

```bash
# From repo root
DB_PATH=/tmp/itdlab-realpath-test/semantic_index.sqlite \
LOG_PATH=/tmp/itdlab-realpath-test/events.jsonl \
go test ./... -run <TestPattern> -v 2>&1 | tee /tmp/itdlab-realpath-test/test_output.txt
```

### Krok 5: Zweryfikuj artefakty dowodowe

Po zakończeniu testów:

```bash
# Check that SQLite was actually used (non-empty, not in-memory artefact)
ls -lh /tmp/itdlab-realpath-test/semantic_index.sqlite

# Check that events were appended
wc -l /tmp/itdlab-realpath-test/events.jsonl

# Check that reports were produced
ls -lh reports/<run_id>/
```

Wszystkie poniższe warunki muszą być spełnione:
- [ ] Plik SQLite istnieje i jest większy niż bazowy pusty schemat
- [ ] Co najmniej 1 wiersz zdarzenia w dzienniku zdarzeń
- [ ] Plik `run_manifest.json` istnieje i zawiera `evidence.complete = true`
- [ ] Plik `run_manifest.json` zawiera `trusted = true`

### Krok 6: Zapisz dowody weryfikacji rzeczywistej ścieżki

W katalogu `reports/<run_id>/`:
```
real_path_verification.md  — zapisz: warstwę, wycinek, polecenie testowe, wynik, znacznik czasu
```

Szablon:
```md
## Real-Path Verification Record

- Layer: <nazwa warstwy>
- Slice: <nazwa wycinka>
- Run ID: <run_id>
- Test command: `<pełne polecenie>`
- SQLite path: <ścieżka>
- Event log path: <ścieżka>
- Result: PASS / FAIL
- Verified at: <data/czas ISO 8601>
- Verified by: <właściciel>
- Notes: <wszelkie odchylenia lub obserwacje>
```

---

## Warunki zatrzymania

Zatrzymaj się i nie rejestruj wyniku PASS, jeśli spełniony jest którykolwiek z poniższych warunków:

- Jakikolwiek komponent `mock-forbidden` był owinięty w mocka lub zaślepkę.
- SQLite było otwarte jako `:memory:`.
- Pliki źródłowe były generowane dynamicznie bez weryfikacji.
- Środowisko testowe było współdzielone z działającym serwerem deweloperskim, który mógł modyfikować stan.
- Plik `run_manifest.json` wskazuje `trusted = false`.

---

## Odzyskiwanie po niepowodzeniu weryfikacji rzeczywistej ścieżki

1. Zapisz niepowodzenie w manifeście uruchomienia (`evidence.complete = false` jeśli brakuje artefaktów, lub `trusted = false` jeśli wykryto naruszenie polityki).
2. Zarejestruj pominięcie w `docs/SKIP_REGISTER.md`, jeśli warstwa nie może zostać wykonana z udokumentowanego powodu.
3. Nie przechodź do oceny bram G4/G5 dopóki warstwa nie przechodzi lub nie jest zarejestrowane pominięcie kategorii 1/2 z właścicielem i datą przeglądu.

---

## Odwołania wewnętrzne
- `docs/POLICY_MOCKS_AND_REAL_PATHS.md`
- `docs/TESTING_STANDARD.md`
- `docs/TEST_CATALOG.md`
- `docs/EVIDENCE_MODEL.md`
- `docs/SKIP_REGISTER.md`

## Metadane przeglądu
- Właściciel: zespół projektu
- Status: szkic
- Ostatni przegląd: 2026-03-30
