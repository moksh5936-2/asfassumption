package fidelity

import (
	"asf-tui/asf/fact"
	"fmt"
	"strings"
)

// HiddenAssumption represents an assumption that is NOT an explicit fact.
// It is something that must be true for the architecture to remain secure,
// but is not explicitly stated in the architecture.
type HiddenAssumption struct {
	ID             string   `json:"id"`
	Description    string   `json:"description"`
	ComponentID    string   `json:"component_id,omitempty"`
	ComponentLabel string   `json:"component_label,omitempty"`
	Category       string   `json:"category"`
	Risk           string   `json:"risk"`
	Confidence     float64  `json:"confidence"`
	SourceType     string   `json:"source_type"` // "fact-derived", "relationship-derived", "trust-boundary-derived", "domain-derived"
	SourceFactID   string   `json:"source_fact_id,omitempty"`
	SourceFactText string   `json:"source_fact_text,omitempty"`
	Reason         string   `json:"reason"`
	Keywords       []string `json:"keywords,omitempty"`
	NoveltyScore   float64  `json:"novelty_score"`
	RelevanceScore float64  `json:"relevance_score"`
	QualityScore   float64  `json:"quality_score"`
}

// HiddenAssumptionEngine generates hidden assumptions.
// It is the core of the ASF v2.2.0+ architectural fidelity engine.
type HiddenAssumptionEngine struct {
	factProtection *fact.ProtectionLayer
	domainPack     string
}

// NewHiddenAssumptionEngine creates a new hidden assumption engine.
func NewHiddenAssumptionEngine(inventory *fact.Inventory, domain string) *HiddenAssumptionEngine {
	return &HiddenAssumptionEngine{
		factProtection: fact.NewProtectionLayer(inventory),
		domainPack:     domain,
	}
}

// Generate generates hidden assumptions from the given facts and components.
func (e *HiddenAssumptionEngine) Generate(facts *fact.Inventory, components []Component, relationships []Relationship) []HiddenAssumption {
	var assumptions []HiddenAssumption

	// Generate from facts: for each fact, what is the hidden assumption?
	assumptions = append(assumptions, e.generateFromFacts(facts)...)

	// Generate from components: for each component, what are the hidden assumptions?
	assumptions = append(assumptions, e.generateFromComponents(components, facts)...)

	// Generate from relationships
	assumptions = append(assumptions, e.generateFromRelationships(relationships, components, facts)...)

	// Generate from domain pack
	assumptions = append(assumptions, e.generateFromDomainPack(facts, components)...)

	// Apply fact protection
	assumptions = e.applyFactProtection(assumptions)

	// Score and filter
	assumptions = e.scoreAndFilter(assumptions)

	return assumptions
}

// generateFromFacts creates hidden assumptions from explicit facts.
// For each fact, we infer what is NOT stated but must be true.
func (e *HiddenAssumptionEngine) generateFromFacts(facts *fact.Inventory) []HiddenAssumption {
	var assumptions []HiddenAssumption

	for _, f := range facts.Facts {
		// If the fact is a control, compliance, or requirement, what hidden assumption is implied?
		if f.FactType == "control" || f.FactType == "configuration" || f.FactType == "compliance" || f.FactType == "requirement" {
			hidden := e.inferHiddenFromFact(f)
			assumptions = append(assumptions, hidden...)
		}
	}

	return assumptions
}

