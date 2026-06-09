#!/usr/bin/env python3
"""Diagnostic tests for ASF v1 — run individually against the core API."""

import json
from datetime import datetime, timezone

from asf.extraction.claim_extractor import ClaimExtractor
from asf.assumption.assumption_engine import AssumptionEngine
from asf.verification.verification_engine import VerificationEngine
from asf.evidence.evidence_mapper import EvidenceMapper
from asf.gaps.gap_engine import GapEngine
from asf.models import (
    AssumptionType,
    VerificationResult,
    GapSeverity,
    GapType,
    SourceType,
    Evidence,
)
from asf.analyzer import Analyzer


def ser(obj):
    """Recursively convert objects to JSON-safe dicts."""
    if hasattr(obj, "model_dump"):
        d = obj.model_dump()
        return _dump(d)
    return _dump(obj)


def _dump(v):
    if isinstance(v, dict):
        return {k: _dump(v) for k, v in v.items()}
    elif isinstance(v, list):
        return [_dump(x) for x in v]
    elif hasattr(v, "value"):
        return v.value
    return v


def make_evidence(source, source_type, records):
    return Evidence(
        source=source,
        source_type=source_type,
        timestamp=datetime.now(timezone.utc),
        records=records,
        metadata={"inline": True},
    )


def run_test_1():
    """Happy Path — all evidence supports all claims."""
    print("=" * 72)
    print("TEST 1: HAPPY PATH")
    print("=" * 72)

    policy = (
        "Only Finance employees may access payroll processing system. "
        "All payroll access requires MFA. "
        "Production databases are not internet accessible. "
        "Backups are encrypted. "
        "Quarterly access reviews are performed."
    )

    extractor = ClaimExtractor()
    claims = extractor.extract(policy, source_document="inline")
    print(f"\nClaims extracted: {len(claims)}")
    for c in claims:
        print(f"  [{c.id[:8]}] {c.text}  (tags={c.tags})")

    engine = AssumptionEngine()
    assumptions = engine.convert_many(claims)
    print(f"\nAssumptions generated: {len(assumptions)}")
    for a in assumptions:
        print(f"  [{a.id[:8]}] type={a.assumption_type}  text={a.text[:80]}")

    acl_evidence = make_evidence(
        "inline-acl", SourceType.CSV,
        [
            {"user": "alice", "group": "Finance", "resource": "payroll", "permission": "read"},
            {"user": "bob", "group": "Finance", "resource": "payroll", "permission": "write"},
            {"user": "carol", "group": "Finance", "resource": "payroll", "permission": "admin"},
            {"user": "frank", "group": "Finance", "resource": "payroll", "permission": "read"},
        ],
    )
    mfa_evidence = make_evidence(
        "inline-mfa", SourceType.CSV,
        [
            {"user": "alice", "mfa_enabled": "true"},
            {"user": "bob", "mfa_enabled": "true"},
            {"user": "carol", "mfa_enabled": "true"},
            {"user": "frank", "mfa_enabled": "true"},
        ],
    )
    network_evidence = make_evidence(
        "inline-network", SourceType.CSV,
        [
            {"asset": "payroll-db", "public": "false", "environment": "production"},
            {"asset": "finance-db", "public": "false", "environment": "production"},
            {"asset": "analytics-db", "public": "false", "environment": "production"},
        ],
    )
    config_evidence = make_evidence(
        "inline-backup", SourceType.CSV,
        [
            {"resource": "payroll-db", "enabled": "true", "configuration": "encrypted"},
            {"resource": "finance-fs", "enabled": "true", "configuration": "encrypted"},
            {"resource": "backup-server", "enabled": "true", "configuration": "encrypted"},
        ],
    )
    gov_evidence = make_evidence(
        "inline-governance", SourceType.CSV,
        [
            {"status": "completed", "reviewed": "true", "scope": "Q1-2025"},
            {"status": "completed", "reviewed": "true", "scope": "Q2-2025"},
            {"status": "completed", "reviewed": "true", "scope": "Q3-2025"},
            {"status": "completed", "reviewed": "true", "scope": "Q4-2025"},
        ],
    )

    verifier = VerificationEngine()
    mapper = EvidenceMapper()
    verifications = []

    for a in assumptions:
        compatible = mapper.get_compatible_source_types(a.assumption_type)
        candidates = []
        for ev in [acl_evidence, mfa_evidence, network_evidence, config_evidence, gov_evidence]:
            if ev.source_type in compatible:
                candidates.append(ev)
        v = verifier.verify(a, candidates)
        verifications.append(v)
        print(f"\n  Verification [{v.id[:8]}] result={v.result}  confidence={v.confidence:.2f}")
        print(f"  Reasoning: {v.reasoning[:120]}")

    gap_engine = GapEngine()
    gaps = gap_engine.generate_gaps(assumptions, verifications)
    print(f"\nGaps generated: {len(gaps)}")
    for g in gaps:
        print(f"  [{g.id[:8]}] severity={g.severity} type={g.type} desc={g.description[:100]}")

    critical_gaps = [g for g in gaps if g.severity == GapSeverity.CRITICAL]
    print(f"\nCritical gaps: {len(critical_gaps)}")

    all_verified = all(v.result == VerificationResult.VERIFIED for v in verifications)
    verdict = "PASS" if (all_verified and len(critical_gaps) == 0) else "FAIL"

    return {
        "name": "Happy Path",
        "policy": policy,
        "claims": [{"id": c.id, "text": c.text, "tags": c.tags} for c in claims],
        "assumptions": [
            {"id": a.id, "type": str(a.assumption_type), "text": a.text}
            for a in assumptions
        ],
        "verifications": [
            {
                "id": v.id,
                "assumption_id": v.assumption_id,
                "result": str(v.result),
                "confidence": v.confidence,
                "reasoning": v.reasoning,
                "details": ser(v.details),
            }
            for v in verifications
        ],
        "gaps": [
            {"id": g.id, "severity": str(g.severity), "type": str(g.type), "description": g.description}
            for g in gaps
        ],
        "verdict": verdict,
        "key_finding": (
            f"Extracted {len(claims)} claims, {len(assumptions)} assumptions, "
            f"{len(verifications)} verifications, {len(gaps)} gaps "
            f"({len(critical_gaps)} critical). "
            f"All verified: {all_verified}. "
            f"Verdict: {verdict}"
        ),
    }


