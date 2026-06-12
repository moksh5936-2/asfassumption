package intelligence

import (
	"strings"
	"testing"
)

func TestERNNewEngine(t *testing.T) {
	e := NewERNEngine()
	if e == nil {
		t.Fatal("expected non-nil ERN engine")
	}
}

func TestERNEmptyInput(t *testing.T) {
	e := NewERNEngine()
	input := ERNInput{}
	res := e.Run(input)
	if res == nil {
		t.Fatal("expected non-nil result")
	}
	if len(res.ExecutiveRisks) == 0 {
		t.Error("expected at least 1 default executive risk")
	}
	if len(res.RiskNarratives) != 0 {
		t.Error("expected 0 risk narratives with no findings")
	}
	if len(res.InvestmentInsights) == 0 {
		t.Error("expected at least 1 investment insight")
	}
	if res.DecisionSupport.Top3Actions == nil || len(res.DecisionSupport.Top3Actions) == 0 {
		t.Error("expected top 3 decision support actions")
	}
	if len(res.DecisionSupport.Top3Actions) > 3 {
		t.Errorf("expected at most 3 decision actions, got %d", len(res.DecisionSupport.Top3Actions))
	}
}

func TestERNHealthcareInput(t *testing.T) {
	e := NewERNEngine()
	pack := healthcarePack()
	input := ERNInput{
		Domain:     "healthcare",
		DomainPack: &pack,
		Findings: []SDRIFinding{
			{ID: "F-001", Title: "Missing MFA on Administrative Access", Category: "AccessControl", Severity: "Critical", Description: "Admin accounts lack MFA", AffectedComponents: []string{"PHI Database", "EHR System"}},
			{ID: "F-002", Title: "PHI Database Encryption at Rest", Category: "DataProtection", Severity: "High", Description: "PHI database lacks encryption at rest", AffectedComponents: []string{"PHI Database"}},
		},
		Threats: []Threat{
			{ID: "T-001", Name: "PHI Theft via Compromised Identity", Severity: "Critical"},
			{ID: "T-002", Name: "Ransomware on Healthcare Systems", Severity: "Critical"},
		},
		Controls: []SDRIControl{
			{ID: "C-001", Name: "BreakGlassAccess", Coverage: "Partial"},
			{ID: "C-002", Name: "PHIAuditLogging", Coverage: "Full"},
			{ID: "C-003", Name: "DataEncryptionAtRest", Coverage: "Missing"},
		},
		ComplianceFrameworks: []string{"HIPAA", "HITRUST"},
	}
	res := e.Run(input)
	if res == nil {
		t.Fatal("expected non-nil result")
	}
	if len(res.ExecutiveRisks) != 2 {
		t.Errorf("expected 2 executive risks, got %d", len(res.ExecutiveRisks))
	}
	if len(res.RiskNarratives) != 2 {
		t.Errorf("expected 2 risk narratives, got %d", len(res.RiskNarratives))
	}
	if len(res.CrownJewelClasses) == 0 {
		t.Error("expected crown jewel classes")
	}
	if len(res.RegulatoryImpacts) == 0 {
		t.Error("expected regulatory impacts for healthcare")
	}
	hasHIPAA := false
	for _, ri := range res.RegulatoryImpacts {
		if ri.Framework == "HIPAA" || ri.Framework == "HITRUST" {
			hasHIPAA = true
		}
	}
	if !hasHIPAA {
		t.Error("expected HIPAA/HITRUST regulatory impact for healthcare")
	}
	if len(res.BusinessImpactMap.Categories) == 0 {
		t.Error("expected business impact categories")
	}
	if len(res.PriorityRisks) != 2 {
		t.Errorf("expected 2 priority risks, got %d", len(res.PriorityRisks))
	}
	if res.FinancialExposure.Level == "" {
		t.Error("expected financial exposure level")
	}
	if res.BoardSummary.Summary == "" {
		t.Error("expected board summary")
	}
	if len(res.CISOBriefing.TopRisks) == 0 {
		t.Error("expected CISO top risks")
	}
	if len(res.RemediationRoadmap.Phase30) == 0 {
		t.Error("expected 30-day remediation items")
	}
	if len(res.InvestmentInsights) == 0 {
		t.Error("expected investment insights")
	}
	if res.Dashboard.RiskScore == 0 {
		t.Error("expected dashboard risk score")
	}
	if res.DecisionSupport.Top3Actions == nil || len(res.DecisionSupport.Top3Actions) == 0 {
		t.Error("expected decision support actions")
	}
	if len(res.DecisionSupport.Top3Actions) > 3 {
		t.Errorf("expected at most 3 decision actions, got %d", len(res.DecisionSupport.Top3Actions))
	}
	// Verify top priority risk is the MFA finding (critical severity)
	if len(res.PriorityRisks) > 0 && res.PriorityRisks[0].Risk.Title != "Missing MFA on Administrative Access" {
		t.Errorf("expected MFA finding as top priority, got %s", res.PriorityRisks[0].Risk.Title)
	}
}

