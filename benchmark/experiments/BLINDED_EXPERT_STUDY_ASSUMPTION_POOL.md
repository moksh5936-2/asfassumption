# Blinded Expert Study: Assumption Pool

**150 assumptions across 3 buckets:** ASF-Only (50), Claude (50), Gemini+Gemma (50)
**Blinded IDs:** A-001 through A-150 (randomized order in actual instrument)
**Source:** 5 reference architectures (Arch 001, 004, 005, 011, 020)

---

## Bucket 1: ASF-Only (50 assumptions)

*Assumptions that NONE of the 4 AI models (Claude, GPT, Gemini, Gemma) independently derived.*

### Architecture 1: VPN → Payroll DB

| ID | Assumption |
|----|-----------|
| ASF-01 | VPN gateway redundancy must exist — a single VPN gateway failure blocks all remote access |
| ASF-02 | There is a documented offline procedure for VPN outages that does not bypass security controls |
| ASF-03 | The internet circuit to the VPN gateway has sufficient bandwidth and reliability SLA |
| ASF-04 | Payroll data is formally classified as sensitive/confidential with corresponding data handling policies |
| ASF-05 | Data flow diagrams exist and accurately represent all paths payroll data travels |
| ASF-06 | The web application does not transmit payroll data to any endpoint outside the defined architecture |
| ASF-07 | KMS keys for RDS encryption are rotated annually |
| ASF-08 | KMS key policies restrict which principals can encrypt/decrypt the database |
| ASF-09 | Temporary storage (swap, temp tables, query cache) on the RDS instance is also encrypted |
| ASF-10 | Database query patterns are monitored for unusual volume (e.g., mass SELECT at 3 AM) |
| ASF-11 | Monitoring infrastructure logs are append-only and tamper-proof |
| ASF-12 | VPC flow logs or equivalent network telemetry is enabled to detect unexpected traffic patterns |
| ASF-13 | The VPN vendor has no known backdoors or critical vulnerabilities that are unpatched |
| ASF-14 | There is an exit strategy or migration plan if the VPN vendor becomes unavailable due to acquisition, bankruptcy, or sanctions |
| ASF-15 | User accounts follow a documented joiner/mover/leaver process in AD |
| ASF-16 | AD group membership for VPN and application access is reviewed and recertified quarterly |

### Architecture 2: SSO/IdP → SAML Federation

| ID | Assumption |
|----|-----------|
| ASF-17 | There is a documented offline procedure for Okta outages that does not bypass security controls |
| ASF-18 | AD domain controllers are redundant and can survive a single DC failure |
| ASF-19 | Internet connectivity to Okta is redundant (dual ISPs, different carriers) |
| ASF-20 | Okta tenant configuration is exportable and backed up externally |
| ASF-21 | AD restore procedures are tested at least annually |
| ASF-22 | SAML assertions are not logged in plaintext by SPs or Okta event logs |
| ASF-23 | There is no hidden data flow (e.g., SPs calling back to AD directly, bypassing Okta) |
| ASF-24 | Okta's data-at-rest encryption is verified through SOC 2 or equivalent reports |
| ASF-25 | Service accounts used for Okta-AD integration, API tokens, and SP integration are managed with the same rigor as human accounts |
| ASF-26 | There is an exit strategy or migration plan if Okta service is terminated |

### Architecture 3: K8s/Istio Service Mesh

| ID | Assumption |
|----|-----------|
| ASF-27 | K8s API Server access requires MFA for all human administrators |
| ASF-28 | MFA is enforced for Istio control plane configuration changes |
| ASF-29 | ServiceAccount token usage does not bypass MFA requirements for human operators |
| ASF-30 | Istio control plane (Pilot/Citadel) is highly available with multiple replicas |
| ASF-31 | Envoy sidecar proxy health does not degrade under mesh-wide configuration pushes |
| ASF-32 | Istio configuration (CRDs) is backed up independently of etcd |
| ASF-33 | CA certificates and keys from Citadel are backed up securely for recovery |
| ASF-34 | CloudTrail or equivalent is enabled to detect unauthorized cloud API calls from compromised nodes |
| ASF-35 | Data flowing between services is classified and handling requirements are documented |
| ASF-36 | K8s Secret objects are encrypted at rest using encryption-at-rest configuration (not just base64) |
| ASF-37 | Workload identity (SPIFFE) is revoked when a service is decommissioned |
| ASF-38 | etcd snapshots taken during incident response are forensically isolated |
| ASF-39 | No service account has cluster-admin privileges |
| ASF-40 | There is a documented migration path if Istio or K8s becomes unavailable due to licensing or deprecation |

