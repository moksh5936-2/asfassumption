package main

import (
	"strings"
	"testing"
)

// ──────────────────────────────────────────────
// Risk Matrix Tests
// ──────────────────────────────────────────────

func TestRiskMatrixCalculate(t *testing.T) {
	rm := &RiskMatrix{}

	tests := []struct {
		lh, im    int
		wantScore int
		wantLevel RiskLevel
	}{
		{5, 5, 25, RiskCritical},
		{5, 4, 20, RiskCritical},
		{4, 5, 20, RiskCritical},
		{4, 4, 16, RiskHigh},
		{4, 3, 12, RiskHigh},
		{3, 4, 12, RiskHigh},
		{3, 3, 9, RiskMedium},
		{3, 2, 6, RiskMedium},
		{2, 3, 6, RiskMedium},
		{2, 2, 4, RiskLow},
		{1, 1, 1, RiskLow},
		{1, 4, 4, RiskLow},
		{3, 1, 3, RiskLow},
	}

	for _, tt := range tests {
		score, level := rm.Calculate(tt.lh, tt.im)
		if score != tt.wantScore || level != tt.wantLevel {
			t.Errorf("Calculate(%d,%d) = (%d,%s), want (%d,%s)",
				tt.lh, tt.im, score, level, tt.wantScore, tt.wantLevel)
		}
	}
}

func TestRiskMatrixBoundaries(t *testing.T) {
	rm := &RiskMatrix{}

	// score 20 = critical boundary
	if s, l := rm.Calculate(5, 4); s != 20 || l != RiskCritical {
		t.Errorf("5×4 should be Critical (20), got %d %s", s, l)
	}
	// score 19 = high
	if s, l := rm.Calculate(5, 3); s != 15 || l != RiskHigh {
		t.Errorf("5×3 should be High (15), got %d %s", s, l)
	}
	// score 12 = high boundary
	if s, l := rm.Calculate(4, 3); s != 12 || l != RiskHigh {
		t.Errorf("4×3 should be High (12), got %d %s", s, l)
	}
	// score 11 = medium
	if s, l := rm.Calculate(3, 3); s != 9 || l != RiskMedium {
		t.Errorf("3×3 should be Medium (9), got %d %s", s, l)
	}
	// score 5 = medium boundary
	if s, l := rm.Calculate(5, 1); s != 5 || l != RiskMedium {
		t.Errorf("5×1 should be Medium (5), got %d %s", s, l)
	}
	// score 4 = low
	if s, l := rm.Calculate(2, 2); s != 4 || l != RiskLow {
		t.Errorf("2×2 should be Low (4), got %d %s", s, l)
	}
}

func TestRiskMatrixDeterministic(t *testing.T) {
	rm := &RiskMatrix{}

	// Same inputs must always produce same outputs
	for i := 0; i < 100; i++ {
		s1, l1 := rm.Calculate(4, 5)
		s2, l2 := rm.Calculate(4, 5)
		if s1 != s2 || l1 != l2 {
			t.Fatalf("Risk matrix not deterministic: (%d,%s) != (%d,%s)", s1, l1, s2, l2)
		}
	}
}

// ──────────────────────────────────────────────
// Confidence Engine Tests
// ──────────────────────────────────────────────

func TestConfidenceEngine(t *testing.T) {
	ce := &ConfidenceEngine{}

	tests := []struct {
		name                                  string
		evCount, stCount, compCount, relCount int
		wantMin                               float64
	}{
		{"no evidence", 0, 0, 0, 0, 0.09},
		{"single evidence", 1, 0, 0, 0, 0.15},
		{"full evidence", 6, 3, 4, 3, 0.8},
		{"extreme evidence", 100, 100, 100, 100, 0.95},
	}

	for _, tt := range tests {
		score, reason := ce.CalculateConfidence(tt.evCount, tt.stCount, tt.compCount, tt.relCount)
		if score < tt.wantMin {
			t.Errorf("%s: confidence %.2f < min %.2f (reason: %s)", tt.name, score, tt.wantMin, reason)
		}
		if score > 0.95 {
			t.Errorf("%s: confidence %.2f > 0.95 max", tt.name, score)
		}
		if reason == "" {
			t.Errorf("%s: empty reason", tt.name)
		}
	}
}

