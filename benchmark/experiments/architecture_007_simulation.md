# ASF Phase 6 Experiment: Architecture #007

**Architecture:** Multi-Region → Active/Passive → Disaster Recovery
**Date:** 2026-06-09
**Simulation Mode:** Human + ASF framework comparison

---

## Architecture Profile

```
[Route53] --Failover--> [Region A (Active)]     [Region B (Passive)]
                              │                              │
                         [App + DB]                   [App + DB (Replica)]
                              │                              │
                         [S3 (Primary)] <--X-Region--> [S3 (Replica)]
```

### Documented Policy
| # | Policy |
|---|--------|
| P1 | RTO: 4 hours |
| P2 | RPO: 15 minutes |
| P3 | Cross-region DB replication enabled |
| P4 | Route53 health checks configured |

### Trust Boundaries
| Boundary | Type |
|----------|------|
| Region A ↔ Region B | Geo boundary |
| Active ↔ Passive promotion | State boundary |
| Cross-region replication | Network boundary |

### Complexity Rating
**Moderate** — multi-region topology, 8 nodes, 3 trust boundaries, cross-region state management.

---

## Step 1: Human-Generated Assumptions

*Role: Senior Security Architect. Listed below are assumptions that MUST remain true for this architecture to remain secure but are NOT stated in the documented policy.*

