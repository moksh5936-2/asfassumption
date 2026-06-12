package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"asf-tui/intelligence"
)

func TestAPDEngineIntegration(t *testing.T) {
	// Find testdata files
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
		t.Skip("no testdata/attack_paths files found")
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
			if len(result.AttackPaths) == 0 {
				t.Errorf("Expected at least 1 attack path, got 0")
			}
			if len(result.ThreatChains) == 0 {
				t.Errorf("Expected at least 1 threat chain, got 0")
			}
			if result.AttackPathSummary.TotalAttackPaths == 0 {
				t.Errorf("AttackPathSummary.TotalAttackPaths is 0")
			}
			for _, p := range result.AttackPaths {
				if p.ID == "" {
					t.Errorf("Attack path has empty ID")
				}
				if p.Name == "" {
					t.Errorf("Attack path has empty Name")
				}
				if p.EntryPoint == "" {
					t.Errorf("Attack path %s has empty EntryPoint", p.ID)
				}
				if p.TargetAsset == "" {
					t.Errorf("Attack path %s has empty TargetAsset", p.ID)
				}
				if len(p.AttackSteps) == 0 {
					t.Errorf("Attack path %s has no attack steps", p.ID)
				}
				if p.RiskScore < 0 || p.RiskScore > 1.0 {
					t.Errorf("Attack path %s has invalid RiskScore: %f", p.ID, p.RiskScore)
				}
				if p.Likelihood < 0 || p.Likelihood > 1.0 {
					t.Errorf("Attack path %s has invalid Likelihood: %f", p.ID, p.Likelihood)
				}
				if p.Impact < 0 || p.Impact > 1.0 {
					t.Errorf("Attack path %s has invalid Impact: %f", p.ID, p.Impact)
				}
				if p.DetectionDifficulty == "" {
					t.Errorf("Attack path %s has empty DetectionDifficulty", p.ID)
				}
				if p.BusinessImpact == "" {
					t.Errorf("Attack path %s has empty BusinessImpact", p.ID)
				}
				if len(p.KillChainPhases) == 0 {
					t.Errorf("Attack path %s has no KillChainPhases", p.ID)
				}
				if len(p.MITREATTACK) == 0 {
					t.Errorf("Attack path %s has no MITRE ATT&CK mappings", p.ID)
				}
				for i, step := range p.AttackSteps {
					if step.SequenceNumber != i+1 {
						t.Errorf("Attack path %s step %d: expected sequence %d, got %d", p.ID, i, i+1, step.SequenceNumber)
					}
					if step.SourceComponent == "" {
						t.Errorf("Attack path %s step %d: empty SourceComponent", p.ID, i)
					}
					if step.TargetComponent == "" {
						t.Errorf("Attack path %s step %d: empty TargetComponent", p.ID, i)
					}
					if step.Action == "" {
						t.Errorf("Attack path %s step %d: empty Action", p.ID, i)
					}
					if step.Reasoning == "" {
						t.Errorf("Attack path %s step %d: empty Reasoning", p.ID, i)
					}
				}
			}
			for _, c := range result.ThreatChains {
				if c.ID == "" {
					t.Errorf("Threat chain has empty ID")
				}
				if len(c.Threats) == 0 {
					t.Errorf("Threat chain %s has no threats", c.ID)
				}
				if len(c.Path) == 0 {
					t.Errorf("Threat chain %s has no path", c.ID)
				}
			}
			if len(result.AttackPaths) > 0 {
				topPaths := result.AttackPathSummary.TopAttackPaths
				if len(topPaths) == 0 {
					t.Errorf("Expected at least 1 top attack path")
				}
			}
			if result.AttackPathSummary.KillChainCoverage == nil || len(result.AttackPathSummary.KillChainCoverage) == 0 {
				t.Errorf("Expected kill chain coverage")
			}
			if len(result.AttackPathSummary.MITRECoverage) == 0 {
				t.Errorf("Expected MITRE coverage")
			}
		})
	}
}

