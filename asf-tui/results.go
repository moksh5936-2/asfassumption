package main

import (
	"fmt"
	"sort"
	"strings"

	"asf-tui/asf/trust"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type resultsModel struct {
	result         *AnalysisResult
	sections       []string
	selected       int
	expanded       map[int]bool
	exportComplete bool
	exportPath     string
	exportFormat   ExportFormat
}

func newResultsModel() resultsModel {
	return resultsModel{
		sections: []string{
			"Assumptions",
			"Critical Assumptions",
			"Risk Matrix",
			"STRIDE Distribution",
			"Recommended Controls",
			"Attack Paths",
			"Security Design Review",
			"Compliance",
			"Compliance Intelligence",
			"Domain Knowledge",
			"Executive Risk Narratives",
			"Portfolio Intelligence",
			"Decision Intelligence",
			"Digital Twin",
			"Architect Narrative",
			"Dependency Graph",
			"Trust Chains",
			"Critical Assumptions",
			"Failure Cascades",
			"Single Points of Trust Failure",
			"Trust Collapse Simulation",
			"Coverage Dashboard",
			"Blind Spot View",
			"Coverage Heatmap",
			"Architect Attention Score",
			"Verification View",
			"Evidence View",
			"Verification Priority",
			"Verification Roadmap",
			"CISO Verification Summary",
			"Review Queue",
			"Review Priority Matrix",
			"Review Campaigns",
			"CISO Review Dashboard",
			"Confidence View",
			"Explainability View",
			"Confidence Breakdown",
		},
		expanded: map[int]bool{},
	}
}

func (m resultsModel) Update(msg tea.Msg) (resultsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.sections)-1 {
				m.selected++
			}
		case "enter":
			m.expanded[m.selected] = !m.expanded[m.selected]
		case "e":
			if m.result != nil {
				return m, func() tea.Msg { return navigateMsg{to: exportView} }
			}
		case "r":
			// Navigate to review mode (handled by mainModel)
		}
	}
	return m, nil
}

func (m mainModel) viewResults() string {
	s := m.styles
	r := m.results

	if r.result == nil {
		return lipgloss.JoinVertical(lipgloss.Left,
			s.Title.Render("Analysis Results"),
			s.Subtitle.Render("No results available. Run an analysis first."),
		)
	}

	result := r.result
	aiBadge := ""
	if result.AnalysisMode == ModeASFAndAI {
		aiCount := 0
		for _, a := range result.Assumptions {
			if strings.HasPrefix(a.ID, "AI-") {
				aiCount++
			}
		}
		if aiCount > 0 {
			aiBadge = " " + s.StatusGood.Render(fmt.Sprintf("+%d AI", aiCount))
		} else {
			aiBadge = " " + s.StatusGood.Render("+AI")
		}
	}
	header := lipgloss.JoinVertical(lipgloss.Left,
		s.Title.Render(fmt.Sprintf("Results: %s", result.ArchitectureName)),
		s.Subtitle.Render(fmt.Sprintf("Mode: %s%s  |  Date: %s", result.AnalysisMode, aiBadge, result.AnalysisDate.Format("Jan 2, 2006 15:04"))),
		s.Subtitle.Render(fmt.Sprintf("Total Assumptions: %d  |  Critical: %d  |  High: %d  |  Medium: %d  |  Low: %d",
			result.TotalAssumptions, result.CriticalCount, result.HighCount, result.MediumCount, result.LowCount)),
	)

	if r.exportComplete {
		header = lipgloss.JoinVertical(lipgloss.Left,
			header,
			s.StatusGood.Render(fmt.Sprintf("✓ Exported: %s", r.exportPath)),
		)
	}

	var sectionViews []string
	for i, section := range r.sections {
		prefix := "  "
		style := s.SectionItem
		if i == r.selected {
			prefix = "▸ "
			style = s.MenuSelected
		}

		expanded := r.expanded[i]
		expandMarker := "[+]"
		if expanded {
			expandMarker = "[-]"
		}

		headerLine := style.Render(fmt.Sprintf("%s %s %s", prefix, expandMarker, section))
		sectionViews = append(sectionViews, headerLine)

		if expanded {
			content := renderSectionContent(s, i, result)
			sectionViews = append(sectionViews,
				s.BorderBox.Render(content),
			)
		}
	}

	if result.AnalysisMode == ModeASFAndAI {
		aiAssumptions := renderAIAssumptions(s, result.Assumptions)
		if aiAssumptions != "" {
			sectionViews = append(sectionViews, "", s.Section.Render("AI-Enhanced Findings"), s.BorderBox.Render(aiAssumptions))
		}
	}

	body := lipgloss.JoinVertical(lipgloss.Left, sectionViews...)

	return lipgloss.JoinVertical(lipgloss.Left, header, body)
}

func renderSectionContent(s StyleSet, section int, result *AnalysisResult) string {
	switch section {
	case 0:
		return renderAssumptions(s, result.Assumptions)
	case 1:
		return renderCriticalAssumptions(s, result.Assumptions)
	case 2:
		return renderRiskMatrix(s, result)
	case 3:
		return renderStrideDist(s, result)
	case 4:
		return renderControls(s, result.Controls)
	case 5:
		return renderAttackPaths(s, result)
	case 6:
		return renderSecurityDesignReview(s, result)
	case 7:
		return renderCompliance(s, result.Compliance)
	case 8:
		return renderComplianceIntelligence(s, result)
	case 9:
		return renderDKPI(s, result)
	case 10:
		return renderERN(s, result)
	case 11:
		return renderSAMPI(s, result)
	case 12:
		return renderSDI(s, result)
	case 13:
		return renderSDT(s, result)
	case 14:
		return renderArchitectNarrative(s, result)
	case 15:
		return renderDependencyGraph(s, result)
	case 16:
		return renderTrustChains(s, result)
	case 17:
		return renderTrustCriticalAssumptions(s, result)
	case 18:
		return renderFailureCascades(s, result)
	case 19:
		return renderSinglePointsOfTrust(s, result)
	case 20:
		return renderTrustCollapse(s, result)
	case 21:
		return renderCoverageDashboard(s, result)
	case 22:
		return renderBlindSpotView(s, result)
	case 23:
		return renderCoverageHeatmap(s, result)
	case 24:
		return renderAttentionScore(s, result)
	case 25:
		return renderVerificationView(s, result)
	case 26:
		return renderEvidenceView(s, result)
	case 27:
		return renderVerificationPriority(s, result)
	case 28:
		return renderVerificationRoadmap(s, result)
	case 29:
		return renderCISOVerificationSummary(s, result)
	case 30:
		return renderReviewQueue(s, result)
	case 31:
		return renderReviewMatrix(s, result)
	case 32:
		return renderReviewCampaigns(s, result)
	case 33:
		return renderCISOReviewDashboard(s, result)
	case 34:
		return renderConfidenceView(s, result)
	case 35:
		return renderExplainabilityView(s, result)
	case 36:
		return renderConfidenceBreakdownView(s, result)
	}
	return ""
}

func renderAIAssumptions(s StyleSet, assumptions []Assumption) string {
	var items []string
	for _, a := range assumptions {
		if !strings.HasPrefix(a.ID, "AI-") {
			continue
		}
		riskStyle := riskStyle(s, a.Risk)
		dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
		strideStr := ""
		if len(a.Stride) > 0 {
			strs := make([]string, len(a.Stride))
			for i, st := range a.Stride {
				strs[i] = string(st)
			}
			strideStr = " [" + strings.Join(strs, ", ") + "]"
		}
		items = append(items, fmt.Sprintf("%s %s%s — %s %s",
			s.StatusGood.Render("AI"),
			riskStyle.Render(string(a.Risk)),
			dim.Render(strideStr),
			a.Description,
			dim.Render(a.Category),
		))
	}
	return strings.Join(items, "\n")
}

func renderAssumptions(s StyleSet, assumptions []Assumption) string {
	var items []string
	for _, a := range assumptions {
		riskStyle := riskStyle(s, a.Risk)
		dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)

		line := fmt.Sprintf("%s [%s] %s — %s",
			riskStyle.Render(string(a.Risk)),
			dim.Render(a.ID),
			a.Description,
			dim.Render(a.Component),
		)

		if a.Confidence > 0 {
			confPct := int(a.Confidence * 100)
			confStyle := s.Value
			if confPct >= 80 {
				confStyle = s.StatusGood
			} else if confPct >= 50 {
				confStyle = s.StatusWarn
			}
			line += fmt.Sprintf(" %s", confStyle.Render(fmt.Sprintf("(%d%% conf)", confPct)))
		}

		if a.Rationale != "" && len(items) < 10 {
			line += "\n  " + dim.Render(a.Rationale)
		}

		items = append(items, line)
	}
	return strings.Join(items, "\n")
}

func renderCriticalAssumptions(s StyleSet, assumptions []Assumption) string {
	var items []string
	for _, a := range assumptions {
		if a.Risk == RiskCritical {
			strideStr := ""
			if len(a.Stride) > 0 {
				strideStr = " [" + string(a.Stride[0]) + "]"
			}
			line := fmt.Sprintf("⚠ %s%s — %s (L:%d I:%d",
				a.Description, strideStr, a.Component, a.Likelihood, a.Impact)
			if a.RiskJustification != nil {
				line += fmt.Sprintf(" Score:%d)", a.RiskJustification.RiskScore)
			} else {
				line += ")"
			}
			if a.Rationale != "" {
				line += "\n  " + lipgloss.NewStyle().Foreground(s.Theme().DimText).Render(a.Rationale)
			}
			items = append(items, line)
		}
	}
	if len(items) == 0 {
		return "No critical assumptions found."
	}
	return strings.Join(items, "\n")
}

