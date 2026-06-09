"""
Assumption Knowledge Base Builder — generates the ASF Assumption Generator Matrix.
Output: CSV (spreadsheet), JSON (machine-readable), Markdown (documentation).

Not code. Knowledge. This produces the methodology artifact.
"""
from __future__ import annotations
import json
import csv
from pathlib import Path
from dataclasses import dataclass, field, asdict
from typing import Any


@dataclass
class HiddenAssumption:
    id: str
    component: str
    depends_on: str
    assumption: str
    pattern: str
    category: str  # trust, identity, dependency, configuration, process, human, architecture, external
    risk: str
    verification_method: str
    example_scenario: str = ""


OUTPUT_DIR = Path(__file__).parent


# ──────────────────────────────────────────────────────────────
# PATTERN 1: Authentication (MFA)
# ──────────────────────────────────────────────────────────────

def derive_mfa() -> list[HiddenAssumption]:
    results = []
    components = [
        ("Primary application login", "MFA provider"),
        ("Admin console access", "MFA provider"),
        ("VPN authentication", "MFA provider"),
        ("SSO identity provider", "MFA provider"),
        ("API endpoint authentication", "MFA provider"),
        ("Cloud console access", "MFA provider"),
        ("Privileged access elevation", "MFA provider"),
        ("Password reset workflow", "MFA provider"),
        ("Third-party application access", "MFA provider"),
        ("CI/CD pipeline access", "MFA provider"),
    ]
    factors = [
        ("TOTP", "time-based one-time password"),
        ("Push notification", "push-based approval"),
        ("SMS code", "SMS-delivered code"),
        ("Hardware token", "physical token"),
        ("Biometric", "biometric verification"),
        ("Recovery code", "recovery workflow"),
    ]
    dimensions = [
        ("availability", "Must be available", "Denial of access during MFA provider outage", "Monitor MFA provider uptime and configure failover"),
        ("enrollment", "Enrollment must be at 100%", "Unenrolled users bypass MFA entirely", "Audit MFA enrollment status across all users monthly"),
        ("enforcement", "MFA is enforced at every auth point", "Partial enforcement creates MFA-free access paths", "Penetration test all authentication points for MFA requirement"),
        ("bypass", "Bypass paths do not exist", "MFA bypass mechanisms are equally exploitable", "Review emergency bypass procedures and access logs"),
        ("fatigue", "MFA fatigue is not exploitable", "Attackers spam MFA prompts until user accepts", "Implement number matching or require entering code from prompt"),
        ("factor_security", "Factor is not interceptable", "Phishing sites can proxy MFA session in real time", "Deploy phishing-resistant MFA (WebAuthn/FIDO2) where possible"),
        ("session", "MFA re-prompt is enforced for sensitive actions", "Long-lived sessions after MFA weaken security", "Require step-up authentication for privilege escalation"),
        ("provisioning", "MFA tokens are provisioned before user needs access", "User cannot authenticate on day one if token not ready", "Automate MFA token provisioning in employee onboarding"),
        ("recovery", "MFA recovery workflow is equally secure", "Recovery via email/SMS is weaker than primary MFA", "Audit MFA recovery method strength against primary factor"),
        ("deprovisioning", "MFA tokens are revoked on offboarding", "Former employee retains MFA token access", "Integrate MFA token revocation into offboarding workflow"),
    ]
    for comp, dep in components:
        for factor_name, factor_desc in factors:
            for dim, dim_cond, dim_risk, dim_verify in dimensions:
                aid = f"MFA-{len(results)+1:04d}"
                assumption = f"{dim_cond} for {comp} using {factor_name} ({factor_desc})."
                risk = dim_risk.format(component=comp)
                verify = dim_verify
                cat = "trust" if dim in ("availability", "reliability") else "identity"
                results.append(HiddenAssumption(
                    id=aid, component=comp, depends_on=dep,
                    assumption=assumption, pattern="Authentication (MFA)",
                    category=cat, risk=risk, verification_method=verify,
                    example_scenario=f"User authenticates to {comp} with {factor_name}"
                ))
    return results


# ──────────────────────────────────────────────────────────────
# PATTERN 2: Authentication (SSO)
# ──────────────────────────────────────────────────────────────

def derive_sso() -> list[HiddenAssumption]:
    results = []
    sso_components = [
        ("Identity Provider (IdP)", "SSO session"),
        ("Service Provider (SP) application", "IdP assertion"),
        ("Federated partner application", "Federation trust"),
        ("Mobile application", "SSO session token"),
        ("API gateway", "SSO bearer token"),
        ("Legacy application", "SSO integration bridge"),
        ("B2B portal", "Federation metadata"),
        ("Directory synchronization", "HR-IdP data feed"),
    ]
    dimensions = [
        ("availability", "IdP must be reachable", "All SSO-dependent applications become unavailable during IdP outage", "Test IdP failover and offline access modes"),
        ("session_lifetime", "SSO session lifetime is appropriate for risk level", "Overly long sessions bypass per-application auth", "Configure session duration based on application risk classification"),
        ("coverage", "SSO covers all applications", "Non-SSO applications use independent credentials", "Audit all applications for SSO integration status quarterly"),
        ("token_security", "SSO tokens are not replayable", "Stolen SAML assertion can be replayed to gain access", "Implement audience restriction, not-before/not-on-or-after validation"),
        ("trust_config", "IdP trust configuration is correct", "Misconfigured IdP trust allows unauthorized access", "Review SP metadata, ACS URL, and entity ID configurations"),
        ("signing", "SAML/OIDC assertions are signed and validated", "Unsigned assertions can be forged by attackers", "Enforce assertion signing and validate signatures on every request"),
        ("partner_onboarding", "Federation partner metadata is up to date", "Expired partner certificate breaks federation authentication", "Monitor partner certificate expiry and rotate before expiration"),
        ("logout", "Single logout propagates to all applications", "Logging out of IdP does not terminate SP sessions", "Test SLO (Single Log Out) across all SPs regularly"),
        ("just_in_time", "JIT provisioning creates correct accounts", "JIT may create accounts with excessive default permissions", "Review JIT provisioning default role assignments"),
        ("attribution", "SSO audit logs attribute actions to specific users", "Shared service accounts after SSO bypass individual attribution", "Configure IdP to send username attribute in assertions"),
    ]
    for comp in sso_components:
        for dim, dim_cond, dim_risk, dim_verify in dimensions:
            results.append(HiddenAssumption(
                id=f"SSO-{len(results)+1:04d}",
                component=comp[0], depends_on=comp[1],
                assumption=f"{dim_cond} for {comp[0]}.",
                pattern="Authentication (SSO)",
                category="identity" if dim != "availability" else "trust",
                risk=dim_risk, verification_method=dim_verify,
                example_scenario=f"{comp[0]} authenticates via SSO"
            ))
    return results


# ──────────────────────────────────────────────────────────────
# PATTERN 3: Network Segmentation
# ──────────────────────────────────────────────────────────────

