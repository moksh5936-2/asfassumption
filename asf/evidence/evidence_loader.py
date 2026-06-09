from __future__ import annotations

from datetime import datetime, timezone
from pathlib import Path
from typing import Optional

from asf.ingestion.pipeline import IngestionPipeline
from asf.models import Evidence, SourceType

from .schema_adapter import EvidenceSchemaAdapter


class EvidenceLoader:
    def __init__(self, field_mappings: dict[str, str] | None = None):
        self.pipeline = IngestionPipeline()
        self.schema_adapter = EvidenceSchemaAdapter(field_mappings)

    def load(self, filepath: str | Path, source_type: SourceType | None = None, auto_map: bool = False) -> Evidence:
        filepath = Path(filepath)
        if source_type is None:
            source_type = self.pipeline.detect_type(filepath)

        records = self.pipeline.parse_to_records(filepath)

        if auto_map and records:
            inferred = self.schema_adapter.infer_field_mapping(records)
            self.schema_adapter.field_mappings.update(inferred)

        evidence = Evidence(
            source=str(filepath),
            source_type=source_type,
            timestamp=datetime.now(timezone.utc),
            content=self.pipeline.parse_text(filepath),
            records=records,
            metadata={
                "filename": filepath.name,
                "size_bytes": filepath.stat().st_size,
                "extension": filepath.suffix,
                "auto_mapped": auto_map,
            },
        )

        evidence = self.schema_adapter.normalize_evidence(evidence)
        return evidence
