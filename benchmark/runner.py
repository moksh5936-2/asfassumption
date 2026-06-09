"""Core benchmark runner — evaluates ASF against ground truth."""
from __future__ import annotations
import re
import json
from pathlib import Path
from typing import Any
from collections import Counter
from concurrent.futures import ProcessPoolExecutor, as_completed

from asf.analyzer import Analyzer
from benchmark.data import GroundTruth, Policy, GroundTruthAssumption, BenchmarkResult


def normalize(text: str) -> str:
    return re.sub(r"[^a-z0-9\s]", "", text.lower()).strip()


def compute_match(gt_assumption: GroundTruthAssumption, asf_assumptions: list[dict]) -> float:
    """Compute overlap score between a ground truth assumption and ASF output."""
    gt_text = normalize(gt_assumption.text)
    gt_keywords = set(k.lower() for k in gt_assumption.keywords)
    gt_type = str(gt_assumption.type)

    best = 0.0
    for asm in asf_assumptions:
        asm_text = normalize(asm.get("text", ""))
        asm_type = asm.get("assumption_type", "")
        asm_keywords = set(k.lower() for k in asm.get("keywords", []))

        text_match = _text_overlap(gt_text, asm_text)
        type_match = 1.0 if asm_type == gt_type else 0.3
        keyword_match = _keyword_overlap(gt_keywords, asm_keywords) if gt_keywords else 0.5

        score = 0.5 * text_match + 0.25 * type_match + 0.25 * keyword_match
        best = max(best, score)

    return best


def _text_overlap(a: str, b: str) -> float:
    a_tokens = set(a.split())
    b_tokens = set(b.split())
    if not a_tokens or not b_tokens:
        return 0.0
    intersection = a_tokens & b_tokens
    return len(intersection) / max(len(a_tokens), len(b_tokens))


def _keyword_overlap(gt_kw: set[str], asm_kw: set[str]) -> float:
    if not gt_kw:
        return 0.0
    intersection = gt_kw & asm_kw
    return len(intersection) / len(gt_kw)


def score_benchmark(
    ground_truth: GroundTruth,
    asf_results: dict[str, list[dict]],
    match_threshold: float = 0.35,
) -> dict[str, Any]:
    """Compute precision, recall, F1 across all policies.
    
    TP = ASF assumptions that match >=1 GT (relevant items retrieved)
    FP = ASF assumptions that match NO GT (irrelevant items retrieved)
    FN = GT assumptions that match NO ASF (relevant items not retrieved)
    """
    gt_by_policy: dict[str, list[GroundTruthAssumption]] = {}
    for a in ground_truth.assumptions:
        gt_by_policy.setdefault(a.policy_id, []).append(a)

    asf_total = 0
    gt_match_count = 0  # GT assumptions that had at least one ASF match
    per_policy: dict[str, dict] = {}

    for policy in ground_truth.policies:
        pid = policy.id
        gt_list = gt_by_policy.get(pid, [])
        asm_list = asf_results.get(pid, [])

        asf_total += len(asm_list)

        # How many GT assumptions does ASF cover
        gt_covered = 0
        for gt in gt_list:
            score = compute_match(gt, asm_list)
            if score >= match_threshold:
                gt_covered += 1
        gt_match_count += gt_covered

        # How many ASF assumptions are relevant (match >=1 GT)
        asf_relevant = 0
        for asm in asm_list:
            for gt in gt_list:
                reverse_score = compute_match(gt, [asm])
                if reverse_score >= match_threshold:
                    asf_relevant += 1
                    break

        asf_irrelevant = len(asm_list) - asf_relevant
        fn = len(gt_list) - gt_covered

        per_policy[pid] = {
            "ground_truth": len(gt_list),
            "asf_found": len(asm_list),
            "gt_covered": gt_covered,
            "asf_relevant": asf_relevant,
            "asf_irrelevant": asf_irrelevant,
            "false_negatives": fn,
            "recall": gt_covered / len(gt_list) if gt_list else 0,
            "precision": asf_relevant / len(asm_list) if asm_list else 0,
        }

    total_gt = len(ground_truth.assumptions)
    total_asf_relevant = sum(p["asf_relevant"] for p in per_policy.values())
    total_asf_irrelevant = sum(p["asf_irrelevant"] for p in per_policy.values())
    total_fn = sum(p["false_negatives"] for p in per_policy.values())

    precision = total_asf_relevant / (total_asf_relevant + total_asf_irrelevant) if (total_asf_relevant + total_asf_irrelevant) > 0 else 0
    recall = gt_match_count / total_gt if total_gt > 0 else 0
    f1 = 2 * precision * recall / (precision + recall) if (precision + recall) > 0 else 0

    return {
        "gt_covered": gt_match_count,
        "asf_relevant": total_asf_relevant,
        "asf_irrelevant": total_asf_irrelevant,
        "false_negatives": total_fn,
        "precision": round(precision, 4),
        "recall": round(recall, 4),
        "f1": round(f1, 4),
        "total_ground_truth": len(ground_truth.assumptions),
        "total_asf_output": asf_total,
        "match_threshold": match_threshold,
        "per_policy": per_policy,
    }


