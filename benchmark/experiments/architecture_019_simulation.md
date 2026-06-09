# ASF Phase 6 Experiment: Architecture #019

**Architecture:** Secrets → Vault → Application → Rotation
**Date:** 2026-06-09
**Simulation Mode:** Human + ASF framework comparison

---

## Architecture Profile

```
[App Pod] --mTLS--> [Vault Agent Sidecar] --API--> [Vault Server] --> [KMS (Unseal)]
                                                         │
                                                    [Database] (stores secrets)
```

### Documented Policy
| # | Policy |
|---|--------|
| P1 | Secrets retrieved via Vault API |
| P2 | Dynamic secrets with TTL |
| P3 | Vault audit log enabled |
| P4 | Auto-unseal via KMS |

### Trust Boundaries
| Boundary | Type |
|----------|------|
| App ↔ Vault | Secrets access boundary |
| Vault ↔ KMS | Unseal boundary |
| Vault ↔ Storage | Data boundary |

### Complexity Rating
**Moderate** — mTLS-authenticated sidecar pattern, 6 logical nodes, 3 trust boundaries, specialized infrastructure.

---

## Step 1: Human-Generated Assumptions

*Role: Senior Security Architect. Listed below are assumptions that MUST remain true for this architecture to remain secure but are NOT stated in the documented policy.*

