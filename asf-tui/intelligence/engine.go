package intelligence

import (
	"fmt"
	"time"
)

// IntelligenceResult holds the full output of the intelligence engine.
type IntelligenceResult struct {
	Assumptions        []Assumption
	Contradictions     []Contradiction
	TrustBoundaries    []TrustBoundary
	Domain             string
	QualityScores      map[string]QualityScore
	Summary            string
	Explainability     string
	Controls           []ControlDetail
	Compliance         []string
	CriticalCount      int
	HighCount          int
	MediumCount        int
	LowCount           int
	TotalAssumptions   int
	StrideDistribution map[StrideCategory]int
}

// IntelligenceEngine composes all sub-engines and orchestrates the analysis.
type IntelligenceEngine struct {
	taxonomy       *TaxonomyEngine
	reasoning      *ReasoningEngine
	domainEngine   *DomainEngine
	contradiction  *ContradictionEngine
	trustBoundary  *TrustBoundaryEngine
	quality        *QualityEngine
	explainability *ExplainabilityEngine
}

// NewIntelligenceEngine creates a new intelligence engine.
func NewIntelligenceEngine() *IntelligenceEngine {
	te := NewTaxonomyEngine()
	return &IntelligenceEngine{
		taxonomy:       te,
		reasoning:      nil,
		domainEngine:   NewDomainEngine(),
		contradiction:  NewContradictionEngine(),
		trustBoundary:  NewTrustBoundaryEngine(),
		quality:        NewQualityEngine(),
		explainability: NewExplainabilityEngine(te),
	}
}

// Run executes the full intelligence pipeline on the given architecture.
func (ie *IntelligenceEngine) Run(arch *ArchDescription) *IntelligenceResult {
	result := &IntelligenceResult{
		QualityScores:      make(map[string]QualityScore),
		StrideDistribution: make(map[StrideCategory]int),
	}

	// Phase 1: Detect domain
	result.Domain = ie.domainEngine.DetectDomain(arch)

	// Phase 2: Apply domain-specific assumptions
	var allAssumptions []Assumption
	if result.Domain != "" {
		domainAssumptions := ie.domainEngine.ApplyDomainPack(result.Domain, arch)
		allAssumptions = append(allAssumptions, domainAssumptions...)
		result.Compliance = ie.domainEngine.GetPack(result.Domain).ComplianceMappings
		result.Controls = ie.domainEngine.GetPack(result.Domain).Controls
	}

	// Phase 3: Topological reasoning
	ie.reasoning = NewReasoningEngine(arch)
	inferred := ie.reasoning.InferAllAssumptions()
	allAssumptions = append(allAssumptions, inferred...)

	// Phase 4: Trust boundary discovery
	boundaries := ie.trustBoundary.DiscoverBoundaries(arch)
	result.TrustBoundaries = boundaries
	boundaryAssumptions := ie.trustBoundary.GenerateAssumptions(boundaries)
	allAssumptions = append(allAssumptions, boundaryAssumptions...)

	// Phase 5: Explainability
	allAssumptions = ie.explainability.ExplainAll(allAssumptions, arch)

	// Phase 6: Quality scoring
	for i := range allAssumptions {
		qs := ie.quality.Score(allAssumptions[i], arch)
		result.QualityScores[allAssumptions[i].ID] = qs
	}
	allAssumptions = ie.quality.Rank(allAssumptions, arch)

	// Phase 7: Contradiction detection
	result.Contradictions = ie.contradiction.DetectContradictions(allAssumptions)

	// Phase 8: Risk counting and STRIDE distribution
	for _, a := range allAssumptions {
		switch a.Risk {
		case RiskCritical:
			result.CriticalCount++
		case RiskHigh:
			result.HighCount++
		case RiskMedium:
			result.MediumCount++
		case RiskLow:
			result.LowCount++
		}
		for _, s := range a.Stride {
			result.StrideDistribution[s]++
		}
	}
	result.TotalAssumptions = len(allAssumptions)
	result.Assumptions = allAssumptions

	// Phase 9: Summary generation
	result.Explainability = ie.explainability.BuildSummaryExplanation(
		allAssumptions,
		result.Contradictions,
		result.TrustBoundaries,
		result.Domain,
	)
	result.Summary = ie.buildSummary(result, arch)

	return result
}

