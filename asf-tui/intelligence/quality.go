package intelligence

import (
	"fmt"
	"sort"
	"strings"
)

// QualityEngine scores assumptions based on quality dimensions.
type QualityEngine struct{}

// NewQualityEngine creates a quality engine.
func NewQualityEngine() *QualityEngine {
	return &QualityEngine{}
}

// Score computes a quality score for an assumption.
func (qe *QualityEngine) Score(a Assumption, arch *ArchDescription) QualityScore {
	qs := QualityScore{}

	// Hiddenness: assumptions that are not obvious from component names score higher
	qs.Hiddenness = qe.scoreHiddenness(a, arch)

	// Impact: higher risk assumptions score higher
	qs.Impact = qe.scoreImpact(a)

	// Novelty: domain-specific or inferred assumptions score higher than generic ones
	qs.Novelty = qe.scoreNovelty(a)

	// ArchitecturalRelevance: assumptions that map to actual components score higher
	qs.ArchitecturalRelevance = qe.scoreArchitecturalRelevance(a, arch)

	// Risk: convert risk level to score
	qs.Risk = qe.scoreRiskLevel(a.Risk)

	// Confidence: higher confidence assumptions score higher
	qs.Confidence = a.Confidence

	// Total is a weighted composite
	qs.Total = qs.Hiddenness*0.20 + qs.Impact*0.20 + qs.Novelty*0.20 + qs.ArchitecturalRelevance*0.15 + qs.Risk*0.15 + qs.Confidence*0.10

	qs.Reason = fmt.Sprintf("hiddenness=%.2f impact=%.2f novelty=%.2f relevance=%.2f risk=%.2f confidence=%.2f",
		qs.Hiddenness, qs.Impact, qs.Novelty, qs.ArchitecturalRelevance, qs.Risk, qs.Confidence)
	return qs
}

// Rank sorts assumptions by quality score (highest first).
func (qe *QualityEngine) Rank(assumptions []Assumption, arch *ArchDescription) []Assumption {
	scored := make([]struct {
		assumption Assumption
		score      QualityScore
	}, len(assumptions))
	for i, a := range assumptions {
		scored[i] = struct {
			assumption Assumption
			score      QualityScore
		}{a, qe.Score(a, arch)}
	}

	sort.Slice(scored, func(i, j int) bool {
		return scored[i].score.Total > scored[j].score.Total
	})

	var result []Assumption
	for _, s := range scored {
		result = append(result, s.assumption)
	}
	return result
}

// scoreHiddenness evaluates how hidden/non-obvious an assumption is.
func (qe *QualityEngine) scoreHiddenness(a Assumption, arch *ArchDescription) float64 {
	// Generic link-encryption assumptions are obvious and score low
	desc := strings.ToLower(a.Description)
	if strings.Contains(desc, "tls") && strings.Contains(desc, "all communication") {
		return 0.2
	}
	if strings.Contains(desc, "encryption in transit") && strings.Contains(desc, "all components") {
		return 0.2
	}
	if strings.Contains(desc, "secure protocol") {
		return 0.2
	}

	// Domain-specific assumptions score higher
	if a.SourceType == "domain-inferred" || a.SourceType == "inferred" {
		return 0.9
	}

	// Inferred assumptions from topology score high
	if strings.HasPrefix(a.ID, "INF-") {
		return 0.85
	}

	// If it mentions missing controls, it's hidden
	if strings.Contains(desc, "does not specify") || strings.Contains(desc, "not specified") || strings.Contains(desc, "not documented") {
		return 0.8
	}

	// If it maps to actual components but is not a direct keyword match
	if arch != nil && len(a.SourceComponents) > 0 {
		return 0.7
	}

	return 0.5
}

// scoreImpact evaluates the impact dimension.
func (qe *QualityEngine) scoreImpact(a Assumption) float64 {
	// Use impact field directly, normalized to 0-1
	score := float64(a.Impact) / 5.0
	if score < 0 {
		score = 0
	}
	if score > 1 {
		score = 1
	}
	return score
}

