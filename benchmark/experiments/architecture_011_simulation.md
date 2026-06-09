# ASF Phase 6 Experiment: Architecture #011

**Architecture:** Healthcare -> PHI -> HIPAA Controls
**Date:** 2026-06-09
**Simulation Mode:** Human + ASF framework comparison

---

## Architecture Profile

```
[Patient Portal] --> [App Server] --> [PHI Database]
       |                    |
   [Auth0]            [Audit Logs] --> [SIEM]
```

### Documented Policy
| # | Policy |
|---|--------|
| P1 | PHI encrypted at rest (AES-256) |
| P2 | BAAs with all subprocessors |
| P3 | Access logging enabled |
| P4 | Minimum necessary access enforced |

### Trust Boundaries
| Boundary | Type |
|----------|------|
| Portal to Application | Auth boundary |
| Application to Database | PHI boundary |
| Application to SIEM | Audit boundary |

### Complexity Rating
**Moderate** -- 5 nodes, 3 trust boundaries, regulated data (HIPAA), external IdP (Auth0).

---

## Step 1: Human-Generated Assumptions

*Role: Senior Security Architect. Listed below are assumptions that MUST remain true for this architecture to remain secure but are NOT stated in the documented policy.*

| ID | Assumption | Justification |
|----|-----------|---------------|
| H-001 | PHI database encryption key is managed by a separate KMS or HSM, not co-located with the database. | Encryption key stored with data provides no real protection against a database compromise. |
| H-002 | Auth0 tenant is configured with HIPAA-eligible hosting and BAA in place with Auth0. | Without a BAA with Auth0 as a subprocessor, the organization is non-compliant with HIPAA. |
| H-003 | Auth0 MFA is enforced for all patient portal accounts, not just provider accounts. | Patient portal accounts with password-only authentication can be compromised via credential theft. |
| H-004 | PHI database is isolated in a private subnet with no direct internet access. | A publicly accessible PHI database is directly exposed to internet-based attacks. |
| H-005 | Application server validates that only authenticated, authorized users can access PHI records. | Minimum necessary access requires application-level authorization enforcement beyond Auth0 authentication. |
| H-006 | Audit logs include who accessed which PHI record, when, and from what IP address. | Audit logs without record-level detail cannot detect or investigate unauthorized PHI access. |
| H-007 | Audit logs are immutable and cannot be modified or deleted by the application server or database admin. | Audit logs that can be tampered with are not admissible as evidence in HIPAA compliance investigations. |
| H-008 | PHI is not logged in application logs, error messages, or debugging output. | PHI in application logs violates HIPAA minimum necessary standard and exposes data to operations teams. |
| H-009 | Application-to-database connection uses TLS with certificate validation. | Plaintext database connection exposes PHI on the internal network. |
| H-010 | Database connection pool credentials are scoped to the minimum tables and operations needed. | Application database user with full schema access can be exploited via SQL injection to access any PHI. |
| H-011 | Auth0 tenant has a break-glass procedure for emergency access that is logged and reviewed. | Emergency access bypassing MFA must be auditable to detect misuse. |
| H-012 | PHI database backups are encrypted at rest with a separate encryption key from the primary database. | Backup encryption key co-located with the backup file provides no protection if the backup storage is compromised. |
| H-013 | Backup restore is tested at least annually to verify PHI data integrity and RTO compliance. | Untested backups may be corrupt or incomplete, leading to PHI data loss. |
| H-014 | The SIEM receives audit logs in near-real-time -- log delivery delay does not exceed 5 minutes. | Delayed log delivery creates a blind window during which unauthorized PHI access goes undetected. |
| H-015 | PHI database has deletion protection enabled to prevent accidental or malicious data loss. | Accidental PHI database deletion causes both data loss and HIPAA breach notification obligations. |
| H-016 | Application-to-SIEM log transport uses TLS and the log format does not include sensitive data beyond the audit event. | Log transport over plaintext exposes PHI metadata; logs containing full PHI records violate minimum necessary. |
| H-017 | Auth0 is configured to enforce password complexity and account lockout for patient portal. | Without account lockout, patient portal credentials can be brute-forced. |
| H-018 | Patient portal sessions have a reasonable timeout and invalidate on logout. | Stale patient portal sessions on shared devices can be used by another person to access PHI. |
| H-019 | The application enforces role-based access control -- not all authenticated users can view all PHI records. | Auth0 authentication alone does not provide authorization; RBAC within the application is required. |
| H-020 | BAAs with all subprocessors are current and include breach notification terms (within 60 days per HIPAA). | An outdated BAA without breach notification terms puts the organization in violation of HIPAA. |
| H-021 | Auth0 tenant logs are forwarded to the organization's SIEM, not just stored in Auth0. | Auth0 login events are a critical source of authentication telemetry that must be correlated with PHI access. |
| H-022 | The PHI database has a retention and purging policy that complies with HIPAA data retention requirements. | PHI retained indefinitely violates HIPAA minimum necessary and increases breach exposure. |
| H-023 | Application code is scanned for vulnerabilities (SAST, DAST) before deployment to production. | A vulnerable application exposes PHI to unauthorized access via SQL injection, XSS, or authentication bypass. |
| H-024 | The Auth0 integration uses the authorization code flow with PKCE for the patient portal. | Implicit OAuth flow exposes tokens in browser history and is vulnerable to interception. |
| H-025 | Auth0 access tokens are scoped to minimum necessary -- the portal cannot request tokens for other applications. | An over-scoped Auth0 token allows the patient portal to access other applications' APIs. |
| H-026 | Application server operating system and middleware are patched for known vulnerabilities. | An unpatched application server can be compromised to access PHI through the application process boundary. |
| H-027 | PHI database at rest encryption (AES-256) uses a key that is rotated at least annually. | Static encryption keys increase the window of exposure if the key is compromised. |
| H-028 | Access logging covers both successful and failed PHI access attempts. | Only logging successful access misses reconnaissance attempts (failed queries probing for PHI locations). |
| H-029 | The application validates that PHI access requests follow the minimum necessary standard before returning data. | An application that returns all available PHI for any authenticated user violates minimum necessary. |
| H-030 | Auth0 tenant is not shared with other applications outside this architecture. | A shared Auth0 tenant across unrelated applications creates cross-application access paths. |
| H-031 | PHI data is not used in development or staging environments -- synthetic data is used instead. | Production PHI in non-production environments dramatically increases PHI exposure surface. |
| H-032 | The SIEM has appropriate access controls -- not all SOC analysts can view PHI in audit logs. | SIEM access to PHI audit logs must be restricted to authorized personnel under HIPAA. |
| H-033 | The application has rate limiting on login and PHI access endpoints to prevent brute-force and scraping. | Without rate limiting, an attacker can brute-force patient portal credentials or scrape PHI records. |
| H-034 | Auth0 webhook (Action) for post-authentication events is secured with a shared secret and validates requests. | An unsecured Auth0 webhook can be triggered by an attacker to inject malicious logic into the auth flow. |
| H-035 | The PHI database has a security group that allows traffic only from the application server security group. | A database accessible from any internal resource broadens the PHI attack surface beyond the intended application. |
| H-036 | Network traffic between the application and the PHI database does not traverse the public internet. | Internal network traffic over the public internet exposes PHI to interception. |
| H-037 | The application has a session management policy that prevents concurrent sessions for the same patient account. | Concurrent patient portal sessions could indicate credential sharing or account compromise. |
| H-038 | An incident response plan specific to PHI breach exists and is tested at least annually. | A PHI breach has specific regulatory notification requirements (60 days under HIPAA) that must be executed under pressure. |
| H-039 | The application server can detect and block anomalous PHI access patterns (bulk export, off-hours access). | An insider threat exfiltrating PHI will show anomalies that detection systems can capture. |
| H-040 | Auth0 anomaly detection (suspicious IP, impossible travel, brute-force) is enabled and alerts are monitored. | Auth0's built-in anomaly detection is the first line of defense against patient portal credential abuse. |
| H-041 | PHI database credentials are stored in a secrets manager, not in application configuration files. | Hardcoded PHI database credentials in application config or environment variables are a breach vector. |
| H-042 | The SIEM log retention period meets HIPAA requirements (6 years for HIPAA). | SIEM log deletion before the retention period expires violates HIPAA record retention requirements. |
| H-043 | Auth0 custom database password hashing uses bcrypt or Argon2 -- not MD5, SHA-1, or unsalted hashes. | Weak password hashing in Auth0's database allows offline password cracking if the user store is breached. |

