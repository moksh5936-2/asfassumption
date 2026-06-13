package verify

type evidenceRule struct {
	categories    []EvidenceCategory
	keywords      []string
	evidence      []EvidenceSource
	actions       []VerificationAction
	stakeholders  []string
	whyVerify     string
	whatToReview  string
	whatEvidence  string
	howToValidate string
	expectedTime  string
}

var evidenceTaxonomy = []evidenceRule{
	{
		categories: []EvidenceCategory{EvCatIdentity},
		keywords:   []string{"mfa", "multi-factor", "2fa", "two-factor"},
		evidence: []EvidenceSource{
			{Type: SourcePolicyDocument, Name: "MFA Policy", Description: "MFA enforcement policy document", Optional: false},
			{Type: SourceConfiguration, Name: "IdP Configuration", Description: "Identity provider MFA configuration", Optional: false},
			{Type: SourceAuditLog, Name: "Access Logs", Description: "MFA authentication audit logs", Optional: false},
		},
		actions: []VerificationAction{
			{Step: 1, Action: "Review MFA policy", Description: "Verify MFA is enforced for all users", Stakeholder: "Security Architect"},
			{Step: 2, Action: "Check IdP configuration", Description: "Verify MFA enforcement settings in identity provider", Stakeholder: "IAM Team"},
			{Step: 3, Action: "Review access logs", Description: "Verify MFA is recorded in authentication logs", Stakeholder: "Security Engineer"},
		},
		stakeholders:  []string{"Security Architect", "IAM Team", "Security Engineer"},
		whyVerify:     "MFA is a critical control for preventing credential-based attacks and account takeover",
		whatToReview:  "MFA enforcement policy, identity provider configuration, and authentication audit logs",
		whatEvidence:  "MFA policy document, IdP configuration screenshots, access log samples showing MFA verification",
		howToValidate: "Confirm MFA is enforced for all users, including administrators and service accounts",
		expectedTime:  "2-3 hours",
	},
	{
		categories: []EvidenceCategory{EvCatIdentity},
		keywords:   []string{"sso", "single sign-on", "federat"},
		evidence: []EvidenceSource{
			{Type: SourcePolicyDocument, Name: "SSO Policy", Description: "SSO implementation policy", Optional: false},
			{Type: SourceConfiguration, Name: "Federation Configuration", Description: "Identity federation setup", Optional: false},
			{Type: SourceArtifact, Name: "SAML/OIDC Metadata", Description: "Federation metadata and certificates", Optional: true},
		},
		actions: []VerificationAction{
			{Step: 1, Action: "Review SSO architecture", Description: "Verify SSO implementation aligns with architecture", Stakeholder: "Security Architect"},
			{Step: 2, Action: "Check federation configuration", Description: "Verify federation trust relationships", Stakeholder: "IAM Team"},
			{Step: 3, Action: "Review metadata", Description: "Verify SAML/OIDC metadata and certificate validity", Stakeholder: "Security Engineer"},
		},
		stakeholders:  []string{"Security Architect", "IAM Team"},
		whyVerify:     "SSO is a centralized authentication control; misconfiguration affects all connected services",
		whatToReview:  "SSO policy, federation configuration, and trust relationships between identity providers",
		whatEvidence:  "SSO policy document, federation metadata, IdP configuration documentation",
		howToValidate: "Verify SSO flow end-to-end and confirm federation trust is properly scoped",
		expectedTime:  "3-4 hours",
	},
	{
		categories: []EvidenceCategory{EvCatAuthorization},
		keywords:   []string{"rbac", "role", "role-based"},
		evidence: []EvidenceSource{
			{Type: SourcePolicyDocument, Name: "RBAC Policy", Description: "Role-based access control policy", Optional: false},
			{Type: SourceReport, Name: "Role Matrix", Description: "Role definitions and assignment matrix", Optional: false},
			{Type: SourceAuditLog, Name: "Access Reviews", Description: "Periodic access review records", Optional: false},
		},
		actions: []VerificationAction{
			{Step: 1, Action: "Review RBAC policy", Description: "Verify RBAC policy covers all roles and permissions", Stakeholder: "Security Architect"},
			{Step: 2, Action: "Review role matrix", Description: "Verify role definitions follow least privilege", Stakeholder: "IAM Team"},
			{Step: 3, Action: "Review access reviews", Description: "Verify access reviews are conducted periodically", Stakeholder: "Audit Team"},
			{Step: 4, Action: "Review privileged access", Description: "Verify privileged role assignments and approvals", Stakeholder: "Security Engineer"},
			{Step: 5, Action: "Review exception process", Description: "Verify exception handling for access requests", Stakeholder: "Security Architect"},
		},
		stakeholders:  []string{"Security Architect", "IAM Team", "Audit Team", "Security Engineer"},
		whyVerify:     "RBAC is foundational to authorization; misconfiguration leads to privilege escalation",
		whatToReview:  "RBAC policy, role matrix, access review records, privilege assignments, exception process",
		whatEvidence:  "RBAC policy document, role definitions, access review reports, privilege assignments",
		howToValidate: "Confirm least privilege principle is applied and access reviews are current",
		expectedTime:  "4-6 hours",
	},
	{
		categories: []EvidenceCategory{EvCatAuthorization},
		keywords:   []string{"least privilege", "principle of least"},
		evidence: []EvidenceSource{
			{Type: SourcePolicyDocument, Name: "Least Privilege Policy", Description: "Least privilege access policy", Optional: false},
			{Type: SourceAuditLog, Name: "Permission Audits", Description: "Permission audit records and reports", Optional: false},
		},
		actions: []VerificationAction{
			{Step: 1, Action: "Review least privilege policy", Description: "Verify least privilege principle is documented", Stakeholder: "Security Architect"},
			{Step: 2, Action: "Review permission audits", Description: "Verify permission audit reports show compliance", Stakeholder: "Audit Team"},
		},
		stakeholders:  []string{"Security Architect", "Audit Team"},
		whyVerify:     "Least privilege limits blast radius of compromised accounts and insider threats",
		whatToReview:  "Least privilege policy, permission audit reports, and access control implementation",
		whatEvidence:  "Least privilege policy, permission audit reports, access control tests",
		howToValidate: "Verify users have minimum required permissions and excessive permissions are identified",
		expectedTime:  "2-3 hours",
	},
	{
		categories: []EvidenceCategory{EvCatCryptography},
		keywords:   []string{"tls", "ssl", "in transit", "https"},
		evidence: []EvidenceSource{
			{Type: SourceConfiguration, Name: "TLS Configuration", Description: "TLS/SSL certificate configuration", Optional: false},
			{Type: SourceToolOutput, Name: "Certificate Scan", Description: "Certificate validity and strength scan", Optional: false},
			{Type: SourceAuditLog, Name: "TLS Audit Log", Description: "TLS handshake and cipher audit logs", Optional: true},
		},
		actions: []VerificationAction{
			{Step: 1, Action: "Review TLS configuration", Description: "Verify TLS version and cipher configuration", Stakeholder: "Security Engineer"},
			{Step: 2, Action: "Run certificate scan", Description: "Verify certificate validity, chain, and expiration", Stakeholder: "Security Engineer"},
			{Step: 3, Action: "Review TLS audit logs", Description: "Verify TLS configuration in audit trail", Stakeholder: "Security Architect"},
		},
		stakeholders:  []string{"Security Engineer", "Security Architect"},
		whyVerify:     "TLS protects data in transit; weak configuration exposes data to interception",
		whatToReview:  "TLS version, cipher suites, certificate validity, and certificate chain of trust",
		whatEvidence:  "TLS configuration files, certificate scan results, cipher suite documentation",
		howToValidate: "Confirm TLS 1.2+ with strong ciphers is enforced and certificates are valid",
		expectedTime:  "1-2 hours",
	},
	{
		categories: []EvidenceCategory{EvCatCryptography},
		keywords:   []string{"kms", "key", "vault", "hsm", "keyvault", "secrets manager", "key rotation", "rotate"},
		evidence: []EvidenceSource{
			{Type: SourceConfiguration, Name: "KMS Configuration", Description: "Key management service configuration", Optional: false},
			{Type: SourcePolicyDocument, Name: "KMS Policy", Description: "Key management policies and procedures", Optional: false},
			{Type: SourceAuditLog, Name: "Rotation Records", Description: "Key rotation audit logs and schedule", Optional: false},
		},
		actions: []VerificationAction{
			{Step: 1, Action: "Review KMS configuration", Description: "Verify key management service configuration", Stakeholder: "Security Engineer"},
			{Step: 2, Action: "Review key rotation schedule", Description: "Verify rotation schedule meets compliance requirements", Stakeholder: "Security Architect"},
			{Step: 3, Action: "Review rotation audit history", Description: "Verify key rotations are logged and audited", Stakeholder: "Audit Team"},
		},
		stakeholders:  []string{"Security Engineer", "Security Architect", "Audit Team"},
		whyVerify:     "Key management is critical for data confidentiality; missing rotation increases exposure window",
		whatToReview:  "KMS configuration, key policies, rotation schedule, and rotation audit records",
		whatEvidence:  "KMS configuration exports, key rotation policy, rotation audit logs, key metadata",
		howToValidate: "Confirm keys are rotated according to policy and rotation events are logged",
		expectedTime:  "3-4 hours",
	},
	{
		categories: []EvidenceCategory{EvCatMonitoring},
		keywords:   []string{"siem", "splunk", "elastic", "log", "audit", "cloudtrail", "cloudwatch"},
		evidence: []EvidenceSource{
			{Type: SourceConfiguration, Name: "SIEM Configuration", Description: "SIEM/logging platform configuration", Optional: false},
			{Type: SourceReport, Name: "Log Samples", Description: "Sample log entries showing coverage", Optional: false},
			{Type: SourcePolicyDocument, Name: "Alert Rules", Description: "Alert rule definitions and thresholds", Optional: false},
			{Type: SourceAuditLog, Name: "Review Records", Description: "Periodic monitoring review records", Optional: true},
		},
		actions: []VerificationAction{
			{Step: 1, Action: "Review SIEM configuration", Description: "Verify logging sources are integrated", Stakeholder: "Security Engineer"},
			{Step: 2, Action: "Review log samples", Description: "Verify log coverage across all components", Stakeholder: "Security Architect"},
			{Step: 3, Action: "Review alert rules", Description: "Verify security alert rules cover key scenarios", Stakeholder: "Security Architect"},
			{Step: 4, Action: "Review monitoring records", Description: "Verify monitoring reviews are conducted", Stakeholder: "Security Engineer"},
		},
		stakeholders:  []string{"Security Engineer", "Security Architect"},
		whyVerify:     "Monitoring detects security incidents; gaps in logging create blind spots",
		whatToReview:  "SIEM configuration, log sources, alert rules, and monitoring review records",
		whatEvidence:  "SIEM integration docs, sample logs, alert rule definitions, review records",
		howToValidate: "Confirm all critical components are logged and alert rules cover key attack scenarios",
		expectedTime:  "4-6 hours",
	},
	{
		categories: []EvidenceCategory{EvCatResilience},
		keywords:   []string{"backup", "snapshot", "restore", "recovery", "disaster"},
		evidence: []EvidenceSource{
			{Type: SourceReport, Name: "Backup Reports", Description: "Backup execution and status reports", Optional: false},
			{Type: SourceReport, Name: "Restore Exercises", Description: "Restore testing results and records", Optional: false},
			{Type: SourcePolicyDocument, Name: "Retention Policies", Description: "Backup retention and lifecycle policies", Optional: false},
		},
		actions: []VerificationAction{
			{Step: 1, Action: "Review backup reports", Description: "Verify backup success rates and coverage", Stakeholder: "Infrastructure Team"},
			{Step: 2, Action: "Review restore exercises", Description: "Verify restore testing is performed regularly", Stakeholder: "Infrastructure Team"},
			{Step: 3, Action: "Review retention policies", Description: "Verify retention meets compliance requirements", Stakeholder: "Security Architect"},
		},
		stakeholders:  []string{"Infrastructure Team", "Security Architect"},
		whyVerify:     "Backup without restore testing is not a backup; untested backups provide false confidence",
		whatToReview:  "Backup success logs, restore test results, retention policies, and DR exercise records",
		whatEvidence:  "Backup reports, restore test results, retention policy documents, DR exercise records",
		howToValidate: "Confirm backups succeed consistently and restore tests are conducted quarterly",
		expectedTime:  "3-4 hours",
	},
	{
		categories: []EvidenceCategory{EvCatThirdParty},
		keywords:   []string{"auth0", "okta", "keycloak", "idp", "third-party", "vendor", "saas"},
		evidence: []EvidenceSource{
			{Type: SourceVendorDocument, Name: "Vendor Security Documentation", Description: "Vendor security posture documentation", Optional: false},
			{Type: SourceReport, Name: "SOC Reports", Description: "Vendor SOC 2/3 audit reports", Optional: false},
			{Type: SourceArtifact, Name: "Configuration Review", Description: "Third-party service configuration review", Optional: false},
		},
		actions: []VerificationAction{
			{Step: 1, Action: "Review vendor security docs", Description: "Verify vendor security controls and certifications", Stakeholder: "Security Architect"},
			{Step: 2, Action: "Review SOC reports", Description: "Verify vendor audit reports are current and complete", Stakeholder: "Audit Team"},
			{Step: 3, Action: "Review configuration", Description: "Verify third-party service security configuration", Stakeholder: "Security Engineer"},
		},
		stakeholders:  []string{"Security Architect", "Audit Team", "Security Engineer"},
		whyVerify:     "Third-party services extend the trust boundary; vendor compromise affects your architecture",
		whatToReview:  "Vendor security documentation, SOC reports, and service configuration",
		whatEvidence:  "Vendor security docs, SOC reports, configuration review results",
		howToValidate: "Confirm vendor meets security requirements and configuration follows best practices",
		expectedTime:  "4-6 hours",
	},
	{
		categories: []EvidenceCategory{EvCatOperational},
		keywords:   []string{"secret", "credential", "vault", "token", "pipeline", "ci/cd", "build"},
		evidence: []EvidenceSource{
			{Type: SourceConfiguration, Name: "Secrets Management Configuration", Description: "Secrets management platform configuration", Optional: false},
			{Type: SourcePolicyDocument, Name: "Secrets Policy", Description: "Secrets management policy and procedures", Optional: false},
			{Type: SourceAuditLog, Name: "Access Audit Logs", Description: "Secrets access audit records", Optional: false},
		},
		actions: []VerificationAction{
			{Step: 1, Action: "Review secrets management config", Description: "Verify secrets management platform integration", Stakeholder: "Security Engineer"},
			{Step: 2, Action: "Review secrets policy", Description: "Verify secrets management procedures and rotation", Stakeholder: "Security Architect"},
			{Step: 3, Action: "Review access audit logs", Description: "Verify secrets access is logged and monitored", Stakeholder: "Security Engineer"},
		},
		stakeholders:  []string{"Security Engineer", "Security Architect"},
		whyVerify:     "Hardcoded secrets are a leading cause of breaches; proper secrets management is essential",
		whatToReview:  "Secrets management platform, policies, rotation procedures, and access logs",
		whatEvidence:  "Secrets management configuration, secrets policy, access audit logs",
		howToValidate: "Confirm no hardcoded secrets exist and secrets management follows security best practices",
		expectedTime:  "3-4 hours",
	},
	{
		categories: []EvidenceCategory{EvCatOperational},
		keywords:   []string{"rate limit", "throttle", "quota", "api gateway"},
		evidence: []EvidenceSource{
			{Type: SourceConfiguration, Name: "Rate Limiting Configuration", Description: "API rate limiting and throttling rules", Optional: false},
			{Type: SourceAuditLog, Name: "Rate Limit Logs", Description: "Rate limiting events and throttling records", Optional: true},
		},
		actions: []VerificationAction{
			{Step: 1, Action: "Review rate limiting config", Description: "Verify rate limiting rules and thresholds", Stakeholder: "Security Engineer"},
			{Step: 2, Action: "Review throttle logs", Description: "Verify rate limiting is enforced and logged", Stakeholder: "Security Architect"},
		},
		stakeholders:  []string{"Security Engineer", "Security Architect"},
		whyVerify:     "Without rate limiting, APIs are vulnerable to DoS and brute-force attacks",
		whatToReview:  "Rate limiting configuration, threshold values, and enforcement logs",
		whatEvidence:  "Rate limiting configuration, throttle logs, API gateway configuration",
		howToValidate: "Confirm rate limiting is enforced at the API gateway and logs are generated",
		expectedTime:  "1-2 hours",
	},
}