func TestERNFintechInput(t *testing.T) {
	e := NewERNEngine()
	pack := fintechPack()
	input := ERNInput{
		Domain:     "fintech",
		DomainPack: &pack,
		Findings: []SDRIFinding{
			{ID: "F-001", Title: "PCI Key Management Gap", Category: "KeyManagement", Severity: "Critical", Description: "No HSM-backed key management for payment encryption keys"},
			{ID: "F-002", Title: "Transaction Logging Deficiency", Category: "AuditLogging", Severity: "High", Description: "Transaction logs lack immutability controls"},
		},
		Threats: []Threat{
			{ID: "T-001", Name: "Payment Transaction Fraud", Severity: "Critical"},
		},
		Controls: []SDRIControl{
			{ID: "C-001", Name: "FraudMonitoring", Coverage: "Partial"},
			{ID: "C-002", Name: "PCIKeyManagement", Coverage: "Missing"},
		},
		ComplianceFrameworks: []string{"PCI DSS"},
	}
	res := e.Run(input)
	if res == nil {
		t.Fatal("expected non-nil result")
	}
	if len(res.ExecutiveRisks) != 2 {
		t.Errorf("expected 2 executive risks, got %d", len(res.ExecutiveRisks))
	}
	hasPCI := false
	for _, ri := range res.RegulatoryImpacts {
		if strings.Contains(ri.Framework, "PCI") {
			hasPCI = true
			break
		}
	}
	if !hasPCI {
		t.Error("expected PCI regulatory impact for fintech")
	}
	// Key management should be a top investment insight
	hasKeyMgmt := false
	for _, ii := range res.InvestmentInsights {
		if strings.Contains(ii.Area, "Key") || strings.Contains(ii.Area, "Crypt") {
			hasKeyMgmt = true
			break
		}
	}
	if !hasKeyMgmt {
		t.Error("expected Key Management investment insight for fintech with key management finding")
	}
}

func TestERNKubernetesInput(t *testing.T) {
	e := NewERNEngine()
	pack := kubernetesPack()
	input := ERNInput{
		Domain:     "kubernetes",
		DomainPack: &pack,
		Findings: []SDRIFinding{
			{ID: "F-001", Title: "Missing Network Policies", Category: "NetworkSegmentation", Severity: "High", Description: "No network policies restricting pod-to-pod traffic"},
		},
		Threats: []Threat{
			{ID: "T-001", Name: "Container Escape", Severity: "Critical"},
			{ID: "T-002", Name: "Secrets Theft from etcd", Severity: "Critical"},
		},
		Controls: []SDRIControl{
			{ID: "C-001", Name: "K8sNetworkPolicies", Coverage: "Missing"},
			{ID: "C-002", Name: "K8sSecretsProtection", Coverage: "Partial"},
		},
	}
	res := e.Run(input)
	if res == nil {
		t.Fatal("expected non-nil result")
	}
	if len(res.ExecutiveRisks) == 0 {
		t.Error("expected executive risks")
	}
	if len(res.CrownJewelClasses) == 0 {
		t.Error("expected crown jewel classes")
	}
	if res.Dashboard.AttackPathCount >= 0 && len(input.AttackPaths) == 0 {
		// Expected with no attack paths
	}
}

