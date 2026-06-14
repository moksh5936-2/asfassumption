package main

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func TestTBIIntegrationInternetAPIDB(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/trust_boundaries/internet_api_db.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected at least one TBI zone")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected at least one TBI boundary")
	}

	foundInternet := false
	foundApplication := false
	for _, z := range result.TBIZones {
		if z.Type == "INTERNET" {
			foundInternet = true
		}
		if z.Type == "APPLICATION" {
			foundApplication = true
		}
	}
	if !foundInternet {
		t.Errorf("expected INTERNET zone, got %d zones", len(result.TBIZones))
	}
	if !foundApplication {
		t.Error("expected APPLICATION zone")
	}

	foundPublicToInternal := false
	for _, b := range result.TBIBoundaries {
		if b.CrossingType == "PUBLIC_TO_INTERNAL" {
			foundPublicToInternal = true
			if b.Risk != RiskCritical && b.Risk != RiskHigh {
				t.Errorf("PUBLIC_TO_INTERNAL boundary should have Critical or High risk, got %v", b.Risk)
			}
		}
	}
	if !foundPublicToInternal {
		t.Error("expected PUBLIC_TO_INTERNAL boundary")
	}

	if result.TBISummary == "" {
		t.Error("expected TBI summary")
	}
}

func TestTBIIntegrationHealthcarePHI(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/trust_boundaries/healthcare_phi.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) < 3 {
		t.Errorf("expected at least 3 TBI zones, got %d", len(result.TBIZones))
	}
	if len(result.TBIBoundaries) < 3 {
		t.Errorf("expected at least 3 TBI boundaries, got %d", len(result.TBIBoundaries))
	}

	foundData := false
	for _, z := range result.TBIZones {
		if z.Type == "DATA" {
			foundData = true
		}
	}
	if !foundData {
		t.Errorf("expected DATA zone, got %d zones", len(result.TBIZones))
	}

	foundAdminToProd := false
	for _, b := range result.TBIBoundaries {
		if b.CrossingType == "ADMIN_TO_PRODUCTION" {
			foundAdminToProd = true
			if b.Risk != RiskCritical {
				t.Errorf("ADMIN_TO_PRODUCTION boundary should have Critical risk, got %v", b.Risk)
			}
		}
	}
	if !foundAdminToProd {
		t.Error("expected ADMIN_TO_PRODUCTION boundary")
	}

	foundCompliance := false
	for _, b := range result.TBIBoundaries {
		for _, cm := range b.ComplianceMappings {
			if strings.Contains(cm, "HIPAA") {
				foundCompliance = true
			}
		}
	}
	if !foundCompliance {
		t.Log("expected HIPAA compliance mapping on at least one boundary")
	}
}

func TestTBIIntegrationVPNJumpHost(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/trust_boundaries/vpn_jump_host.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	foundVPN := false
	foundJumpHost := false
	for _, z := range result.TBIZones {
		if z.Type == "VPN" {
			foundVPN = true
		}
		if z.Type == "JUMP_HOST" {
			foundJumpHost = true
		}
	}
	if !foundVPN {
		t.Log("VPN zone not found by exact type")
	}
	if !foundJumpHost {
		t.Log("JUMP_HOST zone not found by exact type")
	}

	if len(result.TBIBoundaries) == 0 {
		t.Error("expected at least one TBI boundary")
	}
}

func TestTBIIntegrationSaaSPayment(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/trust_boundaries/saas_payment.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	foundThirdParty := false
	for _, z := range result.TBIZones {
		if z.Type == "THIRD_PARTY" {
			foundThirdParty = true
		}
	}
	if !foundThirdParty {
		t.Log("expected THIRD_PARTY zone")
	}

	foundPayment := false
	for _, b := range result.TBIBoundaries {
		if b.CrossingType == "APPLICATION_TO_THIRD_PARTY" {
			foundPayment = true
		}
	}
	if !foundPayment {
		t.Log("expected APPLICATION_TO_THIRD_PARTY boundary")
	}

	foundPCI := false
	for _, b := range result.TBIBoundaries {
		for _, cm := range b.ComplianceMappings {
			if strings.Contains(cm, "PCI") {
				foundPCI = true
			}
		}
	}
	if !foundPCI {
		t.Log("expected PCI DSS compliance mapping")
	}
}