// inferHiddenFromFact infers hidden assumptions from a single fact.
func (e *HiddenAssumptionEngine) inferHiddenFromFact(f fact.Fact) []HiddenAssumption {
	var assumptions []HiddenAssumption
	lower := strings.ToLower(f.Text)

	// MFA enabled -> hidden assumption: admin accounts are restricted
	if strings.Contains(lower, "mfa") && strings.Contains(lower, "enabled") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-mfa-admin", f.ID),
			Description:    "Administrative accounts are restricted to authorized personnel only",
			Category:       "authentication",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "MFA is enabled but administrative account access is not explicitly specified",
			Keywords:       []string{"mfa", "admin", "authentication"},
		})
	}

	// MFA disabled -> hidden assumption: compensating controls exist
	if strings.Contains(lower, "mfa") && strings.Contains(lower, "disabled") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-mfa-comp", f.ID),
			Description:    "Compensating controls exist for MFA-disabled accounts (e.g., IP restrictions, hardware tokens, or policy-based monitoring)",
			Category:       "authentication",
			Risk:           "critical",
			Confidence:     0.7,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "MFA is disabled; compensating controls are needed to maintain security posture",
			Keywords:       []string{"mfa", "compensating", "authentication"},
		})
	}

	// Encryption enabled -> hidden assumption: key management and rotation
	if strings.Contains(lower, "encryption") && strings.Contains(lower, "enabled") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-enc-key", f.ID),
			Description:    "Encryption keys are managed with secure lifecycle including rotation and access controls",
			Category:       "cryptography",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Encryption is enabled but key management practices are not explicitly specified",
			Keywords:       []string{"encryption", "key management", "rotation"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-enc-cert", f.ID),
			Description:    "Certificates are monitored for expiration and rotated before expiry",
			Category:       "cryptography",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Encryption uses certificates but renewal lifecycle is not specified",
			Keywords:       []string{"encryption", "certificate", "rotation"},
		})
	}

	// Encryption disabled -> hidden assumption: compensating controls
	if strings.Contains(lower, "encryption") && strings.Contains(lower, "disabled") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-enc-comp", f.ID),
			Description:    "Data is protected by alternative controls (e.g., physical access, network segmentation, or application-level obfuscation) to compensate for disabled encryption",
			Category:       "cryptography",
			Risk:           "critical",
			Confidence:     0.6,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Encryption is disabled; alternative controls must be explicitly validated",
			Keywords:       []string{"encryption", "compensating", "data protection"},
		})
	}

	// Auth0 used -> hidden assumption: Auth0 admin restricted
	if strings.Contains(lower, "auth0") || strings.Contains(lower, "auth") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-auth-admin", f.ID),
			Description:    "Authentication provider administrators are restricted to authorized personnel and have MFA on their own accounts",
			Category:       "authentication",
			Risk:           "high",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Authentication provider is used but administrator access is not explicitly specified",
			Keywords:       []string{"auth", "admin", "mfa"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-auth-session", f.ID),
			Description:    "Session tokens have short expiration and are invalidated on logout",
			Category:       "authentication",
			Risk:           "high",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Authentication is configured but session management practices are not explicitly specified",
			Keywords:       []string{"auth", "session", "token"},
		})
	}

	// Backups exist -> hidden assumption: restore testing
	if strings.Contains(lower, "backup") && (strings.Contains(lower, "enabled") || strings.Contains(lower, "configured") || strings.Contains(lower, "automated")) {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-backup-test", f.ID),
			Description:    "Backup restore procedures are tested periodically and recovery time is validated",
			Category:       "availability",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Backups exist but restore validation is not explicitly specified",
			Keywords:       []string{"backup", "restore", "testing"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-backup-encrypt", f.ID),
			Description:    "Backup data is encrypted at rest and in transit to backup storage",
			Category:       "cryptography",
			Risk:           "high",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Backups exist but backup encryption is not explicitly specified",
			Keywords:       []string{"backup", "encryption", "data protection"},
		})
	}

	// HIPAA required -> hidden assumption: PHI access controls
	if strings.Contains(lower, "hipaa") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-hipaa-phi", f.ID),
			Description:    "PHI access is logged and audited with immutable audit trails",
			Category:       "compliance",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "HIPAA compliance requires audit logging but is not explicitly stated",
			Keywords:       []string{"hipaa", "phi", "audit", "logging"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-hipaa-breach", f.ID),
			Description:    "Breach notification procedures are documented and tested within 60-day SLA",
			Category:       "compliance",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "HIPAA requires breach notification procedures but they are not explicitly stated",
			Keywords:       []string{"hipaa", "breach", "notification"},
		})
	}

	// WAF enabled -> hidden assumption: rule management
	if strings.Contains(lower, "waf") || strings.Contains(lower, "firewall") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-waf-rules", f.ID),
			Description:    "WAF rules are reviewed and updated regularly to address new threat patterns",
			Category:       "network",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "WAF is enabled but rule management practices are not explicitly specified",
			Keywords:       []string{"waf", "firewall", "rules", "maintenance"},
		})
	}

	// Logging enabled -> hidden assumption: log retention
	if strings.Contains(lower, "log") || strings.Contains(lower, "audit") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-log-retention", f.ID),
			Description:    "Log retention policies meet compliance requirements and are not prematurely deleted",
			Category:       "compliance",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Logging is enabled but retention policies are not explicitly specified",
			Keywords:       []string{"log", "retention", "compliance"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-log-monitor", f.ID),
			Description:    "Logs are monitored for anomalies and security events with automated alerting",
			Category:       "monitoring",
			Risk:           "high",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Logging is enabled but log monitoring is not explicitly specified",
			Keywords:       []string{"log", "monitoring", "alerting"},
		})
	}

	// VPN used -> hidden assumption: VPN access management
	if strings.Contains(lower, "vpn") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-vpn-access", f.ID),
			Description:    "VPN access is restricted by role and requires MFA for all users",
			Category:       "network",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "VPN is used but access management practices are not explicitly specified",
			Keywords:       []string{"vpn", "access", "mfa"},
		})
	}

	// RBAC / Role-based access control -> hidden assumptions
	if strings.Contains(lower, "rbac") || strings.Contains(lower, "role-based") || strings.Contains(lower, "role based") || strings.Contains(lower, "least privilege") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-rbac-review", f.ID),
			Description:    "RBAC roles and permissions are reviewed periodically to remove unused or excessive privileges",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "RBAC is enforced but periodic review is not explicitly specified",
			Keywords:       []string{"rbac", "review", "permissions"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-rbac-service", f.ID),
			Description:    "Service accounts have minimal permissions and are not shared across applications",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "RBAC is enforced but service account practices are not explicitly specified",
			Keywords:       []string{"rbac", "service account", "permissions"},
		})
	}

	// Rate limiting / Throttling / API quota -> hidden assumptions
	if strings.Contains(lower, "rate limit") || strings.Contains(lower, "rate-limit") || strings.Contains(lower, "throttling") || strings.Contains(lower, "quota") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-rate-tune", f.ID),
			Description:    "Rate limiting thresholds are tuned per endpoint and reviewed for DDoS resilience",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Rate limiting is enabled but threshold tuning is not explicitly specified",
			Keywords:       []string{"rate limiting", "threshold", "tuning"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-rate-alert", f.ID),
			Description:    "Rate limit violations are logged and alerted for potential abuse or attacks",
			Category:       "monitoring",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Rate limiting is enabled but violation monitoring is not explicitly specified",
			Keywords:       []string{"rate limiting", "alert", "abuse"},
		})
	}

	// Monitoring / Alerting / Observability -> hidden assumptions
	if strings.Contains(lower, "monitor") || strings.Contains(lower, "alert") || strings.Contains(lower, "observability") || strings.Contains(lower, "cloudwatch") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-mon-config", f.ID),
			Description:    "Monitoring dashboards are configured for security and availability metrics with defined SLAs",
			Category:       "monitoring",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Monitoring is enabled but dashboard configuration is not explicitly specified",
			Keywords:       []string{"monitoring", "dashboard", "sla"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-mon-runbook", f.ID),
			Description:    "Runbooks are documented for responding to critical alerts and incident escalation paths are defined",
			Category:       "monitoring",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Monitoring is enabled but incident response procedures are not explicitly specified",
			Keywords:       []string{"monitoring", "runbook", "incident response"},
		})
	}

	// Penetration testing / Security testing -> hidden assumptions
	if strings.Contains(lower, "penetration") || strings.Contains(lower, "pentest") || strings.Contains(lower, "security test") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-pentest-remediate", f.ID),
			Description:    "Penetration testing findings are tracked to closure with documented remediation timelines and SLA",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Penetration testing is performed but remediation tracking is not explicitly specified",
			Keywords:       []string{"penetration testing", "remediation", "sla"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-pentest-scope", f.ID),
			Description:    "Penetration testing scope covers all external interfaces and critical internal components",
			Category:       "security",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Penetration testing is performed but scope coverage is not explicitly specified",
			Keywords:       []string{"penetration testing", "scope", "coverage"},
		})
	}

	// DLP / Data Loss Prevention -> hidden assumptions
	if strings.Contains(lower, "dlp") || strings.Contains(lower, "data loss") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-dlp-policy", f.ID),
			Description:    "DLP policies are updated to cover new data types, regulatory requirements, and emerging threats",
			Category:       "compliance",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "DLP is enabled but policy update process is not explicitly specified",
			Keywords:       []string{"dlp", "policy", "update"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-dlp-except", f.ID),
			Description:    "DLP exceptions are documented, approved, and reviewed periodically for business justification",
			Category:       "compliance",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "DLP is enabled but exception management is not explicitly specified",
			Keywords:       []string{"dlp", "exception", "review"},
		})
	}

	// Tenant isolation / Multi-tenancy -> hidden assumptions
	if strings.Contains(lower, "tenant isolation") || strings.Contains(lower, "multi-tenant") || strings.Contains(lower, "tenant data") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-tenant-test", f.ID),
			Description:    "Tenant isolation is tested periodically and cross-tenant access is validated as impossible",
			Category:       "security",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Tenant isolation is enforced but testing is not explicitly specified",
			Keywords:       []string{"tenant", "isolation", "testing"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-tenant-key", f.ID),
			Description:    "Tenant data is encrypted with tenant-specific keys and key rotation is managed per tenant",
			Category:       "security",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Tenant isolation is enforced but key management per tenant is not explicitly specified",
			Keywords:       []string{"tenant", "encryption", "key management"},
		})
	}

	// Data retention / Retention policy -> hidden assumptions
	if strings.Contains(lower, "data retention") || strings.Contains(lower, "retention policy") || strings.Contains(lower, "retention") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-retention-auto", f.ID),
			Description:    "Data retention policies are enforced automatically and expired data is purged without manual intervention",
			Category:       "compliance",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Data retention is enforced but automation is not explicitly specified",
			Keywords:       []string{"retention", "automation", "purging"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-retention-legal", f.ID),
			Description:    "Legal hold procedures are documented for suspending retention deletion during litigation or investigation",
			Category:       "compliance",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Data retention is enforced but legal hold process is not explicitly specified",
			Keywords:       []string{"retention", "legal hold", "litigation"},
		})
	}

	// API Gateway / API Management -> hidden assumptions
	if strings.Contains(lower, "api gateway") || strings.Contains(lower, "api management") || strings.Contains(lower, "api rate") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-api-auth", f.ID),
			Description:    "API gateway enforces authentication for all endpoints including internal APIs and health checks",
			Category:       "security",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "API gateway is used but authentication enforcement is not explicitly specified",
			Keywords:       []string{"api gateway", "authentication", "endpoints"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-api-version", f.ID),
			Description:    "API versioning is maintained and deprecated versions are sunset with customer notice",
			Category:       "security",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "API gateway is used but versioning policy is not explicitly specified",
			Keywords:       []string{"api", "versioning", "deprecation"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-api-log", f.ID),
			Description:    "API requests are logged with correlation IDs for tracing and security analysis",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "API gateway is used but request logging is not explicitly specified",
			Keywords:       []string{"api", "logging", "correlation"},
		})
	}

	// GuardDuty / Threat detection -> hidden assumptions
	if strings.Contains(lower, "guardduty") || strings.Contains(lower, "guard duty") || strings.Contains(lower, "threat detection") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-guardduty-remediate", f.ID),
			Description:    "GuardDuty findings are remediated automatically or within documented SLA with ticket tracking",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "GuardDuty is enabled but remediation process is not explicitly specified",
			Keywords:       []string{"guardduty", "remediation", "sla"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-guardduty-false", f.ID),
			Description:    "GuardDuty false positives are tuned to reduce alert fatigue and maintain signal quality",
			Category:       "security",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "GuardDuty is enabled but false positive tuning is not explicitly specified",
			Keywords:       []string{"guardduty", "false positive", "tuning"},
		})
	}

	// CloudTrail / Audit trail -> hidden assumptions
	if strings.Contains(lower, "cloudtrail") || strings.Contains(lower, "cloud trail") || strings.Contains(lower, "audit trail") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-trail-analyze", f.ID),
			Description:    "CloudTrail logs are analyzed for unauthorized API calls, anomalous patterns, and insider threats",
			Category:       "monitoring",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "CloudTrail is enabled but log analysis is not explicitly specified",
			Keywords:       []string{"cloudtrail", "analysis", "api calls"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-trail-integrity", f.ID),
			Description:    "CloudTrail logs are protected from tampering with log file validation and S3 object lock",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "CloudTrail is enabled but integrity protection is not explicitly specified",
			Keywords:       []string{"cloudtrail", "integrity", "tampering"},
		})
	}

	// Security groups / Firewall rules -> hidden assumptions
	if strings.Contains(lower, "security group") || strings.Contains(lower, "firewall") || strings.Contains(lower, "acl") || strings.Contains(lower, "access control list") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-sg-review", f.ID),
			Description:    "Security group rules are reviewed periodically and unused or overly permissive rules are removed",
			Category:       "network",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Security groups are restricted but review process is not explicitly specified",
			Keywords:       []string{"security group", "review", "rules"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-sg-ingress", f.ID),
			Description:    "Ingress rules are restricted to minimum required ports and sources with documented business justification",
			Category:       "network",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Security groups are restricted but ingress justification is not explicitly specified",
			Keywords:       []string{"security group", "ingress", "justification"},
		})
	}

	// S3 / Object storage -> hidden assumptions
	if strings.Contains(lower, "s3") || strings.Contains(lower, "bucket") || strings.Contains(lower, "object storage") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-s3-access", f.ID),
			Description:    "S3 bucket access policies are reviewed and public access is blocked at account and bucket level",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "S3 is encrypted but access policy review is not explicitly specified",
			Keywords:       []string{"s3", "access policy", "public access"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-s3-lifecycle", f.ID),
			Description:    "S3 lifecycle policies are configured for transitioning old data to cheaper storage and eventual deletion",
			Category:       "availability",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "S3 is used but lifecycle management is not explicitly specified",
			Keywords:       []string{"s3", "lifecycle", "storage"},
		})
	}

	// RDS / Database -> hidden assumptions
	if strings.Contains(lower, "rds") || strings.Contains(lower, "database") || strings.Contains(lower, "db") || strings.Contains(lower, "sql") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-db-patch", f.ID),
			Description:    "Database is patched regularly and maintenance windows are configured for minimal downtime",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Database has backups but patching schedule is not explicitly specified",
			Keywords:       []string{"database", "patch", "maintenance"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-db-replica", f.ID),
			Description:    "Database has read replicas for performance and standby for failover with tested switchover",
			Category:       "availability",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Database is used but replication strategy is not explicitly specified",
			Keywords:       []string{"database", "replica", "failover"},
		})
	}

	// AWS Config / Compliance monitoring -> hidden assumptions
	if strings.Contains(lower, "aws config") || strings.Contains(lower, "config") && strings.Contains(lower, "compliance") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-config-remediate", f.ID),
			Description:    "AWS Config rules are remediated automatically and non-compliant resources are flagged for review",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "AWS Config is enabled but remediation is not explicitly specified",
			Keywords:       []string{"config", "remediation", "compliance"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-config-custom", f.ID),
			Description:    "Custom AWS Config rules are defined for organization-specific security requirements",
			Category:       "security",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "AWS Config is enabled but custom rules are not explicitly specified",
			Keywords:       []string{"config", "custom rules", "organization"},
		})
	}

	// Container image scanning / Vulnerability scanning -> hidden assumptions
	if strings.Contains(lower, "container image") || strings.Contains(lower, "image scan") || strings.Contains(lower, "vulnerability scan") || strings.Contains(lower, "scan") && strings.Contains(lower, "container") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-scan-sla", f.ID),
			Description:    "Vulnerability findings are remediated within documented SLA and tracked to closure",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Container images are scanned but remediation SLA is not explicitly specified",
			Keywords:       []string{"vulnerability", "sla", "remediation"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-scan-registry", f.ID),
			Description:    "Container registry has image signing and verification to prevent tampered deployments",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Container images are scanned but registry signing is not explicitly specified",
			Keywords:       []string{"container", "registry", "signing"},
		})
	}

	// Resource quotas / Limits -> hidden assumptions
	if strings.Contains(lower, "resource quota") || strings.Contains(lower, "quota") || strings.Contains(lower, "limit") && strings.Contains(lower, "resource") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-quota-enforce", f.ID),
			Description:    "Resource quotas are enforced for all namespaces and unrestricted namespaces are not allowed",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Resource quotas are set but enforcement scope is not explicitly specified",
			Keywords:       []string{"resource quota", "namespace", "enforcement"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-quota-monitor", f.ID),
			Description:    "Resource usage is monitored against quotas and alerts are configured for approaching limits",
			Category:       "monitoring",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Resource quotas are set but usage monitoring is not explicitly specified",
			Keywords:       []string{"resource quota", "monitoring", "limits"},
		})
	}

	// Auto-scaling / Scaling -> hidden assumptions
	if strings.Contains(lower, "auto-scaling") || strings.Contains(lower, "auto scaling") || strings.Contains(lower, "scaling") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-scale-limit", f.ID),
			Description:    "Auto-scaling limits are configured to prevent resource exhaustion and cost overruns",
			Category:       "availability",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Auto-scaling is enabled but limits are not explicitly specified",
			Keywords:       []string{"auto-scaling", "limits", "cost"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-scale-health", f.ID),
			Description:    "Auto-scaling health checks are configured to ensure only healthy instances are created",
			Category:       "availability",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Auto-scaling is enabled but health check configuration is not explicitly specified",
			Keywords:       []string{"auto-scaling", "health check", "instances"},
		})
	}

	// Node auto-scaling / Node scaling -> hidden assumptions
	if strings.Contains(lower, "node auto-scaling") || strings.Contains(lower, "node scaling") || strings.Contains(lower, "node") && strings.Contains(lower, "auto") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-node-patch", f.ID),
			Description:    "Nodes are patched automatically and unsupported versions are drained and replaced",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Node auto-scaling is enabled but patching process is not explicitly specified",
			Keywords:       []string{"node", "patch", "auto-scaling"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-node-cordon", f.ID),
			Description:    "Node cordon and drain procedures are tested for maintenance and incident response",
			Category:       "availability",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Node auto-scaling is enabled but maintenance procedures are not explicitly specified",
			Keywords:       []string{"node", "cordon", "drain"},
		})
	}

	// Cluster logging / Logging -> hidden assumptions
	if strings.Contains(lower, "cluster logging") || strings.Contains(lower, "logging") && strings.Contains(lower, "cluster") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-cluster-log-forward", f.ID),
			Description:    "Cluster logs are forwarded to a centralized SIEM for correlation and long-term retention",
			Category:       "monitoring",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Cluster logging is enabled but SIEM forwarding is not explicitly specified",
			Keywords:       []string{"cluster logging", "siem", "forwarding"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-cluster-log-retention", f.ID),
			Description:    "Cluster log retention policies meet compliance requirements and are not prematurely truncated",
			Category:       "compliance",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Cluster logging is enabled but retention policies are not explicitly specified",
			Keywords:       []string{"cluster logging", "retention", "compliance"},
		})
	}

	// Cluster monitoring / Monitoring -> hidden assumptions
	if strings.Contains(lower, "cluster monitoring") || strings.Contains(lower, "monitoring") && strings.Contains(lower, "cluster") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-cluster-mon-alert", f.ID),
			Description:    "Cluster monitoring alerts are configured for node health, pod restarts, and resource pressure",
			Category:       "monitoring",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Cluster monitoring is enabled but alert configuration is not explicitly specified",
			Keywords:       []string{"cluster monitoring", "alert", "health"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-cluster-mon-slo", f.ID),
			Description:    "SLOs and SLIs are defined for cluster services and monitored for breach detection",
			Category:       "monitoring",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Cluster monitoring is enabled but SLO definitions are not explicitly specified",
			Keywords:       []string{"cluster monitoring", "slo", "sli"},
		})
	}

	// VPC / Network segmentation -> hidden assumptions
	if strings.Contains(lower, "vpc") || strings.Contains(lower, "network segmentation") || strings.Contains(lower, "segmented") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-vpc-peer", f.ID),
			Description:    "VPC peering connections are restricted and reviewed for unnecessary exposure",
			Category:       "network",
			Risk:           "high",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "VPC is configured but peering connections are not explicitly specified",
			Keywords:       []string{"vpc", "peering", "exposure"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-vpc-flow", f.ID),
			Description:    "VPC flow logs are enabled and monitored for anomalies",
			Category:       "network",
			Risk:           "high",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "VPC is configured but flow logging is not explicitly specified",
			Keywords:       []string{"vpc", "flow logs", "anomalies"},
		})
	}

	// IAM -> hidden assumptions
	if strings.Contains(lower, "iam") || strings.Contains(lower, "identity") && strings.Contains(lower, "access") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-iam-review", f.ID),
			Description:    "IAM policies are reviewed periodically and unused permissions are removed",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "IAM is configured but review process is not explicitly specified",
			Keywords:       []string{"iam", "review", "permissions"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-iam-credential", f.ID),
			Description:    "IAM credentials are rotated and access keys are not hardcoded in applications",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "IAM is configured but credential rotation is not explicitly specified",
			Keywords:       []string{"iam", "credential", "rotation"},
		})
	}

	// KMS / Key management -> hidden assumptions
	if strings.Contains(lower, "kms") || strings.Contains(lower, "key management") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-kms-rotate", f.ID),
			Description:    "KMS keys are rotated regularly and key rotation is logged for audit",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "KMS is used but key rotation schedule is not explicitly specified",
			Keywords:       []string{"kms", "rotation", "audit"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-kms-access", f.ID),
			Description:    "KMS key access is restricted to authorized roles and access is logged for audit",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "KMS is used but access control is not explicitly specified",
			Keywords:       []string{"kms", "access control", "audit"},
		})
	}

	// CDN / CloudFront -> hidden assumptions
	if strings.Contains(lower, "cdn") || strings.Contains(lower, "cloudfront") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-cdn-ssl", f.ID),
			Description:    "CDN enforces HTTPS and does not serve content over HTTP",
			Category:       "network",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "CDN is used but HTTPS enforcement is not explicitly specified",
			Keywords:       []string{"cdn", "https", "ssl"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-cdn-cache", f.ID),
			Description:    "CDN cache invalidation is configured for security updates and content changes",
			Category:       "infrastructure",
			Risk:           "medium",
			Confidence:     0.75,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "CDN is used but cache invalidation is not explicitly specified",
			Keywords:       []string{"cdn", "cache", "invalidation"},
		})
	}

	// DDoS / Denial of Service -> hidden assumptions
	if strings.Contains(lower, "ddos") || strings.Contains(lower, "denial of service") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-ddos-tune", f.ID),
			Description:    "DDoS protection thresholds are tuned and tested regularly for effectiveness",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "DDoS protection is enabled but tuning is not explicitly specified",
			Keywords:       []string{"ddos", "tuning", "testing"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-ddos-response", f.ID),
			Description:    "DDoS incident response procedures are documented and team roles are assigned",
			Category:       "security",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "DDoS protection is enabled but incident response is not explicitly specified",
			Keywords:       []string{"ddos", "incident response", "team"},
		})
	}

	// Fraud detection -> hidden assumptions
	if strings.Contains(lower, "fraud") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-fraud-update", f.ID),
			Description:    "Fraud detection rules are updated based on new threat patterns and incident data",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Fraud detection is enabled but rule updates are not explicitly specified",
			Keywords:       []string{"fraud", "rules", "update"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-fraud-false", f.ID),
			Description:    "Fraud false positives are tuned to reduce customer friction while maintaining detection accuracy",
			Category:       "security",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Fraud detection is enabled but false positive tuning is not explicitly specified",
			Keywords:       []string{"fraud", "false positive", "tuning"},
		})
	}

	// Tokenization / Token vault -> hidden assumptions
	if strings.Contains(lower, "tokenization") || strings.Contains(lower, "token vault") || strings.Contains(lower, "token") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-token-vault", f.ID),
			Description:    "Token vault is secured with access controls and audit logging for all token operations",
			Category:       "security",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Tokenization is enabled but vault security is not explicitly specified",
			Keywords:       []string{"tokenization", "vault", "access control"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-token-rotate", f.ID),
			Description:    "Token vault keys are rotated and token mapping is backed up for disaster recovery",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Tokenization is enabled but key rotation is not explicitly specified",
			Keywords:       []string{"tokenization", "rotation", "backup"},
		})
	}

	// PCI DSS -> hidden assumptions
	if strings.Contains(lower, "pci") || strings.Contains(lower, "payment card") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-pci-segment", f.ID),
			Description:    "Cardholder data environment is segmented and access is restricted to authorized personnel",
			Category:       "compliance",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "PCI DSS is required but CDE segmentation is not explicitly specified",
			Keywords:       []string{"pci", "cde", "segmentation"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-pci-qsa", f.ID),
			Description:    "QSA audits are performed annually and findings are remediated within documented timelines",
			Category:       "compliance",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "PCI DSS is required but QSA audit process is not explicitly specified",
			Keywords:       []string{"pci", "qsa", "audit"},
		})
	}

	// SOC2 -> hidden assumptions
	if strings.Contains(lower, "soc2") || strings.Contains(lower, "soc 2") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-soc2-evidence", f.ID),
			Description:    "SOC2 evidence is collected and maintained for auditor review with documented retention",
			Category:       "compliance",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "SOC2 is required but evidence collection is not explicitly specified",
			Keywords:       []string{"soc2", "evidence", "audit"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-soc2-trust", f.ID),
			Description:    "SOC2 trust service criteria are mapped to controls and tested for compliance",
			Category:       "compliance",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "SOC2 is required but trust criteria mapping is not explicitly specified",
			Keywords:       []string{"soc2", "trust", "criteria"},
		})
	}

	// GDPR -> hidden assumptions
	if strings.Contains(lower, "gdpr") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-gdpr-dsar", f.ID),
			Description:    "Data subject access requests are processed within 30-day SLA with documented workflow",
			Category:       "compliance",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "GDPR is required but DSAR workflow is not explicitly specified",
			Keywords:       []string{"gdpr", "dsar", "sla"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-gdpr-dpo", f.ID),
			Description:    "Data Protection Officer is designated and contact information is published",
			Category:       "compliance",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "GDPR is required but DPO designation is not explicitly specified",
			Keywords:       []string{"gdpr", "dpo", "contact"},
		})
	}

	// ISO27001 -> hidden assumptions
	if strings.Contains(lower, "iso27001") || strings.Contains(lower, "iso 27001") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-iso-isms", f.ID),
			Description:    "ISMS is reviewed and maintained annually with management commitment and resources",
			Category:       "compliance",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "ISO27001 is required but ISMS review is not explicitly specified",
			Keywords:       []string{"iso27001", "isms", "review"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-iso-risk", f.ID),
			Description:    "Information security risk assessments are performed and risk treatment plans are implemented",
			Category:       "compliance",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "ISO27001 is required but risk assessment process is not explicitly specified",
			Keywords:       []string{"iso27001", "risk", "assessment"},
		})
	}

	// NIST -> hidden assumptions
	if strings.Contains(lower, "nist") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-nist-csf", f.ID),
			Description:    "NIST Cybersecurity Framework is mapped to controls and maturity is assessed regularly",
			Category:       "compliance",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "NIST is required but CSF mapping is not explicitly specified",
			Keywords:       []string{"nist", "csf", "maturity"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-nist-800", f.ID),
			Description:    "NIST SP 800-53 controls are implemented and assessed for effectiveness",
			Category:       "compliance",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "NIST is required but SP 800-53 implementation is not explicitly specified",
			Keywords:       []string{"nist", "800-53", "controls"},
		})
	}

	// Split tunneling / VPN client -> hidden assumptions
	if strings.Contains(lower, "split tunnel") || strings.Contains(lower, "split-tunnel") || strings.Contains(lower, "vpn client") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-vpn-split", f.ID),
			Description:    "Split tunneling is disabled and all traffic is routed through VPN for security inspection",
			Category:       "network",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Split tunneling is disabled but traffic inspection is not explicitly specified",
			Keywords:       []string{"split tunnel", "traffic inspection", "vpn"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-vpn-client-update", f.ID),
			Description:    "VPN client software is updated regularly and unsupported versions are blocked",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "VPN client is managed but update process is not explicitly specified",
			Keywords:       []string{"vpn client", "update", "version"},
		})
	}

	// Certificate / Certificate rotation -> hidden assumptions
	if strings.Contains(lower, "certificate") || strings.Contains(lower, "cert rotation") || strings.Contains(lower, "cert") && strings.Contains(lower, "rotate") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-cert-revoke", f.ID),
			Description:    "Certificate revocation lists are maintained and revoked certificates are not accepted",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Certificates are rotated but revocation process is not explicitly specified",
			Keywords:       []string{"certificate", "revocation", "crl"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-cert-auto", f.ID),
			Description:    "Certificate renewal is automated and expiration monitoring is configured with alerts",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Certificates are rotated but automation is not explicitly specified",
			Keywords:       []string{"certificate", "renewal", "automation"},
		})
	}

	// Firewall / Strict firewall -> hidden assumptions
	if strings.Contains(lower, "firewall") && strings.Contains(lower, "strict") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-fw-review", f.ID),
			Description:    "Firewall rules are reviewed periodically and unused rules are removed",
			Category:       "network",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Firewall rules are strict but review process is not explicitly specified",
			Keywords:       []string{"firewall", "rules", "review"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-fw-log", f.ID),
			Description:    "Firewall traffic is logged and analyzed for anomalies and unauthorized access attempts",
			Category:       "monitoring",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Firewall rules are strict but logging is not explicitly specified",
			Keywords:       []string{"firewall", "logging", "anomalies"},
		})
	}

	// Ingress Controller / Ingress -> hidden assumptions
	if strings.Contains(lower, "ingress") || strings.Contains(lower, "ingress controller") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-ingress-tls", f.ID),
			Description:    "Ingress controller terminates TLS with proper certificate management and cipher configuration",
			Category:       "network",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Ingress controller is used but TLS management is not explicitly specified",
			Keywords:       []string{"ingress", "tls", "certificate"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-ingress-restrict", f.ID),
			Description:    "Ingress routes are restricted to necessary paths and default backends are disabled",
			Category:       "network",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Ingress controller is used but route restriction is not explicitly specified",
			Keywords:       []string{"ingress", "routes", "restriction"},
		})
	}

	// Service Mesh -> hidden assumptions
	if strings.Contains(lower, "service mesh") || strings.Contains(lower, "istio") || strings.Contains(lower, "linkerd") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-mesh-mtls", f.ID),
			Description:    "Service mesh enforces mutual TLS for all service-to-service communication",
			Category:       "network",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Service mesh is used but mTLS enforcement is not explicitly specified",
			Keywords:       []string{"service mesh", "mtls", "communication"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-mesh-policy", f.ID),
			Description:    "Service mesh policies are reviewed and updated for new services and traffic patterns",
			Category:       "security",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Service mesh is used but policy review is not explicitly specified",
			Keywords:       []string{"service mesh", "policy", "review"},
		})
	}

	// CI/CD Pipeline -> hidden assumptions
	if strings.Contains(lower, "ci/cd") || strings.Contains(lower, "pipeline") || strings.Contains(lower, "cicd") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-cicd-gate", f.ID),
			Description:    "CI/CD pipeline has security gates for SAST, DAST, and dependency scanning before deployment",
			Category:       "security",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "CI/CD pipeline is used but security gates are not explicitly specified",
			Keywords:       []string{"cicd", "security gates", "sast"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-cicd-secret", f.ID),
			Description:    "CI/CD pipeline does not contain secrets in configuration and uses secret management tools",
			Category:       "security",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "CI/CD pipeline is used but secret management is not explicitly specified",
			Keywords:       []string{"cicd", "secrets", "configuration"},
		})
	}

	// etcd / Key-value store -> hidden assumptions
	if strings.Contains(lower, "etcd") || strings.Contains(lower, "key-value") || strings.Contains(lower, "kv store") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-etcd-backup", f.ID),
			Description:    "etcd is backed up regularly and restore procedures are tested for cluster recovery",
			Category:       "availability",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "etcd is used but backup strategy is not explicitly specified",
			Keywords:       []string{"etcd", "backup", "recovery"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-etcd-encrypt", f.ID),
			Description:    "etcd data is encrypted at rest and backup encryption is configured",
			Category:       "security",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "etcd is used but encryption is not explicitly specified",
			Keywords:       []string{"etcd", "encryption", "rest"},
		})
	}

	// Lambda / Serverless -> hidden assumptions
	if strings.Contains(lower, "lambda") || strings.Contains(lower, "serverless") || strings.Contains(lower, "function") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-lambda-perm", f.ID),
			Description:    "Lambda functions have minimal IAM permissions and execution roles are not overly permissive",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Lambda is used but permission minimization is not explicitly specified",
			Keywords:       []string{"lambda", "permissions", "iam"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-lambda-env", f.ID),
			Description:    "Lambda environment variables do not contain secrets and use AWS Secrets Manager or Parameter Store",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Lambda is used but secret management is not explicitly specified",
			Keywords:       []string{"lambda", "environment", "secrets"},
		})
	}

	// DynamoDB / NoSQL -> hidden assumptions
	if strings.Contains(lower, "dynamodb") || strings.Contains(lower, "nosql") || strings.Contains(lower, "document") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-dynamo-backup", f.ID),
			Description:    "DynamoDB has point-in-time recovery enabled and backups are tested for restore",
			Category:       "availability",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "DynamoDB is used but backup strategy is not explicitly specified",
			Keywords:       []string{"dynamodb", "backup", "recovery"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-dynamo-encrypt", f.ID),
			Description:    "DynamoDB is encrypted at rest and CMK is used for sensitive tables",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "DynamoDB is used but encryption configuration is not explicitly specified",
			Keywords:       []string{"dynamodb", "encryption", "cmk"},
		})
	}

	// EC2 / Compute -> hidden assumptions
	if strings.Contains(lower, "ec2") || strings.Contains(lower, "compute") || strings.Contains(lower, "instance") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-ec2-patch", f.ID),
			Description:    "EC2 instances are patched regularly and unsupported AMIs are not used",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "EC2 is used but patching schedule is not explicitly specified",
			Keywords:       []string{"ec2", "patch", "ami"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-ec2-profile", f.ID),
			Description:    "EC2 instance profiles have minimal permissions and are not shared across instances",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "EC2 is used but instance profile management is not explicitly specified",
			Keywords:       []string{"ec2", "instance profile", "permissions"},
		})
	}

	// EKS / Kubernetes cluster -> hidden assumptions
	if strings.Contains(lower, "eks") || strings.Contains(lower, "kubernetes cluster") || strings.Contains(lower, "k8s cluster") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-eks-version", f.ID),
			Description:    "EKS cluster is kept on a supported Kubernetes version and upgrades are planned",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "EKS is used but version management is not explicitly specified",
			Keywords:       []string{"eks", "version", "upgrade"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-eks-endpoint", f.ID),
			Description:    "EKS cluster endpoint is private or restricted to authorized CIDRs and public access is disabled",
			Category:       "network",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "EKS is used but endpoint access control is not explicitly specified",
			Keywords:       []string{"eks", "endpoint", "access"},
		})
	}

	// Generic catch-all for any control that doesn't match above
	// If the fact is a positive control but didn't match any specific rule above,
	// generate a generic hidden assumption about lifecycle management
	if !f.IsNegative && f.FactType == "control" && len(assumptions) == 0 {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-generic-lifecycle", f.ID),
			Description:    fmt.Sprintf("%s is reviewed periodically and updated to maintain effectiveness against emerging threats", f.Text),
			Category:       "security",
			Risk:           "medium",
			Confidence:     0.7,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Control is enabled but lifecycle management is not explicitly specified",
			Keywords:       []string{"control", "review", "lifecycle"},
		})
	}

	// VPC/Network -> hidden assumptions
	if strings.Contains(lower, "vpc") || strings.Contains(lower, "network") || strings.Contains(lower, "segment") || strings.Contains(lower, "private") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-vpc-peer", f.ID),
			Description:    "VPC peering connections are restricted and reviewed for unnecessary exposure",
			Category:       "network",
			Risk:           "high",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "VPC is configured but peering connections are not explicitly specified",
			Keywords:       []string{"vpc", "network", "peering"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-vpc-flow", f.ID),
			Description:    "VPC flow logs are enabled and monitored for anomalies",
			Category:       "network",
			Risk:           "high",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "VPC is configured but flow logging is not explicitly specified",
			Keywords:       []string{"vpc", "flow logs", "monitoring"},
		})
	}

	// IAM -> hidden assumptions
	if strings.Contains(lower, "iam") || strings.Contains(lower, "role") || strings.Contains(lower, "privilege") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-iam-review", f.ID),
			Description:    "IAM policy changes are reviewed before deployment and regularly audited",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "IAM is configured but policy review process is not explicitly specified",
			Keywords:       []string{"iam", "policy", "review", "audit"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-iam-credential", f.ID),
			Description:    "IAM credentials are rotated and access keys are not hardcoded in applications",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "IAM is configured but credential rotation is not explicitly specified",
			Keywords:       []string{"iam", "credential", "rotation", "access key"},
		})
	}

	// CloudWatch/Monitoring -> hidden assumptions
	if strings.Contains(lower, "monitor") || strings.Contains(lower, "cloudwatch") || strings.Contains(lower, "alert") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-mon-alert", f.ID),
			Description:    "Monitoring alerts are configured for critical security and availability metrics",
			Category:       "monitoring",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Monitoring is enabled but alert configuration is not explicitly specified",
			Keywords:       []string{"monitoring", "alert", "metrics"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-mon-runbook", f.ID),
			Description:    "Runbooks are documented for responding to monitoring alerts and incidents",
			Category:       "monitoring",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Monitoring is enabled but incident response procedures are not explicitly specified",
			Keywords:       []string{"monitoring", "runbook", "incident response"},
		})
	}

	// Kubernetes-specific -> hidden assumptions
	if strings.Contains(lower, "rbac") || strings.Contains(lower, "role") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-k8s-rbac-review", f.ID),
			Description:    "RBAC roles and bindings are reviewed periodically to remove unused permissions",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "RBAC is enforced but periodic review is not explicitly specified",
			Keywords:       []string{"rbac", "review", "permissions"},
		})
	}
	if strings.Contains(lower, "network policy") || strings.Contains(lower, "network policies") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-k8s-net-policy", f.ID),
			Description:    "Network policies are enforced for all namespaces including default and kube-system",
			Category:       "network",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Network policies are configured but namespace coverage is not explicitly specified",
			Keywords:       []string{"network policy", "namespace", "coverage"},
		})
	}
	if strings.Contains(lower, "pod security") || strings.Contains(lower, "admission") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-k8s-admission-review", f.ID),
			Description:    "Admission controller policies are reviewed and updated for new Kubernetes versions",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Admission controllers are enabled but policy review is not explicitly specified",
			Keywords:       []string{"admission", "policy", "review"},
		})
	}
	if strings.Contains(lower, "secret") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-k8s-secret-mount", f.ID),
			Description:    "Secrets are mounted only to pods that need them and are not stored in environment variables",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Secrets are encrypted but mounting practices are not explicitly specified",
			Keywords:       []string{"secret", "mount", "environment variable"},
		})
	}
	if strings.Contains(lower, "container image") || strings.Contains(lower, "scan") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-k8s-vuln-sla", f.ID),
			Description:    "Vulnerability findings are remediated within documented SLA and tracked to closure",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Container images are scanned but remediation SLA is not explicitly specified",
			Keywords:       []string{"vulnerability", "sla", "remediation"},
		})
	}
	if strings.Contains(lower, "resource quota") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-k8s-quota-enforce", f.ID),
			Description:    "Resource quotas are enforced for all namespaces and namespaces are not left unrestricted",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Resource quotas are set but enforcement scope is not explicitly specified",
			Keywords:       []string{"resource quota", "namespace", "enforcement"},
		})
	}
	if strings.Contains(lower, "auto-scaling") || strings.Contains(lower, "auto scaling") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-k8s-scale-limit", f.ID),
			Description:    "Auto-scaling limits are configured to prevent resource exhaustion and cost overruns",
			Category:       "availability",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Auto-scaling is enabled but limits are not explicitly specified",
			Keywords:       []string{"auto-scaling", "limits", "cost"},
		})
	}

	// Cloud-specific -> hidden assumptions
	if strings.Contains(lower, "guardduty") || strings.Contains(lower, "threat detection") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-cloud-guardduty-remediate", f.ID),
			Description:    "GuardDuty findings are remediated automatically or within documented SLA",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "GuardDuty is enabled but remediation process is not explicitly specified",
			Keywords:       []string{"guardduty", "remediation", "sla"},
		})
	}
	if strings.Contains(lower, "cloudtrail") || strings.Contains(lower, "cloud trail") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-cloud-trail-monitor", f.ID),
			Description:    "CloudTrail logs are monitored for unauthorized API calls and anomalous access patterns",
			Category:       "monitoring",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "CloudTrail is enabled but log analysis is not explicitly specified",
			Keywords:       []string{"cloudtrail", "monitoring", "api calls"},
		})
	}
	if strings.Contains(lower, "security group") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-cloud-sg-review", f.ID),
			Description:    "Security group rules are reviewed periodically and unused rules are removed",
			Category:       "network",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Security groups are restricted but review process is not explicitly specified",
			Keywords:       []string{"security group", "review", "rules"},
		})
	}
	if strings.Contains(lower, "s3") || strings.Contains(lower, "bucket") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-cloud-s3-access", f.ID),
			Description:    "S3 bucket access policies are reviewed and public access is blocked at account level",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "S3 is encrypted but access policy review is not explicitly specified",
			Keywords:       []string{"s3", "access policy", "public access"},
		})
	}
	if strings.Contains(lower, "rds") || strings.Contains(lower, "database") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-cloud-rds-patch", f.ID),
			Description:    "RDS database is patched regularly and maintenance windows are configured",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "RDS has backups but patching schedule is not explicitly specified",
			Keywords:       []string{"rds", "patch", "maintenance"},
		})
	}
	if strings.Contains(lower, "config") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-cloud-config-remediate", f.ID),
			Description:    "AWS Config rules are remediated automatically and non-compliant resources are flagged",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "AWS Config is enabled but remediation is not explicitly specified",
			Keywords:       []string{"config", "remediation", "compliance"},
		})
	}

	// SaaS-specific -> hidden assumptions
	if strings.Contains(lower, "tenant") || strings.Contains(lower, "isolation") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-saas-tenant-test", f.ID),
			Description:    "Tenant isolation is tested periodically and cross-tenant access is not possible",
			Category:       "security",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Tenant isolation is enforced but testing is not explicitly specified",
			Keywords:       []string{"tenant", "isolation", "testing"},
		})
	}
	if strings.Contains(lower, "dlp") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-saas-dlp-update", f.ID),
			Description:    "DLP policies are updated to cover new data types and regulatory requirements",
			Category:       "compliance",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "DLP is enabled but policy update process is not explicitly specified",
			Keywords:       []string{"dlp", "policy", "update"},
		})
	}
	if strings.Contains(lower, "data retention") || strings.Contains(lower, "retention") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-saas-retention-enforce", f.ID),
			Description:    "Data retention policies are enforced automatically and expired data is purged",
			Category:       "compliance",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Data retention policies are enforced but automatic purging is not explicitly specified",
			Keywords:       []string{"retention", "purging", "automation"},
		})
	}
	if strings.Contains(lower, "penetration") || strings.Contains(lower, "pentest") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-saas-pentest-remediate", f.ID),
			Description:    "Penetration testing findings are remediated and tracked to closure with documented timelines",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Penetration testing is performed but remediation tracking is not explicitly specified",
			Keywords:       []string{"penetration testing", "remediation", "tracking"},
		})
	}

	// VPN-specific -> hidden assumptions
	if strings.Contains(lower, "split tunnel") || strings.Contains(lower, "split-tunnel") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-vpn-split", f.ID),
			Description:    "Split tunneling is disabled and all traffic is routed through VPN for security inspection",
			Category:       "network",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Split tunneling is disabled but traffic inspection is not explicitly specified",
			Keywords:       []string{"split tunnel", "traffic inspection", "vpn"},
		})
	}
	if strings.Contains(lower, "certificate") || strings.Contains(lower, "cert") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-vpn-cert-revoke", f.ID),
			Description:    "Certificate revocation lists are maintained and revoked certificates are not accepted",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Certificates are rotated but revocation process is not explicitly specified",
			Keywords:       []string{"certificate", "revocation", "crl"},
		})
	}
	if strings.Contains(lower, "firewall") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-vpn-fw-review", f.ID),
			Description:    "Firewall rules are reviewed periodically and unused rules are removed",
			Category:       "network",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Firewall rules are strict but review process is not explicitly specified",
			Keywords:       []string{"firewall", "rules", "review"},
		})
	}
	if strings.Contains(lower, "vpn client") || strings.Contains(lower, "client") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-vpn-client-update", f.ID),
			Description:    "VPN client software is updated regularly and unsupported versions are blocked",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "VPN client is managed but update process is not explicitly specified",
			Keywords:       []string{"vpn client", "update", "version"},
		})
	}

	// PCI DSS -> hidden assumptions
	if strings.Contains(lower, "pci") || strings.Contains(lower, "payment") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-pci-segment", f.ID),
			Description:    "Cardholder data environment is segmented and access is restricted to authorized personnel",
			Category:       "compliance",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "PCI DSS is required but CDE segmentation is not explicitly specified",
			Keywords:       []string{"pci", "cde", "segmentation"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-pi-qsa", f.ID),
			Description:    "QSA (Qualified Security Assessor) audits are performed annually and findings are remediated",
			Category:       "compliance",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "PCI DSS is required but QSA audit process is not explicitly specified",
			Keywords:       []string{"pci", "qsa", "audit"},
		})
	}

	// SOC2 -> hidden assumptions
	if strings.Contains(lower, "soc2") || strings.Contains(lower, "soc 2") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-soc2-evidence", f.ID),
			Description:    "SOC2 evidence is collected and maintained for auditor review with documented retention",
			Category:       "compliance",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "SOC2 is required but evidence collection is not explicitly specified",
			Keywords:       []string{"soc2", "evidence", "audit"},
		})
	}

	// GDPR -> hidden assumptions
	if strings.Contains(lower, "gdpr") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-gdpr-dsar", f.ID),
			Description:    "Data subject access requests (DSAR) are processed within 30-day SLA",
			Category:       "compliance",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "GDPR is required but DSAR process is not explicitly specified",
			Keywords:       []string{"gdpr", "dsar", "sla"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-gdpr-dpo", f.ID),
			Description:    "Data Protection Officer (DPO) is designated and contact information is published",
			Category:       "compliance",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "GDPR is required but DPO designation is not explicitly specified",
			Keywords:       []string{"gdpr", "dpo", "contact"},
		})
	}

	// ISO27001 -> hidden assumptions
	if strings.Contains(lower, "iso27001") || strings.Contains(lower, "iso 27001") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-iso-isms", f.ID),
			Description:    "ISMS (Information Security Management System) is reviewed and maintained annually",
			Category:       "compliance",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "ISO27001 is required but ISMS review is not explicitly specified",
			Keywords:       []string{"iso27001", "isms", "review"},
		})
	}

	// NIST -> hidden assumptions
	if strings.Contains(lower, "nist") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-nist-csf", f.ID),
			Description:    "NIST Cybersecurity Framework is mapped to controls and maturity is assessed regularly",
			Category:       "compliance",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "NIST is required but CSF mapping is not explicitly specified",
			Keywords:       []string{"nist", "csf", "maturity"},
		})
	}

	// Fraud detection -> hidden assumptions
	if strings.Contains(lower, "fraud") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-fraud-update", f.ID),
			Description:    "Fraud detection rules are updated regularly based on new threat patterns and incident data",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Fraud detection is enabled but rule update process is not explicitly specified",
			Keywords:       []string{"fraud", "rules", "update"},
		})
	}

	// Tokenization -> hidden assumptions
	if strings.Contains(lower, "tokenization") || strings.Contains(lower, "token") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-token-vault", f.ID),
			Description:    "Token vault is secured with access controls and audit logging for all token operations",
			Category:       "security",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "Tokenization is enabled but vault security is not explicitly specified",
			Keywords:       []string{"tokenization", "vault", "access control"},
		})
	}

	// DDoS -> hidden assumptions
	if strings.Contains(lower, "ddos") || strings.Contains(lower, "rate limit") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-ddos-tune", f.ID),
			Description:    "DDoS protection thresholds are tuned and tested regularly to ensure effectiveness",
			Category:       "security",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "DDoS protection is enabled but tuning process is not explicitly specified",
			Keywords:       []string{"ddos", "tuning", "testing"},
		})
	}

	// CDN -> hidden assumptions
	if strings.Contains(lower, "cdn") || strings.Contains(lower, "cloudfront") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-cdn-ssl", f.ID),
			Description:    "CDN enforces HTTPS and does not serve content over HTTP",
			Category:       "network",
			Risk:           "high",
			Confidence:     0.85,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "CDN is used but HTTPS enforcement is not explicitly specified",
			Keywords:       []string{"cdn", "https", "ssl"},
		})
	}

	// API gateway -> hidden assumptions
	if strings.Contains(lower, "api gateway") || strings.Contains(lower, "api gateway") {
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-api-auth", f.ID),
			Description:    "API gateway enforces authentication for all endpoints including internal APIs",
			Category:       "security",
			Risk:           "critical",
			Confidence:     0.9,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "API gateway is used but authentication enforcement is not explicitly specified",
			Keywords:       []string{"api gateway", "authentication", "endpoints"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:             fmt.Sprintf("hid-%s-api-version", f.ID),
			Description:    "API versioning is maintained and deprecated versions are sunset with notice",
			Category:       "security",
			Risk:           "medium",
			Confidence:     0.8,
			SourceType:     "fact-derived",
			SourceFactID:   f.ID,
			SourceFactText: f.Text,
			Reason:         "API gateway is used but versioning policy is not explicitly specified",
			Keywords:       []string{"api", "versioning", "deprecation"},
		})
	}

	return assumptions
}

