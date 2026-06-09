# ASF v2 Hidden Assumption Discovery — Experiment Kit

---

## Participant Briefing

### What This Study Is About

Security architects routinely design systems with implicit **hidden assumptions** — things that *must* remain true for the architecture to stay secure, but that nobody writes down. When these assumptions break, security breaks.

This study measures how well a structured methodology (the Assumption Security Framework, or ASF) can surface those hidden assumptions, compared to what experienced architects identify on their own.

### Why It Matters

Most security incidents trace back to an assumption that was silently violated:
- "We assumed the VPN would always be available."
- "We assumed the database was only reachable from the app server."
- "We assumed secrets would be rotated on time."

By making hidden assumptions visible, we can monitor them, test them, and fix them *before* they cause a breach. Your participation helps us build a repeatable method for doing this systematically.

### What You Will Be Asked to Do (45 minutes total)

1. **Review architecture descriptions** — You will be shown 2-3 architecture diagrams described in text. Each is a realistic system topology with a documented security policy.
2. **List hidden assumptions** — For each architecture, you will list every assumption that must be true for the architecture to remain secure, but that is **not** written in the documented policy.
3. **Answer brief follow-up questions** — After the exercise, you will answer 4 short questions about your experience.

### Data Use & Privacy

- Your responses are collected **anonymously**. No personally identifying information is recorded.
- Your participant ID is a random code generated for this session only.
- Raw data will be stored encrypted and accessible only to the research team.
- Published results will contain aggregate statistics only — no individual responses will be identifiable.

### Important Note

**There are no wrong answers.** We are interested in your genuine mental model — how *you* think about security assumptions. If you're unsure whether something counts, include it. We would rather have too many assumptions than too few.

---

## Instructions

### The Task

For each architecture you are given:

> **List every assumption that must be true for this architecture to remain secure. Do NOT list what is already documented in the policy. List only what must remain true but is NOT written down.**

### What Kinds of Assumptions?

Think about:
- **Dependencies** — What external systems or services must behave correctly?
- **Operational practices** — What must operations teams do (or not do)?
- **Environmental conditions** — What must be true about the network, the cloud, the physical world?
- **Human behavior** — What must users, developers, or administrators do (or not do)?
- **Failure modes** — What must never happen, and what must happen when something breaks?
- **Time-sensitive conditions** — What must stay true over time (patching, rotation, certificates)?

### Format

You may use:
- Free-form text
- Bullet points
- Numbered lists
- Any format that is comfortable for you

### Target Volume

Aim for **20-50 assumptions per architecture**. Quality matters, but breadth matters too — hidden assumptions are, by definition, hard to see.

### Time Limit

Spend **no more than 45 minutes total** across all 2-3 architectures.

### A Note on Scope

Focus on assumptions that relate to **security** of the architecture. If something would break confidentiality, integrity, or availability if violated, it belongs on the list. If it's purely a performance or cost concern, it's out of scope.

---

## Architecture Selection

### Recommended Order (Simple → Complex)

The architectures below are ordered by increasing complexity. For first-time participants, we recommend following this order to build familiarity with the task. Each session should cover **2-3 architectures**.

| # | Architecture | Complexity |
|---|--------------|------------|
| 1 | User → VPN → Internal App → DB | Simple |
| 2 | Web App → LB → App Server → RDS | Simple |
| 3 | Mobile App → API → Lambda → DynamoDB | Simple |
| 4 | SaaS → IdP → SSO → SAML Federation | Medium |
| 5 | Microservices → Mesh → K8s → Istio | Medium |
| 6 | E-commerce → Payment → PCI Scope | Medium |
| 7 | Multi-region → Active/Passive → DR | Medium |
| 8 | CI/CD → Artifact → Registry → Deploy | Medium |
| 9 | Vendor SaaS → API → Internal Systems | Medium |
| 10 | Data Pipeline → Kafka → S3 → Redshift | Medium |
| 11 | Healthcare → PHI → HIPAA Controls | Complex |
| 12 | Fintech → Ledger → SOX Controls | Complex |
| 13 | Partner B2B → Federation → API Exchange | Complex |
| 14 | Hybrid Cloud → VPN → Direct Connect | Complex |
| 15 | IoT → MQTT → Gateway → Cloud Backend | Complex |
| 16 | ML Pipeline → Training → Serving → Data | Complex |
| 17 | Multi-tenant SaaS → Tenant Isolation | Complex |
| 18 | Global CDN → WAF → Origin → DB | Complex |
| 19 | Secrets → Vault → App → Rotation | Complex |
| 20 | ERP → SOX → Fin Reporting → Audit | Complex |