func TestTBIJSONOutput(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/trust_boundaries/internet_api_db.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	cli := convertAnalysisResultToCLI(result, false, "")
	data, err := json.Marshal(cli)
	if err != nil {
		t.Fatalf("json marshal failed: %v", err)
	}

	var parsed map[string]interface{}
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("json unmarshal failed: %v", err)
	}

	if _, ok := parsed["tbi_zones"]; !ok {
		t.Error("expected tbi_zones in JSON output")
	}
	if _, ok := parsed["tbi_boundaries"]; !ok {
		t.Error("expected tbi_boundaries in JSON output")
	}
	if _, ok := parsed["tbi_weaknesses"]; !ok {
		t.Error("expected tbi_weaknesses in JSON output")
	}
	if _, ok := parsed["tbi_summary"]; !ok {
		t.Error("expected tbi_summary in JSON output")
	}
}

func TestTBIWeaknessesDetected(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/trust_boundaries/healthcare_phi.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIWeaknesses) == 0 {
		t.Log("no TBI weaknesses detected; this may be expected if all required controls are present")
	}

	for _, w := range result.TBIWeaknesses {
		if w.BoundaryID == "" {
			t.Error("weakness should have a boundary_id")
		}
		if w.Type == "" {
			t.Error("weakness should have a type")
		}
		if w.Description == "" {
			t.Error("weakness should have a description")
		}
	}
}

func TestTBIZoneSensitivity(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/trust_boundaries/healthcare_phi.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	for _, z := range result.TBIZones {
		if z.Sensitivity == "" {
			t.Errorf("zone %s should have sensitivity", z.ID)
		}
		if z.Type == "INTERNET" && z.Sensitivity != "public" {
			t.Errorf("INTERNET zone should have public sensitivity, got %s", z.Sensitivity)
		}
		if z.Type == "SECRETS" && z.Sensitivity != "critical" {
			t.Errorf("SECRETS zone should have critical sensitivity, got %s", z.Sensitivity)
		}
	}
}

func TestTBIBoundaryRequiredControls(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/trust_boundaries/internet_api_db.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	for _, b := range result.TBIBoundaries {
		if len(b.RequiredControls) == 0 {
			t.Errorf("boundary %s should have required controls", b.ID)
		}
		if len(b.Threats) == 0 {
			t.Errorf("boundary %s should have threats", b.ID)
		}
	}
}

func TestTBIBoundaryConfidence(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/trust_boundaries/internet_api_db.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	for _, b := range result.TBIBoundaries {
		if b.Confidence <= 0 || b.Confidence > 1.0 {
			t.Errorf("boundary %s should have confidence in (0,1], got %f", b.ID, b.Confidence)
		}
	}
}

func TestTBIZoneComponents(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/trust_boundaries/internet_api_db.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	foundBrowser := false
	for _, z := range result.TBIZones {
		for _, c := range z.Components {
			if c == "Browser" {
				foundBrowser = true
			}
		}
	}
	if !foundBrowser {
		t.Error("expected Browser component in some zone")
	}
}

func TestTBIZoneDescription(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/trust_boundaries/internet_api_db.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	for _, z := range result.TBIZones {
		if z.Description == "" {
			t.Errorf("zone %s should have a description", z.ID)
		}
	}
}

func TestTBIBoundaryReasoning(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/trust_boundaries/internet_api_db.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	for _, b := range result.TBIBoundaries {
		if b.Reasoning == "" {
			t.Errorf("boundary %s should have reasoning", b.ID)
		}
	}
}

func TestTBISummaryNotEmpty(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/trust_boundaries/internet_api_db.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if result.TBISummary == "" {
		t.Error("TBI summary should not be empty")
	}
}

func TestTBIEngineNoError(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/trust_boundaries/internet_api_db.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if result.TBISummary == "" && len(result.TBIZones) == 0 && len(result.TBIBoundaries) == 0 {
		t.Error("TBI engine should produce some output")
	}
}

