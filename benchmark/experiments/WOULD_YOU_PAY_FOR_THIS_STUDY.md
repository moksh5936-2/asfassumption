# "Would You Pay For This?" Study

**Purpose:** Determine if ASF's unique assumptions are valuable enough to security practitioners that they would change decisions — and pay for the capability.

**Method:** Present 10 Tier A (high-confidence) and 10 Tier C (ASF-unique) assumptions across 5 architectures. For each, ask 5 questions targeting real-world utility.

---

## Study Design

| Element | Detail |
|---------|--------|
| Target participants | 10-20 security architects, CISOs, AppSec engineers, cloud architects |
| Time per participant | 20-30 minutes |
| Format | Structured interview or self-administered survey |
| Architectures covered | 5 (diverse domains and complexity) |
| Assumptions per architecture | 4 (2 Tier A + 2 Tier C) |
| Total assumptions evaluated | 20 |
| Core question | "Would you pay to discover this automatically?" |

### The 5 Questions (per assumption)

1. **Reality** — Is this a real assumption? (Yes / Partial / No)
2. **Risk** — Would missing it create risk? (Critical / High / Medium / Low / None)
3. **Investigation** — Would you investigate it if surfaced? (Definitely / Probably / Maybe / Probably not / No)
4. **Incident history** — Have you seen incidents caused by this? (Yes, multiple / Yes, one / No, but plausible / No)
5. **Willingness to pay** — Would you pay to discover this automatically? (Strongly yes / Yes / Maybe / No / Strongly no)

### Interpretation

| Score Range | Meaning |
|-------------|---------|
| 20-25 | Critical value — builds product conviction |
| 15-19 | Strong value — ASF v2 is worth building |
| 10-14 | Moderate value — needs refinement |
| 5-9 | Low value — methodology needs rethinking |
| 0-4 | No value — pivot required |

### Target Outcomes

| Metric | Success | Stretch |
|--------|---------|---------|
| Mean willingness-to-pay score | >= 3.0/5 | >= 4.0/5 |
| % assumptions rated "Definitely investigate" | >= 50% | >= 70% |
| % participants citing incident history | >= 30% | >= 50% |
| Tier C accepted as "real" | >= 70% | >= 90% |

---

## Architecture 1: VPN → Internal App → Payroll DB

```
[User Laptop] --VPN--> [VPN Gateway] --TLS--> [Internal Web App] --SQL--> [Payroll Database (RDS)]
```

### Tier A Assumptions (everyone agrees — human, ASF, AI)

**A-001: VPN gateway enforces MFA for all remote users**

*Why it's important:* Password-only VPN is the #1 vector for credential-based breaches. Without MFA, a single phished password grants network access.

- Is this a real assumption? ___ / ___ / ___
- Would missing it create risk? ___ / ___ / ___ / ___ / ___
- Would you investigate it? ___ / ___ / ___ / ___ / ___
- Seen incidents caused by this? ___ / ___ / ___ / ___
- Would you pay to discover this? ___ / ___ / ___ / ___ / ___

Notes: ___________________________________________________________________

**A-002: Database is in private subnet with no internet route**

*Why it's important:* A "private" subnet with an implicit NAT gateway or internet route defeats the purpose of isolation. This is the second-most-common cloud misconfiguration after open S3 buckets.

- Is this a real assumption? ___ / ___ / ___
- Would missing it create risk? ___ / ___ / ___ / ___ / ___
- Would you investigate it? ___ / ___ / ___ / ___ / ___
- Seen incidents caused by this? ___ / ___ / ___ / ___
- Would you pay to discover this? ___ / ___ / ___ / ___ / ___

Notes: ___________________________________________________________________

### Tier C Assumptions (ASF-unique — no human or AI produced these)

**A-003: VPN gateway logs are monitored for brute-force attempts**

*Why ASF found this:* The Incident Response pattern in ASF's matrix triggers "if a system exists, detect attacks against it." Neither the human architect nor any AI persona listed VPN brute-force monitoring as an assumption — they assumed detection would happen but didn't make it explicit.

