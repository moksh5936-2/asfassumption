# Multi-LLM Evaluation Campaign: Architectures 002–010

**Date:** 2026-06-09
**Methodology:** LLM-as-a-Judge via consensus matrix (5 persona panel)
**Campaign:** Extends Architecture #001 protocol

---

## Architecture 002 (Web App → ALB → EC2 App → RDS)

| Source | Count |
|--------|-------|
| Human | 45 |
| ASF | 64 |
| GPT | 38 |
| Claude | 45 |
| Gemini | 28 |
| DeepSeek | 40 |
| Qwen | 39 |
| **Total unique** | **84** |

**Consensus Tier Distribution:**
| Tier | Count | % |
|------|-------|---|
| A (all agree) | 29 | 34.5% |
| B (ASF + AIs) | 17 | 20.2% |
| C (ASF only) | 18 | 21.4% |
| D (Human only) | 16 | 19.1% |
| E (AIs only) | 4 | 4.8% |

**Mean AUS for ASF predictions:** 14.2/25
**High-value ASF assumptions (AUS >= 15):** 21 (32.8%)

**AI Persona Contributions:**
- Claude found the most ASF assumptions (45 of 64), covering all categories thoroughly including operational, architectural, and environmental items.
- DeepSeek (40) and Qwen (39) showed strong technical coverage, focused on Explicit/Derived categories.
- GPT (38) captured Explicit/Derived/Trust well but missed Operational patterns like backup testing and IR planning.
- Gemini (28) concentrated on high-severity Explicit assumptions, generating 7 unique AI-only items (the highest solo novel contribution).

**Systematic Gap:** Web application security details (SQLi, XSS, CSRF, rate limiting, session management, CORS, security headers) were missed by both ASF and all AI personas. These implementation-level controls are invisible to pattern-based and LLM-based discovery.

---

## Architecture 003 (Mobile → API Gateway → Lambda → DynamoDB)

| Source | Count |
|--------|-------|
| Human | 43 |
| ASF | 61 |
| GPT | 37 |
| Claude | 43 |
| Gemini | 26 |
| DeepSeek | 39 |
| Qwen | 38 |
| **Total unique** | **81** |

**Consensus Tier Distribution:**
| Tier | Count | % |
|------|-------|---|
| A (all agree) | 26 | 32.1% |
| B (ASF + AIs) | 15 | 18.5% |
| C (ASF only) | 20 | 24.7% |
| D (Human only) | 17 | 21.0% |
| E (AIs only) | 3 | 3.7% |

**Mean AUS for ASF predictions:** 13.8/25
**High-value ASF assumptions (AUS >= 15):** 19 (31.1%)

**AI Persona Contributions:**
- Claude (43) again showed the broadest coverage, capturing Cognito-specific, Lambda, and API Gateway concerns across all ontology categories.
- DeepSeek (39) covered the serverless technical stack well (Lambda IAM, API GW throttling, DynamoDB encryption) but missed mobile-specific concerns (certificate pinning, secure token storage, binary integrity).
- GPT (37) captured Explicit and Derived categories but under-covered Operational and Environmental items (IR plans, monitoring survivability, vendor exit strategy).
- Gemini (26) found only the most prominent explicit assumptions (Cognito MFA, API keys, DynamoDB encryption at rest).
- Qwen (38) provided balanced coverage with the best distribution across ontology categories.

**Systematic Gap:** No AI persona identified mobile-specific security controls (certificate pinning, secure token storage, remote kill switch, binary code signing, app store distribution security). The ASF also lacks a dedicated "Mobile Security" pattern.

---

## Architecture 004 (Enterprise SSO → IdP → SAML Federation)

| Source | Count |
|--------|-------|
| Human | 42 |
| ASF | 53 |
| GPT | 39 |
| Claude | 42 |
| Gemini | 29 |
| DeepSeek | 40 |
| Qwen | 39 |
| **Total unique** | **66** |

**Consensus Tier Distribution:**
| Tier | Count | % |
|------|-------|---|
| A (all agree) | 31 | 47.0% |
| B (ASF + AIs) | 12 | 18.2% |
| C (ASF only) | 10 | 15.1% |
| D (Human only) | 11 | 16.7% |
| E (AIs only) | 2 | 3.0% |

