from __future__ import annotations

import csv
import json
from io import StringIO
from pathlib import Path
from typing import Any


class CSVParser:
    def parse(self, filepath: str | Path) -> str:
        filepath = Path(filepath)
        content = filepath.read_text(encoding="utf-8")
        return content

    def parse_to_records(self, filepath: str | Path) -> list[dict[str, Any]]:
        filepath = Path(filepath)
        with open(filepath, encoding="utf-8") as f:
            reader = csv.DictReader(f)
            return list(reader)

    def parse_to_json(self, filepath: str | Path) -> str:
        records = self.parse_to_records(filepath)
        return json.dumps(records, indent=2)