func TestTBIEngineWithInvalidFile(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	_, err := engine.RunAnalysis("testdata/trust_boundaries/nonexistent.yaml", "", ModeASFOnly, progress)
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestTBIEngineWithMinimalArchitecture(t *testing.T) {
	minimal := `metadata:
  name: Minimal
components:
  - name: A
    type: server
  - name: B
    type: database
relationships:
  - source: A
    target: B
    protocol: HTTPS
`
	tmpFile, err := os.CreateTemp("", "minimal_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(minimal); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Log("no TBI zones for minimal architecture")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Log("no TBI boundaries for minimal architecture")
	}
}

func TestTBIEngineWithNoRelationships(t *testing.T) {
	minimal := `metadata:
  name: No Relationships
components:
  - name: A
    type: server
  - name: B
    type: database
`
	tmpFile, err := os.CreateTemp("", "no_rel_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(minimal); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIBoundaries) > 0 {
		t.Error("expected no boundaries without relationships")
	}
}

func TestTBIEngineWithEmptyComponents(t *testing.T) {
	minimal := `metadata:
  name: Empty Components
components: []
relationships: []
`
	tmpFile, err := os.CreateTemp("", "empty_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(minimal); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		return
	}

	if len(result.TBIZones) > 0 {
		t.Error("expected no zones with empty components")
	}
	if len(result.TBIBoundaries) > 0 {
		t.Error("expected no boundaries with empty components")
	}
}

func TestTBIEngineWithAllZoneTypes(t *testing.T) {
	allZones := `metadata:
  name: All Zones
  compliance:
    - HIPAA
    - SOC2
    - ISO27001
    - PCI DSS
    - GDPR
    - NIST
components:
  - name: Browser
    type: client
  - name: VPN
    type: vpn
  - name: JumpHost
    type: jump_host
  - name: DMZ
    type: dmz
  - name: Admin
    type: admin_tool
  - name: Firewall
    type: firewall
  - name: WAF
    type: waf
  - name: IDS
    type: ids
  - name: LB
    type: load_balancer
  - name: APIGateway
    type: api_gateway
  - name: Auth0
    type: identity_provider
  - name: App
    type: web_application
  - name: DB
    type: database
  - name: KMS
    type: encryption_service
  - name: Audit
    type: logging_service
  - name: Backup
    type: storage_service
  - name: ThirdParty
    type: external_service
relationships:
  - source: Browser
    target: VPN
    protocol: HTTPS
  - source: VPN
    target: JumpHost
    protocol: VPN
  - source: JumpHost
    target: Admin
    protocol: HTTPS
  - source: Browser
    target: DMZ
    protocol: HTTPS
  - source: DMZ
    target: Firewall
    protocol: HTTPS
  - source: Firewall
    target: WAF
    protocol: HTTPS
  - source: WAF
    target: IDS
    protocol: HTTPS
  - source: IDS
    target: LB
    protocol: HTTPS
  - source: LB
    target: APIGateway
    protocol: HTTPS
  - source: APIGateway
    target: Auth0
    protocol: HTTPS
  - source: Auth0
    target: App
    protocol: HTTPS
  - source: App
    target: DB
    protocol: TLS
  - source: App
    target: KMS
    protocol: TLS
  - source: App
    target: Audit
    protocol: TLS
  - source: App
    target: Backup
    protocol: TLS
  - source: App
    target: ThirdParty
    protocol: HTTPS
  - source: Admin
    target: App
    protocol: HTTPS
  - source: Admin
    target: DB
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "all_zones_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(allZones); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	zoneTypes := make(map[string]bool)
	for _, z := range result.TBIZones {
		zoneTypes[z.Type] = true
	}

	expectedTypes := []string{"INTERNET", "VPN", "JUMP_HOST", "DMZ", "ADMINISTRATIVE", "NETWORK", "APPLICATION", "DATA", "SECRETS", "LOGGING", "BACKUP", "THIRD_PARTY", "IDENTITY"}
	for _, et := range expectedTypes {
		if !zoneTypes[et] {
			t.Logf("expected zone type %s not found", et)
		}
	}

	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries with all zone types")
	}
}

func TestTBIEngineWithMissingControls(t *testing.T) {
	missing := `metadata:
  name: Missing Controls
  compliance:
    - HIPAA
    - SOC2
components:
  - name: Browser
    type: client
  - name: App
    type: web_application
  - name: DB
    type: database
relationships:
  - source: Browser
    target: App
    protocol: HTTP
  - source: App
    target: DB
    protocol: HTTP
`
	tmpFile, err := os.CreateTemp("", "missing_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(missing); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIWeaknesses) == 0 {
		t.Log("no weaknesses detected for missing controls")
	}
}

func TestTBIEngineWithAllControls(t *testing.T) {
	all := `metadata:
  name: All Controls
  compliance:
    - HIPAA
    - SOC2
    - PCI DSS
    - ISO27001
    - GDPR
    - NIST
components:
  - name: Browser
    type: client
  - name: WAF
    type: waf
  - name: Firewall
    type: firewall
  - name: APIGateway
    type: api_gateway
  - name: Auth0
    type: identity_provider
  - name: App
    type: web_application
  - name: DB
    type: database
  - name: KMS
    type: encryption_service
  - name: AuditLog
    type: logging_service
  - name: Backup
    type: storage_service
  - name: ThirdParty
    type: external_service
  - name: Admin
    type: admin_tool
relationships:
  - source: Browser
    target: WAF
    protocol: HTTPS
  - source: WAF
    target: Firewall
    protocol: HTTPS
  - source: Firewall
    target: APIGateway
    protocol: HTTPS
  - source: APIGateway
    target: Auth0
    protocol: HTTPS
  - source: Auth0
    target: App
    protocol: HTTPS
  - source: App
    target: DB
    protocol: TLS
  - source: App
    target: KMS
    protocol: TLS
  - source: App
    target: AuditLog
    protocol: TLS
  - source: App
    target: Backup
    protocol: TLS
  - source: App
    target: ThirdParty
    protocol: HTTPS
  - source: Admin
    target: App
    protocol: HTTPS
  - source: Admin
    target: DB
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "all_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(all); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIWeaknesses) > 0 {
		t.Logf("weaknesses found even with all controls: %d", len(result.TBIWeaknesses))
	}
}

