from datetime import datetime, timezone, timedelta

from asf.confidence import ConfidenceEngine
from asf.models import Evidence, SourceType, Verification, VerificationResult


class TestConfidenceEngine:
    def setup_method(self):
        self.engine = ConfidenceEngine()

    def test_high_confidence_for_verified(self):
        v = Verification(
            assumption_id="asm_1",
            result=VerificationResult.VERIFIED,
            confidence=0.9,
            evidence_used=["evd_1", "evd_2"],
        )
        ev = [
            Evidence(source="a.csv", source_type=SourceType.CSV, timestamp=datetime.now(timezone.utc)),
            Evidence(source="b.csv", source_type=SourceType.CSV, timestamp=datetime.now(timezone.utc)),
        ]
        conf = self.engine.compute_verification_confidence(v, ev)
        assert conf > 0.5

    def test_lower_confidence_for_unknown(self):
        v = Verification(
            assumption_id="asm_2",
            result=VerificationResult.UNKNOWN,
            confidence=0.0,
        )
        conf = self.engine.compute_verification_confidence(v, [])
        assert conf == 0.0

    def test_freshness_decay(self):
        old = datetime.now(timezone.utc) - timedelta(days=60)
        ev = [
            Evidence(source="old.csv", source_type=SourceType.CSV, timestamp=old),
        ]
        v = Verification(
            assumption_id="asm_3",
            result=VerificationResult.VERIFIED,
            confidence=0.9,
            evidence_used=["evd_1"],
        )
        conf = self.engine.compute_verification_confidence(v, ev)
        assert conf < 0.9

    def test_assumption_confidence(self):
        vs = [
            Verification(assumption_id="asm_1", result=VerificationResult.VERIFIED, confidence=0.9),
            Verification(assumption_id="asm_1", result=VerificationResult.CONTRADICTED, confidence=0.8),
        ]
        conf = self.engine.compute_assumption_confidence(vs)
        assert 0.0 < conf < 1.0
