package verification

import (
	"testing"

	"asf-tui/asf/models"
)

func makeEvidence(id, source string, st models.SourceType, records []map[string]interface{}) models.Evidence {
	return models.NewEvidence(source, st, records)
}

func makeAssumption(atype models.AssumptionType, text string) models.Assumption {
	return models.NewAssumption("clm_1", text, atype, nil)
}

func TestAccessContradicted(t *testing.T) {
	ve := NewEngine()
	a := makeAssumption(models.AssumptionTypeACCESS, "System assumes access control: Only Finance employees may access the payroll processing system.")
	ev := makeEvidence("evd_1", "payroll_acl.csv", models.SourceTypeCSV, []map[string]interface{}{
		{"user": "alice", "group": "Finance", "permission": "read"},
		{"user": "bob", "group": "Finance", "permission": "write"},
		{"user": "charlie", "group": "Engineering", "permission": "read"},
	})
	v := ve.Verify(a, []models.Evidence{ev})
	if v.Result != models.VerificationResultCONTRADICTED {
		t.Errorf("Expected CONTRADICTED, got %s", v.Result)
	}
	if v.Confidence <= 0 {
		t.Errorf("Expected positive confidence, got %f", v.Confidence)
	}
	if len(v.EvidenceUsed) != 1 {
		t.Errorf("Expected 1 evidence used, got %d", len(v.EvidenceUsed))
	}
}

func TestAccessVerified(t *testing.T) {
	ve := NewEngine()
	a := makeAssumption(models.AssumptionTypeACCESS, "System assumes access control: Only Finance employees may access the system.")
	ev := makeEvidence("evd_1", "acl.csv", models.SourceTypeCSV, []map[string]interface{}{
		{"user": "alice", "group": "Finance", "permission": "read"},
		{"user": "bob", "group": "Finance", "permission": "write"},
	})
	v := ve.Verify(a, []models.Evidence{ev})
	if v.Result != models.VerificationResultVERIFIED {
		t.Errorf("Expected VERIFIED, got %s", v.Result)
	}
}

func TestNetworkContradicted(t *testing.T) {
	ve := NewEngine()
	a := makeAssumption(models.AssumptionTypeNETWORK, "System assumes network posture: The payroll system is isolated from the internet.")
	ev := makeEvidence("evd_1", "network.csv", models.SourceTypeCSV, []map[string]interface{}{
		{"asset": "payroll-db", "public": "false"},
		{"asset": "web-gateway", "public": "true"},
	})
	v := ve.Verify(a, []models.Evidence{ev})
	if v.Result != models.VerificationResultCONTRADICTED {
		t.Errorf("Expected CONTRADICTED, got %s", v.Result)
	}
}

func TestNetworkExposed(t *testing.T) {
	ve := NewEngine()
	a := makeAssumption(models.AssumptionTypeNETWORK, "System assumes network posture: No public access to payroll systems.")
	ev := makeEvidence("evd_1", "network.csv", models.SourceTypeCSV, []map[string]interface{}{
		{"asset": "payroll-db", "exposed": "false"},
		{"asset": "web-gateway", "exposed": "false"},
	})
	v := ve.Verify(a, []models.Evidence{ev})
	if v.Result != models.VerificationResultVERIFIED {
		t.Errorf("Expected VERIFIED, got %s", v.Result)
	}
}

func TestMfaContradicted(t *testing.T) {
	ve := NewEngine()
	a := makeAssumption(models.AssumptionTypeIDENTITY, "All users must have MFA enabled.")
	ev := makeEvidence("evd_1", "mfa.csv", models.SourceTypeCSV, []map[string]interface{}{
		{"user": "alice", "mfa": "false"},
		{"user": "bob", "mfa": "false"},
	})
	v := ve.Verify(a, []models.Evidence{ev})
	if v.Result != models.VerificationResultCONTRADICTED {
		t.Errorf("Expected CONTRADICTED, got %s", v.Result)
	}
}

func TestConfigurationVerified(t *testing.T) {
	ve := NewEngine()
	a := makeAssumption(models.AssumptionTypeCONFIGURATION, "System assumes configuration state: Backups are enabled.")
	ev := makeEvidence("evd_1", "config.csv", models.SourceTypeCSV, []map[string]interface{}{
		{"resource": "payroll-db", "enabled": "true"},
		{"resource": "web-app", "enabled": "true"},
	})
	v := ve.Verify(a, []models.Evidence{ev})
	if v.Result != models.VerificationResultVERIFIED {
		t.Errorf("Expected VERIFIED, got %s", v.Result)
	}
}

func TestNoEvidenceUnknown(t *testing.T) {
	ve := NewEngine()
	a := makeAssumption(models.AssumptionTypeACCESS, "System assumes access control: Only admins can access.")
	v := ve.Verify(a, []models.Evidence{})
	if v.Result != models.VerificationResultUNKNOWN {
		t.Errorf("Expected UNKNOWN, got %s", v.Result)
	}
}

func TestMultipleEvidenceMerge(t *testing.T) {
	ve := NewEngine()
	a := makeAssumption(models.AssumptionTypeACCESS, "System assumes access control: Only Finance employees may access.")
	ev1 := makeEvidence("evd_1", "acl1.csv", models.SourceTypeCSV, []map[string]interface{}{
		{"user": "charlie", "group": "Engineering", "permission": "read"},
	})
	ev2 := makeEvidence("evd_2", "acl2.csv", models.SourceTypeCSV, []map[string]interface{}{
		{"user": "alice", "group": "Finance", "permission": "read"},
	})
	v := ve.Verify(a, []models.Evidence{ev1, ev2})
	if v.Result != models.VerificationResultCONTRADICTED {
		t.Errorf("Expected CONTRADICTED, got %s", v.Result)
	}
	if len(v.EvidenceUsed) != 2 {
		t.Errorf("Expected 2 evidence used, got %d", len(v.EvidenceUsed))
	}
}

func TestGovernanceContradicted(t *testing.T) {
	ve := NewEngine()
	a := makeAssumption(models.AssumptionTypeGOVERNANCE, "System assumes governance compliance: All security audits are completed.")
	ev := makeEvidence("evd_1", "audit.csv", models.SourceTypeCSV, []map[string]interface{}{
		{"status": "pending"},
		{"status": "pending"},
	})
	v := ve.Verify(a, []models.Evidence{ev})
	if v.Result != models.VerificationResultCONTRADICTED {
		t.Errorf("Expected CONTRADICTED, got %s", v.Result)
	}
}
