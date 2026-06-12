package intelligence

import (
	"fmt"
	"strings"
)

// ExplainabilityEngine generates human-readable explanations for assumptions.
type ExplainabilityEngine struct {
	taxonomy *TaxonomyEngine
}

// NewExplainabilityEngine creates an explainability engine.
func NewExplainabilityEngine(taxonomy *TaxonomyEngine) *ExplainabilityEngine {
	return &ExplainabilityEngine{taxonomy: taxonomy}
}

// ExplainWhy generates a human-readable WHY for each assumption.
func (ee *ExplainabilityEngine) ExplainWhy(a *Assumption, arch *ArchDescription) {
	if a == nil {
		return
	}

	// Build evidence list
	var evidence []string
	if arch != nil {
		evidence = ee.gatherEvidence(a, arch)
	}

	// Build missing controls list
	missingControls := ee.identifyMissingControls(a)

	// Build architecture context
	context := ee.buildArchitectureContext(a, arch)

	// Build confidence and category
	confidence := a.Confidence
	category := a.Category

	// Generate the human-readable WHY string
	why := ee.generateWhyString(a, evidence, missingControls, context, confidence, category)

	// Update the assumption
	a.Rationale = why
	if a.EvidenceSources == nil {
		a.EvidenceSources = evidence
	} else {
		a.EvidenceSources = append(a.EvidenceSources, evidence...)
	}
}

// ExplainAll runs ExplainWhy on all assumptions.
func (ee *ExplainabilityEngine) ExplainAll(assumptions []Assumption, arch *ArchDescription) []Assumption {
	var result []Assumption
	for _, a := range assumptions {
		copyA := a
		ee.ExplainWhy(&copyA, arch)
		result = append(result, copyA)
	}
	return result
}

// gatherEvidence finds components, relationships, and text evidence.
func (ee *ExplainabilityEngine) gatherEvidence(a *Assumption, arch *ArchDescription) []string {
	var evidence []string
	lowerDesc := strings.ToLower(a.Description)
	lowerKeywords := strings.ToLower(strings.Join(a.Keywords, " "))
	searchText := lowerDesc + " " + lowerKeywords

	for _, comp := range arch.Components {
		label := strings.ToLower(comp.Label)
		if strings.Contains(searchText, label) {
			evidence = append(evidence, fmt.Sprintf("component: %s", comp.Label))
		}
	}
	for _, rel := range arch.Relationships {
		if strings.Contains(searchText, strings.ToLower(rel.Source)) || strings.Contains(searchText, strings.ToLower(rel.Target)) || strings.Contains(searchText, strings.ToLower(rel.Label)) {
			evidence = append(evidence, fmt.Sprintf("relationship: %s → %s", rel.Source, rel.Target))
		}
	}
	return evidence
}

// identifyMissingControls extracts what controls are missing from the assumption.
func (ee *ExplainabilityEngine) identifyMissingControls(a *Assumption) []string {
	var missing []string
	desc := strings.ToLower(a.Description)

	if strings.Contains(desc, "does not specify") || strings.Contains(desc, "not specified") || strings.Contains(desc, "not documented") {
		// Extract the noun phrase after "does not specify"
		idx := strings.Index(desc, "does not specify")
		if idx >= 0 {
			rest := desc[idx+len("does not specify"):]
			rest = strings.TrimSpace(rest)
			// Take first few words
			words := strings.Fields(rest)
			if len(words) > 0 {
				phrase := strings.Join(words[:minInt(5, len(words))], " ")
				missing = append(missing, phrase)
			}
		}
	}
	if strings.Contains(desc, "missing") {
		idx := strings.Index(desc, "missing")
		if idx >= 0 {
			rest := desc[idx+len("missing"):]
			rest = strings.TrimSpace(rest)
			words := strings.Fields(rest)
			if len(words) > 0 {
				phrase := strings.Join(words[:minInt(5, len(words))], " ")
				missing = append(missing, phrase)
			}
		}
	}
	if strings.Contains(desc, "absence of") {
		idx := strings.Index(desc, "absence of")
		if idx >= 0 {
			rest := desc[idx+len("absence of"):]
			rest = strings.TrimSpace(rest)
			words := strings.Fields(rest)
			if len(words) > 0 {
				phrase := strings.Join(words[:minInt(5, len(words))], " ")
				missing = append(missing, phrase)
			}
		}
	}

	// Fallback: if no missing controls found, infer from category
	if len(missing) == 0 {
		cat := ee.taxonomy.GetCategory(a.Category)
		if cat != nil && len(cat.VerificationRules) > 0 {
			missing = append(missing, cat.VerificationRules[0])
		}
	}
	return missing
}

// buildArchitectureContext generates context about the architecture.
func (ee *ExplainabilityEngine) buildArchitectureContext(a *Assumption, arch *ArchDescription) string {
	if arch == nil {
		return ""
	}
	var parts []string
	if arch.Name != "" {
		parts = append(parts, fmt.Sprintf("architecture: %s", arch.Name))
	}
	if a.Component != "" {
		parts = append(parts, fmt.Sprintf("component: %s", a.Component))
	}
	if len(arch.Components) > 0 {
		parts = append(parts, fmt.Sprintf("total components: %d", len(arch.Components)))
	}
	if len(arch.Relationships) > 0 {
		parts = append(parts, fmt.Sprintf("total relationships: %d", len(arch.Relationships)))
	}
	return strings.Join(parts, "; ")
}

// generateWhyString creates the final human-readable WHY string.
func (ee *ExplainabilityEngine) generateWhyString(a *Assumption, evidence, missingControls []string, context string, confidence float64, category string) string {
	var parts []string

	// Core statement
	parts = append(parts, fmt.Sprintf("Why this assumption exists: %s", a.Description))

	// Evidence
	if len(evidence) > 0 {
		parts = append(parts, fmt.Sprintf("Evidence: %s", strings.Join(evidence, "; ")))
	}

	// Missing controls
	if len(missingControls) > 0 {
		parts = append(parts, fmt.Sprintf("Missing controls: %s", strings.Join(missingControls, "; ")))
	}

	// Architecture context
	if context != "" {
		parts = append(parts, fmt.Sprintf("Architecture context: %s", context))
	}

	// Confidence and category
	parts = append(parts, fmt.Sprintf("Confidence: %.0f%%; Category: %s", confidence*100, category))

	return strings.Join(parts, " | ")
}

// BuildSummaryExplanation returns a high-level explanation of the analysis.
func (ee *ExplainabilityEngine) BuildSummaryExplanation(assumptions []Assumption, contradictions []Contradiction, boundaries []TrustBoundary, domain string) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("Generated %d assumptions from architecture analysis", len(assumptions)))
	if domain != "" {
		parts = append(parts, fmt.Sprintf("Detected domain: %s", domain))
	}
	if len(contradictions) > 0 {
		parts = append(parts, fmt.Sprintf("Found %d contradictions requiring resolution", len(contradictions)))
	}
	if len(boundaries) > 0 {
		parts = append(parts, fmt.Sprintf("Discovered %d trust boundaries", len(boundaries)))
	}
	return strings.Join(parts, "; ")
}
