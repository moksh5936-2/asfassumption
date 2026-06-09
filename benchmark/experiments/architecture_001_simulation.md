# ASF Phase 6 Experiment: Architecture #1

**Architecture:** User → VPN → Internal App → Payroll DB
**Date:** 2026-06-09
**Simulation Mode:** Human + ASF framework comparison

---

## Architecture Profile

```
[User Laptop] --VPN--> [VPN Gateway] --TLS--> [Internal Web App] --SQL--> [Payroll Database (RDS)]
```

### Documented Policy
| # | Policy |
|---|--------|
| P1 | VPN required for remote access |
| P2 | Application authenticates with AD credentials |
| P3 | Database is in private subnet |
| P4 | Backups run nightly |

### Trust Boundaries
| Boundary | Type |
|----------|------|
| User ↔ VPN | Authentication boundary |
| VPN ↔ Application | Network boundary |
| Application ↔ Database | Data boundary |

### Complexity Rating
**Simple** — linear topology, 4 nodes, 3 trust boundaries, well-understood components.

---

## Step 1: Human-Generated Assumptions

*Role: Senior Security Architect. Listed below are assumptions that MUST remain true for this architecture to remain secure but are NOT stated in the documented policy.*

| ID | Assumption | Justification |
|----|-----------|---------------|
| H-001 | The VPN gateway runs fully patched firmware free of known CVEs. | VPN gateways are the primary internet-facing attack surface; an unpatched vulnerability can bypass all downstream controls. |
| H-002 | VPN client software on user laptops is centrally managed and disables split-tunneling. | Split-tunneling creates a bridge between the corporate network and the internet on a potentially compromised laptop. |
| H-003 | VPN credentials are never shared between users. | Shared credentials eliminate individual accountability and prevent targeted revocation. |
| H-004 | The VPN gateway enforces multi-factor authentication for all users. | Password-only VPN access is defeated by credential theft (phishing, password reuse). |
| H-005 | The private subnet containing the database has no route to or from the internet. | A private subnet with an implicit internet gateway or NAT route defeats the purpose of isolation. |
| H-006 | TLS certificates on the VPN gateway and internal web app are valid, not expired, and from a trusted CA. | Expired or untrusted certificates signal misconfiguration or active MITM attacks. |
| H-007 | The VPN gateway enforces modern encryption (AES-256, SHA-256) and disables weak protocols (PPTP, IKEv1). | Weak cryptographic protocols can be broken offline or via downgrade attacks. |
| H-008 | Active Directory is hardened against common attack techniques (Kerberoasting, DCSync, pass-the-hash). | AD compromise gives an attacker control over all authentication and authorization decisions. |
| H-009 | The web application does not store database credentials in source code, config files, or environment variables readable by non-admin users. | Hardcoded credentials are the most common vector for lateral movement after a code repository breach. |
| H-010 | The database user the application uses has least-privilege access (SELECT/INSERT/UPDATE on specific tables only, no DDL). | An application DB user with DDL or cross-schema access turns a SQL injection into full data compromise. |
| H-011 | TLS between the VPN gateway and the web application is properly terminated and validated at both ends (no self-signed or wildcard certs in production). | Improper TLS validation renders the encryption channel vulnerable to MITM within the internal network. |
| H-012 | The internal web application has no public IP address or internet-facing load balancer. | A direct internet path to the web application bypasses the VPN entirely, making the primary control irrelevant. |
| H-013 | The database subnet cannot initiate connections to the VPN gateway or the web application. | A compromised database could be used as a pivot point to attack other internal systems. |
| H-014 | Nightly backups are encrypted at rest. | Unencrypted backup files stored in S3 or backup appliances can be exfiltrated without any authentication against the live database. |
| H-015 | Backup restore procedures are tested at least annually to verify data integrity and RTO feasibility. | Untested backups are indistinguishable from no backups; corrupt or incomplete restores cause data loss. |
| H-016 | Security groups or NACLs restrict database access to only the web application's IP/port and nothing else. | A private subnet without explicit ingress rules still allows any other resource in the VPC to reach the database. |
| H-017 | The web application validates all user inputs and uses parameterized queries to prevent SQL injection. | SQL injection is the most direct path from authenticated user input to payroll data exfiltration. |
| H-018 | The web application implements secure session management (HttpOnly + Secure cookies, short timeouts, invalidation on logout). | Stolen or forged session cookies allow an attacker to impersonate a legitimate user without VPN credentials. |
| H-019 | VPN gateway logs (successful and failed connections, source IPs, durations) are forwarded to a SIEM. | Without VPN logging, brute-force attacks or unauthorized access attempts are invisible to defenders. |
| H-020 | The web application logs authentication events and access to sensitive payroll records. | Insider threats and compromised accounts accessing payroll data go undetected without application-level audit. |
| H-021 | Active Directory enforces account lockout after a threshold of failed login attempts. | Without lockout, an attacker can brute-force AD credentials through the VPN or application login. |
| H-022 | VPN access is revoked within the same shift when an employee is terminated or changes role. | Former employees or role-changed users with active VPN access retain the ability to reach payroll data. |
| H-023 | Database audit logging is enabled for all queries (SELECT, INSERT, UPDATE, DELETE). | Without DB audit, unauthorized reads or writes via compromised application credentials cannot be forensically traced. |
| H-024 | The web application enforces role-based access control so that not all authenticated users can view all payroll records. | AD authentication alone is binary (access vs. no access); RBAC is required for least privilege within the application. |
| H-025 | No SSH, RDP, or bastion host rules allow administrative access to the database from the general corporate network. | Direct administrative access broadens the attack surface; database administration should require a separate jump box with audit. |
| H-026 | The VPN gateway supports sufficient concurrent connections for all remote employees. | Insufficient capacity leads to dropped users who may then resort to insecure workarounds. |
| H-027 | The web application enforces rate limiting on login endpoints to prevent credential stuffing. | Without rate limiting, a credential stuffing attack using breached password databases can succeed silently. |
| H-028 | The TLS version between VPN gateway and web application is restricted to TLS 1.2 or higher. | TLS 1.0 and 1.1 have publicly known cryptographic weaknesses (BEAST, POODLE). |
| H-029 | The payroll database is not shared with other applications. | A vulnerability in a co-tenanted application could be used to access payroll data through the shared database. |
| H-030 | A firewall or security group separates the VPN gateway subnet from the web application subnet, enforcing explicit allow rules. | Defense-in-depth requires that even authenticated VPN traffic be inspected and restricted at the application layer. |
| H-031 | Automatic patching is enabled on the RDS instance for critical database engine vulnerabilities. | Unpatched database engines (e.g., MySQL, PostgreSQL) have known remote code execution vulnerabilities. |
| H-032 | The VPN client has a kill switch that terminates all network traffic if the VPN connection drops. | Without a kill switch, a dropped VPN tunnel leaks corporate-bound traffic over the user's unencrypted internet connection. |
| H-033 | The web application implements defenses against XSS, CSRF, and clickjacking. | Client-side attacks can steal session tokens or perform actions on behalf of authenticated users without their knowledge. |
| H-034 | Payroll data is not written to application logs, error messages, or debugging output. | Sensitive data in logs can be exposed through log aggregation systems, SIEM queries, or support ticket attachments. |
| H-035 | Database credentials are rotated on a regular cadence and immediately upon suspicion of compromise. | Static credentials increase the window of exposure; rotated credentials limit the blast radius of a leak. |
| H-036 | The VPN gateway enforces device compliance checks (OS patch level, antivirus running, disk encryption) before allowing connection. | An unmanaged, compromised laptop on the VPN is an authenticated beachhead for lateral movement. |
| H-037 | The web application does not expose database error details (stack traces, query fragments) to end users. | Leaked schema details in error messages inform and accelerate SQL injection attacks. |
| H-038 | The RDS instance has deletion protection and a final snapshot policy enabled. | Accidental or malicious deletion of the payroll database causes extended downtime even if backups exist. |
| H-039 | The VPN gateway uses perfect forward secrecy so that compromise of the private key does not decrypt past sessions. | Without PFS, an attacker who obtains the VPN server key can decrypt all recorded VPN traffic retroactively. |
| H-040 | AD authentication does not fall back to NTLMv1 or LM hash storage. | NTLMv1 and LM hashes are trivially cracked and expose user credentials to offline brute-force. |

