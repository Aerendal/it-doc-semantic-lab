# Runbook: Export Metadata to Repo 1

## Wymagania wstępne

- Wszystkie bramki jakości (1–5) zostały zaliczone
- Pakiet dowodów jest kompletny
- Uruchomienie ma status `status = 'completed'` w SQLite
- Repozytorium referencyjne jest przełączone na właściwą gałąź

## Kroki

### 1. Zweryfikuj status bramek

```sh
sqlite3 db/semantic_index.sqlite \
  "SELECT run_id, status FROM runs ORDER BY started_at DESC LIMIT 1;"
```

Oczekiwany wynik: `status = 'completed'`

### 2. Uruchom eksport (najpierw dry-run)

```sh
./bin/itdlab export repo1 --target ../IT-Dokumentacja/ --dry-run
```

Przejrzyj listę plików przeznaczonych do promowania.

### 3. Zastosuj eksport

```sh
./bin/itdlab export repo1 --target ../IT-Dokumentacja/
```

### 4. Przejrzyj manifest eksportu

```sh
cat reports/<run_id>/export_manifest.json
```

### 5. Zatwierdź zmiany w repozytorium referencyjnym

```sh
cd ../IT-Dokumentacja
git add .
git status
git commit -m "chore: promote semantic metadata from lab run <run_id>"
```

## Warunki zatrzymania

- Kod wyjścia 2 → błąd bramki jakości; sprawdź `gate_failures.json`
- Pliki już aktualne → brak zmian do promowania
- Katalog docelowy nie istnieje → sprawdź ścieżkę podaną w `--target`
