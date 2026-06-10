# Changelog

All notable changes to ASF are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

## [2.0.0] — 2026-06-11

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
