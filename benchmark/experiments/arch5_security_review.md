# Architecture 5 — Security Review Deliverable
## ERP → SOX → Financial Reporting → Audit

**Reviewer**: Security Architecture Team
**Date**: June 9, 2026
**Scope**: Architecture 5 — Finance Team → ERP Web App → ERP Backend → Financial DB, with Approval Workflow, Audit Logs, Reporting Engine → Auditor Access

---

## 1. Consensus Matrix

| # | Assumption | GPT | Gemini | Gemma | Keep? |
|---|-----------|:---:|:------:|:-----:|:-----:|
| A1 | Finance users are properly authenticated | ✓ | | | ✓ |
| A2 | User credentials are protected (storage, transit, hashing) | ✓ | | | ✓ |
| A3 | MFA or equivalent protections exist for ERP access | ✓ | | | ✓ |
| A4 | Session management is secure (generation, lifecycle, termination) | ✓ | | | ✓ |
| A5 | Session identifiers are cryptographically random and unpredictable | ✓ | | | ✓ |
| A6 | Authorization is implemented correctly at every layer | ✓ | | | ✓ |
| A7 | User roles accurately reflect job responsibilities | ✓ | | | ✓ |
| A8 | Role assignments are reviewed periodically | ✓ | | | ✓ |
| A9 | User provisioning is accurate and timely | ✓ | | | ✓ |
| A10 | User deprovisioning occurs promptly upon role change or termination | ✓ | | | ✓ |
| A11 | Segregation-of-duties (SoD) rules are correctly defined | ✓ | | | ✓ |
| A12 | SoD controls are technically enforced, not merely procedural | ✓ | | | ✓ |
| A13 | No user can both create and approve the same transaction | ✓ | | | ✓ |
| A14 | Approval workflows cannot be bypassed | ✓ | ✓ | | ✓ |
| A15 | Approval decisions are authenticated (non-repudiable) | ✓ | | ✓ | ✓ |
| A16 | Approval records cannot be modified or deleted after creation | ✓ | | | ✓ |
| A17 | Workflow logic is correctly implemented with no edge-case flaws | ✓ | | | ✓ |
| A18 | Emergency override procedures are logged, controlled, and time-bound | ✓ | | | ✓ |
| A19 | ERP application code is free of exploitable vulnerabilities | ✓ | | | ✓ |
| A20 | Input validation is effective against injection attacks | ✓ | ✓ | | ✓ |
| A21 | Backend APIs independently enforce authorization (not just front-end) | ✓ | ✓ | | ✓ |
| A22 | Administrative functions are restricted to authorized roles | ✓ | | | ✓ |
| A23 | Administrative actions are logged and auditable | ✓ | | | ✓ |
| A24 | Financial database access is restricted to authorized services/users | ✓ | | | ✓ |
| A25 | Database permissions follow the principle of least privilege | ✓ | | | ✓ |
| A26 | Database integrity protections exist (constraints, checksums, immutability) | ✓ | | | ✓ |
| A27 | Database software is kept current with security patches | ✓ | | | ✓ |
| A28 | Secrets and credentials (service accounts, API keys) are securely stored | ✓ | | | ✓ |
| A29 | Encryption protects sensitive financial data at rest and in transit | ✓ | | | ✓ |
| A30 | Encryption keys are managed with proper lifecycle controls | ✓ | | | ✓ |
| A31 | Audit logs are complete and capture all SOX-relevant events | ✓ | | | ✓ |
| A32 | Audit logs cannot be altered (immutable storage, append-only) | ✓ | | | ✓ |
| A33 | Audit logs capture all journal-entry activity (create, modify, approve) | ✓ | | | ✓ |
| A34 | System clocks across ERP, Approval Workflow, and Audit Logs are synchronized via NTP | ✓ | ✓ | | ✓ |
| A35 | Reporting engine faithfully reflects source data without transformation errors | ✓ | | | ✓ |
| A36 | Reporting logic is accurate and free of calculation errors | ✓ | | | ✓ |
| A37 | Data transferred from ERP to reporting maintains integrity in transit | ✓ | | | ✓ |
| A38 | Reporting engine access is controlled to authorized users | ✓ | | | ✓ |
| A39 | Auditor accounts are properly authenticated | ✓ | | | ✓ |
| A40 | Read-only access is technically enforced for auditors at the application level | ✓ | | | ✓ |
| A41 | Read-only restrictions for auditors cannot be bypassed via alternate paths | ✓ | | | ✓ |
| A42 | Auditors cannot escalate privileges from read-only to write | ✓ | ✓ | | ✓ |
| A43 | Quarterly recertification of user access is performed and meaningful | ✓ | | | ✓ |
| A44 | Recertification reviews are completed on schedule | ✓ | | | ✓ |
| A45 | Exceptions identified during recertification are remediated in a timely manner | ✓ | | | ✓ |
| A46 | Insider personnel (finance, IT, admin) are trustworthy | ✓ | | | ✓ |
| A47 | Change management prevents unauthorized modifications to financial controls | ✓ | | | ✓ |
| A48 | Backup and recovery processes preserve financial data integrity | ✓ | | | ✓ |
| A49 | Disaster recovery procedures maintain SOX control effectiveness during outages | ✓ | | | ✓ |
| A50 | No undocumented financial processing paths exist | ✓ | | | ✓ |
| A51 | System administrators and DBAs are subject to the same SoD as the finance team | | ✓ | | ✓ |
| A52 | Reporting Engine must have read-only database permissions at the database level | | | ✓ | ✓ |
| A53 | User identity inside the approval workflow cannot be spoofed via session fixation or ID hijacking | | | ✓ | ✓ |
| A54 | Quarterly recertification data must come from a tamper-proof system of record | | | ✓ | ✓ |
| A55 | Low-level users cannot modify session parameters to access admin functions or bypass approval chains | | | ✓ | ✓ |

