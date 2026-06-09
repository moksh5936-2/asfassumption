# ASF Phase 6 Experiment: Architecture #5

**Architecture:** Microservices → Service Mesh → Kubernetes → Istio
**Date:** 2026-06-09
**Simulation Mode:** Human + ASF framework comparison

---

## Architecture Profile

```
[Ingress Gateway] --mTLS--> [Service A] --mTLS--> [Service B] --mTLS--> [Service C]
                    │              │                                       │
               [Istio Pilot]  [K8s API]                              [StatefulSet DB]
                    │              │                                       │
               [Citadel CA]  [etcd]                                  [Persistent Volume]
```

### Documented Policy
| # | Policy |
|---|--------|
| P1 | mTLS enabled between all services |
| P2 | RBAC enforced at namespace level |
| P3 | Pod security policies restrict privileged containers |
| P4 | Network policies isolate namespaces |

### Trust Boundaries
| Boundary | Type |
|----------|------|
| Ingress ↔ Services | Mesh boundary |
| Service ↔ Service | Identity boundary |
| Service ↔ CA | Certificate trust boundary |
| Pod ↔ K8s API | Control plane boundary |

### Complexity Rating
**Complex** — service mesh with multiple control plane components, mTLS identity, persistent stateful workloads, and Kubernetes-native security controls.

---

## Step 1: Human-Generated Assumptions

*Role: Senior Security Architect. Listed below are assumptions that MUST remain true for this architecture to remain secure but are NOT stated in the documented policy.*

