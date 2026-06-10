# ASF Defensibility Analysis

> Version: 1.0.0 | June 2026 | Classification: Internal

## 1. Executive Summary

This document analyzes the defensibility of ASF v1.0.0 — evaluating whether the architecture, dependency choices, code quality, and engineering practices can withstand scrutiny from security professionals, auditors, and enterprise customers.

**Overall Rating: 6/10** — The core architecture is sound and well-documented, but significant gaps in validation, testing, and automation undermine defensibility.

---

## 2. Architecture Defensibility

### Strengths

| Decision | Rationale | Defensible? |
|----------|-----------|-------------|
| Go language | Memory-safe, statically linked, single binary | ✅ Yes |
| Bubble Tea TUI | Mature, well-maintained, active community | ✅ Yes |
| Deterministic analysis | No AI randomness, fully reproducible | ✅ Yes |
| Offline-first | No data exfiltration risk, no network dependency | ✅ Yes |
| HMAC license validation | No phone-home, no external dependency | ✅ Yes |
| YAML config | Human-readable, well-defined format | ✅ Yes |
| 5-format export | Open formats, no proprietary output | ✅ Yes |

### Weaknesses

| Decision | Risk | Defensible? |
|----------|------|-------------|
| Python subprocess for extraction | External dependency, version mismatch risk | ❌ Weak |
| Hardcoded HMAC secret | Extractable from binary | ❌ Weak |
| No input size limits | Potential DoS via large files | ❌ Weak |
| Single-threaded parsing | Performance bottleneck | ⚠️ Acceptable |

### Key Architectural Risks

#### Risk 1: Python Subprocess Dependency

```
ASF ──exec.Command("python3", "-m", "asf.cli.main")──▶ Python Engine
```

**Problem:** The Python ASF engine is a completely separate codebase with its own dependencies, versioning, and potential failure modes. If `python3` is not installed or the ASF package is not installed, ASF produces a cryptic error or silently falls back.

**Defensibility Impact:** Low — this is a well-known architectural pattern (subprocess delegation). The risk is documented in `INSTALLATION_ARCHITECTURE.md` and `ARCHITECTURE.md`.

**Mitigation:** Document the dependency clearly. Add pre-flight checks at startup. Consider a native Go extraction engine for v2.0.

#### Risk 2: Determinism vs. Accuracy

**Claim:** "Every result is reproducible and auditable."

**Defensibility:** Strong — the code is deterministic by design. No random number generators, no AI heuristics, no stochastic processes. Given the same input, ASF always produces the same output.

**Mitigation:** Publish the exact algorithm specifications in `risk_model.md` and `STRIDE_SPECIFICATION.md` (proposed).

#### Risk 3: No Empirical Validation

**Problem:** ASF cannot currently cite precision, recall, FPR, or STRIDE accuracy metrics. An auditor or enterprise customer will reasonably ask "how do we know this works?" — and there is no answer.

**Defensibility Impact:** High — this is the single biggest defensibility gap. Without validation data, every claim about ASF's effectiveness is unsubstantiated.

**Mitigation:** The `EXPERT_VALIDATION_STUDY.md` outlines a rigorous study plan (Phase 7). This must be executed before any enterprise sale.

---

## 3. Dependency Defensibility

### Go Dependencies

| Dependency | Version | Maintained? | License | Risk |
|------------|---------|-------------|---------|------|
| bubbletea | 1.3.10 | ✅ Active | MIT | 🟢 Low |
| lipgloss | 1.1.0 | ✅ Active | MIT | 🟢 Low |
| go-pdf/fpdf | 0.9.0 | ⚠️ Stable (low activity) | MIT | 🟢 Low |
| yaml.v3 | 3.0.1 | ✅ Active | MIT/Apache-2.0 | 🟢 Low |
| 12 transitive | Various | ✅ All active/stable | MIT/BSD | 🟢 Low |

### Dependency Tree Analysis

