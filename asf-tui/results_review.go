package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderReviewQueue(s StyleSet, result *AnalysisResult) string {
	if result.ReviewOutput == nil || result.ReviewOutput.Queue == nil {
		return dimStyle.Render("No review data available.")
	}
	q := result.ReviewOutput.Queue
	rows := []string{
		emphasisStyle.Render(fmt.Sprintf("Review Queue — %d items", len(q.Items))),
		"",
	}
	if len(q.Items) == 0 {
		rows = append(rows, dimStyle.Render("No items in review queue."))
		return strings.Join(rows, "\n")
	}
	rows = append(rows, fmt.Sprintf("  %-5s %-6s %-12s %-16s %s", "Rank", "Score", "Risk", "Category", "Assumption"))
	rows = append(rows, fmt.Sprintf("  %s", strings.Repeat("─", 80)))
	for _, item := range q.Items {
		riskStyle := s.StatusWarn
		switch item.Risk {
		case "Critical":
			riskStyle = s.StatusBad
		case "High":
			riskStyle = s.StatusWarn
		case "Medium":
			riskStyle = s.StatusGood
		}
		text := item.AssumptionText
		if len(text) > 50 {
			text = text[:47] + "..."
		}
		rows = append(rows, fmt.Sprintf("  %-5d %-6.0f %s %-16s %s",
			item.Rank, item.PriorityScore,
			riskStyle.Render(item.Risk),
			item.Category, text))
	}
	return strings.Join(rows, "\n")
}

func renderReviewMatrix(s StyleSet, result *AnalysisResult) string {
	if result.ReviewOutput == nil || result.ReviewOutput.Matrix == nil {
		return dimStyle.Render("No priority matrix available.")
	}
	m := result.ReviewOutput.Matrix
	rows := []string{emphasisStyle.Render("Priority Matrix — 2×2 Quadrant View"), ""}

	if len(m.HighValueLowEffort) > 0 {
		rows = append(rows, s.StatusGood.Render("★ High Value / Low Effort (Do First):"))
		for _, item := range m.HighValueLowEffort {
			rows = append(rows, fmt.Sprintf("    • [%s] %s (Score: %.0f)", item.Risk, item.AssumptionText, item.PriorityScore))
		}
		rows = append(rows, "")
	}
	if len(m.HighValueHighEffort) > 0 {
		rows = append(rows, s.StatusWarn.Render("◆ High Value / High Effort (Plan):"))
		for _, item := range m.HighValueHighEffort {
			rows = append(rows, fmt.Sprintf("    • [%s] %s (Score: %.0f)", item.Risk, item.AssumptionText, item.PriorityScore))
		}
		rows = append(rows, "")
	}
	if len(m.LowValueLowEffort) > 0 {
		rows = append(rows, dimStyle.Render("○ Low Value / Low Effort (Quick Wins):"))
		for _, item := range m.LowValueLowEffort {
			rows = append(rows, fmt.Sprintf("    • [%s] %s (Score: %.0f)", item.Risk, item.AssumptionText, item.PriorityScore))
		}
		rows = append(rows, "")
	}
	if len(m.LowValueHighEffort) > 0 {
		rows = append(rows, lipgloss.NewStyle().Render("○ Low Value / High Effort (Avoid):"))
		for _, item := range m.LowValueHighEffort {
			rows = append(rows, fmt.Sprintf("    • [%s] %s (Score: %.0f)", item.Risk, item.AssumptionText, item.PriorityScore))
		}
		rows = append(rows, "")
	}
	if len(m.HighValueLowEffort)+len(m.HighValueHighEffort)+len(m.LowValueLowEffort)+len(m.LowValueHighEffort) == 0 {
		rows = append(rows, dimStyle.Render("No quadrant classification available."))
	}
	return strings.Join(rows, "\n")
}

