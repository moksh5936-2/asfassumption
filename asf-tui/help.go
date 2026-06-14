package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type helpModel struct{}

func newHelpModel() helpModel {
	return helpModel{}
}

func (m mainModel) viewHelp() string {
	s := m.styles

	sections := []struct {
		title string
		keys  []string
	}{
		{"Navigation", []string{
			"Tab         Toggle sidebar / content focus",
			"↑↓          Navigate content or sidebar tree",
			"Enter       Select item or activate sidebar entry",
			"q / Esc     Go back",
			"Q           Quit",
		}},
		{"New Analysis", []string{
			"Enter       Select / browse architecture file",
			"Esc         Cancel running analysis",
		}},
		{"Case Workspace", []string{
			"↑↓ / j k    Scroll workspace tabs",
			"← →         Switch workspace tab",
			"r           Open Review mode",
			"v           Open Validation",
			"e           Generate Reports",
			"c           Clear case",
			"/           Search within visible content",
		}},
		{"Review Mode", []string{
			"↑↓          Navigate assumptions",
			"Enter       Toggle browse / detail view",
			"s           Accept assumption",
			"r           Reject assumption",
			"m           Mark as modified",
			"n           Edit notes",
			"v           Open Validation",
		}},
		{"Validation", []string{
			"↑↓          Navigate assumptions",
			"Enter       View detailed trace",
		}},
		{"Settings", []string{
			"Enter       Toggle edit mode",
			"← →         Change value during edit",
			"s           Save",
			"Esc         Cancel edit",
		}},
		{"Reports", []string{
			"↑↓          Select export format",
			"Enter       Choose / confirm",
			"Esc         Cancel / back",
		}},
		{"File Picker", []string{
			"↑↓ / j k    Navigate files",
			"Enter       Open directory / select file",
			"Backspace   Go to parent directory",
			".           Toggle hidden files",
			"Tab         Toggle preview panel",
			"/           Search files",
			"Esc         Cancel / close picker",
		}},
		{"Sidebar Tree", []string{
			"CASES",
			"  ➕ New Analysis",
			"  📁 <case files>",
			"WORK",
			"  📋 Review Queue",
			"  ✓ Validation Queue",
			"  📦 Reports",
			"SYSTEM",
			"  ⚙ Settings",
			"  ❓ Help",
			"  ℹ About",
		}},
	}

	header := s.PremiumHeader("Help", m.mainWidth())

	var rows []string
	rows = append(rows, s.DimText.Render(" Press Tab to focus sidebar │ ↑↓ to navigate │ Enter to select"))
	rows = append(rows, "")

	for _, sec := range sections {
		rows = append(rows, s.SectionHeader(sec.title, m.mainWidth()))
		for _, k := range sec.keys {
			rows = append(rows, "  "+s.DimText.Render(k))
		}
		rows = append(rows, "")
	}

	rows = append(rows, s.SectionRule.Render(strings.Repeat("━", max(1, m.mainWidth()))))
	rows = append(rows, s.DimText.Render("  ASF0 v"+ASFVersion+" — Security Assumption Framework"))

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		s.Card("", strings.Join(rows, "\n"), m.mainWidth()-4),
	)
}
