# V18 Confidence & Explainability Engine — Certification

## Engine Components

### 1. Confidence Explainability Engine (`asf/confidencex/engine.go`)
- **`NewExplainabilityEngine(domain string, inputs []ConfidenceInput)`** — initializes engine with domain context and input assumptions.
- **`(e *ExplainabilityEngine) RunAll() *ConfidenceOutput`** — produces per-assumption breakdowns, CISO trust view, architect review view.
- Deterministic: no AI, no LLMs, no hidden scoring. Every confidence score is traced to specific positive/negative factors.

### 2. Confidence Breakdown (`ConfidenceBreakdown`)
Per-assumption explainability output:

| Field | Description |
|-------|-------------|
| `FinalConfidence` | Existing confidence score (unchanged, 0–100) |
| `AdjustedConfidence` | Adjusted score factoring in bias from contributions |
| `StabilityClass` | Very Stable / Stable / Moderate / Weak / Highly Speculative |
| `StabilityReason` | Human-readable explanation of stability classification |
| `PositiveFactors` | Factors that increase confidence with name, impact, description |
| `NegativeFactors` | Factors that decrease confidence with name, impact, description |
| `SupportingFacts` | Fact contributions with fact ID, text, contribution %, positivity |
| `EvidenceContributions` | Evidence present/absent with impact and label (Strong/Moderate/Weak/Missing) |
| `DomainContribution` | Domain pack influence with domain, influence %, reason, strength |
| `TrustContribution` | Trust chain influence, dependency centrality, failure radius influence |
| `WhyExists` | Human-readable explanation of why ASF believes this assumption |
| `WhyUncertain` | Human-readable explanation of why ASF is uncertain |
| `WhatIncreasesConfidence` | Actionable suggestions for increasing confidence |
| `WhatDecreasesConfidence` | Factors that could reduce confidence |

### 3. Positive Factors Collected
- High Baseline Confidence (≥80%: +5, ≥60%: +2)
- Component Traceability (+3)
- Rich Keyword Coverage (≥3 keywords: +2)
- Trust Chain Support (+5)
- Explicit Evidence Source (+4)
- Domain Alignment (+3)
- Detailed Rationale (+2)
- Supporting Facts Present (up to +10)
- High Dependency Centrality (>0.5: +4)

### 4. Negative Factors Collected
- Low Baseline Confidence (<30%: −5, <50%: −2)
- No Evidence Sources (−8)
- Unverified Status (−5)
- Coverage Gap Detected (−6)
- High Blind Spot Risk (>50: −5, >20: −2)
- No Architecture Traceability (−7)
- Missing Rationale (−3)
- No Supporting Facts (−4)
- No Trust Chain (−3)

### 5. Domain Contribution
Domain packs influence confidence based on keyword matching:

| Domain | Strong Influence | Moderate Influence | Weak Influence |
|--------|-----------------|-------------------|----------------|
| Healthcare | +12% (PHI/HIPAA/encryption/audit) | +5% | — |
| Fintech | +12% (PCI/encryption/fraud) | +5% | — |
| Kubernetes | +10% (pod/secret/RBAC/policy) | +4% | — |
| General | — | — | +2% |

### 6. Confidence Stability Classification
| Stability | Conditions |
|-----------|------------|
| Very Stable | Confidence ≥85 and net bias ≤15 |
| Stable | Confidence ≥70 and net bias ≤25 |
| Moderate | Confidence ≥50 |
| Weak | Confidence ≥30 |
| Highly Speculative | Confidence <30 |

### 7. Trust Contribution
- Chain influence: +5 (has trust chain) / −3 (no trust chain)
- Dependency centrality: 0–1 scale
- Failure radius influence: failureRadius × 0.5

## Data Model (`asf/confidencex/model.go`)

| Struct | Purpose |
|--------|---------|
| `ConfidenceInput` | Per-assumption input (assumption data, facts, evidence, trust, domain) |
| `ConfidenceFactor` | Named factor with type (positive/negative), impact, and description |
| `FactContribution` | Fact ID, text, contribution percentage, positivity flag |
| `EvidenceContribution` | Evidence ID, present/absent, impact, label (Strong/Moderate/Weak/Missing) |
| `DomainContribution` | Domain name, influence %, reason, strength (Strong/Moderate/Weak) |
| `TrustContribution` | Trust chain presence, chain influence, dependency centrality, failure radius |
| `ConfidenceBreakdown` | Complete per-assumption explainability output |
| `CISOTrustView` | Most trusted, least trusted, critical low-confidence, highest-risk unknowns |
| `ArchitectReviewView` | Requiring validation, weak support, strong support |
| `ConfidenceOutput` | Top-level container with all breakdowns and views |

## Export Formats (`asf/confidencex/export.go`)

