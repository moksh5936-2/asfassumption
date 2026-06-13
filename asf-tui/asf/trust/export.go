package trust

import (
	"fmt"
	"strings"
)

// ExportMarkdown exports chain analysis as Markdown.
func ExportMarkdown(output *ChainOutput) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("# Trust Chain Analysis\n\n"))
	b.WriteString(fmt.Sprintf("**Generated:** %s\n\n", output.GeneratedAt))

	if output.Domain != "" {
		b.WriteString(fmt.Sprintf("**Domain:** %s\n\n", output.Domain))
	}

	// Dependency Graph Summary
	b.WriteString("## Dependency Graph\n\n")
	b.WriteString(output.DependencyGraph.Summary() + "\n\n")

	// Critical Assumptions
	if len(output.CriticalAssumptions) > 0 {
		b.WriteString("## Critical Assumptions\n\n")
		b.WriteString("| Rank | Assumption | Score | Risk | Centrality | Support | Failure Radius |\n")
		b.WriteString("|------|-----------|-------|------|------------|---------|----------------|\n")
		for i, ca := range output.CriticalAssumptions {
			if i >= 10 {
				break
			}
			b.WriteString(fmt.Sprintf("| %d | %s | %.2f | %s | %.2f | %d | %d |\n",
				i+1, ca.AssumptionText, ca.Score, ca.Risk, ca.Centrality, ca.SupportCount, ca.FailureRadius))
		}
		b.WriteString("\n")
	}

	// Trust Chains
	if len(output.TrustChains) > 0 {
		b.WriteString("## Trust Chains\n\n")
		for i, chain := range output.TrustChains {
			if i >= 10 {
				break
			}
			b.WriteString(fmt.Sprintf("### Chain %s (Length: %d, Confidence: %.2f, Risk: %s)\n\n",
				chain.ID, chain.Length, chain.Confidence, chain.Risk))
			for j, nodeID := range chain.Nodes {
				if node, ok := output.DependencyGraph.Nodes[nodeID]; ok {
					marker := "→"
					if j == 0 {
						marker = "●"
					}
					if j == len(chain.Nodes)-1 {
						marker = "⊘"
					}
					b.WriteString(fmt.Sprintf("%s %s [%s] %s\n", marker, node.ID, node.Risk, node.Text))
				}
			}
			b.WriteString("\n")
		}
	}

	// Failure Cascades
	if len(output.FailureCascades) > 0 {
		b.WriteString("## Failure Cascades\n\n")
		for i, cascade := range output.FailureCascades {
			if i >= 10 {
				break
			}
			b.WriteString(fmt.Sprintf("### %s (Severity: %s, Affected: %d, Max Depth: %d)\n\n",
				cascade.RootAssumptionText, cascade.Severity, cascade.TotalAffected, cascade.MaxDepth))
			for _, step := range cascade.Steps {
				b.WriteString(fmt.Sprintf("%d. **[%s]** %s — %s\n",
					step.Step, step.Severity, step.AssumptionText, step.Reason))
			}
			b.WriteString("\n")
		}
	}

	// Single Points of Trust Failure
	if len(output.SinglePointsOfTrust) > 0 {
		b.WriteString("## Single Points of Trust Failure\n\n")
		for _, spotf := range output.SinglePointsOfTrust {
			b.WriteString(fmt.Sprintf("### %s (Dependents: %d)\n\n", spotf.AssumptionText, spotf.DependentsCount))
			b.WriteString(fmt.Sprintf("**Dependent Assumptions:** %s\n\n", strings.Join(spotf.DependentNodes, ", ")))
			b.WriteString(fmt.Sprintf("**Recommendation:** %s\n\n", spotf.Recommendation))
		}
	}

	// Trust Collapse
	if len(output.TrustCollapseResults) > 0 {
		b.WriteString("## Trust Collapse Simulation\n\n")
		for i, collapse := range output.TrustCollapseResults {
			if i >= 10 {
				break
			}
			b.WriteString(fmt.Sprintf("### %s\n\n", collapse.FailedAssumptionText))
			b.WriteString(fmt.Sprintf("- **Assumptions Lost:** %d\n", len(collapse.AssumptionsLost)))
			b.WriteString(fmt.Sprintf("- **Risk Increase:** %s\n", collapse.RiskIncrease))
			b.WriteString(fmt.Sprintf("- **Affected Components:** %s\n\n", strings.Join(collapse.AffectedComponents, ", ")))
		}
	}

	return b.String()
}

