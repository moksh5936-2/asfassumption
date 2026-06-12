package intelligence

import (
	"testing"
)

func makeSDIInput() SDIInput {
	return SDIInput{
		ArchitectureName: "TestArch",
		Domain:           "fintech",
		RiskScore:        7.5,
		Findings: []SDRIFinding{
			{ID: "F-001", Title: "Missing MFA on Admin Access", Category: "AccessControl", Severity: "Critical", AffectedComponents: []string{"Admin Console"}},
			{ID: "F-002", Title: "Unencrypted Data at Rest", Category: "DataProtection", Severity: "High", AffectedComponents: []string{"Database"}},
			{ID: "F-003", Title: "Missing Audit Logging", Category: "Logging", Severity: "High", AffectedComponents: []string{"API Gateway"}},
			{ID: "F-004", Title: "Weak Password Policy", Category: "Authentication", Severity: "Medium", AffectedComponents: []string{"User Portal"}},
			{ID: "F-005", Title: "No Network Segmentation", Category: "NetworkSecurity", Severity: "Critical", AffectedComponents: []string{"Production Network"}},
			{ID: "F-006", Title: "Hardcoded Secrets in Config", Category: "SecretsManagement", Severity: "Critical", AffectedComponents: []string{"Config Service"}},
		},
		Threats: []Threat{
			{ID: "T-001", Name: "Internet-Facing Attack Surface", Category: "External", Severity: "Critical", Description: "External threat via web exposure"},
			{ID: "T-002", Name: "Insider Threat", Category: "Internal", Severity: "High", Description: "Malicious insider with access"},
		},
		AttackPaths: []AttackPath{
			{ID: "AP-001", Name: "MFA Bypass via Admin Console", Description: "External to admin console through missing MFA authentication", RiskScore: 8.5},
			{ID: "AP-002", Name: "Unencrypted Data Exfiltration", Description: "Data theft through unencrypted database without encryption at rest", RiskScore: 7.0},
			{ID: "AP-003", Name: "Secrets Exfiltration via Config Service", Description: "Extract secrets from config service missing secrets vault", RiskScore: 9.0},
			{ID: "AP-004", Name: "Network Pivot Through Flat Network", Description: "Lateral movement across unsegmented network without segmentation", RiskScore: 8.0},
			{ID: "AP-005", Name: "Privilege Escalation via Weak Access", Description: "Privilege escalation through missing access controls", RiskScore: 6.5},
		},
		Controls: []SDRIControl{
			{ID: "C-AUDIT", Name: "AuditLogging", Category: "Logging"},
			{ID: "C-FW", Name: "FirewallRules", Category: "NetworkSecurity"},
		},
		Compliance: []string{"HIPAA", "SOC 2", "PCI DSS"},
	}
}

func TestSDIPhase1DecisionModel(t *testing.T) {
	e := NewSDIEngine()
	input := makeSDIInput()
	res := e.Run(input)
	if res == nil {
		t.Fatal("expected non-nil result")
	}
	if len(res.Recommendations) == 0 {
		t.Fatal("expected at least 1 recommendation")
	}
	rec := res.Recommendations[0]
	if rec.ID == "" {
		t.Error("expected recommendation ID")
	}
	if rec.Title == "" {
		t.Error("expected recommendation title")
	}
	if rec.Priority == "" {
		t.Error("expected priority")
	}
	if rec.RiskReduction == "" {
		t.Error("expected risk reduction")
	}
}

func TestSDIPhase2RiskReduction(t *testing.T) {
	e := NewSDIEngine()
	input := makeSDIInput()
	res := e.Run(input)

	for _, r := range res.Recommendations {
		if r.RiskReduction == "" {
			t.Errorf("recommendation %s has no risk reduction", r.ID)
		}
		if len(r.AffectedFindings) > 0 && r.RiskReduction == "None" {
			t.Errorf("recommendation %s has findings but no risk reduction", r.ID)
		}
	}

	foundCritical := false
	for _, r := range res.Recommendations {
		if r.RiskReduction == "High" {
			foundCritical = true
			break
		}
	}
	if !foundCritical {
		t.Error("expected at least one High risk reduction recommendation")
	}
}

