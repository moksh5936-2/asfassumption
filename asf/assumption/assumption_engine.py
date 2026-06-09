from __future__ import annotations

import re

from asf.models import Assumption, AssumptionType, Claim, VerificationStatus


class AssumptionEngine:
    TYPE_PATTERNS: dict[AssumptionType, list[str]] = {
        AssumptionType.IDENTITY: [
            r"\b(?:mfa|multi.?factor|identity|authentication|password|credential|role|group)\b",
            r"\b(?:only|just)\s+.+\s+(?:can|may|has)\s+.+\b(?:access|login|authenticate)\b",
        ],
        AssumptionType.ACCESS: [
            r"\b(?:access|permission|acl|allow|deny|grant|read|write|execute|admin)\b",
            r"\b(?:only|just)\s+.+\s+(?:can|may|has)\s+(?:access|permission)\b",
            r"\b(?:restricted?|limited?|blocked?)\s+(?:to|access)\b",
        ],
        AssumptionType.NETWORK: [
            r"\b(?:network|firewall|internet|subnet|vpc|vlan|segment|isolate|expose)\b",
            r"\b(?:not\s+)?(?:accessible|reachable)\s+(?:from|via|over)\b",
            r"\b(?:no\s+)?public\s+(?:access|exposure|facing)\b",
        ],
        AssumptionType.CONFIGURATION: [
            r"\b(?:encrypt|backup|log|audit|monitor|config|setting|parameter)\b",
            r"\b(?:encrypt(?:ed|ion)|backup|log(?:ging|ged)|audit(?:ing|ed)?|monitor(?:ing|ed)?)\b",
            r"\b(?:enabled?|disabled?|configured?)\s+(?:by|with|to|as)\b",
        ],
        AssumptionType.PROCESS: [
            r"\b(?:process|procedure|workflow|review|approve|test|sign.?off|approval)\b",
            r"\b(?:must|shall|should)\s+(?:be\s+)?(?:reviewed|tested|approved|validated)\b",
        ],
        AssumptionType.DOCUMENTATION: [
            r"\b(?:document|policy|runbook|procedure|guide|manual|readme|wiki)\b",
            r"\b(?:as\s+(?:per|described|documented|stated))\b",
        ],
        AssumptionType.DEPENDENCY: [
            r"\b(?:depend|integration|connect|communicate|rel(y|ies)|upstream|downstream)\b",
            r"\b(?:requires?|depends?\s+on|relies?\s+on)\b",
        ],
        AssumptionType.GOVERNANCE: [
            r"\b(?:review|audit|compliance|regulat|policy|standard|framework|govern)\b",
            r"\b(?:annually|quarterly|monthly|regularly|periodically)\b",
            r"\b(?:reviewed|audited|approved|certified)\s+(?:annually|quarterly|monthly|regularly|periodically)\b",
        ],
    }

    def convert(self, claim: Claim) -> Assumption | None:
        atype = self._classify(claim.text)
        if atype is None:
            return None

        text = self._build_assumption_text(claim.text, atype)
        keywords = self._extract_keywords(claim.text)

        return Assumption(
            claim_id=claim.id,
            text=text,
            assumption_type=atype,
            verification_status=VerificationStatus.PENDING,
            keywords=keywords,
        )

    def convert_many(self, claims: list[Claim]) -> list[Assumption]:
        return [a for c in claims if (a := self.convert(c)) is not None]

    def _classify(self, text: str) -> AssumptionType | None:
        scores: dict[AssumptionType, int] = {}
        for atype, patterns in self.TYPE_PATTERNS.items():
            score = 0
            for pattern in patterns:
                matches = re.findall(pattern, text, re.IGNORECASE)
                score += len(matches)
            if score > 0:
                scores[atype] = score

        if not scores:
            return None

        return max(scores, key=scores.get)

    def _build_assumption_text(self, text: str, atype: AssumptionType) -> str:
        type_prefixes = {
            AssumptionType.IDENTITY: "System assumes identity posture: ",
            AssumptionType.ACCESS: "System assumes access control: ",
            AssumptionType.NETWORK: "System assumes network posture: ",
            AssumptionType.CONFIGURATION: "System assumes configuration state: ",
            AssumptionType.PROCESS: "System assumes process compliance: ",
            AssumptionType.DOCUMENTATION: "System assumes documentation accuracy: ",
            AssumptionType.DEPENDENCY: "System assumes dependency relationship: ",
            AssumptionType.GOVERNANCE: "System assumes governance compliance: ",
        }
        prefix = type_prefixes.get(atype, "System assumes: ")
        return f"{prefix}{text}"

    def _extract_keywords(self, text: str) -> list[str]:
        words = re.findall(r"\b[a-zA-Z]{3,}\b", text.lower())
        stopwords = {
            "the", "and", "for", "are", "but", "not", "you", "all", "can", "has",
            "have", "may", "must", "shall", "should", "will", "with", "from",
            "that", "this", "each", "every", "than", "then", "just", "been",
            "were", "was", "its", "also", "per", "via",
        }
        return [w for w in words if w not in stopwords]
