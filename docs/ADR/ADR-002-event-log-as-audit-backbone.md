# ADR-002: Event Log as Audit Backbone

| Field | Value |
|-------|-------|
| Status | Zaakceptowany |
| Date | 2026-03-30 |
| Deciders | zespół projektowy |
| Supersedes | — |
| Superseded by | — |

---

## Kontekst

System przetwarza dokumenty, wnioskuje relacje i normalizuje nazwy w trakcie kolejnych uruchomień. Każde uruchomienie modyfikuje stan SQLite. Bez niezależnego zapisu tego, co się stało, kiedy i w jakiej kolejności, nie ma możliwości:
- odtworzenia historii, w jaki sposób dokument osiągnął swój bieżący stan,
- porównania dwóch uruchomień w celu wykrycia regresji,
- audytu uruchomienia po fakcie bez jego ponownego wykonania,
- wykrycia, czy stan SQLite został zmodyfikowany poza zarządzanym uruchomieniem.

SQLite jest źródłem prawdy dla bieżącego stanu, ale natywnie nie śledzi stanu historycznego ani proweniencji między uruchomieniami. Coś innego musi pełnić rolę szkieletu audytu.

---

## Decyzja

Każda operacja zmieniająca stan w `itdlab` dołącza ustrukturyzowane zdarzenie JSON do append-only pliku logu zdarzeń (`runs/events.jsonl`). Log zdarzeń jest:

1. **Append-only** — istniejące linie nigdy nie są modyfikowane ani usuwane.
2. **Ustrukturyzowany** — każda linia to prawidłowy obiekt JSON z obowiązkowymi polami.
3. **Powiązany z uruchomieniem** — każde zdarzenie niesie `run_id` wiążący je z konkretnym wywołaniem.
4. **Niezależny od SQLite** — log jest zapisywany do osobnego pliku; nie zależy od transakcji SQLite.

### Obowiązkowe pola zdarzeń

```json
{
  "ts": "<ISO 8601 UTC>",
  "run_id": "<string>",
  "step": "<string>",
  "entity": "<string>",
  "entity_id": "<string>",
  "action": "<string>",
  "before": "<any | null>",
  "after": "<any | null>"
}
```

Dodatkowe pola są dozwolone. Obowiązkowy zestaw musi być zawsze obecny.

---

## Rozważane alternatywy

### 1. Tylko SQLite — tabele historii

**Podejście:** Dodaj tabele cieni `_history` przechwytujące poprzednie wartości przy UPDATE/DELETE.

**Odrzucono, ponieważ:**
- Wyzwalacze w SQLite są kruche i trudne do testowania.
- Tabele historii w tym samym pliku co bieżący stan tworzą pojedynczy punkt awarii.
- Odczyt historii uruchomień wymaga złączeń SQL przez wiele tabel; płaski log jest łatwiejszy do strumieniowania i przeszukiwania.
- Tabele historii nie przetrwają, jeśli baza danych zostanie odbudowana od podstaw.

### 2. Write-ahead log (WAL) jako zapis audytowy

**Podejście:** Użyj pliku WAL SQLite jako zapisu zdarzeń.

**Odrzucono, ponieważ:**
- Pliki WAL są checkpointowane i nadpisywane; nie stanowią trwałego śladu audytowego.
- Format WAL jest binarny; nie jest czytelny dla człowieka bez narzędzi.
- WAL jest specyficzny dla implementacji i nie jest częścią publicznego API.

### 3. Wyłącznie logowanie strukturalne (stdout)

**Podejście:** Emituj ustrukturyzowane linie logu na stdout; przekieruj do pliku przez wywołującego.

**Odrzucono, ponieważ:**
- Przechwytywanie stdout zależy od wywołania powłoki; może zostać utracone.
- Stdout miesza wyjście diagnostyczne ze zdarzeniami audytowymi; filtrowanie jest kruche.
- Stdout nie jest z założenia append-only; może zostać obcięty.

---

## Konsekwencje

### Pozytywne

- Kompletna historia wszystkich zmian stanu jest dostępna dla każdego uruchomienia, w tym uruchomień zakończonych niepowodzeniem lub przerwanych.
- Uruchomienie można odtworzyć z logu zdarzeń bez dostępu do bieżącego stanu SQLite.
- Dwa uruchomienia można porównać, zestawiając ich wycinki logu zdarzeń.
- Log zdarzeń można odtworzyć, aby wykryć rozbieżność SQLite z oczekiwanym stanem.
- Analiza śledcza po naruszeniu kontraktu (exit code 3) jest możliwa nawet jeśli SQLite jest uszkodzony.

### Negatywne / zaakceptowane kompromisy

- Log zdarzeń rośnie w nieskończoność. Obecnie nie ma polityki kompaktowania.
- Dołączanie do JSONL przy każdym zdarzeniu jest sekwencyjną operacją I/O; stanowi wąskie gardło przy uruchomieniach o dużej objętości.
- Log zdarzeń i SQLite mogą tymczasowo się rozbiec między zapisem do SQLite a kolejnym dołączeniem do logu. Jest to zaakceptowane; niezmiennik (I1 w `EXECUTION_CONTRACT.md`) wymaga zapisu do SQLite przed dołączeniem do logu, a nie atomowości obu operacji.

### Odroczone

- Narzędzie do odtwarzania między uruchomieniami (`itdlab audit replay`) nie jest zaimplementowane w v1.
- Polityka kompaktowania / archiwizacji jest odroczona.
- Podpisywanie logu zdarzeń / wykrywanie manipulacji jest odroczone.

---

## Uwagi implementacyjne

- Log zdarzeń jest zapisywany przez `internal/adapters/jsonl/event_log.go`.
- Dołączanie chronione muteksem; bezpieczne dla współbieżnych goroutine w ramach jednego uruchomienia.
- Każde uruchomienie tworzy osobny wycinek zdarzeń identyfikowalny przez `run_id`.
- `itdlab audit runs` odczytuje log zdarzeń, aby wylistować historyczne uruchomienia.

---

## Odwołania wewnętrzne
- `docs/EXECUTION_ASSURANCE_PROGRAM.md`
- `docs/EXECUTION_CONTRACT.md`
- `docs/EVIDENCE_MODEL.md`
- `docs/ADR/ADR-001-sqlite-as-source-of-truth.md`

## Metadane przeglądu
- Owner: zespół projektowy
- Status: zaakceptowany
- Last reviewed: 2026-03-30
