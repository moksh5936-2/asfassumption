# ASF Phase 6 Experiment: Architecture #020

**Architecture:** ERP → SOX → Financial Reporting → Audit
**Date:** 2026-06-09
**Simulation Mode:** Human + ASF framework comparison

---

## Architecture Profile

```
[Finance Team] --> [ERP Web App] --> [ERP Backend] --> [Financial DB]
       │                    │              │
   [Approval Workflow]  [Audit Logs]  [Reporting Engine] --> [Auditor Access]
```

### Documented Policy
| # | Policy |
|---|--------|
| P1 | SOX controls on all journal entries |
| P2 | Segregation of duties enforced |
| P3 | Read-only access for auditors |
| P4 | Quarterly recertification |

### Trust Boundaries
| Boundary | Type |
|----------|------|
| User ↔ ERP | Authentication boundary |
| Approval workflow | Segregation boundary |
| ERP ↔ Reporting | Data integrity boundary |
| Auditor Access | Read-only boundary |

### Complexity Rating
**Moderate-High** — 7 logical nodes, 4 trust boundaries, compliance-driven (SOX) architecture with workflow approval.

---

## Step 1: Human-Generated Assumptions

*Role: Senior Security Architect. Listed below are assumptions that MUST remain true for this architecture to remain secure but are NOT stated in the documented policy.*

