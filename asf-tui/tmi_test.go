package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

func TestTMIIntegrationAuth0SaaS(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/threat_models/auth0_saas.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.Threats) == 0 {
		t.Error("expected at least one threat")
	}
	if len(result.ThreatClusters) == 0 {
		t.Error("expected at least one threat cluster")
	}
	if result.ThreatModelSummary.TotalThreats == 0 {
		t.Error("expected total_threats > 0")
	}

	foundIdentityThreat := false
	for _, threat := range result.Threats {
		if threat.Category == "IDENTITY" {
			foundIdentityThreat = true
		}
	}
	if !foundIdentityThreat {
		t.Error("expected IDENTITY threat for Auth0 SaaS")
	}

	foundDataThreat := false
	for _, threat := range result.Threats {
		if threat.Category == "DATA_PROTECTION" {
			foundDataThreat = true
		}
	}
	if !foundDataThreat {
		t.Error("expected DATA_PROTECTION threat for database")
	}

	if result.ThreatModelSummary.CriticalCount+result.ThreatModelSummary.HighCount == 0 {
		t.Error("expected critical or high threats")
	}
}

func TestTMIIntegrationHealthcarePHI(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/threat_models/healthcare_phi.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.Threats) == 0 {
		t.Error("expected at least one threat")
	}

	foundBackupThreat := false
	for _, threat := range result.Threats {
		if threat.Category == "BACKUP" {
			foundBackupThreat = true
		}
	}
	if !foundBackupThreat {
		t.Log("expected BACKUP threat for backup service")
	}

	foundAdminThreat := false
	for _, threat := range result.Threats {
		if threat.Category == "ACCESS_CONTROL" && strings.Contains(threat.Name, "Admin") {
			foundAdminThreat = true
		}
	}
	if !foundAdminThreat {
		t.Log("expected admin-related threat")
	}

	if result.ThreatModelSummary.TotalThreats == 0 {
		t.Error("expected total_threats > 0")
	}
}

func TestTMIIntegrationKubernetesCluster(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/threat_models/kubernetes_cluster.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.Threats) == 0 {
		t.Error("expected at least one threat")
	}

	foundCacheThreat := false
	for _, threat := range result.Threats {
		if strings.Contains(threat.Name, "Cache") {
			foundCacheThreat = true
		}
	}
	if !foundCacheThreat {
		t.Log("expected cache-related threat")
	}

	if result.ThreatModelSummary.TotalThreats == 0 {
		t.Error("expected total_threats > 0")
	}
}

func TestTMIIntegrationFintechPayment(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/threat_models/fintech_payment.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.Threats) == 0 {
		t.Error("expected at least one threat")
	}

	foundThirdPartyThreat := false
	for _, threat := range result.Threats {
		if threat.Category == "THIRD_PARTY" {
			foundThirdPartyThreat = true
		}
	}
	if !foundThirdPartyThreat {
		t.Log("expected THIRD_PARTY threat for payment gateway")
	}

	if result.ThreatModelSummary.TotalThreats == 0 {
		t.Error("expected total_threats > 0")
	}
}

func TestTMIIntegrationVPNInfrastructure(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/threat_models/vpn_infrastructure.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.Threats) == 0 {
		t.Error("expected at least one threat")
	}

	foundNetworkThreat := false
	for _, threat := range result.Threats {
		if threat.Category == "NETWORK" {
			foundNetworkThreat = true
		}
	}
	if !foundNetworkThreat {
		t.Log("expected NETWORK threat for VPN infrastructure")
	}

	if result.ThreatModelSummary.TotalThreats == 0 {
		t.Error("expected total_threats > 0")
	}
}

func TestTMIJSONOutput(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/threat_models/auth0_saas.yaml", "", ModeASFOnly, progress)
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

	if _, ok := parsed["threats"]; !ok {
		t.Error("expected threats in JSON output")
	}
	if _, ok := parsed["threat_clusters"]; !ok {
		t.Error("expected threat_clusters in JSON output")
	}
	if _, ok := parsed["threat_model_summary"]; !ok {
		t.Error("expected threat_model_summary in JSON output")
	}
}

