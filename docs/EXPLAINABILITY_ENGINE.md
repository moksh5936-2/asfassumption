# ASF Explainability Engine

## Overview

The Explainability Engine transforms ASF from a security assumption generator into an explainable security review engine. Every output — assumption, STRIDE category, risk score, confidence level — includes traceable evidence and human-readable reasoning.

## Architecture

```
Architecture Source
    │
    ▼
ParseArchitecture()
    │
    ├── Components [name, type]
    ├── Relationships [source, target, label]
    └── RawText
            │
            ▼
    Python ASF CLI
    (assumption extraction)
            │
            ▼
    ┌─────────────────────────────┐
    │  ExplainabilityPipeline     │
    │                             │
    │  1. EvidenceEngine          │
    │     - Match components      │
    │     - Match relationships   │
    │     - Detect trust bounds   │
    │     - Identify concepts     │
    │                             │
    │  2. AssumptionJustifier     │
    │     - Build rationale       │
    │     - Link to sources       │
    │                             │
    │  3. StrideJustifyEngine     │
    │     - Per-category reason   │
    │     - Rule index tracking   │
    │     - Confidence calc       │
    │                             │
    │  4. LikelihoodAnalyzer      │
    │     - Exposure              │
    │     - Auth dependency       │
    │     - Attack complexity     │
    │                             │
    │  5. ImpactAnalyzer          │
    │     - Data classification   │
    │     - Regulatory exposure   │
    │     - Business criticality  │
    │                             │
    │  6. RiskMatrix              │
    │     - L × I = RiskScore     │
    │     - 5×5 matrix            │
    │                             │
    │  7. ConfidenceEngine        │
    │     - Evidence count        │
    │     - Rule matches          │
    │     - Component matches     │
    └─────────────────────────────┘
            │
            ▼
    Explained Assumption
    ├── Description
    ├── Evidence Sources
    ├── Source Components
    ├── Source Relationships
    ├── Rationale
    ├── STRIDE Justifications
    │   ├── Category
    │   ├── Reason
    │   ├── Matched Rules
    │   ├── Matched Keywords
    │   └── Confidence
    ├── Risk Justification
    │   ├── Likelihood (1-5)
    │   ├── Likelihood Factors
    │   ├── Impact (1-5)
    │   ├── Impact Factors
    │   ├── Risk Score (1-25)
    │   ├── Risk Level
    │   └── Confidence
    ├── Confidence
    ├── Review Status
    └── Review Notes
```

## Component Details

### 1. Evidence Engine (`justify.go:EvidenceEngine`)

Traces every assumption back to the architecture source.

**Input**: `ArchDescription` (components, relationships) + assumption (category, keywords, text)

**Matching Logic**:
| Artifact | Match Method | Example |
|---|---|---|
| Components | Substring match on lowercase label | `"Database"` in assumptions ↔ `Component{Label:"Database"}` |
| Relationships | Source or target label match | `"API Gateway → Database"` matched if either appears in assumption text |
| Trust Boundaries | Contains `internet/external/public/vpn/gateway/dmz` | Component `"Internet"` with relationship to `"API Gateway"` |
| Security Concepts | Keyword groups | authentication, authorization, encryption, network_security, data_protection, audit_logging, dependency, session_management |

**Output**: `EvidenceResult` with matched component/relationship/boundary/concept lists, total count.

### 2. Assumption Justification (`justify.go:JustifyAssumption`)

Generates deterministic, human-readable rationale.

**Logic**:
- If components matched: "detected N relevant component(s): [A, B, C]"
- If relationships matched: "identified N communication path(s) between components"
- If trust boundaries crossed: "crosses N trust boundary/boundaries requiring security verification"
- If security concepts matched: "relates to security concept(s): [X, Y, Z]"
- Fallback: "generated from category CATEGORY with N keyword match(es)"

**No AI. No LLM. Deterministic.**

### 3. STRIDE Justification (`justify.go:StrideJustifyEngine`)