- Is this a real assumption? ___ / ___ / ___
- Would missing it create risk? ___ / ___ / ___ / ___ / ___
- Would you investigate it? ___ / ___ / ___ / ___ / ___
- Seen incidents caused by this? ___ / ___ / ___ / ___
- Would you pay to discover this? ___ / ___ / ___ / ___ / ___

Notes: ___________________________________________________________________

**A-004: VPN gateway is not a single point of failure**

*Why ASF found this:* The Availability & Resilience pattern triggers "if a component is critical, it must be redundant." The documented architecture shows a single VPN gateway — no human or AI flagged this as an assumption worth documenting.

- Is this a real assumption? ___ / ___ / ___
- Would missing it create risk? ___ / ___ / ___ / ___ / ___
- Would you investigate it? ___ / ___ / ___ / ___ / ___
- Seen incidents caused by this? ___ / ___ / ___ / ___
- Would you pay to discover this? ___ / ___ / ___ / ___ / ___

Notes: ___________________________________________________________________

---

## Architecture 2: Enterprise SSO → IdP → SAML Federation

```
[User Browser] --SAML--> [Okta IdP] --SAML Assertion--> [Service Provider Apps (x5)]
                              |
                         [AD Directory]
```

### Tier A Assumptions

**B-001: IdP is available and reachable from all SP applications**

*Why it's important:* IdP downtime blocks authentication to all connected applications. If Okta is unreachable, no one can access any service.

- Is this a real assumption? ___ / ___ / ___
- Would missing it create risk? ___ / ___ / ___ / ___ / ___
- Would you investigate it? ___ / ___ / ___ / ___ / ___
- Seen incidents caused by this? ___ / ___ / ___ / ___
- Would you pay to discover this? ___ / ___ / ___ / ___ / ___

Notes: ___________________________________________________________________

**B-002: MFA is enforced at every federated application, not just the IdP**

*Why it's important:* If an SP application bypasses IdP MFA requirements via direct authentication or a different IdP path, the MFA control is meaningless.

- Is this a real assumption? ___ / ___ / ___
- Would missing it create risk? ___ / ___ / ___ / ___ / ___
- Would you investigate it? ___ / ___ / ___ / ___ / ___
- Seen incidents caused by this? ___ / ___ / ___ / ___
- Would you pay to discover this? ___ / ___ / ___ / ___ / ___

Notes: ___________________________________________________________________

### Tier C Assumptions

**B-003: SAML metadata is refreshed and validated regularly**

*Why ASF found this:* The Dependency pattern triggers "if federation exists, metadata must be current." Stale metadata can cause certificate pinning failures, man-in-the-middle via rogue IdP, or application authentication outages.

- Is this a real assumption? ___ / ___ / ___
- Would missing it create risk? ___ / ___ / ___ / ___ / ___
- Would you investigate it? ___ / ___ / ___ / ___ / ___
- Seen incidents caused by this? ___ / ___ / ___ / ___
- Would you pay to discover this? ___ / ___ / ___ / ___ / ___

Notes: ___________________________________________________________________

**B-004: IdP administrator accounts are protected by hardware-backed MFA and not shared**

*Why ASF found this:* The Least Privilege pattern combines with Identity Lifecycle to examine privileged access to the identity system itself. A compromised IdP admin account can issue SAML assertions for any user to any application.

- Is this a real assumption? ___ / ___ / ___
- Would missing it create risk? ___ / ___ / ___ / ___ / ___
- Would you investigate it? ___ / ___ / ___ / ___ / ___
- Seen incidents caused by this? ___ / ___ / ___ / ___
- Would you pay to discover this? ___ / ___ / ___ / ___ / ___

Notes: ___________________________________________________________________

---

## Architecture 3: Microservices → Service Mesh → Kubernetes → Istio

```
[Ingress Gateway] --mTLS--> [Service A] --mTLS--> [Service B] --mTLS--> [Service C]
                    │              │                                       │
               [Istio Pilot]  [K8s API]                              [StatefulSet DB]
               [Citadel CA]  [etcd]                              [Persistent Volume]
```

### Tier A Assumptions

**C-001: mTLS certificates are rotated before expiry and revocation works**

*Why it's important:* Expired mTLS certificates break service-to-service communication. If the CA is unavailable or revocation lists are stale, compromised certificates remain trusted.