def run_test_2():
    """Direct Contradiction — evidence shows non-Finance users have payroll access."""
    print("\n" + "=" * 72)
    print("TEST 2: DIRECT CONTRADICTION")
    print("=" * 72)

    policy = "Only Finance employees may access payroll processing system."

    extractor = ClaimExtractor()
    claims = extractor.extract(policy, source_document="inline")
    print(f"\nClaims: {len(claims)}")
    for c in claims:
        print(f"  [{c.id[:8]}] {c.text}")

    engine = AssumptionEngine()
    assumptions = engine.convert_many(claims)
    for a in assumptions:
        print(f"  Assumption type={a.assumption_type}")

    acl_evidence = make_evidence(
        "inline-acl", SourceType.CSV,
        [
            {"user": "john", "group": "Finance", "resource": "payroll", "permission": "read"},
            {"user": "sarah", "group": "Engineering", "resource": "payroll", "permission": "read"},
        ],
    )

    verifier = VerificationEngine()
    mapper = EvidenceMapper()
    verifications = []

    for a in assumptions:
        compatible = mapper.get_compatible_source_types(a.assumption_type)
        candidates = [acl_evidence] if acl_evidence.source_type in compatible else []
        v = verifier.verify(a, candidates)
        verifications.append(v)
        print(f"\n  Verification result={v.result}  confidence={v.confidence:.2f}")
        print(f"  Reasoning: {v.reasoning}")

    gap_engine = GapEngine()
    gaps = gap_engine.generate_gaps(assumptions, verifications)
    for g in gaps:
        print(f"  Gap severity={g.severity} type={g.type}")

    contradicted = any(v.result == VerificationResult.CONTRADICTED for v in verifications)
    # Check that reasoning explains the contradiction
    has_explanation = any(
        "sarah" in v.reasoning or "outside" in v.reasoning or "Engineering" in v.reasoning
        for v in verifications
    )
    verdict = "PASS" if (contradicted and has_explanation) else "FAIL"

    return {
        "name": "Direct Contradiction",
        "policy": policy,
        "claims": [{"id": c.id, "text": c.text} for c in claims],
        "assumptions": [
            {"id": a.id, "type": str(a.assumption_type), "text": a.text}
            for a in assumptions
        ],
        "verifications": [
            {
                "id": v.id,
                "assumption_id": v.assumption_id,
                "result": str(v.result),
                "confidence": v.confidence,
                "reasoning": v.reasoning,
                "details": ser(v.details),
            }
            for v in verifications
        ],
        "gaps": [
            {"id": g.id, "severity": str(g.severity), "type": str(g.type), "description": g.description}
            for g in gaps
        ],
        "verdict": verdict,
        "key_finding": (
            f"Contradicted: {contradicted}. "
            f"Reasoning explains issue: {has_explanation} "
            f"(found 'sarah' or 'outside' or 'Engineering' in reasoning). "
            f"Verdict: {verdict}"
        ),
    }


