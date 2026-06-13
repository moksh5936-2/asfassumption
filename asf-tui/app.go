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

type errorMsg string

func (e errorMsg) Error() string { return string(e) }

type viewHistory struct {
	views []view
}

func (v *viewHistory) push(vw view) {
	v.views = append(v.views, vw)
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
	currentView  view
	styles       StyleSet
	config       *Config
	engine       *Engine
	history      viewHistory
	quitting     bool
	err          error
	statusMsg    string
	currentFile  string
	sidebarOpen  bool
	sidebarSel   int
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

	vp      viewport.Model
	scrollY map[view]int
}

type navigateMsg struct {
	to view
}

var sidebarItems = []string{
	"Dashboard",
	"Analyze",
	"Results",
	"File Explorer",
	"AI Models",
	"Settings",
	"About",
	"Help",
}

var sidebarViews = []view{
	dashboardView,
	analyzeView,
	resultsView,
	fileBrowserView,
	localaiView,
	settingsView,
	aboutView,
	helpView,
}

func newMainModel(cfg *Config) *mainModel {
	s := NewStyles(Themes[cfg.Appearance.Theme])
	e := NewEngine(cfg)
	return &mainModel{
		currentView: startupView,
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
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyMsg(msg)

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
		switch m.currentView {
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
		m.history.push(m.currentView)
		m.currentView = resultsView
		m.vp.YOffset = 0
		m.scrollY[resultsView] = 0
		return m, nil

	case fileSelectedMsg:
		m.analyze.setDocPath(string(msg))
		m.currentFile = string(msg)
		m.statusMsg = "File selected: " + string(msg)
		m.addRecentFile(string(msg))
		m.saveScroll()
		m.history.push(m.currentView)
		m.currentView = analyzeView
		m.vp.YOffset = 0
		m.scrollY[analyzeView] = 0
		return m, nil
	}

	switch m.currentView {
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

func (m mainModel) handleKeyMsg(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.searchActive {
		return m.handleSearchInput(msg)
	}

	switch msg.String() {
	case "ctrl+c", "Q":
		m.quitting = true
		return m, tea.Quit
	case "q":
		if m.currentView == startupView {
			m.quitting = true
			return m, tea.Quit
		}
		m.navigateBack()
		return m, nil
	case "?":
		m.navigateTo(helpView)
		return m, nil
	case "esc":
		return m.handleBack()
	case "tab":
		if m.currentView == resultsView {
			break
		}
		m.cycleSidebar(1)
		target := viewForSidebar(m.sidebarSel)
		if target != m.currentView {
			m.navigateTo(target)
		}
		return m, nil
	case "shift+tab":
		if m.currentView == resultsView {
			break
		}
		m.cycleSidebar(-1)
		target := viewForSidebar(m.sidebarSel)
		if target != m.currentView {
			m.navigateTo(target)
		}
		return m, nil
	case "up", "k":
		m.vp.LineUp(1)
		return m, nil
	case "down", "j":
		m.vp.LineDown(1)
		return m, nil
	case "pgup", "b":
		m.vp.HalfViewUp()
		return m, nil
	case "pgdn", " ":
		if m.currentView == resultsView && msg.String() == " " {
			break
		}
		m.vp.HalfViewDown()
		return m, nil
	case "ctrl+u":
		m.vp.ViewUp()
		return m, nil
	case "ctrl+d":
		m.vp.ViewDown()
		return m, nil
	case "home", "g":
		m.vp.GotoTop()
		return m, nil
	case "end", "G":
		m.vp.GotoBottom()
		return m, nil
	case "f":
		m.fileBrowse.path, _ = getDefaultPath(m.config)
		m.navigateTo(fileBrowserView)
		return m, nil
	case "r":
		if m.currentView == resultsView {
			if m.results.result != nil && len(m.results.result.Assumptions) > 0 {
				m.review.assumptions = m.results.result.Assumptions
				m.review.currentIdx = 0
				m.review.mode = "browse"
				m.navigateTo(reviewView)
				return m, nil
			}
		} else {
			m.navigateTo(analyzeView)
			return m, nil
		}
	case "v":
		if m.currentView == resultsView && m.results.result != nil && len(m.results.result.Assumptions) > 0 {
			m.validate.assumptions = m.results.result.Assumptions
			m.validate.currentIdx = 0
			m.navigateTo(validationView)
			return m, nil
		}
		if m.currentView == reviewView && len(m.review.assumptions) > 0 {
			m.validate.assumptions = m.review.assumptions
			m.validate.currentIdx = 0
			m.navigateTo(validationView)
			return m, nil
		}
	case "c":
		if m.currentView == resultsView && m.results.result != nil {
			m.results.result = nil
			m.statusMsg = "Results cleared"
			m.navigateTo(analyzeView)
			return m, nil
		}
	case "e":
		if m.currentView == resultsView && m.results.result != nil {
			return m, func() tea.Msg { return navigateMsg{to: exportView} }
		}
	case "s":
		if m.currentView == settingsView && !m.settings.editing {
			m.config.Save(ConfigPath())
			m.statusMsg = "Settings saved"
		}
	case "/":
		m.searchActive = true
		m.searchQuery = ""
		return m, nil
	}
	return m, nil
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
	m.scrollY[m.currentView] = m.vp.YOffset
}

func (m *mainModel) restoreScroll() {
	if y, ok := m.scrollY[m.currentView]; ok {
		m.vp.YOffset = y
	} else {
		m.vp.YOffset = 0
	}
}

func (m *mainModel) navigateTo(to view) {
	m.saveScroll()
	m.history.push(m.currentView)
	m.currentView = to
	m.restoreScroll()
}

func (m *mainModel) navigateBack() {
	m.saveScroll()
	if prev, ok := m.history.pop(); ok {
		m.currentView = prev
	} else {
		switch m.currentView {
		case fileBrowserView, dashboardView, analyzeView, resultsView, localaiView, settingsView, aboutView, exportView, reviewView, validationView, helpView:
			m.currentView = dashboardView
		}
	}
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

func (m *mainModel) cycleSidebar(dir int) {
	total := len(sidebarItems)
	m.sidebarSel = (m.sidebarSel + dir + total) % total
}

func (m *mainModel) sidebarWidth() int {
	if !m.sidebarOpen {
		return 0
	}
	return 23
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

func (m mainModel) handleBack() (tea.Model, tea.Cmd) {
	m.navigateBack()
	return m, nil
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
	items := sidebarItems
	var rendered []string
	for i, item := range items {
		if i == m.sidebarSel && viewForSidebar(i) == m.currentView {
			rendered = append(rendered, m.styles.SidebarActive.Render(" "+item))
		} else {
			rendered = append(rendered, m.styles.SidebarItem.Render(" "+item))
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

func viewForSidebar(i int) view {
	if i >= len(sidebarViews) {
		return dashboardView
	}
	return sidebarViews[i]
}

func (m mainModel) renderBottomBar() string {
	var hints []string
	hints = append(hints, "Tab:Nav")
	if m.currentView == startupView {
		hints = append(hints, "↑↓:Navigate", "Enter:Select", "q:Quit")
	} else {
		hints = append(hints, "q:Back", "?:Help", "f:File Explorer", "r:Run")
	}
	switch m.currentView {
	case resultsView:
		hints = append(hints, "e:Export", "c:Clear", "r:Review", "v:Validate", "PgUp/Dn:Scroll")
	case dashboardView:
		hints = append(hints, "Enter:Select")
	case analyzeView:
		hints = append(hints, "Enter:Select")
	case reviewView:
		hints = append(hints, "s:Accept", "r:Reject", "m:Modified", "n:Note")
	}

	scrollPct := m.viewportScrollPercent()
	if scrollPct != "" {
		hints = append(hints, scrollPct)
	}

	hintStr := strings.Join(hints, "  •  ")
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
	switch m.currentView {
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
