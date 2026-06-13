package review

import (
	"fmt"
	"strings"
	"testing"
)

func TestEmptyReview(t *testing.T) {
	engine := NewReviewEngine("general", nil)
	output := engine.RunAll()
	if output == nil {
		t.Fatal("expected non-nil output")
	}
}

func TestSingleAssumptionReview(t *testing.T) {
	inputs := []ReviewInput{
		{
			AssumptionID:           "A1",
			AssumptionText:         "MFA is enforced for all users",
			Risk:                   "Critical",
			Category:               "identity",
			Component:              "auth0",
			Centrality:             0.85,
			Criticality:            0.9,
			FailureRadius:          12,
			SupportCount:           18,
			DependencyCount:        12,
			VerificationPriority:   "Critical",
			VerificationConfidence: 45,
			VerificationStatus:     "Unverified",
			CoverageGap:            false,
			BlindSpotScore:         0,
			Domain:                 "general",
		},
	}
	engine := NewReviewEngine("general", inputs)
	output := engine.RunAll()
	if output.Queue == nil || len(output.Queue.Items) != 1 {
		t.Fatalf("expected 1 queue item, got %d", len(output.Queue.Items))
	}
	item := output.Queue.Items[0]
	if item.PriorityScore <= 0 {
		t.Errorf("expected positive priority score, got %.0f", item.PriorityScore)
	}
	if item.Rank != 1 {
		t.Errorf("expected rank 1, got %d", item.Rank)
	}
	if item.WhyReview == "" {
		t.Error("expected why review rationale")
	}
	if item.WhatToReview == "" {
		t.Error("expected what to review")
	}
	if item.ExpectedEvidence == "" {
		t.Error("expected expected evidence")
	}
}

func TestMultipleAssumptions(t *testing.T) {
	inputs := []ReviewInput{
		{
			AssumptionID: "A1", AssumptionText: "MFA enforced", Risk: "Critical",
			Category: "identity", Component: "auth0", Centrality: 0.9, FailureRadius: 15,
			SupportCount: 20, VerificationPriority: "Critical", VerificationConfidence: 30,
		},
		{
			AssumptionID: "A2", AssumptionText: "TLS enabled", Risk: "Critical",
			Category: "cryptography", Component: "", Centrality: 0.5, FailureRadius: 5,
			SupportCount: 8, VerificationPriority: "High", VerificationConfidence: 60,
		},
		{
			AssumptionID: "A3", AssumptionText: "Log retention configured", Risk: "Low",
			Category: "monitoring", Component: "", Centrality: 0.1, FailureRadius: 1,
			SupportCount: 2, VerificationPriority: "Low", VerificationConfidence: 95,
		},
	}
	engine := NewReviewEngine("general", inputs)
	output := engine.RunAll()
	if output.Queue == nil || len(output.Queue.Items) != 3 {
		t.Fatalf("expected 3 items, got %d", len(output.Queue.Items))
	}
	if output.Queue.Items[0].Rank != 1 {
		t.Errorf("expected item 0 to be rank 1, got %d", output.Queue.Items[0].Rank)
	}
	if output.Queue.Items[0].PriorityScore < output.Queue.Items[2].PriorityScore {
		t.Error("expected first item to have highest score")
	}
	if output.Matrix == nil {
		t.Error("expected priority matrix")
	}
}

func TestPriorityScoring(t *testing.T) {
	tests := []struct {
		name     string
		input    ReviewInput
		minScore float64
	}{
		{
			"critical high centrality",
			ReviewInput{AssumptionID: "A1", AssumptionText: "Critical test", Risk: "Critical",
				Category: "identity", Centrality: 0.9, FailureRadius: 15, SupportCount: 20,
				VerificationPriority: "Critical", VerificationConfidence: 20, CoverageGap: true, BlindSpotScore: 90,
				Domain: "healthcare"},
			70,
		},
		{
			"low risk low impact",
			ReviewInput{AssumptionID: "A2", AssumptionText: "Low test", Risk: "Low",
				Category: "operational", Centrality: 0.0, FailureRadius: 0, SupportCount: 0,
				VerificationPriority: "Low", VerificationConfidence: 95, CoverageGap: false, BlindSpotScore: 0,
				Domain: "general"},
			0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := NewReviewEngine(tt.input.Domain, []ReviewInput{tt.input})
			output := engine.RunAll()
			if output.Queue == nil || len(output.Queue.Items) == 0 {
				t.Fatal("expected item")
			}
			if output.Queue.Items[0].PriorityScore < tt.minScore {
				t.Errorf("expected score >= %.0f, got %.0f", tt.minScore, output.Queue.Items[0].PriorityScore)
			}
		})
	}
}

