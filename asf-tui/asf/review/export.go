package review

import (
	"encoding/json"
	"fmt"
	"strings"
)

func ExportMarkdown(output *ReviewOutput) string {
	var b strings.Builder

	b.WriteString("# Security Review Workbench Report\n\n")
	if output.Domain != "" {
		b.WriteString(fmt.Sprintf("**Domain:** %s  \n", output.Domain))
	}
	b.WriteString(fmt.Sprintf("**Generated:** %s  \n\n", output.GeneratedAt))

	if output.Queue == nil || len(output.Queue.Items) == 0 {
		b.WriteString("No review items generated.\n")
		return b.String()
	}

	b.WriteString("## Review Queue Summary\n\n")
	b.WriteString(fmt.Sprintf("| Metric | Value |\n|--------|-------|\n"))
	b.WriteString(fmt.Sprintf("| Total Items | %d |\n", output.Queue.TotalItems))
	b.WriteString(fmt.Sprintf("| Critical | %d |\n", output.Queue.CriticalCount))
	b.WriteString(fmt.Sprintf("| High | %d |\n", output.Queue.HighCount))
	b.WriteString(fmt.Sprintf("| Medium | %d |\n", output.Queue.MediumCount))
	b.WriteString(fmt.Sprintf("| Low | %d |\n", output.Queue.LowCount))
	b.WriteString("\n")

	if output.Matrix != nil {
		b.WriteString("## Priority Matrix\n\n")

		if len(output.Matrix.HighValueLowEffort) > 0 {
			b.WriteString("### High Value / Low Effort (Do First)\n\n")
			for _, it := range output.Matrix.HighValueLowEffort[:min(5, len(output.Matrix.HighValueLowEffort))] {
				b.WriteString(fmt.Sprintf("- **%s** [Score: %.0f] — %s\n", it.AssumptionText, it.PriorityScore, it.WhyReview))
			}
			b.WriteString("\n")
		}
		if len(output.Matrix.HighValueHighEffort) > 0 {
			b.WriteString("### High Value / High Effort (Plan)\n\n")
			for _, it := range output.Matrix.HighValueHighEffort[:min(5, len(output.Matrix.HighValueHighEffort))] {
				b.WriteString(fmt.Sprintf("- **%s** [Score: %.0f] — %s\n", it.AssumptionText, it.PriorityScore, it.EstimatedTime))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("## Prioritized Review Queue\n\n")
	for i, it := range output.Queue.Items {
		if i >= 30 {
			b.WriteString(fmt.Sprintf("... and %d more items\n", len(output.Queue.Items)-30))
			break
		}
		b.WriteString(fmt.Sprintf("### #%d. %s\n\n", it.Rank, it.AssumptionText))
		b.WriteString(fmt.Sprintf("- **Risk:** %s\n", it.Risk))
		b.WriteString(fmt.Sprintf("- **Priority Score:** %.0f/100\n", it.PriorityScore))
		b.WriteString(fmt.Sprintf("- **Category:** %s\n", it.Category))
		b.WriteString(fmt.Sprintf("- **Value:** %s | **Effort:** %s\n", string(it.ReviewValue), string(it.ReviewEffort)))
		b.WriteString(fmt.Sprintf("- **Quadrant:** %s\n", string(it.Quadrant)))
		b.WriteString(fmt.Sprintf("- **Estimated Time:** %s\n", it.EstimatedTime))
		b.WriteString("\n")

		if it.WhyReview != "" {
			b.WriteString(fmt.Sprintf("**Why Review:** %s  \n\n", it.WhyReview))
		}
		if it.WhatToReview != "" {
			b.WriteString(fmt.Sprintf("**What to Review:** %s  \n\n", it.WhatToReview))
		}
		if it.ExpectedEvidence != "" {
			b.WriteString(fmt.Sprintf("**Expected Evidence:** %s  \n\n", it.ExpectedEvidence))
		}
		if it.ExpectedOutcome != "" {
			b.WriteString(fmt.Sprintf("**Expected Outcome:** %s  \n\n", it.ExpectedOutcome))
		}
		if it.ExpectedRiskReduction != "" {
			b.WriteString(fmt.Sprintf("**Expected Risk Reduction:** %s  \n\n", it.ExpectedRiskReduction))
		}

		b.WriteString("---\n\n")
	}

	if len(output.Campaigns) > 0 {
		b.WriteString("## Review Campaigns\n\n")
		for _, c := range output.Campaigns {
			if len(c.Items) == 0 {
				continue
			}
			b.WriteString(fmt.Sprintf("### %s (%s)\n\n", c.Name, c.Duration))
			b.WriteString(fmt.Sprintf("**Items:** %d | **Total Effort:** %s\n\n", c.TotalItems, c.TotalEffort))
			for _, it := range c.Items {
				b.WriteString(fmt.Sprintf("- [%.0f] %s (%s)\n", it.PriorityScore, it.AssumptionText, it.EstimatedTime))
			}
			b.WriteString("\n")
		}
	}

	if output.CISODashboard != nil {
		b.WriteString("## CISO Security Review Dashboard\n\n")
		b.WriteString(fmt.Sprintf("**Total Assumptions:** %d | **Critical:** %d | **High:** %d\n\n",
			output.CISODashboard.TotalAssumptions, output.CISODashboard.CriticalAssumptions, output.CISODashboard.HighAssumptions))

		if len(output.CISODashboard.GreatestRiskReduction) > 0 {
			b.WriteString("### Greatest Risk Reduction Opportunities\n\n")
			for _, it := range output.CISODashboard.GreatestRiskReduction[:min(5, len(output.CISODashboard.GreatestRiskReduction))] {
				b.WriteString(fmt.Sprintf("- **%s** [Score: %.0f] — %s\n", it.AssumptionText, it.PriorityScore, it.ExpectedRiskReduction))
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}

func ExportHTML(output *ReviewOutput) string {
	var b strings.Builder

	b.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	b.WriteString("<meta charset=\"utf-8\">\n")
	b.WriteString("<title>Security Review Workbench Report</title>\n")
	b.WriteString("<style>\n")
	b.WriteString("body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; max-width: 960px; margin: 0 auto; padding: 20px; color: #333; }\n")
	b.WriteString("h1 { color: #111; border-bottom: 2px solid #2980b9; padding-bottom: 10px; }\n")
	b.WriteString("h2 { color: #222; margin-top: 30px; }\n")
	b.WriteString("h3 { color: #333; margin-top: 20px; }\n")
	b.WriteString("table { border-collapse: collapse; width: 100%; margin: 15px 0; }\n")
	b.WriteString("th, td { border: 1px solid #ddd; padding: 8px 12px; text-align: left; }\n")
	b.WriteString("th { background: #f5f5f5; }\n")
	b.WriteString(".item { border: 1px solid #ddd; border-radius: 6px; padding: 15px; margin: 10px 0; background: #fafafa; }\n")
	b.WriteString(".score { font-size: 24px; font-weight: bold; color: #2980b9; }\n")
	b.WriteString(".critical { color: #e74c3c; }\n")
	b.WriteString(".high { color: #e67e22; }\n")
	b.WriteString(".medium { color: #f39c12; }\n")
	b.WriteString(".low { color: #27ae60; }\n")
	b.WriteString(".quadrant { display: inline-block; padding: 3px 8px; border-radius: 4px; font-size: 12px; }\n")
	b.WriteString(".q-hvle { background: #27ae60; color: white; }\n")
	b.WriteString(".q-hvhe { background: #e67e22; color: white; }\n")
	b.WriteString(".q-lvle { background: #95a5a6; color: white; }\n")
	b.WriteString(".q-lvhe { background: #e74c3c; color: white; }\n")
	b.WriteString("</style>\n</head>\n<body>\n")

	b.WriteString("<h1>Security Review Workbench Report</h1>\n")
	if output.Domain != "" {
		b.WriteString(fmt.Sprintf("<p><strong>Domain:</strong> %s</p>\n", output.Domain))
	}
	b.WriteString(fmt.Sprintf("<p><em>Generated: %s</em></p>\n", output.GeneratedAt))

	if output.Queue != nil {
		b.WriteString("<h2>Review Queue Summary</h2>\n")
		b.WriteString("<table>\n")
		b.WriteString("<tr><th>Metric</th><th>Value</th></tr>\n")
		b.WriteString(fmt.Sprintf("<tr><td>Total Items</td><td>%d</td></tr>\n", output.Queue.TotalItems))
		b.WriteString(fmt.Sprintf("<tr><td>Critical</td><td>%d</td></tr>\n", output.Queue.CriticalCount))
		b.WriteString(fmt.Sprintf("<tr><td>High</td><td>%d</td></tr>\n", output.Queue.HighCount))
		b.WriteString(fmt.Sprintf("<tr><td>Medium</td><td>%d</td></tr>\n", output.Queue.MediumCount))
		b.WriteString(fmt.Sprintf("<tr><td>Low</td><td>%d</td></tr>\n", output.Queue.LowCount))
		b.WriteString("</table>\n")

		b.WriteString("<h2>Prioritized Review Queue</h2>\n")
		for i, it := range output.Queue.Items {
			if i >= 20 {
				b.WriteString(fmt.Sprintf("<p><em>... and %d more items</em></p>\n", len(output.Queue.Items)-20))
				break
			}
			quadClass := "q-lvle"
			switch it.Quadrant {
			case QuadHighValueLowEffort:
				quadClass = "q-hvle"
			case QuadHighValueHighEffort:
				quadClass = "q-hvhe"
			case QuadLowValueHighEffort:
				quadClass = "q-lvhe"
			}

			b.WriteString(fmt.Sprintf("<div class=\"item\">\n"))
			b.WriteString(fmt.Sprintf("<h3>#%d. %s</h3>\n", it.Rank, it.AssumptionText))
			b.WriteString(fmt.Sprintf("<p><span class=\"score\">%.0f</span> | Risk: <span class=\"%s\">%s</span> | Category: %s</p>\n",
				it.PriorityScore, strings.ToLower(it.Risk), it.Risk, it.Category))
			b.WriteString(fmt.Sprintf("<p>Value: <strong>%s</strong> | Effort: <strong>%s</strong> | <span class=\"quadrant %s\">%s</span></p>\n",
				string(it.ReviewValue), string(it.ReviewEffort), quadClass, string(it.Quadrant)))
			b.WriteString(fmt.Sprintf("<p><strong>Estimated Time:</strong> %s</p>\n", it.EstimatedTime))

			if it.WhyReview != "" {
				b.WriteString(fmt.Sprintf("<p><strong>Why Review:</strong> %s</p>\n", it.WhyReview))
			}
			if it.WhatToReview != "" {
				b.WriteString(fmt.Sprintf("<p><strong>What to Review:</strong> %s</p>\n", it.WhatToReview))
			}
			if it.ExpectedRiskReduction != "" {
				b.WriteString(fmt.Sprintf("<p><strong>Expected Risk Reduction:</strong> %s</p>\n", it.ExpectedRiskReduction))
			}
			b.WriteString("</div>\n")
		}
	}

	b.WriteString("</body>\n</html>")
	return b.String()
}

func ExportJSON(output *ReviewOutput) string {
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error":"%s"}`, err.Error())
	}
	return string(data)
}
