package intelligence

import (
	"testing"
	"time"
)

func makeSDTInput() SDTInput {
	return SDTInput{
		ArchitectureName: "TestArch",
		Domain:           "fintech",
		RiskScore:        7.5,
		Coverage:         62.0,
		Assets: []TwinAsset{
			{ID: "asset-1", Name: "Admin Console", Type: "Application", Criticality: "Critical"},
			{ID: "asset-2", Name: "Database", Type: "DataStore", Criticality: "High"},
			{ID: "asset-3", Name: "API Gateway", Type: "Service", Criticality: "High"},
			{ID: "asset-4", Name: "Auth Service", Type: "Identity", Criticality: "Critical"},
		},
		Components: []TwinComponent{
			{ID: "comp-1", Name: "Admin Console", Type: "Application", Zone: "Internal"},
			{ID: "comp-2", Name: "Database", Type: "DataStore", Zone: "Data"},
			{ID: "comp-3", Name: "API Gateway", Type: "Service", Zone: "DMZ"},
			{ID: "comp-4", Name: "Auth Service", Type: "Identity", Zone: "Internal"},
		},
		Relationships: []TwinRelationship{
			{ID: "rel-1", SourceID: "comp-1", TargetID: "comp-3", RelationType: "connects", Encrypted: true},
			{ID: "rel-2", SourceID: "comp-3", TargetID: "comp-2", RelationType: "connects", Encrypted: true},
			{ID: "rel-3", SourceID: "comp-4", TargetID: "comp-1", RelationType: "authenticates_via"},
		},
		Findings: []SDRIFinding{
			{ID: "F-001", Title: "Missing MFA on Admin Access", Category: "AccessControl", Severity: "Critical"},
			{ID: "F-002", Title: "Unencrypted Data at Rest", Category: "DataProtection", Severity: "High"},
			{ID: "F-003", Title: "Missing Audit Logging", Category: "Logging", Severity: "High"},
		},
		Threats: []Threat{
			{ID: "T-001", Name: "Internet-Facing Attack Surface", Category: "External", Severity: "Critical", AffectedAssets: []string{"API Gateway"}},
			{ID: "T-002", Name: "Insider Threat", Category: "Internal", Severity: "High", AffectedAssets: []string{"Admin Console"}},
		},
		AttackPaths: []AttackPath{
			{ID: "AP-001", Name: "Web to Database via API", EntryPoint: "External", TargetAsset: "Database", RiskScore: 8.5},
			{ID: "AP-002", Name: "Admin Console Takeover", EntryPoint: "Internal", TargetAsset: "Admin Console", RiskScore: 7.0},
		},
		Controls: []SDRIControl{
			{ID: "C-AUDIT", Name: "AuditLogging", Category: "Logging"},
			{ID: "C-FW", Name: "FirewallRules", Category: "NetworkSecurity"},
		},
		Compliance: []string{"HIPAA", "SOC 2", "PCI DSS"},
		Assumptions: []Assumption{
			{ID: "A-001", Description: "Network segmentation isolates production from development", EvidenceSources: []string{"Network Diagram"}},
			{ID: "A-002", Description: "All API traffic is encrypted", VerificationStatus: "CONTRADICTED"},
		},
		PreviousTwins: []ArchitectureTwin{
			{
				ID: "twin-prev", Version: "0.5", Timestamp: time.Now().Add(-30 * 24 * time.Hour),
				ArchitectureName: "TestArch", RiskScore: 6.0, Coverage: 55.0,
				Compliance: []string{"HIPAA"},
				Threats: []Threat{
					{ID: "T-prev-1", Name: "Legacy Vulnerability", Severity: "High"},
				},
				Controls: []SDRIControl{
					{ID: "C-LEGACY", Name: "LegacyControl", Category: "Governance"},
				},
			},
		},
	}
}

func TestSDTPhase1TwinModel(t *testing.T) {
	e := NewSDTEngine()
	input := makeSDTInput()
	res := e.Run(input)
	if res.Twin.ID == "" {
		t.Error("expected twin ID")
	}
	if len(res.Twin.Assets) == 0 {
		t.Error("expected assets in twin")
	}
	if len(res.Twin.Components) == 0 {
		t.Error("expected components in twin")
	}
	if res.Twin.RiskScore != 7.5 {
		t.Errorf("expected risk score 7.5, got %.1f", res.Twin.RiskScore)
	}
}