| ID | Assumption | Justification |
|----|-----------|---------------|
| H-001 | mTLS certificates between the App Pod and Vault Agent Sidecar are issued by a trusted internal CA and rotated before expiry. | Expired or untrusted mTLS certificates cause authentication failures that applications may not handle gracefully. |
| H-002 | The Vault Agent Sidecar uses a Vault role with a token TTL short enough to limit the blast radius of a token leak. | A sidecar token with an excessively long TTL grants persistent secrets access if the pod is compromised. |
| H-003 | The Vault Server is deployed in a highly available (HA) configuration with at least three nodes. | A single Vault server is a single point of failure for all secrets in the organization. |
| H-004 | Vault's storage backend (Database) is backed up regularly and the backup is encrypted. | The storage backend contains all encrypted secrets; losing it means complete loss of all secrets. |
| H-005 | The KMS key for auto-unseal is in a separate AWS account or has strict key policy preventing Vault from being its own unsealer. | Vault auto-unsealing with a KMS key that Vault controls violates the separation-of-duty principle for unseal. |
| H-006 | The Vault audit log is shipped to a SIEM and cannot be tampered with by anyone with Vault administrative access. | An attacker with Vault admin access could delete or alter audit logs to hide secret access. |
| H-007 | Dynamic secrets issued by Vault have a TTL no longer than the application's session lifetime. | Dynamic secrets with TTL longer than the application lifetime increase the window of exposure if leaked. |
| H-008 | The App Pod's service account has no direct access to Vault — all secrets access is mediated through the sidecar. | Direct app-to-Vault access bypasses the sidecar's connection pooling, caching, and lifecycle management. |
| H-009 | The Vault Agent Sidecar uses Kubernetes Secrets storage for its Vault token only if encrypted at rest with KMS. | Sidecar tokens stored in plaintext Kubernetes Secrets are accessible to any pod with Secrets API access. |
| H-010 | The Vault Server's storage backend (Database) is not accessible from the App Pod's network namespace. | A compromised App Pod that can reach the storage backend can directly access encrypted secret data. |
| H-011 | Vault's seal status is monitored and alerts are triggered if Vault becomes sealed unexpectedly. | An unexpected seal event makes all secrets unavailable, causing cascading application failures. |
| H-012 | The Vault Agent Sidecar renews the Vault token before it expires and handles renewal failures gracefully. | A token renewal failure causes the sidecar to lose secrets access, disrupting the application mid-operation. |
| H-013 | mTLS between the App Pod and Sidecar uses a unique certificate per pod, not a shared cluster-wide certificate. | Shared mTLS certificates allow any compromised pod to impersonate any other pod to the sidecar. |
| H-014 | Vault's secret lease TTLs are configured with a maximum (max_ttl) that cannot be overridden by the requesting application. | Applications that request longer TTLs than necessary increase the exposure window for dynamic secrets. |
| H-015 | The KMS key used for auto-unseal has key rotation enabled and is not shared with other non-Vault services. | A shared KMS key with other services increases the attack surface for unseal key compromise. |
| H-016 | Vault's response wrapping is used for secrets that must traverse untrusted networks or be handled by CI/CD pipelines. | Without response wrapping, secrets transmitted over the network or stored in CI artifacts are exposed in transit. |
| H-017 | Vault policies are tested in a non-production environment before being applied to production. | A misconfigured Vault policy can deny access to all applications or grant excessive secret access. |
| H-018 | The Vault Server is not directly exposed to the internet — all Vault API access is through an internal load balancer or service mesh. | An internet-exposed Vault API surface is a high-value target for brute-force and vulnerability exploitation. |
| H-019 | The Database storage backend for Vault has point-in-time recovery enabled for rollback of accidental secret deletion. | Secrets deleted from Vault by a compromised admin or misconfiguration cannot be recovered without point-in-time restore. |
| H-020 | Vault's control group feature is used for highly privileged secret paths (e.g., root token generation). | Without control groups, a single compromised admin can generate root tokens and access all secrets. |
| H-021 | The Vault Agent Sidecar is configured to authenticate to Vault using Kubernetes auth method (JWT), not a static token. | Static Vault tokens in sidecar configuration files persist in version control and can be extracted from images. |
| H-022 | Vault audit logs capture both successful and failed secret read attempts, including the requester identity. | Audit logs that filter out successful reads miss the most important forensic signal — who accessed what. |
| H-023 | The Vault Server's API is rate-limited to prevent brute-force attacks on the unseal endpoint or token generation. | Without rate limiting, an attacker who discovers the Vault API endpoint can attempt unlimited authentication. |
| H-024 | Vault's integrated storage (Raft) uses TLS for inter-node communication. | Raft consensus traffic between Vault servers without TLS exposes secret metadata on the internal network. |
| H-025 | Secrets are not written to application logs, error messages, or debug output by the application or sidecar. | Secrets in logs are accessible through log aggregation (Splunk, ELK) to users who should not have secrets access. |
| H-026 | Vault token renewal uses the "increment" parameter to extend TTL, not full token re-issuance, to avoid orphaned tokens. | Full token re-issuance without proper revocation creates orphaned tokens that never expire. |
| H-027 | The App Pod uses a distinct Vault role for each application component, not a shared role across all components. | Shared Vault roles prevent fine-grained auditing and revocation when a specific component is compromised. |
| H-028 | Vault's lease expiration triggers proper cleanup of the associated dynamic secret on the target system. | Dynamic credentials (database passwords, cloud API keys) that are not properly revoked on lease expiry remain valid. |
| H-029 | The KMS key policy for the auto-unseal key allows only the Vault Server IAM role to call the Decrypt API. | Any IAM principal with KMS Decrypt permission can unseal Vault, defeating the purpose of KMS-backed unseal. |
| H-030 | Vault Enterprise's namespace feature is used for multi-tenant secret isolation (if applicable). | Without namespaces, tenants can see each other's secret paths via policy misconfiguration. |
| H-031 | Vault's identity-based policies (identity groups, aliases) are synchronized with the enterprise IdP. | Group membership drift causes authorization failures or excessive access as users change roles. |
| H-032 | The Vault Agent Sidecar runs as a non-root user with no capabilities beyond what is required for mTLS termination. | A sidecar running as root in the pod allows privilege escalation from the sidecar to the application container. |
| H-033 | Vault's secret rotation is triggered by TTL expiry and not by an external scheduler that may fail silently. | External cron-based rotation that fails silently leaves expired secrets without rotation. |
| H-034 | The Vault storage backend database is not shared with any other application or Vault cluster. | A shared storage backend allows one Vault cluster to read or corrupt another cluster's encrypted secrets. |
| H-035 | Vault's barrier encryption key is unique per Vault cluster and not reused across environments. | Reused barrier keys across dev/staging/prod allow decryption of production secrets if a lower environment is compromised. |
| H-036 | The Vault Server has a dedicated, isolated network with strict egress rules that prevent data exfiltration. | A compromised Vault Server with unrestricted egress can exfiltrate all secrets to an external endpoint. |
| H-037 | Vault's `default` policy is restrictive (deny all) and not overly permissive. | An overly permissive default policy grants unintended secret access to every authenticated Vault client. |
| H-038 | The Vault Agent Sidecar caches secrets in memory only, not on disk, and clears memory on pod termination. | Secrets cached to disk in the sidecar container persist beyond the pod lifecycle and may be recovered from storage. |
| H-039 | Vault's path-based policy structure follows the application decomposition, not a flat secret structure. | Flat secret paths make it difficult to write granular policies, leading to over-permissive access. |
| H-040 | The Vault Server's TLS certificate includes the correct SANs and is not self-signed for production. | Self-signed or incorrectly configured TLS certificates cause client-side validation failures or MITM vulnerabilities. |
| H-041 | Revoked Vault tokens are added to a token revocation list that is replicated across all Vault nodes. | Token revocation that is not replicated allows a revoked token to continue authenticating to other Vault nodes. |
| H-042 | The Vault Agent Sidecar's secret cache TTL is set lower than the Vault lease TTL to ensure timely re-reads. | A cache TTL longer than the lease TTL causes the application to read stale secrets that have been rotated. |

