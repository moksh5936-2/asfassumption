from asf.models import (
    Assumption,
    AssumptionType,
    Evidence,
    SourceType,
    VerificationResult,
)
from asf.verification import VerificationEngine


class TestVerificationEngine:
    def setup_method(self):
        self.engine = VerificationEngine()

    def test_access_contradicted(self):
        assumption = Assumption(
            claim_id="clm_1",
            text="System assumes access control: Only Finance employees may access the payroll system.",
            assumption_type=AssumptionType.ACCESS,
        )
        evidence = Evidence(
            source="test.csv",
            source_type=SourceType.CSV,
            records=[
                {"user": "alice", "group": "Finance", "permission": "read"},
                {"user": "dave", "group": "Engineering", "permission": "read"},
                {"user": "bob", "group": "Finance", "permission": "write"},
            ],
        )
        ver = self.engine.verify(assumption, [evidence])
        assert ver.result == VerificationResult.CONTRADICTED
        assert ver.confidence > 0.5

    def test_access_verified(self):
        assumption = Assumption(
            claim_id="clm_2",
            text="System assumes access control: Only Finance employees may access the payroll system.",
            assumption_type=AssumptionType.ACCESS,
        )
        evidence = Evidence(
            source="test.csv",
            source_type=SourceType.CSV,
            records=[
                {"user": "alice", "group": "Finance", "permission": "read"},
                {"user": "bob", "group": "Finance", "permission": "write"},
            ],
        )
        ver = self.engine.verify(assumption, [evidence])
        assert ver.result == VerificationResult.VERIFIED

    def test_network_contradicted(self):
        assumption = Assumption(
            claim_id="clm_3",
            text="System assumes network posture: Production databases are not internet accessible.",
            assumption_type=AssumptionType.NETWORK,
        )
        evidence = Evidence(
            source="test.csv",
            source_type=SourceType.CSV,
            records=[
                {"asset": "payroll-db", "public": "false", "internet_facing": "false"},
                {"asset": "customer-portal", "public": "true", "internet_facing": "true"},
            ],
        )
        ver = self.engine.verify(assumption, [evidence])
        assert ver.result == VerificationResult.VERIFIED  # isolated assets found, no exposure

    def test_network_exposed(self):
        assumption = Assumption(
            claim_id="clm_4",
            text="System assumes network posture: No external access is permitted to the finance application.",
            assumption_type=AssumptionType.NETWORK,
        )
        evidence = Evidence(
            source="test.csv",
            source_type=SourceType.CSV,
            records=[
                {"asset": "payroll-app", "public": "false", "internet_facing": "false"},
            ],
        )
        ver = self.engine.verify(assumption, [evidence])
        assert ver.result == VerificationResult.VERIFIED

    def test_mfa_contradicted(self):
        assumption = Assumption(
            claim_id="clm_5",
            text="System assumes identity posture: All administrative access requires MFA.",
            assumption_type=AssumptionType.IDENTITY,
        )
        evidence = Evidence(
            source="test.csv",
            source_type=SourceType.CSV,
            records=[
                {"user": "alice", "mfa": "true"},
                {"user": "bob", "mfa": "false"},
            ],
        )
        ver = self.engine.verify(assumption, [evidence])
        assert ver.result == VerificationResult.PARTIALLY_VERIFIED

    def test_configuration_verified(self):
        assumption = Assumption(
            claim_id="clm_6",
            text="System assumes configuration state: All financial data is backed up daily.",
            assumption_type=AssumptionType.CONFIGURATION,
        )
        evidence = Evidence(
            source="test.csv",
            source_type=SourceType.CSV,
            records=[
                {"resource": "payroll-db", "enabled": "true", "status": "active"},
                {"resource": "finance-fs", "enabled": "true", "status": "active"},
            ],
        )
        ver = self.engine.verify(assumption, [evidence])
        assert ver.result == VerificationResult.VERIFIED

    def test_no_evidence_unknown(self):
        assumption = Assumption(
            claim_id="clm_7",
            text="System assumes: This has no matching evidence.",
            assumption_type=AssumptionType.ACCESS,
        )
        ver = self.engine.verify(assumption, [])
        assert ver.result == VerificationResult.UNKNOWN
