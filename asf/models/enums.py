from enum import Enum, auto


class _StrEnum(str, Enum):
    """Backport of Python 3.11's StrEnum for Python 3.9 compatibility."""

    def __str__(self) -> str:
        return self.value


class AssumptionType(_StrEnum):
    IDENTITY = "IDENTITY"
    ACCESS = "ACCESS"
    NETWORK = "NETWORK"
    CONFIGURATION = "CONFIGURATION"
    PROCESS = "PROCESS"
    DOCUMENTATION = "DOCUMENTATION"
    DEPENDENCY = "DEPENDENCY"
    GOVERNANCE = "GOVERNANCE"


class VerificationStatus(_StrEnum):
    PENDING = auto()
    IN_REVIEW = auto()
    VERIFIED = auto()
    CONTRADICTED = auto()
    UNKNOWN = auto()


class VerificationResult(_StrEnum):
    VERIFIED = "VERIFIED"
    PARTIALLY_VERIFIED = "PARTIALLY_VERIFIED"
    CONTRADICTED = "CONTRADICTED"
    UNKNOWN = "UNKNOWN"


class GapSeverity(_StrEnum):
    CRITICAL = "CRITICAL"
    HIGH = "HIGH"
    MEDIUM = "MEDIUM"
    LOW = "LOW"
    INFO = "INFO"


class GapType(_StrEnum):
    ACCESS_GAP = "ACCESS_GAP"
    IDENTITY_GAP = "IDENTITY_GAP"
    NETWORK_GAP = "NETWORK_GAP"
    CONFIGURATION_GAP = "CONFIGURATION_GAP"
    PROCESS_GAP = "PROCESS_GAP"
    DOCUMENTATION_GAP = "DOCUMENTATION_GAP"
    DEPENDENCY_GAP = "DEPENDENCY_GAP"
    GOVERNANCE_GAP = "GOVERNANCE_GAP"
    EVIDENCE_GAP = "EVIDENCE_GAP"
    VERIFICATION_GAP = "VERIFICATION_GAP"


class SourceType(_StrEnum):
    PDF = "PDF"
    DOCX = "DOCX"
    TXT = "TXT"
    CSV = "CSV"
    JSON = "JSON"
    IAM_EXPORT = "IAM_EXPORT"
    ACL_LIST = "ACL_LIST"
    FIREWALL_RULES = "FIREWALL_RULES"
    ROUTE_TABLES = "ROUTE_TABLES"
    SECURITY_GROUPS = "SECURITY_GROUPS"
    CONFIG_EXPORT = "CONFIG_EXPORT"
    AUDIT_LOG = "AUDIT_LOG"
    POLICY_DOCUMENT = "POLICY_DOCUMENT"
    RUNBOOK = "RUNBOOK"
    ARCHITECTURE_DOC = "ARCHITECTURE_DOC"
    UNKNOWN = "UNKNOWN"
