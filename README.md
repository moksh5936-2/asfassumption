```
                      /^\/^\
                    _|__|  O|
           \/     /~     \_/ \
            \/   |  ||  |   |\
            /\   |  ||  |   |\
           /  \  |  ||  |   |/
      __  /   /  |_||__|   |_
  _  /  \/   /  /         /   \
 / \/  /\_  /  /         /    |
|  |  |  \/  /          |     |
|  |  |    \/            \    |
|  |  |     |             |   |
 \  \ |     |             |  /
  \  \|     |             | /
   |  |     |         __  |/
   |  |      \      /  \ |
   |  |       |    / __  |
    \  \      |   | /  | |
     \  \____/    |/  / /
      \_     |     /__/ /
        \__  |    |    |
           \_|    |____|
              |   |    |
              |   |    |
              |   |    |
              |   |    |
              |   |____|
              |   |    |
              |   |   ||
              |   |   ||
              |   |   ||
              |   |   ||
              |   |   ||
              |   |   ||
               \  |  /
                \ | /
                 \|/

               ASF v1.0.0
   Architecture Security Framework
   Security Assumption Discovery Engine
```

**ASF** is a deterministic, offline-first terminal application that automatically discovers hidden security assumptions in system architecture diagrams and documents. It uses STRIDE threat modeling, risk assessment, and confidence scoring to produce fully explainable security reviews.

---

## Features

- **🔍 Multi-format input** — Draw.io, Mermaid, YAML, JSON, SVG, images (OCR), TXT, PDF, DOCX
- **🧠 Deterministic analysis** — No AI for core analysis. Every result is reproducible and auditable.
- **⚠️ STRIDE threat mapping** — Proprietary rule engine: 17 category rules + 33 keyword patterns
- **📊 5×5 risk matrix** — Deterministic likelihood × impact scoring with full decomposition
- **🎯 Confidence scoring** — 4-metric calculation from evidence, rules, components, relationships
- **🔗 Evidence traceability** — Every assumption traced back to source components and relationships
- **🔄 Architect review mode** — Accept/Reject/Modified status tracking with notes
- **📤 5 export formats** — JSON, Markdown, CSV, PDF, HTML with full explainability data
- **🎨 4 themes** — Dark, Midnight, Cyber, Minimal
- **🤖 Optional local AI** — Ollama integration for enhanced analysis (fully offline)
- **🔑 Enterprise licensing** — HMAC-signed license keys
- **📦 11.85MB binary** — Single-file distribution, no runtime dependencies (except Python)

---

## Architecture

```
┌──────────────┐   ┌──────────────┐   ┌──────────────┐
│  User Input  │──▶│  ASF Engine  │──▶│  Python CLI  │
│  (diagram/   │   │  (Go TUI)    │   │  (asf.cli)   │
│   document)  │   │              │   │              │
└──────────────┘   └──────┬───────┘   └──────────────┘
                          │
                          ▼
                   ┌──────────────┐
                   │  Explainability Pipeline          │
                   │  ┌────────────────────────────┐  │
                   │  │ 1. Evidence Engine          │  │
                   │  │ 2. Assumption Justification │  │
                   │  │ 3. STRIDE Justification     │  │
                   │  │ 4. Likelihood Analysis      │  │
                   │  │ 5. Impact Analysis          │  │
                   │  │ 6. Risk Matrix (5×5)        │  │
                   │  │ 7. Confidence Scoring       │  │
                   │  └────────────────────────────┘  │
                   └──────────────┬───────────────────┘
                                  │
                    ┌─────────────┼─────────────┐
                    ▼             ▼             ▼
             ┌──────────┐  ┌──────────┐  ┌──────────┐
             │  TUI     │  │  Export  │  │  Review  │
             │  Results │  │  5 fmt   │  │  Mode    │
             └──────────┘  └──────────┘  └──────────┘
```

---

## Screenshots

```
+---------------------------------------------------+
|  Welcome to ASF                                   |
|  ┌──────────┐  ┌──────────────────────────┐      |
|  │    /\_/\ │  │  ▸ Analyze Architecture  │      |
|  │   /     \│  │    Results               │      |
|  │  | . .  ||  │    AI Settings           │      |
|  │   \___/  │  │    Settings              │      |
|  │    | |   │  │    About                 │      |
|  │    |O|   │  │    Exit                  │      |
|  │    ---   │  └──────────────────────────┘      |
|  └──────────┘                                     |
|  Architecture Security Framework                  |
+---------------------------------------------------+
  ↑↓: Navigate  •  Enter: Select  •  q: Quit
```

---

## Supported Inputs

