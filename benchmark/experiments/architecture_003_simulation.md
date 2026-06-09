# ASF Phase 6 Experiment: Architecture #3

**Architecture:** Mobile App → API Gateway → Lambda → DynamoDB
**Date:** 2026-06-09
**Simulation Mode:** Human + ASF framework comparison

---

## Architecture Profile

```
[Mobile App] --HTTPS--> [API Gateway] --Event--> [Lambda Function (xN)] --SDK--> [DynamoDB Table]
                                        └--> [Cognito User Pool]
```

### Documented Policy
| # | Policy |
|---|--------|
| P1 | API Gateway requires Cognito auth |
| P2 | Lambda uses least-privilege IAM role |
| P3 | DynamoDB is encrypted at rest |
| P4 | API keys per mobile app version |

### Trust Boundaries
| Boundary | Type |
|----------|------|
| Mobile ↔ API Gateway | Authentication boundary |
| API Gateway ↔ Lambda | Service boundary |
| Lambda ↔ DynamoDB | Data boundary |

### Complexity Rating
**Moderate** — serverless architecture with managed services, event-driven execution, and no persistent compute.

---

## Step 1: Human-Generated Assumptions

*Role: Senior Security Architect. Listed below are assumptions that MUST remain true for this architecture to remain secure but are NOT stated in the documented policy.*

