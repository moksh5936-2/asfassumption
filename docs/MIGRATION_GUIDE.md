# Migration Guide: Explainability Transformation

## Overview

This migration adds explainability to all ASF assumptions. Every assumption now includes traceable evidence, STRIDE justification, risk decomposition, and confidence scoring.

## What Changed

### Data Structures

**`Assumption`** (engine.go) — 10 new fields added:

| New Field | Type | Description |
|-----------|------|-------------|
| `EvidenceSources` | `[]string` | Traceable evidence file paths and matched artifacts |
| `SourceComponents` | `[]string` | Architecture components that triggered the assumption |
| `SourceRelationships` | `[]string` | Architecture relationships that triggered the assumption |
| `Rationale` | `string` | Human-readable explanation of why the assumption exists |
| `StrideJustifications` | `[]StrideJustification` | Per-STRIDE-category justification (reason, rules, confidence) |
| `RiskJustification` | `*RiskJustification` | Complete risk decomposition (likelihood, impact, score, factors) |
| `ReviewStatus` | `string` | Proposed / Accepted / Rejected / Modified |
| `ReviewNotes` | `string` | Architect review notes |
| `ReviewTimestamp` | `time.Time` | When the review decision was made |

**`AnalysisResult`** (engine.go) — 3 new fields:

| New Field | Type | Description |
|-----------|------|-------------|
| `EvidenceSummary` | `EvidenceSummary` | Aggregated evidence across all assumptions |
| `RiskModelVersion` | `string` | Version of the risk model used (currently `asf-risk-model-1.0`) |
| `ConfidenceSummary` | `string` | Average confidence across all assumptions |

### New Files

| File | Purpose |
|------|---------|
| `explain.go` | All new data structures (EvidenceSource, StrideJustification, RiskJustification, ReviewRecord, ValidationRecord) |
| `justify.go` | 7 justification engines + orchestration pipeline |
| `review.go` | TUI review mode with browse/detail views |

### Modified Files

| File | Change |
|------|--------|
| `engine.go` | `buildResult()` now calls `ExplainabilityPipeline.Explain()` for each assumption |
| `stride.go` | Added `GetKeywordRules()` and `GetCategoryRules()` accessors |
| `export.go` | All 5 export formats enhanced with evidence, reasoning, confidence, risk factors |
| `results.go` | Confidence display, 5×5 risk matrix, rationale shown |
| `app.go` | Added `reviewView`, `r` key handler for review navigation |

### Unchanged Files

`main.go`, `parser.go`, `config.go`, `settings.go`, `startup.go`, `dashboard.go`, `localai.go`, `analyze.go`, `about.go`, `styles.go`, `license.go`, `ai.go`, `model.go`

## Backward Compatibility

### JSON Exports
- Old exports will lack the new fields when loaded — `json.Unmarshal` will leave them as zero values
- New exports include all 10 explainability fields per assumption
- ✅ **No breaking changes** — old code can still read new exports, new code can read old exports

### Markdown/HTML/CSV/PDF Exports
- All formats now include expanded content
- Old scripts parsing these formats may need updates for the new structure
- ⚠️ **CSV** — new columns added between existing ones

### Config
- `~/.asf/config.yaml` — unchanged, no migration needed
- New fields are initialized to zero values when loading old configs

### API
- `Engine.RunAnalysis()` — unchanged signature
- `AnalysisResult` — same struct with more fields
- `Assumption` — same struct with more fields

## New Dependencies

None. All new code uses only the standard library.

## Testing Your Setup

1. Run an analysis on any architecture file:
   ```
   ./asf-tui
   ```
2. Navigate to Results → Assumptions to see evidence and confidence
3. Press `r` in Results to enter Review mode
4. Press `e` in Results to export with enhanced formats

## Rollback

To revert to the previous version, restore the old binary and old source files. No config or data migration is needed since the changes are purely additive.
