# ADR-003: Trusted vs. Untrusted Runs

| Field | Value |
|-------|-------|
| Status | Zaakceptowany |
| Date | 2026-03-30 |
| Deciders | zespół projektowy |
| Supersedes | — |
| Superseded by | — |

---

## Kontekst

Nie wszystkie uruchomienia są równoważne. Uruchomienie wykonane w środowisku deweloperskim z zamockowanymi zależnościami, pominiętymi warstwami lub niekompletnym pakietem dowodów nie powinno być traktowane tak samo jak uruchomienie wykonane w czystym środowisku z prawdziwymi danymi wejściowymi i pełnymi dowodami. Traktowanie wszystkich zielonych uruchomień jako równoważnych stwarza fałszywe poczucie jakości.

Problem nie leży w tym, że niekompletne uruchomienia są błędne — często są użyteczne podczas dewelopmentu. Problem pojawia się wtedy, gdy niekompletne uruchomienie jest *cytowane* jako dowód bramkowy lub *używane* jako podstawa decyzji o promocji.

System potrzebuje formalnego rozróżnienia między uruchomieniami, które mogą być cytowane, a tymi, które nie mogą.

---

## Decyzja

Każde uruchomienie jest klasyfikowane jako **trusted** lub **untrusted** w momencie zapisywania manifestu uruchomienia.

### Uruchomienie jest trusted wtedy i tylko wtedy, gdy:

1. Wszystkie warunki wstępne w `docs/EXECUTION_CONTRACT.md` zostały spełnione.
2. Wszystkie warunki końcowe w `docs/EXECUTION_CONTRACT.md` zostały zaspokojone (`exit_code = 0` lub `exit_code = 2`).
3. `evidence.complete = true` w `run_manifest.json`.
4. Żadna warstwa `mock-forbidden` nie została wykonana z mockami.
5. Żaden skip Kategorii 3 ani Kategorii 4 nie był aktywny (zgodnie z `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`).
6. Środowiskiem nie była lokalna maszyna deweloperska z nadpisanymi ścieżkami DB lub logu, które omijają standardową izolację (konkretnie: `--db` i `--log` muszą wskazywać na standardowe ścieżki lub jawnie zadeklarowane alternatywne ścieżki z udokumentowanym uzasadnieniem).

Jeśli którykolwiek z tych warunków jest fałszywy, w manifeście uruchomienia ustawiane jest `trusted = false`. Uruchomienie jest rejestrowane i dostępne do inspekcji, ale nie może być cytowane jako dowód do oceny bramkowej ani promocji.

### Konsekwencje klasyfikacji untrusted

- Uruchomienie nie może być przywołane w raporcie bramkowym jako cytowanie PASS.
- Uruchomienie nie może pojawić się na liście kontrolnej gotowości do promocji jako zakończony element.
- `itdlab audit evidence <run_id>` zgłosi uruchomienie jako untrusted.
- Wyniki testów uruchomienia mogą być nadal używane nieformalnie jako informacja zwrotna podczas dewelopmentu, ale nie wolno mylić tego z dowodem bramkowym.

---

## Rozważane alternatywy

### 1. Wszystkie uruchomienia są równoważne; zaufanie jest deklarowane zewnętrznie

**Podejście:** Nie klasyfikuj uruchomień; pozwól operatorowi zadeklarować, którym ufa.

**Odrzucono, ponieważ:**
- Usuwa to obiektywną egzekwowalność. Operator pod presją zadeklaruje uruchomienia untrusted jako trusted.
- Klasyfikacja musi być deterministyczna i obliczana przez narzędzie, a nie deklarowana przez operatora.

### 2. Zaufanie jest binarne: tylko PASS/FAIL

**Podejście:** Uruchomienie kończące się kodem 0 jest trusted; każdy inny kod — nie.

**Odrzucono, ponieważ:**
- Pozwala to na klasyfikację uruchomień z niekompletnymi pakietami dowodów i aktywnymi naruszeniami mock-forbidden jako trusted, o ile uruchomienie kończy się kodem 0.
- Zielony kod wyjścia jest warunkiem koniecznym, ale niewystarczającym dla zaufania.

### 3. Osobna flaga „trybu audytu"

**Podejście:** Dodaj flagę `--trusted` włączającą ściślejsze sprawdzanie.

**Odrzucono, ponieważ:**
- Tworzy to dwutorowy system, w którym większość uruchomień nigdy nie jest audytowana.
- Deweloperzy nie będą używać `--trusted` podczas normalnego dewelopmentu; nie zapewnia to ochrony przed przypadkowym nadużyciem.
- Klasyfikacja zaufania powinna być automatyczna, a nie opcjonalna.

---

## Konsekwencje

### Pozytywne

- Raporty dowodów bramkowych mogą wyraźnie odróżniać uruchomienia trusted od untrusted.
- Decyzje o promocji mają jasną, obiektywną podstawę: co najmniej jedno uruchomienie trusted na bramkę.
- Przepływ pracy dewelopmentu nie jest blokowany — uruchomienia untrusted są nadal szybkie, użyteczne i rejestrowane.
- Klasyfikacja jest obliczana deterministycznie; brak niejednoznaczności co do statusu uruchomienia.

### Negatywne / zaakceptowane kompromisy

- Pewne fałszywe pozytywy: uruchomienie na maszynie deweloperskiej z pełnymi dowodami, ale niestandardowymi ścieżkami, zostanie sklasyfikowane jako untrusted. Deweloper może to obejść, podając standardowe ścieżki.
- Rozróżnienie trusted/untrusted wprowadza krok klasyfikacji w raporcie audytowym, który musi być utrzymywany wraz z dodawaniem nowych warunków.

### Odroczone

- Profile zaufania per-organizacja (np. CI zawsze trusted niezależnie od ścieżki) są odroczone.
- Kryptograficzne podpisywanie manifestów uruchomień trusted jest odroczone.

---

## Uwagi implementacyjne

- Pole `trusted` jest obliczane i zapisywane do `run_manifest.json` przez logikę finalizacji uruchomienia.
- `itdlab audit runs` wyświetla status zaufania dla każdego uruchomienia.
- `itdlab audit evidence <run_id>` weryfikuje wszystkie warunki i raportuje, które zawiodły, jeśli `trusted = false`.
- Implementacja Go: logika klasyfikacji zaufania powinna znajdować się w `internal/app/audit/trust.go`.

---

## Odwołania wewnętrzne
- `docs/EXECUTION_CONTRACT.md`
- `docs/EVIDENCE_MODEL.md`
- `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`
- `docs/POLICY_MOCKS_AND_REAL_PATHS.md`
- `docs/CONTEXT_VOCABULARY.md`
- `docs/ADR/ADR-001-sqlite-as-source-of-truth.md`
- `docs/ADR/ADR-002-event-log-as-audit-backbone.md`

## Metadane przeglądu
- Owner: zespół projektowy
- Status: zaakceptowany
- Last reviewed: 2026-03-30
