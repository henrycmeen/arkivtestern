# Arkivtestern

> Én AI-drevet arkivtest for oss – rask, leken og utvidbar.

## Overordnet plan for arkivtestern (v0.1 CLI først)

| Fase | Milepæl | Hovedoppgaver |
|------|---------|---------------|
| 0 | Kartlegging | Sjekk lignende verktøy (digital-preservation-toolkit, Siegfried, JHOVE, veraPDF, awesome-digital-preservation). Noter hvilke funksjoner vi kan «wrappe». |
| 1 | CLI-skjelett | `arcaive scan <path>` fungerer. Språk: Go (statisk binær for macOS/Apple Silicon, Linux & Windows). Pakker: Cobra (CLI), Viper (config), Zap (logging). |
| 2 | Fil-identifikasjon | PRONOM-ID & MIME. Integrér Siegfried (Go-lib – null Java-avhengighet) for lynrask identifikasjon. Output som JSON & CSV. |
| 3 | PDF-spesifikke tester | PDF/A- & risikoanalyse. veraPDF for ISO-validering. JHOVE for strukturell kontroll. Ekstra moduler: qpdf (kryptering), pdfinfo (metadata). |
| 4 | Resultat-motor | Én felles rapport (report.json + valgfri report.html). Bruk OpenAPI-skjema så GUI-en senere kan lese samme format. |
| 5 | Lokal AI-assistent | Offline «explain-this». Støtt Ollama eller LM Studio for å kjøre Llama- eller Mistral-modeller lokalt på M-silicon. Prompt: «Gi en kort, ikke-teknisk forklaring på hver feil + forslag til utbedring». |
| 6 | Plugin-API | Eksterne testmoduler. Enkel Go-interface: Run(ctx, inPath) → (results, error). |
| 7 | Pakking | Homebrew-formula + .pkg. Lever én `brew install arcaive` og notariserte universelle 2-binære. |
| 8 | GUI-MVP | Drag-&-drop app. Bruk Tauri: minimal Rust + webview (gir <10 MB app), deler kjerne‐binæren. |

---

## Minimum-funksjoner for v0.1 (CLI)

1. **Input**
   - `arcaive scan /path/to/bag` eller `arcaive scan myfile.pdf`.
2. **Identifikasjon**
   - Siegfried‐rapport (PRONOM ID, versjon, pUID).
3. **PDF-validering (kun PDF)**
   - veraPDF «pass/fail» pr. nivå, med rule-ID-liste.
   - JHOVE status (Well-formed / Valid / etc.).
4. **Exit-koder**
   - 0 = alt ok, 1 = advarsler, 2 = kritiske feil.
5. **Rapporter**
   - `--json`, `--html`, `--quiet`.
6. **AI-hjelp (valgfritt flagg)**
   - `--explain` kaller lokal LLM og skriver «what & how-to-fix».
7. **Auto-update rules**
   - `arcaive update signatures` henter nyeste PRONOM‐signatures (Siegfried) og veraPDF-profiles.

---

## Arkitektur-skisse

```
┌─────────────────┐
│  CLI (Cobra)    │   ← user runs commands
└──────┬──────────┘
       │
┌──────▼──────────┐
│ Core Orchestrator│
└──────┬──────────┘
    ┌──▼───┐  ┌────▼─────┐  ┌────▼────┐
    │ID Lib│  │PDF Suite │  │AI Helper│
    │(Sf)  │  │(vera/J)  │  │(Ollama) │
    └──┬───┘  └────┬─────┘  └────┬────┘
       │ JSON       │ XML        │ prompt
       └────────────┴────────────┴────────→ report builder
```

* Sf = Siegfried, vera/J = veraPDF & JHOVE

---

## Viktige avhengigheter & installasjon (macOS M-series)

| Verktøy    | Installering (Brew)        | Notater |
|------------|----------------------------|---------|
| Go 1.22+   | `brew install go`          | Krysskompilerer lett. |
| Siegfried  | `brew install siegfried`   | Innebygget PRONOM updater. |
| veraPDF    | `brew install verapdf`     | Java 17-runtime følger med. |
| JHOVE      | `brew install jhove`       | Validerer flere formater enn PDF. |
| Ollama     | `brew install ollama`      | Kjør `ollama run llama3:8b` for testen. |

ArcaiVe-binæren vil sjekke PATH og gi vennlige feilmeldinger.

---

## Videre idébank (v0.2+)

- BagIt-validasjon for hele arkivleveranser.
- Checksum-sammenligning (SHA-256 manifest).
- Image-analyse: EXIF, ICC-profiler, JHOVE-TIFF.
- Container-støtte: zip/tar/7z med rekursiv skanning.
- REST-API for integrasjon i CI.
- Electron/Tauri GUI som bruker samme CLI via IPC.

---

## Hvorfor dette kan bli bra

- **Én kommando, null Java-hakking:** Go-kjerne roper bare på verktøy med veldefinerte flags.
- **Offline & privat:** Ingen filer eller metadata trenger å forlate Mac-en.
- **Utvidbart:** PRONOM oppdateres ukentlig; AI-modellen kan byttes med `ollama pull`.
- **Vennlig for både bibliotekarer og dev-ops:** CLI for skripting, GUI for «drag-and-drop».

---

Gi beskjed hvis du vil dykke mer i et bestemt punkt (f.eks. hvordan koble AI-forklaringene inn i rapporten), eller om du heller vil bytte språk (Rust/Swift). Så kan vi begynne å codex-spike selve arcaiVe-skjelettet!

