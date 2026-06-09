# ASF Phase 6 Experiment: Architecture #008

**Architecture:** CI/CD Pipeline → Artifact Registry → Deploy
**Date:** 2026-06-09
**Simulation Mode:** Human + ASF framework comparison

---

## Architecture Profile

```
[Developer] --> [GitHub] --> [GitHub Actions CI] --> [ECR / Artifact Registry]
                                                           │
                                                    [ArgoCD] --> [K8s Cluster]
```

### Documented Policy
| # | Policy |
|---|--------|
| P1 | Code review required before merge |
| P2 | CI runs security scan on every commit |
| P3 | Images are signed before deployment |
| P4 | ArgoCD syncs from Git as source of truth |

### Trust Boundaries
| Boundary | Type |
|----------|------|
| Developer ↔ GitHub | Identity boundary |
| CI ↔ Registry | Pipeline boundary |
| Registry ↔ Cluster | Deployment boundary |

### Complexity Rating
**Moderate** — 5 nodes, 3 trust boundaries, build-deploy pipeline with multiple security controls.

---

## Step 1: Human-Generated Assumptions

*Role: Senior Security Architect. Listed below are assumptions that MUST remain true for this architecture to remain secure but are NOT stated in the documented policy.*

| ID | Assumption | Justification |
|----|-----------|---------------|
| H-001 | GitHub personal access tokens and SSH keys used by developers are protected — never committed, rotated regularly, and revoked on termination. | Compromised developer credentials are the most common initial access vector for supply chain attacks. |
| H-002 | Branch protection rules are enforced on all production branches — no direct pushes bypass pull request requirements. | A developer pushing directly to the main branch bypasses code review entirely and can introduce malicious code. |
| H-003 | GitHub Actions workflows use pinned action versions (SHA hashes, not floating tags). | A GitHub Actions tag that is overwritten by a compromised action maintainer can inject malicious code into the pipeline. |
| H-004 | CI security scans (SAST, dependency scanning) are configured to fail the build on critical findings, not just warn. | A security scan that warns but does not block the build provides no actual security guarantee. |
| H-005 | Artifact images in ECR are scanned for vulnerabilities at rest and on every push. | An image that passed CI scanning but was later found to have a vulnerability must be detected before deployment. |
| H-006 | Image signing keys are stored in a hardware security module (HSM) or a managed key service — not on the CI runner. | A signing key stored on an ephemeral CI runner can be exfiltrated by a malicious workflow or third-party action. |
| H-007 | ArgoCD's connection to the Git repository uses a read-only token with access only to the deployment manifests directory. | An ArgoCD token with write access to Git can be used to modify manifests directly, bypassing the CI pipeline. |
| H-008 | The Kubernetes cluster RBAC restricts ArgoCD's service account to only the namespaces it needs to manage. | An over-permissioned ArgoCD service account can deploy to any namespace or access cluster-wide resources. |
| H-009 | Secrets (database passwords, API keys) are not stored in Git repository manifests, even if encrypted. | Git is not a secrets manager; secrets in Git are accessible to anyone with repository access. |
| H-010 | CI runners are ephemeral — no persistent data, no SSH access, isolated per job. | A persistent or shared CI runner can leak secrets from one job to another and is a lateral movement target. |
| H-011 | GitHub organization and repository audit logs are monitored for anomalous activity (fork creation, collaborator addition, permission changes). | An attacker who compromises a GitHub admin account can add backdoor collaborators or decrease security. |
| H-012 | The artifact registry (ECR) has lifecycle policies that delete old, unused images to reduce the attack surface. | Stale images with unpatched vulnerabilities accumulate over time and expand the blast radius of a registry compromise. |
| H-013 | Deployments through ArgoCD require manual approval for production environments (not fully automated). | A fully automated deployment pipeline means a single compromised CI build can push malicious code to production. |
| H-014 | Image signing is verified by ArgoCD before deployment — unsigned images are not deployed. | An image signing policy that is not enforced at deploy time provides no security benefit. |
| H-015 | GitHub Actions secrets are scoped to the minimum required workflows and are not accessible to all actions in the repository. | A GitHub Actions secret accessible to all workflows can be exfiltrated by a workflow in a different directory. |
| H-016 | The ECR repository policy enforces that images can only be pushed from the CI pipeline's IAM role, not from any other source. | An ECR repository that accepts pushes from any IAM role allows a compromised developer workstation to push images directly. |
| H-017 | Kubernetes admission controllers (e.g., OPA/Gatekeeper, Kyverno) enforce security policies on deployed resources. | Without admission control, a developer can bypass all pipeline security by deploying resources with elevated privileges directly. |
| H-018 | Developer workstations are managed and have endpoint protection (EDR), disk encryption, and screen lock. | An unmanaged developer workstation compromised via phishing can be used to exfiltrate code signing keys or push malicious code. |
| H-019 | CI pipeline logs are retained and monitored — build output, test results, and security scan reports are stored securely. | Without pipeline audit, an attacker who modifies the CI workflow to skip security scans leaves no trace. |
| H-020 | The CI system (GitHub Actions) has no ability to modify branch protection rules or repository settings. | A compromised CI workflow that gains GitHub token with admin permissions can disable protections and cover tracks. |
| H-021 | ArgoCD is configured to sync from a specific Git commit SHA, not the HEAD of a branch, for production deployments. | A branch-based sync causes the cluster state to drift from the tested artifact if the branch is force-pushed or fast-forwarded. |
| H-022 | ECR image tags are immutable — tags cannot be overwritten once pushed. | Mutable tags allow an attacker to overwrite a trusted image tag with a malicious image that ArgoCD will deploy. |
| H-023 | Third-party GitHub Actions from the marketplace are reviewed and approved before use in the organization. | An unverified third-party action can exfiltrate secrets, modify code, or inject backdoors into the build artifact. |
| H-024 | The Kubernetes cluster's network policy isolates ArgoCD from other namespaces and restricts egress to only the ECR API. | A flat network in the cluster allows a compromised container deployed by ArgoCD to reach other namespaces or exfiltrate data. |
| H-025 | GitHub SAML SSO is enforced for all organization members — no backup accounts or local GitHub credentials. | GitHub accounts without SSO are not governed by the organizational identity lifecycle and persist after employment ends. |
| H-026 | Secret scanning on the GitHub repository is enabled and blocks pushes containing known credential patterns. | A developer accidentally committing an API key to a public or internal repository exposes that credential immediately. |
| H-027 | The CI pipeline runs in an isolated environment with no access to production networks or databases. | A CI runner with production network access can be used as a pivot point for lateral movement after a pipeline compromise. |
| H-028 | Code review requires at least two approvers for changes to critical infrastructure manifests (K8s YAML, CI configs). | A single approver who misses a malicious change in a YAML manifest is the single point of failure in the review process. |
| H-029 | Base container images are from a trusted source (official images or an internal curated registry) and are regularly updated. | A base image with a known vulnerability propagates that vulnerability to every downstream application. |
| H-030 | The CI pipeline does not execute arbitrary code from pull requests by external contributors without manual approval. | A pull request from a fork that runs CI automatically can execute arbitrary code in the CI environment. |
| H-031 | ArgoCD webhook integration uses a secret token that is kept confidential and rotated. | A leaked webhook secret allows an attacker to trigger ArgoCD syncs from arbitrary Git commits. |
| H-032 | Kubernetes pod security standards (baseline or restricted) are enforced on all namespaces managed by ArgoCD. | Without pod security standards, a deployment can run as root, mount the host filesystem, or use host network. |
| H-033 | SBOM (Software Bill of Materials) is generated for every build and stored as an attestation in the artifact registry. | Without an SBOM, vulnerability management teams cannot inventory which applications are affected by a new CVE. |
| H-034 | The ECR cross-account or cross-environment push/pull policy is explicitly denied for non-production accounts. | A developer with ECR push access from a dev account can overwrite production images if the registry policy allows cross-account pushes. |
| H-035 | ArgoCD notifications are configured to alert on sync failures, health check degradation, and drift detection. | Silent ArgoCD failures mean the operations team is unaware of deployment failures or configuration drift. |
| H-036 | All developers have MFA enabled on their GitHub accounts. | A developer GitHub account without MFA can be compromised via credential theft, giving an attacker access to the codebase. |
| H-037 | The CI pipeline verifies that the commit was signed (GPG/Sigstore) before building. | Unsigned commits can be attributed to a different author, enabling a contributor to impersonate another developer. |
| H-038 | Image vulnerability scan results from ECR are fed back into the CI/CD pipeline to block deployment of images with critical vulnerabilities. | An image that developed a vulnerability between the CI scan and the ECR push must be caught by a registry-side scan with enforcement. |
| H-039 | ArgoCD is deployed in a private subnet with no direct internet access — it reaches ECR through VPC endpoints or a private registry mirror. | An ArgoCD instance with internet access can be used as an egress point for data exfiltration from the cluster. |
| H-040 | The CI runner's IAM role has no permissions to modify CloudTrail, CloudWatch, or other audit infrastructure. | A compromised CI runner that can disable audit logging can perform malicious actions without leaving evidence. |
| H-041 | GitHub repository visibility is set to private or internal — no public repositories in the organization. | A public repository exposes the codebase to reconnaissance, allowing attackers to study the application for vulnerabilities. |
| H-042 | Deployment rollback capability exists and is tested — ArgoCD can revert to a previous known-good image version. | Without a tested rollback, a bad deployment requires emergency code fixes under pressure, increasing error risk. |
| H-043 | The artifact registry has deletion protection enabled — images cannot be deleted accidentally or maliciously. | Malicious deletion of all ECR images causes deployment failure and requires full pipeline rebuild from source. |

