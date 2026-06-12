package main

import (
	"time"
)

// EvidenceSource tracks where evidence for an assumption came from.
type EvidenceSource struct {
	FilePath                string   `json:"file_path"`
	FileType                string   `json:"file_type"`
	MatchedComponents       []string `json:"matched_components"`
	MatchedRelationships    []string `json:"matched_relationships"`
	MatchedTrustBoundaries  []string `json:"matched_trust_boundaries"`
	MatchedSecurityConcepts []string `json:"matched_security_concepts"`
}

// EvidenceSummary is the top-level evidence container on AnalysisResult.
type EvidenceSummary struct {
	TotalSources       int      `json:"total_sources"`
	TotalComponents    int      `json:"total_components"`
	TotalRelationships int      `json:"total_relationships"`
	SourceFiles        []string `json:"source_files"`
}

// StrideJustification explains why a STRIDE category was assigned.
type StrideJustification struct {
	Category           StrideCategory `json:"category"`
	Reason             string         `json:"reason"`
	MatchedRuleIndexes []int          `json:"matched_rule_indexes"`
	MatchedKeywords    []string       `json:"matched_keywords"`
	MatchedComponents  []string       `json:"matched_components"`
	Confidence         float64        `json:"confidence"`
	ConfidenceReason   string         `json:"confidence_reason"`
}

// LikelihoodFactor describes a factor that contributes to likelihood.
type LikelihoodFactor struct {
	Factor string `json:"factor"`
	Value  int    `json:"value"`
	Reason string `json:"reason"`
}

// ImpactFactor describes a factor that contributes to impact.
type ImpactFactor struct {
	Factor string `json:"factor"`
	Value  int    `json:"value"`
	Reason string `json:"reason"`
}

// RiskJustification explains how risk was calculated.
type RiskJustification struct {
	Likelihood        int                `json:"likelihood"`
	LikelihoodReason  string             `json:"likelihood_reason"`
	LikelihoodFactors []LikelihoodFactor `json:"likelihood_factors"`
	Impact            int                `json:"impact"`
	ImpactReason      string             `json:"impact_reason"`
	ImpactFactors     []ImpactFactor     `json:"impact_factors"`
	RiskScore         int                `json:"risk_score"`
	RiskLevel         RiskLevel          `json:"risk_level"`
	RiskReason        string             `json:"risk_reason"`
	Confidence        float64            `json:"confidence"`
	ConfidenceReason  string             `json:"confidence_reason"`
}

// ReviewRecord tracks the human review status of an assumption.
type ReviewRecord struct {
	Status    string    `json:"status"` // Proposed, Accepted, Rejected, Modified
	Notes     string    `json:"notes"`
	Timestamp time.Time `json:"timestamp"`
	Reviewer  string    `json:"reviewer"`
}

// ValidationRecord stores data needed for future precision/recall studies.
type ValidationRecord struct {
	AssumptionID      string           `json:"assumption_id"`
	Description       string           `json:"description"`
	GeneratedEvidence []string         `json:"generated_evidence"`
	AssignedRisk      RiskLevel        `json:"assigned_risk"`
	RiskScore         int              `json:"risk_score"`
	Confidence        float64          `json:"confidence"`
	STRIDECategories  []StrideCategory `json:"stride_categories"`
	ArchReviewResult  string           `json:"arch_review_result"` // Accepted, Rejected, Modified
	ArchNotes         string           `json:"arch_notes"`
	ReviewTimestamp   time.Time        `json:"review_timestamp"`
}

// The following fields are added to the existing Assumption struct via a helper
// that returns the explainability info. The struct itself lives in engine.go.
// We use composition: Assumption + ExplainabilityExtension.

// ExplainabilityExtension holds all explainability data for an assumption.
// This is set alongside the Assumption in the analysis pipeline.
type ExplainabilityExtension struct {
	EvidenceSources      []string              `json:"evidence_sources"`
	SourceComponents     []string              `json:"source_components"`
	SourceRelationships  []string              `json:"source_relationships"`
	Rationale            string                `json:"rationale"`
	StrideJustifications []StrideJustification `json:"stride_justifications"`
	RiskJustification    *RiskJustification    `json:"risk_justification"`
	Review               ReviewRecord          `json:"review"`
}

// AttachExplainability attaches an ExplainabilityExtension to an Assumption
// by setting the extra fields that exist on the struct.
// ControlDetail links a control to the specific assumptions and threats it mitigates.
type ControlDetail struct {
	ID                     string           `json:"id"`
	Description            string           `json:"description"`
	Rationale              string           `json:"rationale"`
	Category               string           `json:"category"`
	MitigatedAssumptionIDs []string         `json:"mitigated_assumption_ids"`
	MitigatedSTRIDE        []StrideCategory `json:"mitigated_stride"`
	Priority               int              `json:"priority"` // 1=highest
}

// ──────────────────────────────────────────────
// Attack Path Discovery Engine (APD) Types
// ──────────────────────────────────────────────

// AttackPath represents a complete attacker journey from entry point to target asset.
type AttackPath struct {
	ID                  string       `json:"id"`
	Name                string       `json:"name"`
	Description         string       `json:"description"`
	EntryPoint          string       `json:"entry_point"`
	TargetAsset         string       `json:"target_asset"`
	AttackSteps         []AttackStep `json:"attack_steps"`
	RequiredAssumptions []string     `json:"required_assumptions,omitempty"`
	RequiredConditions  []string     `json:"required_conditions,omitempty"`
	ExploitedThreats    []string     `json:"exploited_threats,omitempty"`
	AffectedComponents  []string     `json:"affected_components,omitempty"`
	AffectedBoundaries  []string     `json:"affected_boundaries,omitempty"`
	Likelihood          float64      `json:"likelihood"`
	Impact              float64      `json:"impact"`
	RiskScore           float64      `json:"risk_score"`
	Confidence          float64      `json:"confidence"`
	DetectionDifficulty string       `json:"detection_difficulty"`
	BusinessImpact      string       `json:"business_impact"`
	Recommendations     []string     `json:"recommendations,omitempty"`
	KillChainPhases     []string     `json:"kill_chain_phases,omitempty"`
	MITREATTACK         []string     `json:"mitre_attack,omitempty"`
	STRIDECategories    []string     `json:"stride_categories,omitempty"`
}

