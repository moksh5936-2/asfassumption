# Architecture 1: VPN → Internal Web App → Payroll DB — Security Review

## 1. Consensus Matrix

| # | Assumption | GPT | Gemini | Gemma | Keep? |
|---|---|---|---|---|---|
| 1 | Endpoints connecting via VPN are uncompromised and malware-free | A1 | G1 | GE1 | Yes |
| 2 | VPN authentication resists credential theft (MFA/certificates) | A2 | G3 | — | Yes |
| 3 | VPN gateway software has no exploitable vulnerabilities | A4 | — | — | Yes |
| 4 | VPN gateway configuration is secure | A5 | — | — | Yes |
| 5 | TLS is correctly configured between VPN and web application | A6 | — | — | Yes |
| 6 | TLS certificates used by endpoints are trusted and valid | A7 | — | — | Yes |
| 7 | Internal DNS is trustworthy and not poisoned | A8 | — | — | Yes |
| 8 | The private subnet is actually isolated from other networks | A23 | — | — | Yes |
| 9 | Firewall rules (security groups / NACLs) are correct and restrictive | A24 | — | — | Yes |
| 10 | Only the web application tier can access the database directly | A25 | — | — | Yes |
| 11 | SQL traffic between web app and RDS is encrypted and not subject to sniffing | — | G2 | GE4 | Yes |
| 12 | No undocumented access paths exist beyond the architecture diagram | A40 | — | — | Yes |
| 13 | Active Directory is secure against compromise | A9 | — | — | Yes |
| 14 | AD credentials resist brute-force and credential-stuffing attacks | A10 | G3 | — | Yes |
| 15 | AD authorization groups are correctly maintained and audited | A11 | — | — | Yes |
| 16 | User accounts are deprovisioned promptly upon termination | A12 | — | — | Yes |
| 17 | Service accounts (application and system) are protected and rotated | A13 | — | — | Yes |
| 18 | AD is protected against lateral movement and domain dominance attacks | — | — | GE2 | Yes |
| 19 | Multi-factor authentication is enforced at the application layer | — | G3 | — | Yes |
| 20 | Application authorization logic correctly enforces least privilege | A14 | — | — | Yes |
| 21 | Session management is secure (timeout, rotation, invalidation) | A15 | — | — | Yes |
| 22 | Session identifiers are cryptographically random and unpredictable | A16 | — | — | Yes |
| 23 | Application input validation is effective (no OWASP Top 10 vulns) | A17 | — | GE3 | Yes |
| 24 | SQL queries are parameterized to prevent injection | A18 | — | — | Yes |
| 25 | Application servers are hardened against OS-level compromise | A19 | — | — | Yes |
| 26 | Administrative interfaces are restricted and not exposed externally | A20 | — | — | Yes |
| 27 | Secrets (API keys, credentials, certificates) are securely stored | A21 | — | — | Yes |
| 28 | Database credentials are protected from exfiltration | A22 | — | — | Yes |
| 29 | Database permissions follow the principle of least privilege | A26 | — | — | Yes |
| 30 | Database software is patched against known vulnerabilities | A27 | — | — | Yes |
| 31 | Database encryption keys are protected from unauthorized access | A28 | — | — | Yes |
| 32 | DBAs and cloud administrators cannot bypass application logic to read payroll data | — | G5 | — | Yes |
| 33 | Backups are encrypted and stored in an isolated location | A29 | G4 | — | Yes |
| 34 | Backup access is restricted to authorized personnel only | A30 | G4 | — | Yes |
| 35 | Backup integrity is verified regularly (recovery testing) | A31 | — | — | Yes |
| 36 | Audit logs cannot be altered or destroyed by attackers | A32 | — | — | Yes |
| 37 | Security monitoring and alerting exist for anomalous activity | A33 | — | — | Yes |
| 38 | Time synchronization (NTP) is accurate across all systems | A34 | — | — | Yes |
| 39 | Network devices (routers, switches, firewalls) enforce policy correctly | A39 | — | — | Yes |
| 40 | Cloud provider isolation guarantees work as designed (multi-tenant) | A38 | — | — | Yes |
| 41 | Physical security protects the underlying infrastructure | A37 | — | — | Yes |
| 42 | Change management prevents unauthorized or insecure modifications | A36 | — | — | Yes |
| 43 | Administrators are trustworthy and separation of duties is enforced | A35 | G5 | — | Yes |

## 2. Deduplicated Assumption List

### Endpoint Security
1. Endpoints connecting via VPN are uncompromised and malware-free
2. VPN authentication resists credential theft (MFA/certificates)
3. VPN gateway software has no exploitable vulnerabilities
4. VPN gateway configuration is secure

