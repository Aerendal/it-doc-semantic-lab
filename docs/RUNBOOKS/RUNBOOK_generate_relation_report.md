# Runbook: Generate Relation Report

## Wymagania wstępne

- Normalizacja zakończona pomyślnie
- Tabela `documents` zawiera wiersze ze statusem `status = 'normalized'`
- Reguły relacji istnieją w tabeli `relation_rules`

## Kroki

### 1. Wyświetl wszystkie relacje

```sh
./bin/itdlab relations show
```

### 2. Wyświetl relacje dla konkretnego dokumentu

```sh
./bin/itdlab relations show --doc <document_id>
```

### 3. Wyjaśnij konkretną relację

```sh
./bin/itdlab relations explain --rel <relation_id>
```

### 4. Sprawdź liczbę relacji w bazie danych

```sh
sqlite3 db/semantic_index.sqlite "SELECT type, count(*) FROM relations GROUP BY type;"
```

### 5. Sprawdź relacje bez wyjaśnienia

```sh
sqlite3 db/semantic_index.sqlite \
  "SELECT id, from_id, to_id FROM relations WHERE explanation = '' OR explanation IS NULL;"
```

Oczekiwany wynik: zero wierszy. (G4 blokuje, jeśli jakiekolwiek istnieją.)

### 6. Sprawdź cykle w zależnościach

```sh
# Ręczna weryfikacja dla małych grafów:
sqlite3 db/semantic_index.sqlite \
  "SELECT from_id, to_id FROM relations WHERE type = 'depends_on';"
```

## Warunki zatrzymania

- Relacje z pustym wyjaśnieniem → napraw reguły wnioskowania i uruchom ponownie
- Wykryte cykle w depends_on → przejrzyj reguły i strukturę dokumentów
