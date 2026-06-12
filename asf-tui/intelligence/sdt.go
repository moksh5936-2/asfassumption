package intelligence

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// ═══════════════════════════════════════════════════════
// SDT — Security Digital Twin Engine (ASF V12)
// Phases 1-17
// ═══════════════════════════════════════════════════════

// ── PHASE 1 — DIGITAL TWIN MODEL ──

type TwinAsset struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Criticality string   `json:"criticality"`
	DataTypes   []string `json:"data_types,omitempty"`
}

type TwinComponent struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Zone     string `json:"zone"`
	AssetRef string `json:"asset_ref,omitempty"`
}

type TwinRelationship struct {
	ID            string `json:"id"`
	SourceID      string `json:"source_id"`
	TargetID      string `json:"target_id"`
	RelationType  string `json:"relation_type"`
	Protocol      string `json:"protocol,omitempty"`
	Encrypted     bool   `json:"encrypted"`
	Authenticated bool   `json:"authenticated"`
}

type ArchitectureTwin struct {
	ID               string             `json:"id"`
	Version          string             `json:"version"`
	Timestamp        time.Time          `json:"timestamp"`
	ArchitectureName string             `json:"architecture_name"`
	Domain           string             `json:"domain"`
	Assets           []TwinAsset        `json:"assets,omitempty"`
	Components       []TwinComponent    `json:"components,omitempty"`
	Relationships    []TwinRelationship `json:"relationships,omitempty"`
	Threats          []Threat           `json:"threats,omitempty"`
	Controls         []SDRIControl      `json:"controls,omitempty"`
	Assumptions      []Assumption       `json:"assumptions,omitempty"`
	Compliance       []string           `json:"compliance,omitempty"`
	AttackPaths      []AttackPath       `json:"attack_paths,omitempty"`
	Findings         []SDRIFinding      `json:"findings,omitempty"`
	RiskScore        float64            `json:"risk_score"`
	Coverage         float64            `json:"coverage"`
	SourceHash       string             `json:"source_hash,omitempty"`
}

// ── PHASE 2 — CHANGE IMPACT ──

type ChangeImpact struct {
	Change              string `json:"change"`
	ComponentAffected   string `json:"component_affected"`
	ImpactType          string `json:"impact_type"`
	Severity            string `json:"severity"`
	RisksAffected       int    `json:"risks_affected"`
	AttackPathsAffected int    `json:"attack_paths_affected"`
	ControlsAffected    int    `json:"controls_affected"`
	Description         string `json:"description"`
}

// ── PHASE 3 — SECURITY DIFF ──

type ArchitectureDiff struct {
	Category       string  `json:"category"`
	AddedCount     int     `json:"added_count"`
	RemovedCount   int     `json:"removed_count"`
	ChangedCount   int     `json:"changed_count"`
	RiskScoreDelta float64 `json:"risk_score_delta"`
	CoverageDelta  float64 `json:"coverage_delta"`
	Description    string  `json:"description"`
}

// ── PHASE 4 — EVOLUTION ANALYSIS ──

type EvolutionInsight struct {
	Scenario   string `json:"scenario"`
	Assumption string `json:"assumption"`
	Status     string `json:"status"`
	Rationale  string `json:"rationale"`
}

// ── PHASE 5 — CONTROL DRIFT ──

type ControlDrift struct {
	ControlName   string `json:"control_name"`
	Category      string `json:"category"`
	ExpectedState string `json:"expected_state"`
	CurrentState  string `json:"current_state"`
	DriftType     string `json:"drift_type"`
	Severity      string `json:"severity"`
}

// ── PHASE 6 — ASSUMPTION DECAY ──

type AssumptionDecay struct {
	AssumptionID  string `json:"assumption_id"`
	Description   string `json:"description"`
	Age           string `json:"age"`
	Status        string `json:"status"`
	EvidenceCount int    `json:"evidence_count"`
}

// ── PHASE 7 — SECURITY DEBT ──

type SecurityDebtScore struct {
	TotalDebt      float64 `json:"total_debt"`
	FindingDebt    float64 `json:"finding_debt"`
	ControlDebt    float64 `json:"control_debt"`
	AssumptionDebt float64 `json:"assumption_debt"`
	RiskScore      float64 `json:"risk_score"`
}

// ── PHASE 8 — COMPLIANCE DRIFT ──

type ComplianceDrift struct {
	Framework      string   `json:"framework"`
	Status         string   `json:"status"`
	NewGaps        int      `json:"new_gaps"`
	ResolvedGaps   int      `json:"resolved_gaps"`
	RegressedAreas []string `json:"regressed_areas,omitempty"`
}

// ── PHASE 9 — ATTACK SURFACE EVOLUTION ──

type AttackSurfaceTrend struct {
	InternetExposure int     `json:"internet_exposure"`
	ThirdParties     int     `json:"third_parties"`
	IdentitySystems  int     `json:"identity_systems"`
	CloudServices    int     `json:"cloud_services"`
	AdminPaths       int     `json:"admin_paths"`
	GrowthRate       float64 `json:"growth_rate"`
}

// ── PHASE 10 — ARCHITECTURE TIMELINE ──

type TimelineEntry struct {
	Version     string    `json:"version"`
	Timestamp   time.Time `json:"timestamp"`
	RiskScore   float64   `json:"risk_score"`
	Coverage    float64   `json:"coverage"`
	Findings    int       `json:"findings"`
	Controls    int       `json:"controls"`
	AttackPaths int       `json:"attack_paths"`
}

type ArchitectureTimeline struct {
	Entries   []TimelineEntry `json:"entries,omitempty"`
	Trend     string          `json:"trend"`
	DeltaRisk float64         `json:"delta_risk"`
}

// ── PHASE 11 — WHAT-IF MODELING ──

type WhatIfScenario struct {
	Name            string  `json:"name"`
	Description     string  `json:"description"`
	RiskDelta       float64 `json:"risk_delta"`
	ThreatDelta     int     `json:"threat_delta"`
	AttackPathDelta int     `json:"attack_path_delta"`
	ComplianceDelta string  `json:"compliance_delta"`
	CoverageDelta   float64 `json:"coverage_delta"`
}

// ── PHASE 12 — MERGER & ACQUISITION ──

type MergerAnalysis struct {
	CombinedRiskScore float64  `json:"combined_risk_score"`
	InheritedRisks    int      `json:"inherited_risks"`
	InheritedVendors  []string `json:"inherited_vendors,omitempty"`
	InheritedControls int      `json:"inherited_controls"`
	ComplianceGaps    []string `json:"compliance_gaps,omitempty"`
	SharedRisks       []string `json:"shared_risks,omitempty"`
}

// ── PHASE 13 — ZERO TRUST EVOLUTION ──

