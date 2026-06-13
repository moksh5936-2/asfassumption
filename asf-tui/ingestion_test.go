package main

import (
	"os"
	"strings"
	"testing"

	"asf-tui/asf/models"
)

// ──────────────────────────────────────────────
// Explicit Assumption Classification Tests
// ──────────────────────────────────────────────

func TestClassifyExplicitAssumption(t *testing.T) {
	tests := []struct {
		text string
		want models.AssumptionType
	}{
		{"MFA is enforced for all Auth0 user authentication", models.AssumptionTypeIDENTITY},
		{"Auth0 administrative access is restricted to authorized admins only", models.AssumptionTypeIDENTITY},
		{"All API requests pass through APIGateway for authentication validation", models.AssumptionTypeIDENTITY},
		{"PHI data is encrypted at rest in PHIDatabase", models.AssumptionTypeCONFIGURATION},
		{"PHI data is encrypted in transit between all components", models.AssumptionTypeCONFIGURATION},
		{"Encryption keys are stored and managed in KMS with automatic rotation", models.AssumptionTypeCONFIGURATION},
		{"KMS access is restricted to authorized services only", models.AssumptionTypeACCESS},
		{"Audit logging is immutable and tamper-proof", models.AssumptionTypeCONFIGURATION},
		{"All PHI access events are logged with user and timestamp", models.AssumptionTypeACCESS},
		{"Backup data is encrypted at rest and in transit", models.AssumptionTypeCONFIGURATION},
		{"Backup restore procedures are tested regularly", models.AssumptionTypeCONFIGURATION},
		{"ThirdPartyAnalytics has access only to de-identified PHI", models.AssumptionTypeACCESS},
		{"Third-party provider maintains equivalent security controls", models.AssumptionTypeDEPENDENCY},
		{"AdminConsole requires MFA for all administrative access", models.AssumptionTypeIDENTITY},
		{"Object-level authorization is enforced for PHI record access", models.AssumptionTypeACCESS},
		{"Database connection pooling does not leak data between sessions", models.AssumptionTypeIDENTITY},
		{"TLS certificates are monitored and renewed before expiry", models.AssumptionTypeCONFIGURATION},
		{"API rate limiting prevents abuse of PHI endpoints", models.AssumptionTypeGOVERNANCE},
		{"Network segmentation isolates the PHI database in a private subnet", models.AssumptionTypeNETWORK},
		{"Unauthorized PHI export is detected and alerted", models.AssumptionTypeACCESS},
		{"Session tokens expire and are rotated periodically", models.AssumptionTypeIDENTITY},
		{"Auth0 tenant is configured with breach detection and anomaly alerts", models.AssumptionTypeIDENTITY},
		{"Audit log storage has sufficient retention for compliance", models.AssumptionTypeCONFIGURATION},
		{"PHI data minimization policies are enforced at application layer", models.AssumptionTypeGOVERNANCE},
		{"Incident response plan includes PHI breach notification", models.AssumptionTypePROCESS},
		{"Vendor risk assessments are conducted for ThirdPartyAnalytics", models.AssumptionTypeDEPENDENCY},
		{"Database backups are stored in a separate geographic region", models.AssumptionTypeCONFIGURATION},
		{"KMS key deletion is protected with multi-factor authorization", models.AssumptionTypeIDENTITY},
		{"API Gateway logs are monitored for anomalous access patterns", models.AssumptionTypeACCESS},
		{"System health and availability monitoring covers all components", models.AssumptionTypeCONFIGURATION},
	}

	for _, tt := range tests {
		got := classifyExplicitAssumption(tt.text)
		if got != tt.want {
			t.Errorf("classifyExplicitAssumption(%q) = %s, want %s", tt.text, got, tt.want)
		}
	}
}

func TestClassifyExplicitAssumptionEmpty(t *testing.T) {
	got := classifyExplicitAssumption("")
	if got != models.AssumptionTypeGOVERNANCE {
		t.Errorf("empty text should return GOVERNANCE, got %s", got)
	}
}

func TestClassifyExplicitAssumptionDeterministic(t *testing.T) {
	text := "MFA is enforced for all Auth0 user authentication"
	r1 := classifyExplicitAssumption(text)
	for i := 0; i < 50; i++ {
		r2 := classifyExplicitAssumption(text)
		if r1 != r2 {
			t.Fatalf("classification not deterministic: %s != %s", r1, r2)
		}
	}
}

// ──────────────────────────────────────────────
// Explicit Risk Assessment Tests
// ──────────────────────────────────────────────

