"""
Experiment 1: Human Agreement

Goal: Does ASF think like a security architect?
Procedure: Compare ASF assumption output against ground truth.
Target: 70%+ agreement (recall/precision).
"""
from __future__ import annotations
from benchmark.data import GroundTruth, ExperimentResult, BenchmarkResult
from benchmark.runner import score_benchmark


def run(ground_truth: GroundTruth, benchmark: BenchmarkResult) -> ExperimentResult:
    scores = score_benchmark(ground_truth, benchmark.asf_assumptions, match_threshold=0.35)

    recall = scores["recall"]
    precision = scores["precision"]
    f1 = scores["f1"]

    findings = []
    recommendations = []

    status = "PASS" if recall >= 0.50 else "FAIL"

    if recall >= 0.70:
        findings.append("ASF recall meets 70%+ agreement target with human-generated ground truth.")
    elif recall >= 0.50:
        findings.append(f"ASF recall is {recall:.1%} — approaching but below the 70% target.")
        recommendations.append("Review assumption extraction patterns to improve coverage of implicit and derived assumptions.")
    else:
        findings.append(f"ASF recall is {recall:.1%} — well below the 70% target.")
        recommendations.append("Major expansion needed in assumption extraction engine to capture implicit security assumptions.")

    if precision >= 0.70:
        findings.append("ASF precision is strong — extracted assumptions are relevant and on-topic.")
    elif precision >= 0.50:
        findings.append(f"ASF precision is {precision:.1%} — some irrelevant assumptions detected.")
        recommendations.append("Add filtering to reduce false positive assumptions.")
    else:
        findings.append(f"ASF precision is {precision:.1%} — high false positive rate.")
        recommendations.append("Implement stricter classification and filtering.")

    findings.append(f"F1 Score: {f1:.1%} — harmonic mean of precision and recall.")

    if recall < 0.70 or precision < 0.70:
        recommendations.append("Run Experiment 2 (Novel Discovery) to determine whether missed ground truth assumptions are actually false negatives or legitimate novel discoveries.")

    return ExperimentResult(
        name="Human Agreement",
        status=status,
        metrics={
            "recall": round(recall, 4),
            "precision": round(precision, 4),
            "f1": round(f1, 4),
            "gt_covered": scores["gt_covered"],
            "asf_relevant": scores["asf_relevant"],
            "asf_irrelevant": scores["asf_irrelevant"],
            "match_threshold": scores["match_threshold"],
            "total_ground_truth": scores["total_ground_truth"],
            "total_asf_output": scores["total_asf_output"],
        },
        findings=findings,
        recommendations=recommendations,
    )
