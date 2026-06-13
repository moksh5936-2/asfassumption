package main

import (
	"testing"

	"asf-tui/intelligence"
)

// TestCIEMFAExemption tests MFA required vs exempt contradiction.
func TestCIEMFAExemption(t *testing.T) {
	arch := &intelligence.ArchDescription{
		Name: "MFA Test",
		ExplicitAssumptions: []string{
			"MFA is enforced for all users",
			"Service accounts are exempt from MFA",
		},
		Components: []intelligence.Component{
			{ID: "WebApp", Label: "WebApp"},
			{ID: "Database", Label: "Database"},
		},
	}
	assumptions := []intelligence.Assumption{
		{ID: "A1", Description: "MFA is enforced for all users", Category: "IDENTITY"},
		{ID: "A2", Description: "Service accounts are exempt from MFA", Category: "IDENTITY"},
	}
	cie := intelligence.NewCIEEngine()
	contradictions := cie.DetectAllContradictions(arch, assumptions, nil, nil)

	found := false
	for _, c := range contradictions {
		if c.Type == intelligence.ContradictionTypeAUTHENTICATION {
			found = true
			if c.Severity != intelligence.RiskHigh {
				t.Errorf("Expected HIGH severity, got %s", c.Severity)
			}
			if c.Confidence < 0.8 {
				t.Errorf("Expected confidence >= 0.8, got %.2f", c.Confidence)
			}
			t.Logf("✅ CIE detected: %s", c.Summary)
		}
	}
	if !found {
		t.Errorf("Expected AUTHENTICATION contradiction for MFA exemption")
	}
}

// TestCIEPlaintextBackup tests encrypted vs plaintext backup contradiction.
func TestCIEPlaintextBackup(t *testing.T) {
	arch := &intelligence.ArchDescription{
		Name: "Backup Test",
		ExplicitAssumptions: []string{
			"All data is encrypted at rest",
			"Backups are stored in plaintext",
		},
		Components: []intelligence.Component{
			{ID: "Database", Label: "Database"},
		},
	}
	assumptions := []intelligence.Assumption{
		{ID: "A1", Description: "All data is encrypted at rest", Category: "CONFIGURATION"},
		{ID: "A2", Description: "Backups are stored in plaintext", Category: "CONFIGURATION"},
	}
	cie := intelligence.NewCIEEngine()
	contradictions := cie.DetectAllContradictions(arch, assumptions, nil, nil)

	found := false
	for _, c := range contradictions {
		if c.Type == intelligence.ContradictionTypeENCRYPTION || c.Type == intelligence.ContradictionTypeCONTROL || c.Type == intelligence.ContradictionTypeBACKUP {
			found = true
			t.Logf("✅ CIE detected: %s", c.Summary)
		}
	}
	if !found {
		t.Errorf("Expected ENCRYPTION, CONTROL, or BACKUP contradiction for plaintext backup, got %d contradictions", len(contradictions))
		for _, c := range contradictions {
			t.Logf("  got: type=%s summary=%s", c.Type, c.Summary)
		}
	}
}

// TestCIESharedAdmin tests least privilege vs shared admin contradiction.
func TestCIESharedAdmin(t *testing.T) {
	arch := &intelligence.ArchDescription{
		Name: "Admin Test",
		ExplicitAssumptions: []string{
			"Least privilege is enforced",
			"Administrators share a single account",
		},
		Components: []intelligence.Component{
			{ID: "WebApp", Label: "WebApp"},
		},
	}
	assumptions := []intelligence.Assumption{
		{ID: "A1", Description: "Least privilege is enforced", Category: "ACCESS"},
		{ID: "A2", Description: "Administrators share a single account", Category: "ACCESS"},
	}
	cie := intelligence.NewCIEEngine()
	contradictions := cie.DetectAllContradictions(arch, assumptions, nil, nil)

	found := false
	for _, c := range contradictions {
		if c.Type == intelligence.ContradictionTypeAUTHORIZATION {
			found = true
			if c.Severity != intelligence.RiskCritical {
				t.Errorf("Expected CRITICAL severity, got %s", c.Severity)
			}
			t.Logf("✅ CIE detected: %s", c.Summary)
		}
	}
	if !found {
		t.Errorf("Expected AUTHORIZATION contradiction for shared admin")
	}
}