type ZeroTrustDimension struct {
	Dimension    string  `json:"dimension"`
	CurrentScore float64 `json:"current_score"`
	TargetScore  float64 `json:"target_score"`
	Gap          float64 `json:"gap"`
	Progress     string  `json:"progress"`
}

type ZeroTrustAnalysis struct {
	Dimensions []ZeroTrustDimension `json:"dimensions,omitempty"`
	Overall    float64              `json:"overall"`
	Target     float64              `json:"target"`
	Gap        float64              `json:"gap"`
}

// ── PHASE 14 — RESILIENCE MODELING ──

type ResilienceScenario struct {
	FailurePoint        string   `json:"failure_point"`
	BusinessImpact      string   `json:"business_impact"`
	SecurityImpact      string   `json:"security_impact"`
	AffectedAssets      []string `json:"affected_assets,omitempty"`
	AttackPathsOpened   int      `json:"attack_paths_opened"`
	RecoveryAssumptions []string `json:"recovery_assumptions,omitempty"`
}

// ── PHASE 15 — CROWN JEWEL ANALYSIS ──

type CrownJewelRanking struct {
	AssetName       string  `json:"asset_name"`
	BusinessValue   string  `json:"business_value"`
	AttackValue     string  `json:"attack_value"`
	DependencyCount int     `json:"dependency_count"`
	ThreatCount     int     `json:"threat_count"`
	BlastRadius     string  `json:"blast_radius"`
	OverallScore    float64 `json:"overall_score"`
}

// ── PHASE 16 — EXECUTIVE DIGITAL TWIN REPORT ──

type DigitalTwinReport struct {
	ArchitectureHealth   string               `json:"architecture_health"`
	SecurityDebtScore    float64              `json:"security_debt_score"`
	ControlDriftCount    int                  `json:"control_drift_count"`
	ComplianceDriftCount int                  `json:"compliance_drift_count"`
	RiskTrend            string               `json:"risk_trend"`
	AttackSurfaceTrend   string               `json:"attack_surface_trend"`
	Sections             []DigitalTwinSection `json:"sections,omitempty"`
}

type DigitalTwinSection struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

// ── PHASE 17 — PORTFOLIO DIGITAL TWIN ──

type PortfolioTwinSummary struct {
	ArchitectureCount int      `json:"architecture_count"`
	SharedRisks       []string `json:"shared_risks,omitempty"`
	SharedVendors     []string `json:"shared_vendors,omitempty"`
	SharedControls    []string `json:"shared_controls,omitempty"`
	EnterpriseTrends  []string `json:"enterprise_trends,omitempty"`
	AggregatedDebt    float64  `json:"aggregated_debt"`
}

// ── SDT INPUT / RESULT ──

type SDTInput struct {
	ArchitectureName string
	Domain           string
	Assets           []TwinAsset
	Components       []TwinComponent
	Relationships    []TwinRelationship
	Threats          []Threat
	Controls         []SDRIControl
	Assumptions      []Assumption
	Compliance       []string
	AttackPaths      []AttackPath
	Findings         []SDRIFinding
	RiskScore        float64
	Coverage         float64
	PreviousTwins    []ArchitectureTwin
	PortfolioTwins   []ArchitectureTwin
}

type SDTResult struct {
	Twin               ArchitectureTwin     `json:"twin"`
	ChangeImpacts      []ChangeImpact       `json:"change_impacts,omitempty"`
	ArchitectureDiffs  []ArchitectureDiff   `json:"architecture_diffs,omitempty"`
	EvolutionInsights  []EvolutionInsight   `json:"evolution_insights,omitempty"`
	ControlDrifts      []ControlDrift       `json:"control_drifts,omitempty"`
	AssumptionDecays   []AssumptionDecay    `json:"assumption_decays,omitempty"`
	SecurityDebt       SecurityDebtScore    `json:"security_debt"`
	ComplianceDrifts   []ComplianceDrift    `json:"compliance_drifts,omitempty"`
	AttackSurfaceTrend AttackSurfaceTrend   `json:"attack_surface_trend"`
	Timeline           ArchitectureTimeline `json:"timeline"`
	WhatIfScenarios    []WhatIfScenario     `json:"what_if_scenarios,omitempty"`
	MergerAnalysis     MergerAnalysis       `json:"merger_analysis"`
	ZeroTrust          ZeroTrustAnalysis    `json:"zero_trust"`
	Resilience         []ResilienceScenario `json:"resilience,omitempty"`
	CrownJewels        []CrownJewelRanking  `json:"crown_jewels,omitempty"`
	ExecutiveReport    DigitalTwinReport    `json:"executive_report"`
	PortfolioSummary   PortfolioTwinSummary `json:"portfolio_summary"`
}

// ── SDT ENGINE ──

type SDTEngine struct{}

func NewSDTEngine() *SDTEngine {
	return &SDTEngine{}
}

func (e *SDTEngine) Run(input SDTInput) *SDTResult {
	r := &SDTResult{}

	r.Twin = e.buildTwin(input)
	r.ChangeImpacts = e.analyzeChangeImpacts(input)
	r.ArchitectureDiffs = e.computeArchitectureDiffs(input)
	r.EvolutionInsights = e.analyzeEvolution(input)
	r.ControlDrifts = e.detectControlDrift(input)
	r.AssumptionDecays = e.analyzeAssumptionDecay(input)
	r.SecurityDebt = e.calculateSecurityDebt(input)
	r.ComplianceDrifts = e.analyzeComplianceDrift(input)
	r.AttackSurfaceTrend = e.measureAttackSurface(input)
	r.Timeline = e.buildTimeline(input)
	r.WhatIfScenarios = e.generateWhatIfScenarios(input)
	r.MergerAnalysis = e.analyzeMerger(input)
	r.ZeroTrust = e.measureZeroTrust(input)
	r.Resilience = e.modelResilience(input)
	r.CrownJewels = e.identifyCrownJewels(input)
	r.ExecutiveReport = e.generateExecutiveReport(r)
	r.PortfolioSummary = e.summarizePortfolio(input)

	return r
}

// ── BUILD TWIN ──

func (e *SDTEngine) buildTwin(input SDTInput) ArchitectureTwin {
	assets := input.Assets
	if len(assets) == 0 {
		assets = e.deriveAssets(input)
	}
	components := input.Components
	if len(components) == 0 {
		components = e.deriveComponents(input)
	}
	relationships := input.Relationships
	if len(relationships) == 0 {
		relationships = e.deriveRelationships(input)
	}

	twin := ArchitectureTwin{
		ID:               fmt.Sprintf("twin-%s", input.ArchitectureName),
		Version:          "1.0",
		Timestamp:        time.Now(),
		ArchitectureName: input.ArchitectureName,
		Domain:           input.Domain,
		Assets:           assets,
		Components:       components,
		Relationships:    relationships,
		Threats:          input.Threats,
		Controls:         input.Controls,
		Assumptions:      input.Assumptions,
		Compliance:       input.Compliance,
		AttackPaths:      input.AttackPaths,
		Findings:         input.Findings,
		RiskScore:        input.RiskScore,
		Coverage:         input.Coverage,
	}
	return twin
}

