package intelligence

import "time"

// RiskLevel represents the severity of an assumption.
type RiskLevel string

const (
	RiskCritical RiskLevel = "Critical"
	RiskHigh     RiskLevel = "High"
	RiskMedium   RiskLevel = "Medium"
	RiskLow      RiskLevel = "Low"
)

// StrideCategory represents a STRIDE threat category.
type StrideCategory string

const (
	StrideSpoofing        StrideCategory = "Spoofing"
	StrideTampering       StrideCategory = "Tampering"
	StrideRepudiation     StrideCategory = "Repudiation"
	StrideInfoDisclosure  StrideCategory = "Information Disclosure"
	StrideDenialOfService StrideCategory = "Denial of Service"
	StrideElevationPriv   StrideCategory = "Elevation of Privilege"
)

// Component represents a system component.
type Component struct {
	ID    string
	Label string
}

// Relation represents a relationship between components.
type Relation struct {
	Source string
	Target string
	Label  string
}

// ArchDescription represents a parsed architecture.
type ArchDescription struct {
	Name          string
	Components    []Component
	Relationships []Relation
	Policies      []string
	RawText       string

	ExplicitAssumptions []string
	SecurityControls    map[string][]string
	Compliance          []string
	ExpectedResults     map[string]interface{}
	ValidationCriteria  []string
	Notes               []string
}

// Assumption represents a security assumption.
type Assumption struct {
	ID          string
	Description string
	Component   string
	Category    string
	Risk        RiskLevel
	Stride      []StrideCategory
	Likelihood  int
	Impact      int
	Confidence  float64
	Keywords    []string

	SourceNode string
	SourceLine int

	SourceType    string
	SourceSection string
	SourceIndex   int
	SourceFile    string

	VerificationStatus string

	EvidenceSources      []string
	SourceComponents     []string
	SourceRelationships  []string
	Rationale            string
	StrideJustifications []StrideJustification
	RiskJustification    *RiskJustification
	ReviewStatus         string
	ReviewNotes          string
	ReviewTimestamp      time.Time

	QualityScore float64
}

// StrideJustification explains why a STRIDE category was assigned.
type StrideJustification struct {
	Category           StrideCategory
	Reason             string
	Confidence         float64
	ConfidenceReason   string
	MatchedRuleIndexes []int
	MatchedKeywords    []string
}

// RiskJustification explains the risk assessment.
type RiskJustification struct {
	Likelihood        int
	LikelihoodReason  string
	LikelihoodFactors []LikelihoodFactor
	Impact            int
	ImpactReason      string
	ImpactFactors     []ImpactFactor
	RiskScore         int
	RiskLevel         RiskLevel
	RiskReason        string
	Confidence        float64
	ConfidenceReason  string
}

// LikelihoodFactor represents a likelihood factor.
type LikelihoodFactor struct {
	Factor string
	Value  int
	Reason string
}

// ImpactFactor represents an impact factor.
type ImpactFactor struct {
	Factor string
	Value  int
	Reason string
}

// ControlDetail represents a security control.
type ControlDetail struct {
	ID                   string
	Name                 string
	Description          string
	Category             string
	STRIDECovered        []StrideCategory
	MitigatedAssumptions []string
	Component            string
	Priority             string
}

// Contradiction represents a detected contradiction.
type Contradiction struct {
	ID                  string
	Severity            RiskLevel
	Description         string
	Explanation         string
	AffectedAssumptions []string
	Evidence            []string
	RuleName            string
}

// TrustBoundary represents a discovered trust boundary.
type TrustBoundary struct {
	Type        string
	Components  []string
	RiskLevel   RiskLevel
	Description string
	Assumptions []string
}

// QualityScore represents the quality score of an assumption.
type QualityScore struct {
	Hiddenness             float64
	Impact                 float64
	Novelty                float64
	ArchitecturalRelevance float64
	Risk                   float64
	Confidence             float64
	Overall                float64
	Total                  float64
	Reason                 string
}

// EvidenceSummary represents a summary of evidence.
type EvidenceSummary struct {
	TotalSources       int
	SourceFiles        []string
	TotalComponents    int
	TotalRelationships int
}