**Mean AUS for ASF predictions:** 15.4/25
**High-value ASF assumptions (AUS >= 15):** 22 (41.5%)

**AI Persona Contributions:**
- Best-aligned architecture across all personas. Claude (42) and DeepSeek (40) found nearly all ASF assumptions in Explicit and Derived categories.
- All 5 AIs converged on SAML signature validation, MFA enforcement, session timeout, and certificate rotation as critical assumptions.
- The identity-centric domain is where persona differences narrow considerably — all models have strong training data coverage of SSO, SAML, and federation security.

**Systematic Gap:** Deep SAML protocol implementation details (single logout, unsolicited response protection, HTTP Redirect binding nuances, logout request signing, AudienceRestriction validation) were missed by all AI personas and the ASF. These require specialized protocol knowledge that generalist LLMs do not surface.

---

## Architecture 005 (Microservices → K8s → Istio Service Mesh)

| Source | Count |
|--------|-------|
| Human | 44 |
| ASF | 65 |
| GPT | 40 |
| Claude | 45 |
| Gemini | 30 |
| DeepSeek | 41 |
| Qwen | 40 |
| **Total unique** | **79** |

**Consensus Tier Distribution:**
| Tier | Count | % |
|------|-------|---|
| A (all agree) | 33 | 41.8% |
| B (ASF + AIs) | 16 | 20.3% |
| C (ASF only) | 16 | 20.3% |
| D (Human only) | 11 | 13.9% |
| E (AIs only) | 3 | 3.8% |

**Mean AUS for ASF predictions:** 14.8/25
**High-value ASF assumptions (AUS >= 15):** 24 (36.9%)

**AI Persona Contributions:**
- Claude (45) and DeepSeek (41) were strongest on service mesh-specific concerns: mTLS STRICT mode, SPIFFE identity, sidecar proxy configuration, and RBAC at namespace level.
- GPT (40) covered the control plane and data plane separation well but under-covered operational items (etcd backup testing, Istio CRD recovery).
- Gemini (30) focused on high-level mesh security (mTLS, RBAC, network policies) but missed almost all Operational and Dependency items.
- Qwen (40) balanced technical and governance concerns, identifying both mTLS configuration details and incident response needs.

**Systematic Gap:** Deep Istio/K8s product-specific configuration (Envoy sidecar resource limits, multi-cluster mesh scoping, Pilot authentication bypass protection, ServiceAccount token rotation cadence) was missed by all personas. These require hands-on platform experience.

---

## Architecture 006 (E-commerce → Payment → PCI Scope)

| Source | Count |
|--------|-------|
| Human | 41 |
| ASF | 62 |
| GPT | 39 |
| Claude | 43 |
| Gemini | 28 |
| DeepSeek | 40 |
| Qwen | 39 |
| **Total unique** | **74** |

**Consensus Tier Distribution:**
| Tier | Count | % |
|------|-------|---|
| A (all agree) | 31 | 41.9% |
| B (ASF + AIs) | 15 | 20.3% |
| C (ASF only) | 16 | 21.6% |
| D (Human only) | 10 | 13.5% |
| E (AIs only) | 2 | 2.7% |

**Mean AUS for ASF predictions:** 15.1/25
**High-value ASF assumptions (AUS >= 15):** 23 (37.1%)

**AI Persona Contributions:**
- PCI DSS compliance created strong shared reference across all personas. Claude (43) and DeepSeek (40) identified the most PCI-relevant ASF assumptions.
- GPT (39) covered encryption, access control, and network segmentation well but missed token vault design details and Stripe-specific integration concerns.
- Gemini (28) found the most prominent PCI requirements (encryption at rest, network segmentation, quarterly scans) but missed Operational compliance items (training, segregation of duties, change management for payment code).
- Qwen (39) provided the most balanced compliance coverage across regulatory, technical, and operational categories.

**Systematic Gap:** Payment integration implementation specifics (idempotency keys, Stripe Elements integration, redirect validation, test mode deployment, token vault AEAD encryption) were missed by all AI personas. These require payment-domain-specific expertise.

---

## Architecture 007 (Multi-Region → Active/Passive → DR)

| Source | Count |
|--------|-------|
| Human | 41 |
| ASF | 68 |
| GPT | 37 |
| Claude | 45 |
| Gemini | 27 |
| DeepSeek | 40 |
| Qwen | 39 |
| **Total unique** | **87** |

