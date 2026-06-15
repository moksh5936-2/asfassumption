package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type view int

const (
	analyzeView view = iota
	caseView
	reviewView
	validationView
	reportsView
	settingsView
	helpView
	aboutView
	localAIView
	tourView
)

type layoutManager struct {
	sidebarWidth     int
	headerHeight     int
	breadcrumbHeight int
	hintsHeight      int
	statusBarHeight  int
}

func newLayoutManager() layoutManager {
	return layoutManager{
		sidebarWidth:     26,
		headerHeight:     1,
		breadcrumbHeight: 1,
		hintsHeight:      1,
		statusBarHeight:  1,
	}
}

type errorMsg string

func (e errorMsg) Error() string { return string(e) }

type viewHistory struct {
	views []view
}

const maxHistory = 50

func (v *viewHistory) push(vw view) {
	v.views = append(v.views, vw)
	if len(v.views) > maxHistory {
		v.views = v.views[len(v.views)-maxHistory:]
	}
}

func (v *viewHistory) pop() (view, bool) {
	if len(v.views) == 0 {
		return analyzeView, false
	}
	last := v.views[len(v.views)-1]
	v.views = v.views[:len(v.views)-1]
	return last, true
}

type mainModel struct {
	ready           bool
	startup         bool
	startupMenuIdx  int
	width           int
	height          int
	router          Router
	styles          StyleSet
	config          *Config
	engine          *Engine
	quitting        bool
	err             error
	statusMsg       string
	currentFile     string
	searchActive    bool
	searchQuery     string
	recentFiles     []string
	pickerActive    bool
	filePicker      filePickerState
	lastPickerPaths map[pickerMode]string

	analyze  analyzeModel
	results  resultsModel
	settings settingsModel
	about    aboutModel
	reportsV reportsModel
	review   reviewModel
	validate validationModel
	help     helpModel
	localai  localaiModel
	tour     tourModel

	caseResults map[string]*AnalysisResult
	activeCase  string
	caseTab     int

	vp        viewport.Model
	scrollY   map[view]int
	layoutMgr layoutManager
}

type navigateMsg struct {
	to view
}

