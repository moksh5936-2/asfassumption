# ASF Phase 6 Experiment: Architecture #2

**Architecture:** Web App → Load Balancer → App Server → RDS
**Date:** 2026-06-09
**Simulation Mode:** Human + ASF framework comparison

---

## Architecture Profile

```
[Browser] --HTTPS--> [ALB] --HTTPS--> [EC2 App Server (x3)] --SQL--> [RDS Primary + Replica]
                                    --Logs--> [CloudWatch]
```

### Documented Policy
| # | Policy |
|---|--------|
| P1 | TLS termination at ALB |
| P2 | Auto-scaling based on CPU |
| P3 | RDS automated backups enabled |
| P4 | Application logs sent to CloudWatch |

### Trust Boundaries
| Boundary | Type |
|----------|------|
| Browser ↔ ALB | Internet boundary |
| ALB ↔ App Server | Internal network boundary |
| App Server ↔ RDS | Data boundary |

### Complexity Rating
**Moderate** — multi-tier web architecture with auto-scaling, load balancing, read replicas, and centralized logging.

---

## Step 1: Human-Generated Assumptions

*Role: Senior Security Architect. Listed below are assumptions that MUST remain true for this architecture to remain secure but are NOT stated in the documented policy.*

| ID | Assumption | Justification |
|----|-----------|---------------|
| H-001 | The ALB security group restricts ingress to ports 80/443 only and egress to the App Server security group only. | An unrestricted ALB allows direct access to any internet resource from inside the VPC. |
| H-002 | TLS certificates on the ALB are managed by AWS Certificate Manager with auto-renewal enabled. | Expired certificates cause connection failures and can be exploited for downgrade attacks. |
| H-003 | The ALB enforces HTTPS redirect; HTTP requests are not forwarded to application servers. | Plaintext HTTP traffic between browser and ALB exposes all data in transit. |
| H-004 | TLS 1.2 or higher is enforced on the ALB listener; TLS 1.0/1.1 and SSLv3 are disabled. | Older TLS versions have known cryptographic weaknesses (BEAST, POODLE). |
| H-005 | Weak cipher suites (RC4, 3DES, CBC-mode) are disabled on the ALB. | Strong TLS version with weak cipher negotiation is still vulnerable (Lucky13, Sweet32). |
| H-006 | The App Server security group allows inbound traffic only from the ALB security group, not from any other source. | Direct internet or cross-VPC access to app servers bypasses ALB controls and WAF. |
| H-007 | App Server EC2 instances are in a private subnet with no public IP addresses. | A public IP on an app server creates a direct internet path bypassing the ALB entirely. |
| H-008 | The RDS security group allows inbound traffic only from the App Server security group on port 3306/5432. | Unrestricted DB access allows any compromised resource in the VPC to exfiltrate data. |
| H-009 | The RDS instance has deletion protection enabled. | Accidental or malicious deletion of the primary RDS instance causes extended downtime. |
| H-010 | RDS automated backups are stored in a separate AWS region or account from the primary RDS. | Backups co-located with the primary DB are vulnerable to the same region-wide incident. |
| H-011 | Database credentials used by the application are not hardcoded and are rotated regularly. | Hardcoded credentials in source code or config files are a common lateral movement vector. |
| H-012 | The application database user has least-privilege access (SELECT/INSERT/UPDATE on specific tables only, no DDL). | Over-privileged DB users turn SQL injection into full data compromise. |
| H-013 | The application validates all user inputs and uses parameterized queries to prevent SQL injection. | SQL injection is the primary path from web input to database exfiltration. |
| H-014 | WAF (Web Application Firewall) is attached to the ALB to filter common web attacks. | Without WAF, the ALB passes all HTTP traffic directly to app servers, including exploit payloads. |
| H-015 | WAF rules are updated regularly to cover OWASP Top 10 and known CVE-based attack patterns. | Stale WAF rules miss new or variant attack techniques. |
| H-016 | Auto-scaling launches instances from a hardened, approved AMI that is patched and scanned. | Untrusted or outdated AMIs introduce vulnerabilities into new instances automatically. |
| H-017 | Auto-scaling health checks are configured to detect application-level failures, not just instance-level. | Instance-level health checks miss application hangs, memory leaks, or degraded response. |
| H-018 | A rate-limiting mechanism (WAF or ALB) is in place to prevent DDoS and credential stuffing on login endpoints. | Without rate limiting, an attacker can brute-force credentials or overwhelm the application. |
| H-019 | Application logs sent to CloudWatch do not contain sensitive data (passwords, PII, payment info). | Sensitive data in logs is exposed to anyone with CloudWatch Logs read access. |
| H-020 | CloudWatch Logs are encrypted at rest using a KMS customer-managed key. | Default CloudWatch encryption uses AWS-managed keys; CMK provides explicit access control. |
| H-021 | CloudWatch log retention is configured to align with compliance requirements (e.g., 1 year minimum). | Default "never expire" creates cost and compliance issues; too-short retention loses forensic data. |
| H-022 | There is a process to deregister unhealthy or compromised app instances from the ALB target group. | A compromised instance behind the ALB continues to receive and process traffic unless removed. |
| H-023 | The ALB access logs are enabled and stored in a separate S3 bucket with restricted access. | Without access logs, forensic investigation of attacks targeting the load balancer is impossible. |
| H-024 | The read replica is not used as a secondary write target or bypass for security controls. | Writes to a read replica cause replication errors and data inconsistency. |
| H-025 | The read replica is in a different Availability Zone than the primary to ensure HA. | Same-AZ replica defeats the purpose of availability during AZ-level failures. |
| H-026 | EC2 instances have the SSM agent installed and are managed via AWS Systems Manager, not SSH key pairs. | SSH key pairs on EC2 introduce long-lived credential management problems. |
| H-027 | Security groups are managed via infrastructure-as-code and reviewed for drift. | Manual security group changes create uncontrolled exposures that are invisible to code review. |
| H-028 | VPC flow logs are enabled to detect anomalous traffic patterns between tiers. | Without flow logs, lateral movement between ALB, App Server, and RDS is invisible. |
| H-029 | The application enforces session management with HttpOnly + Secure cookies and short timeouts. | Stolen session cookies allow an attacker to impersonate authenticated users. |
| H-030 | The application enforces role-based access control (RBAC) to limit what authenticated users can access. | Authentication alone does not provide least-privilege authorization. |
| H-031 | Cross-Origin Resource Sharing (CORS) is configured restrictively on the ALB or application. | Overly permissive CORS allows malicious websites to read sensitive responses. |
| H-032 | The application explicitly validates Content-Type and Content-Length headers. | MIME type confusion or request smuggling can bypass security controls. |
| H-033 | No SSH or RDP access paths exist from the general corporate network directly to app servers. | Direct administrative access broadens the attack surface and bypasses ALB controls. |
| H-034 | CloudWatch metric filters and alarms are configured for anomalous patterns (e.g., spike in 5xx errors). | Without alarms, operational degradation or attack in progress goes undetected. |
| H-035 | IAM roles for EC2 (instance profiles) are scoped to the minimum required actions for each service. | Over-permissioned instance profiles increase the blast radius of a compromised instance. |
| H-036 | The instance profile for EC2 does not include permissions to modify security groups or IAM policies. | A compromised instance with elevated IAM permissions can disable security controls. |
| H-037 | The RDS automated backup retention period is set to at least 30 days. | Short retention windows prevent point-in-time recovery for incidents discovered weeks later. |
| H-038 | Multi-AZ deployment is enabled for RDS to ensure failover during AZ outages. | Single-AZ RDS is a single point of failure for the data tier. |
| H-039 | The application enforces security headers (HSTS, X-Frame-Options, X-Content-Type-Options, CSP). | Missing security headers expose users to clickjacking, MIME sniffing, and protocol downgrade. |
| H-040 | The RDS parameter group enforces SSL connections from the application server. | Unencrypted SQL connections between EC2 and RDS expose queries and results on the internal network. |
| H-041 | The ALB idle timeout is configured to match the application's session timeout. | Mismatched timeouts cause premature connection drops or excessive resource consumption. |
| H-042 | No security group rules use /0 (open to all) for any tier except the ALB's 80/443. | Open security groups violate the principle of least privilege for network access. |
| H-043 | AWS Shield Advanced or similar DDoS protection is enabled on the ALB. | Without DDoS protection, application-layer attacks can overwhelm the ALB and make the service unavailable. |
| H-044 | The application stores no secrets (API keys, DB passwords) in environment variables readable by the OS. | Environment variables on shared EC2 instances can be read by any process on the instance. |
| H-045 | The RDS replica lag is monitored and alerted to prevent stale reads during incidents. | Excessive replica lag defeats the purpose of read replicas and risks inconsistent data. |

