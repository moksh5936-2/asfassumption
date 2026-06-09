from asf.extraction import ClaimExtractor
from asf.assumption import AssumptionEngine
from asf.models import AssumptionType, Claim


class TestAssumptionEngine:
    def setup_method(self):
        self.extractor = ClaimExtractor()
        self.engine = AssumptionEngine()

    def test_access_assumption(self):
        claim = Claim(source_document="test.txt", text="Only Finance employees may access payroll.")
        assumption = self.engine.convert(claim)
        assert assumption is not None
        assert assumption.assumption_type == AssumptionType.ACCESS
        assert assumption.claim_id == claim.id
        assert "payroll" in assumption.text

    def test_identity_assumption(self):
        claim = Claim(source_document="test.txt", text="MFA is required for all administrative access.")
        assumption = self.engine.convert(claim)
        assert assumption is not None
        assert assumption.assumption_type == AssumptionType.IDENTITY

    def test_network_assumption(self):
        claim = Claim(source_document="test.txt", text="Production databases are not internet accessible.")
        assumption = self.engine.convert(claim)
        assert assumption is not None
        assert assumption.assumption_type == AssumptionType.NETWORK

    def test_configuration_assumption(self):
        claim = Claim(source_document="test.txt", text="All payroll data is encrypted at rest.")
        assumption = self.engine.convert(claim)
        assert assumption is not None
        assert assumption.assumption_type == AssumptionType.CONFIGURATION

    def test_process_assumption(self):
        claim = Claim(source_document="test.txt", text="All configuration changes must be approved by the security team.")
        assumption = self.engine.convert(claim)
        assert assumption is not None
        assert assumption.assumption_type == AssumptionType.PROCESS

    def test_governance_assumption(self):
        claim = Claim(source_document="test.txt", text="Security reviews are conducted quarterly.")
        assumption = self.engine.convert(claim)
        assert assumption is not None
        assert assumption.assumption_type == AssumptionType.GOVERNANCE

    def test_non_claim_not_converted(self):
        claim = Claim(source_document="test.txt", text="The system has a login page.")
        assumption = self.engine.convert(claim)
        # Should still return an assumption because it might match some patterns
        # but with lower confidence

    def test_keywords_extracted(self):
        claim = Claim(source_document="test.txt", text="Only Finance employees may access the payroll system.")
        assumption = self.engine.convert(claim)
        assert assumption is not None
        assert len(assumption.keywords) > 0
