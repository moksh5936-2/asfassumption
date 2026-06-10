# ASF Developer Guide

## Project Structure

```
cybersec/
├── asf-tui/                  # Go TUI application (main product)
│   ├── main.go               # CLI entry point, flags
│   ├── app.go                # TUI controller, view routing
│   ├── startup.go            # Welcome screen + menu
│   ├── dashboard.go          # System status + quick actions
│   ├── analyze.go            # Analysis setup + progress
│   ├── results.go            # Results display
│   ├── review.go             # Architect review mode
│   ├── export.go             # 5-format export engine
│   ├── settings.go           # Configuration editor
│   ├── about.go              # Version/license info
│   ├── localai.go            # AI model manager TUI
│   ├── engine.go             # Core: Python bridge, result builder
│   ├── parser.go             # Architecture format parsers
│   ├── stride.go             # STRIDE rule engine
│   ├── justify.go            # 6 explainability engines
│   ├── explain.go            # Data structures
│   ├── ai.go                 # AI enhancement pipeline
│   ├── model.go              # Ollama manager
│   ├── license.go            # HMAC license system
│   ├── config.go             # YAML config load/save
│   ├── styles.go             # 4 themes, lipgloss styles
│   ├── explain_test.go       # 20 unit tests
│   ├── go.mod / go.sum       # Go dependencies
│   └── install.sh            # Multi-platform installer
├── asf/                      # Python ASF engine (v1)
│   ├── cli/                  # CLI entry point
│   ├── api/                  # FastAPI server
│   ├── extraction/           # Claim extraction
│   ├── verification/         # Evidence verification
│   ├── assumption/           # Assumption models
│   ├── evidence/             # Evidence parsing
│   ├── gaps/                 # Gap analysis
│   ├── confidence/           # Confidence scoring
│   ├── graph/                # Relationship graph
│   ├── ingestion/            # Document parsing
│   ├── llm/                  # Optional LLM integration
│   ├── db/                   # SQLite persistence
│   └── models/               # Pydantic data models
├── benchmark/                # Research validation
│   ├── experiments/          # 20 architecture simulations
│   ├── report/               # Analyst guides, ontology
│   └── assumption_knowledge_base/  # Gold standard data
├── docs/                     # Documentation
├── sample_data/              # Sample inputs
├── tests/                    # Python tests
├── scripts/                  # Utility scripts
└── templates/                # Report templates
```

## Build Process

### Prerequisites

- Go 1.24+
- Python 3.8+ (for ASF engine)
- Optional: Ollama, Tesseract

### Building

```bash
cd asf-tui

# Build for current platform
go build -o asf-tui .

# Cross-compile
GOOS=linux GOARCH=amd64 go build -o asf-linux-amd64 .
GOOS=linux GOARCH=arm64 go build -o asf-linux-arm64 .
GOOS=darwin GOARCH=amd64 go build -o asf-darwin-amd64 .
GOOS=darwin GOARCH=arm64 go build -o asf-darwin-arm64 .
GOOS=windows GOARCH=amd64 go build -o asf-windows-amd64.exe .
```

### Testing

```bash
cd asf-tui
go test ./... -v    # Run all tests
go test -run TestRiskMatrix ./... -v   # Run specific test
```

### Code Quality

```bash
go vet ./...        # Static analysis
golangci-lint run   # Linting (if configured)
```

## Internal APIs

### Engine API

```go
// Create engine
engine := NewEngine(cfg)

// Run analysis
result, err := engine.RunAnalysis(archPath, evPath, mode, progress)

// Run analysis with progress tracking
progress := make(chan AnalysisProgress, 10)
go func() {
    for p := range progress {
        fmt.Printf("%.0f%%: %s\n", p.Percent, p.Stage)
    }
}()
result, err := engine.RunAnalysis(path, evPath, mode, progress)
```

### Parser API

```go
// Parse any supported format
desc, err := ParseArchitecture(path)
// Returns ArchDescription with Components, Relationships, RawText
```

### STRIDE Engine API

```go
se := NewStrideEngine()
categories := se.MapAssumption(category, text, keywords)
// Returns []StrideCategory matched by rules

// Access rules
rules := se.GetKeywordRules()    // 33 keyword rules
cats := se.GetCategoryRules()    // 17 category rules
```

### Explainability Pipeline API

```go
pipe := NewExplainabilityPipeline(archDesc, sourcePath, strideEngine)

// Process single assumption
pipe.Explain(&assumption)
// Populates: EvidenceSources, SourceComponents, SourceRelationships,
// Rationale, StrideJustifications, RiskJustification, Confidence

// Build evidence summary
summary := pipe.BuildEvidenceSummary(assumptions)
```

### Export API