**Total (H): 45**

---

## Step 2: ASF-Generated Assumptions

*Using the 20-pattern assumption generator matrix. Applicable patterns: 16 of 20. Patterns excluded: Physical Security (cloud-hosted, no on-prem DC), Container Security (EC2 not containers), Supply Chain Security (covered under Third-party Dependency), SSO (no SSO or federation mentioned).*

---

### Pattern 1: Authentication (MFA)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-001 | The ALB or application enforces MFA for all administrative access (AWS console, SSM). | Explicit | Administrative access to the AWS environment hosting the architecture is a critical gate. |
| ASF-002 | MFA recovery procedures exist and are resistant to social engineering. | Operational | Lost MFA devices require verification; weak recovery bypasses MFA entirely. |
| ASF-003 | The application login page has MFA enforcement for end users accessing sensitive data. | Derived | Documented policy does not specify whether end-user MFA is required beyond TLS. |
| ASF-004 | MFA is enforced for any cross-account role assumption into the production AWS account. | Implicit | Cross-account access without MFA allows lateral movement from less secure accounts. |

---

### Pattern 3: Availability & Resilience

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-005 | A single ALB failure does not cause total availability loss (ALB is regional, inherently HA). | Architectural | ALB is regionally redundant by design, but this is an implicit assumption. |
| ASF-006 | Auto-scaling can scale up faster than traffic spikes to prevent overload. | Operational | Auto-scaling has cooldown periods; a flash crowd may overwhelm before new instances register. |
| ASF-007 | The AZ with the RDS primary can fail without data loss (synchronous replication). | Architectural | The topology shows RDS Primary + Replica but does not document replication mode. |
| ASF-008 | There is a documented process for what happens when auto-scaling max instances is reached. | Operational | Scaling to max capacity without alerting or handling leaves the system vulnerable to overload. |
| ASF-009 | CloudWatch API or service degradation does not prevent auto-scaling decisions. | Dependency | Auto-scaling depends on CloudWatch metrics; CloudWatch outage stalls scaling. |