// scoreNovelty evaluates how novel/domain-specific the assumption is.
func (qe *QualityEngine) scoreNovelty(a Assumption) float64 {
	desc := strings.ToLower(a.Description)
	// Domain-specific terms boost novelty
	domainTerms := []string{"phi", "hipaa", "patient", "pci", "payment", "tenant", "fido2", "kubernetes", "pod", "container", "vpn", "federation", "data lineage", "data governance"}
	for _, term := range domainTerms {
		if strings.Contains(desc, term) {
			return 0.9
		}
	}
	// Inferred assumptions are novel
	if strings.HasPrefix(a.ID, "INF-") || strings.HasPrefix(a.ID, "DOM-") || strings.HasPrefix(a.ID, "TB-") {
		return 0.85
	}
	// Generic assumptions score low
	if strings.Contains(desc, "all components") || strings.Contains(desc, "standard policy") || strings.Contains(desc, "all communication") || strings.Contains(desc, "tls encryption") || strings.Contains(desc, "encryption in transit") {
		return 0.2
	}
	return 0.5
}

// scoreArchitecturalRelevance evaluates how relevant the assumption is to the architecture.
func (qe *QualityEngine) scoreArchitecturalRelevance(a Assumption, arch *ArchDescription) float64 {
	if arch == nil {
		return 0.3
	}
	// Check if assumption component exists in architecture
	if a.Component != "" {
		for _, comp := range arch.Components {
			if strings.Contains(strings.ToLower(comp.Label), strings.ToLower(a.Component)) {
				return 0.9
			}
		}
	}
	// Check if keywords match components
	matches := 0
	for _, kw := range a.Keywords {
		for _, comp := range arch.Components {
			if strings.Contains(strings.ToLower(comp.Label), strings.ToLower(kw)) {
				matches++
				break
			}
		}
	}
	if len(a.Keywords) > 0 {
		ratio := float64(matches) / float64(len(a.Keywords))
		if ratio > 0.5 {
			return 0.8
		}
		if ratio > 0.2 {
			return 0.6
		}
	}
	return 0.4
}

// scoreRiskLevel converts risk level to a 0-1 score.
func (qe *QualityEngine) scoreRiskLevel(r RiskLevel) float64 {
	switch r {
	case RiskCritical:
		return 1.0
	case RiskHigh:
		return 0.8
	case RiskMedium:
		return 0.5
	case RiskLow:
		return 0.2
	default:
		return 0.5
	}
}

// GetTopAssumptions returns the top N assumptions by quality.
func (qe *QualityEngine) GetTopAssumptions(assumptions []Assumption, arch *ArchDescription, n int) []Assumption {
	ranked := qe.Rank(assumptions, arch)
	if n > len(ranked) {
		n = len(ranked)
	}
	return ranked[:n]
}

// AverageQuality computes the average quality score across all assumptions.
func (qe *QualityEngine) AverageQuality(assumptions []Assumption, arch *ArchDescription) float64 {
	if len(assumptions) == 0 {
		return 0
	}
	var sum float64
	for _, a := range assumptions {
		sum += qe.Score(a, arch).Total
	}
	return sum / float64(len(assumptions))
}

// QualityReport generates a summary string of quality scores.
func (qe *QualityEngine) QualityReport(assumptions []Assumption, arch *ArchDescription) string {
	if len(assumptions) == 0 {
		return "no assumptions to score"
	}
	avg := qe.AverageQuality(assumptions, arch)
	top := qe.GetTopAssumptions(assumptions, arch, 3)
	var topIDs []string
	for _, a := range top {
		topIDs = append(topIDs, a.ID)
	}
	return fmt.Sprintf("average quality %.2f across %d assumptions; top: %s", avg, len(assumptions), strings.Join(topIDs, ", "))
}