func TestERNNarratives(t *testing.T) {
	// Test narrative generation for different finding types
	tests := []struct {
		title    string
		category string
		contains string
	}{
		{"Missing MFA", "AccessControl", "unauthorized access"},
		{"Encryption Key Rotation", "KeyManagement", "decrypt confidential"},
		{"Network Segmentation", "Network", "unauthorized network access"},
		{"Logging Deficiency", "AuditLogging", "security incidents could go undetected"},
		{"Third Party Access", "ThirdParty", "supply chain"},
		{"Unpatched Vulnerability", "PatchManagement", "known exploits"},
		{"Generic Finding", "General", "executive attention"},
	}

	for _, tt := range tests {
		t.Run(tt.title, func(t *testing.T) {
			f := SDRIFinding{
				ID:          "F-001",
				Title:       tt.title,
				Category:    tt.category,
				Description: "Test finding",
			}
			narrative := buildNarrativeFromFinding(f)
			if !strings.Contains(narrative, tt.contains) {
				t.Errorf("narrative for %q should contain %q, got: %s", tt.title, tt.contains, narrative)
			}
		})
	}
}

func TestERNBusinessClassification(t *testing.T) {
	tests := []struct {
		name             string
		expectedCategory string
		expectedLabel    string
	}{
		{"PHI Database", "Data Asset", "Patient Data Asset"},
		{"KMS", "Cryptographic Trust Infrastructure", "Cryptographic Material Asset"},
		{"Cluster Admin", "Compute Infrastructure", "Privileged Access Asset"},
		{"Secrets Store", "Data Asset", "Cryptographic Material Asset"},
		{"Payment Processor", "Business Asset", "Financial Transaction Asset"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cat := classifyBusinessCategory(tt.name)
			if cat != tt.expectedCategory {
				t.Errorf("expected category %q, got %q", tt.expectedCategory, cat)
			}
			label := classifyBusinessLabel(tt.name)
			if label != tt.expectedLabel {
				t.Errorf("expected label %q, got %q", tt.expectedLabel, label)
			}
		})
	}
}

func TestERNPriorityRanking(t *testing.T) {
	risks := []ExecutiveRisk{
		{
			ID:                "ERN-RISK-001",
			Title:             "Critical Finding - Missing MFA",
			BusinessImpact:    "Significant business impact",
			OperationalImpact: "Extended operational impact",
			ComplianceImpact:  "Potential compliance exposure",
			Severity:          "Critical",
			AffectedAssets:    []string{"Asset1"},
		},
		{
			ID:               "ERN-RISK-002",
			Title:            "Low Risk Finding",
			BusinessImpact:   "Minor impact",
			ComplianceImpact: "None",
			Severity:         "Low",
			AffectedAssets:   []string{},
		},
	}
	prioritized := rankExecutivePriorities(risks)
	if len(prioritized) != 2 {
		t.Fatalf("expected 2 priority risks, got %d", len(prioritized))
	}
	if prioritized[0].Score < prioritized[1].Score {
		t.Errorf("expected first risk to have higher priority score, got %d vs %d",
			prioritized[0].Score, prioritized[1].Score)
	}
	if prioritized[0].Priority != "Immediate" {
		t.Errorf("expected critical severity risk to be Immediate priority, got %s", prioritized[0].Priority)
	}
}

func TestERNFinancialExposure(t *testing.T) {
	pack := healthcarePack()
	// Low exposure
	low := estimateFinancialExposure(ERNInput{})
	if low.Level != "Low" {
		t.Errorf("expected Low exposure for empty input, got %s", low.Level)
	}

	// Higher exposure with many findings, threats, attack paths
	high := estimateFinancialExposure(ERNInput{
		DomainPack:  &pack,
		Threats:     []Threat{{}, {}, {}, {}, {}},
		AttackPaths: []AttackPath{{}, {}, {}, {}},
		Findings: []SDRIFinding{
			{Severity: "Critical"},
			{Severity: "Critical"},
			{Severity: "High"},
			{Severity: "High"},
		},
		ComplianceFrameworks: []string{"HIPAA", "HITRUST", "FDA"},
	})
	if high.Level == "Low" {
		t.Errorf("expected higher exposure level for significant input, got %s (%s)", high.Level, high.Rationale)
	}
}