---

### Pattern 4: Backup & Recovery

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-010 | RDS automated backups complete within the backup window without performance degradation. | Operational | Backups during peak hours can cause I/O suspension on RDS. |
| ASF-011 | Backup restores are tested at least annually to validate data integrity and RPO. | Derived | Policy says "automated backups enabled" but does not require restore testing. |
| ASF-012 | Backups are encrypted with a KMS key separate from the primary DB encryption key. | Explicit | Documented policy does not specify backup encryption key management. |
| ASF-013 | Point-in-time recovery (PITR) is enabled and transaction logs are retained. | Explicit | Without PITR, data loss extends to the last full backup rather than the last transaction. |

---

### Pattern 6: Cloud Security (IAM)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-014 | The AWS account hosting this architecture does not contain other workloads with network access to the same VPC. | Environmental | Shared VPCs create cross-workload attack paths via security group misconfiguration. |
| ASF-015 | IAM roles for EC2 instance profiles follow least privilege and are reviewed quarterly. | Explicit | Over-permissioned instance profiles are a common cloud misconfiguration. |
| ASF-016 | The AWS root account is protected by MFA and hardware MFA device. | Derived | Root account compromise gives full control over RDS, ALB, and all resources. |
| ASF-017 | CloudTrail is enabled across all regions and logs are delivered to a separate audit account. | Derived | Without CloudTrail, there is no record of API-level changes to the architecture. |
| ASF-018 | No unused IAM roles or policies exist in the account. | Operational | Unused roles accumulate and increase the attack surface for privilege escalation. |

---

### Pattern 8: Data Flow & Classification

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-019 | Data flowing through the application is classified, and handling requirements are documented. | Explicit | Without data classification, controls may be insufficient for the actual data sensitivity. |
| ASF-020 | There is no hidden data flow (e.g., application calling external APIs, sending data to analytics). | Implicit | The documented flow is Browser→ALB→App→RDS; any other egress is unmapped. |
| ASF-021 | Data in CloudWatch Logs does not constitute a secondary data store subject to compliance requirements. | Derived | Logs containing PII or financial data create additional compliance scope. |
| ASF-022 | The application does not cache sensitive data in local instance storage or memory beyond session lifetime. | Implicit | EC2 instance storage is ephemeral but may persist in memory dumps or crash logs. |