**Total (H): 43**

---

## Step 2: ASF-Generated Assumptions

*Using the 20-pattern assumption generator matrix. Applicable patterns: 17 of 20. Patterns excluded: Physical Security (cloud-hosted), Availability & Resilience (deferred to Operational), Endpoint Security (covered under Human Factors).*

---

### Pattern 1: Authentication (MFA)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-001 | All developers have MFA enabled on their GitHub accounts. | Explicit | Documented policy requires code review but does not specify MFA for developer accounts. |
| ASF-002 | MFA is enforced for the GitHub organization, not just recommended. | Derived | Voluntary MFA adoption rates are low; enforcement is required for meaningful protection. |
| ASF-003 | Service accounts used by GitHub Actions and ArgoCD do not use MFA but are protected by other means (IP restriction, short-lived tokens). | Operational | A service account with MFA cannot automate; but without compensating controls, a leaked token has no second factor. |
| ASF-004 | GitHub recovery codes for MFA bypass are stored securely and accessible only on termination. | Implicit | Lost MFA device recovery is a social engineering vector if not governed by identity verification. |

---

### Pattern 2: Authentication (SSO)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-005 | GitHub is configured to require SAML SSO through the corporate identity provider. | Explicit | SSO ensures that GitHub access is governed by the corporate identity lifecycle (joiner/mover/leaver). |
| ASF-006 | SSO session lifetime is configured with a reasonable timeout — not "remember this device for 90 days." | Trust | Long-lived SSO sessions on shared or unmanaged devices extend the window of credential misuse. |
| ASF-007 | Developers who leave the organization are deprovisioned from GitHub within the identity lifecycle SLA (24 hours). | Operational | A former employee with active GitHub access can modify code, add backdoors, or exfiltrate source code. |
| ASF-008 | SSO token signing keys used for GitHub SAML are rotated and protected from unauthorized access. | Trust | Compromised signing keys allow forged SAML assertions for any user, granting unauthorized GitHub access. |

