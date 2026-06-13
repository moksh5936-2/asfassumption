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
	"asf-tui/asf/narrative"
	"asf-tui/asf/trust"
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
	Version            string                `json:"version"`
	Architecture       string                `json:"architecture"`
	Summary            cliSummary            `json:"summary"`
	Claims             []cliClaim            `json:"claims,omitempty"`
	Assumptions        []cliAssumption       `json:"assumptions"`
	Verifications      []cliVerification     `json:"verifications"`
	Gaps               []cliGap              `json:"gaps"`
	Graph              *graph.GraphData      `json:"graph,omitempty"`
	Contradictions     []cliContradiction    `json:"contradictions,omitempty"`
	TrustBoundaries    []cliTrustBoundary    `json:"trust_boundaries,omitempty"`
	Domain             string                `json:"domain,omitempty"`
	CIEContradictions  []cliCIEContradiction `json:"cie_contradictions,omitempty"`
	CIESummary         string                `json:"cie_summary,omitempty"`
	TBIZones           []cliTBIZone          `json:"tbi_zones,omitempty"`
	TBIBoundaries      []cliTBIBoundary      `json:"tbi_boundaries,omitempty"`
	TBIWeaknesses      []cliTBIWeakness      `json:"tbi_weaknesses,omitempty"`
	TBISummary         string                `json:"tbi_summary,omitempty"`
	Threats            []cliThreat           `json:"threats,omitempty"`
	ThreatClusters     []cliThreatCluster    `json:"threat_clusters,omitempty"`
	ThreatModelSummary cliThreatModelSummary `json:"threat_model_summary,omitempty"`
	AttackPaths        []cliAttackPath       `json:"attack_paths,omitempty"`
	ThreatChains       []cliThreatChain      `json:"threat_chains,omitempty"`
	AttackPathSummary  cliAttackPathSummary  `json:"attack_path_summary,omitempty"`
	SDRIControls       []cliSDRIControl      `json:"sdri_controls,omitempty"`
	SDRIFindings       []cliSDRIFinding      `json:"sdri_findings,omitempty"`
	SDRIAW             []cliSDRIWeakness     `json:"sdri_weaknesses,omitempty"`
	SDRIRemediations   []cliSDRIRemediation  `json:"sdri_remediations,omitempty"`
	SDRICoverage       []cliSDRICoverage     `json:"sdri_coverage,omitempty"`
	SDRIDashboard      map[string]float64    `json:"sdri_dashboard,omitempty"`
	SDRICompliance     []cliSDRICompliance   `json:"sdri_compliance,omitempty"`
	SDRISummary        string                `json:"sdri_summary,omitempty"`

	// CIARE fields
	CIAREFrameworkCoverages   []cliCIAREFrameworkCoverage   `json:"ciare_framework_coverages,omitempty"`
	CIAREAuditReadiness       []cliCIAREAuditReadiness      `json:"ciare_audit_readiness,omitempty"`
	CIAREEvidenceRequirements []cliCIAREEvidenceRequirement `json:"ciare_evidence_requirements,omitempty"`
	CIAREMissingEvidences     []cliCIAREMissingEvidence     `json:"ciare_missing_evidences,omitempty"`
	CIAREAuditorQuestions     []cliCIAREAuditorQuestion     `json:"ciare_auditor_questions,omitempty"`
	CIAREComplianceGaps       []cliCIAREComplianceGap       `json:"ciare_compliance_gaps,omitempty"`
	CIAREControlMaturities    []cliCIAREControlMaturity     `json:"ciare_control_maturities,omitempty"`
	CIAREComplianceNarratives []cliCIAREComplianceNarrative `json:"ciare_compliance_narratives,omitempty"`
	CIAREAuditPackage         *cliCIAREAuditPackage         `json:"ciare_audit_package,omitempty"`
	CIAREComplianceDashboard  *cliCIAREComplianceDashboard  `json:"ciare_compliance_dashboard,omitempty"`
	CIAREProcurementQuestions []cliCIAREProcurementQuestion `json:"ciare_procurement_questions,omitempty"`

	// DKPI fields
	DKPIDomain          string           `json:"dkpi_domain,omitempty"`
	DKPIConfidence      float64          `json:"dkpi_confidence,omitempty"`
	DKPIRecommendations []string         `json:"dkpi_recommendations,omitempty"`
	DKPIThreats         []cliThreat      `json:"dkpi_threats,omitempty"`
	DKPIControls        []cliSDRIControl `json:"dkpi_controls,omitempty"`
	DKPISummary         string           `json:"dkpi_summary,omitempty"`

	// ERN fields
	ERNExecutiveRisks  []cliERNExecutiveRisk `json:"ern_executive_risks,omitempty"`
	ERNBoardSummary    string                `json:"ern_board_summary,omitempty"`
	ERNExposure        string                `json:"ern_exposure,omitempty"`
	ERNTopRisks        []string              `json:"ern_top_risks,omitempty"`
	ERNRemediation     []string              `json:"ern_remediation,omitempty"`
	ERNInvestmentAreas []string              `json:"ern_investment_areas,omitempty"`
	ERNReportType      string                `json:"ern_report_type,omitempty"`
	ERNBoardReport     string                `json:"ern_board_report,omitempty"`
	ERNExecutiveReport string                `json:"ern_executive_report,omitempty"`
	ERNTechnicalReport string                `json:"ern_technical_report,omitempty"`

	// Trust chain fields
	TrustChains          []cliTrustChain                `json:"trust_chains,omitempty"`
	FailureCascades      []cliFailureCascade            `json:"failure_cascades,omitempty"`
	CriticalAssumptions  []cliCriticalAssumption        `json:"critical_assumptions,omitempty"`
	SinglePointsOfTrust  []cliSinglePointOfTrustFailure `json:"single_points_of_trust_failure,omitempty"`
	TrustCollapseResults []cliTrustCollapseResult       `json:"trust_collapse_results,omitempty"`
	TrustChainSummary    string                         `json:"trust_chain_summary,omitempty"`

	// Security Architect Narrative Engine (SANE) fields
	NarrativeOutput *narrative.NarrativeOutput `json:"narrative_output,omitempty"`
}

