package extraction

import (
	"testing"
)

func TestExtractAccessClaim(t *testing.T) {
	ce := NewClaimExtractor()
	text := "Only Finance employees may access the payroll processing system."
	claims := ce.Extract(text, "test.txt", "/path/test.txt")
	if len(claims) != 1 {
		t.Fatalf("Expected 1 claim, got %d", len(claims))
	}
	if claims[0].Text != text {
		t.Errorf("Wrong text: %s", claims[0].Text)
	}
	if claims[0].ExtractionConfidence < 0.5 {
		t.Errorf("Confidence too low: %f", claims[0].ExtractionConfidence)
	}
}

func TestExtractEncryptionClaim(t *testing.T) {
	ce := NewClaimExtractor()
	text := "All payroll data is encrypted at rest and in transit."
	claims := ce.Extract(text, "test.txt", "")
	if len(claims) != 1 {
		t.Fatalf("Expected 1 claim, got %d", len(claims))
	}
	tags := claims[0].Tags
	hasConfig := false
	for _, tag := range tags {
		if tag == "configuration" {
			hasConfig = true
		}
	}
	if !hasConfig {
		t.Errorf("Expected 'configuration' tag, got %v", tags)
	}
}

func TestExtractMultiClaims(t *testing.T) {
	ce := NewClaimExtractor()
	text := "Only admins can access. All data is encrypted. All users must have MFA."
	claims := ce.Extract(text, "test.txt", "")
	if len(claims) != 3 {
		t.Fatalf("Expected 3 claims, got %d", len(claims))
	}
}

func TestNoClaimsForNarrative(t *testing.T) {
	ce := NewClaimExtractor()
	text := "The system was deployed in 2023. It runs on AWS. The team uses Kubernetes."
	claims := ce.Extract(text, "test.txt", "")
	if len(claims) != 0 {
		t.Errorf("Expected 0 claims for narrative text, got %d", len(claims))
	}
}

func TestDeduplication(t *testing.T) {
	ce := NewClaimExtractor()
	text := "Only admins can access. Only admins can access."
	claims := ce.Extract(text, "test.txt", "")
	if len(claims) != 1 {
		t.Errorf("Expected 1 claim (deduplicated), got %d", len(claims))
	}
}

func TestTagsExtracted(t *testing.T) {
	ce := NewClaimExtractor()
	text := "Only Finance employees may access the payroll processing system. All passwords must use MFA."
	claims := ce.Extract(text, "test.txt", "")
	if len(claims) != 2 {
		t.Fatalf("Expected 2 claims, got %d", len(claims))
	}
	// First claim should have access tag
	hasAccess := false
	for _, tag := range claims[0].Tags {
		if tag == "access" {
			hasAccess = true
		}
	}
	if !hasAccess {
		t.Errorf("First claim should have 'access' tag, got %v", claims[0].Tags)
	}
}

func TestDeclarativePatterns(t *testing.T) {
	ce := NewClaimExtractor()
	tests := []struct {
		text    string
		matched bool
	}{
		{"Only Finance employees may access the system.", true},
		{"All data must be encrypted.", true},
		{"Access is restricted to admins.", true},
		{"The system requires MFA for all users.", true},
		{"Nobody can access without approval.", false},
		{"Access requires MFA authentication.", true},
		{"Backups are encrypted.", true},
		{"The weather is nice today.", false},
		{"System version: 2.4.1", false},
	}
	for _, tc := range tests {
		claims := ce.Extract(tc.text, "test.txt", "")
		gotMatch := len(claims) > 0
		if gotMatch != tc.matched {
			t.Errorf("Text %q: expected matched=%v, got %v", tc.text, tc.matched, gotMatch)
		}
	}
}

func TestFullPolicyExtraction(t *testing.T) {
	ce := NewClaimExtractor()
	policy := `ACCESS CONTROL
Only Finance employees may access the payroll processing system.

ENCRYPTION
All payroll data is encrypted at rest and in transit.

ACCESS APPROVAL
Only the VP of Finance can approve payroll runs.

MFA
All users must have MFA enabled for production access.

BACKUP
System backups are performed daily and stored offsite.

NETWORK
The payroll system is not accessible from the internet.

MONITORING
All access is logged and monitored.

AUDIT
Quarterly security audits are conducted by an external firm.

DEPENDENCY
The payroll system depends on the corporate LDAP for authentication.

REVIEW
Access rights are reviewed on an annual basis.`
	claims := ce.Extract(policy, "test.txt", "")
	if len(claims) < 8 {
		t.Errorf("Expected at least 8 claims from full policy, got %d", len(claims))
	}
}
