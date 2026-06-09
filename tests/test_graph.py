from asf.graph import GraphModel
from asf.models import (
    AnalysisResult,
    Assumption,
    AssumptionType,
    Claim,
    Evidence,
    Gap,
    GapSeverity,
    GapType,
    SourceType,
    Verification,
    VerificationResult,
)


class TestGraphModel:
    def setup_method(self):
        self.graph = GraphModel()

    def test_build_graph(self):
        claim = Claim(source_document="test.txt", text="Only Finance can access payroll.")
        assumption = Assumption(claim_id=claim.id, text="System assumes: Only Finance can access payroll.", assumption_type=AssumptionType.ACCESS)
        evidence = Evidence(source="test.csv", source_type=SourceType.CSV, records=[{"user": "a", "group": "Finance"}])
        verification = Verification(assumption_id=assumption.id, result=VerificationResult.VERIFIED, evidence_used=[evidence.id])
        gap = Gap(assumption_id=assumption.id, severity=GapSeverity.HIGH, type=GapType.ACCESS_GAP, description="test")

        result = AnalysisResult(
            claims=[claim],
            assumptions=[assumption],
            evidence=[evidence],
            verifications=[verification],
            gaps=[gap],
        )

        self.graph.build(result)
        summary = self.graph.summary()
        assert summary["total_nodes"] >= 5
        assert summary["total_edges"] >= 4

    def test_export_json(self):
        claim = Claim(source_document="test.txt", text="All data is encrypted.")
        assumption = Assumption(claim_id=claim.id, text="System assumes: All data is encrypted.", assumption_type=AssumptionType.CONFIGURATION)
        result = AnalysisResult(claims=[claim], assumptions=[assumption])
        self.graph.build(result)
        exported = self.graph.export_json()
        assert "nodes" in exported
        assert "edges" in exported
        assert exported["node_count"] >= 2
