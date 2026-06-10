package models

import "encoding/json"

type AssumptionType string

const (
	AssumptionTypeIDENTITY       AssumptionType = "IDENTITY"
	AssumptionTypeACCESS         AssumptionType = "ACCESS"
	AssumptionTypeNETWORK        AssumptionType = "NETWORK"
	AssumptionTypeCONFIGURATION  AssumptionType = "CONFIGURATION"
	AssumptionTypePROCESS        AssumptionType = "PROCESS"
	AssumptionTypeDOCUMENTATION  AssumptionType = "DOCUMENTATION"
	AssumptionTypeDEPENDENCY     AssumptionType = "DEPENDENCY"
	AssumptionTypeGOVERNANCE     AssumptionType = "GOVERNANCE"
)

type VerificationStatus string

const (
	VerificationStatusPENDING      VerificationStatus = "1"
	VerificationStatusIN_REVIEW    VerificationStatus = "2"
	VerificationStatusVERIFIED     VerificationStatus = "3"
	VerificationStatusCONTRADICTED VerificationStatus = "4"
	VerificationStatusUNKNOWN      VerificationStatus = "5"
)

type VerificationResult string

const (
	VerificationResultVERIFIED            VerificationResult = "VERIFIED"
	VerificationResultPARTIALLY_VERIFIED  VerificationResult = "PARTIALLY_VERIFIED"
	VerificationResultCONTRADICTED        VerificationResult = "CONTRADICTED"
	VerificationResultUNKNOWN             VerificationResult = "UNKNOWN"
)

type GapSeverity string

const (
	GapSeverityCRITICAL GapSeverity = "CRITICAL"
	GapSeverityHIGH     GapSeverity = "HIGH"
	GapSeverityMEDIUM   GapSeverity = "MEDIUM"
	GapSeverityLOW      GapSeverity = "LOW"
	GapSeverityINFO     GapSeverity = "INFO"
)

type GapType string

const (
	GapTypeACCESS         GapType = "ACCESS_GAP"
	GapTypeIDENTITY       GapType = "IDENTITY_GAP"
	GapTypeNETWORK        GapType = "NETWORK_GAP"
	GapTypeCONFIGURATION  GapType = "CONFIGURATION_GAP"
	GapTypePROCESS        GapType = "PROCESS_GAP"
	GapTypeDOCUMENTATION  GapType = "DOCUMENTATION_GAP"
	GapTypeDEPENDENCY     GapType = "DEPENDENCY_GAP"
	GapTypeGOVERNANCE     GapType = "GOVERNANCE_GAP"
	GapTypeEVIDENCE       GapType = "EVIDENCE_GAP"
	GapTypeVERIFICATION   GapType = "VERIFICATION_GAP"
)

type SourceType string

const (
	SourceTypePDF               SourceType = "PDF"
	SourceTypeDOCX              SourceType = "DOCX"
	SourceTypeTXT               SourceType = "TXT"
	SourceTypeCSV               SourceType = "CSV"
	SourceTypeJSON              SourceType = "JSON"
	SourceTypeIAMExport         SourceType = "IAM_EXPORT"
	SourceTypeACLList           SourceType = "ACL_LIST"
	SourceTypeFirewallRules     SourceType = "FIREWALL_RULES"
	SourceTypeRouteTables       SourceType = "ROUTE_TABLES"
	SourceTypeSecurityGroups    SourceType = "SECURITY_GROUPS"
	SourceTypeConfigExport      SourceType = "CONFIG_EXPORT"
	SourceTypeAuditLog          SourceType = "AUDIT_LOG"
	SourceTypePolicyDocument    SourceType = "POLICY_DOCUMENT"
	SourceTypeRunbook           SourceType = "RUNBOOK"
	SourceTypeArchitectureDoc   SourceType = "ARCHITECTURE_DOC"
	SourceTypeUNKNOWN           SourceType = "UNKNOWN"
)

func (st SourceType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(st))
}

func (st *SourceType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*st = SourceType(s)
	return nil
}

func (at AssumptionType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(at))
}

func (at *AssumptionType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*at = AssumptionType(s)
	return nil
}

func (vr VerificationResult) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(vr))
}

func (vr *VerificationResult) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*vr = VerificationResult(s)
	return nil
}

func (gs GapSeverity) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(gs))
}

func (gs *GapSeverity) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*gs = GapSeverity(s)
	return nil
}

func (gt GapType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(gt))
}

func (gt *GapType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*gt = GapType(s)
	return nil
}

func (vs VerificationStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(vs))
}

func (vs *VerificationStatus) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*vs = VerificationStatus(s)
	return nil
}
