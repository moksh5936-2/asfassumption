package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ──────────────────────────────────────────────
// Card System
// ──────────────────────────────────────────────

func (s StyleSet) Card(title string, content string, width int) string {
	return s.cardWithStyle(s.CardBorder, title, content, width)
}

func (s StyleSet) CardAccent(title string, content string, width int) string {
	return s.cardWithStyle(s.CardBorderAccent, title, content, width)
}

func (s StyleSet) cardWithStyle(borderStyle lipgloss.Style, title string, content string, width int) string {
	var b strings.Builder
	bdr := borderStyle.GetBorderStyle()
	innerW := width - lipgloss.Width(bdr.Left) - lipgloss.Width(bdr.Right)
	if innerW < 10 {
		innerW = 10
	}

	b.WriteString(bdr.TopLeft + strings.Repeat(bdr.Top, innerW) + bdr.TopRight + "\n")

	if title != "" {
		titleStr := " " + s.CardTitle.Render(title) + " "
		pad := innerW - lipgloss.Width(titleStr)
		if pad < 0 {
			pad = 0
		}
		b.WriteString(bdr.Left + titleStr + strings.Repeat(" ", pad) + bdr.Right + "\n")
		b.WriteString(bdr.Left + strings.Repeat(" ", innerW) + bdr.Right + "\n")
	}

	for _, line := range strings.Split(content, "\n") {
		clean := lipgloss.Width(line)
		pad := innerW - clean
		if pad < 0 {
			pad = 0
		}
		b.WriteString(bdr.Left + line + strings.Repeat(" ", pad) + bdr.Right + "\n")
	}

	b.WriteString(bdr.BottomLeft + strings.Repeat(bdr.Bottom, innerW) + bdr.BottomRight)

	return s.CardContainer.Render(b.String())
}

func (s StyleSet) CardSimple(content string, width int) string {
	return s.Card("", content, width)
}

// ──────────────────────────────────────────────
// Premium Headers
// ──────────────────────────────────────────────

func (s StyleSet) PremiumHeader(title string, width int) string {
	rule := strings.Repeat("━", width-2)
	return fmt.Sprintf("┏%s┓\n┃%s┃\n┗%s┛",
		rule,
		s.PremiumTitle.Render(lipgloss.PlaceHorizontal(width-2, lipgloss.Center, title)),
		rule)
}

func (s StyleSet) SectionHeader(title string, width int) string {
	rule := strings.Repeat("━", width-2)
	return fmt.Sprintf("━━━ %s %s", s.SectionTitle.Render(title), s.SectionRule.Render(rule))
}

func (s StyleSet) SubHeader(title string) string {
	return " " + s.SubSectionTitle.Render(title)
}

// ──────────────────────────────────────────────
// Subtle Texture / Shadow
// ──────────────────────────────────────────────

func (s StyleSet) ShadowLine(width int) string {
	if width < 2 {
		return ""
	}
	return s.Shadow.Render(" " + strings.Repeat("░", width-2))
}

func (s StyleSet) TextureLabel(label string) string {
	return s.Texture.Render(" ░ " + label + " ░ ")
}

// ──────────────────────────────────────────────
// Risk Heatmap
// ──────────────────────────────────────────────

func (s StyleSet) RiskHeatmap(title string, critical, high, medium, low int) string {
	total := critical + high + medium + low
	if total == 0 {
		return ""
	}
	maxBar := 40
	var b strings.Builder
	if title != "" {
		b.WriteString(s.SubHeader(title) + "\n")
	}
	bar := func(count int, style lipgloss.Style) string {
		if count == 0 {
			return ""
		}
		n := count * maxBar / total
		if n < 1 && count > 0 {
			n = 1
		}
		return style.Render(strings.Repeat("█", n))
	}
	b.WriteString(fmt.Sprintf("  %s %s %d\n", s.BadgeCritical.Render("CRIT"), bar(critical, s.StatusBad), critical))
	b.WriteString(fmt.Sprintf("  %s %s %d\n", s.BadgeHigh.Render("HIGH"), bar(high, s.StatusWarn), high))
	b.WriteString(fmt.Sprintf("  %s %s %d\n", s.BadgeMedium.Render("MED"), bar(medium, s.StatusGood), medium))
	b.WriteString(fmt.Sprintf("  %s %s %d\n", s.BadgeLow.Render("LOW"), bar(low, s.Value), low))
	return b.String()
}

