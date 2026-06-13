package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"asf-tui/asf/analyzer"
	"asf-tui/asf/confidencex"
	"asf-tui/asf/coverage"
	"asf-tui/asf/models"
	"asf-tui/asf/narrative"
	"asf-tui/asf/review"
	"asf-tui/asf/trust"
	"asf-tui/asf/verify"
	"asf-tui/intelligence"
)

const (
	MaxFileSize             = 50 * 1024 * 1024
	MaxAIDisplayAssumptions = 50
	MaxTUIAssumptions       = 500
)

const (
	ModeASFOnly  = "ASF Engine Only"
	ModeASFAndAI = "ASF Engine + Local AI"
)

type RiskLevel string

const (
	RiskCritical RiskLevel = "Critical"
	RiskHigh     RiskLevel = "High"
	RiskMedium   RiskLevel = "Medium"
	RiskLow      RiskLevel = "Low"
)

type StrideCategory string

const (
	StrideSpoofing        StrideCategory = "Spoofing"
	StrideTampering       StrideCategory = "Tampering"
	StrideRepudiation     StrideCategory = "Repudiation"
	StrideInfoDisclosure  StrideCategory = "Information Disclosure"
	StrideDenialOfService StrideCategory = "Denial of Service"
	StrideElevationPriv   StrideCategory = "Elevation of Privilege"
)

type Assumption struct {
	ID          string
	Description string
	Component   string
	Category    string
	Risk        RiskLevel
	Stride      []StrideCategory
	Likelihood  int
	Impact      int
	Confidence  float64
	Keywords    []string

	// Evidence Traceability
	SourceNode string `json:"source_node"`
	SourceLine int    `json:"source_line"`

	// Source metadata (populated for explicit YAML/JSON assumptions)
	SourceType    string `json:"source_type"`    // "explicit", "inferred", "generated"
	SourceSection string `json:"source_section"` // "assumptions", "security_controls", "diagram"
	SourceIndex   int    `json:"source_index"`   // index within that section
	SourceFile    string `json:"source_file"`    // original architecture file

	// Verification status (wired from security controls / evidence)
	VerificationStatus string `json:"verification_status"` // "VERIFIED", "PARTIALLY_VERIFIED", "CONTRADICTED", "UNKNOWN"

	// Explainability fields (added by the explainability engine)
	EvidenceSources      []string              `json:"evidence_sources"`
	SourceComponents     []string              `json:"source_components"`
	SourceRelationships  []string              `json:"source_relationships"`
	Rationale            string                `json:"rationale"`
	StrideJustifications []StrideJustification `json:"stride_justifications"`
	RiskJustification    *RiskJustification    `json:"risk_justification"`
	ReviewStatus         string                `json:"review_status"`
	ReviewNotes          string                `json:"review_notes"`
	ReviewTimestamp      time.Time             `json:"review_timestamp"`

	QualityScore float64 `json:"quality_score,omitempty"`
}

// Contradiction represents a detected logical contradiction.
type Contradiction struct {
	ID                  string    `json:"id"`
	Severity            RiskLevel `json:"severity"`
	Description         string    `json:"description"`
	Explanation         string    `json:"explanation"`
	AffectedAssumptions []string  `json:"affected_assumptions"`
	Evidence            []string  `json:"evidence"`
	RuleName            string    `json:"rule_name"`
}

// TrustBoundary represents a discovered trust boundary.
type TrustBoundary struct {
	Type        string    `json:"type"`
	Components  []string  `json:"components"`
	RiskLevel   RiskLevel `json:"risk_level"`
	Description string    `json:"description"`
}

// QualityScore represents the quality score of an assumption.
type QualityScore struct {
	Hiddenness             float64 `json:"hiddenness"`
	Impact                 float64 `json:"impact"`
	Novelty                float64 `json:"novelty"`
	ArchitecturalRelevance float64 `json:"architectural_relevance"`
	Risk                   float64 `json:"risk"`
	Confidence             float64 `json:"confidence"`
	Overall                float64 `json:"overall"`
}

// CIEContradiction represents a rich contradiction from the Contradiction Intelligence Engine.
type CIEContradiction struct {
	ID                      string       `json:"id"`
	Type                    string       `json:"type"`
	Severity                RiskLevel    `json:"severity"`
	Confidence              float64      `json:"confidence"`
	Summary                 string       `json:"summary"`
	Description             string       `json:"description"`
	StatementA              CIEStatement `json:"statement_a"`
	StatementB              CIEStatement `json:"statement_b"`
	AffectedAssets          []string     `json:"affected_assets,omitempty"`
	AffectedComponents      []string     `json:"affected_components,omitempty"`
	AffectedControls        []string     `json:"affected_controls,omitempty"`
	AffectedTrustBoundaries []string     `json:"affected_trust_boundaries,omitempty"`
	Reasoning               string       `json:"reasoning"`
	Evidence                []string     `json:"evidence,omitempty"`
	Recommendations         []string     `json:"recommendations,omitempty"`
}

// CIEStatement represents a normalized claim in a contradiction.
type CIEStatement struct {
	ID           string  `json:"id"`
	Source       string  `json:"source"`
	OriginalText string  `json:"original_text"`
	Category     string  `json:"category"`
	Confidence   float64 `json:"confidence"`
}

// TBITrustZone represents a trust zone for TBI output.
type TBITrustZone struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Sensitivity string   `json:"sensitivity"`
	Components  []string `json:"components"`
	Description string   `json:"description"`
}

// TBITrustBoundary represents a trust boundary for TBI output.
type TBITrustBoundary struct {
	ID                  string    `json:"id"`
	SourceZone          string    `json:"source_zone"`
	DestinationZone     string    `json:"destination_zone"`
	SourceZoneType      string    `json:"source_zone_type"`
	DestinationZoneType string    `json:"destination_zone_type"`
	CrossingType        string    `json:"crossing_type"`
	Risk                RiskLevel `json:"risk"`
	Confidence          float64   `json:"confidence"`
	RequiredControls    []string  `json:"required_controls"`
	RequiredAssumptions []string  `json:"required_assumptions"`
	Threats             []string  `json:"threats"`
	MissingControls     []string  `json:"missing_controls,omitempty"`
	MissingAssumptions  []string  `json:"missing_assumptions,omitempty"`
	Reasoning           string    `json:"reasoning"`
	Recommendations     []string  `json:"recommendations,omitempty"`
	ComplianceMappings  []string  `json:"compliance_mappings,omitempty"`
}

// TBIWeakness represents a boundary weakness for TBI output.
type TBIWeakness struct {
	ID              string    `json:"id"`
	BoundaryID      string    `json:"boundary_id"`
	Type            string    `json:"type"`
	Severity        RiskLevel `json:"severity"`
	Description     string    `json:"description"`
	Reasoning       string    `json:"reasoning"`
	Recommendations []string  `json:"recommendations,omitempty"`
}

// Threat represents a generated threat from TMI.
type Threat struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	Category           string    `json:"category"`
	Severity           RiskLevel `json:"severity"`
	Likelihood         float64   `json:"likelihood"`
	Impact             float64   `json:"impact"`
	RiskScore          float64   `json:"risk_score"`
	Confidence         float64   `json:"confidence"`
	Description        string    `json:"description"`
	AffectedAssets     []string  `json:"affected_assets,omitempty"`
	AffectedComponents []string  `json:"affected_components,omitempty"`
	AffectedBoundaries []string  `json:"affected_boundaries,omitempty"`
	AffectedData       []string  `json:"affected_data,omitempty"`
	Assumptions        []string  `json:"assumptions,omitempty"`
	Controls           []string  `json:"controls,omitempty"`
	STRIDECategories   []string  `json:"stride_categories,omitempty"`
	Reasoning          string    `json:"reasoning"`
	Recommendations    []string  `json:"recommendations,omitempty"`
	PreventiveControls []string  `json:"preventive_controls,omitempty"`
	DetectiveControls  []string  `json:"detective_controls,omitempty"`
	CorrectiveControls []string  `json:"corrective_controls,omitempty"`
}

// ThreatCluster represents a group of related threats.
type ThreatCluster struct {
	ID              string   `json:"id"`
	Name            string   `json:"name"`
	Category        string   `json:"category"`
	Threats         []string `json:"threats"`
	RiskScore       float64  `json:"risk_score"`
	AffectedAssets  []string `json:"affected_assets,omitempty"`
	Recommendations []string `json:"recommendations,omitempty"`
}

