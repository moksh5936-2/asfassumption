package coverage

import (
	"encoding/json"
	"fmt"
	"strings"
)

func ExportMarkdown(output *CoverageOutput) string {
	if output == nil {
		return ""
	}
	var b strings.Builder

	b.WriteString("# Coverage & Blind Spot Analysis\n\n")
	b.WriteString(fmt.Sprintf("**Domain:** %s  \n", output.Domain))
	b.WriteString(fmt.Sprintf("**Generated:** %s  \n", output.GeneratedAt))

	if output.Assessment != nil && len(output.Assessment.Categories) > 0 {
		b.WriteString("## Coverage by Category\n\n")
		b.WriteString("| Category | Expected | Observed | Coverage | Risk |\n")
		b.WriteString("|----------|----------|----------|----------|------|\n")
		for _, cat := range output.Assessment.Categories {
			b.WriteString(fmt.Sprintf("| %s | %d | %d | %.1f%% | %s |\n",
				cat.Category, cat.ExpectedCount, cat.ObservedCount, cat.CoveragePct, cat.Risk))
		}
		b.WriteString("\n")

		if len(output.Assessment.Gaps) > 0 {
			b.WriteString("## Coverage Gaps\n\n")
			b.WriteString("| Category | Expected | Observed | Missing | Coverage | Risk | Recommendation |\n")
			b.WriteString("|----------|----------|----------|---------|----------|------|----------------|\n")
			for _, gap := range output.Assessment.Gaps {
				b.WriteString(fmt.Sprintf("| %s | %d | %d | %d | %.1f%% | %s | %s |\n",
					gap.Category, gap.ExpectedCount, gap.ObservedCount,
					gap.MissingCount, gap.CoveragePct, gap.Risk, gap.Recommendation))
			}
			b.WriteString("\n")
		}
	}

	if len(output.BlindSpots) > 0 {
		b.WriteString("## Blind Spots\n\n")
		for _, bs := range output.BlindSpots {
			b.WriteString(fmt.Sprintf("### %s [%s] (Score: %.0f)\n\n", bs.Title, bs.Risk, bs.Score))
			b.WriteString(fmt.Sprintf("**Category:** %s  \n", bs.Category))
			b.WriteString(fmt.Sprintf("**Component:** %s  \n", bs.Component))
			b.WriteString(fmt.Sprintf("**Description:** %s  \n", bs.Description))
			b.WriteString(fmt.Sprintf("**Trust Chain Impact:** %s  \n", bs.TrustChainImpact))
			b.WriteString(fmt.Sprintf("**Severity:** %s  \n", bs.ConsequenceSeverity))
			b.WriteString(fmt.Sprintf("**Compliance:** %s  \n", bs.ComplianceRelevance))
			b.WriteString(fmt.Sprintf("**Recommendation:** %s  \n\n", bs.Recommendation))
		}
	}

	if len(output.DomainBlindSpots) > 0 {
		b.WriteString("## Domain-Aware Blind Spots\n\n")
		b.WriteString("| Missing Area | Risk | Description | Recommendation |\n")
		b.WriteString("|--------------|------|-------------|----------------|\n")
		for _, ds := range output.DomainBlindSpots {
			b.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n",
				ds.MissingArea, ds.Risk, ds.Description, ds.Recommendation))
		}
		b.WriteString("\n")
	}

	if output.CISOView != nil {
		b.WriteString("## CISO Summary\n\n")

		if len(output.CISOView.DangerousMissingAssumptions) > 0 {
			b.WriteString("### Most Dangerous Missing Assumptions\n\n")
			for _, bs := range output.CISOView.DangerousMissingAssumptions {
				b.WriteString(fmt.Sprintf("- **%s** (%s) — %s\n", bs.Title, bs.Risk, bs.Description))
			}
			b.WriteString("\n")
		}

		if len(output.CISOView.AreasRequiringReview) > 0 {
			b.WriteString("### Areas Requiring Review\n\n")
			for _, area := range output.CISOView.AreasRequiringReview {
				b.WriteString(fmt.Sprintf("- %s\n", area))
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}

func ExportHTML(output *CoverageOutput) string {
	if output == nil {
		return ""
	}
	md := ExportMarkdown(output)

	html := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>Coverage & Blind Spot Analysis</title>
<style>
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,sans-serif;max-width:900px;margin:40px auto;padding:0 20px;line-height:1.6;color:#333;}
h1{color:#1a1a2e;border-bottom:2px solid #1a1a2e;padding-bottom:8px;}
h2{color:#16213e;margin-top:32px;}
h3{color:#0f3460;}
table{border-collapse:collapse;width:100%%;margin:16px 0;}
th,td{border:1px solid #ddd;padding:10px 12px;text-align:left;}
th{background-color:#1a1a2e;color:white;}
tr:nth-child(even){background-color:#f5f5f5;}
.score-good{color:#27ae60;font-weight:bold;}
.score-warn{color:#f39c12;font-weight:bold;}
.score-bad{color:#e74c3c;font-weight:bold;}
.risk-critical{color:#e74c3c;font-weight:bold;}
.risk-high{color:#e67e22;font-weight:bold;}
.risk-medium{color:#f39c12;}
.risk-low{color:#27ae60;}
</style>
</head>
<body>
%s
</body>
</html>`, mdToHTML(md))

	return html
}

func ExportJSON(output *CoverageOutput) ([]byte, error) {
	if output == nil {
		return []byte("{}"), nil
	}
	return json.MarshalIndent(output, "", "  ")
}

func mdToHTML(md string) string {
	lines := strings.Split(md, "\n")
	var html strings.Builder
	inTable := false
	inList := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "# ") {
			html.WriteString(fmt.Sprintf("<h1>%s</h1>\n", strings.TrimPrefix(trimmed, "# ")))
		} else if strings.HasPrefix(trimmed, "## ") {
			html.WriteString(fmt.Sprintf("<h2>%s</h2>\n", strings.TrimPrefix(trimmed, "## ")))
		} else if strings.HasPrefix(trimmed, "### ") {
			html.WriteString(fmt.Sprintf("<h3>%s</h3>\n", strings.TrimPrefix(trimmed, "### ")))
		} else if strings.HasPrefix(trimmed, "| ") {
			if !inTable {
				inTable = true
				html.WriteString("<table>\n")
			}
			cols := strings.Split(trimmed, "|")
			if strings.Contains(trimmed, "---") {
				continue
			}
			isHeader := strings.Contains(trimmed, "Category") || (len(cols) > 1 && strings.TrimSpace(cols[1]) != "")
			if isHeader && !strings.Contains(trimmed, "---") {
				html.WriteString("<thead><tr>")
				for _, col := range cols {
					c := strings.TrimSpace(col)
					if c != "" {
						html.WriteString(fmt.Sprintf("<th>%s</th>", c))
					}
				}
				html.WriteString("</tr></thead>\n")
			} else {
				html.WriteString("<tr>")
				for _, col := range cols {
					c := strings.TrimSpace(col)
					if c != "" {
						html.WriteString(fmt.Sprintf("<td>%s</td>", c))
					}
				}
				html.WriteString("</tr>\n")
			}
		} else {
			if inTable {
				html.WriteString("</table>\n")
				inTable = false
			}
			if strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") {
				if !inList {
					html.WriteString("<ul>\n")
					inList = true
				}
				html.WriteString(fmt.Sprintf("<li>%s</li>\n", strings.TrimPrefix(strings.TrimPrefix(trimmed, "- "), "* ")))
			} else {
				if inList {
					html.WriteString("</ul>\n")
					inList = false
				}
				if trimmed == "" {
					html.WriteString("<br>\n")
				} else {
					html.WriteString(fmt.Sprintf("<p>%s</p>\n", trimmed))
				}
			}
		}
	}
	if inTable {
		html.WriteString("</table>\n")
	}
	if inList {
		html.WriteString("</ul>\n")
	}

	return html.String()
}