func (s StyleSet) RiskHeatmapCompact(critical, high, medium, low int) string {
	total := critical + high + medium + low
	if total == 0 {
		return ""
	}
	maxBar := 20
	bar := func(count int, style lipgloss.Style) string {
		if count == 0 {
			return ""
		}
		n := count * maxBar / total
		if n < 1 && count > 0 {
			n = 1
		}
		return style.Render(strings.Repeat("█", n))
	}
	return fmt.Sprintf("%s %s  %s %s  %s %s  %s %s",
		s.BadgeCritical.Render(fmt.Sprintf("%d", critical)), bar(critical, s.StatusBad),
		s.BadgeHigh.Render(fmt.Sprintf("%d", high)), bar(high, s.StatusWarn),
		s.BadgeMedium.Render(fmt.Sprintf("%d", medium)), bar(medium, s.StatusGood),
		s.BadgeLow.Render(fmt.Sprintf("%d", low)), bar(low, s.Value),
	)
}

// ──────────────────────────────────────────────
// Progress Bar
// ──────────────────────────────────────────────

func (s StyleSet) ProgressBarVisual(pct float64, width int) string {
	if width < 4 {
		width = 4
	}
	filled := int(float64(width-2) * pct / 100.0)
	if filled > width-2 {
		filled = width - 2
	}
	if filled < 0 {
		filled = 0
	}
	bar := s.ProgressFill.Render(strings.Repeat("█", filled))
	remain := s.ProgressEmpty.Render(strings.Repeat("▒", width-2-filled))
	return bar + remain
}

func (s StyleSet) ProgressWithLabel(pct float64, width int) string {
	bar := s.ProgressBarVisual(pct, width-6)
	return fmt.Sprintf("%s %3.0f%%", bar, pct)
}

// ──────────────────────────────────────────────
// Status Cards
// ──────────────────────────────────────────────

func (s StyleSet) StatusCardLarge(status string, count int, detail string, width int) string {
	var b strings.Builder
	icon := "✓"
	style := s.StatusGood
	switch status {
	case "VERIFIED":
		icon = "✓"
		style = s.StatusGood
	case "PARTIAL":
		icon = "~"
		style = s.StatusWarn
	case "UNVERIFIED":
		icon = "?"
		style = s.DimText
	case "NO EVIDENCE":
		icon = "○"
		style = s.DimText
	case "CRITICAL":
		icon = "!"
		style = s.StatusBad
	}
	b.WriteString(fmt.Sprintf("  %s %s", style.Render(icon), style.Render(status)))
	if count > 0 {
		b.WriteString(fmt.Sprintf("  %s", s.Badge.Render(fmt.Sprintf("%d", count))))
	}
	if detail != "" {
		b.WriteString("\n  " + s.DimText.Render(detail))
	}
	return s.Card("", b.String(), width)
}

// ──────────────────────────────────────────────
// Trust Chain Diagram
// ──────────────────────────────────────────────

func (s StyleSet) TrustDiagram(nodes []string) string {
	if len(nodes) == 0 {
		return ""
	}
	var b strings.Builder
	boxW := 14
	for i, node := range nodes {
		label := node
		if len(label) > boxW-2 {
			label = label[:boxW-4] + ".."
		}
		top := "╭" + strings.Repeat("─", boxW) + "╮"
		mid := "│ " + s.DiagramNode.Render(lipgloss.PlaceHorizontal(boxW-1, lipgloss.Center, label)) + "│"
		bot := "╰" + strings.Repeat("─", boxW) + "╯"
		b.WriteString(s.DiagramBox.Render(top + "\n" + mid + "\n" + bot))
		if i < len(nodes)-1 {
			b.WriteString("\n" + s.DiagramConnector.Render("    │") + "\n" + s.DiagramConnector.Render("    ▼") + "\n")
		}
	}
	return b.String()
}

