package intelligence

import (
	"fmt"
	"sort"
	"strings"
)

// ─────────────────────────────────────────────────────────────
// CIARE — Compliance Intelligence & Audit Readiness Engine
// ASF V7 — Phases 1-15
// ─────────────────────────────────────────────────────────────

type CIAREEngine struct{}

func NewCIAREEngine() *CIAREEngine {
	return &CIAREEngine{}
}

// ─────────────────────────────────────────────────────────────
// PHASE 1 — COMPLIANCE FRAMEWORK ENGINE
// ─────────────────────────────────────────────────────────────

// ComplianceFramework represents a supported compliance framework.
type ComplianceFramework string

const (
	FrameworkHIPAA     ComplianceFramework = "HIPAA"
	FrameworkSOC2      ComplianceFramework = "SOC2"
	FrameworkISO27001  ComplianceFramework = "ISO27001"
	FrameworkPCIDSS    ComplianceFramework = "PCI-DSS"
	FrameworkNIST80053 ComplianceFramework = "NIST800-53"
	FrameworkCIS       ComplianceFramework = "CIS"
	FrameworkGDPR      ComplianceFramework = "GDPR"
)

var allFrameworks = []ComplianceFramework{
	FrameworkHIPAA,
	FrameworkSOC2,
	FrameworkISO27001,
	FrameworkPCIDSS,
	FrameworkNIST80053,
	FrameworkCIS,
	FrameworkGDPR,
}

// ─────────────────────────────────────────────────────────────
// CIARE Input / Output
// ─────────────────────────────────────────────────────────────

type CIAREInput struct {
	Architecture *ArchDescription
	SDRIResult   *SDRIResult
	Domain       string
	Compliance   []string
}

type CIAREResult struct {
	FrameworkCoverages   []CIAREFrameworkCoverage   `json:"framework_coverages"`
	AuditReadinessScores []CIAREAuditReadiness      `json:"audit_readiness_scores"`
	EvidenceRequirements []CIAREEvidenceRequirement `json:"evidence_requirements"`
	MissingEvidences     []CIAREMissingEvidence     `json:"missing_evidences"`
	AuditorQuestions     []CIAREAuditorQuestion     `json:"auditor_questions"`
	ComplianceGaps       []CIAREComplianceGap       `json:"compliance_gaps"`
	ControlMaturities    []CIAREControlMaturity     `json:"control_maturities"`
	ComplianceNarratives []CIAREComplianceNarrative `json:"compliance_narratives"`
	AuditPackage         CIAREAuditPackage          `json:"audit_package"`
	ComplianceDashboard  CIAREComplianceDashboard   `json:"compliance_dashboard"`
	ProcurementQuestions []CIAREProcurementQuestion `json:"procurement_questions"`
}

// Phase 3 — Framework Coverage
type CIAREFrameworkCoverage struct {
	Framework        string   `json:"framework"`
	Required         int      `json:"required"`
	Observed         int      `json:"observed"`
	Missing          int      `json:"missing"`
	CoveragePct      float64  `json:"coverage_pct"`
	Status           string   `json:"status"`
	ObservedControls []string `json:"observed_controls,omitempty"`
	MissingControls  []string `json:"missing_controls,omitempty"`
}

// Phase 4 — Audit Readiness
type CIAREAuditReadiness struct {
	Framework       string   `json:"framework"`
	ReadinessScore  float64  `json:"readiness_score"`
	Status          string   `json:"status"`
	ControlCoverage float64  `json:"control_coverage"`
	EvidenceScore   float64  `json:"evidence_score"`
	ThreatExposure  float64  `json:"threat_exposure"`
	FindingsPenalty float64  `json:"findings_penalty"`
	Factors         []string `json:"factors"`
}

// Phase 5 — Evidence Requirement
type CIAREEvidenceRequirement struct {
	Framework string   `json:"framework"`
	Control   string   `json:"control"`
	Evidence  []string `json:"evidence"`
}

// Phase 6 — Missing Evidence
type CIAREMissingEvidence struct {
	Framework string   `json:"framework"`
	Control   string   `json:"control"`
	Evidences []string `json:"evidences"`
}

// Phase 7 — Auditor Question
type CIAREAuditorQuestion struct {
	Framework string `json:"framework"`
	Control   string `json:"control"`
	Question  string `json:"question"`
}

// Phase 8 — Compliance Gap
type CIAREComplianceGap struct {
	ID          string `json:"id"`
	Framework   string `json:"framework"`
	Requirement string `json:"requirement"`
	Observed    string `json:"observed"`
	Missing     string `json:"missing"`
	Risk        string `json:"risk"`
}

// Phase 9 — Control Maturity
type CIAREControlMaturity struct {
	Domain   string  `json:"domain"`
	Level    int     `json:"level"`
	Label    string  `json:"label"`
	Coverage float64 `json:"coverage"`
}

// Phase 11 — Compliance Narrative
type CIAREComplianceNarrative struct {
	Framework string `json:"framework"`
	Narrative string `json:"narrative"`
}

// Phase 12 — Audit Package
type CIAREAuditPackage struct {
	ExecutiveSummary     string                     `json:"executive_summary"`
	FrameworkCoverages   []CIAREFrameworkCoverage   `json:"framework_coverages"`
	ControlInventory     []SDRIControl              `json:"control_inventory"`
	MissingControls      []CIAREComplianceGap       `json:"missing_controls"`
	EvidenceRequirements []CIAREEvidenceRequirement `json:"evidence_requirements"`
	AuditorQuestions     []CIAREAuditorQuestion     `json:"auditor_questions"`
}

// Phase 13 — Compliance Dashboard
type CIAREComplianceDashboard struct {
	FrameworkCoverages map[string]float64     `json:"framework_coverages"`
	TopGaps            []CIAREComplianceGap   `json:"top_gaps"`
	TopMissingEvidence []CIAREMissingEvidence `json:"top_missing_evidence"`
	TopRisks           []string               `json:"top_risks"`
}

// Phase 14 — Procurement Question
type CIAREProcurementQuestion struct {
	Category string `json:"category"`
	Question string `json:"question"`
}

// ─────────────────────────────────────────────────────────────
// Internal types
// ─────────────────────────────────────────────────────────────

type fwCoverage struct {
	required        int
	observed        int
	missing         int
	pct             float64
	presentControls []string
	missingControls []string
}

// ─────────────────────────────────────────────────────────────
// PHASE 2 — CONTROL TO FRAMEWORK MAPPING
// ─────────────────────────────────────────────────────────────