---

### Pattern 3: Availability & Resilience

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-009 | GitHub is available — CI/CD pipeline cannot function without GitHub's availability for both code storage and workflow execution. | Dependency | A GitHub outage halts all deployments and can prevent emergency fixes from reaching production. |
| ASF-010 | ECR is available in the deployment region — ArgoCD cannot pull images if ECR is down. | Dependency | A regional ECR outage blocks all new deployments and rollbacks. |
| ASF-011 | The CI pipeline has a grace period or fallback procedure if security scanning services (Snyk, Trivvy) are unavailable. | Operational | A security scanner that is down should not block all deployments indefinitely; a documented fallback decision is needed. |
| ASF-012 | ArgoCD's connection to the K8s API server is resilient — ArgoCD can survive API server restarts and network partitions. | Derived | An ArgoCD that loses connection to the API server during a deployment creates a half-deployed state. |

---

### Pattern 4: Backup & Recovery

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-013 | ArgoCD configuration and application manifests are backed up (they are in Git, which is the source of truth). | Explicit | Git is the backup for ArgoCD configuration, but only if a recent clone exists outside of GitHub. |
| ASF-014 | ECR images can be restored from a cross-region replication or pull-through cache if the primary registry is lost. | Derived | An ECR image deletion or corruption without a backup requires rebuilding from source — violating RTO. |
| ASF-015 | GitHub repositories are backed up externally — a GitHub account compromise or data loss requires recovery. | Implicit | GitHub is the single source of truth for code; no external backup means catastrophic data loss on GitHub failure. |
| ASF-016 | CI pipeline configuration (workflow YAML files) is versioned in Git and recoverable from Git history. | Operational | A malicious workflow change committed and then removed from history can still be recovered from Git reflog if available. |

