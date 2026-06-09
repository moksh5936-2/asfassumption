"""
Ground truth assumptions for all 100 policies.
Each policy has 5-15 assumptions totaling 1000+ across the benchmark.
"""
from __future__ import annotations
from benchmark.data import Policy, GroundTruthAssumption, GroundTruth
from asf.models import AssumptionType
from benchmark.policies import get_all_policies

_GT_COUNTER = [0]


def _gt(policy_id: str, text: str, atype: AssumptionType, category: str,
        is_critical: bool = False, commentary: str = "", keywords: list[str] | None = None) -> GroundTruthAssumption:
    _GT_COUNTER[0] += 1
    return GroundTruthAssumption(
        id=f"gt_{_GT_COUNTER[0]:04d}",
        policy_id=policy_id,
        text=text,
        type=atype,
        category=category,
        is_critical=is_critical,
        commentary=commentary,
        keywords=keywords or [],
    )


# ── Common implicit assumption generators ──────────────────────

def _buried_deny(policy_id: str) -> list[GroundTruthAssumption]:
    """Implicit: whatever is restricted implies denial is enforced."""
    return [
        _gt(policy_id, "Deny rules are enforced for all unauthorized access attempts.",
            AssumptionType.ACCESS, "implicit", False,
            "Restriction only works if denial is actually enforced at the control plane."),
        _gt(policy_id, "Authorization decisions are consistent across all enforcement points.",
            AssumptionType.ACCESS, "implicit", False,
            "Policy assumes no gaps in coverage between different access control systems."),
    ]


def _buried_audit(policy_id: str) -> list[GroundTruthAssumption]:
    """Implicit: restricted access is auditable."""
    return [
        _gt(policy_id, "Access attempts are logged for audit purposes.",
            AssumptionType.CONFIGURATION, "implicit", False,
            "Without logging, policy enforcement cannot be verified."),
        _gt(policy_id, "Audit logs are retained and reviewable for compliance.",
            AssumptionType.CONFIGURATION, "implicit", False,
            "Log retention is necessary for post-incident investigation."),
    ]


def _buried_identity(policy_id: str) -> list[GroundTruthAssumption]:
    """Implicit: identity is correctly established."""
    return [
        _gt(policy_id, "User identities are verified before authorization decisions.",
            AssumptionType.IDENTITY, "implicit", False,
            "Authorization presupposes correct identification."),
        _gt(policy_id, "Identity provider is available and correctly configured.",
            AssumptionType.IDENTITY, "implicit", False,
            "IdP downtime breaks all downstream access decisions."),
    ]


def _buried_system_exists(policy_id: str, system: str) -> list[GroundTruthAssumption]:
    """Implicit: the referenced system actually exists and is maintained."""
    return [
        _gt(policy_id, f"The {system} system exists and is operational.",
            AssumptionType.DEPENDENCY, "implicit", False,
            "Policy references a system that must be running and maintained."),
        _gt(policy_id, f"The {system} system is properly configured and maintained.",
            AssumptionType.CONFIGURATION, "implicit", False,
            "System configuration is assumed correct without explicit verification."),
    ]


def _buried_employee(policy_id: str) -> list[GroundTruthAssumption]:
    """Implicit: employee status is accurately tracked."""
    return [
        _gt(policy_id, "Employee status is accurately maintained in the HR system.",
            AssumptionType.IDENTITY, "implicit", False,
            "Access decisions depend on accurate HR data."),
        _gt(policy_id, "HR system integrates with identity provider in real time.",
            AssumptionType.DEPENDENCY, "implicit", False,
            "Delayed HR-IdP sync creates access window violations."),
    ]


def _buried_process(policy_id: str) -> list[GroundTruthAssumption]:
    """Implicit: supporting processes exist."""
    return [
        _gt(policy_id, "Supporting processes for this policy are documented and followed.",
            AssumptionType.PROCESS, "implicit", False,
            "Policy assumes operational processes exist to implement it."),
        _gt(policy_id, "Personnel are trained on this policy consistently.",
            AssumptionType.PROCESS, "implicit", False,
            "Policy compliance depends on human awareness and training."),
    ]


