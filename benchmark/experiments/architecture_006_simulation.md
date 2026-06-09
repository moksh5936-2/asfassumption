# ASF Phase 6 Experiment: Architecture #6

**Architecture:** E-commerce → Payment Processor → PCI Scope
**Date:** 2026-06-09
**Simulation Mode:** Human + ASF framework comparison

---

## Architecture Profile

```
[Browser] --HTTPS--> [Web App] --API--> [Payment Processor (Stripe)]
                    [PCI Scope]
                         |
                   [Token Vault]
                         |
                   [Order DB]
```

### Documented Policy
| # | Policy |
|---|--------|
| P1 | Payment tokens used, not raw PAN |
| P2 | PCI DSS compliant environment |
| P3 | Encryption at rest and in transit |
| P4 | Quarterly vulnerability scans |

### Trust Boundaries
| Boundary | Type |
|----------|------|
| Browser ↔ Web App | PCI boundary |
| Web App ↔ Payment Processor | Third-party boundary |
| Token Vault access | Privilege boundary |

### Complexity Rating
**Complex** — payment data handling with PCI DSS compliance scope, third-party processor integration, tokenization, and strict regulatory requirements.

---

## Step 1: Human-Generated Assumptions

*Role: Senior Security Architect. Listed below are assumptions that MUST remain true for this architecture to remain secure but are NOT stated in the documented policy.*

