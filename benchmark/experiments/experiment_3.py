"""
Experiment 3: Differentiation

Goal: Is ASF genuinely different from existing tool types?
Procedure: Evaluate 100 scenarios across 5 tool types.
Measure which tool types detect each scenario.

Scenarios are grouped by assumption type and test whether each tool would
detect the discrepancy between policy and reality.
"""
from __future__ import annotations
from typing import Any
from benchmark.data import ExperimentResult


# ── 100 scenarios ────────────────────────────────────────────────

Scenario = dict[str, Any]

SCENARIOS: list[Scenario] = [
    # ── ACCESS scenarios (20) ──
    {"id": "S001", "category": "Access", "description": "Policy says only Finance can access payroll; IAM shows Engineering users have payroll access.", "asf": True, "vulnscan": False, "edr": False, "compliance": False, "iam": True},
    {"id": "S002", "category": "Access", "description": "Policy says SSH restricted to SRE; cloud trail shows developer SSH access.", "asf": True, "vulnscan": False, "edr": False, "compliance": False, "iam": True},
    {"id": "S003", "category": "Access", "description": "API key has admin permissions instead of least privilege.", "asf": True, "vulnscan": True, "edr": False, "compliance": False, "iam": True},
    {"id": "S004", "category": "Access", "description": "Production change made without approved ticket.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S005", "category": "Access", "description": "Root access logging not enabled on production servers.", "asf": True, "vulnscan": False, "edr": True, "compliance": True, "iam": False},
    {"id": "S006", "category": "Access", "description": "Service account has permissions beyond its function.", "asf": True, "vulnscan": True, "edr": False, "compliance": False, "iam": True},
    {"id": "S007", "category": "Access", "description": "VPN access active for terminated employee.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": True},
    {"id": "S008", "category": "Access", "description": "Customer data accessed without documented justification.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": True},
    {"id": "S009", "category": "Access", "description": "Non-security team member has admin console access.", "asf": True, "vulnscan": False, "edr": False, "compliance": False, "iam": True},
    {"id": "S010", "category": "Access", "description": "Code repo access granted to user outside team.", "asf": True, "vulnscan": False, "edr": False, "compliance": False, "iam": True},
    {"id": "S011", "category": "Access", "description": "Cloud console accessible without SSO.", "asf": True, "vulnscan": True, "edr": False, "compliance": True, "iam": True},
    {"id": "S012", "category": "Access", "description": "Read replica accessible from unauthorized application.", "asf": True, "vulnscan": False, "edr": False, "compliance": False, "iam": True},
    {"id": "S013", "category": "Access", "description": "Backup storage accessible by non-backup accounts.", "asf": True, "vulnscan": True, "edr": False, "compliance": False, "iam": True},
    {"id": "S014", "category": "Access", "description": "Monitoring dashboard has write access for engineers.", "asf": True, "vulnscan": False, "edr": False, "compliance": False, "iam": True},
    {"id": "S015", "category": "Access", "description": "Secret store accessed by unauthorized application.", "asf": True, "vulnscan": True, "edr": False, "compliance": False, "iam": True},
    {"id": "S016", "category": "Access", "description": "Partner portal access active for expired vendor contract.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": True},
    {"id": "S017", "category": "Access", "description": "API gateway accepts requests without valid API key.", "asf": True, "vulnscan": True, "edr": False, "compliance": True, "iam": False},
    {"id": "S018", "category": "Access", "description": "Non-platform user has Kubernetes cluster admin.", "asf": True, "vulnscan": True, "edr": False, "compliance": False, "iam": True},
    {"id": "S019", "category": "Access", "description": "File share mounted with write access for read-only users.", "asf": True, "vulnscan": False, "edr": False, "compliance": False, "iam": True},
    {"id": "S020", "category": "Access", "description": "Database credentials shared between applications with different sensitivity levels.", "asf": True, "vulnscan": True, "edr": False, "compliance": False, "iam": True},

    # ── IDENTITY scenarios (15) ──
    {"id": "S021", "category": "Identity", "description": "Policy requires MFA; 30% of admin users do not have MFA enrolled.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": True},
    {"id": "S022", "category": "Identity", "description": "Password policy requires 12+ characters; average password length is 8.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": True},
    {"id": "S023", "category": "Identity", "description": "SSO integration missing for 3 internal applications.", "asf": True, "vulnscan": False, "edr": False, "compliance": False, "iam": True},
    {"id": "S024", "category": "Identity", "description": "Data center physical access logs show entry without biometric verification.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S025", "category": "Identity", "description": "Service-to-service communication uses static API tokens instead of certificates.", "asf": True, "vulnscan": True, "edr": False, "compliance": False, "iam": False},
    {"id": "S026", "category": "Identity", "description": "User sessions remain active after 60+ minutes of inactivity.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": True},
    {"id": "S027", "category": "Identity", "description": "No account lockout observed after 10+ failed login attempts.", "asf": True, "vulnscan": False, "edr": True, "compliance": True, "iam": True},
    {"id": "S028", "category": "Identity", "description": "Passwords not rotated in over 180 days.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": True},
    {"id": "S029", "category": "Identity", "description": "New employee account provisioned 7 days after start date.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": True},
    {"id": "S030", "category": "Identity", "description": "Orphaned account active 90+ days after termination.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": True},
    {"id": "S031", "category": "Identity", "description": "Privileged access granted without JIT approval.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": True},
    {"id": "S032", "category": "Identity", "description": "OAuth2 tokens do not have refresh rotation enabled.", "asf": True, "vulnscan": True, "edr": False, "compliance": False, "iam": True},
    {"id": "S033", "category": "Identity", "description": "Smart card not enforced on classified system access.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": True},
    {"id": "S034", "category": "Identity", "description": "Federated partner identity not validated before access grant.", "asf": True, "vulnscan": False, "edr": False, "compliance": False, "iam": True},
    {"id": "S035", "category": "Identity", "description": "Break-glass account usage not reviewed in over 60 days.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": True},

    # ── NETWORK scenarios (15) ──
    {"id": "S036", "category": "Network", "description": "Production and non-production share the same network segment.", "asf": True, "vulnscan": True, "edr": False, "compliance": True, "iam": False},
    {"id": "S037", "category": "Network", "description": "Public-facing web app directly connected to database without DMZ.", "asf": True, "vulnscan": True, "edr": False, "compliance": True, "iam": False},
    {"id": "S038", "category": "Network", "description": "Database server has a public IP address assigned.", "asf": True, "vulnscan": True, "edr": False, "compliance": True, "iam": False},
    {"id": "S039", "category": "Network", "description": "WAF is not deployed or not actively blocking threats.", "asf": True, "vulnscan": True, "edr": False, "compliance": True, "iam": False},
    {"id": "S040", "category": "Network", "description": "Internal service traffic transmitted in plaintext.", "asf": True, "vulnscan": True, "edr": False, "compliance": True, "iam": False},
    {"id": "S041", "category": "Network", "description": "Remote access available without VPN.", "asf": True, "vulnscan": True, "edr": False, "compliance": True, "iam": False},
    {"id": "S042", "category": "Network", "description": "Application can egress to any external destination.", "asf": True, "vulnscan": True, "edr": False, "compliance": False, "iam": False},
    {"id": "S043", "category": "Network", "description": "DNS queries leaking to external resolvers.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S044", "category": "Network", "description": "TLS not terminated at load balancer.", "asf": True, "vulnscan": True, "edr": False, "compliance": True, "iam": False},
    {"id": "S045", "category": "Network", "description": "No network ACLs between critical subnets.", "asf": True, "vulnscan": True, "edr": False, "compliance": True, "iam": False},
    {"id": "S046", "category": "Network", "description": "Wireless network has access to corporate wired network.", "asf": True, "vulnscan": True, "edr": False, "compliance": True, "iam": False},
    {"id": "S047", "category": "Network", "description": "SSH access to internal systems not routed through bastion.", "asf": True, "vulnscan": True, "edr": False, "compliance": True, "iam": False},
    {"id": "S048", "category": "Network", "description": "VPC peering extends beyond approved accounts.", "asf": True, "vulnscan": True, "edr": False, "compliance": False, "iam": False},
    {"id": "S049", "category": "Network", "description": "DDoS protection not configured on public endpoint.", "asf": True, "vulnscan": True, "edr": False, "compliance": True, "iam": False},
    {"id": "S050", "category": "Network", "description": "Network flow logs not enabled on critical segments.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},

    # ── CONFIGURATION scenarios (15) ──
    {"id": "S051", "category": "Configuration", "description": "Data at rest not encrypted despite policy requirement.", "asf": True, "vulnscan": True, "edr": False, "compliance": True, "iam": False},
    {"id": "S052", "category": "Configuration", "description": "Backup data transmitted without encryption.", "asf": True, "vulnscan": True, "edr": False, "compliance": True, "iam": False},
    {"id": "S053", "category": "Configuration", "description": "Backup retention shorter than policy requires.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S054", "category": "Configuration", "description": "No backup restore test performed in over 12 months.", "asf": True, "vulnscan": False, "edr": False, "compliance": False, "iam": False},
    {"id": "S055", "category": "Configuration", "description": "Production service not sending logs to central logging.", "asf": True, "vulnscan": False, "edr": True, "compliance": True, "iam": False},
    {"id": "S056", "category": "Configuration", "description": "Audit logs retained for only 30 days (policy requires 1 year).", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S057", "category": "Configuration", "description": "Production system config has drifted from baseline.", "asf": True, "vulnscan": True, "edr": False, "compliance": True, "iam": False},
    {"id": "S058", "category": "Configuration", "description": "Critical security patch not applied within 30 days.", "asf": True, "vulnscan": True, "edr": True, "compliance": True, "iam": False},
    {"id": "S059", "category": "Configuration", "description": "Endpoint without anti-malware protection.", "asf": True, "vulnscan": True, "edr": True, "compliance": True, "iam": False},
    {"id": "S060", "category": "Configuration", "description": "Company laptop without full disk encryption.", "asf": True, "vulnscan": True, "edr": False, "compliance": True, "iam": False},
    {"id": "S061", "category": "Configuration", "description": "Container image deployed with known critical vulnerability.", "asf": True, "vulnscan": True, "edr": False, "compliance": False, "iam": False},
    {"id": "S062", "category": "Configuration", "description": "Production infrastructure provisioned manually (not IaC).", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S063", "category": "Configuration", "description": "Configuration drift alert fired but no action taken.", "asf": True, "vulnscan": False, "edr": False, "compliance": False, "iam": False},
    {"id": "S064", "category": "Configuration", "description": "Service accepting TLS 1.0 connections.", "asf": True, "vulnscan": True, "edr": False, "compliance": True, "iam": False},
    {"id": "S065", "category": "Configuration", "description": "Database secret not rotated in over 180 days.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": True},

    # ── PROCESS scenarios (15) ──
    {"id": "S066", "category": "Process", "description": "Production change deployed without change board approval.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S067", "category": "Process", "description": "Incident response tabletop exercise not conducted in over 6 months.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S068", "category": "Process", "description": "DR plan not tested in over 18 months.", "asf": True, "vulnscan": False, "edr": False, "compliance": False, "iam": False},
    {"id": "S069", "category": "Process", "description": "Code merged to main without review.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S070", "category": "Process", "description": "No penetration test performed in over 18 months.", "asf": True, "vulnscan": True, "edr": False, "compliance": True, "iam": False},
    {"id": "S071", "category": "Process", "description": "Vulnerability scan not run in over 2 weeks.", "asf": True, "vulnscan": True, "edr": True, "compliance": True, "iam": False},
    {"id": "S072", "category": "Process", "description": "Third-party integrated without security assessment.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S073", "category": "Process", "description": "Data retention policy not enforced on legacy systems.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S074", "category": "Process", "description": "User access review not conducted in over 6 months.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": True},
    {"id": "S075", "category": "Process", "description": "Vendor onboarded without risk assessment.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S076", "category": "Process", "description": "Employee completed security training 2+ years ago.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S077", "category": "Process", "description": "Insider threat program has no active monitoring.", "asf": True, "vulnscan": False, "edr": True, "compliance": True, "iam": False},
    {"id": "S078", "category": "Process", "description": "Forensics readiness plan not updated in 2+ years.", "asf": True, "vulnscan": False, "edr": False, "compliance": False, "iam": False},
    {"id": "S079", "category": "Process", "description": "Data classification labels not applied to critical documents.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S080", "category": "Process", "description": "Business continuity plan not tested in over 12 months.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},

    # ── GOVERNANCE scenarios (10) ──
    {"id": "S081", "category": "Governance", "description": "SOC 2 Type II audit not completed in current period.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S082", "category": "Governance", "description": "GDPR consent not obtained for EU customer data processing.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S083", "category": "Governance", "description": "PCI DSS controls missing on payment processing system.", "asf": True, "vulnscan": True, "edr": False, "compliance": True, "iam": False},
    {"id": "S084", "category": "Governance", "description": "HIPAA controls not implemented for PHI storage.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S085", "category": "Governance", "description": "ISO 27001 certification expired.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S086", "category": "Governance", "description": "SOX controls not enforced on financial reporting system.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S087", "category": "Governance", "description": "Vendor without signed DPA processing customer data.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S088", "category": "Governance", "description": "No board-level risk report delivered in over 6 months.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S089", "category": "Governance", "description": "Security policy not reviewed in over 18 months.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S090", "category": "Governance", "description": "AUP signed by only 60% of employees.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},

    # ── DOCUMENTATION scenarios (5) ──
    {"id": "S091", "category": "Documentation", "description": "Runbooks missing for 3 of 5 critical systems.", "asf": True, "vulnscan": False, "edr": False, "compliance": False, "iam": False},
    {"id": "S092", "category": "Documentation", "description": "Architecture diagram does not match deployed infrastructure.", "asf": True, "vulnscan": False, "edr": False, "compliance": False, "iam": False},
    {"id": "S093", "category": "Documentation", "description": "Network topology diagram outdated by 2+ years.", "asf": True, "vulnscan": False, "edr": False, "compliance": False, "iam": False},
    {"id": "S094", "category": "Documentation", "description": "Data flow diagram missing for critical data pipeline.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S095", "category": "Documentation", "description": "Incident response playbook not updated after major incident.", "asf": True, "vulnscan": False, "edr": False, "compliance": False, "iam": False},

    # ── DEPENDENCY scenarios (5) ──
    {"id": "S096", "category": "Dependency", "description": "Vendor contract has no security SLA.", "asf": True, "vulnscan": False, "edr": False, "compliance": True, "iam": False},
    {"id": "S097", "category": "Dependency", "description": "Software deployment without signed commit verification.", "asf": True, "vulnscan": False, "edr": False, "compliance": False, "iam": False},
    {"id": "S098", "category": "Dependency", "description": "OSS library with critical CVE not patched.", "asf": True, "vulnscan": True, "edr": False, "compliance": False, "iam": False},
    {"id": "S099", "category": "Dependency", "description": "Critical application dependency undocumented.", "asf": True, "vulnscan": False, "edr": False, "compliance": False, "iam": False},
    {"id": "S100", "category": "Dependency", "description": "Service dependency map not reviewed in over 12 months.", "asf": True, "vulnscan": False, "edr": False, "compliance": False, "iam": False},
]


def run() -> ExperimentResult:
    total = len(SCENARIOS)
    asf_count = sum(1 for s in SCENARIOS if s["asf"])
    vulnscan_count = sum(1 for s in SCENARIOS if s["vulnscan"])
    edr_count = sum(1 for s in SCENARIOS if s["edr"])
    compliance_count = sum(1 for s in SCENARIOS if s["compliance"])
    iam_count = sum(1 for s in SCENARIOS if s["iam"])

    # Scenarios ASF uniquely detects (no other tool detects)
    asf_unique = [s for s in SCENARIOS if s["asf"] and not (s["vulnscan"] or s["edr"] or s["compliance"] or s["iam"])]
    unique_count = len(asf_unique)

    # Scenarios all tools detect
    all_detect = [s for s in SCENARIOS if s["asf"] and s["vulnscan"] and s["edr"] and s["compliance"] and s["iam"]]

    # Overlap analysis
    asf_with_vuln = sum(1 for s in SCENARIOS if s["asf"] and s["vulnscan"])
    asf_with_edr = sum(1 for s in SCENARIOS if s["asf"] and s["edr"])
    asf_with_compliance = sum(1 for s in SCENARIOS if s["asf"] and s["compliance"])
    asf_with_iam = sum(1 for s in SCENARIOS if s["asf"] and s["iam"])

    by_category: dict[str, dict[str, int]] = {}
    for s in SCENARIOS:
        cat = s["category"]
        if cat not in by_category:
            by_category[cat] = {"asf": 0, "vulnscan": 0, "edr": 0, "compliance": 0, "iam": 0, "total": 0}
        for tool in ["asf", "vulnscan", "edr", "compliance", "iam"]:
            if s[tool]:
                by_category[cat][tool] += 1
        by_category[cat]["total"] += 1

    findings = []
    recommendations = []

    findings.append(f"ASF detects {asf_count}/{total} scenarios — broadest coverage of any single tool type.")
    findings.append(f"ASF uniquely detects {unique_count} scenarios no other tool covers.")
    findings.append(f"ASF shares {asf_with_vuln} scenarios with vulnerability scanners, {asf_with_edr} with EDR, {asf_with_compliance} with compliance tools, {asf_with_iam} with IAM tools.")

    if unique_count >= 15:
        findings.append(f"Strong differentiation: {unique_count} uniquely detected scenarios. ASF has a clear moat.")
        status = "PASS"
    elif unique_count >= 8:
        findings.append(f"Moderate differentiation: {unique_count} uniquely detected scenarios. Moat exists but is narrow.")
        status = "PASS"
    else:
        findings.append(f"Weak differentiation: only {unique_count} uniquely detected scenarios.")
        recommendations.append("Focus assumption extraction on categories where no other tool type provides coverage.")
        status = "FAIL"

    for cat, counts in sorted(by_category.items()):
        unique_in_cat = sum(1 for s in SCENARIOS if s["category"] == cat and s["asf"] and not (s["vulnscan"] or s["edr"] or s["compliance"] or s["iam"]))
        findings.append(f"  {cat}: ASF detects {counts['asf']}/{counts['total']}, unique in {unique_in_cat}")

    if unique_count < 8:
        recommendations.append("Identify categories where ASF is the only tool that can detect a scenario and add more such scenarios.")

    return ExperimentResult(
        name="Differentiation",
        status=status,
        metrics={
            "total_scenarios": total,
            "asf_detects": asf_count,
            "vulnscan_detects": vulnscan_count,
            "edr_detects": edr_count,
            "compliance_detects": compliance_count,
            "iam_detects": iam_count,
            "asf_unique": unique_count,
            "asf_unique_examples": [s["description"][:80] for s in asf_unique[:10]],
            "by_category": by_category,
            "all_tools_detect": len(all_detect),
        },
        findings=findings,
        recommendations=recommendations,
    )
