package intelligence

import (
	"strings"
	"testing"
)

func TestNewIntelligenceEngine(t *testing.T) {
	ie := NewIntelligenceEngine()
	if ie == nil {
		t.Fatal("expected non-nil intelligence engine")
	}
	if ie.taxonomy == nil {
		t.Error("expected taxonomy engine")
	}
	if ie.domainEngine == nil {
		t.Error("expected domain engine")
	}
	if ie.contradiction == nil {
		t.Error("expected contradiction engine")
	}
	if ie.trustBoundary == nil {
		t.Error("expected trust boundary engine")
	}
	if ie.quality == nil {
		t.Error("expected quality engine")
	}
	if ie.explainability == nil {
		t.Error("expected explainability engine")
	}
}

func TestIntelligenceEngineRun(t *testing.T) {
	ie := NewIntelligenceEngine()
	arch := &ArchDescription{
		Name: "test-arch",
		Components: []Component{
			{ID: "db1", Label: "PatientDatabase"},
			{ID: "idp1", Label: "Auth0"},
			{ID: "gw1", Label: "API Gateway"},
			{ID: "bkp1", Label: "BackupService"},
		},
		Relationships: []Relation{
			{Source: "API Gateway", Target: "PatientDatabase", Label: "SQL"},
			{Source: "Auth0", Target: "API Gateway", Label: "OAuth"},
		},
		RawText: "HIPAA-compliant multi-tenant system with PHI.",
	}
	result := ie.Run(arch)

	if result.TotalAssumptions == 0 {
		t.Fatal("expected assumptions from Run")
	}
	if result.Domain != "Healthcare" {
		t.Errorf("expected Healthcare domain, got %s", result.Domain)
	}
	if len(result.Contradictions) == 0 {
		t.Log("no contradictions detected (may be expected if no conflicting assumptions)")
	}
	if len(result.TrustBoundaries) == 0 {
		t.Error("expected trust boundaries")
	}
	if len(result.QualityScores) == 0 {
		t.Error("expected quality scores")
	}
	if result.Summary == "" {
		t.Error("expected summary")
	}
	if result.Explainability == "" {
		t.Error("expected explainability")
	}
	if len(result.Compliance) == 0 {
		t.Error("expected compliance mappings for healthcare")
	}
	if len(result.Controls) == 0 {
		t.Error("expected controls for healthcare")
	}
	if result.CriticalCount+result.HighCount+result.MediumCount+result.LowCount != result.TotalAssumptions {
		t.Error("risk counts should sum to total assumptions")
	}
}

func TestIntelligenceEngineRunWithExistingAssumptions(t *testing.T) {
	ie := NewIntelligenceEngine()
	arch := &ArchDescription{
		Name: "test-arch",
		Components: []Component{
			{ID: "db1", Label: "Database"},
		},
		RawText: "System with database.",
	}
	existing := []Assumption{
		{ID: "EXIST-001", Description: "Existing explicit assumption", Category: "General", Risk: RiskMedium},
	}
	result := ie.RunWithExistingAssumptions(arch, existing)
	foundExisting := false
	for _, a := range result.Assumptions {
		if a.ID == "EXIST-001" {
			foundExisting = true
		}
	}
	if !foundExisting {
		t.Error("expected existing assumption to be merged")
	}
	if result.TotalAssumptions < len(existing) {
		t.Errorf("expected total >= %d, got %d", len(existing), result.TotalAssumptions)
	}
}

func TestIntelligenceEngineGetAssumptionsByCategory(t *testing.T) {
	ie := NewIntelligenceEngine()
	arch := &ArchDescription{
		Name: "test-arch",
		Components: []Component{
			{ID: "db1", Label: "Database"},
		},
		RawText: "System with database.",
	}
	result := ie.Run(arch)
	byCat := ie.GetAssumptionsByCategory(result, "TrustBoundaries")
	if len(byCat) == 0 {
		t.Log("no TrustBoundaries assumptions found (may depend on topology)")
	}
}