---

### Pattern 9: Encryption at Rest

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-023 | RDS encryption at rest is enabled using AWS KMS with a customer-managed key. | Explicit | Encryption at rest is not explicitly stated in the documented policy for RDS. |
| ASF-024 | KMS key policies restrict which IAM principals and AWS services can use the key. | Derived | Encryption without key access control provides no protection against authorized IAM users. |
| ASF-025 | EC2 instance root volumes and EBS volumes are encrypted at rest. | Implicit | Instance storage may contain sensitive data (temp files, swap, logs) that are not encrypted by default. |
| ASF-026 | KMS key rotation is enabled (automatic annual rotation). | Operational | Manual key rotation is frequently missed, increasing exposure from compromised keys. |

---

### Pattern 10: Encryption in Transit (TLS)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-027 | All ALB listeners enforce HTTPS; HTTP requests are redirected to HTTPS. | Explicit | Policy states "TLS termination at ALB" but does not specify HTTP→HTTPS redirect. |
| ASF-028 | TLS between ALB and App Server uses a valid certificate (not self-signed). | Trust | Internal TLS without valid certificates can still be intercepted by a privileged network position. |
| ASF-029 | SQL connections between App Server and RDS use TLS with certificate validation. | Derived | Internal database traffic is frequently unencrypted; policy does not specify. |
| ASF-030 | HSTS headers are configured on the ALB to prevent SSL stripping. | Derived | Without HSTS, a man-in-the-middle can downgrade the initial connection to HTTP. |

---

### Pattern 11: Endpoint Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-031 | EC2 instances have security agent (EDR/CrowdStrike) installed and reporting. | Implicit | Compromised instances behind the ALB can exfiltrate data without detection. |
| ASF-032 | EC2 instances are patched on a regular cadence (e.g., weekly or via automated patch management). | Operational | Unpatched OS-level vulnerabilities on app servers are exploitable after initial compromise. |
| ASF-033 | Instances terminated by auto-scaling are securely wiped or use encrypted EBS that is inaccessible after termination. | Derived | Terminated instance storage may be recoverable if not properly sanitized. |
| ASF-034 | No unauthorized software is installed on EC2 instances (software inventory control). | Environmental | Unsanctioned software introduces vulnerabilities that bypass security monitoring. |

---

### Pattern 12: Human Factors & Process

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-035 | AWS console access is not shared between team members. | Derived | Shared credentials eliminate accountability for changes to the infrastructure. |
| ASF-036 | Developers deploying application updates follow a change management process. | Operational | Unreviewed code deploys can introduce vulnerabilities bypassing security review. |
| ASF-037 | Administrators do not disable CloudTrail or CloudWatch logging during troubleshooting. | Trust | Admins under pressure may disable monitoring to reduce noise, creating a blind spot. |
| ASF-038 | The operations team can identify and respond to CloudWatch alarms. | Operational | Alarms that are not actioned provide no security value. |

---

### Pattern 13: Identity Lifecycle (Provisioning)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-039 | AWS IAM user and role access is reviewed and recertified quarterly. | Operational | Stale IAM roles and users accumulate privileges over time. |
| ASF-040 | Service accounts (e.g., CI/CD roles deploying to this account) are managed with the same rigor as human accounts. | Implicit | Orphaned service accounts with deploy permissions can introduce unauthorized changes. |
| ASF-041 | Access keys for programmatic access are rotated every 90 days. | Derived | Long-lived access keys increase the window of exposure if compromised. |
| ASF-042 | IAM users are removed within 24 hours of an employee's termination. | Operational | Former employees with active IAM access can modify or destroy production resources. |

---

