package confidence

import (
	"testing"

	"asf-tui/asf/models"
)

func TestVerificationConfidence(t *testing.T) {
	ce := NewEngine()
	v := models.NewVerification("asm_1", []string{"evd_1"}, models.VerificationResultCONTRADICTED, 0.92, "Found unauthorized access", nil)
	ev := []models.Evidence{
		models.NewEvidence("test.csv", models.SourceTypeCSV, []map[string]interface{}{{"user": "alice"}}),
	}
	conf := ce.ComputeVerificationConfidence(v, ev)
	if conf <= 0 {
		t.Errorf("Expected positive confidence, got %f", conf)
	}
}

func TestVerificationConfidenceUnknown(t *testing.T) {
	ce := NewEngine()
	v := models.NewVerification("asm_1", nil, models.VerificationResultUNKNOWN, 0.0, "", nil)
	conf := ce.ComputeVerificationConfidence(v, nil)
	if conf != 0.0 {
		t.Errorf("Expected 0.0 for UNKNOWN, got %f", conf)
	}
}

func TestAssumptionConfidence(t *testing.T) {
	ce := NewEngine()
	verifications := []models.Verification{
		models.NewVerification("asm_1", nil, models.VerificationResultVERIFIED, 0.85, "", nil),
		models.NewVerification("asm_2", nil, models.VerificationResultCONTRADICTED, 0.72, "", nil),
	}
	conf := ce.ComputeAssumptionConfidence(verifications)
	if conf < 0 || conf > 1.0 {
		t.Errorf("Confidence out of range: %f", conf)
	}
}

func TestAssumptionConfidenceEmpty(t *testing.T) {
	ce := NewEngine()
	conf := ce.ComputeAssumptionConfidence(nil)
	if conf != 0.0 {
		t.Errorf("Expected 0.0, got %f", conf)
	}
}

func TestComputeFreshness(t *testing.T) {
	ev := []models.Evidence{
		models.NewEvidence("test.csv", models.SourceTypeCSV, nil),
	}
	score := computeFreshness(ev)
	if score <= 0 || score > 1.0 {
		t.Errorf("Freshness out of range: %f", score)
	}
}

func TestComputeFreshnessEmpty(t *testing.T) {
	score := computeFreshness(nil)
	if score != 0.3 {
		t.Errorf("Expected 0.3 for empty, got %f", score)
	}
}

func TestComputeCoverage(t *testing.T) {
	v := models.NewVerification("asm_1", []string{"evd_1", "evd_2"}, models.VerificationResultVERIFIED, 0.85, "", nil)
	ev := []models.Evidence{
		models.NewEvidence("a.csv", models.SourceTypeCSV, nil),
		models.NewEvidence("b.csv", models.SourceTypeCSV, nil),
		models.NewEvidence("c.csv", models.SourceTypeCSV, nil),
	}
	score := computeCoverage(v, ev)
	if score <= 0 || score > 1.0 {
		t.Errorf("Coverage out of range: %f", score)
	}
}

func TestComputeCompleteness(t *testing.T) {
	v := models.NewVerification("asm_1", nil, models.VerificationResultCONTRADICTED, 0.92, "Found users outside Finance with access", nil)
	score := computeCompleteness(v)
	if score <= 0.3 {
		t.Errorf("Expected completeness > 0.3, got %f", score)
	}
}

func TestComputeCompletenessEmpty(t *testing.T) {
	v := models.NewVerification("asm_1", nil, models.VerificationResultUNKNOWN, 0.0, "", nil)
	score := computeCompleteness(v)
	if score != 0.1 {
		t.Errorf("Expected 0.1 for empty, got %f", score)
	}
}
