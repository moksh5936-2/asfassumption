package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
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
