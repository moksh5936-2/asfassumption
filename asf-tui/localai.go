package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type localaiModel struct {
	models          []ModelInfo
	installedModels []string
	selected        int
	actionSelected  int
	showActions     bool
	downloading     bool
	downloadModel   string
	downloadProg    float64
	downloadStage   string
	downloadErr     string
	statusMsg       string
	modelMgr        *ModelManager
	dlProgress      *syncProgress
	config          *Config
}

func newLocalAIModel(cfg *Config) localaiModel {
	return localaiModel{
		models:          SupportedModels,
		installedModels: cfg.AI.InstalledModels,
		modelMgr:        NewModelManager(),
		config:          cfg,
	}
}

func (m localaiModel) Update(msg tea.Msg) (localaiModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.downloading {
				break
			}
			if m.showActions {
				if m.actionSelected > 0 {
					m.actionSelected--
				}
			} else {
				if m.selected > 0 {
					m.selected--
				}
			}
		case "down", "j":
			if m.downloading {
				break
			}
			if m.showActions {
				if m.actionSelected < len(m.actionsForModel())-1 {
					m.actionSelected++
				}
			} else {
				if m.selected < len(m.models)-1 {
					m.selected++
				}
			}
		case "enter":
			if m.downloading {
				break
			}
			if m.showActions {
				return m.executeAction()
			} else {
				m.showActions = true
				m.actionSelected = 0
			}
		case "esc":
			if m.showActions {
				m.showActions = false
			}
		}

	case aiDownloadTickMsg:
		if !m.downloading || m.dlProgress == nil {
			return m, nil
		}
		pct, stage, done, err := m.dlProgress.Snapshot()
		m.downloadProg = pct
		m.downloadStage = stage
		if err != "" {
			m.downloading = false
			m.downloadErr = err
			m.statusMsg = fmt.Sprintf("Download failed: %s", err)
			return m, nil
		}
		if done {
			m.downloading = false
			m.downloadProg = 100
			m.installedModels = append(m.installedModels, m.downloadModel)
			m.showActions = false
			m.statusMsg = fmt.Sprintf("✓ %s installed", m.downloadModel)
			return m, nil
		}
		return m, m.pollDownloadCmd()
	}
	return m, nil
}

type aiDownloadTickMsg struct{}

func (m localaiModel) pollDownloadCmd() tea.Cmd {
	return tea.Tick(400*time.Millisecond, func(t time.Time) tea.Msg {
		return aiDownloadTickMsg{}
	})
}

func (m localaiModel) actionsForModel() []string {
	actions := []string{"Download", "Set as Active"}
	modelName := m.models[m.selected].Name
	if contains(m.installedModels, modelName) {
		actions = append(actions, "Delete")
	}
	return actions
}

func (m localaiModel) executeAction() (localaiModel, tea.Cmd) {
	modelName := m.models[m.selected].Name
	switch m.actionSelected {
	case 0:
		if !m.modelMgr.CheckAvailable() {
			m.statusMsg = "Ollama not found. Install from https://ollama.ai"
			return m, nil
		}
		m.downloadModel = modelName
		m.downloading = true
		m.downloadProg = 0
		m.downloadStage = "Starting download..."
		m.downloadErr = ""

		sp := &syncProgress{}
		m.dlProgress = sp
		go m.modelMgr.StartDownload(modelName, sp)

		return m, m.pollDownloadCmd()
	case 1:
		if contains(m.installedModels, modelName) {
			m.config.AI.ActiveModel = modelName
			m.config.AI.Enabled = true
			m.config.Save(ConfigPath())
			m.statusMsg = fmt.Sprintf("✓ Active model set to %s", modelName)
		} else {
			m.statusMsg = fmt.Sprintf("%s is not installed", modelName)
		}
		m.showActions = false
	case 2:
		if err := m.modelMgr.DeleteModel(modelName); err != nil {
			m.statusMsg = fmt.Sprintf("Delete failed: %v", err)
		} else {
			m.installedModels = removeStr(m.installedModels, modelName)
			m.statusMsg = fmt.Sprintf("✓ %s deleted", modelName)
		}
		m.showActions = false
	}
	return m, nil
}

func contains(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func removeStr(slice []string, s string) []string {
	var result []string
	for _, v := range slice {
		if v != s {
			result = append(result, v)
		}
	}
	return result
}

func (m mainModel) viewLocalAI() string {
	s := m.styles
	lm := m.localai

	if lm.downloading {
		barWidth := 40
		filled := int(float64(barWidth) * lm.downloadProg / 100.0)
		bar := ""
		for i := 0; i < barWidth; i++ {
			if i < filled {
				bar += "█"
			} else {
				bar += "░"
			}
		}
		return lipgloss.JoinVertical(lipgloss.Left,
			s.Title.Render("Local AI Models"),
			s.BorderBox.Render(
				lipgloss.JoinVertical(lipgloss.Center,
					s.SectionItem.Render(fmt.Sprintf("Downloading: %s", lm.downloadModel)),
					s.ProgressBar.Render(bar),
					s.Value.Render(fmt.Sprintf("%.0f%%", lm.downloadProg)),
					s.SectionItem.Render(lm.downloadStage),
				),
			),
		)
	}

	var items []string
	for i, model := range lm.models {
		style := s.SectionItem
		prefix := "  "
		if !lm.showActions && i == lm.selected {
			style = s.MenuSelected
			prefix = "▸ "
		}
		status := ""
		if contains(lm.installedModels, model.Name) {
			status = s.StatusGood.Render(" ✓ Installed")
		}
		items = append(items, style.Render(fmt.Sprintf("%s%s  (%s)%s", prefix, model.Display, model.Size, status)))
	}

	modelList := lipgloss.JoinVertical(lipgloss.Left, items...)

	rows := []string{
		s.Title.Render("Local AI Models"),
		s.Subtitle.Render("Download and manage local AI models for enhanced analysis"),
		s.BorderBox.Render(modelList),
	}

	if lm.showActions && !lm.downloading {
		var actionItems []string
		actions := lm.actionsForModel()
		for i, action := range actions {
			style := s.SectionItem
			prefix := "  "
			if i == lm.actionSelected {
				style = s.MenuSelected
				prefix = "▸ "
			}
			actionItems = append(actionItems, style.Render(prefix+action))
		}
		rows = append(rows,
			"",
			s.Section.Render(fmt.Sprintf("Actions for %s:", lm.models[lm.selected].Display)),
			lipgloss.JoinVertical(lipgloss.Left, actionItems...),
		)
	}

	if lm.statusMsg != "" {
		rows = append(rows, "", s.StatusGood.Render("  "+lm.statusMsg))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m mainModel) updateLocalAI(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.localai, cmd = m.localai.Update(msg)
	return m, cmd
}
