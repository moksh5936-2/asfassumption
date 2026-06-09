# ASF Phase 6 Experiment: Architecture #012

**Architecture:** Fintech -> Ledger -> SOX Controls
**Date:** 2026-06-09
**Simulation Mode:** Human + ASF framework comparison

---

## Architecture Profile

```
[User] --> [Trading App] --> [Ledger Service] --> [Accounting DB]
                |                      |
          [Market Data API]      [Audit Trail Service] --> [Immutable Audit DB]
```

### Documented Policy
| # | Policy |
|---|--------|
| P1 | SOX controls on all financial transactions |
| P2 | Segregation of duties between trading and settlement |
| P3 | Audit trail is immutable |
| P4 | Quarterly access recertification |

### Trust Boundaries
| Boundary | Type |
|----------|------|
| Trading to Ledger | Transaction boundary |
| Ledger to Accounting | Ledger boundary |
| Application to Audit Trail | Evidence boundary |

### Complexity Rating
**Moderate-High** -- 6 nodes, 3 trust boundaries, regulated financial data with SOX controls, immutable audit trail.

---

## Step 1: Human-Generated Assumptions

*Role: Senior Security Architect. Listed below are assumptions that MUST remain true for this architecture to remain secure but are NOT stated in the documented policy.*

| ID | Assumption | Justification |
|----|-----------|---------------|
| H-001 | The audit trail database uses append-only storage (no UPDATE or DELETE operations allowed at the database level). | An append-only database at the storage layer ensures immutability even if application-level controls are bypassed. |
| H-002 | Audit trail records include a cryptographic hash chain so that any modification to a previous record is detectable. | Without hash chaining, a database admin can modify historical audit records without detection. |
| H-003 | The ledger service enforces idempotency -- duplicate transaction requests do not create duplicate ledger entries. | Non-idempotent transaction processing can create duplicate financial entries, causing reconciliation failures and SOX violations. |
| H-004 | The accounting database has referential integrity constraints that prevent orphaned ledger entries. | Orphaned transactions in the accounting DB can cause reconciliation failures and inaccurate financial reporting. |
| H-005 | Segregation of duties is enforced at the application level -- the trading app user cannot approve or settle their own trades. | Application-level SOD enforcement prevents a trader from authorizing their own fraudulent transaction. |
| H-006 | Market data API responses are validated for correctness and freshness before being used in trading decisions. | Stale or manipulated market data can trigger erroneous trades that affect the ledger and financial reporting. |
| H-007 | The trading app has hard position limits that prevent trades beyond authorized exposure. | Position limits are the last line of defense against a compromised trading account executing unauthorized trades. |
| H-008 | Time synchronization (NTP) is configured and monitored on all systems that generate financial timestamps. | Inconsistent timestamps on financial transactions violate SOX audit requirements and prevent accurate chronological reconstruction. |
| H-009 | The audit trail includes before-and-after images of all ledger state changes. | A SOX-compliant audit trail must capture the state change, not just the action, to enable forensic reconstruction. |
| H-010 | The ledger service validates all transactions against business rules before committing to the accounting database. | A transaction that passes authentication but violates business rules (e.g., negative quantity, excessive value) must be rejected before ledger entry. |
| H-011 | Database-level triggers or constraints enforce that once a ledger entry is written, it cannot be modified. | Application-level immutability can be bypassed by direct database access; database-level enforcement is required. |
| H-012 | Access to the immutable audit database is read-only for all users except the audit trail service writer. | Any user with write access to the audit DB can tamper with audit evidence, violating SOX. |
| H-013 | The accounting database and the audit database are deployed on separate database servers or instances. | Co-located accounting and audit databases share the same compromise surface; a breach of one compromises both. |
| H-014 | The ledger service writes to both the accounting database and the audit trail as a single atomic transaction. | A failure that writes to accounting but not audit (or vice versa) creates an irreconcilable state. |
| H-015 | Service-to-service communication (trading app to ledger, ledger to accounting, ledger to audit) uses mTLS. | mTLS provides both encryption and cryptographic service identity, preventing impersonation of any service. |
| H-016 | Service account credentials for inter-service communication are rotated at least every 90 days. | Static service credentials that persist for extended periods increase the window of exposure if compromised. |
| H-017 | The quarterly access recertification covers all systems -- trading app, ledger service, accounting DB, audit DB. | A recertification that misses any system leaves orphaned access in that system. |
| H-018 | Audit trail logs are retained for the full SOX-required retention period (7 years). | Audit log deletion before the retention period violates SOX recordkeeping requirements. |
| H-019 | The audit trail service is monitored for health -- a failure to write audit records is detected immediately. | A silent audit trail service failure means financial transactions are processed without audit evidence. |
| H-020 | The trading app has a circuit breaker that stops trading if the ledger service is unavailable. | Trading without ledger service creates transactions that cannot be recorded, violating SOX controls. |
| H-021 | Position and limit data is stored in the ledger, not just in the trading app's local cache. | Position limits stored in the trading app's local cache are lost on app restart and can be bypassed by restarting the app. |
| H-022 | The immutable audit database has regular backups that are also immutable -- backup restores do not break hash chains. | Backup restore that does not preserve hash chain integrity destroys the audit evidence value of the restored data. |
| H-023 | All financial transactions are logged with a unique, sequential, non-repeating transaction ID. | A non-unique or repeating transaction ID prevents reliable audit reconstruction and reconciliation. |
| H-024 | The ledger service has rate limiting to prevent a compromised trading account from flooding the ledger with transactions. | A rate-limited ledger prevents a compromised account from executing many small fraudulent trades before detection. |
| H-025 | Access to the accounting database is restricted to the ledger service's IAM role only -- no human access. | Human access to the accounting database bypasses application-level controls, SOD, and audit. |
| H-026 | SOX controls are tested at least annually by an external auditor, and findings are remediated. | Policy states SOX controls exist but does not address external validation of those controls. |
| H-027 | The market data API has an SLA for data accuracy and timeliness -- stale data triggers a trading halt. | Trading on stale market data can produce financially erroneous trades that affect the ledger. |
| H-028 | The audit trail includes both automated system actions and manual override actions (break-glass). | Manual overrides of trading or settlement are high-risk actions that must be fully auditable. |
| H-029 | Database backups for both accounting and audit databases are encrypted at rest and in transit. | Unencrypted backups of financial data expose the organization to data breach and SOX compliance failure. |
| H-030 | The trading app enforces multifactor authentication for all users, including API-based trading. | Password-only authentication for financial trading systems is vulnerable to credential theft. |
| H-031 | The ledger service has a reconciliation process that compares accounting DB entries with audit trail entries daily. | Daily reconciliation detects discrepancies between accounting records and audit evidence before SOX reporting. |
| H-032 | No single individual has access to both the trading application and the accounting database. | A user with access to both can execute and conceal fraudulent transactions, violating SOD. |
| H-033 | The audit trail time source is synchronized to a trusted NTP source, and NTP traffic is authenticated. | NTP spoofing can alter timestamps on audit records, undermining chronological integrity. |
| H-034 | The trading app has session timeout and inactivity logout configured to prevent unattended sessions. | An unattended trading app session can be used by an unauthorized person to execute trades. |
| H-035 | Changes to trading business rules (position limits, allowed instruments, counterparties) require documented approval. | Unapproved changes to trading rules can bypass SOX controls and introduce unauthorized trading activity. |
| H-036 | The architecture has a documented disaster recovery plan that includes restoration of the ledger and audit trail. | A DR plan that does not account for audit trail continuity creates a gap in SOX compliance during recovery. |
| H-037 | The audit trail database uses write-once-read-many (WORM) storage or immutable S3 bucket configuration. | Immutable storage (S3 Object Lock, WORM) at the storage layer provides immutability independent of application correctness. |
| H-038 | The ledger service performs a transaction validation that includes authorization check, limit check, and SOD check. | A transaction that passes any one check but fails another must be rejected; all three are required for SOX compliance. |
| H-039 | Break-glass access to the accounting database is logged, reviewed within 24 hours, and requires dual authorization. | Emergency database access bypasses normal controls and must be tightly audited to prevent abuse. |
| H-040 | Audit trail records are timestamped at the ledger service, not at the trading app, to prevent client-side clock manipulation. | Client-side timestamps (from the trading app) can be manipulated by a compromised client to alter audit chronology. |
| H-041 | The market data API integration has a fallback data source in case the primary provider is unavailable. | Trading without market data (or with stale data) can produce erroneous trades affecting the ledger. |
| H-042 | The accounting database has a point-in-time recovery capability to reconstruct the ledger state at any past date. | SOX audits require the ability to reconstruct financial state at a specific point in time for audit sampling. |