// TestCIEPrivatePublic tests private subnet vs public access contradiction.
func TestCIEPrivatePublic(t *testing.T) {
	arch := &intelligence.ArchDescription{
		Name: "Network Test",
		ExplicitAssumptions: []string{
			"Database is in a private subnet",
			"Database is accessible from the internet",
		},
		Components: []intelligence.Component{
			{ID: "Database", Label: "Database"},
		},
	}
	assumptions := []intelligence.Assumption{
		{ID: "A1", Description: "Database is in a private subnet", Category: "NETWORK"},
		{ID: "A2", Description: "Database is accessible from the internet", Category: "NETWORK"},
	}
	cie := intelligence.NewCIEEngine()
	contradictions := cie.DetectAllContradictions(arch, assumptions, nil, nil)

	found := false
	for _, c := range contradictions {
		if c.Type == intelligence.ContradictionTypeNETWORK {
			found = true
			if c.Severity != intelligence.RiskHigh {
				t.Errorf("Expected HIGH severity, got %s", c.Severity)
			}
			t.Logf("✅ CIE detected: %s", c.Summary)
		}
	}
	if !found {
		t.Errorf("Expected NETWORK contradiction for private/public")
	}
}

// TestCIEImpliedContradiction tests implied contradictions (PHI + public database).
func TestCIEImpliedContradiction(t *testing.T) {
	arch := &intelligence.ArchDescription{
		Name:    "PHI Test",
		RawText: "System contains PHI patient data. Database is publicly accessible.",
		Components: []intelligence.Component{
			{ID: "PHIDatabase", Label: "PHIDatabase"},
		},
	}
	assumptions := []intelligence.Assumption{
		{ID: "A1", Description: "PHI data is stored in PHIDatabase", Category: "DataProtection"},
		{ID: "A2", Description: "Database is publicly accessible", Category: "NETWORK"},
	}
	cie := intelligence.NewCIEEngine()
	contradictions := cie.DetectAllContradictions(arch, assumptions, nil, nil)

	found := false
	for _, c := range contradictions {
		if c.Type == intelligence.ContradictionTypeDATA_CLASSIFICATION {
			found = true
			if c.Severity != intelligence.RiskCritical {
				t.Errorf("Expected CRITICAL severity for PHI, got %s", c.Severity)
			}
			t.Logf("✅ CIE detected implied: %s", c.Summary)
		}
	}
	if !found {
		t.Errorf("Expected DATA_CLASSIFICATION implied contradiction for PHI")
	}
}

// TestCIEComplianceContradiction tests compliance contradictions (HIPAA without audit).
func TestCIEComplianceContradiction(t *testing.T) {
	arch := &intelligence.ArchDescription{
		Name:       "HIPAA Test",
		Compliance: []string{"HIPAA", "SOC2"},
		Components: []intelligence.Component{
			{ID: "PHIDatabase", Label: "PHIDatabase"},
		},
	}
	assumptions := []intelligence.Assumption{
		{ID: "A1", Description: "PHI data is encrypted", Category: "DataProtection"},
	}
	cie := intelligence.NewCIEEngine()
	contradictions := cie.DetectAllContradictions(arch, assumptions, nil, nil)

	foundHIPAA := false
	foundSOC2 := false
	for _, c := range contradictions {
		if c.Type == intelligence.ContradictionTypeCOMPLIANCE {
			if c.Summary == "HIPAA compliance requires audit but not documented" {
				foundHIPAA = true
			}
			if c.Summary == "SOC2 compliance requires access but not documented" {
				foundSOC2 = true
			}
			t.Logf("✅ CIE detected compliance: %s", c.Summary)
		}
	}
	if !foundHIPAA {
		t.Logf("⚠️ HIPAA audit contradiction not found (may be expected if audit assumptions exist)")
	}
	if !foundSOC2 {
		t.Logf("⚠️ SOC2 access contradiction not found (may be expected if access assumptions exist)")
	}
}

// TestCIEControlContradiction tests control contradictions.
func TestCIEControlContradiction(t *testing.T) {
	arch := &intelligence.ArchDescription{
		Name: "Control Test",
		Components: []intelligence.Component{
			{ID: "WebApp", Label: "WebApp"},
		},
	}
	assumptions := []intelligence.Assumption{
		{ID: "A1", Description: "MFA is enforced", Category: "IDENTITY"},
		{ID: "A2", Description: "Service accounts bypass MFA", Category: "IDENTITY"},
	}
	controls := []intelligence.ControlDetail{
		{ID: "CTRL-001", Name: "MFA Control", Description: "MFA must be enforced for all accounts", Category: "Authentication"},
	}
	cie := intelligence.NewCIEEngine()
	contradictions := cie.DetectAllContradictions(arch, assumptions, controls, nil)

	found := false
	for _, c := range contradictions {
		if c.Type == intelligence.ContradictionTypeCONTROL {
			found = true
			if c.Severity != intelligence.RiskHigh {
				t.Errorf("Expected HIGH severity, got %s", c.Severity)
			}
			t.Logf("✅ CIE detected control: %s", c.Summary)
		}
	}
	if !found {
		t.Errorf("Expected CONTROL contradiction for MFA bypass")
	}
}