def _run_asf_on_policy(policy: Policy) -> tuple[str, list[dict]]:
    try:
        analyzer = Analyzer()
        result = analyzer.analyze(document_paths=[])
        from asf.extraction.claim_extractor import ClaimExtractor
        extractor = ClaimExtractor()
        claims = extractor.extract(policy.text, source_document=policy.id)

        from asf.assumption.assumption_engine import AssumptionEngine
        engine = AssumptionEngine()
        assumptions = engine.convert_many(claims)

        asm_list = []
        for a in assumptions:
            asm_list.append({
                "id": a.id,
                "text": a.text,
                "assumption_type": str(a.assumption_type),
                "claim_id": a.claim_id,
                "confidence": a.confidence,
                "keywords": list(a.keywords),
            })
        return policy.id, asm_list
    except Exception as e:
        return policy.id, [{"text": f"ERROR: {e}", "assumption_type": "UNKNOWN", "keywords": []}]


def run_benchmark(ground_truth: GroundTruth, parallel: bool = True, max_workers: int = 8) -> BenchmarkResult:
    asf_results: dict[str, list[dict]] = {}

    if parallel and len(ground_truth.policies) > 1:
        with ProcessPoolExecutor(max_workers=max_workers) as executor:
            futures = {executor.submit(_run_asf_on_policy, p): p for p in ground_truth.policies}
            for future in as_completed(futures):
                pid, asm_list = future.result()
                asf_results[pid] = asm_list
    else:
        for p in ground_truth.policies:
            pid, asm_list = _run_asf_on_policy(p)
            asf_results[pid] = asm_list

    scores = score_benchmark(ground_truth, asf_results)

    return BenchmarkResult(
        ground_truth=ground_truth,
        asf_assumptions=asf_results,
        experiment_results=[],
    )


def print_benchmark_scores(scores: dict[str, Any]) -> None:
    print(f"\n  Benchmark Scores:")
    print(f"  {'Precision:':20s} {scores['precision']:.1%}")
    print(f"  {'Recall:':20s} {scores['recall']:.1%}")
    print(f"  {'F1 Score:':20s} {scores['f1']:.1%}")
    print(f"  {'GT Covered:':20s} {scores['gt_covered']} / {scores['total_ground_truth']}")
    print(f"  {'ASF Relevant:':20s} {scores['asf_relevant']}")
    print(f"  {'ASF Irrelevant:':20s} {scores['asf_irrelevant']}")
    print(f"  {'False Negatives:':20s} {scores['false_negatives']}")
    print(f"  {'Ground Truth Total:':20s} {scores['total_ground_truth']}")
    print(f"  {'ASF Output Total:':20s} {scores['total_asf_output']}")
    print(f"  {'Match Threshold:':20s} {scores['match_threshold']}")


def print_coverage_report(ground_truth: GroundTruth, scores: dict[str, Any]) -> None:
    print(f"\n  Coverage by Type:")
    by_type = ground_truth.by_type()
    for atype in sorted(by_type.keys()):
        count = by_type[atype]
        print(f"    {atype:20s} {count}")

    print(f"\n  Coverage by Category:")
    by_cat = ground_truth.by_category()
    for cat in sorted(by_cat.keys()):
        print(f"    {cat:15s} {by_cat[cat]}")
