# ASF Architecture

## Overview

ASF (Architecture Security Framework) is a terminal-based security analysis tool that automatically discovers hidden security assumptions in system architecture diagrams and documents. It combines a deterministic rule engine (STRIDE threat modeling, risk assessment, confidence scoring) with optional local AI enhancement.

## System Architecture

```
┌────────────────────────────────────────────────────────────┐
│                     ASF TUI (Bubble Tea)                    │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────┐  │
│  │ Startup  │ │Dashboard │ │ Analyze  │ │   Results    │  │
│  │ View     │ │ View     │ │ View     │ │   View       │  │
│  └──────────┘ └──────────┘ └──────────┘ └──────────────┘  │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────────┐  │
│  │ Local AI │ │Settings  │ │ About    │ │   Export     │  │
│  │ View     │ │ View     │ │ View     │ │   View       │  │
│  └──────────┘ └──────────┘ └──────────┘ └──────────────┘  │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                  Review View                          │  │
│  └──────────────────────────────────────────────────────┘  │
└───────────────────────┬────────────────────────────────────┘
                        │
                        ▼
┌────────────────────────────────────────────────────────────┐
│                    Main Controller (app.go)                 │
│  ┌─────────────┐ ┌──────────────┐ ┌────────────────────┐  │
│  │ View Router  │ │ View History │ │ Message Dispatcher │  │
│  └─────────────┘ └──────────────┘ └────────────────────┘  │
└───────────────────────┬────────────────────────────────────┘
                        │
                        ▼
┌────────────────────────────────────────────────────────────┐
│                     Engine Layer                            │
│  ┌──────────────────────────────────────────────────────┐  │
│  │              Engine (engine.go)                       │  │
│  │  ┌─────────────┐  ┌──────────────┐  ┌────────────┐  │  │
│  │  │ Architecture │  │  Python ASF  │  │  Result    │  │  │
│  │  │  Parsing     │─▶│  CLI Bridge  │─▶│  Builder   │  │  │
│  │  └─────────────┘  └──────────────┘  └────────────┘  │  │
│  └──────────────────────────────────────────────────────┘  │
└───────────────────────┬────────────────────────────────────┘
                        │
        ┌───────────────┼───────────────┐
        ▼               ▼               ▼
┌──────────────┐ ┌───────────┐ ┌────────────────┐
│  Parser      │ │  STRIDE   │ │  Explainability│
│  (parser.go) │ │  Engine   │ │  Pipeline      │
│              │ │  (stride) │ │  (justify.go)  │
│  • Draw.io   │ │  • 17 cat │ │  • Evidence    │
│  • Mermaid   │ │  • 33 kw  │ │  • STRIDE Just │
│  • YAML/JSON │ │  • Detrm  │ │  • Likelihood  │
│  • SVG       │ │           │ │  • Impact      │
│  • OCR       │ │           │ │  • Risk Matrix │
│  • TXT/PDF/  │ │           │ │  • Confidence  │
│    DOCX      │ │           │ │                │
└──────────────┘ └───────────┘ └────────────────┘
                        │
                        ▼
┌────────────────────────────────────────────────────────────┐
│                  Export Engine (export.go)                  │
│  ┌────────┐ ┌──────────┐ ┌─────┐ ┌──────┐ ┌────────┐    │
│  │  JSON  │ │ Markdown │ │ CSV │ │ PDF  │ │  HTML  │    │
│  └────────┘ └──────────┘ └─────┘ └──────┘ └────────┘    │
└────────────────────────────────────────────────────────────┘
```

## Component Diagram

```
┌─────────────────────────────────────────────────────────────┐
│                       Engine                                 │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────────┐ │
│  │ Python   │  │ STRIDE   │  │ Explain  │  │ Model Mgr  │ │
│  │ CLI      │  │ Engine   │  │ Pipeline │  │ (model.go) │ │
│  │ Bridge   │  │(stride)  │  │(justify) │  │            │ │
│  └──────────┘  └──────────┘  └──────────┘  └────────────┘ │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────────┐ │
│  │  Parser  │  │  License │  │  AI      │  │  Config    │ │
│  │(parser)  │  │(license) │  │(ai.go)   │  │ (config)   │ │
│  └──────────┘  └──────────┘  └──────────┘  └────────────┘ │
└─────────────────────────────────────────────────────────────┘
```

## Data Flow

```
┌──────────┐     ┌──────────┐     ┌──────────┐     ┌──────────┐
│  User    │────▶│  TUI     │────▶│  Engine  │────▶│  Python  │
│ Input    │     │  Views   │     │  Layer   │     │  ASF CLI │
└──────────┘     └──────────┘     └──────────┘     └──────────┘
                                     │                    │
                                     ▼                    ▼
                              ┌──────────┐        ┌──────────────┐
                              │  Parser  │        │  JSON Result  │
                              │  Layer   │        │  Parse       │
                              └──────────┘        └──────┬───────┘
                                                         │
                                                         ▼
                                              ┌──────────────────┐
                                              │  Explainability  │
                                              │  Pipeline        │
                                              │  6 Engines       │
                                              └────────┬─────────┘
                                                       │
                                                       ▼
                                              ┌──────────────────┐
                                              │  AnalysisResult  │
                                              │  + Assumptions   │
                                              └────────┬─────────┘
                                                       │
                                    ┌──────────────────┼────────────┐
                                    ▼                  ▼            ▼
                              ┌──────────┐      ┌──────────┐ ┌──────────┐
                              │  TUI     │      │  Export  │ │  Review  │
                              │  Results │      │  5 fmt   │ │  Mode    │
                              └──────────┘      └──────────┘ └──────────┘
```