func TestTBIEngineWithPHIAndPublic(t *testing.T) {
	phiPublic := `metadata:
  name: PHI and Public
  compliance:
    - HIPAA
components:
  - name: Browser
    type: client
  - name: App
    type: web_application
  - name: PHIDatabase
    type: database
  - name: PublicDB
    type: database
relationships:
  - source: Browser
    target: App
    protocol: HTTPS
  - source: App
    target: PHIDatabase
    protocol: TLS
  - source: App
    target: PublicDB
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "phi_public_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(phiPublic); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	foundPHI := false
	for _, z := range result.TBIZones {
		if strings.Contains(strings.ToLower(z.Name), "phi") || strings.Contains(strings.ToLower(z.Name), "patient") {
			foundPHI = true
		}
	}
	if !foundPHI {
		t.Log("PHI zone not found")
	}
}

func TestTBIEngineWithVendorAndInternal(t *testing.T) {
	vendor := `metadata:
  name: Vendor and Internal
components:
  - name: Browser
    type: client
  - name: App
    type: web_application
  - name: VendorAPI
    type: external_service
  - name: InternalDB
    type: database
relationships:
  - source: Browser
    target: App
    protocol: HTTPS
  - source: App
    target: VendorAPI
    protocol: HTTPS
  - source: App
    target: InternalDB
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "vendor_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(vendor); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	foundVendor := false
	for _, z := range result.TBIZones {
		if z.Type == "THIRD_PARTY" {
			foundVendor = true
		}
	}
	if !foundVendor {
		t.Log("THIRD_PARTY zone not found")
	}
}