| ID | Assumption | Justification |
|----|-----------|---------------|
| H-001 | Route53 health checks accurately reflect application and database availability in both regions. | A health check reporting a healthy endpoint that is actually serving errors causes Route53 to continue sending traffic to a degraded region. |
| H-002 | DNS failover timeout is within the RTO of 4 hours, not hours longer due to TTL propagation delays. | Long DNS TTLs on Route53 records can extend failover time beyond the stated RTO, violating the SLA. |
| H-003 | Application state (sessions, user data) is either replicated or can be safely lost on failover. | Active/passive failover drops in-memory sessions in the active region; users must re-authenticate or risk data inconsistency. |
| H-004 | Database replication lag never exceeds the 15-minute RPO under normal and peak load. | Replication lag spikes during write-heavy periods can push RPO beyond 15 minutes, causing data loss on failover. |
| H-005 | Cross-region S3 replication has no gaps — all objects written to the primary bucket are replicated within SLA. | S3 replication is eventually consistent; incomplete replication during failover means the passive region has stale data. |
| H-006 | The passive region has sufficient capacity (compute, memory, DB connections) to handle full production load immediately upon promotion. | A passive region running with reduced capacity will be overwhelmed on failover, causing extended downtime or degraded service. |
| H-007 | S3 cross-region replication handles object deletions correctly — delete markers from the primary are replicated. | If delete operations are not replicated, the passive region retains deleted objects (data leak); if they are, malicious deletion propagates. |
| H-008 | Route53 failover is tested at least annually with a full cut-over exercise. | Untested failover is the leading cause of disaster recovery failures; configuration drift makes the documented procedure unreliable. |
| H-009 | The application in the passive region can connect to the replica database and function correctly. | Application code or configuration may hard-code region-specific endpoints that do not work after failover. |
| H-010 | Database credentials, application secrets, and encryption keys are available in the passive region before failover. | Secrets replicated to the passive region may be stale, out of sync, or not replicated at all. |
| H-011 | KMS keys in the passive region can decrypt data from the active region (or a multi-region key is used). | Cross-region encrypted data requires KMS key access in the failover region; key policy misconfiguration blocks decryption. |
| H-012 | The Route53 health check endpoint does not expose sensitive application internals. | A health check URL that leaks stack traces, database status, or configuration data is a reconnaissance vector. |
| H-013 | The database in the passive region accepts read traffic during normal operation without impacting replication. | Read traffic on the passive replica before promotion can increase replication lag and reduce RPO compliance. |
| H-014 | Automated promotion scripts are idempotent and error-handled — partial promotion does not leave the system in a split-brain state. | A failed promotion that activates both regions simultaneously causes data corruption from dual writers. |
| H-015 | Monitoring and alerting systems are also deployed in the passive region and survive region failover. | If monitoring does not fail over with the application, the operations team is blind during the most critical period. |
| H-016 | Route53 DNS changes propagate quickly enough that clients are not stuck resolving to the failed region. | Client-side DNS caching beyond Route53 TTL causes clients to hit the failed region for extended periods. |
| H-017 | TLS certificates for the passive region are valid and not expired at the time of failover. | Expired certificates in the passive region cause TLS handshake failures and block all application traffic. |
| H-018 | The database in the failed region does not corrupt data during an unplanned failover (crash consistency). | A hard crash of the primary database can produce a corrupt final state that propagates to the replica. |
| H-019 | Network ACLs and security groups in the passive region are identical to the active region in allow rules. | A missing security group rule in the passive region blocks application traffic exactly when it is needed most. |
| H-020 | IAM roles and policies work identically in both regions — no region-specific ARN or service endpoint mismatches. | IAM roles in Region B with different trust policies or resource ARNs break application and service connectivity after failover. |
| H-021 | Logs from the active region are replicated or accessible in the passive region for forensic analysis after failover. | Post-incident investigation requires logs from the failed region; unreplicated logs are lost on region failure. |
| H-022 | Route53 weighted routing or latency-based policies do not interfere with the failover record. | A misconfigured routing policy that sends traffic to both regions simultaneously defeats the active/passive design. |
| H-023 | Database replication encryption in transit is enforced (TLS between primary and replica). | Unencrypted cross-region replication exposes database contents to the underlying network provider or intermediaries. |
| H-024 | S3 replication IAM role has least-privilege permissions — read from primary, write to replica, nothing else. | An over-permissioned replication role can be used as a pivot to access other S3 buckets or services. |
| H-025 | Failback (passive → active) is as thoroughly tested as failover, including data reconciliation. | Organizations test failover but not failback; the return-to-normal is often more dangerous than the failover itself. |
| H-026 | Application auto-scaling policies in the passive region trigger correctly during failover to handle surge traffic. | A passive region that does not auto-scale under load will be overwhelmed by the traffic redirected to it. |
| H-027 | Database connection strings in the application use a DNS name (Route53 or internal ALB) not a hard-coded IP or AZ endpoint. | Hard-coded database endpoints in application config prevent seamless failover to the replica. |
| H-028 | The SIEM and security monitoring pipeline receives logs from both regions continuously, not just the active one. | If SIEM only ingests from the active region, a failover event creates a monitoring gap. |
| H-029 | Cross-region replication of S3 objects respects the same encryption and access control policies as the primary. | Objects written to the replica bucket with different KMS keys or ACLs than intended may be inaccessible or over-exposed. |
| H-030 | Route53 health checks are configured for depth (checking application, database, and dependencies, not just the web server). | A shallow health check that only pings the web server does not detect a failed database or upstream service. |
| H-031 | The application in the passive region does not need to write to the primary S3 bucket after failover. | If the promoted region writes to the original primary S3 bucket instead of the local replica, you create a split-brain scenario. |
| H-032 | No single-AZ dependencies exist in either region that would become SPOFs after failover. | A passive region built on a single AZ fails if that AZ is the one experiencing the outage. |
| H-033 | Data at rest encryption keys in the passive region are not subject to a different key rotation schedule. | A key that was rotated in the active region but not in the passive region leaves the replica inaccessible. |
| H-034 | Route53 failover record has a low TTL (60 seconds or less) to enable fast DNS propagation. | Default Route53 TTLs of 300 seconds or more add minutes to the failover time beyond the health check detection. |
| H-035 | The team has run a tabletop exercise for the failover scenario within the last 6 months, not just a technical test. | Technical failover tests do not validate communication chains, decision authority, and manual approval steps. |
| H-036 | CloudTrail and S3 access logs from both regions are aggregated in a single monitorable location. | Split logging across regions makes unified threat detection and post-incident investigation infeasible. |
| H-037 | No direct peering or VPN exists between the two regions that creates an unintended network path for data exfiltration. | Cross-region VPC peering intended for replication but left open may allow an attacker in one region to access resources in the other. |
| H-038 | The application code does not contain region-specific branching logic that behaves incorrectly during failover. | Region-specific conditionals in application code (e.g., "if us-east-1, use endpoint X") fail silently in the other region. |
| H-039 | Database failover scripts account for the replication lag time and do not trigger promotion until lag is within RPO. | Promoting a replica that is hours behind the primary because of lag violates the documented RPO. |
| H-040 | The passive region has a working internet egress path (NAT Gateway, Internet Gateway) for outbound API calls. | Applications that need to call external APIs (Stripe, Slack, monitoring) will fail if the passive region lacks egress. |
| H-041 | Route53 failover action is also triggerable manually and is not fully automated without human approval. | Fully automated failover that triggers on a false-positive health check causes unnecessary downtime and data risk. |