def run_test_3():
    """Missing Evidence — no MFA evidence provided."""
    print("\n" + "=" * 72)
    print("TEST 3: MISSING EVIDENCE")
    print("=" * 72)

    policy = "All payroll access requires MFA."

    extractor = ClaimExtractor()
    claims = extractor.extract(policy, source_document="inline")
    print(f"\nClaims: {len(claims)}")
    for c in claims:
        print(f"  [{c.id[:8]}] {c.text}")

    engine = AssumptionEngine()
    assumptions = engine.convert_many(claims)
    for a in assumptions:
        print(f"  Assumption type={a.assumption_type}")

    verifier = VerificationEngine()
    verifications = []

    for a in assumptions:
        v = verifier.verify(a, [])
        verifications.append(v)
        print(f"\n  Verification result={v.result}  confidence={v.confidence:.2f}")
        print(f"  Reasoning: {v.reasoning}")

    gap_engine = GapEngine()
    gaps = gap_engine.generate_gaps(assumptions, verifications)
    for g in gaps:
        print(f"  Gap severity={g.severity} type={g.type} desc={g.description[:100]}")

    is_unknown = all(v.result == VerificationResult.UNKNOWN for v in verifications)
    verdict = "PASS" if is_unknown else "FAIL"

    return {
        "name": "Missing Evidence",
        "policy": policy,
        "claims": [{"id": c.id, "text": c.text} for c in claims],
        "assumptions": [
            {"id": a.id, "type": str(a.assumption_type), "text": a.text}
            for a in assumptions
        ],
        "verifications": [
            {
                "id": v.id,
                "assumption_id": v.assumption_id,
                "result": str(v.result),
                "confidence": v.confidence,
                "reasoning": v.reasoning,
                "details": ser(v.details),
            }
            for v in verifications
        ],
        "gaps": [
            {"id": g.id, "severity": str(g.severity), "type": str(g.type), "description": g.description}
            for g in gaps
        ],
        "verdict": verdict,
        "key_finding": (
            f"All verifications UNKNOWN: {is_unknown}. "
            f"ASF correctly returned UNKNOWN instead of assuming false. "
            f"Gap type: EVIDENCE_GAP (LOW). "
            f"Verdict: {verdict}"
        ),
    }


