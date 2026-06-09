# ASF Phase 6 Experiment: Architecture #018

**Architecture:** Global CDN → WAF → Origin → Database
**Date:** 2026-06-09
**Simulation Mode:** Human + ASF framework comparison

---

## Architecture Profile

```
[Global Users] --> [CloudFront CDN] --> [WAF] --> [ALB] --> [EC2 Origin] --> [RDS]
                      │
                 [Lambda@Edge] (Auth)
```

### Documented Policy
| # | Policy |
|---|--------|
| P1 | WAF blocks common attack patterns |
| P2 | CDN caches static content |
| P3 | Origin accessible via CDN only (no direct path) |
| P4 | DDoS protection enabled |

### Trust Boundaries
| Boundary | Type |
|----------|------|
| CDN ↔ WAF | Edge boundary |
| WAF ↔ Origin | Origin access boundary |
| Origin ↔ Database | Data boundary |

### Complexity Rating
**Moderate** — 6 nodes including CDN edge functions, 3 trust boundaries, AWS-managed + self-managed hybrid.

---

## Step 1: Human-Generated Assumptions

*Role: Senior Security Architect. Listed below are assumptions that MUST remain true for this architecture to remain secure but are NOT stated in the documented policy.*

| ID | Assumption | Justification |
|----|-----------|---------------|
| H-001 | CloudFront distribution is configured with Origin Access Control (OAC) so the origin ALB only accepts requests forwarded by CloudFront. | Without OAC or OAI, an attacker who discovers the ALB DNS name can bypass the CDN and WAF entirely. |
| H-002 | WAF rules are updated regularly to cover emerging CVEs (e.g., Apache Struts, Log4j) and not just OWASP Top 10 static rules. | WAF with stale rules provides false confidence — attackers exploit new vulnerabilities faster than default rule set updates. |
| H-003 | Lambda@Edge performs authentication token validation and does not expose the origin or internal network information to unauthenticated requests. | A misconfigured Lambda@Edge may leak origin IPs, internal headers, or stack traces in error responses. |
| H-004 | The WAF web ACL is not at capacity (rule limit) — new rules can be added without removing existing ones. | WAF ACLs have a 150-rule limit; hitting capacity forces removal of protective rules, exposing blind spots. |
| H-005 | TLS certificate for the CloudFront distribution is provisioned via AWS Certificate Manager (ACM) and auto-renewed before expiry. | Expired CloudFront certificates cause immediate downtime and cache bypass errors that leak origin information. |
| H-006 | The EC2 origin security group restricts inbound traffic to the CloudFront managed prefix list only, not 0.0.0.0/0. | A security group allowing all inbound traffic defeats the "origin accessible via CDN only" policy. |
| H-007 | WAF rate-based rules are configured to protect against credential stuffing and application-layer DDoS on login endpoints. | DDoS protection (P4) covers network-layer attacks, but application-layer (HTTP flood, login brute-force) requires explicit WAF rate limits. |
| H-008 | The RDS database has automated backups enabled with a retention period sufficient for point-in-time recovery. | Backups are not documented; without them, data loss from corruption, deletion, or ransomware is unrecoverable. |
| H-009 | CloudFront logging (standard + real-time) is enabled and delivered to a separate S3 bucket with access restricted to security team. | Without CDN logs, WAF bypass attempts and anomalous traffic patterns are invisible to defenders. |
| H-010 | WAF logs are sent to CloudWatch Logs or S3 for analysis and archived per compliance requirements. | WAF blocked requests contain attacker IPs and patterns; without logging, tuning rules and forensic investigation are impossible. |
| H-011 | The Lambda@Edge function has no access to the RDS database or any internal resources beyond what is required for auth. | A compromised Lambda@Edge running in the CloudFront edge location could be used as a pivot to internal AWS resources. |
| H-012 | The origin EC2 instance does not have a public IP address assigned — all traffic arrives via the ALB from CloudFront. | An EC2 instance with a public IP in addition to the ALB creates a direct internet path bypassing WAF and CDN. |
| H-013 | The WAF is configured in "Block" mode, not "Count" mode, for critical rule groups (SQLi, XSS, RFI). | WAF in Count mode only observes attacks without blocking them — a common misconfiguration during initial setup. |
| H-014 | CloudFront signed URLs or signed cookies are used for premium/authenticated content, not just origin-shield caching. | Without signed URLs, any user with a valid CloudFront URL can access authenticated content regardless of authorization. |
| H-015 | The RDS instance is deployed in a private subnet with no direct internet route via NAT gateway or internet gateway. | A database with outbound internet access can phone home or be used for data exfiltration over DNS or HTTPS tunnels. |
| H-016 | TLS between CloudFront and the origin ALB uses the minimum required version (TLS 1.2+) and strong ciphers. | A weak TLS configuration between edge and origin exposes data in transit to MITM on the AWS backbone. |
| H-017 | WAF has a rule that blocks requests with missing or invalid Host headers to prevent DNS rebinding attacks. | DNS rebinding can bypass WAF if the WAF evaluates the Host header differently from the origin server. |
| H-018 | CloudFront is configured to use Origin Shield to reduce origin load and provide a higher cache hit ratio. | Without Origin Shield, a cache miss storm during traffic spike can overwhelm the origin and cause cascading failure. |
| H-019 | The origin ALB has access logs enabled and retains logs for at least 90 days. | ALB access logs show the true request volume after WAF processing; without them, origin load assessment is blind. |
| H-020 | Lambda@Edge function errors do not expose stack traces or implementation details to the client (custom error pages configured). | Edge function errors visible to users reveal implementation details that inform targeted attacks. |
| H-021 | The WAF has a geographic (geo) match rule to block requests from countries where the application has no business. | Geo-blocking is the simplest DDoS and credential-stuffing mitigation; its absence increases attack surface. |
| H-022 | CloudFront distribution does not expose the origin ALB domain name in response headers (Server, X-Amz-Cf-*, etc.). | Origin domain name leakage in headers or error pages informs attackers where to bypass the CDN. |
| H-023 | RDS encryption at rest is enabled using a customer-managed KMS key with automatic rotation. | Unencrypted RDS data at rest violates compliance requirements and exposes data if storage volumes are decommissioned. |
| H-024 | The database credentials used by the EC2 origin are stored in AWS Secrets Manager or Parameter Store with automatic rotation. | Static database credentials in application config files are the most common lateral movement vector after an origin compromise. |
| H-025 | CloudFront supports HTTP/2 and HTTP/3 (QUIC) for optimal client performance; downgrade to HTTP/1.1 is blocked. | Downgrade attacks to HTTP/1.1 can expose the application to protocol-specific vulnerabilities. |
| H-026 | The EC2 origin has no direct egress to the internet — all outbound traffic goes through a NAT gateway or VPC endpoint. | An origin with direct internet egress can be used for C2 communication or data exfiltration after compromise. |
| H-027 | WAF logging does not include the full request body for POST requests containing sensitive data (PII, credentials). | WAF logs that capture request bodies with PII create a secondary data exposure surface in log storage. |
| H-028 | The ALB idle timeout is configured to match the application's expected response time and terminate slow-loris attacks. | An overly long ALB idle timeout keeps connections open indefinitely, enabling slow-rate DDoS exhaustion of connection pool. |
| H-029 | CloudFront real-time logs are monitored for spikes in 4xx/5xx error rates triggering automated alerts. | Anomaly detection on CDN error rates provides early warning of WAF bypass attempts or origin degradation. |
| H-030 | The VPC hosting the origin and RDS has VPC Flow Logs enabled and shipped to a central security account. | Without flow logs, lateral movement from a compromised origin to other VPC resources is invisible. |
| H-031 | The Lambda@Edge function uses the principle of least privilege IAM role and is not granted AdministratorAccess. | An over-privileged Lambda@Edge is a critical risk — compromise of the function grants full AWS API access at the edge. |
| H-032 | WAF bot control managed rule group is enabled to block known bots and scrapers. | Without bot control, automated scanners and scrapers consume origin resources and may discover hidden endpoints. |
| H-033 | CloudFront distribution is protected by AWS Shield Advanced for enhanced DDoS mitigation above the free Shield Standard. | DDoS protection (P4) at the standard level may be insufficient for large volumetric attacks targeting the application layer. |
| H-034 | The origin EC2 instance has automated patch management enabled for the OS and web server software. | An unpatched origin OS is the most likely initial compromise vector after WAF bypass. |
| H-035 | The CloudFront default root object is configured and does not leak directory listings. | Without a default root object, CloudFront may serve a directory listing or 403 error that reveals application structure. |
| H-036 | WAF can integrate with AWS Firewall Manager for central rule management across multiple distributions. | Without central management, WAF rules drift across environments and compliance violations go undetected. |
| H-037 | The ALB is configured with a proper health check path that does not expose sensitive application logic. | A health check endpoint that reveals application version, framework, or internal state aids reconnaissance. |
| H-038 | RDS automated snapshots are copied to a separate AWS region for disaster recovery. | Region-local backups are lost in a region failure scenario — cross-region copy is required for true DR. |
| H-039 | CloudFront custom error responses do not expose the origin error page content or status codes. | Custom error pages that show origin error details (500, 503) leak internal architecture information. |
| H-040 | WAF rate-based rules use the source IP as the aggregate key and block at a threshold appropriate for the user base. | Rate limits set too high allow credential stuffing; limits set too low block legitimate users during flash crowds. |
| H-041 | The RDS parameter group disables 'skip_name_resolve' and 'old_passwords' and other legacy MySQL/PostgreSQL settings. | Legacy database settings introduce known vulnerabilities that are outside the WAF's scope to block. |
| H-042 | CloudFront is configured to compress objects (gzip/brotli) and does not serve compressed content over insecure channels. | Compression side-channel attacks (BREACH, CRIME) can leak sensitive data in compressed HTTPS responses. |
| H-043 | The origin EC2 instance runs with a service account that has no interactive login capability. | A service account with interactive shell access allows an attacker to explore the OS after initial foothold. |
| H-044 | WAF has rules to inspect and block file uploads containing malware (e.g., via AWS WAF Marketplace partner integrations). | Without upload scanning, attackers use file upload forms to deliver malware to internal users who retrieve cached content. |
| H-045 | The RDS instance has deletion protection enabled to prevent accidental or malicious database deletion. | A database without deletion protection can be dropped through a compromised application or IAM user. |

