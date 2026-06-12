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

               ASF v2.2.0
   Architecture Security Framework
   Security Assumption Discovery Engine
```

**ASF** is a deterministic, offline-first terminal application that automatically discovers hidden security assumptions in system architecture diagrams and documents. It uses STRIDE threat modeling, risk assessment, and confidence scoring to produce fully explainable security reviews.

---

## Features

- **🔍 Multi-format input** — Draw.io, Mermaid, YAML, JSON, SVG, images (OCR), TXT, MD, PDF, DOCX
- **🧠 Deterministic analysis** — No AI for core analysis. Every result is reproducible and auditable.
- **⚠️ STRIDE threat mapping** — Rule engine: 17 category rules + 34 keyword patterns
- **🛡️ Trust Boundary Intelligence** — Auto-discovers trust zones, boundaries, and control gaps
- **🔍 Contradiction Detection** — Detects conflicting security claims across architecture layers
- **⚔️ Threat Modeling Intelligence** — 12-category threat generation with STRIDE correlation and severity scoring
- **🔗 Attack Path Discovery** — Discovers attacker journeys from entry points to target assets with threat chaining and MITRE mapping
- **📊 5×5 risk matrix** — Deterministic likelihood × impact scoring with full decomposition
- **🎯 Confidence scoring** — 4-metric calculation from evidence, rules, components, relationships
- **🔗 Evidence traceability** — Every assumption traced back to source components and relationships
- **🔄 Architect review mode** — Accept/Reject/Modified status tracking with notes
- **📤 5 export formats** — JSON, Markdown, CSV, PDF, HTML with full explainability data
- **🎨 4 themes** — Dark, Midnight, Cyber, Minimal
- **🤖 Optional local AI** — Ollama integration for enhanced analysis (fully offline)
- **🔑 Demo licensing** — HMAC + Ed25519 license keys (demo-grade, cryptographically extractable)
- **📦 ~8–10MB binary** — Single-file distribution, no runtime dependencies

---

## Architecture

```
┌──────────────┐   ┌───────────────────┐
│  User Input  │──▶│  ASF Engine (Go)  │
│  (diagram/   │   │  Native single    │
│   document)  │   │  binary, no deps  │
└──────────────┘   └────────┬──────────┘
                           │
                           ▼
                    ┌───────────────┐
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

## Intelligence Engines

ASF v2.2.0 includes eleven deterministic intelligence engines that run in sequence during analysis:

### 1. Intelligence Engine V3 (65%)
- **Taxonomy-driven assumption discovery** — 24 categories, 5 severity levels
- **Domain packs** — Context-aware reasoning for healthcare, fintech, SaaS, infrastructure
- **Quality scoring** — Assumption scoring with evidence, confidence, relevance
- **Deduplication** — Normalizes and merges explicit/native assumptions
- **Source metadata** — Every assumption traced to source type, section, index, file

### 2. Contradiction Intelligence Engine (70%)
- **12 claim extraction patterns** — Keyword-based claim identification from assumptions
- **8+ detection rules** — MFA, encryption, access control, networking, compliance contradictions
- **Implied contradictions** — Detects implicit conflicts (e.g., "no MFA" + "privileged access")
- **Trust boundary contradictions** — Cross-zone control and authentication gaps
- **Control contradictions** — Missing controls at declared boundaries
- **Compliance contradictions** — Framework requirement vs. control gaps

### 3. Trust Boundary Intelligence Engine (75%)
- **17 trust zone types** — Internet, DMZ, Internal, Secure, Identity, Data, etc.
- **11 boundary crossing types** — Internet→DMZ, DMZ→Internal, Identity→Internal, etc.
- **Risk scoring** — PHI/PCI sensitivity boost, identity boost, component count boost
- **Weakness detection** — Missing controls, missing assumptions, boundary gaps
- **Compliance enrichment** — HIPAA, SOC2, ISO27001, PCI DSS, GDPR, NIST mappings