func buildControlFrameworkMappings() map[string][]ComplianceFramework {
	return map[string][]ComplianceFramework{
		"MFA":                        {FrameworkHIPAA, FrameworkSOC2, FrameworkISO27001, FrameworkNIST80053, FrameworkPCIDSS},
		"PasswordPolicy":             {FrameworkHIPAA, FrameworkSOC2, FrameworkISO27001},
		"SessionManagement":          {FrameworkSOC2, FrameworkISO27001},
		"ConditionalAccess":          {FrameworkSOC2, FrameworkNIST80053},
		"RBAC":                       {FrameworkHIPAA, FrameworkSOC2, FrameworkISO27001, FrameworkNIST80053, FrameworkPCIDSS, FrameworkGDPR},
		"ABAC":                       {FrameworkNIST80053},
		"JustInTimeAccess":           {FrameworkSOC2, FrameworkNIST80053},
		"PrivilegedAccessManagement": {FrameworkHIPAA, FrameworkSOC2, FrameworkISO27001, FrameworkNIST80053, FrameworkPCIDSS},
		"IdentityGovernance":         {FrameworkSOC2, FrameworkISO27001},
		"AccessReviews":              {FrameworkHIPAA, FrameworkSOC2, FrameworkISO27001, FrameworkPCIDSS, FrameworkNIST80053, FrameworkGDPR},
		"SecretsRotation":            {FrameworkSOC2, FrameworkISO27001, FrameworkPCIDSS},
		"SecretsVault":               {FrameworkSOC2, FrameworkISO27001, FrameworkPCIDSS},
		"SecretsScanning":            {FrameworkSOC2, FrameworkISO27001},
		"KeyRotation":                {FrameworkHIPAA, FrameworkSOC2, FrameworkISO27001, FrameworkPCIDSS, FrameworkNIST80053},
		"KeyAccessLogging":           {FrameworkHIPAA, FrameworkSOC2, FrameworkISO27001, FrameworkPCIDSS, FrameworkNIST80053},
		"SeparationOfDuties":         {FrameworkSOC2, FrameworkISO27001, FrameworkPCIDSS},
		"NetworkSegmentation":        {FrameworkPCIDSS, FrameworkNIST80053, FrameworkCIS},
		"FirewallRules":              {FrameworkPCIDSS, FrameworkNIST80053, FrameworkCIS},
		"TLSEncryption":              {FrameworkHIPAA, FrameworkSOC2, FrameworkISO27001, FrameworkPCIDSS, FrameworkGDPR},
		"VPNAccess":                  {FrameworkHIPAA, FrameworkNIST80053},
		"IntrusionDetection":         {FrameworkPCIDSS, FrameworkNIST80053, FrameworkCIS},
		"AuditLogging":               {FrameworkHIPAA, FrameworkSOC2, FrameworkISO27001, FrameworkPCIDSS, FrameworkNIST80053, FrameworkCIS, FrameworkGDPR},
		"SIEMIntegration":            {FrameworkHIPAA, FrameworkSOC2, FrameworkISO27001, FrameworkPCIDSS, FrameworkNIST80053, FrameworkCIS},
		"RealTimeAlerting":           {FrameworkSOC2, FrameworkPCIDSS, FrameworkNIST80053},
		"IncidentResponsePlan":       {FrameworkHIPAA, FrameworkSOC2, FrameworkISO27001, FrameworkPCIDSS, FrameworkNIST80053, FrameworkGDPR},
		"IncidentResponseTeam":       {FrameworkHIPAA, FrameworkSOC2, FrameworkISO27001},
		"AutomatedBackup":            {FrameworkHIPAA, FrameworkSOC2, FrameworkISO27001, FrameworkPCIDSS, FrameworkNIST80053, FrameworkGDPR, FrameworkCIS},
		"DisasterRecoveryPlan":       {FrameworkSOC2, FrameworkISO27001, FrameworkNIST80053},
		"DataEncryptionAtRest":       {FrameworkHIPAA, FrameworkSOC2, FrameworkISO27001, FrameworkPCIDSS, FrameworkGDPR, FrameworkNIST80053, FrameworkCIS},
		"DataEncryptionInTransit":    {FrameworkHIPAA, FrameworkSOC2, FrameworkISO27001, FrameworkPCIDSS, FrameworkGDPR, FrameworkNIST80053, FrameworkCIS},
		"DataClassification":         {FrameworkHIPAA, FrameworkGDPR, FrameworkSOC2},
		"DLPControls":                {FrameworkHIPAA, FrameworkPCIDSS, FrameworkGDPR},
		"DataRetentionPolicy":        {FrameworkHIPAA, FrameworkGDPR, FrameworkSOC2, FrameworkPCIDSS},
		"PrivacyImpactAssessment":    {FrameworkGDPR, FrameworkHIPAA},
		"ConsentManagement":          {FrameworkGDPR},
		"DataSubjectAccess":          {FrameworkGDPR},
		"ThirdPartyDueDiligence":     {FrameworkSOC2, FrameworkISO27001, FrameworkPCIDSS},
		"VendorRiskManagement":       {FrameworkSOC2, FrameworkISO27001},
		"CloudSecurityPosture":       {FrameworkSOC2, FrameworkNIST80053, FrameworkCIS},
		"CloudAccessBroker":          {FrameworkSOC2, FrameworkNIST80053},
		"ContainerImageScanning":     {FrameworkSOC2, FrameworkNIST80053, FrameworkCIS},
		"ContainerRuntimeSecurity":   {FrameworkSOC2, FrameworkNIST80053, FrameworkCIS},
		"K8sRBAC":                    {FrameworkSOC2, FrameworkNIST80053, FrameworkCIS},
		"K8sAdmissionControllers":    {FrameworkNIST80053, FrameworkCIS},
		"K8sNetworkPolicies":         {FrameworkNIST80053, FrameworkCIS},
		"K8sPodSecurity":             {FrameworkNIST80053, FrameworkCIS},
		"EndpointDetectionResponse":  {FrameworkNIST80053, FrameworkCIS},
		"Antivirus":                  {FrameworkPCIDSS, FrameworkNIST80053, FrameworkCIS},
		"PatchManagement":            {FrameworkHIPAA, FrameworkPCIDSS, FrameworkNIST80053, FrameworkCIS, FrameworkISO27001},
		"ChangeManagement":           {FrameworkSOC2, FrameworkISO27001, FrameworkNIST80053},
		"VulnerabilityScanning":      {FrameworkHIPAA, FrameworkPCIDSS, FrameworkNIST80053, FrameworkCIS, FrameworkISO27001},
		"SecurityAwarenessTraining":  {FrameworkHIPAA, FrameworkSOC2, FrameworkISO27001, FrameworkPCIDSS},
		"AuditTrail":                 {FrameworkHIPAA, FrameworkSOC2, FrameworkISO27001, FrameworkPCIDSS, FrameworkNIST80053, FrameworkCIS, FrameworkGDPR},
		"BreakGlassAccess":           {FrameworkHIPAA, FrameworkSOC2, FrameworkISO27001},
		"FraudMonitoring":            {FrameworkPCIDSS, FrameworkSOC2},
		"SettlementIntegrity":        {FrameworkPCIDSS, FrameworkSOC2, FrameworkISO27001},
		"DualControl":                {FrameworkPCIDSS, FrameworkSOC2, FrameworkISO27001},
		"TransactionLogging":         {FrameworkPCIDSS, FrameworkSOC2},
		"BackupRecoveryTesting":      {FrameworkHIPAA, FrameworkSOC2, FrameworkISO27001, FrameworkPCIDSS, FrameworkGDPR},
		"ConfigurationHardening":     {FrameworkCIS, FrameworkNIST80053, FrameworkPCIDSS},
		"PhishingAwareness":          {FrameworkSOC2, FrameworkISO27001},
		"IncidentResponseAutomation": {FrameworkSOC2, FrameworkNIST80053},
		"ThreatIntelligence":         {FrameworkNIST80053, FrameworkCIS},
		"SupplyChainSecurity":        {FrameworkSOC2, FrameworkISO27001, FrameworkNIST80053},
	}
}