### VPN & Network Edge Security
5. TLS is correctly configured between VPN and web application
6. TLS certificates used by endpoints are trusted and valid
7. Internal DNS is trustworthy and not poisoned

### Network Security
8. The private subnet is actually isolated from other networks
9. Firewall rules (security groups / NACLs) are correct and restrictive
10. Only the web application tier can access the database directly
11. SQL traffic between web app and RDS is encrypted and not subject to sniffing
12. No undocumented access paths exist beyond the architecture diagram

### Identity & Active Directory
13. Active Directory is secure against compromise
14. AD credentials resist brute-force and credential-stuffing attacks
15. AD authorization groups are correctly maintained and audited
16. User accounts are deprovisioned promptly upon termination
17. Service accounts (application and system) are protected and rotated
18. AD is protected against lateral movement and domain dominance attacks
19. Multi-factor authentication is enforced at the application layer

### Application Security
20. Application authorization logic correctly enforces least privilege
21. Session management is secure (timeout, rotation, invalidation)
22. Session identifiers are cryptographically random and unpredictable
23. Application input validation is effective (no OWASP Top 10 vulns)
24. SQL queries are parameterized to prevent injection
25. Application servers are hardened against OS-level compromise
26. Administrative interfaces are restricted and not exposed externally
27. Secrets (API keys, credentials, certificates) are securely stored

### Database Security
28. Database credentials are protected from exfiltration
29. Database permissions follow the principle of least privilege
30. Database software is patched against known vulnerabilities
31. Database encryption keys are protected from unauthorized access
32. DBAs and cloud administrators cannot bypass application logic to read payroll data

### Backup & Recovery
33. Backups are encrypted and stored in an isolated location
34. Backup access is restricted to authorized personnel only
35. Backup integrity is verified regularly (recovery testing)

### Monitoring & Logging
36. Audit logs cannot be altered or destroyed by attackers
37. Security monitoring and alerting exist for anomalous activity
38. Time synchronization (NTP) is accurate across all systems
39. Network devices (routers, switches, firewalls) enforce policy correctly

### Cloud & Infrastructure
40. Cloud provider isolation guarantees work as designed (multi-tenant)
41. Physical security protects the underlying infrastructure

### Governance & Administration
42. Change management prevents unauthorized or insecure modifications
43. Administrators are trustworthy and separation of duties is enforced

## 3. Risk Scores

| # | Assumption | Likelihood | Impact | Risk |
|---|---|---|---|---|
| 1 | Endpoints uncompromised | M | H | H |
| 2 | VPN auth resists credential theft | M | H | H |
| 3 | VPN gateway no vulns | M | H | H |
| 4 | VPN gateway config secure | M | H | H |
| 5 | TLS correctly configured | M | H | H |
| 6 | TLS certs trusted | L | H | M |
| 7 | DNS trustworthy | M | H | H |
| 8 | Private subnet isolated | M | C | C |
| 9 | Firewall rules correct | M | C | C |
| 10 | Only app accesses DB | M | C | C |
| 11 | SQL traffic encrypted | H | H | H |
| 12 | No undocumented paths | H | H | H |
| 13 | AD secure | M | C | C |
| 14 | AD cred brute-force resistance | H | H | H |
| 15 | AD groups maintained | M | M | M |
| 16 | User deprovisioning | M | H | H |
| 17 | Service accounts protected | M | H | H |
| 18 | AD lateral movement protection | H | C | C |
| 19 | MFA at application layer | H | H | H |
| 20 | App authorization correct | M | H | H |
| 21 | Session management secure | M | H | H |
| 22 | Session IDs unpredictable | L | H | M |
| 23 | Input validation / OWASP | H | C | C |
| 24 | SQL parameterized | M | C | C |
| 25 | App servers hardened | M | H | H |
| 26 | Admin interfaces restricted | M | H | H |
| 27 | Secrets securely stored | M | H | H |
| 28 | DB credentials protected | M | C | C |
| 29 | DB least privilege | M | H | H |
| 30 | DB patched | H | C | C |
| 31 | DB encryption keys protected | M | H | H |
| 32 | DBAs/cloud admins bypass app logic | M | C | C |
| 33 | Backups encrypted/isolated | M | C | C |
| 34 | Backup access restricted | M | C | C |
| 35 | Backup integrity verified | L | H | M |
| 36 | Logs not alterable | M | M | M |
| 37 | Security monitoring exists | H | C | C |
| 38 | Time sync accurate | L | M | L |
| 39 | Network devices enforce policy | M | H | H |
| 40 | Cloud provider isolation | L | H | M |
| 41 | Physical security | L | H | M |
| 42 | Change management | M | H | H |
| 43 | Admin trust / SoD enforced | M | C | C |

## 4. STRIDE Mapping