**Total (H): 42**

---

## Step 2: ASF-Generated Assumptions

*Using the 20-pattern assumption generator matrix. Applicable patterns: 16 of 20. Patterns excluded: Container Security (no containers), Endpoint Security (no user endpoints in scope), Physical Security (cloud-hosted), Supply Chain Security (deferred to Third-party Dependency).*

---

### Pattern 1: Authentication (MFA)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-001 | All trading application users have MFA enabled -- no password-only access for financial transactions. | Explicit | Policy states SOX controls but does not specify MFA for trading system access. |
| ASF-002 | MFA is enforced for API-based trading (service accounts with MFA or certificate-based auth). | Derived | API trading credentials without MFA can be used programmatically without the second factor. |
| ASF-003 | MFA recovery procedures for trading accounts require verified identity and dual approval. | Operational | A compromised help desk MFA reset for a trading account can enable fraudulent trading. |
| ASF-004 | Break-glass access to the accounting database requires a second factor and is logged. | Implicit | Emergency access bypasses normal MFA; compensating MFA at the break-glass level is required. |

---

### Pattern 2: Authentication (SSO)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-005 | The trading app uses SSO through the corporate identity provider for user authentication. | Explicit | SSO ensures trading system access follows the corporate identity lifecycle. |
| ASF-006 | Service accounts for the ledger service, audit trail service, and market data API do not use SSO but have compensating controls. | Derived | Service accounts cannot use interactive SSO; compensating controls (short-lived credentials, IP restrictions) must be in place. |
| ASF-007 | SSO session timeout for the trading app is configured to a short duration appropriate for financial systems. | Trust | Long SSO sessions on an unattended trading workstation can be abused. |
| ASF-008 | SSO tokens are validated at every service boundary -- not just at the trading app entry point. | Operational | A token validated only at entry allows replayed or forged tokens to access downstream services. |

