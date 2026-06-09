from __future__ import annotations

from datetime import datetime, timezone
from pathlib import Path
from typing import Optional
from uuid import uuid4

from asf.models import SourceType

from .csv_parser import CSVParser
from .docx_parser import DOCXParser
from .json_parser import JSONParser
from .pdf_parser import PDFParser
from .txt_parser import TXTParser


class IngestionPipeline:
    def __init__(self):
        self.pdf_parser = PDFParser()
        self.docx_parser = DOCXParser()
        self.txt_parser = TXTParser()
        self.csv_parser = CSVParser()
        self.json_parser = JSONParser()

    def detect_type(self, filepath: str | Path) -> SourceType:
        ext = Path(filepath).suffix.lower()
        mapping = {
            ".pdf": SourceType.PDF,
            ".docx": SourceType.DOCX,
            ".txt": SourceType.TXT,
            ".csv": SourceType.CSV,
            ".json": SourceType.JSON,
        }
        return mapping.get(ext, SourceType.UNKNOWN)

    def parse_text(self, filepath: str | Path) -> str:
        ftype = self.detect_type(filepath)
        if ftype == SourceType.PDF:
            return self.pdf_parser.parse(filepath)
        if ftype == SourceType.DOCX:
            return self.docx_parser.parse(filepath)
        if ftype == SourceType.TXT:
            return self.txt_parser.parse(filepath)
        if ftype == SourceType.CSV:
            return self.csv_parser.parse(filepath)
        if ftype == SourceType.JSON:
            return self.json_parser.parse(filepath)
        return self.txt_parser.parse(filepath)

    def parse_to_records(self, filepath: str | Path) -> list[dict]:
        ftype = self.detect_type(filepath)
        if ftype == SourceType.CSV:
            return self.csv_parser.parse_to_records(filepath)
        if ftype == SourceType.JSON:
            data = self.json_parser.parse_to_dict(filepath)
            if isinstance(data, list):
                return data
            return [data]
        return []

    def get_document_metadata(self, filepath: str | Path) -> dict:
        filepath = Path(filepath)
        return {
            "id": f"doc_{uuid4().hex[:12]}",
            "filename": filepath.name,
            "file_type": self.detect_type(filepath).value,
            "size_bytes": filepath.stat().st_size,
            "created_at": datetime.now(timezone.utc).isoformat(),
        }