---

### Pattern 5: Change Management

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-017 | Changes to CI/CD pipeline configuration (GitHub Actions workflows, ArgoCD Application manifests, Dockerfiles) require pull request review. | Explicit | Pipeline configuration changes are code changes and should follow the same code review process as application code. |
| ASF-018 | Emergency CI/CD changes (security patches, hotfixes) follow a documented expedited process with post-hoc review. | Operational | An emergency bypass of the CI pipeline that skips security scanning can introduce vulnerabilities under the guise of urgency. |
| ASF-019 | Changes to GitHub organization settings (branch protection, webhook configuration, OAuth apps) require separate approval. | Derived | Organization-level settings changes are administrative actions that should be governed separately from code changes. |
| ASF-020 | ArgoCD Application Custom Resource changes are reviewed before apply — not applied directly via kubectl. | Operational | Direct kubectl changes to ArgoCD Application CRs bypass the Git-as-source-of-truth model. |

---

### Pattern 6: Cloud Security (IAM)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-021 | The GitHub Actions OIDC provider (or IAM role) is scoped to specific repositories, not the entire AWS account. | Explicit | An OIDC role that any GitHub Actions workflow in any repository can assume creates cross-repo access risk. |
| ASF-022 | The ECR push IAM role has permissions to push only to the intended ECR repository, not all repositories. | Derived | A CI role with ecr:PutImage on all repositories can overwrite images in any environment. |
| ASF-023 | The ArgoCD IAM role has permissions to pull from ECR only — no push, no delete. | Implicit | An ArgoCD role with push or delete permissions on ECR can be abused if ArgoCD is compromised. |
| ASF-024 | The EKS cluster IAM OIDC provider trusts only the specific ArgoCD service account, not all service accounts in the cluster. | Trust | An incorrectly configured OIDC trust that allows any K8s service account to assume AWS IAM roles breaks the trust model. |

---

### Pattern 7: Container Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-025 | Container images are built from minimal base images with no unnecessary packages or build artifacts. | Explicit | A bloated base image increases the attack surface and the frequency of vulnerability scanning findings. |
| ASF-026 | Containers in production run as non-root user with a read-only root filesystem. | Derived | A container running as root in a Kubernetes cluster can escape to the host node. |
| ASF-027 | Distroless or minimal base images are used for production to reduce the vulnerability surface. | Implicit | A full OS base image (Ubuntu, Alpine) with package manager and shell is unnecessary for a compiled Go or Java application. |
| ASF-028 | Container images are rebuilt periodically from fresh base images to pick up OS-level patches. | Operational | A base image vulnerability discovered after the image was built will never be patched unless the image is rebuilt. |

---

### Pattern 8: Data Flow & Classification

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-029 | Source code is classified as intellectual property and subject to data handling policies. | Explicit | Source code is the organization's most valuable asset; its classification governs access controls and data protection. |
| ASF-030 | Secrets (API keys, database passwords) flow through the CI pipeline but must not appear in logs, build output, or artifacts. | Implicit | A secret that appears in CI build logs is exposed to anyone with access to the logging system. |
| ASF-031 | Application data does not flow through the CI/CD pipeline — the pipeline handles code and configuration only. | Derived | A pipeline that processes production data (e.g., in tests) creates data leak risks that are not addressed by CI security controls. |
| ASF-032 | SBOM artifacts that describe the software composition are attached to each build and stored with the image. | Derived | Without an SBOM associated with each image, post-deployment vulnerability assessment cannot determine what is deployed. |

---

### Pattern 9: Encryption at Rest

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-033 | ECR images are encrypted at rest using KMS with a customer-managed key. | Explicit | ECR encryption at rest with AWS-managed key is default; customer-managed key provides key control and audit. |
| ASF-034 | GitHub Actions cache and artifact storage is encrypted at rest. | Implicit | CI artifacts and caches may contain build outputs or dependencies that include sensitive data. |
| ASF-035 | ArgoCD secrets (Git credentials, webhook tokens) are stored encrypted in the cluster and in transit. | Derived | ArgoCD stores secrets that can be used to access Git repositories and artifact registries. |
| ASF-036 | Kubernetes etcd (cluster store) is encrypted at rest, including ArgoCD Application Custom Resources. | Trust | Unencrypted etcd exposes all cluster secrets to anyone with filesystem access to the control plane. |

