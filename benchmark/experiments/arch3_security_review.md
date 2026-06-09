# Architecture 3 — Security Review Deliverable
## K8s/Istio Service Mesh

**Reviewer**: Security Architecture Team
**Date**: June 9, 2026
**Scope**: Architecture 3 — Kubernetes cluster with Istio service mesh, Citadel CA, sidecar proxies, Ingress Gateway, StatefulSet database, namespace RBAC, and mTLS-enforced service-to-service communication

---

## 1. Consensus Matrix

| # | Assumption | GPT | Gemini | Gemma | Keep? |
|---|-----------|:---:|:------:|:-----:|:-----:|
| A1 | Worker nodes are trusted and secure | ✓ | | | ✓ |
| A2 | Container runtime is secure — no escape vulnerabilities | ✓ | ✓ | ✓ | ✓ |
| A3 | Linux kernel namespace isolation is robust | | | ✓ | ✓ |
| A4 | Time synchronization is accurate across the cluster | ✓ | | | ✓ |
| A5 | No alternate networking path bypasses mesh policies | ✓ | ✓ | | ✓ |
| A6 | K8s control plane is secure | ✓ | | | ✓ |
| A7 | K8s control plane is isolated from the data plane | | ✓ | | ✓ |
| A8 | etcd is protected and encrypted | ✓ | | ✓ | ✓ |
| A9 | etcd access is restricted to authorized components only | ✓ | | ✓ | ✓ |
| A10 | etcd data integrity is maintained | ✓ | | | ✓ |
| A11 | K8s API authentication is secure | ✓ | | | ✓ |
| A12 | K8s API authorization is correctly configured | ✓ | | | ✓ |
| A13 | Cluster-admin privileges are tightly controlled | ✓ | | | ✓ |
| A14 | No hidden cluster-wide permissions exist | ✓ | | | ✓ |
| A15 | Admission controllers prevent unauthorized workloads | ✓ | | | ✓ |
| A16 | Operators and administrators are trustworthy | ✓ | | | ✓ |
| A17 | Citadel CA is trustworthy | ✓ | ✓ | | ✓ |
| A18 | CA private keys are protected from compromise | ✓ | ✓ | | ✓ |
| A19 | Certificate issuance is properly controlled | ✓ | | | ✓ |
| A20 | Certificate revocation works effectively | ✓ | | | ✓ |
| A21 | Certificate rotation functions correctly | ✓ | | | ✓ |
| A22 | Istio Pilot (control plane) is trustworthy | ✓ | | | ✓ |
| A23 | Pilot configuration is accurate and not tampered | ✓ | | | ✓ |
| A24 | Sidecar proxies cannot be bypassed by workloads | ✓ | | | ✓ |
| A25 | All service traffic actually traverses the mesh | ✓ | | | ✓ |
| A26 | Service identities are unique and non-confusable | ✓ | | | ✓ |
| A27 | Service certificates correctly represent workloads | ✓ | | | ✓ |
| A28 | mTLS is enforced rather than operating in permissive mode | ✓ | | | ✓ |
| A29 | Services validate peer certificates | ✓ | | | ✓ |
| A30 | Applications do not log or leak mTLS session tokens or decrypted payloads | | | ✓ | ✓ |
| A31 | Pod Security Policies (or equivalent admission controls) are enforced | ✓ | | | ✓ |
| A32 | Service accounts are protected from theft | ✓ | | | ✓ |
| A33 | Service account tokens cannot be stolen from pods | ✓ | | | ✓ |
| A34 | Namespace RBAC is correctly implemented | ✓ | | | ✓ |
| A35 | Network policies are correctly defined for namespace isolation | ✓ | | | ✓ |
| A36 | Network policies are actually enforced by the CNI plugin | ✓ | | | ✓ |
| A37 | DNS inside the cluster is trustworthy | ✓ | | | ✓ |
| A38 | Secrets management is secure | ✓ | | | ✓ |
| A39 | Secrets are encrypted when stored (at rest) | ✓ | | | ✓ |
| A40 | Persistent Volumes are protected from unauthorized access | ✓ | ✓ | | ✓ |
| A41 | Persistent Volume snapshots are secure | ✓ | | | ✓ |
| A42 | PV data is securely encrypted and isolated at the storage layer | | ✓ | | ✓ |
| A43 | Ingress Gateway is hardened against external attack | ✓ | | ✓ | ✓ |
| A44 | Ingress routing rules are secure and not exploitable | ✓ | | | ✓ |
| A45 | Ingress Gateway sanitizes external traffic before routing to internal services | | | ✓ | ✓ |
| A46 | StatefulSet database permissions follow least privilege | ✓ | | | ✓ |
| A47 | Service A trusts only authorized callers (mTLS identity ≠ authorization) | ✓ | | | ✓ |
| A48 | Service B enforces authorization decisions independently | ✓ | | | ✓ |
| A49 | Service C protects sensitive operations from unauthorized access | ✓ | | | ✓ |
| A50 | No undocumented trust relationships exist in the cluster | ✓ | | | ✓ |
| A51 | Monitoring and telemetry systems are trustworthy | ✓ | | | ✓ |
| A52 | Logs cannot be modified or deleted by attackers | ✓ | | | ✓ |
| A53 | Container images are trusted and free of malware | ✓ | | | ✓ |
| A54 | Image registries are secure and not compromised | ✓ | | | ✓ |