- Is this a real assumption? ___ / ___ / ___
- Would missing it create risk? ___ / ___ / ___ / ___ / ___
- Would you investigate it? ___ / ___ / ___ / ___ / ___
- Seen incidents caused by this? ___ / ___ / ___ / ___
- Would you pay to discover this? ___ / ___ / ___ / ___ / ___

Notes: ___________________________________________________________________

**C-002: Network policies are enforced at the data link layer, not just documented**

*Why it's important:* Kubernetes network policies that are defined but not enforced by a CNI plugin (or enforced incorrectly) provide no isolation. This is defense-in-theater.

- Is this a real assumption? ___ / ___ / ___
- Would missing it create risk? ___ / ___ / ___ / ___ / ___
- Would you investigate it? ___ / ___ / ___ / ___ / ___
- Seen incidents caused by this? ___ / ___ / ___ / ___
- Would you pay to discover this? ___ / ___ / ___ / ___ / ___

Notes: ___________________________________________________________________

### Tier C Assumptions

**C-003: Istio control plane (Pilot, Citadel) is itself secured, not just the data plane**

*Why ASF found this:* The Trust pattern triggers "if a component issues identity, its own security must be verified." The mesh secures service-to-service traffic but if the control plane is compromised, all identity is forged. No human or AI listed this.

- Is this a real assumption? ___ / ___ / ___
- Would missing it create risk? ___ / ___ / ___ / ___ / ___
- Would you investigate it? ___ / ___ / ___ / ___ / ___
- Seen incidents caused by this? ___ / ___ / ___ / ___
- Would you pay to discover this? ___ / ___ / ___ / ___ / ___

Notes: ___________________________________________________________________

**C-004: Sidecar proxy resource limits prevent DoS via resource exhaustion**

*Why ASF found this:* The Availability & Resilience pattern triggers "if a shared component exists, resource exhaustion must be prevented." Envoy sidecars share node resources — a memory-leaking sidecar can starve the application container. This is a production operational detail that architecture reviews routinely miss.

- Is this a real assumption? ___ / ___ / ___
- Would missing it create risk? ___ / ___ / ___ / ___ / ___
- Would you investigate it? ___ / ___ / ___ / ___ / ___
- Seen incidents caused by this? ___ / ___ / ___ / ___
- Would you pay to discover this? ___ / ___ / ___ / ___ / ___

Notes: ___________________________________________________________________

---

## Architecture 4: Healthcare → PHI → HIPAA Controls

```
[Patient Portal] --> [App Server] --> [PHI Database]
       │                    │
   [Auth0]            [Audit Logs] --> [SIEM]
```

### Tier A Assumptions

**D-001: PHI data is encrypted at rest using organization-controlled keys, not cloud-provider-default keys**

*Why it's important:* Default cloud-provider encryption keys may not satisfy HIPAA's requirements for customer-managed key control. If the provider can access PHI data, the BAA may not cover it.

- Is this a real assumption? ___ / ___ / ___
- Would missing it create risk? ___ / ___ / ___ / ___ / ___
- Would you investigate it? ___ / ___ / ___ / ___ / ___
- Seen incidents caused by this? ___ / ___ / ___ / ___
- Would you pay to discover this? ___ / ___ / ___ / ___ / ___

Notes: ___________________________________________________________________

**D-002: BAAs with all subprocessors are in place and up to date**

*Why it's important:* If a cloud provider, SIEM vendor, or Auth0 processes PHI data without a current BAA, the organization is out of HIPAA compliance regardless of application-layer controls.

- Is this a real assumption? ___ / ___ / ___
- Would missing it create risk? ___ / ___ / ___ / ___ / ___
- Would you investigate it? ___ / ___ / ___ / ___ / ___
- Seen incidents caused by this? ___ / ___ / ___ / ___
- Would you pay to discover this? ___ / ___ / ___ / ___ / ___

Notes: ___________________________________________________________________

### Tier C Assumptions

**D-003: Audit logs are tamper-proof and cannot be modified by the application or database administrators**

*Why ASF found this:* The Monitoring & Alerting pattern triggers "if evidence exists, it must be immutable." HIPAA requires audit controls that "record and examine" but does not explicitly require tamper-proofing. Administrators who can modify logs can cover up PHI access.

