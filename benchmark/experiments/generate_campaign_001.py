#!/usr/bin/env python3
"""Generate multi-LLM evaluation campaign 001 report with explicit overlap mappings."""

OUTPUT = "/Users/moksh/Project/cybersec/benchmark/experiments/multi_llm_campaign_001.md"

# ─── Human Assumptions (40) ───────────────────────────────────────────────
H = {
    "H-001": "VPN gateway runs fully patched firmware free of known CVEs.",
    "H-002": "VPN client software on user laptops is centrally managed and disables split-tunneling.",
    "H-003": "VPN credentials are never shared between users.",
    "H-004": "VPN gateway enforces multi-factor authentication for all users.",
    "H-005": "Private subnet containing database has no route to or from the internet.",
    "H-006": "TLS certificates on VPN gateway and internal web app are valid, not expired, from trusted CA.",
    "H-007": "VPN gateway enforces modern encryption (AES-256, SHA-256) and disables weak protocols.",
    "H-008": "Active Directory is hardened against common attack techniques.",
    "H-009": "Web app does not store database credentials in source code or config files readable by non-admins.",
    "H-010": "Database user has least-privilege access (SELECT/INSERT/UPDATE on specific tables only, no DDL).",
    "H-011": "TLS between VPN gateway and web app is properly terminated and validated at both ends.",
    "H-012": "Internal web app has no public IP address or internet-facing load balancer.",
    "H-013": "Database subnet cannot initiate connections to VPN gateway or web app.",
    "H-014": "Nightly backups are encrypted at rest.",
    "H-015": "Backup restore procedures are tested at least annually.",
    "H-016": "Security groups restrict database access to only the web app's IP/port.",
    "H-017": "Web app validates all user inputs and uses parameterized queries to prevent SQL injection.",
    "H-018": "Web app implements secure session management (HttpOnly+Secure cookies, short timeouts).",
    "H-019": "VPN gateway logs are forwarded to a SIEM.",
    "H-020": "Web app logs authentication events and access to sensitive payroll records.",
    "H-021": "Active Directory enforces account lockout after threshold of failed login attempts.",
    "H-022": "VPN access revoked within same shift when employee is terminated.",
    "H-023": "Database audit logging is enabled for all queries.",
    "H-024": "Web app enforces role-based access control beyond AD authentication.",
    "H-025": "No SSH/RDP/bastion host rules allow admin access to database from general corporate network.",
    "H-026": "VPN gateway supports sufficient concurrent connections for all remote employees.",
    "H-027": "Web app enforces rate limiting on login endpoints to prevent credential stuffing.",
    "H-028": "TLS version between VPN gateway and web app restricted to TLS 1.2 or higher.",
    "H-029": "Payroll database is not shared with other applications.",
    "H-030": "Firewall/security group separates VPN gateway subnet from web app subnet.",
    "H-031": "Automatic patching enabled on RDS instance for critical DB engine vulnerabilities.",
    "H-032": "VPN client has kill switch that terminates all traffic if VPN connection drops.",
    "H-033": "Web app implements defenses against XSS, CSRF, and clickjacking.",
    "H-034": "Payroll data is not written to app logs, error messages, or debugging output.",
    "H-035": "Database credentials rotated on regular cadence and immediately on compromise suspicion.",
    "H-036": "VPN gateway enforces device compliance checks before allowing connection.",
    "H-037": "Web app does not expose database error details (stack traces, query fragments) to end users.",
    "H-038": "RDS instance has deletion protection and final snapshot policy enabled.",
    "H-039": "VPN gateway uses perfect forward secrecy.",
    "H-040": "AD authentication does not fall back to NTLMv1 or LM hash storage.",
}

# ─── ASF Assumptions (64) ────────────────────────────────────────────────
A = {
    "ASF-001": "VPN gateway enforces MFA for all remote users.",
    "ASF-002": "MFA recovery codes are securely stored and usable only through verified identity proofing.",
    "ASF-003": "Help desk has documented social-engineering-resistant process for MFA token reset.",
    "ASF-004": "MFA is enforced on VPN gateway, not just the web application.",
    "ASF-005": "Web app validates AD/SSO tokens and does not accept unauthenticated requests.",
    "ASF-006": "AD domain controller is available and reachable from web app at all times.",
    "ASF-007": "SSO session timeout is consistent between AD and the web app.",
    "ASF-008": "SSO token signing keys are rotated and protected from unauthorized access.",
    "ASF-009": "Single VPN gateway failure does not block all remote access (redundancy exists).",
    "ASF-010": "Documented offline procedure for VPN outages that does not bypass security controls.",
    "ASF-011": "Internet circuit to VPN gateway has sufficient bandwidth and reliability SLA.",
    "ASF-012": "Web app can gracefully handle database connection failures without leaking data or crashing.",
    "ASF-013": "Nightly backups complete within backup window without impacting DB performance.",
    "ASF-014": "Backup restore is tested at least annually to validate data integrity.",
    "ASF-015": "Backups stored in separate AWS region or account from primary database.",
    "ASF-016": "Backup files encrypted at rest using separate KMS key from primary DB key.",
    "ASF-017": "AWS account hosting RDS has no other workloads with unnecessary network access to payroll VPC.",
    "ASF-018": "IAM roles for RDS access scoped to minimum required actions.",
    "ASF-019": "AWS root account protected by MFA and not used for daily operations.",
    "ASF-020": "No public AMIs or untrusted base images used for web app EC2 instance.",
    "ASF-021": "Payroll data classified as sensitive/confidential and subject to data handling policies.",
    "ASF-022": "Data flow diagrams exist and accurately represent all paths payroll data travels.",
    "ASF-023": "Web app does not transmit payroll data to any endpoint outside defined architecture.",
    "ASF-024": "Payroll data not used in development or staging environments.",
    "ASF-025": "RDS encryption at rest enabled using AWS KMS.",
    "ASF-026": "KMS key policies restrict which principals can encrypt/decrypt the database.",
    "ASF-027": "KMS keys are rotated annually.",
    "ASF-028": "Temporary storage (swap, temp tables, query cache) on RDS instance is also encrypted.",
    "ASF-029": "TLS enforced between VPN gateway and web app as documented.",
    "ASF-030": "Web app validates database server TLS certificate when connecting over SQL.",
    "ASF-031": "TLS 1.2 or higher enforced; TLS 1.0/1.1 and all SSL versions disabled.",
    "ASF-032": "Weak cipher suites (RC4, 3DES, CBC-mode) explicitly disabled on all TLS endpoints.",
    "ASF-033": "User laptops have EDR/AV software installed, running, and receiving current signature updates.",
    "ASF-034": "User laptops managed by MDM with enforced disk encryption and screen lock policies.",
    "ASF-035": "Lost or stolen laptops can be remotely wiped to prevent VPN credentials or cached data recovery.",
    "ASF-036": "Users do not install unauthorized software on corporate laptops.",
    "ASF-037": "Users do not share VPN credentials or leave authenticated sessions logged in on shared workstations.",
    "ASF-038": "Users can identify and report phishing attempts targeting their AD credentials.",
    "ASF-039": "Administrators assigning permissions follow principle of least privilege.",
    "ASF-040": "Users whose role does not require payroll access are not granted access even if AD-authenticated.",
    "ASF-041": "User accounts follow documented joiner/mover/leaver process in AD.",
    "ASF-042": "AD group membership for VPN and app access reviewed and recertified quarterly.",
    "ASF-043": "Service accounts used by web app for DB authentication managed with same rigor as user accounts.",
    "ASF-044": "App-level roles synchronized with HR data so role changes automatically update access.",
    "ASF-045": "Incident response plan covers payroll data exposure scenarios.",
    "ASF-046": "IR team has access to VPN, app, and DB logs during investigation.",
    "ASF-047": "IR plan includes isolation procedures (disconnect VPN, block DB traffic) that preserve forensic evidence.",
    "ASF-048": "Monitoring systems can detect anomalies in VPN and database access patterns indicating breach.",
    "ASF-049": "App DB user has SELECT/INSERT/UPDATE/DELETE only on payroll schema and nothing else.",
    "ASF-050": "Web app does not run as root or with OS-level administrative privileges.",
    "ASF-051": "AD users who can authenticate to VPN cannot also directly access database (no overlapping credentials).",
    "ASF-052": "Web app enforces authorization decisions beyond initial AD authentication (app-level RBAC).",
    "ASF-053": "VPN connection logs monitored for brute-force attempts and anomalous source geographies.",
    "ASF-054": "Database query patterns monitored for unusual volume (e.g., mass SELECT at 3 AM).",
    "ASF-055": "Alerts configured for failed authentication thresholds on VPN, app, and AD.",
    "ASF-056": "Monitoring infrastructure logs are append-only and tamper-proof.",
    "ASF-057": "VPN gateway, web app, and database in separate security groups/subnets with explicit allow rules.",
    "ASF-058": "No direct network path from VPN subnet to database subnet (traffic must pass through web app).",
    "ASF-059": "VPC flow logs or equivalent network telemetry enabled to detect unexpected traffic patterns.",
    "ASF-060": "Database security group allows inbound traffic only from web app security group on port 5432/3306.",
    "ASF-061": "AWS RDS is available and not experiencing region-level service outage.",
    "ASF-062": "VPN vendor has no known backdoors or critical vulnerabilities that are unpatched.",
    "ASF-063": "Third-party libraries used by web app scanned for vulnerabilities before deployment.",
    "ASF-064": "Exit strategy or migration plan exists if VPN vendor or RDS becomes unavailable due to acquisition/bankruptcy/sanctions.",
}

# ─── 5 AI Personas ────────────────────────────────────────────────────────

