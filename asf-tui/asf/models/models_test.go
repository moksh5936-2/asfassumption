package models

import (
	"testing"
	"time"
)

func TestNewClaim(t *testing.T) {
	c := NewClaim("test.txt", "/path/to/test.txt", "All users must use MFA", 0.7, []string{"identity", "access"})
	if c.ID == "" {
		t.Error("Expected non-empty ID")
	}
	if c.SourceDocument != "test.txt" {
		t.Errorf("Expected source_document=test.txt, got %s", c.SourceDocument)
	}
	if c.Text != "All users must use MFA" {
		t.Errorf("Wrong text: %s", c.Text)
	}
	if c.ExtractionConfidence != 0.7 {
		t.Errorf("Wrong confidence: %f", c.ExtractionConfidence)
	}
	if len(c.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(c.Tags))
	}
}

func TestNewAssumption(t *testing.T) {
	a := NewAssumption("clm_abc", "System assumes access control: Only admins can access", AssumptionTypeACCESS, []string{"access", "admin"})
	if a.ClaimID != "clm_abc" {
		t.Error("Wrong claim_id")
	}
	if a.AssumptionType != AssumptionTypeACCESS {
		t.Error("Wrong assumption type")
	}
	if a.VerificationStatus != VerificationStatusPENDING {
		t.Error("Expected PENDING status")
	}
}

func TestNewEvidence(t *testing.T) {
	records := []map[string]interface{}{
		{"user": "alice", "role": "admin"},
	}
	e := NewEvidence("users.csv", SourceTypeCSV, records)
	if e.SourceType != SourceTypeCSV {
		t.Error("Wrong source type")
	}
	if len(e.Records) != 1 {
		t.Error("Expected 1 record")
	}
}

func TestNewVerification(t *testing.T) {
	v := NewVerification("asm_abc", []string{"evd_1", "evd_2"}, VerificationResultVERIFIED, 0.85, "All checks passed", nil)
	if v.AssumptionID != "asm_abc" {
		t.Error("Wrong assumption_id")
	}
	if v.Result != VerificationResultVERIFIED {
		t.Error("Wrong result")
	}
}

func TestNewGap(t *testing.T) {
	g := NewGap("asm_abc", GapSeverityCRITICAL, GapTypeACCESS, "Critical gap found", "Evidence shows contradiction")
	if g.Severity != GapSeverityCRITICAL {
		t.Error("Wrong severity")
	}
	if g.Type != GapTypeACCESS {
		t.Error("Wrong type")
	}
}

func TestAnalysisResult(t *testing.T) {
	r := AnalysisResult{
		Claims: []Claim{
			{ID: "clm_1", Text: "Claim 1"},
		},
		Assumptions: []Assumption{
			{ID: "asm_1", Text: "Assumption 1"},
		},
		Verifications: []Verification{
			{ID: "vrf_1", Result: VerificationResultVERIFIED},
			{ID: "vrf_2", Result: VerificationResultCONTRADICTED},
			{ID: "vrf_3", Result: VerificationResultUNKNOWN},
		},
		Gaps: []Gap{
			{ID: "gap_1", Severity: GapSeverityCRITICAL},
			{ID: "gap_2", Severity: GapSeverityCRITICAL},
			{ID: "gap_3", Severity: GapSeverityHIGH},
		},
	}

	if r.ClaimsFound() != 1 {
		t.Errorf("Expected 1 claim, got %d", r.ClaimsFound())
	}
	if r.AssumptionsFound() != 1 {
		t.Errorf("Expected 1 assumption, got %d", r.AssumptionsFound())
	}
	if r.VerifiedCount() != 1 {
		t.Errorf("Expected 1 verified, got %d", r.VerifiedCount())
	}
	if r.ContradictedCount() != 1 {
		t.Errorf("Expected 1 contradicted, got %d", r.ContradictedCount())
	}
	if r.UnknownCount() != 1 {
		t.Errorf("Expected 1 unknown, got %d", r.UnknownCount())
	}
	if r.CriticalGaps() != 2 {
		t.Errorf("Expected 2 critical, got %d", r.CriticalGaps())
	}

	s := r.BuildSummary()
	if s.ClaimsFound != 1 || s.CriticalGaps != 2 {
		t.Error("BuildSummary wrong")
	}
}

func TestModelDefaults(t *testing.T) {
	c := NewClaim("test.txt", "", "text", 0.5, nil)
	if c.CreatedAt.IsZero() {
		t.Error("CreatedAt should be set")
	}
	if c.CreatedAt.Location() != time.UTC {
		t.Error("CreatedAt should be UTC")
	}
}

func TestVerificationStatusMarshal(t *testing.T) {
	vs := VerificationStatusVERIFIED
	b, err := vs.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != `"3"` {
		t.Errorf("Expected \"3\", got %s", string(b))
	}
}

func TestModelIDs(t *testing.T) {
	c := NewClaim("doc", "", "text", 0.5, nil)
	a := NewAssumption(c.ID, "text", AssumptionTypeACCESS, nil)
	e := NewEvidence("src", SourceTypeCSV, nil)
	v := NewVerification(a.ID, nil, VerificationResultUNKNOWN, 0, "", nil)
	g := NewGap(a.ID, GapSeverityLOW, GapTypeEVIDENCE, "desc", "detail")

	if !hasPrefix(c.ID, "clm_") {
		t.Errorf("Claim ID should start with clm_: %s", c.ID)
	}
	if !hasPrefix(a.ID, "asm_") {
		t.Errorf("Assumption ID should start with asm_: %s", a.ID)
	}
	if !hasPrefix(e.ID, "evd_") {
		t.Errorf("Evidence ID should start with evd_: %s", e.ID)
	}
	if !hasPrefix(v.ID, "vrf_") {
		t.Errorf("Verification ID should start with vrf_: %s", v.ID)
	}
	if !hasPrefix(g.ID, "gap_") {
		t.Errorf("Gap ID should start with gap_: %s", g.ID)
	}
}

func hasPrefix(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	return s[:len(prefix)] == prefix
}
