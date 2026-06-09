from __future__ import annotations

from datetime import datetime, timezone
from typing import Any, Optional
from uuid import uuid4

from pydantic import BaseModel, Field

from .enums import SourceType


class Evidence(BaseModel):
    id: str = Field(default_factory=lambda: f"evd_{uuid4().hex[:12]}")
    source: str
    source_type: SourceType
    timestamp: datetime = Field(default_factory=lambda: datetime.now(timezone.utc))
    content: Any = None
    confidence: float = Field(default=0.8, ge=0.0, le=1.0)
    metadata: dict[str, Any] = Field(default_factory=dict)
    records: list[dict[str, Any]] = Field(default_factory=list)

    def __hash__(self) -> int:
        return hash(self.id)

    def __eq__(self, other: object) -> bool:
        if isinstance(other, Evidence):
            return self.id == other.id
        return NotImplemented
