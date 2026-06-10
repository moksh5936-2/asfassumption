# Explainability Gap Analysis

## Current State (Before Explainability Transformation)

### Assumption Generation
- **Logic**: Python ASF CLI extracts assumptions from architecture text via `asf.cli.main analyze --json`
- **Input**: Pre-processed prose from parser (drawio → text, mermaid → text, etc.)
- **Output**: JSON with assumption ID, text, type, confidence, keywords
- **Limitation**: No link back to which components, relationships, or diagram elements triggered the assumption

### STRIDE Mapping
- **Logic**: `StrideEngine.MapAssumption()` in `stride.go` — rule-based matching
  - 17 category-level rules (IDENTITY → Spoofing+EoP, NETWORK → InfoDisclosure+DoS+Tampering, etc.)
  - 33 keyword-level rules (idor → InfoDisclosure+EoP, audit log → Repudiation+Tampering, etc.)
- **Output**: `[]StrideCategory` — list of matched STRIDE categories
- **Limitation**: No justification. No record of which rules matched. No confidence score. No explanation of why a category was assigned.

### Risk Scoring
- **Logic**: `mapRiskLevel()` in `engine.go` — maps Python ASF severity + verification status to RiskLevel
  - CONTRADICTED → Low
  - CRITICAL/HIGH/MEDIUM severity → corresponding level
  - default → Medium
- **Output**: Single `RiskLevel` enum (Critical/High/Medium/Low)
- **Limitation**: No likelihood/impact separation. No risk score (1-25). No reasoning. Defaults to Medium when severity unknown. `riskToLikelihoodImpact()` hardcodes L/I pairs per risk level (Critical=5/5, High=4/4, etc.) — identical scores, no factor analysis.

### Confidence Scoring
- **Logic**: Confidence comes from Python ASF JSON (`a.Confidence`) — source unknown
- **Output**: Single float64
- **Limitation**: No decomposition. No explanation of what factors contribute. No evidence-based confidence calculation.

### Control Recommendations
- **Logic**: `generateControls()` in `engine.go` — maps assumption category to hardcoded control strings
  - IDENTITY → "Implement strong identity verification with MFA"
  - ACCESS → "Enforce least-privilege access controls"
  - etc.
- **Output**: []string of control descriptions
- **Limitation**: No rationale. No traceability to specific threats.

### Review Support
- **Logic**: None
- **Output**: None
- **Limitation**: No concept of human review. Assumptions are generated and presented as-is. No Accepted/Rejected/Modified workflow.

### Validation Data
- **Logic**: None
- **Output**: None
- **Limitation**: No structures for precision/recall/accuracy studies.

## Evidence Chains

### Missing: Evidence → Assumption
- An assumption exists because of specific architecture elements, but no record of which ones.
- An architect cannot ask "what evidence supports this assumption?"

### Missing: Assumption → STRIDE
- A STRIDE category is assigned, but no record of why.
- An architect cannot ask "why is this Tampering and not Spoofing?"

### Missing: STRIDE → Risk
- A risk level is assigned, but no decomposition into likelihood and impact.
- An architect cannot ask "why is this Critical and not High?"

### Missing: Risk → Confidence
- A confidence score is assigned, but no factors.
- An architect cannot ask "how sure is ASF about this assessment?"

### Missing: Recommended Controls → Threats
- Controls exist independently of assumptions.
- An architect cannot ask "which control mitigates which threat?"

## What Was Added

### `explain.go`
- `EvidenceSource` — tracks file path, file type, matched components, relationships, trust boundaries, security concepts
- `EvidenceSummary` — aggregates evidence across all assumptions
- `StrideJustification` — per-category: reason, matched rule indexes, matched keywords, matched components, confidence
- `LikelihoodFactor` / `ImpactFactor` — decomposition of risk factors
- `RiskJustification` — complete risk calculation with likelihood, impact, score, level, reasoning, confidence
- `ReviewRecord` — Proposed/Accepted/Rejected/Modified with notes and timestamp
- `ValidationRecord` — supports future precision/recall studies

### `justify.go`
- `EvidenceEngine` — traces assumptions back to architecture artifacts (components, relationships, trust boundaries, security concepts)
- `StrideJustifyEngine` — extends StrideEngine to return per-category justifications with rule indexes, keywords, confidence
- `LikelihoodAnalyzer` — evaluates exposure, authentication dependency, attack complexity (1-5 scale)
- `ImpactAnalyzer` — evaluates data classification, regulatory exposure, business criticality (1-5 scale)
- `RiskMatrix` — 5×5 matrix: likelihood × impact = risk score (1-25)
- `ConfidenceEngine` — calculates confidence from evidence count, rule matches, component/relationship matches
- `ExplainabilityPipeline` — orchestrates all engines per assumption

### `review.go`
- `reviewModel` — TUI model for architect review (browse/detail modes)
- `CollectValidationData()` — exports assumptions as `[]ValidationRecord` for future studies
- Keyboard shortcuts: s=Accept, r=Reject, m=Modified, n=Note

### Engine updates (`engine.go`)
- `EvidenceSummary`, `RiskModelVersion`, `ConfidenceSummary` on `AnalysisResult`
- All new explainability fields on `Assumption` struct
- `buildConfidenceSummary()` — average confidence across all assumptions
- Pipeline called from `buildResult()` for every assumption

### Export updates (`export.go`)
- JSON — serializes new fields automatically
- Markdown — detailed per-assumption sections with evidence, rationale, justification, review status
- HTML — expandable detail sections with risk factors, evidence traces, STRIDE justifications
- CSV — new columns: RiskScore, Confidence, ReviewStatus, EvidenceSources, Rationale
- PDF — per-assumption pages with reasoning, evidence sources, risk justification