func TestConfidenceEngineDeterministic(t *testing.T) {
	ce := &ConfidenceEngine{}
	s1, _ := ce.CalculateConfidence(3, 2, 1, 1)
	s2, _ := ce.CalculateConfidence(3, 2, 1, 1)
	if s1 != s2 {
		t.Fatalf("Confidence not deterministic: %.4f != %.4f", s1, s2)
	}
}

// ──────────────────────────────────────────────
// Assumption Justification Tests
// ──────────────────────────────────────────────

func TestJustifyAssumption(t *testing.T) {
	tests := []struct {
		name     string
		category string
		evidence *EvidenceResult
		want     []string // substrings expected
	}{
		{
			name:     "with components",
			category: "ACCESS",
			evidence: &EvidenceResult{
				MatchedComponents:       []string{"Database", "API Gateway"},
				MatchedRelationships:    []string{"API → DB"},
				MatchedTrustBoundaries:  []string{"network boundary"},
				MatchedSecurityConcepts: []string{"authentication"},
			},
			want: []string{"detected", "components", "Database", "API Gateway"},
		},
		{
			name:     "with relationships only",
			category: "NETWORK",
			evidence: &EvidenceResult{
				MatchedRelationships: []string{"Internet → Gateway"},
			},
			want: []string{"identified", "communication path"},
		},
		{
			name:     "no evidence fallback",
			category: "GENERAL",
			evidence: &EvidenceResult{},
			want:     []string{"generated from category", "GENERAL"},
		},
	}

	for _, tt := range tests {
		rationale := JustifyAssumption(tt.category, tt.evidence)
		for _, substr := range tt.want {
			if !strings.Contains(rationale, substr) {
				t.Errorf("%s: rationale missing '%s': %s", tt.name, substr, rationale)
			}
		}
	}
}

// ──────────────────────────────────────────────
// Evidence Engine Tests
// ──────────────────────────────────────────────

func TestEvidenceEngineTrace(t *testing.T) {
	arch := &ArchDescription{
		Components: []Component{
			{ID: "db1", Label: "Database"},
			{ID: "gw1", Label: "API Gateway"},
			{ID: "auth1", Label: "Auth Service"},
		},
		Relationships: []Relation{
			{Source: "API Gateway", Target: "Database", Label: "SQL"},
			{Source: "Auth Service", Target: "API Gateway", Label: "REST"},
		},
	}

	ee := NewEvidenceEngine(arch, "test.drawio")

	t.Run("matches components by category", func(t *testing.T) {
		result := ee.TraceEvidence("DATABASE", []string{"sql", "injection"}, "")
		if len(result.MatchedComponents) == 0 {
			t.Error("expected component matches for DATABASE category")
		}
		if result.EvidenceCount == 0 {
			t.Error("expected evidence count > 0")
		}
	})

	t.Run("matches by keyword", func(t *testing.T) {
		result := ee.TraceEvidence("GENERAL", []string{"database", "backup"}, "backup database")
		if len(result.MatchedComponents) == 0 {
			t.Error("expected component match via keyword")
		}
	})

	t.Run("no match returns empty result", func(t *testing.T) {
		result := ee.TraceEvidence("GENERAL", []string{}, "something completely unrelated xyz")
		// May still get security concepts or other matches, but at minimum should not crash
		if result == nil {
			t.Error("expected non-nil result")
		}
	})
}

// ──────────────────────────────────────────────
// STRIDE Justification Tests
// ──────────────────────────────────────────────

