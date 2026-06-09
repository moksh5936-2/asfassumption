# Independent Derivation Test

**Purpose:** Stronger AI-only validation. Instead of asking "Does the AI agree?", ask "Does the AI independently derive the same assumption?" — a blind, independent reproduction test.

**Method:** For each of 5 architectures, present ONLY the architecture diagram, documented policy, and trust boundaries to GPT-4o, Claude 4 Sonnet, and Gemini 2.5 Pro independently. Do NOT show ASF output. Collect their assumptions. Compare overlap with ASF predictions.

---

## Protocol

### Per-Architecture Steps

1. Copy the architecture prompt (next section) into a clean chat with the AI
2. Capture the full output — all assumptions generated
3. Map each output assumption to the ASF prediction list using the comparison matrix
4. Score: Full match / Partial match / Not found / ASF-unique

### Example Output Mapping

| AI-Generated Assumption | ASF Prediction | Match |
|-------------------------|----------------|-------|
| "VPN requires MFA" | "VPN gateway enforces MFA for all remote users" | Full match |
| "DB is isolated" | "Database is in private subnet with no internet route" | Full match |
| — | "VPN gateway logs are monitored for brute-force" | ASF-unique |

### Independent = No Contamination

Each architecture must be run in a **fresh, clean chat**. Do not tell the AI:
- This is about assumption discovery
- You are comparing against another framework
- You are validating ASF output
- What other AIs answered

The prompt must appear to be a standalone architecture review request.

Use the exact prompts below verbatim.

---

## Architectures

---

## Architecture 1: VPN → Internal App → Payroll DB

### Prompt

```
You are a security architect reviewing the following system.

SYSTEM TOPOLOGY:
[User Laptop] --VPN--> [VPN Gateway] --TLS--> [Internal Web App] --SQL--> [Payroll Database (RDS)]

DOCUMENTED POLICY:
- VPN required for remote access
- Application authenticates with AD credentials
- Database is in private subnet
- Backups run nightly

TRUST BOUNDARIES:
- Between User and VPN (auth boundary)
- Between VPN and Application (network boundary)
- Between Application and Database (data boundary)

Please list ALL implicit assumptions that this architecture depends on for security.
An implicit assumption is something that must be true for the system to be secure, but is not explicitly stated in the documented policy.

For each assumption, explain:
1. What the assumption is
2. Why it matters for security
3. What would happen if the assumption is false

List as many as you can identify.
```

---

## Architecture 2: Enterprise SSO → IdP → SAML Federation

### Prompt

```
You are a security architect reviewing the following system.

SYSTEM TOPOLOGY:
[User Browser] --SAML--> [Okta IdP] --SAML Assertion--> [Service Provider Apps (x5)]
                              |
                         [AD Directory] (User store)

DOCUMENTED POLICY:
- All apps require SSO via Okta
- MFA enforced for all users
- Session timeout after 8 hours
- JIT provisioning enabled

TRUST BOUNDARIES:
- Browser to IdP (auth boundary)
- IdP to SP (federation trust boundary)
- IdP to AD (directory sync boundary)

Please list ALL implicit assumptions that this architecture depends on for security.
An implicit assumption is something that must be true for the system to be secure, but is not explicitly stated in the documented policy.

For each assumption, explain:
1. What the assumption is
2. Why it matters for security
3. What would happen if the assumption is false

List as many as you can identify.
```

---

## Architecture 3: Microservices → Service Mesh → Kubernetes → Istio

### Prompt

```
You are a security architect reviewing the following system.

SYSTEM TOPOLOGY:
[Ingress Gateway] --mTLS--> [Service A] --mTLS--> [Service B] --mTLS--> [Service C]
                    │              │                                       │
               [Istio Pilot]  [K8s API]                              [StatefulSet DB]
                    │              │                                       │
               [Citadel CA]  [etcd]                                  [Persistent Volume]

DOCUMENTED POLICY:
- mTLS enabled between all services
- RBAC enforced at namespace level
- Pod security policies restrict privileged containers
- Network policies isolate namespaces

TRUST BOUNDARIES:
- Ingress to Services (mesh boundary)
- Service to Service (identity boundary)
- Service to CA (certificate trust boundary)
- Pod to K8s API (control plane boundary)

Please list ALL implicit assumptions that this architecture depends on for security.
An implicit assumption is something that must be true for the system to be secure, but is not explicitly stated in the documented policy.

For each assumption, explain:
1. What the assumption is
2. Why it matters for security
3. What would happen if the assumption is false

List as many as you can identify.
```