// ─────────────────────────────────────────────────────────────
// FRAMEWORK REQUIRED CONTROLS
// ─────────────────────────────────────────────────────────────

var frameworkRequiredControls = map[ComplianceFramework][]string{
	FrameworkHIPAA: {
		"MFA", "PasswordPolicy", "RBAC", "PrivilegedAccessManagement",
		"AccessReviews", "AuditLogging", "SIEMIntegration", "IncidentResponsePlan",
		"AutomatedBackup", "DataEncryptionAtRest", "DataEncryptionInTransit",
		"DataClassification", "DLPControls", "DataRetentionPolicy",
		"KeyRotation", "KeyAccessLogging", "VulnerabilityScanning",
		"SecurityAwarenessTraining", "AuditTrail", "BreakGlassAccess",
		"BackupRecoveryTesting", "PatchManagement",
	},
	FrameworkSOC2: {
		"MFA", "PasswordPolicy", "SessionManagement", "ConditionalAccess",
		"RBAC", "JustInTimeAccess", "PrivilegedAccessManagement",
		"IdentityGovernance", "AccessReviews", "AuditLogging",
		"SIEMIntegration", "RealTimeAlerting", "IncidentResponsePlan",
		"IncidentResponseTeam", "AutomatedBackup", "DisasterRecoveryPlan",
		"DataEncryptionAtRest", "DataEncryptionInTransit", "DataClassification",
		"ThirdPartyDueDiligence", "VendorRiskManagement", "CloudSecurityPosture",
		"ChangeManagement", "SecurityAwarenessTraining", "AuditTrail",
		"BackupRecoveryTesting", "PhishingAwareness", "IncidentResponseAutomation",
	},
	FrameworkISO27001: {
		"MFA", "PasswordPolicy", "SessionManagement", "RBAC",
		"IdentityGovernance", "AccessReviews", "SecretsRotation",
		"SecretsVault", "KeyRotation", "KeyAccessLogging",
		"AuditLogging", "SIEMIntegration", "IncidentResponsePlan",
		"IncidentResponseTeam", "AutomatedBackup", "DisasterRecoveryPlan",
		"DataEncryptionAtRest", "DataEncryptionInTransit",
		"ThirdPartyDueDiligence", "VendorRiskManagement",
		"PatchManagement", "VulnerabilityScanning",
		"SecurityAwarenessTraining", "AuditTrail",
		"BackupRecoveryTesting", "SeparationOfDuties",
	},
	FrameworkPCIDSS: {
		"MFA", "RBAC", "PrivilegedAccessManagement", "AccessReviews",
		"NetworkSegmentation", "FirewallRules", "TLSEncryption",
		"IntrusionDetection", "AuditLogging", "SIEMIntegration",
		"RealTimeAlerting", "IncidentResponsePlan", "AutomatedBackup",
		"DataEncryptionAtRest", "DataEncryptionInTransit", "DLPControls",
		"Antivirus", "PatchManagement", "VulnerabilityScanning",
		"SecurityAwarenessTraining", "AuditTrail", "ConfigurationHardening",
		"FraudMonitoring", "DualControl", "TransactionLogging",
		"BackupRecoveryTesting", "ThirdPartyDueDiligence",
	},
	FrameworkNIST80053: {
		"MFA", "RBAC", "ConditionalAccess", "AccessReviews",
		"KeyRotation", "KeyAccessLogging", "NetworkSegmentation",
		"FirewallRules", "VPNAccess", "IntrusionDetection",
		"AuditLogging", "SIEMIntegration", "RealTimeAlerting",
		"IncidentResponsePlan", "AutomatedBackup", "DisasterRecoveryPlan",
		"DataEncryptionAtRest", "DataEncryptionInTransit",
		"CloudSecurityPosture", "PatchManagement", "VulnerabilityScanning",
		"AuditTrail", "EndpointDetectionResponse", "ConfigurationHardening",
		"JustInTimeAccess", "ChangeManagement", "SupplyChainSecurity",
	},
	FrameworkCIS: {
		"NetworkSegmentation", "FirewallRules", "IntrusionDetection",
		"AuditLogging", "SIEMIntegration", "ContainerImageScanning",
		"ContainerRuntimeSecurity", "K8sRBAC", "K8sAdmissionControllers",
		"K8sNetworkPolicies", "K8sPodSecurity", "EndpointDetectionResponse",
		"Antivirus", "PatchManagement", "VulnerabilityScanning",
		"ConfigurationHardening", "AuditTrail", "AutomatedBackup",
		"DataEncryptionAtRest", "DataEncryptionInTransit",
		"CloudSecurityPosture", "ThreatIntelligence",
	},
	FrameworkGDPR: {
		"TLSEncryption", "DataEncryptionAtRest", "DataEncryptionInTransit",
		"DataClassification", "DLPControls", "DataRetentionPolicy",
		"PrivacyImpactAssessment", "ConsentManagement", "DataSubjectAccess",
		"AuditLogging", "AuditTrail", "RBAC",
		"AccessReviews", "AutomatedBackup", "IncidentResponsePlan",
		"BackupRecoveryTesting",
	},
}

// ─────────────────────────────────────────────────────────────
// PHASE 5 — EVIDENCE REQUIREMENTS
// ─────────────────────────────────────────────────────────────