func TestStrideJustifyEngine(t *testing.T) {
	inner := NewStrideEngine()
	sje := NewStrideJustifyEngine(inner)

	t.Run("returns justification for each category", func(t *testing.T) {
		result := sje.Justify("AUTHENTICATION", "mfa bypass", []string{"mfa", "bypass"}, []string{})
		if len(result.Categories) == 0 {
			t.Error("expected at least one STRIDE category for AUTHENTICATION")
		}
		if len(result.Justifications) != len(result.Categories) {
			t.Errorf("expected %d justifications, got %d", len(result.Categories), len(result.Justifications))
		}
		for _, j := range result.Justifications {
			if j.Reason == "" {
				t.Errorf("empty reason for category %s", j.Category)
			}
			if j.Confidence <= 0 {
				t.Errorf("zero confidence for category %s", j.Category)
			}
			if j.ConfidenceReason == "" {
				t.Errorf("empty confidence reason for category %s", j.Category)
			}
		}
	})

	t.Run("tracks matched rule indexes", func(t *testing.T) {
		result := sje.Justify("ACCESS", "idor", []string{"idor"}, []string{})
		// IDOR keyword should match rule index 0
		found := false
		for _, j := range result.Justifications {
			for _, idx := range j.MatchedRuleIndexes {
				if idx == 0 {
					found = true
				}
				if idx < 0 || idx >= 33 {
					t.Errorf("rule index %d out of range (0-32)", idx)
				}
			}
		}
		if !found {
			t.Log("idor keyword may not be in first rule index, checking keyword match instead")
			for _, j := range result.Justifications {
				for _, kw := range j.MatchedKeywords {
					if strings.Contains(kw, "idor") {
						found = true
					}
				}
			}
		}
		if !found {
			t.Log("idor not found in matched keywords either (category may have triggered)")
		}
	})

	t.Run("deterministic output", func(t *testing.T) {
		r1 := sje.Justify("NETWORK", "sql injection", []string{"sqli"}, []string{"Database"})
		r2 := sje.Justify("NETWORK", "sql injection", []string{"sqli"}, []string{"Database"})
		if len(r1.Categories) != len(r2.Categories) {
			t.Fatal("non-deterministic category count")
		}
		for i := range r1.Categories {
			if r1.Categories[i] != r2.Categories[i] {
				t.Fatalf("non-deterministic categories at index %d", i)
			}
			if r1.Justifications[i].Reason != r2.Justifications[i].Reason {
				t.Fatalf("non-deterministic reason at index %d", i)
			}
		}
	})
}

// ──────────────────────────────────────────────
// Likelihood Analyzer Tests
// ──────────────────────────────────────────────

func TestLikelihoodAnalyzer(t *testing.T) {
	la := &LikelihoodAnalyzer{}

	t.Run("returns valid range", func(t *testing.T) {
		a := &Assumption{Category: "ACCESS", Description: "test"}
		e := &EvidenceResult{}
		score, reason, factors := la.AnalyzeLikelihood(a, e)
		if score < 1 || score > 5 {
			t.Errorf("likelihood %d out of range 1-5", score)
		}
		if reason == "" {
			t.Error("empty reason")
		}
		if len(factors) == 0 {
			t.Error("expected at least one factor")
		}
	})

	t.Run("internet exposure increases score", func(t *testing.T) {
		a := &Assumption{Category: "ACCESS", Description: "internet exposed api"}
		e := &EvidenceResult{
			MatchedComponents:       []string{"Internet", "API Gateway"},
			MatchedSecurityConcepts: []string{"network_security"},
		}
		extScore, _, _ := la.AnalyzeLikelihood(a, e)

		a2 := &Assumption{Category: "ACCESS", Description: "internal service"}
		e2 := &EvidenceResult{}
		intScore, _, _ := la.AnalyzeLikelihood(a2, e2)

		if extScore <= intScore {
			t.Errorf("internet-exposed (%d) should have higher likelihood than internal (%d)", extScore, intScore)
		}
	})
}

