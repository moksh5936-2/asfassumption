# Architecture 4 — Security Review Deliverable
## Healthcare → PHI → HIPAA

**Reviewer**: Security Architecture Team
**Date**: June 9, 2026
**Scope**: Architecture 4 — Patient Portal → Auth0 → App Server → PHI Database, with SIEM, Backup, KMS, and Third-Party Integrations

---

## 1. Consensus Matrix

| # | Assumption | GPT | Gemini | Gemma | Keep? |
|---|-----------|:---:|:------:|:-----:|:-----:|
| A1 | Patient devices are trusted | ✓ | | | ✓ |
| A2 | Patient Portal properly validates user identity | ✓ | | | ✓ |
| A3 | Auth0 is securely configured | ✓ | | ✓ | ✓ |
| A4 | Auth0 infrastructure is secure | ✓ | | | ✓ |
| A5 | Administrative access to Auth0 is protected | ✓ | ✓ | | ✓ |
| A6 | Strong authentication methods (MFA) are used | ✓ | | | ✓ |
| A7 | Password recovery processes are secure | ✓ | | | ✓ |
| A8 | TLS protects Portal-to-Application communication | ✓ | | | ✓ |
| A9 | TLS certificates are trusted and valid | ✓ | | | ✓ |
| A10 | Session management is secure | ✓ | | | ✓ |
| A11 | Session identifiers are unpredictable | ✓ | | | ✓ |
| A12 | Session expiration is properly enforced | ✓ | | | ✓ |
| A13 | Application authorization is correctly implemented | ✓ | ✓ | | ✓ |
| A14 | Object-level authorization is enforced (IDOR/BOLA) | ✓ | ✓ | | ✓ |
| A15 | Minimum necessary access is correctly implemented | ✓ | | | ✓ |
| A16 | Employees receive only necessary permissions | ✓ | | | ✓ |
| A17 | Privileged access is tightly controlled | ✓ | | | ✓ |
| A18 | User provisioning is accurate | ✓ | | | ✓ |
| A19 | User deprovisioning is timely | ✓ | | | ✓ |
| A20 | Application input validation is effective | ✓ | | | ✓ |
| A21 | Database queries are parameterized | ✓ | | | ✓ |
| A22 | App servers are hardened | ✓ | | | ✓ |
| A23 | Secrets and credentials are securely managed | ✓ | | | ✓ |
| A24 | Database credentials are protected | ✓ | | | ✓ |
| A25 | PHI Database access is restricted to authorized systems | ✓ | | | ✓ |
| A26 | Database permissions follow least privilege | ✓ | | | ✓ |
| A27 | Database software is patched | ✓ | | | ✓ |
| A28 | AES-256 encryption is implemented correctly | ✓ | | | ✓ |
| A29 | Encryption keys are protected | ✓ | ✓ | ✓ | ✓ |
| A30 | Key rotation procedures exist | ✓ | | | ✓ |
| A31 | Backups are encrypted | ✓ | | | ✓ |
| A32 | Backup access is restricted | ✓ | | | ✓ |
| A33 | Audit logs are complete | ✓ | | | ✓ |
| A34 | Audit logs cannot be altered (immutable) | ✓ | ✓ | | ✓ |
| A35 | Audit logging captures all PHI access events | ✓ | | | ✓ |
| A36 | SIEM receives logs reliably | ✓ | | ✓ | ✓ |
| A37 | SIEM integrity is maintained | ✓ | | ✓ | ✓ |
| A38 | Monitoring personnel review alerts | ✓ | | | ✓ |
| A39 | Business Associate Agreements are honored | ✓ | | | ✓ |
| A40 | Subprocessors maintain equivalent security controls | ✓ | | | ✓ |
| A41 | Insider personnel are trustworthy | ✓ | | ✓ | ✓ |
| A42 | Workforce security training is effective | ✓ | | | ✓ |
| A43 | Physical security protects infrastructure | ✓ | | | ✓ |
| A44 | Time synchronization is accurate | ✓ | | | ✓ |
| A45 | Disaster recovery processes function correctly | ✓ | | | ✓ |
| A46 | Data retention policies are followed | ✓ | | | ✓ |
| A47 | No unauthorized data exports occur | ✓ | | | ✓ |
| A48 | Third-party integrations are secure | ✓ | | | ✓ |
| A49 | No undocumented PHI access paths exist | ✓ | | | ✓ |
| A50 | Regulatory controls are correctly mapped to technical controls | ✓ | | | ✓ |
| A51 | App Server does not cache unencrypted PHI to local disk, temp files, swap, or crash dumps | | | ✓ | ✓ |
| A52 | Insider threats (developers/DevOps) cannot access database encryption keys from production configuration | | | ✓ | ✓ |
| A53 | Auth0 configuration prevents token substitution and replay attacks | | | ✓ | ✓ |
| A54 | Compromised App Server cannot block log delivery to SIEM or alter audit logs | | ✓ | ✓ | ✓ |
| A55 | Encryption keys are stored and managed outside the database in a separate KMS | | ✓ | | ✓ |