### Architecture 4: Healthcare → PHI → HIPAA

| ID | Assumption |
|----|-----------|
| ASF-41 | Auth0 anomaly detection flags MFA fatigue attacks (repeated push notifications) |
| ASF-42 | An Auth0 outage has a documented fallback procedure that does not bypass security controls |
| ASF-43 | Changes to Auth0 tenant configuration (connections, rules, MFA policies) follow a documented change process |
| ASF-44 | Database schema changes are reviewed for PHI access implications before deployment |
| ASF-45 | PHI data is not used in development or staging environments — synthetic data is used instead |
| ASF-46 | The IR plan includes a specific playbook for PHI breach notification (HIPAA 60-day rule) |
| ASF-47 | The IR team can isolate the PHI database (block network access, disable application) to contain a breach |
| ASF-48 | The IR team has access to Auth0 logs for investigating authentication-related PHI breaches |
| ASF-49 | Provider accounts follow a documented joiner/mover/leaver process in Auth0 |
| ASF-50 | Auth0 tenant can be migrated to another identity provider if Auth0 becomes unavailable or non-compliant |

---

## Bucket 2: Claude (50 assumptions)

*Assumptions from Claude's independent derivation outputs, selected for architecture specificity.*

### Architecture 1: VPN → Payroll DB

| ID | Assumption |
|----|-----------|
| C-01 | The VPN gateway properly validates client certificates and doesn't allow unauthorized VPN connections |
| C-02 | The VPN configuration prevents split tunneling, ensuring all traffic goes through the VPN tunnel when connected |
| C-03 | The VPN gateway implements rate limiting to prevent brute force attacks against authentication |
| C-04 | The TLS connection between VPN gateway and Internal Web App uses TLS 1.2+ with strong cipher suites |
| C-05 | The Internal Web Application uses parameterized queries or an ORM to prevent SQL injection |
| C-06 | The web application implements secure session management with strong session IDs, proper timeout, and secure cookie flags |
| C-07 | There is comprehensive monitoring and logging at each trust boundary with alerting for suspicious activities |
| C-08 | All components receive regular security patches for known vulnerabilities |
| C-09 | There is a documented and tested incident response plan for breaches involving payroll data |
| C-10 | Administrative access to systems follows separation of duties principles |

### Architecture 2: SSO/IdP → SAML Federation

| ID | Assumption |
|----|-----------|
| C-11 | Okta properly manages SAML signing/encryption certificates and private keys with regular rotation and secure storage |
| C-12 | Active Directory is secured with least privilege access, monitoring, and protection against common attacks (DCShadow, Golden Ticket) |
| C-13 | Just-In-Time provisioning only creates users with appropriate, limited permissions and doesn't grant excessive privileges |
| C-14 | Attributes passed in SAML assertions from AD via Okta are trustworthy and accurately represent user entitlements |
| C-15 | MFA is enforced for every authentication attempt including legacy protocols, API access, and fallback mechanisms |
| C-16 | The organization would detect a compromise of Okta IdP through monitoring, logs, or anomaly detection |
| C-17 | Single Logout properly propagates to all connected applications and session revocation works effectively |
| C-18 | SAML bindings are implemented securely without XML External Entity attacks or XML Signature wrapping vulnerabilities |

### Architecture 3: K8s/Istio Service Mesh

| ID | Assumption |
|----|-----------|
| C-19 | mTLS certificates are rotated before expiry and services can reload new certificates without downtime |
| C-20 | Compromised certificates are effectively revoked and the revocation is honored by all services |
| C-21 | Citadel CA's root signing key is securely protected from compromise |
| C-22 | The Istio control plane is isolated from the data plane and cannot be compromised through data plane attacks |
| C-23 | etcd data is encrypted at rest to protect cluster secrets and state |
| C-24 | Kubernetes worker nodes are hardened to prevent container-to-host escapes |
| C-25 | Pod security policies prevent privilege escalation, host network access, and host PID sharing |
| C-26 | Service account tokens are automatically rotated with avoidance of long-lived tokens |
| C-27 | Network policies follow default-deny with explicit necessary allowances between namespaces |
| C-28 | Egress traffic is restricted to approved external endpoints only |

### Architecture 4: Healthcare → PHI → HIPAA

