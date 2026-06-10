# Matching Rules for Independent Derivation Analysis

## Purpose

Define explicit, reproducible criteria for classifying whether an independently derived model assumption matches an ASF gold standard concept. Without these rules, the study is not reproducible and results are not comparable across models or raters.

---

## Classification Tiers

### Tier A — Full Match (Y)

Both ASF and the model independently identify the **same underlying security concern at the same level of specificity**.

**Criteria (ALL must be met):**
1. Same **security concern** — the threat, risk, or control objective is identical
2. Same **level of specificity** — both address the mechanism, not just the domain
3. Same **verification focus** — both identify what needs to be true, tested, or enforced

**Examples:**

| ASF Concept | Model Concept | Verdict | Rationale |
|------------|--------------|---------|-----------|
| Backup restore tested quarterly | Backup recovery testing is performed | Tier A | Both identify the need to test backups, not just have them |
| SAML assertion digital signature verification | SAML assertions are digitally signed | Tier A | Both identify signature verification as the mechanism |
| Audit trail integrity (tamper-proof logs) | Logs cannot be altered by attackers | Tier A | Same concern, same level of specificity |
| MFA enrollment enforced for all users | MFA is required for all users | Tier A | Both say MFA is mandatory, not optional |
| Network segmentation between tiers | The private subnet is actually isolated | Tier A | Both identify network isolation at the workload level |
| Segregation of duties enforcement | SoD rules are correctly defined and enforced | Tier A | Same compliance control, same specificity |

---

### Tier B — Partial Match (B)

Both ASF and the model identify concerns in the **same security domain**, but at **different levels of specificity** or from **different angles**.

**Criteria (at least ONE met, but not ALL of Tier A):**
1. Same **domain** (e.g., both address backup, IAM, or TLS) but different **mechanism**
2. One is **specific** and the other is **generic** (e.g., "quarterly restore testing" vs "backups exist")
3. One identifies a **verification method** while the other identifies a **general control objective**
4. The model's concept is **adjacent** — it addresses a neighboring concern but not the exact one

**Examples:**

| ASF Concept | Model Concept | Verdict | Rationale |
|------------|--------------|---------|-----------|
| Backup encryption key separation | Backup encryption at rest | Tier B | Both address backup encryption, but ASF specifies key separation (different key from data encryption); model says only "encrypted at rest" |
| SAML AudienceRestriction validation | SAML assertions are secure | Tier B | Both address SAML security, but ASF names the specific validation; model makes a generic statement |
| VPN brute-force monitoring | Rate limiting on VPN connections | Tier B | Both address VPN auth abuse, but monitoring vs rate limiting are different mechanisms |
| User Deprovisioning within 1 hour | User deprovisioning is timely | Tier B | Both address deprovisioning, but ASF specifies SLA; model is vague |
| KMS key rotation schedule | Encryption keys are properly managed | Tier B | Both address key management, but ASF specifies rotation cadence; model is generic |
| Backup storage in separate region | Backups are protected | Tier B | Both address backup protection, but ASF specifies geographic isolation; model is generic |
| HR system as identity source of truth | Directory synchronization is trustworthy | Tier B | Both address identity data accuracy, but ASF names the HR system as source; model addresses sync channel |

**Special case — Domain match but different mechanism:**
- ASF: "TLS 1.2+ enforcement for VPN Gateway"
- Model: "TLS is correctly configured between browser and VPN"
- Tier B: Both address TLS for the VPN path, but ASF specifies version; model addresses configuration correctness

---

### Tier C/N — No Match (N)

The model **did not independently generate** any concept that addresses the same security concern as the ASF concept. The assumption is entirely absent from the model's output.

**Criteria:**
1. No model assumption relates to the same security concern
2. The model has assumptions in the same domain (e.g., "networking") but none that touch this specific concern
3. The model has assumptions about the same component but for a different security property

**Examples:**

| ASF Concept | Closest Model Concept | Verdict | Rationale |
|------------|----------------------|---------|-----------|
| Vendor exit strategy | (none) | No Match | No model mentioned vendor lock-in, acquisition, or migration |
| MFA recovery process | (none) | No Match | Models say "MFA is required" but never ask about recovery when token is lost |
| Backup restore testing cadence | "Backup integrity is verified" | No Match | "Integrity" and "testing" are different concerns — integrity means not corrupted, testing means functional restore |
| VPN gateway single point of failure | "VPN gateway is hardened" | No Match | "Hardened" addresses configuration, not availability/SPOF |
| Data flow diagram documentation | "Network paths are documented" | No Match | Flow diagrams are about data movement; network paths are about connectivity — different concerns |
| KMS key policy restrictions | "Encryption keys are protected" | No Match | Key policy (who can use/decrypt) vs key protection (storage/access) — related but different security properties |
| Incident response plan for PHI breach with HIPAA 60-day notification | "Incident response plan exists" | No Match | Generic IR vs compliance-specific IR with regulatory timeline — specificity gap is too wide |

---

## Boundary Rules

### When a Generic Statement IS a Full Match

A generic model statement can be Tier A only if:
1. The ASF concept is itself generic (same level of abstraction)
2. The model's generic statement covers the EXACT same concern

**Example:** 
- ASF: "Database permissions follow least privilege" (Generic)
- Model: "Database accounts are restricted to minimum required access" (Generic)
- Tier A: Same concern, same abstraction level

### When a Specific ASF Concept IS a Partial Match with a Generic Model Statement

