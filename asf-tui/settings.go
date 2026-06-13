package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type settingsModel struct {
	settings []settingItem
	selected int
	editing  bool
	config   *Config
}

type settingItem struct {
	label    string
	key      string
	value    string
	values   []string
	valueIdx int
}

func newSettingsModel(cfg *Config) settingsModel {
	themeIdx := 0
	themes := []string{"Dark", "Midnight", "Cyber", "Minimal"}
	for i, t := range themes {
		if t == cfg.Appearance.Theme {
			themeIdx = i
			break
		}
	}

	depthIdx := 0
	depths := []string{"light", "standard", "deep"}
	for i, d := range depths {
		if d == cfg.Analysis.Depth {
			depthIdx = i
			break
		}
	}

	riskIdx := 0
	riskLevels := []string{"low", "medium", "high", "critical"}
	for i, r := range riskLevels {
		if r == cfg.Analysis.RiskThreshold {
			riskIdx = i
			break
		}
	}

	return settingsModel{
		config: cfg,
		settings: []settingItem{
			{label: "Theme", key: "theme", value: themes[themeIdx], values: themes, valueIdx: themeIdx},
			{label: "Fox Style", key: "fox_style", value: cfg.Appearance.FoxStyle, values: []string{"Classic", "Minimal", "None"}, valueIdx: foxIdx(cfg.Appearance.FoxStyle)},
			{label: "Analysis Depth", key: "depth", value: depths[depthIdx], values: depths, valueIdx: depthIdx},
			{label: "Risk Threshold", key: "risk_threshold", value: riskLevels[riskIdx], values: riskLevels, valueIdx: riskIdx},
			{label: "STRIDE Analysis", key: "stride", value: boolStr(cfg.Analysis.Stride), values: []string{"true", "false"}, valueIdx: boolIdx(cfg.Analysis.Stride)},
			{label: "Controls Check", key: "controls", value: boolStr(cfg.Analysis.Controls), values: []string{"true", "false"}, valueIdx: boolIdx(cfg.Analysis.Controls)},
			{label: "Default Export", key: "export", value: cfg.Output.Default, values: []string{"json", "markdown", "html", "csv", "pdf", "narrative-md", "narrative-html"}, valueIdx: exportIdx(cfg.Output.Default)},
			{label: "Export Directory", key: "export_dir", value: cfg.Output.Directory, values: nil, valueIdx: 0},
			{label: "AI Enhancement", key: "ai_enabled", value: boolStr(cfg.AI.Enabled), values: []string{"false", "true"}, valueIdx: boolIdx(cfg.AI.Enabled)},
			{label: "Active Model", key: "active_model", value: cfg.AI.ActiveModel, values: nil, valueIdx: 0},
			{label: "Debug Logging", key: "debug", value: boolStr(cfg.General.Debug), values: []string{"false", "true"}, valueIdx: boolIdx(cfg.General.Debug)},
			{label: "", key: "", value: "", values: nil, valueIdx: 0},
			{label: "Reset to Defaults", key: "reset", value: "", values: []string{"no", "yes"}, valueIdx: 0},
		},
	}
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

func boolIdx(b bool) int {
	if b {
		return 1
	}
	return 0
}

func foxIdx(s string) int {
	switch s {
	case "Classic":
		return 0
	case "Minimal":
		return 1
	case "None":
		return 2
	}
	return 0
}

func exportIdx(s string) int {
	switch s {
	case "json":
		return 0
	case "markdown":
		return 1
	case "html":
		return 2
	case "csv":
		return 3
	case "pdf":
		return 4
	case "narrative-md":
		return 5
	case "narrative-html":
		return 6
	}
	return 1
}

func (m settingsModel) Update(msg tea.Msg) (settingsModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if !m.editing && m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if !m.editing && m.selected < len(m.settings)-1 {
				m.selected++
			}
		case "enter":
			item := &m.settings[m.selected]
			if item.values == nil {
				break
			}
			if !m.editing {
				m.editing = true
			} else {
				m.editing = false
				m.applyChange(m.selected)
			}
		case "left", "h":
			if m.editing {
				item := &m.settings[m.selected]
				if item.values != nil && item.valueIdx > 0 {
					item.valueIdx--
					item.value = item.values[item.valueIdx]
				}
			}
		case "right", "l":
			if m.editing {
				item := &m.settings[m.selected]
				if item.values != nil && item.valueIdx < len(item.values)-1 {
					item.valueIdx++
					item.value = item.values[item.valueIdx]
				}
			}
		case "esc":
			if m.editing {
				m.editing = false
			}
		case "s":
			if !m.editing {
				m.config.Save(ConfigPath())
			}
		}
	}
	return m, nil
}

func (m *settingsModel) applyChange(idx int) {
	item := m.settings[idx]
	switch item.key {
	case "theme":
		m.config.Appearance.Theme = item.value
	case "fox_style":
		m.config.Appearance.FoxStyle = item.value
	case "depth":
		m.config.Analysis.Depth = item.value
	case "risk_threshold":
		m.config.Analysis.RiskThreshold = item.value
	case "stride":
		m.config.Analysis.Stride = item.value == "true"
	case "controls":
		m.config.Analysis.Controls = item.value == "true"
	case "export":
		m.config.Output.Default = item.value
	case "export_dir":
		m.config.Output.Directory = item.value
	case "ai_enabled":
		m.config.AI.Enabled = item.value == "true"
	case "active_model":
		m.config.AI.ActiveModel = item.value
	case "debug":
		m.config.General.Debug = item.value == "true"
	case "reset":
		if item.value == "yes" {
			def := DefaultConfig()
			*m.config = def
			m.config.Save(ConfigPath())
			*m = newSettingsModel(m.config)
		}
		return // return early to avoid double save
	}
	m.config.Save(ConfigPath())
}

func (m mainModel) viewSettings() string {
	s := m.styles

	header := lipgloss.JoinVertical(lipgloss.Left,
		s.Title.Render("Settings"),
		s.Subtitle.Render("Press Enter to edit a value, ←→ to change it, S to save"),
	)

	var rows []string
	for i, item := range m.settings.settings {
		if item.key == "" && item.label == "" {
			rows = append(rows, "")
			continue
		}
		labelStr := item.label
		valueStr := item.value

		editMarker := ""
		if m.settings.editing && i == m.settings.selected {
			if item.values != nil {
				accentStyle := lipgloss.NewStyle().Foreground(s.Theme().Accent)
				editMarker = accentStyle.Render(" < >")
			}
		}

		selectedPrefix := "  "
		if i == m.settings.selected {
			selectedPrefix = "▸ "
		}

		row := fmt.Sprintf("%s%s: %s%s", selectedPrefix, labelStr, valueStr, editMarker)
		if i == m.settings.selected {
			rows = append(rows, s.MenuSelected.Render(row))
		} else {
			rows = append(rows, s.SectionItem.Render(row))
		}
	}

	settingsList := lipgloss.JoinVertical(lipgloss.Left, rows...)

	saveHint := lipgloss.NewStyle().Foreground(s.Theme().DimText).Render("Press 's' to save settings")

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		s.BorderBox.Render(settingsList),
		saveHint,
	)
}

func (m mainModel) updateSettings(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.settings, cmd = m.settings.Update(msg)
	m.styles = NewStyles(Themes[m.config.Appearance.Theme])
	return m, cmd
}
