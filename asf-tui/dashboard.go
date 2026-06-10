package main

import (
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type dashboardModel struct {
	choices  []string
	selected int
}

func newDashboardModel() dashboardModel {
	return dashboardModel{
		choices: []string{
			"Analyze Architecture",
			"Local AI Models",
			"Settings",
			"About",
		},
	}
}

func (m dashboardModel) Update(msg tea.Msg) (dashboardModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.choices)-1 {
				m.selected++
			}
		case "enter":
			targets := []view{analyzeView, localaiView, settingsView, aboutView}
			if m.selected < len(targets) {
				return m, func() tea.Msg { return navigateMsg{to: targets[m.selected]} }
			}
		case "a":
			return m, func() tea.Msg { return navigateMsg{to: analyzeView} }
		case "l":
			return m, func() tea.Msg { return navigateMsg{to: localaiView} }
		case "s":
			return m, func() tea.Msg { return navigateMsg{to: settingsView} }
		case "i":
			return m, func() tea.Msg { return navigateMsg{to: aboutView} }
		}
	}
	return m, nil
}

func (m mainModel) viewDashboard() string {
	s := m.styles

	version := "v0.1.0"
	mode := m.config.Analysis.Depth
	aiStatus := "Offline"
	if m.config.AI.Enabled {
		aiStatus = "Active"
	}
	themeName := m.config.Appearance.Theme

	statusSection := lipgloss.JoinVertical(lipgloss.Left,
		s.Section.Render("System Status"),
		s.SectionItem.Render(fmtStatus(s, "Version", version)),
		s.SectionItem.Render(fmtStatus(s, "Mode", mode)),
		s.SectionItem.Render(fmtStatus(s, "AI", aiStatus)),
		s.SectionItem.Render(fmtStatus(s, "Theme", themeName)),
	)

	var items []string
	for i, choice := range m.dash.choices {
		style := s.MenuItem
		if i == m.dash.selected {
			style = s.MenuSelected
		}
		items = append(items, style.Render(choice))
	}
	navSection := lipgloss.JoinVertical(lipgloss.Left,
		s.Section.Render("Quick Actions"),
		lipgloss.JoinVertical(lipgloss.Left, items...),
	)

	body := lipgloss.JoinHorizontal(lipgloss.Top,
		s.BorderBox.Render(statusSection),
		lipgloss.NewStyle().Width(2).Render(""),
		s.BorderBox.Render(navSection),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		s.Title.Render("Dashboard"),
		body,
	)
}

func fmtStatus(s StyleSet, label, value string) string {
	return s.Label.Render(label) + s.Value.Render(value)
}

func (m mainModel) updateDashboard(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.dash, cmd = m.dash.Update(msg)
	return m, cmd
}
