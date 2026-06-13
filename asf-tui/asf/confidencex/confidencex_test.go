package confidencex

import (
	"strings"
	"testing"
)

func TestEmptyInput(t *testing.T) {
	engine := NewExplainabilityEngine("general", nil)
	output := engine.RunAll()
	if output == nil {
		t.Fatal("expected non-nil output")
	}
	if len(output.Breakdowns) != 0 {
		t.Errorf("expected 0 breakdowns, got %d", len(output.Breakdowns))
	}
}

func TestSingleAssumptionExplainability(t *testing.T) {
	input := ConfidenceInput{
		AssumptionID:         "A1",
		AssumptionText:       "MFA is enforced for all users",
		Component:            "auth0",
		Category:             "identity",
		Risk:                 "Critical",
		Confidence:           72.0,
		EvidenceSources:      []string{"security policy document", "explicit control mapping"},
		SourceComponents:     []string{"auth0", "webapp"},
		Keywords:             []string{"mfa", "authentication", "enforced", "users"},
		Rationale:            "Auth0 is configured as the identity provider with MFA policies enforced at the application level",
		VerificationStatus:   "Unverified",
		Domain:               "general",
		HasTrustChain:        true,
		HasCoverageGap:       false,
		DependencyCentrality: 0.8,
		FailureRadius:        12,
		SupportingFactTexts:  []string{"AWS KMS Used", "PHI Database Present", "Encryption Enabled"},
		SupportingFactIDs:    []string{"F12", "F18", "F25"},
	}
	engine := NewExplainabilityEngine("healthcare", []ConfidenceInput{input})
	output := engine.RunAll()

	if len(output.Breakdowns) != 1 {
		t.Fatalf("expected 1 breakdown, got %d", len(output.Breakdowns))
	}

	bd := output.Breakdowns[0]
	if bd.AssumptionID != "A1" {
		t.Errorf("expected A1, got %s", bd.AssumptionID)
	}
	if bd.FinalConfidence <= 0 {
		t.Errorf("expected positive confidence, got %.2f", bd.FinalConfidence)
	}
	if len(bd.PositiveFactors) == 0 {
		t.Error("expected positive factors")
	}
	if len(bd.NegativeFactors) == 0 {
		t.Error("expected negative factors (unverified)")
	}
	if len(bd.SupportingFacts) == 0 {
		t.Error("expected supporting facts")
	}
	if len(bd.EvidenceContributions) == 0 {
		t.Error("expected evidence contributions")
	}
	if bd.DomainContribution == nil {
		t.Error("expected domain contribution for healthcare")
	}
	if bd.TrustContribution == nil {
		t.Error("expected trust contribution")
	}
	if bd.WhyExists == "" {
		t.Error("expected why exists explanation")
	}
	if bd.WhyUncertain == "" {
		t.Error("expected why uncertain explanation")
	}
	if bd.WhatIncreasesConfidence == "" {
		t.Error("expected what increases confidence")
	}
	if bd.WhatDecreasesConfidence == "" {
		t.Error("expected what decreases confidence")
	}
	if bd.StabilityClass == "" {
		t.Error("expected stability classification")
	}
	if bd.StabilityReason == "" {
		t.Error("expected stability reason")
	}
}

func TestLowConfidenceAssumption(t *testing.T) {
	input := ConfidenceInput{
		AssumptionID:       "A2",
		AssumptionText:     "Backup encryption keys are rotated quarterly",
		Component:          "backup-service",
		Confidence:         15.0,
		EvidenceSources:    nil,
		SourceComponents:   nil,
		Keywords:           []string{"backup", "encryption"},
		VerificationStatus: "",
		Domain:             "general",
		HasTrustChain:      false,
		HasCoverageGap:     true,
		BlindSpotScore:     75,
	}
	engine := NewExplainabilityEngine("general", []ConfidenceInput{input})
	output := engine.RunAll()

	if len(output.Breakdowns) != 1 {
		t.Fatalf("expected 1 breakdown, got %d", len(output.Breakdowns))
	}

	bd := output.Breakdowns[0]
	if bd.FinalConfidence >= 50 {
		t.Errorf("expected low confidence for weak input, got %.2f", bd.FinalConfidence)
	}
	if bd.StabilityClass != StabilityHighlySpeculative && bd.StabilityClass != StabilityWeak {
		t.Errorf("expected highly speculative or weak stability, got %s", bd.StabilityClass)
	}
	foundNoEvidence := false
	for _, f := range bd.NegativeFactors {
		if strings.Contains(f.Name, "No Evidence") {
			foundNoEvidence = true
			break
		}
	}
	if !foundNoEvidence {
		t.Error("expected 'No Evidence Sources' negative factor")
	}
}