func TestPriorityMatrix(t *testing.T) {
	inputs := make([]ReviewInput, 8)
	for i := 0; i < 8; i++ {
		risk := "Low"
		if i < 3 {
			risk = "Critical"
		}
		cat := "identity"
		if i >= 4 {
			cat = "operational"
		}
		inputs[i] = ReviewInput{
			AssumptionID: fmt.Sprintf("A%d", i), AssumptionText: fmt.Sprintf("Test %d", i),
			Risk: risk, Category: cat, Centrality: float64(i) / 10,
			FailureRadius: i * 2, SupportCount: i * 3,
			VerificationPriority: "High", VerificationConfidence: 50,
		}
	}
	engine := NewReviewEngine("general", inputs)
	output := engine.RunAll()
	if output.Matrix == nil {
		t.Fatal("expected matrix")
	}
	if len(output.Matrix.HighValueLowEffort) == 0 && len(output.Matrix.HighValueHighEffort) == 0 {
		t.Error("expected at least some items in matrix quadrants")
	}
}

func TestDomainPrioritization(t *testing.T) {
	inputs := []ReviewInput{
		{
			AssumptionID: "A1", AssumptionText: "PHI access controls", Risk: "Critical",
			Category: "identity", Centrality: 0.8, FailureRadius: 10, SupportCount: 15,
			VerificationPriority: "Critical", VerificationConfidence: 30, Domain: "healthcare",
		},
		{
			AssumptionID: "A2", AssumptionText: "Log retention", Risk: "Low",
			Category: "monitoring", Centrality: 0.2, FailureRadius: 2, SupportCount: 3,
			VerificationPriority: "Low", VerificationConfidence: 90, Domain: "healthcare",
		},
	}
	engine := NewReviewEngine("healthcare", inputs)
	output := engine.RunAll()
	if output.DomainView == nil {
		t.Fatal("expected domain view for healthcare")
	}
	if output.DomainView.Domain != "healthcare" {
		t.Errorf("expected healthcare domain, got %s", output.DomainView.Domain)
	}
	if len(output.DomainView.FocusAreas) == 0 {
		t.Error("expected focus areas for healthcare")
	}
}