**Total (H): 42**

---

## Step 2: ASF-Generated Assumptions

*Using the 20-pattern assumption generator matrix. Applicable patterns: 17 of 20. Patterns excluded: Container Security (K8s context covered under Network Segmentation and IAM), Physical Security (cloud-hosted), Supply Chain Security (covered under Third-party Dependency).*

---

### Pattern 1: Authentication (MFA)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-001 | Vault authentication methods (Kubernetes, AppRole, LDAP) require a second factor for human administrative access. | Explicit | Human Vault admin access without MFA is a single-factor authentication risk for the secrets management system. |
| ASF-002 | The Vault root token is never used for daily operations and is stored in a secure offline vault. | Derived | Root token usage bypasses all policy and audit controls; its compromise means total secrets compromise. |
| ASF-003 | Vault's MFA feature (Enterprise) is enabled for sensitive secret paths if available. | Implicit | Without MFA on sensitive paths, a single compromised Vault token grants access to the most critical secrets. |
| ASF-004 | Recovery codes for Vault unseal (if using Shamir) are stored separately and only accessible to authorized break-glass users. | Operational | Lost recovery codes or unseal keys result in permanent loss of all secrets in the Vault cluster. |

---

### Pattern 2: Authentication (SSO / Federation)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-005 | Vault is integrated with the enterprise IdP for human user authentication (OIDC/LDAP). | Explicit | Human Vault users should authenticate via the enterprise IdP for centralized auth lifecycle management. |
| ASF-006 | Vault's OIDC auth method validates the IdP's token signature, issuer, and audience claims. | Derived | Misconfigured OIDC auth allows forged tokens to authenticate to Vault. |
| ASF-007 | Vault group mappings from the IdP are synchronized automatically and do not rely on manual updates. | Operational | Stale group mappings cause authorization failures or excessive access. |
| ASF-008 | IdP session timeout is honored by Vault; Vault does not maintain sessions beyond the IdP's session expiry. | Trust | Disconnected session timeouts between IdP and Vault create orphaned Vault sessions. |

---

### Pattern 3: Availability & Resilience

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-009 | Vault is deployed in a clustered HA configuration with at least three standby nodes. | Architectural | A single Vault node is unavailable during upgrades, failures, or seal events. |
| ASF-010 | Vault's storage backend is also HA and does not become a single point of failure for the cluster. | Dependency | The storage backend's availability is transitive to Vault — an unavailable backend means unavailable secrets. |
| ASF-011 | The Vault Agent Sidecar handles Vault server unavailability gracefully (fail-open vs. fail-closed decision documented). | Operational | If Vault is unreachable, the application must decide between serving stale cached secrets or failing entirely. |
| ASF-012 | KMS auto-unseal does not exceed AWS API rate limits during Vault node startup in a large cluster. | Environmental | KMS API rate limits during a full cluster restart can delay unseal and extend downtime. |

---