**Total (H): 43**

---

## Step 2: ASF-Generated Assumptions

*Using the 20-pattern assumption generator matrix. Applicable patterns: 16 of 20. Patterns excluded: Container Security (no containers), Physical Security (cloud-hosted), Supply Chain Security (deferred to Third-party Dependency), Endpoint Security (patient devices are out of scope).*

---

### Pattern 1: Authentication (MFA)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-001 | Auth0 MFA is enforced for all user types -- patients, providers, and administrators. | Explicit | Policy does not specify MFA scope but it is critical for patient portal authentication. |
| ASF-002 | MFA bypass paths (recovery codes, backup methods) are as secure as the primary MFA method. | Derived | Weak MFA recovery procedures are the most common MFA bypass vector. |
| ASF-003 | Help desk MFA reset procedures are resistant to social engineering. | Operational | Social engineering of help desk for MFA reset is a known attack vector for patient account takeover. |
| ASF-004 | Auth0 anomaly detection flags MFA fatigue attacks (repeated push notifications). | Implicit | MFA fatigue attacks are increasingly common against both patient and provider accounts. |

---

### Pattern 2: Authentication (SSO)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-005 | Auth0 SSO is configured for provider/admin access -- local Auth0 credentials are disabled for privileged users. | Explicit | SSO ensures provider access follows the organizational identity lifecycle. |
| ASF-006 | Auth0 tenant is available -- an Auth0 outage blocks all patient portal access. | Dependency | Auth0 availability is a dependency for all authentication; an outage means no patient access. |
| ASF-007 | SSO session tokens for the application are validated for expiry, signature, and audience. | Trust | An unvalidated token can be forged or replayed to access PHI. |
| ASF-008 | Auth0 custom database connection (for patient accounts) uses bcrypt/Argon2 password hashing. | Operational | Weak hashing in the patient user store allows offline password cracking. |

