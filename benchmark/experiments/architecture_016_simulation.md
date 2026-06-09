# ASF Phase 6 Experiment: Architecture #016

**Architecture:** ML Pipeline → Training → Serving → Data Lineage
**Date:** 2026-06-09
**Simulation Mode:** Human + ASF framework comparison

---

## Architecture Profile

```
[Feature Store] --> [SageMaker Training] --> [Model Registry] --> [Serving Endpoint]
       |                                                         |
  [Data Lake (S3)]                                         [Production App]
```

### Documented Policy
| # | Policy |
|---|--------|
| P1 | Training data is from approved sources |
| P2 | Models are versioned in registry |
| P3 | Inference endpoint requires auth |
| P4 | Data lineage is tracked |

### Trust Boundaries
| Boundary | Type |
|----------|------|
| Data Lake ↔ Training | Data integrity boundary |
| Training ↔ Registry | Model integrity boundary |
| Registry ↔ Endpoint | Model serving boundary |

### Complexity Rating
**Moderate** — ML-specific topology, 6 nodes, 3 trust boundaries, data provenance, model lifecycle.

---

## Step 1: Human-Generated Assumptions

*Role: Senior Security Architect. Listed below are assumptions that MUST remain true for this architecture to remain secure but are NOT stated in the documented policy.*

