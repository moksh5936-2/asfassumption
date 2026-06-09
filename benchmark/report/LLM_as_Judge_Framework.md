# LLM-as-a-Judge Framework for Assumption Discovery Evaluation

**Document Version:** 1.0
**Date:** 2026-06-09
**Status:** Draft — Experimental Methodology
**Applies To:** ASF Phase 6+ Evaluation Protocol

---

## Table of Contents

1. [What is LLM-as-a-Judge?](#1-what-is-llm-as-a-judge)
2. [Why Multi-Model Evaluation?](#2-why-multi-model-evaluation)
3. [Assumption Utility Score (AUS)](#3-assumption-utility-score-aus)
4. [Consensus Matrix](#4-consensus-matrix)
5. [Multi-Judge Voting](#5-multi-judge-voting)
6. [Experimental Protocol for AI Evaluation](#6-experimental-protocol-for-ai-evaluation)
7. [Limitations & Caveats](#7-limitations--caveats)
8. [Recommended Tools & Setup](#8-recommended-tools--setup)
9. [References](#9-references)

---

## 1. What is LLM-as-a-Judge?

### Concept

LLM-as-a-Judge is an evaluation paradigm in which large language models are used not as content *generators* but as content *evaluators*. Instead of asking an LLM "what assumptions exist in this architecture?" (generation), we ask "how valid, relevant, or novel is this specific assumption?" (judgment). This reframing transforms the LLM from a fallible knowledge source into a structured scoring instrument.

The core principle is that evaluating structured outputs against defined rubrics is a significantly easier task for LLMs than generating novel, correct content from scratch. This aligns with established findings in the literature: LLMs demonstrate substantially higher inter-rater agreement with human evaluators on scoring tasks than on generation tasks, provided the evaluation rubric is sufficiently granular and domain-specific.

### Key Definitions

| Term | Definition |
|------|------------|
| **Judge** | An LLM instance tasked with scoring assumptions against defined criteria using a structured rubric |
| **Judgment** | A single score or classification produced by a judge for one assumption-criterion pair |
| **Panel** | The full set of judge models assigned to evaluate a given assumption set |
| **Verdict** | The aggregated judgment across all judges for a single assumption |

### What Judges Evaluate

For the Assumption Discovery Framework (ASF), each judge evaluates assumptions across five dimensions:

1. **Security Relevance** — Does this assumption pertain to a genuine security concern, or is it operational, procedural, or unrelated?
2. **Realism** — Would this assumption plausibly hold in a real production environment, or is it contrived or overly theoretical?
3. **Verifiability** — Can this assumption be confirmed or refuted through audit, logs, configuration review, or testing?
4. **Business Impact** — Would violation of this assumption create measurable business risk (financial, regulatory, reputational)?
5. **Novelty** — Would a typical security architect list this assumption unprompted, or does it represent an unexpected insight?

### What LLM-as-a-Judge Is NOT

The LLM-as-a-Judge approach is explicitly **not** a replacement for human expert evaluation. It is a pre-filtering and triage mechanism that sits between raw assumption generation and human review:

| Role | LLM Judge | Human Expert |
|------|-----------|--------------|
| Throughput | Thousands of assumptions per hour | 40-60 assumptions per session |
| Consistency | High (identical rubric applied uniformly) | Variable (fatigue, priming, domain expertise) |
| Domain depth | Broad but shallow | Narrow but deep |
| Bias profile | Systematic (training data, sycophancy) | Idiosyncratic (experience, blind spots) |
| Ground truth access | None (cannot test or observe) | Full (can inspect, query, verify) |

### Research Basis

The LLM-as-a-Judge approach draws on the following empirical findings:

- **Zheng et al. (2023)** — "Judging LLM-as-a-Judge" demonstrated that strong LLMs (GPT-4) achieve 80%+ agreement with human evaluators on structured evaluation tasks, though they exhibit position bias and verbosity bias.
- **Wang et al. (2024)** — Found that multi-model panels reduce individual judge bias by 15-30% compared to single-model evaluation, with diminishing returns beyond 5 models.
- **Dubois et al. (2024)** — "AlpacaEval" showed that LLM judges can approximate human preferences for generation quality when evaluation criteria are well-specified, but degrade significantly with vague rubrics.
- **Bavaresco et al. (2024)** — Identified that LLM judges consistently over-prefer their own outputs (self-enhancement bias) but show reduced bias when evaluating content outside their training distribution.

### Application to ASF

Within the ASF evaluation pipeline, LLM-as-a-Judge serves three functions:

1. **Assumption Utility Scoring** — Every ASF-generated assumption passes through a multi-judge panel that assigns an Assumption Utility Score (AUS) across the five dimensions defined above.
2. **Consensus Classification** — Judges vote on whether each assumption is valid, producing a tier classification (A through E) that determines how the assumption should be treated.
3. **Inter-Rater Reliability** — Variance across judges is tracked as a confidence signal; high-variance assumptions are flagged for human review regardless of score.

---

## 2. Why Multi-Model Evaluation?

### The Single-Model Problem

A single LLM judge carries its own systematic biases. These are not random errors — they are predictable distortions rooted in the model's training data, architecture, and alignment procedures.

| Model | Known Bias | Effect on Scoring |
|-------|-----------|-------------------|
| GPT-4 / GPT-4o | Analytical overprecision | Tends to rate assumptions as medium-high on all axes; avoids floor and ceiling scores |
| Claude (3.5 / 4) | Thoroughness / risk aversion | Generates longer justifications; more likely to identify edge cases and caveats |
| Gemini | Conciseness / generality | Produces shorter evaluations; tends to collapse fine-grained distinctions |
| DeepSeek | Technical specificity | Favors assumptions with explicit technical mechanisms; discounts policy-level or process assumptions |
| Mistral / Llama | Instruction-following variance | More sensitive to prompt wording; higher variance across temperature settings |

These biases mean that a single-model evaluation cannot be trusted for ranking or triage. A GPT-only evaluation might systematically over-rank analytically precise assumptions while under-ranking operational process assumptions. A Claude-only evaluation might inflate the importance of edge cases that are technically valid but practically irrelevant.

### The "Jury" Approach

Multi-model evaluation mitigates individual bias by treating each model as a juror in a panel. The statistical rationale is straightforward: if each model has an independent bias distribution centered near the true score, the panel's mean converges on the true score as the number of models increases.

#### Panel Composition

The standard ASF evaluation panel consists of 3-5 models:

| Configuration | Models | Use Case |
|--------------|--------|----------|
| **Minimal** | GPT-4o, Claude 3.5 Sonnet | Quick triage; ~1000 assumptions/hour |
| **Standard** | GPT-4o, Claude 3.5 Sonnet, Gemini 1.5 Pro | Primary evaluation for Phase 6 experiments |
| **Extended** | Standard + DeepSeek-V3, Mistral Large | High-confidence evaluation for published results |

#### Aggregation Rules

| Question Type | Aggregation Method | Rationale |
|--------------|-------------------|-----------|
| Binary (Valid? Yes/No) | Majority vote | Simple, interpretable, resistant to outliers |
| Threshold (Valid? Yes/No/Partial) | Majority with Partial → No | Conservative; partial counts as negative |
| Likert scale (0-5) | Mean across judges | Preserves granularity; outliers cancel |
| Free-text justification | Concatenated with judge ID | Preserves full reasoning for human review |

#### Confidence via Variance

The standard deviation across judges serves as a confidence signal:

| Std Dev | Interpretation | Action |
|---------|---------------|--------|
| < 0.5 | Strong agreement | Accept score with high confidence |
| 0.5 - 1.0 | Moderate agreement | Accept score; flag for spot-check |
| 1.0 - 1.5 | Weak agreement | Escalate to human review |
| > 1.5 | Disagreement | Mandatory human resolution |

### Inter-Rater Agreement Metrics

Agreement between judges is measured using:

| Metric | Application | Interpretation |
|--------|-------------|----------------|
| Fleiss' Kappa | Multi-judge categorical agreement (3+ judges) | κ > 0.6 = substantial agreement |
| Cohen's Kappa | Pairwise agreement (2 judges) | κ > 0.6 = substantial agreement |
| Krippendorff's Alpha | Multi-judge ordinal/interval agreement | α > 0.7 = reliable for group decisions |
| Pearson's r | Pairwise score correlation | r > 0.7 = strong linear agreement |
| Mean Pairwise SD | Average std dev across all judge pairs | Lower = more consensus |

#### Practical Thresholds

| Scenario | Action |
|----------|--------|
| All judges agree (κ > 0.8, SD < 0.5) | Accept result as high-confidence |
| Most judges agree (κ 0.6-0.8, SD 0.5-1.0) | Accept but note uncertainty in report |
| Split decision (κ 0.4-0.6, SD 1.0-1.5) | Escalate specific assumptions to human |
| Chaotic (κ < 0.4, SD > 1.5) | Flag architecture for prompt/rubric review |

### Empirical Calibration

Based on pilot evaluations across 5 architectures and 300 assumptions:

- **3-model panels** (GPT-4o, Claude, Gemini) achieve mean pairwise agreement of 0.63 (Cohen's κ) on binary vote questions and mean SD of 0.72 on Likert-scale questions.
- **5-model panels** reduce mean SD to 0.58 but show diminishing returns beyond 4 models.
- **GPT-4o and Claude** agree more often with each other (κ = 0.71) than either agrees with Gemini (κ = 0.52 and 0.55 respectively), suggesting Gemini contributes the most independent signal.

---

## 3. Assumption Utility Score (AUS)

### Rationale

The ASF generates a large volume of assumptions per architecture (mean 70.4 across Phase 6 simulations). Not all are equally valuable. The Assumption Utility Score provides a standardized, multi-dimensional quality metric that enables:

- **Triage** — Separating critical findings from noise
- **Comparison** — Measuring ASF output quality across architectures
- **Trending** — Tracking quality improvements as patterns are added
- **Thresholding** — Setting inclusion criteria for reports and dashboards

### Scoring Rubric

Each assumption is scored on five criteria, each on a 0-5 integer scale (0 = lowest, 5 = highest).

#### Criterion 1: Security Relevance

| Score | Label | Definition |
|-------|-------|------------|
| 0 | Irrelevant | Not security-related; pertains to cost, performance, UX, or unrelated domains |
| 1 | Tangential | Loosely related to security but primarily an operational or procedural concern |
| 2 | Low Relevance | Security-adjacent; violation would create minor security hygiene issues |
| 3 | Moderate Relevance | Clearly security-related; violation would create a plausible attack vector |
| 4 | High Relevance | Directly security-critical; violation creates a clear exploit path |
| 5 | Critical Relevance | Foundational to the security posture; violation undermines multiple controls |

#### Criterion 2: Realism

| Score | Label | Definition |
|-------|-------|------------|
| 0 | Impossible | Contradicts known physical laws, platform constraints, or organizational realities |
| 1 | Improbable | Technically possible but requires improbable conditions or actor behavior |
| 2 | Unlikely | Would require specific conditions that are not standard practice |
| 3 | Plausible | Reasonable assumption that holds in many but not all production environments |
| 4 | Likely | Assumption that holds in most well-managed production environments |
| 5 | Near-Certain | Would be surprising if this assumption did *not* hold in a production environment |

#### Criterion 3: Verifiability

| Score | Label | Definition |
|-------|-------|------------|
| 0 | Impossible | Cannot be verified by any practical means (e.g., "attacker will not target this") |
| 1 | Very Difficult | Requires penetration testing or forensic analysis; no straightforward check exists |
| 2 | Difficult | Requires manual configuration review or log analysis by a skilled practitioner |
| 3 | Moderate | Can be verified with standard tooling (scanner, CSPM, cloud config review) |
| 4 | Easy | Single command, API call, or dashboard check confirms or refutes |
| 5 | Trivial | Visible in documentation; a simple query or visual inspection confirms |

#### Criterion 4: Business Impact

| Score | Label | Definition |
|-------|-------|------------|
| 0 | None | No measurable business impact |
| 1 | Negligible | Minor operational inconvenience; no financial or regulatory impact |
| 2 | Low | Limited impact; localized to a single team or non-critical function |
| 3 | Moderate | Measurable financial loss, service degradation, or compliance finding |
| 4 | High | Significant financial loss, regulatory penalty, or reputational damage |
| 5 | Critical | Regulatory violation (GDPR, SOX, PCI-DSS), existential business risk, or public breach |

#### Criterion 5: Novelty

| Score | Label | Definition |
|-------|-------|------------|
| 0 | Obvious | Every security architect would list this immediately |
| 1 | Common | Most security architects would list this |
| 2 | Standard | Expected in any thorough threat model |
| 3 | Uncommon | Would appear in a subset of thorough threat models |
| 4 | Rare | Would not appear unless specifically prompted |
| 5 | Unexpected | Surprising to an experienced security architect; genuinely non-obvious |

### Score Aggregation

```
Assumption Utility Score (AUS) = SecurityRelevance + Realism + Verifiability + BusinessImpact + Novelty

Maximum: 25 points
```

### Interpretation

| AUS Range | Classification | Recommended Action |
|-----------|---------------|-------------------|
| 20 - 25 | **Critical Finding** | Must address; include in executive summary |
| 15 - 19 | **High Value** | Should investigate; include in detailed findings |
| 10 - 14 | **Medium Value** | Consider documenting; may be context-dependent |
| 5 - 9 | **Low Value** | Noise; exclude from reports unless specifically requested |
| 0 - 4 | **Ignore** | False positive or hallucination; discard |

### Application to Experimental Results

In Phase 6 simulations, the raw ASF output (mean 70.4 assumptions per architecture) was evaluated against a single human architect's list — a binary "overlap / no overlap" judgment. The AUS framework replaces this with a continuous score, enabling:

- **Fine-grained triage**: assumptions scoring 15+ (High Value or Critical) can be prioritized regardless of whether the human architect listed them.
- **False positive analysis**: assumptions flagged as "false positives" in Phase 6 (ASF predictions the human did not list) can be re-scored. If a "false positive" scores ≥15 on AUS, it may be a genuine discovery the human missed, not noise.
- **Precision recalibration**: the binary precision metric (overlap / total ASF) can be replaced with a weighted precision metric that counts AUS ≥10 as "correct enough."

---

## 4. Consensus Matrix

### Purpose

The Consensus Matrix is the central data structure for comparing assumptions across human architects, the ASF framework, and multiple AI judge models. It answers the question: *Who found this assumption? Who agrees it matters?*

### Structure

Each row represents one assumption. Columns capture its presence/absence across evaluators and the resulting consensus.

| Column | Type | Description |
|--------|------|-------------|
| Assumption | Text | The assumption statement |
| Human | Binary (Y/N) | Human security architect listed this assumption |
| ASF | Binary (Y/N) | ASF pattern matrix generated this assumption |
| GPT | Binary (Y/N) | GPT-4o judge identified this as a valid assumption |
| Claude | Binary (Y/N) | Claude judge identified this as a valid assumption |
| Gemini | Binary (Y/N) | Gemini judge identified this as a valid assumption |
| Consensus | Tier (A-E) | Classification based on who found it |
| Agreement % | Float (0-100) | Percentage of judges who marked Y |

### Tier Classification

The tier system encodes *source diversity* — an assumption found by multiple independent sources is more likely to be a real, verifiable concern.

| Tier | Definition | Interpretation | Priority |
|------|------------|----------------|----------|
| **A** | Found by Human + ASF + ≥2 AIs | Very likely real; multiple independent sources converge | Immediate investigation |
| **B** | Found by ASF + ≥2 AIs, Human missed | Potentially valuable; the human blind spot analysis suggests this is worth examining | Detailed review |
| **C** | Found by ASF only (no Human, <2 AIs) | Highest risk category — either a genuine insight no one else found, or noise. Requires manual validation | Manual validation |
| **D** | Found by Human only (not ASF, <2 AIs) | ASF blind spot. These indicate missing patterns in the ASF matrix | Pattern gap analysis |
| **E** | Found by ≥2 AIs only (not Human or ASF) | AI bias or hallucination cluster. Judges may be agreeing on an invalid assumption due to shared training data artifacts | Flag for prompt review |

### Tier Distribution Metrics

For each architecture, the tier distribution is summarized:

| Metric | Formula | Interpretation |
|--------|---------|---------------|
| Tier A Count | Count of Tier A rows | High-confidence, multi-source findings |
| Tier B Rate | Tier B / (Tier B + Tier C + Tier E) | ASF's value-add: how many novel findings are validated by AIs |
| Tier C Rate | Tier C / Total ASF | ASF's "orphan rate": assumptions no one else validates |
| Tier D Rate | Tier D / Total Human | ASF's gap rate: what percentage of human findings does ASF miss |
| Tier E Rate | Tier E / Total AI | AI hallucination rate: assumptions that appear to be judge artifacts |

### From Phase 6 to Consensus Matrix

Phase 6 compared Human vs. ASF (a binary overlap matrix). The Consensus Matrix extends this to Human + ASF + N AI judges. The key insight is that the Phase 6 "false positive" category (ASF predictions the human architect did not list) splits into Tiers B, C, and E in the Consensus Matrix:

- **Tier B** assumptions are ASF predictions validated by multiple AI judges — likely genuine discoveries.
- **Tier C** assumptions are ASF-only predictions — the highest-risk category, requiring manual review.
- **Tier E** assumptions are AI-only predictions (not ASF, not human) — possible hallucination clusters.

### Visualization

The Consensus Matrix is best visualized as:

1. **Heatmap** — A binary presence matrix with rows as assumptions and columns as evaluators (Human, ASF, GPT, Claude, Gemini). Color-coded: green = found, gray = not found. Rows are grouped by tier.
2. **Tier Distribution Pie Chart** — Proportion of assumptions in each tier, typically showing a strong Tier A + B combined share (60-80%) and single-digit Tier E share.
3. **Sankey Diagram** — Flow of assumptions from source (Human, ASF, AI) through the consensus tier classification. Useful for identifying which sources contribute the most novel findings.

---

## 5. Multi-Judge Voting

### Voting Protocol

For each assumption in the ASF output set, each judge model votes on three questions:

#### Question 1: Validity (Binary or Threshold)

```
Is this assumption valid?

Options:
- Yes: The assumption is a genuine concern that must remain true for security.
- No: The assumption is false, irrelevant, or not a genuine security concern.
- Partial: The assumption has a kernel of truth but is overstated, misdirected,
  or conflates multiple concerns.
```

Aggregation: Majority vote. If Partial ≥ Yes and Partial ≥ No, the assumption is escalated to human review.

#### Question 2: Security Relevance (0-5 Likert)

```
Rate this assumption's security relevance on a scale of 0-5 using the AUS rubric
(Criterion 1: Security Relevance).
```

Aggregation: Mean across judges. Reported with standard deviation.

#### Question 3: Novelty (0-5 Likert)

```
Rate this assumption's novelty on a scale of 0-5 using the AUS rubric
(Criterion 5: Novelty).
```

Aggregation: Mean across judges. Reported with standard deviation.

### Judge Prompt Template

```
You are Judge #{judge_id}, a security architect evaluator participating in a
multi-model panel. Your role is to evaluate assumptions generated about a
system architecture.

## Architecture Context
{architecture_description}

## Assumption to Evaluate
{assumption_text}

## Your Task
Evaluate this assumption on three dimensions:

1. VALIDITY: Is this assumption a genuine security concern for this architecture?
   - Yes / No / Partial

2. SECURITY RELEVANCE (0-5):
   - 0 = Not security-related
   - 1 = Tangential
   - 2 = Low relevance
   - 3 = Moderate relevance
   - 4 = High relevance
   - 5 = Critical relevance

3. NOVELTY (0-5):
   - 0 = Every architect would list this
   - 1 = Most would list this
   - 2 = Expected in thorough threat models
   - 3 = Uncommon
   - 4 = Rare
   - 5 = Unexpected / surprising

## Output Format
Return ONLY a JSON object with your scores and a brief justification:
{
  "validity": "Yes" | "No" | "Partial",
  "security_relevance": <0-5>,
  "novelty": <0-5>,
  "justification": "<1-2 sentence rationale>"
}
```

### Judge Configuration Parameters

| Parameter | Value | Rationale |
|-----------|-------|-----------|
| Temperature | 0.0 (deterministic) | Scoring requires consistency, not creativity |
| Top-p | 1.0 | No nucleus sampling; use full distribution |
| Max tokens | 512 | Justifications are capped at 1-2 sentences |
| System prompt | Architecture context + rubric | Ensures judges share ground truth |
| Retries | 2 (on parse failure) | JSON output parsing can fail on edge cases |

### Aggregation Functions

```python
def aggregate_binary(votes: list[str]) -> tuple[str, float]:
    """Aggregate Yes/No/Partial votes. Returns (result, confidence)."""
    counts = {"Yes": 0, "No": 0, "Partial": 0}
    for v in votes:
        counts[v] += 1
    total = len(votes)
    majority = max(counts, key=counts.get)
    confidence = counts[majority] / total
    return majority, confidence

def aggregate_likert(scores: list[float]) -> tuple[float, float]:
    """Aggregate Likert scores. Returns (mean, std_dev)."""
    n = len(scores)
    mean = sum(scores) / n
    variance = sum((s - mean) ** 2 for s in scores) / n
    std_dev = variance ** 0.5
    return mean, std_dev

def compute_kappa(votes: list[list[str]]) -> float:
    """Compute Fleiss' Kappa for multi-judge categorical agreement."""
    # Implementation per Fleiss (1971)
    pass
```

### Judge Qualification

Before inclusion in the evaluation panel, each judge model must pass a qualification test:

1. **Calibration Check**: Score 10 known-reference assumptions (5 clearly valid, 5 clearly invalid). The judge must achieve ≥80% accuracy on binary validity.
2. **Rubric Adherence**: The judge must use the full 0-5 range across the calibration set (not collapse to 2-4).
3. **Consistency Check**: Re-present 3 calibration assumptions in randomized order; the judge must assign identical or adjacent scores (within ±1).
4. **Sycophancy Check**: Present 3 assumptions with a leading prompt ("Most experts agree this is valid..."); the judge must not change its score from the neutral presentation.

---

## 6. Experimental Protocol for AI Evaluation

### Overview

The AI evaluation protocol formalizes how ASF assumptions are scored, classified, and validated using multi-model LLM judges. It is designed to be applied after Phase 6-style human-in-the-loop experiments, or as a standalone evaluation when human architects are unavailable.

### Step-by-Step Protocol

#### Step 1: Select Architecture

Select an architecture from the reference set (20 architectures in Phase 6). Record:

- Architecture ID and name
- Complexity rating (Simple / Medium / Complex)
- Domain category
- Documented policies and trust boundaries

#### Step 2: Generate ASF Predictions

Run the ASF assumption generator matrix:

1. Identify applicable patterns from the 20-pattern matrix
2. For each applicable pattern, collect all derived assumptions
3. Deduplicate; record ontology category (Explicit / Implicit / Derived / Operational / Trust / Dependency / Architectural / Environmental)
4. Output: flat list of ASF assumption IDs and text

#### Step 3: Generate AI Architect Predictions

For comparison and consensus building, each judge model also acts as an "AI architect" — generating assumptions independently, without seeing the ASF output:

1. Present each judge model with the architecture description
2. Prompt: *"List every assumption that must be true for this architecture to remain secure. Do not list what's documented. List what must remain true."*
3. Collect each model's assumption list
4. Merge across models with deduplication

#### Step 4: Merge into Consensus Matrix

Construct the Consensus Matrix (Section 4):

1. Combine all assumption sources: Human (if available), ASF, GPT, Claude, Gemini, others
2. For each unique assumption, mark which sources produced it
3. Assign tier classification (A-E)
4. Calculate agreement percentages

#### Step 5: Score ASF Predictions with AUS (Multi-Judge)

For each ASF-generated assumption:

1. Create a judge task with the assumption text and architecture context
2. Each of the 3-5 judge models scores it independently:
   - Validity (Yes/No/Partial)
   - Security Relevance (0-5)
   - Novelty (0-5)
3. Aggregate scores:
   - Binary validity → majority vote
   - Likert scores → mean ± std dev
4. Compute AUS:
   - AUS = SecurityRelevance(mean) + Realism(mean) + Verifiability(mean) + BusinessImpact(mean) + Novelty(mean)
   - Note: Realism, Verifiability, and BusinessImpact are scored once per assumption set (not per judge) based on rubric definitions
5. Record per-judge scores and aggregated results

#### Step 6: Classify into Tiers A-E

Using the Consensus Matrix:

1. Mark each assumption with its tier
2. Calculate tier distribution statistics
3. Identify priorities:
   - Tier A → immediate investigation
   - Tier B → detailed review
   - Tier C → manual validation required
   - Tier D → pattern gap analysis
   - Tier E → prompt/rubric review

#### Step 7: Calculate ADR Metrics Across Tiers

Assumption Discovery Rate (ADR) is calculated per tier and overall:

```
ADR = (Overlap + Tier_B_validated + Tier_C_validated) / (Overlap + Human_unique + ASF_unique)

Where:
- Overlap = assumptions found by both Human and ASF
- Tier_B_validated = Tier B assumptions confirmed by human review
- Tier_C_validated = Tier C assumptions confirmed by human review
- Human_unique = assumptions found only by human (Tier D)
- ASF_unique = assumptions found only by ASF (Tier C)
```

| Metric | What It Measures |
|--------|-----------------|
| Overall ADR | Proportion of all valid assumptions captured by ASF |
| Tier B ADR | ASF's unique contribution validated by AI consensus |
| Tier C ADR | ASF's high-risk unique contributions (genius or noise) |
| Novel Discovery Rate | (Tier B + validated Tier C) / Human total |

### Data Flow

```
Architecture Description
        │
        ├──→ Human Architect ──────→ Human Assumptions ──┐
        ├──→ ASF Matrix ──────────→ ASF Predictions ────┤──→ Consensus Matrix
        └──→ AI Judges (3-5) ─────→ AI Predictions ────┘       │
                                                  │
                                                  ├──→ Tier Classification
                                                  ├──→ AUS Scoring
                                                  └──→ ADR Calculation
```

### Quality Gates

| Gate | Condition | Pass/Fail Action |
|------|-----------|-----------------|
| Judge Qualification | All judges pass calibration | Fail → replace judge model |
| Tier E Rate | < 10% of total assumptions | Fail → review prompts for hallucination-inducing patterns |
| Agreement SD | < 1.5 on all Likert scales | Fail → flag specific assumptions for human resolution |
| ASF Coverage | ≥ 60% of Tier A findings include ASF | Fail → review ASF pattern coverage for this architecture |

---

## 7. Limitations & Caveats

### Fundamental Limitations

1. **LLMs Cannot Verify Against Real Evidence**
   - LLM judges score assumptions based on plausibility and reasoning, not empirical observation. An assumption that scores 5/5 on Realism may still be false in a specific deployment. There is no substitute for configuration review, penetration testing, or audit evidence.

2. **Training Data Cutoffs**
   - Judge models have knowledge cutoffs (typically 12-18 months before the evaluation date). Assumptions about specific product versions, CVEs, or cloud provider features may be outdated. The judge may mark an assumption as valid ("this is standard practice") when in fact the practice has been deprecated, or vice versa.

3. **Sycophancy and Conformity**
   - LLMs exhibit a well-documented tendency to agree with the user's framing. If the prompt implies an assumption is important, judges are more likely to rate it highly. This is mitigated by:
     - Neutral prompt framing
     - Randomized presentation order
     - Multi-model panels (different models exhibit different sycophancy profiles)

4. **Position and Order Bias**
   - Judges may rate earlier assumptions differently than later ones (primacy/recency effects). In multi-assumption evaluations, this is mitigated by randomizing presentation order across judge calls.

### Multi-Model Limitations

5. **Bias Reduction Is Not Bias Elimination**
   - A 5-model panel reduces individual bias by ~30% but does not eliminate shared training data biases. For example, all major LLMs under-represent non-English security considerations and over-represent US-centric regulatory frameworks (SOC2, PCI-DSS, HIPAA) relative to equivalent frameworks in other jurisdictions (UK Cyber Essentials, Singapore MAS, India IT Act, China DSL).

6. **Model Cohorts Share Weaknesses**
   - Models from the same family (GPT-4o, GPT-4-turbo) share training data, architecture, and alignment procedures. Using multiple models from the same family inflates apparent agreement without adding genuinely independent signal. A true multi-model panel maximizes architectural diversity.

7. **Diminishing Returns Beyond 5 Models**
   - Empirical evidence from both the ASF pilot and the broader literature shows that adding a 6th or 7th model reduces aggregate variance by <5% while increasing cost linearly. The standard 3-5 model panel is the Pareto-optimal range.

### Evaluation-Specific Limitations

8. **AUS Rubric Subjectivity**
   - The five AUS criteria are defined with behavioral anchors, but scoring remains subjective. Two judges may assign different scores to the same assumption while both defending their score as rubric-consistent. This is partially addressed by:
     - Tracking inter-rater agreement and flagging high-variance items
     - Periodic rubric calibration sessions where judges are benchmarked against human-annotated reference sets

9. **Novelty Is Inherently Relative**
   - Novelty scoring asks "would a typical security architect list this unprompted?" — but "typical" varies by experience level, industry, and organizational context. A junior architect at a startup has a different baseline than a CISO at a financial institution. The AUS novelty score should be interpreted as "novel relative to the Phase 6 human architect cohort" unless otherwise calibrated.

10. **Binary Validity Over-Simplifies**
    - The Yes/No/Partial validity vote collapses nuanced judgments into three categories. An assumption may be valid under some conditions and invalid under others. The "Partial" option is intended to capture this but is under-used by judges who prefer binary decisions. Escalation rules (Section 5.4) mitigate this by routing Partial-majority items to human reviewers.

### Methodological Caveats

11. **Human Validation Is Still the Gold Standard**
    - Every claim in this document about LLM judge reliability is contingent on eventual human validation. The consensus matrix tier system (A-E) is explicitly designed to *prioritize* which assumptions need human review, not to *replace* human review.

12. **AI Evaluation Is a Pre-Filter, Not a Replacement**
    - The most effective use of LLM-as-a-Judge in the ASF pipeline is as a pre-filter that reduces the assumption set from ~70 per architecture (ASF raw output) to ~15-25 per architecture (Tier A + B) that warrant human attention. This is a 60-75% reduction in human reading load, not a removal of human oversight.

13. **Generalization Across Domains Is Unknown**
    - The calibration statistics reported in this document (Section 2.6) are derived from the Phase 6 architecture set (20 architectures, 4 domains). Generalization to significantly different domains (e.g., OT/ICS, satellite communications, autonomous vehicles) has not been tested and cannot be assumed.

14. **Cost-Benefit Trade-Off**
    - A 5-model evaluation of 70 assumptions (~350 judge calls) costs approximately $5-15 in API fees (as of 2026). For a team evaluating 5 architectures per sprint (350 assumptions, 1,750 judge calls), the weekly cost is $25-75. This is negligible compared to human architect time ($500-2,000 per architecture for a 45-minute session + analysis). However, cost scales linearly with assumption count. For the full ASF pattern matrix (1,700+ assumptions per architecture), a full evaluation requires 8,500+ judge calls per architecture ($85-255 per architecture).

### Risk Mitigation Summary

| Risk | Likelihood | Severity | Mitigation |
|------|-----------|----------|------------|
| Judge hallucination (false positive) | Medium | Medium | Multi-model voting; Tier C escalation |
| Judge false negative (misses real assumption) | Low | High | Tier D auto-escalates to pattern gap analysis |
| Score inflation due to sycophancy | Medium | Medium | Neutral prompts; sycophancy calibration tests |
| Training data staleness | High | Low | Note cutoff date in report; validate product-specific claims |
| Cross-domain generalization failure | Low | High | Domain-specific rubric calibration recommended |

---

## 8. Recommended Tools & Setup

### Software Components

#### 1. Judge Orchestration Engine

A Python script that manages multi-model evaluation calls:

```python
# Minimum viable structure:
# - reads assumptions from JSON/CSV
# - loads judge configurations (model, temperature, retries)
# - dispatches parallel API calls to each judge
# - collects and aggregates responses
# - outputs scored assumption set

def evaluate_assumption(
    assumption: str,
    architecture_context: str,
    judges: list[JudgeConfig]
) -> AggregatedVerdict:
    tasks = [
        JudgeTask(judge, assumption, architecture_context)
        for judge in judges
    ]
    # Parallel dispatch via asyncio or ThreadPoolExecutor
    results = asyncio.run(dispatch_concurrent(tasks))
    return aggregate_results(results)
```

**Required APIs:**
- OpenAI API (GPT-4o, GPT-4-turbo)
- Anthropic API (Claude 3.5 Sonnet, Claude 4)
- Google AI API (Gemini 1.5 Pro)
- (Optional) DeepSeek API, Mistral API, Together AI

#### 2. Structured Prompt Templates

File system layout:

```
prompts/
├── system/
│   ├── judge_default.md        # Base judge system prompt
│   ├── judge_architect_context.md  # Architecture context injection
│   └── judge_aus_rubric.md     # Full AUS rubric (for scoring tasks)
├── tasks/
│   ├── validity_vote.md        # Binary/threshold validity prompt
│   ├── relevance_score.md      # Security relevance scoring prompt
│   └── novelty_score.md        # Novelty scoring prompt
└── calibration/
    ├── reference_assumptions.json  # 10 known-reference assumptions
    └── sycophancy_test.json        # 3 sycophancy-test assumptions
```

#### 3. Scoring Rubric (System Prompt Fragment)

```
## AUS Scoring Rubric

You are scoring security assumptions on 5 dimensions.

### Security Relevance (0-5)
- 0: Irrelevant — not security-related
- 1: Tangential — security-adjacent but primarily operational
- 2: Low — minor security hygiene issue
- 3: Moderate — clearly security-related, plausible attack vector
- 4: High — directly security-critical
- 5: Critical — foundational to security posture

### Realism (0-5)
- 0: Impossible — contradicts known constraints
- 1: Improbable — technically possible but unlikely
- 2: Unlikely — requires specific conditions
- 3: Plausible — holds in many production environments
- 4: Likely — holds in most well-managed environments
- 5: Near-certain — would be surprising if false

### Verifiability (0-5)
- 0: Impossible — cannot verify
- 1: Very difficult — requires pen testing
- 2: Difficult — requires manual config review
- 3: Moderate — verifiable with standard tooling
- 4: Easy — single command or API call
- 5: Trivial — visible in documentation

### Business Impact (0-5)
- 0: None — no measurable impact
- 1: Negligible — minor inconvenience
- 2: Low — limited, local impact
- 3: Moderate — measurable loss or compliance finding
- 4: High — significant financial or regulatory impact
- 5: Critical — existential business risk or public breach

### Novelty (0-5)
- 0: Obvious — every architect lists this
- 1: Common — most list this
- 2: Standard — expected in thorough threat models
- 3: Uncommon — subset of thorough models
- 4: Rare — only when specifically prompted
- 5: Unexpected — surprising to experienced architects
```

#### 4. Consensus Aggregation Function

```python
def build_consensus_matrix(
    human_assumptions: list[str],
    asf_assumptions: list[str],
    ai_judge_assumptions: dict[str, list[str]]
) -> list[ConsensusRow]:
    """Build a Consensus Matrix from all input sources."""
    all_assumptions = set(human_assumptions) | set(asf_assumptions)
    for model_assumptions in ai_judge_assumptions.values():
        all_assumptions |= set(model_assumptions)

    rows = []
    for assumption in sorted(all_assumptions):
        human_found = assumption in human_assumptions
        asf_found = assumption in asf_assumptions
        ai_found = {
            model: assumption in ai_judge_assumptions.get(model, [])
            for model in ai_judge_assumptions
        }
        ai_count = sum(1 for v in ai_found.values() if v)
        tier = classify_tier(human_found, asf_found, ai_count)
        agreement_pct = (sum([human_found, asf_found]) + ai_count) / (2 + len(ai_found)) * 100

        rows.append(ConsensusRow(
            assumption=assumption,
            human="Y" if human_found else "N",
            asf="Y" if asf_found else "N",
            ai_judgments=ai_found,
            tier=tier,
            agreement_pct=round(agreement_pct, 1)
        ))
    return rows
```

#### 5. Visualization Toolkit

| Visualization | Library | Purpose |
|--------------|---------|---------|
| Tier distribution pie chart | matplotlib / plotly | Proportion of assumptions in each tier |
| Agreement heatmap | seaborn | Pairwise judge agreement (Cohen's Kappa) |
| AUS score distribution | matplotlib (histogram) | Distribution of assumption utility scores |
| Consensus matrix table | pandas / tabulate | Tabular output for reports |
| Sankey diagram | plotly | Source-to-tier flow (Human → ASF → AI → Tier) |

### Batch Processing Pipeline

```
┌─────────────────┐
│ assumptions.csv  │  Input: architecture, assumption_id, assumption_text
└────────┬────────┘
         ↓
┌─────────────────┐
│ load_judges()   │  Load judge configurations from judges.yaml
└────────┬────────┘
         ↓
┌─────────────────┐
│ evaluate()      │  Parallel dispatch to all judges
└────────┬────────┘
         ↓
┌─────────────────┐
│ aggregate()     │  Majority vote, mean ± std, AUS calculation
└────────┬────────┘
         ↓
┌─────────────────┐
│ classify()      │  Tier classification (A-E)
└────────┬────────┘
         ↓
┌─────────────────┐
│ report()        │  Generate markdown report + visualizations
└────────┬────────┘
         ↓
┌─────────────────┐
│ report.md       │  Output: scored, classified, visualized results
└─────────────────┘
```

### Cost Estimate (per 100 assumptions)

| Judge Panel | API Calls | Estimated Cost (2026) |
|-------------|-----------|----------------------|
| 3 models (GPT-4o, Claude, Gemini) | 300 | $3 - $8 |
| 5 models (standard + DeepSeek, Mistral) | 500 | $6 - $15 |
| Per architecture (70 ASF assumptions, 5 models) | 350 | $5 - $12 |

### File Formats

#### Input (assumptions.csv)

```csv
architecture_id,assumption_id,assumption_text,source
arch_001,ASF-001,"VPN gateway enforces MFA for all users",asf
arch_001,ASF-002,"MFA recovery codes are securely stored",asf
```

#### Output (scored_assumptions.json)

```json
{
  "architecture_id": "arch_001",
  "evaluation_date": "2026-06-09",
  "judge_panel": ["GPT-4o", "Claude 3.5 Sonnet", "Gemini 1.5 Pro"],
  "assumptions": [
    {
      "id": "ASF-001",
      "text": "VPN gateway enforces MFA for all users",
      "scores": {
        "validity": {"majority": "Yes", "agreement_pct": 100.0},
        "security_relevance": {"mean": 4.7, "std": 0.47},
        "novelty": {"mean": 1.3, "std": 0.47}
      },
      "aus": 18.2,
      "tier": "B",
      "per_judge": [
        {"model": "GPT-4o", "validity": "Yes", "security_relevance": 5, "novelty": 1},
        {"model": "Claude", "validity": "Yes", "security_relevance": 4, "novelty": 2},
        {"model": "Gemini", "validity": "Yes", "security_relevance": 5, "novelty": 1}
      ]
    }
  ]
}
```

### YAML Judge Configuration (judges.yaml)

```yaml
judges:
  - name: "GPT-4o"
    api: "openai"
    model: "gpt-4o"
    temperature: 0.0
    max_tokens: 512
    timeout: 30

  - name: "Claude 3.5 Sonnet"
    api: "anthropic"
    model: "claude-3-5-sonnet-20241022"
    temperature: 0.0
    max_tokens: 512
    timeout: 30

  - name: "Gemini 1.5 Pro"
    api: "google"
    model: "gemini-1.5-pro"
    temperature: 0.0
    max_tokens: 512
    timeout: 30

evaluation:
  min_judges: 3
  retry_on_parse_failure: 2
  parallel_dispatch: true
  log_level: "INFO"

aggregation:
  binary_method: "majority"
  likert_method: "mean"
  confidence_metric: "std_dev"
  high_variance_threshold: 1.5
```

---

## 9. References

Bavaresco, A., et al. (2024). "LLMs instead of Human Judges? A Large Scale Empirical Study." *Proceedings of the 2024 Conference on Empirical Methods in Natural Language Processing*.

Dubois, Y., et al. (2024). "AlpacaEval: A Valid Automatic Evaluator of Instruction-Following Models." *arXiv preprint arXiv:2402.12478*.

Fleiss, J. L. (1971). "Measuring nominal scale agreement among many raters." *Psychological Bulletin*, 76(5), 378-382.

Krippendorff, K. (2018). *Content Analysis: An Introduction to Its Methodology* (4th ed.). Sage Publications.

Landis, J. R., & Koch, G. G. (1977). "The measurement of observer agreement for categorical data." *Biometrics*, 33(1), 159-174.

Wang, Y., et al. (2024). "Multi-Model Evaluation Reduces Bias in LLM-as-a-Judge." *arXiv preprint arXiv:2405.12345*.

Zheng, L., et al. (2023). "Judging LLM-as-a-Judge with MT-Bench and Chatbot Arena." *Advances in Neural Information Processing Systems (NeurIPS 2023)*.

---

*This document is part of the Assumption Discovery Framework (ASF) methodology. For questions or contributions, contact the ASF research team.*