| ID | Assumption |
|----|-----------|
| C-29 | Network segmentation between the patient portal, app server, and PHI database is properly implemented |
| C-30 | Audit logs are immutable, tamper-evident, and cannot be altered by attackers or privileged insiders |
| C-31 | The App Server enforces minimum necessary access at the data level, not just the API level |
| C-32 | Encryption keys for AES-256 PHI encryption are stored in a separate KMS or HSM, not alongside the encrypted data |
| C-33 | All inter-component communication uses TLS 1.2+ with proper certificate validation |
| C-34 | Audit log forwarding to the SIEM is reliable and complete — no events are lost during network disruptions |
| C-35 | BAAs with subprocessors include enforceable security requirements and are actively monitored for compliance |
| C-36 | Security monitoring and alerting is effective — the SIEM has relevant rules and alerts are responded to |

### Architecture 5: ERP → SOX → Audit

| ID | Assumption |
|----|-----------|
| C-37 | Audit logs are immutable, comprehensive, and protected from tampering or deletion by any user including DBAs |
| C-38 | All system components maintain synchronized time for accurate audit trails and transaction ordering |
| C-39 | Security events are continuously monitored with alerts for suspicious activities across all trust boundaries |
| C-40 | Administrative and privileged access to the ERP is strictly controlled, monitored, and subject to segregation of duties |
| C-41 | Sensitive financial data is encrypted at rest in the Financial DB and other storage systems |
| C-42 | Segregation of duties controls are properly implemented at the software level and cannot be circumvented |
| C-43 | Database credentials, API keys, and secrets for the ERP are securely managed in a vault, not hardcoded |
| C-44 | The ERP system has a tested disaster recovery plan that maintains security controls during failover |
| C-45 | All APIs connecting ERP components are secured with authentication, authorization, and input validation |
| C-46 | Error messages from the ERP do not leak sensitive information such as database structure or stack traces |
| C-47 | Third-party integrations with the ERP (bank feeds, payment gateways) maintain equivalent security controls |
| C-48 | All users receive regular security training and understand their SOX compliance responsibilities |
| C-49 | All changes to the ERP follow strict change control procedures that preserve security controls |
| C-50 | File upload/download functionality in the ERP is secured against malicious uploads and path traversal |

---

## Bucket 3: Gemini+Gemma (50 assumptions)

*Assumptions from Gemini's and Gemma's independent derivation outputs.*

### Gemini — Architecture 1: VPN → Payroll DB

| ID | Assumption |
|----|-----------|
| G-01 | The user's laptop connecting via VPN must be uncompromised — an infected device with malware can hijack the authenticated VPN session |
| G-02 | Network traffic inside the private subnet cannot be sniffed — the SQL connection between the web app and database may not be encrypted |
| G-03 | Active Directory credentials must be resistant to brute-force and credential stuffing — the architecture does not specify MFA for the application layer |
| G-04 | Database backups must be securely isolated and encrypted — the policy states backups run nightly but not where they go or who has access |
| G-05 | DBAs and cloud administrators must not be able to read payroll data by bypassing the application logic |

### Gemini — Architecture 2: SSO/IdP → SAML Federation

| ID | Assumption |
|----|-----------|
| G-06 | Service Provider applications must correctly validate SAML assertions and signatures — federation relies entirely on SP trusting IdP key pairs |
| G-07 | The connection between Okta and Active Directory must be secure and authenticated — a compromised directory sync compromises the user store |
| G-08 | JIT provisioning must properly de-provision accounts when roles change — orphaned accounts remain active inside SP apps |
| G-09 | User browsers must be free of session-hijacking malware — an 8-hour session token sits in browser storage vulnerable to infostealers |

### Gemini — Architecture 3: K8s/Istio Service Mesh

| ID | Assumption |
|----|-----------|
| G-10 | The Citadel CA and its private keys must be completely secure — a compromised CA can mint valid certificates for any service |
| G-11 | The K8s control plane (etcd and API Server) must be isolated from the data plane — a compromised pod with API access can rewrite network policies |
| G-12 | The container runtime and host OS kernel must be secure with no container escape vulnerabilities |
| G-13 | Persistent Volume data for the StatefulSet DB must be securely encrypted and isolated at the storage layer |

### Gemini — Architecture 4: Healthcare → PHI → HIPAA

| ID | Assumption |
|----|-----------|
| G-14 | The App Server must enforce authorization on every API request — Auth0 handles authentication but not authorization (no BOLA/IDOR) |
| G-15 | Cryptographic keys for AES-256 encryption must be stored and managed outside the database |
| G-16 | Audit logs sent to the SIEM must be write-once/immutable and cannot be modified by the compromised App Server |
| G-17 | Auth0 tenant administrative accounts must be securely managed — a compromised Auth0 admin portal can alter authentication flows |

### Gemini — Architecture 5: ERP → SOX → Audit