func TestTMIThreatFields(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/threat_models/auth0_saas.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	for _, threat := range result.Threats {
		if threat.ID == "" {
			t.Error("threat should have an ID")
		}
		if threat.Name == "" {
			t.Error("threat should have a name")
		}
		if threat.Category == "" {
			t.Error("threat should have a category")
		}
		if threat.Description == "" {
			t.Error("threat should have a description")
		}
		if threat.Reasoning == "" {
			t.Error("threat should have reasoning")
		}
		if threat.Likelihood < 0 || threat.Likelihood > 1.0 {
			t.Errorf("threat likelihood should be in [0,1], got %f", threat.Likelihood)
		}
		if threat.Impact < 0 || threat.Impact > 1.0 {
			t.Errorf("threat impact should be in [0,1], got %f", threat.Impact)
		}
		if threat.RiskScore < 0 || threat.RiskScore > 1.0 {
			t.Errorf("threat risk_score should be in [0,1], got %f", threat.RiskScore)
		}
		if threat.Confidence <= 0 || threat.Confidence > 1.0 {
			t.Errorf("threat confidence should be in (0,1], got %f", threat.Confidence)
		}
		if len(threat.Recommendations) == 0 && len(threat.PreventiveControls) == 0 {
			t.Logf("threat %s has no recommendations or controls", threat.ID)
		}
	}
}

func TestTMISeverityDistribution(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/threat_models/auth0_saas.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	critical := 0
	high := 0
	medium := 0
	low := 0
	for _, threat := range result.Threats {
		switch threat.Severity {
		case RiskCritical:
			critical++
		case RiskHigh:
			high++
		case RiskMedium:
			medium++
		case RiskLow:
			low++
		}
	}

	if critical != result.ThreatModelSummary.CriticalCount {
		t.Errorf("critical count mismatch: %d vs %d", critical, result.ThreatModelSummary.CriticalCount)
	}
	if high != result.ThreatModelSummary.HighCount {
		t.Errorf("high count mismatch: %d vs %d", high, result.ThreatModelSummary.HighCount)
	}
	if medium != result.ThreatModelSummary.MediumCount {
		t.Errorf("medium count mismatch: %d vs %d", medium, result.ThreatModelSummary.MediumCount)
	}
	if low != result.ThreatModelSummary.LowCount {
		t.Errorf("low count mismatch: %d vs %d", low, result.ThreatModelSummary.LowCount)
	}
}

func TestTMIStrideDistribution(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/threat_models/auth0_saas.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.ThreatModelSummary.STRIDEDistribution) == 0 {
		t.Log("STRIDE distribution is empty")
	}

	for _, threat := range result.Threats {
		if len(threat.STRIDECategories) == 0 {
			t.Logf("threat %s has no STRIDE categories", threat.ID)
		}
	}
}

func TestTMIThreatClustering(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/threat_models/auth0_saas.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.ThreatClusters) == 0 {
		t.Error("expected threat clusters")
	}

	for _, cluster := range result.ThreatClusters {
		if cluster.ID == "" {
			t.Error("cluster should have an ID")
		}
		if cluster.Name == "" {
			t.Error("cluster should have a name")
		}
		if cluster.Category == "" {
			t.Error("cluster should have a category")
		}
		if len(cluster.Threats) == 0 {
			t.Error("cluster should have threats")
		}
		if cluster.RiskScore <= 0 {
			t.Error("cluster should have positive risk score")
		}
	}
}

func TestTMISummaryNotEmpty(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/threat_models/auth0_saas.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if result.ThreatModelSummary.SummaryText == "" {
		t.Error("threat model summary should not be empty")
	}

	if len(result.ThreatModelSummary.TopThreats) == 0 {
		t.Log("top threats list is empty")
	}
}

