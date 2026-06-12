package intelligence

import (
	"fmt"
	"sort"
	"strings"
)

// ─────────────────────────────────────────────────────────────
// PHASE 1 — THREAT DATA MODEL
// ─────────────────────────────────────────────────────────────

// ThreatCategory represents the category of a threat.
type ThreatCategory string

const (
	ThreatIdentity       ThreatCategory = "IDENTITY"
	ThreatAccessControl  ThreatCategory = "ACCESS_CONTROL"
	ThreatDataProtection ThreatCategory = "DATA_PROTECTION"
	ThreatKeyManagement  ThreatCategory = "KEY_MANAGEMENT"
	ThreatSecrets        ThreatCategory = "SECRETS_MANAGEMENT"
	ThreatNetwork        ThreatCategory = "NETWORK"
	ThreatMonitoring     ThreatCategory = "MONITORING"
	ThreatBackup         ThreatCategory = "BACKUP"
	ThreatThirdParty     ThreatCategory = "THIRD_PARTY"
	ThreatAvailability   ThreatCategory = "AVAILABILITY"
	ThreatConfiguration  ThreatCategory = "CONFIGURATION"
	ThreatPhysical       ThreatCategory = "PHYSICAL"
)

// Threat represents a generated threat in the threat model.
type Threat struct {
	ID                 string         `json:"id"`
	Name               string         `json:"name"`
	Category           ThreatCategory `json:"category"`
	Severity           RiskLevel      `json:"severity"`
	Likelihood         float64        `json:"likelihood"`
	Impact             float64        `json:"impact"`
	RiskScore          float64        `json:"risk_score"`
	Confidence         float64        `json:"confidence"`
	Description        string         `json:"description"`
	AffectedAssets     []string       `json:"affected_assets,omitempty"`
	AffectedComponents []string       `json:"affected_components,omitempty"`
	AffectedBoundaries []string       `json:"affected_boundaries,omitempty"`
	AffectedData       []string       `json:"affected_data,omitempty"`
	Assumptions        []string       `json:"assumptions,omitempty"`
	Controls           []string       `json:"controls,omitempty"`
	STRIDECategories   []string       `json:"stride_categories,omitempty"`
	Reasoning          string         `json:"reasoning"`
	Recommendations    []string       `json:"recommendations,omitempty"`
	PreventiveControls []string       `json:"preventive_controls,omitempty"`
	DetectiveControls  []string       `json:"detective_controls,omitempty"`
	CorrectiveControls []string       `json:"corrective_controls,omitempty"`
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

// TMIRunResult holds the output of a threat modeling run.
type TMIRunResult struct {
	Threats  []Threat           `json:"threats"`
	Clusters []ThreatCluster    `json:"clusters"`
	Summary  ThreatModelSummary `json:"summary"`
}

// ─────────────────────────────────────────────────────────────
// PHASE 3 — THREAT GENERATION RULES
// ─────────────────────────────────────────────────────────────

// TMIEngine (Threat Modeling Intelligence Engine) generates threats.
type TMIEngine struct {
	componentRules    map[string][]ThreatTemplate
	relationshipRules map[string][]ThreatTemplate
	assumptionRules   map[string][]ThreatTemplate
	boundaryRules     map[string][]ThreatTemplate
	strideMapping     map[ThreatCategory][]string
	controlLibrary    map[string]ControlRecommendations
}

// ThreatTemplate is a template for generating a threat.
type ThreatTemplate struct {
	Name               string
	Category           ThreatCategory
	Description        string
	BaseLikelihood     float64
	BaseImpact         float64
	STRIDECategories   []string
	PreventiveControls []string
	DetectiveControls  []string
	CorrectiveControls []string
}

// ControlRecommendations holds controls for a threat.
type ControlRecommendations struct {
	Preventive []string
	Detective  []string
	Corrective []string
}

// NewTMIEngine creates a new threat modeling engine.
func NewTMIEngine() *TMIEngine {
	return &TMIEngine{
		componentRules: map[string][]ThreatTemplate{
			"identity_provider": {
				{Name: "Admin Account Compromise", Category: ThreatIdentity, Description: "Administrative access to the identity provider may be compromised through credential theft or brute force", BaseLikelihood: 0.3, BaseImpact: 0.9, STRIDECategories: []string{"Spoofing", "Elevation of Privilege"}, PreventiveControls: []string{"MFA", "PAM", "IP Whitelisting"}, DetectiveControls: []string{"Login Monitoring", "SIEM Alerts"}, CorrectiveControls: []string{"Credential Rotation", "Session Revocation"}},
				{Name: "MFA Bypass", Category: ThreatIdentity, Description: "Multi-factor authentication may be bypassed through SIM swapping, social engineering, or compromised MFA infrastructure", BaseLikelihood: 0.2, BaseImpact: 0.9, STRIDECategories: []string{"Spoofing", "Elevation of Privilege"}, PreventiveControls: []string{"Hardware MFA", "Biometric Verification", "Risk-Based Authentication"}, DetectiveControls: []string{"Anomaly Detection", "Login Monitoring"}, CorrectiveControls: []string{"Session Revocation", "Credential Reset"}},
				{Name: "Recovery Flow Abuse", Category: ThreatIdentity, Description: "Account recovery flows may be abused to gain unauthorized access", BaseLikelihood: 0.3, BaseImpact: 0.8, STRIDECategories: []string{"Spoofing"}, PreventiveControls: []string{"Strong Recovery Questions", "Identity Verification", "Time-Delayed Recovery"}, DetectiveControls: []string{"Recovery Attempt Logging", "SIEM Alerts"}, CorrectiveControls: []string{"Account Lockout", "Notification to User"}},
				{Name: "Token Validation Failure", Category: ThreatIdentity, Description: "Tokens may be forged, replayed, or validated incorrectly", BaseLikelihood: 0.25, BaseImpact: 0.85, STRIDECategories: []string{"Spoofing", "Tampering"}, PreventiveControls: []string{"Token Signing", "Clock Skew Tolerance", "Token Binding"}, DetectiveControls: []string{"Token Anomaly Detection", "Replay Detection"}, CorrectiveControls: []string{"Token Revocation", "Key Rotation"}},
				{Name: "Identity Provider Compromise", Category: ThreatIdentity, Description: "The identity provider itself may be compromised, affecting all downstream authentication", BaseLikelihood: 0.15, BaseImpact: 1.0, STRIDECategories: []string{"Spoofing", "Elevation of Privilege"}, PreventiveControls: []string{"IdP Hardening", "MFA", "Regular Audits"}, DetectiveControls: []string{"IdP Monitoring", "SIEM Integration"}, CorrectiveControls: []string{"Failover IdP", "Emergency Access Procedure"}},
			},
			"database": {
				{Name: "Unauthorized Data Access", Category: ThreatDataProtection, Description: "Database may be accessed without proper authorization", BaseLikelihood: 0.4, BaseImpact: 0.9, STRIDECategories: []string{"Information Disclosure", "Elevation of Privilege"}, PreventiveControls: []string{"Least Privilege", "Query Parameterization", "Access Control"}, DetectiveControls: []string{"Database Activity Monitoring", "Audit Logging"}, CorrectiveControls: []string{"Access Revocation", "Data Breach Response"}},
				{Name: "Data Exfiltration", Category: ThreatDataProtection, Description: "Sensitive data may be extracted from the database by authorized or unauthorized actors", BaseLikelihood: 0.3, BaseImpact: 0.95, STRIDECategories: []string{"Information Disclosure"}, PreventiveControls: []string{"DLP", "Access Control", "Encryption at Rest"}, DetectiveControls: []string{"DLP Monitoring", "Anomaly Detection"}, CorrectiveControls: []string{"Account Lockout", "Data Loss Investigation"}},
				{Name: "Privilege Escalation", Category: ThreatAccessControl, Description: "Database privileges may be escalated to gain unauthorized access", BaseLikelihood: 0.25, BaseImpact: 0.85, STRIDECategories: []string{"Elevation of Privilege"}, PreventiveControls: []string{"Least Privilege", "Role-Based Access", "Stored Procedures"}, DetectiveControls: []string{"Privilege Change Logging", "SIEM Alerts"}, CorrectiveControls: []string{"Privilege Revocation", "Access Review"}},
				{Name: "SQL Injection", Category: ThreatDataProtection, Description: "SQL injection may allow unauthorized database access or data manipulation", BaseLikelihood: 0.3, BaseImpact: 0.9, STRIDECategories: []string{"Injection", "Information Disclosure", "Elevation of Privilege"}, PreventiveControls: []string{"Query Parameterization", "Input Validation", "WAF"}, DetectiveControls: []string{"Query Pattern Monitoring", "WAF Logs"}, CorrectiveControls: []string{"Query Blacklisting", "Patching"}},
				{Name: "Data Tampering", Category: ThreatDataProtection, Description: "Database records may be modified without authorization", BaseLikelihood: 0.25, BaseImpact: 0.8, STRIDECategories: []string{"Tampering"}, PreventiveControls: []string{"Integrity Checks", "Access Control", "Immutable Logs"}, DetectiveControls: []string{"Change Detection", "Audit Logging"}, CorrectiveControls: []string{"Backup Restoration", "Investigation"}},
			},
			"encryption_service": {
				{Name: "Key Exposure", Category: ThreatKeyManagement, Description: "Encryption keys may be exposed in memory, logs, or configuration", BaseLikelihood: 0.3, BaseImpact: 0.95, STRIDECategories: []string{"Information Disclosure"}, PreventiveControls: []string{"HSM", "Key Vaulting", "Memory Protection"}, DetectiveControls: []string{"Key Access Logging", "Anomaly Detection"}, CorrectiveControls: []string{"Key Rotation", "Key Revocation"}},
				{Name: "Rotation Failure", Category: ThreatKeyManagement, Description: "Key rotation may fail or be delayed, increasing exposure window", BaseLikelihood: 0.35, BaseImpact: 0.7, STRIDECategories: []string{"Information Disclosure"}, PreventiveControls: []string{"Automated Rotation", "Rotation Policy", "Key Lifecycle Management"}, DetectiveControls: []string{"Rotation Monitoring", "Compliance Scanning"}, CorrectiveControls: []string{"Emergency Rotation", "Key Compromise Procedure"}},
				{Name: "Key Misuse", Category: ThreatKeyManagement, Description: "Keys may be used for unauthorized purposes or by unauthorized entities", BaseLikelihood: 0.25, BaseImpact: 0.8, STRIDECategories: []string{"Elevation of Privilege"}, PreventiveControls: []string{"Key Purpose Separation", "Access Control", "HSM"}, DetectiveControls: []string{"Key Usage Monitoring", "Audit Logging"}, CorrectiveControls: []string{"Key Revocation", "Access Revocation"}},
			},
			"web_application": {
				{Name: "Injection Attack", Category: ThreatDataProtection, Description: "Application may be vulnerable to injection attacks (SQL, NoSQL, OS, LDAP)", BaseLikelihood: 0.35, BaseImpact: 0.85, STRIDECategories: []string{"Injection", "Elevation of Privilege"}, PreventiveControls: []string{"Input Validation", "Query Parameterization", "WAF"}, DetectiveControls: []string{"WAF Monitoring", "SIEM Alerts"}, CorrectiveControls: []string{"Patching", "Input Sanitization"}},
				{Name: "Broken Authentication", Category: ThreatIdentity, Description: "Application authentication may be bypassed or broken", BaseLikelihood: 0.3, BaseImpact: 0.9, STRIDECategories: []string{"Spoofing", "Elevation of Privilege"}, PreventiveControls: []string{"MFA", "Session Management", "Password Policy"}, DetectiveControls: []string{"Login Monitoring", "Brute Force Detection"}, CorrectiveControls: []string{"Account Lockout", "Session Revocation"}},
				{Name: "Sensitive Data Exposure", Category: ThreatDataProtection, Description: "Application may expose sensitive data in transit or at rest", BaseLikelihood: 0.35, BaseImpact: 0.85, STRIDECategories: []string{"Information Disclosure"}, PreventiveControls: []string{"TLS", "Encryption at Rest", "Data Classification"}, DetectiveControls: []string{"SSL/TLS Monitoring", "Data Classification Scanning"}, CorrectiveControls: []string{"Encryption Enforcement", "Data Masking"}},
				{Name: "XML External Entity", Category: ThreatDataProtection, Description: "Application may process XML with external entity references", BaseLikelihood: 0.2, BaseImpact: 0.8, STRIDECategories: []string{"Information Disclosure", "Denial of Service"}, PreventiveControls: []string{"Disable DTD", "XML Parser Hardening", "Input Validation"}, DetectiveControls: []string{"XML Parsing Monitoring", "SIEM Alerts"}, CorrectiveControls: []string{"Parser Update", "Configuration Fix"}},
				{Name: "Broken Access Control", Category: ThreatAccessControl, Description: "Application may fail to enforce access controls properly", BaseLikelihood: 0.35, BaseImpact: 0.85, STRIDECategories: []string{"Elevation of Privilege", "Information Disclosure"}, PreventiveControls: []string{"RBAC", "Object-Level Authorization", "Least Privilege"}, DetectiveControls: []string{"Access Log Monitoring", "Anomaly Detection"}, CorrectiveControls: []string{"Access Revocation", "Policy Update"}},
				{Name: "Security Misconfiguration", Category: ThreatConfiguration, Description: "Application may have insecure default configurations", BaseLikelihood: 0.4, BaseImpact: 0.6, STRIDECategories: []string{"Information Disclosure", "Elevation of Privilege"}, PreventiveControls: []string{"Hardening", "Configuration Management", "CIS Benchmarks"}, DetectiveControls: []string{"Configuration Scanning", "Compliance Monitoring"}, CorrectiveControls: []string{"Configuration Remediation", "Patching"}},
				{Name: "Cross-Site Scripting", Category: ThreatDataProtection, Description: "Application may be vulnerable to XSS attacks", BaseLikelihood: 0.35, BaseImpact: 0.7, STRIDECategories: []string{"Injection", "Information Disclosure"}, PreventiveControls: []string{"Output Encoding", "CSP", "Input Validation"}, DetectiveControls: []string{"XSS Monitoring", "WAF Logs"}, CorrectiveControls: []string{"Code Fix", "WAF Rule Update"}},
				{Name: "Insecure Deserialization", Category: ThreatDataProtection, Description: "Application may deserialize untrusted data insecurely", BaseLikelihood: 0.25, BaseImpact: 0.9, STRIDECategories: []string{"Injection", "Elevation of Privilege"}, PreventiveControls: []string{"Serialization Hardening", "Input Validation", "Type Checking"}, DetectiveControls: []string{"Deserialization Monitoring", "SIEM Alerts"}, CorrectiveControls: []string{"Library Update", "Code Fix"}},
				{Name: "Insufficient Logging", Category: ThreatMonitoring, Description: "Application may not log security-relevant events", BaseLikelihood: 0.45, BaseImpact: 0.5, STRIDECategories: []string{"Repudiation"}, PreventiveControls: []string{"Comprehensive Logging", "Structured Logging", "Log Integrity"}, DetectiveControls: []string{"Log Monitoring", "Log Analysis"}, CorrectiveControls: []string{"Logging Enhancement", "Log Retention"}},
			},
			"api_gateway": {
				{Name: "API Abuse", Category: ThreatNetwork, Description: "API may be abused through excessive requests, scraping, or enumeration", BaseLikelihood: 0.4, BaseImpact: 0.6, STRIDECategories: []string{"Denial of Service", "Information Disclosure"}, PreventiveControls: []string{"Rate Limiting", "API Authentication", "Throttling"}, DetectiveControls: []string{"Traffic Monitoring", "Anomaly Detection"}, CorrectiveControls: []string{"Rate Limit Adjustment", "IP Blocking"}},
				{Name: "Request Smuggling", Category: ThreatNetwork, Description: "HTTP request smuggling may bypass gateway controls", BaseLikelihood: 0.2, BaseImpact: 0.8, STRIDECategories: []string{"Injection", "Elevation of Privilege"}, PreventiveControls: []string{"HTTP/2", "Request Validation", "Gateway Hardening"}, DetectiveControls: []string{"Gateway Monitoring", "SIEM Alerts"}, CorrectiveControls: []string{"Gateway Update", "Configuration Fix"}},
				{Name: "Credential Stuffing", Category: ThreatIdentity, Description: "API authentication may be attacked through credential stuffing", BaseLikelihood: 0.45, BaseImpact: 0.7, STRIDECategories: []string{"Spoofing"}, PreventiveControls: []string{"Rate Limiting", "CAPTCHA", "MFA"}, DetectiveControls: []string{"Login Monitoring", "Brute Force Detection"}, CorrectiveControls: []string{"Account Lockout", "Credential Reset"}},
			},
			"client": {
				{Name: "Client Compromise", Category: ThreatNetwork, Description: "Client devices may be compromised with malware or keyloggers", BaseLikelihood: 0.4, BaseImpact: 0.7, STRIDECategories: []string{"Spoofing", "Information Disclosure"}, PreventiveControls: []string{"Endpoint Protection", "Device Management", "Secure Boot"}, DetectiveControls: []string{"Endpoint Monitoring", "EDR"}, CorrectiveControls: []string{"Device Quarantine", "Credential Reset"}},
				{Name: "Session Hijacking", Category: ThreatIdentity, Description: "Client sessions may be hijacked through XSS, network sniffing, or malware", BaseLikelihood: 0.3, BaseImpact: 0.8, STRIDECategories: []string{"Spoofing", "Information Disclosure"}, PreventiveControls: []string{"Secure Cookies", "Session Binding", "TLS"}, DetectiveControls: []string{"Session Monitoring", "Anomaly Detection"}, CorrectiveControls: []string{"Session Revocation", "Credential Reset"}},
			},
			"logging_service": {
				{Name: "Log Tampering", Category: ThreatMonitoring, Description: "Audit logs may be tampered with or deleted to hide evidence", BaseLikelihood: 0.25, BaseImpact: 0.7, STRIDECategories: []string{"Tampering", "Repudiation"}, PreventiveControls: []string{"Log Integrity", "Immutable Logs", "WORM Storage"}, DetectiveControls: []string{"Log Integrity Monitoring", "Tamper Detection"}, CorrectiveControls: []string{"Log Restoration", "Investigation"}},
				{Name: "Missing Detection", Category: ThreatMonitoring, Description: "Security events may not be detected due to incomplete logging or monitoring", BaseLikelihood: 0.4, BaseImpact: 0.6, STRIDECategories: []string{"Repudiation"}, PreventiveControls: []string{"Comprehensive Logging", "SIEM", "Alert Rules"}, DetectiveControls: []string{"Monitoring Coverage Audit", "Gap Analysis"}, CorrectiveControls: []string{"Monitoring Enhancement", "Rule Update"}},
				{Name: "Alert Fatigue", Category: ThreatMonitoring, Description: "Excessive alerts may cause security events to be ignored", BaseLikelihood: 0.5, BaseImpact: 0.4, STRIDECategories: []string{"Repudiation"}, PreventiveControls: []string{"Alert Tuning", "Prioritization", "SOAR"}, DetectiveControls: []string{"Alert Analysis", "False Positive Tracking"}, CorrectiveControls: []string{"Alert Tuning", "Process Improvement"}},
			},
			"storage_service": {
				{Name: "Backup Exposure", Category: ThreatBackup, Description: "Backup data may be exposed to unauthorized access", BaseLikelihood: 0.3, BaseImpact: 0.85, STRIDECategories: []string{"Information Disclosure"}, PreventiveControls: []string{"Backup Encryption", "Access Control", "Immutable Backups"}, DetectiveControls: []string{"Backup Access Monitoring", "Audit Logging"}, CorrectiveControls: []string{"Backup Rotation", "Access Revocation"}},
				{Name: "Restore Failure", Category: ThreatBackup, Description: "Backups may fail to restore when needed", BaseLikelihood: 0.25, BaseImpact: 0.8, STRIDECategories: []string{"Denial of Service"}, PreventiveControls: []string{"Restore Testing", "Backup Validation", "Geographic Distribution"}, DetectiveControls: []string{"Restore Test Monitoring", "Backup Integrity Checks"}, CorrectiveControls: []string{"Backup Repair", "Alternative Restore"}},
				{Name: "Backup Corruption", Category: ThreatBackup, Description: "Backup data may be corrupted or incomplete", BaseLikelihood: 0.2, BaseImpact: 0.75, STRIDECategories: []string{"Tampering", "Denial of Service"}, PreventiveControls: []string{"Integrity Checks", "Checksums", "Versioning"}, DetectiveControls: []string{"Integrity Monitoring", "Backup Validation"}, CorrectiveControls: []string{"Backup Restoration", "Data Repair"}},
			},
			"external_service": {
				{Name: "Vendor Compromise", Category: ThreatThirdParty, Description: "Third-party service may be compromised, affecting the architecture", BaseLikelihood: 0.25, BaseImpact: 0.85, STRIDECategories: []string{"Spoofing", "Information Disclosure", "Elevation of Privilege"}, PreventiveControls: []string{"Vendor Assessment", "Contractual Controls", "Monitoring"}, DetectiveControls: []string{"Vendor Monitoring", "Third-Party Risk Assessment"}, CorrectiveControls: []string{"Vendor Replacement", "Incident Response"}},
				{Name: "SaaS Breach", Category: ThreatThirdParty, Description: "SaaS provider may experience a data breach affecting customer data", BaseLikelihood: 0.2, BaseImpact: 0.9, STRIDECategories: []string{"Information Disclosure"}, PreventiveControls: []string{"Data Minimization", "Encryption", "DLP"}, DetectiveControls: []string{"Third-Party Monitoring", "Breach Notification"}, CorrectiveControls: []string{"Data Recall", "Service Migration"}},
				{Name: "Identity Provider Failure", Category: ThreatThirdParty, Description: "Third-party identity provider may experience outage or compromise", BaseLikelihood: 0.2, BaseImpact: 0.8, STRIDECategories: []string{"Denial of Service", "Spoofing"}, PreventiveControls: []string{"Failover IdP", "Offline Authentication", "Multiple IdPs"}, DetectiveControls: []string{"IdP Health Monitoring", "Availability Monitoring"}, CorrectiveControls: []string{"Failover Activation", "Emergency Access"}},
			},
			"admin_tool": {
				{Name: "Privilege Abuse", Category: ThreatAccessControl, Description: "Administrative privileges may be abused for unauthorized actions", BaseLikelihood: 0.3, BaseImpact: 0.9, STRIDECategories: []string{"Elevation of Privilege", "Tampering"}, PreventiveControls: []string{"PAM", "Just-In-Time Access", "MFA"}, DetectiveControls: []string{"Admin Activity Monitoring", "Session Recording"}, CorrectiveControls: []string{"Access Revocation", "Investigation"}},
				{Name: "Admin Account Compromise", Category: ThreatIdentity, Description: "Admin accounts may be compromised due to weak credentials or phishing", BaseLikelihood: 0.25, BaseImpact: 0.95, STRIDECategories: []string{"Spoofing", "Elevation of Privilege"}, PreventiveControls: []string{"MFA", "PAM", "Phishing Protection"}, DetectiveControls: []string{"Login Monitoring", "Anomaly Detection"}, CorrectiveControls: []string{"Account Lockout", "Credential Reset"}},
			},
			"vpn": {
				{Name: "VPN Tunnel Compromise", Category: ThreatNetwork, Description: "VPN tunnel may be compromised through weak cryptography or endpoint compromise", BaseLikelihood: 0.2, BaseImpact: 0.85, STRIDECategories: []string{"Spoofing", "Information Disclosure"}, PreventiveControls: []string{"Strong Cryptography", "Certificate Pinning", "Endpoint Verification"}, DetectiveControls: []string{"VPN Monitoring", "Anomaly Detection"}, CorrectiveControls: []string{"Tunnel Reset", "Certificate Revocation"}},
			},
			"jump_host": {
				{Name: "Jump Host Compromise", Category: ThreatNetwork, Description: "Jump host may be compromised, providing lateral movement access", BaseLikelihood: 0.25, BaseImpact: 0.9, STRIDECategories: []string{"Spoofing", "Elevation of Privilege"}, PreventiveControls: []string{"Hardening", "MFA", "Session Recording"}, DetectiveControls: []string{"Session Monitoring", "Activity Logging"}, CorrectiveControls: []string{"Access Revocation", "Investigation"}},
			},
			"firewall": {
				{Name: "Firewall Misconfiguration", Category: ThreatConfiguration, Description: "Firewall rules may be misconfigured, allowing unauthorized traffic", BaseLikelihood: 0.35, BaseImpact: 0.7, STRIDECategories: []string{"Information Disclosure", "Elevation of Privilege"}, PreventiveControls: []string{"Rule Review", "Change Management", "Default Deny"}, DetectiveControls: []string{"Rule Monitoring", "Traffic Analysis"}, CorrectiveControls: []string{"Rule Update", "Configuration Fix"}},
			},
			"waf": {
				{Name: "WAF Bypass", Category: ThreatNetwork, Description: "Web Application Firewall may be bypassed through obfuscation or protocol abuse", BaseLikelihood: 0.3, BaseImpact: 0.75, STRIDECategories: []string{"Injection", "Elevation of Privilege"}, PreventiveControls: []string{"Layered Defenses", "Input Validation", "Protocol Hardening"}, DetectiveControls: []string{"WAF Monitoring", "Bypass Detection"}, CorrectiveControls: []string{"Rule Update", "Signature Update"}},
			},
			"cache": {
				{Name: "Cache Poisoning", Category: ThreatDataProtection, Description: "Cache may be poisoned to serve malicious content to users", BaseLikelihood: 0.25, BaseImpact: 0.7, STRIDECategories: []string{"Tampering", "Information Disclosure"}, PreventiveControls: []string{"Cache Validation", "Input Sanitization", "Origin Authentication"}, DetectiveControls: []string{"Cache Monitoring", "Content Validation"}, CorrectiveControls: []string{"Cache Purge", "Origin Update"}},
			},
			"message_queue": {
				{Name: "Message Queue Poisoning", Category: ThreatDataProtection, Description: "Messages may be poisoned to disrupt processing or inject malicious commands", BaseLikelihood: 0.25, BaseImpact: 0.75, STRIDECategories: []string{"Tampering", "Injection"}, PreventiveControls: []string{"Message Validation", "Schema Enforcement", "Authentication"}, DetectiveControls: []string{"Message Monitoring", "Dead Letter Analysis"}, CorrectiveControls: []string{"Queue Purge", "Consumer Update"}},
			},
		},
		relationshipRules: map[string][]ThreatTemplate{
			"HTTPS": {
				{Name: "TLS Downgrade", Category: ThreatNetwork, Description: "TLS connection may be downgraded to weak cipher or plaintext", BaseLikelihood: 0.2, BaseImpact: 0.8, STRIDECategories: []string{"Information Disclosure", "Tampering"}, PreventiveControls: []string{"TLS 1.3", "HSTS", "Cipher Suite Restrictions"}, DetectiveControls: []string{"TLS Monitoring", "Cipher Analysis"}, CorrectiveControls: []string{"Configuration Update", "Protocol Enforcement"}},
				{Name: "Certificate Validation Failure", Category: ThreatNetwork, Description: "TLS certificate may not be validated properly, allowing MITM", BaseLikelihood: 0.25, BaseImpact: 0.85, STRIDECategories: []string{"Spoofing", "Information Disclosure"}, PreventiveControls: []string{"Certificate Pinning", "CA Validation", "OCSP Stapling"}, DetectiveControls: []string{"Certificate Monitoring", "Validation Testing"}, CorrectiveControls: []string{"Certificate Update", "Pinning Update"}},
			},
			"HTTP": {
				{Name: "Man-in-the-Middle", Category: ThreatNetwork, Description: "Unencrypted HTTP traffic may be intercepted and modified", BaseLikelihood: 0.5, BaseImpact: 0.8, STRIDECategories: []string{"Information Disclosure", "Tampering", "Spoofing"}, PreventiveControls: []string{"TLS Enforcement", "HSTS", "Redirect to HTTPS"}, DetectiveControls: []string{"Traffic Monitoring", "Protocol Analysis"}, CorrectiveControls: []string{"Protocol Update", "Configuration Fix"}},
				{Name: "Credential Exposure", Category: ThreatNetwork, Description: "Credentials may be exposed in transit over HTTP", BaseLikelihood: 0.55, BaseImpact: 0.9, STRIDECategories: []string{"Information Disclosure"}, PreventiveControls: []string{"TLS", "Password Hashing", "Token-Based Auth"}, DetectiveControls: []string{"Network Monitoring", "DLP"}, CorrectiveControls: []string{"Credential Reset", "Protocol Update"}},
			},
			"TLS": {
				{Name: "Certificate Compromise", Category: ThreatNetwork, Description: "TLS certificate may be compromised or fraudulently issued", BaseLikelihood: 0.15, BaseImpact: 0.85, STRIDECategories: []string{"Spoofing", "Information Disclosure"}, PreventiveControls: []string{"Certificate Transparency", "Short-Lived Certificates", "Pinning"}, DetectiveControls: []string{"Certificate Monitoring", "CT Log Monitoring"}, CorrectiveControls: []string{"Certificate Revocation", "Reissue"}},
			},
			"VPN": {
				{Name: "VPN Tunnel Hijacking", Category: ThreatNetwork, Description: "VPN tunnel may be hijacked or intercepted", BaseLikelihood: 0.2, BaseImpact: 0.8, STRIDECategories: []string{"Spoofing", "Information Disclosure"}, PreventiveControls: []string{"Strong Cryptography", "Multi-Factor VPN", "Split Tunneling Control"}, DetectiveControls: []string{"VPN Monitoring", "Anomaly Detection"}, CorrectiveControls: []string{"Tunnel Reset", "Certificate Revocation"}},
			},
		},
		assumptionRules: map[string][]ThreatTemplate{
			"mfa": {
				{Name: "MFA Bypass", Category: ThreatIdentity, Description: "If MFA is not enforced, credential compromise may allow unauthorized access", BaseLikelihood: 0.4, BaseImpact: 0.9, STRIDECategories: []string{"Spoofing", "Elevation of Privilege"}, PreventiveControls: []string{"MFA Enforcement", "Risk-Based Authentication"}, DetectiveControls: []string{"Login Monitoring", "Anomaly Detection"}, CorrectiveControls: []string{"MFA Rollout", "Account Lockout"}},
			},
			"encryption": {
				{Name: "Data Exposure", Category: ThreatDataProtection, Description: "If encryption is not enforced, data may be exposed in transit or at rest", BaseLikelihood: 0.45, BaseImpact: 0.9, STRIDECategories: []string{"Information Disclosure"}, PreventiveControls: []string{"Encryption Enforcement", "TLS Everywhere"}, DetectiveControls: []string{"Traffic Analysis", "Data Classification Scanning"}, CorrectiveControls: []string{"Encryption Rollout", "Data Breach Response"}},
			},
			"backup": {
				{Name: "Backup Data Exposure", Category: ThreatBackup, Description: "If backups are not encrypted, backup data may be exposed", BaseLikelihood: 0.35, BaseImpact: 0.85, STRIDECategories: []string{"Information Disclosure"}, PreventiveControls: []string{"Backup Encryption", "Access Control"}, DetectiveControls: []string{"Backup Access Monitoring", "Audit Logging"}, CorrectiveControls: []string{"Backup Re-encryption", "Access Revocation"}},
			},
			"audit": {
				{Name: "Undetected Activity", Category: ThreatMonitoring, Description: "If audit logging is incomplete, malicious activity may go undetected", BaseLikelihood: 0.4, BaseImpact: 0.7, STRIDECategories: []string{"Repudiation"}, PreventiveControls: []string{"Comprehensive Logging", "Log Integrity"}, DetectiveControls: []string{"Log Coverage Audit", "Gap Analysis"}, CorrectiveControls: []string{"Logging Enhancement", "Investigation"}},
			},
			"least privilege": {
				{Name: "Privilege Escalation", Category: ThreatAccessControl, Description: "If least privilege is not enforced, users may gain excessive access", BaseLikelihood: 0.35, BaseImpact: 0.8, STRIDECategories: []string{"Elevation of Privilege"}, PreventiveControls: []string{"RBAC", "Regular Access Reviews", "JIT Access"}, DetectiveControls: []string{"Privilege Monitoring", "Access Reviews"}, CorrectiveControls: []string{"Privilege Revocation", "Access Review"}},
			},
		},
		boundaryRules: map[string][]ThreatTemplate{
			string(CrossingPUBLIC_TO_INTERNAL): {
				{Name: "DDoS Attack", Category: ThreatAvailability, Description: "Public-facing boundary may be targeted by DDoS attacks", BaseLikelihood: 0.35, BaseImpact: 0.7, STRIDECategories: []string{"Denial of Service"}, PreventiveControls: []string{"DDoS Protection", "Rate Limiting", "CDN"}, DetectiveControls: []string{"Traffic Monitoring", "DDoS Detection"}, CorrectiveControls: []string{"Traffic Filtering", "Scaling"}},
				{Name: "Bot Attack", Category: ThreatNetwork, Description: "Automated bots may attack the public boundary", BaseLikelihood: 0.45, BaseImpact: 0.5, STRIDECategories: []string{"Denial of Service", "Information Disclosure"}, PreventiveControls: []string{"Bot Detection", "CAPTCHA", "Rate Limiting"}, DetectiveControls: []string{"Bot Monitoring", "Traffic Analysis"}, CorrectiveControls: []string{"IP Blocking", "Rule Update"}},
			},
			string(CrossingIDENTITY_TO_APPLICATION): {
				{Name: "Token Replay", Category: ThreatIdentity, Description: "Identity tokens may be replayed across the boundary", BaseLikelihood: 0.25, BaseImpact: 0.8, STRIDECategories: []string{"Spoofing", "Tampering"}, PreventiveControls: []string{"Token Binding", "Replay Detection", "Short-Lived Tokens"}, DetectiveControls: []string{"Token Monitoring", "Anomaly Detection"}, CorrectiveControls: []string{"Token Revocation", "Session Reset"}},
			},
			string(CrossingAPPLICATION_TO_DATA): {
				{Name: "Data Exfiltration", Category: ThreatDataProtection, Description: "Data may be exfiltrated across the application-to-data boundary", BaseLikelihood: 0.3, BaseImpact: 0.9, STRIDECategories: []string{"Information Disclosure"}, PreventiveControls: []string{"DLP", "Access Control", "Query Limiting"}, DetectiveControls: []string{"DLP Monitoring", "Query Analysis"}, CorrectiveControls: []string{"Access Revocation", "Data Breach Response"}},
				{Name: "Privilege Escalation", Category: ThreatAccessControl, Description: "Database privileges may be escalated across the boundary", BaseLikelihood: 0.25, BaseImpact: 0.85, STRIDECategories: []string{"Elevation of Privilege"}, PreventiveControls: []string{"Least Privilege", "Stored Procedures", "Query Parameterization"}, DetectiveControls: []string{"Privilege Monitoring", "Query Analysis"}, CorrectiveControls: []string{"Privilege Revocation", "Access Review"}},
			},
			string(CrossingADMIN_TO_PRODUCTION): {
				{Name: "Insider Threat", Category: ThreatAccessControl, Description: "Administrative access may be used for insider threats", BaseLikelihood: 0.2, BaseImpact: 0.95, STRIDECategories: []string{"Elevation of Privilege", "Tampering"}, PreventiveControls: []string{"PAM", "Session Recording", "Approval Workflows"}, DetectiveControls: []string{"Admin Monitoring", "Session Recording"}, CorrectiveControls: []string{"Access Revocation", "Investigation"}},
			},
			string(CrossingAPPLICATION_TO_THIRD_PARTY): {
				{Name: "Data Leakage", Category: ThreatThirdParty, Description: "Data may leak to third-party services", BaseLikelihood: 0.3, BaseImpact: 0.8, STRIDECategories: []string{"Information Disclosure"}, PreventiveControls: []string{"DLP", "Data Minimization", "Contractual Controls"}, DetectiveControls: []string{"Outbound Monitoring", "Third-Party Assessment"}, CorrectiveControls: []string{"Data Recall", "Contract Enforcement"}},
			},
			string(CrossingSECRETS_ACCESS): {
				{Name: "Secret Extraction", Category: ThreatSecrets, Description: "Secrets may be extracted from the secrets management boundary", BaseLikelihood: 0.25, BaseImpact: 0.9, STRIDECategories: []string{"Information Disclosure"}, PreventiveControls: []string{"HSM", "Access Audit", "Least Privilege"}, DetectiveControls: []string{"Secret Access Monitoring", "Anomaly Detection"}, CorrectiveControls: []string{"Secret Rotation", "Access Revocation"}},
			},
		},
	}
}

// ─────────────────────────────────────────────────────────────
// PHASE 4 — THREAT GENERATION FROM ASSUMPTIONS
// ─────────────────────────────────────────────────────────────

// GenerateThreatsFromAssumptions generates threats by checking "what if this assumption is false?"
func (tmi *TMIEngine) GenerateThreatsFromAssumptions(assumptions []Assumption, arch *ArchDescription) []Threat {
	var threats []Threat
	threatID := 1

	for _, assumption := range assumptions {
		lowerText := strings.ToLower(assumption.Description)
		for keyword, templates := range tmi.assumptionRules {
			if strings.Contains(lowerText, keyword) {
				for _, tmpl := range templates {
					threat := tmi.buildThreatFromTemplate(tmpl, fmt.Sprintf("THREAT-%03d", threatID), assumption.Description)
					threat.Assumptions = []string{assumption.ID}
					threat.AffectedComponents = assumption.SourceComponents
					threats = append(threats, threat)
					threatID++
				}
			}
		}
	}

	return threats
}

// ─────────────────────────────────────────────────────────────
// PHASE 5 — THREAT GENERATION FROM TRUST BOUNDARIES
// ─────────────────────────────────────────────────────────────

// GenerateThreatsFromBoundaries generates threats from trust boundaries.
func (tmi *TMIEngine) GenerateThreatsFromBoundaries(boundaries []TBITrustBoundary, arch *ArchDescription) []Threat {
	var threats []Threat
	threatID := 1

	for _, boundary := range boundaries {
		crossingKey := string(boundary.CrossingType)
		templates, exists := tmi.boundaryRules[crossingKey]
		if !exists {
			continue
		}
		for _, tmpl := range templates {
			threat := tmi.buildThreatFromTemplate(tmpl, fmt.Sprintf("THREAT-%03d", threatID), fmt.Sprintf("Boundary %s: %s", boundary.ID, boundary.Reasoning))
			threat.AffectedBoundaries = []string{boundary.ID}
			threat.AffectedComponents = []string{boundary.SourceZone, boundary.DestinationZone}
			// Boost risk based on boundary risk
			if boundary.Risk == RiskCritical {
				threat.Likelihood += 0.1
				threat.Impact += 0.1
			}
			threats = append(threats, threat)
			threatID++
		}
	}

	return threats
}

// ─────────────────────────────────────────────────────────────
// PHASE 3 — THREAT GENERATION FROM COMPONENTS
// ─────────────────────────────────────────────────────────────

// GenerateThreatsFromComponents generates threats based on component types.
func (tmi *TMIEngine) GenerateThreatsFromComponents(arch *ArchDescription) []Threat {
	var threats []Threat
	threatID := 1

	for _, comp := range arch.Components {
		compLabel := strings.ToLower(comp.Label)
		var templates []ThreatTemplate
		// Try to match by label against component rules
		for ruleType, ruleTemplates := range tmi.componentRules {
			if strings.Contains(compLabel, ruleType) || ruleType == compLabel {
				templates = append(templates, ruleTemplates...)
			}
		}
		for _, tmpl := range templates {
			threat := tmi.buildThreatFromTemplate(tmpl, fmt.Sprintf("THREAT-%03d", threatID), fmt.Sprintf("Component %s: %s", comp.Label, tmpl.Description))
			threat.AffectedComponents = []string{comp.Label}
			threats = append(threats, threat)
			threatID++
		}
	}

	return threats
}

// GenerateThreatsFromRelationships generates threats from relationships.
func (tmi *TMIEngine) GenerateThreatsFromRelationships(arch *ArchDescription) []Threat {
	var threats []Threat
	threatID := 1

	for _, rel := range arch.Relationships {
		protocol := strings.ToUpper(rel.Label)
		templates, exists := tmi.relationshipRules[protocol]
		if !exists {
			continue
		}
		for _, tmpl := range templates {
			threat := tmi.buildThreatFromTemplate(tmpl, fmt.Sprintf("THREAT-%03d", threatID), fmt.Sprintf("Relationship %s -> %s (%s): %s", rel.Source, rel.Target, rel.Label, tmpl.Description))
			threat.AffectedComponents = []string{rel.Source, rel.Target}
			threats = append(threats, threat)
			threatID++
		}
	}

	return threats
}

// ─────────────────────────────────────────────────────────────
// PHASE 8 — THREAT SEVERITY ENGINE
// ─────────────────────────────────────────────────────────────

// buildThreatFromTemplate creates a Threat from a template with scoring.
func (tmi *TMIEngine) buildThreatFromTemplate(tmpl ThreatTemplate, id string, reasoning string) Threat {
	likelihood := tmpl.BaseLikelihood
	impact := tmpl.BaseImpact

	// Apply category-based adjustments
	switch tmpl.Category {
	case ThreatIdentity:
		impact += 0.05
	case ThreatDataProtection:
		impact += 0.05
	case ThreatNetwork:
		likelihood += 0.05
	}

	// Clamp values
	likelihood = clamp(likelihood, 0.1, 1.0)
	impact = clamp(impact, 0.1, 1.0)

	riskScore := likelihood * impact

	// Determine severity
	severity := RiskLow
	if riskScore >= 0.6 {
		severity = RiskCritical
	} else if riskScore >= 0.4 {
		severity = RiskHigh
	} else if riskScore >= 0.2 {
		severity = RiskMedium
	}

	return Threat{
		ID:                 id,
		Name:               tmpl.Name,
		Category:           tmpl.Category,
		Severity:           severity,
		Likelihood:         likelihood,
		Impact:             impact,
		RiskScore:          riskScore,
		Confidence:         0.75,
		Description:        tmpl.Description,
		STRIDECategories:   tmpl.STRIDECategories,
		Reasoning:          reasoning,
		PreventiveControls: tmpl.PreventiveControls,
		DetectiveControls:  tmpl.DetectiveControls,
		CorrectiveControls: tmpl.CorrectiveControls,
		Recommendations:    append(tmpl.PreventiveControls, append(tmpl.DetectiveControls, tmpl.CorrectiveControls...)...),
	}
}

// clamp ensures a value is within [min, max].
func clamp(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}

// ─────────────────────────────────────────────────────────────
// PHASE 6 — STRIDE THREAT CORRELATION
// ─────────────────────────────────────────────────────────────

// BuildSTRIDEDistribution builds a distribution of threats across STRIDE categories.
func (tmi *TMIEngine) BuildSTRIDEDistribution(threats []Threat) map[string]int {
	distribution := make(map[string]int)
	for _, threat := range threats {
		for _, stride := range threat.STRIDECategories {
			distribution[stride]++
		}
	}
	return distribution
}

// ─────────────────────────────────────────────────────────────
// PHASE 9 — THREAT CLUSTERING
// ─────────────────────────────────────────────────────────────

// ClusterThreats groups related threats by category.
func (tmi *TMIEngine) ClusterThreats(threats []Threat) []ThreatCluster {
	clusters := make(map[string]*ThreatCluster)

	for _, threat := range threats {
		catKey := string(threat.Category)
		cluster, exists := clusters[catKey]
		if !exists {
			cluster = &ThreatCluster{
				ID:              fmt.Sprintf("CLUSTER-%s", catKey),
				Name:            fmt.Sprintf("%s Threat Cluster", catKey),
				Category:        catKey,
				Threats:         []string{},
				AffectedAssets:  []string{},
				Recommendations: []string{},
			}
			clusters[catKey] = cluster
		}

		cluster.Threats = append(cluster.Threats, threat.ID)
		cluster.RiskScore += threat.RiskScore
		cluster.AffectedAssets = appendUnique(cluster.AffectedAssets, threat.AffectedComponents...)
		cluster.AffectedAssets = appendUnique(cluster.AffectedAssets, threat.AffectedAssets...)
		cluster.Recommendations = appendUnique(cluster.Recommendations, threat.Recommendations...)
	}

	// Convert to slice
	var result []ThreatCluster
	for _, cluster := range clusters {
		result = append(result, *cluster)
	}

	// Sort by risk score descending
	sort.Slice(result, func(i, j int) bool {
		return result[i].RiskScore > result[j].RiskScore
	})

	return result
}

func appendUnique(slice []string, items ...string) []string {
	seen := make(map[string]bool)
	for _, s := range slice {
		seen[s] = true
	}
	for _, item := range items {
		if !seen[item] {
			seen[item] = true
			slice = append(slice, item)
		}
	}
	return slice
}

// ─────────────────────────────────────────────────────────────
// PHASE 11 — THREAT MODEL SUMMARY
// ─────────────────────────────────────────────────────────────

// BuildSummary creates a summary of the threat model.
func (tmi *TMIEngine) BuildSummary(threats []Threat, clusters []ThreatCluster) ThreatModelSummary {
	counts := map[RiskLevel]int{
		RiskCritical: 0,
		RiskHigh:     0,
		RiskMedium:   0,
		RiskLow:      0,
	}
	for _, t := range threats {
		counts[t.Severity]++
	}

	strideDist := tmi.BuildSTRIDEDistribution(threats)

	// Top 5 threats by risk score
	sorted := make([]Threat, len(threats))
	copy(sorted, threats)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].RiskScore > sorted[j].RiskScore
	})
	var topThreats []string
	for i, t := range sorted {
		if i >= 5 {
			break
		}
		topThreats = append(topThreats, t.Name)
	}

	summaryText := fmt.Sprintf(
		"Threat Model: %d threats identified (%d Critical, %d High, %d Medium, %d Low). %d threat clusters. Top risks: %s.",
		len(threats),
		counts[RiskCritical],
		counts[RiskHigh],
		counts[RiskMedium],
		counts[RiskLow],
		len(clusters),
		strings.Join(topThreats, ", "),
	)

	return ThreatModelSummary{
		TotalThreats:       len(threats),
		CriticalCount:      counts[RiskCritical],
		HighCount:          counts[RiskHigh],
		MediumCount:        counts[RiskMedium],
		LowCount:           counts[RiskLow],
		ClusterCount:       len(clusters),
		STRIDEDistribution: strideDist,
		TopThreats:         topThreats,
		SummaryText:        summaryText,
	}
}

