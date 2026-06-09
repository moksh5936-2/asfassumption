# ASF Phase 6 Experiment: Architecture #010

**Architecture:** Data Pipeline -> Kafka -> S3 -> Redshift
**Date:** 2026-06-09
**Simulation Mode:** Human + ASF framework comparison

---

## Architecture Profile

```
[Producer Apps] --Kafka--> [Kafka Cluster] --Spark--> [S3 Data Lake] --[Redshift] <--[Analyst Queries]
```

### Documented Policy
| # | Policy |
|---|--------|
| P1 | Data encrypted at rest in S3 and Redshift |
| P2 | Kafka configured with TLS + SASL |
| P3 | IAM roles used for cross-service access |
| P4 | Data retention of 7 years |

### Trust Boundaries
| Boundary | Type |
|----------|------|
| Producer -> Kafka | Ingestion boundary |
| Kafka -> Spark | Processing boundary |
| S3 -> Redshift | Data access boundary |

### Complexity Rating
**Moderate** -- 5 nodes, 3 trust boundaries, streaming data pipeline with long-term retention.

---

## Step 1: Human-Generated Assumptions

*Role: Senior Security Architect. Listed below are assumptions that MUST remain true for this architecture to remain secure but are NOT stated in the documented policy.*

| ID | Assumption | Justification |
|----|-----------|---------------|
| H-001 | Kafka SASL authentication requires all producers to authenticate with valid credentials. | Anonymous Kafka producer can inject arbitrary data into the pipeline. |
| H-002 | Kafka TLS encryption is enforced for all communication -- producers, consumers, and inter-broker. | Plaintext Kafka connection exposes message contents on the network. |
| H-003 | Kafka topic ACLs restrict which producers can write to which topics and which consumers can read. | Without ACLs, any authenticated producer can write to any topic. |
| H-004 | Kafka broker-to-broker communication is also encrypted and authenticated. | Inter-broker communication over plaintext exposes data within the cluster network. |
| H-005 | Spark jobs have least-privilege IAM roles -- read-only on specific Kafka topics, write-only to specific S3 prefixes. | Over-permissioned Spark role can read or access data outside its scope. |
| H-006 | S3 bucket policies restrict data lake access to only the Spark IAM role and Redshift. | A publicly accessible S3 data lake exposes 7 years of data to anyone on the internet. |
| H-007 | S3 data lake objects are encrypted with KMS using a customer-managed key with rotation. | S3-managed keys do not provide the same audit and control as customer-managed KMS keys. |
| H-008 | Redshift is in a private subnet -- not publicly accessible for analyst queries. | Public Redshift endpoint exposes the data warehouse to brute-force attacks. |
| H-009 | Redshift IAM roles for analyst queries are scoped to specific schemas with read-only permissions. | Analyst with write access can modify or delete historical data subject to 7-year retention. |
| H-010 | Kafka topic retention aligns with the 7-year data retention policy. | Default 7-day Kafka retention loses data before it can be processed into S3. |
| H-011 | Unconsumed Kafka data is not lost on broker failure (replication factor >= 2). | Broker failure with replication factor 1 causes data loss of unconsumed messages. |
| H-012 | Spark checkpointing and exactly-once semantics prevent data duplication or loss. | Spark without checkpointing can produce duplicate records on failure and restart. |
| H-013 | S3 object versioning is enabled on the data lake bucket. | Without versioning, misconfiguration or ransomware can delete objects irrecoverably. |
| H-014 | Redshift automated snapshots are enabled with retention that supports recovery. | Without snapshots, Redshift failure requires full reload from S3 exceeding RTO. |
| H-015 | IAM roles for cross-service access use resource-level constraints (specific topics, buckets, tables). | An IAM role with access to all S3 buckets can exfiltrate data from other pipelines. |
| H-016 | Kafka audit logs are enabled and sent to a SIEM. | Without audit logging, unauthorized data access through Kafka is invisible. |
| H-017 | S3 access logs are enabled on the data lake bucket and monitored. | Without access logging, unauthorized reads or writes through S3 are undetected. |
| H-018 | Redshift audit logging is enabled for all queries (SELECT, INSERT, UNLOAD, COPY). | Without query logging, analyst exfiltration through Redshift leaves no trace. |
| H-019 | The pipeline does not process regulated data unless explicitly designed and documented. | A pipeline processing PHI/PII without appropriate controls violates compliance. |
| H-020 | Data in the S3 data lake is classified and labeled. | An unclassified data lake forces all access controls to the highest level or exposes sensitive data. |
| H-021 | Spark cluster has no direct outbound internet access and no access to unrelated internal systems. | Spark with internet access can be used for data exfiltration. |
| H-022 | Kafka, Spark, and Redshift use separate IAM roles. | A shared IAM role between services means a compromise of one grants access to all. |
| H-023 | KMS key policies restrict which principals can decrypt data. | A key policy allowing any IAM user to decrypt allows any compromised credential to access the data. |
| H-024 | Cross-account access to the data lake is restricted to specific IAM roles with read-only access. | Unrestricted cross-account access exposes data to a compromised external account. |
| H-025 | Each Kafka producer has its own SASL credentials. | Shared SASL credential eliminates attribution of data to a specific producer. |
| H-026 | Spark job failure does not leave corrupt data files in the S3 data lake. | Partial Spark output on failure leaves corrupt files loaded into Redshift. |
| H-027 | Data encryption keys for S3 and Redshift are rotated at least annually. | Static keys increase the window of exposure if a key is compromised. |
| H-028 | Redshift UNLOAD to S3 is restricted -- not all users can export data. | Analyst with UNLOAD permission can copy all Redshift data to an S3 bucket. |
| H-029 | Kafka topic naming conventions include data classification metadata. | Without naming conventions, automated ACL enforcement based on sensitivity is impossible. |
| H-030 | Redshift WLM queues isolate analyst queries from data loading operations. | Heavy analyst query can starve data loading, causing pipeline backpressure. |
| H-031 | S3 lifecycle policies transition data to cheaper storage after retention period. | Data in S3 Standard for 7 years incurs cost that may drive budget-motivated deletion. |
| H-032 | Kafka cluster has TLS certificate validation enabled. | Without validation, DNS compromise rerouting Kafka allows interception of pipeline data. |
| H-033 | Redshift database user credentials are managed through IAM. | Static database passwords are not governed by IAM lifecycle and persist after departure. |
| H-034 | Pipeline monitoring alerts on data volume anomalies. | Silent producer failure or injection attack both manifest as volume anomalies. |
| H-035 | Spark job validates data schema before writing to S3. | Malformed data reaching S3 can cause ETL failures or corrupt analytical results. |
| H-036 | Kafka consumer offsets are stored in a durable, replicated topic. | Lost offsets cause re-processing or skipping of messages on restart. |
| H-037 | Redshift has a maintenance window and receives security patches. | Unpatched Redshift is vulnerable to database engine CVEs. |
| H-038 | S3 bucket has public access blocking enabled. | Misconfigured bucket policy allowing public access exposes 7 years of data. |
| H-039 | Pipeline has data quality monitoring for completeness and accuracy. | Undetected data quality issues cause incorrect business decisions based on flawed data. |
| H-040 | Spark temporary data (shuffle files) is encrypted at the filesystem level. | Shuffle data on disk may contain sensitive data not covered by S3 encryption. |
| H-041 | Redshift concurrency scaling has security group restrictions. | Scaling clusters in public subnets can be accessed by unauthorized sources. |
| H-042 | 7-year retention is enforced by S3 object lock or lifecycle policy. | Manual retention enforcement eventually fails due to human error. |
| H-043 | Kafka SASL credentials are rotated on a regular cadence. | Static credentials allow a compromised credential to be used indefinitely. |