func renderRiskMatrix(s StyleSet, result *AnalysisResult) string {
	labels := []RiskLevel{RiskCritical, RiskHigh, RiskMedium, RiskLow}
	counts := map[RiskLevel]int{
		RiskCritical: result.CriticalCount,
		RiskHigh:     result.HighCount,
		RiskMedium:   result.MediumCount,
		RiskLow:      result.LowCount,
	}
	var rows []string

	// Header
	rows = append(rows, fmt.Sprintf("Risk Model: %s", s.Value.Render(result.RiskModelVersion)))
	if result.ConfidenceSummary != "" {
		rows = append(rows, s.SectionItem.Render(result.ConfidenceSummary))
	}
	rows = append(rows, "")

	// 5x5 matrix visualization
	rows = append(rows, "  Likelihood × Impact = Risk Score")
	rows = append(rows, "  1-4: Low | 5-11: Medium | 12-19: High | 20-25: Critical")
	rows = append(rows, "")
	matrixHeader := "       | I:1  I:2  I:3  I:4  I:5"
	rows = append(rows, s.Value.Render(matrixHeader))
	rows = append(rows, s.Value.Render("  -----+-----------------------"))
	for lh := 5; lh >= 1; lh-- {
		cell := ""
		for im := 1; im <= 5; im++ {
			score := lh * im
			r := riskForScore(score)
			marker := "■"
			rs := riskStyle(s, r)
			cell += rs.Render(fmt.Sprintf(" %s%-2d", marker, score))
		}
		rows = append(rows, fmt.Sprintf("  L:%d  |%s", lh, cell))
	}
	rows = append(rows, "")

	for _, label := range labels {
		count := counts[label]
		bar := strings.Repeat("■", count)
		if count > 20 {
			bar = strings.Repeat("■", 20)
		}
		style := riskStyle(s, label)
		rows = append(rows, fmt.Sprintf("%s %s (%d)",
			style.Render(padRight(string(label), 10)),
			style.Render(bar),
			count,
		))
	}
	return strings.Join(rows, "\n")
}

func riskForScore(score int) RiskLevel {
	switch {
	case score >= 20:
		return RiskCritical
	case score >= 12:
		return RiskHigh
	case score >= 5:
		return RiskMedium
	default:
		return RiskLow
	}
}

func renderStrideDist(s StyleSet, result *AnalysisResult) string {
	total := 0
	for _, count := range result.StrideDistribution {
		total += count
	}
	if total == 0 {
		return "No STRIDE data available."
	}
	categories := []StrideCategory{
		StrideSpoofing, StrideTampering, StrideRepudiation,
		StrideInfoDisclosure, StrideDenialOfService, StrideElevationPriv,
	}
	var rows []string
	for _, cat := range categories {
		count := result.StrideDistribution[cat]
		pct := float64(count) / float64(total) * 100
		barCount := int(pct / 5)
		if barCount > 20 {
			barCount = 20
		}
		bar := strings.Repeat("▨", barCount)
		if pct > 0 && barCount == 0 {
			bar = "▨"
		}
		rows = append(rows, fmt.Sprintf("%s %s %d (%.1f%%)",
			s.Value.Render(padRight(string(cat), 25)),
			s.Progress.Render(bar),
			count, pct,
		))
	}
	return strings.Join(rows, "\n")
}

func renderControls(s StyleSet, controls []ControlDetail) string {
	var items []string
	for _, c := range controls {
		dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
		line := fmt.Sprintf("✓ %s: %s", c.ID, c.Description)
		items = append(items, line)
		items = append(items, "  "+dim.Render(c.Rationale))
		if len(c.MitigatedAssumptionIDs) > 0 {
			ids := c.MitigatedAssumptionIDs
			if len(ids) > 3 {
				ids = ids[:3]
			}
			items = append(items, "  "+dim.Render("Assumptions: "+strings.Join(ids, ", ")))
		}
		if len(c.MitigatedSTRIDE) > 0 {
			strs := make([]string, len(c.MitigatedSTRIDE))
			for i, s := range c.MitigatedSTRIDE {
				strs[i] = string(s)
			}
			items = append(items, "  "+dim.Render("STRIDE: "+strings.Join(strs, ", ")))
		}
	}
	return strings.Join(items, "\n")
}

func renderAttackPaths(s StyleSet, result *AnalysisResult) string {
	if len(result.AttackPaths) == 0 {
		return "No attack paths discovered."
	}
	aps := result.AttackPathSummary
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	val := s.Value

	header := renderAttackPathSummary(s, aps)
	header += "\n"

	var rows []string
	rows = append(rows, val.Render(fmt.Sprintf("Top %d Attack Paths:", len(result.AttackPaths))))
	rows = append(rows, "")
	for _, p := range result.AttackPaths {
		riskLabel := riskLevelForAPDScore(p.RiskScore)
		riskStyle := riskStyleForLabel(s, riskLabel)
		line := fmt.Sprintf("%s %s → %s  %s  %s",
			val.Render(p.Name),
			dim.Render(p.EntryPoint),
			dim.Render(p.TargetAsset),
			riskStyle.Render(fmt.Sprintf("%.2f (%s)", p.RiskScore, riskLabel)),
			dim.Render(p.DetectionDifficulty),
		)
		rows = append(rows, line)
		if p.BusinessImpact != "" {
			rows = append(rows, "  "+dim.Render("Impact: "+p.BusinessImpact))
		}
		if len(p.KillChainPhases) > 0 {
			rows = append(rows, "  "+dim.Render("Kill Chain: "+strings.Join(p.KillChainPhases, ", ")))
		}
		if len(p.MITREATTACK) > 0 {
			rows = append(rows, "  "+dim.Render("MITRE: "+strings.Join(p.MITREATTACK, ", ")))
		}
		rows = append(rows, "")
	}

	return header + strings.Join(rows, "\n")
}

func renderAttackPathSummary(s StyleSet, aps AttackPathSummary) string {
	val := s.Value
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	var parts []string
	parts = append(parts, val.Render(fmt.Sprintf("Attack Paths: %d", aps.TotalAttackPaths)))
	parts = append(parts, val.Render(fmt.Sprintf("Threat Chains: %d", aps.ThreatChainCount)))
	parts = append(parts, dim.Render(fmt.Sprintf("C:%d H:%d M:%d L:%d", aps.CriticalCount, aps.HighCount, aps.MediumCount, aps.LowCount)))
	if len(aps.KillChainCoverage) > 0 {
		parts = append(parts, dim.Render(fmt.Sprintf("Kill Chain: %d phases", len(aps.KillChainCoverage))))
	}
	if len(aps.MITRECoverage) > 0 {
		parts = append(parts, dim.Render(fmt.Sprintf("MITRE: %d techniques", len(aps.MITRECoverage))))
	}
	return strings.Join(parts, " | ")
}

func riskStyleForLabel(s StyleSet, label string) lipgloss.Style {
	switch label {
	case "Critical":
		return s.StatusBad
	case "High":
		return s.StatusWarn
	case "Medium":
		return s.Value
	case "Low":
		return s.StatusGood
	}
	return s.Value
}

