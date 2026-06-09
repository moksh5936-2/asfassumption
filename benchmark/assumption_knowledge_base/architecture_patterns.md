# 20 Reference Architectures for Assumption Discovery

Each architecture describes:
- **Topology**: the system components and their relationships
- **Documented Policy**: what is stated/written
- **Trust Boundaries**: where assumptions live
- **Prompt**: what to analyze

---

## 1. User → VPN → Internal App → Payroll DB

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
- Between User and VPN (auth boundary)
- Between VPN and Application (network boundary)
- Between Application and Database (data boundary)

**Assumptions to consider:**
- What if VPN is unavailable?
- What if AD authentication is down?
- What if database credentials are leaked?
- What if backup restore is untested?
- What if payroll DB has a public route?
- What if MFA is not enforced on VPN?

---

## 2. Web App → Load Balancer → App Server → RDS

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

## 3. Mobile App → API Gateway → Lambda → DynamoDB

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

## 4. Enterprise SSO → IdP → SAML Federation

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

## 5. Microservices → Service Mesh → Kubernetes → Istio

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

## 6. E-commerce → Payment Processor → PCI Scope

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

## 7. Multi-Region → Active/Passive → Disaster Recovery

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

## 8. CI/CD Pipeline → Artifact Registry → Deploy

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

## 9. Vendor SaaS → API → Internal Systems

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

## 10. Data Pipeline → Kafka → S3 → Redshift

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

## 11. Healthcare → PHI → HIPAA Controls

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

## 12. Fintech → Ledger → SOX Controls

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

## 13. Partner B2B → Federation → API Exchange

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

## 14. Hybrid Cloud → VPN → Direct Connect

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

## 15. IoT → MQTT → Gateway → Cloud Backend

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

## 16. ML Pipeline → Training → Serving → Data Lineage

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

## 17. Multi-tenant SaaS → Tenant Isolation

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

## 18. Global CDN → WAF → Origin → Database

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

## 19. Secrets → Vault → Application → Rotation

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

## 20. ERP → SOX → Financial Reporting → Audit

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
