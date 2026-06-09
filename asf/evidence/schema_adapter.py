from __future__ import annotations

import re
from typing import Any, Optional

from asf.models import Evidence


class EvidenceSchemaAdapter:
    def __init__(self, field_mappings: dict[str, str] | None = None):
        self.field_mappings = field_mappings or {}

    def normalize_evidence(self, evidence: Evidence) -> Evidence:
        if not evidence.records or not self.field_mappings:
            return evidence

        normalized_records: list[dict[str, Any]] = []
        for record in evidence.records:
            normalized: dict[str, Any] = {}
            for key, value in record.items():
                matched = self._map_field(key)
                if matched:
                    normalized[matched] = value
                else:
                    normalized[key] = value
            normalized_records.append(normalized)

        evidence.records = normalized_records
        evidence.metadata["field_mappings_applied"] = list(self.field_mappings.keys())
        return evidence

    def _map_field(self, field_name: str) -> Optional[str]:
        key_lower = field_name.lower().strip()

        exact_match = self.field_mappings.get(key_lower)
        if exact_match:
            return exact_match

        for pattern, target in self.field_mappings.items():
            if pattern in key_lower or key_lower in pattern:
                return target

        return None

    def infer_field_mapping(self, records: list[dict[str, Any]]) -> dict[str, str]:
        if not records:
            return {}

        sample = records[0]
        mappings: dict[str, str] = {}

        canonical_fields = {
            "user": ["user", "username", "employee", "name", "email", "identity", "principal", "login"],
            "group": ["group", "department", "team", "unit", "division", "org", "role_group"],
            "resource": ["resource", "application", "system", "service", "target", "asset_name", "scope"],
            "permission": ["permission", "access", "role", "right", "privilege", "action", "level"],
            "mfa": ["mfa", "mfa_enabled", "multi_factor", "2fa", "totp", "mfa_status"],
            "enabled": ["enabled", "status", "state", "active", "configuration", "value"],
            "public": ["public", "exposed", "internet_facing", "is_public", "exposure", "public_facing"],
            "asset": ["asset", "host", "server", "system", "name", "service_name", "component"],
        }

        for canonical, alternatives in canonical_fields.items():
            for alt in alternatives:
                if alt in sample:
                    mappings[alt] = canonical
                    break

        return mappings
