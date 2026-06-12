package intelligence

import (
	"strings"
	"testing"
)

func TestContradictionEngineMFAExemption(t *testing.T) {
	ce := NewContradictionEngine()
	assumptions := []Assumption{
		{ID: "A1", Description: "MFA is enforced for all users", Category: "Authentication"},
		{ID: "A2", Description: "Service accounts are exempt from MFA", Category: "Authentication"},
	}
	contradictions := ce.DetectContradictions(assumptions)
	found := false
	for _, c := range contradictions {
		if c.RuleName == "MFA_ENFORCED_WITH_EXEMPTION" {
			found = true
			if c.Severity != RiskCritical {
				t.Errorf("expected Critical severity, got %s", c.Severity)
			}
		}
	}
	if !found {
		t.Error("expected MFA_ENFORCED_WITH_EXEMPTION contradiction")
	}
}

func TestContradictionEnginePlaintextBackup(t *testing.T) {
	ce := NewContradictionEngine()
	assumptions := []Assumption{
		{ID: "A1", Description: "All data is encrypted at rest", Category: "DataProtection"},
		{ID: "A2", Description: "Backups are stored in plaintext", Category: "Backups"},
	}
	contradictions := ce.DetectContradictions(assumptions)
	found := false
	for _, c := range contradictions {
		if c.RuleName == "ENCRYPTED_WITH_PLAINTEXT_BACKUP" {
			found = true
		}
	}
	if !found {
		t.Error("expected ENCRYPTED_WITH_PLAINTEXT_BACKUP contradiction")
	}
}

func TestContradictionEngineSharedAdmin(t *testing.T) {
	ce := NewContradictionEngine()
	assumptions := []Assumption{
		{ID: "A1", Description: "Least privilege is enforced", Category: "Authorization"},
		{ID: "A2", Description: "Shared admin account is used for operations", Category: "PrivilegeManagement"},
	}
	contradictions := ce.DetectContradictions(assumptions)
	found := false
	for _, c := range contradictions {
		if c.RuleName == "LEAST_PRIVILEGE_WITH_SHARED_ADMIN" {
			found = true
		}
	}
	if !found {
		t.Error("expected LEAST_PRIVILEGE_WITH_SHARED_ADMIN contradiction")
	}
}

func TestContradictionEngineInternetAccessiblePrivate(t *testing.T) {
	ce := NewContradictionEngine()
	assumptions := []Assumption{
		{ID: "A1", Description: "Database is in private subnet", Category: "NetworkSegmentation"},
		{ID: "A2", Description: "Database is internet accessible", Category: "NetworkSegmentation"},
	}
	contradictions := ce.DetectContradictions(assumptions)
	found := false
	for _, c := range contradictions {
		if c.RuleName == "PRIVATE_SUBNET_INTERNET_ACCESSIBLE" {
			found = true
		}
	}
	if !found {
		t.Error("expected PRIVATE_SUBNET_INTERNET_ACCESSIBLE contradiction")
	}
}

func TestContradictionEngineMutableAudit(t *testing.T) {
	ce := NewContradictionEngine()
	assumptions := []Assumption{
		{ID: "A1", Description: "Audit logs are immutable and tamper-proof", Category: "Auditability"},
		{ID: "A2", Description: "Log deletion is allowed for storage management", Category: "Logging"},
	}
	contradictions := ce.DetectContradictions(assumptions)
	found := false
	for _, c := range contradictions {
		if c.RuleName == "IMMUTABLE_AUDIT_WITH_DELETION" {
			found = true
		}
	}
	if !found {
		t.Error("expected IMMUTABLE_AUDIT_WITH_DELETION contradiction")
	}
}

func TestContradictionEngineHTTPAllowed(t *testing.T) {
	ce := NewContradictionEngine()
	assumptions := []Assumption{
		{ID: "A1", Description: "TLS is required for all communications", Category: "NetworkSegmentation"},
		{ID: "A2", Description: "HTTP is allowed for legacy compatibility", Category: "NetworkSegmentation"},
	}
	contradictions := ce.DetectContradictions(assumptions)
	found := false
	for _, c := range contradictions {
		if c.RuleName == "TLS_REQUIRED_HTTP_ALLOWED" {
			found = true
		}
	}
	if !found {
		t.Error("expected TLS_REQUIRED_HTTP_ALLOWED contradiction")
	}
}