**Total Assumptions: 55** (50 GPT + 1 Gemini-exclusive + 3 Gemma-exclusive + 1 shared Gemini/Gemma)

---

## 2. Deduplicated Assumption List

### 2.1 Patient Portal
1. Patient devices are trusted (A1)
2. Patient Portal properly validates user identity (A2)
3. TLS protects Portal-to-Application communication (A8)
4. TLS certificates are trusted and valid (A9)
5. Session management is secure — generation, lifecycle, and termination (A10)
6. Session identifiers are cryptographically random and unpredictable (A11)
7. Session expiration is properly enforced (A12)
8. Password recovery processes are secure (A7)

### 2.2 Auth0
9. Auth0 is securely configured (A3)
10. Auth0 infrastructure is secure (A4)
11. Administrative access to Auth0 is protected — compromised admin portal can alter authentication flows (A5)
12. Strong authentication (MFA) is enforced for all users (A6)
13. Auth0 configuration prevents token substitution and replay attacks (A53)

### 2.3 Authorization / IDOR
14. Application authorization is correctly implemented on every API request (A13)
15. Object-level authorization (IDOR/BOLA) is enforced — patients cannot access records belonging to others (A14)
16. Minimum necessary access is correctly implemented (A15)
17. Employees receive only necessary permissions for their role (A16)
18. Privileged access is tightly controlled (A17)

### 2.4 App Security
19. User provisioning is accurate and aligned with job function (A18)
20. User deprovisioning occurs promptly upon termination or role change (A19)
21. Application input validation is effective against injection attacks (A20)
22. Database queries are parameterized to prevent SQL injection (A21)
23. App servers are hardened against compromise (A22)
24. Secrets (credentials, tokens, API keys) are securely managed (A23)
25. App Server does not cache unencrypted PHI to local disk, temp files, swap, or crash dumps (A51)

### 2.5 Database / PHI Storage
26. Database credentials are protected from unauthorized access (A24)
27. PHI Database access is restricted to authorized services and users (A25)
28. Database permissions follow least privilege — app accounts lack administrative rights (A26)
29. Database software is kept current with security patches (A27)
30. No undocumented PHI access paths exist — all data flows are accounted for (A49)

### 2.6 Encryption / Key Management
31. AES-256 encryption is implemented correctly with proper algorithm and mode selection (A28)
32. Encryption keys are protected and stored outside the database in a separate KMS (A29, A55)
33. Key rotation procedures exist and are enforced on a regular schedule (A30)
34. Backups are encrypted using strong encryption (A31)

### 2.7 Backup / Disaster Recovery
35. Backup access is restricted to authorized personnel and systems (A32)
36. Disaster recovery processes function correctly and preserve PHI confidentiality (A45)

### 2.8 Logging / Monitoring
37. Audit logs are complete and capture all HIPAA-relevant events (A33)
38. Audit logs are immutable — write-once, append-only, cannot be modified by a compromised App Server (A34, A54)
39. Audit logging captures all PHI access events including view, create, modify, and delete (A35)
40. SIEM receives logs reliably — log delivery cannot be blocked or suppressed (A36)
41. SIEM integrity is maintained — compromised systems cannot inject false log data (A37)
42. Monitoring personnel review alerts in a timely manner (A38)
43. System clocks across Portal, App Server, DB, and SIEM are synchronized via NTP (A44)