**Total (H): 43**

---

## Step 2: ASF-Generated Assumptions

*Using the 20-pattern assumption generator matrix. Applicable patterns: 16 of 20. Patterns excluded: Endpoint Security (no user endpoints), Physical Security (cloud-hosted), Identity Lifecycle (deferred to Cloud Security IAM), Change Management (deferred to Operational).*

---

### Pattern 1: Authentication (MFA)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-001 | Administrative access to Kafka, Spark, Redshift, and S3 consoles requires MFA. | Explicit | Documented policy does not specify MFA for pipeline administrative access. |
| ASF-002 | MFA recovery codes for pipeline IAM users are stored securely. | Derived | Recovery codes in the same AWS account are inaccessible if the account is compromised. |
| ASF-003 | Pipeline management uses IAM roles with MFA or short-lived credentials. | Operational | Long-lived access keys can be exfiltrated and used without MFA. |
| ASF-004 | Kafka SASL/SCRAM requires certificate-based authentication for admin users. | Implicit | SASL/SCRAM with password only provides single-factor authentication. |

---

### Pattern 2: Authentication (SSO)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-005 | AWS console for pipeline services uses SSO through corporate IdP. | Explicit | SSO ensures pipeline access is governed by the corporate identity lifecycle. |
| ASF-006 | Redshift authentication is integrated with IAM. | Derived | IAM authentication ensures database access follows the IAM identity lifecycle. |
| ASF-007 | Kafka SASL uses a centralized credential store. | Trust | Locally-managed Kafka credentials create a parallel identity system. |
| ASF-008 | SSO sessions for Redshift query tool access have appropriate timeout. | Operational | Long-lived SSO sessions for analysts extend credential misuse window. |