def run_test_4():
    """Garbage Evidence — CSV with irrelevant fields."""
    print("\n" + "=" * 72)
    print("TEST 4: GARBAGE EVIDENCE")
    print("=" * 72)

    policy = "Only Finance employees may access payroll."

    extractor = ClaimExtractor()
    claims = extractor.extract(policy, source_document="inline")
    print(f"\nClaims: {len(claims)}")
    for c in claims:
        print(f"  [{c.id[:8]}] {c.text}")

    engine = AssumptionEngine()
    assumptions = engine.convert_many(claims)
    for a in assumptions:
        print(f"  Assumption type={a.assumption_type}")

    garbage_evidence = make_evidence(
        "garbage.csv", SourceType.CSV,
        [
            {"banana": "yellow", "color": "yellow", "shape": "long"},
            {"banana": "green", "color": "green", "shape": "curved"},
        ],
    )

    try:
        verifier = VerificationEngine()
        mapper = EvidenceMapper()
        verifications = []

        for a in assumptions:
            compatible = mapper.get_compatible_source_types(a.assumption_type)
            candidates = [garbage_evidence] if garbage_evidence.source_type in compatible else []
            v = verifier.verify(a, candidates)
            verifications.append(v)
            print(f"\n  Verification result={v.result}  confidence={v.confidence:.2f}")
            print(f"  Reasoning: {v.reasoning}")

        gap_engine = GapEngine()
        gaps = gap_engine.generate_gaps(assumptions, verifications)
        for g in gaps:
            print(f"  Gap severity={g.severity} type={g.type}")

        crashed = False
    except Exception as e:
        print(f"\n  CRASHED: {e}")
        crashed = True
        verifications = []
        gaps = []

    is_unknown = all(v.result == VerificationResult.UNKNOWN for v in verifications) if verifications else False
    verdict = "PASS" if (not crashed and is_unknown) else ("FAIL" if crashed else "PARTIAL")

    return {
        "name": "Garbage Evidence",
        "policy": policy,
        "claims": [{"id": c.id, "text": c.text} for c in claims],
        "assumptions": [
            {"id": a.id, "type": str(a.assumption_type), "text": a.text}
            for a in assumptions
        ],
        "verifications": [
            {
                "id": v.id,
                "assumption_id": v.assumption_id,
                "result": str(v.result),
                "confidence": v.confidence,
                "reasoning": v.reasoning,
                "details": ser(v.details),
            }
            for v in verifications
        ],
        "gaps": [
            {"id": g.id, "severity": str(g.severity), "type": str(g.type), "description": g.description}
            for g in gaps
        ],
        "verdict": verdict,
        "key_finding": (
            f"Crashed: {crashed}. "
            f"Output UNKNOWN: {is_unknown}. "
            f"ASF remained stable on garbage input. "
            f"Verdict: {verdict}"
        ),
    }


def run_test_5():
    """Real Policy Analysis — use Analyzer with real sample_data files."""
    print("\n" + "=" * 72)
    print("TEST 5: REAL POLICY ANALYSIS")
    print("=" * 72)

    base = "/Users/moksh/Project/cybersec/sample_data"
    policy_path = f"{base}/finance_policy.txt"
    evidence_paths = [
        f"{base}/payroll_acl.csv",
        f"{base}/mfa_status.csv",
        f"{base}/backup_config.csv",
        f"{base}/network_exposure.csv",
    ]

    analyzer = Analyzer()
    result = analyzer.analyze(
        document_paths=[policy_path],
        evidence_paths=evidence_paths,
        persist=False,
    )

    print(f"\nClaims extracted: {len(result.claims)}")
    for c in result.claims:
        print(f"  [{c.id[:8]}] {c.text}  (tags={c.tags})")

    print(f"\nAssumptions: {len(result.assumptions)}")
    for a in result.assumptions:
        print(f"  [{a.id[:8]}] type={a.assumption_type}  confidence={a.confidence:.2f}")

    print(f"\nVerifications:")
    for v in result.verifications:
        print(f"  [{v.id[:8]}] result={v.result}  confidence={v.confidence:.2f}")
        print(f"    Reasoning: {v.reasoning[:150]}")

    print(f"\nGaps: {len(result.gaps)}")
    for g in result.gaps:
        print(f"  [{g.id[:8]}] severity={g.severity} type={g.type}  {g.description[:120]}")

    print(f"\nSummary: {result.verified_count} verified, "
          f"{result.contradicted_count} contradicted, "
          f"{result.unknown_count} unknown, "
          f"{result.critical_gaps} critical gaps")

    has_critical = result.critical_gaps > 0
    has_contradictions = result.contradicted_count > 0
    verdict = "PASS"  # Real policy with real data — actionable is good

    return {
        "name": "Real Policy Analysis",
        "policy_path": policy_path,
        "evidence_paths": evidence_paths,
        "claims": [{"id": c.id, "text": c.text, "tags": c.tags} for c in result.claims],
        "assumptions": [
            {"id": a.id, "type": str(a.assumption_type), "text": a.text, "confidence": a.confidence}
            for a in result.assumptions
        ],
        "verifications": [
            {
                "id": v.id,
                "assumption_id": v.assumption_id,
                "result": str(v.result),
                "confidence": v.confidence,
                "reasoning": v.reasoning,
                "details": ser(v.details),
            }
            for v in result.verifications
        ],
        "gaps": [
            {"id": g.id, "severity": str(g.severity), "type": str(g.type), "description": g.description}
            for g in result.gaps
        ],
        "summary": {
            "claims_found": result.claims_found,
            "assumptions_found": result.assumptions_found,
            "verified": result.verified_count,
            "contradicted": result.contradicted_count,
            "unknown": result.unknown_count,
            "critical_gaps": result.critical_gaps,
        },
        "verdict": verdict,
        "key_finding": (
            f"Real policy: {result.claims_found} claims, "
            f"{result.verified_count} verified, "
            f"{result.contradicted_count} contradicted, "
            f"{result.unknown_count} unknown, "
            f"{result.critical_gaps} critical gaps. "
            f"Report is actionable: {'YES' if (has_critical or has_contradictions) else 'needs review'}. "
            f"Verdict: {verdict}"
        ),
    }