### 2.9 HIPAA Compliance
44. Business Associate Agreements are honored by all subprocessors (A39)
45. Subprocessors maintain equivalent security controls (A40)
46. Data retention policies are followed in accordance with HIPAA requirements (A46)
47. No unauthorized data exports occur — PHI does not leave controlled environments (A47)
48. Regulatory controls are correctly mapped to technical controls — compliance documents produce real security (A50)
49. Physical security protects infrastructure hosting PHI (A43)

### 2.10 Third-Party
50. Third-party integrations are secure and do not introduce indirect PHI exposure (A48)

### 2.11 Insider Threat
51. Insider personnel (employees, contractors, admins) are subject to compensating controls (A41)
52. Workforce security training is effective against phishing and accidental disclosure (A42)
53. Insider threats (developers, DevOps) cannot access database encryption keys from production configuration files or environment variables (A52)

---

## 3. Risk Scores

| # | Assumption | Likelihood | Impact | Risk |
|---|-----------|:----------:|:------:|:----:|
| A1 | Patient devices trusted | L | H | H |
| A2 | Portal validates user identity | L | H | H |
| A3 | Auth0 securely configured | M | C | C |
| A4 | Auth0 infrastructure secure | L | C | H |
| A5 | Auth0 admin access protected | M | C | C |
| A6 | MFA enforced | M | H | H |
| A7 | Password recovery secure | M | H | H |
| A8 | TLS Portal-App | L | H | H |
| A9 | TLS certificates trusted | L | H | H |
| A10 | Session management secure | L | H | H |
| A11 | Session IDs unpredictable | L | H | H |
| A12 | Session expiration enforced | L | H | H |
| A13 | Application authorization correct | M | C | C |
| A14 | Object-level authorization (IDOR/BOLA) | M | C | C |
| A15 | Minimum necessary access | M | H | H |
| A16 | Employee permissions correct | M | H | H |
| A17 | Privileged access controlled | M | C | C |
| A18 | User provisioning accurate | M | H | H |
| A19 | User deprovisioning timely | H | H | H |
| A20 | Input validation effective | M | H | H |
| A21 | Queries parameterized | L | H | H |
| A22 | App servers hardened | M | C | C |
| A23 | Secrets securely managed | M | C | C |
| A24 | DB credentials protected | L | C | H |
| A25 | PHI DB access restricted | L | C | H |
| A26 | DB least privilege | M | C | C |
| A27 | DB software patched | M | H | H |
| A28 | AES-256 correct impl | L | H | H |
| A29 | Encryption keys protected | M | C | C |
| A30 | Key rotation exists | M | H | H |
| A31 | Backups encrypted | L | C | H |
| A32 | Backup access restricted | L | C | H |
| A33 | Audit logs complete | M | C | C |
| A34 | Audit logs immutable | M | C | C |
| A35 | Logs capture all PHI access | M | C | C |
| A36 | SIEM receives logs reliably | M | C | C |
| A37 | SIEM integrity maintained | M | C | C |
| A38 | Monitoring reviews alerts | H | H | H |
| A39 | BAAs honored | M | H | H |
| A40 | Subprocessor controls | H | H | H |
| A41 | Insider trustworthy | L | C | H |
| A42 | Security training effective | H | H | H |
| A43 | Physical security | L | H | H |
| A44 | Time sync accurate | L | H | H |
| A45 | DR processes function | M | H | H |
| A46 | Data retention followed | M | H | H |
| A47 | No unauthorized exports | H | C | C |
| A48 | Third-party integrations secure | M | H | H |
| A49 | No undocumented PHI paths | H | C | C |
| A50 | Regulatory→technical mapping | M | H | H |
| A51 | PHI not cached to disk/temp/swap | M | C | C |
| A52 | Dev/DevOps access to keys prevented | M | C | C |
| A53 | Token substitution/replay prevented | L | C | H |
| A54 | Compromised App Server cannot block SIEM | M | C | C |
| A55 | Keys stored outside DB in KMS | M | C | C |

