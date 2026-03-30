# Execution Assurance Program

## Cel

Ten dokument definiuje program zapewnienia wykonania dla repozytorium eksperymentalnego.

Celem nie jest wyłącznie dodawanie nowych możliwości semantycznych, lecz budowa systemu wytwarzania i weryfikacji, który redukuje ryzyko:
- cicho pominiętych testów,
- nieuzasadnionych zielonych przebiegów opartych na mockach,
- niekompletnych ścieżek wykonania,
- nieudokumentowanych skrótów,
- regresji ukrytych za słabymi dowodami,
- dryfu architektonicznego między repozytorium laboratoryjnym a repozytorium stabilnym.

Program jest przeznaczony dla **audytowalnego repozytorium eksperymentalnego**.  
Repozytorium może pozostać eksperymentalne, lecz model wykonania, model dowodowy i dyscyplina testowania muszą być explicite i podlegać przeglądowi.

---

## Zakres

Program dotyczy:
- ingestion źródeł,
- normalizacji,
- klasyfikacji dokumentów,
- mapowania ról sekcji,
- wnioskowania relacji,
- wyrównania autorytetów,
- eksportu/promocji do repozytorium stabilnego,
- przepływów CLI,
- manifestów przebiegów,
- pakietów dowodowych,
- projektowania i wykonywania testów.

---

## Podstawowe zasady inżynieryjne

1. **Brak cichego pomijania**  
   Pominięty test lub pominięty krok wykonania musi być explicite, uzasadniony i dający się zaraportować.

2. **Brak ukrytych mocków dla ścieżek krytycznych**  
   Krytyczne zachowanie w czasie wykonania lub na ścieżce danych nie może być weryfikowane wyłącznie przez mocki, gdy wymagana jest weryfikacja ścieżki rzeczywistej.

3. **Dowód przed zaufaniem**  
   Zielony przebieg nie jest traktowany jako wiarygodny, jeśli nie produkuje manifestu, logów i dających się prześledzić dowodów.

4. **Determinizm przede wszystkim**  
   Te same dane wejściowe powinny dawać te same dane wyjściowe lub jasno wyjaśnioną deltę.

5. **Promocja wyłącznie po udowodnieniu**  
   Komponenty są promowane do repozytorium stabilnego wyłącznie po spełnieniu zdefiniowanych technicznych bramek dowodowych.

6. **Eksperymentowanie jest dozwolone, ale nie nieudokumentowana improwizacja**  
   Każde istotne odchylenie, założenie lub uproszczenie musi być udokumentowane.

---

# Struktura programu

Program zapewnienia wykonania jest podzielony na **4 fazy** i **10 etapów implementacji**.

## Faza A — fundamenty
Etapy 1–3 definiują zasady wykonania i kontrakt dowodowy.

## Faza B — mechanizmy obronne
Etapy 4–6 budują mechanizmy zapobiegające słabym praktykom walidacji.

## Faza C — egzekwowanie i operacje
Etapy 7–9 łączą zasady z narzędziami, dokumentacją i dowodami audytowymi.

## Faza D — wdrożenie i utwardzanie
Etap 10 integruje program z codzienną pracą w repozytorium.

---

# 10 etapów implementacji

## Etap 1 — Model bramek i polityka egzekwowania

### Cel
Zdefiniować, co jest dozwolone, zakazane i warunkowo dozwolone w wykonaniu i testowaniu repozytorium.

### Dane wejściowe
- aktualna struktura repozytorium,
- aktualny przepływ testowania,
- aktualny plan eksperymentu,
- polityka promocji do repozytorium stabilnego.

### Dane wyjściowe
- model bramek,
- polityka pomijania,
- polityka mocków/fake'ów/stubów,
- polityka generowania dowodów,
- polityka decyzji o promocji.

### Wymagane artefakty
- `docs/QUALITY_GATES.md`
- `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`
- `docs/POLICY_MOCKS_AND_REAL_PATHS.md`