// AttackStep represents a single step in an attack path.
type AttackStep struct {
	SequenceNumber     int    `json:"sequence_number"`
	SourceComponent    string `json:"source_component"`
	TargetComponent    string `json:"target_component"`
	Action             string `json:"action"`
	Threat             string `json:"threat"`
	RequiredAssumption string `json:"required_assumption"`
	ControlBypassed    string `json:"control_bypassed"`
	Reasoning          string `json:"reasoning"`
	STRIDECategory     string `json:"stride_category"`
}

// ThreatChain represents a chain of connected threats along an attack path.
type ThreatChain struct {
	ID        string   `json:"id"`
	Threats   []string `json:"threats"`
	Path      []string `json:"path"`
	RiskScore float64  `json:"risk_score"`
	Reasoning string   `json:"reasoning"`
}

// AttackPathSummary summarizes attack path discovery results.
type AttackPathSummary struct {
	TotalAttackPaths  int            `json:"total_attack_paths"`
	CriticalCount     int            `json:"critical_count"`
	HighCount         int            `json:"high_count"`
	MediumCount       int            `json:"medium_count"`
	LowCount          int            `json:"low_count"`
	ThreatChainCount  int            `json:"threat_chain_count"`
	TopAttackPaths    []string       `json:"top_attack_paths"`
	KillChainCoverage map[string]int `json:"kill_chain_coverage"`
	MITRECoverage     []string       `json:"mitre_coverage"`
	SummaryText       string         `json:"summary_text"`
}

func attachExplainability(a *Assumption, ext *ExplainabilityExtension) {
	if ext == nil {
		return
	}
	a.EvidenceSources = ext.EvidenceSources
	a.SourceComponents = ext.SourceComponents
	a.SourceRelationships = ext.SourceRelationships
	a.Rationale = ext.Rationale
	a.StrideJustifications = ext.StrideJustifications
	a.RiskJustification = ext.RiskJustification
	a.ReviewStatus = ext.Review.Status
	a.ReviewNotes = ext.Review.Notes
	a.ReviewTimestamp = ext.Review.Timestamp
}

// ─────────────────────────────────────────────────────────────
// SDRIControl represents a security control with coverage status.
type SDRIControl struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Category    string   `json:"category"`
	Description string   `json:"description"`
	ControlType string   `json:"control_type"`
	Preventive  bool     `json:"preventive"`
	Detective   bool     `json:"detective"`
	Corrective  bool     `json:"corrective"`
	Strength    string   `json:"strength"`
	Evidence    []string `json:"evidence,omitempty"`
	Coverage    string   `json:"coverage"`
	Status      string   `json:"status"`
}