func TestIntelligenceEngineGetAssumptionsByRisk(t *testing.T) {
	ie := NewIntelligenceEngine()
	arch := &ArchDescription{
		Name: "test-arch",
		Components: []Component{
			{ID: "db1", Label: "PatientDatabase"},
		},
		RawText: "HIPAA system with PHI.",
	}
	result := ie.Run(arch)
	byRisk := ie.GetAssumptionsByRisk(result, RiskCritical)
	if len(byRisk) == 0 {
		t.Log("no Critical assumptions found (may depend on topology)")
	}
}

func TestIntelligenceEngineGetContradictionsBySeverity(t *testing.T) {
	ie := NewIntelligenceEngine()
	arch := &ArchDescription{
		Name: "test-arch",
		Components: []Component{
			{ID: "db1", Label: "Database"},
		},
		RawText: "System with database.",
	}
	result := ie.Run(arch)
	bySev := ie.GetContradictionsBySeverity(result, RiskCritical)
	// No assertions on count since contradictions depend on generated assumptions
	_ = bySev
}

func TestIntelligenceEngineGetTopQualityAssumptions(t *testing.T) {
	ie := NewIntelligenceEngine()
	arch := &ArchDescription{
		Name: "test-arch",
		Components: []Component{
			{ID: "db1", Label: "PatientDatabase"},
		},
		RawText: "HIPAA system with PHI.",
	}
	result := ie.Run(arch)
	top := ie.GetTopQualityAssumptions(result, 5)
	if len(top) > 5 {
		t.Errorf("expected at most 5 top assumptions, got %d", len(top))
	}
	if len(top) > 0 && top[0].Rationale == "" {
		t.Error("expected rationale for top assumption")
	}
}

func TestIntelligenceEngineValidateResult(t *testing.T) {
	ie := NewIntelligenceEngine()
	result := &IntelligenceResult{
		TotalAssumptions: 10,
		CriticalCount:    2,
		Contradictions:   []Contradiction{{RuleName: "TEST"}},
	}
	issues := ie.ValidateResult(result)
	if len(issues) != 0 {
		t.Errorf("expected 0 issues, got %v", issues)
	}

	badResult := &IntelligenceResult{
		TotalAssumptions: 0,
	}
	issues = ie.ValidateResult(badResult)
	if len(issues) == 0 {
		t.Error("expected issues for empty result")
	}

	badResult2 := &IntelligenceResult{
		TotalAssumptions: 5,
		CriticalCount:    10,
		Contradictions:   make([]Contradiction, 20),
	}
	issues = ie.ValidateResult(badResult2)
	if len(issues) < 2 {
		t.Errorf("expected at least 2 issues, got %v", issues)
	}
}

func TestIntelligenceEngineRunNilArch(t *testing.T) {
	ie := NewIntelligenceEngine()
	result := ie.Run(nil)
	if result.TotalAssumptions != 0 {
		t.Errorf("expected 0 assumptions for nil arch, got %d", result.TotalAssumptions)
	}
	if result.Domain != "" {
		t.Errorf("expected empty domain for nil arch, got %s", result.Domain)
	}
}

func TestIntelligenceEngineRunWithSaaS(t *testing.T) {
	ie := NewIntelligenceEngine()
	arch := &ArchDescription{
		Name: "saas-arch",
		Components: []Component{
			{ID: "t1", Label: "TenantManager"},
			{ID: "api1", Label: "API"},
			{ID: "db1", Label: "Database"},
		},
		RawText: "Multi-tenant SaaS with tenant isolation.",
	}
	result := ie.Run(arch)
	if result.Domain != "SaaS" {
		t.Errorf("expected SaaS domain, got %s", result.Domain)
	}
	foundTenant := false
	for _, a := range result.Assumptions {
		if strings.Contains(a.Description, "tenant") || strings.Contains(a.Description, "Tenant") {
			foundTenant = true
		}
	}
	if !foundTenant {
		t.Error("expected tenant-related assumptions for SaaS domain")
	}
}

func TestIntelligenceEngineRunWithKubernetes(t *testing.T) {
	ie := NewIntelligenceEngine()
	arch := &ArchDescription{
		Name: "k8s-arch",
		Components: []Component{
			{ID: "pod1", Label: "Pod"},
			{ID: "svc1", Label: "Service"},
			{ID: "ns1", Label: "Namespace"},
		},
		RawText: "Kubernetes cluster with RBAC and network policies.",
	}
	result := ie.Run(arch)
	if result.Domain != "Kubernetes" {
		t.Errorf("expected Kubernetes domain, got %s", result.Domain)
	}
}