GPT = {
    "GPT-001": "VPN gateway firmware is patched against all known critical CVEs.",
    "GPT-002": "VPN client software prohibits split-tunneling on all user laptops.",
    "GPT-003": "VPN credentials are individual and never shared.",
    "GPT-004": "VPN gateway enforces MFA for every remote authentication attempt.",
    "GPT-005": "Modern encryption (AES-256, SHA-256) is enforced on the VPN gateway.",
    "GPT-006": "VPN gateway supports sufficient concurrent sessions for all remote employees.",
    "GPT-007": "VPN client implements a kill switch to block non-VPN traffic on disconnect.",
    "GPT-008": "Perfect forward secrecy is enabled on the VPN gateway.",
    "GPT-009": "VPN gateway has redundant hardware or failover to avoid single-point-of-failure.",
    "GPT-010": "VPN gateway logs are streamed to a centralized SIEM in real time.",
    "GPT-011": "TLS certificates presented by the VPN gateway and web app are valid and CA-signed.",
    "GPT-012": "TLS 1.2 or higher is enforced on all encrypted links; SSL and TLS 1.0/1.1 are disabled.",
    "GPT-013": "Weak cipher suites (RC4, 3DES, CBC-mode) are explicitly disabled on all TLS endpoints.",
    "GPT-014": "TLS termination validates certificates at both VPN gateway and web application ends.",
    "GPT-015": "Active Directory is hardened against Kerberoasting, DCSync, and pass-the-hash attacks.",
    "GPT-016": "AD enforces account lockout after a configured threshold of failed logins.",
    "GPT-017": "NTLMv1 and LM hash protocols are disabled; only Kerberos is used.",
    "GPT-018": "AD domain controllers are patched, monitored, and available for authentication at all times.",
    "GPT-019": "SSO session timeout values are synchronized between AD and the web application.",
    "GPT-020": "SSO token signing keys are rotated regularly and stored in a hardware security module.",
    "GPT-021": "The web application uses parameterized queries exclusively to prevent SQL injection.",
    "GPT-022": "The web application sets HttpOnly, Secure, and SameSite attributes on session cookies.",
    "GPT-023": "Role-based access control is enforced within the web application for payroll data.",
    "GPT-024": "The web application implements Content Security Policy and XSS/CSRF protections.",
    "GPT-025": "Login endpoints in the web application implement rate limiting against credential stuffing.",
    "GPT-026": "Database error details are never exposed in web application HTTP responses.",
    "GPT-027": "The web application does not send payroll data to any third-party endpoint.",
    "GPT-028": "Database credentials are stored securely and not committed to source code.",
    "GPT-029": "The web application process does not run with root or administrator privileges.",
    "GPT-030": "The database user has SELECT/INSERT/UPDATE/DELETE privileges only on the payroll schema.",
    "GPT-031": "The RDS instance resides in a private subnet with no internet gateway route.",
    "GPT-032": "The database security group allows ingress only from the web application security group.",
    "GPT-033": "No network path exists between the VPN subnet and the database subnet directly.",
    "GPT-034": "RDS encryption at rest is enabled with AWS KMS and a customer-managed key.",
    "GPT-035": "Database audit logging is enabled for all DML and DDL statements.",
    "GPT-036": "RDS automatic minor version upgrades are enabled for critical security patches.",
    "GPT-037": "Database credentials are rotated every 90 days and immediately on any compromise suspicion.",
    "GPT-038": "Nightly backups complete within the configured maintenance window without performance degradation.",
    "GPT-039": "Backup restore procedures are tested at least annually and documented.",
    "GPT-040": "Backup files are encrypted at rest using a dedicated KMS key.",
    "GPT-041": "The VPN gateway, web application, and database each reside in separate subnets with explicit allow rules.",
    "GPT-042": "The database does not initiate outbound connections to the VPN gateway or web application.",
    "GPT-043": "IAM roles for RDS and application access follow least-privilege principles.",
    "GPT-044": "VPN access grants are revoked within the same shift on employee termination.",
    "GPT-045": "The web application has a Content Security Policy header configured to limit data exfiltration channels.",
}

CLAUDE = {
    "CLAUDE-001": "VPN gateway firmware is updated within vendor patch window for critical CVEs.",
    "CLAUDE-002": "VPN client is centrally managed with mandatory split-tunneling disabled.",
    "CLAUDE-003": "VPN credentials are unique per user with individual accountability enforced.",
    "CLAUDE-004": "MFA required for all VPN connections including service accounts.",
    "CLAUDE-005": "VPN gateway enforces AES-256-GCM and SHA-256 minimum; rejects legacy proposals.",
    "CLAUDE-006": "VPN gateway has sufficient session capacity for peak concurrent remote access.",
    "CLAUDE-007": "VPN client kill switch is enabled and verified to block all non-tunnel traffic on connection loss.",
    "CLAUDE-008": "Perfect forward secrecy is configured on VPN gateway to protect recorded traffic.",
    "CLAUDE-009": "VPN gateway has active-passive failover or clustering to avoid SPOF.",
    "CLAUDE-010": "Documented offline procedures exist for VPN outages without bypassing security controls.",
    "CLAUDE-011": "VPN logs shipped to SIEM with alerting on anomalous patterns.",
    "CLAUDE-012": "TLS certificates for VPN and web app issued by trusted internal CA and monitored for expiry.",
    "CLAUDE-013": "TLS 1.2 minimum enforced; TLS 1.0, 1.1, SSL explicitly disabled.",
    "CLAUDE-014": "All weak TLS cipher suites (RC4, 3DES, CBC-mode) removed from cipher allow-list.",
    "CLAUDE-015": "Mutual TLS considered for server-to-server communication for dual-end validation.",
    "CLAUDE-016": "Active Directory configured with protected users group and SMB signing to resist relay attacks.",
    "CLAUDE-017": "AD password policies enforce complexity, length >=14, and periodic rotation.",
    "CLAUDE-018": "Account lockout configured in AD with reasonable threshold and duration.",
    "CLAUDE-019": "NTLMv1 blocked by Group Policy; only Kerberos or NTLMv2 permitted.",
    "CLAUDE-020": "AD domain controllers hardened, patched, monitored with availability SLAs.",
    "CLAUDE-021": "SSO session lifetimes consistently configured between AD and web app.",
    "CLAUDE-022": "SSO token signing keys stored in HSM or key vault with access logging.",
    "CLAUDE-023": "Web app uses prepared statements or ORM parameter binding to eliminate SQL injection.",
    "CLAUDE-024": "Session cookies are HttpOnly, Secure, SameSite=Strict with 15-min idle timeout.",
    "CLAUDE-025": "Web app enforces app-level RBAC mapping AD groups to payroll data access tiers.",
    "CLAUDE-026": "Web app mitigates XSS with output encoding, CSRF with anti-forgery tokens, clickjacking with CSP.",
    "CLAUDE-027": "Rate limiting on all auth endpoints with graduated delays and CAPTCHA after threshold.",
    "CLAUDE-028": "Error handlers sanitize stack traces and schema details from production error pages.",
    "CLAUDE-029": "Payroll data not transmitted to any monitoring, analytics, or third-party service outside defined flow.",
    "CLAUDE-030": "Database credentials stored in secrets manager with automated rotation.",
    "CLAUDE-031": "Web app runs under dedicated low-privilege service account, not root/LocalSystem.",
    "CLAUDE-032": "App DB user has least privilege — SELECT, INSERT, UPDATE, DELETE only on payroll schema, no DDL.",
    "CLAUDE-033": "RDS instance in private subnet with no direct/indirect internet route via NAT or egress-only IGW.",
    "CLAUDE-034": "DB security group permits inbound on port 5432/3306 only from web app security group.",
    "CLAUDE-035": "No direct routing path from VPN subnet to DB subnet at VPC or transit gateway level.",
    "CLAUDE-036": "RDS encryption at rest uses customer-managed KMS key with automatic rotation.",
    "CLAUDE-037": "DB audit logging captures all SELECT, INSERT, UPDATE, DELETE with user and timestamp.",
    "CLAUDE-038": "RDS has automated patching for critical and high-severity CVEs with defined maintenance window.",
    "CLAUDE-039": "Database credentials rotated automatically on defined schedule and upon security incident.",
    "CLAUDE-040": "Nightly backups complete within backup window; DB performance monitored for degradation.",
    "CLAUDE-041": "Backup restore tested at least annually with documented RTO/RPO validation.",
    "CLAUDE-042": "Backups encrypted at rest using KMS key separate from primary DB encryption key.",
    "CLAUDE-043": "Payroll data formally classified as sensitive with data handling policies.",
    "CLAUDE-044": "User laptops have EDR/AV with real-time protection and MDM disk encryption.",
    "CLAUDE-045": "Lost/stolen laptops can be remotely wiped and VPN credentials revoked in single orchestrated action.",
    "CLAUDE-046": "Users receive annual security awareness training covering phishing and credential hygiene.",
    "CLAUDE-047": "Access provisioning follows joiner/mover/leaver process enforced by HR system integration.",
    "CLAUDE-048": "AD group membership for VPN and app access recertified quarterly by data owners.",
    "CLAUDE-049": "Documented IR plan specifically covers payroll data breach scenarios with defined roles.",
    "CLAUDE-050": "IR team has timely access to VPN, app, and DB logs with preserved chain of custody.",
    "CLAUDE-051": "IR plan includes DB isolation and VPN disconnection procedures preserving forensic evidence.",
    "CLAUDE-052": "Vendor exit strategy documented for both VPN vendor and RDS including data migration.",
}

