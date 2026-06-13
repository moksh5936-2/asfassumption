package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"asf-tui/asf/confidencex"
	"asf-tui/asf/coverage"
	"asf-tui/asf/narrative"
	"asf-tui/asf/review"
	"asf-tui/asf/trust"
	"asf-tui/asf/verify"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/go-pdf/fpdf"
)

type ExportFormat string

const (
	ExportJSON               ExportFormat = "json"
	ExportMarkdown           ExportFormat = "markdown"
	ExportCSV                ExportFormat = "csv"
	ExportPDF                ExportFormat = "pdf"
	ExportHTML               ExportFormat = "html"
	ExportNarrativeMarkdown  ExportFormat = "narrative-md"
	ExportNarrativeHTML      ExportFormat = "narrative-html"
	ExportTrustMarkdown      ExportFormat = "trust-md"
	ExportTrustHTML          ExportFormat = "trust-html"
	ExportTrustJSON          ExportFormat = "trust-json"
	ExportCoverageMarkdown   ExportFormat = "coverage-md"
	ExportCoverageHTML       ExportFormat = "coverage-html"
	ExportCoverageJSON       ExportFormat = "coverage-json"
	ExportVerifyMarkdown     ExportFormat = "verify-md"
	ExportVerifyHTML         ExportFormat = "verify-html"
	ExportVerifyJSON         ExportFormat = "verify-json"
	ExportReviewMarkdown     ExportFormat = "review-md"
	ExportReviewHTML         ExportFormat = "review-html"
	ExportConfidenceMarkdown ExportFormat = "confidence-md"
	ExportConfidenceHTML     ExportFormat = "confidence-html"
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
	case ExportNarrativeMarkdown:
		path = filepath.Join(outputDir, baseName+"_narrative.md")
		err = exportNarrativeMarkdown(result, path)
	case ExportNarrativeHTML:
		path = filepath.Join(outputDir, baseName+"_narrative.html")
		err = exportNarrativeHTML(result, path)
	case ExportTrustMarkdown:
		path = filepath.Join(outputDir, baseName+"_trust.md")
		err = exportTrustMarkdown(result, path)
	case ExportTrustHTML:
		path = filepath.Join(outputDir, baseName+"_trust.html")
		err = exportTrustHTML(result, path)
	case ExportTrustJSON:
		path = filepath.Join(outputDir, baseName+"_trust.json")
		err = exportTrustJSON(result, path)
	case ExportCoverageMarkdown:
		path = filepath.Join(outputDir, baseName+"_coverage.md")
		err = exportCoverageMarkdown(result, path)
	case ExportCoverageHTML:
		path = filepath.Join(outputDir, baseName+"_coverage.html")
		err = exportCoverageHTML(result, path)
	case ExportCoverageJSON:
		path = filepath.Join(outputDir, baseName+"_coverage.json")
		err = exportCoverageJSON(result, path)
	case ExportVerifyMarkdown:
		path = filepath.Join(outputDir, baseName+"_verify.md")
		err = exportVerifyMarkdown(result, path)
	case ExportVerifyHTML:
		path = filepath.Join(outputDir, baseName+"_verify.html")
		err = exportVerifyHTML(result, path)
	case ExportVerifyJSON:
		path = filepath.Join(outputDir, baseName+"_verify.json")
		err = exportVerifyJSON(result, path)
	case ExportReviewMarkdown:
		path = filepath.Join(outputDir, baseName+"_review.md")
		err = exportReviewMarkdown(result, path)
	case ExportReviewHTML:
		path = filepath.Join(outputDir, baseName+"_review.html")
		err = exportReviewHTML(result, path)
	case ExportConfidenceMarkdown:
		path = filepath.Join(outputDir, baseName+"_confidence.md")
		err = exportConfidenceMarkdown(result, path)
	case ExportConfidenceHTML:
		path = filepath.Join(outputDir, baseName+"_confidence.html")
		err = exportConfidenceHTML(result, path)
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

func exportNarrativeMarkdown(result *AnalysisResult, path string) error {
	if result.NarrativeOutput == nil {
		return fmt.Errorf("narrative output not generated")
	}
	data := narrative.ExportMarkdown(result.NarrativeOutput)
	return os.WriteFile(path, []byte(data), 0644)
}

func exportNarrativeHTML(result *AnalysisResult, path string) error {
	if result.NarrativeOutput == nil {
		return fmt.Errorf("narrative output not generated")
	}
	data := narrative.ExportHTML(result.NarrativeOutput)
	return os.WriteFile(path, []byte(data), 0644)
}

func exportTrustMarkdown(result *AnalysisResult, path string) error {
	if result.TrustOutput == nil {
		return fmt.Errorf("trust chain output not generated")
	}
	data := trust.ExportMarkdown(result.TrustOutput)
	return os.WriteFile(path, []byte(data), 0644)
}

func exportTrustHTML(result *AnalysisResult, path string) error {
	if result.TrustOutput == nil {
		return fmt.Errorf("trust chain output not generated")
	}
	data := trust.ExportHTML(result.TrustOutput)
	return os.WriteFile(path, []byte(data), 0644)
}

func exportTrustJSON(result *AnalysisResult, path string) error {
	if result.TrustOutput == nil {
		return fmt.Errorf("trust chain output not generated")
	}
	data, err := json.MarshalIndent(result.TrustOutput, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func exportCoverageMarkdown(result *AnalysisResult, path string) error {
	if result.CoverageOutput == nil {
		return fmt.Errorf("coverage output not generated")
	}
	data := coverage.ExportMarkdown(result.CoverageOutput)
	return os.WriteFile(path, []byte(data), 0644)
}

func exportCoverageHTML(result *AnalysisResult, path string) error {
	if result.CoverageOutput == nil {
		return fmt.Errorf("coverage output not generated")
	}
	data := coverage.ExportHTML(result.CoverageOutput)
	return os.WriteFile(path, []byte(data), 0644)
}

func exportCoverageJSON(result *AnalysisResult, path string) error {
	if result.CoverageOutput == nil {
		return fmt.Errorf("coverage output not generated")
	}
	data, err := json.MarshalIndent(result.CoverageOutput, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func exportVerifyMarkdown(result *AnalysisResult, path string) error {
	if result.VerificationOutput == nil {
		return fmt.Errorf("verification output not generated")
	}
	data := verify.ExportMarkdown(result.VerificationOutput)
	return os.WriteFile(path, []byte(data), 0644)
}

func exportVerifyHTML(result *AnalysisResult, path string) error {
	if result.VerificationOutput == nil {
		return fmt.Errorf("verification output not generated")
	}
	data := verify.ExportHTML(result.VerificationOutput)
	return os.WriteFile(path, []byte(data), 0644)
}

func exportVerifyJSON(result *AnalysisResult, path string) error {
	if result.VerificationOutput == nil {
		return fmt.Errorf("verification output not generated")
	}
	data, err := json.MarshalIndent(result.VerificationOutput, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func exportReviewMarkdown(result *AnalysisResult, path string) error {
	if result.ReviewOutput == nil {
		return fmt.Errorf("review output not generated")
	}
	data := review.ExportMarkdown(result.ReviewOutput)
	return os.WriteFile(path, []byte(data), 0644)
}

func exportReviewHTML(result *AnalysisResult, path string) error {
	if result.ReviewOutput == nil {
		return fmt.Errorf("review output not generated")
	}
	data := review.ExportHTML(result.ReviewOutput)
	return os.WriteFile(path, []byte(data), 0644)
}

func exportConfidenceMarkdown(result *AnalysisResult, path string) error {
	if result.ConfidenceOutput == nil {
		return fmt.Errorf("confidence output not generated")
	}
	data := confidencex.ExportMarkdown(result.ConfidenceOutput)
	return os.WriteFile(path, []byte(data), 0644)
}

func exportConfidenceHTML(result *AnalysisResult, path string) error {
	if result.ConfidenceOutput == nil {
		return fmt.Errorf("confidence output not generated")
	}
	data := confidencex.ExportHTML(result.ConfidenceOutput)
	return os.WriteFile(path, []byte(data), 0644)
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
	b.WriteString(fmt.Sprintf("- **Medium:** %d\n", result.MediumCount))
	b.WriteString(fmt.Sprintf("- **Low:** %d\n", result.LowCount))
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

	if len(result.AttackPaths) > 0 {
		b.WriteString("## Attack Path Discovery\n\n")
		aps := result.AttackPathSummary
		b.WriteString(fmt.Sprintf("**Total Attack Paths:** %d\n", aps.TotalAttackPaths))
		b.WriteString(fmt.Sprintf("- **Critical:** %d\n", aps.CriticalCount))
		b.WriteString(fmt.Sprintf("- **High:** %d\n", aps.HighCount))
		b.WriteString(fmt.Sprintf("- **Medium:** %d\n", aps.MediumCount))
		b.WriteString(fmt.Sprintf("- **Low:** %d\n", aps.LowCount))
		b.WriteString(fmt.Sprintf("**Threat Chains:** %d\n", aps.ThreatChainCount))
		if aps.SummaryText != "" {
			b.WriteString(fmt.Sprintf("**Summary:** %s\n", aps.SummaryText))
		}
		b.WriteString("\n")

		b.WriteString("### Top Attack Paths\n\n")
		b.WriteString("| Path | Entry | Target | Risk | Detection | Business Impact |\n")
		b.WriteString("|------|-------|--------|------|-----------|----------------|\n")
		for _, p := range result.AttackPaths {
			riskLabel := riskLevelForAPDScore(p.RiskScore)
			b.WriteString(fmt.Sprintf("| %s | %s | %s | %.2f (%s) | %s | %s |\n",
				p.Name, p.EntryPoint, p.TargetAsset, p.RiskScore, riskLabel, p.DetectionDifficulty, p.BusinessImpact))
		}
		b.WriteString("\n")

		if len(aps.KillChainCoverage) > 0 {
			b.WriteString("### Kill Chain Coverage\n\n")
			b.WriteString("| Phase | Paths |\n")
			b.WriteString("|-------|-------|\n")
			phases := make([]string, 0, len(aps.KillChainCoverage))
			for phase := range aps.KillChainCoverage {
				phases = append(phases, phase)
			}
			sort.Strings(phases)
			for _, phase := range phases {
				b.WriteString(fmt.Sprintf("| %s | %d |\n", phase, aps.KillChainCoverage[phase]))
			}
			b.WriteString("\n")
		}

		if len(aps.MITRECoverage) > 0 {
			b.WriteString("### MITRE ATT&CK Coverage\n\n")
			for _, tech := range aps.MITRECoverage {
				b.WriteString(fmt.Sprintf("- %s\n", tech))
			}
			b.WriteString("\n")
		}

		if len(result.ThreatChains) > 0 {
			b.WriteString("### Threat Chains\n\n")
			for _, tc := range result.ThreatChains {
				b.WriteString(fmt.Sprintf("- **%s:** risk %.2f\n", tc.ID, tc.RiskScore))
				if len(tc.Threats) > 0 {
					b.WriteString(fmt.Sprintf("  - Threats: %s\n", strings.Join(tc.Threats, ", ")))
				}
				if len(tc.Path) > 0 {
					b.WriteString(fmt.Sprintf("  - Path: %s\n", strings.Join(tc.Path, " \u2192 ")))
				}
				if tc.Reasoning != "" {
					b.WriteString(fmt.Sprintf("  - Reasoning: %s\n", tc.Reasoning))
				}
			}
			b.WriteString("\n")
		}
	}

	if result.SDRISummary != "" || len(result.SDRIControls) > 0 {
		b.WriteString("## Security Design Review\n\n")
		b.WriteString(result.SDRISummary + "\n\n")

		if len(result.SDRICoverageDashboard) > 0 {
			b.WriteString("### Control Coverage Dashboard\n\n")
			b.WriteString("| Category | Coverage | Level |\n")
			b.WriteString("|----------|----------|-------|\n")
			cats := make([]string, 0, len(result.SDRICoverageDashboard))
			for cat := range result.SDRICoverageDashboard {
				cats = append(cats, cat)
			}
			sort.Strings(cats)
			for _, cat := range cats {
				cov := result.SDRICoverageDashboard[cat]
				b.WriteString(fmt.Sprintf("| %s | %.1f%% | %s |\n", cat, cov, coverageLevelString(cov)))
			}
			b.WriteString("\n")
		}

		if len(result.SDRIDesignFindings) > 0 {
			b.WriteString("### Design Findings\n\n")
			b.WriteString("| ID | Title | Severity | Category | Business Impact |\n")
			b.WriteString("|----|-------|----------|----------|----------------|\n")
			for _, f := range result.SDRIDesignFindings {
				b.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
					f.ID, f.Title, f.Severity, f.Category, f.BusinessImpact))
			}
			b.WriteString("\n")
		}

		if len(result.SDRIAchitecturalWeaknesses) > 0 {
			b.WriteString("### Architectural Weaknesses\n\n")
			b.WriteString("| ID | Pattern | Severity | Impact |\n")
			b.WriteString("|----|---------|----------|--------|\n")
			for _, w := range result.SDRIAchitecturalWeaknesses {
				b.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", w.ID, w.Pattern, w.Severity, w.Impact))
			}
			b.WriteString("\n")
		}

		if len(result.SDRIRemediations) > 0 {
			b.WriteString("### Top Remediations\n\n")
			b.WriteString("| # | Finding | Risk | Effort | Recommendation |\n")
			b.WriteString("|---|---------|------|--------|----------------|\n")
			for _, r := range result.SDRIRemediations {
				b.WriteString(fmt.Sprintf("| %d | %s | %.2f | %s | %s |\n",
					r.Priority, r.Description, r.RiskScore, r.Effort, r.Recommendation))
			}
			b.WriteString("\n")
		}

		if len(result.SDRIComplianceAlignments) > 0 {
			b.WriteString("### Compliance Alignment\n\n")
			b.WriteString("| Framework | Coverage | Status |\n")
			b.WriteString("|-----------|----------|--------|\n")
			for _, m := range result.SDRIComplianceAlignments {
				b.WriteString(fmt.Sprintf("| %s | %.1f%% | %s |\n", m.Framework, m.Coverage, m.Status))
			}
			b.WriteString("\n")
		}
	}

	if len(result.CIAREFrameworkCoverages) > 0 {
		b.WriteString("## Compliance Intelligence & Audit Readiness\n\n")

		b.WriteString("### Framework Coverage\n\n")
		b.WriteString("| Framework | Required | Observed | Missing | Coverage | Status |\n")
		b.WriteString("|-----------|----------|----------|---------|----------|--------|\n")
		for _, c := range result.CIAREFrameworkCoverages {
			b.WriteString(fmt.Sprintf("| %s | %d | %d | %d | %.1f%% | %s |\n",
				c.Framework, c.Required, c.Observed, c.Missing, c.CoveragePct, c.Status))
		}
		b.WriteString("\n")

		if len(result.CIAREAuditReadiness) > 0 {
			b.WriteString("### Audit Readiness\n\n")
			b.WriteString("| Framework | Readiness Score | Status |\n")
			b.WriteString("|-----------|-----------------|--------|\n")
			for _, a := range result.CIAREAuditReadiness {
				b.WriteString(fmt.Sprintf("| %s | %.1f%% | %s |\n", a.Framework, a.ReadinessScore, a.Status))
			}
			b.WriteString("\n")
		}

		if len(result.CIAREComplianceGaps) > 0 {
			b.WriteString("### Compliance Gaps\n\n")
			b.WriteString("| ID | Framework | Requirement | Risk |\n")
			b.WriteString("|----|-----------|-------------|------|\n")
			for _, g := range result.CIAREComplianceGaps {
				b.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", g.ID, g.Framework, g.Requirement, g.Risk))
			}
			b.WriteString("\n")
		}

		if len(result.CIAREMissingEvidences) > 0 {
			b.WriteString("### Missing Evidence\n\n")
			b.WriteString("| Framework | Control | Required Evidence |\n")
			b.WriteString("|-----------|---------|-------------------|\n")
			for _, m := range result.CIAREMissingEvidences {
				ev := strings.Join(m.Evidences, "; ")
				b.WriteString(fmt.Sprintf("| %s | %s | %s |\n", m.Framework, m.Control, ev))
			}
			b.WriteString("\n")
		}

		if len(result.CIAREAuditorQuestions) > 0 {
			b.WriteString("### Likely Auditor Questions\n\n")
			for _, q := range result.CIAREAuditorQuestions {
				b.WriteString(fmt.Sprintf("- **%s** (%s): %s\n", q.Control, q.Framework, q.Question))
			}
			b.WriteString("\n")
		}

		if len(result.CIAREControlMaturities) > 0 {
			b.WriteString("### Control Maturity\n\n")
			b.WriteString("| Domain | Level | Coverage |\n")
			b.WriteString("|--------|-------|----------|\n")
			for _, m := range result.CIAREControlMaturities {
				b.WriteString(fmt.Sprintf("| %s | %s | %.1f%% |\n", m.Domain, m.Label, m.Coverage))
			}
			b.WriteString("\n")
		}

		if len(result.CIAREComplianceNarratives) > 0 {
			b.WriteString("### Compliance Narratives\n\n")
			for _, n := range result.CIAREComplianceNarratives {
				b.WriteString(fmt.Sprintf("**%s:** %s\n\n", n.Framework, n.Narrative))
			}
		}

		if result.CIAREAuditPackage != nil && result.CIAREAuditPackage.ExecutiveSummary != "" {
			b.WriteString("### Audit Package Summary\n\n")
			b.WriteString(result.CIAREAuditPackage.ExecutiveSummary + "\n\n")
		}

		if len(result.CIAREProcurementQuestions) > 0 {
			b.WriteString("### Procurement Review Questions\n\n")
			b.WriteString("| Category | Question |\n")
			b.WriteString("|----------|----------|\n")
			for _, q := range result.CIAREProcurementQuestions {
				b.WriteString(fmt.Sprintf("| %s | %s |\n", q.Category, q.Question))
			}
			b.WriteString("\n")
		}
	}

	if result.DKPI.DomainResult.PrimaryDomain != "" {
		b.WriteString("## Domain Knowledge Intelligence\n\n")
		d := result.DKPI
		b.WriteString(fmt.Sprintf("**Detected Domain:** %s\n", d.DomainResult.PrimaryDomain))
		b.WriteString(fmt.Sprintf("**Confidence:** %.1f%%\n", d.DomainResult.Confidence))
		if len(d.DomainResult.Rationale) > 0 {
			b.WriteString("**Rationale:**\n")
			for _, r := range d.DomainResult.Rationale {
				b.WriteString(fmt.Sprintf("- %s\n", r))
			}
		}
		b.WriteString(fmt.Sprintf("**Summary:** %s\n", d.Summary))
		b.WriteString("\n")

		if len(d.Recommendations) > 0 {
			b.WriteString("### Domain Recommendations\n\n")
			for _, rec := range d.Recommendations {
				b.WriteString(fmt.Sprintf("- %s\n", rec))
			}
			b.WriteString("\n")
		}

		if len(d.InjectedThreats) > 0 {
			b.WriteString("### Domain-Specific Threats\n\n")
			b.WriteString("| ID | Name | Severity | Category | Description |\n")
			b.WriteString("|----|------|----------|----------|-------------|\n")
			for _, t := range d.InjectedThreats {
				b.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n", t.ID, t.Name, t.Severity, t.Category, t.Description))
			}
			b.WriteString("\n")
		}

		if len(d.DomainControls) > 0 {
			b.WriteString("### Domain Controls\n\n")
			b.WriteString("| ID | Name | Category | Coverage | Status |\n")
			b.WriteString("|----|------|----------|----------|--------|\n")
			for _, c := range d.DomainControls {
				b.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n", c.ID, c.Name, c.Category, c.Coverage, c.Status))
			}
			b.WriteString("\n")
		}

		if len(d.DomainCompliance) > 0 {
			b.WriteString("### Domain Compliance Frameworks\n\n")
			for _, f := range d.DomainCompliance {
				b.WriteString(fmt.Sprintf("- %s\n", f))
			}
			b.WriteString("\n")
		}

		if len(d.EvidenceReqs) > 0 {
			b.WriteString("### Evidence Requirements\n\n")
			for _, e := range d.EvidenceReqs {
				b.WriteString(fmt.Sprintf("- **%s:** %s\n", e.Control, strings.Join(e.Evidence, ", ")))
			}
			b.WriteString("\n")
		}
	}

	// ── ERN — Executive Risk Narratives ──
	if len(result.ERN.ExecutiveRisks) > 0 {
		b.WriteString("## Executive Risk Narratives\n\n")
		b.WriteString(fmt.Sprintf("**Financial Exposure:** %s — %s\n\n", result.ERN.FinancialExposure.Level, result.ERN.FinancialExposure.Rationale))
		if result.ERN.BoardSummary.Summary != "" {
			b.WriteString("### Board Summary\n\n")
			b.WriteString(result.ERN.BoardSummary.Summary + "\n\n")
		}
		if len(result.ERN.RiskThemes) > 0 {
			b.WriteString("### Risk Themes\n\n")
			b.WriteString("| Theme | Count | Severity |\n")
			b.WriteString("|-------|-------|----------|\n")
			for _, th := range result.ERN.RiskThemes {
				b.WriteString(fmt.Sprintf("| %s | %d | %s |\n", th.Name, th.RiskCount, th.Severity))
			}
			b.WriteString("\n")
		}
		if len(result.ERN.ExecutiveRisks) > 0 {
			b.WriteString("### Executive Risks\n\n")
			for _, risk := range result.ERN.ExecutiveRisks {
				b.WriteString(fmt.Sprintf("**%s** [%s] — %s\n", risk.ID, risk.Priority, risk.Title))
				b.WriteString(fmt.Sprintf("- Business Impact: %s\n", risk.BusinessImpact))
				b.WriteString(fmt.Sprintf("- Compliance Impact: %s\n", risk.ComplianceImpact))
				if len(risk.RecommendedActions) > 0 {
					b.WriteString("- Recommended Actions:\n")
					for _, a := range risk.RecommendedActions {
						b.WriteString(fmt.Sprintf("  - %s\n", a))
					}
				}
				b.WriteString("\n")
			}
		}
		if len(result.ERN.CISOBriefing.TopRisks) > 0 {
			b.WriteString("### CISO Briefing — Top Risks\n\n")
			for _, r := range result.ERN.CISOBriefing.TopRisks {
				b.WriteString(fmt.Sprintf("- %s\n", r))
			}
			b.WriteString("\n")
		}
		if len(result.ERN.RemediationRoadmap.Phase30) > 0 {
			b.WriteString("### Remediation Roadmap\n\n")
			b.WriteString("**30 Days:**\n")
			for _, item := range result.ERN.RemediationRoadmap.Phase30 {
				b.WriteString(fmt.Sprintf("- [%s] %s\n", item.Priority, item.Action))
			}
			b.WriteString("\n")
			if len(result.ERN.RemediationRoadmap.Phase90) > 0 {
				b.WriteString("**90 Days:**\n")
				for _, item := range result.ERN.RemediationRoadmap.Phase90 {
					b.WriteString(fmt.Sprintf("- [%s] %s\n", item.Priority, item.Action))
				}
				b.WriteString("\n")
			}
			if len(result.ERN.RemediationRoadmap.Phase180) > 0 {
				b.WriteString("**180 Days:**\n")
				for _, item := range result.ERN.RemediationRoadmap.Phase180 {
					b.WriteString(fmt.Sprintf("- [%s] %s\n", item.Priority, item.Action))
				}
				b.WriteString("\n")
			}
			if len(result.ERN.RemediationRoadmap.Phase12m) > 0 {
				b.WriteString("**12 Months:**\n")
				for _, item := range result.ERN.RemediationRoadmap.Phase12m {
					b.WriteString(fmt.Sprintf("- [%s] %s\n", item.Priority, item.Action))
				}
				b.WriteString("\n")
			}
		}
		if len(result.ERN.InvestmentInsights) > 0 {
			b.WriteString("### Security Investment Insights\n\n")
			for _, ii := range result.ERN.InvestmentInsights {
				b.WriteString(fmt.Sprintf("- **%s** [%s]: %s\n", ii.Area, ii.Priority, ii.Rationale))
			}
			b.WriteString("\n")
		}
		if result.ERN.DecisionSupport.Top3Actions != nil && len(result.ERN.DecisionSupport.Top3Actions) > 0 {
			b.WriteString("### CISO Decision Support — Top 3 Actions\n\n")
			for _, da := range result.ERN.DecisionSupport.Top3Actions {
				b.WriteString(fmt.Sprintf("%d. **%s** (%s impact) — %s\n", da.Rank, da.Action, da.Impact, da.Rationale))
			}
			b.WriteString("\n")
		}
		if len(result.ERN.CrownJewelClasses) > 0 {
			b.WriteString("### Crown Jewel Classification\n\n")
			b.WriteString("| Asset | Category | Label |\n")
			b.WriteString("|-------|----------|-------|\n")
			for _, cj := range result.ERN.CrownJewelClasses {
				b.WriteString(fmt.Sprintf("| %s | %s | %s |\n", cj.TechnicalName, cj.BusinessCategory, cj.BusinessLabel))
			}
			b.WriteString("\n")
		}
		b.WriteString("### Regulatory Impact Analysis\n\n")
		if len(result.ERN.RegulatoryImpacts) > 0 {
			for _, ri := range result.ERN.RegulatoryImpacts {
				b.WriteString(fmt.Sprintf("- **%s** (%s): %s\n", ri.Framework, ri.Domain, ri.Rationale))
			}
		} else {
			b.WriteString("No regulatory impacts identified.\n")
		}
		b.WriteString("\n")
	}

	// ── Portfolio Intelligence (Markdown) ──
	sampiMD := renderSAMPIReportMarkdown(result.SAMPI)
	if sampiMD != "" {
		b.WriteString(sampiMD)
	}

	// ── Decision Intelligence (Markdown) ──
	if sdiMD := renderSDIReportMarkdown(result.SDI); sdiMD != "" {
		b.WriteString(sdiMD)
	}

	// ── Digital Twin (Markdown) ──
	if sdtMD := renderSDTReportMarkdown(result.SDT); sdtMD != "" {
		b.WriteString(sdtMD)
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

	// ── Report Packs ──
	if result.ERN.ReportPacks.BoardReport != "" || result.ERN.ReportPacks.ExecutiveReport != "" || result.ERN.ReportPacks.TechnicalReport != "" {
		b.WriteString("## Report Packs\n\n")
		if result.ERN.ReportPacks.BoardReport != "" {
			b.WriteString("### Board Report\n\n")
			b.WriteString(result.ERN.ReportPacks.BoardReport)
			b.WriteString("\n")
		}
		if result.ERN.ReportPacks.ExecutiveReport != "" {
			b.WriteString("### Executive Report\n\n")
			b.WriteString(result.ERN.ReportPacks.ExecutiveReport)
			b.WriteString("\n")
		}
		if result.ERN.ReportPacks.TechnicalReport != "" {
			b.WriteString("### Technical Report\n\n")
			b.WriteString(result.ERN.ReportPacks.TechnicalReport)
			b.WriteString("\n")
		}
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
	b.WriteString(".badge-good{background:#10b981;color:#fff}\n")
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
	b.WriteString(fmt.Sprintf("<tr><td class=\"badge badge-critical\">Critical</td><td>%d</td></tr>\n", result.CriticalCount))
	b.WriteString(fmt.Sprintf("<tr><td class=\"badge badge-high\">High</td><td>%d</td></tr>\n", result.HighCount))
	b.WriteString(fmt.Sprintf("<tr><td class=\"badge badge-medium\">Medium</td><td>%d</td></tr>\n", result.MediumCount))
	b.WriteString(fmt.Sprintf("<tr><td class=\"badge badge-low\">Low</td><td>%d</td></tr>\n", result.LowCount))
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

	if len(result.AttackPaths) > 0 {
		b.WriteString("<h2>Attack Path Discovery</h2>\n")
		aps := result.AttackPathSummary

		b.WriteString("<h3>Summary</h3>\n")
		b.WriteString("<table>\n")
		b.WriteString(fmt.Sprintf("<tr><td>Total Attack Paths</td><td>%d</td></tr>\n", aps.TotalAttackPaths))
		b.WriteString(fmt.Sprintf("<tr><td>Critical</td><td>%d</td></tr>\n", aps.CriticalCount))
		b.WriteString(fmt.Sprintf("<tr><td>High</td><td>%d</td></tr>\n", aps.HighCount))
		b.WriteString(fmt.Sprintf("<tr><td>Medium</td><td>%d</td></tr>\n", aps.MediumCount))
		b.WriteString(fmt.Sprintf("<tr><td>Low</td><td>%d</td></tr>\n", aps.LowCount))
		b.WriteString(fmt.Sprintf("<tr><td>Threat Chains</td><td>%d</td></tr>\n", aps.ThreatChainCount))
		if aps.SummaryText != "" {
			b.WriteString(fmt.Sprintf("<tr><td>Summary</td><td>%s</td></tr>\n", aps.SummaryText))
		}
		b.WriteString("</table>\n")

		b.WriteString("<h3>Top Attack Paths</h3>\n")
		b.WriteString("<table>\n<tr><th>Path</th><th>Entry</th><th>Target</th><th>Risk</th><th>Detection</th><th>Business Impact</th></tr>\n")
		for _, p := range result.AttackPaths {
			riskLabel := riskLevelForAPDScore(p.RiskScore)
			badgeClass := "badge-medium"
			if riskLabel == "Critical" {
				badgeClass = "badge-critical"
			} else if riskLabel == "High" {
				badgeClass = "badge-high"
			} else if riskLabel == "Low" {
				badgeClass = "badge-low"
			}
			cellBadge := fmt.Sprintf("<span class=\"badge %s\">%s</span>", badgeClass, riskLabel)
			b.WriteString(fmt.Sprintf("<tr><td><strong>%s</strong></td><td>%s</td><td>%s</td><td>%.2f %s</td><td>%s</td><td>%s</td></tr>\n",
				p.Name, p.EntryPoint, p.TargetAsset, p.RiskScore, cellBadge, p.DetectionDifficulty, p.BusinessImpact))
		}
		b.WriteString("</table>\n")

		if len(aps.KillChainCoverage) > 0 {
			b.WriteString("<h3>Kill Chain Coverage</h3>\n")
			b.WriteString("<table>\n<tr><th>Phase</th><th>Paths</th></tr>\n")
			phases := make([]string, 0, len(aps.KillChainCoverage))
			for phase := range aps.KillChainCoverage {
				phases = append(phases, phase)
			}
			sort.Strings(phases)
			for _, phase := range phases {
				b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%d</td></tr>\n", phase, aps.KillChainCoverage[phase]))
			}
			b.WriteString("</table>\n")
		}

		if len(aps.MITRECoverage) > 0 {
			b.WriteString("<h3>MITRE ATT&CK Coverage</h3>\n<ul>\n")
			for _, tech := range aps.MITRECoverage {
				b.WriteString(fmt.Sprintf("<li>%s</li>\n", tech))
			}
			b.WriteString("</ul>\n")
		}

		if len(result.ThreatChains) > 0 {
			b.WriteString("<h3>Threat Chains</h3>\n<ul>\n")
			for _, tc := range result.ThreatChains {
				b.WriteString(fmt.Sprintf("<li><strong>%s:</strong> risk %.2f", tc.ID, tc.RiskScore))
				if len(tc.Threats) > 0 {
					b.WriteString(fmt.Sprintf(" — Threats: %s", strings.Join(tc.Threats, ", ")))
				}
				if len(tc.Path) > 0 {
					b.WriteString(fmt.Sprintf(" — Path: %s", strings.Join(tc.Path, " \u2192 ")))
				}
				b.WriteString("</li>\n")
			}
			b.WriteString("</ul>\n")
		}
	}

	if result.SDRISummary != "" || len(result.SDRIControls) > 0 {
		b.WriteString("<h2>Security Design Review</h2>\n")
		b.WriteString(fmt.Sprintf("<p>%s</p>\n", result.SDRISummary))

		if len(result.SDRICoverageDashboard) > 0 {
			b.WriteString("<h3>Control Coverage Dashboard</h3>\n")
			b.WriteString("<table>\n<tr><th>Category</th><th>Coverage</th><th>Level</th></tr>\n")
			cats := make([]string, 0, len(result.SDRICoverageDashboard))
			for cat := range result.SDRICoverageDashboard {
				cats = append(cats, cat)
			}
			sort.Strings(cats)
			for _, cat := range cats {
				cov := result.SDRICoverageDashboard[cat]
				level := coverageLevelString(cov)
				badgeClass := "badge-medium"
				if level == "Excellent" || level == "Strong" {
					badgeClass = "badge-good"
				} else if level == "Fair" || level == "Poor" {
					badgeClass = "badge-critical"
				}
				b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%.1f%%</td><td><span class=\"badge %s\">%s</span></td></tr>\n",
					cat, cov, badgeClass, level))
			}
			b.WriteString("</table>\n")
		}

		if len(result.SDRIDesignFindings) > 0 {
			b.WriteString("<h3>Design Findings</h3>\n")
			b.WriteString("<table>\n<tr><th>ID</th><th>Title</th><th>Severity</th><th>Category</th><th>Business Impact</th></tr>\n")
			for _, f := range result.SDRIDesignFindings {
				sevBadge := fmt.Sprintf("<span class=\"badge badge-%s\">%s</span>", severityCSSClass(f.Severity), f.Severity)
				b.WriteString(fmt.Sprintf("<tr><td>%s</td><td><strong>%s</strong></td><td>%s</td><td>%s</td><td>%s</td></tr>\n",
					f.ID, f.Title, sevBadge, f.Category, f.BusinessImpact))
			}
			b.WriteString("</table>\n")
		}

		if len(result.SDRIAchitecturalWeaknesses) > 0 {
			b.WriteString("<h3>Architectural Weaknesses</h3>\n")
			b.WriteString("<table>\n<tr><th>ID</th><th>Pattern</th><th>Severity</th><th>Impact</th></tr>\n")
			for _, w := range result.SDRIAchitecturalWeaknesses {
				sevBadge := fmt.Sprintf("<span class=\"badge badge-%s\">%s</span>", severityCSSClass(w.Severity), w.Severity)
				b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>\n",
					w.ID, w.Pattern, sevBadge, w.Impact))
			}
			b.WriteString("</table>\n")
		}

		if len(result.SDRIRemediations) > 0 {
			b.WriteString("<h3>Top Remediations</h3>\n")
			b.WriteString("<table>\n<tr><th>#</th><th>Finding</th><th>Risk</th><th>Effort</th><th>Recommendation</th></tr>\n")
			for _, r := range result.SDRIRemediations {
				b.WriteString(fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%.2f</td><td>%s</td><td>%s</td></tr>\n",
					r.Priority, r.Description, r.RiskScore, r.Effort, r.Recommendation))
			}
			b.WriteString("</table>\n")
		}

		if len(result.SDRIComplianceAlignments) > 0 {
			b.WriteString("<h3>Compliance Alignment</h3>\n")
			b.WriteString("<table>\n<tr><th>Framework</th><th>Coverage</th><th>Status</th></tr>\n")
			for _, m := range result.SDRIComplianceAlignments {
				badgeClass := "badge-medium"
				if m.Status == "Excellent" || m.Status == "Strong" {
					badgeClass = "badge-good"
				} else if m.Status == "Fair" || m.Status == "Poor" {
					badgeClass = "badge-critical"
				}
				b.WriteString(fmt.Sprintf("<tr><td><strong>%s</strong></td><td>%.1f%%</td><td><span class=\"badge %s\">%s</span></td></tr>\n",
					m.Framework, m.Coverage, badgeClass, m.Status))
			}
			b.WriteString("</table>\n")
		}
	}

	if len(result.CIAREFrameworkCoverages) > 0 {
		b.WriteString("<h2>Compliance Intelligence & Audit Readiness</h2>\n")

		b.WriteString("<h3>Framework Coverage</h3>\n")
		b.WriteString("<table>\n<tr><th>Framework</th><th>Required</th><th>Observed</th><th>Missing</th><th>Coverage</th><th>Status</th></tr>\n")
		for _, c := range result.CIAREFrameworkCoverages {
			badgeClass := "badge-medium"
			if c.Status == "Excellent" || c.Status == "Strong" {
				badgeClass = "badge-good"
			} else if c.Status == "Fair" || c.Status == "Weak" || c.Status == "Poor" {
				badgeClass = "badge-critical"
			}
			b.WriteString(fmt.Sprintf("<tr><td><strong>%s</strong></td><td>%d</td><td>%d</td><td>%d</td><td>%.1f%%</td><td><span class=\"badge %s\">%s</span></td></tr>\n",
				c.Framework, c.Required, c.Observed, c.Missing, c.CoveragePct, badgeClass, c.Status))
		}
		b.WriteString("</table>\n")

		if len(result.CIAREAuditReadiness) > 0 {
			b.WriteString("<h3>Audit Readiness</h3>\n")
			b.WriteString("<table>\n<tr><th>Framework</th><th>Readiness Score</th><th>Status</th></tr>\n")
			for _, a := range result.CIAREAuditReadiness {
				badgeClass := "badge-medium"
				if a.Status == "Excellent" || a.Status == "Strong" {
					badgeClass = "badge-good"
				} else if a.Status == "Fair" || a.Status == "Weak" || a.Status == "Poor" {
					badgeClass = "badge-critical"
				}
				b.WriteString(fmt.Sprintf("<tr><td><strong>%s</strong></td><td>%.1f%%</td><td><span class=\"badge %s\">%s</span></td></tr>\n",
					a.Framework, a.ReadinessScore, badgeClass, a.Status))
			}
			b.WriteString("</table>\n")
		}

		if len(result.CIAREComplianceGaps) > 0 {
			b.WriteString("<h3>Compliance Gaps</h3>\n")
			b.WriteString("<table>\n<tr><th>ID</th><th>Framework</th><th>Requirement</th><th>Risk</th></tr>\n")
			for _, g := range result.CIAREComplianceGaps {
				sevBadge := fmt.Sprintf("<span class=\"badge badge-%s\">%s</span>", severityCSSClass(g.Risk), g.Risk)
				b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>\n",
					g.ID, g.Framework, g.Requirement, sevBadge))
			}
			b.WriteString("</table>\n")
		}

		if len(result.CIAREMissingEvidences) > 0 {
			b.WriteString("<h3>Missing Evidence</h3>\n")
			b.WriteString("<table>\n<tr><th>Framework</th><th>Control</th><th>Required Evidence</th></tr>\n")
			for _, m := range result.CIAREMissingEvidences {
				ev := strings.Join(m.Evidences, "; ")
				b.WriteString(fmt.Sprintf("<tr><td>%s</td><td><strong>%s</strong></td><td>%s</td></tr>\n",
					m.Framework, m.Control, ev))
			}
			b.WriteString("</table>\n")
		}

		if len(result.CIAREAuditorQuestions) > 0 {
			b.WriteString("<h3>Likely Auditor Questions</h3>\n")
			b.WriteString("<ul>\n")
			for _, q := range result.CIAREAuditorQuestions {
				b.WriteString(fmt.Sprintf("<li><strong>%s</strong> (%s): %s</li>\n", q.Control, q.Framework, q.Question))
			}
			b.WriteString("</ul>\n")
		}

		if len(result.CIAREControlMaturities) > 0 {
			b.WriteString("<h3>Control Maturity</h3>\n")
			b.WriteString("<table>\n<tr><th>Domain</th><th>Level</th><th>Coverage</th></tr>\n")
			for _, m := range result.CIAREControlMaturities {
				b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%.1f%%</td></tr>\n",
					m.Domain, m.Label, m.Coverage))
			}
			b.WriteString("</table>\n")
		}

		if len(result.CIAREComplianceNarratives) > 0 {
			b.WriteString("<h3>Compliance Narratives</h3>\n")
			for _, n := range result.CIAREComplianceNarratives {
				b.WriteString(fmt.Sprintf("<p><strong>%s:</strong> %s</p>\n", n.Framework, n.Narrative))
			}
		}

		if result.CIAREAuditPackage != nil && result.CIAREAuditPackage.ExecutiveSummary != "" {
			b.WriteString("<h3>Audit Package Summary</h3>\n")
			b.WriteString(fmt.Sprintf("<p>%s</p>\n", result.CIAREAuditPackage.ExecutiveSummary))
		}

		if len(result.CIAREProcurementQuestions) > 0 {
			b.WriteString("<h3>Procurement Review Questions</h3>\n")
			b.WriteString("<table>\n<tr><th>Category</th><th>Question</th></tr>\n")
			for _, q := range result.CIAREProcurementQuestions {
				b.WriteString(fmt.Sprintf("<tr><td><strong>%s</strong></td><td>%s</td></tr>\n", q.Category, q.Question))
			}
			b.WriteString("</table>\n")
		}
	}

	if result.DKPI.DomainResult.PrimaryDomain != "" {
		b.WriteString("<h2>Domain Knowledge Intelligence</h2>\n")
		d := result.DKPI
		b.WriteString(fmt.Sprintf("<p><strong>Detected Domain:</strong> %s<br>\n", d.DomainResult.PrimaryDomain))
		b.WriteString(fmt.Sprintf("<strong>Confidence:</strong> %.1f%%<br>\n", d.DomainResult.Confidence))
		b.WriteString(fmt.Sprintf("<strong>Summary:</strong> %s</p>\n", d.Summary))
		if len(d.DomainResult.Rationale) > 0 {
			b.WriteString("<h3>Detection Rationale</h3>\n<ul>\n")
			for _, r := range d.DomainResult.Rationale {
				b.WriteString(fmt.Sprintf("<li>%s</li>\n", r))
			}
			b.WriteString("</ul>\n")
		}
		if len(d.Recommendations) > 0 {
			b.WriteString("<h3>Domain Recommendations</h3>\n<ul>\n")
			for _, rec := range d.Recommendations {
				b.WriteString(fmt.Sprintf("<li>%s</li>\n", rec))
			}
			b.WriteString("</ul>\n")
		}
		if len(d.InjectedThreats) > 0 {
			b.WriteString("<h3>Domain-Specific Threats</h3>\n")
			b.WriteString("<table>\n<tr><th>ID</th><th>Name</th><th>Severity</th><th>Category</th><th>Description</th></tr>\n")
			for _, t := range d.InjectedThreats {
				b.WriteString(fmt.Sprintf("<tr><td>%s</td><td><strong>%s</strong></td><td>%s</td><td>%s</td><td>%s</td></tr>\n",
					t.ID, t.Name, t.Severity, t.Category, t.Description))
			}
			b.WriteString("</table>\n")
		}
		if len(d.DomainControls) > 0 {
			b.WriteString("<h3>Domain Controls</h3>\n")
			b.WriteString("<table>\n<tr><th>ID</th><th>Name</th><th>Category</th><th>Coverage</th><th>Status</th></tr>\n")
			for _, c := range d.DomainControls {
				b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>\n",
					c.ID, c.Name, c.Category, c.Coverage, c.Status))
			}
			b.WriteString("</table>\n")
		}
	}

	// ── ERN — Executive Risk Narratives (HTML) ──
	if len(result.ERN.ExecutiveRisks) > 0 {
		b.WriteString("<h2>Executive Risk Narratives</h2>\n")
		b.WriteString(fmt.Sprintf("<p><strong>Financial Exposure:</strong> %s — %s</p>\n", result.ERN.FinancialExposure.Level, result.ERN.FinancialExposure.Rationale))
		if result.ERN.BoardSummary.Summary != "" {
			b.WriteString("<h3>Board Summary</h3>\n")
			b.WriteString(fmt.Sprintf("<p>%s</p>\n", result.ERN.BoardSummary.Summary))
		}
		if len(result.ERN.RiskThemes) > 0 {
			b.WriteString("<h3>Risk Themes</h3>\n<table>\n<tr><th>Theme</th><th>Count</th><th>Severity</th></tr>\n")
			for _, th := range result.ERN.RiskThemes {
				b.WriteString(fmt.Sprintf("<tr><td>%s</td><td>%d</td><td>%s</td></tr>\n", th.Name, th.RiskCount, th.Severity))
			}
			b.WriteString("</table>\n")
		}
		if len(result.ERN.ExecutiveRisks) > 0 {
			b.WriteString("<h3>Executive Risks</h3>\n")
			for _, risk := range result.ERN.ExecutiveRisks {
				b.WriteString(fmt.Sprintf("<p><strong>%s</strong> [%s] — %s<br>\n", risk.ID, risk.Priority, risk.Title))
				b.WriteString(fmt.Sprintf("Business Impact: %s<br>\n", risk.BusinessImpact))
				b.WriteString(fmt.Sprintf("Compliance Impact: %s</p>\n", risk.ComplianceImpact))
				if len(risk.RecommendedActions) > 0 {
					b.WriteString("<ul>\n")
					for _, a := range risk.RecommendedActions {
						b.WriteString(fmt.Sprintf("<li>%s</li>\n", a))
					}
					b.WriteString("</ul>\n")
				}
			}
		}
		if len(result.ERN.CISOBriefing.TopRisks) > 0 {
			b.WriteString("<h3>CISO Briefing — Top Risks</h3>\n<ul>\n")
			for _, r := range result.ERN.CISOBriefing.TopRisks {
				b.WriteString(fmt.Sprintf("<li>%s</li>\n", r))
			}
			b.WriteString("</ul>\n")
		}
		if len(result.ERN.RemediationRoadmap.Phase30) > 0 {
			b.WriteString("<h3>Remediation Roadmap</h3>\n")
			b.WriteString("<p><strong>30 Days:</strong></p>\n<ul>\n")
			for _, item := range result.ERN.RemediationRoadmap.Phase30 {
				b.WriteString(fmt.Sprintf("<li>[%s] %s</li>\n", item.Priority, item.Action))
			}
			b.WriteString("</ul>\n")
		}
		if len(result.ERN.InvestmentInsights) > 0 {
			b.WriteString("<h3>Security Investment Insights</h3>\n<ul>\n")
			for _, ii := range result.ERN.InvestmentInsights {
				b.WriteString(fmt.Sprintf("<li><strong>%s</strong> [%s]: %s</li>\n", ii.Area, ii.Priority, ii.Rationale))
			}
			b.WriteString("</ul>\n")
		}
		if result.ERN.DecisionSupport.Top3Actions != nil && len(result.ERN.DecisionSupport.Top3Actions) > 0 {
			b.WriteString("<h3>CISO Decision Support — Top 3 Actions</h3>\n<ol>\n")
			for _, da := range result.ERN.DecisionSupport.Top3Actions {
				b.WriteString(fmt.Sprintf("<li><strong>%s</strong> (%s impact) — %s</li>\n", da.Action, da.Impact, da.Rationale))
			}
			b.WriteString("</ol>\n")
		}
	}

	// ── Portfolio Intelligence (HTML) ──
	sampiMD := renderSAMPIReportMarkdown(result.SAMPI)
	if sampiMD != "" {
		b.WriteString(markdownToHTML(sampiMD))
	}

	// ── Decision Intelligence (HTML) ──
	if sdiMD := renderSDIReportMarkdown(result.SDI); sdiMD != "" {
		b.WriteString(markdownToHTML(sdiMD))
	}

	// ── Digital Twin (HTML) ──
	if sdtMD := renderSDTReportMarkdown(result.SDT); sdtMD != "" {
		b.WriteString(markdownToHTML(sdtMD))
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

	// ── Report Packs (HTML) ──
	if result.ERN.ReportPacks.BoardReport != "" || result.ERN.ReportPacks.ExecutiveReport != "" || result.ERN.ReportPacks.TechnicalReport != "" {
		b.WriteString("<h2>Report Packs</h2>\n")
		if result.ERN.ReportPacks.BoardReport != "" {
			b.WriteString("<h3>Board Report</h3>\n")
			b.WriteString("<div class=\"report-pack\">\n")
			b.WriteString(markdownToHTML(result.ERN.ReportPacks.BoardReport))
			b.WriteString("</div>\n")
		}
		if result.ERN.ReportPacks.ExecutiveReport != "" {
			b.WriteString("<h3>Executive Report</h3>\n")
			b.WriteString("<div class=\"report-pack\">\n")
			b.WriteString(markdownToHTML(result.ERN.ReportPacks.ExecutiveReport))
			b.WriteString("</div>\n")
		}
		if result.ERN.ReportPacks.TechnicalReport != "" {
			b.WriteString("<h3>Technical Report</h3>\n")
			b.WriteString("<div class=\"report-pack\">\n")
			b.WriteString(markdownToHTML(result.ERN.ReportPacks.TechnicalReport))
			b.WriteString("</div>\n")
		}
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

	pdf.CellFormat(190, 7, fmt.Sprintf("Critical: %d", result.CriticalCount), "", 1, "L", false, 0, "")
	pdf.CellFormat(190, 7, fmt.Sprintf("High Risk: %d", result.HighCount), "", 1, "L", false, 0, "")
	pdf.CellFormat(190, 7, fmt.Sprintf("Medium Risk: %d", result.MediumCount), "", 1, "L", false, 0, "")
	pdf.CellFormat(190, 7, fmt.Sprintf("Low Risk: %d", result.LowCount), "", 1, "L", false, 0, "")
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

	// Attack Path Discovery
	if len(result.AttackPaths) > 0 {
		pdf.AddPage()
		pdf.SetFont("Helvetica", "B", 16)
		pdf.CellFormat(190, 12, "Attack Path Discovery", "", 1, "L", false, 0, "")
		pdf.Ln(4)

		aps := result.AttackPathSummary
		pdf.SetFont("Helvetica", "", 10)
		pdf.CellFormat(190, 7, fmt.Sprintf("Total Attack Paths: %d", aps.TotalAttackPaths), "", 1, "L", false, 0, "")
		pdf.CellFormat(190, 7, fmt.Sprintf("Critical: %d | High: %d | Medium: %d | Low: %d", aps.CriticalCount, aps.HighCount, aps.MediumCount, aps.LowCount), "", 1, "L", false, 0, "")
		pdf.CellFormat(190, 7, fmt.Sprintf("Threat Chains: %d", aps.ThreatChainCount), "", 1, "L", false, 0, "")
		if aps.SummaryText != "" {
			pdf.SetFont("Helvetica", "I", 9)
			pdf.MultiCell(190, 5, "Summary: "+aps.SummaryText, "", "L", false)
		}
		pdf.Ln(4)

		pdf.SetFont("Helvetica", "B", 12)
		pdf.CellFormat(190, 7, "Top Attack Paths", "", 1, "L", false, 0, "")
		pdf.Ln(2)
		pdf.SetFont("Helvetica", "", 9)
		for _, p := range result.AttackPaths {
			riskLabel := riskLevelForAPDScore(p.RiskScore)
			pdf.SetFont("Helvetica", "B", 9)
			pdf.MultiCell(190, 4, fmt.Sprintf("%s: %s -> %s (risk: %.2f, %s)", p.Name, p.EntryPoint, p.TargetAsset, p.RiskScore, riskLabel), "", "L", false)
			pdf.SetFont("Helvetica", "", 8)
			pdf.CellFormat(190, 4, fmt.Sprintf("  Detection: %s | Impact: %s", p.DetectionDifficulty, p.BusinessImpact), "", 1, "L", false, 0, "")
			if len(p.KillChainPhases) > 0 {
				pdf.CellFormat(190, 4, "  Kill Chain: "+strings.Join(p.KillChainPhases, ", "), "", 1, "L", false, 0, "")
			}
			if len(p.MITREATTACK) > 0 {
				pdf.CellFormat(190, 4, "  MITRE: "+strings.Join(p.MITREATTACK, ", "), "", 1, "L", false, 0, "")
			}
			pdf.Ln(3)
		}
		pdf.Ln(2)

		if len(aps.KillChainCoverage) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "Kill Chain Coverage", "", 1, "L", false, 0, "")
			pdf.Ln(2)
			pdf.SetFont("Helvetica", "", 9)
			phases := make([]string, 0, len(aps.KillChainCoverage))
			for phase := range aps.KillChainCoverage {
				phases = append(phases, phase)
			}
			sort.Strings(phases)
			for _, phase := range phases {
				pdf.CellFormat(190, 5, fmt.Sprintf("  %s: %d paths", phase, aps.KillChainCoverage[phase]), "", 1, "L", false, 0, "")
			}
			pdf.Ln(2)
		}

		if len(aps.MITRECoverage) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "MITRE ATT&CK Coverage", "", 1, "L", false, 0, "")
			pdf.Ln(2)
			pdf.SetFont("Helvetica", "", 9)
			for _, tech := range aps.MITRECoverage {
				pdf.CellFormat(190, 5, "  - "+tech, "", 1, "L", false, 0, "")
			}
			pdf.Ln(2)
		}

		if len(result.ThreatChains) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "Threat Chains", "", 1, "L", false, 0, "")
			pdf.Ln(2)
			pdf.SetFont("Helvetica", "", 9)
			for _, tc := range result.ThreatChains {
				pdf.SetFont("Helvetica", "B", 9)
				pdf.MultiCell(190, 4, fmt.Sprintf("%s (risk: %.2f)", tc.ID, tc.RiskScore), "", "L", false)
				pdf.SetFont("Helvetica", "", 8)
				if len(tc.Threats) > 0 {
					pdf.CellFormat(190, 4, "  Threats: "+strings.Join(tc.Threats, ", "), "", 1, "L", false, 0, "")
				}
				if len(tc.Path) > 0 {
					pdf.CellFormat(190, 4, "  Path: "+strings.Join(tc.Path, " -> "), "", 1, "L", false, 0, "")
				}
				pdf.Ln(2)
			}
		}
	}

	// Security Design Review
	if result.SDRISummary != "" || len(result.SDRIControls) > 0 {
		pdf.AddPage()
		pdf.SetFont("Helvetica", "B", 16)
		pdf.CellFormat(190, 12, "Security Design Review", "", 1, "L", false, 0, "")
		pdf.Ln(4)
		pdf.SetFont("Helvetica", "", 10)
		pdf.MultiCell(190, 5, result.SDRISummary, "", "L", false)
		pdf.Ln(4)

		if len(result.SDRICoverageDashboard) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "Control Coverage Dashboard", "", 1, "L", false, 0, "")
			pdf.Ln(2)
			pdf.SetFont("Helvetica", "", 9)
			cats := make([]string, 0, len(result.SDRICoverageDashboard))
			for cat := range result.SDRICoverageDashboard {
				cats = append(cats, cat)
			}
			sort.Strings(cats)
			for _, cat := range cats {
				cov := result.SDRICoverageDashboard[cat]
				level := coverageLevelString(cov)
				pdf.CellFormat(190, 5, fmt.Sprintf("  %s: %.1f%% (%s)", cat, cov, level), "", 1, "L", false, 0, "")
			}
			pdf.Ln(4)
		}

		if len(result.SDRIDesignFindings) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, fmt.Sprintf("Design Findings (%d)", len(result.SDRIDesignFindings)), "", 1, "L", false, 0, "")
			pdf.Ln(2)
			pdf.SetFont("Helvetica", "", 8)
			for _, f := range result.SDRIDesignFindings {
				pdf.SetFont("Helvetica", "B", 8)
				pdf.CellFormat(190, 4, fmt.Sprintf("%s [%s] %s", f.ID, f.Severity, f.Title), "", 1, "L", false, 0, "")
				pdf.SetFont("Helvetica", "", 7)
				pdf.MultiCell(190, 3, "  Impact: "+f.BusinessImpact, "", "L", false)
				if len(f.AffectedComponents) > 0 {
					pdf.CellFormat(190, 3, "  Components: "+strings.Join(f.AffectedComponents, ", "), "", 1, "L", false, 0, "")
				}
			}
			pdf.Ln(4)
		}

		if len(result.SDRIAchitecturalWeaknesses) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, fmt.Sprintf("Architectural Weaknesses (%d)", len(result.SDRIAchitecturalWeaknesses)), "", 1, "L", false, 0, "")
			pdf.Ln(2)
			pdf.SetFont("Helvetica", "", 8)
			for _, w := range result.SDRIAchitecturalWeaknesses {
				pdf.SetFont("Helvetica", "B", 8)
				pdf.CellFormat(190, 4, fmt.Sprintf("%s [%s] %s", w.ID, w.Severity, w.Pattern), "", 1, "L", false, 0, "")
				pdf.SetFont("Helvetica", "", 7)
				pdf.MultiCell(190, 3, "  "+w.Description, "", "L", false)
			}
			pdf.Ln(4)
		}

		if len(result.SDRIRemediations) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, fmt.Sprintf("Top Remediations (%d)", len(result.SDRIRemediations)), "", 1, "L", false, 0, "")
			pdf.Ln(2)
			pdf.SetFont("Helvetica", "", 8)
			for _, r := range result.SDRIRemediations {
				pdf.SetFont("Helvetica", "B", 8)
				pdf.CellFormat(190, 4, fmt.Sprintf("#%d %s (risk: %.2f, effort: %s)", r.Priority, r.Description, r.RiskScore, r.Effort), "", 1, "L", false, 0, "")
				pdf.SetFont("Helvetica", "", 7)
				pdf.MultiCell(190, 3, "  "+r.Recommendation, "", "L", false)
			}
			pdf.Ln(4)
		}

		if len(result.SDRIComplianceAlignments) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "Compliance Alignment", "", 1, "L", false, 0, "")
			pdf.Ln(2)
			pdf.SetFont("Helvetica", "", 9)
			for _, m := range result.SDRIComplianceAlignments {
				pdf.CellFormat(190, 5, fmt.Sprintf("  %s: %.1f%% (%s) - %d controls", m.Framework, m.Coverage, m.Status, len(m.Controls)), "", 1, "L", false, 0, "")
			}
			pdf.Ln(4)
		}
	}

	// Compliance Intelligence & Audit Readiness
	if len(result.CIAREFrameworkCoverages) > 0 {
		pdf.AddPage()
		pdf.SetFont("Helvetica", "B", 16)
		pdf.CellFormat(190, 12, "Compliance Intelligence & Audit Readiness", "", 1, "L", false, 0, "")
		pdf.Ln(4)

		pdf.SetFont("Helvetica", "B", 12)
		pdf.CellFormat(190, 7, "Framework Coverage", "", 1, "L", false, 0, "")
		pdf.Ln(2)
		pdf.SetFont("Helvetica", "", 8)
		for _, c := range result.CIAREFrameworkCoverages {
			pdf.CellFormat(190, 4, fmt.Sprintf("  %s: %d/%d observed (%.1f%%) - %s",
				c.Framework, c.Observed, c.Required, c.CoveragePct, c.Status), "", 1, "L", false, 0, "")
		}
		pdf.Ln(4)

		if len(result.CIAREAuditReadiness) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "Audit Readiness", "", 1, "L", false, 0, "")
			pdf.Ln(2)
			pdf.SetFont("Helvetica", "", 8)
			for _, a := range result.CIAREAuditReadiness {
				pdf.CellFormat(190, 4, fmt.Sprintf("  %s: %.1f%% - %s", a.Framework, a.ReadinessScore, a.Status), "", 1, "L", false, 0, "")
			}
			pdf.Ln(4)
		}

		if len(result.CIAREComplianceGaps) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, fmt.Sprintf("Compliance Gaps (%d)", len(result.CIAREComplianceGaps)), "", 1, "L", false, 0, "")
			pdf.Ln(2)
			pdf.SetFont("Helvetica", "", 8)
			for i, g := range result.CIAREComplianceGaps {
				if i >= 10 {
					break
				}
				pdf.CellFormat(190, 4, fmt.Sprintf("  [%s] %s - %s (%s)", g.ID, g.Requirement, g.Framework, g.Risk), "", 1, "L", false, 0, "")
			}
			pdf.Ln(4)
		}

		if len(result.CIAREMissingEvidences) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "Missing Evidence", "", 1, "L", false, 0, "")
			pdf.Ln(2)
			pdf.SetFont("Helvetica", "", 8)
			for i, m := range result.CIAREMissingEvidences {
				if i >= 10 {
					break
				}
				ev := strings.Join(m.Evidences, "; ")
				pdf.MultiCell(190, 3, fmt.Sprintf("  %s (%s): %s", m.Control, m.Framework, ev), "", "L", false)
			}
			pdf.Ln(4)
		}

		if len(result.CIAREAuditorQuestions) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "Likely Auditor Questions", "", 1, "L", false, 0, "")
			pdf.Ln(2)
			pdf.SetFont("Helvetica", "", 8)
			for i, q := range result.CIAREAuditorQuestions {
				if i >= 10 {
					break
				}
				pdf.MultiCell(190, 3, fmt.Sprintf("  %s (%s): %s", q.Control, q.Framework, q.Question), "", "L", false)
			}
			pdf.Ln(4)
		}

		if len(result.CIAREControlMaturities) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "Control Maturity", "", 1, "L", false, 0, "")
			pdf.Ln(2)
			pdf.SetFont("Helvetica", "", 8)
			for _, m := range result.CIAREControlMaturities {
				pdf.CellFormat(190, 4, fmt.Sprintf("  %s: %s (%.1f%%)",
					m.Domain, m.Label, m.Coverage), "", 1, "L", false, 0, "")
			}
			pdf.Ln(4)
		}

		if len(result.CIAREComplianceNarratives) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "Compliance Narratives", "", 1, "L", false, 0, "")
			pdf.Ln(2)
			pdf.SetFont("Helvetica", "I", 8)
			for _, n := range result.CIAREComplianceNarratives {
				pdf.MultiCell(190, 3, fmt.Sprintf("%s: %s", n.Framework, n.Narrative), "", "L", false)
			}
			pdf.Ln(4)
		}

		if result.CIAREAuditPackage != nil && result.CIAREAuditPackage.ExecutiveSummary != "" {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "Audit Package", "", 1, "L", false, 0, "")
			pdf.Ln(2)
			pdf.SetFont("Helvetica", "I", 8)
			pdf.MultiCell(190, 3, result.CIAREAuditPackage.ExecutiveSummary, "", "L", false)
			pdf.Ln(4)
		}

		if len(result.CIAREProcurementQuestions) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "Procurement Review Questions", "", 1, "L", false, 0, "")
			pdf.Ln(2)
			pdf.SetFont("Helvetica", "", 8)
			for i, q := range result.CIAREProcurementQuestions {
				if i >= 15 {
					break
				}
				pdf.CellFormat(190, 3, fmt.Sprintf("  [%s] %s", q.Category, q.Question), "", 1, "L", false, 0, "")
			}
			pdf.Ln(4)
		}
	}

	// Domain Knowledge Intelligence
	if result.DKPI.DomainResult.PrimaryDomain != "" {
		pdf.AddPage()
		pdf.SetFont("Helvetica", "B", 16)
		pdf.CellFormat(190, 12, "Domain Knowledge Intelligence", "", 1, "L", false, 0, "")
		pdf.Ln(4)
		d := result.DKPI
		pdf.SetFont("Helvetica", "", 10)
		pdf.CellFormat(190, 7, fmt.Sprintf("Detected Domain: %s", d.DomainResult.PrimaryDomain), "", 1, "L", false, 0, "")
		pdf.CellFormat(190, 7, fmt.Sprintf("Confidence: %.1f%%", d.DomainResult.Confidence), "", 1, "L", false, 0, "")
		pdf.CellFormat(190, 7, "Summary: "+d.Summary, "", 1, "L", false, 0, "")
		pdf.Ln(4)

		if len(d.Recommendations) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "Domain Recommendations", "", 1, "L", false, 0, "")
			pdf.Ln(2)
			pdf.SetFont("Helvetica", "", 8)
			for _, rec := range d.Recommendations {
				pdf.MultiCell(190, 4, "  - "+rec, "", "L", false)
			}
			pdf.Ln(4)
		}

		if len(d.InjectedThreats) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, fmt.Sprintf("Domain Threats (%d)", len(d.InjectedThreats)), "", 1, "L", false, 0, "")
			pdf.Ln(2)
			pdf.SetFont("Helvetica", "", 8)
			for _, t := range d.InjectedThreats {
				pdf.SetFont("Helvetica", "B", 8)
				pdf.CellFormat(190, 4, fmt.Sprintf("%s [%s] %s", t.ID, t.Severity, t.Name), "", 1, "L", false, 0, "")
				pdf.SetFont("Helvetica", "", 7)
				pdf.MultiCell(190, 3, "  "+t.Description, "", "L", false)
			}
			pdf.Ln(4)
		}

		if len(d.DomainControls) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "Domain Controls", "", 1, "L", false, 0, "")
			pdf.Ln(2)
			pdf.SetFont("Helvetica", "", 8)
			for _, c := range d.DomainControls {
				pdf.CellFormat(190, 4, fmt.Sprintf("  %s: %s (%s/%s)", c.ID, c.Name, c.Coverage, c.Status), "", 1, "L", false, 0, "")
			}
			pdf.Ln(4)
		}

		if len(d.DomainCompliance) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "Domain Compliance", "", 1, "L", false, 0, "")
			pdf.Ln(2)
			pdf.SetFont("Helvetica", "", 8)
			for _, f := range d.DomainCompliance {
				pdf.CellFormat(190, 4, "  - "+f, "", 1, "L", false, 0, "")
			}
			pdf.Ln(4)
		}
	}

	// ── ERN — Executive Risk Narratives (PDF) ──
	if len(result.ERN.ExecutiveRisks) > 0 {
		pdf.AddPage()
		pdf.SetFont("Helvetica", "B", 16)
		pdf.CellFormat(190, 12, "Executive Risk Narratives", "", 1, "L", false, 0, "")
		pdf.Ln(4)
		pdf.SetFont("Helvetica", "", 10)
		pdf.MultiCell(190, 5, fmt.Sprintf("Financial Exposure: %s — %s", result.ERN.FinancialExposure.Level, result.ERN.FinancialExposure.Rationale), "", "L", false)
		pdf.Ln(4)
		if result.ERN.BoardSummary.Summary != "" {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "Board Summary", "", 1, "L", false, 0, "")
			pdf.SetFont("Helvetica", "", 9)
			pdf.MultiCell(190, 4, result.ERN.BoardSummary.Summary, "", "L", false)
			pdf.Ln(4)
		}
		if len(result.ERN.CISOBriefing.TopRisks) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, fmt.Sprintf("Top Risks (%d)", len(result.ERN.CISOBriefing.TopRisks)), "", 1, "L", false, 0, "")
			pdf.SetFont("Helvetica", "", 8)
			for _, r := range result.ERN.CISOBriefing.TopRisks {
				pdf.MultiCell(190, 4, "  - "+r, "", "L", false)
			}
			pdf.Ln(4)
		}
		if len(result.ERN.RemediationRoadmap.Phase30) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "Remediation Roadmap - 30 Days", "", 1, "L", false, 0, "")
			pdf.SetFont("Helvetica", "", 8)
			for _, item := range result.ERN.RemediationRoadmap.Phase30 {
				pdf.MultiCell(190, 4, fmt.Sprintf("  [%s] %s", item.Priority, item.Action), "", "L", false)
			}
			pdf.Ln(4)
		}
		if len(result.ERN.InvestmentInsights) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "Security Investment Insights", "", 1, "L", false, 0, "")
			pdf.SetFont("Helvetica", "", 8)
			for _, ii := range result.ERN.InvestmentInsights {
				pdf.MultiCell(190, 4, fmt.Sprintf("  %s [%s]", ii.Area, ii.Priority), "", "L", false)
			}
			pdf.Ln(4)
		}
		if result.ERN.DecisionSupport.Top3Actions != nil && len(result.ERN.DecisionSupport.Top3Actions) > 0 {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "CISO Decision Support - Top 3 Actions", "", 1, "L", false, 0, "")
			pdf.SetFont("Helvetica", "", 8)
			for _, da := range result.ERN.DecisionSupport.Top3Actions {
				pdf.MultiCell(190, 4, fmt.Sprintf("  %d. %s (%s)", da.Rank, da.Action, da.Impact), "", "L", false)
			}
			pdf.Ln(4)
		}
	}

	// ── Report Packs (PDF) ──
	if result.ERN.ReportPacks.BoardReport != "" || result.ERN.ReportPacks.ExecutiveReport != "" || result.ERN.ReportPacks.TechnicalReport != "" {
		pdf.AddPage()
		pdf.SetFont("Helvetica", "B", 16)
		pdf.CellFormat(190, 12, "Report Packs", "", 1, "L", false, 0, "")
		pdf.Ln(4)
		if result.ERN.ReportPacks.BoardReport != "" {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "Board Report", "", 1, "L", false, 0, "")
			pdf.SetFont("Helvetica", "", 8)
			pdf.MultiCell(190, 4, stripMarkdownHeadings(result.ERN.ReportPacks.BoardReport), "", "L", false)
			pdf.Ln(6)
		}
		if result.ERN.ReportPacks.ExecutiveReport != "" {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "Executive Report", "", 1, "L", false, 0, "")
			pdf.SetFont("Helvetica", "", 8)
			pdf.MultiCell(190, 4, stripMarkdownHeadings(result.ERN.ReportPacks.ExecutiveReport), "", "L", false)
			pdf.Ln(6)
		}
		if result.ERN.ReportPacks.TechnicalReport != "" {
			pdf.SetFont("Helvetica", "B", 12)
			pdf.CellFormat(190, 7, "Technical Report", "", 1, "L", false, 0, "")
			pdf.SetFont("Helvetica", "", 8)
			pdf.MultiCell(190, 4, stripMarkdownHeadings(result.ERN.ReportPacks.TechnicalReport), "", "L", false)
			pdf.Ln(6)
		}
	}

	// ── Portfolio Intelligence (PDF) ──
	sampiMD := renderSAMPIReportMarkdown(result.SAMPI)
	if sampiMD != "" {
		pdf.AddPage()
		pdf.SetFont("Helvetica", "B", 16)
		pdf.CellFormat(190, 12, "Portfolio Intelligence", "", 1, "L", false, 0, "")
		pdf.Ln(4)
		pdf.SetFont("Helvetica", "", 8)
		pdf.MultiCell(190, 4, stripMarkdownHeadings(sampiMD), "", "L", false)
		pdf.Ln(6)
	}

	// ── Decision Intelligence (PDF) ──
	if sdiMD := renderSDIReportMarkdown(result.SDI); sdiMD != "" {
		pdf.AddPage()
		pdf.SetFont("Helvetica", "B", 16)
		pdf.CellFormat(190, 12, "Decision Intelligence", "", 1, "L", false, 0, "")
		pdf.Ln(4)
		pdf.SetFont("Helvetica", "", 8)
		pdf.MultiCell(190, 4, stripMarkdownHeadings(sdiMD), "", "L", false)
		pdf.Ln(6)
	}

	// ── Digital Twin (PDF) ──
	if sdtMD := renderSDTReportMarkdown(result.SDT); sdtMD != "" {
		pdf.AddPage()
		pdf.SetFont("Helvetica", "B", 16)
		pdf.CellFormat(190, 12, "Digital Twin", "", 1, "L", false, 0, "")
		pdf.Ln(4)
		pdf.SetFont("Helvetica", "", 8)
		pdf.MultiCell(190, 4, stripMarkdownHeadings(sdtMD), "", "L", false)
		pdf.Ln(6)
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

func stripMarkdownHeadings(md string) string {
	lines := strings.Split(md, "\n")
	var out []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "#") {
			text := strings.TrimLeft(trimmed, "# ")
			out = append(out, text)
			out = append(out, strings.Repeat("-", len(text)))
		} else if strings.HasPrefix(trimmed, "|") {
			if !strings.Contains(trimmed, "---") {
				out = append(out, strings.ReplaceAll(trimmed, "|", " | "))
			}
		} else if strings.HasPrefix(trimmed, "- ") {
			out = append(out, strings.TrimPrefix(trimmed, "- "))
		} else {
			out = append(out, line)
		}
	}
	return strings.Join(out, "\n")
}

