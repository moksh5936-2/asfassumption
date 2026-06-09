from __future__ import annotations

from pathlib import Path


class TXTParser:
    def parse(self, filepath: str | Path) -> str:
        filepath = Path(filepath)
        return filepath.read_text(encoding="utf-8")
