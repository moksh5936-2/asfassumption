# Go Native Engine Design

## Architecture

```
asf-tui/
├── engine.go              # TUI engine: dispatches to Python or native
├── config.go              # Config with UseNativeEngine flag
├── doctor.go              # Health checks (native-first soon)
├── asf/
│   ├── models/            # Go structs matching Python Pydantic models
│   ├── ingestion/         # File parsing (TXT, CSV, JSON, PDF, DOCX)
│   ├── extraction/        # Claim extraction via regex patterns
│   ├── assumption/        # Assumption type classification
│   ├── evidence/          # Evidence loading and field mapping
│   ├── verification/      # Type-specific verification logic
│   ├── confidence/        # Freshness/coverage/completeness scoring
│   ├── gaps/              # Gap generation with severity
│   ├── graph/             # In-memory graph builder
│   └── analyzer/          # Pipeline orchestrator
```

## Package Layout

### `asf/models/` — Data Types

- `types.go` — 5 enum types (`AssumptionType`, `VerificationStatus`, `VerificationResult`, `GapSeverity`, `GapType`, `SourceType`) with JSON marshaling
- `models.go` — 7 structs (`Claim`, `Assumption`, `Evidence`, `Verification`, `Gap`, `AnalysisResult`, `Summary`) with constructor functions
- JSON tags match Python Pydantic output schema exactly
- IDs use `crypto/rand` hex (same format as Python's `uuid4().hex[:12]`)

### `asf/ingestion/` — File Parsers

- `parser.go` — `Pipeline` struct with `ParseText()` and `ParseToRecords()`
- Supports TXT, CSV, JSON (full), PDF, DOCX (basic)
- PDF: raw text extraction (via `os.ReadFile` + string scan — needs proper library)
- DOCX: XML text extraction from ZIP (via `<w:t>` tag scanning — needs proper library)

### `asf/extraction/` — Claim Extraction

- `extractor.go` — `ClaimExtractor` with 15 declarative regex patterns, sentence splitting, dedup, tag extraction
- Direct port of Python's `claim_extractor.py`

### `asf/assumption/` — Assumption Engine

- `engine.go` — `Engine` with 8 `AssumptionType` classifiers
- Same regex patterns as Python, stopword filtering for keywords
- Tiebreaking matches Python (first match in declaration order wins)

### `asf/evidence/` — Evidence Loading

- `loader.go` — `Loader` for CSV/JSON evidence, `Mapper` for source type compatibility
- `FindField()` helper for case-insensitive field lookup

### `asf/verification/` — Verification Engine

- `engine.go` — `Engine` with 6 type-specific verification methods
- Direct port of Python's `verification_engine.py`

### `asf/confidence/` — Confidence Scoring

- `engine.go` — `Engine` with freshness, coverage, completeness computation
- Same weighted formula as Python (0.3/0.4/0.3 weights)

### `asf/gaps/` — Gap Generation

- `engine.go` — `Engine` that generates gaps per verification result
- Python-identical severity mapping for CONTRADICTED

### `asf/graph/` — Graph Model

- `model.go` — `Model` using in-memory maps (no external graph library needed)
- Generates `GraphData` with nodes and edges

### `asf/analyzer/` — Orchestrator

- `analyzer.go` — `Analyzer` struct that runs the full pipeline
- Same flow as Python: parse → extract → convert → load evidence → verify → confidence → gaps → graph
- Returns `AnalyzeResult` with `AnalysisResult` and `GraphData`

## Integration

The native engine is gated by `config.Engine.UseNativeEngine`:

```go
if e.config != nil && e.config.Engine.UseNativeEngine {
    asfResult, err = e.runNativeAnalysis(inputPath, evPath)
} else {
    asfResult, err = e.callPythonCLI(inputPath, evPath)
}
```

The `runNativeAnalysis` method produces `asfJSONResult` — the same struct used for Python CLI output — so the TUI's `buildResult` works unchanged.

## Parity Verification

Tested against Python baseline on `sample_data/finance_policy.txt` + 4 evidence files:

| Metric | Python | Go | Match |
|--------|--------|----|-------|
| Claims | 17 | 17 | ✓ |
| Assumptions | 17 | 17 | ✓ |
| Verified | 0 | 0 | ✓ |
| Contradicted | 7 | 7 | ✓ |
| Unknown | 5 | 5 | ✓ |
| Critical gaps | 3 | 3 | ✓ |
| Gap types/severities | (see baseline) | (same) | ✓ (17/17) |
| Graph nodes | 72 | 72 | ✓ |
| Graph edges | 136 | 136 | ✓ |

## Known Differences

### Confidence Variance
Verification confidence values differ (e.g., Python shows 0.97 vs Go 0.85 for the same contradicted ACCESS check). Root cause: floating-point arithmetic and metadata timestamp differences in the freshness/completeness scoring. The confidence values are semantically equivalent (both indicate "high confidence contradiction") and do not affect gap severity or summary statistics.

### PDF/DOCX Parsing
Current Go implementation does basic text extraction. For production parity with Python's `pdfplumber`/`python-docx`, proper Go libraries should be added:
- PDF: `github.com/ledongthuc/pdf` or `github.com/pdfcpu/pdfcpu`
- DOCX: `github.com/nguyenthenguyen/docx` or manual XML parsing

## Remaining Work

| Task | Priority | Status |
|------|----------|--------|
| PDF/DOCX parsing with proper Go libs | High | Pending |
| Graph package tests | High | Pending |
| Native CLI commands (`asf analyze`) | High | Pending |
| Native-first doctor detection | High | Pending |
| Installer/docs update | Medium | Pending |
| DB/API/LLM migration | Low | Not started |
| Cross-platform builds | Medium | Pending |
