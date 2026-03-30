# Runbook: Rebuild SQLite

## Kiedy używać

- Plik bazy danych jest uszkodzony lub nie istnieje
- Migracja schematu nie zakończyła się pomyślnie
- Świeży start po znaczących zmianach strukturalnych

## Wymagania wstępne

- Istnieje plik `db/schema_v1.sql`
- Dziennik zdarzeń `runs/events.jsonl` jest nienaruszony (do ewentualnego odtworzenia)

## Kroki

### 1. Utwórz kopię zapasową istniejącej bazy (jeśli jest odtwarzalna)

```sh
cp db/semantic_index.sqlite db/semantic_index.sqlite.bak
```

### 2. Usuń uszkodzoną bazę danych

```sh
rm db/semantic_index.sqlite
```

### 3. Ponownie zainicjalizuj schemat

```sh
make db-init
```

### 4. Ponownie przetworz źródła

```sh
./bin/itdlab ingest run --source sources/
```

### 5. Ponownie uruchom normalizację

```sh
./bin/itdlab normalize apply
```

### 6. Ponownie uruchom relacje

```sh
./bin/itdlab relations show
```

### 7. Zweryfikuj zgodność liczby wierszy z kopią zapasową (jeśli dostępna)

```sh
sqlite3 db/semantic_index.sqlite "SELECT count(*) FROM documents;"
sqlite3 db/semantic_index.sqlite.bak "SELECT count(*) FROM documents;"
```

## Uwagi

- Dziennik zdarzeń jest śladem audytowym — nie wymaga odbudowy.
- Jeśli dziennik zdarzeń również zaginął, ponownie przetworz oryginalne pliki źródłowe.
- Nie dodawaj `db/semantic_index.sqlite` do gita (jest objęty gitignore).
