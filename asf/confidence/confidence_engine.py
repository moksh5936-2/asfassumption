from __future__ import annotations

from datetime import datetime, timezone
from typing import Sequence

from asf.models import Evidence, Verification, VerificationResult


class ConfidenceEngine:
    def compute_verification_confidence(self, verification: Verification, evidence_list: Sequence[Evidence]) -> float:
        base = verification.confidence

        evidence_freshness = self._compute_evidence_freshness(evidence_list)
        evidence_coverage = self._compute_evidence_coverage(verification, evidence_list)
        evidence_completeness = self._compute_evidence_completeness(verification)

        factors = [base, evidence_freshness, evidence_coverage, evidence_completeness]
        weights = [0.4, 0.2, 0.2, 0.2]

        final = sum(f * w for f, w in zip(factors, weights))
        return round(max(0.0, min(1.0, final)), 4)

    def compute_assumption_confidence(self, verifications: list[Verification]) -> float:
        if not verifications:
            return 0.0

        scores = []
        for v in verifications:
            multiplier = {
                VerificationResult.VERIFIED: 1.0,
                VerificationResult.PARTIALLY_VERIFIED: 0.5,
                VerificationResult.CONTRADICTED: 0.1,
                VerificationResult.UNKNOWN: 0.0,
            }.get(v.result, 0.0)
            scores.append(v.confidence * multiplier)

        return round(sum(scores) / len(scores), 4) if scores else 0.0

    def _compute_evidence_freshness(self, evidence_list: Sequence[Evidence]) -> float:
        if not evidence_list:
            return 0.0

        now = datetime.now(timezone.utc)
        total_freshness = 0.0
        for ev in evidence_list:
            age_hours = (now - ev.timestamp).total_seconds() / 3600
            freshness = max(0.0, 1.0 - (age_hours / (30 * 24)))
            total_freshness += freshness

        return total_freshness / len(evidence_list)

    def _compute_evidence_coverage(self, verification: Verification, evidence_list: Sequence[Evidence]) -> float:
        if not evidence_list:
            return 0.0
        used_ids = set(verification.evidence_used)
        total = len(evidence_list)
        if total == 0:
            return 0.0
        return len(used_ids) / total

    def _compute_evidence_completeness(self, verification: Verification) -> float:
        if verification.result == VerificationResult.VERIFIED:
            return 1.0
        if verification.result == VerificationResult.PARTIALLY_VERIFIED:
            return 0.5
        if verification.result == VerificationResult.CONTRADICTED:
            return 1.0
        return 0.0