---

## Architecture 4: Healthcare → PHI → HIPAA Controls

### Prompt

```
You are a security architect reviewing the following system.

SYSTEM TOPOLOGY:
[Patient Portal] --> [App Server] --> [PHI Database]
       │                    │
   [Auth0]            [Audit Logs] --> [SIEM]

DOCUMENTED POLICY:
- PHI encrypted at rest (AES-256)
- BAAs with all subprocessors
- Access logging enabled
- Minimum necessary access enforced

TRUST BOUNDARIES:
- Portal to Application (auth boundary)
- Application to Database (PHI boundary)
- Application to SIEM (audit boundary)

Please list ALL implicit assumptions that this architecture depends on for security.
An implicit assumption is something that must be true for the system to be secure, but is not explicitly stated in the documented policy.

For each assumption, explain:
1. What the assumption is
2. Why it matters for security
3. What would happen if the assumption is false

List as many as you can identify.
```

---

## Architecture 5: ERP → SOX → Financial Reporting → Audit

### Prompt

```
You are a security architect reviewing the following system.

SYSTEM TOPOLOGY:
[Finance Team] --> [ERP Web App] --> [ERP Backend] --> [Financial DB]
       │                    │              │
   [Approval Workflow]  [Audit Logs]  [Reporting Engine] --> [Auditor Access]

DOCUMENTED POLICY:
- SOX controls on all journal entries
- Segregation of duties enforced
- Read-only access for auditors
- Quarterly recertification

TRUST BOUNDARIES:
- User to ERP (auth boundary)
- Approval workflow (segregation boundary)
- ERP to Reporting (data integrity boundary)
- Auditor Access (read-only boundary)

Please list ALL implicit assumptions that this architecture depends on for security.
An implicit assumption is something that must be true for the system to be secure, but is not explicitly stated in the documented policy.

For each assumption, explain:
1. What the assumption is
2. Why it matters for security
3. What would happen if the assumption is false

List as many as you can identify.
```

---

## Comparison Framework

### Per-Architecture Mapping

For each AI output, map every assumption to the nearest ASF prediction using this table. Mark:
- **F** = Full match (same intent, same details)
- **P** = Partial match (same general area but different specifics)
- **N** = Not found (AI listed it, ASF did not)
- **U** = ASF-unique (ASF predicted it, AI did not)

### Comparison Matrices

Use the following tables for mapping.

#### Architecture 1: VPN → Payroll DB (Arch 001)

| # | ASF Prediction | GPT | Claude | Gemini |
|---|---------------|-----|--------|--------|
| 1 | VPN gateway enforces MFA for all remote users | | | |
| 2 | Database is in private subnet with no internet route | | | |
| 3 | Backups are encrypted at rest | | | |
| 4 | Application uses parameterized queries to prevent SQLi | | | |
| 5 | User accounts follow joiner/mover/leaver process | | | |
| 6 | Database credentials are rotated regularly | | | |
| 7 | VPN gateway logs are monitored for brute-force attempts | | | |
| 8 | VPN gateway is not a single point of failure | | | |
| 9 | Backup restores are tested | | | |
| 10 | AD authentication server is redundant | | | |

#### Architecture 2: SSO/IdP Federation (Arch 004)

| # | ASF Prediction | GPT | Claude | Gemini |
|---|---------------|-----|--------|--------|
| 1 | IdP is available and reachable from all SP applications | | | |
| 2 | MFA is enforced at every federated application, not just the IdP | | | |
| 3 | SAML assertions are signed to prevent tampering | | | |
| 4 | Session timeout (8hr) is enforced across all SP apps | | | |
| 5 | JIT provisioning does not create orphaned accounts | | | |
| 6 | SAML metadata is refreshed and validated regularly | | | |
| 7 | IdP administrator accounts are protected by hardware-backed MFA | | | |
| 8 | AD to IdP sync is timely and does not drift | | | |
| 9 | SP-initiated SSO is disabled or controlled | | | |
| 10 | Compromised IdP cannot issue assertions for unauthorized apps | | | |

#### Architecture 3: K8s/Istio Service Mesh (Arch 005)