// Component represents a component for the hidden assumption engine.
type Component struct {
	ID    string
	Label string
}

// Relationship represents a relationship for the hidden assumption engine.
type Relationship struct {
	Source string
	Target string
	Label  string
}

// generateFromComponents generates hidden assumptions from component presence.
func (e *HiddenAssumptionEngine) generateFromComponents(components []Component, facts *fact.Inventory) []HiddenAssumption {
	var assumptions []HiddenAssumption

	for _, comp := range components {
		lower := strings.ToLower(comp.Label)

		// Database component
		if strings.Contains(lower, "database") || strings.Contains(lower, "db") || strings.Contains(lower, "data") {
			assumptions = append(assumptions, HiddenAssumption{
				ID:             fmt.Sprintf("hid-comp-%s-db-replica", comp.ID),
				Description:    "Database has replication and failover for high availability",
				ComponentID:    comp.ID,
				ComponentLabel: comp.Label,
				Category:       "availability",
				Risk:           "high",
				Confidence:     0.7,
				SourceType:     "component-derived",
				Reason:         "Database component exists but replication and failover are not explicitly specified",
				Keywords:       []string{"database", "replication", "failover"},
			})
			assumptions = append(assumptions, HiddenAssumption{
				ID:             fmt.Sprintf("hid-comp-%s-db-backup", comp.ID),
				Description:    "Database has automated backups with tested restore procedures",
				ComponentID:    comp.ID,
				ComponentLabel: comp.Label,
				Category:       "availability",
				Risk:           "high",
				Confidence:     0.7,
				SourceType:     "component-derived",
				Reason:         "Database component exists but backup strategy is not explicitly specified",
				Keywords:       []string{"database", "backup", "restore"},
			})
		}

		// API component
		if strings.Contains(lower, "api") || strings.Contains(lower, "gateway") {
			assumptions = append(assumptions, HiddenAssumption{
				ID:             fmt.Sprintf("hid-comp-%s-api-rate", comp.ID),
				Description:    "API has rate limiting and throttling configured to prevent abuse",
				ComponentID:    comp.ID,
				ComponentLabel: comp.Label,
				Category:       "api",
				Risk:           "high",
				Confidence:     0.8,
				SourceType:     "component-derived",
				Reason:         "API component exists but rate limiting is not explicitly specified",
				Keywords:       []string{"api", "rate limiting", "throttling"},
			})
			assumptions = append(assumptions, HiddenAssumption{
				ID:             fmt.Sprintf("hid-comp-%s-api-auth", comp.ID),
				Description:    "API requires strong authentication for all endpoints",
				ComponentID:    comp.ID,
				ComponentLabel: comp.Label,
				Category:       "api",
				Risk:           "high",
				Confidence:     0.8,
				SourceType:     "component-derived",
				Reason:         "API component exists but authentication requirements are not explicitly specified",
				Keywords:       []string{"api", "authentication", "authorization"},
			})
		}

		// Load balancer
		if strings.Contains(lower, "load balancer") || strings.Contains(lower, "lb") {
			assumptions = append(assumptions, HiddenAssumption{
				ID:             fmt.Sprintf("hid-comp-%s-lb-health", comp.ID),
				Description:    "Load balancer has health checks configured to detect and remove unhealthy instances",
				ComponentID:    comp.ID,
				ComponentLabel: comp.Label,
				Category:       "infrastructure",
				Risk:           "high",
				Confidence:     0.85,
				SourceType:     "component-derived",
				Reason:         "Load balancer exists but health check configuration is not explicitly specified",
				Keywords:       []string{"load balancer", "health check", "monitoring"},
			})
		}

		// CDN
		if strings.Contains(lower, "cdn") {
			assumptions = append(assumptions, HiddenAssumption{
				ID:             fmt.Sprintf("hid-comp-%s-cdn-cache", comp.ID),
				Description:    "CDN cache invalidation is configured for security updates and content changes",
				ComponentID:    comp.ID,
				ComponentLabel: comp.Label,
				Category:       "infrastructure",
				Risk:           "medium",
				Confidence:     0.75,
				SourceType:     "component-derived",
				Reason:         "CDN is used but cache invalidation strategy is not explicitly specified",
				Keywords:       []string{"cdn", "cache", "invalidation"},
			})
		}

		// Message queue
		if strings.Contains(lower, "queue") || strings.Contains(lower, "kafka") || strings.Contains(lower, "rabbitmq") {
			assumptions = append(assumptions, HiddenAssumption{
				ID:             fmt.Sprintf("hid-comp-%s-queue-dlq", comp.ID),
				Description:    "Message queue has dead letter queue configured for failed messages",
				ComponentID:    comp.ID,
				ComponentLabel: comp.Label,
				Category:       "availability",
				Risk:           "medium",
				Confidence:     0.75,
				SourceType:     "component-derived",
				Reason:         "Message queue exists but dead letter queue configuration is not explicitly specified",
				Keywords:       []string{"queue", "dead letter", "availability"},
			})
		}

		// Cache
		if strings.Contains(lower, "cache") || strings.Contains(lower, "redis") || strings.Contains(lower, "memcached") {
			assumptions = append(assumptions, HiddenAssumption{
				ID:             fmt.Sprintf("hid-comp-%s-cache-purge", comp.ID),
				Description:    "Cache has secure eviction policies to prevent sensitive data leakage",
				ComponentID:    comp.ID,
				ComponentLabel: comp.Label,
				Category:       "data",
				Risk:           "medium",
				Confidence:     0.75,
				SourceType:     "component-derived",
				Reason:         "Cache component exists but eviction policies are not explicitly specified",
				Keywords:       []string{"cache", "eviction", "data protection"},
			})
		}

		// Object storage
		if strings.Contains(lower, "s3") || strings.Contains(lower, "blob") || strings.Contains(lower, "storage") {
			assumptions = append(assumptions, HiddenAssumption{
				ID:             fmt.Sprintf("hid-comp-%s-storage-version", comp.ID),
				Description:    "Object storage has versioning enabled for data integrity and recovery",
				ComponentID:    comp.ID,
				ComponentLabel: comp.Label,
				Category:       "availability",
				Risk:           "medium",
				Confidence:     0.75,
				SourceType:     "component-derived",
				Reason:         "Object storage exists but versioning is not explicitly specified",
				Keywords:       []string{"storage", "versioning", "data integrity"},
			})
		}
	}

	return assumptions
}