type domainEvidenceRule struct {
	domains      []string
	categories   []EvidenceCategory
	keywords     []string
	evidence     []EvidenceSource
	actions      []VerificationAction
	stakeholders []string
	whyVerify    string
}

var domainEvidenceTaxonomy = []domainEvidenceRule{
	{
		domains:    []string{"healthcare", "hipaa"},
		categories: []EvidenceCategory{EvCatIdentity},
		keywords:   []string{"phi", "access", "patient"},
		evidence: []EvidenceSource{
			{Type: SourcePolicyDocument, Name: "PHI Access Controls", Description: "Protected health information access control policy", Optional: false},
			{Type: SourceAuditLog, Name: "PHI Access Logs", Description: "PHI access audit logs and ePHI disclosure records", Optional: false},
			{Type: SourceArtifact, Name: "Break Glass Procedure", Description: "Emergency access break glass procedure documentation", Optional: false},
		},
		actions: []VerificationAction{
			{Step: 1, Action: "Review PHI access controls", Description: "Verify PHI access is restricted to authorized personnel", Stakeholder: "Security Architect"},
			{Step: 2, Action: "Review PHI access logs", Description: "Verify PHI access is logged in accordance with HIPAA", Stakeholder: "Compliance Team"},
			{Step: 3, Action: "Review break glass procedures", Description: "Verify emergency access procedures are documented and audited", Stakeholder: "Security Engineer"},
		},
		stakeholders: []string{"Security Architect", "Compliance Team", "Security Engineer"},
		whyVerify:    "HIPAA requires strict PHI access controls; violations carry significant penalties",
	},
	{
		domains:    []string{"healthcare", "hipaa"},
		categories: []EvidenceCategory{EvCatMonitoring},
		keywords:   []string{"clinical", "logging", "audit"},
		evidence: []EvidenceSource{
			{Type: SourceConfiguration, Name: "Clinical Logging Config", Description: "Clinical system logging configuration", Optional: false},
			{Type: SourceAuditLog, Name: "Clinical Audit Logs", Description: "Clinical system audit log samples", Optional: false},
		},
		actions: []VerificationAction{
			{Step: 1, Action: "Review clinical logging config", Description: "Verify clinical systems are configured for logging", Stakeholder: "Security Engineer"},
			{Step: 2, Action: "Review audit logs", Description: "Verify clinical audit logs are complete and tamper-proof", Stakeholder: "Compliance Team"},
		},
		stakeholders: []string{"Security Engineer", "Compliance Team"},
		whyVerify:    "Clinical logging is required for HIPAA compliance and patient safety investigations",
	},
	{
		domains:    []string{"fintech", "financial"},
		categories: []EvidenceCategory{EvCatOperational},
		keywords:   []string{"settlement", "payment", "transaction"},
		evidence: []EvidenceSource{
			{Type: SourcePolicyDocument, Name: "Settlement Controls", Description: "Financial settlement control policy", Optional: false},
			{Type: SourceAuditLog, Name: "Settlement Audit Logs", Description: "Settlement reconciliation audit logs", Optional: false},
		},
		actions: []VerificationAction{
			{Step: 1, Action: "Review settlement controls", Description: "Verify settlement reconciliation and controls", Stakeholder: "Security Architect"},
			{Step: 2, Action: "Review settlement audit logs", Description: "Verify settlement transactions are fully auditable", Stakeholder: "Audit Team"},
		},
		stakeholders: []string{"Security Architect", "Audit Team"},
		whyVerify:    "Settlement controls prevent financial loss and ensure transaction integrity",
	},
	{
		domains:    []string{"fintech", "financial"},
		categories: []EvidenceCategory{EvCatMonitoring},
		keywords:   []string{"fraud", "monitor"},
		evidence: []EvidenceSource{
			{Type: SourceConfiguration, Name: "Fraud Monitoring Configuration", Description: "Fraud detection platform configuration", Optional: false},
			{Type: SourceReport, Name: "Fraud Reports", Description: "Fraud detection reports and case records", Optional: false},
		},
		actions: []VerificationAction{
			{Step: 1, Action: "Review fraud monitoring config", Description: "Verify fraud detection rules and thresholds", Stakeholder: "Security Engineer"},
			{Step: 2, Action: "Review fraud reports", Description: "Verify fraud incidents are detected and investigated", Stakeholder: "Security Architect"},
		},
		stakeholders: []string{"Security Engineer", "Security Architect"},
		whyVerify:    "Fraud monitoring is critical for financial services; gaps lead to undetected fraud",
	},
	{
		domains:    []string{"fintech", "financial"},
		categories: []EvidenceCategory{EvCatCryptography},
		keywords:   []string{"key", "custody"},
		evidence: []EvidenceSource{
			{Type: SourcePolicyDocument, Name: "Key Custody Policy", Description: "Cryptographic key custody and management policy", Optional: false},
			{Type: SourceAuditLog, Name: "Key Custody Logs", Description: "Key custody transfer and access logs", Optional: false},
		},
		actions: []VerificationAction{
			{Step: 1, Action: "Review key custody policy", Description: "Verify key custody procedures and segregation of duties", Stakeholder: "Security Architect"},
			{Step: 2, Action: "Review key custody logs", Description: "Verify key custody is logged and audited", Stakeholder: "Audit Team"},
		},
		stakeholders: []string{"Security Architect", "Audit Team"},
		whyVerify:    "Key custody controls prevent unauthorized access to financial signing keys",
	},
	{
		domains:    []string{"cloud", "aws", "azure", "gcp"},
		categories: []EvidenceCategory{EvCatIdentity},
		keywords:   []string{"iam", "federat"},
		evidence: []EvidenceSource{
			{Type: SourceConfiguration, Name: "IAM Configuration", Description: "Cloud IAM roles and policies", Optional: false},
			{Type: SourceAuditLog, Name: "IAM Audit Logs", Description: "Cloud IAM access audit logs and CloudTrail events", Optional: false},
		},
		actions: []VerificationAction{
			{Step: 1, Action: "Review IAM configuration", Description: "Verify IAM policies follow least privilege", Stakeholder: "Security Engineer"},
			{Step: 2, Action: "Review IAM audit logs", Description: "Verify IAM changes are logged and monitored", Stakeholder: "Security Architect"},
		},
		stakeholders: []string{"Security Engineer", "Security Architect"},
		whyVerify:    "Cloud IAM is the primary access control; misconfiguration leads to data exposure",
	},
	{
		domains:    []string{"kubernetes", "k8s"},
		categories: []EvidenceCategory{EvCatAuthorization},
		keywords:   []string{"rbac", "admission", "service account"},
		evidence: []EvidenceSource{
			{Type: SourceConfiguration, Name: "K8s RBAC Configuration", Description: "Kubernetes RBAC roles and bindings", Optional: false},
			{Type: SourceConfiguration, Name: "Admission Controller Config", Description: "Kubernetes admission controller configuration", Optional: false},
			{Type: SourceAuditLog, Name: "Service Account Audit", Description: "Kubernetes service account usage audit", Optional: false},
		},
		actions: []VerificationAction{
			{Step: 1, Action: "Review K8s RBAC", Description: "Verify RBAC roles and bindings follow least privilege", Stakeholder: "Security Engineer"},
			{Step: 2, Action: "Review admission controllers", Description: "Verify admission controllers enforce security policies", Stakeholder: "Security Architect"},
			{Step: 3, Action: "Review service accounts", Description: "Verify service account permissions and usage", Stakeholder: "Security Engineer"},
		},
		stakeholders: []string{"Security Engineer", "Security Architect"},
		whyVerify:    "Kubernetes RBAC misconfiguration is a leading cause of container security incidents",
	},
}