**Total Assumptions: 55** (50 GPT + 1 Gemini-exclusive + 4 Gemma-exclusive)

---

## 2. Deduplicated Assumption List

### 2.1 Authentication
1. Finance users are properly authenticated (A1)
2. User credentials are protected in storage and transit (A2)
3. MFA or equivalent protections exist for ERP access (A3)
4. Session management is secure — generation, lifecycle, and termination (A4)
5. Session identifiers are cryptographically random and unpredictable (A5)
6. Auditor accounts are properly authenticated and managed (A39)

### 2.2 Segregation of Duties
7. User roles accurately reflect job responsibilities (A7)
8. Role assignments are reviewed periodically (A8)
9. User provisioning is accurate and aligned with business function (A9)
10. User deprovisioning occurs promptly upon termination or role change (A10)
11. SoD rules are correctly defined for financial processes (A11)
12. SoD controls are technically enforced, not merely procedural (A12)
13. No user can both create and approve the same transaction (A13)
14. System administrators and DBAs are subject to the same SoD as the finance team (A51)

### 2.3 Approval Workflow
15. Approval workflows cannot be bypassed via direct database or backend API access (A14)
16. Approval decisions are authenticated and identity cannot be spoofed (A15, A53)
17. Approval records are immutable after creation (A16)
18. Workflow logic is correctly implemented with no edge-case flaws (A17)
19. Emergency override procedures are controlled, logged, and time-bound (A18)
20. Low-level users cannot modify session parameters to bypass approval chains or access admin functions (A55)

### 2.4 Application Security
21. Authorization is enforced at every layer including backend APIs (A6, A21)
22. ERP application code is free of exploitable vulnerabilities (A19)
23. Input validation is effective against injection attacks, including SQLi in the Reporting Engine (A20)
24. Administrative functions are restricted to authorized roles (A22)
25. Change management prevents unauthorized modifications to financial controls (A47)
26. Backend APIs independently enforce authorization — front-end controls are not sufficient (A21)

### 2.5 Database Security
27. Financial database access is restricted to authorized services and users (A24)
28. Database permissions follow least privilege (A25)
29. Database integrity protections exist (constraints, checksums, immutability) (A26)
30. Database software is kept current with security patches (A27)
31. Secrets and credentials (service accounts, API keys, connection strings) are securely stored (A28)

### 2.6 Reporting / Audit Access
32. Reporting engine faithfully reflects source data without transformation errors (A35)
33. Reporting logic is accurate and free of calculation errors (A36)
34. Data transferred from ERP to reporting maintains integrity in transit (A37)
35. Reporting engine access is controlled to authorized users (A38)
36. Read-only access for auditors is technically enforced at both application and database level (A40, A52)
37. Read-only restrictions cannot be bypassed via alternate paths (A41)
38. Auditors cannot escalate privileges from read-only to write (A42)
39. Reporting Engine cannot be exploited for SQL injection or privilege escalation (A42 overlap with A20)