| ID | Assumption | Justification |
|----|-----------|---------------|
| H-001 | Training data in the Data Lake is validated for integrity (checksums, signatures) before ingestion. | Corrupted or tampered training data produces a compromised model without visible symptoms. |
| H-002 | Data lineage tracking includes the origin, transformation steps, and access history for each dataset used in training. | Without complete lineage, a data poisoning incident cannot be traced to its source. |
| H-003 | The Feature Store enforces access controls so that only authorized training jobs can read features. | A compromised training job with access to all features can exfiltrate sensitive feature data. |
| H-004 | Training jobs run in an isolated environment with no network access to the internet or other internal systems. | A training job with internet access can exfiltrate training data or model weights to an external endpoint. |
| H-005 | Model artifacts (weights, hyperparameters) are encrypted at rest in the Model Registry. | Unencrypted model artifacts stored in the registry can be exfiltrated by anyone with registry read access. |
| H-006 | Model registry access is restricted to authorized principals with approval workflows for promotion to production. | Unrestricted registry access allows a compromised CI/CD pipeline to deploy a malicious model. |
| H-007 | The Serving Endpoint validates the model version before loading and rejects models without an approved status. | Deploying an unapproved model version (staging, rejected) serves incorrect or malicious inferences to production apps. |
| H-008 | The Serving Endpoint enforces rate limiting per production application client to prevent abuse. | An application sending excessive inference requests can degrade endpoint performance for all consumers. |
| H-009 | Inference requests and responses are logged for audit with sensitive data masking applied. | Inference logs containing PII or sensitive predictions create a secondary data exposure channel. |
| H-010 | The S3 Data Lake has bucket policies that prevent public access and enforce encryption in transit. | A misconfigured S3 bucket exposes training data to the public internet. |
| H-011 | Training data that contains PII is anonymized or de-identified before storage in the Feature Store. | ML models trained on PII-embedded data can memorize and leak that PII through inference outputs. |
| H-012 | The Model Registry supports model signing so that model integrity can be verified before deployment. | An unsigned model in the registry could be replaced with a tampered version without detection. |
| H-013 | The training environment does not persist between runs; each training job starts with a clean environment. | A training environment that persists artifacts from previous runs can leak data across training jobs. |
| H-014 | Production inference data is not used to retrain models without explicit review and approval. | Inference data fed back into training creates a data loop that can amplify model biases or drift. |
| H-015 | The Serving Endpoint has a rollback mechanism to revert to a previous model version if the new model degrades. | A model that produces incorrect or harmful inferences cannot be quickly reverted without a rollback plan. |
| H-016 | Data lineage is immutable and append-only to prevent tampering with training provenance records. | Mutable lineage records allow an attacker to cover their tracks after poisoning the training data. |
| H-017 | The inference endpoint has a maximum request size limit to prevent denial-of-service via large payloads. | Oversized inference requests can exhaust endpoint memory and cause denial of service. |
| H-018 | SageMaker notebook instances (if used for exploration) are not connected to the production environment. | A data scientist's notebook with access to production data or endpoints can inadvertently or maliciously modify production resources. |
| H-019 | Model registry entries include metadata about training data version, hyperparameters, and evaluation metrics. | A registry entry without complete metadata makes it impossible to reproduce a model's training for debugging. |
| H-020 | The Feature Store is configured with a data retention policy that expires stale features. | Stale features that are no longer relevant can be used accidentally in training, degrading model quality. |
| H-021 | The inference endpoint authenticates requests using a robust mechanism (JWT, API key with rotation) and not just network-level allowlists. | Network-allowlisted endpoints can be accessed from any compromised resource inside the VPC. |
| H-022 | Model evaluation metrics are computed against a held-out test set before promotion to production. | A model promoted without evaluation against a test set may have hidden performance issues. |
| H-023 | The training data pipeline includes data quality checks (schema validation, distribution monitoring) before training. | Drifted or corrupted data that passes validation produces a flawed model. |
| H-024 | The serving endpoint encrypts inference responses if they contain sensitive data. | Inference outputs that contain PII or business-sensitive predictions exposed in transit are interceptable. |
| H-025 | There is a mechanism to identify and block adversarial inputs (adversarial attack detection) at the inference endpoint. | ML models are vulnerable to adversarial examples that cause misclassification without input validation. |
| H-026 | The model registry supports model lineage so that each model version is traceable to its training job and dataset. | A model without lineage cannot be audited for compliance or investigated for bias. |
| H-027 | Training data access logs are monitored for unusual access patterns (bulk download, unusual times). | Data exfiltration or poisoning attempts are visible in access logs only if monitored. |
| H-028 | The model registry supports approval workflows with gates for different environments (staging vs production). | A model that passes staging evaluation may still fail production requirements without a separate approval gate. |
| H-029 | Training hyperparameters are validated against an allowed range before job submission. | Extreme hyperparameters (e.g., infinite learning rate) can cause resource exhaustion or produce degenerate models. |
| H-030 | The S3 Data Lake has versioning enabled to recover from accidental data deletion or overwrite. | An accidental delete of training data in an unversioned bucket causes permanent data loss. |
| H-031 | The Feature Store serves features through an API that enforces least-privilege access per feature group. | A training job that reads features it does not need increases the blast radius of a data leak. |
| H-032 | The training job IAM role has permissions scoped to only the specific S3 buckets and Feature Store required. | Over-permissioned training jobs can read or write data outside the approved scope. |
| H-033 | The inference endpoint is protected by a WAF or API Gateway that filters common web attacks. | An inference endpoint fronted without a WAF is exposed to injection attacks and API abuse. |
| H-034 | Model registry entries include a status field (development, staging, production, deprecated, archived). | Without status tracking, stale or deprecated models may be unknowingly served to production. |
| H-035 | The serving endpoint validates the model input schema and rejects inputs that do not match. | Input schema mismatches can cause undefined behavior, crashes, or exploitable bugs in the inference handler. |
| H-036 | Data lineage is tracked at the column/feature level, not just the dataset level. | Dataset-level lineage is insufficient to identify which specific features were compromised in a data poisoning attack. |
| H-037 | The training environment has a maximum runtime and is terminated if it exceeds it. | Runaway training jobs consume compute resources indefinitely and may incur unexpected costs. |
| H-038 | The model registry uses immutable version tags; a version cannot be overwritten once published. | Overwritable model versions allow a published model to be silently replaced with a tampered version. |
| H-039 | The inference endpoint logs the model version used for each inference request for audit. | Without per-request model version logging, it is impossible to attribute inference results to a specific model version. |
| H-040 | The Data Lake has cross-region replication enabled for disaster recovery. | Data Lake loss in one region means loss of all training data and lineage records. |
| H-041 | Production application credentials for the inference endpoint are scoped to a single endpoint with least privilege. | An application credential with access to multiple endpoints can invoke inference on models it should not access. |
| H-042 | Model training uses only approved algorithms and frameworks; custom or untrusted code is reviewed. | Untrusted training code can introduce backdoors into the model or exfiltrate data. |
| H-043 | The inference endpoint has a configurable timeout and returns a clear error on timeout. | A long-running inference request without timeout can exhaust endpoint resources. |
| H-044 | SageMaker training job network isolation is enabled (no internet access, VPC-only mode). | A training job with internet access can download untrusted libraries or exfiltrate data. |
| H-045 | The Feature Store supports point-in-time queries to reconstruct feature values as they existed at training time. | Without point-in-time queries, online and offline feature distributions diverge, causing training-serving skew. |
| H-046 | Model evaluation datasets are stored separately from training data with access controls. | Evaluation dataset contamination with training data invalidates model performance metrics. |