**Total (H): 41**

---

## Step 2: ASF-Generated Assumptions

*Using the 20-pattern assumption generator matrix. Applicable patterns: 15 of 20. Patterns excluded: Container Security (no containers), Endpoint Security (not client-facing), Physical Security (cloud-hosted), Supply Chain Security (deferred to Third-party Dependency).*

---

### Pattern 1: Authentication (MFA)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-001 | Route53 administrative console and API access requires MFA for all users who can modify failover records. | Explicit | Route53 configuration changes can disable failover entirely; MFA on the management plane is critical. |
| ASF-002 | MFA recovery codes for AWS root/IAM users are not stored in a way that is accessible from the passive region alone. | Derived | If the active region is down, recovery codes stored there are inaccessible, creating an identity deadlock. |
| ASF-003 | Break-glass MFA bypass procedures exist for emergency failover when the primary IdP is in the failed region. | Operational | If the IdP is in the failed region, administrators cannot authenticate to promote the passive region. |
| ASF-004 | IAM users assumed through cross-account roles in the passive region also require MFA. | Implicit | Cross-account roles may inherit MFA requirements inconsistently across regions. |

---

### Pattern 2: Authentication (SSO)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-005 | The corporate SSO IdP is deployed in both regions or has a failover configuration independent of the application regions. | Architectural | An IdP co-located in the active region fails with it, preventing all administrative access during failover. |
| ASF-006 | SSO session tokens issued in the active region are valid in the passive region after failover. | Trust | Users authenticated in the active region must be able to access the promoted passive region without re-authentication. |
| ASF-007 | Service accounts for Route53 health checks are not dependent on SSO availability. | Dependency | Health checks that authenticate via SSO will fail during an IdP outage, creating a cascade failure. |
| ASF-008 | SSO token signing certificates are replicated to the passive region or a multi-region certificate is used. | Operational | If signing keys are only in the active region, token validation in the passive region fails. |

---

### Pattern 3: Availability & Resilience

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-009 | Route53 is inherently multi-region and does not itself fail — its control plane is available even if one region fails. | Explicit | Route53 is a global service with its own availability guarantees, but configuration management access depends on the AWS console being reachable. |
| ASF-010 | The passive region does not share a single failure domain (power grid, network backbone, cloud provider control plane) with the active region. | Architectural | Two regions in the same geographic area (e.g., us-east-1 and us-east-2) can share upstream infrastructure failures. |
| ASF-011 | The cross-region replication link has sufficient bandwidth to meet RPO under peak load. | Environmental | Replication bandwidth saturation during peak writes causes lag beyond the 15-minute RPO target. |
| ASF-012 | Application startup time in the passive region plus DNS propagation does not exceed the 4-hour RTO. | Derived | Cold-starting applications and propagating DNS changes consume time from the RTO budget before traffic arrives. |
| ASF-013 | A single-region failure of the cloud provider's KMS service does not block decryption in the remaining region. | Environmental | KMS is regional; if the active region's KMS fails, encrypted data in the passive region must be decryptable with local keys. |
| ASF-014 | There is a documented runbook for manual failover when Route53 health checks fail to trigger automatically. | Operational | Automation failure requires human intervention; an untested runbook delays failover and increases data loss. |

---