---

## 4. STRIDE Mapping

### Spoofing
- A1: Patient device identity (device trust)
- A2: Patient identity validation
- A3, A53: Auth0 configuration — authentication flow integrity
- A6: MFA enforcement
- A7: Password recovery workflow
- A10–A12: Session management / identifier generation

### Tampering
- A8, A9: TLS — man-in-the-middle protection
- A13, A14: Authorization enforcement — modified requests
- A20: Input validation — injection payloads
- A21: Parameterized queries — SQL injection prevention
- A22: App server hardening — binary integrity
- A23: Secret integrity
- A28, A29, A55: Encryption implementation and key protection
- A31: Backup encryption integrity
- A33–A35: Audit log completeness and immutability
- A51: PHI written to disk/temp/swap (unintended persistence)
- A54: Log modification or blocking by compromised App Server

### Repudiation
- A33: Complete audit trail
- A34, A54: Immutable audit logs — non-repudiation
- A35: PHI access event capture
- A37: SIEM data integrity
- A44: NTP-synchronized timestamps

### Information Disclosure
- A1, A2: Device trust and identity
- A8, A9: TLS encryption in transit
- A13, A14: Authorization / IDOR failures
- A15–A17: Access control boundaries
- A24–A27: Database-level PHI protection
- A28, A29, A55: Encryption and key management
- A31, A32: Backup encryption and access
- A47: Unauthorized data exports
- A49: Undocumented PHI data flows
- A51: PHI leakage via disk/temp/swap/crash dumps

### Denial of Service
- A27: Database patch management (availability impact)
- A36: SIEM log delivery reliability
- A45: Disaster recovery processes
- A54: Log blocking — denial of monitoring visibility

### Elevation of Privilege
- A5: Auth0 administrative access
- A13, A14: Authorization / IDOR — horizontal and vertical escalation
- A16, A17: Permission boundaries
- A18, A19: Provisioning / deprovisioning failures
- A22: App server compromise
- A23: Secret theft enabling lateral movement
- A25, A26: Database permission escalation
- A52: Developer/DevOps access to encryption keys
- A55: Key location — keys stored with data are effectively public

---

## 5. Top 10 Critical Assumptions (Ranked)

### 1. Object-Level Authorization Enforcement (IDOR/BOLA) — A14
**Rationale**: Healthcare applications expose PHI by patient identifier (e.g., `/api/patient/123`). If the App Server does not independently verify that the authenticated user owns the requested record ID, any authenticated patient (or attacker with a valid token) can enumerate and exfiltrate every patient record in the system. This is the single most common PHI exposure vector in web-based healthcare systems and directly violates HIPAA Minimum Necessary and Access Control requirements.

### 2. PHI Cached to Local Disk, Temp Files, Swap, or Crash Dumps — A51
**Rationale**: The App Server processes decrypted PHI in memory. If the application framework, JVM/runtime, or OS writes memory pages to disk (swap, crash dumps, temp files, logging frameworks), unencrypted PHI persists on storage outside the control of the encryption layer. A subsequent server compromise, disk forensic analysis, or improper decommissioning exposes PHI in plaintext — bypassing all database-level encryption controls.

### 3. Developer/DevOps Access to Encryption Keys from Production Configuration — A52
**Rationale**: In many deployments, database encryption keys are stored in configuration files, environment variables, or secrets managers accessible to developers and DevOps personnel. Unlike database-level access controls, key material grants the ability to decrypt the entire PHI database offline. A single compromised developer workstation or CI/CD pipeline can lead to bulk PHI exfiltration without triggering database audit logs.

### 4. Auth0 Token Substitution and Replay Attacks — A53
**Rationale**: If Auth0 is misconfigured to accept unsigned, unencrypted, or replayable tokens, an attacker who intercepts a single token (via network capture, browser storage, or log file) can impersonate any user. In a healthcare context, this enables full access to that patient's PHI — and if the intercepted token belongs to a healthcare provider or administrator, access to the entire patient population.