**Total (H): 40**

---

## Step 2: ASF-Generated Assumptions

*Using the 20-pattern assumption generator matrix. Applicable patterns: 15 of 20. Patterns excluded: Container Security (no containers), Physical Security (cloud-hosted, no on-prem DC), Supply Chain Security (deferred to Third-party Dependency as a combined pattern), Change Management (covered under Operational cross-cutting concerns).*

---

### Pattern 1: Authentication (MFA)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-001 | VPN gateway enforces MFA for all remote users. | Explicit | Direct security requirement for remote access protection. |
| ASF-002 | MFA recovery codes are securely stored and usable only through verified identity proofing. | Derived | Lost MFA devices require a recovery path; unverified recovery bypasses MFA entirely. |
| ASF-003 | The help desk has a documented social-engineering-resistant process for MFA token reset. | Operational | Ad-hoc MFA resets are vulnerable to impersonation attacks. |
| ASF-004 | MFA is enforced on the VPN gateway, not just the web application. | Implicit | Policy states "VPN required" but does not specify where MFA is enforced; enforcement might be misapplied. |

---

### Pattern 2: Authentication (SSO)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-005 | The web application validates AD/SSO tokens and does not accept unauthenticated requests. | Explicit | Authentication integration must be correctly implemented at the application level. |
| ASF-006 | The AD domain controller is available and reachable from the web application at all times. | Dependency | AD downtime blocks all application authentication, creating an availability dependency. |
| ASF-007 | SSO session timeout is consistent between AD and the web application. | Trust | Inconsistent timeout configuration can leave application sessions active after AD session expiry. |
| ASF-008 | SSO token signing keys are rotated and protected from unauthorized access. | Operational | Compromised signing keys allow forged authentication tokens for any user. |