func renderReviewCampaigns(s StyleSet, result *AnalysisResult) string {
	if result.ReviewOutput == nil || len(result.ReviewOutput.Campaigns) == 0 {
		return dimStyle.Render("No review campaigns available.")
	}
	r := result.ReviewOutput
	rows := []string{emphasisStyle.Render(fmt.Sprintf("Review Campaigns — %d campaigns", len(r.Campaigns))), ""}

	for _, c := range r.Campaigns {
		rows = append(rows, s.StatusGood.Render(fmt.Sprintf("▶ %s (%d items)", c.Name, len(c.Items))))
		if c.TotalEffort != "" {
			rows = append(rows, fmt.Sprintf("  Duration: %s | Effort: %s", c.Duration, c.TotalEffort))
		}
		for i, item := range c.Items {
			if i >= 5 {
				remaining := len(c.Items) - i
				rows = append(rows, dimStyle.Render(fmt.Sprintf("  ... and %d more", remaining)))
				break
			}
			text := item.AssumptionText
			if len(text) > 60 {
				text = text[:57] + "..."
			}
			rows = append(rows, fmt.Sprintf("    %d. [%s] %s", i+1, item.Risk, text))
		}
		rows = append(rows, "")
	}
	return strings.Join(rows, "\n")
}

func renderCISOReviewDashboard(s StyleSet, result *AnalysisResult) string {
	if result.ReviewOutput == nil || result.ReviewOutput.CISODashboard == nil {
		return dimStyle.Render("No CISO review dashboard available.")
	}
	d := result.ReviewOutput.CISODashboard
	rows := []string{emphasisStyle.Render("CISO Review Dashboard"), ""}

	rows = append(rows, fmt.Sprintf("Total Assumptions: %d", d.TotalAssumptions))
	rows = append(rows, fmt.Sprintf("Critical:          %d", d.CriticalAssumptions))
	rows = append(rows, fmt.Sprintf("High:              %d", d.HighAssumptions))
	rows = append(rows, "")

	if len(d.HighestRiskAssumptions) > 0 {
		rows = append(rows, s.StatusBad.Render("Highest Risk Assumptions:"))
		for _, item := range d.HighestRiskAssumptions {
			text := item.AssumptionText
			if len(text) > 50 {
				text = text[:47] + "..."
			}
			rows = append(rows, fmt.Sprintf("  ⚠ [%s] %s (Score: %.0f)", item.Risk, text, item.PriorityScore))
		}
		rows = append(rows, "")
	}
	if len(d.GreatestRiskReduction) > 0 {
		rows = append(rows, s.StatusGood.Render("Greatest Risk Reduction Opportunities:"))
		for _, item := range d.GreatestRiskReduction {
			text := item.AssumptionText
			if len(text) > 50 {
				text = text[:47] + "..."
			}
			rows = append(rows, fmt.Sprintf("  ✓ [%s] %s (Score: %.0f)", item.Risk, text, item.PriorityScore))
		}
		rows = append(rows, "")
	}
	return strings.Join(rows, "\n")
}

func renderConfidenceView(s StyleSet, result *AnalysisResult) string {
	if result.ConfidenceOutput == nil || len(result.ConfidenceOutput.Breakdowns) == 0 {
		return dimStyle.Render("No confidence data available.")
	}
	c := result.ConfidenceOutput
	rows := []string{
		emphasisStyle.Render(fmt.Sprintf("Confidence View — %d assumptions", len(c.Breakdowns))),
		"",
	}
	for _, bd := range c.Breakdowns {
		confStyle := s.StatusWarn
		switch {
		case bd.FinalConfidence >= 70:
			confStyle = s.StatusGood
		case bd.FinalConfidence < 40:
			confStyle = s.StatusBad
		}
		text := bd.AssumptionText
		if len(text) > 55 {
			text = text[:52] + "..."
		}
		rows = append(rows, fmt.Sprintf("  %s %s %s",
			confStyle.Render(fmt.Sprintf("%.0f%%", bd.FinalConfidence)),
			emphasisStyle.Render(string(bd.StabilityClass)),
			text))
	}
	if c.CISOTrustView != nil {
		rows = append(rows, "")
		rows = append(rows, s.StatusGood.Render(fmt.Sprintf("Most Trusted: %d items", len(c.CISOTrustView.MostTrustedFindings))))
		rows = append(rows, s.StatusBad.Render(fmt.Sprintf("Least Trusted: %d items", len(c.CISOTrustView.LeastTrustedFindings))))
		rows = append(rows, s.StatusWarn.Render(fmt.Sprintf("Critical Low-Confidence: %d items", len(c.CISOTrustView.MostCriticalLowConfidence))))
	}
	return strings.Join(rows, "\n")
}

