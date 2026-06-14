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
	assumptions []Assumption
	currentIdx  int
	mode        string // "browse", "detail"
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
		return s.Card("Validation Queue",
			lipgloss.JoinVertical(lipgloss.Left,
				s.EmptyState.Render("No validations pending."),
				"",
				s.DimText.Render("  All findings have been reviewed."),
				"",
				s.DimText.Render("  Open a case from the sidebar and press 'v' to"),
				s.DimText.Render("  validate its assumptions against evidence."),
			),
			m.mainWidth())
	}

	switch v.mode {
	case "browse":
		return v.renderValidationBrowse(s, m.mainWidth()-4)
	case "detail":
		return v.renderValidationDetail(s, m.mainWidth()-4)
	default:
		return "Unknown validation mode"
	}
}

func (v *validationModel) renderValidationBrowse(s StyleSet, width int) string {
	header := s.PremiumHeader("Validation Queue", width+4)

	var items []string
	for i, a := range v.assumptions {
		confPct := int(a.Confidence * 100)
		done := confPct >= 80
		taskCard := s.TaskCard(a.ID, a.Description, done, width)
		if i == v.currentIdx {
			items = append(items, s.Fox.Render("▶ ")+taskCard)
		} else {
			items = append(items, "  "+taskCard)
		}
	}

	list := strings.Join(items, "\n")
	help := s.SectionRule.Render(strings.Repeat("─", max(1, width)))
	help += "\n" + s.DimText.Render("  ↑↓ Navigate  |  Enter Detail  |  Esc=Back")

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		list,
		"",
		help,
	)
}

func (v *validationModel) renderValidationDetail(s StyleSet, width int) string {
	if v.currentIdx >= len(v.assumptions) {
		return s.Card("", "No assumption selected.", width)
	}
	a := v.assumptions[v.currentIdx]

	confPct := int(a.Confidence * 100)
	confStyle := s.Value
	if confPct >= 80 {
		confStyle = s.StatusGood
	} else if confPct >= 50 {
		confStyle = s.StatusWarn
	} else {
		confStyle = s.StatusBad
	}

	var body []string

	body = append(body, "  "+s.SubSectionTitle.Render("Assumption"))
	body = append(body, "  "+s.DimText.Render(a.Description))
	body = append(body, "")

	body = append(body, "  "+s.SubSectionTitle.Render("Evidence Traceability"))
	if a.SourceNode != "" {
		body = append(body, "  Source Node: "+s.DimText.Render(a.SourceNode))
	}
	if a.SourceLine > 0 {
		body = append(body, "  Source Line: "+s.DimText.Render(fmt.Sprintf("%d", a.SourceLine)))
	}
	for _, ev := range a.EvidenceSources {
		body = append(body, "  • "+s.DimText.Render(ev))
	}
	body = append(body, "")

	strideStrs := make([]string, len(a.Stride))
	for i, st := range a.Stride {
		strideStrs[i] = string(st)
	}
	body = append(body, "  "+s.SubSectionTitle.Render("STRIDE Mapping"))
	body = append(body, "  "+s.Value.Render(strings.Join(strideStrs, ", ")))
	if len(a.StrideJustifications) > 0 {
		for _, sj := range a.StrideJustifications {
			sjConf := int(sj.Confidence * 100)
			body = append(body, "  "+s.StatusGood.Render(string(sj.Category)+":")+" "+sj.Reason)
			if len(sj.MatchedKeywords) > 0 {
				body = append(body, "    Triggered by: "+s.DimText.Render(strings.Join(sj.MatchedKeywords, ", ")))
			}
			body = append(body, "    Confidence: "+s.DimText.Render(fmt.Sprintf("%d%% — %s", sjConf, sj.ConfidenceReason)))
		}
	}
	body = append(body, "")

	body = append(body, "  "+s.SubSectionTitle.Render("Risk Assessment"))
	body = append(body, "  Level: "+riskStyle(s, a.Risk).Render(string(a.Risk)))
	if a.RiskJustification != nil {
		rj := a.RiskJustification
		body = append(body, fmt.Sprintf("  Score: %d/25 (Likelihood %d × Impact %d)", rj.RiskScore, rj.Likelihood, rj.Impact))
		body = append(body, "  Reason: "+s.DimText.Render(rj.RiskReason))
		body = append(body, fmt.Sprintf("  Likelihood: %d/5 — %s", rj.Likelihood, s.DimText.Render(rj.LikelihoodReason)))
		if len(rj.LikelihoodFactors) > 0 {
			for _, lf := range rj.LikelihoodFactors {
				body = append(body, "    • "+s.DimText.Render(fmt.Sprintf("%s: %d — %s", lf.Factor, lf.Value, lf.Reason)))
			}
		}
		body = append(body, fmt.Sprintf("  Impact: %d/5 — %s", rj.Impact, s.DimText.Render(rj.ImpactReason)))
		if len(rj.ImpactFactors) > 0 {
			for _, ifa := range rj.ImpactFactors {
				body = append(body, "    • "+s.DimText.Render(fmt.Sprintf("%s: %d — %s", ifa.Factor, ifa.Value, ifa.Reason)))
			}
		}
	}
	body = append(body, "")

	body = append(body, "  "+s.SubSectionTitle.Render("Confidence"))
	body = append(body, "  Overall: "+confStyle.Render(fmt.Sprintf("%d%%", confPct)))
	if a.RiskJustification != nil && a.RiskJustification.ConfidenceReason != "" {
		body = append(body, "  Factors: "+s.DimText.Render(a.RiskJustification.ConfidenceReason))
	}

	card := s.Card(fmt.Sprintf("Validation: %s — %s", a.ID, a.Risk), strings.Join(body, "\n"), width)
	help := s.SectionRule.Render(strings.Repeat("─", max(1, width)))
	help += "\n" + s.DimText.Render("  Enter=Back to list  |  Esc=Exit validation")

	return lipgloss.JoinVertical(lipgloss.Left,
		card,
		"",
		help,
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
