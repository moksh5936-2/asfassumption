package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-pdf/fpdf"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ExportFormat string

const (
	ExportJSON     ExportFormat = "json"
	ExportMarkdown ExportFormat = "markdown"
	ExportCSV      ExportFormat = "csv"
	ExportPDF      ExportFormat = "pdf"
	ExportHTML     ExportFormat = "html"
)

func ExportResult(result *AnalysisResult, format ExportFormat, outputDir string) (string, error) {
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", fmt.Errorf("cannot create output directory: %w", err)
	}

	timestamp := time.Now().Format("20060102_150405")
	safeName := strings.ReplaceAll(result.ArchitectureName, ".", "_")
	baseName := fmt.Sprintf("%s_%s", safeName, timestamp)

	var path string
	var err error

	switch format {
	case ExportJSON:
		path = filepath.Join(outputDir, baseName+".json")
		err = exportJSON(result, path)
	case ExportMarkdown:
		path = filepath.Join(outputDir, baseName+".md")
		err = exportMarkdown(result, path)
	case ExportCSV:
		path = filepath.Join(outputDir, baseName+".csv")
		err = exportCSV(result, path)
	case ExportPDF:
		path = filepath.Join(outputDir, baseName+".pdf")
		err = exportPDF(result, path)
	case ExportHTML:
		path = filepath.Join(outputDir, baseName+".html")
		err = exportHTML(result, path)
	}

	if err != nil {
		return "", err
	}
	return path, nil
}