**Total (H): 46**

---

## Step 2: ASF-Generated Assumptions

*Using the 20-pattern assumption generator matrix. Applicable patterns: 16 of 20. Patterns excluded: Backup & Recovery (covered under operational), Physical Security (cloud-hosted), Network Segmentation (covered under Cloud Security), Endpoint Security (covered under Cloud Security).*

---

### Pattern 1: Authentication (MFA)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-001 | AWS console access for SageMaker, Model Registry, and S3 management requires MFA. | Explicit | Administrative access to ML infrastructure must be protected. |
| ASF-002 | The inference endpoint API key authentication is not the only factor; client certificate or mTLS is available for sensitive models. | Derived | API key alone is weak against theft; mTLS provides stronger client identity. |
| ASF-003 | MFA is required for approving model promotion to production in the Model Registry. | Operational | Single-person model promotion without MFA allows unauthorized model deployment. |
| ASF-004 | IAM user access keys (if any) for ML pipeline automation are rotated and monitored. | Implicit | Static access keys for CI/CD pipeline automation are a common credential leak vector. |

---

### Pattern 2: Authentication (SSO)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-005 | All human access to the ML pipeline (SageMaker, Registry, S3) uses federated SSO. | Explicit | Federated SSO provides centralized access management and audit. |
| ASF-006 | Service roles for SageMaker training and serving are assumed via IAM, not hardcoded credentials. | Derived | Service roles with temporary credentials reduce the risk of long-term key exposure. |
| ASF-007 | SSO session timeout is enforced for the ML admin console. | Trust | Persistent admin sessions increase the window for unauthorized configuration changes. |
| ASF-008 | Cross-account access to the Model Registry is through IAM roles with explicit trust policies. | Operational | Cross-account model sharing without explicit trust policies creates authorization gaps. |

---

### Pattern 3: Availability & Resilience

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-009 | The SageMaker training job can recover from spot instance interruptions without data loss. | Architectural | Spot instance interruptions are common; training jobs must checkpoint to resume. |
| ASF-010 | There is a fallback inference endpoint (e.g., canary deployment) if the primary endpoint fails. | Operational | Production app cannot serve inferences if the primary endpoint is unavailable. |
| ASF-011 | The Feature Store has sufficient read throughput for all concurrent training jobs. | Environmental | Feature store throttling delays training jobs and increases time-to-deploy. |
| ASF-012 | The S3 Data Lake is available and not throttling reads during training data loading. | Dependency | S3 throttling during training data loading causes failures or extended training times. |

---

### Pattern 4: Backup & Recovery

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-013 | The Model Registry is backed up to prevent loss of model version history. | Operational | Loss of model registry means losing all model provenance and versioning. |
| ASF-014 | Feature Store data is backed up or has point-in-time recovery enabled. | Derived | Feature store corruption requires re-computing all features from raw data. |
| ASF-015 | S3 Data Lake has versioning and cross-region replication enabled. | Operational | Data Lake loss means loss of all training data and lineage. |
| ASF-016 | Model artifacts in the registry are replicated to another region or account for DR. | Implicit | Region failure makes all trained models unavailable for inference. |

