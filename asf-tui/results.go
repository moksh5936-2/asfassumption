package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tabState struct {
	selectedIndex int
	detailOpen    bool
	searchQuery   string
	filterActive  bool
	showHelp      bool
	selectedLine  int // rendered line offset of selected item within tab content
	contentOffset int // number of viewport lines before tab content starts
}

type resultsModel struct {
	result    *AnalysisResult
	resultTab int
	tabs      []resultTabDef
	tabScroll map[int]int
	tabStates map[int]*tabState
	vpReady   bool
	vpWidth   int
	vpHeight  int
}

func (m *resultsModel) tabStateFor(tab int) *tabState {
	if _, ok := m.tabStates[tab]; !ok {
		m.tabStates[tab] = &tabState{}
	}
	return m.tabStates[tab]
}

func (m *resultsModel) tabCount(tab int) int {
	return resultTabCount(m.result, tab)
}

func (m *resultsModel) tabCountString(tab int) string {
	if m.result == nil {
		return ""
	}
	switch tab {
	case 0:
		return ""
	case 1:
		return fmt.Sprintf("%d assumptions", len(m.result.Assumptions))
	case 2:
		if m.result.VerificationOutput == nil || m.result.VerificationOutput.Assessment == nil {
			return "Verification"
		}
		a := m.result.VerificationOutput.Assessment
		parts := []string{}
		if a.VerifiedCount > 0 {
			parts = append(parts, fmt.Sprintf("%d verified", a.VerifiedCount))
		}
		if a.PartialCount > 0 {
			parts = append(parts, fmt.Sprintf("%d partial", a.PartialCount))
		}
		if a.UnverifiedCount > 0 {
			parts = append(parts, fmt.Sprintf("%d unverified", a.UnverifiedCount))
		}
		if a.NoEvidenceCount > 0 {
			parts = append(parts, fmt.Sprintf("%d no evidence", a.NoEvidenceCount))
		}
		if len(parts) == 0 {
			return "Verification"
		}
		return strings.Join(parts, " / ")
	case 3:
		return fmt.Sprintf("%d contradictions", len(m.result.Contradictions))
	case 4:
		t := m.result.TrustOutput
		if t == nil {
			return "Trust"
		}
		parts := []string{}
		if len(t.TrustChains) > 0 {
			parts = append(parts, fmt.Sprintf("%d chain(s)", len(t.TrustChains)))
		}
		if len(t.SinglePointsOfTrust) > 0 {
			parts = append(parts, fmt.Sprintf("%d SPOF(s)", len(t.SinglePointsOfTrust)))
		}
		if m.result.ReviewOutput != nil && m.result.ReviewOutput.Queue != nil && len(m.result.ReviewOutput.Queue.Items) > 0 {
			parts = append(parts, fmt.Sprintf("%d in queue", len(m.result.ReviewOutput.Queue.Items)))
		}
		if len(t.FailureCascades) > 0 {
			parts = append(parts, fmt.Sprintf("%d cascade(s)", len(t.FailureCascades)))
		}
		if len(parts) == 0 {
			return "Trust"
		}
		return strings.Join(parts, " / ")
	case 5:
		return fmt.Sprintf("%d control(s)", len(m.result.Controls))
	case 6:
		count := len(m.result.SDRIControls) + len(m.result.SDRIDesignFindings) + len(m.result.SDRIAchitecturalWeaknesses) + len(m.result.SDRIRemediations) + len(m.result.SDRIComplianceAlignments)
		return fmt.Sprintf("%d items", count)
	default:
		return ""
	}
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
		tabStates: make(map[int]*tabState),
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
	var content string
	tab := m.results.resultTab
	ts := m.results.tabStateFor(tab)
	switch tab {
	case 0:
		content = renderResultSummary(s, r, ts, m.mainWidth()-4)
	case 1:
		content = renderResultAssumptions(s, r, ts, m.mainWidth()-4)
	case 2:
		content = renderResultVerification(s, r, ts, m.mainWidth()-4)
	case 3:
		content = renderResultContradictions(s, r, ts, m.mainWidth()-4)
	case 4:
		content = renderResultTrust(s, r, ts, m.mainWidth()-4)
	case 5:
		content = renderResultControls(s, r, ts, m.mainWidth()-4)
	case 6:
		content = renderResultSDRI(s, r, ts, m.mainWidth()-4)
	}

	ts.selectedLine = 0
	tabBar := m.renderResultTabs()
	breadcrumb := m.renderBreadcrumb(tab, ts)
	sectionName := m.results.tabs[tab].name
	countStr := m.results.tabCountString(tab)
	titleStr := sectionName
	if countStr != "" {
		titleStr = sectionName + " — " + countStr
	}
	header := s.PremiumHeader(titleStr, m.mainWidth())
	ts.contentOffset = strings.Count(lipgloss.JoinVertical(lipgloss.Left, header, tabBar, breadcrumb), "\n") + 1
	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		tabBar,
		breadcrumb,
		content,
	)
}