def derive_network_segmentation() -> list[HiddenAssumption]:
    results = []
    segments = [
        ("Production VPC", "Internet Gateway"),
        ("Production VPC", "NAT Gateway"),
        ("Production subnet", "Route table"),
        ("Non-production VPC", "VPC peering connection"),
        ("Database subnet", "Network ACL"),
        ("Application subnet", "Security group"),
        ("DMZ subnet", "WAF rules"),
        ("Management subnet", "Bastion host"),
        ("Private subnet", "VPC endpoint"),
        ("Transit gateway", "Route propagation"),
    ]
    dimensions = [
        ("isolation", "Segment is isolated from unauthorized networks", "Cross-segment traffic bypasses security controls", "Test network ACLs and security groups with active scanning"),
        ("routing", "Route tables direct traffic to correct destinations", "Misrouted traffic leaves the network boundary", "Audit route tables for unintended default routes"),
        ("acl_bidi", "Network ACLs filter traffic in both directions", "Stateless ACLs only filtering inbound creates outbound bypass", "Verify NACL rules cover both inbound and outbound directions"),
        ("sg_scoping", "Security groups are scoped to specific CIDRs, not 0.0.0.0/0", "Overly permissive security groups expose services to internet", "Review all security groups with 0.0.0.0/0 inbound rules"),
        ("peering", "VPC peering does not extend beyond approved accounts", "Peering across accounts creates unauthorized network paths", "Audit VPC peering connections quarterly"),
        ("endpoint", "VPC endpoints are configured for private access", "Traffic to AWS services traverses internet instead of AWS network", "Verify VPC endpoint policies and route tables"),
        ("subnet_public", "No subnet is unintentionally public", "Subnet with IGW route + public IP is directly internet accessible", "Scan for subnets with 0.0.0.0/0 route to IGW"),
        ("flow_logs", "Network flow logs are enabled for critical segments", "Without flow logs, malicious traffic is invisible", "Enable VPC flow logs for all production subnets"),
    ]
    for comp, dep in segments:
        for dim, dim_cond, dim_risk, dim_verify in dimensions:
            results.append(HiddenAssumption(
                id=f"SEG-{len(results)+1:04d}",
                component=comp, depends_on=dep,
                assumption=f"{dim_cond} for {comp} via {dep}.",
                pattern="Network Segmentation",
                category="network" if dim != "flow_logs" else "configuration",
                risk=dim_risk, verification_method=dim_verify,
                example_scenario=f"Traffic from {comp} traverses {dep}"
            ))
    return results


# ──────────────────────────────────────────────────────────────
# PATTERN 4: Encryption at Rest
# ──────────────────────────────────────────────────────────────

def derive_encryption_at_rest() -> list[HiddenAssumption]:
    results = []
    data_stores = [
        ("RDS database instance", "KMS key"),
        ("S3 bucket", "Bucket encryption policy"),
        ("EBS volume", "Volume encryption"),
        ("DynamoDB table", "DynamoDB encryption"),
        ("ElastiCache cluster", "Encryption at rest config"),
        ("Redshift cluster", "Cluster encryption"),
        ("ECS task definition", "EFS volume encryption"),
        ("Lambda environment variables", "KMS key"),
        ("Secrets Manager secret", "KMS key"),
        ("CloudWatch log group", "Log group encryption"),
    ]
    dimensions = [
        ("key_separation", "Encryption keys are stored separately from encrypted data", "Co-located keys and data renders encryption ineffective", "Verify key management system is separate from data storage"),
        ("key_rotation", "Encryption keys are rotated on schedule", "Long-lived keys increase blast radius of key compromise", "Automate key rotation and verify rotation policy"),
        ("key_access", "Key access is audited", "Unauthorized key usage is undetectable without audit", "Enable key usage logging and monitor for anomalous patterns"),
        ("key_backup", "Key backup exists and is recoverable", "Key loss equals permanent data loss", "Test key backup restoration procedure annually"),
        ("kms_policy", "KMS key policies restrict usage to authorized principals", "Overly permissive key policy allows unauthorized decryption", "Review KMS key policies using IAM Access Analyzer"),
        ("algorithm", "Encryption algorithm is industry standard", "Custom or outdated algorithms may be breakable", "Verify AES-256 or equivalent is used, not DES/RC4/TDEA"),
        ("key_availability", "Key management system is available", "KMS outage blocks all encrypted data access", "Test KMS failover and understand blast radius of KMS region outage"),
        ("key_compromise", "Key compromise incident response exists", "No process for key revocation leaves data exposed", "Document and test key compromise response procedure"),
        ("algorithm_mode", "Encryption mode is secure (GCM, not ECB)", "ECB mode leaks data patterns in encryption output", "Verify encryption algorithm uses authenticated encryption mode"),
        ("hsm_protection", "HSM is not physically compromised", "HSM tampering exposes all keys protected by it", "Verify HSM is in tamper-resistant enclosure with audit trail"),
    ]
    for comp, dep in data_stores:
        for dim, dim_cond, dim_risk, dim_verify in dimensions:
            results.append(HiddenAssumption(
                id=f"ENCR-{len(results)+1:04d}",
                component=comp, depends_on=dep,
                assumption=f"{dim_cond} for {comp} backed by {dep}.",
                pattern="Encryption at Rest",
                category="configuration",
                risk=dim_risk, verification_method=dim_verify,
                example_scenario=f"Data written to {comp} is encrypted via {dep}"
            ))
    return results


# ──────────────────────────────────────────────────────────────
# PATTERN 5: Backup & Recovery
# ──────────────────────────────────────────────────────────────

def derive_backup() -> list[HiddenAssumption]:
    results = []
    backup_components = [
        ("Production database", "Backup service"),
        ("Application configuration", "Backup service"),
        ("File storage volume", "Backup service"),
        ("Container registry", "Backup service"),
        ("DNS zone file", "Backup service"),
        ("SSL certificates", "Backup service"),
        ("IAM configuration", "Backup service"),
        ("Source repository", "Backup service"),
    ]
    dimensions = [
        ("completion", "Backup jobs complete successfully within backup window", "Failed backups silently leave data unprotected", "Monitor backup job completion status and alert on failures"),
        ("integrity", "Backup integrity is verified through restore testing", "Corrupted backups are discovered only during actual disaster", "Perform quarterly restore tests for every critical system"),
        ("rpo", "Recovery Point Objective is achievable", "Backup frequency determines maximum data loss", "Verify backup interval meets RPO requirements"),
        ("rto", "Recovery Time Objective is achievable", "Full restoration takes longer than business can tolerate", "Measure actual restore time against RTO requirement"),
        ("encryption_keys", "Backup encryption keys are accessible during recovery", "Inaccessible encryption keys during disaster equal zero data recovery", "Test key access during recovery drill"),
        ("retention", "Backup retention policy is enforced", "Premature deletion or over-retention both create risk", "Audit backup lifecycle policies quarterly"),
        ("geo_redundancy", "Backups are stored in geographically separate location", "Regional disaster destroys both primary and backup", "Verify cross-region replication is enabled and working"),
        ("air_gap", "Backup storage is logically/physically separate from primary", "Ransomware that reaches primary can also destroy backups", "Implement immutable backup storage with separate access controls"),
        ("sla", "Backup SLA is documented and monitored", "No SLA means no guaranteed recovery capability", "Establish and monitor backup SLAs per system criticality"),
        ("scope", "All critical data is included in backup scope", "Undocumented data sources have no backup protection", "Maintain data inventory and verify backup coverage"),
    ]
    for comp in backup_components:
        for dim, dim_cond, dim_risk, dim_verify in dimensions:
            results.append(HiddenAssumption(
                id=f"BAK-{len(results)+1:04d}",
                component=comp[0], depends_on=comp[1],
                assumption=f"{dim_cond} for {comp[0]}.",
                pattern="Backup & Recovery",
                category="process" if dim in ("completion", "integrity", "rpo", "rto", "sla") else "configuration",
                risk=dim_risk, verification_method=dim_verify,
                example_scenario=f"Disaster recovery scenario for {comp[0]}"
            ))
    return results


# ──────────────────────────────────────────────────────────────
# PATTERN 6: Dependency (Third-party)
# ──────────────────────────────────────────────────────────────