---

### Pattern 3: Availability & Resilience

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-009 | A single VPN gateway failure does not block all remote access (redundancy exists). | Architectural | The documented topology shows a single VPN gateway; its failure would deny all remote access. |
| ASF-010 | There is a documented offline procedure for VPN outages that does not bypass security controls. | Operational | Users under pressure to meet deadlines will resort to insecure alternatives if VPN is unavailable. |
| ASF-011 | The internet circuit to the VPN gateway has sufficient bandwidth and reliability SLA. | Environmental | ISP and connectivity are outside organizational control but critical to the architecture. |
| ASF-012 | The web application can gracefully handle database connection failures without leaking data or crashing. | Derived | An unavailable database should not cause the application to expose error states or cached sensitive data. |

---

### Pattern 4: Backup & Recovery

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-013 | Nightly backups complete within the backup window without impacting database performance. | Operational | Backups that fail silently or degrade production performance have operational consequences. |
| ASF-014 | Backup restore is tested at least annually to validate data integrity. | Derived | Policy documents "backups run nightly" but does not document restore testing — a critical gap. |
| ASF-015 | Backups are stored in a separate AWS region or account from the primary database. | Implicit | Backups co-located with the primary DB are vulnerable to the same region-wide failure or compromise. |
| ASF-016 | Backup files are encrypted at rest using a separate KMS key from the primary database key. | Explicit | Encryption at rest for backups is required but not explicitly stated in documented policy. |

---

### Pattern 6: Cloud Security (IAM)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-017 | The AWS account hosting RDS has no other workloads with unnecessary network access to the payroll VPC. | Environmental | Shared AWS accounts create cross-workload attack paths if VPC peering or transit gateway exists. |
| ASF-018 | IAM roles for RDS access are scoped to the minimum required actions. | Explicit | Over-permissioned IAM roles increase blast radius if the web application is compromised. |
| ASF-019 | The AWS root account is protected by MFA and not used for daily operations. | Derived | Root account compromise gives an attacker full control over the database snapshot, deletion, and access. |
| ASF-020 | No public AMIs or untrusted base images are used for the web application EC2 instance. | Implicit | Untrusted AMIs may contain backdoors or malware introduced at the image level. |

---

### Pattern 8: Data Flow & Classification

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-021 | Payroll data is classified as sensitive/confidential and subject to data handling policies. | Explicit | Security controls presume data sensitivity; without formal classification, controls may be inconsistent. |
| ASF-022 | Data flow diagrams exist and accurately represent all paths payroll data travels. | Implicit | Unmapped data flows (e.g., to monitoring, analytics, or third-party processors) create blind spots. |
| ASF-023 | The web application does not transmit payroll data to any endpoint outside the defined architecture. | Derived | The documented flow shows only App → DB; any other egress (logging service, cache, API) is unaccounted. |
| ASF-024 | Payroll data is not used in development or staging environments. | Environmental | Non-production use of production payroll data violates data classification and increases exposure surface. |

