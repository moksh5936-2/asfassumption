package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type resultsModel struct {
	result    *AnalysisResult
	resultTab int
	tabs      []resultTabDef
	tabScroll map[int]int
	vpReady   bool
	vpWidth   int
	vpHeight  int
}

type resultTabDef struct {
	name string
}

func newResultsModel() resultsModel {
	return resultsModel{
		tabs: []resultTabDef{
			{name: "Summary"},
			{name: "Assumptions"},
			{name: "Verification"},
			{name: "Contradictions"},
			{name: "Trust Chains"},
			{name: "Impact"},
			{name: "Blind Spots"},
			{name: "Controls"},
			{name: "Reports"},
			{name: "SDRI"},
			{name: "Security Design Review"},
			{name: "SPOFs"},
		},
		tabScroll: make(map[int]int),
	}
}

func (m resultsModel) Update(msg tea.Msg) (resultsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.resultTab = (m.resultTab + 1) % len(m.tabs)
		case "shift+tab":
			m.resultTab = (m.resultTab - 1 + len(m.tabs)) % len(m.tabs)
		}
	}
	return m, nil
}

func (m mainModel) viewResults() string {
	s := m.styles

	if m.results.result == nil {
		return lipgloss.JoinVertical(lipgloss.Left,
			s.Title.Render("Analysis Results"),
			s.EmptyState.Render("No results available. Run an analysis first."),
		)
	}

	r := m.results.result
	tabBar := m.renderResultTabs()

	query := m.searchQuery
	var content string
	switch m.results.resultTab {
	case 0:
		content = renderResultSummary(s, r)
	case 1:
		content = renderResultAssumptions(s, r, query)
	case 2:
		content = renderResultVerification(s, r)
	case 3:
		content = renderResultContradictions(s, r, query)
	case 4:
		content = renderResultTrust(s, r, query)
	case 5:
		content = renderResultImpact(s, r)
	case 6:
		content = renderResultBlindSpots(s, r)
	case 7:
		content = renderResultControls(s, r, query)
	case 8:
		content = renderResultReports(s, r)
	case 9:
		content = renderResultSDRI(s, r)
	case 10:
		content = renderResultSecurityDesignReview(s, r)
	case 11:
		content = renderResultSPOFs(s, r)
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		s.Title.Render("Analysis Results"),
		tabBar,
		content,
	)
}

