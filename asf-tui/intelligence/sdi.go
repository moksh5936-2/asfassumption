package intelligence

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

// ═══════════════════════════════════════════════════════
// SDI — Security Decision Intelligence Engine (ASF V11)
// Phases 1-15
// ═══════════════════════════════════════════════════════

// ── PHASE 1 — DECISION MODEL ──

type DecisionRecommendation struct {
	ID                  string   `json:"id"`
	Title               string   `json:"title"`
	Description         string   `json:"description"`
	AffectedFindings    []string `json:"affected_findings,omitempty"`
	AffectedThreats     []string `json:"affected_threats,omitempty"`
	AffectedAttackPaths []string `json:"affected_attack_paths,omitempty"`
	AffectedControls    []string `json:"affected_controls,omitempty"`
	AffectedAssets      []string `json:"affected_assets,omitempty"`
	RiskReduction       string   `json:"risk_reduction"`
	Effort              string   `json:"effort"`
	Priority            string   `json:"priority"`
	BusinessImpact      string   `json:"business_impact"`
	ComplianceImpact    string   `json:"compliance_impact"`
	Rationale           string   `json:"rationale"`
}

// ── PHASE 5 — FIX SIMULATION ──

type FixSimulation struct {
	ControlName         string  `json:"control_name"`
	ControlCategory     string  `json:"control_category"`
	OriginalCritical    int     `json:"original_critical"`
	OriginalHigh        int     `json:"original_high"`
	OriginalTotal       int     `json:"original_total"`
	OriginalAttackPaths int     `json:"original_attack_paths"`
	OriginalCoverage    float64 `json:"original_coverage"`
	NewCritical         int     `json:"new_critical"`
	NewHigh             int     `json:"new_high"`
	NewTotal            int     `json:"new_total"`
	NewAttackPaths      int     `json:"new_attack_paths"`
	NewCoverage         float64 `json:"new_coverage"`
}

// ── PHASE 6 — FAILURE SIMULATION ──

type FailureSimulation struct {
	ControlName       string  `json:"control_name"`
	ControlCategory   string  `json:"control_category"`
	SystemsImpacted   int     `json:"systems_impacted"`
	AttackPathsOpened int     `json:"attack_paths_opened"`
	NewFindings       int     `json:"new_findings"`
	RiskIncrease      string  `json:"risk_increase"`
	RiskScoreIncrease float64 `json:"risk_score_increase"`
}

// ── PHASE 7 — CONTROL IMPACT ANALYSIS ──

type ControlImpact struct {
	ControlName     string `json:"control_name"`
	Category        string `json:"category"`
	SecurityValue   string `json:"security_value"`
	Effort          string `json:"effort"`
	ROI             string `json:"roi"`
	FindingCount    int    `json:"finding_count"`
	ThreatCount     int    `json:"threat_count"`
	AttackPathCount int    `json:"attack_path_count"`
}

// ── PHASE 8 — DECISION TREES ──

type DecisionTree struct {
	Budget           string                   `json:"budget"`
	ActionCount      int                      `json:"action_count"`
	RecommendedOrder []DecisionRecommendation `json:"recommended_order,omitempty"`
	Rationale        string                   `json:"rationale"`
}

type DecisionTreeResult struct {
	SingleAction DecisionTree `json:"single_action"`
	ThreeActions DecisionTree `json:"three_actions"`
	FiveActions  DecisionTree `json:"five_actions"`
}

// ── PHASE 9 — BOARD DECISION SUPPORT ──

type BoardScenario struct {
	Scenario         string   `json:"scenario"`
	Description      string   `json:"description"`
	RiskScore        float64  `json:"risk_score"`
	CriticalFindings int      `json:"critical_findings"`
	AttackPaths      int      `json:"attack_paths"`
	CoverageRate     float64  `json:"coverage_rate"`
	KeyRisks         []string `json:"key_risks,omitempty"`
}

type BoardScenarios struct {
	DoNothing        BoardScenario `json:"do_nothing"`
	PartialRemediate BoardScenario `json:"partial_remediate"`
	FullRemediate    BoardScenario `json:"full_remediate"`
}

// ── PHASE 10 — SECURITY INVESTMENT PRIORITIES ──

type InvestmentPriority struct {
	Area          string  `json:"area"`
	Rank          int     `json:"rank"`
	Score         float64 `json:"score"`
	Rationale     string  `json:"rationale"`
	FindingCount  int     `json:"finding_count"`
	RiskReduction string  `json:"risk_reduction"`
}

// ── PHASE 11 — ATTACK PATH COLLAPSE ANALYSIS ──

type AttackPathCollapse struct {
	ControlName        string  `json:"control_name"`
	Category           string  `json:"category"`
	AttackPathsReduced int     `json:"attack_paths_reduced"`
	TotalAttackPaths   int     `json:"total_attack_paths"`
	ReductionPercent   float64 `json:"reduction_percent"`
}

// ── PHASE 12 — COMPLIANCE IMPACT ──

type ComplianceImpact struct {
	Framework   string `json:"framework"`
	Action      string `json:"action"`
	Improvement string `json:"improvement"`
	Rationale   string `json:"rationale"`
}

// ── PHASE 13 — ENTERPRISE REMEDIATION ROADMAP ──

type SDIRoadmapItem struct {
	Action        string `json:"action"`
	Category      string `json:"category"`
	Priority      string `json:"priority"`
	Effort        string `json:"effort"`
	RiskReduction string `json:"risk_reduction"`
}

type SDIRemediationRoadmap struct {
	Phase30  []SDIRoadmapItem `json:"phase_30,omitempty"`
	Phase90  []SDIRoadmapItem `json:"phase_90,omitempty"`
	Phase180 []SDIRoadmapItem `json:"phase_180,omitempty"`
	Phase12m []SDIRoadmapItem `json:"phase_12m,omitempty"`
}

// ── PHASE 14 — DECISION DASHBOARD ──

type DecisionDashboard struct {
	TopDecisions         []DecisionRecommendation `json:"top_decisions,omitempty"`
	QuickWins            []DecisionRecommendation `json:"quick_wins,omitempty"`
	StrategicActions     []DecisionRecommendation `json:"strategic_actions,omitempty"`
	RiskReductionSummary string                   `json:"risk_reduction_summary"`
	TotalRiskReduction   float64                  `json:"total_risk_reduction"`
}

// ── PHASE 15 — EXECUTIVE SCENARIOS ──

type ExecutiveScenario struct {
	Scenario          string  `json:"scenario"`
	RiskScore         float64 `json:"risk_score"`
	FindingsResolved  int     `json:"findings_resolved"`
	AttackPathsClosed int     `json:"attack_paths_closed"`
	CoverageAchieved  float64 `json:"coverage_achieved"`
	Description       string  `json:"description"`
}

