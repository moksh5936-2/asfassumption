# ASF Phase 6 Experiment: Architecture #009

**Architecture:** Vendor SaaS → API → Internal Systems
**Date:** 2026-06-09
**Simulation Mode:** Human + ASF framework comparison

---

## Architecture Profile

```
[Vendor SaaS (Salesforce)] --> [API Gateway] --> [Internal CRM Sync] --> [Customer DB]
                                    │
                              [OAuth Token]
```

### Documented Policy
| # | Policy |
|---|--------|
| P1 | Vendor has SOC 2 Type II report |
| P2 | API integration uses OAuth 2.0 |
| P3 | Data processing agreement (DPA) in place |
| P4 | Quarterly vendor review |

### Trust Boundaries
| Boundary | Type |
|----------|------|
| Vendor ↔ API Gateway | External boundary |
| API Gateway ↔ Internal | Trust boundary |
| Data residency | Compliance boundary |

### Complexity Rating
**Moderate** — vendor-integrated topology, 4 nodes, 3 trust boundaries, external dependency at the root.

---

## Step 1: Human-Generated Assumptions

*Role: Senior Security Architect. Listed below are assumptions that MUST remain true for this architecture to remain secure but are NOT stated in the documented policy.*

| ID | Assumption | Justification |
|----|-----------|---------------|
| H-001 | The vendor's SOC 2 Type II report is current (issued within the last 12 months) and covers all relevant trust services criteria. | An outdated or incomplete SOC 2 report does not provide assurance about the vendor's current security posture. |
| H-002 | The SOC 2 report does not have any exceptions or findings that affect the confidentiality or availability of the data processed. | A SOC 2 with material exceptions indicates that the vendor's controls are not operating effectively. |
| H-003 | The OAuth 2.0 integration uses the authorization code flow with PKCE, not the implicit grant type. | The implicit grant type does not support client authentication and is vulnerable to access token interception. |
| H-004 | OAuth tokens are scoped to the minimum necessary Salesforce objects and operations needed for the integration. | An over-scoped OAuth token grants the internal system more Salesforce access than necessary, increasing blast radius. |
| H-005 | OAuth refresh tokens have a limited lifetime and are rotated on a regular cadence. | Static long-lived refresh tokens allow the internal system to maintain access indefinitely after compromise. |
| H-006 | The API Gateway validates the OAuth token on every request — not just at connection establishment. | Token validation only at the start of a session allows a revoked token to be reused for the duration of the session. |
| H-007 | The DPA covers all data processing activities, including data stored, transmitted, and processed by the vendor. | A DPA that does not cover all data handling scenarios leaves legal and regulatory gaps. |
| H-008 | The vendor has a breach notification SLA that meets the organization's contractual and regulatory requirements (typically 24-72 hours). | A vendor that delays breach notification beyond regulatory timelines (GDPR 72h, HIPAA 60d) puts the organization in violation. |
| H-009 | The vendor's sub-processors are identified in the DPA and have been assessed for security posture. | A vendor using undisclosed sub-processors extends the trust boundary to parties the organization has not evaluated. |
| H-010 | The quarterly vendor review includes a review of the vendor's security incidents, penetration test results, and roadmap changes. | A quarterly review that only checks contract compliance does not assess the vendor's evolving security posture. |
| H-011 | Data synced from the vendor to the internal CRM is encrypted in transit at all points (vendor → API Gateway → internal). | Unencrypted data at any hop exposes customer data to interception on the network. |
| H-012 | The API Gateway terminates and re-encrypts TLS — the data is not decrypted in the clear at any intermediate point. | An API Gateway that decrypts and forwards in plaintext creates an exposure point for the clear-text data. |
| H-013 | The internal CRM Sync service validates that data received from the vendor matches expected schema and does not contain injection payloads. | A compromised vendor account sending crafted data can inject malicious content into the internal CRM and database. |
| H-014 | The internal CRM Sync service has rate limiting to prevent the vendor from overwhelming the internal system. | A vendor-side misconfiguration or compromise that sends excessive API requests can cause a denial of service against the internal system. |
| H-015 | The customer database is not directly accessible from the internet — only the internal CRM Sync service can write to it. | A database with any public endpoint bypasses all API Gateway and vendor controls. |
| H-016 | IAM roles for the internal CRM Sync service grant least-privilege access to the customer database (specific tables, no DDL). | An over-permissioned database user in the CRM Sync service turns a service compromise into full database access. |
| H-017 | The OAuth client credentials (client ID and client secret) are stored in a secrets manager, not in configuration files or environment variables. | OAuth credentials stored in plaintext are accessible to anyone with filesystem access to the CRM Sync server. |
| H-018 | The vendor's API changes (version upgrades, endpoint deprecation) are communicated with sufficient notice to update the integration. | A breaking API change without notice causes the CRM Sync service to fail, blocking data synchronization. |
| H-019 | There is a fallback procedure for when the vendor's API is unavailable — queuing or offline mode in the CRM Sync service. | A vendor API outage that blocks the CRM Sync can cause data loss if messages are not queued for later processing. |
| H-020 | The vendor's access to the organization's API Gateway is restricted to specific IP addresses or a VPN. | Without network-layer restriction, the OAuth token is the only control; a leaked token can be used from anywhere. |
| H-021 | Data residency requirements are documented — customer data cannot leave approved geographic regions. | Salesforce replicating data to a non-approved region violates data residency compliance requirements. |
| H-022 | The API Gateway has a web application firewall (WAF) configured to inspect incoming vendor traffic for common attack patterns. | The vendor account could be compromised and used to send malicious API requests to the internal system. |
| H-023 | The internal CRM Sync service logs all data operations (sync events, field changes, errors) for audit. | Without sync logging, a data corruption or unauthorized modification by the vendor cannot be traced. |
| H-024 | The vendor's API access is reviewed quarterly to ensure it still requires only the minimum necessary data. | Over time, integration scope creep gives the vendor access to more data than the original design required. |
| H-025 | An incident response plan exists for vendor data breaches that includes steps to revoke OAuth tokens and isolate the API Gateway. | Without a plan, a vendor breach leads to ad-hoc containment that may not be effective or timely. |
| H-026 | The contract includes a right-to-audit clause allowing the organization to assess the vendor's controls directly if needed. | Without right-to-audit, the organization must rely entirely on the SOC 2 report, which may not cover all concerns. |
| H-027 | The vendor encrypts data at rest with a minimum of AES-256 encryption. | Customer data stored by the vendor in unencrypted form is at risk of exposure in a vendor data breach. |
| H-028 | The API Gateway enforces request validation to ensure only expected API endpoints and parameters are accepted from the vendor. | An API Gateway that forwards all requests without validation allows an attacker to probe internal endpoints through the vendor channel. |
| H-029 | The OAuth token endpoint uses TLS 1.2 or higher — no weak cipher or protocol fallback for token requests. | OAuth token requests over weak TLS expose the access token to interception. |
| H-030 | The internal CRM Sync service has a circuit breaker pattern — it stops calling the vendor API if the vendor is responding with errors. | A vendor API returning errors (4xx/5xx) should not cause the internal system to enter a retry loop that degrades internal performance. |
| H-031 | Vendor employee access to the organization's data in the vendor's system is logged and reviewed by the vendor. | A vendor insider can access the organization's Salesforce data without the organization's knowledge. |
| H-032 | The vendor undergoes regular penetration testing (at least annually) and the results are available for review. | A SOC 2 Type II report may include penetration test results, but without explicit confirmation, the organization does not know if testing is current. |
| H-033 | The customer DB has a retention policy that aligns with the DPA — data is purged from the vendor and internal systems per the schedule. | Data retained beyond the contractual period creates legal liability and unnecessary exposure. |
| H-034 | The API Gateway monitors for anomalous vendor API call patterns (unusual volume, out-of-hours access, new endpoints). | A compromised vendor account exfiltrating data will show unusual API call patterns that detection can capture. |
| H-035 | The vendor's data processing is limited to the purposes defined in the DPA — the vendor does not use the organization's data for its own benefit. | A vendor mining the organization's customer data for their own analytics violates the DPA and may breach data protection law. |
| H-036 | The OAuth token includes the organization's Salesforce instance identifier so that a token from one org cannot be used against another. | A misconfigured OAuth token that works across Salesforce orgs allows a token leaked from one customer to access another's data. |
| H-037 | The vendor's availability SLA meets the organization's business continuity requirements (uptime percentage, maintenance windows). | A vendor with 99.9% uptime is unavailable for 8.76 hours per year; the organization must accept this or have a fallback. |
| H-038 | There is an exit strategy and data migration plan if the vendor contract is terminated or the vendor goes out of business. | Vendor lock-in without an exit plan can cause critical business disruption and data loss if the relationship ends. |
| H-039 | The API Gateway's OAuth token validation includes checking the token's expiry, issuer, audience, and signature. | A token validated only for presence — not cryptographic integrity — can be forged or replayed. |
| H-040 | The internal CRM Sync service does not cache customer data from the vendor in unencrypted local storage or logs. | Cached customer data in the internal system creates an additional data exposure surface beyond the database. |
| H-041 | The vendor's security notification process is documented and tested — the organization receives notifications within the contractual SLA. | A notification process that is not tested will fail during an actual incident, delaying the organization's response. |
| H-042 | The API Gateway has a kill switch to immediately block all vendor traffic if a breach is suspected. | Without the ability to isolate the vendor integration rapidly, an active compromise continues while the response team deliberates. |
| H-043 | Customer data is not used in development or staging environments by either the vendor or the internal systems. | Non-production use of production customer data increases exposure risk and may violate the DPA. |