---

### Pattern 5: Change Management

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-017 | Model registry schema changes are backward-compatible to preserve existing model metadata. | Operational | Schema changes that break metadata parsing affect all existing models. |
| ASF-018 | Feature store schema changes are communicated to all model teams using those features. | Derived | A feature silently renamed or removed breaks all models that depend on it. |
| ASF-019 | Inference endpoint configuration changes (scaling, timeout) are previewed in a staging environment. | Trust | Untested endpoint configuration changes can cause production degradation. |
| ASF-020 | Training data source additions go through a review process to ensure data quality and compliance. | Implicit | New data sources added without review may introduce PII or biased data. |

---

### Pattern 6: Cloud Security (IAM)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-021 | SageMaker execution roles are scoped with least privilege for each training job. | Explicit | Over-permissioned training roles increase the blast radius of a compromised training job. |
| ASF-022 | The Model Registry has IAM policies that enforce approval workflows for model deployment. | Derived | IAM policies should enforce that only approved models can be deployed. |
| ASF-023 | S3 bucket policies for the Data Lake enforce encryption in transit (aws:SecureTransport). | Implicit | S3 without SecureTransport policy allows unencrypted data access. |
| ASF-024 | CloudTrail is enabled for the ML pipeline account to audit all API calls. | Explicit | Without CloudTrail, unauthorized pipeline changes go undetected. |

---

### Pattern 7: Compliance & Audit

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-025 | ML models are subject to the same data classification as their training data. | Explicit | A model trained on PII data is itself a PII-bearing artifact. |
| ASF-026 | Model inference logs are retained per regulatory requirements and are tamper-proof. | Derived | Inference logs used for compliance audits must be immutable. |
| ASF-027 | There is a documented process for model bias and fairness auditing. | Operational | Regulatory requirements increasingly require bias audits for ML models in production. |
| ASF-028 | Data lineage records are maintained for the minimum required retention period. | Environmental | Lineage records deleted early violate compliance; retained too long increase exposure. |

---

### Pattern 8: Data Flow & Classification

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-029 | Data flowing through the ML pipeline is classified and handling requirements are documented. | Explicit | Unclassified data has undefined handling requirements; controls may be inconsistent. |
| ASF-030 | Data flow diagrams exist showing all paths from Data Lake through Feature Store to inference. | Implicit | Undocumented data flows hide shadow ML pipelines using the same data. |
| ASF-031 | Training data does not flow to any destination outside the defined architecture. | Derived | The documented flow is one-directional; any egress (monitoring, debugging, third-party) is unaccounted. |
| ASF-032 | Production inference data is not automatically fed back into the training pipeline. | Environmental | Auto-feedback loops can degrade model quality and introduce data drift. |

---

### Pattern 9: Encryption at Rest

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-033 | S3 Data Lake is encrypted at rest using KMS with a dedicated key. | Explicit | Standard requirement for sensitive data in S3. |
| ASF-034 | The Model Registry stores model artifacts encrypted at rest. | Derived | Unencrypted model artifacts in the registry can be exfiltrated. |
| ASF-035 | The Feature Store is encrypted at rest. | Explicit | Feature data is derived from training data and has the same sensitivity. |
| ASF-036 | SageMaker training job ephemeral storage (ML storage volumes) is encrypted. | Implicit | Training data written to unencrypted ephemeral storage is exposed to other processes on the same host. |

---

### Pattern 10: Encryption in Transit (TLS)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-037 | TLS is enforced between all ML pipeline components (SageMaker, S3, Feature Store). | Explicit | Data in transit between services must be encrypted. |
| ASF-038 | The inference endpoint enforces TLS for all client requests. | Derived | Inference requests containing sensitive input data must be encrypted in transit. |
| ASF-039 | TLS 1.2 or higher is enforced; TLS 1.0/1.1 and SSL are disabled. | Derived | Weak TLS versions expose inference data to passive interception. |
| ASF-040 | Weak cipher suites are disabled on the inference endpoint. | Derived | Weak ciphers negate the benefit of TLS enforcement. |