func newMainModel(cfg *Config) *mainModel {
	theme, ok := Themes[cfg.Appearance.Theme]
	if !ok {
		theme = Themes["ASF0"]
	}
	s := NewStyles(theme)
	e := NewEngine(cfg)
	m := &mainModel{
		startup:         true,
		router:          newRouter(),
		styles:          s,
		config:          cfg,
		engine:          e,
		vp:              viewport.New(0, 0),
		scrollY:         make(map[view]int),
		filePicker:      newFilePickerState(),
		lastPickerPaths: make(map[pickerMode]string),
		analyze:         newAnalyzeModel(e),
		results:         newResultsModel(),
		settings:        newSettingsModel(cfg),
		about:           newAboutModel(),
		reportsV:        newReportsModel(),
		review:          newReviewModel(),
		validate:        newValidationModel(),
		help:            newHelpModel(),
		localai:         newLocalAIModel(cfg),
		tour:            newTourModel(),
		layoutMgr:       newLayoutManager(),
		caseResults:     make(map[string]*AnalysisResult),
	}
	m.router.rebuildCaseEntries(m.getCaseLabels())
	return m
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.startup {
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.ready = true
			m.width = msg.Width
			m.height = msg.Height
			return m, nil
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k":
				if m.startupMenuIdx > 0 {
					m.startupMenuIdx--
				}
				return m, nil
			case "down", "j":
				if m.startupMenuIdx < 3 {
					m.startupMenuIdx++
				}
				return m, nil
			case "enter":
				m.startup = false
				switch m.startupMenuIdx {
				case 0: // Start New Analysis
					m.router.SetView(analyzeView)
				case 1: // Open Existing Case
					if len(m.recentFiles) > 0 {
						m.activeCase = m.recentFiles[0]
						m.results.result = m.caseResults[m.activeCase]
						m.results.resultTab = 0
						m.results.tabStates = make(map[int]*tabState)
						m.router.SetView(caseView)
					} else {
						m.router.SetView(analyzeView)
						m.statusMsg = "No existing cases yet. Start a new analysis."
					}
				case 2: // Local AI Settings
					m.router.SetView(localAIView)
				case 3: // Quick Tour
					m.tour.step = 0
					m.router.SetView(tourView)
				}
				return m, nil
			case "o", "O":
				m.startup = false
				if len(m.recentFiles) > 0 {
					m.activeCase = m.recentFiles[0]
					m.results.result = m.caseResults[m.activeCase]
					m.results.resultTab = 0
					m.results.tabStates = make(map[int]*tabState)
					m.router.SetView(caseView)
				} else {
					m.router.SetView(analyzeView)
					m.statusMsg = "No existing cases yet. Start a new analysis."
				}
				return m, nil
			case "a", "A":
				m.startup = false
				m.router.SetView(localAIView)
				return m, nil
			case "?":
				m.startup = false
				m.tour.step = 0
				m.router.SetView(tourView)
				return m, nil
			case "q":
				// 'q' does nothing on startup screen (only Q quits)
				return m, nil
			case "Q", "ctrl+c":
				m.quitting = true
				return m, tea.Quit
			}
		}
		return m, nil
	}

	if m.pickerActive {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			cmd, _ := m.filePicker.handleKey(msg)
			if cmd != nil {
				m.pickerActive = false
				return m, cmd
			}
			if !m.pickerActive {
				return m, nil
			}
			return m, nil
		case filePickedMsg:
			m.pickerActive = false
			m.handleFilePicked(msg)
			return m, nil
		case filePickerCancelledMsg:
			m.pickerActive = false
			return m, nil
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ready = true
		m.width = msg.Width
		m.height = msg.Height
		theme, ok := Themes[m.config.Appearance.Theme]
		if !ok {
			theme = Themes["ASF0"]
		}
		m.styles = NewStyles(theme)
		m.vp.Width = m.mainWidth()
		m.vp.Height = m.mainHeight()
		if m.results.resultTab > 0 && m.results.result != nil {
			ts := m.results.tabStateFor(m.results.resultTab)
			ensureSelectedVisible(ts, m.results.tabCount(m.results.resultTab), m.paneHeight())
		}

	case tea.KeyMsg:
		if m.searchActive {
			return m.handleSearchInput(msg)
		}

		if m.router.focus == focusSidebar {
			switch msg.String() {
			case "tab", "esc":
				m.router.ToggleFocus()
				return m, nil
			case "up", "k":
				m.router.sidebarMoveUp()
				return m, nil
			case "down", "j":
				m.router.sidebarMoveDown()
				return m, nil
			case "enter":
				nodes := m.router.sidebarVisibleNodes()
				if m.router.sidebarSel >= len(nodes) {
					return m, nil
				}
				n := nodes[m.router.sidebarSel]
				if n.isSection {
					return m, nil
				}
				tab := m.router.sidebarSelTab()
				if m.router.sidebarSelIsParent() {
					if m.router.sidebarSelIsExpanded() {
						m.router.sidebarCollapse()
					} else {
						m.router.sidebarExpand()
					}
					return m, nil
				}
				m.router.sidebarActivate()
				if n.vid == caseView && tab >= 0 && tab < len(m.recentFiles) {
					m.activeCase = m.recentFiles[tab]
					m.results.result = m.caseResults[m.activeCase]
					m.results.resultTab = 0
					m.results.tabStates = make(map[int]*tabState)
				} else if tab >= 0 {
					m.results.resultTab = tab
				}
				m.restoreScroll()
				return m, nil
			case "left":
				if m.router.sidebarSelIsParent() && m.router.sidebarSelIsExpanded() {
					m.router.sidebarCollapse()
					return m, nil
				}
			case "right":
				if m.router.sidebarSelIsParent() && !m.router.sidebarSelIsExpanded() {
					m.router.sidebarExpand()
					return m, nil
				}

			}
			return m, nil
		}

		// Tab navigation for case workspace
		if m.router.currentView == caseView && m.results.result != nil {
			switch msg.String() {
			case "left", "h":
				if m.results.resultTab > 0 {
					m.results.tabScroll[m.results.resultTab] = m.vp.YOffset
					m.results.resultTab--
					if y, ok := m.results.tabScroll[m.results.resultTab]; ok {
						m.vp.YOffset = y
					} else {
						m.vp.YOffset = 0
					}
				}
				return m, nil
			case "right", "l":
				if m.results.resultTab < len(m.results.tabs)-1 {
					m.results.tabScroll[m.results.resultTab] = m.vp.YOffset
					m.results.resultTab++
					if y, ok := m.results.tabScroll[m.results.resultTab]; ok {
						m.vp.YOffset = y
					} else {
						m.vp.YOffset = 0
					}
				}
				return m, nil
			case "up", "k", "down", "j", "enter", "/", "n", "N":
				if m.results.resultTab > 0 && m.router.focus == focusContent {
					return m.updateResults(msg)
				}
			case "?":
				if m.router.focus == focusContent {
					ts := m.results.tabStateFor(m.results.resultTab)
					ts.showHelp = !ts.showHelp
					return m, nil
				}
			case "esc":
				if m.results.resultTab > 0 && m.router.focus == focusContent {
					ts := m.results.tabStateFor(m.results.resultTab)
					if m.results.detailFocus || ts.filterActive || ts.showHelp {
						return m.updateResults(msg)
					}
				}
			}
		}

		if m.router.currentView == tourView {
			switch msg.String() {
			case "right", " ", "l":
				if m.tour.step < m.tour.totalSteps-1 {
					m.tour.step++
				}
				return m, nil
			case "left", "h":
				if m.tour.step > 0 {
					m.tour.step--
				}
				return m, nil
			case "q", "esc":
				m.startup = true
				m.router.SetView(analyzeView)
				return m, nil
			}
		}

		handled, model, cmd := m.handleGlobalKey(msg)
		if handled {
			return model, cmd
		}

	case tea.MouseMsg:
		// List tabs handle mouse wheel for selection
		if m.router.currentView == caseView && m.router.focus == focusContent &&
			m.results.result != nil && m.results.resultTab > 0 {
			return m.updateResults(msg)
		}
		switch msg.Type {
		case tea.MouseWheelUp:
			m.vp, _ = m.vp.Update(msg)
			return m, nil
		case tea.MouseWheelDown:
			m.vp, _ = m.vp.Update(msg)
			return m, nil
		}

	case navigateMsg:
		if msg.to == reportsView {
			m.reportsV.selected = 0
			m.reportsV.done = false
			m.reportsV.exportPath = ""
			m.reportsV.showConfirmation = false
			m.reportsV.err = nil
			m.reportsV.result = m.results.result
			m.reportsV.outputDir = m.config.Output.Directory
			if m.reportsV.outputDir == "" {
				m.reportsV.outputDir = "./reports"
			}
			m.reportsV.format = exportFormatFromConfig(m.config)
		}
		if msg.to == caseView && m.results.result != nil {
			m.results.resultTab = 0
		}
		m.navigateTo(msg.to)
		return m, nil

	case errorMsg:
		m.statusMsg = string(msg)
		switch m.router.currentView {
		case analyzeView:
			m.analyze.running = false
			m.analyze.statusMsg = fmt.Sprintf("Error: %s", string(msg))
		}
		return m, nil

	case analysisCompleteMsg:
		m.analyze.running = false
		m.analyze.result = msg.result
		m.analyze.progress = 100
		docPath := m.analyze.docPath()
		m.currentFile = docPath
		m.caseResults[docPath] = msg.result
		m.activeCase = docPath
		m.results.result = msg.result
		m.results.resultTab = 0
		m.results.tabStates = make(map[int]*tabState)
		m.statusMsg = "Analysis complete"
		m.addRecentFile(docPath)
		m.router.rebuildCaseEntries(m.getCaseLabels())
		m.navigateTo(caseView)
		m.scrollY[caseView] = 0
		return m, nil

	case openFilePickerMsg:
		m.filePicker = newFilePickerState()
		m.filePicker.path = m.pickerStartPath(msg.mode)
		m.filePicker.mode = msg.mode
		m.filePicker.refresh()
		m.pickerActive = true
		return m, nil

	case filePickedMsg:
		m.handleFilePicked(msg)
		return m, nil

	case filePickerCancelledMsg:
		m.pickerActive = false
		return m, nil
	}

	switch m.router.currentView {
	case analyzeView:
		return m.updateAnalyze(msg)
	case caseView:
		return m.updateResults(msg)
	case settingsView:
		return m.updateSettings(msg)
	case aboutView:
		return m.updateAbout(msg)
	case reportsView:
		return m.updateReports(msg)
	case reviewView:
		return m.updateReview(msg)
	case validationView:
		return m.updateValidation(msg)
	case helpView:
		return m.updateHelp(msg)
	case localAIView:
		return m.updateLocalAI(msg)
	}
	return m, nil
}

