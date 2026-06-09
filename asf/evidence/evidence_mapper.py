from __future__ import annotations

from asf.models import AssumptionType, SourceType


class EvidenceMapper:
    TYPE_EVIDENCE_MAP: dict[AssumptionType, list[dict]] = {
        AssumptionType.IDENTITY: [
            {"source_type": SourceType.IAM_EXPORT, "description": "IAM user/role export"},
            {"source_type": SourceType.CSV, "description": "Identity listing (users, groups, roles)"},
            {"source_type": SourceType.JSON, "description": "Identity configuration"},
        ],
        AssumptionType.ACCESS: [
            {"source_type": SourceType.ACL_LIST, "description": "Access control list"},
            {"source_type": SourceType.IAM_EXPORT, "description": "IAM policy permissions"},
            {"source_type": SourceType.CSV, "description": "Permission matrix export"},
        ],
        AssumptionType.NETWORK: [
            {"source_type": SourceType.FIREWALL_RULES, "description": "Firewall rule set"},
            {"source_type": SourceType.SECURITY_GROUPS, "description": "Security group configurations"},
            {"source_type": SourceType.ROUTE_TABLES, "description": "Route table entries"},
            {"source_type": SourceType.CSV, "description": "Network exposure data"},
        ],
        AssumptionType.CONFIGURATION: [
            {"source_type": SourceType.CONFIG_EXPORT, "description": "System configuration"},
            {"source_type": SourceType.JSON, "description": "Configuration JSON export"},
            {"source_type": SourceType.CSV, "description": "Configuration audit log"},
        ],
        AssumptionType.PROCESS: [
            {"source_type": SourceType.AUDIT_LOG, "description": "Audit trail evidence"},
            {"source_type": SourceType.CSV, "description": "Process records"},
        ],
        AssumptionType.DOCUMENTATION: [
            {"source_type": SourceType.PDF, "description": "Original document"},
            {"source_type": SourceType.DOCX, "description": "Original document"},
        ],
        AssumptionType.DEPENDENCY: [
            {"source_type": SourceType.JSON, "description": "Service dependency map"},
            {"source_type": SourceType.CSV, "description": "Integration matrix"},
        ],
        AssumptionType.GOVERNANCE: [
            {"source_type": SourceType.AUDIT_LOG, "description": "Review/audit records"},
            {"source_type": SourceType.CSV, "description": "Governance tracking"},
        ],
    }

    def get_required_evidence(self, assumption_type: AssumptionType) -> list[dict]:
        return self.TYPE_EVIDENCE_MAP.get(assumption_type, [])

    def get_compatible_source_types(self, assumption_type: AssumptionType) -> list[SourceType]:
        evidence_defs = self.TYPE_EVIDENCE_MAP.get(assumption_type, [])
        return [e["source_type"] for e in evidence_defs]
