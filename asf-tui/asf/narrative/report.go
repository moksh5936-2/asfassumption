package narrative

import (
	"fmt"
	"sort"
	"strings"
)

// generateExecutiveReport creates the C-level executive report.
func (e *NarrativeEngine) generateExecutiveReport(
	assumptions []Assumption,
	trustBoundaries []TrustBoundary,
	contradictions []Contradiction,
) ExecutiveReport {
	report := ExecutiveReport{}

	// Architecture overview
	report.ArchitectureOverview = e.generateExecutiveOverview(assumptions, trustBoundaries)

	// Key assumptions (all assumptions that are not low risk)
	var keyAssumptions []ExecutiveAssumption
	for _, a := range assumptions {
		if a.Risk != "Low" {
			ea := ExecutiveAssumption{
				Text:        a.Description,
				RiskLevel:   a.Risk,
				Consequence: e.inferConsequence(a),
			}
			keyAssumptions = append(keyAssumptions, ea)
		}
	}
	report.KeyAssumptions = keyAssumptions

	// Most critical assumptions (Critical and High)
	var criticalAssumptions []ExecutiveAssumption
	for _, a := range assumptions {
		if a.Risk == "Critical" || a.Risk == "High" {
			ea := ExecutiveAssumption{
				Text:           a.Description,
				RiskLevel:      a.Risk,
				Consequence:    e.inferConsequence(a),
				BusinessImpact: e.inferBusinessImpact(a),
			}
			criticalAssumptions = append(criticalAssumptions, ea)
		}
	}
	report.MostCriticalAssumptions = criticalAssumptions

	// High impact consequences
	report.HighImpactConsequences = e.generateHighImpactConsequences(assumptions)

	// Trust dependencies
	report.TrustDependencies = e.generateTrustDependencies(trustBoundaries)

	// Single points of failure
	report.SinglePointsOfFailure = e.generateSinglePointsOfFailure(assumptions)

	// Architectural concerns
	report.ArchitecturalConcerns = e.generateArchitecturalConcerns(assumptions, contradictions)

	// Recommended investments
	report.RecommendedInvestments = e.generateRecommendedInvestments(assumptions)

	// Enforce style
	report.ArchitectureOverview = e.enforceStyle(report.ArchitectureOverview)
	for i := range report.HighImpactConsequences {
		report.HighImpactConsequences[i] = e.enforceStyle(report.HighImpactConsequences[i])
	}
	for i := range report.ArchitecturalConcerns {
		report.ArchitecturalConcerns[i] = e.enforceStyle(report.ArchitecturalConcerns[i])
	}
	for i := range report.RecommendedInvestments {
		report.RecommendedInvestments[i] = e.enforceStyle(report.RecommendedInvestments[i])
	}

	return report
}

// generateExecutiveOverview creates the executive-level architecture overview.
func (e *NarrativeEngine) generateExecutiveOverview(assumptions []Assumption, trustBoundaries []TrustBoundary) string {
	var parts []string

	// Count risk levels
	critical, high, medium, low := 0, 0, 0, 0
	for _, a := range assumptions {
		switch a.Risk {
		case "Critical":
			critical++
		case "High":
			high++
		case "Medium":
			medium++
		case "Low":
			low++
		}
	}

	parts = append(parts, fmt.Sprintf("The analysis identified %d assumptions.", len(assumptions)))

	if critical > 0 {
		parts = append(parts, fmt.Sprintf("%d are critical and require immediate attention.", critical))
	}
	if high > 0 {
		parts = append(parts, fmt.Sprintf("%d are high risk.", high))
	}
	if medium > 0 {
		parts = append(parts, fmt.Sprintf("%d are medium risk.", medium))
	}

	if len(trustBoundaries) > 0 {
		parts = append(parts, fmt.Sprintf("The architecture has %d trust boundaries that require review.", len(trustBoundaries)))
	}

	return strings.Join(parts, " ")
}

// generateHighImpactConsequences extracts the highest-impact consequences.
func (e *NarrativeEngine) generateHighImpactConsequences(assumptions []Assumption) []string {
	var consequences []string
	seen := make(map[string]bool)

	for _, a := range assumptions {
		if a.Risk == "Critical" || a.Risk == "High" {
			c := e.inferConsequence(a)
			if !seen[c] {
				seen[c] = true
				consequences = append(consequences, c)
			}
		}
	}

	// Limit to top 5
	if len(consequences) > 5 {
		consequences = consequences[:5]
	}

	return consequences
}