var controlEvidence = map[string][]string{
	"MFA":                        {"MFA policy document", "MFA configuration export", "MFA enrollment records"},
	"PasswordPolicy":             {"Password policy document", "Password complexity configuration", "Password rotation reports"},
	"SessionManagement":          {"Session timeout configuration", "Session management policy", "Session audit logs"},
	"ConditionalAccess":          {"Conditional access policy", "Risk-based access rules", "Access evaluation logs"},
	"RBAC":                       {"Role definitions document", "Access control policy", "User permission audits"},
	"ABAC":                       {"Attribute policy definitions", "Access control rules", "Attribute source configuration"},
	"JustInTimeAccess":           {"JIT access policy", "Privileged access requests", "JIT approval logs"},
	"PrivilegedAccessManagement": {"PAM policy document", "Privileged account inventory", "PAM session recordings"},
	"IdentityGovernance":         {"Identity lifecycle policy", "User provisioning reports", "Identity certification records"},
	"AccessReviews":              {"Access review schedule", "Completed review reports", "Remediation tracking"},
	"SecretsRotation":            {"Secrets rotation policy", "Rotation automation config", "Rotation audit logs"},
	"SecretsVault":               {"Vault configuration", "Access policies", "Vault audit logs"},
	"SecretsScanning":            {"Scanning configuration", "Scan reports", "Remediation records"},
	"KeyRotation":                {"Key rotation policy", "Rotation schedule", "Key rotation audit logs"},
	"KeyAccessLogging":           {"Key access logs", "Key usage monitoring", "Access alert configuration"},
	"SeparationOfDuties":         {"SOD policy document", "Role segregation matrix", "SOD violation reports"},
	"NetworkSegmentation":        {"Network topology diagram", "Segment access rules", "Firewall rule review"},
	"FirewallRules":              {"Firewall rule base", "Rule review records", "Change approval logs"},
	"TLSEncryption":              {"TLS configuration", "Certificate inventory", "Cipher suite documentation"},
	"VPNAccess":                  {"VPN configuration", "Remote access policy", "VPN connection logs"},
	"IntrusionDetection":         {"IDS/IPS configuration", "Alert records", "Signatures update log"},
	"AuditLogging":               {"Audit log configuration", "Log retention policy", "Log review procedures"},
	"SIEMIntegration":            {"SIEM architecture diagram", "Log source inventory", "Alert rule definitions"},
	"RealTimeAlerting":           {"Alerting policy", "Incident alert configuration", "Alert response records"},
	"IncidentResponsePlan":       {"IR plan document", "Incident classification guide", "IR test/tabletop records"},
	"IncidentResponseTeam":       {"Team roster", "Incident contact list", "Training records"},
	"AutomatedBackup":            {"Backup schedule", "Backup monitoring reports", "Restore test results"},
	"DisasterRecoveryPlan":       {"DR plan document", "DR test results", "RTO/RPO documentation"},
	"DataEncryptionAtRest":       {"Encryption policy", "Key management documentation", "Encryption configuration"},
	"DataEncryptionInTransit":    {"TLS/SSL configuration", "Certificate management policy", "Network encryption scan"},
	"DataClassification":         {"Classification policy", "Data sensitivity labels", "Handling procedures"},
	"DLPControls":                {"DLP policy", "DLP rule configuration", "DLP incident reports"},
	"DataRetentionPolicy":        {"Retention schedule", "Disposal procedures", "Retention audit records"},
	"PrivacyImpactAssessment":    {"PIA policy", "Completed PIA reports", "PIA review records"},
	"ConsentManagement":          {"Consent collection records", "Consent management configuration", "Consent audit trail"},
	"DataSubjectAccess":          {"DSAR procedure", "DSAR request logs", "DSAR response records"},
	"ThirdPartyDueDiligence":     {"Vendor assessment policy", "Completed assessments", "Vendor risk ratings"},
	"VendorRiskManagement":       {"Vendor inventory", "Risk monitoring reports", "Vendor review schedule"},
	"CloudSecurityPosture":       {"CSPM configuration", "Cloud security benchmarks", "Compliance scan reports"},
	"CloudAccessBroker":          {"CASB configuration", "Access policy rules", "Shadow IT discovery reports"},
	"ContainerImageScanning":     {"Image scanning policy", "Vulnerability scan reports", "Image approval process"},
	"ContainerRuntimeSecurity":   {"Runtime security policy", "Container monitoring config", "Security event logs"},
	"K8sRBAC":                    {"K8s RBAC configuration", "Service account inventory", "RBAC audit reports"},
	"K8sAdmissionControllers":    {"Admission controller config", "Policy definitions", "Admission audit logs"},
	"K8sNetworkPolicies":         {"Network policy definitions", "Network policy audit", "Traffic flow documentation"},
	"K8sPodSecurity":             {"Pod security standards", "PSP/PSS configuration", "Pod security audit logs"},
	"EndpointDetectionResponse":  {"EDR configuration", "Endpoint inventory", "Threat detection records"},
	"Antivirus":                  {"AV configuration", "Scan schedule", "Malware detection records"},
	"PatchManagement":            {"Patch management policy", "Patch deployment records", "Vulnerability scan reports"},
	"ChangeManagement":           {"Change management policy", "Change request records", "Change approval audit"},
	"VulnerabilityScanning":      {"Scan policy", "Scan results", "Remediation tracking"},
	"SecurityAwarenessTraining":  {"Training policy", "Training completion records", "Phishing simulation results"},
	"AuditTrail":                 {"Audit trail configuration", "Audit log review schedule", "Audit trail retention"},
	"BreakGlassAccess":           {"Emergency access policy", "Break glass procedure", "Break glass usage logs"},
	"FraudMonitoring":            {"Fraud detection policy", "Monitoring rules", "Fraud investigation records"},
	"SettlementIntegrity":        {"Settlement reconciliation policy", "Reconciliation reports", "Exception handling logs"},
	"DualControl":                {"Dual control policy", "Transaction approval logs", "Dual control audit records"},
	"TransactionLogging":         {"Transaction log configuration", "Transaction audit trail", "Log monitoring"},
	"BackupRecoveryTesting":      {"Restore test schedule", "Test results", "Recvery validation reports"},
	"ConfigurationHardening":     {"Hardening standards", "Benchmark compliance reports", "Configuration drift monitoring"},
	"PhishingAwareness":          {"Phishing simulation policy", "Simulation results", "Training improvement tracking"},
	"IncidentResponseAutomation": {"SOAR configuration", "Playbook definitions", "Automation metrics"},
	"ThreatIntelligence":         {"Threat intel feed list", "Intel analysis reports", "Indicators of compromise"},
	"SupplyChainSecurity":        {"Supply chain policy", "SBOM records", "Vendor security assessments"},
}

// ─────────────────────────────────────────────────────────────
// PHASE 7 — AUDITOR QUESTIONS
// ─────────────────────────────────────────────────────────────

func generateAuditorQuestions(fw ComplianceFramework, controls []string, coverage *fwCoverage) []CIAREAuditorQuestion {
	qs := make([]CIAREAuditorQuestion, 0)
	if coverage == nil {
		return qs
	}
	for _, ctrl := range controls {
		q := auditorQuestionForControl(fw, ctrl)
		if q != "" {
			qs = append(qs, CIAREAuditorQuestion{
				Framework: string(fw),
				Control:   ctrl,
				Question:  q,
			})
		}
	}
	return qs
}