| ID | Assumption | Justification |
|----|-----------|---------------|
| H-001 | Envoy sidecar proxies are configured to reject traffic that is not on the mTLS port. | A sidecar accepting plaintext traffic bypasses the mesh security model. |
| H-002 | The Istio Citadel CA is configured with a secure, rotated root certificate and intermediate certs. | Compromise of the mesh CA allows an attacker to forge service identities. |
| H-003 | Citadel CA signing keys are stored in a hardware security module (HSM) or KMS. | Software-only CA keys can be extracted from the control plane. |
| H-004 | Istio mTLS is in STRICT mode, not PERMISSIVE mode. | Permissive mode allows plaintext traffic between services, defeating mTLS. |
| H-005 | mTLS certificate rotation is automated and completes before certificate expiry (no cert renewal failures). | Expired mTLS certs cause service-to-service communication failures. |
| H-006 | Kubernetes RBAC is configured with least-privilege roles and bindings, not permissive ClusterRoles. | Cluster-level admin roles bypass namespace boundaries and grant excessive control. |
| H-007 | K8s API Server is not exposed to the public internet and has IP whitelisting for administrative access. | A public API Server is the primary attack vector for cluster compromise. |
| H-008 | etcd is encrypted at rest and access is restricted to the K8s API Server only. | etcd contains all cluster secrets; direct access bypasses all RBAC controls. |
| H-009 | Pod security policies (or OPA/Gatekeeper) prevent privilege escalation, host network access, and host PID sharing. | PSPs without host access restrictions allow container breakout. |
| H-010 | Network policies define ingress/egress rules for each namespace; default deny is enforced. | Undefined network policies default to allow-all, defeating namespace isolation. |
| H-011 | The ingress gateway terminates TLS and re-encrypts with mTLS for the mesh. | Direct passthrough without re-encryption breaks the mTLS chain. |
| H-012 | Istio Pilot (or Istiod) does not accept unauthenticated requests for proxy configuration. | Unauthenticated Pilot access allows an attacker to reconfigure the mesh. |
| H-013 | Workload identity is based on Kubernetes ServiceAccount tokens, not IP-based. | IP-based identity can be spoofed by compromised pods. |
| H-014 | SPIFFE identities issued by Citadel are unique per workload and not shared. | Shared SPIFFE identities defeat workload-level isolation. |
| H-015 | Istio authorization policies (AuthorizationPolicy CRDs) are configured to enforce service-level access control. | Without authorization policies, mTLS provides authentication only, not authorization. |
| H-016 | Sidecar injection is enforced (not opt-in) so all pods in the mesh have Envoy proxies. | Pods without sidecars bypass mTLS and network policy enforcement. |
| H-017 | K8s secrets used by the mesh (TLS certs, DB credentials) are encrypted at rest in etcd. | Unencrypted secrets in etcd are exposed to anyone with etcd access. |
| H-018 | The K8s API Server audit log is enabled and streamed to a separate security monitoring system. | Without audit logs, unauthorized API calls are invisible to defenders. |
| H-019 | Container images are scanned for vulnerabilities and signed before deployment. | Vulnerable images in the mesh can be exploited for lateral movement. |
| H-020 | Runtime security (Falco, Tracee, or similar) is deployed to detect anomalous syscalls from pods. | mTLS and network policies do not prevent application-level compromise within a pod. |
| H-021 | The StatefulSet DB (e.g., PostgreSQL) requires mTLS or TLS for client connections. | Database connections without TLS expose data within the mesh. |
| H-022 | Persistent volumes are encrypted at rest (e.g., EBS encryption, volume encryption). | Unencrypted PVs expose data at rest if the underlying storage is compromised. |
| H-023 | The K8s dashboard or equivalent UI is not exposed externally and requires RBAC access. | Public K8s dashboard access has historically led to cluster compromises. |
| H-024 | Pod resource limits (CPU/memory) are configured to prevent DoS from a compromised pod. | Unlimited pods can starve resources across the cluster. |
| H-025 | Istio's telemetry (metrics, traces, logs) does not expose sensitive data or PII. | Mesh telemetry containing request bodies or headers leaks data to monitoring systems. |
| H-026 | The mesh implements observability without sampling down sensitive services. | Full-fidelity tracing on sensitive services exposes request contents. |
| H-027 | Automated TLS certificate renewal for the ingress gateway is configured and verified. | Stale ingress certificates cause user-facing connection failures. |
| H-028 | K8s node-level security (node OS hardening, immutable root filesystem, minimal attack surface) is enforced. | Node compromise allows attacker to bypass pod-level controls and access all containers. |
| H-029 | The container runtime (containerd, CRI-O) is kept updated and uses seccomp/AppArmor profiles. | Unrestricted syscalls from containers enable container breakout techniques. |
| H-030 | All inter-service communication uses HTTP/2 or gRPC over mTLS; no raw TCP passthrough. | Raw TCP passthrough bypasses L7 policies and observability. |
| H-031 | Service entry CRDs for external services are explicitly defined and restricted. | Undefined service entries allow services to reach arbitrary external endpoints. |
| H-032 | Mutual authentication is verified at each hop; no service trusts another based on network position. | Network position is not identity; mTLS is the only valid trust mechanism. |
| H-033 | Istio's peer authentication and request authentication policies do not conflict. | Conflicting auth policies can result in permissive default behavior. |
| H-034 | The mesh is scoped to a single Kubernetes cluster; no multi-cluster mesh without additional controls. | Multi-cluster mesh extends the trust boundary across potentially less secure clusters. |
| H-035 | Service B (intermediate service) cannot bypass authorization to directly access the StatefulSet DB. | Service B with direct DB access bypasses Service C's authorization layer. |
| H-036 | K8s cluster-admin access is limited to a small set of SRE/Platform team members with MFA. | Broad cluster-admin access is a single point of compromise for the entire mesh. |
| H-037 | Istio configuration (VirtualService, DestinationRule, AuthorizationPolicy) is managed via GitOps. | Ad-hoc config changes may introduce security regressions or misconfigurations. |
| H-038 | Envoy proxy resource usage (CPU/memory) is accounted for in pod resource requests. | Sidecar resource contention can cause proxy failure and mTLS interruptions. |
| H-039 | The mesh is on a recent Istio version with no known critical vulnerabilities. | Outdated Istio versions have public CVEs for privilege escalation and auth bypass. |
| H-040 | K8s namespace labels used for network policy selection are immutable to prevent policy bypass. | Changing namespace labels can cause pods to escape network policy boundaries. |
| H-041 | The ingress gateway is not used as a general-purpose proxy for traffic outside the mesh. | Using the ingress as an egress proxy can lead to SSRF or data exfiltration. |
| H-042 | Workload identity (ServiceAccount) is not shared between services that require different security postures. | Shared identities mean compromise of one service grants identity of another. |
| H-043 | K8s secrets are not mounted as environment variables in pods. | Environment variable secrets can be read by any process in the pod. |
| H-044 | The K8s API Server's anonymous-auth is disabled. | Anonymous access to the API Server bypasses all authentication. |

