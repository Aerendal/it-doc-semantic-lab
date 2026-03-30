# Quality Gates

## Cel

Niniejszy dokument definiuje **bramki jakości** używane w repozytorium eksperymentalnym.

Celem modelu bramek jest zapewnienie, że prace nad repozytorium nie są traktowane jako zakończone jedynie dlatego, że kod istnieje lub że częściowy przebieg wyprodukował powierzchownie zielony wynik.

Zmiana, przebieg, eksperyment lub kandydat do promocji musi przejść przez właściwy zestaw bramek, zanim zostanie uznany za:
- wiarygodny technicznie,
- operacyjnie reprodukowalny,
- audytowalny,
- kwalifikujący się do promocji,
- lub nadający się do zewnętrznego przeglądu.

Niniejszy dokument jest częścią programu zapewnienia wykonania opisanego w `docs/EXECUTION_ASSURANCE_PROGRAM.md`.

---

## Zależność od standardu testowania

Ocena bramek zależy od standardu testowania repozytorium i katalogu testów.

**Źródło normatywne:**
- `docs/TESTING_STANDARD.md` — filozofia testowania, 6 poziomów, obowiązkowe zasady i konwencje

**Definicje warstw operacyjnych:**
- `docs/TEST_CATALOG.md` — 30-warstwowy katalog testów z mapowaniem bramek per warstwę, polityką mocków i siłą dowodów

Bramka nie może być raportowana jako wiarygodna, jeśli wymagane warstwy testów nie zostały wykonane, zostały po cichu pominięte lub wytworzyły niewystarczające dowody.

Polityki wyjątków i mocków, które określają dopuszczalne odchylenia, są zdefiniowane w:
- `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`
- `docs/POLICY_MOCKS_AND_REAL_PATHS.md`

---

## Zasady

1. **Dowody ponad twierdzenia**  
   Twierdzenie, że dana funkcjonalność działa, musi być poparte logami, manifestami, raportami i powtarzalnym wykonaniem.

2. **Bramki uwzględniające zakres**  
   Nie każda zmiana wymaga każdej bramki. Bramki są stosowane zgodnie z typem i wpływem zmiany.

3. **Brak ukrytego łagodzenia standardów**  
   Pominięcia, mocki, obejścia, tryby degraded oraz ręczne nadpisania muszą być jawne i możliwe do przeglądu.

4. **Promocja wyłącznie po udowodnieniu**  
   Promocja do stabilnego repozytorium jest dozwolona tylko po przejściu bramki promocji.

5. **Zamknięte na błąd tam, gdzie to konieczne**  
   Dla ścieżek krytycznych brakujące dowody lub niejednoznaczny status są traktowane jako błąd, nie sukces.

---

# Poziomy bramek

Repozytorium używa pięciu poziomów bramek.

## Gate G0 — Bramka higieny repozytorium

### Cel
Zapewnienie, że repozytorium jest w stanie możliwym do przeglądu i niechaotyczenym, zanim zostanie podjęta głębsza walidacja.

### Dotyczy
- każdej gałęzi proponowanej do przeglądu,
- każdego kandydata PR,
- każdego kandydata do wydania,
- każdego kandydata do promocji.

### Minimalne kontrole
- drzewo robocze jest czyste dla ocenianego stanu,
- wymagane pliki dokumentacji istnieją,
- żadne zabronione artefakty lokalne nie są zatwierdzone,
- żadne pliki tymczasowe/wyłącznie debugowe nie są częścią zmiany,
- nazewnictwo plików i struktura repozytorium pozostają spójne.

### Wymagane dowody
- podsumowanie statusu repozytorium,
- delta inwentarza plików,
- wynik kontroli zabronionych artefaktów.

### Przykłady błędów
- zatwierdzone błędne archiwa,
- zatwierdzone tymczasowe raporty,
- pliki pomocnicze debugowania pozostawione w śledzonym drzewie,
- nieudokumentowane zmiany struktury repozytorium.

---

## Gate G1 — Bramka kontraktu statycznego

### Cel
Zapewnienie, że kontrakt repozytorium jest wewnętrznie spójny przed wykonaniem w czasie rzeczywistym.

### Dotyczy
- zmian schematów,
- zmian CLI,
- zmian konfiguracji,
- zmian modelu źródłowego,
- zmian modelu relacji,
- zmian modelu eksportu.

### Minimalne kontrole
- wymagane pliki i schematy wczytują się pomyślnie,
- pliki kontraktów YAML/JSON/Markdown parsują się poprawnie,
- nazwy pól, wymagane pola i wyliczenia są prawidłowe,
- nie wprowadzono zabronionych ani przestarzałych pól,
- sygnatury poleceń pozostają interpretowalne.

### Wymagane dowody
- raport walidacji schematu,
- raport statycznego lintera,
- podsumowanie różnicy kontraktu.

### Przykłady błędów
- nieprawidłowy YAML w plikach schematów,
- brakujące wymagane pola schematu,
- mieszane typowanie dla tego samego pola,
- zepsute definicje poleceń,
- nieudokumentowane odchylenia kontraktu.

