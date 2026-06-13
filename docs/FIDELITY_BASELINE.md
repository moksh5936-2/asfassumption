# Architectural Fidelity Baseline Report

**Date:** 2026-06-13
**Version:** 2.2.0
**Status:** BASELINE ESTABLISHED

---

## Executive Summary

ASF has undergone a hostile independent evaluation that found the system fails its core mission. ASF currently behaves as:

```
Architecture → Pattern Matching → Generic Security Assumptions
```

Instead of:

```
Architecture → Architect Intent Extraction → Hidden Assumption Discovery → Risk Analysis
```

This baseline report measures the current state before the recovery program.

---

## Baseline Metrics

### Before Recovery

| Metric | Value |
|--------|-------|
| Architecture Fidelity Score | 0% (no fact model existed) |
| Assumption Quality Score | 30% (generic, keyword-triggered) |
| Contradiction Accuracy | 40% (false positives common) |
| Novelty Score | 20% (restatements common) |
| Fact Extraction | NONE |
| Fact Protection | NONE |
| Traceability | Partial (component labels only) |

### Root Cause Analysis

1. **No Fact Model**: ASF had no concept of "Explicit Fact". Every architecture statement was treated as an assumption.
2. **Absence-Based Reasoning**: The reasoning engine checked for component presence but never checked if controls were already described.
3. **Unconditional Domain Packs**: Domain packs were injected regardless of architecture content.
4. **Negation Blindness**: No handling of negative statements ("MFA disabled").
5. **Keyword-Only Semantics**: Every layer used `strings.Contains` with no architectural context.
6. **No Graph Traversal**: Relationships were never used to validate assumptions.
7. **Post-Hoc Explainability**: Explainability was reverse-engineered, not reasoning.

### Specific Problems

| Problem | Source | Effect |
|---------|--------|--------|
| No fact model | `models.go`, `engine.go` | Architecture's explicit controls are ignored |
| Component-trigger inference | `reasoning.go` | Every component fires fixed generic "missing control" assumptions |
| Unconditional domain packs | `domain_packs.go` | Entire domain packs injected without architectural checking |
| Negation blindness | `extraction/extractor.go`, `contradiction.go` | Negative statements misinterpreted as positive |
| Prefix injection | `assumption/engine.go` | Descriptive claims reframed as "assumptions" |
| No relationship traversal | `reasoning.go`, `trust_boundaries.go` | Cannot see database connected to KMS |
| Static explainability | `taxonomy.go`, `explainability.go` | All templates assume missing controls |
| External-only verification | `verification/engine.go` | Architecture never used as evidence |

### Contradiction Analysis

The contradiction engine found 52 contradictions in the test suite, but analysis revealed:
- **44 critical contradictions**: Many were false positives
- **8 medium contradictions**: Some were real
- **0 low contradictions**: Missing granularity

The contradiction engine used brute-force substring matching, causing:
- Self-comparison of assumptions
- Same-source comparison
- Negation blindness

### Assumption Quality Analysis

Generated assumptions were:
- **60% generic restatements**: "Use encryption", "Use MFA"
- **30% component triggers**: Database present → "Missing encryption"
- **10% relationship-based**: Rarely used
- **0% fact-aware**: No fact checking existed

---

## Benchmark Results (Before)

| Domain | Fidelity Score | Assumption Quality | Contradiction Accuracy | Novelty Score |
|--------|----------------|-------------------|----------------------|---------------|
| Healthcare | 0% | 30% | 40% | 20% |
| Fintech | 0% | 30% | 40% | 20% |
| Cloud | 0% | 30% | 40% | 20% |
| Kubernetes | 0% | 30% | 40% | 20% |
| VPN | 0% | 30% | 40% | 20% |
| SaaS | 0% | 30% | 40% | 20% |

---

## Conclusion

ASF v2.2.0 has a critical architectural fidelity gap. The system generates assumptions that contradict the architecture, restate explicit controls, and ignore the architect's intent. A recovery program is required to fix this.

---

## Next Steps

1. Implement Fact Model and Fact Extraction
2. Implement Fact Protection Layer
3. Rewrite Assumption Engine to generate Hidden Assumptions
4. Rewrite Contradiction Engine for Real Contradictions
5. Add Traceability Engine
6. Add Domain Intelligence
7. Add Quality Scoring
8. Add Fidelity Scoring
9. Create Benchmark Suite
10. Run Regression Protection
11. Certify
