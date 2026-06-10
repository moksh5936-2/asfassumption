# ASF Expert Validation Study

> Version: 1.0.0 | June 2026 | Status: Planned (Not Started)

## 1. Study Overview

### Objective
Measure ASF's precision, recall, false positive rate, and STRIDE accuracy against human security architects to validate (or invalidate) the tool's effectiveness.

### Research Questions

1. **RQ1 (Precision):** What percentage of ASF-generated assumptions are valid security assumptions?
2. **RQ2 (Recall):** What percentage of valid security assumptions in an architecture does ASF discover?
3. **RQ3 (FPR):** What percentage of ASF-generated assumptions are false positives?
4. **RQ4 (STRIDE):** How accurately does ASF map assumptions to STRIDE categories?
5. **RQ5 (Risk):** How well do ASF risk scores correlate with expert risk assessment?
6. **RQ6 (Confidence):** Does ASF's confidence score predict assumption validity?

### Hypothesis

> H0: ASF's assumption discovery is not significantly better than random keyword matching.
> H1: ASF's assumption discovery achieves ≥70% precision and ≥60% recall on real-world architectures.

---

## 2. Study Design

### Methodology

**Double-blind, controlled study** with 10 security architects. Each architect evaluates 20 architecture descriptions (10 real-world, 10 synthetic) with and without ASF assistance.

### Participant Criteria

| Criterion | Requirement |
|-----------|-------------|
| Role | Security architect or equivalent |
| Experience | 5+ years in threat modeling |
| Tool familiarity | Must review ASF outputs during study |
| Sample size | 10 participants (minimum viable) |

### Architecture Sample

| Source | Count | Complexity | Characteristics |
|--------|-------|------------|-----------------|
| Real-world (redacted) | 10 | Medium-High | Cloud, web, IoT, mobile architectures |
| Synthetic | 10 | Low-Medium | Designed with known ground truth |
| Total | 20 | — | — |

### Metrics

| Metric | Formula | Target |
|--------|---------|--------|
| Precision | TP / (TP + FP) | ≥70% |
| Recall | TP / (TP + FN) | ≥60% |
| F1 Score | 2 × (P × R) / (P + R) | ≥65% |
| False Positive Rate | FP / (FP + TN) | ≤30% |
| STRIDE Accuracy | Correct / Total STRIDE assignments | ≥75% |
| Risk Score RMSE | √(Σ(ASF - Expert)² / n) | ≤1.0 (on 1-25 scale) |
| Confidence Calibration | MSE of confidence vs. actual validity | ≤0.1 |

---

## 3. Study Protocol

### Phase 1: Preparation (Week 1-2)

1. Create 20 architecture descriptions with ground truth annotations
2. Build validation infrastructure (annotation tool, data collection)
3. Recruit 10 security architects
4. Create ASF training materials (30-min tutorial)

### Phase 2: Data Collection (Week 3-4)

**Session structure (2 hours per participant):**

| Segment | Duration | Activity |
|---------|----------|----------|
| Training | 20 min | ASF tutorial, study overview |
| Architect baseline | 30 min | Review 10 architectures manually, list assumptions |
| ASF-assisted | 30 min | Review 10 architectures with ASF output |
| Calibration | 20 min | Review 5 architectures from both conditions |
| Debrief | 20 min | Survey, interview, feedback |

### Phase 3: Analysis (Week 5)

1. Calculate precision, recall, FPR for each participant
2. Calculate STRIDE accuracy and risk score correlation
3. Analyze confidence calibration
4. Conduct qualitative thematic analysis of feedback
5. Prepare study report

---

## 4. Ground Truth Annotation

Each architecture will be annotated by 2 senior security architects (not study participants) to establish ground truth:

| Annotation | Format | Detail |
|------------|--------|--------|
| Assumption list | Text | Every security assumption in the architecture |
| Assumption type | Enum | IDENTITY, ACCESS, NETWORK, DATA, COMPUTE, DEPLOYMENT, INTEGRATION, STORAGE |
| STRIDE categories | [STRIDE] | Which threat categories apply |
| Risk level | 1-25 | Consensus risk score |
| Evidence location | Text | Where in the architecture the assumption is found |

Inter-annotator agreement (Cohen's κ) must be ≥0.70 for ground truth to be considered reliable.

---

## 5. Validation Infrastructure

ASF includes prototype validation support in `validation.go`:

```go
type ValidationRecord struct {
    ArchitectureID    string
    Timestamp         time.Time
    Assumptions       []ValidatedAssumption
    TruePositives     int
    FalsePositives    int
    FalseNegatives    int
    TrueNegatives     int
    Precision         float64
    Recall            float64
    F1Score           float64
    STRIDEAccuracy    float64
    RiskScoreRMSE     float64
}
```

The TUI has a validation screen (press `v` from Results or Review) that displays precision, recall, F1, and per-assumption validation status. However, **the validation screen is a display-only view** — there is no mechanism to submit or persist validation data.

### Infrastructure Requirements for Study

1. **Validation data persistence** — Store validation records to JSON/CSV
2. **Batch validation** — Process multiple architectures sequentially
3. **Inter-rater reliability** — Cohen's κ calculation tool
4. **Export** — Validation results to CSV for analysis in R/Python
5. **Blinding** — Participant assignment and architecture ordering

---

## 6. Budget Estimation

| Item | Cost |
|------|------|
| Participant incentives (10 × $200) | $2,000 |
| Ground truth annotation (2 architects × $500) | $1,000 |
| Infrastructure development (40 hours) | $4,000 |
| Analysis and report writing (20 hours) | $2,000 |
| **Total** | **$9,000** |

---

## 7. Risks and Mitigations

| Risk | Likelihood | Impact | Mitigation |
|------|-----------|--------|------------|
| Low recruitment | Medium | High | Increase incentives, extend timeline |
| Low inter-annotator agreement | Medium | High | Refine annotation guidelines, add training |
| ASF performs poorly | Medium | Medium | Publish results regardless (scientific integrity) |
| Architects bias toward ASF | Low | Medium | Counterbalanced design, blinding |
| Synthetic architectures too artificial | Medium | Medium | Pilot test with independent architect |

---

## 8. Timeline

```
Week 1-2: Preparation
  ├── Create architecture samples (20)
  ├── Recruit annotators
  └── Build validation infrastructure

Week 3-4: Data Collection
  ├── Conduct 10 study sessions
  └── Collect qualitative feedback

Week 5: Analysis
  ├── Calculate metrics
  ├── Thematic analysis
  └── Write report

Week 6: Publication
  ├── Internal review
  └── Public release of results
```

---

## 9. Pre-registration

Before data collection begins, the following must be pre-registered:

1. Exact metrics and formulas
2. Sample size (10 participants)
3. Statistical tests (paired t-test for precision/recall, Cohen's κ for agreement)
4. Exclusion criteria (participants who fail attention checks)
5. Minimum viable data (at least 8 participants for publication)

---

## 10. Current Status

**❌ NOT STARTED**

| Item | Status |
|------|--------|
| Study design | ✅ Complete |
| Architectures selection | ❌ Not started |
| Ground truth annotation | ❌ Not started |
| Participant recruitment | ❌ Not started |
| Validation infrastructure | ⚠️ Partial (TUI view exists, no persistence) |
| Funding | ❌ Not allocated |
| Timeline | ❌ Not scheduled |

---

*This document represents the study plan for expert validation of ASF v1.0.0. All metrics, targets, and research questions are preliminary and subject to refinement.*
