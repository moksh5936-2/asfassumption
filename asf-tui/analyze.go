package main

import (
	"fmt"
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
	typ   string // "path", "mode", "action"
}

type analyzeModel struct {
	items    []menuItem
	selected int

	inputMode string
	inputBuf  string

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
			{label: "Document Path", value: "", typ: "path"},
			{label: "Evidence Path", value: "", typ: "path"},
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

		if m.inputMode != "" {
			return m.handleTextInput(msg)
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

func (m *analyzeModel) handleTextInput(msg tea.KeyMsg) (analyzeModel, tea.Cmd) {
	switch msg.String() {
	case "enter":
		for i := range m.items {
			if m.items[i].typ == "path" && m.inputMode == m.items[i].label {
				m.items[i].value = m.inputBuf
				if m.inputMode == "Document Path" {
					m.docPath() // just to sync
				}
				break
			}
		}
		m.inputMode = ""
		m.inputBuf = ""
	case "esc":
		m.inputMode = ""
		m.inputBuf = ""
	case "backspace":
		if len(m.inputBuf) > 0 {
			m.inputBuf = m.inputBuf[:len(m.inputBuf)-1]
		}
	case "space":
		m.inputBuf += " "
	default:
		if len(msg.String()) == 1 {
			m.inputBuf += msg.String()
		}
	}
	return *m, nil
}

func (m *analyzeModel) setDocPath(path string) {
	for i := range m.items {
		if m.items[i].typ == "path" && m.items[i].label == "Document Path" {
			m.items[i].value = path
			return
		}
	}
}

func (m *analyzeModel) docPath() string {
	for _, it := range m.items {
		if it.typ == "path" && it.label == "Document Path" {
			return it.value
		}
	}
	return ""
}

func (m *analyzeModel) evPath() string {
	for _, it := range m.items {
		if it.typ == "path" && it.label == "Evidence Path" {
			return it.value
		}
	}
	return ""
}

func (m *analyzeModel) handleEnter() (analyzeModel, tea.Cmd) {
	item := m.items[m.selected]
	switch item.typ {
	case "path":
		m.inputMode = item.label
		m.inputBuf = item.value
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
		m.statusMsg = "Please set a document path first (select Document Path, press Enter)"
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
	case pct < 20:
		return "Parsing Documents..."
	case pct < 40:
		return "Extracting Claims..."
	case pct < 60:
		return "Converting to Assumptions..."
	case pct < 80:
		return "Verifying Against Evidence..."
	default:
		return "Generating Gap Analysis..."
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

		if i == am.selected && am.inputMode == item.label {
			editStr := fmt.Sprintf("  %s: %s█", item.label, am.inputBuf)
			rows = append(rows, s.MenuSelected.Render("▸ "+editStr))
			continue
		}

		display := item.value
		style := s.SectionItem
		prefix := "  "

		if i == am.selected {
			style = s.MenuSelected
			prefix = "▸ "
		}

		switch item.typ {
		case "path":
			if display == "" {
				display = "Select an architecture file to begin."
			}
			rows = append(rows, style.Render(fmt.Sprintf("%s%s: %s", prefix, item.label, display)))
		case "mode":
			marker := ""
			if am.mode == item.value {
				marker = s.StatusGood.Render(" ✓")
			}
			rows = append(rows, style.Render(fmt.Sprintf("%s%s%s", prefix, item.label, marker)))
		case "action":
			rows = append(rows, style.Render(fmt.Sprintf("%s%s", prefix, item.label)))
		}
	}

	if am.statusMsg != "" {
		rows = append(rows, "", s.StatusWarn.Render("  "+am.statusMsg))
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		s.Title.Render("Analyze Architecture"),
		s.Subtitle.Render("Set document path, select mode, then Start Analysis"),
		s.BorderBox.Render(lipgloss.JoinVertical(lipgloss.Left, rows...)),
	)
}

func (m mainModel) viewAnalyzeProgress() string {
	s := m.styles

	if m.analyze.cancelled {
		return lipgloss.JoinVertical(lipgloss.Center,
			s.Title.Render("Analyze Architecture"),
			s.StatusWarn.Render("Analysis cancelled."),
			s.SectionItem.Render("Press Esc or select a different view."),
		)
	}

	barWidth := 50
	if m.width > 80 {
		barWidth = m.width - 20
	}
	filled := int(float64(barWidth) * m.analyze.progress / 100.0)
	if filled > barWidth {
		filled = barWidth
	}
	bar := strings.Repeat("█", filled) + strings.Repeat("░", barWidth-filled)

	pctStr := fmt.Sprintf("%.0f%%", m.analyze.progress)

	return lipgloss.JoinVertical(lipgloss.Center,
		s.Title.Render("Analyzing Architecture"),
		s.Subtitle.Render("Press Esc to cancel"),
		s.BorderBox.Render(
			lipgloss.JoinVertical(lipgloss.Center,
				s.SectionItem.Render(m.analyze.stage),
				"",
				s.ProgressBar.Render(bar),
				s.Value.Render(pctStr),
			),
		),
	)
}

func (m mainModel) updateAnalyze(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.analyze, cmd = m.analyze.Update(msg)
	return m, cmd
}