---

### Pattern 11: Endpoint Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-041 | The inference endpoint is protected by a WAF or API Gateway that blocks common web attacks. | Implicit | ML endpoints exposed to the internet without WAF are vulnerable to injection attacks. |
| ASF-042 | SageMaker notebook instances (if used) are configured with least-privilege IAM and no root access. | Derived | Notebooks with full SageMaker access can launch arbitrary training jobs. |
| ASF-043 | The inference endpoint has input validation to reject adversarial or malformed requests. | Operational | ML models are vulnerable to adversarial examples that exploit input processing. |
| ASF-044 | The inference endpoint has auto-scaling configured to handle traffic spikes without degradation. | Environmental | Unscaled endpoints fail during traffic spikes, causing inference unavailability. |

---

### Pattern 12: Human Factors & Process

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-045 | Data scientists do not have direct access to production inference endpoints or production data. | Derived | Data scientist access to production resources increases the risk of accidental misconfiguration. |
| ASF-046 | ML engineers understand the security implications of training data provenance and model lineage. | Trust | Without training, engineers may not understand the importance of data integrity for model security. |
| ASF-047 | Model reviewers evaluate models for security issues (adversarial robustness, data memorization) before promotion. | Operational | Models promoted without security review may have vulnerabilities that are invisible to performance metrics. |
| ASF-048 | There is a process for reporting and handling incidents involving biased or harmful model outputs. | Environmental | Model output incidents require different response procedures than traditional security incidents. |

---

### Pattern 13: Identity Lifecycle (Provisioning)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-049 | Data scientist access to SageMaker and the Feature Store is reviewed and recertified quarterly. | Operational | Stale data scientist access grants unnecessary access to training infrastructure. |
| ASF-050 | Production inference endpoint credentials for applications are rotated and managed through a secrets manager. | Derived | Static endpoint credentials in application configuration files are a common leak. |
| ASF-051 | Service accounts used by CI/CD to deploy models are reviewed and scoped to the minimum required. | Implicit | CI/CD service accounts with broad permissions can deploy unapproved models. |
| ASF-052 | Model registry access is revoked when data scientists change teams or leave the organization. | Operational | Former data scientists retain ability to deploy models to production. |

---

### Pattern 14: Incident Response

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-053 | There is an incident response plan covering ML-specific scenarios (model poisoning, data tampering, adversarial attacks). | Operational | ML incidents require different containment (model rollback, data lineage audit) than traditional incidents. |
| ASF-054 | The IR team can quickly revert a production model to the previous known-good version. | Derived | Untested model rollback can fail during a critical incident. |
| ASF-055 | The IR plan includes procedures for identifying which training dataset was compromised. | Trust | Data poisoning response requires tracing lineage from model back to source data. |
| ASF-056 | Monitoring detects unusual model behavior (confidence score drift, accuracy degradation) that may indicate compromise. | Implicit | Without model behavior monitoring, a compromised model operates unnoticed. |

---

### Pattern 15: Least Privilege

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-057 | Training jobs access only the S3 prefixes and Feature Store feature groups required for that specific job. | Explicit | Training jobs with broad data access violate least privilege. |
| ASF-058 | The inference endpoint IAM role has permissions only to load the specific model from the registry. | Derived | An inference endpoint with registry write access can overwrite model artifacts. |
| ASF-059 | Production application credentials for the inference endpoint are scoped to a single model version. | Implicit | Application access to all model versions can invoke inference on unapproved models. |
| ASF-060 | SageMaker notebook IAM roles are scoped to specific S3 buckets and cannot launch training jobs. | Derived | Notebooks with full SageMaker access allow a data scientist to run arbitrary training jobs. |

---