func (s StyleSet) SPOFDiagram(node string, dependentCount int, risk string) string {
	boxW := 14
	label := node
	if len(label) > boxW-2 {
		label = label[:boxW-4] + ".."
	}
	riskStyle := s.StatusBad
	switch risk {
	case "High":
		riskStyle = s.StatusWarn
	case "Medium":
		riskStyle = s.StatusGood
	case "Low":
		riskStyle = s.Value
	}
	var b strings.Builder
	top := "╭" + strings.Repeat("─", boxW) + "╮"
	mid := "│ " + s.DiagramNode.Render(lipgloss.PlaceHorizontal(boxW-1, lipgloss.Center, label)) + "│"
	bot := "╰" + strings.Repeat("─", boxW) + "╯"
	b.WriteString(riskStyle.Render("      ▼") + "\n")
	b.WriteString(s.DiagramBox.Render(top+"\n"+mid+"\n"+bot) + "\n")
	b.WriteString(s.DiagramConnector.Render(strings.Repeat(" ", 5)+"│ │ │") + "\n")
	b.WriteString(s.DiagramConnector.Render(" ────┼─┼─┼────") + "\n")
	b.WriteString(s.DimText.Render(fmt.Sprintf("  %d dependents", dependentCount)) + "\n")
	b.WriteString(riskStyle.Render(fmt.Sprintf("  Risk: %s", risk)))
	return b.String()
}

// ──────────────────────────────────────────────
// Contradiction Side-by-Side
// ──────────────────────────────────────────────

func (s StyleSet) ContradictionView(claimA, claimB, result string, width int) string {
	half := (width - 6) / 2
	if half < 15 {
		half = 15
	}
	aBox := s.Card("Claim A", s.DimText.Render(claimA), half)
	bBox := s.Card("Claim B", s.DimText.Render(claimB), half)
	sideBySide := lipgloss.JoinHorizontal(lipgloss.Top, aBox, strings.Repeat(" ", 2), bBox)
	sep := s.RiskCritical.Render(strings.Repeat("═", width-4))
	resultCard := s.CardAccent("Result", s.RiskCritical.Render("  "+result), width-4)
	return lipgloss.JoinVertical(lipgloss.Left, sideBySide, "", sep, "", resultCard)
}

// ──────────────────────────────────────────────
// Review Card
// ──────────────────────────────────────────────

func (s StyleSet) ReviewCard(id, description, risk, status string, width int) string {
	innerW := width - 4
	if innerW < 20 {
		innerW = 20
	}
	riskBadge := s.RiskBadge(risk)
	statusMarker := ""
	statusStyle := s.Value
	switch status {
	case "Accepted":
		statusMarker = " ✓"
		statusStyle = s.StatusGood
	case "Rejected":
		statusMarker = " ✗"
		statusStyle = s.StatusBad
	case "Modified":
		statusMarker = " ~"
		statusStyle = s.StatusWarn
	default:
		statusMarker = " ?"
		statusStyle = s.DimText
	}
	top := fmt.Sprintf("%s  %s", riskBadge, statusStyle.Render(statusMarker))
	desc := description
	if len(desc) > innerW-4 {
		desc = desc[:innerW-7] + "..."
	}
	content := top + "\n" + s.DimText.Render(id) + "\n" + desc
	return s.CardAccent("", content, width)
}

// ──────────────────────────────────────────────
// Risk Badge
// ──────────────────────────────────────────────

func (s StyleSet) RiskBadge(risk string) string {
	switch risk {
	case "Critical":
		return s.BadgeCritical.Render("CRITICAL")
	case "High":
		return s.BadgeHigh.Render("HIGH")
	case "Medium":
		return s.BadgeMedium.Render("MEDIUM")
	case "Low":
		return s.BadgeLow.Render("LOW")
	default:
		return s.Badge.Render(risk)
	}
}