// RunWithExistingAssumptions runs the intelligence pipeline and merges with existing assumptions.
func (ie *IntelligenceEngine) RunWithExistingAssumptions(arch *ArchDescription, existing []Assumption) *IntelligenceResult {
	result := ie.Run(arch)
	merged := deduplicateAssumptions(append(existing, result.Assumptions...))
	result.Assumptions = merged
	result.TotalAssumptions = len(merged)

	// Re-run contradiction detection on merged set
	result.Contradictions = ie.contradiction.DetectContradictions(merged)

	// Re-score quality on merged set
	result.QualityScores = make(map[string]QualityScore)
	for i := range merged {
		qs := ie.quality.Score(merged[i], arch)
		merged[i].QualityScore = qs.Total
		result.QualityScores[merged[i].ID] = qs
	}
	merged = ie.quality.Rank(merged, arch)

	// Re-count risks after merge
	result.CriticalCount = 0
	result.HighCount = 0
	result.MediumCount = 0
	result.LowCount = 0
	result.StrideDistribution = make(map[StrideCategory]int)
	for _, a := range merged {
		switch a.Risk {
		case RiskCritical:
			result.CriticalCount++
		case RiskHigh:
			result.HighCount++
		case RiskMedium:
			result.MediumCount++
		case RiskLow:
			result.LowCount++
		}
		for _, s := range a.Stride {
			result.StrideDistribution[s]++
		}
	}
	return result
}

// buildSummary generates a human-readable summary of the analysis.
func (ie *IntelligenceEngine) buildSummary(result *IntelligenceResult, arch *ArchDescription) string {
	if arch == nil {
		arch = &ArchDescription{Name: "unknown"}
	}
	return fmt.Sprintf(
		"Intelligence analysis of %s completed at %s: %d assumptions (%d critical, %d high, %d medium, %d low), %d contradictions, %d trust boundaries, domain: %s",
		arch.Name,
		time.Now().Format("2006-01-02 15:04:05"),
		result.TotalAssumptions,
		result.CriticalCount,
		result.HighCount,
		result.MediumCount,
		result.LowCount,
		len(result.Contradictions),
		len(result.TrustBoundaries),
		result.Domain,
	)
}

// GetAssumptionsByCategory returns assumptions filtered by category.
func (ie *IntelligenceEngine) GetAssumptionsByCategory(result *IntelligenceResult, category string) []Assumption {
	var filtered []Assumption
	for _, a := range result.Assumptions {
		if a.Category == category {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

// GetAssumptionsByRisk returns assumptions filtered by risk level.
func (ie *IntelligenceEngine) GetAssumptionsByRisk(result *IntelligenceResult, risk RiskLevel) []Assumption {
	var filtered []Assumption
	for _, a := range result.Assumptions {
		if a.Risk == risk {
			filtered = append(filtered, a)
		}
	}
	return filtered
}

// GetContradictionsBySeverity returns contradictions filtered by severity.
func (ie *IntelligenceEngine) GetContradictionsBySeverity(result *IntelligenceResult, severity RiskLevel) []Contradiction {
	var filtered []Contradiction
	for _, c := range result.Contradictions {
		if c.Severity == severity {
			filtered = append(filtered, c)
		}
	}
	return filtered
}

// GetTopQualityAssumptions returns the top N assumptions by quality.
func (ie *IntelligenceEngine) GetTopQualityAssumptions(result *IntelligenceResult, n int) []Assumption {
	return ie.quality.GetTopAssumptions(result.Assumptions, nil, n)
}

// ValidateResult performs basic sanity checks on the result.
func (ie *IntelligenceEngine) ValidateResult(result *IntelligenceResult) []string {
	var issues []string
	if result.TotalAssumptions == 0 {
		issues = append(issues, "no assumptions generated")
	}
	if result.CriticalCount > result.TotalAssumptions {
		issues = append(issues, "critical count exceeds total assumptions")
	}
	if len(result.Contradictions) > result.TotalAssumptions {
		issues = append(issues, "contradiction count exceeds total assumptions")
	}
	return issues
}