func lookupEvidence(keywords []string, category EvidenceCategory) []EvidenceSource {
	for _, rule := range evidenceTaxonomy {
		if rule.categoryMatches(category) && keywordMatches(rule.keywords, keywords) {
			return rule.evidence
		}
	}
	return nil
}

func lookupActions(keywords []string, category EvidenceCategory) []VerificationAction {
	for _, rule := range evidenceTaxonomy {
		if rule.categoryMatches(category) && keywordMatches(rule.keywords, keywords) {
			return rule.actions
		}
	}
	return nil
}

func lookupStakeholders(keywords []string, category EvidenceCategory) []string {
	for _, rule := range evidenceTaxonomy {
		if rule.categoryMatches(category) && keywordMatches(rule.keywords, keywords) {
			return rule.stakeholders
		}
	}
	return nil
}

func lookupWhyVerify(keywords []string, category EvidenceCategory) string {
	for _, rule := range evidenceTaxonomy {
		if rule.categoryMatches(category) && keywordMatches(rule.keywords, keywords) {
			return rule.whyVerify
		}
	}
	return ""
}

func lookupWhatToReview(keywords []string, category EvidenceCategory) string {
	for _, rule := range evidenceTaxonomy {
		if rule.categoryMatches(category) && keywordMatches(rule.keywords, keywords) {
			return rule.whatToReview
		}
	}
	return ""
}