func exportJSON(result *AnalysisResult, path string) error {
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func exportMarkdown(result *AnalysisResult, path string) error {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("# Architecture Security Analysis: %s\n\n", result.ArchitectureName))
	b.WriteString(fmt.Sprintf("**Analysis Date:** %s\n", result.AnalysisDate.Format(time.RFC1123)))
	b.WriteString(fmt.Sprintf("**Analysis Mode:** %s\n", result.AnalysisMode))
	b.WriteString(fmt.Sprintf("**Total Assumptions:** %d\n", result.TotalAssumptions))
	if result.ConfidenceSummary != "" {
		b.WriteString(fmt.Sprintf("**Confidence:** %s\n", result.ConfidenceSummary))
	}
	b.WriteString("\n")

	b.WriteString("## Summary\n\n")
	b.WriteString(result.Summary + "\n\n")

	if result.EvidenceSummary.TotalSources > 0 {
		b.WriteString("## Evidence Summary\n\n")
		b.WriteString(fmt.Sprintf("- **Sources:** %d\n", result.EvidenceSummary.TotalSources))
		b.WriteString(fmt.Sprintf("- **Components Matched:** %d\n", result.EvidenceSummary.TotalComponents))
		b.WriteString(fmt.Sprintf("- **Relationships Matched:** %d\n", result.EvidenceSummary.TotalRelationships))
		for _, sf := range result.EvidenceSummary.SourceFiles {
			b.WriteString(fmt.Sprintf("- Source: %s\n", sf))
		}
		b.WriteString("\n")
	}

	b.WriteString("## Risk Distribution\n\n")
	b.WriteString(fmt.Sprintf("- **Critical:** %d\n", result.CriticalCount))
	b.WriteString(fmt.Sprintf("- **High:** %d\n", result.HighCount))
	b.WriteString(fmt.Sprintf("- **Medium:** %d\n", result.TotalAssumptions-result.CriticalCount-result.HighCount))
	b.WriteString(fmt.Sprintf("- **Risk Model:** %s\n", result.RiskModelVersion))
	b.WriteString("\n")

	if len(result.StrideDistribution) > 0 {
		b.WriteString("## STRIDE Distribution\n\n")
		total := 0
		for _, count := range result.StrideDistribution {
			total += count
		}
		for cat, count := range result.StrideDistribution {
			pct := float64(count) / float64(total) * 100
			b.WriteString(fmt.Sprintf("- **%s:** %d (%.1f%%)\n", cat, count, pct))
		}
		b.WriteString("\n")
	}

	b.WriteString("## Detailed Assumptions\n\n")
	for _, a := range result.Assumptions {
		b.WriteString(fmt.Sprintf("### %s — %s\n\n", a.ID, a.Description))

		b.WriteString(fmt.Sprintf("- **Risk:** %s\n", a.Risk))
		b.WriteString(fmt.Sprintf("- **STRIDE:** %s\n", strings.Join(convertStride(a.Stride), ", ")))
		b.WriteString(fmt.Sprintf("- **Confidence:** %.0f%%\n", a.Confidence*100))

		if a.Rationale != "" {
			b.WriteString(fmt.Sprintf("- **Rationale:** %s\n", a.Rationale))
		}

		if a.RiskJustification != nil {
			b.WriteString(fmt.Sprintf("- **Likelihood:** %d/5 — %s\n", a.RiskJustification.Likelihood, a.RiskJustification.LikelihoodReason))
			b.WriteString(fmt.Sprintf("- **Impact:** %d/5 — %s\n", a.RiskJustification.Impact, a.RiskJustification.ImpactReason))
			b.WriteString(fmt.Sprintf("- **Risk Score:** %d/25 → **%s**\n", a.RiskJustification.RiskScore, a.RiskJustification.RiskLevel))
		}

		if len(a.EvidenceSources) > 0 {
			b.WriteString("- **Evidence:**\n")
			for _, ev := range a.EvidenceSources {
				b.WriteString(fmt.Sprintf("  - %s\n", ev))
			}
		}

		if len(a.StrideJustifications) > 0 {
			b.WriteString("- **STRIDE Justification:**\n")
			for _, sj := range a.StrideJustifications {
				b.WriteString(fmt.Sprintf("  - **%s:** %s (confidence: %.0f%%)\n", sj.Category, sj.Reason, sj.Confidence*100))
			}
		}

		if a.ReviewStatus != "" {
			b.WriteString(fmt.Sprintf("- **Review:** %s\n", a.ReviewStatus))
			if a.ReviewNotes != "" {
				b.WriteString(fmt.Sprintf("- **Review Notes:** %s\n", a.ReviewNotes))
			}
		}

		b.WriteString("\n")
	}

	if len(result.Controls) > 0 {
		b.WriteString("## Recommended Controls\n\n")
		for _, c := range result.Controls {
			b.WriteString(fmt.Sprintf("- **%s:** %s\n", c.ID, c.Description))
			b.WriteString(fmt.Sprintf("  - Rationale: %s\n", c.Rationale))
			if len(c.MitigatedAssumptionIDs) > 0 {
				ids := c.MitigatedAssumptionIDs
				if len(ids) > 5 {
					ids = ids[:5]
				}
				b.WriteString(fmt.Sprintf("  - Mitigates: %s\n", strings.Join(ids, ", ")))
				if len(c.MitigatedAssumptionIDs) > 5 {
					b.WriteString(fmt.Sprintf("  - (and %d more)\n", len(c.MitigatedAssumptionIDs)-5))
				}
			}
			if len(c.MitigatedSTRIDE) > 0 {
				strs := make([]string, len(c.MitigatedSTRIDE))
				for i, s := range c.MitigatedSTRIDE {
					strs[i] = string(s)
				}
				b.WriteString(fmt.Sprintf("  - Addresses STRIDE: %s\n", strings.Join(strs, ", ")))
			}
			b.WriteString("\n")
		}
	}

	if len(result.Compliance) > 0 {
		b.WriteString("## Compliance Findings\n\n")
		for _, c := range result.Compliance {
			b.WriteString(fmt.Sprintf("- %s\n", c))
		}
		b.WriteString("\n")
	}

	b.WriteString("## Explainability Report\n\n")
	b.WriteString("For each assumption, this section traces every output back to its evidence, rules, and reasoning.\n\n")
	for _, a := range result.Assumptions {
		b.WriteString(fmt.Sprintf("### %s\n\n", a.ID))
		b.WriteString(fmt.Sprintf("**Assumption:** %s\n\n", a.Description))
		b.WriteString("| Property | Value |\n")
		b.WriteString("|----------|-------|\n")

		if a.SourceNode != "" {
			b.WriteString(fmt.Sprintf("| Source Node | %s |\n", a.SourceNode))
		}
		if a.SourceLine > 0 {
			b.WriteString(fmt.Sprintf("| Source Line | %d |\n", a.SourceLine))
		}

		b.WriteString(fmt.Sprintf("| Evidence | %s |\n", strings.Join(a.EvidenceSources, "; ")))
		b.WriteString(fmt.Sprintf("| STRIDE Mapping | %s |\n", strings.Join(convertStride(a.Stride), ", ")))

		if len(a.StrideJustifications) > 0 {
			for _, sj := range a.StrideJustifications {
				b.WriteString(fmt.Sprintf("| STRIDE Reason | %s — %s (confidence: %.0f%%) |\n", sj.Category, sj.Reason, sj.Confidence*100))
			}
		}

		b.WriteString(fmt.Sprintf("| Risk Level | %s |\n", a.Risk))
		if a.RiskJustification != nil {
			b.WriteString(fmt.Sprintf("| Risk Score | %d/25 |\n", a.RiskJustification.RiskScore))
			b.WriteString(fmt.Sprintf("| Risk Reason | %s |\n", a.RiskJustification.RiskReason))
			b.WriteString(fmt.Sprintf("| Likelihood | %d/5 — %s |\n", a.RiskJustification.Likelihood, a.RiskJustification.LikelihoodReason))
			b.WriteString(fmt.Sprintf("| Impact | %d/5 — %s |\n", a.RiskJustification.Impact, a.RiskJustification.ImpactReason))
		}

		b.WriteString(fmt.Sprintf("| Confidence | %.0f%% |\n", a.Confidence*100))

		if len(a.StrideJustifications) > 0 {
			var confReasons []string
			for _, sj := range a.StrideJustifications {
				if sj.ConfidenceReason != "" {
					confReasons = append(confReasons, fmt.Sprintf("%s: %s", sj.Category, sj.ConfidenceReason))
				}
			}
			if len(confReasons) > 0 {
				b.WriteString(fmt.Sprintf("| Confidence Factors | %s |\n", strings.Join(confReasons, "; ")))
			}
		}

		if a.RiskJustification != nil && a.RiskJustification.ConfidenceReason != "" {
			b.WriteString(fmt.Sprintf("| Risk Confidence | %.0f%% — %s |\n", a.RiskJustification.Confidence*100, a.RiskJustification.ConfidenceReason))
		}

		// Recommended controls for this assumption
		var ctrlIDs []string
		for _, c := range result.Controls {
			for _, id := range c.MitigatedAssumptionIDs {
				if id == a.ID {
					ctrlIDs = append(ctrlIDs, fmt.Sprintf("%s (%s)", c.ID, c.Description))
					break
				}
			}
		}
		if len(ctrlIDs) > 0 {
			b.WriteString(fmt.Sprintf("| Recommended Controls | %s |\n", strings.Join(ctrlIDs, "; ")))
		}

		b.WriteString("\n")
	}

	b.WriteString("---\n")
	b.WriteString("*Generated by ASF (Architecture Security Framework) — Deterministic Security Review Engine*\n")

	return os.WriteFile(path, []byte(b.String()), 0644)
}