func auditorQuestionForControl(fw ComplianceFramework, control string) string {
	switch control {
	case "MFA":
		return "How is multi-factor authentication enforced for all user access?"
	case "PasswordPolicy":
		return "What are the password complexity and rotation requirements?"
	case "SessionManagement":
		return "How are sessions managed, timed out, and invalidated?"
	case "ConditionalAccess":
		return "What conditions are evaluated for access decisions?"
	case "RBAC":
		return "How are roles defined and access permissions assigned?"
	case "ABAC":
		return "What attributes are used for access control decisions?"
	case "JustInTimeAccess":
		return "How is just-in-time privileged access managed?"
	case "PrivilegedAccessManagement":
		return "How are privileged accounts managed and monitored?"
	case "IdentityGovernance":
		return "How is the identity lifecycle managed across systems?"
	case "AccessReviews":
		return "How often are access reviews performed and remediated?"
	case "SecretsRotation":
		return "How frequently are secrets and credentials rotated?"
	case "SecretsVault":
		return "How are secrets stored and accessed securely?"
	case "SecretsScanning":
		return "How is code scanned for hardcoded secrets?"
	case "KeyRotation":
		return "How are encryption keys rotated?"
	case "KeyAccessLogging":
		return "How is key access logged and monitored?"
	case "SeparationOfDuties":
		return "How is separation of duties enforced?"
	case "NetworkSegmentation":
		return "How is the network segmented and isolation enforced?"
	case "FirewallRules":
		return "How are firewall rules managed and reviewed?"
	case "TLSEncryption":
		return "How is TLS configured and certificates managed?"
	case "VPNAccess":
		return "How is remote access secured via VPN?"
	case "IntrusionDetection":
		return "How are network intrusions detected and alerted?"
	case "AuditLogging":
		return "How are security events logged and retained?"
	case "SIEMIntegration":
		return "How are logs aggregated and correlated in a SIEM?"
	case "RealTimeAlerting":
		return "How are real-time security alerts generated and handled?"
	case "IncidentResponsePlan":
		return "How is the incident response plan tested and maintained?"
	case "IncidentResponseTeam":
		return "Who is on the incident response team and how are they trained?"
	case "AutomatedBackup":
		return "How are automated backups configured and monitored?"
	case "DisasterRecoveryPlan":
		return "How is disaster recovery tested and RTO/RPO validated?"
	case "DataEncryptionAtRest":
		return "How is data encrypted at rest and keys managed?"
	case "DataEncryptionInTransit":
		return "How is data encrypted in transit between systems?"
	case "DataClassification":
		return "How is data classified by sensitivity level?"
	case "DLPControls":
		return "How is data loss prevention configured and monitored?"
	case "DataRetentionPolicy":
		return "How are data retention and disposal managed?"
	case "PrivacyImpactAssessment":
		return "How are privacy impact assessments conducted?"
	case "ConsentManagement":
		return "How is user consent collected and managed?"
	case "DataSubjectAccess":
		return "How are data subject access requests handled?"
	case "ThirdPartyDueDiligence":
		return "How are third-party vendors assessed for security?"
	case "VendorRiskManagement":
		return "How is vendor risk monitored on an ongoing basis?"
	case "CloudSecurityPosture":
		return "How is cloud security posture continuously assessed?"
	case "CloudAccessBroker":
		return "How is cloud access governed and shadow IT detected?"
	case "ContainerImageScanning":
		return "How are container images scanned for vulnerabilities?"
	case "ContainerRuntimeSecurity":
		return "How is container runtime security monitored?"
	case "K8sRBAC":
		return "How is Kubernetes RBAC configured and audited?"
	case "K8sAdmissionControllers":
		return "What admission controllers are enforced in Kubernetes?"
	case "K8sNetworkPolicies":
		return "How are Kubernetes network policies defined and enforced?"
	case "K8sPodSecurity":
		return "How are Kubernetes pod security standards enforced?"
	case "EndpointDetectionResponse":
		return "How are endpoints monitored for threats and responded to?"
	case "Antivirus":
		return "How is antivirus deployed and kept up to date?"
	case "PatchManagement":
		return "How are security patches deployed and tracked?"
	case "ChangeManagement":
		return "How are changes reviewed and approved before deployment?"
	case "VulnerabilityScanning":
		return "How are vulnerabilities scanned and remediated?"
	case "SecurityAwarenessTraining":
		return "How is security awareness training delivered and measured?"
	case "AuditTrail":
		return "How are audit trails maintained and reviewed?"
	case "BreakGlassAccess":
		return "How is emergency break-glass access controlled and audited?"
	case "FraudMonitoring":
		return "How are fraudulent transactions detected and investigated?"
	case "SettlementIntegrity":
		return "How is settlement integrity verified and reconciled?"
	case "DualControl":
		return "How is dual control enforced for sensitive transactions?"
	case "TransactionLogging":
		return "How are transactions logged and monitored for anomalies?"
	case "BackupRecoveryTesting":
		return "How are backup restores tested and validated?"
	case "ConfigurationHardening":
		return "How are systems hardened and compliance verified?"
	case "PhishingAwareness":
		return "How is phishing awareness tested and improved?"
	case "IncidentResponseAutomation":
		return "How are incident response actions automated?"
	case "ThreatIntelligence":
		return "How is threat intelligence consumed and operationalized?"
	case "SupplyChainSecurity":
		return "How is supply chain security assessed and monitored?"
	default:
		return ""
	}
}

// ─────────────────────────────────────────────────────────────
// PHASE 4 — AUDIT READINESS SCORING
// ─────────────────────────────────────────────────────────────

func computeAuditReadiness(fw ComplianceFramework, coverage *fwCoverage, evidenceScore float64, findingsPenalty float64, threatExposure float64) CIAREAuditReadiness {
	if coverage == nil || coverage.required == 0 {
		return CIAREAuditReadiness{
			Framework:       string(fw),
			ReadinessScore:  0,
			Status:          "Not Assessed",
			ControlCoverage: 0,
			EvidenceScore:   0,
			ThreatExposure:  0,
			FindingsPenalty: 0,
			Factors:         []string{"No controls assessed for framework"},
		}
	}

	controlScore := coverage.pct
	effectiveEvidence := evidenceScore * (coverage.pct / 100.0)
	readiness := (controlScore * 0.45) + (effectiveEvidence * 0.30) - (findingsPenalty * 0.15) - (threatExposure * 0.10)
	if readiness < 0 {
		readiness = 0
	}
	if readiness > 100 {
		readiness = 100
	}

	status := readinessLevelString(readiness)
	factors := []string{
		fmt.Sprintf("Control coverage: %.1f%%", controlScore),
		fmt.Sprintf("Evidence availability: %.1f%%", effectiveEvidence),
		fmt.Sprintf("Findings penalty: %.1f%%", findingsPenalty),
		fmt.Sprintf("Threat exposure: %.1f%%", threatExposure),
	}

	return CIAREAuditReadiness{
		Framework:       string(fw),
		ReadinessScore:  readiness,
		Status:          status,
		ControlCoverage: controlScore,
		EvidenceScore:   effectiveEvidence,
		ThreatExposure:  threatExposure,
		FindingsPenalty: findingsPenalty,
		Factors:         factors,
	}
}

func readinessLevelString(score float64) string {
	switch {
	case score >= 90:
		return "Excellent"
	case score >= 75:
		return "Strong"
	case score >= 50:
		return "Fair"
	case score >= 25:
		return "Weak"
	default:
		return "Poor"
	}
}

// ─────────────────────────────────────────────────────────────
// PHASE 6 — MISSING EVIDENCE DETECTION
// ─────────────────────────────────────────────────────────────