### Pattern 4: Backup & Recovery

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-015 | Cross-region backups are stored in a physically separate geographic location from both primary and replica. | Explicit | Backup data stored in the same region or the same provider as the primary is vulnerable to the same catastrophe. |
| ASF-016 | Backup restore in the passive region is tested independently — the passive DB is not just a replica but has restorable backups. | Derived | Replication and backup serve different purposes; a corrupt schema propagates through replication but not from backup. |
| ASF-017 | Database snapshots are replicated to the passive region on a schedule shorter than the RPO (e.g., every 5 minutes for 15-min RPO). | Implicit | If snapshots are hourly and RPO is 15 minutes, backup restore can never meet the RPO target. |
| ASF-018 | Backup retention policies keep copies in both regions so that a data corruption in one does not destroy all recoverable points. | Operational | A retention policy that only stores backups in the active region loses all recovery points when that region fails. |

---

### Pattern 5: Change Management

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-019 | Changes to Route53 failover configuration require approval from both the application and infrastructure teams. | Explicit | A single-team change to routing rules or health check endpoints can silently break failover. |
| ASF-020 | DR configuration changes (health checks, replication settings, failover scripts) are tested in a non-production environment before deployment. | Operational | Untested DR configuration changes are the root cause of most failover failures during actual incidents. |
| ASF-021 | Changes to IAM policies or KMS key policies in one region are replicated to the other region as part of the change process. | Derived | IAM and KMS changes are regional; applying them only in the active region leaves the passive region with stale permissions. |
| ASF-022 | A configuration drift detection process exists to compare Route53, security group, and IAM configurations between regions. | Operational | Drift between regions accumulates over time; without detection, the passive region silently diverges from the active. |

---

### Pattern 6: Cloud Security (IAM)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-023 | IAM roles used for cross-region replication are scoped to exactly the source and destination buckets, not all S3 buckets. | Explicit | Over-scoped replication roles can be used to exfiltrate any S3 object in the account. |
| ASF-024 | The instance profile for EC2 in the passive region has identical permissions to the active region profile. | Derived | Different instance profile ARNs in the two regions can cause application permission failures after promotion. |
| ASF-025 | Automated failover scripts use an IAM role (not long-lived access keys) that is valid in both regions. | Implicit | Long-lived access keys with failover permissions outside of IAM role rotation create a standing privilege risk. |
| ASF-026 | The AWS root account for each region is protected by hardware MFA and used only for break-glass scenarios. | Trust | Root account credentials that work across all regions are the ultimate backdoor; weak protection on root undermines all IAM controls. |

---

### Pattern 7: Container Security

*Not applicable — no containers in this architecture.*

---

### Pattern 8: Data Flow & Classification

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-027 | Data classification determines which data is replicated cross-region — not all data may need DR protection. | Explicit | Replicating all data, including ephemeral or non-critical data, wastes resources and increases the cross-region data transfer attack surface. |
| ASF-028 | Data flow diagrams exist for the failover state as well as the steady state. | Implicit | Data flows during failover (when the passive region becomes active) are different from steady-state flows and may have undocumented paths. |
| ASF-029 | No data is written only to the active region without replication — no data loss beyond RPO on failover. | Derived | Application components that write data to a local file system not in the database or S3 lose that data on failover. |
| ASF-030 | Cross-region data transfer complies with data residency requirements (data does not leave approved geographies). | Environmental | Replicating data from a European region to a US region may violate GDPR or local data residency laws. |

---

### Pattern 9: Encryption at Rest

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-031 | Both regions use encryption at rest for all data stores (RDS, S3, EBS) with KMS-managed keys. | Explicit | Encryption at rest may be enabled in the active region but overlooked in the passive region or replica resources. |
| ASF-032 | KMS keys in the passive region are not separate keys but are managed as multi-region keys or have a cross-region import process. | Derived | Separate KMS keys in each region mean that data encrypted with the active region key cannot be decrypted in the passive region. |
| ASF-033 | S3 cross-region replication preserves server-side encryption status — objects remain encrypted with the configured key. | Trust | Replication that decrypts objects in the source and re-encrypts in the destination creates a window where plaintext exists in transit. |
| ASF-034 | KMS key rotation does not break decryption of older S3 objects or database backups during failover. | Operational | Key rotation that disables old key versions renders archived data in the passive region permanently unreadable. |