func TestAssessExplicitRisk(t *testing.T) {
	e := &Engine{
		strideEngine: NewStrideEngine(),
	}

	riskOrder := map[RiskLevel]int{
		RiskLow:      0,
		RiskMedium:   1,
		RiskHigh:     2,
		RiskCritical: 3,
	}

	tests := []struct {
		text string
		typ  models.AssumptionType
		min  RiskLevel
		max  RiskLevel
	}{
		{"PHI data is encrypted at rest", models.AssumptionTypeCONFIGURATION, RiskMedium, RiskCritical},
		{"MFA is enforced", models.AssumptionTypeIDENTITY, RiskLow, RiskHigh},
		{"Backup restore procedures are tested", models.AssumptionTypeCONFIGURATION, RiskLow, RiskMedium},
		{"Incident response plan includes PHI breach notification", models.AssumptionTypePROCESS, RiskHigh, RiskCritical},
		{"System health monitoring covers all components", models.AssumptionTypeGOVERNANCE, RiskMedium, RiskCritical},
		{"Unauthorized PHI export is detected and alerted", models.AssumptionTypeCONFIGURATION, RiskMedium, RiskCritical},
		{"Third-party provider maintains security controls", models.AssumptionTypeDEPENDENCY, RiskLow, RiskMedium},
	}

	for _, tt := range tests {
		got := e.assessExplicitRisk(tt.text, tt.typ)
		gotVal := riskOrder[got]
		minVal := riskOrder[tt.min]
		maxVal := riskOrder[tt.max]
		if gotVal < minVal || gotVal > maxVal {
			t.Errorf("assessExplicitRisk(%q, %s) = %s (val %d), want between %s (val %d) and %s (val %d)",
				tt.text, tt.typ, got, gotVal, tt.min, minVal, tt.max, maxVal)
		}
	}
}

func TestAssessExplicitRiskPHI(t *testing.T) {
	e := &Engine{strideEngine: NewStrideEngine()}
	riskOrder := map[RiskLevel]int{
		RiskLow: 0, RiskMedium: 1, RiskHigh: 2, RiskCritical: 3,
	}
	phiRisk := e.assessExplicitRisk("PHI data is encrypted at rest", models.AssumptionTypeCONFIGURATION)
	genRisk := e.assessExplicitRisk("System is monitored", models.AssumptionTypeGOVERNANCE)
	if riskOrder[phiRisk] < riskOrder[genRisk] {
		t.Errorf("PHI risk (%s) should be >= general risk (%s)", phiRisk, genRisk)
	}
}

func TestAssessExplicitRiskScale(t *testing.T) {
	e := &Engine{strideEngine: NewStrideEngine()}
	tests := []struct {
		text string
		typ  models.AssumptionType
	}{
		{"", models.AssumptionTypeGOVERNANCE},
		{"simple text", models.AssumptionTypeGOVERNANCE},
		{"encryption enabled", models.AssumptionTypeCONFIGURATION},
		{"PHI breach response", models.AssumptionTypePROCESS},
		{"critical PHI data access and encryption with KMS", models.AssumptionTypeACCESS},
	}
	for _, tt := range tests {
		got := e.assessExplicitRisk(tt.text, tt.typ)
		if got != RiskCritical && got != RiskHigh && got != RiskMedium && got != RiskLow {
			t.Errorf("assessExplicitRisk returned invalid risk level: %s", got)
		}
	}
}

// ──────────────────────────────────────────────
// Compliance Output Tests
// ──────────────────────────────────────────────

func TestBuildComplianceOutput(t *testing.T) {
	t.Run("with compliance frameworks", func(t *testing.T) {
		e := &Engine{
			archDesc: &ArchDescription{
				Compliance: []string{"HIPAA", "SOC2", "ISO27001"},
			},
		}
		output := e.buildComplianceOutput()
		if len(output) == 0 {
			t.Fatal("expected non-empty compliance output")
		}
		if !strings.Contains(output[0], "Compliance frameworks") {
			t.Errorf("expected header, got: %s", output[0])
		}
		foundHIPAA := false
		for _, line := range output {
			if strings.Contains(line, "HIPAA") {
				foundHIPAA = true
			}
		}
		if !foundHIPAA {
			t.Error("expected HIPAA in compliance output")
		}
	})

	t.Run("no compliance frameworks", func(t *testing.T) {
		e := &Engine{
			archDesc: &ArchDescription{},
		}
		output := e.buildComplianceOutput()
		if len(output) == 0 {
			t.Fatal("expected output even without compliance")
		}
		if !strings.Contains(output[0], "see gap analysis") {
			t.Errorf("expected default message, got: %s", output[0])
		}
	})

	t.Run("nil archDesc", func(t *testing.T) {
		e := &Engine{archDesc: nil}
		output := e.buildComplianceOutput()
		if len(output) == 0 {
			t.Fatal("expected output even with nil archDesc")
		}
	})
}

// ──────────────────────────────────────────────
// Deduplication / normalizeText Tests
// ──────────────────────────────────────────────

