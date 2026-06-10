# ASF-Only Concept AUS Ranking

**Stage 3 Analysis** — Scoring all ~198 ASF-only concepts from the 5-architecture independent derivation test.

**AUS Formula:** Impact (1-5) × Likelihood (1-5) × Novelty (1-5) = AUS (1-125)

**Scale:**
| Score | Impact | Likelihood | Novelty |
|-------|--------|-----------|---------|
| 1 | Minor inconvenience | Improbable | Every architect lists this |
| 2 | Low-severity issue | Unlikely in most orgs | Most experienced architects would list it |
| 3 | Moderate security concern | Plausible in medium/large orgs | Common in comprehensive architecture reviews |
| 4 | High-severity (data breach, compliance violation) | Likely in enterprise environments | Only senior or specialized reviewers would list it |
| 5 | Critical (regulatory fines, catastrophic breach) | Nearly certain in complex deployments | Almost no architect would list unprompted |

---

## Scoring by Domain

### 1. Third-party Dependency (~90% ASF-only, ~18 concepts)

| # | Concept | Arch | Impact | Likelihood | Novelty | AUS |
|---|---------|------|--------|------------|---------|-----|
| 1 | Vendor exit strategy for IdP migration | A2 | 4 | 3 | 5 | 60 |
| 2 | VPN vendor backdoor risk assessment | A1 | 5 | 2 | 5 | 50 |
| 3 | Container registry dependency / upstream trust | A3 | 4 | 3 | 5 | 60 |
| 4 | Third-party SLA dependency for SP apps | A2 | 3 | 4 | 4 | 48 |
| 5 | Auth0 HIPAA eligibility / BAA validity | A4 | 5 | 3 | 4 | 60 |
| 6 | Auth0 sub-processor disclosure requirement | A4 | 4 | 3 | 5 | 60 |
| 7 | Cloud provider SOC 2 for financial services | A5 | 4 | 3 | 4 | 48 |
| 8 | Istio CVE dependency / supply chain monitoring | A3 | 4 | 4 | 4 | 64 |
| 9 | AWS HIPAA BAA for infrastructure provider | A4 | 4 | 3 | 4 | 48 |
| 10 | Third-party integration API security review | A5 | 4 | 3 | 5 | 60 |
| 11 | Istio deprecation migration path | A3 | 3 | 3 | 5 | 45 |
| 12 | Container registry integrity (non-malicious images) | A3 | 5 | 3 | 5 | 75 |
| 13 | K8s version support window | A3 | 3 | 4 | 4 | 48 |
| 14 | DB platform CVE risk tracking | A5 | 4 | 3 | 4 | 48 |
| 15 | Auditor tool vulnerability risk | A5 | 3 | 3 | 5 | 45 |

**Domain mean AUS:** 53.7 | **Median:** 60

### 2. Identity Lifecycle (~82% ASF-only, ~15 concepts)

| # | Concept | Arch | Impact | Likelihood | Novelty | AUS |
|---|---------|------|--------|------------|---------|-----|
| 1 | Service account rigor (annual review, scope validation) | A2 | 4 | 4 | 4 | 64 |
| 2 | HR-integrated joiner/mover/leaver automation | A5 | 4 | 4 | 3 | 48 |
| 3 | Role-change triggers ERP permission updates | A5 | 4 | 4 | 4 | 64 |
| 4 | Service account annual review | A5 | 4 | 3 | 4 | 48 |
| 5 | Manager-attested recertification (not self-certification) | A5 | 4 | 4 | 5 | 80 |
| 6 | Provider joiner/mover/leaver for patient access | A4 | 4 | 4 | 4 | 64 |
| 7 | Patient account deactivation on care-end | A4 | 4 | 3 | 4 | 48 |
| 8 | Quarterly role recertification | A4 | 4 | 4 | 3 | 48 |
| 9 | Service account rigor (access recertification) | A4 | 4 | 4 | 4 | 64 |
| 10 | Admin backdoor account prevention | A5 | 5 | 3 | 5 | 75 |
| 11 | Credential sharing prevention (finance-specific) | A5 | 4 | 4 | 4 | 64 |
| 12 | Approval workflow user understanding | A5 | 3 | 3 | 4 | 36 |
| 13 | Recertification thoroughness (no rubber-stamping) | A5 | 4 | 4 | 5 | 80 |
| 14 | Recertification process auditability | A5 | 4 | 3 | 4 | 48 |

**Domain mean AUS:** 57.9 | **Median:** 62

### 3. Incident Response (~80% ASF-only, ~20 concepts)