func renderSAMPIReportMarkdown(sampi SAMPIIntelligence) string {
	var b strings.Builder
	if sampi.Dashboard.TotalArchitectures == 0 {
		return ""
	}

	b.WriteString("## Portfolio Intelligence\n\n")
	b.WriteString("### Dashboard\n\n")
	b.WriteString("| Metric | Value |\n|--------|-------|\n")
	b.WriteString(fmt.Sprintf("| Architectures | %d |\n", sampi.Dashboard.TotalArchitectures))
	b.WriteString(fmt.Sprintf("| Findings | %d |\n", sampi.Dashboard.TotalFindings))
	b.WriteString(fmt.Sprintf("| Threats | %d |\n", sampi.Dashboard.TotalThreats))
	b.WriteString(fmt.Sprintf("| Attack Paths | %d |\n", sampi.Dashboard.TotalAttackPaths))
	b.WriteString(fmt.Sprintf("| Controls | %d |\n", sampi.Dashboard.TotalControls))
	b.WriteString(fmt.Sprintf("| Average Risk Score | %.1f |\n", sampi.Dashboard.AverageRiskScore))
	b.WriteString(fmt.Sprintf("| Average Coverage | %.1f%% |\n", sampi.Dashboard.AverageCoverage))
	b.WriteString(fmt.Sprintf("| Compliance Count | %d |\n", sampi.Dashboard.ComplianceCount))
	b.WriteString("\n")

	if len(sampi.Heatmaps) > 0 {
		b.WriteString("### Executive Heatmap\n\n")
		b.WriteString("| Architecture | Risk Band | Score | Findings | Controls |\n")
		b.WriteString("|-------------|-----------|-------|----------|----------|\n")
		for _, h := range sampi.Heatmaps {
			b.WriteString(fmt.Sprintf("| %s | %s | %.1f | %d | %d |\n", h.ArchitectureName, h.RiskBand, h.RiskScore, h.FindingCount, h.ControlCount))
		}
		b.WriteString("\n")
	}

	if len(sampi.RepeatedWeaknesses) > 0 {
		b.WriteString("### Repeated Weaknesses\n\n")
		b.WriteString("| Finding | Category | Severity | Occurrences | Systemic |\n")
		b.WriteString("|---------|----------|----------|-------------|----------|\n")
		for _, rw := range sampi.RepeatedWeaknesses {
			sys := ""
			if rw.Systemic {
				sys = "⚠"
			}
			b.WriteString(fmt.Sprintf("| %s | %s | %s | %d | %s |\n", rw.FindingTitle, rw.Category, rw.Severity, rw.OccurrenceCount, sys))
			b.WriteString(fmt.Sprintf("  Affected: %s\n", strings.Join(rw.AffectedArchitectures, ", ")))
		}
		b.WriteString("\n")
	}

	if len(sampi.EnterpriseThemes) > 0 {
		b.WriteString("### Enterprise Risk Themes\n\n")
		b.WriteString("| Theme | Count | Architectures | Severity |\n")
		b.WriteString("|-------|-------|---------------|----------|\n")
		for _, th := range sampi.EnterpriseThemes {
			b.WriteString(fmt.Sprintf("| %s | %d | %d | %s |\n", th.Name, th.RiskCount, th.AffectedArchitectures, th.Severity))
		}
		b.WriteString("\n")
	}

	if len(sampi.ControlCoverage) > 0 {
		b.WriteString("### Control Coverage\n\n")
		b.WriteString("| Control | Coverage | Architectures |\n")
		b.WriteString("|---------|----------|---------------|\n")
		for _, cc := range sampi.ControlCoverage {
			b.WriteString(fmt.Sprintf("| %s | %.1f%% | %d/%d |\n", cc.ControlName, cc.CoveragePercent, cc.ArchitecturesWith, cc.ArchitecturesTotal))
		}
		b.WriteString("\n")
	}

	if sampi.SecurityDebt.Score > 0 {
		b.WriteString("### Security Debt\n\n")
		b.WriteString(fmt.Sprintf("- **Score:** %.1f\n", sampi.SecurityDebt.Score))
		b.WriteString(fmt.Sprintf("- **Longstanding Findings:** %d\n", sampi.SecurityDebt.LongstandingCount))
		b.WriteString(fmt.Sprintf("- **Repeated Findings:** %d\n", sampi.SecurityDebt.RepeatedCount))
		b.WriteString("\n")
		if len(sampi.SecurityDebt.TopDebts) > 0 {
			b.WriteString("| Description | Architecture | Category | Severity |\n")
			b.WriteString("|-------------|--------------|----------|----------|\n")
			for _, d := range sampi.SecurityDebt.TopDebts {
				b.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", d.Description, d.Architecture, d.Category, d.Severity))
			}
			b.WriteString("\n")
		}
	}

	if len(sampi.ProgramInsights) > 0 {
		b.WriteString("### Security Program Insights\n\n")
		b.WriteString("| Area | Priority | Insight |\n")
		b.WriteString("|------|----------|---------|\n")
		for _, pi := range sampi.ProgramInsights {
			b.WriteString(fmt.Sprintf("| %s | %s | %s |\n", pi.Area, pi.Priority, pi.Insight))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func renderSDIReportMarkdown(sdi SDIIntelligence) string {
	var b strings.Builder
	if len(sdi.Recommendations) == 0 {
		return ""
	}
	b.WriteString("## Decision Intelligence\n\n")
	b.WriteString("### Top Recommendations\n\n")
	b.WriteString("| ID | Action | Priority | Risk Reduction | Effort | ROI | Findings |\n")
	b.WriteString("|----|--------|----------|----------------|--------|-----|----------|\n")
	for _, r := range sdi.Recommendations {
		roi := r.RiskReduction + "/" + r.Effort
		b.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s | %s | %d |\n",
			r.ID, r.Title, r.Priority, r.RiskReduction, r.Effort, roi, len(r.AffectedFindings)))
	}
	b.WriteString("\n")

	if len(sdi.Dashboard.QuickWins) > 0 {
		b.WriteString("### Quick Wins\n\n")
		for _, qw := range sdi.Dashboard.QuickWins {
			b.WriteString(fmt.Sprintf("- %s (Low effort, %s priority)\n", qw.Title, qw.Priority))
		}
		b.WriteString("\n")
	}

	if len(sdi.FixSimulations) > 0 {
		b.WriteString("### Fix Simulations\n\n")
		b.WriteString("| Control | Critical Before | Critical After | High Before | High After | Coverage Before | Coverage After |\n")
		b.WriteString("|---------|----------------|----------------|-------------|------------|-----------------|----------------|\n")
		for _, sim := range sdi.FixSimulations {
			b.WriteString(fmt.Sprintf("| %s | %d | %d | %d | %d | %.0f%% | %.0f%% |\n",
				sim.ControlName, sim.OriginalCritical, sim.NewCritical,
				sim.OriginalHigh, sim.NewHigh, sim.OriginalCoverage, sim.NewCoverage))
		}
		b.WriteString("\n")
	}

	if len(sdi.FailureSimulations) > 0 {
		b.WriteString("### Failure Simulations\n\n")
		b.WriteString("| Control | Systems Impacted | New Attack Paths | Risk Increase |\n")
		b.WriteString("|---------|------------------|-----------------|---------------|\n")
		for _, sim := range sdi.FailureSimulations {
			b.WriteString(fmt.Sprintf("| %s | %d | %d | %s |\n",
				sim.ControlName, sim.SystemsImpacted, sim.AttackPathsOpened, sim.RiskIncrease))
		}
		b.WriteString("\n")
	}

	if len(sdi.InvestmentPriorities) > 0 {
		b.WriteString("### Investment Priorities\n\n")
		b.WriteString("| Rank | Area | Score | Findings | Risk Reduction |\n")
		b.WriteString("|------|------|-------|----------|----------------|\n")
		for _, ip := range sdi.InvestmentPriorities {
			b.WriteString(fmt.Sprintf("| %d | %s | %.1f | %d | %s |\n",
				ip.Rank, ip.Area, ip.Score, ip.FindingCount, ip.RiskReduction))
		}
		b.WriteString("\n")
	}

	if len(sdi.AttackPathCollapse) > 0 {
		b.WriteString("### Attack Path Collapse Analysis\n\n")
		b.WriteString("| Control | Attack Paths Reduced | Total | Reduction |\n")
		b.WriteString("|---------|---------------------|-------|-----------|\n")
		for _, apc := range sdi.AttackPathCollapse {
			b.WriteString(fmt.Sprintf("| %s | %d | %d | %.0f%% |\n",
				apc.ControlName, apc.AttackPathsReduced, apc.TotalAttackPaths, apc.ReductionPercent))
		}
		b.WriteString("\n")
	}

	if len(sdi.ComplianceImpacts) > 0 {
		b.WriteString("### Compliance Impact\n\n")
		b.WriteString("| Framework | Action | Improvement |\n")
		b.WriteString("|-----------|--------|-------------|\n")
		for _, ci := range sdi.ComplianceImpacts {
			b.WriteString(fmt.Sprintf("| %s | %s | %s |\n", ci.Framework, ci.Action, ci.Improvement))
		}
		b.WriteString("\n")
	}

	if len(sdi.RemediationRoadmap.Phase30) > 0 {
		b.WriteString("### Remediation Roadmap\n\n")
		b.WriteString("**30 Days:**\n\n")
		for _, item := range sdi.RemediationRoadmap.Phase30 {
			b.WriteString(fmt.Sprintf("- %s [%s]\n", item.Action, item.Priority))
		}
		b.WriteString("\n")
	}
	if len(sdi.RemediationRoadmap.Phase90) > 0 {
		b.WriteString("**90 Days:**\n\n")
		for _, item := range sdi.RemediationRoadmap.Phase90 {
			b.WriteString(fmt.Sprintf("- %s [%s]\n", item.Action, item.Priority))
		}
		b.WriteString("\n")
	}
	if len(sdi.RemediationRoadmap.Phase180) > 0 {
		b.WriteString("**180 Days:**\n\n")
		for _, item := range sdi.RemediationRoadmap.Phase180 {
			b.WriteString(fmt.Sprintf("- %s [%s]\n", item.Action, item.Priority))
		}
		b.WriteString("\n")
	}
	if len(sdi.RemediationRoadmap.Phase12m) > 0 {
		b.WriteString("**12 Months:**\n\n")
		for _, item := range sdi.RemediationRoadmap.Phase12m {
			b.WriteString(fmt.Sprintf("- %s [%s]\n", item.Action, item.Priority))
		}
		b.WriteString("\n")
	}

	if sdi.Dashboard.RiskReductionSummary != "" {
		b.WriteString(fmt.Sprintf("### Summary\n\n%s\n\n", sdi.Dashboard.RiskReductionSummary))
	}

	if sdi.ExecutiveScenarios.BestCase.Scenario != "" {
		b.WriteString("### Executive Scenarios\n\n")
		b.WriteString("| Scenario | Risk Score | Findings Resolved | Attack Paths Closed | Coverage |\n")
		b.WriteString("|----------|------------|-------------------|--------------------|----------|\n")
		b.WriteString(fmt.Sprintf("| %s | %.1f | %d | %d | %.0f%% |\n",
			sdi.ExecutiveScenarios.BestCase.Scenario,
			sdi.ExecutiveScenarios.BestCase.RiskScore,
			sdi.ExecutiveScenarios.BestCase.FindingsResolved,
			sdi.ExecutiveScenarios.BestCase.AttackPathsClosed,
			sdi.ExecutiveScenarios.BestCase.CoverageAchieved))
		b.WriteString(fmt.Sprintf("| %s | %.1f | %d | %d | %.0f%% |\n",
			sdi.ExecutiveScenarios.LikelyCase.Scenario,
			sdi.ExecutiveScenarios.LikelyCase.RiskScore,
			sdi.ExecutiveScenarios.LikelyCase.FindingsResolved,
			sdi.ExecutiveScenarios.LikelyCase.AttackPathsClosed,
			sdi.ExecutiveScenarios.LikelyCase.CoverageAchieved))
		b.WriteString(fmt.Sprintf("| %s | %.1f | %d | %d | %.0f%% |\n",
			sdi.ExecutiveScenarios.WorstCase.Scenario,
			sdi.ExecutiveScenarios.WorstCase.RiskScore,
			sdi.ExecutiveScenarios.WorstCase.FindingsResolved,
			sdi.ExecutiveScenarios.WorstCase.AttackPathsClosed,
			sdi.ExecutiveScenarios.WorstCase.CoverageAchieved))
		b.WriteString("\n")
	}

	return b.String()
}

func renderSDTReportMarkdown(sdt SDTIntelligence) string {
	var b strings.Builder
	if sdt.Twin.ID == "" {
		return ""
	}
	b.WriteString("## Digital Twin\n\n")
	b.WriteString(fmt.Sprintf("**Architecture:** %s (v%s)  \n", sdt.Twin.ArchitectureName, sdt.Twin.Version))
	b.WriteString(fmt.Sprintf("**Risk Score:** %.1f  |  **Coverage:** %.0f%%  |  **Source Hash:** %s\n\n", sdt.Twin.RiskScore, sdt.Twin.Coverage, sdt.Twin.SourceHash))

	if len(sdt.ChangeImpacts) > 0 {
		b.WriteString("### Change Impact Analysis\n\n")
		b.WriteString("| Change | Component | Impact Type | Severity | Risks | Paths |\n")
		b.WriteString("|--------|-----------|-------------|----------|-------|-------|\n")
		for _, ci := range sdt.ChangeImpacts {
			b.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %d | %d |\n",
				ci.Change, ci.ComponentAffected, ci.ImpactType, ci.Severity, ci.RisksAffected, ci.AttackPathsAffected))
		}
		b.WriteString("\n")
	}

	if len(sdt.ControlDrifts) > 0 {
		b.WriteString("### Control Drift\n\n")
		b.WriteString("| Control | Category | Expected | Current | Risk Impact |\n")
		b.WriteString("|---------|----------|----------|---------|-------------|\n")
		for _, cd := range sdt.ControlDrifts {
			b.WriteString(fmt.Sprintf("| %s | %s | %s | %s | %s |\n",
				cd.ControlName, cd.Category, cd.ExpectedState, cd.CurrentState, cd.RiskImpact))
		}
		b.WriteString("\n")
	}

	if sdt.SecurityDebt.TotalDebt > 0 {
		b.WriteString("### Security Debt\n\n")
		b.WriteString(fmt.Sprintf("**Total:** %.0f  |  **Findings:** %.0f  |  **Controls:** %.0f  |  **Assumptions:** %.0f  |  **Risk Score:** %.0f\n\n",
			sdt.SecurityDebt.TotalDebt, sdt.SecurityDebt.FindingDebt,
			sdt.SecurityDebt.ControlDebt, sdt.SecurityDebt.AssumptionDebt, sdt.SecurityDebt.RiskScore))
	}

	if len(sdt.ComplianceDrifts) > 0 {
		b.WriteString("### Compliance Drift\n\n")
		b.WriteString("| Framework | Status | New Gaps | Resolved | Regressed Areas |\n")
		b.WriteString("|-----------|--------|----------|----------|----------------|\n")
		for _, cd := range sdt.ComplianceDrifts {
			areas := strings.Join(cd.RegressedAreas, ", ")
			b.WriteString(fmt.Sprintf("| %s | %s | %d | %d | %s |\n",
				cd.Framework, cd.Status, cd.NewGaps, cd.ResolvedGaps, areas))
		}
		b.WriteString("\n")
	}

	if len(sdt.WhatIfScenarios) > 0 {
		b.WriteString("### What-If Scenarios\n\n")
		b.WriteString("| Scenario | Risk Delta | Coverage Delta | Findings Delta |\n")
		b.WriteString("|----------|------------|----------------|----------------|\n")
		for _, wi := range sdt.WhatIfScenarios {
			b.WriteString(fmt.Sprintf("| %s | %.1f | %.0f%% | %d |\n",
				wi.Name, wi.RiskDelta, wi.CoverageDelta, wi.FindingsDelta))
		}
		b.WriteString("\n")
	}

	if len(sdt.ZeroTrust.Dimensions) > 0 {
		b.WriteString("### Zero Trust Assessment\n\n")
		b.WriteString(fmt.Sprintf("**Overall:** %.1f / %.1f (gap: %.1f)\n\n", sdt.ZeroTrust.Overall, sdt.ZeroTrust.Target, sdt.ZeroTrust.Gap))
		b.WriteString("| Dimension | Score | Target | Gap | Status |\n")
		b.WriteString("|-----------|-------|--------|-----|--------|\n")
		for _, d := range sdt.ZeroTrust.Dimensions {
			b.WriteString(fmt.Sprintf("| %s | %.1f | %.1f | %.1f | %s |\n", d.Dimension, d.Score, d.Target, d.Gap, d.Status))
		}
		b.WriteString("\n")
	}

	if len(sdt.CrownJewels) > 0 {
		b.WriteString("### Crown Jewel Analysis\n\n")
		b.WriteString("| Asset | Business Value | Attack Value | Dependencies | Threats | Blast Radius | Score |\n")
		b.WriteString("|-------|---------------|--------------|--------------|---------|--------------|-------|\n")
		for _, cj := range sdt.CrownJewels {
			b.WriteString(fmt.Sprintf("| %s | %s | %s | %d | %d | %s | %.1f |\n",
				cj.AssetName, cj.BusinessValue, cj.AttackValue, cj.DependencyCount,
				cj.ThreatCount, cj.BlastRadius, cj.OverallScore))
		}
		b.WriteString("\n")
	}

	if sdt.ExecutiveReport.ArchitectureHealth != "" {
		b.WriteString("### Executive Report\n\n")
		b.WriteString(fmt.Sprintf("**Architecture Health:** %s  \n", sdt.ExecutiveReport.ArchitectureHealth))
		b.WriteString(fmt.Sprintf("**Security Debt Score:** %.0f  \n", sdt.ExecutiveReport.SecurityDebtScore))
		b.WriteString(fmt.Sprintf("**Control Drifts:** %d  |  **Compliance Drifts:** %d  \n", sdt.ExecutiveReport.ControlDriftCount, sdt.ExecutiveReport.ComplianceDriftCount))
		b.WriteString(fmt.Sprintf("**Risk Trend:** %s  |  **Attack Surface Trend:** %s\n\n", sdt.ExecutiveReport.RiskTrend, sdt.ExecutiveReport.AttackSurfaceTrend))
	}

	if sdt.PortfolioSummary.ArchitectureCount > 0 {
		b.WriteString("### Portfolio Summary\n\n")
		b.WriteString(fmt.Sprintf("**Architecture Count:** %d  |  **Aggregated Debt:** %.0f\n\n", sdt.PortfolioSummary.ArchitectureCount, sdt.PortfolioSummary.AggregatedDebt))
		if len(sdt.PortfolioSummary.EnterpriseTrends) > 0 {
			b.WriteString("**Enterprise Trends:**\n")
			for _, et := range sdt.PortfolioSummary.EnterpriseTrends {
				b.WriteString(fmt.Sprintf("- %s\n", et))
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}

func markdownToHTML(md string) string {
	lines := strings.Split(md, "\n")
	var html strings.Builder
	inTable := false
	inList := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(trimmed, "### "):
			if inTable {
				html.WriteString("</table>\n")
				inTable = false
			}
			if inList {
				html.WriteString("</ul>\n")
				inList = false
			}
			html.WriteString(fmt.Sprintf("<h4>%s</h4>\n", strings.TrimPrefix(trimmed, "### ")))
		case strings.HasPrefix(trimmed, "## "):
			if inTable {
				html.WriteString("</table>\n")
				inTable = false
			}
			if inList {
				html.WriteString("</ul>\n")
				inList = false
			}
			html.WriteString(fmt.Sprintf("<h3>%s</h3>\n", strings.TrimPrefix(trimmed, "## ")))
		case strings.HasPrefix(trimmed, "# "):
			if inTable {
				html.WriteString("</table>\n")
				inTable = false
			}
			if inList {
				html.WriteString("</ul>\n")
				inList = false
			}
			html.WriteString(fmt.Sprintf("<h2>%s</h2>\n", strings.TrimPrefix(trimmed, "# ")))
		case strings.HasPrefix(trimmed, "|"):
			if !inTable {
				html.WriteString("<table>\n")
				inTable = true
			}
			if strings.Contains(trimmed, "---") {
				continue
			}
			cells := strings.Split(trimmed, "|")
			html.WriteString("<tr>")
			for _, cell := range cells {
				cell = strings.TrimSpace(cell)
				if cell != "" {
					html.WriteString(fmt.Sprintf("<td>%s</td>", cell))
				}
			}
			html.WriteString("</tr>\n")
		case strings.HasPrefix(trimmed, "- "):
			if inTable {
				html.WriteString("</table>\n")
				inTable = false
			}
			if !inList {
				html.WriteString("<ul>\n")
				inList = true
			}
			item := strings.TrimPrefix(trimmed, "- ")
			item = strings.ReplaceAll(item, "**", "<b>")
			item = strings.ReplaceAll(item, "**", "</b>")
			html.WriteString(fmt.Sprintf("<li>%s</li>\n", item))
		default:
			if inTable {
				html.WriteString("</table>\n")
				inTable = false
			}
			if inList {
				html.WriteString("</ul>\n")
				inList = false
			}
			if trimmed != "" {
				html.WriteString(fmt.Sprintf("<p>%s</p>\n", trimmed))
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

func severityCSSClass(severity string) string {
	switch severity {
	case "Critical":
		return "critical"
	case "High":
		return "high"
	case "Medium":
		return "medium"
	case "Low":
		return "low"
	}
	return "medium"
}

func riskLevelForAPDScore(score float64) string {
	switch {
	case score >= 0.8:
		return "Critical"
	case score >= 0.6:
		return "High"
	case score >= 0.4:
		return "Medium"
	default:
		return "Low"
	}
}

type exportModel struct {
	selected         int
	format           ExportFormat
	done             bool
	exportPath       string
	showConfirmation bool
	result           *AnalysisResult
	outputDir        string
	err              error
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
			if m.selected < 13 {
				m.selected++
			}
		case "enter":
			formats := []ExportFormat{ExportJSON, ExportMarkdown, ExportCSV, ExportPDF, ExportHTML, ExportNarrativeMarkdown, ExportNarrativeHTML, ExportTrustMarkdown, ExportTrustHTML, ExportTrustJSON, ExportReviewMarkdown, ExportReviewHTML, ExportConfidenceMarkdown, ExportConfidenceHTML}
			if m.selected < len(formats) {
				m.format = formats[m.selected]
				m.showConfirmation = true
			}
		case "esc":
			m.showConfirmation = false
			m.done = false
			m.err = nil
		case "y":
			if m.showConfirmation && !m.done {
				if m.result != nil {
					path, err := ExportResult(m.result, m.format, m.outputDir)
					if err != nil {
						m.err = err
					} else {
						m.done = true
						m.exportPath = path
					}
				}
			}
		}
	}
	return m, nil
}

func (m mainModel) viewExport() string {
	s := m.styles
	ex := m.exportV

	if ex.err != nil {
		return lipgloss.JoinVertical(lipgloss.Left,
			s.Title.Render("Export Error"),
			s.BorderBox.Render(
				lipgloss.JoinVertical(lipgloss.Center,
					s.StatusBad.Render("✗ Export Failed"),
					s.SectionItem.Render(ex.err.Error()),
					s.SectionItem.Render("Press Esc to return."),
				),
			),
		)
	}

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
		{"Narrative Markdown (.md)", ExportNarrativeMarkdown},
		{"Narrative HTML (.html)", ExportNarrativeHTML},
		{"Trust Chain Markdown (.md)", ExportTrustMarkdown},
		{"Trust Chain HTML (.html)", ExportTrustHTML},
		{"Trust Chain JSON (.json)", ExportTrustJSON},
		{"Review Workbench Markdown (.md)", ExportReviewMarkdown},
		{"Review Workbench HTML (.html)", ExportReviewHTML},
		{"Confidence Report Markdown (.md)", ExportConfidenceMarkdown},
		{"Confidence Report HTML (.html)", ExportConfidenceHTML},
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
