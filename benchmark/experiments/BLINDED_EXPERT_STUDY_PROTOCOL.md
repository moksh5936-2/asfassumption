# Blinded Expert Study: Do ASF's Unique Assumptions Have Real-World Value?

**Purpose:** Determine whether ASF-only assumptions (those that no leading AI independently derived) are rated as valid, important, and actionable by experienced security architects — compared to assumptions from Claude and Gemini.

**Why this matters:** 207 ASF assumptions across 5 architectures were not independently produced by any of 4 AI models. But "unique" does not equal "valuable." This study measures whether those 207 represent genuine blind-spot discovery or methodological over-generation.

---

## Study Design

### Participants

Target 5-10 senior security architects / CISOs / AppSec leads.

Recruitment criteria:
- 5+ years in security architecture or equivalent
- Experience across at least 2 of the 5 architecture domains (network, identity, cloud/K8s, healthcare, financial/SOX)
- Not involved in the ASF project

### Materials

**150 blinded assumptions**, drawn from three source buckets:

| Bucket | Source | Count | Description |
|--------|--------|-------|-------------|
| **ASF-Only** | ASF predictions that none of 4 models (Claude, GPT, Gemini, Gemma) independently derived | 50 | The strongest ASF-unique findings — novel security assumptions no AI independently listed |
| **Claude** | Assumptions from Claude's independent derivation outputs | 50 | Architecture-specific assumptions Claude produced during blind testing |
| **Gemini** | Assumptions from Gemini's independent derivation outputs | 50 | Architecture-specific assumptions Gemini produced during blind testing |

**Total:** 150 assumptions, randomly ordered, source labels removed.

### Architecture Context

Each participant sees 5 architecture diagrams (same 5 used in all experiments):
1. VPN → Payroll DB
2. SSO/IdP → SAML Federation
3. K8s/Istio Service Mesh
4. Healthcare → PHI → HIPAA
5. ERP → SOX → Audit

Assumptions are drawn from across all 5 architectures and mixed.

### Rating Scale (per assumption)

Each assumption is rated on 4 dimensions, 5-point Likert:

| # | Question | 1 | 2 | 3 | 4 | 5 |
|---|----------|---|---|---|---|---|
| Q1 | **Validity:** Is this a real security assumption? | Invalid | Weakly valid | Possibly valid | Valid | Definitely valid |
| Q2 | **Importance:** How important is it? | Trivial | Low | Medium | High | Critical |
| Q3 | **Review value:** Would you want this flagged in an architecture review? | Noise | Low value | Nice-to-have | Important | Essential |
| Q4 | **Risk:** Would missing this create measurable security risk? | No risk | Minimal | Moderate | Significant | Critical |

### Scoring

| Metric | Formula | Interpretation |
|--------|---------|---------------|
| **Mean score** | (Q1+Q2+Q3+Q4) / 4 per assumption | 1-5 scale, higher = more valuable |
| **Total score** | Q1+Q2+Q3+Q4 | Max 20 per assumption |
| **Bucket mean** | Average total score across all 50 assumptions in bucket | Compare ASF-only vs Claude vs Gemini |
| **Pass rate** | % of assumptions scoring ≥ 12/20 | Acceptability threshold |
| **Strong pass rate** | % of assumptions scoring ≥ 16/20 | High-value threshold |
| **Q3 essential rate** | % scoring 4+ on "Would you want this flagged?" | Practitioner demand signal |
| **Inter-rater agreement** | Standard deviation across participants | Reliability check |

### Success Criteria

| Criterion | Target | What It Means |
|-----------|--------|---------------|
| ASF-only mean score | ≥ 12/20 | ASF-unique assumptions are at least moderately valuable |
| ASF-only vs Claude gap | Non-inferior (within 2 pts) | ASF is not noise compared to human-expert-level assumptions |
| ASF-only Q3 essential rate | ≥ 40% | Practitioners actively want these flagged |
| Cross-arch consistency | No single architecture dominates ASF-only high scores | ASF value is general, not architecture-specific |

### Recruiting

1. Identify 10-15 senior security practitioners through professional network
2. Expect 50-70% acceptance rate → 5-10 completed sessions
3. Each session: 30-45 minutes (150 assumptions × ~15 seconds each + setup)
4. Format: Google Form or structured interview with digital score sheet
5. Compensation: Optional ($50-100 gift card for practitioners not affiliated with the project)

---

## Assumption Selection Criteria

### For ASF-Only (50 assumptions)

Select the 50 strongest ASF-unique assumptions from the 207 that no model independently derived. Priority:

1. **High AUS** (from multi-LLM campaign) — prefer assumptions with AUS ≥ 15
2. **Cross-architecture** — prefer concepts that appear across multiple architectures
3. **Domain diversity** — ensure coverage across all 7 ASF-unique domains:
   - Third-party dependency (vendor risk, exit strategy)
   - Identity lifecycle (joiner/mover/leaver, recert, service accounts)
   - Incident response (playbooks, forensic preservation, notification)
   - Availability & resilience (SPOF, DR, offline procedures)
   - Data classification (classification, flow diagrams, dev/staging)
   - Encryption governance (key rotation, key policy, temp storage)
   - Governance/compliance (SoD testing, recert auditability, retention)

### For Claude (50 assumptions)

Select 50 assumptions from Claude's independent derivation outputs:

1. **Architecture-specific** — prefer Claude's strongest specific findings (e.g., "SAML signing keys must be rotated", "Citadel CA root key must be protected")
2. **Diverse sources** — draw from across all 5 architectures
3. **Comparable specificity** — match the specificity level of ASF assumptions

### For Gemini (50 assumptions)

Select 50 from Gemini's outputs:

1. **Architecture-specific** — Gemini's strongest findings
2. **Diverse sources** — draw from across all 5 architectures
3. **Match specificity** — Gemini outputs are shorter (4-5 per arch), so may need to use all available

---

## Data Collection Instrument

### Participant Form

```
## Participant Background

Role: _________________________
Years in security: _________
Primary domains: _________________________
Have you ever conducted architecture security reviews? Y / N
If yes, approximately how many? _________

---

## Rating Instructions

You will see 150 assumptions drawn from 5 architecture diagrams.
Each assumption is something that MUST be true for the architecture to be secure, but is NOT explicitly stated in the documented policy.

For each assumption, rate on 4 dimensions:

Q1 - Is this a real assumption? (1=Invalid, 5=Definitely valid)
Q2 - How important is it? (1=Trivial, 5=Critical)
Q3 - Would you want this in a review? (1=Noise, 5=Essential)
Q4 - Would missing it create risk? (1=No risk, 5=Critical)

There are no right or wrong answers. Your expert judgment is the measurement.

---

## Assumption Ratings

[150 randomized items, format below per item]

### Item A-042

**Assumption:** Database credentials are rotated on a regular cadence and immediately upon suspicion of compromise.

| Q1: Valid? | 1 | 2 | 3 | 4 | 5 |
| Q2: Important? | 1 | 2 | 3 | 4 | 5 |
| Q3: Flag in review? | 1 | 2 | 3 | 4 | 5 |
| Q4: Risk if missed? | 1 | 2 | 3 | 4 | 5 |

---
```

---

## Data Analysis Plan

After collecting 5-10 responses:

1. **Clean data** — remove any participant who scored > 90% uniformly (straight-lining)
2. **Compute per-assumption mean** — average across all raters
3. **Compute bucket means** — ASF-only vs Claude vs Gemini
4. **Compute pass rates** — % scoring ≥ 12/20 per bucket
5. **Compute Q3 essential rate** — % scoring 4+ on "Would you flag this?"
6. **Inter-rater reliability** — Fleiss' kappa or ICC
7. **Domain analysis** — which assumption domains score highest within ASF-only
8. **Top 10 list** — the 10 highest-scoring assumptions across all buckets

### Decision Framework

| If ASF-only score... | Then... |
|---------------------|---------|
| ≥ 14/20 (Strong) | ASF has strong practitioner value. Proceed to ASF v2 design. |
| 12-13.9/20 (Moderate) | ASF has value but needs refinement. Focus on high-scoring domains. |
| 10-11.9/20 (Weak) | ASF assumptions are marginally useful. Reconsider methodology. |
| < 10/20 (Fail) | ASF-unique assumptions are not valuable to practitioners. Pivot required. |

---

## Pre-Selected Assumption Pool

### ASF-Only (50 assumptions — none of 4 AIs independently derived)

*To be filled with specific assumptions from the 207 ASF-only concepts, prioritized by AUS score and domain diversity.*

### Claude (50 assumptions — from Claude's blind derivation outputs)

*To be filled with specific assumptions from Claude's per-architecture outputs.*

### Gemini (50 assumptions — from Gemini's blind derivation outputs)

*To be filled with specific assumptions from Gemini's per-architecture outputs.*

---

## Execution Timeline

| Step | Owner | Duration |
|------|-------|----------|
| 1. Curate 150 assumption pool | Researcher | 1 day |
| 2. Build blinded instrument (randomized, labeled A-001 to A-150) | Researcher | 1 day |
| 3. Recruit 10-15 security architects | Researcher | 5-10 days |
| 4. Run 5-10 sessions (30-45 min each) | Researcher | 3-5 days |
| 5. Analyze data | Researcher | 1 day |
| 6. Write findings | Researcher | 1 day |
| **Total** | | **2-3 weeks** |
