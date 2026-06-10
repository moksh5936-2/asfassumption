package gaps

import (
	"testing"

	"asf-tui/asf/models"
)

func TestGapForContradicted(t *testing.T) {
	ge := NewEngine()
	a := models.NewAssumption("clm_1", "System assumes access control: Only admins can access.", models.AssumptionTypeACCESS, nil)
	v := models.NewVerification(a.ID, []string{"evd_1"}, models.VerificationResultCONTRADICTED, 0.92, "Found unauthorized access", nil)
	gaps := ge.GenerateGaps([]models.Assumption{a}, []models.Verification{v})
	if len(gaps) != 1 {
		t.Fatalf("Expected 1 gap, got %d", len(gaps))
	}
	if gaps[0].Severity != models.GapSeverityCRITICAL {
		t.Errorf("Expected CRITICAL for ACCESS contradicted with confidence 0.92, got %s", gaps[0].Severity)
	}
	if gaps[0].Type != models.GapTypeACCESS {
		t.Errorf("Expected ACCESS_GAP, got %s", gaps[0].Type)
	}
}

func TestGapForUnverified(t *testing.T) {
	ge := NewEngine()
	a := models.NewAssumption("clm_1", "System assumes access control: Only admins.", models.AssumptionTypeACCESS, nil)
	v := models.NewVerification(a.ID, nil, models.VerificationResultUNKNOWN, 0.0, "No evidence", nil)
	gaps := ge.GenerateGaps([]models.Assumption{a}, []models.Verification{v})
	if len(gaps) != 1 {
		t.Fatalf("Expected 1 gap, got %d", len(gaps))
	}
	if gaps[0].Severity != models.GapSeverityLOW {
		t.Errorf("Expected LOW for UNKNOWN, got %s", gaps[0].Severity)
	}
	if gaps[0].Type != models.GapTypeEVIDENCE {
		t.Errorf("Expected EVIDENCE_GAP, got %s", gaps[0].Type)
	}
}

func TestNoGapForVerified(t *testing.T) {
	ge := NewEngine()
	a := models.NewAssumption("clm_1", "System assumes access control: Only admins.", models.AssumptionTypeACCESS, nil)
	v := models.NewVerification(a.ID, []string{"evd_1"}, models.VerificationResultVERIFIED, 0.85, "All checks passed", nil)
	gaps := ge.GenerateGaps([]models.Assumption{a}, []models.Verification{v})
	if len(gaps) != 0 {
		t.Errorf("Expected 0 gaps for VERIFIED, got %d", len(gaps))
	}
}

func TestGapForPartiallyVerified(t *testing.T) {
	ge := NewEngine()
	a := models.NewAssumption("clm_1", "System assumes config state: Backups enabled.", models.AssumptionTypeCONFIGURATION, nil)
	v := models.NewVerification(a.ID, []string{"evd_1"}, models.VerificationResultPARTIALLY_VERIFIED, 0.5, "Partially compliant", nil)
	gaps := ge.GenerateGaps([]models.Assumption{a}, []models.Verification{v})
	if len(gaps) != 1 {
		t.Fatalf("Expected 1 gap, got %d", len(gaps))
	}
	if gaps[0].Severity != models.GapSeverityMEDIUM {
		t.Errorf("Expected MEDIUM for PARTIALLY_VERIFIED, got %s", gaps[0].Severity)
	}
}

func TestGapForUnknown(t *testing.T) {
	ge := NewEngine()
	a := models.NewAssumption("clm_1", "System assumes network: Isolated.", models.AssumptionTypeNETWORK, nil)
	v := models.NewVerification(a.ID, []string{"evd_1"}, models.VerificationResultUNKNOWN, 0.3, "Could not verify", nil)
	gaps := ge.GenerateGaps([]models.Assumption{a}, []models.Verification{v})
	if len(gaps) != 1 {
		t.Fatalf("Expected 1 gap, got %d", len(gaps))
	}
	if gaps[0].Severity != models.GapSeverityLOW {
		t.Errorf("Expected LOW for UNKNOWN, got %s", gaps[0].Severity)
	}
	if gaps[0].Type != models.GapTypeEVIDENCE {
		t.Errorf("Expected EVIDENCE_GAP, got %s", gaps[0].Type)
	}
}

func TestGapSeverityForLowConfidenceContradicted(t *testing.T) {
	ge := NewEngine()
	a := models.NewAssumption("clm_1", "System assumes: Something.", models.AssumptionTypeACCESS, nil)
	// Low confidence contradiction
	v := models.NewVerification(a.ID, []string{"evd_1"}, models.VerificationResultCONTRADICTED, 0.3, "Weak evidence", nil)
	gaps := ge.GenerateGaps([]models.Assumption{a}, []models.Verification{v})
	if len(gaps) != 1 {
		t.Fatalf("Expected 1 gap, got %d", len(gaps))
	}
	if gaps[0].Severity != models.GapSeverityMEDIUM {
		t.Errorf("Expected MEDIUM for low confidence contradicted, got %s", gaps[0].Severity)
	}
}

func TestNoVerificationMissing(t *testing.T) {
	ge := NewEngine()
	a := models.NewAssumption("clm_1", "System assumes: Something.", models.AssumptionTypeACCESS, nil)
	gaps := ge.GenerateGaps([]models.Assumption{a}, nil)
	if len(gaps) != 1 {
		t.Fatalf("Expected 1 gap, got %d", len(gaps))
	}
	if gaps[0].Type != models.GapTypeVERIFICATION {
		t.Errorf("Expected VERIFICATION_GAP, got %s", gaps[0].Type)
	}
}

func TestDetermineSeverity(t *testing.T) {
	tests := []struct {
		atype    models.AssumptionType
		conf     float64
		expected models.GapSeverity
	}{
		{models.AssumptionTypeACCESS, 0.92, models.GapSeverityCRITICAL},
		{models.AssumptionTypeIDENTITY, 0.85, models.GapSeverityCRITICAL},
		{models.AssumptionTypeNETWORK, 0.81, models.GapSeverityCRITICAL},
		{models.AssumptionTypeCONFIGURATION, 0.9, models.GapSeverityHIGH},
		{models.AssumptionTypeGOVERNANCE, 0.88, models.GapSeverityHIGH},
		{models.AssumptionTypePROCESS, 0.8, models.GapSeverityHIGH},
		{models.AssumptionTypeACCESS, 0.6, models.GapSeverityHIGH},
		{models.AssumptionTypeACCESS, 0.4, models.GapSeverityMEDIUM},
		{models.AssumptionTypeDOCUMENTATION, 0.9, models.GapSeverityHIGH},
	}
	for _, tc := range tests {
		v := models.NewVerification("asm_1", nil, models.VerificationResultCONTRADICTED, tc.conf, "", nil)
		sev := determineSeverity(tc.atype, v)
		if sev != tc.expected {
			t.Errorf("determineSeverity(%s, %.2f) = %s, want %s", tc.atype, tc.conf, sev, tc.expected)
		}
	}
}
