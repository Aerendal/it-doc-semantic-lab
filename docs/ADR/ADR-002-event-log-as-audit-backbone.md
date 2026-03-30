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
- odtworzenia historii prowadzącej do bieżącego stanu,
- porównania dwóch uruchomień w celu wykrycia regresji,
- audytu uruchomienia po fakcie bez jego ponownego wykonania,
- wykrycia, czy stan SQLite został zmodyfikowany poza zarządzanym uruchomieniem,
- przeprowadzenia przeglądu po naruszeniu kontraktu lub po awarii częściowej.

SQLite jest źródłem prawdy dla stanu bieżącego, ale nie jest sam w sobie wystarczającym nośnikiem proweniencji między uruchomieniami. Potrzebna jest osobna, append-only warstwa audytowa.

---

## Decyzja

Każda operacja zmieniająca stan w `itdlab` dołącza ustrukturyzowane zdarzenie JSON do append-only pliku logu zdarzeń `runs/events.jsonl`.

Log zdarzeń jest:
1. **append-only** — istniejące linie nigdy nie są modyfikowane ani usuwane,
2. **ustrukturyzowany** — każda linia jest prawidłowym obiektem JSON,
3. **powiązany z uruchomieniem** — każde zdarzenie niesie `run_id`,
4. **niezależny od SQLite** — log nie jest utrzymywany wewnątrz pliku bazy danych,
5. **obowiązkowy dla state-changing runs** — brak logu oznacza brak pełnego evidence model.

### Obowiązkowe pola zdarzenia

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

Dodatkowe pola są dozwolone, lecz zestaw podstawowy jest obowiązkowy dla wszystkich zdarzeń należących do kontraktu audytowego.

---

## Rozważane alternatywy

### Alternatywa A — Tylko SQLite z tabelami historii

**Odrzucona.**

Powody:
- wyzwalacze i shadow tables komplikują projekt i testowanie,
- historia w tym samym pliku co stan bieżący tworzy jeden punkt awarii,
- inspekcja sekwencji zdarzeń jest trudniejsza niż w płaskim logu,
- odbudowa historii po uszkodzeniu bazy jest trudniejsza.

### Alternatywa B — WAL SQLite jako zapis audytowy

**Odrzucona.**

Powody:
- WAL nie jest trwałym, docelowym śladem audytowym,
- jest binarny i trudny do ręcznej inspekcji,
- checkpointing nadpisuje jego znaczenie jako długoterminowego audytu,
- jest artefaktem implementacyjnym, nie kontraktem repozytorium.

### Alternatywa C — Wyłącznie ustrukturyzowany stdout

**Odrzucona.**

Powody:
- zależy od tego, czy operator przechwycił wyjście,
- miesza diagnostykę z audytem,
- nie daje gwarancji append-only,
- łatwo go utracić lub nadpisać.

### Alternatywa D — Zewnętrzny system logowania

**Odrzucona na start.**

Powody:
- łamie local-first discipline,
- zwiększa zależności środowiskowe,
- komplikuje odtwarzalne uruchomienia eksperymentalne.

---

## Konsekwencje

### Pozytywne

- kompletna historia zmian stanu jest dostępna dla każdego uruchomienia,
- możliwe jest porównywanie uruchomień po `run_id`,
- możliwa jest analiza śledcza po awarii lub naruszeniu kontraktu,
- log może służyć jako niezależny punkt odniesienia wobec SQLite,
- łatwiejsza jest detekcja side effects i niejawnych zmian stanu.

### Negatywne / zaakceptowane kompromisy

- log rośnie bez ograniczenia, jeśli nie ma polityki archiwizacji,
- append I/O staje się kosztem obowiązkowym dla state-changing runs,
- między zapisem do SQLite a dołączeniem do logu może istnieć krótkie okno rozbieżności,
- konieczna jest osobna walidacja integralności logu.

### Odroczone

- `itdlab audit replay` jako pełne narzędzie odtwarzania,
- polityka kompaktowania / archiwizacji logu,
- podpisywanie logu i wykrywanie manipulacji kryptograficznie.

---

## Implikacje implementacyjne

1. Adapter logu zdarzeń musi być append-only i blokować modyfikację istniejących wpisów.
2. Każde zdarzenie musi zawierać `run_id` i obowiązkowy zestaw pól.
3. Log zdarzeń musi być uwzględniany w evidence pack i weryfikacji kompletności.
4. Materialization checks powinny porównywać event log z persisted state w SQLite.
5. Naruszenie integralności logu musi wpływać na trust/evidence status przebiegu.
6. Dla state-changing runs brak event logu oznacza brak pełnego, wiarygodnego wyniku.

---

## Review triggers

ADR powinien zostać zrewidowany, jeśli:
- pojawi się potrzeba wielostrumieniowego lub rozproszonego logowania,
- event log stanie się istotnym bottleneckiem wydajnościowym,
- zostanie wdrożony podpis lub zewnętrzny system archiwizacji,
- model evidence pack przestanie traktować log jako obowiązkowy element audytu.

---

## Internal references
- `docs/EXECUTION_ASSURANCE_PROGRAM.md`
- `docs/EXECUTION_CONTRACT.md`
- `docs/EVIDENCE_MODEL.md`
- `docs/ADR/ADR-001-sqlite-as-source-of-truth.md`
- `docs/CONTEXT_VOCABULARY.md`

## Authority anchors
- `docs/REFERENCES.md` — software lifecycle and verification references
- `docs/REFERENCES.md` — requirements engineering and testing references

## Review metadata
- Owner: experimental-repository maintainer
- Status: accepted
- Last reviewed: 2026-03-30
