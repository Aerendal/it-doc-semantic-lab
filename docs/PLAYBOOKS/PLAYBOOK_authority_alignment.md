# Playbook: Authority Alignment

## Cel

Definiuje strategię powiązania artefaktów dokumentacji IT z regulacyjnymi i normalizacyjnymi organami władzy (authorities).

---

## Czym jest odwołanie do organu?

Odwołanie do organu (`authority_ref`) łączy dokument lub sekcję z konkretnym punktem zewnętrznego organu władzy:
- Ramy regulacyjne (np. ISO 27001, GDPR, 21 CFR Part 11)
- Standardy branżowe (np. ITIL, COBIT)
- Wewnętrzne polityki (traktowane jako wewnętrzne organy władzy)

---

## Kiedy tworzyć odwołania do organów

1. Dokument wprost cytuje regulację lub standard w swoim tekście
2. Nagłówek lub treść sekcji odpowiada konkretnemu punktowi
3. Macierz branżowa wymaga określonych dokumentów dla zgodności

---

## Proces wyrównania do organów

### 1. Zdefiniuj organy

Organy są odwoływane przez nazwę + punkt, np.:
- `authority: "ISO 27001"`, `clause: "A.12.6.1"`
- `authority: "GDPR"`, `clause: "Art. 32"`

### 2. Powiąż dokumenty

```
itdlab authority check
```

Polecenie:
1. Skanuje wszystkie zingestowane dokumenty w poszukiwaniu wzorców organów
2. Tworzy wiersze `authority_ref` dla wykrytych dopasowań
3. Raportuje pokrycie dla każdego organu w `authority_coverage_report.json`

### 3. Przejrzyj pokrycie

Pokrycie organów dla zestawu dokumentów to stosunek powiązanych punktów do łącznej liczby punktów danego organu. Cel: ≥ 80% dla każdego wymaganego organu.

---

## Tryby awaryjne

| Objaw | Działanie |
|-------|-----------|
| Niskie pokrycie wymaganego organu | Przejrzyj dokumenty pod kątem brakujących cytowań |
| Zduplikowane odwołania do organów | Sprawdź ponowne uruchomienia bez deduplikacji |
| Niezgodność nazwy organu | Ustandaryzuj nazwy organów przed powiązaniem |
