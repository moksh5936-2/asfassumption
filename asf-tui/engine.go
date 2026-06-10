package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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
	ID           string
	Description  string
	Component    string
	Category     string
	Risk         RiskLevel
	Stride       []StrideCategory
	Likelihood   int
	Impact       int
	Confidence   float64
	Keywords     []string

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
}

type AnalysisResult struct {
	ArchitectureName   string
	AnalysisDate       time.Time
	AnalysisMode       string
	Assumptions        []Assumption
	CriticalCount      int
	HighCount          int
	TotalAssumptions   int
	StrideDistribution map[StrideCategory]int
	Controls           []ControlDetail
	Compliance         []string
	Summary            string
	TrueAssumptions    int
	FalseAssumptions   int
	CriticalGaps       int

	// Explainability fields
	EvidenceSummary      EvidenceSummary `json:"evidence_summary"`
	RiskModelVersion     string          `json:"risk_model_version"`
	ConfidenceSummary    string          `json:"confidence_summary"`
}

type AnalysisProgress struct {
	Percent  float64
	Stage    string
	Complete bool
}

type Engine struct {
	config        *Config
	pythonPath    string
	projectDir    string
	strideEngine  *StrideEngine
	explainPipe   *ExplainabilityPipeline
	archDesc      *ArchDescription
}

func NewEngine(cfg *Config) *Engine {
	return &Engine{
		config:       cfg,
		pythonPath:   "/Users/moksh/Project/cybersec/.venv/bin/python",
		projectDir:   "/Users/moksh/Project/cybersec",
		strideEngine: NewStrideEngine(),
	}
}

func (e *Engine) RunAnalysis(archPath, evPath, mode string, progress chan<- AnalysisProgress) (*AnalysisResult, error) {
	progress <- AnalysisProgress{Percent: 5, Stage: "Parsing Architecture..."}

	inputPath := archPath
	ext := strings.ToLower(filepath.Ext(archPath))
	needsTemp := ext == ".drawio" || ext == ".mmd" || ext == ".md" || ext == ".yaml" || ext == ".yml" || ext == ".json" || ext == ".svg" || ext == ".png" || ext == ".jpg" || ext == ".jpeg"

	if needsTemp {
		desc, err := ParseArchitecture(archPath)
		if err != nil {
			return nil, fmt.Errorf("parse architecture: %w", err)
		}
		e.archDesc = desc
		tmpFile, err := os.CreateTemp("", "asf-*.txt")
		if err != nil {
			return nil, fmt.Errorf("create temp file: %w", err)
		}
		if _, err := tmpFile.WriteString(desc.RawText); err != nil {
			tmpFile.Close()
			os.Remove(tmpFile.Name())
			return nil, fmt.Errorf("write temp file: %w", err)
		}
		tmpFile.Close()
		inputPath = tmpFile.Name()
		defer os.Remove(inputPath)
	} else {
		// Try to parse for evidence even for text files
		desc, err := ParseArchitecture(archPath)
		if err == nil {
			e.archDesc = desc
		}
	}

	progress <- AnalysisProgress{Percent: 20, Stage: "Running ASF Engine..."}

	pythonResult, err := e.callPythonCLI(inputPath, evPath)
	if err != nil {
		return nil, fmt.Errorf("ASF engine error: %w", err)
	}

	progress <- AnalysisProgress{Percent: 60, Stage: "Processing Results..."}

	result := e.buildResult(pythonResult, archPath, mode)

	progress <- AnalysisProgress{Percent: 80, Stage: "Generating STRIDE Mapping..."}
	result.StrideDistribution = e.mapStrideDistribution(result.Assumptions)

	if mode == ModeASFAndAI && e.config != nil && e.config.AI.Enabled && e.config.AI.ActiveModel != "" {
		progress <- AnalysisProgress{Percent: 85, Stage: "Running AI Enhancement..."}
		enhancer := NewAIEnhancer()
		aiResult, err := enhancer.Enhance(result, e.config.AI.ActiveModel)
		if err == nil && aiResult != nil {
			result = mergeAIResults(result, aiResult)
			result.StrideDistribution = e.mapStrideDistribution(result.Assumptions)
		}
	}

	progress <- AnalysisProgress{Percent: 100, Stage: "Complete", Complete: true}

	return result, nil
}

