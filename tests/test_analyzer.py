from pathlib import Path

from asf.analyzer import Analyzer

SAMPLE_DIR = Path(__file__).parent.parent / "sample_data"


class TestAnalyzer:
    def setup_method(self):
        self.analyzer = Analyzer()

    def test_end_to_end_with_sample_data(self):
        docs = [SAMPLE_DIR / "finance_policy.txt"]
        ev = [
            SAMPLE_DIR / "payroll_acl.csv",
            SAMPLE_DIR / "network_exposure.csv",
            SAMPLE_DIR / "mfa_status.csv",
            SAMPLE_DIR / "backup_config.csv",
        ]

        result = self.analyzer.analyze(docs, ev)

        assert len(result.claims) > 0
        assert len(result.assumptions) > 0
        assert len(result.verifications) > 0

        total_verified = sum(1 for v in result.verifications if v.result.value == "VERIFIED")
        total_contradicted = sum(1 for v in result.verifications if v.result.value == "CONTRADICTED")

        # We expect at least some verifications
        assert total_verified + total_contradicted > 0

        # Graph should be built
        graph_summary = self.analyzer.graph_model.summary()
        assert graph_summary["total_nodes"] > 0
        assert graph_summary["total_edges"] > 0

    def test_analyzer_no_evidence(self):
        docs = [SAMPLE_DIR / "finance_policy.txt"]
        result = self.analyzer.analyze(docs)
        assert len(result.claims) > 0
        assert len(result.assumptions) > 0
        # Without evidence, many should be unknown
        unknowns = sum(1 for v in result.verifications if v.result.value == "UNKNOWN")
        assert unknowns > 0