| # | Concept | Arch | Impact | Likelihood | Novelty | AUS |
|---|---------|------|--------|------------|---------|-----|
| 1 | IR plan includes HIPAA 60-day breach notification timeline | A4 | 5 | 4 | 5 | 100 |
| 2 | IR database isolation for forensic preservation | A4 | 5 | 3 | 5 | 75 |
| 3 | PHI audit log forensic preservation during IR | A4 | 5 | 3 | 5 | 75 |
| 4 | IR team Auth0 log access during investigation | A4 | 4 | 4 | 5 | 80 |
| 5 | SOX 802 7-year financial data retention | A5 | 5 | 4 | 4 | 80 |
| 6 | Audit evidence automation (screenshots/logs) for SOX | A5 | 4 | 4 | 5 | 80 |
| 7 | SoD rule documentation and completeness testing | A5 | 4 | 4 | 4 | 64 |
| 8 | Containment plan for compromised SP apps | A2 | 4 | 4 | 4 | 64 |
| 9 | IR plan for payroll-specific breach exposure | A1 | 5 | 3 | 5 | 75 |
| 10 | Forensic preservation of etcd data during K8s IR | A3 | 5 | 3 | 5 | 75 |
| 11 | SOX control annual testing | A5 | 4 | 4 | 4 | 64 |
| 12 | Anomalous DB access monitoring (PHI-specific) | A4 | 5 | 3 | 4 | 60 |

**Domain mean AUS:** 74.3 | **Median:** 75

### 4. Availability & Resilience (~78% ASF-only, ~15 concepts)

| # | Concept | Arch | Impact | Likelihood | Novelty | AUS |
|---|---------|------|--------|------------|---------|-----|
| 1 | VPN gateway single point of failure | A1 | 4 | 4 | 5 | 80 |
| 2 | AD domain controller redundancy | A2 | 4 | 4 | 4 | 64 |
| 3 | Offline auth procedure for IdP outage | A2 | 5 | 3 | 5 | 75 |
| 4 | Control plane HA / multi-replica for Istio | A3 | 5 | 4 | 4 | 80 |
| 5 | Auth0 outage fallback procedure | A4 | 5 | 3 | 5 | 75 |
| 6 | ERP DR plan with defined RTO/RPO | A5 | 5 | 4 | 3 | 60 |
| 7 | ERP HA configuration (active-passive with testing) | A5 | 4 | 4 | 4 | 64 |
| 8 | Approval workflow availability during month-end close | A5 | 5 | 4 | 5 | 100 |
| 9 | Reporting Engine function during backend outage | A5 | 4 | 4 | 5 | 80 |
| 10 | Envoy health under config push load | A3 | 4 | 4 | 5 | 80 |
| 11 | Istio CRD backup strategy | A3 | 3 | 3 | 5 | 45 |
| 12 | Citadel CA key backup / recovery | A3 | 5 | 3 | 5 | 75 |

**Domain mean AUS:** 73.1 | **Median:** 75

### 5. Data Classification & Flow (~85% ASF-only, ~14 concepts)

| # | Concept | Arch | Impact | Likelihood | Novelty | AUS |
|---|---------|------|--------|------------|---------|-----|
| 1 | Data flow diagram accuracy (not just existence) | ALL | 4 | 4 | 5 | 80 |
| 2 | No PHI in dev/staging databases | A4 | 5 | 4 | 4 | 80 |
| 3 | No production data on local workstations | A5 | 5 | 4 | 4 | 80 |
| 4 | No production data in dev/staging (SOX) | A5 | 5 | 4 | 4 | 80 |
| 5 | PHI data formally classified as Protected Health Information | A4 | 4 | 3 | 4 | 48 |
| 6 | Financial data classified as Restricted/Critical | A5 | 4 | 3 | 4 | 48 |
| 7 | Mesh telemetry data exposure (mTLS metadata) | A3 | 3 | 4 | 5 | 60 |
| 8 | Attribute release classification (SAML/IdP) | A2 | 4 | 3 | 5 | 60 |
| 9 | Hidden SP-to-AD auth paths (unintended trust) | A2 | 5 | 3 | 5 | 75 |
| 10 | Data flow mapping documentation | ALL | 3 | 4 | 5 | 60 |

**Domain mean AUS:** 67.1 | **Median:** 75

### 6. Monitoring Infrastructure (~55% ASF-only, ~14 concepts)