func TestSDIPhase3SecurityROI(t *testing.T) {
	tests := []struct {
		effort    string
		reduction string
		expected  string
	}{
		{"Low", "High", "Excellent"},
		{"Low", "Medium", "Good"},
		{"Medium", "High", "Good"},
		{"High", "High", "Fair"},
		{"High", "Low", "Limited"},
	}
	for _, tt := range tests {
		got := computeROI(tt.effort, tt.reduction)
		if got != tt.expected {
			t.Errorf("computeROI(%q, %q) = %q, want %q", tt.effort, tt.reduction, got, tt.expected)
		}
	}
}

func TestSDIPhase4Prioritization(t *testing.T) {
	e := NewSDIEngine()
	input := makeSDIInput()
	res := e.Run(input)

	if len(res.Recommendations) < 2 {
		t.Fatal("expected at least 2 recommendations for priority ordering")
	}

	for i := 0; i < len(res.Recommendations)-1; i++ {
		if priorityScore(res.Recommendations[i].Priority) < priorityScore(res.Recommendations[i+1].Priority) {
			t.Errorf("recommendations not sorted by priority: %s (%s) before %s (%s)",
				res.Recommendations[i].ID, res.Recommendations[i].Priority,
				res.Recommendations[i+1].ID, res.Recommendations[i+1].Priority)
		}
	}
}

func TestSDIPhase5FixSimulation(t *testing.T) {
	e := NewSDIEngine()
	input := makeSDIInput()
	res := e.Run(input)

	if len(res.FixSimulations) == 0 {
		t.Fatal("expected at least 1 fix simulation")
	}

	sim := res.FixSimulations[0]
	if sim.OriginalTotal == 0 {
		t.Error("expected original finding count")
	}
	if sim.ControlName == "" {
		t.Error("expected control name in fix simulation")
	}
	t.Logf("Fix simulation: %s: Critical %d->%d, High %d->%d",
		sim.ControlName, sim.OriginalCritical, sim.NewCritical,
		sim.OriginalHigh, sim.NewHigh)
}

func TestSDIPhase6FailureSimulation(t *testing.T) {
	e := NewSDIEngine()
	input := makeSDIInput()
	res := e.Run(input)

	if len(res.FailureSimulations) == 0 {
		t.Fatal("expected at least 1 failure simulation")
	}

	sim := res.FailureSimulations[0]
	if sim.ControlName == "" {
		t.Error("expected control name in failure simulation")
	}
	if sim.RiskIncrease == "" {
		t.Error("expected risk increase in failure simulation")
	}
	t.Logf("Failure simulation: %s: systems=%d, paths=%d, risk=%s",
		sim.ControlName, sim.SystemsImpacted, sim.AttackPathsOpened, sim.RiskIncrease)
}

func TestSDIPhase7ControlImpact(t *testing.T) {
	e := NewSDIEngine()
	input := makeSDIInput()
	res := e.Run(input)

	if len(res.ControlImpacts) == 0 {
		t.Fatal("expected at least 1 control impact")
	}

	ci := res.ControlImpacts[0]
	if ci.ControlName == "" {
		t.Error("expected control name in impact analysis")
	}
	if ci.SecurityValue == "" {
		t.Error("expected security value")
	}
	if ci.ROI == "" {
		t.Error("expected ROI")
	}
}

func TestSDIPhase8DecisionTrees(t *testing.T) {
	e := NewSDIEngine()
	input := makeSDIInput()
	res := e.Run(input)

	dt := res.DecisionTrees
	if dt.SingleAction.ActionCount != 1 {
		t.Error("expected 1 action in single decision tree")
	}
	if len(dt.SingleAction.RecommendedOrder) != 1 {
		t.Error("expected 1 recommendation in single tree")
	}
	if dt.ThreeActions.ActionCount != 3 {
		t.Error("expected 3 actions in three-action tree")
	}
	if dt.FiveActions.ActionCount != 5 {
		t.Error("expected 5 actions in five-action tree")
	}
}

