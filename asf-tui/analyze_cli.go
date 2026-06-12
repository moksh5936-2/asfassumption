package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"asf-tui/asf/analyzer"
	"asf-tui/asf/graph"
)

type cliClaim struct {
	ID                   string    `json:"id"`
	SourceDocument       string    `json:"source_document"`
	SourceLocation       string    `json:"source_location,omitempty"`
	Text                 string    `json:"text"`
	ExtractionConfidence float64   `json:"extraction_confidence"`
	CreatedAt            time.Time `json:"created_at"`
	Tags                 []string  `json:"tags"`
}

type cliOutput struct {
	Version       string            `json:"version"`
	Architecture  string            `json:"architecture"`
	Summary       cliSummary        `json:"summary"`
	Claims        []cliClaim        `json:"claims,omitempty"`
	Assumptions   []cliAssumption   `json:"assumptions"`
	Verifications []cliVerification `json:"verifications"`
	Gaps          []cliGap          `json:"gaps"`
	Graph         *graph.GraphData  `json:"graph,omitempty"`
}

type cliSummary struct {
	ClaimsFound       int `json:"claims_found"`
	Assumptions       int `json:"assumptions"`
	Verified          int `json:"verified"`
	PartiallyVerified int `json:"partially_verified"`
	Contradicted      int `json:"contradicted"`
	Unknown           int `json:"unknown"`
	CriticalGaps      int `json:"critical_gaps"`
	Critical          int `json:"critical,omitempty"`
	High              int `json:"high,omitempty"`
	Medium            int `json:"medium,omitempty"`
	Low               int `json:"low,omitempty"`
}

type cliAssumption struct {
	ID                 string   `json:"id"`
	Text               string   `json:"text"`
	AssumptionType     string   `json:"assumption_type"`
	VerificationStatus string   `json:"verification_status"`
	Confidence         float64  `json:"confidence"`
	Keywords           []string `json:"keywords"`
	Risk               string   `json:"risk,omitempty"`
}

type cliVerification struct {
	AssumptionID string      `json:"assumption_id"`
	Result       string      `json:"result"`
	Confidence   float64     `json:"confidence"`
	EvidenceUsed interface{} `json:"evidence_used"`
	Reasoning    string      `json:"reasoning"`
}

type cliGap struct {
	AssumptionID   string `json:"assumption_id"`
	Type           string `json:"type"`
	Severity       string `json:"severity"`
	Description    string `json:"description"`
	EvidenceDetail string `json:"evidence_detail"`
}