### 2.7 Logging / Monitoring
40. Administrative actions are logged and auditable (A23)
41. Audit logs are complete and capture all SOX-relevant events (A31)
42. Audit logs cannot be altered — stored in immutable, append-only format (A32)
43. Audit logs capture all journal-entry activity (create, modify, approve) (A33)
44. System clocks across ERP, Approval Workflow, and Audit Logs are synchronized via NTP (A34)

### 2.8 Access Recertification
45. Quarterly recertification of user access is performed and is meaningful (A43)
46. Recertification reviews are completed on schedule (A44)
47. Recertification data must come from a tamper-proof system of record, not from reports generated by the administrators under review (A54)
48. Exceptions identified during recertification are remediated in a timely manner (A45)

### 2.9 Backup / Disaster Recovery
49. Backup and recovery processes preserve financial data integrity (A48)
50. Disaster recovery procedures maintain SOX control effectiveness during outages (A49)

### 2.10 Encryption
51. Encryption protects sensitive financial data at rest and in transit (A29)
52. Encryption keys are managed with proper lifecycle controls (rotation, access, backup) (A30)

### 2.11 Insider Threat
53. Insider personnel (finance, IT, admin) are trustworthy; compensating controls exist for high-risk roles (A46)

### 2.12 Compliance / SOX
54. No undocumented financial processing paths exist — all batch jobs, exports, and manual workflows are accounted for (A50)
55. Administrative actions and system changes are subject to change management controls (A47)

---

## 3. Risk Scores

| # | Assumption | Likelihood | Impact | Risk |
|---|-----------|:----------:|:------:|:----:|
| A1 | Finance users properly authenticated | L | H | H |
| A2 | Credentials protected | M | H | H |
| A3 | MFA exists | M | H | H |
| A4 | Session management secure | L | H | H |
| A5 | Session IDs unpredictable | L | H | H |
| A6 | Authorization correct at every layer | M | C | C |
| A7 | Roles reflect job responsibilities | L | H | H |
| A8 | Roles reviewed periodically | M | H | H |
| A9 | User provisioning accurate | M | H | H |
| A10 | User deprovisioning prompt | H | H | H |
| A11 | SoD rules correctly defined | L | C | H |
| A12 | SoD technically enforced | M | C | C |
| A13 | No user creates and approves same transaction | L | C | H |
| A14 | Approval workflows not bypassable | M | C | C |
| A15 | Approval decisions authenticated | L | C | H |
| A16 | Approval records immutable | L | C | H |
| A17 | Workflow logic correct | M | H | H |
| A18 | Emergency override controlled | M | H | H |
| A19 | ERP app code secure | M | C | C |
| A20 | Input validation effective | M | H | H |
| A21 | Backend APIs enforce authz | M | C | C |
| A22 | Admin functions restricted | M | H | H |
| A23 | Admin actions logged | L | H | H |
| A24 | DB access restricted | M | C | C |
| A25 | DB least privilege | M | C | C |
| A26 | DB integrity protections | L | C | H |
| A27 | DB software patched | M | H | H |
| A28 | Secrets protected | M | C | C |
| A29 | Encryption of financial data | L | H | H |
| A30 | Key management proper | M | C | H |
| A31 | Audit logs complete | M | C | C |
| A32 | Audit logs immutable | L | C | H |
| A33 | Logs capture journal-entry activity | M | H | H |
| A34 | NTP synchronization | L | H | H |
| A35 | Reporting engine faithful | L | H | H |
| A36 | Reporting logic accurate | L | H | H |
| A37 | Data transfer integrity | L | H | H |
| A38 | Reporting engine access controlled | M | H | H |
| A39 | Auditor accounts authenticated | L | H | H |
| A40 | Read-only enforced (app level) | M | H | H |
| A41 | Read-only not bypassable | M | C | C |
| A42 | Auditors cannot elevate privileges | M | C | C |
| A43 | Quarterly recertification meaningful | H | H | H |
| A44 | Recertification performed on schedule | M | H | H |
| A45 | Recertification exceptions remediated | H | H | H |
| A46 | Insider personnel trustworthy | L | C | H |
| A47 | Change management effective | M | H | H |
| A48 | Backup/DR preserves integrity | L | C | H |
| A49 | DR maintains SOX controls | M | H | H |
| A50 | No undocumented processing paths | H | C | C |
| A51 | DBA/Admin subject to SoD | M | C | C |
| A52 | DB-level read-only for Reporting Engine | M | H | H |
| A53 | Approval identity not spoofable | M | C | C |
| A54 | Recertification from tamper-proof system | H | H | H |
| A55 | Low-level users cannot escalate via session parameters | M | C | C |

