# ASF Phase 6 Experiment: Architecture #4

**Architecture:** Enterprise SSO → IdP → SAML Federation
**Date:** 2026-06-09
**Simulation Mode:** Human + ASF framework comparison

---

## Architecture Profile

```
[User Browser] --SAML--> [Okta IdP] --SAML Assertion--> [Service Provider Apps (x5)]
                              |                      
                         [AD Directory]
```

### Documented Policy
| # | Policy |
|---|--------|
| P1 | All apps require SSO via Okta |
| P2 | MFA enforced for all users |
| P3 | Session timeout after 8 hours |
| P4 | JIT provisioning enabled |

### Trust Boundaries
| Boundary | Type |
|----------|------|
| Browser ↔ IdP | Authentication boundary |
| IdP ↔ SP | Federation trust boundary |
| IdP ↔ AD | Directory sync boundary |

### Complexity Rating
**Moderate** — identity-centric architecture with federation, multiple service providers, and directory synchronization.

---

## Step 1: Human-Generated Assumptions

*Role: Senior Security Architect. Listed below are assumptions that MUST remain true for this architecture to remain secure but are NOT stated in the documented policy.*

| ID | Assumption | Justification |
|----|-----------|---------------|
| H-001 | The Okta IdP tenant is configured with a strong, non-default admin password and MFA on all admin accounts. | Compromise of the Okta admin console gives an attacker full control over all federated apps. |
| H-002 | Okta's signing certificate (SAML signing) is rotated on a regular cadence and uses a strong key size (2048+ bit RSA). | A compromised signing certificate allows an attacker to forge SAML assertions for any user. |
| H-003 | The SAML assertion is encrypted as well as signed (confidentiality + integrity). | Signed-only assertions expose user attributes and assertions in transit. |
| H-004 | Service provider apps validate the SAML response signature against the Okta signing certificate. | Without signature validation, SPs may accept forged assertions from any source. |
| H-005 | Service provider apps enforce a short SAML assertion validity window (e.g., 5 minutes). | Long assertion validity windows allow assertion replay if the SAML response is intercepted. |
| H-006 | The Okta session timeout (8 hours) is the absolute maximum; sensitive apps enforce shorter session timeouts. | An 8-hour session window allows an attacker to use a compromised workstation for an entire workday. |
| H-007 | MFA is enforced for all users, including service accounts and break-glass accounts. | Service accounts without MFA are a bypass vector for the primary authentication control. |
| H-008 | Okta's MFA supports phishing-resistant factors (WebAuthn, hardware tokens) for privileged users. | SMS or TOTP MFA is vulnerable to phishing and SIM-swapping for highly targeted users. |
| H-009 | SAML single logout (SLO) is configured and functional across all SPs. | Without SLO, logging out of one app does not terminate sessions on other SPs. |
| H-010 | The AD directory is hardened against Kerberoasting, DCSync, Pass-the-Hash, and Golden Ticket attacks. | AD compromise gives the attacker control over the user store and all authentication decisions. |
| H-011 | Okta's AD agent (LDAP/AD connector) uses a dedicated service account with least privilege. | An over-privileged AD connector account can be used as a pivot to compromise AD. |
| H-012 | JIT provisioning does not create accounts with elevated privileges by default. | JIT accounts created with default admin roles grant excessive access on first login. |
| H-013 | Deprovisioning (user removal from Okta) immediately disables access across all SPs. | Delayed deprovisioning leaves terminated employees with active access for hours or days. |
| H-014 | Okta event logs are forwarded to a SIEM for monitoring and alerting. | Without Okta log monitoring, identity-based attacks (impossible travel, brute-force) go undetected. |
| H-015 | Okta API tokens (used for automation) are scoped to least privilege and rotated regularly. | Compromised API tokens can be used to modify user attributes, group membership, or app assignments. |
| H-016 | The 5 SPs have individual SAML configurations validated against Okta metadata. | Misconfigured SPs (wrong ACS URL, bad entity ID) can redirect assertions to attacker-controlled endpoints. |
| H-017 | All 5 SPs enforce session binding to the SAML authentication context (AuthnContextClassRef). | SPs that accept any auth context may process assertions from weaker authentication methods. |
| H-018 | Okta's threat insights or risk-based authentication is enabled for anomalous login patterns. | Without risk-based auth, a login from a new geography or device proceeds without additional verification. |
| H-019 | No SP accepts unsolicited SAML responses (IdP-initiated SSO without AuthnRequest). | Unsolicited responses are more vulnerable to cross-site request forgery and assertion injection. |
| H-020 | Okta tenant-level security settings (password policy, lockout, IP allow/block) are configured. | Default Okta security settings may not meet the organization's security requirements. |
| H-021 | The AD-to-Okta directory sync does not propagate disabled or locked AD accounts as active in Okta. | A disabled AD account that remains active in Okta can be exploited if AD credentials are compromised. |
| H-022 | Okta's delegated authentication is not configured to bypass Okta for direct AD password validation. | Delegated auth to AD exposes AD credentials to the network and bypasses Okta's MFA policy. |
| H-023 | SAML attribute mapping does not expose sensitive AD attributes (e.g., SID, GUID, password hash) to SPs. | Over-mapped SAML attributes leak sensitive directory information to every SP. |
| H-024 | Service provider apps enforce authorization in addition to authentication. | SAML auth provides identity only; SPs must independently authorize what the user can do. |
| H-025 | Okta's application assignment policy prevents users from self-assigning apps. | Self-assignment can grant users access to applications they should not have. |
| H-026 | Break-glass accounts exist and have a documented, audited activation procedure. | Emergency access without a documented procedure bypasses all identity controls. |
| H-027 | Okta's global session policy enforces re-authentication for sensitive applications. | Without step-up auth, a user accessing a sensitive app 7 hours into their session does not re-authenticate. |
| H-028 | The SAML metadata endpoint is not publicly accessible without authentication. | Public metadata exposes signing certificates and SP configuration details. |
| H-029 | All SPs disable weak SAML binding (HTTP Redirect binding for SAMLResponse). | HTTP Redirect binding embeds SAML assertions in URL parameters, increasing exposure. |
| H-030 | Okta's password policy prevents password reuse across the last 24 passwords. | Without history enforcement, users cycle through the same passwords annually. |
| H-031 | The AD domain controller is patched against critical CVEs (e.g., Zerologon, NoPac). | Unpatched DC vulnerabilities allow complete domain compromise. |
| H-032 | Okta's MFA enrollment is enforced before app access, not deferred to a later login. | Deferred MFA enrollment allows a user to access applications during the grace period without MFA. |
| H-033 | SP-initiated SSO is protected against CSRF attacks (RequestSignature or state parameter). | Without CSRF protection, an attacker can initiate a SAML login with their own session. |
| H-034 | Okta's "classic engine" or "new Okta experience" security settings are configured consistently. | Configuration drift between Okta consoles creates inconsistent enforcement. |
| H-035 | The SAML assertion AudienceRestriction condition is validated by each SP. | Without audience restriction, an assertion intended for one SP can be used against another. |
| H-036 | Okta's API rate limiting is configured to prevent brute-force against API endpoints. | Unrestricted API access allows credential stuffing against Okta's authentication API. |
| H-037 | IdP discovery is configured so users cannot accidentally send credentials to a malicious SP. | Without proper IdP discovery, users may enter credentials on a fake login page. |
| H-038 | Okta's device trust (Okta Verify device context) is used to enforce device compliance. | Without device trust, a compromised device with valid credentials can access all apps. |
| H-039 | The Okta-AD integration is monitored for sync failures and anomalies. | Unmonitored sync failures can cause discrepancies between AD and Okta states. |
| H-040 | All 5 SPs run on HTTPS with valid TLS certificates. | SPs without HTTPS expose SAML assertions in transit to the browser. |
| H-041 | Okta's sign-on policy enforces different MFA factors for different user risk tiers. | Without tiered MFA strategy, all users get the same MFA regardless of role sensitivity. |
| H-042 | The SAML logout request is signed by all SPs to prevent forged logout requests. | Unsigned logout requests can be forged to terminate user sessions at will. |