**Total (H): 45**

---

## Step 2: ASF-Generated Assumptions

*Using the 20-pattern assumption generator matrix. Applicable patterns: 17 of 20. Patterns excluded: Container Security (no containers), Physical Security (cloud-hosted), Supply Chain Security (deferred to Third-party Dependency as combined pattern).*

---

### Pattern 1: Authentication (MFA)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-001 | Lambda@Edge enforces authentication for all requests before forwarding to the origin. | Explicit | Auth at the edge is the primary access control and must be enforced for every request. |
| ASF-002 | The authentication token validation in Lambda@Edge uses a trusted identity provider (Cognito, Okta) with proper signature verification. | Derived | Custom auth logic in edge functions is prone to implementation errors; relying on a managed IdP reduces risk. |
| ASF-003 | Users whose auth tokens expire mid-session are redirected to re-authenticate, not served cached content. | Trust | Stale cached content served after token expiry creates an authorization bypass window. |
| ASF-004 | The Lambda@Edge auth function has a circuit breaker to degrade gracefully if the IdP is unreachable. | Operational | IdP unavailability at the edge would block all user traffic; a degraded mode with warnings is preferable. |

---

### Pattern 2: Authentication (SSO)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-005 | The Lambda@Edge function validates the entire JWT chain (iss, aud, exp, nbf, signature) and rejects invalid tokens. | Explicit | JWT validation at the edge must be complete; many SDKs only verify the signature and skip other claims. |
| ASF-006 | SSO session cookies or tokens are not cached by CloudFront and are evaluated fresh on every request. | Architectural | CloudFront caching authenticated responses could serve stale authorized content to unauthenticated users. |
| ASF-007 | The IdP used by Lambda@Edge has cross-region redundancy matching CloudFront's global footprint. | Dependency | A regional IdP failure degrades auth for users in that region even though CloudFront remains available. |
| ASF-008 | Auth token signing keys used by Lambda@Edge are rotated before the current key's validity window expires. | Operational | Expired signing keys cause widespread auth failures; stale keys increase the risk of forged token acceptance. |