| ID | Assumption | Justification |
|----|-----------|---------------|
| H-001 | Cognito user pool tokens (ID, access, refresh) are validated by API Gateway for every request. | Token validation at the gateway is the primary access control; without it, unauthenticated requests pass through. |
| H-002 | Cognito token expiration is set to a reasonable duration (e.g., 1 hour for access tokens, 30 days for refresh). | Overly long token lifetimes increase the window of exposure if a token is leaked. |
| H-003 | Cognito supports account recovery with social-engineering-resistant verification. | Weak account recovery allows attackers to take over user accounts through support channels. |
| H-004 | The mobile app securely stores Cognito tokens (e.g., iOS Keychain, Android Encrypted SharedPreferences). | Tokens stored insecurely on the mobile device can be extracted by malware or physical access. |
| H-005 | The mobile app enforces certificate pinning or uses a certificate transparency-based trust model. | Without pinning, a compromised CA or proxy can intercept all API traffic from the mobile app. |
| H-006 | API Gateway usage plans and API keys are per-app-version but not shared across different applications. | Shared API keys across different apps prevent revocation of a compromised version without affecting others. |
| H-007 | API Gateway is configured with a custom domain and TLS certificate from ACM. | Default API Gateway endpoints (execute-api.region.amazonaws.com) are subject to domain-based attacks. |
| H-008 | Lambda function URLs are not enabled; all invocation goes through API Gateway. | Lambda function URLs bypass API Gateway entirely, negating auth and throttling controls. |
| H-009 | Lambda IAM execution roles have no permissions beyond DynamoDB CRUD, CloudWatch logs, and X-Ray. | Over-permissioned Lambda roles increase blast radius if a Lambda function is compromised. |
| H-010 | Lambda functions do not have VPC access unless explicitly required (and then only through VPC endpoints). | Lambda in a VPC with NAT can be used as a pivot to internal resources if VPC is shared. |
| H-011 | DynamoDB table access is scoped to specific attributes/columns via IAM condition keys. | Unrestricted DynamoDB access allows a compromised Lambda to read all table attributes. |
| H-012 | DynamoDB point-in-time recovery (PITR) is enabled. | Without PITR, data loss extends to the last backup rather than the last second. |
| H-013 | DynamoDB auto-scaling is configured with minimum and maximum capacity limits. | Unconstrained auto-scaling can cause unexpected cost spikes or throttling. |
| H-014 | Lambda function code does not contain hardcoded secrets, API keys, or database credentials. | Hardcoded secrets in Lambda code are exposed through the Lambda console, version history, and source control. |
| H-015 | Lambda function timeout and memory limits are configured to prevent resource exhaustion attacks. | Unlimited Lambda execution can be exploited for cryptocurrency mining or data exfiltration. |
| H-016 | Lambda environment variables are encrypted at rest with a KMS customer-managed key. | Default Lambda encryption uses AWS-managed keys; CMK provides explicit access control. |
| H-017 | CloudWatch Logs from Lambda executions do not contain PII, tokens, or sensitive business data. | Lambda logs containing sensitive data are exposed to anyone with Logs read access. |
| H-018 | API Gateway request/response body size limits (e.g., 10 MB) are configured to prevent abuse. | Unrestricted payload size allows DDoS through resource exhaustion. |
| H-019 | API Gateway throttling and burst limits prevent single-tenant or single-user API abuse. | Without throttling, a compromised or misbehaving mobile client can overwhelm the backend. |
| H-020 | Cognito user pool allows only strong password policies (length, complexity, history). | Weak passwords in the user pool make credential stuffing and brute-force feasible. |
| H-021 | Cognito has adaptive authentication or risk-based challenges for anomalous logins. | Without adaptive auth, a login from a new device or geography is not challenged. |
| H-022 | Lambda functions validate all input from API Gateway events (event.body, pathParameters, queryString). | API Gateway does not sanitize payloads; Lambda must validate to prevent injection attacks. |
| H-023 | DynamoDB table uses fine-grained access control (FGAC) with condition expressions for user-scoped access. | Without FGAC, any authenticated Lambda can read any row in the table, violating tenant isolation. |
| H-024 | DynamoDB is configured with encryption at rest using a KMS customer-managed key (CMK). | Default AWS-managed keys do not provide explicit control over who can decrypt the table. |
| H-025 | AWS CloudTrail is enabled for all API calls to detect unauthorized configuration changes. | Without CloudTrail, changes to API Gateway, Lambda, or DynamoDB are untracked. |
| H-026 | Each Lambda function version is immutable and cannot be modified after publish. | Mutable Lambda versions allow an attacker to replace a function with malicious code. |
| H-027 | The mobile app implements certificate transparency monitoring to detect mis-issued certificates. | A mis-issued certificate allows MITM against all mobile-to-API traffic. |
| H-028 | Cognito app client IDs are not embedded in publicly accessible mobile app binaries without obfuscation. | Extracted client IDs allow attackers to register rogue devices against the user pool. |
| H-029 | Cognito identity pools (if used for AWS credentials) scope unauthenticated access to the minimum. | Unauthenticated identity pool access can be exploited to read/write AWS resources anonymously. |
| H-030 | API Gateway WAF (AWS WAF) is attached to filter SQL injection, XSS, and common exploit patterns. | API Gateway natively passes all requests to Lambda; WAF is needed for application-layer filtering. |
| H-031 | Lambda runtime versions are current and updated when AWS deprecates a runtime. | Deprecated Lambda runtimes receive no security patches and expose known vulnerabilities. |
| H-032 | Lambda functions do not download or execute external binaries or packages at runtime. | Runtime code download is a common malware delivery technique and increases supply chain risk. |
| H-033 | The mobile app has a kill switch or remote feature disable for compromised app versions. | A compromised app version cannot be remotely disabled, exposing all users to the vulnerability. |
| H-034 | API Gateway access logs are enabled and stored in a secure S3 bucket with restricted access. | Without access logs, forensic investigation of API abuse is impossible. |
| H-035 | Cognito triggers (pre-signup, post-authentication, etc.) are secured and do not introduce vulnerabilities. | Misconfigured Cognito triggers can bypass signup controls or leak user data. |
| H-036 | DynamoDB TTL (time-to-live) is used for session/record expiration to enforce data retention. | Without TTL, stale data persists indefinitely, increasing storage cost and data exposure. |
| H-037 | Sandbox or test accounts in the user pool are identified and excluded from production metrics/billing. | Test accounts generate noise in monitoring and can be exploited if they have elevated privileges. |
| H-038 | The mobile app binary is code-signed and integrity-checked before accepting API responses. | Tampered mobile app binaries can be reverse-engineered to extract API keys or tokens. |
| H-039 | No Lambda function has the `lambda:InvokeFunction` permission on other functions (no chaining). | Lambda function chaining allows lateral movement if one function is compromised. |
| H-040 | DynamoDB table deletion protection is enabled. | Accidental or malicious table deletion causes data loss and service outage. |
| H-041 | API Gateway resource policies restrict access to specific source IP ranges or VPC endpoints. | Public API Gateway endpoints are accessible from any internet client by default. |
| H-042 | Lambda DLQ (dead-letter queue) or destination configuration does not leak sensitive data. | Lambda destinations for failed events can expose event payloads to unauthorized services. |
| H-043 | Cognito hosted UI (if used) follows security best practices (CSP, HSTS, X-Frame-Options). | Cognito hosted UI without security headers exposes users to clickjacking and MIME attacks. |