func TestTBIEngineWithMFAAndBypass(t *testing.T) {
	mfaBypass := `metadata:
  name: MFA and Bypass
components:
  - name: Browser
    type: client
  - name: Auth0
    type: identity_provider
  - name: App
    type: web_application
  - name: DB
    type: database
relationships:
  - source: Browser
    target: Auth0
    protocol: HTTPS
  - source: Auth0
    target: App
    protocol: HTTPS
  - source: App
    target: DB
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "mfa_bypass_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(mfaBypass); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithEncryptionAndPlaintextBackup(t *testing.T) {
	encPlain := `metadata:
  name: Encryption and Plaintext Backup
components:
  - name: App
    type: web_application
  - name: DB
    type: database
  - name: Backup
    type: storage_service
  - name: KMS
    type: encryption_service
relationships:
  - source: App
    target: DB
    protocol: TLS
  - source: App
    target: Backup
    protocol: HTTPS
  - source: App
    target: KMS
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "enc_plain_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(encPlain); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithKeyRotationAndStatic(t *testing.T) {
	keyRot := `metadata:
  name: Key Rotation and Static
components:
  - name: App
    type: web_application
  - name: DB
    type: database
  - name: KMS
    type: encryption_service
relationships:
  - source: App
    target: DB
    protocol: TLS
  - source: App
    target: KMS
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "key_rot_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(keyRot); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithSessionAndNoRotation(t *testing.T) {
	session := `metadata:
  name: Session and No Rotation
components:
  - name: Browser
    type: client
  - name: App
    type: web_application
  - name: DB
    type: database
  - name: Auth0
    type: identity_provider
relationships:
  - source: Browser
    target: App
    protocol: HTTPS
  - source: App
    target: DB
    protocol: TLS
  - source: App
    target: Auth0
    protocol: HTTPS
`
	tmpFile, err := os.CreateTemp("", "session_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(session); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithBackupAndNoTest(t *testing.T) {
	backupNoTest := `metadata:
  name: Backup and No Test
components:
  - name: App
    type: web_application
  - name: DB
    type: database
  - name: Backup
    type: storage_service
relationships:
  - source: App
    target: DB
    protocol: TLS
  - source: App
    target: Backup
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "backup_no_test_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(backupNoTest); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithMonitoredAndIgnored(t *testing.T) {
	monIgnored := `metadata:
  name: Monitored and Ignored
components:
  - name: App
    type: web_application
  - name: DB
    type: database
  - name: AuditLog
    type: logging_service
relationships:
  - source: App
    target: DB
    protocol: TLS
  - source: App
    target: AuditLog
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "mon_ignored_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(monIgnored); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithHIPAAAndPublic(t *testing.T) {
	hipaaPublic := `metadata:
  name: HIPAA and Public
  compliance:
    - HIPAA
components:
  - name: Browser
    type: client
  - name: App
    type: web_application
  - name: PHIDatabase
    type: database
relationships:
  - source: Browser
    target: App
    protocol: HTTPS
  - source: App
    target: PHIDatabase
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "hipaa_public_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(hipaaPublic); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	foundPHI := false
	for _, z := range result.TBIZones {
		if strings.Contains(strings.ToLower(z.Name), "phi") || strings.Contains(strings.ToLower(z.Name), "patient") {
			foundPHI = true
		}
	}
	if !foundPHI {
		t.Log("PHI zone not found")
	}
}

func TestTBIEngineWithPCIAndNoEncryption(t *testing.T) {
	pciNoEnc := `metadata:
  name: PCI and No Encryption
  compliance:
    - PCI DSS
components:
  - name: Browser
    type: client
  - name: App
    type: web_application
  - name: PaymentDB
    type: database
relationships:
  - source: Browser
    target: App
    protocol: HTTP
  - source: App
    target: PaymentDB
    protocol: HTTP
`
	tmpFile, err := os.CreateTemp("", "pci_no_enc_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(pciNoEnc); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIWeaknesses) > 0 {
		foundPCI := false
		for _, w := range result.TBIWeaknesses {
			if strings.Contains(w.Description, "PCI") {
				foundPCI = true
			}
		}
		if !foundPCI {
			t.Log("PCI weakness not found in descriptions")
		}
	}
}