def derive_third_party() -> list[HiddenAssumption]:
    results = []
    vendors = [
        ("Cloud infrastructure provider", "SLA"),
        ("SaaS application", "API availability"),
        ("CDN provider", "Edge network"),
        ("Email security gateway", "Threat intelligence feed"),
        ("SIEM platform", "Log ingestion pipeline"),
        ("Identity provider", "Federation trust"),
        ("Payment processor", "PCI DSS compliance"),
        ("Monitoring platform", "Agent connectivity"),
        ("Incident response retainer", "Contractual response time"),
        ("Security testing vendor", "Testing scope"),
    ]
    dimensions = [
        ("posture", "Vendor security posture remains acceptable over time", "Vendor breach becomes your breach", "Review vendor security assessments annually"),
        ("sla", "Vendor SLA commitments are met", "SLA violations impact business operations without compensation", "Monitor SLA metrics and enforce penalty clauses"),
        ("data_handling", "Vendor data handling complies with policy", "Vendor data processing may violate GDPR/HIPAA/PCI requirements", "Audit vendor data processing agreements and certifications"),
        ("breach_notification", "Vendor breach notification is timely", "Undisclosed vendor breach creates undetected exposure window", "Contractually require 24-72 hour breach notification"),
        ("access_revocation", "Vendor access is promptly revoked when engagement ends", "Former vendor retains access to production systems", "Automate vendor account deprovisioning on contract end"),
        ("dependency_awareness", "Vendor dependencies are documented", "Undocumented vendor dependency creates single point of failure", "Maintain vendor dependency registry and review quarterly"),
        ("exit_plan", "Vendor exit strategy exists", "Inability to leave vendor creates vendor lock-in risk", "Document and test vendor transition plan annually"),
        ("subvendor", "Vendor sub-processors are disclosed and assessed", "Vendor use of sub-processors bypasses security assessment", "Contractually require vendor to disclose all sub-processors"),
        ("compliance_cascade", "Vendor compliance certifications remain valid", "Expired vendor certification creates compliance gap for downstream audits", "Monitor vendor certification expiry dates"),
        ("bcp_testing", "Vendor business continuity plan is tested", "Vendor outage impacts your operations without recourse", "Request vendor BCP test results annually"),
    ]
    for comp, dep in vendors:
        for dim, dim_cond, dim_risk, dim_verify in dimensions:
            results.append(HiddenAssumption(
                id=f"VEN-{len(results)+1:04d}",
                component=comp[0], depends_on=comp[1],
                assumption=f"{dim_cond} for {comp[0]}.",
                pattern="Third-party Dependency",
                category="dependency" if dim != "posture" else "trust",
                risk=dim_risk, verification_method=dim_verify,
                example_scenario=f"Security incident at {comp[0]}"
            ))
    return results


# ──────────────────────────────────────────────────────────────
# PATTERN 7: Access Control (Least Privilege)
# ──────────────────────────────────────────────────────────────

def derive_least_privilege() -> list[HiddenAssumption]:
    results = []
    access_points = [
        ("IAM user", "Permission boundary"),
        ("IAM role", "Trust policy"),
        ("Service account", "Scoped permissions"),
        ("API key", "Usage plan / scope"),
        ("Database user", "GRANT statements"),
        ("Kubernetes ServiceAccount", "RBAC role binding"),
        ("Lambda execution role", "Resource-based policy"),
        ("S3 bucket policy", "Principal restriction"),
        ("Secret store access", "Policy"),
        ("KMS key policy", "Key users"),
    ]
    dimensions = [
        ("accuracy", "Assigned permissions match actual required permissions", "Excessive permissions violate least privilege", "Review IAM Access Advisor metrics and rightsize permissions"),
        ("boundary", "Permission boundaries are enforced", "IAM boundary does not restrict permission escalation paths", "Test that permission boundary stops privilege escalation"),
        ("justification", "Permission requests include business justification", "Unchallenged permission requests accumulate unchecked", "Require ticket number and approval for every permission grant"),
        ("timebound", "Temporary permissions have expiration", "Time-unlimited elevated permissions become standing privileges", "Implement temporary credential mechanisms (STS, JIT)"),
        ("creep_detection", "Permission creep is detected and remediated", "Users accumulate permissions over time without review", "Run access review automation comparing current vs required permissions"),
        ("elevation", "Privilege elevation requires separate auth", "Always-elevated sessions increase blast radius of compromise", "Implement just-in-time elevation with audit trail"),
        ("break_glass", "Break-glass emergency access is monitored", "Emergency access procedures can be abused for unauthorized access", "Audit all break-glass usage within 24 hours and review monthly"),
        ("inheritance", "Group/permission inheritance does not create unintended access", "Nested group membership creates hidden privilege escalation", "Review group hierarchy for unintended inheritance paths"),
        ("cross_account", "Cross-account access is scoped and reviewed", "Compromised account can access resources in other accounts", "Review cross-account IAM roles and trust policies quarterly"),
        ("default_deny", "Default deny is enforced for all access requests", "Allow-by-default access creates standing permissions", "Verify IAM policy evaluation results in explicit deny by default"),
    ]
    for comp, dep in access_points:
        for dim, dim_cond, dim_risk, dim_verify in dimensions:
            results.append(HiddenAssumption(
                id=f"LPO-{len(results)+1:04d}",
                component=comp[0], depends_on=comp[1],
                assumption=f"{dim_cond} for {comp[0]} scoped by {dep}.",
                pattern="Least Privilege",
                category="access",
                risk=dim_risk, verification_method=dim_verify,
                example_scenario=f"{comp[0]} requests access to resource"
            ))
    return results


# ──────────────────────────────────────────────────────────────
# PATTERN 8: Identity Lifecycle (Provisioning)
# ──────────────────────────────────────────────────────────────

def derive_identity_lifecycle() -> list[HiddenAssumption]:
    results = []
    lifecycle_stages = [
        ("Employee onboarding", "HR system"),
        ("Contractor onboarding", "HR system"),
        ("Role change / transfer", "Manager approval"),
        ("Intern onboarding", "HR system"),
        ("Merger / acquisition integration", "Identity federation"),
        ("Service account creation", "Application owner"),
        ("API key generation", "Developer self-service"),
        ("Privileged role assignment", "Manager + security approval"),
    ]
    dimensions = [
        ("timeliness", "Accounts are provisioned before user needs access", "Day-one productivity blocked by access delays", "Measure and monitor time from HR trigger to account creation"),
        ("accuracy", "HR data (department, role, manager) is accurate", "Incorrect HR data propagates to incorrect access grants", "Audit HR data accuracy against employee records quarterly"),
        ("hr_idp_sync", "HR-IdP integration operates in real time", "Delayed sync creates window where access does not match role", "Monitor sync latency between HR system and IdP"),
        ("correct_roles", "Provisioned roles match job function", "Generic onboarding roles grant excessive default permissions", "Review default onboarding role permissions quarterly"),
        ("contractor_expiry", "Contractor accounts have automatic expiry", "Expired contractor accounts remain active indefinitely", "Verify contractor account expiry is enforced within IdP"),
        ("identity_proofing", "Identity proofing is adequate for access level", "Weak identity proofing allows impersonation during account creation", "Review identity proofing requirements based on access level"),
        ("duplicate_detection", "Duplicate identities are detected and merged", "Multiple identities for same person create audit gaps", "Implement duplicate identity detection and reconciliation"),
        ("authoritative_source", "Identity source of truth is authoritative and protected", "Multiple identity sources create conflicts and gaps", "Document identity source of truth for each system"),
    ]
    for comp, dep in lifecycle_stages:
        for dim, dim_cond, dim_risk, dim_verify in dimensions:
            results.append(HiddenAssumption(
                id=f"IDM-{len(results)+1:04d}",
                component=comp, depends_on=dep,
                assumption=f"{dim_cond} for {comp}.",
                pattern="Identity Lifecycle (Provisioning)",
                category="identity",
                risk=dim_risk, verification_method=dim_verify,
                example_scenario=f"New user goes through {comp}"
            ))
    return results


