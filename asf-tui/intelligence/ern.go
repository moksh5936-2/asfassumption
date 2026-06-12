package intelligence

import (
	"fmt"
	"sort"
	"strings"
)

// ─────────────────────────────────────────────────────────────
// ERN — Executive Risk Narratives & CISO Reporting Engine (ASF V9)
// Phases 1-15
// ─────────────────────────────────────────────────────────────

// ── PHASE 1 — EXECUTIVE RISK MODEL ──

// ExecutiveRisk represents a single executive-level risk narrative.
type ExecutiveRisk struct {
	ID                 string   `json:"id"`
	Title              string   `json:"title"`
	Summary            string   `json:"summary"`
	BusinessImpact     string   `json:"business_impact"`
	OperationalImpact  string   `json:"operational_impact"`
	ComplianceImpact   string   `json:"compliance_impact"`
	FinancialImpact    string   `json:"financial_impact"`
	ReputationImpact   string   `json:"reputation_impact"`
	Likelihood         string   `json:"likelihood"`
	Severity           string   `json:"severity"`
	Priority           string   `json:"priority"`
	AffectedAssets     []string `json:"affected_assets"`
	AffectedControls   []string `json:"affected_controls"`
	RecommendedActions []string `json:"recommended_actions"`
}

// ── PHASE 2 — RISK NARRATIVE ──

// RiskNarrative holds a technical-to-business narrative conversion.
type RiskNarrative struct {
	FindingID        string `json:"finding_id"`
	TechnicalSummary string `json:"technical_summary"`
	Narrative        string `json:"narrative"`
}

// ── PHASE 3 — BUSINESS IMPACT MAPPING ──

type BusinessImpactMap struct {
	Categories []BusinessImpactCategory `json:"categories"`
}

type BusinessImpactCategory struct {
	Name     string   `json:"name"`
	Score    int      `json:"score"`
	Findings []string `json:"findings"`
}

// ── PHASE 4 — CROWN JEWEL BUSINESS CLASSIFICATION ──

type CrownJewelClass struct {
	TechnicalName    string `json:"technical_name"`
	BusinessCategory string `json:"business_category"`
	BusinessLabel    string `json:"business_label"`
}

// ── PHASE 5 — FINANCIAL EXPOSURE ──

type FinancialExposure struct {
	Level     string `json:"level"`
	Rationale string `json:"rationale"`
}

// ── PHASE 6 — REGULATORY IMPACT ──

type RegulatoryImpact struct {
	Framework string `json:"framework"`
	Domain    string `json:"domain"`
	Exposure  string `json:"exposure"`
	Rationale string `json:"rationale"`
}

// ── PHASE 7 — EXECUTIVE PRIORITY ──

type PriorityRisk struct {
	Risk     ExecutiveRisk `json:"risk"`
	Priority string        `json:"priority"`
	Score    int           `json:"score"`
}

// ── PHASE 8 — RISK AGGREGATION ──

type RiskTheme struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	RiskCount   int      `json:"risk_count"`
	Severity    string   `json:"severity"`
	Findings    []string `json:"findings"`
}

// ── PHASE 9 — BOARD SUMMARY ──

type BoardSummary struct {
	Summary string `json:"summary"`
}

// ── PHASE 10 — CISO BRIEFING ──

type CISOBriefing struct {
	TopRisks           []string         `json:"top_risks"`
	TopRemediations    []string         `json:"top_remediations"`
	HighRiskAssets     []string         `json:"high_risk_assets"`
	CoverageOverview   CoverageOverview `json:"coverage_overview"`
	ComplianceOverview string           `json:"compliance_overview"`
}

type CoverageOverview struct {
	TotalControls int     `json:"total_controls"`
	Covered       int     `json:"covered"`
	Partial       int     `json:"partial"`
	Missing       int     `json:"missing"`
	CoverageRate  float64 `json:"coverage_rate"`
}

// ── PHASE 11 — REMEDIATION ROADMAP ──

type RemediationRoadmap struct {
	Phase30  []RemediationItem `json:"phase_30"`
	Phase90  []RemediationItem `json:"phase_90"`
	Phase180 []RemediationItem `json:"phase_180"`
	Phase12m []RemediationItem `json:"phase_12m"`
}

type RemediationItem struct {
	Action   string `json:"action"`
	Category string `json:"category"`
	Priority string `json:"priority"`
}

// ── PHASE 12 — RISK TREND MODEL ──

type RiskTrend struct {
	CurrentState RiskTrendState `json:"current_state"`
	TargetState  RiskTrendState `json:"target_state"`
}

type RiskTrendState struct {
	RiskScore           float64 `json:"risk_score"`
	CoverageRate        float64 `json:"coverage_rate"`
	ComplianceReadiness float64 `json:"compliance_readiness"`
	CriticalFindings    int     `json:"critical_findings"`
}

// ── PHASE 13 — SECURITY INVESTMENT INSIGHTS ──

