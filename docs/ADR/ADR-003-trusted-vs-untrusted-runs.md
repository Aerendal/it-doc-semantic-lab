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

Nie wszystkie uruchomienia są równoważne. Uruchomienie wykonane lokalnie, z zamockowanymi zależnościami, pominiętymi warstwami lub niekompletnym pakietem dowodów, nie powinno być traktowane tak samo jak uruchomienie wykonane w kontrolowanym kontekście z pełnym pakietem dowodów i real-path verification.

Problem nie polega na tym, że uruchomienia nieautorytatywne są bezużyteczne. Są one przydatne podczas developmentu, debugowania i szybkiej informacji zwrotnej. Problem pojawia się wtedy, gdy takie uruchomienie jest *cytowane* jako dowód bramkowy lub *używane* jako podstawa decyzji promocyjnej.

Repozytorium potrzebuje formalnego i maszynowo czytelnego rozróżnienia poziomów zaufania uruchomień.

---

## Decyzja

Każde uruchomienie jest klasyfikowane w manifeście przy użyciu pola **`trust_level`**, a nie prostego booleanu `trusted`.

Dozwolone wartości kontraktowe:
- `untrusted_local`
- `trusted_ci`
- `authoritative_verify`

### Znaczenie poziomów

#### `untrusted_local`
Uruchomienie lokalne lub eksploracyjne, które:
- może być użyteczne dla developmentu,
- może generować artefakty robocze,
- ale nie stanowi autorytatywnego dowodu dla promocji,
- i nie może zostać zapieczętowane.

#### `trusted_ci`
Uruchomienie wykonane w kontrolowanym pipeline CI lub innym zaufanym kontekście pośrednim, które:
- może być cytowane jako dowód dla części bramek,
- może wspierać G1–G3, gdy `evidence.complete = true`,
- ale nie może samo w sobie stanowić podstawy promocji do repozytorium stabilnego,
- i nie może być promotion-eligible.

#### `authoritative_verify`
Uruchomienie wykonane w kontekście autorytatywnej weryfikacji, które:
- spełnia wymagania kontekstu autorytatywnego,
- używa real-path verification tam, gdzie wymagane,
- ma kompletny pakiet dowodowy,
- może wspierać ocenę G1–G5,
- może zostać zapieczętowane,
- i może stać się `promotion-eligible`, jeśli spełni pozostałe warunki kontraktowe.

### Pochodne klasyfikacje

Na potrzeby raportowania dopuszczalne jest używanie pojęć zbiorczych:
- **trusted run** — uruchomienie, którego `trust_level` jest `trusted_ci` albo `authoritative_verify`, przy czym dowody są kompletne i nie ma naruszeń polityki mocków/skipów,
- **untrusted run** — uruchomienie o `trust_level = untrusted_local` albo takie, którego poziom zaufania został zdegradowany przez brak dowodów lub naruszenie kontraktu.

Pole `trusted: boolean` nie jest już kontraktem źródłowym manifestu. Jeżeli pojawia się w warstwie raportowania pomocniczej, musi być traktowane jako pochodna od `trust_level` i evidence completeness, a nie jako źródło prawdy.

---

## Reguły klasyfikacji

### `untrusted_local`
Powinno zostać ustawione, jeśli zachodzi którykolwiek z poniższych warunków:
- `evidence.complete = false`,
- `context_validation_result.status = rejected`,
- aktywne jest pominięcie Kategorii 3 lub 4 na warstwie istotnej dla przebiegu,
- warstwa `mock-forbidden` została wykonana z mockami,
- przebieg działa poza dozwolonym kontekstem wykonania,
- uruchomienie jest jawnie eksploracyjne lub debugowe.

### `trusted_ci`
Może zostać ustawione tylko wtedy, gdy:
- kontekst wykonania jest dozwolony,
- pakiet dowodowy jest kompletny,
- przebieg nie narusza polityki `mock-forbidden`,
- nie ma aktywnych pomijań blokujących wiarygodność dla danej bramki,
- przebieg nie rości sobie prawa do promocji.

### `authoritative_verify`
Może zostać ustawione tylko wtedy, gdy:
- kontekst autorytatywny przeszedł walidację,
- przebieg działa na file-backed SQLite i rzeczywistym logu zdarzeń,
- evidence pack jest kompletny,
- wymagane bramki zostały policzone,
- brak aktywnych pomijań Kategorii 3 lub 4,
- brak naruszeń warstw `mock-forbidden`,
- przebieg jest wykonywany jako jawna, kontrolowana weryfikacja.