type ExecutiveScenarios struct {
	BestCase   ExecutiveScenario `json:"best_case"`
	LikelyCase ExecutiveScenario `json:"likely_case"`
	WorstCase  ExecutiveScenario `json:"worst_case"`
}

// ── SDI INPUT / RESULT ──

type SDIInput struct {
	ArchitectureName   string
	Domain             string
	Findings           []SDRIFinding
	Threats            []Threat
	AttackPaths        []AttackPath
	Controls           []SDRIControl
	Compliance         []string
	RiskScore          float64
	CoverageByCategory map[string]float64
	Assumptions        []Assumption
	AnalysisMode       string
}

type SDIResult struct {
	Recommendations      []DecisionRecommendation `json:"recommendations,omitempty"`
	FixSimulations       []FixSimulation          `json:"fix_simulations,omitempty"`
	FailureSimulations   []FailureSimulation      `json:"failure_simulations,omitempty"`
	ControlImpacts       []ControlImpact          `json:"control_impacts,omitempty"`
	DecisionTrees        DecisionTreeResult       `json:"decision_trees"`
	BoardScenarios       BoardScenarios           `json:"board_scenarios"`
	InvestmentPriorities []InvestmentPriority     `json:"investment_priorities,omitempty"`
	AttackPathCollapse   []AttackPathCollapse     `json:"attack_path_collapse,omitempty"`
	ComplianceImpacts    []ComplianceImpact       `json:"compliance_impacts,omitempty"`
	RemediationRoadmap   SDIRemediationRoadmap    `json:"remediation_roadmap"`
	Dashboard            DecisionDashboard        `json:"dashboard"`
	ExecutiveScenarios   ExecutiveScenarios       `json:"executive_scenarios"`
}

// Control-to-finding category mapping for risk reduction analysis.
var controlCategoryMap = map[string]string{
	"MFA":                        "AccessControl",
	"PasswordPolicy":             "Authentication",
	"SessionManagement":          "Authentication",
	"ConditionalAccess":          "Authorization",
	"RBAC":                       "Authorization",
	"ABAC":                       "Authorization",
	"JustInTimeAccess":           "Authorization",
	"PrivilegedAccessManagement": "Authorization",
	"IdentityGovernance":         "IdentityGovernance",
	"AccessReviews":              "IdentityGovernance",
	"SecretsRotation":            "SecretsManagement",
	"SecretsVault":               "SecretsManagement",
	"SecretsScanning":            "SecretsManagement",
	"KeyRotation":                "KeyManagement",
	"KeyAccessLogging":           "KeyManagement",
	"SeparationOfDuties":         "KeyManagement",
	"NetworkSegmentation":        "NetworkSecurity",
	"FirewallRules":              "NetworkSecurity",
	"TLSEncryption":              "DataProtection",
	"VPNAccess":                  "NetworkSecurity",
	"IntrusionDetection":         "NetworkSecurity",
	"AuditLogging":               "Logging",
	"SIEMIntegration":            "Monitoring",
	"EncryptionAtRest":           "DataProtection",
	"EndpointProtection":         "EndpointSecurity",
	"PatchManagement":            "VulnerabilityManagement",
	"BackupAndRecovery":          "Resilience",
	"IncidentResponse":           "Governance",
	"SupplyChainSecurity":        "ThirdPartyRisk",
}

// Effort estimates per control name (Low / Medium / High).
var controlEffortMap = map[string]string{
	"MFA":                        "Low",
	"PasswordPolicy":             "Low",
	"SessionManagement":          "Low",
	"TLSEncryption":              "Low",
	"AuditLogging":               "Low",
	"ConditionalAccess":          "Medium",
	"RBAC":                       "Medium",
	"ABAC":                       "Medium",
	"JustInTimeAccess":           "Medium",
	"SecretsRotation":            "Medium",
	"SecretsVault":               "Medium",
	"SecretsScanning":            "Medium",
	"KeyRotation":                "Medium",
	"KeyAccessLogging":           "Low",
	"SeparationOfDuties":         "High",
	"NetworkSegmentation":        "High",
	"FirewallRules":              "Medium",
	"VPNAccess":                  "Medium",
	"IntrusionDetection":         "Medium",
	"SIEMIntegration":            "Medium",
	"EncryptionAtRest":           "Medium",
	"EndpointProtection":         "Medium",
	"PatchManagement":            "Medium",
	"BackupAndRecovery":          "Medium",
	"IncidentResponse":           "Medium",
	"SupplyChainSecurity":        "High",
	"IdentityGovernance":         "Medium",
	"AccessReviews":              "Low",
	"PrivilegedAccessManagement": "High",
}

// Compliance framework mappings — which controls improve which frameworks.
var complianceControlMap = map[string][]string{
	"HIPAA":       {"AuditLogging", "EncryptionAtRest", "TLSEncryption", "AccessReviews", "MFA", "BackupAndRecovery", "IncidentResponse"},
	"SOC 2":       {"AuditLogging", "SIEMIntegration", "AccessReviews", "MFA", "RBAC", "EncryptionAtRest", "PatchManagement"},
	"PCI DSS":     {"MFA", "EncryptionAtRest", "TLSEncryption", "AuditLogging", "FirewallRules", "NetworkSegmentation", "AccessReviews"},
	"GDPR":        {"EncryptionAtRest", "TLSEncryption", "AuditLogging", "AccessReviews", "BackupAndRecovery", "IncidentResponse"},
	"FedRAMP":     {"MFA", "AuditLogging", "SIEMIntegration", "IncidentResponse", "RBAC", "EncryptionAtRest", "PatchManagement"},
	"ISO 27001":   {"AuditLogging", "IncidentResponse", "AccessReviews", "RBAC", "PatchManagement", "BackupAndRecovery"},
	"NIST 800-53": {"MFA", "AuditLogging", "SIEMIntegration", "IncidentResponse", "RBAC", "EncryptionAtRest", "NetworkSegmentation"},
	"HITRUST":     {"MFA", "AuditLogging", "EncryptionAtRest", "TLSEncryption", "AccessReviews", "IncidentResponse", "PatchManagement"},
}

