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
			{name: "Overview"},
			{name: "Assumptions"},
			{name: "Verification"},
			{name: "Contradictions"},
			{name: "Trust"},
			{name: "Controls"},
			{name: "SDRI"},
		},
		tabScroll: make(map[int]int),
	}
}

func (m resultsModel) Update(msg tea.Msg) (resultsModel, tea.Cmd) {
	return m, nil
}

func (m mainModel) viewResults() string {
	s := m.styles

	if m.results.result == nil {
		return s.Card("Case Workspace",
			s.EmptyState.Render("No results available. Run an analysis first."),
			m.mainWidth())
	}

	r := m.results.result
	query := m.searchQuery
	var content string
	tab := m.results.resultTab
	switch tab {
	case 0:
		content = renderResultSummary(s, r, m.mainWidth()-4)
	case 1:
		content = renderResultAssumptions(s, r, query, m.mainWidth()-4)
	case 2:
		content = renderResultVerification(s, r, m.mainWidth()-4)
	case 3:
		content = renderResultContradictions(s, r, query, m.mainWidth()-4)
	case 4:
		content = renderResultTrust(s, r, query, m.mainWidth()-4)
	case 5:
		content = renderResultControls(s, r, query, m.mainWidth()-4)
	case 6:
		content = renderResultSDRI(s, r, m.mainWidth()-4)
	}

	tabBar := m.renderResultTabs()
	sectionName := m.results.tabs[tab].name
	header := s.PremiumHeader(sectionName, m.mainWidth())
	return lipgloss.JoinVertical(lipgloss.Left,
		header,
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
	sep := s.SectionRule.Render(strings.Repeat("─", max(1, m.mainWidth())))
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
		c := 0
		if r.TrustOutput != nil {
			c += len(r.TrustOutput.TrustChains) + len(r.TrustOutput.SinglePointsOfTrust)
		}
		if r.ReviewOutput != nil && r.ReviewOutput.Queue != nil {
			c += len(r.ReviewOutput.Queue.Items)
		}
		return c
	case 5:
		return len(r.Controls)
	case 6:
		c := 0
		if r.SDRISummary != "" {
			c++
		}
		c += len(r.SDRIControls)
		c += len(r.SDRIDesignFindings) + len(r.SDRIAchitecturalWeaknesses) + len(r.SDRIRemediations) + len(r.SDRIComplianceAlignments)
		return c
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

func renderResultSummary(s StyleSet, r *AnalysisResult, width int) string {
	if width < 30 {
		width = 30
	}

	critical := countRisk(r.Assumptions, "Critical")
	high := countRisk(r.Assumptions, "High")
	medium := countRisk(r.Assumptions, "Medium")
	low := countRisk(r.Assumptions, "Low")

	var infoRows []string
	infoRows = append(infoRows, fmt.Sprintf("  %s  %s", s.SectionTitle.Render("Architecture"), s.Value.Render(r.ArchitectureName)))
	infoRows = append(infoRows, fmt.Sprintf("  %s  %s", s.SectionTitle.Render("Domain"), s.Value.Render(r.Domain)))
	infoRows = append(infoRows, fmt.Sprintf("  %s  %s", s.SectionTitle.Render("Mode"), s.Value.Render(r.AnalysisMode)))
	infoRows = append(infoRows, "")
	infoRows = append(infoRows, fmt.Sprintf("  %s  %s", s.SectionTitle.Render("Assumptions"), s.Value.Render(fmt.Sprintf("%d", len(r.Assumptions)))))
	infoCard := s.Card("Case Info", strings.Join(infoRows, "\n"), width)

	riskCard := s.Card("Risk Distribution", s.RiskHeatmap("", critical, high, medium, low), width)

	var verifyCard string
	if r.VerificationOutput != nil && r.VerificationOutput.Assessment != nil {
		a := r.VerificationOutput.Assessment
		vContent := fmt.Sprintf("  %s  %s  %s",
			s.StatusGood.Render(fmt.Sprintf("✓ %d verified", a.VerifiedCount)),
			s.StatusWarn.Render(fmt.Sprintf("~ %d partial", a.PartialCount)),
			s.DimText.Render(fmt.Sprintf("? %d unverified", a.UnverifiedCount)))
		verifyCard = s.Card("Verification", vContent, width)
	}

	var contradictCard string
	if len(r.Contradictions) > 0 {
		criticalCt := 0
		highCt := 0
		for _, c := range r.Contradictions {
			switch c.Severity {
			case RiskCritical:
				criticalCt++
			case RiskHigh:
				highCt++
			}
		}
		cContent := fmt.Sprintf("  %s Total: %d  %s %d  %s %d",
			s.StatusBad.Render("Critical:"), criticalCt,
			s.StatusWarn.Render("High:"), highCt,
			s.DimText.Render("Total:"), len(r.Contradictions))
		contradictCard = s.CardAccent("Contradictions", cContent, width)
	}

	var trustCard string
	if r.TrustOutput != nil && len(r.TrustOutput.TrustChains) > 0 {
		tContent := fmt.Sprintf("  %s  %s",
			s.Value.Render(fmt.Sprintf("Chains: %d", len(r.TrustOutput.TrustChains))),
			s.RiskBadge(fmt.Sprintf("SPOFs: %d", len(r.TrustOutput.SinglePointsOfTrust))),
		)
		trustCard = s.Card("Trust Analysis", tContent, width)
	}

	var coverageCard string
	if r.CoverageOutput != nil {
		var cvRows []string
		if r.CoverageOutput.BlindSpots != nil {
			cvRows = append(cvRows, fmt.Sprintf("  Blind Spots: %d", len(r.CoverageOutput.BlindSpots)))
		}
		if r.CoverageOutput.Assessment != nil && r.CoverageOutput.Assessment.Gaps != nil {
			cvRows = append(cvRows, fmt.Sprintf("  Gaps: %d", len(r.CoverageOutput.Assessment.Gaps)))
		}
		if len(cvRows) > 0 {
			coverageCard = s.Card("Coverage", strings.Join(cvRows, "\n"), width)
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		infoCard,
		riskCard,
		verifyCard,
		contradictCard,
		trustCard,
		coverageCard,
	)
}

func renderResultAssumptions(s StyleSet, r *AnalysisResult, searchQuery string, width int) string {
	if len(r.Assumptions) == 0 {
		return s.Card("", s.EmptyState.Render("No assumptions found."), width)
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
		return s.Card("", s.EmptyState.Render(fmt.Sprintf("No assumptions match \"%s\".", searchQuery)), width)
	}

	var rows []string
	rows = append(rows, s.SectionRule.Render(strings.Repeat("─", max(1, width-4))))
	rows = append(rows, s.SubSectionTitle.Render("ID  Risk        Confidence  Description"))

	for _, a := range filtered {
		riskStyle := riskStyle(s, a.Risk)
		confPct := int(a.Confidence * 100)
		confStyle := confidenceStyle(s, confPct)
		text := a.Description
		rows = append(rows, fmt.Sprintf("  %s %s  %s  %s",
			s.DimText.Render(a.ID),
			riskStyle.Render(padRight(string(a.Risk), 10)),
			confStyle.Render(fmt.Sprintf("%3d%%", confPct)),
			text))
	}
	return s.Card(sectionTitle, strings.Join(rows, "\n"), width)
}

func renderResultVerification(s StyleSet, r *AnalysisResult, width int) string {
	if r.VerificationOutput == nil || r.VerificationOutput.Assessment == nil {
		return s.Card("", s.EmptyState.Render("No verification data available."), width)
	}
	a := r.VerificationOutput.Assessment

	cardW := (width - 12) / 4
	if cardW < 10 {
		cardW = 10
	}
	statusCards := lipgloss.JoinHorizontal(lipgloss.Top,
		s.StatusCardLarge("VERIFIED", a.VerifiedCount, "", cardW),
		strings.Repeat(" ", 2),
		s.StatusCardLarge("PARTIAL", a.PartialCount, "", cardW),
		strings.Repeat(" ", 2),
		s.StatusCardLarge("UNVERIFIED", a.UnverifiedCount, "", cardW),
		strings.Repeat(" ", 2),
		s.StatusCardLarge("NO EVIDENCE", a.NoEvidenceCount, "", cardW),
	)

	var rows []string
	rows = append(rows, statusCards)

	if a.OverallConfidence > 0 {
		confPct := a.OverallConfidence * 100
		rows = append(rows, "")
		rows = append(rows, s.SubHeader("Overall Confidence"))
		rows = append(rows, "  "+s.ProgressWithLabel(confPct, width-8))
		rows = append(rows, "")
	}

	if r.VerificationOutput.CISOView != nil {
		cv := r.VerificationOutput.CISOView
		if len(cv.TopAssumptionsToVerify) > 0 {
			var items []string
			for _, plan := range cv.TopAssumptionsToVerify {
				items = append(items, "  • "+s.DimText.Render(plan.AssumptionText))
			}
			rows = append(rows, s.Card("Top Items to Verify", strings.Join(items, "\n"), width-4))
		}
		if len(cv.EvidenceGaps) > 0 {
			var gaps []string
			for _, g := range cv.EvidenceGaps {
				gaps = append(gaps, "  • "+s.StatusWarn.Render(g))
			}
			rows = append(rows, s.CardAccent("Evidence Gaps", strings.Join(gaps, "\n"), width-4))
		}
	}

	return s.Card("Verification Assessment", strings.Join(rows, "\n"), width)
}

func renderResultContradictions(s StyleSet, r *AnalysisResult, searchQuery string, width int) string {
	if len(r.Contradictions) == 0 {
		return s.Card("", s.EmptyState.Render("No contradictions detected."), width)
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
		return s.Card("", s.EmptyState.Render(fmt.Sprintf("No contradictions match \"%s\".", searchQuery)), width)
	}

	var cards []string
	for _, c := range filtered {
		severityStyle := s.StatusWarn
		if c.Severity == RiskHigh || c.Severity == RiskCritical {
			severityStyle = s.StatusBad
		} else if c.Severity == RiskLow {
			severityStyle = s.StatusGood
		}
		var detail []string
		detail = append(detail, fmt.Sprintf("  %s", severityStyle.Render(strings.ToUpper(string(c.Severity)))))
		detail = append(detail, "  "+s.Value.Render(c.Description))
		if c.Explanation != "" {
			detail = append(detail, "  "+s.DimText.Render("Reason: "+c.Explanation))
		}
		if len(c.AffectedAssumptions) > 0 {
			detail = append(detail, "  "+s.DimText.Render("Affects: "+strings.Join(c.AffectedAssumptions, ", ")))
		}
		cards = append(cards, s.CardAccent("", strings.Join(detail, "\n"), width))
	}

	title := s.SubSectionTitle.Render(sectionTitle)
	return title + "\n" + strings.Join(cards, "\n")
}

func renderResultTrust(s StyleSet, r *AnalysisResult, searchQuery string, width int) string {
	if r.TrustOutput == nil {
		return s.EmptyState.Render("No trust chain data available.")
	}
	var rows []string
	q := strings.ToLower(searchQuery)

	if len(r.TrustOutput.TrustChains) > 0 {
		total := len(r.TrustOutput.TrustChains)
		matchCount := 0

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

			var nodes []string
			if len(chain.Nodes) > 0 {
				nodes = chain.Nodes
			} else {
				nodes = []string{chain.RootNode, chain.LeafNode}
			}
			diagram := s.TrustDiagram(nodes)

			info := fmt.Sprintf("Chain %s  |  Length: %d  |  Confidence: %d%%",
				s.DimText.Render(chain.ID),
				chain.Length,
				confidence)
			cardContent := diagram + "\n" + info
			rows = append(rows, s.Card("", cardContent, 30))
			rows = append(rows, "")
		}

		sectionTitle := fmt.Sprintf("Trust Chains (%d)", total)
		if searchQuery != "" {
			sectionTitle = fmt.Sprintf("Trust Chains (%d of %d matching \"%s\")", matchCount, total, searchQuery)
		}
		if matchCount == 0 {
			rows = append(rows, s.EmptyState.Render(fmt.Sprintf("No trust chains match \"%s\".", searchQuery)))
		}
		title := s.SubSectionTitle.Render(sectionTitle)
		rows = append([]string{title}, rows...)
	}

	if r.TrustOutput != nil && len(r.TrustOutput.SinglePointsOfTrust) > 0 {
		var spofs []string
		for _, spof := range r.TrustOutput.SinglePointsOfTrust {
			spofs = append(spofs, fmt.Sprintf("  %s", s.DimText.Render(spof.AssumptionText)))
		}
		rows = append(rows, s.CardAccent("Single Points of Trust Failure", strings.Join(spofs, "\n"), width))
	}

	if r.ReviewOutput != nil && r.ReviewOutput.Queue != nil && len(r.ReviewOutput.Queue.Items) > 0 {
		var items []string
		for _, item := range r.ReviewOutput.Queue.Items {
			riskStyle := riskStyle(s, RiskLevel(item.Risk))
			items = append(items, fmt.Sprintf("  #%d %s [%.0f] %s",
				item.Rank, riskStyle.Render(item.Risk), item.PriorityScore, s.DimText.Render(item.AssumptionText)))
		}
		rows = append(rows, s.Card("Priority Queue", strings.Join(items, "\n"), width))
	}

	if r.ReviewOutput != nil && r.ReviewOutput.CISODashboard != nil {
		d := r.ReviewOutput.CISODashboard
		ciso := fmt.Sprintf("  %s %d  %s %d",
			s.BadgeCritical.Render("CRITICAL"), d.CriticalAssumptions,
			s.BadgeHigh.Render("HIGH"), d.HighAssumptions)
		rows = append(rows, s.Card("CISO View", ciso, width))
	}

	if len(rows) == 0 {
		return s.EmptyState.Render("No trust data available.")
	}
	return strings.Join(rows, "\n")
}

func renderResultSPOFs(s StyleSet, r *AnalysisResult) string {
	if r.TrustOutput == nil || len(r.TrustOutput.SinglePointsOfTrust) == 0 {
		return s.EmptyState.Render("No single points of trust failure identified.")
	}
	var rows []string
	rows = append(rows, s.SubSectionTitle.Render(fmt.Sprintf("Single Points of Trust Failure (%d)", len(r.TrustOutput.SinglePointsOfTrust))))
	for _, spof := range r.TrustOutput.SinglePointsOfTrust {
		rows = append(rows, s.SPOFDiagram(spof.AssumptionText, 3, "High"))
		rows = append(rows, "")
	}
	return strings.Join(rows, "\n")
}

func renderResultImpact(s StyleSet, r *AnalysisResult, width int) string {
	var rows []string

	if r.ReviewOutput != nil && r.ReviewOutput.Queue != nil && len(r.ReviewOutput.Queue.Items) > 0 {
		var items []string
		for _, item := range r.ReviewOutput.Queue.Items {
			riskStyle := riskStyle(s, RiskLevel(item.Risk))
			items = append(items, fmt.Sprintf("  #%d %s [%.0f] %s",
				item.Rank, riskStyle.Render(item.Risk), item.PriorityScore, s.DimText.Render(item.AssumptionText)))
		}
		rows = append(rows, s.Card("Priority Queue", strings.Join(items, "\n"), width))
	}

	if r.TrustOutput != nil && len(r.TrustOutput.SinglePointsOfTrust) > 0 {
		var spofs []string
		for _, spof := range r.TrustOutput.SinglePointsOfTrust {
			spofs = append(spofs, fmt.Sprintf("  ⚠ %s", s.DimText.Render(spof.AssumptionText)))
		}
		rows = append(rows, s.CardAccent("Single Points of Trust Failure", strings.Join(spofs, "\n"), width))
	}

	if r.ReviewOutput != nil && r.ReviewOutput.CISODashboard != nil {
		d := r.ReviewOutput.CISODashboard
		ciso := fmt.Sprintf("  %s %d  %s %d",
			s.BadgeCritical.Render("CRITICAL"), d.CriticalAssumptions,
			s.BadgeHigh.Render("HIGH"), d.HighAssumptions)
		rows = append(rows, s.Card("CISO View", ciso, width))
	}

	if len(rows) == 0 {
		return s.Card("", s.EmptyState.Render("No impact data available."), width)
	}
	return strings.Join(rows, "\n")
}

func renderResultBlindSpots(s StyleSet, r *AnalysisResult, width int) string {
	if r.CoverageOutput == nil {
		return s.Card("", s.EmptyState.Render("No blind spot data available."), width)
	}
	var cards []string

	if len(r.CoverageOutput.BlindSpots) > 0 {
		var items []string
		for _, bs := range r.CoverageOutput.BlindSpots {
			items = append(items, "  • "+s.DimText.Render(bs.Title))
		}
		cards = append(cards, s.Card("Blind Spots", strings.Join(items, "\n"), width))
	}

	if len(r.CoverageOutput.DomainBlindSpots) > 0 {
		var items []string
		for _, dbs := range r.CoverageOutput.DomainBlindSpots {
			items = append(items, "  • "+s.DimText.Render(dbs.MissingArea))
		}
		cards = append(cards, s.Card("Domain Blind Spots", strings.Join(items, "\n"), width))
	}

	if r.CoverageOutput.Assessment != nil && len(r.CoverageOutput.Assessment.Gaps) > 0 {
		var items []string
		for _, gap := range r.CoverageOutput.Assessment.Gaps {
			items = append(items, fmt.Sprintf("  %s: %s (observed: %d, expected: %d)",
				s.StatusWarn.Render(string(gap.Category)), gap.Risk, gap.ObservedCount, gap.ExpectedCount))
		}
		cards = append(cards, s.CardAccent("Coverage Gaps", strings.Join(items, "\n"), width))
	}

	if len(cards) == 0 {
		return s.Card("", s.StatusGood.Render("No blind spots or gaps detected."), width)
	}
	return strings.Join(cards, "\n")
}

func renderResultControls(s StyleSet, r *AnalysisResult, searchQuery string, width int) string {
	if len(r.Controls) == 0 {
		return s.Card("", s.EmptyState.Render("No recommended controls."), width)
	}

	q := strings.ToLower(searchQuery)
	matchCount := 0
	var controlCards []string

	for _, c := range r.Controls {
		if searchQuery != "" {
			match := strings.Contains(strings.ToLower(c.Description), q) ||
				strings.Contains(strings.ToLower(c.Rationale), q)
			if !match {
				continue
			}
		}
		matchCount++
		var details []string
		details = append(details, "  "+s.Value.Render(c.Description))
		if c.Rationale != "" {
			details = append(details, "  "+s.DimText.Render(c.Rationale))
		}
		if len(c.MitigatedAssumptionIDs) > 0 {
			details = append(details, "  "+s.DimText.Render("Mitigates: "+strings.Join(c.MitigatedAssumptionIDs, ", ")))
		}
		controlCards = append(controlCards, s.Card("", strings.Join(details, "\n"), width))
	}

	sectionTitle := fmt.Sprintf("Recommended Controls (%d)", len(r.Controls))
	if searchQuery != "" {
		sectionTitle = fmt.Sprintf("Recommended Controls (%d of %d matching \"%s\")", matchCount, len(r.Controls), searchQuery)
	}
	if matchCount == 0 {
		return s.Card("", s.EmptyState.Render(fmt.Sprintf("No controls match \"%s\".", searchQuery)), width)
	}

	return s.SubSectionTitle.Render(sectionTitle) + "\n" + strings.Join(controlCards, "\n")
}

func renderResultReports(s StyleSet, r *AnalysisResult, width int) string {
	var cards []string

	if r.NarrativeOutput != nil {
		cards = append(cards, s.StatusCardLarge("NARRATIVE", 0, "Architect Narrative available — press 'e' to export", width))
	}

	if r.ReviewOutput != nil && r.ReviewOutput.Campaigns != nil && len(r.ReviewOutput.Campaigns) > 0 {
		var items []string
		for _, c := range r.ReviewOutput.Campaigns {
			items = append(items, "  ▶ "+s.Value.Render(c.Name)+" "+s.DimText.Render(fmt.Sprintf("(%d items)", len(c.Items))))
		}
		cards = append(cards, s.Card("Review Campaigns", strings.Join(items, "\n"), width))
	}

	if r.ConfidenceOutput != nil && len(r.ConfidenceOutput.Breakdowns) > 0 {
		highConf := 0
		lowConf := 0
		for _, bd := range r.ConfidenceOutput.Breakdowns {
			if bd.FinalConfidence >= 70 {
				highConf++
			} else if bd.FinalConfidence < 40 {
				lowConf++
			}
		}
		conf := fmt.Sprintf("  %s %d  %s %d",
			s.StatusGood.Render("High Confidence:"), highConf,
			s.StatusBad.Render("Low Confidence:"), lowConf)
		cards = append(cards, s.Card("Confidence Overview", conf, width))
	}

	cards = append(cards, s.Card("Export",
		"  Press 'e' to open export dialog.\n  Formats: JSON, Markdown, HTML, CSV, PDF",
		width))

	return strings.Join(cards, "\n")
}

func renderResultSDRI(s StyleSet, r *AnalysisResult, width int) string {
	if r.SDRISummary == "" && len(r.SDRIControls) == 0 &&
		len(r.SDRIDesignFindings) == 0 && len(r.SDRIAchitecturalWeaknesses) == 0 &&
		len(r.SDRIRemediations) == 0 && len(r.SDRIComplianceAlignments) == 0 {
		return s.Card("", s.EmptyState.Render("No SDRI data available."), width)
	}
	var cards []string

	if r.SDRISummary != "" {
		cards = append(cards, s.Card("Executive Summary", "  "+s.DimText.Render(r.SDRISummary), width))
	}

	if len(r.SDRIControls) > 0 {
		byStatus := map[string]int{}
		for _, c := range r.SDRIControls {
			status := c.Status
			if status == "" {
				status = "unknown"
			}
			byStatus[status]++
		}
		var items []string
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
			items = append(items, fmt.Sprintf("  %s %s: %d", style.Render("●"), st, n))
		}
		cards = append(cards, s.Card("Control Inventory", strings.Join(items, "\n"), width))
	}

	if len(r.SDRICoverageByCategory) > 0 {
		var items []string
		for _, c := range r.SDRICoverageByCategory {
			style := s.StatusGood
			if c.Coverage < 50 {
				style = s.StatusBad
			} else if c.Coverage < 80 {
				style = s.StatusWarn
			}
			items = append(items, fmt.Sprintf("  %s %s: %.0f%% (%d/%d)",
				style.Render("●"), c.Category, c.Coverage*100, c.Observed, c.Expected))
		}
		cards = append(cards, s.Card("Coverage by Category", strings.Join(items, "\n"), width))
	}

	if r.SDRICoverageDashboard != nil && len(r.SDRICoverageDashboard) > 0 {
		var items []string
		for k, v := range r.SDRICoverageDashboard {
			pct := int(v * 100)
			style := confidenceStyle(s, pct)
			items = append(items, fmt.Sprintf("  %s: %s", k, style.Render(fmt.Sprintf("%d%%", pct))))
		}
		cards = append(cards, s.Card("Coverage Dashboard", strings.Join(items, "\n"), width))
	}

	if len(r.SDRIDesignFindings) > 0 {
		var items []string
		for _, f := range r.SDRIDesignFindings {
			severityStyle := s.StatusWarn
			switch f.Severity {
			case "Critical", "High":
				severityStyle = s.StatusBad
			case "Low":
				severityStyle = s.StatusGood
			}
			items = append(items, fmt.Sprintf("  %s %s %s",
				severityStyle.Render("●"), s.RiskBadge(f.Severity), s.Value.Render(f.Title)))
			if f.Description != "" {
				items = append(items, "    "+s.DimText.Render(f.Description))
			}
			if f.Recommendation != "" {
				items = append(items, "    → "+s.DimText.Render(f.Recommendation))
			}
		}
		cards = append(cards, s.Card("Design Findings", strings.Join(items, "\n"), width))
	}

	if len(r.SDRIAchitecturalWeaknesses) > 0 {
		var items []string
		for _, w := range r.SDRIAchitecturalWeaknesses {
			severityStyle := s.StatusWarn
			switch w.Severity {
			case "Critical", "High":
				severityStyle = s.StatusBad
			case "Low":
				severityStyle = s.StatusGood
			}
			items = append(items, fmt.Sprintf("  %s %s %s",
				severityStyle.Render("●"), s.RiskBadge(w.Severity), s.Value.Render(w.Pattern)))
			if w.Description != "" {
				items = append(items, "    "+s.DimText.Render(w.Description))
			}
			if w.Recommendation != "" {
				items = append(items, "    → "+s.DimText.Render(w.Recommendation))
			}
		}
		cards = append(cards, s.CardAccent("Architectural Weaknesses", strings.Join(items, "\n"), width))
	}

	if len(r.SDRIRemediations) > 0 {
		var items []string
		for _, rem := range r.SDRIRemediations {
			items = append(items, fmt.Sprintf("  #%d [%.0f] %s",
				rem.Priority, rem.RiskScore, s.Value.Render(rem.Description)))
			if rem.Recommendation != "" {
				items = append(items, "    → "+s.DimText.Render(rem.Recommendation)+" "+s.DimText.Render("(effort: "+rem.Effort+")"))
			}
		}
		cards = append(cards, s.Card("Remediations", strings.Join(items, "\n"), width))
	}

	if len(r.SDRIComplianceAlignments) > 0 {
		var items []string
		for _, m := range r.SDRIComplianceAlignments {
			style := s.StatusGood
			if m.Coverage < 50 {
				style = s.StatusBad
			} else if m.Coverage < 80 {
				style = s.StatusWarn
			}
			items = append(items, fmt.Sprintf("  %s %s: %.0f%% (%s)",
				style.Render("●"), m.Framework, m.Coverage, m.Status))
		}
		cards = append(cards, s.Card("Compliance Alignments", strings.Join(items, "\n"), width))
	}

	return strings.Join(cards, "\n")
}

func renderResultSecurityDesignReview(s StyleSet, r *AnalysisResult, width int) string {
	if len(r.SDRIDesignFindings) == 0 && len(r.SDRIAchitecturalWeaknesses) == 0 &&
		len(r.SDRIRemediations) == 0 && len(r.SDRIComplianceAlignments) == 0 {
		return s.Card("", s.EmptyState.Render("No security design review data available."), width)
	}
	var cards []string

	if len(r.SDRIDesignFindings) > 0 {
		var items []string
		for _, f := range r.SDRIDesignFindings {
			severityStyle := s.StatusWarn
			switch f.Severity {
			case "Critical", "High":
				severityStyle = s.StatusBad
			case "Low":
				severityStyle = s.StatusGood
			}
			items = append(items, fmt.Sprintf("  %s %s %s",
				severityStyle.Render("●"), s.RiskBadge(f.Severity), s.Value.Render(f.Title)))
			if f.Description != "" {
				items = append(items, "    "+s.DimText.Render(f.Description))
			}
			if f.Recommendation != "" {
				items = append(items, "    → "+s.DimText.Render(f.Recommendation))
			}
		}
		cards = append(cards, s.Card("Design Findings", strings.Join(items, "\n"), width))
	}

	if len(r.SDRIAchitecturalWeaknesses) > 0 {
		var items []string
		for _, w := range r.SDRIAchitecturalWeaknesses {
			severityStyle := s.StatusWarn
			switch w.Severity {
			case "Critical", "High":
				severityStyle = s.StatusBad
			case "Low":
				severityStyle = s.StatusGood
			}
			items = append(items, fmt.Sprintf("  %s %s %s",
				severityStyle.Render("●"), s.RiskBadge(w.Severity), s.Value.Render(w.Pattern)))
			if w.Description != "" {
				items = append(items, "    "+s.DimText.Render(w.Description))
			}
			if w.Recommendation != "" {
				items = append(items, "    → "+s.DimText.Render(w.Recommendation))
			}
		}
		cards = append(cards, s.CardAccent("Architectural Weaknesses", strings.Join(items, "\n"), width))
	}

	if len(r.SDRIRemediations) > 0 {
		var items []string
		for _, rem := range r.SDRIRemediations {
			items = append(items, fmt.Sprintf("  #%d [%.0f] %s",
				rem.Priority, rem.RiskScore, s.Value.Render(rem.Description)))
			if rem.Recommendation != "" {
				items = append(items, "    → "+s.DimText.Render(rem.Recommendation)+" "+s.DimText.Render("(effort: "+rem.Effort+")"))
			}
		}
		cards = append(cards, s.Card("Remediations", strings.Join(items, "\n"), width))
	}

	if len(r.SDRIComplianceAlignments) > 0 {
		var items []string
		for _, m := range r.SDRIComplianceAlignments {
			style := s.StatusGood
			if m.Coverage < 50 {
				style = s.StatusBad
			} else if m.Coverage < 80 {
				style = s.StatusWarn
			}
			items = append(items, fmt.Sprintf("  %s %s: %.0f%% (%s)",
				style.Render("●"), m.Framework, m.Coverage, m.Status))
		}
		cards = append(cards, s.Card("Compliance Alignments", strings.Join(items, "\n"), width))
	}

	return strings.Join(cards, "\n")
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
