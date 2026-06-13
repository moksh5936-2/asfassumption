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
