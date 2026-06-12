package intelligence

import (
	"strings"
	"testing"
)

func TestReasoningEngineInferDatabaseAssumptions(t *testing.T) {
	arch := &ArchDescription{
		Name: "healthcare",
		Components: []Component{
			{ID: "db1", Label: "PatientDatabase"},
			{ID: "app1", Label: "WebApp"},
		},
		RawText: "System stores PHI for patients and uses HIPAA controls.",
	}
	re := NewReasoningEngine(arch)
	assumptions := re.InferAssumptions()

	foundKeyMgmt := false
	foundAudit := false
	foundObjectAuth := false
	for _, a := range assumptions {
		if a.Category == "KeyManagement" && strings.Contains(a.Description, "PHI") {
			foundKeyMgmt = true
		}
		if a.Category == "Auditability" && strings.Contains(a.Description, "PHI") {
			foundAudit = true
		}
		if a.Category == "ObjectLevelAuthorization" && strings.Contains(a.Description, "PHI") {
			foundObjectAuth = true
		}
	}
	if !foundKeyMgmt {
		t.Error("expected key management inference for PHI database")
	}
	if !foundAudit {
		t.Error("expected audit logging inference for PHI database")
	}
	if !foundObjectAuth {
		t.Error("expected object-level authorization inference for PHI database")
	}
}

func TestReasoningEngineInferIdentityProvider(t *testing.T) {
	arch := &ArchDescription{
		Name: "identity",
		Components: []Component{
			{ID: "idp1", Label: "Auth0"},
			{ID: "app1", Label: "WebApp"},
		},
		RawText: "Auth0 handles authentication and SSO.",
	}
	re := NewReasoningEngine(arch)
	assumptions := re.InferAssumptions()

	foundMFA := false
	foundSession := false
	foundToken := false
	foundAdmin := false
	for _, a := range assumptions {
		if a.Category == "Authentication" && strings.Contains(a.Description, "MFA") {
			foundMFA = true
		}
		if a.Category == "SessionSecurity" && strings.Contains(a.Description, "session") {
			foundSession = true
		}
		if a.Category == "Authentication" && strings.Contains(a.Description, "token validation") {
			foundToken = true
		}
		if a.Category == "PrivilegeManagement" && strings.Contains(a.Description, "admin") {
			foundAdmin = true
		}
	}
	if !foundMFA {
		t.Error("expected MFA inference for Auth0/IdP")
	}
	if !foundSession {
		t.Error("expected session security inference for Auth0/IdP")
	}
	if !foundToken {
		t.Error("expected token validation inference for Auth0/IdP")
	}
	if !foundAdmin {
		t.Error("expected admin access inference for Auth0/IdP")
	}
}

func TestReasoningEngineInferAPIGateway(t *testing.T) {
	arch := &ArchDescription{
		Name: "gateway",
		Components: []Component{
			{ID: "gw1", Label: "API Gateway"},
			{ID: "svc1", Label: "BackendService"},
		},
		RawText: "API Gateway routes traffic to backend services.",
	}
	re := NewReasoningEngine(arch)
	assumptions := re.InferAssumptions()

	foundRateLimit := false
	foundAuth := false
	foundLogging := false
	for _, a := range assumptions {
		if a.Category == "APISecurity" && strings.Contains(a.Description, "rate limiting") {
			foundRateLimit = true
		}
		if a.Category == "APISecurity" && strings.Contains(a.Description, "authentication") {
			foundAuth = true
		}
		if a.Category == "Logging" && strings.Contains(a.Description, "logging") {
			foundLogging = true
		}
	}
	if !foundRateLimit {
		t.Error("expected rate limiting inference for API Gateway")
	}
	if !foundAuth {
		t.Error("expected auth validation inference for API Gateway")
	}
	if !foundLogging {
		t.Error("expected logging inference for API Gateway")
	}
}