// ThreatModelSummary represents a summary of threat modeling results.
type ThreatModelSummary struct {
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

type AnalysisResult struct {
	ArchitectureName   string
	AnalysisDate       time.Time
	AnalysisMode       string
	Assumptions        []Assumption
	CriticalCount      int
	HighCount          int
	MediumCount        int
	LowCount           int
	TotalAssumptions   int
	StrideDistribution map[StrideCategory]int
	Controls           []ControlDetail
	Compliance         []string
	Summary            string
	TrueAssumptions    int
	FalseAssumptions   int
	CriticalGaps       int

	// Explainability fields
	EvidenceSummary   EvidenceSummary `json:"evidence_summary"`
	RiskModelVersion  string          `json:"risk_model_version"`
	ConfidenceSummary string          `json:"confidence_summary"`

	// Intelligence engine fields
	Contradictions      []Contradiction         `json:"contradictions,omitempty"`
	TrustBoundaries     []TrustBoundary         `json:"trust_boundaries,omitempty"`
	QualityScores       map[string]QualityScore `json:"quality_scores,omitempty"`
	Domain              string                  `json:"domain,omitempty"`
	IntelligenceSummary string                  `json:"intelligence_summary,omitempty"`

	// Contradiction Intelligence Engine (CIE) fields
	CIEContradictions []CIEContradiction `json:"cie_contradictions,omitempty"`
	CIESummary        string             `json:"cie_summary,omitempty"`

	// Trust Boundary Intelligence Engine (TBI) fields
	TBIZones      []TBITrustZone     `json:"tbi_zones,omitempty"`
	TBIBoundaries []TBITrustBoundary `json:"tbi_boundaries,omitempty"`
	TBIWeaknesses []TBIWeakness      `json:"tbi_weaknesses,omitempty"`
	TBISummary    string             `json:"tbi_summary,omitempty"`

	// Threat Modeling Intelligence Engine (TMI) fields
	Threats            []Threat           `json:"threats,omitempty"`
	ThreatClusters     []ThreatCluster    `json:"threat_clusters,omitempty"`
	ThreatModelSummary ThreatModelSummary `json:"threat_model_summary,omitempty"`

	// Attack Path Discovery Engine (APD) fields
	AttackPaths       []AttackPath      `json:"attack_paths,omitempty"`
	ThreatChains      []ThreatChain     `json:"threat_chains,omitempty"`
	AttackPathSummary AttackPathSummary `json:"attack_path_summary,omitempty"`

	// Security Design Review Intelligence (SDRI) fields
	SDRIControls               []SDRIControl               `json:"sdri_controls,omitempty"`
	SDRIDesignFindings         []SDRIDesignFinding         `json:"sdri_design_findings,omitempty"`
	SDRIAchitecturalWeaknesses []SDRIArchitecturalWeakness `json:"sdri_architectural_weaknesses,omitempty"`
	SDRIRemediations           []SDRIRemediation           `json:"sdri_remediations,omitempty"`
	SDRICoverageByCategory     []SDRICoverageItem          `json:"sdri_coverage_by_category,omitempty"`
	SDRICoverageDashboard      map[string]float64          `json:"sdri_coverage_dashboard,omitempty"`
	SDRIComplianceAlignments   []SDRIComplianceMapping     `json:"sdri_compliance_alignments,omitempty"`
	SDRISummary                string                      `json:"sdri_summary,omitempty"`

	// Compliance Intelligence & Audit Readiness Engine (CIARE) fields
	CIAREFrameworkCoverages   []CIAREFrameworkCoverage   `json:"ciare_framework_coverages,omitempty"`
	CIAREAuditReadiness       []CIAREAuditReadiness      `json:"ciare_audit_readiness,omitempty"`
	CIAREEvidenceRequirements []CIAREEvidenceRequirement `json:"ciare_evidence_requirements,omitempty"`
	CIAREMissingEvidences     []CIAREMissingEvidence     `json:"ciare_missing_evidences,omitempty"`
	CIAREAuditorQuestions     []CIAREAuditorQuestion     `json:"ciare_auditor_questions,omitempty"`
	CIAREComplianceGaps       []CIAREComplianceGap       `json:"ciare_compliance_gaps,omitempty"`
	CIAREControlMaturities    []CIAREControlMaturity     `json:"ciare_control_maturities,omitempty"`
	CIAREComplianceNarratives []CIAREComplianceNarrative `json:"ciare_compliance_narratives,omitempty"`
	CIAREAuditPackage         *CIAREAuditPackage         `json:"ciare_audit_package,omitempty"`
	CIAREComplianceDashboard  *CIAREComplianceDashboard  `json:"ciare_compliance_dashboard,omitempty"`
	CIAREProcurementQuestions []CIAREProcurementQuestion `json:"ciare_procurement_questions,omitempty"`

	// Domain Knowledge Pack Intelligence (DKPI) fields
	DKPI DKPIIntelligence `json:"dkpi,omitempty"`

	// Security Decision Intelligence (SDI) fields
	SDI SDIIntelligence `json:"sdi,omitempty"`

	// Security Digital Twin (SDT) fields
	SDT SDTIntelligence `json:"sdt,omitempty"`

	// Security Architect Narrative Engine (SANE) fields
	NarrativeOutput *narrative.NarrativeOutput `json:"narrative_output,omitempty"`

	// Assumption Dependency & Trust Chain Engine (V14) fields
	TrustOutput *trust.ChainOutput `json:"trust_output,omitempty"`

	CoverageOutput *coverage.CoverageOutput `json:"coverage_output,omitempty"`

	VerificationOutput *verify.VerificationOutput `json:"verification_output,omitempty"`

	ReviewOutput *review.ReviewOutput `json:"review_output,omitempty"`

	ConfidenceOutput *confidencex.ConfidenceOutput `json:"confidence_output,omitempty"`
}

type AnalysisProgress struct {
	Percent  float64
	Stage    string
	Complete bool
}

type Engine struct {
	config       *Config
	strideEngine *StrideEngine
	explainPipe  *ExplainabilityPipeline
	archDesc     *ArchDescription
}

func NewEngine(cfg *Config) *Engine {
	if err := ensureRuntimeDirs(); err != nil {
		debugLog.Printf("runtime dirs: %v", err)
	}
	return &Engine{
		config:       cfg,
		strideEngine: NewStrideEngine(),
	}
}

func (e *Engine) RunAnalysis(archPath, evPath, mode string, progress chan<- AnalysisProgress) (*AnalysisResult, error) {
	progress <- AnalysisProgress{Percent: 5, Stage: "Parsing Architecture..."}
	defer close(progress)

	fi, err := os.Stat(archPath)
	if err != nil {
		asfLog.Printf("analysis error: access arch file: %v", err)
		return nil, fmt.Errorf("access arch file: %w", err)
	}
	if fi.Size() > MaxFileSize {
		asfLog.Printf("analysis error: file too large (%d bytes, max %d)", fi.Size(), MaxFileSize)
		return nil, fmt.Errorf("file too large (%d bytes, max %d bytes)", fi.Size(), MaxFileSize)
	}

	asfLog.Printf("analysis start: %s (mode=%s, size=%d)", archPath, mode, fi.Size())

	inputPath := archPath
	ext := strings.ToLower(filepath.Ext(archPath))
	needsTemp := ext == ".drawio" || ext == ".mmd" || ext == ".md" || ext == ".yaml" || ext == ".yml" || ext == ".json" || ext == ".svg" || ext == ".png" || ext == ".jpg" || ext == ".jpeg"

	if needsTemp {
		desc, err := ParseArchitecture(archPath)
		if err != nil {
			asfLog.Printf("analysis error: parse architecture: %v", err)
			return nil, fmt.Errorf("parse architecture: %w", err)
		}
		e.archDesc = desc
		// Write parsed text to a fixed path in the cache directory for the native analyzer.
		// A file path is required because the native analyzer accepts file paths, not raw text.
		// Eliminating this entirely would require changing the analyzer API (engine change).
		inputPath = filepath.Join(asfCacheDir(), "analysis_input.txt")
		if err := os.WriteFile(inputPath, []byte(desc.RawText), 0644); err != nil {
			asfLog.Printf("analysis error: write input file: %v", err)
			return nil, fmt.Errorf("write input file: %w", err)
		}
	} else {
		desc, err := ParseArchitecture(archPath)
		if err == nil {
			e.archDesc = desc
		}
	}

	progress <- AnalysisProgress{Percent: 20, Stage: "Running ASF Engine..."}

	debugLog.Printf("using native Go engine")
	asfResult, err := e.runNativeAnalysis(inputPath, evPath)
	if err != nil {
		asfLog.Printf("analysis error: ASF engine: %v", err)
		return nil, fmt.Errorf("ASF engine error: %w", err)
	}

	progress <- AnalysisProgress{Percent: 60, Stage: "Processing Results..."}

	result := e.buildResult(asfResult, archPath, mode)

	// Run Intelligence Engine for advanced assumption discovery
	progress <- AnalysisProgress{Percent: 65, Stage: "Running Intelligence Engine..."}
	if e.archDesc != nil {
		ie := intelligence.NewIntelligenceEngine()
		intelArch := convertToIntelArch(e.archDesc)
		existingIntelAssumptions := convertAssumptionsToIntel(result.Assumptions)
		debugLog.Printf("intel: passing %d assumptions to intelligence engine", len(existingIntelAssumptions))
		intelResult := ie.RunWithExistingAssumptions(intelArch, existingIntelAssumptions)
		debugLog.Printf("intel: generated %d assumptions, %d contradictions, %d boundaries", len(intelResult.Assumptions), len(intelResult.Contradictions), len(intelResult.TrustBoundaries))

		// Merge intelligence results
		intelAssumptions := convertIntelAssumptions(intelResult.Assumptions)
		result.Assumptions = mergeAssumptions(result.Assumptions, intelAssumptions)
		result.Contradictions = deduplicateContradictions(convertIntelContradictions(intelResult.Contradictions))
		result.TrustBoundaries = convertIntelTrustBoundaries(intelResult.TrustBoundaries)
		result.Domain = intelResult.Domain
		result.IntelligenceSummary = intelResult.Summary

		// Run Contradiction Intelligence Engine (CIE)
		progress <- AnalysisProgress{Percent: 70, Stage: "Running Contradiction Intelligence Engine..."}
		cie := intelligence.NewCIEEngine()
		cieContradictions := cie.DetectAllContradictions(intelArch, existingIntelAssumptions, convertControlsToIntel(result.Controls), convertTrustBoundariesToIntel(result.TrustBoundaries))
		result.CIEContradictions = deduplicateCIEContradictions(convertCIEContradictions(cieContradictions))
		result.CIESummary = intelligence.BuildContradictionSummary(cieContradictions)
		debugLog.Printf("cie: detected %d contradictions", len(cieContradictions))

		// Run Trust Boundary Intelligence Engine (TBI)
		progress <- AnalysisProgress{Percent: 75, Stage: "Running Trust Boundary Intelligence Engine..."}
		tbi := intelligence.NewTBIEngine()
		tbiResult, err := tbi.Run(intelArch, existingIntelAssumptions)
		if err == nil && tbiResult != nil {
			result.TBIZones = convertTBIZones(tbiResult.Zones)
			result.TBIBoundaries = convertTBIBoundaries(tbiResult.Boundaries)
			result.TBIWeaknesses = convertTBIWeaknesses(tbiResult.Weaknesses)
			result.TBISummary = tbiResult.Summary
			// Merge TBI-generated assumptions
			intelAssumptions = mergeAssumptions(intelAssumptions, convertIntelAssumptions(tbiResult.Assumptions))
			result.Assumptions = intelAssumptions
			debugLog.Printf("tbi: discovered %d zones, %d boundaries, %d weaknesses", len(tbiResult.Zones), len(tbiResult.Boundaries), len(tbiResult.Weaknesses))
		}

		// Run Threat Modeling Intelligence Engine (TMI)
		progress <- AnalysisProgress{Percent: 78, Stage: "Running Threat Modeling Intelligence Engine..."}
		tmi := intelligence.NewTMIEngine()
		tmiResult := tmi.Run(intelArch, existingIntelAssumptions, convertTBIBoundariesToIntel(result.TBIBoundaries))
		if tmiResult != nil {
			result.Threats = convertIntelThreats(tmiResult.Threats)
			result.ThreatClusters = convertIntelThreatClusters(tmiResult.Clusters)
			result.ThreatModelSummary = convertIntelThreatModelSummary(tmiResult.Summary)
			debugLog.Printf("tmi: generated %d threats, %d clusters", len(tmiResult.Threats), len(tmiResult.Clusters))
		}

		// Run Attack Path Discovery Engine (APD)
		progress <- AnalysisProgress{Percent: 82, Stage: "Running Attack Path Discovery Engine..."}
		apd := intelligence.NewAPDEngine()
		apdResult := apd.Run(intelArch, tmiResult.Threats, convertTBIBoundariesToIntel(result.TBIBoundaries), convertTBIZonesToIntel(result.TBIZones), existingIntelAssumptions)
		if apdResult != nil {
			result.AttackPaths = convertAPDAttackPaths(apdResult.AttackPaths)
			result.ThreatChains = convertAPDThreatChains(apdResult.ThreatChains)
			result.AttackPathSummary = convertAPDSummary(apdResult)
			debugLog.Printf("apd: discovered %d attack paths, %d threat chains", len(apdResult.AttackPaths), len(apdResult.ThreatChains))
		}

		// Run Security Design Review Intelligence Engine (SDRI)
		progress <- AnalysisProgress{Percent: 84, Stage: "Running Security Design Review Intelligence Engine..."}
		sdri := intelligence.NewSDRIEngine()
		apdPaths := make([]intelligence.AttackPath, 0)
		tmiThreats := make([]intelligence.Threat, 0)
		if apdResult != nil {
			apdPaths = apdResult.AttackPaths
		}
		if tmiResult != nil {
			tmiThreats = tmiResult.Threats
		}
		sdriResult := sdri.Run(intelArch, existingIntelAssumptions, convertControlsToIntel(result.Controls), apdPaths, tmiThreats, result.Domain)
		if sdriResult != nil {
			result.SDRIControls = convertSDRIControls(sdriResult.Controls)
			result.SDRIDesignFindings = convertSDRIDesignFindings(sdriResult.DesignFindings)
			result.SDRIAchitecturalWeaknesses = convertSDRIWeaknesses(sdriResult.ArchitecturalWeaknesses)
			result.SDRIRemediations = convertSDRIRemediations(sdriResult.Remediations)
			result.SDRICoverageByCategory = convertSDRICoverage(sdriResult.CoverageByCategory)
			result.SDRICoverageDashboard = sdriResult.CoverageDashboard
			result.SDRIComplianceAlignments = convertSDRIComplianceMappings(sdriResult.ComplianceAlignments)
			result.SDRISummary = sdriResult.ExecutiveSummary
			debugLog.Printf("sdri: %d findings, %d weaknesses, %d remediations, coverage %.1f%%",
				len(sdriResult.DesignFindings), len(sdriResult.ArchitecturalWeaknesses),
				len(sdriResult.Remediations), averageCoverage(sdriResult.CoverageByCategory))
		}

		// Run Compliance Intelligence & Audit Readiness Engine (CIARE)
		progress <- AnalysisProgress{Percent: 86, Stage: "Running Compliance Intelligence Engine..."}
		ciare := intelligence.NewCIAREEngine()
		ciareInput := intelligence.CIAREInput{
			Architecture: intelArch,
			SDRIResult:   sdriResult,
			Domain:       result.Domain,
			Compliance:   result.Compliance,
		}
		ciareResult := ciare.Run(ciareInput)
		if ciareResult != nil {
			result.CIAREFrameworkCoverages = convertCIAREFrameworkCoverages(ciareResult.FrameworkCoverages)
			result.CIAREAuditReadiness = convertCIAREAuditReadiness(ciareResult.AuditReadinessScores)
			result.CIAREEvidenceRequirements = convertCIAREEvidenceRequirements(ciareResult.EvidenceRequirements)
			result.CIAREMissingEvidences = convertCIAREMissingEvidences(ciareResult.MissingEvidences)
			result.CIAREAuditorQuestions = convertCIAREAuditorQuestions(ciareResult.AuditorQuestions)
			result.CIAREComplianceGaps = convertCIAREComplianceGaps(ciareResult.ComplianceGaps)
			result.CIAREControlMaturities = convertCIAREControlMaturities(ciareResult.ControlMaturities)
			result.CIAREComplianceNarratives = convertCIAREComplianceNarratives(ciareResult.ComplianceNarratives)
			result.CIAREAuditPackage = convertCIAREAuditPackage(ciareResult.AuditPackage)
			result.CIAREComplianceDashboard = convertCIAREComplianceDashboard(ciareResult.ComplianceDashboard)
			result.CIAREProcurementQuestions = convertCIAREProcurementQuestions(ciareResult.ProcurementQuestions)
			debugLog.Printf("ciare: %d frameworks, %.1f%% avg coverage, %d gaps, %d missing evidence, %d readiness scores",
				len(ciareResult.FrameworkCoverages), ciareAvgCoverage(ciareResult.FrameworkCoverages),
				len(ciareResult.ComplianceGaps), len(ciareResult.MissingEvidences),
				len(ciareResult.AuditReadinessScores))
		}

		// Run Domain Knowledge Pack Intelligence Engine (DKPI)
		progress <- AnalysisProgress{Percent: 88, Stage: "Running Domain Knowledge Pack Intelligence Engine..."}
		dkpi := intelligence.NewDKPIEngine()
		dkpiExistingThreats := make([]intelligence.Threat, 0)
		if tmiResult != nil {
			dkpiExistingThreats = tmiResult.Threats
		}
		dkpiExistingControls := make([]intelligence.SDRIControl, 0)
		dkpiExistingFindings := make([]intelligence.SDRIFinding, 0)
		if sdriResult != nil {
			dkpiExistingControls = sdriResult.Controls
			dkpiExistingFindings = sdriResult.DesignFindings
		}
		dkpiInput := intelligence.DKPIInput{
			Architecture:        intelArch,
			ExistingAssumptions: existingIntelAssumptions,
			ExistingThreats:     dkpiExistingThreats,
			ExistingControls:    dkpiExistingControls,
			ExistingFindings:    dkpiExistingFindings,
			Domain:              result.Domain,
			Compliance:          result.Compliance,
		}
		dkpiResult := dkpi.Run(dkpiInput)
		if dkpiResult != nil {
			result.DKPI = convertDKPIResult(dkpiResult)
			// Update boosted assumptions and enriched controls
			if len(dkpiResult.BoostedAssumptions) > 0 {
				boosted := convertIntelAssumptions(dkpiResult.BoostedAssumptions)
				result.Assumptions = mergeAssumptions(result.Assumptions, boosted)
			}
			if len(dkpiResult.EnrichedControls) > 0 {
				result.SDRIControls = mergeSDRIControls(result.SDRIControls, convertSDRIControls(dkpiResult.EnrichedControls))
			}
			debugLog.Printf("dkpi: domain=%s, confidence=%.1f%%, %d threats, %d recommendations",
				dkpiResult.DetectedDomain.PrimaryDomain, dkpiResult.DetectedDomain.Confidence,
				len(dkpiResult.InjectedThreats), len(dkpiResult.Recommendations))
		}

		// Run Security Decision Intelligence Engine (SDI)
		progress <- AnalysisProgress{Percent: 96, Stage: "Running Security Decision Intelligence Engine..."}
		sdi := intelligence.NewSDIEngine()
		sdiFindings := make([]intelligence.SDRIFinding, 0)
		sdiThreats := make([]intelligence.Threat, 0)
		sdiAttackPaths := make([]intelligence.AttackPath, 0)
		sdiControls := make([]intelligence.SDRIControl, 0)
		if tmiResult != nil {
			sdiThreats = tmiResult.Threats
		}
		if apdResult != nil {
			sdiAttackPaths = apdResult.AttackPaths
		}
		if sdriResult != nil {
			sdiControls = sdriResult.Controls
			sdiFindings = sdriResult.DesignFindings
		}
		sdiInput := intelligence.SDIInput{
			ArchitectureName:   result.ArchitectureName,
			Domain:             result.Domain,
			Findings:           sdiFindings,
			Threats:            sdiThreats,
			AttackPaths:        sdiAttackPaths,
			Controls:           sdiControls,
			Compliance:         result.Compliance,
			RiskScore:          float64(result.CriticalCount)*4 + float64(result.HighCount)*3 + float64(result.MediumCount)*2 + float64(result.LowCount)*1,
			CoverageByCategory: nil,
			Assumptions:        nil,
			AnalysisMode:       result.AnalysisMode,
		}
		sdiResult := sdi.Run(sdiInput)
		if sdiResult != nil {
			result.SDI = convertSDIResult(sdiResult)
			debugLog.Printf("sdi: %d recommendations, %d fix simulations, %d failure simulations",
				len(sdiResult.Recommendations), len(sdiResult.FixSimulations), len(sdiResult.FailureSimulations))
		}

		// Run Security Digital Twin Engine (SDT)
		sdtFindings := make([]intelligence.SDRIFinding, 0)
		sdtThreats := make([]intelligence.Threat, 0)
		sdtAttackPaths := make([]intelligence.AttackPath, 0)
		sdtControls := make([]intelligence.SDRIControl, 0)
		if tmiResult != nil {
			sdtThreats = tmiResult.Threats
		}
		if apdResult != nil {
			sdtAttackPaths = apdResult.AttackPaths
		}
		if sdriResult != nil {
			sdtControls = sdriResult.Controls
			sdtFindings = sdriResult.DesignFindings
		}
		sdtInput := intelligence.SDTInput{
			ArchitectureName: result.ArchitectureName,
			Domain:           result.Domain,
			RiskScore:        float64(result.CriticalCount)*4 + float64(result.HighCount)*3 + float64(result.MediumCount)*2 + float64(result.LowCount)*1,
			Coverage:         0,
			Findings:         sdtFindings,
			Threats:          sdtThreats,
			Controls:         sdtControls,
			Compliance:       result.Compliance,
			AttackPaths:      sdtAttackPaths,
		}
		for _, a := range result.Assumptions {
			sdtInput.Assumptions = append(sdtInput.Assumptions, intelligence.Assumption{ID: a.ID, Description: a.Description, VerificationStatus: a.VerificationStatus})
		}

		sdt := intelligence.NewSDTEngine()
		sdtResult := sdt.Run(sdtInput)
		if sdtResult != nil {
			result.SDT = convertSDTResult(sdtResult)
			debugLog.Printf("sdt: %d change impacts, %d control drifts, %d what-if scenarios",
				len(sdtResult.ChangeImpacts), len(sdtResult.ControlDrifts), len(sdtResult.WhatIfScenarios))
		}

		// Recompute counts after merge
		result.TotalAssumptions = len(result.Assumptions)
		result.CriticalCount = 0
		result.HighCount = 0
		result.MediumCount = 0
		result.LowCount = 0
		for _, a := range result.Assumptions {
			switch a.Risk {
			case RiskCritical:
				result.CriticalCount++
			case RiskHigh:
				result.HighCount++
			case RiskMedium:
				result.MediumCount++
			case RiskLow:
				result.LowCount++
			}
		}
	}

	// Re-apply security control verification to all assumptions (native + intel-generated)
	// The buildResult function applies this, but intel engine replaces the assumption list at line 462,
	// losing the verification status. This ensures every assumption gets control-based verification.
	if e.archDesc != nil && e.archDesc.SecurityControls != nil {
		for i := range result.Assumptions {
			result.Assumptions[i] = applySecurityControlVerification(result.Assumptions[i], e.archDesc.SecurityControls)
		}
	}

	progress <- AnalysisProgress{Percent: 80, Stage: "Generating STRIDE Mapping..."}
	result.StrideDistribution = e.mapStrideDistribution(result.Assumptions)

	// Validate expected results now that StrideDistribution is populated
	if len(e.archDesc.ExpectedResults) > 0 {
		result.Summary = e.buildValidationSummary(result)
	}

	if mode == ModeASFAndAI && e.config != nil && e.config.AI.Enabled && e.config.AI.ActiveModel != "" {
		progress <- AnalysisProgress{Percent: 85, Stage: "Running AI Enhancement..."}
		enhancer := NewAIEnhancer()
		aiResult, err := enhancer.Enhance(result, e.config.AI.ActiveModel)
		if err == nil && aiResult != nil {
			result = mergeAIResults(result, aiResult)
			result.StrideDistribution = e.mapStrideDistribution(result.Assumptions)
		} else if err != nil {
			// AI failed — keep base results, add warning
			warning := fmt.Sprintf("Base ASF analysis completed. Local AI enhancement failed: %s", err.Error())
			result.AnalysisMode = ModeASFOnly
			// Prepend warning to summary so user sees it
			result.Summary = warning + "\n" + result.Summary
		}
	}

	progress <- AnalysisProgress{Percent: 100, Stage: "Complete", Complete: true}

	// Generate Security Architect Narrative
	result.NarrativeOutput = e.generateNarrativeOutput(result)

	// Run Assumption Dependency & Trust Chain Engine (V14)
	progress <- AnalysisProgress{Percent: 94, Stage: "Running Trust Chain Analysis..."}
	result.TrustOutput = e.runTrustChainAnalysis(result)

	// Run Assumption Coverage & Blind Spot Engine (V15)
	progress <- AnalysisProgress{Percent: 96, Stage: "Running Coverage & Blind Spot Analysis..."}
	result.CoverageOutput = e.runCoverageAnalysis(result)

	// Run Assumption Verification Intelligence Engine (V16)
	progress <- AnalysisProgress{Percent: 98, Stage: "Running Verification Intelligence Analysis..."}
	result.VerificationOutput = e.runVerificationAnalysis(result)

	// Run Security Review Workbench (V17)
	progress <- AnalysisProgress{Percent: 99, Stage: "Running Security Review Workbench..."}
	result.ReviewOutput = e.runReviewAnalysis(result)

	// Run Confidence & Explainability Engine (V18)
	progress <- AnalysisProgress{Percent: 100, Stage: "Running Confidence Explainability Analysis..."}
	result.ConfidenceOutput = e.runConfidenceExplainability(result)

	asfLog.Printf("analysis complete: %d assumptions, %d critical, %d high", result.TotalAssumptions, result.CriticalCount, result.HighCount)
	return result, nil
}

func (e *Engine) generateNarrativeOutput(result *AnalysisResult) *narrative.NarrativeOutput {
	// Convert engine assumptions to narrative assumptions
	var narrAssumptions []narrative.Assumption
	for _, a := range result.Assumptions {
		strideCats := make([]string, len(a.Stride))
		for i, s := range a.Stride {
			strideCats[i] = string(s)
		}
		narrAssumptions = append(narrAssumptions, narrative.Assumption{
			ID:                  a.ID,
			Description:         a.Description,
			Component:           a.Component,
			Category:            a.Category,
			Risk:                string(a.Risk),
			STRIDECategories:    strideCats,
			Likelihood:          a.Likelihood,
			Impact:              a.Impact,
			Confidence:          a.Confidence,
			Keywords:            a.Keywords,
			SourceComponents:    a.SourceComponents,
			SourceRelationships: a.SourceRelationships,
			Rationale:           a.Rationale,
			EvidenceSources:     a.EvidenceSources,
		})
	}

	// Convert controls
	var narrControls []narrative.ControlDetail
	for _, c := range result.Controls {
		strideCats := make([]string, len(c.MitigatedSTRIDE))
		for i, s := range c.MitigatedSTRIDE {
			strideCats[i] = string(s)
		}
		narrControls = append(narrControls, narrative.ControlDetail{
			Name:                 c.ID,
			Category:             c.Category,
			Description:          c.Description,
			Rationale:            c.Rationale,
			MitigatedAssumptions: c.MitigatedAssumptionIDs,
			STRIDECategories:     strideCats,
		})
	}

	// Convert trust boundaries
	var narrBoundaries []narrative.TrustBoundary
	for _, tb := range result.TrustBoundaries {
		narrBoundaries = append(narrBoundaries, narrative.TrustBoundary{
			Type:        tb.Type,
			Components:  tb.Components,
			RiskLevel:   string(tb.RiskLevel),
			Description: tb.Description,
		})
	}

	// Convert contradictions
	var narrContradictions []narrative.Contradiction
	for _, c := range result.Contradictions {
		narrContradictions = append(narrContradictions, narrative.Contradiction{
			ID:                  c.ID,
			Severity:            string(c.Severity),
			Description:         c.Description,
			Explanation:         c.Explanation,
			AffectedAssumptions: c.AffectedAssumptions,
		})
	}

	// Build stride and risk distributions
	strideDist := make(map[string]int)
	for cat, count := range result.StrideDistribution {
		strideDist[string(cat)] = count
	}
	riskDist := map[string]int{
		"Critical": result.CriticalCount,
		"High":     result.HighCount,
		"Medium":   result.MediumCount,
		"Low":      result.LowCount,
	}

	// Extract components
	var components []string
	seen := make(map[string]bool)
	for _, a := range result.Assumptions {
		if a.Component != "" && !seen[a.Component] {
			seen[a.Component] = true
			components = append(components, a.Component)
		}
	}

	engine := narrative.NewNarrativeEngine(result.Domain, components, nil)
	return engine.GenerateNarrative(
		result.ArchitectureName,
		narrAssumptions,
		narrControls,
		narrBoundaries,
		narrContradictions,
		result.Domain,
		strideDist,
		riskDist,
	)
}

func (e *Engine) runTrustChainAnalysis(result *AnalysisResult) *trust.ChainOutput {
	if len(result.Assumptions) == 0 {
		return nil
	}

	// Convert assumptions to trust engine input
	inputs := make([]trust.AssumptionInput, len(result.Assumptions))
	for i, a := range result.Assumptions {
		inputs[i] = trust.AssumptionInput{
			ID:         a.ID,
			Text:       a.Description,
			Component:  a.Component,
			Category:   a.Category,
			Risk:       string(a.Risk),
			Confidence: a.Confidence,
			Keywords:   a.Keywords,
			Source:     a.SourceType,
		}
	}

	// Extract unique components
	compSet := make(map[string]bool)
	for _, a := range result.Assumptions {
		if a.Component != "" {
			compSet[a.Component] = true
		}
	}
	components := make([]string, 0, len(compSet))
	for c := range compSet {
		components = append(components, c)
	}

	domain := result.Domain
	if domain == "" {
		domain = result.DKPI.DomainResult.PrimaryDomain
	}

	discovery := trust.NewDiscoveryEngine(domain, components)
	graph := discovery.DiscoverDependencies(inputs)
	engine := trust.NewTrustChainEngine(graph)
	output := engine.RunAll()
	output.Domain = domain

	return output
}

func (e *Engine) runCoverageAnalysis(result *AnalysisResult) *coverage.CoverageOutput {
	if len(result.Assumptions) == 0 {
		return nil
	}

	inputs := make([]coverage.AssumptionInput, len(result.Assumptions))
	for i, a := range result.Assumptions {
		inputs[i] = coverage.AssumptionInput{
			ID:          a.ID,
			Description: a.Description,
			Component:   a.Component,
			Category:    a.Category,
			Keywords:    a.Keywords,
			Risk:        string(a.Risk),
		}
	}

	compSet := make(map[string]bool)
	for _, a := range result.Assumptions {
		if a.Component != "" {
			compSet[a.Component] = true
		}
	}
	components := make([]string, 0, len(compSet))
	for c := range compSet {
		components = append(components, c)
	}

	domain := result.Domain
	if domain == "" {
		domain = result.DKPI.DomainResult.PrimaryDomain
	}

	engine := coverage.NewCoverageEngine(domain, components, inputs)
	return engine.RunAll()
}

func (e *Engine) runVerificationAnalysis(result *AnalysisResult) *verify.VerificationOutput {
	if len(result.Assumptions) == 0 {
		return nil
	}

	inputs := make([]verify.VerificationInput, len(result.Assumptions))
	for i, a := range result.Assumptions {
		inputs[i] = verify.VerificationInput{
			ID:          a.ID,
			Description: a.Description,
			Component:   a.Component,
			Category:    a.Category,
			Risk:        string(a.Risk),
			Keywords:    a.Keywords,
		}
	}

	compSet := make(map[string]bool)
	for _, a := range result.Assumptions {
		if a.Component != "" {
			compSet[a.Component] = true
		}
	}
	components := make([]string, 0, len(compSet))
	for c := range compSet {
		components = append(components, c)
	}

	domain := result.Domain
	if domain == "" {
		domain = result.DKPI.DomainResult.PrimaryDomain
	}

	engine := verify.NewVerificationEngine(domain, components, inputs)
	return engine.RunAll()
}

func (e *Engine) runReviewAnalysis(result *AnalysisResult) *review.ReviewOutput {
	if len(result.Assumptions) == 0 {
		return nil
	}

	inputs := make([]review.ReviewInput, len(result.Assumptions))
	for i, a := range result.Assumptions {
		verConf := a.Confidence * 100
		coverageGap := len(a.EvidenceSources) == 0

		// Look up blind spot score from coverage analysis
		blindSpot := 0.0
		if result.CoverageOutput != nil {
			for _, bs := range result.CoverageOutput.BlindSpots {
				if bs.Component == a.Component || bs.Category == coverage.CoverageCategory(a.Category) {
					blindSpot = bs.Score
					break
				}
			}
		}

		centrality := a.QualityScore
		if a.Confidence > 0 && centrality == 0 {
			centrality = a.Confidence * 0.5
		}

		supportCount := len(a.SourceComponents)
		if supportCount == 0 {
			supportCount = len(a.Keywords)
		}

		depCount := len(a.SourceRelationships)

		verPriority := string(a.Risk)
		if a.VerificationStatus == "" {
			verPriority = "Unverified"
		}

		inputs[i] = review.ReviewInput{
			AssumptionID:           a.ID,
			AssumptionText:         a.Description,
			Risk:                   string(a.Risk),
			Category:               a.Category,
			Component:              a.Component,
			Centrality:             centrality,
			Criticality:            riskToCriticality(a.Risk),
			FailureRadius:          a.Impact,
			SupportCount:           supportCount,
			DependencyCount:        depCount,
			VerificationPriority:   verPriority,
			VerificationConfidence: verConf,
			VerificationStatus:     a.VerificationStatus,
			CoverageGap:            coverageGap,
			BlindSpotScore:         blindSpot,
			Domain:                 result.Domain,
		}
	}

	engine := review.NewReviewEngine(result.Domain, inputs)
	return engine.RunAll()
}

func riskToCriticality(r RiskLevel) float64 {
	switch r {
	case RiskCritical:
		return 0.95
	case RiskHigh:
		return 0.80
	case RiskMedium:
		return 0.50
	default:
		return 0.20
	}
}

func (e *Engine) runConfidenceExplainability(result *AnalysisResult) *confidencex.ConfidenceOutput {
	if len(result.Assumptions) == 0 {
		return nil
	}

	inputs := make([]confidencex.ConfidenceInput, len(result.Assumptions))
	for i, a := range result.Assumptions {
		hasTrustChain := false
		if result.TrustOutput != nil && len(result.TrustOutput.TrustChains) > 0 {
			for _, tc := range result.TrustOutput.TrustChains {
				for _, n := range tc.Nodes {
					if n == a.ID {
						hasTrustChain = true
						break
					}
				}
				if hasTrustChain {
					break
				}
			}
		}

		hasCoverageGap := false
		blindSpotScore := 0.0
		if result.CoverageOutput != nil {
			for _, bs := range result.CoverageOutput.BlindSpots {
				if bs.Component == a.Component || bs.Category == coverage.CoverageCategory(a.Category) {
					hasCoverageGap = true
					blindSpotScore = bs.Score
					break
				}
			}
		}

		inputs[i] = confidencex.ConfidenceInput{
			AssumptionID:         a.ID,
			AssumptionText:       a.Description,
			Component:            a.Component,
			Category:             a.Category,
			Risk:                 string(a.Risk),
			Confidence:           a.Confidence * 100,
			EvidenceSources:      a.EvidenceSources,
			SourceComponents:     a.SourceComponents,
			SourceRelationships:  a.SourceRelationships,
			Keywords:             a.Keywords,
			Rationale:            a.Rationale,
			VerificationStatus:   a.VerificationStatus,
			Domain:               result.Domain,
			HasTrustChain:        hasTrustChain,
			HasCoverageGap:       hasCoverageGap,
			BlindSpotScore:       blindSpotScore,
			DependencyCentrality: float64(len(a.SourceRelationships)) / 10.0,
			FailureRadius:        a.Impact,
		}
	}

	engine := confidencex.NewExplainabilityEngine(result.Domain, inputs)
	return engine.RunAll()
}

func (e *Engine) runNativeAnalysis(docPath, evPath string) (*asfJSONResult, error) {
	docs := []string{docPath}
	var evs []string
	if evPath != "" {
		if fi, err := os.Stat(evPath); err != nil {
			debugLog.Printf("evidence path inaccessible: %s: %v", evPath, err)
		} else if fi.IsDir() {
			entries, err := os.ReadDir(evPath)
			if err != nil {
				debugLog.Printf("cannot read evidence dir: %s: %v", evPath, err)
			} else {
				evs = append(evs, evPath)
				if len(entries) == 0 {
					debugLog.Printf("warning: evidence directory is empty: %s", evPath)
				}
			}
		} else {
			evs = append(evs, evPath)
		}
	}
	if _, err := os.Stat(docPath); err != nil {
		debugLog.Printf("document path inaccessible: %s: %v", docPath, err)
	}

	an := analyzer.New()
	ar, err := an.Analyze(docs, evs)
	if err != nil {
		return nil, fmt.Errorf("native analysis: %w", err)
	}

	result := &asfJSONResult{}
	s := ar.Result.BuildSummary()
	result.Summary.ClaimsFound = s.ClaimsFound
	result.Summary.Assumptions = s.Assumptions
	result.Summary.Verified = s.Verified
	result.Summary.Contradicted = s.Contradicted
	result.Summary.Unknown = s.Unknown
	result.Summary.CriticalGaps = s.CriticalGaps

	for _, a := range ar.Result.Assumptions {
		result.Assumptions = append(result.Assumptions, struct {
			ID                 string   `json:"id"`
			Text               string   `json:"text"`
			AssumptionType     string   `json:"assumption_type"`
			VerificationStatus string   `json:"verification_status"`
			Confidence         float64  `json:"confidence"`
			Keywords           []string `json:"keywords"`
		}{
			ID:                 a.ID,
			Text:               a.Text,
			AssumptionType:     string(a.AssumptionType),
			VerificationStatus: string(a.VerificationStatus),
			Confidence:         a.Confidence,
			Keywords:           a.Keywords,
		})
	}

	for _, v := range ar.Result.Verifications {
		result.Verifications = append(result.Verifications, struct {
			AssumptionID string      `json:"assumption_id"`
			Result       string      `json:"result"`
			Confidence   float64     `json:"confidence"`
			EvidenceUsed interface{} `json:"evidence_used"`
			Reasoning    string      `json:"reasoning"`
		}{
			AssumptionID: v.AssumptionID,
			Result:       string(v.Result),
			Confidence:   v.Confidence,
			EvidenceUsed: v.EvidenceUsed,
			Reasoning:    v.Reasoning,
		})
	}

	for _, g := range ar.Result.Gaps {
		result.Gaps = append(result.Gaps, struct {
			AssumptionID   string `json:"assumption_id"`
			Type           string `json:"type"`
			Severity       string `json:"severity"`
			Description    string `json:"description"`
			EvidenceDetail string `json:"evidence_detail"`
		}{
			AssumptionID:   g.AssumptionID,
			Type:           string(g.Type),
			Severity:       string(g.Severity),
			Description:    g.Description,
			EvidenceDetail: g.EvidenceDetail,
		})
	}

	debugLog.Printf("runNativeAnalysis: %d assumptions, %d verifications, %d gaps",
		len(result.Assumptions), len(result.Verifications), len(result.Gaps))

	return result, nil
}

type asfJSONResult struct {
	Summary struct {
		ClaimsFound  int `json:"claims_found"`
		Assumptions  int `json:"assumptions"`
		Verified     int `json:"verified"`
		Contradicted int `json:"contradicted"`
		Unknown      int `json:"unknown"`
		CriticalGaps int `json:"critical_gaps"`
	} `json:"summary"`
	Assumptions []struct {
		ID                 string   `json:"id"`
		Text               string   `json:"text"`
		AssumptionType     string   `json:"assumption_type"`
		VerificationStatus string   `json:"verification_status"`
		Confidence         float64  `json:"confidence"`
		Keywords           []string `json:"keywords"`
	} `json:"assumptions"`
	Verifications []struct {
		AssumptionID string      `json:"assumption_id"`
		Result       string      `json:"result"`
		Confidence   float64     `json:"confidence"`
		EvidenceUsed interface{} `json:"evidence_used"`
		Reasoning    string      `json:"reasoning"`
	} `json:"verifications"`
	Gaps []struct {
		AssumptionID   string `json:"assumption_id"`
		Type           string `json:"type"`
		Severity       string `json:"severity"`
		Description    string `json:"description"`
		EvidenceDetail string `json:"evidence_detail"`
	} `json:"gaps"`
}

func (e *Engine) buildResult(r *asfJSONResult, archPath, mode string) *AnalysisResult {
	result := &AnalysisResult{
		ArchitectureName:   fileBase(archPath),
		AnalysisDate:       time.Now(),
		AnalysisMode:       mode,
		TotalAssumptions:   r.Summary.Assumptions,
		TrueAssumptions:    r.Summary.Verified,
		FalseAssumptions:   r.Summary.Contradicted,
		CriticalGaps:       r.Summary.CriticalGaps,
		StrideDistribution: make(map[StrideCategory]int),
		RiskModelVersion:   "asf-risk-model-1.0",
		Summary: fmt.Sprintf("ASF processed %s and found %d assumptions (%d verified, %d contradicted, %d unknown). %d critical gaps identified.",
			fileBase(archPath), r.Summary.Assumptions, r.Summary.Verified, r.Summary.Contradicted, r.Summary.Unknown, r.Summary.CriticalGaps),
	}

	verificationMap := make(map[string]string)
	for _, v := range r.Verifications {
		verificationMap[v.AssumptionID] = v.Result
	}

	gapMap := make(map[string]string)
	for _, g := range r.Gaps {
		gapMap[g.AssumptionID] = g.Severity
	}

	// Initialize explainability pipeline if we have architecture data
	if e.archDesc != nil && e.explainPipe == nil {
		e.explainPipe = NewExplainabilityPipeline(e.archDesc, archPath, e.strideEngine)
	}

	for _, a := range r.Assumptions {
		sev := gapMap[a.ID]
		risk := mapRiskLevel(sev, verificationMap[a.ID])

		stride := e.strideEngine.MapAssumption(a.AssumptionType, a.Text, a.Keywords)
		component := extractComponent(a.Keywords, a.Text)

		lh, im := riskToLikelihoodImpact(risk)
		desc := cleanAssumptionText(a.Text)
		assumption := Assumption{
			ID:          a.ID,
			Description: desc,
			Component:   component,
			Category:    a.AssumptionType,
			Risk:        risk,
			Stride:      stride,
			Likelihood:  lh,
			Impact:      im,
			Confidence:  a.Confidence,
			Keywords:    a.Keywords,
		}

		// Run through explainability pipeline
		if e.explainPipe != nil {
			e.explainPipe.Explain(&assumption)
		}

		// Apply fact protection: transform assumptions that contradict explicit architecture facts
		if e.archDesc != nil && len(e.archDesc.SecurityControls) > 0 {
			assumption = transformAssumptionForFacts(assumption, e.archDesc.SecurityControls)
		}

		result.Assumptions = append(result.Assumptions, assumption)

		switch assumption.Risk {
		case RiskCritical:
			result.CriticalCount++
		case RiskHigh:
			result.HighCount++
		case RiskMedium:
			result.MediumCount++
		case RiskLow:
			result.LowCount++
		}
	}

	// Process explicit assumptions from YAML/JSON
	if len(e.archDesc.ExplicitAssumptions) > 0 {
		explicitSet := e.processExplicitAssumptions(result.Assumptions)
		for _, ea := range explicitSet {
			result.Assumptions = append(result.Assumptions, ea)
			switch ea.Risk {
			case RiskCritical:
				result.CriticalCount++
			case RiskHigh:
				result.HighCount++
			case RiskMedium:
				result.MediumCount++
			case RiskLow:
				result.LowCount++
			}
		}
	}

	// Update totals to include explicit assumptions (post-dedup)
	result.TotalAssumptions = len(result.Assumptions)

	// Populate compliance from architecture description
	result.Compliance = e.buildComplianceOutput()

	// Apply security controls to all assumptions (native + explicit)
	if e.archDesc.SecurityControls != nil {
		for i := range result.Assumptions {
			result.Assumptions[i] = applySecurityControlVerification(result.Assumptions[i], e.archDesc.SecurityControls)
		}
	}

	// Re-map risk after security control verification so CONTRADICTED status is reflected
	for i, a := range result.Assumptions {
		if a.VerificationStatus == "CONTRADICTED" {
			result.Assumptions[i].Risk = RiskLow
		} else if a.VerificationStatus == "" || a.VerificationStatus == "UNKNOWN" {
			// Calibrate risk for known insecure patterns in assumptions
			textLower := strings.ToLower(a.Description)
			insecureEscalation := false
			insecureHigh := false

			criticalPatterns := []string{
				"shared admin", "default credential", "no encryption",
				"plaintext", "unencrypted", "no authentication",
				"no authorization", "single factor",
			}
			for _, p := range criticalPatterns {
				if strings.Contains(textLower, p) {
					insecureEscalation = true
					break
				}
			}

			highPatterns := []string{
				"flat network", "no logging", "no monitoring",
				"unencrypted backup", "weak cipher",
			}
			for _, p := range highPatterns {
				if strings.Contains(textLower, p) {
					insecureHigh = true
					break
				}
			}

			if insecureEscalation && a.Risk < RiskCritical {
				result.Assumptions[i].Risk = RiskCritical
			} else if insecureHigh && a.Risk < RiskHigh {
				result.Assumptions[i].Risk = RiskHigh
			}
		}
	}

	// Generate architecture-specific controls (not just generic templates)
	result.Controls = generateControls(result.Assumptions, e.archDesc.Components)
	if len(e.archDesc.SecurityControls) > 0 {
		result.Controls = enhanceControlsWithSecurityControls(result.Controls, e.archDesc.SecurityControls)
	}

	// Build evidence summary
	if e.explainPipe != nil {
		result.EvidenceSummary = e.explainPipe.BuildEvidenceSummary(result.Assumptions)
		confSummary := buildConfidenceSummary(result.Assumptions)
		result.ConfidenceSummary = confSummary
	}

	// Build validation summary if expected results are defined
	if e.archDesc != nil && len(e.archDesc.ExpectedResults) > 0 {
		result.Summary = e.buildValidationSummary(result)
	}

	return result
}

func mapRiskLevel(severity, verificationStatus string) RiskLevel {
	if verificationStatus == "CONTRADICTED" {
		return RiskLow
	}
	switch severity {
	case "CRITICAL":
		return RiskCritical
	case "HIGH":
		return RiskHigh
	case "MEDIUM":
		return RiskMedium
	default:
		return RiskMedium
	}
}

func riskToLikelihoodImpact(r RiskLevel) (int, int) {
	switch r {
	case RiskCritical:
		return 5, 5
	case RiskHigh:
		return 4, 4
	case RiskMedium:
		return 3, 3
	default:
		return 2, 2
	}
}

func (e *Engine) mapStrideDistribution(assumptions []Assumption) map[StrideCategory]int {
	dist := make(map[StrideCategory]int)
	for _, a := range assumptions {
		for _, s := range a.Stride {
			dist[s]++
		}
	}
	return dist
}

type controlTemplate struct {
	Category  string
	BaseDesc  string
	Rationale string
	STRIDE    []StrideCategory
	Priority  int
}

func controlTemplates() []controlTemplate {
	return []controlTemplate{
		{Category: "IDENTITY", BaseDesc: "Implement strong identity verification with MFA",
			Rationale: "Identity-related assumptions require robust authentication to prevent spoofing and unauthorized access.",
			STRIDE:    []StrideCategory{StrideSpoofing, StrideElevationPriv}, Priority: 1},
		{Category: "AUTHENTICATION", BaseDesc: "Enforce multi-factor authentication for all access",
			Rationale: "Authentication assumptions require verified identity to prevent credential-based attacks.",
			STRIDE:    []StrideCategory{StrideSpoofing, StrideElevationPriv}, Priority: 1},
		{Category: "AUTHORIZATION", BaseDesc: "Implement role-based access control with principle of least privilege",
			Rationale: "Authorization assumptions require strict access boundaries to prevent privilege escalation.",
			STRIDE:    []StrideCategory{StrideElevationPriv, StrideInfoDisclosure}, Priority: 1},
		{Category: "ACCESS", BaseDesc: "Enforce least-privilege access controls across all components",
			Rationale: "Access control assumptions limit blast radius and prevent lateral movement.",
			STRIDE:    []StrideCategory{StrideElevationPriv, StrideInfoDisclosure}, Priority: 1},
		{Category: "NETWORK", BaseDesc: "Implement network segmentation and encryption in transit",
			Rationale: "Network assumptions require boundary protection to prevent data exposure and DoS.",
			STRIDE:    []StrideCategory{StrideInfoDisclosure, StrideDenialOfService, StrideTampering}, Priority: 1},
		{Category: "ENCRYPTION", BaseDesc: "Implement encryption at rest and in transit for all sensitive data",
			Rationale: "Encryption assumptions protect confidentiality against data disclosure attacks.",
			STRIDE:    []StrideCategory{StrideInfoDisclosure}, Priority: 1},
		{Category: "CONFIGURATION", BaseDesc: "Use infrastructure-as-code with automated configuration validation",
			Rationale: "Configuration assumptions prevent tampering through misconfiguration and drift.",
			STRIDE:    []StrideCategory{StrideTampering}, Priority: 2},
		{Category: "DEPENDENCY", BaseDesc: "Implement dependency verification and supply chain security",
			Rationale: "Dependency assumptions protect against supply chain attacks and third-party compromise.",
			STRIDE:    []StrideCategory{StrideDenialOfService, StrideTampering}, Priority: 2},
		{Category: "PROCESS", BaseDesc: "Implement audit logging and process verification",
			Rationale: "Process assumptions ensure accountability and non-repudiation of security-relevant actions.",
			STRIDE:    []StrideCategory{StrideRepudiation, StrideTampering}, Priority: 2},
		{Category: "DATABASE", BaseDesc: "Implement database access controls and encryption",
			Rationale: "Database assumptions protect the confidentiality and integrity of stored data.",
			STRIDE:    []StrideCategory{StrideTampering, StrideInfoDisclosure}, Priority: 2},
		{Category: "LOGGING", BaseDesc: "Implement immutable audit logging with tamper detection",
			Rationale: "Logging assumptions prevent repudiation and enable forensic investigation.",
			STRIDE:    []StrideCategory{StrideRepudiation, StrideTampering}, Priority: 2},
		{Category: "BACKUP", BaseDesc: "Implement encrypted backup with tested restore procedures",
			Rationale: "Backup assumptions ensure data availability and recovery against ransomware and data loss.",
			STRIDE:    []StrideCategory{StrideInfoDisclosure, StrideDenialOfService}, Priority: 2},
		{Category: "SESSION", BaseDesc: "Implement secure session management with rotation and timeout",
			Rationale: "Session assumptions prevent session hijacking and credential reuse attacks.",
			STRIDE:    []StrideCategory{StrideSpoofing, StrideElevationPriv}, Priority: 2},
		{Category: "THIRD_PARTY", BaseDesc: "Implement third-party security assessment and monitoring",
			Rationale: "Third-party assumptions require vendor risk management to prevent supply chain attacks.",
			STRIDE:    []StrideCategory{StrideTampering, StrideInfoDisclosure}, Priority: 2},
		{Category: "DOCUMENTATION", BaseDesc: "Maintain accurate and version-controlled architecture documentation",
			Rationale: "Documentation assumptions ensure knowledge continuity and accurate threat modeling.",
			STRIDE:    []StrideCategory{StrideRepudiation}, Priority: 3},
		{Category: "GOVERNANCE", BaseDesc: "Establish security governance framework with regular reviews",
			Rationale: "Governance assumptions require oversight to maintain security posture over time.",
			STRIDE:    []StrideCategory{StrideRepudiation, StrideTampering}, Priority: 3},
	}
}

func generateControls(assumptions []Assumption, components []Component) []ControlDetail {
	templates := controlTemplates()
	tmplMap := make(map[string]*controlTemplate)
	for i := range templates {
		tmplMap[templates[i].Category] = &templates[i]
	}

	// Collect assumption IDs per category
	catAssumptions := make(map[string][]string)
	catStride := make(map[string]map[StrideCategory]bool)
	catSeen := make(map[string]bool)

	for _, a := range assumptions {
		cat := a.Category
		catAssumptions[cat] = append(catAssumptions[cat], a.ID)
		if catStride[cat] == nil {
			catStride[cat] = make(map[StrideCategory]bool)
		}
		for _, s := range a.Stride {
			catStride[cat][s] = true
		}
		catSeen[cat] = true
	}

	var controls []ControlDetail
	controlIdx := 0

	// Priority order
	for priority := 1; priority <= 3; priority++ {
		for _, tmpl := range templates {
			if tmpl.Priority != priority {
				continue
			}
			if !catSeen[tmpl.Category] && !hasMatchingStride(tmpl.STRIDE, catStride) {
				continue
			}
			if catSeen[tmpl.Category] && len(catAssumptions[tmpl.Category]) == 0 {
				continue
			}

			// If no direct category match but STRIDE matches, still include
			if !catSeen[tmpl.Category] {
				continue
			}

			controlIdx++
			ctrl := ControlDetail{
				ID:          fmt.Sprintf("CTRL-%03d", controlIdx),
				Description: tmpl.BaseDesc,
				Rationale:   tmpl.Rationale,
				Category:    tmpl.Category,
				Priority:    tmpl.Priority,
			}

			if ids, ok := catAssumptions[tmpl.Category]; ok {
				ctrl.MitigatedAssumptionIDs = ids
			}

			var strideList []StrideCategory
			for s := range catStride[tmpl.Category] {
				strideList = append(strideList, s)
			}
			if len(strideList) == 0 {
				strideList = tmpl.STRIDE
			}
			ctrl.MitigatedSTRIDE = strideList

			controls = append(controls, ctrl)
		}
	}

	// Generate architecture-specific controls based on actual components
	controls = generateArchitectureSpecificControls(controls, components, catAssumptions, catStride, &controlIdx)

	return controls
}

// generateArchitectureSpecificControls adds component-specific controls
// based on the actual architecture components (e.g., PHIDatabase, Auth0).
func generateArchitectureSpecificControls(controls []ControlDetail, components []Component, catAssumptions map[string][]string, catStride map[string]map[StrideCategory]bool, controlIdx *int) []ControlDetail {
	if len(components) == 0 {
		return controls
	}

	componentSpecific := map[string]struct {
		category  string
		desc      string
		rationale string
		stride    []StrideCategory
	}{
		"database": {
			category:  "DATABASE",
			desc:      "Implement database-specific encryption, access logging, and connection pooling safeguards",
			rationale: "Database components require tailored controls for data integrity and confidentiality.",
			stride:    []StrideCategory{StrideTampering, StrideInfoDisclosure},
		},
		"api_gateway": {
			category:  "NETWORK",
			desc:      "Implement API Gateway-specific rate limiting, request validation, and TLS termination policies",
			rationale: "API Gateway components require targeted controls for ingress security and availability.",
			stride:    []StrideCategory{StrideDenialOfService, StrideInfoDisclosure},
		},
		"identity_provider": {
			category:  "IDENTITY",
			desc:      "Implement identity-provider-specific MFA enforcement, session hardening, and breach detection",
			rationale: "Identity provider components require specialized controls for authentication resilience.",
			stride:    []StrideCategory{StrideSpoofing, StrideElevationPriv},
		},
		"web_application": {
			category:  "ACCESS",
			desc:      "Implement web-application-specific input validation, CSRF protection, and secure session handling",
			rationale: "Web application components require application-layer controls for common attack vectors.",
			stride:    []StrideCategory{StrideTampering, StrideElevationPriv},
		},
		"external_service": {
			category:  "THIRD_PARTY",
			desc:      "Implement external-service-specific vendor risk monitoring, data flow audits, and contractual controls",
			rationale: "External service components require supply-chain and data-exposure controls.",
			stride:    []StrideCategory{StrideTampering, StrideInfoDisclosure},
		},
		"admin_tool": {
			category:  "ACCESS",
			desc:      "Implement admin-tool-specific privileged access monitoring, command logging, and just-in-time access",
			rationale: "Admin tools require elevated controls due to high privilege usage.",
			stride:    []StrideCategory{StrideElevationPriv, StrideRepudiation},
		},
		"storage_service": {
			category:  "BACKUP",
			desc:      "Implement storage-service-specific encrypted backup, cross-region replication, and restore validation",
			rationale: "Storage service components require resilience and recovery controls.",
			stride:    []StrideCategory{StrideInfoDisclosure, StrideDenialOfService},
		},
		"encryption_service": {
			category:  "ENCRYPTION",
			desc:      "Implement encryption-service-specific key rotation, access auditing, and deletion protection",
			rationale: "Encryption service components require key-management lifecycle controls.",
			stride:    []StrideCategory{StrideInfoDisclosure, StrideTampering},
		},
		"logging_service": {
			category:  "LOGGING",
			desc:      "Implement logging-service-specific immutability, tamper detection, and retention compliance",
			rationale: "Logging service components require audit integrity and non-repudiation controls.",
			stride:    []StrideCategory{StrideRepudiation, StrideTampering},
		},
	}

	seenCategory := make(map[string]bool)
	for _, c := range controls {
		seenCategory[c.Category] = true
	}

	for _, comp := range components {
		compLower := strings.ToLower(comp.Label)
		compType := strings.ToLower(comp.Label)
		for key, ctrl := range componentSpecific {
			if strings.Contains(compLower, key) || strings.Contains(compType, key) {
				if seenCategory[ctrl.category] {
					continue
				}
				seenCategory[ctrl.category] = true
				*controlIdx++
				newCtrl := ControlDetail{
					ID:              fmt.Sprintf("CTRL-%03d", *controlIdx),
					Description:     fmt.Sprintf("[%s] %s", comp.Label, ctrl.desc),
					Rationale:       ctrl.rationale,
					Category:        ctrl.category,
					Priority:        1,
					MitigatedSTRIDE: ctrl.stride,
				}
				if ids, ok := catAssumptions[ctrl.category]; ok {
					newCtrl.MitigatedAssumptionIDs = ids
				}
				controls = append(controls, newCtrl)
				break
			}
		}
	}

	return controls
}

// enhanceControlsWithSecurityControls extends generated controls with specifics from YAML security_controls.
// This is called from generateControls with the archDesc.SecurityControls map.
func enhanceControlsWithSecurityControls(controls []ControlDetail, securityControls map[string][]string) []ControlDetail {
	if len(securityControls) == 0 {
		return controls
	}

	// Category mapping from security_controls categories to template categories
	catMap := map[string]string{
		"authentication": "AUTHENTICATION",
		"authorization":  "AUTHORIZATION",
		"encryption":     "ENCRYPTION",
		"logging":        "LOGGING",
		"backup":         "BACKUP",
		"network":        "NETWORK",
		"monitoring":     "PROCESS",
		"third_party":    "THIRD_PARTY",
		"session":        "SESSION",
	}

	for scCategory, scControls := range securityControls {
		tmplCat := catMap[scCategory]
		if tmplCat == "" {
			continue
		}

		// Find matching control
		found := false
		for i := range controls {
			if controls[i].Category == tmplCat {
				// Enrich description with specific controls
				if len(scControls) > 0 {
					controls[i].Description = fmt.Sprintf("%s: %s", controls[i].Description, strings.Join(scControls, ", "))
				}
				found = true
				break
			}
		}

		if !found {
			// Create a new control for this security category
			ctrlID := fmt.Sprintf("CTRL-%03d", len(controls)+1)
			ctrl := ControlDetail{
				ID:          ctrlID,
				Description: fmt.Sprintf("Implement %s controls: %s", scCategory, strings.Join(scControls, ", ")),
				Rationale:   fmt.Sprintf("Architecture defines %s controls that must be verified.", scCategory),
				Category:    tmplCat,
				Priority:    1,
			}
			controls = append(controls, ctrl)
		}
	}

	return controls
}

func hasMatchingStride(target []StrideCategory, catStride map[string]map[StrideCategory]bool) bool {
	for _, cat := range catStride {
		for _, t := range target {
			if cat[t] {
				return true
			}
		}
	}
	return false
}

// processExplicitAssumptions processes explicit assumptions from YAML/JSON,
// deduplicates against existing assumptions, and enriches them.
// Deep deduplication: strips bullet prefixes and merges source metadata.
func (e *Engine) processExplicitAssumptions(existing []Assumption) []Assumption {
	existingMap := make(map[string]int) // normalized -> index in existing
	for i, a := range existing {
		existingMap[normalizeText(a.Description)] = i
	}

	var result []Assumption
	nextID := 1

	for i, raw := range e.archDesc.ExplicitAssumptions {
		normalized := normalizeText(raw)
		if normalized == "" {
			continue
		}
		if existingIdx, ok := existingMap[normalized]; ok {
			if existingIdx >= 0 {
				// Deep dedup: merge source metadata into existing assumption
				existing[existingIdx] = mergeSourceMetadata(existing[existingIdx], raw, "assumptions", i, e.archDesc.Name)
			}
			continue
		}
		existingMap[normalized] = -1 // mark as processed

		atype := classifyExplicitAssumption(raw)
		keywords := extractKeywords(raw)
		component := extractComponent(keywords, raw)

		id := fmt.Sprintf("ASM-%03d", nextID)
		nextID++

		risk := e.assessExplicitRisk(raw, atype)
		lh, im := riskToLikelihoodImpact(risk)
		stride := e.strideEngine.MapAssumption(string(atype), raw, keywords)
		confidence := computeExplicitConfidence(raw, atype, e.archDesc.SecurityControls)

		assumption := Assumption{
			ID:                 id,
			Description:        raw,
			Component:          component,
			Category:           string(atype),
			Risk:               risk,
			Stride:             stride,
			Likelihood:         lh,
			Impact:             im,
			Confidence:         confidence,
			Keywords:           keywords,
			SourceType:         "explicit",
			SourceSection:      "assumptions",
			SourceIndex:        i,
			SourceFile:         e.archDesc.Name,
			VerificationStatus: "UNKNOWN",
		}

		// Wire security controls into verification
		if e.archDesc.SecurityControls != nil {
			assumption = applySecurityControlVerification(assumption, e.archDesc.SecurityControls)
		}

		if e.explainPipe != nil {
			e.explainPipe.Explain(&assumption)
		}

		result = append(result, assumption)
	}
	return result
}

// mergeSourceMetadata merges source metadata from a duplicate explicit assumption
// into an existing assumption, preserving the original source.
func mergeSourceMetadata(a Assumption, raw, section string, index int, sourceFile string) Assumption {
	if a.SourceType == "" {
		a.SourceType = "merged"
	}
	if a.SourceSection == "" {
		a.SourceSection = section
	} else {
		a.SourceSection += "," + section
	}
	if a.SourceFile == "" {
		a.SourceFile = sourceFile
	} else if a.SourceFile != sourceFile {
		a.SourceFile += "," + sourceFile
	}
	if a.SourceIndex == 0 {
		a.SourceIndex = index
	}
	return a
}

// computeExplicitConfidence computes confidence for explicit assumptions.
// Base 0.75; boosted to 0.80+ if supported by declared security controls.
func computeExplicitConfidence(text string, atype models.AssumptionType, securityControls map[string][]string) float64 {
	confidence := 0.75
	lower := strings.ToLower(text)

	// Check if assumption keywords match security controls
	if securityControls != nil {
		categoryBoosts := map[string][]string{
			"authentication": {"mfa", "password", "session", "sso", "oauth", "oidc", "auth0", "login", "identity", "authentication"},
			"authorization":  {"rbac", "access", "authorized", "permission", "least privilege", "admin", "acl"},
			"encryption":     {"encrypt", "tls", "ssl", "cipher", "key", "kms", "aes", "cryptographic"},
			"logging":        {"log", "audit", "monitor", "alert", "detect", "immutable"},
			"backup":         {"backup", "restore", "recovery", "replicate"},
			"network":        {"network", "subnet", "firewall", "segment", "private", "internet", "rate limit"},
			"monitoring":     {"monitor", "health", "availability", "anomaly", "detect"},
			"third_party":    {"third-party", "third party", "vendor", "external", "supplier", "de-identified"},
			"session":        {"session", "token", "jwt", "expire", "rotate"},
		}
		for category, kws := range categoryBoosts {
			controls, hasControls := securityControls[category]
			if !hasControls || len(controls) == 0 {
				continue
			}
			for _, kw := range kws {
				if strings.Contains(lower, kw) {
					confidence = 0.85
					break
				}
			}
			if confidence >= 0.85 {
				break
			}
		}
	}

	if confidence > 0.95 {
		confidence = 0.95
	}
	return confidence
}

// applySecurityControlVerification checks if an assumption is covered by
// declared security controls and marks it PARTIALLY_VERIFIED when applicable.
var insecureControlExact = map[string]bool{
	"none":     true,
	"disabled": true,
}

var insecureControlContains = []string{
	"none", "disabled", "plaintext", "flatnetwork", "directinternet",
	"basic", "shared", "admin_by_default",
}

func hasInsecureControl(controls []string) bool {
	for _, ctrl := range controls {
		ctrlKey := strings.ToLower(strings.ReplaceAll(ctrl, "_", ""))
		if insecureControlExact[ctrlKey] {
			return true
		}
		for _, prefix := range insecureControlContains {
			if strings.Contains(ctrlKey, strings.ReplaceAll(prefix, "_", "")) {
				return true
			}
		}
	}
	return false
}

// normalizedControlName maps common assumption keywords to canonical control names
// so that, e.g., "multi-factor" in text matches "Admin_MFA" in security controls.
// normalizedControlName maps common assumption keywords to canonical control names
// so that, e.g., "multi-factor" in text matches "Admin_MFA" in security controls.
var normalizedControlName = map[string]string{
	"mfa":             "mfa",
	"multi-factor":    "mfa",
	"multi factor":    "mfa",
	"two-factor":      "mfa",
	"two factor":      "mfa",
	"2fa":             "mfa",
	"totp":            "mfa",
	"rbac":            "rbac",
	"role-based":      "rbac",
	"role based":      "rbac",
	"abac":            "abac",
	"least privilege": "least_privilege",
	"least-privilege": "least_privilege",
	"tls":             "tls",
	"https":           "https",
	"ssl":             "tls",
	"aes-256":         "aes256",
	"aes 256":         "aes256",
	"aes256":          "aes256",
	"kms":             "kms",
	"encrypted":       "encrypted_control",
	"encryption":      "encrypted_control",
	"audit":           "audit_logging",
	"siem":            "siem",
	"alert":           "alerting",
	"backup":          "backup",
	"restore":         "restore_testing",
}

// controlCategoryConcept maps control categories to normalized concept keywords
// that appear in assumption text. For example, an assumption saying "encrypted"
// matches any control in the "encryption" category, even if no specific control
// name (TLS, AES256, KMS) appears in the text.
var controlCategoryConcept = map[string][]string{
	"authentication": {"mfa", "authentication", "login", "identity", "auth", "password"},
	"authorization":  {"authorization", "rbac", "access control", "permission", "role", "privilege"},
	"encryption":     {"encrypted", "encryption", "tls", "https", "ssl", "cipher", "crypto"},
	"backup":         {"backup", "restore", "recovery"},
	"monitoring":     {"audit", "log", "monitor", "alert", "siem", "detection"},
	"logging":        {"audit", "log", "monitor"},
	"network":        {"network", "firewall", "segment", "vpc", "subnet"},
}

func applySecurityControlVerification(a Assumption, securityControls map[string][]string) Assumption {
	lower := strings.ToLower(a.Description)
	categoryMap := map[string][]string{
		"IDENTITY":                 {"authentication", "authorization", "session"},
		"ACCESS":                   {"authorization", "authentication"},
		"CONFIGURATION":            {"encryption", "logging", "backup", "network"},
		"NETWORK":                  {"network"},
		"PROCESS":                  {"monitoring", "logging"},
		"GOVERNANCE":               {"monitoring", "logging"},
		"DEPENDENCY":               {"third_party"},
		"DOCUMENTATION":            {},
		"AUTHENTICATION":           {"authentication"},
		"AUTHORIZATION":            {"authorization"},
		"ENCRYPTION":               {"encryption"},
		"BACKUPS":                  {"backup"},
		"LOGGING":                  {"logging"},
		"MONITORING":               {"monitoring"},
		"KEYMANAGEMENT":            {"encryption"},
		"AUDITABILITY":             {"monitoring", "logging"},
		"PRIVILEGEMANAGEMENT":      {"authorization", "authentication"},
		"SESSIONSECURITY":          {"authentication"},
		"DATAPROTECTION":           {"encryption", "backup"},
		"NETWORKSEGMENTATION":      {"network"},
		"DISASTERRECOVERY":         {"backup"},
		"THIRDPARTYRISK":           {"third_party"},
		"VENDORRISK":               {"third_party"},
		"APISECURITY":              {"authentication", "authorization", "encryption"},
		"SECRETSACCESS":            {"authentication"},
		"SECRETS_ACCESS":           {"authentication"},
		"DATARETENTION":            {"backup"},
		"OBJECTLEVELAUTHORIZATION": {"authorization"},
		"TRUSTBOUNDARIES":          {"authentication", "authorization", "encryption", "monitoring"},
		"IDENTITY_TO_APPLICATION":  {"authentication", "authorization", "encryption"},
		"APPLICATION_TO_DATA":      {"encryption", "backup", "authorization"},
		"THIRD_PARTY_TO_INTERNAL":  {"third_party", "authentication"},
		"THIRD_PARTY_TO_DATA":      {"third_party", "encryption", "backup"},
		"COMPLIANCE":               {"monitoring", "logging"},
		"PRIVACY":                  {"encryption", "backup"},
	}

	categories, ok := categoryMap[strings.ToUpper(a.Category)]
	if !ok {
		categories = []string{}
	}

	// Phase 1: Check for negative/insecure controls that CONTRADICT the assumption.
	// When a mapped security control category has explicitly insecure values
	// (e.g. encryption: [None], network: [Flat_Network]), any assumption in
	// that category is contradicted regardless of its text content.
	insecureFound := false
	hasSecureControls := false
	for _, cat := range categories {
		controls, hasControls := securityControls[cat]
		if !hasControls || len(controls) == 0 {
			continue
		}
		if hasInsecureControl(controls) {
			insecureFound = true
		} else {
			hasSecureControls = true
		}
	}

	if insecureFound {
		a.VerificationStatus = "CONTRADICTED"
		if a.Confidence < 0.90 {
			a.Confidence = 0.90
		}
		a.Rationale = fmt.Sprintf("Contradicted by declared security controls: insecure configuration found for category %s", strings.Join(categories, ", "))
		return a
	}

	// Phase 2: Positive verification — if secure controls exist for the assumption's domain.
	// VERIFIED when the assumption text explicitly mentions a matching control or
	// a normalized variant. PARTIALLY_VERIFIED when controls exist but text does not
	// reference them.
	if hasSecureControls {
		matched := false
		var matchedControl string
		for _, cat := range categories {
			controls, hasControls := securityControls[cat]
			if !hasControls || len(controls) == 0 {
				continue
			}
			for _, ctrl := range controls {
				ctrlLower := strings.ToLower(strings.ReplaceAll(ctrl, "_", " "))
				// Direct text match: does the assumption contain the control name?
				if strings.Contains(lower, ctrlLower) || strings.Contains(lower, strings.ReplaceAll(ctrlLower, " ", "")) {
					matched = true
					matchedControl = ctrl
					break
				}
				// Normalized name match: check if any keyword in the text maps to this control
				for textKW, canonical := range normalizedControlName {
					canonicalLower := strings.ToLower(strings.ReplaceAll(canonical, "_", " "))
					if strings.Contains(lower, textKW) && (canonicalLower == ctrlLower || canonicalLower == strings.ToLower(ctrl)) {
						matched = true
						matchedControl = ctrl
						break
					}
				}
				if matched {
					break
				}
			}
			if matched {
				break
			}
		}

		if matched {
			a.VerificationStatus = "VERIFIED"
			if a.Confidence < 0.85 {
				a.Confidence = 0.85
			}
			a.Rationale = fmt.Sprintf("Verified: assumption matches declared control '%s'. Secure design confirmed for category %s.", matchedControl, strings.Join(categories, ", "))
		} else {
			// Check concept-based match: does the text mention a concept that maps to a control category?
			conceptMatched := false
			var matchedConcept string
			for _, cat := range categories {
				concepts, hasConcepts := controlCategoryConcept[cat]
				if !hasConcepts {
					continue
				}
				for _, concept := range concepts {
					if strings.Contains(lower, concept) {
						conceptMatched = true
						matchedConcept = concept
						break
					}
				}
				if conceptMatched {
					break
				}
			}
			if conceptMatched {
				a.VerificationStatus = "VERIFIED"
				if a.Confidence < 0.75 {
					a.Confidence = 0.75
				}
				a.Rationale = fmt.Sprintf("Verified: assumption mentions '%s' which is covered by declared controls in category %s.", matchedConcept, strings.Join(categories, ", "))
			} else {
				a.VerificationStatus = "PARTIALLY_VERIFIED"
				if a.Confidence < 0.70 {
					a.Confidence = 0.70
				}
				a.Rationale = fmt.Sprintf("Partially verified: security controls exist for category %s but assumption text does not explicitly reference them.", strings.Join(categories, ", "))
			}
		}
	}
	return a
}

// classifyExplicitAssumption classifies an explicit assumption into a type.
func classifyExplicitAssumption(text string) models.AssumptionType {
	lower := strings.ToLower(text)

	// Identity/authentication keywords
	identityKws := []string{"mfa", "multi-factor", "authentication", "password", "credential", "sso", "oauth", "oidc", "auth0", "login", "identity"}
	for _, kw := range identityKws {
		if strings.Contains(lower, kw) {
			return models.AssumptionTypeIDENTITY
		}
	}

	// Access/authorization keywords
	accessKws := []string{"access", "authorized", "permission", "rbac", "acl", "restricted", "least privilege", "admin"}
	for _, kw := range accessKws {
		if strings.Contains(lower, kw) {
			return models.AssumptionTypeACCESS
		}
	}

	// Encryption keywords
	encKws := []string{"encrypt", "tls", "ssl", "cipher", "key", "kms", "cryptographic", "aes"}
	for _, kw := range encKws {
		if strings.Contains(lower, kw) {
			return models.AssumptionTypeCONFIGURATION
		}
	}

	// Network keywords
	netKws := []string{"network", "subnet", "firewall", "segment", "tls termination", "private", "internet"}
	for _, kw := range netKws {
		if strings.Contains(lower, kw) {
			return models.AssumptionTypeNETWORK
		}
	}

	// Logging/monitoring keywords
	logKws := []string{"log", "audit", "monitor", "alert", "detect"}
	for _, kw := range logKws {
		if strings.Contains(lower, kw) {
			return models.AssumptionTypeCONFIGURATION
		}
	}

	// Backup keywords
	backupKws := []string{"backup", "restore", "recovery", "replicate"}
	for _, kw := range backupKws {
		if strings.Contains(lower, kw) {
			return models.AssumptionTypeCONFIGURATION
		}
	}

	// Session keywords
	sessionKws := []string{"session", "token", "jwt", "expire", "rotate"}
	for _, kw := range sessionKws {
		if strings.Contains(lower, kw) {
			return models.AssumptionTypeIDENTITY
		}
	}

	// Third-party keywords
	thirdKws := []string{"third-party", "third party", "vendor", "external", "supplier"}
	for _, kw := range thirdKws {
		if strings.Contains(lower, kw) {
			return models.AssumptionTypeDEPENDENCY
		}
	}

	// Process/governance keywords
	procKws := []string{"incident", "response", "breach", "procedure", "policy", "review", "assessment"}
	for _, kw := range procKws {
		if strings.Contains(lower, kw) {
			return models.AssumptionTypePROCESS
		}
	}

	return models.AssumptionTypeGOVERNANCE
}

// assessExplicitRisk determines the initial risk level for an explicit assumption.
func (e *Engine) assessExplicitRisk(text string, atype models.AssumptionType) RiskLevel {
	lower := strings.ToLower(text)
	score := 0

	// PHI/healthcare keywords boost risk significantly
	phiKws := []string{"phi", "health", "hipaa", "patient", "medical", "phi data", "protected health", "clinic", "hospital", "electronic health record", "ehr"}
	for _, kw := range phiKws {
		if strings.Contains(lower, kw) {
			score += 4
			break
		}
	}

	// High-severity keywords
	highKws := []string{"critical", "compromise", "breach", "unauthorized", "restricted", "immutable", "forbidden", "mandatory"}
	for _, kw := range highKws {
		if strings.Contains(lower, kw) {
			score += 3
			break
		}
	}

	// Medium-severity keywords
	medKws := []string{"encrypt", "kms", "access", "authenticate", "mfa", "token", "key", "audit", "backup", "monitor", "rate limit", "session", "tls", "ssl", "certificate"}
	for _, kw := range medKws {
		if strings.Contains(lower, kw) {
			score += 1
			break
		}
	}

	// Type-based boost
	switch atype {
	case models.AssumptionTypeACCESS, models.AssumptionTypeIDENTITY:
		score += 2
	case models.AssumptionTypeNETWORK, models.AssumptionTypeCONFIGURATION:
		score += 1
	}

	switch {
	case score >= 6:
		return RiskCritical
	case score >= 4:
		return RiskHigh
	case score >= 2:
		return RiskMedium
	default:
		return RiskLow
	}
}

// complianceFrameworkDetails maps known frameworks to specific areas covered.
var complianceFrameworkDetails = map[string]struct {
	Label   string
	Areas   []string
	Control string
}{
	"HIPAA": {
		Label:   "HIPAA (Health Insurance Portability and Accountability Act)",
		Areas:   []string{"PHI access controls (164.312(a))", "Encryption of PHI at rest and in transit (164.312(a)(2)(iv))", "Audit controls for PHI access (164.312(b))", "Integrity controls (164.312(c)(1))", "Emergency access procedure (164.312(a)(2)(ii))", "Automatic logoff (164.312(a)(2)(iii))"},
		Control: "HIPAA Security Rule — Administrative, Physical, and Technical Safeguards",
	},
	"SOC2": {
		Label:   "SOC 2 (Service Organization Control 2)",
		Areas:   []string{"Security —保护 against unauthorized access", "Availability —监控 and capacity planning", "Processing Integrity — data processing accuracy", "Confidentiality — data classification and handling", "Privacy — PII collection and use"},
		Control: "SOC 2 Trust Services Criteria — Security, Availability, Processing Integrity, Confidentiality, Privacy",
	},
	"ISO27001": {
		Label:   "ISO/IEC 27001 (Information Security Management)",
		Areas:   []string{"A.9 Access control — identity and authorization", "A.10 Cryptography — encryption key management", "A.12 Operations security — malware protection, backup, monitoring", "A.16 Incident management — reporting and response", "A.18 Compliance — regulatory and contractual obligations"},
		Control: "ISO/IEC 27001:2022 Annex A controls",
	},
	"PCIDSS": {
		Label:   "PCI DSS (Payment Card Industry Data Security Standard)",
		Areas:   []string{"Requirement 3 — Protect stored cardholder data", "Requirement 4 — Encrypt transmission of cardholder data", "Requirement 7 — Restrict access to cardholder data", "Requirement 10 — Track and monitor access to data", "Requirement 12 — Maintain information security policy"},
		Control: "PCI DSS v4.0 requirements",
	},
	"GDPR": {
		Label:   "GDPR (General Data Protection Regulation)",
		Areas:   []string{"Art. 5 — Principles of processing (integrity and confidentiality)", "Art. 25 — Data protection by design and default", "Art. 32 — Security of processing", "Art. 33 — Breach notification", "Art. 35 — Data protection impact assessment"},
		Control: "GDPR Chapter IV — Controller and Processor",
	},
	"FedRAMP": {
		Label:   "FedRAMP (Federal Risk and Authorization Management Program)",
		Areas:   []string{"AC — Access Control", "AU — Audit and Accountability", "IA — Identification and Authentication", "SC — System and Communications Protection", "SI — System and Information Integrity"},
		Control: "NIST SP 800-53 rev 5 controls",
	},
}

// buildComplianceOutput builds the compliance section from architecture description.
func (e *Engine) buildComplianceOutput() []string {
	if e.archDesc == nil || len(e.archDesc.Compliance) == 0 {
		return []string{
			"ASF analysis completed — see gap analysis for compliance mapping",
		}
	}

	compliance := e.archDesc.Compliance
	output := make([]string, 0, len(compliance)*6+4)
	output = append(output, "Compliance frameworks identified in architecture definition:")
	for _, c := range compliance {
		if details, ok := complianceFrameworkDetails[c]; ok {
			output = append(output, fmt.Sprintf(""))
			output = append(output, fmt.Sprintf("--- %s ---", details.Label))
			output = append(output, fmt.Sprintf("Framework: %s", details.Control))
			output = append(output, fmt.Sprintf("Relevant areas:"))
			for _, area := range details.Areas {
				output = append(output, fmt.Sprintf("  - %s", area))
			}
		} else {
			output = append(output, fmt.Sprintf("- %s (custom framework — review specific requirements)", c))
		}
	}
	output = append(output, "")
	output = append(output, "Review gap analysis for detailed compliance mapping against these requirements.")
	return output
}

// buildValidationSummary validates expected results against actual analysis results.
func (e *Engine) buildValidationSummary(result *AnalysisResult) string {
	expected := e.archDesc.ExpectedResults
	var violations []string
	var met []string

	// Check minimum assumptions count
	if minAssump, ok := expected["minimum_assumptions"]; ok {
		if min, ok := toFloat(minAssump); ok {
			if result.TotalAssumptions < int(min) {
				violations = append(violations, fmt.Sprintf("expected ≥%.0f assumptions, got %d", min, result.TotalAssumptions))
			} else {
				met = append(met, fmt.Sprintf("minimum assumptions met: %d (≥%.0f)", result.TotalAssumptions, min))
			}
		}
	}

	// Check minimum critical count
	if minCrit, ok := expected["minimum_critical"]; ok {
		if min, ok := toFloat(minCrit); ok {
			if result.CriticalCount < int(min) {
				violations = append(violations, fmt.Sprintf("expected ≥%.0f critical findings, got %d", min, result.CriticalCount))
			} else {
				met = append(met, fmt.Sprintf("minimum critical findings met: %d (≥%.0f)", result.CriticalCount, min))
			}
		}
	}

	// Check minimum high count
	if minHigh, ok := expected["minimum_high"]; ok {
		if min, ok := toFloat(minHigh); ok {
			if result.HighCount < int(min) {
				violations = append(violations, fmt.Sprintf("expected ≥%.0f high findings, got %d", min, result.HighCount))
			} else {
				met = append(met, fmt.Sprintf("minimum high findings met: %d (≥%.0f)", result.HighCount, min))
			}
		}
	}

	// Check expected STRIDE categories
	if expStride, ok := expected["expected_stride_categories"]; ok {
		if strideList, ok := expStride.([]interface{}); ok {
			present := make(map[string]bool)
			for cat, count := range result.StrideDistribution {
				if count > 0 {
					present[string(cat)] = true
				}
			}
			for _, s := range strideList {
				str := fmt.Sprintf("%v", s)
				if !present[str] {
					violations = append(violations, fmt.Sprintf("expected STRIDE category %q not found", str))
				} else {
					met = append(met, fmt.Sprintf("STRIDE category %q present", str))
				}
			}
		}
	}

	// Build summary
	var b strings.Builder
	b.WriteString(fmt.Sprintf("ASF processed architecture"))
	if e.archDesc != nil && e.archDesc.Name != "" {
		b.WriteString(fmt.Sprintf(" %q", e.archDesc.Name))
	}
	b.WriteString(fmt.Sprintf(" and found %d assumptions (%d critical, %d high, %d medium, %d low). ",
		result.TotalAssumptions, result.CriticalCount, result.HighCount, result.MediumCount, result.LowCount))

	if len(violations) > 0 {
		b.WriteString(fmt.Sprintf("Validation: %d violation(s) found — %s.",
			len(violations), strings.Join(violations, "; ")))
	} else {
		b.WriteString("Validation: all expected criteria met.")
	}
	if len(met) > 0 {
		b.WriteString(fmt.Sprintf(" Criteria met: %s.", strings.Join(met, "; ")))
	}

	return b.String()
}

// toFloat converts an interface{} to float64 for numeric comparisons.
func toFloat(v interface{}) (float64, bool) {
	switch n := v.(type) {
	case float64:
		return n, true
	case int:
		return float64(n), true
	case int32:
		return float64(n), true
	case int64:
		return float64(n), true
	case uint:
		return float64(n), true
	case uint64:
		return float64(n), true
	default:
		return 0, false
	}
}

// normalizeText normalizes text for deduplication comparison.
// Strips bullet markers, trailing periods, and collapses whitespace.
func normalizeText(text string) string {
	s := strings.ToLower(text)
	// Strip leading bullet markers: "- ", "* ", "  - ", etc.
	s = regexp.MustCompile(`^[\s\-*•·]+`).ReplaceAllString(s, "")
	// Strip trailing period
	s = strings.TrimSuffix(s, ".")
	// Collapse whitespace and strip periods from each token
	parts := strings.Fields(s)
	cleaned := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSuffix(p, ".")
		if p != "" {
			cleaned = append(cleaned, p)
		}
	}
	return strings.Join(cleaned, " ")
}