### Markdown (`ExportMarkdown`)
- `# Confidence & Explainability Report` header
- Per-assumption sections: confidence, stability, why exists, why uncertain, increases/decreases
- Tables for positive factors, negative factors, supporting facts, evidence contributions
- Domain contribution section, trust contribution section
- CISO Trust View and Architect Review View sections

### HTML (`ExportHTML`)
- Same content rendered with HTML5
- Color-coded confidence (green ≥70, yellow 40–69, red <40)
- CSS-styled stability badges
- Tables, sections, lists for all data

## TUI Integration (`results.go`)

Three new sections in the results view (indices 34–36):

| Index | Section | Render Function |
|-------|---------|----------------|
| 34 | Confidence View | `renderConfidenceView` — per-assumption confidence with color coding and CISO summary |
| 35 | Explainability View | `renderExplainabilityView` — why exists, why uncertain, increases per assumption |
| 36 | Confidence Breakdown | `renderConfidenceBreakdownView` — full factor analysis with facts, domain, trust, CISO/architect views |

## Integration Points (`engine.go`)

### New Field on `AnalysisResult`
```go
ConfidenceOutput *confidencex.ConfidenceOutput `json:"confidence_output,omitempty"`
```

### New Method
```go
func (e *Engine) runConfidenceExplainability(result *AnalysisResult) *confidencex.ConfidenceOutput
```
Maps `engine.Assumption` → `confidencex.ConfidenceInput` for each assumption, translating:
- `Confidence × 100` → `Confidence` (0–100 scale)
- `EvidenceSources`, `SourceComponents`, `Keywords`, `Rationale` — passed directly
- `VerificationStatus` — passed directly
- `TrustOutput.TrustChains` → `HasTrustChain` (matched by assumption ID)
- `CoverageOutput.BlindSpots` → `HasCoverageGap`, `BlindSpotScore`
- `len(SourceRelationships) / 10` → `DependencyCentrality`
- `Impact` → `FailureRadius`

## Export CLI Integration (`export.go`)

Two new format constants:

| Constant | Value | File Extension |
|----------|-------|----------------|
| `ExportConfidenceMarkdown` | `"confidence-md"` | `_confidence.md` |
| `ExportConfidenceHTML` | `"confidence-html"` | `_confidence.html` |

Added to export view format list (14 formats total) and dispatch in `ExportResult()`.

## Test Coverage (`confidencex_test.go`)

| Test | What It Verifies |
|------|-----------------|
| `TestEmptyInput` | Empty input produces non-nil output with 0 breakdowns |
| `TestSingleAssumptionExplainability` | Full breakdown: factors, facts, evidence, domain, trust, explanations, stability |
| `TestLowConfidenceAssumption` | Low-confidence input classified as Highly Speculative/Weak, no-evidence factor detected |
| `TestDomainContributions` | 6 sub-tests: healthcare PHI (+12), healthcare generic (+5), fintech PCI (+12), fintech generic (+5), K8s RBAC (+10), general (+2) |
| `TestStabilityClassifications` | 5 sub-tests: Very Stable, Stable, Moderate, Weak, Highly Speculative |
| `TestCISOTrustView` | Most trusted, least trusted, critical low-confidence lists |
| `TestArchitectReviewView` | Requiring validation, weak support, strong support categorization |
| `TestExportMarkdown` | Report contains title, assumption text, why section, uncertainty section |
| `TestExportHTML` | Report contains title, HTML tag, assumption text |
| `TestExplainabilityDetail` | Full detailed output logged for visual inspection |
| `BenchmarkExplainabilityEngine` | 100 items: ~372k ns/op, 498kB, 4738 allocs |
| `BenchmarkExplainabilityLarge` | 500 items: ~1.58ms/op, 2.8MB, 19527 allocs |

## Regression

- `go build ./...` — passes (21 packages)
- `go vet ./...` — passes
- `go test ./...` — all 21 packages pass (including V1–V17 tests)
- V1–V17 engines unchanged: Facts, Assumptions, Narratives, Verification, Coverage, Review Workbench
- CLI unchanged
- TUI unchanged (new sections appended at end)

## Final Verdict

**CONFIDENCE_ENGINE_CERTIFIED**

The Confidence & Explainability Engine (V18) meets all certification criteria:

1. ✅ **Deterministic** — no AI, no LLMs, no hidden scoring
2. ✅ **Explainable** — every confidence score traced to specific factors
3. ✅ **Transparent** — why exists, why uncertain, what increases/decreases confidence
4. ✅ **Complete** — fact, evidence, domain, trust contribution analysis
5. ✅ **Stable** — confidence stability classification with reasoning
6. ✅ **CISO-ready** — most trusted, least trusted, critical low-confidence, highest-risk unknowns
7. ✅ **Architect-ready** — requiring validation, weak support, strong support views
8. ✅ **Exportable** — Markdown and HTML Confidence & Explainability Reports
9. ✅ **Integrated** — TUI sections, export CLI, analysis pipeline
10. ✅ **Non-breaking** — all 21 packages pass regression tests
