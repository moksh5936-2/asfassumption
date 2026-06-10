# ASF Technical Reference

## Struct Reference

### Config (`config.go`)

```go
type Config struct {
    General struct {
        Theme    string   // UI theme name
        FoxStyle string   // Fox art style
    }
    Analysis struct {
        Depth         string // "light" | "standard" | "deep"
        Stride        bool   // Enable STRIDE mapping
        Controls      bool   // Generate controls
        RiskThreshold string // "low" | "medium" | "high" | "critical"
    }
    AI struct {
        Enabled         bool     // Enable AI enhancement
        ActiveModel     string   // Currently selected Ollama model
        InstalledModels []string // Downloaded models
    }
    Output struct {
        Default   string // Default export format
        Directory string // Export output directory
    }
    Appearance struct {
        Theme    string // UI theme name
        FoxStyle string // Fox art style
    }
}
```

### Engine (`engine.go`)

```go
type Engine struct {
    config       *Config
    pythonPath   string   // Path to Python interpreter
    projectDir   string   // ASF project directory
    strideEngine *StrideEngine
    explainPipe  *ExplainabilityPipeline
    archDesc     *ArchDescription
}
```

### AnalysisResult (`engine.go`)

```go
type AnalysisResult struct {
    ArchitectureName   string
    AnalysisDate       time.Time
    AnalysisMode       string                 // ModeASFOnly | ModeASFAndAI
    Assumptions        []Assumption
    CriticalCount      int
    HighCount          int
    TotalAssumptions   int
    StrideDistribution map[StrideCategory]int
    Controls           []ControlDetail
    Compliance         []string
    Summary            string
    TrueAssumptions    int
    FalseAssumptions   int
    CriticalGaps       int
    EvidenceSummary    EvidenceSummary
    RiskModelVersion   string                 // "asf-risk-model-1.0"
    ConfidenceSummary  string
}
```

### Assumption (`engine.go`)

```go
type Assumption struct {
    ID                 string
    Description        string
    Component          string
    Category           string
    Risk               RiskLevel
    Stride             []StrideCategory
    Likelihood         int       // 1-5
    Impact             int       // 1-5
    Confidence         float64   // 0.0-1.0
    Keywords           []string
    EvidenceSources    []string
    SourceComponents   []string
    SourceRelationships []string
    Rationale          string
    StrideJustifications []StrideJustification
    RiskJustification  *RiskJustification
    ReviewStatus       string    // "Proposed"|"Accepted"|"Rejected"|"Modified"
    ReviewNotes        string
    ReviewTimestamp    time.Time
}
```

### StrideEngine (`stride.go`)

```go
type StrideEngine struct {
    categoryRules map[string][]StrideCategory  // 17 categories
    keywordRules  []keywordRule                // 33 patterns
}

type keywordRule struct {
    keywords []string
    stride   []StrideCategory
}
```

### Evidence Sources & Justification (`explain.go`)

```go
type EvidenceSource struct {
    FilePath                string
    FileType                string
    MatchedComponents       []string
    MatchedRelationships    []string
    MatchedTrustBoundaries  []string
    MatchedSecurityConcepts []string
}

type EvidenceSummary struct {
    TotalSources       int
    TotalComponents    int
    TotalRelationships int
    SourceFiles        []string
}

type StrideJustification struct {
    Category          StrideCategory
    Reason            string
    MatchedRuleIndexes []int       // Indexes into stride.go keyword rules
    MatchedKeywords   []string
    MatchedComponents []string
    Confidence        float64
    ConfidenceReason  string
}

type RiskJustification struct {
    Likelihood        int              // 1-5
    LikelihoodReason  string
    LikelihoodFactors []LikelihoodFactor
    Impact            int              // 1-5
    ImpactReason      string
    ImpactFactors     []ImpactFactor
    RiskScore         int              // 1-25
    RiskLevel         RiskLevel
    RiskReason        string
    Confidence        float64
    ConfidenceReason  string
}

type LikelihoodFactor struct {
    Factor string  // "Exposure Level" | "Authentication Dependency" | "Attack Surface Complexity"
    Value  int     // 1-5
    Reason string
}

type ImpactFactor struct {
    Factor string  // "Data Classification" | "Regulatory Exposure" | "Business Criticality"
    Value  int     // 1-5
    Reason string
}

type ControlDetail struct {
    ID                   string
    Description          string
    Rationale            string
    Category             string
    MitigatedAssumptionIDs []string
    MitigatedSTRIDE      []StrideCategory
    Priority             int  // 1=highest, 3=lowest
}

type ReviewRecord struct {
    Status    string    // Proposed, Accepted, Rejected, Modified
    Notes     string
    Timestamp time.Time
    Reviewer  string
}

type ValidationRecord struct {
    AssumptionID      string
    Description       string
    GeneratedEvidence  []string
    AssignedRisk      RiskLevel
    RiskScore         int
    Confidence        float64
    STRIDECategories  []StrideCategory
    ArchReviewResult  string   // Accepted, Rejected, Modified
    ArchNotes         string
    ReviewTimestamp   time.Time
}
```

