from __future__ import annotations
from dataclasses import dataclass, field
from typing import Optional
from datetime import datetime, timezone
from asf.models import AssumptionType


@dataclass
class Policy:
    id: str
    text: str
    domain: str
    tags: list[str] = field(default_factory=list)
    source_type: str = "policy"

    def to_dict(self) -> dict:
        return {"id": self.id, "text": self.text, "domain": self.domain, "tags": self.tags}


@dataclass
class GroundTruthAssumption:
    id: str
    policy_id: str
    text: str
    type: AssumptionType
    category: str  # "explicit", "implicit", "derived"
    is_critical: bool = False
    commentary: str = ""
    keywords: list[str] = field(default_factory=list)

    def to_dict(self) -> dict:
        return {
            "id": self.id,
            "policy_id": self.policy_id,
            "text": self.text,
            "type": str(self.type),
            "category": self.category,
            "is_critical": self.is_critical,
            "commentary": self.commentary,
            "keywords": self.keywords,
        }


@dataclass
class GroundTruth:
    policies: list[Policy] = field(default_factory=list)
    assumptions: list[GroundTruthAssumption] = field(default_factory=list)

    @property
    def total_assumptions(self) -> int:
        return len(self.assumptions)

    @property
    def total_policies(self) -> int:
        return len(self.policies)

    def get_assumptions_for(self, policy_id: str) -> list[GroundTruthAssumption]:
        return [a for a in self.assumptions if a.policy_id == policy_id]

    def by_type(self) -> dict[str, int]:
        counts: dict[str, int] = {}
        for a in self.assumptions:
            t = str(a.type)
            counts[t] = counts.get(t, 0) + 1
        return counts

    def by_category(self) -> dict[str, int]:
        counts: dict[str, int] = {}
        for a in self.assumptions:
            counts[a.category] = counts.get(a.category, 0) + 1
        return counts

    def by_domain(self) -> dict[str, int]:
        counts: dict[str, int] = {}
        for p in self.policies:
            counts[p.domain] = counts.get(p.domain, 0) + 1
        return counts

    def to_dict(self) -> dict:
        return {
            "policies": [p.to_dict() for p in self.policies],
            "assumptions": [a.to_dict() for a in self.assumptions],
            "summary": {
                "total_policies": self.total_policies,
                "total_assumptions": self.total_assumptions,
                "by_type": self.by_type(),
                "by_category": self.by_category(),
                "by_domain": self.by_domain(),
            },
        }


@dataclass
class ExperimentResult:
    name: str
    status: str  # PASS, FAIL, INCONCLUSIVE
    metrics: dict = field(default_factory=dict)
    findings: list[str] = field(default_factory=list)
    recommendations: list[str] = field(default_factory=list)
    raw_data: dict = field(default_factory=dict)

    def to_dict(self) -> dict:
        return {
            "name": self.name,
            "status": self.status,
            "metrics": self.metrics,
            "findings": self.findings,
            "recommendations": self.recommendations,
        }


@dataclass
class BenchmarkResult:
    ground_truth: GroundTruth
    asf_assumptions: dict[str, list[dict]]  # policy_id -> list of assumption dicts
    experiment_results: list[ExperimentResult] = field(default_factory=list)
    created_at: str = field(default_factory=lambda: datetime.now(timezone.utc).isoformat())

    def summary(self) -> dict:
        asf_total = sum(len(v) for v in self.asf_assumptions.values())
        return {
            "total_policies": self.ground_truth.total_policies,
            "total_ground_truth": self.ground_truth.total_assumptions,
            "total_asf_assumptions": asf_total,
            "experiments_run": len(self.experiment_results),
            "experiments_passed": sum(1 for e in self.experiment_results if e.status == "PASS"),
        }
