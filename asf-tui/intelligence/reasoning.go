package intelligence

import (
	"fmt"
	"strings"
)

// ReasoningEngine generates implicit assumptions from architecture topology.
type ReasoningEngine struct {
	arch *ArchDescription
}

// NewReasoningEngine creates a reasoning engine for the given architecture.
func NewReasoningEngine(arch *ArchDescription) *ReasoningEngine {
	return &ReasoningEngine{arch: arch}
}

// InferAssumptions generates implicit assumptions based on topology rules.
func (re *ReasoningEngine) InferAssumptions() []Assumption {
	if re.arch == nil {
		return nil
	}
	var assumptions []Assumption
	assumptions = append(assumptions, re.inferDatabaseAssumptions()...)
	assumptions = append(assumptions, re.inferIdentityProviderAssumptions()...)
	assumptions = append(assumptions, re.inferAPIGatewayAssumptions()...)
	assumptions = append(assumptions, re.inferKMSAssumptions()...)
	assumptions = append(assumptions, re.inferBackupServiceAssumptions()...)
	assumptions = append(assumptions, re.inferThirdPartyAssumptions()...)
	assumptions = append(assumptions, re.inferAdminConsoleAssumptions()...)
	assumptions = append(assumptions, re.inferAuditLogAssumptions()...)
	return assumptions
}

// inferDatabaseAssumptions: Database + PHI → key management, audit logging, object-level auth.
func (re *ReasoningEngine) inferDatabaseAssumptions() []Assumption {
	var results []Assumption
	hasDatabase := false
	hasPHI := false
	for _, comp := range re.arch.Components {
		label := strings.ToLower(comp.Label)
		if strings.Contains(label, "database") || strings.Contains(label, "db") || strings.Contains(label, "data store") {
			hasDatabase = true
		}
		if strings.Contains(label, "phi") || strings.Contains(label, "health") || strings.Contains(label, "patient") || strings.Contains(label, "medical") || strings.Contains(label, "ehr") {
			hasPHI = true
		}
	}
	raw := strings.ToLower(re.arch.RawText)
	if strings.Contains(raw, "phi") || strings.Contains(raw, "hipaa") || strings.Contains(raw, "patient") || strings.Contains(raw, "medical") {
		hasPHI = true
	}

	if hasDatabase && hasPHI {
		results = append(results, Assumption{
			ID:          "INF-DB-001",
			Description: "Database stores PHI but architecture does not specify key management controls for encrypted PHI.",
			Component:   "Database",
			Category:    "KeyManagement",
			Risk:        RiskCritical,
			Likelihood:  4,
			Impact:      5,
			Confidence:  0.85,
			Keywords:    []string{"database", "phi", "key management", "encryption"},
			Rationale:   "PHI requires encryption with dedicated key management; absence of key management controls creates regulatory and confidentiality risk.",
		})
		results = append(results, Assumption{
			ID:          "INF-DB-002",
			Description: "Database stores PHI but architecture does not specify audit logging for data access and modifications.",
			Component:   "Database",
			Category:    "Auditability",
			Risk:        RiskCritical,
			Likelihood:  4,
			Impact:      5,
			Confidence:  0.85,
			Keywords:    []string{"database", "phi", "audit log", "hipaa"},
			Rationale:   "HIPAA and similar regulations require immutable audit logging for PHI access; absence creates compliance and repudiation risk.",
		})
		results = append(results, Assumption{
			ID:          "INF-DB-003",
			Description: "Database stores PHI but architecture does not specify object-level authorization controls.",
			Component:   "Database",
			Category:    "ObjectLevelAuthorization",
			Risk:        RiskCritical,
			Likelihood:  4,
			Impact:      5,
			Confidence:  0.85,
			Keywords:    []string{"database", "phi", "object-level authorization", "bola"},
			Rationale:   "PHI access requires object-level authorization to prevent unauthorized record access; absence creates BOLA/IDOR risk.",
		})
	} else if hasDatabase {
		results = append(results, Assumption{
			ID:          "INF-DB-004",
			Description: "Database is present but architecture does not specify encryption at rest or key management controls.",
			Component:   "Database",
			Category:    "KeyManagement",
			Risk:        RiskHigh,
			Likelihood:  3,
			Impact:      4,
			Confidence:  0.75,
			Keywords:    []string{"database", "encryption at rest", "key management"},
			Rationale:   "Database without encryption at rest or key management exposes stored data to tampering and disclosure.",
		})
	}
	return results
}