**Total Assumptions: 54** (50 GPT + 2 Gemini-exclusive + 2 Gemma-exclusive)

---

## 2. Deduplicated Assumption List

### 2.1 Cluster Infrastructure
1. Worker nodes are trusted and secure (A1)
2. Container runtime is secure — no escape vulnerabilities (A2)
3. Linux kernel namespace isolation is robust (A3)
4. Time synchronization is accurate across the cluster (A4)
5. No alternate networking path bypasses mesh policies (A5)

### 2.2 Control Plane
6. K8s control plane is secure (A6)
7. K8s control plane is isolated from the data plane (A7)
8. etcd is protected and encrypted (A8)
9. etcd access is restricted to authorized components only (A9)
10. etcd data integrity is maintained (A10)
11. K8s API authentication is secure (A11)
12. K8s API authorization is correctly configured (A12)
13. Cluster-admin privileges are tightly controlled (A13)
14. No hidden cluster-wide permissions exist (A14)
15. Admission controllers prevent unauthorized workloads (A15)
16. Operators and administrators are trustworthy (A16)

### 2.3 Service Mesh / mTLS
17. Citadel CA is trustworthy (A17)
18. CA private keys are protected from compromise (A18)
19. Certificate issuance is properly controlled (A19)
20. Certificate revocation works effectively (A20)
21. Certificate rotation functions correctly (A21)
22. Istio Pilot (control plane) is trustworthy (A22)
23. Pilot configuration is accurate and not tampered (A23)
24. Sidecar proxies cannot be bypassed by workloads (A24)
25. All service traffic actually traverses the mesh (A25)
26. Service identities are unique and non-confusable (A26)
27. Service certificates correctly represent workloads (A27)
28. mTLS is enforced rather than operating in permissive mode (A28)
29. Services validate peer certificates (A29)
30. Applications do not log or leak mTLS session tokens or decrypted payloads (A30)

### 2.4 Pod Security
31. Pod Security Policies (or equivalent admission controls) are enforced (A31)
32. Service accounts are protected from theft (A32)
33. Service account tokens cannot be stolen from pods (A33)
34. Namespace RBAC is correctly implemented (A34)

### 2.5 Network Policies
35. Network policies are correctly defined for namespace isolation (A35)
36. Network policies are actually enforced by the CNI plugin (A36)
37. DNS inside the cluster is trustworthy (A37)

### 2.6 Secrets / Encryption
38. Secrets management is secure (A38)
39. Secrets are encrypted when stored (at rest) (A39)
40. Persistent Volumes are protected from unauthorized access (A40)
41. Persistent Volume snapshots are secure (A41)
42. PV data is securely encrypted and isolated at the storage layer (A42)

