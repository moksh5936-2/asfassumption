package main

import (
	"fmt"
	"strings"

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
			"Compliance",
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
		return renderCompliance(s, result.Compliance)
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