func detectMissingEvidence(fw ComplianceFramework, reqControls []string, observedControls []SDRIControl, evidenceReqs []CIAREEvidenceRequirement) []CIAREMissingEvidence {
	missing := make([]CIAREMissingEvidence, 0)
	controlMap := make(map[string]SDRIControl)
	for _, c := range observedControls {
		controlMap[normalizeControlName(c.Name)] = c
	}

	for _, req := range evidenceReqs {
		if string(fw) != req.Framework {
			continue
		}
		ctrl, ok := controlMap[normalizeControlName(req.Control)]
		if !ok {
			continue
		}
		hasEvidence := len(ctrl.Evidence) > 0
		if hasEvidence {
			continue
		}
		missing = append(missing, CIAREMissingEvidence{
			Framework: string(fw),
			Control:   req.Control,
			Evidences: req.Evidence,
		})
	}

	return missing
}

// ─────────────────────────────────────────────────────────────
// PHASE 8 — COMPLIANCE GAPS
// ─────────────────────────────────────────────────────────────

func generateComplianceGaps(fw ComplianceFramework, reqControls []string, observedControls []SDRIControl) []CIAREComplianceGap {
	gaps := make([]CIAREComplianceGap, 0)
	gapID := 1
	observed := make(map[string]bool)
	for _, c := range observedControls {
		if c.Status != "Missing" {
			observed[normalizeControlName(c.Name)] = true
		}
	}

	for _, ctrl := range reqControls {
		if observed[normalizeControlName(ctrl)] {
			continue
		}
		risk := "High"
		highRiskControls := map[string]bool{
			"MFA": true, "DataEncryptionAtRest": true, "DataEncryptionInTransit": true,
			"AuditLogging": true, "RBAC": true, "AccessReviews": true,
			"IncidentResponsePlan": true, "VulnerabilityScanning": true,
			"PatchManagement": true, "NetworkSegmentation": true,
			"SecretsVault": true, "PrivilegedAccessManagement": true,
		}
		if highRiskControls[ctrl] {
			risk = "Critical"
		}
		gaps = append(gaps, CIAREComplianceGap{
			ID:          fmt.Sprintf("CG-%s-%03d", fw, gapID),
			Framework:   string(fw),
			Requirement: ctrl,
			Observed:    "Not observed",
			Missing:     ctrl + " control not implemented",
			Risk:        risk,
		})
		gapID++
	}

	return gaps
}

// ─────────────────────────────────────────────────────────────
// PHASE 9 — CONTROL MATURITY MODEL
// ─────────────────────────────────────────────────────────────

func maturityLevel(coverage float64) (int, string) {
	switch {
	case coverage >= 90:
		return 5, "Optimized"
	case coverage >= 75:
		return 4, "Managed"
	case coverage >= 50:
		return 3, "Defined"
	case coverage >= 25:
		return 2, "Repeatable"
	default:
		return 1, "Ad Hoc"
	}
}

func estimateControlMaturity(coverageByCategory []SDRICoverage) []CIAREControlMaturity {
	maturities := make([]CIAREControlMaturity, 0)
	for _, cov := range coverageByCategory {
		level, label := maturityLevel(cov.Coverage)
		maturities = append(maturities, CIAREControlMaturity{
			Domain:   cov.Category,
			Level:    level,
			Label:    fmt.Sprintf("Level %d - %s", level, label),
			Coverage: cov.Coverage,
		})
	}
	sort.Slice(maturities, func(i, j int) bool {
		return maturities[i].Domain < maturities[j].Domain
	})
	return maturities
}

// ─────────────────────────────────────────────────────────────
// PHASE 10 — DOMAIN-SPECIFIC COMPLIANCE PACKS
// ─────────────────────────────────────────────────────────────

var domainPacks = map[string]struct {
	Frameworks  []string
	DomainLabel string
	Priorities  []string
	Description string
}{
	"healthcare": {
		Frameworks:  []string{"HIPAA", "SOC2", "NIST800-53"},
		DomainLabel: "Healthcare",
		Priorities:  []string{"PHI Access Controls", "Audit Logging", "Encryption", "Break Glass Access"},
		Description: "Healthcare architecture handling PHI requires HIPAA compliance with strong access controls, audit trails, and encryption.",
	},
	"fintech": {
		Frameworks:  []string{"PCI-DSS", "SOC2", "ISO27001"},
		DomainLabel: "Financial Technology",
		Priorities:  []string{"Transaction Integrity", "Fraud Monitoring", "Dual Control", "Settlement"},
		Description: "Financial services architecture requires PCI-DSS for payment data, with fraud monitoring and transaction integrity controls.",
	},
	"saas": {
		Frameworks:  []string{"SOC2", "ISO27001", "GDPR"},
		DomainLabel: "SaaS / Cloud Service",
		Priorities:  []string{"Availability", "Data Protection", "Access Control", "Incident Response"},
		Description: "SaaS architecture requires SOC2 reporting with focus on availability, confidentiality, and data protection.",
	},
	"kubernetes": {
		Frameworks:  []string{"CIS", "NIST800-53", "SOC2"},
		DomainLabel: "Kubernetes / Container Platform",
		Priorities:  []string{"Admission Control", "Network Policies", "Pod Security", "Image Scanning"},
		Description: "Container platform requires CIS Kubernetes benchmarks with strong admission controls and runtime security.",
	},
	"enterprise": {
		Frameworks:  []string{"ISO27001", "NIST800-53", "SOC2", "GDPR"},
		DomainLabel: "Enterprise",
		Priorities:  []string{"Identity Management", "Access Control", "Audit Logging", "Incident Response"},
		Description: "Enterprise architecture requires broad compliance coverage across information security management.",
	},
}

// ─────────────────────────────────────────────────────────────
// PHASE 11 — COMPLIANCE RISK NARRATIVES
// ─────────────────────────────────────────────────────────────

func generateComplianceNarrative(fw ComplianceFramework, coverage *fwCoverage, gaps int, missingEvidences int) CIAREComplianceNarrative {
	if coverage == nil || coverage.required == 0 {
		return CIAREComplianceNarrative{
			Framework: string(fw),
			Narrative: fmt.Sprintf("The architecture does not declare compliance with %s. No controls were assessed.", fw),
		}
	}

	base := fmt.Sprintf("The architecture declares %s compliance.", fw)
	obs := fmt.Sprintf(" Of %d required controls, %d are observed (%.1f%% coverage).",
		coverage.required, coverage.observed, coverage.pct)
	miss := ""
	if coverage.missing > 0 {
		miss = fmt.Sprintf(" %d controls are missing.", coverage.missing)
	}
	evidenceNarr := ""
	if missingEvidences > 0 {
		evidenceNarr = fmt.Sprintf(" %d controls lack supporting evidence.", missingEvidences)
	}
	gapNarr := ""
	if gaps > 0 {
		gapNarr = fmt.Sprintf(" %d compliance gaps were identified.", gaps)
	}
	reco := " Review and remediate gaps to improve audit readiness."

	return CIAREComplianceNarrative{
		Framework: string(fw),
		Narrative: base + obs + miss + evidenceNarr + gapNarr + reco,
	}
}

