# Playbook: Source Ingest

## Cel

Definiuje strategię przyjmowania (ingestowania) surowych źródeł dokumentacji IT do laboratorium semantycznego.

---

## Kiedy stosować ten Playbook

- Dodawanie nowych dokumentów źródłowych do `sources/`
- Ponowne ingestowanie po aktualizacji plików źródłowych
- Diagnozowanie nieudanego uruchomienia ingestowania

---

## Strategia ingestowania

### 1. Umieszczanie źródeł

Umieść surowe pliki Markdown w odpowiednim podkatalogu:

```
sources/
  matrices/    — macierze dokumentów branżowych
  metagraph/   — definicje metagrafu
  plans/       — plany projektu i mapy drogowe
```

Pliki muszą być zakodowane w UTF-8 Markdown. Bez BOM. Bez załączników binarnych.

### 2. Walidacja przed ingestowaniem

Przed uruchomieniem ingestowania:
- Uruchom testy Layer 1–5 (walidacja kontraktu i danych wejściowych)
- Zweryfikuj kodowania plików
- Upewnij się, że każdy plik zawiera co najmniej jeden nagłówek

### 3. Uruchomienie ingestowania

```
itdlab ingest run --source sources/
```

Polecenie:
1. Wykrywa wszystkie pliki `.md` pod podaną ścieżką źródłową
2. Parsuje każdy plik (nagłówki, sekcje, metadane)
3. Zapisuje wiersz `Document` do SQLite dla każdego pliku
4. Dołącza zdarzenia `ingested` do `runs/events.jsonl`
5. Zapisuje `source_manifest.json` i `parse_report.json` do `reports/<run_id>/`

### 4. Weryfikacja po ingestowaniu

```
itdlab ingest inspect <path>   # przegląd wyniku parsowania konkretnego pliku
```

Sprawdź:
- Wszystkie oczekiwane dokumenty są widoczne w SQLite
- Brak błędów parsowania w `parse_report.json`
- Dziennik zdarzeń zawiera wpisy dla wszystkich zingestowanych plików

---

## Tryby awaryjne

| Objaw | Prawdopodobna przyczyna | Działanie |
|-------|------------------------|-----------|
| Plik nieobecny w manifeście źródłowym | Plik nie znajduje się w ścieżce źródłowej | Zweryfikuj umieszczenie |
| Błąd kodowania | BOM lub kodowanie inne niż UTF-8 | Konwertuj za pomocą `iconv` |
| Brak nagłówków | Plik nie zawiera nagłówków `#` | Sprawdź strukturę Markdown |
| Niepowodzenie Gate 1 | Naruszenie kontraktu | Popraw plik źródłowy, uruchom ponownie |