A specific ASF concept is Tier B (not N) when:
1. The model's generic statement would logically include the ASF concern IF someone asked a follow-up question
2. The model's statement AND the ASF concept belong to the same parent security domain

**Example:**
- ASF: "Backup tested quarterly via restore drill" (Specific)
- Model: "Backup integrity is verified" (Generic)
- Tier B: Both address backup verification. The model didn't specify testing cadence, but "integrity is verified" implies some form of validation. A follow-up question like "how verified?" would likely elicit restore testing.

### When a Generic Model Statement is NOT Even a Partial Match

A generic model statement is No Match when:
1. The ASF concept addresses a concern in an entirely different sub-domain
2. The model's statement is so generic that it could apply to any component and doesn't specifically address the ASF concern

**Example:**
- ASF: "VPN vendor backdoor risk assessment" (Vendor/Supply chain)
- Model: "VPN gateway software is not vulnerable" (Technical vulnerability)
- No Match: Vendor backdoor risk and software vulnerability are different threat categories. "Not vulnerable" doesn't cover supply chain trust.

---

## Semantic Similarity Guidelines

The following transformations are considered **semantically equivalent** (Tier A):

| Transformation | Example |
|---------------|---------|
| Passive ↔ Active voice | "MFA must be enforced" ↔ "Administration enforces MFA" |
| Positive ↔ Negative form | "TLS must be enabled" ↔ "No unencrypted connections permitted" |
| MUST ↔ assumes | "Backups MUST be tested" ↔ "Assumes backups are tested" |
| Component synonym | "Load balancer" ↔ "Application load balancer" ↔ "ALB" |
| Action synonym | "Restrict" ↔ "Limit" ↔ "Control" ↔ "Isolate" |
| Compliance synonym | "HIPAA" ↔ "Healthcare privacy" ↔ "PHI protection" |

The following transformations are **NOT equivalent** (Tier B or N):

| Non-Equivalent Pair | Gap |
|--------------------|-----|
| "Backup encryption" vs "Backup key separation" | Key management vs encryption |
| "MFA required" vs "MFA recovery process" | Enforcement vs recovery |
| "Incident response plan" vs "IR plan with 60-day breach notification" | Generic vs compliance-specific |
| "Network segmentation" vs "Micro-segmentation workload isolation" | Perimeter vs workload-level |
| "Logs cannot be altered" vs "Log retention for 7 years (SOX 802)" | Immutability vs retention duration |
| "Dependency scanning" vs "SBOM with fail-on-critical" | Tooling vs policy/enforcement |

---

## Per-Domain Specificity Rules

Some domains consistently show a specificity gap. These rules codify when a model's generic statement in that domain counts as a match:

**Backup & Recovery**
- If ASF says "tested quarterly" and model says "backup exists": **Tier B** (both address backup existence, ASF adds cadence)
- If ASF says "tested quarterly" and model says "backup integrity verified": **Tier B** (both address some verification)
- If ASF says "key separation" and model says "encrypted": **Tier B** (both address encryption, ASF adds key mgmt)
- If ASF says "cross-region storage" and model says "protected": **Tier B**

**Identity Lifecycle**
- If ASF says "deprovisioning within 1 hour" and model says "deprovisioning is timely": **Tier B**
- If ASF says "provisioning within 24 hours" and model says "provisioning is accurate": **Tier B**
- If ASF says "recertification quarterly" and model says "role reviews occur": **Tier B**

**Incident Response**
- If ASF says "IR plan with HIPAA 60-day timeline" and model says "IR plan exists": **Tier B**
- If ASF says "forensic evidence preservation" and model says "logs are retained": **Tier B** (evidence vs logs)
- If ASF says "component isolation procedure" and model says "containment plan exists": **Tier B**

**Third-party Dependency**
- If ASF says "vendor exit strategy" and model says "vendor security is assessed": **No Match** — different concepts
- If ASF says "SLA dependency" and model says "third-party availability matters": **Tier B**

**Availability & Resilience**
- If ASF says "single point of failure" and model says "multi-AZ deployment": **Tier A** — same concept, different framing
- If ASF says "offline procedure" and model says "DR plan exists": **Tier B**

---

## Decision Flowchart

```
Is the security concern the same?
├── YES → Is the specificity level the same?
│   ├── YES → Tier A (Full Match)
│   └── NO → Tier B (Partial Match - same concern, different detail)
└── NO → Is the concern in the same domain with some overlap?
    ├── YES → Tier B (Partial Match - adjacent concept)
    └── NO → Tier C/N (No Match)
```

## Reproducibility Requirements

To ensure reproducibility:

1. **Two independent raters** must classify each ASF-concept vs model-assumption pair
2. **Disagreements** resolved by third rater or consensus discussion
3. **Tier B examples** maintained as precedent for future matching
4. **Ambiguous cases** documented with rationale (not silently classified)
5. **Rules reviewed** after every 50 classifications to catch drift

---

## Appendix: Rejected Classification Approaches

**"Bag of words" overlap** — Two concepts matching because they share keywords like "backup" and "encryption" is not sufficient. The security concern must match.

**"If the model would agree if asked"** — Hypothetical follow-up questions do not count. Only what the model actually said in its independent derivation.

**"Negative evidence"** — The absence of a contradiction does not constitute a match. "Model didn't say encryption is unnecessary" is not evidence that model matches "encryption is required."

**"Lower bound only"** — Tier A is a lower bound on convergence. If the model has a more specific version of the same concept (model says "SAML AudienceRestriction validated by each SP" and ASF says "SAML assertion validation"), it is still Tier A — specificity in either direction is acceptable.