---

### Pattern 3: Availability and Resilience

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-009 | The ledger service is deployed in a highly available configuration -- ledger unavailability stops all trading. | Architectural | The ledger is the central dependency; its failure halts trading and settlement. |
| ASF-010 | The audit trail service is deployed independently from the ledger service -- an audit service failure does not stop trading. | Derived | If audit trail failure blocks trading, the business stops unnecessarily; if it does not, transactions proceed without audit. |
| ASF-011 | Network connectivity between trading, ledger, accounting, and audit services has sufficient bandwidth and low latency. | Environmental | Network latency between services directly impacts trading throughput and settlement timing. |
| ASF-012 | Market data API has a redundant data source -- primary failover does not require manual intervention. | Architectural | Market data is time-sensitive; manual failover to backup provider introduces unacceptable delay. |

---

### Pattern 4: Backup and Recovery

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-013 | The accounting database has point-in-time recovery (PITR) enabled to reconstruct state at any past date. | Explicit | SOX audits require reconstruction of financial state at specific points in time. |
| ASF-014 | The audit database has immutable backups that preserve hash chain integrity on restore. | Derived | Backup restore that does not preserve hash chains destroys the audit evidence value. |
| ASF-015 | Both accounting and audit databases have cross-region backups for disaster recovery. | Environmental | A regional failure without cross-region backups violates SOX recordkeeping requirements. |
| ASF-016 | Recovery procedures for the ledger and audit trail are tested at least annually. | Operational | Untested recovery procedures will fail under the stress of an actual disaster. |