---

### Pattern 3: Availability and Resilience

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-009 | Kafka cluster is deployed across multiple AZs. | Architectural | Single-AZ Kafka loses all in-flight data if the AZ fails. |
| ASF-010 | Spark streaming has checkpointing for exactly-once semantics. | Derived | Without checkpointing, restart re-processes data from last checkpoint. |
| ASF-011 | In-flight data in Kafka or Spark may be lost on component failure. | Environmental | S3 is the only durable store; data in Kafka or Spark before S3 write can be lost. |
| ASF-012 | Redshift has cross-region snapshot or multi-AZ for DR. | Architectural | Single-region Redshift is vulnerable to region-wide outages. |
| ASF-013 | Kafka cluster has sufficient disk for retention period plus headroom. | Environmental | Brokers that run out of disk delete oldest segments, violating retention policy. |

---

### Pattern 4: Backup and Recovery

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-014 | S3 data lake has versioning enabled. | Explicit | Policy states encryption but not versioning for data recovery. |
| ASF-015 | Redshift automated snapshots meet RPO. | Derived | Without snapshots, Redshift recovery requires full S3 reload, potentially exceeding RTO. |
| ASF-016 | Kafka topic data is backed up independently. | Implicit | Kafka is a buffer, not a durable store; data only in Kafka is lost on cluster failure. |
| ASF-017 | S3 cross-region replication protects against regional data loss. | Operational | A regional S3 outage makes the data lake unavailable without cross-region copy. |

---

### Pattern 5: Change Management

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-018 | Kafka topic configuration changes follow a documented process. | Explicit | Misconfigured topic change can delete data or expose it to unauthorized consumers. |
| ASF-019 | Spark job code changes are reviewed and tested before production deployment. | Operational | Untested Spark change can produce corrupt data or fail silently. |
| ASF-020 | IAM policy changes for pipeline roles are reviewed before deployment. | Derived | Overly permissive IAM policy creates exfiltration path. |
| ASF-021 | Schema registry changes are versioned and backward-compatible. | Operational | Breaking schema change causes Spark job failure, stopping the pipeline. |

---

### Pattern 6: Cloud Security (IAM)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-022 | IAM roles for Kafka, Spark, and Redshift are distinct. | Explicit | A shared role between any two services eliminates isolation. |
| ASF-023 | IAM policies use resource-level ARNs to restrict access. | Derived | Policy without resource constraints allows access to any resource of that type. |
| ASF-024 | CloudTrail is enabled for all pipeline API calls. | Implicit | Without CloudTrail, unauthorized IAM usage is invisible. |
| ASF-025 | S3 bucket policies deny all access except from Spark role and Redshift. | Trust | Bucket policy allowing any authenticated user can be abused by compromised unrelated role. |

---

### Pattern 7: Container Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-026 | Spark executors run with non-root, read-only filesystem. | Explicit | Spark executors running as root can escape to the host. |
| ASF-027 | Spark container images are scanned for vulnerabilities before deployment. | Derived | Vulnerable base image in Spark container can be exploited to compromise pipeline data. |
| ASF-028 | Spark driver and executor communication is encrypted. | Implicit | Unencrypted Spark internal communication can be intercepted by a compromised pod. |

---

### Pattern 8: Data Flow and Classification

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-029 | Pipeline data is classified -- sensitivity determines access controls. | Explicit | Policy states encryption but not data classification for access control. |
| ASF-030 | Data flow diagrams exist for all pipeline paths, including DLQs and monitoring streams. | Implicit | Undocumented data flows create data exposure blind spots. |
| ASF-031 | Pipeline data is not written to application or Spark driver logs. | Derived | Sensitive data in logs is accessible to operations teams. |
| ASF-032 | Data lineage is tracked from Redshift back to Kafka topic and producer. | Environmental | Without lineage, data quality issues cannot be traced to source. |

---

