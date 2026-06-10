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
	ID                  string    `json:"id"`
	SourceDocument      string    `json:"source_document"`
	SourceLocation      string    `json:"source_location,omitempty"`
	Text                string    `json:"text"`
	ExtractionConfidence float64  `json:"extraction_confidence"`
	CreatedAt           time.Time `json:"created_at"`
	Tags                []string  `json:"tags"`
}

type cliOutput struct {
	Version       string              `json:"version"`
	Architecture  string              `json:"architecture"`
	Summary       cliSummary          `json:"summary"`
	Claims        []cliClaim          `json:"claims,omitempty"`
	Assumptions   []cliAssumption     `json:"assumptions"`
	Verifications []cliVerification   `json:"verifications"`
	Gaps          []cliGap            `json:"gaps"`
	Graph         *graph.GraphData    `json:"graph,omitempty"`
}

type cliSummary struct {
	ClaimsFound       int `json:"claims_found"`
	Assumptions       int `json:"assumptions"`
	Verified          int `json:"verified"`
	PartiallyVerified int `json:"partially_verified"`
	Contradicted      int `json:"contradicted"`
	Unknown           int `json:"unknown"`
	CriticalGaps      int `json:"critical_gaps"`
}

type cliAssumption struct {
	ID                 string   `json:"id"`
	Text               string   `json:"text"`
	AssumptionType     string   `json:"assumption_type"`
	VerificationStatus string   `json:"verification_status"`
	Confidence         float64  `json:"confidence"`
	Keywords           []string `json:"keywords"`
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

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--graph":
			graphFlag = true
		case "--json":
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
		fmt.Fprintf(os.Stderr, "Usage: asf analyze <file> [-e evidence ...] [--graph]\n")
		os.Exit(1)
	}

	info, err := os.Stat(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	an := analyzer.New()
	var docs []string
	if info.IsDir() {
		entries, err := os.ReadDir(filePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading directory: %s\n", err)
			os.Exit(1)
		}
		docExts := map[string]bool{".txt": true, ".pdf": true, ".docx": true}
		for _, entry := range entries {
			if !entry.IsDir() && docExts[strings.ToLower(filepath.Ext(entry.Name()))] {
				docs = append(docs, filepath.Join(filePath, entry.Name()))
			}
		}
		if len(docs) == 0 {
			fmt.Fprintf(os.Stderr, "Error: no supported documents (.txt, .pdf, .docx) found in %s\n", filePath)
			os.Exit(1)
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
		os.Exit(1)
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
			ID:                  c.ID,
			SourceDocument:      c.SourceDocument,
			SourceLocation:      c.SourceLocation,
			Text:                c.Text,
			ExtractionConfidence: c.ExtractionConfidence,
			CreatedAt:           c.CreatedAt,
			Tags:                c.Tags,
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
		os.Exit(1)
	}
}
