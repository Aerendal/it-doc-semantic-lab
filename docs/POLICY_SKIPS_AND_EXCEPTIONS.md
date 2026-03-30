# Policy: Skips and Exceptions

Dokument definiuje, kiedy warstwa testowa, sprawdzenie bramki jakości lub wymóg dowodowy może zostać pominięty, jak zarejestrować wyjątek oraz jak pominięcia wpływają na wiarygodność bramki.

---

## Zasada

**Pominięcie to udokumentowana decyzja, nie wygoda.**

Ciche pominięcie — niezarejestrowane, nieprzejrzane i nieograniczone czasowo — jest traktowane równoważnie z niepowodzeniem bramki. Usuwa wiarygodność każdej bramki zależącej od pominiętej warstwy.

---

## Kategorie pominięć

### Kategoria 1: Dopuszczalne pominięcie — infrastruktura niedostępna

Warstwa może zostać pominięta bez wpływu na bramkę, jeśli:
- infrastruktura, od której zależy, jest rzeczywiście niedostępna w bieżącym środowisku (np. brak sieci, brak usługi zewnętrznej),
- pominięcie jest zarejestrowane zgodnie z poniższą procedurą,
- pominięta warstwa **nie** leży na ścieżce krytycznej ocenianej bramki.

**Przykład:** Layer 28 (Performance Budget Tests) może zostać pominięty w środowisku CI o ograniczonych zasobach, jeśli budżet wydajnościowy nie jest jeszcze formalnie określony.

### Kategoria 2: Dopuszczalne pominięcie — warstwa jeszcze niezaimplementowana

Warstwa może zostać pominięta bez wpływu na bramkę, jeśli:
- jest jawnie oznaczona jako `not yet implemented` w `docs/TEST_CATALOG.md`,
- bramka zależna od niej nie jest oceniana w tym przebiegu,
- pominięcie jest zarejestrowane.

### Kategoria 3: Ograniczone pominięcie — wymaga wyraźnej zgody

Warstwa może zostać pominięta z potwierdzonym wpływem na bramkę, jeśli:
- warstwa leży na ścieżce krytycznej aktywnej bramki,
- pominięcie jest zarejestrowane z właścicielem i datą przeglądu,
- wynik bramki jest oznaczany jako `degraded` (nie `PASS`) do czasu przywrócenia warstwy.

**Warstwy, których nie można pominąć bez zgody:**
- Layer 1 (File Presence), Layer 5 (Source Schema) — G1
- Layer 8 (Golden Extraction), Layer 10 (Determinism) — G2
- Layer 18 (Explainability), Layer 19 (Acyclicity) — G4
- Layer 27 (Evidence Pack) — G5

### Kategoria 4: Zakazane pominięcie — nigdy niedopuszczalne

Poniższe elementy **nigdy** nie mogą być pominięte pod żadnym pozorem:

| Element | Powód |
|---------|-------|
| Produkcja paczki dowodowej (Layer 27) | Bez dowodów przebieg nie może być audytowany |
| Rejestrowanie kodu wyjścia | Wymagane do wszystkich ocen bramek |
| Dopisywanie do dziennika zdarzeń | Wymagane dla odtwarzalności i audytu |
| Tworzenie rekordu przebiegu w SQLite | Wymagane do śledzenia przebiegów |

Próba pominięcia elementu Kategorii 4 skutkuje automatycznym niepowodzeniem bramki dla wszystkich bramek w przebiegu.

---

## Rejestracja wyjątku

Każde pominięcie (Kategorie 1–3) musi być zarejestrowane. Niezarejestrowane pominięcie jest traktowane jako Kategoria 4.

### Wymagane pola

```
skip_id:       <unique identifier, e.g., SKIP-2026-001>
layer:         <layer number and name, e.g., Layer 28 — Performance Budget Tests>
category:      <1 | 2 | 3>
reason:        <brief factual description — what is missing and why>
owner:         <name or role responsible for this exception>
registered_at: <ISO 8601 date>
review_date:   <ISO 8601 date — when this skip must be reviewed>
gate_impact:   <none | degraded:<gate_id>>
```

### Miejsce rejestracji

Pominięcia są rejestrowane w `docs/SKIP_REGISTER.md`. Każde pominięcie to jeden wpis w tym pliku.

---

## Wygaśnięcie pominięcia

Pominięcie wygasa w dniu `review_date`. Po wygaśnięciu:

1. Właściciel musi przywrócić warstwę lub odnowić pominięcie z nową datą `review_date`.
2. Wygasłe pominięcie bez odnowienia jest automatycznie traktowane jako Kategoria 4 (zakazane).
3. Wpis w `SKIP_REGISTER.md` musi zostać zaktualizowany o wynik.

---

## Wpływ pominięć na bramki

| Kategoria | Wpływ na bramkę |
|-----------|----------------|
| 1 (infrastruktura niedostępna, niekrytyczna) | Brak — bramka ocenia pozostałe warstwy |
| 2 (jeszcze niezaimplementowana, niekrytyczna) | Brak — bramka ocenia pozostałe warstwy |
| 3 (zatwierdzone, ścieżka krytyczna) | Bramka oznaczona jako `degraded` — niedopuszczalna do promocji |
| 4 (zakazane) lub niezarejestrowane | Automatyczne niepowodzenie bramki dla wszystkich aktywnych bramek |

Wynik bramki `degraded` oznacza:
- przebieg nie może być promowany do stabilnego repozytorium,
- status bramki musi być zapisany jako `degraded` w podsumowaniu przebiegu,
- stan degradacji musi być widoczny w wyjściu `audit evidence`.

---

## Co oznacza „wiarygodna zieleń"

Wynik bramki `PASS` jest wiarygodny tylko wtedy, gdy:
1. Wszystkie wymagane warstwy testowe zostały wykonane (nie pominięte).
2. Wszystkie zarejestrowane pominięcia są Kategorii 1 lub 2 (niekrytyczne).
3. Żaden element Kategorii 4 nie został ominięty.
4. Paczka dowodowa jest kompletna.

Wynik niespełniający tych warunków musi być raportowany jako `degraded`, nie `PASS`.
