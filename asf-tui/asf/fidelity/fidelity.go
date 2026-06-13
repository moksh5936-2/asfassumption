package fidelity

import (
	"asf-tui/asf/fact"
	"fmt"
)

// FidelityScore represents the architectural fidelity score.
type FidelityScore struct {
	TotalFacts            int     `json:"total_facts"`
	RespectedFacts        int     `json:"respected_facts"`
	ContradictedFacts     int     `json:"contradicted_facts"`
	IgnoredFacts          int     `json:"ignored_facts"`
	UnmappedFacts         int     `json:"unmapped_facts"`
	Score                 float64 `json:"score"`
	AssumptionQuality     float64 `json:"assumption_quality"`
	ContradictionAccuracy float64 `json:"contradiction_accuracy"`
	NoveltyScore          float64 `json:"novelty_score"`
	Overall               string  `json:"overall"`
}

// FidelityScorer computes the architectural fidelity score.
type FidelityScorer struct {
	factInventory *fact.Inventory
}

// NewFidelityScorer creates a new fidelity scorer.
func NewFidelityScorer(inventory *fact.Inventory) *FidelityScorer {
	return &FidelityScorer{factInventory: inventory}
}

// Compute computes the fidelity score.
func (fs *FidelityScorer) Compute(assumptions []HiddenAssumption, contradictions []RealContradiction, traceability []TraceabilityRecord) FidelityScore {
	total := len(fs.factInventory.Facts)
	if total == 0 {
		return FidelityScore{
			TotalFacts: 0,
			Score:      1.0,
			Overall:    "PASS",
		}
	}

	respected := 0
	contradicted := 0
	ignored := 0
	unmapped := 0

	for _, f := range fs.factInventory.Facts {
		// Check if any assumption contradicts this fact
		factContradicted := false
		for _, c := range contradictions {
			if c.FactA.ID == f.ID || (c.FactB != nil && c.FactB.ID == f.ID) {
				factContradicted = true
				break
			}
		}

		if factContradicted {
			contradicted++
			continue
		}

		// Check if the fact is referenced in traceability
		factReferenced := false
		for _, t := range traceability {
			if t.SourceFactID == f.ID {
				factReferenced = true
				break
			}
		}

		if factReferenced {
			respected++
		} else {
			unmapped++
		}
	}

	// Compute scores
	score := float64(respected) / float64(total)
	if contradicted > 0 {
		// Each contradiction reduces score
		score -= float64(contradicted) * 0.1
	}
	if score < 0 {
		score = 0
	}

	// Assumption quality
	assumptionQuality := fs.computeAssumptionQuality(assumptions)

	// Contradiction accuracy
	contradictionAccuracy := fs.computeContradictionAccuracy(contradictions)

	// Novelty
	novelty := fs.computeNovelty(assumptions)

	overall := "NOT_CERTIFIED"
	if score >= 0.9 && assumptionQuality >= 0.7 && contradictionAccuracy >= 0.9 && novelty >= 0.6 {
		overall = "ARCHITECTURAL_FIDELITY_CERTIFIED"
	} else if score >= 0.7 {
		overall = "CONDITIONAL"
	}

	return FidelityScore{
		TotalFacts:            total,
		RespectedFacts:        respected,
		ContradictedFacts:     contradicted,
		IgnoredFacts:          ignored,
		UnmappedFacts:         unmapped,
		Score:                 score,
		AssumptionQuality:     assumptionQuality,
		ContradictionAccuracy: contradictionAccuracy,
		NoveltyScore:          novelty,
		Overall:               overall,
	}
}

// computeAssumptionQuality computes the average quality score.
func (fs *FidelityScorer) computeAssumptionQuality(assumptions []HiddenAssumption) float64 {
	if len(assumptions) == 0 {
		return 0
	}

	total := 0.0
	for _, a := range assumptions {
		total += a.QualityScore
	}

	return total / float64(len(assumptions))
}

// computeContradictionAccuracy computes the accuracy of contradictions.
func (fs *FidelityScorer) computeContradictionAccuracy(contradictions []RealContradiction) float64 {
	if len(contradictions) == 0 {
		return 1.0
	}

	// All contradictions should be real (fact-fact or fact-assumption)
	valid := 0
	for _, c := range contradictions {
		if c.Type == "fact-fact" || c.Type == "fact-assumption" {
			valid++
		}
	}

	return float64(valid) / float64(len(contradictions))
}

// computeNovelty computes the average novelty score.
func (fs *FidelityScorer) computeNovelty(assumptions []HiddenAssumption) float64 {
	if len(assumptions) == 0 {
		return 0
	}

	total := 0.0
	for _, a := range assumptions {
		total += a.NoveltyScore
	}

	return total / float64(len(assumptions))
}

// FormatReport formats the fidelity report.
func (fs *FidelityScorer) FormatReport(score FidelityScore) string {
	var report string
	report += fmt.Sprintf("# Architectural Fidelity Report\n\n")
	report += fmt.Sprintf("## Score\n\n")
	report += fmt.Sprintf("- **Fidelity Score**: %.1f%%\n", score.Score*100)
	report += fmt.Sprintf("- **Assumption Quality**: %.1f%%\n", score.AssumptionQuality*100)
	report += fmt.Sprintf("- **Contradiction Accuracy**: %.1f%%\n", score.ContradictionAccuracy*100)
	report += fmt.Sprintf("- **Novelty Score**: %.1f%%\n", score.NoveltyScore*100)
	report += fmt.Sprintf("- **Overall**: %s\n\n", score.Overall)

	report += fmt.Sprintf("## Facts\n\n")
	report += fmt.Sprintf("- **Total Facts**: %d\n", score.TotalFacts)
	report += fmt.Sprintf("- **Respected**: %d\n", score.RespectedFacts)
	report += fmt.Sprintf("- **Contradicted**: %d\n", score.ContradictedFacts)
	report += fmt.Sprintf("- **Unmapped**: %d\n\n", score.UnmappedFacts)

	report += fmt.Sprintf("## Thresholds\n\n")
	report += fmt.Sprintf("- **Goal**: 90%%+\n")
	report += fmt.Sprintf("- **Current**: %.1f%%\n", score.Score*100)

	return report
}