func exportHTML(result *AnalysisResult, path string) error {
	var b strings.Builder
	b.WriteString("<!DOCTYPE html>\n<html>\n<head>\n")
	b.WriteString("<meta charset=\"UTF-8\">\n")
	b.WriteString(fmt.Sprintf("<title>ASF Analysis: %s</title>\n", result.ArchitectureName))
	b.WriteString("<style>\n")
	b.WriteString("body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Helvetica,Arial,sans-serif;max-width:960px;margin:40px auto;padding:0 20px;color:#1a1a2e;background:#f8f9fa}\n")
	b.WriteString("h1{color:#1a1a2e;border-bottom:3px solid #e94560;padding-bottom:10px}\n")
	b.WriteString("h2{color:#16213e;margin-top:30px}\n")
	b.WriteString("h3{color:#1a1a2e;margin-top:24px;padding:8px;background:#e2e8f0;border-radius:4px}\n")
	b.WriteString(".critical{color:#e94560;font-weight:700}\n")
	b.WriteString(".high{color:#f59e0b;font-weight:700}\n")
	b.WriteString(".medium{color:#3b82f6}\n")
	b.WriteString(".low{color:#10b981}\n")
	b.WriteString(".badge{display:inline-block;padding:2px 8px;border-radius:4px;font-size:12px;font-weight:700;margin-right:6px}\n")
	b.WriteString(".badge-critical{background:#e94560;color:#fff}\n")
	b.WriteString(".badge-high{background:#f59e0b;color:#fff}\n")
	b.WriteString(".badge-medium{background:#3b82f6;color:#fff}\n")
	b.WriteString(".badge-low{background:#10b981;color:#fff}\n")
	b.WriteString(".stride-bar{display:flex;height:24px;border-radius:4px;overflow:hidden;margin:10px 0}\n")
	b.WriteString(".stride-segment{padding:2px 8px;font-size:11px;color:#fff;font-weight:700;white-space:nowrap}\n")
	b.WriteString("table{width:100%;border-collapse:collapse;margin:10px 0}\n")
	b.WriteString("th,td{text-align:left;padding:8px 12px;border-bottom:1px solid #e2e8f0;vertical-align:top}\n")
	b.WriteString("th{background:#1a1a2e;color:#fff;font-size:13px}\n")
	b.WriteString("tr:hover{background:#e2e8f0}\n")
	b.WriteString(".evidence{font-size:12px;color:#64748b;margin:4px 0}\n")
	b.WriteString(".reasoning{font-size:13px;color:#334155;margin:6px 0;padding:6px;background:#f1f5f9;border-radius:4px}\n")
	b.WriteString(".confidence-bar{display:inline-block;height:8px;border-radius:4px;margin:2px 4px}\n")
	b.WriteString(".detail-section{margin:8px 0 8px 20px;padding:6px 10px;border-left:3px solid #cbd5e1;font-size:13px}\n")
	b.WriteString(".footer{margin-top:40px;padding-top:20px;border-top:1px solid #e2e8f0;font-size:13px;color:#64748b}\n")
	b.WriteString("</style>\n</head>\n<body>\n")

	b.WriteString(fmt.Sprintf("<h1>ASF Analysis Report</h1>\n"))
	b.WriteString(fmt.Sprintf("<p><strong>Architecture:</strong> %s<br>\n", result.ArchitectureName))
	b.WriteString(fmt.Sprintf("<strong>Mode:</strong> %s<br>\n", result.AnalysisMode))
	b.WriteString(fmt.Sprintf("<strong>Date:</strong> %s<br>\n", result.AnalysisDate.Format("Jan 2, 2006 15:04")))
	b.WriteString(fmt.Sprintf("<strong>Total Assumptions:</strong> %d", result.TotalAssumptions))
	if result.ConfidenceSummary != "" {
		b.WriteString(fmt.Sprintf("<br><strong>Confidence:</strong> %s", result.ConfidenceSummary))
	}
	b.WriteString("</p>\n")

	b.WriteString("<h2>Summary</h2>\n")
	b.WriteString(fmt.Sprintf("<p>%s</p>\n", result.Summary))

	if result.EvidenceSummary.TotalSources > 0 {
		b.WriteString("<h2>Evidence Sources</h2>\n<ul>\n")
		for _, sf := range result.EvidenceSummary.SourceFiles {
			b.WriteString(fmt.Sprintf("<li>%s</li>\n", sf))
		}
		b.WriteString(fmt.Sprintf("<li>Total components matched: %d</li>\n", result.EvidenceSummary.TotalComponents))
		b.WriteString(fmt.Sprintf("<li>Total relationships matched: %d</li>\n", result.EvidenceSummary.TotalRelationships))
		b.WriteString("</ul>\n")
	}

	b.WriteString("<h2>Risk Distribution</h2>\n")
	b.WriteString("<table>\n<tr><th>Level</th><th>Count</th></tr>\n")
	highRisk := 0
	for _, a := range result.Assumptions {
		if a.Risk == RiskHigh || a.Risk == RiskCritical {
			highRisk++
		}
	}
	medRisk := result.TotalAssumptions - highRisk
	b.WriteString(fmt.Sprintf("<tr><td class=\"badge badge-critical\">Critical</td><td>%d</td></tr>\n", result.CriticalCount))
	b.WriteString(fmt.Sprintf("<tr><td class=\"badge badge-high\">High</td><td>%d</td></tr>\n", result.HighCount))
	b.WriteString(fmt.Sprintf("<tr><td class=\"badge badge-medium\">Medium</td><td>%d</td></tr>\n", medRisk))
	b.WriteString("</table>\n")
	b.WriteString(fmt.Sprintf("<p><em>Risk Model: %s</em></p>\n", result.RiskModelVersion))

	if len(result.StrideDistribution) > 0 {
		b.WriteString("<h2>STRIDE Distribution</h2>\n")
		b.WriteString("<div class=\"stride-bar\">\n")
		colors := []string{"#e94560", "#f59e0b", "#3b82f6", "#10b981", "#8b5cf6", "#ec4899"}
		labels := []StrideCategory{StrideSpoofing, StrideTampering, StrideRepudiation, StrideInfoDisclosure, StrideDenialOfService, StrideElevationPriv}
		total := 0
		for _, count := range result.StrideDistribution {
			total += count
		}
		if total > 0 {
			for i, cat := range labels {
				if count := result.StrideDistribution[cat]; count > 0 {
					pct := float64(count) / float64(total) * 100
					color := colors[i%len(colors)]
					b.WriteString(fmt.Sprintf("<div class=\"stride-segment\" style=\"width:%.1f%%;background:%s\">%s %d</div>\n", pct, color, cat, count))
				}
			}
		}
		b.WriteString("</div>\n")
	}

	b.WriteString("<h2>Detailed Assumptions</h2>\n")
	for _, a := range result.Assumptions {
		badgeClass := "badge-medium"
		if a.Risk == RiskCritical {
			badgeClass = "badge-critical"
		} else if a.Risk == RiskHigh {
			badgeClass = "badge-high"
		} else if a.Risk == RiskLow {
			badgeClass = "badge-low"
		}
		strideStrs := make([]string, len(a.Stride))
		for i, s := range a.Stride {
			strideStrs[i] = string(s)
		}
		aiTag := ""
		if strings.HasPrefix(a.ID, "AI-") {
			aiTag = " <span style=\"color:#10b981;font-size:11px\">[AI]</span>"
		}
		reviewTag := ""
		if a.ReviewStatus != "" {
			reviewTag = fmt.Sprintf(" <span class=\"badge\">%s</span>", a.ReviewStatus)
		}

		b.WriteString(fmt.Sprintf("<h3>%s%s%s</h3>\n", a.ID, aiTag, reviewTag))
		b.WriteString(fmt.Sprintf("<p><strong>%s</strong></p>\n", a.Description))
		b.WriteString(fmt.Sprintf("<p><span class=\"badge %s\">%s</span>", badgeClass, a.Risk))
		b.WriteString(fmt.Sprintf(" <strong>STRIDE:</strong> %s", strings.Join(strideStrs, ", ")))
		b.WriteString(fmt.Sprintf(" <strong>Confidence:</strong> %.0f%%</p>\n", a.Confidence*100))

		if a.Rationale != "" {
			b.WriteString(fmt.Sprintf("<div class=\"reasoning\">%s</div>\n", a.Rationale))
		}

		if a.RiskJustification != nil {
			b.WriteString("<div class=\"detail-section\">\n")
			b.WriteString(fmt.Sprintf("<strong>Likelihood:</strong> %d/5 — %s<br>\n", a.RiskJustification.Likelihood, a.RiskJustification.LikelihoodReason))
			if len(a.RiskJustification.LikelihoodFactors) > 0 {
				b.WriteString("<ul>\n")
				for _, lf := range a.RiskJustification.LikelihoodFactors {
					b.WriteString(fmt.Sprintf("<li>%s: %d — %s</li>\n", lf.Factor, lf.Value, lf.Reason))
				}
				b.WriteString("</ul>\n")
			}
			b.WriteString(fmt.Sprintf("<strong>Impact:</strong> %d/5 — %s<br>\n", a.RiskJustification.Impact, a.RiskJustification.ImpactReason))
			if len(a.RiskJustification.ImpactFactors) > 0 {
				b.WriteString("<ul>\n")
				for _, ifa := range a.RiskJustification.ImpactFactors {
					b.WriteString(fmt.Sprintf("<li>%s: %d — %s</li>\n", ifa.Factor, ifa.Value, ifa.Reason))
				}
				b.WriteString("</ul>\n")
			}
			b.WriteString(fmt.Sprintf("<strong>Risk Score:</strong> %d/25 → <strong>%s</strong><br>\n", a.RiskJustification.RiskScore, a.RiskJustification.RiskLevel))
			b.WriteString("</div>\n")
		}

		if len(a.EvidenceSources) > 0 {
			b.WriteString("<div class=\"evidence\">\n<strong>Evidence:</strong><br>\n")
			for _, ev := range a.EvidenceSources {
				b.WriteString(fmt.Sprintf("&nbsp;&nbsp;%s<br>\n", ev))
			}
			b.WriteString("</div>\n")
		}

		if len(a.StrideJustifications) > 0 {
			b.WriteString("<div class=\"evidence\">\n<strong>STRIDE Justification:</strong><br>\n")
			for _, sj := range a.StrideJustifications {
				b.WriteString(fmt.Sprintf("&nbsp;&nbsp;<strong>%s:</strong> %s (conf: %.0f%%)<br>\n", sj.Category, sj.Reason, sj.Confidence*100))
			}
			b.WriteString("</div>\n")
		}
	}

	if len(result.Controls) > 0 {
		b.WriteString("<h2>Recommended Controls</h2>\n")
		for _, c := range result.Controls {
			b.WriteString(fmt.Sprintf("<div style=\"margin:8px 0;padding:10px;background:#f1f5f9;border-radius:4px;border-left:4px solid #3b82f6\">\n"))
			b.WriteString(fmt.Sprintf("<strong>%s:</strong> %s<br>\n", c.ID, c.Description))
			b.WriteString(fmt.Sprintf("<em>%s</em><br>\n", c.Rationale))
			if len(c.MitigatedAssumptionIDs) > 0 {
				ids := c.MitigatedAssumptionIDs
				if len(ids) > 5 {
					ids = ids[:5]
				}
				b.WriteString(fmt.Sprintf("<span style=\"font-size:12px;color:#64748b\">Assumptions: %s</span><br>\n", strings.Join(ids, ", ")))
			}
			if len(c.MitigatedSTRIDE) > 0 {
				strs := make([]string, len(c.MitigatedSTRIDE))
				for i, s := range c.MitigatedSTRIDE {
					strs[i] = string(s)
				}
				b.WriteString(fmt.Sprintf("<span style=\"font-size:12px;color:#64748b\">STRIDE: %s</span>\n", strings.Join(strs, ", ")))
			}
			b.WriteString("</div>\n")
		}
	}

	b.WriteString("<h2>Explainability Report</h2>\n")
	b.WriteString("<p>For each assumption, this section traces every output back to its evidence, rules, and reasoning.</p>\n")
	for _, a := range result.Assumptions {
		b.WriteString(fmt.Sprintf("<h3>%s</h3>\n", a.ID))
		b.WriteString(fmt.Sprintf("<p><strong>Assumption:</strong> %s</p>\n", a.Description))
		b.WriteString("<table>\n<tr><th>Property</th><th>Value</th></tr>\n")
		if a.SourceNode != "" {
			b.WriteString(fmt.Sprintf("<tr><td>Source Node</td><td>%s</td></tr>\n", a.SourceNode))
		}
		if a.SourceLine > 0 {
			b.WriteString(fmt.Sprintf("<tr><td>Source Line</td><td>%d</td></tr>\n", a.SourceLine))
		}
		b.WriteString(fmt.Sprintf("<tr><td>Evidence</td><td>%s</td></tr>\n", strings.Join(a.EvidenceSources, "; ")))
		b.WriteString(fmt.Sprintf("<tr><td>STRIDE Mapping</td><td>%s</td></tr>\n", strings.Join(convertStride(a.Stride), ", ")))
		if len(a.StrideJustifications) > 0 {
			for _, sj := range a.StrideJustifications {
				b.WriteString(fmt.Sprintf("<tr><td>STRIDE Reason</td><td>%s — %s (conf: %.0f%%)</td></tr>\n", sj.Category, sj.Reason, sj.Confidence*100))
			}
		}
		b.WriteString(fmt.Sprintf("<tr><td>Risk Level</td><td>%s</td></tr>\n", a.Risk))
		if a.RiskJustification != nil {
			b.WriteString(fmt.Sprintf("<tr><td>Risk Score</td><td>%d/25</td></tr>\n", a.RiskJustification.RiskScore))
			b.WriteString(fmt.Sprintf("<tr><td>Risk Reason</td><td>%s</td></tr>\n", a.RiskJustification.RiskReason))
			b.WriteString(fmt.Sprintf("<tr><td>Likelihood</td><td>%d/5 — %s</td></tr>\n", a.RiskJustification.Likelihood, a.RiskJustification.LikelihoodReason))
			b.WriteString(fmt.Sprintf("<tr><td>Impact</td><td>%d/5 — %s</td></tr>\n", a.RiskJustification.Impact, a.RiskJustification.ImpactReason))
		}
		b.WriteString(fmt.Sprintf("<tr><td>Confidence</td><td>%.0f%%</td></tr>\n", a.Confidence*100))
		// Recommended controls for this assumption
		var ctrlIDs []string
		for _, c := range result.Controls {
			for _, id := range c.MitigatedAssumptionIDs {
				if id == a.ID {
					ctrlIDs = append(ctrlIDs, fmt.Sprintf("%s (%s)", c.ID, c.Description))
					break
				}
			}
		}
		if len(ctrlIDs) > 0 {
			b.WriteString(fmt.Sprintf("<tr><td>Recommended Controls</td><td>%s</td></tr>\n", strings.Join(ctrlIDs, "; ")))
		}
		b.WriteString("</table>\n")
	}

	b.WriteString("<div class=\"footer\">\n")
	b.WriteString("Generated by ASF (Architecture Security Framework) — Deterministic Security Review Engine<br>\n")
	b.WriteString("Risk Model: asf-risk-model-1.0 | 5x5 matrix | Likelihood × Impact = Risk Score<br>\n")
	b.WriteString("All outputs are deterministic and reproducible.\n")
	b.WriteString("</div>\n</body>\n</html>\n")
	return os.WriteFile(path, []byte(b.String()), 0644)
}