---

### Pattern 9: Encryption at Rest

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-025 | RDS encryption at rest is enabled using AWS KMS. | Explicit | Standard requirement for sensitive data on RDS. |
| ASF-026 | KMS key policies restrict which principals can encrypt/decrypt the database. | Derived | Encryption without access control on the key provides no real protection against authorized IAM users. |
| ASF-027 | KMS keys are rotated annually. | Operational | Key rotation limits the window of exposure if a KMS key is compromised. |
| ASF-028 | Temporary storage (swap, temp tables, query cache) on the RDS instance is also encrypted. | Implicit | Encryption at rest typically covers persistent storage; temp files may be written unencrypted by default. |

---

### Pattern 10: Encryption in Transit (TLS)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-029 | TLS is enforced between the VPN gateway and the web application as documented. | Explicit | Directly from documented topology. |
| ASF-030 | The web application validates the database server TLS certificate when connecting over SQL. | Trust | Without certificate validation, a MITM on the internal network can intercept database credentials. |
| ASF-031 | TLS 1.2 or higher is enforced; TLS 1.0/1.1 and all SSL versions are disabled at both the VPN and the application. | Derived | Older TLS versions have known cryptographic attacks that allow passive decryption. |
| ASF-032 | Weak cipher suites (RC4, 3DES, CBC-mode ciphers) are explicitly disabled on all TLS endpoints. | Derived | Strong TLS version with weak cipher negotiation is still vulnerable (e.g., Lucky13, Sweet32). |

---

### Pattern 11: Endpoint Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-033 | User laptops have EDR/AV software installed, running, and receiving current signature updates. | Implicit | An infected laptop connected via VPN can be used to pivot to internal resources. |
| ASF-034 | User laptops are managed by MDM with enforced disk encryption and screen lock policies. | Derived | Unmanaged endpoints have unknown security posture and cannot be trusted on the corporate network. |
| ASF-035 | Lost or stolen laptops can be remotely wiped to prevent VPN credentials or cached data from being recovered. | Operational | VPN client configuration files and cached AD tickets on a stolen laptop permit network access. |
| ASF-036 | Users do not install unauthorized software on corporate laptops. | Environmental | Unsanctioned software introduces vulnerabilities that EDR may not be configured to detect. |

---

### Pattern 12: Human Factors & Process

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-037 | Users do not share their VPN credentials or leave authenticated sessions logged in on shared workstations. | Derived | Credential sharing undermines individual accountability and access control. |
| ASF-038 | Users can identify and report phishing attempts targeting their AD credentials. | Trust | Phishing is the primary vector for credential theft; user detection is the last line of defense. |
| ASF-039 | Administrators assigning permissions follow the principle of least privilege. | Implicit | Without training and enforcement, administrators default to granting more access than necessary. |
| ASF-040 | Users whose role does not require payroll access are not granted access even if they are AD-authenticated. | Operational | Application-level access decisions depend on administrators correctly provisioning application roles. |

---

### Pattern 13: Identity Lifecycle (Provisioning)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-041 | User accounts follow a documented joiner/mover/leaver process in AD. | Operational | Stale accounts are the most common identity risk; without lifecycle management, accounts persist indefinitely. |
| ASF-042 | AD group membership for VPN and application access is reviewed and recertified quarterly. | Derived | Group membership drift leads to privilege creep over time. |
| ASF-043 | Service accounts used by the web application for database authentication are managed with the same rigor as user accounts. | Implicit | Orphaned or over-permissioned service accounts are frequently overlooked in access reviews. |
| ASF-044 | Application-level roles are synchronized with HR data so role changes automatically update access. | Environmental | Manual role provisioning lags behind HR events, creating windows of excessive or insufficient access. |

---

### Pattern 14: Incident Response

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-045 | There is an incident response plan that covers payroll data exposure scenarios. | Operational | Without a plan, response is ad-hoc and delays containment during a critical data breach. |
| ASF-046 | The incident response team has access to VPN, application, and database logs during an investigation. | Derived | Inaccessible logs prevent root cause analysis and attacker attribution. |
| ASF-047 | The IR plan includes isolation procedures (disconnect VPN, block DB traffic) that preserve forensic evidence. | Trust | Hasty isolation (e.g., powering down a DB server) destroys in-memory evidence of attacker activity. |
| ASF-048 | Monitoring systems can detect anomalies in VPN and database access patterns that indicate a breach. | Implicit | Detection is a prerequisite for incident response; undetected breaches have no response. |