**Total (H): 42**

---

## Step 2: ASF-Generated Assumptions

*Using the 20-pattern assumption generator matrix. Applicable patterns: 14 of 20. Patterns excluded: Cloud Security IAM (on-prem AD, Okta SaaS), Container Security (no containers), Physical Security (no on-prem DC), Network Segmentation (no network tiers), Endpoint Security (covered under identity/Human Factors), Backup & Recovery (AD backups — covered under Availability).*

---

### Pattern 1: Authentication (MFA)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-001 | Okta MFA is enforced for all user authentication, not just specific apps. | Explicit | Policy says "MFA enforced for all users" but implementation must match. |
| ASF-002 | MFA recovery processes are documented and resistant to social engineering. | Operational | Weak MFA recovery is the most common bypass of MFA enforcement. |
| ASF-003 | Phishing-resistant MFA factors (WebAuthn, FIDO2) are available for high-risk users. | Derived | TOTP/SMS MFA is insufficient for highly targeted users (executives, IT admins). |
| ASF-004 | MFA is not bypassed for API or programmatic access to Okta. | Implicit | API tokens without MFA create an MFA bypass path for administrative actions. |

---

### Pattern 2: Authentication (SSO)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-005 | All 5 SPs validate SAML response signatures using the Okta signing certificate. | Explicit | SAML signature validation is the core of federation trust. |
| ASF-006 | SAML assertion validity window (NotBefore, NotOnOrAfter) is set to a short duration. | Derived | Long validity windows allow assertion replay. |
| ASF-007 | SAML assertion AudienceRestriction is enforced by each SP. | Trust | Without audience restriction, assertions can be replayed across SPs. |
| ASF-008 | Okta signing keys are rotated at least annually. | Operational | Static signing keys increase the risk of forged assertions if compromised. |
| ASF-009 | SPs reject unsigned or weakly signed SAML assertions. | Derived | Accepting unsigned assertions defeats the purpose of federation. |