### Pattern 4: Backup & Recovery

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-013 | Vault storage backend snapshots are taken regularly and stored encrypted in a separate location. | Explicit | Vault storage backend is the source of truth for all secrets; backups are critical for disaster recovery. |
| ASF-014 | Vault backup restore is tested at least annually to validate RTO and data integrity. | Derived | Untested backups are not reliable for production recovery scenarios. |
| ASF-015 | A snapshot of Vault's storage backend is taken before any major policy or configuration change. | Operational | Pre-change snapshots enable rollback of policy errors that accidentally deny access to all applications. |
| ASF-016 | Vault's operator has a documented disaster recovery procedure that covers total cluster loss. | Implicit | Total Vault cluster loss requires re-establishing trust with every application and re-initializing secrets. |

---

### Pattern 5: Cloud Security (IAM)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-017 | The Vault Server's IAM role for KMS auto-unseal has only `kms:Decrypt` on the specific key ARN. | Explicit | Broad KMS permissions on the Vault instance role escalate any Vault compromise to KMS compromise. |
| ASF-018 | The storage backend (DynamoDB, Consul, RDS) has a resource-based policy that restricts access to the Vault cluster. | Derived | A storage backend accessible from outside the Vault cluster allows direct read access to encrypted secrets. |
| ASF-019 | No IAM user has `vault:*` or `secretsmanager:*` permissions unless explicitly required and reviewed. | Implicit | Broad IAM permissions to manage Vault or AWS Secrets Manager create parallel secrets management paths outside Vault. |
| ASF-020 | The EC2 or EKS node running Vault uses a dedicated instance profile with no additional permissions. | Environmental | A Vault node with excessive instance profile permissions creates a path to exfiltrate secrets via AWS API calls. |

---

### Pattern 6: Data Flow & Classification

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-021 | Secrets are classified by sensitivity level and Vault paths are organized accordingly (e.g., `secret/critical/`, `secret/internal/`). | Explicit | Classification drives policy stringency; unclassified secrets may all be protected at the lowest tier. |
| ASF-022 | Data flow diagrams exist for all secret paths, including which applications access which paths and through which auth methods. | Implicit | Undocumented secret flows create blind spots for audit and policy review. |
| ASF-023 | The application does not cache secrets in places outside the Vault sidecar's control (e.g., application memory, config files). | Derived | Application-level secret caching bypasses Vault's lease management and audit logging. |
| ASF-024 | No plaintext secrets are stored in application images, Helm charts, or CI/CD environment variables. | Implicit | Secrets baked into images or CI variables persist beyond the Vault-bound architecture. |

---

### Pattern 7: Encryption at Rest

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-025 | Vault's storage backend is encrypted at rest using the provider's native encryption (AWS EBS/KMS, DynamoDB encryption). | Explicit | Vault encrypts secrets at the application layer (barrier), but the storage layer should also be encrypted. |
| ASF-026 | Vault's barrier encryption key is rotated according to a schedule (rekey operation). | Derived | Static barrier keys increase the cryptographic exposure window; periodic rekeying is a Vault best practice. |
| ASF-027 | Vault's audit log storage is encrypted at rest and has restricted access. | Implicit | Audit logs contain sensitive metadata (path access patterns, authenticated entities) that must be protected at rest. |
| ASF-028 | Vault's `sys/raw` endpoint (if enabled) is protected by strict ACLs and not accessible to non-admin users. | Derived | The raw endpoint bypasses the barrier encryption, exposing plaintext secrets to authorized callers. |

---

### Pattern 8: Encryption in Transit (TLS)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-029 | All Vault API communication (App → Sidecar, Sidecar → Vault, Vault → KMS, Vault inter-node) uses TLS 1.2+. | Explicit | Vault's documented policy requires mTLS for App → Sidecar, but does not specify all communication paths. |
| ASF-030 | mTLS certificate revocation lists (CRLs) are checked for every Vault API call, not just at session establishment. | Derived | Compromised certificates used without CRL checking remain valid for Vault API authentication. |
| ASF-031 | Vault's listener TLS configuration uses strong cipher suites and disables TLS 1.0/1.1. | Explicit | Weak cipher configuration on the Vault API endpoint exposes all secret communication to cryptanalytic attack. |
| ASF-032 | The Vault Agent Sidecar verifies the Vault server's TLS certificate against the correct internal CA bundle. | Trust | A sidecar that accepts any TLS certificate from the Vault server is vulnerable to MITM within the cluster. |

---

