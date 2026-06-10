# ASF Explainability Architecture

## Overview

ASF's explainability system makes every output traceable, explainable, defensible, and auditable. No STRIDE mapping, risk score, or confidence value exists without justification.

The system consists of four engines that form a justification chain:

```
Architecture Input
       ↓
  ┌─────────────┐
  │  Evidence   │─── Identifies source nodes, components, relationships
  │   Engine    │─── Matches assumptions to architecture artifacts
  └─────────────┘
       ↓
  ┌─────────────┐
  │   STRIDE    │─── Maps assumptions to threat categories
  │ Justifier   │─── Explains which rule triggered and why
  └─────────────┘
       ↓
  ┌─────────────┐
  │    Risk     │─── Computes likelihood × impact
  │  Analyzer   │─── Factors: exposure, auth, data, regulatory
  └─────────────┘
       ↓
  ┌─────────────┐
  │ Confidence  │─── Deterministic score from evidence
  │   Engine    │─── STRIDE matches + component matches
  └─────────────┘
       ↓
  ┌─────────────┐
  │    TUI /    │─── Results, Review, Validation views
  │   Export    │─── Markdown, HTML, PDF, CSV, JSON
  └─────────────┘
```

Every output from ASF can be traced back through this chain to its original evidence.

---

## 1. STRIDE Rules

### Category Rules

STRIDE categories are assigned based on the assumption's `Category` field (determined by the Python ASF engine's classification):

| Category | Assigned STRIDE |
|----------|----------------|
| IDENTITY | Spoofing, Elevation of Privilege |
| AUTHENTICATION | Spoofing, Elevation of Privilege |
| AUTHORIZATION | Elevation of Privilege, Information Disclosure |
| ACCESS | Elevation of Privilege, Information Disclosure |
| NETWORK | Information Disclosure, Denial of Service, Tampering |
| ENCRYPTION | Information Disclosure |
| CONFIGURATION | Tampering |
| DEPENDENCY | Denial of Service, Tampering |
| PROCESS | Repudiation, Tampering |
| DATABASE | Tampering, Information Disclosure |
| LOGGING | Repudiation, Tampering |
| BACKUP | Information Disclosure, Denial of Service |
| SESSION | Spoofing, Elevation of Privilege |
| THIRD_PARTY | Tampering, Information Disclosure |
| DOCUMENTATION | Repudiation |
| GOVERNANCE | Repudiation, Tampering |
| GENERAL | _(none)_ |

### Keyword Rules

Keywords are matched against the assumption text, description, and existing keywords. Each keyword rule maps to one or more STRIDE categories.

**Algorithm:**

1. Build search text: `strings.ToLower(category + " " + text + " " + keywords)`
2. For each keyword rule, check if any keyword is contained in the search text
3. If matched, add the rule's STRIDE categories
4. De-duplicate by category

**Keyword-to-STRIDE mapping** (30 rules):

| Keyword(s) | STRIDE |
|------------|--------|
| idor, insecure direct object | Information Disclosure, Elevation of Privilege |
| bola, broken object level | Information Disclosure, Elevation of Privilege |
| session hijack, session fixat, session predict | Spoofing, Elevation of Privilege |
| audit log, audit trail, log immutable, log tamper | Repudiation, Tampering |
| backup, data loss, data recover | Information Disclosure, Denial of Service |
| mfa, multi factor, two factor, 2fa | Spoofing |
| sql injection, sqli, nosql injec | Tampering, Information Disclosure |
| key management, key rotat, key stor | Information Disclosure |
| buffer overflow, memory corrupt | Tampering, Denial of Service |
| cross site script, xss | Tampering, Information Disclosure |
| csrf, cross site request | Tampering, Elevation of Privilege |
| ssrf, server side request | Information Disclosure, Elevation of Privilege |
| privilege escal | Elevation of Privilege |
| denial of serv, dos, ddos | Denial of Service |
| man in the middl, mitm | Spoofing, Tampering, Information Disclosure |
| replay attack | Spoofing, Elevation of Privilege |
| tls, ssl, https | Information Disclosure, Tampering |
| auth bypass, authn bypass, authentication bypass | Spoofing, Elevation of Privilege |
| rate limit | Denial of Service |
| supply chain | Tampering, Denial of Service |
| secret, credential, password, token | Spoofing, Information Disclosure |
| firewall, acl, network segment | Denial of Service, Information Disclosure |
| encrypt, decrypt, cipher | Information Disclosure, Tampering |
| signing, signature | Spoofing, Tampering |
| certificate, cert | Spoofing, Information Disclosure |
| oauth, saml, oidc | Spoofing, Elevation of Privilege |
| rbac, abac, access control | Elevation of Privilege, Information Disclosure |
| monitoring, alert, detect | Repudiation |
| patch, update | Tampering, Denial of Service |