// SDRIDesignFinding represents a security design review finding.
type SDRIDesignFinding struct {
	ID                 string   `json:"id"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	Severity           string   `json:"severity"`
	Category           string   `json:"category"`
	AffectedComponents []string `json:"affected_components,omitempty"`
	AffectedControls   []string `json:"affected_controls,omitempty"`
	BusinessImpact     string   `json:"business_impact"`
	Recommendation     string   `json:"recommendation"`
	Reasoning          string   `json:"reasoning"`
}

// SDRIArchitecturalWeakness represents an architectural weakness.
type SDRIArchitecturalWeakness struct {
	ID             string   `json:"id"`
	Pattern        string   `json:"pattern"`
	Description    string   `json:"description"`
	Severity       string   `json:"severity"`
	Components     []string `json:"components,omitempty"`
	Impact         string   `json:"impact"`
	Recommendation string   `json:"recommendation"`
}

// SDRIRemediation represents a prioritized remediation.
type SDRIRemediation struct {
	ID                 string   `json:"id"`
	Priority           int      `json:"priority"`
	Description        string   `json:"description"`
	RiskScore          float64  `json:"risk_score"`
	BusinessImpact     string   `json:"business_impact"`
	Effort             string   `json:"effort"`
	Category           string   `json:"category"`
	Recommendation     string   `json:"recommendation"`
	AffectedComponents []string `json:"affected_components,omitempty"`
}

// SDRICoverageItem represents control coverage for a category.
type SDRICoverageItem struct {
	Category string  `json:"category"`
	Expected int     `json:"expected"`
	Observed int     `json:"observed"`
	Coverage float64 `json:"coverage"`
	Level    string  `json:"level"`
}

// SDRIComplianceMapping represents control coverage for a compliance framework.
type SDRIComplianceMapping struct {
	Framework string   `json:"framework"`
	Coverage  float64  `json:"coverage"`
	Controls  []string `json:"controls,omitempty"`
	Status    string   `json:"status"`
}

// ── CIARE — Compliance Intelligence & Audit Readiness Types ──

// CIAREFrameworkCoverage represents per-framework coverage analysis.
type CIAREFrameworkCoverage struct {
	Framework        string   `json:"framework"`
	Required         int      `json:"required"`
	Observed         int      `json:"observed"`
	Missing          int      `json:"missing"`
	CoveragePct      float64  `json:"coverage_pct"`
	Status           string   `json:"status"`
	ObservedControls []string `json:"observed_controls,omitempty"`
	MissingControls  []string `json:"missing_controls,omitempty"`
}

// CIAREAuditReadiness represents audit readiness scoring per framework.
type CIAREAuditReadiness struct {
	Framework       string   `json:"framework"`
	ReadinessScore  float64  `json:"readiness_score"`
	Status          string   `json:"status"`
	ControlCoverage float64  `json:"control_coverage"`
	EvidenceScore   float64  `json:"evidence_score"`
	ThreatExposure  float64  `json:"threat_exposure"`
	FindingsPenalty float64  `json:"findings_penalty"`
	Factors         []string `json:"factors,omitempty"`
}

// CIAREEvidenceRequirement represents evidence needed for a control under a framework.
type CIAREEvidenceRequirement struct {
	Framework string   `json:"framework"`
	Control   string   `json:"control"`
	Evidence  []string `json:"evidence"`
}

// CIAREMissingEvidence represents a control that exists but lacks supporting evidence.
type CIAREMissingEvidence struct {
	Framework string   `json:"framework"`
	Control   string   `json:"control"`
	Evidences []string `json:"evidences"`
}

// CIAREAuditorQuestion represents a likely auditor question.
type CIAREAuditorQuestion struct {
	Framework string `json:"framework"`
	Control   string `json:"control"`
	Question  string `json:"question"`
}

// CIAREComplianceGap represents a compliance gap per framework.
type CIAREComplianceGap struct {
	ID          string `json:"id"`
	Framework   string `json:"framework"`
	Requirement string `json:"requirement"`
	Observed    string `json:"observed"`
	Missing     string `json:"missing"`
	Risk        string `json:"risk"`
}

// CIAREControlMaturity represents maturity level for a control domain.
type CIAREControlMaturity struct {
	Domain   string  `json:"domain"`
	Level    int     `json:"level"`
	Label    string  `json:"label"`
	Coverage float64 `json:"coverage"`
}

// CIAREComplianceNarrative represents executive narrative per framework.
type CIAREComplianceNarrative struct {
	Framework string `json:"framework"`
	Narrative string `json:"narrative"`
}

// CIAREAuditPackage represents a complete audit readiness package.
type CIAREAuditPackage struct {
	ExecutiveSummary     string                     `json:"executive_summary"`
	FrameworkCoverages   []CIAREFrameworkCoverage   `json:"framework_coverages,omitempty"`
	ControlInventory     []SDRIControl              `json:"control_inventory,omitempty"`
	MissingControls      []CIAREComplianceGap       `json:"missing_controls,omitempty"`
	EvidenceRequirements []CIAREEvidenceRequirement `json:"evidence_requirements,omitempty"`
	AuditorQuestions     []CIAREAuditorQuestion     `json:"auditor_questions,omitempty"`
}

// CIAREComplianceDashboard represents the compliance dashboard view.
type CIAREComplianceDashboard struct {
	FrameworkCoverages map[string]float64     `json:"framework_coverages"`
	TopGaps            []CIAREComplianceGap   `json:"top_gaps,omitempty"`
	TopMissingEvidence []CIAREMissingEvidence `json:"top_missing_evidence,omitempty"`
	TopRisks           []string               `json:"top_risks,omitempty"`
}

// CIAREProcurementQuestion represents a vendor security review question.
type CIAREProcurementQuestion struct {
	Category string `json:"category"`
	Question string `json:"question"`
}

// ── DKPI — Domain Knowledge Pack Intelligence Types ──

type DKPIDomainResult struct {
	PrimaryDomain string            `json:"primary_domain"`
	Confidence    float64           `json:"confidence"`
	Rationale     []string          `json:"rationale,omitempty"`
	Matches       []DKPIDomainMatch `json:"matches,omitempty"`
}

type DKPIDomainMatch struct {
	PackID     string   `json:"pack_id"`
	PackName   string   `json:"pack_name"`
	Score      int      `json:"score"`
	Confidence float64  `json:"confidence"`
	Reasons    []string `json:"reasons,omitempty"`
}

type DKPIKnowledgePack struct {
	ID                   string                     `json:"id"`
	Name                 string                     `json:"name"`
	Industry             string                     `json:"industry"`
	Description          string                     `json:"description"`
	CrownJewels          []string                   `json:"crown_jewels,omitempty"`
	ExpectedControls     []DKPIKnowledgePackControl `json:"expected_controls,omitempty"`
	ThreatPatterns       []DKPIKnowledgePackThreat  `json:"threat_patterns,omitempty"`
	ComplianceFrameworks []string                   `json:"compliance_frameworks,omitempty"`
}

type DKPIKnowledgePackControl struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Priority    string `json:"priority"`
}

type DKPIKnowledgePackEvidence struct {
	Control  string   `json:"control"`
	Evidence []string `json:"evidence,omitempty"`
}

type DKPIKnowledgePackThreat struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Category    string `json:"category"`
}

type DKPIKnowledgePackAttackPath struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
	Target      string   `json:"target"`
}

type DKPIIntelligence struct {
	DomainResult       DKPIDomainResult            `json:"domain_result"`
	ActivePack         *DKPIKnowledgePack          `json:"active_pack,omitempty"`
	Recommendations    []string                    `json:"recommendations,omitempty"`
	InjectedThreats    []Threat                    `json:"injected_threats,omitempty"`
	DomainControls     []SDRIControl               `json:"domain_controls,omitempty"`
	DomainCompliance   []string                    `json:"domain_compliance,omitempty"`
	EvidenceReqs       []DKPIKnowledgePackEvidence `json:"evidence_requirements,omitempty"`
	BoostedAssumptions []Assumption                `json:"boosted_assumptions,omitempty"`
	Summary            string                      `json:"summary,omitempty"`
}

// ── ERN — Executive Risk Narratives Types ──

type ERNExecutiveRisk struct {
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

type ERNBoardSummary struct {
	Summary string `json:"summary"`
}

type ERNCISOBriefing struct {
	TopRisks           []string            `json:"top_risks"`
	TopRemediations    []string            `json:"top_remediations"`
	HighRiskAssets     []string            `json:"high_risk_assets"`
	CoverageOverview   ERNCoverageOverview `json:"coverage_overview"`
	ComplianceOverview string              `json:"compliance_overview"`
}

type ERNCoverageOverview struct {
	TotalControls int     `json:"total_controls"`
	Covered       int     `json:"covered"`
	Partial       int     `json:"partial"`
	Missing       int     `json:"missing"`
	CoverageRate  float64 `json:"coverage_rate"`
}

type ERNRemediationItem struct {
	Action   string `json:"action"`
	Category string `json:"category"`
	Priority string `json:"priority"`
}

type ERNRemediationRoadmap struct {
	Phase30  []ERNRemediationItem `json:"phase_30"`
	Phase90  []ERNRemediationItem `json:"phase_90"`
	Phase180 []ERNRemediationItem `json:"phase_180"`
	Phase12m []ERNRemediationItem `json:"phase_12m"`
}

type ERNRegulatoryImpact struct {
	Framework string `json:"framework"`
	Domain    string `json:"domain"`
	Exposure  string `json:"exposure"`
	Rationale string `json:"rationale"`
}

type ERNInvestmentInsight struct {
	Area        string `json:"area"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Rationale   string `json:"rationale"`
}