### 2.7 Workload Security
43. Ingress Gateway is hardened against external attack (A43)
44. Ingress routing rules are secure and not exploitable (A44)
45. Ingress Gateway sanitizes external traffic before routing to internal services (A45)
46. StatefulSet database permissions follow least privilege (A46)
47. Service A trusts only authorized callers (mTLS identity ≠ authorization) (A47)
48. Service B enforces authorization decisions independently (A48)
49. Service C protects sensitive operations from unauthorized access (A49)
50. No undocumented trust relationships exist in the cluster (A50)

### 2.8 Monitoring
51. Monitoring and telemetry systems are trustworthy (A51)
52. Logs cannot be modified or deleted by attackers (A52)

### 2.9 Supply Chain
53. Container images are trusted and free of malware (A53)
54. Image registries are secure and not compromised (A54)

### 2.10 Operations
(Assumptions A13–A16 under Control Plane cover operations; no additional operations-specific assumptions were identified)

---

## 3. Risk Scores

| # | Assumption | Likelihood | Impact | Risk |
|---|-----------|:----------:|:------:|:----:|
| A1 | Worker nodes trusted | M | C | C |
| A2 | Container runtime secure | M | C | C |
| A3 | Kernel namespace isolation robust | M | C | C |
| A4 | Time sync accurate | L | H | H |
| A5 | No alternate network path | M | C | C |
| A6 | Control plane secure | L | C | H |
| A7 | Control plane isolated from data plane | M | C | C |
| A8 | etcd protected and encrypted | L | C | H |
| A9 | etcd access restricted | L | C | H |
| A10 | etcd data integrity | L | C | H |
| A11 | K8s API auth secure | L | C | H |
| A12 | K8s API authz correct | M | C | C |
| A13 | Cluster-admin tightly controlled | M | C | C |
| A14 | No hidden permissions | H | H | H |
| A15 | Admission controllers effective | M | C | C |
| A16 | Operators/administrators trustworthy | L | C | H |
| A17 | Citadel CA trustworthy | L | C | H |
| A18 | CA private keys protected | L | C | H |
| A19 | Certificate issuance controlled | M | C | C |
| A20 | Certificate revocation effective | M | C | C |
| A21 | Certificate rotation correct | M | H | H |
| A22 | Istio Pilot trustworthy | L | C | H |
| A23 | Pilot config accurate | M | C | C |
| A24 | Sidecar proxies not bypassable | M | C | C |
| A25 | All traffic traverses mesh | M | C | C |
| A26 | Service identities unique | L | H | H |
| A27 | Certs represent workloads | L | C | H |
| A28 | mTLS enforced (not permissive) | M | C | C |
| A29 | Services validate peer certs | M | C | C |
| A30 | mTLS tokens not leaked in logs | H | C | C |
| A31 | Pod Security Policies enforced | M | C | C |
| A32 | Service accounts protected | M | C | C |
| A33 | SA tokens not stolen | M | C | C |
| A34 | Namespace RBAC correct | M | H | H |
| A35 | Network policies defined correctly | M | C | C |
| A36 | Network policies enforced | M | C | C |
| A37 | DNS trustworthy | L | H | H |
| A38 | Secrets management secure | M | C | C |
| A39 | Secrets encrypted at rest | L | C | H |
| A40 | Persistent Volumes protected | M | C | C |
| A41 | PV snapshots secure | M | H | H |
| A42 | PV encrypted and isolated at storage layer | M | C | C |
| A43 | Ingress Gateway hardened | H | C | C |
| A44 | Ingress routing rules secure | M | H | H |
| A45 | Ingress sanitizes external traffic | H | C | C |
| A46 | StatefulSet DB least privilege | M | H | H |
| A47 | Service A trusts authorized callers only | M | C | C |
| A48 | Service B enforces authz independently | M | C | C |
| A49 | Service C protects sensitive ops | M | C | C |
| A50 | No undocumented trust relationships | H | C | C |
| A51 | Monitoring systems trustworthy | M | C | C |
| A52 | Logs not modifiable by attackers | M | C | C |
| A53 | Container images trusted | M | C | C |
| A54 | Image registries secure | M | C | C |