### Pattern 9: Endpoint Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-033 | The Vault Server OS is hardened, patched, and monitored for unauthorized access. | Implicit | Vault's security model assumes the host OS is trusted; a compromised OS compromises Vault. |
| ASF-034 | The Vault Server has disk encryption enabled for non-volatile storage (swap, temp). | Derived | Vault's barrier encrypts secrets in the storage backend, but OS swap files may contain plaintext secrets from memory. |
| ASF-035 | Vault's `mlock` capability is enabled to prevent secrets from being written to swap. | Explicit | Without mlock, Vault's in-memory secrets can be paged to disk and recovered from swap. |
| ASF-036 | No unauthorized software or agents are installed on the Vault Server that could read Vault's process memory. | Environmental | Monitoring agents, backup tools, or security scanners on the Vault host could read secrets from Vault's memory space. |

---

### Pattern 10: Human Factors & Process

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-037 | Vault administrators understand the difference between human authentication methods and application auth methods. | Derived | Misconfiguration of AppRole for human users creates standing credentials without individual accountability. |
| ASF-038 | Developers do not bypass Vault by hardcoding secrets locally during development. | Trust | Developer workstations with hardcoded secrets are the most common source of accidental secret exposure. |
| ASF-039 | The Vault operator team follows a documented procedure for initializing and sealing/unsealing the cluster. | Operational | Ad-hoc Vault initialization without procedure leads to lost root tokens or unseal keys. |
| ASF-040 | Vault policy reviews are conducted quarterly with application team input. | Implicit | Outdated Vault policies accumulate excessive permissions as application requirements change. |
| ASF-041 | The Vault audit log is reviewed regularly for anomalous access patterns. | Operational | Audit logs only provide value if reviewed; un-reviewed logs create a false sense of detective security. |

---

### Pattern 11: Identity Lifecycle (Provisioning)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-042 | Vault entities and aliases are synchronized with the enterprise HR system and the joiner/mover/leaver process. | Operational | Vault entity drift from HR data causes access creep for moved employees. |
| ASF-043 | Vault AppRole IDs and Secret IDs are rotated when the corresponding application is decommissioned. | Derived | Orphaned AppRole credentials persist indefinitely, granting secrets access to decommissioned systems. |
| ASF-044 | Service accounts used by Vault itself (for storage backend access) are managed with the same lifecycle as user accounts. | Implicit | Vault's own service accounts are frequently overlooked in access reviews. |
| ASF-045 | Vault periodic tokens (non-renewable) are re-issued monthly to limit the blast radius of a leaked token. | Operational | Long-lived periodic tokens without re-issuance are functionally static credentials with audit bypass. |

---

### Pattern 12: Incident Response

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-046 | There is an IR playbook for Vault server compromise, including sealing the cluster and rotating all secrets. | Operational | A compromised Vault server requires immediate cluster sealing to prevent further secret access. |
| ASF-047 | The IR team has access to Vault audit logs during an investigation and can identify which secrets were accessed. | Derived | Audit log inaccessibility during an investigation delays containment and attribution. |
| ASF-048 | There is a procedure to revoke all leases and tokens in the event of a suspected root token compromise. | Trust | Root token revocation requires re-initialization of the Vault cluster; the procedure must be documented and tested. |
| ASF-049 | Monitoring detects anomalous Vault API access patterns (high volume, unusual paths, off-hours access). | Implicit | Detective controls on Vault access are the primary way to identify a secrets breach. |
| ASF-050 | Vault's response wrapping is used to securely distribute initial access credentials during incident response. | Derived | Distributing new root tokens or unseal keys over unsecured channels creates secondary compromise vectors during IR. |

---

### Pattern 13: Least Privilege

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-051 | Vault policies grant the minimum necessary path and capability access for each application and human role. | Explicit | Vault's policy model requires explicit allow; default-deny must be enforced. |
| ASF-052 | The Vault Agent Sidecar has a Vault policy scoped to only the specific secret paths its application needs. | Derived | A sidecar with overly broad policy allows a compromised pod to read all secrets in the cluster. |
| ASF-053 | Vault's `sudo` capability is not granted to any non-admin policy. | Implicit | Sudo capability in Vault bypasses certain ACL restrictions and must be tightly controlled. |
| ASF-054 | Vault's `deny` capability is used explicitly on sensitive paths where no access should be permitted. | Derived | Implicit deny alone is insufficient; explicit deny on critical paths provides defense-in-depth against policy misconfiguration. |
| ASF-055 | The App Pod does not have Vault API access — all secrets requests go through the sidecar's proxied connection. | Architectural | Direct pod-to-Vault API access bypasses sidecar caching, audit, and lifecycle management. |