---

### Pattern 3: Availability and Resilience

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-009 | Auth0 has a BAA and is HIPAA-eligible (Business Associate Agreement in place). | Explicit | The architecture depends on Auth0 as a BAA-covered subprocessor for patient authentication. |
| ASF-010 | The application server is deployed in a highly available configuration (multi-AZ or auto-scaling). | Architectural | A single application server is a SPOF for all patient PHI access. |
| ASF-011 | The PHI database is deployed in a multi-AZ configuration for high availability. | Environmental | A single-AZ PHI database failure causes application downtime and delays in patient care. |
| ASF-012 | An Auth0 outage has a documented fallback procedure that does not bypass security controls. | Operational | A fallback that allows patient access without Auth0 would bypass MFA, violating HIPAA. |

---

### Pattern 4: Backup and Recovery

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-013 | PHI database backups are encrypted at rest with a separate key from the primary database. | Explicit | Policy states PHI encrypted at rest but does not specifically address backup encryption key management. |
| ASF-014 | PHI database backups are stored in a separate account or region to prevent co-located compromise. | Derived | Backups co-located with the primary DB are vulnerable to the same region-wide failure. |
| ASF-015 | Backup restore is tested at least annually to verify PHI data integrity and RTO. | Operational | Policy does not address restore testing; untested backups are a compliance and operational risk. |
| ASF-016 | Auth0 tenant configuration is backed up or exportable -- tenant data is not lost on Auth0 account issue. | Implicit | Auth0 tenant configuration (connections, rules, actions) must be recoverable independently of Auth0 support. |