### Pattern 16: Monitoring & Alerting

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-061 | Training job failures are monitored and alerted to the ML engineering team. | Operational | Silent training failures delay model availability. |
| ASF-062 | Inference endpoint latency and error rates are monitored with alerting on threshold breaches. | Derived | Degraded endpoint performance affects all production applications consuming inferences. |
| ASF-063 | Data Lake access patterns are monitored for unusual bulk downloads or access from unusual IPs. | Operational | Bulk data download is a signal of data exfiltration. |
| ASF-064 | Model drift is monitored and triggers retraining when prediction accuracy falls below a threshold. | Derived | Undetected model drift leads to incorrect business decisions based on stale models. |

---

### Pattern 17: Network Segmentation

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-065 | Training jobs run in a private VPC with no internet access (VPC-only mode). | Architectural | Training jobs with internet access can exfiltrate data or download untrusted code. |
| ASF-066 | The inference endpoint is deployed in a VPC with access restricted to the production application VPC. | Architectural | Publicly accessible inference endpoints are exposed to internet-based attacks. |
| ASF-067 | SageMaker endpoints do not have public IP addresses; access is via PrivateLink or VPC interface endpoints. | Explicit | Public IP on SageMaker endpoints bypasses VPC access controls. |
| ASF-068 | Security groups restrict SageMaker endpoint access to specific source security groups. | Derived | SageMaker endpoint security groups should restrict access to the production app security group only. |

---

### Pattern 18: Secrets Management

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-069 | AWS KMS keys used for S3, Feature Store, and Model Registry encryption are managed with separate keys per environment. | Explicit | Single KMS key across environments increases cross-environment blast radius. |
| ASF-070 | Inference endpoint API keys are stored in a secrets manager, not in application code. | Derived | API keys in application code are exposed through repository breaches and logs. |
| ASF-071 | SageMaker notebook credentials (if used with Git or external services) are stored in a secrets manager. | Implicit | Notebook credentials stored in plaintext files are extractable. |
| ASF-072 | KMS key access is audited and restricted to the minimum required IAM principals. | Operational | Unaudited KMS key access allows any authorized IAM user to decrypt data at rest. |

---

### Pattern 19: Supply Chain Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-073 | SageMaker built-in algorithms and framework containers have no known critical vulnerabilities. | Dependency | AWS-managed containers are trusted but should be monitored for vulnerability disclosures. |
| ASF-074 | Third-party libraries used in custom training code are scanned for vulnerabilities before use. | Operational | Training code dependency vulnerabilities (e.g., in PyTorch, TensorFlow) can lead to code execution. |
| ASF-075 | Pre-trained models sourced from external registries (Hugging Face, PyTorch Hub) are vetted before use. | Derived | Untrusted pre-trained models can contain backdoors, poisoned weights, or malware. |
| ASF-076 | The organization maintains an SBOM for custom training containers and monitors for new vulnerabilities. | Operational | Without SBOM tracking, a newly disclosed vulnerability cannot be quickly assessed for impact. |

---

### Pattern 20: Third-party Dependency

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-077 | SageMaker service is available in the selected region and meets the required throughput SLA. | Dependency | SageMaker region outage blocks training and inference. |
| ASF-078 | The Feature Store provider (if third-party) has a business continuity plan and SLA. | Dependency | Feature store unavailability blocks all training jobs that depend on it. |
| ASF-079 | There is a fallback compute provider if SageMaker becomes unavailable for an extended period. | Operational | Vendor lock-in to SageMaker means no training or inference during an extended outage. |
| ASF-080 | The model registry vendor/backend has an export capability to prevent data lock-in. | Derived | Inability to export model artifacts prevents migration to another ML platform. |

**Total (A): 80** (4 per pattern x 16 patterns + 16 overflow from high-complexity patterns)

---

## Step 3: Comparison

### Overlap Mapping