type InvestmentInsight struct {
	Area        string `json:"area"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Rationale   string `json:"rationale"`
}

// ── PHASE 14 — EXECUTIVE DASHBOARD ──

type ExecutiveDashboard struct {
	RiskScore           float64  `json:"risk_score"`
	PriorityFindings    int      `json:"priority_findings"`
	ComplianceReadiness float64  `json:"compliance_readiness"`
	CoverageRate        float64  `json:"coverage_rate"`
	AttackPathCount     int      `json:"attack_path_count"`
	CriticalAssets      []string `json:"critical_assets"`
}

// ── PHASE 15 — CISO DECISION SUPPORT ──

type DecisionSupport struct {
	Top3Actions []DecisionAction `json:"top_3_actions"`
}

type DecisionAction struct {
	Rank      int    `json:"rank"`
	Action    string `json:"action"`
	Impact    string `json:"impact"`
	Rationale string `json:"rationale"`
}

// ── PHASE 16 — REPORT PACKS ──

type ReportPackType int

const (
	ReportPackBoard ReportPackType = iota
	ReportPackExecutive
	ReportPackTechnical
)

type ReportPacks struct {
	BoardReport     string `json:"board_report"`
	ExecutiveReport string `json:"executive_report"`
	TechnicalReport string `json:"technical_report"`
}

func (r ReportPackType) String() string {
	switch r {
	case ReportPackBoard:
		return "Board"
	case ReportPackExecutive:
		return "Executive"
	case ReportPackTechnical:
		return "Technical"
	default:
		return "Unknown"
	}
}

func ParseReportPackType(s string) ReportPackType {
	switch strings.ToLower(s) {
	case "board":
		return ReportPackBoard
	case "executive":
		return ReportPackExecutive
	case "technical":
		return ReportPackTechnical
	default:
		return ReportPackTechnical
	}
}

func generateReportPacks(result *ERNRunResult) ReportPacks {
	return ReportPacks{
		BoardReport:     generateBoardReport(result),
		ExecutiveReport: generateExecutiveReport(result),
		TechnicalReport: generateTechnicalReport(result),
	}
}

func generateBoardReport(result *ERNRunResult) string {
	var b strings.Builder

	b.WriteString("# Board Report\n\n")
	b.WriteString("## Executive Summary\n\n")
	if result.BoardSummary.Summary != "" {
		b.WriteString(result.BoardSummary.Summary + "\n\n")
	}

	b.WriteString(fmt.Sprintf("## Financial Exposure\n\n**Level:** %s\n\n**Rationale:** %s\n\n",
		result.FinancialExposure.Level, result.FinancialExposure.Rationale))

	if len(result.RiskThemes) > 0 {
		b.WriteString("## Key Risk Themes\n\n")
		for _, th := range result.RiskThemes {
			b.WriteString(fmt.Sprintf("- **%s** (%d findings, %s severity)\n", th.Name, th.RiskCount, th.Severity))
		}
		b.WriteString("\n")
	}

	if result.DecisionSupport.Top3Actions != nil && len(result.DecisionSupport.Top3Actions) > 0 {
		b.WriteString("## Recommended Actions\n\n")
		for _, da := range result.DecisionSupport.Top3Actions {
			b.WriteString(fmt.Sprintf("%d. **%s** (%s impact)\n", da.Rank, da.Action, da.Impact))
		}
		b.WriteString("\n")
	}

	if result.Dashboard.RiskScore > 0 {
		b.WriteString("## Risk Overview\n\n")
		b.WriteString(fmt.Sprintf("- **Risk Score:** %.1f\n", result.Dashboard.RiskScore))
		b.WriteString(fmt.Sprintf("- **Priority Findings:** %d\n", result.Dashboard.PriorityFindings))
		b.WriteString(fmt.Sprintf("- **Coverage Rate:** %.1f%%\n", result.Dashboard.CoverageRate))
		b.WriteString(fmt.Sprintf("- **Attack Paths:** %d\n", result.Dashboard.AttackPathCount))
		if len(result.Dashboard.CriticalAssets) > 0 {
			b.WriteString(fmt.Sprintf("- **Critical Assets:** %s\n", strings.Join(result.Dashboard.CriticalAssets, ", ")))
		}
		b.WriteString("\n")
	}

	if len(result.RegulatoryImpacts) > 0 {
		b.WriteString("## Regulatory Considerations\n\n")
		for _, ri := range result.RegulatoryImpacts {
			b.WriteString(fmt.Sprintf("- **%s** (%s domain) — %s\n", ri.Framework, ri.Domain, ri.Exposure))
		}
		b.WriteString("\n")
	}

	return b.String()
}

func generateExecutiveReport(result *ERNRunResult) string {
	var b strings.Builder

	b.WriteString("# Executive Risk Report\n\n")

	if result.BoardSummary.Summary != "" {
		b.WriteString("## Board Summary\n\n")
		b.WriteString(result.BoardSummary.Summary + "\n\n")
	}

	b.WriteString(fmt.Sprintf("## Financial Exposure\n\n**Level:** %s — %s\n\n",
		result.FinancialExposure.Level, result.FinancialExposure.Rationale))

	if len(result.RiskThemes) > 0 {
		b.WriteString("## Risk Themes\n\n")
		b.WriteString("| Theme | Findings | Severity |\n")
		b.WriteString("|-------|----------|----------|\n")
		for _, th := range result.RiskThemes {
			b.WriteString(fmt.Sprintf("| %s | %d | %s |\n", th.Name, th.RiskCount, th.Severity))
		}
		b.WriteString("\n")
	}

	if len(result.ExecutiveRisks) > 0 {
		b.WriteString("## Executive Risks\n\n")
		for i, risk := range result.ExecutiveRisks {
			if i >= 5 {
				b.WriteString(fmt.Sprintf("... and %d more risks\n\n", len(result.ExecutiveRisks)-5))
				break
			}
			b.WriteString(fmt.Sprintf("### %s [%s] — %s\n\n", risk.ID, risk.Priority, risk.Title))
			b.WriteString(fmt.Sprintf("- **Business Impact:** %s\n", risk.BusinessImpact))
			b.WriteString(fmt.Sprintf("- **Operational Impact:** %s\n", risk.OperationalImpact))
			b.WriteString(fmt.Sprintf("- **Compliance Impact:** %s\n", risk.ComplianceImpact))
			if len(risk.RecommendedActions) > 0 {
				b.WriteString("- **Actions:**\n")
				for _, a := range risk.RecommendedActions {
					b.WriteString(fmt.Sprintf("  - %s\n", a))
				}
			}
			b.WriteString("\n")
		}
	}

	if len(result.CISOBriefing.TopRisks) > 0 {
		b.WriteString("## CISO Briefing — Top Risks\n\n")
		for _, r := range result.CISOBriefing.TopRisks {
			b.WriteString(fmt.Sprintf("- %s\n", r))
		}
		b.WriteString("\n")
	}

	if len(result.CISOBriefing.TopRemediations) > 0 {
		b.WriteString("## Top Remediations\n\n")
		for _, r := range result.CISOBriefing.TopRemediations {
			b.WriteString(fmt.Sprintf("- %s\n", r))
		}
		b.WriteString("\n")
	}

	b.WriteString("## Coverage Overview\n\n")
	b.WriteString(fmt.Sprintf("- **Controls:** %d total, %d covered, %d partial, %d missing\n",
		result.CISOBriefing.CoverageOverview.TotalControls,
		result.CISOBriefing.CoverageOverview.Covered,
		result.CISOBriefing.CoverageOverview.Partial,
		result.CISOBriefing.CoverageOverview.Missing))
	b.WriteString(fmt.Sprintf("- **Coverage Rate:** %.1f%%\n\n", result.CISOBriefing.CoverageOverview.CoverageRate))

	if len(result.RemediationRoadmap.Phase30) > 0 {
		b.WriteString("## Remediation — 30 Days\n\n")
		for _, item := range result.RemediationRoadmap.Phase30 {
			b.WriteString(fmt.Sprintf("- [%s] %s\n", item.Priority, item.Action))
		}
		b.WriteString("\n")
	}
	if len(result.RemediationRoadmap.Phase90) > 0 {
		b.WriteString("## Remediation — 90 Days\n\n")
		for _, item := range result.RemediationRoadmap.Phase90 {
			b.WriteString(fmt.Sprintf("- [%s] %s\n", item.Priority, item.Action))
		}
		b.WriteString("\n")
	}

	if len(result.InvestmentInsights) > 0 {
		b.WriteString("## Security Investment Insights\n\n")
		for _, ii := range result.InvestmentInsights {
			b.WriteString(fmt.Sprintf("- **%s** [%s]: %s\n", ii.Area, ii.Priority, ii.Rationale))
		}
		b.WriteString("\n")
	}

	if result.DecisionSupport.Top3Actions != nil && len(result.DecisionSupport.Top3Actions) > 0 {
		b.WriteString("## Decision Support — Top 3 Actions\n\n")
		for _, da := range result.DecisionSupport.Top3Actions {
			b.WriteString(fmt.Sprintf("%d. **%s** (%s impact) — %s\n", da.Rank, da.Action, da.Impact, da.Rationale))
		}
		b.WriteString("\n")
	}

	if len(result.CrownJewelClasses) > 0 {
		b.WriteString("## Crown Jewel Classification\n\n")
		for _, cj := range result.CrownJewelClasses {
			b.WriteString(fmt.Sprintf("- **%s** → *%s* (%s)\n", cj.TechnicalName, cj.BusinessLabel, cj.BusinessCategory))
		}
		b.WriteString("\n")
	}

	if len(result.RegulatoryImpacts) > 0 {
		b.WriteString("## Regulatory Impact\n\n")
		for _, ri := range result.RegulatoryImpacts {
			b.WriteString(fmt.Sprintf("- **%s** (%s): %s\n", ri.Framework, ri.Domain, ri.Rationale))
		}
		b.WriteString("\n")
	}

	b.WriteString("## Dashboard\n\n")
	b.WriteString(fmt.Sprintf("- **Risk Score:** %.1f\n", result.Dashboard.RiskScore))
	b.WriteString(fmt.Sprintf("- **Priority Findings:** %d\n", result.Dashboard.PriorityFindings))
	b.WriteString(fmt.Sprintf("- **Compliance Readiness:** %.1f%%\n", result.Dashboard.ComplianceReadiness))
	b.WriteString(fmt.Sprintf("- **Coverage Rate:** %.1f%%\n", result.Dashboard.CoverageRate))
	b.WriteString(fmt.Sprintf("- **Attack Paths:** %d\n", result.Dashboard.AttackPathCount))

	return b.String()
}

func generateTechnicalReport(result *ERNRunResult) string {
	var b strings.Builder

	b.WriteString("# Technical Security Report\n\n")

	b.WriteString("## Executive Summary\n\n")
	if result.BoardSummary.Summary != "" {
		b.WriteString(result.BoardSummary.Summary + "\n\n")
	}

	b.WriteString(fmt.Sprintf("## Financial Exposure\n\n**Level:** %s\n\n**Rationale:** %s\n\n",
		result.FinancialExposure.Level, result.FinancialExposure.Rationale))

	if len(result.RiskThemes) > 0 {
		b.WriteString("## Risk Themes\n\n")
		b.WriteString("| Theme | Description | Findings | Severity |\n")
		b.WriteString("|-------|-------------|----------|----------|\n")
		for _, th := range result.RiskThemes {
			b.WriteString(fmt.Sprintf("| %s | %s | %d | %s |\n", th.Name, th.Description, th.RiskCount, th.Severity))
		}
		b.WriteString("\n")
	}

	if len(result.ExecutiveRisks) > 0 {
		b.WriteString("## Executive Risks\n\n")
		for _, risk := range result.ExecutiveRisks {
			b.WriteString(fmt.Sprintf("### %s [%s] — %s\n\n", risk.ID, risk.Priority, risk.Title))
			b.WriteString(fmt.Sprintf("- **Business Impact:** %s\n", risk.BusinessImpact))
			b.WriteString(fmt.Sprintf("- **Operational Impact:** %s\n", risk.OperationalImpact))
			b.WriteString(fmt.Sprintf("- **Compliance Impact:** %s\n", risk.ComplianceImpact))
			b.WriteString(fmt.Sprintf("- **Financial Impact:** %s\n", risk.FinancialImpact))
			b.WriteString(fmt.Sprintf("- **Reputation Impact:** %s\n", risk.ReputationImpact))
			b.WriteString(fmt.Sprintf("- **Likelihood:** %s | **Severity:** %s\n", risk.Likelihood, risk.Severity))
			if len(risk.AffectedAssets) > 0 {
				b.WriteString(fmt.Sprintf("- **Affected Assets:** %s\n", strings.Join(risk.AffectedAssets, ", ")))
			}
			if len(risk.AffectedControls) > 0 {
				b.WriteString(fmt.Sprintf("- **Affected Controls:** %s\n", strings.Join(risk.AffectedControls, ", ")))
			}
			if len(risk.RecommendedActions) > 0 {
				b.WriteString("- **Recommended Actions:**\n")
				for _, a := range risk.RecommendedActions {
					b.WriteString(fmt.Sprintf("  - %s\n", a))
				}
			}
			b.WriteString("\n")
		}
	}

	b.WriteString("## CISO Briefing\n\n")
	if len(result.CISOBriefing.TopRisks) > 0 {
		b.WriteString("### Top Risks\n\n")
		for _, r := range result.CISOBriefing.TopRisks {
			b.WriteString(fmt.Sprintf("- %s\n", r))
		}
		b.WriteString("\n")
	}
	if len(result.CISOBriefing.TopRemediations) > 0 {
		b.WriteString("### Top Remediations\n\n")
		for _, r := range result.CISOBriefing.TopRemediations {
			b.WriteString(fmt.Sprintf("- %s\n", r))
		}
		b.WriteString("\n")
	}
	if len(result.CISOBriefing.HighRiskAssets) > 0 {
		b.WriteString("### Highest Risk Assets\n\n")
		for _, a := range result.CISOBriefing.HighRiskAssets {
			b.WriteString(fmt.Sprintf("- %s\n", a))
		}
		b.WriteString("\n")
	}
	b.WriteString("### Coverage Overview\n\n")
	b.WriteString(fmt.Sprintf("| Metric | Value |\n|--------|-------|\n"))
	b.WriteString(fmt.Sprintf("| Total Controls | %d |\n", result.CISOBriefing.CoverageOverview.TotalControls))
	b.WriteString(fmt.Sprintf("| Covered | %d |\n", result.CISOBriefing.CoverageOverview.Covered))
	b.WriteString(fmt.Sprintf("| Partial | %d |\n", result.CISOBriefing.CoverageOverview.Partial))
	b.WriteString(fmt.Sprintf("| Missing | %d |\n", result.CISOBriefing.CoverageOverview.Missing))
	b.WriteString(fmt.Sprintf("| Coverage Rate | %.1f%% |\n", result.CISOBriefing.CoverageOverview.CoverageRate))
	b.WriteString("\n")
	if result.CISOBriefing.ComplianceOverview != "" {
		b.WriteString(fmt.Sprintf("### Compliance Overview\n\n%s\n\n", result.CISOBriefing.ComplianceOverview))
	}

	b.WriteString("## Remediation Roadmap\n\n")
	phases := []struct {
		title string
		items []RemediationItem
	}{
		{"30 Days", result.RemediationRoadmap.Phase30},
		{"90 Days", result.RemediationRoadmap.Phase90},
		{"180 Days", result.RemediationRoadmap.Phase180},
		{"12 Months", result.RemediationRoadmap.Phase12m},
	}
	for _, phase := range phases {
		if len(phase.items) > 0 {
			b.WriteString(fmt.Sprintf("### %s\n\n", phase.title))
			for _, item := range phase.items {
				b.WriteString(fmt.Sprintf("- [%s] [%s] %s\n", item.Priority, item.Category, item.Action))
			}
			b.WriteString("\n")
		}
	}

	if len(result.InvestmentInsights) > 0 {
		b.WriteString("## Security Investment Insights\n\n")
		b.WriteString("| Area | Priority | Rationale |\n|------|----------|-----------|\n")
		for _, ii := range result.InvestmentInsights {
			b.WriteString(fmt.Sprintf("| %s | %s | %s |\n", ii.Area, ii.Priority, ii.Rationale))
		}
		b.WriteString("\n")
	}

	if result.DecisionSupport.Top3Actions != nil && len(result.DecisionSupport.Top3Actions) > 0 {
		b.WriteString("## Decision Support — Top 3 Actions\n\n")
		for _, da := range result.DecisionSupport.Top3Actions {
			b.WriteString(fmt.Sprintf("### Action %d: %s\n\n- **Impact:** %s\n- **Rationale:** %s\n\n",
				da.Rank, da.Action, da.Impact, da.Rationale))
		}
	}

	b.WriteString("## Risk Trend Analysis\n\n")
	b.WriteString("### Current State\n\n")
	b.WriteString(fmt.Sprintf("| Metric | Value |\n|--------|-------|\n"))
	b.WriteString(fmt.Sprintf("| Risk Score | %.1f |\n", result.RiskTrend.CurrentState.RiskScore))
	b.WriteString(fmt.Sprintf("| Coverage Rate | %.1f%% |\n", result.RiskTrend.CurrentState.CoverageRate))
	b.WriteString(fmt.Sprintf("| Compliance Readiness | %.1f%% |\n", result.RiskTrend.CurrentState.ComplianceReadiness))
	b.WriteString(fmt.Sprintf("| Critical Findings | %d |\n", result.RiskTrend.CurrentState.CriticalFindings))
	b.WriteString("\n### Target State\n\n")
	b.WriteString(fmt.Sprintf("| Metric | Value |\n|--------|-------|\n"))
	b.WriteString(fmt.Sprintf("| Risk Score | %.1f |\n", result.RiskTrend.TargetState.RiskScore))
	b.WriteString(fmt.Sprintf("| Coverage Rate | %.1f%% |\n", result.RiskTrend.TargetState.CoverageRate))
	b.WriteString(fmt.Sprintf("| Compliance Readiness | %.1f%% |\n", result.RiskTrend.TargetState.ComplianceReadiness))
	b.WriteString(fmt.Sprintf("| Critical Findings | %d |\n", result.RiskTrend.TargetState.CriticalFindings))
	b.WriteString("\n")

	if len(result.CrownJewelClasses) > 0 {
		b.WriteString("## Crown Jewel Classification\n\n")
		b.WriteString("| Technical Asset | Business Label | Category |\n|----------------|---------------|----------|\n")
		for _, cj := range result.CrownJewelClasses {
			b.WriteString(fmt.Sprintf("| %s | %s | %s |\n", cj.TechnicalName, cj.BusinessLabel, cj.BusinessCategory))
		}
		b.WriteString("\n")
	}

	if len(result.RegulatoryImpacts) > 0 {
		b.WriteString("## Regulatory Impact Analysis\n\n")
		b.WriteString("| Framework | Domain | Exposure | Rationale |\n|-----------|--------|----------|-----------|\n")
		for _, ri := range result.RegulatoryImpacts {
			b.WriteString(fmt.Sprintf("| %s | %s | %s | %s |\n", ri.Framework, ri.Domain, ri.Exposure, ri.Rationale))
		}
		b.WriteString("\n")
	}

	b.WriteString("## Dashboard\n\n")
	b.WriteString("| Metric | Value |\n|--------|-------|\n")
	b.WriteString(fmt.Sprintf("| Risk Score | %.1f |\n", result.Dashboard.RiskScore))
	b.WriteString(fmt.Sprintf("| Priority Findings | %d |\n", result.Dashboard.PriorityFindings))
	b.WriteString(fmt.Sprintf("| Compliance Readiness | %.1f%% |\n", result.Dashboard.ComplianceReadiness))
	b.WriteString(fmt.Sprintf("| Coverage Rate | %.1f%% |\n", result.Dashboard.CoverageRate))
	b.WriteString(fmt.Sprintf("| Attack Path Count | %d |\n", result.Dashboard.AttackPathCount))
	if len(result.Dashboard.CriticalAssets) > 0 {
		b.WriteString(fmt.Sprintf("| Critical Assets | %s |\n", strings.Join(result.Dashboard.CriticalAssets, ", ")))
	}
	b.WriteString("\n")

	if len(result.RiskNarratives) > 0 {
		b.WriteString("## Risk Narratives\n\n")
		for _, n := range result.RiskNarratives {
			b.WriteString(fmt.Sprintf("### Finding: %s\n\n%s\n\n", n.FindingID, n.Narrative))
		}
	}

	return b.String()
}

// ── ERN INPUT ──

// ERNInput consolidates all prior engine outputs for executive processing.
type ERNInput struct {
	Domain               string
	DomainPack           *KnowledgePack
	Threats              []Threat
	AttackPaths          []AttackPath
	Controls             []SDRIControl
	Findings             []SDRIFinding
	Assumptions          []Assumption
	ComplianceFrameworks []string
	Architecture         *ArchDescription
}

// ── ERN RUN RESULT ──

type ERNRunResult struct {
	ExecutiveRisks     []ExecutiveRisk     `json:"executive_risks"`
	RiskNarratives     []RiskNarrative     `json:"risk_narratives"`
	BusinessImpactMap  BusinessImpactMap   `json:"business_impact_map"`
	CrownJewelClasses  []CrownJewelClass   `json:"crown_jewel_classes"`
	FinancialExposure  FinancialExposure   `json:"financial_exposure"`
	RegulatoryImpacts  []RegulatoryImpact  `json:"regulatory_impacts"`
	PriorityRisks      []PriorityRisk      `json:"priority_risks"`
	RiskThemes         []RiskTheme         `json:"risk_themes"`
	BoardSummary       BoardSummary        `json:"board_summary"`
	CISOBriefing       CISOBriefing        `json:"ciso_briefing"`
	RemediationRoadmap RemediationRoadmap  `json:"remediation_roadmap"`
	RiskTrend          RiskTrend           `json:"risk_trend"`
	InvestmentInsights []InvestmentInsight `json:"investment_insights"`
	Dashboard          ExecutiveDashboard  `json:"dashboard"`
	DecisionSupport    DecisionSupport     `json:"decision_support"`
	ReportPacks        ReportPacks         `json:"report_packs,omitempty"`
}

// ── ERN ENGINE ──

type ERNEngine struct{}

func NewERNEngine() *ERNEngine {
	return &ERNEngine{}
}

func (e *ERNEngine) Run(input ERNInput) *ERNRunResult {
	result := &ERNRunResult{}

	// Phase 2-3: Generate risk narratives and business impact from findings
	result.RiskNarratives = generateRiskNarratives(input.Findings)
	result.BusinessImpactMap = mapBusinessImpact(input.Findings, input.Threats)

	// Phase 4: Classify crown jewels
	result.CrownJewelClasses = classifyCrownJewels(input.DomainPack)

	// Phase 1+2+3+7: Build executive risks from all inputs
	result.ExecutiveRisks = buildExecutiveRisks(input)

	// Phase 5: Estimate financial exposure
	result.FinancialExposure = estimateFinancialExposure(input)

	// Phase 6: Regulatory impact analysis
	result.RegulatoryImpacts = analyzeRegulatoryImpact(input)

	// Phase 7: Priority ranking
	result.PriorityRisks = rankExecutivePriorities(result.ExecutiveRisks)

	// Phase 8: Risk aggregation into themes
	result.RiskThemes = aggregateRiskThemes(result.ExecutiveRisks, input.Findings)

	// Phase 9: Board summary
	result.BoardSummary = generateBoardSummary(result)

	// Phase 10: CISO briefing
	result.CISOBriefing = generateCISOBriefing(result, input)

	// Phase 11: Remediation roadmap
	result.RemediationRoadmap = generateRemediationRoadmap(input)

	// Phase 12: Risk trend model
	result.RiskTrend = generateRiskTrend(result, input)

	// Phase 13: Security investment insights
	result.InvestmentInsights = generateInvestmentInsights(result, input)

	// Phase 14: Executive dashboard
	result.Dashboard = buildExecutiveDashboard(result, input)

	// Phase 15: CISO decision support
	result.DecisionSupport = generateDecisionSupport(result)

	// Phase 16: Generate report packs
	result.ReportPacks = generateReportPacks(result)

	return result
}

// ── PHASE 2 — RISK NARRATIVE GENERATION ──

func generateRiskNarratives(findings []SDRIFinding) []RiskNarrative {
	narratives := make([]RiskNarrative, 0)
	for _, f := range findings {
		narratives = append(narratives, RiskNarrative{
			FindingID:        f.ID,
			TechnicalSummary: f.Title + ": " + f.Description,
			Narrative:        buildNarrativeFromFinding(f),
		})
	}
	return narratives
}

func buildNarrativeFromFinding(f SDRIFinding) string {
	cat := strings.ToLower(f.Category)
	title := strings.ToLower(f.Title)

	switch {
	case containsAny(title, []string{"mfa", "authentication", "password", "credential", "identity", "sso", "login"}):
		return fmt.Sprintf(
			"%s may be vulnerable to unauthorized access. If compromised, attackers could gain elevated privileges, modify security settings, and access sensitive business systems. This could lead to operational disruption, compliance violations, and increased breach risk.",
			f.Title,
		)
	case containsAny(title, []string{"encryption", "key", "crypto", "tls", "ssl", "cipher", "secret"}):
		return fmt.Sprintf(
			"%s may expose sensitive data in transit or at rest. If exploited, attackers could decrypt confidential information, leading to data breaches, regulatory penalties, and loss of customer trust.",
			f.Title,
		)
	case containsAny(title, []string{"network", "firewall", "segment", "perimeter", "dmz", "ingress", "egress"}):
		return fmt.Sprintf(
			"%s may allow unauthorized network access. If exploited, attackers could move laterally across the network, access restricted systems, and exfiltrate sensitive data.",
			f.Title,
		)
	case containsAny(title, []string{"logging", "audit", "monitor", "detect", "alert"}):
		return fmt.Sprintf(
			"%s may reduce visibility into security events. Without adequate monitoring, security incidents could go undetected, allowing attackers to maintain persistence and escalate privileges over extended periods.",
			f.Title,
		)
	case containsAny(cat, []string{"third", "vendor", "supply", "partner"}):
		return fmt.Sprintf(
			"%s may introduce risk through external dependencies. A compromise in the supply chain could lead to unauthorized access to internal systems, data exfiltration, and operational disruption.",
			f.Title,
		)
	case containsAny(title, []string{"patch", "update", "vulnerability", "cve", "exploit"}):
		return fmt.Sprintf(
			"%s may leave systems exposed to known exploits. Unpatched vulnerabilities provide attackers with reliable entry points, potentially leading to full system compromise and data breach.",
			f.Title,
		)
	default:
		return fmt.Sprintf(
			"%s represents a security concern that requires executive attention. If left unaddressed, this issue could impact business operations, regulatory compliance, and data protection.",
			f.Title,
		)
	}
}

// ── PHASE 3 — BUSINESS IMPACT MAPPING ──

func mapBusinessImpact(findings []SDRIFinding, threats []Threat) BusinessImpactMap {
	impactCategories := []string{"Financial", "Operational", "Compliance", "Reputational", "Strategic"}
	impactMap := make(map[string][]string)
	for _, cat := range impactCategories {
		impactMap[cat] = []string{}
	}

	for _, f := range findings {
		title := strings.ToLower(f.Title)
		cat := strings.ToLower(f.Category)
		_ = strings.ToLower(f.Description)

		switch {
		case containsAny(title, []string{"phi", "pii", "personal", "patient", "credit card", "cardholder", "bank"}):
			impactMap["Compliance"] = append(impactMap["Compliance"], f.ID)
			impactMap["Reputational"] = append(impactMap["Reputational"], f.ID)
			impactMap["Financial"] = append(impactMap["Financial"], f.ID)
		case containsAny(title, []string{"identity", "sso", "auth", "oauth", "idp", "federat"}):
			impactMap["Operational"] = append(impactMap["Operational"], f.ID)
			impactMap["Financial"] = append(impactMap["Financial"], f.ID)
		case containsAny(title, []string{"ransomware", "backup", "disaster", "bcdr", "availability", "downtime"}):
			impactMap["Operational"] = append(impactMap["Operational"], f.ID)
			impactMap["Financial"] = append(impactMap["Financial"], f.ID)
			impactMap["Reputational"] = append(impactMap["Reputational"], f.ID)
		case containsAny(cat, []string{"compliance", "regulatory", "audit"}):
			impactMap["Compliance"] = append(impactMap["Compliance"], f.ID)
		case containsAny(cat, []string{"third", "vendor", "supply"}):
			impactMap["Strategic"] = append(impactMap["Strategic"], f.ID)
			impactMap["Operational"] = append(impactMap["Operational"], f.ID)
		case containsAny(title, []string{"encrypt", "key", "crypto", "data protection"}):
			impactMap["Compliance"] = append(impactMap["Compliance"], f.ID)
			impactMap["Financial"] = append(impactMap["Financial"], f.ID)
		default:
			impactMap["Operational"] = append(impactMap["Operational"], f.ID)
		}
	}

	categories := make([]BusinessImpactCategory, 0)
	for _, cat := range impactCategories {
		if len(impactMap[cat]) > 0 {
			categories = append(categories, BusinessImpactCategory{
				Name:     cat,
				Score:    len(impactMap[cat]),
				Findings: impactMap[cat],
			})
		}
	}

	sort.Slice(categories, func(i, j int) bool {
		return categories[i].Score > categories[j].Score
	})

	return BusinessImpactMap{Categories: categories}
}

// ── PHASE 4 — CROWN JEWEL BUSINESS CLASSIFICATION ──

func classifyCrownJewels(pack *KnowledgePack) []CrownJewelClass {
	classes := make([]CrownJewelClass, 0)
	if pack == nil {
		return classes
	}
	for _, j := range pack.CrownJewels {
		classes = append(classes, CrownJewelClass{
			TechnicalName:    j,
			BusinessCategory: classifyBusinessCategory(j),
			BusinessLabel:    classifyBusinessLabel(j),
		})
	}
	return classes
}

func classifyBusinessCategory(name string) string {
	lower := strings.ToLower(name)
	switch {
	case containsAny(lower, []string{"database", "data", "record", "store", "warehouse"}):
		return "Data Asset"
	case containsAny(lower, []string{"auth", "identity", "idp", "credential", "saml", "federat", "login"}):
		return "Identity Infrastructure"
	case containsAny(lower, []string{"kms", "key", "crypto", "certificate"}):
		return "Cryptographic Trust Infrastructure"
	case containsAny(lower, []string{"gateway", "api", "proxy", "load balancer"}):
		return "Network Infrastructure"
	case containsAny(lower, []string{"pod", "cluster", "container", "node", "orchestrat"}):
		return "Compute Infrastructure"
	case containsAny(lower, []string{"network", "firewall", "segment", "vpn"}):
		return "Network Security Infrastructure"
	case containsAny(lower, []string{"controller", "scada", "plc", "hmi", "sensor"}):
		return "Industrial Control Asset"
	default:
		return "Business Asset"
	}
}

func classifyBusinessLabel(name string) string {
	lower := strings.ToLower(name)
	switch {
	case containsAny(lower, []string{"phi", "patient", "ehr", "clinical", "health"}):
		return "Patient Data Asset"
	case containsAny(lower, []string{"payment", "cardholder", "settlement", "transaction", "fraud"}):
		return "Financial Transaction Asset"
	case containsAny(lower, []string{"citizen", "government", "federal", "classified"}):
		return "Government Data Asset"
	case containsAny(lower, []string{"customer", "user", "profile", "account"}):
		return "Customer Identity Asset"
	case containsAny(lower, []string{"admin", "cluster", "privileged"}):
		return "Privileged Access Asset"
	case containsAny(lower, []string{"secret", "key", "vault", "kms"}):
		return "Cryptographic Material Asset"
	case containsAny(lower, []string{"safety", "control", "historian"}):
		return "Safety & Monitoring Asset"
	default:
		return "Critical Business Asset"
	}
}

// ── PHASE 5 — FINANCIAL EXPOSURE ESTIMATION ──

func estimateFinancialExposure(input ERNInput) FinancialExposure {
	score := 0

	// Asset sensitivity
	if input.DomainPack != nil {
		score += len(input.DomainPack.CrownJewels) * 2
	}

	// Attack paths
	score += len(input.AttackPaths) * 3

	// Threat count
	score += len(input.Threats) * 1

	// Compliance exposure
	score += len(input.ComplianceFrameworks) * 2

	// Finding severity
	criticalCount := 0
	highCount := 0
	for _, f := range input.Findings {
		switch strings.ToLower(f.Severity) {
		case "critical":
			criticalCount++
		case "high":
			highCount++
		}
	}
	score += criticalCount * 4
	score += highCount * 2

	level := "Low"
	rationale := "Limited financial exposure based on assessed risk factors."
	switch {
	case score >= 30:
		level = "Severe"
		rationale = "Significant financial exposure due to high asset sensitivity, broad attack surface, and regulatory compliance requirements."
	case score >= 20:
		level = "Significant"
		rationale = "Substantial financial exposure driven by critical assets, multiple attack paths, and compliance obligations."
	case score >= 10:
		level = "Moderate"
		rationale = "Moderate financial exposure from a combination of sensitive assets, threats, and compliance factors."
	}

	return FinancialExposure{Level: level, Rationale: rationale}
}

// ── PHASE 6 — REGULATORY IMPACT ANALYSIS ──

func analyzeRegulatoryImpact(input ERNInput) []RegulatoryImpact {
	impacts := make([]RegulatoryImpact, 0)
	added := make(map[string]bool)

	if input.DomainPack != nil {
		for _, fw := range input.DomainPack.ComplianceFrameworks {
			key := strings.ToLower(fw)
			if added[key] {
				continue
			}
			added[key] = true
			domain := input.DomainPack.Industry
			exposure := "Potential"
			rationale := fmt.Sprintf(
				"Architecture falls within %s domain which is typically subject to %s requirements. Review recommended to confirm applicability.",
				domain, fw,
			)
			impacts = append(impacts, RegulatoryImpact{
				Framework: fw,
				Domain:    domain,
				Exposure:  exposure,
				Rationale: rationale,
			})
		}
	}

	for _, fw := range input.ComplianceFrameworks {
		key := strings.ToLower(fw)
		if added[key] {
			continue
		}
		added[key] = true
		impacts = append(impacts, RegulatoryImpact{
			Framework: fw,
			Domain:    "Cross-Domain",
			Exposure:  "Potential",
			Rationale: fmt.Sprintf(
				"Architecture references %s compliance framework. Potential regulatory impact exists and should be validated with compliance team.",
				fw,
			),
		})
	}

	return impacts
}

// ── PHASE 1+2+3+7 — BUILD EXECUTIVE RISKS ──

func buildExecutiveRisks(input ERNInput) []ExecutiveRisk {
	risks := make([]ExecutiveRisk, 0)

	for i, f := range input.Findings {
		narrative := buildNarrativeFromFinding(f)
		bizImpact := buildBusinessImpactComponent(f)
		opsImpact := buildOperationalImpactComponent(f)
		compImpact := buildComplianceImpactComponent(f)
		finImpact := buildFinancialImpactComponent(f)
		repImpact := buildReputationImpactComponent(f)

		likelihood := estimateLikelihood(f)
		severity := f.Severity
		if severity == "" {
			severity = "Medium"
		}
		priority := computePriority(bizImpact, opsImpact, compImpact, severity)

		assets := findAffectedAssets(f, input.DomainPack)
		controls := findAffectedControls(f, input.Controls)

		risks = append(risks, ExecutiveRisk{
			ID:                 fmt.Sprintf("ERN-RISK-%03d", i+1),
			Title:              f.Title,
			Summary:            narrative,
			BusinessImpact:     bizImpact,
			OperationalImpact:  opsImpact,
			ComplianceImpact:   compImpact,
			FinancialImpact:    finImpact,
			ReputationImpact:   repImpact,
			Likelihood:         likelihood,
			Severity:           severity,
			Priority:           priority,
			AffectedAssets:     assets,
			AffectedControls:   controls,
			RecommendedActions: buildRecommendedActions(f),
		})
	}

	if len(risks) == 0 {
		risks = append(risks, ExecutiveRisk{
			ID:                 "ERN-RISK-001",
			Title:              "Architecture Security Posture Review",
			Summary:            "The architecture demonstrates foundational security controls. Continued monitoring and periodic review are recommended to maintain security posture.",
			BusinessImpact:     "No significant business impact identified at this time.",
			OperationalImpact:  "Routine operational security measures are in place.",
			ComplianceImpact:   "Compliance posture appears aligned with common standards.",
			FinancialImpact:    "No material financial exposure identified.",
			ReputationImpact:   "No material reputation risk identified.",
			Likelihood:         "Low",
			Severity:           "Low",
			Priority:           "Low",
			AffectedAssets:     []string{},
			AffectedControls:   []string{},
			RecommendedActions: []string{"Continue monitoring security posture", "Periodic review of controls and compliance"},
		})
	}

	return risks
}

func buildBusinessImpactComponent(f SDRIFinding) string {
	title := strings.ToLower(f.Title)
	switch {
	case containsAny(title, []string{"phi", "pii", "patient", "health"}):
		return "Exposure of protected health information could result in regulatory action, loss of patient trust, and operational disruption."
	case containsAny(title, []string{"payment", "card", "transaction", "settlement"}):
		return "Compromise of financial transaction systems could result in unauthorized transactions, financial loss, and regulatory penalties."
	case containsAny(title, []string{"identity", "auth", "sso", "login"}):
		return "Compromise of identity infrastructure could result in unauthorized access to all business systems and data."
	case containsAny(title, []string{"key", "encrypt", "crypto", "kms"}):
		return "Weakness in cryptographic controls could result in loss of confidentiality and integrity for all protected data."
	case containsAny(title, []string{"third", "vendor", "supply"}):
		return "Third-party dependencies may introduce risks that affect business continuity and data protection."
	default:
		return "This finding represents a security concern that could impact business operations if exploited."
	}
}

func buildOperationalImpactComponent(f SDRIFinding) string {
	title := strings.ToLower(f.Title)
	switch {
	case containsAny(title, []string{"backup", "disaster", "bcdr", "availability", "redundant"}):
		return "Operational resilience may be compromised, potentially leading to extended downtime and service disruption."
	case containsAny(title, []string{"monitor", "logging", "audit", "detect"}):
		return "Limited visibility into security events may delay incident detection and response, prolonging operational impact."
	case containsAny(title, []string{"network", "firewall", "segment"}):
		return "Network segmentation weaknesses could allow lateral movement, increasing the scope of any compromise."
	default:
		return "May affect day-to-day security operations and require additional resources to address."
	}
}

func buildComplianceImpactComponent(f SDRIFinding) string {
	title := strings.ToLower(f.Title)
	cat := strings.ToLower(f.Category)
	switch {
	case containsAny(title, []string{"phi", "hipaa", "patient", "ehr"}):
		return "Potential HIPAA Privacy and Security Rule exposure. May require breach notification and regulatory reporting."
	case containsAny(title, []string{"pci", "cardholder", "payment"}):
		return "Potential PCI DSS compliance exposure. May affect scope of assessment and cardholder data environment validation."
	case containsAny(title, []string{"fedramp", "nist", "government", "fisma"}):
		return "Potential FedRAMP / FISMA compliance exposure. May affect authorization and continuous monitoring status."
	case containsAny(cat, []string{"compliance", "regulatory", "audit"}):
		return "Direct compliance exposure identified. Requires immediate review and remediation planning."
	default:
		return "May have compliance implications depending on applicable regulatory frameworks."
	}
}

func buildFinancialImpactComponent(f SDRIFinding) string {
	title := strings.ToLower(f.Title)
	switch {
	case containsAny(title, []string{"payment", "fraud", "settlement"}):
		return "Direct financial loss risk from fraudulent transactions or settlement manipulation."
	case containsAny(title, []string{"ransomware", "backup"}):
		return "Potential financial loss from ransom demands, recovery costs, and operational downtime."
	case containsAny(title, []string{"phi", "pii", "breach"}):
		return "Potential financial exposure from breach notification costs, regulatory fines, and legal expenses."
	default:
		return "Indirect financial exposure through operational disruption or remediation costs."
	}
}

func buildReputationImpactComponent(f SDRIFinding) string {
	title := strings.ToLower(f.Title)
	switch {
	case containsAny(title, []string{"phi", "pii", "patient", "health"}):
		return "High reputation risk from potential data breach involving sensitive personal information."
	case containsAny(title, []string{"breach", "data", "exfiltrat"}):
		return "Significant reputation damage from potential data breach and customer trust erosion."
	case containsAny(title, []string{"third", "vendor", "supply"}):
		return "Reputation risk from potential third-party breach reflecting on the organization."
	default:
		return "Moderate reputation risk if the finding is exploited and becomes public."
	}
}

func estimateLikelihood(f SDRIFinding) string {
	title := strings.ToLower(f.Title)
	switch {
	case containsAny(title, []string{"missing", "absent", "none", "lack"}):
		return "High"
	case containsAny(title, []string{"weak", "insufficient", "partial", "limited"}):
		return "Medium"
	default:
		return "Low"
	}
}

func findAffectedAssets(f SDRIFinding, pack *KnowledgePack) []string {
	assets := make([]string, 0)
	if pack != nil {
		title := strings.ToLower(f.Title)
		for _, j := range pack.CrownJewels {
			if containsAny(title, []string{strings.ToLower(j)}) {
				assets = append(assets, j)
			}
		}
	}
	if len(assets) == 0 {
		assets = append(assets, f.AffectedComponents...)
	}
	if len(assets) == 0 {
		assets = append(assets, "Undetermined")
	}
	return assets
}

func findAffectedControls(f SDRIFinding, controls []SDRIControl) []string {
	affected := make([]string, 0)
	title := strings.ToLower(f.Title)
	for _, c := range controls {
		if containsAny(title, []string{strings.ToLower(c.Name)}) {
			affected = append(affected, c.ID)
		}
	}
	if len(affected) == 0 && len(f.AffectedControls) > 0 {
		affected = f.AffectedControls
	}
	return affected
}

func buildRecommendedActions(f SDRIFinding) []string {
	actions := make([]string, 0)
	if f.Recommendation != "" {
		actions = append(actions, f.Recommendation)
	}
	title := strings.ToLower(f.Title)
	switch {
	case containsAny(title, []string{"mfa", "authentication"}):
		actions = append(actions, "Implement multi-factor authentication for all administrative and privileged access.")
	case containsAny(title, []string{"encrypt", "key", "crypto"}):
		actions = append(actions, "Deploy encryption for data at rest and in transit using approved algorithms.")
		actions = append(actions, "Implement automated key rotation and HSM-backed key management.")
	case containsAny(title, []string{"network", "firewall", "segment"}):
		actions = append(actions, "Review and harden network segmentation boundaries.")
		actions = append(actions, "Implement least-privilege network access controls.")
	case containsAny(title, []string{"logging", "audit", "monitor"}):
		actions = append(actions, "Deploy centralized logging and monitoring solution.")
		actions = append(actions, "Establish incident detection and response procedures.")
	case containsAny(title, []string{"backup", "disaster", "bcdr"}):
		actions = append(actions, "Implement automated backup and disaster recovery procedures.")
		actions = append(actions, "Test backup restoration and business continuity plans.")
	case containsAny(title, []string{"third", "vendor", "supply"}):
		actions = append(actions, "Establish vendor risk management program.")
		actions = append(actions, "Conduct third-party security assessments.")
	default:
		actions = append(actions, "Review and remediate based on security best practices.")
	}
	return actions
}

// ── PHASE 7 — EXECUTIVE PRIORITY ENGINE ──

func rankExecutivePriorities(risks []ExecutiveRisk) []PriorityRisk {
	prioritized := make([]PriorityRisk, 0)

	for _, r := range risks {
		score := 0

		// Business impact weight
		if strings.Contains(strings.ToLower(r.BusinessImpact), "critical") ||
			strings.Contains(strings.ToLower(r.BusinessImpact), "significant") {
			score += 15
		}

		// Severity weight
		switch strings.ToLower(r.Severity) {
		case "critical":
			score += 20
		case "high":
			score += 12
		case "medium":
			score += 6
		}

		// Likelihood weight
		switch strings.ToLower(r.Likelihood) {
		case "high":
			score += 10
		case "medium":
			score += 5
		}

		// Asset sensitivity (inferred from affected assets count)
		score += len(r.AffectedAssets) * 3
		score += len(r.AffectedControls) * 2

		// Compliance impact weight
		if strings.Contains(strings.ToLower(r.ComplianceImpact), "potential") {
			score += 8
		}

		priority := "Low"
		switch {
		case score >= 30:
			priority = "Immediate"
		case score >= 20:
			priority = "High"
		case score >= 10:
			priority = "Medium"
		}

		prioritized = append(prioritized, PriorityRisk{
			Risk:     r,
			Priority: priority,
			Score:    score,
		})
	}

	sort.Slice(prioritized, func(i, j int) bool {
		if prioritized[i].Score != prioritized[j].Score {
			return prioritized[i].Score > prioritized[j].Score
		}
		return prioritized[i].Risk.ID < prioritized[j].Risk.ID
	})

	return prioritized
}

// ── PHASE 8 — RISK AGGREGATION ──

func aggregateRiskThemes(risks []ExecutiveRisk, findings []SDRIFinding) []RiskTheme {
	themeMap := map[string]*RiskTheme{
		"Identity Risk": {
			Name:        "Identity Risk",
			Description: "Risks related to identity management, authentication, authorization, and privileged access controls.",
		},
		"Data Protection Risk": {
			Name:        "Data Protection Risk",
			Description: "Risks related to data encryption, key management, data classification, and data loss prevention.",
		},
		"Third Party Risk": {
			Name:        "Third Party Risk",
			Description: "Risks related to vendor management, supply chain security, and third-party dependencies.",
		},
		"Operational Resilience Risk": {
			Name:        "Operational Resilience Risk",
			Description: "Risks related to business continuity, disaster recovery, availability, and incident response.",
		},
		"Compliance Risk": {
			Name:        "Compliance Risk",
			Description: "Risks related to regulatory compliance, audit readiness, and framework alignment.",
		},
		"Network Security Risk": {
			Name:        "Network Security Risk",
			Description: "Risks related to network segmentation, firewall management, and perimeter controls.",
		},
	}

	for _, f := range findings {
		title := strings.ToLower(f.Title)
		cat := strings.ToLower(f.Category)

		mapped := false
		for themeName, theme := range themeMap {
			matched := false
			switch themeName {
			case "Identity Risk":
				matched = containsAny(title, []string{"mfa", "auth", "identity", "sso", "credential", "login", "password", "privileged"})
			case "Data Protection Risk":
				matched = containsAny(title, []string{"encrypt", "key", "crypto", "data", "phi", "pii", "classification", "dlp"})
			case "Third Party Risk":
				matched = containsAny(cat, []string{"third", "vendor", "supply"}) ||
					containsAny(title, []string{"third", "vendor", "supply", "partner"})
			case "Operational Resilience Risk":
				matched = containsAny(title, []string{"backup", "disaster", "bcdr", "availability", "incident", "monitor", "logging"})
			case "Compliance Risk":
				matched = containsAny(cat, []string{"compliance", "regulatory", "audit"}) ||
					containsAny(title, []string{"hipaa", "pci", "fedramp", "sox", "gdpr"})
			case "Network Security Risk":
				matched = containsAny(title, []string{"network", "firewall", "segment", "perimeter", "dmz"})
			}
			if matched {
				theme.RiskCount++
				theme.Findings = append(theme.Findings, f.ID)
				mapped = true
				break
			}
		}

		if !mapped {
			theme := themeMap["Operational Resilience Risk"]
			theme.RiskCount++
			theme.Findings = append(theme.Findings, f.ID)
		}
	}

	for _, f := range findings {
		for _, r := range risks {
			title := strings.ToLower(r.Title)
			fTitle := strings.ToLower(f.Title)
			if strings.Contains(title, fTitle) || strings.Contains(title, strings.ToLower(f.ID)) {
				break
			}
		}
	}

	themes := make([]RiskTheme, 0)
	for _, theme := range themeMap {
		theme.Severity = computeThemeSeverity(theme, findings)
		if theme.RiskCount > 0 {
			themes = append(themes, *theme)
		}
	}

	sort.Slice(themes, func(i, j int) bool {
		return themes[i].RiskCount > themes[j].RiskCount
	})

	return themes
}

func computeThemeSeverity(theme *RiskTheme, findings []SDRIFinding) string {
	severityScore := 0
	for _, fid := range theme.Findings {
		for _, f := range findings {
			if f.ID == fid {
				switch strings.ToLower(f.Severity) {
				case "critical":
					severityScore += 10
				case "high":
					severityScore += 6
				case "medium":
					severityScore += 3
				default:
					severityScore += 1
				}
			}
		}
	}
	switch {
	case severityScore >= 20:
		return "Critical"
	case severityScore >= 10:
		return "High"
	case severityScore >= 5:
		return "Medium"
	default:
		return "Low"
	}
}

// ── PHASE 9 — BOARD SUMMARY GENERATOR ──

func generateBoardSummary(result *ERNRunResult) BoardSummary {
	totalRisks := len(result.ExecutiveRisks)
	criticalCount := 0
	highCount := 0
	for _, pr := range result.PriorityRisks {
		switch pr.Priority {
		case "Immediate":
			criticalCount++
		case "High":
			highCount++
		}
	}

	summary := fmt.Sprintf(
		"The architecture assessment identified %d executive-level risk items, ",
		totalRisks,
	)

	if criticalCount > 0 || highCount > 0 {
		summary += fmt.Sprintf(
			"with %d requiring immediate attention and %d at high priority. ",
			criticalCount, highCount,
		)
	}

	if len(result.RiskThemes) > 0 {
		topTheme := result.RiskThemes[0]
		summary += fmt.Sprintf(
			"The primary risk area is %s, with %d related findings. ",
			topTheme.Name, topTheme.RiskCount,
		)
		if len(result.RiskThemes) > 1 {
			secondaryAreas := make([]string, 0)
			for i, t := range result.RiskThemes {
				if i > 0 && i < 4 {
					secondaryAreas = append(secondaryAreas, t.Name)
				}
			}
			if len(secondaryAreas) > 0 {
				summary += fmt.Sprintf(
					"Additional risk areas include %s. ",
					strings.Join(secondaryAreas, ", "),
				)
			}
		}
	}

	summary += "The architecture demonstrates foundational security controls. However, the identified risk areas represent elevated concerns that may affect sensitive data protection, regulatory readiness, and operational resilience."

	if result.FinancialExposure.Level != "Low" {
		summary += fmt.Sprintf(
			" Financial exposure is classified as %s, warranting board-level awareness and prioritization.",
			result.FinancialExposure.Level,
		)
	}

	return BoardSummary{Summary: summary}
}

// ── PHASE 10 — CISO BRIEFING ──

func generateCISOBriefing(result *ERNRunResult, input ERNInput) CISOBriefing {
	topRisks := make([]string, 0)
	topRemediations := make([]string, 0)
	highRiskAssets := make([]string, 0)

	// Top 5 risks
	for i, pr := range result.PriorityRisks {
		if i >= 5 {
			break
		}
		topRisks = append(topRisks, fmt.Sprintf(
			"%s [%s] %s", pr.Risk.ID, pr.Priority, pr.Risk.Title,
		))
	}

	if len(topRisks) == 0 {
		topRisks = append(topRisks, "No critical risks identified.")
	}

	// Top 5 remediations (deduplicated)
	seenRem := make(map[string]bool)
	for _, pr := range result.PriorityRisks {
		for _, action := range pr.Risk.RecommendedActions {
			if len(topRemediations) >= 5 {
				break
			}
			if !seenRem[action] {
				seenRem[action] = true
				topRemediations = append(topRemediations, action)
			}
		}
		if len(topRemediations) >= 5 {
			break
		}
	}

	if len(topRemediations) == 0 {
		topRemediations = append(topRemediations, "Continue monitoring and periodic review.")
	}

	// High risk assets from crown jewels
	if input.DomainPack != nil {
		for _, j := range input.DomainPack.CrownJewels {
			highRiskAssets = append(highRiskAssets, j)
		}
	}

	// Coverage overview
	coverage := CoverageOverview{}
	for _, c := range input.Controls {
		coverage.TotalControls++
		switch strings.ToLower(c.Coverage) {
		case "full", "complete", "implemented":
			coverage.Covered++
		case "partial":
			coverage.Partial++
		case "missing", "":
			coverage.Missing++
		default:
			coverage.Covered++
		}
	}
	if coverage.TotalControls > 0 {
		coverage.CoverageRate = float64(coverage.Covered) / float64(coverage.TotalControls) * 100
	}
	if coverage.TotalControls == 0 && input.DomainPack != nil {
		coverage.TotalControls = len(input.DomainPack.ExpectedControls)
	}

	// Compliance overview
	complianceOverview := "Compliance readiness assessment not available."
	if len(result.RegulatoryImpacts) > 0 {
		frameworks := make([]string, 0)
		for _, ri := range result.RegulatoryImpacts {
			frameworks = append(frameworks, ri.Framework)
		}
		complianceOverview = fmt.Sprintf(
			"Architecture intersects with %d regulatory frameworks: %s. Each framework requires validation of applicable controls and evidence collection.",
			len(result.RegulatoryImpacts), strings.Join(frameworks, ", "),
		)
	}

	return CISOBriefing{
		TopRisks:           topRisks,
		TopRemediations:    topRemediations,
		HighRiskAssets:     highRiskAssets,
		CoverageOverview:   coverage,
		ComplianceOverview: complianceOverview,
	}
}

// ── PHASE 11 — REMEDIATION ROADMAP ──

func generateRemediationRoadmap(input ERNInput) RemediationRoadmap {
	roadmap := RemediationRoadmap{}

	seen := make(map[string]bool)

	addIfNotSeen := func(action string, category string, priority string, phase *[]RemediationItem) {
		key := strings.ToLower(action)
		if !seen[key] {
			seen[key] = true
			*phase = append(*phase, RemediationItem{
				Action:   action,
				Category: category,
				Priority: priority,
			})
		}
	}

	for _, f := range input.Findings {
		title := strings.ToLower(f.Title)
		actions := buildRecommendedActions(f)

		// 30-day: critical severity, authentication, encryption gaps
		if strings.ToLower(f.Severity) == "critical" ||
			containsAny(title, []string{"mfa", "auth", "credential", "encrypt"}) {
			for _, a := range actions {
				addIfNotSeen(a, f.Category, "Immediate", &roadmap.Phase30)
			}
		}

		// 90-day: high severity, network, logging
		if strings.ToLower(f.Severity) == "high" && len(roadmap.Phase30) == 0 {
			for _, a := range actions {
				addIfNotSeen(a, f.Category, "High", &roadmap.Phase90)
			}
		}
	}

	// Ensure we have baseline items if empty
	if len(roadmap.Phase30) == 0 {
		roadmap.Phase30 = append(roadmap.Phase30,
			RemediationItem{Action: "Implement multi-factor authentication for administrative access", Category: "Identity", Priority: "Immediate"},
			RemediationItem{Action: "Deploy encryption for sensitive data at rest and in transit", Category: "Data Protection", Priority: "Immediate"},
		)
	}
	if len(roadmap.Phase90) == 0 {
		roadmap.Phase90 = append(roadmap.Phase90,
			RemediationItem{Action: "Establish centralized logging and monitoring", Category: "Detection", Priority: "High"},
			RemediationItem{Action: "Implement network segmentation and access controls", Category: "Network", Priority: "High"},
		)
	}
	if len(roadmap.Phase180) == 0 {
		roadmap.Phase180 = append(roadmap.Phase180,
			RemediationItem{Action: "Deploy key management and certificate automation", Category: "Cryptography", Priority: "Medium"},
			RemediationItem{Action: "Establish vendor assurance program", Category: "Third Party", Priority: "Medium"},
		)
	}
	if len(roadmap.Phase12m) == 0 {
		roadmap.Phase12m = append(roadmap.Phase12m,
			RemediationItem{Action: "Conduct full security architecture review", Category: "Governance", Priority: "Low"},
			RemediationItem{Action: "Implement continuous compliance monitoring", Category: "Compliance", Priority: "Low"},
		)
	}

	return roadmap
}

// ── PHASE 12 — RISK TREND MODEL ──

func generateRiskTrend(result *ERNRunResult, input ERNInput) RiskTrend {
	riskScore := 0.0
	criticalCount := 0
	totalControls := len(input.Controls)
	coveredControls := 0

	for _, pr := range result.PriorityRisks {
		switch pr.Priority {
		case "Immediate":
			riskScore += 30
		case "High":
			riskScore += 15
		case "Medium":
			riskScore += 5
		}
	}

	if len(result.PriorityRisks) > 0 {
		riskScore = riskScore / float64(len(result.PriorityRisks))
	}

	for _, pr := range result.PriorityRisks {
		if strings.EqualFold(pr.Priority, "Immediate") {
			criticalCount++
		}
	}

	for _, c := range input.Controls {
		if strings.EqualFold(c.Coverage, "Full") || strings.EqualFold(c.Coverage, "Complete") {
			coveredControls++
		}
	}

	coverageRate := 0.0
	if totalControls > 0 {
		coverageRate = float64(coveredControls) / float64(totalControls) * 100
	}

	complianceReadiness := 0.0
	if len(result.RegulatoryImpacts) > 0 {
		complianceReadiness = float64(len(input.ComplianceFrameworks)) / float64(len(result.RegulatoryImpacts)+len(input.ComplianceFrameworks)) * 100
		if complianceReadiness > 100 {
			complianceReadiness = 100
		}
	}

	// Target state assumes improvement
	targetRiskScore := riskScore * 0.5
	targetCoverage := coverageRate + 30
	if targetCoverage > 100 {
		targetCoverage = 100
	}
	targetCompliance := complianceReadiness + 25
	if targetCompliance > 100 {
		targetCompliance = 100
	}
	targetCritical := criticalCount / 2
	if targetCritical < 0 {
		targetCritical = 0
	}

	return RiskTrend{
		CurrentState: RiskTrendState{
			RiskScore:           riskScore,
			CoverageRate:        coverageRate,
			ComplianceReadiness: complianceReadiness,
			CriticalFindings:    criticalCount,
		},
		TargetState: RiskTrendState{
			RiskScore:           targetRiskScore,
			CoverageRate:        targetCoverage,
			ComplianceReadiness: targetCompliance,
			CriticalFindings:    targetCritical,
		},
	}
}

// ── PHASE 13 — SECURITY INVESTMENT INSIGHTS ──

func generateInvestmentInsights(result *ERNRunResult, input ERNInput) []InvestmentInsight {
	insights := make([]InvestmentInsight, 0)
	areaMap := make(map[string]int)

	scoreArea := func(area string, score int) {
		areaMap[area] += score
	}

	for _, pr := range result.PriorityRisks {
		title := strings.ToLower(pr.Risk.Title)
		switch {
		case containsAny(title, []string{"mfa", "auth", "identity", "sso", "credential", "login"}):
			scoreArea("Identity Hardening", 10)
		case containsAny(title, []string{"encrypt", "key", "crypto", "kms"}):
			scoreArea("Key Management", 9)
		case containsAny(title, []string{"log", "monitor", "detect", "audit"}):
			scoreArea("Centralized Logging", 8)
		case containsAny(title, []string{"third", "vendor", "supply"}):
			scoreArea("Third-Party Governance", 7)
		case containsAny(title, []string{"network", "segment", "firewall"}):
			scoreArea("Network Segmentation", 6)
		case containsAny(title, []string{"backup", "disaster", "bcdr"}):
			scoreArea("Business Continuity", 5)
		}
	}

	type areaScore struct {
		name  string
		score int
	}
	sorted := make([]areaScore, 0)
	for name, score := range areaMap {
		sorted = append(sorted, areaScore{name, score})
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].score > sorted[j].score
	})

	for i, as := range sorted {
		priority := "Medium"
		switch {
		case i == 0:
			priority = "Highest"
		case i <= 2:
			priority = "High"
		}
		insights = append(insights, InvestmentInsight{
			Area:     as.name,
			Priority: priority,
			Rationale: fmt.Sprintf(
				"Scored %d points based on frequency and severity of related findings across assessment.",
				as.score,
			),
		})
	}

	if len(insights) == 0 {
		insights = append(insights, InvestmentInsight{
			Area:      "Security Baseline",
			Priority:  "Medium",
			Rationale: "Standard security controls appear adequate. Focus on continuous improvement and monitoring.",
		})
	}

	return insights
}

// ── PHASE 14 — EXECUTIVE DASHBOARD ──

func buildExecutiveDashboard(result *ERNRunResult, input ERNInput) ExecutiveDashboard {
	riskScore := 0.0
	priorityFindings := 0

	for _, pr := range result.PriorityRisks {
		if pr.Priority == "Immediate" || pr.Priority == "High" {
			priorityFindings++
		}
		riskScore += float64(pr.Score)
	}
	if len(result.PriorityRisks) > 0 {
		riskScore = riskScore / float64(len(result.PriorityRisks))
	}

	complianceReadiness := 0.0
	if len(result.RegulatoryImpacts) > 0 {
		complianceReadiness = float64(len(input.ComplianceFrameworks)) / float64(len(result.RegulatoryImpacts)) * 100
		if complianceReadiness > 100 {
			complianceReadiness = 100
		}
	}

	coverageRate := 0.0
	if len(input.Controls) > 0 {
		covered := 0
		for _, c := range input.Controls {
			if c.Coverage == "Full" || c.Coverage == "Complete" || c.Coverage == "Enhanced" || c.Coverage == "Implemented" {
				covered++
			}
		}
		coverageRate = float64(covered) / float64(len(input.Controls)) * 100
	}

	attackPathCount := len(input.AttackPaths)

	criticalAssets := make([]string, 0)
	if input.DomainPack != nil {
		criticalAssets = input.DomainPack.CrownJewels
	}

	return ExecutiveDashboard{
		RiskScore:           riskScore,
		PriorityFindings:    priorityFindings,
		ComplianceReadiness: complianceReadiness,
		CoverageRate:        coverageRate,
		AttackPathCount:     attackPathCount,
		CriticalAssets:      criticalAssets,
	}
}

// ── PHASE 15 — CISO DECISION SUPPORT ──

func generateDecisionSupport(result *ERNRunResult) DecisionSupport {
	actions := make([]DecisionAction, 0)
	seen := make(map[string]bool)

	addAction := func(rank int, action string, impact string, rationale string) {
		key := strings.ToLower(action)
		if !seen[key] {
			seen[key] = true
			actions = append(actions, DecisionAction{
				Rank:      rank,
				Action:    action,
				Impact:    impact,
				Rationale: rationale,
			})
		}
	}

	// Collect top actions from highest priority risks
	for _, pr := range result.PriorityRisks {
		if len(actions) >= 3 {
			break
		}
		for _, rec := range pr.Risk.RecommendedActions {
			if len(actions) >= 3 {
				break
			}
			impact := "High"
			if pr.Priority == "Immediate" {
				impact = "Critical"
			}
			addAction(len(actions)+1, rec, impact,
				fmt.Sprintf("Addresses %s priority risk: %s", pr.Priority, pr.Risk.Title))
		}
	}

	// Fallback if no actions from risks
	if len(actions) == 0 {
		addAction(1, "Implement multi-factor authentication for all privileged access",
			"Critical",
			"Single highest-impact control for reducing breach risk across all attack paths.")
		addAction(2, "Deploy centralized logging and security monitoring",
			"High",
			"Essential for incident detection, response, and compliance reporting.")
		addAction(3, "Establish encryption and key management program",
			"High",
			"Protects sensitive data at rest and in transit, reduces breach impact.")
	}

	return DecisionSupport{Top3Actions: actions}
}

// ── CONVENIENCE ──

func containsAny(s string, items []string) bool {
	for _, item := range items {
		if strings.Contains(s, item) {
			return true
		}
	}
	return false
}

func computePriority(bizImpact, opsImpact, compImpact, severity string) string {
	score := 0

	if strings.Contains(strings.ToLower(bizImpact), "critical") ||
		strings.Contains(strings.ToLower(bizImpact), "significant") {
		score += 10
	}
	if strings.Contains(strings.ToLower(opsImpact), "extended") ||
		strings.Contains(strings.ToLower(opsImpact), "significant") {
		score += 8
	}
	if strings.Contains(strings.ToLower(compImpact), "potential") {
		score += 6
	}
	switch strings.ToLower(severity) {
	case "critical":
		score += 12
	case "high":
		score += 8
	case "medium":
		score += 4
	}

	switch {
	case score >= 20:
		return "Immediate"
	case score >= 12:
		return "High"
	case score >= 6:
		return "Medium"
	default:
		return "Low"
	}
}
