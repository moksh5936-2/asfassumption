from pydantic import BaseModel, Field

from .claim import Claim
from .assumption import Assumption
from .evidence import Evidence
from .verification import Verification
from .gap import Gap


class AnalysisResult(BaseModel):
    claims: list[Claim] = Field(default_factory=list)
    assumptions: list[Assumption] = Field(default_factory=list)
    evidence: list[Evidence] = Field(default_factory=list)
    verifications: list[Verification] = Field(default_factory=list)
    gaps: list[Gap] = Field(default_factory=list)

    @property
    def claims_found(self) -> int:
        return len(self.claims)

    @property
    def assumptions_found(self) -> int:
        return len(self.assumptions)

    @property
    def verified_count(self) -> int:
        from .enums import VerificationResult
        return sum(1 for v in self.verifications if v.result == VerificationResult.VERIFIED)

    @property
    def contradicted_count(self) -> int:
        from .enums import VerificationResult
        return sum(1 for v in self.verifications if v.result == VerificationResult.CONTRADICTED)

    @property
    def unknown_count(self) -> int:
        from .enums import VerificationResult
        return sum(1 for v in self.verifications if v.result == VerificationResult.UNKNOWN)

    @property
    def critical_gaps(self) -> int:
        from .enums import GapSeverity
        return sum(1 for g in self.gaps if g.severity == GapSeverity.CRITICAL)