// ─────────────────────────────────────────────────────────────
// PHASE 13 — TMI ENGINE RUN
// ─────────────────────────────────────────────────────────────

// Run executes the full threat modeling pipeline.
func (tmi *TMIEngine) Run(arch *ArchDescription, assumptions []Assumption, boundaries []TBITrustBoundary) *TMIRunResult {
	var allThreats []Threat

	// Generate threats from components
	allThreats = append(allThreats, tmi.GenerateThreatsFromComponents(arch)...)

	// Generate threats from relationships
	allThreats = append(allThreats, tmi.GenerateThreatsFromRelationships(arch)...)

	// Generate threats from assumptions
	allThreats = append(allThreats, tmi.GenerateThreatsFromAssumptions(assumptions, arch)...)

	// Generate threats from trust boundaries
	allThreats = append(allThreats, tmi.GenerateThreatsFromBoundaries(boundaries, arch)...)

	// Deduplicate by name
	allThreats = deduplicateThreats(allThreats)

	// Renumber IDs
	for i := range allThreats {
		allThreats[i].ID = fmt.Sprintf("THREAT-%03d", i+1)
	}

	// Cluster threats
	clusters := tmi.ClusterThreats(allThreats)

	// Build summary
	summary := tmi.BuildSummary(allThreats, clusters)

	return &TMIRunResult{
		Threats:  allThreats,
		Clusters: clusters,
		Summary:  summary,
	}
}

// deduplicateThreats removes duplicate threats by name.
func deduplicateThreats(threats []Threat) []Threat {
	seen := make(map[string]bool)
	var result []Threat
	for _, t := range threats {
		key := t.Name + "-" + strings.Join(t.AffectedComponents, ",")
		if !seen[key] {
			seen[key] = true
			result = append(result, t)
		}
	}
	return result
}