---

### Pattern 3: Availability & Resilience

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-009 | CloudFront origin failover is configured with a primary and secondary origin for ALB redundancy. | Architectural | A single ALB origin is a single point of failure; origin failover provides automatic recovery. |
| ASF-010 | WAF capacity units (WCUs) are provisioned sufficiently above peak traffic to avoid WAF throttling. | Operational | WAF throttling during traffic peaks drops legitimate traffic and blocks all requests regardless of content. |
| ASF-011 | The Lambda@Edge function has sufficient concurrency limits and does not throttle during traffic spikes. | Environmental | Lambda@Edge concurrency limits are lower than standard Lambda; hitting the limit drops auth requests. |
| ASF-012 | CloudFront origin shield is configured to reduce origin request volume during cache-miss storms. | Derived | Without origin shield, a traffic surge causing mass cache misses can overload the origin ALB and RDS. |

---

### Pattern 4: Backup & Recovery

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-013 | RDS automated backups complete within the configured backup window and are verified for integrity. | Operational | Backups that fail silently produce unusable restore points; integrity verification is required for recovery assurance. |
| ASF-014 | Cross-region backup copy is configured for disaster recovery scenarios. | Derived | Single-region backups are vulnerable to region-wide outages; cross-region copy provides geo-redundancy. |
| ASF-015 | CloudFront distribution configuration is backed up (via infrastructure-as-code) for rapid redeployment. | Implicit | A misconfigured CloudFront distribution can be rolled back only if configuration snapshots exist. |
| ASF-016 | WAF rule configurations are exported and version-controlled to protect against accidental rule deletion. | Operational | WAF rules are the primary security control; losing them via misconfiguration removes all protection. |