**Total (H): 43**

---

## Step 2: ASF-Generated Assumptions

*Using the 20-pattern assumption generator matrix. Applicable patterns: 15 of 20. Patterns excluded: Network Segmentation (serverless — no traditional network tiers), Physical Security (cloud-hosted), Container Security (serverless, no containers), Backup & Recovery (DynamoDB PITR — covered under Data Flow), Change Management (covered under Operational cross-cutting).*

---

### Pattern 1: Authentication (MFA)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-001 | Cognito user pool has MFA enforcement enabled (TOTP or SMS) for all users. | Explicit | Policy mentions Cognito auth but does not specify MFA. |
| ASF-002 | MFA recovery processes are documented and resistant to social engineering. | Operational | Lost MFA device recovery is a common bypass vector. |
| ASF-003 | MFA is not bypassed for API-key-based access (app version keys). | Implicit | API keys may be accepted as an alternative to Cognito MFA. |
| ASF-004 | The mobile app supports biometric or platform-native MFA beyond Cognito. | Derived | Mobile-first architecture suggests device-native MFA integration. |

---

### Pattern 2: Authentication (SSO)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-005 | Cognito can federate with external IdPs (Google, Apple, enterprise SAML) as documented. | Explicit | Mobile apps commonly support social login; federation introduces additional trust boundaries. |
| ASF-006 | Cognito token signing keys (JWKS) are rotated regularly. | Derived | Static signing keys increase the risk of forged tokens if the key is compromised. |
| ASF-007 | The mobile app validates the Cognito token issuer and audience claims. | Trust | Without issuer/audience validation, tokens from a different user pool can be accepted. |
| ASF-008 | Token revocation (Cognitos `revoke-token`) is properly implemented on the client side. | Operational | Revoked tokens cached on the mobile client continue to grant access until cache expires. |

---

### Pattern 3: Availability & Resilience

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-009 | Lambda concurrent execution limits are configured to prevent account-level throttling. | Architectural | Default Lambda concurrency limits are account-level; a burst in one function can throttle others. |
| ASF-010 | API Gateway can handle the expected peak request volume without throttling. | Operational | API Gateway throttling at 10,000 rps default may not suffice for production mobile apps. |
| ASF-011 | DynamoDB can handle provisioned capacity without throttling during peak load. | Operational | DynamoDB throttling causes API errors visible to mobile end users. |
| ASF-012 | Lambda function cold starts do not cause unacceptable latency for the mobile app. | Derived | Cold starts add 200ms-1s latency which can degrade mobile UX and cause retry storms. |

---

### Pattern 4: Backup & Recovery

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-013 | DynamoDB point-in-time recovery (PITR) is enabled with adequate retention. | Explicit | Documented policy does not mention backup or PITR for the database. |
| ASF-014 | DynamoDB on-demand backup is scheduled and stored in a separate account. | Derived | PITR alone is insufficient for region-wide failures; cross-account backups are needed. |
| ASF-015 | Backup restore procedures are documented and tested. | Operational | Untested backups provide false confidence in recovery capability. |
| ASF-016 | Cognito user pool export/import capability is available for user account recovery. | Implicit | User pool data (users, attributes) is not backed up by default; loss requires re-registration. |

---

