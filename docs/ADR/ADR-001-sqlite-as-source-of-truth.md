# ADR-001: SQLite as Source of Truth

**Status:** Zaakceptowany  
**Deciders:** Zespół projektowy  
**Date:** 2026-03-29

---

## Kontekst

Laboratorium semantyczne potrzebuje odpytywalnego, trwałego magazynu dla dokumentów, sekcji, relacji, normalizacji i historii uruchomień.  
Rozważane opcje: pliki płaskie (YAML/JSON), PostgreSQL, Neo4j, SQLite.

## Decyzja

Użyj **SQLite** jako jedynego źródła prawdy dla całego stanu laboratorium.

Uzupełnij go **append-only logiem zdarzeń JSONL** (`runs/events.jsonl`) na potrzeby audytu, odtwarzalności i historii uruchomień — lecz SQLite jest autorytatywny dla bieżącego stanu.

## Uzasadnienie

| Kryterium | SQLite | Alternatywy |
|-----------|--------|-------------|
| Zero-config | ✅ plikowy | ❌ PostgreSQL wymaga serwera |
| Odpytywalny | ✅ pełny SQL | ❌ płaski JSON nieodpytywalny |
| Relacje grafowe | ✅ oparte na JOIN | Neo4j to przerost w tej skali |
| CGO-free | ✅ `modernc.org/sqlite` | `mattn/go-sqlite3` wymaga CGO |
| Przenośny binarny | ✅ | ❌ bazy oparte na serwerze |
| Audytowalność | ✅ tryb WAL + log zdarzeń | |

## Konsekwencje

- Wszystkie zapisy trafiają najpierw do SQLite, a następnie do logu zdarzeń.
- Log JSONL może odtworzyć stan SQLite od podstaw (gwarancja odtwarzalności).
- Migracje schematu korzystają z tabeli `schema_version`. Każda wersja to addytywny plik DDL.
- Neo4j lub inne magazyny grafowe mogą zostać dodane później jako **widoki tylko do odczytu** na eksportach SQLite — nigdy jako główny magazyn.

## Odrzucone alternatywy

- **PostgreSQL** — wymaga zewnętrznego serwera; zbędny dla lokalnych narzędzi laboratoryjnych.
- **Neo4j** — odpowiedni do przechodzenia grafów w skali; tutaj przedwczesny.
- **Płaskie YAML/JSON** — nieodpytywalne, brak integralności relacyjnej.