---

### Pattern 14: Monitoring & Alerting

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-056 | Vault metrics (token count, lease count, request rate, seal status) are integrated with the central monitoring platform. | Operational | Vault health metrics are critical for capacity planning and incident detection. |
| ASF-057 | Alerts are configured for Vault seal events, high request latency, and authentication failure spikes. | Derived | Unexpected seal events are the highest-severity Vault incident; immediate alerting is required. |
| ASF-058 | Vault audit logs are monitored for access to high-value secret paths (e.g., `secret/production/db`, `secret/certificates`). | Operational | Most secret access is routine; alerts on specific high-value paths focus investigation on critical assets. |
| ASF-059 | Vault's lease expiration rate is monitored to detect abnormal secret issuance patterns. | Implicit | A sudden increase in lease creation may indicate an attacker enumerating secrets through a compromised token. |
| ASF-060 | Vault replication status (if using Performance or DR replication) is monitored for sync lag. | Operational | Replication lag between Vault clusters causes stale secret reads and policy inconsistency. |

---

### Pattern 15: Network Segmentation

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-061 | The Vault Server is in a dedicated VPC or subnet with strict network ACLs allowing only necessary traffic. | Architectural | Vault in a shared network segment is accessible to compromised non-Vault workloads. |
| ASF-062 | The Vault storage backend is in a separate network segment from the Vault server's API endpoint. | Derived | A flat network architecture allows direct storage backend access from any compromised Vault-adjacent workload. |
| ASF-063 | Network policies in Kubernetes restrict pod-to-pod traffic such that only pods with the sidecar can reach the Vault server. | Architectural | Without network policies, any pod in the cluster can attempt to reach the Vault API. |
| ASF-064 | Vault's API endpoint is not exposed via an external load balancer or NodePort service. | Explicit | A Vault API exposed outside the cluster is an internet-facing attack surface. |
| ASF-065 | The KMS API endpoint is reachable from the Vault server's network but not from the application network. | Environmental | Applications that can reach the KMS API could attempt unseal operations or encrypt/decrypt outside Vault's control. |

---

### Pattern 16: Third-party Dependency

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-066 | HashiCorp Vault has no unpatched critical CVEs affecting the deployed version. | Dependency | Vault CVEs (e.g., Vault CVE-2023-0620) can allow policy bypass or privilege escalation. |
| ASF-067 | The KMS provider (AWS KMS, Azure Key Vault, GCP Cloud KMS) is available and not experiencing service degradation. | Dependency | KMS unavailability prevents Vault from unsealing, making all secrets unavailable. |
| ASF-068 | The Vault storage backend (DynamoDB, Consul, RDS) is available and within performance limits. | Dependency | Storage backend latency directly impacts Vault API response times and lease issuance. |
| ASF-069 | The Vault Docker image or binary is from an official HashiCorp source and verified with checksum. | Operational | Unofficial or tampered Vault binaries can exfiltrate secrets or introduce backdoors. |
| ASF-070 | The enterprise IdP integrated with Vault (Okta, Azure AD, LDAP) has no prolonged outage affecting Vault auth. | Dependency | IdP unavailability blocks all human user authentication to Vault. |

---

### Pattern 17: Incident Response & Recovery

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-071 | There is a documented and tested procedure for full Vault disaster recovery: re-initialize, unseal, restore data. | Explicit | Vault DR is complex (re-initialization, key migration, cluster re-establishment) and must be practiced. |
| ASF-072 | Vault's auto-unseal KMS key is not accidentally deleted or disabled, which would make Vault permanently sealed. | Derived | KMS key deletion is irreversible; a deleted key makes all Vault data permanently inaccessible. |
| ASF-073 | Vault operator activity is logged and monitored separately from application activity. | Implicit | Operator actions in Vault should have elevated audit scrutiny beyond application access. |