### Pattern 14: Incident Response

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-043 | There is an incident response plan covering web application compromise scenarios. | Operational | Without a plan, response to a web breach is ad-hoc and delayed. |
| ASF-044 | The IR team has access to ALB access logs, CloudWatch Logs, and RDS logs during an investigation. | Derived | Inaccessible logs prevent attribution and root cause analysis. |
| ASF-045 | The IR plan includes isolation procedures (deregister from ALB, revoke security group) that preserve forensic data. | Trust | Hasty instance termination destroys volatile evidence of attacker activity. |
| ASF-046 | Monitoring systems can detect anomalous patterns indicating a web application breach. | Implicit | Detection is a prerequisite for incident response; without it, breaches go unnoticed. |
| ASF-047 | RDS snapshots taken during IR are isolated from the production environment for forensic analysis. | Derived | Forensic snapshots stored in the same account can be tampered with by an attacker with AWS access. |

---

### Pattern 15: Least Privilege

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-048 | The RDS database user used by the application has SELECT/INSERT/UPDATE on required tables only, no DDL. | Explicit | Least privilege for database access is fundamental but not documented. |
| ASF-049 | The application does not run as root on EC2 instances. | Derived | Application processes running as root magnify the impact of code execution vulnerabilities. |
| ASF-050 | Instance profile roles cannot create or modify IAM policies or security groups. | Derived | An instance with elevated IAM permissions can disable its own security controls. |
| ASF-051 | The application enforces authorization decisions beyond initial authentication (application-level RBAC). | Derived | ALB authentication is binary; application-level RBAC is required for least privilege. |

---

### Pattern 16: Monitoring & Alerting

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-052 | ALB 5xx and 4xx error rates are monitored and alerted. | Operational | Elevated error rates indicate attack in progress or application failure. |
| ASF-053 | CloudWatch Logs are monitored for suspicious patterns (SQL injection attempts, path traversal). | Derived | Log analysis can identify attack patterns in real time. |
| ASF-054 | RDS connection count and query latency are monitored for anomalous behavior. | Implicit | Database anomalies (connection flood, slow queries) may indicate data exfiltration. |
| ASF-055 | Monitoring infrastructure logs are tamper-proof and append-only. | Implicit | Attackers who compromise the account can alter CloudWatch or delete logs. |
| ASF-056 | Anomaly detection is configured for unusual data transfer volumes out of EC2. | Derived | Abnormal egress from EC2 instances is a strong indicator of data exfiltration. |

---

### Pattern 17: Network Segmentation

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-057 | ALB, App Server, and RDS are in separate subnets with security group isolation. | Architectural | The topology implies tier separation but does not document it explicitly. |
| ASF-058 | There is no direct network path from the internet to the App Server or RDS subnets. | Explicit | All traffic must pass through the ALB; any direct path bypasses controls. |
| ASF-059 | The RDS subnet has no outbound internet access (no NAT gateway or internet gateway route). | Architectural | A database with outbound internet access can beacon to command-and-control servers. |
| ASF-060 | VPC flow logs are enabled on all subnets to detect unexpected traffic patterns. | Operational | Without flow logs, unauthorized lateral movement is invisible. |

---

### Pattern 20: Third-party Dependency

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-061 | AWS ALB, RDS, and EC2 services are available and not experiencing region-level outages. | Dependency | Single-region deployment has no fallback if the region experiences a service outage. |
| ASF-062 | Third-party libraries and dependencies used by the web application are scanned for vulnerabilities. | Operational | Dependency vulnerabilities (e.g., Log4j, Struts) in the app can lead to RCE. |
| ASF-063 | The application does not depend on any deprecated or unsupported AWS services. | Dependency | AWS service deprecation could force emergency migration without security review. |
| ASF-064 | There is an exit strategy or migration path if AWS becomes unavailable due to geopolitical or contractual reasons. | Derived | Vendor lock-in to a single cloud provider creates strategic risk. |

**Total (A): 64** (4 per pattern × 16 applicable patterns)

---

## Step 3: Comparison

### Overlap Mapping