// generateTrustDependencies extracts trust dependencies from boundaries.
func (e *NarrativeEngine) generateTrustDependencies(trustBoundaries []TrustBoundary) []string {
	var deps []string
	seen := make(map[string]bool)

	for _, tb := range trustBoundaries {
		if tb.Description != "" && !seen[tb.Description] {
			seen[tb.Description] = true
			deps = append(deps, tb.Description)
		}
	}

	return deps
}

// generateSinglePointsOfFailure identifies single points of failure.
func (e *NarrativeEngine) generateSinglePointsOfFailure(assumptions []Assumption) []string {
	var spofs []string
	seen := make(map[string]bool)

	for _, a := range assumptions {
		if a.Risk == "Critical" {
			// Check if this is a single point of failure
			if isSinglePointOfFailure(a) {
				text := fmt.Sprintf("%s: %s", a.Component, a.Description)
				if !seen[text] {
					seen[text] = true
					spofs = append(spofs, text)
				}
			}
		}
	}

	return spofs
}

// isSinglePointOfFailure checks if an assumption represents a single point of failure.
func isSinglePointOfFailure(a Assumption) bool {
	desc := strings.ToLower(a.Description)
	keywords := []string{
		"single point", "single source", "only", "sole", "primary",
		"main", "central", "critical path", "bottleneck",
	}
	for _, kw := range keywords {
		if strings.Contains(desc, kw) {
			return true
		}
	}
	return false
}

// generateArchitecturalConcerns identifies architectural concerns.
func (e *NarrativeEngine) generateArchitecturalConcerns(assumptions []Assumption, contradictions []Contradiction) []string {
	var concerns []string
	seen := make(map[string]bool)

	// Contradictions are the highest concern
	for _, c := range contradictions {
		if c.Explanation != "" && !seen[c.Explanation] {
			seen[c.Explanation] = true
			concerns = append(concerns, c.Explanation)
		}
	}

	// Critical assumptions with specific patterns
	for _, a := range assumptions {
		if a.Risk == "Critical" {
			desc := strings.ToLower(a.Description)
			if strings.Contains(desc, "not implemented") || strings.Contains(desc, "missing") || strings.Contains(desc, "absence") {
				concern := fmt.Sprintf("Missing control: %s", a.Description)
				if !seen[concern] {
					seen[concern] = true
					concerns = append(concerns, concern)
				}
			}
		}
	}

	return concerns
}

// generateRecommendedInvestments generates investment recommendations.
func (e *NarrativeEngine) generateRecommendedInvestments(assumptions []Assumption) []string {
	var investments []string
	seen := make(map[string]bool)

	for _, a := range assumptions {
		if a.Risk == "Critical" || a.Risk == "High" {
			controls := e.inferControls(a)
			for _, c := range controls {
				if !seen[c] {
					seen[c] = true
					investments = append(investments, c)
				}
			}
		}
	}

	return investments
}

// inferBusinessImpact infers business impact from an assumption.
func (e *NarrativeEngine) inferBusinessImpact(a Assumption) string {
	desc := strings.ToLower(a.Description)

	if strings.Contains(desc, "authentication") || strings.Contains(desc, "mfa") {
		return "Unauthorized access to systems and data."
	}
	if strings.Contains(desc, "encrypt") {
		return "Data breach and regulatory non-compliance."
	}
	if strings.Contains(desc, "access") {
		return "Unauthorized data modification or exfiltration."
	}
	if strings.Contains(desc, "logging") {
		return "Inability to detect or respond to security incidents."
	}
	if strings.Contains(desc, "backup") {
		return "Permanent data loss and business continuity failure."
	}
	if strings.Contains(desc, "network") {
		return "Lateral movement and network compromise."
	}
	if strings.Contains(desc, "api") {
		return "API abuse and data exfiltration."
	}
	if strings.Contains(desc, "compliance") || strings.Contains(desc, "hipaa") || strings.Contains(desc, "pci") {
		return "Regulatory fines and legal liability."
	}

	return "Security control failure with operational impact."
}