### Reviewer Note

You may pick any architecture in any order. If you have a specific domain expertise (e.g., healthcare, fintech, IoT), choose architectures from that domain. The study design accommodates any subset.

---

## Quick Reference: Architecture Summaries

---

### Architecture 1: User → VPN → Internal App → Payroll DB

**Topology:**
```
[User Laptop] --VPN--> [VPN Gateway] --TLS--> [Internal Web App] --SQL--> [Payroll Database (RDS)]
```

**Documented Policy:**
- VPN required for remote access
- Application authenticates with AD credentials
- Database is in private subnet
- Backups run nightly

**Trust Boundaries:**
- User ↔ VPN (authentication boundary)
- VPN ↔ Application (network boundary)
- Application ↔ Database (data boundary)

---

### Architecture 2: Web App → Load Balancer → App Server → RDS

**Topology:**
```
[Browser] --HTTPS--> [ALB] --HTTPS--> [EC2 App Server (x3)] --SQL--> [RDS Primary + Replica]
                                    --Logs--> [CloudWatch]
```

**Documented Policy:**
- TLS termination at ALB
- Auto-scaling based on CPU
- RDS automated backups enabled
- Application logs sent to CloudWatch

**Trust Boundaries:**
- Browser to ALB (internet boundary)
- ALB to App Server (internal network boundary)
- App Server to RDS (data boundary)

---

### Architecture 3: Mobile App → API Gateway → Lambda → DynamoDB

**Topology:**
```
[Mobile App] --HTTPS--> [API Gateway] --Event--> [Lambda Function (xN)] --SDK--> [DynamoDB Table]
                                          └--> [Cognito User Pool] (Auth)
```

**Documented Policy:**
- API Gateway requires Cognito auth
- Lambda uses least-privilege IAM role
- DynamoDB is encrypted at rest
- API keys per mobile app version

**Trust Boundaries:**
- Mobile to API Gateway (auth boundary)
- API Gateway to Lambda (service boundary)
- Lambda to DynamoDB (data boundary)

---

### Architecture 4: Enterprise SSO → IdP → SAML Federation

**Topology:**
```
[User Browser] --SAML--> [Okta IdP] --SAML Assertion--> [Service Provider Apps (x5)]
                              |
                         [AD Directory] (User store)
```

**Documented Policy:**
- All apps require SSO via Okta
- MFA enforced for all users
- Session timeout after 8 hours
- JIT provisioning enabled

**Trust Boundaries:**
- Browser to IdP (auth boundary)
- IdP to SP (federation trust boundary)
- IdP to AD (directory sync boundary)

---

### Architecture 5: Microservices → Service Mesh → Kubernetes → Istio

**Topology:**
```
[Ingress Gateway] --mTLS--> [Service A] --mTLS--> [Service B] --mTLS--> [Service C]
                    │              │                                       │
               [Istio Pilot]  [K8s API]                              [StatefulSet DB]
                    │              │                                       │
               [Citadel CA]  [etcd]                                  [Persistent Volume]
```

**Documented Policy:**
- mTLS enabled between all services
- RBAC enforced at namespace level
- Pod security policies restrict privileged containers
- Network policies isolate namespaces

**Trust Boundaries:**
- Ingress to Services (mesh boundary)
- Service to Service (identity boundary)
- Service to CA (certificate trust boundary)
- Pod to K8s API (control plane boundary)

---

### Architecture 6: E-commerce → Payment Processor → PCI Scope

**Topology:**
```
[Browser] --HTTPS--> [Web App] --API--> [Payment Processor (Stripe)]
                    [PCI Scope]
                         |
                   [Token Vault]
                         |
                   [Order DB]
```

**Documented Policy:**
- Payment tokens used, not raw PAN
- PCI DSS compliant environment
- Encryption at rest and in transit
- Quarterly vulnerability scans

