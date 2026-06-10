# Confidence Parity Report

**Phase 5 — June 2026**

## Summary

The Go and Python engines use **different confidence scoring algorithms by design**. Numerical confidence values diverge systematically (~0.27 avg absolute difference), but this does **not** affect any analytic output — verification statuses, gap types, gap severities, and claim/assumption counts are all identical.

## Methodology

- Compared confidence values across all 24 shared-format samples in `testdata/parity/go/` vs `testdata/parity/python/`
- Fields compared: `assumptions[*].confidence`, `verifications[*].confidence`
- Tolerance: ±0.001
- Assumptions matched by whitespace-normalized text content; verifications matched by array index (same count/order confirmed)

## Results

| Metric | Value |
|--------|-------|
| Total confidence values compared | 164 |
| Matches within 0.001 | 36 (22.0%) |
| Average per-file absolute diff | 0.2720 |
| Maximum absolute diff | 0.7507 |
| Zero-confidence (UNKNOWN) values | 100% match (both return 0.0) |

## Algorithm Differences

### Verification Confidence

| Component | Python (Additive Weighted Average) | Go (Multiplicative) |
|-----------|-----------------------------------|---------------------|
| **Formula** | `base*0.40 + freshness*0.20 + coverage*0.20 + completeness*0.20` | `base * (freshness*0.30 + coverage*0.40 + completeness*0.30)` |
| **Freshness** | Continuous linear decay: `max(0, 1 - age_hours/720)` | Discrete banded decay: <24h→1.0, <7d→0.9, <30d→0.7, <90d→0.5, <365d→0.3, else→0.1 |
| **Coverage** | Simple ratio: `used_evidence / total_evidence` (can be 0.0) | `0.2 + ratio*0.8` (floor at 0.2, even with no evidence) |
| **Completeness** | Result-based: VERIFIED→1.0, PARTIAL→0.5, CONTRADICTED→1.0, UNKNOWN→0.0 | NLP heuristic: `0.3 + (indicator_count/8)*0.7` based on reasoning text keywords |
| **UNKNOWN special case** | Applies full weighted formula with all factors | Returns base confidence as-is, skips all computation |

### Assumption Confidence

| Aspect | Python | Go |
|--------|--------|----|
| **Formula** | `avg(v.confidence * result_multiplier)` | `avg(v.confidence)` — simple mean |
| **Result multipliers** | VERIFIED×1.0, PARTIAL×0.5, CONTRADICTED×0.1, UNKNOWN×0.0 | None (all verifications weighted equally) |
| **Effect** | CONTRADICTED assumptions heavily penalized (×0.1) | CONTRADICTED and VERIFIED contribute equally |

### Extraction Confidence (Claim Level)

Both engines use the **identical** formula: start at 0.5, +0.1 per strong indicator pattern matched, capped at 0.95. **No difference.**

## Why This Is Acceptable

1. **Structural parity is maintained** — verification results (CONTRADICTED, UNKNOWN, VERIFIED, PARTIALLY_VERIFIED) are identical, which drives gap severity assignment
2. **Gap severities are identical** — the gap engine uses confidence thresholds (≥0.8, ≥0.5) calibrated per-engine, producing identical classifications
3. **Confidence is advisory** — it is documented as a secondary metric, not a primary analytic output
4. **Monotonic ordering preserved** — both formulas produce higher confidence for more/better evidence

## Conclusion

The confidence algorithms are **different by design**. The divergence is documented, understood, and accepted. Confidence values should not be compared numerically between engines; only relative ordering and the derived gap severities are meaningful.