func renderSecurityDesignReview(s StyleSet, result *AnalysisResult) string {
	if len(result.SDRIDesignFindings) == 0 && len(result.SDRIControls) == 0 {
		return "No security design review data available."
	}
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	val := s.Value
	var rows []string

	rows = append(rows, val.Render("=== Security Design Review ==="))
	rows = append(rows, "")

	if result.SDRISummary != "" {
		rows = append(rows, dim.Render(result.SDRISummary))
		rows = append(rows, "")
	}

	if len(result.SDRICoverageDashboard) > 0 {
		rows = append(rows, val.Render("Control Coverage:"))
		cats := make([]string, 0, len(result.SDRICoverageDashboard))
		for cat := range result.SDRICoverageDashboard {
			cats = append(cats, cat)
		}
		sort.Strings(cats)
		for _, cat := range cats {
			cov := result.SDRICoverageDashboard[cat]
			level := coverageLevelString(cov)
			style := s.StatusGood
			if level == "Fair" || level == "Poor" {
				style = s.StatusWarn
			}
			rows = append(rows, fmt.Sprintf("  %s: %s %.1f%%",
				dim.Render(padRight(cat, 25)),
				style.Render(level),
				cov,
			))
		}
		rows = append(rows, "")
	}

	if len(result.SDRIDesignFindings) > 0 {
		rows = append(rows, val.Render(fmt.Sprintf("Findings (%d):", len(result.SDRIDesignFindings))))
		for _, f := range result.SDRIDesignFindings {
			sevStyle := riskStyleForLabel(s, f.Severity)
			rows = append(rows, fmt.Sprintf("  %s %s — %s",
				sevStyle.Render(f.Severity),
				s.SectionItem.Render(f.Title),
				dim.Render(f.Category),
			))
			rows = append(rows, "    "+dim.Render(f.BusinessImpact))
			if len(f.AffectedComponents) > 0 {
				rows = append(rows, "    "+dim.Render("Components: "+strings.Join(f.AffectedComponents, ", ")))
			}
		}
		rows = append(rows, "")
	}

	if len(result.SDRIAchitecturalWeaknesses) > 0 {
		rows = append(rows, val.Render(fmt.Sprintf("Architectural Weaknesses (%d):", len(result.SDRIAchitecturalWeaknesses))))
		for _, w := range result.SDRIAchitecturalWeaknesses {
			sevStyle := riskStyleForLabel(s, w.Severity)
			rows = append(rows, fmt.Sprintf("  %s %s — %s",
				sevStyle.Render(w.Severity),
				s.SectionItem.Render(w.Pattern),
				dim.Render(w.Description),
			))
		}
		rows = append(rows, "")
	}

	if len(result.SDRIRemediations) > 0 {
		rows = append(rows, val.Render(fmt.Sprintf("Top Remediations (%d):", len(result.SDRIRemediations))))
		for i, r := range result.SDRIRemediations {
			if i >= 5 {
				break
			}
			rows = append(rows, fmt.Sprintf("  #%d %s (risk: %.2f) — %s",
				r.Priority, s.SectionItem.Render(r.Description), r.RiskScore, dim.Render(r.Effort)))
		}
		rows = append(rows, "")
	}

	if len(result.SDRIComplianceAlignments) > 0 {
		rows = append(rows, val.Render("Compliance Alignment:"))
		for _, m := range result.SDRIComplianceAlignments {
			statusStyle := s.StatusGood
			if m.Status == "Fair" || m.Status == "Poor" {
				statusStyle = s.StatusWarn
			}
			rows = append(rows, fmt.Sprintf("  %s: %s %.1f%% (%d controls)",
				dim.Render(padRight(m.Framework, 12)),
				statusStyle.Render(m.Status),
				m.Coverage,
				len(m.Controls),
			))
		}
	}

	return strings.Join(rows, "\n")
}

func coverageLevelString(cov float64) string {
	switch {
	case cov >= 90:
		return "Excellent"
	case cov >= 75:
		return "Strong"
	case cov >= 50:
		return "Good"
	case cov >= 25:
		return "Fair"
	default:
		return "Poor"
	}
}

func renderCompliance(s StyleSet, compliance []string) string {
	if len(compliance) == 0 {
		return "No compliance frameworks specified."
	}
	var items []string
	for i, c := range compliance {
		if i == 0 {
			items = append(items, s.Value.Render(c))
		} else {
			items = append(items, "▸ "+c)
		}
	}
	return strings.Join(items, "\n")
}

func renderComplianceIntelligence(s StyleSet, result *AnalysisResult) string {
	if len(result.CIAREFrameworkCoverages) == 0 {
		return "No compliance intelligence data available."
	}
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	val := s.Value
	var rows []string

	rows = append(rows, val.Render("=== Compliance Intelligence ==="))
	rows = append(rows, "")

	// Framework Coverage
	rows = append(rows, val.Render("Framework Coverage:"))
	for _, c := range result.CIAREFrameworkCoverages {
		style := s.StatusGood
		if c.Status == "Fair" || c.Status == "Weak" || c.Status == "Poor" {
			style = s.StatusWarn
		}
		rows = append(rows, fmt.Sprintf("  %s: %s %.1f%% (%d/%d controls)",
			dim.Render(padRight(c.Framework, 12)),
			style.Render(c.Status),
			c.CoveragePct, c.Observed, c.Required,
		))
	}
	rows = append(rows, "")

	// Audit Readiness
	if len(result.CIAREAuditReadiness) > 0 {
		rows = append(rows, val.Render("Audit Readiness:"))
		for _, a := range result.CIAREAuditReadiness {
			arStyle := s.StatusGood
			if a.Status == "Fair" || a.Status == "Weak" || a.Status == "Poor" {
				arStyle = s.StatusWarn
			}
			rows = append(rows, fmt.Sprintf("  %s: %s %.1f%%",
				dim.Render(padRight(a.Framework, 12)),
				arStyle.Render(a.Status),
				a.ReadinessScore,
			))
		}
		rows = append(rows, "")
	}

	// Top Compliance Gaps
	if len(result.CIAREComplianceGaps) > 0 {
		rows = append(rows, val.Render(fmt.Sprintf("Top Compliance Gaps (%d shown):", minInt(5, len(result.CIAREComplianceGaps)))))
		for i, g := range result.CIAREComplianceGaps {
			if i >= 5 {
				break
			}
			gStyle := s.StatusWarn
			if g.Risk == "Critical" {
				gStyle = s.StatusBad
			}
			rows = append(rows, fmt.Sprintf("  %s %s — %s",
				gStyle.Render(g.Risk),
				s.SectionItem.Render(g.Requirement),
				dim.Render(g.Framework),
			))
		}
		rows = append(rows, "")
	}

	// Missing Evidence
	if len(result.CIAREMissingEvidences) > 0 {
		rows = append(rows, val.Render(fmt.Sprintf("Missing Evidence (%d):", len(result.CIAREMissingEvidences))))
		for i, m := range result.CIAREMissingEvidences {
			if i >= 5 {
				break
			}
			rows = append(rows, fmt.Sprintf("  %s — %s",
				s.SectionItem.Render(m.Control),
				dim.Render(m.Framework),
			))
		}
		rows = append(rows, "")
	}

	// Narratives
	if len(result.CIAREComplianceNarratives) > 0 {
		rows = append(rows, val.Render("Compliance Narratives:"))
		for _, n := range result.CIAREComplianceNarratives {
			rows = append(rows, dim.Render(n.Narrative))
		}
	}

	return strings.Join(rows, "\n")
}

func renderDKPI(s StyleSet, result *AnalysisResult) string {
	if result.DKPI.DomainResult.PrimaryDomain == "" {
		return "No domain knowledge data available."
	}
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	val := s.Value
	d := result.DKPI
	var rows []string

	rows = append(rows, val.Render("=== Domain Knowledge Intelligence ==="))
	rows = append(rows, "")
	rows = append(rows, fmt.Sprintf("Domain: %s  Confidence: %.1f%%",
		val.Render(d.DomainResult.PrimaryDomain), d.DomainResult.Confidence))
	rows = append(rows, dim.Render(d.Summary))
	rows = append(rows, "")

	if len(d.DomainResult.Rationale) > 0 {
		rows = append(rows, val.Render("Detection Rationale:"))
		for _, r := range d.DomainResult.Rationale {
			rows = append(rows, "  "+dim.Render("▸ "+r))
		}
		rows = append(rows, "")
	}

	if len(d.Recommendations) > 0 {
		rows = append(rows, val.Render(fmt.Sprintf("Domain Recommendations (%d):", len(d.Recommendations))))
		for _, rec := range d.Recommendations {
			rows = append(rows, "  "+dim.Render("▸ "+rec))
		}
		rows = append(rows, "")
	}

	if len(d.InjectedThreats) > 0 {
		rows = append(rows, val.Render(fmt.Sprintf("Domain Threats (%d):", len(d.InjectedThreats))))
		for _, t := range d.InjectedThreats {
			sevStyle := riskStyleForLabel(s, string(t.Severity))
			rows = append(rows, fmt.Sprintf("  %s %s — %s",
				sevStyle.Render(string(t.Severity)),
				s.SectionItem.Render(t.Name),
				dim.Render(t.Description)))
		}
		rows = append(rows, "")
	}

	if len(d.DomainControls) > 0 {
		rows = append(rows, val.Render(fmt.Sprintf("Domain Controls (%d):", len(d.DomainControls))))
		for _, c := range d.DomainControls {
			rows = append(rows, fmt.Sprintf("  %s — %s (%s/%s)",
				s.SectionItem.Render(c.Name), dim.Render(c.Category), c.Coverage, c.Status))
		}
		rows = append(rows, "")
	}

	if len(d.DomainCompliance) > 0 {
		rows = append(rows, val.Render("Domain Compliance:"))
		for _, f := range d.DomainCompliance {
			rows = append(rows, "  "+dim.Render("▸ "+f))
		}
	}

	return strings.Join(rows, "\n")
}

