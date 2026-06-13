package verify

import (
	"fmt"
	"strings"
	"testing"
)

func TestEmptyVerification(t *testing.T) {
	engine := NewVerificationEngine("general", nil, nil)
	output := engine.RunAll()
	if output == nil {
		t.Fatal("expected non-nil output")
	}
}

func TestVerificationSingleMFA(t *testing.T) {
	inputs := []VerificationInput{
		{ID: "A1", Description: "MFA is enforced for all users", Category: "identity", Risk: "Critical", Keywords: []string{"mfa", "multi-factor"}},
	}
	engine := NewVerificationEngine("general", []string{"auth0"}, inputs)
	output := engine.RunAll()
	if output.Assessment == nil || len(output.Assessment.Plans) != 1 {
		t.Fatalf("expected 1 plan, got %d", len(output.Assessment.Plans))
	}
	p := output.Assessment.Plans[0]
	if p.Confidence <= 0 {
		t.Errorf("expected positive confidence, got %.0f", p.Confidence)
	}
	if len(p.EvidenceRequired) == 0 {
		t.Error("expected evidence requirements")
	}
	if len(p.Actions) == 0 {
		t.Error("expected verification actions")
	}
	if p.WhyVerify == "" {
		t.Error("expected why verify rationale")
	}
	if p.Priority != VpHigh {
		t.Errorf("expected High priority for verified MFA, got %s", p.Priority)
	}
	if p.Status != VsPartiallyVerified && p.Status != VsVerified {
		t.Errorf("expected PartiallyVerified or Verified for MFA, got %s", p.Status)
	}
}

func TestVerificationMultipleAssumptions(t *testing.T) {
	inputs := []VerificationInput{
		{ID: "A1", Description: "MFA is enforced for all users", Category: "identity", Risk: "Critical", Keywords: []string{"mfa"}},
		{ID: "A2", Description: "RBAC is configured for role-based access", Category: "authorization", Risk: "High", Keywords: []string{"rbac", "role"}},
		{ID: "A3", Description: "TLS is enabled for all API traffic", Category: "cryptography", Risk: "Critical", Keywords: []string{"tls"}},
	}
	engine := NewVerificationEngine("general", nil, inputs)
	output := engine.RunAll()
	if output.Assessment == nil || len(output.Assessment.Plans) != 3 {
		t.Fatalf("expected 3 plans, got %d", len(output.Assessment.Plans))
	}
	if output.Assessment.TotalAssumptions != 3 {
		t.Errorf("expected 3 total assumptions, got %d", output.Assessment.TotalAssumptions)
	}
	if output.Assessment.OverallConfidence <= 0 {
		t.Errorf("expected positive overall confidence, got %.1f", output.Assessment.OverallConfidence)
	}
}

func TestVerificationConfidenceLevels(t *testing.T) {
	tests := []struct {
		name    string
		risk    string
		minConf float64
	}{
		{"critical risk", "Critical", 30},
		{"high risk", "High", 30},
		{"low risk", "Low", 30},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inputs := []VerificationInput{
				{ID: "A1", Description: "Test assumption with " + tt.name, Category: "operational", Risk: tt.risk, Keywords: []string{"secret"}},
			}
			engine := NewVerificationEngine("general", nil, inputs)
			output := engine.RunAll()
			if output.Assessment == nil || len(output.Assessment.Plans) == 0 {
				t.Fatal("expected plan")
			}
			p := output.Assessment.Plans[0]
			if p.Confidence < tt.minConf {
				t.Errorf("expected confidence >= %.0f for %s, got %.0f", tt.minConf, tt.name, p.Confidence)
			}
		})
	}
}

func TestVerificationPriority(t *testing.T) {
	tests := []struct {
		name     string
		risk     string
		status   VerificationStatus
		category EvidenceCategory
		domain   string
		expected VerificationPriority
	}{
		{"critical no evidence", "Critical", VsNoEvidence, EvCatOperational, "general", VpCritical},
		{"high no evidence", "High", VsNoEvidence, EvCatOperational, "general", VpHigh},
		{"medium no evidence", "Medium", VsNoEvidence, EvCatOperational, "general", VpMedium},
		{"low no evidence", "Low", VsNoEvidence, EvCatOperational, "general", VpMedium},
		{"critical verified identity", "Critical", VsVerified, EvCatIdentity, "general", VpHigh},
		{"critical verified", "Critical", VsVerified, EvCatOperational, "general", VpMedium},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			priority := computePriority(tt.risk, tt.status, tt.category, tt.domain)
			if priority != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, priority)
			}
		})
	}
}