def _derived_risk(policy_id: str, domain: str) -> list[GroundTruthAssumption]:
    """Derived: violations cause specific risks."""
    risk_map = {
        "access": ("unauthorized data access or exfiltration", AssumptionType.ACCESS),
        "identity": ("credential compromise or identity takeover", AssumptionType.IDENTITY),
        "network": ("network breach or lateral movement", AssumptionType.NETWORK),
        "configuration": ("security control bypass due to misconfiguration", AssumptionType.CONFIGURATION),
        "process": ("process failure leading to operational gap", AssumptionType.PROCESS),
        "governance": ("compliance violation or regulatory penalty", AssumptionType.GOVERNANCE),
        "documentation": ("knowledge loss or operational error", AssumptionType.DOCUMENTATION),
        "dependency": ("supply chain failure or service outage", AssumptionType.DEPENDENCY),
    }
    risk_text, risk_type = risk_map.get(domain, ("security gap", AssumptionType.GOVERNANCE))
    return [
        _gt(policy_id, f"Violation of this policy leads to {risk_text}.",
            risk_type, "derived", True,
            "The risk consequence justifies the policy's existence."),
        _gt(policy_id, "No compensating controls exist if this policy fails.",
            AssumptionType.CONFIGURATION, "derived", False,
            "Defense-in-depth requires compensating controls; policy assumes they exist."),
    ]


def _derived_measurement(policy_id: str) -> list[GroundTruthAssumption]:
    """Derived: compliance must be measurable."""
    return [
        _gt(policy_id, "Compliance with this policy is measured and reported.",
            AssumptionType.GOVERNANCE, "derived", False,
            "Unmeasured policies are effectively optional."),
        _gt(policy_id, "Exceptions to this policy are tracked and approved.",
            AssumptionType.GOVERNANCE, "derived", False,
            "Policy assumes exceptions follow formal waiver process."),
    ]


def _derived_scope(policy_id: str, system: str) -> list[GroundTruthAssumption]:
    """Derived: scope of coverage applies universally."""
    return [
        _gt(policy_id, f"This policy covers all instances of {system} without exception.",
            AssumptionType.GOVERNANCE, "derived", False,
            "Partial coverage creates blind spots."),
    ]


# ── Build all ground truth ──────────────────────────────────────

def build_ground_truth() -> GroundTruth:
    policies = get_all_policies()
    all_gt: list[GroundTruthAssumption] = []

    for p in policies:
        pid = p.id
        domain = p.domain
        tags = p.tags
        text = p.text

        # ── Explicit assumption ──
        all_gt.append(_gt(pid, text, _domain_to_type(domain), "explicit", True,
                          "Directly stated security requirement.", tags))

        # ── Buried: denial enforcement ──
        if domain in ("access", "network", "identity"):
            all_gt.extend(_buried_deny(pid))
            all_gt.extend(_buried_audit(pid))

        # ── Buried: identity assumptions ──
        if domain in ("access", "identity", "network"):
            all_gt.extend(_buried_identity(pid))

        # ── Buried: employee / HR accuracy ──
        if "employee" in tags or "provisioning" in tags or "offboarding" in tags:
            all_gt.extend(_buried_employee(pid))

        # ── Buried: system exists ──
        system_tag = _find_system_tag(tags)
        if system_tag:
            all_gt.extend(_buried_system_exists(pid, system_tag))

        # ── Buried: process assumptions ──
        all_gt.extend(_buried_process(pid))

        # ── Derived: risk consequence ──
        all_gt.extend(_derived_risk(pid, domain))

        # ── Derived: measurement ──
        all_gt.extend(_derived_measurement(pid))

        # ── Derived: scope ──
        system_tag = _find_system_tag(tags)
        if system_tag:
            all_gt.extend(_derived_scope(pid, system_tag))

        # ── Domain-specific assumptions ──
        more = _domain_specific_assumptions(pid, domain, text, tags)
        all_gt.extend(more)

        # ── Derived: enforcement reality ──
        more2 = _enforcement_assumptions(pid, domain, text)
        all_gt.extend(more2)

        # ── Concrete scenario: evidence would need to exist ──
        more3 = _evidence_assumptions(pid, domain, text, tags)
        all_gt.extend(more3)

    return GroundTruth(policies=policies, assumptions=all_gt)


def _domain_to_type(domain: str) -> AssumptionType:
    mapping = {
        "access": AssumptionType.ACCESS,
        "identity": AssumptionType.IDENTITY,
        "network": AssumptionType.NETWORK,
        "configuration": AssumptionType.CONFIGURATION,
        "process": AssumptionType.PROCESS,
        "governance": AssumptionType.GOVERNANCE,
        "documentation": AssumptionType.DOCUMENTATION,
        "dependency": AssumptionType.DEPENDENCY,
    }
    return mapping.get(domain, AssumptionType.GOVERNANCE)


