package main

import (
	"os"
	"testing"
)

func TestExportAllFormats(t *testing.T) {
	path := "testdata/asftest.yaml"
	cfg := &Config{}
	e := NewEngine(cfg)
	progress := make(chan AnalysisProgress, 100)
	go func() { for range progress {} }()
	result, err := e.RunAnalysis(path, "", ModeASFOnly, progress)
	if err != nil {
		t.Fatalf("RunAnalysis: %v", err)
	}

	for _, format := range []ExportFormat{ExportJSON, ExportMarkdown, ExportCSV, ExportHTML, ExportPDF} {
		outPath, err := ExportResult(result, format, "/tmp")
		if err != nil {
			t.Errorf("Export %s: %v", format, err)
			continue
		}
		info, err := os.Stat(outPath)
		if err != nil {
			t.Errorf("Export %s: file not found at %s", format, outPath)
			continue
		}
		if info.Size() == 0 {
			t.Errorf("Export %s: file is empty", format)
			continue
		}
		t.Logf("Export %s: %s (%d bytes)", format, outPath, info.Size())
	}
}