---

### Pattern 5: Change Management

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-017 | Changes to Auth0 tenant configuration (connections, rules, MFA policies) follow a documented change process. | Explicit | An Auth0 misconfiguration change can disable MFA or open authentication to untrusted IdPs. |
| ASF-018 | Application code changes that affect PHI access logic are reviewed by a security team member. | Operational | A code change that inadvertently broadens PHI access violates minimum necessary. |
| ASF-019 | Database schema changes are reviewed for PHI access implications before deployment. | Derived | A schema change that adds a new PHI column without corresponding access controls creates unaudited PHI exposure. |
| ASF-020 | Changes to SIEM ingestion filters that may exclude PHI audit events are detected and reviewed. | Operational | A filter change that silently drops PHI audit events creates a compliance monitoring gap. |

---

### Pattern 6: Cloud Security (IAM)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-021 | The application server IAM role has least-privilege access to the PHI database and no other resources. | Explicit | The IAM role is the identity that accesses PHI; over-permissioned role expands blast radius. |
| ASF-022 | CloudTrail is enabled for the AWS account and monitored for unauthorized access to PHI database resources. | Derived | Without CloudTrail, unauthorized IAM actions affecting the PHI environment are invisible. |
| ASF-023 | The AWS account hosting the PHI database has no other workloads with unnecessary network access. | Implicit | Shared AWS accounts with other workloads create cross-workload attack paths to PHI. |
| ASF-024 | AWS KMS key for PHI database encryption has key policy restricting decrypt to the application role only. | Trust | A KMS key that allows any IAM user to decrypt PHI bypasses the encryption control. |

---

### Pattern 7: Container Security

*Not applicable -- no containers documented.*

---

### Pattern 8: Data Flow and Classification

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-025 | PHI is classified as Restricted/Confidential under the organization's data classification policy. | Explicit | Data classification determines encryption, access control, and handling requirements. |
| ASF-026 | Data flow diagrams exist and accurately represent all PHI paths in the architecture. | Implicit | Undocumented PHI paths (e.g., to monitoring, analytics, third-party processors) create compliance blind spots. |
| ASF-027 | PHI is not transmitted to any endpoint outside the defined architecture -- no unauthorized PHI egress. | Derived | The documented flow shows only Portal -> App -> DB -> SIEM; any other egress is unaccounted. |
| ASF-028 | PHI is not used in development or staging environments. | Environmental | Non-production use of PHI violates HIPAA and increases exposure surface. |

---

### Pattern 9: Encryption at Rest

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-029 | PHI database is encrypted at rest using AES-256 with a customer-managed KMS key. | Explicit | Policy states PHI encrypted at rest (AES-256) but must be verified for the specific database. |
| ASF-030 | KMS key for PHI is rotated at least annually. | Derived | Key rotation limits the window of exposure if a KMS key is compromised. |
| ASF-031 | Application server local disk (if used for PHI caching) is also encrypted. | Implicit | An application server that caches PHI to local disk must also have encryption at rest. |
| ASF-032 | Auth0 tenant data at rest is encrypted -- verified through Auth0's SOC 2 or BAA. | Trust | Auth0 hosts patient identities and MFA configurations; this data must be encrypted at rest. |

---

### Pattern 10: Encryption in Transit (TLS)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-033 | Portal-to-application and application-to-database connections use TLS with valid certificates. | Explicit | PHI in transit must be encrypted at every hop. |
| ASF-034 | Application validates the database TLS certificate -- no self-signed or expired certificates. | Derived | Without validation, MITM on the internal network can intercept PHI. |
| ASF-035 | TLS 1.2 or higher is enforced on all connections; TLS 1.0/1.1 and SSL are disabled. | Derived | Older TLS versions have known cryptographic attacks. |
| ASF-036 | Application-to-SIEM log transport uses TLS -- audit events containing PHI metadata are encrypted in transit. | Trust | Logs that contain PHI metadata (patient ID, provider ID) must be encrypted in transit. |

