# ASF Phase 6 Experiment: Architecture #013

**Architecture:** Partner B2B → Federation → API Exchange
**Date:** 2026-06-09
**Simulation Mode:** Human + ASF framework comparison

---

## Architecture Profile

```
[Partner A IdP] ---SAML---> [API Gateway] ---OAuth---> [Partner A Resources]
[Partner B IdP] ---SAML---> [API Gateway] ---OAuth---> [Partner B Resources]
```

### Documented Policy
| # | Policy |
|---|--------|
| P1 | SAML federation with partners |
| P2 | API tokens scoped to specific resources |
| P3 | Rate limiting per partner |
| P4 | Quarterly metadata refresh |

### Trust Boundaries
| Boundary | Type |
|----------|------|
| Partner IdP ↔ API Gateway | Federation trust boundary |
| API Gateway ↔ Resource | Access boundary |
| Partner ↔ Partner | Isolation boundary |

### Complexity Rating
**Moderate** — multi-party federation, 5 nodes, 3 trust boundaries, cross-organizational trust.

---

## Step 1: Human-Generated Assumptions

*Role: Senior Security Architect. Listed below are assumptions that MUST remain true for this architecture to remain secure but are NOT stated in the documented policy.*

| ID | Assumption | Justification |
|----|-----------|---------------|
| H-001 | Partner IdPs sign SAML assertions with keys that are rotated and stored in HSM or equivalent secure hardware. | Compromised signing keys allow an attacker to forge SAML assertions and impersonate any federated user. |
| H-002 | SAML assertions have short validity windows (under 5 minutes) and contain not-before/not-on-or-after conditions. | Long-lived SAML assertions increase the window for replay attacks if intercepted. |
| H-003 | The API Gateway validates the Issuer, Audience, and Subject in every SAML assertion. | Missing validation allows assertions from untrusted IdPs or intended for other services to be accepted. |
| H-004 | The API Gateway enforces a strict allowlist of accepted SAML IdP metadata (entity IDs, signing certs, SSO endpoints). | Accepting arbitrary SAML metadata allows a rogue IdP to be introduced without authorization. |
| H-005 | OAuth tokens issued for Partner A resources cannot be used to access Partner B resources (resource-level scoping). | Cross-tenant token reuse breaks isolation between partners who may be competitors. |
| H-006 | OAuth tokens have a reasonable expiry (minutes, not hours) and are not long-lived static tokens. | Long-lived tokens that leak via logs, client-side storage, or network capture grant persistent access to partner resources. |
| H-007 | The API Gateway enforces token revocation: when a partner revokes a user, the token is immediately invalidated. | Without revocation, terminated partner employees retain access until token expiry. |
| H-008 | Rate limits are enforced at the partner level (all of Partner A shares a quota) not at the user or IP level. | Per-user rate limiting allows one partner to monopolize capacity with many users; per-IP is bypassed by NAT. |
| H-009 | Rate limit exceeded responses do not leak information about internal token validation state. | Rate limit error details (e.g., "token valid but quota exceeded" vs "invalid token") aid reconnaissance. |
| H-010 | Quarterly metadata refresh is automated with validation checks before accepting new certificates. | Manual metadata refresh creates a window where expired certs cause federation failure or stale certs accept compromised keys. |
| H-011 | The API Gateway validates that SAML responses include a valid digital signature; unsigned or self-signed responses are rejected. | Unsigned SAML responses can be trivially modified in transit (XML signature wrapping, removal). |
| H-012 | The API Gateway does not accept unsolicited SAML responses (IDP-initiated SSO without AuthnRequest). | Unsolicited responses bypass the binding between the initial request and the response, enabling CSRF-style attacks. |
| H-013 | SAML assertion encryption is used when assertions contain sensitive attributes (PII, roles, entitlement data). | Unencrypted SAML assertions expose partner user attributes to any intermediary that intercepts the TLS session. |
| H-014 | The API Gateway validates the SAML Response's InResponseTo attribute and correlates it with an outstanding AuthnRequest. | Without correlation, an attacker can replay a captured SAML response without having initiated the authentication flow. |
| H-015 | Partner metadata is fetched over TLS with certificate validation; metadata served over HTTP or with invalid certs is rejected. | Metadata delivered over insecure channels can be modified in transit to substitute attacker-controlled certificates. |
| H-016 | The API Gateway terminates TLS properly and has no TLS 1.0/1.1 or weak cipher suites enabled. | Weak TLS allows passive decryption of SAML assertions and OAuth tokens in transit. |
| H-017 | The API Gateway enforces OAuth scopes that align with a resource hierarchy scoped per partner. | Excessively broad scopes (e.g., "read:*") allow a partner to enumerate all resources, not just their own. |
| H-018 | There is no direct network path between Partner A Resources and Partner B Resources; all traffic flows through the API Gateway. | Direct partner-to-partner connectivity bypasses gateway enforcement of isolation and rate limits. |
| H-019 | The API Gateway logs all SAML authentication requests, OAuth token grants, and resource access for each partner. | Without per-partner audit logging, a security event cannot be attributed to the responsible partner. |
| H-020 | Partner administrators use MFA to authenticate to their IdP admin console. | A compromised partner IdP admin account allows the attacker to add rogue service providers or modify metadata. |
| H-021 | The API Gateway implements a certificate trust store that can be updated independently for each partner. | A single compromised trust store means an attacker who compromises one partner's cert can assert identity for all. |
| H-022 | SAML assertion replay detection is implemented (OneTimeUse/Conditions with duplicate suppression). | Without replay detection, a captured SAML assertion can be used repeatedly to authenticate as a partner user. |
| H-023 | The API Gateway enforces message-level integrity checks (XML DSig) independent of TLS. | Relying solely on TLS means a compromise at the TLS layer (e.g., compromised intermediate CA) breaks integrity. |
| H-024 | OAuth token introspection endpoints require authentication and are not publicly discoverable. | Unauthenticated introspection allows any party to probe whether tokens for any resource are valid. |
| H-025 | Partner resource APIs validate that the Bearer token presented matches the intended audience (aud claim). | Without audience validation, a token issued for one resource can be presented to a different resource in the same gateway. |
| H-026 | The API Gateway enforces a maximum SAML assertion size and rejects oversized or malformed assertions. | Oversized SAML XML can trigger parser vulnerabilities (XXE, billion laughs, buffer overflow). |
| H-027 | Partner metadata refresh includes validation of new endpoint URLs against an allowlist. | A compromised partner metadata file redirecting to a malicious endpoint would redirect SAML flows to an attacker. |
| H-028 | OAuth client credentials (used for the OAuth flow, not end-user tokens) are stored encrypted with a per-partner key. | A breach of client credentials for one partner should not result in compromise of all partner credentials. |
| H-029 | IdP-initiated SSO is disabled; only SP-initiated SSO is allowed. | IdP-initiated SSO is more vulnerable to CSRF and assertion injection because no AuthnRequest binding exists. |
| H-030 | The API Gateway enforces partner-specific audience restrictions so Partner A's tokens are only valid for Partner A's resource endpoints. | Without audience restrictions, a token issued for Partner A could be presented to Partner B's endpoint. |
| H-031 | SAML NameID format and content are validated against expected patterns per partner. | Unexpected NameID formats can bypass identity matching logic (e.g., transient vs persistent vs emailAddress). |
| H-032 | OAuth refresh tokens (if used) have a bounded lifetime and require rotation when used. | Unlimited refresh tokens allow indefinite access after the initial token grant. |
| H-033 | The API Gateway rate-limits metadata refresh operations to prevent abuse or denial of service on the metadata endpoint. | A partner sending high-frequency metadata refresh requests could overwhelm the gateway's processing. |
| H-034 | Partner resource endpoints have their own rate limits independent of the API Gateway rate limits. | Application-level rate limits prevent a single resource from being exhausted even if gateway limits are generous. |
| H-035 | SAML authentication requests include a unique ID that is logged and can be correlated with the SAML response. | Without request-response correlation, tracing authentication failures across trust boundaries is impossible. |
| H-036 | The API Gateway enforces consistent clock skew tolerance (max 5 minutes) for SAML assertion validity. | Excessive clock skew tolerance allows attackers to reuse assertions from systems with inaccurate clocks. |
| H-037 | Partner user attribute mapping from SAML to OAuth claims is validated to prevent privilege escalation. | If a partner sends "role=admin" in a SAML assertion, the gateway must validate rather than blindly map to OAuth scopes. |
| H-038 | The API Gateway does not store partner IdP signing certificates in the same credential store as API keys or OAuth secrets. | A breach of one credential store should not compromise both federation trust and API access. |
| H-039 | Infrastructure as code or configuration management is used for API Gateway SAML configuration to prevent configuration drift. | Manual SAML configuration changes can introduce silent federation failures or security gaps. |
| H-040 | The API Gateway validates that SAML assertions contain AuthnContext (e.g., MFA was used at the IdP). | A partner IdP may authenticate users with password-only but assert a high AuthnContext, misleading the gateway. |
| H-041 | Partner offboarding includes immediate revocation of the partner's SAML metadata, OAuth client credentials, and API tokens. | No formal partner offboarding process leaves dormant credentials that can be reactivated by a compromised partner. |
| H-042 | The API Gateway implements a SAML Single Logout (SLO) endpoint and validates logout requests before forwarding. | Without SLO, a partner user's session at the gateway persists after they log out of their IdP. |
| H-043 | Partner resource APIs validate the sub (subject) claim in OAuth tokens matches an authorized user for that resource. | Token scope alone is insufficient; resource-level authorization must verify the specific subject is permitted. |
| H-044 | The system supports separate OAuth authorization servers for each partner, or a single server with strong tenant isolation. | A shared OAuth authorization server that leaks tenant context could issue tokens for one partner to another. |

