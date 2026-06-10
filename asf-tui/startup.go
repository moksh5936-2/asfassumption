package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type startupModel struct {
	choices  []string
	selected int
}

func newStartupModel() startupModel {
	return startupModel{
		choices: []string{
			"Analyze Architecture",
			"Results",
			"AI Settings",
			"Settings",
			"About",
			"Exit",
		},
	}
}

func foxArt() string {
	return `                      /^\/^\
                    _|__|  O|
           \/     /~     \_/ \
            \/   |  ||  |   |\
            /\   |  ||  |   |\
           /  \  |  ||  |   |/
      __  /   /  |_||__|   |_
  _  /  \/   /  /         /   \
 / \/  /\_  /  /         /    |
|  |  |  \/  /          |     |
|  |  |    \/            \    |
|  |  |     |             |   |
 \  \ |     |             |  /
  \  \|     |             | /
   |  |     |         __  |/
   |  |      \      /  \ |
   |  |       |    / __  |
    \  \      |   | /  | |
     \  \____/    |/  / /
      \_     |     /__/ /
        \__  |    |    |
           \_|    |____|
              |   |    |
              |   |    |
              |   |    |
              |   |    |
              |   |____|
              |   |    |
              |   |   ||
              |   |   ||
              |   |   ||
              |   |   ||
              |   |   ||
              |   |   ||
               \  |  /
                \ | /
                 \|/`
}

func foxArtSmall() string {
	return `    /\___/\
   /       \
  |  .   .  |
  \   ___   /    ASF
   |     |
   |  O  |
   |     |
   ------`
}

func (m startupModel) Update(msg tea.Msg) (startupModel, tea.Cmd) {
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
			switch m.selected {
			case 0:
				return m, func() tea.Msg { return navigateMsg{to: analyzeView} }
			case 1:
				return m, func() tea.Msg { return navigateMsg{to: resultsView} }
			case 2:
				return m, func() tea.Msg { return navigateMsg{to: localaiView} }
			case 3:
				return m, func() tea.Msg { return navigateMsg{to: settingsView} }
			case 4:
				return m, func() tea.Msg { return navigateMsg{to: aboutView} }
			case 5:
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m mainModel) viewStartup() string {
	s := m.styles

	fox := foxArtSmall()

	leftHalf := lipgloss.JoinVertical(lipgloss.Center,
		s.Fox.Render(fox),
		s.Subtitle.Render("Architecture Security Framework"),
		s.Subtitle.Render("Security Assumption Discovery Engine"),
	)

	var items []string
	for i, choice := range m.startup.choices {
		style := s.MenuItem
		if i == m.startup.selected {
			style = s.MenuSelected
		}
		prefix := "  "
		if i == m.startup.selected {
			prefix = "▸ "
		}
		items = append(items, style.Render(prefix+choice))
	}
	menu := lipgloss.JoinVertical(lipgloss.Left, items...)

	rightHalf := s.BorderBox.Render(menu)

	header := s.Title.Render("Welcome to ASF")

	titleStyle := lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width - 4)
	header = titleStyle.Render(header)

	content := lipgloss.JoinVertical(lipgloss.Center,
		header,
		lipgloss.JoinHorizontal(lipgloss.Top, leftHalf, rightHalf),
	)

	v := lipgloss.NewStyle().Align(lipgloss.Center).Width(m.width)
	return v.Render(content)
}

func (m mainModel) updateStartup(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.startup, cmd = m.startup.Update(msg)
	return m, cmd
}