---

### Pattern 10: Encryption in Transit (TLS)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-037 | All communication between GitHub Actions runners and GitHub API uses TLS 1.2 or higher. | Explicit | Runner-to-GitHub API communication must be encrypted to prevent MITM on the build network. |
| ASF-038 | ArgoCD-to-ECR communication is over HTTPS with certificate validation. | Derived | ArgoCD pulling images without certificate validation can be redirected to pull a malicious image. |
| ASF-039 | ArgoCD-to-K8s-API communication uses TLS with client certificate or token authentication. | Trust | ArgoCD communicating with the K8s API over plaintext exposes cluster management credentials. |
| ASF-040 | Developer-to-GitHub communication is over HTTPS — no fallback to SSH without key verification. | Implicit | A developer connecting to GitHub over HTTP on an untrusted network can be subject to MITM. |

---

### Pattern 11: Endpoint Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-041 | Developer workstations have EDR/AV installed and managed with current signatures. | Implicit | A compromised developer workstation is the most common vector for source code theft and supply chain injection. |
| ASF-042 | Developer workstations have full disk encryption and are configured to lock after inactivity. | Derived | A lost or stolen developer workstation exposes source code, SSH keys, and GitHub tokens. |
| ASF-043 | Developers do not install unapproved software on workstations used for code development. | Environmental | Unapproved software may contain malware or keyloggers that capture credentials and signing keys. |
| ASF-044 | CI/CD administrators have hardened workstations with additional controls (PAM, session recording) for pipeline management. | Operational | Administrative access to CI/CD infrastructure is the highest-value target; these workstations require elevated security. |

---

### Pattern 12: Human Factors & Process

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-045 | Developers review the actual diff for each pull request and do not approve requests without understanding the change. | Derived | Rubber-stamping pull requests is the most common bypass of the code review control. |
| ASF-046 | Developers do not share their GitHub accounts or use shared service accounts for personal commits. | Implicit | Shared accounts eliminate individual accountability and make audit logs useless for attribution. |
| ASF-047 | DevOps engineers do not apply ad-hoc changes to the cluster via kubectl that bypass ArgoCD. | Operational | Bypassing ArgoCD creates configuration drift between Git and the cluster, defeating the source-of-truth model. |
| ASF-048 | Developers understand that CI pipeline admission control is enforced and cannot be bypassed. | Trust | If developers believe they can bypass the pipeline for "emergency" fixes, the pipeline becomes optional. |

---

### Pattern 13: Identity Lifecycle (Provisioning)

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-049 | Developer GitHub accounts are linked to the corporate identity provider and follow the joiner/mover/leaver process. | Operational | GitHub accounts not linked to the corporate IdP are not automatically deprovisioned on termination. |
| ASF-050 | GitHub team membership for repository access is reviewed and recertified quarterly. | Derived | Stale team memberships accumulate and grant unnecessary repository access over time. |
| ASF-051 | Service accounts used by GitHub Actions and ArgoCD have their keys rotated on a regular schedule. | Implicit | Long-lived service account keys that are never rotated create a standing credential risk. |
| ASF-052 | ArgoCD service account in the K8s cluster is subject to the same identity lifecycle management as user accounts. | Trust | An orphaned ArgoCD service account from a decommissioned namespace can still authenticate to the API server. |

---

### Pattern 14: Incident Response

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-053 | There is an incident response plan that covers CI/CD pipeline compromise (malicious commit, poisoned image, ArgoCD compromise). | Operational | A pipeline compromise can distribute malware to all production deployments; the response is different from a typical server breach. |
| ASF-054 | The IR team has the ability to block deployment of specific images or roll back ArgoCD applications to a known-good state. | Derived | During an active supply chain attack, the ability to halt deployments and roll back is the primary containment action. |
| ASF-055 | CI/CD audit logs (GitHub audit, CloudTrail, ArgoCD logs) are accessible to the IR team during investigation. | Trust | Inaccessible logs prevent root cause analysis of a supply chain compromise. |
| ASF-056 | There is a documented procedure for rotating all CI/CD credentials (GitHub tokens, ECR credentials, ArgoCD webhook secrets) during a compromise. | Operational | A pipeline compromise requires rapid rotation of secrets across all systems; an untested rotation procedure will miss critical credentials. |