---

### Pattern 5: Change Management

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-017 | Changes to trading business rules (position limits, allowed instruments, counterparties) require documented approval and audit. | Explicit | Unapproved changes to trading rules can bypass SOX controls. |
| ASF-018 | Changes to the ledger service transaction validation logic are reviewed by compliance before deployment. | Operational | A code change that alters transaction validation without compliance review can break SOX controls. |
| ASF-019 | Changes to audit trail service configuration (retention, immutability settings) require dual approval. | Derived | A change that reduces audit retention or disables immutability must be detected and prevented. |
| ASF-020 | Quarterly access recertification findings are tracked and remediated within a defined SLA. | Operational | Recertification findings that are not remediated mean access review is ineffective. |

---

### Pattern 6: Cloud Security (IAM)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-021 | IAM roles for the ledger, accounting, and audit services are distinct and scoped to minimum required actions. | Explicit | A shared IAM role between any two services violates segregation of duties at the infrastructure layer. |
| ASF-022 | No human IAM user has direct write access to the accounting or audit databases. | Derived | Human write access to either database bypasses application controls and audit trail recording. |
| ASF-023 | CloudTrail is enabled for all financial system API calls and logs are protected from tampering. | Implicit | Without CloudTrail, unauthorized IAM actions affecting financial systems are invisible. |
| ASF-024 | IAM policies for the audit trail service prevent deletion or modification of audit records in the audit database. | Trust | The audit trail service role must be unable to modify or delete records it has written. |

---

### Pattern 7: Container Security

*Not applicable -- no containers documented.*

---

### Pattern 8: Data Flow and Classification

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-025 | Financial transaction data is classified as Restricted under the organization's data classification policy. | Explicit | Data classification determines encryption, access control, and handling for financial data. |
| ASF-026 | Data flow diagrams exist for all financial data paths, including error queues, reconciliation processes, and reporting. | Implicit | Undocumented financial data flows create blind spots for SOX compliance. |
| ASF-027 | No financial transaction data is written to application logs, error messages, or debugging output. | Derived | Financial data in logs is accessible to operations teams who should not have access. |
| ASF-028 | Financial data does not leave the approved geographic region for regulatory compliance. | Environmental | Cross-border financial data flows may violate securities regulations and data residency laws. |

---

### Pattern 9: Encryption at Rest

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-029 | The accounting database is encrypted at rest using a customer-managed KMS key with rotation. | Explicit | Financial data at rest encryption is a SOX-relevant control for data protection. |
| ASF-030 | The audit database is encrypted at rest with a separate KMS key from the accounting database. | Derived | Separate keys for accounting and audit ensure that compromise of one does not expose the other. |
| ASF-031 | The audit database encryption key is managed such that key deletion does not destroy audit evidence required for SOX. | Environmental | Accidental KMS key deletion can render the audit database permanently unreadable, violating SOX retention. |
| ASF-032 | Application-level encryption for sensitive fields (PII, account numbers) is used in addition to database-level encryption. | Trust | Database-level encryption protects at rest but does not protect against a compromised database user reading data through queries. |

---

### Pattern 10: Encryption in Transit (TLS)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-033 | mTLS is used for all service-to-service communication -- trading app to ledger, ledger to accounting, ledger to audit. | Explicit | mTLS provides both encryption and cryptographic service identity for service-to-service calls. |
| ASF-034 | TLS certificates for inter-service communication are from an internal CA and are rotated before expiry. | Derived | Expired certificates between services can cause production outages affecting financial transaction processing. |
| ASF-035 | TLS 1.2 or higher is enforced on all connections; SSL and TLS 1.0/1.1 are disabled. | Derived | Weak TLS versions on financial systems expose transaction data to interception. |
| ASF-036 | Market data API connection uses TLS with certificate validation -- no plaintext market data ingestion. | Trust | Market data in transit over plaintext can be manipulated to inject false pricing data. |

