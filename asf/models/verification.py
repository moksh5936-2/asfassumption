from __future__ import annotations

from datetime import datetime, timezone
from typing import Optional
from uuid import uuid4

from pydantic import BaseModel, Field

from .enums import VerificationResult


class Verification(BaseModel):
    id: str = Field(default_factory=lambda: f"vrf_{uuid4().hex[:12]}")
    assumption_id: str
    evidence_used: list[str] = Field(default_factory=list)
    result: VerificationResult = VerificationResult.UNKNOWN
    confidence: float = Field(default=0.0, ge=0.0, le=1.0)
    reasoning: str = ""
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    details: dict = Field(default_factory=dict)

    def __hash__(self) -> int:
        return hash(self.id)

    def __eq__(self, other: object) -> bool:
        if isinstance(other, Verification):
            return self.id == other.id
        return NotImplemented