// generateFromRelationships generates hidden assumptions from relationships.
func (e *HiddenAssumptionEngine) generateFromRelationships(relationships []Relationship, components []Component, facts *fact.Inventory) []HiddenAssumption {
	var assumptions []HiddenAssumption

	for _, rel := range relationships {
		// Database -> Application
		if (strings.Contains(strings.ToLower(rel.Target), "database") || strings.Contains(strings.ToLower(rel.Target), "db")) &&
			(strings.Contains(strings.ToLower(rel.Source), "app") || strings.Contains(strings.ToLower(rel.Source), "service") || strings.Contains(strings.ToLower(rel.Source), "api")) {
			assumptions = append(assumptions, HiddenAssumption{
				ID:          fmt.Sprintf("hid-rel-%s-%s-db-conn", rel.Source, rel.Target),
				Description: "Database connections use parameterized queries or ORM to prevent injection",
				Category:    "data",
				Risk:        "critical",
				Confidence:  0.85,
				SourceType:  "relationship-derived",
				Reason:      "Application connects to database but query safety is not explicitly specified",
				Keywords:    []string{"database", "injection", "query safety"},
			})
		}

		// Internet -> Application
		if strings.Contains(strings.ToLower(rel.Source), "internet") || strings.Contains(strings.ToLower(rel.Source), "web") || strings.Contains(strings.ToLower(rel.Source), "public") {
			assumptions = append(assumptions, HiddenAssumption{
				ID:          fmt.Sprintf("hid-rel-%s-%s-public-ddos", rel.Source, rel.Target),
				Description: "Public-facing endpoints have DDoS protection and rate limiting configured",
				Category:    "network",
				Risk:        "high",
				Confidence:  0.85,
				SourceType:  "relationship-derived",
				Reason:      "Public-facing relationship exists but DDoS protection is not explicitly specified",
				Keywords:    []string{"public", "ddos", "rate limiting"},
			})
		}

		// API -> Database
		if (strings.Contains(strings.ToLower(rel.Source), "api") || strings.Contains(strings.ToLower(rel.Source), "gateway")) &&
			(strings.Contains(strings.ToLower(rel.Target), "database") || strings.Contains(strings.ToLower(rel.Target), "db")) {
			assumptions = append(assumptions, HiddenAssumption{
				ID:          fmt.Sprintf("hid-rel-%s-%s-api-db-auth", rel.Source, rel.Target),
				Description: "API-to-database connections authenticate with least-privilege credentials and rotate regularly",
				Category:    "authentication",
				Risk:        "critical",
				Confidence:  0.85,
				SourceType:  "relationship-derived",
				Reason:      "API connects to database but connection authentication is not explicitly specified",
				Keywords:    []string{"api", "database", "authentication", "least privilege"},
			})
		}

		// Third-party -> Application
		if strings.Contains(strings.ToLower(rel.Source), "third") || strings.Contains(strings.ToLower(rel.Source), "external") || strings.Contains(strings.ToLower(rel.Source), "vendor") {
			assumptions = append(assumptions, HiddenAssumption{
				ID:          fmt.Sprintf("hid-rel-%s-%s-thirdparty-contract", rel.Source, rel.Target),
				Description: "Third-party integrations have security requirements in contracts and SLA monitoring",
				Category:    "governance",
				Risk:        "high",
				Confidence:  0.8,
				SourceType:  "relationship-derived",
				Reason:      "Third-party connection exists but security contract terms are not explicitly specified",
				Keywords:    []string{"third-party", "contract", "sla", "governance"},
			})
		}

		// Admin -> Database
		if (strings.Contains(strings.ToLower(rel.Source), "admin") || strings.Contains(strings.ToLower(rel.Source), "management")) &&
			(strings.Contains(strings.ToLower(rel.Target), "database") || strings.Contains(strings.ToLower(rel.Target), "db")) {
			assumptions = append(assumptions, HiddenAssumption{
				ID:          fmt.Sprintf("hid-rel-%s-%s-admin-db", rel.Source, rel.Target),
				Description: "Administrative access to database uses dedicated accounts with MFA and session logging",
				Category:    "authentication",
				Risk:        "critical",
				Confidence:  0.9,
				SourceType:  "relationship-derived",
				Reason:      "Admin access to database exists but access controls are not explicitly specified",
				Keywords:    []string{"admin", "database", "mfa", "access control"},
			})
		}

		// Load balancer -> Application
		if (strings.Contains(strings.ToLower(rel.Source), "load balancer") || strings.Contains(strings.ToLower(rel.Source), "lb")) &&
			(strings.Contains(strings.ToLower(rel.Target), "app") || strings.Contains(strings.ToLower(rel.Target), "service")) {
			assumptions = append(assumptions, HiddenAssumption{
				ID:          fmt.Sprintf("hid-rel-%s-%s-lb-ssl", rel.Source, rel.Target),
				Description: "Load balancer terminates SSL/TLS with proper certificate management and cipher suite configuration",
				Category:    "network",
				Risk:        "high",
				Confidence:  0.85,
				SourceType:  "relationship-derived",
				Reason:      "Load balancer connects to application but SSL termination details are not explicitly specified",
				Keywords:    []string{"load balancer", "ssl", "tls", "certificate"},
			})
		}
	}

	return assumptions
}

