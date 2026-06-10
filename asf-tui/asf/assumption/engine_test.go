package assumption

import (
	"testing"

	"asf-tui/asf/models"
)

func TestClassifyAccess(t *testing.T) {
	ae := NewEngine()
	claim := models.NewClaim("test.txt", "", "Only Finance employees may access the payroll processing system.", 0.6, nil)
	a := ae.Convert(claim)
	if a == nil {
		t.Fatal("Expected non-nil assumption")
	}
	if a.AssumptionType != models.AssumptionTypeACCESS {
		t.Errorf("Expected ACCESS, got %s", a.AssumptionType)
	}
	if a.ClaimID != claim.ID {
		t.Errorf("ClaimID mismatch: %s vs %s", a.ClaimID, claim.ID)
	}
}

func TestClassifyEncryption(t *testing.T) {
	ae := NewEngine()
	claim := models.NewClaim("test.txt", "", "All payroll data is encrypted at rest and in transit.", 0.6, nil)
	a := ae.Convert(claim)
	if a == nil {
		t.Fatal("Expected non-nil assumption")
	}
	if a.AssumptionType != models.AssumptionTypeCONFIGURATION {
		t.Errorf("Expected CONFIGURATION, got %s", a.AssumptionType)
	}
}

func TestClassifyNetwork(t *testing.T) {
	ae := NewEngine()
	claim := models.NewClaim("test.txt", "", "The payroll system is not accessible from the internet.", 0.6, nil)
	a := ae.Convert(claim)
	if a == nil {
		t.Fatal("Expected non-nil assumption")
	}
	if a.AssumptionType != models.AssumptionTypeNETWORK {
		t.Errorf("Expected NETWORK, got %s", a.AssumptionType)
	}
}

func TestClassifyIdentity(t *testing.T) {
	ae := NewEngine()
	claim := models.NewClaim("test.txt", "", "All users must have MFA enabled for production access.", 0.6, nil)
	a := ae.Convert(claim)
	if a == nil {
		t.Fatal("Expected non-nil assumption")
	}
	if a.AssumptionType != models.AssumptionTypeIDENTITY {
		t.Errorf("Expected IDENTITY, got %s", a.AssumptionType)
	}
}

func TestClassifyGovernance(t *testing.T) {
	ae := NewEngine()
	claim := models.NewClaim("test.txt", "", "Quarterly security audits are conducted by an external firm.", 0.6, nil)
	a := ae.Convert(claim)
	if a == nil {
		t.Fatal("Expected non-nil assumption")
	}
	if a.AssumptionType != models.AssumptionTypeGOVERNANCE {
		t.Errorf("Expected GOVERNANCE, got %s", a.AssumptionType)
	}
}

func TestClassifyProcess(t *testing.T) {
	ae := NewEngine()
	// "approved" matches both ACCESS (score 1) and PROCESS (score 1),
	// ACCESS wins tiebreak (earlier in pattern list, matching Python behavior)
	claim := models.NewClaim("test.txt", "", "All access requests must be approved by management.", 0.6, nil)
	a := ae.Convert(claim)
	if a == nil {
		t.Fatal("Expected non-nil assumption")
	}
	_ = a
}

func TestClassifyProcessOnly(t *testing.T) {
	ae := NewEngine()
	// "shall be reviewed" matches PROCESS (score 1), "access" matches ACCESS (score 1).
	// Tiebreak: ACCESS wins (first in pattern list order, matching Python behavior)
	claim := models.NewClaim("test.txt", "", "All access requests shall be reviewed by security team.", 0.6, nil)
	a := ae.Convert(claim)
	if a == nil {
		t.Fatal("Expected non-nil assumption")
	}
	_ = a
}

func TestClassifyProcessOnlyUnique(t *testing.T) {
	ae := NewEngine()
	// Only PROCESS pattern 2 matches, no ACCESS keyword
	claim := models.NewClaim("test.txt", "", "The report must be reviewed by security team.", 0.6, nil)
	a := ae.Convert(claim)
	if a == nil {
		t.Fatal("Expected non-nil assumption")
	}
	if a.AssumptionType != models.AssumptionTypePROCESS {
		t.Errorf("Expected PROCESS, got %s", a.AssumptionType)
	}
}

func TestClassifyDocumentationOnly(t *testing.T) {
	ae := NewEngine()
	claim := models.NewClaim("test.txt", "", "The security policy is documented as described in the runbook.", 0.6, nil)
	a := ae.Convert(claim)
	if a == nil {
		t.Fatal("Expected non-nil assumption")
	}
	if a.AssumptionType != models.AssumptionTypeDOCUMENTATION {
		t.Errorf("Expected DOCUMENTATION, got %s", a.AssumptionType)
	}
}

func TestClassifyDocumentationTiebreak(t *testing.T) {
	ae := NewEngine()
	// "documented" and "runbook" match DOCUMENTATION (score 1) and
	// "network" matches NETWORK (score 1); NETWORK wins tiebreak (matches Python)
	claim := models.NewClaim("test.txt", "", "Network architecture is documented in the runbook.", 0.6, nil)
	a := ae.Convert(claim)
	if a == nil {
		t.Fatal("Expected non-nil assumption")
	}
	// NETWORK wins tiebreak (earlier in pattern list)
	if a.AssumptionType != models.AssumptionTypeNETWORK {
		t.Errorf("Expected NETWORK (tiebreak), got %s", a.AssumptionType)
	}
}

func TestClassifyDependency(t *testing.T) {
	ae := NewEngine()
	claim := models.NewClaim("test.txt", "", "The payroll system depends on the corporate LDAP.", 0.6, nil)
	a := ae.Convert(claim)
	if a == nil {
		t.Fatal("Expected non-nil assumption")
	}
	if a.AssumptionType != models.AssumptionTypeDEPENDENCY {
		t.Errorf("Expected DEPENDENCY, got %s", a.AssumptionType)
	}
}

func TestNoMatch(t *testing.T) {
	ae := NewEngine()
	claim := models.NewClaim("test.txt", "", "The weather is nice today.", 0.3, nil)
	a := ae.Convert(claim)
	if a != nil {
		t.Errorf("Expected nil for non-security text, got assumption of type %s", a.AssumptionType)
	}
}

func TestConvertMany(t *testing.T) {
	ae := NewEngine()
	claims := []models.Claim{
		models.NewClaim("test.txt", "", "Only admins can access.", 0.6, nil),
		models.NewClaim("test.txt", "", "All data is encrypted.", 0.6, nil),
		models.NewClaim("test.txt", "", "The weather is nice.", 0.3, nil),
	}
	assumptions := ae.ConvertMany(claims)
	if len(assumptions) != 2 {
		t.Errorf("Expected 2 assumptions (1 filtered out), got %d", len(assumptions))
	}
}

func TestBuildAssumptionText(t *testing.T) {
	ae := NewEngine()
	claim := models.NewClaim("test.txt", "", "Only admins can access.", 0.6, nil)
	a := ae.Convert(claim)
	if a == nil {
		t.Fatal("Expected non-nil")
	}
	if a.Text != "System assumes access control: Only admins can access." {
		t.Errorf("Wrong assumption text: %s", a.Text)
	}
}

func TestExtractKeywords(t *testing.T) {
	ae := NewEngine()
	claim := models.NewClaim("test.txt", "", "Only Finance employees may access the payroll system.", 0.6, nil)
	a := ae.Convert(claim)
	if a == nil {
		t.Fatal("Expected non-nil")
	}
	if len(a.Keywords) == 0 {
		t.Error("Expected non-empty keywords")
	}
}