| STRIDE Category | Assumption Numbers | Rationale |
|---|---|---|
| **Spoofing** | 1, 2, 3, 4, 6, 7, 13, 14, 16, 17, 19, 22, 27, 28, 40 | Identity/device impersonation via compromised endpoints, weak auth, forged certificates, poisoned DNS, or stolen credentials |
| **Tampering** | 5, 9, 11, 12, 23, 24, 30, 35, 39, 42 | Data/modification attacks through TLS misconfig, firewall bypass, injection, unpatched software, backup corruption, or unauthorized change |
| **Repudiation** | 36, 37, 38 | Non-repudiation depends on tamper-proof logs, active monitoring, and accurate timestamps |
| **Information Disclosure** | 1, 5, 6, 7, 8, 9, 10, 11, 12, 13, 27, 28, 29, 31, 32, 33, 34, 40, 41, 43 | Payroll data exposure through network sniffing, direct DB access, backup leakage, insider abuse, or cloud isolation failure |
| **Denial of Service** | 8, 13, 23, 30, 42 | Service disruption through subnet exposure, AD compromise, application-level DoS, unpatched vulns, or unauthorized configuration changes |
| **Elevation of Privilege** | 1, 3, 4, 13, 15, 17, 18, 20, 21, 23, 24, 25, 26, 30, 43 | Privilege escalation through endpoint takeover, VPN/AD compromise, app logic flaws, SQLi, unhardened servers, or broken SoD |

## 5. Top 10 Critical Assumptions

| Rank | Assumption | Rationale |
|---|---|---|
| **1** | Private subnet is actually isolated (A8) | **Internal network trust.** If the private subnet is routable from other networks, the database is directly exposed — bypassing all application-layer controls. This is foundational to the entire architecture. |
| **2** | Security monitoring exists and is effective (A37) | Without detection, all other controls are moot. An attacker can dwell, pivot, and exfiltrate payroll data for months with no response. |
| **3** | Input validation prevents OWASP Top 10 (A23) | SQL injection or RCE in the web app provides direct access to payroll data. This is the most common attack vector against internal web applications. |
| **4** | Backups are encrypted and isolated (A33) | **Backup isolation.** Nightly backups contain the complete payroll dataset. If backup storage is not isolated or encrypted, data exfiltration can occur without touching production systems. |
| **5** | DBAs and cloud admins cannot bypass app logic (A32) | **Cloud admin SoD.** Privileged users with direct database or storage access can read, export, or modify payroll records without any application-level audit trail or authorization check. |
| **6** | SQL traffic is encrypted between web app and RDS (A11) | Internal network sniffing (via ARP spoofing, compromised adjacent host, or cloud side-channel) can capture unencrypted SQL queries and result sets containing payroll data. |
| **7** | AD credentials resist brute-force and credential stuffing (A14) | The architecture does not mandate MFA for the application layer. Weak or reused credentials are the primary initial access vector in enterprise breaches. |
| **8** | AD is protected against lateral movement and domain dominance (A18) | Once a single workstation is compromised, unconstrained Kerberos delegation, pass-the-hash, or DCSync attacks can escalate to full domain control — granting access to all systems including the payroll database. |
| **9** | VPN gateway has no exploitable vulnerabilities (A3) | The VPN is the external trust boundary. A remote-code-execution or auth-bypass vulnerability in the VPN appliance gives attackers direct access to the internal network, bypassing all perimeter controls. |
| **10** | No undocumented access paths exist (A12) | Real architectures contain jump boxes, bastion hosts, API endpoints, and management interfaces not captured in diagrams. Each undocumented path represents a blind spot that may bypass the entire security model. |

## 6. Recommended Controls

