from __future__ import annotations

import json as _json
from pathlib import Path
from typing import Optional

import click
from rich.console import Console
from rich.panel import Panel
from rich.table import Table
from rich import box

from asf import __version__
from asf.analyzer import Analyzer
from asf.config import ASFConfig
from asf.models import AnalysisResult, GapSeverity, VerificationResult
from asf.settings import load_config, write_default_config

console = Console()


@click.group()
@click.option("--config", "-c", type=click.Path(path_type=Path), help="Path to config file")
@click.pass_context
def cli(ctx, config):
    """ASF Validator v0.1 — Assumption Security Framework Validator

    Experimental research platform for testing security assumptions against evidence.
    """
    cfg = load_config(config) if config else ASFConfig.default()
    ctx.ensure_object(dict)
    ctx.obj["config"] = cfg


@cli.command()
@click.argument("paths", nargs=-1, type=click.Path(exists=True, path_type=Path), required=False)
@click.option("-e", "--evidence", multiple=True, type=click.Path(exists=True, path_type=Path), help="Evidence files")
@click.option("--json", "json_output", is_flag=True, help="Output as JSON")
@click.option("--graph", "graph_output", is_flag=True, help="Export graph as JSON")
@click.option("--persist", is_flag=True, help="Persist results to database")
@click.option("--auto-map", is_flag=True, help="Auto-map evidence column names")
@click.pass_context
def analyze(ctx, paths, evidence, json_output, graph_output, persist, auto_map):
    """Analyze documents and evidence for security assumptions.

    PATHS: Policy files, architecture docs, runbooks (PDF, DOCX, TXT).

    If PATHS is a directory, all supported files within are processed.

    EVIDENCE: IAM exports, ACL files, configuration exports (CSV, JSON).
    """
    config: ASFConfig = ctx.obj["config"]

    doc_list = _resolve_paths(paths, config)
    if not doc_list:
        click.echo("Error: no document paths provided and none found in config", err=True)
        raise click.Abort()

    ev_list = _resolve_evidence(evidence, config)

    analyzer = Analyzer(config)

    with console.status("[bold green]Running ASF analysis..."):
        result = analyzer.analyze(doc_list, ev_list, persist=persist)

    if json_output:
        output = _json.dumps(_serialize_result(result), indent=2, default=str, ensure_ascii=False)
        click.echo(output)
        return

    if graph_output:
        click.echo(analyzer.graph_model.export_json_string())
        return

    _print_summary(result)
    _print_findings(result)
    _print_gaps(result)


@cli.command()
def init():
    """Create a default asf.config.yaml in the current directory."""
    path = write_default_config()
    click.echo(f"Created default config: {path}")
    click.echo("Edit this file to customize evidence schemas, patterns, and LLM settings.")


def _resolve_paths(paths: tuple[Path, ...], config: ASFConfig) -> list[Path]:
    result: list[Path] = []
    supported = {".pdf", ".docx", ".txt", ".csv", ".json"}

    for p in paths:
        if p.is_dir():
            for ext in supported:
                result.extend(sorted(p.glob(f"*{ext}")))
        elif p.suffix.lower() in supported:
            result.append(p)

    return result


def _resolve_evidence(evidence: tuple[Path, ...], config: ASFConfig) -> list[Path]:
    result = list(evidence)
    for d in config.evidence_dirs:
        dpath = Path(d)
        if dpath.is_dir():
            for ext in {".csv", ".json"}:
                result.extend(sorted(dpath.glob(f"*{ext}")))
    return result


def _print_summary(result: AnalysisResult) -> None:
    panel = Panel(
        f"[bold]Documents Processed:[/] {result.claims_found} claims extracted\n"
        f"[bold]Assumptions:[/] {result.assumptions_found}\n"
        f"[green]Verified:[/] {result.verified_count}\n"
        f"[red]Contradicted:[/] {result.contradicted_count}\n"
        f"[yellow]Unknown/Partial:[/] {result.unknown_count}\n"
        f"[bold red]Critical Gaps:[/] {result.critical_gaps}",
        title="ASF Analysis Summary",
        box=box.ROUNDED,
    )
    console.print(panel)


def _print_findings(result: AnalysisResult) -> None:
    if not result.verifications:
        return

    table = Table(title="Findings", box=box.SIMPLE)
    table.add_column("Finding", style="cyan", no_wrap=False)
    table.add_column("Status", style="bold")
    table.add_column("Confidence", justify="right")
    table.add_column("Evidence")
    table.add_column("Explanation", no_wrap=False)

    for v in result.verifications:
        assumption = next(
            (a for a in result.assumptions if a.id == v.assumption_id), None
        )
        text = (
            assumption.text[:60] + "..."
            if assumption and len(assumption.text) > 60
            else (assumption.text if assumption else "N/A")
        )

        status_style = {
            VerificationResult.VERIFIED: "green",
            VerificationResult.PARTIALLY_VERIFIED: "yellow",
            VerificationResult.CONTRADICTED: "red",
            VerificationResult.UNKNOWN: "dim",
        }.get(v.result, "white")

        explanation = _build_explanation(v, result)

        table.add_row(
            text,
            f"[{status_style}]{v.result.value}[/]",
            f"{v.confidence:.0%}",
            str(len(v.evidence_used)),
            explanation[:80] + "..." if len(explanation) > 80 else explanation,
        )

    console.print(table)


def _build_explanation(v, result):
    for gap in result.gaps:
        if gap.assumption_id == v.assumption_id:
            return gap.description
    return v.reasoning[:80] if v.reasoning else "No explanation available"


def _print_gaps(result: AnalysisResult) -> None:
    if not result.gaps:
        return

    table = Table(title="Gap Analysis", box=box.SIMPLE)
    table.add_column("Type", style="bold")
    table.add_column("Severity")
    table.add_column("Description", no_wrap=False)

    sev_style = {
        GapSeverity.CRITICAL: "red",
        GapSeverity.HIGH: "orange1",
        GapSeverity.MEDIUM: "yellow",
        GapSeverity.LOW: "dim",
        GapSeverity.INFO: "blue",
    }

    for gap in result.gaps:
        style = sev_style.get(gap.severity, "white")
        table.add_row(
            f"[{style}]{gap.type.value}[/]",
            f"[{style}]{gap.severity.value}[/]",
            gap.description[:80] + "..." if len(gap.description) > 80 else gap.description,
        )

    console.print(table)


def _serialize_result(result: AnalysisResult) -> dict:
    return {
        "summary": {
            "claims_found": result.claims_found,
            "assumptions": result.assumptions_found,
            "verified": result.verified_count,
            "contradicted": result.contradicted_count,
            "unknown": result.unknown_count,
            "critical_gaps": result.critical_gaps,
        },
        "claims": [c.model_dump() for c in result.claims],
        "assumptions": [a.model_dump() for a in result.assumptions],
        "verifications": [
            {
                "assumption_id": v.assumption_id,
                "result": v.result.value,
                "confidence": v.confidence,
                "evidence_used": v.evidence_used,
                "reasoning": v.reasoning,
            }
            for v in result.verifications
        ],
        "gaps": [
            {
                "assumption_id": g.assumption_id,
                "type": g.type.value,
                "severity": g.severity.value,
                "description": g.description,
                "evidence_detail": g.evidence_detail,
            }
            for g in result.gaps
        ],
    }


if __name__ == "__main__":
    cli()
