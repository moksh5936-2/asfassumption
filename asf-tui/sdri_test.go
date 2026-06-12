package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSDRIIntegration(t *testing.T) {
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

			if len(result.SDRIControls) == 0 {
				t.Errorf("Expected at least 1 SDRI control, got 0")
			}

			if result.SDRISummary == "" {
				t.Errorf("Expected non-empty SDRISummary")
			}

			if len(result.SDRICoverageDashboard) == 0 {
				t.Errorf("Expected at least 1 coverage dashboard entry")
			}

			if len(result.SDRIDesignFindings) > 0 {
				for _, f := range result.SDRIDesignFindings {
					if f.ID == "" {
						t.Errorf("Design finding has empty ID")
					}
					if f.Title == "" {
						t.Errorf("Design finding %s has empty Title", f.ID)
					}
					if f.Severity == "" {
						t.Errorf("Design finding %s has empty Severity", f.ID)
					}
				}
			}

			if len(result.SDRIRemediations) > 0 {
				for _, r := range result.SDRIRemediations {
					if r.ID == "" {
						t.Errorf("Remediation has empty ID")
					}
					if r.Priority <= 0 {
						t.Errorf("Remediation %s has invalid Priority: %d", r.ID, r.Priority)
					}
					if r.RiskScore < 0 {
						t.Errorf("Remediation %s has negative RiskScore: %f", r.ID, r.RiskScore)
					}
				}
			}

			if len(result.SDRIComplianceAlignments) > 0 {
				for _, m := range result.SDRIComplianceAlignments {
					if m.Framework == "" {
						t.Errorf("Compliance mapping has empty Framework")
					}
					if m.Coverage < 0 || m.Coverage > 100 {
						t.Errorf("Compliance mapping %s has invalid Coverage: %f", m.Framework, m.Coverage)
					}
				}
			}
		})
	}
}

func TestSDRIDeterminism(t *testing.T) {
	archFile := "testdata/attack_paths/healthcare_phi.yaml"
	if _, err := os.Stat(archFile); os.IsNotExist(err) {
		t.Skip("testdata file not found")
	}

	cfg := &Config{}
	engine := NewEngine(cfg)

	var firstSummary string
	for i := 0; i < 3; i++ {
		progress := make(chan AnalysisProgress, 100)
		go func() {
			for range progress {
			}
		}()
		result, err := engine.RunAnalysis(archFile, "", ModeASFOnly, progress)
		if err != nil {
			t.Fatalf("Run %d failed: %v", i, err)
		}
		if i == 0 {
			firstSummary = result.SDRISummary
		} else if result.SDRISummary != firstSummary {
			t.Errorf("Run %d: SDRISummary changed (determinism failure)\nGot:  %s\nWant: %s",
				i, result.SDRISummary, firstSummary)
		}
	}
}

func TestSDRIDesignFindings(t *testing.T) {
	archFile := "testdata/attack_paths/healthcare_phi.yaml"
	if _, err := os.Stat(archFile); os.IsNotExist(err) {
		t.Skip("testdata file not found")
	}

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

	if len(result.SDRIDesignFindings) == 0 {
		t.Error("Expected at least 1 design finding")
	}

	findingMap := make(map[string]int)
	for _, f := range result.SDRIDesignFindings {
		findingMap[f.Severity]++
	}

	if findingMap["Critical"] > 0 || findingMap["High"] > 0 {
		t.Logf("Design findings by severity: Critical=%d, High=%d, Medium=%d, Low=%d",
			findingMap["Critical"], findingMap["High"], findingMap["Medium"], findingMap["Low"])
	}
}

func TestSDRICoverage(t *testing.T) {
	archFile := "testdata/attack_paths/healthcare_phi.yaml"
	if _, err := os.Stat(archFile); os.IsNotExist(err) {
		t.Skip("testdata file not found")
	}

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

	if len(result.SDRICoverageByCategory) == 0 {
		t.Error("Expected at least 1 coverage category")
	}

	for _, c := range result.SDRICoverageByCategory {
		if c.Category == "" {
			t.Errorf("Coverage item has empty Category")
		}
		if c.Coverage < 0 || c.Coverage > 100 {
			t.Errorf("Coverage %s has invalid Coverage: %f", c.Category, c.Coverage)
		}
		if c.Level == "" {
			t.Errorf("Coverage %s has empty Level", c.Category)
		}
	}

	if len(result.SDRICoverageDashboard) != len(result.SDRICoverageByCategory) {
		t.Errorf("Coverage dashboard count %d != coverage by category count %d",
			len(result.SDRICoverageDashboard), len(result.SDRICoverageByCategory))
	}
}

