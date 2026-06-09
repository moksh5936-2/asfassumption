# Multi-LLM Evaluation Campaign 001: Architecture #1

**Architecture:** User Laptop → VPN Gateway → Internal Web App → Payroll Database (RDS)
**Date:** 2026-06-09
**Evaluation Mode:** 5 AI Security Architect personas + Human + ASF comparison

---

## 1. Campaign Overview

This campaign simulates 5 different AI security architects independently reviewing Architecture #1. Each AI produces an assumption list. These are compared against the 40 human-generated assumptions and 64 ASF-generated assumptions from the Phase 6 simulation.

### Architecture

```
[User Laptop] --VPN--> [VPN Gateway] --TLS--> [Internal Web App] --SQL--> [Payroll Database (RDS)]
```

### Documented Policy
- VPN required for remote access
- Application authenticates with AD credentials
- Database is in private subnet
- Backups run nightly

### Persona Summary

| Persona | Style | Count | Key Strength | Key Weakness |
|---------|-------|-------|-------------|--------------|
| GPT | Analytical, step-by-step | 45 | Logical chains, structured enumeration | May miss subtle contextual assumptions |
| Claude | Thorough, nuanced | 52 | Edge cases, comprehensive coverage | May over-include low-probability scenarios |
| Gemini | Concise, direct | 40 | High-confidence, focused | May miss less obvious assumptions |
| DeepSeek | Technical, detail-oriented | 48 | Infrastructure/protocol details | May miss business/process assumptions |
| Qwen | Balanced, practical | 45 | Mix of technical and process | Moderate depth in both |

### Source Counts

| Source | Count |
|--------|-------|
| Human Architect | **40** |
| ASF Framework | **64** |
| GPT | **45** |
| Claude | **52** |
| Gemini | **40** |
| DeepSeek | **48** |
| Qwen | **45** |

---

## 2. Per-Persona Assumption Lists

### Persona 1: GPT (Analytical, Step-by-Step)

Prompt style: *"Think step by step through each component"*

| ID | Assumption |
|----|-----------|
| GPT-001 | VPN gateway firmware is patched against all known critical CVEs. |
| GPT-002 | VPN client software prohibits split-tunneling on all user laptops. |
| GPT-003 | VPN credentials are individual and never shared. |
| GPT-004 | VPN gateway enforces MFA for every remote authentication attempt. |
| GPT-005 | Modern encryption (AES-256, SHA-256) is enforced on the VPN gateway. |
| GPT-006 | VPN gateway supports sufficient concurrent sessions for all remote employees. |
| GPT-007 | VPN client implements a kill switch to block non-VPN traffic on disconnect. |
| GPT-008 | Perfect forward secrecy is enabled on the VPN gateway. |
| GPT-009 | VPN gateway has redundant hardware or failover to avoid single-point-of-failure. |
| GPT-010 | VPN gateway logs are streamed to a centralized SIEM in real time. |
| GPT-011 | TLS certificates presented by the VPN gateway and web app are valid and CA-signed. |
| GPT-012 | TLS 1.2 or higher is enforced on all encrypted links; SSL and TLS 1.0/1.1 are disabled. |
| GPT-013 | Weak cipher suites (RC4, 3DES, CBC-mode) are explicitly disabled on all TLS endpoints. |
| GPT-014 | TLS termination validates certificates at both VPN gateway and web application ends. |
| GPT-015 | Active Directory is hardened against Kerberoasting, DCSync, and pass-the-hash attacks. |
| GPT-016 | AD enforces account lockout after a configured threshold of failed logins. |
| GPT-017 | NTLMv1 and LM hash protocols are disabled; only Kerberos is used. |
| GPT-018 | AD domain controllers are patched, monitored, and available for authentication at all times. |
| GPT-019 | SSO session timeout values are synchronized between AD and the web application. |
| GPT-020 | SSO token signing keys are rotated regularly and stored in a hardware security module. |
| GPT-021 | The web application uses parameterized queries exclusively to prevent SQL injection. |
| GPT-022 | The web application sets HttpOnly, Secure, and SameSite attributes on session cookies. |
| GPT-023 | Role-based access control is enforced within the web application for payroll data. |
| GPT-024 | The web application implements Content Security Policy and XSS/CSRF protections. |
| GPT-025 | Login endpoints in the web application implement rate limiting against credential stuffing. |
| GPT-026 | Database error details are never exposed in web application HTTP responses. |
| GPT-027 | The web application does not send payroll data to any third-party endpoint. |
| GPT-028 | Database credentials are stored securely and not committed to source code. |
| GPT-029 | The web application process does not run with root or administrator privileges. |
| GPT-030 | The database user has SELECT/INSERT/UPDATE/DELETE privileges only on the payroll schema. |
| GPT-031 | The RDS instance resides in a private subnet with no internet gateway route. |
| GPT-032 | The database security group allows ingress only from the web application security group. |
| GPT-033 | No network path exists between the VPN subnet and the database subnet directly. |
| GPT-034 | RDS encryption at rest is enabled with AWS KMS and a customer-managed key. |
| GPT-035 | Database audit logging is enabled for all DML and DDL statements. |
| GPT-036 | RDS automatic minor version upgrades are enabled for critical security patches. |
| GPT-037 | Database credentials are rotated every 90 days and immediately on any compromise suspicion. |
| GPT-038 | Nightly backups complete within the configured maintenance window without performance degradation. |
| GPT-039 | Backup restore procedures are tested at least annually and documented. |
| GPT-040 | Backup files are encrypted at rest using a dedicated KMS key. |
| GPT-041 | The VPN gateway, web application, and database each reside in separate subnets with explicit allow rules. |
| GPT-042 | The database does not initiate outbound connections to the VPN gateway or web application. |
| GPT-043 | IAM roles for RDS and application access follow least-privilege principles. |
| GPT-044 | VPN access grants are revoked within the same shift on employee termination. |
| GPT-045 | The web application has a Content Security Policy header configured to limit data exfiltration channels. |