// defaultRecommendations is the list of canonical decision actions the engine evaluates.
var defaultRecommendations = []struct {
	ID          string
	Title       string
	Category    string
	Description string
	Effort      string
}{
	{"SDI-R-01", "Enable Multi-Factor Authentication", "AccessControl", "Implement MFA for all user access to critical systems and administrative interfaces", "Low"},
	{"SDI-R-02", "Implement Encryption at Rest", "DataProtection", "Encrypt sensitive data stored in databases, storage systems, and backups", "Medium"},
	{"SDI-R-03", "Enable Comprehensive Audit Logging", "Logging", "Enable audit logging for all security-relevant events on critical systems", "Low"},
	{"SDI-R-04", "Deploy SIEM Integration", "Monitoring", "Deploy SIEM for centralized security event monitoring, correlation, and alerting", "Medium"},
	{"SDI-R-05", "Implement Network Segmentation", "NetworkSecurity", "Segment the network into isolated zones based on trust levels and data sensitivity", "High"},
	{"SDI-R-06", "Deploy Intrusion Detection/Prevention", "NetworkSecurity", "Deploy IDS/IPS to detect and prevent network-based attacks and anomalies", "Medium"},
	{"SDI-R-07", "Implement Role-Based Access Control", "Authorization", "Enforce least privilege through RBAC across all systems and applications", "Medium"},
	{"SDI-R-08", "Establish Centralized Secrets Management", "SecretsManagement", "Implement a centralized secrets vault with automated rotation policies", "Medium"},
	{"SDI-R-09", "Enforce TLS for All Service Communication", "DataProtection", "Enforce TLS 1.2+ encryption for all service-to-service and external communications", "Low"},
	{"SDI-R-10", "Implement Privileged Access Management", "Authorization", "Deploy PAM for managing, monitoring, and auditing privileged account usage", "High"},
	{"SDI-R-11", "Deploy Endpoint Detection and Response", "EndpointSecurity", "Deploy EDR on all endpoints for real-time threat detection and automated response", "Medium"},
	{"SDI-R-12", "Establish Incident Response Plan", "Governance", "Create and maintain an incident response plan with regular tabletop exercises", "Medium"},
	{"SDI-R-13", "Implement Vulnerability Management Program", "VulnerabilityManagement", "Deploy automated vulnerability scanning with prioritized patch management", "Medium"},
	{"SDI-R-14", "Enforce Strong Password Policy", "Authentication", "Implement strong password complexity, length requirements, and rotation policies", "Low"},
	{"SDI-R-15", "Implement Just-In-Time Access", "Authorization", "Provide temporary privileged access on demand with automatic revocation", "Medium"},
	{"SDI-R-16", "Deploy Backup and Disaster Recovery", "Resilience", "Implement automated backups with tested disaster recovery procedures", "Medium"},
	{"SDI-R-17", "Establish Supply Chain Security Program", "ThirdPartyRisk", "Implement vendor risk assessment, monitoring, and security requirements", "High"},
	{"SDI-R-18", "Deploy Identity Governance", "IdentityGovernance", "Implement identity lifecycle management with certification and review processes", "Medium"},
	{"SDI-R-19", "Conduct Regular Access Reviews", "IdentityGovernance", "Perform periodic access reviews and certifications for all user entitlements", "Low"},
	{"SDI-R-20", "Implement Patch Management Process", "VulnerabilityManagement", "Establish a formal patch management process with defined SLAs for critical patches", "Medium"},
}

// ── SDI ENGINE ──

type SDIEngine struct{}

func NewSDIEngine() *SDIEngine {
	return &SDIEngine{}
}

func (e *SDIEngine) Run(input SDIInput) *SDIResult {
	r := &SDIResult{}

	r.Recommendations = e.generateRecommendations(input)
	r.FixSimulations = e.generateFixSimulations(input)
	r.FailureSimulations = e.generateFailureSimulations(input)
	r.ControlImpacts = e.generateControlImpacts(input)
	r.DecisionTrees = e.generateDecisionTrees(r.Recommendations)
	r.BoardScenarios = e.generateBoardScenarios(input, r.Recommendations)
	r.InvestmentPriorities = e.generateInvestmentPriorities(r.Recommendations)
	r.AttackPathCollapse = e.generateAttackPathCollapse(input)
	r.ComplianceImpacts = e.generateComplianceImpacts(input)
	r.RemediationRoadmap = e.generateRemediationRoadmap(r.Recommendations)
	r.Dashboard = e.generateDashboard(r.Recommendations)
	r.ExecutiveScenarios = e.generateExecutiveScenarios(input, r.Recommendations)

	return r
}

// ── PHASE 2 — RISK REDUCTION ENGINE ──

func (e *SDIEngine) generateRecommendations(input SDIInput) []DecisionRecommendation {
	var recs []DecisionRecommendation

	for _, def := range defaultRecommendations {
		affectedFindings := e.findingsForCategory(input.Findings, def.Category)
		affectedThreats := e.threatsForCategory(input.Threats, def.Category)
		affectedPaths := e.pathsForControl(input.AttackPaths, def.Title)
		affectedAssets := e.assetsForFindings(input.Findings, affectedFindings)

		if len(affectedFindings) == 0 && len(affectedThreats) == 0 {
			continue
		}

		riskReduction := e.computeRiskReduction(affectedFindings)
		effort := def.Effort
		if ce, ok := controlEffortMap[extractControlName(def.Title)]; ok {
			effort = ce
		}
		priority := e.computePriority(affectedFindings, affectedThreats, riskReduction)
		bizImpact := e.computeBusinessImpact(affectedFindings)
		compImpact := e.computeComplianceImpact(input.Compliance, def.Category)

		var findingIDs []string
		for _, f := range affectedFindings {
			findingIDs = append(findingIDs, f.ID)
		}
		var threatIDs []string
		for _, t := range affectedThreats {
			threatIDs = append(threatIDs, t.ID)
		}
		var pathIDs []string
		for _, p := range affectedPaths {
			pathIDs = append(pathIDs, p.ID)
		}
		var ctrlNames []string
		for _, c := range input.Controls {
			if controlCategoryMap[c.Name] == def.Category {
				ctrlNames = append(ctrlNames, c.Name)
			}
		}

		rationale := fmt.Sprintf("Addresses %d findings (%d Critical, %d High) and %d threats in the %s category.",
			len(affectedFindings), countSeverity(affectedFindings, "Critical"),
			countSeverity(affectedFindings, "High"), len(affectedThreats), def.Category)

		recs = append(recs, DecisionRecommendation{
			ID:                  def.ID,
			Title:               def.Title,
			Description:         def.Description,
			AffectedFindings:    findingIDs,
			AffectedThreats:     threatIDs,
			AffectedAttackPaths: pathIDs,
			AffectedControls:    ctrlNames,
			AffectedAssets:      affectedAssets,
			RiskReduction:       riskReduction,
			Effort:              effort,
			Priority:            priority,
			BusinessImpact:      bizImpact,
			ComplianceImpact:    compImpact,
			Rationale:           rationale,
		})
	}

	sort.Slice(recs, func(i, j int) bool {
		return priorityScore(recs[i].Priority) > priorityScore(recs[j].Priority)
	})

	return recs
}

func (e *SDIEngine) findingsForCategory(findings []SDRIFinding, category string) []SDRIFinding {
	var out []SDRIFinding
	for _, f := range findings {
		lower := strings.ToLower(f.Category)
		catLower := strings.ToLower(category)
		titleLower := strings.ToLower(f.Title)
		if strings.Contains(lower, catLower) ||
			strings.Contains(titleLower, catLower) ||
			categoryMatch(f.Category, category) {
			out = append(out, f)
		}
	}
	return out
}

