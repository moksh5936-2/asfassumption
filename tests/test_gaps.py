from asf.gaps import GapEngine
from asf.models import Assumption, AssumptionType, GapSeverity, GapType, Verification, VerificationResult


class TestGapEngine:
    def setup_method(self):
        self.engine = GapEngine()

    def test_gap_for_contradicted(self):
        assumption = Assumption(
            claim_id="clm_1",
            text="Only Finance can access payroll.",
            assumption_type=AssumptionType.ACCESS,
        )
        verification = Verification(
            assumption_id=assumption.id,
            result=VerificationResult.CONTRADICTED,
            confidence=0.92,
        )
        gaps = self.engine.generate_gaps([assumption], [verification])
        assert len(gaps) == 1
        assert gaps[0].type == GapType.ACCESS_GAP
        assert gaps[0].severity == GapSeverity.CRITICAL

    def test_gap_for_unverified(self):
        assumption = Assumption(
            claim_id="clm_2",
            text="Some assumption",
            assumption_type=AssumptionType.NETWORK,
        )
        gaps = self.engine.generate_gaps([assumption], [])
        assert len(gaps) == 1
        assert gaps[0].type == GapType.VERIFICATION_GAP
        assert gaps[0].severity == GapSeverity.MEDIUM

    def test_no_gap_for_verified(self):
        assumption = Assumption(
            claim_id="clm_3",
            text="Verified assumption.",
            assumption_type=AssumptionType.CONFIGURATION,
        )
        verification = Verification(
            assumption_id=assumption.id,
            result=VerificationResult.VERIFIED,
            confidence=0.9,
        )
        gaps = self.engine.generate_gaps([assumption], [verification])
        assert len(gaps) == 0

    def test_gap_for_partially_verified(self):
        assumption = Assumption(
            claim_id="clm_4",
            text="Partial assumption.",
            assumption_type=AssumptionType.IDENTITY,
        )
        verification = Verification(
            assumption_id=assumption.id,
            result=VerificationResult.PARTIALLY_VERIFIED,
            confidence=0.5,
        )
        gaps = self.engine.generate_gaps([assumption], [verification])
        assert len(gaps) == 1
        assert gaps[0].severity == GapSeverity.MEDIUM

    def test_gap_for_unknown(self):
        assumption = Assumption(
            claim_id="clm_5",
            text="Unknown assumption.",
            assumption_type=AssumptionType.DOCUMENTATION,
        )
        verification = Verification(
            assumption_id=assumption.id,
            result=VerificationResult.UNKNOWN,
            confidence=0.0,
        )
        gaps = self.engine.generate_gaps([assumption], [verification])
        assert len(gaps) == 1
        assert gaps[0].type == GapType.EVIDENCE_GAP
