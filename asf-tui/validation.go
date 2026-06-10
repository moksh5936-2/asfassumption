package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ValidationMode is a developer-focused view that displays every assumption
// with its complete evidence, STRIDE justification, risk justification, and
// confidence breakdown. This exists solely to help evaluate and improve model
// quality — no additional functionality.

type validationModel struct {
	assumptions  []Assumption
	currentIdx   int
	mode         string // "browse", "detail"
}

func newValidationModel() validationModel {
	return validationModel{
		mode: "browse",
	}
}

func (m mainModel) viewValidation() string {
	s := m.styles
	v := m.validate

	if len(v.assumptions) == 0 {
		return lipgloss.JoinVertical(lipgloss.Left,
			s.Title.Render("Validation Mode"),
			s.Subtitle.Render("No assumptions to validate. Run an analysis first."),
		)
	}

	switch v.mode {
	case "browse":
		return v.renderValidationBrowse(s)
	case "detail":
		return v.renderValidationDetail(s)
	default:
		return "Unknown validation mode"
	}
}

func (v *validationModel) renderValidationBrowse(s StyleSet) string {
	header := s.Title.Render("Validation Mode")
	sub := s.Subtitle.Render("Developer evaluation view — every assumption shown with full traceability.")

	var items []string
	for i, a := range v.assumptions {
		prefix := "  "
		style := s.SectionItem
		if i == v.currentIdx {
			prefix = "▸ "
			style = s.MenuSelected
		}

		confPct := int(a.Confidence * 100)
		confStyle := s.Value
		if confPct >= 80 {
			confStyle = s.StatusGood
		} else if confPct >= 50 {
			confStyle = s.StatusWarn
		} else {
			confStyle = s.StatusBad
		}

		riskStyle := riskStyle(s, a.Risk)

		label := fmt.Sprintf("%s[%s] %s %s %s",
			prefix,
			a.ID,
			riskStyle.Render(padRight(string(a.Risk), 10)),
			truncateStr(a.Description, 50),
			confStyle.Render(fmt.Sprintf("(%d%% conf)", confPct)),
		)
		items = append(items, style.Render(label))
	}

	list := lipgloss.JoinVertical(lipgloss.Left, items...)

	return lipgloss.JoinVertical(lipgloss.Left,
		header, sub, "",
		s.BorderBox.Render(list), "",
		s.SectionItem.Render("↑↓ navigate | Enter detail | Esc back"),
	)
}