**Trust Boundaries:**
- Browser to Web App (PCI boundary)
- Web App to Payment Processor (third-party boundary)
- Token Vault access (privilege boundary)

---

### Architecture 7: Multi-Region → Active/Passive → Disaster Recovery

**Topology:**
```
[Route53] --Failover--> [Region A (Active)]         [Region B (Passive)]
                              │                              │
                         [App + DB]                   [App + DB (Replica)]
                              │                              │
                         [S3 (Primary)] <--Cross-Region--> [S3 (Replica)]
```

**Documented Policy:**
- RTO: 4 hours
- RPO: 15 minutes
- Cross-region DB replication enabled
- Route53 health checks configured

**Trust Boundaries:**
- Between Region A and B (geo boundary)
- Active to Passive promotion (state boundary)
- Cross-region replication (network boundary)

---

### Architecture 8: CI/CD Pipeline → Artifact Registry → Deploy

**Topology:**
```
[Developer] --> [GitHub] --> [CI (GitHub Actions)] --> [Artifact Registry (ECR)]
                                                            │
                                                     [CD (ArgoCD)] --> [K8s Cluster]
```

**Documented Policy:**
- Code review required before merge
- CI runs security scan on every commit
- Images are signed before deployment
- ArgoCD syncs from Git as source of truth

**Trust Boundaries:**
- Developer to GitHub (identity boundary)
- CI to Registry (pipeline boundary)
- Registry to Cluster (deployment boundary)

---

### Architecture 9: Vendor SaaS → API → Internal Systems

**Topology:**
```
[Vendor SaaS (Salesforce)] --> [API Gateway] --> [Internal CRM Sync] --> [Customer DB]
                                    │
                              [OAuth Token]
```

**Documented Policy:**
- Vendor has SOC 2 Type II report
- API integration uses OAuth 2.0
- Data processing agreement in place
- Quarterly vendor review

**Trust Boundaries:**
- Vendor to API Gateway (external boundary)
- API Gateway to Internal (trust boundary)
- Data residency (compliance boundary)

---

### Architecture 10: Data Pipeline → Kafka → S3 → Redshift

**Topology:**
```
[Producer Apps] --Produce--> [Kafka Cluster] --Consume--> [Spark Job] --Write--> [S3 Data Lake]
                                                                                    │
                                                                             [Redshift] <--[Analyst Queries]
```

**Documented Policy:**
- Data encrypted at rest in S3 and Redshift
- Kafka configured with TLS + SASL
- IAM roles used for cross-service access
- Data retention of 7 years

**Trust Boundaries:**
- Producer to Kafka (ingestion boundary)
- Kafka to Spark (processing boundary)
- S3 to Redshift (data access boundary)

---

### Architecture 11: Healthcare → PHI → HIPAA Controls

**Topology:**
```
[Patient Portal] --> [App Server] --> [PHI Database]
       │                    │
   [Auth0]            [Audit Logs] --> [SIEM]
```

**Documented Policy:**
- PHI encrypted at rest (AES-256)
- BAAs with all subprocessors
- Access logging enabled
- Minimum necessary access enforced

**Trust Boundaries:**
- Portal to Application (auth boundary)
- Application to Database (PHI boundary)
- Application to SIEM (audit boundary)

---

### Architecture 12: Fintech → Ledger → SOX Controls

**Topology:**
```
[User] --> [Trading App] --> [Ledger Service] --> [Accounting DB]
                │                      │
          [Market Data API]      [Audit Trail Service] --> [Immutable Audit DB]
```

**Documented Policy:**
- SOX controls on all financial transactions
- Segregation of duties between trading and settlement
- Audit trail is immutable
- Quarterly access recertification

**Trust Boundaries:**
- Trading to Ledger (transaction boundary)
- Ledger to Accounting (ledger boundary)
- Application to Audit Trail (evidence boundary)

---

### Architecture 13: Partner B2B → Federation → API Exchange

**Topology:**
```
[Partner A IdP] ---SAML---> [API Gateway] ---OAuth---> [Partner A Resources]
[Partner B IdP] ---SAML---> [API Gateway] ---OAuth---> [Partner B Resources]
```

**Documented Policy:**
- SAML federation with partners
- API tokens scoped to specific resources
- Rate limiting per partner
- Quarterly metadata refresh