func TestReasoningEngineInferKMS(t *testing.T) {
	arch := &ArchDescription{
		Name: "kms",
		Components: []Component{
			{ID: "kms1", Label: "AWS KMS"},
			{ID: "db1", Label: "Database"},
		},
		RawText: "KMS is used for key management.",
	}
	re := NewReasoningEngine(arch)
	assumptions := re.InferAssumptions()

	foundRotation := false
	foundAccess := false
	foundBackup := false
	for _, a := range assumptions {
		if a.Category == "KeyManagement" && strings.Contains(a.Description, "rotation") {
			foundRotation = true
		}
		if a.Category == "KeyManagement" && strings.Contains(a.Description, "access restriction") {
			foundAccess = true
		}
		if a.Category == "DisasterRecovery" && strings.Contains(a.Description, "backup") {
			foundBackup = true
		}
	}
	if !foundRotation {
		t.Error("expected key rotation inference for KMS")
	}
	if !foundAccess {
		t.Error("expected access restriction inference for KMS")
	}
	if !foundBackup {
		t.Error("expected backup inference for KMS")
	}
}

func TestReasoningEngineInferBackupService(t *testing.T) {
	arch := &ArchDescription{
		Name: "backup",
		Components: []Component{
			{ID: "bkp1", Label: "BackupService"},
			{ID: "db1", Label: "Database"},
		},
		RawText: "BackupService creates snapshots of the database.",
	}
	re := NewReasoningEngine(arch)
	assumptions := re.InferAssumptions()

	foundEncryption := false
	foundTesting := false
	foundGeo := false
	for _, a := range assumptions {
		if a.Category == "Backups" && strings.Contains(a.Description, "encryption") {
			foundEncryption = true
		}
		if a.Category == "Backups" && strings.Contains(a.Description, "restore testing") {
			foundTesting = true
		}
		if a.Category == "Backups" && strings.Contains(a.Description, "geographic distribution") {
			foundGeo = true
		}
	}
	if !foundEncryption {
		t.Error("expected backup encryption inference")
	}
	if !foundTesting {
		t.Error("expected restore testing inference")
	}
	if !foundGeo {
		t.Error("expected geographic distribution inference")
	}
}

func TestReasoningEngineInferThirdParty(t *testing.T) {
	arch := &ArchDescription{
		Name: "thirdparty",
		Components: []Component{
			{ID: "tp1", Label: "ThirdPartyAnalytics"},
			{ID: "db1", Label: "Database"},
		},
		RawText: "Third-party analytics has database access.",
	}
	re := NewReasoningEngine(arch)
	assumptions := re.InferAssumptions()

	foundVendorRisk := false
	foundControls := false
	foundMinimization := false
	for _, a := range assumptions {
		if a.Category == "ThirdPartyRisk" && strings.Contains(a.Description, "vendor risk") {
			foundVendorRisk = true
		}
		if a.Category == "VendorRisk" && strings.Contains(a.Description, "contract") {
			foundControls = true
		}
		if a.Category == "ThirdPartyRisk" && strings.Contains(a.Description, "data minimization") {
			foundMinimization = true
		}
	}
	if !foundVendorRisk {
		t.Error("expected vendor risk inference for third-party")
	}
	if !foundControls {
		t.Error("expected equivalent controls inference for third-party")
	}
	if !foundMinimization {
		t.Error("expected data minimization inference for third-party")
	}
}

func TestReasoningEngineInferAdminConsole(t *testing.T) {
	arch := &ArchDescription{
		Name: "admin",
		Components: []Component{
			{ID: "adm1", Label: "AdminConsole"},
			{ID: "db1", Label: "Database"},
		},
		RawText: "AdminConsole provides management access.",
	}
	re := NewReasoningEngine(arch)
	assumptions := re.InferAssumptions()

	foundMFA := false
	foundBreakGlass := false
	foundAudit := false
	for _, a := range assumptions {
		if a.Category == "Authentication" && strings.Contains(a.Description, "MFA") {
			foundMFA = true
		}
		if a.Category == "PrivilegeManagement" && strings.Contains(a.Description, "break-glass") {
			foundBreakGlass = true
		}
		if a.Category == "Auditability" && strings.Contains(a.Description, "audit") {
			foundAudit = true
		}
	}
	if !foundMFA {
		t.Error("expected MFA inference for admin console")
	}
	if !foundBreakGlass {
		t.Error("expected break-glass inference for admin console")
	}
	if !foundAudit {
		t.Error("expected audit inference for admin console")
	}
}