func lookupWhatEvidence(keywords []string, category EvidenceCategory) string {
	for _, rule := range evidenceTaxonomy {
		if rule.categoryMatches(category) && keywordMatches(rule.keywords, keywords) {
			return rule.whatEvidence
		}
	}
	return ""
}

func lookupHowToValidate(keywords []string, category EvidenceCategory) string {
	for _, rule := range evidenceTaxonomy {
		if rule.categoryMatches(category) && keywordMatches(rule.keywords, keywords) {
			return rule.howToValidate
		}
	}
	return ""
}

func lookupExpectedTime(keywords []string, category EvidenceCategory) string {
	for _, rule := range evidenceTaxonomy {
		if rule.categoryMatches(category) && keywordMatches(rule.keywords, keywords) {
			return rule.expectedTime
		}
	}
	return "TBD"
}

func lookupDomainEvidence(domain string, category EvidenceCategory, keywords []string) []EvidenceSource {
	dl := toLower(domain)
	for _, rule := range domainEvidenceTaxonomy {
		domainMatch := false
		for _, d := range rule.domains {
			if dl == toLower(d) {
				domainMatch = true
				break
			}
		}
		if !domainMatch {
			continue
		}
		if rule.categories[0] == category || matchesAnyCategory(rule.categories, category) {
			if keywordMatch(rule.keywords, keywords) {
				return rule.evidence
			}
		}
	}
	return nil
}