---

### Pattern 10: Encryption in Transit (TLS)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-035 | Cross-region database replication uses TLS encryption between the primary and replica databases. | Explicit | Unencrypted cross-region replication exposes all database writes to any entity with network access to the replication link. |
| ASF-036 | Cross-region S3 replication is encrypted in transit (SSL/TLS between regions). | Derived | S3 replication uses HTTPS by default, but configuration options may allow plaintext transfer for performance. |
| ASF-037 | Application-to-application TLS between regions (if any) uses certificates validated against a common CA. | Trust | Self-signed or region-specific CA certificates cause TLS handshake failures during cross-region communication after failover. |
| ASF-038 | TLS 1.2 or higher is enforced on all cross-region communication channels, including replication endpoints. | Derived | Cross-region links using older TLS are vulnerable to downgrade attacks on the inter-region network. |

---

### Pattern 11: Endpoint Security

*Not applicable in this architecture — no user endpoints in scope.*

---

### Pattern 12: Human Factors & Process

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-039 | Engineers performing failover exercises understand the difference between active and passive region configuration. | Derived | A mistake during a failover exercise (e.g., writing to the wrong region) can cause a real production incident. |
| ASF-040 | On-call engineers have access to both region's AWS consoles and can authenticate even if the primary region is down. | Operational | If the on-call engineer's IdP or MFA device depends on the failed region, they cannot respond. |
| ASF-041 | DR runbooks are stored in a location accessible from both regions — not just in the active region's wiki or S3 bucket. | Implicit | A runbook stored in the active region S3 bucket is unavailable when that region fails. |
| ASF-042 | Incident response team members do not make unaudited ad-hoc changes during failover without following the runbook. | Trust | Ad-hoc changes during the stress of a failover event are the most common source of post-failover configuration errors. |

---

### Pattern 13: Identity Lifecycle (Provisioning)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-043 | IAM user and role lifecycle management is synchronized across both regions — disabled in one means disabled in both. | Operational | A terminated employee's IAM access in Region B remains active if identity deprovisioning is not cross-region. |
| ASF-044 | Cross-account roles used for DR are reviewed and recertified quarterly. | Derived | Stale cross-account trust relationships create standing privilege that bypasses region boundaries. |
| ASF-045 | Service accounts used by the DR automation scripts are not shared between teams or environments. | Implicit | A shared service account used for DR automation lacks accountability and creates a massive blast radius. |
| ASF-046 | DR team membership is current — former team members do not retain administrative access to Route53 or failover scripts. | Operational | A team member who changed roles but still has Route53 administrative access could accidentally or maliciously trigger failover. |

---

### Pattern 14: Incident Response

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-047 | The incident response plan includes a specific playbook for multi-region failover scenarios. | Operational | A generic IR plan does not provide the specific technical steps, communication templates, and rollback criteria for a region failover. |
| ASF-048 | The IR team can determine whether to failover or wait based on real-time data about the nature of the region impairment. | Derived | Premature failover causes unnecessary disruption; delayed failover exceeds RTO. The decision requires real-time telemetry. |
| ASF-049 | Forensic data from the failed region is preserved for post-incident analysis even after failover is complete. | Trust | Restoring service in the passive region should not destroy evidence in the failed region needed for root cause analysis. |
| ASF-050 | There is a documented decision tree for partial vs. full region failover (e.g., DB healthy but app unhealthy vs. full region loss). | Explicit | Not all impairments require a full failover; a nuanced response reduces unnecessary risk and data loss. |

---

### Pattern 15: Least Privilege

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-051 | The IAM role used by Route53 health checks has read-only permissions — cannot modify Route53 records. | Explicit | A health check Lambda with write access to Route53 could be used to redirect traffic maliciously. |
| ASF-052 | No single IAM user has permission to both modify Route53 failover records and disable CloudTrail logging. | Derived | Segregation of duties prevents a single compromised credential from both manipulating failover and covering tracks. |
| ASF-053 | The database replica user in the passive region has only read permissions until promotion. | Implicit | A replica user with write access to the passive database can cause split-brain data corruption if writes are accepted before promotion. |
| ASF-054 | S3 replication IAM role cannot delete objects from the source or destination bucket. | Derived | A replication role with delete capability propagates malicious or accidental deletions from primary to replica. |