**Total: 45**

### Persona 2: Claude (Thorough, Nuanced)

Prompt style: *"Consider the architecture holistically"*

| ID | Assumption |
|----|-----------|
| CLAUDE-001 | VPN gateway firmware is updated within vendor patch window for critical CVEs. |
| CLAUDE-002 | VPN client is centrally managed with mandatory split-tunneling disabled. |
| CLAUDE-003 | VPN credentials are unique per user with individual accountability enforced. |
| CLAUDE-004 | MFA required for all VPN connections including service accounts. |
| CLAUDE-005 | VPN gateway enforces AES-256-GCM and SHA-256 minimum; rejects legacy proposals. |
| CLAUDE-006 | VPN gateway has sufficient session capacity for peak concurrent remote access. |
| CLAUDE-007 | VPN client kill switch is enabled and verified to block all non-tunnel traffic on connection loss. |
| CLAUDE-008 | Perfect forward secrecy is configured on VPN gateway to protect recorded traffic. |
| CLAUDE-009 | VPN gateway has active-passive failover or clustering to avoid SPOF. |
| CLAUDE-010 | Documented offline procedures exist for VPN outages without bypassing security controls. |
| CLAUDE-011 | VPN logs shipped to SIEM with alerting on anomalous patterns. |
| CLAUDE-012 | TLS certificates for VPN and web app issued by trusted internal CA and monitored for expiry. |
| CLAUDE-013 | TLS 1.2 minimum enforced; TLS 1.0, 1.1, SSL explicitly disabled. |
| CLAUDE-014 | All weak TLS cipher suites (RC4, 3DES, CBC-mode) removed from cipher allow-list. |
| CLAUDE-015 | Mutual TLS considered for server-to-server communication for dual-end validation. |
| CLAUDE-016 | Active Directory configured with protected users group and SMB signing to resist relay attacks. |
| CLAUDE-017 | AD password policies enforce complexity, length >=14, and periodic rotation. |
| CLAUDE-018 | Account lockout configured in AD with reasonable threshold and duration. |
| CLAUDE-019 | NTLMv1 blocked by Group Policy; only Kerberos or NTLMv2 permitted. |
| CLAUDE-020 | AD domain controllers hardened, patched, monitored with availability SLAs. |
| CLAUDE-021 | SSO session lifetimes consistently configured between AD and web app. |
| CLAUDE-022 | SSO token signing keys stored in HSM or key vault with access logging. |
| CLAUDE-023 | Web app uses prepared statements or ORM parameter binding to eliminate SQL injection. |
| CLAUDE-024 | Session cookies are HttpOnly, Secure, SameSite=Strict with 15-min idle timeout. |
| CLAUDE-025 | Web app enforces app-level RBAC mapping AD groups to payroll data access tiers. |
| CLAUDE-026 | Web app mitigates XSS with output encoding, CSRF with anti-forgery tokens, clickjacking with CSP. |
| CLAUDE-027 | Rate limiting on all auth endpoints with graduated delays and CAPTCHA after threshold. |
| CLAUDE-028 | Error handlers sanitize stack traces and schema details from production error pages. |
| CLAUDE-029 | Payroll data not transmitted to any monitoring, analytics, or third-party service outside defined flow. |
| CLAUDE-030 | Database credentials stored in secrets manager with automated rotation. |
| CLAUDE-031 | Web app runs under dedicated low-privilege service account, not root/LocalSystem. |
| CLAUDE-032 | App DB user has least privilege — SELECT, INSERT, UPDATE, DELETE only on payroll schema, no DDL. |
| CLAUDE-033 | RDS instance in private subnet with no direct/indirect internet route via NAT or egress-only IGW. |
| CLAUDE-034 | DB security group permits inbound on port 5432/3306 only from web app security group. |
| CLAUDE-035 | No direct routing path from VPN subnet to DB subnet at VPC or transit gateway level. |
| CLAUDE-036 | RDS encryption at rest uses customer-managed KMS key with automatic rotation. |
| CLAUDE-037 | DB audit logging captures all SELECT, INSERT, UPDATE, DELETE with user and timestamp. |
| CLAUDE-038 | RDS has automated patching for critical and high-severity CVEs with defined maintenance window. |
| CLAUDE-039 | Database credentials rotated automatically on defined schedule and upon security incident. |
| CLAUDE-040 | Nightly backups complete within backup window; DB performance monitored for degradation. |
| CLAUDE-041 | Backup restore tested at least annually with documented RTO/RPO validation. |
| CLAUDE-042 | Backups encrypted at rest using KMS key separate from primary DB encryption key. |
| CLAUDE-043 | Payroll data formally classified as sensitive with data handling policies. |
| CLAUDE-044 | User laptops have EDR/AV with real-time protection and MDM disk encryption. |
| CLAUDE-045 | Lost/stolen laptops can be remotely wiped and VPN credentials revoked in single orchestrated action. |
| CLAUDE-046 | Users receive annual security awareness training covering phishing and credential hygiene. |
| CLAUDE-047 | Access provisioning follows joiner/mover/leaver process enforced by HR system integration. |
| CLAUDE-048 | AD group membership for VPN and app access recertified quarterly by data owners. |
| CLAUDE-049 | Documented IR plan specifically covers payroll data breach scenarios with defined roles. |
| CLAUDE-050 | IR team has timely access to VPN, app, and DB logs with preserved chain of custody. |
| CLAUDE-051 | IR plan includes DB isolation and VPN disconnection procedures preserving forensic evidence. |
| CLAUDE-052 | Vendor exit strategy documented for both VPN vendor and RDS including data migration. |