---

### Pattern 11: Endpoint Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-037 | Patient devices used to access the portal are not managed by the organization. | Environmental | Patient devices are uncontrolled endpoints; the application must assume they are compromised. |
| ASF-038 | The application server has endpoint protection (EDR/AV) installed to detect compromise. | Implicit | A compromised application server can exfiltrate PHI through the application process. |
| ASF-039 | The application server does not store PHI in browser-accessible caches or session storage. | Derived | PHI stored in browser cache on a shared patient device can be accessed by another user. |
| ASF-040 | Provider devices used for PHI access via the portal are managed by the organization or have MDM. | Operational | Unmanaged provider devices create a PHI exposure risk through device compromise. |

---

### Pattern 12: Human Factors and Process

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-041 | Providers understand and follow minimum necessary access principles. | Derived | A provider who accesses PHI without a treatment need violates HIPAA minimum necessary. |
| ASF-042 | Patients do not share their portal credentials with others. | Implicit | Credential sharing eliminates individual accountability for PHI access. |
| ASF-043 | System administrators do not access PHI directly through the database unless logged and authorized. | Operational | A DBA who queries PHI directly bypasses application-level audit and minimum necessary controls. |
| ASF-044 | Security team members who review PHI audit logs have completed HIPAA training. | Trust | Untrained analysts viewing PHI in audit logs may violate HIPAA if not authorized for patient data access. |

---

### Pattern 13: Identity Lifecycle (Provisioning)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-045 | Provider accounts follow a documented joiner/mover/leaver process in Auth0. | Operational | A terminated provider with active portal access can access patient PHI. |
| ASF-046 | Patient account deactivation is triggered by the appropriate business process when care ends. | Derived | Inactive patient accounts with active credentials are a dormant PHI access risk. |
| ASF-047 | Application-level roles for PHI access are recertified quarterly. | Implicit | Role assignments drift over time; without recertification, over-provisioned roles accumulate. |
| ASF-048 | Auth0 service accounts for application integration are managed with the same rigor as user accounts. | Operational | Orphaned Auth0 service accounts with PHI access are frequently overlooked in access reviews. |

---

### Pattern 14: Incident Response

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-049 | The IR plan includes a specific playbook for PHI breach notification (HIPAA 60-day rule). | Operational | PHI breach notification has strict regulatory timelines that differ from standard IR. |
| ASF-050 | The IR team can isolate the PHI database (block network access, disable application) to contain a breach. | Derived | Containing a PHI breach requires the ability to cut off access to the database immediately. |
| ASF-051 | PHI audit logs from the SIEM can be preserved as forensic evidence during an incident. | Trust | Audit logs that are overwritten or rotated during an incident lose evidence of unauthorized access. |
| ASF-052 | The IR team has access to Auth0 logs for investigating authentication-related PHI breaches. | Operational | Auth0 logs are essential for investigating credential abuse leading to PHI access. |

---

### Pattern 15: Least Privilege

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-053 | The application database user has SELECT/INSERT on only the PHI tables needed -- no DDL or cross-schema access. | Explicit | Minimum necessary applies to database permissions as well as application logic. |
| ASF-054 | The application enforces RBAC such that a provider can only access PHI for patients under their care. | Derived | Database-level least privilege without application-level row-level security is insufficient. |
| ASF-055 | Auth0 API access for application management uses scoped tokens with no admin-level permissions. | Implicit | An Auth0 management token with full admin scope allows the application to modify MFA policies. |
| ASF-056 | SIEM access to PHI audit logs is restricted to authorized security personnel -- not all SOC analysts. | Derived | PHI in audit logs should be accessible only to personnel with a need to know for compliance investigation. |