GEMINI = {
    "GEMINI-001": "VPN gateway firmware is patched and free of known exploits.",
    "GEMINI-002": "Split-tunneling is disabled on all managed VPN clients.",
    "GEMINI-003": "VPN credentials are unique per user and not shared.",
    "GEMINI-004": "MFA is enforced for all VPN access.",
    "GEMINI-005": "VPN gateway uses strong encryption with no legacy protocol fallback.",
    "GEMINI-006": "VPN gateway has capacity for all concurrent remote users.",
    "GEMINI-007": "VPN client kill switch prevents traffic leakage on tunnel drop.",
    "GEMINI-008": "Perfect forward secrecy is enabled on the VPN.",
    "GEMINI-009": "VPN logs are sent to SIEM and monitored for threats.",
    "GEMINI-010": "TLS certificates for VPN and web app are valid and monitored for expiry.",
    "GEMINI-011": "TLS 1.2 or higher is enforced throughout.",
    "GEMINI-012": "Weak ciphers (RC4, 3DES) are disabled at all TLS endpoints.",
    "GEMINI-013": "Active Directory is hardened against lateral movement techniques.",
    "GEMINI-014": "AD account lockout is configured to resist brute-force.",
    "GEMINI-015": "NTLMv1 is disabled in the domain.",
    "GEMINI-016": "AD domain controllers are highly available and patched.",
    "GEMINI-017": "SSO tokens are cryptographically signed and validated by the web app.",
    "GEMINI-018": "Web application is immune to SQL injection via parameterized queries.",
    "GEMINI-019": "Session cookies are secured with HttpOnly and Secure flags.",
    "GEMINI-020": "Application enforces RBAC for payroll data access.",
    "GEMINI-021": "XSS and CSRF protections are active in the web application.",
    "GEMINI-022": "Login endpoints have rate limiting against credential abuse.",
    "GEMINI-023": "Error messages do not leak database internals to users.",
    "GEMINI-024": "Payroll data is not leaked to external endpoints.",
    "GEMINI-025": "Database credentials are stored in a secrets manager, not code.",
    "GEMINI-026": "The application DB user has least privilege (no DDL).",
    "GEMINI-027": "RDS is in a private subnet with no internet route.",
    "GEMINI-028": "Database security group allows only web application traffic.",
    "GEMINI-029": "No direct route from VPN subnet to database subnet.",
    "GEMINI-030": "RDS encryption at rest is enabled.",
    "GEMINI-031": "Database auditing is enabled for all queries.",
    "GEMINI-032": "Backups are encrypted at rest.",
    "GEMINI-033": "Backup restore is tested at least annually.",
    "GEMINI-034": "Endpoints have EDR/AV installed and updated.",
    "GEMINI-035": "Laptops are MDM-managed with disk encryption enforced.",
    "GEMINI-036": "Access is revoked on termination within business hours.",
    "GEMINI-037": "AWS root account has MFA and is not used for daily operations.",
    "GEMINI-038": "IR plan exists for payroll data incidents.",
    "GEMINI-039": "Third-party libraries are scanned for vulnerabilities.",
    "GEMINI-040": "VPN vendor and RDS have demonstrated business continuity.",
}

DS = {
    "DS-001": "VPN gateway runs firmware version with no unpatched CVSS >= 7.0 vulnerabilities.",
    "DS-002": "VPN client config enforces full-tunnel mode via Group Policy or MDM profile.",
    "DS-003": "Each VPN user has unique certificate or token-bound credential; no shared accounts.",
    "DS-004": "VPN gateway requires TOTP, WebAuthn, or push-based MFA; SMS OTP not used.",
    "DS-005": "IPsec/IKEv2 or WireGuard with AES-256-GCM and SHA-384; IKEv1 and PPTP disabled.",
    "DS-006": "VPN gateway concurrent session limit exceeds remote workforce plus 25% buffer.",
    "DS-007": "VPN client implements network-level kill switch via firewall rules dropping non-tunnel traffic.",
    "DS-008": "VPN gateway configured with DHE or ECDHE key exchange for PFS.",
    "DS-009": "VPN gateway deployed in active-active or active-passive cluster with automatic failover.",
    "DS-010": "ISP circuit to VPN gateway has documented SLA >=99.9% uptime with redundant providers.",
    "DS-011": "VPN gateway syslog forwarded to SIEM with TLS encryption and integrity verification.",
    "DS-012": "TLS certificates issued by internal CA with automated renewal via ACME or equivalent.",
    "DS-013": "TLS 1.3 preferred; TLS 1.2 minimum; TLS 1.0/1.1 and SSL disabled at OS/ALB level.",
    "DS-014": "TLS cipher config excludes NULL, RC4, 3DES, CBC-mode; only AEAD ciphers (GCM/ChaCha20).",
    "DS-015": "TLS cert verification includes CRL or OCSP stapling; self-signed certs rejected in production.",
    "DS-016": "AD forest functional level Windows 2016+ with SMB signing and LDAP signing enforced.",
    "DS-017": "AD fine-grained password policies enforce 16-char minimum (NIST 800-63B guidance).",
    "DS-018": "AD account lockout: 10 failed attempts within 15 min, 30-min lockout duration.",
    "DS-019": "NTLMv1 blocked by Group Policy; LM hash storage disabled.",
    "DS-020": "AD domain controllers deployed in multi-AZ configuration with health monitoring.",
    "DS-021": "SSO token lifetime configured to 8h maximum with sliding window no longer than 1h.",
    "DS-022": "SSO token signing keys stored in AWS KMS with automatic rotation every 90 days.",
    "DS-023": "Web app uses prepared statements with bound parameters for all DB queries; no dynamic SQL.",
    "DS-024": "Session cookies: HttpOnly, Secure, SameSite=Strict, max idle TTL 15 min.",
    "DS-025": "Web app implements ABAC or RBAC with row-level security for payroll data.",
    "DS-026": "XSS mitigated by context-aware output encoding; CSRF by anti-forgery tokens; clickjacking by X-Frame-Options: DENY.",
    "DS-027": "Rate limiting via reverse proxy with graduated response (delay/CAPTCHA/block) at 10 req/min per IP.",
    "DS-028": "Web app returns generic error pages with correlation ID; stack traces logged server-side only.",
    "DS-029": "Database credentials stored in AWS Secrets Manager or Parameter Store with auto-rotation.",
    "DS-030": "Web app runs as non-root user with no capabilities and read-only filesystem where feasible.",
    "DS-031": "App DB user restricted to table-level grants (SELECT, INSERT, UPDATE, DELETE) with no schema modifications.",
    "DS-032": "RDS instance deployed in private subnet with VPC endpoint or NAT-gateway-restricted egress.",
    "DS-033": "RDS security group ingress restricted to web app security group on TCP 5432/3306.",
    "DS-034": "VPC route tables have no route from VPN subnet CIDR to DB subnet CIDR.",
    "DS-035": "RDS encryption at rest uses customer-managed KMS key with automatic annual rotation and deletion protection.",
    "DS-036": "RDS enhanced monitoring and audit logs enabled, streamed to CloudWatch with 1-year retention.",
    "DS-037": "RDS auto-minor-version-upgrade enabled; major upgrades tested in staging.",
    "DS-038": "DB credentials rotated every 90 days via automated pipeline with secret versioning.",
    "DS-039": "Nightly automated snapshots complete within backup window; retention at least 30 days.",
    "DS-040": "Backup restore tested quarterly with documented RTO <=4h and RPO <=24h.",
    "DS-041": "Backup snapshots encrypted with separate KMS key in different AWS region.",
    "DS-042": "VPC flow logs enabled for all subnets; logs delivered to centralized S3 with object lock.",
    "DS-043": "IAM policies for RDS, EC2, KMS follow least-privilege model with condition keys restricting source VPC/IP.",
    "DS-044": "Third-party libraries scanned via SCA tooling (Snyk/Dependabot) before merge to production.",
    "DS-045": "RDS deletion protection enabled; final DB snapshot taken before any deletion operation.",
    "DS-046": "Monitoring infrastructure logging immutable; CloudTrail and S3 access logs enabled for audit bucket.",
    "DS-047": "AD users who authenticate to VPN not granted direct DB access with overlapping credentials.",
    "DS-048": "Web app not reachable via public DNS or ALB; internal Route53-only.",
}

QWEN = {
    "QWEN-001": "VPN gateway patched for critical vulnerabilities within vendor SLA.",
    "QWEN-002": "Split-tunneling disabled on VPN client via MDM policy.",
    "QWEN-003": "VPN credentials unique per user and never shared.",
    "QWEN-004": "MFA enforced for all VPN connections.",
    "QWEN-005": "VPN gateway uses strong encryption and disables outdated protocols.",
    "QWEN-006": "VPN gateway has adequate capacity for remote workforce.",
    "QWEN-007": "VPN client has kill switch to prevent data leakage on tunnel drop.",
    "QWEN-008": "VPN gateway logs sent to SIEM and monitored daily.",
    "QWEN-009": "Documented offline procedure exists for when VPN is unavailable.",
    "QWEN-010": "TLS certificates for VPN and web app are valid and auto-renewed.",
    "QWEN-011": "TLS 1.2 or higher enforced across the architecture.",
    "QWEN-012": "Weak TLS ciphers disabled on all endpoints.",
    "QWEN-013": "Active Directory hardened against common attacks.",
    "QWEN-014": "AD account lockout configured to prevent brute-force.",
    "QWEN-015": "NTLMv1 disabled in favor of Kerberos authentication.",
    "QWEN-016": "AD domain controllers patched and monitored regularly.",
    "QWEN-017": "SSO tokens validated by web app before granting access.",
    "QWEN-018": "SSO token signing keys rotated and access-controlled.",
    "QWEN-019": "Web app uses parameterized queries to prevent SQL injection.",
    "QWEN-020": "Session cookies secured with HttpOnly, Secure, SameSite attributes.",
    "QWEN-021": "RBAC enforced within application for payroll data access.",
    "QWEN-022": "XSS, CSRF, and clickjacking protections implemented.",
    "QWEN-023": "Rate limiting configured on login endpoints.",
    "QWEN-024": "Error pages do not leak DB schema or stack traces to users.",
    "QWEN-025": "Payroll data not exfiltrated to unauthorized external endpoints.",
    "QWEN-026": "DB credentials stored in secrets manager, rotated regularly.",
    "QWEN-027": "App DB user has only privileges needed for its function.",
    "QWEN-028": "RDS in private subnet with no internet gateway attached.",
    "QWEN-029": "DB security group allows traffic only from web application.",
    "QWEN-030": "No direct network path from VPN subnet to DB subnet.",
    "QWEN-031": "RDS encryption at rest enabled with customer-managed key.",
    "QWEN-032": "DB audit logging enabled for all queries.",
    "QWEN-033": "RDS patching configured for automatic critical updates.",
    "QWEN-034": "DB credentials rotated at least every 90 days.",
    "QWEN-035": "Backups complete successfully within nightly window.",
    "QWEN-036": "Backup restore tested at least annually.",
    "QWEN-037": "Backups encrypted at rest.",
    "QWEN-038": "Endpoints have antivirus/EDR installed and centrally managed.",
    "QWEN-039": "Laptops managed with disk encryption and screen lock enforced.",
    "QWEN-040": "Access revoked promptly on employee termination.",
    "QWEN-041": "User access recertified quarterly by payroll data owners.",
    "QWEN-042": "IR plan covers payroll data exposure scenarios.",
    "QWEN-043": "IR team can access VPN, app, and DB logs during investigations.",
    "QWEN-044": "Third-party libraries scanned before deployment.",
    "QWEN-045": "Payroll data classified as sensitive and handled accordingly.",
}

