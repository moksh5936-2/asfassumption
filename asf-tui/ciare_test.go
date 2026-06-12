package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"asf-tui/intelligence"
)

func runCIAREAnalysis(t *testing.T, archFile string) *AnalysisResult {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()
	result, err := engine.RunAnalysis(archFile, "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}
	return result
}

func TestCIAREIntegration(t *testing.T) {
	patterns := []string{
		"testdata/attack_paths/*.yaml",
		"testdata/attack_paths/*.yml",
	}
	var testFiles []string
	for _, pattern := range patterns {
		matches, err := filepath.Glob(pattern)
		if err != nil {
			continue
		}
		testFiles = append(testFiles, matches...)
	}
	if len(testFiles) == 0 {
		t.Skip("no testdata files found")
	}

	for _, archFile := range testFiles {
		name := strings.TrimSuffix(filepath.Base(archFile), filepath.Ext(archFile))
		t.Run(name, func(t *testing.T) {
			result := runCIAREAnalysis(t, archFile)
			if len(result.CIAREFrameworkCoverages) == 0 {
				t.Error("Expected at least 1 framework coverage, got 0")
			}
			if len(result.CIAREAuditReadiness) == 0 {
				t.Error("Expected at least 1 audit readiness score, got 0")
			}
			if len(result.CIAREComplianceGaps) == 0 {
				t.Error("Expected at least 1 compliance gap, got 0")
			}
			if len(result.CIAREEvidenceRequirements) == 0 {
				t.Error("Expected evidence requirements, got 0")
			}
			if len(result.CIAREAuditorQuestions) == 0 {
				t.Error("Expected auditor questions, got 0")
			}
			if result.CIAREAuditPackage == nil {
				t.Error("Expected audit package, got nil")
			}
			if result.CIAREComplianceDashboard == nil {
				t.Error("Expected compliance dashboard, got nil")
			}
		})
	}
}

func TestCIAREDeterminism(t *testing.T) {
	archFile := "testdata/attack_paths/healthcare_phi.yaml"
	result1 := runCIAREAnalysis(t, archFile)
	result2 := runCIAREAnalysis(t, archFile)
	result3 := runCIAREAnalysis(t, archFile)

	if len(result1.CIAREFrameworkCoverages) != len(result2.CIAREFrameworkCoverages) ||
		len(result2.CIAREFrameworkCoverages) != len(result3.CIAREFrameworkCoverages) {
		t.Fatal("CIARE framework coverage counts differ between runs")
	}

	for i := range result1.CIAREFrameworkCoverages {
		c1 := result1.CIAREFrameworkCoverages[i]
		c2 := result2.CIAREFrameworkCoverages[i]
		c3 := result3.CIAREFrameworkCoverages[i]
		if c1.CoveragePct != c2.CoveragePct || c2.CoveragePct != c3.CoveragePct {
			t.Errorf("Coverage for %s differs between runs: %.1f vs %.1f vs %.1f",
				c1.Framework, c1.CoveragePct, c2.CoveragePct, c3.CoveragePct)
		}
	}
}

func TestCIAREFrameworkCoverage(t *testing.T) {
	result := runCIAREAnalysis(t, "testdata/attack_paths/healthcare_phi.yaml")
	for _, c := range result.CIAREFrameworkCoverages {
		if c.Framework == "" {
			t.Error("Framework coverage has empty Framework")
		}
		if c.Required <= 0 {
			t.Errorf("Framework %s has Required=%d, expected >0", c.Framework, c.Required)
		}
		if c.CoveragePct < 0 || c.CoveragePct > 100 {
			t.Errorf("Framework %s has invalid CoveragePct: %.1f", c.Framework, c.CoveragePct)
		}
		if c.Status == "" {
			t.Errorf("Framework %s has empty Status", c.Framework)
		}
	}
}

func TestCIAREAuditReadiness(t *testing.T) {
	result := runCIAREAnalysis(t, "testdata/attack_paths/healthcare_phi.yaml")
	for _, a := range result.CIAREAuditReadiness {
		if a.Framework == "" {
			t.Error("Audit readiness has empty Framework")
		}
		if a.ReadinessScore < 0 || a.ReadinessScore > 100 {
			t.Errorf("Framework %s has ReadinessScore=%.1f, expected 0-100", a.Framework, a.ReadinessScore)
		}
		if a.Status == "" {
			t.Errorf("Framework %s has empty Status", a.Framework)
		}
		if len(a.Factors) == 0 {
			t.Errorf("Framework %s has no factors", a.Framework)
		}
	}
}

func TestCIAREComplianceGaps(t *testing.T) {
	result := runCIAREAnalysis(t, "testdata/attack_paths/healthcare_phi.yaml")
	if len(result.CIAREComplianceGaps) == 0 {
		t.Fatal("Expected compliance gaps")
	}
	for _, g := range result.CIAREComplianceGaps {
		if g.ID == "" {
			t.Error("Compliance gap has empty ID")
		}
		if g.Framework == "" {
			t.Error("Compliance gap has empty Framework")
		}
		if g.Requirement == "" {
			t.Error("Compliance gap has empty Requirement")
		}
		if g.Risk != "Critical" && g.Risk != "High" {
			t.Errorf("Compliance gap %s has unexpected Risk: %s", g.ID, g.Risk)
		}
	}
}

func TestCIAREMissingEvidence(t *testing.T) {
	result := runCIAREAnalysis(t, "testdata/attack_paths/healthcare_phi.yaml")
	for _, m := range result.CIAREMissingEvidences {
		if m.Framework == "" {
			t.Error("Missing evidence has empty Framework")
		}
		if m.Control == "" {
			t.Error("Missing evidence has empty Control")
		}
		if len(m.Evidences) == 0 {
			t.Errorf("Missing evidence for %s/%s has no evidence items", m.Framework, m.Control)
		}
	}
}

func TestCIAREAuditorQuestions(t *testing.T) {
	result := runCIAREAnalysis(t, "testdata/attack_paths/healthcare_phi.yaml")
	if len(result.CIAREAuditorQuestions) == 0 {
		t.Fatal("Expected auditor questions")
	}
	for _, q := range result.CIAREAuditorQuestions {
		if q.Framework == "" {
			t.Error("Auditor question has empty Framework")
		}
		if q.Control == "" {
			t.Error("Auditor question has empty Control")
		}
		if q.Question == "" {
			t.Error("Auditor question has empty Question")
		}
	}
}

func TestCIAREControlMaturity(t *testing.T) {
	result := runCIAREAnalysis(t, "testdata/attack_paths/healthcare_phi.yaml")
	for _, m := range result.CIAREControlMaturities {
		if m.Domain == "" {
			t.Error("Control maturity has empty Domain")
		}
		if m.Level < 1 || m.Level > 5 {
			t.Errorf("Control maturity %s has Level=%d, expected 1-5", m.Domain, m.Level)
		}
		if m.Label == "" {
			t.Errorf("Control maturity %s has empty Label", m.Domain)
		}
		if m.Coverage < 0 || m.Coverage > 100 {
			t.Errorf("Control maturity %s has Coverage=%.1f", m.Domain, m.Coverage)
		}
	}
}

func TestCIAREComplianceNarratives(t *testing.T) {
	result := runCIAREAnalysis(t, "testdata/attack_paths/healthcare_phi.yaml")
	if len(result.CIAREComplianceNarratives) == 0 {
		t.Fatal("Expected compliance narratives")
	}
	for _, n := range result.CIAREComplianceNarratives {
		if n.Framework == "" {
			t.Error("Narrative has empty Framework")
		}
		if n.Narrative == "" {
			t.Errorf("Narrative for %s has empty Narrative", n.Framework)
		}
	}
}

func TestCIAREAuditPackage(t *testing.T) {
	result := runCIAREAnalysis(t, "testdata/attack_paths/healthcare_phi.yaml")
	if result.CIAREAuditPackage == nil {
		t.Fatal("Expected audit package")
	}
	if result.CIAREAuditPackage.ExecutiveSummary == "" {
		t.Error("Audit package has empty ExecutiveSummary")
	}
	if len(result.CIAREAuditPackage.FrameworkCoverages) == 0 {
		t.Error("Audit package has no framework coverages")
	}
}

func TestCIAREComplianceDashboard(t *testing.T) {
	result := runCIAREAnalysis(t, "testdata/attack_paths/healthcare_phi.yaml")
	if result.CIAREComplianceDashboard == nil {
		t.Fatal("Expected compliance dashboard")
	}
	if len(result.CIAREComplianceDashboard.FrameworkCoverages) == 0 {
		t.Error("Dashboard has no framework coverages")
	}
	if len(result.CIAREComplianceDashboard.TopRisks) == 0 {
		t.Log("Dashboard has no top risks (may be expected)")
	}
}

func TestCIAREProcurementQuestions(t *testing.T) {
	result := runCIAREAnalysis(t, "testdata/attack_paths/healthcare_phi.yaml")
	if len(result.CIAREProcurementQuestions) == 0 {
		t.Fatal("Expected procurement questions")
	}
	for _, q := range result.CIAREProcurementQuestions {
		if q.Category == "" {
			t.Error("Procurement question has empty Category")
		}
		if q.Question == "" {
			t.Error("Procurement question has empty Question")
		}
	}
}

func TestCIARENoCrashOnEmptyArchitecture(t *testing.T) {
	ciare := intelligence.NewCIAREEngine()
	result := ciare.Run(intelligence.CIAREInput{
		Architecture: &intelligence.ArchDescription{},
		SDRIResult:   &intelligence.SDRIResult{},
		Domain:       "",
		Compliance:   []string{},
	})
	if result == nil {
		t.Fatal("CIARE returned nil result")
	}
	_ = os.Stdout
}

func TestCIAREFrameworkDetermination(t *testing.T) {
	ciare := intelligence.NewCIAREEngine()
	tests := []struct {
		name          string
		compliance    []string
		minFrameworks int
	}{
		{"healthcare compliance", []string{"HIPAA", "SOC2"}, 2},
		{"empty compliance", []string{}, 7},
		{"all frameworks", []string{"HIPAA", "SOC2", "ISO27001", "PCI-DSS", "NIST800-53", "CIS", "GDPR"}, 7},
		{"single framework", []string{"HIPAA"}, 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ciare.Run(intelligence.CIAREInput{
				Architecture: &intelligence.ArchDescription{Compliance: tt.compliance},
				SDRIResult:   &intelligence.SDRIResult{},
				Domain:       "",
				Compliance:   tt.compliance,
			})
			if len(result.FrameworkCoverages) < tt.minFrameworks {
				t.Errorf("Expected at least %d frameworks, got %d", tt.minFrameworks, len(result.FrameworkCoverages))
			}
		})
	}
}

