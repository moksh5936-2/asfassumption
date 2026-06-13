package main

import (
	"strings"
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
		{"Global", []string{
			"Ctrl+C / Q  Force quit",
			"q           Go back / navigate to previous view",
			"Esc         Go back / Cancel analysis",
			"?           Toggle help",
			"Tab         Cycle sidebar sections",
			"Shift+Tab   Cycle sidebar sections (reverse)",
			"f           Open file explorer",
			"r           Run analysis (open Analyze view)",
		}},
		{"Navigation", []string{
			"↑ / k       Move up / Scroll up",
			"↓ / j       Move down / Scroll down",
			"PgUp / b    Page up",
			"PgDn / Space Page down",
			"Home / g    Go to top",
			"End / G     Go to bottom",
			"Ctrl+U      Half page up",
			"Ctrl+D      Half page down",
		}},
		{"Dashboard", []string{
			"↑↓          Select quick action",
			"Enter       Open selected action",
			"a           Analyze architecture",
			"l           Open AI models",
			"s           Open settings",
			"i           Open about",
			"1-9         Open recent file (by number)",
		}},
		{"Analyze", []string{
			"↑↓          Select field (path, mode, start)",
			"Enter       Edit path / Select mode / Start analysis",
			"In path: type text, Enter to confirm",
			"Esc         Cancel running analysis",
			"f           Open file browser to pick document",
		}},
		{"Results", []string{
			"Tab         Next result tab",
			"Shift+Tab   Previous result tab",
			"/           Filter current tab by keyword",
			"n/N         Next/prev match in search",
			"e           Export results",
			"c           Clear results",
			"r           Open review mode",
			"v           Open validation mode",
		}},
		{"File Explorer", []string{
			"↑↓          Navigate files",
			"Enter       Open folder / Select file",
			"Backspace   Go to parent directory",
			".           Toggle hidden files",
			"Tab         Toggle preview panel",
			"/           Search filename",
		}},
		{"AI Models", []string{
			"↑↓          Select model",
			"Enter       Show actions for model",
			"Esc         Close action menu",
		}},
		{"Review Mode", []string{
			"Enter       Toggle browse/detail",
			"s           Accept assumption",
			"r           Reject assumption",
			"m           Mark as modified",
			"n           Edit note",
			"v           Open validation for this assumption",
		}},
		{"Export", []string{
			"↑↓          Select format",
			"Enter       Confirm selection",
			"Esc         Back to results",
		}},
		{"Settings", []string{
			"↑↓          Select setting",
			"Enter       Start editing",
			"← →         Change value",
			"s           Save settings",
			"Esc         Cancel edit",
		}},
		{"Search", []string{
			"/           Start search (current result tab)",
			"Esc / Enter Exit search",
		}},
	}

	var rows []string
	rows = append(rows, s.Title.Render("Help — Keyboard Shortcuts"))
	rows = append(rows, s.Subtitle.Render("Press ? to close help, Esc to go back"))
	rows = append(rows, "")

	for _, sec := range sections {
		rows = append(rows, s.Section.Render(sec.title))
		for _, k := range sec.keys {
			rows = append(rows, s.SectionItem.Render("  "+k))
		}
		rows = append(rows, "")
	}

	rows = append(rows, s.Section.Render("Supported File Types"))
	rows = append(rows, s.SectionItem.Render("  YAML (.yaml, .yml), JSON (.json), Markdown (.md)"))
	rows = append(rows, s.SectionItem.Render("  Mermaid (.mmd), Draw.io (.drawio), SVG (.svg)"))
	rows = append(rows, s.SectionItem.Render("  PDF (.pdf), DOCX (.docx), TXT (.txt)"))

	rows = append(rows, "")
	rows = append(rows, s.DimText.Render("  ASF v"+ASFVersion+" — Architecture Security Framework"))
	rows = append(rows, s.DimText.Render("  Local-first, offline-capable, no telemetry"))

	return strings.Join(rows, "\n")
}