func TestReviewCampaigns(t *testing.T) {
	inputs := make([]ReviewInput, 20)
	for i := 0; i < 20; i++ {
		risk := "Medium"
		if i < 5 {
			risk = "Critical"
		}
		cat := []string{"identity", "authorization", "cryptography", "monitoring", "resilience", "third_party", "operational"}
		inputs[i] = ReviewInput{
			AssumptionID: fmt.Sprintf("A%d", i), AssumptionText: fmt.Sprintf("Assumption %d", i),
			Risk: risk, Category: cat[i%7], Centrality: float64(i%10) / 10,
			FailureRadius: i, SupportCount: i * 2,
			VerificationPriority: "High", VerificationConfidence: float64(100 - i*5),
		}
	}
	engine := NewReviewEngine("general", inputs)
	output := engine.RunAll()
	if len(output.Campaigns) == 0 {
		t.Fatal("expected campaigns")
	}
	found := false
	for _, c := range output.Campaigns {
		if c.Name == "30 Minute Review Plan" && len(c.Items) > 0 {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 30 minute review plan with items")
	}
}

func TestCISODashboard(t *testing.T) {
	inputs := make([]ReviewInput, 15)
	for i := 0; i < 15; i++ {
		inputs[i] = ReviewInput{
			AssumptionID: fmt.Sprintf("A%d", i), AssumptionText: fmt.Sprintf("Assumption %d", i),
			Risk: "High", Category: "identity", Centrality: float64(i%5) / 10,
			FailureRadius: i, SupportCount: i * 2,
			VerificationPriority: "High", VerificationConfidence: float64(50 - i*3),
		}
	}
	engine := NewReviewEngine("general", inputs)
	output := engine.RunAll()
	if output.CISODashboard == nil {
		t.Fatal("expected CISO dashboard")
	}
	if len(output.CISODashboard.HighestRiskAssumptions) == 0 {
		t.Error("expected highest risk assumptions")
	}
	if len(output.CISODashboard.GreatestRiskReduction) == 0 {
		t.Error("expected risk reduction opportunities")
	}
}

func TestExportMarkdown(t *testing.T) {
	inputs := []ReviewInput{
		{
			AssumptionID: "A1", AssumptionText: "MFA enforced", Risk: "Critical",
			Category: "identity", Centrality: 0.8, FailureRadius: 10, SupportCount: 15,
			VerificationPriority: "Critical", VerificationConfidence: 30,
		},
	}
	engine := NewReviewEngine("general", inputs)
	output := engine.RunAll()
	md := ExportMarkdown(output)
	if md == "" {
		t.Error("expected non-empty markdown")
	}
	if !strings.Contains(md, "Security Review Workbench Report") {
		t.Error("expected report title in markdown")
	}
	if !strings.Contains(md, "MFA enforced") {
		t.Error("expected assumption text in markdown")
	}
}

func TestExportHTML(t *testing.T) {
	inputs := []ReviewInput{
		{
			AssumptionID: "A1", AssumptionText: "RBAC configured", Risk: "High",
			Category: "authorization", Centrality: 0.6, FailureRadius: 5, SupportCount: 8,
			VerificationPriority: "High", VerificationConfidence: 60,
		},
	}
	engine := NewReviewEngine("general", inputs)
	output := engine.RunAll()
	html := ExportHTML(output)
	if html == "" {
		t.Error("expected non-empty HTML")
	}
	if !strings.Contains(html, "Security Review Workbench Report") {
		t.Error("expected report title in HTML")
	}
}

func TestReviewPrecision(t *testing.T) {
	inputs := make([]ReviewInput, 10)
	for i := 0; i < 10; i++ {
		cat := []string{"identity", "authorization", "cryptography", "monitoring", "resilience", "third_party", "operational"}
		risk := []string{"Critical", "High", "Medium", "Low"}
		inputs[i] = ReviewInput{
			AssumptionID: fmt.Sprintf("A%d", i), AssumptionText: fmt.Sprintf("Precision test %d", i),
			Risk: risk[i%4], Category: cat[i%7], Centrality: float64(i%5+1) / 10,
			FailureRadius: i + 1, SupportCount: (i + 1) * 2,
			VerificationPriority: risk[i%4], VerificationConfidence: float64(100 - i*10),
		}
	}
	engine := NewReviewEngine("fintech", inputs)
	output := engine.RunAll()
	if output.Queue == nil || len(output.Queue.Items) != 10 {
		t.Fatalf("expected 10 items, got %d", len(output.Queue.Items))
	}
	for _, it := range output.Queue.Items {
		t.Logf("rank=%d score=%.0f risk=%s cat=%s value=%s effort=%s quadrant=%s",
			it.Rank, it.PriorityScore, it.Risk, it.Category,
			string(it.ReviewValue), string(it.ReviewEffort), string(it.Quadrant))
		if it.PriorityScore < 0 || it.PriorityScore > 100 {
			t.Errorf("score out of range: %.0f", it.PriorityScore)
		}
		if it.Rank < 1 || it.Rank > 10 {
			t.Errorf("rank out of range: %d", it.Rank)
		}
	}
}

func BenchmarkReviewEngine(b *testing.B) {
	inputs := make([]ReviewInput, 100)
	for i := 0; i < 100; i++ {
		cat := []string{"identity", "authorization", "cryptography", "monitoring", "resilience", "third_party", "operational"}
		risk := []string{"Critical", "High", "Medium", "Low"}
		inputs[i] = ReviewInput{
			AssumptionID: fmt.Sprintf("A%d", i), AssumptionText: fmt.Sprintf("Bench %d", i),
			Risk: risk[i%4], Category: cat[i%7], Centrality: float64(i%10) / 10,
			FailureRadius: i % 15, SupportCount: i % 20,
			VerificationPriority: risk[i%4], VerificationConfidence: float64(100 - i),
			Domain: "healthcare",
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine := NewReviewEngine("healthcare", inputs)
		engine.RunAll()
	}
}

func BenchmarkReviewLarge(b *testing.B) {
	inputs := make([]ReviewInput, 500)
	for i := 0; i < 500; i++ {
		cat := []string{"identity", "authorization", "cryptography", "monitoring", "resilience", "third_party", "operational"}
		inputs[i] = ReviewInput{
			AssumptionID: fmt.Sprintf("A%d", i), AssumptionText: fmt.Sprintf("Large %d", i),
			Risk: "High", Category: cat[i%7], Centrality: 0.5,
			FailureRadius: 5, SupportCount: 10,
			VerificationPriority: "High", VerificationConfidence: 50,
		}
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine := NewReviewEngine("kubernetes", inputs)
		engine.RunAll()
	}
}
