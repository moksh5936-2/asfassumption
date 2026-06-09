# Architecture 2 — Security Assumption Review

## 1. Consensus Matrix

| # | Assumption | GPT | Gemini | Gemma | Keep? |
|---|-----------|:---:|:------:|:-----:|:-----:|
| 1 | User devices are free of malware | ✓ | ✓ |  | ✅ |
| 2 | Browsers properly validate TLS certificates | ✓ |  |  | ✅ |
| 3 | TLS is correctly configured between browsers and Okta | ✓ |  |  | ✅ |
| 4 | Okta infrastructure is secure | ✓ |  |  | ✅ |
| 5 | Administrative access to Okta is protected | ✓ |  |  | ✅ |
| 6 | MFA methods are resistant to phishing | ✓ |  |  | ✅ |
| 7 | MFA enrollment is secure | ✓ |  |  | ✅ |
| 8 | MFA recovery procedures are secure | ✓ |  |  | ✅ |
| 9 | MFA is resistant to push fatigue / MFA spamming |  |  | ✓ | ✅ |
| 10 | Authentication policies are correctly configured | ✓ |  |  | ✅ |
| 11 | SAML assertions are digitally signed | ✓ |  |  | ✅ |
| 12 | Service providers verify SAML assertion signatures | ✓ | ✓ |  | ✅ |
| 13 | Service providers validate assertion expiration | ✓ |  |  | ✅ |
| 14 | Service providers validate audience restrictions | ✓ |  |  | ✅ |
| 15 | Service providers validate issuer identity | ✓ | ✓ |  | ✅ |
| 16 | SAML signing private keys remain strictly confidential | ✓ |  | ✓ | ✅ |
| 17 | SAML certificate rotation occurs correctly | ✓ |  |  | ✅ |
| 18 | Clock synchronization is accurate across all systems | ✓ |  |  | ✅ |
| 19 | Federation metadata is accurate | ✓ |  |  | ✅ |
| 20 | Federation metadata updates are authenticated | ✓ |  |  | ✅ |
| 21 | Active Directory is secure | ✓ |  |  | ✅ |
| 22 | AD-to-Okta connection is authenticated and secure |  | ✓ |  | ✅ |
| 23 | Directory synchronization is trustworthy | ✓ |  |  | ✅ |
| 24 | Directory synchronization channels are protected | ✓ | ✓ |  | ✅ |
| 25 | JIT provisioning only creates legitimate users | ✓ |  |  | ✅ |
| 26 | JIT provisioning assigns appropriate roles | ✓ |  |  | ✅ |
| 27 | JIT deprovisioning is timely when membership or roles change | ✓ | ✓ |  | ✅ |
| 28 | JIT provisioning includes automated attribute verification |  |  | ✓ | ✅ |
| 29 | Group memberships are accurate | ✓ |  |  | ✅ |
| 30 | Session timeout (8 hours) is enforced consistently | ✓ |  |  | ✅ |
| 31 | Session termination truly invalidates tokens | ✓ |  |  | ✅ |
| 32 | Session tokens cannot be exfiltrated from browser storage | ✓ | ✓ |  | ✅ |
| 33 | Session tokens cannot be reused from a different network context |  |  | ✓ | ✅ |
| 34 | Session identifiers are unpredictable | ✓ |  |  | ✅ |
| 35 | Single Logout works correctly across all SPs | ✓ |  |  | ✅ |
| 36 | Service providers enforce authorization independently | ✓ |  |  | ✅ |
| 37 | Applications correctly consume SAML attributes for authZ | ✓ |  |  | ✅ |
| 38 | Applications do not trust unvalidated attributes | ✓ |  |  | ✅ |
| 39 | No application has a local authentication bypass | ✓ |  |  | ✅ |
| 40 | Legacy protocols (e.g., LDAP bind, basic auth) are disabled | ✓ |  |  | ✅ |
| 41 | Identity data integrity is preserved across the pipeline | ✓ |  |  | ✅ |
| 42 | No undocumented trust relationships exist | ✓ |  |  | ✅ |
| 43 | Administrative actions are logged | ✓ |  |  | ✅ |
| 44 | Logs are tamper-proof | ✓ |  |  | ✅ |
| 45 | Monitoring and alerting detect identity attacks | ✓ |  |  | ✅ |
| 46 | Help-desk / recovery processes are secure against social engineering | ✓ |  |  | ✅ |
| 47 | Insider administrators are trustworthy | ✓ |  |  | ✅ |
| 48 | Service provider applications are individually secure | ✓ |  |  | ✅ |

