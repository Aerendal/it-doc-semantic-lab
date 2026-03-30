# Runbook: Bootstrap Local Environment

## Wymagania wstępne

- Go 1.21+
- `make`
- `sqlite3` CLI (do `db-init` i inspekcji)
- Git

## Kroki

### 1. Sklonuj repozytorium i wejdź do katalogu

```sh
cd /path/to/it-doc-semantic-lab
```

### 2. Uporządkuj zależności

```sh
make tidy
```

### 3. Zbuduj CLI

```sh
make build
```

Plik binarny zostanie umieszczony w `bin/itdlab`.

### 4. Zainicjalizuj bazę danych

```sh
make db-init
```

Tworzy `db/semantic_index.sqlite` ze schematem v1.

### 5. Utwórz katalogi robocze (jeśli nie istnieją)

```sh
mkdir -p runs reports normalized
```

### 6. Zweryfikuj poprawność kompilacji

```sh
./bin/itdlab --help
```

Oczekiwany wynik: tekst pomocy zawierający polecenia `ingest`, `normalize`, `classify`, `relations`, `export`, `audit`.

### 7. Uruchom skrócone testy

```sh
make test-short
```

Oczekiwany wynik: wszystkie testy zaliczone.

## Warunki zatrzymania

- Jeśli `make tidy` zakończy się błędem: sprawdź połączenie z internetem potrzebne do pobrania modułów Go
- Jeśli `make build` zakończy się błędem: sprawdź wersję Go (`go version` musi wynosić 1.21+)
- Jeśli `make db-init` zakończy się błędem: sprawdź, czy `sqlite3` jest zainstalowane (`which sqlite3`)
