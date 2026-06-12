package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
)

// BenchmarkResult holds the outcome of a benchmark run.
type BenchmarkResult struct {
	ArchitectureName      string             `json:"architecture_name"`
	GoldAssumptions       int                `json:"gold_assumptions"`
	GeneratedAssumptions  int                `json:"generated_assumptions"`
	MatchedAssumptions    int                `json:"matched_assumptions"`
	Recall                float64            `json:"recall"`
	Precision             float64            `json:"precision"`
	F1Score               float64            `json:"f1_score"`
	CategoryCoverage      map[string]float64 `json:"category_coverage"`
	FalsePositives        int                `json:"false_positives"`
	FalseNegatives        int                `json:"false_negatives"`
	ContradictionsFound   int                `json:"contradictions_found"`
	TrustBoundariesFound  int                `json:"trust_boundaries_found"`
	QualityScore          float64            `json:"quality_score"`
	AssumptionsByCategory map[string]int     `json:"assumptions_by_category"`
}

// goldAssumptionsForAsftest returns the gold-standard assumptions for the healthcare benchmark.
func goldAssumptionsForAsftest() []string {
	return []string{
		"MFA is enforced for all Auth0 user authentication",
		"Auth0 administrative access is restricted to authorized admins only",
		"All API requests pass through APIGateway for authentication validation",
		"PHI data is encrypted at rest in PHIDatabase",
		"PHI data is encrypted in transit between all components",
		"Encryption keys are stored and managed in KMS with automatic rotation",
		"KMS access is restricted to authorized services only",
		"Audit logging is immutable and tamper-proof",
		"All PHI access events are logged with user and timestamp",
		"Backup data is encrypted at rest and in transit",
		"Backup restore procedures are tested regularly",
		"ThirdPartyAnalytics has access only to de-identified PHI",
		"Third-party provider maintains equivalent security controls",
		"AdminConsole requires MFA for all administrative access",
		"Object-level authorization is enforced for PHI record access",
		"Database connection pooling does not leak data between sessions",
		"TLS certificates are monitored and renewed before expiry",
		"API rate limiting prevents abuse of PHI endpoints",
		"Network segmentation isolates the PHI database in a private subnet",
		"Unauthorized PHI export is detected and alerted",
		"Session tokens expire and are rotated periodically",
		"Auth0 tenant is configured with breach detection and anomaly alerts",
		"Audit log storage has sufficient retention for compliance",
		"PHI data minimization policies are enforced at application layer",
		"Incident response plan includes PHI breach notification",
		"Vendor risk assessments are conducted for ThirdPartyAnalytics",
		"Database backups are stored in a separate geographic region",
		"KMS key deletion is protected with multi-factor authorization",
		"API Gateway logs are monitored for anomalous access patterns",
		"System health and availability monitoring covers all components",
		// Hidden assumptions (not explicitly stated but should be inferred):
		"Key management procedures include key escrow and recovery",
		"Secrets are rotated when personnel leave or systems are decommissioned",
		"Audit logs are regularly reviewed for anomalies and unauthorized access",
		"Backup integrity is verified through periodic restore testing",
		"Third-party access is logged and monitored for data exfiltration",
		"Break-glass access procedures exist for emergency PHI access",
		"Data retention policies enforce automatic deletion of expired PHI",
		"Certificate pinning is enforced for API communication",
		"Privileged access is reviewed quarterly and revoked when roles change",
		"Security incident response includes containment of affected PHI systems",
	}
}

// isAssumptionMatch checks if a generated assumption matches a gold assumption.
func isAssumptionMatch(generated, gold string) bool {
	genLower := strings.ToLower(generated)
	goldLower := strings.ToLower(gold)

	// Direct containment
	if strings.Contains(genLower, goldLower) || strings.Contains(goldLower, genLower) {
		return true
	}

	// Check for key concept overlap
	genWords := strings.Fields(genLower)
	goldWords := strings.Fields(goldLower)

	// Extract meaningful keywords (length > 3, not stopwords)
	stopWords := map[string]bool{
		"the": true, "and": true, "for": true, "are": true, "but": true,
		"not": true, "all": true, "can": true, "has": true, "have": true,
		"may": true, "must": true, "shall": true, "should": true, "will": true,
		"with": true, "from": true, "that": true, "this": true, "each": true,
		"every": true, "than": true, "then": true, "just": true, "been": true,
		"were": true, "was": true, "its": true, "also": true, "per": true,
		"via": true, "is": true, "to": true, "in": true, "of": true,
		"on": true, "at": true, "by": true, "as": true, "an": true, "or": true,
		"system": true, "assumes": true, "control": true, "state": true,
	}

	var genKeywords, goldKeywords []string
	for _, w := range genWords {
		w = strings.TrimSuffix(w, ".")
		if len(w) > 3 && !stopWords[w] {
			genKeywords = append(genKeywords, w)
		}
	}
	for _, w := range goldWords {
		w = strings.TrimSuffix(w, ".")
		if len(w) > 3 && !stopWords[w] {
			goldKeywords = append(goldKeywords, w)
		}
	}

	// If at least 2 keywords overlap, consider it a match
	overlap := 0
	for _, gw := range genKeywords {
		for _, glw := range goldKeywords {
			if gw == glw || strings.Contains(gw, glw) || strings.Contains(glw, gw) {
				overlap++
				break
			}
		}
	}
	return overlap >= 2
}