// extractKeywords extracts significant keywords from text for explicit assumptions.
func extractKeywords(text string) []string {
	re := regexp.MustCompile(`\b[a-zA-Z]{3,}\b`)
	words := re.FindAllString(strings.ToLower(text), -1)
	var result []string
	stop := map[string]bool{
		"the": true, "and": true, "for": true, "are": true, "but": true,
		"not": true, "you": true, "all": true, "can": true, "has": true,
		"have": true, "may": true, "must": true, "shall": true, "should": true,
		"will": true, "with": true, "from": true, "that": true, "this": true,
		"each": true, "every": true, "than": true, "then": true, "just": true,
		"been": true, "were": true, "was": true, "its": true, "also": true,
		"per": true, "via": true, "is": true, "to": true, "in": true,
		"of": true, "on": true, "at": true, "by": true, "as": true,
		"an": true, "or": true,
	}
	for _, w := range words {
		if !stop[w] {
			result = append(result, w)
		}
	}
	return result
}

func fileBase(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '/' || path[i] == '\\' {
			return path[i+1:]
		}
	}
	return path
}

func cleanAssumptionText(text string) string {
	colonIdx := strings.Index(text, ": ")
	if colonIdx >= 0 && colonIdx < 80 {
		after := strings.TrimSpace(text[colonIdx+2:])
		parts := strings.SplitN(after, "\n", 2)
		firstLine := strings.TrimSpace(parts[0])
		if len(firstLine) > 5 {
			return firstLine
		}
	}
	lines := strings.SplitN(text, "\n", 2)
	firstLine := strings.TrimSpace(lines[0])
	if len(firstLine) > 10 {
		return firstLine
	}
	if len(text) > 120 {
		return text[:120] + "..."
	}
	return text
}