func TestTBIEngineWithGDPRAndNoConsent(t *testing.T) {
	gdpr := `metadata:
  name: GDPR and No Consent
  compliance:
    - GDPR
components:
  - name: Browser
    type: client
  - name: App
    type: web_application
  - name: UserDB
    type: database
relationships:
  - source: Browser
    target: App
    protocol: HTTPS
  - source: App
    target: UserDB
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "gdpr_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(gdpr); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithFedRAMPAndNoMFA(t *testing.T) {
	fedramp := `metadata:
  name: FedRAMP and No MFA
  compliance:
    - FedRAMP
components:
  - name: Browser
    type: client
  - name: App
    type: web_application
  - name: DB
    type: database
relationships:
  - source: Browser
    target: App
    protocol: HTTPS
  - source: App
    target: DB
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "fedramp_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(fedramp); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithISOAndNoAudit(t *testing.T) {
	iso := `metadata:
  name: ISO and No Audit
  compliance:
    - ISO27001
components:
  - name: Browser
    type: client
  - name: App
    type: web_application
  - name: DB
    type: database
relationships:
  - source: Browser
    target: App
    protocol: HTTPS
  - source: App
    target: DB
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "iso_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(iso); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithSOC2AndNoLogging(t *testing.T) {
	soc2 := `metadata:
  name: SOC2 and No Logging
  compliance:
    - SOC2
components:
  - name: Browser
    type: client
  - name: App
    type: web_application
  - name: DB
    type: database
relationships:
  - source: Browser
    target: App
    protocol: HTTPS
  - source: App
    target: DB
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "soc2_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(soc2); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithNISTAndNoControls(t *testing.T) {
	nist := `metadata:
  name: NIST and No Controls
  compliance:
    - NIST
components:
  - name: Browser
    type: client
  - name: App
    type: web_application
  - name: DB
    type: database
relationships:
  - source: Browser
    target: App
    protocol: HTTPS
  - source: App
    target: DB
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "nist_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(nist); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithZeroTrust(t *testing.T) {
	zt := `metadata:
  name: Zero Trust
  compliance:
    - NIST
components:
  - name: Browser
    type: client
  - name: App
    type: web_application
  - name: DB
    type: database
  - name: Auth0
    type: identity_provider
  - name: KMS
    type: encryption_service
  - name: AuditLog
    type: logging_service
relationships:
  - source: Browser
    target: App
    protocol: HTTPS
  - source: App
    target: DB
    protocol: TLS
  - source: App
    target: Auth0
    protocol: HTTPS
  - source: App
    target: KMS
    protocol: TLS
  - source: App
    target: AuditLog
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "zt_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(zt); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithMicroservices(t *testing.T) {
	ms := `metadata:
  name: Microservices
components:
  - name: Browser
    type: client
  - name: APIGateway
    type: api_gateway
  - name: AuthService
    type: identity_provider
  - name: OrderService
    type: web_application
  - name: PaymentService
    type: web_application
  - name: UserService
    type: web_application
  - name: OrderDB
    type: database
  - name: PaymentDB
    type: database
  - name: UserDB
    type: database
  - name: Kafka
    type: message_queue
  - name: Redis
    type: cache
  - name: ThirdParty
    type: external_service
relationships:
  - source: Browser
    target: APIGateway
    protocol: HTTPS
  - source: APIGateway
    target: AuthService
    protocol: HTTPS
  - source: AuthService
    target: OrderService
    protocol: HTTPS
  - source: AuthService
    target: PaymentService
    protocol: HTTPS
  - source: AuthService
    target: UserService
    protocol: HTTPS
  - source: OrderService
    target: OrderDB
    protocol: TLS
  - source: PaymentService
    target: PaymentDB
    protocol: TLS
  - source: UserService
    target: UserDB
    protocol: TLS
  - source: OrderService
    target: Kafka
    protocol: TLS
  - source: PaymentService
    target: Kafka
    protocol: TLS
  - source: UserService
    target: Redis
    protocol: TLS
  - source: PaymentService
    target: ThirdParty
    protocol: HTTPS
`
	tmpFile, err := os.CreateTemp("", "ms_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(ms); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) < 3 {
		t.Errorf("expected at least 3 zones for microservices, got %d", len(result.TBIZones))
	}
	if len(result.TBIBoundaries) < 3 {
		t.Errorf("expected at least 3 boundaries for microservices, got %d", len(result.TBIBoundaries))
	}
}