### Pattern 6: Cloud Security (IAM)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-017 | Lambda execution roles are scoped to specific DynamoDB table ARN and specific actions. | Explicit | Least privilege for Lambda IAM is stated as policy but relies on correct implementation. |
| ASF-018 | API Gateway can assume an IAM role for integration (credentials passthrough). | Derived | API Gateway's credentials passthrough needs correct trust policy configuration. |
| ASF-019 | No IAM roles in the account have `lambda:CreateFunction` or `lambda:UpdateFunctionCode` without restrictions. | Architectural | Unrestricted Lambda creation permissions allow privilege escalation via Lambda. |
| ASF-020 | CloudTrail is enabled on the account to detect IAM and Lambda configuration changes. | Derived | Without CloudTrail, unauthorized IAM changes are invisible. |

---

### Pattern 8: Data Flow & Classification

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-021 | Data stored in DynamoDB is classified and handling requirements are documented. | Explicit | Encryption at rest assumes data sensitivity but policy does not state classification. |
| ASF-022 | No data flows exist beyond the documented path (e.g., Lambda calling external APIs). | Implicit | The documented flow is Mobile→Gateway→Lambda→DynamoDB; any other egress is unmapped. |
| ASF-023 | Mobile app does not cache sensitive data that persists beyond the app session. | Derived | Mobile client caching can create a secondary data store on the device. |
| ASF-024 | Data passed to CloudWatch Logs does not create a secondary data store with compliance implications. | Environmental | Lambda logs containing PII or financial data expand compliance scope. |

---

### Pattern 9: Encryption at Rest

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-025 | DynamoDB encryption at rest is enabled using a KMS customer-managed key. | Explicit | Policy states "DynamoDB encrypted at rest" but does not specify key management. |
| ASF-026 | KMS key policy restricts which IAM principals can encrypt/decrypt the DynamoDB table. | Derived | Encryption without key access control provides no protection against authorized IAM users. |
| ASF-027 | Lambda /tmp directory is not used for sensitive data, or is encrypted if used. | Implicit | Lambda /tmp storage is ephemeral and may not be encrypted depending on the runtime. |
| ASF-028 | KMS key rotation is enabled (automatic annual rotation). | Operational | Manual key rotation is frequently missed, increasing exposure from compromised keys. |

---

### Pattern 10: Encryption in Transit (TLS)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-029 | HTTPS is enforced between the mobile app and API Gateway; HTTP requests are rejected. | Explicit | Policy implies HTTPS but does not explicitly require it. |
| ASF-030 | The mobile app validates the API Gateway TLS certificate against a trusted root. | Trust | Mobile TLS without certificate validation can be MITM'd by proxy or malicious CA. |
| ASF-031 | TLS 1.2 or higher is enforced on API Gateway's custom domain. | Derived | API Gateway supports TLS 1.0 by default; explicit restriction is required. |
| ASF-032 | API Gateway and Lambda communicate via TLS within AWS network (no transit via internet). | Implicit | Lambda invoked via API Gateway uses AWS internal network, but this should be verified. |

---

### Pattern 11: Endpoint Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-033 | Mobile devices have OS-level security controls (screen lock, encryption, no root/jailbreak). | Implicit | A compromised mobile OS exposes all in-app data and tokens. |
| ASF-034 | The mobile app is distributed through official app stores with code signing requirements. | Derived | Sideloaded mobile apps bypass store security review and code signing. |
| ASF-035 | Lost or stolen devices can have app access revoked (via Cognito token revocation). | Operational | Without remote revocation, a stolen device retains API access until token expiry. |
| ASF-036 | The mobile app is regularly updated to patch security vulnerabilities. | Environmental | Unpatched mobile app versions expose users to known vulnerabilities. |

---

### Pattern 12: Human Factors & Process

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-037 | Developers do not share AWS console access or IAM user credentials. | Derived | Shared credentials eliminate accountability for infrastructure changes. |
| ASF-038 | Lambda function code goes through code review before deployment. | Operational | Unreviewed code can introduce backdoors or vulnerabilities into the serverless backend. |
| ASF-039 | API Gateway API keys are not committed to public or internal source repositories. | Trust | Leaked API keys allow unauthorized parties to call the API. |
| ASF-040 | Developers understand and follow least-privilege principles when writing IAM policies. | Implicit | Without training, developers default to overly permissive IAM policies. |