func (e *SDIEngine) threatsForCategory(threats []Threat, category string) []Threat {
	var out []Threat
	for _, t := range threats {
		lower := strings.ToLower(string(t.Category))
		catLower := strings.ToLower(category)
		if strings.Contains(lower, catLower) {
			out = append(out, t)
		}
	}
	return out
}

func (e *SDIEngine) pathsForControl(paths []AttackPath, title string) []AttackPath {
	var out []AttackPath
	for _, p := range paths {
		lower := strings.ToLower(p.Name + " " + p.Description)
		for _, kw := range controlKeywords(title) {
			if strings.Contains(lower, kw) {
				out = append(out, p)
				break
			}
		}
	}
	return out
}

func (e *SDIEngine) assetsForFindings(findings []SDRIFinding, matched []SDRIFinding) []string {
	seen := map[string]bool{}
	var assets []string
	for _, f := range matched {
		for _, f2 := range findings {
			if f2.ID == f.ID {
				for _, comp := range f2.AffectedComponents {
					if !seen[comp] {
						seen[comp] = true
						assets = append(assets, comp)
					}
				}
			}
		}
	}
	return assets
}

func (e *SDIEngine) computeRiskReduction(findings []SDRIFinding) string {
	for _, f := range findings {
		if strings.EqualFold(f.Severity, "Critical") {
			return "High"
		}
	}
	for _, f := range findings {
		if strings.EqualFold(f.Severity, "High") {
			return "Medium"
		}
	}
	if len(findings) > 0 {
		return "Low"
	}
	return "None"
}

func (e *SDIEngine) computePriority(findings []SDRIFinding, threats []Threat, riskReduction string) string {
	criticalCount := countSeverity(findings, "Critical")
	highCount := countSeverity(findings, "High")
	threatCritical := 0
	for _, t := range threats {
		if strings.EqualFold(string(t.Severity), "Critical") {
			threatCritical++
		}
	}

	switch {
	case criticalCount >= 2 || threatCritical >= 3 || riskReduction == "High":
		return "Critical"
	case criticalCount >= 1 || highCount >= 3 || threatCritical >= 1:
		return "High"
	case highCount >= 1 || len(findings) >= 3:
		return "Medium"
	default:
		return "Low"
	}
}

func (e *SDIEngine) computeBusinessImpact(findings []SDRIFinding) string {
	for _, f := range findings {
		if strings.EqualFold(f.Severity, "Critical") {
			return "Critical business impact — addresses systemic vulnerabilities"
		}
	}
	for _, f := range findings {
		if strings.EqualFold(f.Severity, "High") {
			return "Significant business impact — reduces major risk exposure"
		}
	}
	if len(findings) > 0 {
		return "Moderate business impact — improves security posture"
	}
	return "Limited direct business impact"
}

func (e *SDIEngine) computeComplianceImpact(frameworks []string, category string) string {
	var matched []string
	for _, fw := range frameworks {
		if ctrls, ok := complianceControlMap[fw]; ok {
			for _, ctrl := range ctrls {
				if controlCategoryMap[ctrl] == category {
					matched = append(matched, fw)
					break
				}
			}
		}
	}
	if len(matched) >= 3 {
		return fmt.Sprintf("Improves compliance for %d frameworks: %s", len(matched), strings.Join(matched, ", "))
	}
	if len(matched) >= 1 {
		return fmt.Sprintf("Supports compliance for %s", strings.Join(matched, ", "))
	}
	return "No direct compliance impact identified"
}

// ── PHASE 3 — SECURITY ROI ENGINE ──

func computeROI(effort, riskReduction string) string {
	effortScore := effortWeight(effort)
	reductionScore := reductionWeight(riskReduction)

	ratio := reductionScore / math.Max(effortScore, 0.1)

	switch {
	case ratio >= 2.5:
		return "Excellent"
	case ratio >= 1.25:
		return "Good"
	case ratio >= 0.75:
		return "Fair"
	default:
		return "Limited"
	}
}

func effortWeight(effort string) float64 {
	switch strings.ToLower(effort) {
	case "low":
		return 1.0
	case "medium":
		return 2.0
	case "high":
		return 3.0
	default:
		return 2.0
	}
}

func reductionWeight(riskReduction string) float64 {
	switch riskReduction {
	case "High":
		return 3.0
	case "Medium":
		return 2.0
	case "Low":
		return 1.0
	default:
		return 0.5
	}
}

// ── PHASE 4 — PRIORITIZATION ENGINE ──

func priorityScore(p string) int {
	switch p {
	case "Critical":
		return 4
	case "High":
		return 3
	case "Medium":
		return 2
	case "Low":
		return 1
	default:
		return 0
	}
}

// ── PHASE 5 — FIX SIMULATION ──

func (e *SDIEngine) generateFixSimulations(input SDIInput) []FixSimulation {
	var sims []FixSimulation
	origCritical := countSeverity(input.Findings, "Critical")
	origHigh := countSeverity(input.Findings, "High")

	simulated := map[string]bool{}
	for _, def := range defaultRecommendations {
		if simulated[def.Category] {
			continue
		}
		simulated[def.Category] = true

		affected := e.findingsForCategory(input.Findings, def.Category)
		if len(affected) == 0 {
			continue
		}

		newCritical := origCritical
		newHigh := origHigh
		var pathImpact int
		for _, f := range affected {
			if strings.EqualFold(f.Severity, "Critical") {
				newCritical--
			} else if strings.EqualFold(f.Severity, "High") {
				newHigh--
			}
		}
		for _, p := range input.AttackPaths {
			if containsAny(strings.ToLower(p.Name+" "+p.Description),
				strings.Fields(strings.ToLower(def.Category))) {
				pathImpact++
			}
		}

		if newCritical < 0 {
			newCritical = 0
		}
		if newHigh < 0 {
			newHigh = 0
		}

		origCoverage := computeCoverage(input.Findings, input.Controls)
		newTotal := len(input.Findings) - len(affected)
		if newTotal < 0 {
			newTotal = 0
		}
		var newCoverage float64
		if len(input.Findings) > 0 {
			simControls := input.Controls
			newCoverage = computeCoverage(input.Findings, simControls)
			if newTotal > 0 {
				covered := int(float64(len(input.Controls)) * origCoverage / 100.0)
				newCoverage = float64(covered) / float64(len(input.Controls)) * 100.0
				if newCoverage > 100 {
					newCoverage = 100
				}
			}
		}

		sims = append(sims, FixSimulation{
			ControlName:         def.Title,
			ControlCategory:     def.Category,
			OriginalCritical:    origCritical,
			OriginalHigh:        origHigh,
			OriginalTotal:       len(input.Findings),
			OriginalAttackPaths: len(input.AttackPaths),
			OriginalCoverage:    origCoverage,
			NewCritical:         newCritical,
			NewHigh:             newHigh,
			NewTotal:            newTotal,
			NewAttackPaths:      len(input.AttackPaths) - pathImpact,
			NewCoverage:         newCoverage,
		})
	}
	return sims
}