// factPolarityRule describes a known security control pattern and how to
// transform an assumption that contradicts it.
type factPolarityRule struct {
	category     string   // security control category, lowercased
	negValues    []string // negative indicator values
	triggerWords []string // words in assumption description that trigger protection
	transformed  string   // replacement description prefix for the transformed assumption
}

var defaultPolarityRules = []factPolarityRule{
	{category: "encryption", negValues: []string{"none", "disabled", "false"},
		triggerWords: []string{"encrypt", "tls", "ssl", "https"},
		transformed:  "Plaintext communication is expected; compensating controls or accepted risk exists"},
	{category: "authentication", negValues: []string{"basic", "none", "disabled", "single_factor", "password_only"},
		triggerWords: []string{"mfa", "multi-factor", "two-factor", "2fa", "strong authentication", "passwordless"},
		transformed:  "Single-factor or weak authentication is in use; compensating controls or accepted risk exists"},
	{category: "authorization", negValues: []string{"disabled", "none"},
		triggerWords: []string{"authorization", "access control", "least privilege", "rbac"},
		transformed:  "Authorization controls are disabled; risk of unauthorized access accepted"},
	{category: "network", negValues: []string{"flat", "flatnetwork", "open", "none"},
		triggerWords: []string{"network segmentation", "firewall", "isolated", "network control", "vlan"},
		transformed:  "Flat network topology in use; lateral movement risk accepted"},
	{category: "backup", negValues: []string{"unencrypted", "none", "disabled", "false"},
		triggerWords: []string{"backup encryption", "encrypted backup", "backup security"},
		transformed:  "Backups are unencrypted; data exposure risk accepted"},
	{category: "monitoring", negValues: []string{"disabled", "none", "false"},
		triggerWords: []string{"monitoring", "alerting", "logging", "siem", "audit"},
		transformed:  "Monitoring and alerting are disabled; detection gaps accepted"},
}