func TestIntelligenceEngineRunWithVPN(t *testing.T) {
	ie := NewIntelligenceEngine()
	arch := &ArchDescription{
		Name: "vpn-arch",
		Components: []Component{
			{ID: "vpn1", Label: "VPNGateway"},
			{ID: "rc1", Label: "RemoteClient"},
		},
		RawText: "IPSec VPN for remote access.",
	}
	result := ie.Run(arch)
	if result.Domain != "VPN" {
		t.Errorf("expected VPN domain, got %s", result.Domain)
	}
}

func TestIntelligenceEngineRunWithFintech(t *testing.T) {
	ie := NewIntelligenceEngine()
	arch := &ArchDescription{
		Name: "fintech-arch",
		Components: []Component{
			{ID: "pay1", Label: "PaymentGateway"},
			{ID: "fraud1", Label: "FraudDetection"},
		},
		RawText: "PCI DSS payment processing with fraud detection.",
	}
	result := ie.Run(arch)
	if result.Domain != "Fintech" {
		t.Errorf("expected Fintech domain, got %s", result.Domain)
	}
	if len(result.Compliance) == 0 {
		t.Error("expected compliance mappings for Fintech")
	}
}

func TestIntelligenceEngineRunWithCloudNative(t *testing.T) {
	ie := NewIntelligenceEngine()
	arch := &ArchDescription{
		Name: "cloud-arch",
		Components: []Component{
			{ID: "fn1", Label: "Lambda"},
			{ID: "db1", Label: "DynamoDB"},
		},
		RawText: "Cloud-native serverless architecture.",
	}
	result := ie.Run(arch)
	if result.Domain != "CloudNative" {
		t.Errorf("expected CloudNative domain, got %s", result.Domain)
	}
}

func TestIntelligenceEngineRunWithIdentityPlatform(t *testing.T) {
	ie := NewIntelligenceEngine()
	arch := &ArchDescription{
		Name: "idp-arch",
		Components: []Component{
			{ID: "sso1", Label: "SSO"},
			{ID: "fed1", Label: "Federation"},
		},
		RawText: "Identity platform with SSO and SAML federation.",
	}
	result := ie.Run(arch)
	if result.Domain != "IdentityPlatform" {
		t.Errorf("expected IdentityPlatform domain, got %s", result.Domain)
	}
}

func TestIntelligenceEngineRunWithDataPlatform(t *testing.T) {
	ie := NewIntelligenceEngine()
	arch := &ArchDescription{
		Name: "data-arch",
		Components: []Component{
			{ID: "etl1", Label: "ETLPipeline"},
			{ID: "dw1", Label: "DataWarehouse"},
		},
		RawText: "Data platform with lineage and governance.",
	}
	result := ie.Run(arch)
	if result.Domain != "DataPlatform" {
		t.Errorf("expected DataPlatform domain, got %s", result.Domain)
	}
}

func TestIntelligenceEngineRunWithEnterprise(t *testing.T) {
	ie := NewIntelligenceEngine()
	arch := &ArchDescription{
		Name: "ent-arch",
		Components: []Component{
			{ID: "ad1", Label: "ActiveDirectory"},
			{ID: "app1", Label: "CorporateApp"},
		},
		RawText: "Enterprise identity lifecycle and access reviews.",
	}
	result := ie.Run(arch)
	if result.Domain != "Enterprise" {
		t.Errorf("expected Enterprise domain, got %s", result.Domain)
	}
}

func TestIntelligenceEngineSummaryContainsExpected(t *testing.T) {
	ie := NewIntelligenceEngine()
	arch := &ArchDescription{
		Name: "test-arch",
		Components: []Component{
			{ID: "db1", Label: "Database"},
		},
		RawText: "System.",
	}
	result := ie.Run(arch)
	if !strings.Contains(result.Summary, "test-arch") {
		t.Error("expected architecture name in summary")
	}
	if !strings.Contains(result.Summary, "Intelligence analysis") {
		t.Error("expected 'Intelligence analysis' in summary")
	}
}