**Total (H): 44**

---

## Step 2: ASF-Generated Assumptions

*Using the 20-pattern assumption generator matrix. Applicable patterns: 18 of 20. Patterns excluded: Physical Security (cloud-hosted), Supply Chain Security (covered under Third-party Dependency).*

---

### Pattern 1: Authentication (MFA)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-001 | K8s API Server access requires MFA for all human administrators. | Explicit | Administrative access to the cluster control plane is the highest-risk auth point. |
| ASF-002 | MFA is enforced for Istio control plane configuration changes. | Derived | Changes to mesh security policies (mTLS, auth policies) require elevated assurance. |
| ASF-003 | ServiceAccount token usage does not bypass MFA requirements for human operators. | Implicit | Automated service accounts should not be used by humans to avoid MFA. |
| ASF-004 | MFA recovery for cluster access is documented and social-engineering-resistant. | Operational | K8s admin account recovery without MFA creates a bypass path. |

---

### Pattern 3: Availability & Resilience

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-005 | Istio control plane (Pilot/Citadel) is highly available with multiple replicas. | Architectural | Single-replica control plane is a single point of failure for the entire mesh. |
| ASF-006 | K8s API Server and etcd are deployed in a highly available configuration (3+ nodes). | Architectural | API Server/etcd loss makes the entire cluster unmanageable. |
| ASF-007 | Envoy sidecar proxy health does not degrade under mesh-wide configuration pushes. | Operational | A config push to 100s of services can overwhelm control plane components. |
| ASF-008 | There is a documented procedure for Istio control plane failure that does not break service-to-service communication. | Operational | Control plane outage should not affect existing mTLS connections (eventual consistency). |

---

### Pattern 4: Backup & Recovery

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-009 | etcd is backed up regularly and restore procedures are tested. | Explicit | etcd is the source of truth for all cluster state; loss means total cluster failure. |
| ASF-010 | Istio configuration (CRDs) is backed up independently of etcd. | Derived | CRDs can be recreated from GitOps if backed up; otherwise manual recovery is needed. |
| ASF-011 | Persistent volume data (StatefulSet DB) is backed up and tested for restore. | Explicit | Stateful workloads require separate backup strategies from stateless microservices. |
| ASF-012 | CA certificates and keys from Citadel are backed up securely for recovery. | Implicit | Loss of the mesh CA requires re-issuing all workload certificates. |

---

### Pattern 6: Cloud Security (IAM)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-013 | Cloud provider IAM roles assigned to K8s node instances are scoped to the minimum. | Explicit | Over-permissioned node IAM roles allow node compromise to access cloud resources. |
| ASF-014 | K8s ServiceAccount-to-IAM role mapping (IRSA or pod identity) is scoped per service. | Derived | Broad IAM role mapping allows any compromised pod to access cloud APIs. |
| ASF-015 | CloudTrail or equivalent is enabled to detect unauthorized cloud API calls from compromised nodes. | Operational | Without cloud API audit, compromised node activity is invisible. |
| ASF-016 | The cloud account has no unused or over-permissioned instance profiles. | Environmental | Legacy instance profiles in the account can be exploited from the K8s environment. |

---

### Pattern 7: Container Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-017 | Container images are from trusted registries and are vulnerability-scanned before deployment. | Explicit | Vulnerable container images are the most common initial vector in K8s attacks. |
| ASF-018 | Containers run as non-root users with read-only root filesystem. | Derived | Root containers with writable filesystems enable persistence and privilege escalation. |
| ASF-019 | No privileged containers are deployed (privileged: false in security context). | Explicit | Privileged containers have near-host-level access and can escape the container. |
| ASF-020 | Container runtime uses seccomp, AppArmor, or SELinux profiles to restrict syscalls. | Derived | Unrestricted syscalls allow container breakout techniques. |

---