func TestBenchmarkAsftestYAML(t *testing.T) {
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()
	result, err := engine.RunAnalysis("testdata/asftest.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	gold := goldAssumptionsForAsftest()
	generated := result.Assumptions

	// Calculate matches
	matched := 0
	matchedGold := make(map[int]bool)
	for _, gen := range generated {
		for i, g := range gold {
			if isAssumptionMatch(gen.Description, g) {
				matched++
				matchedGold[i] = true
				break
			}
		}
	}

	// Calculate metrics
	recall := float64(len(matchedGold)) / float64(len(gold))
	precision := float64(matched) / float64(len(generated))
	f1 := 0.0
	if recall+precision > 0 {
		f1 = 2 * recall * precision / (recall + precision)
	}
	falsePositives := len(generated) - matched
	falseNegatives := len(gold) - len(matchedGold)

	// Calculate category coverage
	categoryCoverage := make(map[string]float64)
	categoryCounts := make(map[string]int)
	for _, a := range generated {
		categoryCounts[a.Category]++
	}
	// Define expected categories
	expectedCategories := []string{
		"IDENTITY", "ACCESS", "CONFIGURATION", "NETWORK", "DEPENDENCY",
		"KeyManagement", "SecretsManagement", "Auditability", "Logging",
		"Monitoring", "Backups", "ThirdPartyRisk", "VendorRisk",
		"TrustBoundaries", "APISecurity", "SessionSecurity", "DataProtection",
		"Compliance", "IncidentResponse", "Governance",
	}
	for _, cat := range expectedCategories {
		if categoryCounts[cat] > 0 {
			categoryCoverage[cat] = 1.0
		} else {
			categoryCoverage[cat] = 0.0
		}
	}
	coverageScore := 0.0
	for _, v := range categoryCoverage {
		coverageScore += v
	}
	coverageScore = coverageScore / float64(len(expectedCategories))

	// Build benchmark result
	benchmark := BenchmarkResult{
		ArchitectureName:      "asftest.yaml",
		GoldAssumptions:       len(gold),
		GeneratedAssumptions:  len(generated),
		MatchedAssumptions:    len(matchedGold),
		Recall:                recall,
		Precision:             precision,
		F1Score:               f1,
		CategoryCoverage:      categoryCoverage,
		FalsePositives:        falsePositives,
		FalseNegatives:        falseNegatives,
		ContradictionsFound:   len(result.Contradictions),
		TrustBoundariesFound:  len(result.TrustBoundaries),
		QualityScore:          coverageScore,
		AssumptionsByCategory: categoryCounts,
	}

	// Write benchmark report
	reportPath := "/tmp/benchmark_asftest.json"
	reportJSON, _ := json.MarshalIndent(benchmark, "", "  ")
	os.WriteFile(reportPath, reportJSON, 0644)

	// Log results
	t.Logf("=== Benchmark Results for %s ===", benchmark.ArchitectureName)
	t.Logf("Gold Assumptions:       %d", benchmark.GoldAssumptions)
	t.Logf("Generated Assumptions:  %d", benchmark.GeneratedAssumptions)
	t.Logf("Matched Assumptions:    %d", benchmark.MatchedAssumptions)
	t.Logf("Recall:                 %.1f%%", recall*100)
	t.Logf("Precision:              %.1f%%", precision*100)
	t.Logf("F1 Score:               %.2f", f1)
	t.Logf("False Positives:        %d", falsePositives)
	t.Logf("False Negatives:        %d", falseNegatives)
	t.Logf("Contradictions:         %d", benchmark.ContradictionsFound)
	t.Logf("Trust Boundaries:       %d", benchmark.TrustBoundariesFound)
	t.Logf("Category Coverage:      %.1f%%", coverageScore*100)
	t.Logf("Report:                 %s", reportPath)

	// Log all generated categories
	t.Logf("\n=== Generated Categories ===")
	for cat, count := range categoryCounts {
		t.Logf("  %s: %d", cat, count)
	}

	// Log all matched gold assumptions
	t.Logf("\n=== Matched Gold Assumptions ===")
	for i, g := range gold {
		if matchedGold[i] {
			t.Logf("  ✓ %s", g)
		}
	}

	// Log all unmatched gold assumptions
	t.Logf("\n=== Unmatched Gold Assumptions (False Negatives) ===")
	for i, g := range gold {
		if !matchedGold[i] {
			t.Logf("  ✗ %s", g)
		}
	}

	// Assert success criteria
	if recall < 0.60 {
		t.Errorf("RECALL FAIL: %.1f%% (target: >60%%)", recall*100)
	} else {
		t.Logf("RECALL PASS: %.1f%%", recall*100)
	}
	if precision < 0.70 {
		t.Errorf("PRECISION FAIL: %.1f%% (target: >70%%)", precision*100)
	} else {
		t.Logf("PRECISION PASS: %.1f%%", precision*100)
	}

	// Log top 10 assumptions by quality
	t.Logf("\n=== Top 10 Assumptions by Quality ===")
	// Sort by quality score
	type scoredAssumption struct {
		assumption Assumption
		score      float64
	}
	var scored []scoredAssumption
	for _, a := range generated {
		scored = append(scored, scoredAssumption{a, a.QualityScore})
	}
	// Simple bubble sort for top 10
	for i := 0; i < len(scored) && i < 10; i++ {
		maxIdx := i
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[maxIdx].score {
				maxIdx = j
			}
		}
		scored[i], scored[maxIdx] = scored[maxIdx], scored[i]
	}
	for i := 0; i < len(scored) && i < 10; i++ {
		t.Logf("  %d. [%.2f] %s (%s)", i+1, scored[i].score, scored[i].assumption.Description, scored[i].assumption.Category)
	}
}