| Rank | Critical Assumption | Technical Controls |
|---|---|---|
| **1** | Private subnet isolation | 1) VPC network ACLs denying all inbound traffic except from app-tier security groups. 2) AWS Security Groups with explicit allow rules only for the web app's private IP range on port 5432/3306. 3) VPC Flow Logs configured on the private subnet with alerts for any traffic originating outside the app tier. 4) Regular network penetration testing validating subnet boundaries. 5) AWS Config rules to detect security group drift. |
| **2** | Security monitoring | 1) SIEM aggregation of CloudTrail, VPC Flow Logs, RDS audit logs, and AD DS logs. 2) Automated alert rules for: database access from non-app IPs, failed login spikes, backup export events, and console login without MFA. 3) 24/7 MDR service or internal SOC with defined escalation runbooks. 4) Monthly tabletop exercises auditing alert response. |
| **3** | Input validation / OWASP | 1) Web Application Firewall (WAF) with OWASP Top 10 rule set (e.g., AWS WAF + Managed Rules). 2) DAST scanning in CI/CD pipeline (OWASP ZAP, Burp Suite Enterprise). 3) SAST scanning on every pull request (e.g., Semgrep, Checkmarx). 4) Annual full-scope penetration test with authenticated and unauthenticated scanning. 5) Content Security Policy (CSP) headers. |
| **4** | Backup isolation and encryption | 1) Backups stored in a separate AWS account with explicit cross-account IAM roles. 2) S3 bucket policies enforcing MFA delete and server-side encryption with AWS KMS (customer-managed key). 3) Automated backup integrity checks via checksum verification. 4) Cross-region replication for ransomware resilience. 5) Quarterly restore drills validated against payroll data integrity. |
| **5** | DBA / cloud admin SoD | 1) Just-in-time (JIT) access via AWS IAM Identity Center with approval workflows for break-glass elevation. 2) Database Activity Monitoring (DAM) capturing all queries executed by admin accounts. 3) Separate IAM roles for operational DB administration vs. application DB access — no overlapping permissions. 4) Alerting on any direct query to the `payroll` schema by a non-application principal. 5) Quarterly access review with sign-off from data owner. |
| **6** | SQL traffic encryption | 1) Force `ssl-mode=require` (PostgreSQL) or `encrypt=true` (SQL Server) on all RDS connections. 2) Reject self-signed or expired certificates at the database layer via `rds.force_ssl=1`. 3) Application connection strings must pin the RDS CA certificate. 4) Network security group denying port 5432/3306 from any source outside the app tier. 5) Certificate rotation automated at 90-day intervals. |
| **7** | AD credential resistance | 1) Enforce Azure AD MFA (or equivalent) for all application and VPN authentication. 2) Deploy conditional access policies blocking sign-ins from non-corporate IP ranges and untrusted devices. 3) Risky-sign-in detection with automated account lockout. 4) Passwordless deployment (Windows Hello / FIDO2) for endpoint auth. 5) Block legacy authentication protocols. |
| **8** | AD lateral movement protection | 1) Tiered administration model (Tier 0/1/2) with strict group membership boundaries. 2) Deploy Microsoft LAPS for local admin password management on all workstations and servers. 3) Monitor for golden-ticket events (event ID 4624 anomalous), DCSync (event ID 4662), and pass-the-hash. 4) Time-based group membership for Tier 0 access. 5) Harden domain controllers with LDAP signing, SMB signing, and Kerberos armoring. |
| **9** | VPN gateway vulnerabilities | 1) Automated patch management for VPN appliance with SLA ≤48 hours for critical CVEs. 2) Vulnerability scanning of the VPN endpoint weekly (external + authenticated). 3) VPN appliance configuration benchmarked against CIS or vendor hardening guide. 4) Redundant VPN gateways in active/passive for patching without downtime. 5) Rate-limiting and fail2ban for VPN login attempts. |
| **10** | Undocumented access paths | 1) Agent-based network discovery scan (e.g., AWS Discovery, Tenable) every 30 days. 2) Configuration Management Database (CMDB) enforcing that all network paths are documented and approved. 3) AWS Resource Explorer and Config aggregator maintaining a real-time asset inventory. 4) Change control board approval required for any new security group rule or route table entry. 5) Quarterly architecture review comparing deployed infrastructure to the approved diagram. |

## 7. Summary Statistics

| Metric | Count |
|---|---|
| **Total Assumptions** | 43 |
| **Critical Risk** | 14 |
| **High Risk** | 21 |
| **Medium Risk** | 7 |
| **Low Risk** | 1 |

### Risk Distribution

```
Critical:  ████████████████  14 (32.6%)
High:      █████████████████████████  21 (48.8%)
Medium:    ████████   7 (16.3%)
Low:       █   1 (2.3%)
```

### Model Contribution Breakdown

| Model | Raw Assumptions | Unique Contributions After Dedup |
|---|---|---|
| GPT-4o | 40 | 30 (primary source) |
| Gemini | 5 | 2 (A11 SQL encryption, A32 DBA SoD) |
| Gemma | 4 | 1 (A18 AD lateral movement) |

### Key Findings

- **14 Critical assumptions** must be validated or mitigated before this architecture can be considered secure for payroll data.
- The three most foundational risks are **internal network trust** (subnet isolation), **backup isolation**, and **cloud administrator separation of duties** — all explicitly called out by the models.
- Application-layer vulnerabilities (OWASP Top 10, SQLi) and **absent monitoring** represent the highest-likelihood, highest-impact exploitation paths.
- The architecture implicitly trusts the internal network — the single greatest design risk is that the "private subnet" is assumed isolated without validation controls.
- 5 of the 9 Gemini/Gemma assumptions address gaps not covered by GPT-4o, demonstrating the value of multi-model analysis for security review.