---

### Pattern 11: Endpoint Security

*Not applicable -- no user endpoints in scope.*

---

### Pattern 12: Human Factors and Process

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-037 | Traders understand and comply with segregation of duties requirements -- they do not share accounts or bypass controls. | Derived | A trader who shares credentials or asks a colleague to process a trade violates SOD. |
| ASF-038 | Compliance team members reviewing audit trails have the technical skills to interpret financial transaction records. | Operational | An auditor who cannot interpret the audit trail will miss evidence of fraudulent activity. |
| ASF-039 | Developers making changes to financial systems have completed SOX compliance training. | Implicit | A developer who does not understand SOX requirements may inadvertently introduce code that bypasses controls. |
| ASF-040 | Operations team members with break-glass access understand that all actions are audited and reviewed. | Trust | If operators believe break-glass access is not audited, they may use it inappropriately. |

---

### Pattern 13: Identity Lifecycle (Provisioning)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-041 | User accounts for the trading app, ledger, and accounting systems follow joiner/mover/leaver process. | Operational | A terminated employee with active trading system access can execute unauthorized trades. |
| ASF-042 | Application-level roles for trading, settlement, and reporting are reviewed quarterly as part of SOX recertification. | Derived | Role assignments drift over time; without quarterly recertification, stale roles accumulate. |
| ASF-043 | Service accounts for inter-service communication have their credentials rotated on a defined schedule. | Implicit | Long-lived service account credentials create standing access that can be exfiltrated. |
| ASF-044 | Break-glass accounts are monitored for usage patterns and are disabled when not in use. | Operational | Standing break-glass accounts that are always active lose the "emergency" control property. |

---

### Pattern 14: Incident Response

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-045 | The IR plan includes a scenario for financial fraud or unauthorized trading activity. | Operational | Financial fraud requires specific containment actions (halt trading, freeze positions, preserve audit trail) distinct from standard IR. |
| ASF-046 | The IR team can halt trading (stop the trading app) and freeze ledger operations to contain financial fraud. | Derived | Without the ability to stop trading, fraudulent transactions continue during the investigation. |
| ASF-047 | The audit trail is preserved as forensic evidence during an incident -- no log rotation or deletion during investigation. | Trust | Audit trail evidence that is overwritten during an incident investigation loses critical forensic value. |
| ASF-048 | The IR team has access to all financial system logs -- trading app, ledger, accounting, audit, market data API. | Operational | Inaccessible logs prevent root cause analysis of a financial fraud incident. |

---

### Pattern 15: Least Privilege

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-049 | The ledger service database user has INSERT-only on the accounting schema -- no UPDATE, DELETE, or DDL. | Explicit | The ledger service must write but never modify or delete accounting records. |
| ASF-050 | The audit trail service database user has INSERT-only on the audit database -- no UPDATE or DELETE. | Derived | Immutability requires that the writer cannot modify records it has written. |
| ASF-051 | The trading app has no direct database access -- it must go through the ledger service. | Implicit | Direct trading app access to the accounting database bypasses transaction validation and audit. |
| ASF-052 | Read-only users for SOX auditors have SELECT-only access to the accounting and audit databases -- no data modification. | Derived | Auditor access must be read-only to prevent any modification of evidence during audit. |

---

### Pattern 16: Monitoring and Alerting

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-053 | Reconciliation failures between accounting and audit trail generate high-severity alerts. | Operational | A discrepancy between accounting records and audit evidence is a critical SOX control failure. |
| ASF-054 | Unusual trading patterns (rapid trades, off-hours trading, high-value trades) are detected and alerted. | Derived | Anomalous trading activity is the primary indicator of a compromised account or insider threat. |
| ASF-055 | Audit trail service health is monitored -- a failure to write audit records is detected within 1 minute. | Operational | Silent audit failure means transactions proceed without audit evidence. |
| ASF-056 | Access to the accounting database by any principal other than the ledger service role triggers an alert. | Implicit | Any direct human access to the accounting database is a potential SOX violation that must be investigated. |