# ──────────────────────────────────────────────────────────────
# PATTERN 9: Monitoring & Alerting
# ──────────────────────────────────────────────────────────────

def derive_monitoring() -> list[HiddenAssumption]:
    results = []
    monitoring_targets = [
        ("Production application logs", "Log agent + pipeline"),
        ("CloudTrail management events", "CloudTrail trail"),
        ("VPC flow logs", "Flow log delivery"),
        ("OS-level syslog", "Syslog forwarder"),
        ("Database audit logs", "Database audit feature"),
        ("Network device logs", "Syslog server"),
        ("Endpoint detection events", "EDR agent"),
        ("WAF access logs", "WAF logging config"),
        ("Load balancer logs", "Access log delivery"),
        ("API gateway execution logs", "CloudWatch logging"),
    ]
    dimensions = [
        ("coverage", "All critical systems are sending logs", "Unlogged systems are invisible during incident investigation", "Audit log source coverage against asset inventory quarterly"),
        ("capacity", "Log storage has sufficient capacity for retention period", "Logs rotated before compliance period due to capacity limits", "Monitor log storage utilization and auto-scale or alert"),
        ("tamper_proof", "Logs are tamper-proof (immutable)", "Attacker covering tracks can modify logs to hide activity", "Implement immutable log storage with separate admin access controls"),
        ("alert_config", "Alerts are configured for key detection events", "Security events without alerts are invisible", "Review alert rule coverage against ATT&CK framework"),
        ("alert_triage", "Alerts are triaged within defined SLA", "Un-triaged alerts accumulate and critical events are missed", "Monitor mean time to triage (MTTT) and alert backlog"),
        ("retention", "Log retention meets compliance requirements", "Logs deleted before compliance period create audit failures", "Verify log retention settings against compliance calendar"),
        ("coverage_env", "Monitoring covers all environments (prod, staging, dev)", "Security events in non-prod environments go undetected", "Extend monitoring agent deployment to all environments"),
        ("correlation", "Log correlation across sources is operational", "Siloed logs miss multi-stage attack sequences", "Test cross-source correlation scenarios quarterly"),
        ("baseline", "Anomaly detection has baseline of normal behavior", "Without baseline, anomaly detection generates excessive noise", "Establish and tune behavioral baselines per system"),
        ("response_integration", "Alerts integrate with response workflow (PagerDuty, Slack, etc.)", "Alerts that nobody sees might as well not exist", "Test alert notification delivery path monthly"),
    ]
    for comp, dep in monitoring_targets:
        for dim, dim_cond, dim_risk, dim_verify in dimensions:
            results.append(HiddenAssumption(
                id=f"MON-{len(results)+1:04d}",
                component=comp, depends_on=dep,
                assumption=f"{dim_cond} for {comp} delivered via {dep}.",
                pattern="Monitoring & Alerting",
                category="configuration",
                risk=dim_risk, verification_method=dim_verify,
                example_scenario=f"Security event occurs in {comp}"
            ))
    return results


# ──────────────────────────────────────────────────────────────
# PATTERN 10: Container Security
# ──────────────────────────────────────────────────────────────

def derive_container_security() -> list[HiddenAssumption]:
    results = []
    container_components = [
        ("Container base image", "Container registry"),
        ("Running container", "Container runtime"),
        ("Kubernetes pod", "Pod security policy"),
        ("Kubernetes namespace", "RBAC binding"),
        ("Helm chart", "Chart repository"),
        ("Container registry", "Registry auth"),
        ("Service mesh sidecar", "mTLS cert"),
        ("Secrets store CSI driver", "KMS key"),
        ("Pod network policy", "CNI plugin"),
        ("Cluster autoscaler", "Cloud provider API"),
    ]
    dimensions = [
        ("base_image_trust", "Container base images are from trusted sources", "Compromised base image injects malware into all derived images", "Scan base image provenance and enforce allowed registries"),
        ("image_scanning", "Images are scanned for vulnerabilities before deployment", "Known vulnerabilities in images reach production", "Block deployment of images with critical vulnerabilities"),
        ("runtime_scanning", "Running containers are scanned for vulnerabilities", "New vulnerabilities discovered after deployment go undetected", "Deploy runtime vulnerability scanning in production cluster"),
        ("registry_access", "Container registry access is controlled", "Unauthorized access to registry allows malicious image injection", "Restrict registry push access and enable image signing"),
        ("k8s_rbac", "Kubernetes RBAC follows least privilege", "Overly permissive RBAC allows cluster-wide compromise from single pod", "Audit RBAC bindings and remove cluster-admin where not needed"),
        ("pod_security", "Pod security standards (PSP/PSS) are enforced", "Privileged containers escape container isolation", "Enforce Pod Security Standards at the namespace level"),
        ("runtime_secure", "Container runtime is not vulnerable to escape", "Container escape vulnerability breaks isolation boundary", "Keep container runtime updated and monitor for CVE disclosures"),
        ("network_policy", "Kubernetes network policies restrict pod communication", "Flat pod network allows lateral movement from compromised pod", "Implement default-deny network policies for all namespaces"),
        ("secret_mgmt", "Kubernetes secrets are encrypted at rest", "Secrets stored as plaintext in etcd are exposed in backup", "Enable etcd encryption and use external secrets store"),
        ("immutable", "Container images are immutable and reproducible", "Mutable tags (latest) create non-reproducible deployments", "Enforce unique image tags (commit SHA) for all deployments"),
    ]
    for comp, dep in container_components:
        for dim, dim_cond, dim_risk, dim_verify in dimensions:
            results.append(HiddenAssumption(
                id=f"CNT-{len(results)+1:04d}",
                component=comp, depends_on=dep,
                assumption=f"{dim_cond} for {comp} via {dep}.",
                pattern="Container Security",
                category="configuration" if dim not in ("k8s_rbac", "registry_access", "immutable") else "access",
                risk=dim_risk, verification_method=dim_verify,
                example_scenario=f"Container deployed with {comp}"
            ))
    return results


# ──────────────────────────────────────────────────────────────
# PATTERN 11: Encryption in Transit (TLS)
# ──────────────────────────────────────────────────────────────

