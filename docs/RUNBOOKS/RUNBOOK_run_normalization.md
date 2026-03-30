# Runbook: Run Normalization

## Wymagania wstępne

- Ingest zakończony pomyślnie
- Tabela `documents` zawiera wiersze ze statusem `status = 'ingested'`

## Kroki

### 1. Podgląd normalizacji

```sh
./bin/itdlab normalize preview
```

Przejrzyj wynik — sprawdź nieoczekiwane zmiany canonical ID lub kolizje.

### 2. Zastosuj normalizację

```sh
./bin/itdlab normalize apply
```

### 3. Zweryfikuj canonical IDs

```sh
sqlite3 db/semantic_index.sqlite "SELECT id, canonical_id, raw_name FROM documents;"
```

### 4. Sprawdź kolizje

```sh
sqlite3 db/semantic_index.sqlite \
  "SELECT canonical_id, count(*) as n FROM documents GROUP BY canonical_id HAVING n > 1;"
```

Oczekiwany wynik: zero wierszy. Jeśli kolizje istnieją, rozwiąż je zgodnie z PLAYBOOK_normalization.md przed kontynuacją.

### 5. Przejrzyj raport normalizacji

```sh
cat reports/<run_id>/normalization_report.json
```

## Warunki zatrzymania

- Liczba kolizji > 0 → rozwiąż przed eksportem (G3)
- Nieoczekiwane zmiany canonical ID → sprawdź ponownie reguły normalizacji