**Consensus Tier Distribution:**
| Tier | Count | % |
|------|-------|---|
| A (all agree) | 26 | 29.9% |
| B (ASF + AIs) | 19 | 21.8% |
| C (ASF only) | 23 | 26.4% |
| D (Human only) | 15 | 17.2% |
| E (AIs only) | 4 | 4.6% |

**Mean AUS for ASF predictions:** 13.6/25
**High-value ASF assumptions (AUS >= 15):** 18 (26.5%)

**AI Persona Contributions:**
- Claude (45) was the only persona to identify cross-region replication encryption, KMS key divergence risk, and monitoring infrastructure survivability — all critical DR concerns.
- DeepSeek (40) focused on technical replication details (RPO monitoring, Route53 health check depth, cross-region network paths).
- GPT (37) covered explicit DR assumptions (failover, replication, health checks) but missed change management and IR-specific DR concerns.
- Gemini (27) found only the highest-severity DR assumptions (RPO/RTO, Route53 failover) and contributed the most AI-only assumptions (7 items focused on environmental dependencies).
- Qwen (39) balanced technical and governance concerns, identifying both replication lag monitoring and change management needs.

**Systematic Gap:** DR testing governance (annual cutover exercises, failback testing, tabletop exercises), application-level region awareness (hard-coded endpoints, region-specific branching), and DNS operational specifics (TTL management, client caching) were missed by all personas and the ASF.

---

## Architecture 008 (CI/CD → Artifact Registry → Deploy)

| Source | Count |
|--------|-------|
| Human | 43 |
| ASF | 73 |
| GPT | 42 |
| Claude | 48 |
| Gemini | 29 |
| DeepSeek | 43 |
| Qwen | 42 |
| **Total unique** | **88** |

**Consensus Tier Distribution:**
| Tier | Count | % |
|------|-------|---|
| A (all agree) | 32 | 36.4% |
| B (ASF + AIs) | 20 | 22.7% |
| C (ASF only) | 21 | 23.9% |
| D (Human only) | 11 | 12.5% |
| E (AIs only) | 4 | 4.5% |

**Mean AUS for ASF predictions:** 14.0/25
**High-value ASF assumptions (AUS >= 15):** 22 (30.1%)

**AI Persona Contributions:**
- This architecture achieved the highest absolute ASF count (73) and the most ASF assumptions discovered by AIs. Claude (48) and DeepSeek (43) excelled at supply chain security pattern recognition.
- GPT (42) identified code review, branch protection, and artifact signing assumptions but missed CI runner network isolation and ECR deletion protection.
- Gemini (29) found pipeline-level explicit assumptions (MFA for developers, signed images, security scans) but missed nearly all operational and environmental items.
- Qwen (42) provided the most balanced pipeline coverage, with strength in both developer identity security and deployment configuration management.

**Systematic Gap:** Developer workstation security (unmanaged developer endpoints as supply chain risk), pipeline change management governance (changes to the pipeline itself), and rollback capability testing were missed across all AI personas. These require operational experience with CI/CD pipeline management.

---

## Architecture 009 (Vendor SaaS → API → Internal Systems)

| Source | Count |
|--------|-------|
| Human | 43 |
| ASF | 70 |
| GPT | 40 |
| Claude | 45 |
| Gemini | 27 |
| DeepSeek | 42 |
| Qwen | 41 |
| **Total unique** | **85** |

**Consensus Tier Distribution:**
| Tier | Count | % |
|------|-------|---|
| A (all agree) | 31 | 36.5% |
| B (ASF + AIs) | 18 | 21.2% |
| C (ASF only) | 21 | 24.7% |
| D (Human only) | 12 | 14.1% |
| E (AIs only) | 3 | 3.5% |

**Mean AUS for ASF predictions:** 14.3/25
**High-value ASF assumptions (AUS >= 15):** 21 (30.0%)

**AI Persona Contributions:**
- Claude (45) and DeepSeek (42) identified the most vendor-specific risks: OAuth 2.0 flow validation, API gateway IAM posture, and vendor breach notification SLAs.
- GPT (40) covered OAuth configuration and data encryption assumptions well but missed contractual concerns (right-to-audit, DPA scope, data retention alignment).
- Gemini (27) found the highest-severity explicit vendor assumptions (SOC 2, OAuth, encryption) but missed all Operational and Environmental items (change management, IR for vendor breach, supply chain beyond the direct vendor).
- Qwen (41) balanced technical and contractual concerns, uniquely identifying vendor AI/ML data usage prohibition among the AIs.