// generateFromDomainPack generates domain-specific hidden assumptions.
func (e *HiddenAssumptionEngine) generateFromDomainPack(facts *fact.Inventory, components []Component) []HiddenAssumption {
	var assumptions []HiddenAssumption

	// Healthcare domain
	if e.domainPack == "healthcare" {
		assumptions = append(assumptions, HiddenAssumption{
			ID:          "hid-domain-healthcare-breakglass",
			Description: "Break-glass procedures are documented and tested for emergency PHI access",
			Category:    "compliance",
			Risk:        "critical",
			Confidence:  0.9,
			SourceType:  "domain-derived",
			Reason:      "Healthcare domain requires break-glass procedures for emergency access",
			Keywords:    []string{"healthcare", "break-glass", "phi", "emergency"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:          "hid-domain-healthcare-clinical",
			Description: "Clinical data access is logged with user identity and timestamp for audit trails",
			Category:    "compliance",
			Risk:        "critical",
			Confidence:  0.9,
			SourceType:  "domain-derived",
			Reason:      "Healthcare domain requires clinical data access logging",
			Keywords:    []string{"healthcare", "clinical", "logging", "audit"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:          "hid-domain-healthcare-patient",
			Description: "Patient safety mechanisms are in place for data integrity and availability during incidents",
			Category:    "compliance",
			Risk:        "critical",
			Confidence:  0.9,
			SourceType:  "domain-derived",
			Reason:      "Healthcare domain requires patient safety mechanisms",
			Keywords:    []string{"healthcare", "patient safety", "data integrity"},
		})
	}

	// Fintech domain
	if e.domainPack == "fintech" || e.domainPack == "payment" {
		assumptions = append(assumptions, HiddenAssumption{
			ID:          "hid-domain-fintech-fraud",
			Description: "Fraud detection and prevention mechanisms are configured for payment processing",
			Category:    "security",
			Risk:        "critical",
			Confidence:  0.9,
			SourceType:  "domain-derived",
			Reason:      "Fintech domain requires fraud detection for payment processing",
			Keywords:    []string{"fintech", "fraud", "payment", "detection"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:          "hid-domain-fintech-aml",
			Description: "AML (Anti-Money Laundering) and KYC (Know Your Customer) procedures are enforced for all transactions",
			Category:    "compliance",
			Risk:        "critical",
			Confidence:  0.9,
			SourceType:  "domain-derived",
			Reason:      "Fintech domain requires AML/KYC compliance",
			Keywords:    []string{"fintech", "aml", "kyc", "compliance"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:          "hid-domain-fintech-settlement",
			Description: "Settlement reconciliation processes are automated and auditable for financial accuracy",
			Category:    "compliance",
			Risk:        "high",
			Confidence:  0.85,
			SourceType:  "domain-derived",
			Reason:      "Fintech domain requires settlement reconciliation",
			Keywords:    []string{"fintech", "settlement", "reconciliation", "audit"},
		})
	}

	// Kubernetes domain
	if e.domainPack == "kubernetes" || e.domainPack == "k8s" {
		assumptions = append(assumptions, HiddenAssumption{
			ID:          "hid-domain-k8s-rbac",
			Description: "Kubernetes RBAC is enforced with least-privilege roles and service accounts are restricted",
			Category:    "security",
			Risk:        "critical",
			Confidence:  0.9,
			SourceType:  "domain-derived",
			Reason:      "Kubernetes domain requires RBAC enforcement",
			Keywords:    []string{"kubernetes", "rbac", "least privilege"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:          "hid-domain-k8s-admission",
			Description: "Admission controllers are configured to enforce pod security policies and resource limits",
			Category:    "security",
			Risk:        "critical",
			Confidence:  0.9,
			SourceType:  "domain-derived",
			Reason:      "Kubernetes domain requires admission control",
			Keywords:    []string{"kubernetes", "admission", "pod security"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:          "hid-domain-k8s-secrets",
			Description: "Kubernetes secrets are encrypted at rest and never stored in environment variables or config maps",
			Category:    "security",
			Risk:        "critical",
			Confidence:  0.9,
			SourceType:  "domain-derived",
			Reason:      "Kubernetes domain requires secure secret management",
			Keywords:    []string{"kubernetes", "secrets", "encryption", "security"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:          "hid-domain-k8s-network",
			Description: "Network policies are enforced to restrict pod-to-pod communication and prevent lateral movement",
			Category:    "network",
			Risk:        "high",
			Confidence:  0.9,
			SourceType:  "domain-derived",
			Reason:      "Kubernetes domain requires network policies",
			Keywords:    []string{"kubernetes", "network policy", "lateral movement"},
		})
	}

	// Cloud domain
	if e.domainPack == "cloud" || e.domainPack == "aws" || e.domainPack == "azure" || e.domainPack == "gcp" {
		assumptions = append(assumptions, HiddenAssumption{
			ID:          "hid-domain-cloud-iam",
			Description: "Cloud IAM policies are enforced with least-privilege access and MFA for all privileged roles",
			Category:    "security",
			Risk:        "critical",
			Confidence:  0.9,
			SourceType:  "domain-derived",
			Reason:      "Cloud domain requires IAM enforcement",
			Keywords:    []string{"cloud", "iam", "least privilege", "mfa"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:          "hid-domain-cloud-kms",
			Description: "Cloud KMS keys are rotated regularly and access is logged for audit trails",
			Category:    "security",
			Risk:        "high",
			Confidence:  0.9,
			SourceType:  "domain-derived",
			Reason:      "Cloud domain requires KMS key management",
			Keywords:    []string{"cloud", "kms", "key rotation", "audit"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:          "hid-domain-cloud-federation",
			Description: "Identity federation is configured for cross-account and cross-organization access with conditional access policies",
			Category:    "security",
			Risk:        "high",
			Confidence:  0.85,
			SourceType:  "domain-derived",
			Reason:      "Cloud domain requires identity federation",
			Keywords:    []string{"cloud", "identity federation", "conditional access"},
		})
		assumptions = append(assumptions, HiddenAssumption{
			ID:          "hid-domain-cloud-logging",
			Description: "Cloud resource logs are centralized and monitored for security anomalies with automated alerting",
			Category:    "monitoring",
			Risk:        "high",
			Confidence:  0.9,
			SourceType:  "domain-derived",
			Reason:      "Cloud domain requires centralized logging",
			Keywords:    []string{"cloud", "logging", "monitoring", "alerting"},
		})
	}

	return assumptions
}

// applyFactProtection filters out assumptions that contradict facts.
func (e *HiddenAssumptionEngine) applyFactProtection(assumptions []HiddenAssumption) []HiddenAssumption {
	var filtered []HiddenAssumption
	for _, a := range assumptions {
		result := e.factProtection.CheckAssumption(a.Description)
		if result.Allowed {
			filtered = append(filtered, a)
		}
		// If suppressed, we could log it but don't include it
	}
	return filtered
}

// scoreAndFilter scores assumptions and filters out low-quality ones.
func (e *HiddenAssumptionEngine) scoreAndFilter(assumptions []HiddenAssumption) []HiddenAssumption {
	var scored []HiddenAssumption
	for _, a := range assumptions {
		// Compute quality score
		a.NoveltyScore = e.computeNovelty(a)
		a.RelevanceScore = e.computeRelevance(a)
		a.QualityScore = e.computeQuality(a)

		// Filter out generic assumptions
		if a.QualityScore >= 0.5 {
			scored = append(scored, a)
		}
	}
	return scored
}

// computeNovelty computes how novel (not generic) the assumption is.
func (e *HiddenAssumptionEngine) computeNovelty(a HiddenAssumption) float64 {
	genericPhrases := []string{
		"use encryption", "use mfa", "use tls", "use ssl", "use firewall",
		"use waf", "use vpn", "use backup", "use logging", "use audit",
		"enable encryption", "enable mfa", "enable tls", "enable ssl",
		"implement encryption", "implement mfa", "implement tls",
	}

	lower := strings.ToLower(a.Description)
	for _, phrase := range genericPhrases {
		if strings.Contains(lower, phrase) {
			return 0.1
		}
	}

	// If it's specific (mentions a component, has specific detail)
	if a.ComponentID != "" {
		return 0.9
	}

	// If it has specific technical detail
	if len(a.Description) > 50 {
		return 0.8
	}

	return 0.6
}

// computeRelevance computes how relevant the assumption is.
func (e *HiddenAssumptionEngine) computeRelevance(a HiddenAssumption) float64 {
	score := 0.5

	// Domain-derived are highly relevant
	if a.SourceType == "domain-derived" {
		score += 0.3
	}

	// Fact-derived are highly relevant
	if a.SourceType == "fact-derived" {
		score += 0.3
	}

	// Component-specific is relevant
	if a.ComponentID != "" {
		score += 0.2
	}

	// High risk is more relevant
	if a.Risk == "critical" {
		score += 0.2
	} else if a.Risk == "high" {
		score += 0.1
	}

	return min(score, 1.0)
}

// computeQuality computes overall quality.
func (e *HiddenAssumptionEngine) computeQuality(a HiddenAssumption) float64 {
	novelty := a.NoveltyScore
	relevance := a.RelevanceScore
	confidence := a.Confidence

	// Weighted average
	quality := (novelty*0.4 + relevance*0.3 + confidence*0.3)
	return min(quality, 1.0)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}
