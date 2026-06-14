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

	version := "v" + ASFVersion

	title := s.PremiumHeader("ASF0 Security Assumption Framework", m.mainWidth())

	details := fmt.Sprintf("  %s  %s\n", s.SectionTitle.Render("Version:"), s.Value.Render(version))
	details += fmt.Sprintf("  %s  %s\n", s.SectionTitle.Render("License:"), s.Value.Render("MIT"))
	details += fmt.Sprintf("  %s  %s\n", s.SectionTitle.Render("Engine:"), s.DimText.Render("Deterministic Risk Analysis"))
	details += fmt.Sprintf("  %s  %s\n", s.SectionTitle.Render("Risk Model:"), s.DimText.Render("5×5 Likelihood × Impact Matrix"))
	details += fmt.Sprintf("  %s  %s\n", s.SectionTitle.Render("Repository:"), s.DimText.Render("github.com/moksh5936-2/asfassumption"))

	fox := s.Fox.Render(FoxLogoLarge())
	if m.config.Appearance.FoxStyle == "Minimal" {
		fox = s.Fox.Render(FoxLogoCompact())
	} else if m.config.Appearance.FoxStyle == "None" {
		fox = ""
	}

	body := s.Card("About", fox+"\n\n"+details, m.mainWidth()-4)
	tagline := s.Card("",
		"  "+s.DimText.Render("Terminal-based security analysis workbench")+"\n"+
			"  "+s.DimText.Render("Architecture assumption validation, threat modeling, risk assessment")+"\n"+
			"  "+s.DimText.Render("Deterministic outputs • Evidence traceability • Review workflow"),
		m.mainWidth()-4)

	built := s.DimText.Render("  Built with Bubble Tea + Lipgloss  │  Q=Quit  ?=Help")

	return lipgloss.JoinVertical(lipgloss.Left,
		title, "",
		body, "",
		tagline, "",
		built,
	)
}

func (m mainModel) updateAbout(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.about, cmd = m.about.Update(msg)
	return m, cmd
}
