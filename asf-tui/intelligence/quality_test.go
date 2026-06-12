package intelligence

import (
	"strings"
	"testing"
)

func TestQualityEngineScore(t *testing.T) {
	qe := NewQualityEngine()
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
		Risk:        RiskHigh,
		Impact:      4,
		Confidence:  0.80,
		SourceType:  "inferred",
	}
	qs := qe.Score(a, arch)

	if qs.Hiddenness < 0.5 {
		t.Errorf("expected hiddenness > 0.5 for inferred assumption, got %.2f", qs.Hiddenness)
	}
	if qs.Impact < 0.5 {
		t.Errorf("expected impact > 0.5 for impact=4, got %.2f", qs.Impact)
	}
	if qs.Novelty < 0.5 {
		t.Errorf("expected novelty > 0.5 for inferred assumption, got %.2f", qs.Novelty)
	}
	if qs.ArchitecturalRelevance < 0.5 {
		t.Errorf("expected relevance > 0.5 for matching component, got %.2f", qs.ArchitecturalRelevance)
	}
	if qs.Risk < 0.5 {
		t.Errorf("expected risk > 0.5 for High risk, got %.2f", qs.Risk)
	}
	if qs.Confidence < 0.5 {
		t.Errorf("expected confidence > 0.5 for 0.80 confidence, got %.2f", qs.Confidence)
	}
	if qs.Total <= 0 {
		t.Errorf("expected total > 0, got %.2f", qs.Total)
	}
	if qs.Reason == "" {
		t.Error("expected non-empty reason")
	}
}

func TestQualityEngineScoreGenericLow(t *testing.T) {
	qe := NewQualityEngine()
	arch := &ArchDescription{
		Name: "test",
		Components: []Component{
			{ID: "app1", Label: "App"},
		},
	}
	a := Assumption{
		ID:          "A1",
		Description: "All communication uses TLS encryption",
		Component:   "General",
		Category:    "NetworkSegmentation",
		Risk:        RiskLow,
		Impact:      2,
		Confidence:  0.60,
		SourceType:  "explicit",
	}
	qs := qe.Score(a, arch)
	if qs.Hiddenness > 0.3 {
		t.Errorf("expected low hiddenness for generic TLS, got %.2f", qs.Hiddenness)
	}
	if qs.Novelty > 0.3 {
		t.Errorf("expected low novelty for generic TLS, got %.2f", qs.Novelty)
	}
}

func TestQualityEngineRank(t *testing.T) {
	qe := NewQualityEngine()
	arch := &ArchDescription{
		Name: "test",
		Components: []Component{
			{ID: "db1", Label: "Database"},
		},
	}
	assumptions := []Assumption{
		{ID: "A1", Description: "Generic TLS", Risk: RiskLow, Impact: 2, Confidence: 0.60, SourceType: "explicit"},
		{ID: "A2", Description: "PHI key management missing", Risk: RiskCritical, Impact: 5, Confidence: 0.90, SourceType: "inferred", Component: "Database"},
		{ID: "A3", Description: "Session rotation missing", Risk: RiskHigh, Impact: 4, Confidence: 0.80, SourceType: "inferred"},
	}
	ranked := qe.Rank(assumptions, arch)
	if ranked[0].ID != "A2" {
		t.Errorf("expected A2 (Critical) first, got %s", ranked[0].ID)
	}
	if ranked[len(ranked)-1].ID != "A1" {
		t.Errorf("expected A1 (Low, generic) last, got %s", ranked[len(ranked)-1].ID)
	}
}

func TestQualityEngineGetTopAssumptions(t *testing.T) {
	qe := NewQualityEngine()
	arch := &ArchDescription{
		Name: "test",
		Components: []Component{
			{ID: "db1", Label: "Database"},
		},
	}
	assumptions := []Assumption{
		{ID: "A1", Description: "Generic TLS", Risk: RiskLow, Impact: 2, Confidence: 0.60, SourceType: "explicit"},
		{ID: "A2", Description: "PHI key management missing", Risk: RiskCritical, Impact: 5, Confidence: 0.90, SourceType: "inferred", Component: "Database"},
		{ID: "A3", Description: "Session rotation missing", Risk: RiskHigh, Impact: 4, Confidence: 0.80, SourceType: "inferred"},
		{ID: "A4", Description: "Backup encryption missing", Risk: RiskHigh, Impact: 4, Confidence: 0.85, SourceType: "inferred"},
	}
	top := qe.GetTopAssumptions(assumptions, arch, 2)
	if len(top) != 2 {
		t.Errorf("expected 2 top assumptions, got %d", len(top))
	}
	if top[0].ID != "A2" {
		t.Errorf("expected A2 first, got %s", top[0].ID)
	}
}

func TestQualityEngineAverageQuality(t *testing.T) {
	qe := NewQualityEngine()
	arch := &ArchDescription{
		Name: "test",
		Components: []Component{
			{ID: "db1", Label: "Database"},
		},
	}
	assumptions := []Assumption{
		{ID: "A1", Description: "Generic TLS", Risk: RiskLow, Impact: 2, Confidence: 0.60, SourceType: "explicit"},
		{ID: "A2", Description: "PHI key management missing", Risk: RiskCritical, Impact: 5, Confidence: 0.90, SourceType: "inferred", Component: "Database"},
	}
	avg := qe.AverageQuality(assumptions, arch)
	if avg <= 0 {
		t.Errorf("expected average quality > 0, got %.2f", avg)
	}
	if avg >= 1 {
		t.Errorf("expected average quality < 1, got %.2f", avg)
	}
}

func TestQualityEngineAverageQualityEmpty(t *testing.T) {
	qe := NewQualityEngine()
	avg := qe.AverageQuality(nil, nil)
	if avg != 0 {
		t.Errorf("expected 0 for empty assumptions, got %.2f", avg)
	}
}

func TestQualityEngineQualityReport(t *testing.T) {
	qe := NewQualityEngine()
	arch := &ArchDescription{
		Name: "test",
		Components: []Component{
			{ID: "db1", Label: "Database"},
		},
	}
	assumptions := []Assumption{
		{ID: "A1", Description: "Generic TLS", Risk: RiskLow, Impact: 2, Confidence: 0.60, SourceType: "explicit"},
		{ID: "A2", Description: "PHI key management missing", Risk: RiskCritical, Impact: 5, Confidence: 0.90, SourceType: "inferred", Component: "Database"},
	}
	report := qe.QualityReport(assumptions, arch)
	if report == "no assumptions to score" {
		t.Error("expected non-empty report")
	}
	if !strings.Contains(report, "average quality") {
		t.Error("expected 'average quality' in report")
	}
	if !strings.Contains(report, "A2") {
		t.Error("expected top assumption ID in report")
	}
}

func TestQualityEngineScoreRiskLevel(t *testing.T) {
	qe := NewQualityEngine()
	tests := []struct {
		risk RiskLevel
		want float64
	}{
		{RiskCritical, 1.0},
		{RiskHigh, 0.8},
		{RiskMedium, 0.5},
		{RiskLow, 0.2},
		{"Unknown", 0.5},
	}
	for _, tt := range tests {
		got := qe.scoreRiskLevel(tt.risk)
		if got != tt.want {
			t.Errorf("scoreRiskLevel(%s) = %.2f, want %.2f", tt.risk, got, tt.want)
		}
	}
}