def _find_system_tag(tags: list[str]) -> str | None:
    system_keywords = {
        "payroll", "production", "database", "backup", "vpn", "monitoring",
        "kubernetes", "secrets", "vault", "partner", "code-repo", "cloud-console",
        "admin-console", "customer-data", "file-share", "container", "laptops",
        "wireless", "bastion", "vpc", "dns", "load-balancer",
    }
    for t in tags:
        if t in system_keywords:
            return t
    return None


def _domain_specific_assumptions(pid: str, domain: str, text: str, tags: list[str]) -> list[GroundTruthAssumption]:
    more: list[GroundTruthAssumption] = []

    if domain == "access":
        if "least-privilege" in tags:
            more.append(_gt(pid, "Least privilege principle is correctly implemented, not just documented.",
                            AssumptionType.ACCESS, "implicit", True,
                            "Documented principle vs. actual implementation gap is a common risk."))
            more.append(_gt(pid, "Permission creep is detected and remediated.",
                            AssumptionType.PROCESS, "derived", False,
                            "Credentials accumulate permissions over time without review."))
        if "api-keys" in tags:
            more.append(_gt(pid, "API keys are rotated regularly and revoked when no longer needed.",
                            AssumptionType.CONFIGURATION, "implicit", True,
                            "Static API keys are a common credential risk."))
        if "root" in tags:
            more.append(_gt(pid, "Root access is not shared or used for routine tasks.",
                            AssumptionType.PROCESS, "implicit", True,
                            "Shared root credentials defeat accountability."))
        if "audit" in tags or "logging" in tags:
            more.append(_gt(pid, "Audit logs for access events are tamper-proof.",
                            AssumptionType.CONFIGURATION, "implicit", False,
                            "Mutable logs cannot serve as evidence."))
        if "team-membership" in tags or "security-team" in tags:
            more.append(_gt(pid, "Team membership is accurately reflected in access control groups.",
                            AssumptionType.IDENTITY, "implicit", False,
                            "Stale group membership creates unauthorized access."))
        if "sso" in tags:
            more.append(_gt(pid, "SSO sessions are validated for each access request, not cached indefinitely.",
                            AssumptionType.IDENTITY, "implicit", False,
                            "Long-lived SSO sessions bypass individual access controls."))
        if "read-only" in tags:
            more.append(_gt(pid, "Read-only access cannot be escalated to write access.",
                            AssumptionType.CONFIGURATION, "derived", True,
                            "Read-only guarantees are only as strong as the underlying permission model."))

    elif domain == "identity":
        if "mfa" in tags:
            more.append(_gt(pid, "MFA enrollment is enforced for all users, not just optional.",
                            AssumptionType.IDENTITY, "implicit", True,
                            "Optional MFA means most users will not enroll."))
            more.append(_gt(pid, "MFA bypass mechanisms (recovery codes, backup methods) are equally secure.",
                            AssumptionType.CONFIGURATION, "implicit", True,
                            "MFA recovery paths are often weaker than primary method."))
            more.append(_gt(pid, "MFA is enforced at every authentication point, not just the primary login.",
                            AssumptionType.IDENTITY, "derived", False,
                            "MFA gaps in sub-applications bypass the control."))
        if "password" in tags and "complexity" in tags:
            more.append(_gt(pid, "Password complexity requirements do not encourage insecure workarounds.",
                            AssumptionType.PROCESS, "implicit", False,
                            "Overly complex requirements drive password reuse and sticky notes."))
            more.append(_gt(pid, "Passwords are not stored in plaintext or transmitted insecurely.",
                            AssumptionType.CONFIGURATION, "implicit", True,
                            "Policy assumes secure credential storage."))
        if "sso" in tags or "federation" in tags:
            more.append(_gt(pid, "Federated identity providers enforce equivalent authentication standards.",
                            AssumptionType.IDENTITY, "implicit", False,
                            "Weak federation partner authentication weakens overall security."))
        if "certificate" in tags:
            more.append(_gt(pid, "Certificate authorities are properly managed and certificates are renewed before expiry.",
                            AssumptionType.CONFIGURATION, "implicit", False,
                            "Expired certificates cause service disruption."))
        if "jit" in tags:
            more.append(_gt(pid, "Just-in-time approval workflow is responsive enough for legitimate needs.",
                            AssumptionType.PROCESS, "derived", False,
                            "Overly slow JIT approval leads to standing privilege requests."))
        if "session" in tags:
            more.append(_gt(pid, "Session timeout is enforced on all client types including mobile.",
                            AssumptionType.CONFIGURATION, "implicit", False,
                            "Mobile apps often have different session handling."))
        if "orphaned-accounts" in tags:
            more.append(_gt(pid, "Termination notifications are reliably sent from HR to IT.",
                            AssumptionType.PROCESS, "implicit", True,
                            "Manual offboarding processes have high failure rates."))

    elif domain == "network":
        if "segmentation" in tags or "isolation" in tags:
            more.append(_gt(pid, "Network segmentation is enforced at the data link layer, not just documented.",
                            AssumptionType.NETWORK, "implicit", True,
                            "Logical segmentation without enforcement is security theater."))
            more.append(_gt(pid, "Segmentation boundaries are tested for leaks.",
                            AssumptionType.PROCESS, "derived", False,
                            "Segmentation must be verified through active testing."))
        if "internet" in tags or "public" in tags:
            more.append(_gt(pid, "No undocumented direct internet paths exist to protected resources.",
                            AssumptionType.NETWORK, "implicit", True,
                            "Shadow IT and undocumented connections bypass network controls."))
        if "encryption" in tags and "internal" in tags:
            more.append(_gt(pid, "Internal traffic encryption keys are managed separately from external TLS.",
                            AssumptionType.CONFIGURATION, "implicit", False,
                            "Key management for mTLS differs from public certificate management."))
        if "waf" in tags or "ddos" in tags:
            more.append(_gt(pid, "WAF rules are kept current against emerging threat patterns.",
                            AssumptionType.CONFIGURATION, "implicit", False,
                            "Static WAF rules become ineffective against new attack techniques."))
        if "dns" in tags:
            more.append(_gt(pid, "Internal DNS resolution does not leak queries to external resolvers.",
                            AssumptionType.NETWORK, "implicit", False,
                            "DNS leaks reveal internal network structure."))
        if "egress" in tags:
            more.append(_gt(pid, "Egress filtering rules are restrictive enough to prevent data exfiltration.",
                            AssumptionType.NETWORK, "implicit", True,
                            "Permissive egress rules are a common data exfiltration vector."))

    elif domain == "configuration":
        if "encryption" in tags and "at-rest" in tags:
            more.append(_gt(pid, "Encryption keys are stored separately from the encrypted data.",
                            AssumptionType.CONFIGURATION, "implicit", True,
                            "Co-located keys and data renders encryption useless."))
            more.append(_gt(pid, "Encryption is verified through active testing, not just configuration review.",
                            AssumptionType.PROCESS, "derived", False,
                            "Configuration says encrypted; testing proves it."))
        if "backup" in tags:
            more.append(_gt(pid, "Backup integrity is verified through restore testing.",
                            AssumptionType.PROCESS, "derived", True,
                            "Untested backups are not backups."))
            more.append(_gt(pid, "Backup media is stored in a separate physical location.",
                            AssumptionType.CONFIGURATION, "implicit", True,
                            "Co-located backups fail alongside primary data."))
        if "logging" in tags:
            more.append(_gt(pid, "Log storage has sufficient capacity and retention is enforced.",
                            AssumptionType.CONFIGURATION, "implicit", False,
                            "Logs that are rotated away before review are equivalent to no logs."))
        if "patching" in tags:
            more.append(_gt(pid, "Patch deployment does not introduce regressions or outages.",
                            AssumptionType.PROCESS, "implicit", False,
                            "Untested patches can break production systems."))
            more.append(_gt(pid, "Emergency patches are applied outside the standard cycle when needed.",
                            AssumptionType.PROCESS, "derived", False,
                            "Zero-day exploits require immediate response."))
        if "drift-monitoring" in tags or "drift-alerts" in tags:
            more.append(_gt(pid, "Configuration drift alerts are acted upon, not ignored.",
                            AssumptionType.PROCESS, "derived", False,
                            "Alert fatigue leads to missed drift notifications."))
        if "iac" in tags:
            more.append(_gt(pid, "Infrastructure as code state matches actual deployed infrastructure.",
                            AssumptionType.CONFIGURATION, "implicit", True,
                            "Terraform state drift is a known reliability risk."))
        if "container" in tags:
            more.append(_gt(pid, "Container base images are from trusted, verified sources.",
                            AssumptionType.DEPENDENCY, "implicit", True,
                            "Supply chain attacks via base images are an increasing threat."))
        if "tls" in tags:
            more.append(_gt(pid, "TLS certificates are valid and not expired.",
                            AssumptionType.CONFIGURATION, "implicit", True,
                            "Certificate expiry is one of the most common production incidents."))

    elif domain == "process":
        if "testing" in tags:
            more.append(_gt(pid, "Test results are reviewed and acted upon, not just collected.",
                            AssumptionType.PROCESS, "derived", False,
                            "Tests that nobody reviews do not improve security."))
            more.append(_gt(pid, "Test scenarios reflect real-world failure conditions accurately.",
                            AssumptionType.PROCESS, "implicit", False,
                            "Happy-path testing misses real failure modes."))
        if "change-management" in tags or "approval" in tags:
            more.append(_gt(pid, "The change approval process is not bypassed during emergencies.",
                            AssumptionType.PROCESS, "implicit", True,
                            "Emergency changes often skip required approvals."))
            more.append(_gt(pid, "Change records are accurate and complete.",
                            AssumptionType.DOCUMENTATION, "implicit", False,
                            "Incomplete change records hinder incident investigation."))
        if "access-review" in tags:
            more.append(_gt(pid, "Access review findings are remediated within the review cycle.",
                            AssumptionType.PROCESS, "derived", False,
                            "Identified excessive access must be revoked to be effective."))
        if "security-training" in tags:
            more.append(_gt(pid, "Security training content is current with the threat landscape.",
                            AssumptionType.PROCESS, "implicit", False,
                            "Outdated training creates false confidence."))

    elif domain == "governance":
        if "compliance" in tags or "sox" in tags or "hipaa" in tags or "pci-dss" in tags:
            more.append(_gt(pid, "Compliance is validated through independent audit, not self-assessment.",
                            AssumptionType.GOVERNANCE, "implicit", True,
                            "Self-assessments have inherent bias and blind spots."))
            more.append(_gt(pid, "Compliance scope covers all relevant systems without exception.",
                            AssumptionType.GOVERNANCE, "derived", False,
                            "Scope gaps create compliance liabilities."))
        if "gdpr" in tags:
            more.append(_gt(pid, "User consent records are stored and retrievable for audit.",
                            AssumptionType.CONFIGURATION, "implicit", False,
                            "Consent cannot be proven without records."))
        if "policy-review" in tags:
            more.append(_gt(pid, "Policy updates are communicated and acknowledged by all staff.",
                            AssumptionType.PROCESS, "derived", False,
                            "Uncommunicated policy changes are not enforceable."))

    elif domain == "documentation":
        more.append(_gt(pid, "Documentation matches actual deployed state.",
                        AssumptionType.DOCUMENTATION, "implicit", True,
                        "Outdated documentation is worse than no documentation."))
        more.append(_gt(pid, "Documentation is accessible to those who need it.",
                        AssumptionType.DOCUMENTATION, "derived", False,
                        "Inaccessible documentation is equivalent to missing documentation."))

    elif domain == "dependency":
        if "vendors" in tags or "third-party" in tags:
            more.append(_gt(pid, "Vendor security incidents are communicated in a timely manner.",
                            AssumptionType.DEPENDENCY, "implicit", True,
                            "Undisclosed vendor breaches become your breaches."))
            more.append(_gt(pid, "Vendor SLAs are enforceable and audited.",
                            AssumptionType.GOVERNANCE, "derived", False,
                            "Unenforced SLAs provide no actual protection."))
        if "supply-chain" in tags:
            more.append(_gt(pid, "Signed commits are verified before deployment.",
                            AssumptionType.PROCESS, "implicit", False,
                            "Unverified signed commits defeat the purpose of signing."))
        if "oss" in tags:
            more.append(_gt(pid, "Open source vulnerabilities are patched within the scan cycle.",
                            AssumptionType.PROCESS, "derived", False,
                            "Scans without remediation create false security confidence."))

    return more


