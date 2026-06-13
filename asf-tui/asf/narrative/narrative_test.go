package narrative

import (
	"strings"
	"testing"
)

func TestNarrativeEngineGenerateNarrative(t *testing.T) {
	engine := NewNarrativeEngine("healthcare", []string{"Auth0", "API", "Database"}, []string{"Auth0->API", "API->Database"})

	assumptions := []Assumption{
		{
			ID:                  "A1",
			Description:         "Admin access is restricted to authorized personnel",
			Component:           "Auth0",
			Category:            "access",
			Risk:                "Critical",
			STRIDECategories:    []string{"Elevation of Privilege", "Spoofing"},
			Likelihood:          3,
			Impact:              4,
			Confidence:          0.85,
			Keywords:            []string{"admin", "access", "restricted"},
			SourceComponents:    []string{"Auth0"},
			SourceRelationships: []string{"Auth0->API"},
			Rationale:           "Auth0 provides centralized authentication for multiple downstream services",
			EvidenceSources:     []string{"architecture.yaml"},
		},
		{
			ID:                  "A2",
			Description:         "MFA is enforced for all administrative accounts",
			Component:           "Auth0",
			Category:            "identity",
			Risk:                "High",
			STRIDECategories:    []string{"Spoofing"},
			Likelihood:          3,
			Impact:              3,
			Confidence:          0.90,
			Keywords:            []string{"MFA", "multi-factor", "admin"},
			SourceComponents:    []string{"Auth0"},
			SourceRelationships: []string{"Auth0->API", "User->Auth0"},
			Rationale:           "The architecture specifies MFA for admin accounts",
			EvidenceSources:     []string{"security_controls.yaml"},
		},
		{
			ID:                  "A3",
			Description:         "API traffic is encrypted in transit",
			Component:           "API",
			Category:            "network",
			Risk:                "High",
			STRIDECategories:    []string{"Information Disclosure", "Tampering"},
			Likelihood:          2,
			Impact:              3,
			Confidence:          0.75,
			Keywords:            []string{"encryption", "TLS", "API"},
			SourceComponents:    []string{"API"},
			SourceRelationships: []string{"Auth0->API", "API->Database"},
			Rationale:           "TLS is specified for all API traffic",
			EvidenceSources:     []string{"architecture.yaml"},
		},
	}

	controls := []ControlDetail{
		{
			Name:                 "MFA",
			Category:             "identity",
			Description:          "Multi-factor authentication",
			Rationale:            "Prevents credential compromise",
			MitigatedAssumptions: []string{"A2"},
			STRIDECategories:     []string{"Spoofing"},
		},
	}

	trustBoundaries := []TrustBoundary{
		{
			Type:        "network",
			Components:  []string{"Auth0", "API"},
			RiskLevel:   "High",
			Description: "Authentication boundary between Auth0 and internal services",
		},
	}

	contradictions := []Contradiction{
		{
			ID:                  "C1",
			Severity:            "High",
			Description:         "Admin access requires MFA but break-glass bypasses MFA",
			Explanation:         "The break-glass procedure contradicts the MFA requirement",
			AffectedAssumptions: []string{"A2"},
		},
	}

	strideDist := map[string]int{"Spoofing": 2, "Elevation of Privilege": 1, "Information Disclosure": 1, "Tampering": 1}
	riskDist := map[string]int{"Critical": 1, "High": 2, "Medium": 0, "Low": 0}

	output := engine.GenerateNarrative(
		"Healthcare API",
		assumptions,
		controls,
		trustBoundaries,
		contradictions,
		"healthcare",
		strideDist,
		riskDist,
	)

	if output == nil {
		t.Fatal("expected non-nil output")
	}

	// Verify architecture overview
	if output.ArchitectureOverview.Name != "Healthcare API" {
		t.Errorf("expected name 'Healthcare API', got %s", output.ArchitectureOverview.Name)
	}
	if output.ArchitectureOverview.Domain != "healthcare" {
		t.Errorf("expected domain 'healthcare', got %s", output.ArchitectureOverview.Domain)
	}
	if output.ArchitectureOverview.TotalAssumptions != 3 {
		t.Errorf("expected 3 assumptions, got %d", output.ArchitectureOverview.TotalAssumptions)
	}
	if output.ArchitectureOverview.CriticalCount != 1 {
		t.Errorf("expected 1 critical, got %d", output.ArchitectureOverview.CriticalCount)
	}
	if output.ArchitectureOverview.HighCount != 2 {
		t.Errorf("expected 2 high, got %d", output.ArchitectureOverview.HighCount)
	}
	if len(output.ArchitectureOverview.KeyComponents) != 2 {
		t.Errorf("expected 2 key components, got %d", len(output.ArchitectureOverview.KeyComponents))
	}

	// Verify assumption narratives
	if len(output.AssumptionNarratives) != 3 {
		t.Fatalf("expected 3 narratives, got %d", len(output.AssumptionNarratives))
	}

	// Check first assumption (Critical)
	n1 := output.AssumptionNarratives[0]
	if n1.AssumptionID != "A1" {
		t.Errorf("expected A1, got %s", n1.AssumptionID)
	}
	if n1.RiskLevel != "Critical" {
		t.Errorf("expected Critical, got %s", n1.RiskLevel)
	}
	if n1.Context == "" {
		t.Error("expected non-empty context")
	}
	if n1.WhyASFIdentifiedIt == "" {
		t.Error("expected non-empty why identified")
	}
	if n1.ArchitecturalImportance == "" {
		t.Error("expected non-empty architectural importance")
	}
	if n1.FailureConsequence == "" {
		t.Error("expected non-empty failure consequence")
	}
	if n1.SecurityRecommendation == "" {
		t.Error("expected non-empty security recommendation")
	}

	// Verify executive report
	if output.ExecutiveReport.ArchitectureOverview == "" {
		t.Error("expected non-empty executive overview")
	}
	if len(output.ExecutiveReport.MostCriticalAssumptions) == 0 {
		t.Error("expected critical assumptions")
	}
	if len(output.ExecutiveReport.ArchitecturalConcerns) == 0 {
		t.Error("expected architectural concerns from contradictions")
	}
	if len(output.ExecutiveReport.RecommendedInvestments) == 0 {
		t.Error("expected recommended investments")
	}

	// Verify technical summary
	if output.TechnicalSummary.ArchitectureSummary == "" {
		t.Error("expected non-empty technical summary")
	}
	if len(output.TechnicalSummary.AssumptionDetails) != 3 {
		t.Errorf("expected 3 technical details, got %d", len(output.TechnicalSummary.AssumptionDetails))
	}
	if len(output.TechnicalSummary.Recommendations) == 0 {
		t.Error("expected recommendations")
	}
	if len(output.TechnicalSummary.Dependencies) == 0 {
		t.Error("expected dependencies")
	}

	// Verify architect narrative
	if output.ArchitectNarrative == "" {
		t.Error("expected non-empty architect narrative")
	}
	if !strings.Contains(output.ArchitectNarrative, "Healthcare API") {
		t.Error("expected narrative to contain architecture name")
	}
	if !strings.Contains(output.ArchitectNarrative, "Context") {
		t.Error("expected narrative to contain Context section")
	}
	if !strings.Contains(output.ArchitectNarrative, "Why ASF Identified This") {
		t.Error("expected narrative to contain Why ASF Identified This section")
	}
	if !strings.Contains(output.ArchitectNarrative, "Architectural Importance") {
		t.Error("expected narrative to contain Architectural Importance section")
	}
	if !strings.Contains(output.ArchitectNarrative, "Failure Consequence") {
		t.Error("expected narrative to contain Failure Consequence section")
	}
	if !strings.Contains(output.ArchitectNarrative, "Security Recommendation") {
		t.Error("expected narrative to contain Security Recommendation section")
	}
}