---

### Pattern 16: Monitoring & Alerting

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-055 | Route53 health check status changes trigger alerts within 1 minute of detection. | Operational | Delayed alerting on health check transitions reduces the effective time window for human decision-making before RTO expires. |
| ASF-056 | Replication lag between primary and replica is monitored in near-real-time with threshold alerts. | Operational | Unmonitored replication lag leads to surprise RPO violations during failover. |
| ASF-057 | Monitoring infrastructure has cross-region failover capability — monitoring does not go dark with the active region. | Architectural | If monitoring is deployed only in the active region, the team loses visibility exactly when it is most critical. |
| ASF-058 | Alerts are configured for a passive region promotion (intentional or accidental) as a high-severity security event. | Derived | An unauthorized or accidental promotion of the passive region is a critical security event that must be detected immediately. |

---

### Pattern 17: Network Segmentation

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-059 | VPCs in each region are fully isolated from each other except through the documented replication channels. | Explicit | Direct VPC peering between regions creates unintended network paths that bypass security controls. |
| ASF-060 | Cross-region replication traffic does not traverse the public internet — it uses AWS private network backbone. | Implicit | Replication over the public internet exposes data to additional network-level threats; the AWS backbone provides inherent encryption. |
| ASF-061 | Security groups in the passive region allow traffic from the Route53 health checkers' IP ranges after failover. | Trust | Route53 health checkers use specific IP ranges; security groups that restrict inbound traffic may block health checks post-failover. |
| ASF-062 | There is no network path from the passive region's application subnet to the active region's database (or vice versa). | Architectural | Cross-region network paths that bypass the replication channel create split-brain write risks. |

---

### Pattern 18: Physical Security

*Not applicable — cloud-hosted.*

---

### Pattern 19: Supply Chain Security

*Deferred to Third-party Dependency.*

---

### Pattern 20: Third-party Dependency

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-063 | AWS Route53 service availability — a Route53 outage prevents both health checking and DNS failover. | Dependency | Route53 is the single point of failure for the entire failover mechanism; its control plane must remain available. |
| ASF-064 | AWS KMS availability in both regions — a KMS outage blocks encryption/decryption of replicated data. | Dependency | If KMS in the surviving region is also impaired, encrypted data becomes permanently inaccessible. |
| ASF-065 | The cloud provider's cross-region network SLA meets the RTO and RPO requirements. | Dependency | Cross-region network latency or bandwidth constraints are outside organizational control but directly impact RTO/RPO. |
| ASF-066 | There is an exit strategy for migrating to another cloud provider or region pair if the current provider cannot meet DR requirements. | Derived | Vendor lock-in to a single provider for DR means the organization has no recourse if the provider's multi-region guarantees fail. |
| ASF-067 | Third-party monitoring tools (PagerDuty, Datadog) are accessible from both regions and do not depend on the active region's egress. | Dependency | Monitoring SaaS tools that are only reachable from the active region's egress IP become unreachable after failover. |
| ASF-068 | DNS resolver services used by internal applications (Route53 Resolver, custom DNS) are available in both regions. | Dependency | Internal DNS failure in the promoted region causes all service-to-service communication to break. |

**Total (A): 68** (4 per pattern × 17 applicable patterns)

---

## Step 3: Comparison

### Overlap Mapping

