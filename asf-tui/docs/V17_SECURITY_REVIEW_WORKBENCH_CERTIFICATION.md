# V17 Security Review Workbench — Certification

## Engine Components

### 1. Review Priority Engine (`asf/review/scoring.go`)
- **`NewReviewEngine(domain string, inputs []ReviewInput) *ReviewEngine`** — initializes engine with domain context and assumptions.
- **`(e *ReviewEngine) RunAll() *ReviewOutput`** — produces full output: queue, matrix, campaigns, CISO dashboard, domain view.
- Priority formula (0–100): `min(riskContribution + centralityBoost + failureContrib + supportContrib + blindSpotBoost - confidenceDiscount + domainBoost, 100)`
  - `riskContribution`: Critical=30, High=20, Medium=10, Low=0
  - `centralityBoost`: centrality × 20
  - `failureContrib`: failureRadius × 1.0
  - `supportContrib`: supportCount × 0.5
  - `blindSpotBoost`: coverageGap ? (25 + 0.3×blindSpotScore) : 0
  - `confidenceDiscount`: (100 - verificationConfidence) × 0.1
  - `domainBoost`: "healthcare" → +10

### 2. Priority Matrix (`ReviewMatrix`)
- 2×2 quadrant classification based on `ReviewValue` (Very High/High/Medium/Low) and `ReviewEffort` (Low/Medium/High).
- Four quadrants: High Value / Low Effort, High Value / High Effort, Low Value / Low Effort, Low Value / High Effort.

### 3. Review Queue (`ReviewQueue`)
- All items sorted by `PriorityScore` descending, ranked 1..N.
- Each item includes: why-review rationale, what-to-review guidance, expected evidence, estimated time.

### 4. Review Campaigns
- **30 Minute Review Plan** — top items fitting in ~30 minutes of estimated time.
- **High Value / Low Effort** — items in the top-left quadrant.
- **Deep Dive Required** — high-value/high-effort items.
- **Full Campaign** — all items if ≤ 20, else top 20.

### 5. CISO Dashboard (`CISOReviewDashboard`)
- Summary counts and two curated lists:
  - `HighestRiskAssumptions` — top 5 highest-scoring items.
  - `GreatestRiskReduction` — top 5 items with highest `ExpectedRiskReduction`.

### 6. Domain Prioritization (`DomainPrioritization`)
- Domain-specific focus areas and top 5 priorities.
- Healthcare domain adds specific focus areas (PHI access controls, audit logging, etc.).

## Data Model (`asf/review/model.go`)

| Struct | Purpose |
|--------|---------|
| `ReviewInput` | Per-assumption input (risk, centrality, failure radius, support count, coverage gap, domain) |
| `ReviewPriority` | Scored & ranked queue item with quadrant, rationale, evidence guidance |
| `ReviewQueue` | Sorted list of items with summary counts |
| `ReviewMatrix` | 2×2 quadrants |
| `ReviewCampaign` | Named time-boxed review campaign |
| `CISOReviewDashboard` | Executive summary of highest risks and reduction opportunities |
| `DomainPrioritization` | Domain-specific prioritization output |
| `ReviewOutput` | Top-level container for all workbench output |

## Export Formats (`asf/review/export.go`)

### Markdown (`ExportMarkdown`)
- `# Security Review Workbench Report` header
- Priority Queue table (Rank | Score | Risk | Category | Assumption | Why Review)
- Quadrant sections with item details
- Campaign summaries
- CISO Dashboard executive summary

### HTML (`ExportHTML`)
- Same content as Markdown rendered with HTML5 structure
- `<title>Security Review Workbench Report</title>`
- Semantic sections, tables, and styled elements
- Inline CSS for basic presentation

## TUI Integration (`results.go`)

Four new sections in the results view (indices 30–33):

| Index | Section | Render Function |
|-------|---------|----------------|
| 30 | Review Queue | `renderReviewQueue` — ranked table with color-coded risk |
| 31 | Review Priority Matrix | `renderReviewMatrix` — 2×2 quadrant view |
| 32 | Review Campaigns | `renderReviewCampaigns` — campaign list with items |
| 33 | CISO Review Dashboard | `renderCISOReviewDashboard` — executive summary |

Each section renders through `renderSectionContent()` dispatch integrated into the existing `resultsModel` view framework.

## Integration Points (`engine.go`)

### New Field on `AnalysisResult`
```go
ReviewOutput *review.ReviewOutput `json:"review_output,omitempty"`
```

### New Method
```go
func (e *Engine) runReviewAnalysis(result *AnalysisResult) *review.ReviewOutput
```
Maps `engine.Assumption` → `review.ReviewInput` for each assumption, translating:
- `QualityScore` / `Confidence` → `Centrality`
- `Impact` → `FailureRadius`
- `len(SourceComponents)` → `SupportCount`
- `len(SourceRelationships)` → `DependencyCount`
- `len(EvidenceSources) == 0` → `CoverageGap`

## Export CLI Integration (`export.go`)

Two new format constants:

| Constant | Value | File Extension |
|----------|-------|----------------|
| `ExportReviewMarkdown` | `"review-md"` | `_review.md` |
| `ExportReviewHTML` | `"review-html"` | `_review.html` |

Added to export view format list and dispatch in `ExportResult()`.

## Test Coverage (`review_test.go`)

| Test | What It Verifies |
|------|-----------------|
| `TestEmptyReview` | Empty input produces non-nil output |
| `TestSingleAssumptionReview` | Single item gets rank 1, positive score, populated rationale fields |
| `TestMultipleAssumptions` | 3 items sorted by score, ranked correctly, matrix generated |
| `TestPriorityScoring` | Critical high-centrality item scores ≥70; low-risk item scores 0 |
| `TestPriorityMatrix` | 8 items classified into quadrants correctly |
| `TestDomainPrioritization` | Healthcare domain generates domain view with focus areas |
| `TestReviewCampaigns` | 20 items generate campaigns with "30 Minute Review Plan" |
| `TestCISODashboard` | Dashboard has highest risk and risk reduction lists |
| `TestExportMarkdown` | Markdown output contains report title and assumption text |
| `TestExportHTML` | HTML output contains report title and opens with `<html>` |
| `TestReviewPrecision` | 10 mixed items scored 0–100, ranked 1–10, valid quadrants |
| `BenchmarkReviewEngine` | 100 items: ~222k ns/op, 334kB, 1111 allocs |
| `BenchmarkReviewLarge` | 500 items: ~1ms/op, 1.4MB, 6093 allocs |

## Regression

- `go build ./...` — passes
- `go vet ./...` — passes
- `go test ./...` — all 20 packages pass (including existing V1–V16 tests)
