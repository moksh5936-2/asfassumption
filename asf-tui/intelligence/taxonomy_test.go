package intelligence

import (
	"strings"
	"testing"
)

func TestTaxonomyEngineCount(t *testing.T) {
	te := NewTaxonomyEngine()
	if te.Count() < 40 {
		t.Errorf("expected at least 40 categories, got %d", te.Count())
	}
}

func TestTaxonomyEngineAllCategories(t *testing.T) {
	te := NewTaxonomyEngine()
	names := te.AllCategories()
	if len(names) != te.Count() {
		t.Errorf("AllCategories returned %d names but Count is %d", len(names), te.Count())
	}
	required := []string{
		"Identity", "Authentication", "Authorization", "PrivilegeManagement",
		"SecretsManagement", "KeyManagement", "CertificateManagement", "Cryptography",
		"DataProtection", "DataRetention", "Privacy", "Logging", "Monitoring",
		"Alerting", "Auditability", "Backups", "DisasterRecovery", "Availability",
		"Resilience", "ThirdPartyRisk", "VendorRisk", "SupplyChain", "InfrastructureSecurity",
		"NetworkSegmentation", "CloudSecurity", "ContainerSecurity", "KubernetesSecurity",
		"OperationalSecurity", "Compliance", "Governance", "ChangeManagement",
		"IncidentResponse", "DetectionEngineering", "TrustBoundaries", "HumanProcess",
		"InsiderThreat", "SessionSecurity", "APISecurity", "ObjectLevelAuthorization",
		"BusinessContinuity", "DataGovernance", "EncryptionAtRest",
	}
	seen := make(map[string]bool)
	for _, n := range names {
		seen[n] = true
	}
	for _, r := range required {
		if !seen[r] {
			t.Errorf("missing required category: %s", r)
		}
	}
}

func TestTaxonomyEngineGetCategory(t *testing.T) {
	te := NewTaxonomyEngine()
	cat := te.GetCategory("KeyManagement")
	if cat == nil {
		t.Fatal("expected KeyManagement category")
	}
	if cat.Name != "KeyManagement" {
		t.Errorf("expected name KeyManagement, got %s", cat.Name)
	}
	if len(cat.Keywords) == 0 {
		t.Error("expected keywords for KeyManagement")
	}
	if len(cat.Patterns) == 0 {
		t.Error("expected patterns for KeyManagement")
	}
	if len(cat.RiskMappings) == 0 {
		t.Error("expected risk mappings for KeyManagement")
	}
	if len(cat.VerificationRules) == 0 {
		t.Error("expected verification rules for KeyManagement")
	}
	if cat.ExplainabilityTemplate == "" {
		t.Error("expected explainability template for KeyManagement")
	}
}

func TestTaxonomyEngineMatchCategory(t *testing.T) {
	te := NewTaxonomyEngine()
	tests := []struct {
		text string
		want string
	}{
		{"We use Auth0 for identity and SSO", "Authentication"},
		{"Database encryption with key rotation via KMS", "KeyManagement"},
		{"TLS certificates are managed by our CA", "CertificateManagement"},
		{"Multi-tenant SaaS with tenant isolation", "ObjectLevelAuthorization"},
		{"HIPAA compliance for PHI storage", "Compliance"},
		{"Audit logging and tamper detection", "Auditability"},
		{"Container image scanning and pod security", "ContainerSecurity"},
		{"API rate limiting and gateway validation", "APISecurity"},
	}
	for _, tt := range tests {
		matches := te.MatchCategory(tt.text)
		found := false
		for _, m := range matches {
			if m == tt.want {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("text %q: expected match for %s, got %v", tt.text, tt.want, matches)
		}
	}
}

func TestTaxonomyEngineRiskMappings(t *testing.T) {
	te := NewTaxonomyEngine()
	cat := te.GetCategory("DataProtection")
	if cat == nil {
		t.Fatal("expected DataProtection category")
	}
	if cat.RiskMappings["default"] != RiskHigh {
		t.Errorf("expected default risk High, got %s", cat.RiskMappings["default"])
	}
	if cat.RiskMappings["phi"] != RiskCritical {
		t.Errorf("expected phi risk Critical, got %s", cat.RiskMappings["phi"])
	}
}

func TestTaxonomyEngineVerificationRules(t *testing.T) {
	te := NewTaxonomyEngine()
	cat := te.GetCategory("Authentication")
	if cat == nil {
		t.Fatal("expected Authentication category")
	}
	if len(cat.VerificationRules) < 2 {
		t.Errorf("expected at least 2 verification rules, got %d", len(cat.VerificationRules))
	}
	for _, rule := range cat.VerificationRules {
		if strings.TrimSpace(rule) == "" {
			t.Error("empty verification rule found")
		}
	}
}

func TestTaxonomyEngineExplainabilityTemplate(t *testing.T) {
	te := NewTaxonomyEngine()
	for name, cat := range te.Categories {
		if cat.ExplainabilityTemplate == "" {
			t.Errorf("category %s has empty explainability template", name)
		}
	}
}

func TestTaxonomyEnginePatternsCompile(t *testing.T) {
	te := NewTaxonomyEngine()
	for name, cat := range te.Categories {
		for i, pat := range cat.Patterns {
			if pat == nil {
				t.Errorf("category %s has nil pattern at index %d", name, i)
			}
		}
	}
}