func TestContradictionEngineEncryptionWithoutKeyManagement(t *testing.T) {
	ce := NewContradictionEngine()
	assumptions := []Assumption{
		{ID: "A1", Description: "All sensitive data is encrypted", Category: "DataProtection"},
	}
	contradictions := ce.DetectContradictions(assumptions)
	found := false
	for _, c := range contradictions {
		if c.RuleName == "ENCRYPTION_WITHOUT_KEY_MANAGEMENT" {
			found = true
			if c.Severity != RiskHigh {
				t.Errorf("expected High severity, got %s", c.Severity)
			}
		}
	}
	if !found {
		t.Error("expected ENCRYPTION_WITHOUT_KEY_MANAGEMENT contradiction")
	}
}

func TestContradictionEngineSessionWithoutRotation(t *testing.T) {
	ce := NewContradictionEngine()
	assumptions := []Assumption{
		{ID: "A1", Description: "Session tokens are used for authentication", Category: "SessionSecurity"},
	}
	contradictions := ce.DetectContradictions(assumptions)
	found := false
	for _, c := range contradictions {
		if c.RuleName == "SESSION_WITHOUT_ROTATION" {
			found = true
			if c.Severity != RiskHigh {
				t.Errorf("expected High severity, got %s", c.Severity)
			}
		}
	}
	if !found {
		t.Error("expected SESSION_WITHOUT_ROTATION contradiction")
	}
}

func TestContradictionEngineNoContradictions(t *testing.T) {
	ce := NewContradictionEngine()
	assumptions := []Assumption{
		{ID: "A1", Description: "MFA is enforced for all users", Category: "Authentication"},
		{ID: "A2", Description: "Admin access requires break-glass approval", Category: "PrivilegeManagement"},
	}
	contradictions := ce.DetectContradictions(assumptions)
	if len(contradictions) != 0 {
		t.Errorf("expected 0 contradictions, got %d", len(contradictions))
	}
}

func TestContradictionEngineCountBySeverity(t *testing.T) {
	ce := NewContradictionEngine()
	assumptions := []Assumption{
		{ID: "A1", Description: "MFA is enforced", Category: "Authentication"},
		{ID: "A2", Description: "Service accounts are exempt", Category: "Authentication"},
		{ID: "A3", Description: "Backups are plaintext", Category: "Backups"},
		{ID: "A4", Description: "All data is encrypted", Category: "DataProtection"},
	}
	contradictions := ce.DetectContradictions(assumptions)
	counts := ce.CountBySeverity(contradictions)
	if counts[RiskCritical] < 1 {
		t.Errorf("expected at least 1 Critical contradiction, got %d", counts[RiskCritical])
	}
	if counts[RiskHigh] < 1 {
		t.Errorf("expected at least 1 High contradiction, got %d", counts[RiskHigh])
	}
}

func TestGetAffectedAssumptionIDs(t *testing.T) {
	contradictions := []Contradiction{
		{
			RuleName:            "TEST",
			AffectedAssumptions: []string{"A1", "A2"},
		},
		{
			RuleName:            "TEST2",
			AffectedAssumptions: []string{"A2", "A3"},
		},
	}
	ids := GetAffectedAssumptionIDs(contradictions)
	if len(ids) != 3 {
		t.Errorf("expected 3 unique affected IDs, got %d", len(ids))
	}
	seen := make(map[string]bool)
	for _, id := range ids {
		if seen[id] {
			t.Errorf("duplicate ID: %s", id)
		}
		seen[id] = true
	}
}

func TestFormatContradiction(t *testing.T) {
	c := Contradiction{
		Severity:            RiskCritical,
		Explanation:         "Test contradiction",
		RuleName:            "TEST_RULE",
		AffectedAssumptions: []string{"A1"},
	}
	formatted := FormatContradiction(c)
	if formatted == "" {
		t.Error("expected non-empty formatted contradiction")
	}
	if !strings.Contains(formatted, "TEST_RULE") {
		t.Error("expected rule name in formatted contradiction")
	}
}