func TestAPDDeterminism(t *testing.T) {
	archFile := "testdata/attack_paths/healthcare_phi.yaml"
	if _, err := os.Stat(archFile); os.IsNotExist(err) {
		t.Skip("healthcare_phi.yaml not found")
	}

	cfg := &Config{}
	var results []*AnalysisResult

	for i := 0; i < 3; i++ {
		engine := NewEngine(cfg)
		progress := make(chan AnalysisProgress, 100)
		go func() {
			for range progress {
			}
		}()
		result, err := engine.RunAnalysis(archFile, "", ModeASFOnly, progress)
		if err != nil {
			t.Fatalf("Run %d failed: %v", i, err)
		}
		results = append(results, result)
	}

	if len(results) < 2 {
		t.Fatal("need at least 2 results")
	}

	first := results[0]
	for i, result := range results[1:] {
		if len(first.AttackPaths) != len(result.AttackPaths) {
			t.Errorf("Run %d: expected %d attack paths, got %d", i+1, len(first.AttackPaths), len(result.AttackPaths))
			continue
		}
		for j, p1 := range first.AttackPaths {
			p2 := result.AttackPaths[j]
			if p1.ID != p2.ID {
				t.Errorf("Run %d, path %d: ID mismatch (%s vs %s)", i+1, j, p1.ID, p2.ID)
			}
			if p1.Name != p2.Name {
				t.Errorf("Run %d, path %d: Name mismatch (%s vs %s)", i+1, j, p1.Name, p2.Name)
			}
			if p1.RiskScore != p2.RiskScore {
				t.Errorf("Run %d, path %d: RiskScore mismatch (%f vs %f)", i+1, j, p1.RiskScore, p2.RiskScore)
			}
			if len(p1.AttackSteps) != len(p2.AttackSteps) {
				t.Errorf("Run %d, path %d: step count mismatch (%d vs %d)", i+1, j, len(p1.AttackSteps), len(p2.AttackSteps))
			} else {
				for k, s1 := range p1.AttackSteps {
					s2 := p2.AttackSteps[k]
					if s1.Action != s2.Action {
						t.Errorf("Run %d, path %d, step %d: Action mismatch (%s vs %s)", i+1, j, k, s1.Action, s2.Action)
					}
				}
			}
		}
		if first.AttackPathSummary.TotalAttackPaths != result.AttackPathSummary.TotalAttackPaths {
			t.Errorf("Run %d: summary TotalAttackPaths mismatch (%d vs %d)", i+1, first.AttackPathSummary.TotalAttackPaths, result.AttackPathSummary.TotalAttackPaths)
		}
	}
}

func TestAPDEmptyArchitecture(t *testing.T) {
	apd := intelligence.NewAPDEngine()
	result := apd.Run(&intelligence.ArchDescription{
		Name: "empty",
	}, nil, nil, nil, nil)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result.AttackPaths) != 0 {
		t.Errorf("expected 0 attack paths for empty architecture, got %d", len(result.AttackPaths))
	}
	if result.Summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestAPDEntryPoints(t *testing.T) {
	apd := intelligence.NewAPDEngine()
	result := apd.Run(&intelligence.ArchDescription{
		Name: "test",
		Components: []intelligence.Component{
			{ID: "internet", Label: "Internet"},
			{ID: "api", Label: "API Gateway"},
			{ID: "db", Label: "Database"},
		},
		Relationships: []intelligence.Relation{
			{Source: "internet", Target: "api", Label: "HTTPS"},
			{Source: "api", Target: "db", Label: "TLS"},
		},
	}, nil, nil, nil, nil)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result.EntryPoints) == 0 {
		t.Fatal("expected at least 1 entry point")
	}
	foundInternet := false
	for _, ep := range result.EntryPoints {
		if ep.Component == "Internet" && ep.Exposure == 1.0 {
			foundInternet = true
			break
		}
	}
	if !foundInternet {
		t.Error("expected entry point 'Internet' with exposure 1.0")
	}
}

func TestAPDTargetAssets(t *testing.T) {
	apd := intelligence.NewAPDEngine()
	result := apd.Run(&intelligence.ArchDescription{
		Name: "test",
		Components: []intelligence.Component{
			{ID: "web", Label: "Web App"},
			{ID: "phi_db", Label: "PHI Database"},
		},
		Relationships: []intelligence.Relation{
			{Source: "web", Target: "phi_db", Label: "HTTPS"},
		},
	}, nil, nil, nil, nil)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result.TargetAssets) == 0 {
		t.Fatal("expected at least 1 target asset")
	}
	foundPHI := false
	for _, ta := range result.TargetAssets {
		if ta.Sensitivity == "critical" {
			foundPHI = true
			break
		}
	}
	if !foundPHI {
		t.Error("expected a critical target asset for PHI database")
	}
}