| Human ID | ASF ID | Match Rationale |
|----------|--------|-----------------|
| H-001 | ASF-055 | Both address Route53 health check accuracy and alerting on status changes. |
| H-004 | ASF-056 | Both address replication lag monitoring relative to RPO. |
| H-005 | ASF-029 | Both assume S3 replication completes and no data is lost beyond RPO. |
| H-006 | ASF-006 | Both assume the passive region has capacity to handle production load after promotion. |
| H-008 | ASF-014 | Both require failover testing and runbooks. |
| H-009 | ASF-009 | Both assume the application in the passive region connects correctly. |
| H-010 | ASF-010 | Both require secrets and credentials to be available in the passive region. |
| H-011 | ASF-032 | Both require cross-region KMS key access for decryption after failover. |
| H-013 | ASF-053 | Both address the passive database replica's read-only state before promotion. |
| H-014 | ASF-025 | Both concern promotion scripts being correct and avoiding split-brain. |
| H-015 | ASF-057 | Both require monitoring to survive region failover. |
| H-017 | ASF-017 | Both require valid TLS certificates in the passive region. |
| H-019 | ASF-061 | Both require security group parity between regions. |
| H-020 | ASF-024 | Both require IAM role/permission parity between active and passive regions. |
| H-023 | ASF-035 | Both require cross-region replication encryption in transit. |
| H-024 | ASF-023 | Both require least-privilege IAM for S3 replication role. |
| H-028 | ASF-057 | Both require log aggregation from both regions to SIEM. |
| H-030 | ASF-050 | Both address depth and appropriateness of Route53 health checks. |
| H-032 | ASF-010 | Both assume no shared failure domains between regions. |
| H-033 | ASF-034 | Both require KMS key rotation to be consistent and not break access. |
| H-036 | ASF-057 | Both require cross-region log aggregation for monitoring. |
| H-037 | ASF-059 | Both assume no unintended network paths between regions. |
| H-038 | ASF-009 | Both assume application code is region-agnostic and behaves correctly on failover. |
| H-039 | ASF-056 | Both require replication lag measurement before promotion to meet RPO. |
| H-040 | ASF-012 | Both require the passive region to have working egress for external APIs. |
| H-041 | ASF-014 | Both address manual failover capability and human-in-the-loop decision-making. |

**Overlap (O): 26**

### Metrics

| Metric | Value | Formula |
|--------|-------|---------|
| Human assumptions (H) | 41 | Count of unique human-generated assumptions |
| ASF assumptions (A) | 68 | Count of unique ASF-generated assumptions |
| Overlap (O) | 26 | Count appearing in both lists |
| **Precision** | **38.2%** | O / A = 26/68 |
| **Recall** | **63.4%** | O / H = 26/41 |
| **F1 Score** | **47.7%** | 2 × (P × R) / (P + R) |
| Novel findings (A - O) | 42 | Assumptions ASF found that human missed (61.8% of ASF total) |
| Missed findings (H - O) | 15 | Assumptions human found that ASF missed (36.6% of human total) |

### Success Criteria Assessment

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Recall | >= 70% | 63.4% | ❌ Not met |
| Precision | >= 50% | 38.2% | ❌ Not met |
| Novel discoveries | >= 10% of total (A+O) | 38.2% (42/110) | ✅ Exceeded |
| Expert agreement (F1 proxy) | > 60% | 47.7% | ❌ Not met |

---

## Step 4: Analysis

### Ontology Overlap Analysis

| Ontology Category | Overlaps | ASF Total | Overlap Rate |
|-------------------|----------|-----------|-------------|
| Explicit | 7 | 12 | 58.3% |
| Derived | 8 | 18 | 44.4% |
| Operational | 5 | 16 | 31.3% |
| Implicit | 3 | 8 | 37.5% |
| Trust | 2 | 6 | 33.3% |
| Architectural | 0 | 6 | 0.0% |
| Dependency | 1 | 8 | 12.5% |
| Environmental | 0 | 4 | 0.0% |

**Best overlap:** Explicit and Derived categories showed the strongest agreement. Both the human and the ASF immediately recognize that cross-region KMS keys, replication encryption, and IAM scoping are critical assumptions.

**Worst overlap:** Architectural and Environmental categories had zero overlap. The ASF identified architectural concerns (region isolation, health check network paths, cross-region monitoring deployment) and environmental concerns (replication bandwidth, KMS regional failure) that the human did not list as assumptions — possibly because the human treats these as design constraints rather than hidden assumptions.

### What Humans Caught That ASF Missed (Missed Findings = 15)