// inferIdentityProviderAssumptions: Auth0/IdentityProvider → MFA, session security, token validation, admin access.
func (re *ReasoningEngine) inferIdentityProviderAssumptions() []Assumption {
	var results []Assumption
	hasIdP := false
	for _, comp := range re.arch.Components {
		label := strings.ToLower(comp.Label)
		if strings.Contains(label, "auth0") || strings.Contains(label, "identity") || strings.Contains(label, "idp") || strings.Contains(label, "sso") || strings.Contains(label, "oauth") || strings.Contains(label, "oidc") || strings.Contains(label, "saml") {
			hasIdP = true
		}
	}
	raw := strings.ToLower(re.arch.RawText)
	if strings.Contains(raw, "auth0") || strings.Contains(raw, "identity provider") || strings.Contains(raw, "idp") {
		hasIdP = true
	}

	if hasIdP {
		results = append(results, Assumption{
			ID:          "INF-IDP-001",
			Description: "Identity provider is present but architecture does not specify MFA enforcement for all user populations.",
			Component:   "IdentityProvider",
			Category:    "Authentication",
			Risk:        RiskHigh,
			Likelihood:  4,
			Impact:      4,
			Confidence:  0.80,
			Keywords:    []string{"identity provider", "mfa", "authentication"},
			Rationale:   "Identity provider without MFA enforcement creates a single point of compromise for account takeover.",
		})
		results = append(results, Assumption{
			ID:          "INF-IDP-002",
			Description: "Identity provider is present but architecture does not specify session security controls (rotation, timeout, concurrent limits).",
			Component:   "IdentityProvider",
			Category:    "SessionSecurity",
			Risk:        RiskHigh,
			Likelihood:  4,
			Impact:      4,
			Confidence:  0.80,
			Keywords:    []string{"identity provider", "session", "token rotation", "timeout"},
			Rationale:   "Session security gaps increase risk of session hijacking and credential reuse.",
		})
		results = append(results, Assumption{
			ID:          "INF-IDP-003",
			Description: "Identity provider is present but architecture does not specify token validation rules (signature, expiry, revocation).",
			Component:   "IdentityProvider",
			Category:    "Authentication",
			Risk:        RiskHigh,
			Likelihood:  4,
			Impact:      4,
			Confidence:  0.80,
			Keywords:    []string{"identity provider", "token validation", "jwt", "revocation"},
			Rationale:   "Token validation gaps allow replay attacks and use of revoked or expired credentials.",
		})
		results = append(results, Assumption{
			ID:          "INF-IDP-004",
			Description: "Identity provider is present but architecture does not specify admin access controls (MFA, break-glass, session restrictions).",
			Component:   "IdentityProvider",
			Category:    "PrivilegeManagement",
			Risk:        RiskCritical,
			Likelihood:  4,
			Impact:      5,
			Confidence:  0.85,
			Keywords:    []string{"identity provider", "admin", "break-glass", "privileged access"},
			Rationale:   "Admin access to identity provider is high-value; absence of controls creates privilege escalation and account takeover risk.",
		})
	}
	return results
}