func (m *mainModel) pickerStartPath(mode pickerMode) string {
	if p, ok := m.lastPickerPaths[mode]; ok && p != "" {
		return p
	}
	switch mode {
	case pickerArchitecture:
		if m.analyze.docPath() != "" {
			return filepath.Dir(m.analyze.docPath())
		}
	case pickerEvidence:
		if m.analyze.evPath() != "" {
			return filepath.Dir(m.analyze.evPath())
		}
	}
	cwd, err := os.Getwd()
	if err == nil {
		return cwd
	}
	home, err := os.UserHomeDir()
	if err == nil {
		return home
	}
	return "."
}

func (m *mainModel) handleFilePicked(msg filePickedMsg) {
	m.pickerActive = false
	m.lastPickerPaths[msg.mode] = filepath.Dir(msg.path)
	m.addRecentFile(msg.path)
	if msg.mode == pickerArchitecture {
		m.analyze.setDocPath(msg.path)
		m.currentFile = msg.path
		m.statusMsg = "Architecture file selected: " + filepath.Base(msg.path)
	} else {
		m.analyze.addEvidence(msg.path)
		m.statusMsg = "Evidence added: " + filepath.Base(msg.path)
	}
}

func (m mainModel) handleGlobalKey(msg tea.KeyMsg) (bool, tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "Q":
		m.quitting = true
		return true, m, tea.Quit
	case "q":
		m.navigateBack()
		return true, m, nil
	case "?":
		if m.router.currentView == tourView {
			return false, m, nil
		}
		m.navigateTo(helpView)
		return true, m, nil
	case "tab":
		m.router.ToggleFocus()
		return true, m, nil
	case "esc":
		switch m.router.currentView {
		case analyzeView:
			if m.analyze.running {
				return false, m, nil
			}
		case settingsView:
			if m.settings.editing {
				return false, m, nil
			}
		case reportsView:
			if m.reportsV.showConfirmation || m.reportsV.done {
				return false, m, nil
			}
		case reviewView:
			if m.review.editing {
				return false, m, nil
			}
		}
		m.navigateBack()
		return true, m, nil
	case "up", "k":
		switch m.router.currentView {
		case caseView, helpView, aboutView:
			m.vp.LineUp(1)
			return true, m, nil
		}
		return false, m, nil
	case "down", "j":
		switch m.router.currentView {
		case caseView, helpView, aboutView:
			m.vp.LineDown(1)
			return true, m, nil
		}
		return false, m, nil
	case "pgup", "b":
		m.vp.HalfViewUp()
		return true, m, nil
	case "pgdown", " ":
		if m.router.currentView == caseView && msg.String() == " " {
			return false, m, nil
		}
		m.vp.HalfViewDown()
		return true, m, nil
	case "ctrl+u":
		m.vp.ViewUp()
		return true, m, nil
	case "ctrl+d":
		m.vp.ViewDown()
		return true, m, nil
	case "home", "g":
		m.vp.GotoTop()
		return true, m, nil
	case "end", "G":
		m.vp.GotoBottom()
		return true, m, nil
	case "r":
		if m.router.currentView == reviewView {
			return false, m, nil
		}
		if m.router.currentView == caseView {
			if m.results.result != nil && len(m.results.result.Assumptions) > 0 {
				m.review.assumptions = m.results.result.Assumptions
				m.review.currentIdx = 0
				m.review.mode = "browse"
				m.navigateTo(reviewView)
				return true, m, nil
			}
		} else {
			m.navigateTo(analyzeView)
			return true, m, nil
		}
	case "v":
		if m.router.currentView == caseView && m.results.result != nil && len(m.results.result.Assumptions) > 0 {
			m.validate.assumptions = m.results.result.Assumptions
			m.validate.currentIdx = 0
			m.navigateTo(validationView)
			return true, m, nil
		}
		if m.router.currentView == reviewView && len(m.review.assumptions) > 0 {
			m.validate.assumptions = m.review.assumptions
			m.validate.currentIdx = 0
			m.navigateTo(validationView)
			return true, m, nil
		}
	case "c":
		if m.router.currentView == caseView && m.activeCase != "" {
			delete(m.caseResults, m.activeCase)
			m.results.result = nil
			m.activeCase = ""
			m.results.tabStates = make(map[int]*tabState)
			m.router.rebuildCaseEntries(m.getCaseLabels())
			m.statusMsg = "Case cleared"
			m.navigateTo(analyzeView)
			return true, m, nil
		}
	case "e":
		if m.router.currentView == caseView && m.results.result != nil {
			return true, m, func() tea.Msg { return navigateMsg{to: reportsView} }
		}
	case "s":
		if m.router.currentView == settingsView && !m.settings.editing {
			m.config.Save(ConfigPath())
			m.statusMsg = "Settings saved"
			return true, m, nil
		}
	case "/":
		m.searchActive = true
		m.searchQuery = ""
		return true, m, nil
	}
	return false, m, nil
}