### 5. Audit Log Tampering by Compromised App Server — A34, A54
**Rationale**: If the App Server has write access to audit logs or can block log delivery to the SIEM, a server compromise enables the attacker to erase all forensic evidence of PHI access. HIPAA requires audit controls that record and examine all PHI access; if logs can be modified or suppressed by the very system being audited, the audit requirement is functionally defeated. Compensating controls require log shipping with append-only semantics and an independent logging path.

### 6. Encryption Keys Stored Alongside Encrypted Data — A29, A55
**Rationale**: If AES-256 encryption keys are stored in the same database as the encrypted PHI (e.g., in a `keys` table or configuration row), encryption provides no实质性protection. An attacker who breaches the database (via SQL injection, compromised credentials, or backup theft) retrieves both the ciphertext and the decryption key in a single operation. Keys must reside in a separate KMS with independent access controls.

### 7. Auth0 Administrative Account Compromise — A5
**Rationale**: Auth0 tenant administrators can create users, modify authentication policies, reset passwords, and alter MFA requirements. A compromised Auth0 admin account (via credential theft, session hijacking, or insider threat) provides persistent, undetected access to every application user's account — including healthcare providers and patients — without triggering any application-level alerts.

### 8. Database Least Privilege Violation — A26
**Rationale**: If the App Server's database user has excessive privileges (e.g., `INSERT`, `UPDATE`, `DELETE` on all tables, or DDL rights), a SQL injection vulnerability or compromised application server can be used to modify or delete PHI at the database level. Under HIPAA, this impacts integrity controls (45 CFR §164.312(c)(1)) and can result in widespread data corruption, malicious record modification, or complete database destruction.

### 9. Insider Threat — Authorized Personnel PHI Abuse — A41
**Rationale**: Healthcare employees with legitimate database access can browse, export, or sell PHI without triggering technical controls if access monitoring relies on audit logs that they can also modify. The 2015 UCLA Health System breach ($7.5M settlement), where employees accessed celebrity medical records, demonstrates that technical controls against insider abuse remain one of the most difficult and costly HIPAA compliance challenges.

### 10. SIEM Log Delivery Blocked by Compromised App Server — A36, A54
**Rationale**: If the App Server controls log delivery (e.g., the logging agent runs on the App Server and pushes to SIEM), a compromised App Server can stop log shipment without triggering an alert. This creates a blind window during which the attacker can access PHI, exfiltrate data, and cover their tracks — all while the SIEM reports "no events" and monitoring personnel observe no anomalies.

---

## 6. Recommended Controls

### 1. Object-Level Authorization (IDOR/BOLA)
- Implement centralized authorization middleware that validates resource ownership on every API request
- Use opaque, unguessable resource identifiers (UUIDs) instead of sequential integers
- Enforce authorization at the API gateway and again at the controller layer (defense in depth)
- Conduct automated IDOR scanning in CI/CD with tools like Authz0 or custom property-based tests
- Perform quarterly penetration testing focused on horizontal privilege escalation between patient accounts

### 2. PHI Temp File / Swap / Crash Dump Prevention
- Configure the App Server runtime to disable swapping and core dumps for the PHI-processing process
- Use `mlock()` / `VirtualLock()` to prevent memory pages containing PHI from being paged to disk
- Mount temp directories with `noexec` and encryption; clear temp files after each request
- Disable crash dump generation for production App Server processes, or encrypt crash dumps
- Audit filesystem writes from the App Server process to detect unapproved PHI persistence

### 3. Developer/DevOps Key Access Prevention
- Deploy a dedicated Hardware Security Module (HSM) or cloud KMS (AWS KMS, Azure Key Vault, GCP Cloud KMS) with role-based access controls
- Grant key usage permissions to the App Server's runtime role only — never to human users
- Implement just-in-time (JIT) key access for break-glass scenarios with approval workflow
- Audit all key access events and alert on any key usage by a human principal
- Store key access policies in a separate administrative domain from application configuration