func exportCSV(result *AnalysisResult, path string) error {
	var b strings.Builder
	b.WriteString("ID,Description,Component,Category,Risk,STRIDE,Likelihood,Impact,RiskScore,Confidence,ReviewStatus,\"EvidenceSources\",SourceNode,SourceLine,\"Rationale\",\"StrideReason\",\"RiskReason\",\"MitigatingControls\"\n")
	for _, a := range result.Assumptions {
		strideStr := strings.Join(convertStride(a.Stride), ";")
		evidenceStr := strings.Join(a.EvidenceSources, " | ")
		riskScore := 0
		riskReason := ""
		if a.RiskJustification != nil {
			riskScore = a.RiskJustification.RiskScore
			riskReason = a.RiskJustification.RiskReason
		}
		strideReason := ""
		if len(a.StrideJustifications) > 0 {
			var parts []string
			for _, sj := range a.StrideJustifications {
				parts = append(parts, fmt.Sprintf("%s: %s", sj.Category, sj.Reason))
			}
			strideReason = strings.Join(parts, " | ")
		}
		// Find controls that mitigate this assumption
		var ctrlIDs []string
		for _, c := range result.Controls {
			for _, id := range c.MitigatedAssumptionIDs {
				if id == a.ID {
					ctrlIDs = append(ctrlIDs, c.ID)
					break
				}
			}
		}
		ctrlStr := strings.Join(ctrlIDs, ";")
		b.WriteString(fmt.Sprintf("%s,\"%s\",%s,%s,%s,%s,%d,%d,%d,%.2f,%s,\"%s\",%s,%d,\"%s\",\"%s\",\"%s\",\"%s\"\n",
			a.ID, a.Description, a.Component, a.Category, a.Risk,
			strideStr, a.Likelihood, a.Impact, riskScore, a.Confidence,
			a.ReviewStatus, evidenceStr, a.SourceNode, a.SourceLine,
			a.Rationale, strideReason, riskReason, ctrlStr))
	}
	return os.WriteFile(path, []byte(b.String()), 0644)
}

