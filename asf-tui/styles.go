package main

import "github.com/charmbracelet/lipgloss"

type Theme struct {
	Name      lipgloss.Color
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
	SidebarBg lipgloss.Color

	RiskCritical lipgloss.Color
	RiskHigh     lipgloss.Color
	RiskMedium   lipgloss.Color
	RiskLow      lipgloss.Color
	RiskInfo     lipgloss.Color
	CardBg       lipgloss.Color
	Glow         lipgloss.Color
}

var Themes = map[string]Theme{
	"ASF0": {
		Name:      "ASF0",
		Primary:   lipgloss.Color("#E8590C"),
		Secondary: lipgloss.Color("#F59E0B"),
		Accent:    lipgloss.Color("#FBBF24"),
		Success:   lipgloss.Color("#10B981"),
		Warning:   lipgloss.Color("#F59E0B"),
		Error:     lipgloss.Color("#EF4444"),
		Info:      lipgloss.Color("#E8590C"),
		Text:      lipgloss.Color("#F0F0F0"),
		DimText:   lipgloss.Color("#7878A0"),
		Border:    lipgloss.Color("#2D2D44"),
		Highlight: lipgloss.Color("#F97316"),
		Bg:        lipgloss.Color("#0F0F1A"),
		SidebarBg: lipgloss.Color("#0A0A14"),

		RiskCritical: lipgloss.Color("#EF4444"),
		RiskHigh:     lipgloss.Color("#F97316"),
		RiskMedium:   lipgloss.Color("#F59E0B"),
		RiskLow:      lipgloss.Color("#10B981"),
		RiskInfo:     lipgloss.Color("#38BDF8"),
		CardBg:       lipgloss.Color("#15152A"),
		Glow:         lipgloss.Color("#E8590C"),
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

	Sidebar       lipgloss.Style
	SidebarItem   lipgloss.Style
	SidebarActive lipgloss.Style
	SidebarParent lipgloss.Style
	SidebarChild  lipgloss.Style
	HeaderBar     lipgloss.Style
	HintsBar      lipgloss.Style
	StatusBar     lipgloss.Style
	Tab           lipgloss.Style
	TabActive     lipgloss.Style
	ScrollHint    lipgloss.Style
	EmptyState    lipgloss.Style
	Badge         lipgloss.Style
	BadgeCritical lipgloss.Style
	BadgeHigh     lipgloss.Style
	BadgeMedium   lipgloss.Style
	BadgeLow      lipgloss.Style
	DimText       lipgloss.Style
	Accent        lipgloss.Style

	CardBorder       lipgloss.Style
	CardBorderAccent lipgloss.Style
	CardContainer    lipgloss.Style
	CardTitle        lipgloss.Style
	StartupBox       lipgloss.Style
	PremiumTitle     lipgloss.Style
	SectionTitle     lipgloss.Style
	SectionRule      lipgloss.Style
	SubSectionTitle  lipgloss.Style
	Shadow           lipgloss.Style
	Texture          lipgloss.Style
	ProgressFill     lipgloss.Style
	ProgressEmpty    lipgloss.Style
	DiagramNode      lipgloss.Style
	DiagramBox       lipgloss.Style
	DiagramConnector lipgloss.Style
	RiskCritical     lipgloss.Style
	RiskHigh         lipgloss.Style
	RiskMedium       lipgloss.Style
	RiskLow          lipgloss.Style
}

func NewStyles(t Theme) StyleSet {
	sb := lipgloss.Color("#131518")
	if t.SidebarBg != "" {
		sb = t.SidebarBg
	}
	return StyleSet{
		t: t,

		App: lipgloss.NewStyle().
			Background(t.Bg),

		Title: lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true).
			MarginBottom(1),

		Subtitle: lipgloss.NewStyle().
			Foreground(t.Secondary).
			Italic(true).
			MarginBottom(1),

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
			Italic(true),

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
			Foreground(t.Primary).
			Bold(true),

		Sidebar: lipgloss.NewStyle().
			Background(sb).
			Padding(0, 1).
			Width(26),

		SidebarItem: lipgloss.NewStyle().
			Background(sb).
			Foreground(t.DimText).
			Padding(0, 1).
			Width(24),

		SidebarActive: lipgloss.NewStyle().
			Background(t.Primary).
			Foreground(t.Bg).
			Bold(true).
			Padding(0, 1).
			Width(24),

		SidebarParent: lipgloss.NewStyle().
			Background(sb).
			Foreground(t.Text).
			Bold(true).
			Padding(0, 1).
			Width(24),

		SidebarChild: lipgloss.NewStyle().
			Background(sb).
			Foreground(t.DimText).
			Padding(0, 1).
			Width(22),

		HeaderBar: lipgloss.NewStyle().
			Background(t.Primary).
			Foreground(t.Bg).
			Bold(true).
			Padding(0, 1),

		HintsBar: lipgloss.NewStyle().
			Background(t.Border).
			Foreground(t.DimText).
			Padding(0, 1),

		StatusBar: lipgloss.NewStyle().
			Background(sb).
			Foreground(t.DimText).
			Padding(0, 1),

		Tab: lipgloss.NewStyle().
			Foreground(t.DimText).
			Padding(0, 1),

		TabActive: lipgloss.NewStyle().
			Foreground(t.Text).
			Background(t.Border).
			Bold(true).
			Padding(0, 1),

		ScrollHint: lipgloss.NewStyle().
			Foreground(t.DimText).
			Italic(true),

		EmptyState: lipgloss.NewStyle().
			Foreground(t.DimText).
			Italic(true).
			Padding(1, 2),

		Badge: lipgloss.NewStyle().
			Foreground(t.DimText).
			Padding(0, 1),

		BadgeCritical: lipgloss.NewStyle().
			Foreground(t.Bg).
			Background(t.Error).
			Bold(true).
			Padding(0, 1),

		BadgeHigh: lipgloss.NewStyle().
			Foreground(t.Bg).
			Background(t.Warning).
			Bold(true).
			Padding(0, 1),

		BadgeMedium: lipgloss.NewStyle().
			Foreground(t.Bg).
			Background(t.Accent).
			Bold(true).
			Padding(0, 1),

		BadgeLow: lipgloss.NewStyle().
			Foreground(t.Bg).
			Background(t.Secondary).
			Bold(true).
			Padding(0, 1),

		DimText: lipgloss.NewStyle().
			Foreground(t.DimText).
			Italic(true),

		Accent: lipgloss.NewStyle().
			Foreground(t.Accent).
			Bold(true),

		CardBorder: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(t.Primary).
			Padding(0, 1),

		CardBorderAccent: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(t.Accent).
			Padding(0, 1),

		CardContainer: lipgloss.NewStyle().
			MarginTop(1).
			MarginBottom(1),

		CardTitle: lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true),

		StartupBox: lipgloss.NewStyle().
			Border(lipgloss.DoubleBorder()).
			BorderForeground(t.Primary).
			Padding(1, 4).
			Margin(1, 2),

		PremiumTitle: lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true),

		SectionTitle: lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true),

		SectionRule: lipgloss.NewStyle().
			Foreground(t.Border),

		SubSectionTitle: lipgloss.NewStyle().
			Foreground(t.Accent).
			Bold(true).
			MarginTop(1).
			MarginBottom(1),

		Shadow: lipgloss.NewStyle().
			Foreground(t.Border),

		Texture: lipgloss.NewStyle().
			Foreground(t.Primary).
			Bold(true),

		ProgressFill: lipgloss.NewStyle().
			Foreground(t.Primary),

		ProgressEmpty: lipgloss.NewStyle().
			Foreground(t.Border),

		DiagramNode: lipgloss.NewStyle().
			Foreground(t.Text).
			Bold(true),

		DiagramBox: lipgloss.NewStyle().
			Foreground(t.Primary),

		DiagramConnector: lipgloss.NewStyle().
			Foreground(t.DimText),

		RiskCritical: lipgloss.NewStyle().
			Foreground(t.RiskCritical).
			Bold(true),

		RiskHigh: lipgloss.NewStyle().
			Foreground(t.RiskHigh).
			Bold(true),

		RiskMedium: lipgloss.NewStyle().
			Foreground(t.RiskMedium).
			Bold(true),

		RiskLow: lipgloss.NewStyle().
			Foreground(t.RiskLow).
			Bold(true),
	}
}

func (s StyleSet) ThemeName() string { return string(s.t.Name) }
func (s StyleSet) Theme() Theme      { return s.t }