- Is this a real assumption? ___ / ___ / ___
- Would missing it create risk? ___ / ___ / ___ / ___ / ___
- Would you investigate it? ___ / ___ / ___ / ___ / ___
- Seen incidents caused by this? ___ / ___ / ___ / ___
- Would you pay to discover this? ___ / ___ / ___ / ___ / ___

Notes: ___________________________________________________________________

**D-004: Minimum necessary access is enforced at the application level, not just documented as policy**

*Why ASF found this:* The Least Privilege and Governance patterns combine to examine whether documented policies have corresponding enforcement mechanisms. A policy stating "minimum necessary access" without application-level RBAC or attribute-based access control is aspirational, not operational.

- Is this a real assumption? ___ / ___ / ___
- Would missing it create risk? ___ / ___ / ___ / ___ / ___
- Would you investigate it? ___ / ___ / ___ / ___ / ___
- Seen incidents caused by this? ___ / ___ / ___ / ___
- Would you pay to discover this? ___ / ___ / ___ / ___ / ___

Notes: ___________________________________________________________________

---

## Architecture 5: ERP → SOX → Financial Reporting → Audit

```
[Finance Team] --> [ERP Web App] --> [ERP Backend] --> [Financial DB]
       │                    │              │
   [Approval Workflow]  [Audit Logs]  [Reporting Engine] --> [Auditor Access]
```

### Tier A Assumptions

**E-001: Segregation of duties between transaction creation, approval, and reporting**

*Why it's important:* SOX requires that no single person can create, approve, and report a financial transaction. If the same user can do all three, the control is absent and SOX auditors will flag it.

- Is this a real assumption? ___ / ___ / ___
- Would missing it create risk? ___ / ___ / ___ / ___ / ___
- Would you investigate it? ___ / ___ / ___ / ___ / ___
- Seen incidents caused by this? ___ / ___ / ___ / ___
- Would you pay to discover this? ___ / ___ / ___ / ___ / ___

Notes: ___________________________________________________________________

**E-002: Audit log is immutable and covers all financial transactions**

*Why it's important:* SOX auditors require an unbroken chain of evidence. If the audit log can be truncated, filtered, or modified, the entire SOX control framework is undermined.

- Is this a real assumption? ___ / ___ / ___
- Would missing it create risk? ___ / ___ / ___ / ___ / ___
- Would you investigate it? ___ / ___ / ___ / ___ / ___
- Seen incidents caused by this? ___ / ___ / ___ / ___
- Would you pay to discover this? ___ / ___ / ___ / ___ / ___

Notes: ___________________________________________________________________

### Tier C Assumptions

**E-003: Journal entries cannot be deleted — only reversed with an audited reversing entry**

*Why ASF found this:* The Governance pattern triggers "if financial integrity is required, data deletion must be prevented." Standard ERP systems allow deletion of unposted journals. If a user can delete (rather than reverse) a journal entry, the audit trail is broken and SOX controls fail. No human or AI listed this.

- Is this a real assumption? ___ / ___ / ___
- Would missing it create risk? ___ / ___ / ___ / ___ / ___
- Would you investigate it? ___ / ___ / ___ / ___ / ___
- Seen incidents caused by this? ___ / ___ / ___ / ___
- Would you pay to discover this? ___ / ___ / ___ / ___ / ___

Notes: ___________________________________________________________________

**E-004: Auditor access is truly read-only and cannot be escalated**

*Why ASF found this:* The Trust + Least Privilege patterns combine to examine whether documented restrictions have enforcement gaps. "Read-only" auditor access is commonly implemented via a database user with SELECT only — but if the auditor can connect through the application, row-level security may not be enforced. If they connect directly, they may see all schemas.

- Is this a real assumption? ___ / ___ / ___
- Would missing it create risk? ___ / ___ / ___ / ___ / ___
- Would you investigate it? ___ / ___ / ___ / ___ / ___
- Seen incidents caused by this? ___ / ___ / ___ / ___
- Would you pay to discover this? ___ / ___ / ___ / ___ / ___

Notes: ___________________________________________________________________

---

## Scoring Summary Sheet

### Assumption Scores (to be filled after data collection)

