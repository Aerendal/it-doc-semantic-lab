MORFOLOGIA JEZYKA POLSKIEGO
Taksonomia regul i siec algorytmow ekstrakcji

Specyfikacja techniczna v1.0




# 1. Ile regul ma jezyk polski?
To pytanie nie ma jednej liczby jako odpowiedzi - wszystko zalazy od granularnosci z jaka definiujemy 'regule'. Najrozsadniejszy podzial wyodrebnia piecy niezalezne warstwy, z ktorych kazda ma swoj wlasny aparat formalny.



## 1.1  Warstwa L1 - Fleksja nominalna (deklinacja)
Jezyk polski posiada 7 przypadkow gramatycznych: mianownik, dopelniacz, celownik, biernik, narzednik, miejscownik i wolacz. Kazdy rzeczownik, przymiotnik, zaimek i liczebnik odmienia sie przez te przypadki w dwoch liczbach (pojedyncza i mnoga), co daje do 14 form dla jednego paradygmatu. Fleksja nominalna jest najrozleglejszym podsystemem jezyka.


### Uwaga: przyslowki jako klasa graniczna L1/L2
Przyslowki sa czescia mowy nieodmiennej przez przypadki ani liczby, co wyklucza je z klasycznej fleksji nominalnej. Jednak posiadaja wlasny system form - stopniowanie - i dlatego wymagaja wlasnego paradygmatu w schemacie. Algorytm traktuje je jako klase L1* (warstwa pomiedzy L1 i L2): nie maja paradygmatu deklinacyjnego, ale maja 3 formy stopnia. Wyzwalaja one takze alternacje w rdzeniu (szybko -> szybciej; dobrze -> lepiej - supletyw), co czyni je istotnym przypadkiem dla modulu MOD-4.

### Rodzaje gramatyczne
Jezyk polski wyroznia trzy podstawowe rodzaje, przy czym rodzaj meski dzieli sie na trzy podtypy w zaleznosci od kategorii animatycznosci i osobowosci - co ma bezposredni wplyw na forme biernika i mianownika liczby mnogiej:
Meski osobowy (mezkoosobowy) - dotyczy mezczyzn i grup mieszanych
Meski zywotny - dotyczy zwierzat i niektorych rzeczy (np. 'mam samochoda')
Meski niezywotny - dotyczy przedmiotow (np. 'mam samochod')
Zenski - rzeczowniki rodzaju zenskiego
Nijaki - rzeczowniki rodzaju nijakiego
W liczbie mnogiej rodzaje meskie zywotny i niezywotny oraz zenski i nijaki zlewaja sie w jeden rodzaj niemezkoosabowy, co oznacza, ze system rodajowy w pluralis jest 2-elementowy (mezkoosobowy vs. reszta).

## 1.2  Warstwa L2 - Fleksja werbalna (koniugacja)
Czasownik polski jest pod wzgledem morfologicznym najbardziej rozbudowana czescia mowy. Bez imieslow jedne czasownik moze przybierac do 47 roznych form; z imiesloowami liczba ta przekracza 100.


## 1.3  Warstwa L3 - Alternacje morfonologiczne
Alternacje to wymiany glosek zachodzace na granicach morfemow - przede wszystkim na granicy rdzenia i konczowki. Sa to reguly zrodlowo fonetyczne (wywodza sie z historycznych zmian dzwiekowych), ale w dzisiejszym jezyku maja charakter morfologicznie uwarunkowany: ta sama konczowka wyzwala alternacje w jednej klasie slow, a nie wyzwala jej w innej.

### Alternacje samogloskowe

### Alternacje spolgoskowe - palatalizacja historyczna
Palatalizacja to zmiekczenie twardej spogloski pod wplywem nastepujacej po niej samogloski przedniej lub gloski /j/. W jezyku polskim wykrystalizowaly sie z niej trwale alternacje morfologiczne. Ponizej lista alternacji produktywnych:

### Alternacje spolgoskowe - upodobnienia (dzwoncznosc)
Upodobnienia nie sa alternacjami morfologicznymi sensu stricto, ale musza byc uwzgledniowane przy rekonstrukcji rdzenia ze slowa z upodobnieniem wstecznym lub ubezdwieznieniem wyglosowym:
Ubezdwieznienie wyglosowe: chleb [b] -> chle[p], sad [d] -> sa[t]
Upodobnienie wsteczne (regresywne): nozka -> [no[sk]a], babka -> [ba[pk]a]
Ubezdwieznienie przed bezdwiezna: trawka -> tra[fk]a