---

### Pattern 15: Least Privilege

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-049 | The application database user has SELECT/INSERT/UPDATE/ DELETE only on the payroll schema and nothing else. | Explicit | Least privilege for database access is required but not explicitly documented. |
| ASF-050 | The web application does not run as root or with OS-level administrative privileges. | Derived | A web application running with elevated OS privileges magnifies the impact of any code execution vulnerability. |
| ASF-051 | AD users who can authenticate to the VPN cannot also directly access the database (no overlapping credentials). | Implicit | Users with both VPN and direct DB access can bypass application-layer controls and audit. |
| ASF-052 | The web application enforces authorization decisions beyond the initial AD authentication (application-level RBAC). | Derived | AD auth only confirms identity, not authorization; the application must independently enforce what each user can do. |

---

### Pattern 16: Monitoring & Alerting

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-053 | VPN connection logs are monitored for brute-force attempts and anomalous source geographies. | Operational | VPN brute-force is a common initial access technique; without monitoring, it proceeds undetected. |
| ASF-054 | Database query patterns are monitored for unusual volume (e.g., mass SELECT at 3 AM). | Derived | Bulk data exfiltration is visible in query volume anomalies if monitoring is configured. |
| ASF-055 | Alerts are configured for failed authentication thresholds on VPN, application, and AD. | Operational | High failed-auth rates indicate credential stuffing or brute-force in progress. |
| ASF-056 | Monitoring infrastructure logs are append-only and tamper-proof. | Implicit | Attackers who compromise monitoring can hide their activity by altering or deleting logs. |

---

### Pattern 17: Network Segmentation

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-057 | The VPN gateway, web application, and database are in separate security groups/subnets with explicit allow rules. | Architectural | Flat network architecture allows lateral movement; segmentation is required for defense-in-depth. |
| ASF-058 | There is no direct network path from the VPN subnet to the database subnet (traffic must pass through the web application). | Architectural | A direct path bypasses application-layer controls and audit. |
| ASF-059 | VPC flow logs or equivalent network telemetry is enabled to detect unexpected traffic patterns. | Operational | Without flow logs, unauthorized lateral movement is invisible. |
| ASF-060 | The database security group allows inbound traffic only from the web application security group on port 5432/3306. | Explicit | Restricting ingress to a specific source security group limits the blast radius of a compromised non-application resource. |

---

### Pattern 20: Third-party Dependency

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-061 | AWS RDS is available and not experiencing a region-level service outage. | Dependency | The architecture has no documented multi-region fallback; an RDS outage means complete data unavailability. |
| ASF-062 | The VPN vendor has no known backdoors or critical vulnerabilities that are unpatched. | Dependency | A compromised VPN appliance exposes all network traffic traversing it. |
| ASF-063 | Third-party libraries used by the web application are scanned for vulnerabilities before deployment. | Operational | Dependency vulnerabilities (e.g., Log4j) in the web application can lead to remote code execution. |
| ASF-064 | There is an exit strategy or migration plan if the VPN vendor or RDS becomes unavailable due to acquisition, bankruptcy, or sanctions. | Derived | Vendor discontinuity forces emergency migration that may bypass security review. |

**Total (A): 64** (4 per pattern × 15 patterns + 4 from Backup)

*Note: Backup & Recovery generated 4 assumptions (ASF-013 through ASF-016), making the total 64.*

---

## Step 3: Comparison

### Overlap Mapping

