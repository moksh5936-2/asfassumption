# Python Engine Inventory

Complete inventory of the Python ASF engine at `asf/` (commit `6baef2a`, v1.1.0).

## Directory Structure

```
asf/
├── __init__.py              # Package marker (empty)
├── analyzer.py              # Main orchestrator (159 lines)
├── config.py                # Config class + YAML defaults (120 lines)
├── settings.py              # Config loader with env var overrides (46 lines)
├── cli/
│   └── main.py              # Click CLI (analyze command, 200+ lines)
├── models/
│   ├── __init__.py          # Exports all models (22 lines)
│   ├── enums.py             # 5 enum classes (74 lines)
│   ├── claim.py             # Claim Pydantic model (25 lines)
│   ├── assumption.py        # Assumption Pydantic model (28 lines)
│   ├── evidence.py          # Evidence Pydantic model (28 lines)
│   ├── verification.py      # Verification Pydantic model (28 lines)
│   ├── gap.py               # Gap Pydantic model (26 lines)
│   └── analysis.py          # AnalysisResult container (43 lines)
├── extraction/
│   ├── __init__.py          # Exports ClaimExtractor
│   └── claim_extractor.py   # 15 regex patterns, sentence splitting, dedup, tags (112 lines)
├── assumption/
│   ├── __init__.py          # Exports AssumptionEngine
│   └── assumption_engine.py # 8 AssumptionType classifiers, scoring, keyword extraction (104 lines)
├── evidence/
│   ├── __init__.py          # Exports EvidenceLoader, EvidenceMapper
│   ├── evidence_loader.py   # CSV/JSON loading, column auto-mapping (80+ lines)
│   ├── evidence_mapper.py   # Source type compatibility mapping (60+ lines)
│   └── schema_adapter.py    # Field name normalization (50+ lines)
├── verification/
│   ├── __init__.py          # Exports VerificationEngine
│   └── verification_engine.py # 6 type-specific verifiers, field matching (379 lines)
├── confidence/
│   ├── __init__.py          # Exports ConfidenceEngine
│   └── confidence_engine.py # Freshness/coverage/completeness scoring (100+ lines)
├── gaps/
│   ├── __init__.py          # Exports GapEngine
│   └── gap_engine.py        # Gap generation with severity mapping (89 lines)
├── graph/
│   ├── __init__.py          # Exports GraphModel
│   └── graph_model.py       # NetworkX-based graph builder (129 lines)
├── ingestion/
│   ├── __init__.py          # Exports IngestionPipeline
│   ├── pipeline.py          # Format detection, text extraction dispatch (69 lines)
│   ├── csv_parser.py        # CSV reading
│   ├── json_parser.py       # JSON reading
│   ├── txt_parser.py        # Plain text reading
│   ├── pdf_parser.py        # PDF text extraction (pdfplumber)
│   └── docx_parser.py       # DOCX text extraction (python-docx)
├── db/
│   ├── __init__.py
│   └── database.py          # SQLite persistence layer
├── api/
│   ├── __init__.py
│   └── server.py            # FastAPI REST API
└── llm/
    ├── __init__.py
    ├── client.py            # OpenAI/Ollama LLM client
    └── prompts.py           # Prompt templates
```

## Tests

```
tests/
├── test_analyzer.py         # Integration tests
├── test_assumption.py       # Assumption engine tests
├── test_confidence.py       # Confidence engine tests
├── test_extraction.py       # Claim extraction tests
├── test_gaps.py             # Gap engine tests
├── test_graph.py            # Graph model tests
├── test_ingestion.py        # Ingestion pipeline tests
├── test_evidence.py         # Evidence loader tests
├── test_models.py           # Model creation/validation tests
├── test_verification.py     # Verification engine tests
├── test_cli.py              # CLI tests
├── validation_harness.py    # Cross-validation utilities
└── __init__.py
```

**59 tests total, all passing.**

## Key Dependency Versions

| Package | Version | Purpose |
|---------|---------|---------|
| pydantic | 2.x | Data models with validation |
| click | 8.x | CLI framework |
| rich | 13.x | Console formatting |
| networkx | 3.x | Graph model |
| pdfplumber | 0.11.x | PDF text extraction |
| python-docx | 1.x | DOCX text extraction |
| fastapi | 0.x | REST API server |
| pyyaml | 6.x | YAML config parsing |

## Migration Status

| Module | Go Port | Parity |
|--------|---------|--------|
| models/ | `asf/models/` | 100% (7 models + 5 enums) |
| extraction/ | `asf/extraction/` | 100% (15 patterns, dedup, tags) |
| assumption/ | `asf/assumption/` | 100% (8 classifiers, scoring, keywords) |
| evidence/ | `asf/evidence/` | 100% (loader + mapper) |
| verification/ | `asf/verification/` | 100% (6 type-specific checks) |
| confidence/ | `asf/confidence/` | 100% (freshness/coverage/completeness) |
| gaps/ | `asf/gaps/` | 100% (severity mapping matches) |
| graph/ | `asf/graph/` | 100% (72 nodes, 136 edges match) |
| ingestion/ | `asf/ingestion/` | Basic — needs proper PDF/DOCX libs |
| analyzer/ | `asf/analyzer/` | 100% (integration matches) |
| cli/ | Not ported | TUI still handles interaction |
| db/ | Not ported | SQLite not yet migrated |
| api/ | Not ported | FastAPI not yet migrated |
| llm/ | Not ported | AI enhancement not yet migrated |
| config.py/settings.py | `config.go` | Native config handles this |