Extends `StrideEngine` with per-category justification.

**Output per category**:
| Field | Source |
|---|---|
| `Category` | From `StrideEngine.MapAssumption()` |
| `Reason` | Template per category + matched keywords |
| `MatchedRuleIndexes` | Indexes into the 33 keyword rules in `stride.go` |
| `MatchedKeywords` | The specific keywords/text that triggered the match |
| `Confidence` | 0.3 base + 0.1 per keyword match (max 0.4) + 0.08 per component (max 0.3), capped at 0.95 |

**Confidence formula**:
```
confidence = 0.3 + min(0.1 × keywordMatches, 0.4) + min(0.08 × componentMatches, 0.3)
```

### 4. Likelihood Analysis (`justify.go:LikelihoodAnalyzer`)

Decomposes likelihood into three factors (each 1-5):

| Factor | Base | Modifiers |
|---|---|---|
| Exposure Level | 1 | +2 if network-accessible, +4 if internet-exposed, +3 if network/internet category |
| Auth Dependency | 2 | +2 if auth/authorization concepts matched, +2 if auth component present |
| Attack Complexity | 2 | +1 for >3 components, +2 for >5 relationships |

**Score**: Sum of (factor - base) for each factor, clamped to 1-5.

**Reason**: "likelihood X/5 based on exposure(N), auth dependency(N), attack complexity(N)"

### 5. Impact Analysis (`justify.go:ImpactAnalyzer`)

Decomposes impact into three factors (each 1-5):

| Factor | Base | Modifiers |
|---|---|---|
| Data Classification | 2 | +2 if data_protection concept present, +2 if database/storage/backup, +3 if PII/PHI/financial/payment |
| Regulatory Exposure | 1 | +4 for healthcare/PII, +3 for financial/SOX, +3 for GDPR/privacy |
| Business Criticality | 2 | +2 for core/main/primary components, +2 for >8 relationships |

**Score**: Sum of (factor - base) factors, clamped to 1-5.

### 6. Risk Matrix (`justify.go:RiskMatrix`)

Standard 5×5 matrix:

```
         | Impact 1 | 2 | 3 | 4 | 5 |
---------|---------|---|---|---|----|
Like 5   |    5     | 10| 15| 20| 25 |
Like 4   |    4     | 8 | 12| 16| 20 |
Like 3   |    3     | 6 | 9 | 12| 15 |
Like 2   |    2     | 4 | 6 | 8 | 10 |
Like 1   |    1     | 2 | 3 | 4 | 5  |
```

**Risk Level thresholds**:
- Critical: score >= 20
- High: score >= 12
- Medium: score >= 5
- Low: score < 5

### 7. Confidence Engine (`justify.go:ConfidenceEngine`)

Calculates overall confidence from available evidence:

```
confidence = 0.1
           + min(0.05 × evidenceCount, 0.3)
           + min(0.08 × strideRuleMatches, 0.25)
           + min(0.06 × componentMatches, 0.2)
           + min(0.04 × relationshipMatches, 0.15)
           capped at 0.95
```

### 8. Review Mode (`review.go`)

TUI-based architect review with status tracking.

**Statuses**:
- `Proposed` — default, not yet reviewed
- `Accepted` — architect confirms assumption is valid
- `Rejected` — architect considers assumption invalid
- `Modified` — architect modified the assumption

**Keyboard shortcuts**:
- `s` — Accept
- `r` — Reject
- `m` — Mark as Modified
- `n` — Toggle note
- `Enter` — Toggle detail view
- `↑↓` — Navigate assumptions

## Validation Data

`CollectValidationData()` exports all assumptions (with review results) as `[]ValidationRecord` for future studies:

- Precision: Accepted / (Accepted + Rejected)
- Recall: Novel findings / total expected
- False positive rate: Rejected / Total
- STRIDE accuracy: Agreement rate with expert labels
- Human agreement: Cohen's kappa between reviewers
