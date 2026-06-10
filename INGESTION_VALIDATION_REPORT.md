# Ingestion Validation Report: Python vs Go Text Extraction Engines

**Date:** 2026-06-10
**Scope:** All 5 shared document formats (TXT, PDF, DOCX, CSV, JSON)
**Corpus:** `/Users/moksh/Project/cybersec/asf-tui/testdata/parity/samples/`
**Python Engine:** `asf.ingestion.pipeline.IngestionPipeline` (asf/ingestion/)
**Go Engine:** `asf-tui/asf/ingestion/parser.go` (compiled as `asf analyze`)

---

## 1. TXT — Plain Text

### Engines
| Aspect | Python | Go |
|---|---|---|
| File | `txt_parser.py:7-9` | `parser.go:72-78` |
| Function | `Path.read_text(encoding="utf-8")` | `os.ReadFile(path)` → `string(data)` |
| Approach | Standard UTF-8 file read | Standard UTF-8 file read |

### Result
- **Raw Content Identical?** YES
- **Words/Content Identical?** YES (identical byte-for-byte)
- **Length:** 1,375 bytes in both engines

### Analysis Output Parity
| Metric | Python | Go |
|---|---|---|
| Claims found | 17 | 17 |
| Assumptions | 17 | 17 |
| Verifications (UNKNOWN) | 17 | 17 |
| Gaps | 0 | 0 |

Assumption texts are identical word-for-word; only whitespace in the concatenated `section_title + "\n\n" + sentence` formatting varies trivially.

**Verdict:** FULL PARITY — identical raw read, identical analysis output.

---

## 2. PDF — Portable Document Format

### Engines
| Aspect | Python | Go |
|---|---|---|
| File | `pdf_parser.py:7-20` | `parser.go:162-180` |
| Library | `pdfplumber` | `github.com/ledongthuc/pdf` |
| Approach | `page.extract_text()` per page → `"\n\n".join(text_parts)` | `r.GetPlainText()` → `io.ReadAll` continuous text |

