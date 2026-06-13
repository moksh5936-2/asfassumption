package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func init() {
	asfLog = log.New(os.Stderr, "[test] ", log.Ltime)
}

func TestBaselineAsftestYAML(t *testing.T) {
	path := "testdata/asftest.yaml"

	// First, parse the architecture to see explicit assumptions
	desc, err := ParseArchitecture(path)
	if err != nil {
		t.Fatalf("ParseArchitecture: %v", err)
	}

	t.Logf("Name: %s", desc.Name)
	t.Logf("Components: %d", len(desc.Components))
	t.Logf("Relationships: %d", len(desc.Relationships))
	t.Logf("ExplicitAssumptions: %d", len(desc.ExplicitAssumptions))
	t.Logf("SecurityControls categories: %d", len(desc.SecurityControls))
	t.Logf("Compliance: %v", desc.Compliance)
	t.Logf("RawText length: %d bytes", len(desc.RawText))

	// Write raw text to inspect
	os.WriteFile("/tmp/baseline_rawtext_test.txt", []byte(desc.RawText), 0644)

	// Now run the full engine analysis
	cfg := &Config{}
	e := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for p := range progress {
			t.Logf("Progress: %.0f%% - %s", p.Percent, p.Stage)
		}
	}()

	result, err := e.RunAnalysis(path, "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis: %v", err)
	}

	t.Logf("ArchitectureName: %s", result.ArchitectureName)
	t.Logf("TotalAssumptions: %d", result.TotalAssumptions)
	t.Logf("CriticalCount: %d", result.CriticalCount)
	t.Logf("HighCount: %d", result.HighCount)
	t.Logf("MediumCount: %d", result.MediumCount)
	t.Logf("LowCount: %d", result.LowCount)

	t.Logf("Assumptions:")
	for _, a := range result.Assumptions {
		t.Logf("  [%s] Risk=%-8s Conf=%.0f%% Stride=%v Cat=%-15s %s",
			a.ID, a.Risk, a.Confidence*100, a.Stride, a.Category, truncateStr(a.Description, 80))
	}

	t.Logf("Controls:")
	for _, c := range result.Controls {
		t.Logf("  [%s] %s", c.ID, c.Description)
	}

	t.Logf("StrideDistribution:")
	for k, v := range result.StrideDistribution {
		t.Logf("  %s: %d", k, v)
	}

	t.Logf("Summary: %s", result.Summary)
	t.Logf("Compliance: %v", result.Compliance)

	// Save as JSON
	jenc := json.NewEncoder(os.Stdout)
	jenc.SetIndent("", "  ")
	fmt.Fprintln(os.Stderr, "\n=== FULL JSON OUTPUT ===")
	if err := jenc.Encode(result); err != nil {
		t.Logf("JSON encode: %v", err)
	}

	// Save to file
	f, err := os.Create("/tmp/baseline_engine_result.json")
	if err != nil {
		t.Fatalf("create output file: %v", err)
	}
	defer f.Close()
	jenc2 := json.NewEncoder(f)
	jenc2.SetIndent("", "  ")
	jenc2.Encode(result)

	// Summary assertions
	if result.TotalAssumptions < 25 {
		t.Errorf("expected >=25 assumptions, got %d", result.TotalAssumptions)
	}

	// Verification distribution assertion
	verificationStates := map[string]int{}
	for _, a := range result.Assumptions {
		state := a.VerificationStatus
		if state == "" {
			state = "EMPTY"
		}
		verificationStates[state]++
	}
	t.Logf("Verification distribution: %v", verificationStates)
	if verificationStates["PARTIALLY_VERIFIED"] == 0 && verificationStates["CONTRADICTED"] == 0 && verificationStates["VERIFIED"] == 0 {
		t.Errorf("verification: expected at least one assumption with non-empty verification status")
	} else {
		t.Logf("PASS: verification engine produces meaningful statuses")
	}
}