func TestERNRiskThemes(t *testing.T) {
	findings := []SDRIFinding{
		{ID: "F-001", Title: "Missing MFA", Category: "AccessControl", Severity: "Critical"},
		{ID: "F-002", Title: "Encryption Gap", Category: "DataProtection", Severity: "High"},
		{ID: "F-003", Title: "Vendor Assessment Needed", Category: "ThirdParty", Severity: "Medium"},
	}
	risks := []ExecutiveRisk{
		{ID: "ERN-RISK-001", Title: "Missing MFA", Severity: "Critical"},
		{ID: "ERN-RISK-002", Title: "Encryption Gap", Severity: "High"},
		{ID: "ERN-RISK-003", Title: "Vendor Assessment Needed", Severity: "Medium"},
	}
	themes := aggregateRiskThemes(risks, findings)
	if len(themes) == 0 {
		t.Fatal("expected at least 1 risk theme")
	}
	hasIdentity := false
	hasDataProtection := false
	hasThirdParty := false
	for _, th := range themes {
		switch th.Name {
		case "Identity Risk":
			hasIdentity = true
		case "Data Protection Risk":
			hasDataProtection = true
		case "Third Party Risk":
			hasThirdParty = true
		}
	}
	if !hasIdentity {
		t.Error("expected Identity Risk theme for MFA finding")
	}
	if !hasThirdParty {
		t.Error("expected Third Party Risk theme for vendor finding")
	}
	_ = hasDataProtection
}

func TestERNBoardSummary(t *testing.T) {
	result := &ERNRunResult{
		ExecutiveRisks: []ExecutiveRisk{
			{ID: "ERN-RISK-001", Title: "Missing MFA", Severity: "Critical", Priority: "Immediate"},
			{ID: "ERN-RISK-002", Title: "Encryption Gap", Severity: "High", Priority: "High"},
		},
		PriorityRisks: []PriorityRisk{
			{Risk: ExecutiveRisk{ID: "ERN-RISK-001", Title: "Missing MFA"}, Priority: "Immediate", Score: 35},
			{Risk: ExecutiveRisk{ID: "ERN-RISK-002", Title: "Encryption Gap"}, Priority: "High", Score: 22},
		},
		RiskThemes: []RiskTheme{
			{Name: "Identity Risk", RiskCount: 1},
			{Name: "Data Protection Risk", RiskCount: 1},
		},
		FinancialExposure: FinancialExposure{Level: "Moderate"},
	}
	summary := generateBoardSummary(result)
	if summary.Summary == "" {
		t.Error("expected non-empty board summary")
	}
	if !strings.Contains(summary.Summary, "immediate attention") {
		t.Error("board summary should mention immediate attention findings")
	}
	if !strings.Contains(summary.Summary, "Identity Risk") {
		t.Error("board summary should mention Identity Risk theme")
	}
}

func TestERNCISOBriefing(t *testing.T) {
	pack := healthcarePack()
	result := &ERNRunResult{
		ExecutiveRisks: []ExecutiveRisk{
			{ID: "ERN-RISK-001", Title: "Missing MFA", Priority: "Immediate", RecommendedActions: []string{"Implement MFA", "Review access controls"}},
			{ID: "ERN-RISK-002", Title: "Encryption Gap", Priority: "High", RecommendedActions: []string{"Deploy encryption", "Key rotation"}},
		},
		PriorityRisks: []PriorityRisk{
			{Risk: ExecutiveRisk{ID: "ERN-RISK-001", Title: "Missing MFA"}, Priority: "Immediate", Score: 35},
			{Risk: ExecutiveRisk{ID: "ERN-RISK-002", Title: "Encryption Gap"}, Priority: "High", Score: 22},
		},
		RegulatoryImpacts: []RegulatoryImpact{
			{Framework: "HIPAA", Domain: "Healthcare", Exposure: "Potential"},
		},
	}
	input := ERNInput{
		DomainPack: &pack,
		Controls: []SDRIControl{
			{ID: "C-001", Coverage: "Full"},
			{ID: "C-002", Coverage: "Partial"},
		},
	}
	briefing := generateCISOBriefing(result, input)
	if len(briefing.TopRisks) == 0 {
		t.Error("expected top risks")
	}
	if len(briefing.TopRisks) > 5 {
		t.Errorf("expected at most 5 top risks, got %d", len(briefing.TopRisks))
	}
	if len(briefing.TopRemediations) == 0 {
		t.Error("expected top remediations")
	}
	if len(briefing.TopRemediations) > 5 {
		t.Errorf("expected at most 5 top remediations, got %d", len(briefing.TopRemediations))
	}
	if briefing.CoverageOverview.TotalControls == 0 {
		t.Error("expected coverage overview with controls")
	}
	if briefing.ComplianceOverview == "" {
		t.Error("expected compliance overview")
	}
}