1. **DR-specific testing and exercise details (H-008, H-025, H-035):** The human emphasized annual cutover tests, failback testing, and tabletop exercises. The ASF patterns lack a dedicated "DR Testing" sub-pattern.

2. **Application-level failover correctness (H-009, H-027, H-038):** Hard-coded database endpoints, region-specific branching logic, and application configuration that ties services to specific region endpoints are application-code-level concerns that the ASF does not surface.

3. **DNS TTL and propagation specifics (H-002, H-016, H-034):** DNS TTL values, client-side caching behavior, and Route53-specific timing details are operational concerns at a level of specificity below the ASF pattern resolution.

4. **Split-brain prevention (H-014, H-031):** The risk of both regions accepting writes simultaneously is a DR-specific architectural concern that the ASF patterns touch on (ASF-025) but not comprehensively.

5. **Data at rest in the passive region (H-033):** The human worried about KMS key rotation schedules diverging between regions — a detail of KMS operational management that the ASF treats more abstractly.

### What ASF Caught That Humans Missed (Novel Findings = 42)

1. **Change management for DR configuration (ASF-019 through ASF-022):** The human generated zero assumptions about the change process for Route53 or DR configurations. The ASF contributed a full pattern on change approval, drift detection, and cross-region configuration synchronization.

2. **Incident response specifics (ASF-047 through ASF-050):** The human assumed monitoring existed but did not consider the specific IR playbook for failover, the decision tree for partial vs. full failover, or forensic preservation in the failed region.

3. **Third-party dependencies beyond the cloud provider (ASF-063 through ASF-068):** The human focused on AWS service dependencies. The ASF surfaced monitoring SaaS dependency, internal DNS dependency, and vendor exit strategy — all risks outside the direct architecture diagram.

4. **Identity lifecycle cross-region (ASF-043 through ASF-046):** The human did not consider that IAM deprovisioning must be synchronized across regions, or that DR team membership must be current.

5. **Data residency and compliance (ASF-030):** The human did not consider that cross-region data replication might violate data residency laws — a critical compliance gap.

6. **Monitoring infrastructure survival (ASF-057):** The human assumed logs reach SIEM but did not consider that the monitoring infrastructure itself must survive a region failover.

### Architecture Complexity Assessment

Architecture #007 (Multi-Region DR) was rated **Moderate**. Key findings:
- The **human/ASF gap** (63.4% recall, 38.2% precision) is similar to the Simple architecture, suggesting that pattern coverage, not complexity, drives recall.
- The high **novelty rate (61.8%)** confirms that the ASF adds substantial value even for moderate-complexity architectures.
- The **missed findings** were concentrated in DR-specific testing practices, application-level correctness, and DNS operational details — areas that may benefit from a dedicated "Disaster Recovery Testing" pattern.

### Key Insight

The ASF pattern matrix has strong coverage of infrastructure, network, and operational concerns for DR, but lacks explicit patterns for:
- **DR testing and exercise governance** (tabletop, failback, annual certification)
- **Application-level region awareness** (hard-coded endpoints, region-specific branching)
- **DNS operational specifics** (TTL management, client caching, Route53-specific behaviors)

Adding a "Disaster Recovery Governance" pattern (or sub-pattern under Availability & Resilience) would likely close the recall gap to above 70%.

---

## Conclusions

| Metric | Target | Actual | Verdict |
|--------|--------|--------|---------|
| Recall | >= 70% | 63.4% | ❌ Below target — missing DR testing and application-level patterns |
| Precision | >= 50% | 38.2% | ❌ Below target — ASF is deliberately broad/generative |
| Novel discoveries | >= 10% | 38.2% | ✅ ASF adds substantial value beyond human reasoning |
| Expert agreement (F1) | > 60% | 47.7% | ❌ Below target — driven by low precision |

The ASF framework applied to Architecture #007 demonstrates strong exploration breadth (finding assumptions the human missed) but lower precision. For DR architecture risk identification, this trade-off is acceptable — false positives are preferable to assumptions missed that could lead to data loss during a real failover event. The primary actionable finding is the need for a **DR Testing & Governance** pattern and an **Application Region-Awareness** pattern to close the recall gap.