func (m *mainModel) handleSearchInput(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "enter":
		m.searchActive = false
		return m, nil
	case "n":
		m.vp.LineDown(1)
		return m, nil
	case "N":
		m.vp.LineUp(1)
		return m, nil
	case "backspace":
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
		}
	default:
		if len(msg.String()) == 1 {
			m.searchQuery += msg.String()
		}
	}
	return m, nil
}

func (m *mainModel) saveScroll() {
	m.scrollY[m.router.currentView] = m.vp.YOffset
}

func (m *mainModel) restoreScroll() {
	if y, ok := m.scrollY[m.router.currentView]; ok {
		m.vp.YOffset = y
	} else {
		m.vp.YOffset = 0
	}
}

func (m *mainModel) navigateTo(to view) {
	m.saveScroll()
	m.router.NavigateTo(to)
	m.restoreScroll()
}

func (m *mainModel) navigateBack() {
	m.saveScroll()
	m.router.NavigateBack()
	m.restoreScroll()
}

func (m *mainModel) getCaseLabels() []string {
	var labels []string
	for _, f := range m.recentFiles {
		if _, ok := m.caseResults[f]; ok {
			labels = append(labels, filepath.Base(f))
		}
	}
	return labels
}

func (m *mainModel) addRecentFile(path string) {
	if path == "" {
		return
	}
	dedup := make([]string, 0, len(m.recentFiles)+1)
	dedup = append(dedup, path)
	for _, f := range m.recentFiles {
		if f != path {
			dedup = append(dedup, f)
			if len(dedup) >= 10 {
				break
			}
		}
	}
	m.recentFiles = dedup
}