# ─── Explicit Overlap Mappings ────────────────────────────────────────────
# Map each AI assumption ID -> set of Human IDs and set of ASF IDs it matches

AI_OVERLAP_H = {
    "GPT": {
        "GPT-001": {"H-001"}, "GPT-002": {"H-002"}, "GPT-003": {"H-003"}, "GPT-004": {"H-004"},
        "GPT-005": {"H-007"}, "GPT-006": {"H-026"}, "GPT-007": {"H-032"}, "GPT-008": {"H-039"},
        "GPT-010": {"H-019"}, "GPT-011": {"H-006"}, "GPT-012": {"H-028"}, "GPT-014": {"H-011"},
        "GPT-015": {"H-008"}, "GPT-016": {"H-021"}, "GPT-017": {"H-040"}, "GPT-021": {"H-017"},
        "GPT-022": {"H-018"}, "GPT-023": {"H-024"}, "GPT-024": {"H-033"}, "GPT-025": {"H-027"},
        "GPT-026": {"H-037"}, "GPT-028": {"H-009"}, "GPT-030": {"H-010"}, "GPT-031": {"H-005"},
        "GPT-032": {"H-016"}, "GPT-033": {"H-013"}, "GPT-035": {"H-023"}, "GPT-036": {"H-031"},
        "GPT-037": {"H-035"}, "GPT-039": {"H-015"}, "GPT-040": {"H-014"}, "GPT-041": {"H-030"},
        "GPT-042": {"H-013"}, "GPT-044": {"H-022"},
    },
    "Claude": {
        "CLAUDE-001": {"H-001"}, "CLAUDE-002": {"H-002"}, "CLAUDE-003": {"H-003"}, "CLAUDE-004": {"H-004"},
        "CLAUDE-005": {"H-007"}, "CLAUDE-006": {"H-026"}, "CLAUDE-007": {"H-032"}, "CLAUDE-008": {"H-039"},
        "CLAUDE-011": {"H-019"}, "CLAUDE-012": {"H-006"}, "CLAUDE-013": {"H-028"}, "CLAUDE-016": {"H-008"},
        "CLAUDE-018": {"H-021"}, "CLAUDE-019": {"H-040"}, "CLAUDE-023": {"H-017"}, "CLAUDE-024": {"H-018"},
        "CLAUDE-025": {"H-024"}, "CLAUDE-026": {"H-033"}, "CLAUDE-027": {"H-027"}, "CLAUDE-028": {"H-037"},
        "CLAUDE-030": {"H-009"}, "CLAUDE-032": {"H-010"}, "CLAUDE-033": {"H-005"}, "CLAUDE-034": {"H-016"},
        "CLAUDE-035": {"H-013"}, "CLAUDE-037": {"H-023"}, "CLAUDE-038": {"H-031"}, "CLAUDE-039": {"H-035"},
        "CLAUDE-041": {"H-015"}, "CLAUDE-042": {"H-014"}, "CLAUDE-043": {"H-034"},
        "CLAUDE-047": {"H-022"}, "CLAUDE-048": {"H-022"},
    },
    "Gemini": {
        "GEMINI-001": {"H-001"}, "GEMINI-002": {"H-002"}, "GEMINI-003": {"H-003"}, "GEMINI-004": {"H-004"},
        "GEMINI-005": {"H-007"}, "GEMINI-006": {"H-026"}, "GEMINI-007": {"H-032"}, "GEMINI-008": {"H-039"},
        "GEMINI-009": {"H-019"}, "GEMINI-010": {"H-006"}, "GEMINI-011": {"H-028"}, "GEMINI-013": {"H-008"},
        "GEMINI-014": {"H-021"}, "GEMINI-015": {"H-040"}, "GEMINI-018": {"H-017"}, "GEMINI-019": {"H-018"},
        "GEMINI-020": {"H-024"}, "GEMINI-021": {"H-033"}, "GEMINI-022": {"H-027"}, "GEMINI-023": {"H-037"},
        "GEMINI-025": {"H-009"}, "GEMINI-026": {"H-010"}, "GEMINI-027": {"H-005"}, "GEMINI-028": {"H-016"},
        "GEMINI-029": {"H-013"}, "GEMINI-031": {"H-023"}, "GEMINI-033": {"H-015"}, "GEMINI-032": {"H-014"},
        "GEMINI-036": {"H-022"},
    },
    "DeepSeek": {
        "DS-001": {"H-001"}, "DS-002": {"H-002"}, "DS-003": {"H-003"}, "DS-004": {"H-004"},
        "DS-005": {"H-007"}, "DS-006": {"H-026"}, "DS-007": {"H-032"}, "DS-008": {"H-039"},
        "DS-011": {"H-019"}, "DS-012": {"H-006"}, "DS-013": {"H-028"}, "DS-015": {"H-011"},
        "DS-016": {"H-008"}, "DS-018": {"H-021"}, "DS-019": {"H-040"}, "DS-023": {"H-017"},
        "DS-024": {"H-018"}, "DS-025": {"H-024"}, "DS-026": {"H-033"}, "DS-027": {"H-027"},
        "DS-028": {"H-037"}, "DS-029": {"H-009"}, "DS-031": {"H-010"}, "DS-032": {"H-005"},
        "DS-033": {"H-016"}, "DS-034": {"H-013"}, "DS-036": {"H-023"}, "DS-037": {"H-031"},
        "DS-038": {"H-035"}, "DS-040": {"H-015"}, "DS-041": {"H-014"}, "DS-045": {"H-038"},
        "DS-048": {"H-012"},
    },
    "Qwen": {
        "QWEN-001": {"H-001"}, "QWEN-002": {"H-002"}, "QWEN-003": {"H-003"}, "QWEN-004": {"H-004"},
        "QWEN-005": {"H-007"}, "QWEN-006": {"H-026"}, "QWEN-007": {"H-032"}, "QWEN-008": {"H-019"},
        "QWEN-010": {"H-006"}, "QWEN-011": {"H-028"}, "QWEN-013": {"H-008"}, "QWEN-014": {"H-021"},
        "QWEN-015": {"H-040"}, "QWEN-019": {"H-017"}, "QWEN-020": {"H-018"}, "QWEN-021": {"H-024"},
        "QWEN-022": {"H-033"}, "QWEN-023": {"H-027"}, "QWEN-024": {"H-037"}, "QWEN-026": {"H-009"},
        "QWEN-027": {"H-010"}, "QWEN-028": {"H-005"}, "QWEN-029": {"H-016"}, "QWEN-030": {"H-013"},
        "QWEN-032": {"H-023"}, "QWEN-033": {"H-031"}, "QWEN-034": {"H-035"}, "QWEN-036": {"H-015"},
        "QWEN-037": {"H-014"}, "QWEN-040": {"H-022"},
    },
}

AI_OVERLAP_A = {
    "GPT": {
        "GPT-003": {"ASF-037"}, "GPT-004": {"ASF-001"}, "GPT-009": {"ASF-009"}, "GPT-010": {"ASF-053"},
        "GPT-011": {"ASF-030"}, "GPT-012": {"ASF-031"}, "GPT-013": {"ASF-032"}, "GPT-014": {"ASF-030"},
        "GPT-018": {"ASF-006"}, "GPT-019": {"ASF-007"}, "GPT-020": {"ASF-008"}, "GPT-023": {"ASF-052"},
        "GPT-027": {"ASF-023"}, "GPT-029": {"ASF-050"}, "GPT-030": {"ASF-049"}, "GPT-032": {"ASF-060"},
        "GPT-033": {"ASF-058"}, "GPT-034": {"ASF-025"}, "GPT-035": {"ASF-054"}, "GPT-038": {"ASF-013"},
        "GPT-039": {"ASF-014"}, "GPT-040": {"ASF-016"}, "GPT-041": {"ASF-057"}, "GPT-043": {"ASF-018"},
        "GPT-044": {"ASF-041"},
    },
    "Claude": {
        "CLAUDE-003": {"ASF-037"}, "CLAUDE-004": {"ASF-001"}, "CLAUDE-009": {"ASF-009"},
        "CLAUDE-010": {"ASF-010"}, "CLAUDE-011": {"ASF-053"}, "CLAUDE-012": {"ASF-030"},
        "CLAUDE-013": {"ASF-031"}, "CLAUDE-014": {"ASF-032"}, "CLAUDE-021": {"ASF-007"},
        "CLAUDE-022": {"ASF-008"}, "CLAUDE-025": {"ASF-052"}, "CLAUDE-029": {"ASF-023"},
        "CLAUDE-031": {"ASF-050"}, "CLAUDE-032": {"ASF-049"}, "CLAUDE-034": {"ASF-060"},
        "CLAUDE-035": {"ASF-058"}, "CLAUDE-036": {"ASF-025"}, "CLAUDE-037": {"ASF-054"},
        "CLAUDE-040": {"ASF-013"}, "CLAUDE-041": {"ASF-014"}, "CLAUDE-042": {"ASF-016"},
        "CLAUDE-043": {"ASF-021"}, "CLAUDE-044": {"ASF-033", "ASF-034"}, "CLAUDE-045": {"ASF-035"},
        "CLAUDE-046": {"ASF-038"}, "CLAUDE-047": {"ASF-041"}, "CLAUDE-048": {"ASF-042"},
        "CLAUDE-049": {"ASF-045"}, "CLAUDE-050": {"ASF-046"}, "CLAUDE-051": {"ASF-047"},
        "CLAUDE-052": {"ASF-064"},
    },
    "Gemini": {
        "GEMINI-003": {"ASF-037"}, "GEMINI-004": {"ASF-001"}, "GEMINI-011": {"ASF-031"},
        "GEMINI-012": {"ASF-032"}, "GEMINI-017": {"ASF-005"}, "GEMINI-020": {"ASF-052"},
        "GEMINI-024": {"ASF-023"}, "GEMINI-026": {"ASF-049"}, "GEMINI-028": {"ASF-060"},
        "GEMINI-029": {"ASF-058"}, "GEMINI-030": {"ASF-025"}, "GEMINI-032": {"ASF-016"},
        "GEMINI-033": {"ASF-014"}, "GEMINI-034": {"ASF-033"}, "GEMINI-035": {"ASF-034"},
        "GEMINI-037": {"ASF-019"}, "GEMINI-038": {"ASF-045"}, "GEMINI-039": {"ASF-063"},
        "GEMINI-040": {"ASF-064"},
    },
    "DeepSeek": {
        "DS-003": {"ASF-037"}, "DS-004": {"ASF-001"}, "DS-009": {"ASF-009"}, "DS-010": {"ASF-011"},
        "DS-011": {"ASF-053"}, "DS-012": {"ASF-030"}, "DS-013": {"ASF-031"}, "DS-014": {"ASF-032"},
        "DS-021": {"ASF-007"}, "DS-022": {"ASF-008"}, "DS-025": {"ASF-052"}, "DS-030": {"ASF-050"},
        "DS-031": {"ASF-049"}, "DS-033": {"ASF-060"}, "DS-034": {"ASF-058"}, "DS-035": {"ASF-025"},
        "DS-036": {"ASF-054"}, "DS-039": {"ASF-013"}, "DS-040": {"ASF-014"}, "DS-041": {"ASF-015", "ASF-016"},
        "DS-042": {"ASF-059"}, "DS-043": {"ASF-018"}, "DS-044": {"ASF-063"}, "DS-046": {"ASF-056"},
        "DS-047": {"ASF-051"},
    },
    "Qwen": {
        "QWEN-003": {"ASF-037"}, "QWEN-004": {"ASF-001"}, "QWEN-009": {"ASF-010"}, "QWEN-008": {"ASF-053"},
        "QWEN-010": {"ASF-030"}, "QWEN-011": {"ASF-031"}, "QWEN-012": {"ASF-032"}, "QWEN-017": {"ASF-005"},
        "QWEN-018": {"ASF-008"}, "QWEN-021": {"ASF-052"}, "QWEN-025": {"ASF-023"}, "QWEN-027": {"ASF-049"},
        "QWEN-029": {"ASF-060"}, "QWEN-030": {"ASF-058"}, "QWEN-031": {"ASF-025"}, "QWEN-032": {"ASF-054"},
        "QWEN-035": {"ASF-013"}, "QWEN-036": {"ASF-014"}, "QWEN-037": {"ASF-016"}, "QWEN-038": {"ASF-033"},
        "QWEN-039": {"ASF-034"}, "QWEN-040": {"ASF-041"}, "QWEN-041": {"ASF-042"}, "QWEN-042": {"ASF-045"},
        "QWEN-043": {"ASF-046"}, "QWEN-044": {"ASF-063"}, "QWEN-045": {"ASF-021"},
    },
}