```
asf-tui (Go 1.24.0)
├── github.com/charmbracelet/bubbletea v1.3.10
│   ├── github.com/charmbracelet/x/term v0.2.1
│   ├── github.com/erikgeiser/coninput v0.0.0-20211004153227
│   ├── github.com/mattn/go-localereader v0.0.1
│   ├── github.com/muesli/ansi v0.0.0-20230316100256
│   ├── github.com/muesli/cancelreader v0.2.2
│   └── golang.org/x/sys v0.36.0
├── github.com/charmbracelet/lipgloss v1.1.0
│   ├── github.com/charmbracelet/colorprofile v0.2.3-0.20250311203215
│   ├── github.com/charmbracelet/x/ansi v0.10.1
│   ├── github.com/charmbracelet/x/cellbuf v0.0.13-0.20250311204145
│   ├── github.com/lucasb-eyer/go-colorful v1.2.0
│   ├── github.com/mattn/go-isatty v0.0.20
│   ├── github.com/mattn/go-runewidth v0.0.16
│   ├── github.com/muesli/termenv v0.16.0
│   ├── github.com/rivo/uniseg v0.4.7
│   └── github.com/xo/terminfo v0.0.0-20220910002029
├── github.com/go-pdf/fpdf v0.9.0
└── gopkg.in/yaml.v3 v3.0.1
```

**Depth:** Maximum 2 levels of transitive dependencies.
**Total:** 4 direct + 12 transitive = 16 unique dependencies.
**Risk:** Low — all dependencies are from trusted, well-maintained projects.

### Supply Chain Defensibility

| Practice | Status | Defensible? |
|----------|--------|-------------|
| `go.sum` verification | ✅ In place | ✅ Yes |
| Version pinning | ✅ In place | ✅ Yes |
| Vendor directory | ❌ Not used | ⚠️ Weak (requires network for build) |
| SBOM generation | ❌ Not done | ❌ Weak |
| Dependency scanning | ❌ Not configured | ❌ Weak |
| Binary provenance | ❌ Not signed | ❌ Weak |

**Recommendation:** Add `.goreleaser` configuration with SBOM generation and binary signing for production releases.

---

## 4. Testing Defensibility

### Current Coverage

| Package | Tests | Coverage (est.) |
|---------|-------|----------------|
| justify.go (explain_test.go) | 20 | ~80% of explainability code |
| All other packages | 0 | ~0% |
| **Total** | **20** | **~15%** |

### Critical Untested Code

| File | Lines | Risk |
|------|-------|------|
| `parser.go` | ~200 | 🟡 Medium — format-specific parsing, edge cases |
| `stride.go` | ~200 | 🟡 Medium — STRIDE mapping correctness |
| `export.go` | ~300 | 🟡 Medium — export format correctness |
| `engine.go` | ~400 | 🔴 High — Python bridge, result building |
| `app.go` | ~300 | 🟡 Medium — TUI navigation, state management |
| `review.go` | ~200 | 🟡 Medium — review workflow correctness |
| `license.go` | ~150 | 🟡 Medium — license validation logic |

**Defensibility Impact:** Medium-High. An auditor would reasonably expect unit tests for the core analysis pipeline, STRIDE mapping, and export formats.

**Mitigation:** Prioritize tests for `stride.go` and `export.go` before v1.0.0 final.

---

## 5. Documentation Defensibility

| Document | Completeness | Accuracy | Honesty | Defensible? |
|----------|-------------|----------|---------|-------------|
| ARCHITECTURE.md | ✅ High | ✅ Current | ✅ Honest | ✅ Yes |
| EXECUTIVE_SUMMARY.md | ✅ High | ✅ Current | ✅ Honest | ✅ Yes |
| USER_MANUAL.md | ✅ High | ✅ Current | ✅ Honest | ✅ Yes |
| DEVELOPER_GUIDE.md | ✅ High | ✅ Current | ✅ Honest | ✅ Yes |
| TECHNICAL_REFERENCE.md | ✅ High | ✅ Current | ✅ Honest | ✅ Yes |
| VALIDATION_STATUS.md | ✅ High | ✅ Current | ✅ Brutally honest | ✅ Yes |
| README.md | ✅ High | ✅ Current | ✅ Honest | ✅ Yes |
| SECURITY_REVIEW.md | ✅ New | ✅ Current | ✅ Honest | ✅ Yes |
| DEFENSIBILITY_ANALYSIS.md | ✅ New | ✅ Current | ✅ Honest | ✅ Yes (this document) |