func (e *SDTEngine) deriveAssets(input SDTInput) []TwinAsset {
	seen := map[string]bool{}
	var assets []TwinAsset
	for _, f := range input.Findings {
		for _, comp := range f.AffectedComponents {
			if !seen[comp] {
				seen[comp] = true
				assets = append(assets, TwinAsset{
					ID:   fmt.Sprintf("asset-%d", len(assets)+1),
					Name: comp,
					Type: classifyComponent(comp),
				})
			}
		}
	}
	for _, t := range input.Threats {
		for _, a := range t.AffectedAssets {
			if !seen[a] {
				seen[a] = true
				assets = append(assets, TwinAsset{
					ID:   fmt.Sprintf("asset-%d", len(assets)+1),
					Name: a,
					Type: classifyComponent(a),
				})
			}
		}
		for _, comp := range t.AffectedComponents {
			if !seen[comp] {
				seen[comp] = true
				assets = append(assets, TwinAsset{
					ID:   fmt.Sprintf("asset-%d", len(assets)+1),
					Name: comp,
					Type: classifyComponent(comp),
				})
			}
		}
	}
	for _, p := range input.AttackPaths {
		if p.EntryPoint != "" && !seen[p.EntryPoint] {
			seen[p.EntryPoint] = true
			assets = append(assets, TwinAsset{
				ID: fmt.Sprintf("asset-%d", len(assets)+1), Name: p.EntryPoint, Type: "EntryPoint",
			})
		}
		if p.TargetAsset != "" && !seen[p.TargetAsset] {
			seen[p.TargetAsset] = true
			assets = append(assets, TwinAsset{
				ID: fmt.Sprintf("asset-%d", len(assets)+1), Name: p.TargetAsset, Type: "Target",
			})
		}
	}
	if len(assets) == 0 {
		assets = append(assets, TwinAsset{
			ID: "asset-1", Name: input.ArchitectureName, Type: "Architecture",
		})
	}
	return assets
}

func (e *SDTEngine) deriveComponents(input SDTInput) []TwinComponent {
	var comps []TwinComponent
	for i, a := range e.deriveAssets(input) {
		comps = append(comps, TwinComponent{
			ID: fmt.Sprintf("comp-%d", i+1), Name: a.Name,
			Type: a.Type, AssetRef: a.ID,
		})
	}
	return comps
}

func (e *SDTEngine) deriveRelationships(input SDTInput) []TwinRelationship {
	var rels []TwinRelationship
	comps := e.deriveComponents(input)
	if len(comps) >= 2 {
		for i := 0; i < len(comps)-1; i++ {
			rels = append(rels, TwinRelationship{
				ID: fmt.Sprintf("rel-%d", i+1), SourceID: comps[i].ID,
				TargetID: comps[i+1].ID, RelationType: "connects",
			})
		}
	}
	for _, p := range input.AttackPaths {
		if len(p.AffectedComponents) >= 2 {
			for i := 0; i < len(p.AffectedComponents)-1; i++ {
				srcID := findComponentID(comps, p.AffectedComponents[i])
				tgtID := findComponentID(comps, p.AffectedComponents[i+1])
				if srcID != "" && tgtID != "" {
					rels = append(rels, TwinRelationship{
						ID:       fmt.Sprintf("rel-ap-%d", len(rels)+1),
						SourceID: srcID, TargetID: tgtID,
						RelationType: "attacks",
					})
				}
			}
		}
	}
	return rels
}

// ── PHASE 2 — CHANGE IMPACT ENGINE ──

func (e *SDTEngine) analyzeChangeImpacts(input SDTInput) []ChangeImpact {
	twin := e.buildTwin(input)
	var impacts []ChangeImpact

	if hasComponent(twin.Components, "auth") || hasComponent(twin.Components, "identity") {
		impacts = append(impacts, ChangeImpact{
			Change: "Replace Identity Provider", ComponentAffected: "Authentication",
			ImpactType: "High", Severity: "High",
			RisksAffected:       countThreatsByCategory(input.Threats, "Authentication"),
			AttackPathsAffected: len(input.AttackPaths), ControlsAffected: len(input.Controls),
			Description: "Replacing the identity provider impacts authentication, session management, and MFA flows.",
		})
	}
	if hasComponent(twin.Components, "database") || hasComponent(twin.Components, "data") {
		impacts = append(impacts, ChangeImpact{
			Change: "Replace Database Platform", ComponentAffected: "Data Storage",
			ImpactType: "High", Severity: "High",
			RisksAffected:       countSeverity(input.Findings, "Critical") + countSeverity(input.Findings, "High"),
			AttackPathsAffected: len(input.AttackPaths), ControlsAffected: countControlsByCategory(input.Controls, "DataProtection"),
			Description: "Database migration impacts encryption at rest, access controls, and backup procedures.",
		})
	}
	if hasComponent(twin.Components, "api") || hasComponent(twin.Components, "gateway") {
		impacts = append(impacts, ChangeImpact{
			Change: "Modify API Gateway", ComponentAffected: "API Layer",
			ImpactType: "Medium", Severity: "Medium",
			RisksAffected: len(input.Threats), AttackPathsAffected: len(input.AttackPaths),
			ControlsAffected: len(input.Controls),
			Description:      "API gateway changes affect authentication, rate limiting, and request validation.",
		})
	}
	if len(input.Compliance) > 0 {
		impacts = append(impacts, ChangeImpact{
			Change: "Update Compliance Framework", ComponentAffected: "Compliance",
			ImpactType: "Medium", Severity: "Medium",
			RisksAffected: len(input.Findings), AttackPathsAffected: 0, ControlsAffected: len(input.Controls),
			Description: fmt.Sprintf("Compliance changes for %s may require new controls and evidence collection.", strings.Join(input.Compliance, ", ")),
		})
	}
	return impacts
}

// ── PHASE 3 — SECURITY DIFF ENGINE ──