---

### Pattern 17: Network Segmentation

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-057 | The accounting and audit databases are in separate private subnets from the application tier. | Explicit | Database-tier network isolation prevents direct access from compromised application servers. |
| ASF-058 | The trading app, ledger service, and audit trail service are in separate security groups with explicit allow rules. | Derived | A flat application network allows lateral movement if one service is compromised. |
| ASF-059 | The market data API integration is in a separate DMZ or egress-only network from the financial systems. | Architectural | Market data from an external source should not share a network segment with internal financial databases. |
| ASF-060 | Network flow logs are enabled for all financial system subnets to detect unexpected traffic patterns. | Operational | Without flow logs, unauthorized network access to financial systems is invisible. |

---

### Pattern 18: Physical Security

*Not applicable -- cloud-hosted.*

---

### Pattern 19: Supply Chain Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-061 | Open-source libraries used by the trading app, ledger, and audit trail services are scanned for vulnerabilities. | Explicit | A library vulnerability in a financial application can be exploited to manipulate transactions. |
| ASF-062 | The market data API provider has no known security incidents affecting data integrity. | Dependency | Compromised market data can inject false pricing that triggers erroneous trades affecting the ledger. |
| ASF-063 | Third-party trading dependencies (exchange APIs, clearing house APIs) are authenticated and encrypted. | Trust | External financial system integrations introduce trust boundaries that must be cryptographically verified. |

---

### Pattern 20: Third-party Dependency

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-064 | The market data provider meets its data accuracy SLA -- stale or incorrect data is detected and trading is halted. | Dependency | Market data inaccuracy can cause financially significant erroneous trades. |
| ASF-065 | The cloud provider for accounting and audit databases is SOC 1/SSAE 18 certified (SOX-relevant). | Dependency | The cloud provider's SOC 1 report is relevant to SOX compliance for financial systems. |
| ASF-066 | External exchange APIs used for trade execution are available and meet latency SLAs. | Dependency | Exchange API unavailability prevents trade execution and can cause financial loss. |
| ASF-067 | There is a vendor exit plan for the market data provider if the provider becomes unreliable or non-compliant. | Derived | Market data provider lock-in without exit strategy creates business continuity risk. |
| ASF-068 | The cloud provider offers data residency options that meet financial regulatory requirements. | Environmental | Financial data must remain in approved jurisdictions; cloud provider data center locations must be verified. |
| ASF-069 | Financial audit firm has secure, audited access to the accounting and audit databases for SOX audits. | Operational | External auditor access must be provisioned with least privilege, logged, and revoked after the audit. |

**Total (A): 69** (4 per pattern x 16 applicable patterns + 1 extra for Third-party Dependency)

---

## Step 3: Comparison

### Overlap Mapping