---

## 2. Deduplicated Assumption List

### Browser / Endpoint
1. User devices are free of malware (session-hijacking, infostealers)
2. Browsers properly validate TLS certificates
3. TLS is correctly configured between browsers and Okta

### IdP / Security
4. Okta infrastructure is secure
5. Administrative access to Okta is protected

### MFA
6. MFA methods are resistant to phishing
7. MFA enrollment is secure
8. MFA recovery procedures are secure
9. MFA is resistant to push fatigue / MFA spamming attacks
10. Authentication policies are correctly configured

### SAML Protocol
11. SAML assertions are digitally signed
12. Service providers verify SAML assertion signatures
13. Service providers validate assertion expiration
14. Service providers validate audience restrictions
15. Service providers validate issuer identity
16. SAML signing private keys remain strictly confidential
17. SAML certificate rotation occurs correctly
18. Clock synchronization is accurate across all systems
19. Federation metadata is accurate
20. Federation metadata updates are authenticated

### AD / Directory
21. Active Directory is secure
22. AD-to-Okta connection is authenticated and secure
23. Directory synchronization is trustworthy
24. Directory synchronization channels are protected

### JIT Provisioning
25. JIT provisioning only creates legitimate users
26. JIT provisioning assigns appropriate roles
27. JIT deprovisioning is timely when membership or roles change
28. JIT provisioning includes automated attribute verification
29. Group memberships are accurate

### Session Management
30. Session timeout (8 hours) is enforced consistently
31. Session termination truly invalidates tokens
32. Session tokens cannot be exfiltrated from browser storage
33. Session tokens cannot be reused from a different network context
34. Session identifiers are unpredictable
35. Single Logout works correctly across all SPs

### Authorization
36. Service providers enforce authorization independently
37. Applications correctly consume SAML attributes for authorization
38. Applications do not trust unvalidated attributes
39. No application has a local authentication bypass
40. Legacy protocols (e.g., LDAP bind, basic auth) are disabled
48. Service provider applications are individually secure

### Governance
41. Identity data integrity is preserved across the pipeline
42. No undocumented trust relationships exist
47. Insider administrators are trustworthy

### Monitoring / Logging
43. Administrative actions are logged
44. Logs are tamper-proof
45. Monitoring and alerting detect identity attacks

### Operations
46. Help-desk / recovery processes are secure against social engineering

---

## 3. Risk Scores