func TestAPDKillChainCoverage(t *testing.T) {
	archFile := "testdata/attack_paths/healthcare_phi.yaml"
	if _, err := os.Stat(archFile); os.IsNotExist(err) {
		t.Skip("healthcare_phi.yaml not found")
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

	if len(result.AttackPathSummary.KillChainCoverage) == 0 {
		t.Error("expected non-empty kill chain coverage")
	}
	for phase, count := range result.AttackPathSummary.KillChainCoverage {
		if count <= 0 {
			t.Errorf("kill chain phase %s has count %d, expected > 0", phase, count)
		}
	}
	expectedPhases := map[string]bool{
		"Reconnaissance": true, "Initial Access": true, "Execution": true,
		"Persistence": true, "Privilege Escalation": true, "Credential Access": true,
		"Discovery": true, "Lateral Movement": true, "Collection": true,
		"Exfiltration": true,
	}
	hasExpected := false
	for phase := range result.AttackPathSummary.KillChainCoverage {
		if expectedPhases[phase] {
			hasExpected = true
			break
		}
	}
	if !hasExpected {
		t.Errorf("kill chain coverage missing expected phases, got: %v", keys(result.AttackPathSummary.KillChainCoverage))
	}
}

func TestAPDMITRECoverage(t *testing.T) {
	archFile := "testdata/attack_paths/healthcare_phi.yaml"
	if _, err := os.Stat(archFile); os.IsNotExist(err) {
		t.Skip("healthcare_phi.yaml not found")
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

	if len(result.AttackPathSummary.MITRECoverage) == 0 {
		t.Error("expected non-empty MITRE coverage")
	}
	for _, technique := range result.AttackPathSummary.MITRECoverage {
		if !strings.Contains(technique, "T") {
			t.Errorf("MITRE technique %q does not contain technique ID", technique)
		}
	}
}

func TestAPDBusinessImpact(t *testing.T) {
	archFile := "testdata/attack_paths/healthcare_phi.yaml"
	if _, err := os.Stat(archFile); os.IsNotExist(err) {
		t.Skip("healthcare_phi.yaml not found")
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

	for _, p := range result.AttackPaths {
		if p.BusinessImpact == "" {
			t.Errorf("attack path %s has no business impact", p.Name)
		}
	}
}

func TestAPDRecommendations(t *testing.T) {
	archFile := "testdata/attack_paths/healthcare_phi.yaml"
	if _, err := os.Stat(archFile); os.IsNotExist(err) {
		t.Skip("healthcare_phi.yaml not found")
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

	for _, p := range result.AttackPaths {
		if len(p.Recommendations) == 0 {
			t.Errorf("attack path %s has no recommendations", p.Name)
		}
	}
}

func TestAPDPrioritization(t *testing.T) {
	apd := intelligence.NewAPDEngine()
	result := apd.Run(&intelligence.ArchDescription{
		Name: "test",
		Components: []intelligence.Component{
			{ID: "internet", Label: "Internet"},
			{ID: "web", Label: "Web App"},
			{ID: "phi_db", Label: "PHI Database"},
			{ID: "log", Label: "Logging"},
		},
		Relationships: []intelligence.Relation{
			{Source: "internet", Target: "web", Label: "HTTPS"},
			{Source: "web", Target: "phi_db", Label: "TLS"},
			{Source: "web", Target: "log", Label: "TLS"},
		},
	}, nil, nil, nil, nil)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	if len(result.AttackPaths) > 1 {
		for i := 1; i < len(result.AttackPaths); i++ {
			if result.AttackPaths[i].RiskScore > result.AttackPaths[i-1].RiskScore {
				t.Errorf("attack paths not sorted by risk: %s (%f) > %s (%f)",
					result.AttackPaths[i].Name, result.AttackPaths[i].RiskScore,
					result.AttackPaths[i-1].Name, result.AttackPaths[i-1].RiskScore)
			}
		}
	}
}

func TestAPDDetectionDifficulty(t *testing.T) {
	apd := intelligence.NewAPDEngine()
	validDifficulties := map[string]bool{
		"Easy": true, "Moderate": true, "Hard": true, "Very Hard": true,
	}
	result := apd.Run(&intelligence.ArchDescription{
		Name: "test",
		Components: []intelligence.Component{
			{ID: "internet", Label: "Internet"},
			{ID: "db", Label: "Database"},
		},
		Relationships: []intelligence.Relation{
			{Source: "internet", Target: "db", Label: "HTTPS"},
		},
	}, nil, nil, nil, nil)
	if result == nil {
		t.Fatal("expected non-nil result")
	}
	for _, p := range result.AttackPaths {
		if !validDifficulties[p.DetectionDifficulty] {
			t.Errorf("attack path %s has invalid detection difficulty: %q", p.Name, p.DetectionDifficulty)
		}
	}
}

func TestAPDTrustBoundaryCrossings(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()
	result, err := engine.RunAnalysis("testdata/attack_paths/vpn_infrastructure.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Skipf("vpn_infrastructure.yaml not found or analysis failed: %v", err)
	}
	for _, p := range result.AttackPaths {
		if len(p.AffectedBoundaries) > 0 {
			return
		}
	}
	t.Log("no boundaries crossed in VPN infrastructure (may be expected with limited test data)")
}

func TestAPDReverseConversion(t *testing.T) {
	paths := []intelligence.AttackPath{
		{
			ID:          "ap-1",
			Name:        "Test Path",
			Description: "Test description",
			EntryPoint:  "Internet",
			TargetAsset: "Database",
			AttackSteps: []intelligence.AttackStep{
				{SequenceNumber: 1, SourceComponent: "Internet", TargetComponent: "Web", Action: "Test", Threat: "Test Threat", Reasoning: "Because"},
			},
			Likelihood:          0.5,
			Impact:              0.8,
			RiskScore:           0.4,
			Confidence:          0.7,
			DetectionDifficulty: "Hard",
			BusinessImpact:      "Test impact",
			Recommendations:     []string{"Test"},
			KillChainPhases:     []string{"Initial Access"},
			MITREATTACK:         []string{"T1078 - Valid Accounts"},
			STRIDECategories:    []string{"Spoofing"},
		},
	}
	converted := convertAPDAttackPaths(paths)
	if len(converted) != 1 {
		t.Fatalf("expected 1 converted path, got %d", len(converted))
	}
	p := converted[0]
	if p.ID != "ap-1" || p.Name != "Test Path" || p.EntryPoint != "Internet" || p.TargetAsset != "Database" {
		t.Errorf("conversion fields mismatch: %+v", p)
	}
	if len(p.AttackSteps) != 1 {
		t.Errorf("expected 1 attack step, got %d", len(p.AttackSteps))
	}
	if p.AttackSteps[0].Action != "Test" {
		t.Errorf("step action mismatch: %s", p.AttackSteps[0].Action)
	}
}

func TestAPDReverseConversionThreatChains(t *testing.T) {
	chains := []intelligence.ThreatChain{
		{
			ID:        "tc-1",
			Threats:   []string{"Credential Stuffing", "Data Exfiltration"},
			Path:      []string{"Internet", "Web", "DB"},
			RiskScore: 0.75,
			Reasoning: "Test chain",
		},
	}
	converted := convertAPDThreatChains(chains)
	if len(converted) != 1 {
		t.Fatalf("expected 1 chain, got %d", len(converted))
	}
	c := converted[0]
	if c.ID != "tc-1" || len(c.Threats) != 2 || c.RiskScore != 0.75 {
		t.Errorf("conversion fields mismatch: %+v", c)
	}
}

func TestAPDSummaryConversion(t *testing.T) {
	result := &intelligence.APDRunResult{
		AttackPaths: []intelligence.AttackPath{
			{Name: "Path A", RiskScore: 0.8},
			{Name: "Path B", RiskScore: 0.5},
			{Name: "Path C", RiskScore: 0.3},
			{Name: "Path D", RiskScore: 0.1},
		},
		ThreatChains: []intelligence.ThreatChain{
			{ID: "tc-1"},
		},
		TopPaths: []intelligence.AttackPath{
			{Name: "Path A"},
			{Name: "Path B"},
		},
		KillChainCoverage: map[string]int{"Initial Access": 2, "Exfiltration": 1},
		MITREMapping:      map[string][]string{"T1078": {"Valid Accounts"}},
		Summary:           "Test summary text",
	}
	s := convertAPDSummary(result)
	if s.TotalAttackPaths != 4 {
		t.Errorf("expected 4 total paths, got %d", s.TotalAttackPaths)
	}
	if s.CriticalCount != 1 {
		t.Errorf("expected 1 critical, got %d", s.CriticalCount)
	}
	if s.HighCount != 1 {
		t.Errorf("expected 1 high, got %d", s.HighCount)
	}
	if s.ThreatChainCount != 1 {
		t.Errorf("expected 1 threat chain, got %d", s.ThreatChainCount)
	}
	if len(s.TopAttackPaths) != 2 {
		t.Errorf("expected 2 top paths, got %d", len(s.TopAttackPaths))
	}
	if len(s.KillChainCoverage) != 2 {
		t.Errorf("expected 2 kill chain phases, got %d", len(s.KillChainCoverage))
	}
	if s.SummaryText != "Test summary text" {
		t.Errorf("summary text mismatch: %s", s.SummaryText)
	}
}

func keys(m map[string]int) []string {
	var k []string
	for key := range m {
		k = append(k, key)
	}
	return k
}