def derive_tls() -> list[HiddenAssumption]:
    results = []
    tls_endpoints = [
        ("Public-facing web application", "TLS certificate"),
        ("Internal API endpoint", "Internal CA"),
        ("Database connection", "TLS certificate"),
        ("Load balancer listener", "TLS termination"),
        ("Service mesh sidecar", "mTLS certificate"),
        ("CDN origin", "Origin certificate"),
        ("Email server", "TLS certificate"),
        ("LDAP directory", "LDAPS certificate"),
        ("Webhook receiver", "TLS certificate"),
        ("IoT device connection", "Device certificate"),
    ]
    dimensions = [
        ("validity", "TLS certificates are valid and not expired", "Expired certificate causes service outage or insecure fallback", "Monitor certificate expiry and auto-renew with 30-day notice"),
        ("enforcement", "TLS is enforced on ALL connections, not optional", "Downgrade attacks can force plaintext fallback", "Enforce HSTS and reject non-TLS connections at load balancer"),
        ("min_version", "Minimum TLS version 1.2 is enforced", "TLS 1.0/1.1 are vulnerable to protocol downgrade attacks", "Configure server to reject TLS versions below 1.2"),
        ("ciphers", "Cipher suites are strong (no RC4, 3DES, CBC)", "Weak ciphers enable decryption by sophisticated attackers", "Audit cipher suite configuration against Mozilla SSL guidelines"),
        ("mtls", "mTLS is used for service-to-service communication", "Without mTLS, server cannot verify client identity", "Implement mTLS via service mesh or certificate-based auth"),
        ("crl_check", "Certificate revocation is checked", "Compromised certificate continues to be trusted", "Enable OCSP stapling or CRL checking on all endpoints"),
        ("private_key", "Private keys are protected (not in code, not world-readable)", "Exposed private key allows traffic decryption and impersonation", "Store private keys in HSM or secrets manager with access audit"),
        ("chain", "Certificate chain is complete and trusted by clients", "Incomplete cert chain causes browser/client trust warnings", "Verify full certificate chain is served including intermediates"),
        ("sni", "SNI routing is correctly configured for multi-domain", "Incorrect SNI routing sends traffic to wrong backend", "Test SNI routing for all hostnames on shared IP"),
        ("ocsp", "OCSP responder is reachable for real-time validation", "Unreachable OCSP responder fails open (accepts all) or closed", "Configure OCSP stapling to avoid OCSP verification failures"),
    ]
    for comp, dep in tls_endpoints:
        for dim, dim_cond, dim_risk, dim_verify in dimensions:
            results.append(HiddenAssumption(
                id=f"TLS-{len(results)+1:04d}",
                component=comp, depends_on=dep,
                assumption=f"{dim_cond} for {comp} using {dep}.",
                pattern="Encryption in Transit (TLS)",
                category="configuration",
                risk=dim_risk, verification_method=dim_verify,
                example_scenario=f"Client connects to {comp}"
            ))
    return results


# ──────────────────────────────────────────────────────────────
# PATTERN 12: Change Management
# ──────────────────────────────────────────────────────────────

def derive_change_management() -> list[HiddenAssumption]:
    results = []
    change_types = [
        ("Production deployment", "Change approval"),
        ("Emergency hotfix", "Emergency change process"),
        ("Infrastructure change (Terraform)", "IaC pipeline"),
        ("Database migration", "Migration script"),
        ("Configuration change", "Change ticket"),
        ("DNS change", "Change ticket"),
        ("Firewall rule change", "Change ticket"),
        ("Certificate rotation", "Change ticket"),
        ("Secret rotation", "Change ticket"),
        ("Network ACL change", "Change ticket"),
    ]
    dimensions = [
        ("process_adherence", "All changes follow the defined change process", "Unauthorized changes create undocumented configuration drift", "Audit change tickets against actual infrastructure changes"),
        ("emergency_review", "Emergency changes are reviewed post-hoc", "Emergency changes bypassing review accumulate risk", "Conduct post-emergency change review within 48 hours"),
        ("documentation", "Change documentation accurately describes what changed", "Inaccurate documentation misleads incident responders", "Verify change documentation completeness within 24 hours of change"),
        ("testing", "Changes are tested in non-production before deployment", "Untested changes in production cause preventable incidents", "Enforce change promotion from dev → staging → prod pipeline"),
        ("rollback", "Rollback procedure exists and is tested", "Failed change without rollback leads to extended outage", "Include and test rollback plan in every change ticket"),
        ("change_window", "Changes occur within approved change windows", "Out-of-window changes bypass normal review", "Alert on infrastructure changes outside approved windows"),
        ("drift_detection", "Configuration drift is detected after change", "Manual changes outside IaC create unreproducible environment", "Run drift detection after every change and remediate automatically"),
        ("approval", "Change approval comes from authorized approver", "Rubber-stamped approvals defeat change management purpose", "Rotate change approver assignments and audit approval patterns"),
        ("notification", "Change notifications reach all affected stakeholders", "Undisclosed changes surprise downstream teams", "Implement automatic notification distribution based on change impact"),
        ("scheduling", "Changes do not conflict with other planned changes", "Conflicting changes cause cascading failures", "Require change scheduling coordination for dependent systems"),
    ]
    for comp, dep in change_types:
        for dim, dim_cond, dim_risk, dim_verify in dimensions:
            results.append(HiddenAssumption(
                id=f"CHG-{len(results)+1:04d}",
                component=comp, depends_on=dep,
                assumption=f"{dim_cond} for {comp} requiring {dep}.",
                pattern="Change Management",
                category="process",
                risk=dim_risk, verification_method=dim_verify,
                example_scenario=f"Team performs {comp}"
            ))
    return results


# ──────────────────────────────────────────────────────────────
# PATTERN 13: Incident Response
# ──────────────────────────────────────────────────────────────

def derive_incident_response() -> list[HiddenAssumption]:
    results = []
    ir_phases = [
        ("Detection alert", "SIEM rule"),
        ("Initial triage", "On-call responder"),
        ("Containment action", "Runbook"),
        ("Forensic acquisition", "Forensic tooling"),
        ("Evidence preservation", "Chain of custody"),
        ("Eradication step", "Remediation script"),
        ("Recovery verification", "Validation test"),
        ("Post-incident review", "Meeting"),
        ("Stakeholder notification", "Communication plan"),
        ("Lessons learned", "Action tracker"),
    ]
    dimensions = [
        ("plan_current", "IR plan is current and reflects current architecture", "Outdated IR plan references decommissioned systems and contacts", "Review and test IR plan at least annually"),
        ("roles_assigned", "IR team roles are assigned with backups", "Unclear role assignment causes confusion during active incident", "Maintain on-call rotation with documented escalation paths"),
        ("communication", "Communication channels are established and tested", "Failure to reach responders delays containment", "Test IR communication channels (phone, Slack, radio) quarterly"),
        ("tooling_ready", "Forensic and containment tooling is pre-deployed", "Time spent acquiring tools during incident is time lost for containment", "Pre-deploy forensic agents and containment automation"),
        ("evidence_chain", "Evidence chain of custody is maintained", "Inadmissible evidence prevents legal action against attacker", "Document evidence handling procedures and train IR team"),
        ("playbooks_current", "Playbooks cover current attack scenarios", "Incident requiring un-practiced response has longer containment time", "Test and update playbooks for top 5 attack scenarios annually"),
        ("remediation_testing", "Remediation steps are tested before execution", "Untested remediation causes additional damage", "Sandbox-test remediation scripts before production execution"),
        ("third_party_coordination", "Third-party (LE, PR, legal) coordination is documented", "Third-party coordination delays create public relations damage", "Pre-establish contacts with legal counsel and law enforcement"),
        ("recovery_verification", "Recovery verification confirms attacker access is removed", "Premature recovery declaration leads to reinfection", "Define recovery verification criteria per incident type"),
        ("lessons_implemented", "Lessons learned from incidents are implemented as improvements", "Repeated incidents recur because corrective actions were not completed", "Track lessons-learned action items with accountability"),
    ]
    for comp, dep in ir_phases:
        for dim, dim_cond, dim_risk, dim_verify in dimensions:
            results.append(HiddenAssumption(
                id=f"IR-{len(results)+1:04d}",
                component=comp, depends_on=dep,
                assumption=f"{dim_cond} for {comp} via {dep}.",
                pattern="Incident Response",
                category="process",
                risk=dim_risk, verification_method=dim_verify,
                example_scenario=f"Security incident triggers {comp}"
            ))
    return results


# ──────────────────────────────────────────────────────────────
# PATTERN 14: Cloud Security (IAM)
# ──────────────────────────────────────────────────────────────