func TestSDRIComplianceAlignment(t *testing.T) {
	archFile := "testdata/attack_paths/healthcare_phi.yaml"
	if _, err := os.Stat(archFile); os.IsNotExist(err) {
		t.Skip("testdata file not found")
	}

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

	// healthcare_phi.yaml declares HIPAA, SOC2, ISO27001
	for _, m := range result.SDRIComplianceAlignments {
		if m.Status == "" {
			t.Errorf("Compliance %s has empty Status", m.Framework)
		}
		if m.Coverage < 0 || m.Coverage > 100 {
			t.Errorf("Compliance %s has invalid Coverage: %f", m.Framework, m.Coverage)
		}
	}
}

func TestSDRIRemediationPriorities(t *testing.T) {
	archFile := "testdata/attack_paths/healthcare_phi.yaml"
	if _, err := os.Stat(archFile); os.IsNotExist(err) {
		t.Skip("testdata file not found")
	}

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

	if len(result.SDRIRemediations) == 0 {
		t.Fatal("Expected at least 1 remediation")
	}

	for i := 1; i < len(result.SDRIRemediations); i++ {
		if result.SDRIRemediations[i].RiskScore > result.SDRIRemediations[i-1].RiskScore {
			t.Errorf("Remediations not sorted by risk: %s (%.2f) > %s (%.2f) at index %d",
				result.SDRIRemediations[i].Description, result.SDRIRemediations[i].RiskScore,
				result.SDRIRemediations[i-1].Description, result.SDRIRemediations[i-1].RiskScore,
				i)
		}
	}
}

func TestSDRIArchitecturalWeaknesses(t *testing.T) {
	archFile := "testdata/attack_paths/healthcare_phi.yaml"
	if _, err := os.Stat(archFile); os.IsNotExist(err) {
		t.Skip("testdata file not found")
	}

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

	if len(result.SDRIAchitecturalWeaknesses) == 0 {
		t.Error("Expected at least 1 architectural weakness (flat network, missing RBAC, etc.)")
	}

	for _, w := range result.SDRIAchitecturalWeaknesses {
		if w.ID == "" {
			t.Errorf("Weakness has empty ID")
		}
		if w.Pattern == "" {
			t.Errorf("Weakness %s has empty Pattern", w.ID)
		}
		if w.Severity == "" {
			t.Errorf("Weakness %s has empty Severity", w.ID)
		}
		if w.Impact == "" {
			t.Errorf("Weakness %s has empty Impact", w.ID)
		}
	}
}

func TestSDRINoCrashOnEmptyArchitecture(t *testing.T) {
	engine := NewEngine(&Config{})
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()
	result, err := engine.RunAnalysis("", "", ModeASFOnly, progress)
	if err != nil {
		return
	}
	if result.SDRISummary == "" {
		t.Log("SDRI summary empty for empty architecture (expected)")
	}
}

func TestSDRIControlsHaveRequiredFields(t *testing.T) {
	archFile := "testdata/attack_paths/healthcare_phi.yaml"
	if _, err := os.Stat(archFile); os.IsNotExist(err) {
		t.Skip("testdata file not found")
	}

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

	for _, c := range result.SDRIControls {
		if c.ID == "" {
			t.Errorf("Control has empty ID")
		}
		if c.Name == "" {
			t.Errorf("Control has empty Name")
		}
		if c.Category == "" {
			t.Errorf("Control %s has empty Category", c.ID)
		}
		if c.Status == "" {
			t.Errorf("Control %s has empty Status", c.ID)
		}
		if c.Strength == "" {
			t.Errorf("Control %s has empty Strength", c.ID)
		}
	}
}