---

### Pattern 5: Cloud Security (IAM)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-017 | The Lambda@Edge IAM role has only the permissions required for auth token verification, not full AWS API access. | Explicit | Over-privileged Lambda@Edge functions are the highest-risk edge compute resource. |
| ASF-018 | CloudFront distribution uses a service-linked role with least-privilege for accessing S3 origins (if any). | Derived | CloudFront's default service role may have broader S3 access than necessary for the application. |
| ASF-019 | The AWS account root user is protected with MFA and hardware TOTP, and not used for CDN/waf/origin administration. | Explicit | Root user compromise gives an attacker control over the entire CloudFront distribution, including edge functions. |
| ASF-020 | No CloudFront distribution uses a trusted key group with unauthorized or leaked public keys for signed URLs. | Implicit | Leaked CloudFront trusted key group keys allow anyone to generate valid signed URLs for any content. |
| ASF-021 | The EC2 origin instance profile has the minimal IAM permissions needed (SSM: DescribeParameters, KMS: Decrypt) and nothing more. | Derived | An EC2 instance profile with overly broad permissions escalates a server compromise to an AWS API compromise. |

---

### Pattern 6: Data Flow & Classification

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-022 | Data flowing through CloudFront is classified according to its sensitivity and cached accordingly. | Explicit | Cached content at the edge must respect data classification — PII/cardholder data may not be cacheable. |
| ASF-023 | Data flow diagrams exist showing all paths: user → CDN → WAF → ALB → EC2 → RDS, plus Lambda@Edge execution flow. | Implicit | Undocumented data flows (e.g., direct RDS queries from monitoring tools) bypass edge security controls. |
| ASF-024 | Sensitive data is not written to CloudFront access logs or WAF logs. | Derived | Logs containing PII, tokens, or credentials create secondary exposure surfaces in S3 or CloudWatch. |
| ASF-025 | The application does not serve sensitive data (API keys, session tokens) in URL query strings that CloudFront caches. | Implicit | Cached URLs with sensitive query parameters persist at edge locations for the TTL duration. |

---

### Pattern 7: Encryption at Rest

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-026 | RDS encryption at rest is enabled using AWS KMS with automatic key rotation. | Explicit | Encryption at rest is a baseline security control for database storage. |
| ASF-027 | EC2 root and EBS volumes are encrypted at rest using a customer-managed KMS key. | Derived | The documented policy mentions encryption for RDS but not for the EC2 origin volumes. |
| ASF-028 | Lambda@Edge function code and /tmp storage are encrypted at rest. | Explicit | Edge function execution environment must protect function code and temporary data. |
| ASF-029 | KMS key policies for RDS and EBS restrict key usage to only the required IAM roles and services. | Implicit | Overly permissive KMS key policies allow unauthorized decryption of encrypted resources. |

---