type ERNExecutiveDashboard struct {
	RiskScore           float64  `json:"risk_score"`
	PriorityFindings    int      `json:"priority_findings"`
	ComplianceReadiness float64  `json:"compliance_readiness"`
	CoverageRate        float64  `json:"coverage_rate"`
	AttackPathCount     int      `json:"attack_path_count"`
	CriticalAssets      []string `json:"critical_assets"`
}

type ERNDecisionAction struct {
	Rank      int    `json:"rank"`
	Action    string `json:"action"`
	Impact    string `json:"impact"`
	Rationale string `json:"rationale"`
}

type ERNDecisionSupport struct {
	Top3Actions []ERNDecisionAction `json:"top_3_actions"`
}

type ERNRiskTheme struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	RiskCount   int    `json:"risk_count"`
	Severity    string `json:"severity"`
}

type ERNFinancialExposure struct {
	Level     string `json:"level"`
	Rationale string `json:"rationale"`
}

type ERNCrownJewelClass struct {
	TechnicalName    string `json:"technical_name"`
	BusinessCategory string `json:"business_category"`
	BusinessLabel    string `json:"business_label"`
}

type ERNReportPacks struct {
	BoardReport     string `json:"board_report,omitempty"`
	ExecutiveReport string `json:"executive_report,omitempty"`
	TechnicalReport string `json:"technical_report,omitempty"`
}

type ERNIntelligence struct {
	ExecutiveRisks     []ERNExecutiveRisk     `json:"executive_risks,omitempty"`
	BoardSummary       ERNBoardSummary        `json:"board_summary,omitempty"`
	CISOBriefing       ERNCISOBriefing        `json:"ciso_briefing,omitempty"`
	RemediationRoadmap ERNRemediationRoadmap  `json:"remediation_roadmap,omitempty"`
	RegulatoryImpacts  []ERNRegulatoryImpact  `json:"regulatory_impacts,omitempty"`
	InvestmentInsights []ERNInvestmentInsight `json:"investment_insights,omitempty"`
	Dashboard          ERNExecutiveDashboard  `json:"dashboard,omitempty"`
	DecisionSupport    ERNDecisionSupport     `json:"decision_support,omitempty"`
	RiskThemes         []ERNRiskTheme         `json:"risk_themes,omitempty"`
	FinancialExposure  ERNFinancialExposure   `json:"financial_exposure,omitempty"`
	CrownJewelClasses  []ERNCrownJewelClass   `json:"crown_jewel_classes,omitempty"`
	ReportPacks        ERNReportPacks         `json:"report_packs,omitempty"`
}

// ── SAMPI — Security Architecture Memory & Portfolio Intelligence ──

type SAMPIDashboard struct {
	TotalArchitectures int            `json:"total_architectures,omitempty"`
	TotalFindings      int            `json:"total_findings,omitempty"`
	TotalThreats       int            `json:"total_threats,omitempty"`
	TotalAttackPaths   int            `json:"total_attack_paths,omitempty"`
	TotalControls      int            `json:"total_controls,omitempty"`
	AverageCoverage    float64        `json:"average_coverage,omitempty"`
	AverageRiskScore   float64        `json:"average_risk_score,omitempty"`
	ComplianceCount    int            `json:"compliance_count,omitempty"`
	RiskDistribution   map[string]int `json:"risk_distribution,omitempty"`
}

type SAMPIHeatmap struct {
	ArchitectureName string  `json:"architecture_name,omitempty"`
	RiskScore        float64 `json:"risk_score,omitempty"`
	FindingCount     int     `json:"finding_count,omitempty"`
	ControlCount     int     `json:"control_count,omitempty"`
	ComplianceCount  int     `json:"compliance_count,omitempty"`
	RiskBand         string  `json:"risk_band,omitempty"`
}

type SAMPIRepeatedWeakness struct {
	FindingTitle          string   `json:"finding_title,omitempty"`
	Category              string   `json:"category,omitempty"`
	Severity              string   `json:"severity,omitempty"`
	AffectedArchitectures []string `json:"affected_architectures,omitempty"`
	OccurrenceCount       int      `json:"occurrence_count,omitempty"`
	Systemic              bool     `json:"systemic,omitempty"`
}

type SAMPIEnterpriseTheme struct {
	Name                  string `json:"name,omitempty"`
	Description           string `json:"description,omitempty"`
	RiskCount             int    `json:"risk_count,omitempty"`
	AffectedArchitectures int    `json:"affected_architectures,omitempty"`
	Severity              string `json:"severity,omitempty"`
}

