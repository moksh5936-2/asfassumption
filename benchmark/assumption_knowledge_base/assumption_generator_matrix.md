# ASF Assumption Generator Matrix

Total hidden assumptions: **2188**

## Architecture Patterns

| Pattern | Assumptions |
|---------|------------|
| Authentication (MFA) | 600 |
| Authentication (SSO) | 80 |
| Availability & Resilience | 64 |
| Backup & Recovery | 80 |
| Change Management | 100 |
| Cloud Security (IAM) | 100 |
| Container Security | 100 |
| Data Flow & Classification | 64 |
| Encryption at Rest | 100 |
| Encryption in Transit (TLS) | 100 |
| Endpoint Security | 64 |
| Human Factors & Process | 64 |
| Identity Lifecycle (Provisioning) | 64 |
| Incident Response | 100 |
| Least Privilege | 100 |
| Monitoring & Alerting | 100 |
| Network Segmentation | 80 |
| Physical Security | 64 |
| Supply Chain Security | 64 |
| Third-party Dependency | 100 |

## Categories

| Category | Assumptions |
|----------|------------|
| access | 230 |
| architecture | 64 |
| configuration | 484 |
| dependency | 154 |
| governance | 64 |
| human | 64 |
| identity | 676 |
| network | 70 |
| process | 240 |
| trust | 142 |

## Pattern Details

### Authentication (MFA)

**600 hidden assumptions**