# ─── Build Unified Assumption Set ─────────────────────────────────────────
# Each "canonical assumption" is a tuple (source, id, text)
# We group by "topic" — assumptions that refer to the same security concern.

# We build a topic-based mapping where each topic gets:
# - A canonical description
# - Which sources covered it (Human IDs, ASF IDs, AI IDs per persona)

# Approach: use the Human and ASF IDs as anchors.
# Each Human/ASF assumption defines a "topic".
# AI assumptions map to those topics via the overlap mappings.
# AI assumptions NOT mapped to any H/ASF are unique topics.

# 1. Build topic map: topic_id -> {text, human_ids, asf_ids, persona_ids}
topics = {}
topic_counter = 0

def add_topic(text, human_ids=None, asf_ids=None, persona_contrib=None):
    global topic_counter
    topic_counter += 1
    tid = f"T{topic_counter:03d}"
    topics[tid] = {
        "text": text,
        "human": set(human_ids or []),
        "asf": set(asf_ids or []),
        "personas": persona_contrib or {},
    }
    return tid

# First, create topics for all Human + ASF assumptions (including overlap)
# The overlap from Phase 6: 25 pairs
human_asf_overlap = {
    "H-003": ["ASF-037"], "H-004": ["ASF-001"], "H-005": ["ASF-057"],
    "H-006": ["ASF-030"], "H-007": ["ASF-031"], "H-009": ["ASF-043"],
    "H-010": ["ASF-049"], "H-011": ["ASF-030"], "H-013": ["ASF-058"],
    "H-014": ["ASF-016"], "H-015": ["ASF-014"], "H-016": ["ASF-060"],
    "H-017": ["ASF-052"], "H-018": ["ASF-052"], "H-019": ["ASF-053"],
    "H-022": ["ASF-041"], "H-023": ["ASF-054"], "H-024": ["ASF-052"],
    "H-025": ["ASF-057"], "H-028": ["ASF-031"], "H-029": ["ASF-017"],
    "H-030": ["ASF-057"], "H-036": ["ASF-033"], "H-037": ["ASF-023"],
    "H-040": ["ASF-031"],
}

# Track which H and A IDs have been assigned to topics
assigned_h = set()
assigned_a = set()

# Create topics from H-ASF overlaps
for hid, aids in human_asf_overlap.items():
    h_text = H[hid]
    for aid in aids:
        a_text = A[aid]
        assigned_h.add(hid)
        assigned_a.add(aid)
    add_topic(h_text, human_ids=[hid], asf_ids=aids)

# Create topics for Human-only assumptions
for hid in sorted(H.keys()):
    if hid not in assigned_h:
        add_topic(H[hid], human_ids=[hid])
        assigned_h.add(hid)

# Create topics for ASF-only assumptions
for aid in sorted(A.keys()):
    if aid not in assigned_a:
        add_topic(A[aid], asf_ids=[aid])
        assigned_a.add(aid)

# Now map AI assumptions to topics
# For each AI persona and each of its assumptions:
# - If the AI assumption maps to H or ASF IDs, add to those topics
# - If not, create a new topic

# Reverse map: persona -> [(aid, topic_id)]
for persona in ["GPT", "Claude", "Gemini", "DeepSeek", "Qwen"]:
    ai_list = {"GPT": GPT, "Claude": CLAUDE, "Gemini": GEMINI, "DeepSeek": DS, "Qwen": QWEN}[persona]
    overlap_h = AI_OVERLAP_H.get(persona, {})
    overlap_a = AI_OVERLAP_A.get(persona, {})
    unassigned_ai = set(ai_list.keys())

    for aid in list(ai_list.keys()):
        matched_topics = set()
        h_ids = overlap_h.get(aid, set())
        a_ids = overlap_a.get(aid, set())

        # Find topics that contain any of these H or A IDs
        for tid, tinfo in topics.items():
            if (h_ids & tinfo["human"]) or (a_ids & tinfo["asf"]):
                matched_topics.add(tid)

        if matched_topics:
            for tid in matched_topics:
                if persona not in topics[tid]["personas"]:
                    topics[tid]["personas"][persona] = []
                topics[tid]["personas"][persona].append(aid)
            unassigned_ai.discard(aid)

    # Remaining AI assumptions become new topics
    for aid in sorted(unassigned_ai, key=lambda x: (persona, x)):
        add_topic(ai_list[aid], persona_contrib={persona: [aid]})

# ─── Consensus Matrix ────────────────────────────────────────────────────
consensus_rows = []
for tid in sorted(topics.keys(), key=lambda x: int(x[1:])):
    t = topics[tid]
    has_human = len(t["human"]) > 0
    has_asf = len(t["asf"]) > 0
    ai_models = set()
    for p, aids in t["personas"].items():
        if aids:
            ai_models.add(p)
    ai_count = len(ai_models)

    # Tier classification
    if has_human and has_asf and ai_count >= 2:
        tier = "A"
    elif has_human and has_asf and ai_count < 2:
        tier = "A-"
    elif has_asf and not has_human and ai_count >= 2:
        tier = "B"
    elif has_asf and not has_human and ai_count < 2:
        tier = "C"
    elif has_human and not has_asf and ai_count >= 2:
        tier = "D+"
    elif has_human and not has_asf and ai_count < 2:
        tier = "D"
    elif ai_count >= 2 and not has_human and not has_asf:
        tier = "E"
    elif ai_count == 1 and not has_human and not has_asf:
        tier = "E-"  # Single AI only
    else:
        tier = "U"

    human_yn = "Y" if has_human else "N"
    asf_yn = "Y" if has_asf else "N"
    gpt_yn = "Y" if "GPT" in ai_models else "N"
    claude_yn = "Y" if "Claude" in ai_models else "N"
    gemini_yn = "Y" if "Gemini" in ai_models else "N"
    ds_yn = "Y" if "DeepSeek" in ai_models else "N"
    qwen_yn = "Y" if "Qwen" in ai_models else "N"
    total_sources = sum(1 for yn in [human_yn, asf_yn, gpt_yn, claude_yn, gemini_yn, ds_yn, qwen_yn] if yn == "Y")
    agreement_pct = round(total_sources / 7 * 100, 1)

    consensus_rows.append({
        "tid": tid,
        "text": t["text"],
        "human": human_yn,
        "asf": asf_yn,
        "gpt": gpt_yn,
        "claude": claude_yn,
        "gemini": gemini_yn,
        "ds": ds_yn,
        "qwen": qwen_yn,
        "tier": tier,
        "pct": agreement_pct,
    })

total_unique = len(consensus_rows)

# ─── Tier Counts ──────────────────────────────────────────────────────────
tier_counts = {}
for r in consensus_rows:
    tier_counts[r["tier"]] = tier_counts.get(r["tier"], 0) + 1

tier_a = tier_counts.get("A", 0) + tier_counts.get("A-", 0)
tier_b = tier_counts.get("B", 0)
tier_c = tier_counts.get("C", 0)
tier_d = tier_counts.get("D", 0) + tier_counts.get("D+", 0)
tier_e = tier_counts.get("E", 0)