type SAMPIControlCoverage struct {
	ControlName        string  `json:"control_name,omitempty"`
	Category           string  `json:"category,omitempty"`
	ArchitecturesWith  int     `json:"architectures_with,omitempty"`
	ArchitecturesTotal int     `json:"architectures_total,omitempty"`
	CoveragePercent    float64 `json:"coverage_percent,omitempty"`
}

type SAMPIComparison struct {
	ArchitectureA     string  `json:"architecture_a,omitempty"`
	ArchitectureB     string  `json:"architecture_b,omitempty"`
	SharedAssumptions int     `json:"shared_assumptions,omitempty"`
	SharedThreats     int     `json:"shared_threats,omitempty"`
	SharedControls    int     `json:"shared_controls,omitempty"`
	RiskScoreA        float64 `json:"risk_score_a,omitempty"`
	RiskScoreB        float64 `json:"risk_score_b,omitempty"`
	RiskDelta         float64 `json:"risk_delta,omitempty"`
	SimilarityScore   float64 `json:"similarity_score,omitempty"`
}

type SAMPIRiskTrend struct {
	ArchitectureID string  `json:"architecture_id,omitempty"`
	Name           string  `json:"name,omitempty"`
	PreviousScore  float64 `json:"previous_score,omitempty"`
	CurrentScore   float64 `json:"current_score,omitempty"`
	Direction      string  `json:"direction,omitempty"`
}

type SAMPISecurityDebt struct {
	Score             float64         `json:"score,omitempty"`
	LongstandingCount int             `json:"longstanding_count,omitempty"`
	RepeatedCount     int             `json:"repeated_count,omitempty"`
	TopDebts          []SAMPIDebtItem `json:"top_debts,omitempty"`
}

type SAMPIDebtItem struct {
	Description  string `json:"description,omitempty"`
	Architecture string `json:"architecture,omitempty"`
	Category     string `json:"category,omitempty"`
	Severity     string `json:"severity,omitempty"`
	Age          string `json:"age,omitempty"`
}

type SAMPIExposure struct {
	InternetExposure   int `json:"internet_exposure,omitempty"`
	ThirdPartyExposure int `json:"third_party_exposure,omitempty"`
	IdentityExposure   int `json:"identity_exposure,omitempty"`
	CloudExposure      int `json:"cloud_exposure,omitempty"`
	TotalExposure      int `json:"total_exposure,omitempty"`
}

type SAMPICrownJewel struct {
	Name             string `json:"name,omitempty"`
	ArchitectureName string `json:"architecture_name,omitempty"`
	Category         string `json:"category,omitempty"`
	ThreatCount      int    `json:"threat_count,omitempty"`
	RiskLevel        string `json:"risk_level,omitempty"`
}

type SAMISharedDependency struct {
	DependencyName      string   `json:"dependency_name,omitempty"`
	Category            string   `json:"category,omitempty"`
	UsedByArchitectures []string `json:"used_by_architectures,omitempty"`
	UsageCount          int      `json:"usage_count,omitempty"`
	RiskLevel           string   `json:"risk_level,omitempty"`
}

type SAMPIBlastRadius struct {
	ComponentName    string   `json:"component_name,omitempty"`
	ArchitectureName string   `json:"architecture_name,omitempty"`
	FailureImpact    string   `json:"failure_impact,omitempty"`
	AffectedSystems  []string `json:"affected_systems,omitempty"`
	Severity         string   `json:"severity,omitempty"`
}

type SAMPIComplianceFramework struct {
	Framework          string  `json:"framework,omitempty"`
	ArchitecturesWith  int     `json:"architectures_with,omitempty"`
	TotalArchitectures int     `json:"total_architectures,omitempty"`
	Coverage           float64 `json:"coverage,omitempty"`
}

type SAMPIComplianceView struct {
	Frameworks         []SAMPIComplianceFramework `json:"frameworks,omitempty"`
	TotalArchitectures int                        `json:"total_architectures,omitempty"`
}

type SAMPIProgramInsight struct {
	Area      string `json:"area,omitempty"`
	Insight   string `json:"insight,omitempty"`
	Priority  string `json:"priority,omitempty"`
	Rationale string `json:"rationale,omitempty"`
}

type SAMPIIntelligence struct {
	Dashboard          SAMPIDashboard          `json:"dashboard,omitempty"`
	Heatmaps           []SAMPIHeatmap          `json:"heatmaps,omitempty"`
	RepeatedWeaknesses []SAMPIRepeatedWeakness `json:"repeated_weaknesses,omitempty"`
	EnterpriseThemes   []SAMPIEnterpriseTheme  `json:"enterprise_themes,omitempty"`
	ControlCoverage    []SAMPIControlCoverage  `json:"control_coverage,omitempty"`
	Comparisons        []SAMPIComparison       `json:"comparisons,omitempty"`
	RiskTrends         []SAMPIRiskTrend        `json:"risk_trends,omitempty"`
	SecurityDebt       SAMPISecurityDebt       `json:"security_debt,omitempty"`
	AttackSurface      SAMPIExposure           `json:"attack_surface,omitempty"`
	CrownJewels        []SAMPICrownJewel       `json:"crown_jewels,omitempty"`
	SharedDependencies []SAMISharedDependency  `json:"shared_dependencies,omitempty"`
	BlastRadii         []SAMPIBlastRadius      `json:"blast_radii,omitempty"`
	ComplianceView     SAMPIComplianceView     `json:"compliance_view,omitempty"`
	ProgramInsights    []SAMPIProgramInsight   `json:"program_insights,omitempty"`
}

// ── SDI — Security Decision Intelligence ──