---

## 2. STRIDE Justification

When a STRIDE category is assigned, a `StrideJustification` is generated containing:

```go
type StrideJustification struct {
    Category           StrideCategory  // e.g., "Spoofing"
    Reason             string          // Human-readable explanation
    MatchedRuleIndexes []int           // Which keyword rules fired
    MatchedKeywords    []string        // Which specific keywords matched
    MatchedComponents  []string        // Which architecture components
    Confidence         float64         // 0.0–1.0
    ConfidenceReason   string          // Why this confidence
}
```

**Reason template** (per STRIDE category):

| Category | Reason Template |
|----------|----------------|
| Spoofing | `identity verification required — <category> (matched: <keywords>)` |
| Tampering | `data integrity risk — <category> (matched: <keywords>)` |
| Repudiation | `non-repudiation concern — <category> (matched: <keywords>)` |
| Information Disclosure | `information disclosure risk — <category> (matched: <keywords>)` |
| Denial of Service | `availability risk — <category> (matched: <keywords>)` |
| Elevation of Privilege | `privilege escalation risk — <category> (matched: <keywords>)` |

**STRIDE Confidence formula:**

```
base = 0.3
+ min(keywordMatches × 0.10, 0.40)  // up to 4 keyword matches
+ min(componentMatches × 0.08, 0.30) // up to ~4 component matches
capped at 0.95
```

---

## 3. Risk Scoring

### Likelihood Analysis

Three factors are evaluated:

| Factor | Base | Conditions | Max |
|--------|------|-----------|-----|
| Exposure Level | 1 (internal) | 3 (network) / 4 (internet category) / 5 (internet-exposed component) | 5 |
| Authentication Dependency | 2 (standard) | 4 (auth component or auth/authorization concept) | 4 |
| Attack Surface Complexity | 2 (moderate) | 3 (>3 components) / 4 (>5 relationships) | 4 |

**Formula:** `likelihood = 1 + (exposure - 1) + (auth - 2) + (complexity - 2)`, clamped 1–5.

### Impact Analysis

Three factors are evaluated:

| Factor | Base | Conditions | Max |
|--------|------|-----------|-----|
| Data Classification | 2 (standard) | 4 (data_protection concept or db component) / 5 (PII/PHI/financial) | 5 |
| Regulatory Exposure | 1 (none) | 4 (GDPR/SOX/PII) / 5 (HIPAA/PCI DSS) | 5 |
| Business Criticality | 2 (standard) | 4 (core component or >8 relationships) | 4 |

**Formula:** `impact = 1 + (data - 2) + (regulatory - 1) + (criticality - 2)`, clamped 1–5.

### Risk Matrix

```
riskScore = likelihood × impact  (range: 1–25)

  1–4:   Low
  5–11:  Medium
  12–19: High
  20–25: Critical
```

**RiskJustification output:**

```go
type RiskJustification struct {
    Likelihood        int               // 1–5
    LikelihoodReason  string            // Human-readable
    LikelihoodFactors []LikelihoodFactor // Per-factor breakdown
    Impact            int               // 1–5
    ImpactReason      string            // Human-readable
    ImpactFactors     []ImpactFactor    // Per-factor breakdown
    RiskScore         int               // 1–25
    RiskLevel         RiskLevel         // Low/Medium/High/Critical
    RiskReason        string            // "risk score N/25 (L × I) = Level"
    Confidence        float64           // From confidence engine
    ConfidenceReason  string            // Why this confidence
}
```