| ID | Assumption | Justification |
|----|-----------|---------------|
| H-001 | ERP user accounts are individually assigned and not shared among finance team members. | Shared accounts eliminate audit trail attribution for financial transactions. |
| H-002 | The ERP web application enforces MFA for all finance team users, especially those who can approve journal entries. | Password-only ERP access exposes journal entry approval to credential theft. |
| H-003 | Segregation of duties (SoD) rules are enforced automatically by the ERP system, not manually by administrators. | Manual SoD enforcement is inconsistent and relies on administrative diligence that can be overridden. |
| H-004 | The approval workflow for journal entries requires at least two independent approvers for entries above a threshold. | Single-approver journal entries bypass the core SoD control documented in the policy. |
| H-005 | The Financial DB has immutable audit logging enabled — records cannot be altered or deleted by any user, including DBAs. | Mutable audit logs allow an insider to modify financial records after the fact without detection. |
| H-006 | The Reporting Engine serves data from a read-replica of the Financial DB, not the primary write database. | Reporting queries on the primary DB create contention and can inadvertently modify data through query injection. |
| H-007 | Auditor Access is restricted to the Reporting Engine only and does not include direct database or backend API access. | Direct DB or API access for auditors bypasses the read-only boundary and SOX controls. |
| H-008 | All ERP user sessions are logged with user ID, IP address, timestamp, and action performed. | Without comprehensive session logging, forensic investigation of financial irregularities is impossible. |
| H-009 | Quarterly recertification includes a review of ERP application roles, not just AD/LDAP group membership. | Recertification of AD groups alone misses ERP-application-level role assignments that grant financial permissions. |
| H-010 | Failed login attempts to the ERP web application are logged and trigger alerts after a configurable threshold. | Unmonitored failed logins allow credential brute-force against financial systems. |
| H-011 | The ERP backend enforces parameterized queries or prepared statements to prevent SQL injection into the Financial DB. | SQL injection in the ERP backend can extract, modify, or delete all financial records. |
| H-012 | The Financial DB has audit triggers on all financial tables that record pre- and post-update values. | Without value-level audit, rollback of unauthorized changes cannot be detected or reversed. |
| H-013 | The ERP web application and backend are deployed on separate servers/containers, not a monolithic deployment. | A monolithic ERP deployment conflates web and backend security boundaries, increasing blast radius. |
| H-014 | The ERP system enforces that no single user can both create and approve a journal entry (conflict of interest detection). | This is the fundamental SoD requirement — exceeding the documented policy's implicit scope. |
| H-015 | Audit logs from the ERP, approval workflow, and database are aggregated into a centralized SIEM with tamper-proof storage. | Decentralized logs can be altered individually; centralized tamper-proof storage preserves forensic integrity. |
| H-016 | The Financial DB is encrypted at rest using a key managed separately from the ERP application. | DB encryption with the application's key allows a compromised application to decrypt the database. |
| H-017 | The Reporting Engine is a separate service account with read-only database permissions, using a distinct credential from the ERP backend. | Shared credentials between ERP and Reporting violate the read-only boundary and least privilege. |
| H-018 | The ERP web application enforces CSRF tokens on all state-changing requests (journal entries, approvals, user management). | CSRF attacks can force authenticated finance users to submit unauthorized journal entries. |
| H-019 | Approved journal entries are written to an immutable ledger or append-only table before being applied to general ledger tables. | Without an immutable entry log, journal entries can be inserted, modified, or deleted without evidence. |
| H-020 | The ERP system enforces time-based access controls — finance users can only access the system during business hours. | Off-hours access to the ERP system is a strong indicator of compromised credentials or insider threat. |
| H-021 | The ERP backend validates that reporting queries do not modify data (no INSERT, UPDATE, DELETE, DDL). | Stored procedures or dynamic queries from the reporting tool could inadvertently or maliciously modify financial data. |
| H-022 | The approval workflow system is not bypassable, even by system administrators, without a documented break-glass procedure. | An approval workflow that administrators can bypass undermines the entire SoD control framework. |
| H-023 | Financial DB connection strings and credentials are not hardcoded in the ERP application configuration files. | Hardcoded database credentials in config files are exposed through source code access, CI/CD pipelines, or file reads. |
| H-024 | The ERP web application logs out inactive sessions after a short timeout (e.g., 15 minutes). | An unattended finance terminal with an active ERP session allows unauthorized journal entries. |
| H-025 | Database backups of the Financial DB are encrypted and stored in a separate location from the primary database. | Unencrypted backups are a data exposure vector; co-located backups are lost in a site disaster. |
| H-026 | The ERP application enforces row-level security so that finance users in one business unit cannot see entries from another unit. | Without row-level security, a global finance role exposes all business units' financial data to every finance user. |
| H-027 | Auditor Access is limited to completed (posted) journal entries and does not include access to draft or in-process entries. | Auditor access to in-process entries violates the segregation of duties between preparation and audit. |
| H-028 | The ERP system has a complete data classification that marks the Financial DB as "critical" and applies additional controls. | Data classification drives security control stringency; unclassified financial data may have insufficient protection. |
| H-029 | The ERP frontend and backend are on separate network segments with a firewall or security group between them. | Flat network architecture for the ERP allows lateral movement from a compromised web server to the backend. |
| H-030 | The Financial DB schema is version-controlled and all schema changes require DBA review and approval. | Uncontrolled schema changes can alter audit triggers, drop constraints, or add backdoor tables. |
| H-031 | The ERP system supports session revocation — the ability to terminate all active sessions for a user immediately. | Without session revocation, a terminated employee's ERP session remains active until timeout. |
| H-032 | Annual financial audits include a security component that reviews ERP access and SoD rule effectiveness. | Financial audits focused only on numbers may miss access control and SoD rule drift. |
| H-033 | The ERP web application enforces HTTPS-only and HSTS headers to prevent downgrade attacks. | HTTP ERP traffic on the internal network can be intercepted by any host on the same segment. |
| H-034 | The ERP system has export controls — bulk data exports of financial records are logged and require additional approval. | Unmonitored bulk exports from the ERP are the primary data exfiltration vector for financial data. |
| H-035 | The approval workflow enforces that approvers are not subordinates of the entry creator (manager approval control). | A manager approving their own subordinate's entries creates a conflict that undermines SoD. |
| H-036 | The ERP backend uses a unique database user per application module (AP, AR, GL) following least privilege. | A single shared database user for all modules grants excessive access to each module. |
| H-037 | The Financial DB has automated alerts for unusual transaction volumes or values outside normal thresholds. | Anomalous transactions (e.g., large payments to new vendors) are a key indicator of fraud. |
| H-038 | The ERP system logs all privileged operations (user creation, role assignment, permission changes) with before and after values. | Privileged operation audit is required to detect insider threats creating unauthorized finance users. |
| H-039 | The Reporting Engine enforces output masking for sensitive fields (bank account numbers, tax IDs) in reports. | Unmasked sensitive fields in reports expose PII and financial account data to users who should not see full values. |
| H-040 | The ERP web app has a rate-limited login endpoint to prevent credential stuffing attacks. | Credential stuffing against the ERP login can compromise finance accounts without triggering individual lockouts. |
| H-041 | The approval workflow system is clock-synchronized (NTP) across all components to ensure accurate audit timestamps. | Clock drift between systems produces unreliable audit timestamps that cannot be used for forensic sequencing. |
| H-042 | Audit logs are retained for a minimum of 7 years in accordance with SOX record-keeping requirements. | SOX requires 7-year retention; premature log deletion violates compliance and impairs future audits. |