// ExportHTML exports chain analysis as HTML.
func ExportHTML(output *ChainOutput) string {
	var b strings.Builder

	b.WriteString(`<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Trust Chain Analysis</title>
<style>
body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; line-height: 1.6; max-width: 900px; margin: 40px auto; padding: 20px; color: #333; }
h1 { border-bottom: 2px solid #2c3e50; padding-bottom: 10px; color: #2c3e50; }
h2 { color: #34495e; border-bottom: 1px solid #bdc3c7; padding-bottom: 5px; margin-top: 30px; }
h3 { color: #7f8c8d; }
table { border-collapse: collapse; width: 100%; margin: 20px 0; }
th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
th { background: #f8f9fa; font-weight: 600; }
tr:nth-child(even) { background: #f8f9fa; }
.badge-critical { background: #e74c3c; color: white; padding: 2px 8px; border-radius: 3px; font-size: 12px; font-weight: bold; }
.badge-high { background: #e67e22; color: white; padding: 2px 8px; border-radius: 3px; font-size: 12px; font-weight: bold; }
.badge-medium { background: #f39c12; color: white; padding: 2px 8px; border-radius: 3px; font-size: 12px; font-weight: bold; }
.chain { border-left: 3px solid #3498db; padding-left: 15px; margin: 20px 0; }
.cascade { border-left: 3px solid #e74c3c; padding-left: 15px; margin: 20px 0; }
.spotf { border-left: 3px solid #f39c12; padding-left: 15px; margin: 20px 0; }
</style>
</head>
<body>
`)

	b.WriteString(fmt.Sprintf("<h1>Trust Chain Analysis</h1>\n"))
	b.WriteString(fmt.Sprintf("<p><strong>Generated:</strong> %s</p>\n", output.GeneratedAt))

	if output.Domain != "" {
		b.WriteString(fmt.Sprintf("<p><strong>Domain:</strong> %s</p>\n", output.Domain))
	}

	// Critical Assumptions
	if len(output.CriticalAssumptions) > 0 {
		b.WriteString("<h2>Critical Assumptions</h2>\n")
		b.WriteString("<table><tr><th>Rank</th><th>Assumption</th><th>Score</th><th>Risk</th><th>Centrality</th><th>Support</th><th>Failure Radius</th></tr>\n")
		for i, ca := range output.CriticalAssumptions {
			if i >= 10 {
				break
			}
			b.WriteString(fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%.2f</td><td><span class=\"badge-%s\">%s</span></td><td>%.2f</td><td>%d</td><td>%d</td></tr>\n",
				i+1, htmlEscape(ca.AssumptionText), ca.Score, strings.ToLower(ca.Risk), ca.Risk, ca.Centrality, ca.SupportCount, ca.FailureRadius))
		}
		b.WriteString("</table>\n")
	}

	// Trust Chains
	if len(output.TrustChains) > 0 {
		b.WriteString("<h2>Trust Chains</h2>\n")
		for i, chain := range output.TrustChains {
			if i >= 10 {
				break
			}
			b.WriteString(fmt.Sprintf("<div class=\"chain\"><h3>Chain %s</h3><p>Length: %d | Confidence: %.2f | Risk: <span class=\"badge-%s\">%s</span></p>\n",
				chain.ID, chain.Length, chain.Confidence, strings.ToLower(chain.Risk), chain.Risk))
			for j, nodeID := range chain.Nodes {
				if node, ok := output.DependencyGraph.Nodes[nodeID]; ok {
					marker := "→"
					if j == 0 {
						marker = "●"
					}
					if j == len(chain.Nodes)-1 {
						marker = "⊘"
					}
					b.WriteString(fmt.Sprintf("<p>%s %s [%s] %s</p>\n", marker, htmlEscape(node.ID), htmlEscape(node.Risk), htmlEscape(node.Text)))
				}
			}
			b.WriteString("</div>\n")
		}
	}

	// Failure Cascades
	if len(output.FailureCascades) > 0 {
		b.WriteString("<h2>Failure Cascades</h2>\n")
		for i, cascade := range output.FailureCascades {
			if i >= 10 {
				break
			}
			b.WriteString(fmt.Sprintf("<div class=\"cascade\"><h3>%s</h3><p>Severity: <span class=\"badge-%s\">%s</span> | Affected: %d | Max Depth: %d</p>\n",
				htmlEscape(cascade.RootAssumptionText), strings.ToLower(cascade.Severity), cascade.Severity, cascade.TotalAffected, cascade.MaxDepth))
			b.WriteString("<ol>\n")
			for _, step := range cascade.Steps {
				b.WriteString(fmt.Sprintf("<li><strong>[%s]</strong> %s — %s</li>\n",
					htmlEscape(step.Severity), htmlEscape(step.AssumptionText), htmlEscape(step.Reason)))
			}
			b.WriteString("</ol></div>\n")
		}
	}

	// SPOTF
	if len(output.SinglePointsOfTrust) > 0 {
		b.WriteString("<h2>Single Points of Trust Failure</h2>\n")
		for _, spotf := range output.SinglePointsOfTrust {
			b.WriteString(fmt.Sprintf("<div class=\"spotf\"><h3>%s</h3><p>Dependents: %d</p>\n", htmlEscape(spotf.AssumptionText), spotf.DependentsCount))
			b.WriteString(fmt.Sprintf("<p><strong>Dependent Assumptions:</strong> %s</p>\n", htmlEscape(strings.Join(spotf.DependentNodes, ", "))))
			b.WriteString(fmt.Sprintf("<p><strong>Recommendation:</strong> %s</p></div>\n", htmlEscape(spotf.Recommendation)))
		}
	}

	b.WriteString("</body>\n</html>")

	return b.String()
}

func htmlEscape(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	return s
}