// transformAssumptionForFacts checks an assumption against the architecture's
// security controls. If the assumption asserts a security property that is
// explicitly negated by a control, the assumption description is replaced with
// a risk-aware statement that respects the architect's declared fact.
func transformAssumptionForFacts(a Assumption, sc map[string][]string) Assumption {
	descLower := strings.ToLower(a.Description)

	for _, rule := range defaultPolarityRules {
		// Check if the assumption text triggers this rule
		triggered := false
		for _, tw := range rule.triggerWords {
			if strings.Contains(descLower, tw) {
				triggered = true
				break
			}
		}
		if !triggered {
			continue
		}

		// Check if the architecture has a security control for this category
		catLower := strings.ToLower(rule.category)
		values, ok := sc[catLower]
		if !ok {
			// Also check the original case (YAML keys may vary)
			values, ok = sc[rule.category]
			if !ok {
				continue
			}
		}

		// Check if any control value is a negative/disabled value
		hasNegative := false
		for _, v := range values {
			vLower := strings.ToLower(v)
			for _, nv := range rule.negValues {
				if vLower == nv || strings.ReplaceAll(vLower, "_", "") == nv {
					hasNegative = true
					break
				}
			}
			if hasNegative {
				break
			}
		}
		if !hasNegative {
			continue
		}

		// Transform the assumption: replace description with risk-aware statement
		a.Description = rule.transformed + " [" + rule.category + ": " + strings.Join(values, ", ") + "]"
		a.VerificationStatus = "CONTRADICTED"
		if a.Confidence < 0.90 {
			a.Confidence = 0.90
		}
		break
	}

	return a
}