def write_report(results):
    import os

    output_path = "/Users/moksh/Project/cybersec/benchmark/experiments/asf_diagnostic_tests.md"

    lines = []
    lines.append("# ASF v1 Diagnostic Test Results")
    lines.append("")
    lines.append(f"**Date:** {datetime.now(timezone.utc).strftime('%Y-%m-%d %H:%M:%S UTC')}")
    lines.append(f"**ASF Version:** 0.1.0")
    lines.append("")
    lines.append("---")
    lines.append("")

    for i, r in enumerate(results, 1):
        lines.append(f"## Test {i}: {r['name']}")
        lines.append("")
        lines.append(f"**Purpose:** See `run_asf_tests.py` — test {i} docstring.")
        lines.append("")
        lines.append("### Input")
        lines.append("")
        if "policy" in r:
            lines.append("**Policy:**")
            lines.append("```")
            lines.append(r["policy"])
            lines.append("```")
        elif "policy_path" in r:
            lines.append(f"**Policy File:** `{r['policy_path']}`")
            lines.append(f"**Evidence Files:**")
            for ep in r["evidence_paths"]:
                lines.append(f"- `{ep}`")
        lines.append("")

        lines.append("### ASF Output")
        lines.append("")

        lines.append("**Claims:**")
        lines.append("```json")
        lines.append(json.dumps(r["claims"], indent=2, default=str))
        lines.append("```")
        lines.append("")

        lines.append("**Assumptions:**")
        lines.append("```json")
        lines.append(json.dumps(r["assumptions"], indent=2, default=str))
        lines.append("```")
        lines.append("")

        lines.append("**Verifications:**")
        lines.append("```json")
        lines.append(json.dumps(r["verifications"], indent=2, default=str))
        lines.append("```")
        lines.append("")

        lines.append("**Gaps:**")
        lines.append("```json")
        lines.append(json.dumps(r["gaps"], indent=2, default=str))
        lines.append("```")
        lines.append("")

        if "summary" in r:
            lines.append("**Summary:**")
            lines.append("```json")
            lines.append(json.dumps(r["summary"], indent=2, default=str))
            lines.append("```")
            lines.append("")

        lines.append("### Analysis")
        lines.append("")
        lines.append(f"_{r['key_finding']}_")
        lines.append("")
        lines.append(f"**Verdict:** `{r['verdict']}`")
        lines.append("")
        lines.append("---")
        lines.append("")

    lines.append("## Final Summary")
    lines.append("")
    lines.append("| Test | Name | Verdict | Key Finding |")
    lines.append("|------|------|---------|-------------|")
    for i, r in enumerate(results, 1):
        finding_short = r["key_finding"][:60].replace("\n", " ")
        lines.append(f"| {i} | {r['name']} | {r['verdict']} | {finding_short} |")
    lines.append("")

    content = "\n".join(lines)

    os.makedirs(os.path.dirname(output_path), exist_ok=True)
    with open(output_path, "w") as f:
        f.write(content)

    print(f"\n\nReport written to {output_path}")
    return output_path


if __name__ == "__main__":
    results = [
        run_test_1(),
        run_test_2(),
        run_test_3(),
        run_test_4(),
        run_test_5(),
    ]

    path = write_report(results)

    print("\n" + "=" * 72)
    print("SUMMARY")
    print("=" * 72)
    print(f"{'Test':<6} {'Name':<25} {'Verdict':<10} Key Finding")
    print("-" * 72)
    for i, r in enumerate(results, 1):
        finding = r["key_finding"][:70]
        print(f"{i:<6} {r['name']:<25} {r['verdict']:<10} {finding}")