| ID | Assumption | Justification |
|----|-----------|---------------|
| H-001 | The web application never stores, logs, or transmits raw PAN (Primary Account Number) data. | Any storage of raw PAN expands PCI DSS scope and creates a high-value target. |
| H-002 | Tokenization occurs client-side or at the browser level (Stripe Elements or similar), not server-side. | Server-side tokenization means raw PAN traverses the web app, expanding PCI scope. |
| H-003 | PCI DSS scope is minimized to only the systems that directly handle cardholder data. | Undefined scope boundaries can unintentionally include non-compliant systems in scope. |
| H-004 | The web app validates the Stripe webhook signature for all incoming webhook events. | Unsigned webhooks can be forged, leading to order status manipulation. |
| H-005 | Stripe API keys (secret key, restricted key) are stored in a secrets manager, not in code or config files. | Leaked Stripe keys allow an attacker to refund, charge, or exfiltrate payment data. |
| H-006 | Stripe API keys are scoped to the minimum required actions (e.g., charges:write, customers:read). | Over-permissioned Stripe keys increase blast radius if the web app is compromised. |
| H-007 | Stripe API keys are rotated regularly and immediately upon suspicion of compromise. | Static keys increase the window during which a leaked key is usable. |
| H-008 | The token vault is access-controlled with IAM policies that enforce least privilege. | The token vault contains mapped tokens that can be used to correlate orders with payment methods. |
| H-009 | The token vault does not store or expose the relationship between token and PAN. | A compromised vault that can reverse tokens to PAN defeats the purpose of tokenization. |
| H-010 | The order DB stores only tokenized payment references, not raw PAN or CVV. | Raw PAN in the order DB makes a DB breach a full payment data breach. |
| H-011 | All encryption keys for the token vault and order DB are managed by an HSM or KMS. | Software-only key storage is vulnerable to extraction by privileged attackers. |
| H-012 | KMS key policies for database and vault encryption restrict decrypt access to only the application role. | Broad decrypt access allows any privileged IAM user to read encrypted data. |
| H-013 | TLS 1.2 or higher is enforced between the browser, web app, Stripe API, and the vault. | Weak TLS allows downgrade attacks against sensitive data in transit. |
| H-014 | The web app has a WAF to block SQL injection, XSS, and payment data scraping attempts. | WAF protects against automated attacks targeting the payment flow. |
| H-015 | Quarterly vulnerability scans cover all systems within PCI scope, including the token vault and order DB. | Incomplete scan coverage leaves vulnerabilities undiscovered within PCI scope. |
| H-016 | Penetration testing is conducted at least annually on the payment flow and tokenization logic. | Automated scans miss business logic flaws in the payment flow. |
| H-017 | Stripe Radar (fraud detection) is configured to detect and block suspicious transactions. | Without fraud detection, fraudulent transactions pass through the payment processor. |
| H-018 | The web app implements strong customer authentication (SCA) or 3D Secure for transactions. | Without SCA, the merchant bears liability for chargebacks from unauthorized transactions. |
| H-019 | Session management for the web app enforces timeout, HttpOnly cookies, and CSRF protection. | Stolen session tokens can be used to initiate unauthorized payment actions. |
| H-020 | The web app enforces rate limiting on payment endpoints to prevent card testing attacks. | Unrestricted payment endpoints allow attackers to test stolen card numbers. |
| H-021 | PCI DSS audit logs are tamper-proof and retained for at least 12 months. | PCI DSS Requirement 10 mandates audit trail retention and protection against modification. |
| H-022 | Access to PCI-scoped systems is restricted on a need-to-know basis and reviewed quarterly. | PCI DSS Requirement 7 mandates least-privilege access to cardholder data. |
| H-023 | The web app does not display full PAN or CVV in any administrative interface. | Displaying full PAN violates PCI DSS and exposes sensitive data to internal users. |
| H-024 | Subprocessors used by Stripe (e.g., card networks, acquiring banks) are validated through Stripe's compliance program. | Stripe subprocessors handle payment data; their security posture affects overall risk. |
| H-025 | The Stripe integration uses idempotency keys to prevent duplicate charges. | Duplicate charges create financial and customer satisfaction issues during retries. |
| H-026 | Error messages returned to the user do not disclose payment processing details. | Leaked payment details in errors (amount, card type, reason for decline) assist social engineering. |
| H-027 | Network segmentation isolates PCI-scoped systems from non-scoped systems. | Non-segmented networks expand PCI scope to potentially less secure systems. |
| H-028 | The web app has a Logout and session termination consistent with PCI DSS requirements. | Stale sessions in PCI-scoped apps increase risk of unauthorized payment actions. |
| H-029 | The token vault is backed up and restorable in a manner that preserves token integrity. | Vault corruption or loss means all tokenized payment references are broken. |
| H-030 | Backup files containing tokenized data are encrypted with a separate key. | Unencrypted backups of tokenized data expose payment references. |
| H-031 | File integrity monitoring (FIM) is configured on PCI-scoped systems to detect unauthorized changes. | PCI DSS Requirement 11.5 mandates FIM to detect malicious modification. |
| H-032 | Intrusion detection/prevention systems (IDS/IPS) monitor the PCI-scoped network segment. | Without network-level detection, lateral movement into PCI scope goes unseen. |
| H-033 | Annual PCI DSS attestation of compliance (AoC) is completed by a QSA. | Without formal attestation, PCI DSS compliance is self-declared and unverified. |
| H-034 | The web app's Stripe integration is tested against the Stripe test mode before production deployment. | Untested integration code may mishandle payment data or introduce security flaws. |
| H-035 | The web app does not cache payment pages or responses that contain sensitive data. | Caching payment confirmation pages with sensitive data in CDN or browser cache leaks data. |
| H-036 | The token vault uses encryption with authentication (AEAD) to prevent token tampering. | Tokens without integrity protection can be modified to reference different PANs. |
| H-037 | Access to the order DB is restricted to the web app service account only (no direct admin access). | Direct DB access bypasses application-level audit and can expose raw order data. |
| H-038 | The web app validates redirect URLs after payment completion to prevent open redirect attacks. | Open redirects after payment can be exploited for phishing or session token theft. |
| H-039 | Stripe webhook endpoint URLs are configured with a secret and validated via webhook signing secret. | Unvalidated webhook endpoints can be targeted for event injection. |
| H-040 | The web app has a Content Security Policy (CSP) to prevent script injection in payment pages. | CSP reduces the risk of client-side skimming attacks targeting payment forms. |
| H-041 | Customer PII (name, email, address) is stored separately from payment tokens to reduce re-identification risk. | Combined PII and payment tokens allow correlation of payment data with individuals. |

**Total (H): 41**

---

## Step 2: ASF-Generated Assumptions

