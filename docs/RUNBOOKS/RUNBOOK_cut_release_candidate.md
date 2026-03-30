# Runbook: Cut Release Candidate

## Kiedy używać

Gdy stabilny zestaw metadanych semantycznych jest gotowy do promowania do repozytorium referencyjnego i potrzebna jest wersjonowana migawka.

## Wymagania wstępne

- Wszystkie bramki jakości zaliczone
- Eksport do repo 1 zakończony pomyślnie
- Wszystkie testy zaliczone (`make test`)

## Kroki

### 1. Uruchom pełny zestaw testów

```sh
make test
```

Wszystkie testy muszą przejść. Zero błędów.

### 2. Zweryfikuj pakiet dowodów

```sh
./bin/itdlab audit evidence <run_id>
```

Wszystkie wymagane artefakty muszą być obecne.

### 3. Oznacz wydanie w tym repozytorium

```sh
git tag -a "rc/<date>-<run_id>" -m "Release candidate: <date> run <run_id>"
git push origin "rc/<date>-<run_id>"
```

### 4. Zapisz informację o wydaniu w manifeście eksportu

```sh
cat reports/<run_id>/export_manifest.json
```

Zapisz kopię manifestu do `reports/rc/<date>/manifest.json`.

### 5. Oznacz repozytorium referencyjne

W `IT-Dokumentacja`:

```sh
git tag -a "semantic-rc/<date>" -m "Semantic metadata from lab rc/<date>-<run_id>"
git push origin "semantic-rc/<date>"
```

## Warunki zatrzymania

- Jakikolwiek błąd testu → napraw przed tworzeniem RC
- Błąd bramki → napraw przed tworzeniem RC
- Niekompletny pakiet dowodów → najpierw go uzupełnij
