package verify

import (
	"encoding/json"
	"fmt"
	"strings"
)

func ExportMarkdown(output *VerificationOutput) string {
	var b strings.Builder

	b.WriteString("# Verification Intelligence Report\n\n")
	if output.Domain != "" {
		b.WriteString(fmt.Sprintf("**Domain:** %s  \n", output.Domain))
	}
	b.WriteString(fmt.Sprintf("**Generated:** %s  \n\n", output.GeneratedAt))

	if output.Assessment == nil || len(output.Assessment.Plans) == 0 {
		b.WriteString("No verification plans generated.\n")
		return b.String()
	}

	b.WriteString("## Executive Summary\n\n")
	b.WriteString(fmt.Sprintf("| Metric | Value |\n|--------|-------|\n"))
	b.WriteString(fmt.Sprintf("| Total Assumptions | %d |\n", output.Assessment.TotalAssumptions))
	b.WriteString(fmt.Sprintf("| Verified | %d |\n", output.Assessment.VerifiedCount))
	b.WriteString(fmt.Sprintf("| Partially Verified | %d |\n", output.Assessment.PartialCount))
	b.WriteString(fmt.Sprintf("| Unverified | %d |\n", output.Assessment.UnverifiedCount))
	b.WriteString(fmt.Sprintf("| No Evidence | %d |\n", output.Assessment.NoEvidenceCount))
	b.WriteString(fmt.Sprintf("| Overall Confidence | %.1f%% |\n", output.Assessment.OverallConfidence))
	b.WriteString("\n")

	if output.CISOView != nil {
		b.WriteString("## CISO Review\n\n")
		b.WriteString(fmt.Sprintf("**Verification Backlog:** %d items  \n", len(output.CISOView.VerificationBacklog)))
		b.WriteString(fmt.Sprintf("**Evidence Gaps:** %d  \n", len(output.CISOView.EvidenceGaps)))
		b.WriteString(fmt.Sprintf("**Priority Distribution:** Critical=%d, High=%d, Medium=%d, Low=%d  \n\n",
			output.CISOView.CriticalCount, output.CISOView.HighCount,
			output.CISOView.MediumCount, output.CISOView.LowCount))

		if len(output.CISOView.HighestRiskUnverified) > 0 {
			b.WriteString("### Highest Risk Unverified\n\n")
			for _, p := range output.CISOView.HighestRiskUnverified {
				b.WriteString(fmt.Sprintf("- **%s** [%s] — %s\n", p.AssumptionText, string(p.Priority), p.Risk))
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("## Verification Plans\n\n")
	for i, p := range output.Assessment.Plans {
		b.WriteString(fmt.Sprintf("### %d. %s\n", i+1, p.AssumptionText))
		b.WriteString(fmt.Sprintf("- **Category:** %s\n", string(p.Category)))
		b.WriteString(fmt.Sprintf("- **Risk:** %s\n", p.Risk))
		b.WriteString(fmt.Sprintf("- **Confidence:** %.0f/100\n", p.Confidence))
		b.WriteString(fmt.Sprintf("- **Priority:** %s\n", string(p.Priority)))
		b.WriteString(fmt.Sprintf("- **Status:** %s\n", string(p.Status)))
		b.WriteString(fmt.Sprintf("- **Effort:** %s\n", string(p.Effort)))
		b.WriteString(fmt.Sprintf("- **Expected Time:** %s\n", p.ExpectedTime))
		b.WriteString("\n")

		if p.WhyVerify != "" {
			b.WriteString(fmt.Sprintf("**Why Verify:** %s  \n\n", p.WhyVerify))
		}

		if len(p.EvidenceRequired) > 0 {
			b.WriteString("**Required Evidence:**\n\n")
			for _, ev := range p.EvidenceRequired {
				opt := ""
				if ev.Optional {
					opt = " (optional)"
				}
				b.WriteString(fmt.Sprintf("- %s — %s%s\n", ev.Name, ev.Description, opt))
			}
			b.WriteString("\n")
		}

		if len(p.Actions) > 0 {
			b.WriteString("**Verification Actions:**\n\n")
			for _, a := range p.Actions {
				b.WriteString(fmt.Sprintf("%d. %s — %s (%s)\n", a.Step, a.Action, a.Description, a.Stakeholder))
			}
			b.WriteString("\n")
		}

		if p.WhatToReview != "" {
			b.WriteString(fmt.Sprintf("**What to Review:** %s  \n\n", p.WhatToReview))
		}
		if p.WhatEvidenceToCollect != "" {
			b.WriteString(fmt.Sprintf("**What Evidence to Collect:** %s  \n\n", p.WhatEvidenceToCollect))
		}
		if p.HowToValidate != "" {
			b.WriteString(fmt.Sprintf("**How to Validate:** %s  \n\n", p.HowToValidate))
		}

		if len(p.Stakeholders) > 0 {
			b.WriteString(fmt.Sprintf("**Stakeholders:** %s  \n\n", strings.Join(p.Stakeholders, ", ")))
		}
	}

	if len(output.Roadmaps) > 0 {
		b.WriteString("## Verification Roadmaps\n\n")
		for i, r := range output.Roadmaps {
			if i >= 20 {
				b.WriteString(fmt.Sprintf("... and %d more roadmaps\n", len(output.Roadmaps)-20))
				break
			}
			b.WriteString(fmt.Sprintf("### Roadmap: %s\n", r.AssumptionText))
			b.WriteString(fmt.Sprintf("- **Priority:** %s\n", string(r.Priority)))
			b.WriteString(fmt.Sprintf("- **Effort:** %s\n", string(r.Effort)))
			b.WriteString(fmt.Sprintf("- **Stakeholders:** %s\n\n", strings.Join(r.Stakeholders, ", ")))
			for _, s := range r.Steps {
				b.WriteString(fmt.Sprintf("%d. %s\n", s.Step, s.Action))
				if s.Description != "" {
					b.WriteString(fmt.Sprintf("   - %s\n", s.Description))
				}
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}

func ExportHTML(output *VerificationOutput) string {
	var b strings.Builder

	b.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	b.WriteString("<meta charset=\"utf-8\">\n")
	b.WriteString("<title>Verification Intelligence Report</title>\n")
	b.WriteString("<style>\n")
	b.WriteString("body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; max-width: 960px; margin: 0 auto; padding: 20px; color: #333; }\n")
	b.WriteString("h1 { color: #111; border-bottom: 2px solid #e74c3c; padding-bottom: 10px; }\n")
	b.WriteString("h2 { color: #222; margin-top: 30px; }\n")
	b.WriteString("h3 { color: #333; margin-top: 20px; }\n")
	b.WriteString("table { border-collapse: collapse; width: 100%; margin: 15px 0; }\n")
	b.WriteString("th, td { border: 1px solid #ddd; padding: 8px 12px; text-align: left; }\n")
	b.WriteString("th { background: #f5f5f5; }\n")
	b.WriteString(".critical { color: #e74c3c; font-weight: bold; }\n")
	b.WriteString(".high { color: #e67e22; font-weight: bold; }\n")
	b.WriteString(".medium { color: #f1c40f; font-weight: bold; }\n")
	b.WriteString(".low { color: #27ae60; }\n")
	b.WriteString(".verified { color: #27ae60; }\n")
	b.WriteString(".unverified { color: #e74c3c; }\n")
	b.WriteString(".plan { border: 1px solid #ddd; border-radius: 6px; padding: 15px; margin: 10px 0; background: #fafafa; }\n")
	b.WriteString(".evidence { margin: 10px 0; }\n")
	b.WriteString(".evidence li { margin: 4px 0; }\n")
	b.WriteString(".actions { margin: 10px 0; padding-left: 20px; }\n")
	b.WriteString(".actions li { margin: 4px 0; }\n")
	b.WriteString(".stakeholder { color: #666; font-style: italic; }\n")
	b.WriteString(".summary-box { background: #f8f9fa; border: 1px solid #dee2e6; border-radius: 8px; padding: 20px; margin: 20px 0; }\n")
	b.WriteString(".gap-list li { margin: 4px 0; }\n")
	b.WriteString("</style>\n</head>\n<body>\n")

	b.WriteString("<h1>Verification Intelligence Report</h1>\n")
	if output.Domain != "" {
		b.WriteString(fmt.Sprintf("<p><strong>Domain:</strong> %s</p>\n", output.Domain))
	}
	b.WriteString(fmt.Sprintf("<p><em>Generated: %s</em></p>\n", output.GeneratedAt))

	if output.Assessment != nil && len(output.Assessment.Plans) > 0 {
		b.WriteString("<h2>Executive Summary</h2>\n")
		b.WriteString("<div class=\"summary-box\">\n")
		b.WriteString("<table>\n")
		b.WriteString("<tr><th>Metric</th><th>Value</th></tr>\n")
		b.WriteString(fmt.Sprintf("<tr><td>Total Assumptions</td><td>%d</td></tr>\n", output.Assessment.TotalAssumptions))
		b.WriteString(fmt.Sprintf("<tr><td>Verified</td><td>%d</td></tr>\n", output.Assessment.VerifiedCount))
		b.WriteString(fmt.Sprintf("<tr><td>Partially Verified</td><td>%d</td></tr>\n", output.Assessment.PartialCount))
		b.WriteString(fmt.Sprintf("<tr><td>Unverified</td><td>%d</td></tr>\n", output.Assessment.UnverifiedCount))
		b.WriteString(fmt.Sprintf("<tr><td>No Evidence</td><td>%d</td></tr>\n", output.Assessment.NoEvidenceCount))
		b.WriteString(fmt.Sprintf("<tr><td>Overall Confidence</td><td>%.1f%%</td></tr>\n", output.Assessment.OverallConfidence))
		b.WriteString("</table>\n</div>\n")

		if output.CISOView != nil {
			b.WriteString("<h2>CISO Review</h2>\n")
			b.WriteString(fmt.Sprintf("<p>Verification Backlog: <strong>%d</strong> items</p>\n", len(output.CISOView.VerificationBacklog)))
			b.WriteString(fmt.Sprintf("<p>Evidence Gaps: <strong>%d</strong></p>\n", len(output.CISOView.EvidenceGaps)))
			b.WriteString(fmt.Sprintf("<p>Priority Distribution: Critical=%d, High=%d, Medium=%d, Low=%d</p>\n",
				output.CISOView.CriticalCount, output.CISOView.HighCount,
				output.CISOView.MediumCount, output.CISOView.LowCount))

			if len(output.CISOView.HighestRiskUnverified) > 0 {
				b.WriteString("<h3>Highest Risk Unverified</h3>\n<ul>\n")
				for _, p := range output.CISOView.HighestRiskUnverified {
					b.WriteString(fmt.Sprintf("<li class=\"%s\">%s [%s]</li>\n",
						toLower(string(p.Priority)), p.AssumptionText, string(p.Priority)))
				}
				b.WriteString("</ul>\n")
			}
		}

		b.WriteString("<h2>Verification Plans</h2>\n")
		for i, p := range output.Assessment.Plans {
			statusClass := "unverified"
			if p.Status == VsVerified {
				statusClass = "verified"
			}
			priorityClass := toLower(string(p.Priority))

			b.WriteString(fmt.Sprintf("<div class=\"plan\">\n"))
			b.WriteString(fmt.Sprintf("<h3>%d. %s</h3>\n", i+1, p.AssumptionText))
			b.WriteString(fmt.Sprintf("<p><strong>Category:</strong> %s | <strong>Risk:</strong> %s | ", string(p.Category), p.Risk))
			b.WriteString(fmt.Sprintf("<strong>Confidence:</strong> <span class=\"%s\">%.0f/100</span> | ", statusClass, p.Confidence))
			b.WriteString(fmt.Sprintf("<strong>Priority:</strong> <span class=\"%s\">%s</span> | ", priorityClass, string(p.Priority)))
			b.WriteString(fmt.Sprintf("<strong>Status:</strong> <span class=\"%s\">%s</span> | ", statusClass, string(p.Status)))
			b.WriteString(fmt.Sprintf("<strong>Effort:</strong> %s | <strong>Time:</strong> %s</p>\n", string(p.Effort), p.ExpectedTime))

			if p.WhyVerify != "" {
				b.WriteString(fmt.Sprintf("<p><strong>Why Verify:</strong> %s</p>\n", p.WhyVerify))
			}

			if len(p.EvidenceRequired) > 0 {
				b.WriteString("<p><strong>Required Evidence:</strong></p>\n<ul class=\"evidence\">\n")
				for _, ev := range p.EvidenceRequired {
					opt := ""
					if ev.Optional {
						opt = " (optional)"
					}
					b.WriteString(fmt.Sprintf("<li>%s — %s%s</li>\n", ev.Name, ev.Description, opt))
				}
				b.WriteString("</ul>\n")
			}

			if len(p.Actions) > 0 {
				b.WriteString("<p><strong>Verification Actions:</strong></p>\n<ol class=\"actions\">\n")
				for _, a := range p.Actions {
					b.WriteString(fmt.Sprintf("<li>%s — %s <span class=\"stakeholder\">(%s)</span></li>\n", a.Action, a.Description, a.Stakeholder))
				}
				b.WriteString("</ol>\n")
			}

			if p.WhatToReview != "" {
				b.WriteString(fmt.Sprintf("<p><strong>What to Review:</strong> %s</p>\n", p.WhatToReview))
			}
			if p.WhatEvidenceToCollect != "" {
				b.WriteString(fmt.Sprintf("<p><strong>What Evidence to Collect:</strong> %s</p>\n", p.WhatEvidenceToCollect))
			}
			if p.HowToValidate != "" {
				b.WriteString(fmt.Sprintf("<p><strong>How to Validate:</strong> %s</p>\n", p.HowToValidate))
			}

			b.WriteString("</div>\n")
		}
	}

	b.WriteString("</body>\n</html>")
	return b.String()
}

func ExportJSON(output *VerificationOutput) string {
	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"error":"%s"}`, err.Error())
	}
	return string(data)
}