| # | Assumption | Tier | Reality (1-3) | Risk (1-5) | Investigate (1-5) | Incidents (1-4) | Pay (1-5) | Total |
|---|-----------|------|--------------|-----------|------------------|----------------|----------|-------|
| A-001 | VPN MFA | A | | | | | | /22 |
| A-002 | Private subnet | A | | | | | | /22 |
| A-003 | VPN brute-force monitoring | C | | | | | | /22 |
| A-004 | VPN single point of failure | C | | | | | | /22 |
| B-001 | IdP availability | A | | | | | | /22 |
| B-002 | MFA at every app | A | | | | | | /22 |
| B-003 | SAML metadata refresh | C | | | | | | /22 |
| B-004 | IdP admin MFA | C | | | | | | /22 |
| C-001 | mTLS certificate rotation | A | | | | | | /22 |
| C-002 | Network policy enforcement | A | | | | | | /22 |
| C-003 | Control plane security | C | | | | | | /22 |
| C-004 | Sidecar resource limits | C | | | | | | /22 |
| D-001 | PHI encryption keys | A | | | | | | /22 |
| D-002 | BAA currency | A | | | | | | /22 |
| D-003 | Tamper-proof audit logs | C | | | | | | /22 |
| D-004 | Minimum necessary enforcement | C | | | | | | /22 |
| E-001 | Segregation of duties | A | | | | | | /22 |
| E-002 | Immutable audit trail | A | | | | | | /22 |
| E-003 | Journal entry deletion | C | | | | | | /22 |
| E-004 | Auditor read-only enforcement | C | | | | | | /22 |

### Participant Summary

| Question | Score |
|----------|-------|
| Mean Tier A "pay" score | ___ / 5 |
| Mean Tier C "pay" score | ___ / 5 |
| % assumptions rated "Definitely investigate" | ___ % |
| % participants reporting incident history on ≥1 assumption | ___ % |
| % Tier C accepted as "real" | ___ % |
| Overall mean total score (across all 20) | ___ / 22 |

### Participant Background

- Role: _________________________
- Years in security: _________
- Organization size: _________
- Domains worked: _________________________
- Current tooling for architecture review: _________________________

### Open-Ended Questions

1. What surprised you most about these assumptions?

2. Have you ever missed an assumption like this in a real review?

3. How do you currently validate assumptions in your architecture?

4. What would a tool need to do differently from SAST/DAST/Cloud Security Posture tools for you to adopt it?

5. Any other feedback?

---

## Facilitator Instructions

### Recruitment

Target 10-20 participants across these roles:
- 3-5 Security Architects
- 2-3 CISOs / Security Directors
- 2-3 Cloud Architects
- 2-3 AppSec Engineers
- 1-2 Compliance/Audit professionals

### Session Flow

1. **Briefing (5 min):** Explain the project. "We've built a framework that discovers hidden security assumptions from architecture. We want to know if the assumptions it finds are useful to practitioners like you."

2. **Architecture walkthrough (3 min each, 15 min total):** Show each architecture diagram. Read the documented policy. Contextualize.

3. **Assumption evaluation (1 min each, 20 min total):** For each assumption, ask the 5 questions. Take notes on any reactions ("that's obvious", "I've never thought of that", "we had an incident exactly like this").

4. **Debrief (5 min):** Open-ended questions. What surprised them? Would they use this? What's missing?

### What to Look For

**Positive signals:**
- Participant says "I've never thought of that" for a Tier C assumption
- Participant cites a specific incident matching an ASF-unique finding
- Participant asks "When can I get this tool?"
- Participant says "We spend hours in architecture reviews and still miss these"

**Negative signals:**
- Participant says "Everyone knows this" for Tier C (means it's not novel)
- Participant dismisses multiple assumptions as unrealistic
- Participant cannot see the value over existing tools

### Success Determination

Pass the study if:
- Mean pay score >= 3.0/5 across all participants
- >= 70% of Tier C accepted as real
- >= 30% of participants cite incident history on ≥1 ASF finding

Strong pass if:
- Mean pay score >= 4.0/5
- >= 90% of Tier C accepted as real
- Multiple participants independently say some variation of "I would use this"