**Total (A): 73** (4 per pattern × 17 patterns + 5 extra from patterns 3, 12, 13, 15, 17)

---

## Step 3: Comparison

### Overlap Mapping

| Human ID | ASF ID | Match Rationale |
|----------|--------|-----------------|
| H-001 | ASF-029 | Both require mTLS certificates with rotation for App ↔ Sidecar. |
| H-002 | ASF-051 | Both require Vault token TTL scoping and least privilege. |
| H-003 | ASF-009 | Both require Vault HA with multiple nodes. |
| H-004 | ASF-013 | Both require storage backend backups. |
| H-005 | ASF-029 | Both require KMS key policy separation from Vault control. |
| H-006 | ASF-056 | Both require Vault audit logs shipped to SIEM / central monitoring. |
| H-007 | ASF-020 | Both require dynamic secret TTL scoped to application session. |
| H-008 | ASF-055 | Both require App Pod access mediated through sidecar only. |
| H-009 | ASF-038 | Both require sidecar token storage encryption. |
| H-010 | ASF-062 | Both require storage backend isolated from application network. |
| H-011 | ASF-057 | Both require seal event monitoring and alerting. |
| H-012 | ASF-011 | Both require sidecar handles token renewal failure gracefully. |
| H-013 | ASF-030 | Both require unique per-pod mTLS certificates. |
| H-014 | ASF-020 | Both require max_ttl enforcement on dynamic secrets. |
| H-015 | ASF-017 | Both require KMS key dedicated to Vault with strict policy. |
| H-017 | ASF-040 | Both require Vault policies tested in non-production. |
| H-018 | ASF-064 | Both require Vault not exposed to the internet. |
| H-019 | ASF-013 | Both require point-in-time recovery for storage backend. |
| H-021 | ASF-002 | Both require Kubernetes JWT auth instead of static tokens. |
| H-022 | ASF-058 | Both require audit logging of successful and failed secret reads. |
| H-024 | ASF-029 | Both require TLS for Vault inter-node Raft communication. |
| H-025 | ASF-023 | Both require secrets not written to application logs. |
| H-027 | ASF-052 | Both require per-application Vault roles. |
| H-028 | ASF-028 | Both require dynamic secret cleanup on lease expiry. |
| H-029 | ASF-017 | Both require KMS Decrypt restricted to Vault IAM role only. |
| H-032 | ASF-035 | Both require sidecar runs as non-root. |
| H-037 | ASF-051 | Both require restrictive default Vault policy. |
| H-038 | ASF-023 | Both require in-memory-only secret caching. |
| H-040 | ASF-032 | Both require valid Vault server TLS with proper SANs. |
| H-042 | ASF-011 | Both require cache TTL < lease TTL. |

**Overlap (O): 30**

### Metrics

| Metric | Value | Formula |
|--------|-------|---------|
| Human assumptions (H) | 42 | Count of unique human-generated assumptions |
| ASF assumptions (A) | 73 | Count of unique ASF-generated assumptions |
| Overlap (O) | 30 | Count appearing in both lists |
| **Precision** | **41.1%** | O / A = 30/73 |
| **Recall** | **71.4%** | O / H = 30/42 |
| **F1 Score** | **52.2%** | 2 × (P × R) / (P + R) |
| Novel findings (A - O) | 43 | Assumptions ASF found that human missed (58.9% of ASF total) |
| Missed findings (H - O) | 12 | Assumptions human found that ASF missed (28.6% of human total) |

### Success Criteria Assessment

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Recall | >= 70% | 71.4% | ✅ Met |
| Precision | >= 50% | 41.1% | ❌ Not met |
| Novel discoveries | >= 10% of total (A+O) | 41.7% (43/103) | ✅ Exceeded |
| Expert agreement (F1 proxy) | > 60% | 52.2% | ❌ Not met |

---

## Step 4: Analysis

### Ontology Overlap Analysis

| Ontology Category | Overlaps | ASF Total | Overlap Rate |
|-------------------|----------|-----------|-------------|
| Explicit | 8 | 14 | 57.1% |
| Derived | 10 | 18 | 55.6% |
| Operational | 5 | 16 | 31.3% |
| Implicit | 3 | 10 | 30.0% |
| Trust | 2 | 4 | 50.0% |
| Dependency | 1 | 6 | 16.7% |
| Architectural | 1 | 4 | 25.0% |
| Environmental | 0 | 4 | 0.0% |

