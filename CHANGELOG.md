# Changelog

All notable changes to ASF are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [2.1.0] — 2026-06-12

### Added

- YAML/JSON archive definition now parses 13 fields (up from 5): assumptions, security_controls, compliance, expected_results, validation_criteria, notes
- Explicit assumptions from YAML/JSON classified into 8 security types, assessed for risk and STRIDE-mapped
- Security controls enrichment — generated controls enhanced with YAML-defined specific controls
- Compliance framework output — structured display of compliance targets (HIPAA, SOC2, ISO27001, etc.)
- Expected results validation — compares actual analysis against declared benchmarks
- Low and Medium risk tracking separated from Critical/High for accurate distribution
- `processExplicitAssumptions`, `classifyExplicitAssumption`, `assessExplicitRisk`, `buildComplianceOutput`, `buildValidationSummary`, `enhanceControlsWithSecurityControls` methods
- `ingestion_test.go` with 40 tests covering all new parsing, classification, risk, compliance, and validation paths
- `testdata/asftest.yaml` — comprehensive healthcare architecture test file with 30 assumptions, 9 components, 7 control categories, 3 compliance frameworks
- `docs/STRUCTURED_YAML_INGESTION_IMPROVEMENT_REPORT.md`

### Changed

- `ArchDescription` struct: 6 new optional fields (backward compatible)
- `archDefinition` struct: 7 new optional YAML/JSON fields (backward compatible)
- `AnalysisResult` struct: added `MediumCount`, `LowCount` fields
- `buildResult` processes explicit assumptions through dedup, classification, risk, and explainability pipeline
- All export formats (Markdown, HTML, PDF, CSV) include Medium/Low risk distribution
- TUI results header shows Medium and Low risk counts
- `mergeAIResults` tracks Medium/Low for AI-generated assumptions

### Fixed

- `normalizeText` now correctly strips trailing periods from individual tokens
- `toFloat` handles `int32` and `uint` types
- `buildValidationSummary` reports "all expected criteria met" for empty expected results
- `buildResult` computes TotalAssumptions from actual deduplicated count

### Changed

- Removed Python ASF engine bridge entirely — ASF is now a true Go-native single-binary
- All analysis now uses native Go engine (no `python3`, `pip`, `venv`, or `PYTHONPATH` required)
- Replaced Python-based install with pure Go binary downloads
- Cleaned up docter output to remove Python engine references
- Increased binary size to ~9MB (up from ~8MB) due to native engine inclusion
- Version bumped from 1.1.0 to 2.0.0 (breaking change: Python dependency removed)

### Removed

- `callPythonCLI`, `discoverPythonPath`, `preFlightCheck` from engine.go
- Python engine section from doctor.go (findPython, downloadEngineBundle, etc.)
- `PythonPath` field from config.go
- `asf-python-engine-*.tar.gz` release artifacts
- `scripts/package-python-engine.sh`
- `package-python-engine` job from CI/CD

### Fixed

- Directory input crash: `asf analyze <directory>` now expands supported files
- Help text: `asf doctor --fix` no longer mentions "install Python engine"
- 5 certification blockers resolved for Go-native single-binary certification

### Added

- `claims[]` array in native Go JSON output matching Python schema
- Cross-platform release assets at ~33% smaller due to stripped debug symbols

## [1.0.0] — 2026-06-10

### Added

- Complete TUI application with 9 views (Startup, Dashboard, Analyze, Results, Review, Settings, Local AI, About, Export)
- Architecture diagram parser supporting 10+ formats (Draw.io, Mermaid, YAML, JSON, SVG, Images/OCR, TXT, PDF, DOCX)
- Python ASF CLI bridge for deterministic assumption extraction
- Explainability pipeline with 7 engines:
  - Evidence Engine (component/relationship/concept matching)
  - Assumption Justification (human-readable rationale)
  - STRIDE Justification (per-category reasoning with confidence)
  - Likelihood Analyzer (3-factor scoring 1-5)
  - Impact Analyzer (3-factor scoring 1-5)
  - Risk Matrix (5×5 deterministic matrix)
  - Confidence Engine (4-metric scoring)
- Proprietary STRIDE rule engine: 17 category rules + 33 keyword patterns
- 5 export formats (JSON, Markdown, CSV, PDF, HTML) with full explainability data
- Architect review mode with Accept/Reject/Modified status tracking
- Validation data collection for precision/recall studies
- Control detail generation with mitigated assumption and STRIDE tracking
- 4 themes (Dark, Midnight, Cyber, Minimal)
- Settings screen with 9 configurable options
- Optional local AI enhancement via Ollama API
- Model manager (download/list/delete/activate Ollama models)
- HMAC-based license system (ASF-XXXX-XXXX-XXXX-XXXX format)
- Auto-config migration from legacy path
- 20 unit tests covering all 7 explainability engines
- Cross-platform builds (Linux/macOS/Windows, AMD64/ARM64)
- Install script (install.sh)
- Multi-platform documentation suite

### Changed

- Pivot from mock assumptions to real Python ASF engine bridge
- Risk model from flat Medium to 5×5 matrix distribution
- Dashboard version display updated to v1.0.0
- Config path from `~/.config/asf/config.yaml` to `~/.asf/config.yaml`

### Fixed

- Empty evidence array handling (`interface{}` in JSON struct)
- Architecture description initialization ordering
- Various edge cases in confidence engine and risk matrix
- Progress bar display during analysis
- View history backstack navigation

### Known Issues

- Python ASF CLI dependency required for analysis
- Image OCR requires Tesseract
- AI features require Ollama
- No human validation study completed
- Precision and recall not yet measured
