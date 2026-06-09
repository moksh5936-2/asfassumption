from __future__ import annotations

from datetime import datetime, timezone
from uuid import uuid4

from pydantic import BaseModel, Field

from .enums import GapSeverity, GapType


class Gap(BaseModel):
    id: str = Field(default_factory=lambda: f"gap_{uuid4().hex[:12]}")
    assumption_id: str
    severity: GapSeverity
    type: GapType
    description: str
    evidence_detail: str = ""
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))

    def __hash__(self) -> int:
        return hash(self.id)

    def __eq__(self, other: object) -> bool:
        if isinstance(other, Gap):
            return self.id == other.id
        return NotImplemented
