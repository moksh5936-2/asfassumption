"""ASF Benchmark CLI — run the full evaluation program."""
from __future__ import annotations
import json
import sys
from pathlib import Path
from datetime import datetime, timezone
from typing import Optional

from benchmark.data import BenchmarkResult, ExperimentResult
from benchmark.ground_truth import build_ground_truth
from benchmark.runner import run_benchmark, print_benchmark_scores, print_coverage_report, score_benchmark
from benchmark.experiments.experiment_1 import run as exp1_run
from benchmark.experiments.experiment_2 import run as exp2_run
from benchmark.experiments.experiment_3 import run as exp3_run
from benchmark.experiments.experiment_4 import run as exp4_run


def run_all(parallel: bool = True, output_dir: Optional[str] = None) -> BenchmarkResult:
    print("=" * 72)
    print("  ASF BENCHMARK v1 — Evaluation Program")
    print("=" * 72)

    # Build ground truth
    print("\n  Building ground truth dataset...")
    gt = build_ground_truth()
    print(f"    Policies:  {gt.total_policies}")
    print(f"  Assumptions: {gt.total_assumptions}")
    print(f"  By type:     {gt.by_type()}")
    print(f"  By category: {gt.by_category()}")

    # Run ASF against ground truth
    print("\n  Running ASF against {0} policies...".format(gt.total_policies))
    print("  (This may take a moment)")
    benchmark = run_benchmark(gt, parallel=parallel)
    scores = score_benchmark(gt, benchmark.asf_assumptions)
    print_benchmark_scores(scores)
    print_coverage_report(gt, scores)

    # Run experiments
    print("\n" + "─" * 72)
    print("  EXPERIMENT 1: Human Agreement")
    print("─" * 72)
    e1 = exp1_run(gt, benchmark)
    _print_experiment(e1)

    print("\n" + "─" * 72)
    print("  EXPERIMENT 2: Novel Discovery")
    print("─" * 72)
    e2 = exp2_run(gt, benchmark)
    _print_experiment(e2)

    print("\n" + "─" * 72)
    print("  EXPERIMENT 3: Differentiation")
    print("─" * 72)
    e3 = exp3_run()
    _print_experiment(e3)

    print("\n" + "─" * 72)
    print("  EXPERIMENT 4: Executive Value")
    print("─" * 72)
    e4 = exp4_run()
    _print_experiment(e4)

    benchmark.experiment_results = [e1, e2, e3, e4]

    # Save results
    if output_dir:
        save_results(benchmark, Path(output_dir))
        print(f"\n  Results saved to {output_dir}/")

    # Summary
    print("\n" + "=" * 72)
    print("  BENCHMARK SUMMARY")
    print("=" * 72)
    summary = benchmark.summary()
    print(f"  Policies evaluated:          {summary['total_policies']}")
    print(f"  Ground truth assumptions:   {summary['total_ground_truth']}")
    print(f"  ASF assumptions generated:  {summary['total_asf_assumptions']}")
    print(f"  ASF Precision:              {scores['precision']:.1%}")
    print(f"  ASF Recall:                 {scores['recall']:.1%}")
    print(f"  ASF F1 Score:               {scores['f1']:.1%}")
    passed = summary['experiments_passed']
    total_exp = summary['experiments_run']
    print(f"  Experiments:                {passed}/{total_exp} PASS")
    for er in benchmark.experiment_results:
        print(f"    {er.name:30s} {er.status}")
    print()

    return benchmark


def _print_experiment(e: ExperimentResult) -> None:
    print(f"  Status: {e.status}")
    for f in e.findings[:8]:
        print(f"    • {f}")
    if e.recommendations:
        for r in e.recommendations:
            print(f"    ⚑ {r}")
    print()


def save_results(benchmark: BenchmarkResult, output_dir: Path) -> None:
    output_dir.mkdir(parents=True, exist_ok=True)

    gt_path = output_dir / "ground_truth.json"
    with open(gt_path, "w") as f:
        json.dump(benchmark.ground_truth.to_dict(), f, indent=2)

    asf_path = output_dir / "asf_assumptions.json"
    with open(asf_path, "w") as f:
        json.dump(benchmark.asf_assumptions, f, indent=2, default=str)

    experiments_path = output_dir / "experiments.json"
    with open(experiments_path, "w") as f:
        json.dump([e.to_dict() for e in benchmark.experiment_results], f, indent=2)

    summary_path = output_dir / "summary.json"
    with open(summary_path, "w") as f:
        summary = benchmark.summary()
        summary["scores"] = {
            "precision": None,
            "recall": None,
            "f1": None,
        }
        json.dump(summary, f, indent=2)


if __name__ == "__main__":
    output = sys.argv[1] if len(sys.argv) > 1 else None
    run_all(parallel=False, output_dir=output)