### Kryteria zamknięcia
- zakazane i dozwolone wzorce są explicite,
- repozytorium może klasyfikować przebieg jako prawidłowy / nieprawidłowy / niekompletny,
- definicja ścieżki krytycznej istnieje.

---

## Etap 2 — Katalog warstw testowych i mapowanie ryzyka

### Cel
Zdefiniować pełną mapę warstw testowych i co każda z nich ma wykrywać.

### Dane wejściowe
- cele potoku semantycznego,
- znane tryby awarii,
- planowany model relacji,
- planowane przepływy ingestion / normalizacji / eksportu.

### Dane wyjściowe
- katalog warstw testowych,
- mapowanie ryzyk na testy,
- uzasadnienie dla każdej warstwy testowej.

### Wymagane artefakty
- `docs/TESTING_STANDARD.md`
- `docs/TEST_CATALOG.md`
- `docs/RISK_TO_TEST_MATRIX.md`

### Kryteria zamknięcia
- każda warstwa testowa ma nazwany cel,
- każde ryzyko krytyczne jest zmapowane na co najmniej jedną warstwę testową,
- repozytorium może wyjaśnić, dlaczego dana klasa testów istnieje.

---

## Etap 3 — Kontrakt wykonania i model dowodowy

### Cel
Zdefiniować, co liczy się jako rzeczywisty przebieg wykonania i jakie dowody muszą być generowane.

### Dane wejściowe
- model wykonania CLI,
- aktualne raporty,
- zamierzone wymagania audytowe.

### Dane wyjściowe
- schemat manifestu przebiegu,
- schemat pakietu dowodowego,
- kontrakt wykonania.

### Wymagane artefakty
- `docs/EVIDENCE_MODEL.md`
- `docs/RUN_MANIFEST_SCHEMA.md`
- `docs/EXECUTION_CONTRACT.md`

### Kryteria zamknięcia
- każdy istotny przebieg może wygenerować manifest,
- manifesty zawierają wykonane kroki, pominięte kroki, konfigurację i wynik,
- pakiety dowodowe są reprodukowalne i podlegają przeglądowi.

---

## Etap 4 — Mechanizmy anty-skip

### Cel
Zapobiegać cichemu osłabianiu zestawu testów przez nieuzasadnione pomijania lub miękkie obejścia.

### Dane wejściowe
- istniejące testy,
- użycie markerów,
- aktualne wzorce skip/xfail.

### Dane wyjściowe
- lintowanie skip/xfail,
- polityka listy dozwolonych pomijań,
- raport delta dla nowo wprowadzonych pomijań.

### Wymagane artefakty
- reguły lintera testów,
- raport audytu pomijań,
- rejestr wyjątków.

### Kryteria zamknięcia
- nowe pomijania są widoczne,
- nieuzasadnione pomijania nie przechodzą bramek,
- użycie pomijań może być audytowane historycznie.

---

## Etap 5 — Mechanizmy anty-mock

### Cel
Zapobiegać fałszywie zielonym przebiegom powodowanym przez ukryte mockowanie ścieżek krytycznych.

### Dane wejściowe
- aktualne adaptery,
- aktualny styl testów,
- definicja ścieżki krytycznej z Etapu 1.

### Dane wyjściowe
- polityka mocków per warstwa,
- klasyfikacja ścieżka rzeczywista vs symulowana,
- reguły wykrywania niedozwolonego patchowania.

### Wymagane artefakty
- `docs/POLICY_MOCKS_AND_REAL_PATHS.md`
- rejestr ścieżek krytycznych,
- raport audytu użycia mocków.

### Kryteria zamknięcia
- ścieżki krytyczne mają explicite reguły dotyczące mocków,
- niedozwolone mockowanie jest wykrywalne,
- repozytorium może odróżnić weryfikację rzeczywistą od symulowanej.

---

## Etap 6 — Deterministyczne fixtures i rzeczywiste środowisko testowe

### Cel
Tworzyć powtarzalne, audytowane fixtures i pliki golden, które zmniejszają presję na „po prostu zamockuj to".