func TestAssumptionNarrativeSections(t *testing.T) {
	engine := NewNarrativeEngine("fintech", []string{"PaymentGateway"}, []string{})

	assumptions := []Assumption{
		{
			ID:               "F1",
			Description:      "Payment data is encrypted using AES-256",
			Component:        "PaymentGateway",
			Category:         "network",
			Risk:             "High",
			STRIDECategories: []string{"Information Disclosure"},
			Likelihood:       2,
			Impact:           4,
			Confidence:       0.80,
			Keywords:         []string{"encryption", "AES-256", "payment"},
			SourceComponents: []string{"PaymentGateway"},
			Rationale:        "The architecture specifies encryption for payment data",
			EvidenceSources:  []string{"security.yaml"},
		},
	}

	output := engine.GenerateNarrative(
		"Fintech Platform",
		assumptions,
		[]ControlDetail{},
		[]TrustBoundary{},
		[]Contradiction{},
		"fintech",
		map[string]int{},
		map[string]int{},
	)

	n := output.AssumptionNarratives[0]

	// Context should mention the component
	if !strings.Contains(n.Context, "PaymentGateway") {
		t.Errorf("expected context to mention PaymentGateway, got: %s", n.Context)
	}

	// Why identified should use rationale
	if !strings.Contains(n.WhyASFIdentifiedIt, "architecture specifies") {
		t.Errorf("expected why identified to reference rationale, got: %s", n.WhyASFIdentifiedIt)
	}

	// Architectural importance should mention risk
	if !strings.Contains(n.ArchitecturalImportance, "significant") {
		t.Errorf("expected architectural importance to mention significance, got: %s", n.ArchitecturalImportance)
	}

	// Failure consequence should mention encryption
	if !strings.Contains(n.FailureConsequence, "encryption") {
		t.Errorf("expected failure consequence to mention encryption, got: %s", n.FailureConsequence)
	}

	// Recommendation should mention encryption
	if !strings.Contains(n.SecurityRecommendation, "encryption") {
		t.Errorf("expected recommendation to mention encryption, got: %s", n.SecurityRecommendation)
	}
}