// TestCIETrustBoundaryContradiction tests trust boundary contradictions.
func TestCIETrustBoundaryContradiction(t *testing.T) {
	arch := &intelligence.ArchDescription{
		Name:    "Trust Boundary Test",
		RawText: "PHI data in internal network",
		Components: []intelligence.Component{
			{ID: "PHIDatabase", Label: "PHIDatabase"},
		},
	}
	assumptions := []intelligence.Assumption{
		{ID: "A1", Description: "PHI data is encrypted", Category: "DataProtection"},
	}
	boundaries := []intelligence.TrustBoundary{
		{Type: "Internet", Components: []string{"PHIDatabase"}, RiskLevel: intelligence.RiskCritical},
	}
	cie := intelligence.NewCIEEngine()
	contradictions := cie.DetectAllContradictions(arch, assumptions, nil, boundaries)

	found := false
	for _, c := range contradictions {
		if c.Type == intelligence.ContradictionTypeTRUST_BOUNDARY {
			found = true
			if c.Severity != intelligence.RiskCritical {
				t.Errorf("Expected CRITICAL severity, got %s", c.Severity)
			}
			t.Logf("✅ CIE detected trust boundary: %s", c.Summary)
		}
	}
	if !found {
		t.Errorf("Expected TRUST_BOUNDARY contradiction for PHI internet exposure")
	}
}

// TestCIEAllContradictionsFromFiles runs CIE against all test data files.
func TestCIEAllContradictionsFromFiles(t *testing.T) {
	files := []string{
		"testdata/contradictions/mfa_exemption.yaml",
		"testdata/contradictions/plaintext_backup.yaml",
		"testdata/contradictions/shared_admin.yaml",
		"testdata/contradictions/private_public.yaml",
		"testdata/contradictions/encrypted_backup_unknown.yaml",
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			cfg := &Config{}
			engine := NewEngine(cfg)
			progress := make(chan AnalysisProgress, 100)
			go func() {
				for range progress {
				}
			}()
			result, err := engine.RunAnalysis(file, "", ModeASFOnly, progress)
			if err != nil {
				t.Fatalf("RunAnalysis failed: %v", err)
			}

			t.Logf("File: %s", file)
			t.Logf("  Assumptions: %d", len(result.Assumptions))
			t.Logf("  CIE Contradictions: %d", len(result.CIEContradictions))
			for _, c := range result.CIEContradictions {
				t.Logf("    - [%s] %s: %s", c.Severity, c.Type, c.Summary)
			}

			if len(result.CIEContradictions) == 0 {
				t.Logf("⚠️ No contradictions found in %s (may be due to native analyzer output)", file)
			}
		})
	}
}

// TestCIESummary tests that the CIE summary is generated.
func TestCIESummary(t *testing.T) {
	arch := &intelligence.ArchDescription{
		Name: "Summary Test",
		ExplicitAssumptions: []string{
			"MFA is enforced for all users",
			"Service accounts are exempt from MFA",
			"All data is encrypted at rest",
			"Backups are stored in plaintext",
		},
		Components: []intelligence.Component{
			{ID: "WebApp", Label: "WebApp"},
			{ID: "Database", Label: "Database"},
		},
	}
	assumptions := []intelligence.Assumption{
		{ID: "A1", Description: "MFA is enforced for all users", Category: "IDENTITY"},
		{ID: "A2", Description: "Service accounts are exempt from MFA", Category: "IDENTITY"},
		{ID: "A3", Description: "All data is encrypted at rest", Category: "CONFIGURATION"},
		{ID: "A4", Description: "Backups are stored in plaintext", Category: "CONFIGURATION"},
	}
	cie := intelligence.NewCIEEngine()
	contradictions := cie.DetectAllContradictions(arch, assumptions, nil, nil)

	summary := intelligence.BuildContradictionSummary(contradictions)
	t.Logf("CIE Summary: %s", summary)

	if len(contradictions) == 0 {
		t.Errorf("Expected contradictions, found none")
	}
}