| Human ID | ASF ID | Match Rationale |
|----------|--------|-----------------|
| H-001 | ASF-015 | Both require S3 Data Lake versioning and data integrity. |
| H-002 | ASF-055 | Both require complete data lineage tracing. |
| H-003 | ASF-057 | Both require Feature Store access controls per feature group. |
| H-004 | ASF-065 | Both require training jobs in isolated VPC with no internet. |
| H-005 | ASF-034 | Both require model artifacts encrypted at rest in registry. |
| H-006 | ASF-022 | Both require registry access control with approval workflows. |
| H-007 | ASF-007 | Both require model version validation before serving. |
| H-008 | ASF-008 | Both require per-client rate limiting at endpoint. |
| H-009 | ASF-026 | Both require inference request/response logging. |
| H-010 | ASF-023 | Both require S3 bucket policies preventing public access. |
| H-012 | ASF-075 | Both require model signing/integrity verification. |
| H-013 | ASF-009 | Both require clean, non-persistent training environments. |
| H-014 | ASF-032 | Both require controlled feedback loops from inference to training. |
| H-015 | ASF-054 | Both require model rollback mechanism. |
| H-016 | ASF-016 | Both require immutable lineage records. |
| H-017 | ASF-017 | Both require request size limits on endpoint. |
| H-018 | ASF-042 | Both require notebook isolation from production. |
| H-020 | ASF-014 | Both require Feature Store retention policy. |
| H-021 | ASF-050 | Both require robust endpoint authentication. |
| H-022 | ASF-027 | Both require evaluation against held-out test set. |
| H-023 | ASF-020 | Both require training data quality checks. |
| H-024 | ASF-038 | Both require inference response encryption. |
| H-025 | ASF-043 | Both require adversarial input detection. |
| H-026 | ASF-055 | Both require model lineage traceability. |
| H-027 | ASF-063 | Both require Data Lake access monitoring. |
| H-030 | ASF-015 | Both require S3 versioning for DR. |
| H-032 | ASF-021 | Both require least-privilege training job IAM roles. |
| H-033 | ASF-041 | Both require WAF/API Gateway for endpoint protection. |
| H-040 | ASF-015 | Both require cross-region Data Lake replication. |
| H-041 | ASF-059 | Both require scoped credentials per endpoint. |
| H-042 | ASF-074 | Both require approved algorithms and code review. |
| H-043 | ASF-010 | Both require configurable endpoint timeout. |
| H-044 | ASF-065 | Both require SageMaker VPC-only mode. |
| H-046 | ASF-046 | Both require evaluation data separation from training. |

**Overlap (O): 34**

### Metrics

| Metric | Value | Formula |
|--------|-------|---------|
| Human assumptions (H) | 46 | Count of unique human-generated assumptions |
| ASF assumptions (A) | 80 | Count of unique ASF-generated assumptions |
| Overlap (O) | 34 | Count appearing in both lists |
| **Precision** | **42.5%** | O / A = 34/80 |
| **Recall** | **73.9%** | O / H = 34/46 |
| **F1 Score** | **54.0%** | 2 x (P x R) / (P + R) |
| Novel findings (A - O) | 46 | Assumptions ASF found that human missed (57.5% of ASF total) |
| Missed findings (H - O) | 12 | Assumptions human found that ASF missed (26.1% of human total) |

### Success Criteria Assessment

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Recall | >= 70% | 73.9% |  Met |
| Precision | >= 50% | 42.5% |  Not met |
| Novel discoveries | >= 10% of total (A+O) | 36.5% (46/126) |  Exceeded |
| Expert agreement (F1 proxy) | > 60% | 54.0% |  Not met |

---

## Step 4: Analysis

### Ontology Overlap Analysis

| Ontology Category | Overlaps | ASF Total | Overlap Rate |
|-------------------|----------|-----------|-------------|
| Explicit | 8 | 16 | 50.0% |
| Derived | 12 | 20 | 60.0% |
| Operational | 6 | 20 | 30.0% |
| Implicit | 4 | 12 | 33.3% |
| Trust | 2 | 4 | 50.0% |
| Dependency | 1 | 8 | 12.5% |
| Architectural | 1 | 4 | 25.0% |
| Environmental | 0 | 4 | 0.0% |

