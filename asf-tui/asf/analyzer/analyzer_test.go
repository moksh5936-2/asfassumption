package analyzer

import (
	"testing"
)

func TestAnalyzeWithSampleData(t *testing.T) {
	a := New()

	docs := []string{"../../../sample_data/finance_policy.txt"}
	evs := []string{
		"../../../sample_data/mfa_status.csv",
		"../../../sample_data/payroll_acl.csv",
		"../../../sample_data/network_exposure.csv",
		"../../../sample_data/backup_config.csv",
	}

	result, err := a.Analyze(docs, evs)
	if err != nil {
		t.Fatalf("Analysis failed: %v", err)
	}

	if result == nil {
		t.Fatal("Expected non-nil result")
	}

	s := result.Result.BuildSummary()
	t.Logf("Summary: claims=%d assumptions=%d verified=%d contradicted=%d unknown=%d critical=%d",
		s.ClaimsFound, s.Assumptions, s.Verified, s.Contradicted, s.Unknown, s.CriticalGaps)
	t.Logf("Evidence: %d records", len(result.Result.Evidence))
	t.Logf("Graph: %d nodes, %d edges", result.Graph.NodeCount, result.Graph.EdgeCount)

	if s.ClaimsFound == 0 {
		t.Error("Expected at least 1 claim")
	}
	if s.Assumptions == 0 {
		t.Error("Expected at least 1 assumption")
	}

	for i, a := range result.Result.Assumptions {
		t.Logf("Assumption %d: type=%s confidence=%.2f status=%s",
			i, a.AssumptionType, a.Confidence, a.VerificationStatus)
	}

	for i, v := range result.Result.Verifications {
		t.Logf("Verification %d: result=%s confidence=%.2f evidence=%v reasoning=%s",
			i, v.Result, v.Confidence, v.EvidenceUsed, v.Reasoning)
	}
}

func TestAnalyzeEmptyDocument(t *testing.T) {
	a := New()
	result, err := a.Analyze([]string{}, []string{})
	if err != nil {
		t.Fatalf("Analysis failed: %v", err)
	}
	if result.Result.ClaimsFound() != 0 {
		t.Errorf("Expected 0 claims, got %d", result.Result.ClaimsFound())
	}
}