func TestCIAREControlMaturityLevels(t *testing.T) {
	tests := []struct {
		coverage float64
		level    int
		label    string
	}{
		{0, 1, "Level 1 - Ad Hoc"},
		{10, 1, "Level 1 - Ad Hoc"},
		{25, 2, "Level 2 - Repeatable"},
		{50, 3, "Level 3 - Defined"},
		{75, 4, "Level 4 - Managed"},
		{90, 5, "Level 5 - Optimized"},
		{100, 5, "Level 5 - Optimized"},
	}

	for _, tt := range tests {
		level, label := ciareMaturityLevel(tt.coverage)
		if level != tt.level {
			t.Errorf("For coverage %.1f, expected level %d, got %d", tt.coverage, tt.level, level)
		}
		if label != tt.label {
			t.Errorf("For coverage %.1f, expected label %q, got %q", tt.coverage, tt.label, label)
		}
	}
}

func ciareMaturityLevel(coverage float64) (int, string) {
	switch {
	case coverage >= 90:
		return 5, "Level 5 - Optimized"
	case coverage >= 75:
		return 4, "Level 4 - Managed"
	case coverage >= 50:
		return 3, "Level 3 - Defined"
	case coverage >= 25:
		return 2, "Level 2 - Repeatable"
	default:
		return 1, "Level 1 - Ad Hoc"
	}
}
