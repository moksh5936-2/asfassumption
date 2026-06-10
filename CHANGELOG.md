# Changelog

All notable changes to ASF are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/),
and this project adheres to [Semantic Versioning](https://semver.org/).

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