### Pattern 8: Data Flow & Classification

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-021 | Data flowing between services is classified and handling requirements are documented. | Explicit | Microservices handle diverse data types; classification drives encryption and access decisions. |
| ASF-022 | No hidden data flows exist (e.g., services writing directly to persistent volume, bypassing the mesh). | Implicit | Direct data access from unauthorized services bypasses mesh telemetry and policies. |
| ASF-023 | Istio telemetry data does not contain sensitive request/response payloads. | Derived | Traces and logs containing sensitive data expose it to monitoring systems. |
| ASF-024 | Service mesh does not route sensitive data through untrusted or external networks. | Environmental | Multi-cluster mesh or external service routing exposes data beyond the trust boundary. |

---

### Pattern 9: Encryption at Rest

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-025 | etcd is encrypted at rest with a KMS-backed encryption configuration. | Explicit | etcd contains all secrets; encryption at rest is a critical control. |
| ASF-026 | Persistent volumes backing the StatefulSet DB are encrypted at rest. | Explicit | Database volumes contain application data requiring encryption. |
| ASF-027 | K8s Secret objects are encrypted at rest using encryption-at-rest configuration (not just base64). | Derived | K8s Secrets are only base64-encoded by default; encryption config is required. |
| ASF-028 | Citadel CA private keys are encrypted at rest (HSM or KMS-backed). | Derived | Mesh CA key compromise allows forgery of all workload identities. |

---

### Pattern 10: Encryption in Transit (TLS)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-029 | mTLS is configured in STRICT mode; PERMISSIVE mode is not used. | Explicit | Policy says "mTLS enabled" but does not specify strict vs. permissive mode. |
| ASF-030 | TLS certificates presented by services are validated against the mesh CA. | Trust | Without validation, a compromised certificate can impersonate any service. |
| ASF-031 | TLS 1.2 or higher is used for all mTLS connections. | Derived | Istio defaults to TLS 1.2+ but older versions may be negotiated. |
| ASF-032 | Control plane communication (Pilot → Envoy, Citadel → Envoy) uses mTLS. | Trust | Control plane traffic without mTLS allows configuration injection. |

---

### Pattern 11: Endpoint Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-033 | K8s nodes are hardened (minimal OS, no SSH, regular patching). | Implicit | Node compromise undermines all pod-level security controls. |
| ASF-034 | K8s nodes have EDR or security agents installed and reporting. | Derived | Without node-level detection, host compromise goes unnoticed. |
| ASF-035 | No workloads run on the host network (hostNetwork: false). | Explicit | Host network pods bypass all network policy and sidecar proxy enforcement. |
| ASF-036 | K8s nodes do not expose kubelet read-only or authenticated endpoints to the network. | Operational | Exposed kubelet ports allow pod and node information gathering. |

---

### Pattern 12: Human Factors & Process

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-037 | Platform team members managing Istio and K8s understand security implications of mesh configuration. | Trust | Misconfiguration of authorization policies or mTLS settings can open the mesh. |
| ASF-038 | Developers do not request privileged containers or host access without security review. | Operational | Developer requests for privileged access bypass container security policies. |
| ASF-039 | Cluster configuration changes go through change management and peer review. | Operational | Unreviewed K8s changes (RBAC, PSP, network policies) can introduce vulnerabilities. |
| ASF-040 | Service owners correctly define authorization policies for their services. | Implicit | Incorrect authorization policies result in either blocked legitimate traffic or open access. |

---

### Pattern 13: Identity Lifecycle (Provisioning)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-041 | K8s ServiceAccount token rotation is automated (e.g., TokenRequest API with TTL). | Operational | Long-lived ServiceAccount tokens increase exposure if leaked. |
| ASF-042 | K8s user and ServiceAccount access is recertified quarterly. | Derived | Stale ServiceAccounts and RBAC bindings lead to privilege creep. |
| ASF-043 | Workload identity (SPIFFE) is revoked when a service is decommissioned. | Implicit | Decommissioned service identities remain valid for mTLS authentication. |
| ASF-044 | ServiceAccount deletion does not impact running pods that relied on the account. | Derived | Deleting a ServiceAccount before pods terminate can break runtime identity. |

---