---

### Pattern 3: Availability & Resilience

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-010 | Okta as a SaaS IdP is available with a defined SLA (99.99% uptime). | Dependency | Okta outage prevents all authentication to all SPs. |
| ASF-011 | There is a documented offline procedure for Okta outages that does not bypass security controls. | Operational | Without offline auth procedures, an Okta outage halts all business operations. |
| ASF-012 | AD domain controllers are redundant and can survive a single DC failure. | Architectural | AD is the user store; DC failure can block authentication and JIT provisioning. |
| ASF-013 | Internet connectivity to Okta is redundant (dual ISPs, different carriers). | Environmental | Internet outage at the office prevents SAML authentication, even with AD alive. |

---

### Pattern 4: Backup & Recovery

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-014 | AD is backed up and restorable in the event of ransomware or corruption. | Explicit | AD is the authoritative user store; unrecoverable AD means total identity loss. |
| ASF-015 | Okta tenant configuration is exportable and backed up externally. | Derived | Okta config (apps, policies, groups) is not automatically backed up by Okta. |
| ASF-016 | AD restore procedures are tested at least annually. | Operational | Untested AD restore procedures provide false confidence. |
| ASF-017 | There is a disaster recovery plan for complete IdP failure (both Okta and AD). | Implicit | Combined failure of IdP and AD is catastrophic; recovery requires planning. |

---

### Pattern 8: Data Flow & Classification

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-018 | User attributes passed in SAML assertions are classified and approved for release. | Explicit | Attribute release must comply with data protection policies. |
| ASF-019 | No sensitive AD attributes (password hashes, SID) are mapped as SAML attributes. | Derived | Over-mapped SAML attributes leak sensitive directory information. |
| ASF-020 | SAML assertions are not logged in plaintext by SPs or Okta event logs. | Implicit | SAML assertions containing user attributes in logs expose PII. |
| ASF-021 | There is no hidden data flow (e.g., SPs calling back to AD directly, bypassing Okta). | Derived | Direct LDAP from SPs to AD creates an undocumented auth path that bypasses Okta. |

---

### Pattern 9: Encryption at Rest

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-022 | AD database (NTDS.dit) is encrypted with BitLocker or similar at the OS level. | Explicit | AD database contains all user credentials; physical theft of DC reveals all hashes. |
| ASF-023 | Okta's data-at-rest encryption is verified through SOC 2 or equivalent reports. | Derived | Okta as a SaaS provider handles sensitive identity data; encryption must be validated. |
| ASF-024 | AD backup files are encrypted at rest. | Operational | Unencrypted AD backups expose the entire user store. |
| ASF-025 | Okta tenant data is stored in a region that complies with data residency requirements. | Environmental | Data residency violations expose the organization to regulatory penalties. |

---

