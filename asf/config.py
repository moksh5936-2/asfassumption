from __future__ import annotations

import os
from pathlib import Path
from typing import Any, Optional

DEFAULT_CONFIG_PATH = Path("asf.config.yaml")


class ASFConfig:
    def __init__(self, data: dict[str, Any] | None = None):
        d = data or {}
        self.evidence_schema: dict[str, str] = d.get("evidence_schema", {})
        self.claim_patterns: list[str] = d.get("claim_patterns", [])
        self.assumption_weights: dict[str, int] = d.get("assumption_weights", {})
        self.severity_rules: list[dict[str, Any]] = d.get("severity_rules", [])
        self.llm: dict[str, Any] = d.get("llm", {})
        self.db_path: str = d.get("db_path", "asf_validator.db")
        self.evidence_dirs: list[str] = d.get("evidence_dirs", [])
        self.field_mappings: dict[str, str] = d.get("field_mappings", {})

    @classmethod
    def load(cls, path: str | Path | None = None) -> "ASFConfig":
        from asf.settings import load_config
        return load_config(path)

    @classmethod
    def default(cls) -> "ASFConfig":
        return cls()


DEFAULT_CONFIG_YAML = """# ASF Validator v0.1 Configuration
# See: https://github.com/anomalyco/asf-validator

# Database path
db_path: "asf_validator.db"

# Evidence column name mappings
# Maps expected field names to actual column names in your evidence files
evidence_schema:
  user: "user"
  group: "group"
  resource: "resource"
  permission: "permission"
  department: "department"
  role: "role"
  mfa: "mfa_enabled"
  public: "public"
  internet_facing: "internet_facing"
  enabled: "enabled"
  status: "status"
  asset: "asset"
  environment: "environment"

# Additional claim extraction patterns (regex)
claim_patterns: []

# Assumption type scoring weights (higher = more likely to match)
assumption_weights:
  IDENTITY: 1.0
  ACCESS: 1.0
  NETWORK: 1.0
  CONFIGURATION: 1.0
  PROCESS: 1.0
  DOCUMENTATION: 1.0
  DEPENDENCY: 1.0
  GOVERNANCE: 1.0

# Custom severity rules
# Override gap severity based on assumption type + verification confidence
severity_rules: []

# LLM provider configuration
llm:
  provider: ""            # "openai" or "ollama"
  api_key: ""             # env var: ASF_LLM_API_KEY
  base_url: ""            # env var: ASF_LLM_BASE_URL
  model: ""               # env var: ASF_LLM_MODEL
  temperature: 0.1

# Default evidence directories to search
evidence_dirs: []

# Field mappings for evidence schema adapter
# Maps column name patterns found in user files -> standard field names
field_mappings:
  "employee": "user"
  "username": "user"
  "name": "user"
  "email": "user"
  "identity": "user"
  "principal": "user"
  "member": "user"
  "department": "group"
  "team": "group"
  "unit": "group"
  "division": "group"
  "org": "group"
  "application": "resource"
  "system": "resource"
  "service": "resource"
  "target": "resource"
  "access": "permission"
  "right": "permission"
  "privilege": "permission"
  "action": "permission"
  "public_facing": "public"
  "is_public": "public"
  "exposed": "public"
  "exposure": "public"
  "multi_factor": "mfa"
  "mfa_enabled": "mfa"
  "2fa": "mfa"
  "totp": "mfa"
  "state": "enabled"
  "active": "enabled"
  "host": "asset"
  "server": "asset"
  "system_name": "asset"
"""
