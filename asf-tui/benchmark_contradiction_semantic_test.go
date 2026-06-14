package main

import (
	"testing"

	"asf-tui/intelligence"
)

func TestSemanticContradictionHealthcare(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()
	result, err := engine.RunAnalysis("testdata/benchmarks/healthcare_contradictions.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	t.Logf("=== Healthcare Architecture Contradictions ===")
	t.Logf("Total CIEContradictions (after dedup): %d", len(result.CIEContradictions))
	for _, c := range result.CIEContradictions {
		t.Logf("  [%s] %s (Conf=%.2f): %s", c.Type, c.Severity, c.Confidence, c.Summary)
	}

	types := make(map[string]bool)
	for _, c := range result.CIEContradictions {
		types[c.Type] = true
	}

	requiredTypes := []string{"ENCRYPTION", "AUTHENTICATION", "AUTHORIZATION", "NETWORK", "COMPLIANCE"}

	for _, rt := range requiredTypes {
		if !types[rt] {
			t.Errorf("Missing contradiction type: %s", rt)
		} else {
			t.Logf("  ✓ Detected: %s", rt)
		}
	}

	if len(result.CIEContradictions) < 5 {
		t.Errorf("Healthcare: expected >=5 contradictions, got %d", len(result.CIEContradictions))
	}
}

func TestSemanticContradictionPayroll(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()
	result, err := engine.RunAnalysis("testdata/benchmarks/payroll_contradictions.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	t.Logf("=== Payroll Architecture Contradictions ===")
	t.Logf("Total CIEContradictions (after dedup): %d", len(result.CIEContradictions))
	for _, c := range result.CIEContradictions {
		t.Logf("  [%s] %s (Conf=%.2f): %s", c.Type, c.Severity, c.Confidence, c.Summary)
	}

	types := make(map[string]bool)
	for _, c := range result.CIEContradictions {
		types[c.Type] = true
	}

	requiredTypes := []string{"ENCRYPTION", "AUTHENTICATION", "AUTHORIZATION"}

	for _, rt := range requiredTypes {
		if !types[rt] {
			t.Errorf("Missing contradiction type: %s", rt)
		} else {
			t.Logf("  ✓ Detected: %s", rt)
		}
	}

	if len(result.CIEContradictions) < 3 {
		t.Errorf("Payroll: expected >=3 contradictions, got %d", len(result.CIEContradictions))
	}
}

func TestSemanticContradictionCloud(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()
	result, err := engine.RunAnalysis("testdata/benchmarks/cloud_contradictions.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	t.Logf("=== Cloud Architecture Contradictions ===")
	t.Logf("Total CIEContradictions (after dedup): %d", len(result.CIEContradictions))
	for _, c := range result.CIEContradictions {
		t.Logf("  [%s] %s (Conf=%.2f): %s", c.Type, c.Severity, c.Confidence, c.Summary)
	}

	types := make(map[string]bool)
	for _, c := range result.CIEContradictions {
		types[c.Type] = true
	}

	requiredTypes := []string{"ENCRYPTION", "NETWORK", "AUTHORIZATION", "LOGGING"}

	for _, rt := range requiredTypes {
		if !types[rt] {
			t.Errorf("Missing contradiction type: %s", rt)
		} else {
			t.Logf("  ✓ Detected: %s", rt)
		}
	}

	if len(result.CIEContradictions) < 4 {
		t.Errorf("Cloud: expected >=4 contradictions, got %d", len(result.CIEContradictions))
	}
}

func TestSemanticContradictionCleanArchitecture(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()
	result, err := engine.RunAnalysis("testdata/asftest.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	t.Logf("=== Clean Architecture (False Positive Check) ===")
	t.Logf("Total CIEContradictions: %d", len(result.CIEContradictions))
	for _, c := range result.CIEContradictions {
		t.Logf("  [%s] %s: %s", c.Type, c.Severity, c.Summary)
	}

	if len(result.CIEContradictions) > 0 {
		t.Errorf("Clean architecture should have 0 contradictions, got %d", len(result.CIEContradictions))
	}
}

func TestSemanticContradictionPrecisionPreserved(t *testing.T) {
	result := runBenchmarkFixture(t, "fixture_c_true_contradictions.yaml")

	t.Logf("=== Fixture C Contradiction Precision (Regression) ===")
	t.Logf("CIEContradictions: %d", len(result.CIEContradictions))
	for _, c := range result.CIEContradictions {
		t.Logf("  [%s] %s (Conf=%.2f): %s", c.Type, c.Severity, c.Confidence, c.Summary)
	}

	for _, c := range result.CIEContradictions {
		if c.StatementA.ID == c.StatementB.ID {
			t.Errorf("self-comparison contradiction detected: %s (A=%s == B=%s)", c.ID, c.StatementA.ID, c.StatementB.ID)
		}
	}

	seen := make(map[string]bool)
	for _, c := range result.CIEContradictions {
		key := c.StatementA.ID + "|" + c.StatementB.ID
		revKey := c.StatementB.ID + "|" + c.StatementA.ID
		if seen[key] || seen[revKey] {
			t.Errorf("duplicate contradiction pair: %s (%s <-> %s)", c.ID, c.StatementA.ID, c.StatementB.ID)
		}
		seen[key] = true
	}

	if len(result.CIEContradictions) < 4 {
		t.Errorf("too few contradictions: expected >=4, got %d", len(result.CIEContradictions))
	}
	if len(result.CIEContradictions) > 16 {
		t.Errorf("too many contradictions: expected <=16, got %d", len(result.CIEContradictions))
	}
}

func TestSemanticEngineUnit(t *testing.T) {
	se := intelligence.NewSemanticEngine()

	tests := []struct {
		name      string
		claims    []intelligence.Statement
		minWant   int
		wantTypes []intelligence.ContradictionType
	}{
		{
			name: "encryption plaintext",
			claims: []intelligence.Statement{
				{ID: "A1", OriginalText: "All data is encrypted at rest"},
				{ID: "A2", OriginalText: "Backups are stored in plaintext"},
			},
			minWant:   1,
			wantTypes: []intelligence.ContradictionType{intelligence.ContradictionTypeENCRYPTION},
		},
		{
			name: "mfa vs single factor",
			claims: []intelligence.Statement{
				{ID: "B1", OriginalText: "MFA is enforced for all authentication"},
				{ID: "B2", OriginalText: "Some accounts use password only"},
			},
			minWant:   1,
			wantTypes: []intelligence.ContradictionType{intelligence.ContradictionTypeAUTHENTICATION},
		},
		{
			name: "restricted vs open access",
			claims: []intelligence.Statement{
				{ID: "C1", OriginalText: "Access is restricted to authorized staff only"},
				{ID: "C2", OriginalText: "System allows unrestricted access to all"},
			},
			minWant:   1,
			wantTypes: []intelligence.ContradictionType{intelligence.ContradictionTypeAUTHORIZATION},
		},
		{
			name: "private vs public network",
			claims: []intelligence.Statement{
				{ID: "D1", OriginalText: "System is in a private network"},
				{ID: "D2", OriginalText: "Service is exposed to the internet"},
			},
			minWant:   1,
			wantTypes: []intelligence.ContradictionType{intelligence.ContradictionTypeNETWORK},
		},
		{
			name: "HA vs single instance",
			claims: []intelligence.Statement{
				{ID: "E1", OriginalText: "High availability is required for all services"},
				{ID: "E2", OriginalText: "Database runs as a single instance"},
			},
			minWant:   1,
			wantTypes: []intelligence.ContradictionType{intelligence.ContradictionTypeAVAILABILITY},
		},
		{
			name: "zero trust vs implicit trust",
			claims: []intelligence.Statement{
				{ID: "F1", OriginalText: "Zero trust architecture is enforced"},
				{ID: "F2", OriginalText: "Internal traffic is trusted by default"},
			},
			minWant:   1,
			wantTypes: []intelligence.ContradictionType{intelligence.ContradictionTypeTRUST},
		},
		{
			name: "no contradiction",
			claims: []intelligence.Statement{
				{ID: "G1", OriginalText: "All data is encrypted at rest"},
				{ID: "G2", OriginalText: "TLS 1.3 is enforced for all traffic"},
			},
			minWant:   0,
			wantTypes: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := se.DetectSemanticContradictions(tt.claims)
			t.Logf("  Found %d contradictions", len(got))
			for _, c := range got {
				t.Logf("    [%s] %s", c.Type, c.Summary)
			}
			if len(got) < tt.minWant {
				t.Errorf("expected >=%d contradictions, got %d", tt.minWant, len(got))
			}
			for _, wt := range tt.wantTypes {
				found := false
				for _, c := range got {
					if c.Type == wt {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected contradiction type %s not found", wt)
				}
			}
		})
	}
}