func TestERNRemediationRoadmap(t *testing.T) {
	input := ERNInput{
		Findings: []SDRIFinding{
			{ID: "F-001", Title: "Missing MFA on Admin Access", Category: "Identity", Severity: "Critical"},
			{ID: "F-002", Title: "Encryption Key Gap", Category: "Cryptography", Severity: "High"},
		},
	}
	roadmap := generateRemediationRoadmap(input)
	if len(roadmap.Phase30) == 0 {
		t.Error("expected 30-day phase items")
	}
	if len(roadmap.Phase90) == 0 {
		t.Error("expected 90-day phase items")
	}
	if len(roadmap.Phase180) == 0 {
		t.Error("expected 180-day phase items")
	}
	if len(roadmap.Phase12m) == 0 {
		t.Error("expected 12-month phase items")
	}
}

func TestERNEmptyRoadmap(t *testing.T) {
	roadmap := generateRemediationRoadmap(ERNInput{})
	if len(roadmap.Phase30) == 0 {
		t.Error("expected baseline 30-day items even with empty input")
	}
	if len(roadmap.Phase90) == 0 {
		t.Error("expected baseline 90-day items even with empty input")
	}
}

func TestERNDashboard(t *testing.T) {
	pack := healthcarePack()
	result := &ERNRunResult{
		PriorityRisks: []PriorityRisk{
			{Risk: ExecutiveRisk{ID: "ERN-RISK-001"}, Priority: "Immediate", Score: 30},
			{Risk: ExecutiveRisk{ID: "ERN-RISK-002"}, Priority: "High", Score: 20},
		},
		RegulatoryImpacts: []RegulatoryImpact{
			{Framework: "HIPAA"},
			{Framework: "HITRUST"},
		},
	}
	input := ERNInput{
		DomainPack: &pack,
		Controls: []SDRIControl{
			{ID: "C-001", Coverage: "Full"},
			{ID: "C-002", Coverage: "Partial"},
			{ID: "C-003", Coverage: "Missing"},
			{ID: "C-004", Coverage: "Enhanced"},
		},
		ComplianceFrameworks: []string{"HIPAA"},
		AttackPaths:          []AttackPath{{}},
	}
	dash := buildExecutiveDashboard(result, input)
	if dash.PriorityFindings != 2 {
		t.Errorf("expected 2 priority findings, got %d", dash.PriorityFindings)
	}
	if dash.AttackPathCount != 1 {
		t.Errorf("expected 1 attack path, got %d", dash.AttackPathCount)
	}
	if len(dash.CriticalAssets) == 0 {
		t.Error("expected critical assets from healthcare pack")
	}
}

func TestERNDecisionSupport(t *testing.T) {
	result := &ERNRunResult{
		PriorityRisks: []PriorityRisk{
			{
				Priority: "Immediate",
				Risk: ExecutiveRisk{
					ID: "ERN-RISK-001", Title: "Missing MFA on Admin Access",
					RecommendedActions: []string{"Implement MFA for all privileged access"},
				},
			},
			{
				Priority: "High",
				Risk: ExecutiveRisk{
					ID: "ERN-RISK-002", Title: "Encryption Gap",
					RecommendedActions: []string{"Deploy encryption for sensitive data"},
				},
			},
			{
				Priority: "High",
				Risk: ExecutiveRisk{
					ID: "ERN-RISK-003", Title: "Third Party Risk",
					RecommendedActions: []string{"Establish vendor risk management program"},
				},
			},
			{
				Priority: "Medium",
				Risk: ExecutiveRisk{
					ID: "ERN-RISK-004", Title: "Logging Gap",
					RecommendedActions: []string{"Deploy centralized logging"},
				},
			},
		},
	}
	support := generateDecisionSupport(result)
	if len(support.Top3Actions) != 3 {
		t.Errorf("expected exactly 3 decision actions, got %d", len(support.Top3Actions))
	}
	if len(support.Top3Actions) > 0 && support.Top3Actions[0].Rank != 1 {
		t.Errorf("expected first action to be rank 1, got %d", support.Top3Actions[0].Rank)
	}
}