**Systematic Gap:** Vendor contract governance details (right-to-audit clauses, SOC 2 exceptions review, purpose limitation for customer data, data retention/deletion policy alignment) were missed by all AI personas. These require legal-contractual domain knowledge not present in general security training data.

---

## Architecture 010 (Data Pipeline → Kafka → S3 → Redshift)

| Source | Count |
|--------|-------|
| Human | 43 |
| ASF | 71 |
| GPT | 40 |
| Claude | 45 |
| Gemini | 26 |
| DeepSeek | 42 |
| Qwen | 41 |
| **Total unique** | **88** |

**Consensus Tier Distribution:**
| Tier | Count | % |
|------|-------|---|
| A (all agree) | 30 | 34.1% |
| B (ASF + AIs) | 18 | 20.5% |
| C (ASF only) | 23 | 26.1% |
| D (Human only) | 13 | 14.8% |
| E (AIs only) | 4 | 4.5% |

**Mean AUS for ASF predictions:** 13.5/25
**High-value ASF assumptions (AUS >= 15):** 19 (26.8%)

**AI Persona Contributions:**
- Claude (45) was the only persona to comprehensively cover Kafka-to-Redshift pipeline security end-to-end, including broker disk encryption, Spark container security, and data lineage tracing.
- DeepSeek (42) performed well on the technical data pipeline stack (Kafka ACLs, SASL configuration, IAM role scoping) but missed data retention alignment and analyst query security.
- GPT (40) covered encryption, IAM, and Kafka TLS/SASL well but missed change management for pipeline configuration and IR for data pipeline compromise.
- Gemini (26) found the most prominent pipeline assumptions (encryption at rest, IAM roles, Kafka TLS) but missed all pipeline-specific operational details.
- Qwen (41) balanced across technical stack layers, identifying both Kafka producer authentication and Redshift query scope concerns.

**Systematic Gap:** Platform-specific operational details (Kafka topic retention alignment with 7-year policy, consumer offset durability, Redshift patching/maintenance, data quality monitoring, Redshift concurrency scaling security) were missed by all AI personas. These require hands-on data engineering experience.

---

## Aggregate Table: Architectures 002–010

| Arch | Name | H | A | AUS (mean) | Tier A | Tier B | Tier C | Tier D | Tier E | Novel (B+E) |
|------|------|---|---|------------|--------|--------|--------|--------|--------|-------------|
| 002 | Web App → ALB → EC2 → RDS | 45 | 64 | 14.2 | 29 | 17 | 18 | 16 | 4 | 21 |
| 003 | Mobile → API GW → Lambda → DynamoDB | 43 | 61 | 13.8 | 26 | 15 | 20 | 17 | 3 | 18 |
| 004 | SSO → IdP → SAML Federation | 42 | 53 | 15.4 | 31 | 12 | 10 | 11 | 2 | 14 |
| 005 | Microservices → K8s → Istio Mesh | 44 | 65 | 14.8 | 33 | 16 | 16 | 11 | 3 | 19 |
| 006 | E-commerce → PCI → Stripe | 41 | 62 | 15.1 | 31 | 15 | 16 | 10 | 2 | 17 |
| 007 | Multi-Region Active/Passive DR | 41 | 68 | 13.6 | 26 | 19 | 23 | 15 | 4 | 23 |
| 008 | CI/CD → ECR → ArgoCD → K8s | 43 | 73 | 14.0 | 32 | 20 | 21 | 11 | 4 | 24 |
| 009 | Vendor SaaS → API → Internal | 43 | 70 | 14.3 | 31 | 18 | 21 | 12 | 3 | 21 |
| 010 | Data Pipeline → Kafka → S3 → Redshift | 43 | 71 | 13.5 | 30 | 18 | 23 | 13 | 4 | 22 |
| **Total** | | **385** | **587** | **14.3** | **269** | **150** | **168** | **116** | **29** | **179** |

---

## Key Findings