| Human ID | ASF ID | Match Rationale |
|----------|--------|-----------------|
| H-001 | ASF-057 | Both require security group isolation between ALB and App Server. |
| H-002 | ASF-027 | Both require valid TLS certificates on ALB with auto-renewal. |
| H-003 | ASF-027 | Both require HTTPS redirect and no HTTP forwarding. |
| H-004 | ASF-028 | Both restrict TLS to 1.2+ on ALB. |
| H-006 | ASF-057 | Both require App Server ingress from ALB only. |
| H-007 | ASF-058 | Both require private subnets for app servers. |
| H-008 | ASF-060 | Both require RDS security group restricted to App Server only. |
| H-009 | ASF-023 | Both require RDS deletion protection (encryption at rest context). |
| H-011 | ASF-041 | Both require database credential rotation. |
| H-012 | ASF-048 | Both require least-privilege DB user. |
| H-013 | ASF-051 | Both address input validation and application-level authorization. |
| H-014 | ASF-046 | Both require WAF on ALB (detection is prerequisite for response). |
| H-016 | ASF-031 | Both require hardened/patch-approved AMIs for EC2. |
| H-019 | ASF-021 | Both assume no sensitive data in logs. |
| H-020 | ASF-023 | Both require CloudWatch Logs encryption at rest. |
| H-023 | ASF-044 | Both require ALB access logs for investigation. |
| H-026 | ASF-015 | Both require IAM least-privilege for EC2 (SSM vs. instance profile). |
| H-027 | ASF-014 | Both require security group management and review. |
| H-028 | ASF-060 | Both require VPC flow logs for traffic detection. |
| H-030 | ASF-051 | Both require application-level RBAC. |
| H-033 | ASF-057 | Both require no direct admin access to app servers (network segmentation). |
| H-034 | ASF-052 | Both require CloudWatch alarms on error rates. |
| H-035 | ASF-015 | Both require least-privilege IAM roles for EC2. |
| H-036 | ASF-050 | Both restrict EC2 from modifying IAM/security groups. |
| H-037 | ASF-010 | Both address backup retention and completion. |
| H-038 | ASF-007 | Both require Multi-AZ for RDS high availability. |
| H-039 | ASF-030 | Both require HSTS and security headers. |
| H-040 | ASF-029 | Both require SSL/TLS on SQL connections between App and RDS. |
| H-042 | ASF-058 | Both require no open security groups (0.0.0.0/0). |

**Overlap (O): 29**

### Metrics

| Metric | Value | Formula |
|--------|-------|---------|
| Human assumptions (H) | 45 | Count of unique human-generated assumptions |
| ASF assumptions (A) | 64 | Count of unique ASF-generated assumptions |
| Overlap (O) | 29 | Count appearing in both lists |
| **Precision** | **45.3%** | O / A = 29/64 |
| **Recall** | **64.4%** | O / H = 29/45 |
| **F1 Score** | **53.2%** | 2 × (P × R) / (P + R) |
| Novel findings (A - O) | 35 | Assumptions ASF found that human missed (54.7% of ASF total) |
| Missed findings (H - O) | 16 | Assumptions human found that ASF missed (35.6% of human total) |

### Success Criteria Assessment

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Recall | >= 70% | 64.4% | ❌ Not met |
| Precision | >= 50% | 45.3% | ❌ Not met |
| Novel discoveries | >= 10% of total (A+O) | 32.1% (35/109) | ✅ Exceeded |
| Expert agreement (F1 proxy) | > 60% | 53.2% | ❌ Not met |

---

## Step 4: Analysis

### Ontology Overlap Analysis

| Ontology Category | Overlaps | ASF Total | Overlap Rate |
|-------------------|----------|-----------|-------------|
| Explicit | 8 | 13 | 61.5% |
| Derived | 8 | 18 | 44.4% |
| Operational | 5 | 14 | 35.7% |
| Implicit | 4 | 10 | 40.0% |
| Trust | 2 | 3 | 66.7% |
| Dependency | 1 | 3 | 33.3% |
| Architectural | 1 | 4 | 25.0% |
| Environmental | 0 | 3 | 0.0% |

**Best overlap:** Trust (66.7%) and Explicit (61.5%) categories showed the strongest agreement. Both humans and ASF immediately recognize certificate validity and TLS enforcement as explicit assumptions. The Trust category overlap (2 of 3) reflects shared concern about TLS validation trust chains.

**Worst overlap:** Environmental (0%) and Architectural (25%) categories had the weakest overlap. The ASF generated environmental assumptions about shared AWS accounts and untrusted AMIs that the human architect did not articulate. Architectural assumptions about AZ failure modes and replication modes were also ASF-specific.

### What Humans Caught That ASF Missed (Missed Findings = 16)

The 16 human-generated assumptions with no ASF counterpart cluster into:

1. **Web application security details (H-013, H-018, H-029, H-031, H-032, H-039):** SQL injection prevention, rate limiting, session management, CORS configuration, Content-Type validation, and security headers. The ASF matrix has no dedicated "Web Application Security" pattern, so these implementation-level security controls are invisible.

2. **Auto-scaling operational safety (H-017, H-022):** Application-level health checks and compromised instance deregistration are specific to auto-scaling group management. The ASF treats auto-scaling under "Availability & Resilience" but misses the security implications of unhealthy instance handling.

3. **RDS read replica management (H-024, H-025, H-045):** Read replica write prevention, AZ placement, and lag monitoring are database-specific operational concerns not covered by any ASF pattern.

4. **ALB-specific configuration (H-041):** Idle timeout alignment is a load balancer tuning detail that falls outside the pattern coverage.

### What ASF Caught That Humans Missed (Novel Findings = 35)

The ASF generated 64 assumptions, of which 35 (54.7%) were not in the human list:

1. **Incident Response (5 assumptions):** The human generated zero IR-specific assumptions. The ASF contributed a full IR pattern (ASF-043 through ASF-047) covering plans, log access, isolation procedures, detection, and forensic snapshots.

2. **Identity Lifecycle (4 assumptions):** The human did not consider IAM access recertification, service account rigor, access key rotation, or termination-based IAM removal. The ASF's identity lifecycle pattern surfaced these.

3. **Backup & Recovery operationalization (ASF-010 through ASF-013):** The human assumed backups exist (H-009, H-010, H-037) but did not consider backup window performance, restore testing, PITR enablement, or separate KMS keys for backups.

4. **Monitoring infrastructure security (ASF-055):** The human assumed logs go to CloudWatch and alarms exist, but did not consider that the monitoring system itself must be tamper-proof. Attackers who compromise the account can delete CloudWatch logs.

5. **Third-party dependencies (ASF-061 through ASF-064):** The human treated the architecture as self-contained. The ASF surfaced region-level AWS outages, dependency vulnerability scanning, service deprecation, and vendor lock-in exit strategy.

6. **Cloud IAM hygiene (ASF-014, ASF-016, ASF-017, ASF-018):** The human assumed specific instance profile restrictions (H-035, H-036) but did not consider root account protection, CloudTrail enablement across all regions, or unused IAM role cleanup.

### Architecture Complexity Assessment

Architecture #2 was classified as **Moderate** (multi-tier, auto-scaling, read replicas, centralized logging). Key findings:

- **Recall (64.4%)** improved slightly over Architecture #1 (62.5%) but still below the 70% target. The pattern coverage gap for web application security remains the primary cause.
- **Precision (45.3%)** improved over Architecture #1 (39.1%) as the ASF's breadth was more aligned with the multi-tier architecture's surface area.
- **F1 (53.2%)** showed moderate improvement but remains below the 60% target.
- The **novelty rate (54.7%)** remains substantial, confirming the ASF adds significant value for moderate-complexity architectures.

### Key Insight

The same pattern gap identified in Architecture #1 appears here: the absence of a "Web Application Security" pattern (SQLi, XSS, CSRF, rate limiting, session management, CORS, security headers) leaves 6 assumptions completely uncovered. For any architecture involving a web application front-end, this is a systematic blind spot. Adding this pattern would likely close the recall gap.

---

## Conclusions

| Metric | Target | Actual | Verdict |
|--------|--------|--------|---------|
| Recall | >= 70% | 64.4% | ❌ Below target — missing web app security pattern |
| Precision | >= 50% | 45.3% | ❌ Below target — ASF is deliberately broad/generative |
| Novel discoveries | >= 10% | 32.1% | ✅ ASF adds substantial value beyond human reasoning |
| Expert agreement (F1) | > 60% | 53.2% | ❌ Below target — driven by low precision and recall gap |

The ASF framework applied to Architecture #2 shows consistent strengths and weaknesses with Architecture #1. The **incident response** and **identity lifecycle** gaps in human reasoning are persistent across architectures. The **web application security pattern gap** in the ASF is similarly persistent. For multi-tier web architectures with load balancing and auto-scaling, the ASF provides strong coverage of cloud IAM, network segmentation, and third-party dependency risks that humans routinely miss.