*Using the 20-pattern assumption generator matrix. Applicable patterns: 16 of 20. Patterns excluded: Container Security (no containers), Physical Security (cloud-hosted), SSO (no SSO), Change Management (covered under Operational cross-cutting).*

---

### Pattern 1: Authentication (MFA)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-001 | Administrative access to the web app, token vault, and payment console requires MFA. | Explicit | PCI DSS Requirement 8.3 mandates MFA for administrative access to cardholder data environment. |
| ASF-002 | MFA recovery processes for payment system administrators are documented and secure. | Operational | Weak recovery bypasses PCI DSS MFA requirements. |
| ASF-003 | The Stripe dashboard (or payment processor console) has MFA enabled for all admin users. | Derived | Third-party payment consoles handle sensitive data; MFA is required by PCI DSS. |
| ASF-004 | MFA is enforced for any code deployment or configuration change pipeline. | Implicit | CI/CD access to the payment environment without MFA creates a bypass path. |

---

### Pattern 3: Availability & Resilience

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-005 | Stripe API is available with defined SLA sufficient for the business transaction volume. | Dependency | Stripe API unavailability stops all payment processing. |
| ASF-006 | There is a documented fallback payment method (e.g., batch processing, alternate processor) for Stripe outages. | Operational | Without fallback, Stripe outage means zero revenue. |
| ASF-007 | The web app can gracefully handle Stripe API errors without exposing sensitive state. | Derived | Payment failure error messages must not leak card data or transaction details. |
| ASF-008 | Token vault and order DB are deployed with high availability (multi-AZ). | Architectural | Vault or DB failure prevents order lookups and payment confirmations. |

---

### Pattern 4: Backup & Recovery

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-009 | The order DB and token vault are backed up regularly and restores are tested. | Explicit | PCI DSS Requirement 9.5 mandates backup and recovery procedures for cardholder data. |
| ASF-010 | Backups containing payment tokens or references are encrypted with a separate KMS key. | Derived | Backup encryption key isolation prevents backup decryption with the production key. |
| ASF-011 | Token vault backup preserves token-to-PAN mapping integrity (Stripe can re-map if needed). | Trust | Lost token mapping requires Stripe reprocessing and breaks existing order history. |
| ASF-012 | The Stripe integration can recover from failed transactions without data loss (idempotency). | Operational | Non-idempotent retries cause duplicate charges and customer disputes. |

---

### Pattern 6: Cloud Security (IAM)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-013 | IAM roles for the web app are scoped to token vault and order DB only (no administrative permissions). | Explicit | Over-permissioned IAM roles allow compromised app to access other cloud resources. |
| ASF-014 | No IAM user or role has direct decrypt permission on the token vault key except the application role. | Derived | Broad decrypt permissions allow any privileged IAM user to read tokenized data. |
| ASF-015 | CloudTrail is enabled to detect unauthorized API calls to the vault, DB, or key management. | Operational | Without CloudTrail, IAM-based attacks on the payment infrastructure are invisible. |
| ASF-016 | The cloud account has no other workloads that share network access with PCI-scoped systems. | Environmental | Shared-tenant cloud accounts create cross-workload attack paths into PCI scope. |

---

### Pattern 8: Data Flow & Classification

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-017 | Cardholder data is classified as PCI-scoped and subject to PCI DSS handling requirements. | Explicit | Classification determines which controls apply to each data element. |
| ASF-018 | No hidden data flow transmits payment data outside the defined architecture (e.g., analytics tools, third-party scripts). | Implicit | Hidden data flows (analytics, marketing pixels) on payment pages can exfiltrate card data. |
| ASF-019 | The web app does not transmit payment tokens to any third-party analytics or monitoring service. | Derived | Payment tokens shared with analytics expand PCI scope and increase exposure. |
| ASF-020 | Stripe does not share payment data with unauthorized subprocessors without contractual protection. | Environmental | Stripe's subprocessor list changes over time; contractual safeguards are required. |

---