func (m mainModel) renderResultTabs() string {
	s := m.styles
	var tabs []string
	for i, tab := range m.results.tabs {
		style := s.Tab
		if i == m.results.resultTab {
			style = s.TabActive
		}
		count := resultTabCount(m.results.result, i)
		label := tab.name
		if count > 0 {
			label = fmt.Sprintf("%s %d", label, count)
		}
		tabs = append(tabs, style.Render(" "+label+" "))
	}
	bar := lipgloss.JoinHorizontal(lipgloss.Top, tabs...)
	sep := lipgloss.NewStyle().
		Foreground(s.Theme().Border).
		Render(strings.Repeat("─", max(1, lipgloss.Width(bar))))
	return lipgloss.JoinVertical(lipgloss.Left, bar, sep)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func resultTabCount(r *AnalysisResult, tab int) int {
	switch tab {
	case 1:
		return len(r.Assumptions)
	case 2:
		if r.VerificationOutput != nil && r.VerificationOutput.Assessment != nil {
			a := r.VerificationOutput.Assessment
			return a.VerifiedCount + a.PartialCount + a.UnverifiedCount + a.NoEvidenceCount
		}
	case 3:
		return len(r.Contradictions)
	case 4:
		if r.TrustOutput != nil {
			return len(r.TrustOutput.TrustChains)
		}
	case 5:
		c := 0
		if r.TrustOutput != nil {
			c += len(r.TrustOutput.SinglePointsOfTrust)
		}
		if r.ReviewOutput != nil && r.ReviewOutput.Queue != nil {
			c += len(r.ReviewOutput.Queue.Items)
		}
		return c
	case 6:
		if r.CoverageOutput != nil && r.CoverageOutput.BlindSpots != nil {
			return len(r.CoverageOutput.BlindSpots)
		}
	case 7:
		return len(r.Controls)
	case 8:
		c := 0
		if r.NarrativeOutput != nil {
			c++
		}
		if r.ReviewOutput != nil {
			c++
		}
		return c
	case 9:
		c := 0
		if r.SDRISummary != "" {
			c++
		}
		c += len(r.SDRIControls)
		return c
	case 10:
		c := len(r.SDRIDesignFindings) + len(r.SDRIAchitecturalWeaknesses) + len(r.SDRIRemediations)
		if c == 0 {
			c = len(r.SDRIComplianceAlignments)
		}
		return c
	case 11:
		if r.TrustOutput != nil {
			return len(r.TrustOutput.SinglePointsOfTrust)
		}
	}
	return 0
}

func (m mainModel) updateResults(msg tea.Msg) (tea.Model, tea.Cmd) {
	oldTab := m.results.resultTab
	var cmd tea.Cmd
	m.results, cmd = m.results.Update(msg)
	if m.results.resultTab != oldTab {
		m.results.tabScroll[oldTab] = m.vp.YOffset
		if y, ok := m.results.tabScroll[m.results.resultTab]; ok {
			m.vp.YOffset = y
		} else {
			m.vp.YOffset = 0
		}
	}
	return m, cmd
}

func renderResultSummary(s StyleSet, r *AnalysisResult) string {
	var rows []string

	rows = append(rows, s.Section.Render("Summary"))
	rows = append(rows, fmt.Sprintf("  Architecture: %s", r.ArchitectureName))
	rows = append(rows, fmt.Sprintf("  Domain:       %s", r.Domain))
	rows = append(rows, fmt.Sprintf("  Mode:         %s", r.AnalysisMode))
	rows = append(rows, "")

	totalAssumptions := len(r.Assumptions)
	critical := countRisk(r.Assumptions, "Critical")
	high := countRisk(r.Assumptions, "High")
	medium := countRisk(r.Assumptions, "Medium")
	low := countRisk(r.Assumptions, "Low")

	rows = append(rows, s.Section.Render("Assumptions"))
	rows = append(rows, fmt.Sprintf("  Total:     %d", totalAssumptions))
	if critical > 0 {
		rows = append(rows, fmt.Sprintf("  %s %d", s.BadgeCritical.Render("CRITICAL"), critical))
	}
	if high > 0 {
		rows = append(rows, fmt.Sprintf("  %s %d", s.BadgeHigh.Render("HIGH"), high))
	}
	if medium > 0 {
		rows = append(rows, fmt.Sprintf("  %s %d", s.BadgeMedium.Render("MEDIUM"), medium))
	}
	if low > 0 {
		rows = append(rows, fmt.Sprintf("  %s %d", s.BadgeLow.Render("LOW"), low))
	}
	rows = append(rows, "")

	if r.VerificationOutput != nil && r.VerificationOutput.Assessment != nil {
		a := r.VerificationOutput.Assessment
		rows = append(rows, s.Section.Render("Verification"))
		rows = append(rows, fmt.Sprintf("  Verified:  %d", a.VerifiedCount))
		rows = append(rows, fmt.Sprintf("  Partial:   %d", a.PartialCount))
		rows = append(rows, fmt.Sprintf("  Unknown:   %d", a.UnverifiedCount))
		rows = append(rows, "")
	}

	if len(r.Contradictions) > 0 {
		rows = append(rows, s.Section.Render("Contradictions"))
		rows = append(rows, fmt.Sprintf("  Found:     %d", len(r.Contradictions)))
		rows = append(rows, "")
	}

	if r.TrustOutput != nil && len(r.TrustOutput.TrustChains) > 0 {
		rows = append(rows, s.Section.Render("Trust Chains"))
		rows = append(rows, fmt.Sprintf("  Chains:    %d", len(r.TrustOutput.TrustChains)))
		rows = append(rows, fmt.Sprintf("  SPOFs:     %d", len(r.TrustOutput.SinglePointsOfTrust)))
		rows = append(rows, "")
	}

	if r.CoverageOutput != nil {
		rows = append(rows, s.Section.Render("Coverage"))
		if r.CoverageOutput.BlindSpots != nil {
			rows = append(rows, fmt.Sprintf("  Blind Spots: %d", len(r.CoverageOutput.BlindSpots)))
		}
		if r.CoverageOutput.Assessment != nil && r.CoverageOutput.Assessment.Gaps != nil {
			rows = append(rows, fmt.Sprintf("  Gaps:     %d", len(r.CoverageOutput.Assessment.Gaps)))
		}
		rows = append(rows, "")
	}

	rows = append(rows, s.Section.Render("Exports"))
	rows = append(rows, "  Press 'e' to export results.")
	return strings.Join(rows, "\n")
}

func renderResultAssumptions(s StyleSet, r *AnalysisResult, searchQuery string) string {
	if len(r.Assumptions) == 0 {
		return s.EmptyState.Render("No assumptions found.")
	}

	var filtered []Assumption
	for _, a := range r.Assumptions {
		if searchQuery != "" {
			q := strings.ToLower(searchQuery)
			match := strings.Contains(strings.ToLower(a.ID), q) ||
				strings.Contains(strings.ToLower(a.Description), q) ||
				strings.Contains(strings.ToLower(string(a.Risk)), q)
			if !match {
				continue
			}
		}
		filtered = append(filtered, a)
	}

	sectionTitle := fmt.Sprintf("Assumptions (%d)", len(r.Assumptions))
	if searchQuery != "" {
		sectionTitle = fmt.Sprintf("Assumptions (%d of %d matching \"%s\")", len(filtered), len(r.Assumptions), searchQuery)
	}
	if len(filtered) == 0 {
		return s.EmptyState.Render(fmt.Sprintf("No assumptions match \"%s\".", searchQuery))
	}

	var rows []string
	rows = append(rows, s.Section.Render(sectionTitle))
	rows = append(rows, s.DimText.Render("  ID  Risk       Confidence  Description"))

	for _, a := range filtered {
		riskStyle := riskStyle(s, a.Risk)
		confPct := int(a.Confidence * 100)
		confStyle := confidenceStyle(s, confPct)
		text := a.Description
		rows = append(rows, fmt.Sprintf("  %s %s  %s  %s",
			a.ID,
			riskStyle.Render(padRight(string(a.Risk), 10)),
			confStyle.Render(fmt.Sprintf("%3d%%", confPct)),
			text))
	}
	return strings.Join(rows, "\n")
}

func renderResultVerification(s StyleSet, r *AnalysisResult) string {
	if r.VerificationOutput == nil || r.VerificationOutput.Assessment == nil {
		return s.EmptyState.Render("No verification data available.")
	}
	a := r.VerificationOutput.Assessment
	var rows []string
	rows = append(rows, s.Section.Render(fmt.Sprintf("Verification — %d total assumptions",
		a.TotalAssumptions)))

	rows = append(rows, "")
	rows = append(rows, fmt.Sprintf("  %s  %d verified", s.StatusGood.Render("✓"), a.VerifiedCount))
	rows = append(rows, fmt.Sprintf("  %s  %d partial", s.StatusWarn.Render("~"), a.PartialCount))
	rows = append(rows, fmt.Sprintf("  %s  %d unverified", s.DimText.Render("?"), a.UnverifiedCount))
	rows = append(rows, fmt.Sprintf("  %s  %d no evidence", s.DimText.Render("○"), a.NoEvidenceCount))
	rows = append(rows, "")

	if a.OverallConfidence > 0 {
		rows = append(rows, fmt.Sprintf("  Overall Confidence: %.1f%%", a.OverallConfidence*100))
		rows = append(rows, "")
	}

	if r.VerificationOutput.CISOView != nil {
		cv := r.VerificationOutput.CISOView
		if len(cv.TopAssumptionsToVerify) > 0 {
			rows = append(rows, s.Section.Render("Top Items to Verify"))
			for _, plan := range cv.TopAssumptionsToVerify {
				text := plan.AssumptionText
				rows = append(rows, fmt.Sprintf("  • %s", text))
			}
			rows = append(rows, "")
		}
		if len(cv.EvidenceGaps) > 0 {
			rows = append(rows, s.StatusWarn.Render("Evidence Gaps:"))
			for _, g := range cv.EvidenceGaps {
				rows = append(rows, fmt.Sprintf("  • %s", g))
			}
			rows = append(rows, "")
		}
	}

	return strings.Join(rows, "\n")
}

func renderResultContradictions(s StyleSet, r *AnalysisResult, searchQuery string) string {
	if len(r.Contradictions) == 0 {
		return s.EmptyState.Render("No contradictions detected.")
	}

	var filtered []Contradiction
	for _, c := range r.Contradictions {
		if searchQuery != "" {
			q := strings.ToLower(searchQuery)
			match := strings.Contains(strings.ToLower(c.Description), q) ||
				strings.Contains(strings.ToLower(c.Explanation), q)
			if !match {
				continue
			}
		}
		filtered = append(filtered, c)
	}

	sectionTitle := fmt.Sprintf("Contradictions (%d)", len(r.Contradictions))
	if searchQuery != "" {
		sectionTitle = fmt.Sprintf("Contradictions (%d of %d matching \"%s\")", len(filtered), len(r.Contradictions), searchQuery)
	}
	if len(filtered) == 0 {
		return s.EmptyState.Render(fmt.Sprintf("No contradictions match \"%s\".", searchQuery))
	}

	var rows []string
	rows = append(rows, s.Section.Render(sectionTitle))

	for _, c := range filtered {
		severityStyle := s.StatusWarn
		if c.Severity == RiskHigh || c.Severity == RiskCritical {
			severityStyle = s.StatusBad
		} else if c.Severity == RiskLow {
			severityStyle = s.StatusGood
		}
		rows = append(rows, fmt.Sprintf("  %s %s", severityStyle.Render(strings.ToUpper(string(c.Severity))), c.Description))
		if c.Explanation != "" {
			rows = append(rows, fmt.Sprintf("    Reason: %s", c.Explanation))
		}
		if len(c.AffectedAssumptions) > 0 {
			rows = append(rows, fmt.Sprintf("    Affects: %s", strings.Join(c.AffectedAssumptions, ", ")))
		}
		rows = append(rows, "")
	}
	return strings.Join(rows, "\n")
}

func renderResultTrust(s StyleSet, r *AnalysisResult, searchQuery string) string {
	if r.TrustOutput == nil {
		return s.EmptyState.Render("No trust chain data available.")
	}
	var rows []string
	q := strings.ToLower(searchQuery)

	if len(r.TrustOutput.TrustChains) > 0 {
		total := len(r.TrustOutput.TrustChains)
		matchCount := 0
		var chainRows []string

		for _, chain := range r.TrustOutput.TrustChains {
			if searchQuery != "" {
				match := strings.Contains(strings.ToLower(chain.ID), q) ||
					strings.Contains(strings.ToLower(chain.RootNode), q) ||
					strings.Contains(strings.ToLower(chain.LeafNode), q)
				if !match {
					continue
				}
			}
			matchCount++
			confidence := int(chain.Confidence * 100)
			route := chain.RootNode
			if chain.LeafNode != "" && chain.LeafNode != chain.RootNode {
				route += " → " + chain.LeafNode
			}
			chainRows = append(chainRows, fmt.Sprintf("  Chain %s: %s (len=%d, conf=%d%%)",
				chain.ID, route, chain.Length, confidence))
			if len(chain.Nodes) > 0 {
				path := strings.Join(chain.Nodes, " → ")
				chainRows = append(chainRows, fmt.Sprintf("    Path: %s", path))
			}
			chainRows = append(chainRows, "")
		}

		sectionTitle := fmt.Sprintf("Trust Chains (%d)", total)
		if searchQuery != "" {
			sectionTitle = fmt.Sprintf("Trust Chains (%d of %d matching \"%s\")", matchCount, total, searchQuery)
		}
		if matchCount == 0 {
			chainRows = append(chainRows, s.EmptyState.Render(fmt.Sprintf("No trust chains match \"%s\".", searchQuery)))
		}
		rows = append(rows, s.Section.Render(sectionTitle))
		rows = append(rows, chainRows...)
	}

	if len(rows) == 0 {
		return s.EmptyState.Render("No trust chain data found.")
	}
	return strings.Join(rows, "\n")
}

func renderResultSPOFs(s StyleSet, r *AnalysisResult) string {
	if r.TrustOutput == nil || len(r.TrustOutput.SinglePointsOfTrust) == 0 {
		return s.EmptyState.Render("No single points of trust failure identified.")
	}
	var rows []string
	rows = append(rows, s.Section.Render(fmt.Sprintf("Single Points of Trust Failure (%d)", len(r.TrustOutput.SinglePointsOfTrust))))
	for _, spof := range r.TrustOutput.SinglePointsOfTrust {
		rows = append(rows, fmt.Sprintf("  ⚠ %s", spof.AssumptionText))
	}
	return strings.Join(rows, "\n")
}

func renderResultImpact(s StyleSet, r *AnalysisResult) string {
	var rows []string

	if r.ReviewOutput != nil && r.ReviewOutput.Queue != nil && len(r.ReviewOutput.Queue.Items) > 0 {
		rows = append(rows, s.Section.Render(fmt.Sprintf("Priority Queue (%d)", len(r.ReviewOutput.Queue.Items))))
		for _, item := range r.ReviewOutput.Queue.Items {
			riskStyle := riskStyle(s, RiskLevel(item.Risk))
			text := item.AssumptionText
			rows = append(rows, fmt.Sprintf("  #%d %s [%.0f] %s",
				item.Rank, riskStyle.Render(item.Risk), item.PriorityScore, text))
		}
		rows = append(rows, "")
	}

	if r.TrustOutput != nil && len(r.TrustOutput.SinglePointsOfTrust) > 0 {
		rows = append(rows, s.Section.Render("Single Points of Trust Failure"))
		for _, spof := range r.TrustOutput.SinglePointsOfTrust {
			rows = append(rows, fmt.Sprintf("  ⚠ %s", spof.AssumptionText))
		}
		rows = append(rows, "")
	}

	if r.ReviewOutput != nil && r.ReviewOutput.CISODashboard != nil {
		d := r.ReviewOutput.CISODashboard
		rows = append(rows, s.Section.Render("CISO View"))
		rows = append(rows, fmt.Sprintf("  Critical: %d", d.CriticalAssumptions))
		rows = append(rows, fmt.Sprintf("  High:     %d", d.HighAssumptions))
		rows = append(rows, "")
	}

	if len(rows) == 0 {
		return s.EmptyState.Render("No impact data available.")
	}
	return strings.Join(rows, "\n")
}

func renderResultBlindSpots(s StyleSet, r *AnalysisResult) string {
	if r.CoverageOutput == nil {
		return s.EmptyState.Render("No blind spot data available.")
	}
	var rows []string

	if len(r.CoverageOutput.BlindSpots) > 0 {
		rows = append(rows, s.Section.Render(fmt.Sprintf("Blind Spots (%d)", len(r.CoverageOutput.BlindSpots))))
		for _, bs := range r.CoverageOutput.BlindSpots {
			rows = append(rows, fmt.Sprintf("  • %s", bs.Title))
		}
		rows = append(rows, "")
	}

	if len(r.CoverageOutput.DomainBlindSpots) > 0 {
		rows = append(rows, s.Section.Render(fmt.Sprintf("Domain Blind Spots (%d)", len(r.CoverageOutput.DomainBlindSpots))))
		for _, dbs := range r.CoverageOutput.DomainBlindSpots {
			rows = append(rows, fmt.Sprintf("  • %s", dbs.MissingArea))
		}
		rows = append(rows, "")
	}

	if r.CoverageOutput.Assessment != nil && len(r.CoverageOutput.Assessment.Gaps) > 0 {
		rows = append(rows, s.Section.Render(fmt.Sprintf("Coverage Gaps (%d)", len(r.CoverageOutput.Assessment.Gaps))))
		for _, gap := range r.CoverageOutput.Assessment.Gaps {
			rows = append(rows, fmt.Sprintf("  %s: %s (observed: %d, expected: %d)",
				gap.Category, gap.Risk, gap.ObservedCount, gap.ExpectedCount))
		}
	}

	if len(rows) == 0 {
		rows = append(rows, s.StatusGood.Render("No blind spots or gaps detected."))
	}
	return strings.Join(rows, "\n")
}

func renderResultControls(s StyleSet, r *AnalysisResult, searchQuery string) string {
	if len(r.Controls) == 0 {
		return s.EmptyState.Render("No recommended controls.")
	}

	q := strings.ToLower(searchQuery)
	matchCount := 0
	var controlRows []string

	for _, c := range r.Controls {
		if searchQuery != "" {
			match := strings.Contains(strings.ToLower(c.Description), q) ||
				strings.Contains(strings.ToLower(c.Rationale), q)
			if !match {
				continue
			}
		}
		matchCount++
		controlRows = append(controlRows, fmt.Sprintf("  • %s", c.Description))
		if c.Rationale != "" {
			controlRows = append(controlRows, fmt.Sprintf("    %s", c.Rationale))
		}
		if len(c.MitigatedAssumptionIDs) > 0 {
			controlRows = append(controlRows, fmt.Sprintf("    Mitigates: %s", strings.Join(c.MitigatedAssumptionIDs, ", ")))
		}
		controlRows = append(controlRows, "")
	}

	sectionTitle := fmt.Sprintf("Recommended Controls (%d)", len(r.Controls))
	if searchQuery != "" {
		sectionTitle = fmt.Sprintf("Recommended Controls (%d of %d matching \"%s\")", matchCount, len(r.Controls), searchQuery)
	}
	if matchCount == 0 {
		return s.EmptyState.Render(fmt.Sprintf("No controls match \"%s\".", searchQuery))
	}

	var rows []string
	rows = append(rows, s.Section.Render(sectionTitle))
	rows = append(rows, controlRows...)
	return strings.Join(rows, "\n")
}

func renderResultReports(s StyleSet, r *AnalysisResult) string {
	var rows []string
	rows = append(rows, s.Section.Render("Reports & Exports"))
	rows = append(rows, "")

	if r.NarrativeOutput != nil {
		rows = append(rows, s.StatusGood.Render("Architect Narrative Available"))
		rows = append(rows, s.DimText.Render("  Press 'e' → select narrative-md or narrative-html"))
		rows = append(rows, "")
	}

	if r.ReviewOutput != nil && r.ReviewOutput.Campaigns != nil && len(r.ReviewOutput.Campaigns) > 0 {
		rows = append(rows, s.Section.Render(fmt.Sprintf("Review Campaigns (%d)", len(r.ReviewOutput.Campaigns))))
		for _, c := range r.ReviewOutput.Campaigns {
			rows = append(rows, fmt.Sprintf("  ▶ %s (%d items)", c.Name, len(c.Items)))
		}
		rows = append(rows, "")
	}

	if r.ConfidenceOutput != nil && len(r.ConfidenceOutput.Breakdowns) > 0 {
		rows = append(rows, s.Section.Render("Confidence"))
		highConf := 0
		lowConf := 0
		for _, bd := range r.ConfidenceOutput.Breakdowns {
			if bd.FinalConfidence >= 70 {
				highConf++
			} else if bd.FinalConfidence < 40 {
				lowConf++
			}
		}
		rows = append(rows, fmt.Sprintf("  High Confidence: %d", highConf))
		rows = append(rows, fmt.Sprintf("  Low Confidence:  %d", lowConf))
		rows = append(rows, "")
	}

	rows = append(rows, s.Section.Render("Export"))
	rows = append(rows, "  Press 'e' to open export dialog.")
	rows = append(rows, "  Available formats: JSON, Markdown, HTML, CSV, PDF")

	return strings.Join(rows, "\n")
}

func renderResultSDRI(s StyleSet, r *AnalysisResult) string {
	if r.SDRISummary == "" && len(r.SDRIControls) == 0 {
		return s.EmptyState.Render("No SDRI data available.")
	}
	var rows []string
	rows = append(rows, s.Section.Render("Security Design Review Intelligence (SDRI)"))
	rows = append(rows, "")

	if r.SDRISummary != "" {
		rows = append(rows, s.Subtitle.Render("Executive Summary"))
		rows = append(rows, "  "+r.SDRISummary)
		rows = append(rows, "")
	}

	if len(r.SDRIControls) > 0 {
		rows = append(rows, s.Subtitle.Render(fmt.Sprintf("Control Inventory (%d)", len(r.SDRIControls))))
		byStatus := map[string]int{}
		for _, c := range r.SDRIControls {
			status := c.Status
			if status == "" {
				status = "unknown"
			}
			byStatus[status]++
		}
		for st, n := range byStatus {
			style := s.Value
			switch st {
			case "implemented", "partial":
				style = s.StatusGood
			case "planned", "in-progress":
				style = s.StatusWarn
			case "missing", "none":
				style = s.StatusBad
			}
			rows = append(rows, fmt.Sprintf("  %s %s: %d", style.Render("●"), st, n))
		}
		rows = append(rows, "")
	}

	if len(r.SDRICoverageByCategory) > 0 {
		rows = append(rows, s.Subtitle.Render("Coverage by Category"))
		for _, c := range r.SDRICoverageByCategory {
			style := s.StatusGood
			if c.Coverage < 50 {
				style = s.StatusBad
			} else if c.Coverage < 80 {
				style = s.StatusWarn
			}
			rows = append(rows, fmt.Sprintf("  %s %s: %.0f%% (%d/%d)",
				style.Render("●"), c.Category, c.Coverage*100, c.Observed, c.Expected))
		}
		rows = append(rows, "")
	}

	if r.SDRICoverageDashboard != nil && len(r.SDRICoverageDashboard) > 0 {
		rows = append(rows, s.Subtitle.Render("Coverage Dashboard"))
		for k, v := range r.SDRICoverageDashboard {
			pct := int(v * 100)
			style := confidenceStyle(s, pct)
			rows = append(rows, fmt.Sprintf("  %s: %s", k, style.Render(fmt.Sprintf("%d%%", pct))))
		}
		rows = append(rows, "")
	}

	return strings.Join(rows, "\n")
}

func renderResultSecurityDesignReview(s StyleSet, r *AnalysisResult) string {
	if len(r.SDRIDesignFindings) == 0 && len(r.SDRIAchitecturalWeaknesses) == 0 &&
		len(r.SDRIRemediations) == 0 && len(r.SDRIComplianceAlignments) == 0 {
		return s.EmptyState.Render("No security design review data available.")
	}
	var rows []string
	rows = append(rows, s.Section.Render("Security Design Review"))
	rows = append(rows, "")

	if len(r.SDRIDesignFindings) > 0 {
		rows = append(rows, s.Subtitle.Render(fmt.Sprintf("Design Findings (%d)", len(r.SDRIDesignFindings))))
		for _, f := range r.SDRIDesignFindings {
			severityStyle := s.StatusWarn
			switch f.Severity {
			case "Critical", "High":
				severityStyle = s.StatusBad
			case "Low":
				severityStyle = s.StatusGood
			}
			rows = append(rows, fmt.Sprintf("  %s [%s] %s", severityStyle.Render("●"), f.Severity, f.Title))
			if f.Description != "" {
				rows = append(rows, fmt.Sprintf("    %s", f.Description))
			}
			if f.Recommendation != "" {
				rows = append(rows, fmt.Sprintf("    → %s", f.Recommendation))
			}
			rows = append(rows, "")
		}
	}

	if len(r.SDRIAchitecturalWeaknesses) > 0 {
		rows = append(rows, s.Subtitle.Render(fmt.Sprintf("Architectural Weaknesses (%d)", len(r.SDRIAchitecturalWeaknesses))))
		for _, w := range r.SDRIAchitecturalWeaknesses {
			severityStyle := s.StatusWarn
			switch w.Severity {
			case "Critical", "High":
				severityStyle = s.StatusBad
			case "Low":
				severityStyle = s.StatusGood
			}
			rows = append(rows, fmt.Sprintf("  %s [%s] %s", severityStyle.Render("●"), w.Severity, w.Pattern))
			if w.Description != "" {
				rows = append(rows, fmt.Sprintf("    %s", w.Description))
			}
			if w.Recommendation != "" {
				rows = append(rows, fmt.Sprintf("    → %s", w.Recommendation))
			}
			rows = append(rows, "")
		}
	}

	if len(r.SDRIRemediations) > 0 {
		rows = append(rows, s.Subtitle.Render(fmt.Sprintf("Remediations (%d)", len(r.SDRIRemediations))))
		for _, rem := range r.SDRIRemediations {
			rows = append(rows, fmt.Sprintf("  #%d [%.0f] %s", rem.Priority, rem.RiskScore, rem.Description))
			if rem.Recommendation != "" {
				rows = append(rows, fmt.Sprintf("    → %s (effort: %s)", rem.Recommendation, rem.Effort))
			}
			rows = append(rows, "")
		}
	}

	if len(r.SDRIComplianceAlignments) > 0 {
		rows = append(rows, s.Subtitle.Render(fmt.Sprintf("Compliance Alignments (%d)", len(r.SDRIComplianceAlignments))))
		for _, m := range r.SDRIComplianceAlignments {
			style := s.StatusGood
			if m.Coverage < 50 {
				style = s.StatusBad
			} else if m.Coverage < 80 {
				style = s.StatusWarn
			}
			rows = append(rows, fmt.Sprintf("  %s %s: %.0f%% (%s)",
				style.Render("●"), m.Framework, m.Coverage, m.Status))
		}
		rows = append(rows, "")
	}

	return strings.Join(rows, "\n")
}

func countRisk(assumptions []Assumption, risk RiskLevel) int {
	n := 0
	for _, a := range assumptions {
		if a.Risk == risk {
			n++
		}
	}
	return n
}

func riskStyle(s StyleSet, risk RiskLevel) lipgloss.Style {
	switch risk {
	case "Critical":
		return s.StatusBad
	case "High":
		return s.StatusWarn
	case "Medium":
		return s.StatusGood
	default:
		return s.Value
	}
}

func confidenceStyle(s StyleSet, pct int) lipgloss.Style {
	switch {
	case pct >= 80:
		return s.StatusGood
	case pct >= 50:
		return s.StatusWarn
	default:
		return s.StatusBad
	}
}

func padRight(s string, n int) string {
	if len(s) >= n {
		return s
	}
	return s + strings.Repeat(" ", n-len(s))
}

func exportFormatFromConfig(cfg *Config) ExportFormat {
	switch cfg.Output.Default {
	case "json":
		return ExportJSON
	case "markdown":
		return ExportMarkdown
	case "html":
		return ExportHTML
	case "csv":
		return ExportCSV
	case "pdf":
		return ExportPDF
	case "narrative-md":
		return ExportNarrativeMarkdown
	case "narrative-html":
		return ExportNarrativeHTML
	default:
		return ExportMarkdown
	}
}