type SDIIntelligence struct {
	Recommendations      []SDIDecisionRecommendation `json:"recommendations,omitempty"`
	FixSimulations       []SDIFixSimulation          `json:"fix_simulations,omitempty"`
	FailureSimulations   []SDIFailureSimulation      `json:"failure_simulations,omitempty"`
	ControlImpacts       []SDIControlImpact          `json:"control_impacts,omitempty"`
	DecisionTrees        SDIDecisionTreeResult       `json:"decision_trees,omitempty"`
	BoardScenarios       SDIBoardScenarios           `json:"board_scenarios,omitempty"`
	InvestmentPriorities []SDIInvestmentPriority     `json:"investment_priorities,omitempty"`
	AttackPathCollapse   []SDIAttackPathCollapse     `json:"attack_path_collapse,omitempty"`
	ComplianceImpacts    []SDIComplianceImpact       `json:"compliance_impacts,omitempty"`
	RemediationRoadmap   SDIRemediationRoadmap       `json:"remediation_roadmap,omitempty"`
	Dashboard            SDIDecisionDashboard        `json:"dashboard,omitempty"`
	ExecutiveScenarios   SDIExecutiveScenarios       `json:"executive_scenarios,omitempty"`
}

type SDIDecisionRecommendation struct {
	ID                  string   `json:"id,omitempty"`
	Title               string   `json:"title,omitempty"`
	Description         string   `json:"description,omitempty"`
	AffectedFindings    []string `json:"affected_findings,omitempty"`
	AffectedThreats     []string `json:"affected_threats,omitempty"`
	AffectedAttackPaths []string `json:"affected_attack_paths,omitempty"`
	AffectedControls    []string `json:"affected_controls,omitempty"`
	AffectedAssets      []string `json:"affected_assets,omitempty"`
	RiskReduction       string   `json:"risk_reduction,omitempty"`
	Effort              string   `json:"effort,omitempty"`
	Priority            string   `json:"priority,omitempty"`
	BusinessImpact      string   `json:"business_impact,omitempty"`
	ComplianceImpact    string   `json:"compliance_impact,omitempty"`
	Rationale           string   `json:"rationale,omitempty"`
}

type SDIFixSimulation struct {
	ControlName         string  `json:"control_name,omitempty"`
	ControlCategory     string  `json:"control_category,omitempty"`
	OriginalCritical    int     `json:"original_critical,omitempty"`
	OriginalHigh        int     `json:"original_high,omitempty"`
	OriginalTotal       int     `json:"original_total,omitempty"`
	OriginalAttackPaths int     `json:"original_attack_paths,omitempty"`
	OriginalCoverage    float64 `json:"original_coverage,omitempty"`
	NewCritical         int     `json:"new_critical,omitempty"`
	NewHigh             int     `json:"new_high,omitempty"`
	NewTotal            int     `json:"new_total,omitempty"`
	NewAttackPaths      int     `json:"new_attack_paths,omitempty"`
	NewCoverage         float64 `json:"new_coverage,omitempty"`
}

type SDIFailureSimulation struct {
	ControlName       string  `json:"control_name,omitempty"`
	ControlCategory   string  `json:"control_category,omitempty"`
	SystemsImpacted   int     `json:"systems_impacted,omitempty"`
	AttackPathsOpened int     `json:"attack_paths_opened,omitempty"`
	NewFindings       int     `json:"new_findings,omitempty"`
	RiskIncrease      string  `json:"risk_increase,omitempty"`
	RiskScoreIncrease float64 `json:"risk_score_increase,omitempty"`
}

type SDIControlImpact struct {
	ControlName     string `json:"control_name,omitempty"`
	Category        string `json:"category,omitempty"`
	SecurityValue   string `json:"security_value,omitempty"`
	Effort          string `json:"effort,omitempty"`
	ROI             string `json:"roi,omitempty"`
	FindingCount    int    `json:"finding_count,omitempty"`
	ThreatCount     int    `json:"threat_count,omitempty"`
	AttackPathCount int    `json:"attack_path_count,omitempty"`
}

type SDIDecisionTree struct {
	Budget           string                      `json:"budget,omitempty"`
	ActionCount      int                         `json:"action_count,omitempty"`
	RecommendedOrder []SDIDecisionRecommendation `json:"recommended_order,omitempty"`
	Rationale        string                      `json:"rationale,omitempty"`
}

type SDIDecisionTreeResult struct {
	SingleAction SDIDecisionTree `json:"single_action,omitempty"`
	ThreeActions SDIDecisionTree `json:"three_actions,omitempty"`
	FiveActions  SDIDecisionTree `json:"five_actions,omitempty"`
}

type SDIBoardScenario struct {
	Scenario         string   `json:"scenario,omitempty"`
	Description      string   `json:"description,omitempty"`
	RiskScore        float64  `json:"risk_score,omitempty"`
	CriticalFindings int      `json:"critical_findings,omitempty"`
	AttackPaths      int      `json:"attack_paths,omitempty"`
	CoverageRate     float64  `json:"coverage_rate,omitempty"`
	KeyRisks         []string `json:"key_risks,omitempty"`
}

type SDIBoardScenarios struct {
	DoNothing        SDIBoardScenario `json:"do_nothing,omitempty"`
	PartialRemediate SDIBoardScenario `json:"partial_remediate,omitempty"`
	FullRemediate    SDIBoardScenario `json:"full_remediate,omitempty"`
}

type SDIInvestmentPriority struct {
	Area          string  `json:"area,omitempty"`
	Rank          int     `json:"rank,omitempty"`
	Score         float64 `json:"score,omitempty"`
	Rationale     string  `json:"rationale,omitempty"`
	FindingCount  int     `json:"finding_count,omitempty"`
	RiskReduction string  `json:"risk_reduction,omitempty"`
}

type SDIAttackPathCollapse struct {
	ControlName        string  `json:"control_name,omitempty"`
	Category           string  `json:"category,omitempty"`
	AttackPathsReduced int     `json:"attack_paths_reduced,omitempty"`
	TotalAttackPaths   int     `json:"total_attack_paths,omitempty"`
	ReductionPercent   float64 `json:"reduction_percent,omitempty"`
}