### Parser (`parser.go`)

```go
type Component struct {
    ID    string
    Label string
}

type Relation struct {
    Source string
    Target string
    Label  string
}

type ArchDescription struct {
    Name          string
    Components    []Component
    Relationships []Relation
    Policies      []string
    RawText       string   // Generated prose for ASF engine
}
```

### Evidence Engine (`justify.go`)

```go
type EvidenceEngine struct {
    arch       *ArchDescription
    sourcePath string
    sourceType string
}

type EvidenceResult struct {
    MatchedComponents       []string
    MatchedRelationships    []string
    MatchedTrustBoundaries  []string
    MatchedSecurityConcepts []string
    EvidenceCount           int
}
```

### Explainability Pipeline (`justify.go`)

```go
type ExplainabilityPipeline struct {
    evidenceEngine     *EvidenceEngine
    strideJustify      *StrideJustifyEngine
    likelihoodAnalyzer *LikelihoodAnalyzer
    impactAnalyzer     *ImpactAnalyzer
    riskMatrix         *RiskMatrix
    confidenceEngine   *ConfidenceEngine
}
```

### AI (`ai.go`, `model.go`)

```go
type AIEnhancer struct {
    model *ModelManager
}

type AIEnhancedResult struct {
    AdditionalAssumptions []AIAssumption
    RefinedRisks         []AIRiskRefinement
    MissingThreats       []string
    Recommendations      []string
    RawResponse          string
}

type ModelManager struct {
    ollamaCmd string
}

type OllamaGenerateRequest struct {
    Model  string
    Prompt string
    Stream bool
}

type OllamaGenerateResponse struct {
    Response string
    Done     bool
}
```

### License (`license.go`)

```go
type LicenseInfo struct {
    Key     string
    Valid   bool
    Tier    string
    Message string
}
```

## Interfaces & Constants

### Risk Levels

```go
const (
    RiskCritical RiskLevel = "Critical"
    RiskHigh     RiskLevel = "High"
    RiskMedium   RiskLevel = "Medium"
    RiskLow      RiskLevel = "Low"
)
```

### STRIDE Categories

```go
const (
    StrideSpoofing        StrideCategory = "Spoofing"
    StrideTampering       StrideCategory = "Tampering"
    StrideRepudiation     StrideCategory = "Repudiation"
    StrideInfoDisclosure  StrideCategory = "Information Disclosure"
    StrideDenialOfService StrideCategory = "Denial of Service"
    StrideElevationPriv   StrideCategory = "Elevation of Privilege"
)
```

### Analysis Modes

```go
const (
    ModeASFOnly  = "ASF Engine Only"
    ModeASFAndAI = "ASF Engine + Local AI"
)
```

### Export Formats

```go
const (
    ExportJSON     ExportFormat = "json"
    ExportMarkdown ExportFormat = "markdown"
    ExportCSV      ExportFormat = "csv"
    ExportPDF      ExportFormat = "pdf"
    ExportHTML     ExportFormat = "html"
)
```

### License Constants

```go
const (
    ASFVersion    = "1.0.0"
    LicensePrefix = "ASF"
)
```

## Services

### Engine Services

| Service | Method | Input | Output | Description |
|---------|--------|-------|--------|-------------|
| NewEngine | Constructor | *Config | *Engine | Create engine with config |
| RunAnalysis | Method | archPath, evPath, mode, progress | *AnalysisResult, error | Full analysis pipeline |
| callPythonCLI | Method | docPath, evPath | *asfJSONResult, error | Bridge to Python ASF |
| buildResult | Method | *asfJSONResult, archPath, mode | *AnalysisResult | Build structured result |
| mapStrideDistribution | Method | []Assumption | map[StrideCategory]int | Count STRIDE categories |
| generateControls | Function | []Assumption | []ControlDetail | Map controls to assumptions |