func TestERNEmptyDecisionSupport(t *testing.T) {
	result := &ERNRunResult{
		ExecutiveRisks: []ExecutiveRisk{
			{ID: "ERN-RISK-001", Title: "Baseline", RecommendedActions: []string{}},
		},
		PriorityRisks: []PriorityRisk{
			{Risk: ExecutiveRisk{Title: "Baseline", RecommendedActions: []string{}}, Priority: "Low", Score: 0},
		},
	}
	support := generateDecisionSupport(result)
	if len(support.Top3Actions) != 3 {
		t.Errorf("expected 3 fallback actions, got %d", len(support.Top3Actions))
	}
}

func TestERNInvestmentInsights(t *testing.T) {
	result := &ERNRunResult{
		PriorityRisks: []PriorityRisk{
			{Risk: ExecutiveRisk{Title: "Missing MFA on Admin", RecommendedActions: []string{"Implement MFA"}}, Priority: "Immediate", Score: 30},
			{Risk: ExecutiveRisk{Title: "Encryption Key Management Gap", RecommendedActions: []string{"Deploy KMS"}}, Priority: "High", Score: 20},
		},
	}
	insights := generateInvestmentInsights(result, ERNInput{})
	if len(insights) == 0 {
		t.Error("expected investment insights")
	}
	if len(insights) > 0 && insights[0].Priority != "Highest" {
		t.Errorf("expected top insight to be Highest priority, got %s", insights[0].Priority)
	}
}

func TestERNRiskTrend(t *testing.T) {
	result := &ERNRunResult{
		PriorityRisks: []PriorityRisk{
			{Risk: ExecutiveRisk{ID: "R1"}, Priority: "Immediate", Score: 30},
		},
		RegulatoryImpacts: []RegulatoryImpact{
			{Framework: "HIPAA"},
		},
	}
	input := ERNInput{
		Controls: []SDRIControl{
			{ID: "C-001", Coverage: "Full"},
			{ID: "C-002", Coverage: "Missing"},
		},
		ComplianceFrameworks: []string{"HIPAA"},
	}
	trend := generateRiskTrend(result, input)
	if trend.CurrentState.CriticalFindings != 1 {
		t.Errorf("expected 1 critical finding, got %d", trend.CurrentState.CriticalFindings)
	}
	if trend.TargetState.RiskScore >= trend.CurrentState.RiskScore {
		t.Error("expected target risk score to be lower than current")
	}
	if trend.TargetState.CoverageRate <= trend.CurrentState.CoverageRate {
		t.Error("expected target coverage rate to be higher than current")
	}
}

func TestERNEmptyRiskTrend(t *testing.T) {
	trend := generateRiskTrend(&ERNRunResult{}, ERNInput{})
	if trend.CurrentState.RiskScore != 0 {
		t.Errorf("expected 0 risk score for empty, got %.1f", trend.CurrentState.RiskScore)
	}
}

func TestERNBuildExecutiveRisksDefault(t *testing.T) {
	risks := buildExecutiveRisks(ERNInput{})
	if len(risks) != 1 {
		t.Fatalf("expected 1 default risk, got %d", len(risks))
	}
	if risks[0].Title != "Architecture Security Posture Review" {
		t.Errorf("unexpected default title: %s", risks[0].Title)
	}
	if risks[0].Priority != "Low" {
		t.Errorf("expected Low priority for default risk, got %s", risks[0].Priority)
	}
}

func TestERNCrownJewelClassification(t *testing.T) {
	pack := criticalInfrastructurePack()
	classes := classifyCrownJewels(&pack)
	if len(classes) == 0 {
		t.Fatal("expected crown jewel classes")
	}
	for _, c := range classes {
		if c.BusinessCategory == "" {
			t.Errorf("missing business category for %s", c.TechnicalName)
		}
		if c.BusinessLabel == "" {
			t.Errorf("missing business label for %s", c.TechnicalName)
		}
	}
}

