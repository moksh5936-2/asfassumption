package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errorMsg string

func (e errorMsg) Error() string { return string(e) }

type view int

const (
	startupView view = iota
	dashboardView
	analyzeView
	resultsView
	localaiView
	settingsView
	aboutView
	exportView
	reviewView
	validationView
)

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
	ready       bool
	width       int
	height      int
	currentView view
	styles      StyleSet
	config      *Config
	engine      *Engine
	history     viewHistory

	startup  startupModel
	dash     dashboardModel
	analyze  analyzeModel
	results  resultsModel
	localai  localaiModel
	settings settingsModel
	about    aboutModel
	exportV  exportModel
	review   reviewModel
	validate validationModel

	quitting bool
	err      error
	scrollY  int
}

type navigateMsg struct {
	to view
}

func newMainModel(cfg *Config) *mainModel {
	s := NewStyles(Themes[cfg.Appearance.Theme])
	e := NewEngine(cfg)
	return &mainModel{
		currentView: startupView,
		styles:      s,
		config:      cfg,
		engine:      e,
		startup:     newStartupModel(),
		dash:        newDashboardModel(),
		analyze:     newAnalyzeModel(e),
		results:     newResultsModel(),
		localai:     newLocalAIModel(cfg),
		settings:    newSettingsModel(cfg),
		about:       newAboutModel(),
		exportV:     newExportModel(),
		review:      newReviewModel(),
		validate:    newValidationModel(),
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
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "Q":
			if m.currentView == startupView {
				m.quitting = true
				return m, tea.Quit
			}
		case "r":
			if m.currentView == resultsView && len(m.results.result.Assumptions) > 0 {
				m.review.assumptions = m.results.result.Assumptions
				m.review.currentIdx = 0
				m.review.mode = "browse"
				m.history.push(m.currentView)
				m.currentView = reviewView
				return m, nil
			}
		case "v":
			if m.currentView == resultsView && len(m.results.result.Assumptions) > 0 {
				m.validate.assumptions = m.results.result.Assumptions
				m.validate.currentIdx = 0
				m.history.push(m.currentView)
				m.currentView = validationView
				return m, nil
			}
			if m.currentView == reviewView && len(m.review.assumptions) > 0 {
				m.validate.assumptions = m.review.assumptions
				m.validate.currentIdx = 0
				m.history.push(m.currentView)
				m.currentView = validationView
				return m, nil
			}
		case "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "esc":
			if m.currentView == analyzeView && m.analyze.inputMode != "" {
				m.analyze.inputMode = ""
				m.analyze.inputBuf = ""
				return m, nil
			}
			m.scrollY = 0
			return m.handleBack()
		case "pgup":
			m.scrollY -= m.height - 3
			if m.scrollY < 0 {
				m.scrollY = 0
			}
		case "pgdn":
			m.scrollY += m.height - 3
		case "ctrl+u":
			m.scrollY -= m.height / 2
			if m.scrollY < 0 {
				m.scrollY = 0
			}
		case "ctrl+d":
			m.scrollY += m.height / 2
		}

	case navigateMsg:
		m.scrollY = 0
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
		m.history.push(m.currentView)
		m.currentView = msg.to
		return m, nil

	case errorMsg:
		switch m.currentView {
		case analyzeView:
			m.analyze.running = false
			m.analyze.statusMsg = fmt.Sprintf("ASF Engine error: %s", string(msg))
			return m, nil
		default:
			m.err = error(msg)
			return m, nil
		}
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
	}

	return m, nil
}

func (m mainModel) handleBack() (tea.Model, tea.Cmd) {
	if prev, ok := m.history.pop(); ok {
		m.currentView = prev
		return m, nil
	}
	switch m.currentView {
	case dashboardView:
		m.currentView = startupView
	case analyzeView:
		if m.analyze.running {
			return m, nil
		}
		m.currentView = dashboardView
	case resultsView:
		m.currentView = dashboardView
	case localaiView, settingsView, aboutView, exportView, reviewView, validationView:
		m.currentView = dashboardView
	}
	return m, nil
}