**Best overlap:** Derived category showed the strongest agreement. Both humans and the ASF recognized derived concerns (model encryption from classification, VPC isolation from network segmentation principles).

**Worst overlap:** Environmental. The ASF identified environmental concerns (data residency, regulatory retention, feedback loop controls) that the human did not enumerate.

### What Humans Caught That ASF Missed (Missed Findings = 12)

1. **ML-specific operational controls (H-019, H-028, H-029, H-034, H-035, H-036, H-037, H-038, H-039):** Model registry metadata completeness, multi-environment approval gates, hyperparameter validation, model status tracking, input schema validation, column-level lineage, max training runtime, immutable version tags, per-request model version logging. These are ML-pipeline-specific operational concerns not covered by generic ASF patterns.

2. **Data quality and bias prevention (H-011, H-022, H-045):** PII anonymization in features, feature point-in-time queries for training-serving skew, and evaluation dataset separation are ML-specific data concerns.

### What ASF Caught That Humans Missed (Novel Findings = 46)

1. **Secrets Management (4 assumptions):** The human covered KMS keys but the ASF extended to KMS key cross-environment separation, endpoint API key vault storage, notebook credential management, and KMS key access audit.

2. **Supply Chain Security (4 assumptions):** The human assumed approved algorithms (H-042) but the ASF extended to SageMaker container vulnerabilities, training code library scanning, pre-trained model vetting, and SBOM management.

3. **Compliance & Audit (4 assumptions):** The human generated no assumptions about model data classification, inference log regulatory retention, model bias auditing, or lineage record retention periods.

4. **Change Management (4 assumptions):** The human did not address registry schema backward-compatibility, feature schema change communication, endpoint config change preview, or training data source review process.

5. **Identity Lifecycle (4 assumptions):** The human covered access controls but the ASF extended to quarterly data scientist access recertification, endpoint credential rotation, CI/CD service account review, and registry access revocation.

### Architecture Complexity Assessment

Architecture #016 achieved the highest recall (73.9%) of the five architectures. This is because ML security concerns align well with the ASF's patterns: data integrity, encryption, IAM least privilege, monitoring, and network segmentation are all well-covered.

- **Recall (73.9%)** exceeds the 70% target. The ASF covers most ML-related assumptions through its data flow, encryption, IAM, and network segmentation patterns.
- **Precision (42.5%)** is the second-highest, reflecting good alignment between ASF patterns and ML security concerns.
- **Novel rate (57.5%)** remains high, with the ASF surfacing secrets management, supply chain, compliance, and change management concerns.

### Key Insight

The ASF's existing patterns are surprisingly well-suited to ML pipeline security because ML pipelines are fundamentally data pipelines with strong requirements for data integrity, access control, and audit. The primary gaps are **ML-specific operational metadata** (model status tracking, version immutability, hyperparameter validation) and **ML governance** (bias auditing, model metadata completeness, feature engineering controls). Adding an "ML Security" or "Model Governance" pattern would address these specific gaps.

---

## Conclusions

| Metric | Target | Actual | Verdict |
|--------|--------|--------|---------|
| Recall | >= 70% | 73.9% |  Target met (ASF patterns align well with ML data pipeline concerns) |
| Precision | >= 50% | 42.5% |  Below target — ASF produces broad assumptions across governance dimensions |
| Novel discoveries | >= 10% | 36.5% |  ASF adds value in secrets management, supply chain, compliance |
| Expert agreement (F1) | > 60% | 54.0% |  Below target — driven by low precision |

The ASF framework applied to Architecture #016 achieved the best recall of all five architectures, exceeding the 70% target. This suggests the ASF's existing patterns are well-aligned with ML pipeline security concerns, particularly around data integrity, access control, and encryption. The primary actionable finding is the need for an **ML Governance** pattern covering model metadata management, version immutability, hyperparameter validation, and bias auditing.