func TestNormalizeText(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"MFA is enforced.", "mfa is enforced"},
		{"MFA is enforced", "mfa is enforced"},
		{"  MFA   is  enforced.  ", "mfa is enforced"},
		{"", ""},
		{".", ""},
	}

	for _, tt := range tests {
		got := normalizeText(tt.input)
		if got != tt.want {
			t.Errorf("normalizeText(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestNormalizeTextDedupEquivalent(t *testing.T) {
	a := normalizeText("PHI data is encrypted at rest in PHIDatabase.")
	b := normalizeText("PHI data is encrypted at rest in PHIDatabase")
	if a != b {
		t.Errorf("equivalent texts should normalize identically: %q vs %q", a, b)
	}
}

// ──────────────────────────────────────────────
// Keyword Extraction Tests
// ──────────────────────────────────────────────

func TestExtractKeywords(t *testing.T) {
	tests := []struct {
		text string
		min  int
	}{
		{"MFA is enforced for all Auth0 user authentication", 3},
		{"PHI data is encrypted at rest", 2},
		{"", 0},
		{"the and for are but not you all", 0},
	}

	for _, tt := range tests {
		got := extractKeywords(tt.text)
		if len(got) < tt.min {
			t.Errorf("extractKeywords(%q) = %v, expected >= %d keywords", tt.text, got, tt.min)
		}
		for _, kw := range got {
			if len(kw) < 3 {
				t.Errorf("keyword %q too short in %q", kw, tt.text)
			}
		}
	}
}

func TestExtractKeywordsNoStopwords(t *testing.T) {
	got := extractKeywords("the and for are but not you all")
	if len(got) != 0 {
		t.Errorf("expected no keywords from stopwords, got %v", got)
	}
}

// ──────────────────────────────────────────────
// Security Controls Enhancement Tests
// ──────────────────────────────────────────────

func TestEnhanceControlsWithSecurityControls(t *testing.T) {
	existing := []ControlDetail{
		{ID: "CTRL-001", Description: "Enforce multi-factor authentication for all access", Category: "AUTHENTICATION", Priority: 1},
		{ID: "CTRL-002", Description: "Implement network segmentation and encryption in transit", Category: "NETWORK", Priority: 1},
	}

	securityControls := map[string][]string{
		"authentication": {"MFA", "Password_Policy", "Session_Management"},
		"encryption":     {"KMS_Key_Management", "TLS_1.3"},
		"backup":         {"Encrypted_Backups", "Cross_Region_Replication"},
	}

	enhanced := enhanceControlsWithSecurityControls(existing, securityControls)

	if len(enhanced) < len(existing) {
		t.Errorf("expected at least %d controls, got %d", len(existing), len(enhanced))
	}

	foundAuth := false
	for _, c := range enhanced {
		if c.Category == "AUTHENTICATION" && strings.Contains(c.Description, "MFA") {
			foundAuth = true
			break
		}
	}
	if !foundAuth {
		t.Error("expected AUTHENTICATION control to be enriched with MFA")
	}

	foundEnc := false
	for _, c := range enhanced {
		if strings.Contains(c.Description, "encryption") {
			foundEnc = true
			break
		}
	}
	if !foundEnc {
		t.Error("expected encryption control to be added")
	}

	foundBackup := false
	for _, c := range enhanced {
		if strings.Contains(c.Description, "backup") || strings.Contains(c.Description, "Backup") {
			foundBackup = true
			break
		}
	}
	if !foundBackup {
		t.Error("expected backup control to be added")
	}
}

func TestEnhanceControlsEmpty(t *testing.T) {
	existing := []ControlDetail{{ID: "CTRL-001", Category: "AUTHENTICATION", Description: "test"}}
	result := enhanceControlsWithSecurityControls(existing, nil)
	if len(result) != len(existing) {
		t.Errorf("expected unchanged controls with nil security controls")
	}

	result2 := enhanceControlsWithSecurityControls(existing, map[string][]string{})
	if len(result2) != len(existing) {
		t.Errorf("expected unchanged controls with empty security controls")
	}
}

func TestEnhanceControlsUnknownCategory(t *testing.T) {
	existing := []ControlDetail{}
	controls := map[string][]string{
		"unknown_category": {"Some_Control"},
	}
	result := enhanceControlsWithSecurityControls(existing, controls)
	if len(result) != 0 {
		t.Errorf("expected no controls for unknown category, got %d", len(result))
	}
}

// ──────────────────────────────────────────────
// toFloat Helper Tests
// ──────────────────────────────────────────────

func TestToFloat(t *testing.T) {
	tests := []struct {
		input interface{}
		want  float64
		ok    bool
	}{
		{float64(25), 25, true},
		{int(10), 10, true},
		{int64(100), 100, true},
		{uint64(50), 50, true},
		{"string", 0, false},
		{nil, 0, false},
		{[]int{1, 2, 3}, 0, false},
	}

	for _, tt := range tests {
		got, ok := toFloat(tt.input)
		if ok != tt.ok {
			t.Errorf("toFloat(%v) ok = %v, want %v", tt.input, ok, tt.ok)
		}
		if ok && got != tt.want {
			t.Errorf("toFloat(%v) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestToFloatInt32(t *testing.T) {
	got, ok := toFloat(int32(42))
	if !ok || got != 42 {
		t.Errorf("int32: got %v, %v", got, ok)
	}
}

func TestToFloatUint(t *testing.T) {
	got, ok := toFloat(uint(99))
	if !ok || got != 99 {
		t.Errorf("uint: got %v, %v", got, ok)
	}
}

// ──────────────────────────────────────────────
// buildValidationSummary Tests
// ──────────────────────────────────────────────

func TestBuildValidationSummaryAllMet(t *testing.T) {
	e := &Engine{
		archDesc: &ArchDescription{
			ExpectedResults: map[string]interface{}{
				"minimum_assumptions": float64(1),
				"minimum_critical":    float64(0),
				"minimum_high":        float64(0),
			},
		},
	}
	result := &AnalysisResult{
		TotalAssumptions: 5,
		CriticalCount:    1,
		HighCount:        2,
		MediumCount:      1,
		LowCount:         1,
		StrideDistribution: map[StrideCategory]int{
			StrideSpoofing: 2,
		},
	}
	summary := e.buildValidationSummary(result)
	if strings.Contains(summary, "violation") {
		t.Errorf("expected no violations, got: %s", summary)
	}
	if !strings.Contains(summary, "all expected criteria met") {
		t.Errorf("expected 'all expected criteria met', got: %s", summary)
	}
}

func TestBuildValidationSummaryViolations(t *testing.T) {
	e := &Engine{
		archDesc: &ArchDescription{
			ExpectedResults: map[string]interface{}{
				"minimum_assumptions": float64(100),
				"minimum_critical":    float64(5),
			},
		},
	}
	result := &AnalysisResult{
		TotalAssumptions: 10,
		CriticalCount:    1,
		HighCount:        2,
		MediumCount:      5,
		LowCount:         2,
	}
	summary := e.buildValidationSummary(result)
	if !strings.Contains(summary, "violation") {
		t.Errorf("expected violations, got: %s", summary)
	}
	if !strings.Contains(summary, "100") || !strings.Contains(summary, "5") {
		t.Errorf("expected threshold values in summary: %s", summary)
	}
}

func TestBuildValidationSummarySTRIDECategories(t *testing.T) {
	e := &Engine{
		archDesc: &ArchDescription{
			ExpectedResults: map[string]interface{}{
				"expected_stride_categories": []interface{}{
					"Spoofing", "Tampering",
				},
			},
		},
	}
	result := &AnalysisResult{
		TotalAssumptions: 3,
		StrideDistribution: map[StrideCategory]int{
			StrideSpoofing:  2,
			StrideTampering: 1,
		},
	}
	summary := e.buildValidationSummary(result)
	if strings.Contains(summary, "violation") {
		t.Errorf("expected no violations when all STRIDE categories present, got: %s", summary)
	}
}

func TestBuildValidationSummarySTRIDEMissing(t *testing.T) {
	e := &Engine{
		archDesc: &ArchDescription{
			ExpectedResults: map[string]interface{}{
				"expected_stride_categories": []interface{}{
					"Spoofing", "Denial of Service",
				},
			},
		},
	}
	result := &AnalysisResult{
		TotalAssumptions: 2,
		StrideDistribution: map[StrideCategory]int{
			StrideSpoofing: 1,
		},
	}
	summary := e.buildValidationSummary(result)
	if !strings.Contains(summary, "violation") {
		t.Errorf("expected violation for missing STRIDE category, got: %s", summary)
	}
	if !strings.Contains(summary, "Denial of Service") {
		t.Errorf("expected missing category in summary: %s", summary)
	}
}

func TestBuildValidationSummaryEmptyExpected(t *testing.T) {
	e := &Engine{archDesc: &ArchDescription{ExpectedResults: map[string]interface{}{}}}
	result := &AnalysisResult{TotalAssumptions: 1}
	summary := e.buildValidationSummary(result)
	if !strings.Contains(summary, "all expected criteria met") {
		t.Errorf("expected 'all expected criteria met' for empty expected results, got: %s", summary)
	}
}

// ──────────────────────────────────────────────
// YAML Parsing with New Fields
// ──────────────────────────────────────────────

func TestParseYAMLWithAssumptions(t *testing.T) {
	yaml := `name: test
description: Test arch
components:
  - name: WebApp
    type: web
    description: Web application
relationships:
  - source: WebApp
    target: Database
    protocol: TLS
    description: API calls
assumptions:
  - MFA is enforced for all user authentication
  - Data is encrypted at rest
security_controls:
  authentication:
    - MFA
  encryption:
    - AES256
metadata:
  compliance:
    - HIPAA
expected_results:
  minimum_assumptions: 1
validation_criteria:
  - All access must be authenticated
notes:
  - Test note`

	tmpFile, err := os.CreateTemp("", "asf-test-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	tmpName := tmpFile.Name()
	defer os.Remove(tmpName)

	if _, err := tmpFile.Write([]byte(yaml)); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	desc, err := parseYAMLArch(tmpName)
	if err != nil {
		t.Fatalf("parseYAMLArch failed: %v", err)
	}

	if desc.Name == "" {
		t.Error("expected non-empty name")
	}
	if len(desc.Components) == 0 {
		t.Error("expected at least 1 component")
	}
	if len(desc.ExplicitAssumptions) != 2 {
		t.Errorf("expected 2 explicit assumptions, got %d", len(desc.ExplicitAssumptions))
	}
	if len(desc.SecurityControls) != 2 {
		t.Errorf("expected 2 security control categories, got %d", len(desc.SecurityControls))
	}
	if len(desc.Compliance) != 1 || desc.Compliance[0] != "HIPAA" {
		t.Errorf("expected HIPAA compliance, got %v", desc.Compliance)
	}
	if desc.ExpectedResults == nil {
		t.Error("expected expected_results map")
	}
	if len(desc.ValidationCriteria) != 1 {
		t.Errorf("expected 1 validation criteria, got %d", len(desc.ValidationCriteria))
	}
	if len(desc.Notes) != 1 {
		t.Errorf("expected 1 note, got %d", len(desc.Notes))
	}
	if desc.RawText == "" {
		t.Error("expected non-empty RawText")
	}
}

func TestParseYAMLNoAssumptions(t *testing.T) {
	yaml := `name: minimal
description: No assumptions
components:
  - name: App
    type: app
    description: Application`

	tmpFile, err := os.CreateTemp("", "asf-test-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	tmpName := tmpFile.Name()
	defer os.Remove(tmpName)

	if _, err := tmpFile.Write([]byte(yaml)); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	desc, err := parseYAMLArch(tmpName)
	if err != nil {
		t.Fatalf("parseYAMLArch failed: %v", err)
	}

	if len(desc.ExplicitAssumptions) != 0 {
		t.Errorf("expected no explicit assumptions for minimal YAML, got %v", desc.ExplicitAssumptions)
	}
	if len(desc.SecurityControls) != 0 {
		t.Errorf("expected no security controls for minimal YAML, got %d", len(desc.SecurityControls))
	}
}

func TestParseYAMLEmptyAssumptions(t *testing.T) {
	yaml := `name: empty
description: Empty assumptions list
components:
  - name: App
    type: app
    description: Application
assumptions: []`

	tmpFile, err := os.CreateTemp("", "asf-test-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	tmpName := tmpFile.Name()
	defer os.Remove(tmpName)

	if _, err := tmpFile.Write([]byte(yaml)); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	desc, err := parseYAMLArch(tmpName)
	if err != nil {
		t.Fatalf("parseYAMLArch failed: %v", err)
	}

	if len(desc.ExplicitAssumptions) != 0 {
		t.Errorf("expected 0 explicit assumptions for empty list, got %d", len(desc.ExplicitAssumptions))
	}
}

func TestParseJSONWithAssumptions(t *testing.T) {
	json := `{
		"name": "test-json",
		"description": "Test JSON arch",
		"components": [{"name": "App", "type": "app", "description": "App"}],
		"relationships": [{"source": "App", "target": "DB", "protocol": "TLS", "description": "conn"}],
		"assumptions": ["MFA is enforced"],
		"security_controls": {"authentication": ["MFA"]},
		"metadata": {"compliance": ["SOC2"]},
		"expected_results": {"minimum_assumptions": 1},
		"notes": ["JSON test"]
	}`

	tmpFile, err := os.CreateTemp("", "asf-test-*.json")
	if err != nil {
		t.Fatal(err)
	}
	tmpName := tmpFile.Name()
	defer os.Remove(tmpName)

	if _, err := tmpFile.Write([]byte(json)); err != nil {
		t.Fatal(err)
	}
	tmpFile.Close()

	desc, err := parseJSONArch(tmpName)
	if err != nil {
		t.Fatalf("parseJSONArch failed: %v", err)
	}

	if len(desc.ExplicitAssumptions) != 1 {
		t.Errorf("expected 1 explicit assumption, got %d", len(desc.ExplicitAssumptions))
	}
	if len(desc.SecurityControls) != 1 {
		t.Errorf("expected 1 security control category, got %d", len(desc.SecurityControls))
	}
	if len(desc.Compliance) != 1 || desc.Compliance[0] != "SOC2" {
		t.Errorf("expected SOC2 compliance, got %v", desc.Compliance)
	}
	if len(desc.Notes) != 1 {
		t.Errorf("expected 1 note, got %d", len(desc.Notes))
	}
}

// ──────────────────────────────────────────────
// Explicit Assumptions Integration Test
// ──────────────────────────────────────────────

func TestProcessExplicitAssumptions(t *testing.T) {
	e := &Engine{
		strideEngine: NewStrideEngine(),
		archDesc: &ArchDescription{
			ExplicitAssumptions: []string{
				"MFA is enforced for all Auth0 user authentication",
				"PHI data is encrypted at rest in PHIDatabase",
			},
		},
	}
	e.explainPipe = NewExplainabilityPipeline(e.archDesc, "test.yaml", e.strideEngine)

	existing := []Assumption{
		{ID: "asm_001", Description: "Database access is restricted"},
	}

	explicit := e.processExplicitAssumptions(existing)
	if len(explicit) == 0 {
		t.Fatal("expected explicit assumptions")
	}

	for _, ea := range explicit {
		if ea.ID == "" {
			t.Error("expected non-empty ID")
		}
		if ea.Description == "" {
			t.Error("expected non-empty description")
		}
		if ea.Category == "" {
			t.Error("expected non-empty category")
		}
		if ea.Risk == "" {
			t.Error("expected non-empty risk level")
		}
		if ea.Confidence <= 0 {
			t.Errorf("expected positive confidence, got %.2f", ea.Confidence)
		}
		if ea.Keywords == nil {
			t.Error("expected non-nil keywords")
		}
		if len(ea.Stride) == 0 {
			t.Error("expected at least one STRIDE category")
		}
		if ea.Likelihood < 1 || ea.Likelihood > 5 {
			t.Errorf("likelihood %d out of range", ea.Likelihood)
		}
		if ea.Impact < 1 || ea.Impact > 5 {
			t.Errorf("impact %d out of range", ea.Impact)
		}
	}
}

func TestProcessExplicitAssumptionsDedup(t *testing.T) {
	e := &Engine{
		strideEngine: NewStrideEngine(),
		archDesc: &ArchDescription{
			ExplicitAssumptions: []string{
				"MFA is enforced for all Auth0 user authentication",
				"MFA is enforced for all Auth0 user authentication.",
			},
		},
	}

	existing := []Assumption{}

	explicit := e.processExplicitAssumptions(existing)
	if len(explicit) > 1 {
		t.Errorf("expected dedup to remove duplicate, got %d assumptions", len(explicit))
	}
}

func TestProcessExplicitAssumptionsDedupAgainstExisting(t *testing.T) {
	e := &Engine{
		strideEngine: NewStrideEngine(),
		archDesc: &ArchDescription{
			ExplicitAssumptions: []string{
				"MFA is enforced for all Auth0 user authentication",
			},
		},
	}

	existing := []Assumption{
		{ID: "asm_001", Description: "MFA is enforced for all Auth0 user authentication."},
	}

	explicit := e.processExplicitAssumptions(existing)
	if len(explicit) != 0 {
		t.Errorf("expected 0 new assumptions (all deduped), got %d", len(explicit))
	}
}

func TestProcessExplicitAssumptionsEmpty(t *testing.T) {
	e := &Engine{
		strideEngine: NewStrideEngine(),
		archDesc:     &ArchDescription{},
	}
	explicit := e.processExplicitAssumptions(nil)
	if len(explicit) != 0 {
		t.Errorf("expected 0 explicit assumptions for empty archDesc, got %d", len(explicit))
	}
}

// ──────────────────────────────────────────────
// Full Pipeline via buildResult
// ──────────────────────────────────────────────

func TestBuildResultWithExplicitAssumptions(t *testing.T) {
	e := &Engine{
		strideEngine: NewStrideEngine(),
		archDesc: &ArchDescription{
			Name: "test",
			ExplicitAssumptions: []string{
				"MFA is enforced for all Auth0 user authentication",
			},
			SecurityControls: map[string][]string{
				"authentication": {"MFA"},
			},
			Compliance: []string{"HIPAA"},
			ExpectedResults: map[string]interface{}{
				"minimum_assumptions": float64(1),
			},
		},
	}

	r := &asfJSONResult{}
	r.Summary.Assumptions = 0

	result := e.buildResult(r, "test.yaml", ModeASFOnly)

	if result.TotalAssumptions != 1 {
		t.Errorf("expected 1 total assumption, got %d", result.TotalAssumptions)
	}
	if len(result.Assumptions) != 1 {
		t.Errorf("expected 1 assumption in list, got %d", len(result.Assumptions))
	}
	if len(result.Compliance) == 0 {
		t.Error("expected compliance output")
	}
	if len(result.Controls) == 0 {
		t.Error("expected at least one control")
	}
	if !strings.Contains(result.Summary, "criteria met") {
		t.Errorf("expected validation in summary, got: %s", result.Summary)
	}
}

// ──────────────────────────────────────────────
// Edge Cases
// ──────────────────────────────────────────────

func TestProcessExplicitAssumptionsNoExplainPipe(t *testing.T) {
	e := &Engine{
		strideEngine: NewStrideEngine(),
		archDesc: &ArchDescription{
			ExplicitAssumptions: []string{"MFA is enforced"},
		},
		explainPipe: nil,
	}

	explicit := e.processExplicitAssumptions(nil)
	if len(explicit) != 1 {
		t.Fatalf("expected 1 assumption, got %d", len(explicit))
	}
}

func TestExtractKeywordsEmpty(t *testing.T) {
	got := extractKeywords("")
	if len(got) != 0 {
		t.Errorf("expected empty keywords, got %v", got)
	}
}

func TestNormalizeTextWhitespace(t *testing.T) {
	a := normalizeText("  PHI   data   is   encrypted  .  ")
	b := normalizeText("PHI data is encrypted.")
	if a != b {
		t.Errorf("whitespace normalization mismatch: %q vs %q", a, b)
	}
}

func TestClassifyExplicitAssumptionEdgeCases(t *testing.T) {
	tests := []struct {
		text string
		want models.AssumptionType
	}{
		{"", models.AssumptionTypeGOVERNANCE},
		{"   ", models.AssumptionTypeGOVERNANCE},
		{"completely unrelated text with no security meaning", models.AssumptionTypeGOVERNANCE},
	}
	for _, tt := range tests {
		got := classifyExplicitAssumption(tt.text)
		if got != tt.want {
			t.Errorf("classifyExplicitAssumption(%q) = %s, want %s", tt.text, got, tt.want)
		}
	}
}

// ──────────────────────────────────────────────
// Determinism Tests
// ──────────────────────────────────────────────

func TestExplicitAssumptionsDeterministic(t *testing.T) {
	e := &Engine{
		strideEngine: NewStrideEngine(),
		archDesc: &ArchDescription{
			ExplicitAssumptions: []string{
				"MFA is enforced for all Auth0 user authentication",
				"PHI data is encrypted at rest",
				"Audit logging is immutable",
			},
		},
	}
	e.explainPipe = NewExplainabilityPipeline(e.archDesc, "test.yaml", e.strideEngine)

	r1 := e.processExplicitAssumptions(nil)
	r2 := e.processExplicitAssumptions(nil)

	if len(r1) != len(r2) {
		t.Fatalf("non-deterministic count: %d vs %d", len(r1), len(r2))
	}
	for i := range r1 {
		if r1[i].ID != r2[i].ID {
			t.Errorf("non-deterministic ID at index %d: %s vs %s", i, r1[i].ID, r2[i].ID)
		}
		if r1[i].Category != r2[i].Category {
			t.Errorf("non-deterministic category at index %d: %s vs %s", i, r1[i].Category, r2[i].Category)
		}
		if r1[i].Risk != r2[i].Risk {
			t.Errorf("non-deterministic risk at index %d: %s vs %s", i, r1[i].Risk, r2[i].Risk)
		}
	}
}

// ──────────────────────────────────────────────
// asftest.yaml End-to-End Parsing Test (CI)
// ──────────────────────────────────────────────

func TestParseASFTestYAML(t *testing.T) {
	desc, err := parseYAMLArch("testdata/asftest.yaml")
	if err != nil {
		t.Fatalf("parseYAMLArch(testdata/asftest.yaml) failed: %v", err)
	}
	if desc.Name == "" {
		t.Error("expected non-empty name")
	}
	if len(desc.ExplicitAssumptions) != 30 {
		t.Errorf("expected 30 explicit assumptions, got %d", len(desc.ExplicitAssumptions))
	}
	if len(desc.SecurityControls) < 5 {
		t.Errorf("expected 5+ security control categories, got %d", len(desc.SecurityControls))
	}
	if len(desc.Compliance) < 2 {
		t.Errorf("expected 2+ compliance frameworks, got %d", len(desc.Compliance))
	}
	if desc.ExpectedResults == nil {
		t.Error("expected expected_results map")
	}
	if len(desc.ValidationCriteria) == 0 {
		t.Error("expected validation criteria")
	}
	if len(desc.Notes) == 0 {
		t.Error("expected notes")
	}
	if desc.RawText == "" {
		t.Error("expected non-empty RawText")
	}
	if !strings.Contains(desc.RawText, "MFA") {
		t.Error("expected RawText to contain explicit assumptions")
	}
}

// ──────────────────────────────────────────────
// Phase 8: Deep Deduplication with Source Merging
// ──────────────────────────────────────────────

func TestMergeSourceMetadata(t *testing.T) {
	a := Assumption{
		ID:          "ASM-001",
		Description: "MFA is enforced",
		SourceType:  "explicit",
		SourceFile:  "test.yaml",
	}
	merged := mergeSourceMetadata(a, "MFA is enforced.", "assumptions", 1, "test.yaml")
	if merged.SourceType != "explicit" {
		t.Errorf("expected source type preserved, got %s", merged.SourceType)
	}
	if merged.SourceFile != "test.yaml" {
		t.Errorf("expected source file test.yaml, got %s", merged.SourceFile)
	}

	// Cross-file merge
	a2 := Assumption{ID: "ASM-002", Description: "MFA is enforced"}
	merged2 := mergeSourceMetadata(a2, "MFA is enforced", "assumptions", 0, "other.yaml")
	if merged2.SourceType != "merged" {
		t.Errorf("expected source type merged, got %s", merged2.SourceType)
	}
	if !strings.Contains(merged2.SourceFile, "other.yaml") {
		t.Errorf("expected source file to contain other.yaml, got %s", merged2.SourceFile)
	}
}

func TestNormalizeTextBulletPrefix(t *testing.T) {
	a := normalizeText("- MFA is enforced")
	b := normalizeText("MFA is enforced")
	if a != b {
		t.Errorf("bullet prefix should be stripped for dedup: %q vs %q", a, b)
	}

	c := normalizeText("* MFA is enforced")
	if c != b {
		t.Errorf("asterisk bullet should be stripped for dedup: %q vs %q", c, b)
	}
}

// ──────────────────────────────────────────────
// Phase 9: Security Controls Verification Wiring
// ──────────────────────────────────────────────

func TestApplySecurityControlVerification(t *testing.T) {
	securityControls := map[string][]string{
		"authentication": {"MFA", "Password_Policy"},
		"encryption":     {"AES256", "TLS_1.3"},
	}

	a := Assumption{
		Description: "MFA is enforced for all user authentication",
		Category:    "IDENTITY",
		Confidence:  0.75,
	}
	result := applySecurityControlVerification(a, securityControls)
	if result.VerificationStatus != "VERIFIED" {
		t.Errorf("expected VERIFIED, got %s", result.VerificationStatus)
	}
	if result.Confidence < 0.80 {
		t.Errorf("expected confidence >= 0.80 after control match, got %.2f", result.Confidence)
	}
}

func TestApplySecurityControlVerificationNoMatch(t *testing.T) {
	securityControls := map[string][]string{
		"authentication": {"MFA"},
	}

	a := Assumption{
		Description: "Network segmentation isolates the database",
		Category:    "NETWORK",
		Confidence:  0.75,
	}
	result := applySecurityControlVerification(a, securityControls)
	if result.VerificationStatus != "" {
		t.Errorf("expected empty status when no control match, got %s", result.VerificationStatus)
	}
}

// ──────────────────────────────────────────────
// Phase 10: Dynamic Confidence for Explicit Assumptions
// ──────────────────────────────────────────────

func TestComputeExplicitConfidenceBase(t *testing.T) {
	conf := computeExplicitConfidence("simple text", models.AssumptionTypeGOVERNANCE, nil)
	if conf != 0.75 {
		t.Errorf("expected base confidence 0.75, got %.2f", conf)
	}
}

func TestComputeExplicitConfidenceBoosted(t *testing.T) {
	securityControls := map[string][]string{
		"authentication": {"MFA", "Password_Policy"},
	}
	conf := computeExplicitConfidence("MFA is enforced for all users", models.AssumptionTypeIDENTITY, securityControls)
	if conf < 0.80 {
		t.Errorf("expected boosted confidence >= 0.80, got %.2f", conf)
	}
	if conf > 0.95 {
		t.Errorf("confidence should not exceed 0.95, got %.2f", conf)
	}
}

func TestComputeExplicitConfidenceEncryption(t *testing.T) {
	securityControls := map[string][]string{
		"encryption": {"AES256", "TLS_1.3"},
	}
	conf := computeExplicitConfidence("PHI data is encrypted using AES256", models.AssumptionTypeCONFIGURATION, securityControls)
	if conf < 0.80 {
		t.Errorf("expected boosted confidence for encryption match, got %.2f", conf)
	}
}

// ──────────────────────────────────────────────
// Phase 11: Architecture-Specific Controls
// ──────────────────────────────────────────────

func TestGenerateArchitectureSpecificControls(t *testing.T) {
	assumptions := []Assumption{
		{ID: "ASM-001", Category: "IDENTITY", Risk: RiskHigh},
		{ID: "ASM-002", Category: "CONFIGURATION", Risk: RiskHigh},
	}
	components := []Component{
		{ID: "Auth0", Label: "Auth0"},
		{ID: "DB", Label: "PHIDatabase"},
	}
	controls := generateControls(assumptions, components)
	if len(controls) == 0 {
		t.Fatal("expected controls")
	}

	foundDB := false
	foundAuth := false
	for _, c := range controls {
		if strings.Contains(c.Description, "PHIDatabase") || strings.Contains(c.Category, "DATABASE") {
			foundDB = true
		}
		if strings.Contains(c.Description, "Auth0") || strings.Contains(c.Category, "IDENTITY") {
			foundAuth = true
		}
	}
	if !foundDB {
		t.Error("expected database-specific control for PHIDatabase")
	}
	if !foundAuth {
		t.Error("expected identity-specific control for Auth0")
	}
}

func TestGenerateArchitectureSpecificControlsEmpty(t *testing.T) {
	assumptions := []Assumption{
		{ID: "ASM-001", Category: "ACCESS", Risk: RiskMedium},
	}
	controls := generateControls(assumptions, nil)
	if len(controls) == 0 {
		t.Error("expected controls even with nil components")
	}
}

// ──────────────────────────────────────────────
// Phase 12: Expected Results Validation Summary
// ──────────────────────────────────────────────

func TestBuildResultValidationSummary(t *testing.T) {
	e := &Engine{
		strideEngine: NewStrideEngine(),
		archDesc: &ArchDescription{
			Name: "test",
			ExplicitAssumptions: []string{
				"MFA is enforced for all Auth0 user authentication",
			},
			ExpectedResults: map[string]interface{}{
				"minimum_assumptions": float64(1),
			},
		},
	}

	r := &asfJSONResult{}
	r.Summary.Assumptions = 0

	result := e.buildResult(r, "test.yaml", ModeASFOnly)
	if !strings.Contains(result.Summary, "criteria met") {
		t.Errorf("expected validation summary with 'criteria met', got: %s", result.Summary)
	}
	if result.TotalAssumptions < 1 {
		t.Errorf("expected at least 1 assumption, got %d", result.TotalAssumptions)
	}
}

func TestBuildResultValidationSummaryViolation(t *testing.T) {
	e := &Engine{
		strideEngine: NewStrideEngine(),
		archDesc: &ArchDescription{
			Name: "test",
			ExplicitAssumptions: []string{
				"MFA is enforced for all Auth0 user authentication",
			},
			ExpectedResults: map[string]interface{}{
				"minimum_assumptions": float64(100),
			},
		},
	}

	r := &asfJSONResult{}
	result := e.buildResult(r, "test.yaml", ModeASFOnly)
	if !strings.Contains(result.Summary, "violation") {
		t.Errorf("expected violation in summary, got: %s", result.Summary)
	}
}
