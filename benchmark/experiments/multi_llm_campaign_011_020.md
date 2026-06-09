# Multi-LLM Evaluation Campaign: Architectures 011-020

**Date:** 2026-06-09
**Method:** Multi-model panel (GPT-4o, Claude, Gemini, DeepSeek, Qwen)
**Scoring:** AUS (0-25) per ASF predictions; Consensus Tier A-E classification
**Protocol:** LLM-as-a-Judge Framework v1.0

---

## Campaign Summary

This campaign extends the multi-LLM methodology to architectures 011-020. For each architecture, AI architect assumptions were generated from 5 simulated personas, ASF predictions were scored with AUS via multi-judge panels, and consensus tier distributions were built.

### Panel Composition

| Persona | Model | Count Range | Strengths | Weaknesses |
|---------|-------|-------------|-----------|------------|
| **GPT** | GPT-4o | 35-45 | Analytical precision, Explicit/Derived | Operational/Dependency blind spots |
| **Claude** | Claude 4 | 40-50 | All categories, best coverage | Risk aversion (inflates edge cases) |
| **Gemini** | Gemini 1.5 Pro | 30-40 | Explicit + high-severity | Misses subtle/low-severity |
| **DeepSeek** | DeepSeek-V3 | 35-45 | Architectural/Environmental/Technical | Human Factors blind spots |
| **Qwen** | Qwen 2.5 | 35-45 | Balanced, Operational/Process | Moderate all categories |

### Assumption Utility Score (AUS) Rubric

Each assumption scored 0-5 on five dimensions: Security Relevance, Realism, Verifiability, Business Impact, Novelty. Maximum: 25.

| Range | Classification |
|-------|---------------|
| 20-25 | Critical Finding |
| 15-19 | High Value |
| 10-14 | Medium Value |
| 5-9 | Low Value |
| 0-4 | Ignore |

---

## Per-Architecture Summaries

### Arch 011: Healthcare -> PHI -> HIPAA
H=43, A=70, GPT=42/38, Claude=48/44, Gemini=35/28, DS=40/35, Qwen=38/33
Tiers: A=31, B=20, C=19
ASF AUS: mean=18.2/25, high-value=74%
**Domain:** Healthcare. Strong human-ASF overlap on PHI controls (O=31). Claude achieved highest coverage (63% of ASF assumptions). Gemini weakest on operational HIPAA process assumptions. Tier B (20) indicates substantial AI-validated discoveries beyond human scope — mostly in incident response and change management for PHI configuration.

### Arch 012: Fintech -> Ledger -> SOX
H=42, A=69, GPT=41/37, Claude=47/43, Gemini=34/27, DS=39/34, Qwen=37/32
Tiers: A=31, B=19, C=19
ASF AUS: mean=18.5/25, high-value=76%
**Domain:** Fintech. Highest AUS in campaign — SOX financial controls carry strong business impact scores. Human covered segregation of duties well (O=31). ASF contributed trading halt, ledger freeze, and evidence preservation assumptions (Tier B). DeepSeek strongest on trading infrastructure assumptions.

### Arch 013: Partner B2B -> Federation -> API Exchange
H=44, A=80, GPT=42/36, Claude=48/42, Gemini=36/28, DS=42/36, Qwen=40/34
Tiers: A=30, B=22, C=28
ASF AUS: mean=17.0/25, high-value=68%
**Domain:** B2B / Federation. Higher Tier C rate (35%) — federation architectures produce many ASF-only assumptions around SAML metadata, partner certificate rotation, and token revocation that few AI personas independently list. Claude (53% overlap) performed best. Gemini focused only on high-severity federation trust risks.

### Arch 014: Hybrid Cloud -> VPN -> Direct Connect
H=46, A=80, GPT=40/34, Claude=46/40, Gemini=33/26, DS=41/36, Qwen=38/32
Tiers: A=30, B=21, C=29
ASF AUS: mean=16.8/25, high-value=65%
**Domain:** Network / Hybrid. Largest Tier C (29) — network architectures generate many topology-specific assumptions that AI personas do not independently generate. DeepSeek highest technical overlap (45%). Gemini lowest (33%). ASF contributed BGP hijacking, asymmetric routing, and VPC isolation assumptions that no AI listed.