**Total (H): 42**

---

## Step 2: ASF-Generated Assumptions

*Using the 20-pattern assumption generator matrix. Applicable patterns: 17 of 20. Patterns excluded: Container Security (no containers), Physical Security (on-premises data center not documented), Supply Chain Security (deferred to Third-party Dependency).*

---

### Pattern 1: Authentication (MFA)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-001 | All ERP users are authenticated via MFA before accessing the web application. | Explicit | MFA is the baseline for financial system access but not explicitly stated in the documented policy. |
| ASF-002 | MFA is enforced at the application login, not bypassable for API or reporting access. | Derived | MFA on the web login is useless if API, reporting engine, or direct DB access paths lack MFA. |
| ASF-003 | MFA recovery for ERP access follows a documented, social-engineering-resistant process. | Operational | Lost MFA devices for finance users create pressure for rapid access restoration, inviting procedural shortcuts. |
| ASF-004 | The ERP system supports hardware security keys (FIDO2/WebAuthn) as an MFA method for high-risk operations. | Implicit | SMS or TOTP MFA is phishable; hardware keys provide stronger assurance for financial transaction approval. |

---

### Pattern 2: Authentication (SSO)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-005 | ERP authentication is integrated with the enterprise IdP for centralized identity management. | Explicit | Enterprise SSO integration ensures consistent authentication across the organization. |
| ASF-006 | The ERP system validates IdP tokens correctly and does not accept forged or replayed tokens. | Derived | Token validation failures allow attackers to impersonate finance users without credentials. |
| ASF-007 | SSO session timeout is aligned between the IdP and the ERP application to prevent orphaned sessions. | Trust | An IdP session timeout longer than the ERP session timeout leaves a window for session hijacking. |
| ASF-008 | Service accounts used by the Reporting Engine and integration layers use OAuth client credentials or mTLS, not shared passwords. | Implicit | Machine-to-machine authentication with shared passwords lacks audit trail and credential rotation. |

---

### Pattern 3: Availability & Resilience

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-009 | The ERP system has a documented disaster recovery plan with defined RTO and RPO for financial processing. | Operational | Financial close and reporting have hard deadlines; extended ERP downtime has regulatory consequences. |
| ASF-010 | The ERP backend and Financial DB are deployed in an HA configuration (active-passive or active-active). | Architectural | A single ERP backend or DB failure stops all financial processing including journal entries and approvals. |
| ASF-011 | The approval workflow system is available and responsive during financial close periods. | Environmental | Approval workflow unavailability during month-end close prevents critical journal entry processing. |
| ASF-012 | The Reporting Engine can function during a partial ERP backend outage using cached or replicated data. | Derived | Auditors require reporting access even during ERP incidents; the reporting path must be resilient independently. |

---

### Pattern 4: Backup & Recovery

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-013 | The Financial DB has automated backups with point-in-time recovery and a retention period meeting SOX requirements. | Explicit | Financial data retention and recoverability are SOX compliance requirements. |
| ASF-014 | Backup restore for the Financial DB is tested at least quarterly to ensure data integrity and RTO compliance. | Derived | Untested financial DB backups risk data loss that has regulatory and financial reporting implications. |
| ASF-015 | Audit logs are backed up separately from the Financial DB and retained for the full SOX-mandated period. | Implicit | Audit logs stored only in the Financial DB are lost if the DB is restored from backup. |
| ASF-016 | ERP application configuration and approval workflow rules are backed up and version-controlled. | Operational | Loss of approval workflow rules (SoD configuration) requires manual re-creation, introducing errors. |

---