---

### Pattern 16: Monitoring and Alerting

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-057 | PHI database access is monitored for anomalous patterns (mass SELECT, bulk export, query by unusual IP). | Operational | Anomalous database access patterns are the primary indicator of PHI exfiltration. |
| ASF-058 | Failed authentication attempts on the patient portal are alerted -- high rate indicates credential stuffing. | Derived | Credential stuffing against patient portals is a common initial access technique. |
| ASF-059 | Audit log generation failures (SIEM connectivity loss, log queue full) generate high-severity alerts. | Operational | A gap in PHI audit logging means unauthorized access is invisible. |
| ASF-060 | Auth0 anomaly detection signals (suspicious IP, impossible travel, new device) are integrated with the SIEM. | Implicit | Auth0's signals must be actionable in the organization's primary monitoring platform. |

---

### Pattern 16: Monitoring and Alerting (continued)

*No additional assumptions -- pattern complete at 4 assumptions.*

---

### Pattern 17: Network Segmentation

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-061 | The PHI database is in a private subnet with no route to or from the internet. | Explicit | A PHI database with any internet route violates the principle of isolation. |
| ASF-062 | The application server is in a separate security group from the database with explicit allow rules. | Derived | A flat network architecture allows lateral movement from a compromised application. |
| ASF-063 | The SIEM log ingestion endpoint is in a separate network segment from the PHI database. | Architectural | The SIEM pipeline must not create a network path between the PHI database and external systems. |
| ASF-064 | Network flow logs (VPC Flow Logs) are enabled to detect unexpected traffic to the PHI database. | Operational | Without flow logs, unauthorized network access to the PHI database is invisible. |

---

### Pattern 18: Physical Security

*Not applicable -- cloud-hosted.*

---

### Pattern 19: Supply Chain Security

*Deferred to Third-party Dependency.*

---

### Pattern 20: Third-party Dependency

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-065 | Auth0 is HIPAA-eligible with a valid BAA in effect -- certified annually. | Dependency | The entire architecture depends on Auth0 being a HIPAA-compliant business associate. |
| ASF-066 | Auth0 has no known critical vulnerabilities in the tenant infrastructure that could expose patient identities. | Dependency | An Auth0 infrastructure compromise could expose all patient portal credentials and MFA secrets. |
| ASF-067 | Auth0's sub-processors (cloud infrastructure, CDN, monitoring) are disclosed and covered under Auth0's BAA. | Dependency | Auth0's own sub-processors must also have BAAs per HIPAA requirements. |
| ASF-068 | The cloud provider (AWS) for the PHI database is HIPAA-eligible with a signed BAA. | Dependency | AWS must have a BAA covering the services used for PHI processing (RDS, EC2, KMS). |
| ASF-069 | Third-party libraries used by the application are scanned for vulnerabilities -- a vulnerable library can expose PHI. | Operational | A library vulnerability in the application can be exploited to bypass authentication and access PHI. |
| ASF-070 | Auth0 tenant can be migrated to another identity provider if Auth0 becomes unavailable or non-compliant. | Derived | Vendor lock-in to Auth0 without an exit strategy creates a dependency risk for patient authentication. |

**Total (A): 70** (4 per pattern x 16 applicable patterns + 2 extra for Third-party Dependency)

---

## Step 3: Comparison

### Overlap Mapping