// ─────────────────────────────────────────────────────────────
// PHASE 12 — AUDIT PACKAGE GENERATOR
// ─────────────────────────────────────────────────────────────

func buildAuditPackage(
	coverages []CIAREFrameworkCoverage,
	readiness []CIAREAuditReadiness,
	controls []SDRIControl,
	evidenceReqs []CIAREEvidenceRequirement,
	questions []CIAREAuditorQuestion,
	gaps []CIAREComplianceGap,
) CIAREAuditPackage {
	if len(coverages) == 0 {
		return CIAREAuditPackage{
			ExecutiveSummary: "No compliance frameworks assessed.",
		}
	}

	summary := fmt.Sprintf("Compliance assessment covers %d frameworks.", len(coverages))
	totalReqs := 0
	totalObs := 0
	totalGaps := len(gaps)
	for _, c := range coverages {
		totalReqs += c.Required
		totalObs += c.Observed
	}
	if totalReqs > 0 {
		pct := float64(totalObs) / float64(totalReqs) * 100
		summary += fmt.Sprintf(" Total controls required: %d, observed: %d (%.1f%%).", totalReqs, totalObs, pct)
	}
	if totalGaps > 0 {
		summary += fmt.Sprintf(" Compliance gaps: %d.", totalGaps)
	}

	return CIAREAuditPackage{
		ExecutiveSummary:     summary,
		FrameworkCoverages:   coverages,
		ControlInventory:     controls,
		MissingControls:      gaps,
		EvidenceRequirements: evidenceReqs,
		AuditorQuestions:     questions,
	}
}

// ─────────────────────────────────────────────────────────────
// PHASE 13 — COMPLIANCE DASHBOARD
// ─────────────────────────────────────────────────────────────

func buildComplianceDashboard(
	coverages []CIAREFrameworkCoverage,
	readiness []CIAREAuditReadiness,
	gaps []CIAREComplianceGap,
	missingEvidences []CIAREMissingEvidence,
) CIAREComplianceDashboard {
	fwMap := make(map[string]float64)
	for _, c := range coverages {
		fwMap[c.Framework] = c.CoveragePct
	}

	topGaps := gaps
	if len(topGaps) > 10 {
		topGaps = topGaps[:10]
	}

	topMissing := missingEvidences
	if len(topMissing) > 10 {
		topMissing = topMissing[:10]
	}

	topRisks := make([]string, 0)
	for _, g := range gaps {
		if g.Risk == "Critical" {
			topRisks = append(topRisks, g.Framework+": "+g.Requirement)
		}
	}
	if len(topRisks) > 5 {
		topRisks = topRisks[:5]
	}

	return CIAREComplianceDashboard{
		FrameworkCoverages: fwMap,
		TopGaps:            topGaps,
		TopMissingEvidence: topMissing,
		TopRisks:           topRisks,
	}
}

// ─────────────────────────────────────────────────────────────
// PHASE 14 — PROCUREMENT REVIEW MODE
// ─────────────────────────────────────────────────────────────

var procurementQuestions = []CIAREProcurementQuestion{
	{Category: "Access Control", Question: "How does the vendor enforce multi-factor authentication?"},
	{Category: "Access Control", Question: "What role-based access control model is used?"},
	{Category: "Access Control", Question: "How are privileged accounts managed?"},
	{Category: "Access Control", Question: "How often are access reviews conducted?"},
	{Category: "Encryption", Question: "What encryption standards are used for data at rest and in transit?"},
	{Category: "Encryption", Question: "How are encryption keys managed and rotated?"},
	{Category: "Data Protection", Question: "How is data classified and handled based on sensitivity?"},
	{Category: "Data Protection", Question: "What data retention and disposal policies are in place?"},
	{Category: "Data Protection", Question: "How is data loss prevention configured?"},
	{Category: "Audit Logging", Question: "What security events are logged and retained?"},
	{Category: "Audit Logging", Question: "How are logs integrated with a SIEM?"},
	{Category: "Incident Response", Question: "What is the incident response process and SLA?"},
	{Category: "Incident Response", Question: "How are incidents detected and escalated?"},
	{Category: "Backup & Recovery", Question: "What is the backup frequency and retention period?"},
	{Category: "Backup & Recovery", Question: "How are backups tested for restorability?"},
	{Category: "Backup & Recovery", Question: "What is the RTO and RPO?"},
	{Category: "Compliance", Question: "What compliance certifications does the vendor hold?"},
	{Category: "Compliance", Question: "When was the last SOC 2 or equivalent audit?"},
	{Category: "Compliance", Question: "How are compliance gaps tracked and remediated?"},
	{Category: "Third Party", Question: "What sub-processors does the vendor use?"},
	{Category: "Third Party", Question: "How are sub-processors assessed for security?"},
	{Category: "Third Party", Question: "What is the vendor's incident notification SLA?"},
	{Category: "Network Security", Question: "How is the vendor's network segmented?"},
	{Category: "Network Security", Question: "What intrusion detection measures are in place?"},
	{Category: "Vulnerability Management", Question: "How are vulnerabilities identified and remediated?"},
	{Category: "Vulnerability Management", Question: "What is the patch management process?"},
	{Category: "Personnel Security", Question: "How are employees vetted and trained?"},
	{Category: "Personnel Security", Question: "What security awareness training is provided?"},
	{Category: "Physical Security", Question: "How is physical access to data centers controlled?"},
	{Category: "Business Continuity", Question: "What is the business continuity and disaster recovery plan?"},
}

// ─────────────────────────────────────────────────────────────
// HELPERS
// ─────────────────────────────────────────────────────────────

func determineFrameworks(compliance []string) []ComplianceFramework {
	if len(compliance) == 0 {
		return allFrameworks
	}
	fws := make([]ComplianceFramework, 0)
	seen := make(map[ComplianceFramework]bool)
	for _, c := range compliance {
		fw := ComplianceFramework(strings.ToUpper(c))
		for _, af := range allFrameworks {
			if string(af) == string(fw) && !seen[af] {
				fws = append(fws, af)
				seen[af] = true
				break
			}
		}
	}
	if len(fws) == 0 {
		return allFrameworks
	}
	return fws
}

func normalizeControlName(name string) string {
	normalized := strings.ToLower(name)
	normalized = strings.NewReplacer(" ", "", "-", "", "_", "", ".", "").Replace(normalized)
	return normalized
}

func controlObserved(name string, observedControls []SDRIControl) bool {
	n := normalizeControlName(name)
	for _, c := range observedControls {
		if normalizeControlName(c.Name) == n && c.Status != "Missing" {
			return true
		}
	}
	return false
}

func countEvidenceForFramework(fw ComplianceFramework, evidenceReqs []CIAREEvidenceRequirement, observedControls []SDRIControl) float64 {
	total := 0
	withEvidence := 0
	for _, req := range evidenceReqs {
		if req.Framework != string(fw) {
			continue
		}
		total++
		for _, c := range observedControls {
			if normalizeControlName(c.Name) == normalizeControlName(req.Control) && len(c.Evidence) > 0 {
				withEvidence++
				break
			}
		}
	}
	if total == 0 {
		return 0
	}
	return float64(withEvidence) / float64(total) * 100
}

