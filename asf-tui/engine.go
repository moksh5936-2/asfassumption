package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"asf-tui/asf/analyzer"
	"asf-tui/asf/models"
)

var debugLog = log.New(os.Stderr, "[asf-debug] ", log.Ltime|log.Lshortfile)

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

	// Evidence Traceability
	SourceNode string `json:"source_node"`
	SourceLine int    `json:"source_line"`

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
		desc, err := ParseArchitecture(archPath)
		if err == nil {
			e.archDesc = desc
		}
	}

	progress <- AnalysisProgress{Percent: 20, Stage: "Running ASF Engine..."}

	debugLog.Printf("using native Go engine")
	asfResult, err := e.runNativeAnalysis(inputPath, evPath)
	if err != nil {
		return nil, fmt.Errorf("ASF engine error: %w", err)
	}

	progress <- AnalysisProgress{Percent: 60, Stage: "Processing Results..."}

	result := e.buildResult(asfResult, archPath, mode)

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

func (e *Engine) runNativeAnalysis(docPath, evPath string) (*asfJSONResult, error) {
	docs := []string{docPath}
	var evs []string
	if evPath != "" {
		if _, err := os.Stat(evPath); err == nil {
			evs = append(evs, evPath)
		}
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

		_ = i
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

	// Validate expected results if present
	if len(e.archDesc.ExpectedResults) > 0 {
		result.Summary = e.buildValidationSummary(result)
	}

	result.Controls = generateControls(result.Assumptions)
	if len(e.archDesc.SecurityControls) > 0 {
		result.Controls = enhanceControlsWithSecurityControls(result.Controls, e.archDesc.SecurityControls)
	}

	// Build evidence summary
	if e.explainPipe != nil {
		result.EvidenceSummary = e.explainPipe.BuildEvidenceSummary(result.Assumptions)
		confSummary := buildConfidenceSummary(result.Assumptions)
		result.ConfidenceSummary = confSummary
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
func (e *Engine) processExplicitAssumptions(existing []Assumption) []Assumption {
	existingSet := make(map[string]bool)
	for _, a := range existing {
		existingSet[normalizeText(a.Description)] = true
	}

	var result []Assumption
	nextID := 1

	for _, raw := range e.archDesc.ExplicitAssumptions {
		normalized := normalizeText(raw)
		if normalized == "" {
			continue
		}
		if existingSet[normalized] {
			continue
		}
		existingSet[normalized] = true

		atype := classifyExplicitAssumption(raw)
		keywords := extractKeywords(raw)
		component := extractComponent(keywords, raw)

		id := fmt.Sprintf("ASM-%03d", nextID)
		nextID++

		risk := e.assessExplicitRisk(raw, atype)
		lh, im := riskToLikelihoodImpact(risk)
		stride := e.strideEngine.MapAssumption(string(atype), raw, keywords)

		assumption := Assumption{
			ID:          id,
			Description: raw,
			Component:   component,
			Category:    string(atype),
			Risk:        risk,
			Stride:      stride,
			Likelihood:  lh,
			Impact:      im,
			Confidence:  0.75,
			Keywords:    keywords,
		}

		if e.explainPipe != nil {
			e.explainPipe.Explain(&assumption)
		}

		result = append(result, assumption)
	}
	return result
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

	// PHI/healthcare keywords boost risk
	phiKws := []string{"phi", "health", "hipaa", "patient", "medical", "phi data", "protected health"}
	for _, kw := range phiKws {
		if strings.Contains(lower, kw) {
			score += 3
			break
		}
	}

	// High-severity keywords
	highKws := []string{"critical", "compromise", "breach", "unauthorized", "restricted", "immutable"}
	for _, kw := range highKws {
		if strings.Contains(lower, kw) {
			score += 2
			break
		}
	}

	// Medium-severity keywords
	medKws := []string{"encrypt", "kms", "access", "authenticate", "mfa", "token", "key", "audit", "backup", "monitor", "rate limit", "session"}
	for _, kw := range medKws {
		if strings.Contains(lower, kw) {
			score += 1
			break
		}
	}

	// Type-based boost
	switch atype {
	case models.AssumptionTypeACCESS:
		score += 1
	case models.AssumptionTypeIDENTITY:
		score += 1
	}

	switch {
	case score >= 5:
		return RiskCritical
	case score >= 3:
		return RiskHigh
	case score >= 2:
		return RiskMedium
	default:
		return RiskLow
	}
}

// buildComplianceOutput builds the compliance section from architecture description.
func (e *Engine) buildComplianceOutput() []string {
	if e.archDesc == nil || len(e.archDesc.Compliance) == 0 {
		return []string{
			"ASF analysis completed — see gap analysis for compliance mapping",
		}
	}

	compliance := e.archDesc.Compliance
	output := make([]string, 0, len(compliance)+2)
	output = append(output, "Compliance frameworks identified in architecture definition:")
	for _, c := range compliance {
		output = append(output, fmt.Sprintf("- %s related findings reviewed in this analysis", c))
	}
	output = append(output, "Review gap analysis for detailed compliance mapping.")
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
func normalizeText(text string) string {
	lower := strings.ToLower(text)
	// Collapse whitespace and strip periods from each token
	parts := strings.Fields(lower)
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
