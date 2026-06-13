# Benchmark 1 — Initial Baseline (v2.2)

## Executive Summary

ASF v2.2 showed foundational weaknesses across all critical evaluation criteria. The verification engine was non-functional (returning 0 VERIFIED on any input), the parser leaked YAML structure as assumptions, trust chain generation was absent, contradiction detection did not exist, and recall was approximately 15%. The engine was a topology linter rather than a security assumption discovery platform.

---

## Score

**3.7 / 10**

## Verdict

**NOT_READY**

---

## Major Findings

### Top Strengths

1. **Pipeline execution** — All 5 benchmark datasets execute without crashes or panics
2. **Basic assumption generation** — 48 assumptions generated (though low quality)
3. **Deterministic output** — No AI/cloud dependency, fully reproducible results

### Top Weaknesses

1. **Parser pollution** — YAML comments, headers, and filenames leaked into assumption text
2. **Verification failure** — 0 VERIFIED statuses produced for any input
3. **Trust chain failure** — No trust chain generation
4. **Contradiction failure** — 0 contradictions detected
5. **Low recall** — ~15%, missing the majority of architectural insights
6. **Only 8 assumption categories** — limited domain coverage

### Blocking Issues

| Issue | Severity | Impact |
|-------|----------|--------|
| Parser pollution | Critical | Assumptions contain irrelevant text, eroding trust |
| Verification non-functional | Critical | Cannot determine if a design is recognized as secure |
| Trust chain absent | Critical | No identity/crypto trust root analysis |
| Contradiction absent | High | No inconsistency detection |
| Low recall | High | Misses >80% of valid architectural concerns |

### Actions Taken

- Full intelligence engine overhaul (Go native rewrite)
- Parser hardening to filter non-architectural text
- Trust chain discovery engine implemented
- Contradiction detection engine implemented
- Assumption category expansion from 8 to 24
- Fact protection logic added

### Outcome

The v2.2 baseline established that ASF could execute a pipeline but could not produce trustworthy security analysis. The 3.7/10 score triggered a full engine overhaul that became the hardening sprint to v2.8.

---

## Source Documents

- `docs/BENCHMARK_REPORT.md` — Initial benchmark validation
- `docs/BENCHMARK_IMPROVEMENT_REPORT.md` — Before/after improvement analysis
- `asf-tui/docs/V29_TARGETED_REBENCHMARK_REPORT.md` — Section 6 (previous vs current comparison)