type cliERNExecutiveRisk struct {
	ID                 string   `json:"id"`
	Title              string   `json:"title"`
	Priority           string   `json:"priority"`
	Severity           string   `json:"severity"`
	BusinessImpact     string   `json:"business_impact"`
	RecommendedActions []string `json:"recommended_actions,omitempty"`
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

type cliContradiction struct {
	ID                  string   `json:"id"`
	Severity            string   `json:"severity"`
	Description         string   `json:"description"`
	Explanation         string   `json:"explanation"`
	AffectedAssumptions []string `json:"affected_assumptions"`
	RuleName            string   `json:"rule_name"`
}

type cliTrustBoundary struct {
	Type        string   `json:"type"`
	Components  []string `json:"components"`
	RiskLevel   string   `json:"risk_level"`
	Description string   `json:"description"`
}

type cliTBIZone struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Sensitivity string   `json:"sensitivity"`
	Components  []string `json:"components"`
	Description string   `json:"description"`
}

type cliTBIBoundary struct {
	ID                  string   `json:"id"`
	SourceZone          string   `json:"source_zone"`
	DestinationZone     string   `json:"destination_zone"`
	SourceZoneType      string   `json:"source_zone_type"`
	DestinationZoneType string   `json:"destination_zone_type"`
	CrossingType        string   `json:"crossing_type"`
	Risk                string   `json:"risk"`
	Confidence          float64  `json:"confidence"`
	RequiredControls    []string `json:"required_controls"`
	RequiredAssumptions []string `json:"required_assumptions"`
	Threats             []string `json:"threats"`
	MissingControls     []string `json:"missing_controls,omitempty"`
	MissingAssumptions  []string `json:"missing_assumptions,omitempty"`
	Reasoning           string   `json:"reasoning"`
	Recommendations     []string `json:"recommendations,omitempty"`
	ComplianceMappings  []string `json:"compliance_mappings,omitempty"`
}

type cliTBIWeakness struct {
	ID              string   `json:"id"`
	BoundaryID      string   `json:"boundary_id"`
	Type            string   `json:"type"`
	Severity        string   `json:"severity"`
	Description     string   `json:"description"`
	Reasoning       string   `json:"reasoning"`
	Recommendations []string `json:"recommendations,omitempty"`
}

type cliThreat struct {
	ID                 string   `json:"id"`
	Name               string   `json:"name"`
	Category           string   `json:"category"`
	Severity           string   `json:"severity"`
	Likelihood         float64  `json:"likelihood"`
	Impact             float64  `json:"impact"`
	RiskScore          float64  `json:"risk_score"`
	Confidence         float64  `json:"confidence"`
	Description        string   `json:"description"`
	AffectedAssets     []string `json:"affected_assets,omitempty"`
	AffectedComponents []string `json:"affected_components,omitempty"`
	AffectedBoundaries []string `json:"affected_boundaries,omitempty"`
	AffectedData       []string `json:"affected_data,omitempty"`
	Assumptions        []string `json:"assumptions,omitempty"`
	Controls           []string `json:"controls,omitempty"`
	STRIDECategories   []string `json:"stride_categories,omitempty"`
	Reasoning          string   `json:"reasoning"`
	Recommendations    []string `json:"recommendations,omitempty"`
	PreventiveControls []string `json:"preventive_controls,omitempty"`
	DetectiveControls  []string `json:"detective_controls,omitempty"`
	CorrectiveControls []string `json:"corrective_controls,omitempty"`
}

type cliThreatCluster struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Category        string   `json:"category"`
	Threats         []string `json:"threats"`
	RiskScore       float64  `json:"risk_score"`
	AffectedAssets  []string `json:"affected_assets,omitempty"`
	Recommendations []string `json:"recommendations,omitempty"`
}

type cliThreatModelSummary struct {
	TotalThreats       int            `json:"total_threats"`
	CriticalCount      int            `json:"critical_count"`
	HighCount          int            `json:"high_count"`
	MediumCount        int            `json:"medium_count"`
	LowCount           int            `json:"low_count"`
	ClusterCount       int            `json:"cluster_count"`
	STRIDEDistribution map[string]int `json:"stride_distribution"`
	TopThreats         []string       `json:"top_threats"`
	SummaryText        string         `json:"summary_text"`
}