| # | Assumption | Likelihood | Impact | Risk |
|---|-----------|:----------:|:------:|:----:|
| 1 | User devices are free of malware | H | C | C |
| 2 | Browsers properly validate TLS certificates | L | H | M |
| 3 | TLS is correctly configured between browsers and Okta | L | H | M |
| 4 | Okta infrastructure is secure | L | C | H |
| 5 | Administrative access to Okta is protected | L | C | H |
| 6 | MFA methods are resistant to phishing | H | C | C |
| 7 | MFA enrollment is secure | M | H | H |
| 8 | MFA recovery procedures are secure | H | C | C |
| 9 | MFA is resistant to push fatigue / MFA spamming | H | C | C |
| 10 | Authentication policies are correctly configured | M | H | H |
| 11 | SAML assertions are digitally signed | L | C | H |
| 12 | Service providers verify SAML assertion signatures | M | C | C |
| 13 | Service providers validate assertion expiration | M | H | H |
| 14 | Service providers validate audience restrictions | M | H | H |
| 15 | Service providers validate issuer identity | M | C | C |
| 16 | SAML signing private keys remain strictly confidential | M | C | C |
| 17 | SAML certificate rotation occurs correctly | M | H | H |
| 18 | Clock synchronization is accurate across all systems | M | H | H |
| 19 | Federation metadata is accurate | M | H | H |
| 20 | Federation metadata updates are authenticated | M | C | C |
| 21 | Active Directory is secure | H | C | C |
| 22 | AD-to-Okta connection is authenticated and secure | M | C | C |
| 23 | Directory synchronization is trustworthy | M | C | C |
| 24 | Directory synchronization channels are protected | M | H | H |
| 25 | JIT provisioning only creates legitimate users | M | C | C |
| 26 | JIT provisioning assigns appropriate roles | M | H | H |
| 27 | JIT deprovisioning is timely when membership or roles change | H | H | H |
| 28 | JIT provisioning includes automated attribute verification | H | H | H |
| 29 | Group memberships are accurate | M | H | H |
| 30 | Session timeout (8 hours) is enforced consistently | M | H | H |
| 31 | Session termination truly invalidates tokens | M | H | H |
| 32 | Session tokens cannot be exfiltrated from browser storage | H | C | C |
| 33 | Session tokens cannot be reused from a different network context | H | H | H |
| 34 | Session identifiers are unpredictable | L | H | M |
| 35 | Single Logout works correctly across all SPs | M | M | M |
| 36 | Service providers enforce authorization independently | M | H | H |
| 37 | Applications correctly consume SAML attributes for authZ | M | H | H |
| 38 | Applications do not trust unvalidated attributes | M | H | H |
| 39 | No application has a local authentication bypass | M | C | C |
| 40 | Legacy protocols (e.g., LDAP bind, basic auth) are disabled | H | C | C |
| 41 | Identity data integrity is preserved across the pipeline | M | H | H |
| 42 | No undocumented trust relationships exist | M | C | C |
| 43 | Administrative actions are logged | L | H | M |
| 44 | Logs are tamper-proof | M | H | H |
| 45 | Monitoring and alerting detect identity attacks | M | H | H |
| 46 | Help-desk / recovery processes are secure against social engineering | H | C | C |
| 47 | Insider administrators are trustworthy | L | C | H |
| 48 | Service provider applications are individually secure | H | C | C |

---

## 4. STRIDE Mapping