func renderExplainabilityView(s StyleSet, result *AnalysisResult) string {
	if result.ConfidenceOutput == nil || len(result.ConfidenceOutput.Breakdowns) == 0 {
		return dimStyle.Render("No explainability data available.")
	}
	c := result.ConfidenceOutput
	rows := []string{
		emphasisStyle.Render("Explainability View — Per-Assumption Rationale"),
		"",
	}
	for _, bd := range c.Breakdowns {
		text := bd.AssumptionText
		if len(text) > 50 {
			text = text[:47] + "..."
		}
		rows = append(rows, emphasisStyle.Render(fmt.Sprintf("  %s [%.0f%%]", text, bd.FinalConfidence)))
		rows = append(rows, fmt.Sprintf("    Why: %s", bd.WhyExists))
		rows = append(rows, fmt.Sprintf("    Uncertain: %s", bd.WhyUncertain))
		if bd.WhatIncreasesConfidence != "" {
			rows = append(rows, fmt.Sprintf("    Increases: %s", bd.WhatIncreasesConfidence))
		}
		rows = append(rows, "")
	}
	return strings.Join(rows, "\n")
}

func renderConfidenceBreakdownView(s StyleSet, result *AnalysisResult) string {
	if result.ConfidenceOutput == nil || len(result.ConfidenceOutput.Breakdowns) == 0 {
		return dimStyle.Render("No confidence breakdown available.")
	}
	c := result.ConfidenceOutput
	rows := []string{
		emphasisStyle.Render("Confidence Breakdown — Factor Analysis"),
		"",
	}
	for _, bd := range c.Breakdowns {
		text := bd.AssumptionText
		if len(text) > 50 {
			text = text[:47] + "..."
		}
		rows = append(rows, emphasisStyle.Render(fmt.Sprintf("  %s [%.0f%%] (%s)", text, bd.FinalConfidence, bd.StabilityClass)))

		for _, f := range bd.PositiveFactors {
			rows = append(rows, fmt.Sprintf("    + %s: +%.0f", f.Name, f.Impact))
		}
		for _, f := range bd.NegativeFactors {
			rows = append(rows, fmt.Sprintf("    - %s: %.0f", f.Name, f.Impact))
		}
		for _, fc := range bd.SupportingFacts {
			sign := "+"
			if !fc.IsPositive {
				sign = ""
			}
			rows = append(rows, fmt.Sprintf("    Fact %s: %s (%s%.1f%%)", fc.FactID, fc.FactText, sign, fc.Contribution))
		}
		if bd.DomainContribution != nil {
			rows = append(rows, fmt.Sprintf("    Domain (%s): +%.0f%% (%s)", bd.DomainContribution.Domain, bd.DomainContribution.Influence, bd.DomainContribution.Strength))
		}
		rows = append(rows, "")
	}

	if c.CISOTrustView != nil {
		rows = append(rows, emphasisStyle.Render("  CISO Trust View"))
		if len(c.CISOTrustView.MostTrustedFindings) > 0 {
			rows = append(rows, s.StatusGood.Render("  Most Trusted:"))
			for _, f := range c.CISOTrustView.MostTrustedFindings[:min(3, len(c.CISOTrustView.MostTrustedFindings))] {
				ft := f.AssumptionText
				if len(ft) > 50 {
					ft = ft[:47] + "..."
				}
				rows = append(rows, fmt.Sprintf("    [%.0f%%] %s", f.FinalConfidence, ft))
			}
		}
		if len(c.CISOTrustView.LeastTrustedFindings) > 0 {
			rows = append(rows, s.StatusBad.Render("  Least Trusted:"))
			for _, f := range c.CISOTrustView.LeastTrustedFindings[:min(3, len(c.CISOTrustView.LeastTrustedFindings))] {
				ft := f.AssumptionText
				if len(ft) > 50 {
					ft = ft[:47] + "..."
				}
				rows = append(rows, fmt.Sprintf("    [%.0f%%] %s", f.FinalConfidence, ft))
			}
		}
	}

	if c.ArchitectReviewView != nil {
		rows = append(rows, "")
		rows = append(rows, emphasisStyle.Render("  Architect Review View"))
		rows = append(rows, fmt.Sprintf("    Requiring Validation: %d", len(c.ArchitectReviewView.RequiringValidation)))
		rows = append(rows, fmt.Sprintf("    Weak Support: %d", len(c.ArchitectReviewView.WeakSupport)))
		rows = append(rows, fmt.Sprintf("    Strong Support: %d", len(c.ArchitectReviewView.StrongSupport)))
	}

	return strings.Join(rows, "\n")
}

var (
	dimStyle      = lipgloss.NewStyle().Faint(true)
	emphasisStyle = lipgloss.NewStyle().Bold(true)
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