func (m mainModel) viewStartup() string {
	s := m.styles

	sel := m.startupMenuIdx

	fox := s.Fox.Render(` /\_/\  `)
	title := s.PremiumTitle.Render("ASF0")
	subtitle := s.DimText.Render("Assumption Security Framework Zero")
	tagline := s.Accent.Render(`"Find hidden assumptions, contradictions,`) + "\n" +
		s.Accent.Render(` trust failures, and control gaps."`)
	sep := s.SectionRule.Render(strings.Repeat("━", 50))

	menuItems := []struct {
		key  string
		desc string
	}{
		{"Enter", "Start New Analysis"},
		{"O", "Open Existing Case"},
		{"A", "Local AI Settings"},
		{"?", "Quick Tour  (recommended for first use)"},
	}

	var menuRows []string
	for i, item := range menuItems {
		active := i == sel
		menuKey := s.Accent.Render("  [" + item.key + "]  ")
		if active {
			menuKey = s.SidebarActive.Render("  [" + item.key + "]  ")
		}
		menuDesc := s.DimText.Render(item.desc)
		menuRows = append(menuRows, "  "+menuKey+"  "+menuDesc)
	}
	menu := lipgloss.JoinVertical(lipgloss.Left, menuRows...)

	quitHint := s.DimText.Render("  [Q] Quit")

	tipTitle := s.SubSectionTitle.Render("  Tip:")
	tipText := s.DimText.Render("Results appear inside case workspaces.")
	tipBody := s.DimText.Render("Critical findings: Contradictions · Trust Chains · SDRI")

	content := lipgloss.JoinVertical(lipgloss.Center,
		"",
		fox,
		"",
		title,
		"",
		subtitle,
		"",
		tagline,
		"",
		sep,
		"",
		menu,
		"",
		sep,
		"",
		tipTitle,
		tipText,
		tipBody,
		"",
		quitHint,
		"",
		s.DimText.Render("v"+ASFVersion),
	)

	return lipgloss.Place(m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		s.StartupBox.Render(content),
	)
}

func (m *mainModel) sidebarWidth() int {
	return m.layoutMgr.sidebarWidth
}

func (m *mainModel) mainWidth() int {
	w := m.width - m.sidebarWidth()
	if w < 20 {
		w = 20
	}
	return w
}

func (m *mainModel) caseTabName() string {
	names := []string{"Overview", "Assumptions", "Verification", "Contradictions", "Trust", "Controls", "SDRI"}
	if m.router.currentView == caseView && m.results.result != nil {
		idx := m.results.resultTab
		if idx >= 0 && idx < len(names) {
			return names[idx]
		}
	}
	return ""
}