func runAnalyzeCLI(args []string) {
	graphFlag := false
	var evPaths []string
	filePath := ""

	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: asf analyze <file> [-e evidence ...] [--json] [--graph]")
			fmt.Println()
			fmt.Println("Analyze a policy document or architecture diagram for security assumptions.")
			fmt.Println()
			fmt.Println("Arguments:")
			fmt.Println("  <file>                    Policy file, architecture doc, or directory")
			fmt.Println("  -e, --evidence <path>     Evidence files/directories (CSV, JSON, YAML)")
			fmt.Println("  --json                    Output as JSON (default)")
			fmt.Println("  --graph                   Include dependency graph in JSON output")
			fmt.Println("  --help, -h                Show this help")
			os.Exit(ExitSuccess)
		}
	}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--json":
		case "--graph":
			graphFlag = true
		case "-e", "--evidence":
			i++
			if i < len(args) {
				evPaths = append(evPaths, args[i])
			}
		default:
			if filePath == "" && !strings.HasPrefix(args[i], "-") {
				filePath = args[i]
			}
		}
	}

	if filePath == "" {
		fmt.Fprintf(os.Stderr, "Error: no input file specified\n")
		fmt.Fprintf(os.Stderr, "Usage: asf analyze <file> [-e evidence ...] [--json] [--graph]\n")
		os.Exit(ExitInvalidCmd)
	}

	info, err := os.Stat(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(ExitAnalysisErr)
	}

	// Check if file is a structured architecture doc that should use the Engine pipeline
	ext := strings.ToLower(filepath.Ext(filePath))
	needsEngine := ext == ".yaml" || ext == ".yml" || ext == ".json" || ext == ".md" || ext == ".mmd" || ext == ".drawio" || ext == ".svg"

	if needsEngine && !info.IsDir() {
		// Use the full Engine pipeline for structured files
		cfg := &Config{}
		engine := NewEngine(cfg)
		progress := make(chan AnalysisProgress, 100)
		go func() {
			for range progress {
			}
		}()
		result, err := engine.RunAnalysis(filePath, "", ModeASFOnly, progress)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(ExitAnalysisErr)
		}
		out := convertAnalysisResultToCLI(result, graphFlag)
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		if err := enc.Encode(out); err != nil {
			fmt.Fprintf(os.Stderr, "Error encoding output: %v\n", err)
			os.Exit(ExitExportErr)
		}
		os.Exit(ExitSuccess)
	}

	// Fallback to raw analyzer for plain text files and directories
	an := analyzer.New()
	var docs []string
	if info.IsDir() {
		entries, err := os.ReadDir(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading directory: %s\n", err)
			os.Exit(ExitAnalysisErr)
		}
		docExts := map[string]bool{".txt": true, ".pdf": true, ".docx": true}
		for _, entry := range entries {
			if !entry.IsDir() && docExts[strings.ToLower(filepath.Ext(entry.Name()))] {
				docs = append(docs, filepath.Join(filePath, entry.Name()))
			}
		}
		if len(docs) == 0 {
			fmt.Fprintf(os.Stderr, "Error: no supported documents (.txt, .pdf, .docx) found in %s\n", filePath)
			os.Exit(ExitAnalysisErr)
		}
	} else {
		docs = []string{filePath}
	}
	var evs []string
	evidenceExts := map[string]bool{".csv": true, ".json": true, ".yaml": true, ".yml": true}
	for _, path := range evPaths {
		if info, err := os.Stat(path); err == nil {
			if info.IsDir() {
				entries, err := os.ReadDir(path)
				if err == nil {
					for _, entry := range entries {
						if !entry.IsDir() && evidenceExts[strings.ToLower(filepath.Ext(entry.Name()))] {
							evs = append(evs, filepath.Join(path, entry.Name()))
						}
					}
				}
			} else if evidenceExts[strings.ToLower(filepath.Ext(path))] {
				evs = append(evs, path)
			}
		}
	}

	ar, err := an.Analyze(docs, evs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(ExitAnalysisErr)
	}

	s := ar.Result.BuildSummary()

	out := cliOutput{
		Version:      ASFVersion,
		Architecture: filepath.Base(filePath),
	}

	out.Summary = cliSummary{
		ClaimsFound:       s.ClaimsFound,
		Assumptions:       s.Assumptions,
		Verified:          s.Verified,
		PartiallyVerified: s.PartiallyVerified,
		Contradicted:      s.Contradicted,
		Unknown:           s.Unknown,
		CriticalGaps:      s.CriticalGaps,
	}

	for _, c := range ar.Result.Claims {
		out.Claims = append(out.Claims, cliClaim{
			ID:                   c.ID,
			SourceDocument:       c.SourceDocument,
			SourceLocation:       c.SourceLocation,
			Text:                 c.Text,
			ExtractionConfidence: c.ExtractionConfidence,
			CreatedAt:            c.CreatedAt,
			Tags:                 c.Tags,
		})
	}

	for _, a := range ar.Result.Assumptions {
		out.Assumptions = append(out.Assumptions, cliAssumption{
			ID:                 a.ID,
			Text:               a.Text,
			AssumptionType:     string(a.AssumptionType),
			VerificationStatus: string(a.VerificationStatus),
			Confidence:         a.Confidence,
			Keywords:           a.Keywords,
		})
	}

	for _, v := range ar.Result.Verifications {
		out.Verifications = append(out.Verifications, cliVerification{
			AssumptionID: v.AssumptionID,
			Result:       string(v.Result),
			Confidence:   v.Confidence,
			EvidenceUsed: v.EvidenceUsed,
			Reasoning:    v.Reasoning,
		})
	}

	for _, g := range ar.Result.Gaps {
		out.Gaps = append(out.Gaps, cliGap{
			AssumptionID:   g.AssumptionID,
			Type:           string(g.Type),
			Severity:       string(g.Severity),
			Description:    g.Description,
			EvidenceDetail: g.EvidenceDetail,
		})
	}

	if graphFlag {
		out.Graph = &ar.Graph
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	if err := enc.Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "Error encoding output: %v\n", err)
		os.Exit(ExitExportErr)
	}
}