// inferAPIGatewayAssumptions: API Gateway → rate limiting, auth validation, logging.
func (re *ReasoningEngine) inferAPIGatewayAssumptions() []Assumption {
	var results []Assumption
	hasGateway := false
	for _, comp := range re.arch.Components {
		label := strings.ToLower(comp.Label)
		if strings.Contains(label, "api gateway") || strings.Contains(label, "gateway") || strings.Contains(label, "apigw") || strings.Contains(label, "ingress") {
			hasGateway = true
		}
	}
	raw := strings.ToLower(re.arch.RawText)
	if strings.Contains(raw, "api gateway") || strings.Contains(raw, "apigw") {
		hasGateway = true
	}

	if hasGateway {
		results = append(results, Assumption{
			ID:          "INF-GW-001",
			Description: "API Gateway is present but architecture does not specify rate limiting and throttling policies.",
			Component:   "APIGateway",
			Category:    "APISecurity",
			Risk:        RiskHigh,
			Likelihood:  4,
			Impact:      4,
			Confidence:  0.80,
			Keywords:    []string{"api gateway", "rate limiting", "throttling", "dos"},
			Rationale:   "API Gateway without rate limiting is vulnerable to brute force, scraping, and denial-of-service attacks.",
		})
		results = append(results, Assumption{
			ID:          "INF-GW-002",
			Description: "API Gateway is present but architecture does not specify authentication and authorization validation at the edge.",
			Component:   "APIGateway",
			Category:    "APISecurity",
			Risk:        RiskHigh,
			Likelihood:  4,
			Impact:      4,
			Confidence:  0.80,
			Keywords:    []string{"api gateway", "auth validation", "edge security"},
			Rationale:   "API Gateway should validate authentication and authorization before forwarding requests to backend services.",
		})
		results = append(results, Assumption{
			ID:          "INF-GW-003",
			Description: "API Gateway is present but architecture does not specify request/response logging and audit trail generation.",
			Component:   "APIGateway",
			Category:    "Logging",
			Risk:        RiskMedium,
			Likelihood:  3,
			Impact:      3,
			Confidence:  0.75,
			Keywords:    []string{"api gateway", "logging", "audit trail"},
			Rationale:   "API Gateway logging is essential for incident response, forensic analysis, and compliance audit trails.",
		})
	}
	return results
}

// inferKMSAssumptions: KMS → key rotation, access restriction, backup.
func (re *ReasoningEngine) inferKMSAssumptions() []Assumption {
	var results []Assumption
	hasKMS := false
	for _, comp := range re.arch.Components {
		label := strings.ToLower(comp.Label)
		if strings.Contains(label, "kms") || strings.Contains(label, "key management") || strings.Contains(label, "hsm") || strings.Contains(label, "vault") {
			hasKMS = true
		}
	}
	raw := strings.ToLower(re.arch.RawText)
	if strings.Contains(raw, "kms") || strings.Contains(raw, "key management service") {
		hasKMS = true
	}

	if hasKMS {
		results = append(results, Assumption{
			ID:          "INF-KMS-001",
			Description: "KMS is present but architecture does not specify key rotation schedule and automation.",
			Component:   "KMS",
			Category:    "KeyManagement",
			Risk:        RiskHigh,
			Likelihood:  4,
			Impact:      4,
			Confidence:  0.80,
			Keywords:    []string{"kms", "key rotation", "automation"},
			Rationale:   "Key rotation reduces exposure from compromised keys; absence of rotation policy increases cryptographic risk.",
		})
		results = append(results, Assumption{
			ID:          "INF-KMS-002",
			Description: "KMS is present but architecture does not specify access restriction and least-privilege key policies.",
			Component:   "KMS",
			Category:    "KeyManagement",
			Risk:        RiskCritical,
			Likelihood:  4,
			Impact:      5,
			Confidence:  0.85,
			Keywords:    []string{"kms", "access restriction", "least privilege", "key policy"},
			Rationale:   "Broad access to KMS keys allows decryption by unauthorized parties; least-privilege key policies are required.",
		})
		results = append(results, Assumption{
			ID:          "INF-KMS-003",
			Description: "KMS is present but architecture does not specify key backup and disaster recovery procedures.",
			Component:   "KMS",
			Category:    "DisasterRecovery",
			Risk:        RiskHigh,
			Likelihood:  3,
			Impact:      4,
			Confidence:  0.75,
			Keywords:    []string{"kms", "backup", "disaster recovery", "key escrow"},
			Rationale:   "KMS key loss leads to permanent data loss; backup and escrow procedures are required for business continuity.",
		})
	}
	return results
}