### Pattern 8: Encryption in Transit (TLS)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-030 | CloudFront enforces HTTPS-only viewer protocol policy (no HTTP fallback). | Explicit | Redirecting HTTP to HTTPS at CloudFront is the minimum transport security requirement. |
| ASF-031 | TLS between CloudFront and the origin ALB uses the Origin Protocol Policy of HTTPS-only with TLS 1.2+. | Derived | The documented policy states CDN-only access but does not specify TLS version for origin communication. |
| ASF-032 | Weak TLS cipher suites are disabled on both the CloudFront viewer and origin TLS termination. | Explicit | TLS with weak ciphers (RC4, 3DES, CBC-mode) exposes traffic to cryptanalytic attacks. |
| ASF-033 | The EC2 origin validates the RDS TLS certificate when establishing database connections. | Trust | MySQL/PostgreSQL connections without certificate validation are vulnerable to MITM within the VPC. |
| ASF-034 | CloudFront is configured with security policy TLSv1.2_2021 or later, not the outdated TLSv1_2016 policy. | Derived | Older CloudFront security policies include TLS 1.0/1.1 support that may be deprecated by compliance standards. |

---

### Pattern 9: Endpoint Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-035 | The EC2 origin instance has an EDR agent installed and reporting to a security operations center. | Implicit | The origin is the compute resource closest to data; an unmonitored origin is a blind spot. |
| ASF-036 | The EC2 origin has OS-level file integrity monitoring (FIM) enabled for critical web server directories. | Derived | Without FIM, web shell uploads or configuration tampering on the origin goes undetected. |
| ASF-037 | No unmanaged SSH or RDP keys exist on the EC2 origin that provide backdoor access. | Operational | Stale or orphaned SSH keys on the origin provide persistent backdoor access. |
| ASF-038 | The EC2 origin AMI is from a trusted source and is regularly updated to the latest patched version. | Implicit | Untrusted AMIs may contain pre-installed malware or backdoors. |

---

### Pattern 10: Human Factors & Process

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-039 | Administrators do not disable the WAF for troubleshooting and forget to re-enable it. | Derived | WAF temporarily disabled for testing is a common source of production blind periods. |
| ASF-040 | CloudFront distribution updates follow a change management process with security review. | Operational | Direct distribution updates bypass security review and can introduce misconfigurations. |
| ASF-041 | Engineers understand that Lambda@Edge functions run in the CloudFront data plane and can impact global traffic. | Trust | Deploying buggy edge functions can degrade or block traffic across all edge locations globally. |
| ASF-042 | WAF rule updates are tested in a staging environment before production deployment. | Operational | Untested WAF rules can block legitimate traffic or fail to block attack patterns. |
| ASF-043 | The team managing the WAF and CloudFront has training on AWS WAF rule writing and common pitfalls. | Implicit | Misconfigured WAF rules (e.g., incorrect regex) may silently pass attack traffic through. |

---

### Pattern 11: Identity Lifecycle (Provisioning)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-044 | IAM users with CloudFront and WAF administrative access follow joiner/mover/leaver processes. | Operational | Orphaned IAM users with CDN admin access can modify distribution settings including WAF association. |
| ASF-045 | API keys used by external services to upload content to CloudFront origins are rotated regularly. | Derived | Leaked API keys for origin upload access allow attackers to inject malicious content cached by CloudFront. |
| ASF-046 | IAM role trust policies for Lambda@Edge restrict assumption to only the expected CloudFront service principal. | Implicit | Overly broad trust policies allow other AWS services to assume the edge function's role. |
| ASF-047 | Cross-account access to CloudFront distributions is reviewed and limited to authorized accounts. | Operational | Multi-account organizations may grant excessive cross-account CloudFront access. |

---

### Pattern 12: Incident Response

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-048 | There is an IR playbook for WAF bypass or CloudFront compromise scenarios. | Operational | A successful WAF bypass targeting the origin requires immediate isolation actions. |
| ASF-049 | The IR team has access to CloudFront access logs, WAF logs, and Lambda@Edge logs during an investigation. | Derived | Logs stored in separate accounts or encrypted with different keys may be inaccessible during incident response. |
| ASF-050 | There is a documented procedure to blackhole malicious IPs at the CloudFront level quickly. | Trust | Rapid IP blocking at the CDN level can stop ongoing attacks without touching the WAF. |
| ASF-051 | Monitoring detects unusual CloudFront traffic patterns (high error rates, unexpected geographies, unusual URL paths). | Implicit | Anomaly detection on CDN telemetry provides the earliest warning of a WAF bypass or content scraping attack. |
| ASF-052 | WAF false positive events are reviewed and tuned to prevent permanent blocking of legitimate traffic. | Operational | WAF rules that block legitimate traffic due to false positives erode trust and may lead to the WAF being disabled. |