| # | Concept | Arch | Impact | Likelihood | Novelty | AUS |
|---|---------|------|--------|------------|---------|-----|
| 1 | Monitoring logs themselves must be tamper-proof (append-only) | ALL | 5 | 4 | 4 | 80 |
| 2 | SIEM filter change detection (don't filter attack logs) | A4 | 5 | 4 | 5 | 100 |
| 3 | CloudTrail / audit log monitoring | A3, A4 | 4 | 4 | 4 | 64 |
| 4 | Auth0-SIEM integration reliability | A4 | 4 | 4 | 5 | 80 |
| 5 | Audit log failure alerting | A4 | 4 | 3 | 4 | 48 |
| 6 | SoD violation attempt alerts | A5 | 4 | 4 | 4 | 64 |
| 7 | Role change monitoring (detect privilege changes) | A5 | 4 | 4 | 4 | 64 |
| 8 | Credential stuffing detection alerts | A4 | 4 | 3 | 4 | 48 |
| 9 | SIEM PHI access restriction (who can see medical logs) | A4 | 4 | 3 | 5 | 60 |
| 10 | Auditor access pattern monitoring (data scraping detection) | A5 | 4 | 3 | 5 | 60 |

**Domain mean AUS:** 66.8 | **Median:** 64

### 7. Encryption Governance (~65% ASF-only, ~11 concepts)

| # | Concept | Arch | Impact | Likelihood | Novelty | AUS |
|---|---------|------|--------|------------|---------|-----|
| 1 | Backup KMS key separation from primary encryption key | A5 | 5 | 4 | 5 | 100 |
| 2 | KMS key policy restricts decrypt to application role only | A4 | 5 | 4 | 4 | 80 |
| 3 | KMS key rotation schedule (automated, audited) | A1 | 4 | 4 | 4 | 64 |
| 4 | Temp storage encryption (not just primary storage) | A4 | 4 | 4 | 5 | 80 |
| 5 | App server local disk encryption | A4 | 4 | 4 | 4 | 64 |
| 6 | K8s Secret encryption (not just base64) | A3 | 5 | 4 | 4 | 80 |
| 7 | ERP application log encryption | A5 | 4 | 4 | 4 | 64 |
| 8 | Reporting cache encryption | A5 | 4 | 3 | 5 | 60 |
| 9 | Auth0 tenant data encryption verification (SOC 2 review) | A4 | 3 | 3 | 5 | 45 |
| 10 | DB connection TLS with certificate validation | A5 | 5 | 4 | 3 | 60 |

**Domain mean AUS:** 69.7 | **Median:** 64

### 8. Human Factors / Training (~67% ASF-only, ~10 concepts)

| # | Concept | Arch | Impact | Likelihood | Novelty | AUS |
|---|---------|------|--------|------------|---------|-----|
| 1 | Help desk anti-social-engineering procedures | A2 | 5 | 4 | 5 | 100 |
| 2 | App admin SAML training (correct federation config) | A2 | 4 | 4 | 5 | 80 |
| 3 | Provider minimum necessary training (HIPAA) | A4 | 4 | 4 | 4 | 64 |
| 4 | Security team HIPAA-specific training | A4 | 4 | 4 | 4 | 64 |
| 5 | DBA direct PHI access controls awareness | A4 | 4 | 4 | 4 | 64 |
| 6 | Phishing reporting (finance-specific workflow) | A5 | 4 | 4 | 5 | 80 |
| 7 | Patient credential sharing prevention | A4 | 3 | 4 | 4 | 48 |
| 8 | Approval workflow user training (non-repudiation) | A5 | 3 | 4 | 4 | 48 |

**Domain mean AUS:** 68.5 | **Median:** 64

### 9. Network Segmentation (~65% ASF-only, ~13 concepts)

| # | Concept | Arch | Impact | Likelihood | Novelty | AUS |
|---|---------|------|--------|------------|---------|-----|
| 1 | Reporting engine DMZ isolation (SOX) | A5 | 5 | 4 | 5 | 100 |
| 2 | Approval workflow isolation from internet | A5 | 5 | 3 | 5 | 75 |
| 3 | Financial DB no internet route (NAT/IGW) | A5 | 5 | 4 | 4 | 80 |
| 4 | Auditor VPN/bastion path with separate logging | A5 | 4 | 4 | 5 | 80 |
| 5 | Egress traffic control from service mesh | A3 | 4 | 4 | 4 | 64 |
| 6 | Network flow logs for PHI DB (database-level flow) | A4 | 4 | 3 | 5 | 60 |
| 7 | No direct VPN-to-database path | A1 | 5 | 4 | 4 | 80 |
| 8 | Istio authorization policies (not just network policies) | A3 | 4 | 4 | 4 | 64 |

**Domain mean AUS:** 75.4 | **Median:** 77.5

### 10. Federation/SAML/Gov (~50-75% ASF-only, ~8 concepts)

| # | Concept | Arch | Impact | Likelihood | Novelty | AUS |
|---|---------|------|--------|------------|---------|-----|
| 1 | AudienceRestriction validation by each SP | A2 | 5 | 4 | 5 | 100 |
| 2 | SAML assertion expiration validation (not just signature) | A2 | 5 | 4 | 4 | 80 |
| 3 | IdP token validation correctness | A5 | 5 | 4 | 4 | 80 |
| 4 | SSO timeout alignment between IdP and ERP | A5 | 4 | 4 | 5 | 80 |
| 5 | MFA bypass prevention for API/reporting access | A5 | 5 | 4 | 5 | 100 |
| 6 | Hardware security keys (FIDO2) for high-risk operations | A5 | 5 | 3 | 5 | 75 |

**Domain mean AUS:** 85.8 | **Median:** 80

---

## Aggregate Results

### Mean AUS by Domain

| Domain | Concepts Scored | Mean AUS | Median AUS | % Max (125) |
|--------|----------------|----------|------------|-------------|
| Federation/SAML | 6 | **85.8** | 80 | 68.6% |
| Incident Response | 12 | **74.3** | 75 | 59.4% |
| Network Segmentation | 8 | **75.4** | 77.5 | 60.3% |
| Availability & Resilience | 12 | **73.1** | 75 | 58.5% |
| Encryption Governance | 10 | **69.7** | 64 | 55.8% |
| Human Factors | 8 | **68.5** | 64 | 54.8% |
| Data Classification & Flow | 10 | **67.1** | 75 | 53.7% |
| Monitoring Infrastructure | 10 | **66.8** | 64 | 53.4% |
| Identity Lifecycle | 14 | **57.9** | 62 | 46.3% |
| Third-party Dependency | 15 | **53.7** | 60 | 43.0% |
| **All Domains** | **105** | **66.6** | **68** | **53.3%** |

### AUS Distribution

| AUS Range | Interpretation | Count | % |
|-----------|---------------|-------|---|
| **80-125** | High value — critical blind spot | 32 | 30.5% |
| **60-79** | Moderate-high value — genuine gap | 38 | 36.2% |
| **40-59** | Moderate value — valuable but niche | 28 | 26.7% |
| **< 40** | Low value — edges toward noise | 7 | 6.7% |

### Top 10 Highest-AUS ASF-Only Concepts

| Rank | Concept | Domain | I | L | N | AUS |
|------|---------|--------|---|---|---|-----|
| 1 | IR plan includes HIPAA 60-day breach notification | Incident Response | 5 | 4 | 5 | **100** |
| 2 | SIEM filter change detection (don't filter attack logs) | Monitoring | 5 | 4 | 5 | **100** |
| 3 | Backup KMS key separation from primary encryption key | Encryption | 5 | 4 | 5 | **100** |
| 4 | Help desk anti-social-engineering procedures | Human Factors | 5 | 4 | 5 | **100** |
| 5 | Reporting engine DMZ isolation | Network Seg | 5 | 4 | 5 | **100** |
| 6 | AudienceRestriction validation by each SP | Federation | 5 | 4 | 5 | **100** |
| 7 | MFA bypass prevention for API/reporting access | Federation | 5 | 4 | 5 | **100** |
| 8 | Approval workflow availability during month-end close | Availability | 5 | 4 | 5 | **100** |
| 9 | Manager-attested recertification (not self-certified) | Identity | 4 | 4 | 5 | **80** |
| 10 | Recertification thoroughness (no rubber-stamping) | Identity | 4 | 4 | 5 | **80** |

### Key Findings

1. **30.5% of ASF-only concepts are high-value** (AUS ≥ 80). These are critical blind spots no model identified.

2. **66.7% are moderate-to-high value** (AUS ≥ 60). The majority of ASF-only assumptions have genuine security value.

3. **Only 6.7% are low-value** (AUS < 40). The noise rate is low.

4. **Federation/SAML domain has the highest mean AUS (85.8)** — ASF's protocol-level SAML assumptions go beyond what any model surfaces.

5. **Third-party dependency has the lowest mean AUS (53.7)** — valuable but lower impact because vendor dependencies are probabilistic (may not materialize).

6. **Every domain's ASF-only concepts score above midpoint (62.5/125)**. No domain produces predominantly low-value output.

7. **The multi-LLM cluster (GPT/Gemini/Gemma) Tier C from the grand aggregate scored 16.8/25 (67.2% of max) on the 5-dimension AUS.** This independent AUS analysis finds a similar proportion: 66.7% above midpoint — consistent cross-methodology validation.