// generateTechnicalSummary creates the technical summary.
func (e *NarrativeEngine) generateTechnicalSummary(
	assumptions []Assumption,
	depMap map[string][]string,
	strideDist map[string]int,
	riskDist map[string]int,
) TechnicalSummary {
	summary := TechnicalSummary{
		AssumptionDetails:  make([]TechnicalAssumption, 0, len(assumptions)),
		STRIDEDistribution: strideDist,
		RiskDistribution:   riskDist,
		Dependencies:       make([]TechnicalDependency, 0, len(assumptions)),
		Recommendations:    make([]TechnicalRecommendation, 0, len(assumptions)),
	}

	// Architecture summary
	summary.ArchitectureSummary = e.generateTechnicalArchitectureSummary(assumptions)

	// Per-assumption details
	for _, a := range assumptions {
		ta := TechnicalAssumption{
			ID:                a.ID,
			Description:       a.Description,
			Component:         a.Component,
			Category:          a.Category,
			RiskLevel:         a.Risk,
			STRIDECategories:  a.STRIDECategories,
			Likelihood:        a.Likelihood,
			Impact:            a.Impact,
			Confidence:        a.Confidence,
			EvidenceSources:   a.EvidenceSources,
			Rationale:         a.Rationale,
			DownstreamSystems: depMap[a.ID],
			FailureScenario:   e.inferConsequence(a),
			Recommendation:    e.generateRecommendation(a),
		}
		summary.AssumptionDetails = append(summary.AssumptionDetails, ta)
	}

	// Dependencies
	for _, a := range assumptions {
		if len(depMap[a.ID]) > 0 {
			td := TechnicalDependency{
				AssumptionID:        a.ID,
				AssumptionText:      a.Description,
				DependentComponents: depMap[a.ID],
				DependencyType:      "architectural",
			}
			summary.Dependencies = append(summary.Dependencies, td)
		}
	}

	// Recommendations
	for _, a := range assumptions {
		if a.Risk == "Critical" || a.Risk == "High" {
			controls := e.inferControls(a)
			for _, c := range controls {
				tr := TechnicalRecommendation{
					AssumptionID:         a.ID,
					AssumptionText:       a.Description,
					Recommendation:       c,
					Priority:             a.Risk,
					ImplementationEffort: e.inferEffort(a, c),
					MitigatesSTRIDE:      a.STRIDECategories,
				}
				summary.Recommendations = append(summary.Recommendations, tr)
			}
		}
	}

	// Sort recommendations by priority
	sort.Slice(summary.Recommendations, func(i, j int) bool {
		priorityOrder := map[string]int{"Critical": 0, "High": 1, "Medium": 2, "Low": 3}
		return priorityOrder[summary.Recommendations[i].Priority] < priorityOrder[summary.Recommendations[j].Priority]
	})

	return summary
}

// generateTechnicalArchitectureSummary creates the technical architecture summary.
func (e *NarrativeEngine) generateTechnicalArchitectureSummary(assumptions []Assumption) string {
	var parts []string

	critical, high := 0, 0
	for _, a := range assumptions {
		if a.Risk == "Critical" {
			critical++
		} else if a.Risk == "High" {
			high++
		}
	}

	parts = append(parts, fmt.Sprintf("Analysis identified %d assumptions.", len(assumptions)))
	if critical > 0 {
		parts = append(parts, fmt.Sprintf("%d critical findings require immediate remediation.", critical))
	}
	if high > 0 {
		parts = append(parts, fmt.Sprintf("%d high-risk findings require planned remediation.", high))
	}

	return strings.Join(parts, " ")
}

// inferEffort estimates implementation effort.
func (e *NarrativeEngine) inferEffort(a Assumption, control string) string {
	if strings.Contains(control, "MFA") || strings.Contains(control, "multi-factor") {
		return "Low (configuration change)"
	}
	if strings.Contains(control, "encryption") {
		return "Medium (infrastructure change)"
	}
	if strings.Contains(control, "segmentation") || strings.Contains(control, "network") {
		return "High (architecture change)"
	}
	if strings.Contains(control, "backup") || strings.Contains(control, "recovery") {
		return "Medium (process and infrastructure)"
	}
	if strings.Contains(control, "audit") || strings.Contains(control, "logging") {
		return "Low (process and tooling)"
	}
	if strings.Contains(control, "WAF") || strings.Contains(control, "DDoS") {
		return "Medium (infrastructure deployment)"
	}
	return "Medium"
}

