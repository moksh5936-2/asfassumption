package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type tabState struct {
	selectedIndex  int
	detailOpen     bool
	searchQuery    string
	SearchActive   bool
	filterActive   bool
	showHelp       bool
	selectedLine   int
	contentOffset  int
	ViewportOffset int
	DetailOffset   int
	filteredCount  int
}

type resultsModel struct {
	result      *AnalysisResult
	resultTab   int
	tabs        []resultTabDef
	tabScroll   map[int]int
	tabStates   map[int]*tabState
	detailFocus bool
	vpReady     bool
	vpWidth     int
	vpHeight    int
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
		tabScroll:   make(map[int]int),
		tabStates:   make(map[int]*tabState),
		detailFocus: false,
	}
}

func (m resultsModel) Update(msg tea.Msg) (resultsModel, tea.Cmd) {
	return m, nil
}

func ensureSelectedVisible(ts *tabState, itemCount int, visibleHeight int) {
	if itemCount <= 0 {
		ts.selectedIndex = 0
		ts.ViewportOffset = 0
		return
	}
	if ts.selectedIndex < 0 {
		ts.selectedIndex = 0
	}
	if ts.selectedIndex >= itemCount {
		ts.selectedIndex = itemCount - 1
	}
	if visibleHeight < 1 {
		visibleHeight = 1
	}
	if ts.selectedIndex < ts.ViewportOffset {
		ts.ViewportOffset = ts.selectedIndex
	}
	if ts.selectedIndex >= ts.ViewportOffset+visibleHeight {
		ts.ViewportOffset = ts.selectedIndex - visibleHeight + 1
	}
	maxOffset := itemCount - visibleHeight
	if maxOffset < 0 {
		maxOffset = 0
	}
	if ts.ViewportOffset > maxOffset {
		ts.ViewportOffset = maxOffset
	}
	if ts.ViewportOffset < 0 {
		ts.ViewportOffset = 0
	}
}