| # | Assumption | STRIDE Category |
|---|-----------|:---------------:|
| 1 | User devices are free of malware | Tampering |
| 2 | Browsers properly validate TLS certificates | Spoofing |
| 3 | TLS is correctly configured between browsers and Okta | Tampering |
| 4 | Okta infrastructure is secure | Tampering |
| 5 | Administrative access to Okta is protected | Elevation of Privilege |
| 6 | MFA methods are resistant to phishing | Spoofing |
| 7 | MFA enrollment is secure | Spoofing |
| 8 | MFA recovery procedures are secure | Spoofing |
| 9 | MFA is resistant to push fatigue / MFA spamming | Spoofing |
| 10 | Authentication policies are correctly configured | Elevation of Privilege |
| 11 | SAML assertions are digitally signed | Tampering |
| 12 | Service providers verify SAML assertion signatures | Spoofing |
| 13 | Service providers validate assertion expiration | Repudiation |
| 14 | Service providers validate audience restrictions | Elevation of Privilege |
| 15 | Service providers validate issuer identity | Spoofing |
| 16 | SAML signing private keys remain strictly confidential | Spoofing |
| 17 | SAML certificate rotation occurs correctly | Denial of Service |
| 18 | Clock synchronization is accurate across all systems | Repudiation |
| 19 | Federation metadata is accurate | Spoofing |
| 20 | Federation metadata updates are authenticated | Spoofing |
| 21 | Active Directory is secure | Elevation of Privilege |
| 22 | AD-to-Okta connection is authenticated and secure | Tampering |
| 23 | Directory synchronization is trustworthy | Tampering |
| 24 | Directory synchronization channels are protected | Tampering |
| 25 | JIT provisioning only creates legitimate users | Spoofing |
| 26 | JIT provisioning assigns appropriate roles | Elevation of Privilege |
| 27 | JIT deprovisioning is timely when membership or roles change | Elevation of Privilege |
| 28 | JIT provisioning includes automated attribute verification | Elevation of Privilege |
| 29 | Group memberships are accurate | Elevation of Privilege |
| 30 | Session timeout (8 hours) is enforced consistently | Elevation of Privilege |
| 31 | Session termination truly invalidates tokens | Spoofing |
| 32 | Session tokens cannot be exfiltrated from browser storage | Information Disclosure |
| 33 | Session tokens cannot be reused from a different network context | Spoofing |
| 34 | Session identifiers are unpredictable | Spoofing |
| 35 | Single Logout works correctly across all SPs | Repudiation |
| 36 | Service providers enforce authorization independently | Elevation of Privilege |
| 37 | Applications correctly consume SAML attributes for authZ | Elevation of Privilege |
| 38 | Applications do not trust unvalidated attributes | Elevation of Privilege |
| 39 | No application has a local authentication bypass | Spoofing |
| 40 | Legacy protocols (e.g., LDAP bind, basic auth) are disabled | Spoofing |
| 41 | Identity data integrity is preserved across the pipeline | Tampering |
| 42 | No undocumented trust relationships exist | Elevation of Privilege |
| 43 | Administrative actions are logged | Repudiation |
| 44 | Logs are tamper-proof | Tampering |
| 45 | Monitoring and alerting detect identity attacks | Repudiation |
| 46 | Help-desk / recovery processes are secure against social engineering | Spoofing |
| 47 | Insider administrators are trustworthy | Elevation of Privilege |
| 48 | Service provider applications are individually secure | Elevation of Privilege |

---

## 5. Top 10 Critical Assumptions

| Rank | Assumption | Rationale |
|:----:|-----------|-----------|
| 1 | Service providers verify SAML assertion signatures | The entire federation trust model collapses if SPs accept unsigned or forged assertions. Without signature verification, an attacker can mint arbitrary identities without compromising any upstream system. |
| 2 | SAML signing private keys remain strictly confidential | Leaked signing keys allow forging valid SAML assertions for any user without ever touching Okta. Key compromise defeats all upstream controls including MFA, policies, and monitoring. |
| 3 | Active Directory is secure | AD is the authoritative user store. Compromise enables attacker creation of users, group manipulation, and privilege escalation that propagates through sync into every SP application. |
| 4 | Directory synchronization is trustworthy | The sync bridge between AD and Okta is a critical trust boundary. A compromised or misconfigured sync can inject rogue users, modify attributes, or propagate garbage data that drives authZ decisions in every federated app. |
| 5 | No application has a local authentication bypass | If any SP maintains a local login path (backdoor admin creds, separate auth source), attackers can bypass Okta MFA, policies, and monitoring entirely — making federation irrelevant for that app. |
| 6 | Legacy protocols (e.g., LDAP bind, basic auth) are disabled | Legacy auth protocols often lack MFA enforcement and operate outside the SAML trust model. An active LDAP bind or HTTP basic auth path on any SP undermines the entire SSO architecture. |
| 7 | JIT provisioning only creates legitimate users | Automatic account creation extends implicit trust to any user who can authenticate via Okta. Without gating logic, an attacker who gains an Okta session also gets accounts in every JIT-enabled SP. |
| 8 | MFA is resistant to push fatigue / MFA spamming | Attackers increasingly bypass MFA by flooding users with push notifications until they accidentally approve. This is the most common real-world MFA bypass and directly undermines the architecture's primary authentication control. |
| 9 | Session tokens cannot be exfiltrated from browser storage | The 8-hour session window means a stolen token gives the attacker prolonged access. Infostealer malware and browser extensions can silently dump session tokens, making this a high-exposure attack surface. |
| 10 | Monitoring and alerting detect identity attacks | Without effective detection, every other control failure goes unnoticed. SSO systems are high-value targets, and the absence of monitoring allows attackers to iterate, persist, and escalate over extended periods. |

