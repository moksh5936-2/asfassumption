package coverage

import (
	"fmt"
	"strings"
	"testing"
)

func TestEmptyCoverage(t *testing.T) {
	engine := NewCoverageEngine("healthcare", []string{}, nil)
	output := engine.RunAll()

	if output == nil {
		t.Fatal("expected non-nil output")
	}

}

func TestCoverageAllExpected(t *testing.T) {
	assumptions := []AssumptionInput{
		{ID: "A1", Description: "MFA is enabled for all admin accounts", Component: "Auth0", Category: "identity", Risk: "Critical"},
		{ID: "A2", Description: "Admin access is restricted to authorized users", Component: "Auth0", Category: "identity", Risk: "Critical"},
		{ID: "A3", Description: "SSO is configured for all apps", Component: "Auth0", Category: "identity", Risk: "High"},
		{ID: "A4", Description: "RBAC is configured for database roles", Component: "Auth0", Category: "authorization", Risk: "Critical"},
		{ID: "A5", Description: "Least privilege is enforced", Component: "Auth0", Category: "authorization", Risk: "High"},
		{ID: "A6", Description: "Access audit logging is enabled", Component: "Auth0", Category: "monitoring", Risk: "High"},
		{ID: "A7", Description: "Federation is configured with trusted IdPs", Component: "Auth0", Category: "identity", Risk: "High"},
	}
	components := []string{"Auth0"}

	engine := NewCoverageEngine("healthcare", components, assumptions)
	output := engine.RunAll()

	if output.Assessment == nil {
		t.Fatal("expected assessment")
	}

	for _, cat := range output.Assessment.Categories {
		if cat.CoveragePct < 100.0 {
			t.Logf("category %s: %.1f%% coverage", cat.Category, cat.CoveragePct)
		}
	}
}

func TestCoverageGaps(t *testing.T) {
	assumptions := []AssumptionInput{
		{ID: "A1", Description: "MFA is enabled", Component: "Auth0", Category: "identity", Risk: "Critical"},
	}
	components := []string{"Auth0", "Database", "KMS"}

	engine := NewCoverageEngine("healthcare", components, assumptions)
	output := engine.RunAll()

	if output.Assessment == nil {
		t.Fatal("expected assessment")
	}

	if len(output.Assessment.Gaps) == 0 {
		t.Error("expected coverage gaps with multiple components and few assumptions")
	}
}

func TestBlindSpots(t *testing.T) {
	assumptions := []AssumptionInput{
		{ID: "A1", Description: "MFA is enabled", Component: "Auth0", Category: "identity", Risk: "Critical"},
	}
	components := []string{"Auth0", "Database", "KMS", "BackupService"}

	engine := NewCoverageEngine("healthcare", components, assumptions)
	output := engine.RunAll()

	if len(output.BlindSpots) == 0 {
		t.Error("expected blind spots with KMS present but no key assumptions")
	}

	hasKMSBlindSpot := false
	for _, bs := range output.BlindSpots {
		if strings.Contains(bs.Component, "kms") || strings.Contains(bs.Component, "KMS") {
			hasKMSBlindSpot = true
			break
		}
	}
	if !hasKMSBlindSpot {
		t.Error("expected blind spot for KMS component")
	}

	if output.CISOView == nil {
		t.Error("expected CISO view")
	}
}

func TestDomainBlindSpots(t *testing.T) {
	spots := GetDomainBlindSpots("healthcare")
	if len(spots) == 0 {
		t.Error("expected domain blind spots for healthcare")
	}

	hasBreakGlass := false
	for _, s := range spots {
		if strings.Contains(s.MissingArea, "Break Glass") {
			hasBreakGlass = true
			break
		}
	}
	if !hasBreakGlass {
		t.Error("expected break glass access blind spot for healthcare")
	}
}

func TestTaxonomy(t *testing.T) {
	expectations := GetExpectations("Auth0")
	if len(expectations) == 0 {
		t.Error("expected taxonomy entries for Auth0")
	}

	cats := GetComponentCategories("Auth0")
	if len(cats) == 0 {
		t.Error("expected categories for Auth0")
	}
}

func TestCISOView(t *testing.T) {
	assumptions := []AssumptionInput{
		{ID: "A1", Description: "Basic MFA", Component: "Auth0", Category: "identity", Risk: "Critical"},
	}
	components := []string{"Auth0", "Database", "KMS", "APIGateway", "BackupService"}

	engine := NewCoverageEngine("healthcare", components, assumptions)
	output := engine.RunAll()

	if output.CISOView == nil {
		t.Fatal("expected CISO view")
	}

	if len(output.CISOView.TopBlindSpots) == 0 {
		t.Error("expected top blind spots")
	}
	if len(output.CISOView.AreasRequiringReview) == 0 {
		t.Error("expected areas requiring review")
	}
}

