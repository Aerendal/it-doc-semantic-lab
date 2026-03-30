# Playbook: Schema Evolution

## Cel

Definiuje strategię ewolucji schematu SQLite bez przerywania pracy na istniejących danych ani uruchomieniach.

---

## Zasady

1. Zmiany schematu są **zawsze addytywne** — brak usuwania kolumn, brak zmiany nazw w miejscu
2. Każda wersja schematu posiada własny plik DDL: `db/schema_v<N>.sql`
3. Tabela `schema_version` śledzi zastosowane wersje
4. Testy migracji (Layer 15) muszą przejść przed zastosowaniem jakiejkolwiek migracji

---

## Dodawanie nowej kolumny

```sql
-- db/schema_v2.sql (przykład)
ALTER TABLE documents ADD COLUMN confidence REAL NOT NULL DEFAULT 0.0;
INSERT OR IGNORE INTO schema_version (version, applied_at, description)
VALUES (2, datetime('now'), 'add confidence column to documents');
```

Zasady:
- Nowe kolumny muszą posiadać wartość `DEFAULT`
- Nigdy `NOT NULL` bez wartości domyślnej w istniejącej tabeli
- Dodaj indeks, jeśli kolumna będzie używana w zapytaniach

---

## Dodawanie nowej tabeli

Dodaj nową instrukcję `CREATE TABLE IF NOT EXISTS` do nowego pliku schematu. Nowe tabele nie wymagają migracji danych.

---

## Proces migracji

1. Napisz `db/schema_v<N>.sql`
2. Napisz test migracji (Layer 15) używając migawki SQLite sprzed migracji
3. Uruchom `make test` — zweryfikuj, że test migracji przechodzi
4. Zastosuj: `sqlite3 db/semantic_index.sqlite < db/schema_v<N>.sql`
5. Zaktualizuj `internal/adapters/sqlite/schema.go` tak, aby zawierał nowe DDL

---

## Czego nigdy nie robić

- Nie używaj `DROP COLUMN` ani `DROP TABLE` na istniejących wersjach schematu
- Nie zmieniaj nazw kolumn (użyj nowej kolumny + zdeprecjonuj starą)
- Nie zmieniaj ograniczeń `CHECK` na istniejących kolumnach w miejscu