---

## 6. Recommended Controls

| Rank | Assumption | Recommended Controls |
|:----:|-----------|---------------------|
| 1 | SPs verify SAML assertion signatures | Mandate signed assertions in Okta policy; automate SP onboarding validation with a test harness that sends unsigned/altered assertions and verifies rejection; perform regular penetration testing of SAML consumer endpoints. |
| 2 | SAML signing private key confidentiality | Store signing keys in a hardware security module (HSM) or key management service (KMS); rotate keys on a scheduled cadence (≤12 months); restrict key export permissions to a break-glass process with approval workflow and audit trail. |
| 3 | Active Directory security | Deploy dedicated AD tier-0 administration model; enable advanced audit logging (e.g., Microsoft 365 Defender for Identity); enforce privilege access workstations (PAWs) for AD administrators; implement regular AD security baseline scanning. |
| 4 | Directory sync trust | Use mutual TLS or authenticated API calls between AD and Okta; restrict sync service accounts to least privilege (read-only, scoped OUs); monitor sync logs for unexpected user creations or attribute modifications; deploy Okta's Directory Integration with strict signing. |
| 5 | No local auth bypass on SPs | Mandate a quarterly SP audit that checks for local accounts, alternative auth methods, and direct database access; use an external crawl tool to probe each SP's login page for non-SSO paths; enforce via policy that SPs must disable local login. |
| 6 | Legacy protocol hardening | Inventory all SPs and disable legacy protocols (LDAP simple bind, NTLM, HTTP basic auth); where legacy is unavoidable, front with a reverse proxy that enforces MFA; scan the environment quarterly for unexpected legacy endpoints. |
| 7 | JIT user gating | Implement a pre-provisioning approval workflow for sensitive SPs; scope JIT group-to-role mappings at the most restrictive level; deploy a reconciliation job that cross-checks JIT-created accounts against HR data daily. |
| 8 | MFA spamming resistance | Enforce number matching or FIDO2/WebAuthn in Okta MFA policies over push-only; configure push notification timeouts and lockout thresholds; educate users to never approve unexpected MFA prompts. |
| 9 | Session token protection | Set session token as HttpOnly and SameSite=Strict; reduce session lifetime below 8 hours for high-risk applications; implement Okta device posture (e.g., Jamf/Intune) to block sessions from non-compliant or unknown devices; monitor for anomalous token reuse (geo/ASN jumps). |
| 10 | Detection and monitoring | Stream Okta System Log and AD audit logs to a SIEM (e.g., Splunk, Sentinel); build detection rules for: MFA push fatigue spikes, new admin account creation, failed SAML assertion validation, sync attribute anomalies, and login-from-unfamiliar-location events; conduct a tabletop exercise for SSO incident response quarterly. |

---

## 7. Summary Stats

| Category | Count |
|:--------:|:-----:|
| **Total Assumptions** | 48 |
| **Critical (C)** | 18 |
| **High (H)** | 25 |
| **Medium (M)** | 5 |
| **Low (L)** | 0 |

**Breakdown by Risk Rating:**
- **Critical (C):** Assumptions 1, 6, 8, 9, 12, 15, 16, 20, 21, 22, 23, 25, 32, 39, 40, 42, 46, 48 = 18
- **High (H):** Assumptions 4, 5, 7, 10, 11, 13, 14, 17, 18, 19, 24, 26, 27, 28, 29, 30, 31, 33, 36, 37, 38, 41, 44, 45, 47 = 25
- **Medium (M):** Assumptions 2, 3, 34, 35, 43 = 5
- **Low (L):** None