def derive_cloud_iam() -> list[HiddenAssumption]:
    results = []
    cloud_components = [
        ("AWS IAM user", "IAM policy"),
        ("AWS IAM role", "Trust policy"),
        ("AWS S3 bucket", "Bucket policy"),
        ("AWS KMS key", "Key policy"),
        ("AWS Lambda function", "Execution role"),
        ("AWS EC2 instance", "Instance profile"),
        ("AWS RDS instance", "IAM auth config"),
        ("AWS VPC endpoint", "Endpoint policy"),
        ("AWS CloudTrail", "Trail configuration"),
        ("AWS Organization", "SCP policy"),
    ]
    dimensions = [
        ("least_privilege", "IAM policies allow only required actions", "Overly permissive IAM policies enable privilege escalation", "Use IAM Access Analyzer to identify unused permissions"),
        ("role_instead_user", "IAM roles are used instead of long-lived user credentials", "Long-lived access keys create rotation and exposure risk", "Migrate workloads to instance profiles/IRSA with temporary credentials"),
        ("trust_restrictive", "IAM role trust policies are restrictive", "Overly broad trust policy allows unauthorized account to assume role", "Restrict trust policy to specific accounts and external IDs"),
        ("key_rotation", "IAM access keys are rotated regularly", "Stale access keys increase blast radius of credential leak", "Automate access key rotation and monitor key age"),
        ("cloudtrail_coverage", "CloudTrail is enabled across all regions", "Activity in un-trailed regions is invisible to security team", "Enable multi-region CloudTrail trail with organization trail"),
        ("guardduty_review", "GuardDuty findings are reviewed and triaged", "Critical findings in GuardDuty backlog are missed", "Integrate GuardDuty with SIEM and set up alert routing"),
        ("s3_restrictive", "S3 bucket policies are restrictive (no public access)", "Public S3 buckets leak sensitive data continuously", "Enable S3 Block Public Access at account level"),
        ("config_monitoring", "AWS Config rules monitor for compliance drift", "Unmonitored configuration changes create undetected policy violations", "Implement AWS Config conformance packs for security benchmarks"),
        ("service_control", "SCPs enforce guardrails across all accounts", "Member accounts can opt out of security controls without SCP enforcement", "Apply SCPs to prevent security control disabling"),
        ("access_analyzer", "IAM Access Analyzer is enabled and findings reviewed", "Unintended resource exposure goes undetected without Access Analyzer", "Enable IAM Access Analyzer and review findings weekly"),
    ]
    for comp, dep in cloud_components:
        for dim, dim_cond, dim_risk, dim_verify in dimensions:
            results.append(HiddenAssumption(
                id=f"CLD-{len(results)+1:04d}",
                component=comp, depends_on=dep,
                assumption=f"{dim_cond} for {comp} via {dep}.",
                pattern="Cloud Security (IAM)",
                category="access",
                risk=dim_risk, verification_method=dim_verify,
                example_scenario=f"Cloud resource access via {comp}"
            ))
    return results


# ──────────────────────────────────────────────────────────────
# PATTERN 15: Endpoint Security
# ──────────────────────────────────────────────────────────────

def derive_endpoint() -> list[HiddenAssumption]:
    results = []
    endpoint_components = [
        ("Corporate laptop", "EDR agent"),
        ("Server endpoint", "Anti-malware"),
        ("Mobile device", "MDM profile"),
        ("Virtual desktop", "AV agent"),
        ("CI/CD runner", "Hardened image"),
        ("Network device", "Firmware"),
        ("Printer/IoT", "Network segmentation"),
        ("Developer workstation", "Local admin control"),
    ]
    dimensions = [
        ("edr_installed", "EDR agent is installed on all endpoints", "Unprotected endpoint is invisible to security team during incident", "Audit EDR agent coverage against endpoint inventory"),
        ("edr_running", "EDR agent is running and reporting", "Disabled EDR agent provides no protection", "Monitor EDR agent heartbeat and alert on check-in failures"),
        ("antimalware_updated", "Anti-malware signatures are current", "Outdated signatures miss recent malware variants", "Verify signature age < 24 hours across all endpoints"),
        ("patching", "OS and application patches are applied within SLA", "Unpatched vulnerabilities are actively exploited in the wild", "Deploy patch management with automated enforcement"),
        ("disk_encryption", "Full disk encryption is enabled", "Lost/stolen device exposes all data without FDE", "Verify FDE status via MDM or EDR reporting"),
        ("local_admin", "Local admin rights are restricted", "Users with local admin can disable security controls", "Remove local admin rights and implement LAPS for admin passwords"),
        ("app_control", "Application allowlisting/blacklisting is enforced", "Users can install unapproved applications with malware risk", "Deploy application control policies via MDM/GPO"),
        ("usb_control", "USB device control is enforced", "Malware introduced via USB bypasses network controls", "Restrict USB mass storage and enable audit logging"),
    ]
    for comp in endpoint_components:
        for dim, dim_cond, dim_risk, dim_verify in dimensions:
            results.append(HiddenAssumption(
                id=f"END-{len(results)+1:04d}",
                component=comp[0], depends_on=comp[1],
                assumption=f"{dim_cond} for {comp[0]} via {comp[1]}.",
                pattern="Endpoint Security",
                category="configuration",
                risk=dim_risk, verification_method=dim_verify,
                example_scenario=f"Endpoint {comp[0]} is compromised"
            ))
    return results


# ──────────────────────────────────────────────────────────────
# PATTERN 16: Data Flow & Classification
# ──────────────────────────────────────────────────────────────

def derive_data_flow() -> list[HiddenAssumption]:
    results = []
    data_contexts = [
        ("PII data in transit", "Data classification label"),
        ("Payment card data (PCI)", "PCI scope boundary"),
        ("Protected health information (PHI)", "HIPAA controls"),
        ("Intellectual property data", "DLP policy"),
        ("Customer financial data", "SOX controls"),
        ("EU citizen data", "GDPR controls"),
        ("Employee HR data", "HR data policy"),
        ("Authentication secrets", "Secrets management"),
    ]
    dimensions = [
        ("classification", "Data is correctly classified according to policy", "Misclassified data receives inadequate protection", "Audit data classification labels against data content sampling"),
        ("encryption_matching", "Encryption level matches data classification", "Sensitive data without encryption matching classification level is exposed", "Verify encryption configurations align with data classification requirements"),
        ("flow_mapping", "Data flows are mapped and documented", "Undocumented data flows bypass security controls", "Maintain data flow diagrams and update on architecture changes"),
        ("dlp_enforcement", "DLP controls are enforced at data boundaries", "Data exfiltration via email, USB, cloud upload goes undetected", "Deploy DLP at key egress points and monitor violations"),
        ("retention_enforcement", "Data retention policies are enforced per classification", "Over-retained data increases breach impact and compliance liability", "Automate data lifecycle management based on classification"),
        ("purpose_limitation", "Data is used only for stated purpose", "Data used beyond stated purpose violates privacy commitments", "Implement data usage auditing and access justification"),
        ("consent_records", "User consent records are stored and retrievable", "Inability to prove consent creates GDPR violation risk", "Store consent records with timestamp and scope metadata"),
        ("data_sovereignty", "Data remains in approved geographic regions", "Data stored in restricted region violates data sovereignty laws", "Implement data residency controls via cloud provider SCP"),
    ]
    for comp, dep in data_contexts:
        for dim, dim_cond, dim_risk, dim_verify in dimensions:
            results.append(HiddenAssumption(
                id=f"DFL-{len(results)+1:04d}",
                component=comp[0], depends_on=comp[1],
                assumption=f"{dim_cond} for {comp[0]} based on {comp[1]}.",
                pattern="Data Flow & Classification",
                category="governance",
                risk=dim_risk, verification_method=dim_verify,
                example_scenario=f"Data flows through system handling {comp[0]}"
            ))
    return results