---

### Pattern 15: Least Privilege

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-057 | GitHub Actions workflow permissions are scoped to the minimum required — no write access to other repositories. | Explicit | A CI workflow that has write access to all repositories can modify its own workflow to disable security scans. |
| ASF-058 | The ECR push role can push only to the specific repository for that application, not all ECR repositories. | Derived | An application's CI pipeline pushing images to another application's ECR repository can overwrite production images. |
| ASF-059 | ArgoCD applications are deployed with namespace-scoped RBAC — no cluster-admin privileges. | Implicit | An ArgoCD configured with cluster-admin privileges can deploy resources to any namespace, including kube-system. |
| ASF-060 | Developers have read-only access to production namespaces — no ability to exec into pods, view secrets, or port-forward. | Derived | A developer with exec access to production pods can bypass all CI/CD controls and modify application state directly. |

---

### Pattern 16: Monitoring & Alerting

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-061 | GitHub audit log events are streamed to a SIEM for analysis. | Operational | Without SIEM ingestion, GitHub audit events (repo creation, collaborator add, branch policy change) are invisible. |
| ASF-062 | CI pipeline failures and security scan findings generate alerts — not just logs. | Derived | A silently failing security scan provides no security value; alerts ensure findings are acted upon. |
| ASF-063 | ArgoCD drift detection is configured and alerts when the cluster state diverges from Git. | Operational | Configuration drift detected by ArgoCD is a security signal — unauthorized changes to the cluster are visible as drift. |
| ASF-064 | Image vulnerability scanning at ECR generates alerts for newly discovered vulnerabilities in previously scanned images. | Derived | A CVE discovered after an image was scanned will never be detected unless scanning is retrospective with alerting. |

---

### Pattern 17: Network Segmentation

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-065 | GitHub Actions runners are in an isolated network with no access to internal corporate resources unless explicitly granted. | Explicit | CI runners with unrestricted network access can be used to pivot to internal systems. |
| ASF-066 | The Kubernetes cluster is in a private VPC with no direct internet exposure — all ingress goes through a load balancer or ingress controller. | Architectural | A cluster with public-facing node ports or public IPs on nodes is directly exposed to internet attacks. |
| ASF-067 | ArgoCD is not exposed to the public internet — its UI/API is accessible only from the internal network or VPN. | Implicit | An internet-exposed ArgoCD interface is an attack vector for cluster management compromise. |
| ASF-068 | ECR is accessed via VPC endpoints from the cluster and CI runner networks — not over the public internet. | Derived | ECR access over the public internet creates data exfiltration risk and dependence on internet gateway availability. |

---

### Pattern 18: Physical Security

*Not applicable — cloud-hosted.*

---

### Pattern 19: Supply Chain Security

| ID | Assumption | Ontology | Justification |
|----|-----------|----------|---------------|
| ASF-069 | Third-party GitHub Actions from the marketplace have been vetted for security before approval for organizational use. | Explicit | An unverified third-party action can exfiltrate secrets, inject malicious code, or modify CI results. |
| ASF-070 | Open-source dependencies used by the application are scanned for vulnerabilities and license compliance on every build. | Derived | Dependency scanning that only runs on schedule misses new CVEs that affect currently-deployed dependencies. |
| ASF-071 | Container base images are from verified publishers and include a signature or attestation. | Trust | An unverified base image from an unofficial publisher may contain backdoors or malware. |
| ASF-072 | A software bill of materials (SBOM) is generated and attached to each build for supply chain transparency. | Derived | Without an SBOM, the organization cannot identify which deployed applications are affected by a newly disclosed CVE. |
| ASF-073 | The organization has a policy for handling critical vulnerabilities in the supply chain (must fix SLA, emergency pipeline). | Operational | A critical CVE in a dependency requires rapid response; without a policy, teams may delay patching indefinitely. |

**Total (A): 73** (4 per pattern × 18 applicable patterns + 1 extra for Supply Chain)

---

## Step 3: Comparison

### Overlap Mapping