| Human ID | ASF ID | Match Rationale |
|----------|--------|-----------------|
| H-003 | ASF-037 | Both address credential sharing / shared credentials. |
| H-004 | ASF-001 | Both require MFA on VPN. |
| H-005 | ASF-057 | Both concern network isolation / segmentation of private subnet. |
| H-006 | ASF-030 | Both address TLS certificate validity and validation. |
| H-007 | ASF-031 | Both address modern encryption and disabling weak protocols. |
| H-009 | ASF-043 | Both address management of service/DB credentials (hardcoded or orphaned). |
| H-010 | ASF-049 | Both require least-privilege database user permissions. |
| H-011 | ASF-030 | Both require proper TLS termination and certificate validation. |
| H-013 | ASF-058 | Both assume no direct network path from VPN to database. |
| H-014 | ASF-016 | Both require backup encryption at rest. |
| H-015 | ASF-014 | Both require regular restore testing. |
| H-016 | ASF-060 | Both require security groups restricting DB access to app only. |
| H-017 | ASF-052 | Both address input validation and application-level authorization (overlapping concern). |
| H-018 | ASF-052 | Both address session management and application-level access control. |
| H-019 | ASF-053 | Both require VPN log monitoring. |
| H-022 | ASF-041 | Both require account lifecycle management (termination → revocation). |
| H-023 | ASF-054 | Both require database audit logging and query monitoring. |
| H-024 | ASF-052 | Both require application-level RBAC beyond AD authentication. |
| H-025 | ASF-057 | Both require network segmentation / restricted administrative access to DB. |
| H-028 | ASF-031 | Both require TLS 1.2+. |
| H-029 | ASF-017 | Both assume no other workloads access the payroll database or its environment. |
| H-030 | ASF-057 | Both require firewall/security group separation between tiers. |
| H-036 | ASF-033 | Both require endpoint security posture checks before VPN connection. |
| H-037 | ASF-023 | Both assume payroll data does not leak to unauthorized endpoints (logs, errors). |
| H-040 | ASF-031 | Both address disabling weak authentication protocols and enforcing modern crypto. |

**Overlap (O): 25**

### Metrics

| Metric | Value | Formula |
|--------|-------|---------|
| Human assumptions (H) | 40 | Count of unique human-generated assumptions |
| ASF assumptions (A) | 64 | Count of unique ASF-generated assumptions |
| Overlap (O) | 25 | Count appearing in both lists |
| **Precision** | **39.1%** | O / A = 25/64 |
| **Recall** | **62.5%** | O / H = 25/40 |
| **F1 Score** | **48.1%** | 2 × (P × R) / (P + R) |
| Novel findings (A - O) | 39 | Assumptions ASF found that human missed (61.0% of ASF total) |
| Missed findings (H - O) | 15 | Assumptions human found that ASF missed (37.5% of human total) |

### Success Criteria Assessment

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Recall | >= 70% | 62.5% | ❌ Not met |
| Precision | >= 50% | 39.1% | ❌ Not met |
| Novel discoveries | >= 10% of total (A+O) | 37.5% (39/104) | ✅ Exceeded |
| Expert agreement (F1 proxy) | > 60% | 48.1% | ❌ Not met |

---

## Step 4: Analysis

### Ontology Overlap Analysis

The cross-tabulation from the ASF ontology (1697 assumptions) shows that **Derived** (34.3%), **Operational** (20.0%), and **Implicit** (16.6%) are the largest ontology categories. In this experiment, the distribution of overlapping assumptions tilted toward Explicit and Derived:

| Ontology Category | Overlaps | ASF Total | Overlap Rate |
|-------------------|----------|-----------|-------------|
| Explicit | 7 | 12 | 58.3% |
| Derived | 8 | 16 | 50.0% |
| Operational | 5 | 16 | 31.3% |
| Implicit | 3 | 8 | 37.5% |
| Trust | 1 | 4 | 25.0% |
| Dependency | 1 | 4 | 25.0% |
| Architectural | 0 | 4 | 0.0% |
| Environmental | 0 | 4 | 0.0% |

**Best overlap:** Explicit and Derived categories showed the strongest agreement. Both humans and the ASF immediately recognize that "MFA must be enforced" (Explicit) and "backups must be tested" (Derived) are critical.

**Worst overlap:** Architectural and Environmental categories had zero overlap. The ASF identified architectural concerns (VPN redundancy, network segmentation) and environmental concerns (shared AWS accounts, non-production data usage) that the human did not list as assumptions — possibly because a human architect treats these as contextual design decisions rather than hidden assumptions.

### What Humans Caught That ASF Missed (Missed Findings = 15)

The 15 human-generated assumptions with no ASF counterpart fall into three clusters:

1. **Application security implementation details (H-017, H-018, H-027, H-033, H-037):** The ASF patterns have no dedicated "Secure Coding" or "Web Application Security" pattern. SQL injection prevention, CSRF/XSS defenses, rate limiting, and error message sanitation are all invisible to the ASF's current 20-pattern matrix. These are implementation-level concerns rather than architectural or policy-level.