func (e *SDTEngine) computeArchitectureDiffs(input SDTInput) []ArchitectureDiff {
	if len(input.PreviousTwins) == 0 {
		return []ArchitectureDiff{
			{Category: "Baseline", Description: "No previous twin available. Current state is the baseline.", RiskScoreDelta: 0, CoverageDelta: 0},
		}
	}
	prev := input.PreviousTwins[len(input.PreviousTwins)-1]
	diffs := []ArchitectureDiff{
		{
			Category:       "Threats",
			AddedCount:     lenDiff(len(input.Threats), len(prev.Threats), true),
			RemovedCount:   lenDiff(len(input.Threats), len(prev.Threats), false),
			RiskScoreDelta: input.RiskScore - prev.RiskScore,
			Description:    fmt.Sprintf("Threat count changed from %d to %d", len(prev.Threats), len(input.Threats)),
		},
		{
			Category:      "Controls",
			AddedCount:    lenDiff(len(input.Controls), len(prev.Controls), true),
			RemovedCount:  lenDiff(len(input.Controls), len(prev.Controls), false),
			CoverageDelta: input.Coverage - prev.Coverage,
			Description:   fmt.Sprintf("Control count changed from %d to %d", len(prev.Controls), len(input.Controls)),
		},
		{
			Category:       "Attack Paths",
			AddedCount:     lenDiff(len(input.AttackPaths), len(prev.AttackPaths), true),
			RemovedCount:   lenDiff(len(input.AttackPaths), len(prev.AttackPaths), false),
			RiskScoreDelta: input.RiskScore - prev.RiskScore,
			Description:    fmt.Sprintf("Attack path count changed from %d to %d", len(prev.AttackPaths), len(input.AttackPaths)),
		},
		{
			Category:     "Findings",
			AddedCount:   lenDiff(len(input.Findings), len(prev.Findings), true),
			RemovedCount: lenDiff(len(input.Findings), len(prev.Findings), false),
			Description:  fmt.Sprintf("Finding count changed from %d to %d", len(prev.Findings), len(input.Findings)),
		},
	}
	return diffs
}

// ── PHASE 4 — ARCHITECTURE EVOLUTION ANALYSIS ──

func (e *SDTEngine) analyzeEvolution(input SDTInput) []EvolutionInsight {
	var insights []EvolutionInsight
	twin := e.buildTwin(input)

	if len(twin.Components) > 0 {
		insights = append(insights, EvolutionInsight{
			Scenario: "10x Growth in Services", Assumption: "Current trust boundaries scale linearly",
			Status: "Needs Review", Rationale: fmt.Sprintf("With %d components, a 10x increase to %d may invalidate network segmentation and identity assumptions.", len(twin.Components), len(twin.Components)*10),
		})
	}
	if len(twin.Relationships) > 5 {
		insights = append(insights, EvolutionInsight{
			Scenario: "Service Count Doubles", Assumption: "Current authentication model supports growth",
			Status: "Likely Invalid", Rationale: fmt.Sprintf("%d relationships suggest complex interdependencies that may not scale linearly.", len(twin.Relationships)),
		})
	}
	if hasComponent(twin.Components, "vendor") || hasComponent(twin.Components, "third") {
		insights = append(insights, EvolutionInsight{
			Scenario: "Vendor Count Increases", Assumption: "Vendor risk remains constant",
			Status: "Needs Review", Rationale: "Additional vendors introduce supply chain risk and require expanded monitoring capabilities.",
		})
	}
	if len(input.Compliance) > 0 {
		insights = append(insights, EvolutionInsight{
			Scenario: "New Region Expansion", Assumption: "Compliance requirements are uniform across regions",
			Status: "Needs Review", Rationale: fmt.Sprintf("Expanding to new regions may introduce additional compliance requirements beyond %s.", strings.Join(input.Compliance, ", ")),
		})
	}
	if len(insights) == 0 {
		insights = append(insights, EvolutionInsight{
			Scenario: "Baseline", Assumption: "Current architecture stable",
			Status: "Valid", Rationale: "No growth indicators detected at this time.",
		})
	}
	return insights
}

// ── PHASE 5 — CONTROL DRIFT DETECTION ──