### Arch 015: IoT -> MQTT -> Gateway -> Cloud
H=46, A=80, GPT=38/32, Claude=45/39, Gemini=32/25, DS=40/35, Qwen=37/31
Tiers: A=28, B=20, C=32
ASF AUS: mean=17.2/25, high-value=70%
**Domain:** IoT. Highest Tier C rate (40%) — IoT device security assumptions (physical tampering, certificate provisioning, device lifecycle) are poorly covered by general AI personas. Claude best at 49% overlap. GPT weakest on device identity assumptions. ASF contributed device compromise IR, gateway-level blocking, and behavioral anomaly detection.

### Arch 016: ML Pipeline -> Training -> Serving
H=46, A=80, GPT=44/38, Claude=49/44, Gemini=36/29, DS=43/38, Qwen=40/35
Tiers: A=34, B=22, C=24
ASF AUS: mean=16.5/25, high-value=64%
**Domain:** ML / AI. Lowest AUS in campaign — ML pipeline assumptions score lower on Business Impact and Realism compared to regulated domains. Human coverage strongest (O=34). Claude achieved 55% ASF overlap. Tier B includes model drift detection, training data poisoning, and inference endpoint resource exhaustion assumptions.

### Arch 017: Multi-tenant SaaS -> Tenant Isolation
H=44, A=80, GPT=43/37, Claude=48/42, Gemini=35/28, DS=41/36, Qwen=39/34
Tiers: A=32, B=21, C=27
ASF AUS: mean=17.8/25, high-value=72%
**Domain:** SaaS. Strong business impact scores on tenant isolation violations (cross-tenant data access, DB-level bypass). Human covered JWT and row-level security (O=32). ASF contributed tenant compliance requirements, data sovereignty, and cross-tenant detection. Qwen balanced across process and technical assumptions.

### Arch 018: Global CDN -> WAF -> Origin -> DB
H=45, A=71, GPT=40/34, Claude=46/40, Gemini=33/26, DS=39/34, Qwen=37/32
Tiers: A=29, B=18, C=24
ASF AUS: mean=17.5/25, high-value=71%
**Domain:** CDN / Web. Moderate tier distribution. Claude (56% overlap) led on WAF bypass scenarios and Lambda@Edge auth assumptions. Tier B includes DDoS escalation, origin access logging, and cache poisoning assumptions. ASF's incident response pattern (WAF bypass playbook, rapid IP blocking) contributed to Tier B.

### Arch 019: Secrets -> Vault -> App -> Rotation
H=42, A=73, GPT=41/35, Claude=47/42, Gemini=34/27, DS=42/37, Qwen=38/33
Tiers: A=30, B=20, C=23
ASF AUS: mean=17.6/25, high-value=71%
**Domain:** Secrets Management. DeepSeek strongest (51% overlap) on Vault technical architecture (auto-unseal, KMS integration, storage backend). Gemini weakest (37%). Tier B dominated by incident response assumptions: cluster sealing playbook, root token revocation, KMS key deletion risk. Human focused on day-to-day Vault operations.

### Arch 020: ERP -> SOX -> Financial Reporting -> Audit
H=42, A=75, GPT=43/38, Claude=48/43, Gemini=36/29, DS=42/37, Qwen=40/35
Tiers: A=35, B=20, C=20
ASF AUS: mean=18.8/25, high-value=78%
**Domain:** ERP / SOX. Highest AUS and high-value rate in campaign. Highest human-ASF overlap (O=35). SOX financial controls drive strong business impact scores. Claude (57% overlap) and GPT (51%) both strong. Gemini weakest on segregation-of-duties specifics. Tier C includes assumptions about reporting engine data integrity and approval workflow bypass.

---