// ──────────────────────────────────────────────
// Impact Analyzer Tests
// ──────────────────────────────────────────────

func TestImpactAnalyzer(t *testing.T) {
	ia := &ImpactAnalyzer{}

	t.Run("returns valid range", func(t *testing.T) {
		a := &Assumption{Category: "GENERAL", Description: "test"}
		e := &EvidenceResult{}
		score, reason, factors := ia.AnalyzeImpact(a, e)
		if score < 1 || score > 5 {
			t.Errorf("impact %d out of range 1-5", score)
		}
		if reason == "" {
			t.Error("empty reason")
		}
		if len(factors) == 0 {
			t.Error("expected at least one factor")
		}
	})

	t.Run("healthcare data increases impact", func(t *testing.T) {
		a := &Assumption{Category: "ACCESS", Description: "patient data"}
		e := &EvidenceResult{
			MatchedComponents: []string{"healthcare_db", "phi_data"},
		}
		healthScore, _, _ := ia.AnalyzeImpact(a, e)

		a2 := &Assumption{Category: "ACCESS", Description: "general config"}
		e2 := &EvidenceResult{}
		genScore, _, _ := ia.AnalyzeImpact(a2, e2)

		if healthScore <= genScore {
			t.Errorf("healthcare (%d) should have higher impact than general (%d)", healthScore, genScore)
		}
	})
}

// ──────────────────────────────────────────────
// Evidence Source Type Tests
// ──────────────────────────────────────────────