func (v *validationModel) renderValidationDetail(s StyleSet) string {
	if v.currentIdx >= len(v.assumptions) {
		return "No assumption selected."
	}
	a := v.assumptions[v.currentIdx]

	header := s.Title.Render(fmt.Sprintf("Validation: %s — %s", a.ID, a.Risk))

	dim := lipgloss.NewStyle().Foreground(s.Theme().DimText)
	sectionLabel := lipgloss.NewStyle().Foreground(s.Theme().Primary).Bold(true)
	confPct := int(a.Confidence * 100)
	confStyle := s.Value
	if confPct >= 80 {
		confStyle = s.StatusGood
	} else if confPct >= 50 {
		confStyle = s.StatusWarn
	} else {
		confStyle = s.StatusBad
	}

	var b strings.Builder

	// ── Assumption ──
	b.WriteString(sectionLabel.Render("Assumption") + "\n")
	b.WriteString(fmt.Sprintf("  %s\n\n", a.Description))

	// ── Evidence Traceability ──
	b.WriteString(sectionLabel.Render("Evidence") + "\n")
	if a.SourceNode != "" {
		b.WriteString(fmt.Sprintf("  Source Node: %s\n", dim.Render(a.SourceNode)))
	}
	if a.SourceLine > 0 {
		b.WriteString(fmt.Sprintf("  Source Line: %d\n", a.SourceLine))
	}
	for _, ev := range a.EvidenceSources {
		b.WriteString(fmt.Sprintf("  %s\n", dim.Render(ev)))
	}
	b.WriteString("\n")

	// ── STRIDE Mapping ──
	strideStrs := make([]string, len(a.Stride))
	for i, st := range a.Stride {
		strideStrs[i] = string(st)
	}
	b.WriteString(sectionLabel.Render("STRIDE Mapping") + "\n")
	b.WriteString(fmt.Sprintf("  %s\n", strings.Join(strideStrs, ", ")))
	if len(a.StrideJustifications) > 0 {
		for _, sj := range a.StrideJustifications {
			sjConf := int(sj.Confidence * 100)
			b.WriteString(fmt.Sprintf("  %s %s\n", s.StatusGood.Render(string(sj.Category)+":"), sj.Reason))
			if len(sj.MatchedKeywords) > 0 {
				b.WriteString(fmt.Sprintf("    Triggered by: %s\n", dim.Render(strings.Join(sj.MatchedKeywords, ", "))))
			}
			b.WriteString(fmt.Sprintf("    Confidence: %d%% — %s\n", sjConf, dim.Render(sj.ConfidenceReason)))
		}
	}
	b.WriteString("\n")

	// ── Risk Assessment ──
	b.WriteString(sectionLabel.Render("Risk Assessment") + "\n")
	b.WriteString(fmt.Sprintf("  Level: %s\n", riskStyle(s, a.Risk).Render(string(a.Risk))))
	if a.RiskJustification != nil {
		rj := a.RiskJustification
		b.WriteString(fmt.Sprintf("  Score: %d/25 (Likelihood %d × Impact %d)\n", rj.RiskScore, rj.Likelihood, rj.Impact))
		b.WriteString(fmt.Sprintf("  Reason: %s\n", dim.Render(rj.RiskReason)))
		b.WriteString(fmt.Sprintf("  Likelihood: %d/5 — %s\n", rj.Likelihood, dim.Render(rj.LikelihoodReason)))
		if len(rj.LikelihoodFactors) > 0 {
			for _, lf := range rj.LikelihoodFactors {
				b.WriteString(fmt.Sprintf("    • %s: %d — %s\n", lf.Factor, lf.Value, dim.Render(lf.Reason)))
			}
		}
		b.WriteString(fmt.Sprintf("  Impact: %d/5 — %s\n", rj.Impact, dim.Render(rj.ImpactReason)))
		if len(rj.ImpactFactors) > 0 {
			for _, ifa := range rj.ImpactFactors {
				b.WriteString(fmt.Sprintf("    • %s: %d — %s\n", ifa.Factor, ifa.Value, dim.Render(ifa.Reason)))
			}
		}
	}
	b.WriteString("\n")

	// ── Confidence ──
	b.WriteString(sectionLabel.Render("Confidence") + "\n")
	b.WriteString(fmt.Sprintf("  Overall: %s\n", confStyle.Render(fmt.Sprintf("%d%%", confPct))))
	if a.RiskJustification != nil && a.RiskJustification.ConfidenceReason != "" {
		b.WriteString(fmt.Sprintf("  Factors: %s\n", dim.Render(a.RiskJustification.ConfidenceReason)))
	}
	b.WriteString("\n")

	// ── Recommended Controls ──
	// (Controls are not stored per-assumption in the result, skip here)

	body := s.BorderBox.Render(b.String())

	return lipgloss.JoinVertical(lipgloss.Left,
		header, body, "",
		s.SectionItem.Render("Enter: Back to list | Esc: Exit validation"),
	)
}

func (m mainModel) updateValidation(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.validate.mode == "browse" && m.validate.currentIdx > 0 {
				m.validate.currentIdx--
			}
		case "down", "j":
			if m.validate.mode == "browse" && m.validate.currentIdx < len(m.validate.assumptions)-1 {
				m.validate.currentIdx++
			}
		case "enter":
			if m.validate.mode == "browse" {
				m.validate.mode = "detail"
			} else {
				m.validate.mode = "browse"
			}
		}
	}
	return m, nil
}