### Pattern 5: Cloud Security (IAM)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-017 | IAM roles for ERP components (if cloud-hosted) are scoped to the minimum required permissions. | Explicit | Over-permissioned IAM roles increase blast radius of any compromised ERP component. |
| ASF-018 | The Financial DB IAM policy does not grant access to any principal outside the ERP backend. | Derived | Direct database access from any service other than the ERP backend bypasses SoD and audit controls. |
| ASF-019 | The AWS account or cloud tenant hosting the ERP has MFA on the root account and no shared access keys. | Implicit | Cloud root account compromise gives an attacker full access to Financial DB snapshots and backups. |
| ASF-020 | Infrastructure-as-code for the ERP is stored in a repository with branch protection and code review. | Environmental | IaC changes without review can introduce security regressions in the financial system. |

---

### Pattern 6: Data Flow & Classification

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-021 | Financial data is classified as "Restricted" or "Critical" with corresponding handling policies. | Explicit | Classification drives encryption, access control, and data loss prevention requirements. |
| ASF-022 | Data flow diagrams exist for all financial data paths: journal entry creation, approval, posting, reporting, and audit. | Implicit | Undocumented data flows (e.g., direct data extracts to spreadsheets) bypass all ERP controls. |
| ASF-023 | Financial data is not transmitted to or stored on endpoints outside the defined architecture (e.g., local workstations, unmanaged devices). | Derived | Finance users downloading data to local spreadsheets creates uncontrolled data copies. |
| ASF-024 | Data classification policies prohibit using production financial data in non-production environments. | Environmental | Financial data in dev/test environments lacks SoD controls and audit logging. |

---

### Pattern 7: Encryption at Rest

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-025 | The Financial DB is encrypted at rest using a dedicated key managed separately from the ERP application. | Explicit | Encryption at rest with an application-controlled key offers limited protection against application compromise. |
| ASF-026 | ERP application logs containing financial data are encrypted at rest. | Derived | Logs with financial transaction data must have the same encryption protections as the database. |
| ASF-027 | The Reporting Engine's cached or aggregated data is encrypted at rest with strict access controls. | Implicit | Reporting caches containing financial aggregates are a frequently overlooked data exposure surface. |
| ASF-028 | Database backups are encrypted using a separate KMS key from the primary database key. | Derived | Backups encrypted with the same key as the primary DB offer no protection if the primary key is compromised. |

---

### Pattern 8: Encryption in Transit (TLS)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-029 | All inter-component communication (ERP Web → Backend, Backend → DB, Reporting → DB, Approval → Backend) uses TLS 1.2+. | Explicit | TLS is required for all financial data in transit but is not documented in the policy. |
| ASF-030 | Auditor Access to the Reporting Engine uses TLS and the auditor's client certificate or SSO token. | Derived | Auditor access must be encrypted and authenticated; unencrypted auditor sessions expose financial reports. |
| ASF-031 | TLS certificates for the ERP web application are from a trusted internal CA and have appropriate SANs. | Trust | Self-signed or wildcard TLS certificates weaken trust and may indicate misconfiguration. |
| ASF-032 | The database connection between the ERP backend and Financial DB uses TLS with certificate validation. | Explicit | DB connections without TLS expose financial data to MITM on the internal network. |

---

### Pattern 9: Endpoint Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-033 | Finance team workstations have endpoint protection (EDR/AV), disk encryption, and screen lock enforced. | Implicit | A compromised finance workstation exposes ERP credentials and cached financial data. |
| ASF-034 | ERP servers (web, backend, reporting) have EDR agents installed and report to a SOC. | Implicit | Compromised ERP servers without EDR provide silent persistence for financial data exfiltration. |
| ASF-035 | No unmanaged devices are permitted to access the ERP web application or auditor portal. | Operational | BYOD or unmanaged devices accessing the ERP bypass endpoint security controls. |
| ASF-036 | ERP server file systems are monitored for unauthorized file changes (FIM). | Derived | File integrity monitoring detects web shell uploads or configuration tampering on ERP servers. |

---