func (m *mainModel) renderBreadcrumbBar() string {
	if m.router.currentView != caseView || m.results.result == nil {
		return ""
	}
	s := m.styles
	parts := []string{}
	parts = append(parts, s.Breadcrumb.Render("ASF0"))
	parts = append(parts, s.BreadcrumbSep.Render(" / "))
	fileLabel := filepath.Base(m.activeCase)
	if fileLabel == "" {
		fileLabel = "case"
	}
	parts = append(parts, s.Breadcrumb.Render(fileLabel))
	tab := m.results.resultTab
	tabName := m.caseTabName()
	if tabName != "" {
		parts = append(parts, s.BreadcrumbSep.Render(" / "))
		parts = append(parts, s.DimText.Render(tabName))
	}
	ts := m.results.tabStateFor(tab)
	if tab > 0 && ts.selectedIndex >= 0 && ts.selectedIndex < m.results.tabCount(tab) {
		parts = append(parts, s.BreadcrumbSep.Render(" / "), s.DimText.Render(fmt.Sprintf("#%d", ts.selectedIndex+1)))
	}
	if m.results.detailFocus {
		parts = append(parts, s.BreadcrumbSep.Render(" / "), s.DimText.Render("detail"))
	}
	if ts.filterActive || ts.searchQuery != "" {
		parts = append(parts, s.BreadcrumbSep.Render(" / "), s.SearchActive.Render("filter:"+ts.searchQuery))
	}
	return s.HeaderBar.Render(strings.Join(parts, ""))
}

func (m *mainModel) mainHeight() int {
	extra := 0
	if m.router.currentView == caseView && m.results.result != nil {
		extra = m.layoutMgr.breadcrumbHeight
	}
	h := m.height - m.layoutMgr.headerHeight - extra - m.layoutMgr.hintsHeight - m.layoutMgr.statusBarHeight
	if h < 5 {
		h = 5
	}
	return h
}

func (m mainModel) View() string {
	if !m.ready {
		return m.styles.BrandedLoading("Initializing...", 0)
	}
	if m.startup {
		return m.viewStartup()
	}
	if m.quitting {
		return ""
	}
	if m.err != nil {
		return m.styles.ErrorText.Render(fmt.Sprintf("Fatal Error: %v", m.err))
	}
	if m.width < 60 || m.height < 10 {
		return fmt.Sprintf("Terminal too small.\nMinimum: 60x10\nCurrent: %dx%d", m.width, m.height)
	}

	content := m.renderContent()

	if m.pickerActive {
		overlay := m.renderFilePicker(m.mainWidth(), m.mainHeight())
		content = overlay
	}

	if m.searchActive {
		searchPrompt := fmt.Sprintf("Search: %s█", m.searchQuery)
		content = m.styles.StatusWarn.Render(searchPrompt) + "\n" + content
	}

	sidebar := m.renderSidebar()
	m.vp.Width = m.mainWidth()
	m.vp.Height = m.mainHeight()
	m.vp.SetContent(content)

	mainArea := m.styles.App.Render(m.vp.View())

	body := lipgloss.JoinHorizontal(lipgloss.Top, sidebar, mainArea)

	headerBar := m.renderHeaderBar()
	breadcrumbBar := m.renderBreadcrumbBar()
	hintsBar := m.renderHintsBar()
	statusBar := m.renderStatusBar()

	views := []string{headerBar}
	if breadcrumbBar != "" {
		views = append(views, breadcrumbBar)
	}
	views = append(views, body, hintsBar, statusBar)

	return lipgloss.JoinVertical(lipgloss.Top, views...)
}

func (m mainModel) renderHeaderBar() string {
	s := m.styles
	version := "v" + ASFVersion
	fox := s.Fox.Render(" /\\_/\\  ")
	left := fmt.Sprintf(" %sASF0  %s", fox, version)
	right := " Security Assumption Framework "
	fill := m.width - lipgloss.Width(left) - lipgloss.Width(right) - 2
	if fill < 1 {
		fill = 1
	}
	return s.HeaderBar.Render(left + strings.Repeat(" ", fill) + right)
}