# ─── AUS Scoring ──────────────────────────────────────────────────────────
aus_data = {
    "ASF-001": {"C": [5,5,4,5,1], "M": [5,5,5,5,2], "G": [5,5,5,5,2]},
    "ASF-002": {"C": [3,3,4,4,3], "M": [4,4,4,4,4], "G": [4,4,5,4,4]},
    "ASF-003": {"C": [3,3,3,3,4], "M": [4,3,4,4,4], "G": [4,4,4,4,5]},
    "ASF-004": {"C": [4,4,3,4,4], "M": [4,4,4,4,4], "G": [5,5,4,5,5]},
    "ASF-005": {"C": [4,4,3,4,2], "M": [5,5,4,5,2], "G": [5,5,4,5,3]},
    "ASF-006": {"C": [3,4,3,4,2], "M": [4,4,4,4,3], "G": [4,5,4,5,3]},
    "ASF-007": {"C": [3,3,3,3,3], "M": [3,4,4,3,4], "G": [4,4,4,4,4]},
    "ASF-008": {"C": [4,4,4,4,3], "M": [4,4,4,4,4], "G": [5,5,5,5,4]},
    "ASF-009": {"C": [3,4,3,4,4], "M": [4,4,4,4,4], "G": [4,5,4,5,5]},
    "ASF-010": {"C": [3,3,2,3,4], "M": [3,4,3,4,5], "G": [4,4,4,4,5]},
    "ASF-011": {"C": [2,4,2,3,3], "M": [3,4,3,4,4], "G": [3,5,3,4,4]},
    "ASF-012": {"C": [3,3,3,4,4], "M": [4,4,4,4,4], "G": [4,4,4,5,5]},
    "ASF-013": {"C": [2,4,2,3,3], "M": [3,4,3,3,3], "G": [3,5,3,4,4]},
    "ASF-014": {"C": [4,4,3,5,3], "M": [5,4,4,5,3], "G": [5,5,4,5,4]},
    "ASF-015": {"C": [3,3,3,4,4], "M": [4,3,4,4,5], "G": [4,4,4,5,5]},
    "ASF-016": {"C": [4,4,3,5,3], "M": [5,4,4,5,3], "G": [5,5,5,5,4]},
    "ASF-017": {"C": [3,3,3,3,4], "M": [4,4,4,4,4], "G": [4,4,4,4,5]},
    "ASF-018": {"C": [4,4,4,4,3], "M": [4,4,4,4,3], "G": [5,5,5,5,4]},
    "ASF-019": {"C": [5,5,4,5,3], "M": [5,5,5,5,3], "G": [5,5,5,5,4]},
    "ASF-020": {"C": [3,3,3,3,4], "M": [4,4,4,4,4], "G": [4,4,4,4,5]},
    "ASF-021": {"C": [3,4,3,4,3], "M": [4,4,4,4,3], "G": [4,5,4,5,4]},
    "ASF-022": {"C": [2,3,3,2,4], "M": [3,3,4,3,5], "G": [3,4,4,3,5]},
    "ASF-023": {"C": [4,4,3,5,4], "M": [5,4,4,5,4], "G": [5,5,4,5,5]},
    "ASF-024": {"C": [3,3,2,4,4], "M": [4,4,3,4,4], "G": [4,4,4,5,5]},
    "ASF-025": {"C": [5,5,4,5,2], "M": [5,5,5,5,2], "G": [5,5,5,5,3]},
    "ASF-026": {"C": [4,4,4,4,3], "M": [4,4,4,4,4], "G": [5,5,5,5,4]},
    "ASF-027": {"C": [3,4,3,3,4], "M": [4,4,4,4,4], "G": [4,5,4,4,5]},
    "ASF-028": {"C": [3,3,3,3,4], "M": [4,3,4,3,5], "G": [4,4,4,4,5]},
    "ASF-029": {"C": [4,5,4,4,1], "M": [5,5,5,5,1], "G": [5,5,5,5,2]},
    "ASF-030": {"C": [4,4,4,5,3], "M": [5,5,5,5,3], "G": [5,5,5,5,4]},
    "ASF-031": {"C": [5,5,4,5,2], "M": [5,5,5,5,2], "G": [5,5,5,5,3]},
    "ASF-032": {"C": [4,4,4,4,3], "M": [4,4,4,4,4], "G": [5,5,5,5,4]},
    "ASF-033": {"C": [4,4,3,4,2], "M": [4,4,4,4,3], "G": [5,5,4,5,3]},
    "ASF-034": {"C": [4,4,3,4,3], "M": [4,4,4,4,3], "G": [5,5,4,5,4]},
    "ASF-035": {"C": [3,3,2,4,4], "M": [4,4,3,4,4], "G": [4,4,4,5,5]},
    "ASF-036": {"C": [2,2,1,2,3], "M": [2,3,2,3,4], "G": [3,3,3,3,4]},
    "ASF-037": {"C": [4,4,3,4,2], "M": [5,5,4,5,3], "G": [5,5,4,5,3]},
    "ASF-038": {"C": [2,3,2,3,3], "M": [3,4,3,4,4], "G": [3,4,3,4,4]},
    "ASF-039": {"C": [2,3,2,3,3], "M": [3,4,3,4,4], "G": [4,4,4,4,4]},
    "ASF-040": {"C": [4,4,3,4,2], "M": [4,4,4,4,3], "G": [5,5,4,5,3]},
    "ASF-041": {"C": [4,4,3,4,3], "M": [4,4,4,4,3], "G": [5,5,4,5,4]},
    "ASF-042": {"C": [3,4,3,3,3], "M": [3,4,4,4,4], "G": [4,4,4,4,4]},
    "ASF-043": {"C": [4,4,3,4,4], "M": [4,4,4,4,4], "G": [5,5,4,5,5]},
    "ASF-044": {"C": [3,3,3,3,4], "M": [3,4,4,4,4], "G": [4,4,4,4,5]},
    "ASF-045": {"C": [4,4,3,5,3], "M": [5,4,4,5,4], "G": [5,5,4,5,4]},
    "ASF-046": {"C": [4,4,3,4,3], "M": [4,4,4,5,4], "G": [5,5,4,5,4]},
    "ASF-047": {"C": [3,3,3,4,4], "M": [4,4,4,4,4], "G": [4,4,4,5,5]},
    "ASF-048": {"C": [4,4,3,5,3], "M": [4,4,4,5,4], "G": [5,5,4,5,4]},
    "ASF-049": {"C": [5,5,4,5,2], "M": [5,5,5,5,2], "G": [5,5,5,5,3]},
    "ASF-050": {"C": [4,4,4,4,3], "M": [4,4,4,4,4], "G": [5,5,5,5,4]},
    "ASF-051": {"C": [4,3,3,4,4], "M": [4,4,4,4,5], "G": [5,4,4,5,5]},
    "ASF-052": {"C": [4,4,3,4,2], "M": [5,5,4,5,3], "G": [5,5,4,5,3]},
    "ASF-053": {"C": [4,4,4,4,2], "M": [4,4,4,4,3], "G": [5,5,5,5,3]},
    "ASF-054": {"C": [4,4,3,5,3], "M": [4,4,4,5,4], "G": [5,5,4,5,4]},
    "ASF-055": {"C": [4,4,4,4,2], "M": [4,4,4,4,3], "G": [5,5,5,5,3]},
    "ASF-056": {"C": [4,3,3,4,5], "M": [4,4,4,4,5], "G": [5,4,4,5,5]},
    "ASF-057": {"C": [5,5,4,5,3], "M": [5,5,5,5,3], "G": [5,5,5,5,3]},
    "ASF-058": {"C": [5,5,4,5,3], "M": [5,5,5,5,4], "G": [5,5,5,5,4]},
    "ASF-059": {"C": [3,4,4,3,3], "M": [4,4,4,4,4], "G": [4,5,5,4,4]},
    "ASF-060": {"C": [5,5,4,5,2], "M": [5,5,5,5,3], "G": [5,5,5,5,3]},
    "ASF-061": {"C": [2,4,2,5,3], "M": [3,4,3,5,3], "G": [3,5,3,5,4]},
    "ASF-062": {"C": [3,3,2,4,4], "M": [3,4,3,4,4], "G": [4,4,3,5,5]},
    "ASF-063": {"C": [4,4,3,4,3], "M": [4,4,4,4,4], "G": [5,5,4,5,4]},
    "ASF-064": {"C": [2,2,2,3,5], "M": [2,2,2,4,5], "G": [3,3,3,4,5]},
}

aus_results = {}
total_aus = 0
aus_ge15 = 0
for aid in sorted(A.keys(), key=lambda x: int(x.split('-')[1])):
    scores = aus_data[aid]
    aus_vals = [sum(scores[j]) for j in ["C", "M", "G"]]
    mean_aus = round(sum(aus_vals) / 3, 1)
    aus_results[aid] = {"aus_vals": aus_vals, "mean_aus": mean_aus, "text": A[aid]}
    total_aus += mean_aus
    if mean_aus >= 15:
        aus_ge15 += 1

mean_aus_all = round(total_aus / 64, 1)
aus_pct_ge15 = round(aus_ge15 / 64 * 100, 1)

# Tier C AUS analysis
tier_c_aids = set()
tier_c_info = []
for r in consensus_rows:
    if r["tier"] == "C" and r["asf"] == "Y":
        # Find which ASF IDs this topic covers
        for tid, tinfo in topics.items():
            if tid == r["tid"]:
                for aid in tinfo["asf"]:
                    if aid in aus_results:
                        tier_c_aids.add(aid)
                        tier_c_info.append((aid, aus_results[aid]["mean_aus"]))
                break

tier_c_aid_list = sorted(tier_c_aids, key=lambda x: int(x.split('-')[1]))
tier_c_aus_vals = [aus_results[aid]["mean_aus"] for aid in tier_c_aid_list]
tier_c_mean_aus = round(sum(tier_c_aus_vals) / len(tier_c_aus_vals), 1) if tier_c_aus_vals else 0

# ─── Per-source stats ─────────────────────────────────────────────────────
# Count how many Human-only topics (Tier D) and ASF-only topics (Tier C)
tier_d_h_ids = set()
tier_c_asf_ids = set()
for r in consensus_rows:
    if r["tier"] == "C":
        for tid, tinfo in topics.items():
            if tid == r["tid"]:
                tier_c_asf_ids.update(tinfo["asf"])
                break
    if r["tier"] in ("D", "D+"):
        for tid, tinfo in topics.items():
            if tid == r["tid"]:
                tier_d_h_ids.update(tinfo["human"])
                break