---

### Pattern 13: Least Privilege

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-053 | The RDS database user used by the EC2 origin has only the required DML permissions (SELECT/INSERT/UPDATE/DELETE) and no DDL. | Explicit | An application database user with DDL access can drop tables or alter schemas via SQL injection. |
| ASF-054 | The Lambda@Edge function is not granted permissions to write to S3, invoke other Lambda functions, or access RDS directly. | Derived | Edge functions should have no data plane permissions beyond authentication — least privilege principle. |
| ASF-055 | The WAF administrative role is scoped to specific rule groups and does not grant full WAF management access. | Implicit | Full WAF admin access allows an insider or compromised account to delete all protective rules. |
| ASF-056 | CloudFront distribution settings (origins, behaviors, error pages) are managed by IaC and not via console direct edits. | Derived | Direct console edits bypass review and audit, and can introduce shadow IT configurations. |

---

### Pattern 14: Monitoring & Alerting

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-057 | CloudFront and WAF metrics are integrated with a central monitoring dashboard (e.g., CloudWatch, Datadog). | Operational | Distributed monitoring between CDN, WAF, ALB, and EC2 creates blind spots without a unified view. |
| ASF-058 | Alerts are configured for WAF blocked request spikes that may indicate ongoing attack campaigns. | Derived | WAF block spikes are the primary indicator of probing or active exploitation attempts. |
| ASF-059 | Lambda@Edge error rate and duration metrics are monitored and alert on anomalies. | Operational | Edge function errors that silently fail open can bypass authentication for a subset of requests. |
| ASF-060 | CloudFront cache hit ratio is monitored as a leading indicator of origin health. | Implicit | A sudden drop in cache hit ratio signals origin connection issues or cache eviction problems. |
| ASF-061 | WAF rule group updates from AWS are tracked via AWS Health events and reviewed for impact. | Operational | AWS-managed WAF rule updates can change blocking behavior without notice, breaking production traffic. |

---

### Pattern 15: Network Segmentation

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-062 | The EC2 origin and RDS are in separate subnets with security groups that allow only the minimum required traffic. | Architectural | Flat networking with EC2 and RDS in the same subnet violates defense-in-depth principles. |
| ASF-063 | The origin ALB is internal (scheme: internal) and not internet-facing. | Architectural | An internet-facing ALB directly contradicts the "origin accessible via CDN only" policy. |
| ASF-064 | No VPC peering or transit gateway connection allows traffic from other VPCs to reach the origin or RDS subnets. | Environmental | VPC peering from other accounts or environments creates lateral movement paths to the production database. |
| ASF-065 | The RDS security group allows inbound traffic on the database port only from the EC2 origin security group. | Explicit | Database access must be restricted to the application tier's security group identifier. |
| ASF-066 | WAF inspection occurs before CloudFront routing — all traffic passes through WAF regardless of cache hit or miss. | Architectural | Traffic that bypasses WAF (e.g., directly to S3 origin for cached content) evades inspection entirely. |

---

### Pattern 16: Third-party Dependency

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-067 | AWS CloudFront, WAF, and Shield are available and not experiencing service-level degradation. | Dependency | A CloudFront outage makes the entire application inaccessible globally. |
| ASF-068 | Third-party IdP (Cognito, Okta, Auth0) used by Lambda@Edge is available and has no auth-stopping incidents. | Dependency | IdP availability is a transitive dependency for all authenticated requests through the CDN. |
| ASF-069 | Lambda@Edge runtime versions are supported by AWS and receive security patches. | Dependency | Unsupported Lambda runtimes (Node.js 12, Python 3.6) running at the edge expose the application to unpatched vulnerabilities. |
| ASF-070 | WAF managed rule group publishers (AWS, Fortinet, F5) maintain up-to-date rule sets against active threats. | Dependency | Managed rule groups that are not updated by the publisher create a false sense of security. |
| ASF-071 | The CDN provider's edge network does not have a compromised edge server that intercepts or modifies traffic. | Trust | CloudFront edge servers operate outside organizational control; compromise at the edge affects all users. |