### Dane wejściowe
- corpora źródłowe,
- znormalizowane dane wyjściowe,
- przypadki złote,
- kandydaci na relacje,
- referencje autorytetów.

### Dane wyjściowe
- deterministyczne fixtures,
- pliki golden,
- przypadki negatywne,
- przypadki z korupcją danych,
- przypadki brzegowe.

### Wymagane artefakty
- `testdata/fixtures/`
- `testdata/golden/`
- `docs/FIXTURE_POLICY.md`

### Kryteria zamknięcia
- główne przepływy mają stabilne fixtures,
- oczekiwane dane wyjściowe są wersjonowane,
- dobór fixtures jest explicite i reprodukowalny.

---

## Etap 7 — Bramki jakości CLI i CI

### Cel
Połączyć model zapewnienia z wykonywalnymi poleceniami i automatycznymi bramkami.

### Dane wejściowe
- dane wyjściowe Etapów 1–6,
- projekt CLI repozytorium,
- strategia przepływu CI.

### Dane wyjściowe
- polecenia weryfikacji,
- polecenia bramek,
- integracja z potokiem CI,
- semantyka fail/warn/pass.

### Wymagane artefakty
- `make verify`
- `make audit-run`
- definicje przepływów CI
- raport podsumowania bramek

### Kryteria zamknięcia
- repozytorium może spójnie uruchamiać sprawdzenia bramek,
- CI odzwierciedla rzeczywisty status bramek,
- awarie są możliwe do podjęcia działań i interpretowalnych wniosków.

---

## Etap 8 — Playbooks i runbooks

### Cel
Dokumentować, jak praca powinna być wykonywana, powtarzana i audytowana.

### Dane wejściowe
- rzeczywiście zaimplementowane przepływy,
- kontrakt dowodowy,
- model bramek.

### Dane wyjściowe
- playbooks,
- runbooks,
- procedury operatorów,
- logika troubleshootingu.

### Wymagane artefakty
- `docs/PLAYBOOKS/`
- `docs/RUNBOOKS/`
- `docs/TROUBLESHOOTING.md`

### Kryteria zamknięcia
- główne przepływy są udokumentowane krok po kroku,
- recenzent może odtworzyć przebieg bez wyjaśnień ustnych,
- niejednoznaczność operacyjna jest zredukowana.

---

## Etap 9 — Warstwa audytowa i pakiety dowodowe

### Cel
Zapewnić, że każdy ważny przebieg pozostawia weryfikowalny ślad.

### Dane wejściowe
- model manifestu przebiegu,
- wykonanie bramek,
- dane wyjściowe fixtures i raportów.

### Dane wyjściowe
- pakiety dowodowe,
- podsumowania przebiegów,
- sumy kontrolne/odciski palca,
- raporty audytowe.

### Wymagane artefakty
- `runs/`
- `reports/`
- `evidence/`
- ustandaryzowane pliki podsumowań audytowych

### Kryteria zamknięcia
- każdy ważny przebieg ma przejrzysty pakiet dowodowy,
- dowody są wystarczające do wyjaśnienia PASS/FAIL/WARN,
- ślady wykonania można porównywać między przebiegami.

---

## Etap 10 — Wdrożenie i utwardzanie

### Cel
Zastosować system zapewnienia w codziennej pracy i zamknąć pozostałe obejścia.

### Dane wejściowe
- wszystkie poprzednie etapy,
- rzeczywiste użycie deweloperskie,
- zaobserwowane wzorce awarii/obejść.

### Dane wyjściowe
- utwardzone bramki,
- zredukowana powierzchnia obejść,
- stabilna polityka promocji,
- zaktualizowana dyscyplina operacyjna.

### Wymagane artefakty
- lista kontrolna wdrożenia,
- backlog utwardzania,
- raport przeglądu po wdrożeniu.

### Kryteria zamknięcia
- model zapewnienia jest używany w normalnym wytwarzaniu,
- obejścia są zredukowane i widoczne,
- promocja do repozytorium stabilnego jest oparta na dowodach.

---

# 30 warstw testowych/kontrolnych