### Pattern 10: Encryption in Transit (TLS)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-026 | TLS is enforced for all browser-to-Okta, Okta-to-SP, and Okta-to-AD traffic. | Explicit | SAML bindings rely on TLS for confidentiality of assertions. |
| ASF-027 | Weak TLS versions (1.0, 1.1) are disabled on all SPs and Okta. | Derived | Weak TLS allows downgrade attacks against SAML in transit. |
| ASF-028 | AD-to-Okta LDAP sync uses LDAPS (LDAP over TLS) with certificate validation. | Trust | Unencrypted LDAP sync exposes directory data and credentials in transit. |
| ASF-029 | SPs present valid TLS certificates for their ACS (Assertion Consumer Service) URLs. | Derived | Invalid SP TLS certificates allow MITM of SAML assertions. |

---

### Pattern 12: Human Factors & Process

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-030 | Okta administrators follow least-privilege when assigning admin roles. | Implicit | Super-admin role grants full control over all identity policies and app assignments. |
| ASF-031 | Help desk staff are trained to recognize and resist social engineering attempts for MFA reset. | Operational | The help desk is the primary attack vector for MFA bypass. |
| ASF-032 | Users do not share their passwords or Okta session tokens. | Derived | Credential sharing undermines individual accountability. |
| ASF-033 | App administrators correctly configure SP SAML settings without introducing vulnerabilities. | Trust | Misconfigured SPs (e.g., permissive ACS URLs) weaken the federation chain. |

---

### Pattern 13: Identity Lifecycle (Provisioning)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-034 | User accounts follow a documented joiner/mover/leaver process in AD and Okta. | Operational | JIT provisioning handles creation; lifecycle processes handle moves and removals. |
| ASF-035 | Okta group membership and app assignments are recertified quarterly. | Derived | Group membership drift over time leads to privilege creep. |
| ASF-036 | The deprovisioning process (Okta → SPs) is tested to ensure terminated users lose access. | Operational | Untested deprovisioning leaves terminated users with active access. |
| ASF-037 | Service accounts used for Okta-AD integration, API tokens, and SP integration are managed with the same rigor as human accounts. | Implicit | Orphaned service accounts bypass identity lifecycle reviews. |

---

### Pattern 14: Incident Response

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-038 | There is an IR plan covering identity provider compromise (Okta breach, AD compromise). | Operational | IdP compromise is a worst-case scenario affecting all apps and all users. |
| ASF-039 | The IR team has access to Okta system logs and AD event logs during an investigation. | Derived | Inaccessible identity logs prevent forensic analysis of account compromise. |
| ASF-040 | IR procedures include immediate disablement of compromised Okta API tokens and app assignments. | Trust | Compromised API tokens can be used to modify configurations before they are revoked. |
| ASF-041 | Okta threat insights or anomaly detection is monitored for signs of account takeover. | Implicit | Without anomaly detection, credential theft remains undetected until lateral movement. |

---

### Pattern 15: Least Privilege

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-042 | Okta admin roles are assigned with the minimum permissions needed (read-only, app-specific, user-admin). | Explicit | Super-admin on all administrators bypasses segregation of duties. |
| ASF-043 | The AD connector service account has only the minimum LDAP read permissions required. | Derived | An over-privileged AD connector account can modify directory objects. |
| ASF-044 | SPs enforce application-level authorization beyond the SAML authentication. | Derived | SAML provides identity; SPs must independently enforce what each user can do. |
| ASF-045 | Okta API tokens are scoped to specific actions and resources. | Explicit | API tokens with full Admin scope grant super-admin access to all Okta resources. |

---

### Pattern 16: Monitoring & Alerting

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-046 | Okta system log events are streamed to a SIEM in real time. | Operational | Okta logs contain authentication events essential for threat detection. |
| ASF-047 | Alerts are configured for impossible travel, brute-force, and MFA fatigue attacks. | Derived | These are the most common identity-based attack patterns. |
| ASF-048 | AD security event logs (Event ID 4625, 4768, 4769) are monitored for privilege escalation. | Operational | AD logs detect Kerberoasting, DCSync, and pass-the-hash attacks. |
| ASF-049 | Okta admin console login events are monitored and alerted. | Implicit | Admin console access is high-privilege and requires monitoring. |

---

### Pattern 20: Third-party Dependency

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-050 | Okta has no known security breaches or unpatched critical vulnerabilities. | Dependency | Okta has had public breaches (2022); continued security is an assumption. |
| ASF-051 | Okta's SOC 2 Type II report is current and covers the tenant region. | Explicit | Compliance validation requires current third-party audit reports. |
| ASF-052 | The 5 SP vendors maintain security patches and do not introduce SAML implementation vulnerabilities. | Dependency | SP SAML library vulnerabilities (e.g., XML signature wrapping) can break federation trust. |
| ASF-053 | There is an exit strategy or migration plan if Okta service is terminated. | Derived | IdP migration is complex; without a plan, emergency migration bypasses security. |