func TestStyleEnforcement(t *testing.T) {
	engine := NewNarrativeEngine("cloud", []string{}, []string{})

	tests := []struct {
		input  string
		banned string
	}{
		{"Leverage multi-factor authentication", "leverage"},
		{"This is a robust solution", "robust"},
		{"Ensure that all access is restricted", "ensure that"},
		{"It is recommended to implement MFA", "it is recommended"},
		{"Consider implementing WAF", "consider implementing"},
		{"This is cutting-edge technology", "cutting-edge"},
	}

	for _, tc := range tests {
		result := engine.enforceStyle(tc.input)
		if strings.Contains(strings.ToLower(result), tc.banned) {
			t.Errorf("style enforcement failed for '%s': still contains '%s' in result '%s'", tc.input, tc.banned, result)
		}
	}
}

func TestInferControls(t *testing.T) {
	engine := NewNarrativeEngine("cloud", []string{}, []string{})

	tests := []struct {
		desc     string
		expected string
	}{
		{"MFA is enforced", "multi-factor authentication"},
		{"Data is encrypted", "encryption at rest and in transit"},
		{"Access is restricted", "role-based access control"},
		{"Logging is enabled", "centralized logging and monitoring"},
		{"Backups are taken", "automated backup and recovery"},
		{"WAF is deployed", "web application firewall"},
		{"Rate limiting is applied", "rate limiting"},
		{"Secrets are managed", "secrets management"},
	}

	for _, tc := range tests {
		a := Assumption{Description: tc.desc}
		controls := engine.inferControls(a)
		found := false
		for _, c := range controls {
			if strings.Contains(c, tc.expected) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected control containing '%s' for '%s', got: %v", tc.expected, tc.desc, controls)
		}
	}
}

func TestInferConsequence(t *testing.T) {
	engine := NewNarrativeEngine("cloud", []string{}, []string{})

	tests := []struct {
		desc     string
		expected string
	}{
		{"MFA is enforced", "Without MFA"},
		{"Data is encrypted", "Without encryption"},
		{"Access is restricted", "Without access restrictions"},
		{"Logging is enabled", "Without logging"},
		{"Backups are taken", "Without backups"},
	}

	for _, tc := range tests {
		a := Assumption{Description: tc.desc}
		consequence := engine.inferConsequence(a)
		if !strings.Contains(consequence, tc.expected) {
			t.Errorf("expected consequence containing '%s' for '%s', got: %s", tc.expected, tc.desc, consequence)
		}
	}
}

