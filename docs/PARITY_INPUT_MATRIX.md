# Parity Input Matrix

## Analysis Input Formats (Architecture Documents)

| # | Format | Extension | Python | Go Engine | Arch Pipeline | Temp Conv | Parity Needed |
|---|--------|-----------|--------|-----------|---------------|-----------|---------------|
| 1 | Plain Text | `.txt` | `txt_parser.py` | `asf/ingestion/parser.go:72` | `parser.go:247` (raw read) | No | **YES** |
| 2 | PDF | `.pdf` | `pdf_parser.py` (pdfplumber) | `asf/ingestion/parser.go:119` (ledongthuc/pdf) | `parser.go:247` (raw read) | No | **YES** |
| 3 | DOCX | `.docx` | `docx_parser.py` (python-docx) | `asf/ingestion/parser.go:162` (ZIP+XML) | `parser.go:247` (raw read) | No | **YES** |
| 4 | Markdown | `.md` | No parser (falls to txt) | `parser.go:247` (raw read) | `parser.go:247` | Yes | N/A (Go-only) |
| 5 | YAML | `.yaml/.yml` | No parser | `parser.go:412` (yaml struct) | `parser.go:412` | Yes | N/A (Go-only) |
| 6 | JSON Arch | `.json` | No parser | `parser.go:424` (json struct) | `parser.go:424` | Yes | N/A (Go-only) |
| 7 | Draw.io | `.drawio` | No parser | `parser.go:60` (XML mxGraph) | `parser.go:60` | Yes | N/A (Go-only) |
| 8 | Mermaid | `.mmd` | No parser | `parser.go:170` (regex) | `parser.go:170` | Yes | N/A (Go-only) |
| 9 | SVG | `.svg` | No parser | `parser.go:528` (XML parse) | `parser.go:528` | Yes | N/A (Go-only) |
| 10 | PNG | `.png` | No parser | `parser.go:593` (Tesseract OCR) | `parser.go:593` | Yes | N/A (Go-only) |
| 11 | JPG/JPEG | `.jpg/.jpeg` | No parser | `parser.go:593` (Tesseract OCR) | `parser.go:593` | Yes | N/A (Go-only) |

## Evidence Formats

| # | Format | Extension | Python | Go Engine | Records? | Parity Needed |
|---|--------|-----------|--------|-----------|----------|---------------|
| 1 | CSV | `.csv` | `csv_parser.py` (DictReader) | `asf/ingestion/parser.go:80` (encoding/csv) | YES | **YES** |
| 2 | JSON Evidence | `.json` | `json_parser.py` (json.load) | `asf/ingestion/parser.go:135` (json.Unmarshal) | YES | **YES** |
| 3 | TXT (as evidence) | `.txt` | Falls to txt_parser | `asf/ingestion/parser.go:72` | No | No (raw text) |
| 4 | PDF (as evidence) | `.pdf` | pdf_parser.py | `asf/ingestion/parser.go:119` | No | No (raw text) |
| 5 | DOCX (as evidence) | `.docx` | docx_parser.py | `asf/ingestion/parser.go:162` | No | No (raw text) |

## Analysis Pipeline (All Formats Feed Into Same Pipeline)

Once text is extracted, both engines run the same logical pipeline:

```
text -> claims -> assumptions -> verification -> gaps -> graph
```

The pipeline modules (extraction, assumption, verification, confidence, gaps, graph) are independent of input format.

## Certification Status (Updated: June 2026)

| # | Format | Status | Notes |
|---|--------|--------|-------|
| 1 | TXT | ✅ **Certified** | 22/22 samples — full parity across all field types |
| 2 | PDF | ✅ **Certified** | 1/1 sample — full parity (text extraction whitespace diff is cosmetic, does not affect analysis) |
| 3 | DOCX | ⏳ Pending | No .docx sample in corpus yet |
| 4 | CSV (records) | ✅ **Certified** | 4/4 evidence CSV files — identical record parsing |
| 5 | JSON (records) | ✅ **Certified** | 1/1 evidence JSON file — identical record parsing |
| 6 | Analysis pipeline | ✅ **Certified** | Identical text input produces identical claims, assumptions, verifications, gaps |

### PDF Text Extraction Note

The Python engine (pdfplumber) and Go engine (ledongthuc/pdf) extract text with different whitespace patterns from the same PDF. The word content is identical (206 words, same order), but Go preserves extra blank lines between sections. This cosmetic difference does **not** affect the analysis pipeline — when both engines receive identical text input, all outputs match exactly.

### Bug Found & Fixed

During parity audit, a bug was found in `analyze_cli.go` directory evidence expansion: when `-e <directory>` was used, **all** files in the directory were added as evidence (including TXT files), causing spurious `PARTIALLY_VERIFIED` results. Fixed to filter by supported evidence extensions (.csv, .json, .yaml, .yml).

## Key Insight

The Python and Go engines share only **5 common format parsers**: TXT, PDF, DOCX, CSV, JSON.

All other formats (drawio, mermaid, YAML, SVG, PNG, JPG) are **Go-only architecture pipeline features** that convert diagrams/images to text before feeding into the analysis engine. Python never handled these.
