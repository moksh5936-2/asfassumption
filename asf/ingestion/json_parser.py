from __future__ import annotations

import json
from pathlib import Path
from typing import Any


class JSONParser:
    def parse(self, filepath: str | Path) -> str:
        filepath = Path(filepath)
        data = json.loads(filepath.read_text(encoding="utf-8"))
        if isinstance(data, list):
            return json.dumps(data, indent=2)
        if isinstance(data, dict):
            return json.dumps(data, indent=2)
        return str(data)

    def parse_to_dict(self, filepath: str | Path) -> Any:
        filepath = Path(filepath)
        return json.loads(filepath.read_text(encoding="utf-8"))