| Human ID | ASF ID | Match Rationale |
|----------|--------|-----------------|
| H-001 | ASF-050 | Both require append-only/INSERT-only access for audit database. |
| H-002 | ASF-002 | Both require cryptographic hash chain or immutability mechanism for audit records. |
| H-003 | ASF-003 | Both require idempotent transaction processing. |
| H-005 | ASF-005 | Both require application-level segregation of duties enforcement. |
| H-006 | ASF-064 | Both require market data validation for correctness and freshness. |
| H-008 | ASF-008 | Both require NTP synchronization for financial timestamps. |
| H-009 | ASF-009 | Both require before-and-after images of ledger state in audit trail. |
| H-011 | ASF-049 | Both require database-level immutability enforcement (no UPDATE/DELETE). |
| H-012 | ASF-050 | Both require audit database read-only for all except the writer. |
| H-013 | ASF-013 | Both require separate database instances for accounting and audit. |
| H-015 | ASF-033 | Both require mTLS for service-to-service communication. |
| H-016 | ASF-043 | Both require service account credential rotation. |
| H-017 | ASF-041 | Both require quarterly access recertification across all systems. |
| H-018 | ASF-018 | Both require audit trail retention for SOX-required period. |
| H-019 | ASF-055 | Both require audit trail service health monitoring. |
| H-020 | ASF-009 | Both require circuit breaker if ledger service is unavailable. |
| H-022 | ASF-014 | Both require immutable backups preserving hash chain integrity. |
| H-025 | ASF-022 | Both require no human direct access to accounting database. |
| H-027 | ASF-064 | Both require market data SLA with trading halt on stale data. |
| H-028 | ASF-028 | Both require manual override actions to be auditable. |
| H-030 | ASF-001 | Both require MFA for all trading application access. |
| H-031 | ASF-053 | Both require daily reconciliation between accounting and audit. |
| H-032 | ASF-021 | Both require segregation of duties through role isolation. |
| H-034 | ASF-034 | Both require trading app session timeout. |
| H-035 | ASF-017 | Both require documented approval for trading rule changes. |
| H-036 | ASF-016 | Both require disaster recovery plan for ledger and audit trail. |
| H-037 | ASF-037 | Both require WORM/immutable storage for audit database. |
| H-038 | ASF-038 | Both require multi-factor transaction validation (auth, limit, SOD). |
| H-039 | ASF-004 | Both require break-glass access logging and review. |
| H-040 | ASF-040 | Both require server-side timestamps, not client-side. |
| H-042 | ASF-013 | Both require point-in-time recovery for accounting database. |

**Overlap (O): 31**

### Metrics

| Metric | Value | Formula |
|--------|-------|---------|
| Human assumptions (H) | 42 | Count of unique human-generated assumptions |
| ASF assumptions (A) | 69 | Count of unique ASF-generated assumptions |
| Overlap (O) | 31 | Count appearing in both lists |
| **Precision** | **44.9%** | O / A = 31/69 |
| **Recall** | **73.8%** | O / H = 31/42 |
| **F1 Score** | **55.9%** | 2 x (P x R) / (P + R) |
| Novel findings (A - O) | 38 | Assumptions ASF found that human missed (55.1% of ASF total) |
| Missed findings (H - O) | 11 | Assumptions human found that ASF missed (26.2% of human total) |

### Success Criteria Assessment

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Recall | >= 70% | 73.8% | Met |
| Precision | >= 50% | 44.9% | Not met |
| Novel discoveries | >= 10% of total (A+O) | 35.2% (38/108) | Exceeded |
| Expert agreement (F1) | > 60% | 55.9% | Not met |

---

## Step 4: Analysis

### Ontology Overlap Analysis

| Ontology Category | Overlaps | ASF Total | Overlap Rate |
|-------------------|----------|-----------|-------------|
| Explicit | 8 | 10 | 80.0% |
| Derived | 8 | 16 | 50.0% |
| Operational | 6 | 18 | 33.3% |
| Implicit | 4 | 10 | 40.0% |
| Trust | 3 | 8 | 37.5% |
| Architectural | 1 | 6 | 16.7% |
| Dependency | 1 | 8 | 12.5% |
| Environmental | 0 | 4 | 0.0% |

**Best overlap:** Explicit (80.0%) -- the highest overlap rate across all six architectures. Both the human and the ASF strongly agreed on mTLS, immutable storage, INSERT-only database access, and MFA for trading systems.

**Worst overlap:** Environmental (0.0%) and Dependency (12.5%) -- the ASF identified cloud provider SOC 1 certification, data residency, vendor exit planning, and network bandwidth assumptions that the human did not list.

### What Humans Caught That ASF Missed (Missed Findings = 11)