Program używa **30 warstw testowych/kontrolnych**, pogrupowanych w 6 poziomów.

## Level A — kontrakt źródłowy i wejściowy (1–5)
1. testy obecności plików  
2. testy czytelności plików  
3. testy kodowania  
4. testy struktury markdown  
5. testy schematu źródłowego

## Level B — parser i ekstrakcja (6–10)
6. testy jednostkowe parsera  
7. testy parsowania fixtures  
8. testy ekstrakcji golden  
9. testy częściowej korupcji  
10. testy determinizmu

## Level C — normalizacja i model kanoniczny (11–15)
11. testy kanonicznego ID  
12. testy wykrywania kolizji  
13. testy rozwiązywania aliasów  
14. testy typowania  
15. testy migracji

## Level D — relacje i logika semantyczna (16–20)
16. testy jednostkowe reguł relacji  
17. testy spójności relacji  
18. testy wyjaśnialności  
19. testy cykli / acykliczności  
20. testy wpływu sekcji

## Level E — interfejsy i przepływy wykonania (21–25)
21. testy kontraktu CLI  
22. testy end-to-end slice  
23. testy resume / restart  
24. testy integralności dziennika zdarzeń  
25. testy materializacji SQLite

## Level F — audytowalność i kontrola wydania (26–30)
26. testy reprodukowalności  
27. testy pakietów dowodowych  
28. testy budżetu wydajnościowego  
29. testy trybów awarii  
30. testy bramek wydania

---

# Oczekiwana dokumentacja repozytorium

Repozytorium powinno docelowo zawierać co najmniej:

- `docs/QUALITY_GATES.md`
- `docs/TESTING_STANDARD.md` *(normatywny)* — filozofia testowania, 6 poziomów, obowiązkowe zasady, polityki mock/skip/promocji
- `docs/TEST_CATALOG.md` *(operacyjny)* — 30-warstwowy katalog testów z blokującą bramką, polityką mocków i siłą dowodową per warstwa
- `docs/EVIDENCE_MODEL.md`
- `docs/EXECUTION_CONTRACT.md`
- `docs/PLAYBOOKS/`
- `docs/RUNBOOKS/`
- `docs/POLICY_SKIPS_AND_EXCEPTIONS.md`
- `docs/POLICY_MOCKS_AND_REAL_PATHS.md`
- `docs/RISK_TO_TEST_MATRIX.md`
- `docs/FIXTURE_POLICY.md`
- `docs/ADR/ADR-001-sqlite-as-source-of-truth.md`

---

# Zalecana kolejność implementacji

Zalecana kolejność prac:

1. Etap 1 — model bramek i polityka egzekwowania
2. Etap 2 — katalog warstw testowych i mapowanie ryzyka
3. Etap 3 — kontrakt wykonania i model dowodowy
4. Etap 4 — mechanizmy anty-skip
5. Etap 5 — mechanizmy anty-mock
6. Etap 6 — deterministyczne fixtures i rzeczywiste środowisko testowe
7. Etap 7 — bramki CLI i CI
8. Etap 8 — playbooks i runbooks
9. Etap 9 — warstwa audytowa i pakiety dowodowe
10. Etap 10 — wdrożenie i utwardzanie

---

# Relacja do repozytorium eksperymentalnego

Program jest przeznaczony dla **repozytorium eksperymentalnego**, w którym możliwości semantyczne są badane i implementowane w kontrolowanych warunkach.

Repozytorium stabilne pozostaje celem promocji.  
Repozytorium eksperymentalne jest miejscem, gdzie:
- projektowane są nowe mechanizmy semantyczne,
- utwardzana jest dyscyplina wykonania,
- rozwijane są modele dowodowe,
- logika weryfikacji jest walidowana przed promocją.

---

# Uwaga końcowa

Celem tego programu nie jest eliminacja eksperymentowania.
Celem jest zapewnienie, że eksperymentowanie jest:
- explicite,
- podlegające przeglądowi,
- powtarzalne,
- oparte na dowodach,
- i odporne na ciche osłabianie standardów jakości.
