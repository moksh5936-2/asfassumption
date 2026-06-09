from __future__ import annotations

from datetime import datetime, timezone
from typing import Optional
from uuid import uuid4

from pydantic import BaseModel, Field


class Claim(BaseModel):
    id: str = Field(default_factory=lambda: f"clm_{uuid4().hex[:12]}")
    source_document: str
    source_location: Optional[str] = None
    text: str
    extraction_confidence: float = Field(default=0.5, ge=0.0, le=1.0)
    created_at: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    tags: list[str] = Field(default_factory=list)

    def __hash__(self) -> int:
        return hash(self.id)

    def __eq__(self, other: object) -> bool:
        if isinstance(other, Claim):
            return self.id == other.id
        return NotImplemented