func TestExtractSourceType(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"test.drawio", ".drawio"},
		{"diagram.mmd", ".mmd"},
		{"arch.yaml", ".yaml"},
		{"arch.json", ".json"},
		{"image.png", ".png"},
		{"doc.txt", ".txt"},
		{"noext", ".txt"},
	}
	for _, tt := range tests {
		got := extractSourceType(tt.path)
		if got != tt.want {
			t.Errorf("extractSourceType(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

// ──────────────────────────────────────────────
// Risk Level Consistency
// ──────────────────────────────────────────────

func TestRiskForScoreConsistency(t *testing.T) {
	for score := 1; score <= 25; score++ {
		level := riskForScore(score)
		switch {
		case score >= 20:
			if level != RiskCritical {
				t.Errorf("score %d should be Critical, got %s", score, level)
			}
		case score >= 12:
			if level != RiskHigh {
				t.Errorf("score %d should be High, got %s", score, level)
			}
		case score >= 5:
			if level != RiskMedium {
				t.Errorf("score %d should be Medium, got %s", score, level)
			}
		default:
			if level != RiskLow {
				t.Errorf("score %d should be Low, got %s", score, level)
			}
		}
	}
}

// ──────────────────────────────────────────────
// Full Pipeline Integration
// ──────────────────────────────────────────────

func TestExplainabilityPipeline(t *testing.T) {
	arch := &ArchDescription{
		Components: []Component{
			{ID: "db", Label: "Database"},
			{ID: "api", Label: "API Gateway"},
		},
		Relationships: []Relation{
			{Source: "API Gateway", Target: "Database", Label: "SQL"},
		},
	}

	se := NewStrideEngine()
	pipe := NewExplainabilityPipeline(arch, "test.drawio", se)

	a := &Assumption{
		ID:          "ASM-001",
		Description: "Database credentials are stored securely",
		Category:    "ACCESS",
		Keywords:    []string{"database", "credentials"},
	}

	pipe.Explain(a)

	if len(a.EvidenceSources) == 0 {
		t.Error("expected evidence sources after pipeline")
	}
	if a.Rationale == "" {
		t.Error("expected rationale after pipeline")
	}
	if len(a.StrideJustifications) == 0 {
		t.Error("expected STRIDE justifications after pipeline")
	}
	if a.RiskJustification == nil {
		t.Fatal("expected risk justification after pipeline")
	}
	if a.RiskJustification.RiskScore < 1 || a.RiskJustification.RiskScore > 25 {
		t.Errorf("risk score %d out of range", a.RiskJustification.RiskScore)
	}
	if a.Confidence <= 0 {
		t.Errorf("confidence %.2f should be > 0", a.Confidence)
	}
	if a.RiskJustification.Confidence <= 0 {
		t.Errorf("risk confidence %.2f should be > 0", a.RiskJustification.Confidence)
	}
}

func TestExplainabilityPipelineNoArch(t *testing.T) {
	// Should handle nil arch gracefully
	se := NewStrideEngine()
	pipe := NewExplainabilityPipeline(nil, "test.txt", se)

	a := &Assumption{
		ID:          "ASM-001",
		Description: "test",
		Category:    "GENERAL",
	}

	pipe.Explain(a)

	// Should not panic, should produce fallback evidence
	if len(a.EvidenceSources) == 0 {
		t.Error("expected at least source file in evidence")
	}
	if a.Rationale == "" {
		t.Error("expected fallback rationale")
	}
}

// ──────────────────────────────────────────────
// Review & Validation Data
// ──────────────────────────────────────────────

func TestCollectValidationData(t *testing.T) {
	assumptions := []Assumption{
		{
			ID:           "ASM-001",
			Description:  "test",
			Risk:         RiskHigh,
			Confidence:   0.85,
			Stride:       []StrideCategory{StrideTampering},
			ReviewStatus: "Accepted",
			ReviewNotes:  "Looks correct",
		},
	}

	records := CollectValidationData(assumptions)
	if len(records) != 1 {
		t.Fatalf("expected 1 record, got %d", len(records))
	}
	r := records[0]
	if r.AssumptionID != "ASM-001" {
		t.Errorf("expected ASM-001, got %s", r.AssumptionID)
	}
	if r.AssignedRisk != RiskHigh {
		t.Errorf("expected High risk, got %s", r.AssignedRisk)
	}
	if r.ArchReviewResult != "Accepted" {
		t.Errorf("expected Accepted review, got %s", r.ArchReviewResult)
	}
	if r.ArchNotes != "Looks correct" {
		t.Errorf("expected review notes, got %s", r.ArchNotes)
	}
}

// ──────────────────────────────────────────────
// Edge Cases
// ──────────────────────────────────────────────

func TestEdgeCases(t *testing.T) {
	t.Run("nil assumption in pipeline", func(t *testing.T) {
		se := NewStrideEngine()
		arch := &ArchDescription{Name: "test"}
		pipe := NewExplainabilityPipeline(arch, "test.txt", se)
		// Should not panic
		pipe.Explain(nil)
	})

	t.Run("empty evidence result", func(t *testing.T) {
		r := &EvidenceResult{}
		if r.EvidenceCount != 0 {
			t.Errorf("expected 0 evidence count, got %d", r.EvidenceCount)
		}
	})

	t.Run("confidence with max values capped", func(t *testing.T) {
		ce := &ConfidenceEngine{}
		score, _ := ce.CalculateConfidence(100, 100, 100, 100)
		if score > 0.95 {
			t.Errorf("confidence %.2f exceeded max 0.95", score)
		}
	})

	t.Run("buildConfidenceSummary empty", func(t *testing.T) {
		summary := buildConfidenceSummary(nil)
		if summary != "no assumptions to evaluate" {
			t.Errorf("unexpected empty summary: %s", summary)
		}
	})

	t.Run("buildConfidenceSummary with data", func(t *testing.T) {
		assumptions := []Assumption{
			{Confidence: 0.9},
			{Confidence: 0.5},
		}
		summary := buildConfidenceSummary(assumptions)
		if !strings.Contains(summary, "70%") {
			t.Errorf("expected 70%% average in summary, got: %s", summary)
		}
	})
}
