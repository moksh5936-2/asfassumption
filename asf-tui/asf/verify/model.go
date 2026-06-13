package verify

type EvidenceCategory string

const (
	EvCatIdentity      EvidenceCategory = "identity"
	EvCatAuthorization EvidenceCategory = "authorization"
	EvCatCryptography  EvidenceCategory = "cryptography"
	EvCatMonitoring    EvidenceCategory = "monitoring"
	EvCatResilience    EvidenceCategory = "resilience"
	EvCatThirdParty    EvidenceCategory = "third_party"
	EvCatOperational   EvidenceCategory = "operational"
)

type VerificationPriority string

const (
	VpCritical VerificationPriority = "Critical"
	VpHigh     VerificationPriority = "High"
	VpMedium   VerificationPriority = "Medium"
	VpLow      VerificationPriority = "Low"
)

type VerificationStatus string

const (
	VsVerified          VerificationStatus = "Verified"
	VsPartiallyVerified VerificationStatus = "Partially Verified"
	VsUnverified        VerificationStatus = "Unverified"
	VsNoEvidence        VerificationStatus = "No Evidence"
)

type VerificationEffort string

const (
	EffortLow    VerificationEffort = "Low"
	EffortMedium VerificationEffort = "Medium"
	EffortHigh   VerificationEffort = "High"
)

type EvidenceSourceType string

const (
	SourcePolicyDocument EvidenceSourceType = "policy_document"
	SourceConfiguration  EvidenceSourceType = "configuration"
	SourceAuditLog       EvidenceSourceType = "audit_log"
	SourceReport         EvidenceSourceType = "report"
	SourceInterview      EvidenceSourceType = "interview"
	SourceToolOutput     EvidenceSourceType = "tool_output"
	SourceArtifact       EvidenceSourceType = "artifact"
	SourceVendorDocument EvidenceSourceType = "vendor_document"
)

type EvidenceRequirement struct {
	Category    EvidenceCategory     `json:"category"`
	Title       string               `json:"title"`
	Description string               `json:"description,omitempty"`
	SourceTypes []EvidenceSourceType `json:"source_types"`
	Risk        string               `json:"risk"`
}

type EvidenceSource struct {
	Type        EvidenceSourceType `json:"type"`
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	Location    string             `json:"location,omitempty"`
	Optional    bool               `json:"optional,omitempty"`
}

type VerificationAction struct {
	Step        int    `json:"step"`
	Action      string `json:"action"`
	Description string `json:"description,omitempty"`
	Stakeholder string `json:"stakeholder,omitempty"`
}

type VerificationPlan struct {
	AssumptionID          string               `json:"assumption_id"`
	AssumptionText        string               `json:"assumption_text"`
	Category              EvidenceCategory     `json:"category"`
	Risk                  string               `json:"risk"`
	EvidenceRequired      []EvidenceSource     `json:"evidence_required,omitempty"`
	EvidencePresent       []EvidenceSource     `json:"evidence_present,omitempty"`
	EvidenceMissing       []EvidenceSource     `json:"evidence_missing,omitempty"`
	Confidence            float64              `json:"confidence"`
	Priority              VerificationPriority `json:"priority"`
	Status                VerificationStatus   `json:"status"`
	Effort                VerificationEffort   `json:"effort"`
	Actions               []VerificationAction `json:"actions,omitempty"`
	Stakeholders          []string             `json:"stakeholders,omitempty"`
	WhyVerify             string               `json:"why_verify,omitempty"`
	WhatToReview          string               `json:"what_to_review,omitempty"`
	WhatEvidenceToCollect string               `json:"what_evidence_to_collect,omitempty"`
	HowToValidate         string               `json:"how_to_validate,omitempty"`
	ExpectedTime          string               `json:"expected_time,omitempty"`
}

type VerificationRoadmap struct {
	AssumptionID   string               `json:"assumption_id"`
	AssumptionText string               `json:"assumption_text"`
	Priority       VerificationPriority `json:"priority"`
	Steps          []VerificationAction `json:"steps"`
	Effort         VerificationEffort   `json:"effort"`
	Stakeholders   []string             `json:"stakeholders,omitempty"`
}

type VerificationAssessment struct {
	Plans             []VerificationPlan `json:"plans"`
	VerifiedCount     int                `json:"verified_count"`
	PartialCount      int                `json:"partial_count"`
	UnverifiedCount   int                `json:"unverified_count"`
	NoEvidenceCount   int                `json:"no_evidence_count"`
	TotalAssumptions  int                `json:"total_assumptions"`
	OverallConfidence float64            `json:"overall_confidence"`
}

type CISOReviewView struct {
	TopAssumptionsToVerify []VerificationPlan `json:"top_assumptions_to_verify"`
	HighestRiskUnverified  []VerificationPlan `json:"highest_risk_unverified"`
	EvidenceGaps           []string           `json:"evidence_gaps"`
	VerificationBacklog    []VerificationPlan `json:"verification_backlog"`
	CriticalCount          int                `json:"critical_count"`
	HighCount              int                `json:"high_count"`
	MediumCount            int                `json:"medium_count"`
	LowCount               int                `json:"low_count"`
}

type VerificationOutput struct {
	Assessment  *VerificationAssessment `json:"assessment,omitempty"`
	CISOView    *CISOReviewView         `json:"ciso_view,omitempty"`
	Roadmaps    []VerificationRoadmap   `json:"roadmaps,omitempty"`
	Domain      string                  `json:"domain"`
	GeneratedAt string                  `json:"generated_at"`
}