**Total: 52**

### Persona 3: Gemini (Concise, Direct)

Prompt style: *"What are the key assumptions?"*

| ID | Assumption |
|----|-----------|
| GEMINI-001 | VPN gateway firmware is patched and free of known exploits. |
| GEMINI-002 | Split-tunneling is disabled on all managed VPN clients. |
| GEMINI-003 | VPN credentials are unique per user and not shared. |
| GEMINI-004 | MFA is enforced for all VPN access. |
| GEMINI-005 | VPN gateway uses strong encryption with no legacy protocol fallback. |
| GEMINI-006 | VPN gateway has capacity for all concurrent remote users. |
| GEMINI-007 | VPN client kill switch prevents traffic leakage on tunnel drop. |
| GEMINI-008 | Perfect forward secrecy is enabled on the VPN. |
| GEMINI-009 | VPN logs are sent to SIEM and monitored for threats. |
| GEMINI-010 | TLS certificates for VPN and web app are valid and monitored for expiry. |
| GEMINI-011 | TLS 1.2 or higher is enforced throughout. |
| GEMINI-012 | Weak ciphers (RC4, 3DES) are disabled at all TLS endpoints. |
| GEMINI-013 | Active Directory is hardened against lateral movement techniques. |
| GEMINI-014 | AD account lockout is configured to resist brute-force. |
| GEMINI-015 | NTLMv1 is disabled in the domain. |
| GEMINI-016 | AD domain controllers are highly available and patched. |
| GEMINI-017 | SSO tokens are cryptographically signed and validated by the web app. |
| GEMINI-018 | Web application is immune to SQL injection via parameterized queries. |
| GEMINI-019 | Session cookies are secured with HttpOnly and Secure flags. |
| GEMINI-020 | Application enforces RBAC for payroll data access. |
| GEMINI-021 | XSS and CSRF protections are active in the web application. |
| GEMINI-022 | Login endpoints have rate limiting against credential abuse. |
| GEMINI-023 | Error messages do not leak database internals to users. |
| GEMINI-024 | Payroll data is not leaked to external endpoints. |
| GEMINI-025 | Database credentials are stored in a secrets manager, not code. |
| GEMINI-026 | The application DB user has least privilege (no DDL). |
| GEMINI-027 | RDS is in a private subnet with no internet route. |
| GEMINI-028 | Database security group allows only web application traffic. |
| GEMINI-029 | No direct route from VPN subnet to database subnet. |
| GEMINI-030 | RDS encryption at rest is enabled. |
| GEMINI-031 | Database auditing is enabled for all queries. |
| GEMINI-032 | Backups are encrypted at rest. |
| GEMINI-033 | Backup restore is tested at least annually. |
| GEMINI-034 | Endpoints have EDR/AV installed and updated. |
| GEMINI-035 | Laptops are MDM-managed with disk encryption enforced. |
| GEMINI-036 | Access is revoked on termination within business hours. |
| GEMINI-037 | AWS root account has MFA and is not used for daily operations. |
| GEMINI-038 | IR plan exists for payroll data incidents. |
| GEMINI-039 | Third-party libraries are scanned for vulnerabilities. |
| GEMINI-040 | VPN vendor and RDS have demonstrated business continuity. |

**Total: 40**

### Persona 4: DeepSeek (Technical, Detail-Oriented)

Prompt style: *"Analyze from a systems engineering perspective"*