**Trust Boundaries:**
- Partner IdP to API Gateway (federation boundary)
- API Gateway to Resource (access boundary)
- Partner to Partner (isolation boundary)

---

### Architecture 14: Hybrid Cloud → VPN → Direct Connect

**Topology:**
```
[On-Prem DC] --Direct Connect--> [AWS VPC]
                │                       │
           [Corporate FW]         [Transit Gateway]
                                       │
                              [Prod, Staging, Dev VPCs]
```

**Documented Policy:**
- Direct Connect is primary, VPN is backup
- On-prem to cloud routing through TGW
- Security groups restrict cross-environment traffic
- Flow logs enabled on all VPCs

**Trust Boundaries:**
- On-Prem to Cloud (hybrid network boundary)
- VPC to VPC (environment boundary)
- Direct Connect to TGW (routing boundary)

---

### Architecture 15: IoT → MQTT → Gateway → Cloud Backend

**Topology:**
```
[IoT Device] --MQTT/TLS--> [IoT Gateway] --Stream--> [Kinesis] --> [Processing Lambda] --> [TimeSeries DB]
                  │
            [Device Registry / CA]
```

**Documented Policy:**
- Device certificates for authentication
- TLS 1.2 minimum for MQTT
- Device registry manages identity
- Data encrypted at rest

**Trust Boundaries:**
- Device to Gateway (device identity boundary)
- Gateway to Cloud (network boundary)
- Device Registry (certificate trust boundary)

---

### Architecture 16: ML Pipeline → Training → Serving → Data Lineage

**Topology:**
```
[Feature Store] --> [Training Job (SageMaker)] --> [Model Registry] --> [Serving Endpoint]
       │                                                                    │
  [Data Lake (S3)]                                                    [Production App]
```

**Documented Policy:**
- Training data is from approved sources
- Models are versioned in registry
- Inference endpoint requires auth
- Data lineage is tracked

**Trust Boundaries:**
- Data Lake to Training (data integrity boundary)
- Training to Registry (model integrity boundary)
- Registry to Endpoint (model serving boundary)

---

### Architecture 17: Multi-tenant SaaS → Tenant Isolation

**Topology:**
```
[Tenant A] --> [API Gateway] --> [App Service] --> [Tenant A DB]
[Tenant B] --> [API Gateway] --> [App Service] --> [Tenant B DB]
                (Shared)              (Shared)        (Isolated)
```

**Documented Policy:**
- Tenant data is isolated at database level
- Row-level security enforced
- Tenant context propagated through JWT claims
- Cross-tenant access is blocked

**Trust Boundaries:**
- Tenant to Gateway (auth boundary)
- Shared Service to Isolated DB (isolation boundary)
- Tenant to Tenant (data boundary)

---

### Architecture 18: Global CDN → WAF → Origin → Database

**Topology:**
```
[Global Users] --> [CloudFront CDN] --> [WAF] --> [ALB] --> [EC2 Origin] --> [RDS]
                      │
                 [Lambda@Edge] (Auth)
```

**Documented Policy:**
- WAF blocks common attack patterns
- CDN caches static content
- Origin only accessible via CDN (not directly)
- DDoS protection enabled

**Trust Boundaries:**
- CDN to WAF (edge boundary)
- WAF to Origin (origin access boundary)
- Origin to Database (data boundary)

---

### Architecture 19: Secrets → Vault → Application → Rotation

**Topology:**
```
[App Pod] --mTLS--> [Vault Agent Sidecar] --API--> [Vault Server] --> [KMS (Unseal)]
                                                         │
                                                    [Database] (stores secrets)
```

**Documented Policy:**
- Secrets retrieved via Vault API
- Dynamic secrets with TTL
- Vault audit log enabled
- Auto-unseal via KMS

**Trust Boundaries:**
- App to Vault (secrets access boundary)
- Vault to KMS (unseal boundary)
- Vault to Storage (data boundary)

---

### Architecture 20: ERP → SOX → Financial Reporting → Audit

**Topology:**
```
[Finance Team] --> [ERP Web App] --> [ERP Backend] --> [Financial DB]
       │                    │              │
   [Approval Workflow]  [Audit Logs]  [Reporting Engine] --> [Auditor Access]
```