func (m mainModel) renderBreadcrumb(tab int, ts *tabState) string {
	s := m.styles
	if m.activeCase == "" {
		return ""
	}
	caseName := fileBase(m.activeCase)
	tabName := m.results.tabs[tab].name
	parts := []string{
		s.Breadcrumb.Render("ASF0"),
		s.BreadcrumbSep.Render(" / "),
		s.BreadcrumbSep.Render("Case: "),
		s.Breadcrumb.Render(caseName),
		s.BreadcrumbSep.Render(" / "),
		s.Breadcrumb.Render(tabName),
	}
	if tab > 0 && ts.selectedIndex >= 0 && ts.selectedIndex < m.results.tabCount(tab) {
		itemStr := fmt.Sprintf(" / #%d", ts.selectedIndex+1)
		parts = append(parts, s.BreadcrumbSep.Render(itemStr))
	}
	if ts.detailOpen {
		parts = append(parts, s.BreadcrumbSep.Render(" / "), s.BreadcrumbSep.Render("detail"))
	}
	if ts.filterActive || ts.searchQuery != "" {
		parts = append(parts, s.BreadcrumbSep.Render(" / "), s.SearchActive.Render("filter:"+ts.searchQuery))
	}
	return strings.Join(parts, "")
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
	switch msg := msg.(type) {
	case tea.MouseMsg:
		tab := m.results.resultTab
		if tab > 0 && m.router.focus == focusContent {
			ts := m.results.tabStateFor(tab)
			maxIdx := m.results.tabCount(tab) - 1
			switch msg.Type {
			case tea.MouseWheelUp:
				if ts.selectedIndex > 0 {
					ts.selectedIndex--
				}
				return m, nil
			case tea.MouseWheelDown:
				if ts.selectedIndex < maxIdx {
					ts.selectedIndex++
				}
				return m, nil
			}
		}
		return m, nil
	case tea.KeyMsg:
		tab := m.results.resultTab
		if tab > 0 && m.router.focus == focusContent {
			ts := m.results.tabStateFor(tab)
			maxIdx := m.results.tabCount(tab) - 1
			switch msg.String() {
			case "up", "k":
				if !ts.filterActive && ts.selectedIndex > 0 {
					ts.selectedIndex--
				}
				return m, nil
			case "down", "j":
				if !ts.filterActive && ts.selectedIndex < maxIdx {
					ts.selectedIndex++
				}
				return m, nil
			case "enter":
				if !ts.filterActive && maxIdx >= 0 {
					ts.detailOpen = !ts.detailOpen
				}
				return m, nil
			case "/":
				ts.filterActive = true
				ts.searchQuery = ""
				return m, nil
			case "n":
				if ts.filterActive && ts.searchQuery != "" && ts.selectedIndex < maxIdx {
					ts.selectedIndex++
				}
				return m, nil
			case "N":
				if ts.filterActive && ts.searchQuery != "" && ts.selectedIndex > 0 {
					ts.selectedIndex--
				}
				return m, nil
			case "backspace":
				if ts.filterActive && len(ts.searchQuery) > 0 {
					ts.searchQuery = ts.searchQuery[:len(ts.searchQuery)-1]
				}
				return m, nil
			case "esc":
				if ts.showHelp {
					ts.showHelp = false
					return m, nil
				}
				if ts.detailOpen {
					ts.detailOpen = false
					return m, nil
				}
				if ts.filterActive {
					ts.filterActive = false
					ts.searchQuery = ""
					if ts.selectedIndex >= m.results.tabCount(tab) {
						ts.selectedIndex = 0
					}
					return m, nil
				}
			default:
				if ts.filterActive && len(msg.String()) == 1 && msg.String() != "\t" {
					ts.searchQuery += msg.String()
					return m, nil
				}
			}
		} else if tab == 0 {
			switch msg.String() {
			case "up", "k", "down", "j":
				return m, nil
			}
		}
	}

	tabSwitch := m.results.resultTab
	var cmd tea.Cmd
	m.results, cmd = m.results.Update(msg)
	if m.results.resultTab != tabSwitch {
		m.results.tabScroll[tabSwitch] = m.vp.YOffset
		if y, ok := m.results.tabScroll[m.results.resultTab]; ok {
			m.vp.YOffset = y
		} else {
			m.vp.YOffset = 0
		}
	}
	return m, cmd
}