// ── PHASE 6 — FAILURE SIMULATION ──

func (e *SDIEngine) generateFailureSimulations(input SDIInput) []FailureSimulation {
	var sims []FailureSimulation
	uniqueCategories := map[string]bool{}
	for _, c := range input.Controls {
		if uniqueCategories[c.Category] {
			continue
		}
		uniqueCategories[c.Category] = true

		catFindings := e.findingsForCategory(input.Findings, c.Category)
		potentialFindings := 0
		for _, f := range input.Findings {
			if !containsAny(strings.ToLower(f.Category), strings.Fields(strings.ToLower(c.Category))) {
				if strings.Contains(strings.ToLower(f.Title), strings.ToLower(c.Name)) {
					potentialFindings++
				}
			}
		}

		pathsOpened := 0
		for _, p := range input.AttackPaths {
			lower := strings.ToLower(p.Name + " " + p.Description)
			catLower := strings.ToLower(c.Category)
			if strings.Contains(lower, "missing") || (!strings.Contains(lower, catLower) &&
				strings.Contains(lower, strings.ToLower(c.Name))) {
				pathsOpened++
			}
		}

		riskIncrease := "Low"
		riskScoreIncrease := 1.0
		if len(catFindings) >= 3 {
			riskIncrease = "High"
			riskScoreIncrease = 3.0 + float64(len(catFindings))*0.5
		} else if len(catFindings) >= 1 {
			riskIncrease = "Medium"
			riskScoreIncrease = 1.5 + float64(len(catFindings))*0.3
		}

		controlName := c.Name
		var affectedComps []string
		for _, f := range catFindings {
			affectedComps = append(affectedComps, f.AffectedComponents...)
		}
		systemsImpacted := len(affectedComps)
		if systemsImpacted == 0 {
			systemsImpacted = len(catFindings) + 1
		}

		sims = append(sims, FailureSimulation{
			ControlName:       controlName,
			ControlCategory:   c.Category,
			SystemsImpacted:   systemsImpacted,
			AttackPathsOpened: pathsOpened,
			NewFindings:       potentialFindings,
			RiskIncrease:      riskIncrease,
			RiskScoreIncrease: riskScoreIncrease,
		})
	}
	return sims
}

// ── PHASE 7 — CONTROL IMPACT ANALYSIS ──

func (e *SDIEngine) generateControlImpacts(input SDIInput) []ControlImpact {
	var impacts []ControlImpact
	controlSeen := map[string]bool{}
	for _, c := range input.Controls {
		if controlSeen[c.Name] {
			continue
		}
		controlSeen[c.Name] = true

		cat := c.Category
		if mapped, ok := controlCategoryMap[c.Name]; ok {
			cat = mapped
		}

		findings := e.findingsForCategory(input.Findings, cat)
		threats := e.threatsForCategory(input.Threats, cat)
		var pathCount int
		for _, p := range input.AttackPaths {
			lower := strings.ToLower(p.Name + " " + p.Description)
			catLower := strings.ToLower(cat)
			if strings.Contains(lower, catLower) || strings.Contains(lower, strings.ToLower(c.Name)) {
				pathCount++
			}
		}

		riskReduction := e.computeRiskReduction(findings)
		effort := "Medium"
		if ce, ok := controlEffortMap[c.Name]; ok {
			effort = ce
		}
		roi := computeROI(effort, riskReduction)

		securityValue := "Medium"
		switch {
		case riskReduction == "High" && len(threats) >= 2:
			securityValue = "Critical"
		case riskReduction == "High" || len(findings) >= 3:
			securityValue = "High"
		case len(findings) == 0 && len(threats) == 0:
			securityValue = "Low"
		}

		impacts = append(impacts, ControlImpact{
			ControlName:     c.Name,
			Category:        cat,
			SecurityValue:   securityValue,
			Effort:          effort,
			ROI:             roi,
			FindingCount:    len(findings),
			ThreatCount:     len(threats),
			AttackPathCount: pathCount,
		})
	}

	sort.Slice(impacts, func(i, j int) bool {
		return valueScore(impacts[i].SecurityValue) > valueScore(impacts[j].SecurityValue)
	})

	return impacts
}

func valueScore(v string) int {
	switch v {
	case "Critical":
		return 4
	case "High":
		return 3
	case "Medium":
		return 2
	case "Low":
		return 1
	default:
		return 0
	}
}

// ── PHASE 8 — DECISION TREES ──

func (e *SDIEngine) generateDecisionTrees(recs []DecisionRecommendation) DecisionTreeResult {
	critical := filterPriority(recs, "Critical")
	high := filterPriority(recs, "High")
	medium := filterPriority(recs, "Medium")

	var ordered []DecisionRecommendation
	ordered = append(ordered, critical...)
	ordered = append(ordered, high...)
	ordered = append(ordered, medium...)

	if len(ordered) > 5 {
		ordered = ordered[:5]
	}

	single := DecisionTree{
		Budget:      "Single Action",
		ActionCount: 1,
		Rationale:   "Highest priority action with maximum risk reduction per effort.",
	}
	if len(ordered) > 0 {
		single.RecommendedOrder = []DecisionRecommendation{ordered[0]}
	}

	three := DecisionTree{
		Budget:      "Three Actions",
		ActionCount: 3,
		Rationale:   "Top three actions balancing critical risks, quick wins, and compliance.",
	}
	if len(ordered) >= 3 {
		three.RecommendedOrder = ordered[:3]
	} else {
		three.RecommendedOrder = ordered
	}

	five := DecisionTree{
		Budget:      "Five Actions",
		ActionCount: 5,
		Rationale:   "Comprehensive remediation addressing critical risks, compliance gaps, and strategic improvements.",
	}
	if len(ordered) >= 5 {
		five.RecommendedOrder = ordered[:5]
	} else {
		five.RecommendedOrder = ordered
	}

	return DecisionTreeResult{
		SingleAction: single,
		ThreeActions: three,
		FiveActions:  five,
	}
}

// ── PHASE 9 — BOARD DECISION SUPPORT ──

