# Playbook: Export to Repo 1

## Cel

Definiuje strategię promowania stabilnych, zwalidowanych metadanych semantycznych z tego laboratorium do repozytorium referencyjnego (`IT-Dokumentacja`).

---

## Co jest eksportowane?

Promowane są wyłącznie **stabilne, zwalidowane przez bramy** metadane. Nigdy surowy ani pośredni stan.

Eksportowane artefakty:
- Kanoniczne identyfikatory i klasy dokumentów
- Mapowania ról sekcji (wyłącznie stabilne)
- Graf relacji (wyłącznie zwalidowane krawędzie)
- Raport pokrycia organów

---

## Wymagania przed eksportem

Wszystkie poniższe warunki muszą być spełnione przed eksportem:

1. Gate 1–5 przeszły pomyślnie dla bieżącego uruchomienia
2. Wszystkie relacje posiadają niepuste wyjaśnienia (Layer 18)
3. Pakiet dowodów jest kompletny (Layer 27)
4. Brak nierozwiązanych kolizji kanonicznych identyfikatorów (Gate 3)

Jeśli którykolwiek warunek nie jest spełniony, `itdlab export repo1` kończy działanie z kodem `2`.

---

## Proces eksportu

```
itdlab export repo1 --target ../IT-Dokumentacja/
```

Polecenie:
1. Waliduje wszystkie bramy jakości
2. Generuje artefakty eksportu w `normalized/`
3. Kopiuje stabilne pliki do docelowego repozytorium
4. Zapisuje `export_manifest.json` z listą wszystkich promowanych plików
5. Dołącza zdarzenia `exported` do dziennika JSONL

---

## Co NIE jest eksportowane

- Surowe pliki źródłowe
- Pośredni stan normalizacji
- Eksperymentalne lub mało pewne relacje (confidence < 0.8)
- Baza danych SQLite (pozostaje wyłącznie w laboratorium)
- Dziennik zdarzeń (pozostaje wyłącznie w laboratorium)

---

## Wycofanie zmian

Eksport nie jest destrukcyjny dla laboratorium. Aby cofnąć promocję do repo 1:
- Użyj `git revert` w repozytorium referencyjnym
- Stan laboratorium pozostaje nienaruszony