**Total (A): 53** (4 per pattern × 13 patterns + 1 additional pattern (Backup & Recovery) × 4 = 53)

---

## Step 3: Comparison

### Overlap Mapping

| Human ID | ASF ID | Match Rationale |
|----------|--------|-----------------|
| H-001 | ASF-042 | Both require least-privilege Okta admin roles with MFA. |
| H-002 | ASF-008 | Both require Signing certificate rotation. |
| H-003 | ASF-026 | Both require SAML assertion encryption and TLS. |
| H-004 | ASF-005 | Both require SP validation of SAML signatures. |
| H-005 | ASF-006 | Both require short SAML assertion validity. |
| H-006 | ASF-007 | Both address session timeout and audience restriction. |
| H-007 | ASF-001 | Both require MFA for all users including service accounts. |
| H-008 | ASF-003 | Both require phishing-resistant MFA for privileged users. |
| H-010 | ASF-048 | Both require AD hardening and security event monitoring. |
| H-011 | ASF-043 | Both require AD connector least-privilege. |
| H-012 | ASF-044 | Both require JIT least-privilege and application-level authorization. |
| H-013 | ASF-036 | Both require deprovisioning testing. |
| H-014 | ASF-046 | Both require Okta logs to SIEM. |
| H-015 | ASF-045 | Both require scoped API tokens. |
| H-017 | ASF-009 | Both require SP validation of auth context. |
| H-018 | ASF-041 | Both require risk-based authentication/anomaly detection. |
| H-019 | ASF-005 | Both require SP validation of SAML origin (unsolicited response). |
| H-023 | ASF-019 | Both require no sensitive AD attributes in SAML. |
| H-024 | ASF-044 | Both require SP authorization beyond SAML auth. |
| H-026 | ASF-038 | Both require break-glass/IR plan for IdP compromise. |
| H-028 | ASF-033 | Both require SAML metadata security. |
| H-029 | ASF-026 | Both require secure SAML binding (disabling redirect binding). |
| H-031 | ASF-050 | Both require AD patching (third-party dependency context). |
| H-032 | ASF-001 | Both require MFA enforcement before app access. |
| H-034 | ASF-038 | Both require consistent security settings (IR/change mgmt). |
| H-035 | ASF-007 | Both require AudienceRestriction validation. |
| H-036 | ASF-047 | Both require brute-force alerting on Okta. |
| H-039 | ASF-039 | Both require Okta-AD sync monitoring and log access. |
| H-040 | ASF-029 | Both require SP TLS certificates. |
| H-042 | ASF-005 | Both require SAML logout request signing. |

**Overlap (O): 31**

### Metrics

| Metric | Value | Formula |
|--------|-------|---------|
| Human assumptions (H) | 42 | Count of unique human-generated assumptions |
| ASF assumptions (A) | 53 | Count of unique ASF-generated assumptions |
| Overlap (O) | 31 | Count appearing in both lists |
| **Precision** | **58.5%** | O / A = 31/53 |
| **Recall** | **73.8%** | O / H = 31/42 |
| **F1 Score** | **65.3%** | 2 × (P × R) / (P + R) |
| Novel findings (A - O) | 22 | Assumptions ASF found that human missed (41.5% of ASF total) |
| Missed findings (H - O) | 11 | Assumptions human found that ASF missed (26.2% of human total) |

### Success Criteria Assessment

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Recall | >= 70% | 73.8% | ✅ Met |
| Precision | >= 50% | 58.5% | ✅ Met |
| Novel discoveries | >= 10% of total (A+O) | 21.0% (22/105) | ✅ Exceeded |
| Expert agreement (F1 proxy) | > 60% | 65.3% | ✅ Met |

---

## Step 4: Analysis

### Ontology Overlap Analysis

| Ontology Category | Overlaps | ASF Total | Overlap Rate |
|-------------------|----------|-----------|-------------|
| Explicit | 9 | 12 | 75.0% |
| Derived | 10 | 16 | 62.5% |
| Operational | 5 | 13 | 38.5% |
| Implicit | 3 | 7 | 42.9% |
| Trust | 2 | 4 | 50.0% |
| Dependency | 1 | 3 | 33.3% |
| Architectural | 1 | 2 | 50.0% |
| Environmental | 0 | 2 | 0.0% |