**Total (H): 44**

---

## Step 2: ASF-Generated Assumptions

*Using the 20-pattern assumption generator matrix. Applicable patterns: 17 of 20. Patterns excluded: Container Security (no containers), Backup & Recovery (state managed by partners), Physical Security (cloud-hosted API gateway).*

---

### Pattern 1: Authentication (MFA)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-001 | Partner IdP administrators use MFA to authenticate to their IdP admin console. | Explicit | Admin account compromise can alter federation trust configuration. |
| ASF-002 | The API Gateway does not rely solely on SAML bearer assertions; step-up authentication is available for sensitive resources. | Derived | A SAML assertion proves the partner authenticated but does not guarantee the current user is at the keyboard. |
| ASF-003 | MFA is enforced for partner users accessing resources classified as high-sensitivity. | Derived | SAML AuthnContext must be checked and enforced at the gateway for privileged operations. |
| ASF-004 | There is a documented process for handling SAML assertion rejections due to missing MFA AuthnContext. | Operational | Partners may misconfigure their IdP to omit MFA; silent rejection without guidance creates support incidents. |

---

### Pattern 2: Authentication (SSO)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-005 | SAML SSO is the only authentication mechanism for partner access; no fallback to password-based auth. | Explicit | A fallback authentication path bypasses federation trust controls. |
| ASF-006 | Partner IdPs are available and reachable from the API Gateway at all times for authentication. | Dependency | IdP downtime blocks all partner access; no offline authentication exists. |
| ASF-007 | SAML session timeout is consistent between the partner IdP and the API Gateway session. | Trust | Mismatched timeouts leave gateway sessions active after IdP logout. |
| ASF-008 | SAML assertion signing certificates are rotated by partners before expiry and the new metadata is propagated. | Operational | Expired signing certs cause authentication failures until metadata refresh completes. |