func TestSDIPhase9BoardScenarios(t *testing.T) {
	e := NewSDIEngine()
	input := makeSDIInput()
	res := e.Run(input)

	bs := res.BoardScenarios
	if bs.DoNothing.Scenario != "Do Nothing" {
		t.Error("expected Do Nothing scenario")
	}
	if bs.PartialRemediate.Scenario != "Partial Remediation (Top 3 Actions)" {
		t.Error("expected Partial Remediation scenario")
	}
	if bs.FullRemediate.Scenario != "Full Remediation (All Actions)" {
		t.Error("expected Full Remediation scenario")
	}
	if bs.DoNothing.RiskScore <= 0 {
		t.Error("expected positive risk score")
	}
}

func TestSDIPhase10InvestmentPriorities(t *testing.T) {
	e := NewSDIEngine()
	input := makeSDIInput()
	res := e.Run(input)

	if len(res.InvestmentPriorities) == 0 {
		t.Fatal("expected at least 1 investment priority")
	}

	ip := res.InvestmentPriorities[0]
	if ip.Area == "" {
		t.Error("expected investment area")
	}
	if ip.Rank <= 0 {
		t.Error("expected positive rank")
	}
	if ip.Score <= 0 {
		t.Error("expected positive score")
	}
}

func TestSDIPhase11AttackPathCollapse(t *testing.T) {
	e := NewSDIEngine()
	input := makeSDIInput()
	res := e.Run(input)

	if len(res.AttackPathCollapse) == 0 {
		t.Fatal("expected at least 1 attack path collapse analysis")
	}

	apc := res.AttackPathCollapse[0]
	if apc.ControlName == "" {
		t.Error("expected control name in collapse analysis")
	}
	if apc.TotalAttackPaths <= 0 {
		t.Error("expected positive total attack paths")
	}
}

func TestSDIPhase12ComplianceImpact(t *testing.T) {
	e := NewSDIEngine()
	input := makeSDIInput()
	res := e.Run(input)

	if len(res.ComplianceImpacts) == 0 {
		t.Fatal("expected at least 1 compliance impact")
	}

	ci := res.ComplianceImpacts[0]
	if ci.Framework == "" {
		t.Error("expected framework in compliance impact")
	}
	if ci.Action == "" {
		t.Error("expected action in compliance impact")
	}
	if ci.Improvement == "" {
		t.Error("expected improvement level")
	}
}

func TestSDIPhase13RemediationRoadmap(t *testing.T) {
	e := NewSDIEngine()
	input := makeSDIInput()
	res := e.Run(input)

	rm := res.RemediationRoadmap
	if len(rm.Phase30) == 0 {
		t.Error("expected 30-day phase items")
	}
	if len(rm.Phase90) == 0 {
		t.Error("expected 90-day phase items")
	}
	if len(rm.Phase180) == 0 {
		t.Error("expected 180-day phase items")
	}

	for _, item := range rm.Phase30 {
		if item.Action == "" {
			t.Error("expected action name in roadmap item")
		}
		if item.Priority == "" {
			t.Error("expected priority in roadmap item")
		}
	}
}

func TestSDIPhase14DecisionDashboard(t *testing.T) {
	e := NewSDIEngine()
	input := makeSDIInput()
	res := e.Run(input)

	db := res.Dashboard
	if len(db.TopDecisions) == 0 {
		t.Error("expected top decisions")
	}
	if len(db.QuickWins) == 0 {
		t.Error("expected quick wins")
	}
	if db.RiskReductionSummary == "" {
		t.Error("expected risk reduction summary")
	}
	if db.TotalRiskReduction <= 0 {
		t.Error("expected positive total risk reduction")
	}
}

