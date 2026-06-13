package narrative

import (
	"fmt"
	"strings"
)

// ExportMarkdown generates a Markdown executive report.
func ExportMarkdown(output *NarrativeOutput) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# Security Architect Narrative: %s\n\n", output.ArchitectureOverview.Name))
	b.WriteString(fmt.Sprintf("**Generated:** %s\n\n", output.GeneratedAt.Format("2006-01-02 15:04:05")))

	// Executive Summary
	b.WriteString("## Executive Summary\n\n")
	b.WriteString(output.ExecutiveReport.ArchitectureOverview + "\n\n")

	if len(output.ExecutiveReport.MostCriticalAssumptions) > 0 {
		b.WriteString("### Most Critical Assumptions\n\n")
		for _, a := range output.ExecutiveReport.MostCriticalAssumptions {
			b.WriteString(fmt.Sprintf("- **[%s]** %s\n", a.RiskLevel, a.Text))
			b.WriteString(fmt.Sprintf("  - Consequence: %s\n", a.Consequence))
			b.WriteString(fmt.Sprintf("  - Business Impact: %s\n", a.BusinessImpact))
		}
		b.WriteString("\n")
	}

	if len(output.ExecutiveReport.HighImpactConsequences) > 0 {
		b.WriteString("### High Impact Consequences\n\n")
		for _, c := range output.ExecutiveReport.HighImpactConsequences {
			b.WriteString(fmt.Sprintf("- %s\n", c))
		}
		b.WriteString("\n")
	}

	if len(output.ExecutiveReport.TrustDependencies) > 0 {
		b.WriteString("### Trust Dependencies\n\n")
		for _, d := range output.ExecutiveReport.TrustDependencies {
			b.WriteString(fmt.Sprintf("- %s\n", d))
		}
		b.WriteString("\n")
	}

	if len(output.ExecutiveReport.SinglePointsOfFailure) > 0 {
		b.WriteString("### Single Points of Failure\n\n")
		for _, s := range output.ExecutiveReport.SinglePointsOfFailure {
			b.WriteString(fmt.Sprintf("- %s\n", s))
		}
		b.WriteString("\n")
	}

	if len(output.ExecutiveReport.ArchitecturalConcerns) > 0 {
		b.WriteString("### Architectural Concerns\n\n")
		for _, c := range output.ExecutiveReport.ArchitecturalConcerns {
			b.WriteString(fmt.Sprintf("- %s\n", c))
		}
		b.WriteString("\n")
	}

	if len(output.ExecutiveReport.RecommendedInvestments) > 0 {
		b.WriteString("### Recommended Investments\n\n")
		for _, i := range output.ExecutiveReport.RecommendedInvestments {
			b.WriteString(fmt.Sprintf("- %s\n", i))
		}
		b.WriteString("\n")
	}

	// Technical Summary
	b.WriteString("## Technical Summary\n\n")
	b.WriteString(output.TechnicalSummary.ArchitectureSummary + "\n\n")

	if len(output.TechnicalSummary.Recommendations) > 0 {
		b.WriteString("### Prioritized Recommendations\n\n")
		b.WriteString("| Priority | Recommendation | Assumption | Effort |\n")
		b.WriteString("|----------|---------------|------------|--------|\n")
		for _, r := range output.TechnicalSummary.Recommendations {
			b.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				r.Priority, r.Recommendation, r.AssumptionText, r.ImplementationEffort))
		}
		b.WriteString("\n")
	}

	if len(output.TechnicalSummary.Dependencies) > 0 {
		b.WriteString("### Dependencies\n\n")
		for _, d := range output.TechnicalSummary.Dependencies {
			b.WriteString(fmt.Sprintf("- **%s** affects: %s\n", d.AssumptionText, strings.Join(d.DependentComponents, ", ")))
		}
		b.WriteString("\n")
	}

	// Architect Narrative
	b.WriteString("## Architect Narrative\n\n")
	b.WriteString(output.ArchitectNarrative)

	return b.String()
}

