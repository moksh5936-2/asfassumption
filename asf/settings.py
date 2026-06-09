from __future__ import annotations

import os
from pathlib import Path
from typing import Any, Optional

import yaml

from asf.config import ASFConfig, DEFAULT_CONFIG_PATH, DEFAULT_CONFIG_YAML


def _merge_config(env: dict, defaults: dict) -> dict:
    result = dict(defaults)
    for k, v in env.items():
        if k in result and isinstance(result[k], dict) and isinstance(v, dict):
            result[k] = _merge_config(v, result[k])
        elif v is not None and v != "":
            result[k] = v
    return result


def load_config(path: str | Path | None = None) -> ASFConfig:
    config_path = Path(path) if path else DEFAULT_CONFIG_PATH
    config_data: dict[str, Any] = yaml.safe_load(DEFAULT_CONFIG_YAML) or {}

    if config_path.exists():
        with open(config_path) as f:
            user_data = yaml.safe_load(f) or {}
            config_data = _merge_config(user_data, config_data)

    llm = config_data.get("llm", {})
    llm["api_key"] = llm.get("api_key") or os.environ.get("ASF_LLM_API_KEY", "")
    llm["base_url"] = llm.get("base_url") or os.environ.get("ASF_LLM_BASE_URL", "")
    llm["model"] = llm.get("model") or os.environ.get("ASF_LLM_MODEL", "")
    config_data["llm"] = llm

    config_data["db_path"] = config_data.get("db_path") or os.environ.get("ASF_DB_PATH", "asf_validator.db")

    return ASFConfig(config_data)


def write_default_config(path: str | Path = DEFAULT_CONFIG_PATH) -> Path:
    path = Path(path)
    if not path.exists():
        path.write_text(DEFAULT_CONFIG_YAML)
    return path