func TestVerificationDomainSpecific(t *testing.T) {
	inputs := []VerificationInput{
		{ID: "A1", Description: "PHI access is restricted to authorized personnel", Category: "identity", Risk: "Critical", Keywords: []string{"phi", "access"}},
	}
	engine := NewVerificationEngine("healthcare", nil, inputs)
	output := engine.RunAll()
	if output.Assessment == nil || len(output.Assessment.Plans) != 1 {
		t.Fatalf("expected 1 plan, got %d", len(output.Assessment.Plans))
	}
	p := output.Assessment.Plans[0]
	foundDomainEvidence := false
	for _, ev := range p.EvidenceRequired {
		if strings.Contains(ev.Name, "PHI") {
			foundDomainEvidence = true
			break
		}
	}
	if !foundDomainEvidence {
		t.Error("expected PHI-specific evidence for healthcare domain")
	}
	foundBreakGlass := false
	for _, ev := range p.EvidenceRequired {
		if strings.Contains(ev.Name, "Break Glass") {
			foundBreakGlass = true
			break
		}
	}
	if !foundBreakGlass {
		t.Error("expected Break Glass procedure evidence for healthcare domain")
	}
}

func TestVerificationRoadmap(t *testing.T) {
	inputs := []VerificationInput{
		{ID: "A1", Description: "Key rotation is configured for KMS", Category: "cryptography", Risk: "Critical", Keywords: []string{"kms", "key rotation"}},
		{ID: "A2", Description: "Backup restore testing is performed", Category: "resilience", Risk: "High", Keywords: []string{"backup", "restore"}},
	}
	engine := NewVerificationEngine("general", nil, inputs)
	output := engine.RunAll()
	if len(output.Roadmaps) == 0 {
		t.Fatal("expected roadmaps")
	}
	if output.Roadmaps[0].Priority != VpHigh {
		t.Errorf("expected roadmap[0] to be High, got %s", output.Roadmaps[0].Priority)
	}
	if len(output.Roadmaps[0].Steps) == 0 {
		t.Error("expected steps in roadmap")
	}
}

func TestCISOView(t *testing.T) {
	inputs := []VerificationInput{
		{ID: "A1", Description: "MFA is enforced", Category: "identity", Risk: "Critical", Keywords: []string{"mfa"}},
		{ID: "A2", Description: "SSO is configured", Category: "identity", Risk: "High", Keywords: []string{"sso"}},
		{ID: "A3", Description: "TLS is enabled", Category: "cryptography", Risk: "Critical", Keywords: []string{"tls"}},
	}
	engine := NewVerificationEngine("general", nil, inputs)
	output := engine.RunAll()
	if output.CISOView == nil {
		t.Fatal("expected CISO view")
	}
	if len(output.CISOView.TopAssumptionsToVerify) == 0 {
		t.Error("expected top assumptions to verify")
	}
}

func TestEvidenceMatching(t *testing.T) {
	keywords := []string{"mfa", "multi-factor"}
	evidence := lookupEvidence(keywords, EvCatIdentity)
	if len(evidence) == 0 {
		t.Error("expected evidence for MFA keywords")
	}
	foundMfa := false
	for _, ev := range evidence {
		if strings.Contains(ev.Name, "MFA") {
			foundMfa = true
			break
		}
	}
	if !foundMfa {
		t.Error("expected MFA Policy evidence for MFA keywords")
	}

	actions := lookupActions(keywords, EvCatIdentity)
	if len(actions) == 0 {
		t.Error("expected actions for MFA keywords")
	}
}

func TestExtractKeywords(t *testing.T) {
	desc := "MFA is enforced for all users accessing the admin portal"
	keywords := extractKeywords(desc)
	if len(keywords) == 0 {
		t.Error("expected extracted keywords")
	}
	found := false
	for _, kw := range keywords {
		if kw == "mfa" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected 'mfa' in extracted keywords")
	}
}

func TestExportMarkdown(t *testing.T) {
	inputs := []VerificationInput{
		{ID: "A1", Description: "MFA is enforced", Category: "identity", Risk: "Critical", Keywords: []string{"mfa"}},
	}
	engine := NewVerificationEngine("general", nil, inputs)
	output := engine.RunAll()
	md := ExportMarkdown(output)
	if md == "" {
		t.Error("expected non-empty markdown")
	}
	if !strings.Contains(md, "Verification Intelligence Report") {
		t.Error("expected report title in markdown")
	}
	if !strings.Contains(md, "MFA is enforced") {
		t.Error("expected assumption text in markdown")
	}
}

func TestExportHTML(t *testing.T) {
	inputs := []VerificationInput{
		{ID: "A1", Description: "RBAC is configured", Category: "authorization", Risk: "High", Keywords: []string{"rbac"}},
	}
	engine := NewVerificationEngine("general", nil, inputs)
	output := engine.RunAll()
	html := ExportHTML(output)
	if html == "" {
		t.Error("expected non-empty HTML")
	}
	if !strings.Contains(html, "Verification Intelligence Report") {
		t.Error("expected report title in HTML")
	}
	if !strings.Contains(html, "RBAC") {
		t.Error("expected RBAC reference in HTML")
	}
}

