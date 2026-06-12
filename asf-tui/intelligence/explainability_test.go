package intelligence

import (
	"strings"
	"testing"
)

func TestExplainabilityEngineExplainWhy(t *testing.T) {
	te := NewTaxonomyEngine()
	ee := NewExplainabilityEngine(te)
	arch := &ArchDescription{
		Name: "test",
		Components: []Component{
			{ID: "db1", Label: "Database"},
		},
	}
	a := Assumption{
		ID:          "A1",
		Description: "Database does not specify encryption at rest",
		Component:   "Database",
		Category:    "DataProtection",
		Confidence:  0.80,
		Keywords:    []string{"database", "encryption"},
	}
	copyA := a
	ee.ExplainWhy(&copyA, arch)

	if copyA.Rationale == "" {
		t.Error("expected rationale after ExplainWhy")
	}
	if !strings.Contains(copyA.Rationale, "Database does not specify encryption at rest") {
		t.Error("expected description in rationale")
	}
	if !strings.Contains(copyA.Rationale, "80%") {
		t.Error("expected confidence in rationale")
	}
	if !strings.Contains(copyA.Rationale, "DataProtection") {
		t.Error("expected category in rationale")
	}
	if len(copyA.EvidenceSources) == 0 {
		t.Error("expected evidence sources after ExplainWhy")
	}
}

func TestExplainabilityEngineExplainWhyNil(t *testing.T) {
	te := NewTaxonomyEngine()
	ee := NewExplainabilityEngine(te)
	// Should not panic
	ee.ExplainWhy(nil, nil)
}

func TestExplainabilityEngineExplainAll(t *testing.T) {
	te := NewTaxonomyEngine()
	ee := NewExplainabilityEngine(te)
	arch := &ArchDescription{
		Name: "test",
		Components: []Component{
			{ID: "db1", Label: "Database"},
		},
	}
	assumptions := []Assumption{
		{ID: "A1", Description: "Database does not specify encryption", Component: "Database", Category: "DataProtection", Confidence: 0.80},
		{ID: "A2", Description: "App does not specify auth", Component: "App", Category: "Authentication", Confidence: 0.70},
	}
	result := ee.ExplainAll(assumptions, arch)
	if len(result) != 2 {
		t.Fatalf("expected 2 assumptions, got %d", len(result))
	}
	for _, a := range result {
		if a.Rationale == "" {
			t.Errorf("expected rationale for %s", a.ID)
		}
	}
}

func TestExplainabilityEngineGatherEvidence(t *testing.T) {
	te := NewTaxonomyEngine()
	ee := NewExplainabilityEngine(te)
	arch := &ArchDescription{
		Name: "test",
		Components: []Component{
			{ID: "db1", Label: "Database"},
			{ID: "gw1", Label: "Gateway"},
		},
		Relationships: []Relation{
			{Source: "Gateway", Target: "Database", Label: "SQL"},
		},
	}
	a := Assumption{
		ID:          "A1",
		Description: "Database does not specify encryption",
		Component:   "Database",
		Keywords:    []string{"database", "sql"},
	}
	evidence := ee.gatherEvidence(&a, arch)
	foundDB := false
	foundRel := false
	for _, e := range evidence {
		if strings.Contains(e, "Database") {
			foundDB = true
		}
		if strings.Contains(e, "Gateway") && strings.Contains(e, "Database") {
			foundRel = true
		}
	}
	if !foundDB {
		t.Error("expected Database component in evidence")
	}
	if !foundRel {
		t.Error("expected Gateway->Database relationship in evidence")
	}
}

func TestExplainabilityEngineIdentifyMissingControls(t *testing.T) {
	te := NewTaxonomyEngine()
	ee := NewExplainabilityEngine(te)
	a := Assumption{
		Description: "Database does not specify encryption at rest",
		Category:    "DataProtection",
	}
	missing := ee.identifyMissingControls(&a)
	if len(missing) == 0 {
		t.Error("expected missing controls")
	}
	found := false
	for _, m := range missing {
		if strings.Contains(m, "encryption") || strings.Contains(m, "specify") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected encryption-related missing control, got %v", missing)
	}
}

func TestExplainabilityEngineIdentifyMissingControlsFallback(t *testing.T) {
	te := NewTaxonomyEngine()
	ee := NewExplainabilityEngine(te)
	a := Assumption{
		Description: "Something is configured",
		Category:    "Authentication",
	}
	missing := ee.identifyMissingControls(&a)
	if len(missing) == 0 {
		t.Error("expected fallback missing controls from taxonomy")
	}
}

func TestExplainabilityEngineBuildArchitectureContext(t *testing.T) {
	te := NewTaxonomyEngine()
	ee := NewExplainabilityEngine(te)
	arch := &ArchDescription{
		Name: "test",
		Components: []Component{
			{ID: "db1", Label: "Database"},
		},
		Relationships: []Relation{
			{Source: "App", Target: "Database"},
		},
	}
	a := Assumption{Component: "Database"}
	ctx := ee.buildArchitectureContext(&a, arch)
	if !strings.Contains(ctx, "test") {
		t.Error("expected architecture name in context")
	}
	if !strings.Contains(ctx, "Database") {
		t.Error("expected component in context")
	}
	if !strings.Contains(ctx, "total components: 1") {
		t.Error("expected component count in context")
	}
	if !strings.Contains(ctx, "total relationships: 1") {
		t.Error("expected relationship count in context")
	}
}

func TestExplainabilityEngineBuildSummaryExplanation(t *testing.T) {
	te := NewTaxonomyEngine()
	ee := NewExplainabilityEngine(te)
	summary := ee.BuildSummaryExplanation(
		[]Assumption{{ID: "A1"}, {ID: "A2"}},
		[]Contradiction{{RuleName: "C1"}},
		[]TrustBoundary{{Type: "Internet"}},
		"Healthcare",
	)
	if !strings.Contains(summary, "2 assumptions") {
		t.Error("expected assumption count in summary")
	}
	if !strings.Contains(summary, "Healthcare") {
		t.Error("expected domain in summary")
	}
	if !strings.Contains(summary, "1 contradictions") {
		t.Error("expected contradiction count in summary")
	}
	if !strings.Contains(summary, "1 trust boundaries") {
		t.Error("expected boundary count in summary")
	}
}