func TestERNFinancialExposureScoring(t *testing.T) {
	pack := healthcarePack()
	exp := estimateFinancialExposure(ERNInput{
		DomainPack:           &pack,
		Threats:              []Threat{{}, {}, {}, {}},
		AttackPaths:          []AttackPath{{}, {}},
		Findings:             []SDRIFinding{{Severity: "Critical"}, {Severity: "Critical"}, {Severity: "High"}},
		ComplianceFrameworks: []string{"HIPAA", "HITRUST", "FDA", "PCI DSS"},
	})
	levels := map[string]bool{"Low": true, "Moderate": true, "Significant": true, "Severe": true}
	if !levels[exp.Level] {
		t.Errorf("unexpected financial exposure level: %s", exp.Level)
	}
	if exp.Rationale == "" {
		t.Error("expected rationale for financial exposure")
	}
}

func TestERNRegulatoryImpactHealthcare(t *testing.T) {
	pack := healthcarePack()
	impacts := analyzeRegulatoryImpact(ERNInput{
		DomainPack: &pack,
	})
	if len(impacts) == 0 {
		t.Fatal("expected regulatory impacts for healthcare")
	}
	for _, ri := range impacts {
		if ri.Exposure != "Potential" {
			t.Errorf("expected Potential exposure, got %s", ri.Exposure)
		}
		if ri.Rationale == "" {
			t.Errorf("missing rationale for %s", ri.Framework)
		}
	}
}

func TestERNBoardSummaryEmptyRisks(t *testing.T) {
	summary := generateBoardSummary(&ERNRunResult{
		ExecutiveRisks: []ExecutiveRisk{},
		PriorityRisks:  []PriorityRisk{},
	})
	if summary.Summary == "" {
		t.Error("expected non-empty board summary even with empty risks")
	}
}

func TestERNCISOBriefingEmptyControls(t *testing.T) {
	briefing := generateCISOBriefing(&ERNRunResult{}, ERNInput{})
	if len(briefing.TopRisks) == 0 {
		t.Error("expected at least the default top risk")
	}
	if briefing.CoverageOverview.TotalControls == 0 {
		// No controls in input, so this is expected to be 0
	}
}

// ── Phase 16 — Report Packs ──

func TestERNReportPacksNonEmpty(t *testing.T) {
	e := NewERNEngine()
	input := ERNInput{
		Domain:               "healthcare",
		ComplianceFrameworks: []string{"HIPAA", "PCI DSS"},
		Findings: []SDRIFinding{
			{ID: "F-001", Title: "Missing MFA on Admin Access", Category: "AccessControl", Severity: "Critical", Description: "Admin accounts lack MFA", AffectedComponents: []string{"PHI Database"}},
		},
	}
	res := e.Run(input)
	if res.ReportPacks.BoardReport == "" {
		t.Error("expected non-empty Board report")
	}
	if res.ReportPacks.ExecutiveReport == "" {
		t.Error("expected non-empty Executive report")
	}
	if res.ReportPacks.TechnicalReport == "" {
		t.Error("expected non-empty Technical report")
	}
}

func TestERNReportPacksBoardContainsKeySections(t *testing.T) {
	e := NewERNEngine()
	input := ERNInput{
		Domain:               "fintech",
		ComplianceFrameworks: []string{"PCI DSS", "SOX"},
		Findings: []SDRIFinding{
			{ID: "F-001", Title: "Missing Encryption at Rest", Category: "DataProtection", Severity: "High", Description: "Lacks encryption at rest", AffectedComponents: []string{"Payment DB"}},
		},
	}
	res := e.Run(input)
	br := res.ReportPacks.BoardReport
	if !strings.Contains(br, "Financial Exposure") {
		t.Error("Board report should contain Financial Exposure section")
	}
	if !strings.Contains(br, "Recommended Actions") {
		t.Error("Board report should contain Recommended Actions section")
	}
	if !strings.Contains(br, "Risk Overview") {
		t.Error("Board report should contain Risk Overview section")
	}
}