func exportPDF(result *AnalysisResult, path string) error {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetCreator("ASF-TUI", true)
	pdf.SetTitle("Architecture Security Framework - Analysis Report", true)

	// Title page
	pdf.AddPage()
	pdf.SetFont("Helvetica", "B", 24)
	pdf.CellFormat(190, 20, "ASF Analysis Report", "", 1, "C", false, 0, "")
	pdf.Ln(10)
	pdf.SetFont("Helvetica", "", 11)
	pdf.CellFormat(190, 8, "Architecture: "+result.ArchitectureName, "", 1, "C", false, 0, "")
	pdf.CellFormat(190, 8, "Analysis Mode: "+result.AnalysisMode, "", 1, "C", false, 0, "")
	pdf.CellFormat(190, 8, "Date: "+time.Now().Format("2006-01-02 15:04"), "", 1, "C", false, 0, "")
	pdf.Ln(5)
	pdf.CellFormat(190, 8, fmt.Sprintf("Total Assumptions: %d", result.TotalAssumptions), "", 1, "C", false, 0, "")
	if result.ConfidenceSummary != "" {
		pdf.CellFormat(190, 8, "Confidence: "+result.ConfidenceSummary, "", 1, "C", false, 0, "")
	}

	// Summary page
	pdf.AddPage()
	pdf.SetFont("Helvetica", "B", 16)
	pdf.CellFormat(190, 12, "Summary", "", 1, "L", false, 0, "")
	pdf.Ln(4)
	pdf.SetFont("Helvetica", "", 11)
	pdf.CellFormat(190, 7, fmt.Sprintf("Total Assumptions: %d", result.TotalAssumptions), "", 1, "L", false, 0, "")

	if result.EvidenceSummary.TotalSources > 0 {
		pdf.CellFormat(190, 7, fmt.Sprintf("Evidence Sources: %d", result.EvidenceSummary.TotalSources), "", 1, "L", false, 0, "")
		pdf.CellFormat(190, 7, fmt.Sprintf("Components Matched: %d", result.EvidenceSummary.TotalComponents), "", 1, "L", false, 0, "")
	}
	pdf.Ln(2)

	highRisk, medRisk, lowRisk := 0, 0, 0
	for _, a := range result.Assumptions {
		switch a.Risk {
		case RiskHigh:
			highRisk++
		case RiskMedium:
			medRisk++
		case RiskLow:
			lowRisk++
		}
	}
	pdf.CellFormat(190, 7, fmt.Sprintf("Critical: %d", result.CriticalCount), "", 1, "L", false, 0, "")
	pdf.CellFormat(190, 7, fmt.Sprintf("High Risk: %d", highRisk), "", 1, "L", false, 0, "")
	pdf.CellFormat(190, 7, fmt.Sprintf("Medium Risk: %d", medRisk), "", 1, "L", false, 0, "")
	pdf.CellFormat(190, 7, fmt.Sprintf("Low Risk: %d", lowRisk), "", 1, "L", false, 0, "")
	pdf.CellFormat(190, 7, "Risk Model: "+result.RiskModelVersion, "", 1, "L", false, 0, "")
	pdf.Ln(4)

	pdf.SetFont("Helvetica", "B", 12)
	pdf.CellFormat(190, 8, "STRIDE Distribution", "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 11)
	strideLabels := []StrideCategory{StrideSpoofing, StrideTampering, StrideRepudiation, StrideInfoDisclosure, StrideDenialOfService, StrideElevationPriv}
	for _, s := range strideLabels {
		if count, ok := result.StrideDistribution[s]; ok {
			pdf.CellFormat(190, 7, fmt.Sprintf("  %s: %d", string(s), count), "", 1, "L", false, 0, "")
		}
	}

	// Controls section
	if len(result.Controls) > 0 {
		pdf.AddPage()
		pdf.SetFont("Helvetica", "B", 16)
		pdf.CellFormat(190, 12, "Recommended Controls", "", 1, "L", false, 0, "")
		pdf.Ln(4)
		for _, c := range result.Controls {
			pdf.SetFont("Helvetica", "B", 10)
			pdf.CellFormat(190, 6, fmt.Sprintf("%s: %s", c.ID, c.Description), "", 1, "L", false, 0, "")
			pdf.SetFont("Helvetica", "I", 8)
			pdf.MultiCell(190, 4, "Rationale: "+c.Rationale, "", "L", false)
			if len(c.MitigatedAssumptionIDs) > 0 {
				ids := c.MitigatedAssumptionIDs
				if len(ids) > 5 {
					ids = ids[:5]
				}
				pdf.SetFont("Helvetica", "", 7)
				pdf.CellFormat(190, 3, "Mitigates: "+strings.Join(ids, ", "), "", 1, "L", false, 0, "")
			}
			if len(c.MitigatedSTRIDE) > 0 {
				strs := make([]string, len(c.MitigatedSTRIDE))
				for i, s := range c.MitigatedSTRIDE {
					strs[i] = string(s)
				}
				pdf.SetFont("Helvetica", "", 7)
				pdf.CellFormat(190, 3, "STRIDE: "+strings.Join(strs, ", "), "", 1, "L", false, 0, "")
			}
			pdf.Ln(4)
		}
	}

	// Detailed Assumptions
	for i, a := range result.Assumptions {
		// New page every 3 assumptions
		if i%3 == 0 {
			pdf.AddPage()
			pdf.SetFont("Helvetica", "B", 16)
			pdf.CellFormat(190, 12, fmt.Sprintf("Assumption %s", a.ID), "", 1, "L", false, 0, "")
		} else {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 8, fmt.Sprintf("Assumption %s", a.ID), "", 1, "L", false, 0, "")
		}

		strideStrs := make([]string, len(a.Stride))
		for i, s := range a.Stride {
			strideStrs[i] = string(s)
		}
		pdf.SetFont("Helvetica", "", 10)
		pdf.MultiCell(190, 5, a.Description, "", "L", false)
		pdf.Ln(1)

		pdf.SetFont("Helvetica", "B", 9)
		pdf.CellFormat(190, 5, fmt.Sprintf("Risk: %s | STRIDE: %s | Confidence: %.0f%%", string(a.Risk), strings.Join(strideStrs, ", "), a.Confidence*100), "", 1, "L", false, 0, "")

		if a.RiskJustification != nil {
			pdf.SetFont("Helvetica", "", 8)
			pdf.CellFormat(190, 4, fmt.Sprintf("Likelihood %d/5: %s", a.RiskJustification.Likelihood, a.RiskJustification.LikelihoodReason), "", 1, "L", false, 0, "")
			pdf.CellFormat(190, 4, fmt.Sprintf("Impact %d/5: %s", a.RiskJustification.Impact, a.RiskJustification.ImpactReason), "", 1, "L", false, 0, "")
			pdf.CellFormat(190, 4, fmt.Sprintf("Risk Score: %d/25 -> %s", a.RiskJustification.RiskScore, a.RiskJustification.RiskLevel), "", 1, "L", false, 0, "")
		}

		if a.Rationale != "" {
			pdf.SetFont("Helvetica", "I", 8)
			pdf.MultiCell(190, 4, "Rationale: "+a.Rationale, "", "L", false)
		}

		if len(a.EvidenceSources) > 0 {
			pdf.SetFont("Helvetica", "", 7)
			for _, ev := range a.EvidenceSources {
				pdf.CellFormat(190, 3, "  "+ev, "", 1, "L", false, 0, "")
			}
		}

		if len(a.StrideJustifications) > 0 {
			pdf.SetFont("Helvetica", "B", 8)
			pdf.CellFormat(190, 4, "STRIDE Justification:", "", 1, "L", false, 0, "")
			pdf.SetFont("Helvetica", "", 7)
			for _, sj := range a.StrideJustifications {
				pdf.CellFormat(190, 3, fmt.Sprintf("  %s: %s (conf: %.0f%%)", sj.Category, sj.Reason, sj.Confidence*100), "", 1, "L", false, 0, "")
			}
		}

		if a.SourceNode != "" || a.SourceLine > 0 {
			pdf.SetFont("Helvetica", "", 7)
			if a.SourceNode != "" {
				pdf.CellFormat(190, 3, "  Source Node: "+a.SourceNode, "", 1, "L", false, 0, "")
			}
			if a.SourceLine > 0 {
				pdf.CellFormat(190, 3, fmt.Sprintf("  Source Line: %d", a.SourceLine), "", 1, "L", false, 0, "")
			}
		}

		pdf.Ln(4)
	}

	// Explainability Report
	pdf.AddPage()
	pdf.SetFont("Helvetica", "B", 16)
	pdf.CellFormat(190, 12, "Explainability Report", "", 1, "L", false, 0, "")
	pdf.SetFont("Helvetica", "", 8)
	pdf.MultiCell(190, 4, "For each assumption, this section traces every output back to its evidence, rules, and reasoning.", "", "L", false)
	pdf.Ln(4)
	for _, a := range result.Assumptions {
		pdf.SetFont("Helvetica", "B", 10)
		pdf.CellFormat(190, 6, fmt.Sprintf("Assumption %s", a.ID), "", 1, "L", false, 0, "")
		pdf.SetFont("Helvetica", "", 8)
		pdf.MultiCell(190, 4, "Assumption: "+a.Description, "", "L", false)
		if a.SourceNode != "" {
			pdf.CellFormat(190, 4, "  Source Node: "+a.SourceNode, "", 1, "L", false, 0, "")
		}
		if a.SourceLine > 0 {
			pdf.CellFormat(190, 4, fmt.Sprintf("  Source Line: %d", a.SourceLine), "", 1, "L", false, 0, "")
		}
		pdf.CellFormat(190, 4, "  STRIDE: "+strings.Join(convertStride(a.Stride), ", "), "", 1, "L", false, 0, "")
		if len(a.StrideJustifications) > 0 {
			for _, sj := range a.StrideJustifications {
				pdf.CellFormat(190, 4, fmt.Sprintf("  STRIDE Reason: %s - %s (conf: %.0f%%)", sj.Category, sj.Reason, sj.Confidence*100), "", 1, "L", false, 0, "")
			}
		}
		pdf.CellFormat(190, 4, fmt.Sprintf("  Risk: %s", a.Risk), "", 1, "L", false, 0, "")
		if a.RiskJustification != nil {
			pdf.CellFormat(190, 4, fmt.Sprintf("  Risk Score: %d/25  Reason: %s", a.RiskJustification.RiskScore, a.RiskJustification.RiskReason), "", 1, "L", false, 0, "")
		}
		pdf.CellFormat(190, 4, fmt.Sprintf("  Confidence: %.0f%%", a.Confidence*100), "", 1, "L", false, 0, "")
		var ctrlIDs []string
		for _, c := range result.Controls {
			for _, id := range c.MitigatedAssumptionIDs {
				if id == a.ID {
					ctrlIDs = append(ctrlIDs, c.ID)
					break
				}
			}
		}
		if len(ctrlIDs) > 0 {
			pdf.CellFormat(190, 4, "  Controls: "+strings.Join(ctrlIDs, ", "), "", 1, "L", false, 0, "")
		}
		pdf.Ln(3)
	}

	return pdf.OutputFileAndClose(path)
}