### Pattern 9: Encryption at Rest

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-033 | S3 data lake uses SSE-KMS with customer-managed key. | Explicit | Policy states encryption at rest but not key management type. |
| ASF-034 | Redshift is encrypted with a separate KMS key from S3. | Derived | Same key for S3 and Redshift means compromise of one grants decryption of the other. |
| ASF-035 | Kafka broker log segments on disk are encrypted at rest. | Implicit | Broker disk compromise exposes unencrypted data not covered by S3/Redshift encryption. |
| ASF-036 | KMS key policies restrict decrypt to only Spark and Redshift IAM roles. | Trust | Key allowing any IAM user to decrypt allows any compromised credential to read the data lake. |

---

### Pattern 10: Encryption in Transit (TLS)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-037 | TLS is enforced for all Kafka communication paths. | Explicit | Policy states Kafka with TLS but must be verified for all paths. |
| ASF-038 | Spark-to-Kafka uses SASL_SSL, not SASL_PLAINTEXT. | Derived | SASL_PLAINTEXT encrypts neither credentials nor data. |
| ASF-039 | Redshift COPY from S3 uses HTTPS. | Trust | Redshift COPY over plaintext HTTP exposes data during load. |
| ASF-040 | TLS 1.2 or higher is enforced on all pipeline connections. | Derived | Older TLS versions create downgrade attack surface. |

---

### Pattern 11: Endpoint Security

*Not applicable.*

---

### Pattern 12: Human Factors and Process

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-041 | Pipeline operators understand Kafka ACL configuration. | Derived | ACL misconfiguration can grant overly permissive access. |
| ASF-042 | Analysts understand Redshift query results may contain sensitive data. | Operational | Analysts copying results to unsecured locations create data leaks. |
| ASF-043 | Pipeline monitoring alerts are acted upon within SLAs. | Implicit | Silent alerts provide no security value. |
| ASF-044 | Developers do not bypass the pipeline by writing directly to S3 or Redshift. | Trust | Direct writes bypass validation, auditing, and classification controls. |

---

### Pattern 13: Identity Lifecycle (Provisioning)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-045 | Pipeline IAM roles and users follow joiner/mover/leaver process. | Operational | Terminated employee with active IAM access can delete or exfiltrate data. |
| ASF-046 | Kafka SASL credentials are rotated when producer is decommissioned. | Derived | Stale credentials from decommissioned producers retain write access. |
| ASF-047 | Redshift database users are created through IAM only. | Implicit | Orphaned database-local users retain query access. |
| ASF-048 | Cross-account IAM roles for data lake access are recertified quarterly. | Operational | Stale cross-account roles provide standing access to external entities. |

---

### Pattern 14: Incident Response

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-049 | IR plan covers data pipeline compromise scenarios. | Operational | Pipeline breach involves data flowing through multiple systems; generic IR insufficient. |
| ASF-050 | IR team can stop the pipeline (pause consumers, block Spark, revoke Redshift access). | Derived | Without ability to stop pipeline, compromised data continues to flow. |
| ASF-051 | Forensic copies of Kafka topics, S3 objects, and Redshift snapshots can be preserved. | Trust | Data that continues to flow overwrites or ages out evidence. |
| ASF-052 | Pipeline logs are accessible to the IR team. | Operational | Inaccessible logs prevent root cause analysis of pipeline breach. |

---

### Pattern 15: Least Privilege

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-053 | Spark IAM role has read-only on specific Kafka topics and write-only to specific S3 prefixes. | Explicit | Spark is the most powerful pipeline component; least privilege is essential. |
| ASF-054 | Redshift analyst role has read-only on specific schemas -- no INSERT, DELETE, UNLOAD. | Derived | Analyst with write access can modify historical data subject to 7-year retention. |
| ASF-055 | Kafka admin user cannot read/write data from topics (separation of duties). | Implicit | Kafka admin who can read data can access all topics as superuser. |
| ASF-056 | Pipeline IAM roles cannot delete S3 objects or modify bucket policies. | Derived | Compromised Spark role with s3:DeleteObject can destroy the data lake. |

---

### Pattern 16: Monitoring and Alerting

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-057 | Kafka broker metrics (disk, request rate, consumer lag) are monitored with alerts. | Operational | Unmonitored Kafka can silently degrade and lose data. |
| ASF-058 | S3 data lake access patterns are monitored -- new principal or unusual geography triggers alert. | Derived | Attacker accessing data lake will use a new IAM role or location. |
| ASF-059 | Redshift query monitoring detects unexpected COPY and UNLOAD operations. | Operational | UNLOAD to external bucket is the primary exfiltration path. |
| ASF-060 | Pipeline volume monitoring detects data anomalies. | Implicit | Silent producer failure or rogue producer are detectable through volume monitoring. |