func buildConfidenceSummary(assumptions []Assumption) string {
	total := len(assumptions)
	if total == 0 {
		return "no assumptions to evaluate"
	}
	var sum float64
	for _, a := range assumptions {
		sum += a.Confidence
	}
	avg := sum / float64(total)
	high := 0
	for _, a := range assumptions {
		if a.Confidence >= 0.7 {
			high++
		}
	}
	return fmt.Sprintf("average confidence %.0f%% across %d assumptions (%d high-confidence)", avg*100, total, high)
}

func extractComponent(keywords []string, text string) string {
	// Look for capitalized words that look like component names
	// Generated text uses PascalCase/TitleCase names (WebApp, Auth0, PHIDatabase, etc.)
	words := strings.Fields(text)
	genericUpper := map[string]bool{
		"THE": true, "ALL": true, "ONLY": true, "EVERY": true, "EACH": true,
		"THIS": true, "THAT": true, "THESE": true, "THOSE": true, "WHAT": true,
		"WHEN": true, "WHICH": true, "WHERE": true, "WHILE": true, "WITH": true,
		"FROM": true, "INTO": true, "OVER": true, "UNDER": true, "BETWEEN": true,
		"THROUGH": true, "BEFORE": true, "AFTER": true, "ABOVE": true, "BELOW": true,
		"SYSTEM": true, "SYSTEMS": true, "COMPONENT": true, "COMPONENTS": true,
	}
	for _, w := range words {
		cleaned := strings.Trim(w, ".,;:!?\"'()[]{}-")
		if len(cleaned) >= 2 {
			upper := strings.ToUpper(cleaned)
			if !genericUpper[upper] && cleaned[0] >= 'A' && cleaned[0] <= 'Z' {
				// Check if the rest is lowercase (camelCase/PascalCase name)
				hasLower := false
				for i := 1; i < len(cleaned); i++ {
					if cleaned[i] >= 'a' && cleaned[i] <= 'z' {
						hasLower = true
						break
					}
				}
				if hasLower {
					return cleaned
				}
			}
		}
	}

	// Fall back to first meaningful keyword
	for _, kw := range keywords {
		lower := strings.ToLower(kw)
		genericLower := map[string]bool{
			"encrypt": true, "encryption": true, "access": true, "authentication": true,
			"authorization": true, "network": true, "security": true, "communication": true,
			"configuration": true, "connection": true, "data": true, "information": true,
		}
		if !genericLower[lower] && len(kw) > 2 {
			return kw
		}
	}
	if len(keywords) > 0 {
		return keywords[0]
	}
	return "general"
}

// ──────────────────────────────────────────────
// Intelligence Engine Conversion Helpers
// ──────────────────────────────────────────────

func convertToIntelArch(arch *ArchDescription) *intelligence.ArchDescription {
	if arch == nil {
		return nil
	}
	var components []intelligence.Component
	for _, c := range arch.Components {
		components = append(components, intelligence.Component{ID: c.ID, Label: c.Label})
	}
	var relationships []intelligence.Relation
	for _, r := range arch.Relationships {
		relationships = append(relationships, intelligence.Relation{Source: r.Source, Target: r.Target, Label: r.Label})
	}
	return &intelligence.ArchDescription{
		Name:                arch.Name,
		Components:          components,
		Relationships:       relationships,
		Policies:            arch.Policies,
		RawText:             arch.RawText,
		ExplicitAssumptions: arch.ExplicitAssumptions,
		SecurityControls:    arch.SecurityControls,
		Compliance:          arch.Compliance,
		ExpectedResults:     arch.ExpectedResults,
		ValidationCriteria:  arch.ValidationCriteria,
		Notes:               arch.Notes,
	}
}

func convertAssumptionsToIntel(assumptions []Assumption) []intelligence.Assumption {
	var result []intelligence.Assumption
	for _, a := range assumptions {
		var stride []intelligence.StrideCategory
		for _, s := range a.Stride {
			stride = append(stride, intelligence.StrideCategory(s))
		}
		result = append(result, intelligence.Assumption{
			ID:          a.ID,
			Description: a.Description,
			Component:   a.Component,
			Category:    a.Category,
			Risk:        intelligence.RiskLevel(a.Risk),
			Stride:      stride,
			Likelihood:  a.Likelihood,
			Impact:      a.Impact,
			Confidence:  a.Confidence,
			Keywords:    a.Keywords,
			SourceType:  a.SourceType,
			SourceFile:  a.SourceFile,
		})
	}
	return result
}

func convertIntelAssumptions(assumptions []intelligence.Assumption) []Assumption {
	var result []Assumption
	for _, a := range assumptions {
		var stride []StrideCategory
		for _, s := range a.Stride {
			stride = append(stride, StrideCategory(s))
		}
		result = append(result, Assumption{
			ID:           a.ID,
			Description:  a.Description,
			Component:    a.Component,
			Category:     a.Category,
			Risk:         RiskLevel(a.Risk),
			Stride:       stride,
			Likelihood:   a.Likelihood,
			Impact:       a.Impact,
			Confidence:   a.Confidence,
			Keywords:     a.Keywords,
			SourceType:   "intelligence",
			SourceFile:   a.SourceFile,
			QualityScore: a.QualityScore,
			Rationale:    a.Rationale,
		})
	}
	return result
}

func convertIntelContradictions(contradictions []intelligence.Contradiction) []Contradiction {
	var result []Contradiction
	for _, c := range contradictions {
		result = append(result, Contradiction{
			ID:                  c.ID,
			Severity:            RiskLevel(c.Severity),
			Description:         c.Description,
			Explanation:         c.Explanation,
			AffectedAssumptions: c.AffectedAssumptions,
			Evidence:            c.Evidence,
			RuleName:            c.RuleName,
		})
	}
	return result
}

func convertIntelTrustBoundaries(boundaries []intelligence.TrustBoundary) []TrustBoundary {
	var result []TrustBoundary
	for _, b := range boundaries {
		result = append(result, TrustBoundary{
			Type:        b.Type,
			Components:  b.Components,
			RiskLevel:   RiskLevel(b.RiskLevel),
			Description: b.Description,
		})
	}
	return result
}

func mergeAssumptions(existing, inferred []Assumption) []Assumption {
	seen := make(map[string]bool)
	var result []Assumption
	for _, a := range existing {
		key := normalizeText(a.Description)
		if !seen[key] {
			seen[key] = true
			result = append(result, a)
		}
	}
	for _, a := range inferred {
		key := normalizeText(a.Description)
		if !seen[key] {
			seen[key] = true
			result = append(result, a)
		}
	}
	return result
}

// CIE conversion helpers
func convertControlsToIntel(controls []ControlDetail) []intelligence.ControlDetail {
	var result []intelligence.ControlDetail
	for _, c := range controls {
		var stride []intelligence.StrideCategory
		for _, s := range c.MitigatedSTRIDE {
			stride = append(stride, intelligence.StrideCategory(s))
		}
		result = append(result, intelligence.ControlDetail{
			ID:                   c.ID,
			Name:                 c.ID,
			Description:          c.Description,
			Category:             c.Category,
			STRIDECovered:        stride,
			MitigatedAssumptions: c.MitigatedAssumptionIDs,
			Component:            "",
			Priority:             fmt.Sprintf("%d", c.Priority),
		})
	}
	return result
}

func convertTrustBoundariesToIntel(boundaries []TrustBoundary) []intelligence.TrustBoundary {
	var result []intelligence.TrustBoundary
	for _, b := range boundaries {
		result = append(result, intelligence.TrustBoundary{
			Type:        b.Type,
			Components:  b.Components,
			RiskLevel:   intelligence.RiskLevel(b.RiskLevel),
			Description: b.Description,
		})
	}
	return result
}

func convertCIEContradictions(cieContradictions []intelligence.CIEContradiction) []CIEContradiction {
	var result []CIEContradiction
	for _, c := range cieContradictions {
		result = append(result, CIEContradiction{
			ID:                      c.ID,
			Type:                    string(c.Type),
			Severity:                RiskLevel(c.Severity),
			Confidence:              c.Confidence,
			Summary:                 c.Summary,
			Description:             c.Description,
			StatementA:              CIEStatement{ID: c.StatementA.ID, Source: c.StatementA.Source, OriginalText: c.StatementA.OriginalText, Category: c.StatementA.Category, Confidence: c.StatementA.Confidence},
			StatementB:              CIEStatement{ID: c.StatementB.ID, Source: c.StatementB.Source, OriginalText: c.StatementB.OriginalText, Category: c.StatementB.Category, Confidence: c.StatementB.Confidence},
			AffectedAssets:          c.AffectedAssets,
			AffectedComponents:      c.AffectedComponents,
			AffectedControls:        c.AffectedControls,
			AffectedTrustBoundaries: c.AffectedTrustBoundaries,
			Reasoning:               c.Reasoning,
			Evidence:                c.Evidence,
			Recommendations:         c.Recommendations,
		})
	}
	return result
}

