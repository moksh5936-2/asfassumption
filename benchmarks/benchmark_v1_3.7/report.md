# Benchmark 1 — Detailed Report (v2.2)

## Evaluation Context

| Field | Value |
|---|---|
| **Version** | v2.2 |
| **Date** | 2026-06-13 |
| **Platform** | darwin/arm64 |
| **Evaluator** | Hostile Principal Security Architect / Staff Go Engineer |
| **Mode** | Local CLI, deterministic, no AI/cloud |

---

## Scoring Breakdown

| Criterion | Score | Evidence |
|-----------|-------|----------|
| Architectural Fidelity | 4/10 | Parser leaks, no fact protection |
| Assumption Discovery | 5/10 | 48 assumptions, 8 categories |
| Contradiction Accuracy | 2/10 | 0 contradictions detected |
| Consequence Quality | 3/10 | No trust chains, no cascades |
| Trust Chains | 2/10 | Not implemented |
| Blind Spot Detection | 2/10 | Not implemented |
| Verification Intelligence | 2/10 | 0 VERIFIED on any input |
| Review Workbench | 4/10 | CLI output minimal |
| Confidence Explainability | 7/10 | Confidence scores present |
| Enterprise Utility | 4/10 | Too noisy for direct use |

**Weighted Score:** (4+5+2+3+2+2+2+4+7+4) / 10 = **3.7 / 10**

---

## Fixture Results

| Fixture | Result |
|---------|--------|
| Healthcare PHI | 48 assumptions, 3 Critical, 4 High, 9 Medium, 32 Low |
| — Recall | ~15% |
| — Precision | ~40-60% |
| — Categories | 8 |
| — Contradictions | 0 |
| — Trust Boundaries | 0 |

## Original Objective

The v2.2 benchmark was conducted as a baseline to establish whether ASF could:
1. Execute the analysis pipeline without crashes
2. Generate security assumptions from architecture YAML
3. Detect contradictions
4. Identify trust boundaries
5. Produce actionable output for security architects

Only criterion 1 and 2 (partially) were met.

## Certifications Achieved

- **Pipeline execution** — PASS (basic)

## Certifications Failed

- **Intelligence Engine** — FAIL (3.7/10 baseline)
- **Verification** — FAIL (0 VERIFIED)
- **Contradiction** — FAIL (0 detected)
- **Trust Chain** — FAIL (absent)