**Documented Policy:**
- SOX controls on all journal entries
- Segregation of duties enforced
- Read-only access for auditors
- Quarterly recertification

**Trust Boundaries:**
- User to ERP (auth boundary)
- Approval workflow (segregation boundary)
- ERP to Reporting (data integrity boundary)
- Auditor Access (read-only boundary)

---

## Scoring Sheet

Print one copy per architecture reviewed.

```
┌─────────────────────────────────────────────────────┐
│              ASSUMPTION DISCOVERY SCORING            │
├─────────────────────────────────────────────────────┤
│                                                      │
│  Architecture #:  ___                                │
│                                                      │
│  Architecture Name: _______________________________  │
│                                                      │
│  Participant ID:  ___                                │
│                                                      │
│  Date:  ____________________                         │
│                                                      │
│  Start Time:  ______          End Time:  ______      │
│                                                      │
│  ────────────────────────────────────────────────    │
│                                                      │
│  Total assumptions listed:  ___                      │
│                                                      │
│  Participant confidence (1-5):  ___                  │
│  (1 = guessing, 5 = very confident)                  │
│                                                      │
│  Notes:                                              │
│  _________________________________________________  │
│  _________________________________________________  │
│  _________________________________________________  │
│  _________________________________________________  │
│                                                      │
└─────────────────────────────────────────────────────┘
```

### Facilitator Scoring (post-session, for research use)

```
┌─────────────────────────────────────────────────────┐
│              FACILITATOR SCORING (RESEARCH ONLY)     │
├─────────────────────────────────────────────────────┤
│                                                      │
│  Unique assumptions (H):  _____                      │
│                                                      │
│  Unique ASF assumptions (A):  _____                  │
│                                                      │
│  Overlap (O):  _____                                 │
│                                                      │
│  Precision (O/A):  _____%                            │
│                                                      │
│  Recall (O/H):  _____%                               │
│                                                      │
│  Novel findings (A-O):  _____                        │
│                                                      │
│  Missed findings (H-O):  _____                       │
│                                                      │
│  F1 Score:  _____%                                   │
│                                                      │
└─────────────────────────────────────────────────────┘
```

---

## Assumption Recording Sheet

Print multiple copies per architecture. Use one row per assumption.

```
┌─────┬─────────────────────────────────────────────────────────────┬──────────────────┬──────────────┐
│  #  │                        Assumption                          │   Category       │  Confidence  │
│     │  (must be true for architecture to remain secure,           │  (see key below) │  (1-5,       │
│     │   NOT in documented policy)                                 │                  │   1=guess)   │
├─────┼─────────────────────────────────────────────────────────────┼──────────────────┼──────────────┤
│  1  │                                                             │                  │              │
├─────┼─────────────────────────────────────────────────────────────┼──────────────────┼──────────────┤
│  2  │                                                             │                  │              │
├─────┼─────────────────────────────────────────────────────────────┼──────────────────┼──────────────┤
│  3  │                                                             │                  │              │
├─────┼─────────────────────────────────────────────────────────────┼──────────────────┼──────────────┤
│  4  │                                                             │                  │              │
├─────┼─────────────────────────────────────────────────────────────┼──────────────────┼──────────────┤
│  5  │                                                             │                  │              │
├─────┼─────────────────────────────────────────────────────────────┼──────────────────┼──────────────┤
│  6  │                                                             │                  │              │
├─────┼─────────────────────────────────────────────────────────────┼──────────────────┼──────────────┤
│  7  │                                                             │                  │              │
├─────┼─────────────────────────────────────────────────────────────┼──────────────────┼──────────────┤
│  8  │                                                             │                  │              │
├─────┼─────────────────────────────────────────────────────────────┼──────────────────┼──────────────┤
│  9  │                                                             │                  │              │
├─────┼─────────────────────────────────────────────────────────────┼──────────────────┼──────────────┤
│ 10  │                                                             │                  │              │
├─────┼─────────────────────────────────────────────────────────────┼──────────────────┼──────────────┤
│ 11  │                                                             │                  │              │
├─────┼─────────────────────────────────────────────────────────────┼──────────────────┼──────────────┤
│ 12  │                                                             │                  │              │
├─────┼─────────────────────────────────────────────────────────────┼──────────────────┼──────────────┤
│ 13  │                                                             │                  │              │
├─────┼─────────────────────────────────────────────────────────────┼──────────────────┼──────────────┤
│ 14  │                                                             │                  │              │
├─────┼─────────────────────────────────────────────────────────────┼──────────────────┼──────────────┤
│ 15  │                                                             │                  │              │
├─────┼─────────────────────────────────────────────────────────────┼──────────────────┼──────────────┤
│ 16  │                                                             │                  │              │
├─────┼─────────────────────────────────────────────────────────────┼──────────────────┼──────────────┤
│ 17  │                                                             │                  │              │
├─────┼─────────────────────────────────────────────────────────────┼──────────────────┼──────────────┤
│ 18  │                                                             │                  │              │
├─────┼─────────────────────────────────────────────────────────────┼──────────────────┼──────────────┤
│ 19  │                                                             │                  │              │
├─────┼─────────────────────────────────────────────────────────────┼──────────────────┼──────────────┤
│ 20  │                                                             │                  │              │
└─────┴─────────────────────────────────────────────────────────────┴──────────────────┴──────────────┘

                     ┌──────────────┬─────────────────────────────────────┐
                     │   Category   │         Description                 │
                     │     Key      │                                     │
                     ├──────────────┼─────────────────────────────────────┤
                     │     AUTH     │ Authentication & identity           │
                     │     NET      │ Network security & segmentation     │
                     │     DATA     │ Data protection & encryption        │
                     │     IAM      │ Access control & entitlements       │
                     │     OPS      │ Operational & change mgmt           │
                     │     DEP      │ Dependencies & third parties        │
                     │     MON      │ Monitoring & incident response      │
                     │     RES      │ Resilience & availability           │
                     │     PHY      │ Physical & environmental            │
                     │     COMP     │ Compliance & regulatory             │
                     └──────────────┴─────────────────────────────────────┘
```