### Text Extraction Difference
Python builds text page-by-page, joining with double newlines. Go reads all pages as a single continuous stream. This means Python inserts `\n\n` between pages while Go does not (or relies on the PDF library's internal pagination handling).

On the `finance_policy.pdf` sample (single page), Python's `pdfplumber` additionally extracts a header line `"Finance Access Control Policy"` (likely from PDF document metadata/properties) that precedes the visible body text, while Go's library does not extract this metadata header.

### Result
- **Words/Content Identical?** YES (ignoring whitespace and the extra metadata header)
- **Whitespace Difference?** YES — Python uses `\n\n` between pages; Go uses continuous text. Python may include extra metadata from PDF properties.
- **Claims/Assumptions Match?** YES — both engines found 17 claims/assumptions with identical content

### Analysis Output Parity
| Metric | Python | Go |
|---|---|---|
| Claims found | 17 | 17 |
| Assumptions | 17 | 17 |
| Verifications | 17 UNKNOWN | 17 UNKNOWN |
| Gaps | 0 | 0 |

**Verdict:** CONTENT PARITY — words and claims identical; whitespace and an extra metadata header differ as expected by library choice.

---

## 3. DOCX — Office Open XML Document

### Engines
| Aspect | Python | Go |
|---|---|---|
| File | `docx_parser.py:7-16` | `parser.go:182-239` |
| Library | `python-docx` | Manual ZIP + XML (`encoding/xml`) |
| Approach | `Document(str(path))`, iterate `.paragraphs`, filter `p.text.strip()`, `"\n\n".join()` | Open ZIP, parse `word/document.xml`, iterate `<w:p>`→`<w:r>`→`<w:t>`, `"\n".join()` |

### Key Implementation Differences

**Python (python-docx):**
- Filters out empty/whitespace-only paragraphs (`if p.text.strip()`)
- Joins surviving paragraphs with `"\n\n"` (double newline)
- Each paragraph's runs are concatenated by the library

**Go (ZIP+XML):**
- Includes ALL paragraphs (including empty ones that represent blank lines)
- Joins paragraphs with `"\n"` (single newline)
- Extracts `<w:t>` text from each `<w:r>` run within each `<w:p>` paragraph
- Does not handle text formatting, images, headers/footers, or embedded content

### Result
- **Words/Content Identical?** YES (ignoring whitespace)
- **Whitespace Difference?** YES — Python produces `"\n\n"` between paragraphs and omits blank-line-only paragraphs; Go produces `"\n"` between ALL paragraphs including empties.
- **Claims/Assumptions Match?** YES — both engines found 17 claims/assumptions with identical text content

### Analysis Output Parity
| Metric | Python | Go |
|---|---|---|
| Claims found | 17 | 17 |
| Assumptions | 17 | 17 |
| Verifications | 17 UNKNOWN | 17 UNKNOWN |
| Gaps | 0 | 0 |

**Verdict:** CONTENT PARITY — words and claims identical; interior whitespace/paragraph separation differs as expected by implementation approach.

---

## 4. CSV — Comma-Separated Values

### Engines
| Aspect | Python | Go |
|---|---|---|
| File | `csv_parser.py:11-20` | `parser.go:80-133` |
| Library | `csv.DictReader` | `encoding/csv` |
| Text approach | Raw file read (`filepath.read_text()`) | CSV parse + reformat: `"headers: " + row0 + "\n" + rows1..N + "\n"` |
| Records approach | `csv.DictReader` → `list[dict[str, Any]]` | `csv.ReadAll` → manual header map → `[]map[string]interface{}` |

### Text Extraction Difference
| Aspect | Python text output | Go text output |
|---|---|---|
| Format | Raw CSV content as-is | Reformatted with "headers: " prefix |
| Example (backup_config.csv) | `resource,configuration,enabled,status,frequency\npayroll-db,encrypted,true,active,daily\n...` | `headers: resource, configuration, enabled, status, frequency\npayroll-db, encrypted, true, active, daily\n...` |
| Commas | No spaces after commas | Space after each comma (from `strings.Join(row, ", ")`) |
| Header row | Included verbatim as first line | Prefixed with `"headers: "` |

### Records Parity
| Sample | Python rows | Go rows | Content match? |
|---|---|---|---|
| backup_config.csv | 10 | 10 | YES — identical key-value pairs |
| mfa_status.csv | 10 | 10 | YES |
| network_exposure.csv | 10 | 10 | YES |
| payroll_acl.csv | 10 | 10 | YES |

Record field names and values are identical. Python returns `str` values; Go returns `interface{}` (concrete type `string`) — logically equivalent.

### Analysis Output Parity
Both engines produce identical claims from CSV files when they exist (some small CSVs produce 0–1 claims depending on whether the content matches declarative patterns).

| Sample | Python claims | Go claims |
|---|---|---|
| backup_config.csv | 1 | 1 |
| mfa_status.csv | 0 | 0 |
| network_exposure.csv | 1 | 1 |
| payroll_acl.csv | 0 | 0 |

**Verdict:** RECORD PARITY — records are identical. Text output differs (raw vs reformatted) but this is a design choice for the `parse_text` method. Analysis output parity is maintained.

---

## 5. JSON — JavaScript Object Notation

### Engines
| Aspect | Python | Go |
|---|---|---|
| File | `json_parser.py:9-20` | `parser.go:135-160` |
| Library | `json.loads` / `json.dumps` | `encoding/json` / `json.Unmarshal` |
| Text approach | `json.loads` → `json.dumps(data, indent=2)` | `os.ReadFile` → raw string |
| Records approach | `json.loads` → return parsed data | `json.Unmarshal` → array (or wrapped single object) |

### Text Extraction Difference
| Aspect | Python text output | Go text output |
|---|---|---|
| Format | Pretty-printed with `json.dumps(indent=2)` | Raw file content (no reformatting) |
| Example (iam_export.json) | Re-indented with 2-space indent | Raw content as written in source file |

### Records Parity
Both engines produce identical parsed JSON:
- Python: `list[dict]` containing `{"users": [...], "groups": {...}}`
- Go: `[]map[string]interface{}` containing `{"users": [...], "groups": {...}}`
- Structure and values are identical

### Analysis Output Parity
Both engines find 0 claims/assumptions from `iam_export.json` because JSON data does not contain natural-language declarative sentences meeting the extractor threshold.

**Verdict:** RECORD PARITY — records are identical. Text output differs (pretty-printed vs raw) but this is a design choice. Analysis output parity is maintained.

---

## Summary Comparison Matrix

| Format | Words/Content Identical? | Whitespace Identical? | Records Identical? | Analysis (Claims) Identical? | Analysis (Assumptions) Identical? |
|---|---|---|---|---|---|
| TXT | YES | YES | N/A | YES (17) | YES (17) |
| PDF | YES | NO¹ | N/A | YES (17) | YES (17) |
| DOCX | YES | NO² | N/A | YES (17) | YES (17) |
| CSV | NO³ | NO³ | YES | YES | YES |
| JSON | NO⁴ | NO⁴ | YES | YES (0) | YES (0) |

¹ Python pdfplumber adds `"Finance Access Control Policy"` metadata header and uses `\n\n` page separation; Go ledongthuc/pdf does not extract metadata headers.
² Python python-docx filters empty paragraphs and joins with `\n\n`; Go ZIP+XML keeps all paragraphs and joins with `\n`.
³ Python returns raw CSV text; Go reformats with `"headers: "` prefix and space after commas.
⁴ Python pretty-prints with `json.dumps(indent=2)`; Go returns raw file content.

## Bugs Found

1. **Go `parseCSVToText` reformatting changes text semantics** (`parser.go:93-99`): The `parseCSVToText` function reformats CSV content with `"headers: "` prefix and spaces after commas. This means Go's `ParseText` output for CSV is NOT a valid CSV file, unlike Python which returns raw valid CSV. This could affect downstream consumers that expect valid CSV from `ParseText`.

2. **Go `parseCSVToRecords` potential panic on ragged rows** (`parser.go:124-129`): If a data row has fewer columns than the header row, accessing `row[i]` will panic with an index-out-of-bounds error. Python's `csv.DictReader` handles this gracefully by returning `None` for missing fields. No mitigation is present.

3. **Go `parseJSONToText` does not normalize whitespace** (`parser.go:135-141`): Unlike Python which round-trips through `json.loads`/`json.dumps` to produce consistent formatting, Go returns the raw file content. This means Go's text output depends on the original file's formatting, while Python normalizes it.

4. **Go DOCX parser does not handle missing `word/document.xml` gracefully** (`parser.go:222-224`): If the DOCX file is malformed and lacks the required XML entry, the function returns an error. However, it does not close the ZIP reader if `break` is never reached before a potential early return — there is a `defer rc.Close()` inside the loop body which could leak if the file is not found (minor resource leak).

5. **Python DOCX parser silently drops empty paragraphs** (`docx_parser.py:15`): The `if p.text.strip()` filter removes paragraphs that contain only whitespace. This means the Python output loses blank-line formatting from the original document. While this is often desirable for analysis, it is a behavioral difference from Go's approach.

## Recommendations

1. Standardize `ParseText` for CSV: Either both return raw CSV or both reformat. Raw CSV is preferred for round-trip compatibility.
2. Standardize `ParseText` for JSON: Both should normalize through parse/serialize to produce consistent formatting.
3. Add bounds checking in Go `parseCSVToRecords` to prevent panic on ragged CSV rows.
4. Consider whether the DOCX paragraph filtering difference (Python filters empties, Go keeps all) should be reconciled.
5. The PDF metadata header difference (`"Finance Access Control Policy"` extracted by pdfplumber but not by ledongthuc/pdf) should be noted as a known library behavior difference — it does not affect analysis output.