func deduplicateContradictions(contradictions []Contradiction) []Contradiction {
	seen := make(map[string]bool)
	var result []Contradiction
	for _, c := range contradictions {
		assumptions := make([]string, len(c.AffectedAssumptions))
		copy(assumptions, c.AffectedAssumptions)
		sort.Strings(assumptions)
		key := c.RuleName + "|" + strings.Join(assumptions, ",")
		if !seen[key] {
			seen[key] = true
			result = append(result, c)
		}
	}
	return result
}

func deduplicateCIEContradictions(contradictions []CIEContradiction) []CIEContradiction {
	// Two-phase dedup:
	// Phase 1: Dedup by exact statement text pair (handles same claim with different IDs)
	// Phase 2: Dedup by type + summary (handles semantically same claims with slightly different text)
	textSeen := make(map[string]bool)
	var phase1 []CIEContradiction
	for _, c := range contradictions {
		aText := strings.ToLower(strings.TrimSpace(c.StatementA.OriginalText))
		bText := strings.ToLower(strings.TrimSpace(c.StatementB.OriginalText))

		if aText == bText {
			continue
		}
		if aText > bText {
			aText, bText = bText, aText
		}
		key := c.Type + "|" + aText + "|" + bText
		if !textSeen[key] {
			textSeen[key] = true
			phase1 = append(phase1, c)
		}
	}

	// Phase 2: Dedup by type + summary (semantic dedup)
	summarySeen := make(map[string]bool)
	var result []CIEContradiction
	for _, c := range phase1 {
		summary := strings.ToLower(strings.TrimSpace(c.Summary))
		key := c.Type + "|" + summary
		if !summarySeen[key] {
			summarySeen[key] = true
			result = append(result, c)
		}
	}
	return result
}

// TBI conversion helpers
func convertTBIZones(zones []intelligence.TrustZone) []TBITrustZone {
	var result []TBITrustZone
	for _, z := range zones {
		result = append(result, TBITrustZone{
			ID:          z.ID,
			Name:        z.Name,
			Type:        string(z.Type),
			Sensitivity: z.Sensitivity,
			Components:  z.Components,
			Description: z.Description,
		})
	}
	return result
}