### 4. Threat Modeling Intelligence Engine (78%)
- **12 threat categories** — Injection, authentication, data exposure, DoS, cryptography, etc.
- **4 rule engines** — Components (20+ types), relationships (4 protocols), assumptions (5 keywords), trust boundaries (11 crossing types)
- **STRIDE correlation** — 6 STRIDE categories mapped to every threat
- **Severity engine** — Likelihood × impact with category-based adjustments
- **Threat clustering** — Groups threats by category with aggregated risk scores
- **Control recommendations** — Preventive, detective, corrective controls per threat

### 5. Attack Path Discovery Engine (82%)
- **Entry point discovery** — 5 rule groups (internet, third-party, API gateway, VPN, admin) with exposure scoring
- **Target asset identification** — 4 sensitivity levels (critical/high/medium/low) via keyword detection
- **Path construction** — DFS traversal from entry points to crown jewels with cycle detection
- **Threat chaining** — Links isolated threats into connected attacker journeys
- **Trust boundary traversal** — Every boundary crossing generates attack opportunities
- **Risk scoring** — Likelihood × impact with boundary/threat count adjustments
- **Business impact** — Maps technical risk to business narratives (HIPAA, PCI, SSO)
- **Detection difficulty** — 4 levels (Easy/Moderate/Hard/Very Hard) based on assumptions
- **Kill chain mapping** — 12-phase coverage (Reconnaissance through Impact)
- **MITRE ATT&CK mapping** — 30+ deterministic technique mappings

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
| PDF documents | `.pdf` | Text extraction via pdftxt library |
| Word documents | `.docx` | Text extraction via docx XML parsing |

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

### macOS / Linux — Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash
```

No repository clone required. The script detects your OS and architecture, downloads the correct binary from the latest GitHub release, verifies the SHA-256 checksum, and installs to `~/.local/bin/asf` with a symlink at `~/.asf/asf`.

Flags:
- **`--upgrade`, `-u`** — Upgrade existing installation (backs up config)
- **`--repair`** — Fix broken symlink/install without re-downloading
- **`--clean`** — Force clean reinstall (removes binary, keeps config)

```bash
# Upgrade to latest
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash -s -- --upgrade

# Repair a broken installation
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash -s -- --repair

# Clean reinstall
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash -s -- --clean
```

> **Private repository:** If the repository is private, set `GITHUB_TOKEN` before running the installer:
> ```bash
> export GITHUB_TOKEN=ghp_xxxxx
> curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash
> ```
> The script also auto-detects tokens from `gh auth token` if the GitHub CLI is installed.

### Windows — PowerShell

```powershell
powershell -c "irm https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.ps1 | iex"
```

### Manual Install

Download the binary for your platform from the [latest release](https://github.com/moksh5936-2/asfassumption/releases/latest):

| Platform | Download |
|----------|---------|
| macOS Apple Silicon | `ASF-v2.2.0-darwin-arm64` |
| macOS Intel | `ASF-v2.2.0-darwin-amd64` |
| Linux AMD64 | `ASF-v2.2.0-linux-amd64` |
| Linux ARM64 | `ASF-v2.2.0-linux-arm64` |
| Windows AMD64 | `ASF-v2.2.0-windows-amd64.exe` |

```bash
# Example: macOS Apple Silicon
curl -sfLO https://github.com/moksh5936-2/asfassumption/releases/download/v2.2.0/ASF-v2.2.0-darwin-arm64
chmod +x ASF-v2.2.0-darwin-arm64
mkdir -p ~/.local/bin ~/.asf
cp ASF-v2.2.0-darwin-arm64 ~/.asf/asf
ln -sf ~/.asf/asf ~/.local/bin/asf
```

### Verify Installation

```bash
asf --version
# Expected: ASF v2.2.0
```

### Uninstall

```bash
# Remove the binary and symlink
rm -f ~/.asf/asf ~/.local/bin/asf