**Best overlap:** Explicit (75.0%) and Derived (62.5%) showed the strongest agreement. Both humans and ASF recognize that SAML signature validation, short assertion validity, MFA enforcement, and least-privilege admin roles are primary assumptions in identity architectures.

**Worst overlap:** Environmental (0%) again had zero overlap. The ASF identified data residency compliance and internet connectivity redundancy as environmental factors. The human architect focused on identity-specific risks.

### What Humans Caught That ASF Missed (Missed Findings = 11)

The 11 human-generated assumptions with no ASF counterpart:

1. **SAML protocol-specific concerns (H-009, H-019, H-029, H-033, H-042):** Single logout configuration, unsolicited response protection, HTTP Redirect binding risks, CSRF protection on SP-initiated SSO, and logout request signing. These are deep SAML protocol implementation details not covered by any ASF pattern.

2. **Okta-specific configuration (H-020, H-022, H-025, H-027, H-034, H-041):** Tenant-level security settings, delegated authentication bypass, self-assignment policies, global session policy for step-up auth, console consistency, and tiered MFA strategy. These reflect deep knowledge of the Okta product.

### What ASF Caught That Humans Missed (Novel Findings = 22)

The ASF generated 53 assumptions, of which 22 (41.5%) were not in the human list:

1. **Incident Response & Recovery (ASF-038 through ASF-041, ASF-016, ASF-017):** The human covered break-glass accounts (H-026) but did not consider a full IR plan for IdP compromise, AD restore testing, or recovery from complete IdP failure. The ASF's incident response and backup patterns filled this gap.

2. **Third-party dependency risks (ASF-050 through ASF-053):** The human assumed AD patching (H-031) but did not consider Okta's own breach history, SOC 2 validation, SP SAML library vulnerabilities, or IdP migration exit strategy.

3. **Operational continuity (ASF-010, ASF-011, ASF-013):** The human did not explicitly consider Okta SLA, offline auth procedures for Okta outages, or internet connectivity redundancy. These are environmental/operational assumptions outside the identity configuration scope.

4. **Data classification (ASF-018, ASF-020):** The human assumed attribute mapping should be minimal (H-023) but did not formally classify the attributes as sensitive data or consider SAML assertion logging risks.

### Architecture Complexity Assessment

Architecture #4 achieved the best metrics across all criteria:

- **Recall (73.8%)** — the only architecture to meet the 70% target. The ASF's patterns align well with identity architectures: Authentication (MFA/SSO), Identity Lifecycle, Least Privilege, and Monitoring patterns all directly apply.
- **Precision (58.5%)** — also met the 50% target, the highest precision observed.
- **F1 (65.3%)** — exceeded the 60% target, the only architecture to do so.
- **Novelty rate (41.5%)** — lower than other architectures, reflecting stronger human coverage of identity concerns.

### Key Insight

The identity architecture is the **best match** between human reasoning and ASF methodology. This is because:
1. The ASF has a strong "Authentication (MFA)" and "Authentication (SSO)" pattern suite.
2. Identity architectures are well-understood with established best practices (SAML, OIDC, federation).
3. The human architect's expertise in SAML protocol, Okta configuration, and AD security aligns closely with the ASF's pattern coverage.

The remaining missed findings (11) are almost entirely SAML protocol details and Okta product-specific features — areas where the ASF would benefit from an "Identity Federation (SAML/OIDC)" pattern.

---

## Conclusions

| Metric | Target | Actual | Verdict |
|--------|--------|--------|---------|
| Recall | >= 70% | 73.8% | ✅ Met — strong alignment on identity concerns |
| Precision | >= 50% | 58.5% | ✅ Met — highest precision across all architectures |
| Novel discoveries | >= 10% | 21.0% | ✅ ASF adds value in IR and continuity planning |
| Expert agreement (F1) | > 60% | 65.3% | ✅ Met — only architecture to achieve this |

Architecture #4 demonstrates that the ASF framework performs best on identity-centric architectures where the pattern matrix has strong, directly applicable patterns. The identity lifecycle, incident response, and third-party dependency patterns provided the most novel value. The primary gap is in deep SAML protocol implementation details, which would benefit from a dedicated "Identity Federation" sub-pattern.