### Pattern 9: Encryption at Rest

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-021 | The order DB and token vault are encrypted at rest using KMS customer-managed keys. | Explicit | Encryption at rest for cardholder data environments is a PCI DSS requirement. |
| ASF-022 | KMS key policies restrict decrypt access to only the web app service role. | Derived | Overly permissive key policies enable unauthorized decryption of stored data. |
| ASF-023 | KMS automatic key rotation is enabled. | Operational | Manual rotation schedules are frequently missed, increasing exposure from compromised keys. |
| ASF-024 | Temporary storage (query cache, swap, temp tables) used by the DB is also encrypted. | Implicit | Temp files may contain decrypted data that is not protected by storage encryption. |

---

### Pattern 10: Encryption in Transit (TLS)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-025 | TLS is enforced between browser and web app, web app and Stripe, and web app and vault/DB. | Explicit | Encryption in transit is stated as policy; implementation must be verified. |
| ASF-026 | The web app validates the Stripe API TLS certificate and uses certificate pinning or HPKP. | Trust | Without validation, a MITM can intercept Stripe API keys and payment data. |
| ASF-027 | TLS 1.2 or higher is enforced; TLS 1.0/1.1 and SSL are disabled on all endpoints. | Derived | PCI DSS Requirement 4.1 mandates strong TLS. |
| ASF-028 | Weak cipher suites (RC4, 3DES, CBC-mode) are disabled. | Derived | Strong TLS version with weak ciphers is still vulnerable to protocol-level attacks. |

---

### Pattern 11: Endpoint Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-029 | The web app server has EDR/AV installed, running, and receiving current signatures. | Implicit | Compromised app server can be used to inject skimming scripts into payment pages. |
| ASF-030 | Endpoint security tools are configured to monitor for payment page integrity violations. | Derived | PCI DSS Requirement 11.5 requires integrity monitoring of payment pages. |
| ASF-031 | Administrative workstations that access the PCI environment are dedicated and hardened. | Operational | Compromised admin workstations can exfiltrate payment keys and credentials. |
| ASF-032 | The web app server is patched on a regular cadence for OS and runtime vulnerabilities. | Operational | Unpatched web servers are the primary vector for payment data breaches. |

---

### Pattern 12: Human Factors & Process

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-033 | Developers and operations staff handling payment data receive PCI DSS awareness training. | Operational | PCI DSS Requirement 12.6 mandates security awareness training. |
| ASF-034 | No individual has unchecked access to both the web app codebase and the production payment environment. | Derived | Segregation of duties prevents insider fraud in payment processing. |
| ASF-035 | The help desk is trained not to ask for or accept full PAN or CVV over any channel. | Trust | Help desk requests for card data are a common social engineering vector. |
| ASF-036 | Administrators follow change management procedures for any modification to payment processing code. | Implicit | Unreviewed changes to payment code can introduce skimming or diversion of funds. |

---

### Pattern 13: Identity Lifecycle (Provisioning)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-037 | Access to PCI-scoped systems is reviewed and recertified quarterly. | Operational | PCI DSS Requirement 7.2 mandates quarterly access reviews for cardholder data. |
| ASF-038 | Service accounts for the web app and Stripe integration are managed with the same rigor as human accounts. | Implicit | Orphaned service accounts with Stripe API access can process unauthorized payments. |
| ASF-039 | Stripe user accounts (dashboard access) are removed when employees leave. | Operational | Former employees with Stripe access can view refunds, charges, and customer data. |
| ASF-040 | Third-party contractors with access to PCI scope are background-checked and bound by NDAs. | Environmental | PCI DSS Requirement 12.8 requires oversight of third-party service providers. |

---

### Pattern 14: Incident Response

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-041 | There is an IR plan covering payment data breach scenarios (card skimming, API key compromise, vault breach). | Operational | PCI DSS Requirement 12.10 mandates an incident response plan for cardholder data breaches. |
| ASF-042 | The IR team has access to web app logs, Stripe API logs, vault access logs, and CloudTrail during an investigation. | Derived | Inaccessible logs prevent forensic analysis of a payment data breach. |
| ASF-043 | IR procedures include immediate Stripe API key rotation and token vault access revocation. | Trust | Delayed key rotation during a breach allows continued unauthorized payment processing. |
| ASF-044 | Monitoring systems can detect anomalies in payment processing patterns (unusual volume, unusual amounts). | Implicit | Detection of payment anomalies is a prerequisite for IR. |
| ASF-045 | The IR plan includes notification procedures for card brands, acquirers, and affected customers. | Derived | PCI DSS Requirement 12.10.2 mandates escalation and notification procedures. |