---

## Gate G2 — Bramka integralności wykonania

### Cel
Zapewnienie, że przebieg faktycznie nastąpił w wiarygodny technicznie sposób.

### Dotyczy
- przebiegów ingestu,
- przebiegów normalizacji,
- przebiegów wnioskowania relacji,
- przebiegów eksportu,
- twierdzeń o wynikach eksperymentów.

### Minimalne kontrole
- manifest przebiegu jest wygenerowany,
- wykonane kroki są wymienione,
- pominięte kroki są jawnie wymienione,
- konfiguracja i zestaw wejść są zarejestrowane,
- logi istnieją i odpowiadają manifestowi,
- status wyników jest jawny (`PASS`, `FAIL`, `WARN`, `INCOMPLETE`).

### Wymagane dowody
- manifest przebiegu,
- logi wykonania,
- migawka konfiguracji,
- sumy kontrolne lub odciski palców krytycznych wejść,
- raport podsumowujący.

### Przykłady błędów
- brak manifestu,
- niejednoznaczny status przebiegu,
- wykonane kroki, ale niezarejestrowane,
- paczka dowodów brakuje obowiązkowych plików,
- niezgodność między manifestem a wytworzonymi wynikami.

---

## Gate G3 — Bramka wiarygodności weryfikacji

### Cel
Zapewnienie, że walidacja była znacząca i nieosłabiona przez ciche pominięcia, nieuzasadnione mocki lub dowody o niskiej wartości.

### Dotyczy
- wszystkich wycinków eksperymentów,
- wszystkich twierdzeń o funkcjach semantycznych,
- wszystkich kandydatów do promocji,
- wszystkich zewnętrznie recenzowanych przebiegów.

### Minimalne kontrole
- mające zastosowanie testy zostały wykonane,
- użycie pomijania jest jawne i uzasadnione,
- użycie xfail jest jawne i uzasadnione,
- mockowanie na ścieżkach krytycznych jest zgodne z polityką,
- golden outputs są aktualne i możliwe do przeglądu,
- oczekiwania dotyczące determinizmu są spełnione lub wyjaśnione.

### Wymagane dowody
- podsumowanie testów,
- raport audytu pomijania/xfail,
- raport ścieżki mock vs ścieżki rzeczywistej,
- raport determinizmu,
- raport weryfikacji golden.

### Przykłady błędów
- zachowanie na ścieżce krytycznej walidowane wyłącznie przez ukryte mocki,
- nowe pominięcie dodane bez rejestracji,
- zielony wynik z brakującą weryfikacją golden,
- niestabilne wyniki bez wyjaśnienia,
- testy przechodzące tylko w trybie degraded.

---

## Gate G4 — Bramka promocji

### Cel
Decydowanie, czy eksperymentalna funkcjonalność kwalifikuje się do promocji do stabilnego repozytorium.

### Dotyczy
- funkcji przeznaczonych dla `IT-Dokumentacja`,
- eksportów mających stać się kontraktem repozytorium,
- pól metadanych przeznaczonych dla stabilnych szablonów,
- mechanizmów wielokrotnego użytku przeznaczonych dla stabilnego środowiska uruchomieniowego.

### Minimalne kontrole
- G0–G3 przeszły,
- złoty standard istnieje,
- wszystkie wymagane przypadki corpus przechodzą,
- idempotentność jest potwierdzona,
- dry-run integracji jest zielony,
- walidacja stabilnego repozytorium pozostaje zielona po integracji.

### Wymagane dowody
- raport kandydata do promocji,
- raport pokrycia corpus,
- raport idempotentności,
- raport weryfikacji integracji,
- notatka o wycofaniu lub oświadczenie o możliwości wycofania.

### Przykłady błędów
- przejście wyłącznie na jednym przypadku happy-path,
- brak złotego standardu,
- niestabilne powtarzalne wykonanie,
- stabilne repozytorium zepsute po próbie scalenia,
- niejasne konsekwencje migracji.

---

# Macierz zastosowania bramek

## Typy zmian i wymagana minimalna bramka

| Typ zmiany | Minimalna bramka |
|---|---|
| Wyłącznie dokumentacyjna notatka, bez wpływu na kontrakt | G0 |
| Zmiana definicji schematu / YAML / kontraktu | G1 |
| Zmiana parsera / normalizatora / wykonania relacji | G2 |
| Zmiana wpływająca na testy lub weryfikację | G3 |
| Promocja do stabilnego repozytorium | G4 |

## Reguła eskalacji
Jeśli zmiana dotyczy więcej niż jednej kategorii, **obowiązuje najwyższa mająca zastosowanie bramka**.

Przykład:
- refaktoryzacja parsera, która zmienia również reguły schematu i oczekiwania testów, to **nie** tylko G2; musi spełniać **G3**.

---

# Model statusu bramki

Każdy wynik bramki musi używać jednego z następujących statusów:

- `PASS` — bramka spełniona,
- `WARN` — bramka spełniona z udokumentowanym problemem nieblokującym,
- `FAIL` — bramka niespełniona,
- `INCOMPLETE` — ocena nie mogła zostać zakończona z powodu brakujących wymaganych dowodów.

## Reguły interpretacji

### PASS
Używane tylko gdy wszystkie obowiązkowe kontrole dla bramki są spełnione.

### WARN
Używane tylko gdy:
- bramka została wykonana,
- dowody są kompletne,
- problem jest udokumentowany,
- problem jest jawnie sklasyfikowany jako nieblokujący.

### FAIL
Używane gdy jakiekolwiek wymaganie blokujące jest naruszone.

### INCOMPLETE
Używane gdy sama ocena nie jest wiarygodna z powodu brakujących obowiązkowych dowodów.

> `INCOMPLETE` nigdy nie może być raportowane tak, jakby było równoważne z `PASS`.

---

# Zabronione wzorce

Następujące wzorce są zabronione, chyba że zostały jawnie udokumentowane i zatwierdzone w ramach polityki wyjątków.

## Związane z pomijaniem
- ciche wprowadzanie `skip` lub `xfail`,
- szerokie warunkowe pomijanie bez jawnego powodu,
- pomijanie używane jako substytut naprawy zepsutej ścieżki.

## Związane z mockami
- mockowanie krytycznych ścieżek środowiska uruchomieniowego/danych bez jawnego zezwolenia polityki,
- zastępowanie rzeczywistej weryfikacji fałszywymi adapterami przy raportowaniu wyników jako prawdziwej walidacji,
- monkeypatching zachowania w sposób ukrywający rzeczywiste ryzyko integracji.

## Związane z dowodami
- twierdzenie o sukcesie bez manifestu i logów,
- usuwanie lub pomijanie dowodów dla twierdzonego przebiegu,
- publikowanie podsumowania PASS z brakującymi obowiązkowymi artefaktami.

## Związane z kontraktem
- wprowadzanie nowych wymaganych pól bez aktualizacji schematu,
- używanie tego samego pola z konfliktującą semantyką,
- zmiana semantyki eksportu bez dokumentowania konsekwencji migracji.

---

# Wyjątki i odchylenia

Wyjątki są dozwolone tylko gdy wszystkie poniższe warunki są prawdziwe:

1. wyjątek jest udokumentowany,
2. powód jest konkretny i techniczny,
3. zakres jest ograniczony,
4. czas trwania jest ograniczony lub wyzwalany przeglądem,
5. wynikające ryzyko jest opisane,
6. repozytorium nadal uczciwie raportuje status degraded.

## Wymagane pola rekordu wyjątku
- ID wyjątku,
- data,
- właściciel,
- dotknięta bramka,
- dotknięty komponent,
- powód,
- ryzyko,
- data wygaśnięcia/przeglądu,
- plan mitygacji.

Zalecany plik:
- `docs/EXCEPTION_REGISTER.md`

---

# Wymagane dowody per bramka

| Bramka | Minimalne dowody |
|---|---|
| G0 | podsumowanie statusu repozytorium, kontrola zabronionych artefaktów |
| G1 | raport walidacji schematu, raport statycznego lintera |
| G2 | manifest przebiegu, logi, migawka konfiguracji, podsumowanie |
| G3 | podsumowanie testów, audyt pomijania, audyt mocków, raport determinizmu |
| G4 | raport promocji, pokrycie corpus, raport idempotentności, weryfikacja integracji stabilnej |

---

# Reguła wydania i promocji

Komponent może zostać promowany z repozytorium eksperymentalnego do stabilnego tylko jeśli:

1. zamierzony stabilny zakres jest jawny,
2. G0–G4 przechodzą,
3. funkcja ma złoty standard,
4. pokrycie corpus jest wystarczające dla deklarowanego zakresu,
5. integracja do stabilnego repozytorium pozostaje zielona,
6. wycofanie jest możliwe lub jawnie przeanalizowane.

Jeśli którykolwiek z tych warunków jest niespełniony, funkcja pozostaje eksperymentalna.

---

# Mapowanie do programu zapewnienia wykonania

Niniejszy dokument implementuje **Etap 1 — Model bramek i polityka egzekwowania** z `docs/EXECUTION_ASSURANCE_PROGRAM.md`.

Dokumenty następcze oczekiwane po tym:
- `docs/TESTING_STANDARD.md`
- `docs/TEST_CATALOG.md`
- `docs/EVIDENCE_MODEL.md`
- `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`
- `docs/POLICY_MOCKS_AND_REAL_PATHS.md`

---

# Uwaga końcowa

Celem bramek jakości nie jest uniemożliwienie eksperymentowania.

Celem jest zapewnienie, że eksperymentowanie pozostaje:
- jawne,
- zdyscyplinowane,
- możliwe do przeglądu,
- poparte dowodami,
- i odporne na przypadkowe lub celowe osłabianie standardów technicznych.