```go
path, err := ExportResult(result, ExportMarkdown, "./reports")
// Formats: ExportJSON, ExportMarkdown, ExportCSV, ExportPDF, ExportHTML
```

### Config API

```go
cfg, err := LoadConfig(path)
cfg.Save(path)
path := ConfigPath()   // Returns ~/.asf/config.yaml
```

### License API

```go
info := ValidateLicense("ASF-XXXX-XXXX-XXXX-XXXX")
info := LoadLicense()           // Load from ~/.asf/license.key
SaveLicense("ASF-XXXX-...")     // Save to ~/.asf/license.key
GenerateLicenseKey(data)        // Generate new license
```

### AI Model API

```go
mm := NewModelManager()
mm.CheckAvailable()     // Is Ollama running?
mm.ListInstalled()      // List installed models
mm.StartDownload(name, progress) // Download model
mm.DeleteModel(name)    // Remove model
response, err := mm.Generate(prompt, model) // Generate text
```

### AI Enhancement API

```go
enhancer := NewAIEnhancer()
enhanced, err := enhancer.Enhance(result, "llama3.2:3b")
result = mergeAIResults(original, enhanced)  // Merge AI findings
```

### Review & Validation API

```go
records := CollectValidationData(assumptions)
// Returns []ValidationRecord for precision/recall studies
```

## Architecture

### Model-View-Update (Elm Architecture)

ASF-TUI uses Bubble Tea's Elm-style architecture:

```
Model ←── Update(msg) ──→ Cmd
  │                         │
  └────── View() ──────────→┘
```

Each view is a separate model struct with its own Update and View methods. The main controller (`app.go`) routes messages to the active view.

### Data Flow

1. User input → `tea.KeyMsg` → `mainModel.Update()`
2. Router dispatches to active view's Update method
3. View returns new model + optional command
4. Commands execute async (Python CLI, AI, etc.)
5. Results returned as messages → Update processes them
6. View renders updated model

### Thread Safety

- All model state is owned by the main event loop
- Async operations (Python CLI, AI, download) run in goroutines
- Results delivered via tea.Msg through channels
- `syncProgress` struct uses `sync.Mutex` for download progress

## Extension Points

### Adding a New Parser

1. Add file extension check in `ParseArchitecture()` (parser.go)
2. Implement parse function returning `*ArchDescription`
3. Add prose generation in `buildTextFromDiagram()` if needed

### Adding a New STRIDE Rule

1. Add category mapping in `buildCategoryRules()` (stride.go)
2. Or add keyword pattern in `buildKeywordRules()` (stride.go)

### Adding a New Export Format

1. Add format constant in export.go
2. Add case in `ExportResult()` switch
3. Implement export function following existing patterns

### Adding a New Theme

1. Add theme entry in `Themes` map (styles.go)
2. Add to settings options in `settings.go`

### Adding a New Setting

1. Add field to Config struct (config.go)
2. Add setting item in `newSettingsModel()` (settings.go)
3. Add case in `applyChange()` (settings.go)

## Release Process

```bash
# 1. Update version in license.go (ASFVersion constant)
# 2. Update version in install.sh
# 3. Update CHANGELOG.md
# 4. Build all platforms
# 5. Generate checksums
# 6. Tag release
# 7. Push tag
# 8. Create GitHub release with binaries
```

## Testing Strategy

- **Unit tests** in `explain_test.go` (20 tests)
  - Risk matrix calculation and boundary conditions
  - Confidence engine (determinism, max cap, empty input)
  - Assumption justification (component/relationship/fallback)
  - Evidence engine tracing (keyword matching, no-match)
  - STRIDE justification (category mapping, rule indices, determinism)
  - Likelihood and impact analyzers (range, exposure sensitivity)
  - Full pipeline integration
  - Edge cases (nil assumption, empty evidence, max confidence)
  - Validation data collection

- **Missing tests** (needs contribution):
  - Parser unit tests (Draw.io, Mermaid, YAML, JSON, SVG, OCR)
  - Export format tests (JSON, Markdown, CSV, PDF, HTML)
  - AI enhancement tests
  - Integration tests with real Python ASF CLI
  - TUI view tests
  - Config migration tests
  - License edge cases

## Dependencies

| Dependency | Version | Purpose |
|-----------|---------|---------|
| github.com/charmbracelet/bubbletea | v1.3.10 | TUI framework |
| github.com/charmbracelet/lipgloss | v1.1.0 | Terminal styling |
| github.com/go-pdf/fpdf | v0.9.0 | PDF generation |
| gopkg.in/yaml.v3 | v3.0.1 | YAML parsing |
| Python 3.8+ | — | ASF engine (external) |
| Ollama | — | Local AI (optional) |
| Tesseract | — | Image OCR (optional) |

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md).