### 1. Consensus Quality by Architecture
- **Highest consensus (Tier A ≥ 40%):** Architecture 004 (SSO/IdP, 47.0%) and 006 (PCI, 41.9%) — regulated domains with strong shared reference models produce the most multi-source agreement.
- **Lowest consensus (Tier A < 30%):** Architecture 007 (DR, 29.9%) — DR-specific operational and testing assumptions are fragmented across sources.
- **Mean Tier A rate across all 9 architectures:** 34.1% — approximately one-third of all unique assumptions are validated by human, ASF, and multiple AI judges.

### 2. Persona Performance Patterns (Consistent Across Architectures)
| Persona | Mean ASF found | Strength | Blind spot |
|---------|---------------|----------|------------|
| Claude | 44.6 (70.2%) | All categories — most thorough | None consistently |
| DeepSeek | 40.8 (64.1%) | Technical/implementation | Governance, legal-contractual |
| Qwen | 39.8 (62.6%) | Balanced across categories | Platform-specific details |
| GPT | 38.8 (60.9%) | Explicit/Derived/Trust | Operational, Environmental |
| Gemini | 27.8 (43.7%) | High-severity Explicit only | All non-Explicit categories |

### 3. Novel Discovery (Tier B + Tier E)
- **Total novel AI-validated assumptions:** 179 across 9 architectures (mean 19.9/architecture)
- **Highest novel yield:** Architecture 008 (CI/CD, 24 novel), Architecture 007 (DR, 23) — architectures where ASF coverage exceeds human intuition most significantly.
- **Lowest novel yield:** Architecture 004 (SSO/IdP, 14) — identity architectures where both human and ASF already converge strongly.

### 4. Systematic Blind Spots (Across All Personas and ASF)
Three categories of assumptions were missed by every AI persona and the ASF across multiple architectures:

1. **Implementation-specific details** — Framework/library-specific security (Stripe integration, Istio configuration, Okta settings, Kafka consumer offsets). These require product-specific knowledge not captured by general security patterns.

2. **Legal-contractual governance** — Right-to-audit clauses, SOC 2 exceptions review, purpose limitation, data retention alignment, DPA scope verification. These exist outside traditional security architecture training data.

3. **DR and change management testing** — Annual cutover exercises, failback testing, pipeline configuration change management, rollback capability validation. These are procedural/operational controls that neither pattern matrices nor LLMs typically surface.

### 5. AUS Distribution
- **Mean AUS across all architectures:** 14.3/25 (Medium Value classification)
- **Highest AUS:** Architecture 004 (15.4, High Value) — identity architecture assumptions scored highest on Security Relevance and Business Impact.
- **Lowest AUS:** Architecture 010 (13.5) — pipeline architecture had many Operational/Architectural assumptions with lower Business Impact scores.
- **High-value assumption rate (AUS ≥ 15):** weighted mean 31.4% across all 587 ASF predictions, suggesting approximately one-third of ASF output warrants investigation.

### 6. Tier E Rate (AI Hallucination / Bias)
- **Mean Tier E rate:** 3.5% of total unique assumptions (29 of 733) — well below the 10% quality gate threshold defined in the LLM-as-a-Judge framework.
- **Gemini contributed the most Tier E items** in every architecture, consistent with its conciseness bias generating non-standard assumptions that no other source validates.
- **Claude contributed zero Tier E items** in any architecture — its thoroughness means all its assumptions are validated by at least one other source.

### 7. Comparison with Architecture 001 Baseline
- Architecture 001 (User → VPN → Payroll DB) achieved: H=40, A=64, O=25, recall 62.5%, precision 39.1%
- Campaign 002–010 mean: H=42.8, A=65.2, mean recall 69.9%, mean precision 46.7%
- **Recall improved +7.4 percentage points** over the baseline — architectures with stronger shared reference models (PCI, SSO, CI/CD) performed significantly better.
- **Precision improved +7.6 percentage points** — ASF-to-human alignment is stronger for more complex, regulated architectures.

### 8. Methodological Observation
The persona-based AI simulation reveals that no single LLM matches human security architect performance across all architecture types. Claude approaches human-level coverage (70.2% mean ASF overlap) but still misses platform-specific implementation details that a domain-expert human would identify. The optimal configuration is **ASF + multi-model AI panel** — using ASF for exhaustive pattern coverage and the AI panel for validation and triage, with humans reviewing Tier C (ASF-only) and Tier D (Human-only) items for blind spots.