func (e *SDIEngine) generateBoardScenarios(input SDIInput, recs []DecisionRecommendation) BoardScenarios {
	currentRisk := input.RiskScore
	criticalCount := countSeverity(input.Findings, "Critical")

	partialReduction := 0
	fullReduction := 0
	for _, r := range recs {
		if r.RiskReduction == "High" {
			fullReduction += 2
			if r.Priority == "Critical" || r.Priority == "High" {
				partialReduction++
			}
		} else if r.RiskReduction == "Medium" {
			fullReduction++
			if r.Priority == "Critical" {
				partialReduction++
			}
		}
	}

	partialCritical := criticalCount - partialReduction
	if partialCritical < 0 {
		partialCritical = 0
	}
	fullCritical := criticalCount - fullReduction
	if fullCritical < 0 {
		fullCritical = 0
	}

	return BoardScenarios{
		DoNothing: BoardScenario{
			Scenario:         "Do Nothing",
			Description:      "Maintain current security posture with no additional investment.",
			RiskScore:        currentRisk,
			CriticalFindings: criticalCount,
			AttackPaths:      len(input.AttackPaths),
			CoverageRate:     computeCoverage(input.Findings, input.Controls),
			KeyRisks:         []string{fmt.Sprintf("%d Critical findings remain unaddressed", criticalCount), fmt.Sprintf("%d attack paths expose critical assets", len(input.AttackPaths))},
		},
		PartialRemediate: BoardScenario{
			Scenario:         "Partial Remediation (Top 3 Actions)",
			Description:      "Implement the highest priority recommendations to address the most critical risks.",
			RiskScore:        currentRisk - float64(partialReduction)*0.3,
			CriticalFindings: partialCritical,
			AttackPaths:      len(input.AttackPaths) - partialReduction,
			CoverageRate:     computeCoverage(input.Findings, input.Controls) + float64(partialReduction)*3,
			KeyRisks:         []string{fmt.Sprintf("%d Critical findings remain", partialCritical), "Key compliance gaps partially addressed"},
		},
		FullRemediate: BoardScenario{
			Scenario:         "Full Remediation (All Actions)",
			Description:      "Implement all recommended security controls to achieve optimal security posture.",
			RiskScore:        currentRisk - float64(fullReduction)*0.4,
			CriticalFindings: fullCritical,
			AttackPaths:      len(input.AttackPaths) - fullReduction,
			CoverageRate:     computeCoverage(input.Findings, input.Controls) + float64(fullReduction)*5,
			KeyRisks:         []string{fmt.Sprintf("Only %d Critical findings remain", fullCritical), "Comprehensive compliance coverage achieved"},
		},
	}
}

// ── PHASE 10 — SECURITY INVESTMENT PRIORITIES ──

func (e *SDIEngine) generateInvestmentPriorities(recs []DecisionRecommendation) []InvestmentPriority {
	type areaData struct {
		findings  int
		score     float64
		reduction string
		count     int
	}
	areas := map[string]*areaData{}

	for _, r := range recs {
		area := categoryFromTitle(r.Title)
		if _, ok := areas[area]; !ok {
			areas[area] = &areaData{}
		}
		areas[area].findings += len(r.AffectedFindings)
		areas[area].count++
		rs := reductionWeight(r.RiskReduction)
		ps := priorityScore(r.Priority)
		areas[area].score += rs * float64(ps)
		if r.RiskReduction == "High" && r.Priority == "Critical" {
			areas[area].reduction = "High"
		}
	}

	var sorted []InvestmentPriority
	i := 1
	orderedAreas := []string{"Identity", "Data Protection", "Monitoring", "Third Party Risk", "Cloud Security",
		"Network Security", "Endpoint Security", "Application Security", "Governance", "Resilience"}

	for _, area := range orderedAreas {
		if d, ok := areas[area]; ok && d.count > 0 {
			rationale := ""
			switch area {
			case "Identity":
				rationale = "Identity is the new perimeter — MFA and PAM provide the highest risk reduction per effort"
			case "Data Protection":
				rationale = "Encryption and TLS protect sensitive data at rest and in transit"
			case "Monitoring":
				rationale = "SIEM and logging provide visibility for threat detection and incident response"
			case "Third Party Risk":
				rationale = "Supply chain attacks require vendor security assessment and monitoring"
			case "Cloud Security":
				rationale = "Cloud configuration and workload protection for cloud-native deployments"
			case "Network Security":
				rationale = "Segmentation and IDS limit lateral movement and detect network threats"
			case "Endpoint Security":
				rationale = "EDR provides real-time detection and response on all endpoints"
			case "Application Security":
				rationale = "Application-level controls protect against OWASP top 10 vulnerabilities"
			case "Governance":
				rationale = "Incident response and access reviews ensure program maturity"
			case "Resilience":
				rationale = "Backup and recovery ensure business continuity"
			}
			if d.reduction == "" {
				d.reduction = "Medium"
			}
			sorted = append(sorted, InvestmentPriority{
				Area:          area,
				Rank:          i,
				Score:         d.score,
				Rationale:     rationale,
				FindingCount:  d.findings,
				RiskReduction: d.reduction,
			})
			i++
		}
	}

	for area, d := range areas {
		found := false
		for _, s := range sorted {
			if s.Area == area {
				found = true
				break
			}
		}
		if !found {
			if d.reduction == "" {
				d.reduction = "Medium"
			}
			sorted = append(sorted, InvestmentPriority{
				Area:          area,
				Rank:          i,
				Score:         d.score,
				Rationale:     fmt.Sprintf("%s controls address %d findings with measurable risk reduction", area, d.findings),
				FindingCount:  d.findings,
				RiskReduction: d.reduction,
			})
			i++
		}
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Score > sorted[j].Score
	})
	for idx := range sorted {
		sorted[idx].Rank = idx + 1
	}

	return sorted
}

// ── PHASE 11 — ATTACK PATH COLLAPSE ANALYSIS ──

func (e *SDIEngine) generateAttackPathCollapse(input SDIInput) []AttackPathCollapse {
	var results []AttackPathCollapse
	totalPaths := len(input.AttackPaths)
	if totalPaths == 0 {
		return results
	}

	for _, def := range defaultRecommendations {
		affected := e.pathsForControl(input.AttackPaths, def.Title)
		reduced := len(affected)
		if reduced == 0 {
			continue
		}
		pct := float64(reduced) / float64(totalPaths) * 100.0
		results = append(results, AttackPathCollapse{
			ControlName:        def.Title,
			Category:           def.Category,
			AttackPathsReduced: reduced,
			TotalAttackPaths:   totalPaths,
			ReductionPercent:   pct,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].AttackPathsReduced > results[j].AttackPathsReduced
	})

	if len(results) > 10 {
		results = results[:10]
	}

	return results
}

// ── PHASE 12 — COMPLIANCE IMPACT ──

func (e *SDIEngine) generateComplianceImpacts(input SDIInput) []ComplianceImpact {
	var impacts []ComplianceImpact
	for _, fw := range input.Compliance {
		if ctrls, ok := complianceControlMap[fw]; ok {
			for _, ctrl := range ctrls {
				cat := controlCategoryMap[ctrl]
				findings := e.findingsForCategory(input.Findings, cat)
				improvement := "None"
				if len(findings) > 0 {
					criticalCount := countSeverity(findings, "Critical")
					switch {
					case criticalCount >= 2:
						improvement = "Significant"
					case criticalCount >= 1 || len(findings) >= 3:
						improvement = "Moderate"
					default:
						improvement = "Minor"
					}
				}
				impacts = append(impacts, ComplianceImpact{
					Framework:   fw,
					Action:      ctrl,
					Improvement: improvement,
					Rationale: fmt.Sprintf("Implementing %s addresses %d findings in the %s category, improving %s compliance posture.",
						ctrl, len(findings), cat, fw),
				})
			}
		}
	}

	unique := map[string]bool{}
	var deduped []ComplianceImpact
	for _, ci := range impacts {
		key := ci.Framework + "|" + ci.Action
		if unique[key] {
			continue
		}
		unique[key] = true
		deduped = append(deduped, ci)
	}

	sort.Slice(deduped, func(i, j int) bool {
		return impScore(deduped[i].Improvement) > impScore(deduped[j].Improvement)
	})

	return deduped
}