func (m mainModel) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}

	if m.err != nil {
		return m.styles.ErrorText.Render(fmt.Sprintf("Error: %v", m.err))
	}

	var content string
	switch m.currentView {
	case startupView:
		content = m.viewStartup()
	case dashboardView:
		content = m.viewDashboard()
	case analyzeView:
		content = m.viewAnalyze()
	case resultsView:
		content = m.viewResults()
	case localaiView:
		content = m.viewLocalAI()
	case settingsView:
		content = m.viewSettings()
	case aboutView:
		content = m.viewAbout()
	case exportView:
		content = m.viewExport()
	case reviewView:
		content = m.viewReview()
	case validationView:
		content = m.viewValidation()
	}

	help := m.renderHelp()
	helpLines := strings.Count(help, "\n") + 1
	availLines := m.height - helpLines - 2

	allLines := strings.Split(content, "\n")
	if m.scrollY >= len(allLines) {
		m.scrollY = 0
	}
	if len(allLines) > availLines && availLines > 0 {
		start := m.scrollY
		end := m.scrollY + availLines
		if end > len(allLines) {
			end = len(allLines)
		}
		visible := allLines[start:end]
		dim := lipgloss.NewStyle().Foreground(m.styles.Theme().DimText)
		var prefix, suffix string
		if m.scrollY > 0 {
			prefix = dim.Render(fmt.Sprintf("(↑ %d more — PgUp) ", m.scrollY))
		}
		if m.scrollY+len(visible) < len(allLines) {
			remaining := len(allLines) - (m.scrollY + len(visible))
			suffix = "\n" + dim.Render(fmt.Sprintf("(↓ %d more — PgDn) ", remaining))
		}
		content = prefix + strings.Join(visible, "\n") + suffix
	} else {
		m.scrollY = 0
	}

	return m.styles.App.Render(
		lipgloss.JoinVertical(lipgloss.Top,
			content,
			help,
		),
	)
}

func (m mainModel) renderHelp() string {
	if m.quitting {
		return ""
	}
	helpKeys := map[view][]string{
		startupView:    {"↑↓: Navigate", "Enter: Select", "q: Quit"},
		dashboardView:  {"↑↓: Navigate", "Enter: Select", "Esc: Back", "q: Quit"},
		analyzeView:    {"↑↓: Navigate", "Enter: Edit/Select", "Esc: Back"},
		resultsView:    {"↑↓: Navigate", "Enter: Toggle", "PgDn/PgUp: Scroll", "e: Export", "r: Review", "v: Validate", "Esc: Back"},
		reviewView:     {"↑↓: Navigate", "Enter: Toggle", "PgDn/PgUp: Scroll", "s: Accept", "r: Reject", "m: Modified", "n: Note", "v: Validate", "Esc: Back"},
		validationView: {"↑↓: Navigate", "Enter: Detail", "PgDn/PgUp: Scroll", "Esc: Back"},
		localaiView:    {"↑↓: Navigate", "Enter: Select action", "PgDn/PgUp: Scroll", "Esc: Back"},
		settingsView:   {"↑↓: Navigate", "Enter: Change value", "PgDn/PgUp: Scroll", "Esc: Back"},
		aboutView:      {"Esc: Back"},
		exportView:     {"↑↓: Navigate", "Enter: Select format", "y: Confirm export", "Esc: Back"},
	}

	keys := helpKeys[m.currentView]
	if keys == nil {
		keys = []string{"Esc: Back"}
	}

	helpStr := strings.Join(keys, "  •  ")
	separator := strings.Repeat(" ", 1)
	if m.width > len(helpStr)+4 {
		fill := m.width - lipgloss.Width(helpStr) - 4
		if fill > 0 {
			separator = strings.Repeat(" ", fill/2)
		}
	}

	return m.styles.Help.Render(fmt.Sprintf("%s%s%s",
		separator, helpStr, separator,
	))
}