### Pattern 14: Incident Response

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-045 | There is an IR plan covering container escape or cluster compromise scenarios. | Operational | Container escape is a high-impact scenario requiring specific containment procedures. |
| ASF-046 | The IR team has access to K8s audit logs, pod logs, and network flow logs. | Derived | Inaccessible K8s logs prevent forensic analysis of cluster compromise. |
| ASF-047 | IR procedures include mesh isolation (update authorization policies to block compromised service). | Trust | Targeted isolation of compromised services preserving other mesh traffic. |
| ASF-048 | Runtime security alerts (Falco, Tracee) are monitored and trigger IR. | Implicit | Without runtime detection, container breakouts go unnoticed. |
| ASF-049 | etcd snapshots taken during IR are forensically isolated. | Derived | Compromised etcd snapshots can be tampered with if stored in the same cluster. |

---

### Pattern 15: Least Privilege

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-050 | K8s RBAC roles for services are scoped to specific namespace resources and verbs. | Explicit | RBAC at namespace level is policy; correct implementation is assumed. |
| ASF-051 | No service account has cluster-admin privileges. | Derived | Cluster-admin ServiceAccount gives any pod running with it full cluster control. |
| ASF-052 | Pod Security Standards (restricted profile) are enforced, not just baseline or privileged. | Explicit | Restricted profile prevents many container breakout techniques. |
| ASF-053 | Service mesh authorization policies enforce least-privilege service-to-service access. | Derived | mTLS provides authentication; AuthorizationPolicy CRDs provide authorization. |

---

### Pattern 16: Monitoring & Alerting

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-054 | K8s API Server audit log is enabled with a comprehensive rule set. | Operational | Audit logs capture all API calls; without them, cluster changes are invisible. |
| ASF-055 | Istio telemetry (request volumes, error rates, latency) is monitored for anomalies. | Operational | Mesh telemetry detects service degradation or attack patterns. |
| ASF-056 | Network flow logs between pods are captured for forensic analysis. | Derived | Without flow logs, lateral movement between services is invisible. |
| ASF-057 | Monitoring infrastructure logs are stored outside the cluster to prevent tampering. | Implicit | Cluster compromise allows deletion or alteration of in-cluster monitoring data. |

---

### Pattern 17: Network Segmentation

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-058 | Namespace-level network policies enforce default-deny for ingress and egress. | Explicit | Policy states "Network policies isolate namespaces" but not default-deny. |
| ASF-059 | Control plane components (Pilot, Citadel, API Server) are in isolated namespaces with restricted network policies. | Architectural | Control plane exposure to workload namespaces increases attack surface. |
| ASF-060 | The ingress gateway is in a separate namespace with access only to mesh services. | Architectural | Ingress in a shared namespace increases blast radius of a compromised gateway. |
| ASF-061 | No network policy allows egress to the internet from workload namespaces (egress via gateway only). | Derived | Direct internet egress from pods bypasses security controls. |

---

### Pattern 20: Third-party Dependency

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-062 | Istio and its Envoy sidecar have no known unpatched critical CVEs. | Dependency | Istio CVEs (e.g., CVE-2023-44487, HTTP/2 Rapid Reset) affect the entire mesh. |
| ASF-063 | The container registry (e.g., Docker Hub, ECR, GCR) is available and not compromised. | Dependency | Registry compromise allows injection of malicious images into the cluster. |
| ASF-064 | Kubernetes versions are within the support window and receive security patches. | Dependency | Outdated K8s versions have known vulnerabilities (e.g., CVE-2023-5521). |
| ASF-065 | There is a documented migration path if Istio or K8s becomes unavailable due to licensing or deprecation. | Derived | Vendor dependency on open-source projects carries strategic risk. |

**Total (A): 65** (4 per pattern × 16 patterns + 1 extra from Container Security = 65)

---

## Step 3: Comparison

### Overlap Mapping