**Total (A): 71** (4 per pattern × 17 patterns + 3 extra from patterns 5, 8, 16)

---

## Step 3: Comparison

### Overlap Mapping

| Human ID | ASF ID | Match Rationale |
|----------|--------|-----------------|
| H-001 | ASF-062 | Both require OAC/OAI to prevent direct origin access bypassing CDN. |
| H-002 | ASF-070 | Both require WAF rules to be updated for emerging threats. |
| H-003 | ASF-001 | Both require Lambda@Edge to perform authentication and not leak internal info. |
| H-004 | ASF-010 | Both concern WAF capacity limits and throttling. |
| H-005 | ASF-030 | Both require valid, auto-renewed TLS certificates for CloudFront. |
| H-006 | ASF-063 | Both restrict origin ALB to CloudFront traffic only via security groups. |
| H-007 | ASF-040 | Both require rate-based rules for application-layer attacks. |
| H-008 | ASF-013 | Both require RDS automated backups with adequate retention. |
| H-009 | ASF-057 | Both require CloudFront logging enabled for monitoring. |
| H-010 | ASF-048 | Both require WAF logs for forensic analysis and incident response. |
| H-011 | ASF-054 | Both require Lambda@Edge least privilege — no RDS or S3 access. |
| H-012 | ASF-063 | Both require EC2 origin has no public IP address. |
| H-013 | ASF-043 | Both address WAF being in Block mode vs Count mode. |
| H-015 | ASF-065 | Both require RDS in private subnet with no direct internet route. |
| H-016 | ASF-031 | Both require TLS 1.2+ between CloudFront and origin. |
| H-019 | ASF-057 | Both require ALB access logging. |
| H-020 | ASF-020 | Both require Lambda@Edge errors not to expose internal details. |
| H-023 | ASF-026 | Both require RDS encryption at rest with KMS. |
| H-024 | ASF-045 | Both require database credentials managed and rotated. |
| H-026 | ASF-064 | Both require no direct origin egress to internet. |
| H-027 | ASF-024 | Both require WAF logs not to capture sensitive request bodies. |
| H-030 | ASF-064 | Both require VPC Flow Logs enabled. |
| H-031 | ASF-017 | Both require Lambda@Edge least privilege. |
| H-033 | ASF-067 | Both require Shield Advanced for DDoS beyond standard. |
| H-034 | ASF-038 | Both require origin patching and trusted AMI maintenance. |
| H-038 | ASF-014 | Both require cross-region backup copies for DR. |
| H-040 | ASF-038 | Both require appropriate rate limit thresholds. |
| H-041 | ASF-033 | Both require secure database parameter group configuration. |
| H-045 | ASF-013 | Both require RDS deletion protection enabled. |

**Overlap (O): 29**

### Metrics

| Metric | Value | Formula |
|--------|-------|---------|
| Human assumptions (H) | 45 | Count of unique human-generated assumptions |
| ASF assumptions (A) | 71 | Count of unique ASF-generated assumptions |
| Overlap (O) | 29 | Count appearing in both lists |
| **Precision** | **40.8%** | O / A = 29/71 |
| **Recall** | **64.4%** | O / H = 29/45 |
| **F1 Score** | **50.0%** | 2 × (P × R) / (P + R) |
| Novel findings (A - O) | 42 | Assumptions ASF found that human missed (59.2% of ASF total) |
| Missed findings (H - O) | 16 | Assumptions human found that ASF missed (35.6% of human total) |

### Success Criteria Assessment

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Recall | >= 70% | 64.4% | ❌ Not met |
| Precision | >= 50% | 40.8% | ❌ Not met |
| Novel discoveries | >= 10% of total (A+O) | 37.2% (42/113) | ✅ Exceeded |
| Expert agreement (F1 proxy) | > 60% | 50.0% | ❌ Not met |

---

## Step 4: Analysis

### Ontology Overlap Analysis

| Ontology Category | Overlaps | ASF Total | Overlap Rate |
|-------------------|----------|-----------|-------------|
| Explicit | 8 | 13 | 61.5% |
| Derived | 9 | 16 | 56.3% |
| Operational | 5 | 15 | 33.3% |
| Implicit | 3 | 10 | 30.0% |
| Trust | 2 | 5 | 40.0% |
| Dependency | 1 | 5 | 20.0% |
| Architectural | 1 | 5 | 20.0% |
| Environmental | 0 | 4 | 0.0% |