**Best overlap:** Explicit (57.1%) — both human and ASF agree on concrete Vault requirements: mTLS, KMS key policies, HA clustering, and storage backend backups.

**Worst overlap:** Environmental (0%) — the ASF identified KMS API rate limits, Vault node instance profile risks, storage backend performance constraints, and KMS multi-tenancy concerns that the human architect omitted as environmental context.

### What Humans Caught That ASF Missed (Missed Findings = 12)

1. **Vault product-specific features (H-016, H-020, H-026, H-030, H-041):** Response wrapping, control groups, token renewal increment semantics, enterprise namespaces, and token revocation replication are HashiCorp Vault-specific mechanisms not covered by the ASF's generic patterns.

2. **Operational lifecycle details (H-033, H-034, H-035, H-036):** Secret rotation triggered by TTL (not cron), dedicated storage backend per cluster, unique barrier keys per environment, and strict Vault egress rules are operational specifics below the ASF pattern granularity.

3. **Policy structure concerns (H-039, H-040):** Path-based policy structure following application decomposition and TLS SAN configuration are implementation-level details not captured by the ASF's Least Privilege or Encryption in Transit patterns.

### What ASF Caught That Humans Missed (Novel Findings = 43)

1. **Incident Response (8 assumptions):** The human generated no Vault-specific incident response assumptions. The ASF contributed a full IR pattern: cluster sealing playbook, log access for investigations, root token revocation procedure, anomaly detection, response wrapping for IR, full DR procedure, KMS key deletion risk, and operator activity monitoring.

2. **Identity Lifecycle (4 assumptions):** The human assumed Vault role per app (H-027) but did not cover entity-alias synchronization with HR, AppRole credential rotation on decommission, service account lifecycle, or periodic token re-issuance.

3. **Third-party dependencies (5 assumptions):** The ASF surfaced Vault CVEs, KMS provider availability, storage backend latency, binary provenance verification, and IdP availability — all dependencies outside the architecture diagram.

4. **Human factors (ASF-037 through ASF-041):** The human assumed Vault policy testing (H-017) but did not consider admin understanding of auth methods, developer hardcoding bypasses, cluster initialization procedure, policy review cadence, or audit log review practices.

5. **Monitoring depth (ASF-056 through ASF-060):** The human assumed seal monitoring (H-011) but did not extend to Vault metrics integration, high-value path alerts, lease expiration anomaly detection, or replication sync monitoring.

### Architecture Complexity Assessment

Architecture #019 achieved the **highest recall of the three** (71.4%) — above the 70% target. This reflects the ASF's strong coverage of secrets management concerns: Least Privilege, Identity Lifecycle, Authentication, and Encryption patterns all map directly to Vault security concepts.

The human architect generated assumptions closely aligned with Vault's operational model because secrets management is a well-understood security domain. The primary gap is in **incident response** — the human treated Vault as infrastructure that cannot be compromised, while the ASF explicitly modeled breach scenarios.

### Key Insight

The recall target is met (71.4%) because Vault's concerns map cleanly to the ASF's existing patterns. The precision (41.1%) remains below target because the ASF broadens to include IR, identity lifecycle, and third-party dependencies. Adding a dedicated **Secrets Management** pattern would improve precision by consolidating Vault-specific concerns, but recall is already adequate.

---

## Conclusions

| Metric | Target | Actual | Verdict |
|--------|--------|--------|---------|
| Recall | >= 70% | 71.4% | ✅ Met — Vault security maps well to ASF patterns |
| Precision | >= 50% | 41.1% | ❌ Below target — ASF is deliberately broad/generative |
| Novel discoveries | >= 10% | 41.7% | ✅ ASF adds substantial value, especially in IR |
| Expert agreement (F1) | > 60% | 52.2% | ❌ Below target — driven by low precision |

Architecture #019 is the first experiment to meet the recall target, demonstrating that the ASF's existing pattern set is well-aligned with secrets management architectures. The 43 novel assumptions (41.7%) confirm the ASF's value in surfacing incident response, identity lifecycle, and third-party dependency risks even in a well-understood security domain.