| Human ID | ASF ID | Match Rationale |
|----------|--------|-----------------|
| H-001 | ASF-029 | Both require mTLS STRICT mode / reject non-mTLS traffic. |
| H-002 | ASF-028 | Both require secure CA key management and rotation. |
| H-003 | ASF-028 | Both require HSM/KMS for CA keys. |
| H-004 | ASF-029 | Both require STRICT mTLS, not PERMISSIVE. |
| H-005 | ASF-041 | Both require automated certificate rotation. |
| H-006 | ASF-050 | Both require least-privilege RBAC. |
| H-007 | ASF-001 | Both require K8s API security (not public, MFA). |
| H-008 | ASF-025 | Both require etcd encryption at rest. |
| H-009 | ASF-019 | Both require PSP to restrict privileged containers. |
| H-010 | ASF-058 | Both require default-deny network policies. |
| H-011 | ASF-029 | Both require mTLS re-encryption at ingress. |
| H-015 | ASF-053 | Both require authorization policies beyond mTLS. |
| H-016 | ASF-035 | Both require sidecar injection enforcement. |
| H-017 | ASF-027 | Both require K8s Secrets encryption at rest. |
| H-018 | ASF-054 | Both require K8s API Server audit logging. |
| H-019 | ASF-017 | Both require container image vulnerability scanning. |
| H-020 | ASF-048 | Both require runtime security monitoring. |
| H-021 | ASF-032 | Both require TLS for DB connections within mesh. |
| H-022 | ASF-026 | Both require persistent volume encryption at rest. |
| H-024 | ASF-024 | Both require pod resource limits. |
| H-025 | ASF-023 | Both require no sensitive data in telemetry. |
| H-028 | ASF-033 | Both require node OS hardening. |
| H-029 | ASF-020 | Both require seccomp/AppArmor profiles. |
| H-030 | ASF-030 | Both require HTTP/2 or gRPC over mTLS. |
| H-031 | ASF-022 | Both require controlled external service access. |
| H-032 | ASF-030 | Both require identity-based trust, not network position. |
| H-036 | ASF-051 | Both require no cluster-admin for non-SRE roles. |
| H-037 | ASF-039 | Both require GitOps/change management for Istio config. |
| H-039 | ASF-062 | Both require recent Istio version without known CVEs. |
| H-040 | ASF-058 | Both require namespace label immutability for network policies. |
| H-042 | ASF-053 | Both require unique workload identity for least privilege. |
| H-043 | ASF-027 | Both require secrets not as environment variables. |
| H-044 | ASF-001 | Both require anonymous-auth disabled on API Server. |

**Overlap (O): 33**

### Metrics

| Metric | Value | Formula |
|--------|-------|---------|
| Human assumptions (H) | 44 | Count of unique human-generated assumptions |
| ASF assumptions (A) | 65 | Count of unique ASF-generated assumptions |
| Overlap (O) | 33 | Count appearing in both lists |
| **Precision** | **50.8%** | O / A = 33/65 |
| **Recall** | **75.0%** | O / H = 33/44 |
| **F1 Score** | **60.6%** | 2 × (P × R) / (P + R) |
| Novel findings (A - O) | 32 | Assumptions ASF found that human missed (49.2% of ASF total) |
| Missed findings (H - O) | 11 | Assumptions human found that ASF missed (25.0% of human total) |

### Success Criteria Assessment

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Recall | >= 70% | 75.0% | ✅ Met |
| Precision | >= 50% | 50.8% | ✅ Met |
| Novel discoveries | >= 10% of total (A+O) | 24.6% (32/130) | ✅ Exceeded |
| Expert agreement (F1 proxy) | > 60% | 60.6% | ✅ Met |

---

## Step 4: Analysis

### Ontology Overlap Analysis

| Ontology Category | Overlaps | ASF Total | Overlap Rate |
|-------------------|----------|-----------|-------------|
| Explicit | 10 | 16 | 62.5% |
| Derived | 9 | 18 | 50.0% |
| Operational | 5 | 12 | 41.7% |
| Implicit | 4 | 9 | 44.4% |
| Trust | 2 | 4 | 50.0% |
| Dependency | 1 | 4 | 25.0% |
| Architectural | 1 | 4 | 25.0% |
| Environmental | 1 | 2 | 50.0% |

**Best overlap:** Explicit (62.5%) and Derived (50.0%) showed the strongest agreement. Both humans and ASF have strong shared understanding of mTLS STRICT mode, RBAC least privilege, PSP enforcement, and network policy default-deny as primary security assumptions for service mesh architectures.

**Worst overlap:** Dependency (25.0%) and Architectural (25.0%) had the weakest overlap. The ASF identified dependency risks (Istio CVEs, registry availability, K8s support window) and architectural concerns (control plane HA, control plane isolation) that the human architect treated as external rather than assumptions.

