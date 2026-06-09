from asf.models import (
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
    VerificationStatus,
    AnalysisResult,
)


class TestModels:
    def test_claim_creation(self):
        c = Claim(source_document="test.pdf", text="Only Finance can access payroll.")
        assert c.id.startswith("clm_")
        assert c.source_document == "test.pdf"
        assert c.extraction_confidence == 0.5

    def test_assumption_creation(self):
        a = Assumption(claim_id="clm_123", text="System assumes...", assumption_type=AssumptionType.ACCESS)
        assert a.id.startswith("asm_")
        assert a.verification_status == VerificationStatus.PENDING
        assert a.confidence == 0.0

    def test_evidence_creation(self):
        e = Evidence(source="iam.csv", source_type=SourceType.CSV)
        assert e.id.startswith("evd_")
        assert e.confidence == 0.8

    def test_verification_creation(self):
        v = Verification(assumption_id="asm_123")
        assert v.id.startswith("vrf_")
        assert v.result == VerificationResult.UNKNOWN

    def test_gap_creation(self):
        g = Gap(assumption_id="asm_123", severity=GapSeverity.HIGH, type=GapType.ACCESS_GAP, description="test")
        assert g.id.startswith("gap_")

    def test_analysis_result(self):
        r = AnalysisResult()
        assert r.claims_found == 0
        assert r.assumptions_found == 0
        assert r.verified_count == 0
        assert r.contradicted_count == 0
        assert r.unknown_count == 0
        assert r.critical_gaps == 0

    def test_enums(self):
        assert AssumptionType.ACCESS.value == "ACCESS"
        assert VerificationResult.VERIFIED.value == "VERIFIED"
        assert GapSeverity.CRITICAL.value == "CRITICAL"
        assert GapType.ACCESS_GAP.value == "ACCESS_GAP"
        assert SourceType.CSV.value == "CSV"
