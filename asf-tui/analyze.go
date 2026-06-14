package main

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type analysisCompleteMsg struct {
	result *AnalysisResult
}

type menuItem struct {
	label string
	value string
	typ   string // "path", "mode", "evidence", "action"
}

type analyzeModel struct {
	items    []menuItem
	selected int

	architectureFile string
	evidenceFiles    []string

	requestPicker bool
	pickerMode    pickerMode

	mode      string
	running   bool
	cancelled bool
	progress  float64
	stage     string
	result    *AnalysisResult
	engine    *Engine
	statusMsg string
}

func newAnalyzeModel(engine *Engine) analyzeModel {
	return analyzeModel{
		items: []menuItem{
			{label: "Architecture File", value: "", typ: "path"},
			{label: "Evidence Files", value: "0 files", typ: "evidence"},
			{label: "", value: "", typ: "sep"},
			{label: ModeASFOnly, value: ModeASFOnly, typ: "mode"},
			{label: ModeASFAndAI, value: ModeASFAndAI, typ: "mode"},
			{label: "", value: "", typ: "sep"},
			{label: "▶ Start Analysis", value: "", typ: "action"},
		},
		mode:   ModeASFOnly,
		engine: engine,
	}
}

type progressTickMsg struct{}

func (m *analyzeModel) Update(msg tea.Msg) (analyzeModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.running {
			switch msg.String() {
			case "esc":
				m.running = false
				m.cancelled = true
				m.statusMsg = "Analysis cancelled"
			}
			return *m, nil
		}

		switch msg.String() {
		case "up", "k":
			m.moveSel(-1)
		case "down", "j":
			m.moveSel(1)
		case "enter":
			return m.handleEnter()
		}

	case progressTickMsg:
		if !m.running {
			return *m, nil
		}
		if m.progress < 99 {
			m.progress += 10
			m.stage = analyzeStage(int(m.progress))
		}
		return *m, m.progressCmd()

	case errorMsg:
		m.running = false
		m.statusMsg = fmt.Sprintf("Error: %s", string(msg))
		return *m, nil
	}

	return *m, nil
}

func (m *analyzeModel) moveSel(dir int) {
	n := len(m.items)
	for i := 0; i < n; i++ {
		m.selected = (m.selected + dir + n) % n
		if m.items[m.selected].typ != "sep" {
			return
		}
	}
}

func (m *analyzeModel) setDocPath(path string) {
	m.architectureFile = path
	for i := range m.items {
		if m.items[i].typ == "path" && m.items[i].label == "Architecture File" {
			if path != "" {
				m.items[i].value = filepath.Base(path)
			} else {
				m.items[i].value = ""
			}
			return
		}
	}
}

func (m *analyzeModel) docPath() string {
	return m.architectureFile
}

func (m *analyzeModel) addEvidence(path string) {
	for _, existing := range m.evidenceFiles {
		if existing == path {
			return
		}
	}
	m.evidenceFiles = append(m.evidenceFiles, path)
	for i := range m.items {
		if m.items[i].typ == "evidence" {
			m.items[i].value = fmt.Sprintf("%d files", len(m.evidenceFiles))
			return
		}
	}
}

func (m *analyzeModel) removeEvidence(idx int) {
	if idx >= 0 && idx < len(m.evidenceFiles) {
		m.evidenceFiles = append(m.evidenceFiles[:idx], m.evidenceFiles[idx+1:]...)
		for i := range m.items {
			if m.items[i].typ == "evidence" {
				if len(m.evidenceFiles) == 0 {
					m.items[i].value = "0 files"
				} else {
					m.items[i].value = fmt.Sprintf("%d files", len(m.evidenceFiles))
				}
				return
			}
		}
	}
}

func (m *analyzeModel) clearEvidence() {
	m.evidenceFiles = nil
	for i := range m.items {
		if m.items[i].typ == "evidence" {
			m.items[i].value = "0 files"
			return
		}
	}
}

func (m *analyzeModel) evPath() string {
	if len(m.evidenceFiles) == 0 {
		return ""
	}
	return m.evidenceFiles[0]
}

func (m *analyzeModel) handleEnter() (analyzeModel, tea.Cmd) {
	item := m.items[m.selected]
	switch item.typ {
	case "path":
		m.requestPicker = true
		m.pickerMode = pickerArchitecture
	case "evidence":
		m.requestPicker = true
		m.pickerMode = pickerEvidence
	case "mode":
		m.mode = item.value
	case "action":
		return m.startAnalysis()
	}
	return *m, nil
}

func (m *analyzeModel) startAnalysis() (analyzeModel, tea.Cmd) {
	path := m.docPath()
	if path == "" {
		m.statusMsg = "Please select an architecture file first"
		return *m, nil
	}
	m.running = true
	m.cancelled = false
	m.progress = 0
	m.stage = "Initializing..."
	m.statusMsg = ""
	m.result = nil
	return *m, tea.Batch(m.progressCmd(), m.runAnalysisCmd())
}

func (m *analyzeModel) progressCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return progressTickMsg{}
	})
}