func lookupDomainActions(domain string, category EvidenceCategory, keywords []string) []VerificationAction {
	dl := toLower(domain)
	for _, rule := range domainEvidenceTaxonomy {
		domainMatch := false
		for _, d := range rule.domains {
			if dl == toLower(d) {
				domainMatch = true
				break
			}
		}
		if !domainMatch {
			continue
		}
		if rule.categories[0] == category || matchesAnyCategory(rule.categories, category) {
			if keywordMatch(rule.keywords, keywords) {
				return rule.actions
			}
		}
	}
	return nil
}

func lookupDomainStakeholders(domain string, category EvidenceCategory, keywords []string) []string {
	dl := toLower(domain)
	for _, rule := range domainEvidenceTaxonomy {
		domainMatch := false
		for _, d := range rule.domains {
			if dl == toLower(d) {
				domainMatch = true
				break
			}
		}
		if !domainMatch {
			continue
		}
		if rule.categories[0] == category || matchesAnyCategory(rule.categories, category) {
			if keywordMatch(rule.keywords, keywords) {
				return rule.stakeholders
			}
		}
	}
	return nil
}

func lookupDomainWhyVerify(domain string, category EvidenceCategory, keywords []string) string {
	dl := toLower(domain)
	for _, rule := range domainEvidenceTaxonomy {
		domainMatch := false
		for _, d := range rule.domains {
			if dl == toLower(d) {
				domainMatch = true
				break
			}
		}
		if !domainMatch {
			continue
		}
		if rule.categories[0] == category || matchesAnyCategory(rule.categories, category) {
			if keywordMatch(rule.keywords, keywords) {
				return rule.whyVerify
			}
		}
	}
	return ""
}

func (r *evidenceRule) categoryMatches(cat EvidenceCategory) bool {
	for _, c := range r.categories {
		if c == cat {
			return true
		}
	}
	return false
}

func matchesAnyCategory(cats []EvidenceCategory, cat EvidenceCategory) bool {
	for _, c := range cats {
		if c == cat {
			return true
		}
	}
	return false
}

func keywordMatches(ruleKeywords []string, assumptionKeywords []string) bool {
	for _, rk := range ruleKeywords {
		for _, ak := range assumptionKeywords {
			if contains(toLower(rk), toLower(ak)) || contains(toLower(ak), toLower(rk)) {
				return true
			}
		}
	}
	return false
}

func keywordMatch(ruleKeywords []string, assumptionKeywords []string) bool {
	return keywordMatches(ruleKeywords, assumptionKeywords)
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && containsStr(s, substr)
}

func containsStr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	b := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 32
		}
		b[i] = c
	}
	return string(b)
}