# ──────────────────────────────────────────────────────────────
# PATTERN 17: Availability & Resilience
# ──────────────────────────────────────────────────────────────

def derive_availability() -> list[HiddenAssumption]:
    results = []
    availability_components = [
        ("Primary application", "Load balancer"),
        ("Stateful database", "Replica/failover"),
        ("Stateless service", "Auto-scaling group"),
        ("DNS resolution", "DNS failover"),
        ("API gateway", "API Gateway HA"),
        ("Message queue", "Queue replication"),
        ("Cache layer", "Cache replication"),
        ("Storage backend", "Storage replication"),
    ]
    dimensions = [
        ("failover", "Failover mechanism is operational and tested", "Untested failover fails when needed", "Test failover at least quarterly with production-like load"),
        ("capacity", "System has sufficient capacity for peak + failover load", "Overloaded system during failover cascades to complete outage", "Load test to 2x expected peak and monitor utilization"),
        ("single_point", "No single point of failure exists in critical path", "Single component failure causes complete system outage", "Architecture review for single points of failure"),
        ("load_balancing", "Traffic is distributed across healthy instances", "Sticky sessions pin traffic to failing instance", "Verify health check configuration and load balancer draining"),
        ("circuit_breaker", "Circuit breakers prevent cascade failures", "Failing downstream service degrades upstream consumers", "Implement circuit breakers with appropriate thresholds"),
        ("rate_limiting", "Rate limiting protects against traffic spikes", "Unmitigated traffic spike overwhelms backend services", "Configure rate limiting at API gateway and application layers"),
        ("backpressure", "Backpressure mechanisms prevent producer overwhelm", "Producer overwhelming consumer causes memory exhaustion", "Implement message queue backpressure or reactive streams"),
        ("graceful_degradation", "System degrades gracefully under failure", "Non-critical failure causes complete system outage", "Define and test graceful degradation modes for each component"),
    ]
    for comp, dep in availability_components:
        for dim, dim_cond, dim_risk, dim_verify in dimensions:
            results.append(HiddenAssumption(
                id=f"AVL-{len(results)+1:04d}",
                component=comp[0], depends_on=comp[1],
                assumption=f"{dim_cond} for {comp[0]} relying on {comp[1]}.",
                pattern="Availability & Resilience",
                category="architecture",
                risk=dim_risk, verification_method=dim_verify,
                example_scenario=f"Failure scenario for {comp[0]}"
            ))
    return results


# ──────────────────────────────────────────────────────────────
# PATTERN 18: Supply Chain Security
# ──────────────────────────────────────────────────────────────

def derive_supply_chain() -> list[HiddenAssumption]:
    results = []
    supply_chain_components = [
        ("Open source library", "Package repository"),
        ("Container base image", "Image registry"),
        ("CI/CD pipeline", "SCM integration"),
        ("Artifact repository", "Repository access control"),
        ("Signed commit", "Developer GPG key"),
        ("Third-party SDK", "Vendor distribution"),
        ("Infrastructure module (Terraform)", "Module registry"),
        ("Helm chart dependency", "Chart repository"),
    ]
    dimensions = [
        ("oss_vuln_scan", "Open source dependencies are scanned for vulnerabilities", "Known CVEs in dependencies are exploited in production", "Run SCA scanning in CI/CD pipeline and block critical vulnerabilities"),
        ("base_image_trust", "Base images are from verified sources with known provenance", "Compromised base image from untrusted source supplies chain attack", "Enforce image provenance attestation and trusted registry"),
        ("pipeline_integrity", "CI/CD pipeline is hardened against injection", "Compromised CI/CD pipeline deploys malicious code", "Implement CI/CD pipeline security controls (least privilege, signed commits)"),
        ("artifact_signing", "Artifacts are signed and signatures verified before deployment", "Unsigned artifact deployment enables supply chain substitution", "Implement Sigstore/cosign for artifact signing and verification"),
        ("commit_verification", "Signed commits are verified before merge", "Unverified signed commits provide no authenticity guarantee", "Enforce commit signature verification in SCM branch rules"),
        ("sdk_trust", "Third-party SDKs are from verified distribution channels", "Typosquatted SDK with malicious code executed in production", "Maintain software bill of materials (SBOM) for all dependencies"),
        ("dependency_pinning", "Dependency versions are pinned (no floating tags)", "Floating tag automatically pulls malicious updated version", "Pin dependencies to specific versions and audit changes"),
        ("supplier_assessment", "Software suppliers are assessed for security practices", "Supplier with poor security introduces vulnerabilities", "Conduct supplier security assessment before procurement"),
    ]
    for comp, dep in supply_chain_components:
        for dim, dim_cond, dim_risk, dim_verify in dimensions:
            results.append(HiddenAssumption(
                id=f"SCS-{len(results)+1:04d}",
                component=comp[0], depends_on=comp[1],
                assumption=f"{dim_cond} for {comp[0]} sourced from {comp[1]}.",
                pattern="Supply Chain Security",
                category="dependency",
                risk=dim_risk, verification_method=dim_verify,
                example_scenario=f"Software supply chain attack via {comp[0]}"
            ))
    return results


# ──────────────────────────────────────────────────────────────
# PATTERN 19: Physical Security
# ──────────────────────────────────────────────────────────────

def derive_physical() -> list[HiddenAssumption]:
    results = []
    physical_components = [
        ("Data center entrance", "Badge access system"),
        ("Server room door", "Biometric scanner"),
        ("Network closet", "Locking mechanism"),
        ("Wiring/cabling", "Cable management"),
        ("Cooling system", "HVAC redundancy"),
        ("Power supply", "UPS + generator"),
        ("Fire suppression", "Detection system"),
        ("Environmental monitoring", "Sensors"),
    ]
    dimensions = [
        ("access_control", "Physical access is restricted to authorized personnel only", "Unauthorized physical access leads to data theft or destruction", "Audit physical access logs and revoke departed employee badges"),
        ("tailgating_prevention", "Tailgating (piggybacking) is prevented", "Tailgating bypasses badge access control completely", "Implement mantraps or turnstiles at data center entrances"),
        ("visitor_log", "Visitor access is logged and escorted", "Unescorted visitor has unrestricted physical access", "Enforce visitor check-in process and escort policy"),
        ("biometric_redundancy", "Biometric system has alternate authentication method", "Biometric failure denies access or forces security-compromising bypass", "Maintain alternative access method and emergency access protocol"),
        ("environmental_monitoring", "Temperature and humidity are monitored", "Cooling failure causes hardware damage and data loss", "Implement environmental monitoring with automated alerting"),
        ("power_backup", "UPS and generator provide continuous power", "Power loss causes abrupt system shutdown and data corruption", "Test UPS battery capacity and generator under load quarterly"),
        ("fire_suppression", "Fire suppression system is operational", "Fire destroys hardware and data permanently", "Inspect fire suppression system and test detection annually"),
        ("maintenance_access", "Vendor maintenance access is logged and supervised", "Vendor maintenance provides opportunity for data exfiltration", "Supervise all third-party maintenance and log activities"),
    ]
    for comp, dep in physical_components:
        for dim, dim_cond, dim_risk, dim_verify in dimensions:
            results.append(HiddenAssumption(
                id=f"PHY-{len(results)+1:04d}",
                component=comp[0], depends_on=comp[1],
                assumption=f"{dim_cond} for {comp[0]} via {comp[1]}.",
                pattern="Physical Security",
                category="trust",
                risk=dim_risk, verification_method=dim_verify,
                example_scenario=f"Physical access attempt at {comp[0]}"
            ))
    return results