// ──────────────────────────────────────────────
// Loading / Branding
// ──────────────────────────────────────────────

func FoxLogoSmall() string {
	return ` /\_/\
( o.o )
 > ^ <
 ASF0`
}

func FoxLogoLarge() string {
	return `        ╭──────────────────────╮
        │   /\_/\               │
        │  ( o.o )   ASF0       │
        │   > ^ <    v` + ASFVersion + `          │
        │   Security Assumption  │
        │   Framework Zero      │
        ╰──────────────────────╯`
}

func FoxLogoCompact() string {
	return ` /\_/\  ASF0 v` + ASFVersion + `
( o.o ) Security Assumption
 > ^ <  Framework Zero`
}

func (s StyleSet) BrandedLoading(stage string, pct float64) string {
	fox := s.Fox.Render(FoxLogoSmall())
	bar := s.ProgressWithLabel(pct, 40)
	return lipgloss.JoinVertical(lipgloss.Center,
		fox,
		"",
		s.Title.Render("ASF0 Security Assumption Framework"),
		s.DimText.Render(stage),
		"",
		bar,
	)
}

// ──────────────────────────────────────────────
// Executive Summary Widget
// ──────────────────────────────────────────────

func (s StyleSet) ExecutiveSummaryWidget(critical, high, medium, low,
	verified, partial, unverified, chains, spofs int, width int) string {
	var b strings.Builder
	b.WriteString(s.SubHeader("Risk Distribution") + "\n")
	b.WriteString(s.RiskHeatmap("", critical, high, medium, low))
	b.WriteString("\n")
	b.WriteString(s.SubHeader("Verification Status") + "\n")
	b.WriteString(fmt.Sprintf("  %s  %s  %s",
		s.StatusGood.Render(fmt.Sprintf("✓%d", verified)),
		s.StatusWarn.Render(fmt.Sprintf("~%d", partial)),
		s.DimText.Render(fmt.Sprintf("?%d", unverified))))
	b.WriteString("\n\n")
	b.WriteString(s.SubHeader("Trust Analysis") + "\n")
	b.WriteString(fmt.Sprintf("  Chains: %d  SPOFs: %s", chains,
		s.RiskBadge(fmt.Sprintf("%d", spofs))))
	return s.Card("Executive Summary", b.String(), width)
}

// ──────────────────────────────────────────────
// Workflow Lifecycle Widget
// ──────────────────────────────────────────────

func (s StyleSet) renderWorkflow(r *AnalysisResult) string {
	steps := []struct {
		label string
		done  bool
	}{
		{"New Analysis", true},
		{"Run Analysis", true},
		{"Review", len(r.Assumptions) > 0},
		{"Validate", r.VerificationOutput != nil},
		{"Reports", false},
	}

	var lines []string
	for _, st := range steps {
		if st.done {
			lines = append(lines, fmt.Sprintf("  %s  %s", s.StatusGood.Render("✓"), s.DimText.Render(st.label)))
		} else {
			lines = append(lines, fmt.Sprintf("  %s  %s", s.DimText.Render("○"), s.DimText.Render(st.label)))
		}
	}

	flowStr := strings.Join(lines, "\n")
	return s.Card("Workflow", flowStr, 40)
}

// ──────────────────────────────────────────────
// Task Card (for validation queue)
// ──────────────────────────────────────────────

func (s StyleSet) TaskCard(id, description string, done bool, width int) string {
	innerW := width - 4
	if innerW < 10 {
		innerW = 10
	}
	check := "□"
	style := s.DimText
	if done {
		check = "■"
		style = s.StatusGood
	}
	desc := description
	if len(desc) > innerW-6 {
		desc = desc[:innerW-9] + "..."
	}
	content := fmt.Sprintf("%s %s  %s", style.Render(check), s.DimText.Render(id), desc)
	return s.Card("", content, width)
}