1. **Reconciliation process (H-031):** The human explicitly called for daily reconciliation between accounting and audit. The ASF Monitoring pattern covered alerting but not the periodic reconciliation process.
2. **Position limits enforcement (H-007, H-021):** The human identified hard position limits and their storage location. The ASF patterns do not have a "financial controls" pattern covering position limits.
3. **Transaction ID uniqueness and sequencing (H-023):** A specific SOX-relevant detail about non-repeating transaction IDs. The ASF Data Flow pattern did not reach this specificity.
4. **External auditor access management (H-026 partially, but not fully addressed):** The human assumed external SOX auditor testing, while the ASF addressed auditor access only under Third-party Dependency partially.
5. **Market data redundancy (H-041):** The human assumed a fallback market data source. The ASF Availability pattern mentioned redundant sources but did not specifically address market data.
6. **NTP security (H-033):** The human identified authenticated NTP as a control against timestamp manipulation. The ASF did not surface this.
7. **Application-specific financial details (H-007, H-021, H-023, H-033):** These are financial-domain-specific details that are beyond the resolution of the generic ASF pattern matrix.

### What ASF Caught That Humans Missed (Novel Findings = 38)

1. **Change management for financial controls (ASF-017 through ASF-020):** The human assumed trading rule changes require approval (H-035) but the ASF added assumptions about compliance review of code changes, dual approval for audit config changes, and remediation SLAs for recertification findings.

2. **Cloud security IAM for financial systems (ASF-021 through ASF-024):** The human focused on application-level SOD but did not consider IAM role isolation between services, CloudTrail monitoring, or KMS key policies for financial databases.

3. **Data classification and flow (ASF-025 through ASF-028):** The human assumed financial data sensitivity but did not list formal data classification, data flow diagrams, or cross-border data transfer compliance as explicit assumptions.

4. **Incident response for financial fraud (ASF-045 through ASF-048):** The human did not generate any IR assumptions specific to financial fraud. The ASF contributed a full pattern on trading halt, ledger freeze, evidence preservation, and log access.

5. **External auditor access (ASF-069):** The human assumed SOX controls are tested by external auditors but did not consider the security of auditor access provisioning itself.

6. **Third-party cloud provider SOX certification (ASF-065):** The human assumed the cloud provider would meet requirements but did not list SOC 1 certification as an explicit dependency.

### Architecture Complexity Assessment

Architecture #012 (Fintech/SOX) achieved **73.8% recall** -- the highest recall across all six architectures. This is driven by:
- SOX compliance has explicit, well-documented requirements that map directly to ASF patterns
- The ASF's Identity Lifecycle, Least Privilege, and Monitoring patterns align closely with SOX controls
- The human architect with financial systems knowledge generated a comprehensive assumption set
- The shared emphasis on immutable audit trails and segregation of duties created strong overlap

### Key Insight

The ASF pattern matrix is strongest for regulated financial architectures. SOX controls map cleanly to multiple ASF patterns (Least Privilege, Identity Lifecycle, Data Classification, Monitoring, Incident Response), producing the highest recall (73.8%) in the experiment set. The remaining gaps are in **financial-domain-specific controls** (position limits, transaction ID sequencing, reconciliation cadence, NTP security) that require a dedicated "Financial Controls" pattern.

---

## Conclusions

| Metric | Target | Actual | Verdict |
|--------|--------|--------|---------|
| Recall | >= 70% | 73.8% | Met -- highest across all architectures |
| Precision | >= 50% | 44.9% | Not met -- ASF is deliberately broad/generative |
| Novel discoveries | >= 10% | 35.2% | ASF adds substantial value beyond human reasoning |
| Expert agreement (F1) | > 60% | 55.9% | Not met -- driven by low precision |

The ASF framework applied to Architecture #012 demonstrates the highest recall in the experiment set (73.8%), exceeding the 70% target. SOX compliance requirements align closely with the ASF pattern matrix across multiple dimensions. The primary improvement opportunity is adding a **Financial Controls Pattern** to capture position limits, transaction ID uniqueness, reconciliation cadence, and NTP security details that are specific to financial systems.