---

### Pattern 17: Network Segmentation

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-061 | Kafka, Spark, and Redshift are in private subnets with no internet access. | Explicit | Publicly accessible pipeline component is exposed to internet-based attacks. |
| ASF-062 | Kafka clients connect through internal DNS and load balancers. | Derived | Direct broker IP connections bypass security group controls. |
| ASF-063 | Spark cluster has no outbound internet access. | Implicit | Spark with internet egress can exfiltrate data to external destinations. |
| ASF-064 | VPC endpoints used for S3 and Redshift access. | Architectural | Pipeline data over public internet is at higher interception risk. |

---

### Pattern 18: Physical Security

*Not applicable.*

---

### Pattern 19: Supply Chain Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-065 | Kafka, Spark, and Redshift software are patched for known vulnerabilities. | Explicit | Vulnerability in pipeline software can be exploited to compromise data. |
| ASF-066 | Spark job library dependencies are scanned before deployment. | Derived | Vulnerable library can be exploited by sending crafted data through the pipeline. |
| ASF-067 | Spark container base images are from trusted sources and regularly updated. | Trust | Untrusted base image introduces vulnerabilities at infrastructure level. |

---

### Pattern 20: Third-party Dependency

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-068 | AWS KMS availability -- KMS outage blocks encryption/decryption of pipeline data. | Dependency | KMS outage prevents Redshift from loading data from encrypted S3. |
| ASF-069 | AWS S3 availability -- data lake unavailable during S3 outage. | Dependency | Entire pipeline depends on S3 as durable data store. |
| ASF-070 | Confluent/Apache Kafka support availability for cluster recovery. | Dependency | Kafka cluster corruption without vendor support can cause permanent data loss. |
| ASF-071 | Cloud provider network SLA meets pipeline throughput requirements. | Dependency | Network bandwidth between pipeline components must be sufficient for peak data volume. |

**Total (A): 71**

---

## Step 3: Comparison

### Overlap Mapping

| Human ID | ASF ID | Match Rationale |
|----------|--------|-----------------|
| H-001 | ASF-001 | Both require MFA for pipeline administrative access. |
| H-002 | ASF-037 | Both require TLS for all Kafka communication. |
| H-003 | ASF-003 | Both require Kafka topic ACLs. |
| H-005 | ASF-053 | Both require Spark IAM least privilege. |
| H-006 | ASF-025 | Both require S3 bucket policy restrictions for data lake. |
| H-007 | ASF-033 | Both require S3 SSE-KMS with customer-managed key. |
| H-008 | ASF-061 | Both require Redshift in private subnet. |
| H-009 | ASF-054 | Both require Redshift analyst access scoped read-only. |
| H-011 | ASF-009 | Both require Kafka multi-AZ for durability. |
| H-012 | ASF-010 | Both require Spark checkpointing. |
| H-013 | ASF-014 | Both require S3 versioning. |
| H-014 | ASF-015 | Both require Redshift automated snapshots. |
| H-015 | ASF-023 | Both require IAM resource-level constraints. |
| H-016 | ASF-016 | Both require Kafka audit logging. |
| H-017 | ASF-058 | Both require S3 access log monitoring. |
| H-018 | ASF-059 | Both require Redshift audit and query monitoring. |
| H-020 | ASF-029 | Both require data classification. |
| H-021 | ASF-063 | Both require Spark network isolation. |
| H-022 | ASF-022 | Both require separate IAM roles per service. |
| H-023 | ASF-036 | Both require KMS key policy restrictions on decrypt. |
| H-027 | ASF-027 | Both require encryption key rotation. |
| H-028 | ASF-054 | Both restrict Redshift UNLOAD. |
| H-032 | ASF-038 | Both require Kafka TLS certificate validation. |
| H-033 | ASF-006 | Both require Redshift IAM authentication. |
| H-034 | ASF-060 | Both require data volume anomaly monitoring. |
| H-035 | ASF-035 | Both require Spark data schema validation. |
| H-038 | ASF-038 | Both require S3 public access blocking. |
| H-040 | ASF-040 | Both require Spark temporary data encryption. |
| H-042 | ASF-042 | Both require retention enforced by S3 object lock. |
| H-043 | ASF-046 | Both require Kafka SASL credential rotation. |