func (m *analyzeModel) runAnalysisCmd() tea.Cmd {
	return func() tea.Msg {
		progress := make(chan AnalysisProgress, 10)
		go func() {
			for range progress {
			}
		}()
		if m.engine == nil {
			close(progress)
			return errorMsg("engine not initialized")
		}
		result, err := m.engine.RunAnalysis(m.docPath(), m.evPath(), m.mode, progress)
		if err != nil {
			return errorMsg(err.Error())
		}
		return analysisCompleteMsg{result: result}
	}
}

func analyzeStage(pct int) string {
	switch {
	case pct < 10:
		return "Loading architecture document..."
	case pct < 20:
		return "Parsing document structure..."
	case pct < 30:
		return "Extracting security claims..."
	case pct < 40:
		return "Classifying claim types..."
	case pct < 50:
		return "Converting to assumptions..."
	case pct < 60:
		return "Applying STRIDE categorization..."
	case pct < 70:
		return "Verifying against evidence..."
	case pct < 80:
		return "Assessing risk levels..."
	case pct < 90:
		return "Building trust chains..."
	default:
		return "Generating gap analysis..."
	}
}

func (m mainModel) viewAnalyze() string {
	s := m.styles
	am := m.analyze

	if am.running {
		return m.viewAnalyzeProgress()
	}

	var rows []string
	for i, item := range am.items {
		if item.typ == "sep" {
			rows = append(rows, "")
			continue
		}

		display := item.value
		style := s.SectionItem
		prefix := "  "

		if i == am.selected {
			style = s.MenuSelected
			prefix = s.Fox.Render("▶ ")
		}

		switch item.typ {
		case "path":
			if display == "" {
				display = s.DimText.Render("[ Select Architecture File ]")
			}
			rows = append(rows, style.Render(fmt.Sprintf("%s%s: %s", prefix, item.label, display)))
		case "evidence":
			rows = append(rows, style.Render(fmt.Sprintf("%s%s: %s", prefix, item.label, display)))
			if i == am.selected && len(am.evidenceFiles) > 0 {
				for _, ef := range am.evidenceFiles {
					rows = append(rows, s.DimText.Render("    "+s.StatusGood.Render("✓")+" "+filepath.Base(ef)))
				}
			}
		case "mode":
			marker := ""
			if am.mode == item.value {
				marker = " " + s.StatusGood.Render("●")
			}
			rows = append(rows, style.Render(fmt.Sprintf("%s%s%s", prefix, item.label, marker)))
		case "action":
			rows = append(rows, style.Render(fmt.Sprintf("%s%s", prefix, item.label)))
		}
	}

	if len(am.evidenceFiles) > 0 {
		if am.selected < 0 || am.items[am.selected].typ != "evidence" {
			for _, ef := range am.evidenceFiles {
				rows = append(rows, s.DimText.Render("  "+s.StatusGood.Render("✓")+" "+filepath.Base(ef)))
			}
		}
	}

	if am.statusMsg != "" {
		rows = append(rows, "", s.CardAccent("", "  "+s.StatusWarn.Render(am.statusMsg), 40))
	}

	header := s.PremiumHeader("New Analysis", m.mainWidth())
	breadcrumb := s.Breadcrumb.Render("ASF0") +
		s.BreadcrumbSep.Render(" / ") +
		s.BreadcrumbSep.Render("New Analysis / ") +
		s.BreadcrumbSep.Render("Select Architecture")

	guidance := s.DimText.Render("  Select an architecture file (YAML, JSON, MD, DrawIO, SVG) to analyze.") +
		"\n" + s.DimText.Render("  Optionally add evidence files (CSV, JSON, YAML) for verification.")

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		breadcrumb,
		"",
		guidance,
		"",
		s.Card("Configuration", strings.Join(rows, "\n"), m.mainWidth()-4),
	)
}

func (m mainModel) viewAnalyzeProgress() string {
	s := m.styles

	if m.analyze.cancelled {
		return s.Card("", lipgloss.JoinVertical(lipgloss.Center,
			s.StatusWarn.Render("Analysis cancelled."),
			s.DimText.Render("Press Esc or select a different view."),
		), m.mainWidth()-4)
	}

	progress := m.analyze.progress
	if progress < 0 {
		progress = 0
	}
	if progress > 100 {
		progress = 100
	}

	fox := s.Fox.Render(FoxLogoSmall())

	return lipgloss.JoinVertical(lipgloss.Center,
		"",
		fox,
		"",
		s.PremiumHeader("Analysis in Progress", m.mainWidth()),
		"",
		s.Card("", lipgloss.JoinVertical(lipgloss.Center,
			"",
			s.DimText.Render(m.analyze.stage),
			"",
			"  "+s.ProgressWithLabel(progress, 40)+"  ",
			"",
			s.Accent.Render(fmt.Sprintf("%.0f%%", progress)),
		), m.mainWidth()-4),
		"",
		s.DimText.Render("  Press Esc to cancel  "),
	)
}

func (m mainModel) updateAnalyze(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.analyze, cmd = m.analyze.Update(msg)
	if m.analyze.requestPicker {
		m.analyze.requestPicker = false
		m.filePicker = newFilePickerState()
		m.filePicker.path = m.pickerStartPath(m.analyze.pickerMode)
		m.filePicker.mode = m.analyze.pickerMode
		m.filePicker.refresh()
		m.pickerActive = true
	}
	return m, cmd
}
