package confidencex

import (
	"fmt"
	"strings"
)

func ExportMarkdown(output *ConfidenceOutput) string {
	var b strings.Builder
	b.WriteString("# Confidence & Explainability Report\n\n")
	b.WriteString(fmt.Sprintf("**Generated:** %s\n", output.GeneratedAt))
	b.WriteString(fmt.Sprintf("**Domain:** %s\n\n", output.Domain))

	if len(output.Breakdowns) == 0 {
		b.WriteString("_No confidence data available._\n")
		return b.String()
	}

	b.WriteString(fmt.Sprintf("**Total Assumptions:** %d\n\n", len(output.Breakdowns)))

	for _, bd := range output.Breakdowns {
		b.WriteString(fmt.Sprintf("## %s\n", bd.AssumptionText))
		b.WriteString(fmt.Sprintf("**Assumption ID:** %s  \n", bd.AssumptionID))
		b.WriteString(fmt.Sprintf("**Confidence:** %.1f%%  \n", bd.FinalConfidence))
		b.WriteString(fmt.Sprintf("**Stability:** %s  \n", bd.StabilityClass))
		b.WriteString(fmt.Sprintf("**Stability Reason:** %s  \n\n", bd.StabilityReason))

		b.WriteString("### Why ASF Believes This\n")
		b.WriteString(bd.WhyExists + "\n\n")

		b.WriteString("### Why ASF Is Uncertain\n")
		b.WriteString(bd.WhyUncertain + "\n\n")

		b.WriteString("### What Would Increase Confidence\n")
		b.WriteString(bd.WhatIncreasesConfidence + "\n\n")

		b.WriteString("### What Would Decrease Confidence\n")
		b.WriteString(bd.WhatDecreasesConfidence + "\n\n")

		if len(bd.PositiveFactors) > 0 {
			b.WriteString("### Positive Factors\n")
			b.WriteString("| Factor | Impact | Description |\n")
			b.WriteString("|--------|--------|-------------|\n")
			for _, f := range bd.PositiveFactors {
				b.WriteString(fmt.Sprintf("| %s | +%.1f | %s |\n", f.Name, f.Impact, f.Description))
			}
			b.WriteString("\n")
		}

		if len(bd.NegativeFactors) > 0 {
			b.WriteString("### Negative Factors\n")
			b.WriteString("| Factor | Impact | Description |\n")
			b.WriteString("|--------|--------|-------------|\n")
			for _, f := range bd.NegativeFactors {
				b.WriteString(fmt.Sprintf("| %s | %.1f | %s |\n", f.Name, f.Impact, f.Description))
			}
			b.WriteString("\n")
		}

		if len(bd.SupportingFacts) > 0 {
			b.WriteString("### Supporting Facts\n")
			b.WriteString("| Fact ID | Fact | Contribution |\n")
			b.WriteString("|---------|------|-------------|\n")
			for _, fc := range bd.SupportingFacts {
				mark := "+"
				if !fc.IsPositive {
					mark = ""
				}
				b.WriteString(fmt.Sprintf("| %s | %s | %s%.1f%% |\n", fc.FactID, fc.FactText, mark, fc.Contribution))
			}
			b.WriteString("\n")
		}

		if len(bd.EvidenceContributions) > 0 {
			b.WriteString("### Evidence Contributions\n")
			b.WriteString("| Evidence | Present | Impact | Label |\n")
			b.WriteString("|----------|---------|--------|-------|\n")
			for _, ec := range bd.EvidenceContributions {
				present := "No"
				if ec.Present {
					present = "Yes"
				}
				b.WriteString(fmt.Sprintf("| %s | %s | %.1f | %s |\n", ec.EvidenceID, present, ec.Impact, ec.Label))
			}
			b.WriteString("\n")
		}

		if bd.DomainContribution != nil {
			b.WriteString("### Domain Contribution\n")
			b.WriteString(fmt.Sprintf("**Domain:** %s  \n", bd.DomainContribution.Domain))
			b.WriteString(fmt.Sprintf("**Influence:** +%.1f%%  \n", bd.DomainContribution.Influence))
			b.WriteString(fmt.Sprintf("**Strength:** %s  \n", bd.DomainContribution.Strength))
			b.WriteString(fmt.Sprintf("**Reason:** %s  \n\n", bd.DomainContribution.Reason))
		}

		if bd.TrustContribution != nil {
			b.WriteString("### Trust Contribution\n")
			b.WriteString(fmt.Sprintf("**Has Trust Chain:** %v  \n", bd.TrustContribution.HasTrustChain))
			b.WriteString(fmt.Sprintf("**Chain Influence:** %.1f  \n", bd.TrustContribution.ChainInfluence))
			b.WriteString(fmt.Sprintf("**Dependency Centrality:** %.2f  \n", bd.TrustContribution.DependencyCentrality))
			b.WriteString(fmt.Sprintf("**Failure Radius Influence:** %.1f  \n\n", bd.TrustContribution.FailureRadiusInfluence))
		}
	}

	if output.CISOTrustView != nil {
		b.WriteString("## CISO Trust View\n\n")

		if len(output.CISOTrustView.MostTrustedFindings) > 0 {
			b.WriteString("### Most Trusted Findings\n")
			for _, f := range output.CISOTrustView.MostTrustedFindings {
				b.WriteString(fmt.Sprintf("- [%.1f%%] %s\n", f.FinalConfidence, f.AssumptionText))
			}
			b.WriteString("\n")
		}
		if len(output.CISOTrustView.LeastTrustedFindings) > 0 {
			b.WriteString("### Least Trusted Findings\n")
			for _, f := range output.CISOTrustView.LeastTrustedFindings {
				b.WriteString(fmt.Sprintf("- [%.1f%%] %s\n", f.FinalConfidence, f.AssumptionText))
			}
			b.WriteString("\n")
		}
		if len(output.CISOTrustView.MostCriticalLowConfidence) > 0 {
			b.WriteString("### Most Critical Low-Confidence Findings\n")
			for _, f := range output.CISOTrustView.MostCriticalLowConfidence {
				b.WriteString(fmt.Sprintf("- [%.1f%%] %s\n", f.FinalConfidence, f.AssumptionText))
			}
			b.WriteString("\n")
		}
		if len(output.CISOTrustView.HighestRiskUnknowns) > 0 {
			b.WriteString("### Highest-Risk Unknowns\n")
			for _, f := range output.CISOTrustView.HighestRiskUnknowns {
				b.WriteString(fmt.Sprintf("- [%.1f%%] %s\n", f.FinalConfidence, f.AssumptionText))
			}
			b.WriteString("\n")
		}
	}

	if output.ArchitectReviewView != nil {
		b.WriteString("## Architect Review View\n\n")

		if len(output.ArchitectReviewView.RequiringValidation) > 0 {
			b.WriteString("### Assumptions Requiring Validation\n")
			for _, f := range output.ArchitectReviewView.RequiringValidation {
				b.WriteString(fmt.Sprintf("- [%.1f%%] %s\n", f.FinalConfidence, f.AssumptionText))
			}
			b.WriteString("\n")
		}
		if len(output.ArchitectReviewView.WeakSupport) > 0 {
			b.WriteString("### Assumptions With Weak Support\n")
			for _, f := range output.ArchitectReviewView.WeakSupport {
				b.WriteString(fmt.Sprintf("- [%.1f%%] %s\n", f.FinalConfidence, f.AssumptionText))
			}
			b.WriteString("\n")
		}
		if len(output.ArchitectReviewView.StrongSupport) > 0 {
			b.WriteString("### Assumptions With Strong Support\n")
			for _, f := range output.ArchitectReviewView.StrongSupport {
				b.WriteString(fmt.Sprintf("- [%.1f%%] %s\n", f.FinalConfidence, f.AssumptionText))
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}

func ExportHTML(output *ConfidenceOutput) string {
	var b strings.Builder
	b.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	b.WriteString("<title>Confidence & Explainability Report</title>\n")
	b.WriteString("<style>\n")
	b.WriteString("body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; max-width: 960px; margin: 40px auto; padding: 0 20px; color: #1a1a2e; line-height: 1.6; }\n")
	b.WriteString("h1 { color: #1a1a2e; border-bottom: 2px solid #e94560; padding-bottom: 10px; }\n")
	b.WriteString("h2 { color: #16213e; margin-top: 30px; border-bottom: 1px solid #ccc; padding-bottom: 5px; }\n")
	b.WriteString("h3 { color: #0f3460; margin-top: 20px; }\n")
	b.WriteString("h4 { color: #e94560; margin: 10px 0; }\n")
	b.WriteString("table { border-collapse: collapse; width: 100%; margin: 10px 0; }\n")
	b.WriteString("th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }\n")
	b.WriteString("th { background-color: #16213e; color: white; }\n")
	b.WriteString("tr:nth-child(even) { background-color: #f8f8f8; }\n")
	b.WriteString(".confidence-high { color: #27ae60; font-weight: bold; }\n")
	b.WriteString(".confidence-mid { color: #f39c12; font-weight: bold; }\n")
	b.WriteString(".confidence-low { color: #e74c3c; font-weight: bold; }\n")
	b.WriteString(".stability { display: inline-block; padding: 2px 8px; border-radius: 3px; font-size: 0.9em; }\n")
	b.WriteString(".stability-very-stable { background: #27ae60; color: white; }\n")
	b.WriteString(".stability-stable { background: #2ecc71; color: white; }\n")
	b.WriteString(".stability-moderate { background: #f39c12; color: white; }\n")
	b.WriteString(".stability-weak { background: #e67e22; color: white; }\n")
	b.WriteString(".stability-highly-speculative { background: #e74c3c; color: white; }\n")
	b.WriteString(".positive { color: #27ae60; }\n")
	b.WriteString(".negative { color: #e74c3c; }\n")
	b.WriteString(".section { background: #f5f6fa; padding: 15px; border-radius: 5px; margin: 15px 0; }\n")
	b.WriteString("</style>\n</head>\n<body>\n")

	b.WriteString("<h1>Confidence & Explainability Report</h1>\n")
	b.WriteString(fmt.Sprintf("<p><strong>Generated:</strong> %s</p>\n", output.GeneratedAt))
	b.WriteString(fmt.Sprintf("<p><strong>Domain:</strong> %s</p>\n", output.Domain))

	if len(output.Breakdowns) == 0 {
		b.WriteString("<p><em>No confidence data available.</em></p>\n")
		b.WriteString("</body>\n</html>\n")
		return b.String()
	}

	b.WriteString(fmt.Sprintf("<p><strong>Total Assumptions:</strong> %d</p>\n", len(output.Breakdowns)))

	for _, bd := range output.Breakdowns {
		confClass := "confidence-mid"
		if bd.FinalConfidence >= 70 {
			confClass = "confidence-high"
		} else if bd.FinalConfidence < 40 {
			confClass = "confidence-low"
		}

		stabClass := "stability-moderate"
		switch bd.StabilityClass {
		case StabilityVeryStable:
			stabClass = "stability-very-stable"
		case StabilityStable:
			stabClass = "stability-stable"
		case StabilityWeak:
			stabClass = "stability-weak"
		case StabilityHighlySpeculative:
			stabClass = "stability-highly-speculative"
		}

		b.WriteString(fmt.Sprintf("<h2>%s</h2>\n", bd.AssumptionText))
		b.WriteString(fmt.Sprintf("<p><strong>Assumption ID:</strong> %s<br>\n", bd.AssumptionID))
		b.WriteString(fmt.Sprintf("<strong>Confidence:</strong> <span class=\"%s\">%.1f%%</span><br>\n", confClass, bd.FinalConfidence))
		b.WriteString(fmt.Sprintf("<strong>Stability:</strong> <span class=\"stability %s\">%s</span><br>\n", stabClass, bd.StabilityClass))
		b.WriteString(fmt.Sprintf("<strong>Stability Reason:</strong> %s</p>\n", bd.StabilityReason))

		b.WriteString("<div class=\"section\">\n")
		b.WriteString("<h3>Why ASF Believes This</h3>\n")
		b.WriteString(fmt.Sprintf("<p>%s</p>\n", bd.WhyExists))
		b.WriteString("</div>\n")

		b.WriteString("<div class=\"section\">\n")
		b.WriteString("<h3>Why ASF Is Uncertain</h3>\n")
		b.WriteString(fmt.Sprintf("<p>%s</p>\n", bd.WhyUncertain))
		b.WriteString("</div>\n")

		b.WriteString("<div class=\"section\">\n")
		b.WriteString("<h3>What Would Increase Confidence</h3>\n")
		b.WriteString(fmt.Sprintf("<p>%s</p>\n", bd.WhatIncreasesConfidence))
		b.WriteString("</div>\n")

		b.WriteString("<div class=\"section\">\n")
		b.WriteString("<h3>What Would Decrease Confidence</h3>\n")
		b.WriteString(fmt.Sprintf("<p>%s</p>\n", bd.WhatDecreasesConfidence))
		b.WriteString("</div>\n")

		if len(bd.PositiveFactors) > 0 {
			b.WriteString("<h3>Positive Factors</h3>\n")
			b.WriteString("<table><tr><th>Factor</th><th>Impact</th><th>Description</th></tr>\n")
			for _, f := range bd.PositiveFactors {
				b.WriteString(fmt.Sprintf("<tr><td>%s</td><td class=\"positive\">+%.1f</td><td>%s</td></tr>\n", f.Name, f.Impact, f.Description))
			}
			b.WriteString("</table>\n")
		}

		if len(bd.NegativeFactors) > 0 {
			b.WriteString("<h3>Negative Factors</h3>\n")
			b.WriteString("<table><tr><th>Factor</th><th>Impact</th><th>Description</th></tr>\n")
			for _, f := range bd.NegativeFactors {
				b.WriteString(fmt.Sprintf("<tr><td>%s</td><td class=\"negative\">%.1f</td><td>%s</td></tr>\n", f.Name, f.Impact, f.Description))
			}
			b.WriteString("</table>\n")
		}

		if len(bd.SupportingFacts) > 0 {
			b.WriteString("<h3>Supporting Facts</h3>\n")
			b.WriteString("<table><tr><th>Fact ID</th><th>Fact</th><th>Contribution</th></tr>\n")
			for _, fc := range bd.SupportingFacts {
				sign := "+"
				cls := "positive"
				if !fc.IsPositive {
					sign = ""
					cls = "negative"
				}
				b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td class=\"%s\">%s%.1f%%</td></tr>\n", fc.FactID, fc.FactText, cls, sign, fc.Contribution))
			}
			b.WriteString("</table>\n")
		}

		if len(bd.EvidenceContributions) > 0 {
			b.WriteString("<h3>Evidence Contributions</h3>\n")
			b.WriteString("<table><tr><th>Evidence</th><th>Present</th><th>Impact</th><th>Label</th></tr>\n")
			for _, ec := range bd.EvidenceContributions {
				present := "No"
				cls := "negative"
				if ec.Present {
					present = "Yes"
					cls = "positive"
				}
				b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td class=\"%s\">%.1f</td><td>%s</td></tr>\n", ec.EvidenceID, present, cls, ec.Impact, ec.Label))
			}
			b.WriteString("</table>\n")
		}

		if bd.DomainContribution != nil {
			b.WriteString("<h3>Domain Contribution</h3>\n")
			b.WriteString("<table><tr><th>Domain</th><th>Influence</th><th>Strength</th><th>Reason</th></tr>\n")
			b.WriteString(fmt.Sprintf("<tr><td>%s</td><td class=\"positive\">+%.1f%%</td><td>%s</td><td>%s</td></tr>\n",
				bd.DomainContribution.Domain, bd.DomainContribution.Influence, bd.DomainContribution.Strength, bd.DomainContribution.Reason))
			b.WriteString("</table>\n")
		}

		if bd.TrustContribution != nil {
			b.WriteString("<h3>Trust Contribution</h3>\n")
			b.WriteString("<table><tr><th>Property</th><th>Value</th></tr>\n")
			b.WriteString(fmt.Sprintf("<tr><td>Has Trust Chain</td><td>%v</td></tr>\n", bd.TrustContribution.HasTrustChain))
			b.WriteString(fmt.Sprintf("<tr><td>Chain Influence</td><td>%.1f</td></tr>\n", bd.TrustContribution.ChainInfluence))
			b.WriteString(fmt.Sprintf("<tr><td>Dependency Centrality</td><td>%.2f</td></tr>\n", bd.TrustContribution.DependencyCentrality))
			b.WriteString(fmt.Sprintf("<tr><td>Failure Radius Influence</td><td>%.1f</td></tr>\n", bd.TrustContribution.FailureRadiusInfluence))
			b.WriteString("</table>\n")
		}
	}

	if output.CISOTrustView != nil {
		b.WriteString("<h2>CISO Trust View</h2>\n")

		if len(output.CISOTrustView.MostTrustedFindings) > 0 {
			b.WriteString("<h3>Most Trusted Findings</h3>\n<ul>\n")
			for _, f := range output.CISOTrustView.MostTrustedFindings {
				b.WriteString(fmt.Sprintf("<li><span class=\"confidence-high\">[%.1f%%]</span> %s</li>\n", f.FinalConfidence, f.AssumptionText))
			}
			b.WriteString("</ul>\n")
		}
		if len(output.CISOTrustView.LeastTrustedFindings) > 0 {
			b.WriteString("<h3>Least Trusted Findings</h3>\n<ul>\n")
			for _, f := range output.CISOTrustView.LeastTrustedFindings {
				b.WriteString(fmt.Sprintf("<li><span class=\"confidence-low\">[%.1f%%]</span> %s</li>\n", f.FinalConfidence, f.AssumptionText))
			}
			b.WriteString("</ul>\n")
		}
		if len(output.CISOTrustView.MostCriticalLowConfidence) > 0 {
			b.WriteString("<h3>Most Critical Low-Confidence Findings</h3>\n<ul>\n")
			for _, f := range output.CISOTrustView.MostCriticalLowConfidence {
				b.WriteString(fmt.Sprintf("<li><span class=\"confidence-low\">[%.1f%%]</span> %s</li>\n", f.FinalConfidence, f.AssumptionText))
			}
			b.WriteString("</ul>\n")
		}
		if len(output.CISOTrustView.HighestRiskUnknowns) > 0 {
			b.WriteString("<h3>Highest-Risk Unknowns</h3>\n<ul>\n")
			for _, f := range output.CISOTrustView.HighestRiskUnknowns {
				b.WriteString(fmt.Sprintf("<li><span class=\"confidence-low\">[%.1f%%]</span> %s</li>\n", f.FinalConfidence, f.AssumptionText))
			}
			b.WriteString("</ul>\n")
		}
	}

	if output.ArchitectReviewView != nil {
		b.WriteString("<h2>Architect Review View</h2>\n")

		if len(output.ArchitectReviewView.RequiringValidation) > 0 {
			b.WriteString("<h3>Assumptions Requiring Validation</h3>\n<ul>\n")
			for _, f := range output.ArchitectReviewView.RequiringValidation {
				b.WriteString(fmt.Sprintf("<li><span class=\"confidence-low\">[%.1f%%]</span> %s</li>\n", f.FinalConfidence, f.AssumptionText))
			}
			b.WriteString("</ul>\n")
		}
		if len(output.ArchitectReviewView.WeakSupport) > 0 {
			b.WriteString("<h3>Assumptions With Weak Support</h3>\n<ul>\n")
			for _, f := range output.ArchitectReviewView.WeakSupport {
				b.WriteString(fmt.Sprintf("<li><span class=\"confidence-mid\">[%.1f%%]</span> %s</li>\n", f.FinalConfidence, f.AssumptionText))
			}
			b.WriteString("</ul>\n")
		}
		if len(output.ArchitectReviewView.StrongSupport) > 0 {
			b.WriteString("<h3>Assumptions With Strong Support</h3>\n<ul>\n")
			for _, f := range output.ArchitectReviewView.StrongSupport {
				b.WriteString(fmt.Sprintf("<li><span class=\"confidence-high\">[%.1f%%]</span> %s</li>\n", f.FinalConfidence, f.AssumptionText))
			}
			b.WriteString("</ul>\n")
		}
	}

	b.WriteString("</body>\n</html>\n")
	return b.String()
}
