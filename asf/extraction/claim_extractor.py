from __future__ import annotations

import re
from typing import Optional

from asf.models import Claim


class ClaimExtractor:
    DECLARATIVE_PATTERNS = [
        r"(?:only|just)\s+.+\s+(?:can|may|has|have|should|must|will)",
        r"(?:all|every|each)\s+.+\s+(?:are|is|shall|must|will|should)",
        r"(?:is|are)\s+(?:not|never)\s+.+",
        r"(?:is|are)\s+\w+\s+(?:encrypted|backed\s+up|logged|audited|monitored|reviewed|tested|approved|validated|isolated|segmented|restricted|limited)",
        r"(?:cannot|can\s+not|must\s+not|shall\s+not|should\s+not)\s+.+",
        r"(?:requir(?:e|es|ed|ing)|ensures?|guarantees?|provides?|protects?|prevents?|blocks?|restricts?|limits?)\s+.+",
        r"(?:ensure[s]?|guarantee[s]?)\s+that\s+.+",
        r"(?:accessed?|accessible|available)\s+(?:only|exclusively|solely)\s+.+",
        r"(?:encrypt(?:ed|s|ion)|back(?:ed|s)?\s+ups?|backups?|log(?:ged|s|ging)?|audit(?:ed|s|ing)?|monitor(?:ed|s|ing)?)",
        r"(?:configured?|set\s+up|deployed?|implemented?)\s+(?:to|with|as)\s+.+",
        r"(?:review|test|approv|validat|certif)\w+\s+(?:are|is|shall|must|should|will|performed|conducted)",
        r"(?:separated?|isolated?|segmented?|partitioned?)\s+.+",
        r"(?:is|are)\s+restricted\s+to\s+.+",
        r"manage[sd]?\s+.+\s+(?:access|permissions)",
        r"(?:security\s+)?(?:review|audit|assessment)s?\s+.+\s+(?:are|is|conducted|performed|scheduled)",
    ]

    def extract(self, text: str, source_document: str = "", source_location: Optional[str] = None) -> list[Claim]:
        claims: list[Claim] = []
        seen: set[str] = set()

        lines = self._split_into_sentences(text)
        for line in lines:
            cleaned = line.strip()
            if not cleaned or len(cleaned) < 15:
                continue
            if self._is_declarative(cleaned):
                normalized = cleaned.lower().strip()
                if normalized not in seen:
                    seen.add(normalized)
                    confidence = self._compute_confidence(cleaned)
                    claim = Claim(
                        source_document=source_document,
                        source_location=source_location,
                        text=cleaned,
                        extraction_confidence=confidence,
                        tags=self._extract_tags(cleaned),
                    )
                    claims.append(claim)

        return claims

    def _split_into_sentences(self, text: str) -> list[str]:
        sentences = re.split(r"(?<=[.!?])\s+", text)
        result = []
        for s in sentences:
            s = s.strip()
            if s:
                result.append(s)
        return result

    def _is_declarative(self, text: str) -> bool:
        for pattern in self.DECLARATIVE_PATTERNS:
            if re.search(pattern, text, re.IGNORECASE):
                return True
        return False

    def _compute_confidence(self, text: str) -> float:
        score = 0.5
        strong_indicators = [
            r"\b(?:only|never|always|all|every|must|shall)\b",
            r"\b(?:encrypt|backup|audit|require|ensure|guarantee|manage)\b",
            r"\b(?:isolated|segmented|restricted|limited|conducted)\b",
        ]
        for pat in strong_indicators:
            if re.search(pat, text, re.IGNORECASE):
                score += 0.1
        return min(score, 0.95)

    def _extract_tags(self, text: str) -> list[str]:
        tags = []
        keyword_map = {
            "access": "access",
            "permission": "access",
            "identity": "identity",
            "mfa": "identity",
            "authentication": "identity",
            "network": "network",
            "firewall": "network",
            "internet": "network",
            "encrypt": "configuration",
            "backup": "configuration",
            "log": "configuration",
            "audit": "configuration",
            "review": "governance",
            "approve": "governance",
            "restrict": "access",
            "manage": "access",
            "process": "process",
            "procedure": "process",
            "test": "process",
            "document": "documentation",
            "policy": "documentation",
            "depend": "dependency",
            "integration": "dependency",
        }
        lowered = text.lower()
        for keyword, tag in keyword_map.items():
            if keyword in lowered:
                if tag not in tags:
                    tags.append(tag)
        return tags