// convertAnalysisResultToCLI converts a full Engine AnalysisResult to the
// backward-compatible cliOutput format used by the CLI analyze command.
func convertAnalysisResultToCLI(result *AnalysisResult, graphFlag bool) cliOutput {
	verified := 0
	partiallyVerified := 0
	contradicted := 0
	unknown := 0
	for _, a := range result.Assumptions {
		switch a.VerificationStatus {
		case "VERIFIED":
			verified++
		case "PARTIALLY_VERIFIED":
			partiallyVerified++
		case "CONTRADICTED":
			contradicted++
		default:
			unknown++
		}
	}
	out := cliOutput{
		Version:      ASFVersion,
		Architecture: result.ArchitectureName,
		Summary: cliSummary{
			Assumptions:       result.TotalAssumptions,
			Verified:          verified,
			PartiallyVerified: partiallyVerified,
			Unknown:           unknown,
			Contradicted:      contradicted,
			CriticalGaps:      result.CriticalGaps,
			Critical:          result.CriticalCount,
			High:              result.HighCount,
			Medium:            result.MediumCount,
			Low:               result.LowCount,
		},
	}

	for _, a := range result.Assumptions {
		vStatus := a.VerificationStatus
		if vStatus == "" {
			vStatus = "UNKNOWN"
		}
		out.Assumptions = append(out.Assumptions, cliAssumption{
			ID:                 a.ID,
			Text:               a.Description,
			AssumptionType:     a.Category,
			VerificationStatus: vStatus,
			Confidence:         a.Confidence,
			Keywords:           a.Keywords,
			Risk:               string(a.Risk),
		})
	}

	// Build verifications from explicit assumptions that have verification status
	for _, a := range result.Assumptions {
		if a.VerificationStatus == "" {
			continue
		}
		out.Verifications = append(out.Verifications, cliVerification{
			AssumptionID: a.ID,
			Result:       a.VerificationStatus,
			Confidence:   a.Confidence,
			Reasoning:    a.Rationale,
		})
	}

	// Build gaps from assumptions that are not verified
	for _, a := range result.Assumptions {
		if a.VerificationStatus == "VERIFIED" || a.VerificationStatus == "PARTIALLY_VERIFIED" {
			continue
		}
		severity := "LOW"
		switch a.Risk {
		case RiskCritical:
			severity = "CRITICAL"
		case RiskHigh:
			severity = "HIGH"
		case RiskMedium:
			severity = "MEDIUM"
		}
		out.Gaps = append(out.Gaps, cliGap{
			AssumptionID:   a.ID,
			Type:           "EVIDENCE_GAP",
			Severity:       severity,
			Description:    a.Description,
			EvidenceDetail: "No matching evidence available for verification",
		})
	}

	// Add a summary claim to preserve the claims_found count
	out.Summary.ClaimsFound = len(out.Assumptions)
	out.Claims = append(out.Claims, cliClaim{
		ID:                   "clm_analysis",
		SourceDocument:       result.ArchitectureName,
		Text:                 result.Summary,
		ExtractionConfidence: 0.95,
		CreatedAt:            result.AnalysisDate,
		Tags:                 []string{"analysis", "summary"},
	})

	return out
}