func convertTBIBoundaries(boundaries []intelligence.TBITrustBoundary) []TBITrustBoundary {
	var result []TBITrustBoundary
	for _, b := range boundaries {
		result = append(result, TBITrustBoundary{
			ID:                  b.ID,
			SourceZone:          b.SourceZone,
			DestinationZone:     b.DestinationZone,
			SourceZoneType:      string(b.SourceZoneType),
			DestinationZoneType: string(b.DestinationZoneType),
			CrossingType:        string(b.CrossingType),
			Risk:                RiskLevel(b.Risk),
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
	return result
}

func convertTBIWeaknesses(weaknesses []intelligence.BoundaryWeakness) []TBIWeakness {
	var result []TBIWeakness
	for _, w := range weaknesses {
		result = append(result, TBIWeakness{
			ID:              w.ID,
			BoundaryID:      w.BoundaryID,
			Type:            w.Type,
			Severity:        RiskLevel(w.Severity),
			Description:     w.Description,
			Reasoning:       w.Reasoning,
			Recommendations: w.Recommendations,
		})
	}
	return result
}

// Convert TBI boundaries back to intelligence format for TMI engine.
func convertTBIBoundariesToIntel(boundaries []TBITrustBoundary) []intelligence.TBITrustBoundary {
	var result []intelligence.TBITrustBoundary
	for _, b := range boundaries {
		result = append(result, intelligence.TBITrustBoundary{
			ID:                  b.ID,
			SourceZone:          b.SourceZone,
			DestinationZone:     b.DestinationZone,
			SourceZoneType:      intelligence.TrustZoneType(b.SourceZoneType),
			DestinationZoneType: intelligence.TrustZoneType(b.DestinationZoneType),
			CrossingType:        intelligence.CrossingType(b.CrossingType),
			Risk:                intelligence.RiskLevel(b.Risk),
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
	return result
}

// TMI conversion helpers
func convertIntelThreats(threats []intelligence.Threat) []Threat {
	var result []Threat
	for _, t := range threats {
		result = append(result, Threat{
			ID:                 t.ID,
			Name:               t.Name,
			Category:           string(t.Category),
			Severity:           RiskLevel(t.Severity),
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
	return result
}

func convertIntelThreatClusters(clusters []intelligence.ThreatCluster) []ThreatCluster {
	var result []ThreatCluster
	for _, c := range clusters {
		result = append(result, ThreatCluster{
			ID:              c.ID,
			Name:            c.Name,
			Category:        c.Category,
			Threats:         c.Threats,
			RiskScore:       c.RiskScore,
			AffectedAssets:  c.AffectedAssets,
			Recommendations: c.Recommendations,
		})
	}
	return result
}

func convertIntelThreatModelSummary(summary intelligence.ThreatModelSummary) ThreatModelSummary {
	return ThreatModelSummary{
		TotalThreats:       summary.TotalThreats,
		CriticalCount:      summary.CriticalCount,
		HighCount:          summary.HighCount,
		MediumCount:        summary.MediumCount,
		LowCount:           summary.LowCount,
		ClusterCount:       summary.ClusterCount,
		STRIDEDistribution: summary.STRIDEDistribution,
		TopThreats:         summary.TopThreats,
		SummaryText:        summary.SummaryText,
	}
}

// APD conversion helpers
func convertTBIZonesToIntel(zones []TBITrustZone) []intelligence.TrustZone {
	var result []intelligence.TrustZone
	for _, z := range zones {
		result = append(result, intelligence.TrustZone{
			ID:          z.ID,
			Name:        z.Name,
			Type:        intelligence.TrustZoneType(z.Type),
			Sensitivity: z.Sensitivity,
			Components:  z.Components,
			Description: z.Description,
		})
	}
	return result
}

func convertAPDAttackPaths(paths []intelligence.AttackPath) []AttackPath {
	var result []AttackPath
	for _, p := range paths {
		steps := make([]AttackStep, len(p.AttackSteps))
		for i, s := range p.AttackSteps {
			steps[i] = AttackStep{
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
		result = append(result, AttackPath{
			ID:                  p.ID,
			Name:                p.Name,
			Description:         p.Description,
			EntryPoint:          p.EntryPoint,
			TargetAsset:         p.TargetAsset,
			AttackSteps:         steps,
			RequiredAssumptions: p.RequiredAssumptions,
			RequiredConditions:  p.RequiredConditions,
			ExploitedThreats:    p.ExploitedThreats,
			AffectedComponents:  p.AffectedComponents,
			AffectedBoundaries:  p.AffectedBoundaries,
			Likelihood:          p.Likelihood,
			Impact:              p.Impact,
			RiskScore:           p.RiskScore,
			Confidence:          p.Confidence,
			DetectionDifficulty: p.DetectionDifficulty,
			BusinessImpact:      p.BusinessImpact,
			Recommendations:     p.Recommendations,
			KillChainPhases:     p.KillChainPhases,
			MITREATTACK:         p.MITREATTACK,
			STRIDECategories:    p.STRIDECategories,
		})
	}
	return result
}

func convertAPDThreatChains(chains []intelligence.ThreatChain) []ThreatChain {
	var result []ThreatChain
	for _, c := range chains {
		result = append(result, ThreatChain{
			ID:        c.ID,
			Threats:   c.Threats,
			Path:      c.Path,
			RiskScore: c.RiskScore,
			Reasoning: c.Reasoning,
		})
	}
	return result
}

func convertAPDSummary(result *intelligence.APDRunResult) AttackPathSummary {
	if result == nil {
		return AttackPathSummary{}
	}
	critical := 0
	high := 0
	medium := 0
	low := 0
	for _, p := range result.AttackPaths {
		switch {
		case p.RiskScore >= 0.6:
			critical++
		case p.RiskScore >= 0.4:
			high++
		case p.RiskScore >= 0.2:
			medium++
		default:
			low++
		}
	}
	topPaths := make([]string, 0, len(result.TopPaths))
	for _, p := range result.TopPaths {
		topPaths = append(topPaths, p.Name)
	}
	mitreCoverage := make([]string, 0)
	for technique := range result.MITREMapping {
		mitreCoverage = append(mitreCoverage, technique)
	}
	sort.Strings(mitreCoverage)
	return AttackPathSummary{
		TotalAttackPaths:  len(result.AttackPaths),
		CriticalCount:     critical,
		HighCount:         high,
		MediumCount:       medium,
		LowCount:          low,
		ThreatChainCount:  len(result.ThreatChains),
		TopAttackPaths:    topPaths,
		KillChainCoverage: result.KillChainCoverage,
		MITRECoverage:     mitreCoverage,
		SummaryText:       result.Summary,
	}
}

// SDRI conversion helpers
func convertSDRIControls(controls []intelligence.SDRIControl) []SDRIControl {
	var result []SDRIControl
	for _, c := range controls {
		result = append(result, SDRIControl{
			ID:          c.ID,
			Name:        c.Name,
			Category:    c.Category,
			Description: c.Description,
			ControlType: string(c.ControlType),
			Preventive:  c.Preventive,
			Detective:   c.Detective,
			Corrective:  c.Corrective,
			Strength:    string(c.Strength),
			Evidence:    c.Evidence,
			Coverage:    c.Coverage,
			Status:      c.Status,
		})
	}
	return result
}

func convertSDRIDesignFindings(findings []intelligence.SDRIFinding) []SDRIDesignFinding {
	var result []SDRIDesignFinding
	for _, f := range findings {
		result = append(result, SDRIDesignFinding{
			ID:                 f.ID,
			Title:              f.Title,
			Description:        f.Description,
			Severity:           f.Severity,
			Category:           f.Category,
			AffectedComponents: f.AffectedComponents,
			AffectedControls:   f.AffectedControls,
			BusinessImpact:     f.BusinessImpact,
			Recommendation:     f.Recommendation,
			Reasoning:          f.Reasoning,
		})
	}
	return result
}

func convertSDRIWeaknesses(weaknesses []intelligence.SDRIArchitecturalWeakness) []SDRIArchitecturalWeakness {
	var result []SDRIArchitecturalWeakness
	for _, w := range weaknesses {
		result = append(result, SDRIArchitecturalWeakness{
			ID:             w.ID,
			Pattern:        w.Pattern,
			Description:    w.Description,
			Severity:       w.Severity,
			Components:     w.Components,
			Impact:         w.Impact,
			Recommendation: w.Recommendation,
		})
	}
	return result
}

func convertSDRIRemediations(remediations []intelligence.SDRIRemediation) []SDRIRemediation {
	var result []SDRIRemediation
	for _, r := range remediations {
		result = append(result, SDRIRemediation{
			ID:                 r.ID,
			Priority:           r.Priority,
			Description:        r.Description,
			RiskScore:          r.RiskScore,
			BusinessImpact:     r.BusinessImpact,
			Effort:             r.Effort,
			Category:           r.Category,
			Recommendation:     r.Recommendation,
			AffectedComponents: r.AffectedComponents,
		})
	}
	return result
}

func convertSDRICoverage(coverage []intelligence.SDRICoverage) []SDRICoverageItem {
	var result []SDRICoverageItem
	for _, c := range coverage {
		result = append(result, SDRICoverageItem{
			Category: c.Category,
			Expected: c.Expected,
			Observed: c.Observed,
			Coverage: c.Coverage,
			Level:    c.Level,
		})
	}
	return result
}

func convertSDRIComplianceMappings(mappings []intelligence.SDRIComplianceMapping) []SDRIComplianceMapping {
	var result []SDRIComplianceMapping
	for _, m := range mappings {
		result = append(result, SDRIComplianceMapping{
			Framework: m.Framework,
			Coverage:  m.Coverage,
			Controls:  m.Controls,
			Status:    m.Status,
		})
	}
	return result
}

func averageCoverage(coverage []intelligence.SDRICoverage) float64 {
	if len(coverage) == 0 {
		return 0
	}
	total := 0.0
	for _, c := range coverage {
		total += c.Coverage
	}
	return total / float64(len(coverage))
}

// ── CIARE Conversion Helpers ──

func convertCIAREFrameworkCoverages(c []intelligence.CIAREFrameworkCoverage) []CIAREFrameworkCoverage {
	var r []CIAREFrameworkCoverage
	for _, v := range c {
		r = append(r, CIAREFrameworkCoverage{
			Framework:        v.Framework,
			Required:         v.Required,
			Observed:         v.Observed,
			Missing:          v.Missing,
			CoveragePct:      v.CoveragePct,
			Status:           v.Status,
			ObservedControls: v.ObservedControls,
			MissingControls:  v.MissingControls,
		})
	}
	return r
}

func convertCIAREAuditReadiness(a []intelligence.CIAREAuditReadiness) []CIAREAuditReadiness {
	var r []CIAREAuditReadiness
	for _, v := range a {
		r = append(r, CIAREAuditReadiness{
			Framework:       v.Framework,
			ReadinessScore:  v.ReadinessScore,
			Status:          v.Status,
			ControlCoverage: v.ControlCoverage,
			EvidenceScore:   v.EvidenceScore,
			ThreatExposure:  v.ThreatExposure,
			FindingsPenalty: v.FindingsPenalty,
			Factors:         v.Factors,
		})
	}
	return r
}

func convertCIAREEvidenceRequirements(e []intelligence.CIAREEvidenceRequirement) []CIAREEvidenceRequirement {
	var r []CIAREEvidenceRequirement
	for _, v := range e {
		r = append(r, CIAREEvidenceRequirement{
			Framework: v.Framework,
			Control:   v.Control,
			Evidence:  v.Evidence,
		})
	}
	return r
}

func convertCIAREMissingEvidences(m []intelligence.CIAREMissingEvidence) []CIAREMissingEvidence {
	var r []CIAREMissingEvidence
	for _, v := range m {
		r = append(r, CIAREMissingEvidence{
			Framework: v.Framework,
			Control:   v.Control,
			Evidences: v.Evidences,
		})
	}
	return r
}

func convertCIAREAuditorQuestions(q []intelligence.CIAREAuditorQuestion) []CIAREAuditorQuestion {
	var r []CIAREAuditorQuestion
	for _, v := range q {
		r = append(r, CIAREAuditorQuestion{
			Framework: v.Framework,
			Control:   v.Control,
			Question:  v.Question,
		})
	}
	return r
}

func convertCIAREComplianceGaps(g []intelligence.CIAREComplianceGap) []CIAREComplianceGap {
	var r []CIAREComplianceGap
	for _, v := range g {
		r = append(r, CIAREComplianceGap{
			ID:          v.ID,
			Framework:   v.Framework,
			Requirement: v.Requirement,
			Observed:    v.Observed,
			Missing:     v.Missing,
			Risk:        v.Risk,
		})
	}
	return r
}

func convertCIAREControlMaturities(m []intelligence.CIAREControlMaturity) []CIAREControlMaturity {
	var r []CIAREControlMaturity
	for _, v := range m {
		r = append(r, CIAREControlMaturity{
			Domain:   v.Domain,
			Level:    v.Level,
			Label:    v.Label,
			Coverage: v.Coverage,
		})
	}
	return r
}

func convertCIAREComplianceNarratives(n []intelligence.CIAREComplianceNarrative) []CIAREComplianceNarrative {
	var r []CIAREComplianceNarrative
	for _, v := range n {
		r = append(r, CIAREComplianceNarrative{
			Framework: v.Framework,
			Narrative: v.Narrative,
		})
	}
	return r
}

func convertCIAREAuditPackage(p intelligence.CIAREAuditPackage) *CIAREAuditPackage {
	return &CIAREAuditPackage{
		ExecutiveSummary:     p.ExecutiveSummary,
		FrameworkCoverages:   convertCIAREFrameworkCoverages(p.FrameworkCoverages),
		ControlInventory:     convertSDRIControls(p.ControlInventory),
		MissingControls:      convertCIAREComplianceGaps(p.MissingControls),
		EvidenceRequirements: convertCIAREEvidenceRequirements(p.EvidenceRequirements),
		AuditorQuestions:     convertCIAREAuditorQuestions(p.AuditorQuestions),
	}
}

func convertCIAREComplianceDashboard(d intelligence.CIAREComplianceDashboard) *CIAREComplianceDashboard {
	return &CIAREComplianceDashboard{
		FrameworkCoverages: d.FrameworkCoverages,
		TopGaps:            convertCIAREComplianceGaps(d.TopGaps),
		TopMissingEvidence: convertCIAREMissingEvidences(d.TopMissingEvidence),
		TopRisks:           d.TopRisks,
	}
}

func convertCIAREProcurementQuestions(q []intelligence.CIAREProcurementQuestion) []CIAREProcurementQuestion {
	var r []CIAREProcurementQuestion
	for _, v := range q {
		r = append(r, CIAREProcurementQuestion{
			Category: v.Category,
			Question: v.Question,
		})
	}
	return r
}

func ciareAvgCoverage(coverages []intelligence.CIAREFrameworkCoverage) float64 {
	if len(coverages) == 0 {
		return 0
	}
	total := 0.0
	for _, c := range coverages {
		total += c.CoveragePct
	}
	return total / float64(len(coverages))
}

// ── DKPI Conversion Helpers ──

func convertDKPIResult(r *intelligence.DKPIEngineResult) DKPIIntelligence {
	if r == nil {
		return DKPIIntelligence{}
	}
	dkpi := DKPIIntelligence{
		DomainResult: DKPIDomainResult{
			PrimaryDomain: r.DetectedDomain.PrimaryDomain,
			Confidence:    r.DetectedDomain.Confidence,
			Rationale:     r.DetectedDomain.Rationale,
		},
		Recommendations:    r.Recommendations,
		InjectedThreats:    convertIntelThreats(r.InjectedThreats),
		DomainControls:     convertSDRIControls(r.DomainControls),
		DomainCompliance:   r.DomainCompliance,
		EvidenceReqs:       convertDKPIEvidenceReqs(r.DomainEvidence),
		BoostedAssumptions: convertIntelAssumptions(r.BoostedAssumptions),
	}
	for _, m := range r.DetectedDomain.Matches {
		dkpi.DomainResult.Matches = append(dkpi.DomainResult.Matches, DKPIDomainMatch{
			PackID:     m.PackID,
			PackName:   m.PackName,
			Score:      m.Score,
			Confidence: m.Confidence,
			Reasons:    m.Reasons,
		})
	}
	if r.ActivePack != nil {
		dkpi.ActivePack = convertDKPIPack(r.ActivePack)
	}
	if r.DetectedDomain.PrimaryDomain != "" {
		dkpi.Summary = fmt.Sprintf("Domain detected: %s (%.1f%% confidence). %d recommendations, %d domain threats.",
			r.DetectedDomain.PrimaryDomain, r.DetectedDomain.Confidence,
			len(r.Recommendations), len(r.InjectedThreats))
	} else {
		dkpi.Summary = "No domain detected — generic security analysis."
	}
	return dkpi
}

func convertDKPIPack(p *intelligence.KnowledgePack) *DKPIKnowledgePack {
	if p == nil {
		return nil
	}
	pack := &DKPIKnowledgePack{
		ID:                   p.ID,
		Name:                 p.Name,
		Industry:             p.Industry,
		Description:          p.Description,
		CrownJewels:          p.CrownJewels,
		ComplianceFrameworks: p.ComplianceFrameworks,
	}
	for _, c := range p.ExpectedControls {
		pack.ExpectedControls = append(pack.ExpectedControls, DKPIKnowledgePackControl{
			Name: c.Name, Description: c.Description,
			Category: c.Category, Priority: c.Priority,
		})
	}
	for _, t := range p.ThreatPatterns {
		pack.ThreatPatterns = append(pack.ThreatPatterns, DKPIKnowledgePackThreat{
			Name: t.Name, Description: t.Description,
			Severity: t.Severity, Category: t.Category,
		})
	}
	return pack
}

func convertDKPIEvidenceReqs(ev []intelligence.KnowledgePackEvidence) []DKPIKnowledgePackEvidence {
	var r []DKPIKnowledgePackEvidence
	for _, e := range ev {
		r = append(r, DKPIKnowledgePackEvidence{
			Control: e.Control, Evidence: e.Evidence,
		})
	}
	return r
}

func mergeSDRIControls(existing, enriched []SDRIControl) []SDRIControl {
	seen := make(map[string]int)
	for i, c := range existing {
		seen[c.ID] = i
	}
	for _, c := range enriched {
		if idx, ok := seen[c.ID]; ok {
			if c.Category != "" {
				existing[idx].Category = c.Category
			}
			if c.Coverage != "" {
				existing[idx].Coverage = c.Coverage
			}
		} else {
			existing = append(existing, c)
			seen[c.ID] = len(existing) - 1
		}
	}
	return existing
}

func convertSDIResult(r *intelligence.SDIResult) SDIIntelligence {
	if r == nil {
		return SDIIntelligence{}
	}
	sdi := SDIIntelligence{}
	for _, rec := range r.Recommendations {
		sdi.Recommendations = append(sdi.Recommendations, SDIDecisionRecommendation{
			ID: rec.ID, Title: rec.Title, Description: rec.Description,
			AffectedFindings: rec.AffectedFindings, AffectedThreats: rec.AffectedThreats,
			AffectedAttackPaths: rec.AffectedAttackPaths, AffectedControls: rec.AffectedControls,
			AffectedAssets: rec.AffectedAssets, RiskReduction: rec.RiskReduction,
			Effort: rec.Effort, Priority: rec.Priority,
			BusinessImpact: rec.BusinessImpact, ComplianceImpact: rec.ComplianceImpact,
			Rationale: rec.Rationale,
		})
	}
	for _, sim := range r.FixSimulations {
		sdi.FixSimulations = append(sdi.FixSimulations, SDIFixSimulation{
			ControlName: sim.ControlName, ControlCategory: sim.ControlCategory,
			OriginalCritical: sim.OriginalCritical, OriginalHigh: sim.OriginalHigh,
			OriginalTotal: sim.OriginalTotal, OriginalAttackPaths: sim.OriginalAttackPaths,
			OriginalCoverage: sim.OriginalCoverage,
			NewCritical:      sim.NewCritical, NewHigh: sim.NewHigh,
			NewTotal: sim.NewTotal, NewAttackPaths: sim.NewAttackPaths,
			NewCoverage: sim.NewCoverage,
		})
	}
	for _, sim := range r.FailureSimulations {
		sdi.FailureSimulations = append(sdi.FailureSimulations, SDIFailureSimulation{
			ControlName: sim.ControlName, ControlCategory: sim.ControlCategory,
			SystemsImpacted: sim.SystemsImpacted, AttackPathsOpened: sim.AttackPathsOpened,
			NewFindings: sim.NewFindings, RiskIncrease: sim.RiskIncrease,
			RiskScoreIncrease: sim.RiskScoreIncrease,
		})
	}
	for _, ci := range r.ControlImpacts {
		sdi.ControlImpacts = append(sdi.ControlImpacts, SDIControlImpact{
			ControlName: ci.ControlName, Category: ci.Category,
			SecurityValue: ci.SecurityValue, Effort: ci.Effort, ROI: ci.ROI,
			FindingCount: ci.FindingCount, ThreatCount: ci.ThreatCount,
			AttackPathCount: ci.AttackPathCount,
		})
	}
	sdi.DecisionTrees = SDIDecisionTreeResult{
		SingleAction: sdiConvertTree(r.DecisionTrees.SingleAction),
		ThreeActions: sdiConvertTree(r.DecisionTrees.ThreeActions),
		FiveActions:  sdiConvertTree(r.DecisionTrees.FiveActions),
	}
	sdi.BoardScenarios = SDIBoardScenarios{
		DoNothing:        sdiConvertBoardScenario(r.BoardScenarios.DoNothing),
		PartialRemediate: sdiConvertBoardScenario(r.BoardScenarios.PartialRemediate),
		FullRemediate:    sdiConvertBoardScenario(r.BoardScenarios.FullRemediate),
	}
	for _, ip := range r.InvestmentPriorities {
		sdi.InvestmentPriorities = append(sdi.InvestmentPriorities, SDIInvestmentPriority{
			Area: ip.Area, Rank: ip.Rank, Score: ip.Score,
			Rationale: ip.Rationale, FindingCount: ip.FindingCount,
			RiskReduction: ip.RiskReduction,
		})
	}
	for _, apc := range r.AttackPathCollapse {
		sdi.AttackPathCollapse = append(sdi.AttackPathCollapse, SDIAttackPathCollapse{
			ControlName: apc.ControlName, Category: apc.Category,
			AttackPathsReduced: apc.AttackPathsReduced,
			TotalAttackPaths:   apc.TotalAttackPaths,
			ReductionPercent:   apc.ReductionPercent,
		})
	}
	for _, ci := range r.ComplianceImpacts {
		sdi.ComplianceImpacts = append(sdi.ComplianceImpacts, SDIComplianceImpact{
			Framework: ci.Framework, Action: ci.Action,
			Improvement: ci.Improvement, Rationale: ci.Rationale,
		})
	}
	sdi.RemediationRoadmap = SDIRemediationRoadmap{
		Phase30:  sdiConvertRoadmapItems(r.RemediationRoadmap.Phase30),
		Phase90:  sdiConvertRoadmapItems(r.RemediationRoadmap.Phase90),
		Phase180: sdiConvertRoadmapItems(r.RemediationRoadmap.Phase180),
		Phase12m: sdiConvertRoadmapItems(r.RemediationRoadmap.Phase12m),
	}
	sdi.Dashboard = SDIDecisionDashboard{
		RiskReductionSummary: r.Dashboard.RiskReductionSummary,
		TotalRiskReduction:   r.Dashboard.TotalRiskReduction,
	}
	for _, td := range r.Dashboard.TopDecisions {
		sdi.Dashboard.TopDecisions = append(sdi.Dashboard.TopDecisions, sdiConvertRec(td))
	}
	for _, qw := range r.Dashboard.QuickWins {
		sdi.Dashboard.QuickWins = append(sdi.Dashboard.QuickWins, sdiConvertRec(qw))
	}
	for _, sa := range r.Dashboard.StrategicActions {
		sdi.Dashboard.StrategicActions = append(sdi.Dashboard.StrategicActions, sdiConvertRec(sa))
	}
	sdi.ExecutiveScenarios = SDIExecutiveScenarios{
		BestCase:   sdiConvertExecScenario(r.ExecutiveScenarios.BestCase),
		LikelyCase: sdiConvertExecScenario(r.ExecutiveScenarios.LikelyCase),
		WorstCase:  sdiConvertExecScenario(r.ExecutiveScenarios.WorstCase),
	}
	return sdi
}

func sdiConvertTree(dt intelligence.DecisionTree) SDIDecisionTree {
	out := SDIDecisionTree{Budget: dt.Budget, ActionCount: dt.ActionCount, Rationale: dt.Rationale}
	for _, r := range dt.RecommendedOrder {
		out.RecommendedOrder = append(out.RecommendedOrder, sdiConvertRec(r))
	}
	return out
}

func sdiConvertBoardScenario(bs intelligence.BoardScenario) SDIBoardScenario {
	return SDIBoardScenario{
		Scenario: bs.Scenario, Description: bs.Description,
		RiskScore: bs.RiskScore, CriticalFindings: bs.CriticalFindings,
		AttackPaths: bs.AttackPaths, CoverageRate: bs.CoverageRate,
		KeyRisks: bs.KeyRisks,
	}
}

func sdiConvertRoadmapItems(items []intelligence.SDIRoadmapItem) []SDIRoadmapItem {
	var out []SDIRoadmapItem
	for _, item := range items {
		out = append(out, SDIRoadmapItem{
			Action: item.Action, Category: item.Category,
			Priority: item.Priority, Effort: item.Effort,
			RiskReduction: item.RiskReduction,
		})
	}
	return out
}

func sdiConvertRec(r intelligence.DecisionRecommendation) SDIDecisionRecommendation {
	return SDIDecisionRecommendation{
		ID: r.ID, Title: r.Title, Description: r.Description,
		AffectedFindings: r.AffectedFindings, AffectedThreats: r.AffectedThreats,
		AffectedAttackPaths: r.AffectedAttackPaths, AffectedControls: r.AffectedControls,
		AffectedAssets: r.AffectedAssets, RiskReduction: r.RiskReduction,
		Effort: r.Effort, Priority: r.Priority,
		BusinessImpact: r.BusinessImpact, ComplianceImpact: r.ComplianceImpact,
		Rationale: r.Rationale,
	}
}

func sdiConvertExecScenario(es intelligence.ExecutiveScenario) SDIExecutiveScenario {
	return SDIExecutiveScenario{
		Scenario: es.Scenario, RiskScore: es.RiskScore,
		FindingsResolved: es.FindingsResolved, AttackPathsClosed: es.AttackPathsClosed,
		CoverageAchieved: es.CoverageAchieved, Description: es.Description,
	}
}

func convertSDTResult(r *intelligence.SDTResult) SDTIntelligence {
	if r == nil {
		return SDTIntelligence{}
	}
	sdt := SDTIntelligence{
		Twin: ArchitectureTwinPR{
			ID: r.Twin.ID, Version: r.Twin.Version,
			ArchitectureName: r.Twin.ArchitectureName, Domain: r.Twin.Domain,
			RiskScore: r.Twin.RiskScore, Coverage: r.Twin.Coverage,
			SourceHash: r.Twin.SourceHash,
		},
	}
	for _, ci := range r.ChangeImpacts {
		sdt.ChangeImpacts = append(sdt.ChangeImpacts, ChangeImpactPR{
			Change: ci.Change, ComponentAffected: ci.ComponentAffected,
			ImpactType: ci.ImpactType, Severity: ci.Severity,
			RisksAffected: ci.RisksAffected, AttackPathsAffected: ci.AttackPathsAffected,
			ControlsAffected: ci.ControlsAffected, Description: ci.Description,
		})
	}
	for _, ad := range r.ArchitectureDiffs {
		sdt.ArchitectureDiffs = append(sdt.ArchitectureDiffs, ArchitectureDiffPR{
			Category: ad.Category, AddedCount: ad.AddedCount, RemovedCount: ad.RemovedCount,
			ChangedCount: ad.ChangedCount, RiskScoreDelta: ad.RiskScoreDelta,
			CoverageDelta: ad.CoverageDelta, Description: ad.Description,
		})
	}
	for _, ei := range r.EvolutionInsights {
		sdt.EvolutionInsights = append(sdt.EvolutionInsights, EvolutionInsightPR{
			Scenario: ei.Scenario, Assumption: ei.Assumption,
			Status: ei.Status, Rationale: ei.Rationale,
		})
	}
	for _, cd := range r.ControlDrifts {
		sdt.ControlDrifts = append(sdt.ControlDrifts, ControlDriftPR{
			ControlName: cd.ControlName, Category: cd.Category,
			ExpectedState: cd.ExpectedState, CurrentState: cd.CurrentState,
			RiskImpact: cd.DriftType,
		})
	}
	for _, ad := range r.AssumptionDecays {
		sdt.AssumptionDecays = append(sdt.AssumptionDecays, AssumptionDecayPR{
			AssumptionID: ad.AssumptionID, Description: ad.Description,
			TimeElapsed: ad.Age, Status: ad.Status,
			Recommendation: fmt.Sprintf("evidence: %d sources", ad.EvidenceCount),
		})
	}
	sdt.SecurityDebt = SecurityDebtScorePR{
		TotalDebt: r.SecurityDebt.TotalDebt, FindingDebt: r.SecurityDebt.FindingDebt,
		ControlDebt: r.SecurityDebt.ControlDebt, AssumptionDebt: r.SecurityDebt.AssumptionDebt,
		RiskScore: r.SecurityDebt.RiskScore,
	}
	for _, cd := range r.ComplianceDrifts {
		sdt.ComplianceDrifts = append(sdt.ComplianceDrifts, ComplianceDriftPR{
			Framework: cd.Framework, Status: cd.Status,
			NewGaps: cd.NewGaps, ResolvedGaps: cd.ResolvedGaps,
			RegressedAreas: cd.RegressedAreas,
		})
	}
	sdt.AttackSurfaceTrend = AttackSurfaceTrendPR{
		InternetExposure: r.AttackSurfaceTrend.InternetExposure,
		ThirdParties:     r.AttackSurfaceTrend.ThirdParties,
		IdentitySystems:  r.AttackSurfaceTrend.IdentitySystems,
		CloudServices:    r.AttackSurfaceTrend.CloudServices,
		AdminPaths:       r.AttackSurfaceTrend.AdminPaths,
		GrowthRate:       r.AttackSurfaceTrend.GrowthRate,
	}
	sdt.Timeline = ArchitectureTimelinePR{
		Trend: r.Timeline.Trend, DeltaRisk: r.Timeline.DeltaRisk,
	}
	for _, wi := range r.WhatIfScenarios {
		sdt.WhatIfScenarios = append(sdt.WhatIfScenarios, WhatIfScenarioPR{
			Name: wi.Name, RiskDelta: wi.RiskDelta,
			CoverageDelta: wi.CoverageDelta, FindingsDelta: wi.ThreatDelta,
			Description: wi.Description,
		})
	}
	sdt.MergerAnalysis = MergerAnalysisPR{
		CombinedRiskScore: r.MergerAnalysis.CombinedRiskScore,
		InheritedRisks:    r.MergerAnalysis.InheritedRisks,
		InheritedControls: r.MergerAnalysis.InheritedControls,
		SharedRisks:       r.MergerAnalysis.SharedRisks,
	}
	ztd := ZeroTrustAnalysisPR{Overall: r.ZeroTrust.Overall, Target: r.ZeroTrust.Target, Gap: r.ZeroTrust.Gap}
	for _, d := range r.ZeroTrust.Dimensions {
		ztd.Dimensions = append(ztd.Dimensions, ZeroTrustDimensionPR{
			Dimension: d.Dimension, Score: d.CurrentScore,
			Target: d.TargetScore, Gap: d.Gap, Status: d.Progress,
		})
	}
	sdt.ZeroTrust = ztd
	for _, rs := range r.Resilience {
		sdt.Resilience = append(sdt.Resilience, ResilienceScenarioPR{
			FailurePoint: rs.FailurePoint, BusinessImpact: rs.BusinessImpact,
			SecurityImpact: rs.SecurityImpact, AffectedAssets: rs.AffectedAssets,
			AttackPathsOpened:   rs.AttackPathsOpened,
			RecoveryAssumptions: rs.RecoveryAssumptions,
		})
	}
	for _, cj := range r.CrownJewels {
		sdt.CrownJewels = append(sdt.CrownJewels, CrownJewelRankingPR{
			AssetName: cj.AssetName, BusinessValue: cj.BusinessValue,
			AttackValue: cj.AttackValue, DependencyCount: cj.DependencyCount,
			ThreatCount: cj.ThreatCount, BlastRadius: cj.BlastRadius,
			OverallScore: cj.OverallScore,
		})
	}
	sdt.ExecutiveReport = DigitalTwinReportPR{
		ArchitectureHealth:   r.ExecutiveReport.ArchitectureHealth,
		SecurityDebtScore:    r.ExecutiveReport.SecurityDebtScore,
		ControlDriftCount:    r.ExecutiveReport.ControlDriftCount,
		ComplianceDriftCount: r.ExecutiveReport.ComplianceDriftCount,
		RiskTrend:            r.ExecutiveReport.RiskTrend,
		AttackSurfaceTrend:   r.ExecutiveReport.AttackSurfaceTrend,
	}
	sdt.PortfolioSummary = PortfolioTwinSummaryPR{
		ArchitectureCount: r.PortfolioSummary.ArchitectureCount,
		SharedRisks:       r.PortfolioSummary.SharedRisks,
		SharedVendors:     r.PortfolioSummary.SharedVendors,
		SharedControls:    r.PortfolioSummary.SharedControls,
		EnterpriseTrends:  r.PortfolioSummary.EnterpriseTrends,
		AggregatedDebt:    r.PortfolioSummary.AggregatedDebt,
	}
	return sdt
}
