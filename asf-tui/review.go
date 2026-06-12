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
		return lipgloss.JoinVertical(lipgloss.Left,
			s.Title.Render("Review Mode"),
			s.Subtitle.Render("No assumptions to review. Run an analysis first."),
		)
	}

	switch rv.mode {
	case "browse":
		return rv.renderBrowse(s)
	case "detail":
		return rv.renderDetail(s)
	default:
		return "Unknown review mode"
	}
}

func (rv *reviewModel) renderBrowse(s StyleSet) string {
	header := s.Title.Render("Architect Review")
	sub := s.Subtitle.Render("Select an assumption to review. Press Enter for details, R to mark status.")

	var items []string
	for i, a := range rv.assumptions {
		prefix := "  "
		style := s.SectionItem
		if i == rv.currentIdx {
			prefix = "▸ "
			style = s.MenuSelected
		}

		statusMarker := ""
		statusStyle := s.Value
		switch a.ReviewStatus {
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
		}

		label := fmt.Sprintf("%s%s [%s] %s%s", prefix, a.ID, a.Risk, truncateStr(a.Description, 60), statusStyle.Render(statusMarker))
		items = append(items, style.Render(label))
	}

	list := lipgloss.JoinVertical(lipgloss.Left, items...)

	help := s.SectionItem.Render("↑↓ navigate | Enter detail | R toggle status | Esc back")

	return lipgloss.JoinVertical(lipgloss.Left,
		header, sub, "",
		s.BorderBox.Render(list), "",
		help,
	)
}

func (rv *reviewModel) renderDetail(s StyleSet) string {
	if rv.currentIdx >= len(rv.assumptions) {
		return "No assumption selected."
	}
	a := rv.assumptions[rv.currentIdx]

	header := s.Title.Render(fmt.Sprintf("Review: %s", a.ID))

	statusStr := a.ReviewStatus
	if statusStr == "" {
		statusStr = "Proposed"
	}

	detail := fmt.Sprintf("Description: %s\n", a.Description)
	detail += fmt.Sprintf("Risk: %s (L:%d I:%d Score:%d)\n", a.Risk, a.Likelihood, a.Impact, riskScoreFromJust(a.RiskJustification))
	if len(a.Stride) > 0 {
		strideStrs := make([]string, len(a.Stride))
		for i, s := range a.Stride {
			strideStrs[i] = string(s)
		}
		detail += fmt.Sprintf("STRIDE: %s\n", strings.Join(strideStrs, ", "))
	}
	detail += fmt.Sprintf("Status: %s\n", statusStr)

	if rv.editing {
		detail += fmt.Sprintf("Notes: [EDITING] %s█\n", rv.note)
	} else if a.ReviewNotes != "" {
		detail += fmt.Sprintf("Notes: %s\n", a.ReviewNotes)
	}

	if len(a.EvidenceSources) > 0 {
		detail += "\nEvidence:\n"
		for _, ev := range a.EvidenceSources {
			detail += fmt.Sprintf("  %s\n", ev)
		}
	}

	if a.Rationale != "" {
		detail += fmt.Sprintf("\nRationale: %s\n", a.Rationale)
	}

	if a.RiskJustification != nil {
		detail += fmt.Sprintf("\nRisk Justification:\n")
		detail += fmt.Sprintf("  Likelihood %d/5: %s\n", a.RiskJustification.Likelihood, a.RiskJustification.LikelihoodReason)
		detail += fmt.Sprintf("  Impact %d/5: %s\n", a.RiskJustification.Impact, a.RiskJustification.ImpactReason)
		detail += fmt.Sprintf("  Score: %d/25 → %s\n", a.RiskJustification.RiskScore, a.RiskJustification.RiskLevel)
	}

	body := s.BorderBox.Render(s.SectionItem.Render(detail))

	help := s.SectionItem.Render("S:Accept | R:Reject | M:Modified | N:Edit note | Enter back | Esc:cancel edit")

	return lipgloss.JoinVertical(lipgloss.Left,
		header, body, "",
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