func (m mainModel) viewResults() string {
	s := m.styles

	if m.results.result == nil {
		return lipgloss.JoinVertical(lipgloss.Left,
			s.PremiumHeader("Case Workspace", m.mainWidth()),
			"",
			s.Card("",
				lipgloss.JoinVertical(lipgloss.Left,
					s.SubSectionTitle.Render("  No results yet."),
					"",
					s.DimText.Render("  Start a new analysis or open an existing case to see results here."),
					"",
					s.DimText.Render("  Press  n for New Analysis"),
					s.DimText.Render("  Press  ? for Quick Tour"),
				),
				m.mainWidth()-4,
			),
		)
	}

	r := m.results.result
	tab := m.results.resultTab
	ts := m.results.tabStateFor(tab)
	totalWidth := m.mainWidth() - 4

	tabBar := m.renderResultTabs()
	sectionName := m.results.tabs[tab].name
	countStr := m.results.tabCountString(tab)
	titleStr := sectionName
	if countStr != "" {
		titleStr = sectionName + " — " + countStr
	}
	header := s.PremiumHeader(titleStr, m.mainWidth())
	ts.contentOffset = strings.Count(lipgloss.JoinVertical(lipgloss.Left, header, tabBar), "\n") + 1

	if tab == 0 {
		content := renderResultSummary(s, r, ts, totalWidth)
		return lipgloss.JoinVertical(lipgloss.Left,
			header,
			tabBar,
			content,
		)
	}

	// Split pane for tabs 1-6
	paneHeight := m.mainHeight() - ts.contentOffset - 1
	if paneHeight < 10 {
		paneHeight = 10
	}

	const minSplitWidth = 100
	listWidth := totalWidth * 40 / 100
	detailWidth := totalWidth - listWidth - 3
	if totalWidth < minSplitWidth {
		listWidth = totalWidth
		detailWidth = 0
	}
	if listWidth < 30 {
		listWidth = 30
	}
	if detailWidth < 30 && totalWidth >= minSplitWidth {
		detailWidth = 30
		listWidth = totalWidth - detailWidth - 3
	}

	var listContent string
	switch tab {
	case 1:
		listContent = renderAssumptionsList(s, r, ts, listWidth)
	case 2:
		listContent = renderVerificationList(s, r, ts, listWidth)
	case 3:
		listContent = renderContradictionsList(s, r, ts, listWidth)
	case 4:
		listContent = renderTrustList(s, r, ts, listWidth)
	case 5:
		listContent = renderControlsList(s, r, ts, listWidth)
	case 6:
		listContent = renderSDRISummaryList(s, r, ts, listWidth)
	}

	listLines := strings.Split(listContent, "\n")
	totalListLines := len(listLines)

	// Count header lines (everything before the first item line)
	headerLines := totalListLines - ts.filteredCount
	if headerLines < 0 {
		headerLines = 0
	}
	// Item visible height = paneHeight minus header area
	itemVisHeight := paneHeight - headerLines
	if itemVisHeight < 1 {
		itemVisHeight = 1
	}
	// Clamp viewport in item-space using correct visible height
	ensureSelectedVisible(ts, ts.filteredCount, itemVisHeight)
	// Convert item-space ViewportOffset to line-space for slicing
	lineOffset := headerLines + ts.ViewportOffset
	if lineOffset < 0 {
		lineOffset = 0
	}
	if lineOffset >= totalListLines && totalListLines > 0 {
		lineOffset = max(0, totalListLines-1)
	}
	visFrom := lineOffset
	visTo := visFrom + paneHeight
	if visTo > totalListLines {
		visTo = totalListLines
	}
	if visFrom < totalListLines {
		listLines = listLines[visFrom:visTo]
	} else {
		listLines = nil
	}
	for len(listLines) < paneHeight {
		listLines = append(listLines, "")
	}
	renderedList := strings.Join(listLines, "\n")

	listBorder := s.DetailBox
	if !m.results.detailFocus {
		listBorder = s.SelectedItem
	}
	listPane := listBorder.Render(renderedList)

	if totalWidth < minSplitWidth {
		return lipgloss.JoinVertical(lipgloss.Left,
			header,
			tabBar,
			listPane,
		)
	}

	var detailContent string
	switch tab {
	case 1:
		detailContent = renderAssumptionsDetail(s, r, ts, detailWidth)
	case 2:
		detailContent = renderVerificationDetail(s, r, ts, detailWidth)
	case 3:
		detailContent = renderContradictionsDetail(s, r, ts, detailWidth)
	case 4:
		detailContent = renderTrustDetail(s, r, ts, detailWidth)
	case 5:
		detailContent = renderControlsDetail(s, r, ts, detailWidth)
	case 6:
		detailContent = renderSDRIDetail(s, r, ts, detailWidth)
	}

	detailLines := strings.Split(detailContent, "\n")
	totalDetailLines := len(detailLines)
	if ts.DetailOffset >= totalDetailLines && ts.DetailOffset > 0 {
		ts.DetailOffset = max(0, totalDetailLines-1)
	}

	detailLines2 := detailLines
	if ts.DetailOffset < len(detailLines2) {
		detailLines2 = detailLines2[ts.DetailOffset:]
	}
	if len(detailLines2) > paneHeight {
		detailLines2 = detailLines2[:paneHeight]
	}
	for len(detailLines2) < paneHeight {
		detailLines2 = append(detailLines2, "")
	}
	renderedDetail := strings.Join(detailLines2, "\n")

	detailBorder := s.DetailBox
	if m.results.detailFocus {
		detailBorder = s.SelectedItem
	}
	detailPane := detailBorder.Render(renderedDetail)

	splitContent := lipgloss.JoinHorizontal(lipgloss.Top,
		listPane,
		s.DimText.Render(" │ "),
		detailPane,
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		tabBar,
		splitContent,
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
	if m.results.detailFocus {
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

func (m mainModel) paneHeight() int {
	h := m.mainHeight() - 4 // subtract header + tab bar
	if h < 10 {
		h = 10
	}
	return h
}

func (m mainModel) updateResults(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		tab := m.results.resultTab
		if tab > 0 && m.router.focus == focusContent {
			ts := m.results.tabStateFor(tab)
			switch msg.Type {
			case tea.MouseWheelUp:
				if m.results.detailFocus {
					if ts.DetailOffset > 0 {
						ts.DetailOffset--
					}
				} else {
					if ts.selectedIndex > 0 {
						ts.selectedIndex--
					}
					ensureSelectedVisible(ts, m.results.tabCount(tab), m.paneHeight())
				}
				return m, nil
			case tea.MouseWheelDown:
				if m.results.detailFocus {
					ts.DetailOffset++
				} else if ts.selectedIndex < m.results.tabCount(tab)-1 {
					ts.selectedIndex++
					ensureSelectedVisible(ts, m.results.tabCount(tab), m.paneHeight())
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

			// Search input mode
			if ts.filterActive {
				switch msg.String() {
				case "esc":
					ts.filterActive = false
					ts.searchQuery = ""
					if ts.selectedIndex >= m.results.tabCount(tab) {
						ts.selectedIndex = 0
					}
					return m, nil
				case "backspace":
					if len(ts.searchQuery) > 0 {
						ts.searchQuery = ts.searchQuery[:len(ts.searchQuery)-1]
					}
					return m, nil
				case "n":
					if ts.searchQuery != "" && ts.selectedIndex < maxIdx {
						ts.selectedIndex++
						ensureSelectedVisible(ts, m.results.tabCount(tab), m.paneHeight())
					}
					return m, nil
				case "N":
					if ts.searchQuery != "" && ts.selectedIndex > 0 {
						ts.selectedIndex--
						ensureSelectedVisible(ts, m.results.tabCount(tab), m.paneHeight())
					}
					return m, nil
				default:
					if len(msg.String()) == 1 && msg.String() != "\t" {
						ts.searchQuery += msg.String()
						return m, nil
					}
				}
				return m, nil
			}

			switch msg.String() {
			case "up", "k":
				if m.results.detailFocus {
					if ts.DetailOffset > 0 {
						ts.DetailOffset--
					}
				} else {
					if ts.selectedIndex > 0 {
						ts.selectedIndex--
					}
					ensureSelectedVisible(ts, m.results.tabCount(tab), m.paneHeight())
				}
				return m, nil
			case "down", "j":
				if m.results.detailFocus {
					ts.DetailOffset++
				} else {
					if ts.selectedIndex < maxIdx {
						ts.selectedIndex++
					}
					ensureSelectedVisible(ts, m.results.tabCount(tab), m.paneHeight())
				}
				return m, nil
			case "pgup":
				if m.results.detailFocus {
					ts.DetailOffset -= m.paneHeight()
					if ts.DetailOffset < 0 {
						ts.DetailOffset = 0
					}
				} else {
					ts.selectedIndex -= m.paneHeight()
					if ts.selectedIndex < 0 {
						ts.selectedIndex = 0
					}
					ensureSelectedVisible(ts, m.results.tabCount(tab), m.paneHeight())
				}
				return m, nil
			case "pgdown":
				if m.results.detailFocus {
					ts.DetailOffset += m.paneHeight()
				} else {
					ts.selectedIndex += m.paneHeight()
					if ts.selectedIndex > maxIdx {
						ts.selectedIndex = maxIdx
					}
					ensureSelectedVisible(ts, m.results.tabCount(tab), m.paneHeight())
				}
				return m, nil
			case "home":
				if m.results.detailFocus {
					ts.DetailOffset = 0
				} else {
					ts.selectedIndex = 0
					ensureSelectedVisible(ts, m.results.tabCount(tab), m.paneHeight())
				}
				return m, nil
			case "end":
				if m.results.detailFocus {
					ts.DetailOffset = 1 << 30
				} else {
					ts.selectedIndex = maxIdx
					ensureSelectedVisible(ts, m.results.tabCount(tab), m.paneHeight())
				}
				return m, nil
			case "enter":
				if !m.results.detailFocus && maxIdx >= 0 {
					m.results.detailFocus = true
					ts.DetailOffset = 0
				}
				return m, nil
			case "/":
				m.results.detailFocus = false
				ts.filterActive = true
				ts.searchQuery = ""
				return m, nil
			case "esc":
				if ts.showHelp {
					ts.showHelp = false
					return m, nil
				}
				if m.results.detailFocus {
					m.results.detailFocus = false
					return m, nil
				}
			default:
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
		m.results.detailFocus = false
		ts := m.results.tabStateFor(m.results.resultTab)
		ensureSelectedVisible(ts, m.results.tabCount(m.results.resultTab), m.paneHeight())
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

	var topContradictionsCard string
	if len(r.Contradictions) > 0 {
		var cRows []string
		maxC := 3
		if maxC > len(r.Contradictions) {
			maxC = len(r.Contradictions)
		}
		for i := 0; i < maxC; i++ {
			c := r.Contradictions[i]
			sev := riskStyle(s, c.Severity).Render(string(c.Severity))
			desc := c.Description
			if len(desc) > 60 {
				desc = desc[:57] + "..."
			}
			cRows = append(cRows, fmt.Sprintf("  %s  %s", sev, s.Value.Render(desc)))
		}
		if len(r.Contradictions) > maxC {
			cRows = append(cRows, fmt.Sprintf("  %s  (+%d more — see Tab 3)",
				s.DimText.Render("⋯"),
				len(r.Contradictions)-maxC))
		}
		topContradictionsCard = s.CardAccent("Critical Contradictions (Top)", strings.Join(cRows, "\n"), width)
	}

	var topSPOFsCard string
	if r.TrustOutput != nil && len(r.TrustOutput.SinglePointsOfTrust) > 0 {
		var sRows []string
		maxS := 5
		if maxS > len(r.TrustOutput.SinglePointsOfTrust) {
			maxS = len(r.TrustOutput.SinglePointsOfTrust)
		}
		for i := 0; i < maxS; i++ {
			spof := r.TrustOutput.SinglePointsOfTrust[i]
			sRows = append(sRows, fmt.Sprintf("  %s  %s  (%d dependents)",
				s.StatusBad.Render("SPOF"),
				s.Value.Render(spof.NodeID),
				spof.DependentsCount))
		}
		if len(r.TrustOutput.SinglePointsOfTrust) > maxS {
			sRows = append(sRows, fmt.Sprintf("  %s  (+%d more — see Tab 4)",
				s.DimText.Render("⋯"),
				len(r.TrustOutput.SinglePointsOfTrust)-maxS))
		}
		topSPOFsCard = s.Card("Single Points of Trust Failure", strings.Join(sRows, "\n"), width)
	}

	var sdriFindingsCard string
	if len(r.SDRIDesignFindings) > 0 || len(r.SDRIAchitecturalWeaknesses) > 0 {
		var fRows []string
		count := 0
		maxF := 3
		for i := 0; i < len(r.SDRIDesignFindings) && count < maxF; i++ {
			f := r.SDRIDesignFindings[i]
			if f.Severity == "Critical" || f.Severity == "High" {
				sev := s.StatusBad.Render(f.Severity)
				desc := f.Title
				if len(desc) > 55 {
					desc = desc[:52] + "..."
				}
				fRows = append(fRows, fmt.Sprintf("  %s  %s", sev, s.Value.Render(desc)))
				count++
			}
		}
		for i := 0; i < len(r.SDRIAchitecturalWeaknesses) && count < maxF; i++ {
			w := r.SDRIAchitecturalWeaknesses[i]
			if w.Severity == "Critical" || w.Severity == "High" {
				sev := s.StatusWarn.Render(w.Severity)
				desc := w.Pattern
				if len(desc) > 55 {
					desc = desc[:52] + "..."
				}
				fRows = append(fRows, fmt.Sprintf("  %s  %s", sev, s.Value.Render(desc)))
				count++
			}
		}
		if len(fRows) > 0 {
			totalSDRI := len(r.SDRIDesignFindings) + len(r.SDRIAchitecturalWeaknesses)
			if totalSDRI > count {
				fRows = append(fRows, fmt.Sprintf("  %s  (+%d more — see Tab 6)",
					s.DimText.Render("⋯"),
					totalSDRI-count))
			}
			sdriFindingsCard = s.Card("SDRI — Critical Findings", strings.Join(fRows, "\n"), width)
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
		topContradictionsCard,
		topSPOFsCard,
		sdriFindingsCard,
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

func renderAssumptionsList(s StyleSet, r *AnalysisResult, ts *tabState, width int) string {
	if len(r.Assumptions) == 0 {
		return s.EmptyState.Render("No assumptions found.")
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
	if len(filtered) == 0 {
		return s.EmptyState.Render(fmt.Sprintf("No assumptions match \"%s\".", ts.searchQuery))
	}
	var rows []string
	if ts.showHelp {
		rows = append(rows, s.DimText.Render("  ↑↓ Select  |  Enter Detail  |  r Review  |  v Validate  |  / Search"))
		rows = append(rows, "")
	}
	rows = append(rows, s.DimText.Render(fmt.Sprintf("  %d of %d assumptions", len(filtered), len(r.Assumptions))))
	rows = append(rows, s.SectionRule.Render(strings.Repeat("─", max(1, width-4))))
	rows = append(rows, s.DimText.Render("  ID   Risk      Conf %  Description"))
	ts.filteredCount = len(filtered)
	for i, a := range filtered {
		sel := i == ts.selectedIndex
		prefix := "  "
		style := s.Value
		if sel {
			prefix = "▸ "
			style = s.SelectedItem
		}
		riskSt := riskStyle(s, a.Risk)
		confPct := int(a.Confidence * 100)
		confSt := confidenceStyle(s, confPct)
		rows = append(rows, style.Render(fmt.Sprintf("%s%s %s %s  %s",
			prefix, a.ID,
			riskSt.Render(padRight(string(a.Risk), 10)),
			confSt.Render(fmt.Sprintf("%3d%%", confPct)),
			a.Description)))
	}
	return strings.Join(rows, "\n")
}

func renderAssumptionsDetail(s StyleSet, r *AnalysisResult, ts *tabState, width int) string {
	if len(r.Assumptions) == 0 {
		return ""
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
	if ts.selectedIndex < 0 || ts.selectedIndex >= len(filtered) {
		return s.DimText.Render("No item selected")
	}
	a := filtered[ts.selectedIndex]
	var detail []string
	detail = append(detail, s.SubSectionTitle.Render(fmt.Sprintf("  %s — Detail", a.ID)))
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
	return strings.Join(detail, "\n")
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

func renderVerificationList(s StyleSet, r *AnalysisResult, ts *tabState, width int) string {
	if r.VerificationOutput == nil || r.VerificationOutput.Assessment == nil {
		return s.EmptyState.Render("No verification data.")
	}
	a := r.VerificationOutput.Assessment
	var rows []string
	if ts.showHelp {
		rows = append(rows, s.DimText.Render("  ↑↓ Select  |  Enter Detail  |  / Search"))
		rows = append(rows, "")
	}
	cardW := (width - 12) / 4
	if cardW < 10 {
		cardW = 10
	}
	rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top,
		s.StatusCardLarge("VERIFIED", a.VerifiedCount, "", cardW),
		"  ",
		s.StatusCardLarge("PARTIAL", a.PartialCount, "", cardW),
		"  ",
		s.StatusCardLarge("UNVERIFIED", a.UnverifiedCount, "", cardW),
		"  ",
		s.StatusCardLarge("NO EVIDENCE", a.NoEvidenceCount, "", cardW),
	))
	if a.OverallConfidence > 0 {
		rows = append(rows, "")
		rows = append(rows, s.SubHeader("Overall Confidence"))
		rows = append(rows, "  "+s.ProgressWithLabel(a.OverallConfidence*100, width-8))
	}
	if r.VerificationOutput.CISOView == nil {
		return strings.Join(rows, "\n")
	}
	cv := r.VerificationOutput.CISOView
	q := strings.ToLower(ts.searchQuery)
	type vItem struct {
		label           string
		howToValidate   string
		evidenceMissing string
		group           string
		isWarn          bool
	}
	var filtered []vItem
	for _, plan := range cv.TopAssumptionsToVerify {
		label := plan.AssumptionText
		if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(label), q) {
			continue
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
		filtered = append(filtered, vItem{label: label, howToValidate: plan.HowToValidate, evidenceMissing: emStr, group: group})
	}
	for _, g := range cv.EvidenceGaps {
		label := g
		if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(label), q) {
			continue
		}
		filtered = append(filtered, vItem{label: "⚠ " + label, group: "Evidence Gaps", isWarn: true})
	}
	if len(filtered) == 0 {
		return strings.Join(rows, "\n")
	}
	ts.filteredCount = len(filtered)
	prevGroup := ""
	for i, item := range filtered {
		sel := i == ts.selectedIndex
		prefix := "  • "
		style := s.Value
		if sel {
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
		if sel {
			ts.selectedLine = len(rows)
		}
		rows = append(rows, style.Render(prefix+item.label))
	}
	return strings.Join(rows, "\n")
}

func renderVerificationDetail(s StyleSet, r *AnalysisResult, ts *tabState, width int) string {
	if r.VerificationOutput == nil || r.VerificationOutput.CISOView == nil {
		return s.DimText.Render("No detail available.")
	}
	cv := r.VerificationOutput.CISOView
	q := strings.ToLower(ts.searchQuery)
	type vItem struct {
		label           string
		howToValidate   string
		evidenceMissing string
		group           string
		isWarn          bool
	}
	var filtered []vItem
	for _, plan := range cv.TopAssumptionsToVerify {
		label := plan.AssumptionText
		if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(label), q) {
			continue
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
		filtered = append(filtered, vItem{label: label, howToValidate: plan.HowToValidate, evidenceMissing: emStr, group: group})
	}
	for _, g := range cv.EvidenceGaps {
		label := g
		if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(label), q) {
			continue
		}
		filtered = append(filtered, vItem{label: "⚠ " + label, group: "Evidence Gaps", isWarn: true})
	}
	if ts.selectedIndex < 0 || ts.selectedIndex >= len(filtered) {
		return s.DimText.Render("No item selected")
	}
	item := filtered[ts.selectedIndex]
	var detail []string
	detail = append(detail, s.SubSectionTitle.Render("  Detail"))
	detail = append(detail, "")
	if item.isWarn {
		detail = append(detail, s.DetailBox.Render("  Evidence Gap — requires additional evidence collection."))
	} else {
		detail = append(detail, s.DetailBox.Render(fmt.Sprintf("  Assumption: %s", item.label)))
		if item.howToValidate != "" {
			detail = append(detail, "")
			detail = append(detail, s.DetailBox.Render(fmt.Sprintf("  How to Validate: %s", item.howToValidate)))
		}
		if item.evidenceMissing != "" {
			detail = append(detail, "")
			detail = append(detail, s.DetailBox.Render(fmt.Sprintf("  Evidence Missing: %s", item.evidenceMissing)))
		}
	}
	return strings.Join(detail, "\n")
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

func renderContradictionsList(s StyleSet, r *AnalysisResult, ts *tabState, width int) string {
	if len(r.Contradictions) == 0 {
		return s.EmptyState.Render("No contradictions detected.")
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
	if len(filtered) == 0 {
		return s.EmptyState.Render(fmt.Sprintf("No contradictions match \"%s\".", ts.searchQuery))
	}
	ts.filteredCount = len(filtered)
	var rows []string
	if ts.showHelp {
		rows = append(rows, s.DimText.Render("  ↑↓ Select  |  Enter Detail  |  / Search"))
		rows = append(rows, "")
	}
	rows = append(rows, s.DimText.Render(fmt.Sprintf("  %d of %d contradictions", len(filtered), len(r.Contradictions))))
	for i, c := range filtered {
		sel := i == ts.selectedIndex
		prefix := "  "
		style := s.Value
		if sel {
			prefix = "▸ "
			style = s.SelectedItem
		}
		sevStyle := s.StatusWarn
		if c.Severity == RiskHigh || c.Severity == RiskCritical {
			sevStyle = s.StatusBad
		} else if c.Severity == RiskLow {
			sevStyle = s.StatusGood
		}
		if sel {
			ts.selectedLine = len(rows)
		}
		rows = append(rows, style.Render(fmt.Sprintf("%s%s %s",
			prefix, sevStyle.Render("["+strings.ToUpper(string(c.Severity))+"]"), c.Description)))
	}
	return strings.Join(rows, "\n")
}

func renderContradictionsDetail(s StyleSet, r *AnalysisResult, ts *tabState, width int) string {
	if len(r.Contradictions) == 0 {
		return ""
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
	if ts.selectedIndex < 0 || ts.selectedIndex >= len(filtered) {
		return s.DimText.Render("No item selected")
	}
	c := filtered[ts.selectedIndex]
	var detail []string
	detail = append(detail, s.SubSectionTitle.Render(fmt.Sprintf("  [%s] — Detail", strings.ToUpper(string(c.Severity)))))
	detail = append(detail, "")
	detail = append(detail, s.DetailBox.Render(fmt.Sprintf("  Description: %s", c.Description)))
	if c.RuleName != "" {
		detail = append(detail, s.DetailBox.Render(fmt.Sprintf("  Rule: %s", c.RuleName)))
	}
	if c.Explanation != "" {
		detail = append(detail, s.DetailBox.Render(fmt.Sprintf("  Reason: %s", c.Explanation)))
	}
	if len(c.AffectedAssumptions) > 0 {
		detail = append(detail, s.DetailBox.Render(fmt.Sprintf("  Affects: %s", strings.Join(c.AffectedAssumptions, ", "))))
	}
	if len(c.Evidence) > 0 {
		for _, ev := range c.Evidence {
			detail = append(detail, s.DetailBox.Render(fmt.Sprintf("  Evidence: %s", ev)))
		}
	}
	return strings.Join(detail, "\n")
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

func renderTrustList(s StyleSet, r *AnalysisResult, ts *tabState, width int) string {
	if r.TrustOutput == nil {
		return s.EmptyState.Render("No trust chain data.")
	}
	type tItem struct {
		section string
		label   string
	}
	var items []tItem
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
		items = append(items, tItem{section: "chain",
			label: fmt.Sprintf("Chain %s | Length: %d | %s | %.0f%%",
				chain.ID, chain.Length, chain.Risk, chain.Confidence*100)})
	}
	if r.TrustOutput != nil {
		for _, spof := range r.TrustOutput.SinglePointsOfTrust {
			if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(spof.AssumptionText), q) {
				continue
			}
			items = append(items, tItem{section: "spof", label: "⚠ SPOF: " + spof.AssumptionText})
		}
		for _, fc := range r.TrustOutput.FailureCascades {
			if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(fc.RootAssumptionText), q) {
				continue
			}
			items = append(items, tItem{section: "cascade",
				label: fmt.Sprintf("Cascade: %s | Steps: %d | %s", fc.RootAssumptionText, len(fc.Steps), fc.Severity)})
		}
	}
	if r.ReviewOutput != nil && r.ReviewOutput.Queue != nil {
		for _, item := range r.ReviewOutput.Queue.Items {
			if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(item.AssumptionText), q) {
				continue
			}
			items = append(items, tItem{section: "queue",
				label: fmt.Sprintf("#%d [%.0f] %s", item.Rank, item.PriorityScore, item.AssumptionText)})
		}
	}
	if len(items) == 0 {
		return s.EmptyState.Render("No trust data.")
	}
	ts.filteredCount = len(items)
	var rows []string
	if ts.showHelp {
		rows = append(rows, s.DimText.Render("  ↑↓ Select  |  Enter Detail  |  / Search"))
		rows = append(rows, "")
	}
	rows = append(rows, s.DimText.Render(fmt.Sprintf("  %d trust items", len(items))))
	for i, it := range items {
		sel := i == ts.selectedIndex
		prefix := "  "
		style := s.Value
		if sel {
			prefix = "▸ "
			style = s.SelectedItem
		}
		tag := ""
		switch it.section {
		case "chain":
			tag = "[Chain] "
		case "spof":
			tag = "[SPOF] "
		case "cascade":
			tag = "[Cascade] "
		case "queue":
			tag = "[Queue] "
		}
		if sel {
			ts.selectedLine = len(rows)
		}
		rows = append(rows, style.Render(prefix+tag+it.label))
	}
	return strings.Join(rows, "\n")
}

func renderTrustDetail(s StyleSet, r *AnalysisResult, ts *tabState, width int) string {
	if r.TrustOutput == nil {
		return ""
	}
	type tItem struct {
		section string
		label   string
		detail  string
	}
	var items []tItem
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
		var nodes []string
		if len(chain.Nodes) > 0 {
			nodes = chain.Nodes
		} else {
			nodes = []string{chain.RootNode, chain.LeafNode}
		}
		d := fmt.Sprintf("Nodes: %s", strings.Join(nodes, " → "))
		if chain.Risk != "" {
			d += fmt.Sprintf("\nRisk: %s", chain.Risk)
		}
		items = append(items, tItem{section: "chain",
			label: fmt.Sprintf("Chain %s | Length: %d | %s | %.0f%%",
				chain.ID, chain.Length, chain.Risk, chain.Confidence*100),
			detail: d})
	}
	if r.TrustOutput != nil {
		for _, spof := range r.TrustOutput.SinglePointsOfTrust {
			if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(spof.AssumptionText), q) {
				continue
			}
			var dl []string
			dl = append(dl, "Single Point of Trust Failure")
			if spof.DependentsCount > 0 {
				dl = append(dl, fmt.Sprintf("Dependents: %d", spof.DependentsCount))
			}
			if len(spof.DependentNodes) > 0 {
				dl = append(dl, "Affected: "+strings.Join(spof.DependentNodes, ", "))
			}
			if spof.Recommendation != "" {
				dl = append(dl, "Recommendation: "+spof.Recommendation)
			}
			items = append(items, tItem{section: "spof",
				label: "⚠ SPOF: " + spof.AssumptionText, detail: strings.Join(dl, "\n")})
		}
		for _, fc := range r.TrustOutput.FailureCascades {
			if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(fc.RootAssumptionText), q) {
				continue
			}
			d := fmt.Sprintf("Root: %s (%s)\nSeverity: %s\nTotal Affected: %d\nMax Depth: %d",
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
				d += "\nCascade Chain:\n" + strings.Join(steps, "\n")
			}
			items = append(items, tItem{section: "cascade",
				label:  fmt.Sprintf("Cascade: %s | Steps: %d | %s", fc.RootAssumptionText, len(fc.Steps), fc.Severity),
				detail: d})
		}
	}
	if r.ReviewOutput != nil && r.ReviewOutput.Queue != nil {
		for _, item := range r.ReviewOutput.Queue.Items {
			if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(item.AssumptionText), q) {
				continue
			}
			items = append(items, tItem{section: "queue",
				label:  fmt.Sprintf("#%d [%.0f] %s", item.Rank, item.PriorityScore, item.AssumptionText),
				detail: fmt.Sprintf("Risk: %s | Priority: %.0f", item.Risk, item.PriorityScore)})
		}
	}
	if ts.selectedIndex < 0 || ts.selectedIndex >= len(items) {
		return s.DimText.Render("No item selected")
	}
	it := items[ts.selectedIndex]
	var detail []string
	detail = append(detail, s.SubSectionTitle.Render("  Detail"))
	detail = append(detail, "")
	if it.detail != "" {
		detail = append(detail, s.DetailBox.Render("  "+it.detail))
	} else {
		detail = append(detail, s.DimText.Render("  No additional details."))
	}
	return strings.Join(detail, "\n")
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

func renderControlsList(s StyleSet, r *AnalysisResult, ts *tabState, width int) string {
	if len(r.Controls) == 0 {
		return s.EmptyState.Render("No recommended controls.")
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
		return s.EmptyState.Render(fmt.Sprintf("No controls match \"%s\".", ts.searchQuery))
	}
	ts.filteredCount = len(filtered)
	var rows []string
	if ts.showHelp {
		rows = append(rows, s.DimText.Render("  ↑↓ Select  |  Enter Detail  |  / Search"))
		rows = append(rows, "")
	}
	rows = append(rows, s.DimText.Render(fmt.Sprintf("  %d of %d controls", len(filtered), len(r.Controls))))
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
	for i, c := range filtered {
		sel := i == ts.selectedIndex
		prefix := "  "
		style := s.Value
		if sel {
			prefix = "▸ "
			style = s.SelectedItem
			ts.selectedLine = len(rows)
		}
		rows = append(rows, style.Render(prefix+c.Description))
	}
	return strings.Join(rows, "\n")
}

func renderControlsDetail(s StyleSet, r *AnalysisResult, ts *tabState, width int) string {
	if len(r.Controls) == 0 {
		return ""
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
	if ts.selectedIndex < 0 || ts.selectedIndex >= len(filtered) {
		return s.DimText.Render("No item selected")
	}
	c := filtered[ts.selectedIndex]
	var detail []string
	detail = append(detail, s.SubSectionTitle.Render("  Control Detail"))
	detail = append(detail, "")
	if c.Rationale != "" {
		detail = append(detail, s.DetailBox.Render("  Rationale: "+c.Rationale))
	}
	if c.Category != "" {
		detail = append(detail, s.DetailBox.Render("  Category: "+c.Category))
	}
	if len(c.MitigatedAssumptionIDs) > 0 {
		detail = append(detail, s.DetailBox.Render("  Mitigates: "+strings.Join(c.MitigatedAssumptionIDs, ", ")))
	}
	return strings.Join(detail, "\n")
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

func renderSDRISummaryList(s StyleSet, r *AnalysisResult, ts *tabState, width int) string {
	if r.SDRISummary == "" && len(r.SDRIControls) == 0 &&
		len(r.SDRIDesignFindings) == 0 && len(r.SDRIAchitecturalWeaknesses) == 0 &&
		len(r.SDRIRemediations) == 0 && len(r.SDRIComplianceAlignments) == 0 {
		return s.EmptyState.Render("No SDRI data available.")
	}
	type sdriItem struct {
		section string
		label   string
	}
	var items []sdriItem
	q := strings.ToLower(ts.searchQuery)
	if r.SDRISummary != "" {
		items = append(items, sdriItem{section: "summary", label: "Executive Summary"})
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
		items = append(items, sdriItem{section: "controls", label: label})
	}
	for _, f := range r.SDRIDesignFindings {
		label := fmt.Sprintf("[%s] %s", f.Severity, f.Title)
		if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(label), q) && !strings.Contains(strings.ToLower(f.Description), q) {
			continue
		}
		items = append(items, sdriItem{section: "findings", label: label})
	}
	for _, w := range r.SDRIAchitecturalWeaknesses {
		label := fmt.Sprintf("Weakness: %s [%s]", w.Pattern, w.Severity)
		if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(label), q) && !strings.Contains(strings.ToLower(w.Description), q) {
			continue
		}
		items = append(items, sdriItem{section: "weaknesses", label: label})
	}
	for _, rem := range r.SDRIRemediations {
		label := fmt.Sprintf("Remediation #%d [%.0f]", rem.Priority, rem.RiskScore)
		if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(rem.Description), q) {
			continue
		}
		items = append(items, sdriItem{section: "remediations", label: label})
	}
	for _, m := range r.SDRIComplianceAlignments {
		label := fmt.Sprintf("Compliance: %s (%.0f%% %s)", m.Framework, m.Coverage, m.Status)
		if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(label), q) {
			continue
		}
		items = append(items, sdriItem{section: "compliance", label: label})
	}
	if len(items) == 0 {
		return s.EmptyState.Render("No SDRI data matches filter.")
	}
	ts.filteredCount = len(items)
	var rows []string
	if ts.showHelp {
		rows = append(rows, s.DimText.Render("  ↑↓ Select  |  Enter Detail  |  / Search"))
		rows = append(rows, "")
	}
	rows = append(rows, s.DimText.Render(fmt.Sprintf("  %d SDRI items", len(items))))
	for i, it := range items {
		sel := i == ts.selectedIndex
		prefix := "  "
		style := s.Value
		if sel {
			prefix = "▸ "
			style = s.SelectedItem
			ts.selectedLine = len(rows)
		}
		rows = append(rows, style.Render(prefix+it.label))
	}
	return s.SubSectionTitle.Render(fmt.Sprintf("SDRI (%d)", len(items))) + "\n" + strings.Join(rows, "\n")
}

func renderSDRIDetail(s StyleSet, r *AnalysisResult, ts *tabState, width int) string {
	if r.SDRISummary == "" && len(r.SDRIControls) == 0 &&
		len(r.SDRIDesignFindings) == 0 && len(r.SDRIAchitecturalWeaknesses) == 0 &&
		len(r.SDRIRemediations) == 0 && len(r.SDRIComplianceAlignments) == 0 {
		return ""
	}
	type sdriItem struct {
		section string
		label   string
		detail  string
	}
	var items []sdriItem
	q := strings.ToLower(ts.searchQuery)
	if r.SDRISummary != "" {
		items = append(items, sdriItem{section: "summary", label: "Executive Summary", detail: r.SDRISummary})
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
		items = append(items, sdriItem{section: "controls", label: label,
			detail: fmt.Sprintf("Status: %s | Category: %s", c.Status, c.Category)})
	}
	for _, f := range r.SDRIDesignFindings {
		label := fmt.Sprintf("Finding: %s [%s]", f.Title, f.Severity)
		if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(label), q) && !strings.Contains(strings.ToLower(f.Description), q) {
			continue
		}
		d := f.Description
		if f.Recommendation != "" {
			d += " → " + f.Recommendation
		}
		items = append(items, sdriItem{section: "findings", label: label, detail: d})
	}
	for _, w := range r.SDRIAchitecturalWeaknesses {
		label := fmt.Sprintf("Weakness: %s [%s]", w.Pattern, w.Severity)
		if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(label), q) && !strings.Contains(strings.ToLower(w.Description), q) {
			continue
		}
		d := w.Description
		if w.Recommendation != "" {
			d += " → " + w.Recommendation
		}
		items = append(items, sdriItem{section: "weaknesses", label: label, detail: d})
	}
	for _, rem := range r.SDRIRemediations {
		label := fmt.Sprintf("Remediation #%d [%.0f]", rem.Priority, rem.RiskScore)
		if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(rem.Description), q) {
			continue
		}
		d := rem.Description
		if rem.Recommendation != "" {
			d += " → " + rem.Recommendation + " (effort: " + rem.Effort + ")"
		}
		items = append(items, sdriItem{section: "remediations", label: label, detail: d})
	}
	for _, m := range r.SDRIComplianceAlignments {
		label := fmt.Sprintf("Compliance: %s (%.0f%% %s)", m.Framework, m.Coverage, m.Status)
		if ts.filterActive && q != "" && !strings.Contains(strings.ToLower(label), q) {
			continue
		}
		items = append(items, sdriItem{section: "compliance", label: label,
			detail: fmt.Sprintf("Framework: %s | Coverage: %.0f%% | Status: %s", m.Framework, m.Coverage, m.Status)})
	}
	if ts.selectedIndex < 0 || ts.selectedIndex >= len(items) {
		return s.DimText.Render("No item selected")
	}
	it := items[ts.selectedIndex]
	var detail []string
	detail = append(detail, s.SubSectionTitle.Render("  SDRI Detail"))
	detail = append(detail, "")
	detail = append(detail, s.DetailBox.Render("  "+it.detail))
	return strings.Join(detail, "\n")
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