| Type | Extensions | Parser |
|------|-----------|--------|
| Draw.io diagrams | `.drawio` | XML-based component/relationship extraction |
| Mermaid diagrams | `.mmd` | Regex-based node/edge parsing |
| YAML definitions | `.yaml`, `.yml` | Structured architecture definition |
| JSON definitions | `.json` | Structured architecture definition |
| SVG diagrams | `.svg` | XML text extraction |
| Images | `.png`, `.jpg`, `.jpeg` | Tesseract OCR |
| Text documents | `.txt`, `.md` | Raw text analysis |
| PDF documents | `.pdf` | Raw text analysis |
| Word documents | `.docx` | Raw text analysis |

## Supported Outputs

| Format | Extension | Content |
|--------|-----------|---------|
| JSON | `.json` | Full structured result with all evidence |
| Markdown | `.md` | Readable report with risk/STRIDE breakdown |
| CSV | `.csv` | Flat table for spreadsheet analysis |
| PDF | `.pdf` | Formal multi-page report |
| HTML | `.html` | Styled single-page dashboard |

---

## Installation

### Quick Install (macOS/Linux)

```bash
curl -sfL https://raw.githubusercontent.com/asfsecurity/asf/main/asf-tui/install.sh | bash
```

### Manual Install

```bash
# Download for your platform
curl -sfL https://github.com/asfsecurity/asf/releases/download/v1.0.0/asf-darwin-arm64 -o /usr/local/bin/asf
chmod +x /usr/local/bin/asf

# Install Python ASF engine (required)
cd /path/to/asf
pip install -e .

# Run
asf
```

### Prerequisites

| Component | Required | Install |
|-----------|----------|---------|
| Python 3.8+ | ✅ Core | `pip install -e .` |
| Ollama | Optional (AI) | https://ollama.ai |
| Tesseract | Optional (OCR) | `brew install tesseract` |

---

## Quick Start

```bash
# Launch ASF
asf

# Or run analysis directly from the TUI:
# 1. Select "Analyze Architecture"
# 2. Enter path to architecture file
# 3. Select analysis mode
# 4. Start analysis
# 5. Review results
# 6. Export report
```

---

## Example: Analyze a Draw.io Diagram

1. Launch ASF: `asf`
2. Select **Analyze Architecture**
3. Set Document Path: `~/project/architecture.drawio`
4. (Optional) Set Evidence Path: `~/project/access_controls.csv`
5. Select **ASF Engine Only**
6. Select **▶ Start Analysis**
7. Review results in the TUI
8. Press `e` to export as Markdown
9. Press `r` to enter Review Mode
10. Use `s` (Accept), `r` (Reject), `m` (Modified) to mark assumptions

---

## AI Integration

ASF supports optional local AI enhancement via Ollama. No cloud dependency.

```bash
# 1. Install Ollama
brew install ollama
ollama serve

# 2. Download a model from ASF's AI Settings
#    (or manually: ollama pull llama3.2:3b)

# 3. Enable AI in Settings or select "ASF Engine + Local AI" mode
```

When AI is enabled, after the deterministic analysis completes:
1. ASF builds a structured prompt from the analysis results
2. Calls `http://localhost:11434/api/generate`
3. Parses the response for additional assumptions, risk refinements, and recommendations
4. Merges AI findings (prefixed with `AI-`) into the results

---

## Configuration

Configuration is stored in `~/.asf/config.yaml`. It auto-migrates from the legacy path `~/.config/asf/config.yaml`.

```yaml
general:
  theme: Dark
  fox_style: Classic
analysis:
  depth: deep
  stride: true
  controls: true
  risk_threshold: low
ai:
  enabled: false
  active_model: ""
  installed_models: []
output:
  default: markdown
  directory: ./reports
appearance:
  theme: Dark
  fox_style: Classic
```

---

## Licensing

ASF supports enterprise licensing with HMAC-signed keys.

```bash
# Activate a license
echo 'ASF-XXXX-XXXX-XXXX-XXXX' > ~/.asf/license.key

# Check license
asf --license
```

License format: `ASF-XXXX-XXXX-XXXX-XXXX` (16 hex characters + 8-char HMAC signature)

---

## Roadmap

### Short-term (next)

- Expert validation study (10 architects × 20 architectures)
- CI/CD pipeline (GitHub Actions)
- Code signing and notarization

### Medium-term

- Multi-architecture batch analysis
- Results comparison mode
- Results database persistence
- Team collaboration features

### Long-term

- Native Go assumption extraction (remove Python dependency)
- Cloud AI provider integration (optional)
- REST API server mode
- VS Code extension

---

## Validation Status

