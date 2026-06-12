# ASF Validation Status

**Date:** June 2026

This document is a brutally honest assessment of what is validated and what is not. It is intended for researchers, evaluators, and anyone considering ASF for production use.

---

## What Is Validated

### Technical Validation

| Item | Status | Evidence |
|------|--------|----------|
| Go compilation | ✅ Passed | `go build` clean, binary 11.85MB |
| Static analysis | ✅ Passed | `go vet` clean, zero warnings |
| Unit tests | ✅ Passed | 20 tests covering risk matrix, confidence, evidence, STRIDE, pipeline |
| Risk matrix calculation | ✅ Validated | All 25 score combinations tested, boundaries verified |
| Confidence calculation | ✅ Validated | Deterministic, capped at 0.95, all inputs tested |
| STRIDE mapping | ✅ Validated | 17 category rules + 34 keyword patterns, deterministic |
| Evidence engine | ✅ Validated | Component/relationship/concept matching tested |
| Assumption justification | ✅ Validated | Component/relationship/fallback rationale tested |
| Export format output | ✅ Validated | All 5 formats produce valid output |
| Cross-platform build | ✅ Validated | Linux/macOS AMD64/ARM64, Windows AMD64 |
| Config migration | ✅ Validated | Auto-migration from legacy config path |
| License validation | ✅ Validated | HMAC signing/verification, format enforcement |
| Benchmark run | ✅ Validated | 2158 assumptions across 20 architectures processed |

### Benchmark Results (Descriptive)

The benchmark suite ran on 20 reference architectures with the following aggregate results:

- **Total assumptions processed:** 2,158
- **Risk distribution:** 0% Critical, 2% High, 9.5% Medium, 88.5% Low
- **Evidence points:** 5,428
- **Average confidence:** 38.3%
- **Risk model:** `asf-risk-model-1.0`

Note: These are descriptive statistics, not validation metrics. They describe what ASF found, not whether ASF is correct.

---

## What Is NOT Validated

These are the most critical gaps. None of these have been measured.

### Fundamental Accuracy Metrics

| Metric | Status | Impact |
|--------|--------|--------|
| **Precision** | ❌ Not measured | We don't know what fraction of ASF's findings are true positives |
| **Recall** | ❌ Not measured | We don't know what fraction of real assumptions ASF misses |
| **False positive rate** | ❌ Not measured | We don't know how many findings are noise |
| **False negative rate** | ❌ Not measured | We don't know what ASF misses |
| **F1 score** | ❌ Not measured | Combined precision/recall unknown |

### STRIDE Accuracy