---

## 4. STRIDE Mapping

### Spoofing
- A1: User authentication
- A2: Credential protection
- A3: MFA enforcement
- A15, A53: Approval identity authentication
- A39: Auditor authentication

### Tampering
- A12: SoD technical enforcement
- A13: Create/approve separation
- A14: Workflow bypass prevention
- A16: Approval record immutability
- A19: Application code integrity
- A24: Database access restriction
- A26: Database integrity
- A32: Audit log immutability
- A41: Read-only bypass
- A47: Change management integrity
- A48: Backup/DR integrity
- A52: Database-level read-only enforcement
- A55: Session parameter tampering

### Repudiation
- A15: Authenticated approval decisions
- A23: Administrative action logging
- A31–A33: Complete, immutable audit logs
- A34: NTP-synchronized timestamps

### Information Disclosure
- A24: Database access restriction
- A29: Encryption at rest and in transit
- A30: Key management
- A38: Reporting engine access control
- A39: Auditor authentication
- A40, A52: Read-only enforcement
- A28: Secrets protection

### Denial of Service
- A27: Database patching (availability)
- A34: NTP synchronization (availability of reliable timestamps)
- A48: Backup/DR processes
- A49: DR control maintenance

### Elevation of Privilege
- A6: Authorization at every layer
- A11, A12: SoD definition and enforcement
- A18: Emergency override controls
- A21: Backend API authorization
- A22: Admin function restriction
- A25: Database least privilege
- A28: Secrets protection (credential escalation)
- A41: Read-only bypass (privilege escalation)
- A42: Auditor privilege escalation
- A51: DBA SoD
- A53: Approval identity spoofing
- A55: Session parameter manipulation

---

## 5. Top 10 Critical Assumptions (Ranked)

### 1. DBA/Administrator SoD Enforcement (A51)
**Rationale**: DBAs possess unrestricted access to the Financial DB. Without SoD enforcement mirroring the finance team, a DBA could create, approve, and conceal financial transactions with no technical barrier. This single role concentrates the highest privilege with the lowest accountability.

### 2. Approval Workflow Bypass via Direct Database or Backend API Access (A14 / A21)
**Rationale**: If the Approval Workflow engine can be bypassed by calling backend APIs directly or writing to the Financial DB, every transaction-level control in the architecture is rendered moot. The approval chain is the primary preventive control under SOX.

### 3. Approval Identity Spoofing via Session Fixation or ID Hijacking (A53)
**Rationale**: If user identity within the approval workflow can be spoofed, a single attacker can fabricate an entire chain of approvals — creation, review, and sign-off — without detection. Self-approval of fraudulent transactions becomes trivially achievable.

### 4. Read-Only Access Enforced at Application Level Only, Not Database Level (A52)
**Rationale**: Application-level read-only enforcement is bypassable by connecting directly to the database with stolen or misconfigured credentials. Database-level permissions provide a defense-in-depth layer that survives web application compromise.

### 5. Audit Log Tampering / Immutability Failure (A32)
**Rationale**: SOX audits depend entirely on the integrity of audit logs. If logs can be modified or deleted, fraudulent activity can be systematically concealed. Without immutable, append-only storage, the entire audit framework is unreliable.

### 6. Backend API Authorization Not Independently Enforced (A21)
**Rationale**: Front-end UI controls are trivially bypassed. If backend APIs do not independently verify authorization for every request, an attacker with network access can submit transactions that bypass approval workflows, SoD checks, and logging.

### 7. No Undocumented Financial Processing Paths (A50)
**Rationale**: Architecture diagrams rarely capture every batch job, ETL process, spreadsheet macro, or manual journal entry. These undocumented paths represent blind spots where transactions can enter or leave the financial system without passing through any SOX control.

### 8. Low-Level User Escalation via Session Parameter Manipulation (A55)
**Rationale**: If session parameters (e.g., role, user ID, approval authority) are stored client-side or in tamperable server-side state, a low-level accountant can escalate to administrative functions or bypass approval chains. This is a common web application vulnerability with severe financial implications.

### 9. User Deprovisioning Failures (A10)
**Rationale**: Former employees, contractors, and role-changed users often retain access longer than policy permits. In a SOX-controlled financial system, each orphaned account represents a persistent unauthorized access path that can be exploited internally or externally.