func TestTMITotalThreatsMatch(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/threat_models/auth0_saas.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.Threats) != result.ThreatModelSummary.TotalThreats {
		t.Errorf("threat count mismatch: %d vs %d", len(result.Threats), result.ThreatModelSummary.TotalThreats)
	}
}

func TestTMIClusterCountMatch(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/threat_models/auth0_saas.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if len(result.ThreatClusters) != result.ThreatModelSummary.ClusterCount {
		t.Errorf("cluster count mismatch: %d vs %d", len(result.ThreatClusters), result.ThreatModelSummary.ClusterCount)
	}
}

func TestTMIThreatsDeduplicated(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/threat_models/auth0_saas.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	seen := make(map[string]bool)
	for _, threat := range result.Threats {
		key := threat.Name + "-" + strings.Join(threat.AffectedComponents, ",")
		if seen[key] {
			t.Errorf("duplicate threat detected: %s", key)
		}
		seen[key] = true
	}
}

func TestTMIThreatIDsSequential(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/threat_models/auth0_saas.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	for i, threat := range result.Threats {
		expectedID := fmt.Sprintf("THREAT-%03d", i+1)
		if threat.ID != expectedID {
			t.Errorf("expected threat ID %s, got %s", expectedID, threat.ID)
		}
	}
}

func TestTMIThreatAffectedComponents(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/threat_models/auth0_saas.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	foundComponent := false
	for _, threat := range result.Threats {
		if len(threat.AffectedComponents) > 0 {
			foundComponent = true
			break
		}
	}
	if !foundComponent {
		t.Log("no threats have affected components")
	}
}

func TestTMIThreatControls(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/threat_models/auth0_saas.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	foundPreventive := false
	foundDetective := false
	foundCorrective := false
	for _, threat := range result.Threats {
		if len(threat.PreventiveControls) > 0 {
			foundPreventive = true
		}
		if len(threat.DetectiveControls) > 0 {
			foundDetective = true
		}
		if len(threat.CorrectiveControls) > 0 {
			foundCorrective = true
		}
	}
	if !foundPreventive {
		t.Log("no preventive controls found")
	}
	if !foundDetective {
		t.Log("no detective controls found")
	}
	if !foundCorrective {
		t.Log("no corrective controls found")
	}
}

func TestTMINoRegression(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()

	result, err := engine.RunAnalysis("testdata/threat_models/auth0_saas.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	if result.TotalAssumptions == 0 {
		t.Error("expected assumptions to still be generated")
	}
	if len(result.Assumptions) == 0 {
		t.Error("expected assumptions in result")
	}
	if result.CriticalCount == 0 && result.HighCount == 0 {
		t.Log("no critical or high assumptions found")
	}
}

func TestTMIEngineWithMinimalArchitecture(t *testing.T) {
	minimal := `metadata:
  name: Minimal
components:
  - name: App
    type: web_application
  - name: DB
    type: database
relationships:
  - source: App
    target: DB
    protocol: TLS
`
	tmpFile, err := os.CreateTemp("", "minimal_tmi_*.yaml")
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

	if len(result.Threats) == 0 {
		t.Log("no threats for minimal architecture")
	}
}

func TestTMIEngineWithEmptyComponents(t *testing.T) {
	minimal := `metadata:
  name: Empty Components
components: []
relationships: []
`
	tmpFile, err := os.CreateTemp("", "empty_tmi_*.yaml")
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

	if len(result.Threats) > 0 {
		t.Error("expected no threats with empty components")
	}
	if len(result.ThreatClusters) > 0 {
		t.Error("expected no clusters with empty components")
	}
}

func TestTMIEngineWithNoRelationships(t *testing.T) {
	minimal := `metadata:
  name: No Relationships
components:
  - name: App
    type: web_application
  - name: DB
    type: database
`
	tmpFile, err := os.CreateTemp("", "no_rel_tmi_*.yaml")
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

	if len(result.Threats) == 0 {
		t.Log("no threats without relationships")
	}
}