func TestDomainContributions(t *testing.T) {
	tests := []struct {
		name            string
		domain          string
		text            string
		expectInfluence float64
	}{
		{"healthcare phi", "healthcare", "PHI data requires encryption at rest", 12.0},
		{"healthcare generic", "healthcare", "System uses standard logging", 5.0},
		{"fintech pci", "fintech", "PCI DSS compliance for transactions", 12.0},
		{"fintech generic", "fintech", "Standard authentication used", 5.0},
		{"kubernetes rbac", "kubernetes", "RBAC policies for pod access", 10.0},
		{"general", "general", "Standard assumption", 2.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewExplainabilityEngine(tt.domain, []ConfidenceInput{
				{AssumptionID: "A1", AssumptionText: tt.text, Domain: tt.domain, Confidence: 50},
			})
			output := engine.RunAll()
			if len(output.Breakdowns) == 0 {
				t.Fatal("expected breakdown")
			}
			dc := output.Breakdowns[0].DomainContribution
			if dc == nil {
				t.Fatal("expected domain contribution")
			}
			if dc.Influence != tt.expectInfluence {
				t.Errorf("expected influence %.1f, got %.1f", tt.expectInfluence, dc.Influence)
			}
		})
	}
}

func TestStabilityClassifications(t *testing.T) {
	tests := []struct {
		name       string
		confidence float64
		posCount   int
		negCount   int
		expected   ConfidenceStability
	}{
		{"very stable", 90.0, 5, 0, StabilityVeryStable},
		{"stable", 75.0, 3, 1, StabilityStable},
		{"moderate", 60.0, 2, 2, StabilityModerate},
		{"weak", 40.0, 1, 3, StabilityWeak},
		{"highly speculative", 20.0, 0, 5, StabilityHighlySpeculative},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var posFactors []ConfidenceFactor
			var negFactors []ConfidenceFactor
			for i := 0; i < tt.posCount; i++ {
				posFactors = append(posFactors, ConfidenceFactor{Name: "P", Type: "positive", Impact: 3.0})
			}
			for i := 0; i < tt.negCount; i++ {
				negFactors = append(negFactors, ConfidenceFactor{Name: "N", Type: "negative", Impact: -3.0})
			}
			engine := &ExplainabilityEngine{}
			stability := engine.classifyStability(tt.confidence, posFactors, negFactors)
			if stability != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, stability)
			}
		})
	}
}

func TestCISOTrustView(t *testing.T) {
	inputs := make([]ConfidenceInput, 10)
	for i := 0; i < 10; i++ {
		risk := "Medium"
		if i < 3 {
			risk = "Critical"
		}
		inputs[i] = ConfidenceInput{
			AssumptionID:   string(rune('A' + i)),
			AssumptionText: "Test assumption " + string(rune('A'+i)),
			Confidence:     float64(100 - i*10),
			Risk:           risk,
		}
	}
	engine := NewExplainabilityEngine("general", inputs)
	output := engine.RunAll()

	if output.CISOTrustView == nil {
		t.Fatal("expected CISO trust view")
	}
	if len(output.CISOTrustView.MostTrustedFindings) == 0 {
		t.Error("expected most trusted findings")
	}
	if len(output.CISOTrustView.LeastTrustedFindings) == 0 {
		t.Error("expected least trusted findings")
	}
	if len(output.CISOTrustView.MostCriticalLowConfidence) == 0 {
		t.Log("no critical low-confidence findings (expected for this test)")
	}
}

func TestArchitectReviewView(t *testing.T) {
	inputs := make([]ConfidenceInput, 8)
	for i := 0; i < 8; i++ {
		inputs[i] = ConfidenceInput{
			AssumptionID:   string(rune('A' + i)),
			AssumptionText: "Arch test " + string(rune('A'+i)),
			Confidence:     float64(30 + i*10),
		}
	}
	engine := NewExplainabilityEngine("general", inputs)
	output := engine.RunAll()

	if output.ArchitectReviewView == nil {
		t.Fatal("expected architect review view")
	}
	if len(output.ArchitectReviewView.RequiringValidation)+
		len(output.ArchitectReviewView.WeakSupport)+
		len(output.ArchitectReviewView.StrongSupport) == 0 {
		t.Error("expected at least some items in architect view")
	}
}

func TestExportMarkdown(t *testing.T) {
	input := ConfidenceInput{
		AssumptionID:        "A1",
		AssumptionText:      "MFA is enforced for all users",
		Confidence:          72.0,
		EvidenceSources:     []string{"security policy"},
		SourceComponents:    []string{"auth0"},
		Keywords:            []string{"mfa", "authentication"},
		Domain:              "healthcare",
		SupportingFactTexts: []string{"AWS KMS Used"},
		SupportingFactIDs:   []string{"F12"},
	}
	engine := NewExplainabilityEngine("healthcare", []ConfidenceInput{input})
	output := engine.RunAll()
	md := ExportMarkdown(output)
	if md == "" {
		t.Fatal("expected non-empty markdown")
	}
	if !strings.Contains(md, "Confidence & Explainability Report") {
		t.Error("expected report title")
	}
	if !strings.Contains(md, "MFA is enforced") {
		t.Error("expected assumption text")
	}
	if !strings.Contains(md, "Why ASF Believes This") {
		t.Error("expected why section")
	}
	if !strings.Contains(md, "Why ASF Is Uncertain") {
		t.Error("expected uncertainty section")
	}
}