type SDIComplianceImpact struct {
	Framework   string `json:"framework,omitempty"`
	Action      string `json:"action,omitempty"`
	Improvement string `json:"improvement,omitempty"`
	Rationale   string `json:"rationale,omitempty"`
}

type SDIRemediationRoadmap struct {
	Phase30  []SDIRoadmapItem `json:"phase_30,omitempty"`
	Phase90  []SDIRoadmapItem `json:"phase_90,omitempty"`
	Phase180 []SDIRoadmapItem `json:"phase_180,omitempty"`
	Phase12m []SDIRoadmapItem `json:"phase_12m,omitempty"`
}

type SDIRoadmapItem struct {
	Action        string `json:"action,omitempty"`
	Category      string `json:"category,omitempty"`
	Priority      string `json:"priority,omitempty"`
	Effort        string `json:"effort,omitempty"`
	RiskReduction string `json:"risk_reduction,omitempty"`
}

type SDIDecisionDashboard struct {
	TopDecisions         []SDIDecisionRecommendation `json:"top_decisions,omitempty"`
	QuickWins            []SDIDecisionRecommendation `json:"quick_wins,omitempty"`
	StrategicActions     []SDIDecisionRecommendation `json:"strategic_actions,omitempty"`
	RiskReductionSummary string                      `json:"risk_reduction_summary,omitempty"`
	TotalRiskReduction   float64                     `json:"total_risk_reduction,omitempty"`
}

type SDIExecutiveScenario struct {
	Scenario          string  `json:"scenario,omitempty"`
	RiskScore         float64 `json:"risk_score,omitempty"`
	FindingsResolved  int     `json:"findings_resolved,omitempty"`
	AttackPathsClosed int     `json:"attack_paths_closed,omitempty"`
	CoverageAchieved  float64 `json:"coverage_achieved,omitempty"`
	Description       string  `json:"description,omitempty"`
}

type SDIExecutiveScenarios struct {
	BestCase   SDIExecutiveScenario `json:"best_case,omitempty"`
	LikelyCase SDIExecutiveScenario `json:"likely_case,omitempty"`
	WorstCase  SDIExecutiveScenario `json:"worst_case,omitempty"`
}

// ── SDT (Security Digital Twin) ──

type SDTIntelligence struct {
	Twin               ArchitectureTwinPR     `json:"twin,omitempty"`
	ChangeImpacts      []ChangeImpactPR       `json:"change_impacts,omitempty"`
	ArchitectureDiffs  []ArchitectureDiffPR   `json:"architecture_diffs,omitempty"`
	EvolutionInsights  []EvolutionInsightPR   `json:"evolution_insights,omitempty"`
	ControlDrifts      []ControlDriftPR       `json:"control_drifts,omitempty"`
	AssumptionDecays   []AssumptionDecayPR    `json:"assumption_decays,omitempty"`
	SecurityDebt       SecurityDebtScorePR    `json:"security_debt,omitempty"`
	ComplianceDrifts   []ComplianceDriftPR    `json:"compliance_drifts,omitempty"`
	AttackSurfaceTrend AttackSurfaceTrendPR   `json:"attack_surface_trend,omitempty"`
	Timeline           ArchitectureTimelinePR `json:"timeline,omitempty"`
	WhatIfScenarios    []WhatIfScenarioPR     `json:"what_if_scenarios,omitempty"`
	MergerAnalysis     MergerAnalysisPR       `json:"merger_analysis,omitempty"`
	ZeroTrust          ZeroTrustAnalysisPR    `json:"zero_trust,omitempty"`
	Resilience         []ResilienceScenarioPR `json:"resilience,omitempty"`
	CrownJewels        []CrownJewelRankingPR  `json:"crown_jewels,omitempty"`
	ExecutiveReport    DigitalTwinReportPR    `json:"executive_report,omitempty"`
	PortfolioSummary   PortfolioTwinSummaryPR `json:"portfolio_summary,omitempty"`
}