// inferBackupServiceAssumptions: BackupService → encryption, testing, geographic distribution.
func (re *ReasoningEngine) inferBackupServiceAssumptions() []Assumption {
	var results []Assumption
	hasBackup := false
	for _, comp := range re.arch.Components {
		label := strings.ToLower(comp.Label)
		if strings.Contains(label, "backup") || strings.Contains(label, "snapshot") || strings.Contains(label, "restore") || strings.Contains(label, "archive") {
			hasBackup = true
		}
	}
	raw := strings.ToLower(re.arch.RawText)
	if strings.Contains(raw, "backup") || strings.Contains(raw, "snapshot") {
		hasBackup = true
	}

	if hasBackup {
		results = append(results, Assumption{
			ID:          "INF-BKP-001",
			Description: "Backup service is present but architecture does not specify backup encryption for all copies.",
			Component:   "BackupService",
			Category:    "Backups",
			Risk:        RiskHigh,
			Likelihood:  3,
			Impact:      4,
			Confidence:  0.75,
			Keywords:    []string{"backup", "encryption", "snapshot"},
			Rationale:   "Unencrypted backups expose sensitive data if storage media is lost or stolen.",
		})
		results = append(results, Assumption{
			ID:          "INF-BKP-002",
			Description: "Backup service is present but architecture does not specify restore testing and validation procedures.",
			Component:   "BackupService",
			Category:    "Backups",
			Risk:        RiskHigh,
			Likelihood:  4,
			Impact:      4,
			Confidence:  0.80,
			Keywords:    []string{"backup", "restore test", "validation"},
			Rationale:   "Backups without tested restore procedures may fail during recovery, leading to extended downtime.",
		})
		results = append(results, Assumption{
			ID:          "INF-BKP-003",
			Description: "Backup service is present but architecture does not specify geographic distribution and cross-region replication.",
			Component:   "BackupService",
			Category:    "Backups",
			Risk:        RiskMedium,
			Likelihood:  3,
			Impact:      3,
			Confidence:  0.70,
			Keywords:    []string{"backup", "cross-region", "geographic distribution", "replication"},
			Rationale:   "Backups in a single region are vulnerable to regional disasters; cross-region replication improves resilience.",
		})
	}
	return results
}

// inferThirdPartyAssumptions: ThirdParty → vendor risk, equivalent controls, data minimization.
func (re *ReasoningEngine) inferThirdPartyAssumptions() []Assumption {
	var results []Assumption
	hasThirdParty := false
	for _, comp := range re.arch.Components {
		label := strings.ToLower(comp.Label)
		if strings.Contains(label, "third party") || strings.Contains(label, "third-party") || strings.Contains(label, "external") || strings.Contains(label, "vendor") || strings.Contains(label, "saas") || strings.Contains(label, "analytics") || strings.Contains(label, "stripe") || strings.Contains(label, "sendgrid") {
			hasThirdParty = true
		}
	}
	raw := strings.ToLower(re.arch.RawText)
	if strings.Contains(raw, "third party") || strings.Contains(raw, "third-party") || strings.Contains(raw, "external service") || strings.Contains(raw, "vendor") {
		hasThirdParty = true
	}

	if hasThirdParty {
		results = append(results, Assumption{
			ID:          "INF-TP-001",
			Description: "Third-party integration is present but architecture does not specify vendor risk assessment and ongoing monitoring.",
			Component:   "ThirdParty",
			Category:    "ThirdPartyRisk",
			Risk:        RiskHigh,
			Likelihood:  4,
			Impact:      4,
			Confidence:  0.80,
			Keywords:    []string{"third-party", "vendor risk", "assessment", "monitoring"},
			Rationale:   "Third-party components without risk assessment introduce unknown security posture and supply chain risk.",
		})
		results = append(results, Assumption{
			ID:          "INF-TP-002",
			Description: "Third-party integration is present but architecture does not specify equivalent security controls requirement in contracts.",
			Component:   "ThirdParty",
			Category:    "VendorRisk",
			Risk:        RiskHigh,
			Likelihood:  3,
			Impact:      4,
			Confidence:  0.75,
			Keywords:    []string{"third-party", "contract", "security controls", "sla"},
			Rationale:   "Contracts without security control requirements do not enforce third-party accountability.",
		})
		results = append(results, Assumption{
			ID:          "INF-TP-003",
			Description: "Third-party integration has data access but architecture does not specify data minimization and egress controls.",
			Component:   "ThirdParty",
			Category:    "ThirdPartyRisk",
			Risk:        RiskCritical,
			Likelihood:  4,
			Impact:      5,
			Confidence:  0.85,
			Keywords:    []string{"third-party", "data minimization", "egress", "data access"},
			Rationale:   "Third-party analytics with database access without data minimization increases data exposure and regulatory risk.",
		})
	}
	return results
}