func (m mainModel) renderSidebar() string {
	s := m.styles
	nodes := m.router.sidebarVisibleNodes()
	availHeight := m.mainHeight()
	totalNodes := len(nodes)

	// Apply sidebar viewport offset so selection stays visible
	m.router.MaintainSidebarOffset(availHeight)
	offset := m.router.sidebarOffset
	if offset >= totalNodes && totalNodes > 0 {
		offset = totalNodes - 1
	}
	end := offset + availHeight
	if end > totalNodes {
		end = totalNodes
	}
	visibleNodes := nodes[offset:end]

	var rendered []string
	selNodeIdx := m.router.sidebarSel
	for i, n := range visibleNodes {
		actualIdx := offset + i
		if n.isSection {
			labelStr := " " + n.label + " "
			labelWidth := lipgloss.Width(labelStr)
			ruleWidth := s.sidebarInnerWidth() - labelWidth
			if ruleWidth < 1 {
				ruleWidth = 1
			}
			rendered = append(rendered, s.Texture.Render(labelStr)+s.SectionRule.Render(strings.Repeat("━", ruleWidth)))
			continue
		}

		isParent := len(n.children) > 0
		active := actualIdx == selNodeIdx && m.router.focus == focusSidebar
		viewActive := n.vid == m.router.currentView
		if n.vid == caseView && n.tab >= 0 && n.tab < len(m.recentFiles) {
			viewActive = m.recentFiles[n.tab] == m.activeCase
		}

		var prefix string
		if isParent {
			if n.expanded {
				prefix = "▾ "
			} else {
				prefix = "▸ "
			}
		} else {
			prefix = "  "
		}

		indent := ""
		if len(m.router.sidebarTree) > 0 && n != m.router.sidebarTree[0] {
			for _, parent := range m.router.sidebarTree {
				for _, child := range parent.children {
					if child == n {
						indent = "  "
						break
					}
				}
				if indent != "" {
					break
				}
			}
		}

		label := indent + prefix + n.label

		switch {
		case active:
			rendered = append(rendered, s.SidebarActive.Render(label))
		case viewActive:
			rendered = append(rendered, s.SidebarParent.Render(label))
		case isParent:
			rendered = append(rendered, s.SidebarParent.Render(label))
		default:
			rendered = append(rendered, s.SidebarItem.Render(label))
		}
	}
	sidebarContent := lipgloss.JoinVertical(lipgloss.Left, rendered...)
	lines := strings.Count(sidebarContent, "\n") + 1
	if lines < availHeight {
		sidebarContent += strings.Repeat("\n", availHeight-lines)
	}

	if totalNodes > availHeight {
		scrollPct := int(float64(offset+availHeight) / float64(totalNodes) * 100)
		if scrollPct > 100 {
			scrollPct = 100
		}
		scrollBar := fmt.Sprintf(" %s %d%%", s.DimText.Render("▐█▌"), scrollPct)
		sidebarContent += "\n" + s.DimText.Render(scrollBar)
	}
	return s.Sidebar.Render(sidebarContent)
}

func (s StyleSet) sidebarInnerWidth() int {
	return s.Sidebar.GetWidth() - 2
}

func (m mainModel) renderHintsBar() string {
	s := m.styles
	var hints []string

	guidance := ""
	switch m.router.currentView {
	case analyzeView:
		guidance = "New Analysis — Select an architecture document to begin analysis"
		if m.analyze.running {
			hints = append(hints, s.DimText.Render("Esc Cancel"))
		} else {
			hints = append(hints, s.DimText.Render("Enter Select"))
		}
	case caseView:
		guidance = "Case Workspace — Explore findings across tabs"
		if m.results.resultTab > 0 {
			ts := m.results.tabStateFor(m.results.resultTab)
			if m.results.detailFocus {
				hints = append(hints, s.DimText.Render("↑↓ Scroll Detail"))
				hints = append(hints, s.DimText.Render("Esc List"))
			} else {
				hints = append(hints, s.DimText.Render("↑↓ Select"))
				hints = append(hints, s.DimText.Render("Enter Detail"))
				hints = append(hints, s.DimText.Render("/ Search"))
			}
			if ts.filterActive || ts.searchQuery != "" {
				hints = append(hints, s.Accent.Render(fmt.Sprintf("filter: %s", ts.searchQuery)))
			}
		} else {
			hints = append(hints, s.DimText.Render("↑↓ Scroll"))
		}
		hints = append(hints, s.DimText.Render("←→ Tabs"))
		hints = append(hints, s.DimText.Render("r Review"))
		hints = append(hints, s.DimText.Render("v Validate"))
		hints = append(hints, s.DimText.Render("e Reports"))
		hints = append(hints, s.DimText.Render("c Clear"))
	case reviewView:
		guidance = "Review Queue — Human analyst approval workflow for assumptions"
		if m.review.editing {
			hints = append(hints, s.DimText.Render("Enter Save"))
			hints = append(hints, s.DimText.Render("Esc Cancel"))
		} else {
			hints = append(hints, s.DimText.Render("↑↓ Navigate"))
			hints = append(hints, s.DimText.Render("Enter Detail"))
			hints = append(hints, s.Accent.Render("s Accept"))
			hints = append(hints, s.DimText.Render("r Reject"))
			hints = append(hints, s.DimText.Render("m Modify"))
			hints = append(hints, s.DimText.Render("n Notes"))
			hints = append(hints, s.DimText.Render("v Validate"))
			hints = append(hints, s.DimText.Render("Tab Sidebar"))
		}
	case validationView:
		guidance = "Validation Queue — Evidence-backed verification workflow for assumptions"
		hints = append(hints, s.DimText.Render("↑↓ Navigate"))
		hints = append(hints, s.DimText.Render("Enter Detail"))
	case settingsView:
		guidance = "Settings — Configure analysis engine, output, and preferences"
		if m.settings.editing {
			hints = append(hints, s.DimText.Render("←→ Change"))
			hints = append(hints, s.DimText.Render("Esc Done"))
		} else {
			hints = append(hints, s.DimText.Render("Enter Edit"))
			hints = append(hints, s.DimText.Render("s Save"))
		}
	case reportsView:
		guidance = "Reports — Generate and export analysis results (PDF, HTML, JSON, CSV, Markdown)"
		if m.reportsV.showConfirmation || m.reportsV.done {
			hints = append(hints, s.DimText.Render("Esc Back"))
		} else {
			hints = append(hints, s.DimText.Render("↑↓ Select"))
			hints = append(hints, s.DimText.Render("Enter Choose"))
		}
	case helpView:
		guidance = "Help — Keyboard shortcuts, workflow guide, and documentation"
		hints = append(hints, s.DimText.Render("↑↓ Scroll"))
		hints = append(hints, s.DimText.Render("/ Search"))
	case aboutView:
		guidance = "About — Version, license, and system information"
		hints = append(hints, s.DimText.Render("Q Quit"))
	case localAIView:
		guidance = "Local AI — Manage Ollama models for AI-assisted analysis"
		hints = append(hints, s.DimText.Render("↑↓ Select"))
		hints = append(hints, s.DimText.Render("Enter Action"))
		hints = append(hints, s.DimText.Render("Esc Cancel"))
	}
	if guidance != "" {
		hints = append(hints, s.Accent.Render(guidance))
	}
	if m.router.focus == focusSidebar {
		hints = append(hints, s.Accent.Render(" [Sidebar]"))
		hints = append(hints, s.DimText.Render("Tab Content"))
	} else {
		hints = append(hints, s.DimText.Render("Tab Sidebar"))
	}
	hints = append(hints, s.DimText.Render("? Help"))
	hints = append(hints, s.DimText.Render("q Back"))
	hints = append(hints, s.DimText.Render("Q Quit"))

	scrollPct := m.viewportScrollPercent()
	if scrollPct != "" {
		hints = append(hints, s.DimText.Render(scrollPct))
	}

	hintStr := strings.Join(hints, "  │  ")
	avail := m.width - 2
	if w := lipgloss.Width(hintStr); w < avail {
		hintStr += strings.Repeat(" ", avail-w)
	}
	return s.HintsBar.Render(hintStr)
}