func TestSDTPhase2ChangeImpact(t *testing.T) {
	e := NewSDTEngine()
	input := makeSDTInput()
	res := e.Run(input)
	if len(res.ChangeImpacts) == 0 {
		t.Fatal("expected change impacts")
	}
	hasIdentity := false
	for _, ci := range res.ChangeImpacts {
		if ci.Change == "Replace Identity Provider" {
			hasIdentity = true
		}
		if ci.RisksAffected <= 0 && ci.AttackPathsAffected <= 0 {
			t.Errorf("change %s has zero affected items", ci.Change)
		}
	}
	if !hasIdentity {
		t.Error("expected Replace Identity Provider impact (Auth Service present)")
	}
}

func TestSDTPhase3SecurityDiff(t *testing.T) {
	e := NewSDTEngine()
	input := makeSDTInput()
	res := e.Run(input)
	if len(res.ArchitectureDiffs) == 0 {
		t.Fatal("expected architecture diffs")
	}
	for _, d := range res.ArchitectureDiffs {
		if d.Category == "" {
			t.Error("expected diff category")
		}
	}
}

func TestSDTPhase4EvolutionAnalysis(t *testing.T) {
	e := NewSDTEngine()
	input := makeSDTInput()
	res := e.Run(input)
	if len(res.EvolutionInsights) == 0 {
		t.Fatal("expected evolution insights")
	}
	for _, ei := range res.EvolutionInsights {
		if ei.Scenario == "" {
			t.Error("expected scenario name")
		}
		if ei.Status == "" {
			t.Error("expected status")
		}
	}
}

func TestSDTPhase5ControlDrift(t *testing.T) {
	e := NewSDTEngine()
	input := makeSDTInput()
	res := e.Run(input)
	if len(res.ControlDrifts) == 0 {
		t.Fatal("expected control drifts")
	}
	foundMissing := false
	for _, cd := range res.ControlDrifts {
		if cd.CurrentState == "Missing" {
			foundMissing = true
			break
		}
	}
	if !foundMissing {
		t.Error("expected at least one missing control drift")
	}
}

func TestSDTPhase6AssumptionDecay(t *testing.T) {
	e := NewSDTEngine()
	input := makeSDTInput()
	res := e.Run(input)
	if len(res.AssumptionDecays) == 0 {
		t.Fatal("expected assumption decays")
	}
	foundInvalid := false
	for _, ad := range res.AssumptionDecays {
		if ad.Status == "Likely Invalid" {
			foundInvalid = true
		}
		if ad.AssumptionID == "" {
			t.Error("expected assumption ID")
		}
	}
	if !foundInvalid {
		t.Error("expected at least one Likely Invalid assumption")
	}
}

func TestSDTPhase7SecurityDebt(t *testing.T) {
	e := NewSDTEngine()
	input := makeSDTInput()
	res := e.Run(input)
	if res.SecurityDebt.TotalDebt <= 0 {
		t.Error("expected positive security debt")
	}
	if res.SecurityDebt.FindingDebt <= 0 {
		t.Error("expected positive finding debt")
	}
}

func TestSDTPhase8ComplianceDrift(t *testing.T) {
	e := NewSDTEngine()
	input := makeSDTInput()
	res := e.Run(input)
	if len(res.ComplianceDrifts) == 0 {
		t.Fatal("expected compliance drifts")
	}
	for _, cd := range res.ComplianceDrifts {
		if cd.Framework == "" {
			t.Error("expected framework name")
		}
	}
}

func TestSDTPhase9AttackSurface(t *testing.T) {
	e := NewSDTEngine()
	input := makeSDTInput()
	res := e.Run(input)
	if res.AttackSurfaceTrend.InternetExposure <= 0 && res.AttackSurfaceTrend.IdentitySystems <= 0 {
		t.Error("expected at least some attack surface exposure")
	}
}

func TestSDTPhase10Timeline(t *testing.T) {
	e := NewSDTEngine()
	input := makeSDTInput()
	res := e.Run(input)
	if len(res.Timeline.Entries) == 0 {
		t.Fatal("expected timeline entries")
	}
	if res.Timeline.Trend == "" {
		t.Error("expected timeline trend")
	}
}

func TestSDTPhase11WhatIf(t *testing.T) {
	e := NewSDTEngine()
	input := makeSDTInput()
	res := e.Run(input)
	if len(res.WhatIfScenarios) == 0 {
		t.Fatal("expected what-if scenarios")
	}
	for _, wi := range res.WhatIfScenarios {
		if wi.Name == "" {
			t.Error("expected scenario name")
		}
		if wi.RiskDelta == 0 {
			t.Error("expected non-zero risk delta")
		}
	}
}