func TestSDIPhase15ExecutiveScenarios(t *testing.T) {
	e := NewSDIEngine()
	input := makeSDIInput()
	res := e.Run(input)

	es := res.ExecutiveScenarios
	if es.BestCase.Scenario != "Best Case" {
		t.Error("expected Best Case scenario")
	}
	if es.LikelyCase.Scenario != "Likely Case" {
		t.Error("expected Likely Case scenario")
	}
	if es.WorstCase.Scenario != "Worst Case" {
		t.Error("expected Worst Case scenario")
	}
	if es.BestCase.FindingsResolved <= 0 {
		t.Error("expected Best Case to resolve some findings")
	}
	if es.WorstCase.FindingsResolved != 0 {
		t.Error("expected Worst Case to resolve zero findings")
	}
	if es.BestCase.RiskScore > input.RiskScore {
		t.Error("expected Best Case risk score to be lower than current")
	}
	if es.WorstCase.RiskScore < input.RiskScore {
		t.Error("expected Worst Case risk score to be higher than current")
	}
}

func TestSDIEmptyInput(t *testing.T) {
	e := NewSDIEngine()
	input := SDIInput{}
	res := e.Run(input)
	if res == nil {
		t.Fatal("expected non-nil result from empty input")
	}
	if len(res.Recommendations) != 0 {
		t.Error("expected 0 recommendations for empty input")
	}
	if len(res.FixSimulations) != 0 {
		t.Error("expected 0 fix simulations for empty input")
	}
	if len(res.FailureSimulations) != 0 {
		t.Error("expected 0 failure simulations for empty input")
	}
	if len(res.ControlImpacts) != 0 {
		t.Error("expected 0 control impacts for empty input")
	}
}

func TestSDIAllRecommendationsHaveIDs(t *testing.T) {
	e := NewSDIEngine()
	input := makeSDIInput()
	res := e.Run(input)

	for _, r := range res.Recommendations {
		if r.ID == "" {
			t.Error("found recommendation with empty ID")
		}
		if r.Title == "" {
			t.Error("found recommendation with empty title")
		}
	}
}

func TestSDIFixSimulationReducesFindings(t *testing.T) {
	e := NewSDIEngine()
	input := makeSDIInput()
	res := e.Run(input)

	for _, sim := range res.FixSimulations {
		if sim.NewTotal > sim.OriginalTotal {
			t.Errorf("fix simulation for %s increased findings: %d -> %d", sim.ControlName, sim.OriginalTotal, sim.NewTotal)
		}
		if sim.NewCritical > sim.OriginalCritical {
			t.Errorf("fix simulation for %s increased critical findings: %d -> %d", sim.ControlName, sim.OriginalCritical, sim.NewCritical)
		}
	}
}

func TestSDIRoadmapPhasesOrderedByPriority(t *testing.T) {
	e := NewSDIEngine()
	input := makeSDIInput()
	res := e.Run(input)
	rm := res.RemediationRoadmap

	for _, phase := range [][]SDIRoadmapItem{rm.Phase30, rm.Phase90, rm.Phase180} {
		for _, item := range phase {
			if item.Priority == "" {
				t.Error("found roadmap item with empty priority")
			}
		}
	}
}

func TestSDIRecommendationConsistency(t *testing.T) {
	e := NewSDIEngine()
	input := makeSDIInput()
	res := e.Run(input)

	for _, r := range res.Recommendations {
		affectedCount := len(r.AffectedFindings) + len(r.AffectedThreats) + len(r.AffectedAttackPaths)
		if affectedCount == 0 {
			t.Errorf("recommendation %s has no affected items", r.ID)
		}
	}
}

func TestSDIDashboardConsistency(t *testing.T) {
	e := NewSDIEngine()
	input := makeSDIInput()
	res := e.Run(input)

	db := res.Dashboard
	for _, qw := range db.QuickWins {
		if qw.Effort != "Low" {
			t.Errorf("quick win %s has effort %s, expected Low", qw.ID, qw.Effort)
		}
		if qw.Priority != "Critical" && qw.Priority != "High" {
			t.Errorf("quick win %s has priority %s, expected Critical or High", qw.ID, qw.Priority)
		}
	}
}