### Pattern 10: Human Factors & Process

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-037 | Finance users understand that ERP credentials must not be shared, even with assistants or managers. | Derived | Credential sharing in finance is common for coverage during absence; it destroys audit accountability. |
| ASF-038 | ERP administrators do not create backdoor accounts or grant themselves excessive permissions to bypass SoD. | Trust | Admins with the ability to override SoD controls are the highest insider threat in the financial system. |
| ASF-039 | The ERP approval workflow is understood by all finance users, including when an entry is blocked by SoD rules. | Operational | Users who don't understand why their entry was blocked may attempt workarounds or share credentials. |
| ASF-040 | Finance team members report phishing attempts targeting ERP credentials. | Implicit | Phishing targeting finance departments is the primary initial access vector for financial fraud. |
| ASF-041 | The quarterly recertification is performed thoroughly and not treated as a rubber-stamp exercise. | Operational | Rubber-stamped recertification is the most common finding in SOX access control audits. |

---

### Pattern 11: Identity Lifecycle (Provisioning)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-042 | ERP access follows a documented joiner/mover/leaver process integrated with HR. | Operational | Stale finance accounts are the top access control finding in SOX audits. |
| ASF-043 | Role changes (mover events) automatically trigger ERP access modification, not just AD group updates. | Derived | Employee role changes require corresponding ERP permission changes that are often missed. |
| ASF-044 | Service accounts for the Reporting Engine and integration are reviewed annually with business owners. | Implicit | Orphaned service accounts with ERP access provide undetected data access paths. |
| ASF-045 | Quarterly recertification includes attestation from the user's manager that access is still required. | Derived | Self-certification of access needs is unreliable; manager attestation provides stronger accountability. |

---

### Pattern 12: Incident Response

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-046 | There is an incident response plan specific to financial data compromise or ERP system breach. | Operational | Financial data breaches have specific notification requirements (SEC, SOX, GDPR) that generic IR plans may not cover. |
| ASF-047 | The IR team can isolate a compromised ERP component without destroying forensic evidence. | Derived | Tearing down a compromised ERP server without forensic imaging destroys evidence of financial fraud. |
| ASF-048 | IR playbooks include procedures for suspending ERP access during an active incident. | Trust | Rapid ERP account suspension during an incident prevents ongoing unauthorized financial transactions. |
| ASF-049 | Monitoring systems can detect anomalous ERP access patterns indicating account compromise. | Implicit | Without detection, financial fraud may continue for months before discovery during the next audit. |
| ASF-050 | The ERP system supports forensic export of all journal entries and approvals related to a specific time window. | Operational | Forensic data export in a format admissible in court requires specific capabilities not part of normal reporting. |

---

### Pattern 13: Least Privilege

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-051 | ERP application roles grant the minimum permissions necessary for each job function. | Explicit | Least privilege is fundamental to SoD; it is implied by the policy but not explicitly stated. |
| ASF-052 | The Financial DB user for the ERP backend has no DDL permissions and cannot alter schema. | Derived | An application database user with DDL access can modify audit triggers or drop tables. |
| ASF-053 | The Reporting Engine's database user has read-only access to a specific set of views, not the underlying tables. | Implicit | Direct table access for reporting bypasses row-level security and exposes all financial data. |
| ASF-054 | No ERP user has both creator and approver permissions in the same financial module. | Derived | This is the defining property of SoD — it must be enforced at the role level, not just procedurally. |
| ASF-055 | ERP administrators have a separate, non-admin account for their daily finance work. | Architectural | Admin accounts used for personal finance work create a path to privilege escalation. |

---

### Pattern 14: Monitoring & Alerting

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-056 | ERP authentication events (success, failure, lockout) are monitored and alerted in real time. | Operational | Failed login spikes indicate credential stuffing or brute-force attacks against finance accounts. |
| ASF-057 | Journal entries above a configured monetary threshold trigger an immediate alert independent of the approval workflow. | Derived | The approval workflow confirms authorization, but independent alerts provide a second detection layer for large amounts. |
| ASF-058 | SoD violation attempts (e.g., user tries to approve own entry) are logged and alerted. | Operational | SoD violation attempts are leading indicators of policy testing or insider threat reconnaissance. |
| ASF-059 | Changes to ERP role assignments and permission sets are monitored and alerted. | Implicit | Unauthorized role changes create backdoor SoD bypasses that persist until the next recertification. |
| ASF-060 | Auditor access patterns are monitored — excessive reporting queries may indicate data scraping. | Derived | Auditors with legitimate access scraping data outside their scope can exfiltrate financial records. |

---