---

## 4. Confidence Scoring

### Deterministic Formula

Confidence is computed from observable evidence only — no randomness or heuristics:

```
confidence = 0.10 (base)
+ min(evidenceCount × 0.05, 0.30)        // evidence points
+ min(strideMatchCount × 0.08, 0.25)      // STRIDE keyword matches
+ min(componentCount × 0.06, 0.20)        // matched component names
+ min(relationshipCount × 0.04, 0.15)      // matched relationships
= final (capped at 0.95)
```

| Factor | Contribution | Cap |
|--------|-------------|-----|
| Base | +0.10 | 0.10 |
| Evidence points | +0.05 per point | 0.30 |
| STRIDE keyword matches | +0.08 per match | 0.25 |
| Component matches | +0.06 per match | 0.20 |
| Relationship matches | +0.04 per match | 0.15 |
| **Maximum** | | **0.95** |

The formula is deterministic — identical inputs always produce identical confidence scores.

For display, multiply by 100: `0.74 → 74%`.

---

## 5. Evidence Extraction

### EvidenceEngine

The `EvidenceEngine` traces assumptions back to architecture artifacts:

1. **Component matching:** For each component in the parsed architecture, check if the component label appears in the assumption's category, text, or keywords.
2. **Relationship matching:** For each relationship, check if the source or target component name appears in the assumption text.
3. **Trust boundary detection:** Components with trust-relevant labels (internet, external, public, vpn, gateway, dmz) that appear alongside matched components are flagged as trust boundaries.
4. **Security concept matching:** Assumption text is checked against 8 concept groups:
   - authentication, authorization, encryption, network_security, data_protection, audit_logging, dependency, session_management
5. **Source line identification:** The raw architecture text is searched line-by-line for the matched component or assumption text to determine the source line number.

### EvidenceResult

```go
type EvidenceResult struct {
    MatchedComponents       []string  // Architecture component labels
    MatchedRelationships    []string  // "A → B" format
    MatchedTrustBoundaries  []string  // Trust boundary descriptions
    MatchedSecurityConcepts []string  // Security domain tags
    PrimarySourceNode       string    // Best matching component
    SourceLine              int       // Line number in raw text
    EvidenceCount           int       // Total evidence items
}
```

### Assumption Traceability

Every `Assumption` now carries:

```go
SourceNode string  // Primary architecture node the assumption relates to
SourceLine int     // Line number in the source architecture document
```

These are populated by the ExplainabilityPipeline during analysis.

---

## 6. Justification Chain Architecture

```
ExplainabilityPipeline.Explain(assumption)
  │
  ├─ 1. EvidenceEngine.TraceEvidence()
  │      Returns: EvidenceResult with matched components, relationships, etc.
  │
  ├─ 2. EvidenceEngine.FindSourceLine()
  │      Returns: line number in raw architecture text
  │
  ├─ 3. Build evidence sources
  │      EvidenceEngine.BuildEvidenceSources()
  │      → assumption.EvidenceSources
  │      → assumption.SourceNode
  │      → assumption.SourceLine
  │
  ├─ 4. Generate rationale
  │      JustifyAssumption()
  │      → assumption.Rationale
  │
  ├─ 5. STRIDE justification
  │      StrideJustifyEngine.Justify()
  │      → assumption.Stride (categories)
  │      → assumption.StrideJustifications (with reasons + confidence)
  │
  ├─ 6. Risk justification
  │      LikelihoodAnalyzer.AnalyzeLikelihood()
  │      ImpactAnalyzer.AnalyzeImpact()
  │      RiskMatrix.Calculate()
  │      → assumption.Likelihood, assumption.Impact
  │      → assumption.RiskJustification (with factors + reasons)
  │
  └─ 7. Confidence scoring
        ConfidenceEngine.CalculateConfidence()
        → assumption.Confidence
        → assumption.RiskJustification.Confidence
```

