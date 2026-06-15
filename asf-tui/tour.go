package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type tourModel struct {
	step       int
	totalSteps int
}

func newTourModel() tourModel {
	return tourModel{
		step:       0,
		totalSteps: 7,
	}
}

type tourSlide struct {
	title   string
	icon    string
	content string
}

func tourSlides() []tourSlide {
	return []tourSlide{
		{
			icon: "📄", title: "Import Architecture",
			content: `ASF0 reads YAML architecture files that describe your
system's components, connections, and security claims.

A good architecture file includes:
  • Components — services, databases, APIs, users
  • Connections — who talks to whom and over what protocol
  • Security claims — "TLS enforced", "MFA required", "RBAC"

ASF0 extracts every assumption from your file automatically.

Supported formats: YAML · JSON · Markdown · DrawIO · SVG`,
		},
		{
			icon: "📎", title: "Add Evidence",
			content: `Evidence files contain real-world data that ASF0 uses to
verify or challenge assumptions.

Examples:
  • Network scan results showing actual TLS versions
  • IAM policy exports showing who has what access
  • Audit logs showing authentication events
  • Configuration files showing encryption settings

ASF0 cross-references claims against evidence to find gaps.

Evidence formats: CSV · JSON · YAML`,
		},
		{
			icon: "⚡", title: "Run Analysis",
			content: `ASF0 runs 12 parallel analysis engines:

  • Assumption Extraction   — finds every claim in your file
  • Contradiction Engine    — detects logical conflicts
  • Trust Chain Mapper      — maps who-trusts-whom
  • SDRI                    — Security Design Review Intelligence
  • Threat Modeler          — identifies attack paths
  • CIARE                   — Compliance audit prep (168 questions)
  • Coverage Analysis       — finds blind spots and gaps

Most analyses complete in 1–5 seconds.`,
		},
		{
			icon: "📊", title: "Review Findings",
			content: `Results appear in tabbed case workspaces.

  Tab 0  Overview      — Summary of all findings
  Tab 1  Assumptions   — Every claim extracted
  Tab 2  Verification  — Evidence-backed status check
  Tab 3  Contradictions — Logical conflicts (key output!)
  Tab 4  Trust         — Trust chains + SPOFs
  Tab 5  Controls      — Recommended mitigations
  Tab 6  SDRI          — Design review intelligence

Use  ← →  to switch tabs. Use  ↑ ↓  to select items.`,
		},
		{
			icon: "⚔️", title: "Contradictions",
			content: `A contradiction is a logical conflict between two claims.

Example:
  "MFA is enforced for all user authentication"
  vs.
  "Service accounts use password only with no MFA"

ASF0 found 21 such contradictions in the healthcare
benchmark. These are the most valuable output — they
reveal hidden risks that a manual reviewer might miss.

The Contradictions tab (Tab 3) shows every conflict
with severity, affected assumptions, and evidence.`,
		},
		{
			icon: "🔗", title: "Trust Chains & SPOFs",
			content: `ASF0 maps every trust relationship in your architecture.

  API → Database → Backup → Storage → Vendor

It then identifies:
  • Single Points of Trust Failure (SPOFs)
    — one component whose compromise breaks everything
  • Failure Cascades
    — how a single failure propagates through the system
  • Collapse Simulations
    — what happens when critical assumptions fail

A typical 24-component architecture generates ~100 trust
chains and 15+ single points of failure.`,
		},
		{
			icon: "🚀", title: "Export & Next Steps",
			content: `Once your analysis is complete:

  • Press  e  to export reports (JSON, Markdown, HTML, PDF)
  • Press  r  to review and approve/reject assumptions
  • Press  v  to validate assumptions against evidence
  • Use the WORK section in the sidebar for queues

Sample workflow:
  1. Import YAML architecture
  2. Run analysis (1–5 seconds)
  3. Review contradictions (Tab 3)
  4. Examine trust chains (Tab 4)
  5. Export report (press  e)
  6. Fix your architecture and re-run

Press  q  to return to the startup screen.`,
		},
	}
}

func (m mainModel) viewTour() string {
	s := m.styles
	slides := tourSlides()
	step := m.tour.step
	slide := slides[step]

	header := s.PremiumHeader(
		fmt.Sprintf("🦊 Quick Tour (%d/%d)", step+1, len(slides)),
		m.mainWidth(),
	)

	title := s.SubSectionTitle.Render(fmt.Sprintf("  %s  %s", slide.icon, slide.title))

	var bodyLines []string
	for _, line := range strings.Split(slide.content, "\n") {
		bodyLines = append(bodyLines, "  "+s.Value.Render(line))
	}

	var navHint string
	if step > 0 {
		navHint += s.DimText.Render("  ←  Prev")
	}
	if step < len(slides)-1 {
		navHint += s.Accent.Render("  →  Next  ")
	} else {
		navHint += s.DimText.Render("        ")
	}
	navHint += s.DimText.Render("  q  Back to Start")

	scrollBar := fmt.Sprintf("%s  Slide %d of %d",
		s.DimText.Render("▐█▌"), step+1, len(slides))

	sep := s.SectionRule.Render(strings.Repeat("━", max(1, m.mainWidth()-4)))

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		"",
		s.Card("", lipgloss.JoinVertical(lipgloss.Left,
			title,
			"",
			strings.Join(bodyLines, "\n"),
		), m.mainWidth()-4),
		"",
		sep,
		"",
		lipgloss.JoinHorizontal(lipgloss.Left,
			navHint,
			strings.Repeat(" ", max(1, m.mainWidth()-lipgloss.Width(navHint)-lipgloss.Width(scrollBar)-6)),
			scrollBar,
		),
	)
}