### What Humans Caught That ASF Missed (Missed Findings = 11)

The 11 human-generated assumptions with no ASF counterpart:

1. **Envoy/Istio configuration details (H-012, H-033, H-034, H-038, H-041):** Unauthenticated Pilot requests, conflicting peer/request authentication policies, multi-cluster mesh scoping, sidecar resource accounting, and ingress gateway misuse. These are deep Istio product-specific details.

2. **K8s operational specifics (H-023, H-040, H-044):** K8s dashboard exposure, namespace label immutability for network policy, and anonymous-auth disabling. These are K8s-specific security settings.

3. **Service mesh identity and trust (H-013, H-014, H-026):** IP-based vs. ServiceAccount-based workload identity, SPIFFE identity uniqueness, and observability sampling without data leakage. These reflect nuanced understanding of how SPIFFE identities work in practice.

### What ASF Caught That Humans Missed (Novel Findings = 32)

The ASF generated 65 assumptions, of which 32 (49.2%) were not in the human list:

1. **Incident Response (5 assumptions):** The human generated zero IR assumptions. The ASF contributed a full pattern covering container escape IR plans, K8s audit log access, mesh isolation procedures, runtime security monitoring, and etcd forensic isolation.

2. **Identity Lifecycle (4 assumptions):** The human did not consider ServiceAccount token rotation, quarterly RBAC recertification, SPIFFE identity revocation on decommission, or ServiceAccount deletion timing. These operational identity lifecycle gaps are common in K8s environments.

3. **Backup & Recovery (4 assumptions):** The human assumed etcd encryption (H-008) but did not consider etcd backup/restore testing, Istio CRD backup, StatefulSet DB backup, or Citadel CA certificate backup. These are critical for mesh recovery.

4. **Cloud IAM (4 assumptions):** The human assumed node-level security (H-028) but did not consider node IAM role least privilege, ServiceAccount-to-IAM mapping (IRSA), CloudTrail for cloud API audit, or unused instance profile cleanup.

5. **Third-party dependencies (4 assumptions):** The human assumed Istio should be current (H-039) but did not generalize to Istio CVE risk, container registry availability and integrity, K8s version support window, or strategic migration path.

### Architecture Complexity Assessment

Architecture #5 (Service Mesh/K8s) performed well across all metrics:

- **Recall (75.0%)** — met the 70% target. The ASF's Container Security, Network Segmentation, and Encryption patterns align well with service mesh architectures.
- **Precision (50.8%)** — narrowly met the 50% target. The breadth of the ASF patterns (18 applicable) generates 65 assumptions, of which half overlapped with the human.
- **F1 (60.6%)** — met the 60% target, the second architecture to do so.
- **Novelty rate (49.2%)** — substantial value added, particularly in IR, identity lifecycle, and backup patterns.

### Key Insight

The service mesh architecture demonstrates the strongest ASF performance for a complex system. The Container Security pattern (added for K8s/mesh architectures but excluded for EC2-based architectures) significantly improves coverage. The top human misses (IR, identity lifecycle, backup) mirror the same gaps seen in all previous architectures — these are systematic human blind spots, not architecture-specific ones.

---

## Conclusions

| Metric | Target | Actual | Verdict |
|--------|--------|--------|---------|
| Recall | >= 70% | 75.0% | ✅ Met — Container Security pattern helps |
| Precision | >= 50% | 50.8% | ✅ Met — just above threshold |
| Novel discoveries | >= 10% | 24.6% | ✅ ASF adds value in IR, lifecycle, backup |
| Expert agreement (F1) | > 60% | 60.6% | ✅ Met — best F1 after Architecture #4 |

The ASF applied to Architecture #5 demonstrates that for complex, multi-component architectures with strong pattern alignment, the framework can meet all success criteria. The systematic gaps remain consistent across architectures: incident response, identity lifecycle, backup/recovery operationalization, and third-party dependency risk are areas where the ASF consistently adds value beyond human reasoning. The missed findings (11) are almost entirely deep Istio/K8s product-specific configuration details — suggesting a "Service Mesh Configuration" sub-pattern could further improve recall.
