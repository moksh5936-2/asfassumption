package confidencex

type ConfidenceInput struct {
	AssumptionID         string
	AssumptionText       string
	Component            string
	Category             string
	Risk                 string
	Confidence           float64
	EvidenceSources      []string
	SourceComponents     []string
	SourceRelationships  []string
	Keywords             []string
	Rationale            string
	VerificationStatus   string
	Domain               string
	SupportingFactTexts  []string
	SupportingFactIDs    []string
	FactCategories       []string
	DependencyCentrality float64
	FailureRadius        int
	HasTrustChain        bool
	HasCoverageGap       bool
	BlindSpotScore       float64
}

type ConfidenceStability string

const (
	StabilityVeryStable        ConfidenceStability = "Very Stable"
	StabilityStable            ConfidenceStability = "Stable"
	StabilityModerate          ConfidenceStability = "Moderate"
	StabilityWeak              ConfidenceStability = "Weak"
	StabilityHighlySpeculative ConfidenceStability = "Highly Speculative"
)

type ConfidenceFactor struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"` // "positive" or "negative"
	Impact      float64 `json:"impact"`
	Description string  `json:"description"`
}

type FactContribution struct {
	FactID       string  `json:"fact_id"`
	FactText     string  `json:"fact_text"`
	Contribution float64 `json:"contribution"`
	IsPositive   bool    `json:"is_positive"`
}

type EvidenceContribution struct {
	EvidenceID string  `json:"evidence_id"`
	Present    bool    `json:"present"`
	Impact     float64 `json:"impact"`
	Label      string  `json:"label"`
}

type ContributionStrength string

const (
	StrengthStrong   ContributionStrength = "Strong"
	StrengthModerate ContributionStrength = "Moderate"
	StrengthWeak     ContributionStrength = "Weak"
)

type DomainContribution struct {
	Domain    string               `json:"domain"`
	Influence float64              `json:"influence"`
	Reason    string               `json:"reason"`
	Strength  ContributionStrength `json:"strength"`
}

type TrustContribution struct {
	HasTrustChain          bool    `json:"has_trust_chain"`
	ChainInfluence         float64 `json:"chain_influence"`
	DependencyCentrality   float64 `json:"dependency_centrality"`
	FailureRadiusInfluence float64 `json:"failure_radius_influence"`
}

type ConfidenceBreakdown struct {
	AssumptionID            string                 `json:"assumption_id"`
	AssumptionText          string                 `json:"assumption_text"`
	FinalConfidence         float64                `json:"final_confidence"`
	AdjustedConfidence      float64                `json:"adjusted_confidence,omitempty"`
	StabilityClass          ConfidenceStability    `json:"stability_class"`
	StabilityReason         string                 `json:"stability_reason"`
	PositiveFactors         []ConfidenceFactor     `json:"positive_factors"`
	NegativeFactors         []ConfidenceFactor     `json:"negative_factors"`
	SupportingFacts         []FactContribution     `json:"supporting_facts"`
	EvidenceContributions   []EvidenceContribution `json:"evidence_contributions"`
	DomainContribution      *DomainContribution    `json:"domain_contribution,omitempty"`
	TrustContribution       *TrustContribution     `json:"trust_contribution,omitempty"`
	WhyExists               string                 `json:"why_exists"`
	WhyUncertain            string                 `json:"why_uncertain"`
	WhatIncreasesConfidence string                 `json:"what_increases_confidence"`
	WhatDecreasesConfidence string                 `json:"what_decreases_confidence"`
}

type CISOTrustView struct {
	MostTrustedFindings       []ConfidenceBreakdown `json:"most_trusted_findings"`
	LeastTrustedFindings      []ConfidenceBreakdown `json:"least_trusted_findings"`
	MostCriticalLowConfidence []ConfidenceBreakdown `json:"most_critical_low_confidence"`
	HighestRiskUnknowns       []ConfidenceBreakdown `json:"highest_risk_unknowns"`
}

type ArchitectReviewView struct {
	RequiringValidation []ConfidenceBreakdown `json:"requiring_validation"`
	WeakSupport         []ConfidenceBreakdown `json:"weak_support"`
	StrongSupport       []ConfidenceBreakdown `json:"strong_support"`
}

type ConfidenceOutput struct {
	Breakdowns          []ConfidenceBreakdown `json:"breakdowns"`
	CISOTrustView       *CISOTrustView        `json:"ciso_trust_view,omitempty"`
	ArchitectReviewView *ArchitectReviewView  `json:"architect_review_view,omitempty"`
	Domain              string                `json:"domain"`
	GeneratedAt         string                `json:"generated_at"`
}