// inferAdminConsoleAssumptions: AdminConsole → MFA, break-glass, audit.
func (re *ReasoningEngine) inferAdminConsoleAssumptions() []Assumption {
	var results []Assumption
	hasAdmin := false
	for _, comp := range re.arch.Components {
		label := strings.ToLower(comp.Label)
		if strings.Contains(label, "admin") || strings.Contains(label, "console") || strings.Contains(label, "dashboard") || strings.Contains(label, "management") || strings.Contains(label, "portal") {
			hasAdmin = true
		}
	}
	raw := strings.ToLower(re.arch.RawText)
	if strings.Contains(raw, "admin console") || strings.Contains(raw, "management portal") || strings.Contains(raw, "admin dashboard") {
		hasAdmin = true
	}

	if hasAdmin {
		results = append(results, Assumption{
			ID:          "INF-ADM-001",
			Description: "Admin console is present but architecture does not specify MFA enforcement for all administrative access.",
			Component:   "AdminConsole",
			Category:    "Authentication",
			Risk:        RiskCritical,
			Likelihood:  4,
			Impact:      5,
			Confidence:  0.85,
			Keywords:    []string{"admin", "mfa", "console", "privileged access"},
			Rationale:   "Admin console without MFA is a high-value target for compromise; MFA is mandatory for privileged access.",
		})
		results = append(results, Assumption{
			ID:          "INF-ADM-002",
			Description: "Admin console is present but architecture does not specify break-glass procedures and emergency access logging.",
			Component:   "AdminConsole",
			Category:    "PrivilegeManagement",
			Risk:        RiskHigh,
			Likelihood:  3,
			Impact:      4,
			Confidence:  0.80,
			Keywords:    []string{"admin", "break-glass", "emergency access", "console"},
			Rationale:   "Break-glass procedures ensure emergency access is controlled, time-limited, and fully audited.",
		})
		results = append(results, Assumption{
			ID:          "INF-ADM-003",
			Description: "Admin console is present but architecture does not specify comprehensive audit logging for all administrative actions.",
			Component:   "AdminConsole",
			Category:    "Auditability",
			Risk:        RiskHigh,
			Likelihood:  3,
			Impact:      4,
			Confidence:  0.80,
			Keywords:    []string{"admin", "audit", "console", "logging"},
			Rationale:   "Admin actions must be fully audited to detect abuse, support investigations, and meet compliance requirements.",
		})
	}
	return results
}

// inferAuditLogAssumptions: AuditLog → immutability, retention, tamper detection.
func (re *ReasoningEngine) inferAuditLogAssumptions() []Assumption {
	var results []Assumption
	hasAuditLog := false
	for _, comp := range re.arch.Components {
		label := strings.ToLower(comp.Label)
		if strings.Contains(label, "audit") || strings.Contains(label, "log") || strings.Contains(label, "siem") || strings.Contains(label, "logging") || strings.Contains(label, "monitoring") {
			hasAuditLog = true
		}
	}
	raw := strings.ToLower(re.arch.RawText)
	if strings.Contains(raw, "audit log") || strings.Contains(raw, "audit trail") || strings.Contains(raw, "immutability") {
		hasAuditLog = true
	}

	if hasAuditLog {
		results = append(results, Assumption{
			ID:          "INF-AUD-001",
			Description: "Audit logging is present but architecture does not specify log immutability and write-once guarantees.",
			Component:   "AuditLog",
			Category:    "Auditability",
			Risk:        RiskCritical,
			Likelihood:  4,
			Impact:      5,
			Confidence:  0.85,
			Keywords:    []string{"audit", "immutability", "write-once", "tamper"},
			Rationale:   "Mutable audit logs can be altered to hide malicious activity; immutability is required for non-repudiation.",
		})
		results = append(results, Assumption{
			ID:          "INF-AUD-002",
			Description: "Audit logging is present but architecture does not specify retention policy aligned with regulatory and forensic requirements.",
			Component:   "AuditLog",
			Category:    "DataRetention",
			Risk:        RiskHigh,
			Likelihood:  3,
			Impact:      4,
			Confidence:  0.75,
			Keywords:    []string{"audit", "retention", "forensic", "compliance"},
			Rationale:   "Audit log retention must meet regulatory minimums and forensic investigation windows; absence creates compliance gaps.",
		})
		results = append(results, Assumption{
			ID:          "INF-AUD-003",
			Description: "Audit logging is present but architecture does not specify tamper detection and integrity verification mechanisms.",
			Component:   "AuditLog",
			Category:    "Auditability",
			Risk:        RiskHigh,
			Likelihood:  4,
			Impact:      4,
			Confidence:  0.80,
			Keywords:    []string{"audit", "tamper detection", "integrity", "hash"},
			Rationale:   "Tamper detection ensures audit log integrity; without it, forensic evidence may be challenged or invalidated.",
		})
	}
	return results
}