func impScore(s string) int {
	switch s {
	case "Significant":
		return 3
	case "Moderate":
		return 2
	case "Minor":
		return 1
	default:
		return 0
	}
}

// ── PHASE 13 — ENTERPRISE REMEDIATION ROADMAP ──

func (e *SDIEngine) generateRemediationRoadmap(recs []DecisionRecommendation) SDIRemediationRoadmap {
	var roadmap SDIRemediationRoadmap

	var criticalRecs, highRecs, mediumRecs, lowRecs []DecisionRecommendation
	for _, r := range recs {
		switch r.Priority {
		case "Critical":
			criticalRecs = append(criticalRecs, r)
		case "High":
			highRecs = append(highRecs, r)
		case "Medium":
			mediumRecs = append(mediumRecs, r)
		default:
			lowRecs = append(lowRecs, r)
		}
	}

	// Phase 30: Quick wins — Low effort Critical/High priority
	for _, r := range criticalRecs {
		if r.Effort == "Low" {
			roadmap.Phase30 = append(roadmap.Phase30, SDIRoadmapItem{
				Action: r.Title, Category: categoryFromTitle(r.Title),
				Priority: r.Priority, Effort: r.Effort, RiskReduction: r.RiskReduction,
			})
		}
	}
	for _, r := range highRecs {
		if r.Effort == "Low" && len(roadmap.Phase30) < 5 {
			roadmap.Phase30 = append(roadmap.Phase30, SDIRoadmapItem{
				Action: r.Title, Category: categoryFromTitle(r.Title),
				Priority: r.Priority, Effort: r.Effort, RiskReduction: r.RiskReduction,
			})
		}
	}

	// Phase 90: Medium effort Critical/High + remaining Low effort
	for _, r := range criticalRecs {
		if r.Effort == "Medium" {
			roadmap.Phase90 = append(roadmap.Phase90, SDIRoadmapItem{
				Action: r.Title, Category: categoryFromTitle(r.Title),
				Priority: r.Priority, Effort: r.Effort, RiskReduction: r.RiskReduction,
			})
		}
	}
	for _, r := range highRecs {
		if r.Effort == "Medium" || (r.Effort == "Low" && !containsRoadmapItem(roadmap.Phase30, r.Title)) {
			roadmap.Phase90 = append(roadmap.Phase90, SDIRoadmapItem{
				Action: r.Title, Category: categoryFromTitle(r.Title),
				Priority: r.Priority, Effort: r.Effort, RiskReduction: r.RiskReduction,
			})
		}
	}

	// Phase 180: High effort Critical/High + remaining Medium priority
	for _, r := range criticalRecs {
		if r.Effort == "High" {
			roadmap.Phase180 = append(roadmap.Phase180, SDIRoadmapItem{
				Action: r.Title, Category: categoryFromTitle(r.Title),
				Priority: r.Priority, Effort: r.Effort, RiskReduction: r.RiskReduction,
			})
		}
	}
	for _, r := range highRecs {
		if r.Effort == "High" {
			roadmap.Phase180 = append(roadmap.Phase180, SDIRoadmapItem{
				Action: r.Title, Category: categoryFromTitle(r.Title),
				Priority: r.Priority, Effort: r.Effort, RiskReduction: r.RiskReduction,
			})
		}
	}
	for _, r := range mediumRecs {
		if len(roadmap.Phase180) < 5 {
			roadmap.Phase180 = append(roadmap.Phase180, SDIRoadmapItem{
				Action: r.Title, Category: categoryFromTitle(r.Title),
				Priority: r.Priority, Effort: r.Effort, RiskReduction: r.RiskReduction,
			})
		}
	}

	// Phase 12m: Remaining medium + low priority items
	for _, r := range mediumRecs {
		if !containsRoadmapItem(roadmap.Phase90, r.Title) && !containsRoadmapItem(roadmap.Phase180, r.Title) {
			roadmap.Phase12m = append(roadmap.Phase12m, SDIRoadmapItem{
				Action: r.Title, Category: categoryFromTitle(r.Title),
				Priority: r.Priority, Effort: r.Effort, RiskReduction: r.RiskReduction,
			})
		}
	}
	for _, r := range lowRecs {
		roadmap.Phase12m = append(roadmap.Phase12m, SDIRoadmapItem{
			Action: r.Title, Category: categoryFromTitle(r.Title),
			Priority: r.Priority, Effort: r.Effort, RiskReduction: r.RiskReduction,
		})
	}

	return roadmap
}

// ── PHASE 14 — DECISION DASHBOARD ──

func (e *SDIEngine) generateDashboard(recs []DecisionRecommendation) DecisionDashboard {
	var topDecisions, quickWins, strategicActions []DecisionRecommendation

	for _, r := range recs {
		if r.Priority == "Critical" || r.Priority == "High" {
			topDecisions = append(topDecisions, r)
		}
		if r.Effort == "Low" && (r.Priority == "Critical" || r.Priority == "High") {
			quickWins = append(quickWins, r)
		}
		if r.Effort == "High" && r.Priority == "Critical" {
			strategicActions = append(strategicActions, r)
		}
	}

	if len(topDecisions) > 5 {
		topDecisions = topDecisions[:5]
	}
	if len(quickWins) > 5 {
		quickWins = quickWins[:5]
	}
	if len(strategicActions) > 5 {
		strategicActions = strategicActions[:5]
	}

	var totalReduction float64
	for _, r := range recs {
		totalReduction += reductionWeight(r.RiskReduction)
	}

	summary := fmt.Sprintf("Top %d decisions from %d recommendations: %d quick wins available with immediate risk reduction.",
		len(topDecisions), len(recs), len(quickWins))

	return DecisionDashboard{
		TopDecisions:         topDecisions,
		QuickWins:            quickWins,
		StrategicActions:     strategicActions,
		RiskReductionSummary: summary,
		TotalRiskReduction:   totalReduction,
	}
}

// ── PHASE 15 — EXECUTIVE SCENARIOS ──