func renderResultSummary(s StyleSet, r *AnalysisResult, ts *tabState, width int) string {
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

	wflow := s.renderWorkflow(r)

	result := lipgloss.JoinVertical(lipgloss.Left,
		infoCard,
		riskCard,
		verifyCard,
		contradictCard,
		trustCard,
		coverageCard,
		wflow,
	)

	if ts.showHelp {
		helpText := lipgloss.JoinVertical(lipgloss.Left,
			s.DimText.Render("  Overview shows a summary of the analysis results, risk distribution, and quick stats."),
			s.DimText.Render("  Use ←→/hl to switch tabs. Press r to Review, v to Validate, e to Export."),
			"",
		)
		return lipgloss.JoinVertical(lipgloss.Left, helpText, result)
	}
	return result
}

func renderResultAssumptions(s StyleSet, r *AnalysisResult, ts *tabState, width int) string {
	if len(r.Assumptions) == 0 {
		return s.Card("", s.EmptyState.Render("No assumptions found."), width)
	}

	q := strings.ToLower(ts.searchQuery)
	var filtered []Assumption
	for _, a := range r.Assumptions {
		if ts.filterActive && q != "" {
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
	if ts.filterActive && q != "" {
		sectionTitle = fmt.Sprintf("Assumptions (%d of %d matching \"%s\")", len(filtered), len(r.Assumptions), q)
	}
	if len(filtered) == 0 {
		return s.Card("", s.EmptyState.Render(fmt.Sprintf("No assumptions match \"%s\".", ts.searchQuery)), width)
	}

	var rows []string
	if ts.showHelp {
		rows = append(rows, s.DimText.Render("  Assumptions are claims extracted from the architecture that may impact security."))
		rows = append(rows, s.DimText.Render("  ↑↓ Select  |  Enter Detail  |  r Review  |  v Validate  |  / Search"))
		rows = append(rows, "")
	}
	rows = append(rows, s.SectionRule.Render(strings.Repeat("─", max(1, width-4))))
	rows = append(rows, s.SubSectionTitle.Render("ID  Risk        Confidence  Description"))

	for i, a := range filtered {
		selected := i == ts.selectedIndex
		prefix := "  "
		style := s.Value
		if selected {
			prefix = "▸ "
			style = s.SelectedItem
			ts.selectedLine = len(rows)
		}
		riskStyle := riskStyle(s, a.Risk)
		confPct := int(a.Confidence * 100)
		confStyle := confidenceStyle(s, confPct)
		text := a.Description
		rows = append(rows, style.Render(fmt.Sprintf("%s%s %s  %s  %s",
			prefix, a.ID,
			riskStyle.Render(padRight(string(a.Risk), 10)),
			confStyle.Render(fmt.Sprintf("%3d%%", confPct)),
			text)))
	}

	if ts.detailOpen && ts.selectedIndex >= 0 && ts.selectedIndex < len(filtered) {
		a := filtered[ts.selectedIndex]
		var detail []string
		detail = append(detail, "")
		detail = append(detail, s.DetailBox.Render(fmt.Sprintf("  ID: %s", a.ID)))
		detail = append(detail, s.DetailBox.Render(fmt.Sprintf("  Risk: %s", a.Risk)))
		detail = append(detail, s.DetailBox.Render(fmt.Sprintf("  Confidence: %.0f%%", a.Confidence*100)))
		if a.Component != "" {
			detail = append(detail, s.DetailBox.Render(fmt.Sprintf("  Component: %s", a.Component)))
		}
		if a.Category != "" {
			detail = append(detail, s.DetailBox.Render(fmt.Sprintf("  Category: %s", a.Category)))
		}
		if a.VerificationStatus != "" {
			detail = append(detail, s.DetailBox.Render(fmt.Sprintf("  Status: %s", a.VerificationStatus)))
		}
		if len(a.Stride) > 0 {
			strideStrs := make([]string, len(a.Stride))
			for i, st := range a.Stride {
				strideStrs[i] = string(st)
			}
			detail = append(detail, s.DetailBox.Render(fmt.Sprintf("  STRIDE: %s", strings.Join(strideStrs, ", "))))
		}
		if len(a.EvidenceSources) > 0 {
			detail = append(detail, s.DetailBox.Render(fmt.Sprintf("  Evidence: %s", strings.Join(a.EvidenceSources, ", "))))
		}
		if a.Rationale != "" {
			detail = append(detail, s.DetailBox.Render(fmt.Sprintf("  Rationale: %s", a.Rationale)))
		}
		rows = append(rows, detail...)
	}

	return s.Card(sectionTitle, strings.Join(rows, "\n"), width)
}

func renderResultVerification(s StyleSet, r *AnalysisResult, ts *tabState, width int) string {
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
	if ts.showHelp {
		rows = append(rows, s.DimText.Render("  Verification assesses how well each assumption is supported by evidence."))
		rows = append(rows, s.DimText.Render("  ↑↓ Select  |  Enter Detail  |  r Review  |  v Validate  |  / Search"))
		rows = append(rows, "")
	}
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
		q := strings.ToLower(ts.searchQuery)
		var filteredItems []struct {
			label           string
			text            string
			howToValidate   string
			evidenceMissing string
			group           string
		}
		for _, plan := range cv.TopAssumptionsToVerify {
			label := plan.AssumptionText
			if ts.filterActive && q != "" {
				if !strings.Contains(strings.ToLower(label), q) {
					continue
				}
			}
			var emStr string
			if len(plan.EvidenceMissing) > 0 {
				var missing []string
				for _, em := range plan.EvidenceMissing {
					name := em.Name
					if name == "" {
						name = em.Description
					}
					if name != "" {
						missing = append(missing, name)
					}
				}
				emStr = strings.Join(missing, ", ")
			}
			group := string(plan.Status)
			if group == "" {
				group = "Unknown"
			}
			filteredItems = append(filteredItems, struct {
				label           string
				text            string
				howToValidate   string
				evidenceMissing string
				group           string
			}{label: label, text: "", howToValidate: plan.HowToValidate, evidenceMissing: emStr, group: group})
		}
		for _, g := range cv.EvidenceGaps {
			label := g
			if ts.filterActive && q != "" {
				if !strings.Contains(strings.ToLower(label), q) {
					continue
				}
			}
			filteredItems = append(filteredItems, struct {
				label           string
				text            string
				howToValidate   string
				evidenceMissing string
				group           string
			}{label: "⚠ " + label, text: "warn", group: "Evidence Gaps"})
		}

		if len(filteredItems) > 0 {
			prevGroup := ""
			for i, item := range filteredItems {
				selected := i == ts.selectedIndex
				prefix := "  • "
				style := s.Value
				if selected {
					prefix = "▸ "
					style = s.SelectedItem
				}

				if item.group != prevGroup {
					if prevGroup != "" {
						rows = append(rows, "")
					}
					rows = append(rows, s.SubSectionTitle.Render("  "+item.group))
					prevGroup = item.group
				}

				if selected {
					ts.selectedLine = len(rows)
				}
				rows = append(rows, style.Render(prefix+item.label))

				if ts.detailOpen && selected {
					if item.text == "warn" {
						rows = append(rows, "", s.DetailBox.Render("  Evidence Gap — requires additional evidence collection."))
					} else {
						var detail []string
						detail = append(detail, fmt.Sprintf("  Top Priority Assumption to Verify: %s", item.label))
						if item.howToValidate != "" {
							detail = append(detail, "")
							detail = append(detail, fmt.Sprintf("  How to Validate: %s", item.howToValidate))
						}
						if item.evidenceMissing != "" {
							detail = append(detail, "")
							detail = append(detail, fmt.Sprintf("  Evidence Missing: %s", item.evidenceMissing))
						}
						rows = append(rows, "", s.DetailBox.Render(strings.Join(detail, "\n")))
					}
				}
			}
		}
	}

	return s.Card("Verification Assessment", strings.Join(rows, "\n"), width)
}

func renderResultContradictions(s StyleSet, r *AnalysisResult, ts *tabState, width int) string {
	if len(r.Contradictions) == 0 {
		return s.Card("", s.EmptyState.Render("No contradictions detected."), width)
	}

	q := strings.ToLower(ts.searchQuery)
	var filtered []Contradiction
	for _, c := range r.Contradictions {
		if ts.filterActive && q != "" {
			match := strings.Contains(strings.ToLower(c.Description), q) ||
				strings.Contains(strings.ToLower(c.Explanation), q)
			if !match {
				continue
			}
		}
		filtered = append(filtered, c)
	}

	sectionTitle := fmt.Sprintf("Contradictions (%d)", len(r.Contradictions))
	if ts.filterActive && q != "" {
		sectionTitle = fmt.Sprintf("Contradictions (%d of %d matching \"%s\")", len(filtered), len(r.Contradictions), q)
	}
	if len(filtered) == 0 {
		return s.Card("", s.EmptyState.Render(fmt.Sprintf("No contradictions match \"%s\".", ts.searchQuery)), width)
	}

	var items []string
	if ts.showHelp {
		items = append(items, s.DimText.Render("  Contradictions are logical conflicts between assumptions or with known facts."))
		items = append(items, s.DimText.Render("  ↑↓ Select  |  Enter Detail  |  / Search"))
		items = append(items, "")
	}
	for i, c := range filtered {
		selected := i == ts.selectedIndex
		prefix := "  "
		style := s.Value
		if selected {
			prefix = "▸ "
			style = s.SelectedItem
		}

		severityStyle := s.StatusWarn
		if c.Severity == RiskHigh || c.Severity == RiskCritical {
			severityStyle = s.StatusBad
		} else if c.Severity == RiskLow {
			severityStyle = s.StatusGood
		}
		if selected {
			ts.selectedLine = 1 + len(items)
		}
		items = append(items, style.Render(fmt.Sprintf("%s%s %s",
			prefix,
			severityStyle.Render("["+strings.ToUpper(string(c.Severity))+"]"),
			c.Description)))
		if ts.detailOpen && selected {
			var detail []string
			detail = append(detail, "")
			if c.RuleName != "" {
				detail = append(detail, s.DetailBox.Render("  Rule: "+c.RuleName))
			}
			if c.Explanation != "" {
				detail = append(detail, s.DetailBox.Render("  Reason: "+c.Explanation))
			}
			if len(c.AffectedAssumptions) > 0 {
				detail = append(detail, s.DetailBox.Render("  Affects: "+strings.Join(c.AffectedAssumptions, ", ")))
			}
			if len(c.Evidence) > 0 {
				for _, ev := range c.Evidence {
					detail = append(detail, s.DetailBox.Render("  Evidence: "+ev))
				}
			}
			items = append(items, detail...)
		}
	}

	title := s.SubSectionTitle.Render(sectionTitle)
	return title + "\n" + strings.Join(items, "\n")
}

func renderResultTrust(s StyleSet, r *AnalysisResult, ts *tabState, width int) string {
	if r.TrustOutput == nil {
		return s.EmptyState.Render("No trust chain data available.")
	}

	type trustItem struct {
		section string
		label   string
		detail  string
	}

	var items []trustItem
	q := strings.ToLower(ts.searchQuery)

	for _, chain := range r.TrustOutput.TrustChains {
		if ts.filterActive && q != "" {
			match := strings.Contains(strings.ToLower(chain.ID), q) ||
				strings.Contains(strings.ToLower(chain.RootNode), q) ||
				strings.Contains(strings.ToLower(chain.LeafNode), q)
			if !match {
				continue
			}
		}
		confidence := int(chain.Confidence * 100)
		var nodes []string
		if len(chain.Nodes) > 0 {
			nodes = chain.Nodes
		} else {
			nodes = []string{chain.RootNode, chain.LeafNode}
		}
		label := fmt.Sprintf("Chain %s  |  Length: %d  |  Risk: %s  |  Confidence: %d%%", chain.ID, chain.Length, chain.Risk, confidence)
		diagram := s.TrustDiagram(nodes)
		riskNote := ""
		if chain.Risk != "" {
			riskNote = fmt.Sprintf("\nRisk: %s", chain.Risk)
		}
		detail := fmt.Sprintf("Nodes: %s%s", strings.Join(nodes, " → "), riskNote)
		items = append(items, trustItem{section: "chains", label: label, detail: diagram + "\n" + detail})
	}

	if r.TrustOutput != nil && len(r.TrustOutput.SinglePointsOfTrust) > 0 {
		for _, spof := range r.TrustOutput.SinglePointsOfTrust {
			if ts.filterActive && q != "" {
				if !strings.Contains(strings.ToLower(spof.AssumptionText), q) {
					continue
				}
			}
			var detailLines []string
			detailLines = append(detailLines, "Single Point of Trust Failure")
			if spof.DependentsCount > 0 {
				detailLines = append(detailLines, fmt.Sprintf("Dependents: %d", spof.DependentsCount))
			}
			if len(spof.DependentNodes) > 0 {
				detailLines = append(detailLines, "Affected: "+strings.Join(spof.DependentNodes, ", "))
			}
			if spof.Recommendation != "" {
				detailLines = append(detailLines, "Recommendation: "+spof.Recommendation)
			}
			items = append(items, trustItem{section: "spof", label: "⚠ " + spof.AssumptionText, detail: strings.Join(detailLines, "\n")})
		}
	}

	if r.TrustOutput != nil && len(r.TrustOutput.FailureCascades) > 0 {
		for _, fc := range r.TrustOutput.FailureCascades {
			if ts.filterActive && q != "" {
				if !strings.Contains(strings.ToLower(fc.RootAssumptionText), q) {
					continue
				}
			}
			stepCount := len(fc.Steps)
			label := fmt.Sprintf("Cascade: %s  |  Steps: %d  |  Severity: %s", fc.RootAssumptionText, stepCount, fc.Severity)
			detail := fmt.Sprintf("Root: %s (%s)\nSeverity: %s\nTotal Affected: %d\nMax Depth: %d",
				fc.RootAssumptionText, fc.RootAssumptionID, fc.Severity, fc.TotalAffected, fc.MaxDepth)
			if len(fc.Steps) > 0 {
				var steps []string
				for _, step := range fc.Steps {
					s := fmt.Sprintf("  Step %d: %s", step.Step, step.AssumptionText)
					if step.Reason != "" {
						s += " — " + step.Reason
					}
					steps = append(steps, s)
				}
				detail += "\nCascade Chain:\n" + strings.Join(steps, "\n")
			}
			items = append(items, trustItem{section: "cascade", label: label, detail: detail})
		}
	}

	if r.ReviewOutput != nil && r.ReviewOutput.Queue != nil {
		for _, item := range r.ReviewOutput.Queue.Items {
			if ts.filterActive && q != "" {
				if !strings.Contains(strings.ToLower(item.AssumptionText), q) {
					continue
				}
			}
			items = append(items, trustItem{
				section: "queue",
				label:   fmt.Sprintf("#%d [%.0f] %s", item.Rank, item.PriorityScore, item.AssumptionText),
				detail:  fmt.Sprintf("Risk: %s | Priority: %.0f", item.Risk, item.PriorityScore),
			})
		}
	}

	if r.ReviewOutput != nil && r.ReviewOutput.CISODashboard != nil {
		d := r.ReviewOutput.CISODashboard
		items = append(items, trustItem{
			section: "ciso",
			label:   fmt.Sprintf("CISO Dashboard — Critical: %d  High: %d", d.CriticalAssumptions, d.HighAssumptions),
			detail:  "CISO-level summary of critical and high assumptions requiring attention.",
		})
	}

	if len(items) == 0 {
		return s.EmptyState.Render("No trust data available.")
	}

	title := s.SubSectionTitle.Render(fmt.Sprintf("Trust Analysis (%d)", len(items)))
	if ts.filterActive && q != "" {
		title = s.SubSectionTitle.Render(fmt.Sprintf("Trust Analysis (%d matching \"%s\")", len(items), q))
	}

	var rows []string
	if ts.showHelp {
		rows = append(rows, s.DimText.Render("  Trust Analysis maps dependency chains between assumptions and identifies single points of failure."))
		rows = append(rows, s.DimText.Render("  ↑↓ Select  |  Enter Detail  |  / Search"))
		rows = append(rows, "")
	}
	rows = append(rows, title)

	for i, it := range items {
		selected := i == ts.selectedIndex
		prefix := "  "
		style := s.Value
		if selected {
			prefix = "▸ "
			style = s.SelectedItem
		}
		sectionTag := ""
		switch it.section {
		case "chains":
			sectionTag = "[Chain] "
		case "spof":
			sectionTag = "[SPOF] "
		case "cascade":
			sectionTag = "[Cascade] "
		case "queue":
			sectionTag = "[Queue] "
		case "ciso":
			sectionTag = ""
		}
		if selected {
			ts.selectedLine = len(rows)
		}
		rows = append(rows, style.Render(prefix+sectionTag+it.label))

		if ts.detailOpen && selected {
			rows = append(rows, "")
			rows = append(rows, s.DetailBox.Render("  "+it.detail))
			rows = append(rows, "")
		}
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

func renderResultControls(s StyleSet, r *AnalysisResult, ts *tabState, width int) string {
	if len(r.Controls) == 0 {
		return s.Card("", s.EmptyState.Render("No recommended controls."), width)
	}

	q := strings.ToLower(ts.searchQuery)
	var filtered []ControlDetail
	for _, c := range r.Controls {
		if ts.filterActive && q != "" {
			match := strings.Contains(strings.ToLower(c.Description), q) ||
				strings.Contains(strings.ToLower(c.Rationale), q)
			if !match {
				continue
			}
		}
		filtered = append(filtered, c)
	}

	if len(filtered) == 0 {
		return s.Card("", s.EmptyState.Render(fmt.Sprintf("No controls match \"%s\".", ts.searchQuery)), width)
	}

	sectionTitle := fmt.Sprintf("Recommended Controls (%d)", len(r.Controls))
	if ts.filterActive && q != "" {
		sectionTitle = fmt.Sprintf("Recommended Controls (%d of %d matching \"%s\")", len(filtered), len(r.Controls), q)
	}

	var rows []string
	if ts.showHelp {
		rows = append(rows, s.DimText.Render("  Controls are recommended mitigations that address identified risks and assumptions."))
		rows = append(rows, s.DimText.Render("  ↑↓ Select  |  Enter Detail  |  / Search"))
		rows = append(rows, "")
	}

	if r.CoverageOutput != nil {
		var cov []string
		if len(r.CoverageOutput.BlindSpots) > 0 {
			cov = append(cov, fmt.Sprintf("  Blind Spots: %d", len(r.CoverageOutput.BlindSpots)))
		}
		if len(r.CoverageOutput.DomainBlindSpots) > 0 {
			cov = append(cov, fmt.Sprintf("  Domain Blind Spots: %d", len(r.CoverageOutput.DomainBlindSpots)))
		}
		if r.CoverageOutput.Assessment != nil && len(r.CoverageOutput.Assessment.Gaps) > 0 {
			cov = append(cov, fmt.Sprintf("  Coverage Gaps: %d", len(r.CoverageOutput.Assessment.Gaps)))
		}
		if len(cov) > 0 {
			rows = append(rows, s.Card("Coverage Overview", strings.Join(cov, "\n"), width-4))
			rows = append(rows, "")
		}
	}

		rows = append(rows, s.SubSectionTitle.Render(sectionTitle))

	for i, c := range filtered {
		selected := i == ts.selectedIndex
		prefix := "  "
		style := s.Value
		if selected {
			prefix = "▸ "
			style = s.SelectedItem
			ts.selectedLine = len(rows)
		}
		rows = append(rows, style.Render(prefix+c.Description))

		if ts.detailOpen && selected {
			var detail []string
			detail = append(detail, "")
			if c.Rationale != "" {
				detail = append(detail, s.DetailBox.Render("  Rationale: "+c.Rationale))
			}
			if len(c.MitigatedAssumptionIDs) > 0 {
				detail = append(detail, s.DetailBox.Render("  Mitigates: "+strings.Join(c.MitigatedAssumptionIDs, ", ")))
			}
			if c.Category != "" {
				detail = append(detail, s.DetailBox.Render("  Category: "+c.Category))
			}
			rows = append(rows, detail...)
		}
	}

	return strings.Join(rows, "\n")
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

func renderResultSDRI(s StyleSet, r *AnalysisResult, ts *tabState, width int) string {
	if r.SDRISummary == "" && len(r.SDRIControls) == 0 &&
		len(r.SDRIDesignFindings) == 0 && len(r.SDRIAchitecturalWeaknesses) == 0 &&
		len(r.SDRIRemediations) == 0 && len(r.SDRIComplianceAlignments) == 0 {
		return s.Card("", s.EmptyState.Render("No SDRI data available."), width)
	}

	type sdriItem struct {
		section string
		label   string
		detail  string
	}

	var items []sdriItem
	q := strings.ToLower(ts.searchQuery)

	if r.SDRISummary != "" {
		items = append(items, sdriItem{
			section: "summary",
			label:   "Executive Summary",
			detail:  r.SDRISummary,
		})
	}

	for _, c := range r.SDRIControls {
		status := c.Status
		if status == "" {
			status = "unknown"
		}
		label := fmt.Sprintf("Control: %s (%s)", c.Name, status)
		if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(label), q) {
			continue
		}
		items = append(items, sdriItem{
			section: "controls",
			label:   label,
			detail:  fmt.Sprintf("Status: %s | Category: %s", c.Status, c.Category),
		})
	}

	for _, f := range r.SDRIDesignFindings {
		label := fmt.Sprintf("Finding: %s [%s]", f.Title, f.Severity)
		if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(label), q) && !strings.Contains(strings.ToLower(f.Description), q) {
			continue
		}
		detail := f.Description
		if f.Recommendation != "" {
			detail += " → " + f.Recommendation
		}
		items = append(items, sdriItem{
			section: "findings",
			label:   label,
			detail:  detail,
		})
	}

	for _, w := range r.SDRIAchitecturalWeaknesses {
		label := fmt.Sprintf("Weakness: %s [%s]", w.Pattern, w.Severity)
		if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(label), q) && !strings.Contains(strings.ToLower(w.Description), q) {
			continue
		}
		detail := w.Description
		if w.Recommendation != "" {
			detail += " → " + w.Recommendation
		}
		items = append(items, sdriItem{
			section: "weaknesses",
			label:   label,
			detail:  detail,
		})
	}

	for _, rem := range r.SDRIRemediations {
		label := fmt.Sprintf("Remediation #%d [%.0f]", rem.Priority, rem.RiskScore)
		if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(rem.Description), q) {
			continue
		}
		detail := rem.Description
		if rem.Recommendation != "" {
			detail += " → " + rem.Recommendation + " (effort: " + rem.Effort + ")"
		}
		items = append(items, sdriItem{
			section: "remediations",
			label:   label,
			detail:  detail,
		})
	}

	for _, m := range r.SDRIComplianceAlignments {
		label := fmt.Sprintf("Compliance: %s (%.0f%% %s)", m.Framework, m.Coverage, m.Status)
		if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(label), q) {
			continue
		}
		items = append(items, sdriItem{
			section: "compliance",
			label:   label,
			detail:  fmt.Sprintf("Framework: %s | Coverage: %.0f%% | Status: %s", m.Framework, m.Coverage, m.Status),
		})
	}

	if len(items) == 0 {
		return s.Card("", s.EmptyState.Render("No SDRI data matches filter."), width)
	}

	var rows []string
	if ts.showHelp {
		rows = append(rows, s.DimText.Render("  SDRI (Security Design Review Intelligence) evaluates design findings, weaknesses, and remediations."))
		rows = append(rows, s.DimText.Render("  ↑↓ Select  |  Enter Detail  |  / Search"))
		rows = append(rows, "")
	}
	for i, it := range items {
		selected := i == ts.selectedIndex
		prefix := "  "
		style := s.Value
		if selected {
			prefix = "▸ "
			style = s.SelectedItem
			ts.selectedLine = 1 + len(rows)
		}
		rows = append(rows, style.Render(prefix+it.label))
		if ts.detailOpen && selected {
			rows = append(rows, "")
			rows = append(rows, s.DetailBox.Render("  "+it.detail))
			rows = append(rows, "")
		}
	}

	sectionTitle := fmt.Sprintf("SDRI (%d)", len(items))
	if ts.filterActive && q != "" {
		sectionTitle = fmt.Sprintf("SDRI (%d matching \"%s\")", len(items), q)
	}
	return s.SubSectionTitle.Render(sectionTitle) + "\n" + strings.Join(rows, "\n")
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