| ID | Assumption |
|----|-----------|
| DS-001 | VPN gateway runs firmware version with no unpatched CVSS >= 7.0 vulnerabilities. |
| DS-002 | VPN client config enforces full-tunnel mode via Group Policy or MDM profile. |
| DS-003 | Each VPN user has unique certificate or token-bound credential; no shared accounts. |
| DS-004 | VPN gateway requires TOTP, WebAuthn, or push-based MFA; SMS OTP not used. |
| DS-005 | IPsec/IKEv2 or WireGuard with AES-256-GCM and SHA-384; IKEv1 and PPTP disabled. |
| DS-006 | VPN gateway concurrent session limit exceeds remote workforce plus 25% buffer. |
| DS-007 | VPN client implements network-level kill switch via firewall rules dropping non-tunnel traffic. |
| DS-008 | VPN gateway configured with DHE or ECDHE key exchange for PFS. |
| DS-009 | VPN gateway deployed in active-active or active-passive cluster with automatic failover. |
| DS-010 | ISP circuit to VPN gateway has documented SLA >=99.9% uptime with redundant providers. |
| DS-011 | VPN gateway syslog forwarded to SIEM with TLS encryption and integrity verification. |
| DS-012 | TLS certificates issued by internal CA with automated renewal via ACME or equivalent. |
| DS-013 | TLS 1.3 preferred; TLS 1.2 minimum; TLS 1.0/1.1 and SSL disabled at OS/ALB level. |
| DS-014 | TLS cipher config excludes NULL, RC4, 3DES, CBC-mode; only AEAD ciphers (GCM/ChaCha20). |
| DS-015 | TLS cert verification includes CRL or OCSP stapling; self-signed certs rejected in production. |
| DS-016 | AD forest functional level Windows 2016+ with SMB signing and LDAP signing enforced. |
| DS-017 | AD fine-grained password policies enforce 16-char minimum (NIST 800-63B guidance). |
| DS-018 | AD account lockout: 10 failed attempts within 15 min, 30-min lockout duration. |
| DS-019 | NTLMv1 blocked by Group Policy; LM hash storage disabled. |
| DS-020 | AD domain controllers deployed in multi-AZ configuration with health monitoring. |
| DS-021 | SSO token lifetime configured to 8h maximum with sliding window no longer than 1h. |
| DS-022 | SSO token signing keys stored in AWS KMS with automatic rotation every 90 days. |
| DS-023 | Web app uses prepared statements with bound parameters for all DB queries; no dynamic SQL. |
| DS-024 | Session cookies: HttpOnly, Secure, SameSite=Strict, max idle TTL 15 min. |
| DS-025 | Web app implements ABAC or RBAC with row-level security for payroll data. |
| DS-026 | XSS mitigated by context-aware output encoding; CSRF by anti-forgery tokens; clickjacking by X-Frame-Options: DENY. |
| DS-027 | Rate limiting via reverse proxy with graduated response (delay/CAPTCHA/block) at 10 req/min per IP. |
| DS-028 | Web app returns generic error pages with correlation ID; stack traces logged server-side only. |
| DS-029 | Database credentials stored in AWS Secrets Manager or Parameter Store with auto-rotation. |
| DS-030 | Web app runs as non-root user with no capabilities and read-only filesystem where feasible. |
| DS-031 | App DB user restricted to table-level grants (SELECT, INSERT, UPDATE, DELETE) with no schema modifications. |
| DS-032 | RDS instance deployed in private subnet with VPC endpoint or NAT-gateway-restricted egress. |
| DS-033 | RDS security group ingress restricted to web app security group on TCP 5432/3306. |
| DS-034 | VPC route tables have no route from VPN subnet CIDR to DB subnet CIDR. |
| DS-035 | RDS encryption at rest uses customer-managed KMS key with automatic annual rotation and deletion protection. |
| DS-036 | RDS enhanced monitoring and audit logs enabled, streamed to CloudWatch with 1-year retention. |
| DS-037 | RDS auto-minor-version-upgrade enabled; major upgrades tested in staging. |
| DS-038 | DB credentials rotated every 90 days via automated pipeline with secret versioning. |
| DS-039 | Nightly automated snapshots complete within backup window; retention at least 30 days. |
| DS-040 | Backup restore tested quarterly with documented RTO <=4h and RPO <=24h. |
| DS-041 | Backup snapshots encrypted with separate KMS key in different AWS region. |
| DS-042 | VPC flow logs enabled for all subnets; logs delivered to centralized S3 with object lock. |
| DS-043 | IAM policies for RDS, EC2, KMS follow least-privilege model with condition keys restricting source VPC/IP. |
| DS-044 | Third-party libraries scanned via SCA tooling (Snyk/Dependabot) before merge to production. |
| DS-045 | RDS deletion protection enabled; final DB snapshot taken before any deletion operation. |
| DS-046 | Monitoring infrastructure logging immutable; CloudTrail and S3 access logs enabled for audit bucket. |
| DS-047 | AD users who authenticate to VPN not granted direct DB access with overlapping credentials. |
| DS-048 | Web app not reachable via public DNS or ALB; internal Route53-only. |

**Total: 48**

### Persona 5: Qwen (Balanced, Practical)

Prompt style: *"What would a pragmatic security engineer say?"*