---

## Post-Session Questions

Please answer after completing all architecture reviews.

---

**Q1:** What was the hardest part of this exercise?

```
___________________________________________________________________________
___________________________________________________________________________
___________________________________________________________________________
___________________________________________________________________________
```

---

**Q2:** Did any of the assumptions you listed surprise you? If so, which ones?

```
___________________________________________________________________________
___________________________________________________________________________
___________________________________________________________________________
___________________________________________________________________________
```

---

**Q3:** Would you normally document these assumptions in your day-to-day work? Why or why not?

```
___________________________________________________________________________
___________________________________________________________________________
___________________________________________________________________________
___________________________________________________________________________
```

---

**Q4:** What percentage of the assumptions you listed do you think are actually verified or monitored in production?

```
    0%    10%    20%    30%    40%    50%    60%    70%    80%    90%    100%
    ┌──────┬──────┬──────┬──────┬──────┬──────┬──────┬──────┬──────┬──────┐
    │      │      │      │      │      │      │      │      │      │      │
    └──────┴──────┴──────┴──────┴──────┴──────┴──────┴──────┴──────┴──────┘

    Actual percentage:  _____%
```

---

**Q5 (optional):** Any other feedback about the exercise or the ASF methodology?

```
___________________________________________________________________________
___________________________________________________________________________
___________________________________________________________________________
___________________________________________________________________________
```

---

## Facilitator Checklist

```
Before Session:
  ☐ Print architecture summaries (pages ___ – ___)
  ☐ Print scoring sheet(s) — 1 per architecture
  ☐ Print assumption recording sheets — 1+ per architecture
  ☐ Print post-session questions
  ☐ Assign participant ID
  ☐ Prepare timer (45 minutes)

During Session:
  ☐ Read briefing aloud or provide to participant
  ☐ Hand participant instructions
  ☐ Hand first architecture summary
  ☐ Record start time
  ☐ After 15 min: check progress, offer to move to next architecture
  ☐ After 30 min: remind of remaining time
  ☐ At 45 min: collect sheets
  ☐ Distribute post-session questions

After Session:
  ☐ Transfer assumptions to digital format
  ☐ Run ASF generator for each architecture reviewed
  ☐ Compute overlap, precision, recall, F1
  ☐ Log results in experiment tracker
```

---

*End of Experiment Kit*