| Metric | Status | Impact |
|--------|--------|--------|
| STRIDE category assignment accuracy | ❌ Not measured | Rules may misclassify |
| Per-category precision | ❌ Not measured | Some categories may be noisier than others |
| Per-category recall | ❌ Not measured | Some categories may be systematically missed |
| Expert agreement (Cohen's kappa) | ❌ Not measured | Unknown if STRIDE mapping matches human experts |

### Risk Assessment Accuracy

| Metric | Status | Impact |
|--------|--------|--------|
| Risk level calibration | ❌ Not measured | Likelihood/impact scores may be systematically wrong |
| Likelihood scoring vs actual incidents | ❌ Not measured | No correlation with real-world breach data |
| Impact scoring vs actual damage | ❌ Not measured | No correlation with real-world damage data |

### Human Validation

| Study | Status | Impact |
|-------|--------|--------|
| **Expert validation study** (10 architects × 20 architectures) | ❌ Not started | This is the single most important missing item |
| Independent derivation test (AI blind reproduction) | ❌ Not started | Would measure ASF-unique findings vs AI |
| "Would You Pay For This?" human survey | ❌ Not started | Would measure perceived value by practitioners |
| Inter-rater reliability (multiple experts) | ❌ Not started | Would measure consistency of human judgments |
| User experience study | ❌ Not started | Would measure usability and workflow fit |

### Specific Capability Validation

| Capability | Status | Notes |
|-----------|--------|-------|
| Draw.io parsing accuracy | ❌ Not measured | May miss or misparse complex diagrams |
| Mermaid parsing accuracy | ❌ Not measured | Limited to subset of Mermaid syntax |
| OCR accuracy on real diagrams | ❌ Not measured | Tesseract quality varies with image quality |
| PDF/DOCX text extraction | ❌ Not measured | May fail on complex layouts |
| AI enhancement quality | ❌ Not measured | AI-generated findings may be hallucinated |
| AI hallucination rate | ❌ Not measured | Unknown false discovery rate from AI layer |
| Python ASF CLI assumption extraction accuracy | ❌ Not measured | The foundation of all analysis — unvalidated |

### Performance Benchmarks

| Metric | Status |
|--------|--------|
| Analysis time per architecture | ❌ Not measured |
| Memory usage | ❌ Not measured |
| Scaling with architecture size | ❌ Not measured |
| Concurrent analysis performance | ❌ Not measured |

---

## Known Assumptions

These are assumptions built into ASF that have not been validated:

1. **Component labels in diagrams correspond to security-relevant entities** — A component named "Database" is assumed to be a real database. Could be a logical grouping or mislabeled.
2. **Relationships imply trust boundaries** — A line between "Internet" and "API Gateway" is assumed to be a trust boundary. May not be true in all architectures.
3. **STRIDE category rules are correct** — Category → STRIDE mappings are based on threat modeling literature but not validated against expert judgments.
4. **34 keyword patterns are comprehensive** — The keyword set may miss important security concepts.
5. **Likelihood factors are appropriate** — Exposure, authentication dependency, and attack complexity may not capture all relevant factors.
6. **Impact factors are appropriate** — Data classification, regulatory exposure, and business criticality may miss other impact dimensions.
7. **5×5 risk matrix thresholds are correct** — The 20/12/5 boundaries for Critical/High/Medium/Low are conventional but untested for this domain.
8. **Confidence formula is well-calibrated** — The 0.1 base + weighted components formula has not been calibrated against human confidence judgments.
9. **Evidence tracing is meaningful** — Substring matching of component labels against assumption text assumes that the text contains the component name.
10. **The Python ASF engine produces valid assumptions** — The foundation of all analysis depends on the Python CLI's extraction quality, which has not been independently validated.

---

## Known Limitations

1. ~~**Python dependency** — Requires Python 3.8+ with ASF package installed.~~ **Resolved in v2.0.0** — ASF is now a pure Go binary with no Python dependency.
2. **No CI/CD** — No automated testing pipeline. Changes may regress without detection.
3. **No code signing** — Binaries are not signed or notarized. Users must trust the download source.
4. **Tesseract dependency** — Image OCR requires external Tesseract installation. Not bundled.
5. **Ollama dependency** — AI features require external Ollama installation. Not bundled.
6. **Limited diagram format support** — Draw.io parsing handles only basic XML structures. Complex diagrams may fail.
7. **Mermaid parsing is basic** — Only handles `node[label]` and `-->` syntax. No subgraph, styling, or advanced syntax.
8. **OCR quality unvalidated** — No tests for image-based diagram accuracy. Results vary significantly.
9. **No cloud AI option** — By design, but limits use case for organizations that prefer cloud AI.
10. **Windows TUI untested** — Should work with Windows Terminal but not validated.
11. **No multi-architecture batch analysis** — Analyzes one architecture at a time.
12. **No comparison mode** — Cannot compare results across architectures or analysts.
13. **No persistent storage** — Results are not saved to a database. Export is the only persistence.
14. **No team features** — Single-user TUI. No collaboration, sharing, or server mode.

---

## What Would Need to Happen for Production Use

### Minimum acceptable validation (v1.0):

1. Expert validation study with ≥10 security architects reviewing ≥10 architectures each
2. Precision and recall measurement from the study
3. False positive rate < 30%
4. At least 3 of 5 export formats verified for accuracy by an independent evaluator

### Stretch goals:

1. Independent derivation test (ASF-unique findings vs AI)
2. STRIDE accuracy measurement (Cohen's kappa ≥ 0.6 with expert labels)
3. Risk calibration against real breach data
4. User experience study with ≥20 participants
5. Performance benchmarks on architectures ranging from 10-500 components

---

## Recommendations

1. **Before using ASF for any critical review**: Run a small validation study internally. Have 2-3 security architects manually review 50-100 assumptions from ASF and measure agreement rate.
2. **Before making procurement decisions**: Wait for the expert validation study results.
3. **Before citing ASF in research**: Conduct and publish independent validation.
4. **Before building on ASF**: Understand that the Python engine dependency is the weakest link and may need to be rewritten in Go for production reliability.