---

## 4. STRIDE Mapping

### Spoofing
- A11: K8s API authentication
- A17–A18: Citadel CA trust and key protection
- A19: Certificate issuance control
- A26: Service identity uniqueness
- A27: Certificate-to-workload binding
- A29: Peer certificate validation
- A32–A33: Service account and token protection

### Tampering
- A6–A7: Control plane integrity and isolation
- A8–A10: etcd protection, access, and integrity
- A12: Authorization correctness
- A14: Hidden permission detection
- A15: Admission controller bypass
- A22–A23: Pilot trust and configuration integrity
- A24: Sidecar proxy bypass
- A28: mTLS mode enforcement (permissive → plaintext)
- A35–A36: Network policy tampering
- A39: Secrets at-rest encryption
- A40–A42: PV and snapshot protection
- A44: Ingress routing rule tampering
- A52: Log immutability

### Repudiation
- A4: NTP time synchronization
- A51: Monitoring system trust
- A52: Log integrity

### Information Disclosure
- A1: Worker node compromise → data exposure
- A2–A3: Container escape → cross-tenant access
- A8: etcd unencrypted access
- A30: mTLS token/decrypted payload leakage in logs
- A38–A39: Secrets exposure
- A40–A42: PV data exposure
- A45: Ingress traffic lacking sanitization
- A53–A54: Malicious image deployment

### Denial of Service
- A4: NTP failure → cert validation failures
- A6: Control plane DoS
- A22: Pilot DoS → mesh routing collapse
- A24: Sidecar proxy resource exhaustion
- A35–A36: Network policy misconfiguration → connectivity loss

### Elevation of Privilege
- A12–A13: RBAC misconfiguration / cluster-admin abuse
- A14: Hidden cluster-wide permissions
- A16: Insider threat from operators
- A31: Pod Security Policy bypass → host access
- A32–A33: Service account token theft
- A34: Namespace RBAC boundary crossing
- A46: StatefulSet DB excessive permissions
- A47–A49: Missing application-layer authorization
- A50: Undocumented trust relationships

---

## 5. Top 10 Critical Assumptions (Ranked)

### 1. CA Root Compromise — Citadel CA and Private Key Protection (A17 / A18)
**Rationale**: The entire mesh trust model collapses if the Citadel CA or its private keys are compromised. An attacker with CA access can mint valid certificates for any service identity, impersonate any workload, decrypt all mTLS traffic, and move laterally across the entire cluster without detection. This is the single root of trust for the Istio security model. No other control can compensate for a compromised CA.

### 2. Decrypted Data Leakage via Application Logs (A30)
**Rationale**: Internal services within the mesh trust mTLS for transport encryption, but the decrypted payload is visible to the application process. If applications log decrypted request bodies, session tokens, or sensitive data in debug output, those logs become a primary vector for credential and data leakage — accessible via log aggregators, monitoring systems, or compromised pods with log read access. This assumption is unique to Gemma and addresses a gap overlooked by GPT and Gemini.

### 3. Ingress Gateway Sanitization of External Traffic (A45)
**Rationale**: Internal services that rely on mTLS for peer authentication may skip input validation on the assumption that all incoming traffic is from trusted mesh peers. The Ingress Gateway terminates external TLS and re-encrypts with mTLS for internal routing, but if it does not also sanitize and validate the payload, attackers can inject malicious requests that internal services trust implicitly. This is the primary external → internal attack path.

### 4. Sidecar Proxy Bypass / Traffic Not Traversing the Mesh (A24 / A25)
**Rationale**: Istio's security model depends on all workload traffic being intercepted by the Envoy sidecar proxy. If a workload can communicate directly (e.g., via host networking, iptables bypass, or a missing sidecar injection), traffic may travel unencrypted outside the mesh, bypassing mTLS, authorization policies, and telemetry. Partial mesh adoption or misconfigured sidecar injection creates invisible security gaps.

