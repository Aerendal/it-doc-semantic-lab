# Playbook: Normalization

## Cel

Definiuje strategię normalizacji nazw dokumentów i przypisywania kanonicznych identyfikatorów.

---

## Zasady normalizacji

### 1. Generowanie kanonicznego identyfikatora

Kanoniczny identyfikator jest wyprowadzany z surowej nazwy dokumentu przez:
1. Zamianę wszystkich znaków na małe litery
2. Zastąpienie spacji, myślników i ukośników podkreślnikami
3. Usunięcie znaków niealfanumerycznych (z wyjątkiem podkreślników)
4. Zwinięcie wielu podkreślników w jeden

Przykłady:
- `"Risk Register"` → `risk_register`
- `"IT Security Policy v2"` → `it_security_policy_v2`
- `"Change Management Process (CMP)"` → `change_management_process_cmp`

### 2. Obsługa kolizji

Jeśli dwa różne dokumenty generują ten sam kanoniczny identyfikator:
1. Wpis o kolizji jest zapisywany do SQLite
2. Uruchomienie NIE jest blokowane — normalizacja jest kontynuowana
3. Gate 3 blokuje promocję do czasu rozwiązania wszystkich kolizji

Opcje rozwiązania:
- Zmień nazwę jednego z dokumentów
- Dodaj odróżniający sufiks (np. `_v2`, `_legacy`)
- Połącz dokumenty, jeśli są rzeczywiście duplikatami

### 3. Rejestracja aliasów

Dokument może mieć wiele znanych aliasów. Aliasy:
- Są przechowywane w tabeli `normalizations`
- Zawsze wskazują dokładnie jeden kanoniczny identyfikator
- Nie mogą być współdzielone między kanonicznymi identyfikatorami

---

## Podgląd przed zastosowaniem

Zawsze wykonaj podgląd przed zastosowaniem zmian:

```
itdlab normalize preview
```

Przejrzyj dane wyjściowe `normalization_report.json`. Zastosuj zmiany tylko wtedy, gdy wszystkie modyfikacje są zamierzone:

```
itdlab normalize apply
```

---

## Kiedy ponownie przeprowadzić normalizację

- Po dodaniu nowych dokumentów źródłowych
- Po zmianie nazw plików źródłowych
- Po zmianie reguł normalizacji (wymaga testu migracji — Layer 15)