// ExportHTML generates an HTML executive report.
func ExportHTML(output *NarrativeOutput) string {
	var b strings.Builder

	b.WriteString(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Security Architect Narrative: ` + htmlEscape(output.ArchitectureOverview.Name) + `</title>
<style>
body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; line-height: 1.6; max-width: 900px; margin: 40px auto; padding: 20px; color: #333; }
h1 { border-bottom: 2px solid #2c3e50; padding-bottom: 10px; color: #2c3e50; }
h2 { color: #34495e; border-bottom: 1px solid #bdc3c7; padding-bottom: 5px; margin-top: 30px; }
h3 { color: #7f8c8d; }
.badge-critical { background: #e74c3c; color: white; padding: 2px 8px; border-radius: 3px; font-size: 12px; font-weight: bold; }
.badge-high { background: #e67e22; color: white; padding: 2px 8px; border-radius: 3px; font-size: 12px; font-weight: bold; }
.badge-medium { background: #f39c12; color: white; padding: 2px 8px; border-radius: 3px; font-size: 12px; font-weight: bold; }
.badge-low { background: #27ae60; color: white; padding: 2px 8px; border-radius: 3px; font-size: 12px; font-weight: bold; }
table { border-collapse: collapse; width: 100%; margin: 20px 0; }
th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
th { background: #f8f9fa; font-weight: 600; }
tr:nth-child(even) { background: #f8f9fa; }
ul { padding-left: 20px; }
li { margin: 5px 0; }
.section { margin: 30px 0; }
.assumption { border-left: 3px solid #3498db; padding-left: 15px; margin: 20px 0; }
.assumption-critical { border-left-color: #e74c3c; }
.assumption-high { border-left-color: #e67e22; }
.assumption-medium { border-left-color: #f39c12; }
</style>
</head>
<body>
`)

	b.WriteString(fmt.Sprintf("<h1>Security Architect Narrative: %s</h1>\n", htmlEscape(output.ArchitectureOverview.Name)))
	b.WriteString(fmt.Sprintf("<p><strong>Generated:</strong> %s</p>\n", output.GeneratedAt.Format("2006-01-02 15:04:05")))

	// Executive Summary
	b.WriteString("<h2>Executive Summary</h2>\n")
	b.WriteString(fmt.Sprintf("<p>%s</p>\n", htmlEscape(output.ExecutiveReport.ArchitectureOverview)))

	if len(output.ExecutiveReport.MostCriticalAssumptions) > 0 {
		b.WriteString("<h3>Most Critical Assumptions</h3>\n<ul>\n")
		for _, a := range output.ExecutiveReport.MostCriticalAssumptions {
			b.WriteString(fmt.Sprintf("<li><span class=\"badge-%s\">%s</span> %s<br>\n", strings.ToLower(a.RiskLevel), a.RiskLevel, htmlEscape(a.Text)))
			b.WriteString(fmt.Sprintf("<small>Consequence: %s<br>Business Impact: %s</small></li>\n", htmlEscape(a.Consequence), htmlEscape(a.BusinessImpact)))
		}
		b.WriteString("</ul>\n")
	}

	if len(output.ExecutiveReport.HighImpactConsequences) > 0 {
		b.WriteString("<h3>High Impact Consequences</h3>\n<ul>\n")
		for _, c := range output.ExecutiveReport.HighImpactConsequences {
			b.WriteString(fmt.Sprintf("<li>%s</li>\n", htmlEscape(c)))
		}
		b.WriteString("</ul>\n")
	}

	if len(output.ExecutiveReport.TrustDependencies) > 0 {
		b.WriteString("<h3>Trust Dependencies</h3>\n<ul>\n")
		for _, d := range output.ExecutiveReport.TrustDependencies {
			b.WriteString(fmt.Sprintf("<li>%s</li>\n", htmlEscape(d)))
		}
		b.WriteString("</ul>\n")
	}

	if len(output.ExecutiveReport.SinglePointsOfFailure) > 0 {
		b.WriteString("<h3>Single Points of Failure</h3>\n<ul>\n")
		for _, s := range output.ExecutiveReport.SinglePointsOfFailure {
			b.WriteString(fmt.Sprintf("<li>%s</li>\n", htmlEscape(s)))
		}
		b.WriteString("</ul>\n")
	}

	if len(output.ExecutiveReport.ArchitecturalConcerns) > 0 {
		b.WriteString("<h3>Architectural Concerns</h3>\n<ul>\n")
		for _, c := range output.ExecutiveReport.ArchitecturalConcerns {
			b.WriteString(fmt.Sprintf("<li>%s</li>\n", htmlEscape(c)))
		}
		b.WriteString("</ul>\n")
	}

	if len(output.ExecutiveReport.RecommendedInvestments) > 0 {
		b.WriteString("<h3>Recommended Investments</h3>\n<ul>\n")
		for _, i := range output.ExecutiveReport.RecommendedInvestments {
			b.WriteString(fmt.Sprintf("<li>%s</li>\n", htmlEscape(i)))
		}
		b.WriteString("</ul>\n")
	}

	// Technical Summary
	b.WriteString("<h2>Technical Summary</h2>\n")
	b.WriteString(fmt.Sprintf("<p>%s</p>\n", htmlEscape(output.TechnicalSummary.ArchitectureSummary)))

	if len(output.TechnicalSummary.Recommendations) > 0 {
		b.WriteString("<h3>Prioritized Recommendations</h3>\n")
		b.WriteString("<table><tr><th>Priority</th><th>Recommendation</th><th>Assumption</th><th>Effort</th></tr>\n")
		for _, r := range output.TechnicalSummary.Recommendations {
			b.WriteString(fmt.Sprintf("<tr><td><span class=\"badge-%s\">%s</span></td><td>%s</td><td>%s</td><td>%s</td></tr>\n",
				strings.ToLower(r.Priority), r.Priority, htmlEscape(r.Recommendation), htmlEscape(r.AssumptionText), htmlEscape(r.ImplementationEffort)))
		}
		b.WriteString("</table>\n")
	}

	if len(output.TechnicalSummary.Dependencies) > 0 {
		b.WriteString("<h3>Dependencies</h3>\n<ul>\n")
		for _, d := range output.TechnicalSummary.Dependencies {
			b.WriteString(fmt.Sprintf("<li><strong>%s</strong> affects: %s</li>\n", htmlEscape(d.AssumptionText), htmlEscape(strings.Join(d.DependentComponents, ", "))))
		}
		b.WriteString("</ul>\n")
	}

	// Architect Narrative
	b.WriteString("<h2>Architect Narrative</h2>\n")
	b.WriteString(markdownToHTML(output.ArchitectNarrative))

	b.WriteString("</body>\n</html>")

	return b.String()
}

// htmlEscape escapes HTML special characters.
func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}

// markdownToHTML is a simple markdown-to-HTML converter for the narrative text.
func markdownToHTML(md string) string {
	lines := strings.Split(md, "\n")
	var result []string
	inList := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Headers
		if strings.HasPrefix(trimmed, "# ") {
			if inList {
				result = append(result, "</ul>")
				inList = false
			}
			text := strings.TrimPrefix(trimmed, "# ")
			result = append(result, fmt.Sprintf("<h2>%s</h2>", htmlEscape(text)))
			continue
		}
		if strings.HasPrefix(trimmed, "## ") {
			if inList {
				result = append(result, "</ul>")
				inList = false
			}
			text := strings.TrimPrefix(trimmed, "## ")
			result = append(result, fmt.Sprintf("<h3>%s</h3>", htmlEscape(text)))
			continue
		}
		if strings.HasPrefix(trimmed, "### ") {
			if inList {
				result = append(result, "</ul>")
				inList = false
			}
			text := strings.TrimPrefix(trimmed, "### ")
			result = append(result, fmt.Sprintf("<h4>%s</h4>", htmlEscape(text)))
			continue
		}

		// List items
		if strings.HasPrefix(trimmed, "- ") {
			if !inList {
				result = append(result, "<ul>")
				inList = true
			}
			text := strings.TrimPrefix(trimmed, "- ")
			result = append(result, fmt.Sprintf("<li>%s</li>", htmlEscape(text)))
			continue
		}

		// Bold
		if inList && !strings.HasPrefix(trimmed, "- ") && trimmed != "" {
			result = append(result, "</ul>")
			inList = false
		}

		// Horizontal rule
		if trimmed == "---" {
			result = append(result, "<hr>")
			continue
		}

		// Empty line
		if trimmed == "" {
			if inList {
				result = append(result, "</ul>")
				inList = false
			}
			result = append(result, "<p></p>")
			continue
		}

		// Regular paragraph
		if inList {
			result = append(result, "</ul>")
			inList = false
		}
		result = append(result, fmt.Sprintf("<p>%s</p>", htmlEscape(trimmed)))
	}

	if inList {
		result = append(result, "</ul>")
	}

	return strings.Join(result, "\n")
}