# ──────────────────────────────────────────────────────────────
# PATTERN 20: Human Factors & Process
# ──────────────────────────────────────────────────────────────

def derive_human_factors() -> list[HiddenAssumption]:
    results = []
    human_contexts = [
        ("Security training completion", "Training platform"),
        ("Phishing simulation", "SIM platform"),
        ("Insider threat detection", "UEBA tool"),
        ("Policy acknowledgment", "HR system"),
        ("Access request approval", "Manager judgment"),
        ("Incident reporting", "Employee awareness"),
        ("Clean desk policy", "Employee compliance"),
        ("Social engineering prevention", "Employee awareness"),
    ]
    dimensions = [
        ("training_effectiveness", "Security training changes employee behavior", "Training completion without behavior change provides no risk reduction", "Measure training effectiveness via phishing simulation click rates"),
        ("phishing_resilience", "Employees can identify and report phishing attempts", "Phishing remains primary initial access vector if employees cannot identify", "Conduct quarterly phishing simulations with training for clickers"),
        ("insider_detection", "Insider threat detection detects anomalous behavior", "Malicious insider activity goes undetected without behavioral analytics", "Implement UEBA with baseline of normal behavior per role"),
        ("policy_awareness", "Employees are aware of and acknowledge security policies", "Unaware employees cannot comply with policies", "Require annual policy acknowledgment with quiz on key policies"),
        ("approval_accuracy", "Access request approvers make informed decisions", "Rubber-stamped approvals grant unnecessary access", "Require approver to confirm access justification before approval"),
        ("reporting_culture", "Employees report security incidents without fear of blame", "Unreported incidents delay containment and increase damage", "Establish non-punitive incident reporting culture and anonymous channel"),
        ("procedure_compliance", "Employees follow security procedures in practice", "Security procedures that impede work are routinely bypassed", "Conduct process compliance audits and remove friction points"),
        ("social_engineering_resilience", "Employees can resist social engineering tactics", "Social engineering bypasses technical controls via human manipulation", "Conduct social engineering simulations (phone, in-person, email)"),
    ]
    for comp, dep in human_contexts:
        for dim, dim_cond, dim_risk, dim_verify in dimensions:
            results.append(HiddenAssumption(
                id=f"HUM-{len(results)+1:04d}",
                component=comp[0], depends_on=comp[1],
                assumption=f"{dim_cond} for {comp[0]} via {comp[1]}.",
                pattern="Human Factors & Process",
                category="human",
                risk=dim_risk, verification_method=dim_verify,
                example_scenario=f"Human error scenario: {comp[0]}"
            ))
    return results


# ──────────────────────────────────────────────────────────────
# BUILD: Assemble all assumptions
# ──────────────────────────────────────────────────────────────

def build_all() -> list[HiddenAssumption]:
    all_assumptions: list[HiddenAssumption] = []
    generators = [
        derive_mfa, derive_sso, derive_network_segmentation,
        derive_encryption_at_rest, derive_backup, derive_third_party,
        derive_least_privilege, derive_identity_lifecycle, derive_monitoring,
        derive_container_security, derive_tls, derive_change_management,
        derive_incident_response, derive_cloud_iam, derive_endpoint,
        derive_data_flow, derive_availability, derive_supply_chain,
        derive_physical, derive_human_factors,
    ]
    for gen in generators:
        all_assumptions.extend(gen())
    return all_assumptions


# ──────────────────────────────────────────────────────────────
# EXPORT: CSV, JSON, Markdown
# ──────────────────────────────────────────────────────────────

def export_csv(assumptions: list[HiddenAssumption], path: Path):
    with open(path, "w", newline="") as f:
        w = csv.writer(f)
        w.writerow(["ID", "Pattern", "Component", "Depends On", "Assumption", "Category", "Risk", "Verification Method", "Example Scenario"])
        for a in assumptions:
            w.writerow([a.id, a.pattern, a.component, a.depends_on, a.assumption, a.category, a.risk, a.verification_method, a.example_scenario])
    print(f"  CSV: {path} ({len(assumptions)} assumptions)")


def export_json(assumptions: list[HiddenAssumption], path: Path):
    data = {
        "metadata": {
            "name": "ASF Assumption Generator Matrix",
            "version": "1.0",
            "total_assumptions": len(assumptions),
            "patterns": list(sorted(set(a.pattern for a in assumptions))),
            "categories": list(sorted(set(a.category for a in assumptions))),
        },
        "assumptions": [{
            "id": a.id,
            "pattern": a.pattern,
            "component": a.component,
            "depends_on": a.depends_on,
            "assumption": a.assumption,
            "category": a.category,
            "risk": a.risk,
            "verification_method": a.verification_method,
            "example_scenario": a.example_scenario,
        } for a in assumptions],
    }
    with open(path, "w") as f:
        json.dump(data, f, indent=2)
    print(f"  JSON: {path} ({len(assumptions)} assumptions)")


def export_markdown_summary(assumptions: list[HiddenAssumption], path: Path):
    from collections import Counter
    pattern_counts = Counter(a.pattern for a in assumptions)
    cat_counts = Counter(a.category for a in assumptions)

    lines = []
    lines.append("# ASF Assumption Generator Matrix")
    lines.append("")
    lines.append(f"Total hidden assumptions: **{len(assumptions)}**")
    lines.append("")
    lines.append("## Architecture Patterns")
    lines.append("")
    lines.append("| Pattern | Assumptions |")
    lines.append("|---------|------------|")
    for p, c in sorted(pattern_counts.items()):
        lines.append(f"| {p} | {c} |")
    lines.append("")
    lines.append("## Categories")
    lines.append("")
    lines.append("| Category | Assumptions |")
    lines.append("|----------|------------|")
    for c, count in sorted(cat_counts.items()):
        lines.append(f"| {c} | {count} |")
    lines.append("")
    lines.append("## Pattern Details")
    lines.append("")
    for p in sorted(set(a.pattern for a in assumptions)):
        pattern_as = [a for a in assumptions if a.pattern == p]
        lines.append(f"### {p}")
        lines.append("")
        lines.append(f"**{len(pattern_as)} hidden assumptions**")
        lines.append("")
        lines.append("| ID | Component | Depends On | Assumption | Risk | Verification |")
        lines.append("|----|-----------|------------|------------|------|-------------|")
        for a in pattern_as[:10]:
            lines.append(f"| {a.id} | {a.component} | {a.depends_on} | {a.assumption[:60]}... | {a.risk[:60]}... | {a.verification_method[:60]}... |")
        if len(pattern_as) > 10:
            lines.append(f"| ... | *{len(pattern_as) - 10} more* | | | | |")
        lines.append("")

    with open(path, "w") as f:
        f.write("\n".join(lines))
    print(f"  Markdown: {path}")


# ──────────────────────────────────────────────────────────────
# MAIN
# ──────────────────────────────────────────────────────────────

if __name__ == "__main__":
    print("Building ASF Assumption Knowledge Base...")
    assumptions = build_all()
    print(f"  Generated {len(assumptions)} hidden assumptions from 20 architecture patterns")

    templates_dir = OUTPUT_DIR / "templates"
    templates_dir.mkdir(exist_ok=True)

    export_csv(assumptions, OUTPUT_DIR / "assumption_generator_matrix.csv")
    export_json(assumptions, OUTPUT_DIR / "assumption_generator_matrix.json")
    export_markdown_summary(assumptions, OUTPUT_DIR / "assumption_generator_matrix.md")

    print(f"\nDone. Files in {OUTPUT_DIR}/")