| # | ASF Prediction | GPT | Claude | Gemini |
|---|---------------|-----|--------|--------|
| 1 | mTLS certificates are rotated before expiry | | | |
| 2 | Network policies are enforced at data link layer, not just documented | | | |
| 3 | RBAC is enforced at pod level, not just namespace level | | | |
| 4 | Pod security policies are actively enforced, not audited | | | |
| 5 | mTLS identity is bound to workload identity, not just node identity | | | |
| 6 | Istio control plane (Pilot, Citadel) is itself secured | | | |
| 7 | Sidecar proxy resource limits prevent DoS via resource exhaustion | | | |
| 8 | etcd is encrypted at rest | | | |
| 9 | K8s API is not exposed to mesh-internal traffic | | | |
| 10 | Certificate revocation works before expiry | | | |

#### Architecture 4: Healthcare/PHI/HIPAA (Arch 011)

| # | ASF Prediction | GPT | Claude | Gemini |
|---|---------------|-----|--------|--------|
| 1 | PHI encryption uses organization-controlled keys, not provider-default | | | |
| 2 | BAAs with all subprocessors are current | | | |
| 3 | Access logging covers all PHI access, not just auth events | | | |
| 4 | Minimum necessary access is enforced at application level, not policy | | | |
| 5 | PHI data is classified and tagged in the database | | | |
| 6 | Audit logs are tamper-proof (cannot be modified by admins) | | | |
| 7 | PHI is masked in non-production environments | | | |
| 8 | SIEM is monitored, not just collecting logs | | | |
| 9 | Data retention and purging schedules exist for PHI | | | |
| 10 | Authorization server (Auth0) does not log PHI in debug/error logs | | | |

#### Architecture 5: ERP/SOX/Audit (Arch 020)

| # | ASF Prediction | GPT | Claude | Gemini |
|---|---------------|-----|--------|--------|
| 1 | Segregation of duties between create, approve, and report | | | |
| 2 | Audit log is immutable and covers all financial transactions | | | |
| 3 | Read-only access for auditors is enforced at DB level | | | |
| 4 | Quarterly recertification is enforced, not just documented | | | |
| 5 | Approval workflow cannot be bypassed | | | |
| 6 | Journal entries cannot be deleted — only reversed with audited reversal | | | |
| 7 | Auditor access is read-only and cannot be escalated | | | |
| 8 | Reporting engine does not cache sensitive financial data | | | |
| 9 | Approval workflow has escalation path for reviewer unavailability | | | |
| 10 | ERP backend validates authorization on every request (not just frontend) | | | |

---

## Scoring

### Per-AI Metrics

| Metric | Calculation |
|--------|------------|
| Full match rate | F count / 10 |
| ASF-unique rate | U count / 10 |
| Novel findings | N count (AI found something ASF missed) |
| Overlap efficacy | (F + P) / total distinct assumptions |

### Cross-AI Consensus

| Tier | Definition | Count |
|------|------------|-------|
| A | All 3 AIs independently derived it | |
| B | 2 of 3 AIs independently derived it | |
| C | 1 of 3 AIs independently derived it | |
| D | 0 of 3 AIs independently derived it (ASF-unique across all 3) | |

### Target Outcomes (vs Multi-LLM Campaign)

| Metric | Prior Campaign (Agreement) | Independent Test |
|--------|---------------------------|------------------|
| Tier A (all agree) | ~69.1% validated | ___ |
| Tier D (ASF-unique vs all AI) | ~31.6% (Tier C in prior) | ___ |
| Novel findings (AI found something ASF missed) | Not measured | ___ |

---

## Execution Log

| Arch | GPT-4o | Claude 4 Sonnet | Gemini 2.5 Pro | Date |
|------|--------|-----------------|----------------|------|
| 1. VPN → Payroll | | | | |
| 2. SSO/IdP | | | | |
| 3. K8s/Istio | | | | |
| 4. Healthcare/PHI | | | | |
| 5. ERP/SOX | | | | |

---

## Instructions for Running

1. Open a **fresh, clean chat** for each architecture × AI combination (15 total runs)
2. Paste the exact architecture prompt (do NOT modify, do NOT add context)
3. Wait for full response
4. Copy the response into a new file: `independent_derivation_{arch_num}_{ai_name}.md`
5. Fill in the comparison matrix above
6. Calculate per-architecture and aggregate metrics

Do NOT tell the AI this is a comparison test. Each session must appear to be a standalone architecture security review.