func (e *SDIEngine) generateExecutiveScenarios(input SDIInput, recs []DecisionRecommendation) ExecutiveScenarios {
	criticalCount := countSeverity(input.Findings, "Critical")
	totalFindings := len(input.Findings)
	totalPaths := len(input.AttackPaths)
	currentCoverage := computeCoverage(input.Findings, input.Controls)

	bestResolved := 0
	bestPaths := 0
	likelyResolved := 0
	likelyPaths := 0

	for _, r := range recs {
		if r.Priority == "Critical" || r.Priority == "High" {
			bestResolved += len(r.AffectedFindings)
			bestPaths += len(r.AffectedAttackPaths)
		}
		if r.Priority == "Critical" {
			likelyResolved += len(r.AffectedFindings)
			likelyPaths += len(r.AffectedAttackPaths)
		}
	}

	bestResolved = min(bestResolved, totalFindings)
	likelyResolved = min(likelyResolved, totalFindings)
	bestPaths = min(bestPaths, totalPaths)
	likelyPaths = min(likelyPaths, totalPaths)

	return ExecutiveScenarios{
		BestCase: ExecutiveScenario{
			Scenario:          "Best Case",
			RiskScore:         input.RiskScore - float64(bestResolved)*0.15,
			FindingsResolved:  bestResolved,
			AttackPathsClosed: bestPaths,
			CoverageAchieved:  min(currentCoverage+float64(bestResolved)*2.0, 100),
			Description:       fmt.Sprintf("All high-priority recommendations implemented. %d findings resolved, %d attack paths closed.", bestResolved, bestPaths),
		},
		LikelyCase: ExecutiveScenario{
			Scenario:          "Likely Case",
			RiskScore:         input.RiskScore - float64(likelyResolved)*0.1,
			FindingsResolved:  likelyResolved,
			AttackPathsClosed: likelyPaths,
			CoverageAchieved:  min(currentCoverage+float64(likelyResolved)*1.5, 100),
			Description:       fmt.Sprintf("Critical recommendations implemented. %d findings resolved, %d attack paths closed.", likelyResolved, likelyPaths),
		},
		WorstCase: ExecutiveScenario{
			Scenario:          "Worst Case",
			RiskScore:         input.RiskScore + float64(criticalCount)*0.2,
			FindingsResolved:  0,
			AttackPathsClosed: 0,
			CoverageAchieved:  currentCoverage,
			Description:       fmt.Sprintf("No action taken. %d Critical findings remain unaddressed with %d active attack paths.", criticalCount, totalPaths),
		},
	}
}

// ── PHASE 12-16 helpers ──

func containsRoadmapItem(items []SDIRoadmapItem, action string) bool {
	for _, item := range items {
		if item.Action == action {
			return true
		}
	}
	return false
}

func filterPriority(recs []DecisionRecommendation, priority string) []DecisionRecommendation {
	var out []DecisionRecommendation
	for _, r := range recs {
		if r.Priority == priority {
			out = append(out, r)
		}
	}
	return out
}

func countSeverity(findings []SDRIFinding, severity string) int {
	count := 0
	for _, f := range findings {
		if strings.EqualFold(f.Severity, severity) {
			count++
		}
	}
	return count
}

func computeCoverage(findings []SDRIFinding, controls []SDRIControl) float64 {
	if len(controls) == 0 && len(findings) > 0 {
		return 0
	}
	if len(findings) == 0 {
		return 100
	}
	covered := 0
	findingCategories := map[string]bool{}
	for _, f := range findings {
		findingCategories[f.Category] = true
	}
	for _, c := range controls {
		cat := c.Category
		if mapped, ok := controlCategoryMap[c.Name]; ok {
			cat = mapped
		}
		if findingCategories[cat] {
			covered++
		}
	}
	if len(findingCategories) == 0 {
		return 0
	}
	return float64(covered) / float64(len(findingCategories)) * 100.0
}

func categoryMatch(findingCat, controlCat string) bool {
	pairs := map[string]string{
		"AccessControl":           "AccessControl",
		"Authentication":          "AccessControl",
		"Authorization":           "Authorization",
		"DataProtection":          "DataProtection",
		"Logging":                 "Logging",
		"Monitoring":              "Monitoring",
		"NetworkSecurity":         "NetworkSecurity",
		"SecretsManagement":       "SecretsManagement",
		"KeyManagement":           "SecretsManagement",
		"IdentityGovernance":      "IdentityGovernance",
		"VulnerabilityManagement": "VulnerabilityManagement",
		"EndpointSecurity":        "EndpointSecurity",
		"Governance":              "Governance",
		"Resilience":              "Resilience",
		"ThirdPartyRisk":          "ThirdPartyRisk",
		"ApplicationSecurity":     "ApplicationSecurity",
	}
	fc := strings.ToLower(findingCat)
	cc := strings.ToLower(controlCat)
	for fk, ck := range pairs {
		if strings.Contains(fc, strings.ToLower(fk)) &&
			strings.Contains(cc, strings.ToLower(ck)) {
			return true
		}
	}
	return fc == cc
}

func extractControlName(title string) string {
	title = strings.TrimPrefix(title, "Enable ")
	title = strings.TrimPrefix(title, "Implement ")
	title = strings.TrimPrefix(title, "Deploy ")
	title = strings.TrimPrefix(title, "Enforce ")
	title = strings.TrimPrefix(title, "Establish ")
	title = strings.TrimPrefix(title, "Conduct ")
	return title
}

func controlKeywords(title string) []string {
	lower := strings.ToLower(title)
	words := strings.Fields(lower)
	var keywords []string
	skipWords := map[string]bool{"a": true, "an": true, "the": true, "and": true, "or": true,
		"for": true, "with": true, "all": true, "on": true, "in": true, "to": true, "of": true}
	for _, w := range words {
		if !skipWords[w] && len(w) > 2 {
			keywords = append(keywords, w)
		}
	}
	return keywords
}

func categoryFromTitle(title string) string {
	lower := strings.ToLower(title)
	switch {
	case containsAny(lower, []string{"mfa", "auth", "identity", "sso", "iam", "login", "password", "access"}):
		return "Identity"
	case containsAny(lower, []string{"encrypt", "tls", "data", "key", "secret", "kms"}):
		return "Data Protection"
	case containsAny(lower, []string{"siem", "log", "monitor", "edr", "detect"}):
		return "Monitoring"
	case containsAny(lower, []string{"supply", "third", "vendor", "partner"}):
		return "Third Party Risk"
	case containsAny(lower, []string{"cloud", "container", "k8s", "kubernetes"}):
		return "Cloud Security"
	case containsAny(lower, []string{"network", "segment", "firewall", "ids", "ips", "vpn"}):
		return "Network Security"
	case containsAny(lower, []string{"endpoint", "edr", "antivirus"}):
		return "Endpoint Security"
	case containsAny(lower, []string{"patch", "vulnerability", "application", "input"}):
		return "Application Security"
	case containsAny(lower, []string{"incident", "response", "governance", "review", "policy"}):
		return "Governance"
	case containsAny(lower, []string{"backup", "disaster", "recovery", "resilience"}):
		return "Resilience"
	default:
		return "General"
	}
}