func renderERN(s StyleSet, result *AnalysisResult) string {
	if len(result.ERN.ExecutiveRisks) == 0 {
		return "No executive risk data available."
	}
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	val := s.Value
	ern := result.ERN
	var rows []string

	rows = append(rows, val.Render("=== Executive Risk Narratives ==="))
	rows = append(rows, "")
	rows = append(rows, fmt.Sprintf("Financial Exposure: %s", val.Render(ern.FinancialExposure.Level)))
	rows = append(rows, dim.Render(ern.FinancialExposure.Rationale))
	rows = append(rows, "")

	if ern.BoardSummary.Summary != "" {
		rows = append(rows, val.Render("Board Summary:"))
		rows = append(rows, "  "+dim.Render(ern.BoardSummary.Summary))
		rows = append(rows, "")
	}

	if len(ern.RiskThemes) > 0 {
		rows = append(rows, val.Render(fmt.Sprintf("Risk Themes (%d):", len(ern.RiskThemes))))
		for _, th := range ern.RiskThemes {
			sevStyle := riskStyleForLabel(s, th.Severity)
			rows = append(rows, fmt.Sprintf("  %s %s (%d findings, %s)",
				sevStyle.Render("▸"),
				th.Name, th.RiskCount, th.Severity))
		}
		rows = append(rows, "")
	}

	if len(ern.ExecutiveRisks) > 0 {
		rows = append(rows, val.Render(fmt.Sprintf("Executive Risks (%d):", len(ern.ExecutiveRisks))))
		for _, risk := range ern.ExecutiveRisks {
			priStyle := riskStyleForLabel(s, risk.Priority)
			rows = append(rows, fmt.Sprintf("  %s %s — %s",
				priStyle.Render("["+risk.Priority+"]"),
				s.SectionItem.Render(risk.Title),
				dim.Render(risk.BusinessImpact)))
		}
		rows = append(rows, "")
	}

	if len(ern.CISOBriefing.TopRisks) > 0 {
		rows = append(rows, val.Render("Top Risks (CISO Briefing):"))
		for _, r := range ern.CISOBriefing.TopRisks {
			rows = append(rows, "  "+dim.Render("▸ "+r))
		}
		rows = append(rows, "")
	}

	if len(ern.RemediationRoadmap.Phase30) > 0 {
		rows = append(rows, val.Render("Remediation Roadmap — 30 Days:"))
		for _, item := range ern.RemediationRoadmap.Phase30 {
			rows = append(rows, "  "+dim.Render("▸ ["+item.Priority+"] "+item.Action))
		}
		rows = append(rows, "")
	}

	if len(ern.InvestmentInsights) > 0 {
		rows = append(rows, val.Render("Security Investment Insights:"))
		for _, ii := range ern.InvestmentInsights {
			rows = append(rows, fmt.Sprintf("  %s — %s",
				s.SectionItem.Render(ii.Area),
				dim.Render("["+ii.Priority+"] "+ii.Rationale)))
		}
		rows = append(rows, "")
	}

	if ern.DecisionSupport.Top3Actions != nil && len(ern.DecisionSupport.Top3Actions) > 0 {
		rows = append(rows, val.Render("CISO Decision Support — Top 3 Actions:"))
		for _, da := range ern.DecisionSupport.Top3Actions {
			rows = append(rows, fmt.Sprintf("  %d. %s (%s)",
				da.Rank, s.SectionItem.Render(da.Action), da.Impact))
		}
		rows = append(rows, "")
	}

	rows = append(rows, val.Render("Report Packs:"))
	if ern.ReportPacks.BoardReport != "" {
		rows = append(rows, dim.Render("  Board Report: ✓ available"))
	}
	if ern.ReportPacks.ExecutiveReport != "" {
		rows = append(rows, dim.Render("  Executive Report: ✓ available"))
	}
	if ern.ReportPacks.TechnicalReport != "" {
		rows = append(rows, dim.Render("  Technical Report: ✓ available"))
	}

	return strings.Join(rows, "\n")
}

func renderSAMPI(s StyleSet, result *AnalysisResult) string {
	sampi := result.SAMPI
	if sampi.Dashboard.TotalArchitectures == 0 {
		return "No portfolio data available. Use 'asf portfolio add <file>' to build a portfolio."
	}
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	val := s.Value
	var rows []string

	rows = append(rows, val.Render("=== Portfolio Intelligence ==="))
	rows = append(rows, "")
	rows = append(rows, fmt.Sprintf("Architectures: %d", sampi.Dashboard.TotalArchitectures))
	rows = append(rows, fmt.Sprintf("Total Findings: %d", sampi.Dashboard.TotalFindings))
	rows = append(rows, fmt.Sprintf("Total Threats: %d", sampi.Dashboard.TotalThreats))
	rows = append(rows, fmt.Sprintf("Total Controls: %d", sampi.Dashboard.TotalControls))
	rows = append(rows, fmt.Sprintf("Average Risk Score: %.1f", sampi.Dashboard.AverageRiskScore))
	rows = append(rows, fmt.Sprintf("Average Coverage: %.1f%%", sampi.Dashboard.AverageCoverage))
	rows = append(rows, "")

	if len(sampi.Heatmaps) > 0 {
		rows = append(rows, val.Render("Executive Heatmap:"))
		for _, h := range sampi.Heatmaps {
			bandStyle := riskStyleForLabel(s, h.RiskBand)
			rows = append(rows, fmt.Sprintf("  %s %s [%s] score=%.1f",
				bandStyle.Render("▸"), h.ArchitectureName, h.RiskBand, h.RiskScore))
		}
		rows = append(rows, "")
	}

	if len(sampi.RepeatedWeaknesses) > 0 {
		rows = append(rows, val.Render(fmt.Sprintf("Repeated Weaknesses (%d):", len(sampi.RepeatedWeaknesses))))
		for _, rw := range sampi.RepeatedWeaknesses {
			mark := " "
			if rw.Systemic {
				mark = dim.Render("⚠")
			}
			rows = append(rows, fmt.Sprintf("  %s %s (%d architectures)",
				mark, rw.FindingTitle, rw.OccurrenceCount))
		}
		rows = append(rows, "")
	}

	if len(sampi.EnterpriseThemes) > 0 {
		rows = append(rows, val.Render("Enterprise Risk Themes:"))
		for _, th := range sampi.EnterpriseThemes {
			sevStyle := riskStyleForLabel(s, th.Severity)
			rows = append(rows, fmt.Sprintf("  %s %s (%d) — %s",
				sevStyle.Render("▸"), th.Name, th.RiskCount, th.Severity))
		}
		rows = append(rows, "")
	}

	if len(sampi.ControlCoverage) > 0 {
		rows = append(rows, val.Render("Control Coverage:"))
		for _, cc := range sampi.ControlCoverage {
			mark := dim.Render(" ")
			if cc.CoveragePercent < 50 {
				mark = s.StatusWarn.Render("!")
			}
			rows = append(rows, fmt.Sprintf("  %s %s: %.1f%%", mark, cc.ControlName, cc.CoveragePercent))
		}
		rows = append(rows, "")
	}

	if sampi.SecurityDebt.Score > 0 {
		rows = append(rows, val.Render(fmt.Sprintf("Security Debt Score: %.1f", sampi.SecurityDebt.Score)))
		rows = append(rows, dim.Render(fmt.Sprintf("  Longstanding: %d  Repeated: %d",
			sampi.SecurityDebt.LongstandingCount, sampi.SecurityDebt.RepeatedCount)))
		rows = append(rows, "")
	}

	if len(sampi.ProgramInsights) > 0 {
		rows = append(rows, val.Render("Security Program Insights:"))
		for _, pi := range sampi.ProgramInsights {
			priStyle := riskStyleForLabel(s, pi.Priority)
			rows = append(rows, fmt.Sprintf("  %s [%s] %s",
				priStyle.Render("▸"), pi.Priority, pi.Insight))
		}
		rows = append(rows, "")
	}

	return strings.Join(rows, "\n")
}

