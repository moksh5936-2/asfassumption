# ASF v2 Discovery Experiment Protocol

## Purpose
Validate whether the ASF Assumption Generator Matrix produces a set of hidden
assumptions that agrees with human security architects at >= 70% overlap.

## Procedure

### Step 1: Select Architecture
Pick one from the 20 reference architectures (see `architecture_patterns.md`).
Each is a realistic system diagram described in text.

### Step 2: Human Session (45 min)

Give one architecture to a security engineer. Ask:

> "List every assumption that must be true for this architecture to remain
> secure. Do not list what's documented. List what must remain true."

Collect their assumptions as a flat list.

### Step 3: ASF Methodology

Run the assumption generator matrix for the same architecture:

1. Identify which of the 20 patterns apply to this architecture
2. For each applicable pattern, collect the derived assumptions
3. Merge into a single flat list with duplicates removed

### Step 4: Compare

| Metric | Formula |
|--------|---------|
| Human assumptions | H = count of unique human-generated assumptions |
| ASF assumptions | A = count of unique ASF-generated assumptions |
| Overlap | O = count of assumptions appearing in both lists |
| Precision | O / A |
| Recall | O / H |
| F1 | 2 * (P * R) / (P + R) |
| Novel findings | A - O (assumptions ASF found that human missed) |
| Missed findings | H - O (assumptions human found that ASF missed) |

### Step 5: Target

| Metric | Target |
|--------|--------|
| Recall | >= 70% (ASF captures most human-identified assumptions) |
| Precision | >= 50% (not excessive false positives) |
| Novel findings | >= 10% of total (ASF adds value beyond human) |
| Missed findings | <= 30% of human total (acceptable gap) |

### Step 6: Iterate

After each architecture, review missed findings:
- Are they due to missing patterns? Add to assumption generator matrix.
- Are they architecture-specific? Document as edge case.
- Are they false positives in human list? Note for next comparison.

Target: After 10 architectures, recall >= 70% should be stable.

## Scoring Sheet

```
Architecture: ___________________________________
Date: ___________________________________________
Reviewer: _______________________________________

Human assumptions (H):   _____
ASF assumptions (A):     _____
Overlap (O):             _____
Precision (O/A):         _____%
Recall (O/H):            _____%
Novel findings (A-O):    _____
Missed findings (H-O):   _____
```

## Architecture Order (recommended)

| # | Architecture | Complexity | Patterns |
|---|---|---|---|
| 1 | User → VPN → Internal App → DB | Simple | MFA, VPN, DB access, encryption |
| 2 | Web App → LB → App Server → RDS | Simple | TLS, WAF, backup, monitoring |
| 3 | Mobile App → API → Lambda → DynamoDB | Simple | API auth, IAM, encryption, monitoring |
| 4 | SaaS → IdP → SSO → SAML Federation | Medium | SSO, federation, provisioning |
| 5 | Microservices → Mesh → K8s → Istio | Medium | mTLS, RBAC, network policy, secrets |
| 6 | E-commerce → Payment → PCI Scope | Medium | PCI, encryption, logging, dependency |
| 7 | Multi-region → Active/Passive → DR | Medium | Failover, backup, RTO, RPO |
| 8 | CI/CD → Artifact → Registry → Deploy | Medium | Supply chain, signing, scan |
| 9 | Vendor SaaS → API → Internal Systems | Medium | Third-party, SLA, data handling |
| 10 | Data Pipeline → Kafka → S3 → Redshift | Medium | Encryption, IAM, retention, classification |
| 11 | Healthcare → PHI → HIPAA Controls | Complex | PHI, encryption, audit, BAA |
| 12 | Fintech → Ledger → SOX Controls | Complex | SOX, audit, change management, segregation |
| 13 | Partner B2B → Federation → API Exchange | Complex | Federation, trust, token, SLA |
| 14 | Hybrid Cloud → VPN → Direct Connect | Complex | Network segmentation, routing, encryption |
| 15 | IoT → MQTT → Gateway → Cloud Backend | Complex | Device identity, certificate, protocol |
| 16 | ML Pipeline → Training → Serving → Data | Complex | Data lineage, model security, access |
| 17 | Multi-tenant SaaS → Tenant Isolation | Complex | Isolation, RBAC, data segregation |
| 18 | Global CDN → WAF → Origin → DB | Complex | CDN, WAF, TLS, DDoS, origin security |
| 19 | Secrets → Vault → App → Rotation | Complex | Secrets lifecycle, rotation, access |
| 20 | ERP → SOX → Fin Reporting → Audit | Complex | SOX, segregation, change audit, review |
