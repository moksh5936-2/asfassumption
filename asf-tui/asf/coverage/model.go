package coverage

type CoverageCategory string

const (
	CatIdentity      CoverageCategory = "identity"
	CatAuthorization CoverageCategory = "authorization"
	CatCryptography  CoverageCategory = "cryptography"
	CatMonitoring    CoverageCategory = "monitoring"
	CatResilience    CoverageCategory = "resilience"
	CatThirdParty    CoverageCategory = "third_party"
	CatOperational   CoverageCategory = "operational"
)

var AllCategories = []CoverageCategory{
	CatIdentity, CatAuthorization, CatCryptography,
	CatMonitoring, CatResilience, CatThirdParty, CatOperational,
}

type ExpectedAssumption struct {
	Category    CoverageCategory `json:"category"`
	Title       string           `json:"title"`
	Description string           `json:"description,omitempty"`
	Risk        string           `json:"risk"`
	Source      string           `json:"source"`
}

type ComponentExpectations struct {
	Component    string               `json:"component"`
	Expectations []ExpectedAssumption `json:"expectations"`
}

type CoverageMetric struct {
	Category       CoverageCategory `json:"category"`
	ExpectedCount  int              `json:"expected_count"`
	ObservedCount  int              `json:"observed_count"`
	CoveragePct    float64          `json:"coverage_percentage"`
	Risk           string           `json:"risk_level"`
	AttentionDelta float64          `json:"attention_delta"`
	Reason         string           `json:"reason,omitempty"`
}

type CoverageGap struct {
	Category       CoverageCategory `json:"category"`
	ExpectedCount  int              `json:"expected_count"`
	ObservedCount  int              `json:"observed_count"`
	MissingCount   int              `json:"missing_count"`
	CoveragePct    float64          `json:"coverage_percentage"`
	Risk           string           `json:"risk"`
	Recommendation string           `json:"recommendation"`
}

type DomainBlindSpot struct {
	Domain         string `json:"domain"`
	MissingArea    string `json:"missing_area"`
	Description    string `json:"description"`
	Risk           string `json:"risk"`
	Recommendation string `json:"recommendation"`
}

type BlindSpot struct {
	Category            CoverageCategory `json:"category"`
	BlindSpotID         string           `json:"blind_spot_id"`
	Title               string           `json:"title"`
	Description         string           `json:"description"`
	Risk                string           `json:"risk"`
	Score               float64          `json:"score"`
	Domain              string           `json:"domain,omitempty"`
	Component           string           `json:"component,omitempty"`
	TrustChainImpact    string           `json:"trust_chain_impact,omitempty"`
	ConsequenceSeverity string           `json:"consequence_severity,omitempty"`
	ComplianceRelevance string           `json:"compliance_relevance,omitempty"`
	Recommendation      string           `json:"recommendation"`
}

type CISOView struct {
	TopBlindSpots               []BlindSpot `json:"top_blind_spots"`
	DangerousMissingAssumptions []BlindSpot `json:"dangerous_missing_assumptions"`
	AreasRequiringReview        []string    `json:"areas_requiring_review"`
	HighestRiskUnknowns         []BlindSpot `json:"highest_risk_unknowns"`
}

type CoverageAssessment struct {
	Categories       []CoverageMetric        `json:"categories"`
	ComponentResults []ComponentExpectations `json:"component_results"`
	Gaps             []CoverageGap           `json:"gaps"`
}

type CoverageOutput struct {
	Assessment       *CoverageAssessment `json:"assessment,omitempty"`
	BlindSpots       []BlindSpot         `json:"blind_spots,omitempty"`
	DomainBlindSpots []DomainBlindSpot   `json:"domain_blind_spots,omitempty"`
	CISOView         *CISOView           `json:"ciso_view,omitempty"`

	Domain      string `json:"domain"`
	GeneratedAt string `json:"generated_at"`
}