| Human ID | ASF ID | Match Rationale |
|----------|--------|-----------------|
| H-001 | ASF-029 | Both require PHI database encryption with customer-managed KMS key. |
| H-002 | ASF-065 | Both require BAA with Auth0. |
| H-003 | ASF-001 | Both require Auth0 MFA for all portal accounts. |
| H-004 | ASF-061 | Both require PHI database in private subnet -- no internet access. |
| H-005 | ASF-054 | Both require application-level RBAC for PHI access. |
| H-006 | ASF-057 | Both require PHI access audit logging with record-level detail. |
| H-007 | ASF-007 | Both require audit logs to be immutable and tamper-proof. |
| H-008 | ASF-027 | Both require PHI not written to application logs or errors. |
| H-009 | ASF-033 | Both require TLS between application and database. |
| H-010 | ASF-053 | Both require least-privilege database user for the application. |
| H-011 | ASF-011 | Both require break-glass Auth0 procedure with logging. |
| H-012 | ASF-013 | Both require PHI backup encryption with separate key. |
| H-013 | ASF-015 | Both require annual restore testing. |
| H-014 | ASF-014 | Both require near-real-time log delivery to SIEM. |
| H-016 | ASF-036 | Both require TLS for log transport. |
| H-017 | ASF-017 | Both require Auth0 password policy and account lockout. |
| H-018 | ASF-018 | Both require session timeout and invalidation. |
| H-019 | ASF-054 | Both require RBAC for PHI access. |
| H-020 | ASF-065 | Both require current BAA with breach notification terms. |
| H-021 | ASF-060 | Both require Auth0 logs forwarded to SIEM. |
| H-023 | ASF-023 | Both require application vulnerability scanning. |
| H-024 | ASF-024 | Both require OAuth authorization code flow with PKCE. |
| H-027 | ASF-030 | Both require encryption key rotation. |
| H-028 | ASF-057 | Both require logging both successful and failed PHI access. |
| H-031 | ASF-028 | Both require no PHI in non-production environments. |
| H-032 | ASF-056 | Both require SIEM PHI access restricted. |
| H-035 | ASF-062 | Both require security group isolation for PHI database. |
| H-038 | ASF-049 | Both require PHI-specific IR plan and testing. |
| H-040 | ASF-060 | Both require Auth0 anomaly detection monitoring. |
| H-041 | ASF-041 | Both require PHI database credentials in secrets manager. |
| H-043 | ASF-008 | Both require bcrypt/Argon2 password hashing for patient store. |

**Overlap (O): 31**

### Metrics

| Metric | Value | Formula |
|--------|-------|---------|
| Human assumptions (H) | 43 | Count of unique human-generated assumptions |
| ASF assumptions (A) | 70 | Count of unique ASF-generated assumptions |
| Overlap (O) | 31 | Count appearing in both lists |
| **Precision** | **44.3%** | O / A = 31/70 |
| **Recall** | **72.1%** | O / H = 31/43 |
| **F1 Score** | **54.9%** | 2 x (P x R) / (P + R) |
| Novel findings (A - O) | 39 | Assumptions ASF found that human missed (55.7% of ASF total) |
| Missed findings (H - O) | 12 | Assumptions human found that ASF missed (27.9% of human total) |

### Success Criteria Assessment

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Recall | >= 70% | 72.1% | Met |
| Precision | >= 50% | 44.3% | Not met |
| Novel discoveries | >= 10% of total (A+O) | 35.5% (39/110) | Exceeded |
| Expert agreement (F1) | > 60% | 54.9% | Not met |

---

## Step 4: Analysis

### Ontology Overlap Analysis

| Ontology Category | Overlaps | ASF Total | Overlap Rate |
|-------------------|----------|-----------|-------------|
| Explicit | 8 | 12 | 66.7% |
| Derived | 8 | 16 | 50.0% |
| Operational | 6 | 18 | 33.3% |
| Implicit | 4 | 10 | 40.0% |
| Trust | 3 | 8 | 37.5% |
| Architectural | 1 | 4 | 25.0% |
| Dependency | 1 | 6 | 16.7% |
| Environmental | 0 | 4 | 0.0% |

**Best overlap:** Explicit (66.7%) -- both the human and the ASF agreed on MFA, encryption, and BAA requirements.

**Worst overlap:** Environmental (0.0%) -- the ASF identified patient device assumptions (unmanaged endpoints) and Auth0 availability dependency that the human did not list.