2. **VPN-specific operational hardening (H-002, H-026, H-032, H-039):** Split-tunneling prevention, concurrent connection capacity, kill switch behavior, and perfect forward secrecy are VPN product configuration details. The ASF treats VPN under "Network Segmentation" and "Encryption in Transit," but these VPN-specific operational settings fall through the cracks.

3. **Database platform security (H-031, H-038, H-044):** RDS patching cadence, deletion protection, and final snapshot policies are specific to managed database platforms. The ASF's "Cloud Security (IAM)" and "Backup & Recovery" patterns cover general cloud and backup concerns but miss database-platform-specific safety nets.

### What ASF Caught That Humans Missed (Novel Findings = 39)

The ASF generated 64 assumptions, of which 39 (61.0%) were not in the human list. The most significant novel categories:

1. **Incident Response (4 assumptions):** The human architect generated zero IR assumptions. The ASF contributed a full pattern on IR planning, log access, isolation procedures, and detection. This is the single largest gap in the human list.

2. **Identity Lifecycle (4 assumptions):** The human assumed termination revocation (H-022) but did not extend to joiner/mover/leaver process, quarterly recertification, or service account management. The ASF's identity lifecycle pattern surfaced these operational identity concerns.

3. **Operational resilience (ASF-009, ASF-010, ASF-011, ASF-012):** The human assumed VPN capacity (H-026) but did not consider VPN redundancy, documented offline procedures, ISP reliability, or graceful application degradation during database failures.

4. **Monitoring infrastructure security (ASF-056):** The human assumed logs are sent to SIEM (H-019) but did not consider that the monitoring infrastructure itself must be tamper-proof. This is a classic blind spot — securing the security tools.

5. **Third-party dependency risk (ASF-061 through ASF-064):** The human treated the architecture as self-contained. The ASF surfaced dependencies on AWS RDS availability, VPN vendor security posture, library vulnerability scanning, and vendor exit strategy — all risks outside the architecture diagram.

### Architecture Complexity Assessment

Architecture #1 was classified as **Simple** (linear, 4 nodes, 3 boundaries, well-understood components). However:

- The **human precision/recall gap** (62.5% recall, 39.1% precision) suggests that even for simple architectures, the ASF generates a much broader set of assumptions than a human architect working from first principles.
- The high **novelty rate (61.0% of ASF output)** indicates that even for simple architectures, the ASF adds substantial value beyond unaided human reasoning.
- The **missed rate (37.5%)** is above the 30% target, driven primarily by the absence of a "Web Application Security" or "Secure Coding" pattern in the ASF matrix.

The human architect focused heavily on the **attack path** (VPN → App → DB) and generated assumptions along that path. The ASF, by contrast, generated assumptions **orthogonal** to the data path — identity lifecycle, incident response, monitoring infrastructure, third-party dependencies — that a human might not consider unless prompted.

### Key Insight

The biggest root cause of the missed findings is **pattern coverage**: the ASF's 20-pattern matrix has strong coverage of infrastructure, network, cloud, and operational concerns but lacks explicit patterns for:
- **Web application security** (SQLi, XSS, CSRF, rate limiting)
- **Secure coding practices** (input validation, output encoding, error handling)
- **Platform-specific database security** (RDS deletion protection, patching, snapshots)

Adding a "Web Application Security" pattern (pattern 21) and a "Managed Database Security" sub-pattern under "Cloud Security (IAM)" would likely close the recall gap to above 70%.

---

## Conclusions

| Metric | Target | Actual | Verdict |
|--------|--------|--------|---------|
| Recall | >= 70% | 62.5% | ❌ Below target — missing web app security pattern |
| Precision | >= 50% | 39.1% | ❌ Below target — ASF is deliberately broad/generative |
| Novel discoveries | >= 10% | 37.5% | ✅ ASF adds substantial value beyond human reasoning |
| Expert agreement (F1) | > 60% | 48.1% | ❌ Below target — driven by low precision |

The ASF framework applied to Architecture #1 demonstrates strong **exploration breadth** (finding assumptions the human missed) but lower **precision** (generating many assumptions the human did not consider relevant). For architecture-level risk identification, this trade-off is acceptable — false positives are preferable to false negatives. The primary actionable finding is the need for a **Web Application Security** pattern to close the recall gap for SQL injection, CSRF, XSS, and rate-limiting concerns.