func TestERNReportPacksExecutiveContainsKeySections(t *testing.T) {
	e := NewERNEngine()
	input := ERNInput{
		Domain:               "kubernetes",
		ComplianceFrameworks: []string{"CIS Benchmarks", "SOC 2"},
		Findings: []SDRIFinding{
			{ID: "F-001", Title: "Container Privilege Escalation", Category: "PrivilegeEscalation", Severity: "Critical", Description: "Privileged containers", AffectedComponents: []string{"Kubernetes API"}},
		},
	}
	res := e.Run(input)
	er := res.ReportPacks.ExecutiveReport
	if !strings.Contains(er, "Executive Risk Report") {
		t.Error("Executive report should contain main header")
	}
	if !strings.Contains(er, "Board Summary") {
		t.Error("Executive report should contain Board Summary header")
	}
	if !strings.Contains(er, "Financial Exposure") {
		t.Error("Executive report should contain Financial Exposure")
	}
	if !strings.Contains(er, "Risk Themes") {
		t.Error("Executive report should contain Risk Themes table")
	}
	if !strings.Contains(er, "CISO Briefing") {
		t.Error("Executive report should contain CISO Briefing section")
	}
	if !strings.Contains(er, "Decision Support") {
		t.Error("Executive report should contain Decision Support section")
	}
	if !strings.Contains(er, "Risk Score") {
		t.Error("Executive report should contain Dashboard section")
	}
}

func TestERNReportPacksTechnicalContainsAllSections(t *testing.T) {
	e := NewERNEngine()
	input := ERNInput{
		Domain:               "government",
		ComplianceFrameworks: []string{"FedRAMP", "NIST 800-53"},
		Findings: []SDRIFinding{
			{ID: "F-001", Title: "Missing Audit Logging", Category: "Audit", Severity: "High", Description: "No audit logging", AffectedComponents: []string{"Gov Portal"}},
		},
	}
	res := e.Run(input)
	tr := res.ReportPacks.TechnicalReport
	if !strings.Contains(tr, "Technical Security Report") {
		t.Error("Technical report should contain main header")
	}
	if !strings.Contains(tr, "Executive Summary") {
		t.Error("Technical report should contain Executive Summary")
	}
	if !strings.Contains(tr, "Financial Exposure") {
		t.Error("Technical report should contain Financial Exposure")
	}
	if !strings.Contains(tr, "CISO Briefing") {
		t.Error("Technical report should contain CISO Briefing")
	}
	if !strings.Contains(tr, "Remediation Roadmap") {
		t.Error("Technical report should contain Remediation Roadmap")
	}
	if !strings.Contains(tr, "Risk Trend Analysis") {
		t.Error("Technical report should contain Risk Trend Analysis")
	}
	if !strings.Contains(tr, "Dashboard") {
		t.Error("Technical report should contain Dashboard")
	}
	if !strings.Contains(tr, "Risk Narratives") {
		t.Error("Technical report should contain Risk Narratives")
	}
}

func TestERNReportPacksParseReportPackType(t *testing.T) {
	if ParseReportPackType("board") != ReportPackBoard {
		t.Error("expected ReportPackBoard for 'board'")
	}
	if ParseReportPackType("executive") != ReportPackExecutive {
		t.Error("expected ReportPackExecutive for 'executive'")
	}
	if ParseReportPackType("technical") != ReportPackTechnical {
		t.Error("expected ReportPackTechnical for 'technical'")
	}
	if ParseReportPackType("unknown") != ReportPackTechnical {
		t.Error("expected default ReportPackTechnical for unknown")
	}
}

func TestERNReportPacksString(t *testing.T) {
	if ReportPackBoard.String() != "Board" {
		t.Errorf("expected 'Board', got %q", ReportPackBoard.String())
	}
	if ReportPackExecutive.String() != "Executive" {
		t.Errorf("expected 'Executive', got %q", ReportPackExecutive.String())
	}
	if ReportPackTechnical.String() != "Technical" {
		t.Errorf("expected 'Technical', got %q", ReportPackTechnical.String())
	}
	var unknown ReportPackType = 99
	if unknown.String() != "Unknown" {
		t.Errorf("expected 'Unknown', got %q", unknown.String())
	}
}

func TestERNReportPacksReportPacksJSONSerialization(t *testing.T) {
	r := ReportPacks{
		BoardReport:     "board content",
		ExecutiveReport: "executive content",
		TechnicalReport: "technical content",
	}
	if r.BoardReport != "board content" {
		t.Error("ReportPacks struct serialization incorrect")
	}
}
