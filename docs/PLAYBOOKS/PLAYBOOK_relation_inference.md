# Playbook: Relation Inference

## Cel

Definiuje strategię wnioskowania relacji semantycznych między artefaktami dokumentacji IT.

---

## Typy relacji

| Typ | Znaczenie |
|-----|-----------|
| `depends_on` | Dokument A wymaga istnienia i aktualności Dokumentu B |
| `implements` | Dokument A nadaje konkretną formę polityce zdefiniowanej w B |
| `references` | Dokument A cytuje Dokument B |
| `supersedes` | Dokument A zastępuje Dokument B |
| `complements` | Dokumenty A i B razem obejmują wspólne zagadnienie |
| `derived_from` | Dokument A jest pochodną lub specjalizacją B |

---

## Wnioskowanie oparte na regułach

Relacje są wnioskowane przez reguły przechowywane w tabeli `relation_rules`. Reguła posiada:
- `from_class`: klasa dokumentu źródłowego
- `to_class`: klasa dokumentu docelowego
- `rel_type`: typ relacji do wywnioskowania
- `description`: wyjaśnienie czytelne dla człowieka

### Przykładowa reguła

> Dokument klasy `procedure` posiada relację `depends_on` do dokumentu klasy `policy` w tej samej domenie.

### Zastosowanie reguł

1. Dla każdej reguły zapytaj wszystkie pary dokumentów pasujące do `(from_class, to_class)`
2. Dla każdej pary utwórz wiersz `Relation` z `source = 'rule_engine'`
3. Ustaw pole `explanation` na opis reguły + dopasowane nazwy dokumentów
4. Dołącz zdarzenie do dziennika JSONL

---

## Wymóg wyjaśnialności

Każda wywnioskowana relacja musi posiadać niepuste pole `explanation`. Jest to egzekwowane przez Layer 18 (testy wyjaśnialności) i Gate 4.

Wyjaśnienie musi być czytelne dla człowieka i zawierać:
- Regułę, która się uruchomiła
- Nazwy dokumentów źródłowego i docelowego
- Wynik pewności (confidence score), jeśli jest poniżej 1.0

---

## Kiedy ponownie uruchomić wnioskowanie relacji

- Po zingestowaniu i sklasyfikowaniu nowych dokumentów
- Po modyfikacji reguł relacji
- Po aktualizacji kanonicznych identyfikatorów (najpierw uruchom normalizację ponownie)