// inferFromRawText generates additional assumptions from raw text patterns.
func (re *ReasoningEngine) inferFromRawText() []Assumption {
	if re.arch == nil || re.arch.RawText == "" {
		return nil
	}
	var results []Assumption
	raw := strings.ToLower(re.arch.RawText)

	if (strings.Contains(raw, "encryption") || strings.Contains(raw, "encrypted")) && !strings.Contains(raw, "key management") && !strings.Contains(raw, "kms") {
		results = append(results, Assumption{
			ID:          "INF-TXT-001",
			Description: "Architecture mentions encryption but does not specify key management or key rotation controls.",
			Component:   "General",
			Category:    "KeyManagement",
			Risk:        RiskHigh,
			Likelihood:  3,
			Impact:      4,
			Confidence:  0.70,
			Keywords:    []string{"encryption", "key management", "rotation"},
			Rationale:   "Encryption without key management is incomplete; key lifecycle must be defined and enforced.",
		})
	}
	if strings.Contains(raw, "tls") && strings.Contains(raw, "http") {
		results = append(results, Assumption{
			ID:          "INF-TXT-002",
			Description: "Architecture contains both TLS and HTTP references but does not mandate TLS for all communications.",
			Component:   "General",
			Category:    "NetworkSegmentation",
			Risk:        RiskHigh,
			Likelihood:  4,
			Impact:      4,
			Confidence:  0.75,
			Keywords:    []string{"tls", "http", "mandate", "encryption"},
			Rationale:   "Coexistence of TLS and HTTP without enforcement creates downgrade and interception risks.",
		})
	}
	return results
}

// InferAllAssumptions runs all inference rules including raw text.
func (re *ReasoningEngine) InferAllAssumptions() []Assumption {
	all := re.InferAssumptions()
	all = append(all, re.inferFromRawText()...)
	return deduplicateAssumptions(all)
}

func deduplicateAssumptions(assumptions []Assumption) []Assumption {
	seen := make(map[string]bool)
	var result []Assumption
	for _, a := range assumptions {
		key := a.ID + "|" + strings.ToLower(a.Description)
		if !seen[key] {
			seen[key] = true
			result = append(result, a)
		}
	}
	return result
}

// minInt returns the minimum of two integers.
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// maxInt returns the maximum of two integers.
func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// containsAnyLower checks if any target string is contained in items (case-insensitive).
func containsAnyLower(items []string, targets ...string) bool {
	for _, item := range items {
		il := strings.ToLower(item)
		for _, t := range targets {
			if strings.Contains(il, t) {
				return true
			}
		}
	}
	return false
}

// buildExplainability generates a human-readable explainability string for an inferred assumption.
func buildExplainability(template, context string, missing []string, confidence float64, category string) string {
	parts := []string{template}
	if context != "" {
		parts = append(parts, fmt.Sprintf("Context: %s", context))
	}
	if len(missing) > 0 {
		parts = append(parts, fmt.Sprintf("Missing controls: %s", strings.Join(missing, ", ")))
	}
	parts = append(parts, fmt.Sprintf("Confidence: %.0f%%", confidence*100))
	parts = append(parts, fmt.Sprintf("Category: %s", category))
	return strings.Join(parts, "; ")
}