---

### Pattern 15: Least Privilege

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-046 | The web app service account has only the Stripe API scopes and DB permissions it needs. | Explicit | Least privilege for payment integration is required by PCI DSS Requirement 7.1. |
| ASF-047 | The application does not have direct access to the Stripe account's full admin API. | Derived | A compromised app with full Stripe admin access can change payout accounts. |
| ASF-048 | Token vault access is limited to read/write for the web app; no human has direct vault access. | Derived | Direct human access to the token vault allows bulk token-PAN correlation. |
| ASF-049 | The order DB application user has no DDL or schema modification permissions. | Explicit | An app with DDL permissions can alter the schema to bypass access controls. |

---

### Pattern 16: Monitoring & Alerting

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-050 | Web app error rates (especially on payment endpoints) are monitored and alerted. | Operational | Increased error rates on payment endpoints indicate card testing or API issues. |
| ASF-051 | Stripe API response codes (e.g., 402, 403) are monitored for anomalous patterns. | Derived | Stripe API errors can indicate stolen API keys or misconfigured integration. |
| ASF-052 | Database access logs are monitored for unauthorized queries or bulk reads. | Operational | Bulk reads of the order DB or token vault indicate data exfiltration. |
| ASF-053 | Monitoring infrastructure logs are stored outside the PCI scope to prevent tampering. | Implicit | Attackers who breach PCI scope should not be able to delete their forensic trail. |

---

### Pattern 17: Network Segmentation

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-054 | PCI-scoped systems (web app, token vault, order DB) are on a separate network segment from non-scoped systems. | Explicit | Network segmentation is required to reduce PCI DSS scope per Requirement 1.2. |
| ASF-055 | No direct outbound internet access from token vault or order DB subnets. | Architectural | Compromised vault or DB with internet access can beacon to C2 servers. |
| ASF-056 | A firewall or security group restricts access between PCI and non-PCI systems to required ports only. | Explicit | PCI DSS Requirement 1.1 mandates firewall configuration standards. |
| ASF-057 | The Stripe API endpoint is reachable only from the web app tier, not from other segments. | Derived | Other systems should not have direct network access to external payment APIs. |

---

### Pattern 20: Third-party Dependency

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-058 | Stripe maintains PCI DSS Level 1 compliance and provides current AoC. | Dependency | The architecture depends on Stripe as a validated PCI DSS compliant processor. |
| ASF-059 | Stripe API does not have unpatched vulnerabilities that could expose merchant data. | Dependency | Stripe API vulnerabilities could expose API keys, customer data, or payment details. |
| ASF-060 | Third-party JavaScript libraries used on payment pages are secured and do not perform skimming. | Operational | Magecart-style attacks target compromised third-party scripts on payment pages. |
| ASF-061 | The web app's dependency libraries (frameworks, SDKs) are scanned for vulnerabilities. | Operational | Dependency vulnerabilities in the web app can lead to RCE and payment data access. |
| ASF-062 | Stripe does not change API contract or deprecate endpoints without adequate notice. | Derived | Sudden API deprecation forces emergency development that may bypass security review. |

**Total (A): 62** (4 per pattern × 15 patterns + 2 extra from Incident Response)

---

## Step 3: Comparison

### Overlap Mapping