### Output Targets

| Output | Evidence | STRIDE Justif. | Risk Justif. | Confidence | Source Trace |
|--------|----------|----------------|--------------|------------|-------------|
| TUI Results | ✅ | ❌ (not displayed) | ❌ (not displayed) | ✅ (%) | ❌ |
| TUI Review | ✅ | ❌ (not displayed) | ✅ (score only) | ✅ | ❌ |
| TUI Validation | ✅ | ✅ | ✅ | ✅ (with factors) | ✅ |
| Markdown Export | ✅ | ✅ | ✅ (with factors) | ✅ | ✅ |
| HTML Export | ✅ | ✅ | ✅ (with factors) | ✅ | ✅ |
| PDF Export | ✅ | ✅ (new) | ✅ | ✅ | ✅ (new) |
| CSV Export | ✅ | ✅ (new) | ✅ (new) | ✅ | ✅ (new) |
| JSON Export | ✅ | ✅ | ✅ | ✅ | ✅ |

---

## 7. Reproducibility Guarantee

All ASF outputs are deterministic:

- **STRIDE mapping:** Same keywords + category → same STRIDE categories every time.
- **Risk scoring:** Same likelihood/impact factors → same risk score every time.
- **Confidence scoring:** Same evidence counts → same confidence value every time.
- **Evidence matching:** Same architecture + assumption → same evidence every time.

There is no randomness, no machine learning inference, and no non-deterministic component in the justification chain. The optional AI enhancement (Local AI mode) runs separately and is clearly marked with AI-prefixed assumption IDs.

---

## 8. Architecture of Key Components

```
StrideEngine
  ├── MapAssumption(category, text, keywords) → []StrideCategory
  ├── GetKeywordRules() → []keywordRule
  └── GetCategoryRules() → map[string][]StrideCategory

StrideJustifyEngine
  └── Justify(category, text, keywords, components) → StrideResult
        ├── Categories []StrideCategory
        └── Justifications []StrideJustification

LikelihoodAnalyzer
  └── AnalyzeLikelihood(assumption, evidence) → (score, reason, factors)

ImpactAnalyzer
  └── AnalyzeImpact(assumption, evidence) → (score, reason, factors)

RiskMatrix
  ├── Calculate(likelihood, impact) → (score, level)
  └── RiskReason(likelihood, impact, score, level) → string

ConfidenceEngine
  └── CalculateConfidence(evidenceCount, strideCount, compCount, relCount) → (score, reason)

EvidenceEngine
  ├── TraceEvidence(category, keywords, text) → EvidenceResult
  ├── FindSourceLine(searchText, evidence) → int
  ├── BuildEvidenceSources(evidence) → []string
  └── BuildEvidenceSummary(assumptions) → EvidenceSummary

ExplainabilityPipeline
  └── Explain(assumption) → populates all Assumption fields
```

---

## 9. Validation Mode

The Validation Mode (accessible by pressing `v` from Results or Review views) provides a developer-focused display showing every assumption with:

- Full assumption text
- Source node and source line
- All evidence sources
- STRIDE categories with per-category justification and triggered keywords
- Risk level, score, likelihood/impact breakdown with all factors
- Overall confidence with factor breakdown

This mode exists solely to evaluate and improve model quality. It does not provide additional analysis functionality.

---

## 10. Limitations

1. **Evidence line numbers** are best-effort — they search the raw text for the best matching content. For diagrams (Draw.io, Mermaid, SVG), line numbers refer to the serialized text representation, not the original diagram coordinates.
2. **Component matching** is case-insensitive substring matching. Component labels that are substrings of other words may produce false matches.
3. **Confidence scoring** is capped at 95% maximum, even with perfect evidence. This prevents overconfidence in the deterministic engine.
4. **STRIDE justification confidence** uses a simplified formula (keyword match count + component match count) that does not distinguish between strong and weak keyword matches.
