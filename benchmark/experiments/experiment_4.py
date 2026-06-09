"""
Experiment 4: Executive Value

Goal: Do decision makers care about ASF findings?
Procedure: Generate sample ASF reports and provide framework for
executive evaluation sessions. This experiment requires human
participants and cannot be fully automated.
"""
from __future__ import annotations
from typing import Any
from benchmark.data import BenchmarkResult, ExperimentResult, GroundTruth


# 10 report templates that simulate ASF output for executive review
EXECUTIVE_REPORT_TEMPLATES: list[dict[str, Any]] = [
    {
        "id": "ER-001",
        "title": "Payroll Access Contradiction",
        "finding": "Policy states only Finance employees may access payroll, but evidence shows 3 Engineering users have payroll access.",
        "risk": "Unauthorized access to salary data — potential insider threat and GDPR violation.",
        "severity": "CRITICAL",
        "would_change_decision": True,
    },
    {
        "id": "ER-002",
        "title": "MFA Enrollment Gap",
        "finding": "Policy requires MFA for all admin access, but 40% of admin users lack MFA enrollment.",
        "risk": "Account takeover via credential compromise — most common initial access vector in breaches.",
        "severity": "CRITICAL",
        "would_change_decision": True,
    },
    {
        "id": "ER-003",
        "title": "Untested Backup Assumption",
        "finding": "Policy states backups are performed daily, but no restore test has been performed in 18 months.",
        "risk": "Data loss in disaster scenario — backups may be non-functional and organization is operating on an untested assumption.",
        "severity": "HIGH",
        "would_change_decision": True,
    },
    {
        "id": "ER-004",
        "title": "Network Isolation Contradiction",
        "finding": "Policy states database servers are not internet accessible, but a production database has a public IP.",
        "risk": "Direct data exposure to internet — database accessible to any attacker scanning the IP range.",
        "severity": "CRITICAL",
        "would_change_decision": True,
    },
    {
        "id": "ER-005",
        "title": "Access Review Not Performed",
        "finding": "Policy requires quarterly access reviews, but no review conducted in current quarter.",
        "risk": "Permission creep undetected — users accumulate access privileges over time without review.",
        "severity": "HIGH",
        "would_change_decision": True,
    },
    {
        "id": "ER-006",
        "title": "DR Plan Assumption",
        "finding": "Disaster recovery plan exists but has not been tested in 24 months.",
        "risk": "DR plan may be non-functional — organization assuming protection that does not exist.",
        "severity": "HIGH",
        "would_change_decision": True,
    },
    {
        "id": "ER-007",
        "title": "Service Account Over-Provisioned",
        "finding": "Backup service account has admin privileges across all production systems.",
        "risk": "Compromised backup service account provides unrestricted access to all production data.",
        "severity": "CRITICAL",
        "would_change_decision": True,
    },
    {
        "id": "ER-008",
        "title": "Configuration Drift",
        "finding": "3 production servers have configuration that differs from the security baseline.",
        "risk": "Security controls may be ineffective on non-compliant systems — undetected drift undermines baseline.",
        "severity": "MEDIUM",
        "would_change_decision": False,
    },
    {
        "id": "ER-009",
        "title": "Dependency on Deprecated Service",
        "finding": "Critical application depends on a service that is in end-of-life with no migration plan.",
        "risk": "Unpatched vulnerabilities in EOL software — no security patches available.",
        "severity": "HIGH",
        "would_change_decision": True,
    },
    {
        "id": "ER-010",
        "title": "Documentation Out of Date",
        "finding": "Network topology diagram describes architecture that was decommissioned 18 months ago.",
        "risk": "Incident response teams operating on inaccurate information — delayed containment.",
        "severity": "MEDIUM",
        "would_change_decision": False,
    },
]


def run() -> ExperimentResult:
    critical = sum(1 for r in EXECUTIVE_REPORT_TEMPLATES if r["severity"] == "CRITICAL")
    high = sum(1 for r in EXECUTIVE_REPORT_TEMPLATES if r["severity"] == "HIGH")
    medium = sum(1 for r in EXECUTIVE_REPORT_TEMPLATES if r["severity"] == "MEDIUM")
    would_change = sum(1 for r in EXECUTIVE_REPORT_TEMPLATES if r["would_change_decision"])
    total = len(EXECUTIVE_REPORT_TEMPLATES)

    findings = []
    recommendations = []

    findings.append(f"Generated {total} executive-ready ASF findings ({critical} CRITICAL, {high} HIGH, {medium} MEDIUM).")
    findings.append(f"{would_change}/{total} findings are designed to change decisions — the test of executive value.")

    if would_change / total >= 0.7:
        findings.append("70%+ of findings target decision-changing severity — strong signal for executive value.")
    else:
        findings.append("Less than 70% of findings target decision-changing severity — may need sharper focus.")

    findings.append("")
    findings.append("EXPERIMENT REQUIRES HUMAN PARTICIPATION — Cannot be fully automated.")
    findings.append("")
    findings.append("Recommended Protocol:")
    findings.append("  1. Present reports ER-001 through ER-010 to security managers/architects.")
    findings.append("  2. For each report, ask: 'Would this finding change a decision or investment?'")
    findings.append("  3. Ask: 'Would you pay for a tool that produced this finding automatically?'")
    findings.append("  4. Ask: 'Is this finding available from any existing tool you use?'")
    findings.append("  5. Score: YES = 1, MAYBE = 0.5, NO = 0")
    findings.append("")
    findings.append("Target: average score >= 0.7 across all reviewers.")

    recommendations.append("Schedule 30-min sessions with 3+ security decision makers.")
    recommendations.append("Document verbatim responses — qualitative feedback is as important as scores.")
    recommendations.append("If executives consistently say 'we already know this', pivot to Experiment 2 (novel discovery) results.")
    recommendations.append("If executives consistently say 'this is valuable but we have no budget', probe for what would unlock budget.")

    return ExperimentResult(
        name="Executive Value",
        status="INCONCLUSIVE",
        metrics={
            "total_reports": total,
            "critical_count": critical,
            "high_count": high,
            "medium_count": medium,
            "would_change_count": would_change,
            "reports": EXECUTIVE_REPORT_TEMPLATES,
        },
        findings=findings,
        recommendations=recommendations,
    )
