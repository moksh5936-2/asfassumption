package fidelity

import (
	"asf-tui/asf/fact"
	"testing"
)

// Test healthcare architecture with explicit facts
func TestHiddenAssumptionEngine_Healthcare(t *testing.T) {
	// Create fact inventory
	inv := fact.NewInventory()
	inv.Add(fact.Fact{ID: "f1", Text: "MFA is enabled for all users", Source: "yaml", FactType: "control", Category: "security", IsNegative: false})
	inv.Add(fact.Fact{ID: "f2", Text: "Encryption is enabled for data at rest", Source: "yaml", FactType: "control", Category: "security", IsNegative: false})
	inv.Add(fact.Fact{ID: "f3", Text: "HIPAA compliance is required", Source: "yaml", FactType: "compliance", Category: "compliance", IsNegative: false})
	inv.Add(fact.Fact{ID: "f4", Text: "Backups are automated daily", Source: "yaml", FactType: "control", Category: "availability", IsNegative: false})
	inv.Add(fact.Fact{ID: "f5", Text: "WAF is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false})
	inv.Add(fact.Fact{ID: "f6", Text: "Audit logging is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false})
	inv.Add(fact.Fact{ID: "f7", Text: "VPN is used for admin access", Source: "yaml", FactType: "control", Category: "security", IsNegative: false})

	// Create components
	components := []Component{
		{ID: "comp1", Label: "Patient Database"},
		{ID: "comp2", Label: "API Gateway"},
		{ID: "comp3", Label: "Auth0 Service"},
		{ID: "comp4", Label: "Audit Log"},
		{ID: "comp5", Label: "Load Balancer"},
	}

	relationships := []Relationship{
		{Source: "API Gateway", Target: "Patient Database", Label: "queries"},
		{Source: "Auth0 Service", Target: "API Gateway", Label: "authenticates"},
		{Source: "Load Balancer", Target: "API Gateway", Label: "routes"},
	}

	engine := NewHiddenAssumptionEngine(inv, "healthcare")
	assumptions := engine.Generate(inv, components, relationships)

	// Verify: no assumptions should contradict facts
	protection := fact.NewProtectionLayer(inv)
	for _, a := range assumptions {
		result := protection.CheckAssumption(a.Description)
		if !result.Allowed {
			t.Errorf("Assumption contradicts fact: %s - %s", a.Description, result.Reason)
		}
	}

	// Verify: no generic assumptions
	for _, a := range assumptions {
		if a.QualityScore < 0.5 {
			t.Errorf("Low quality assumption: %s (score: %.2f)", a.Description, a.QualityScore)
		}
	}

	// Verify: traceability
	if len(assumptions) == 0 {
		t.Error("No assumptions generated")
	}

	// Verify: no restated facts
	for _, a := range assumptions {
		for _, f := range inv.Facts {
			if a.Description == f.Text {
				t.Errorf("Assumption restates fact: %s", a.Description)
			}
		}
	}

	// Verify: healthcare domain assumptions
	hasBreakGlass := false
	hasClinicalLogging := false
	for _, a := range assumptions {
		if a.SourceType == "domain-derived" && a.Category == "compliance" {
			if contains(a.Keywords, "break-glass") || contains(a.Keywords, "break glass") {
				hasBreakGlass = true
			}
			if contains(a.Keywords, "clinical") || contains(a.Keywords, "patient safety") {
				hasClinicalLogging = true
			}
		}
	}
	if !hasBreakGlass {
		t.Error("Missing healthcare break-glass assumption")
	}
	if !hasClinicalLogging {
		t.Error("Missing healthcare clinical logging assumption")
	}
}

// Test fintech architecture with explicit facts
func TestHiddenAssumptionEngine_Fintech(t *testing.T) {
	inv := fact.NewInventory()
	inv.Add(fact.Fact{ID: "f1", Text: "MFA is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false})
	inv.Add(fact.Fact{ID: "f2", Text: "Encryption is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false})
	inv.Add(fact.Fact{ID: "f3", Text: "PCI DSS compliance is required", Source: "yaml", FactType: "compliance", Category: "compliance", IsNegative: false})
	inv.Add(fact.Fact{ID: "f4", Text: "Backups are enabled", Source: "yaml", FactType: "control", Category: "availability", IsNegative: false})

	components := []Component{
		{ID: "comp1", Label: "Payment Processor"},
		{ID: "comp2", Label: "API Gateway"},
	}

	engine := NewHiddenAssumptionEngine(inv, "fintech")
	assumptions := engine.Generate(inv, components, nil)

	// Verify: no contradictions
	protection := fact.NewProtectionLayer(inv)
	for _, a := range assumptions {
		result := protection.CheckAssumption(a.Description)
		if !result.Allowed {
			t.Errorf("Assumption contradicts fact: %s - %s", a.Description, result.Reason)
		}
	}

	// Verify: fintech domain assumptions
	hasFraud := false
	hasAML := false
	for _, a := range assumptions {
		if a.SourceType == "domain-derived" {
			if contains(a.Keywords, "fraud") {
				hasFraud = true
			}
			if contains(a.Keywords, "aml") {
				hasAML = true
			}
		}
	}
	if !hasFraud {
		t.Error("Missing fintech fraud detection assumption")
	}
	if !hasAML {
		t.Error("Missing fintech AML/KYC assumption")
	}
}

// Test negative facts (MFA disabled)
func TestHiddenAssumptionEngine_NegativeFacts(t *testing.T) {
	inv := fact.NewInventory()
	inv.Add(fact.Fact{ID: "f1", Text: "MFA is disabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: true})
	inv.Add(fact.Fact{ID: "f2", Text: "Encryption is disabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: true})

	components := []Component{
		{ID: "comp1", Label: "App Server"},
	}

	engine := NewHiddenAssumptionEngine(inv, "saas")
	assumptions := engine.Generate(inv, components, nil)

	// Verify: no assumptions that say MFA is required
	for _, a := range assumptions {
		if contains(a.Keywords, "mfa") && contains(a.Keywords, "required") {
			t.Errorf("Assumption contradicts negative fact: %s", a.Description)
		}
		if contains(a.Keywords, "encryption") && contains(a.Keywords, "required") {
			t.Errorf("Assumption contradicts negative fact: %s", a.Description)
		}
	}

	// Verify: compensating assumptions exist
	hasCompensating := false
	for _, a := range assumptions {
		if contains(a.Keywords, "compensating") || contains(a.Keywords, "alternative") {
			hasCompensating = true
		}
	}
	if !hasCompensating {
		t.Error("Missing compensating control assumption for disabled MFA")
	}
}

// Test contradiction detection
func TestRealContradictionEngine(t *testing.T) {
	inv := fact.NewInventory()
	inv.Add(fact.Fact{ID: "f1", Text: "MFA is required", Source: "yaml", FactType: "control", Category: "security", IsNegative: false})
	inv.Add(fact.Fact{ID: "f2", Text: "MFA is disabled for service accounts", Source: "yaml", FactType: "control", Category: "security", IsNegative: true})
	inv.Add(fact.Fact{ID: "f3", Text: "Encryption is required", Source: "yaml", FactType: "control", Category: "security", IsNegative: false})
	inv.Add(fact.Fact{ID: "f4", Text: "Encryption is disabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: true})

	engine := NewRealContradictionEngine(inv)
	contradictions := engine.Detect()

	// Should find 2 contradictions
	if len(contradictions) != 2 {
		t.Errorf("Expected 2 contradictions, got %d", len(contradictions))
	}

	// Verify: MFA contradiction
	mfaContradiction := false
	for _, c := range contradictions {
		if c.Description == "MFA is required but disabled" {
			mfaContradiction = true
		}
	}
	if !mfaContradiction {
		t.Error("Missing MFA contradiction")
	}

	// Verify: no self-comparison
	for _, c := range contradictions {
		if c.FactA.ID == c.FactB.ID {
			t.Error("Self-comparison detected")
		}
	}
}

// Test fact protection
func TestFactProtection(t *testing.T) {
	inv := fact.NewInventory()
	inv.Add(fact.Fact{ID: "f1", Text: "MFA is disabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: true})
	inv.Add(fact.Fact{ID: "f2", Text: "Encryption is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false})

	protection := fact.NewProtectionLayer(inv)

	// Should reject: MFA is required
	result := protection.CheckAssumption("MFA is required for all users")
	if result.Allowed {
		t.Error("Should reject assumption that contradicts negative fact")
	}

	// Should reject: restates fact
	result = protection.CheckAssumption("Encryption is enabled")
	if result.Allowed {
		t.Error("Should reject assumption that restates fact")
	}

	// Should allow: hidden assumption
	result = protection.CheckAssumption("Certificates are rotated before expiry")
	if !result.Allowed {
		t.Error("Should allow hidden assumption")
	}
}

// Test fidelity scoring
func TestFidelityScorer(t *testing.T) {
	inv := fact.NewInventory()
	inv.Add(fact.Fact{ID: "f1", Text: "MFA is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false})
	inv.Add(fact.Fact{ID: "f2", Text: "Encryption is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false})
	inv.Add(fact.Fact{ID: "f3", Text: "Backups are enabled", Source: "yaml", FactType: "control", Category: "availability", IsNegative: false})

	scorer := NewFidelityScorer(inv)

	assumptions := []HiddenAssumption{
		{ID: "a1", Description: "Certificates are rotated", QualityScore: 0.8, NoveltyScore: 0.9, RelevanceScore: 0.8, SourceType: "fact-derived"},
		{ID: "a2", Description: "Restore testing is performed", QualityScore: 0.8, NoveltyScore: 0.9, RelevanceScore: 0.8, SourceType: "fact-derived"},
	}

	contradictions := []RealContradiction{}
	traceability := []TraceabilityRecord{
		{AssumptionID: "a1", SourceFactID: "f2", SourceType: "fact-derived"},
		{AssumptionID: "a2", SourceFactID: "f3", SourceType: "fact-derived"},
	}

	score := scorer.Compute(assumptions, contradictions, traceability)

	if score.TotalFacts != 3 {
		t.Errorf("Expected 3 facts, got %d", score.TotalFacts)
	}
	if score.RespectedFacts != 2 {
		t.Errorf("Expected 2 respected facts, got %d", score.RespectedFacts)
	}
	if score.Score < 0.6 {
		t.Errorf("Score too low: %.2f", score.Score)
	}
	if score.Overall != "ARCHITECTURAL_FIDELITY_CERTIFIED" && score.Overall != "CONDITIONAL" && score.Overall != "NOT_CERTIFIED" {
		t.Errorf("Unexpected overall: %s", score.Overall)
	}
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
