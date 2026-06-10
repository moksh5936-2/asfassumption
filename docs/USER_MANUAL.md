# ASF User Manual

## Table of Contents

1. [Installation](#installation)
2. [Startup](#startup)
3. [Navigation](#navigation)
4. [Main Menu](#main-menu)
5. [Analyze Architecture](#analyze-architecture)
6. [Results](#results)
7. [Review Mode](#review-mode)
8. [Settings](#settings)
9. [AI Settings](#ai-settings)
10. [About](#about)
11. [Exporting Reports](#exporting-reports)
12. [Configuration](#configuration)
13. [License Management](#license-management)
14. [Troubleshooting](#troubleshooting)
15. [Example Workflows](#example-workflows)

---

## Installation

### macOS

```bash
# Download the binary
curl -sfL https://github.com/moksh5936-2/asfassumption/releases/download/v1.0.0/asf-darwin-arm64 -o /usr/local/bin/asf

# Make executable
chmod +x /usr/local/bin/asf

# Run
asf
```

Or use the install script:

```bash
curl -sfL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/asf-tui/install.sh | bash
```

### Linux

```bash
curl -sfL https://github.com/moksh5936-2/asfassumption/releases/download/v1.0.0/asf-linux-amd64 -o /usr/local/bin/asf
chmod +x /usr/local/bin/asf
asf
```

### Windows

Download `asf-windows-amd64.exe` from the releases page.

### Prerequisites

| Feature | Requirement |
|---------|-------------|
| Core analysis | Python 3.8+ with ASF package (`pip install -e .`) |
| AI Enhancement | Ollama (https://ollama.ai) |
| Image OCR | Tesseract (`brew install tesseract`) |
| Full functionality | All of the above |

---

## Startup

Run the ASF binary:

```bash
asf
```

You will see the startup screen with the ASF fox logo and main menu.

**CLI Flags:**

```bash
asf --version     # Show version
asf -v            # Show version
asf --license     # Show license status
```

---

## Navigation

ASF is entirely keyboard-navigated. No commands to memorize.

| Key | Action |
|-----|--------|
| `↑` / `k` | Move selection up |
| `↓` / `j` | Move selection down |
| `Enter` | Select / Confirm |
| `Esc` | Back to previous view |
| `q` | Quit (from startup) |
| `Ctrl+C` | Force quit |

All available keys are shown in the help bar at the bottom of each screen.

---

## Main Menu

The startup screen presents 6 options:

```
▸ Analyze Architecture
  Results
  AI Settings
  Settings
  About
  Exit
```

| Option | Description |
|--------|-------------|
| **Analyze Architecture** | Run a security analysis on an architecture document |
| **Results** | View results from the most recent analysis |
| **AI Settings** | Manage local AI models (Ollama) |
| **Settings** | Configure ASF preferences |
| **About** | Version, license, and system information |
| **Exit** | Quit ASF |

After entering a view, you can access the Dashboard (press Esc from most views) with Quick Actions:

```
▸ Analyze Architecture
  Local AI Models
  Settings
  About
```

The Dashboard also shows system status (version, mode, AI status, theme).

---

## Analyze Architecture

Select "Analyze Architecture" to run a security analysis.

**Step 1: Set Document Path**

Select "Document Path" and press Enter. Type the full path to your architecture file, then press Enter again.

Supported formats:
- `.drawio` — Draw.io diagrams
- `.mmd` — Mermaid diagrams
- `.yaml` / `.yml` — Structured YAML definitions
- `.json` — Structured JSON definitions
- `.svg` — SVG diagrams
- `.png` / `.jpg` / `.jpeg` — Images (requires Tesseract)
- `.txt` / `.md` / `.pdf` / `.docx` — Text documents

**Step 2: Set Evidence Path (Optional)**

Select "Evidence Path", press Enter, type the path to your evidence CSV/JSON file. Evidence is used to verify assumptions.

**Step 3: Select Analysis Mode**

Two modes:

| Mode | Description |
|------|-------------|
| **ASF Engine Only** | Deterministic analysis only. No AI. Fully reproducible. |
| **ASF Engine + Local AI** | ASF analysis + local Ollama AI enhancement. Requires AI setup. |

**Step 4: Start Analysis**

Select "▶ Start Analysis" and press Enter.

A progress bar shows the analysis stages:

1. **Parsing Architecture** — Reading and parsing your input file
2. **Running ASF Engine** — Calling the Python ASF CLI for assumption extraction
3. **Processing Results** — Building the analysis result
4. **Generating STRIDE Mapping** — Applying STRIDE rules
5. **(Optional) Running AI Enhancement** — If AI mode is selected
6. **Complete** — Analysis finished

When complete, you are automatically taken to the Results view.

---

## Results

The Results view shows your analysis output with collapsible sections.

**Sections:**

| Section | Content |
|---------|---------|
| **Assumptions** | Full list of discovered assumptions with risk level, confidence, STRIDE |
| **Critical Assumptions** | Only Critical-rated assumptions |
| **Risk Matrix** | 5×5 visual matrix + risk level distribution |
| **STRIDE Distribution** | Bar chart of STRIDE category distribution |
| **Recommended Controls** | Mitigation controls mapped to assumptions and STRIDE |
| **Compliance** | Compliance findings |

**Keyboard shortcuts:**

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate sections |
| `Enter` | Expand/collapse section |
| `e` | Export results (default format) |
| `r` | Enter Review Mode |
| `Esc` | Back to Dashboard |

**Assumption Display:**

Each assumption shows:
- Risk level (colored: Critical=red, High=yellow, Medium=blue, Low=green)
- ID and description
- Component reference
- Confidence percentage
- Rationale (when expanded)

---

## Review Mode

Press `r` from Results view to enter Review Mode.

Review Mode allows architects to evaluate each assumption and record their judgment.

**Browse View:**

Shows a list of all assumptions with status markers:
- `?` — Not reviewed (Proposed)
- `✓` — Accepted
- `✗` — Rejected
- `~` — Modified

**Detail View:**

Press Enter on an assumption to see full detail:
- Description
- Risk level with likelihood/impact breakdown
- STRIDE categories
- Current review status
- Review notes
- Evidence sources
- Rationale
- Risk justification

**Review Keyboard Shortcuts:**

| Key | Action |
|-----|--------|
| `s` | Accept assumption |
| `r` | Reject assumption |
| `m` | Mark as Modified |
| `n` | Toggle note |
| `Enter` | Toggle between browse/detail view |
| `↑` / `↓` | Navigate assumptions |
| `v` | Toggle validation data view |

---

## Settings

Select "Settings" to configure ASF.

**Available Settings:**

| Setting | Values | Default | Description |
|---------|--------|---------|-------------|
| Theme | Dark, Midnight, Cyber, Minimal | Dark | UI color theme |
| Fox Style | Classic, Minimal, None | Classic | Startup fox art style |
| Analysis Depth | light, standard, deep | deep | Analysis depth |
| Risk Threshold | low, medium, high, critical | low | Minimum risk to highlight |
| STRIDE Analysis | true, false | true | Enable STRIDE mapping |
| Controls Check | true, false | true | Generate control recommendations |
| Default Export | json, markdown, html, csv, pdf | markdown | Default export format |
| Export Directory | (text) | ./reports | Where exports are saved |
| AI Enhancement | false, true | false | Enable AI enhancement |

**Keyboard shortcuts:**

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate settings |
| `Enter` | Start/confirm editing a value |
| `←` / `→` | Change value (when editing) |
| `s` | Save all settings |
| `Esc` | Back to Dashboard |

Changes apply immediately after editing. Press `s` to persist to disk.

---

## AI Settings

Select "AI Settings" to manage local AI models (Ollama required).

**Supported Models:**

| Model | Size |
|-------|------|
| Llama 3.2 (3B) | 2.0 GB |
| Qwen 3 (4B) | 2.5 GB |
| Mistral (7B) | 4.1 GB |
| Phi-4 (14B) | 9.1 GB |

**Actions per model:**

| Action | Description |
|--------|-------------|
| Download | Pull model from Ollama registry (requires internet) |
| Set as Active | Use this model for AI-enhanced analysis |
| Delete | Remove downloaded model |

**Keyboard shortcuts:**

| Key | Action |
|-----|--------|
| `↑` / `↓` | Navigate models/actions |
| `Enter` | Select model / execute action |
| `Esc` | Back (from actions submenu) |

**Note:** AI features are entirely optional. ASF's core analysis works without any AI model installed.

---

## About

Displays:
- Version number
- License status
- Description of ASF
- Technology stack
- Keyboard reference

---

## Exporting Reports

There are two ways to export:

### From Results View

Press `e` to export using the default format (configured in Settings).

### From Export View

Accessible from the Dashboard, the Export view lets you choose:

| Format | Extension | Use Case |
|--------|-----------|----------|
| JSON | `.json` | Machine processing, integration |
| Markdown | `.md` | Readable reports, documentation |
| CSV | `.csv` | Spreadsheets, data analysis |
| PDF | `.pdf` | Formal reports, presentations |
| HTML | `.html` | Browser viewing, dashboards |

**Export Contents:**

All export formats include:
- Architecture name and analysis date
- Analysis mode
- Total assumptions count
- Evidence summary (sources, components, relationships)
- Risk distribution (Critical, High, Medium, Low)
- STRIDE distribution
- Each assumption with:
  - Risk level and score
  - STRIDE categories
  - Confidence percentage
  - Rationale
  - Evidence sources
  - Risk justification (likelihood/impact breakdown)
  - STRIDE justification per category
  - Review status
- Recommended controls with mitigated assumption IDs and STRIDE
- Compliance findings

---

## Configuration

ASF reads configuration from `~/.asf/config.yaml`. It automatically migrates from the legacy path `~/.config/asf/config.yaml` if found.

**Example config:**

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

## License Management

ASF supports commercial licensing with HMAC-signed keys.

**License format:** `ASF-XXXX-XXXX-XXXX-XXXX`

**To activate a license:**

```bash
echo 'ASF-XXXX-XXXX-XXXX-XXXX' > ~/.asf/license.key
```

**To check license status:**

```bash
asf --license
```

Or view from the About screen in the TUI.

**License validation:** Offline. No internet required for license check.

**Tiers:**
- **Community** — No license file, full functionality
- **Enterprise** — Valid license key with HMAC signature

---

## Troubleshooting

### Common Issues

| Problem | Likely Cause | Solution |
|---------|-------------|----------|
| "ASF Engine error" | Python ASF CLI not installed | `pip install -e .` in project directory |
| "Tesseract not found" | OCR not installed | `brew install tesseract` (macOS) |
| "Ollama not found" | AI runtime not installed | Install from https://ollama.ai |
| "No results" | Analysis produced no assumptions | Check input file format and content |
| Config not saving | Permissions on `~/.asf/` | Check directory permissions |
| Wrong Python version | `.venv` path mismatch | Update `engine.go:NewEngine()` or use system Python |

### Known Limitations

- Python ASF CLI must be installed separately
- Image OCR quality depends on image clarity
- AI enhancement requires Ollama running locally
- Windows TUI not thoroughly tested (but should work with Windows Terminal)

---

## Example Workflows

### Basic Workflow

```
1. asf
2. Select "Analyze Architecture"
3. Document Path: ~/project/architecture.drawio
4. Select "ASF Engine Only"
5. Select "▶ Start Analysis"
6. Review results
7. Press 'e' to export as Markdown
8. Press 'r' to start review
9. Use s/r/m to mark assumptions
10. Esc → Esc → q to quit
```

### AI-Enhanced Workflow

```
1. asf
2. Select "AI Settings"
3. Select "Llama 3.2 (3B)" → Enter → Download
4. After download, "Set as Active"
5. Esc → "Analyze Architecture"
6. Select "ASF Engine + Local AI"
7. Start Analysis
8. Review AI-enhanced results (AI-prefixed IDs)
```

### Document Analysis Workflow

```
1. Place architecture document in ~/architectures/
2. Place optional evidence in ~/evidence/
3. asf → Analyze Architecture
4. Document Path: ~/architectures/system.drawio
5. Evidence Path: ~/evidence/access_control.csv
6. Start Analysis
7. Review contradicted assumptions
8. Export comprehensive report as HTML
```

### Team Review Workflow

```
1. Architect runs analysis
2. Exports to CSV for spreadsheet review
3. Loads results in Review Mode
4. Each architect reviews subset of assumptions
5. Uses s/r/m to mark review decisions
6. Collects validation data for study
7. Review outcomes inform security roadmap
```