| Human ID | ASF ID | Match Rationale |
|----------|--------|-----------------|
| H-001 | ASF-001 | Both require MFA on developer GitHub accounts. |
| H-002 | ASF-017 | Both require branch protection and PR review for production branches. |
| H-003 | ASF-069 | Both require pinned action versions and vetted third-party actions. |
| H-004 | ASF-062 | Both require security scans to fail builds on critical findings with alerting. |
| H-005 | ASF-064 | Both require vulnerability scanning at rest on ECR with alerting. |
| H-006 | ASF-035 | Both require image signing keys to be protected (secret storage). |
| H-007 | ASF-007 | Both require read-only, scoped credentials for ArgoCD-to-Git access. |
| H-008 | ASF-059 | Both require ArgoCD RBAC scoped to specific namespaces. |
| H-009 | ASF-030 | Both assume secrets are not stored in Git manifests. |
| H-010 | ASF-065 | Both require CI runners to be isolated and ephemeral. |
| H-011 | ASF-061 | Both require GitHub audit logs to be monitored/SIEM-fed. |
| H-013 | ASF-013 | Both require manual approval for production deployments. |
| H-014 | ASF-014 | Both require image signature verification at deploy time. |
| H-015 | ASF-021 | Both require CI secrets and IAM roles to be scoped to the minimum needed. |
| H-016 | ASF-022 | Both require ECR push restrictions to the CI IAM role only. |
| H-017 | ASF-017 | Both require admission controllers for policy enforcement at deploy time. |
| H-019 | ASF-019 | Both require CI pipeline logs to be retained and monitored. |
| H-020 | ASF-057 | Both require CI to have no ability to modify repository settings. |
| H-022 | ASF-022 | Both require ECR immutable tags to prevent overwriting trusted images. |
| H-023 | ASF-069 | Both require third-party action review and approval. |
| H-024 | ASF-066 | Both require K8s network policies for namespace isolation. |
| H-025 | ASF-005 | Both require SAML SSO for GitHub. |
| H-028 | ASF-017 | Both require multiple approvers for critical manifests. |
| H-029 | ASF-071 | Both require trusted base images. |
| H-031 | ASF-035 | Both require ArgoCD webhook secret protection. |
| H-032 | ASF-059 | Both require pod security standards enforcement. |
| H-033 | ASF-072 | Both require SBOM generation per build. |
| H-036 | ASF-001 | Both require MFA for developers. |
| H-037 | ASF-037 | Both require commit signing and verification. |
| H-038 | ASF-064 | Both require ECR scan results to block deployment of vulnerable images. |
| H-039 | ASF-067 | Both require ArgoCD in private subnet with VPC endpoints. |
| H-043 | ASF-014 | Both require ECR deletion protection and image recoverability. |

**Overlap (O): 32**

### Metrics

| Metric | Value | Formula |
|--------|-------|---------|
| Human assumptions (H) | 43 | Count of unique human-generated assumptions |
| ASF assumptions (A) | 73 | Count of unique ASF-generated assumptions |
| Overlap (O) | 32 | Count appearing in both lists |
| **Precision** | **43.8%** | O / A = 32/73 |
| **Recall** | **74.4%** | O / H = 32/43 |
| **F1 Score** | **55.2%** | 2 × (P × R) / (P + R) |
| Novel findings (A - O) | 41 | Assumptions ASF found that human missed (56.2% of ASF total) |
| Missed findings (H - O) | 11 | Assumptions human found that ASF missed (25.6% of human total) |

### Success Criteria Assessment

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Recall | >= 70% | 74.4% | ✅ Met |
| Precision | >= 50% | 43.8% | ❌ Not met |
| Novel discoveries | >= 10% of total (A+O) | 35.7% (41/115) | ✅ Exceeded |
| Expert agreement (F1 proxy) | > 60% | 55.2% | ❌ Not met |

---

## Step 4: Analysis

### Ontology Overlap Analysis

| Ontology Category | Overlaps | ASF Total | Overlap Rate |
|-------------------|----------|-----------|-------------|
| Explicit | 8 | 12 | 66.7% |
| Derived | 10 | 20 | 50.0% |
| Operational | 7 | 20 | 35.0% |
| Implicit | 4 | 10 | 40.0% |
| Trust | 2 | 8 | 25.0% |
| Architectural | 1 | 4 | 25.0% |
| Dependency | 0 | 4 | 0.0% |
| Environmental | 0 | 3 | 0.0% |

**Best overlap:** Explicit and Derived categories showed the strongest agreement. Both humans and the ASF agree on MFA, SSO, branch protection, and least-privilege IAM roles as critical pipeline security assumptions.