func (e *Engine) callPythonCLI(docPath, evPath string) (*asfJSONResult, error) {
	args := []string{"-m", "asf.cli.main", "analyze", "--json", docPath}
	if evPath != "" {
		if _, err := os.Stat(evPath); err == nil {
			args = append(args, "-e", evPath)
		}
	}
	cmd := exec.Command(e.pythonPath, args...)
	cmd.Dir = e.projectDir

	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		stderrStr := stderr.String()
		if stderrStr != "" {
			return nil, fmt.Errorf("%s (stderr: %s)", err, stderrStr)
		}
		return nil, err
	}

	var result asfJSONResult
	if err := json.Unmarshal([]byte(stdout.String()), &result); err != nil {
		return nil, fmt.Errorf("parse error: %w\nRaw: %s", err, stdout.String()[:min(len(stdout.String()), 200)])
	}

	return &result, nil
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
		ID                 string  `json:"id"`
		Text               string  `json:"text"`
		AssumptionType     string  `json:"assumption_type"`
		VerificationStatus string  `json:"verification_status"`
		Confidence         float64 `json:"confidence"`
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
		ArchitectureName: fileBase(archPath),
		AnalysisDate:     time.Now(),
		AnalysisMode:     mode,
		TotalAssumptions: r.Summary.Assumptions,
		TrueAssumptions:  r.Summary.Verified,
		FalseAssumptions: r.Summary.Contradicted,
		CriticalGaps:     r.Summary.CriticalGaps,
		StrideDistribution: make(map[StrideCategory]int),
		RiskModelVersion:  "asf-risk-model-1.0",
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

	for i, a := range r.Assumptions {
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

		result.Assumptions = append(result.Assumptions, assumption)

		if assumption.Risk == RiskCritical {
			result.CriticalCount++
		} else if assumption.Risk == RiskHigh {
			result.HighCount++
		}

		_ = i
	}

	// Build evidence summary
	if e.explainPipe != nil {
		result.EvidenceSummary = e.explainPipe.BuildEvidenceSummary(result.Assumptions)
		confSummary := buildConfidenceSummary(result.Assumptions)
		result.ConfidenceSummary = confSummary
	}

	result.Controls = generateControls(result.Assumptions)
	result.Compliance = []string{
		"ASF analysis completed — see gap analysis for compliance mapping",
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
	Category   string
	BaseDesc   string
	Rationale  string
	STRIDE     []StrideCategory
	Priority   int
}