func renderSDI(s StyleSet, result *AnalysisResult) string {
	sdi := result.SDI
	if len(sdi.Recommendations) == 0 {
		return "No decision intelligence data available. Run an analysis to generate recommendations."
	}
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	val := s.Value
	prio := s.StatusBad
	warn := s.StatusWarn
	var rows []string

	rows = append(rows, val.Render("=== Decision Intelligence ==="))
	rows = append(rows, "")

	rows = append(rows, val.Render(fmt.Sprintf("Top %d Recommendations:", len(sdi.Recommendations))))
	for _, r := range sdi.Recommendations {
		var ps lipgloss.Style
		switch r.Priority {
		case "Critical":
			ps = prio
		case "High":
			ps = warn
		default:
			ps = val
		}
		roi := r.RiskReduction + "/" + r.Effort
		rows = append(rows, fmt.Sprintf("  %s [%s] %s (%s) — ROI: %s",
			ps.Render("▸"), r.Priority, r.Title, r.ID, roi))
		rows = append(rows, dim.Render(fmt.Sprintf("    → %d findings, %d threats, Impact: %s",
			len(r.AffectedFindings), len(r.AffectedThreats), r.BusinessImpact)))
	}
	rows = append(rows, "")

	if len(sdi.Dashboard.QuickWins) > 0 {
		rows = append(rows, val.Render("Quick Wins:"))
		for _, qw := range sdi.Dashboard.QuickWins {
			rows = append(rows, fmt.Sprintf("  %s %s (Low effort, %s priority)",
				s.StatusGood.Render("✓"), qw.Title, qw.Priority))
		}
		rows = append(rows, "")
	}

	if len(sdi.FixSimulations) > 0 {
		rows = append(rows, val.Render("Fix Simulations (what if implemented?):"))
		for _, sim := range sdi.FixSimulations {
			rows = append(rows, fmt.Sprintf("  %s: Critical %d→%d, High %d→%d, Coverage %.0f%%→%.0f%%",
				sim.ControlName, sim.OriginalCritical, sim.NewCritical,
				sim.OriginalHigh, sim.NewHigh, sim.OriginalCoverage, sim.NewCoverage))
		}
		rows = append(rows, "")
	}

	if len(sdi.FailureSimulations) > 0 {
		rows = append(rows, warn.Render("Failure Simulations (what if control fails?):"))
		for _, sim := range sdi.FailureSimulations {
			rows = append(rows, fmt.Sprintf("  %s: %d systems impacted, %d new paths, Risk: %s",
				sim.ControlName, sim.SystemsImpacted, sim.AttackPathsOpened, sim.RiskIncrease))
		}
		rows = append(rows, "")
	}

	if len(sdi.InvestmentPriorities) > 0 {
		rows = append(rows, val.Render("Investment Priorities:"))
		for _, ip := range sdi.InvestmentPriorities {
			rows = append(rows, fmt.Sprintf("  #%d %s (score=%.1f, %d findings)",
				ip.Rank, ip.Area, ip.Score, ip.FindingCount))
		}
		rows = append(rows, "")
	}

	if sdi.RemediationRoadmap.Phase30 != nil {
		rows = append(rows, val.Render("Remediation Roadmap:"))
		if len(sdi.RemediationRoadmap.Phase30) > 0 {
			rows = append(rows, dim.Render("  30 Days:"))
			for _, item := range sdi.RemediationRoadmap.Phase30 {
				rows = append(rows, fmt.Sprintf("    %s [%s]", item.Action, item.Priority))
			}
		}
		if len(sdi.RemediationRoadmap.Phase90) > 0 {
			rows = append(rows, dim.Render("  90 Days:"))
			for _, item := range sdi.RemediationRoadmap.Phase90 {
				rows = append(rows, fmt.Sprintf("    %s [%s]", item.Action, item.Priority))
			}
		}
		rows = append(rows, "")
	}

	if sdi.Dashboard.RiskReductionSummary != "" {
		rows = append(rows, dim.Render(sdi.Dashboard.RiskReductionSummary))
		rows = append(rows, "")
	}

	if sdi.ExecutiveScenarios.BestCase.Scenario != "" {
		rows = append(rows, val.Render("Executive Scenarios:"))
		rows = append(rows, fmt.Sprintf("  %s: %.1f risk score, %d findings resolved",
			s.StatusGood.Render("Best Case"), sdi.ExecutiveScenarios.BestCase.RiskScore,
			sdi.ExecutiveScenarios.BestCase.FindingsResolved))
		rows = append(rows, fmt.Sprintf("  %s: %.1f risk score, %d findings resolved",
			warn.Render("Likely Case"), sdi.ExecutiveScenarios.LikelyCase.RiskScore,
			sdi.ExecutiveScenarios.LikelyCase.FindingsResolved))
		rows = append(rows, fmt.Sprintf("  %s: %.1f risk score, %d findings resolved",
			prio.Render("Worst Case"), sdi.ExecutiveScenarios.WorstCase.RiskScore,
			sdi.ExecutiveScenarios.WorstCase.FindingsResolved))
		rows = append(rows, "")
	}

	return strings.Join(rows, "\n")
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func riskStyle(s StyleSet, risk RiskLevel) lipgloss.Style {
	switch risk {
	case RiskCritical:
		return s.StatusBad
	case RiskHigh:
		return s.StatusWarn
	case RiskMedium:
		return s.Value
	case RiskLow:
		return s.StatusGood
	}
	return s.Value
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}

func exportFormatFromConfig(cfg *Config) ExportFormat {
	switch cfg.Output.Default {
	case "html":
		return ExportHTML
	case "pdf":
		return ExportPDF
	case "csv":
		return ExportCSV
	case "json":
		return ExportJSON
	case "narrative-md":
		return ExportNarrativeMarkdown
	case "narrative-html":
		return ExportNarrativeHTML
	default:
		return ExportMarkdown
	}
}

func (m mainModel) updateResults(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	if m.analyze.result != nil {
		m.results.result = m.analyze.result
		m.results.exportFormat = exportFormatFromConfig(m.config)
		m.results.exportComplete = false
		m.results.exportPath = ""
		m.results.expanded = map[int]bool{}
		m.analyze.result = nil
	}
	m.results, cmd = m.results.Update(msg)
	return m, cmd
}

func renderSDT(s StyleSet, result *AnalysisResult) string {
	sdt := result.SDT
	if sdt.Twin.ID == "" {
		return "No digital twin data available. Run an analysis to generate a twin."
	}
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	val := s.Value
	warn := s.StatusWarn
	good := s.StatusGood
	var rows []string

	rows = append(rows, val.Render("=== Security Digital Twin ==="))
	rows = append(rows, "")

	rows = append(rows, val.Render(fmt.Sprintf("Twin: %s (v%s)", sdt.Twin.ArchitectureName, sdt.Twin.Version)))
	rows = append(rows, dim.Render(fmt.Sprintf("  Risk: %.1f  |  Coverage: %.0f%%  |  Hash: %s",
		sdt.Twin.RiskScore, sdt.Twin.Coverage, sdt.Twin.SourceHash)))
	rows = append(rows, "")

	if len(sdt.ChangeImpacts) > 0 {
		rows = append(rows, val.Render("Change Impact Analysis:"))
		for _, ci := range sdt.ChangeImpacts {
			cs := val
			if ci.Severity == "Critical" || ci.Severity == "High" {
				cs = warn
			}
			rows = append(rows, fmt.Sprintf("  %s %s → %s (%s) — %d risks, %d paths",
				cs.Render("▸"), ci.Change, ci.ComponentAffected, ci.Severity, ci.RisksAffected, ci.AttackPathsAffected))
		}
		rows = append(rows, "")
	}

	if len(sdt.ControlDrifts) > 0 {
		rows = append(rows, warn.Render("Control Drift:"))
		for _, cd := range sdt.ControlDrifts {
			rows = append(rows, fmt.Sprintf("  %s: %s → %s (%s)",
				cd.ControlName, cd.ExpectedState, cd.CurrentState, cd.RiskImpact))
		}
		rows = append(rows, "")
	}

	if sdt.SecurityDebt.TotalDebt > 0 {
		rows = append(rows, val.Render(fmt.Sprintf("Security Debt: %.0f (findings: %.0f, controls: %.0f, assumptions: %.0f)",
			sdt.SecurityDebt.TotalDebt, sdt.SecurityDebt.FindingDebt,
			sdt.SecurityDebt.ControlDebt, sdt.SecurityDebt.AssumptionDebt)))
		rows = append(rows, "")
	}

	if len(sdt.ComplianceDrifts) > 0 {
		rows = append(rows, val.Render("Compliance Drift:"))
		for _, cd := range sdt.ComplianceDrifts {
			cs := good
			if cd.Status == "Regressed" || cd.NewGaps > 0 {
				cs = warn
			}
			rows = append(rows, fmt.Sprintf("  %s %s — %d new gaps, %d resolved",
				cs.Render(cd.Framework), cd.Status, cd.NewGaps, cd.ResolvedGaps))
		}
		rows = append(rows, "")
	}

	if len(sdt.WhatIfScenarios) > 0 {
		rows = append(rows, val.Render("What-If Scenarios:"))
		for _, wi := range sdt.WhatIfScenarios {
			ds := good
			if wi.RiskDelta > 0 {
				ds = warn
			}
			rows = append(rows, fmt.Sprintf("  %s %s: risk %.1f, coverage %.0f%%, %d findings",
				ds.Render("▸"), wi.Name, wi.RiskDelta, wi.CoverageDelta, wi.FindingsDelta))
		}
		rows = append(rows, "")
	}

	if len(sdt.ZeroTrust.Dimensions) > 0 {
		rows = append(rows, val.Render("Zero Trust Assessment:"))
		rows = append(rows, dim.Render(fmt.Sprintf("  Overall: %.1f / %.1f (gap: %.1f)",
			sdt.ZeroTrust.Overall, sdt.ZeroTrust.Target, sdt.ZeroTrust.Gap)))
		for _, d := range sdt.ZeroTrust.Dimensions {
			rows = append(rows, fmt.Sprintf("  %s: %.1f / %.1f", d.Dimension, d.Score, d.Target))
		}
		rows = append(rows, "")
	}

	if len(sdt.CrownJewels) > 0 {
		rows = append(rows, val.Render("Crown Jewel Analysis:"))
		for _, cj := range sdt.CrownJewels {
			rows = append(rows, fmt.Sprintf("  %s — Business: %s, Attack: %s, Blast: %s (score: %.1f)",
				cj.AssetName, cj.BusinessValue, cj.AttackValue, cj.BlastRadius, cj.OverallScore))
		}
		rows = append(rows, "")
	}

	if sdt.ExecutiveReport.ArchitectureHealth != "" {
		rows = append(rows, val.Render("Executive Report:"))
		rows = append(rows, fmt.Sprintf("  Health: %s | Debt: %.0f | Drifts: %d control, %d compliance",
			sdt.ExecutiveReport.ArchitectureHealth, sdt.ExecutiveReport.SecurityDebtScore,
			sdt.ExecutiveReport.ControlDriftCount, sdt.ExecutiveReport.ComplianceDriftCount))
		rows = append(rows, fmt.Sprintf("  Trend: %s / %s",
			sdt.ExecutiveReport.RiskTrend, sdt.ExecutiveReport.AttackSurfaceTrend))
		rows = append(rows, "")
	}

	if sdt.PortfolioSummary.ArchitectureCount > 1 {
		rows = append(rows, val.Render(fmt.Sprintf("Portfolio: %d architectures, debt: %.0f",
			sdt.PortfolioSummary.ArchitectureCount, sdt.PortfolioSummary.AggregatedDebt)))
		if len(sdt.PortfolioSummary.EnterpriseTrends) > 0 {
			for _, et := range sdt.PortfolioSummary.EnterpriseTrends {
				rows = append(rows, fmt.Sprintf("  • %s", et))
			}
		}
		rows = append(rows, "")
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func renderArchitectNarrative(s StyleSet, result *AnalysisResult) string {
	if result.NarrativeOutput == nil {
		return "Narrative output not generated."
	}

	var rows []string
	n := result.NarrativeOutput

	// Overview
	rows = append(rows, s.Section.Render("Architecture Overview"))
	rows = append(rows, fmt.Sprintf("Domain: %s | Components: %d | Assumptions: %d (Critical: %d, High: %d)",
		n.ArchitectureOverview.Domain, n.ArchitectureOverview.TotalComponents,
		n.ArchitectureOverview.TotalAssumptions, n.ArchitectureOverview.CriticalCount,
		n.ArchitectureOverview.HighCount))
	rows = append(rows, "")

	// Key components
	if len(n.ArchitectureOverview.KeyComponents) > 0 {
		rows = append(rows, fmt.Sprintf("Key Components: %s", strings.Join(n.ArchitectureOverview.KeyComponents, ", ")))
		rows = append(rows, "")
	}

	// Most critical assumptions
	if len(n.ExecutiveReport.MostCriticalAssumptions) > 0 {
		rows = append(rows, s.Section.Render("Most Critical Assumptions"))
		for _, a := range n.ExecutiveReport.MostCriticalAssumptions {
			riskStyle := riskStyle(s, RiskLevel(a.RiskLevel))
			rows = append(rows, fmt.Sprintf("%s %s", riskStyle.Render(string(a.RiskLevel)), a.Text))
			rows = append(rows, fmt.Sprintf("  Impact: %s", a.BusinessImpact))
		}
		rows = append(rows, "")
	}

	// Assumption narratives (first 3)
	if len(n.AssumptionNarratives) > 0 {
		rows = append(rows, s.Section.Render("Assumption Narratives"))
		showCount := len(n.AssumptionNarratives)
		if showCount > 3 {
			showCount = 3
		}
		for i := 0; i < showCount; i++ {
			an := n.AssumptionNarratives[i]
			riskStyle := riskStyle(s, RiskLevel(an.RiskLevel))
			rows = append(rows, fmt.Sprintf("%s %s", riskStyle.Render(an.RiskLevel), an.AssumptionText))
			rows = append(rows, fmt.Sprintf("  Context: %s", an.Context))
			rows = append(rows, fmt.Sprintf("  Consequence: %s", an.FailureConsequence))
			rows = append(rows, fmt.Sprintf("  Recommendation: %s", an.SecurityRecommendation))
			rows = append(rows, "")
		}
		if len(n.AssumptionNarratives) > 3 {
			rows = append(rows, fmt.Sprintf("... and %d more narratives", len(n.AssumptionNarratives)-3))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func renderDependencyGraph(s StyleSet, result *AnalysisResult) string {
	if result.TrustOutput == nil || result.TrustOutput.DependencyGraph == nil {
		return "No dependency graph data available. Run an analysis first."
	}
	g := result.TrustOutput.DependencyGraph
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	val := s.Value
	var rows []string

	rows = append(rows, val.Render("=== Dependency Graph ==="))
	rows = append(rows, "")
	rows = append(rows, fmt.Sprintf("Nodes: %d  |  Edges: %d", len(g.Nodes), len(g.Edges)))
	rows = append(rows, "")

	typeCounts := make(map[trust.DependencyType]int)
	for _, e := range g.Edges {
		typeCounts[e.DependencyType]++
	}
	for _, dt := range trust.AllDependencyTypes {
		if count := typeCounts[dt]; count > 0 {
			rows = append(rows, fmt.Sprintf("  %s: %d", dim.Render(string(dt)), count))
		}
	}
	rows = append(rows, "")

	var criticalNodes []string
	for _, node := range g.Nodes {
		if node.Criticality >= 0.8 || node.Centrality >= 0.5 {
			criticalNodes = append(criticalNodes, fmt.Sprintf("%s (%s, cent: %.2f)", node.ID, node.Risk, node.Centrality))
		}
	}
	if len(criticalNodes) > 0 {
		rows = append(rows, val.Render(fmt.Sprintf("High-Value Nodes (%d):", len(criticalNodes))))
		for _, n := range criticalNodes {
			rows = append(rows, "  "+dim.Render(n))
		}
	}

	return strings.Join(rows, "\n")
}

func renderTrustChains(s StyleSet, result *AnalysisResult) string {
	if result.TrustOutput == nil || len(result.TrustOutput.TrustChains) == 0 {
		return "No trust chains discovered."
	}
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	val := s.Value
	var rows []string

	rows = append(rows, val.Render(fmt.Sprintf("=== Trust Chains (%d) ===", len(result.TrustOutput.TrustChains))))
	rows = append(rows, "")

	for i, chain := range result.TrustOutput.TrustChains {
		if i >= 10 {
			rows = append(rows, dim.Render(fmt.Sprintf("... and %d more chains", len(result.TrustOutput.TrustChains)-10)))
			break
		}
		chainRS := riskStyle(s, RiskLevel(chain.Risk))
		rows = append(rows, fmt.Sprintf("%s Chain %s (len=%d, conf=%.2f, risk=%s)",
			val.Render("◆"), chain.ID, chain.Length, chain.Confidence, chainRS.Render(chain.Risk)))

		for j, nodeID := range chain.Nodes {
			if node, ok := result.TrustOutput.DependencyGraph.Nodes[nodeID]; ok {
				arrow := "→"
				if j == 0 {
					arrow = "●"
				}
				nr := riskStyle(s, RiskLevel(node.Risk))
				rows = append(rows, fmt.Sprintf("  %s %s %s", arrow, nr.Render(node.Risk), dim.Render(node.Text)))
			}
		}
		rows = append(rows, "")
	}

	return strings.Join(rows, "\n")
}

func renderTrustCriticalAssumptions(s StyleSet, result *AnalysisResult) string {
	if result.TrustOutput == nil || len(result.TrustOutput.CriticalAssumptions) == 0 {
		return "No critical assumptions detected."
	}
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	val := s.Value
	var rows []string

	rows = append(rows, val.Render(fmt.Sprintf("=== Critical Assumptions (%d) ===", len(result.TrustOutput.CriticalAssumptions))))
	rows = append(rows, "")

	for i, ca := range result.TrustOutput.CriticalAssumptions {
		if i >= 20 {
			break
		}
		riskSty := riskStyle(s, RiskLevel(ca.Risk))
		rows = append(rows, fmt.Sprintf("%s %s [centrality=%.2f, support=%d, failure_radius=%d, trust_radius=%d]",
			riskSty.Render("⚠"), ca.AssumptionText, ca.Centrality, ca.SupportCount, ca.FailureRadius, ca.TrustRadius))
		rows = append(rows, dim.Render(fmt.Sprintf("   Score: %.1f | Dependencies: %s", ca.Score, ca.DependencyTypes)))
		rows = append(rows, "")
	}

	return strings.Join(rows, "\n")
}

func renderFailureCascades(s StyleSet, result *AnalysisResult) string {
	if result.TrustOutput == nil || len(result.TrustOutput.FailureCascades) == 0 {
		return "No failure cascades generated."
	}
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	val := s.Value
	var rows []string

	rows = append(rows, val.Render(fmt.Sprintf("=== Failure Cascades (%d) ===", len(result.TrustOutput.FailureCascades))))
	rows = append(rows, "")

	for i, cascade := range result.TrustOutput.FailureCascades {
		if i >= 10 {
			rows = append(rows, dim.Render(fmt.Sprintf("... and %d more cascades", len(result.TrustOutput.FailureCascades)-10)))
			break
		}
		sevStyle := riskStyle(s, RiskLevel(cascade.Severity))
		rows = append(rows, fmt.Sprintf("%s %s", sevStyle.Render("⬡"), val.Render(cascade.RootAssumptionText)))
		rows = append(rows, dim.Render(fmt.Sprintf("  Severity: %s | Affected: %d | Max Depth: %d",
			cascade.Severity, cascade.TotalAffected, cascade.MaxDepth)))
		for _, step := range cascade.Steps {
			stepStyle := riskStyle(s, RiskLevel(step.Severity))
			rows = append(rows, fmt.Sprintf("  %d. %s %s", step.Step, stepStyle.Render("["+step.Severity+"]"), dim.Render(step.AssumptionText)))
		}
		rows = append(rows, "")
	}

	return strings.Join(rows, "\n")
}

func renderSinglePointsOfTrust(s StyleSet, result *AnalysisResult) string {
	if result.TrustOutput == nil || len(result.TrustOutput.SinglePointsOfTrust) == 0 {
		return "No single points of trust failure detected."
	}
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	val := s.Value
	var rows []string

	rows = append(rows, val.Render(fmt.Sprintf("=== Single Points of Trust Failure (%d) ===", len(result.TrustOutput.SinglePointsOfTrust))))
	rows = append(rows, "")

	for _, spotf := range result.TrustOutput.SinglePointsOfTrust {
		rows = append(rows, fmt.Sprintf("%s %s", s.StatusBad.Render("⚠"), val.Render(spotf.AssumptionText)))
		rows = append(rows, dim.Render(fmt.Sprintf("  Dependents: %d", spotf.DependentsCount)))
		rows = append(rows, dim.Render(fmt.Sprintf("  Recommendation: %s", spotf.Recommendation)))
		rows = append(rows, "")
	}

	return strings.Join(rows, "\n")
}

func renderTrustCollapse(s StyleSet, result *AnalysisResult) string {
	if result.TrustOutput == nil || len(result.TrustOutput.TrustCollapseResults) == 0 {
		return "No trust collapse simulations available."
	}
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	val := s.Value
	var rows []string

	rows = append(rows, val.Render(fmt.Sprintf("=== Trust Collapse Simulations (%d) ===", len(result.TrustOutput.TrustCollapseResults))))
	rows = append(rows, "")

	for i, collapse := range result.TrustOutput.TrustCollapseResults {
		if i >= 10 {
			rows = append(rows, dim.Render(fmt.Sprintf("... and %d more simulations", len(result.TrustOutput.TrustCollapseResults)-10)))
			break
		}
		rows = append(rows, fmt.Sprintf("%s %s", s.StatusBad.Render("◆"), val.Render(collapse.FailedAssumptionText)))
		rows = append(rows, dim.Render(fmt.Sprintf("  Assumptions Lost: %d", len(collapse.AssumptionsLost))))
		rows = append(rows, dim.Render(fmt.Sprintf("  Risk: %.2f → %.2f (%s)", collapse.RiskScoreBefore, collapse.RiskScoreAfter, collapse.RiskIncrease)))
		if len(collapse.AffectedComponents) > 0 {
			rows = append(rows, dim.Render(fmt.Sprintf("  Affected Components: %s", strings.Join(collapse.AffectedComponents, ", "))))
		}
		rows = append(rows, "")
	}

	return strings.Join(rows, "\n")
}

func renderCoverageDashboard(s StyleSet, result *AnalysisResult) string {
	if result.CoverageOutput == nil {
		return "No coverage analysis available."
	}
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	val := s.Value
	var rows []string

	rows = append(rows, val.Render("=== Coverage Dashboard ==="))
	rows = append(rows, "")

	cov := result.CoverageOutput
	assessment := cov.Assessment
	if assessment == nil {
		return "No coverage assessment data available."
	}

	totalExpected := 0
	totalCovered := 0
	for _, cat := range assessment.Categories {
		totalExpected += cat.ExpectedCount
		totalCovered += cat.ObservedCount
	}
	if totalExpected == 0 {
		return "No coverage expectations defined. Check component taxonomy configuration."
	}
	coveragePct := float64(totalCovered) / float64(totalExpected) * 100

	rows = append(rows, fmt.Sprintf("Total Expected Checks: %d", totalExpected))
	rows = append(rows, fmt.Sprintf("Assumption Coverage: %d/%d (%.1f%%)", totalCovered, totalExpected, coveragePct))
	rows = append(rows, "")
	rows = append(rows, fmt.Sprintf("Coverage Gaps: %d", len(assessment.Gaps)))
	rows = append(rows, fmt.Sprintf("Blind Spots: %d", len(cov.BlindSpots)))
	rows = append(rows, fmt.Sprintf("Domain Blind Spots: %d", len(cov.DomainBlindSpots)))
	rows = append(rows, fmt.Sprintf("Architect Attention Score: %.0f/100", cov.AttentionScore))
	rows = append(rows, "")

	rows = append(rows, dim.Render("By Category:"))
	for _, cat := range assessment.Categories {
		label := string(cat.Category)
		if cat.ObservedCount > cat.ExpectedCount {
			label += " (exceeds)"
		}
		rows = append(rows, fmt.Sprintf("  %-20s %d/%d (%.0f%%)", label, cat.ObservedCount, cat.ExpectedCount, cat.CoveragePct))
	}

	return strings.Join(rows, "\n")
}

func renderBlindSpotView(s StyleSet, result *AnalysisResult) string {
	if result.CoverageOutput == nil {
		return "No coverage analysis available."
	}
	if len(result.CoverageOutput.BlindSpots) == 0 && len(result.CoverageOutput.DomainBlindSpots) == 0 {
		return "No blind spots detected."
	}
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	val := s.Value
	var rows []string

	if len(result.CoverageOutput.BlindSpots) > 0 {
		rows = append(rows, val.Render(fmt.Sprintf("=== Critical Blind Spots (%d) ===", len(result.CoverageOutput.BlindSpots))))
		rows = append(rows, "")
		for _, bs := range result.CoverageOutput.BlindSpots {
			check := fmt.Sprintf("Expected: %s", bs.Title)
			if bs.Component != "" {
				check += fmt.Sprintf(" (component: %s)", bs.Component)
			}
			rows = append(rows, fmt.Sprintf("%s %s", s.StatusBad.Render("✗"), check))
			rows = append(rows, dim.Render(fmt.Sprintf("  Missing: %s", bs.Description)))
			rows = append(rows, dim.Render(fmt.Sprintf("  Risk: %s (Score %.0f)", bs.Risk, bs.Score)))
			rows = append(rows, dim.Render(fmt.Sprintf("  Recommendation: %s", bs.Recommendation)))
			rows = append(rows, "")
		}
	}

	if len(result.CoverageOutput.DomainBlindSpots) > 0 {
		rows = append(rows, val.Render(fmt.Sprintf("=== Domain-Specific Blind Spots (%d) ===", len(result.CoverageOutput.DomainBlindSpots))))
		rows = append(rows, "")
		for _, dbs := range result.CoverageOutput.DomainBlindSpots {
			rows = append(rows, fmt.Sprintf("%s %s", s.StatusWarn.Render("◆"), dbs.MissingArea))
			rows = append(rows, dim.Render(fmt.Sprintf("  Rationale: %s", dbs.Description)))
			rows = append(rows, "")
		}
	}

	return strings.Join(rows, "\n")
}

func renderCoverageHeatmap(s StyleSet, result *AnalysisResult) string {
	if result.CoverageOutput == nil || result.CoverageOutput.Assessment == nil || len(result.CoverageOutput.Assessment.Categories) == 0 {
		return "No coverage data available for heatmap."
	}
	val := s.Value
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	var rows []string

	rows = append(rows, val.Render("=== Coverage Heatmap ==="))
	rows = append(rows, "")

	for _, cat := range result.CoverageOutput.Assessment.Categories {
		pct := cat.CoveragePct

		var bar string
		n := int(pct / 10)
		if n > 10 {
			n = 10
		}
		var barStyle lipgloss.Style
		if pct >= 80 {
			barStyle = s.StatusGood
		} else if pct >= 50 {
			barStyle = s.StatusWarn
		} else {
			barStyle = s.StatusBad
		}
		bar = barStyle.Render(strings.Repeat("█", n)) + dim.Render(strings.Repeat("░", 10-n))

		rows = append(rows, fmt.Sprintf("%-22s %s %5.0f%% (%d/%d)", string(cat.Category), bar, pct, cat.ObservedCount, cat.ExpectedCount))
	}
	rows = append(rows, "")
	rows = append(rows, dim.Render("█ ≥80%  █ 50-79%  █ <50%"))

	return strings.Join(rows, "\n")
}

func renderAttentionScore(s StyleSet, result *AnalysisResult) string {
	if result.CoverageOutput == nil {
		return "No coverage analysis available."
	}
	val := s.Value
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	var rows []string

	score := result.CoverageOutput.AttentionScore

	rows = append(rows, val.Render(fmt.Sprintf("=== Architect Attention Score: %.0f/100 ===", score)))
	rows = append(rows, "")

	var level, explanation string
	switch {
	case score >= 80:
		level = "Low Attention Required"
		explanation = "Coverage is strong across all categories. The architecture is well-documented and assumptions are broadly covered."
	case score >= 60:
		level = "Moderate Attention Required"
		explanation = "Coverage is acceptable but has gaps in specific areas. Review flagged categories for improvement."
	case score >= 40:
		level = "Elevated Attention Required"
		explanation = "Significant coverage gaps exist. Priority attention needed for uncovered high-risk categories."
	case score >= 20:
		level = "High Attention Required"
		explanation = "Critical coverage gaps detected across multiple categories. Architectural review strongly recommended."
	default:
		level = "Critical Attention Required"
		explanation = "Coverage is severely deficient. Immediate architectural review is necessary."
	}

	rows = append(rows, val.Render(level))
	rows = append(rows, "")
	rows = append(rows, dim.Render(explanation))

	if result.CoverageOutput.CISOView != nil {
		rows = append(rows, "")
		gapCount := 0
		if result.CoverageOutput.Assessment != nil {
			gapCount = len(result.CoverageOutput.Assessment.Gaps)
		}
		rows = append(rows, dim.Render(fmt.Sprintf("Gaps: %d | Blind Spots: %d | Domain Blind Spots: %d | Categories: %d",
			gapCount,
			len(result.CoverageOutput.BlindSpots),
			len(result.CoverageOutput.DomainBlindSpots),
			len(result.CoverageOutput.Assessment.Categories))))

		if len(result.CoverageOutput.CISOView.AreasRequiringReview) > 0 {
			rows = append(rows, "")
			rows = append(rows, dim.Render("Key CISO Findings:"))
			for i, f := range result.CoverageOutput.CISOView.AreasRequiringReview {
				if i >= 5 {
					break
				}
				rows = append(rows, fmt.Sprintf("  %s %s", s.StatusWarn.Render("•"), f))
			}
		}
	}

	return strings.Join(rows, "\n")
}

func renderVerificationView(s StyleSet, result *AnalysisResult) string {
	if result.VerificationOutput == nil || result.VerificationOutput.Assessment == nil || len(result.VerificationOutput.Assessment.Plans) == 0 {
		return "No verification plans available."
	}
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	val := s.Value
	var rows []string

	rows = append(rows, val.Render("=== Verification Plans ==="))
	rows = append(rows, "")

	for i, p := range result.VerificationOutput.Assessment.Plans {
		if i >= 20 {
			rows = append(rows, dim.Render(fmt.Sprintf("... and %d more plans", len(result.VerificationOutput.Assessment.Plans)-20)))
			break
		}
		priorityStyle := s.StatusWarn
		if p.Priority == "Critical" || p.Priority == "High" {
			priorityStyle = s.StatusBad
		} else if p.Priority == "Low" {
			priorityStyle = s.StatusGood
		}
		statusStyle := s.StatusWarn
		if p.Status == "Verified" {
			statusStyle = s.StatusGood
		} else if p.Status == "No Evidence" {
			statusStyle = s.StatusBad
		}

		rows = append(rows, fmt.Sprintf("%d. [%s] [%s] %s",
			i+1,
			priorityStyle.Render(string(p.Priority)),
			statusStyle.Render(string(p.Status)),
			val.Render(p.AssumptionText)))
		if p.WhyVerify != "" {
			rows = append(rows, dim.Render(fmt.Sprintf("   Why: %s", p.WhyVerify)))
		}
		rows = append(rows, dim.Render(fmt.Sprintf("   Confidence: %.0f%% | Effort: %s | Time: %s",
			p.Confidence, string(p.Effort), p.ExpectedTime)))
		rows = append(rows, "")
	}

	return strings.Join(rows, "\n")
}

func renderEvidenceView(s StyleSet, result *AnalysisResult) string {
	if result.VerificationOutput == nil || result.VerificationOutput.Assessment == nil || len(result.VerificationOutput.Assessment.Plans) == 0 {
		return "No verification plans available."
	}
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	val := s.Value
	var rows []string

	rows = append(rows, val.Render("=== Evidence View ==="))
	rows = append(rows, "")

	for i, p := range result.VerificationOutput.Assessment.Plans {
		if i >= 15 {
			rows = append(rows, dim.Render(fmt.Sprintf("... and %d more plans", len(result.VerificationOutput.Assessment.Plans)-15)))
			break
		}
		rows = append(rows, val.Render(fmt.Sprintf("%d. %s", i+1, p.AssumptionText)))
		rows = append(rows, dim.Render(fmt.Sprintf("   Category: %s | Risk: %s", string(p.Category), p.Risk)))

		if len(p.EvidenceRequired) > 0 {
			rows = append(rows, dim.Render("   Required Evidence:"))
			for _, ev := range p.EvidenceRequired {
				opt := ""
				if ev.Optional {
					opt = " (optional)"
				}
				rows = append(rows, fmt.Sprintf("     - %s: %s%s", ev.Name, ev.Description, opt))
			}
		}

		if p.WhatEvidenceToCollect != "" {
			rows = append(rows, dim.Render(fmt.Sprintf("   What to Collect: %s", p.WhatEvidenceToCollect)))
		}
		rows = append(rows, "")
	}

	return strings.Join(rows, "\n")
}

func renderVerificationPriority(s StyleSet, result *AnalysisResult) string {
	if result.VerificationOutput == nil || result.VerificationOutput.CISOView == nil {
		return "No verification priority data available."
	}
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	val := s.Value
	var rows []string

	cv := result.VerificationOutput.CISOView
	rows = append(rows, val.Render("=== Verification Priority ==="))
	rows = append(rows, "")
	rows = append(rows, dim.Render(fmt.Sprintf("Critical: %d | High: %d | Medium: %d | Low: %d",
		cv.CriticalCount, cv.HighCount, cv.MediumCount, cv.LowCount)))
	rows = append(rows, "")

	if len(cv.TopAssumptionsToVerify) > 0 {
		rows = append(rows, val.Render("Top Assumptions to Verify:"))
		for i, p := range cv.TopAssumptionsToVerify {
			if i >= 10 {
				break
			}
			priorityStyle := s.StatusWarn
			if p.Priority == "Critical" || p.Priority == "High" {
				priorityStyle = s.StatusBad
			} else if p.Priority == "Low" {
				priorityStyle = s.StatusGood
			}
			statusStyle := s.StatusWarn
			if p.Status == "Verified" {
				statusStyle = s.StatusGood
			} else if p.Status == "No Evidence" {
				statusStyle = s.StatusBad
			}
			rows = append(rows, fmt.Sprintf("  %s [%s] %s — Conf: %.0f%%",
				priorityStyle.Render(string(p.Priority)),
				statusStyle.Render(string(p.Status)),
				p.AssumptionText, p.Confidence))
		}
		rows = append(rows, "")
	}

	if len(cv.HighestRiskUnverified) > 0 {
		rows = append(rows, val.Render("Highest Risk Unverified:"))
		for _, p := range cv.HighestRiskUnverified {
			rows = append(rows, fmt.Sprintf("  %s [%s] %s — %s",
				s.StatusBad.Render("✗"), string(p.Priority), p.AssumptionText, p.Risk))
		}
		rows = append(rows, "")
	}

	return strings.Join(rows, "\n")
}

func renderVerificationRoadmap(s StyleSet, result *AnalysisResult) string {
	if result.VerificationOutput == nil || len(result.VerificationOutput.Roadmaps) == 0 {
		return "No verification roadmaps available."
	}
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	val := s.Value
	var rows []string

	rows = append(rows, val.Render("=== Verification Roadmaps ==="))
	rows = append(rows, "")

	for i, r := range result.VerificationOutput.Roadmaps {
		if i >= 10 {
			rows = append(rows, dim.Render(fmt.Sprintf("... and %d more roadmaps", len(result.VerificationOutput.Roadmaps)-10)))
			break
		}
		priorityStyle := s.StatusWarn
		if r.Priority == "Critical" || r.Priority == "High" {
			priorityStyle = s.StatusBad
		} else if r.Priority == "Low" {
			priorityStyle = s.StatusGood
		}
		rows = append(rows, fmt.Sprintf("%d. [%s] %s",
			i+1, priorityStyle.Render(string(r.Priority)), val.Render(r.AssumptionText)))
		rows = append(rows, dim.Render(fmt.Sprintf("   Effort: %s | Stakeholders: %s",
			string(r.Effort), strings.Join(r.Stakeholders, ", "))))
		for _, step := range r.Steps {
			rows = append(rows, fmt.Sprintf("   %d. %s", step.Step, step.Action))
		}
		rows = append(rows, "")
	}

	return strings.Join(rows, "\n")
}

func renderCISOVerificationSummary(s StyleSet, result *AnalysisResult) string {
	if result.VerificationOutput == nil || result.VerificationOutput.CISOView == nil {
		return "No CISO verification summary available."
	}
	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	val := s.Value
	var rows []string

	cv := result.VerificationOutput.CISOView
	rows = append(rows, val.Render("=== CISO Verification Summary ==="))
	rows = append(rows, "")

	rows = append(rows, dim.Render("Priority Distribution:"))
	rows = append(rows, fmt.Sprintf("  Critical: %d", cv.CriticalCount))
	rows = append(rows, fmt.Sprintf("  High: %d", cv.HighCount))
	rows = append(rows, fmt.Sprintf("  Medium: %d", cv.MediumCount))
	rows = append(rows, fmt.Sprintf("  Low: %d", cv.LowCount))
	rows = append(rows, "")

	rows = append(rows, dim.Render(fmt.Sprintf("Verification Backlog: %d items", len(cv.VerificationBacklog))))
	rows = append(rows, dim.Render(fmt.Sprintf("Evidence Gaps: %d", len(cv.EvidenceGaps))))
	rows = append(rows, "")

	if len(cv.EvidenceGaps) > 0 {
		rows = append(rows, val.Render("Evidence Gaps:"))
		seen := make(map[string]bool)
		for _, gap := range cv.EvidenceGaps {
			if !seen[gap] {
				seen[gap] = true
				rows = append(rows, fmt.Sprintf("  %s %s", s.StatusBad.Render("✗"), gap))
			}
		}
		rows = append(rows, "")
	}

	if len(cv.HighestRiskUnverified) > 0 {
		rows = append(rows, val.Render("Highest Risk Unverified:"))
		for _, p := range cv.HighestRiskUnverified {
			rows = append(rows, fmt.Sprintf("  %s [%s] %s",
				s.StatusBad.Render("⚠"), string(p.Priority), p.AssumptionText))
		}
		rows = append(rows, "")
	}

	if len(cv.VerificationBacklog) > 0 {
		rows = append(rows, dim.Render(fmt.Sprintf("Backlog: %d items — review verification plans for details", len(cv.VerificationBacklog))))
	}

	return strings.Join(rows, "\n")
}