def _enforcement_assumptions(pid: str, domain: str, text: str) -> list[GroundTruthAssumption]:
    """Assumptions about how policy is enforced in practice."""
    return [
        _gt(pid, "This policy is consistently enforced across all environments.",
            AssumptionType.GOVERNANCE, "derived", False,
            "Inconsistent enforcement creates exploitable gaps."),
        _gt(pid, "Violations of this policy are detected and reported.",
            AssumptionType.CONFIGURATION, "derived", False,
            "Undetected violations are indistinguishable from compliance."),
        _gt(pid, "This policy does not conflict with other security policies.",
            AssumptionType.GOVERNANCE, "implicit", False,
            "Conflicting policies create unresolvable compliance ambiguity."),
    ]


def _evidence_assumptions(pid: str, domain: str, text: str, tags: list[str]) -> list[GroundTruthAssumption]:
    """Assumptions about evidence that would verify this policy."""
    return [
        _gt(pid, "Evidence exists to verify compliance with this policy.",
            AssumptionType.DOCUMENTATION, "derived", True,
            "Unverifiable policies cannot be audited."),
        _gt(pid, "Evidence collection does not introduce new security vulnerabilities.",
            AssumptionType.CONFIGURATION, "implicit", False,
            "Monitoring and logging agents expand the attack surface."),
    ]