### 4. Auth0 Token Security
- Enforce RS256 or ES256 signing with a private key stored in Auth0's secure key store
- Enable short token expiration (15–30 minutes) and implement refresh token rotation
- Configure token binding to client TLS session to prevent replay
- Disable token forwarding and implement audience/issuer validation on every API request
- Regularly audit Auth0 tenant configuration for misconfigured token policies

### 5. Immutable Audit Logs
- Ship logs via a separate, hardened log shipper running on a dedicated log-collection host (not the App Server)
- Store logs in WORM/immutable storage (AWS S3 Object Lock, Azure Immutable Blob, write-once filesystem)
- Implement cryptographic hash chaining (Linked Timestamping) to detect log tampering
- Alert on any gap or delay in log delivery exceeding 60 seconds
- Maintain a separate, air-gapped log monitoring console with independent administrative credentials

### 6. Key-Data Separation
- Never store encryption keys in the same database, filesystem, or backup as the encrypted PHI
- Use a separate KMS or HSM with distinct authentication and authorization domains
- Implement key hierarchy: data encryption keys (DEKs) encrypted by a key encryption key (KEK) stored in the KMS
- Rotate DEKs on a 90-day schedule and KEKs annually
- Audit all key creation, usage, rotation, and destruction events

### 7. Auth0 Administrative Access Protection
- Enforce MFA for all Auth0 tenant administrative accounts
- Implement IP allowlisting for the Auth0 admin portal
- Use separate administrator accounts for each environment (dev, staging, prod)
- Enable Auth0's built-in breach detection and anomalous login alerts
- Conduct quarterly reviews of administrative access with identity team sign-off

### 8. Database Least Privilege
- Create separate database users for read, write, and admin operations with minimal required grants
- Use row-level security (RLS) policies to enforce patient-scoped access at the database level
- Implement database firewall rules restricting the App Server's source IP and query patterns
- Conduct automated permission reviews using tools like `pg_permissions` or `sql_audit`
- Remove all DDL grants from application database users

### 9. Insider Threat Compensating Controls
- Implement user behavior analytics (UBA) on PHI access patterns — flag bulk exports, off-hours access, and unusual query patterns
- Deploy database activity monitoring (DAM) with real-time alerting on anomalous queries
- Enforce dual-control for PHI export operations (two authorized users must approve)
- Require documented business justification for all PHI access to patient records not assigned to the user
- Conduct quarterly access recertification with automated deprovisioning of unconfirmed access

### 10. SIEM Log Delivery Reliability
- Deploy a log shipper on a separate, hardened bastion host with network connectivity to both App Server and SIEM
- Implement a dead-letter queue (DLQ) or local buffer to cache logs if SIEM is unreachable
- Monitor log shipping heartbeat — alert on any log delivery gap exceeding 30 seconds
- Use a separate, out-of-band management network for log transport to the SIEM
- Conduct monthly log delivery failover testing

---

## 7. Summary Statistics

| Category | Count |
|----------|:-----:|
| **Total Assumptions** | **55** |
| **Critical Risk** | 20 |
| **High Risk** | 35 |
| **Medium Risk** | 0 |
| **Low Risk** | 0 |
| **Sources** | 3 models (GPT-4o, Gemini, Gemma) |
| **Exclusive (GPT only)** | 41 |
| **Exclusive (Gemini only)** | 1 |
| **Exclusive (Gemma only)** | 3 |
| **Shared (Gemini + Gemma, not GPT)** | 1 |

**Critical Assumptions** (20): A3, A5, A13, A14, A17, A22, A23, A26, A29, A33, A34, A35, A36, A37, A47, A49, A51, A52, A54, A55

**Top Risk Drivers**: IDOR/BOLA (patients accessing each other's PHI), PHI leakage via App Server disk caching/temp files/swap, developer/DevOps access to encryption keys from production configuration, Auth0 token substitution/replay attacks, and audit log tampering by a compromised App Server — these five failure modes represent the highest concentration of unmitigated PHI exposure risk in Architecture 4. Each individually defeats core HIPAA Privacy and Security Rule requirements; their combined failure would render the entire PHI protection framework inoperable.

---

*End of Security Review — Architecture 4*