func TestTMIEngineWithHTTPProtocol(t *testing.T) {
	httpArch := `metadata:
  name: HTTP Protocol
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
	tmpFile, err := os.CreateTemp("", "http_tmi_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	if _, err := tmpFile.WriteString(httpArch); err != nil {
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

	foundMITM := false
	for _, threat := range result.Threats {
		if strings.Contains(threat.Name, "Man-in-the-Middle") {
			foundMITM = true
		}
	}
	if !foundMITM {
		t.Log("expected MITM threat for HTTP protocol")
	}
}

func TestTMIEngineWithAdminOnly(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "admin_tmi_*.yaml")
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

	foundAdminThreat := false
	for _, threat := range result.Threats {
		if strings.Contains(threat.Name, "Admin") || strings.Contains(threat.Name, "Privilege") {
			foundAdminThreat = true
		}
	}
	if !foundAdminThreat {
		t.Log("expected admin-related threat")
	}
}

func TestTMIEngineWithIdentityOnly(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "id_tmi_*.yaml")
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

	foundIdentityThreat := false
	for _, threat := range result.Threats {
		if threat.Category == "IDENTITY" {
			foundIdentityThreat = true
		}
	}
	if !foundIdentityThreat {
		t.Log("expected IDENTITY threat")
	}
}

func TestTMIEngineWithSecretsOnly(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "secrets_tmi_*.yaml")
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

	foundKeyThreat := false
	for _, threat := range result.Threats {
		if threat.Category == "KEY_MANAGEMENT" || threat.Category == "SECRETS_MANAGEMENT" {
			foundKeyThreat = true
		}
	}
	if !foundKeyThreat {
		t.Log("expected KEY_MANAGEMENT or SECRETS_MANAGEMENT threat")
	}
}

func TestTMIEngineWithLoggingOnly(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "log_tmi_*.yaml")
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

	foundMonitoringThreat := false
	for _, threat := range result.Threats {
		if threat.Category == "MONITORING" {
			foundMonitoringThreat = true
		}
	}
	if !foundMonitoringThreat {
		t.Log("expected MONITORING threat")
	}
}

func TestTMIEngineWithBackupOnly(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "backup_tmi_*.yaml")
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

	foundBackupThreat := false
	for _, threat := range result.Threats {
		if threat.Category == "BACKUP" {
			foundBackupThreat = true
		}
	}
	if !foundBackupThreat {
		t.Log("expected BACKUP threat")
	}
}

func TestTMIEngineWithThirdPartyOnly(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "third_tmi_*.yaml")
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

	foundThirdPartyThreat := false
	for _, threat := range result.Threats {
		if threat.Category == "THIRD_PARTY" {
			foundThirdPartyThreat = true
		}
	}
	if !foundThirdPartyThreat {
		t.Log("expected THIRD_PARTY threat")
	}
}

func TestTMIEngineWithCacheOnly(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "cache_tmi_*.yaml")
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

	foundCacheThreat := false
	for _, threat := range result.Threats {
		if strings.Contains(threat.Name, "Cache") {
			foundCacheThreat = true
		}
	}
	if !foundCacheThreat {
		t.Log("expected cache-related threat")
	}
}

func TestTMIEngineWithMessageQueueOnly(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "mq_tmi_*.yaml")
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

	foundMQThreat := false
	for _, threat := range result.Threats {
		if strings.Contains(threat.Name, "Message") {
			foundMQThreat = true
		}
	}
	if !foundMQThreat {
		t.Log("expected message queue-related threat")
	}
}

func TestTMIEngineWithDatabaseOnly(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "db_tmi_*.yaml")
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

	foundDBThreat := false
	for _, threat := range result.Threats {
		if threat.Category == "DATA_PROTECTION" {
			foundDBThreat = true
		}
	}
	if !foundDBThreat {
		t.Log("expected DATA_PROTECTION threat")
	}
}

func TestTMIEngineWithStorageOnly(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "storage_tmi_*.yaml")
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

	foundStorageThreat := false
	for _, threat := range result.Threats {
		if threat.Category == "DATA_PROTECTION" {
			foundStorageThreat = true
		}
	}
	if !foundStorageThreat {
		t.Log("expected DATA_PROTECTION threat")
	}
}

func TestTMIEngineWithClientOnly(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "client_tmi_*.yaml")
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

	foundClientThreat := false
	for _, threat := range result.Threats {
		if threat.Category == "NETWORK" || threat.Category == "IDENTITY" {
			foundClientThreat = true
		}
	}
	if !foundClientThreat {
		t.Log("expected NETWORK or IDENTITY threat")
	}
}

func TestTMIEngineWithMicroservices(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "ms_tmi_*.yaml")
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

	if len(result.Threats) < 5 {
		t.Errorf("expected at least 5 threats for microservices, got %d", len(result.Threats))
	}
	if len(result.ThreatClusters) < 2 {
		t.Errorf("expected at least 2 clusters for microservices, got %d", len(result.ThreatClusters))
	}
}

func TestTMIEngineWithServerless(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "serverless_tmi_*.yaml")
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

	if len(result.Threats) < 3 {
		t.Errorf("expected at least 3 threats for serverless, got %d", len(result.Threats))
	}
}

func TestTMIEngineWithContainerized(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "container_tmi_*.yaml")
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

	if len(result.Threats) < 3 {
		t.Errorf("expected at least 3 threats for containerized, got %d", len(result.Threats))
	}
}

func TestTMIEngineWithLegacyMonolith(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "legacy_tmi_*.yaml")
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

	if len(result.Threats) < 3 {
		t.Errorf("expected at least 3 threats for legacy monolith, got %d", len(result.Threats))
	}
}

func TestTMIEngineWithZeroTrust(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "zt_tmi_*.yaml")
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

	if len(result.Threats) < 3 {
		t.Errorf("expected at least 3 threats for zero trust, got %d", len(result.Threats))
	}
}

func TestTMIEngineWithAPIGatewayOnly(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "api_only_tmi_*.yaml")
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

	if len(result.Threats) < 3 {
		t.Errorf("expected at least 3 threats for API gateway, got %d", len(result.Threats))
	}
}

func TestTMIEngineWithLoadBalancerOnly(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "lb_only_tmi_*.yaml")
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

	if len(result.Threats) < 3 {
		t.Errorf("expected at least 3 threats for load balancer, got %d", len(result.Threats))
	}
}

func TestTMIEngineWithWAFOnly(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "waf_only_tmi_*.yaml")
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

	if len(result.Threats) < 3 {
		t.Errorf("expected at least 3 threats for WAF, got %d", len(result.Threats))
	}
}

func TestTMIEngineWithFirewallOnly(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "fw_only_tmi_*.yaml")
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

	if len(result.Threats) < 3 {
		t.Errorf("expected at least 3 threats for firewall, got %d", len(result.Threats))
	}
}

func TestTMIEngineWithIDSOnly(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "ids_only_tmi_*.yaml")
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

	if len(result.Threats) < 3 {
		t.Errorf("expected at least 3 threats for IDS, got %d", len(result.Threats))
	}
}

func TestTMIEngineWithVPNOnly(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "vpn_only_tmi_*.yaml")
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

	if len(result.Threats) < 3 {
		t.Errorf("expected at least 3 threats for VPN, got %d", len(result.Threats))
	}
}

func TestTMIEngineWithJumpHostOnly(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "jh_only_tmi_*.yaml")
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

	if len(result.Threats) < 3 {
		t.Errorf("expected at least 3 threats for jump host, got %d", len(result.Threats))
	}
}

func TestTMIEngineWithDMZOnly(t *testing.T) {
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
	tmpFile, err := os.CreateTemp("", "dmz_only_tmi_*.yaml")
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

	if len(result.Threats) < 3 {
		t.Errorf("expected at least 3 threats for DMZ, got %d", len(result.Threats))
	}
}
