package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type aboutModel struct{}

func newAboutModel() aboutModel {
	return aboutModel{}
}

func (m aboutModel) Update(msg tea.Msg) (aboutModel, tea.Cmd) {
	return m, nil
}

func (m mainModel) viewAbout() string {
	s := m.styles

	l := LoadLicense()
	licenseStr := "No license"
	if l != nil {
		if l.Valid {
			licenseStr = s.StatusGood.Render(l.Message)
		} else {
			licenseStr = s.StatusWarn.Render(l.Message)
		}
	}

	content := lipgloss.JoinVertical(lipgloss.Left,
		s.Title.Render("About ASF"),
		"",
		s.Section.Render("Architecture Security Framework"),
		s.SectionItem.Render(fmt.Sprintf("Version: v%s", ASFVersion)),
		s.SectionItem.Render("Security Assumption Discovery Engine"),
		"",
		s.Section.Render("License"),
		s.SectionItem.Render(licenseStr),
		"",
		s.Section.Render("Description"),
		s.SectionItem.Render("ASF automatically discovers hidden security assumptions"),
		s.SectionItem.Render("in system architecture diagrams and documents."),
		s.SectionItem.Render("Uses STRIDE threat modeling methodology and"),
		s.SectionItem.Render("automatic risk assessment to identify assumptions"),
		s.SectionItem.Render("that could lead to security vulnerabilities."),
		"",
		s.Section.Render("Technology"),
		s.SectionItem.Render("Built with Go, Bubble Tea, Lipgloss"),
		s.SectionItem.Render("Fully offline, no cloud dependency"),
		s.SectionItem.Render("Optional local AI enhancement layer"),
		"",
		s.Section.Render("Keyboard"),
		s.SectionItem.Render("↑↓/jk: Navigate    Enter: Select    Esc: Back"),
		s.SectionItem.Render("q: Quit from startup    Ctrl+C: Force quit"),
		"",
		s.SectionItem.Render(lipgloss.NewStyle().Foreground(s.Theme().DimText).Render("© 2026 ASF Project")),
	)

	return lipgloss.JoinVertical(lipgloss.Left,
		s.Title.Render("About"),
		s.BorderBox.Render(content),
	)
}

func (m mainModel) updateAbout(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.about, cmd = m.about.Update(msg)
	return m, cmd
}