## Processing Pipeline

```
Document Input
    │
    ▼
┌────────────────────────────────────────┐
│ 1. Parse Architecture (parser.go)       │
│    • Draw.io → XML → components/rels    │
│    • Mermaid → regex → components/rels  │
│    • YAML/JSON → structured def         │
│    • SVG → XML text extraction          │
│    • PNG/JPG → Tesseract OCR            │
│    • TXT/MD/PDF/DOCX → raw text         │
│    Output: ArchDescription{             │
│      Components, Relationships,         │
│      RawText (prose)                    │
│    }                                    │
└────────────────────────────────────────┘
    │
    ▼
┌────────────────────────────────────────┐
│ 2. Run Python ASF CLI (engine.go)       │
│    • exec: python -m asf.cli.main       │
│    • Input: architecture prose text     │
│    • Evidence: optional CSV/JSON        │
│    • Output: JSON with assumptions,     │
│      verifications, gaps                │
└────────────────────────────────────────┘
    │
    ▼
┌────────────────────────────────────────┐
│ 3. Build Result (engine.go)             │
│    • Parse ASF JSON                     │
│    • Map risk levels                    │
│    • Apply STRIDE engine rules          │
│    • Initialize explainability pipeline │
└────────────────────────────────────────┘
    │
    ▼
┌────────────────────────────────────────┐
│ 4. Explainability Pipeline (justify.go) │
│    Phase 1: EvidenceEngine              │
│      • Match components                 │
│      • Match relationships              │
│      • Detect trust boundaries          │
│      • Identify security concepts       │
│                                         │
│    Phase 2: JustifyAssumption           │
│      • Build human-readable rationale   │
│                                         │
│    Phase 3: StrideJustifyEngine         │
│      • Match STRIDE categories          │
│      • Track matched rules/keywords     │
│      • Calculate per-category confidence│
│                                         │
│    Phase 4: LikelihoodAnalyzer          │
│      • Exposure level (1-5)             │
│      • Auth dependency (1-5)            │
│      • Attack complexity (1-5)          │
│                                         │
│    Phase 5: ImpactAnalyzer              │
│      • Data classification (1-5)        │
│      • Regulatory exposure (1-5)        │
│      • Business criticality (1-5)       │
│                                         │
│    Phase 6: RiskMatrix                  │
│      • L × I = score (1-25)             │
│      • Map to risk level                │
│                                         │
│    Phase 7: ConfidenceEngine            │
│      • Evidence points                  │
│      • STRIDE rule matches              │
│      • Component/relationship matches   │
│      • Combined score (0.1-0.95)        │
└────────────────────────────────────────┘
    │
    ▼
┌────────────────────────────────────────┐
│ 5. Optional AI Enhancement (ai.go)      │
│    • Only if AI enabled + configured    │
│    • Build prompt from analysis results │
│    • Call Ollama API locally             │
│    • Parse AI response                  │
│    • Merge additional assumptions       │
│    • Add AI-generated controls          │
│    • Tag with AI- prefix on IDs         │
└────────────────────────────────────────┘
    │
    ▼
┌────────────────────────────────────────┐
│ 6. Generate Controls (engine.go)        │
│    • 16 control templates               │
│    • Map by assumption category/STRIDE  │
│    • Priority 1-3 ordering              │
│    • Link mitigated assumption IDs      │
└────────────────────────────────────────┘
    │
    ▼
┌────────────────────────────────────────┐
│ 7. Display / Export                     │
│    • TUI Results view                   │
│    • Review mode with status tracking   │
│    • Export to JSON/MD/CSV/PDF/HTML     │
│    • Collect validation data for studies│
└────────────────────────────────────────┘
```

## Module Descriptions