## 1.4  Warstwa L4 - Slowotworstwo
Slowotworstwo jest najrozleglejsza warstwa regul i obejmuje sposob tworzenia nowych wyrazow z istniejacych. Dla celow algorytmicznych mozna je podzielic na cztery mechanizmy:


## 1.5  Warstwa L5 - Kongruencja i aspekt
Kongruencja (zgoda gramatyczna) to reguly nakazujace dopasowanie form gramatycznych miedzy slowami w zdaniu. Chocia z perspektywy analizy pojedynczego slowa kongruencja objawia sie tylko w jego wlasnych formach, algorytm musi wiedziec, jakie cechy slowo eksportuje do kongruencji i na co reaguje.

# 2. Siec algorytmow - architektura ogolna
Rozbicie jednego slowa na wszystkie stosujace sie do niego reguly to problem nietrywialny, bo ten sam ciag znakow moze byc roznymi formami roznych wyrazow (ambiguacja morfolgoiczna). Przykladowo 'lotu' to dopelniacz l. poj. rzeczownika 'lot' ale take celownik liczby mnogiej (rzadko uzywany). Algorytm musi wiec operowac na rozkladach prawdopodobienstwa, a nie na jednej deterministycznej odpowiedzi.
Proponowana architektura to siec pieciu wyspecjalizowanych modulow, polaczonych potokowo. Kazdy modul produkuje output uzywany przez nastepny:


## 2.1  Zasada dzialania sieci
Siec pracuje w trybie 'best-k hipotez': zamiast wybierac jedyna poprawna analize, kazdy modul utrzymuje liste k-najlepszych hipotez z wagami (prawdopodobienstwami). Wynikowy zbior regul dla slowa to suma zbiorow regul ze wszystkich hipotez, z wagami. Pozwala to na prace w trybie probabilistycznym bez korpusu kontekstowego.
Dla slow rozpoznanych jako jednoznaczne (hapax legomena lub slowa o stalym paradygmacie jak spojniki, przyimki) siec moze skrocic sciezke do 2-3 modulow. Dla slow bardzo nieregularnych (np. 'byc', 'isc') modul MOD-3 siega do slownika wyjatkow przed uruchomieniem regul ogolnych.

## 2.2  Model probabilistyczny - definicja funkcji oceniania
Ponizej zdefiniowano formalnie aparat matematyczny, na ktorym opieraja sie 'wagi' i 'prawdopodobienstwa' w calym potoku. Zamiast trenowac zewnetrzny model ML, system korzysta z log-liniowego modelu Bayesowskiego, gdzie parametry sa estymowane z korpusu NKJP lub inicjalizowane priorami z produktywnosci morfemow (kolumna productivity w tabelach morphemes i wf_patterns).
### 2.2.1  Definicja score() w MOD-1
Funkcja score(prefix, suffix) uzywana w pseudokodzie segmentacji to iloczyn trzech skladnikow:
score(p, s) = P_morph(p) * P_morph(s) * P_cooc(p, s)

P_morph(m)   = productivity[m] z tabeli morphemes / wf_patterns
= (liczba slow w NKJP zawierajacych morfem m)
/ (laczna liczba analizowanych slow)

