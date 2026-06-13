package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type view int

const (
	startupView view = iota
	dashboardView
	analyzeView
	resultsView
	fileBrowserView
	localaiView
	settingsView
	aboutView
	exportView
	reviewView
	validationView
	helpView
)

var viewNames = map[view]string{
	startupView:     "Startup",
	dashboardView:   "Dashboard",
	analyzeView:     "Analyze",
	resultsView:     "Results",
	fileBrowserView: "File Explorer",
	localaiView:     "AI Models",
	settingsView:    "Settings",
	aboutView:       "About",
	helpView:        "Help",
}

type sidebarEntry struct {
	name string
	vid  view
	tab  int
}

type focusManager struct {
	activeView view
	subFocus   string
}

type layoutManager struct {
	sidebarWidth    int
	topBarHeight    int
	bottomBarHeight int
}

func newFocusManager() focusManager {
	return focusManager{}
}

func newLayoutManager() layoutManager {
	return layoutManager{
		sidebarWidth:    23,
		topBarHeight:    1,
		bottomBarHeight: 1,
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
		return startupView, false
	}
	last := v.views[len(v.views)-1]
	v.views = v.views[:len(v.views)-1]
	return last, true
}

type mainModel struct {
	ready        bool
	width        int
	height       int
	router       Router
	styles       StyleSet
	config       *Config
	engine       *Engine
	quitting     bool
	err          error
	statusMsg    string
	currentFile  string
	sidebarOpen  bool
	searchActive bool
	searchQuery  string
	recentFiles  []string

	startup    startupModel
	dash       dashboardModel
	analyze    analyzeModel
	results    resultsModel
	fileBrowse fileBrowserModel
	localai    localaiModel
	settings   settingsModel
	about      aboutModel
	exportV    exportModel
	review     reviewModel
	validate   validationModel
	help       helpModel

	vp        viewport.Model
	scrollY   map[view]int
	focusMgr  focusManager
	layoutMgr layoutManager
}

type navigateMsg struct {
	to view
}

var sidebarEntries = []sidebarEntry{
	{"Dashboard", dashboardView, -1},
	{"File Explorer", fileBrowserView, -1},
	{"Analyze", analyzeView, -1},
	{"Assumptions", resultsView, 1},
	{"Verification", resultsView, 2},
	{"Contradictions", resultsView, 3},
	{"Trust Chains", resultsView, 4},
	{"Single Points of Trust", resultsView, 4},
	{"Assumption Impact Analysis", resultsView, 5},
	{"Blind Spots", resultsView, 6},
	{"SDRI", resultsView, 9},
	{"Recommended Controls", resultsView, 7},
	{"Security Design Review", resultsView, 10},
	{"Reports / Exports", resultsView, 8},
	{"Settings", settingsView, -1},
	{"Help", helpView, -1},
}

func newMainModel(cfg *Config) *mainModel {
	s := NewStyles(Themes[cfg.Appearance.Theme])
	e := NewEngine(cfg)
	return &mainModel{
		router:      newRouter(),
		styles:      s,
		config:      cfg,
		engine:      e,
		sidebarOpen: true,
		vp:          viewport.New(0, 0),
		scrollY:     make(map[view]int),
		startup:     newStartupModel(),
		dash:        newDashboardModel(),
		analyze:     newAnalyzeModel(e),
		results:     newResultsModel(),
		fileBrowse:  newFileBrowserModel(),
		localai:     newLocalAIModel(cfg),
		settings:    newSettingsModel(cfg),
		about:       newAboutModel(),
		exportV:     newExportModel(),
		review:      newReviewModel(),
		validate:    newValidationModel(),
		help:        newHelpModel(),
		focusMgr:    newFocusManager(),
		layoutMgr:   newLayoutManager(),
	}
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.ready = true
		m.width = msg.Width
		m.height = msg.Height
		m.styles = NewStyles(m.styles.Theme())
		m.vp.Width = m.mainWidth()
		m.vp.Height = m.mainHeight()

	case tea.KeyMsg:
		if m.searchActive {
			return m.handleSearchInput(msg)
		}
		handled, model, cmd := m.handleGlobalKey(msg)
		if handled {
			return model, cmd
		}

	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			m.vp, _ = m.vp.Update(msg)
			return m, nil
		case tea.MouseWheelDown:
			m.vp, _ = m.vp.Update(msg)
			return m, nil
		}

	case navigateMsg:
		if msg.to == exportView {
			m.exportV.selected = 0
			m.exportV.done = false
			m.exportV.exportPath = ""
			m.exportV.showConfirmation = false
			m.exportV.err = nil
			m.exportV.result = m.results.result
			m.exportV.outputDir = m.config.Output.Directory
			if m.exportV.outputDir == "" {
				m.exportV.outputDir = "./reports"
			}
			m.exportV.format = exportFormatFromConfig(m.config)
		}
		if msg.to == resultsView && m.results.result != nil {
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
		m.results.result = msg.result
		m.results.resultTab = 0
		m.currentFile = m.analyze.docPath()
		m.statusMsg = "Analysis complete"
		m.addRecentFile(m.currentFile)
		m.saveScroll()
		m.router.NavigateTo(resultsView)
		m.vp.YOffset = 0
		m.scrollY[resultsView] = 0
		return m, nil

	case fileSelectedMsg:
		m.analyze.setDocPath(string(msg))
		m.currentFile = string(msg)
		m.statusMsg = "File selected: " + string(msg)
		m.addRecentFile(string(msg))
		m.saveScroll()
		m.router.NavigateTo(analyzeView)
		m.vp.YOffset = 0
		m.scrollY[analyzeView] = 0
		return m, nil
	}

	switch m.router.currentView {
	case startupView:
		return m.updateStartup(msg)
	case dashboardView:
		return m.updateDashboard(msg)
	case analyzeView:
		return m.updateAnalyze(msg)
	case resultsView:
		return m.updateResults(msg)
	case fileBrowserView:
		return m.updateFileBrowser(msg)
	case localaiView:
		return m.updateLocalAI(msg)
	case settingsView:
		return m.updateSettings(msg)
	case aboutView:
		return m.updateAbout(msg)
	case exportView:
		return m.updateExport(msg)
	case reviewView:
		return m.updateReview(msg)
	case validationView:
		return m.updateValidation(msg)
	case helpView:
		return m.updateHelp(msg)
	}
	return m, nil
}