func TestReasoningEngineInferAuditLog(t *testing.T) {
	arch := &ArchDescription{
		Name: "audit",
		Components: []Component{
			{ID: "log1", Label: "AuditLog"},
			{ID: "db1", Label: "Database"},
		},
		RawText: "AuditLog records all database access.",
	}
	re := NewReasoningEngine(arch)
	assumptions := re.InferAssumptions()

	foundImmutable := false
	foundRetention := false
	foundTamper := false
	for _, a := range assumptions {
		if a.Category == "Auditability" && strings.Contains(a.Description, "immutability") {
			foundImmutable = true
		}
		if a.Category == "DataRetention" && strings.Contains(a.Description, "retention") {
			foundRetention = true
		}
		if a.Category == "Auditability" && strings.Contains(a.Description, "tamper detection") {
			foundTamper = true
		}
	}
	if !foundImmutable {
		t.Error("expected immutability inference for audit log")
	}
	if !foundRetention {
		t.Error("expected retention inference for audit log")
	}
	if !foundTamper {
		t.Error("expected tamper detection inference for audit log")
	}
}

func TestReasoningEngineInferFromRawText(t *testing.T) {
	arch := &ArchDescription{
		Name: "rawtext",
		Components: []Component{
			{ID: "app1", Label: "App"},
		},
		RawText: "All data is encrypted. Communication uses HTTP and TLS.",
	}
	re := NewReasoningEngine(arch)
	assumptions := re.inferFromRawText()

	foundKeyMgmt := false
	foundHTTP := false
	for _, a := range assumptions {
		if a.Category == "KeyManagement" && strings.Contains(a.Description, "key management") {
			foundKeyMgmt = true
		}
		if a.Category == "NetworkSegmentation" && strings.Contains(a.Description, "TLS") {
			foundHTTP = true
		}
	}
	if !foundKeyMgmt {
		t.Error("expected key management inference from raw text mentioning encryption")
	}
	if !foundHTTP {
		t.Error("expected HTTP/TLS contradiction inference from raw text")
	}
}

func TestReasoningEngineDeduplication(t *testing.T) {
	arch := &ArchDescription{
		Name: "dedup",
		Components: []Component{
			{ID: "db1", Label: "Database"},
		},
		RawText: "Database with PHI.",
	}
	re := NewReasoningEngine(arch)
	assumptions := re.InferAllAssumptions()
	seen := make(map[string]bool)
	for _, a := range assumptions {
		key := a.ID + "|" + strings.ToLower(a.Description)
		if seen[key] {
			t.Errorf("duplicate assumption: %s", a.ID)
		}
		seen[key] = true
	}
}

func TestReasoningEngineNilArch(t *testing.T) {
	re := NewReasoningEngine(nil)
	assumptions := re.InferAssumptions()
	if len(assumptions) != 0 {
		t.Errorf("expected 0 assumptions for nil arch, got %d", len(assumptions))
	}
}

func TestBuildExplainability(t *testing.T) {
	why := buildExplainability("Template text", "Context here", []string{"control1", "control2"}, 0.85, "KeyManagement")
	if !strings.Contains(why, "Template text") {
		t.Error("expected template text in explainability")
	}
	if !strings.Contains(why, "Context here") {
		t.Error("expected context in explainability")
	}
	if !strings.Contains(why, "control1") {
		t.Error("expected control1 in explainability")
	}
	if !strings.Contains(why, "85%") {
		t.Error("expected 85% confidence in explainability")
	}
	if !strings.Contains(why, "KeyManagement") {
		t.Error("expected category in explainability")
	}
}
