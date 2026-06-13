package narrative

import (
	"fmt"
	"strings"
)

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
	{
		var parts []string
		critical, high, medium := 0, 0, 0
		for _, a := range assumptions {
			switch a.Risk {
			case "Critical":
				critical++
			case "High":
				high++
			case "Medium":
				medium++
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
		overview.Summary = strings.Join(parts, " ")
	}

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