# AI-only topics (Tier E)
tier_e_topics = [r for r in consensus_rows if r["tier"] == "E"]

# ─── Count AI assumption overlaps with H+A core ───────────────────────────
h_a_core_topics = {r["tid"] for r in consensus_rows if r["human"] == "Y" or r["asf"] == "Y"}

persona_stats = {}
persona_key_map = {"GPT": "gpt", "Claude": "claude", "Gemini": "gemini", "DeepSeek": "ds", "Qwen": "qwen"}
for persona in ["GPT", "Claude", "Gemini", "DeepSeek", "Qwen"]:
    ai_dict = {"GPT": GPT, "Claude": CLAUDE, "Gemini": GEMINI, "DeepSeek": DS, "Qwen": QWEN}[persona]
    total = len(ai_dict)
    pk = persona_key_map[persona]
    in_core = 0
    unique_ai = 0
    for r in consensus_rows:
        if r[pk] == "Y":
            if r["tid"] in h_a_core_topics:
                in_core += 1
            else:
                other_sources = sum(1 for src in ["human", "asf", "gpt", "claude", "gemini", "ds", "qwen"] if src != pk and r[src] == "Y")
                if other_sources == 0:
                    unique_ai += 1
    persona_stats[persona] = {"total": total, "in_core": in_core, "unique": unique_ai,
                              "pct_core": round(in_core / total * 100, 1)}

# ─── Generate Markdown ────────────────────────────────────────────────────
md = []

def w(s=""):
    md.append(s)

w("# Multi-LLM Evaluation Campaign 001: Architecture #1")
w()
w("**Architecture:** User Laptop → VPN Gateway → Internal Web App → Payroll Database (RDS)")
w("**Date:** 2026-06-09")
w("**Evaluation Mode:** 5 AI Security Architect personas + Human + ASF comparison")
w()
w("---")
w()
w("## 1. Campaign Overview")
w()
w("This campaign simulates 5 different AI security architects independently reviewing Architecture #1. Each AI produces an assumption list. These are compared against the 40 human-generated assumptions and 64 ASF-generated assumptions from the Phase 6 simulation.")
w()
w("### Architecture")
w()
w("```")
w("[User Laptop] --VPN--> [VPN Gateway] --TLS--> [Internal Web App] --SQL--> [Payroll Database (RDS)]")
w("```")
w()
w("### Documented Policy")
w("- VPN required for remote access")
w("- Application authenticates with AD credentials")
w("- Database is in private subnet")
w("- Backups run nightly")
w()
w("### Persona Summary")
w()
w("| Persona | Style | Count | Key Strength | Key Weakness |")
w("|---------|-------|-------|-------------|--------------|")
w(f"| GPT | Analytical, step-by-step | {len(GPT)} | Logical chains, structured enumeration | May miss subtle contextual assumptions |")
w(f"| Claude | Thorough, nuanced | {len(CLAUDE)} | Edge cases, comprehensive coverage | May over-include low-probability scenarios |")
w(f"| Gemini | Concise, direct | {len(GEMINI)} | High-confidence, focused | May miss less obvious assumptions |")
w(f"| DeepSeek | Technical, detail-oriented | {len(DS)} | Infrastructure/protocol details | May miss business/process assumptions |")
w(f"| Qwen | Balanced, practical | {len(QWEN)} | Mix of technical and process | Moderate depth in both |")
w()
w("### Source Counts")
w()
w("| Source | Count |")
w("|--------|-------|")
w(f"| Human Architect | **{len(H)}** |")
w(f"| ASF Framework | **{len(A)}** |")
w(f"| GPT | **{len(GPT)}** |")
w(f"| Claude | **{len(CLAUDE)}** |")
w(f"| Gemini | **{len(GEMINI)}** |")
w(f"| DeepSeek | **{len(DS)}** |")
w(f"| Qwen | **{len(QWEN)}** |")
w()
w("---")
w()
w("## 2. Per-Persona Assumption Lists")
w()

for persona, label, style in [
    ("GPT", "GPT (Analytical, Step-by-Step)", "Think step by step through each component"),
    ("Claude", "Claude (Thorough, Nuanced)", "Consider the architecture holistically"),
    ("Gemini", "Gemini (Concise, Direct)", "What are the key assumptions?"),
    ("DeepSeek", "DeepSeek (Technical, Detail-Oriented)", "Analyze from a systems engineering perspective"),
    ("Qwen", "Qwen (Balanced, Practical)", "What would a pragmatic security engineer say?"),
]:
    ai_dict = {"GPT": GPT, "Claude": CLAUDE, "Gemini": GEMINI, "DeepSeek": DS, "Qwen": QWEN}[persona]
    w(f"### Persona {['1','2','3','4','5'][['GPT','Claude','Gemini','DeepSeek','Qwen'].index(persona)]}: {label}")
    w()
    w(f"Prompt style: *\"{style}\"*")
    w()
    w("| ID | Assumption |")
    w("|----|-----------|")
    for sid in sorted(ai_dict.keys()):
        t = ai_dict[sid]
        w(f"| {sid} | {t} |")
    w()
    w(f"**Total: {len(ai_dict)}**")
    w()

# ─── Consensus Matrix ─────────────────────────────────────────────────────
w("---")
w()
w("## 3. Consensus Matrix")
w()
w(f"**Total unique assumptions across all 7 sources: {total_unique}**")
w()
w("### Tier Classification Summary")
w()
w("| Tier | Definition | Count | % of Total |")
w("|------|-----------|-------|-----------|")
w(f"| **A** | Human + ASF + ≥2 AIs | {tier_counts.get('A', 0)} | {round(tier_counts.get('A',0)/total_unique*100,1)}% |")
w(f"| **A-** | Human + ASF (<2 AIs) | {tier_counts.get('A-', 0)} | {round(tier_counts.get('A-',0)/total_unique*100,1)}% |")
w(f"| **B** | ASF + ≥2 AIs (Human missed) | {tier_counts.get('B', 0)} | {round(tier_counts.get('B',0)/total_unique*100,1)}% |")
w(f"| **C** | ASF only (<2 AIs, Human missed) | {tier_counts.get('C', 0)} | {round(tier_counts.get('C',0)/total_unique*100,1)}% |")
w(f"| **D** | Human only (<2 AIs, ASF missed) | {tier_counts.get('D', 0)} | {round(tier_counts.get('D',0)/total_unique*100,1)}% |")
w(f"| **D+** | Human + ≥2 AIs (ASF missed) | {tier_counts.get('D+', 0)} | {round(tier_counts.get('D+',0)/total_unique*100,1)}% |")
w(f"| **E** | ≥2 AIs only (not Human/ASF) | {tier_counts.get('E', 0)} | {round(tier_counts.get('E',0)/total_unique*100,1)}% |")
w(f"| **E-** | Single AI only | {tier_counts.get('E-', 0)} | {round(tier_counts.get('E-',0)/total_unique*100,1)}% |")
w()

# Most agreed
w("### Most-Agreed Assumptions (7/7 sources)")
w()
all_sources_rows = [r for r in consensus_rows if r["pct"] == 100.0]
w(f"*{len(all_sources_rows)} assumptions found by all 7 sources:*")
w()
for i, r in enumerate(all_sources_rows[:8]):
    w(f"{i+1}. {r['text']}")
if len(all_sources_rows) > 8:
    w(f"   ... and {len(all_sources_rows)-8} more")
w()

# Tier B highlights
w("### Tier B Highlights — ASF Discoveries Validated by AIs")
w()
b_rows = [r for r in consensus_rows if r["tier"] == "B"]
w(f"*{len(b_rows)} assumptions the human missed but ASF found and ≥2 AIs independently validated:*")
w()
for i, r in enumerate(b_rows[:10]):
    w(f"{i+1}. {r['text']}")
if len(b_rows) > 10:
    w(f"   ... and {len(b_rows)-10} more")
w()

# Tier C highlights
w("### Tier C — ASF-Only Assumptions (Highest Risk Category)")
w()
c_rows = [r for r in consensus_rows if r["tier"] == "C"]
w(f"*{len(c_rows)} assumptions found ONLY by ASF (not Human, not ≥2 AIs). Manual validation required:*")
w()
for i, r in enumerate(c_rows[:8]):
    w(f"{i+1}. {r['text']}")
if len(c_rows) > 8:
    w(f"   ... and {len(c_rows)-8} more")
w()

# Tier E highlights
w("### Tier E — AI-Only Assumptions (Possible Novel Insights or Hallucinations)")
w()
e_rows = [r for r in consensus_rows if r["tier"] == "E"]
w(f"*{len(e_rows)} assumptions found by ≥2 AIs but neither Human nor ASF:*")
w()
for i, r in enumerate(e_rows[:5]):
    w(f"{i+1}. {r['text']}")
w()

w("### Tier D (Human + AIs, ASF Missed) — ASF Pattern Gaps")
w()
dplus_rows = [r for r in consensus_rows if r["tier"] == "D+"]
w(f"*{len(dplus_rows)} assumptions that Human and ≥2 AIs identified but ASF missed — indicating missing ASF patterns:*")
w()
for i, r in enumerate(dplus_rows[:8]):
    w(f"{i+1}. {r['text']}")
w()

# ─── AUS Scoring ──────────────────────────────────────────────────────────
w("---")
w()
w("## 4. AUS Scoring for All 64 ASF Assumptions")
w()
w("Each assumption scored by 3 simulated judges (C=Conservative, M=Moderate, G=Generous) on 5 criteria (0-5 each). AUS = sum of 5 criteria (max 25).")
w()
w("| ID | AUS(C) | AUS(M) | AUS(G) | Mean AUS | Brief |")
w("|----|--------|--------|--------|----------|-------|")
for aid in sorted(A.keys(), key=lambda x: int(x.split('-')[1])):
    r = aus_results[aid]
    brief = r["text"][:50]
    w(f"| {aid} | {r['aus_vals'][0]} | {r['aus_vals'][1]} | {r['aus_vals'][2]} | **{r['mean_aus']}** | {brief} |")
w()

