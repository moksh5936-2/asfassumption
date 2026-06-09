from __future__ import annotations

from asf.models import (
    Assumption,
    AssumptionType,
    Gap,
    GapSeverity,
    GapType,
    Verification,
    VerificationResult,
)

ASSUMPTION_TYPE_GAP_MAP = {
    AssumptionType.IDENTITY: GapType.IDENTITY_GAP,
    AssumptionType.ACCESS: GapType.ACCESS_GAP,
    AssumptionType.NETWORK: GapType.NETWORK_GAP,
    AssumptionType.CONFIGURATION: GapType.CONFIGURATION_GAP,
    AssumptionType.PROCESS: GapType.PROCESS_GAP,
    AssumptionType.DOCUMENTATION: GapType.DOCUMENTATION_GAP,
    AssumptionType.DEPENDENCY: GapType.DEPENDENCY_GAP,
    AssumptionType.GOVERNANCE: GapType.GOVERNANCE_GAP,
}


class GapEngine:
    def generate_gaps(self, assumptions: list[Assumption], verifications: list[Verification]) -> list[Gap]:
        gaps: list[Gap] = []
        ver_map = {v.assumption_id: v for v in verifications}

        for assumption in assumptions:
            verification = ver_map.get(assumption.id)
            if verification is None:
                gap = Gap(
                    assumption_id=assumption.id,
                    severity=GapSeverity.MEDIUM,
                    type=GapType.VERIFICATION_GAP,
                    description=f"Assumption '{assumption.text[:80]}...' has not been verified",
                    evidence_detail="No verification performed",
                )
                gaps.append(gap)
                continue

            if verification.result == VerificationResult.CONTRADICTED:
                gap_type = ASSUMPTION_TYPE_GAP_MAP.get(assumption.assumption_type, GapType.VERIFICATION_GAP)
                severity = self._determine_severity(assumption.assumption_type, verification)
                gap = Gap(
                    assumption_id=assumption.id,
                    severity=severity,
                    type=gap_type,
                    description=f"Assumption contradicted: {assumption.text[:120]}",
                    evidence_detail=verification.reasoning,
                )
                gaps.append(gap)

            elif verification.result == VerificationResult.PARTIALLY_VERIFIED:
                gap_type = ASSUMPTION_TYPE_GAP_MAP.get(assumption.assumption_type, GapType.VERIFICATION_GAP)
                gap = Gap(
                    assumption_id=assumption.id,
                    severity=GapSeverity.MEDIUM,
                    type=gap_type,
                    description=f"Assumption only partially verified: {assumption.text[:120]}",
                    evidence_detail=verification.reasoning,
                )
                gaps.append(gap)

            elif verification.result == VerificationResult.UNKNOWN:
                gap = Gap(
                    assumption_id=assumption.id,
                    severity=GapSeverity.LOW,
                    type=GapType.EVIDENCE_GAP,
                    description=f"Insufficient evidence to verify: {assumption.text[:120]}",
                    evidence_detail=verification.reasoning,
                )
                gaps.append(gap)

        return gaps

    def _determine_severity(self, assumption_type: AssumptionType, verification: Verification) -> GapSeverity:
        if verification.confidence >= 0.8:
            if assumption_type in (AssumptionType.ACCESS, AssumptionType.IDENTITY, AssumptionType.NETWORK):
                return GapSeverity.CRITICAL
            if assumption_type in (AssumptionType.CONFIGURATION, AssumptionType.GOVERNANCE):
                return GapSeverity.HIGH
            return GapSeverity.HIGH

        if verification.confidence >= 0.5:
            return GapSeverity.HIGH

        return GapSeverity.MEDIUM