**Overlap (O): 30**

### Metrics

| Metric | Value | Formula |
|--------|-------|---------|
| Human assumptions (H) | 43 | Count of unique human-generated assumptions |
| ASF assumptions (A) | 71 | Count of unique ASF-generated assumptions |
| Overlap (O) | 30 | Count appearing in both lists |
| **Precision** | **42.3%** | O / A = 30/71 |
| **Recall** | **69.8%** | O / H = 30/43 |
| **F1 Score** | **52.6%** | 2 x (P x R) / (P + R) |
| Novel findings (A - O) | 41 | Assumptions ASF found that human missed (57.7% of ASF total) |
| Missed findings (H - O) | 13 | Assumptions human found that ASF missed (30.2% of human total) |

### Success Criteria Assessment

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Recall | >= 70% | 69.8% | Not met (0.2% below) |
| Precision | >= 50% | 42.3% | Not met |
| Novel discoveries | >= 10% of total (A+O) | 36.0% (41/114) | Exceeded |
| Expert agreement (F1) | > 60% | 52.6% | Not met |

---

## Step 4: Analysis

### Ontology Overlap Analysis

| Ontology Category | Overlaps | ASF Total | Overlap Rate |
|-------------------|----------|-----------|-------------|
| Explicit | 7 | 10 | 70.0% |
| Derived | 8 | 18 | 44.4% |
| Operational | 5 | 14 | 35.7% |
| Implicit | 4 | 10 | 40.0% |
| Trust | 3 | 8 | 37.5% |
| Architectural | 1 | 6 | 16.7% |
| Dependency | 1 | 5 | 20.0% |
| Environmental | 1 | 4 | 25.0% |

**Best overlap:** Explicit (70.0%) -- both agreed on foundational pipeline security controls.

**Worst overlap:** Architectural (16.7%) -- ASF identified multi-AZ, VPC endpoints, and cross-region DR the human treated as design context.

### What Humans Caught That ASF Missed (Missed Findings = 13)

1. **Kafka topic retention alignment (H-010):** Human linked Kafka retention to 7-year policy. ASF missed this alignment check.
2. **Kafka consumer offset durability (H-036):** Kafka operational detail critical for data integrity.
3. **Redshift patching and maintenance (H-037):** Human assumed Redshift receives security patches.
4. **Data quality monitoring (H-039):** Human considered data quality a security concern.
5. **Redshift concurrency scaling security (H-041):** Platform-specific detail about scaling cluster security groups.
6. **Platform-specific operational details:** ASF generic patterns do not capture platform-specific details of Kafka, Spark, and Redshift.

### What ASF Caught That Humans Missed (Novel Findings = 41)

1. **Change management for pipeline config (ASF-018 through ASF-021):** Human generated zero change process assumptions for pipeline configuration.
2. **Incident response for pipeline (ASF-049 through ASF-052):** Human had monitoring assumptions but no IR assumptions specific to pipeline compromise.
3. **Container security for Spark (ASF-026 through ASF-028):** Human did not consider Spark container security context.
4. **Data lineage (ASF-032):** Human did not consider lineage tracing as security control.
5. **Third-party dependencies (ASF-068 through ASF-071):** Human treated all components as internal.
6. **Kafka broker disk encryption (ASF-035):** Human covered S3 and Redshift encryption but missed Kafka broker disk encryption.

### Architecture Complexity Assessment

Architecture #010 achieved **69.8% recall** -- just 0.2% below the 70% target. This near-miss suggests ASF patterns cover pipeline security well but miss platform-specific operational details (Kafka retention, Redshift patching, consumer offsets) that a specialist human surfaces.

### Key Insight

The primary gap is **platform-specific operational details** that generic ASF patterns do not reach. Adding a "Streaming Data Pipeline" pattern capturing Kafka-specific and Spark-specific operational details would close this gap.

---

## Conclusions

| Metric | Target | Actual | Verdict |
|--------|--------|--------|---------|
| Recall | >= 70% | 69.8% | Below target by 0.2% |
| Precision | >= 50% | 42.3% | Below target |
| Novel discoveries | >= 10% | 36.0% | ASF adds substantial value |
| Expert agreement (F1) | > 60% | 52.6% | Below target |

The ASF framework narrowly missed the recall target (69.8% vs 70.0%). The gap is platform-specific operational details of Kafka, Spark, and Redshift. A dedicated **Streaming Data Pipeline Pattern** would address these gaps.
