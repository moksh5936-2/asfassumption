package intelligence

import (
	"testing"
	"time"
)

func makeTestArchRecord(id, name, domain string, score float64, findings int, threats int) ArchitectureRecord {
	rec := ArchitectureRecord{
		ArchitectureID: id,
		Name:           name,
		Domain:         domain,
		AnalysisDate:   time.Now(),
		Version:        "1.0",
		RiskScore:      score,
	}
	for i := 0; i < findings; i++ {
		rec.Findings = append(rec.Findings, SDRIFinding{
			ID:       "F-" + id + "-" + string(rune('0'+i)),
			Title:    "Missing MFA on Admin Access",
			Category: "AccessControl",
			Severity: "Critical",
		})
	}
	if findings == 0 {
		rec.Findings = append(rec.Findings, SDRIFinding{
			ID: "F-" + id + "-0", Title: "Missing MFA on Admin Access",
			Category: "AccessControl", Severity: "Critical",
		})
	}
	for i := 0; i < threats; i++ {
		rec.Threats = append(rec.Threats, Threat{
			ID: "T-" + id, Name: "Internet-Facing Attack Surface",
			Severity: "Critical", Description: "External threat",
		})
	}
	if threats == 0 {
		rec.Threats = append(rec.Threats, Threat{
			ID: "T-" + id, Name: "External Threat",
			Severity: "Medium", Description: "Generic threat",
		})
	}
	rec.Controls = append(rec.Controls, SDRIControl{
		ID: "C-MFA", Name: "Multi-Factor Authentication",
		Category: "AccessControl", ControlType: "Preventive",
	})
	rec.Controls = append(rec.Controls, SDRIControl{
		ID: "C-ENC", Name: "Encryption at Rest",
		Category: "DataProtection", ControlType: "Preventive",
	})
	rec.Compliance = append(rec.Compliance, "HIPAA", "SOC 2")
	return rec
}

func TestSAMPIEmptyPortfolio(t *testing.T) {
	e := NewSAMPIEngine()
	input := SAMPIInput{Portfolio: NewPortfolio()}
	res := e.Run(input)
	if res == nil {
		t.Fatal("expected non-nil result")
	}
	if len(res.RepeatedWeaknesses) != 0 {
		t.Error("expected 0 repeated weaknesses from empty portfolio")
	}
}