# Remove configuration and data (optional)
rm -rf ~/.asf
```

### Prerequisites

| Component | Required | Install |
|-----------|----------|---------|
| Python 3.8+ | ❌ Removed | No Python required — native Go engine |
| Ollama | Optional (AI) | `brew install ollama` / https://ollama.ai |
| Tesseract | Optional (OCR) | `brew install tesseract` |

### Troubleshooting

| Problem | Solution |
|---------|----------|
| `asf: command not found` | Ensure `~/.local/bin` is in your PATH: `export PATH="$PATH:$HOME/.local/bin"`. Or run `~/.asf/asf` directly |
| `Permission denied` | Run `chmod +x ~/.asf/asf` |
| Download fails with 404 | The version may not have a release for your platform, or the repo is private. Try setting `GITHUB_TOKEN`: `export GITHUB_TOKEN=ghp_xxx && curl ... \| bash` |
| Checksum mismatch | The download may be corrupted. Re-run the installer or download manually from the GitHub release page |

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
3. Parses the response for additional assumptions and recommendations
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

ASF supports demo licensing with HMAC and Ed25519 signing. Both are demo-grade — keys are derived from compile-time constants and are extractable from the binary. Not suitable for production security.

```bash
# Activate an HMAC license
echo 'ASF-XXXX-XXXX-XXXX-XXXX' > ~/.asf/license.key

# Activate an Ed25519 license
echo 'ASF-ED25519-identifier-<hex_signature>' > ~/.asf/license.key

# Check license
asf --license
```

Legacy format: `ASF-XXXX-XXXX-XXXX-XXXX` (12 hex chars + 8-char HMAC signature)  
Ed25519 format: `ASF-ED25519-<identifier>-<128-char-hex-signature>` (Ed25519 signed)

Use `ReplacePublicKey()` in production deployments to swap the compile-time public key with a securely generated one.

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

- Cloud AI provider integration (optional)
- REST API server mode
- VS Code extension

---

## Validation Status

| Metric | Status |
|--------|--------|
| Unit tests | ✅ 257 tests across 11 packages (all pass) |
| Code quality | ✅ `go vet` clean |
| Benchmark run | ✅ 2158+ assumptions processed across 25+ architecture benchmarks |
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
  engine.go                # Analysis engine
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
  *_test.go               # 400+ tests across 13 packages
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
asf/                        # Python ASF engine (v1, archived)
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
A: Draw.io (.drawio), Mermaid (.mmd), YAML (.yaml/.yml), JSON (.json), SVG (.svg), images (.png/.jpg/.jpeg via OCR), and text documents (.txt/.md/.pdf/.docx). All formats have structured text extraction.

**Q: How is this different from a vulnerability scanner?**  
A: ASF finds implicit security assumptions, not known vulnerabilities. It answers "what did my team assume about the system?" not "what CVEs exist?"

**Q: Can I use ASF without Python?**  
A: Yes. ASF v2.1.1+ is a self-contained Go binary with no Python dependency. The analysis engine is fully native.

**Q: How accurate is ASF?**  
A: We don't know yet. Precision, recall, and false positive rate have not been measured. See [VALIDATION_STATUS.md](docs/VALIDATION_STATUS.md).

---

## Security Notice

ASF processes architecture documents entirely locally. No data is sent to external services. AI enhancement uses a local Ollama instance — no data leaves your machine.

License validation uses HMAC + Ed25519 and is performed locally. No phone-home mechanism exists. Demo keys are extractable from the binary (compile-time constants).

## Limitations

- Image OCR requires external Tesseract installation
- AI features require external Ollama installation
- Windows TUI not thoroughly tested
- Precision, recall, and false positive rate not measured
- No external validation study completed
- License system is demo-grade (HMAC + Ed25519 keys are derived from compile-time constants, extractable from binary)
- No CI/CD pipeline — releases must be built manually
- No code signing or macOS notarization

---

## Acknowledgements

- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and [Lipgloss](https://github.com/charmbracelet/lipgloss)
- PDF generation via [go-pdf/fpdf](https://github.com/go-pdf/fpdf)
- STRIDE methodology based on Microsoft's threat modeling framework

---

## License

Demo use only. Not for production or commercial use.

```
© 2026 ASF Project
```