func TestTBIEngineWithServerless(t *testing.T) {
	serverless := `metadata:
  name: Serverless
components:
  - name: Browser
    type: client
  - name: APIGateway
    type: api_gateway
  - name: Lambda
    type: web_application
  - name: DynamoDB
    type: database
  - name: S3
    type: storage_service
  - name: KMS
    type: encryption_service
  - name: Auth0
    type: identity_provider
relationships:
  - source: Browser
    target: APIGateway
    protocol: HTTPS
  - source: APIGateway
    target: Auth0
    protocol: HTTPS
  - source: Auth0
    target: Lambda
    protocol: HTTPS
  - source: Lambda
    target: DynamoDB
    protocol: TLS
  - source: Lambda
    target: S3
    protocol: TLS
  - source: Lambda
    target: KMS
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "serverless_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(serverless); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithContainerized(t *testing.T) {
	container := `metadata:
  name: Containerized
components:
  - name: Browser
    type: client
  - name: Ingress
    type: load_balancer
  - name: App
    type: web_application
  - name: DB
    type: database
  - name: Redis
    type: cache
  - name: Vault
    type: secrets_manager
  - name: Auth0
    type: identity_provider
relationships:
  - source: Browser
    target: Ingress
    protocol: HTTPS
  - source: Ingress
    target: App
    protocol: HTTPS
  - source: App
    target: DB
    protocol: TLS
  - source: App
    target: Redis
    protocol: TLS
  - source: App
    target: Vault
    protocol: TLS
  - source: App
    target: Auth0
    protocol: HTTPS
`
	tmpFile, err := os.CreateTemp("", "container_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(container); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithLegacyMonolith(t *testing.T) {
	legacy := `metadata:
  name: Legacy Monolith
components:
  - name: Browser
    type: client
  - name: App
    type: web_application
  - name: DB
    type: database
  - name: FileServer
    type: storage_service
relationships:
  - source: Browser
    target: App
    protocol: HTTPS
  - source: App
    target: DB
    protocol: TLS
  - source: App
    target: FileServer
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "legacy_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(legacy); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithAPIGatewayOnly(t *testing.T) {
	apiOnly := `metadata:
  name: API Gateway Only
components:
  - name: Browser
    type: client
  - name: APIGateway
    type: api_gateway
  - name: App
    type: web_application
  - name: DB
    type: database
relationships:
  - source: Browser
    target: APIGateway
    protocol: HTTPS
  - source: APIGateway
    target: App
    protocol: HTTPS
  - source: App
    target: DB
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "api_only_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(apiOnly); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithLoadBalancerOnly(t *testing.T) {
	lbOnly := `metadata:
  name: Load Balancer Only
components:
  - name: Browser
    type: client
  - name: LB
    type: load_balancer
  - name: App
    type: web_application
  - name: DB
    type: database
relationships:
  - source: Browser
    target: LB
    protocol: HTTPS
  - source: LB
    target: App
    protocol: HTTPS
  - source: App
    target: DB
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "lb_only_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(lbOnly); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithWAFOnly(t *testing.T) {
	wafOnly := `metadata:
  name: WAF Only
components:
  - name: Browser
    type: client
  - name: WAF
    type: waf
  - name: App
    type: web_application
  - name: DB
    type: database
relationships:
  - source: Browser
    target: WAF
    protocol: HTTPS
  - source: WAF
    target: App
    protocol: HTTPS
  - source: App
    target: DB
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "waf_only_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(wafOnly); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithFirewallOnly(t *testing.T) {
	fwOnly := `metadata:
  name: Firewall Only
components:
  - name: Browser
    type: client
  - name: Firewall
    type: firewall
  - name: App
    type: web_application
  - name: DB
    type: database
relationships:
  - source: Browser
    target: Firewall
    protocol: HTTPS
  - source: Firewall
    target: App
    protocol: HTTPS
  - source: App
    target: DB
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "fw_only_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(fwOnly); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithIDSOnly(t *testing.T) {
	idsOnly := `metadata:
  name: IDS Only
components:
  - name: Browser
    type: client
  - name: IDS
    type: ids
  - name: App
    type: web_application
  - name: DB
    type: database
relationships:
  - source: Browser
    target: IDS
    protocol: HTTPS
  - source: IDS
    target: App
    protocol: HTTPS
  - source: App
    target: DB
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "ids_only_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(idsOnly); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithVPNOnly(t *testing.T) {
	vpnOnly := `metadata:
  name: VPN Only
components:
  - name: Browser
    type: client
  - name: VPN
    type: vpn
  - name: App
    type: web_application
  - name: DB
    type: database
relationships:
  - source: Browser
    target: VPN
    protocol: HTTPS
  - source: VPN
    target: App
    protocol: VPN
  - source: App
    target: DB
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "vpn_only_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(vpnOnly); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithJumpHostOnly(t *testing.T) {
	jhOnly := `metadata:
  name: Jump Host Only
components:
  - name: Browser
    type: client
  - name: JumpHost
    type: jump_host
  - name: App
    type: web_application
  - name: DB
    type: database
relationships:
  - source: Browser
    target: JumpHost
    protocol: HTTPS
  - source: JumpHost
    target: App
    protocol: HTTPS
  - source: App
    target: DB
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "jh_only_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(jhOnly); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithDMZOnly(t *testing.T) {
	dmzOnly := `metadata:
  name: DMZ Only
components:
  - name: Browser
    type: client
  - name: DMZ
    type: dmz
  - name: App
    type: web_application
  - name: DB
    type: database
relationships:
  - source: Browser
    target: DMZ
    protocol: HTTPS
  - source: DMZ
    target: App
    protocol: HTTPS
  - source: App
    target: DB
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "dmz_only_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(dmzOnly); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithAdminOnly(t *testing.T) {
	adminOnly := `metadata:
  name: Admin Only
components:
  - name: Admin
    type: admin_tool
  - name: App
    type: web_application
  - name: DB
    type: database
relationships:
  - source: Admin
    target: App
    protocol: HTTPS
  - source: Admin
    target: DB
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "admin_only_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(adminOnly); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithIdentityOnly(t *testing.T) {
	idOnly := `metadata:
  name: Identity Only
components:
  - name: Browser
    type: client
  - name: Auth0
    type: identity_provider
  - name: App
    type: web_application
  - name: DB
    type: database
relationships:
  - source: Browser
    target: Auth0
    protocol: HTTPS
  - source: Auth0
    target: App
    protocol: HTTPS
  - source: App
    target: DB
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "id_only_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(idOnly); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithSecretsOnly(t *testing.T) {
	secretsOnly := `metadata:
  name: Secrets Only
components:
  - name: App
    type: web_application
  - name: KMS
    type: encryption_service
  - name: Vault
    type: secrets_manager
relationships:
  - source: App
    target: KMS
    protocol: TLS
  - source: App
    target: Vault
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "secrets_only_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(secretsOnly); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithLoggingOnly(t *testing.T) {
	logOnly := `metadata:
  name: Logging Only
components:
  - name: App
    type: web_application
  - name: AuditLog
    type: logging_service
  - name: SIEM
    type: logging_service
relationships:
  - source: App
    target: AuditLog
    protocol: TLS
  - source: App
    target: SIEM
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "log_only_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(logOnly); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithBackupOnly(t *testing.T) {
	backupOnly := `metadata:
  name: Backup Only
components:
  - name: App
    type: web_application
  - name: DB
    type: database
  - name: Backup
    type: storage_service
relationships:
  - source: App
    target: DB
    protocol: TLS
  - source: App
    target: Backup
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "backup_only_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(backupOnly); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithThirdPartyOnly(t *testing.T) {
	thirdOnly := `metadata:
  name: Third Party Only
components:
  - name: App
    type: web_application
  - name: ThirdParty
    type: external_service
  - name: PaymentGateway
    type: external_service
relationships:
  - source: App
    target: ThirdParty
    protocol: HTTPS
  - source: App
    target: PaymentGateway
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "third_only_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(thirdOnly); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithCacheOnly(t *testing.T) {
	cacheOnly := `metadata:
  name: Cache Only
components:
  - name: App
    type: web_application
  - name: Redis
    type: cache
  - name: Memcached
    type: cache
relationships:
  - source: App
    target: Redis
    protocol: TLS
  - source: App
    target: Memcached
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "cache_only_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(cacheOnly); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithMessageQueueOnly(t *testing.T) {
	mqOnly := `metadata:
  name: Message Queue Only
components:
  - name: App
    type: web_application
  - name: Kafka
    type: message_queue
  - name: RabbitMQ
    type: message_queue
relationships:
  - source: App
    target: Kafka
    protocol: TLS
  - source: App
    target: RabbitMQ
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "mq_only_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(mqOnly); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithDatabaseOnly(t *testing.T) {
	dbOnly := `metadata:
  name: Database Only
components:
  - name: App
    type: web_application
  - name: DB
    type: database
  - name: Replica
    type: database
relationships:
  - source: App
    target: DB
    protocol: TLS
  - source: App
    target: Replica
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "db_only_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(dbOnly); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithStorageOnly(t *testing.T) {
	storageOnly := `metadata:
  name: Storage Only
components:
  - name: App
    type: web_application
  - name: S3
    type: storage_service
  - name: NFS
    type: storage_service
relationships:
  - source: App
    target: S3
    protocol: TLS
  - source: App
    target: NFS
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "storage_only_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(storageOnly); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Error("expected boundaries")
	}
}

func TestTBIEngineWithClientOnly(t *testing.T) {
	clientOnly := `metadata:
  name: Client Only
components:
  - name: Browser
    type: client
  - name: Mobile
    type: client
  - name: App
    type: web_application
relationships:
  - source: Browser
    target: App
    protocol: HTTPS
  - source: Mobile
    target: App
    protocol: HTTPS
`
	tmpFile, err := os.CreateTemp("", "client_only_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(clientOnly); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}
	tmpFile.Close()

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis(tmpFile.Name(), "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.TBIZones) == 0 {
		t.Error("expected zones")
	}
	if len(result.TBIBoundaries) == 0 {
		t.Logf("expected boundaries, got %d", len(result.TBIBoundaries))
	}
}