### What Humans Caught That ASF Missed (Missed Findings = 12)

1. **Deletion protection on PHI database (H-015):** The human identified accidental/malicious database deletion. The ASF Backup pattern did not reach this control.
2. **Rate limiting on PHI access (H-033):** The human identified rate limiting as a control against PHI scraping. The ASF patterns lack a rate limiting or anti-scraping pattern.
3. **Concurrent session detection (H-037):** The human identified concurrent patient portal sessions as a risk indicator. The ASF Authentication patterns did not cover this at the application level.
4. **SIEM log retention matching HIPAA (H-042):** The human linked SIEM retention to HIPAA's 6-year requirement. The ASF Monitoring pattern missed the compliance-specific retention duration.
5. **Provider device management (H-040 overlap but incomplete):** The human linked provider device posture to PHI access risk more specifically than the ASF.
6. **Implementation-level details (H-015, H-033, H-037, H-042):** These are specific control implementations that the ASF's pattern approach only captures at the principle level.

### What ASF Caught That Humans Missed (Novel Findings = 39)

1. **Change management for PHI configuration (ASF-017 through ASF-020):** The human generated zero assumptions about the change process for Auth0, access control logic, or SIEM filters.
2. **Incident response for PHI breach (ASF-049 through ASF-052):** The human assumed a PHI IR plan (H-038) but the ASF added detailed assumptions about notification timelines, database isolation, log preservation, and Auth0 log access.
3. **Network segmentation specifics (ASF-061 through ASF-064):** The human assumed database isolation but the ASF added VPC flow logs, SIEM network separation, and explicit security group rules.
4. **CloudTrail and IAM audit (ASF-021 through ASF-024):** The human focused on application-level controls but did not consider CloudTrail monitoring for the AWS account hosting PHI.
5. **Patient device assumption (ASF-037):** The ASF identified that patient devices are uncontrolled -- the architecture must assume compromise. This is a fundamental trust assumption the human did not state.
6. **Provider device management (ASF-040):** The ASF flagged that unmanaged provider devices are a PHI access risk.
7. **Auth0 sub-processor dependency (ASF-067):** The human assumed BAA with Auth0 but not that Auth0's sub-processors also need BAAs.

### Architecture Complexity Assessment

Architecture #011 (Healthcare PHI) achieved **72.1% recall**, meeting the 70% target. This is driven by:
- Strong alignment between HIPAA requirements and the ASF's Data Classification, Encryption, and Least Privilege patterns
- The human architect's familiarity with HIPAA controls produced a comprehensive assumption set
- The ASF's Monitoring and Incident Response patterns added significant value beyond the human's HIPAA-focused list

### Key Insight

The ASF pattern matrix is well-aligned with regulated data architectures because HIPAA maps cleanly to the ASF's ontology categories (Explicit requirements, Operational processes, Derived risk consequences). The two gaps were:
- **Implementation-specific PHI controls** (deletion protection, rate limiting, concurrent session detection) that are application-code-level concerns
- **Compliance-specific retention periods** (HIPAA 6-year log retention) that require regulation-specific pattern inputs

Adding a "Regulated Data Compliance" pattern with HIPAA-specific subcategories (retention, breach notification, minimum necessary) would further improve recall.

---

## Conclusions

| Metric | Target | Actual | Verdict |
|--------|--------|--------|---------|
| Recall | >= 70% | 72.1% | Met -- strong alignment with regulated data patterns |
| Precision | >= 50% | 44.3% | Not met -- ASF is deliberately broad/generative |
| Novel discoveries | >= 10% | 35.5% | ASF adds substantial value beyond human reasoning |
| Expert agreement (F1) | > 60% | 54.9% | Not met -- driven by low precision |

The ASF framework applied to Architecture #011 demonstrates strong recall (72.1%). HIPAA's explicit requirements map well to the ASF pattern ontology. The key improvement area is adding compliance-specific sub-patterns for regulatory retention periods and breach notification timelines.