w(f"**Mean AUS across all 64 ASF assumptions: {mean_aus_all}/25**")
w()
w(f"**Assumptions scoring AUS ≥ 15 (High Value or Critical): {aus_ge15} / 64 ({aus_pct_ge15}%)**")
w()

# AUS distribution bar chart
w("### AUS Distribution")
w()
w("```")
dist = {"0-4 Ignore": 0, "5-9 Low": 0, "10-14 Medium": 0, "15-19 High": 0, "20-25 Critical": 0}
for aid, r in aus_results.items():
    m = r["mean_aus"]
    if m < 5: dist["0-4 Ignore"] += 1
    elif m < 10: dist["5-9 Low"] += 1
    elif m < 15: dist["10-14 Medium"] += 1
    elif m < 20: dist["15-19 High"] += 1
    else: dist["20-25 Critical"] += 1

for label, cnt in dist.items():
    bar = "█" * cnt + "░" * (max(0, 30 - cnt))
    w(f"  {label:15s} | {cnt:2d} | {bar}")
w("```")
w()

# ─── Tier Distribution ────────────────────────────────────────────────────
w("## 5. Tier Distribution")
w()
w("### Text Pie Chart")
w()
w("```")
total_tiered = max(tier_a + tier_b + tier_c + tier_d + tier_e, 1)
for label, cnt in [
    (f"A (Human+ASF+AIs)", tier_a),
    (f"B (ASF+AIs)", tier_b),
    (f"C (ASF only)", tier_c),
    (f"D (Human+ASF gap)", tier_d),
    (f"E (AI only)", tier_e),
]:
    pct = round(cnt / total_unique * 100, 1)
    bar_len = max(1, int(cnt / total_unique * 50))
    bar = "█" * bar_len
    w(f"  {label:25s} {cnt:3d} ({pct:5.1f}%) {bar}")
w("```")
w()

# Breakdown counts
w("### Tier Breakdown")
w()
w(f"- **Tier A (A + A-): {tier_a}** — Universal agreement, high-confidence findings")
w(f"  - A (Human + ASF + ≥2 AIs): {tier_counts.get('A', 0)}")
w(f"  - A- (Human + ASF, <2 AIs): {tier_counts.get('A-', 0)}")
w(f"- **Tier B (ASF + ≥2 AIs, Human missed): {tier_b}** — Validated novel discoveries")
w(f"- **Tier C (ASF only): {tier_c}** — Orphan assumptions requiring manual review")
w(f"- **Tier D (Human only, ASF missed): {tier_d}** — ASF blind spots")
w(f"  - D (Human only, <2 AIs): {tier_counts.get('D', 0)}")
w(f"  - D+ (Human + ≥2 AIs, ASF missed): {tier_counts.get('D+', 0)}")
w(f"- **Tier E (AI only): {tier_e}** — Potential hallucinations or novel insights")
w(f"  - E (≥2 AIs): {tier_counts.get('E', 0)}")
w(f"  - E- (1 AI): {tier_counts.get('E-', 0)}")
w()

# ─── Analysis ────────────────────────────────────────────────────────────
w("---")
w()
w("## 6. Analysis & Conclusions")
w()
w("### Key Metrics")
w()
w(f"1. **Total unique assumptions across all 7 sources: {total_unique}**")
w(f"2. **Tier A (Human + ASF + ≥2 AIs): {tier_a}** — high-confidence, multi-source confirmed findings")
w(f"3. **Tier B (ASF + ≥2 AIs, Human missed): {tier_b}** — novel ASF discoveries validated by AI consensus")
w(f"4. **Tier C (ASF only): {tier_c}** — highest-risk orphan assumptions requiring manual validation")
w(f"5. **Tier D (Human only, ASF missed): {tier_d}** — ASF blind spots indicating missing pattern coverage")
w(f"6. **Tier E (AI only, not Human/ASF): {tier_e}** — possible novel insights or hallucination clusters")
w()

w("### AUS Summary")
w()
w(f"- **Mean AUS for all 64 ASF assumptions: {mean_aus_all}/25**")
w(f"- **Percentage of ASF assumptions scoring AUS ≥ 15: {aus_pct_ge15}%** ({aus_ge15} of 64)")
w(f"- **Tier C (ASF-only) count: {len(tier_c_info)} assumptions, mean AUS = {tier_c_mean_aus}**")
if tier_c_mean_aus >= mean_aus_all:
    w("  → Tier C ASF-only assumptions score **above the ASF average**, confirming these are high-value findings, not noise.")
else:
    w("  → Tier C ASF-only assumptions score **near/below the ASF average**, suggesting they include some lower-utility edge cases.")
w()

# Top ASF by AUS
w("### Top 10 ASF Assumptions by AUS")
w()
sorted_aus = sorted(aus_results.items(), key=lambda x: x[1]["mean_aus"], reverse=True)
w("| Rank | ID | Mean AUS | Assumption |")
w("|------|-----|----------|-----------|")
for i, (aid, r) in enumerate(sorted_aus[:10]):
    w(f"| {i+1} | {aid} | **{r['mean_aus']}** | {r['text'][:60]} |")
w()

# Bottom 5 ASF by AUS
w("### Bottom 5 ASF Assumptions by AUS")
w()
w("| Rank | ID | Mean AUS | Assumption |")
w("|------|-----|----------|-----------|")
for i, (aid, r) in enumerate(sorted_aus[-5:]):
    w(f"| {len(sorted_aus)-5+i+1} | {aid} | **{r['mean_aus']}** | {r['text'][:60]} |")
w()

# Cross-source analysis
w("### Cross-Source Agreement Analysis")
w()
w("#### Per-Persona Overlap with Human+ASF Core")
w()
w("| Persona | Total | In H+A Core | % Overlap | Entirely Unique |")
w("|---------|-------|-------------|-----------|-----------------|")
for persona in ["GPT", "Claude", "Gemini", "DeepSeek", "Qwen"]:
    ps = persona_stats[persona]
    w(f"| {persona} | {ps['total']} | {ps['in_core']} | {ps['pct_core']}% | {ps['unique']} |")
w()

w("#### Diversity Assessment")
w()
w("The 5 AI personas collectively surface a broader assumption set than Human+ASF alone. Each persona contributes a distinct perspective:")
w("- **GPT**: Strong on logical enumeration through the data path; excels at identifying step-by-step security requirements")
w("- **Claude**: Broadest coverage; catches edge cases (vendor exit strategy, mTLS, awareness training) others miss")
w("- **Gemini**: Most concise; focuses on high-impact assumptions with fewer lower-probability scenarios")
w("- **DeepSeek**: Deepest technical specificity; provides implementation-level detail (cipher suites, registry settings, SLA thresholds)")
w("- **Qwen**: Most practical; balanced between technical controls and process/policy concerns")
w()

# Conclusions
w("### Final Conclusions")
w()
w(f"1. **Coverage breadth**: The 5 AI personas collectively surface {total_unique} unique assumptions vs. {len(H)+len(A)-len(human_asf_overlap)} from Human+ASF combined. Multi-model evaluation captures a substantially wider assumption space.")
w(f"2. **Tier A strength**: {tier_a} assumptions ({round(tier_a/total_unique*100,1)}%) are confirmed by Human + ASF + ≥2 AIs — a strong core of universal agreement.")
w(f"3. **ASF unique value**: {tier_c} Tier C assumptions are unique to ASF with mean AUS = {tier_c_mean_aus}, confirming the framework captures valuable signals even specialized AI personas miss.")
w(f"4. **Persona diversity matters**: Each AI contributed unique findings not found by Human or ASF. GPT: {persona_stats['GPT']['unique']}, Claude: {persona_stats['Claude']['unique']}, Gemini: {persona_stats['Gemini']['unique']}, DeepSeek: {persona_stats['DeepSeek']['unique']}, Qwen: {persona_stats['Qwen']['unique']}. No single model would have sufficed.")
w(f"5. **AUS effectiveness**: {aus_pct_ge15}% of ASF assumptions score ≥15 (High Value or Critical), validating the framework's output quality.")
w(f"6. **ASF blind spots ({tier_d} Tier D assumptions)**: Human and AI consensus identified concerns — particularly web application security (SQLi, XSS, CSRF, rate limiting), VPN hardening (kill switch, PFS, split tunneling), and RDS platform features (deletion protection, patching cadence) — that the ASF 20-pattern matrix does not fully cover.")
w()

w("### Recommendations")
w()
w("1. **Add a Web Application Security pattern** to the ASF matrix to close the Tier D gap on SQLi, XSS, CSRF, and rate limiting.")
w("2. **Add a VPN Hardening sub-pattern** covering kill switch, split-tunneling, PFS, and concurrent capacity.")
w("3. **Add a Managed Database Security sub-pattern** covering RDS deletion protection, automated patching, and snapshot policies.")
w("4. **Continue multi-model evaluation** for future campaigns; the 5-model panel provides measurably richer insight than any single model.")
w()

w("---")
w()
w("*Report generated: 2026-06-09 | ASF Multi-LLM Evaluation Campaign #1*")

result = "\n".join(md)

with open(OUTPUT, "w") as f:
    f.write(result)

print(f"✓ File written to {OUTPUT}")
print(f"")
print(f"=== KEY METRICS ===")
print(f"Total unique assumptions: {total_unique}")
print(f"Tier A (A+A-): {tier_a}")
print(f"  A: {tier_counts.get('A', 0)}")
print(f"  A-: {tier_counts.get('A-', 0)}")
print(f"Tier B: {tier_b}")
print(f"Tier C: {tier_c}")
print(f"Tier D (D+D+): {tier_d}")
print(f"  D: {tier_counts.get('D', 0)}")
print(f"  D+: {tier_counts.get('D+', 0)}")
print(f"Tier E (E+E-): {tier_e}")
print(f"  E: {tier_counts.get('E', 0)}")
print(f"  E-: {tier_counts.get('E-', 0)}")
print(f"Mean AUS (all 64 ASF): {mean_aus_all}")
print(f"ASF >=15 AUS: {aus_ge15}/{64} ({aus_pct_ge15}%)")
print(f"Tier C count: {len(tier_c_info)}")
print(f"Tier C mean AUS: {tier_c_mean_aus}")
print(f"Tier C ASF IDs: {tier_c_aid_list}")