| Metric | Status |
|--------|--------|
| Unit tests | ✅ 20 passing |
| Code quality | ✅ `go vet` clean (unverifiable — Go not on this machine) |
| Benchmark run | ✅ 2158 assumptions processed |
| Precision | ❌ Not measured |
| Recall | ❌ Not measured |
| False positive rate | ❌ Not measured |
| STRIDE accuracy | ❌ Not measured |
| Expert validation study | ❌ Not started |

See [docs/VALIDATION_STATUS.md](docs/VALIDATION_STATUS.md) for a complete, honest assessment.

---

## Project Structure

```
asf-tui/                   # Go TUI application
  main.go                  # CLI entry point
  app.go                   # TUI controller
  engine.go                # Python bridge, result builder
  parser.go                # All input format parsers
  stride.go                # STRIDE rule engine
  justify.go               # Explainability pipeline (7 engines)
  explain.go               # Data structures
  validation.go            # TUI validation mode
  review.go                # Review mode
  ai.go                    # AI enhancement
  model.go                 # Ollama manager
  localai.go               # AI model manager
  export.go                # 5 export formats
  config.go                # Configuration
  license.go               # License system
  styles.go                # 4 themes
  startup.go               # Welcome screen
  dashboard.go             # Quick actions
  analyze.go               # Analysis setup
  results.go               # Results display
  settings.go              # Settings editor
  about.go                 # About screen
  explain_test.go          # 20 unit tests
  install.sh               # Installer script
  go.mod / go.sum          # Dependencies
docs/                      # Documentation
  ARCHITECTURE.md
  EXECUTIVE_SUMMARY.md
  USER_MANUAL.md
  DEVELOPER_GUIDE.md
  TECHNICAL_REFERENCE.md
  VALIDATION_STATUS.md
  risk_model.md
  EXPLAINABILITY_ENGINE.md
  EXPLAINABILITY.md
  EXPLAINABILITY_GAP_ANALYSIS.md
  EXPLAINABILITY_READINESS_REPORT.md
  MIGRATION_GUIDE.md
  BUILD_SYSTEM.md
  INSTALLATION_ARCHITECTURE.md
  LICENSE_ARCHITECTURE.md
  SECURITY_REVIEW.md
  EXPERT_VALIDATION_STUDY.md
  MARKET_POSITIONING.md
  DEFENSIBILITY_ANALYSIS.md
release/                    # Release artifacts
  README.md
  VERSION
  checksums.txt
  install.sh
scripts/                    # Build automation scripts
  build-release.sh
  build-release.ps1
benchmark/                  # Research validation
asf/                        # Python ASF engine (v1)
PROJECT_AUDIT_REPORT.md     # Comprehensive codebase audit
RELEASE_CHECKLIST.md        # Release verification checklist
EXECUTIVE_RELEASE_REPORT.md # Executive release report
```

---

## Contributing

See [CONTRIBUTING.md](CONTRIBUTING.md).

---

## FAQ

**Q: Does ASF require internet?**  
A: No. The core analysis is fully offline. Only model downloads and AI enhancement require local network access to Ollama.

**Q: Is ASF an AI tool?**  
A: No. AI is optional and purely additive. The core analysis is deterministic rule-based Go code.

**Q: What file formats are supported?**  
A: Draw.io (.drawio), Mermaid (.mmd), YAML (.yaml/.yml), JSON (.json), SVG (.svg), images (.png/.jpg/.jpeg via OCR), and text documents (.txt/.md/.pdf/.docx).

**Q: How is this different from a vulnerability scanner?**  
A: ASF finds implicit security assumptions, not known vulnerabilities. It answers "what did my team assume about the system?" not "what CVEs exist?"

**Q: Can I use ASF without Python?**  
A: Currently no — the assumption extraction is done by the Python ASF CLI. A future version may replace this with native Go code.

**Q: How accurate is ASF?**  
A: We don't know yet. Precision, recall, and false positive rate have not been measured. See [VALIDATION_STATUS.md](docs/VALIDATION_STATUS.md).

---

## Security Notice

ASF processes architecture documents entirely locally. No data is sent to external services. AI enhancement uses a local Ollama instance — no data leaves your machine.

License validation uses HMAC and is performed locally. No phone-home mechanism exists.

## Limitations

- Python ASF CLI must be installed separately
- Image OCR requires Tesseract
- AI features require Ollama
- No human validation study yet
- Windows TUI not thoroughly tested

---

## Acknowledgements

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss)
- PDF generation via [go-pdf/fpdf](https://github.com/go-pdf/fpdf)
- STRIDE methodology based on Microsoft's threat modeling framework

---

## License

Research and educational use. Enterprise licensing available.

```
© 2026 ASF Project
```