---

## Konsekwencje poziomu zaufania

### Dla `untrusted_local`
- nie może być podstawą promotion decision,
- nie może zostać oznaczone jako `seal_status: sealed`,
- nie może być cytowane jako authoritative evidence,
- może być zachowane dla debugowania i analizy rozwojowej.

### Dla `trusted_ci`
- może wspierać ocenę wybranych bramek technicznych,
- może być cytowane w raportach jakościowych jako dowód pośredni,
- nie może samodzielnie odblokować promocji,
- nie może być promotion-eligible.

### Dla `authoritative_verify`
- może wspierać ocenę wszystkich wymaganych bramek,
- może zostać zapieczętowane,
- może być promotion-eligible po spełnieniu wszystkich warunków kontraktowych,
- stanowi najwyższy poziom zaufania przewidziany dla tego repozytorium.

---

## Rozważane alternatywy

### Alternatywa A — wszystkie uruchomienia są równoważne

**Odrzucona.**

Powody:
- usuwa egzekwowalność,
- pozwala operatorowi deklarować zaufanie arbitralnie,
- miesza development feedback z dowodem bramkowym.

### Alternatywa B — prosty boolean `trusted`

**Odrzucona jako kontrakt źródłowy.**

Powody:
- gubi semantykę poziomów pośrednich,
- nie rozróżnia CI od autorytatywnej weryfikacji,
- utrudnia powiązanie z promotion rules i sealing,
- jest zbyt uboga względem obecnego modelu manifestu i execution context.

### Alternatywa C — ręczna flaga operatora `--trusted`

**Odrzucona.**

Powody:
- przenosi zaufanie z modelu egzekwowanego przez narzędzie na deklarację człowieka,
- zwiększa ryzyko nadużycia lub pomyłki,
- osłabia mechaniczne rozróżnienie oficjalnych i nieoficjalnych przebiegów.

---

## Konsekwencje

### Pozytywne

- manifest uruchomienia staje się bardziej precyzyjny semantycznie,
- raporty mogą rozróżniać rodzaj zaufania, nie tylko jego obecność/nieobecność,
- decyzje promocyjne mają jednoznaczny fundament kontraktowy,
- lokalne uruchomienia pozostają użyteczne, ale nie są mylone z dowodem oficjalnym.

### Negatywne / zaakceptowane kompromisy

- wzrasta liczba warunków klasyfikacyjnych,
- raportowanie staje się nieco bardziej złożone,
- konieczne jest utrzymanie spójności między manifestem, execution contract i policy documents.

### Odroczone

- kryptograficzne podpisywanie authoritative manifests,
- bardziej rozbudowane profile zaufania per środowisko lub organizacja,
- dodatkowe poziomy zaufania, jeśli repozytorium istotnie zwiększy złożoność wykonawczą.

---

## Implikacje implementacyjne

1. `trust_level` musi być polem obowiązkowym manifestu.
2. Klasyfikacja poziomu zaufania musi być liczona przez narzędzie, nie deklarowana ręcznie przez operatora.
3. `seal_status = sealed` jest dozwolone tylko dla `authoritative_verify`.
4. Promotion guard musi odmawiać promocji dla uruchomień innych niż `authoritative_verify`.
5. Warstwa raportowania może wyprowadzać pomocniczy boolean `trusted`, ale tylko jako pochodną od `trust_level` i kompletności dowodów.

---

## Review triggers

ADR powinien zostać zrewidowany, jeśli:
- model manifestu zmieni kontrakt `trust_level`,
- pojawią się nowe oficjalne konteksty wykonania,
- promotion rules zaczną wymagać dodatkowych klas zaufania,
- podpisy kryptograficzne lub zewnętrzny verifier staną się częścią execution contract.

---

## Internal references
- `docs/EXECUTION_CONTRACT.md`
- `docs/EVIDENCE_MODEL.md`
- `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`
- `docs/POLICY_MOCKS_AND_REAL_PATHS.md`
- `docs/CONTEXT_VOCABULARY.md`
- `docs/RUN_MANIFEST_SCHEMA.md`
- `docs/ADR/ADR-001-sqlite-as-source-of-truth.md`
- `docs/ADR/ADR-002-event-log-as-audit-backbone.md`

## Authority anchors
- `docs/REFERENCES.md` — verification, validation and testing references
- `docs/REFERENCES.md` — requirements engineering and standards language references

## Review metadata
- Owner: experimental-repository maintainer
- Status: accepted
- Last reviewed: 2026-03-30