| Module | File | Lines | Purpose |
|--------|------|-------|---------|
| Main | `main.go` | 57 | CLI entry point, flags, config load |
| App | `app.go` | 266 | TUI controller, view routing, message dispatch |
| Startup | `startup.go` | 154 | Welcome screen with fox art + menu |
| Dashboard | `dashboard.go` | 106 | System status + quick actions |
| Analyze | `analyze.go` | 319 | Analysis setup, progress, file input |
| Results | `results.go` | 424 | Assumption list, risk matrix, STRIDE, controls |
| Review | `review.go` | 253 | Architect review mode with status tracking |
| Export | `export.go` | 672 | 5 export formats with full evidence/reasoning |
| Settings | `settings.go` | 238 | Config editor with cycle values |
| About | `about.go` | 72 | Version, license, technology info |
| Local AI | `localai.go` | 267 | Model list, download, delete, activate |
| Config | `config.go` | 101 | YAML config load/save with migration |
| Engine | `engine.go` | 559 | Core: parsing, Python bridge, result building |
| Parser | `parser.go` | 627 | All format parsers + prose generator |
| STRIDE | `stride.go` | 125 | 17 category + 33 keyword rules |
| Explain | `explain.go` | 129 | Data structures for explainability |
| Justify | `justify.go` | 717 | 6 engines: evidence, STRIDE, likelihood, impact, risk, confidence |
| License | `license.go` | 98 | HMAC license validation, generation |
| AI | `ai.go` | 227 | AI prompt builder, response parser, merge |
| Model | `model.go` | 216 | Ollama manager: download, list, delete, generate |
| Styles | `styles.go` | 227 | 4 themes, lipgloss style set |

## Configuration Flow

```
~/.asf/config.yaml                          Engine Config
┌──────────────────────────┐                ┌──────────────────┐
│ general:                 │    LoadConfig   │ Engine{           │
│   theme: Dark            │──────────────▶│   config: cfg      │
│   fox_style: Classic     │                │   pythonPath      │
│ analysis:                │                │   projectDir      │
│   depth: deep            │                │   strideEngine    │
│   stride: true           │                │   explainPipe     │
│   controls: true         │                │   archDesc        │
│   risk_threshold: low    │                └──────────────────┘
│ ai:                      │
│   enabled: false         │    ConfigPath() auto-migrates from
│   active_model: ""       │    ~/.config/asf/config.yaml
│ output:                  │
│   default: markdown      │
│   directory: ./reports   │
│ appearance:              │
│   theme: Dark            │
│   fox_style: Classic     │
└──────────────────────────┘
```

## STRIDE Engine Architecture

```
StrideEngine (stride.go)
│
├── categoryRules: map[string][]StrideCategory
│   ├── IDENTITY      → {Spoofing, EoP}
│   ├── AUTHENTICATION → {Spoofing, EoP}
│   ├── AUTHORIZATION  → {EoP, InfoDisclosure}
│   ├── ACCESS         → {EoP, InfoDisclosure}
│   ├── NETWORK        → {InfoDisclosure, DoS, Tampering}
│   ├── ENCRYPTION     → {InfoDisclosure}
│   ├── CONFIGURATION  → {Tampering}
│   ├── DEPENDENCY     → {DoS, Tampering}
│   ├── PROCESS        → {Repudiation, Tampering}
│   ├── DATABASE       → {Tampering, InfoDisclosure}
│   ├── LOGGING        → {Repudiation, Tampering}
│   ├── BACKUP         → {InfoDisclosure, DoS}
│   ├── SESSION        → {Spoofing, EoP}
│   ├── THIRD_PARTY    → {Tampering, InfoDisclosure}
│   ├── DOCUMENTATION  → {Repudiation}
│   ├── GOVERNANCE     → {Repudiation, Tampering}
│   └── GENERAL        → {}
│
└── keywordRules: 33 patterns
    ├── "idor"       → {InfoDisclosure, EoP}
    ├── "mfa"        → {Spoofing}
    ├── "sqli"       → {Tampering, InfoDisclosure}
    ├── "xss"        → {Tampering, InfoDisclosure}
    ├── "privilege escal" → {EoP}
    ├── "dos"/"ddos" → {DoS}
    ├── ... 27 more patterns
    └── "container escape" → {EoP}
```

## Export Pipeline

```
ExportResult(result, format, outputDir)
    │
    ├── ExportJSON     → json.MarshalIndent → .json
    ├── ExportMarkdown → template builder   → .md
    ├── ExportCSV      → CSV writer         → .csv
    ├── ExportPDF      → go-pdf/fpdf        → .pdf
    └── ExportHTML     → HTML template      → .html
```

## AI Integration (Optional)

```
AI Enhancement (ai.go + model.go)
─────────────────────────────────
Required: Ollama running locally (http://localhost:11434)

Flow:
  1. ASF analysis completes
  2. AIEnhancer.buildPrompt() creates structured prompt
  3. ModelManager.Generate() → POST /api/generate
  4. parseResponse() extracts 4 sections:
     • Additional assumptions
     • Risk refinements
     • Missing threats
     • Recommendations
  5. mergeAIResults() adds AI- prefixed assumptions + AI-CTRL controls
```

## Verification Results

| Result | Meaning |
|--------|---------|
| `VERIFIED` | Evidence supports the assumption |
| `CONTRADICTED` | Evidence contradicts the assumption |
| `PARTIALLY_VERIFIED` | Mixed evidence |
| `UNKNOWN` | No matching evidence |

## Risk Model

`asf-risk-model-1.0` — Deterministic 5×5 matrix:
- Likelihood factors: Exposure, Auth Dependency, Attack Complexity
- Impact factors: Data Classification, Regulatory Exposure, Business Criticality
- Score = L × I (1-25)
- Levels: Low (1-4), Medium (5-11), High (12-19), Critical (20-25)