**Documentation is a strength.** All documents clearly state limitations, unknowns, and risks. This honesty is itself a defensibility asset — it demonstrates awareness and transparency.

---

## 6. Process Defensibility

### Current Process

```
Code → Manual build → Manual test → Manual release
```

**Problems:**
- No CI/CD pipeline — every build is a manual `go build`
- No automated testing — "tests pass" means "I ran them once"
- No release automation — cut-and-release is manual
- No changelog enforcement — conventional commits not required

### Required Process for Enterprise

```
Code → PR review → CI (lint, vet, test, build) → Staging → QA → Release (goreleaser)
```

**Gap:** All enterprise process requirements are missing.

---

## 7. Legal & Compliance Defensibility

| Aspect | Status | Defensible? |
|--------|--------|-------------|
| License (code) | ⚠️ "Research and educational use" | ❌ Weak — not a standard OSI license |
| Enterprise license | ✅ HMAC-signed keys | ✅ Yes |
| Open source dependencies | ✅ All MIT/BSD/Apache-2.0 | ✅ Yes |
| Contribution agreement | ❌ None (single contributor) | ⚠️ Acceptable for now |
| Privacy policy | ❌ Not published | ❌ Weak for enterprises |
| Terms of service | ❌ Not published | ❌ Weak for enterprises |

---

## 8. Competitor Defensibility Comparison

| Attribute | ASF | MS TMT | IriusRisk | Threat Dragon |
|-----------|-----|--------|-----------|---------------|
| Open source | ✅ | ❌ | ❌ | ✅ |
| Deterministic | ✅ | ⚠️ Manual | ⚠️ Manual | ⚠️ Manual |
| Offline | ✅ | ❌ | ❌ | ✅ |
| Automated | ✅ | ❌ | ⚠️ Partial | ❌ |
| Tested | ❌ | Unknown | Unknown | Unknown |
| Validated | ❌ | Unknown | Unknown | ❌ |
| Documented | ✅ | ✅ | ✅ | ⚠️ |

**Key Insight:** ASF's open-source, deterministic, offline, automated nature gives it a defensibility advantage over proprietary competitors — even with its testing/validation gaps.

---

## 9. Recommendations by Priority

### Must Fix (Defensibility Blockers)

1. **Run expert validation study** — Without metrics, ASF cannot be taken seriously
2. **Add CI/CD (GitHub Actions)** — Automated builds and tests are table stakes
3. **Increase test coverage** — Minimum: stride.go, export.go, engine.go
4. **Standardize license** — Choose OSI-approved license for open source release

### Should Fix (Defensibility Weaknesses)

5. Add `go vet` and `golangci-lint` to CI
6. Add SBOM generation to release workflow
7. Add macOS code signing
8. Add input size limits to file parsers
9. Create contributor agreement template

### Nice to Fix (Defensibility Enhancements)

10. Add audit logging
11. Add binary provenance (SLSA Level 1+)
12. Add vulnerability scanning (Dependabot)
13. Publish privacy policy and terms of service

---

## 10. Conclusion

ASF v1.0.0 has a defensible core architecture but lacks the empirical validation and process maturity needed for enterprise adoption. The documentation is excellent and honestly communicates limitations.

**Defensibility Score: 6/10**

| Category | Score | Notes |
|----------|-------|-------|
| Architecture | 8/10 | Sound design, Python dependency is main weakness |
| Dependencies | 8/10 | Well-maintained, low risk |
| Testing | 3/10 | 20 tests for explainability, nothing else |
| Documentation | 9/10 | Comprehensive and honest |
| Process | 2/10 | No CI/CD, no automation |
| Legal | 4/10 | Non-standard license, no privacy/terms |
| **Overall** | **6/10** | |

**Bottom Line:** ASF's technical defensibility (determinism, offline-first, Go, clean architecture) is strong. Its process defensibility (testing, CI/CD, validation) is weak. The gap is fixable but requires dedicated effort.

---

*Generated by `docs/DEFENSIBILITY_ANALYSIS.md` — ASF Release Engineering Phase 11*