func TestSAMPIRepeatedWeaknesses(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	p.AddArchitecture(makeTestArchRecord("arch1", "Payment Gateway", "fintech", 7.5, 2, 1))
	p.AddArchitecture(makeTestArchRecord("arch2", "Customer Portal", "fintech", 6.0, 3, 2))
	p.AddArchitecture(makeTestArchRecord("arch3", "Admin Console", "fintech", 5.0, 1, 1))
	res := e.Run(SAMPIInput{Portfolio: p})
	if len(res.RepeatedWeaknesses) == 0 {
		t.Fatal("expected repeated weaknesses")
	}
	found := false
	for _, rw := range res.RepeatedWeaknesses {
		if rw.OccurrenceCount >= 3 && rw.Systemic {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected systemic weakness (occurrence >= 3)")
	}
}

func TestSAMPIEnterpriseRiskThemes(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	rec := makeTestArchRecord("arch1", "Test Arch", "fintech", 5.0, 1, 1)
	rec.Findings = append(rec.Findings, SDRIFinding{
		ID: "F-001", Title: "Third Party Vendor Access Without Review",
		Category: "ThirdPartyRisk", Severity: "High", Description: "Third party access",
	})
	p.AddArchitecture(rec)
	rec2 := makeTestArchRecord("arch2", "Test Arch 2", "fintech", 6.0, 1, 1)
	rec2.Findings = append(rec2.Findings, SDRIFinding{
		ID: "F-002", Title: "Cloud Storage Bucket Publicly Accessible",
		Category: "CloudSecurity", Severity: "Critical", Description: "Cloud exposure",
	})
	p.AddArchitecture(rec2)
	res := e.Run(SAMPIInput{Portfolio: p})
	if len(res.EnterpriseRiskThemes) == 0 {
		t.Fatal("expected enterprise risk themes")
	}
	hasThirdParty := false
	hasCloud := false
	for _, theme := range res.EnterpriseRiskThemes {
		if theme.Name == "Third Party Risk" {
			hasThirdParty = true
		}
		if theme.Name == "Cloud Risk" {
			hasCloud = true
		}
	}
	if !hasThirdParty {
		t.Error("expected Third Party Risk theme")
	}
	if !hasCloud {
		t.Error("expected Cloud Risk theme")
	}
}

func TestSAMPIControlReuse(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	p.AddArchitecture(makeTestArchRecord("arch1", "Arch 1", "fintech", 5.0, 1, 1))
	p.AddArchitecture(makeTestArchRecord("arch2", "Arch 2", "healthcare", 6.0, 1, 1))
	res := e.Run(SAMPIInput{Portfolio: p})
	if len(res.ControlCoverage) == 0 {
		t.Fatal("expected control coverage items")
	}
	for _, cc := range res.ControlCoverage {
		if cc.ControlName == "multi-factor_authentication" && cc.CoveragePercent != 100 {
			t.Errorf("expected 100%% MFA coverage, got %.1f%%", cc.CoveragePercent)
		}
	}
}

func TestSAMPIComparison(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	p.AddArchitecture(makeTestArchRecord("arch1", "Arch A", "fintech", 7.0, 1, 1))
	p.AddArchitecture(makeTestArchRecord("arch2", "Arch B", "fintech", 5.0, 1, 1))
	res := e.Run(SAMPIInput{Portfolio: p})
	if len(res.Comparisons) != 1 {
		t.Fatalf("expected 1 comparison, got %d", len(res.Comparisons))
	}
	cmp := res.Comparisons[0]
	if cmp.ArchitectureA != "Arch A" || cmp.ArchitectureB != "Arch B" {
		t.Errorf("unexpected architecture names: %s vs %s", cmp.ArchitectureA, cmp.ArchitectureB)
	}
	if cmp.SimilarityScore < 0 || cmp.SimilarityScore > 100 {
		t.Errorf("similarity score out of range: %.1f", cmp.SimilarityScore)
	}
}

func TestSAMPIArchitectureComparisonEdge(t *testing.T) {
	p := NewPortfolio()
	p.AddArchitecture(makeTestArchRecord("arch1", "Only Arch", "fintech", 5.0, 1, 1))
	res := NewSAMPIEngine().Run(SAMPIInput{Portfolio: p})
	if len(res.Comparisons) != 0 {
		t.Error("expected 0 comparisons with single architecture")
	}
}

func TestSAMPIRiskTrends(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	rec1 := makeTestArchRecord("same-id", "App v1", "fintech", 8.0, 2, 1)
	rec1.AnalysisDate = time.Now().Add(-30 * 24 * time.Hour)
	rec2 := makeTestArchRecord("same-id", "App v2", "fintech", 4.0, 2, 1)
	rec2.AnalysisDate = time.Now()
	p.Architectures = append(p.Architectures, rec1, rec2)
	res := e.Run(SAMPIInput{Portfolio: p})
	if len(res.RiskTrends) == 0 {
		t.Fatal("expected risk trends")
	}
	if res.RiskTrends[0].Direction != "Improving" {
		t.Errorf("expected Improving trend, got %s", res.RiskTrends[0].Direction)
	}
	if res.RiskTrends[0].PreviousScore != 8.0 || res.RiskTrends[0].CurrentScore != 4.0 {
		t.Error("unexpected risk scores in trend")
	}
}

func TestSAMPISecurityDebt(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	for i := 0; i < 4; i++ {
		rec := makeTestArchRecord("arch"+string(rune('0'+i)), "Arch "+string(rune('0'+i)), "fintech", 5.0, 1, 1)
		p.AddArchitecture(rec)
	}
	res := e.Run(SAMPIInput{Portfolio: p})
	if res.SecurityDebt.RepeatedCount == 0 {
		t.Error("expected at least 1 repeated finding across 4 architectures")
	}
	if res.SecurityDebt.Score <= 0 {
		t.Error("expected positive security debt score")
	}
}

func TestSAMPIEmptySecurityDebt(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	rec := makeTestArchRecord("arch1", "Only Arch", "fintech", 5.0, 0, 0)
	rec.Findings = nil
	p.AddArchitecture(rec)
	res := e.Run(SAMPIInput{Portfolio: p})
	if res.SecurityDebt.Score != 0 {
		t.Errorf("expected 0 security debt with no findings, got %.1f", res.SecurityDebt.Score)
	}
	if res.SecurityDebt.LongstandingCount != 0 {
		t.Errorf("expected 0 longstanding with no findings, got %d", res.SecurityDebt.LongstandingCount)
	}
}

func TestSAMPIEmptyRiskTrend(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	rec := makeTestArchRecord("arch1", "Only Arch", "fintech", 5.0, 1, 1)
	p.AddArchitecture(rec)
	res := e.Run(SAMPIInput{Portfolio: p})
	if len(res.RiskTrends) != 0 {
		t.Error("expected 0 risk trends with single version")
	}
}

func TestSAMPIWorseningRiskTrend(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	rec1 := makeTestArchRecord("id", "App v1", "fintech", 3.0, 1, 1)
	rec1.AnalysisDate = time.Now().Add(-30 * 24 * time.Hour)
	rec2 := makeTestArchRecord("id", "App v2", "fintech", 8.0, 2, 2)
	rec2.AnalysisDate = time.Now()
	p.Architectures = append(p.Architectures, rec1, rec2)
	res := e.Run(SAMPIInput{Portfolio: p})
	if len(res.RiskTrends) == 0 {
		t.Fatal("expected risk trends")
	}
	if res.RiskTrends[0].Direction != "Worsening" {
		t.Errorf("expected Worsening trend, got %s", res.RiskTrends[0].Direction)
	}
}

func TestSAMPIStableRiskTrend(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	rec1 := makeTestArchRecord("id", "App v1", "fintech", 5.0, 1, 1)
	rec1.AnalysisDate = time.Now().Add(-30 * 24 * time.Hour)
	rec2 := makeTestArchRecord("id", "App v2", "fintech", 5.2, 1, 1)
	rec2.AnalysisDate = time.Now()
	p.Architectures = append(p.Architectures, rec1, rec2)
	res := e.Run(SAMPIInput{Portfolio: p})
	if len(res.RiskTrends) == 0 {
		t.Fatal("expected risk trends")
	}
	if res.RiskTrends[0].Direction != "Stable" {
		t.Errorf("expected Stable trend, got %s", res.RiskTrends[0].Direction)
	}
}

func TestSAMPIRepeatedWeaknessSorting(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	p.AddArchitecture(makeTestArchRecord("a1", "A1", "fintech", 5.0, 3, 1))
	p.AddArchitecture(makeTestArchRecord("a2", "A2", "fintech", 5.0, 2, 1))
	p.AddArchitecture(makeTestArchRecord("a3", "A3", "fintech", 5.0, 1, 1))
	p.AddArchitecture(makeTestArchRecord("a4", "A4", "fintech", 5.0, 1, 1))
	res := e.Run(SAMPIInput{Portfolio: p})
	if len(res.RepeatedWeaknesses) == 0 {
		t.Fatal("expected repeated weaknesses")
	}
	for i := 1; i < len(res.RepeatedWeaknesses); i++ {
		if res.RepeatedWeaknesses[i].OccurrenceCount > res.RepeatedWeaknesses[i-1].OccurrenceCount {
			t.Error("expected descending sort by occurrence count")
		}
	}
}

func TestSAMPIEnterpriseThemeSorting(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	rec := makeTestArchRecord("a1", "A1", "fintech", 5.0, 1, 1)
	rec.Findings = append(rec.Findings,
		SDRIFinding{ID: "F-1", Title: "Cloud Misconfiguration", Category: "Cloud", Severity: "High", Description: "cloud aws"},
		SDRIFinding{ID: "F-2", Title: "IAM Policy Violation", Category: "Identity", Severity: "Critical", Description: "identity access"},
		SDRIFinding{ID: "F-3", Title: "Third Party Vendor Risk", Category: "ThirdParty", Severity: "High", Description: "vendor third party"},
		SDRIFinding{ID: "F-4", Title: "Missing Encryption", Category: "DataProtection", Severity: "High", Description: "data protection"},
		SDRIFinding{ID: "F-5", Title: "Logging Deficiency", Category: "Monitoring", Severity: "Medium", Description: "monitoring logging"},
	)
	p.AddArchitecture(rec)
	res := e.Run(SAMPIInput{Portfolio: p})
	for i := 1; i < len(res.EnterpriseRiskThemes); i++ {
		if res.EnterpriseRiskThemes[i].RiskCount > res.EnterpriseRiskThemes[i-1].RiskCount {
			t.Error("expected descending sort by risk count")
		}
	}
}

func TestSAMPIHeatmapGeneration(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	p.AddArchitecture(makeTestArchRecord("a1", "High Risk Arch", "fintech", 9.0, 3, 2))
	p.AddArchitecture(makeTestArchRecord("a2", "Low Risk Arch", "fintech", 2.0, 1, 1))
	res := e.Run(SAMPIInput{Portfolio: p})
	if len(res.Heatmaps) != 2 {
		t.Fatalf("expected 2 heatmaps, got %d", len(res.Heatmaps))
	}
	if res.Heatmaps[0].RiskScore != 9.0 || res.Heatmaps[1].RiskScore != 2.0 {
		t.Error("expected descending sort by risk score")
	}
}

func TestSAMPIHeatmapRiskBand(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	p.AddArchitecture(makeTestArchRecord("a1", "Critical", "fintech", 9.0, 3, 2))
	p.AddArchitecture(makeTestArchRecord("a2", "High", "fintech", 6.0, 2, 1))
	p.AddArchitecture(makeTestArchRecord("a3", "Medium", "fintech", 4.0, 1, 1))
	p.AddArchitecture(makeTestArchRecord("a4", "Low", "fintech", 1.0, 0, 0))
	res := e.Run(SAMPIInput{Portfolio: p})
	bands := make(map[string]bool)
	for _, h := range res.Heatmaps {
		bands[h.RiskBand] = true
	}
	if !bands["Critical"] || !bands["High"] || !bands["Medium"] || !bands["Low"] {
		t.Errorf("expected all 4 risk bands, got %v", bands)
	}
}

func TestSAMPIEnterpriseComplianceView(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	p.AddArchitecture(makeTestArchRecord("a1", "A1", "fintech", 5.0, 1, 1))
	p.AddArchitecture(makeTestArchRecord("a2", "A2", "healthcare", 5.0, 1, 1))
	res := e.Run(SAMPIInput{Portfolio: p})
	if len(res.ComplianceView.Frameworks) == 0 {
		t.Fatal("expected compliance frameworks")
	}
	if res.ComplianceView.TotalArchitectures != 2 {
		t.Errorf("expected 2 total architectures, got %d", res.ComplianceView.TotalArchitectures)
	}
}

func TestSAMPIProgramInsights(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	for i := 0; i < 5; i++ {
		rec := makeTestArchRecord("a"+string(rune('0'+i)), "A"+string(rune('0'+i)), "fintech", 5.0, 2, 1)
		rec.Findings = append(rec.Findings,
			SDRIFinding{ID: "F-id", Title: "Third Party Vendor Risk", Category: "ThirdParty", Severity: "High", Description: "vendor third party"},
		)
		p.AddArchitecture(rec)
	}
	res := e.Run(SAMPIInput{Portfolio: p})
	if len(res.ProgramInsights) == 0 {
		t.Fatal("expected program insights")
	}
	hasSystemic := false
	hasTheme := false
	for _, pi := range res.ProgramInsights {
		if pi.Area == "Systemic Weakness Remediation" {
			hasSystemic = true
		}
		if pi.Area == "Third Party Risk" {
			hasTheme = true
		}
	}
	if !hasSystemic {
		t.Error("expected Systemic Weakness Remediation insight")
	}
	if !hasTheme {
		t.Error("expected Third Party Risk insight")
	}
}

func TestSAMPIProgramInsightsEmpty(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	p.AddArchitecture(makeTestArchRecord("a1", "A1", "fintech", 5.0, 0, 0))
	res := e.Run(SAMPIInput{Portfolio: p})
	if len(res.ProgramInsights) == 0 {
		t.Error("expected at least default program insights from repeated weaknesses analysis")
	}
}

func TestSAMPIAttackSurface(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	rec := makeTestArchRecord("a1", "A1", "fintech", 5.0, 1, 0)
	rec.Threats = []Threat{
		{ID: "T-1", Name: "Internet-Facing API Exposed", Severity: "Critical", Description: "internet exposure"},
		{ID: "T-2", Name: "Third Party Vendor Data Leak", Severity: "High", Description: "third party vendor risk"},
	}
	p.AddArchitecture(rec)
	res := e.Run(SAMPIInput{Portfolio: p})
	if res.AttackSurface.InternetExposure == 0 || res.AttackSurface.ThirdPartyExposure == 0 {
		t.Error("expected internet and third party exposure")
	}
	if res.AttackSurface.TotalExposure <= 0 {
		t.Error("expected positive total exposure")
	}
}

func TestSAMPIEmptyAttackSurface(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	rec := makeTestArchRecord("a1", "A1", "fintech", 5.0, 0, 0)
	rec.Threats = []Threat{}
	res := e.Run(SAMPIInput{Portfolio: p})
	if res.AttackSurface.TotalExposure != 0 {
		t.Error("expected 0 exposure with no threats")
	}
}

func TestSAMPIBlastRadii(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	rec := makeTestArchRecord("a1", "A1", "fintech", 5.0, 1, 0)
	rec.Threats = []Threat{
		{ID: "T-1", Name: "Database Breach", Severity: "Critical", AffectedAssets: []string{"Payment DB"}, AffectedComponents: []string{"Payment DB", "API Gateway"}},
	}
	p.AddArchitecture(rec)
	res := e.Run(SAMPIInput{Portfolio: p})
	if len(res.BlastRadii) == 0 {
		t.Fatal("expected blast radii")
	}
	if res.BlastRadii[0].ComponentName != "Payment DB" {
		t.Errorf("expected Payment DB, got %s", res.BlastRadii[0].ComponentName)
	}
}

func TestSAMPIBlastRadiiEmpty(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	rec := makeTestArchRecord("a1", "A1", "fintech", 5.0, 0, 0)
	rec.Threats = []Threat{}
	p.AddArchitecture(rec)
	res := e.Run(SAMPIInput{Portfolio: p})
	if len(res.BlastRadii) != 0 {
		t.Error("expected 0 blast radii with no threats")
	}
}

func TestSAMPIBlastRadiiSeveritySorting(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	rec := makeTestArchRecord("a1", "A1", "fintech", 5.0, 0, 0)
	rec.Threats = []Threat{
		{ID: "T-1", Name: "Critical", Severity: "Critical", AffectedAssets: []string{"Critical DB"}, AffectedComponents: []string{"Critical DB"}},
		{ID: "T-2", Name: "Low", Severity: "Low", AffectedAssets: []string{"Low Svc"}, AffectedComponents: []string{"Low Svc"}},
	}
	p.AddArchitecture(rec)
	res := e.Run(SAMPIInput{Portfolio: p})
	if len(res.BlastRadii) >= 2 {
		if res.BlastRadii[0].Severity != "Critical" {
			t.Error("expected Critical severity first in sorted blast radii")
		}
	}
}

func TestSAMPIRepeatedWeaknessSystemic(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	for i := 0; i < 5; i++ {
		p.AddArchitecture(makeTestArchRecord("a"+string(rune('0'+i)), "A"+string(rune('0'+i)), "fintech", 5.0, 1, 1))
	}
	res := e.Run(SAMPIInput{Portfolio: p})
	allSystemic := true
	for _, rw := range res.RepeatedWeaknesses {
		if rw.OccurrenceCount >= 3 && !rw.Systemic {
			allSystemic = false
			break
		}
	}
	if !allSystemic {
		t.Error("expected all high-occurrence weaknesses to be marked systemic")
	}
}

func TestSAMPIRepeatedWeaknessNoSystemic(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	rec := makeTestArchRecord("a1", "A1", "fintech", 5.0, 1, 1)
	rec.Findings = []SDRIFinding{
		{ID: "F-1", Title: "Unique Finding Alpha", Category: "Unique", Severity: "Low"},
	}
	p.AddArchitecture(rec)
	rec2 := makeTestArchRecord("a2", "A2", "fintech", 5.0, 1, 1)
	rec2.Findings = []SDRIFinding{
		{ID: "F-2", Title: "Unique Finding Beta", Category: "Unique", Severity: "Low"},
	}
	p.AddArchitecture(rec2)
	res := e.Run(SAMPIInput{Portfolio: p})
	for _, rw := range res.RepeatedWeaknesses {
		if rw.Systemic {
			t.Error("expected no systemic weaknesses with unique findings across architectures")
		}
	}
}

func TestSAMPIEmptyComplianceView(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	rec := makeTestArchRecord("a1", "A1", "fintech", 5.0, 0, 0)
	rec.Compliance = []string{}
	p.AddArchitecture(rec)
	res := e.Run(SAMPIInput{Portfolio: p})
	if len(res.ComplianceView.Frameworks) != 0 {
		t.Error("expected 0 compliance frameworks with empty compliance list")
	}
	if res.ComplianceView.TotalArchitectures != 1 {
		t.Errorf("expected 1 total architecture, got %d", res.ComplianceView.TotalArchitectures)
	}
}

func TestSAMPIEmptyDashboard(t *testing.T) {
	e := NewSAMPIEngine()
	res := e.Run(SAMPIInput{Portfolio: nil})
	if res == nil {
		t.Fatal("expected non-nil result with nil portfolio")
	}
	if res.Dashboard.TotalArchitectures != 0 {
		t.Error("expected 0 in dashboard with nil portfolio")
	}
}

func TestSAMPIDashboardAggregation(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	p.AddArchitecture(makeTestArchRecord("a1", "A1", "fintech", 7.0, 3, 2))
	p.AddArchitecture(makeTestArchRecord("a2", "A2", "fintech", 5.0, 1, 1))
	res := e.Run(SAMPIInput{Portfolio: p})
	if res.Dashboard.TotalArchitectures != 2 {
		t.Errorf("expected 2 architectures, got %d", res.Dashboard.TotalArchitectures)
	}
	if res.Dashboard.TotalFindings == 0 {
		t.Error("expected non-zero total findings")
	}
	if res.Dashboard.AverageRiskScore == 0 {
		t.Error("expected non-zero average risk score")
	}
}

func TestSAMPIProgramInsightsControlGaps(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	rec := makeTestArchRecord("a1", "A1", "fintech", 5.0, 1, 1)
	rec.Controls = []SDRIControl{
		{ID: "C-1", Name: "Unique Control Alpha"},
	}
	p.AddArchitecture(rec)
	rec2 := makeTestArchRecord("a2", "A2", "fintech", 5.0, 1, 1)
	rec2.Controls = []SDRIControl{
		{ID: "C-2", Name: "Unique Control Beta"},
	}
	p.AddArchitecture(rec2)
	res := e.Run(SAMPIInput{Portfolio: p})
	hasGapInsight := false
	for _, pi := range res.ProgramInsights {
		if pi.Area == "Control Coverage Gaps" {
			hasGapInsight = true
			break
		}
	}
	if !hasGapInsight {
		t.Error("expected Control Coverage Gaps insight with low-coverage controls")
	}
}

func TestSAMPIProgramInsightsAttackSurface(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	rec := makeTestArchRecord("a1", "A1", "fintech", 5.0, 1, 0)
	rec.Threats = []Threat{
		{ID: "T-1", Name: "Internet Exposure", Severity: "Critical", Description: "internet exposure"},
		{ID: "T-2", Name: "Cloud Misconfiguration", Severity: "High", Description: "cloud aws"},
	}
	p.AddArchitecture(rec)
	res := e.Run(SAMPIInput{Portfolio: p})
	hasSurfaceInsight := false
	for _, pi := range res.ProgramInsights {
		if pi.Area == "Attack Surface Reduction" {
			hasSurfaceInsight = true
			break
		}
	}
	if !hasSurfaceInsight {
		t.Error("expected Attack Surface Reduction insight with exposure")
	}
}

func TestSAMPIProgramInsightsSharedDependency(t *testing.T) {
	e := NewSAMPIEngine()
	p := NewPortfolio()
	for i := 0; i < 15; i++ {
		p.AddArchitecture(makeTestArchRecord("a"+string(rune('0'+i)), "A"+string(rune('0'+i)), "fintech", 5.0, 1, 1))
	}
	res := e.Run(SAMPIInput{Portfolio: p})
	hasDepInsight := false
	for _, pi := range res.ProgramInsights {
		if pi.Area == "Shared Dependency Risk" {
			hasDepInsight = true
			break
		}
	}
	if !hasDepInsight {
		t.Error("expected Shared Dependency Risk insight with 15 architectures sharing controls")
	}
}

func TestSAMPIPortfolioStorage(t *testing.T) {
	dir := t.TempDir()
	p := NewPortfolio()
	p.AddArchitecture(makeTestArchRecord("arch1", "Payment Gateway", "fintech", 7.5, 2, 1))
	p.AddArchitecture(makeTestArchRecord("arch2", "Customer Portal", "fintech", 6.0, 3, 2))
	path := dir + "/portfolio.json"
	if err := p.Save(path); err != nil {
		t.Fatalf("failed to save portfolio: %v", err)
	}
	loaded := NewPortfolio()
	if err := loaded.Load(path); err != nil {
		t.Fatalf("failed to load portfolio: %v", err)
	}
	if len(loaded.Architectures) != 2 {
		t.Errorf("expected 2 architectures after load, got %d", len(loaded.Architectures))
	}
	if loaded.Architectures[0].Name != "Payment Gateway" {
		t.Errorf("expected Payment Gateway, got %s", loaded.Architectures[0].Name)
	}
}

func TestSAMPIAddArchitectureReplacesExisting(t *testing.T) {
	p := NewPortfolio()
	p.AddArchitecture(makeTestArchRecord("same-id", "Original", "fintech", 5.0, 1, 1))
	p.AddArchitecture(makeTestArchRecord("same-id", "Updated", "fintech", 7.0, 2, 2))
	if len(p.Architectures) != 1 {
		t.Errorf("expected 1 architecture after replace, got %d", len(p.Architectures))
	}
	if p.Architectures[0].Name != "Updated" {
		t.Errorf("expected Updated name, got %s", p.Architectures[0].Name)
	}
	if p.Architectures[0].RiskScore != 7.0 {
		t.Errorf("expected risk score 7.0, got %.1f", p.Architectures[0].RiskScore)
	}
}

func TestSAMPIRemoveArchitecture(t *testing.T) {
	p := NewPortfolio()
	p.AddArchitecture(makeTestArchRecord("a1", "A1", "fintech", 5.0, 1, 1))
	p.AddArchitecture(makeTestArchRecord("a2", "A2", "fintech", 5.0, 1, 1))
	p.RemoveArchitecture("a1")
	if len(p.Architectures) != 1 {
		t.Errorf("expected 1 architecture after remove, got %d", len(p.Architectures))
	}
	if p.Architectures[0].ArchitectureID != "a2" {
		t.Error("expected a2 to remain after removal")
	}
}

func TestSAMIGetArchitecture(t *testing.T) {
	p := NewPortfolio()
	p.AddArchitecture(makeTestArchRecord("a1", "A1", "fintech", 5.0, 1, 1))
	rec := p.GetArchitecture("a1")
	if rec == nil {
		t.Fatal("expected non-nil architecture")
	}
	if rec.Name != "A1" {
		t.Errorf("expected A1, got %s", rec.Name)
	}
	missing := p.GetArchitecture("nonexistent")
	if missing != nil {
		t.Error("expected nil for nonexistent architecture")
	}
}

func TestSAMPIEmptyEdgeCases(t *testing.T) {
	e := NewSAMPIEngine()
	// All empty portfolio
	res := e.Run(SAMPIInput{Portfolio: NewPortfolio()})
	if res.Dashboard.TotalArchitectures != 0 {
		t.Error("expected 0 with empty portfolio")
	}
	// Single record portfolio with no findings or threats
	p := NewPortfolio()
	p.AddArchitecture(ArchitectureRecord{
		ArchitectureID: "empty",
		Name:           "Empty",
		Domain:         "test",
		AnalysisDate:   time.Now(),
	})
	res = e.Run(SAMPIInput{Portfolio: p})
	if res.Dashboard.TotalArchitectures != 1 {
		t.Errorf("expected 1 architecture, got %d", res.Dashboard.TotalArchitectures)
	}
	if len(res.EnterpriseRiskThemes) != 0 {
		t.Error("expected 0 risk themes with no findings")
	}
	if len(res.ControlCoverage) != 0 {
		t.Error("expected 0 control coverage with no controls")
	}
}