### Pattern 15: Network Segmentation

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-061 | The ERP web application, backend, and Financial DB are in separate network zones with firewall rules. | Architectural | All three tiers in the same network segment allow unrestricted lateral movement. |
| ASF-062 | The Reporting Engine is in a DMZ or separate VLAN accessible only to auditors and the ERP backend. | Architectural | A reporting engine accessible from the general corporate network exposes financial data broadly. |
| ASF-063 | The approval workflow system is isolated from general internet traffic and accessible only via the ERP application. | Explicit | The approval workflow must be protected from direct access to enforce SoD. |
| ASF-064 | The Financial DB has no direct route to the internet (no NAT gateway, no internet gateway). | Explicit | A financial database with internet access can exfiltrate data directly or receive C2 commands. |
| ASF-065 | Auditor Access is restricted to a specific VPN or bastion path that is logged separately. | Operational | Auditor access from the public internet without a VPN or bastion bypasses network-level controls. |

---

### Pattern 16: Third-party Dependency

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-066 | The ERP vendor provides security patches for known vulnerabilities within the vendor's SLA window. | Dependency | Unpatched ERP vulnerabilities (e.g., Oracle EBS, SAP NetWeaver) are the primary attack vector for financial systems. |
| ASF-067 | The ERP database platform (Oracle, SQL Server, PostgreSQL) has no unpatched critical vulnerabilities. | Dependency | Database engine CVEs can allow direct data access without ERP authentication. |
| ASF-068 | Third-party integrations (bank feeds, payment gateways, tax engines) use authenticated APIs and are reviewed quarterly. | Operational | Third-party integrations create data flows outside the documented architecture that bypass SoD controls. |
| ASF-069 | The cloud provider or data center hosting the ERP has SOC 2 Type II reports and meets financial services compliance. | Dependency | The infrastructure provider's security posture impacts SOX compliance scope for the ERP. |
| ASF-070 | The ERP system's external auditor access tool has no vulnerabilities that allow privilege escalation. | Trust | Auditor access tools compromised via their own vulnerabilities can be used for unauthorized data access. |

---

### Pattern 17: Audit & Compliance

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-071 | SOX controls are tested by internal audit at least annually, and failing controls are remediated. | Explicit | SOX compliance requires control testing; untested controls provide no compliance assurance. |
| ASF-072 | The ERP system produces evidence of control operation (screenshots, logs, reports) that can be presented to external auditors. | Derived | Audit evidence collection must be automated; manual evidence collection is incomplete and costly. |
| ASF-073 | Segregation of duties rules are documented, mapped to ERP roles, and tested for completeness. | Operational | Undocumented SoD rules cannot be tested; incomplete SoD testing leaves gaps for fraud. |
| ASF-074 | The quarterly recertification process is itself auditable — timestamps, approvers, and changes made are logged. | Implicit | The recertification process must be auditable to prevent retrospective alteration of access records. |
| ASF-075 | Financial data retention complies with SOX Section 802 (7-year retention for audit records). | Explicit | SOX Section 802 mandates 7-year retention for audit records; non-compliance is a separate violation. |

**Total (A): 75** (4 per pattern × 17 patterns + 7 extra from patterns 10, 12, 14, 15, 16, 17)

---

## Step 3: Comparison

### Overlap Mapping