func TestDependencyMapping(t *testing.T) {
	engine := NewNarrativeEngine("cloud", []string{"Auth0", "API", "DB"}, []string{"Auth0->API", "API->DB"})

	assumptions := []Assumption{
		{
			ID:                  "D1",
			Description:         "Auth0 manages authentication",
			Component:           "Auth0",
			SourceRelationships: []string{"Auth0->API"},
		},
		{
			ID:                  "D2",
			Description:         "API handles requests",
			Component:           "API",
			SourceRelationships: []string{"Auth0->API", "API->DB"},
		},
		{
			ID:                  "D3",
			Description:         "DB stores data",
			Component:           "DB",
			SourceRelationships: []string{"API->DB"},
		},
	}

	depMap := engine.buildDependencyMap(assumptions)

	// D1 and D2 share Auth0->API relationship
	if len(depMap["D1"]) == 0 {
		t.Error("expected D1 to have downstream dependencies")
	}
	// D2 and D3 share API->DB relationship
	if len(depMap["D2"]) == 0 {
		t.Error("expected D2 to have downstream dependencies")
	}
}

func TestExportMarkdown(t *testing.T) {
	engine := NewNarrativeEngine("cloud", []string{"Auth0"}, []string{})
	assumptions := []Assumption{
		{
			ID:               "E1",
			Description:      "Admin access is restricted",
			Component:        "Auth0",
			Risk:             "Critical",
			Confidence:       0.90,
			STRIDECategories: []string{"Elevation of Privilege"},
		},
	}

	output := engine.GenerateNarrative(
		"Test",
		assumptions,
		[]ControlDetail{},
		[]TrustBoundary{},
		[]Contradiction{},
		"cloud",
		map[string]int{},
		map[string]int{},
	)

	md := ExportMarkdown(output)

	if !strings.Contains(md, "# Security Architect Narrative") {
		t.Error("expected markdown to contain title")
	}
	if !strings.Contains(md, "## Executive Summary") {
		t.Error("expected markdown to contain Executive Summary")
	}
	if !strings.Contains(md, "## Technical Summary") {
		t.Error("expected markdown to contain Technical Summary")
	}
	if !strings.Contains(md, "## Architect Narrative") {
		t.Error("expected markdown to contain Architect Narrative")
	}
	if !strings.Contains(md, "Admin access is restricted") {
		t.Error("expected markdown to contain assumption description")
	}
}

func TestExportHTML(t *testing.T) {
	engine := NewNarrativeEngine("cloud", []string{"Auth0"}, []string{})
	assumptions := []Assumption{
		{
			ID:               "E1",
			Description:      "Admin access is restricted",
			Component:        "Auth0",
			Risk:             "Critical",
			Confidence:       0.90,
			STRIDECategories: []string{"Elevation of Privilege"},
		},
	}

	output := engine.GenerateNarrative(
		"Test",
		assumptions,
		[]ControlDetail{},
		[]TrustBoundary{},
		[]Contradiction{},
		"cloud",
		map[string]int{},
		map[string]int{},
	)

	html := ExportHTML(output)

	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("expected HTML to contain DOCTYPE")
	}
	if !strings.Contains(html, "Security Architect Narrative") {
		t.Error("expected HTML to contain title")
	}
	if !strings.Contains(html, "Executive Summary") {
		t.Error("expected HTML to contain Executive Summary")
	}
	if !strings.Contains(html, "badge-critical") {
		t.Error("expected HTML to contain critical badge")
	}
}