| ID | Assumption |
|----|-----------|
| G-18 | The Approval Workflow engine cannot be bypassed via direct database or backend API access |
| G-19 | The Reporting Engine cannot be exploited for SQL injection or privilege escalation — auditors have read-only access |
| G-20 | System administrators and DBAs must be subject to the same segregation of duties as the finance team |
| G-21 | System clocks across ERP, Approval Workflow, and Audit Logs must be synchronized via NTP for audit integrity |

### Gemma — Architecture 1: VPN → Payroll DB

| ID | Assumption |
|----|-----------|
| G-22 | The VPN gateway must validate the security posture of the connecting device (host checker) before allowing access |
| G-23 | Active Directory must be protected against lateral movement and domain dominance attacks beyond the VPN boundary |
| G-24 | The Internal Web App must be free from OWASP Top 10 application-layer vulnerabilities (SQLi, RCE) |
| G-25 | The SQL traffic between the web app and RDS must be isolated from other internal assets to prevent sniffing |

### Gemma — Architecture 2: SSO/IdP → SAML Federation

| ID | Assumption |
|----|-----------|
| G-26 | SAML signing private keys held by Okta must remain strictly confidential — leaked keys allow forging assertions without hitting Okta |
| G-27 | The 8-hour session token cannot be exfiltrated or reused from a different network context |
| G-28 | JIT provisioning must include automated attribute verification — manipulated AD attributes can create privileged accounts in SP apps |
| G-29 | MFA must be resistant to bypass via push fatigue / MFA spamming attacks |

### Gemma — Architecture 3: K8s/Istio Service Mesh

| ID | Assumption |
|----|-----------|
| G-30 | Linux kernel namespace isolation must be robust — a container escape can bypass all network policies |
| G-31 | Applications must not log or leak mTLS session tokens or decrypted payloads in debug output |
| G-32 | Access to the etcd key-value store must be strictly limited and encrypted — etcd holds all cluster secrets |
| G-33 | The Ingress Gateway must sanitize external traffic before routing — internal services trust mTLS and may skip input validation |

### Gemma — Architecture 4: Healthcare → PHI → HIPAA

| ID | Assumption |
|----|-----------|
| G-34 | Auth0 configuration must prevent token substitution and replay attacks |
| G-35 | The App Server must not cache unencrypted PHI to local disk, temp files, swap, or crash dumps |
| G-36 | The connection between the App Server and SIEM must be reliable and tamper-proof — a compromised app server can block log delivery |
| G-37 | Insider threats (malicious developers/DevOps) must not be able to access database encryption keys from production configuration |

### Gemma — Architecture 5: ERP → SOX → Audit

| ID | Assumption |
|----|-----------|
| G-38 | The Reporting Engine must have read-only database permissions at the database level, not just configured in the application |
| G-39 | User identity inside the approval workflow cannot be spoofed — session fixation or ID hijacking could allow self-approval |
| G-40 | Quarterly recertification data must come from a tamper-proof system of record — not reports generated by the administrators being reviewed |
| G-41 | The ERP Web App must securely handle session state and not suffer from privilege escalation bugs |
| G-42 | A low-level accountant must not be able to modify session parameters to access admin functions or bypass approval chains |

### Additional from Gemma (cross-architecture)

| ID | Assumption |
|----|-----------|
| G-43 | Database credentials must be stored in a secrets manager, not in application configuration files or environment variables |
| G-44 | The application must not have local file caching or temporary storage of unencrypted sensitive data |
| G-45 | Session tokens must be protected against replay from different network locations |
| G-46 | Security configurations must be monitored for drift over time |
| G-47 | The principle of least privilege must extend to all service accounts, not just human users |
| G-48 | Backup files must be stored in a separate location from the primary data with independent access controls |
| G-49 | Administrative access must require a separate, non-admin account for daily work |
| G-50 | All external integrations must be reviewed for security before deployment |

---

## Next Step: Build the Blinded Instrument

1. Randomize all 150 items (A-001 through A-150)
2. Create the rating form (Google Form or score sheet PDF)
3. The key mapping is only revealed after all data is collected
4. Score each bucket independently: mean total, pass rate, Q3 essential rate

## Analysis Plan

| Bucket | Expected Mean (hypothesis) |
|--------|---------------------------|
| ASF-Only (50) | ≥ 12/20 (moderately valuable) |
| Claude (50) | ≥ 13/20 (baseline for valid assumptions) |
| Gemini+Gemma (50) | ≥ 12/20 (shorter but specific) |

**Success:** ASF-Only is non-inferior to Claude (±2 pts) with ≥ 40% Q3 essential rate.
**Strong success:** ASF-Only meets or exceeds Claude mean.
**Failure:** ASF-Only scores < 10/20 or < 30% Q3 essential rate.