| ID | Component | Depends On | Assumption | Risk | Verification |
|----|-----------|------------|------------|------|-------------|
| MFA-0001 | Primary application login | MFA provider | Must be available for Primary application login using TOTP (... | Denial of access during MFA provider outage... | Monitor MFA provider uptime and configure failover... |
| MFA-0002 | Primary application login | MFA provider | Enrollment must be at 100% for Primary application login usi... | Unenrolled users bypass MFA entirely... | Audit MFA enrollment status across all users monthly... |
| MFA-0003 | Primary application login | MFA provider | MFA is enforced at every auth point for Primary application ... | Partial enforcement creates MFA-free access paths... | Penetration test all authentication points for MFA requireme... |
| MFA-0004 | Primary application login | MFA provider | Bypass paths do not exist for Primary application login usin... | MFA bypass mechanisms are equally exploitable... | Review emergency bypass procedures and access logs... |
| MFA-0005 | Primary application login | MFA provider | MFA fatigue is not exploitable for Primary application login... | Attackers spam MFA prompts until user accepts... | Implement number matching or require entering code from prom... |
| MFA-0006 | Primary application login | MFA provider | Factor is not interceptable for Primary application login us... | Phishing sites can proxy MFA session in real time... | Deploy phishing-resistant MFA (WebAuthn/FIDO2) where possibl... |
| MFA-0007 | Primary application login | MFA provider | MFA re-prompt is enforced for sensitive actions for Primary ... | Long-lived sessions after MFA weaken security... | Require step-up authentication for privilege escalation... |
| MFA-0008 | Primary application login | MFA provider | MFA tokens are provisioned before user needs access for Prim... | User cannot authenticate on day one if token not ready... | Automate MFA token provisioning in employee onboarding... |
| MFA-0009 | Primary application login | MFA provider | MFA recovery workflow is equally secure for Primary applicat... | Recovery via email/SMS is weaker than primary MFA... | Audit MFA recovery method strength against primary factor... |
| MFA-0010 | Primary application login | MFA provider | MFA tokens are revoked on offboarding for Primary applicatio... | Former employee retains MFA token access... | Integrate MFA token revocation into offboarding workflow... |
| ... | *590 more* | | | | |

### Authentication (SSO)

**80 hidden assumptions**

| ID | Component | Depends On | Assumption | Risk | Verification |
|----|-----------|------------|------------|------|-------------|
| SSO-0001 | Identity Provider (IdP) | SSO session | IdP must be reachable for Identity Provider (IdP).... | All SSO-dependent applications become unavailable during IdP... | Test IdP failover and offline access modes... |
| SSO-0002 | Identity Provider (IdP) | SSO session | SSO session lifetime is appropriate for risk level for Ident... | Overly long sessions bypass per-application auth... | Configure session duration based on application risk classif... |
| SSO-0003 | Identity Provider (IdP) | SSO session | SSO covers all applications for Identity Provider (IdP).... | Non-SSO applications use independent credentials... | Audit all applications for SSO integration status quarterly... |
| SSO-0004 | Identity Provider (IdP) | SSO session | SSO tokens are not replayable for Identity Provider (IdP).... | Stolen SAML assertion can be replayed to gain access... | Implement audience restriction, not-before/not-on-or-after v... |
| SSO-0005 | Identity Provider (IdP) | SSO session | IdP trust configuration is correct for Identity Provider (Id... | Misconfigured IdP trust allows unauthorized access... | Review SP metadata, ACS URL, and entity ID configurations... |
| SSO-0006 | Identity Provider (IdP) | SSO session | SAML/OIDC assertions are signed and validated for Identity P... | Unsigned assertions can be forged by attackers... | Enforce assertion signing and validate signatures on every r... |
| SSO-0007 | Identity Provider (IdP) | SSO session | Federation partner metadata is up to date for Identity Provi... | Expired partner certificate breaks federation authentication... | Monitor partner certificate expiry and rotate before expirat... |
| SSO-0008 | Identity Provider (IdP) | SSO session | Single logout propagates to all applications for Identity Pr... | Logging out of IdP does not terminate SP sessions... | Test SLO (Single Log Out) across all SPs regularly... |
| SSO-0009 | Identity Provider (IdP) | SSO session | JIT provisioning creates correct accounts for Identity Provi... | JIT may create accounts with excessive default permissions... | Review JIT provisioning default role assignments... |
| SSO-0010 | Identity Provider (IdP) | SSO session | SSO audit logs attribute actions to specific users for Ident... | Shared service accounts after SSO bypass individual attribut... | Configure IdP to send username attribute in assertions... |
| ... | *70 more* | | | | |

### Availability & Resilience

**64 hidden assumptions**

| ID | Component | Depends On | Assumption | Risk | Verification |
|----|-----------|------------|------------|------|-------------|
| AVL-0001 | P | r | Failover mechanism is operational and tested for P relying o... | Untested failover fails when needed... | Test failover at least quarterly with production-like load... |
| AVL-0002 | P | r | System has sufficient capacity for peak + failover load for ... | Overloaded system during failover cascades to complete outag... | Load test to 2x expected peak and monitor utilization... |
| AVL-0003 | P | r | No single point of failure exists in critical path for P rel... | Single component failure causes complete system outage... | Architecture review for single points of failure... |
| AVL-0004 | P | r | Traffic is distributed across healthy instances for P relyin... | Sticky sessions pin traffic to failing instance... | Verify health check configuration and load balancer draining... |
| AVL-0005 | P | r | Circuit breakers prevent cascade failures for P relying on r... | Failing downstream service degrades upstream consumers... | Implement circuit breakers with appropriate thresholds... |
| AVL-0006 | P | r | Rate limiting protects against traffic spikes for P relying ... | Unmitigated traffic spike overwhelms backend services... | Configure rate limiting at API gateway and application layer... |
| AVL-0007 | P | r | Backpressure mechanisms prevent producer overwhelm for P rel... | Producer overwhelming consumer causes memory exhaustion... | Implement message queue backpressure or reactive streams... |
| AVL-0008 | P | r | System degrades gracefully under failure for P relying on r.... | Non-critical failure causes complete system outage... | Define and test graceful degradation modes for each componen... |
| AVL-0009 | S | t | Failover mechanism is operational and tested for S relying o... | Untested failover fails when needed... | Test failover at least quarterly with production-like load... |
| AVL-0010 | S | t | System has sufficient capacity for peak + failover load for ... | Overloaded system during failover cascades to complete outag... | Load test to 2x expected peak and monitor utilization... |
| ... | *54 more* | | | | |

### Backup & Recovery

**80 hidden assumptions**

| ID | Component | Depends On | Assumption | Risk | Verification |
|----|-----------|------------|------------|------|-------------|
| BAK-0001 | Production database | Backup service | Backup jobs complete successfully within backup window for P... | Failed backups silently leave data unprotected... | Monitor backup job completion status and alert on failures... |
| BAK-0002 | Production database | Backup service | Backup integrity is verified through restore testing for Pro... | Corrupted backups are discovered only during actual disaster... | Perform quarterly restore tests for every critical system... |
| BAK-0003 | Production database | Backup service | Recovery Point Objective is achievable for Production databa... | Backup frequency determines maximum data loss... | Verify backup interval meets RPO requirements... |
| BAK-0004 | Production database | Backup service | Recovery Time Objective is achievable for Production databas... | Full restoration takes longer than business can tolerate... | Measure actual restore time against RTO requirement... |
| BAK-0005 | Production database | Backup service | Backup encryption keys are accessible during recovery for Pr... | Inaccessible encryption keys during disaster equal zero data... | Test key access during recovery drill... |
| BAK-0006 | Production database | Backup service | Backup retention policy is enforced for Production database.... | Premature deletion or over-retention both create risk... | Audit backup lifecycle policies quarterly... |
| BAK-0007 | Production database | Backup service | Backups are stored in geographically separate location for P... | Regional disaster destroys both primary and backup... | Verify cross-region replication is enabled and working... |
| BAK-0008 | Production database | Backup service | Backup storage is logically/physically separate from primary... | Ransomware that reaches primary can also destroy backups... | Implement immutable backup storage with separate access cont... |
| BAK-0009 | Production database | Backup service | Backup SLA is documented and monitored for Production databa... | No SLA means no guaranteed recovery capability... | Establish and monitor backup SLAs per system criticality... |
| BAK-0010 | Production database | Backup service | All critical data is included in backup scope for Production... | Undocumented data sources have no backup protection... | Maintain data inventory and verify backup coverage... |
| ... | *70 more* | | | | |

### Change Management

**100 hidden assumptions**

| ID | Component | Depends On | Assumption | Risk | Verification |
|----|-----------|------------|------------|------|-------------|
| CHG-0001 | Production deployment | Change approval | All changes follow the defined change process for Production... | Unauthorized changes create undocumented configuration drift... | Audit change tickets against actual infrastructure changes... |
| CHG-0002 | Production deployment | Change approval | Emergency changes are reviewed post-hoc for Production deplo... | Emergency changes bypassing review accumulate risk... | Conduct post-emergency change review within 48 hours... |
| CHG-0003 | Production deployment | Change approval | Change documentation accurately describes what changed for P... | Inaccurate documentation misleads incident responders... | Verify change documentation completeness within 24 hours of ... |
| CHG-0004 | Production deployment | Change approval | Changes are tested in non-production before deployment for P... | Untested changes in production cause preventable incidents... | Enforce change promotion from dev → staging → prod pipeline... |
| CHG-0005 | Production deployment | Change approval | Rollback procedure exists and is tested for Production deplo... | Failed change without rollback leads to extended outage... | Include and test rollback plan in every change ticket... |
| CHG-0006 | Production deployment | Change approval | Changes occur within approved change windows for Production ... | Out-of-window changes bypass normal review... | Alert on infrastructure changes outside approved windows... |
| CHG-0007 | Production deployment | Change approval | Configuration drift is detected after change for Production ... | Manual changes outside IaC create unreproducible environment... | Run drift detection after every change and remediate automat... |
| CHG-0008 | Production deployment | Change approval | Change approval comes from authorized approver for Productio... | Rubber-stamped approvals defeat change management purpose... | Rotate change approver assignments and audit approval patter... |
| CHG-0009 | Production deployment | Change approval | Change notifications reach all affected stakeholders for Pro... | Undisclosed changes surprise downstream teams... | Implement automatic notification distribution based on chang... |
| CHG-0010 | Production deployment | Change approval | Changes do not conflict with other planned changes for Produ... | Conflicting changes cause cascading failures... | Require change scheduling coordination for dependent systems... |
| ... | *90 more* | | | | |

### Cloud Security (IAM)

**100 hidden assumptions**

| ID | Component | Depends On | Assumption | Risk | Verification |
|----|-----------|------------|------------|------|-------------|
| CLD-0001 | AWS IAM user | IAM policy | IAM policies allow only required actions for AWS IAM user vi... | Overly permissive IAM policies enable privilege escalation... | Use IAM Access Analyzer to identify unused permissions... |
| CLD-0002 | AWS IAM user | IAM policy | IAM roles are used instead of long-lived user credentials fo... | Long-lived access keys create rotation and exposure risk... | Migrate workloads to instance profiles/IRSA with temporary c... |
| CLD-0003 | AWS IAM user | IAM policy | IAM role trust policies are restrictive for AWS IAM user via... | Overly broad trust policy allows unauthorized account to ass... | Restrict trust policy to specific accounts and external IDs... |
| CLD-0004 | AWS IAM user | IAM policy | IAM access keys are rotated regularly for AWS IAM user via I... | Stale access keys increase blast radius of credential leak... | Automate access key rotation and monitor key age... |
| CLD-0005 | AWS IAM user | IAM policy | CloudTrail is enabled across all regions for AWS IAM user vi... | Activity in un-trailed regions is invisible to security team... | Enable multi-region CloudTrail trail with organization trail... |
| CLD-0006 | AWS IAM user | IAM policy | GuardDuty findings are reviewed and triaged for AWS IAM user... | Critical findings in GuardDuty backlog are missed... | Integrate GuardDuty with SIEM and set up alert routing... |
| CLD-0007 | AWS IAM user | IAM policy | S3 bucket policies are restrictive (no public access) for AW... | Public S3 buckets leak sensitive data continuously... | Enable S3 Block Public Access at account level... |
| CLD-0008 | AWS IAM user | IAM policy | AWS Config rules monitor for compliance drift for AWS IAM us... | Unmonitored configuration changes create undetected policy v... | Implement AWS Config conformance packs for security benchmar... |
| CLD-0009 | AWS IAM user | IAM policy | SCPs enforce guardrails across all accounts for AWS IAM user... | Member accounts can opt out of security controls without SCP... | Apply SCPs to prevent security control disabling... |
| CLD-0010 | AWS IAM user | IAM policy | IAM Access Analyzer is enabled and findings reviewed for AWS... | Unintended resource exposure goes undetected without Access ... | Enable IAM Access Analyzer and review findings weekly... |
| ... | *90 more* | | | | |

### Container Security

**100 hidden assumptions**

| ID | Component | Depends On | Assumption | Risk | Verification |
|----|-----------|------------|------------|------|-------------|
| CNT-0001 | Container base image | Container registry | Container base images are from trusted sources for Container... | Compromised base image injects malware into all derived imag... | Scan base image provenance and enforce allowed registries... |
| CNT-0002 | Container base image | Container registry | Images are scanned for vulnerabilities before deployment for... | Known vulnerabilities in images reach production... | Block deployment of images with critical vulnerabilities... |
| CNT-0003 | Container base image | Container registry | Running containers are scanned for vulnerabilities for Conta... | New vulnerabilities discovered after deployment go undetecte... | Deploy runtime vulnerability scanning in production cluster... |
| CNT-0004 | Container base image | Container registry | Container registry access is controlled for Container base i... | Unauthorized access to registry allows malicious image injec... | Restrict registry push access and enable image signing... |
| CNT-0005 | Container base image | Container registry | Kubernetes RBAC follows least privilege for Container base i... | Overly permissive RBAC allows cluster-wide compromise from s... | Audit RBAC bindings and remove cluster-admin where not neede... |
| CNT-0006 | Container base image | Container registry | Pod security standards (PSP/PSS) are enforced for Container ... | Privileged containers escape container isolation... | Enforce Pod Security Standards at the namespace level... |
| CNT-0007 | Container base image | Container registry | Container runtime is not vulnerable to escape for Container ... | Container escape vulnerability breaks isolation boundary... | Keep container runtime updated and monitor for CVE disclosur... |
| CNT-0008 | Container base image | Container registry | Kubernetes network policies restrict pod communication for C... | Flat pod network allows lateral movement from compromised po... | Implement default-deny network policies for all namespaces... |
| CNT-0009 | Container base image | Container registry | Kubernetes secrets are encrypted at rest for Container base ... | Secrets stored as plaintext in etcd are exposed in backup... | Enable etcd encryption and use external secrets store... |
| CNT-0010 | Container base image | Container registry | Container images are immutable and reproducible for Containe... | Mutable tags (latest) create non-reproducible deployments... | Enforce unique image tags (commit SHA) for all deployments... |
| ... | *90 more* | | | | |

### Data Flow & Classification

**64 hidden assumptions**

| ID | Component | Depends On | Assumption | Risk | Verification |
|----|-----------|------------|------------|------|-------------|
| DFL-0001 | P | I | Data is correctly classified according to policy for P based... | Misclassified data receives inadequate protection... | Audit data classification labels against data content sampli... |
| DFL-0002 | P | I | Encryption level matches data classification for P based on ... | Sensitive data without encryption matching classification le... | Verify encryption configurations align with data classificat... |
| DFL-0003 | P | I | Data flows are mapped and documented for P based on I.... | Undocumented data flows bypass security controls... | Maintain data flow diagrams and update on architecture chang... |
| DFL-0004 | P | I | DLP controls are enforced at data boundaries for P based on ... | Data exfiltration via email, USB, cloud upload goes undetect... | Deploy DLP at key egress points and monitor violations... |
| DFL-0005 | P | I | Data retention policies are enforced per classification for ... | Over-retained data increases breach impact and compliance li... | Automate data lifecycle management based on classification... |
| DFL-0006 | P | I | Data is used only for stated purpose for P based on I.... | Data used beyond stated purpose violates privacy commitments... | Implement data usage auditing and access justification... |
| DFL-0007 | P | I | User consent records are stored and retrievable for P based ... | Inability to prove consent creates GDPR violation risk... | Store consent records with timestamp and scope metadata... |
| DFL-0008 | P | I | Data remains in approved geographic regions for P based on I... | Data stored in restricted region violates data sovereignty l... | Implement data residency controls via cloud provider SCP... |
| DFL-0009 | P | a | Data is correctly classified according to policy for P based... | Misclassified data receives inadequate protection... | Audit data classification labels against data content sampli... |
| DFL-0010 | P | a | Encryption level matches data classification for P based on ... | Sensitive data without encryption matching classification le... | Verify encryption configurations align with data classificat... |
| ... | *54 more* | | | | |

### Encryption at Rest

**100 hidden assumptions**

| ID | Component | Depends On | Assumption | Risk | Verification |
|----|-----------|------------|------------|------|-------------|
| ENCR-0001 | RDS database instance | KMS key | Encryption keys are stored separately from encrypted data fo... | Co-located keys and data renders encryption ineffective... | Verify key management system is separate from data storage... |
| ENCR-0002 | RDS database instance | KMS key | Encryption keys are rotated on schedule for RDS database ins... | Long-lived keys increase blast radius of key compromise... | Automate key rotation and verify rotation policy... |
| ENCR-0003 | RDS database instance | KMS key | Key access is audited for RDS database instance backed by KM... | Unauthorized key usage is undetectable without audit... | Enable key usage logging and monitor for anomalous patterns... |
| ENCR-0004 | RDS database instance | KMS key | Key backup exists and is recoverable for RDS database instan... | Key loss equals permanent data loss... | Test key backup restoration procedure annually... |
| ENCR-0005 | RDS database instance | KMS key | KMS key policies restrict usage to authorized principals for... | Overly permissive key policy allows unauthorized decryption... | Review KMS key policies using IAM Access Analyzer... |
| ENCR-0006 | RDS database instance | KMS key | Encryption algorithm is industry standard for RDS database i... | Custom or outdated algorithms may be breakable... | Verify AES-256 or equivalent is used, not DES/RC4/TDEA... |
| ENCR-0007 | RDS database instance | KMS key | Key management system is available for RDS database instance... | KMS outage blocks all encrypted data access... | Test KMS failover and understand blast radius of KMS region ... |
| ENCR-0008 | RDS database instance | KMS key | Key compromise incident response exists for RDS database ins... | No process for key revocation leaves data exposed... | Document and test key compromise response procedure... |
| ENCR-0009 | RDS database instance | KMS key | Encryption mode is secure (GCM, not ECB) for RDS database in... | ECB mode leaks data patterns in encryption output... | Verify encryption algorithm uses authenticated encryption mo... |
| ENCR-0010 | RDS database instance | KMS key | HSM is not physically compromised for RDS database instance ... | HSM tampering exposes all keys protected by it... | Verify HSM is in tamper-resistant enclosure with audit trail... |
| ... | *90 more* | | | | |

### Encryption in Transit (TLS)

**100 hidden assumptions**

| ID | Component | Depends On | Assumption | Risk | Verification |
|----|-----------|------------|------------|------|-------------|
| TLS-0001 | Public-facing web application | TLS certificate | TLS certificates are valid and not expired for Public-facing... | Expired certificate causes service outage or insecure fallba... | Monitor certificate expiry and auto-renew with 30-day notice... |
| TLS-0002 | Public-facing web application | TLS certificate | TLS is enforced on ALL connections, not optional for Public-... | Downgrade attacks can force plaintext fallback... | Enforce HSTS and reject non-TLS connections at load balancer... |
| TLS-0003 | Public-facing web application | TLS certificate | Minimum TLS version 1.2 is enforced for Public-facing web ap... | TLS 1.0/1.1 are vulnerable to protocol downgrade attacks... | Configure server to reject TLS versions below 1.2... |
| TLS-0004 | Public-facing web application | TLS certificate | Cipher suites are strong (no RC4, 3DES, CBC) for Public-faci... | Weak ciphers enable decryption by sophisticated attackers... | Audit cipher suite configuration against Mozilla SSL guideli... |
| TLS-0005 | Public-facing web application | TLS certificate | mTLS is used for service-to-service communication for Public... | Without mTLS, server cannot verify client identity... | Implement mTLS via service mesh or certificate-based auth... |
| TLS-0006 | Public-facing web application | TLS certificate | Certificate revocation is checked for Public-facing web appl... | Compromised certificate continues to be trusted... | Enable OCSP stapling or CRL checking on all endpoints... |
| TLS-0007 | Public-facing web application | TLS certificate | Private keys are protected (not in code, not world-readable)... | Exposed private key allows traffic decryption and impersonat... | Store private keys in HSM or secrets manager with access aud... |
| TLS-0008 | Public-facing web application | TLS certificate | Certificate chain is complete and trusted by clients for Pub... | Incomplete cert chain causes browser/client trust warnings... | Verify full certificate chain is served including intermedia... |
| TLS-0009 | Public-facing web application | TLS certificate | SNI routing is correctly configured for multi-domain for Pub... | Incorrect SNI routing sends traffic to wrong backend... | Test SNI routing for all hostnames on shared IP... |
| TLS-0010 | Public-facing web application | TLS certificate | OCSP responder is reachable for real-time validation for Pub... | Unreachable OCSP responder fails open (accepts all) or close... | Configure OCSP stapling to avoid OCSP verification failures... |
| ... | *90 more* | | | | |

### Endpoint Security

**64 hidden assumptions**

| ID | Component | Depends On | Assumption | Risk | Verification |
|----|-----------|------------|------------|------|-------------|
| END-0001 | Corporate laptop | EDR agent | EDR agent is installed on all endpoints for Corporate laptop... | Unprotected endpoint is invisible to security team during in... | Audit EDR agent coverage against endpoint inventory... |
| END-0002 | Corporate laptop | EDR agent | EDR agent is running and reporting for Corporate laptop via ... | Disabled EDR agent provides no protection... | Monitor EDR agent heartbeat and alert on check-in failures... |
| END-0003 | Corporate laptop | EDR agent | Anti-malware signatures are current for Corporate laptop via... | Outdated signatures miss recent malware variants... | Verify signature age < 24 hours across all endpoints... |
| END-0004 | Corporate laptop | EDR agent | OS and application patches are applied within SLA for Corpor... | Unpatched vulnerabilities are actively exploited in the wild... | Deploy patch management with automated enforcement... |
| END-0005 | Corporate laptop | EDR agent | Full disk encryption is enabled for Corporate laptop via EDR... | Lost/stolen device exposes all data without FDE... | Verify FDE status via MDM or EDR reporting... |
| END-0006 | Corporate laptop | EDR agent | Local admin rights are restricted for Corporate laptop via E... | Users with local admin can disable security controls... | Remove local admin rights and implement LAPS for admin passw... |
| END-0007 | Corporate laptop | EDR agent | Application allowlisting/blacklisting is enforced for Corpor... | Users can install unapproved applications with malware risk... | Deploy application control policies via MDM/GPO... |
| END-0008 | Corporate laptop | EDR agent | USB device control is enforced for Corporate laptop via EDR ... | Malware introduced via USB bypasses network controls... | Restrict USB mass storage and enable audit logging... |
| END-0009 | Server endpoint | Anti-malware | EDR agent is installed on all endpoints for Server endpoint ... | Unprotected endpoint is invisible to security team during in... | Audit EDR agent coverage against endpoint inventory... |
| END-0010 | Server endpoint | Anti-malware | EDR agent is running and reporting for Server endpoint via A... | Disabled EDR agent provides no protection... | Monitor EDR agent heartbeat and alert on check-in failures... |
| ... | *54 more* | | | | |

### Human Factors & Process

**64 hidden assumptions**

| ID | Component | Depends On | Assumption | Risk | Verification |
|----|-----------|------------|------------|------|-------------|
| HUM-0001 | S | e | Security training changes employee behavior for S via e.... | Training completion without behavior change provides no risk... | Measure training effectiveness via phishing simulation click... |
| HUM-0002 | S | e | Employees can identify and report phishing attempts for S vi... | Phishing remains primary initial access vector if employees ... | Conduct quarterly phishing simulations with training for cli... |
| HUM-0003 | S | e | Insider threat detection detects anomalous behavior for S vi... | Malicious insider activity goes undetected without behaviora... | Implement UEBA with baseline of normal behavior per role... |
| HUM-0004 | S | e | Employees are aware of and acknowledge security policies for... | Unaware employees cannot comply with policies... | Require annual policy acknowledgment with quiz on key polici... |
| HUM-0005 | S | e | Access request approvers make informed decisions for S via e... | Rubber-stamped approvals grant unnecessary access... | Require approver to confirm access justification before appr... |
| HUM-0006 | S | e | Employees report security incidents without fear of blame fo... | Unreported incidents delay containment and increase damage... | Establish non-punitive incident reporting culture and anonym... |
| HUM-0007 | S | e | Employees follow security procedures in practice for S via e... | Security procedures that impede work are routinely bypassed... | Conduct process compliance audits and remove friction points... |
| HUM-0008 | S | e | Employees can resist social engineering tactics for S via e.... | Social engineering bypasses technical controls via human man... | Conduct social engineering simulations (phone, in-person, em... |
| HUM-0009 | P | h | Security training changes employee behavior for P via h.... | Training completion without behavior change provides no risk... | Measure training effectiveness via phishing simulation click... |
| HUM-0010 | P | h | Employees can identify and report phishing attempts for P vi... | Phishing remains primary initial access vector if employees ... | Conduct quarterly phishing simulations with training for cli... |
| ... | *54 more* | | | | |

### Identity Lifecycle (Provisioning)

**64 hidden assumptions**

| ID | Component | Depends On | Assumption | Risk | Verification |
|----|-----------|------------|------------|------|-------------|
| IDM-0001 | Employee onboarding | HR system | Accounts are provisioned before user needs access for Employ... | Day-one productivity blocked by access delays... | Measure and monitor time from HR trigger to account creation... |
| IDM-0002 | Employee onboarding | HR system | HR data (department, role, manager) is accurate for Employee... | Incorrect HR data propagates to incorrect access grants... | Audit HR data accuracy against employee records quarterly... |
| IDM-0003 | Employee onboarding | HR system | HR-IdP integration operates in real time for Employee onboar... | Delayed sync creates window where access does not match role... | Monitor sync latency between HR system and IdP... |
| IDM-0004 | Employee onboarding | HR system | Provisioned roles match job function for Employee onboarding... | Generic onboarding roles grant excessive default permissions... | Review default onboarding role permissions quarterly... |
| IDM-0005 | Employee onboarding | HR system | Contractor accounts have automatic expiry for Employee onboa... | Expired contractor accounts remain active indefinitely... | Verify contractor account expiry is enforced within IdP... |
| IDM-0006 | Employee onboarding | HR system | Identity proofing is adequate for access level for Employee ... | Weak identity proofing allows impersonation during account c... | Review identity proofing requirements based on access level... |
| IDM-0007 | Employee onboarding | HR system | Duplicate identities are detected and merged for Employee on... | Multiple identities for same person create audit gaps... | Implement duplicate identity detection and reconciliation... |
| IDM-0008 | Employee onboarding | HR system | Identity source of truth is authoritative and protected for ... | Multiple identity sources create conflicts and gaps... | Document identity source of truth for each system... |
| IDM-0009 | Contractor onboarding | HR system | Accounts are provisioned before user needs access for Contra... | Day-one productivity blocked by access delays... | Measure and monitor time from HR trigger to account creation... |
| IDM-0010 | Contractor onboarding | HR system | HR data (department, role, manager) is accurate for Contract... | Incorrect HR data propagates to incorrect access grants... | Audit HR data accuracy against employee records quarterly... |
| ... | *54 more* | | | | |

### Incident Response

**100 hidden assumptions**

| ID | Component | Depends On | Assumption | Risk | Verification |
|----|-----------|------------|------------|------|-------------|
| IR-0001 | Detection alert | SIEM rule | IR plan is current and reflects current architecture for Det... | Outdated IR plan references decommissioned systems and conta... | Review and test IR plan at least annually... |
| IR-0002 | Detection alert | SIEM rule | IR team roles are assigned with backups for Detection alert ... | Unclear role assignment causes confusion during active incid... | Maintain on-call rotation with documented escalation paths... |
| IR-0003 | Detection alert | SIEM rule | Communication channels are established and tested for Detect... | Failure to reach responders delays containment... | Test IR communication channels (phone, Slack, radio) quarter... |
| IR-0004 | Detection alert | SIEM rule | Forensic and containment tooling is pre-deployed for Detecti... | Time spent acquiring tools during incident is time lost for ... | Pre-deploy forensic agents and containment automation... |
| IR-0005 | Detection alert | SIEM rule | Evidence chain of custody is maintained for Detection alert ... | Inadmissible evidence prevents legal action against attacker... | Document evidence handling procedures and train IR team... |
| IR-0006 | Detection alert | SIEM rule | Playbooks cover current attack scenarios for Detection alert... | Incident requiring un-practiced response has longer containm... | Test and update playbooks for top 5 attack scenarios annuall... |
| IR-0007 | Detection alert | SIEM rule | Remediation steps are tested before execution for Detection ... | Untested remediation causes additional damage... | Sandbox-test remediation scripts before production execution... |
| IR-0008 | Detection alert | SIEM rule | Third-party (LE, PR, legal) coordination is documented for D... | Third-party coordination delays create public relations dama... | Pre-establish contacts with legal counsel and law enforcemen... |
| IR-0009 | Detection alert | SIEM rule | Recovery verification confirms attacker access is removed fo... | Premature recovery declaration leads to reinfection... | Define recovery verification criteria per incident type... |
| IR-0010 | Detection alert | SIEM rule | Lessons learned from incidents are implemented as improvemen... | Repeated incidents recur because corrective actions were not... | Track lessons-learned action items with accountability... |
| ... | *90 more* | | | | |

### Least Privilege

**100 hidden assumptions**

| ID | Component | Depends On | Assumption | Risk | Verification |
|----|-----------|------------|------------|------|-------------|
| LPO-0001 | I | A | Assigned permissions match actual required permissions for I... | Excessive permissions violate least privilege... | Review IAM Access Advisor metrics and rightsize permissions... |
| LPO-0002 | I | A | Permission boundaries are enforced for I scoped by Permissio... | IAM boundary does not restrict permission escalation paths... | Test that permission boundary stops privilege escalation... |
| LPO-0003 | I | A | Permission requests include business justification for I sco... | Unchallenged permission requests accumulate unchecked... | Require ticket number and approval for every permission gran... |
| LPO-0004 | I | A | Temporary permissions have expiration for I scoped by Permis... | Time-unlimited elevated permissions become standing privileg... | Implement temporary credential mechanisms (STS, JIT)... |
| LPO-0005 | I | A | Permission creep is detected and remediated for I scoped by ... | Users accumulate permissions over time without review... | Run access review automation comparing current vs required p... |
| LPO-0006 | I | A | Privilege elevation requires separate auth for I scoped by P... | Always-elevated sessions increase blast radius of compromise... | Implement just-in-time elevation with audit trail... |
| LPO-0007 | I | A | Break-glass emergency access is monitored for I scoped by Pe... | Emergency access procedures can be abused for unauthorized a... | Audit all break-glass usage within 24 hours and review month... |
| LPO-0008 | I | A | Group/permission inheritance does not create unintended acce... | Nested group membership creates hidden privilege escalation... | Review group hierarchy for unintended inheritance paths... |
| LPO-0009 | I | A | Cross-account access is scoped and reviewed for I scoped by ... | Compromised account can access resources in other accounts... | Review cross-account IAM roles and trust policies quarterly... |
| LPO-0010 | I | A | Default deny is enforced for all access requests for I scope... | Allow-by-default access creates standing permissions... | Verify IAM policy evaluation results in explicit deny by def... |
| ... | *90 more* | | | | |

### Monitoring & Alerting

**100 hidden assumptions**

| ID | Component | Depends On | Assumption | Risk | Verification |
|----|-----------|------------|------------|------|-------------|
| MON-0001 | Production application logs | Log agent + pipeline | All critical systems are sending logs for Production applica... | Unlogged systems are invisible during incident investigation... | Audit log source coverage against asset inventory quarterly... |
| MON-0002 | Production application logs | Log agent + pipeline | Log storage has sufficient capacity for retention period for... | Logs rotated before compliance period due to capacity limits... | Monitor log storage utilization and auto-scale or alert... |
| MON-0003 | Production application logs | Log agent + pipeline | Logs are tamper-proof (immutable) for Production application... | Attacker covering tracks can modify logs to hide activity... | Implement immutable log storage with separate admin access c... |
| MON-0004 | Production application logs | Log agent + pipeline | Alerts are configured for key detection events for Productio... | Security events without alerts are invisible... | Review alert rule coverage against ATT&CK framework... |
| MON-0005 | Production application logs | Log agent + pipeline | Alerts are triaged within defined SLA for Production applica... | Un-triaged alerts accumulate and critical events are missed... | Monitor mean time to triage (MTTT) and alert backlog... |
| MON-0006 | Production application logs | Log agent + pipeline | Log retention meets compliance requirements for Production a... | Logs deleted before compliance period create audit failures... | Verify log retention settings against compliance calendar... |
| MON-0007 | Production application logs | Log agent + pipeline | Monitoring covers all environments (prod, staging, dev) for ... | Security events in non-prod environments go undetected... | Extend monitoring agent deployment to all environments... |
| MON-0008 | Production application logs | Log agent + pipeline | Log correlation across sources is operational for Production... | Siloed logs miss multi-stage attack sequences... | Test cross-source correlation scenarios quarterly... |
| MON-0009 | Production application logs | Log agent + pipeline | Anomaly detection has baseline of normal behavior for Produc... | Without baseline, anomaly detection generates excessive nois... | Establish and tune behavioral baselines per system... |
| MON-0010 | Production application logs | Log agent + pipeline | Alerts integrate with response workflow (PagerDuty, Slack, e... | Alerts that nobody sees might as well not exist... | Test alert notification delivery path monthly... |
| ... | *90 more* | | | | |

### Network Segmentation

**80 hidden assumptions**

| ID | Component | Depends On | Assumption | Risk | Verification |
|----|-----------|------------|------------|------|-------------|
| SEG-0001 | Production VPC | Internet Gateway | Segment is isolated from unauthorized networks for Productio... | Cross-segment traffic bypasses security controls... | Test network ACLs and security groups with active scanning... |
| SEG-0002 | Production VPC | Internet Gateway | Route tables direct traffic to correct destinations for Prod... | Misrouted traffic leaves the network boundary... | Audit route tables for unintended default routes... |
| SEG-0003 | Production VPC | Internet Gateway | Network ACLs filter traffic in both directions for Productio... | Stateless ACLs only filtering inbound creates outbound bypas... | Verify NACL rules cover both inbound and outbound directions... |
| SEG-0004 | Production VPC | Internet Gateway | Security groups are scoped to specific CIDRs, not 0.0.0.0/0 ... | Overly permissive security groups expose services to interne... | Review all security groups with 0.0.0.0/0 inbound rules... |
| SEG-0005 | Production VPC | Internet Gateway | VPC peering does not extend beyond approved accounts for Pro... | Peering across accounts creates unauthorized network paths... | Audit VPC peering connections quarterly... |
| SEG-0006 | Production VPC | Internet Gateway | VPC endpoints are configured for private access for Producti... | Traffic to AWS services traverses internet instead of AWS ne... | Verify VPC endpoint policies and route tables... |
| SEG-0007 | Production VPC | Internet Gateway | No subnet is unintentionally public for Production VPC via I... | Subnet with IGW route + public IP is directly internet acces... | Scan for subnets with 0.0.0.0/0 route to IGW... |
| SEG-0008 | Production VPC | Internet Gateway | Network flow logs are enabled for critical segments for Prod... | Without flow logs, malicious traffic is invisible... | Enable VPC flow logs for all production subnets... |
| SEG-0009 | Production VPC | NAT Gateway | Segment is isolated from unauthorized networks for Productio... | Cross-segment traffic bypasses security controls... | Test network ACLs and security groups with active scanning... |
| SEG-0010 | Production VPC | NAT Gateway | Route tables direct traffic to correct destinations for Prod... | Misrouted traffic leaves the network boundary... | Audit route tables for unintended default routes... |
| ... | *70 more* | | | | |

### Physical Security

**64 hidden assumptions**

| ID | Component | Depends On | Assumption | Risk | Verification |
|----|-----------|------------|------------|------|-------------|
| PHY-0001 | D | a | Physical access is restricted to authorized personnel only f... | Unauthorized physical access leads to data theft or destruct... | Audit physical access logs and revoke departed employee badg... |
| PHY-0002 | D | a | Tailgating (piggybacking) is prevented for D via a.... | Tailgating bypasses badge access control completely... | Implement mantraps or turnstiles at data center entrances... |
| PHY-0003 | D | a | Visitor access is logged and escorted for D via a.... | Unescorted visitor has unrestricted physical access... | Enforce visitor check-in process and escort policy... |
| PHY-0004 | D | a | Biometric system has alternate authentication method for D v... | Biometric failure denies access or forces security-compromis... | Maintain alternative access method and emergency access prot... |
| PHY-0005 | D | a | Temperature and humidity are monitored for D via a.... | Cooling failure causes hardware damage and data loss... | Implement environmental monitoring with automated alerting... |
| PHY-0006 | D | a | UPS and generator provide continuous power for D via a.... | Power loss causes abrupt system shutdown and data corruption... | Test UPS battery capacity and generator under load quarterly... |
| PHY-0007 | D | a | Fire suppression system is operational for D via a.... | Fire destroys hardware and data permanently... | Inspect fire suppression system and test detection annually... |
| PHY-0008 | D | a | Vendor maintenance access is logged and supervised for D via... | Vendor maintenance provides opportunity for data exfiltratio... | Supervise all third-party maintenance and log activities... |
| PHY-0009 | S | e | Physical access is restricted to authorized personnel only f... | Unauthorized physical access leads to data theft or destruct... | Audit physical access logs and revoke departed employee badg... |
| PHY-0010 | S | e | Tailgating (piggybacking) is prevented for S via e.... | Tailgating bypasses badge access control completely... | Implement mantraps or turnstiles at data center entrances... |
| ... | *54 more* | | | | |

### Supply Chain Security

**64 hidden assumptions**

| ID | Component | Depends On | Assumption | Risk | Verification |
|----|-----------|------------|------------|------|-------------|
| SCS-0001 | O | p | Open source dependencies are scanned for vulnerabilities for... | Known CVEs in dependencies are exploited in production... | Run SCA scanning in CI/CD pipeline and block critical vulner... |
| SCS-0002 | O | p | Base images are from verified sources with known provenance ... | Compromised base image from untrusted source supplies chain ... | Enforce image provenance attestation and trusted registry... |
| SCS-0003 | O | p | CI/CD pipeline is hardened against injection for O sourced f... | Compromised CI/CD pipeline deploys malicious code... | Implement CI/CD pipeline security controls (least privilege,... |
| SCS-0004 | O | p | Artifacts are signed and signatures verified before deployme... | Unsigned artifact deployment enables supply chain substituti... | Implement Sigstore/cosign for artifact signing and verificat... |
| SCS-0005 | O | p | Signed commits are verified before merge for O sourced from ... | Unverified signed commits provide no authenticity guarantee... | Enforce commit signature verification in SCM branch rules... |
| SCS-0006 | O | p | Third-party SDKs are from verified distribution channels for... | Typosquatted SDK with malicious code executed in production... | Maintain software bill of materials (SBOM) for all dependenc... |
| SCS-0007 | O | p | Dependency versions are pinned (no floating tags) for O sour... | Floating tag automatically pulls malicious updated version... | Pin dependencies to specific versions and audit changes... |
| SCS-0008 | O | p | Software suppliers are assessed for security practices for O... | Supplier with poor security introduces vulnerabilities... | Conduct supplier security assessment before procurement... |
| SCS-0009 | C | o | Open source dependencies are scanned for vulnerabilities for... | Known CVEs in dependencies are exploited in production... | Run SCA scanning in CI/CD pipeline and block critical vulner... |
| SCS-0010 | C | o | Base images are from verified sources with known provenance ... | Compromised base image from untrusted source supplies chain ... | Enforce image provenance attestation and trusted registry... |
| ... | *54 more* | | | | |

### Third-party Dependency

**100 hidden assumptions**

| ID | Component | Depends On | Assumption | Risk | Verification |
|----|-----------|------------|------------|------|-------------|
| VEN-0001 | C | l | Vendor security posture remains acceptable over time for C.... | Vendor breach becomes your breach... | Review vendor security assessments annually... |
| VEN-0002 | C | l | Vendor SLA commitments are met for C.... | SLA violations impact business operations without compensati... | Monitor SLA metrics and enforce penalty clauses... |
| VEN-0003 | C | l | Vendor data handling complies with policy for C.... | Vendor data processing may violate GDPR/HIPAA/PCI requiremen... | Audit vendor data processing agreements and certifications... |
| VEN-0004 | C | l | Vendor breach notification is timely for C.... | Undisclosed vendor breach creates undetected exposure window... | Contractually require 24-72 hour breach notification... |
| VEN-0005 | C | l | Vendor access is promptly revoked when engagement ends for C... | Former vendor retains access to production systems... | Automate vendor account deprovisioning on contract end... |
| VEN-0006 | C | l | Vendor dependencies are documented for C.... | Undocumented vendor dependency creates single point of failu... | Maintain vendor dependency registry and review quarterly... |
| VEN-0007 | C | l | Vendor exit strategy exists for C.... | Inability to leave vendor creates vendor lock-in risk... | Document and test vendor transition plan annually... |
| VEN-0008 | C | l | Vendor sub-processors are disclosed and assessed for C.... | Vendor use of sub-processors bypasses security assessment... | Contractually require vendor to disclose all sub-processors... |
| VEN-0009 | C | l | Vendor compliance certifications remain valid for C.... | Expired vendor certification creates compliance gap for down... | Monitor vendor certification expiry dates... |
| VEN-0010 | C | l | Vendor business continuity plan is tested for C.... | Vendor outage impacts your operations without recourse... | Request vendor BCP test results annually... |
| ... | *90 more* | | | | |