func TestAdversarialBaseline(t *testing.T) {
	path := "testdata/adversarial_insecure.yaml"

	desc, err := ParseArchitecture(path)
	if err != nil {
		t.Fatalf("ParseArchitecture: %v", err)
	}

	t.Logf("Name: %s", desc.Name)
	t.Logf("Components: %d", len(desc.Components))
	t.Logf("Relationships: %d", len(desc.Relationships))
	t.Logf("ExplicitAssumptions: %d", len(desc.ExplicitAssumptions))
	t.Logf("SecurityControls categories: %d", len(desc.SecurityControls))
	t.Logf("RawText length: %d bytes", len(desc.RawText))

	cfg := &Config{}
	e := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() {
		for p := range progress {
			t.Logf("Progress: %.0f%% - %s", p.Percent, p.Stage)
		}
	}()

	result, err := e.RunAnalysis(path, "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis: %v", err)
	}

	t.Logf("ArchitectureName: %s", result.ArchitectureName)
	t.Logf("TotalAssumptions: %d", result.TotalAssumptions)
	t.Logf("CriticalCount: %d", result.CriticalCount)
	t.Logf("HighCount: %d", result.HighCount)
	t.Logf("MediumCount: %d", result.MediumCount)
	t.Logf("LowCount: %d", result.LowCount)

	t.Logf("Adversarial Assumptions:")
	for _, a := range result.Assumptions {
		t.Logf("  [%s] Risk=%-8s Conf=%.0f%% Ver=%-15s Stride=%v Cat=%-15s %s",
			a.ID, a.Risk, a.Confidence*100, a.VerificationStatus, a.Stride, a.Category, truncateStr(a.Description, 80))
	}

	t.Logf("Adversarial Controls:")
	for _, c := range result.Controls {
		t.Logf("  [%s] %s", c.ID, c.Description)
	}

	// Verify no parser pollution: no assumptions from markdown headers
	for _, a := range result.Assumptions {
		if a.Category == "DOCUMENTATION" || a.Category == "" {
			if strings.Contains(a.Description, "Architecture") || strings.Contains(a.Description, "##") {
				t.Errorf("parser pollution: assumption %s has DOCUMENTATION type with header text: %s", a.ID, truncateStr(a.Description, 80))
			}
		}
		if strings.Contains(a.Description, "# Architecture:") {
			t.Errorf("parser pollution: assumption %s contains # Architecture: header: %s", a.ID, truncateStr(a.Description, 80))
		}
	}

	// ── Required Assertions ──────────────────────────────────────────

	// 1. No parser pollution: no DOCUMENTATION-type with header text
	for _, a := range result.Assumptions {
		if a.Category == "DOCUMENTATION" || a.Category == "" {
			if strings.Contains(a.Description, "Architecture") || strings.Contains(a.Description, "##") {
				t.Errorf("parser pollution: assumption %s has DOCUMENTATION type with header text: %s", a.ID, truncateStr(a.Description, 80))
			}
		}
		if strings.Contains(a.Description, "# Architecture:") {
			t.Errorf("parser pollution: assumption %s contains # Architecture: header: %s", a.ID, truncateStr(a.Description, 80))
		}
	}
	t.Logf("PASS: No parser pollution detected")

	// 2. Fact protection: assumptions that contradict explicit facts must be transformed
	transformedCount := 0
	for _, a := range result.Assumptions {
		if strings.Contains(a.Description, "Plaintext communication is expected") ||
			strings.Contains(a.Description, "Single-factor or weak authentication") ||
			strings.Contains(a.Description, "Flat network") ||
			strings.Contains(a.Description, "compensating controls or accepted risk") {
			transformedCount++
			if a.VerificationStatus != "CONTRADICTED" {
				t.Errorf("fact protection: transformed assumption %s should be CONTRADICTED, got %s", a.ID, a.VerificationStatus)
			}
		}
	}
	if transformedCount == 0 {
		t.Errorf("fact protection: expected at least 1 transformed assumption for insecure architecture")
	} else {
		t.Logf("PASS: %d assumptions transformed by fact protection", transformedCount)
	}

	// 3. Verification distribution must not be all empty/UNKNOWN
	verificationStates := map[string]int{}
	for _, a := range result.Assumptions {
		state := a.VerificationStatus
		if state == "" {
			state = "EMPTY"
		}
		verificationStates[state]++
	}
	t.Logf("Verification distribution: %v", verificationStates)
	if verificationStates["CONTRADICTED"] == 0 {
		t.Errorf("verification: expected at least 1 CONTRADICTED assumption, got 0")
	} else {
		t.Logf("PASS: verification engine produces CONTRADICTED status")
	}

	// 4. Trust chains must be non-zero for multi-component architecture
	if result.TrustOutput == nil {
		t.Errorf("trust chain: TrustOutput is nil")
	} else {
		t.Logf("TrustChains: %d", len(result.TrustOutput.TrustChains))
		t.Logf("FailureCascades: %d", len(result.TrustOutput.FailureCascades))
		t.Logf("CriticalAssumptions: %d", len(result.TrustOutput.CriticalAssumptions))
		t.Logf("SPOTF: %d", len(result.TrustOutput.SinglePointsOfTrust))
		if len(result.TrustOutput.TrustChains) == 0 {
			t.Errorf("trust chain: expected >0 trust chains, got 0")
		} else {
			t.Logf("PASS: trust chain engine produces %d chains", len(result.TrustOutput.TrustChains))
		}
		for _, tc := range result.TrustOutput.TrustChains {
			t.Logf("  Chain: %s (risk=%s, len=%d, conf=%.2f)", tc.ID, tc.Risk, tc.Length, tc.Confidence)
		}
	}

	// 5. Risk severity: insecure patterns must produce Critical/High
	if result.CriticalCount < 3 {
		t.Errorf("risk calibration: expected >=3 critical-risk assumptions, got %d", result.CriticalCount)
	} else {
		t.Logf("PASS: risk calibration produces %d Critical assumptions", result.CriticalCount)
	}
	if result.HighCount < 5 {
		t.Errorf("risk calibration: expected >=5 high-risk assumptions, got %d", result.HighCount)
	} else {
		t.Logf("PASS: risk calibration produces %d High assumptions", result.HighCount)
	}

	// 6. Contradiction deduplication: check no duplicate contradiction keys
	if result.Contradictions != nil {
		seenContra := map[string]bool{}
		for _, c := range result.Contradictions {
			key := c.RuleName + "|" + strings.Join(c.AffectedAssumptions, ",")
			if seenContra[key] {
				t.Errorf("contradiction dedup: duplicate contradiction key %s", key)
			}
			seenContra[key] = true
		}
		t.Logf("PASS: %d unique contradictions (no duplicates)", len(result.Contradictions))
	}
	if result.CIEContradictions != nil {
		seenCIEContra := map[string]bool{}
		for _, c := range result.CIEContradictions {
			key := c.Type + "|" + c.StatementA.ID + "|" + c.StatementB.ID
			if seenCIEContra[key] {
				t.Errorf("CIE contradiction dedup: duplicate CIE contradiction key %s", key)
			}
			seenCIEContra[key] = true
		}
		t.Logf("PASS: %d unique CIE contradictions (no duplicates)", len(result.CIEContradictions))
	}

	// 7. Coverage output must have components
	if result.CoverageOutput == nil {
		t.Errorf("coverage: CoverageOutput is nil")
	} else {
		componentsInCoverage := 0
		if result.CoverageOutput.Assessment != nil {
			componentsInCoverage = len(result.CoverageOutput.Assessment.ComponentResults)
		}
		blindSpotCount := len(result.CoverageOutput.BlindSpots)
		t.Logf("Coverage: %d component results, %d blind spots", componentsInCoverage, blindSpotCount)
		if blindSpotCount == 0 {
			t.Errorf("coverage: expected >0 blind spots for insecure architecture")
		} else {
			t.Logf("PASS: coverage engine detects %d blind spots", blindSpotCount)
		}
	}

	// 8. Review output must be populated
	if result.ReviewOutput == nil {
		t.Errorf("review: ReviewOutput is nil")
	} else {
		queueLen := result.ReviewOutput.Queue.TotalItems
		t.Logf("Review queue: %d items", queueLen)
		if queueLen == 0 {
			t.Errorf("review: expected >0 items in review queue")
		} else {
			t.Logf("PASS: review workbench has %d review queue items", queueLen)
		}
	}
}

func TestMarkdownParser(t *testing.T) {
	path := "testdata/asftest.md"

	desc, err := ParseArchitecture(path)
	if err != nil {
		t.Fatalf("ParseArchitecture: %v", err)
	}

	t.Logf("Name: %s", desc.Name)
	t.Logf("ExplicitAssumptions: %d", len(desc.ExplicitAssumptions))
	t.Logf("SecurityControls categories: %d", len(desc.SecurityControls))
	for cat, controls := range desc.SecurityControls {
		t.Logf("  SecurityControls[%s] = %v", cat, controls)
	}
	t.Logf("Compliance: %v", desc.Compliance)
	for i, c := range desc.Compliance {
		t.Logf("  Compliance[%d] = %q", i, c)
	}

	if len(desc.ExplicitAssumptions) < 25 {
		t.Errorf("expected >=25 explicit assumptions from Markdown, got %d", len(desc.ExplicitAssumptions))
	}

	for i, a := range desc.ExplicitAssumptions {
		t.Logf("  [%d] %s", i, a)
	}
}