func controlTemplates() []controlTemplate {
	return []controlTemplate{
		{Category: "IDENTITY", BaseDesc: "Implement strong identity verification with MFA",
			Rationale: "Identity-related assumptions require robust authentication to prevent spoofing and unauthorized access.",
			STRIDE: []StrideCategory{StrideSpoofing, StrideElevationPriv}, Priority: 1},
		{Category: "AUTHENTICATION", BaseDesc: "Enforce multi-factor authentication for all access",
			Rationale: "Authentication assumptions require verified identity to prevent credential-based attacks.",
			STRIDE: []StrideCategory{StrideSpoofing, StrideElevationPriv}, Priority: 1},
		{Category: "AUTHORIZATION", BaseDesc: "Implement role-based access control with principle of least privilege",
			Rationale: "Authorization assumptions require strict access boundaries to prevent privilege escalation.",
			STRIDE: []StrideCategory{StrideElevationPriv, StrideInfoDisclosure}, Priority: 1},
		{Category: "ACCESS", BaseDesc: "Enforce least-privilege access controls across all components",
			Rationale: "Access control assumptions limit blast radius and prevent lateral movement.",
			STRIDE: []StrideCategory{StrideElevationPriv, StrideInfoDisclosure}, Priority: 1},
		{Category: "NETWORK", BaseDesc: "Implement network segmentation and encryption in transit",
			Rationale: "Network assumptions require boundary protection to prevent data exposure and DoS.",
			STRIDE: []StrideCategory{StrideInfoDisclosure, StrideDenialOfService, StrideTampering}, Priority: 1},
		{Category: "ENCRYPTION", BaseDesc: "Implement encryption at rest and in transit for all sensitive data",
			Rationale: "Encryption assumptions protect confidentiality against data disclosure attacks.",
			STRIDE: []StrideCategory{StrideInfoDisclosure}, Priority: 1},
		{Category: "CONFIGURATION", BaseDesc: "Use infrastructure-as-code with automated configuration validation",
			Rationale: "Configuration assumptions prevent tampering through misconfiguration and drift.",
			STRIDE: []StrideCategory{StrideTampering}, Priority: 2},
		{Category: "DEPENDENCY", BaseDesc: "Implement dependency verification and supply chain security",
			Rationale: "Dependency assumptions protect against supply chain attacks and third-party compromise.",
			STRIDE: []StrideCategory{StrideDenialOfService, StrideTampering}, Priority: 2},
		{Category: "PROCESS", BaseDesc: "Implement audit logging and process verification",
			Rationale: "Process assumptions ensure accountability and non-repudiation of security-relevant actions.",
			STRIDE: []StrideCategory{StrideRepudiation, StrideTampering}, Priority: 2},
		{Category: "DATABASE", BaseDesc: "Implement database access controls and encryption",
			Rationale: "Database assumptions protect the confidentiality and integrity of stored data.",
			STRIDE: []StrideCategory{StrideTampering, StrideInfoDisclosure}, Priority: 2},
		{Category: "LOGGING", BaseDesc: "Implement immutable audit logging with tamper detection",
			Rationale: "Logging assumptions prevent repudiation and enable forensic investigation.",
			STRIDE: []StrideCategory{StrideRepudiation, StrideTampering}, Priority: 2},
		{Category: "BACKUP", BaseDesc: "Implement encrypted backup with tested restore procedures",
			Rationale: "Backup assumptions ensure data availability and recovery against ransomware and data loss.",
			STRIDE: []StrideCategory{StrideInfoDisclosure, StrideDenialOfService}, Priority: 2},
		{Category: "SESSION", BaseDesc: "Implement secure session management with rotation and timeout",
			Rationale: "Session assumptions prevent session hijacking and credential reuse attacks.",
			STRIDE: []StrideCategory{StrideSpoofing, StrideElevationPriv}, Priority: 2},
		{Category: "THIRD_PARTY", BaseDesc: "Implement third-party security assessment and monitoring",
			Rationale: "Third-party assumptions require vendor risk management to prevent supply chain attacks.",
			STRIDE: []StrideCategory{StrideTampering, StrideInfoDisclosure}, Priority: 2},
		{Category: "DOCUMENTATION", BaseDesc: "Maintain accurate and version-controlled architecture documentation",
			Rationale: "Documentation assumptions ensure knowledge continuity and accurate threat modeling.",
			STRIDE: []StrideCategory{StrideRepudiation}, Priority: 3},
		{Category: "GOVERNANCE", BaseDesc: "Establish security governance framework with regular reviews",
			Rationale: "Governance assumptions require oversight to maintain security posture over time.",
			STRIDE: []StrideCategory{StrideRepudiation, StrideTampering}, Priority: 3},
	}
}

func generateControls(assumptions []Assumption) []ControlDetail {
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
	if len(keywords) > 0 {
		return strings.Join(keywords[:min(3, len(keywords))], ", ")
	}
	words := strings.Fields(text)
	if len(words) > 5 {
		return words[0]
	}
	return "general"
}