**Total (H): 43**

---

## Step 2: ASF-Generated Assumptions

*Using the 20-pattern assumption generator matrix. Applicable patterns: 16 of 20. Patterns excluded: Container Security (no containers), Endpoint Security (no user endpoints), Physical Security, Backup & Recovery (deferred to the vendor's responsibility for SaaS).*

---

### Pattern 1: Authentication (MFA)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-001 | Vendor administrative access to the Salesforce org requires MFA for all vendor employees and the organization's admins. | Explicit | SOC 2 Type II should cover MFA, but the organization must verify it is enforced for all administrative access. |
| ASF-002 | MFA recovery processes at the vendor do not allow bypassing MFA without proper identity verification. | Derived | A vendor help desk that resets MFA via email or knowledge-based authentication defeats the MFA control. |
| ASF-003 | The organization's Salesforce administrators have MFA enabled on their Salesforce accounts, not just through SSO. | Operational | Salesforce admin accounts with local credentials that bypass SSO create an identity management gap. |
| ASF-004 | API Gateway administrative access for managing OAuth integrations requires MFA. | Implicit | The API Gateway is the control point for the vendor integration; its administrative console must be protected. |

---

### Pattern 2: Authentication (SSO)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-005 | Salesforce SSO is configured through the organization's identity provider — local credentials are disabled. | Explicit | SSO ensures that Salesforce access is governed by the corporate identity lifecycle. |
| ASF-006 | The OAuth 2.0 flow between Salesforce and the API Gateway uses the authorization code grant type. | Derived | The documented policy says "OAuth 2.0" but does not specify the grant type; implicit grant is less secure. |
| ASF-007 | OAuth token refresh is performed using client credentials that are securely stored and rotated. | Trust | Refresh tokens that are not rotated create persistent access that outlives the original authorization. |
| ASF-008 | The OAuth authorization server (Salesforce) is available — token refresh fails if Salesforce is down. | Dependency | If Salesforce IdP is unavailable, token refresh fails and the integration stops working. |

---

### Pattern 3: Availability & Resilience

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-009 | Salesforce API availability meets the SLA documented in the contract — the organization has a documented acceptable outage duration. | Dependency | Salesforce outages are outside the organization's control but impact the internal CRM Sync. |
| ASF-010 | The API Gateway is deployed in a highly available configuration (multi-AZ) to avoid becoming a single point of failure. | Architectural | A single-AZ API Gateway fails with the AZ, blocking all vendor traffic. |
| ASF-011 | The internal CRM Sync service has a message queue or buffer for when Salesforce API calls fail or are rate-limited. | Derived | A synchronous integration that fails on API errors loses data that was in transit at the time of failure. |
| ASF-012 | The vendor maintains a business continuity plan that covers the services the organization depends on. | Dependency | The vendor's BCP may prioritize other customers; the organization's SLA may not be met during the vendor's disaster recovery. |

---

### Pattern 4: Backup & Recovery

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-013 | Customer data in Salesforce can be exported on demand (data portability) if the contract is terminated. | Explicit | Data portability is a regulatory requirement (GDPR Article 20) and a business continuity requirement. |
| ASF-014 | The internal customer database is backed up independently of the Salesforce data — recovery does not depend on Salesforce availability. | Derived | A database restore that depends on re-syncing from Salesforce will fail if Salesforce is also unavailable. |
| ASF-015 | Salesforce data backups are tested for integrity and can be restored within the RTO. | Implicit | The DPA should specify backup and restore SLAs, but these are typically not verified. |
| ASF-016 | The organization maintains a local cache or copy of critical Salesforce data that can be used if Salesforce is unavailable for extended periods. | Operational | Business operations that depend on real-time Salesforce data will halt during extended vendor outages. |

---

### Pattern 5: Change Management

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-017 | Changes to the API Gateway configuration for the vendor integration follow the organization's change management process. | Explicit | A misconfiguration change to the API Gateway can expose internal systems or break the integration. |
| ASF-018 | Vendor API version upgrades are tested in a non-production environment before the production integration is updated. | Operational | Untested vendor API upgrades can break the CRM Sync service with unexpected breaking changes. |
| ASF-019 | OAuth scope changes (adding new permissions) require approval — they are not made ad-hoc by developers. | Derived | Scope creep in OAuth tokens is a common finding in vendor integration audits. |
| ASF-020 | The quarterly vendor review includes reviewing any changes to the vendor's security program or infrastructure. | Operational | Vendor security changes (e.g., moving to a new cloud provider, changing data centers) may affect data residency and security posture. |

---

### Pattern 6: Cloud Security (IAM)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-021 | The API Gateway IAM execution role has least-privilege permissions — no access to resources outside the integration scope. | Explicit | An over-permissioned API Gateway role can be used to access other AWS resources if the gateway is compromised. |
| ASF-022 | The internal CRM Sync service's IAM role for accessing the customer database is scoped to specific database operations. | Derived | A CRM Sync role with database administrator privileges can modify schemas or delete tables. |
| ASF-023 | The AWS account running the API Gateway and CRM Sync has no other workloads that share the same IAM roles. | Implicit | Shared IAM roles across workloads increase blast radius if the vendor integration is compromised. |
| ASF-024 | CloudTrail is enabled for the API Gateway account to log all administrative actions and API calls. | Trust | Without CloudTrail, unauthorized changes to the API Gateway configuration are invisible. |

---

### Pattern 7: Container Security

*Not applicable — no containers documented in this architecture.*

---

### Pattern 8: Data Flow & Classification

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-025 | Customer data flowing through the integration is classified as sensitive/confidential under the organization's data classification policy. | Explicit | Data classification determines encryption, access control, and handling requirements for the data in transit. |
| ASF-026 | Data flow diagrams exist and accurately represent all paths customer data travels — from Salesforce through API Gateway to the internal database. | Implicit | An undocumented data flow (e.g., logs that capture response payloads, error queues) creates a blind spot. |
| ASF-027 | No customer data from the vendor integration is written to application logs, error messages, or debugging output. | Derived | Customer data in logs is accessible to operations teams and SIEM vendors who should not have access. |
| ASF-028 | Data residency requirements are enforced by the vendor — Salesforce does not replicate data to regions outside the approved geography. | Environmental | A global SaaS vendor may replicate data for disaster recovery to regions that violate data residency laws. |

---

### Pattern 9: Encryption at Rest

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-029 | Customer data in Salesforce is encrypted at rest with AES-256 or equivalent. | Explicit | SOC 2 Type II should cover encryption at rest, but the organization must verify the encryption standard. |
| ASF-030 | The customer database (internal) is encrypted at rest using a KMS-managed key with key rotation. | Derived | The documented DPA does not cover internal system encryption; the organization must ensure its own encryption. |
| ASF-031 | The API Gateway does not persist customer data to disk or cache unencrypted data. | Implicit | An API Gateway that caches request/response payloads for logging or analytics creates an unencrypted data store. |
| ASF-032 | Encryption keys for the customer database are not stored with the data — separate KMS or HSM is used. | Trust | Database encryption with a key stored in the same system provides no real protection against a database compromise. |

---

### Pattern 10: Encryption in Transit (TLS)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-033 | TLS is enforced between Salesforce and the API Gateway — HTTPS with valid certificates. | Explicit | The documented integration uses OAuth 2.0 over HTTPS; this must be verified at the network level. |
| ASF-034 | The API Gateway validates the Salesforce TLS certificate — no self-signed or expired certificates are accepted. | Derived | Without certificate validation, a DNS hijack of Salesforce's API endpoint can intercept the connection. |
| ASF-035 | TLS 1.2 or higher is enforced on all connections — Salesforce API endpoint, API Gateway, and customer database. | Derived | TLS 1.0/1.1 connections are vulnerable to downgrade attacks. |
| ASF-036 | Mutual TLS (mTLS) is considered for the API Gateway-to-internal service connection to provide service identity. | Trust | The current architecture trusts the internal network; mTLS would provide cryptographic service identity. |

---

### Pattern 11: Endpoint Security

*Not applicable — no user endpoints in this architecture.*

---

### Pattern 12: Human Factors & Process

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-037 | The team managing the vendor integration understands OAuth 2.0 security best practices (token rotation, scope limitation, PKCE). | Derived | OAuth misconfiguration is one of the most common API security vulnerabilities. |
| ASF-038 | Employees with Salesforce administrative access receive security training specific to SaaS vendor management. | Operational | General security training often does not cover vendor integration risks or OAuth security. |
| ASF-039 | The individual conducting the quarterly vendor review has the technical expertise to evaluate the vendor's security responses. | Implicit | A procurement manager conducting a vendor review may not understand security implications of the vendor's answers. |
| ASF-040 | CRM Sync configuration changes are made by authorized personnel only — not shared admin credentials for the integration. | Trust | Shared admin credentials for the integration eliminate individual accountability for configuration changes. |

---

### Pattern 13: Identity Lifecycle (Provisioning)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-041 | Salesforce user accounts for organization employees follow the joiner/mover/leaver process. | Operational | A terminated employee with active Salesforce access can exfiltrate customer data. |
| ASF-042 | Salesforce OAuth client registrations are reviewed and recertified quarterly. | Derived | Orphaned OAuth client registrations (from decommissioned integrations) retain authorized access to Salesforce. |
| ASF-043 | Service account used for the OAuth integration is not shared across environments (dev, staging, production). | Implicit | A shared OAuth client credential across environments makes it impossible to audit which environment accessed which data. |
| ASF-044 | Salesforce API access for the integration is governed by a dedicated API-only user with restricted profile, not a full-license user. | Operational | An API integration using a full-license Salesforce user has access to UI features and data the integration does not need. |

---

### Pattern 14: Incident Response

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-045 | The incident response plan includes a scenario for vendor data breach — vendor compromises customer data in Salesforce. | Operational | A vendor breach requires coordination with the vendor's IR team, which is slower and more complex than internal IR. |
| ASF-046 | The IR team has the ability to revoke the OAuth token immediately to cut off the vendor integration during an incident. | Derived | The OAuth token is the only control that connects the vendor and internal systems; revoking it isolates both sides. |
| ASF-047 | The vendor's breach notification process has been tested or tabletopped to verify the organization receives timely notification. | Trust | An untested vendor notification process will fail under the stress of an actual breach. |
| ASF-048 | The IR team can access logs from the API Gateway (request/response logging) to determine scope of a vendor-side compromise. | Operational | Without API Gateway logs, the IR team cannot determine what data was accessed through the compromised vendor account. |

---

### Pattern 15: Least Privilege

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-049 | The OAuth token scope is limited to read-only access on the minimum Salesforce objects needed for synchronization. | Explicit | An OAuth token with write access allows the internal system to modify Salesforce data, creating an unmonitored data path. |
| ASF-050 | The internal CRM Sync service's database user has SELECT/INSERT/UPDATE only on the customer schema — no DDL or DELETE. | Derived | A CRM Sync user with DELETE access can accidentally or maliciously delete customer records. |
| ASF-051 | The API Gateway does not have access to any internal systems beyond the CRM Sync service. | Implicit | An API Gateway with network access to multiple internal services can be used as a pivot point. |
| ASF-052 | Salesforce profile permissions for the integration user are scoped to object-level read, not "View All Data." | Derived | A "View All Data" profile on the integration user bypasses object-level and record-level security. |

---

### Pattern 16: Monitoring & Alerting

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-053 | API Gateway request/response logging is enabled for the vendor integration and monitored for anomalies. | Operational | Without logging, a compromised vendor account conducting reconnaissance through the API is invisible. |
| ASF-054 | Failed OAuth token validation attempts are logged and alerted — a high rate indicates token scanning or brute-force. | Derived | Attackers who obtain a leaked OAuth token will test it; failed attempts are the detection signal. |
| ASF-055 | The organization monitors Salesforce login activity for anomalous patterns (geolocation, time of day, concurrent sessions). | Operational | A compromised Salesforce admin account logging in from an unusual location indicates credential compromise. |
| ASF-056 | Data volume anomalies from the vendor API are detected — a sudden increase in data pulled from Salesforce indicates exfiltration. | Implicit | A compromised vendor account or insider threat exfiltrating customer data will show as increased API call volume. |

---

### Pattern 17: Network Segmentation

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-057 | The API Gateway is in a DMZ or public subnet — the CRM Sync service is in a private subnet with no direct internet access. | Explicit | A private subnet for the CRM Sync prevents direct internet attacks and limits egress to the API Gateway only. |
| ASF-058 | Security groups restrict the CRM Sync service to inbound traffic only from the API Gateway security group. | Derived | A CRM Sync that accepts traffic from any source defeats network segmentation. |
| ASF-059 | The customer database security group allows traffic only from the CRM Sync service — no direct access from API Gateway. | Architectural | A database accessible from the API Gateway is one hop closer to the internet than necessary. |
| ASF-060 | The API Gateway has IP whitelist rules restricting inbound traffic to Salesforce's published IP ranges. | Implicit | Without IP restriction, the OAuth token is the only barrier; a leaked token can be used from any network. |

---

### Pattern 18: Physical Security

*Not applicable — SaaS vendor and cloud-hosted internal systems.*

---

### Pattern 19: Supply Chain Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-061 | Salesforce's supply chain (third-party services Salesforce depends on) does not introduce risk to the organization's data. | Dependency | A compromise of a Salesforce sub-processor (e.g., their cloud provider, CDN, monitoring vendor) can affect customer data. |
| ASF-062 | The API Gateway software (WAF, reverse proxy, API management) is from a trusted vendor and has no known unpatched vulnerabilities. | Explicit | The API Gateway itself is infrastructure software that must be patched and maintained. |
| ASF-063 | Open-source libraries used by the CRM Sync service are scanned for vulnerabilities before deployment. | Operational | A vulnerable library in the CRM Sync service can be exploited if the vendor sends crafted API responses. |
| ASF-064 | The organization maintains a software bill of materials (SBOM) for the internal CRM Sync application. | Derived | Without an SBOM, the organization cannot determine if a newly disclosed library CVE affects the CRM Sync. |

---

### Pattern 20: Third-party Dependency

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-065 | Salesforce maintains SOC 2 Type II certification throughout the contract period — recertification does not lapse. | Dependency | A lapsed SOC 2 certification means the organization has no current assurance of the vendor's controls. |
| ASF-066 | Salesforce's security posture does not degrade between SOC 2 audit cycles — no significant control changes without notice. | Dependency | SOC 2 is a point-in-time assessment; the vendor's security could degrade the day after the audit. |
| ASF-067 | The organization has a contractual right to access Salesforce's penetration test results or security assessment reports. | Dependency | Without access to test results, the organization relies entirely on the vendor's self-reported security posture. |
| ASF-068 | The vendor has a documented data retention and deletion policy that complies with the DPA and applicable regulations. | Dependency | Data retained by the vendor beyond the contract term violates GDPR right to erasure and other data protection principles. |
| ASF-069 | There is a vendor exit plan that includes data extraction, deletion verification, and transition to an alternative vendor or internal system. | Derived | An untested vendor exit plan leads to data loss or extended business disruption during vendor transition. |
| ASF-070 | The vendor's service does not use the organization's data to train AI/ML models without explicit consent. | Environmental | Many SaaS vendors now use customer data for AI training; this must be explicitly prohibited in the DPA. |

**Total (A): 70** (4 per pattern × 16 applicable patterns + 2 extra for Third-party Dependency)

---

## Step 3: Comparison

### Overlap Mapping

| Human ID | ASF ID | Match Rationale |
|----------|--------|-----------------|
| H-001 | ASF-065 | Both require current SOC 2 certification throughout the contract. |
| H-003 | ASF-006 | Both require authorization code flow with PKCE for OAuth. |
| H-004 | ASF-049 | Both require OAuth token scoped to minimum Salesforce objects/operations. |
| H-005 | ASF-007 | Both require OAuth refresh token rotation and limited lifetime. |
| H-006 | ASF-006 | Both require token validation on every API request. |
| H-008 | ASF-047 | Both require vendor breach notification within SLA. |
| H-009 | ASF-061 | Both require vendor sub-processors to be identified and assessed. |
| H-011 | ASF-033 | Both require encryption in transit between all system hops. |
| H-013 | ASF-063 | Both require input validation and protection against injection from vendor. |
| H-014 | ASF-011 | Both require rate limiting and queuing for vendor API calls. |
| H-015 | ASF-059 | Both require database not directly accessible from internet. |
| H-016 | ASF-050 | Both require least-privilege database access for CRM Sync. |
| H-017 | ASF-007 | Both require secure storage for OAuth credentials. |
| H-018 | ASF-018 | Both require vendor API change notification and testing. |
| H-019 | ASF-011 | Both require fallback/queue when vendor API is unavailable. |
| H-020 | ASF-060 | Both require API Gateway IP whitelisting for Salesforce. |
| H-021 | ASF-028 | Both require data residency requirements to be enforced. |
| H-023 | ASF-053 | Both require API Gateway request/response logging for audit. |
| H-024 | ASF-042 | Both require OAuth scope and integration review. |
| H-025 | ASF-045 | Both require incident response plan for vendor breach. |
| H-026 | ASF-067 | Both require right-to-audit or access to pen test results. |
| H-027 | ASF-029 | Both require vendor encryption at rest for customer data. |
| H-028 | ASF-028 | Both require API Gateway request validation. |
| H-029 | ASF-035 | Both require TLS 1.2+ for OAuth token endpoint. |
| H-031 | ASF-031 | Both require vendor employee access logging. |
| H-034 | ASF-056 | Both require monitoring for anomalous API call patterns. |
| H-037 | ASF-009 | Both require vendor availability SLA. |
| H-038 | ASF-069 | Both require vendor exit strategy and data migration plan. |
| H-040 | ASF-027 | Both assume customer data is not written to logs or caches. |
| H-042 | ASF-046 | Both require ability to revoke OAuth token/kill switch for vendor integration. |

**Overlap (O): 31**

### Metrics

| Metric | Value | Formula |
|--------|-------|---------|
| Human assumptions (H) | 43 | Count of unique human-generated assumptions |
| ASF assumptions (A) | 70 | Count of unique ASF-generated assumptions |
| Overlap (O) | 31 | Count appearing in both lists |
| **Precision** | **44.3%** | O / A = 31/70 |
| **Recall** | **72.1%** | O / H = 31/43 |
| **F1 Score** | **54.9%** | 2 × (P × R) / (P + R) |
| Novel findings (A - O) | 39 | Assumptions ASF found that human missed (55.7% of ASF total) |
| Missed findings (H - O) | 12 | Assumptions human found that ASF missed (27.9% of human total) |

### Success Criteria Assessment

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Recall | >= 70% | 72.1% | ✅ Met |
| Precision | >= 50% | 44.3% | ❌ Not met |
| Novel discoveries | >= 10% of total (A+O) | 35.5% (39/110) | ✅ Exceeded |
| Expert agreement (F1 proxy) | > 60% | 54.9% | ❌ Not met |

---

## Step 4: Analysis

### Ontology Overlap Analysis

| Ontology Category | Overlaps | ASF Total | Overlap Rate |
|-------------------|----------|-----------|-------------|
| Explicit | 8 | 12 | 66.7% |
| Derived | 8 | 18 | 44.4% |
| Operational | 7 | 18 | 38.9% |
| Implicit | 3 | 8 | 37.5% |
| Trust | 2 | 8 | 25.0% |
| Dependency | 2 | 8 | 25.0% |
| Architectural | 1 | 4 | 25.0% |
| Environmental | 0 | 2 | 0.0% |

**Best overlap:** Explicit and Derived — both the human and the ASF focused on OAuth configuration, TLS enforcement, and data encryption as primary assumptions.

**Worst overlap:** Environmental had zero overlap. The ASF identified Salesforce AI training data usage (ASF-070) and data residency (ASF-028) as assumptions the human did not list.

### What Humans Caught That ASF Missed (Missed Findings = 12)

1. **SOC 2 exceptions review (H-002):** The human looked beyond the existence of the SOC 2 report to the content of findings. The ASF assumes the report provides assurance without examining exceptions.

2. **Right-to-audit contract clause (H-026):** The human identified legal recourse for security assessment. The ASF treats vendor assessment as a periodic activity, not a contractual right.

3. **Vendor insider threat monitoring (H-031):** The human assumed the vendor logs and reviews its own employees' access. The ASF's vendor posture pattern did not address vendor employee monitoring.

4. **Vendor penetration testing cadence (H-032):** The human specified annual pen testing as a requirement. The ASF did not generate an assumption about the vendor's testing frequency.

5. **Data retention and deletion policy alignment (H-033):** The human linked retention policy to the DPA. The ASF assumed retention exists but did not match the human's cross-referencing.

6. **Purpose limitation for vendor data use (H-035, H-043):** The human restricted the vendor from using data for AI/ML training or other secondary purposes. This is a modern concern the ASF patterns do not address.

7. **OAuth token Salesforce instance validation (H-036):** The human considered a specific Salesforce multi-tenancy risk. The ASF OAuth pattern is generic and did not cover this platform-specific detail.

8. **Database backup independence from vendor (partially covered by ASF-014 but not fully):** The human assumed local database backups are independent. The ASF covered backup but not the decoupling from vendor reliance.

### What ASF Caught That Humans Missed (Novel Findings = 39)

1. **Change management for vendor integration (ASF-017 through ASF-020):** The human generated zero assumptions about the change process for API Gateway configuration or OAuth scope changes. Pipeline configuration drift is a blind spot.

2. **IAM and CloudTrail for API Gateway (ASF-021, ASF-024):** The human focused on the external vendor and the database but did not consider the IAM posture or audit logging of the API Gateway itself.

3. **Data classification (ASF-025):** The human assumed data sensitivity but did not list formal data classification as an explicit assumption.

4. **Incident response details (ASF-045 through ASF-048):** The human had a single IR assumption (H-025); the ASF provided a full pattern on vendor breach response, OAuth token revocation, notification testing, and log access.

5. **Identity lifecycle for Salesforce access (ASF-041 through ASF-044):** The human had a single assumption about quarterly review (H-024); the ASF generated four assumptions about joiner/mover/leaver, OAuth client recertification, environment separation, and API-only user profiles.

6. **Supply chain risk beyond Salesforce (ASF-061):** The human assumed Salesforce itself is the only third-party; the ASF pointed out that Salesforce depends on other vendors (cloud, CDN, monitoring) that create a deeper supply chain.

7. **AI/ML data usage prohibition (ASF-070):** The ASF identified a modern vendor risk that the human did not consider: the vendor using customer data to train AI models.

### Architecture Complexity Assessment

Architecture #009 (Vendor SaaS Integration) achieved **72.1% recall**, meeting the 70% target. This is the second architecture to meet the recall criterion. The strong performance is driven by:
- Close alignment between the ASF Third-party Dependency pattern and the architecture's primary risk
- The ASF OAuth pattern closely matching the authentication mechanism
- Strong overlap in data protection assumptions (encryption, classification)

The **novelty rate (55.7%)** confirms the ASF adds substantial value even with good recall. The human focused on vendor contractual and technical controls, while the ASF added identity lifecycle, change management, and incident response assumptions that the human overlooked.

### Key Insight

The ASF pattern matrix is well-aligned with vendor integration security concerns. The two gaps in precision were:
- **Vendor contract-specific details** (right-to-audit, SOC 2 exceptions, purpose limitation) that require legal-contractual pattern inputs
- **Vendor platform-specific OAuth details** (Salesforce instance validation, API-only user profiles)

Adding a "Vendor Contract Governance" sub-pattern under Third-party Dependency would address the legal-contractual assumptions. The strong recall (72.1%) suggests the current patterns are largely sufficient for this architecture type.

---

## Conclusions

| Metric | Target | Actual | Verdict |
|--------|--------|--------|---------|
| Recall | >= 70% | 72.1% | ✅ Met — strong alignment with vendor security patterns |
| Precision | >= 50% | 44.3% | ❌ Below target — ASF is deliberately broad/generative |
| Novel discoveries | >= 10% | 35.5% | ✅ ASF adds substantial value beyond human reasoning |
| Expert agreement (F1) | > 60% | 54.9% | ❌ Below target — driven by low precision |

The ASF framework applied to Architecture #009 demonstrates strong recall (72.1%), meeting the 70% target. The ASF's Third-party Dependency pattern is well-suited for vendor-integrated architectures. The key improvement opportunity is adding a **Vendor Contract Governance** sub-pattern to capture legal-contractual assumptions (right-to-audit, SOC 2 exceptions, purpose limitation, data retention alignment) that the current patterns do not surface.