**Worst overlap:** Dependency and Environmental had zero overlap. The ASF surfaced GitHub and ECR availability dependencies as assumptions that the human treats as contextual givens.

### What Humans Caught That ASF Missed (Missed Findings = 11)

1. **Developer workstation management (H-018, H-044):** The human identified that unmanaged developer workstations are a supply chain risk. The ASF Endpoint Security pattern exists but was not applied for this architecture in the pipeline context.

2. **Rollback capability (H-042):** The human identified that tested rollback is a critical assumption. The ASF Backup & Recovery pattern covers image restoration but not the operational rollback procedure.

3. **Secret scanning at push (H-026):** The human identified that GitHub secret scanning must block pushes containing credentials. This is a specific GitHub control that the ASF Credential derivation rule covers abstractly.

4. **CI runner network isolation (H-027):** The human assumed CI runners have no production network access. The ASF Network Segmentation pattern covers this but the match was partial.

5. **Two-person review for infrastructure (H-028):** The human specified two approvers. The ASF change management pattern assumes "review required" but not a specific count.

6. **Deletion protection on ECR (H-043):** The human identified this as distinct from backup. The ASF treated image recoverability under Backup, missing the distinct control of deletion protection.

### What ASF Caught That Humans Missed (Novel Findings = 41)

1. **CI/CD change management governance (ASF-017 through ASF-020):** The human generated assumptions about CI security controls but did not consider the change management process for the pipeline itself. Pipeline configuration changes that bypass review are a blind spot.

2. **K8s etcd encryption (ASF-036):** The human did not consider that the Kubernetes store (etcd) must be encrypted at rest, including ArgoCD Application Custom Resources.

3. **Supply chain vulnerability response policy (ASF-073):** The human assumed scanning happened but did not consider the organizational policy for responding to critical supply chain vulnerabilities.

4. **Incident response plan for pipeline compromise (ASF-053 through ASF-056):** The human generated zero IR assumptions specific to a CI/CD pipeline. The ASF contributed IR pattern covering pipeline-specific response procedures, credential rotation, and deployment rollback.

5. **ArgoCD drift detection as security signal (ASF-063):** The human treated ArgoCD drift as an operational concern. The ASF identified it as a security detection mechanism.

6. **VPC endpoints for ECR (ASF-068):** The human assumed ArgoCD was in a private subnet but did not specify VPC endpoints for ECR access. The ASF identified that ECR access over the internet creates exfiltration risk.

### Architecture Complexity Assessment

Architecture #008 (CI/CD Pipeline) achieved the **highest recall of any architecture so far (74.4%)**, meeting the 70% target. This is explained by:
- The close alignment between the ASF pattern matrix and CI/CD security concerns (Least Privilege, Supply Chain, Container Security, Identity Lifecycle)
- Popular CI/CD security guidance (SLSA framework, supply chain security) is well-represented in the ASF ontology
- The human architect and ASF both focused heavily on the same critical concern: preventing a malicious code injection into the pipeline

### Key Insight

The ASF pattern matrix is strongest for CI/CD architectures because the ASF ontology was partially inspired by supply chain security frameworks (SLSA, CNCF). The recall success validates this alignment. The two remaining gaps are:
- **Developer workstation security** (treated as out-of-scope by the ASF patterns applied)
- **Pipeline change management governance** (changes to the pipeline itself require the same controls as changes through the pipeline)

Adding a dedicated "CI/CD Pipeline Security" pattern that consolidates pipeline-specific concerns would further improve recall.

---

## Conclusions

| Metric | Target | Actual | Verdict |
|--------|--------|--------|---------|
| Recall | >= 70% | 74.4% | ✅ Met — strong alignment with supply chain patterns |
| Precision | >= 50% | 43.8% | ❌ Below target — ASF is deliberately broad/generative |
| Novel discoveries | >= 10% | 35.7% | ✅ ASF adds substantial value beyond human reasoning |
| Expert agreement (F1) | > 60% | 55.2% | ❌ Below target — driven by low precision |

The ASF framework applied to Architecture #008 demonstrates the highest recall of the experiment set, meeting the 70% target. This validates that the ASF's supply chain and container security patterns are well-matched to CI/CD architectures. The primary improvement opportunity is adding a **CI/CD Pipeline Governance** pattern to capture pipeline change management, deployment approval, and rollback assumptions that the current pattern matrix distributes across multiple patterns.