func (m mainModel) handleGlobalKey(msg tea.KeyMsg) (bool, tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "Q":
		m.quitting = true
		return true, m, tea.Quit
	case "q":
		if m.router.currentView == startupView {
			m.quitting = true
			return true, m, tea.Quit
		}
		m.navigateBack()
		return true, m, nil
	case "?":
		m.navigateTo(helpView)
		return true, m, nil
	case "esc":
		switch m.router.currentView {
		case analyzeView:
			if m.analyze.running || m.analyze.inputMode != "" {
				return false, m, nil
			}
		case settingsView:
			if m.settings.editing {
				return false, m, nil
			}
		case localaiView:
			if m.localai.showActions {
				return false, m, nil
			}
		case exportView:
			if m.exportV.showConfirmation || m.exportV.done {
				return false, m, nil
			}
		case reviewView:
			if m.review.editing {
				return false, m, nil
			}
		}
		m.navigateBack()
		return true, m, nil
	case "tab":
		switch m.router.currentView {
		case resultsView, fileBrowserView:
			return false, m, nil
		}
		m.saveScroll()
		m.router.CycleSidebar(1)
		m.router.ActivateSidebar()
		tab := m.router.ActivateSidebarTab()
		if tab >= 0 {
			m.results.resultTab = tab
		}
		m.restoreScroll()
		return true, m, nil
	case "shift+tab":
		switch m.router.currentView {
		case resultsView, fileBrowserView:
			return false, m, nil
		}
		m.saveScroll()
		m.router.CycleSidebar(-1)
		m.router.ActivateSidebar()
		tab := m.router.ActivateSidebarTab()
		if tab >= 0 {
			m.results.resultTab = tab
		}
		m.restoreScroll()
		return true, m, nil
	case "up", "k":
		switch m.router.currentView {
		case resultsView, helpView, aboutView:
			m.vp.LineUp(1)
			return true, m, nil
		}
		return false, m, nil
	case "down", "j":
		switch m.router.currentView {
		case resultsView, helpView, aboutView:
			m.vp.LineDown(1)
			return true, m, nil
		}
		return false, m, nil
	case "pgup", "b":
		m.vp.HalfViewUp()
		return true, m, nil
	case "pgdown", " ":
		if m.router.currentView == resultsView && msg.String() == " " {
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
	case "f":
		m.fileBrowse.path, _ = getDefaultPath(m.config)
		m.navigateTo(fileBrowserView)
		return true, m, nil
	case "r":
		if m.router.currentView == reviewView {
			return false, m, nil
		}
		if m.router.currentView == resultsView {
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
		if m.router.currentView == resultsView && m.results.result != nil && len(m.results.result.Assumptions) > 0 {
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
		if m.router.currentView == resultsView && m.results.result != nil {
			m.results.result = nil
			m.statusMsg = "Results cleared"
			m.navigateTo(analyzeView)
			return true, m, nil
		}
	case "e":
		if m.router.currentView == resultsView && m.results.result != nil {
			return true, m, func() tea.Msg { return navigateMsg{to: exportView} }
		}
	case "s":
		if m.router.currentView == settingsView && !m.settings.editing {
			m.config.Save(ConfigPath())
			m.statusMsg = "Settings saved"
			return true, m, nil
		}
	case "/":
		if m.router.currentView == fileBrowserView {
			return false, m, nil
		}
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

func (m *mainModel) sidebarWidth() int {
	if !m.sidebarOpen {
		return 0
	}
	return m.layoutMgr.sidebarWidth
}

func (m *mainModel) mainWidth() int {
	w := m.width - m.sidebarWidth() - 2
	if w < 10 {
		w = 10
	}
	return w
}

func (m *mainModel) mainHeight() int {
	h := m.height - 3
	if h < 3 {
		h = 3
	}
	return h
}

func (m mainModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	if m.quitting {
		return ""
	}
	if m.err != nil {
		return m.styles.ErrorText.Render(fmt.Sprintf("Fatal Error: %v", m.err))
	}
	if m.width < 60 || m.height < 12 {
		return fmt.Sprintf("Terminal too small.\nMinimum: 60x12\nCurrent: %dx%d", m.width, m.height)
	}

	content := m.renderContent()
	sidebar := m.renderSidebar()

	m.vp.Width = m.mainWidth()
	m.vp.Height = m.mainHeight()
	m.vp.SetContent(content)

	mainArea := m.styles.App.Render(m.vp.View())

	var body string
	if m.sidebarOpen {
		body = lipgloss.JoinHorizontal(lipgloss.Top,
			sidebar,
			mainArea,
		)
	} else {
		body = mainArea
	}

	topBar := m.renderTopBar()
	bottomBar := m.renderBottomBar()

	return lipgloss.JoinVertical(lipgloss.Top,
		topBar,
		body,
		bottomBar,
	)
}

func (m mainModel) renderTopBar() string {
	version := "v" + ASFVersion
	file := m.currentFile
	if file == "" {
		file = "no file"
	} else {
		if len(file) > 40 {
			file = "..." + file[len(file)-37:]
		}
	}
	status := m.statusMsg
	if status == "" {
		status = "ready"
	}
	left := fmt.Sprintf(" ASF %s  │  %s", version, file)
	right := fmt.Sprintf(" %s ", status)
	fill := m.width - lipgloss.Width(left) - lipgloss.Width(right) - 4
	if fill < 1 {
		fill = 1
	}
	return m.styles.TopBar.Render(left + strings.Repeat(" ", fill) + right)
}

func (m mainModel) renderSidebar() string {
	var rendered []string
	for i, e := range sidebarEntries {
		active := i == m.router.sidebarSel && e.vid == m.router.currentView
		if e.vid == resultsView && m.router.currentView == resultsView {
			for j, se := range sidebarEntries {
				if se.vid == resultsView && se.tab == m.results.resultTab {
					active = i == j
					break
				}
			}
		}
		if active {
			rendered = append(rendered, m.styles.SidebarActive.Render(" "+e.name))
		} else {
			rendered = append(rendered, m.styles.SidebarItem.Render(" "+e.name))
		}
	}
	sidebarContent := lipgloss.JoinVertical(lipgloss.Left, rendered...)
	availHeight := m.mainHeight()
	lines := strings.Count(sidebarContent, "\n") + 1
	if lines < availHeight {
		sidebarContent += strings.Repeat("\n", availHeight-lines)
	}
	return m.styles.Sidebar.Render(sidebarContent)
}

func (m mainModel) renderBottomBar() string {
	var hints []string
	if m.router.currentView == startupView {
		hints = append(hints, "↑↓=Navigate", "Enter=Select", "Q=Quit")
	} else {
		hints = append(hints, "F=Files", "R=Analyze", "/=Search", "?=Help", "Q=Quit")
		switch m.router.currentView {
		case resultsView:
			hints = append(hints, "Tab=Tabs", "E=Export", "C=Clear")
		case fileBrowserView:
			hints = append(hints, "Tab=Preview", ".=Hidden")
		case settingsView:
			hints = append(hints, "Enter=Edit", "S=Save")
		case reviewView:
			hints = append(hints, "S=Accept", "R=Reject", "M=Mod", "N=Note")
		}
	}

	scrollPct := m.viewportScrollPercent()
	if scrollPct != "" {
		hints = append(hints, scrollPct)
	}

	hintStr := strings.Join(hints, "  ")
	fill := m.width - lipgloss.Width(hintStr) - 4
	if fill > 0 {
		hintStr = hintStr + strings.Repeat(" ", fill)
	}
	return m.styles.BottomBar.Render(hintStr)
}

func (m mainModel) viewportScrollPercent() string {
	total := m.vp.TotalLineCount()
	visible := m.vp.Height
	offset := m.vp.YOffset
	if total <= visible || total == 0 {
		return "All"
	}
	pct := int(float64(offset+visible) / float64(total) * 100)
	if pct > 100 {
		pct = 100
	}
	from := offset + 1
	to := offset + visible
	if to > total {
		to = total
	}
	return fmt.Sprintf("%d-%d/%d (%d%%)", from, to, total, pct)
}

func (m mainModel) renderContent() string {
	switch m.router.currentView {
	case startupView:
		return m.viewStartup()
	case dashboardView:
		return m.viewDashboard()
	case analyzeView:
		return m.viewAnalyze()
	case resultsView:
		return m.viewResults()
	case fileBrowserView:
		return m.viewFileBrowser()
	case localaiView:
		return m.viewLocalAI()
	case settingsView:
		return m.viewSettings()
	case aboutView:
		return m.viewAbout()
	case exportView:
		return m.viewExport()
	case reviewView:
		return m.viewReview()
	case validationView:
		return m.viewValidation()
	case helpView:
		return m.viewHelp()
	}
	return ""
}

func (m mainModel) updateHelp(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}