### 5. Container Escape via Runtime or Kernel Vulnerability (A2 / A3)
**Rationale**: A container escape from any pod grants an attacker host-level access to the worker node. From a compromised node, the attacker can bypass all network policies, access other pods' traffic via the network namespace, steal node-level credentials, and potentially pivot to the control plane. This renders every mesh-level control moot.

### 6. etcd Compromise (A8 / A9 / A10)
**Rationale**: etcd stores all cluster state including secrets, service account tokens, RBAC definitions, and Istio configuration. Read access to etcd is equivalent to full cluster compromise. If etcd is not encrypted at rest, access is not restricted via TLS client certificates, or its integrity is not monitored, an attacker with network access to the etcd port can extract every credential in the cluster.

### 7. Secrets Management Failure (A38 / A39)
**Rationale**: Kubernetes Secrets are base64-encoded by default and not encrypted at rest unless EncryptionConfiguration is applied. Database credentials, API keys, and service account tokens stored as Secrets are accessible to any user or component with read access to the Secrets API or to etcd. Secrets management failures cascade into database compromise, service impersonation, and lateral movement.

### 8. Missing Application-Layer Authorization (A47 / A48 / A49)
**Rationale**: mTLS proves identity but does not grant authorization. If Service A trusts any mTLS-authenticated caller (rather than checking a specific allowlist), or if Services B and C do not independently enforce authorization for each operation, an attacker who compromises any mesh workload can access all mesh services. Lateral movement becomes trivial once inside the mesh.

### 9. Network Policy Bypass / Misconfiguration (A35 / A36)
**Rationale**: Namespace isolation in the architecture depends on correctly defined and enforced Kubernetes Network Policies. If policies are overly permissive, missing default-deny rules, or not enforced by the CNI plugin, workloads in different namespaces may communicate in violation of the intended security boundary. This undermines the multi-tenant isolation model.

### 10. RBAC / Cluster-Admin Privilege Creep (A13 / A14 / A34)
**Rationale**: Complex RBAC environments accumulate excessive permissions over time. If cluster-admin privileges are broadly granted, hidden bindings exist, or namespace RBAC boundaries are misconfigured, a single compromised workload or user account can escalate to cluster-wide administrative access. Regular RBAC audits are essential but frequently omitted.

---

## 6. Recommended Controls

### 1. CA Root Compromise Prevention
- Store Citadel CA private keys in a hardware security module (HSM) or cloud KMS with key release policies
- Implement key ceremony procedures for CA key generation and rotation
- Enable Istio's pluggable CA key integration for external KMS-backed keys
- Monitor and alert on all certificate issuance events; audit CA access logs
- Set short-lived certificates (24h TTL) to limit blast radius of CA compromise

### 2. mTLS Token / Payload Leakage Prevention
- Implement application-level log sanitization — strip mTLS session tokens and sensitive payload fields
- Deploy structured logging with auto-redaction for known sensitive patterns (tokens, keys, PII)
- Configure log shippers with filter rules to drop matching log entries before egress
- Conduct code reviews targeting debug log statements in production paths
- Use static analysis tools (e.g., Semgrep, CodeQL) to detect logging of sensitive data

### 3. Ingress Gateway Traffic Sanitization
- Deploy a Web Application Firewall (WAF) in front of or within the Ingress Gateway
- Validate request payload, headers, and parameters at the Gateway before routing
- Use Istio's Envoy filter chain to apply authentication, rate limiting, and input validation
- Implement strict mTLS between Ingress Gateway and internal services (ISTIO_MUTUAL)
- Do not rely on internal service input validation as a substitute for edge sanitization

### 4. Sidecar Proxy Bypass Prevention
- Enforce sidecar injection via namespace label injection and validating webhook
- Use Istio PeerAuthentication with STRICT mTLS mode globally
- Deploy network policies that block non-mesh traffic at the pod level
- Implement egress traffic controls to detect direct outbound connections bypassing the sidecar
- Run mesh health checks to verify all workloads have proxies injected and receiving traffic