### 10. Recertification Data from Non-Tamper-Proof Sources (A54)
**Rationale**: If recertification reviews are conducted using reports generated by the very administrators whose access is being reviewed, the data lacks independence. Administrators could omit their own privileged accounts from reports, allowing privilege creep to go undetected.

---

## 6. Recommended Controls

### 1. DBA/Administrator SoD Enforcement
- Implement break-glass DBA access with approval workflow, time-bound grants, and full audit logging
- Deploy a privileged access management (PAM) solution for all database administrative sessions
- Require dual-authorization for all schema changes and direct data modifications
- Conduct quarterly reviews of all DBA-level access with sign-off from finance and audit

### 2. Approval Workflow Bypass Prevention
- Implement API gateway with centralized authorization checks for every backend endpoint
- Database triggers or stored procedures that reject transactions not originating from the approval workflow engine
- Write-ahead logging for all transaction attempts, including rejected direct-write attempts
- Runtime application self-protection (RASP) to detect and block direct database calls

### 3. Approval Identity Spoofing Prevention
- Bind approval session to cryptographically generated, server-side session tokens with HttpOnly, Secure, and SameSite attributes
- Require step-up authentication (MFA re-prompt) for approval actions
- Implement approval signing with digital signatures or server-side HMAC
- Audit and alert on any case where approval identity differs from session identity

### 4. Database-Level Read-Only Enforcement
- Configure database roles with explicit read-only permissions for the Reporting Engine service account
- Remove all write grants (INSERT, UPDATE, DELETE, DDL) from the Reporting Engine database user
- Implement database firewall rules restricting the Reporting Engine to SELECT-only statements
- Periodically audit database-level permissions using automated scanners

### 5. Audit Log Immutability
- Deploy append-only, write-once-read-many (WORM) storage for audit logs (e.g., AWS S3 Object Lock, immutable file system)
- Centralize logs to a SIEM platform with separate administrative domain from the ERP
- Implement cryptographic log chaining (hash-chain) to detect log tampering
- Alert on any attempt to modify or delete audit log entries

### 6. Backend API Authorization
- Implement attribute-based access control (ABAC) evaluated at the API gateway
- Apply the `authorization` decorator/check to every backend controller method, not just the web UI
- Conduct automated API fuzzing and authorization testing in CI/CD pipeline
- Perform regular penetration testing targeting backend API endpoints

### 7. Undocumented Processing Path Discovery
- Conduct architecture discovery interviews with finance, IT, and operations teams
- Implement network flow monitoring to identify unapproved data paths
- Inventory all batch jobs, cron tasks, ETL processes, spreadsheet imports, and manual procedures
- Establish a policy requiring architecture review for any new financial data flow

### 8. Session Parameter Integrity
- Store all authorization attributes server-side; never accept role or user ID from client
- Implement session integrity checks (tamper-evident session tokens with server-side HMAC)
- Conduct automated scans for privilege escalation vulnerabilities
- Perform regular code review focusing on session and authorization logic

### 9. User Deprovisioning Automation
- Integrate ERP with HR identity feed (HRIS) for automatic deprovisioning on termination date
- Implement daily reconciliation between active directory/IdP and ERP user accounts
- Deploy a quarterly recertification workflow with automated removal of unconfirmed access
- Generate and report on orphaned account metrics to SOX compliance team

### 10. Tamper-Proof Recertification
- Generate recertification data directly from audit logs and database-level access controls, not from application reports
- Implement a separate, read-only reporting database for access certification
- Require recertification reviewers to be independent of the systems being reviewed
- Maintain an immutable history of all recertification actions

---

## 7. Summary Statistics

| Category | Count |
|----------|:-----:|
| **Total Assumptions** | **55** |
| **Critical Risk** | 15 |
| **High Risk** | 40 |
| **Medium Risk** | 0 |
| **Low Risk** | 0 |
| **Sources** | 3 models (GPT-4o, Gemini, Gemma) |
| **Exclusive (GPT only)** | 50 |
| **Exclusive (Gemini only)** | 1 |
| **Exclusive (Gemma only)** | 4 |

**Critical Assumptions** (15): A6, A12, A14, A19, A21, A24, A25, A28, A31, A41, A42, A50, A51, A53, A55

**Top Risk Drivers**: DBA/administrator SoD bypass, approval workflow subversion via direct API/DB access, approval identity spoofing, database-level permission gaps, and audit log tampering represent the highest concentration of critical risk in Architecture 5. These five failure modes would individually defeat the SOX control framework; their correlated failure would render all financial reporting untrustworthy.

---

*End of Security Review — Architecture 5*
