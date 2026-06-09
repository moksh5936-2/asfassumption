from __future__ import annotations

from pathlib import Path


class DOCXParser:
    def parse(self, filepath: str | Path) -> str:
        filepath = Path(filepath)
        try:
            from docx import Document
        except ImportError:
            raise ImportError("python-docx is required for DOCX parsing: pip install python-docx")

        doc = Document(str(filepath))
        paragraphs = [p.text for p in doc.paragraphs if p.text.strip()]
        return "\n\n".join(paragraphs)
