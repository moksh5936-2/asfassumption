"""
Experiment 2: Novel Discovery

Goal: Does ASF discover assumptions humans missed?
Procedure: Identify ASF-generated assumptions that have NO corresponding ground truth.
These are candidates for "novel discovery" — assumptions a human didn't list
but ASF identified.
"""
from __future__ import annotations
from typing import Any
from benchmark.data import GroundTruth, BenchmarkResult, ExperimentResult, GroundTruthAssumption
from benchmark.runner import normalize


def run(ground_truth: GroundTruth, benchmark: BenchmarkResult) -> ExperimentResult:
    gt_by_policy: dict[str, list[GroundTruthAssumption]] = {}
    for a in ground_truth.assumptions:
        gt_by_policy.setdefault(a.policy_id, []).append(a)

    novel_total = 0
    novel_by_type: dict[str, int] = {}
    novel_findings: list[dict[str, Any]] = []

    for policy in ground_truth.policies:
        pid = policy.id
        asm_list = benchmark.asf_assumptions.get(pid, [])
        gt_list = gt_by_policy.get(pid, [])

        gt_texts = [normalize(a.text) for a in gt_list]
        gt_keywords = [set(k.lower() for k in a.keywords) for a in gt_list]

        for asm in asm_list:
            asm_text = normalize(asm.get("text", ""))
            asm_type = asm.get("assumption_type", "UNKNOWN")

            # Check if this ASF assumption has any match in ground truth
            best_match = 0.0
            for i, gt_text in enumerate(gt_texts):
                atokens = set(asm_text.split())
                gtokens = set(gt_text.split())
                if not atokens or not gtokens:
                    continue
                overlap = len(atokens & gtokens) / max(len(atokens), len(gtokens))
                best_match = max(best_match, overlap)

            # If the best match is below threshold, it's a novel discovery candidate
            if best_match < 0.20:
                novel_total += 1
                novel_by_type[asm_type] = novel_by_type.get(asm_type, 0) + 1
                if len(novel_findings) < 20:
                    novel_findings.append({
                        "policy_id": pid,
                        "policy_text": policy.text[:100],
                        "asf_text": asm.get("text", ""),
                        "type": asm_type,
                        "best_gt_overlap": round(best_match, 3),
                    })

    total_asf = sum(len(v) for v in benchmark.asf_assumptions.values())
    novel_rate = novel_total / total_asf if total_asf > 0 else 0

    findings = []
    recommendations = []

    if novel_rate >= 0.30:
        findings.append(f"High novel discovery rate: {novel_rate:.1%} of ASF assumptions ({novel_total}/{total_asf}) are not in ground truth.")
        recommendations.append("Review novel assumptions with security practitioners to validate their correctness.")
        recommendations.append("This is the experiment that matters most — schedule review sessions with 3+ security engineers.")
    elif novel_rate >= 0.15:
        findings.append(f"Moderate novel discovery rate: {novel_rate:.1%} of ASF assumptions ({novel_total}/{total_asf}) are not in ground truth.")
        recommendations.append("Sample novel assumptions for human review — some may be valuable, some may be false positives.")
    else:
        findings.append(f"Low novel discovery rate: {novel_rate:.1%} of ASF assumptions ({novel_total}/{total_asf}) are novel.")
        recommendations.append("ASF may be overfitting to ground truth patterns rather than discovering genuinely new assumptions.")

    findings.append(f"Novel discoveries by type: {dict(sorted(novel_by_type.items(), key=lambda x: -x[1]))}")

    if novel_total > 0:
        findings.append(f"Top {min(20, len(novel_findings))} novel discovery candidates identified for human review.")

    status = "PASS" if novel_total >= 50 else "FAIL"

    return ExperimentResult(
        name="Novel Discovery",
        status=status,
        metrics={
            "total_novel": novel_total,
            "total_asf": total_asf,
            "novel_rate": round(novel_rate, 4),
            "novel_by_type": novel_by_type,
            "top_candidates": novel_findings[:10],
        },
        findings=findings,
        recommendations=recommendations,
    )