func TestVerificationPrecision(t *testing.T) {
	inputs := []VerificationInput{
		{ID: "A1", Description: "MFA enforced", Category: "identity", Risk: "Critical", Keywords: []string{"mfa"}},
		{ID: "A2", Description: "KMS key rotation", Category: "cryptography", Risk: "Critical", Keywords: []string{"kms", "rotation"}},
		{ID: "A3", Description: "SIEM alerting", Category: "monitoring", Risk: "High", Keywords: []string{"siem", "alert"}},
		{ID: "A4", Description: "Backup restore testing", Category: "resilience", Risk: "High", Keywords: []string{"backup", "restore"}},
		{ID: "A5", Description: "Rate limiting", Category: "operational", Risk: "Medium", Keywords: []string{"rate limit"}},
		{ID: "A6", Description: "Auth0 vendor security", Category: "third_party", Risk: "High", Keywords: []string{"auth0", "vendor"}},
	}
	engine := NewVerificationEngine("general", nil, inputs)
	output := engine.RunAll()
	if output.Assessment == nil || len(output.Assessment.Plans) != 6 {
		t.Fatalf("expected 6 plans, got %d", len(output.Assessment.Plans))
	}

	for _, p := range output.Assessment.Plans {
		t.Logf("assumption %s (%s): confidence=%.0f priority=%s status=%s evidence=%d actions=%d",
			p.AssumptionID, p.Risk, p.Confidence, p.Priority, p.Status,
			len(p.EvidenceRequired), len(p.Actions))

		if p.Confidence < 0 || p.Confidence > 100 {
			t.Errorf("confidence out of range for %s: %.0f", p.AssumptionID, p.Confidence)
		}
		if len(p.EvidenceRequired) == 0 {
			t.Errorf("no evidence for %s", p.AssumptionID)
		}
		if len(p.Actions) == 0 {
			t.Errorf("no actions for %s", p.AssumptionID)
		}
	}
}

func TestBenchmarkVerification(t *testing.T) {
	inputs := make([]VerificationInput, 100)
	for i := 0; i < 100; i++ {
		cat := []string{"identity", "authorization", "cryptography", "monitoring", "resilience", "third_party", "operational"}
		risks := []string{"Critical", "High", "Medium", "Low"}
		kws := [][]string{
			{"mfa"}, {"rbac"}, {"tls"}, {"siem"}, {"backup"}, {"auth0"}, {"secret"},
		}
		inputs[i] = VerificationInput{
			ID:          fmt.Sprintf("A-%d", i),
			Description: fmt.Sprintf("Test assumption %d for %s", i, cat[i%7]),
			Category:    cat[i%7],
			Risk:        risks[i%4],
			Keywords:    kws[i%7],
		}
	}
	engine := NewVerificationEngine("healthcare", []string{"auth0", "database", "kms"}, inputs)
	output := engine.RunAll()
	if output.Assessment == nil || len(output.Assessment.Plans) != 100 {
		t.Fatalf("expected 100 plans, got %d", len(output.Assessment.Plans))
	}
	if output.CISOView == nil {
		t.Fatal("expected CISO view")
	}
	if len(output.Roadmaps) == 0 {
		t.Fatal("expected roadmaps")
	}
}

func BenchmarkVerificationEngine(b *testing.B) {
	inputs := make([]VerificationInput, 50)
	for i := 0; i < 50; i++ {
		cat := []string{"identity", "authorization", "cryptography", "monitoring", "resilience", "third_party", "operational"}
		risks := []string{"Critical", "High", "Medium", "Low"}
		kws := [][]string{
			{"mfa"}, {"rbac"}, {"tls", "key"}, {"siem", "log"}, {"backup", "restore"}, {"auth0", "vendor"}, {"secret", "credential"},
		}
		inputs[i] = VerificationInput{
			ID:          fmt.Sprintf("A-%d", i),
			Description: fmt.Sprintf("Benchmark assumption %d", i),
			Category:    cat[i%7],
			Risk:        risks[i%4],
			Keywords:    kws[i%7],
		}
	}
	engine := NewVerificationEngine("cloud", []string{"auth0", "kms", "database"}, inputs)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.RunAll()
	}
}

func BenchmarkVerificationLarge(b *testing.B) {
	inputs := make([]VerificationInput, 500)
	for i := 0; i < 500; i++ {
		cat := []string{"identity", "authorization", "cryptography", "monitoring", "resilience", "third_party", "operational"}
		inputs[i] = VerificationInput{
			ID:          fmt.Sprintf("A-%d", i),
			Description: fmt.Sprintf("Benchmark large assumption %d", i),
			Category:    cat[i%7],
			Risk:        "High",
			Keywords:    []string{"test", "verify"},
		}
	}
	engine := NewVerificationEngine("kubernetes", []string{"auth0", "kms", "database", "siem", "backup"}, inputs)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine.RunAll()
	}
}