### Parser Services

| Service | Method | Input | Output |
|---------|--------|-------|--------|
| ParseArchitecture | Function | path | *ArchDescription, error |
| parseDrawio | Function | path | *ArchDescription, error |
| parseMermaid | Function | path | *ArchDescription, error |
| parseYAMLArch | Function | path | *ArchDescription, error |
| parseJSONArch | Function | path | *ArchDescription, error |
| parseSVG | Function | path | *ArchDescription, error |
| parseImageOCR | Function | path | *ArchDescription, error |
| buildTextFromDiagram | Function | name, components, relations | string |

### STRIDE Services

| Service | Method | Input | Output |
|---------|--------|-------|--------|
| NewStrideEngine | Constructor | — | *StrideEngine |
| MapAssumption | Method | category, text, keywords | []StrideCategory |
| GetKeywordRules | Method | — | []keywordRule |
| GetCategoryRules | Method | — | map[string][]StrideCategory |

### Explainability Services

| Service | Method | Input | Output |
|---------|--------|-------|--------|
| NewExplainabilityPipeline | Constructor | arch, sourcePath, strideEngine | *ExplainabilityPipeline |
| Explain | Method | *Assumption | void (mutates assumption) |
| BuildEvidenceSummary | Method | []Assumption | EvidenceSummary |
| NewEvidenceEngine | Constructor | arch, sourcePath | *EvidenceEngine |
| TraceEvidence | Method | category, keywords, text | *EvidenceResult |
| BuildEvidenceSources | Method | *EvidenceResult | []string |
| JustifyAssumption | Function | category, evidence | string |
| AnalyzeLikelihood | Method | *Assumption, *EvidenceResult | int, string, []LikelihoodFactor |
| AnalyzeImpact | Method | *Assumption, *EvidenceResult | int, string, []ImpactFactor |
| RiskMatrix.Calculate | Method | likelihood, impact | int, RiskLevel |
| RiskMatrix.RiskReason | Method | lh, im, score, level | string |
| ConfidenceEngine.Calculate | Method | evCount, stCount, compCount, relCount | float64, string |

### Export Services

| Service | Method | Input | Output |
|---------|--------|-------|--------|
| ExportResult | Function | result, format, outputDir | string, error |
| exportJSON | Function | result, path | error |
| exportMarkdown | Function | result, path | error |
| exportCSV | Function | result, path | error |
| exportPDF | Function | result, path | error |
| exportHTML | Function | result, path | error |

### Config Services

| Service | Method | Input | Output |
|---------|--------|-------|--------|
| LoadConfig | Function | path | *Config, error |
| Save | Method | path | error |
| ConfigPath | Function | — | string |
| DefaultConfig | Function | — | Config |

### License Services

| Service | Method | Input | Output |
|---------|--------|-------|--------|
| LoadLicense | Function | — | *LicenseInfo |
| ValidateLicense | Function | key | *LicenseInfo |
| SaveLicense | Function | key | error |
| GenerateLicenseKey | Function | data | string |

### AI Services

| Service | Method | Input | Output |
|---------|--------|-------|--------|
| NewAIEnhancer | Constructor | — | *AIEnhancer |
| Enhance | Method | result, modelName | *AIEnhancedResult, error |
| buildPrompt | Method | *AnalysisResult | string |
| parseResponse | Method | response, result | *AIEnhancedResult |

### Model Services

| Service | Method | Input | Output |
|---------|--------|-------|--------|
| NewModelManager | Constructor | — | *ModelManager |
| CheckAvailable | Method | — | bool |
| ListInstalled | Method | — | []string, error |
| StartDownload | Method | modelName, *syncProgress | void (async) |
| DeleteModel | Method | modelName | error |
| Generate | Method | prompt, model | string, error |

### Review Services

| Service | Method | Input | Output |
|---------|--------|-------|--------|
| CollectValidationData | Function | []Assumption | []ValidationRecord |

## YAML Configuration