func convertStride(s []StrideCategory) []string {
	r := make([]string, len(s))
	for i, v := range s {
		r[i] = string(v)
	}
	return r
}

type exportModel struct {
	selected         int
	format           ExportFormat
	done             bool
	exportPath       string
	showConfirmation bool
}

func newExportModel() exportModel {
	return exportModel{}
}

func (m exportModel) Update(msg tea.Msg) (exportModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < 4 {
				m.selected++
			}
		case "enter":
			formats := []ExportFormat{ExportJSON, ExportMarkdown, ExportCSV, ExportPDF, ExportHTML}
			if m.selected < len(formats) {
				m.format = formats[m.selected]
				m.showConfirmation = true
			}
		case "esc":
			m.showConfirmation = false
			m.done = false
		case "y":
			if m.showConfirmation && !m.done {
				m.done = true
			}
		}
	}
	return m, nil
}

func (m mainModel) viewExport() string {
	s := m.styles
	ex := m.exportV

	if ex.done {
		return lipgloss.JoinVertical(lipgloss.Left,
			s.Title.Render("Export"),
			s.BorderBox.Render(
				lipgloss.JoinVertical(lipgloss.Center,
					s.StatusGood.Render("✓ Export Complete"),
					s.SectionItem.Render(ex.exportPath),
					s.SectionItem.Render("Press Esc to return."),
				),
			),
		)
	}

	if ex.showConfirmation {
		return lipgloss.JoinVertical(lipgloss.Left,
			s.Title.Render("Export"),
			s.BorderBox.Render(
				lipgloss.JoinVertical(lipgloss.Center,
					s.SectionItem.Render(fmt.Sprintf("Export as %s?", ex.format)),
					s.SectionItem.Render("Press Y to confirm, Esc to cancel."),
				),
			),
		)
	}

	formats := []struct {
		name   string
		format ExportFormat
	}{
		{"JSON (.json)", ExportJSON},
		{"Markdown (.md)", ExportMarkdown},
		{"HTML (.html)", ExportHTML},
		{"CSV (.csv)", ExportCSV},
		{"PDF (.pdf)", ExportPDF},
	}

	var items []string
	for i, f := range formats {
		style := s.SectionItem
		if i == ex.selected {
			style = s.MenuSelected
		}
		items = append(items, style.Render(f.name))
	}

	return lipgloss.JoinVertical(lipgloss.Left,
		s.Title.Render("Export Results"),
		s.Subtitle.Render("Select export format:"),
		lipgloss.JoinVertical(lipgloss.Left, items...),
	)
}

func (m mainModel) updateExport(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.exportV, cmd = m.exportV.Update(msg)
	return m, cmd
}