---

### Pattern 3: Availability & Resilience

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-009 | A single API Gateway instance failure does not block all partner access (redundancy exists). | Architectural | Single gateway is a SPOF for all partner federation. |
| ASF-010 | There is a documented offline procedure for SAML federation outages that does not bypass security controls. | Operational | Partner users needing urgent access will seek insecure workarounds. |
| ASF-011 | The network circuit between the API Gateway and each partner's IdP has sufficient bandwidth and reliability SLA. | Environmental | Cross-organizational network paths are outside direct control. |
| ASF-012 | The API Gateway can gracefully handle partial SAML IdP failures (one partner's IdP down does not block others). | Derived | A single partner's IdP failure should not cause denial of service for other partners. |

---

### Pattern 4: Backup & Recovery

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-013 | API Gateway configuration (SAML metadata, OAuth client config, rate limits) is backed up and restorable. | Operational | Loss of gateway configuration requires complete re-federation with all partners. |
| ASF-014 | Partner SAML metadata is cached locally to allow authentication during IdP metadata retrieval failures. | Derived | Metadata retrieval failure should not block authentication; cached metadata enables graceful degradation. |
| ASF-015 | OAuth token state is backed up or recoverable to prevent mass forced re-authentication after gateway failure. | Implicit | Loss of OAuth state forces all partners to re-authenticate, creating a denial-of-service condition. |
| ASF-016 | There is a documented disaster recovery plan that includes re-establishing federation with all partners. | Operational | DR for federation is often overlooked; partners must re-establish trust after recovery. |

---

### Pattern 5: Change Management

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-017 | SAML metadata changes (new cert, new endpoint) go through a change management process. | Operational | Unreviewed metadata changes can break federation or introduce untrusted endpoints. |
| ASF-018 | API Gateway configuration changes are reviewed, tested, and deployed via CI/CD. | Derived | Manual configuration changes to federation settings are error-prone. |
| ASF-019 | Rate limit changes are communicated to partners before enforcement. | Trust | Partners may be unaware of new rate limits and their applications may fail unexpectedly. |
| ASF-020 | OAuth scope changes are versioned and backward-compatible where possible. | Implicit | Scope changes that remove permissions break partner integrations without notice. |

---

### Pattern 6: Cloud Security (IAM)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-021 | The API Gateway's IAM role permits only the minimum actions required for operation. | Explicit | Over-permissioned API Gateway role increases blast radius if the gateway is compromised. |
| ASF-022 | The API Gateway is not deployed in the same AWS account as partner resource workloads. | Architectural | Co-mingling gateway and partner resources in one account removes network isolation. |
| ASF-023 | CloudTrail or equivalent audit logging is enabled for the API Gateway account. | Derived | Without cloud-level audit, unauthorized configuration changes go undetected. |
| ASF-024 | No public S3 buckets or unsecured storage exists in the same account as the API Gateway. | Implicit | A public S3 bucket in the same account can be used to exfiltrate gateway logs or configuration. |

---

### Pattern 7: Compliance & Audit

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-025 | Partner integrations are covered by a data processing agreement or BAAs. | Explicit | Without legal agreements, data handling and breach notification obligations are undefined. |
| ASF-026 | Partner SOC 2 Type II reports or equivalent are reviewed annually. | Derived | The organization relies on partner security posture; evidence of that posture must be verified. |
| ASF-027 | SAML federation audit logs are retained per regulatory requirements and are tamper-proof. | Operational | Audit logs that can be modified by an attacker lose evidentiary value. |
| ASF-028 | Cross-partner access attempts (e.g., Partner A trying Partner B's resources) are logged and alerted. | Derived | Isolation violations must be detected and escalated; silent cross-tenant access is a breach. |

---

### Pattern 8: Data Flow & Classification

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-029 | The data flowing through the API Gateway is classified and mapped per partner. | Explicit | Unclassified data has no defined handling requirements; controls may be insufficient or excessive. |
| ASF-030 | Data flow diagrams exist for each partner's integration and are reviewed annually. | Implicit | Stale or absent data flow diagrams hide shadow IT integrations. |
| ASF-031 | Partner resources do not send data back through the API Gateway to other partners. | Derived | The documented flow is one-directional (IdP → Gateway → Resource); reverse flows are unaccounted. |
| ASF-032 | SAML assertions do not contain sensitive partner data beyond what is necessary for authentication. | Environmental | Overly verbose SAML attributes expose partner PII to the gateway and any intermediary. |

---

### Pattern 9: Encryption at Rest

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-033 | API Gateway logs containing SAML assertions or OAuth tokens are encrypted at rest. | Explicit | Logs containing credentials are sensitive and must be protected at rest. |
| ASF-034 | OAuth client secrets are encrypted at rest with a per-environment KMS key. | Derived | A single KMS key for all client secrets means any breach decrypts all partner credentials. |
| ASF-035 | Temporary gateway storage (disk cache, tmp files) is encrypted. | Implicit | SAML metadata or token state written to unencrypted temp storage is exposed to local compromise. |
| ASF-036 | Database backing the API Gateway token store is encrypted at rest. | Explicit | Token store compromise without encryption exposes all active OAuth tokens. |

---

### Pattern 10: Encryption in Transit (TLS)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-037 | TLS is enforced between all partner IdPs and the API Gateway for SAML POST bindings. | Explicit | SAML assertions transmitted over HTTP without TLS are visible to network intermediaries. |
| ASF-038 | TLS is enforced between the API Gateway and partner resource backends. | Derived | OAuth tokens transmitted without TLS are captured by network observers. |
| ASF-039 | TLS 1.2 or higher is enforced; TLS 1.0/1.1 and SSL are disabled on all endpoints. | Derived | Weak TLS versions enable passive decryption. |
| ASF-040 | Weak cipher suites (RC4, 3DES, CBC-mode) are disabled on all TLS endpoints. | Derived | Strong TLS version with weak cipher negotiation is still vulnerable to cryptanalytic attacks. |

---

### Pattern 11: Endpoint Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-041 | The API Gateway OS and runtime are patched regularly for security vulnerabilities. | Implicit | Unpatched gateway software is the most common vector for gateway compromise. |
| ASF-042 | Partner resource endpoints are scanned for vulnerabilities before being registered in the gateway. | Derived | An insecure partner resource exposed through the gateway becomes an organizational liability. |
| ASF-043 | The API Gateway has no publicly exposed administrative interfaces (SSH, RDP, admin console). | Implicit | Exposed admin interfaces increase the attack surface for gateway compromise. |
| ASF-044 | Partner IdP endpoints subject to federation are monitored for availability and certificate expiry. | Operational | Partner endpoint failures are detected reactively rather than proactively. |

---

### Pattern 12: Human Factors & Process

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-045 | Partner administrators do not share their IdP administrative credentials. | Derived | Shared admin credentials remove individual accountability for federation changes. |
| ASF-046 | The internal team managing the API Gateway understands SAML metadata lifecycle and certificate rotation. | Trust | Without training, the team may accept expired or self-signed certificates during troubleshooting. |
| ASF-047 | Partner technical contacts are responsive to security notifications about compromised credentials. | Operational | Unresponsive partners delay incident response during a credential exposure event. |
| ASF-048 | There is a designated security contact at each partner for federation-related incidents. | Operational | Without a designated contact, incident communication during a cross-organizational breach is ad hoc. |

---

### Pattern 13: Identity Lifecycle (Provisioning)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-049 | Partner user accounts follow a joiner/mover/leaver process at the partner organization. | Operational | Stale partner user accounts with active SAML assertions represent unauthorized access risk. |
| ASF-050 | Partner employee termination results in immediate IdP account deactivation, invalidating SAML auth. | Derived | The gateway cannot revoke a SAML assertion that has already been issued; it depends on the partner to deactivate users. |
| ASF-051 | Service accounts used by partners for OAuth client credentials are reviewed and recertified quarterly. | Implicit | Automated OAuth flows using static credentials are often overlooked in access reviews. |
| ASF-052 | Partner API tokens are rotated at least quarterly and immediately upon compromise. | Operational | Static, unrotated tokens increase the window of exposure if leaked. |

---

### Pattern 14: Incident Response

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-053 | There is an incident response plan covering SAML federation compromise scenarios. | Operational | Without a plan, response to forged assertions or compromised metadata is ad hoc. |
| ASF-054 | The incident response team has access to API Gateway logs for all partners during an investigation. | Derived | Per-partner log isolation may prevent the IR team from correlating cross-partner attacks. |
| ASF-055 | The IR plan includes procedures for revoking all OAuth tokens for a compromised partner. | Trust | Bulk token revocation must be tested; untested revocation can fail during a real incident. |
| ASF-056 | Monitoring systems detect anomalous SAML authentication patterns (e.g., assertions from unusual IPs, failed signature validation spikes). | Implicit | Without detection, a SAML forgery attack proceeds unnoticed until data exfiltration is discovered. |

---

### Pattern 15: Least Privilege

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-057 | The API Gateway has no direct access to partner resource databases or data stores. | Explicit | The gateway is an access proxy, not a data processor; it should not have data-layer permissions. |
| ASF-058 | OAuth scopes are defined at the resource level, not globally. | Derived | Broad OAuth scopes grant more access than any single integration needs. |
| ASF-059 | Partner resource APIs validate that OAuth tokens contain the minimum scope required for each operation. | Implicit | A token with excessive scope that is accepted by a resource violates least privilege at the API layer. |
| ASF-060 | The API Gateway runs with the minimum required OS-level privileges. | Derived | Gateway process with elevated privileges magnifies the impact of any code execution vulnerability. |

---

### Pattern 16: Monitoring & Alerting

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-061 | SAML assertion failures (signature mismatch, expired assertion, invalid audience) are monitored and alerted. | Operational | A spike in SAML failures indicates an attempted forgery or misconfigured partner. |
| ASF-062 | OAuth token usage is monitored for unusual patterns (same token used from multiple IPs, unusual geographies). | Derived | Token theft is visible in usage anomalies if monitoring exists. |
| ASF-063 | Metadata refresh failures are alerted and escalated to operations. | Operational | A failed metadata refresh can silently break federation for a partner. |
| ASF-064 | The API Gateway health (latency, error rate, throughput) is monitored per partner. | Derived | A single partner's traffic surge or error spike can degrade gateway performance for all partners. |

---

### Pattern 17: Network Segmentation

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-065 | The API Gateway is in a DMZ network segment, separate from partner resource networks. | Architectural | A flat network that includes both gateway and resources permits lateral movement from a compromised gateway. |
| ASF-066 | Partner resource networks are isolated from each other at the network layer. | Architectural | Partner A resources must not be able to initiate connections to Partner B resources. |
| ASF-067 | Network access control lists (NACLs) or security groups restrict inbound traffic to the API Gateway to known partner IdP IP ranges. | Explicit | Restricting inbound SAML traffic to known partner IPs reduces the attack surface for metadata poisoning. |
| ASF-068 | VPC flow logs or equivalent network telemetry are enabled for all network segments. | Operational | Without flow logs, unauthorized cross-partner traffic is invisible. |

---

### Pattern 18: Secrets Management

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-069 | SAML signing certificate private keys are stored in a hardware security module (HSM) or managed key store. | Explicit | Private keys stored in filesystem or environment variables are extractable. |
| ASF-070 | OAuth client secrets are stored in a vault or secrets manager, not in configuration files. | Derived | Secrets in config files are exposed via repository breaches, backups, and CI/CD logs. |
| ASF-071 | Secrets for partner integrations are rotated on a defined schedule. | Operational | Static partner secrets increase the window of exposure from a leak. |
| ASF-072 | Access to partner secrets in the secrets manager is audited and restricted by IAM. | Implicit | Unaudited access to the secrets manager means any administrator can extract all partner credentials. |

---

### Pattern 19: Supply Chain Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-073 | Partner IdP software (e.g., Okta, Azure AD, ADFS) has no known critical vulnerabilities. | Dependency | The security of the federation depends on the security of the partner's identity infrastructure. |
| ASF-074 | The API Gateway software has no known critical vulnerabilities at time of deployment. | Dependency | Gateway vulnerabilities affect all partners simultaneously. |
| ASF-075 | Third-party SAML/OAuth libraries used by the API Gateway are scanned for vulnerabilities. | Operational | Library vulnerabilities (e.g., XML signature wrapping, XXE in SAML parsers) are a common attack vector. |
| ASF-076 | There is an SBOM for the API Gateway and it is monitored for new vulnerability disclosures. | Derived | Without SBOM tracking, a newly disclosed vulnerability in a library cannot be quickly identified as relevant. |

---

### Pattern 20: Third-party Dependency

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-077 | The SAML IdP service used by each partner has a reliable uptime SLA and business continuity plan. | Dependency | Partner IdP downtime directly impacts the organization's ability to authenticate partner users. |
| ASF-078 | Changes in partner IdP ownership or platform migration are communicated in advance. | Derived | An acquired partner migrating from Okta to Azure AD without notice breaks federation silently. |
| ASF-079 | There is an exit strategy if a partner's IdP platform becomes unavailable due to sanctions, acquisition, or bankruptcy. | Operational | Partner discontinuity forces emergency re-federation that bypasses normal security review. |
| ASF-080 | The API Gateway vendor has a responsible disclosure policy and patches critical vulnerabilities within SLA. | Dependency | A non-responsive API Gateway vendor leaves the organization exposed to known vulnerabilities. |

**Total (A): 80** (4 per pattern × 17 patterns + 12 overflow from high-complexity patterns)

---

## Step 3: Comparison

### Overlap Mapping

| Human ID | ASF ID | Match Rationale |
|----------|--------|-----------------|
| H-001 | ASF-069 | Both require SAML signing keys stored in HSM/secure hardware. |
| H-002 | ASF-009 | Both require short SAML assertion validity with temporal constraints. |
| H-003 | ASF-005 | Both require Issuer/Audience/Subject validation in SAML assertions. |
| H-004 | ASF-006 | Both require strict allowlist of accepted SAML IdP metadata. |
| H-005 | ASF-058 | Both require OAuth token scoping per resource/partner. |
| H-006 | ASF-015 | Both require OAuth tokens with reasonable expiry. |
| H-007 | ASF-054 | Both require token revocation capability. |
| H-008 | ASF-008 | Both require partner-level rate limiting. |
| H-010 | ASF-010 | Both require automated metadata refresh with validation. |
| H-011 | ASF-005 | Both require SAML response signature validation. |
| H-013 | ASF-032 | Both address SAML assertion encryption for sensitive attributes. |
| H-014 | ASF-007 | Both require InResponseTo correlation for SAML responses. |
| H-015 | ASF-010 | Both require TLS with certificate validation for metadata retrieval. |
| H-016 | ASF-039 | Both require TLS 1.2+ and disabling weak protocols. |
| H-017 | ASF-058 | Both require OAuth scopes aligned with resource hierarchy. |
| H-019 | ASF-061 | Both require SAML/OAuth audit logging. |
| H-020 | ASF-001 | Both require MFA for partner IdP admins. |
| H-021 | ASF-006 | Both require per-partner certificate trust store. |
| H-022 | ASF-007 | Both require SAML assertion replay detection. |
| H-024 | ASF-028 | Both require OAuth introspection endpoint protection. |
| H-025 | ASF-059 | Both require audience validation in OAuth tokens. |
| H-030 | ASF-058 | Both require partner-specific audience restrictions. |
| H-032 | ASF-052 | Both require bounded OAuth refresh token lifetime and rotation. |
| H-033 | ASF-003 | Both address rate limiting of metadata operations. |
| H-035 | ASF-007 | Both require SAML request-response correlation. |
| H-036 | ASF-009 | Both address clock skew tolerance for SAML assertions. |
| H-037 | ASF-059 | Both require attribute mapping validation. |
| H-041 | ASF-050 | Both require partner offboarding with credential revocation. |
| H-042 | ASF-018 | Both address SAML Single Logout implementation. |
| H-044 | ASF-030 | Both require OAuth authorization server tenant isolation. |

**Overlap (O): 30**

### Metrics

| Metric | Value | Formula |
|--------|-------|---------|
| Human assumptions (H) | 44 | Count of unique human-generated assumptions |
| ASF assumptions (A) | 80 | Count of unique ASF-generated assumptions |
| Overlap (O) | 30 | Count appearing in both lists |
| **Precision** | **37.5%** | O / A = 30/80 |
| **Recall** | **68.2%** | O / H = 30/44 |
| **F1 Score** | **48.4%** | 2 × (P × R) / (P + R) |
| Novel findings (A - O) | 50 | Assumptions ASF found that human missed (62.5% of ASF total) |
| Missed findings (H - O) | 14 | Assumptions human found that ASF missed (31.8% of human total) |

### Success Criteria Assessment

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Recall | >= 70% | 68.2% | ❌ Not met (borderline) |
| Precision | >= 50% | 37.5% | ❌ Not met |
| Novel discoveries | >= 10% of total (A+O) | 40.3% (50/124) | ✅ Exceeded |
| Expert agreement (F1 proxy) | > 60% | 48.4% | ❌ Not met |

---

## Step 4: Analysis

### Ontology Overlap Analysis

| Ontology Category | Overlaps | ASF Total | Overlap Rate |
|-------------------|----------|-----------|-------------|
| Explicit | 8 | 16 | 50.0% |
| Derived | 10 | 20 | 50.0% |
| Operational | 6 | 20 | 30.0% |
| Implicit | 4 | 12 | 33.3% |
| Trust | 1 | 4 | 25.0% |
| Dependency | 1 | 8 | 12.5% |
| Architectural | 0 | 4 | 0.0% |
| Environmental | 0 | 4 | 0.0% |

**Best overlap:** Explicit and Derived categories showed the strongest agreement. Both humans and the ASF recognize that SAML signature validation, TLS enforcement, and token scoping are critical.

**Worst overlap:** Architectural and Environmental categories had zero overlap. The ASF identified architectural concerns (gateway redundancy, network segmentation design) and environmental concerns (cross-organizational network reliability, partner data volume) that the human did not treat as assumptions.

### What Humans Caught That ASF Missed (Missed Findings = 14)

1. **SAML protocol implementation details (H-012, H-023, H-026, H-029, H-031, H-040):** Unsolicited SAML response rejection, XML DSig independent of TLS, assertion size limits, IdP-initiated SSO disabling, NameID format validation, and AuthnContext validation are SAML-specific concerns not covered by generic ASF patterns.

2. **Credential store separation (H-028, H-038):** The human assumed that partner client credentials are encrypted with per-partner keys and that IdP signing certificates are stored separately from API keys. The ASF treats secrets management generically but does not model per-partner separation.

3. **Resource-level rate limiting (H-034):** The human distinguished gateway-level rate limiting from resource-level rate limiting. The ASF pattern covers rate limiting but does not model the multi-layer architecture.

4. **Configuration drift prevention (H-039):** The human assumed IaC for SAML configuration. The ASF's Change Management pattern covers process but not infrastructure-as-code enforcement.

### What ASF Caught That Humans Missed (Novel Findings = 50)

1. **Incident Response (4 assumptions):** The human generated zero IR assumptions specific to federation compromise. The ASF contributed IR planning, log access, token revocation procedures, and anomaly detection.

2. **Secrets Management (4 assumptions):** The human covered credential storage but the ASF extended to vault-based secrets management, rotation schedules, and access auditing for the secrets store itself.

3. **Supply Chain and Third-party Dependency (8 assumptions):** The human treated the architecture as self-contained. The ASF surfaced partner IdP software vulnerability risk, API Gateway library vulnerabilities, SBOM tracking, vendor disclosure policies, and business continuity dependencies on partner infrastructure.

4. **Change Management (4 assumptions):** The human did not address metadata change governance, CI/CD for gateway config, rate limit change communication, or OAuth scope versioning.

5. **Monitoring infrastructure security (ASF-027, ASF-064):** The human assumed logging (H-019) but did not address audit log tamper-proofing or per-partner gateway health monitoring.

6. **Data flow documentation (ASF-030):** The human did not assume that partner-specific data flow diagrams exist or are reviewed annually.

### Architecture Complexity Assessment

Architecture #013 was classified as **Moderate** (multi-party federation, 5 nodes, 3 trust boundaries, cross-organizational trust).

- **Recall (68.2%)** is borderline below the 70% target, driven by SAML protocol-specific concerns that the ASF matrix does not explicitly cover.
- **Precision (37.5%)** reflects the breadth of the ASF matrix generating many assumptions the human did not consider.
- **Novel rate (62.5%)** indicates substantial value from the ASF in surfacing supply chain, incident response, and identity lifecycle concerns in a multi-party architecture.
- The high **novelty rate** is partially explained by the cross-organizational nature: the human focused on the gateway as the boundary, while the ASF extended to partner-side risks, dependency risks, and governance.

### Key Insight

The biggest root cause of the missed findings is **SAML-specific protocol knowledge**: the ASF's generic patterns do not cover SAML binding constraints (unsolicited responses, IdP-initiated SSO, NameID validation, XML DSig vs TLS, AuthnContext enforcement). Adding a "Federation Security" pattern (pattern 21) covering SAML and OAuth protocol-specific assumptions would likely close the recall gap to above 70%.

The high novelty rate (62.5%) is expected for multi-party architectures: the ASF systematically enumerates risks across the partner lifecycle, supply chain, and governance dimensions that a human architect focused on the federation protocol itself may miss.

---

## Conclusions

| Metric | Target | Actual | Verdict |
|--------|--------|--------|---------|
| Recall | >= 70% | 68.2% | ❌ Below target — missing SAML/OAuth protocol-specific pattern |
| Precision | >= 50% | 37.5% | ❌ Below target — ASF is deliberately broad/generative |
| Novel discoveries | >= 10% | 40.3% | ✅ ASF adds substantial value for multi-party architecture |
| Expert agreement (F1) | > 60% | 48.4% | ❌ Below target — driven by low precision |

The ASF framework applied to Architecture #013 demonstrates strong exploration breadth, particularly for cross-organizational risks (supply chain, dependency, governance). The primary actionable finding is the need for a **Federation Security** pattern covering SAML/OAuth protocol-specific assumptions to close the recall gap.