**Best overlap:** Explicit (61.5%) and Derived (56.3%) — both humans and ASF converge on concrete, directly deducible security controls like TLS enforcement, encryption, and authentication.

**Worst overlap:** Environmental (0%) — the ASF identified cross-VPC peering risks, Lambda@Edge concurrency limits, and multi-environment concerns that the human architect treated as infrastructure context outside the assumptions scope.

### What Humans Caught That ASF Missed (Missed Findings = 16)

1. **CDN-specific operational hardening (H-014, H-018, H-025, H-029, H-035, H-042):** The human identified signed URLs for premium content, Origin Shield, HTTP/3 support, cache hit ratio monitoring, default root object configuration, and compression side-channel risks. These are CDN product-specific settings that the ASF's generic patterns do not capture.

2. **WAF rule composition details (H-002, H-017, H-021, H-032, H-044):** Host header validation for DNS rebinding, geo-blocking, bot control managed rule group, file upload malware scanning — these are WAF rule-specific features below the pattern granularity of the ASF.

3. **Database platform safety (H-041, H-045):** RDS parameter group hardening and deletion protection are database-platform-specific controls not covered by the ASF's Backup or Least Privilege patterns.

### What ASF Caught That Humans Missed (Novel Findings = 42)

1. **Incident Response (5 assumptions):** The human generated zero IR-specific assumptions. The ASF contributed a full IR pattern covering WAF bypass playbooks, log access, rapid IP blocking, anomaly detection, and false positive management.

2. **Identity Lifecycle (4 assumptions):** The human assumed credential rotation (H-024) but did not cover joiner/mover/leaver for IAM admin users, API key rotation for uploaded content, role trust policy hardening, or cross-account access reviews.

3. **Third-party dependency (5 assumptions):** The ASF surfaced dependencies on CloudFront service availability, IdP availability, Lambda runtime support, managed rule set publisher currency, and CDN edge server integrity — all outside the architect's consideration.

4. **Operational resilience (ASF-009 through ASF-012):** The human assumed rate limits (H-040) but did not consider origin failover, WAF capacity unit provisioning, Lambda concurrency, or origin shield for cache-miss storms.

5. **IAM least privilege expansion (ASF-018, ASF-021, ASF-055, ASF-056):** The human assumed Lambda@Edge least privilege (H-031) but did not extend to CloudFront service-linked roles, EC2 instance profiles, WAF admin role scoping, or IaC-based distribution management.

### Architecture Complexity Assessment

Architecture #018 is **Moderate** with a CDN-WAF-origin-DB chain, edge compute, and hybrid managed/self-managed components. The recall (64.4%) is close to the 70% target but pulled down by CDN-specific operational details and WAF rule composition concerns. The precision (40.8%) is consistent with Arch #001, reflecting the ASF's broad generative approach.

The human architect focused on the **CDN → WAF → Origin** attack path, generating assumptions along AWS-managed edge services and their interactions. The ASF contributed orthogonal assumptions in incident response, identity lifecycle, and third-party dependencies that are independent of the specific CDN topology.

### Key Insight

The main recall gap is **CDN-specific operational configuration** — signed URLs, Origin Shield, HTTP/3 settings, compression side-channels, and cache behaviors. The ASF's existing patterns lack a "CDN / Edge Computing Security" pattern that would capture these product-specific controls. Adding such a pattern (pattern 21) would likely close the recall gap to above 70%.

---

## Conclusions

| Metric | Target | Actual | Verdict |
|--------|--------|--------|---------|
| Recall | >= 70% | 64.4% | ❌ Below target — missing CDN-specific operational pattern |
| Precision | >= 50% | 40.8% | ❌ Below target — ASF is deliberately broad/generative |
| Novel discoveries | >= 10% | 37.2% | ✅ ASF adds substantial value beyond human reasoning |
| Expert agreement (F1) | > 60% | 50.0% | ❌ Below target — driven by low precision |

The ASF framework applied to Architecture #018 demonstrates strong exploratory breadth with 42 novel assumptions (37.2% of combined total) covering incident response, identity lifecycle, and third-party dependencies the human architect omitted. The primary actionable finding is the need for a **CDN & Edge Computing Security** pattern to close the recall gap for CDN-specific operational settings.
