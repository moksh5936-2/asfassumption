package main

import (
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	Name      string
	Primary   lipgloss.Color
	Secondary lipgloss.Color
	Accent    lipgloss.Color
	Success   lipgloss.Color
	Warning   lipgloss.Color
	Error     lipgloss.Color
	Info      lipgloss.Color
	Text      lipgloss.Color
	DimText   lipgloss.Color
	Border    lipgloss.Color
	Highlight lipgloss.Color
	Bg        lipgloss.Color
}

var Themes = map[string]Theme{
	"Dark": {
		Name:      "Dark",
		Primary:   lipgloss.Color("#7C3AED"),
		Secondary: lipgloss.Color("#10B981"),
		Accent:    lipgloss.Color("#F59E0B"),
		Success:   lipgloss.Color("#22C55E"),
		Warning:   lipgloss.Color("#EAB308"),
		Error:     lipgloss.Color("#EF4444"),
		Info:      lipgloss.Color("#3B82F6"),
		Text:      lipgloss.Color("#E2E8F0"),
		DimText:   lipgloss.Color("#64748B"),
		Border:    lipgloss.Color("#334155"),
		Highlight: lipgloss.Color("#6366F1"),
		Bg:        lipgloss.Color("#1E293B"),
	},
	"Midnight": {
		Name:      "Midnight",
		Primary:   lipgloss.Color("#818CF8"),
		Secondary: lipgloss.Color("#34D399"),
		Accent:    lipgloss.Color("#F472B6"),
		Success:   lipgloss.Color("#22C55E"),
		Warning:   lipgloss.Color("#FBBF24"),
		Error:     lipgloss.Color("#F87171"),
		Info:      lipgloss.Color("#60A5FA"),
		Text:      lipgloss.Color("#C7D2FE"),
		DimText:   lipgloss.Color("#6B7280"),
		Border:    lipgloss.Color("#374151"),
		Highlight: lipgloss.Color("#A78BFA"),
		Bg:        lipgloss.Color("#111827"),
	},
	"Cyber": {
		Name:      "Cyber",
		Primary:   lipgloss.Color("#00FF41"),
		Secondary: lipgloss.Color("#00BFFF"),
		Accent:    lipgloss.Color("#FF00FF"),
		Success:   lipgloss.Color("#00FF41"),
		Warning:   lipgloss.Color("#FFFF00"),
		Error:     lipgloss.Color("#FF0000"),
		Info:      lipgloss.Color("#00BFFF"),
		Text:      lipgloss.Color("#00FF41"),
		DimText:   lipgloss.Color("#008F20"),
		Border:    lipgloss.Color("#004D14"),
		Highlight: lipgloss.Color("#FF00FF"),
		Bg:        lipgloss.Color("#0D0D0D"),
	},
	"Minimal": {
		Name:      "Minimal",
		Primary:   lipgloss.Color("#FFFFFF"),
		Secondary: lipgloss.Color("#888888"),
		Accent:    lipgloss.Color("#AAAAAA"),
		Success:   lipgloss.Color("#66BB6A"),
		Warning:   lipgloss.Color("#FFD54F"),
		Error:     lipgloss.Color("#E57373"),
		Info:      lipgloss.Color("#64B5F6"),
		Text:      lipgloss.Color("#FFFFFF"),
		DimText:   lipgloss.Color("#666666"),
		Border:    lipgloss.Color("#444444"),
		Highlight: lipgloss.Color("#FFFFFF"),
		Bg:        lipgloss.Color("#1A1A1A"),
	},
}

type StyleSet struct {
	t Theme

	App          lipgloss.Style
	Title        lipgloss.Style
	Subtitle     lipgloss.Style
	MenuItem     lipgloss.Style
	MenuSelected lipgloss.Style
	Label        lipgloss.Style
	Value        lipgloss.Style
	StatusGood   lipgloss.Style
	StatusWarn   lipgloss.Style
	StatusBad    lipgloss.Style
	BorderBox    lipgloss.Style
	Help         lipgloss.Style
	ErrorText    lipgloss.Style
	Button       lipgloss.Style
	ButtonFocus  lipgloss.Style
	Progress     lipgloss.Style
	ProgressBar  lipgloss.Style
	Section      lipgloss.Style
	SectionItem  lipgloss.Style
	Header       lipgloss.Style
	Fox          lipgloss.Style
}

func NewStyles(t Theme) StyleSet {
	return StyleSet{
		t: t,

		App: lipgloss.NewStyle().
			Background(t.Bg).
			Padding(0, 2),

		Title: lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true).
			MarginTop(1).
			MarginBottom(1),

		Subtitle: lipgloss.NewStyle().
			Foreground(t.Secondary).
			Italic(true),

		MenuItem: lipgloss.NewStyle().
			Foreground(t.Text).
			Padding(0, 2).
			MarginTop(1),

		MenuSelected: lipgloss.NewStyle().
			Foreground(t.Bg).
			Background(t.Primary).
			Bold(true).
			Padding(0, 2).
			MarginTop(1),

		Label: lipgloss.NewStyle().
			Foreground(t.DimText).
			Width(20).
			Align(lipgloss.Right).
			Padding(0, 1),

		Value: lipgloss.NewStyle().
			Foreground(t.Text),

		StatusGood: lipgloss.NewStyle().
			Foreground(t.Success).
			Bold(true),

		StatusWarn: lipgloss.NewStyle().
			Foreground(t.Warning).
			Bold(true),

		StatusBad: lipgloss.NewStyle().
			Foreground(t.Error).
			Bold(true),

		BorderBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(t.Border).
			Padding(1, 2).
			MarginTop(1).
			MarginBottom(1),

		Help: lipgloss.NewStyle().
			Foreground(t.DimText).
			Italic(true).
			MarginTop(1).
			MarginBottom(1),

		ErrorText: lipgloss.NewStyle().
			Foreground(t.Error).
			Bold(true),

		Button: lipgloss.NewStyle().
			Foreground(t.Text).
			Background(t.Border).
			Padding(0, 3).
			MarginTop(1).
			MarginRight(1).
			Align(lipgloss.Center),

		ButtonFocus: lipgloss.NewStyle().
			Foreground(t.Bg).
			Background(t.Primary).
			Bold(true).
			Padding(0, 3).
			MarginTop(1).
			MarginRight(1).
			Align(lipgloss.Center),

		Progress: lipgloss.NewStyle().
			Foreground(t.Accent).
			Width(40),

		ProgressBar: lipgloss.NewStyle().
			Background(t.Border),

		Section: lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true).
			MarginTop(1).
			Underline(true),

		SectionItem: lipgloss.NewStyle().
			Foreground(t.Text).
			Padding(0, 2).
			MarginTop(1),

		Header: lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true).
			MarginBottom(1),

		Fox: lipgloss.NewStyle().
			Foreground(t.Secondary).
			Bold(true),
	}
}

func (s StyleSet) ThemeName() string { return s.t.Name }
func (s StyleSet) Theme() Theme      { return s.t }
