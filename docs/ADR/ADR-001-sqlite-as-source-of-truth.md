# ADR-001: SQLite as Source of Truth

| Field | Value |
|-------|-------|
| Status | Zaakceptowany |
| Date | 2026-03-29 |
| Deciders | zespół projektowy |
| Supersedes | — |
| Superseded by | — |

---

## Kontekst

Laboratorium semantyczne potrzebuje lokalnego, trwałego i odpytywalnego magazynu stanu dla:
- rekordów przebiegów,
- znormalizowanych dokumentów,
- sekcji,
- relacji,
- wyników bramek,
- statusu zapieczętowania,
- oraz metadanych istotnych dla promocji do repozytorium stabilnego.

Repozytorium ma pozostać local-first i działać bez obowiązkowej infrastruktury zewnętrznej. Model stanu musi wspierać:
- transakcje,
- deterministyczne odczyty,
- jawną wersję schematu,
- łatwe snapshotowanie,
- integrację z append-only event logiem,
- oraz prostą inspekcję podczas debugowania, przeglądu i audytu.

Celem nie jest optymalizacja pod rozproszoną przepustowość ani wielowęzłowy runtime. Celem jest:
- poprawność,
- śledzalność,
- niski koszt operacyjny,
- oraz odtwarzalne uruchomienia na pojedynczej maszynie.

---

## Decyzja

Repozytorium przyjmuje **SQLite jako główne źródło prawdy (Source of Truth)** dla stanu wykonawczego i trwałego stanu semantycznego.

SQLite jest autorytatywnym magazynem dla:
- `runs`,
- wersji schematu,
- znormalizowanych dokumentów,
- sekcji,
- relacji,
- wyników bramek,
- statusu `sealed`,
- i metadanych używanych przy ocenie promotion eligibility.

Append-only log zdarzeń JSONL pozostaje obowiązkowy, lecz nie zastępuje SQLite jako odpytywalnego stanu bieżącego. Model jest następujący:
- SQLite przechowuje stan materializowany,
- JSONL przechowuje chronologiczny ślad audytowy,
- manifest i raporty przechowują widoki dowodowe przebiegu.

---

## Czynniki decyzyjne

1. **Local-first operation**  
   Narzędzie musi działać na pojedynczej maszynie bez serwera bazy danych.

2. **Transactional integrity**  
   Wieloetapowe zmiany stanu muszą być atomowe i odwracalne przy błędzie.

3. **Auditability**  
   Recenzent musi móc porównać stan trwały z event logiem i evidence packiem.

4. **Low operational burden**  
   Repozytorium eksperymentalne nie powinno wymagać orkiestracji usług.

5. **Schema evolution**  
   Wersje schematu i migracje muszą być jawne i testowalne.

6. **Deterministic execution support**  
   Zachowanie warstwy stanu powinno sprzyjać odtwarzalnym uruchomieniom lokalnym.

---

## Rozważane alternatywy

### Alternatywa A — PostgreSQL jako SoT

**Odrzucona na start.**

Powody:
- wymaga zewnętrznego serwera,
- zwiększa koszt operacyjny bez proporcjonalnej korzyści dla local-first workflow,
- utrudnia przenośność i szybkie uruchomienie eksperymentalnego środowiska.

### Alternatywa B — Neo4j jako główne SoT

**Odrzucona na start.**

Powody:
- jest atrakcyjna dla eksploracji grafowej, ale przedwcześnie zwiększa złożoność systemu,
- wprowadza usługozależność i cięższy setup,
- nie jest konieczna jako pierwszy autorytatywny magazyn stanu.

Neo4j może zostać rozważony później jako **projection/read-side** dla eksploracji grafu, nigdy jako pierwszy fundament wykonawczy.

### Alternatywa C — RocksDB lub inny embedded KV store

**Odrzucona na start.**

Powody:
- daje zbyt niski poziom abstrakcji względem potrzeb modelu,
- utrudnia inspekcję, query semantics i reasoning relacyjny,
- wymagałaby zbudowania dodatkowej warstwy integralności i interpretowalności.

### Alternatywa D — Wyłącznie pliki YAML/JSON/Markdown

**Odrzucona jako autorytatywny stan strukturalny.**

Powody:
- brak mocnych gwarancji transakcyjnych,
- trudniejsza kontrola spójności między bytami,
- wyższe ryzyko driftu między plikami,
- słabsza egzekwowalność niezmienników.

Pliki pozostają ważne jako nośnik dowodów i raportów, ale nie jako jedyny magazyn stanu.

---

## Konsekwencje

### Pozytywne

- brak obowiązkowej zależności od zewnętrznego serwera bazy danych,
- transakcyjny zapis i rollback,
- prosta inspekcja stanu przy użyciu standardowych narzędzi SQLite,
- łatwe snapshoty plikowe,
- naturalne miejsce dla `schema_version`, `runs`, `gate state` i promotion metadata,
- dobra zgodność z local-first audytowalnością.

### Negatywne / zaakceptowane kompromisy

- oczekiwania współbieżności muszą pozostać umiarkowane i jawnie kontrolowane,
- ciężkie traversale grafowe mogą wymagać dodatkowych projekcji lub zapytań pomocniczych,
- wzrost skali ponad workflow lokalny może w przyszłości wymagać dodatkowej warstwy odczytowej.

### Odroczone

- dodatkowy magazyn tylko-do-odczytu dla eksploracji grafu,
- polityka rozbudowanej współbieżności,
- projekcje do zewnętrznych narzędzi analitycznych.

---

## Implikacje implementacyjne

1. SQLite dla authoritative runs musi działać w trybie plikowym, nie in-memory.
2. Wszystkie wieloetapowe zmiany stanu muszą być objęte transakcją.
3. Wersja schematu musi być jawnie przechowywana i walidowana.
4. Checkpoint WAL musi być kontrolowany przed obliczeniem checksum dla authoritative evidence.
5. Materialization checks powinny porównywać stan SQLite z event logiem dla krytycznych przebiegów.
6. JSONL event log pozostaje obowiązkowy; SQLite nie zastępuje warstwy audytowej.

---

## Review triggers

ADR powinien zostać zrewidowany, jeśli zajdzie którykolwiek z poniższych warunków:
- repozytorium zacznie wymagać wysokiej współbieżności multi-writer,
- zapytania grafowe staną się dominującym głównym workloadem,
- file-backed SQLite stanie się mierzalnym bottleneckiem dla authoritative runs,
- wymagania promotion-critical przekroczą gwarancje obecnego modelu SQLite.

---

## Internal references
- `docs/EXECUTION_CONTRACT.md`
- `docs/EVIDENCE_MODEL.md`
- `docs/TESTING_STANDARD.md`
- `docs/CONTEXT_VOCABULARY.md`
- `docs/ADR/ADR-002-event-log-as-audit-backbone.md`

## Authority anchors
- `docs/REFERENCES.md` — software lifecycle, architecture description, and requirements engineering
- `docs/REFERENCES.md` — verification and validation references

## Review metadata
- Owner: experimental-repository maintainer
- Status: accepted
- Last reviewed: 2026-03-30