func TestExecutiveReport(t *testing.T) {
	engine := NewNarrativeEngine("cloud", []string{"Auth0", "API", "DB"}, []string{})
	assumptions := []Assumption{
		{
			ID:          "X1",
			Description: "Auth0 is the single source of authentication",
			Component:   "Auth0",
			Risk:        "Critical",
			Category:    "identity",
		},
		{
			ID:          "X2",
			Description: "API validates all tokens",
			Component:   "API",
			Risk:        "High",
			Category:    "access",
		},
		{
			ID:          "X3",
			Description: "Database has no encryption",
			Component:   "DB",
			Risk:        "Critical",
			Category:    "network",
		},
	}

	output := engine.GenerateNarrative(
		"Test",
		assumptions,
		[]ControlDetail{},
		[]TrustBoundary{},
		[]Contradiction{},
		"cloud",
		map[string]int{},
		map[string]int{},
	)

	report := output.ExecutiveReport

	// Should identify critical assumptions
	if len(report.MostCriticalAssumptions) != 3 {
		t.Errorf("expected 3 most critical, got %d", len(report.MostCriticalAssumptions))
	}

	// Should identify single points of failure
	found := false
	for _, s := range report.SinglePointsOfFailure {
		if strings.Contains(s, "Auth0") {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected Auth0 to be identified as single point of failure, got: %v", report.SinglePointsOfFailure)
	}

	// Should generate recommendations
	if len(report.RecommendedInvestments) == 0 {
		t.Error("expected recommended investments")
	}
}

func TestBannedPhrases(t *testing.T) {
	if len(BannedPhrases) == 0 {
		t.Error("expected banned phrases to be defined")
	}

	// Verify all banned phrases have reasons
	for _, rule := range BannedPhrases {
		if rule.Pattern == "" {
			t.Error("expected pattern to be non-empty")
		}
		if rule.Reason == "" {
			t.Errorf("expected reason for pattern '%s'", rule.Pattern)
		}
	}
}

func TestInferBusinessImpact(t *testing.T) {
	engine := NewNarrativeEngine("cloud", []string{}, []string{})

	tests := []struct {
		desc     string
		contains string
	}{
		{"MFA is enforced", "Unauthorized access"},
		{"Data is encrypted", "Data breach"},
		{"Access is restricted", "Unauthorized data"},
		{"Logging is enabled", "Inability to detect"},
		{"Backups are taken", "Permanent data loss"},
		{"API validates tokens", "API abuse"},
		{"HIPAA compliance", "Regulatory fines"},
	}

	for _, tc := range tests {
		a := Assumption{Description: tc.desc}
		impact := engine.inferBusinessImpact(a)
		if !strings.Contains(impact, tc.contains) {
			t.Errorf("expected business impact containing '%s' for '%s', got: %s", tc.contains, tc.desc, impact)
		}
	}
}

func TestInferRole(t *testing.T) {
	engine := NewNarrativeEngine("cloud", []string{}, []string{})

	tests := []struct {
		desc     string
		contains string
	}{
		{"authentication is handled", "authentication"},
		{"access control", "access control"},
		{"encryption is used", "data protection"},
		{"logging is enabled", "observability"},
		{"network security", "network security"},
		{"database storage", "data storage"},
		{"API gateway", "API management"},
	}

	for _, tc := range tests {
		a := Assumption{Description: tc.desc}
		role := engine.inferRole(a)
		if !strings.Contains(role, tc.contains) {
			t.Errorf("expected role containing '%s' for '%s', got: %s", tc.contains, tc.desc, role)
		}
	}
}

func TestInferEffort(t *testing.T) {
	engine := NewNarrativeEngine("cloud", []string{}, []string{})

	tests := []struct {
		control  string
		contains string
	}{
		{"multi-factor authentication", "Low"},
		{"encryption", "Medium"},
		{"network segmentation", "High"},
		{"backup", "Medium"},
		{"audit logging", "Low"},
		{"WAF", "Medium"},
	}

	for _, tc := range tests {
		a := Assumption{Description: "Test"}
		effort := engine.inferEffort(a, tc.control)
		if !strings.Contains(effort, tc.contains) {
			t.Errorf("expected effort containing '%s' for '%s', got: %s", tc.contains, tc.control, effort)
		}
	}
}

func TestIsSinglePointOfFailure(t *testing.T) {
	spof := Assumption{Description: "Auth0 is the single source of authentication"}
	if !isSinglePointOfFailure(spof) {
		t.Error("expected single point of failure for 'single source'")
	}

	notSpof := Assumption{Description: "MFA is enabled"}
	if isSinglePointOfFailure(notSpof) {
		t.Error("expected not single point of failure for 'MFA is enabled'")
	}
}

func TestEmptyAssumptions(t *testing.T) {
	engine := NewNarrativeEngine("cloud", []string{}, []string{})
	output := engine.GenerateNarrative(
		"Empty",
		[]Assumption{},
		[]ControlDetail{},
		[]TrustBoundary{},
		[]Contradiction{},
		"cloud",
		map[string]int{},
		map[string]int{},
	)

	if len(output.AssumptionNarratives) != 0 {
		t.Errorf("expected 0 narratives, got %d", len(output.AssumptionNarratives))
	}
	if output.ArchitectureOverview.TotalAssumptions != 0 {
		t.Errorf("expected 0 total assumptions, got %d", output.ArchitectureOverview.TotalAssumptions)
	}
}