func runAnalyzeCLI(args []string) {
	graphFlag := false
	reportType := ""
	var evPaths []string
	filePath := ""

	for _, a := range args {
		if a == "--help" || a == "-h" {
			fmt.Println("Usage: asf analyze <file> [-e evidence ...] [--json] [--graph] [--report-type board|executive|technical|architect-narrative|executive-summary|technical-summary]")
			fmt.Println()
			fmt.Println("Analyze a policy document or architecture diagram for security assumptions.")
			fmt.Println()
			fmt.Println("Arguments:")
			fmt.Println("  <file>                    Policy file, architecture doc, or directory")
			fmt.Println("  -e, --evidence <path>     Evidence files/directories (CSV, JSON, YAML)")
			fmt.Println("  --json                    Output as JSON (default)")
			fmt.Println("  --graph                   Include dependency graph in JSON output")
			fmt.Println("  --report-type <type>      Filter output to specific report pack (board, executive, technical, architect-narrative, executive-summary, technical-summary)")
			fmt.Println("  --help, -h                Show this help")
			os.Exit(ExitSuccess)
		}
	}

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--json":
		case "--graph":
			graphFlag = true
		case "--report-type":
			i++
			if i < len(args) {
				reportType = args[i]
			}
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
		out := convertAnalysisResultToCLI(result, graphFlag, reportType)
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

func convertDependencyTypes(dts []trust.DependencyType) []string {
	out := make([]string, len(dts))
	for i, dt := range dts {
		out[i] = string(dt)
	}
	return out
}

// convertAnalysisResultToCLI converts a full Engine AnalysisResult to the
// backward-compatible cliOutput format used by the CLI analyze command.
func convertAnalysisResultToCLI(result *AnalysisResult, graphFlag bool, reportType string) cliOutput {
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

	// Add contradictions
	for _, c := range result.Contradictions {
		out.Contradictions = append(out.Contradictions, cliContradiction{
			ID:                  c.ID,
			Severity:            string(c.Severity),
			Description:         c.Description,
			Explanation:         c.Explanation,
			AffectedAssumptions: c.AffectedAssumptions,
			RuleName:            c.RuleName,
		})
	}

	// Add trust boundaries
	for _, b := range result.TrustBoundaries {
		out.TrustBoundaries = append(out.TrustBoundaries, cliTrustBoundary{
			Type:        b.Type,
			Components:  b.Components,
			RiskLevel:   string(b.RiskLevel),
			Description: b.Description,
		})
	}

	// Add domain
	out.Domain = result.Domain

	// Add CIE contradictions
	for _, c := range result.CIEContradictions {
		out.CIEContradictions = append(out.CIEContradictions, cliCIEContradiction{
			ID:                      c.ID,
			Type:                    c.Type,
			Severity:                string(c.Severity),
			Confidence:              c.Confidence,
			Summary:                 c.Summary,
			Description:             c.Description,
			StatementA:              cliCIEStatement{ID: c.StatementA.ID, Source: c.StatementA.Source, OriginalText: c.StatementA.OriginalText, Category: c.StatementA.Category, Confidence: c.StatementA.Confidence},
			StatementB:              cliCIEStatement{ID: c.StatementB.ID, Source: c.StatementB.Source, OriginalText: c.StatementB.OriginalText, Category: c.StatementB.Category, Confidence: c.StatementB.Confidence},
			AffectedAssets:          c.AffectedAssets,
			AffectedComponents:      c.AffectedComponents,
			AffectedControls:        c.AffectedControls,
			AffectedTrustBoundaries: c.AffectedTrustBoundaries,
			Reasoning:               c.Reasoning,
			Evidence:                c.Evidence,
			Recommendations:         c.Recommendations,
		})
	}

	// Add CIE summary
	out.CIESummary = result.CIESummary

	// Add TBI zones
	for _, z := range result.TBIZones {
		out.TBIZones = append(out.TBIZones, cliTBIZone{
			ID:          z.ID,
			Name:        z.Name,
			Type:        z.Type,
			Sensitivity: z.Sensitivity,
			Components:  z.Components,
			Description: z.Description,
		})
	}

	// Add TBI boundaries
	for _, b := range result.TBIBoundaries {
		out.TBIBoundaries = append(out.TBIBoundaries, cliTBIBoundary{
			ID:                  b.ID,
			SourceZone:          b.SourceZone,
			DestinationZone:     b.DestinationZone,
			SourceZoneType:      b.SourceZoneType,
			DestinationZoneType: b.DestinationZoneType,
			CrossingType:        b.CrossingType,
			Risk:                string(b.Risk),
			Confidence:          b.Confidence,
			RequiredControls:    b.RequiredControls,
			RequiredAssumptions: b.RequiredAssumptions,
			Threats:             b.Threats,
			MissingControls:     b.MissingControls,
			MissingAssumptions:  b.MissingAssumptions,
			Reasoning:           b.Reasoning,
			Recommendations:     b.Recommendations,
			ComplianceMappings:  b.ComplianceMappings,
		})
	}

	// Add TBI weaknesses
	for _, w := range result.TBIWeaknesses {
		out.TBIWeaknesses = append(out.TBIWeaknesses, cliTBIWeakness{
			ID:              w.ID,
			BoundaryID:      w.BoundaryID,
			Type:            w.Type,
			Severity:        string(w.Severity),
			Description:     w.Description,
			Reasoning:       w.Reasoning,
			Recommendations: w.Recommendations,
		})
	}

	// Add TBI summary
	out.TBISummary = result.TBISummary

	// Add trust chain data
	if result.TrustOutput != nil {
		for _, tc := range result.TrustOutput.TrustChains {
			out.TrustChains = append(out.TrustChains, cliTrustChain{
				ID:              tc.ID,
				Nodes:           tc.Nodes,
				Length:          tc.Length,
				Confidence:      tc.Confidence,
				Risk:            tc.Risk,
				DependencyCount: tc.DependencyCount,
				RootNode:        tc.RootNode,
				LeafNode:        tc.LeafNode,
			})
		}
		for _, fc := range result.TrustOutput.FailureCascades {
			steps := make([]cliCascadeResult, len(fc.Steps))
			for i, s := range fc.Steps {
				steps[i] = cliCascadeResult{
					Step:             s.Step,
					AssumptionID:     s.AssumptionID,
					AssumptionText:   s.AssumptionText,
					Severity:         s.Severity,
					AffectedAssets:   s.AffectedAssets,
					AffectedControls: s.AffectedControls,
					Reason:           s.Reason,
				}
			}
			out.FailureCascades = append(out.FailureCascades, cliFailureCascade{
				RootAssumptionID:   fc.RootAssumptionID,
				RootAssumptionText: fc.RootAssumptionText,
				Steps:              steps,
				TotalAffected:      fc.TotalAffected,
				Severity:           fc.Severity,
				MaxDepth:           fc.MaxDepth,
			})
		}
		for _, ca := range result.TrustOutput.CriticalAssumptions {
			out.CriticalAssumptions = append(out.CriticalAssumptions, cliCriticalAssumption{
				AssumptionID:    ca.AssumptionID,
				AssumptionText:  ca.AssumptionText,
				Centrality:      ca.Centrality,
				SupportCount:    ca.SupportCount,
				FailureRadius:   ca.FailureRadius,
				TrustRadius:     ca.TrustRadius,
				Risk:            ca.Risk,
				Score:           ca.Score,
				DependencyTypes: convertDependencyTypes(ca.DependencyTypes),
			})
		}
		for _, sp := range result.TrustOutput.SinglePointsOfTrust {
			out.SinglePointsOfTrust = append(out.SinglePointsOfTrust, cliSinglePointOfTrustFailure{
				NodeID:          sp.NodeID,
				AssumptionText:  sp.AssumptionText,
				DependentsCount: sp.DependentsCount,
				DependentNodes:  sp.DependentNodes,
				DependencyTypes: convertDependencyTypes(sp.DependencyTypes),
				Recommendation:  sp.Recommendation,
			})
		}
		for _, tc := range result.TrustOutput.TrustCollapseResults {
			out.TrustCollapseResults = append(out.TrustCollapseResults, cliTrustCollapseResult{
				FailedAssumptionID:   tc.FailedAssumptionID,
				FailedAssumptionText: tc.FailedAssumptionText,
				AssumptionsLost:      tc.AssumptionsLost,
				ControlsLost:         tc.ControlsLost,
				AssetsExposed:        tc.AssetsExposed,
				RiskIncrease:         tc.RiskIncrease,
				RiskScoreBefore:      tc.RiskScoreBefore,
				RiskScoreAfter:       tc.RiskScoreAfter,
				AffectedComponents:   tc.AffectedComponents,
			})
		}
		out.TrustChainSummary = fmt.Sprintf("%d trust chains, %d failure cascades, %d critical assumptions, %d single points of trust failure, %d collapse simulations",
			len(result.TrustOutput.TrustChains),
			len(result.TrustOutput.FailureCascades),
			len(result.TrustOutput.CriticalAssumptions),
			len(result.TrustOutput.SinglePointsOfTrust),
			len(result.TrustOutput.TrustCollapseResults))
	}

	// Add TMI threats
	for _, t := range result.Threats {
		out.Threats = append(out.Threats, cliThreat{
			ID:                 t.ID,
			Name:               t.Name,
			Category:           t.Category,
			Severity:           string(t.Severity),
			Likelihood:         t.Likelihood,
			Impact:             t.Impact,
			RiskScore:          t.RiskScore,
			Confidence:         t.Confidence,
			Description:        t.Description,
			AffectedAssets:     t.AffectedAssets,
			AffectedComponents: t.AffectedComponents,
			AffectedBoundaries: t.AffectedBoundaries,
			AffectedData:       t.AffectedData,
			Assumptions:        t.Assumptions,
			Controls:           t.Controls,
			STRIDECategories:   t.STRIDECategories,
			Reasoning:          t.Reasoning,
			Recommendations:    t.Recommendations,
			PreventiveControls: t.PreventiveControls,
			DetectiveControls:  t.DetectiveControls,
			CorrectiveControls: t.CorrectiveControls,
		})
	}

	// Add TMI threat clusters
	for _, c := range result.ThreatClusters {
		out.ThreatClusters = append(out.ThreatClusters, cliThreatCluster{
			ID:              c.ID,
			Name:            c.Name,
			Category:        c.Category,
			Threats:         c.Threats,
			RiskScore:       c.RiskScore,
			AffectedAssets:  c.AffectedAssets,
			Recommendations: c.Recommendations,
		})
	}

	// Add TMI threat model summary
	out.ThreatModelSummary = cliThreatModelSummary{
		TotalThreats:       result.ThreatModelSummary.TotalThreats,
		CriticalCount:      result.ThreatModelSummary.CriticalCount,
		HighCount:          result.ThreatModelSummary.HighCount,
		MediumCount:        result.ThreatModelSummary.MediumCount,
		LowCount:           result.ThreatModelSummary.LowCount,
		ClusterCount:       result.ThreatModelSummary.ClusterCount,
		STRIDEDistribution: result.ThreatModelSummary.STRIDEDistribution,
		TopThreats:         result.ThreatModelSummary.TopThreats,
		SummaryText:        result.ThreatModelSummary.SummaryText,
	}

	// Add APD attack paths
	for _, p := range result.AttackPaths {
		steps := make([]cliAttackStep, len(p.AttackSteps))
		for i, s := range p.AttackSteps {
			steps[i] = cliAttackStep{
				SequenceNumber:     s.SequenceNumber,
				SourceComponent:    s.SourceComponent,
				TargetComponent:    s.TargetComponent,
				Action:             s.Action,
				Threat:             s.Threat,
				RequiredAssumption: s.RequiredAssumption,
				ControlBypassed:    s.ControlBypassed,
				Reasoning:          s.Reasoning,
				STRIDECategory:     s.STRIDECategory,
			}
		}
		out.AttackPaths = append(out.AttackPaths, cliAttackPath{
			ID:                  p.ID,
			Name:                p.Name,
			Description:         p.Description,
			EntryPoint:          p.EntryPoint,
			TargetAsset:         p.TargetAsset,
			AttackSteps:         steps,
			Likelihood:          p.Likelihood,
			Impact:              p.Impact,
			RiskScore:           p.RiskScore,
			DetectionDifficulty: p.DetectionDifficulty,
			BusinessImpact:      p.BusinessImpact,
			Recommendations:     p.Recommendations,
			KillChainPhases:     p.KillChainPhases,
			MITREATTACK:         p.MITREATTACK,
		})
	}

	// Add APD threat chains
	for _, c := range result.ThreatChains {
		out.ThreatChains = append(out.ThreatChains, cliThreatChain{
			ID:        c.ID,
			Threats:   c.Threats,
			Path:      c.Path,
			RiskScore: c.RiskScore,
			Reasoning: c.Reasoning,
		})
	}

	// Add APD attack path summary
	out.AttackPathSummary = cliAttackPathSummary{
		TotalAttackPaths:  result.AttackPathSummary.TotalAttackPaths,
		CriticalCount:     result.AttackPathSummary.CriticalCount,
		HighCount:         result.AttackPathSummary.HighCount,
		MediumCount:       result.AttackPathSummary.MediumCount,
		LowCount:          result.AttackPathSummary.LowCount,
		ThreatChainCount:  result.AttackPathSummary.ThreatChainCount,
		TopAttackPaths:    result.AttackPathSummary.TopAttackPaths,
		KillChainCoverage: result.AttackPathSummary.KillChainCoverage,
		MITRECoverage:     result.AttackPathSummary.MITRECoverage,
		SummaryText:       result.AttackPathSummary.SummaryText,
	}

	// Add SDRI data
	for _, c := range result.SDRIControls {
		out.SDRIControls = append(out.SDRIControls, cliSDRIControl{
			ID:          c.ID,
			Name:        c.Name,
			Category:    c.Category,
			Description: c.Description,
			ControlType: c.ControlType,
			Preventive:  c.Preventive,
			Detective:   c.Detective,
			Corrective:  c.Corrective,
			Strength:    c.Strength,
			Coverage:    c.Coverage,
			Status:      c.Status,
		})
	}
	for _, f := range result.SDRIDesignFindings {
		out.SDRIFindings = append(out.SDRIFindings, cliSDRIFinding{
			ID:                 f.ID,
			Title:              f.Title,
			Description:        f.Description,
			Severity:           f.Severity,
			Category:           f.Category,
			AffectedComponents: f.AffectedComponents,
			AffectedControls:   f.AffectedControls,
			BusinessImpact:     f.BusinessImpact,
			Recommendation:     f.Recommendation,
		})
	}
	for _, w := range result.SDRIAchitecturalWeaknesses {
		out.SDRIAW = append(out.SDRIAW, cliSDRIWeakness{
			ID:             w.ID,
			Pattern:        w.Pattern,
			Description:    w.Description,
			Severity:       w.Severity,
			Components:     w.Components,
			Impact:         w.Impact,
			Recommendation: w.Recommendation,
		})
	}
	for _, r := range result.SDRIRemediations {
		out.SDRIRemediations = append(out.SDRIRemediations, cliSDRIRemediation{
			ID:             r.ID,
			Priority:       r.Priority,
			Description:    r.Description,
			RiskScore:      r.RiskScore,
			BusinessImpact: r.BusinessImpact,
			Effort:         r.Effort,
			Category:       r.Category,
			Recommendation: r.Recommendation,
		})
	}
	for _, c := range result.SDRICoverageByCategory {
		out.SDRICoverage = append(out.SDRICoverage, cliSDRICoverage{
			Category: c.Category,
			Expected: c.Expected,
			Observed: c.Observed,
			Coverage: c.Coverage,
			Level:    c.Level,
		})
	}
	out.SDRIDashboard = result.SDRICoverageDashboard
	for _, m := range result.SDRIComplianceAlignments {
		out.SDRICompliance = append(out.SDRICompliance, cliSDRICompliance{
			Framework: m.Framework,
			Coverage:  m.Coverage,
			Controls:  m.Controls,
			Status:    m.Status,
		})
	}
	out.SDRISummary = result.SDRISummary

	// Add CIARE data
	for _, c := range result.CIAREFrameworkCoverages {
		out.CIAREFrameworkCoverages = append(out.CIAREFrameworkCoverages, cliCIAREFrameworkCoverage{
			Framework:        c.Framework,
			Required:         c.Required,
			Observed:         c.Observed,
			Missing:          c.Missing,
			CoveragePct:      c.CoveragePct,
			Status:           c.Status,
			ObservedControls: c.ObservedControls,
			MissingControls:  c.MissingControls,
		})
	}
	for _, a := range result.CIAREAuditReadiness {
		out.CIAREAuditReadiness = append(out.CIAREAuditReadiness, cliCIAREAuditReadiness{
			Framework:       a.Framework,
			ReadinessScore:  a.ReadinessScore,
			Status:          a.Status,
			ControlCoverage: a.ControlCoverage,
			EvidenceScore:   a.EvidenceScore,
			ThreatExposure:  a.ThreatExposure,
			FindingsPenalty: a.FindingsPenalty,
			Factors:         a.Factors,
		})
	}
	for _, e := range result.CIAREEvidenceRequirements {
		out.CIAREEvidenceRequirements = append(out.CIAREEvidenceRequirements, cliCIAREEvidenceRequirement{
			Framework: e.Framework,
			Control:   e.Control,
			Evidence:  e.Evidence,
		})
	}
	for _, m := range result.CIAREMissingEvidences {
		out.CIAREMissingEvidences = append(out.CIAREMissingEvidences, cliCIAREMissingEvidence{
			Framework: m.Framework,
			Control:   m.Control,
			Evidences: m.Evidences,
		})
	}
	for _, q := range result.CIAREAuditorQuestions {
		out.CIAREAuditorQuestions = append(out.CIAREAuditorQuestions, cliCIAREAuditorQuestion{
			Framework: q.Framework,
			Control:   q.Control,
			Question:  q.Question,
		})
	}
	for _, g := range result.CIAREComplianceGaps {
		out.CIAREComplianceGaps = append(out.CIAREComplianceGaps, cliCIAREComplianceGap{
			ID:          g.ID,
			Framework:   g.Framework,
			Requirement: g.Requirement,
			Observed:    g.Observed,
			Missing:     g.Missing,
			Risk:        g.Risk,
		})
	}
	for _, m := range result.CIAREControlMaturities {
		out.CIAREControlMaturities = append(out.CIAREControlMaturities, cliCIAREControlMaturity{
			Domain:   m.Domain,
			Level:    m.Level,
			Label:    m.Label,
			Coverage: m.Coverage,
		})
	}
	for _, n := range result.CIAREComplianceNarratives {
		out.CIAREComplianceNarratives = append(out.CIAREComplianceNarratives, cliCIAREComplianceNarrative{
			Framework: n.Framework,
			Narrative: n.Narrative,
		})
	}
	if result.CIAREAuditPackage != nil {
		auditPkg := &cliCIAREAuditPackage{
			ExecutiveSummary: result.CIAREAuditPackage.ExecutiveSummary,
		}
		for _, c := range result.CIAREAuditPackage.FrameworkCoverages {
			auditPkg.FrameworkCoverages = append(auditPkg.FrameworkCoverages, cliCIAREFrameworkCoverage{
				Framework: c.Framework, Required: c.Required, Observed: c.Observed,
				Missing: c.Missing, CoveragePct: c.CoveragePct, Status: c.Status,
			})
		}
		for _, c := range result.CIAREAuditPackage.ControlInventory {
			auditPkg.ControlInventory = append(auditPkg.ControlInventory, cliSDRIControl{
				ID: c.ID, Name: c.Name, Category: c.Category, Description: c.Description,
				ControlType: c.ControlType, Preventive: c.Preventive, Detective: c.Detective,
				Corrective: c.Corrective, Strength: c.Strength, Coverage: c.Coverage, Status: c.Status,
			})
		}
		for _, g := range result.CIAREAuditPackage.MissingControls {
			auditPkg.MissingControls = append(auditPkg.MissingControls, cliCIAREComplianceGap{
				ID: g.ID, Framework: g.Framework, Requirement: g.Requirement,
				Observed: g.Observed, Missing: g.Missing, Risk: g.Risk,
			})
		}
		for _, e := range result.CIAREAuditPackage.EvidenceRequirements {
			auditPkg.EvidenceRequirements = append(auditPkg.EvidenceRequirements, cliCIAREEvidenceRequirement{
				Framework: e.Framework, Control: e.Control, Evidence: e.Evidence,
			})
		}
		for _, q := range result.CIAREAuditPackage.AuditorQuestions {
			auditPkg.AuditorQuestions = append(auditPkg.AuditorQuestions, cliCIAREAuditorQuestion{
				Framework: q.Framework, Control: q.Control, Question: q.Question,
			})
		}
		out.CIAREAuditPackage = auditPkg
	}
	if result.CIAREComplianceDashboard != nil {
		dash := &cliCIAREComplianceDashboard{
			FrameworkCoverages: result.CIAREComplianceDashboard.FrameworkCoverages,
			TopRisks:           result.CIAREComplianceDashboard.TopRisks,
		}
		for _, g := range result.CIAREComplianceDashboard.TopGaps {
			dash.TopGaps = append(dash.TopGaps, cliCIAREComplianceGap{
				ID: g.ID, Framework: g.Framework, Requirement: g.Requirement,
				Observed: g.Observed, Missing: g.Missing, Risk: g.Risk,
			})
		}
		for _, m := range result.CIAREComplianceDashboard.TopMissingEvidence {
			dash.TopMissingEvidence = append(dash.TopMissingEvidence, cliCIAREMissingEvidence{
				Framework: m.Framework, Control: m.Control, Evidences: m.Evidences,
			})
		}
		out.CIAREComplianceDashboard = dash
	}
	for _, q := range result.CIAREProcurementQuestions {
		out.CIAREProcurementQuestions = append(out.CIAREProcurementQuestions, cliCIAREProcurementQuestion{
			Category: q.Category,
			Question: q.Question,
		})
	}

	// Add DKPI data
	out.DKPIDomain = result.DKPI.DomainResult.PrimaryDomain
	out.DKPIConfidence = result.DKPI.DomainResult.Confidence
	out.DKPIRecommendations = result.DKPI.Recommendations
	out.DKPISummary = result.DKPI.Summary
	for _, t := range result.DKPI.InjectedThreats {
		out.DKPIThreats = append(out.DKPIThreats, cliThreat{
			ID:                 t.ID,
			Name:               t.Name,
			Category:           t.Category,
			Severity:           string(t.Severity),
			Likelihood:         t.Likelihood,
			Impact:             t.Impact,
			RiskScore:          t.RiskScore,
			Confidence:         t.Confidence,
			Description:        t.Description,
			AffectedAssets:     t.AffectedAssets,
			AffectedComponents: t.AffectedComponents,
			AffectedBoundaries: t.AffectedBoundaries,
			STRIDECategories:   t.STRIDECategories,
			Reasoning:          t.Reasoning,
			Recommendations:    t.Recommendations,
		})
	}
	for _, c := range result.DKPI.DomainControls {
		out.DKPIControls = append(out.DKPIControls, cliSDRIControl{
			ID: c.ID, Name: c.Name, Category: c.Category,
			Description: c.Description, ControlType: c.ControlType,
			Preventive: c.Preventive, Detective: c.Detective, Corrective: c.Corrective,
			Strength: c.Strength, Coverage: c.Coverage, Status: c.Status,
		})
	}

	// Add ERN data
	for _, risk := range result.ERN.ExecutiveRisks {
		out.ERNExecutiveRisks = append(out.ERNExecutiveRisks, cliERNExecutiveRisk{
			ID:                 risk.ID,
			Title:              risk.Title,
			Priority:           risk.Priority,
			Severity:           risk.Severity,
			BusinessImpact:     risk.BusinessImpact,
			RecommendedActions: risk.RecommendedActions,
		})
	}
	if result.ERN.BoardSummary.Summary != "" {
		out.ERNBoardSummary = result.ERN.BoardSummary.Summary
	}
	out.ERNExposure = result.ERN.FinancialExposure.Level
	out.ERNTopRisks = result.ERN.CISOBriefing.TopRisks
	for _, item := range result.ERN.RemediationRoadmap.Phase30 {
		out.ERNRemediation = append(out.ERNRemediation, item.Action)
	}
	for _, item := range result.ERN.RemediationRoadmap.Phase90 {
		out.ERNRemediation = append(out.ERNRemediation, item.Action)
	}
	for _, ii := range result.ERN.InvestmentInsights {
		out.ERNInvestmentAreas = append(out.ERNInvestmentAreas, ii.Area)
	}
	if reportType != "" {
		out.ERNReportType = reportType
	}
	if result.ERN.ReportPacks.BoardReport != "" {
		out.ERNBoardReport = result.ERN.ReportPacks.BoardReport
	}
	if result.ERN.ReportPacks.ExecutiveReport != "" {
		out.ERNExecutiveReport = result.ERN.ReportPacks.ExecutiveReport
	}
	if result.ERN.ReportPacks.TechnicalReport != "" {
		out.ERNTechnicalReport = result.ERN.ReportPacks.TechnicalReport
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

	// Include narrative output if available
	if result.NarrativeOutput != nil {
		out.NarrativeOutput = result.NarrativeOutput
	}

	return out
}

// cliCIEContradiction represents a CIE contradiction in CLI output.
type cliCIEContradiction struct {
	ID                      string          `json:"id"`
	Type                    string          `json:"type"`
	Severity                string          `json:"severity"`
	Confidence              float64         `json:"confidence"`
	Summary                 string          `json:"summary"`
	Description             string          `json:"description"`
	StatementA              cliCIEStatement `json:"statement_a"`
	StatementB              cliCIEStatement `json:"statement_b"`
	AffectedAssets          []string        `json:"affected_assets,omitempty"`
	AffectedComponents      []string        `json:"affected_components,omitempty"`
	AffectedControls        []string        `json:"affected_controls,omitempty"`
	AffectedTrustBoundaries []string        `json:"affected_trust_boundaries,omitempty"`
	Reasoning               string          `json:"reasoning"`
	Evidence                []string        `json:"evidence,omitempty"`
	Recommendations         []string        `json:"recommendations,omitempty"`
}

// cliCIEStatement represents a claim statement in CLI output.
type cliCIEStatement struct {
	ID           string  `json:"id"`
	Source       string  `json:"source"`
	OriginalText string  `json:"original_text"`
	Category     string  `json:"category"`
	Confidence   float64 `json:"confidence"`
}

// Trust chain CLI types
type cliTrustChain struct {
	ID              string   `json:"id"`
	Nodes           []string `json:"nodes"`
	Length          int      `json:"length"`
	Confidence      float64  `json:"confidence"`
	Risk            string   `json:"risk"`
	DependencyCount int      `json:"dependency_count"`
	RootNode        string   `json:"root_node"`
	LeafNode        string   `json:"leaf_node"`
}

type cliCascadeResult struct {
	Step             int      `json:"step"`
	AssumptionID     string   `json:"assumption_id"`
	AssumptionText   string   `json:"assumption_text"`
	Severity         string   `json:"severity"`
	AffectedAssets   []string `json:"affected_assets,omitempty"`
	AffectedControls []string `json:"affected_controls,omitempty"`
	Reason           string   `json:"reason"`
}

type cliFailureCascade struct {
	RootAssumptionID   string             `json:"root_assumption_id"`
	RootAssumptionText string             `json:"root_assumption_text"`
	Steps              []cliCascadeResult `json:"steps"`
	TotalAffected      int                `json:"total_affected"`
	Severity           string             `json:"severity"`
	MaxDepth           int                `json:"max_depth"`
}

type cliCriticalAssumption struct {
	AssumptionID    string   `json:"assumption_id"`
	AssumptionText  string   `json:"assumption_text"`
	Centrality      float64  `json:"centrality"`
	SupportCount    int      `json:"support_count"`
	FailureRadius   int      `json:"failure_radius"`
	TrustRadius     int      `json:"trust_radius"`
	Risk            string   `json:"risk"`
	Score           float64  `json:"score"`
	DependencyTypes []string `json:"dependency_types"`
}

type cliSinglePointOfTrustFailure struct {
	NodeID          string   `json:"node_id"`
	AssumptionText  string   `json:"assumption_text"`
	DependentsCount int      `json:"dependents_count"`
	DependentNodes  []string `json:"dependent_nodes"`
	DependencyTypes []string `json:"dependency_types"`
	Recommendation  string   `json:"recommendation"`
}

type cliTrustCollapseResult struct {
	FailedAssumptionID   string   `json:"failed_assumption_id"`
	FailedAssumptionText string   `json:"failed_assumption_text"`
	AssumptionsLost      []string `json:"assumptions_lost"`
	ControlsLost         []string `json:"controls_lost,omitempty"`
	AssetsExposed        []string `json:"assets_exposed,omitempty"`
	RiskIncrease         string   `json:"risk_increase"`
	RiskScoreBefore      float64  `json:"risk_score_before"`
	RiskScoreAfter       float64  `json:"risk_score_after"`
	AffectedComponents   []string `json:"affected_components,omitempty"`
}

// APD CLI types
type cliAttackPath struct {
	ID                  string          `json:"id"`
	Name                string          `json:"name"`
	Description         string          `json:"description"`
	EntryPoint          string          `json:"entry_point"`
	TargetAsset         string          `json:"target_asset"`
	AttackSteps         []cliAttackStep `json:"attack_steps"`
	Likelihood          float64         `json:"likelihood"`
	Impact              float64         `json:"impact"`
	RiskScore           float64         `json:"risk_score"`
	DetectionDifficulty string          `json:"detection_difficulty"`
	BusinessImpact      string          `json:"business_impact"`
	Recommendations     []string        `json:"recommendations,omitempty"`
	KillChainPhases     []string        `json:"kill_chain_phases,omitempty"`
	MITREATTACK         []string        `json:"mitre_attack,omitempty"`
}

type cliAttackStep struct {
	SequenceNumber     int    `json:"sequence_number"`
	SourceComponent    string `json:"source_component"`
	TargetComponent    string `json:"target_component"`
	Action             string `json:"action"`
	Threat             string `json:"threat"`
	RequiredAssumption string `json:"required_assumption"`
	ControlBypassed    string `json:"control_bypassed"`
	Reasoning          string `json:"reasoning"`
	STRIDECategory     string `json:"stride_category"`
}

type cliThreatChain struct {
	ID        string   `json:"id"`
	Threats   []string `json:"threats"`
	Path      []string `json:"path"`
	RiskScore float64  `json:"risk_score"`
	Reasoning string   `json:"reasoning"`
}

type cliAttackPathSummary struct {
	TotalAttackPaths  int            `json:"total_attack_paths"`
	CriticalCount     int            `json:"critical_count"`
	HighCount         int            `json:"high_count"`
	MediumCount       int            `json:"medium_count"`
	LowCount          int            `json:"low_count"`
	ThreatChainCount  int            `json:"threat_chain_count"`
	TopAttackPaths    []string       `json:"top_attack_paths"`
	KillChainCoverage map[string]int `json:"kill_chain_coverage"`
	MITRECoverage     []string       `json:"mitre_coverage"`
	SummaryText       string         `json:"summary_text"`
}

// SDRI CLI types
type cliSDRIControl struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
	ControlType string `json:"control_type"`
	Preventive  bool   `json:"preventive"`
	Detective   bool   `json:"detective"`
	Corrective  bool   `json:"corrective"`
	Strength    string `json:"strength"`
	Coverage    string `json:"coverage"`
	Status      string `json:"status"`
}

type cliSDRIFinding struct {
	ID                 string   `json:"id"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	Severity           string   `json:"severity"`
	Category           string   `json:"category"`
	AffectedComponents []string `json:"affected_components,omitempty"`
	AffectedControls   []string `json:"affected_controls,omitempty"`
	BusinessImpact     string   `json:"business_impact"`
	Recommendation     string   `json:"recommendation"`
}

type cliSDRIWeakness struct {
	ID             string   `json:"id"`
	Pattern        string   `json:"pattern"`
	Description    string   `json:"description"`
	Severity       string   `json:"severity"`
	Components     []string `json:"components,omitempty"`
	Impact         string   `json:"impact"`
	Recommendation string   `json:"recommendation"`
}

type cliSDRIRemediation struct {
	ID             string  `json:"id"`
	Priority       int     `json:"priority"`
	Description    string  `json:"description"`
	RiskScore      float64 `json:"risk_score"`
	BusinessImpact string  `json:"business_impact"`
	Effort         string  `json:"effort"`
	Category       string  `json:"category"`
	Recommendation string  `json:"recommendation"`
}

type cliSDRICoverage struct {
	Category string  `json:"category"`
	Expected int     `json:"expected"`
	Observed int     `json:"observed"`
	Coverage float64 `json:"coverage"`
	Level    string  `json:"level"`
}

type cliSDRICompliance struct {
	Framework string   `json:"framework"`
	Coverage  float64  `json:"coverage"`
	Controls  []string `json:"controls,omitempty"`
	Status    string   `json:"status"`
}

// ── CIARE CLI Types ──

type cliCIAREFrameworkCoverage struct {
	Framework        string   `json:"framework"`
	Required         int      `json:"required"`
	Observed         int      `json:"observed"`
	Missing          int      `json:"missing"`
	CoveragePct      float64  `json:"coverage_pct"`
	Status           string   `json:"status"`
	ObservedControls []string `json:"observed_controls,omitempty"`
	MissingControls  []string `json:"missing_controls,omitempty"`
}

type cliCIAREAuditReadiness struct {
	Framework       string   `json:"framework"`
	ReadinessScore  float64  `json:"readiness_score"`
	Status          string   `json:"status"`
	ControlCoverage float64  `json:"control_coverage"`
	EvidenceScore   float64  `json:"evidence_score"`
	ThreatExposure  float64  `json:"threat_exposure"`
	FindingsPenalty float64  `json:"findings_penalty"`
	Factors         []string `json:"factors,omitempty"`
}

type cliCIAREEvidenceRequirement struct {
	Framework string   `json:"framework"`
	Control   string   `json:"control"`
	Evidence  []string `json:"evidence"`
}

type cliCIAREMissingEvidence struct {
	Framework string   `json:"framework"`
	Control   string   `json:"control"`
	Evidences []string `json:"evidences"`
}

type cliCIAREAuditorQuestion struct {
	Framework string `json:"framework"`
	Control   string `json:"control"`
	Question  string `json:"question"`
}

type cliCIAREComplianceGap struct {
	ID          string `json:"id"`
	Framework   string `json:"framework"`
	Requirement string `json:"requirement"`
	Observed    string `json:"observed"`
	Missing     string `json:"missing"`
	Risk        string `json:"risk"`
}

type cliCIAREControlMaturity struct {
	Domain   string  `json:"domain"`
	Level    int     `json:"level"`
	Label    string  `json:"label"`
	Coverage float64 `json:"coverage"`
}

type cliCIAREComplianceNarrative struct {
	Framework string `json:"framework"`
	Narrative string `json:"narrative"`
}

type cliCIAREAuditPackage struct {
	ExecutiveSummary     string                        `json:"executive_summary"`
	FrameworkCoverages   []cliCIAREFrameworkCoverage   `json:"framework_coverages,omitempty"`
	ControlInventory     []cliSDRIControl              `json:"control_inventory,omitempty"`
	MissingControls      []cliCIAREComplianceGap       `json:"missing_controls,omitempty"`
	EvidenceRequirements []cliCIAREEvidenceRequirement `json:"evidence_requirements,omitempty"`
	AuditorQuestions     []cliCIAREAuditorQuestion     `json:"auditor_questions,omitempty"`
}

type cliCIAREComplianceDashboard struct {
	FrameworkCoverages map[string]float64        `json:"framework_coverages"`
	TopGaps            []cliCIAREComplianceGap   `json:"top_gaps,omitempty"`
	TopMissingEvidence []cliCIAREMissingEvidence `json:"top_missing_evidence,omitempty"`
	TopRisks           []string                  `json:"top_risks,omitempty"`
}

type cliCIAREProcurementQuestion struct {
	Category string `json:"category"`
	Question string `json:"question"`
}