func countFindingsForFramework(fw ComplianceFramework, controls []SDRIControl, findings []SDRIFinding) float64 {
	relevant := 0
	for _, f := range findings {
		for _, c := range controls {
			for _, ac := range f.AffectedControls {
				if normalizeControlName(ac) == normalizeControlName(c.Name) {
					relevant++
					break
				}
			}
		}
	}
	return float64(relevant) * 5.0
}

// ─────────────────────────────────────────────────────────────
// MAIN — CIARE ENGINE RUN
// ─────────────────────────────────────────────────────────────

func (e *CIAREEngine) Run(input CIAREInput) *CIAREResult {
	frameworks := determineFrameworks(input.Compliance)
	mappings := buildControlFrameworkMappings()
	_ = mappings

	observedControls := make([]SDRIControl, 0)
	for _, c := range input.SDRIResult.Controls {
		if c.Status != "Missing" {
			observedControls = append(observedControls, c)
		}
	}

	coverageByFW := make(map[ComplianceFramework]*fwCoverage)
	for _, fw := range frameworks {
		reqControls, ok := frameworkRequiredControls[fw]
		if !ok {
			continue
		}
		observed := 0
		var present []string
		var missing []string
		for _, ctrl := range reqControls {
			if controlObserved(ctrl, observedControls) {
				observed++
				present = append(present, ctrl)
			} else {
				missing = append(missing, ctrl)
			}
		}
		total := len(reqControls)
		pct := 0.0
		if total > 0 {
			pct = float64(observed) / float64(total) * 100
		}
		status := readinessLevelString(pct)
		coverageByFW[fw] = &fwCoverage{
			required:        total,
			observed:        observed,
			missing:         total - observed,
			pct:             pct,
			presentControls: present,
			missingControls: missing,
		}
		_ = status
	}

	// Phase 3: Build framework coverages
	coverages := make([]CIAREFrameworkCoverage, 0, len(frameworks))
	for _, fw := range frameworks {
		cov, ok := coverageByFW[fw]
		if !ok {
			continue
		}
		status := readinessLevelString(cov.pct)
		coverages = append(coverages, CIAREFrameworkCoverage{
			Framework:        string(fw),
			Required:         cov.required,
			Observed:         cov.observed,
			Missing:          cov.missing,
			CoveragePct:      cov.pct,
			Status:           status,
			ObservedControls: cov.presentControls,
			MissingControls:  cov.missingControls,
		})
	}

	// Phase 5: Evidence requirements
	evidenceReqs := make([]CIAREEvidenceRequirement, 0)
	for _, fw := range frameworks {
		reqControls, ok := frameworkRequiredControls[fw]
		if !ok {
			continue
		}
		for _, ctrl := range reqControls {
			ev, ok := controlEvidence[ctrl]
			if !ok {
				ev = []string{"Evidence documentation for " + ctrl}
			}
			evidenceReqs = append(evidenceReqs, CIAREEvidenceRequirement{
				Framework: string(fw),
				Control:   ctrl,
				Evidence:  ev,
			})
		}
	}

	// Phase 6: Missing evidence
	missingEvidences := make([]CIAREMissingEvidence, 0)
	for _, fw := range frameworks {
		reqControls, ok := frameworkRequiredControls[fw]
		if !ok {
			continue
		}
		fwEvidenceReqs := make([]CIAREEvidenceRequirement, 0)
		for _, req := range evidenceReqs {
			if req.Framework == string(fw) {
				fwEvidenceReqs = append(fwEvidenceReqs, req)
			}
		}
		miss := detectMissingEvidence(fw, reqControls, observedControls, fwEvidenceReqs)
		missingEvidences = append(missingEvidences, miss...)
	}

	// Phase 7: Auditor questions
	questions := make([]CIAREAuditorQuestion, 0)
	for _, fw := range frameworks {
		reqControls, ok := frameworkRequiredControls[fw]
		if !ok {
			continue
		}
		cov := coverageByFW[fw]
		fwQs := generateAuditorQuestions(fw, reqControls, cov)
		questions = append(questions, fwQs...)
	}

	// Phase 8: Compliance gaps
	gaps := make([]CIAREComplianceGap, 0)
	for _, fw := range frameworks {
		reqControls, ok := frameworkRequiredControls[fw]
		if !ok {
			continue
		}
		fwGaps := generateComplianceGaps(fw, reqControls, observedControls)
		gaps = append(gaps, fwGaps...)
	}

	// Phase 4: Audit readiness
	readiness := make([]CIAREAuditReadiness, 0)
	for _, fw := range frameworks {
		cov := coverageByFW[fw]
		evScore := countEvidenceForFramework(fw, evidenceReqs, observedControls)
		findingsPenalty := countFindingsForFramework(fw, observedControls, input.SDRIResult.DesignFindings)
		threatExp := 0.0
		if cov != nil {
			threatExp = 100.0 - cov.pct
		}
		ar := computeAuditReadiness(fw, cov, evScore, findingsPenalty, threatExp)
		readiness = append(readiness, ar)
	}

	// Phase 9: Control maturity
	maturities := estimateControlMaturity(input.SDRIResult.CoverageByCategory)

	// Phase 11: Compliance narratives
	narratives := make([]CIAREComplianceNarrative, 0)
	for _, fw := range frameworks {
		cov := coverageByFW[fw]
		gapCount := 0
		missEvCount := 0
		for _, g := range gaps {
			if g.Framework == string(fw) {
				gapCount++
			}
		}
		for _, m := range missingEvidences {
			if m.Framework == string(fw) {
				missEvCount++
			}
		}
		narr := generateComplianceNarrative(fw, cov, gapCount, missEvCount)
		narratives = append(narratives, narr)
	}

	// Phase 12: Audit package
	auditPackage := buildAuditPackage(coverages, readiness, input.SDRIResult.Controls, evidenceReqs, questions, gaps)

	// Phase 13: Compliance dashboard
	dashboard := buildComplianceDashboard(coverages, readiness, gaps, missingEvidences)

	// Phase 14: Procurement questions
	procurementQs := make([]CIAREProcurementQuestion, len(procurementQuestions))
	copy(procurementQs, procurementQuestions)

	// Ensure a valid SDRIResult if input was nil
	if input.SDRIResult == nil {
		input.SDRIResult = &SDRIResult{}
	}

	return &CIAREResult{
		FrameworkCoverages:   coverages,
		AuditReadinessScores: readiness,
		EvidenceRequirements: evidenceReqs,
		MissingEvidences:     missingEvidences,
		AuditorQuestions:     questions,
		ComplianceGaps:       gaps,
		ControlMaturities:    maturities,
		ComplianceNarratives: narratives,
		AuditPackage:         auditPackage,
		ComplianceDashboard:  dashboard,
		ProcurementQuestions: procurementQs,
	}
}