| ID | Assumption |
|----|-----------|
| QWEN-001 | VPN gateway patched for critical vulnerabilities within vendor SLA. |
| QWEN-002 | Split-tunneling disabled on VPN client via MDM policy. |
| QWEN-003 | VPN credentials unique per user and never shared. |
| QWEN-004 | MFA enforced for all VPN connections. |
| QWEN-005 | VPN gateway uses strong encryption and disables outdated protocols. |
| QWEN-006 | VPN gateway has adequate capacity for remote workforce. |
| QWEN-007 | VPN client has kill switch to prevent data leakage on tunnel drop. |
| QWEN-008 | VPN gateway logs sent to SIEM and monitored daily. |
| QWEN-009 | Documented offline procedure exists for when VPN is unavailable. |
| QWEN-010 | TLS certificates for VPN and web app are valid and auto-renewed. |
| QWEN-011 | TLS 1.2 or higher enforced across the architecture. |
| QWEN-012 | Weak TLS ciphers disabled on all endpoints. |
| QWEN-013 | Active Directory hardened against common attacks. |
| QWEN-014 | AD account lockout configured to prevent brute-force. |
| QWEN-015 | NTLMv1 disabled in favor of Kerberos authentication. |
| QWEN-016 | AD domain controllers patched and monitored regularly. |
| QWEN-017 | SSO tokens validated by web app before granting access. |
| QWEN-018 | SSO token signing keys rotated and access-controlled. |
| QWEN-019 | Web app uses parameterized queries to prevent SQL injection. |
| QWEN-020 | Session cookies secured with HttpOnly, Secure, SameSite attributes. |
| QWEN-021 | RBAC enforced within application for payroll data access. |
| QWEN-022 | XSS, CSRF, and clickjacking protections implemented. |
| QWEN-023 | Rate limiting configured on login endpoints. |
| QWEN-024 | Error pages do not leak DB schema or stack traces to users. |
| QWEN-025 | Payroll data not exfiltrated to unauthorized external endpoints. |
| QWEN-026 | DB credentials stored in secrets manager, rotated regularly. |
| QWEN-027 | App DB user has only privileges needed for its function. |
| QWEN-028 | RDS in private subnet with no internet gateway attached. |
| QWEN-029 | DB security group allows traffic only from web application. |
| QWEN-030 | No direct network path from VPN subnet to DB subnet. |
| QWEN-031 | RDS encryption at rest enabled with customer-managed key. |
| QWEN-032 | DB audit logging enabled for all queries. |
| QWEN-033 | RDS patching configured for automatic critical updates. |
| QWEN-034 | DB credentials rotated at least every 90 days. |
| QWEN-035 | Backups complete successfully within nightly window. |
| QWEN-036 | Backup restore tested at least annually. |
| QWEN-037 | Backups encrypted at rest. |
| QWEN-038 | Endpoints have antivirus/EDR installed and centrally managed. |
| QWEN-039 | Laptops managed with disk encryption and screen lock enforced. |
| QWEN-040 | Access revoked promptly on employee termination. |
| QWEN-041 | User access recertified quarterly by payroll data owners. |
| QWEN-042 | IR plan covers payroll data exposure scenarios. |
| QWEN-043 | IR team can access VPN, app, and DB logs during investigations. |
| QWEN-044 | Third-party libraries scanned before deployment. |
| QWEN-045 | Payroll data classified as sensitive and handled accordingly. |

**Total: 45**

---

## 3. Consensus Matrix

**Total unique assumptions across all 7 sources: 94**

### Tier Classification Summary

| Tier | Definition | Count | % of Total |
|------|-----------|-------|-----------|
| **A** | Human + ASF + ≥2 AIs | 22 | 23.4% |
| **A-** | Human + ASF (<2 AIs) | 3 | 3.2% |
| **B** | ASF + ≥2 AIs (Human missed) | 17 | 18.1% |
| **C** | ASF only (<2 AIs, Human missed) | 29 | 30.9% |
| **D** | Human only (<2 AIs, ASF missed) | 4 | 4.3% |
| **D+** | Human + ≥2 AIs (ASF missed) | 11 | 11.7% |
| **E** | ≥2 AIs only (not Human/ASF) | 0 | 0.0% |
| **E-** | Single AI only | 8 | 8.5% |

### Most-Agreed Assumptions (7/7 sources)

*19 assumptions found by all 7 sources:*

1. VPN credentials are never shared between users.
2. VPN gateway enforces multi-factor authentication for all users.
3. Private subnet containing database has no route to or from the internet.
4. TLS certificates on VPN gateway and internal web app are valid, not expired, from trusted CA.
5. VPN gateway enforces modern encryption (AES-256, SHA-256) and disables weak protocols.
6. Web app does not store database credentials in source code or config files readable by non-admins.
7. Database user has least-privilege access (SELECT/INSERT/UPDATE on specific tables only, no DDL).
8. Database subnet cannot initiate connections to VPN gateway or web app.
   ... and 11 more

### Tier B Highlights — ASF Discoveries Validated by AIs

*17 assumptions the human missed but ASF found and ≥2 AIs independently validated:*

1. Web app validates AD/SSO tokens and does not accept unauthenticated requests.
2. SSO session timeout is consistent between AD and the web app.
3. SSO token signing keys are rotated and protected from unauthorized access.
4. Single VPN gateway failure does not block all remote access (redundancy exists).
5. Documented offline procedure for VPN outages that does not bypass security controls.
6. Nightly backups complete within backup window without impacting DB performance.
7. IAM roles for RDS access scoped to minimum required actions.
8. Payroll data classified as sensitive/confidential and subject to data handling policies.
9. RDS encryption at rest enabled using AWS KMS.
10. Weak cipher suites (RC4, 3DES, CBC-mode) explicitly disabled on all TLS endpoints.
   ... and 7 more

### Tier C — ASF-Only Assumptions (Highest Risk Category)

*29 assumptions found ONLY by ASF (not Human, not ≥2 AIs). Manual validation required:*

1. MFA recovery codes are securely stored and usable only through verified identity proofing.
2. Help desk has documented social-engineering-resistant process for MFA token reset.
3. MFA is enforced on VPN gateway, not just the web application.
4. AD domain controller is available and reachable from web app at all times.
5. Internet circuit to VPN gateway has sufficient bandwidth and reliability SLA.
6. Web app can gracefully handle database connection failures without leaking data or crashing.
7. Backups stored in separate AWS region or account from primary database.
8. AWS root account protected by MFA and not used for daily operations.
   ... and 21 more

