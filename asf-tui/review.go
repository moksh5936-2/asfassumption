package main

import (
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ──────────────────────────────────────────────
// Review Model
// ──────────────────────────────────────────────

type reviewModel struct {
	assumptions    []Assumption
	currentIdx     int
	mode           string // "browse", "detail", "edit"
	editState      string // "status", "notes"
	statusOptions  []string
	selectedStatus int
	notesBuffer    string
	showValidation bool
	validationData []ValidationRecord
	editing        bool
	note           string
}

func newReviewModel() reviewModel {
	return reviewModel{
		statusOptions: []string{"Proposed", "Accepted", "Rejected", "Modified"},
		mode:          "browse",
	}
}

func (m mainModel) viewReview() string {
	s := m.styles
	rv := m.review

	if rv.mode == "" || len(rv.assumptions) == 0 {
		return s.Card("Review Queue",
			lipgloss.JoinVertical(lipgloss.Left,
				s.EmptyState.Render("No review items."),
				"",
				s.DimText.Render("  Analysis is fully reviewed."),
				"",
				s.DimText.Render("  Open a case from the sidebar and press 'r' to"),
				s.DimText.Render("  review its assumptions."),
			),
			m.mainWidth())
	}

	switch rv.mode {
	case "browse":
		return rv.renderBrowse(s, m.mainWidth()-4)
	case "detail":
		return rv.renderDetail(s, m.mainWidth()-4)
	default:
		return "Unknown review mode"
	}
}

func (rv *reviewModel) renderBrowse(s StyleSet, width int) string {
	header := s.PremiumHeader("Review Queue", width+4)

	var cards []string
	for i, a := range rv.assumptions {
		prefix := "  "
		if i == rv.currentIdx {
			prefix = s.Fox.Render("▶ ")
		}
		cardContent := s.ReviewCard(a.ID, a.Description, string(a.Risk), a.ReviewStatus, width)
		cards = append(cards, prefix+cardContent)
	}

	list := strings.Join(cards, "\n")

	help := s.SectionRule.Render(strings.Repeat("─", max(1, width)))
	help += "\n" + s.DimText.Render("  ↑↓ Navigate  |  Enter Detail  |  s=Accept  r=Reject  m=Modify  n=Notes  v=Validate  Esc=Back")

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		list,
		"",
		help,
	)
}

func (rv *reviewModel) renderDetail(s StyleSet, width int) string {
	if rv.currentIdx >= len(rv.assumptions) {
		return s.Card("", "No assumption selected.", width)
	}
	a := rv.assumptions[rv.currentIdx]

	statusStr := a.ReviewStatus
	if statusStr == "" {
		statusStr = "Proposed"
	}

	var bodyRows []string
	bodyRows = append(bodyRows, "  "+s.SubSectionTitle.Render("Description"))
	bodyRows = append(bodyRows, "  "+s.DimText.Render(a.Description))
	bodyRows = append(bodyRows, "")

	riskLine := fmt.Sprintf("  %s  L:%d  I:%d  Score:%d",
		s.RiskBadge(string(a.Risk)), a.Likelihood, a.Impact, riskScoreFromJust(a.RiskJustification))
	bodyRows = append(bodyRows, "  "+s.SubSectionTitle.Render("Risk"))
	bodyRows = append(bodyRows, "  "+riskLine)
	bodyRows = append(bodyRows, "")

	if len(a.Stride) > 0 {
		strideStrs := make([]string, len(a.Stride))
		for i, st := range a.Stride {
			strideStrs[i] = string(st)
		}
		bodyRows = append(bodyRows, "  "+s.SubSectionTitle.Render("STRIDE"))
		bodyRows = append(bodyRows, "  "+s.Value.Render(strings.Join(strideStrs, ", ")))
		bodyRows = append(bodyRows, "")
	}

	statusStyle := s.Value
	switch a.ReviewStatus {
	case "Accepted":
		statusStyle = s.StatusGood
	case "Rejected":
		statusStyle = s.StatusBad
	case "Modified":
		statusStyle = s.StatusWarn
	}
	bodyRows = append(bodyRows, "  "+s.SubSectionTitle.Render("Status"))
	if rv.editing {
		bodyRows = append(bodyRows, "  "+statusStyle.Render(statusStr)+"  "+s.DimText.Render("[EDITING] "+rv.note+"█"))
	} else {
		bodyRows = append(bodyRows, "  "+statusStyle.Render(statusStr))
		if a.ReviewNotes != "" {
			bodyRows = append(bodyRows, "  "+s.DimText.Render("Notes: "+a.ReviewNotes))
		}
	}
	bodyRows = append(bodyRows, "")

	if len(a.EvidenceSources) > 0 {
		bodyRows = append(bodyRows, "  "+s.SubSectionTitle.Render("Evidence"))
		for _, ev := range a.EvidenceSources {
			bodyRows = append(bodyRows, "  • "+s.DimText.Render(ev))
		}
		bodyRows = append(bodyRows, "")
	}

	if a.Rationale != "" {
		bodyRows = append(bodyRows, "  "+s.SubSectionTitle.Render("Rationale"))
		bodyRows = append(bodyRows, "  "+s.DimText.Render(a.Rationale))
		bodyRows = append(bodyRows, "")
	}

	if a.RiskJustification != nil {
		bodyRows = append(bodyRows, "  "+s.SubSectionTitle.Render("Risk Justification"))
		rj := a.RiskJustification
		bodyRows = append(bodyRows, fmt.Sprintf("  Likelihood %d/5: %s", rj.Likelihood, s.DimText.Render(rj.LikelihoodReason)))
		bodyRows = append(bodyRows, fmt.Sprintf("  Impact %d/5: %s", rj.Impact, s.DimText.Render(rj.ImpactReason)))
		bodyRows = append(bodyRows, fmt.Sprintf("  %s → %s",
			s.Value.Render(fmt.Sprintf("Score: %d/25", rj.RiskScore)),
			s.RiskBadge(string(rj.RiskLevel))))
	}

	card := s.Card(fmt.Sprintf("Review: %s", a.ID), strings.Join(bodyRows, "\n"), width)
	help := s.SectionRule.Render(strings.Repeat("─", max(1, width)))
	help += "\n" + s.DimText.Render("  s=Accept  r=Reject  m=Modify  n=Edit Note  Enter=Back  Esc=Cancel")

	return lipgloss.JoinVertical(lipgloss.Left,
		card, "",
		help,
	)
}