func TestSDTPhase12MergerAnalysis(t *testing.T) {
	e := NewSDTEngine()
	input := makeSDTInput()
	res := e.Run(input)
	if res.MergerAnalysis.CombinedRiskScore <= 0 {
		t.Error("expected positive combined risk score")
	}
	if res.MergerAnalysis.InheritedRisks <= 0 {
		t.Error("expected inherited risks from previous twin")
	}
}

func TestSDTPhase13ZeroTrust(t *testing.T) {
	e := NewSDTEngine()
	input := makeSDTInput()
	res := e.Run(input)
	if len(res.ZeroTrust.Dimensions) == 0 {
		t.Fatal("expected zero trust dimensions")
	}
	if res.ZeroTrust.Overall <= 0 {
		t.Error("expected positive overall score")
	}
	for _, d := range res.ZeroTrust.Dimensions {
		if d.Dimension == "" {
			t.Error("expected dimension name")
		}
	}
}

func TestSDTPhase14Resilience(t *testing.T) {
	e := NewSDTEngine()
	input := makeSDTInput()
	res := e.Run(input)
	if len(res.Resilience) == 0 {
		t.Fatal("expected resilience scenarios")
	}
	for _, rs := range res.Resilience {
		if rs.FailurePoint == "" {
			t.Error("expected failure point name")
		}
		if rs.BusinessImpact == "" {
			t.Error("expected business impact")
		}
	}
}

func TestSDTPhase15CrownJewel(t *testing.T) {
	e := NewSDTEngine()
	input := makeSDTInput()
	res := e.Run(input)
	if len(res.CrownJewels) == 0 {
		t.Fatal("expected crown jewel rankings")
	}
	if res.CrownJewels[0].OverallScore <= 0 {
		t.Error("expected positive overall score")
	}
}

func TestSDTPhase16ExecutiveReport(t *testing.T) {
	e := NewSDTEngine()
	input := makeSDTInput()
	res := e.Run(input)
	if res.ExecutiveReport.ArchitectureHealth == "" {
		t.Error("expected architecture health")
	}
	if res.ExecutiveReport.SecurityDebtScore <= 0 {
		t.Error("expected positive debt score")
	}
}

func TestSDTPhase17PortfolioSummary(t *testing.T) {
	e := NewSDTEngine()
	input := makeSDTInput()
	res := e.Run(input)
	if res.PortfolioSummary.ArchitectureCount <= 0 {
		t.Error("expected positive architecture count")
	}
	if len(res.PortfolioSummary.EnterpriseTrends) == 0 {
		t.Error("expected enterprise trends")
	}
}

func TestSDTEmptyInput(t *testing.T) {
	e := NewSDTEngine()
	input := SDTInput{}
	res := e.Run(input)
	if res.Twin.ID == "" {
		t.Error("expected twin ID even with empty input")
	}
	if len(res.EvolutionInsights) == 0 {
		t.Error("expected evolution insights even with empty input")
	}
	if len(res.ControlDrifts) == 0 {
		t.Error("expected control drifts even with empty input")
	}
}

func TestSDTTwinDerivation(t *testing.T) {
	e := NewSDTEngine()
	input := SDTInput{
		ArchitectureName: "Minimal",
		Findings: []SDRIFinding{
			{ID: "F-001", Title: "Missing MFA", Category: "AccessControl", Severity: "Critical", AffectedComponents: []string{"Console"}},
		},
	}
	res := e.Run(input)
	if len(res.Twin.Assets) == 0 {
		t.Error("expected assets derived from findings")
	}
	if len(res.Twin.Components) == 0 {
		t.Error("expected components derived from assets")
	}
	if len(res.CrownJewels) == 0 {
		t.Error("expected crown jewel from derived data")
	}
}

func TestSDTAllPhasesPresent(t *testing.T) {
	e := NewSDTEngine()
	input := makeSDTInput()
	res := e.Run(input)
	if res.Twin.ID == "" {
		t.Error("missing twin")
	}
	if res.ChangeImpacts == nil {
		t.Error("missing change impacts")
	}
	if res.ArchitectureDiffs == nil {
		t.Error("missing architecture diffs")
	}
	if res.EvolutionInsights == nil {
		t.Error("missing evolution insights")
	}
	if res.ControlDrifts == nil {
		t.Error("missing control drifts")
	}
	if res.AssumptionDecays == nil {
		t.Error("missing assumption decays")
	}
	if res.ComplianceDrifts == nil {
		t.Error("missing compliance drifts")
	}
	if res.WhatIfScenarios == nil {
		t.Error("missing what-if scenarios")
	}
	if res.Resilience == nil {
		t.Error("missing resilience")
	}
	if res.CrownJewels == nil {
		t.Error("missing crown jewels")
	}
}
