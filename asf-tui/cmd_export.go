package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func runExportCLI(args []string) {
	format := "json"
	outputDir := "./reports"
	includeNarrative := false
	includeTrust := false
	var filePath string

	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: asf export <file> [-f format] [-o output_dir] [--narrative] [--trust]")
			fmt.Println()
			fmt.Println("Export analysis results to various formats.")
			fmt.Println()
			fmt.Println("Arguments:")
			fmt.Println("  <file>                    Analysis result file (JSON)")
			fmt.Println("  -f, --format <format>     Output format: json, markdown, csv, html, pdf, jsonl")
			fmt.Println("  -o, --output <dir>        Output directory (default: ./reports)")
			fmt.Println("  --narrative               Include narrative output")
			fmt.Println("  --trust                   Include trust chain output")
			fmt.Println("  --help, -h                Show this help")
			os.Exit(ExitSuccess)
		}
	}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "-f", "--format":
			if i+1 < len(args) {
				i++
				format = args[i]
			}
		case "-o", "--output":
			if i+1 < len(args) {
				i++
				outputDir = args[i]
			}
		case "--narrative":
			includeNarrative = true
		case "--trust":
			includeTrust = true
		default:
			if filePath == "" && !strings.HasPrefix(args[i], "-") {
				filePath = args[i]
			}
		}
	}

	if filePath == "" {
		fmt.Fprintf(os.Stderr, "Error: no input file specified\n")
		fmt.Fprintf(os.Stderr, "Usage: asf export <file> [-f format] [-o output_dir]\n")
		os.Exit(ExitInvalidCmd)
	}

	data, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input file: %v\n", err)
		os.Exit(ExitExportErr)
	}

	var result AnalysisResult
	if err := json.Unmarshal(data, &result); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing analysis result: %v\n", err)
		os.Exit(ExitExportErr)
	}

	exportFormat := ExportFormat(format)
	switch exportFormat {
	case ExportJSON:
	case ExportMarkdown:
	case ExportCSV:
	case ExportPDF:
	case ExportHTML:
	case ExportJSONL:
	default:
		fmt.Fprintf(os.Stderr, "Unsupported format: %s\n", format)
		fmt.Fprintf(os.Stderr, "Supported formats: json, markdown, csv, html, pdf, jsonl\n")
		os.Exit(ExitInvalidCmd)
	}

	path, err := ExportResult(&result, exportFormat, outputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Export failed: %v\n", err)
		os.Exit(ExitExportErr)
	}

	fmt.Printf("Exported to: %s\n", path)

	if includeNarrative && result.NarrativeOutput != nil {
		narrPath := filepath.Join(outputDir, fmt.Sprintf("%s_narrative.md",
			strings.ReplaceAll(result.ArchitectureName, ".", "_")))
		ExportResult(&result, ExportNarrativeMarkdown, outputDir)
		fmt.Printf("Narrative:   %s\n", narrPath)
	}

	if includeTrust && result.TrustOutput != nil {
		ExportResult(&result, ExportTrustMarkdown, outputDir)
	}
}
