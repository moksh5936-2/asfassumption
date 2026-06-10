# ASF Deterministic Risk Model

## Version

`asf-risk-model-1.0`

## Overview

ASF uses a standard 5×5 risk matrix with separate likelihood and impact analysis. Every score is deterministic and reproducible. No AI or randomization is involved.

## Likelihood Scale (1–5)

| Score | Label | Description |
|-------|-------|-------------|
| 1 | Very Low | Internal component, no external exposure, authenticated access required |
| 2 | Low | Internal with limited network exposure, standard authentication |
| 3 | Medium | Network-accessible component, moderate attack surface |
| 4 | High | Internet-exposed, authentication-dependent, multiple attack paths |
| 5 | Very High | Directly internet-facing, no authentication, known active threats |

### Likelihood Factors

Each assumption is evaluated on three factors:

#### 1. Exposure Level (base: 1)
| Condition | Score |
|-----------|-------|
| Internal, no external exposure | 1 |
| Network-accessible component | 3 |
| Internet-exposed (detected via gateway/internet components) | 5 |
| Network/internet category assumption | 4 |

#### 2. Authentication Dependency (base: 2)
| Condition | Score |
|-----------|-------|
| Standard security boundary | 2 |
| Authentication/authorization concepts detected | 4 |
| Auth/login/sso/mfa/identity component present | 4 |

#### 3. Attack Surface Complexity (base: 2)
| Condition | Score |
|-----------|-------|
| ≤3 components in evidence | 2 |
| >3 components in evidence | 3 |
| >5 relationships in evidence | 4 |

### Likelihood Score Calculation
```
score = clamp_to_1_5(1 + (exposure - 1) + (auth_dep - 2) + (complexity - 2))
```

## Impact Scale (1–5)

| Score | Label | Description |
|-------|-------|-------------|
| 1 | Very Low | No data loss, no regulatory impact, no business disruption |
| 2 | Low | Minor data exposure, no regulatory penalties |
| 3 | Medium | Moderate data exposure, potential compliance findings |
| 4 | High | Sensitive data exposure, regulatory penalties possible |
| 5 | Very High | PII/PHI/financial data exposure, guaranteed regulatory action, business-critical |

### Impact Factors

Each assumption is evaluated on three factors:

#### 1. Data Classification (base: 2)
| Condition | Score |
|-----------|-------|
| Standard business data | 2 |
| Data protection concepts detected | 4 |
| Database/storage/backup components | 4 |
| PII/PHI/financial/payment components | 5 |

#### 2. Regulatory Exposure (base: 1)
| Condition | Score |
|-----------|-------|
| No regulatory keywords detected | 1 |
| Healthcare (HIPAA) keywords | 5 |
| Payment (PCI DSS) keywords | 5 |
| Financial/SOX keywords | 4 |
| GDPR/privacy keywords | 4 |

#### 3. Business Criticality (base: 2)
| Condition | Score |
|-----------|-------|
| Standard component | 2 |
| Core/main/primary component | 4 |
| >8 relationships (cascading failure risk) | 4 |

### Impact Score Calculation
```
score = clamp_to_1_5(1 + (data_class - 2) + (regulatory - 1) + (criticality - 2))
```

## Risk Matrix (5 × 5)

```
         | Impact 1  |  2      |  3      |  4      |  5      |
---------|-----------|---------|---------|---------|---------|
Like 5   |   5       | 10      | 15      | 20      | 25      |
         |   Low     | Medium  | High    | Critical| Critical|
---------|-----------|---------|---------|---------|---------|
Like 4   |   4       |  8      | 12      | 16      | 20      |
         |   Low     | Medium  | High    | High    | Critical|
---------|-----------|---------|---------|---------|---------|
Like 3   |   3       |  6      |  9      | 12      | 15      |
         |   Low     | Medium  | Medium  | High    | High    |
---------|-----------|---------|---------|---------|---------|
Like 2   |   2       |  4      |  6      |  8      | 10      |
         |   Low     | Low     | Medium  | Medium  | Medium  |
---------|-----------|---------|---------|---------|---------|
Like 1   |   1       |  2      |  3      |  4      |  5      |
         |   Low     | Low     | Low     | Low     | Medium  |
```

### Risk Score = Likelihood × Impact

### Risk Level Thresholds
| Score Range | Level |
|-------------|-------|
| 20–25 | **Critical** |
| 12–19 | **High** |
| 5–11 | **Medium** |
| 1–4 | **Low** |

## Reproducibility

Every risk score is fully deterministic:

1. Same architecture input → same components/relationships extracted
2. Same Python ASF output → same assumption categories and keywords
3. Same evidence matches → same likelihood and impact factors
4. Same likelihood × impact → same risk score and level

No random seeds. No AI inference. No API calls. All logic is in:
- `justify.go:LikelihoodAnalyzer`
- `justify.go:ImpactAnalyzer`
- `justify.go:RiskMatrix`

## Confidence Scoring

Overall assumption confidence is calculated from four metrics:

| Metric | Weight (max) | Formula |
|--------|-------------|---------|
| Evidence points | 0.30 | min(0.05 × count, 0.30) |
| STRIDE rule matches | 0.25 | min(0.08 × matches, 0.25) |
| Component matches | 0.20 | min(0.06 × matches, 0.20) |
| Relationship matches | 0.15 | min(0.04 × matches, 0.15) |

```
Base: 0.10
Max:  0.95
```

**Interpretation**:
| Confidence | Interpretation |
|------------|----------------|
| < 0.5 | Low — limited evidence |
| 0.5–0.7 | Medium — moderate evidence |
| 0.7–0.85 | High — strong evidence |
| 0.85–0.95 | Very High — multiple corroborating evidence sources |