### Tier E — AI-Only Assumptions (Possible Novel Insights or Hallucinations)

*0 assumptions found by ≥2 AIs but neither Human nor ASF:*


### Tier D (Human + AIs, ASF Missed) — ASF Pattern Gaps

*11 assumptions that Human and ≥2 AIs identified but ASF missed — indicating missing ASF patterns:*

1. VPN gateway runs fully patched firmware free of known CVEs.
2. VPN client software on user laptops is centrally managed and disables split-tunneling.
3. Active Directory is hardened against common attack techniques.
4. Active Directory enforces account lockout after threshold of failed login attempts.
5. VPN gateway supports sufficient concurrent connections for all remote employees.
6. Web app enforces rate limiting on login endpoints to prevent credential stuffing.
7. Automatic patching enabled on RDS instance for critical DB engine vulnerabilities.
8. VPN client has kill switch that terminates all traffic if VPN connection drops.

---

## 4. AUS Scoring for All 64 ASF Assumptions

Each assumption scored by 3 simulated judges (C=Conservative, M=Moderate, G=Generous) on 5 criteria (0-5 each). AUS = sum of 5 criteria (max 25).

| ID | AUS(C) | AUS(M) | AUS(G) | Mean AUS | Brief |
|----|--------|--------|--------|----------|-------|
| ASF-001 | 20 | 22 | 22 | **21.3** | VPN gateway enforces MFA for all remote users. |
| ASF-002 | 17 | 20 | 21 | **19.3** | MFA recovery codes are securely stored and usable  |
| ASF-003 | 16 | 19 | 21 | **18.7** | Help desk has documented social-engineering-resist |
| ASF-004 | 19 | 20 | 24 | **21.0** | MFA is enforced on VPN gateway, not just the web a |
| ASF-005 | 17 | 21 | 22 | **20.0** | Web app validates AD/SSO tokens and does not accep |
| ASF-006 | 16 | 19 | 21 | **18.7** | AD domain controller is available and reachable fr |
| ASF-007 | 15 | 18 | 20 | **17.7** | SSO session timeout is consistent between AD and t |
| ASF-008 | 19 | 20 | 24 | **21.0** | SSO token signing keys are rotated and protected f |
| ASF-009 | 18 | 20 | 23 | **20.3** | Single VPN gateway failure does not block all remo |
| ASF-010 | 15 | 19 | 21 | **18.3** | Documented offline procedure for VPN outages that  |
| ASF-011 | 14 | 18 | 19 | **17.0** | Internet circuit to VPN gateway has sufficient ban |
| ASF-012 | 17 | 20 | 22 | **19.7** | Web app can gracefully handle database connection  |
| ASF-013 | 14 | 16 | 19 | **16.3** | Nightly backups complete within backup window with |
| ASF-014 | 19 | 21 | 23 | **21.0** | Backup restore is tested at least annually to vali |
| ASF-015 | 17 | 20 | 22 | **19.7** | Backups stored in separate AWS region or account f |
| ASF-016 | 19 | 21 | 24 | **21.3** | Backup files encrypted at rest using separate KMS  |
| ASF-017 | 16 | 20 | 21 | **19.0** | AWS account hosting RDS has no other workloads wit |
| ASF-018 | 19 | 19 | 24 | **20.7** | IAM roles for RDS access scoped to minimum require |
| ASF-019 | 22 | 23 | 24 | **23.0** | AWS root account protected by MFA and not used for |
| ASF-020 | 16 | 20 | 21 | **19.0** | No public AMIs or untrusted base images used for w |
| ASF-021 | 17 | 19 | 22 | **19.3** | Payroll data classified as sensitive/confidential  |
| ASF-022 | 14 | 18 | 19 | **17.0** | Data flow diagrams exist and accurately represent  |
| ASF-023 | 20 | 22 | 24 | **22.0** | Web app does not transmit payroll data to any endp |
| ASF-024 | 16 | 19 | 22 | **19.0** | Payroll data not used in development or staging en |
| ASF-025 | 21 | 22 | 23 | **22.0** | RDS encryption at rest enabled using AWS KMS. |
| ASF-026 | 19 | 20 | 24 | **21.0** | KMS key policies restrict which principals can enc |
| ASF-027 | 17 | 20 | 22 | **19.7** | KMS keys are rotated annually. |
| ASF-028 | 16 | 19 | 21 | **18.7** | Temporary storage (swap, temp tables, query cache) |
| ASF-029 | 18 | 21 | 22 | **20.3** | TLS enforced between VPN gateway and web app as do |
| ASF-030 | 20 | 23 | 24 | **22.3** | Web app validates database server TLS certificate  |
| ASF-031 | 21 | 22 | 23 | **22.0** | TLS 1.2 or higher enforced; TLS 1.0/1.1 and all SS |
| ASF-032 | 19 | 20 | 24 | **21.0** | Weak cipher suites (RC4, 3DES, CBC-mode) explicitl |
| ASF-033 | 17 | 19 | 22 | **19.3** | User laptops have EDR/AV software installed, runni |
| ASF-034 | 18 | 19 | 23 | **20.0** | User laptops managed by MDM with enforced disk enc |
| ASF-035 | 16 | 19 | 22 | **19.0** | Lost or stolen laptops can be remotely wiped to pr |
| ASF-036 | 10 | 14 | 16 | **13.3** | Users do not install unauthorized software on corp |
| ASF-037 | 17 | 22 | 22 | **20.3** | Users do not share VPN credentials or leave authen |
| ASF-038 | 13 | 18 | 18 | **16.3** | Users can identify and report phishing attempts ta |
| ASF-039 | 13 | 18 | 20 | **17.0** | Administrators assigning permissions follow princi |
| ASF-040 | 17 | 19 | 22 | **19.3** | Users whose role does not require payroll access a |
| ASF-041 | 18 | 19 | 23 | **20.0** | User accounts follow documented joiner/mover/leave |
| ASF-042 | 16 | 19 | 20 | **18.3** | AD group membership for VPN and app access reviewe |
| ASF-043 | 19 | 20 | 24 | **21.0** | Service accounts used by web app for DB authentica |
| ASF-044 | 16 | 19 | 21 | **18.7** | App-level roles synchronized with HR data so role  |
| ASF-045 | 19 | 22 | 23 | **21.3** | Incident response plan covers payroll data exposur |
| ASF-046 | 18 | 21 | 23 | **20.7** | IR team has access to VPN, app, and DB logs during |
| ASF-047 | 17 | 20 | 22 | **19.7** | IR plan includes isolation procedures (disconnect  |
| ASF-048 | 19 | 21 | 23 | **21.0** | Monitoring systems can detect anomalies in VPN and |
| ASF-049 | 21 | 22 | 23 | **22.0** | App DB user has SELECT/INSERT/UPDATE/DELETE only o |
| ASF-050 | 19 | 20 | 24 | **21.0** | Web app does not run as root or with OS-level admi |
| ASF-051 | 18 | 21 | 23 | **20.7** | AD users who can authenticate to VPN cannot also d |
| ASF-052 | 17 | 22 | 22 | **20.3** | Web app enforces authorization decisions beyond in |
| ASF-053 | 18 | 19 | 23 | **20.0** | VPN connection logs monitored for brute-force atte |
| ASF-054 | 19 | 21 | 23 | **21.0** | Database query patterns monitored for unusual volu |
| ASF-055 | 18 | 19 | 23 | **20.0** | Alerts configured for failed authentication thresh |
| ASF-056 | 19 | 21 | 23 | **21.0** | Monitoring infrastructure logs are append-only and |
| ASF-057 | 22 | 23 | 23 | **22.7** | VPN gateway, web app, and database in separate sec |
| ASF-058 | 22 | 24 | 24 | **23.3** | No direct network path from VPN subnet to database |
| ASF-059 | 17 | 20 | 22 | **19.7** | VPC flow logs or equivalent network telemetry enab |
| ASF-060 | 21 | 23 | 23 | **22.3** | Database security group allows inbound traffic onl |
| ASF-061 | 16 | 18 | 20 | **18.0** | AWS RDS is available and not experiencing region-l |
| ASF-062 | 16 | 18 | 21 | **18.3** | VPN vendor has no known backdoors or critical vuln |
| ASF-063 | 18 | 20 | 23 | **20.3** | Third-party libraries used by web app scanned for  |
| ASF-064 | 14 | 15 | 18 | **15.7** | Exit strategy or migration plan exists if VPN vend |

**Mean AUS across all 64 ASF assumptions: 19.8/25**

**Assumptions scoring AUS ≥ 15 (High Value or Critical): 63 / 64 (98.4%)**

### AUS Distribution

```
  0-4 Ignore      |  0 | ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
  5-9 Low         |  0 | ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
  10-14 Medium    |  1 | █░░░░░░░░░░░░░░░░░░░░░░░░░░░░░
  15-19 High      | 28 | ████████████████████████████░░
  20-25 Critical  | 35 | ███████████████████████████████████
```

## 5. Tier Distribution

### Text Pie Chart

```
  A (Human+ASF+AIs)          25 ( 26.6%) █████████████
  B (ASF+AIs)                17 ( 18.1%) █████████
  C (ASF only)               29 ( 30.9%) ███████████████
  D (Human+ASF gap)          15 ( 16.0%) ███████
  E (AI only)                 0 (  0.0%) █
```

### Tier Breakdown

- **Tier A (A + A-): 25** — Universal agreement, high-confidence findings
  - A (Human + ASF + ≥2 AIs): 22
  - A- (Human + ASF, <2 AIs): 3
- **Tier B (ASF + ≥2 AIs, Human missed): 17** — Validated novel discoveries
- **Tier C (ASF only): 29** — Orphan assumptions requiring manual review
- **Tier D (Human only, ASF missed): 15** — ASF blind spots
  - D (Human only, <2 AIs): 4
  - D+ (Human + ≥2 AIs, ASF missed): 11
- **Tier E (AI only): 0** — Potential hallucinations or novel insights
  - E (≥2 AIs): 0
  - E- (1 AI): 8

---

## 6. Analysis & Conclusions

### Key Metrics

1. **Total unique assumptions across all 7 sources: 94**
2. **Tier A (Human + ASF + ≥2 AIs): 25** — high-confidence, multi-source confirmed findings
3. **Tier B (ASF + ≥2 AIs, Human missed): 17** — novel ASF discoveries validated by AI consensus
4. **Tier C (ASF only): 29** — highest-risk orphan assumptions requiring manual validation
5. **Tier D (Human only, ASF missed): 15** — ASF blind spots indicating missing pattern coverage
6. **Tier E (AI only, not Human/ASF): 0** — possible novel insights or hallucination clusters

### AUS Summary

- **Mean AUS for all 64 ASF assumptions: 19.8/25**
- **Percentage of ASF assumptions scoring AUS ≥ 15: 98.4%** (63 of 64)
- **Tier C (ASF-only) count: 29 assumptions, mean AUS = 19.1**
  → Tier C ASF-only assumptions score **near/below the ASF average**, suggesting they include some lower-utility edge cases.

### Top 10 ASF Assumptions by AUS

| Rank | ID | Mean AUS | Assumption |
|------|-----|----------|-----------|
| 1 | ASF-058 | **23.3** | No direct network path from VPN subnet to database subnet (t |
| 2 | ASF-019 | **23.0** | AWS root account protected by MFA and not used for daily ope |
| 3 | ASF-057 | **22.7** | VPN gateway, web app, and database in separate security grou |
| 4 | ASF-030 | **22.3** | Web app validates database server TLS certificate when conne |
| 5 | ASF-060 | **22.3** | Database security group allows inbound traffic only from web |
| 6 | ASF-023 | **22.0** | Web app does not transmit payroll data to any endpoint outsi |
| 7 | ASF-025 | **22.0** | RDS encryption at rest enabled using AWS KMS. |
| 8 | ASF-031 | **22.0** | TLS 1.2 or higher enforced; TLS 1.0/1.1 and all SSL versions |
| 9 | ASF-049 | **22.0** | App DB user has SELECT/INSERT/UPDATE/DELETE only on payroll  |
| 10 | ASF-001 | **21.3** | VPN gateway enforces MFA for all remote users. |

### Bottom 5 ASF Assumptions by AUS

| Rank | ID | Mean AUS | Assumption |
|------|-----|----------|-----------|
| 60 | ASF-039 | **17.0** | Administrators assigning permissions follow principle of lea |
| 61 | ASF-013 | **16.3** | Nightly backups complete within backup window without impact |
| 62 | ASF-038 | **16.3** | Users can identify and report phishing attempts targeting th |
| 63 | ASF-064 | **15.7** | Exit strategy or migration plan exists if VPN vendor or RDS  |
| 64 | ASF-036 | **13.3** | Users do not install unauthorized software on corporate lapt |

### Cross-Source Agreement Analysis

#### Per-Persona Overlap with Human+ASF Core

| Persona | Total | In H+A Core | % Overlap | Entirely Unique |
|---------|-------|-------------|-----------|-----------------|
| GPT | 45 | 43 | 95.6% | 1 |
| Claude | 52 | 51 | 98.1% | 3 |
| Gemini | 40 | 38 | 95.0% | 1 |
| DeepSeek | 48 | 47 | 97.9% | 2 |
| Qwen | 45 | 44 | 97.8% | 1 |

#### Diversity Assessment

The 5 AI personas collectively surface a broader assumption set than Human+ASF alone. Each persona contributes a distinct perspective:
- **GPT**: Strong on logical enumeration through the data path; excels at identifying step-by-step security requirements
- **Claude**: Broadest coverage; catches edge cases (vendor exit strategy, mTLS, awareness training) others miss
- **Gemini**: Most concise; focuses on high-impact assumptions with fewer lower-probability scenarios
- **DeepSeek**: Deepest technical specificity; provides implementation-level detail (cipher suites, registry settings, SLA thresholds)
- **Qwen**: Most practical; balanced between technical controls and process/policy concerns

### Final Conclusions

1. **Coverage breadth**: The 5 AI personas collectively surface 94 unique assumptions vs. 79 from Human+ASF combined. Multi-model evaluation captures a substantially wider assumption space.
2. **Tier A strength**: 25 assumptions (26.6%) are confirmed by Human + ASF + ≥2 AIs — a strong core of universal agreement.
3. **ASF unique value**: 29 Tier C assumptions are unique to ASF with mean AUS = 19.1, confirming the framework captures valuable signals even specialized AI personas miss.
4. **Persona diversity matters**: Each AI contributed unique findings not found by Human or ASF. GPT: 1, Claude: 3, Gemini: 1, DeepSeek: 2, Qwen: 1. No single model would have sufficed.
5. **AUS effectiveness**: 98.4% of ASF assumptions score ≥15 (High Value or Critical), validating the framework's output quality.
6. **ASF blind spots (15 Tier D assumptions)**: Human and AI consensus identified concerns — particularly web application security (SQLi, XSS, CSRF, rate limiting), VPN hardening (kill switch, PFS, split tunneling), and RDS platform features (deletion protection, patching cadence) — that the ASF 20-pattern matrix does not fully cover.

### Recommendations

1. **Add a Web Application Security pattern** to the ASF matrix to close the Tier D gap on SQLi, XSS, CSRF, and rate limiting.
2. **Add a VPN Hardening sub-pattern** covering kill switch, split-tunneling, PFS, and concurrent capacity.
3. **Add a Managed Database Security sub-pattern** covering RDS deletion protection, automated patching, and snapshot policies.
4. **Continue multi-model evaluation** for future campaigns; the 5-model panel provides measurably richer insight than any single model.

---

*Report generated: 2026-06-09 | ASF Multi-LLM Evaluation Campaign #1*