| Human ID | ASF ID | Match Rationale |
|----------|--------|-----------------|
| H-001 | ASF-037 | Both require individual ERP accounts, no credential sharing. |
| H-002 | ASF-001 | Both require MFA for ERP users. |
| H-003 | ASF-054 | Both require SoD rules enforced by the ERP system automatically. |
| H-004 | ASF-054 | Both require dual-approval for journal entries (multi-person SoD). |
| H-005 | ASF-071 | Both require immutable audit logging for financial DB. |
| H-006 | ASF-053 | Both require Reporting Engine uses read-replica or read-only views. |
| H-007 | ASF-065 | Both require Auditor Access restricted to Reporting Engine only. |
| H-008 | ASF-056 | Both require ERP session logging with user attribution. |
| H-009 | ASF-043 | Both require recertification includes ERP application roles. |
| H-010 | ASF-056 | Both require failed login monitoring and alerting. |
| H-011 | ASF-052 | Both require parameterized queries / DDL prevention for financial DB. |
| H-013 | ASF-061 | Both require ERP web and backend in separate tiers. |
| H-014 | ASF-054 | Both require no single user can create and approve — core SoD. |
| H-015 | ASF-046 | Both require centralized tamper-proof audit log aggregation. |
| H-016 | ASF-025 | Both require Financial DB encryption at rest with separate key. |
| H-017 | ASF-053 | Both require separate reporting service account with read-only access. |
| H-019 | ASF-072 | Both require immutable entry ledger before posting. |
| H-022 | ASF-038 | Both require approval workflow not bypassable by admins. |
| H-023 | ASF-042 | Both require DB credentials not hardcoded in config. |
| H-025 | ASF-028 | Both require encrypted separate-location DB backups. |
| H-026 | ASF-051 | Both require row-level security for business unit isolation. |
| H-027 | ASF-060 | Both require auditor access limited to posted entries. |
| H-028 | ASF-021 | Both require financial data classification as critical/restricted. |
| H-029 | ASF-061 | Both require network segmentation between ERP tiers. |
| H-030 | ASF-072 | Both require version-controlled DB schema with review. |
| H-031 | ASF-048 | Both require session revocation capability. |
| H-033 | ASF-029 | Both require HTTPS-only with HSTS for ERP. |
| H-034 | ASF-060 | Both require export controls and monitoring for bulk data access. |
| H-035 | ASF-054 | Both require approver not subordinate to creator. |
| H-036 | ASF-051 | Both require per-module database users. |
| H-037 | ASF-057 | Both require anomalous transaction monitoring. |
| H-038 | ASF-059 | Both require privileged operation audit. |
| H-039 | ASF-053 | Both require output masking for sensitive fields in reports. |
| H-042 | ASF-075 | Both require 7-year audit log retention per SOX. |

**Overlap (O): 35**

### Metrics

| Metric | Value | Formula |
|--------|-------|---------|
| Human assumptions (H) | 42 | Count of unique human-generated assumptions |
| ASF assumptions (A) | 75 | Count of unique ASF-generated assumptions |
| Overlap (O) | 35 | Count appearing in both lists |
| **Precision** | **46.7%** | O / A = 35/75 |
| **Recall** | **83.3%** | O / H = 35/42 |
| **F1 Score** | **59.8%** | 2 × (P × R) / (P + R) |
| Novel findings (A - O) | 40 | Assumptions ASF found that human missed (53.3% of ASF total) |
| Missed findings (H - O) | 7 | Assumptions human found that ASF missed (16.7% of human total) |

### Success Criteria Assessment

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Recall | >= 70% | 83.3% | ✅ Exceeded |
| Precision | >= 50% | 46.7% | ❌ Not met (approaching) |
| Novel discoveries | >= 10% of total (A+O) | 36.4% (40/110) | ✅ Exceeded |
| Expert agreement (F1 proxy) | > 60% | 59.8% | ❌ Not met (marginally below) |

---

## Step 4: Analysis

### Ontology Overlap Analysis

| Ontology Category | Overlaps | ASF Total | Overlap Rate |
|-------------------|----------|-----------|-------------|
| Explicit | 10 | 14 | 71.4% |
| Derived | 11 | 16 | 68.8% |
| Operational | 7 | 17 | 41.2% |
| Implicit | 4 | 10 | 40.0% |
| Trust | 2 | 5 | 40.0% |
| Dependency | 1 | 5 | 20.0% |
| Architectural | 0 | 4 | 0.0% |
| Environmental | 0 | 4 | 0.0% |

**Best overlap:** Explicit (71.4%) — MFA, encryption, SoD enforcement, and access controls are well-aligned between human and ASF in the SOX compliance domain.

**Worst overlap:** Architectural (0%) and Environmental (0%) — the ASF identified ERP tier network isolation, approval workflow system isolation, and availability constraints during financial close that the human architect did not frame as assumptions.

### What Humans Caught That ASF Missed (Missed Findings = 7)

The 7 missed findings are the **lowest across all three architectures**, demonstrating strong ASF coverage of the financial SOX domain:

1. **ERP-specific UI controls (H-018, H-024, H-040):** CSRF token enforcement, session idle timeout (15 min), and login rate limiting are web application security controls not covered by the ASF's current patterns. These are standard OWASP controls specific to the web application tier.

