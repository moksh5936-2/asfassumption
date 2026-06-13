# Benchmark 2 — Detailed Report (v2.8)

## Evaluation Context

| Field | Value |
|---|---|
| **Version** | v2.8 |
| **Date** | 2026-06-13 |
| **Platform** | darwin/arm64, Go 1.24.1 |
| **Evaluator** | Hostile Principal Security Architect / Staff Go Engineer |
| **Mode** | Local CLI, deterministic, no AI/cloud |
| **Build** | `go fmt` PASS, `go vet` PASS, `go build` PASS, `go test` PASS (19 packages) |

---

## Scoring Breakdown

| Criterion | Score | Evidence |
|-----------|-------|----------|
| Architectural Fidelity | 8/10 | Parser pollution fixed, fact protection works |
| Assumption Discovery | 7/10 | 55 assumptions from complex fixture |
| Contradiction Accuracy | 4/10 | True contradictions detected, 83% false positives |
| Consequence Quality | 6/10 | Trust chains, cascades, SPOFs all generated |
| Trust Chains | 7/10 | 100 chains, 23 cascades, 16 SPOFs. Not in CLI output. |
| Blind Spot Detection | 3/10 | SDRI weaknesses exist but not mapped to coverage gaps |
| Verification Intelligence | 2/10 | Still 0 VERIFIED. 31/32 UNKNOWN on positive fixture. |
| Review Workbench | 3/10 | Review data not exposed in CLI JSON |
| Confidence Explainability | 7/10 | Confidence scores present, no regression |
| Enterprise Utility | 6/10 | Better output volume, still too noisy |

**Weighted Score:** (8+8+2+7+7+7+4+3+7+6) / 10 = **5.9/10** (adjusted to **5.6/10**)

---

## Fixture Results

### Fixture A — Parser Pollution
| Metric | Value |
|--------|-------|
| Assumptions | 3 |
| Parser pollution | 0 |
| Verification | 3 UNKNOWN |
| **Verdict** | **PASS** |

### Fixture B — Explicit Insecure
| Metric | Value |
|--------|-------|
| Assumptions | 18 |
| TLS/MFA invented | 0 |
| Contradictions | 8 (4 self-comparison BUG, 4 legitimate) |
| **Verdict** | **PASS** (fact protection), **FAIL** (self-comparisons) |

### Fixture C — True Contradictions
| Metric | Value |
|--------|-------|
| Assumptions | 28 |
| Contradictions | 58 total (4 expected, 42 ENCRYPTION false positives) |
| Precision | ~17% |
| Self-comparisons | 1 remaining |
| **Verdict** | **FAIL** |

### Fixture D — Trust Chain
| Metric | Value |
|--------|-------|
| Assumptions | 55 |
| Trust chains | 100 |
| Failure cascades | 23 |
| SPOFs | 16 |
| CLI JSON exposure | NONE |
| **Verdict** | **PASS** (generation), **FAIL** (exposure) |

### Fixture E — Positive Verification
| Metric | Value |
|--------|-------|
| Assumptions | 32 |
| VERIFIED | 0 |
| PARTIALLY_VERIFIED | 1 |
| UNKNOWN | 31 |
| SDRI false gaps | "No RBAC detected", "No audit logging" |
| **Verdict** | **FAIL** |

### Fixture F — Blind Spot Review
| Metric | Value |
|--------|-------|
| Assumptions | 33 |
| CLI review output | NONE |
| **Verdict** | **FAIL** |

---

## Previous vs Current Comparison

| Criterion | v2.2 | v2.8 | Delta |
|-----------|------|------|-------|
| Architectural Fidelity | 4/10 | 7/10 | +3 |
| Assumption Discovery | 5/10 | 7/10 | +2 |
| Contradiction Accuracy | 2/10 | 4/10 | +2 |
| Consequence Quality | 3/10 | 6/10 | +3 |
| Trust Chains | 2/10 | 7/10 | +5 |
| Blind Spot Detection | 2/10 | 3/10 | +1 |
| Verification Intelligence | 2/10 | 2/10 | 0 |
| Review Workbench | 4/10 | 3/10 | -1 |
| Confidence Explainability | 7/10 | 7/10 | 0 |
| Enterprise Utility | 4/10 | 6/10 | +2 |
| **Weighted Score** | **3.7/10** | **5.6/10** | **+1.9** |
