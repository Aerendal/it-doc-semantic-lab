# Fixture Policy

Dokument ten definiuje sposób tworzenia, przeglądania, przechowywania i utrzymywania fixtures testowych w tym repozytorium.

---

## Cel

Fixtures są kontrolowanymi danymi wejściowymi, względem których uruchamiane są testy. Jakość fixtures determinuje jakość dowodów. Test działający na słabo zrozumianej fixture produkuje bezsensowne dowody — nawet jeśli sam test jest dobrze napisany.

---

## Definicje

| Termin | Definicja |
|--------|-----------|
| **fixture** | Statyczny plik lub struktura danych używana jako kontrolowane dane wejściowe testu |
| **real-file fixture** | Fixture wywodząca się z lub równoważna rzeczywistemu materiałowi źródłowemu |
| **synthetic fixture** | Fixture skonstruowana od podstaw dla konkretnego przypadku testowego |
| **golden output** | Oczekiwane wyjście dla danej fixture, przechowywane obok niej |
| **fixture version** | Identyfikator adresowany treścią dla fixture (SHA-256 zawartości) |
| **fixture inventory** | Rejestr wszystkich fixtures w `internal/testkit/fixtures/` |

---

## Lokalizacje fixtures

| Typ | Ścieżka | Przeznaczenie |
|-----|---------|---------------|
| Fixtures wejściowe | `internal/testkit/fixtures/` | Kontrolowane dane wejściowe testów |
| Golden outputs | `internal/testkit/golden/` | Oczekiwane wyjścia dla testów regresji |
| Fixtures źródłowe integracji | `sources/` | Rzeczywisty materiał źródłowy (używany również jako fixtures integracyjne) |
| Pomocnicy builderu | `internal/testkit/builders/` | Programatyczna konstrukcja fixtures dla testów jednostkowych |

---

## Kategorie

### Kategoria A — Fixtures z prawdziwych plików

Fixtures wywodzące się bezpośrednio z rzeczywistych dokumentów źródłowych w `sources/`. Są to fixtures najwyższej jakości i muszą być używane do testów integracyjnych i akceptacyjnych.

**Zasady:**
- Nie mogą być modyfikowane dla wygody testowania. Jeśli prawdziwy plik jest nieprawidłowym wejściem, należy to udokumentować w komentarzu w teście.
- Jeśli prawdziwy plik się zmienia, testy, które go używają, muszą zostać przejrzane i zaktualizowane.
- Golden outputs oparte na fixtures z prawdziwych plików muszą być regenerowane i przeglądane po każdej zmianie parsera.

### Kategoria B — Fixtures syntetyczne

Fixtures skonstruowane specjalnie dla przypadku testowego. Używane do testów jednostkowych i kontraktowych, gdzie prawdziwy plik byłby zbyt zaszumiony lub niedostępny.

**Zasady:**
- Muszą zawierać komentarz w pliku fixture lub w teście wyjaśniający reprezentowany scenariusz.
- Muszą być reprezentatywne dla prawdziwej klasy danych wejściowych, które zastępują.
- Nie mogą pomijać właściwości strukturalnie istotnych (np. nie pomijać pola `phase:`, które parser musi obsługiwać).

### Kategoria C — Fixtures budowane przez builder

Fixtures konstruowane programatycznie za pomocą `internal/testkit/builders/`. Są to struktury Go, nie pliki.

**Zasady:**
- Buildery muszą domyślnie produkować strukturalnie prawidłowe dane wejściowe.
- Buildery muszą udostępniać pomocniki mutacji do testowania przypadków brzegowych (np. `WithMissingPhase()`, `WithDuplicateName()`).
- Buildery muszą być przeglądane przy zmianie modelu dziedzinowego.

---

## Zasady tworzenia

1. **Każda nowa fixture musi mieć udokumentowane przeznaczenie.** Fixture bez wyjaśnienia, dlaczego istnieje, jest obciążeniem utrzymaniowym.

2. **Fixtures muszą testować konkretną, nazwaną właściwość.** Nazwij plik odpowiednio: `risk_register_missing_phase.md`, nie `test_fixture_42.md`.

3. **Fixtures nie mogą być tworzone przez kopiowanie aktywnych danych produkcyjnych bez przeglądu zawartości.** Jeśli zawartość jest wrażliwa, musi zostać zanonimizowana przed zatwierdzeniem.

4. **Golden outputs muszą być generowane przez deterministyczne polecenie.** Ręcznie napisane pliki golden są dopuszczalne tylko dla wersji początkowej; kolejne wersje muszą być generowane.

5. **Polecenie użyte do wygenerowania pliku golden musi być zapisane** w komentarzu na początku pliku golden lub w teście, który go używa.

---

## Zasady przeglądu

1. **Nowe fixtures muszą być przeglądane przed scaleniem.** Recenzent musi potwierdzić: przeznaczenie jest udokumentowane, zawartość jest reprezentatywna, brak wrażliwych danych.

2. **Aktualizacje plików golden muszą być przeglądane tak starannie jak zmiany kodu.** Zmieniony plik golden oznacza zmianę zachowania. Musi to być celowe i udokumentowane.

3. **Fixtures dla warstw `mock-forbidden` muszą być fixtures z prawdziwych plików (Kategoria A) lub przeglądanymi fixtures syntetycznymi (Kategoria B).** Buildery Kategorii C nie są wystarczające dla testów warstwy `mock-forbidden`.

---

## Wersjonowanie

Fixtures są wersjonowane według zawartości (SHA-256). Gdy fixture jest aktualizowana:
1. Stary SHA musi być zapisany w komentarzu lub historii zmian, jeśli zmiana wpływa na golden outputs.
2. Każdy test, który odwołuje się do fixture przez skrót zawartości, musi zostać zaktualizowany.
3. Każdy golden output wywodzący się z fixture musi zostać wygenerowany ponownie.

Nie istnieje formalny numer wersji fixture poza adresem zawartości SHA-256.

---

## Inwentarz

Inwentarz fixtures jest utrzymywany nieformalnie: każdy plik fixture musi mieć komentarz nagłówkowy wyjaśniający jego przeznaczenie. Formalny indeks (`internal/testkit/fixtures/INVENTORY.md`) jest zalecany, gdy liczba fixtures przekroczy 20.

---

## Zabronione praktyki

| Praktyka | Dlaczego jest zabroniona |
|----------|--------------------------|
| Fixtures zawsze produkujące PASS | Nie testują niczego. Fixture, która nigdy nie wyzwala trybu awarii, nie jest fixture; jest martwym kodem. |
| Anonimowe fixtures (`fixture1.md`, `testdata.json`) | Recenzent nie może określić, jaka właściwość jest testowana |
| Wygenerowane fixtures zatwierdzone bez przeglądu | Nieprzeglądana wygenerowana zawartość stanowi tylne wejście dla nieprawidłowych oczekiwanych wyjść |
| Fixtures z aktywnymi poświadczeniami lub danymi osobowymi | Naruszenie bezpieczeństwa |
| Mutowalne fixtures współdzielone między testami | Testy stają się zależne od kolejności i nieodtwarzalne |

---

## Wewnętrzne odniesienia
- `docs/TESTING_STANDARD.md`
- `docs/TEST_CATALOG.md`
- `docs/POLICY_MOCKS_AND_REAL_PATHS.md`
- `docs/EVIDENCE_MODEL.md`

## Metadane przeglądu
- Właściciel: zespół projektowy
- Status: szkic
- Last reviewed: 2026-03-30