func (e *SDTEngine) detectControlDrift(input SDTInput) []ControlDrift {
	var drifts []ControlDrift

	expectedControls := []struct {
		name     string
		category string
		keywords []string
	}{
		{"Multi-Factor Authentication", "AccessControl", []string{"mfa", "auth", "login", "sso", "identity"}},
		{"Encryption at Rest", "DataProtection", []string{"encrypt", "data", "database", "storage"}},
		{"Audit Logging", "Logging", []string{"log", "audit", "monitor"}},
		{"Network Segmentation", "NetworkSecurity", []string{"network", "segment", "firewall"}},
		{"Access Control", "Authorization", []string{"rbac", "access", "permission", "role"}},
		{"TLS Encryption", "DataProtection", []string{"tls", "ssl", "https", "encrypt"}},
		{"Secrets Management", "SecretsManagement", []string{"secret", "vault", "credential", "key"}},
		{"Incident Response", "Governance", []string{"incident", "response", "recovery"}},
		{"Vulnerability Management", "VulnerabilityManagement", []string{"vuln", "patch", "cve", "scan"}},
		{"Backup and Recovery", "Resilience", []string{"backup", "recovery", "disaster"}},
	}

	for _, ec := range expectedControls {
		found := false
		for _, c := range input.Controls {
			if strings.Contains(strings.ToLower(c.Name), ec.keywords[0]) {
				found = true
				break
			}
		}
		if !found {
			for _, f := range input.Findings {
				lower := strings.ToLower(f.Title + " " + f.Description)
				for _, kw := range ec.keywords {
					if strings.Contains(lower, kw) {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
		}
		if !found {
			for _, a := range input.Assumptions {
				lower := strings.ToLower(a.Description)
				for _, kw := range ec.keywords {
					if strings.Contains(lower, kw) {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
		}

		state := "Present"
		driftType := "None"
		severity := "Low"
		if !found {
			state = "Missing"
			driftType = "Control Absent"
			severity = "Medium"
			for _, f := range input.Findings {
				if strings.EqualFold(f.Severity, "Critical") {
					for _, kw := range ec.keywords {
						if strings.Contains(strings.ToLower(f.Title), kw) {
							severity = "Critical"
							break
						}
					}
				}
			}
		}

		drifts = append(drifts, ControlDrift{
			ControlName: ec.name, Category: ec.category,
			ExpectedState: "Present", CurrentState: state,
			DriftType: driftType, Severity: severity,
		})
	}
	return drifts
}

// ── PHASE 6 — ASSUMPTION DECAY ENGINE ──

func (e *SDTEngine) analyzeAssumptionDecay(input SDTInput) []AssumptionDecay {
	var decays []AssumptionDecay
	for _, a := range input.Assumptions {
		age := "Unknown"
		status := "Valid"
		evidenceCount := len(a.EvidenceSources)

		if evidenceCount == 0 {
			status = "Needs Review"
		}
		if a.VerificationStatus == "CONTRADICTED" {
			status = "Likely Invalid"
		}

		if a.SourceSection != "" {
			age = "Current"
		}

		decays = append(decays, AssumptionDecay{
			AssumptionID:  a.ID,
			Description:   truncate(a.Description, 80),
			Age:           age,
			Status:        status,
			EvidenceCount: evidenceCount,
		})
	}
	if len(decays) == 0 {
		decays = append(decays, AssumptionDecay{
			AssumptionID: "N/A", Description: "No assumptions recorded",
			Age: "N/A", Status: "Valid", EvidenceCount: 0,
		})
	}
	return decays
}

// ── PHASE 7 — SECURITY DEBT ENGINE ──

func (e *SDTEngine) calculateSecurityDebt(input SDTInput) SecurityDebtScore {
	findingDebt := 0.0
	for _, f := range input.Findings {
		switch strings.ToLower(f.Severity) {
		case "critical":
			findingDebt += 8.0
		case "high":
			findingDebt += 5.0
		case "medium":
			findingDebt += 3.0
		default:
			findingDebt += 1.0
		}
	}

	controlDebt := 0.0
	expectedCount := 10
	actualCount := len(input.Controls)
	if actualCount < expectedCount {
		controlDebt = float64(expectedCount-actualCount) * 4.0
	}

	assumptionDebt := 0.0
	for _, a := range input.Assumptions {
		if a.VerificationStatus == "CONTRADICTED" {
			assumptionDebt += 5.0
		} else if len(a.EvidenceSources) == 0 {
			assumptionDebt += 2.0
		}
	}

	totalDebt := findingDebt + controlDebt + assumptionDebt
	riskScore := totalDebt / 20.0
	if riskScore > 10 {
		riskScore = 10
	}

	return SecurityDebtScore{
		TotalDebt:      totalDebt,
		FindingDebt:    findingDebt,
		ControlDebt:    controlDebt,
		AssumptionDebt: assumptionDebt,
		RiskScore:      riskScore,
	}
}

// ── PHASE 8 — COMPLIANCE DRIFT ──

func (e *SDTEngine) analyzeComplianceDrift(input SDTInput) []ComplianceDrift {
	var drifts []ComplianceDrift

	if len(input.PreviousTwins) == 0 {
		for _, fw := range input.Compliance {
			drifts = append(drifts, ComplianceDrift{
				Framework: fw, Status: "Baseline",
				NewGaps: 0, ResolvedGaps: 0,
			})
		}
		if len(drifts) == 0 {
			drifts = append(drifts, ComplianceDrift{
				Framework: "N/A", Status: "None",
				NewGaps: 0, ResolvedGaps: 0,
			})
		}
		return drifts
	}

	prev := input.PreviousTwins[len(input.PreviousTwins)-1]
	prevCompliance := map[string]bool{}
	for _, fw := range prev.Compliance {
		prevCompliance[fw] = true
	}
	currentCompliance := map[string]bool{}
	for _, fw := range input.Compliance {
		currentCompliance[fw] = true
	}

	allFrameworks := map[string]bool{}
	for fw := range prevCompliance {
		allFrameworks[fw] = true
	}
	for fw := range currentCompliance {
		allFrameworks[fw] = true
	}

	for fw := range allFrameworks {
		had := prevCompliance[fw]
		has := currentCompliance[fw]
		drift := ComplianceDrift{Framework: fw}

		if had && !has {
			drift.Status = "Regressed"
			drift.RegressedAreas = []string{"Compliance coverage lost"}
		} else if !had && has {
			drift.Status = "New"
			drift.NewGaps = len(input.Findings)
		} else {
			drift.Status = "Stable"
		}
		drifts = append(drifts, drift)
	}

	return drifts
}

// ── PHASE 9 — ATTACK SURFACE EVOLUTION ──

func (e *SDTEngine) measureAttackSurface(input SDTInput) AttackSurfaceTrend {
	twin := e.buildTwin(input)
	trend := AttackSurfaceTrend{}

	for _, comp := range twin.Components {
		lower := strings.ToLower(comp.Name)
		if containsAny(lower, []string{"internet", "external", "public", "dmz", "web"}) {
			trend.InternetExposure++
		}
		if containsAny(lower, []string{"vendor", "partner", "third", "supply"}) {
			trend.ThirdParties++
		}
		if containsAny(lower, []string{"auth", "identity", "sso", "iam", "login"}) {
			trend.IdentitySystems++
		}
		if containsAny(lower, []string{"cloud", "aws", "azure", "gcp"}) {
			trend.CloudServices++
		}
		if containsAny(lower, []string{"admin", "manage", "dashboard", "console"}) {
			trend.AdminPaths++
		}
	}

	if len(input.PreviousTwins) > 0 {
		prev := input.PreviousTwins[len(input.PreviousTwins)-1]
		prevCount := len(prev.Components)
		if prevCount > 0 {
			trend.GrowthRate = float64(len(twin.Components)-prevCount) / float64(prevCount) * 100
		}
	}

	return trend
}

// ── PHASE 10 — ARCHITECTURE TIMELINE ──

func (e *SDTEngine) buildTimeline(input SDTInput) ArchitectureTimeline {
	var entries []TimelineEntry

	current := TimelineEntry{
		Version: "1.0", Timestamp: time.Now(),
		RiskScore: input.RiskScore, Coverage: input.Coverage,
		Findings: len(input.Findings), Controls: len(input.Controls),
		AttackPaths: len(input.AttackPaths),
	}
	entries = append(entries, current)

	for i, prev := range input.PreviousTwins {
		entries = append(entries, TimelineEntry{
			Version: fmt.Sprintf("0.%d", i+1), Timestamp: prev.Timestamp,
			RiskScore: prev.RiskScore, Coverage: prev.Coverage,
			Findings: len(prev.Findings), Controls: len(prev.Controls),
			AttackPaths: len(prev.AttackPaths),
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.Before(entries[j].Timestamp)
	})

	trend := "Stable"
	deltaRisk := 0.0
	if len(entries) >= 2 {
		first := entries[0]
		last := entries[len(entries)-1]
		deltaRisk = last.RiskScore - first.RiskScore
		if deltaRisk > 1.0 {
			trend = "Worsening"
		} else if deltaRisk < -1.0 {
			trend = "Improving"
		}
	}

	return ArchitectureTimeline{
		Entries:   entries,
		Trend:     trend,
		DeltaRisk: deltaRisk,
	}
}

// ── PHASE 11 — WHAT-IF MODELING ──

func (e *SDTEngine) generateWhatIfScenarios(input SDTInput) []WhatIfScenario {
	var scenarios []WhatIfScenario
	twin := e.buildTwin(input)

	// Scenario: Remove MFA
	mfaThreatDelta := 0
	for _, t := range input.Threats {
		lower := strings.ToLower(string(t.Category) + " " + t.Name + " " + t.Description)
		if containsAny(lower, []string{"auth", "credential", "session", "mfa"}) {
			mfaThreatDelta++
		}
	}
	scenarios = append(scenarios, WhatIfScenario{
		Name: "Remove MFA", Description: "If multi-factor authentication is removed from all systems",
		RiskDelta:   2.5 + float64(mfaThreatDelta)*0.5,
		ThreatDelta: mfaThreatDelta, AttackPathDelta: len(twin.AttackPaths) / 2,
		ComplianceDelta: "Negative — PCI DSS, HIPAA, SOC 2 non-compliant",
		CoverageDelta:   -15.0,
	})

	// Scenario: Replace Vendor
	vendorThreats := 0
	for _, t := range input.Threats {
		lower := strings.ToLower(string(t.Category) + " " + t.Name + " " + t.Description)
		if containsAny(lower, []string{"vendor", "supply", "third"}) {
			vendorThreats++
		}
	}
	scenarios = append(scenarios, WhatIfScenario{
		Name: "Replace Vendor", Description: "If a critical vendor is replaced",
		RiskDelta:   1.5 + float64(vendorThreats)*0.3,
		ThreatDelta: vendorThreats, AttackPathDelta: len(twin.AttackPaths) / 3,
		ComplianceDelta: "Neutral — requires re-certification",
		CoverageDelta:   -5.0,
	})

	// Scenario: Add New Region
	scenarios = append(scenarios, WhatIfScenario{
		Name: "Add New Region", Description: "If architecture expands to a new geographic region",
		RiskDelta: 2.0, ThreatDelta: len(input.Threats) / 2,
		AttackPathDelta: len(twin.AttackPaths) / 2,
		ComplianceDelta: fmt.Sprintf("Requires additional compliance for %s", strings.Join(input.Compliance, ", ")),
		CoverageDelta:   -10.0,
	})

	// Scenario: Add New Cloud
	scenarios = append(scenarios, WhatIfScenario{
		Name: "Add New Cloud Provider", Description: "If a new cloud provider is added",
		RiskDelta: 1.5, ThreatDelta: len(input.Threats) / 3,
		AttackPathDelta: len(twin.AttackPaths) / 3,
		ComplianceDelta: "Neutral — requires cloud-specific controls",
		CoverageDelta:   -8.0,
	})

	// Scenario: Merge Networks
	netRisks := 0
	for _, f := range input.Findings {
		if containsAny(strings.ToLower(f.Category), []string{"network", "segment"}) {
			netRisks++
		}
	}
	scenarios = append(scenarios, WhatIfScenario{
		Name: "Merge Networks", Description: "If segmented networks are merged into a flat network",
		RiskDelta:   3.0 + float64(netRisks)*0.5,
		ThreatDelta: len(input.Threats) / 2, AttackPathDelta: len(twin.AttackPaths),
		ComplianceDelta: "Negative — PCI DSS network segmentation requirement violated",
		CoverageDelta:   -20.0,
	})

	return scenarios
}

// ── PHASE 12 — MERGER & ACQUISITION ANALYSIS ──

func (e *SDTEngine) analyzeMerger(input SDTInput) MergerAnalysis {
	analysis := MergerAnalysis{CombinedRiskScore: input.RiskScore}

	if len(input.PreviousTwins) > 0 {
		for _, pt := range input.PreviousTwins {
			analysis.InheritedRisks += len(pt.Findings) + len(pt.Threats)
			analysis.InheritedControls += len(pt.Controls)
			for _, fw := range pt.Compliance {
				has := false
				for _, cfw := range input.Compliance {
					if cfw == fw {
						has = true
						break
					}
				}
				if !has {
					analysis.ComplianceGaps = append(analysis.ComplianceGaps, fw)
				}
			}
			for _, c := range pt.Controls {
				dup := false
				for _, cc := range input.Controls {
					if cc.Name == c.Name {
						dup = true
						break
					}
				}
				if dup {
					analysis.SharedRisks = append(analysis.SharedRisks, c.Name)
				}
			}
			analysis.CombinedRiskScore += pt.RiskScore
		}
		analysis.CombinedRiskScore /= float64(1 + len(input.PreviousTwins))
	}

	return analysis
}

// ── PHASE 13 — ZERO TRUST EVOLUTION ──

func (e *SDTEngine) measureZeroTrust(input SDTInput) ZeroTrustAnalysis {
	analysis := ZeroTrustAnalysis{
		Dimensions: []ZeroTrustDimension{
			e.ztDimension("Identity", input, []string{"mfa", "auth", "sso", "iam", "identity"}, 8.0),
			e.ztDimension("Devices", input, []string{"edr", "endpoint", "device", "certificate"}, 7.0),
			e.ztDimension("Network", input, []string{"segment", "firewall", "micro", "network", "vpn"}, 8.0),
			e.ztDimension("Applications", input, []string{"app", "api", "service", "workload"}, 7.0),
			e.ztDimension("Data", input, []string{"encrypt", "dlp", "tokenize", "classify"}, 8.0),
		},
		Target: 10.0,
	}

	var total float64
	for _, d := range analysis.Dimensions {
		total += d.CurrentScore
	}
	analysis.Overall = total / float64(len(analysis.Dimensions))
	analysis.Gap = analysis.Target - analysis.Overall

	return analysis
}

func (e *SDTEngine) ztDimension(name string, input SDTInput, keywords []string, maxScore float64) ZeroTrustDimension {
	score := 0.0
	increment := maxScore / 4.0

	for _, c := range input.Controls {
		lower := strings.ToLower(c.Name + " " + c.Category + " " + c.Description)
		for _, kw := range keywords {
			if strings.Contains(lower, kw) {
				score += increment * 0.5
				break
			}
		}
	}
	for _, f := range input.Findings {
		lower := strings.ToLower(f.Title + " " + f.Category + " " + f.Description)
		for _, kw := range keywords {
			if strings.Contains(lower, kw) && strings.EqualFold(f.Severity, "Critical") {
				score -= increment * 0.3
				break
			}
		}
	}
	for _, a := range input.Assumptions {
		lower := strings.ToLower(a.Description)
		for _, kw := range keywords {
			if strings.Contains(lower, kw) {
				score += increment * 0.2
				break
			}
		}
	}

	if score < 0 {
		score = 0
	}
	if score > maxScore {
		score = maxScore
	}

	gap := maxScore - score
	progress := "Not Started"
	if gap < maxScore*0.25 {
		progress = "Advanced"
	} else if gap < maxScore*0.5 {
		progress = "In Progress"
	} else if gap < maxScore*0.75 {
		progress = "Early"
	}

	return ZeroTrustDimension{
		Dimension: name, CurrentScore: score,
		TargetScore: maxScore, Gap: gap, Progress: progress,
	}
}

// ── PHASE 14 — RESILIENCE MODELING ──

func (e *SDTEngine) modelResilience(input SDTInput) []ResilienceScenario {
	var scenarios []ResilienceScenario
	twin := e.buildTwin(input)

	vendorAssets := findAssetsByKeyword(twin.Assets, []string{"vendor", "third", "partner", "supply"})
	scenarios = append(scenarios, ResilienceScenario{
		FailurePoint: "Vendor Failure", BusinessImpact: "High",
		SecurityImpact:      "High — loss of vendor-provided security controls",
		AffectedAssets:      assetNames(vendorAssets),
		AttackPathsOpened:   len(twin.AttackPaths) / 3,
		RecoveryAssumptions: []string{"Backup vendor available", "Manual failover procedures documented"},
	})

	identityAssets := findAssetsByKeyword(twin.Assets, []string{"auth", "identity", "sso", "iam", "login"})
	scenarios = append(scenarios, ResilienceScenario{
		FailurePoint: "Identity Provider Failure", BusinessImpact: "Critical",
		SecurityImpact:      "Critical — all authentication and authorization unavailable",
		AffectedAssets:      assetNames(identityAssets),
		AttackPathsOpened:   len(twin.AttackPaths) / 2,
		RecoveryAssumptions: []string{"Local admin accounts available", "Offline authentication fallback configured"},
	})

	dbAssets := findAssetsByKeyword(twin.Assets, []string{"database", "data", "storage", "warehouse"})
	scenarios = append(scenarios, ResilienceScenario{
		FailurePoint: "Database Failure", BusinessImpact: "Critical",
		SecurityImpact:      "High — data unavailable, potential data loss",
		AffectedAssets:      assetNames(dbAssets),
		AttackPathsOpened:   len(twin.AttackPaths) / 4,
		RecoveryAssumptions: []string{"Backups available and tested", "Disaster recovery procedures documented"},
	})

	kmsAssets := findAssetsByKeyword(twin.Assets, []string{"kms", "key", "vault", "hsm", "certificate"})
	scenarios = append(scenarios, ResilienceScenario{
		FailurePoint: "KMS Failure", BusinessImpact: "High",
		SecurityImpact:      "High — unable to decrypt data, access to encrypted systems lost",
		AffectedAssets:      assetNames(kmsAssets),
		AttackPathsOpened:   len(twin.AttackPaths) / 3,
		RecoveryAssumptions: []string{"Key backup available", "Key escrow procedures in place"},
	})

	cloudAssets := findAssetsByKeyword(twin.Assets, []string{"cloud", "aws", "azure", "gcp", "k8s", "container"})
	scenarios = append(scenarios, ResilienceScenario{
		FailurePoint: "Cloud Provider Failure", BusinessImpact: "Critical",
		SecurityImpact:      "Critical — widespread service disruption",
		AffectedAssets:      assetNames(cloudAssets),
		AttackPathsOpened:   len(twin.AttackPaths) / 2,
		RecoveryAssumptions: []string{"Multi-region deployment", "Cloud-agnostic architecture", "Backup provider configured"},
	})

	return scenarios
}

// ── PHASE 15 — CROWN JEWEL ANALYSIS ──

func (e *SDTEngine) identifyCrownJewels(input SDTInput) []CrownJewelRanking {
	var rankings []CrownJewelRanking
	twin := e.buildTwin(input)

	assetThreatMap := map[string]int{}
	for _, t := range input.Threats {
		for _, a := range t.AffectedAssets {
			assetThreatMap[a]++
		}
		for _, comp := range t.AffectedComponents {
			assetThreatMap[comp]++
		}
	}

	for _, asset := range twin.Assets {
		threatCount := assetThreatMap[asset.Name]
		deps := countDependencies(twin.Relationships, asset.ID)
		bizValue := classifyBusinessValue(asset.Name)
		attackValue := classifyAttackValue(asset.Name, threatCount, deps)
		blastRadius := classifyBlastRadius(deps, threatCount)

		score := float64(threatCount)*2.0 + float64(deps)*1.5
		if bizValue == "Critical" {
			score += 3.0
		}
		if attackValue == "Critical" {
			score += 2.0
		}

		rankings = append(rankings, CrownJewelRanking{
			AssetName: asset.Name, BusinessValue: bizValue,
			AttackValue: attackValue, DependencyCount: deps,
			ThreatCount: threatCount, BlastRadius: blastRadius,
			OverallScore: score,
		})
	}

	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].OverallScore > rankings[j].OverallScore
	})

	if len(rankings) > 10 {
		rankings = rankings[:10]
	}

	return rankings
}

// ── PHASE 16 — EXECUTIVE DIGITAL TWIN REPORT ──

func (e *SDTEngine) generateExecutiveReport(r *SDTResult) DigitalTwinReport {
	health := "Healthy"
	debtScore := r.SecurityDebt.RiskScore
	driftCount := 0
	for _, d := range r.ControlDrifts {
		if d.DriftType != "None" {
			driftCount++
		}
	}
	compDriftCount := len(r.ComplianceDrifts)
	riskTrend := r.Timeline.Trend

	attackSurfaceDesc := "Stable"
	if r.AttackSurfaceTrend.GrowthRate > 10 {
		attackSurfaceDesc = "Growing"
	} else if r.AttackSurfaceTrend.GrowthRate < -5 {
		attackSurfaceDesc = "Shrinking"
	}

	if debtScore > 7 || driftCount > 3 || riskTrend == "Worsening" {
		health = "Needs Attention"
	}
	if debtScore > 9 || driftCount > 6 {
		health = "Critical"
	}

	var sections []DigitalTwinSection

	if len(r.CrownJewels) > 0 {
		top := r.CrownJewels[0]
		sections = append(sections, DigitalTwinSection{
			Title:   "Top Crown Jewel",
			Content: fmt.Sprintf("%s (score=%.1f, threats=%d)", top.AssetName, top.OverallScore, top.ThreatCount),
		})
	}
	if len(r.Resilience) > 0 {
		worst := r.Resilience[0]
		for _, rs := range r.Resilience {
			if rs.BusinessImpact == "Critical" && worst.BusinessImpact != "Critical" {
				worst = rs
			}
		}
		sections = append(sections, DigitalTwinSection{
			Title:   "Highest Resilience Risk",
			Content: fmt.Sprintf("%s — Business Impact: %s, Opens %d attack paths", worst.FailurePoint, worst.BusinessImpact, worst.AttackPathsOpened),
		})
	}
	if len(r.EvolutionInsights) > 0 {
		needsReview := 0
		for _, ei := range r.EvolutionInsights {
			if ei.Status == "Needs Review" || ei.Status == "Likely Invalid" {
				needsReview++
			}
		}
		sections = append(sections, DigitalTwinSection{
			Title:   "Evolution Risks",
			Content: fmt.Sprintf("%d assumptions may be invalidated by growth", needsReview),
		})
	}
	if len(r.ZeroTrust.Dimensions) > 0 {
		lowest := r.ZeroTrust.Dimensions[0]
		for _, d := range r.ZeroTrust.Dimensions {
			if d.CurrentScore < lowest.CurrentScore {
				lowest = d
			}
		}
		sections = append(sections, DigitalTwinSection{
			Title:   "Zero Trust Gap",
			Content: fmt.Sprintf("Lowest dimension: %s (%.1f/%.1f)", lowest.Dimension, lowest.CurrentScore, lowest.TargetScore),
		})
	}

	return DigitalTwinReport{
		ArchitectureHealth:   health,
		SecurityDebtScore:    debtScore,
		ControlDriftCount:    driftCount,
		ComplianceDriftCount: compDriftCount,
		RiskTrend:            riskTrend,
		AttackSurfaceTrend:   attackSurfaceDesc,
		Sections:             sections,
	}
}

// ── PHASE 17 — PORTFOLIO DIGITAL TWIN ──

func (e *SDTEngine) summarizePortfolio(input SDTInput) PortfolioTwinSummary {
	summary := PortfolioTwinSummary{
		ArchitectureCount: 1 + len(input.PreviousTwins) + len(input.PortfolioTwins),
	}

	allTwins := append(input.PreviousTwins, input.PortfolioTwins...)
	allTwins = append(allTwins, e.buildTwin(input))

	totalDebt := 0.0
	for _, twin := range allTwins {
		sdi := SDTInput{
			Findings: twin.Findings, Controls: twin.Controls,
			Assumptions: twin.Assumptions,
		}
		debt := e.calculateSecurityDebt(sdi)
		totalDebt += debt.TotalDebt

		for _, t := range twin.Threats {
			if containsAny(strings.ToLower(t.Name), []string{"vendor", "supply", "third"}) {
				if !contains(summary.SharedRisks, t.Name) {
					summary.SharedRisks = append(summary.SharedRisks, t.Name)
				}
			}
		}
		for _, c := range twin.Controls {
			shared := false
			for _, sc := range summary.SharedControls {
				if sc == c.Name {
					shared = true
					break
				}
			}
			if !shared {
				for _, other := range allTwins {
					if other.ID != twin.ID {
						for _, oc := range other.Controls {
							if oc.Name == c.Name {
								summary.SharedControls = append(summary.SharedControls, c.Name)
								shared = true
								break
							}
						}
					}
					if shared {
						break
					}
				}
			}
		}
	}
	summary.AggregatedDebt = totalDebt

	trends := []string{
		fmt.Sprintf("%d total architectures in digital twin portfolio", summary.ArchitectureCount),
		fmt.Sprintf("Aggregated security debt: %.1f", totalDebt),
		fmt.Sprintf("Shared controls across portfolio: %d", len(summary.SharedControls)),
	}
	if len(summary.SharedRisks) > 0 {
		trends = append(trends, fmt.Sprintf("Shared risks identified: %d", len(summary.SharedRisks)))
	}
	summary.EnterpriseTrends = trends

	return summary
}

// ── HELPERS ──

func classifyComponent(name string) string {
	lower := strings.ToLower(name)
	switch {
	case containsAny(lower, []string{"database", "data", "store", "warehouse", "bucket"}):
		return "DataStore"
	case containsAny(lower, []string{"api", "gateway", "service", "endpoint"}):
		return "Service"
	case containsAny(lower, []string{"auth", "identity", "iam", "login", "sso"}):
		return "Identity"
	case containsAny(lower, []string{"network", "firewall", "vpn", "proxy", "load"}):
		return "Network"
	case containsAny(lower, []string{"container", "pod", "cluster", "vm", "node", "server"}):
		return "Infrastructure"
	case containsAny(lower, []string{"web", "app", "ui", "portal", "dashboard"}):
		return "Application"
	case containsAny(lower, []string{"vendor", "third", "partner", "supply"}):
		return "External"
	default:
		return "Component"
	}
}

func classifyBusinessValue(name string) string {
	lower := strings.ToLower(name)
	if containsAny(lower, []string{"payment", "pii", "credential", "key", "auth", "identity"}) {
		return "Critical"
	}
	if containsAny(lower, []string{"database", "api", "gateway", "service"}) {
		return "High"
	}
	if containsAny(lower, []string{"log", "monitor", "cache", "queue"}) {
		return "Medium"
	}
	return "Low"
}

func classifyAttackValue(name string, threatCount, deps int) string {
	lower := strings.ToLower(name)
	if containsAny(lower, []string{"admin", "vault", "key", "auth", "identity"}) && threatCount >= 2 {
		return "Critical"
	}
	if containsAny(lower, []string{"database", "api", "gateway"}) || threatCount >= 3 {
		return "High"
	}
	if threatCount >= 1 || deps >= 3 {
		return "Medium"
	}
	return "Low"
}

func classifyBlastRadius(deps, threatCount int) string {
	if deps >= 5 || threatCount >= 4 {
		return "Critical"
	}
	if deps >= 3 || threatCount >= 2 {
		return "High"
	}
	if deps >= 1 || threatCount >= 1 {
		return "Medium"
	}
	return "Low"
}

func findComponentID(comps []TwinComponent, name string) string {
	for _, c := range comps {
		if strings.EqualFold(c.Name, name) {
			return c.ID
		}
	}
	return ""
}

func hasComponent(comps []TwinComponent, keyword string) bool {
	for _, c := range comps {
		if strings.Contains(strings.ToLower(c.Name), keyword) {
			return true
		}
	}
	return false
}

func findAssetsByKeyword(assets []TwinAsset, keywords []string) []TwinAsset {
	var out []TwinAsset
	for _, a := range assets {
		for _, kw := range keywords {
			if strings.Contains(strings.ToLower(a.Name), kw) {
				out = append(out, a)
				break
			}
		}
	}
	return out
}

func assetNames(assets []TwinAsset) []string {
	var names []string
	for _, a := range assets {
		names = append(names, a.Name)
	}
	return names
}

func countDependencies(rels []TwinRelationship, assetID string) int {
	count := 0
	for _, r := range rels {
		if r.SourceID == assetID || r.TargetID == assetID {
			count++
		}
	}
	return count
}

func countThreatsByCategory(threats []Threat, category string) int {
	count := 0
	for _, t := range threats {
		if strings.Contains(strings.ToLower(string(t.Category)), strings.ToLower(category)) {
			count++
		}
	}
	return count
}

func countControlsByCategory(controls []SDRIControl, category string) int {
	count := 0
	for _, c := range controls {
		if strings.Contains(strings.ToLower(c.Category), strings.ToLower(category)) ||
			strings.Contains(strings.ToLower(c.Name), strings.ToLower(category)) {
			count++
		}
	}
	return count
}

func lenDiff(current, previous int, added bool) int {
	diff := current - previous
	if added && diff > 0 {
		return diff
	}
	if !added && diff < 0 {
		return -diff
	}
	return 0
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
