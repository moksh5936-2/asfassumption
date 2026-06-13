package narrative

import (
	"time"
)

// NarrativeOutput is the top-level container for all narrative content.
type NarrativeOutput struct {
	// Generated when the narrative was created
	GeneratedAt time.Time `json:"generated_at"`

	// Architecture-level narrative
	ArchitectureOverview ArchitectureOverview `json:"architecture_overview"`

	// Per-assumption narratives
	AssumptionNarratives []AssumptionNarrative `json:"assumption_narratives"`

	// Executive report
	ExecutiveReport ExecutiveReport `json:"executive_report"`

	// Technical summary
	TechnicalSummary TechnicalSummary `json:"technical_summary"`

	// Architect narrative (full text)
	ArchitectNarrative string `json:"architect_narrative"`
}

// ArchitectureOverview provides the high-level architectural context.
type ArchitectureOverview struct {
	Name              string   `json:"name"`
	Domain            string   `json:"domain"`
	TotalComponents   int      `json:"total_components"`
	TotalAssumptions  int      `json:"total_assumptions"`
	CriticalCount     int      `json:"critical_count"`
	HighCount         int      `json:"high_count"`
	KeyComponents     []string `json:"key_components"`
	TrustDependencies []string `json:"trust_dependencies"`
	Summary           string   `json:"summary"`
}

// AssumptionNarrative is the architect-style explanation for a single assumption.
type AssumptionNarrative struct {
	// Reference to the original assumption
	AssumptionID   string `json:"assumption_id"`
	AssumptionText string `json:"assumption_text"`

	// The five narrative sections
	Context                 string `json:"context"`
	WhyASFIdentifiedIt      string `json:"why_asf_identified_it"`
	ArchitecturalImportance string `json:"architectural_importance"`
	FailureConsequence      string `json:"failure_consequence"`
	SecurityRecommendation  string `json:"security_recommendation"`

	// Derived metadata
	RiskLevel        string   `json:"risk_level"`
	STRIDECategories []string `json:"stride_categories"`
	DependsOn        []string `json:"depends_on"`
	DownstreamImpact []string `json:"downstream_impact"`
	Confidence       float64  `json:"confidence"`
}

// ExecutiveReport is the C-level summary.
type ExecutiveReport struct {
	ArchitectureOverview    string                `json:"architecture_overview"`
	KeyAssumptions          []ExecutiveAssumption `json:"key_assumptions"`
	MostCriticalAssumptions []ExecutiveAssumption `json:"most_critical_assumptions"`
	HighImpactConsequences  []string              `json:"high_impact_consequences"`
	TrustDependencies       []string              `json:"trust_dependencies"`
	SinglePointsOfFailure   []string              `json:"single_points_of_failure"`
	ArchitecturalConcerns   []string              `json:"architectural_concerns"`
	RecommendedInvestments  []string              `json:"recommended_investments"`
}

// ExecutiveAssumption is a condensed view for executives.
type ExecutiveAssumption struct {
	Text           string `json:"text"`
	RiskLevel      string `json:"risk_level"`
	Consequence    string `json:"consequence"`
	BusinessImpact string `json:"business_impact"`
}

// TechnicalSummary is the detailed technical view.
type TechnicalSummary struct {
	ArchitectureSummary string                    `json:"architecture_summary"`
	AssumptionDetails   []TechnicalAssumption     `json:"assumption_details"`
	STRIDEDistribution  map[string]int            `json:"stride_distribution"`
	RiskDistribution    map[string]int            `json:"risk_distribution"`
	Dependencies        []TechnicalDependency     `json:"dependencies"`
	Recommendations     []TechnicalRecommendation `json:"recommendations"`
}

// TechnicalAssumption is the technical view of a single assumption.
type TechnicalAssumption struct {
	ID                string   `json:"id"`
	Description       string   `json:"description"`
	Component         string   `json:"component"`
	Category          string   `json:"category"`
	RiskLevel         string   `json:"risk_level"`
	STRIDECategories  []string `json:"stride_categories"`
	Likelihood        int      `json:"likelihood"`
	Impact            int      `json:"impact"`
	Confidence        float64  `json:"confidence"`
	EvidenceSources   []string `json:"evidence_sources"`
	Rationale         string   `json:"rationale"`
	DownstreamSystems []string `json:"downstream_systems"`
	FailureScenario   string   `json:"failure_scenario"`
	Recommendation    string   `json:"recommendation"`
}

// TechnicalDependency maps what depends on an assumption.
type TechnicalDependency struct {
	AssumptionID        string   `json:"assumption_id"`
	AssumptionText      string   `json:"assumption_text"`
	DependentComponents []string `json:"dependent_components"`
	DependentSystems    []string `json:"dependent_systems"`
	DependencyType      string   `json:"dependency_type"`
}

// TechnicalRecommendation is a technical recommendation.
type TechnicalRecommendation struct {
	AssumptionID         string   `json:"assumption_id"`
	AssumptionText       string   `json:"assumption_text"`
	Recommendation       string   `json:"recommendation"`
	Priority             string   `json:"priority"`
	ImplementationEffort string   `json:"implementation_effort"`
	MitigatesSTRIDE      []string `json:"mitigates_stride"`
}

// StyleRule defines a banned phrase or pattern.
type StyleRule struct {
	Pattern     string `json:"pattern"`
	Reason      string `json:"reason"`
	Replacement string `json:"replacement"`
}

// BannedPhrases are the style enforcement rules.
var BannedPhrases = []StyleRule{
	{"leverage", "Marketing language", "use"},
	{"synergy", "Marketing language", "cooperation"},
	{"holistic", "Marketing language", "comprehensive"},
	{"robust", "Marketing language", "strong"},
	{"seamless", "Marketing language", "transparent"},
	{"cutting-edge", "Marketing language", "modern"},
	{"state-of-the-art", "Marketing language", "current"},
	{"best practice", "Generic advice", "specific control"},
	{"it is recommended", "Generic advice", "implement"},
	{"consider implementing", "Generic advice", "implement"},
	{"may want to", "Generic advice", "should"},
	{"should consider", "Generic advice", "must"},
	{"AI-powered", "AI fluff", "automated"},
	{"machine learning", "AI fluff", "algorithm"},
	{"intelligent", "AI fluff", "automated"},
	{"smart", "AI fluff", "automated"},
	{"proactive", "Generic advice", "preventive"},
	{"ensure that", "Wordy", "verify"},
	{"in order to", "Wordy", "to"},
	{"due to the fact that", "Wordy", "because"},
	{"at this point in time", "Wordy", "now"},
	{"in the event that", "Wordy", "if"},
	{"for the purpose of", "Wordy", "for"},
	{"with regard to", "Wordy", "regarding"},
	{"it is important to note", "Filler", ""},
	{"it should be noted", "Filler", ""},
	{"please note", "Filler", ""},
	{"as mentioned above", "Filler", ""},
	{"as discussed", "Filler", ""},
}
