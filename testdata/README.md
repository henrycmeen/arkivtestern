# testdata

Denne mappen inneholder nå [OpenPreserve Format Corpus](https://github.com/openpreserve/format-corpus), et omfattende sett med testfiler for digital bevaringsformater.

## Hva er Format Corpus?
Format Corpus er en åpen samling av tusenvis av filer i ulike formater (tekst, bilde, lyd, video, dokument, arkiv, mm.).
Den brukes av digital bevaringsmiljøet for å teste og validere filidentifikasjon, konvertering og bevaringsverktøy.

## Bruk med arkivtestern
Ved å kjøre:

```sh
arkivtestern scan testdata/format-corpus
```

vil du teste Siegfried-integrasjonen mot et bredt spekter av ekte og representative filer. Dette gir en realistisk sjekk på hvor godt filidentifikasjon fungerer, og hvor mange filer som blir korrekt identifisert, ukjent eller feilet.

## Tips
- Du kan fortsatt legge inn egne testfiler eller undermapper i `testdata/` for egne behov.
- Resultatene fra skanningen kan brukes til å sammenligne mot andre verktøy eller til å feilsøke spesifikke formater.

---

**Kilde:** https://github.com/openpreserve/format-corpus