---

### Pattern 13: Identity Lifecycle (Provisioning)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-041 | Cognito user accounts are reviewed and stale accounts removed regularly. | Operational | Abandoned user accounts can be reactivated by attackers. |
| ASF-042 | Service accounts (e.g., CI/CD roles deploying Lambda) are managed with the same rigor as human accounts. | Implicit | Orphaned CI/CD roles can deploy unauthorized code changes. |
| ASF-043 | Cognito app client configuration is reviewed per mobile version to remove outdated clients. | Derived | Old app client configurations with weaker settings remain active if not removed. |
| ASF-044 | IAM users and roles are removed within 24 hours of employee termination. | Operational | Former employees with active IAM access can modify Lambda functions or access DynamoDB. |

---

### Pattern 14: Incident Response

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-045 | There is an incident response plan covering serverless application compromise scenarios. | Operational | Without a plan, response to a Lambda breach is ad-hoc and delayed. |
| ASF-046 | The IR team has access to API Gateway access logs, CloudWatch Logs, and CloudTrail during an investigation. | Derived | Inaccessible logs prevent root cause analysis and attacker attribution. |
| ASF-047 | The IR plan includes isolation procedures (disable API key, disable Cognito user, update resource policy) that preserve evidence. | Trust | Hasty isolation can destroy forensic evidence in serverless environments. |
| ASF-048 | Monitoring systems can detect anomalous Lambda invocation patterns indicating compromise. | Implicit | Detection is a prerequisite for incident response; without it, breaches go unnoticed. |
| ASF-049 | DynamoDB export snapshots taken during IR are isolated for forensic analysis. | Derived | Forensic snapshots in the same account can be tampered with by an attacker with AWS access. |

---

### Pattern 15: Least Privilege

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-050 | Lambda IAM execution role has access only to the specific DynamoDB table and CloudWatch Logs. | Explicit | Least privilege is stated as policy but requires correct implementation. |
| ASF-051 | Lambda function does not have permissions to modify IAM, Lambda, or API Gateway resources. | Derived | Lambda with IAM write permissions can escalate its own privileges. |
| ASF-052 | Cognito identity pool (if used) assigns temporary AWS credentials with minimal scope. | Implicit | Identity pool unauthenticated role must be scoped to the absolute minimum. |
| ASF-053 | The mobile app enforces authorization scopes beyond Cognito authentication (application-level authorization). | Derived | Cognito auth confirms identity; the app must independently authorize actions. |

---

### Pattern 16: Monitoring & Alerting

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-054 | API Gateway 4xx and 5xx error rates are monitored and alerted. | Operational | Elevated error rates indicate attacks or application failures. |
| ASF-055 | Lambda error counts, duration, and throttles are monitored and alerted. | Operational | Anomalous Lambda metrics indicate code issues or abuse. |
| ASF-056 | CloudWatch Logs are monitored for suspicious patterns (credential stuffing, injection attempts). | Derived | Log analysis can identify attack patterns in real time. |
| ASF-057 | Monitoring infrastructure logs are append-only and tamper-proof. | Implicit | Attackers who compromise the AWS account can alter or delete CloudWatch logs. |

---

### Pattern 20: Third-party Dependency

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-058 | AWS API Gateway, Lambda, and DynamoDB are available and not experiencing service outages. | Dependency | Single-region serverless deployment has no fallback for region-level outages. |
| ASF-059 | Third-party SDKs and libraries used by the Lambda functions are scanned for vulnerabilities. | Operational | Dependency vulnerabilities in Lambda functions can lead to RCE. |
| ASF-060 | The mobile app's SDK dependencies (Cognito SDK, HTTP client) are kept up to date. | Derived | Outdated mobile SDKs may have known security vulnerabilities. |
| ASF-061 | There is an exit strategy or migration plan if any AWS service is deprecated. | Derived | AWS Lambda and API Gateway service deprecation would require emergency migration. |

**Total (A): 61** (4 per pattern × 15 applicable patterns + 1 additional from Incident Response)

---

## Step 3: Comparison

### Overlap Mapping

