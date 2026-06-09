from __future__ import annotations

from pathlib import Path


class PDFParser:
    def parse(self, filepath: str | Path) -> str:
        filepath = Path(filepath)
        try:
            import pdfplumber
        except ImportError:
            raise ImportError("pdfplumber is required for PDF parsing: pip install pdfplumber")

        text_parts: list[str] = []
        with pdfplumber.open(str(filepath)) as pdf:
            for page in pdf.pages:
                text = page.extract_text()
                if text:
                    text_parts.append(text)
        return "\n\n".join(text_parts)
