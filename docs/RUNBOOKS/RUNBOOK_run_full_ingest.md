# Runbook: Run Full Ingest

## Wymagania wstępne

- `bin/itdlab` zbudowane (`make build`)
- `db/semantic_index.sqlite` zainicjalizowane (`make db-init`)
- Pliki źródłowe obecne w katalogu `sources/`

## Kroki

### 1. Sprawdź, czy pliki źródłowe istnieją

```sh
ls sources/
```

Oczekiwany wynik: co najmniej jeden plik `.md` w podkatalogu.

### 2. Uruchom ingest

```sh
./bin/itdlab ingest run --source sources/
```

### 3. Zweryfikuj wpisy w dzienniku zdarzeń

```sh
grep '"step":"ingest"' runs/events.jsonl | wc -l
```

Oczekiwany wynik: liczba równa liczbie przetworzonych plików.

### 4. Sprawdź wybrany dokument

```sh
./bin/itdlab ingest inspect sources/<path-to-file>.md
```

### 5. Sprawdź bazę danych

```sh
sqlite3 db/semantic_index.sqlite "SELECT count(*) FROM documents;"
```

Oczekiwany wynik: niezerowa liczba wierszy.

### 6. Przejrzyj raport parsowania

```sh
cat reports/<run_id>/parse_report.json
```

## Warunki zatrzymania

- Niezerowy kod wyjścia → sprawdź `reports/<run_id>/stdout.txt` i dziennik zdarzeń
- Zero wierszy w tabeli `documents` → nie znaleziono plików `.md`; sprawdź ścieżkę źródłową
- Błąd G1 → napraw pliki źródłowe zgodnie z `gate_failures.json` i uruchom ponownie