type TwinAssetPR struct {
	ID          string `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Type        string `json:"type,omitempty"`
	Criticality string `json:"criticality,omitempty"`
}

type TwinComponentPR struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Type string `json:"type,omitempty"`
	Zone string `json:"zone,omitempty"`
}

type TwinRelationshipPR struct {
	ID           string `json:"id,omitempty"`
	SourceID     string `json:"source_id,omitempty"`
	TargetID     string `json:"target_id,omitempty"`
	RelationType string `json:"relation_type,omitempty"`
	Encrypted    bool   `json:"encrypted,omitempty"`
}

type ArchitectureTwinPR struct {
	ID               string  `json:"id,omitempty"`
	Version          string  `json:"version,omitempty"`
	ArchitectureName string  `json:"architecture_name,omitempty"`
	Domain           string  `json:"domain,omitempty"`
	RiskScore        float64 `json:"risk_score,omitempty"`
	Coverage         float64 `json:"coverage,omitempty"`
	SourceHash       string  `json:"source_hash,omitempty"`
}

type ChangeImpactPR struct {
	Change              string `json:"change,omitempty"`
	ComponentAffected   string `json:"component_affected,omitempty"`
	ImpactType          string `json:"impact_type,omitempty"`
	Severity            string `json:"severity,omitempty"`
	RisksAffected       int    `json:"risks_affected,omitempty"`
	AttackPathsAffected int    `json:"attack_paths_affected,omitempty"`
	ControlsAffected    int    `json:"controls_affected,omitempty"`
	Description         string `json:"description,omitempty"`
}

type ArchitectureDiffPR struct {
	Category       string  `json:"category,omitempty"`
	AddedCount     int     `json:"added_count,omitempty"`
	RemovedCount   int     `json:"removed_count,omitempty"`
	ChangedCount   int     `json:"changed_count,omitempty"`
	RiskScoreDelta float64 `json:"risk_score_delta,omitempty"`
	CoverageDelta  float64 `json:"coverage_delta,omitempty"`
	Description    string  `json:"description,omitempty"`
}

type EvolutionInsightPR struct {
	Scenario   string `json:"scenario,omitempty"`
	Assumption string `json:"assumption,omitempty"`
	Status     string `json:"status,omitempty"`
	Rationale  string `json:"rationale,omitempty"`
}

type ControlDriftPR struct {
	ControlName   string `json:"control_name,omitempty"`
	Category      string `json:"category,omitempty"`
	ExpectedState string `json:"expected_state,omitempty"`
	CurrentState  string `json:"current_state,omitempty"`
	RiskImpact    string `json:"risk_impact,omitempty"`
}

type AssumptionDecayPR struct {
	AssumptionID   string `json:"assumption_id,omitempty"`
	Description    string `json:"description,omitempty"`
	TimeElapsed    string `json:"time_elapsed,omitempty"`
	Status         string `json:"status,omitempty"`
	Recommendation string `json:"recommendation,omitempty"`
}

type SecurityDebtScorePR struct {
	TotalDebt      float64 `json:"total_debt,omitempty"`
	FindingDebt    float64 `json:"finding_debt,omitempty"`
	ControlDebt    float64 `json:"control_debt,omitempty"`
	AssumptionDebt float64 `json:"assumption_debt,omitempty"`
	RiskScore      float64 `json:"risk_score,omitempty"`
}

type ComplianceDriftPR struct {
	Framework      string   `json:"framework,omitempty"`
	Status         string   `json:"status,omitempty"`
	NewGaps        int      `json:"new_gaps,omitempty"`
	ResolvedGaps   int      `json:"resolved_gaps,omitempty"`
	RegressedAreas []string `json:"regressed_areas,omitempty"`
}

type AttackSurfaceTrendPR struct {
	InternetExposure int     `json:"internet_exposure,omitempty"`
	ThirdParties     int     `json:"third_parties,omitempty"`
	IdentitySystems  int     `json:"identity_systems,omitempty"`
	CloudServices    int     `json:"cloud_services,omitempty"`
	AdminPaths       int     `json:"admin_paths,omitempty"`
	GrowthRate       float64 `json:"growth_rate,omitempty"`
}

type ArchitectureTimelinePR struct {
	Trend     string  `json:"trend,omitempty"`
	DeltaRisk float64 `json:"delta_risk,omitempty"`
}

type WhatIfScenarioPR struct {
	Name          string  `json:"name,omitempty"`
	RiskDelta     float64 `json:"risk_delta,omitempty"`
	CoverageDelta float64 `json:"coverage_delta,omitempty"`
	FindingsDelta int     `json:"findings_delta,omitempty"`
	Description   string  `json:"description,omitempty"`
}

type MergerAnalysisPR struct {
	CombinedRiskScore float64  `json:"combined_risk_score,omitempty"`
	InheritedRisks    int      `json:"inherited_risks,omitempty"`
	InheritedControls int      `json:"inherited_controls,omitempty"`
	SharedRisks       []string `json:"shared_risks,omitempty"`
	NewRisks          []string `json:"new_risks,omitempty"`
	GapFindings       []string `json:"gap_findings,omitempty"`
}

type ZeroTrustDimensionPR struct {
	Dimension string  `json:"dimension,omitempty"`
	Score     float64 `json:"score,omitempty"`
	Target    float64 `json:"target,omitempty"`
	Gap       float64 `json:"gap,omitempty"`
	Status    string  `json:"status,omitempty"`
}

type ZeroTrustAnalysisPR struct {
	Dimensions []ZeroTrustDimensionPR `json:"dimensions,omitempty"`
	Overall    float64                `json:"overall,omitempty"`
	Target     float64                `json:"target,omitempty"`
	Gap        float64                `json:"gap,omitempty"`
}

type ResilienceScenarioPR struct {
	FailurePoint        string   `json:"failure_point,omitempty"`
	BusinessImpact      string   `json:"business_impact,omitempty"`
	SecurityImpact      string   `json:"security_impact,omitempty"`
	AffectedAssets      []string `json:"affected_assets,omitempty"`
	AttackPathsOpened   int      `json:"attack_paths_opened,omitempty"`
	RecoveryAssumptions []string `json:"recovery_assumptions,omitempty"`
}

type CrownJewelRankingPR struct {
	AssetName       string  `json:"asset_name,omitempty"`
	BusinessValue   string  `json:"business_value,omitempty"`
	AttackValue     string  `json:"attack_value,omitempty"`
	DependencyCount int     `json:"dependency_count,omitempty"`
	ThreatCount     int     `json:"threat_count,omitempty"`
	BlastRadius     string  `json:"blast_radius,omitempty"`
	OverallScore    float64 `json:"overall_score,omitempty"`
}

type DigitalTwinReportPR struct {
	ArchitectureHealth   string  `json:"architecture_health,omitempty"`
	SecurityDebtScore    float64 `json:"security_debt_score,omitempty"`
	ControlDriftCount    int     `json:"control_drift_count,omitempty"`
	ComplianceDriftCount int     `json:"compliance_drift_count,omitempty"`
	RiskTrend            string  `json:"risk_trend,omitempty"`
	AttackSurfaceTrend   string  `json:"attack_surface_trend,omitempty"`
}

type PortfolioTwinSummaryPR struct {
	ArchitectureCount int      `json:"architecture_count,omitempty"`
	SharedRisks       []string `json:"shared_risks,omitempty"`
	SharedVendors     []string `json:"shared_vendors,omitempty"`
	SharedControls    []string `json:"shared_controls,omitempty"`
	EnterpriseTrends  []string `json:"enterprise_trends,omitempty"`
	AggregatedDebt    float64  `json:"aggregated_debt,omitempty"`
}
