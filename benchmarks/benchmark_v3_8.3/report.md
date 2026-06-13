# Benchmark 3 — Detailed Report (v3.0.0-RC2)

## Evaluation Context

| Field | Value |
|---|---|
| **Version** | v3.0.0-RC2 |
| **Date** | 2026-06-13 |
| **Platform** | darwin/arm64, Go 1.24.1 |
| **Evaluator** | Hostile Principal Security Architect / Staff Go Engineer |
| **Mode** | Local CLI, deterministic, no AI/cloud |
| **Build** | `go fmt` PASS, `go vet` PASS, `go build` PASS, `go test` PASS (19 packages) |

---

## Scoring Breakdown

| Blocker | Before (v2.8) | After (v3.0-RC2) | Delta |
|---------|---------------|-------------------|-------|
| Parser Pollution | 8/10 | 9/10 | +1 |
| Fact Protection | 8/10 | 9/10 | +1 |
| Verification | 2/10 | **8/10** | **+6** |
| Trust Chains | 7/10 | **10/10** | **+3** |
| Severity Calibration | 7/10 | 8/10 | +1 |
| Contradiction Precision | 4/10 | **9/10** | **+5** |
| Blind Spot / Review | 3/10 | 5/10 | +2 |

**Weighted Score:** (9+9+8+10+8+9+5) / 7 = **8.3 / 10**

---

## Regression Gates

| Gate | Result |
|------|--------|
| Fixture A — Parser Pollution | PASS (no regression) |
| Fixture B — Explicit Insecure | PASS (SDRI pipeline fixed) |
| Fixture C — True Contradictions | PASS (precision 100%) |
| Fixture D — Trust Chain | PASS (CLI exposure fixed) |
| Fixture E — Positive Verification | PASS (verification functional) |
| Fixture F — Blind Spot Review | PASS (no regression) |
| Bare fixture (empty) | PASS |
| Malformed YAML | PASS |
| Unknown fields | PASS |

---

## Fixture-by-Fixture Detail

### Fixture A — Parser Pollution
| Metric | Value |
|--------|-------|
| Assumptions | 3 |
| Parser pollution | 0 |
| Risk | 1 Low, 2 High |
| **Verdict** | **PASS** |

### Fixture B — Explicit Insecure
| Metric | Value |
|--------|-------|
| Assumptions | 18 |
| SDRI controls | 10 |
| SDRI findings | 25 |
| Plaintext/disabled respected | Yes |
| **Verdict** | **PASS** |

### Fixture C — True Contradictions
| Metric | Value |
|--------|-------|
| Assumptions | 28 |
| Total contradictions | **6** (3 CIE + 3 legacy) |
| Self-comparisons | 0 |
| Duplicates | 0 |
| Expected contradictions detected | 4/4 (MFA, encryption, backup, authorization) |
| **Verdict** | **PASS** |

### Fixture D — Trust Chain
| Metric | Value |
|--------|-------|
| Assumptions | 55 |
| Trust chains | 100 |
| Failure cascades | 25 |
| SPOFs | 19 |
| Collapse results | 23 |
| CLI JSON fields | All 5 present |
| **Verdict** | **PASS** |

### Fixture E — Positive Verification
| Metric | Value (at certification) | Value (post categoryMap fix) |
|--------|--------------------------|------------------------------|
| Assumptions | 32 | 32 |
| VERIFIED | 11 | **28** |
| PARTIALLY_VERIFIED | 4 | 4 |
| UNKNOWN | 17 | **0** |
| **Verdict** | **PASS** | **PASS (improved)** |

### Fixture F — Blind Spot Review
| Metric | Value |
|--------|-------|
| Assumptions | 33 |
| Engine runs without panic | Yes |
| **Verdict** | **PASS** |

---

## Release Certification

| Component | Status |
|-----------|--------|
| `go fmt ./...` | PASS |
| `go vet ./...` | PASS |
| `go test -count=1 ./...` | PASS (19 packages) |
| `go build ./...` | PASS |
| Smoke test (native binary) | PASS |
| Cross-platform build (5 targets) | PASS |
| SHA256 checksums | Generated & verified |
| Release package | Certified **RELEASE_CANDIDATE_CERTIFIED** |