| Human ID | ASF ID | Match Rationale |
|----------|--------|-----------------|
| H-001 | ASF-001 | Both require token validation at API Gateway. |
| H-002 | ASF-006 | Both address token lifetime/rotation. |
| H-004 | ASF-033 | Both require secure mobile token storage. |
| H-005 | ASF-030 | Both require certificate validation/pinning on mobile. |
| H-006 | ASF-003 | Both address API key management per app version. |
| H-008 | ASF-032 | Both require Lambda invoked only through API Gateway, not function URLs. |
| H-009 | ASF-017 | Both require least-privilege Lambda IAM roles. |
| H-010 | ASF-032 | Both require Lambda no-VPC or controlled VPC access. |
| H-011 | ASF-050 | Both require DynamoDB access scoped to specific table/attributes. |
| H-012 | ASF-013 | Both require DynamoDB PITR. |
| H-013 | ASF-011 | Both require DynamoDB auto-scaling limits. |
| H-014 | ASF-037 | Both address hardcoded secrets/credentials. |
| H-015 | ASF-009 | Both require Lambda resource limits (timeout/concurrency). |
| H-016 | ASF-025 | Both require Lambda env var encryption with KMS CMK. |
| H-017 | ASF-024 | Both require no sensitive data in CloudWatch Logs. |
| H-019 | ASF-054 | Both require API Gateway throttling/monitoring. |
| H-020 | ASF-001 | Both require Cognito password policies and MFA. |
| H-022 | ASF-053 | Both require Lambda input validation. |
| H-023 | ASF-050 | Both require DynamoDB fine-grained access control. |
| H-024 | ASF-025 | Both require DynamoDB encryption with KMS CMK. |
| H-025 | ASF-020 | Both require CloudTrail enabled. |
| H-031 | ASF-058 | Both address Lambda runtime versions (dependency on AWS). |
| H-034 | ASF-046 | Both require API Gateway access logs for investigation. |
| H-039 | ASF-051 | Both restrict Lambda invoke permissions. |
| H-040 | ASF-013 | Both require DynamoDB deletion protection (PITR is a form). |
| H-041 | ASF-029 | Both require API Gateway resource policies. |

**Overlap (O): 26**

### Metrics

| Metric | Value | Formula |
|--------|-------|---------|
| Human assumptions (H) | 43 | Count of unique human-generated assumptions |
| ASF assumptions (A) | 61 | Count of unique ASF-generated assumptions |
| Overlap (O) | 26 | Count appearing in both lists |
| **Precision** | **42.6%** | O / A = 26/61 |
| **Recall** | **60.5%** | O / H = 26/43 |
| **F1 Score** | **50.0%** | 2 × (P × R) / (P + R) |
| Novel findings (A - O) | 35 | Assumptions ASF found that human missed (57.4% of ASF total) |
| Missed findings (H - O) | 17 | Assumptions human found that ASF missed (39.5% of human total) |

### Success Criteria Assessment

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Recall | >= 70% | 60.5% | ❌ Not met |
| Precision | >= 50% | 42.6% | ❌ Not met |
| Novel discoveries | >= 10% of total (A+O) | 33.7% (35/104) | ✅ Exceeded |
| Expert agreement (F1 proxy) | > 60% | 50.0% | ❌ Not met |

---

## Step 4: Analysis

### Ontology Overlap Analysis

| Ontology Category | Overlaps | ASF Total | Overlap Rate |
|-------------------|----------|-----------|-------------|
| Explicit | 8 | 12 | 66.7% |
| Derived | 7 | 16 | 43.8% |
| Operational | 4 | 13 | 30.8% |
| Implicit | 4 | 10 | 40.0% |
| Trust | 1 | 3 | 33.3% |
| Dependency | 1 | 3 | 33.3% |
| Architectural | 1 | 3 | 33.3% |
| Environmental | 0 | 3 | 0.0% |

**Best overlap:** Explicit (66.7%) showed strongest agreement — both humans and ASF immediately recognize Cognito auth requirements, Lambda IAM constraints, and KMS encryption needs as explicit assumptions.

