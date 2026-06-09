"""CLI entry point for ASF Benchmark."""
from __future__ import annotations
import sys
import click
from benchmark.main import run_all


@click.group()
def cli():
    """ASF Benchmark v1 — Evaluation Program"""


@cli.command()
@click.option("--output", "-o", default="benchmark_report", help="Output directory for results")
@click.option("--sequential", is_flag=True, help="Disable parallel processing")
def run(output: str, sequential: bool):
    """Run the full ASF Benchmark evaluation suite."""
    run_all(parallel=not sequential, output_dir=output)


@cli.command()
def info():
    """Show benchmark dataset summary."""
    from benchmark.ground_truth import build_ground_truth
    gt = build_ground_truth()
    click.echo(f"Policies:   {gt.total_policies}")
    click.echo(f"Assumptions: {gt.total_assumptions}")
    click.echo(f"By type:     {gt.by_type()}")
    click.echo(f"By category: {gt.by_category()}")
    click.echo(f"By domain:   {gt.by_domain()}")


if __name__ == "__main__":
    cli()