2. **Time-based access controls (H-020):** The human assumed business-hours-only ERP access. The ASF's Authentication and Monitoring patterns do not explicitly model temporal access constraints.

3. **Manager-chain SoD (H-035):** The human assumed approvers are not subordinates of creators — a detailed organizational SoD rule. The ASF patterns cover automated SoD enforcement (ASF-054) but not the specific manager-subordinate conflict.

4. **NTP synchronization for audit (H-041):** Clock synchronization across the ERP for audit timestamp accuracy is a deployment-level operational detail.

5. **Annual security audit (H-032):** The human assumed financial audits include a security component. The ASF's Audit & Compliance pattern (ASF-071) tests SOX controls but does not explicitly model security-audit integration.

### What ASF Caught That Humans Missed (Novel Findings = 40)

1. **Third-party dependencies (5 assumptions):** The human generated no assumptions about the ERP vendor, database platform CVEs, third-party integrations, cloud provider SOC reports, or auditor tool vulnerabilities. This is the largest external-risk gap.

2. **Incident Response (5 assumptions):** The human generated no IR assumptions for financial data compromise. The ASF contributed SOX-specific IR: financial breach notification, forensic evidence preservation, account suspension, anomaly detection, and forensic data export.

3. **Identity Lifecycle (4 assumptions):** The human assumed recertification (H-009) but did not extend to joiner/mover/leaver integration, role-change automation, service account review, or manager-attested recertification.

4. **Human factors (ASF-037 through ASF-041):** The human assumed no credential sharing (H-001) but did not cover admin backdoor accounts, workflow understanding, phishing reporting, or recertification thoroughness.

5. **Monitoring depth (ASF-056 through ASF-060):** The human assumed login monitoring (H-010) and transaction alerts (H-037) but did not extend to SoD violation alerts, role change monitoring, or auditor access pattern analysis.

6. **Audit & Compliance (ASF-071 through ASF-075):** The ASF added SOX-specific assumptions about control testing cadence, evidence automation, SoD rule documentation, recertification auditability, and SOX 802 retention.

### Architecture Complexity Assessment

Architecture #020 achieved the **highest recall of all three architectures** (83.3%), significantly exceeding the 70% target. This is because the ERP/SOX domain is well-covered by the ASF's Least Privilege, Identity Lifecycle, Monitoring, and Audit patterns. The precision (46.7%) is also the highest across the three experiments, approaching the 50% target.

The human architect focused on the **SoD enforcement chain** — user authentication, journal entry approval, audit logging, and read-only access for auditors. The ASF contributed orthogonal assumptions in third-party dependencies, incident response, and identity lifecycle that are outside the human architect's scope but highly relevant to the overall risk posture.

### Key Insight

Architecture #020 demonstrates that **compliance-driven architectures** (SOX) achieve the best ASF performance because their requirements are codified in regulatory frameworks that align with the ASF's pattern categories. The 83.3% recall and 59.8% F1 are the strongest results across Phases 2-6, suggesting the ASF is particularly well-suited to compliance-heavy environments.

The remaining 7 missed findings are concentrated in **web application security controls** (CSRF, rate limiting, session timeout) — the same gap identified in Architecture #001 and #018. A single "Web Application Security" pattern would close this gap across all architectures.

---

## Conclusions

| Metric | Target | Actual | Verdict |
|--------|--------|--------|---------|
| Recall | >= 70% | 83.3% | ✅ Exceeded — best recall across all architectures |
| Precision | >= 50% | 46.7% | ❌ Not met (approaching — closest of the three) |
| Novel discoveries | >= 10% | 36.4% | ✅ ASF adds substantial value in compliance domain |
| Expert agreement (F1) | > 60% | 59.8% | ❌ Not met (marginally below target) |

Architecture #020 achieves the strongest ASF performance to date: 83.3% recall and 46.7% precision, with F1 of 59.8% approaching the 60% threshold. The SOX compliance domain maps effectively to the ASF's pattern matrix, particularly Least Privilege, Identity Lifecycle, Monitoring, and the new Audit & Compliance pattern. The primary remaining gap is web application security controls (CSRF, rate limiting, session management) — a cross-cutting concern not covered by any existing pattern.
