from __future__ import annotations

from datetime import datetime, timezone
from typing import Optional
from uuid import uuid4

from pydantic import BaseModel, Field

from .enums import AssumptionType, VerificationStatus


class Assumption(BaseModel):
    id: str = Field(default_factory=lambda: f"asm_{uuid4().hex[:12]}")
    claim_id: str
    text: str
    assumption_type: AssumptionType
    verification_status: VerificationStatus = VerificationStatus.PENDING
    confidence: float = Field(default=0.0, ge=0.0, le=1.0)
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    keywords: list[str] = Field(default_factory=list)

    def __hash__(self) -> int:
        return hash(self.id)

    def __eq__(self, other: object) -> bool:
        if isinstance(other, Assumption):
            return self.id == other.id
        return NotImplemented