func (m mainModel) renderStatusBar() string {
	s := m.styles
	version := "v" + ASFVersion
	mode := "ASF Engine"
	if m.config.AI.Enabled {
		mode = s.Accent.Render("AI Enhanced")
	}
	file := m.activeCase
	if file == "" {
		file = s.DimText.Render("no case")
	} else {
		file = s.Value.Render(filepath.Base(file))
	}
	state := m.statusMsg
	if state == "" {
		state = s.DimText.Render("ready")
	} else {
		state = s.StatusGood.Render(state)
	}

	left := fmt.Sprintf("  %s  %s  %s", s.DimText.Render(version), mode, file)
	right := fmt.Sprintf(" %s ", state)
	fill := m.width - lipgloss.Width(left) - lipgloss.Width(right) - 2
	if fill < 1 {
		fill = 1
	}
	return s.StatusBar.Render(left + strings.Repeat(" ", fill) + right)
}

func (m mainModel) viewportScrollPercent() string {
	total := m.vp.TotalLineCount()
	visible := m.vp.Height
	offset := m.vp.YOffset
	if total <= visible || total == 0 {
		return ""
	}
	pct := int(float64(offset+visible) / float64(total) * 100)
	if pct > 100 {
		pct = 100
	}
	first := offset + 1
	last := offset + visible
	if last > total {
		last = total
	}
	return fmt.Sprintf("Line %d–%d / %d  (%d%%)", first, last, total, pct)
}

func (m mainModel) updateHelp(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m mainModel) renderContent() string {
	switch m.router.currentView {
	case analyzeView:
		return m.viewAnalyze()
	case caseView:
		return m.viewResults()
	case settingsView:
		return m.viewSettings()
	case aboutView:
		return m.viewAbout()
	case reportsView:
		return m.viewReports()
	case reviewView:
		return m.viewReview()
	case validationView:
		return m.viewValidation()
	case helpView:
		return m.viewHelp()
	case localAIView:
		return m.viewLocalAI()
	case tourView:
		return m.viewTour()
	}
	return ""
}