func (m mainModel) updateReview(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.review.mode == "browse" && m.review.currentIdx > 0 {
				m.review.currentIdx--
			}
		case "down", "j":
			if m.review.mode == "browse" && m.review.currentIdx < len(m.review.assumptions)-1 {
				m.review.currentIdx++
			}
		case "enter":
			if m.review.mode == "browse" {
				m.review.mode = "detail"
			} else {
				m.review.mode = "browse"
			}
		case "s":
			if m.review.mode == "detail" {
				m.review.assumptions[m.review.currentIdx].ReviewStatus = "Accepted"
				m.review.assumptions[m.review.currentIdx].ReviewTimestamp = time.Now()
			}
		case "r":
			if m.review.mode == "detail" {
				m.review.assumptions[m.review.currentIdx].ReviewStatus = "Rejected"
				m.review.assumptions[m.review.currentIdx].ReviewTimestamp = time.Now()
			}
		case "m":
			if m.review.mode == "detail" {
				m.review.assumptions[m.review.currentIdx].ReviewStatus = "Modified"
				m.review.assumptions[m.review.currentIdx].ReviewTimestamp = time.Now()
			}
		case "n":
			if m.review.mode == "detail" {
				if !m.review.editing {
					m.review.editing = true
					m.review.note = m.review.assumptions[m.review.currentIdx].ReviewNotes
				} else {
					m.review.editing = false
					if m.review.currentIdx >= 0 && m.review.currentIdx < len(m.review.assumptions) {
						m.review.assumptions[m.review.currentIdx].ReviewNotes = m.review.note
					}
				}
			}
		case "backspace":
			if m.review.mode == "detail" && m.review.editing && len(m.review.note) > 0 {
				m.review.note = m.review.note[:len(m.review.note)-1]
			}
		case "esc":
			if m.review.editing {
				m.review.editing = false
			}
		default:
			if m.review.mode == "detail" && m.review.editing && len(msg.String()) == 1 {
				m.review.note += msg.String()
			}
		case "v":
			if m.review.mode == "browse" {
				m.review.showValidation = !m.review.showValidation
			}
		}
	}
	return m, nil
}

// ──────────────────────────────────────────────
// Validation Data Export
// ──────────────────────────────────────────────

// CollectValidationData gathers all assumptions into validation records.
func CollectValidationData(assumptions []Assumption) []ValidationRecord {
	var records []ValidationRecord
	for _, a := range assumptions {
		rec := ValidationRecord{
			AssumptionID:      a.ID,
			Description:       a.Description,
			GeneratedEvidence: a.EvidenceSources,
			AssignedRisk:      a.Risk,
			Confidence:        a.Confidence,
			STRIDECategories:  a.Stride,
		}
		if a.RiskJustification != nil {
			rec.RiskScore = a.RiskJustification.RiskScore
		}
		if a.ReviewStatus != "" {
			rec.ArchReviewResult = a.ReviewStatus
			rec.ArchNotes = a.ReviewNotes
			rec.ReviewTimestamp = a.ReviewTimestamp
		}
		records = append(records, rec)
	}
	return records
}

// ──────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────

func riskScoreFromJust(rj *RiskJustification) int {
	if rj == nil {
		return 0
	}
	return rj.RiskScore
}

func truncateStr(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}