```yaml
# ~/.asf/config.yaml
general:
  theme: Dark              # "Dark" | "Midnight" | "Cyber" | "Minimal"
  fox_style: Classic       # "Classic" | "Minimal" | "None"
analysis:
  depth: deep              # "light" | "standard" | "deep"
  stride: true             # true | false
  controls: true           # true | false
  risk_threshold: low      # "low" | "medium" | "high" | "critical"
ai:
  enabled: false           # true | false
  active_model: ""         # Ollama model name
  installed_models: []     # List of installed models
output:
  default: markdown        # "json" | "markdown" | "html" | "csv" | "pdf"
  directory: ./reports     # Export output directory
appearance:
  theme: Dark              # "Dark" | "Midnight" | "Cyber" | "Minimal"
  fox_style: Classic       # "Classic" | "Minimal" | "None"
```

## CLI Flags

| Flag | Description |
|------|-------------|
| `--version`, `-v` | Show version and exit |
| `--license` | Show license status and exit |
| (no args) | Launch TUI |

## Export Formats

### JSON

Full `AnalysisResult` struct serialized with `json.MarshalIndent`. Contains all fields including evidence, reasoning, confidence, controls, review status.

### Markdown

Structured markdown with sections for summary, evidence, risk distribution, STRIDE distribution, detailed assumptions, controls, compliance.

### CSV

Flat CSV with columns: ID, Description, Component, Category, Risk, STRIDE, Likelihood, Impact, RiskScore, Confidence, ReviewStatus, EvidenceSources, Rationale, MitigatingControls.

### PDF

Multi-page PDF with title page, summary, risk distribution, STRIDE distribution, recommended controls, detailed assumptions. Uses Helvetica font.

### HTML

Styled single-page HTML with CSS for risk badges, STRIDE bar chart, expandable details, responsive layout.

## Parser Support

| Format | Extension | Parser | Engine Used |
|--------|-----------|--------|-------------|
| Draw.io | .drawio | XML/gzip | Go encoding/xml |
| Mermaid | .mmd | Regex | Go regexp |
| YAML | .yaml/.yml | Structured | gopkg.in/yaml.v3 |
| JSON | .json | Structured | encoding/json |
| SVG | .svg | XML | encoding/xml |
| PNG/JPG | .png/.jpg/.jpeg | OCR | Tesseract CLI |
| Text | .txt | Raw text | os.ReadFile |
| Markdown | .md | Raw text | os.ReadFile |
| PDF | .pdf | Raw text | os.ReadFile |
| DOCX | .docx | Raw text | os.ReadFile |

## Dependencies

### Go (direct)

| Module | Version | Purpose |
|--------|---------|---------|
| github.com/charmbracelet/bubbletea | v1.3.10 | TUI framework |
| github.com/charmbracelet/lipgloss | v1.1.0 | Terminal styling |
| github.com/go-pdf/fpdf | v0.9.0 | PDF generation |
| gopkg.in/yaml.v3 | v3.0.1 | YAML parsing |

### External (runtime)

| Dependency | Required For | Install |
|-----------|-------------|---------|
| Python 3.8+ | ASF engine | `pip install -e .` |
| Ollama | AI enhancement | https://ollama.ai |
| Tesseract | Image OCR | `brew install tesseract` |

### Go (indirect)

| Module | Version |
|--------|---------|
| github.com/aymanbagabas/go-osc52/v2 | v2.0.1 |
| github.com/charmbracelet/colorprofile | v0.2.3 |
| github.com/charmbracelet/x/ansi | v0.10.1 |
| github.com/charmbracelet/x/cellbuf | v0.0.13 |
| github.com/charmbracelet/x/term | v0.2.1 |
| github.com/erikgeiser/coninput | v0.0.0-20211004153227 |
| github.com/lucasb-eyer/go-colorful | v1.2.0 |
| github.com/mattn/go-isatty | v0.0.20 |
| github.com/mattn/go-localereader | v0.0.1 |
| github.com/mattn/go-runewidth | v0.0.16 |
| github.com/muesli/ansi | v0.0.0-20230316100256 |
| github.com/muesli/cancelreader | v0.2.2 |
| github.com/muesli/termenv | v0.16.0 |
| github.com/rivo/uniseg | v0.4.7 |
| github.com/xo/terminfo | v0.0.0-20220910002029 |
| golang.org/x/sys | v0.36.0 |
| golang.org/x/text | v0.3.8 |

## Environment

There are no environment variables. All configuration is via `~/.asf/config.yaml`.