// generateArchitectureOverview creates the architecture overview.
func (e *NarrativeEngine) generateArchitectureOverview(
	archName string,
	assumptions []Assumption,
	controls []ControlDetail,
	trustBoundaries []TrustBoundary,
	domain string,
) ArchitectureOverview {
	overview := ArchitectureOverview{
		Name:             archName,
		Domain:           domain,
		TotalAssumptions: len(assumptions),
	}

	// Count risk levels
	for _, a := range assumptions {
		switch a.Risk {
		case "Critical":
			overview.CriticalCount++
		case "High":
			overview.HighCount++
		}
	}

	// Key components
	seen := make(map[string]bool)
	for _, a := range assumptions {
		if a.Component != "" && !seen[a.Component] {
			seen[a.Component] = true
			overview.KeyComponents = append(overview.KeyComponents, a.Component)
		}
	}
	overview.TotalComponents = len(overview.KeyComponents)

	// Trust dependencies
	seen = make(map[string]bool)
	for _, tb := range trustBoundaries {
		if tb.Description != "" && !seen[tb.Description] {
			seen[tb.Description] = true
			overview.TrustDependencies = append(overview.TrustDependencies, tb.Description)
		}
	}

	// Summary
	overview.Summary = e.generateExecutiveOverview(assumptions, trustBoundaries)

	return overview
}

// generateFullArchitectNarrative generates the complete architect narrative text.
func (e *NarrativeEngine) generateFullArchitectNarrative(
	overview ArchitectureOverview,
	narratives []AssumptionNarrative,
) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# Architecture Security Narrative: %s\n\n", overview.Name))
	b.WriteString(fmt.Sprintf("**Domain:** %s\n", overview.Domain))
	b.WriteString(fmt.Sprintf("**Components:** %d\n", overview.TotalComponents))
	b.WriteString(fmt.Sprintf("**Assumptions:** %d (Critical: %d, High: %d)\n\n",
		overview.TotalAssumptions, overview.CriticalCount, overview.HighCount))

	b.WriteString("## Architecture Overview\n\n")
	b.WriteString(overview.Summary + "\n\n")

	if len(overview.KeyComponents) > 0 {
		b.WriteString("**Key Components:** ")
		b.WriteString(strings.Join(overview.KeyComponents, ", "))
		b.WriteString("\n\n")
	}

	if len(overview.TrustDependencies) > 0 {
		b.WriteString("**Trust Dependencies:** ")
		b.WriteString(strings.Join(overview.TrustDependencies, ", "))
		b.WriteString("\n\n")
	}

	b.WriteString("## Assumption Narratives\n\n")

	for _, n := range narratives {
		b.WriteString(fmt.Sprintf("### %s\n\n", n.AssumptionText))
		b.WriteString(fmt.Sprintf("**Risk:** %s | **Confidence:** %.0f%%\n\n", n.RiskLevel, n.Confidence*100))

		if len(n.STRIDECategories) > 0 {
			b.WriteString(fmt.Sprintf("**STRIDE:** %s\n\n", strings.Join(n.STRIDECategories, ", ")))
		}

		b.WriteString("**Context**\n\n")
		b.WriteString(n.Context + "\n\n")

		b.WriteString("**Why ASF Identified This**\n\n")
		b.WriteString(n.WhyASFIdentifiedIt + "\n\n")

		b.WriteString("**Architectural Importance**\n\n")
		b.WriteString(n.ArchitecturalImportance + "\n\n")

		b.WriteString("**Failure Consequence**\n\n")
		b.WriteString(n.FailureConsequence + "\n\n")

		b.WriteString("**Security Recommendation**\n\n")
		b.WriteString(n.SecurityRecommendation + "\n\n")

		if len(n.DependsOn) > 0 {
			b.WriteString(fmt.Sprintf("**Depends On:** %s\n\n", strings.Join(n.DependsOn, ", ")))
		}
		if len(n.DownstreamImpact) > 0 {
			b.WriteString(fmt.Sprintf("**Downstream Impact:** %s\n\n", strings.Join(n.DownstreamImpact, ", ")))
		}

		b.WriteString("---\n\n")
	}

	return b.String()
}