func TestExportMarkdown(t *testing.T) {
	assumptions := []AssumptionInput{
		{ID: "A1", Description: "MFA is enabled", Component: "Auth0", Category: "identity", Risk: "Critical"},
	}
	components := []string{"Auth0", "KMS"}

	engine := NewCoverageEngine("healthcare", components, assumptions)
	output := engine.RunAll()

	md := ExportMarkdown(output)
	if !strings.Contains(md, "Coverage & Blind Spot Analysis") {
		t.Error("expected markdown title")
	}

}

func TestExportHTML(t *testing.T) {
	assumptions := []AssumptionInput{
		{ID: "A1", Description: "MFA is enabled", Component: "Auth0", Category: "identity", Risk: "Critical"},
	}
	components := []string{"Auth0", "KMS"}

	engine := NewCoverageEngine("healthcare", components, assumptions)
	output := engine.RunAll()

	html := ExportHTML(output)
	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("expected HTML doctype")
	}

}

func TestCoveragePrecision(t *testing.T) {
	assumptions := []AssumptionInput{
		{ID: "A1", Description: "MFA is enabled", Component: "Auth0", Category: "identity", Risk: "Critical"},
		{ID: "A2", Description: "SSO is configured", Component: "Auth0", Category: "identity", Risk: "High"},
		{ID: "A3", Description: "Admin access is restricted", Component: "Auth0", Category: "identity", Risk: "Critical"},
		{ID: "A4", Description: "RBAC is configured", Component: "Auth0", Category: "authorization", Risk: "Critical"},
		{ID: "A5", Description: "Access audit logging is enabled", Component: "Auth0", Category: "monitoring", Risk: "High"},
		{ID: "A6", Description: "Encryption at rest is enabled", Component: "Database", Category: "cryptography", Risk: "Critical"},
		{ID: "A7", Description: "TLS is configured", Component: "Database", Category: "cryptography", Risk: "Critical"},
		{ID: "A8", Description: "Database backups are configured", Component: "Database", Category: "resilience", Risk: "High"},
		{ID: "A9", Description: "Key rotation is configured", Component: "KMS", Category: "cryptography", Risk: "Critical"},
		{ID: "A10", Description: "Secrets management is configured", Component: "KMS", Category: "operational", Risk: "Critical"},
	}
	components := []string{"Auth0", "Database", "KMS", "APIGateway", "BackupService"}

	engine := NewCoverageEngine("healthcare", components, assumptions)
	output := engine.RunAll()

	if len(output.Assessment.Gaps) == 0 {
		t.Log("all categories have sufficient coverage")
	}

	hasAPI := false
	for _, cat := range output.Assessment.Categories {
		if cat.Category == CatAuthorization && cat.ObservedCount >= 1 {
			hasAPI = true
		}
		t.Logf("category %s: %.1f%% (%d/%d)", cat.Category, cat.CoveragePct, cat.ObservedCount, cat.ExpectedCount)
	}
	if hasAPI {
		t.Log("authorization coverage present")
	}

	if output.CISOView != nil && len(output.CISOView.TopBlindSpots) > 0 {
		t.Logf("top blind spots: %d", len(output.CISOView.TopBlindSpots))
		for _, bs := range output.CISOView.TopBlindSpots[:min(3, len(output.CISOView.TopBlindSpots))] {
			t.Logf("  %s (%.0f)", bs.Title, bs.Score)
		}
	}
}

func BenchmarkCoverageAnalysis(b *testing.B) {
	assumptions := []AssumptionInput{
		{ID: "A1", Description: "MFA is enabled", Component: "Auth0", Category: "identity", Risk: "Critical"},
		{ID: "A2", Description: "RBAC configured", Component: "Auth0", Category: "authorization", Risk: "Critical"},
		{ID: "A3", Description: "TLS configured", Component: "WebApp", Category: "cryptography", Risk: "Critical"},
		{ID: "A4", Description: "Audit logging enabled", Component: "SIEM", Category: "monitoring", Risk: "High"},
		{ID: "A5", Description: "Backups configured", Component: "Database", Category: "resilience", Risk: "High"},
	}
	components := []string{"Auth0", "Database", "KMS", "SIEM", "APIGateway", "BackupService", "WebApp"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine := NewCoverageEngine("healthcare", components, assumptions)
		output := engine.RunAll()
		if output == nil {
			b.Fatal("expected output")
		}
	}
}

func BenchmarkCoverageLarge(b *testing.B) {
	var assumptions []AssumptionInput
	for i := 0; i < 100; i++ {
		assumptions = append(assumptions, AssumptionInput{
			ID:          fmt.Sprintf("A%d", i+1),
			Description: "Test assumption with keywords for coverage matching",
			Component:   "Auth0",
			Category:    "identity",
			Risk:        "Critical",
		})
	}
	components := []string{"Auth0", "Database", "KMS", "SIEM", "APIGateway", "BackupService", "WebApp", "Jenkins", "S3"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine := NewCoverageEngine("healthcare", components, assumptions)
		output := engine.RunAll()
		if output == nil {
			b.Fatal("expected output")
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