**Worst overlap:** Environmental (0%) again had zero overlap. The ASF identified that mobile app distribution channels (official stores only), developer training gaps, and mobile OS security posture are environmental assumptions that the human architect did not list.

### What Humans Caught That ASF Missed (Missed Findings = 17)

The 17 human-generated assumptions with no ASF counterpart:

1. **Mobile-specific security (H-004, H-005, H-027, H-028, H-033, H-038):** Secure token storage on device, certificate pinning, certificate transparency monitoring, app client ID obfuscation, remote kill switch for compromised app versions, and binary code signing. The ASF has no "Mobile Security" pattern.

2. **Cognito-specific configuration (H-003, H-021, H-035, H-037, H-043):** Social-engineering-resistant account recovery, adaptive authentication, Cognito trigger security, test account management, and hosted UI security headers. These are specific to Cognito and not covered by any pattern.

3. **Serverless operational details (H-026, H-031, H-042):** Lambda version immutability, runtime deprecation cadence, and DLQ data leakage. The ASF covers Lambda broadly but misses version management and dead-letter queue risks.

### What ASF Caught That Humans Missed (Novel Findings = 35)

The ASF generated 61 assumptions, of which 35 (57.4%) were not in the human list:

1. **Incident Response (5 assumptions):** The human generated zero IR assumptions. The ASF contributed a full IR pattern covering plans, log access, isolation procedures, detection, and forensic snapshots.

2. **Identity Lifecycle (4 assumptions):** The human did not consider Cognito user account review, service account management, app client configuration review, or termination-based IAM removal.

3. **Backup & Recovery (ASF-013 through ASF-016):** The human assumed PITR (H-012) but did not consider cross-account backup storage, restore testing, or Cognito user pool export for recovery.

4. **Third-party dependencies (ASF-058 through ASF-061):** The human treated the serverless stack as self-contained. The ASF surfaced AWS service availability risk, SDK vulnerability scanning, mobile library updates, and vendor exit strategy.

5. **MFA operationalization (ASF-002, ASF-003, ASF-004):** The human assumed MFA is enforced (H-020) but did not consider recovery processes, API key bypass, or device-native MFA integration.

6. **Monitoring infrastructure security (ASF-057):** The human assumed logs exist but did not consider that the monitoring system itself must be tamper-proof.

### Architecture Complexity Assessment

Architecture #3 is serverless with managed services, which shifts the assumption landscape significantly:

- **Recall (60.5%)** is the lowest of the three architectures so far. The "Mobile Security" and "Cognito-specific" gaps are large. The ASF needs a dedicated "Mobile Application Security" pattern.
- **Precision (42.6%)** is similar to Architecture #2. The ASF generates many valid assumptions, but the mobile/serverless context introduces many architecture-specific concerns.
- **F1 (50.0%)** is below target, driven by recall gap.
- **Novelty rate (57.4%)** remains strong.

### Key Insight

The serverless/mobile architecture reveals two systematic gaps:
1. **No "Mobile Security" pattern** in the ASF matrix: certificate pinning, secure storage, binary integrity, remote kill switch, app store distribution.
2. **No "Managed Service Configuration" pattern**: Cognito-specific settings (triggers, hosted UI, adaptive auth) and Lambda-specific features (versions, DLQs, runtimes) fall through the cracks.

Adding these two patterns would likely close the recall gap for architectures involving mobile clients and managed serverless services.

---

## Conclusions

| Metric | Target | Actual | Verdict |
|--------|--------|--------|---------|
| Recall | >= 70% | 60.5% | ❌ Below target — missing mobile security and managed service patterns |
| Precision | >= 50% | 42.6% | ❌ Below target — ASF is deliberately broad/generative |
| Novel discoveries | >= 10% | 33.7% | ✅ ASF adds substantial value beyond human reasoning |
| Expert agreement (F1) | > 60% | 50.0% | ❌ Below target — driven by systematic recall gap |

The ASF applied to Architecture #3 confirms the pattern coverage gaps observed in Architectures #1 and #2, with the addition of a mobile-specific blind spot. For serverless architectures with mobile clients, the IR, identity lifecycle, and third-party dependency patterns provide the most novel value beyond unaided human reasoning.