### 5. Container Escape Mitigation
- Deploy a pod security admission controller (Pod Security Admission, OPA/Gatekeeper, or Kyverno)
- Drop all container capabilities, run as non-root, use read-only root filesystems
- Use seccomp, AppArmor, or SELinux profiles to restrict syscall access
- Keep worker node kernels updated with latest security patches
- Use container-optimized OS (Flatcar, Bottlerocket, GKE Container-Optimized OS)
- Deploy runtime security tools (Falco, Tracee) for container escape detection

### 6. etcd Security Hardening
- Enable etcd TLS with mutual authentication for all peer and client connections
- Encrypt etcd data at rest using etcd encryption or KMS-backed envelope encryption
- Restrict network access to etcd to only the K8s API server (firewall rule / network policy)
- Enable etcd audit logging and monitor for unauthorized access attempts
- Regularly back up etcd snapshots and test restoration procedures

### 7. Secrets Management
- Enable K8s EncryptionConfiguration with KMS provider for Secrets at rest
- Migrate from K8s Secrets to external secrets management (HashiCorp Vault, AWS Secrets Manager, GCP Secret Manager) with the CSI Secrets Store driver
- Implement least-privilege RBAC for Secrets access — deny watch/list/get on Secrets by default
- Rotate database credentials and service account tokens automatically
- Audit Secrets access logs for anomalous patterns

### 8. Application-Layer Authorization
- Implement a service-to-service authorization framework (e.g., Istio AuthorizationPolicy with principal allowlists)
- Enforce authorization checks in each service for every request — do not rely solely on mTLS identity
- Use Istio RequestAuthentication and AuthorizationPolicy for mesh-wide access control
- Adopt the Zero Trust principle: verify every request regardless of source identity
- Conduct penetration testing targeting service-to-service authorization gaps

### 9. Network Policy Guardrails
- Implement default-deny ingress and egress Network Policies for all namespaces
- Verify CNI plugin supports and enforces Network Policies (Calico, Cilium, Antrea)
- Use policy-as-code tools (OPA/Gatekeeper) to validate Network Policy compliance
- Run continuous network policy validation with kube-scan or similar tools
- Block hostNetwork pods except for explicitly approved system components

### 10. RBAC / Permission Hygiene
- Conduct quarterly RBAC audits with automated tools (kubectl who-can, rbac-manager, Polar)
- Implement just-in-time (JIT) cluster-admin access with approver workflow and time-bound grants
- Remove all cluster-admin bindings to regular user and service accounts
- Use impersonation for break-glass administrative access with full audit logging
- Deploy a privileged access management (PAM) solution for cluster administration

---

## 7. Summary Statistics

| Category | Count |
|----------|:-----:|
| **Total Assumptions** | **54** |
| **Critical Risk** | 37 |
| **High Risk** | 17 |
| **Medium Risk** | 0 |
| **Low Risk** | 0 |
| **Sources** | 3 models (GPT-4o, Gemini, Gemma) |
| **GPT-originated** | 50 |
| **Gemini-exclusive** | 2 |
| **Gemma-exclusive** | 2 |

**Critical Assumptions** (37): A1, A2, A3, A5, A7, A12, A13, A15, A17, A18, A19, A20, A22, A23, A24, A25, A27, A28, A29, A30, A31, A32, A33, A35, A36, A38, A40, A42, A43, A45, A47, A48, A49, A50, A51, A52, A53, A54

**Top Risk Drivers**: CA root compromise (A17/A18), decrypted payload leakage via application logs (A30), and Ingress Gateway traffic sanitization (A45) represent the three highest-severity assumption failures in Architecture 3. The CA root is the single point of failure for the entire mesh trust model; mTLS token leakage is a stealthy data exfiltration path that standard monitoring may miss; and the Ingress Gateway is the primary external attack surface where a sanitization gap can exploit internal services' implicit trust in mTLS-authenticated peers. These three failure modes — a compromised root of trust, a silent data leak, and an unguarded perimeter — capture the distinct threat vectors unique to a service mesh architecture.

---

*End of Security Review — Architecture 3*