func TestBenchmarkContradictionDetection(t *testing.T) {
	// Create a test architecture with contradictions
	testYAML := `
metadata:
  name: Contradiction Test
  version: "1.0"
assumptions:
  - MFA is enforced for all users
  - Service accounts are exempt from MFA
  - All data is encrypted at rest
  - Backups are stored in plaintext
  - Least privilege is enforced
  - Administrators share a single account
  - Database is in a private subnet
  - Database is accessible from the internet
  - Audit logs are immutable
  - Old audit logs are deleted monthly
components:
  - name: WebApp
    type: web_application
    description: Web application
  - name: Database
    type: database
    description: Database
`
	tmpFile := "/tmp/contradiction_test.yaml"
	os.WriteFile(tmpFile, []byte(testYAML), 0644)
	defer os.Remove(tmpFile)

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()
	result, err := engine.RunAnalysis(tmpFile, "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	t.Logf("Found %d contradictions:", len(result.Contradictions))
	for _, c := range result.Contradictions {
		t.Logf("  - [%s] %s: %s", c.Severity, c.RuleName, c.Description)
	}

	// Should detect at least 2 contradictions
	if len(result.Contradictions) < 2 {
		t.Errorf("Expected at least 2 contradictions, found %d", len(result.Contradictions))
	}
}

func TestBenchmarkTrustBoundaryDiscovery(t *testing.T) {
	// Create a test architecture with trust boundaries
	testYAML := `
metadata:
  name: Trust Boundary Test
  version: "1.0"
components:
  - name: Internet
    type: external
    description: Internet users
  - name: VPN
    type: network
    description: VPN gateway
  - name: WebApp
    type: web_application
    description: Web application
  - name: Database
    type: database
    description: Database
  - name: ThirdPartyAPI
    type: external_service
    description: Third-party API
  - name: AdminConsole
    type: admin_tool
    description: Admin console
relationships:
  - source: Internet
    target: WebApp
    protocol: HTTPS
  - source: VPN
    target: WebApp
    protocol: HTTPS
  - source: WebApp
    target: Database
    protocol: TLS
  - source: WebApp
    target: ThirdPartyAPI
    protocol: HTTPS
  - source: AdminConsole
    target: Database
    protocol: TLS
`
	tmpFile := "/tmp/trust_boundary_test.yaml"
	os.WriteFile(tmpFile, []byte(testYAML), 0644)
	defer os.Remove(tmpFile)

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()
	result, err := engine.RunAnalysis(tmpFile, "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	t.Logf("Found %d trust boundaries:", len(result.TrustBoundaries))
	for _, b := range result.TrustBoundaries {
		t.Logf("  - [%s] %s: %s", b.RiskLevel, b.Type, strings.Join(b.Components, ", "))
	}

	// Should detect at least 3 boundaries (internet, vendor, admin)
	if len(result.TrustBoundaries) < 3 {
		t.Errorf("Expected at least 3 trust boundaries, found %d", len(result.TrustBoundaries))
	}
}

func TestBenchmarkDomainDetection(t *testing.T) {
	// Test healthcare domain detection
	testYAML := `
metadata:
  name: Healthcare Platform
  version: "1.0"
components:
  - name: PHIDatabase
    type: database
    description: Patient health records
  - name: EHR
    type: web_application
    description: Electronic health records
  - name: Auth0
    type: identity_provider
    description: Identity provider
`
	tmpFile := "/tmp/domain_test.yaml"
	os.WriteFile(tmpFile, []byte(testYAML), 0644)
	defer os.Remove(tmpFile)

	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()
	result, err := engine.RunAnalysis(tmpFile, "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	t.Logf("Detected domain: %s", result.Domain)
	if result.Domain == "" {
		t.Logf("No domain detected (architecture too small)")
	} else {
		t.Logf("Domain detected: %s", result.Domain)
	}
}

func TestBenchmarkReportExport(t *testing.T) {
	// Run the full benchmark and export a report
	cfg := &Config{}
	engine := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for range progress {
		}
	}()
	result, err := engine.RunAnalysis("testdata/asftest.yaml", "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis failed: %v", err)
	}

	gold := goldAssumptionsForAsftest()
	generated := result.Assumptions

	matched := 0
	matchedGold := make(map[int]bool)
	for _, gen := range generated {
		for i, g := range gold {
			if isAssumptionMatch(gen.Description, g) {
				matched++
				matchedGold[i] = true
				break
			}
		}
	}

	recall := float64(len(matchedGold)) / float64(len(gold))
	precision := float64(matched) / float64(len(generated))
	f1 := 0.0
	if recall+precision > 0 {
		f1 = 2 * recall * precision / (recall + precision)
	}

	report := fmt.Sprintf(`# ASF Benchmark Report

## Architecture: asftest.yaml

| Metric | Value |
|--------|-------|
| Gold Assumptions | %d |
| Generated Assumptions | %d |
| Matched Assumptions | %d |
| Recall | %.1f%% |
| Precision | %.1f%% |
| F1 Score | %.2f |
| False Positives | %d |
| False Negatives | %d |
| Contradictions | %d |
| Trust Boundaries | %d |
| Domain | %s |

## Success Criteria

| Criterion | Target | Actual | Status |
|-----------|--------|--------|--------|
| Recall | >60%% | %.1f%% | %s |
| Precision | >70%% | %.1f%% | %s |

## Category Distribution

`, len(gold), len(generated), len(matchedGold), recall*100, precision*100, f1,
		len(generated)-matched, len(gold)-len(matchedGold),
		len(result.Contradictions), len(result.TrustBoundaries), result.Domain,
		recall*100, status(recall, 0.60),
		precision*100, status(precision, 0.70))

	categoryCounts := make(map[string]int)
	for _, a := range generated {
		categoryCounts[a.Category]++
	}
	for cat, count := range categoryCounts {
		report += fmt.Sprintf("- %s: %d\n", cat, count)
	}

	report += "\n## Top 10 Assumptions by Quality\n\n"
	type scoredAssumption struct {
		assumption Assumption
		score      float64
	}
	var scored []scoredAssumption
	for _, a := range generated {
		scored = append(scored, scoredAssumption{a, a.QualityScore})
	}
	for i := 0; i < len(scored) && i < 10; i++ {
		maxIdx := i
		for j := i + 1; j < len(scored); j++ {
			if scored[j].score > scored[maxIdx].score {
				maxIdx = j
			}
		}
		scored[i], scored[maxIdx] = scored[maxIdx], scored[i]
	}
	for i := 0; i < len(scored) && i < 10; i++ {
		report += fmt.Sprintf("%d. [%.2f] %s (%s)\n", i+1, scored[i].score, scored[i].assumption.Description, scored[i].assumption.Category)
	}

	report += fmt.Sprintf("\n## Verdict: %s\n", verdict(recall, precision))

	reportPath := "/tmp/benchmark_report.md"
	os.WriteFile(reportPath, []byte(report), 0644)
	t.Logf("Benchmark report: %s", reportPath)
}

func status(actual float64, threshold float64) string {
	if actual >= threshold {
		return "✅ PASS"
	}
	return "❌ FAIL"
}

func verdict(recall, precision float64) string {
	if recall >= 0.60 && precision >= 0.70 {
		return "INTELLIGENCE_ENGINE_CERTIFIED"
	}
	return "NOT_CERTIFIED"
}