P_cooc(p,s)  = P(p i s wspolwystepuja w tym samym wyrazie)
= count(p+s w NKJP) / count(p w NKJP)
= score_prior z tabeli wf_patterns (jesli para znana)
= 0.01 (prior dla nieznanych par - regularyzacja Laplace'a)
Wynik jest logarytmowany, aby uniknac underflow przy dlugich slowach (dlugie lancuchy morfemow daja bardzo male iloczyny prawdopodobienstw):
log_score(p, s) = log P_morph(p) + log P_morph(s) + log P_cooc(p, s)

### 2.2.2  Propagacja hipotez miedzy modulami (best-k beam search)
Na wyjsciu kazdego modulu utrzymywana jest posortowana lista co najwyzej k hipotez. Domyslnie k=5, co oznacza ze w kazdym momencie pipeline przetwarza do 5 alternatywnych analiz rownolegle. Dla kazdej hipotezy h_i przechowywana jest jej skumulowana log-waga:
W(h_i) = sum_j [ log P_j(h_i | evidence_j) ]

gdzie j przebiega po modulach 1..N, a evidence_j to
obserwacje dostepne dla modulu j (konczowki, rdzenie, alternacje).

Przy wyjsciu z MOD-6 hipotezy sa normalizowane do rozkladu:
P(h_i) = exp(W(h_i)) / sum_j exp(W(h_j))   [softmax]
Wynikowy wektor regul zawiera reguly wszystkich hipotez, ale kazda regula jest opatrzona waga P(h_i) hipotezy, z ktorej pochodzi. Reguly z waga < 0.05 (prog odciecia) sa opcjonalnie pomijane w wyjsciu.

### 2.2.3  Zrodlo danych dla estymacji prawdopodobienstw

# 3. Opis modulow
## Modul MOD-1: Normalizacja i segmentacja

### Zadania modulu
Normalizacja Unicode: rozroznienie i/y, liter diakrytycznych (a/a-ogonkowe, c/c-miekkie itp.), wielkich i malych liter.
Detekcja granic morfemow - Metoda inkrementalna: algorytm przeglada wszystkie mozliwe prefiks-rdzen-sufiks podzialy slowa i punktuje je na podstawie listy znanych morfemow. Preferowane sa podzialy, ktore daja rozpoznane morfemy z obu stron granicy.
Budowanie grafu morfemowego: slowo jest reprezentowane jako DAG (acykliczny skierowany graf), gdzie krawedzie to mozliwe granice morfemow.

### Pseudokod MOD-1
def segment(word: str) -> MorphGraph:
candidates = []
for i in range(1, len(word)):
prefix = word[:i]
suffix = word[i:]
if is_known_morpheme(prefix) and is_known_morpheme(suffix):
candidates.append((i, score(prefix, suffix)))
return build_dag(word, sorted(candidates, key=score, reverse=True))

## Modul MOD-2: Klasyfikator POS i lematyzator

### Klasyfikacja POS - algorytm
Klasyfikacja czesci mowy dla izolowanego slowa (bez kontekstu) opiera sie na sufiks-drzewie: algorytm szuka najdluzszego sufiksu slowa, ktory jest jednoznacznie lub probabilistycznie zwiazany z klasa POS. Priorytety:
Sprawdzenie w slowniku wyjatkow (nieregularne formy jak 'byl', 'beda', 'mnie').
Sprawdzenie w slowniku lemmatow (forma bazowa).
Lookup w drzewie sufiksow: sufiksy o dlugosci 1-6 znakow plus ich distribucja POS-owa z korpusu.
Fallback: reguly heurystyczne (np. slowo konczace sie na '-owanie' -> rzeczownik odsloadownikowy nijaki).

### Lematyzacja - reguly
Lematyzacja to odtworzenie formy podstawowej. Dla kazdej klasy POS reguly sa inne:

## Modul MOD-3: Dopasowanie paradygmatu fleksyjnego

### Paradygmaty deklinacyjne - hierarchia
Paradygmaty sa zorganizowane hierarchicznie: od klasy ogolnej do podklasy szczegolowej. Kazdy paradygmat to macierz [przypadek x liczba] -> konczowka, z zaznaczeniem mozliwych alternacji:

### Paradygmaty koniugacyjne - klasy
Tradycyjna gramatyka wyroznia 4 klasy koniugacyjne (lub ich warianty), ale dla celow obliczeniowych bardziej uzyteczny jest podzial na osnowy (stem classes):

## Modul MOD-4: Detekcja alternacji morfonologicznych

### Algorytm detekcji
Detektor porownuje lemat z jego rekonstruowanymi formami fleksyjnymi i identyfikuje miejsca, w ktorych rdzen ulega zmianie. Kluczowe kroki:
Wygeneruj wszystkie 14 (lub wiecej) form slowa za pomoca paradygmatu z MOD-3.
Wyrownaj pary form (np. lemat vs. forma w D.l.poj.) za pomoca algorytmu LCS (Longest Common Subsequence) dla znalezienia najdluzszego wspolnego podrdzenia.
Stwierdz, jakie gloski ulegly zmianie i w jakim kontekscie (nastepujaca konczowka, poprzedzajacy rdzen).
Dopasuj zidentyfikowana zmiane do katalogu znanych alternacji (tabela regul).
Zbuduj wektor alternacji: {lokalizacja w rdzeniu, typ alternacji, warunek triggerujacy}.

### Format wyjscia MOD-4 - przyklad
Slowo: 'rzece'
Lemat: 'rzeka'
Alternacje: [
{ typ: 'k->c', pozycja: -2 (od konca rdzenia), trigger: 'konczowka -e' },
{ typ: 'a->e', pozycja: -1 (samogloska przed alternacja k), trigger: 'msc.l.poj.' }
]

## Modul MOD-5: Analiza slowotworcza

### Algorytm rekurencyjny
Analiza slowotworcza jest z natury rekurencyjna: wyraz pochodny sam moze byc podstawa dalszej derywacji. Algorytm pracuje w petli, az rdzen nie moze byc dalej rozkladany lub osiagnie rdzen niemotywowany (wyraz prosty).
def analyze_wf(lemat: str, depth=0) -> DerivTree:
if is_root_form(lemat) or depth > MAX_DEPTH:
return DerivNode(lemat, typ='rdzen')
best = None
for prefix in KNOWN_PREFIXES:
if lemat.startswith(prefix):
base = lemat[len(prefix):]
if is_valid_base(base):
candidate = DerivNode(lemat, morph=prefix, typ='pre',
child=analyze_wf(base, depth+1))
best = max_score(best, candidate)
for suffix in KNOWN_SUFFIXES:
if lemat.endswith(suffix):
base = lemat[:-len(suffix)]
base_restored = restore_alternations(base, suffix)
if is_valid_base(base_restored):
candidate = DerivNode(lemat, morph=suffix, typ='suf',
child=analyze_wf(base_restored, depth+1))
best = max_score(best, candidate)
return best or DerivNode(lemat, typ='nieprzezroczysty')

### Typy sufiksow - wybor kluczowych klas

## Modul MOD-6: Kongruencja i aspekt

### Cechy kongruencyjne
Kazde slowo eksportuje pewne cechy gramatyczne (te, z ktorymi musza sie zgadzac inne slowa) i importuje inne (te, ktore samo musi dostosowac). Rzeczownik eksportuje rod, liczbe i przypadek do przymiotnikow; czasownik importuje osobe i liczbe od podmiotu.

### Aspekt czasownika
Algorytm aspektowy dziala w dwoch krokach: najpierw klasyfikuje czasownik jako dokonany lub niedokonany na podstawie prefiksow i obecnosci w slownikowej bazie par aspektowych; potem lokalizuje odpowiednika aspektowego (jezeli jest znany). Jesli baza nie zawiera slowa, stosowane sa heurystyki prefiksalne (przedrostek na-/za-/po-/u-/wy- bez zmiany znaczenia leksykalnego -> czesto tworzy dokonany od niedokonanego).

# 4. Format wyjscia - wektor regul
Wyjsciem calego potoku jest jeden obiekt JSON nazywany 'wektorem regul', ktory agreguje wszystkie ustalenia wszystkich modulow dla danego slowa. Ponizej schemat:

{
'slowo': 'pisalem',
'lemat': 'pisac',
'pos': 'czasownik',
'pewnosc': 0.97,
'fleksja': {
'klasa_koniugacyjna': 'K1',
'czas': 'przeszly',
'tryb': 'oznajmujacy',
'osoba': '1',
'liczba': 'pojedyncza',
'rodzaj': 'meski',
'aspekt': 'niedokonany',
'para_aspektowa': 'napisac'
},
'alternacje': [
{ 'typ': 's->sz', 'pozycja_w_rdzeniu': -2, 'trigger': 'czas_przeszly_meski' }
],
'slowotworstwo': {
'rdzen_etymologiczny': 'pis-',
'morfemy': [{ 'morph': '-a-', 'typ': 'temat' }, { 'morph': 'l', 'typ': 'suf_czasu_przeszlego' },
{ 'morph': 'em', 'typ': 'konczowka_osobowa_meski_1sg' }]
},
'kongruencja': {
'eksportuje': [],
'importuje': ['podmiot.osoba=1', 'podmiot.liczba=sg', 'podmiot.rodzaj=m']
},
'aktywne_reguly': [
'KONJ_K1_PRZESZLY_MESKI_1SG',
'ALT_NOSOWE_PRZED_L',
'ASPEKT_NIEDOKONANY',
'KONGR_CZASOWNIK_IMPORTUJE_PODMIOT'
]
}


# 5. Baza regul - schemat danych
Wszystkie reguly sa przechowywane w relacyjnej bazie SQLite, co umozliwia ich latwa inspekcje, aktualizacje i wersjonowanie. Ponizej minimalny schemat:

-- Tabela glowna regul
CREATE TABLE rules (
rule_id    TEXT PRIMARY KEY,   -- np. 'KONJ_K1_PRZESZLY_MESKI_1SG'
layer      TEXT NOT NULL,       -- L1/L2/L3/L4/L5
category   TEXT NOT NULL,       -- np. 'koniugacja', 'alternacja'
pos        TEXT,                -- czesc mowy ktorej dotyczy
description TEXT,               -- opis slowny
formal     TEXT,                -- formalny zapis reguly
coverage   INTEGER DEFAULT 0,   -- liczba slow w korpusie objeta regula
is_active  BOOLEAN DEFAULT 1
);

-- Paradygmaty fleksyjne
CREATE TABLE paradigms (
paradigm_id TEXT PRIMARY KEY,
pos         TEXT,
gender      TEXT,
ending_m_sg TEXT,   -- konczowka mianownik l. poj.
ending_d_sg TEXT,   -- konczowka dopelniacz l. poj.
-- ... (14 pol dla 7 przypadkow x 2 liczby)
requires_alt TEXT   -- FK do alternacji wymaganych
);

-- Morfemy slowotworeze
CREATE TABLE morphemes (
morpheme    TEXT PRIMARY KEY,
typ         TEXT,   -- 'prefix', 'suffix', 'infix'
pos_in      TEXT,   -- wymagana POS podstawy
pos_out     TEXT,   -- POS wyniku
semantic    TEXT,   -- kategoria semantyczna
productivity FLOAT  -- szacowana produktywnosc [0-1]
);

-- Alternacje morfonologiczne (uzywane przez MOD-4)
CREATE TABLE alternations (
alt_id      TEXT PRIMARY KEY,  -- np. 'k_c_msc_sg'
from_seg    TEXT NOT NULL,      -- segment wyjsciowy, np. 'k'
to_seg      TEXT NOT NULL,      -- segment docelowy, np. 'c'
trigger     TEXT NOT NULL,      -- kontekst wyzwalajacy, np. 'konczowka=-e'
direction   TEXT DEFAULT 'right', -- 'right' (nastepujaca) / 'left' (poprzedzajaca)
pos_scope   TEXT,               -- ograniczenie POS; NULL = wszystkie
layer       TEXT DEFAULT 'L3',
examples    TEXT,               -- przyklady slow (JSON array)
UNIQUE(from_seg, to_seg, trigger)
);

-- Wzorce slowotworcze (uzywane przez MOD-5)
CREATE TABLE wf_patterns (
pattern_id    TEXT PRIMARY KEY,  -- np. 'VERB_NACT_ANIE'
morpheme_id   TEXT NOT NULL REFERENCES morphemes(morpheme),
pos_in        TEXT NOT NULL,     -- POS podstawy wymagany
pos_out       TEXT NOT NULL,     -- POS wyniku
semantic_role TEXT,              -- 'nomen_agentis', 'nomen_actionis', 'deminutivum'...
stem_change   TEXT,              -- FK do alternations wymaganych przy dolaczaniu
productivity  FLOAT,             -- [0-1], estymowana z korpusu NKJP
score_prior   FLOAT DEFAULT 0.5, -- prior bayesowski uzywany przez MOD-1/MOD-2
examples      TEXT               -- przyklady (JSON array)
);

-- Slownik wyjatkow i form nieregularnych (uzywany przez MOD-2 i MOD-3)
CREATE TABLE exceptions (
form        TEXT NOT NULL,       -- forma nieregularna, np. 'byl', 'mnie', 'sa'
lemat       TEXT NOT NULL,       -- lemat formy, np. 'byc', 'ja', 'byc'
pos         TEXT NOT NULL,       -- czesc mowy
case_       TEXT,                -- przypadek (dla form deklinacyjnych)
number_     TEXT,                -- liczba: 'sg' / 'pl'
gender      TEXT,                -- rodzaj
person      TEXT,                -- osoba (dla czasownikow)
tense       TEXT,                -- czas
priority    INTEGER DEFAULT 100, -- im wyzszy, tym silniejszy override regul
source      TEXT DEFAULT 'sgjp', -- zrodlo: 'sgjp' / 'polimorf' / 'manual'
PRIMARY KEY (form, lemat, pos)
);


# 6. Podsumowanie: kompletna mapa regul


## Dalsze kroki
Implementacja MOD-1 i MOD-2 jako deterministycznych automatow skonczonych (FST - Finite State Transducer) - to standardowy paradygmat w NLP dla morfologii.
Populacja bazy SQLite pelnym katalogiem morfemow i paradygmatow (zrodla: SGJP - Slownik gramatyczny jezyka polskiego, PoliMorf).
Testy regresyjne: dla kazdej reguly - co najmniej 5 przykladow slow, ktore ja aktywuja i 5 slow, ktore jej nie aktywuja.
Benchmarkowanie na korpusie NKJP (Narodowy Korpus Jezyka Polskiego).



Dokument wygenerowany automatycznie na podstawie zrodel lingwistycznych | v1.0