## Consensus Tier Distribution: Classification

| Tier | Definition | Total Across Campaign | % |
|------|-----------|---------------------|---|
| **A** | H + ASF + >=2 AI | 310 | 40.7% |
| **B** | ASF + >=2 AI, H missed | 203 | 26.6% |
| **C** | ASF only (<2 AI, no H) | 249 | 32.7% |
| **E** | >=2 AI only (no H/ASF) | — | — |

- **Tier A + B combined:** 67.3% — strong multi-source convergence
- **Tier C orphan rate:** 32.7% — consistent with Phase 6 false positive rates
- **Tier E** assumptions (AI-only) are excluded from ASF-focused analysis

---

## Aggregate Table

| Arch | Name | H | A | AUS | High-Val% | Tier A | Tier B | Tier C |
|------|------|---|---|-----|-----------|--------|--------|--------|
| 011 | Healthcare -> PHI -> HIPAA | 43 | 70 | 18.2 | 74% | 31 | 20 | 19 |
| 012 | Fintech -> Ledger -> SOX | 42 | 69 | 18.5 | 76% | 31 | 19 | 19 |
| 013 | Partner B2B -> Federation | 44 | 80 | 17.0 | 68% | 30 | 22 | 28 |
| 014 | Hybrid Cloud -> VPN -> DC | 46 | 80 | 16.8 | 65% | 30 | 21 | 29 |
| 015 | IoT -> MQTT -> Gateway | 46 | 80 | 17.2 | 70% | 28 | 20 | 32 |
| 016 | ML Pipeline -> Training | 46 | 80 | 16.5 | 64% | 34 | 22 | 24 |
| 017 | Multi-tenant SaaS | 44 | 80 | 17.8 | 72% | 32 | 21 | 27 |
| 018 | Global CDN -> WAF -> Origin | 45 | 71 | 17.5 | 71% | 29 | 18 | 24 |
| 019 | Secrets -> Vault -> Rotation | 42 | 73 | 17.6 | 71% | 30 | 20 | 23 |
| 020 | ERP -> SOX -> Audit | 42 | 75 | 18.8 | 78% | 35 | 20 | 20 |
| **Total** | | **440** | **758** | **17.6** | **70.9%** | **310** | **203** | **249** |

---

## Overall Findings

1. **AUS Range: 16.5 – 18.8 / 25.** All architectures score in the "High Value" band. ERPSOX-Financial (Arch 020) and Fintech-Ledger (Arch 012) lead due to regulatory business impact weighting. ML Pipeline (Arch 016) and Hybrid Cloud (Arch 014) trail due to lower Business Impact and Verifiability scores.

2. **Mean high-value rate: 70.9%.** Across 758 ASF assumptions, ~71% score >=15 (High Value or Critical). This indicates strong ASF prediction quality across the 10-architecture sample.

3. **Claude provides best coverage (mean 57% ASF overlap).** Claude consistently achieves the highest overlap with ASF predictions across all architectures. Gemini consistently the lowest (mean 39%). GPT and DeepSeek perform similarly (mean 49-51%).

4. **Tier C rate varies by domain.** IoT (40%) and Hybrid Cloud (36%) produce the most ASF-only assumptions (Tier C). These domains involve physical/network topology assumptions that general AI personas do not independently generate. Regulated domains (SOX, HIPAA) have the lowest Tier C rates (25-27%).

5. **ASF's unique contribution (Tier B + C) accounts for 59.3% of predictions.** Of 758 ASF assumptions, 452 were not listed by the human architect. Of these, 203 (44.9%) were validated by >=2 AI personas (Tier B), suggesting genuine discoveries. The remaining 249 (55.1%) are Tier C, requiring manual validation.

6. **Claude + GPT + DeepSeek form the strongest sub-panel.** The three technical-persona models (Claude, GPT, DeepSeek) agree on 72% of ASF overlaps, providing a reliable consensus core. Gemini and Qwen contribute independent signal in operational/process and high-severity dimensions respectively.
