# PDF/DOCX Support Fix Report (B7)

## Problem

The README at `README.md` claimed support for PDF and DOCX documents with "Raw text analysis" — implying structured text extraction was functional:

> | PDF documents | `.pdf` | Raw text analysis |
> | Word documents | `.docx` | Raw text analysis |

However, the implementation in `asf-tui/parser.go:255-264` (`parseTextFile`) simply calls `os.ReadFile()` and treats the raw bytes as a Go string. PDF and DOCX are binary container formats — raw reads produce garbled, non‑printable text that is useless for downstream analysis.

## Decision

Rather than adding real PDF/DOCX parsing (e.g., `unidoc`, `pdfcpu`, or `gooxml` libraries — which would introduce significant dependencies and complexity), we chose **Option B: Remove/qualify the claim**. The feature was always aspirational and the false claim could erode user trust.

## Changes Made

### 1. README.md

- **Supported Inputs table** — Changed both PDF and DOCX rows from "Raw text analysis" to "Binary raw text (limited — text extraction not implemented)"
- **Footnote** — Added a footnote below the table:
  > PDF and DOCX support reads raw binary content. Structured text extraction is not yet implemented. Use .txt, .yaml, .json, .drawio, .mmd, or .svg inputs for best results.
- **Features section** — Appended `\*` to `PDF` and `DOCX` in the multi-format bullet point to signal the limitation.
- **FAQ** — Updated the supported-formats answer to note the limitation.

### 2. asf-tui/parser.go

- **New helper `isPrintableText`** — Scans the byte slice and computes the ratio of printable ASCII characters (0x20–0x7E plus newline, carriage‑return, tab). Returns `true` if the ratio exceeds 0.5.
- **Modified `parseTextFile`** — When the file extension is `.pdf` or `.docx` **and** the content fails the printable-text check, the extracted text is prepended with:
  ```
  [WARNING: Raw binary content — text extraction may produce garbled output]
  ```
  This ensures users of these formats see a clear indication that the output may be garbage.

### 3. New Report

This file — `docs/PDF_DOCX_SUPPORT_FIX_REPORT.md` — documents the issue, decision, and changes.

## Future Recommendation

If real PDF/DOCX support is desired in the future, consider:

| Format | Recommended Library | Notes |
|--------|-------------------|-------|
| PDF | `github.com/pdfcpu/pdfcpu` or `github.com/unidoc/unipdf` | pdfcpu is OSS; unidoc requires a commercial license for some features |
| DOCX | `github.com/unidoc/unioffice` or `github.com/nguyenthenguyen/docx` | unioffice is commercial; `docx` is a simpler MIT reader |

These should be integrated by adding dedicated parse functions (e.g., `parsePDF(path)`, `parseDOCX(path)`) dispatched from `ParseArchitecture`, rather than routing through `parseTextFile`.