func TestExportHTML(t *testing.T) {
	input := ConfidenceInput{
		AssumptionID:   "A1",
		AssumptionText: "RBAC is properly configured",
		Confidence:     65.0,
		Domain:         "general",
	}
	engine := NewExplainabilityEngine("general", []ConfidenceInput{input})
	output := engine.RunAll()
	html := ExportHTML(output)
	if html == "" {
		t.Fatal("expected non-empty HTML")
	}
	if !strings.Contains(html, "Confidence & Explainability Report") {
		t.Error("expected report title in HTML")
	}
	if !strings.Contains(html, "<html>") {
		t.Error("expected HTML tag")
	}
	if !strings.Contains(html, "RBAC is properly configured") {
		t.Error("expected assumption text in HTML")
	}
}

func TestExplainabilityDetail(t *testing.T) {
	input := ConfidenceInput{
		AssumptionID:         "A1",
		AssumptionText:       "Key rotation exists for encryption keys",
		Component:            "kms",
		Confidence:           72.0,
		EvidenceSources:      []string{"KMS Detected", "Encryption Detected"},
		SourceComponents:     []string{"kms", "phi-database"},
		Keywords:             []string{"key", "rotation", "encryption", "kms"},
		Rationale:            "KMS service detected with encryption keys; healthcare domain requires key management",
		VerificationStatus:   "Unverified",
		Domain:               "healthcare",
		HasTrustChain:        true,
		DependencyCentrality: 0.85,
		FailureRadius:        15,
		SupportingFactTexts:  []string{"AWS KMS Used", "PHI Database Present", "Encryption Enabled"},
		SupportingFactIDs:    []string{"F12", "F18", "F25"},
		FactCategories:       []string{"security", "infrastructure", "security"},
	}
	engine := NewExplainabilityEngine("healthcare", []ConfidenceInput{input})
	output := engine.RunAll()

	if len(output.Breakdowns) != 1 {
		t.Fatalf("expected 1 breakdown, got %d", len(output.Breakdowns))
	}

	bd := output.Breakdowns[0]

	t.Logf("Confidence: %.1f%%", bd.FinalConfidence)
	t.Logf("Stability: %s — %s", bd.StabilityClass, bd.StabilityReason)
	t.Logf("Why Exists: %s", bd.WhyExists)
	t.Logf("Why Uncertain: %s", bd.WhyUncertain)
	t.Logf("Increases: %s", bd.WhatIncreasesConfidence)
	t.Logf("Decreases: %s", bd.WhatDecreasesConfidence)
	t.Logf("Domain: %s +%.1f%%", bd.DomainContribution.Domain, bd.DomainContribution.Influence)
	t.Logf("Trust Chain: %v (influence=%.1f)", bd.TrustContribution.HasTrustChain, bd.TrustContribution.ChainInfluence)

	for _, f := range bd.PositiveFactors {
		t.Logf("  + %s: +%.1f", f.Name, f.Impact)
	}
	for _, f := range bd.NegativeFactors {
		t.Logf("  - %s: %.1f", f.Name, f.Impact)
	}
	for _, fc := range bd.SupportingFacts {
		t.Logf("  Fact %s: %s (%.1f%%)", fc.FactID, fc.FactText, fc.Contribution)
	}
	for _, ec := range bd.EvidenceContributions {
		status := "missing"
		if ec.Present {
			status = "present"
		}
		t.Logf("  Evidence %s: %s (impact=%.1f)", ec.EvidenceID, status, ec.Impact)
	}
}

func BenchmarkExplainabilityEngine(b *testing.B) {
	inputs := make([]ConfidenceInput, 100)
	for i := 0; i < 100; i++ {
		inputs[i] = ConfidenceInput{
			AssumptionID:         string(rune('A' + i%26)),
			AssumptionText:       "Benchmark assumption " + string(rune('A'+i%26)),
			Confidence:           float64(30 + i%70),
			EvidenceSources:      []string{"evidence1", "evidence2"},
			SourceComponents:     []string{"comp1", "comp2"},
			Keywords:             []string{"key1", "key2", "key3"},
			HasTrustChain:        i%2 == 0,
			DependencyCentrality: float64(i) / 100.0,
			Domain:               "healthcare",
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine := NewExplainabilityEngine("healthcare", inputs)
		engine.RunAll()
	}
}

func BenchmarkExplainabilityLarge(b *testing.B) {
	inputs := make([]ConfidenceInput, 500)
	for i := 0; i < 500; i++ {
		inputs[i] = ConfidenceInput{
			AssumptionID:   string(rune('A' + i%26)),
			AssumptionText: "Large benchmark assumption",
			Confidence:     50.0,
			Domain:         "general",
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine := NewExplainabilityEngine("general", inputs)
		engine.RunAll()
	}
}