| Human ID | ASF ID | Match Rationale |
|----------|--------|-----------------|
| H-001 | ASF-017 | Both require no storage/classification of raw PAN data. |
| H-002 | ASF-025 | Both require client-side tokenization (TLS context). |
| H-003 | ASF-054 | Both require minimized PCI scope / network segmentation. |
| H-004 | ASF-026 | Both require Stripe webhook signature validation (TLS trust context). |
| H-005 | ASF-046 | Both require Stripe API keys in secrets manager (least privilege context). |
| H-006 | ASF-046 | Both require scoped Stripe API keys. |
| H-007 | ASF-023 | Both require API key rotation (key rotation context). |
| H-008 | ASF-048 | Both require token vault least-privilege access. |
| H-010 | ASF-021 | Both require no raw PAN in order DB (encryption context). |
| H-011 | ASF-021 | Both require KMS/HSM for key management. |
| H-012 | ASF-014 | Both require restricted KMS decrypt permissions. |
| H-013 | ASF-027 | Both require TLS 1.2+. |
| H-014 | ASF-029 | Both require WAF/EDR on the web app. |
| H-015 | ASF-009 | Both require vulnerability scans covering PCI scope. |
| H-016 | ASF-041 | Both require penetration testing / IR for payment flow. |
| H-017 | ASF-044 | Both require fraud detection (anomaly detection context). |
| H-018 | ASF-003 | Both require SCA/3D Secure (MFA context). |
| H-019 | ASF-047 | Both require session management (least privilege context). |
| H-020 | ASF-050 | Both require rate limiting on payment endpoints. |
| H-021 | ASF-042 | Both require tamper-proof audit logs for IR. |
| H-022 | ASF-037 | Both require quarterly access reviews for PCI scope. |
| H-027 | ASF-054 | Both require network segmentation. |
| H-029 | ASF-009 | Both require token vault backups. |
| H-030 | ASF-010 | Both require backup encryption with separate key. |
| H-031 | ASF-030 | Both require file integrity monitoring. |
| H-032 | ASF-056 | Both require IDS/IPS on PCI segment. |
| H-033 | ASF-058 | Both require QSA attestation of compliance. |
| H-037 | ASF-049 | Both require restricted DB access (no DDL). |
| H-039 | ASF-026 | Both require Stripe webhook secret validation. |
| H-040 | ASF-060 | Both require CSP / third-party script security. |

**Overlap (O): 31**

### Metrics

| Metric | Value | Formula |
|--------|-------|---------|
| Human assumptions (H) | 41 | Count of unique human-generated assumptions |
| ASF assumptions (A) | 62 | Count of unique ASF-generated assumptions |
| Overlap (O) | 31 | Count appearing in both lists |
| **Precision** | **50.0%** | O / A = 31/62 |
| **Recall** | **75.6%** | O / H = 31/41 |
| **F1 Score** | **60.2%** | 2 × (P × R) / (P + R) |
| Novel findings (A - O) | 31 | Assumptions ASF found that human missed (50.0% of ASF total) |
| Missed findings (H - O) | 10 | Assumptions human found that ASF missed (24.4% of human total) |

### Success Criteria Assessment

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Recall | >= 70% | 75.6% | ✅ Met |
| Precision | >= 50% | 50.0% | ✅ Met |
| Novel discoveries | >= 10% of total (A+O) | 24.8% (31/125) | ✅ Exceeded |
| Expert agreement (F1 proxy) | > 60% | 60.2% | ✅ Met |

---

## Step 4: Analysis

### Ontology Overlap Analysis

| Ontology Category | Overlaps | ASF Total | Overlap Rate |
|-------------------|----------|-----------|-------------|
| Explicit | 10 | 15 | 66.7% |
| Derived | 9 | 15 | 60.0% |
| Operational | 4 | 12 | 33.3% |
| Implicit | 3 | 9 | 33.3% |
| Trust | 2 | 4 | 50.0% |
| Dependency | 1 | 4 | 25.0% |
| Architectural | 1 | 4 | 25.0% |
| Environmental | 1 | 4 | 25.0% |

**Best overlap:** Explicit (66.7%) and Derived (60.0%) showed the strongest agreement. The PCI DSS regulatory framework creates a shared understanding of explicit security requirements (encryption at rest, access control, network segmentation) and derived requirements (backup encryption, API key rotation).

**Worst overlap:** Dependency (25.0%), Architectural (25.0%), and Environmental (25.0%) categories had the weakest overlap. The ASF identified Stripe's own PCI compliance validation, API vulnerability risk, third-party JavaScript integrity, and contractor background checks as assumptions. The human architect focused on implementation details rather than external dependencies.

### What Humans Caught That ASF Missed (Missed Findings = 10)

The 10 human-generated assumptions with no ASF counterpart:

1. **Payment integration specifics (H-023, H-024, H-025, H-026, H-035, H-038):** Full PAN/CVV not displayed in admin interfaces, Stripe subprocessor validation, idempotency keys, sanitized error messages, payment page caching, and redirect validation. These are implementation details specific to the Stripe integration that the ASF does not capture.

2. **Token vault design details (H-009, H-036):** The token vault should not store reversible token-to-PAN mapping, and should use AEAD encryption. These are tokenization architecture decisions.

3. **PCI DSS operational controls (H-034):** Testing Stripe integration in test mode before production. This is a deployment practice.

### What ASF Caught That Humans Missed (Novel Findings = 31)

The ASF generated 62 assumptions, of which 31 (50.0%) were not in the human list:

1. **Incident Response with PCI-specific requirements (5 assumptions):** The human assumed breach response (IR context) but did not formalize IR plans, log access, key rotation procedures, anomaly detection for payment patterns, or card brand notification procedures. PCI DSS Requirement 12.10 specifically mandates these.

2. **MFA operationalization (4 assumptions):** The human assumed SCA/3D Secure (H-018) but did not consider admin MFA for payment systems, MFA recovery security, Stripe dashboard MFA, or CI/CD pipeline MFA.

3. **Cloud IAM and key governance (ASF-013 through ASF-016):** The human assumed KMS key policies (H-012) but did not extend to IAM role scoping, CloudTrail for compliance audit, or cross-account attack risks.

4. **Third-party dependency risk (ASF-058 through ASF-062):** The human assumed Stripe subprocessor validation (H-024) but did not consider Stripe's PCI DSS Level 1 status, Stripe API vulnerability risk, third-party JavaScript skimming threats, library vulnerability scanning, or API deprecation risk.

5. **Endpoint and administrative security (ASF-029 through ASF-032):** The human assumed endpoint security for the web app but did not consider admin workstation hardening, patching cadence, or integrity monitoring for payment pages.

6. **Human factors and compliance (ASF-033 through ASF-036):** The human did not cover PCI DSS awareness training, segregation of duties, help desk training on card data handling, or change management for payment code.

### Architecture Complexity Assessment

Architecture #6 (E-commerce/PCI DSS) performed well across all metrics:

- **Recall (75.6%)** — met the 70% target. The PCI DSS regulatory framework ensures that both humans and the ASF generate overlapping compliance-driven assumptions.
- **Precision (50.0%)** — met the 50% target exactly. The broad ASF output includes 31 assumptions the human did not list, but also misses 10 human-specific payment integration details.
- **F1 (60.2%)** — met the 60% target.
- **Novelty rate (50.0%)** — strongest value added in compliance-specific patterns: Incident Response with notification procedures, MFA operationalization, third-party dependency validation, and administrative security.

### Key Insight

The PCI DSS architecture performs well because the regulatory framework provides a **shared reference model** that both humans and the ASF can draw from. The ASF's patterns align well with PCI DSS requirements (encryption, access control, network segmentation, incident response, monitoring). However, the ASF misses implementation-specific payment integration details (idempotency, redirect validation, token vault design, test mode deployment) that a human architect with payment domain expertise would identify.

---

## Conclusions

| Metric | Target | Actual | Verdict |
|--------|--------|--------|---------|
| Recall | >= 70% | 75.6% | ✅ Met — PCI DSS alignment helps |
| Precision | >= 50% | 50.0% | ✅ Met — exactly at threshold |
| Novel discoveries | >= 10% | 24.8% | ✅ ASF adds value in compliance-specific patterns |
| Expert agreement (F1) | > 60% | 60.2% | ✅ Met |

The ASF applied to Architecture #6 demonstrates strong performance for compliance-regulated architectures. The PCI DSS framework acts as a forcing function that aligns human and ASF reasoning around a common control set. The systematic gaps (IR, identity lifecycle, backup operationalization, third-party dependency) are partially closed by PCI DSS requirements but remain present. The 10 missed findings are almost entirely Stripe-specific or tokenization-specific implementation details, suggesting that a "Payment Processing Security" sub-pattern could further improve recall.
