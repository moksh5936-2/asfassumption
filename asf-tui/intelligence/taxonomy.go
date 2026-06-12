package intelligence

import (
	"regexp"
	"strings"
)

// TaxonomyCategory defines a complete assumption category with
// matching rules, risk mappings, and explainability templates.
type TaxonomyCategory struct {
	Name                   string
	Keywords               []string
	Patterns               []*regexp.Regexp
	RiskMappings           map[string]RiskLevel
	VerificationRules      []string
	ExplainabilityTemplate string
}

// TaxonomyEngine holds the complete assumption taxonomy.
type TaxonomyEngine struct {
	Categories map[string]*TaxonomyCategory
}

// NewTaxonomyEngine creates a taxonomy engine with all registered categories.
func NewTaxonomyEngine() *TaxonomyEngine {
	te := &TaxonomyEngine{
		Categories: make(map[string]*TaxonomyCategory),
	}
	te.registerAll()
	return te
}

// GetCategory returns a category by name.
func (te *TaxonomyEngine) GetCategory(name string) *TaxonomyCategory {
	return te.Categories[name]
}

// MatchCategory returns all category names that match the given text.
func (te *TaxonomyEngine) MatchCategory(text string) []string {
	lower := strings.ToLower(text)
	var matches []string
	seen := make(map[string]bool)
	for name, cat := range te.Categories {
		if seen[name] {
			continue
		}
		matched := false
		for _, kw := range cat.Keywords {
			if strings.Contains(lower, kw) {
				matched = true
				break
			}
		}
		if !matched {
			for _, pat := range cat.Patterns {
				if pat.MatchString(text) {
					matched = true
					break
				}
			}
		}
		if matched {
			matches = append(matches, name)
			seen[name] = true
		}
	}
	return matches
}

// registerAll populates the 40+ taxonomy categories.
func (te *TaxonomyEngine) registerAll() {
	categories := []*TaxonomyCategory{
		{
			Name:     "Identity",
			Keywords: []string{"identity", "identifier", "user identity", "digital identity", "identity proofing"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\bidentity\s+(provider|verification|proofing|federation)\b`),
				regexp.MustCompile(`(?i)\buser\s+identity\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":   RiskHigh,
				"federated": RiskCritical,
			},
			VerificationRules: []string{
				"Verify identity uniqueness across the system",
				"Confirm identity lifecycle management (provision, update, deprovision)",
				"Validate identity proofing strength against assurance level",
			},
			ExplainabilityTemplate: "Architecture references identity services but does not specify identity verification strength, lifecycle controls, or federation risks.",
		},
		{
			Name:     "Authentication",
			Keywords: []string{"authentication", "authn", "login", "password", "credential", "authenticate", "sso", "oauth", "oidc", "saml", "auth0"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(mfa|2fa|multi.factor|two.factor)\b`),
				regexp.MustCompile(`(?i)\b(single.sign.on|sso)\b`),
				regexp.MustCompile(`(?i)\bauth0\b`),
				regexp.MustCompile(`(?i)\b(password|credential|api.key|token)\s+(policy|rotation|expir)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":  RiskHigh,
				"mfa":      RiskMedium,
				"no_mfa":   RiskCritical,
				"password": RiskHigh,
			},
			VerificationRules: []string{
				"Verify MFA enrollment rate and enforcement",
				"Confirm authentication event logging and anomaly detection",
				"Validate credential storage (hashed, salted, peppered)",
			},
			ExplainabilityTemplate: "Architecture contains authentication mechanisms but does not specify MFA enforcement, credential lifecycle, or anomaly detection coverage.",
		},
		{
			Name:     "Authorization",
			Keywords: []string{"authorization", "authz", "access control", "rbac", "abac", "permission", "role", "acl"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(role.based|attribute.based|policy.based)\s+access\b`),
				regexp.MustCompile(`(?i)\bleast\s+privilege\b`),
				regexp.MustCompile(`(?i)\baccess\s+(control|matrix|policy)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":         RiskHigh,
				"rbac":            RiskMedium,
				"least_privilege": RiskMedium,
				"no_acl":          RiskCritical,
			},
			VerificationRules: []string{
				"Verify role definitions and privilege assignments",
				"Confirm authorization decisions are logged and auditable",
				"Validate separation of duties across critical roles",
			},
			ExplainabilityTemplate: "Architecture contains authorization references but does not specify role granularity, access reviews, or dynamic authorization policies.",
		},
		{
			Name:     "PrivilegeManagement",
			Keywords: []string{"privilege", "privileged", "admin", "root", "superuser", "sudo", "elevation", "impersonation"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\bprivilege\s+(escalation|elevation|management)\b`),
				regexp.MustCompile(`(?i)\badmin\s+(access|console|account|role)\b`),
				regexp.MustCompile(`(?i)\bbreak.glass\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":      RiskCritical,
				"break_glass":  RiskHigh,
				"just_in_time": RiskMedium,
			},
			VerificationRules: []string{
				"Verify privileged accounts are catalogued and monitored",
				"Confirm just-in-time access is enforced for critical privileges",
				"Validate break-glass procedures are documented and tested",
			},
			ExplainabilityTemplate: "Architecture references privileged access but does not specify privilege elevation controls, monitoring, or break-glass procedures.",
		},
		{
			Name:     "SecretsManagement",
			Keywords: []string{"secret", "secrets", "vault", "password", "api key", "apikey", "token", "credential", "hmac"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\bsecret\s+(management|store|rotation|engine)\b`),
				regexp.MustCompile(`(?i)\b(hashi|vault|aws secrets|azure key)\b`),
				regexp.MustCompile(`(?i)\bapi\s+(key|secret)\s+(storage|rotation|scope)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":   RiskHigh,
				"vault":     RiskMedium,
				"hardcoded": RiskCritical,
			},
			VerificationRules: []string{
				"Verify secrets are never hardcoded or committed to version control",
				"Confirm secrets rotation policy is enforced and measured",
				"Validate least-privilege scoping for API keys and tokens",
			},
			ExplainabilityTemplate: "Architecture references secrets but does not specify rotation frequency, storage hardening, or scope minimization.",
		},
		{
			Name:     "KeyManagement",
			Keywords: []string{"key management", "kms", "key rotation", "encryption key", "master key", "data key", "hsm"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(key|kms|hsm)\s+(rotation|management|lifecycle|escrow)\b`),
				regexp.MustCompile(`(?i)\baws\s+kms\b`),
				regexp.MustCompile(`(?i)\bazure\s+key\s+vault\b`),
				regexp.MustCompile(`(?i)\bgcp\s+cloud\s+kms\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":  RiskCritical,
				"rotation": RiskMedium,
				"hsm":      RiskMedium,
			},
			VerificationRules: []string{
				"Verify key rotation schedule aligns with data classification",
				"Confirm key access is restricted to least-privilege roles",
				"Validate key backup and escrow procedures are tested",
			},
			ExplainabilityTemplate: "Architecture contains encrypted data but does not specify key rotation cadence, access restriction, or key escrow.",
		},
		{
			Name:     "CertificateManagement",
			Keywords: []string{"certificate", "cert", "tls", "ssl", "x509", "ca", "pki", "public key"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(cert|certificate|tls|ssl)\s+(management|rotation|pinning|transparency)\b`),
				regexp.MustCompile(`(?i)\b(ca|pki|certificate.authority)\b`),
				regexp.MustCompile(`(?i)\bmutual\s+tls\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default": RiskHigh,
				"mtls":    RiskMedium,
				"expired": RiskCritical,
			},
			VerificationRules: []string{
				"Verify certificate expiry monitoring and renewal automation",
				"Confirm certificate transparency logging is enabled",
				"Validate mTLS enforcement for service-to-service communication",
			},
			ExplainabilityTemplate: "Architecture references TLS/SSL but does not specify certificate rotation, expiry monitoring, or mTLS enforcement.",
		},
		{
			Name:     "Cryptography",
			Keywords: []string{"cryptography", "crypto", "cipher", "aes", "rsa", "ecdsa", "sha", "hash", "random", "nonce"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(aes|rsa|ecdsa|sha.256|sha.384|sha.512|hmac)\b`),
				regexp.MustCompile(`(?i)\bcryptographic\s+(algorithm|protocol|primitive)\b`),
				regexp.MustCompile(`(?i)\bquantum\s+(safe|resistant)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default": RiskMedium,
				"weak":    RiskCritical,
				"legacy":  RiskHigh,
			},
			VerificationRules: []string{
				"Verify approved cipher suites and algorithm versions",
				"Confirm cryptographic agility plan for algorithm deprecation",
				"Validate randomness source for key generation and nonces",
			},
			ExplainabilityTemplate: "Architecture contains cryptographic references but does not specify approved algorithms, agility planning, or entropy validation.",
		},
		{
			Name:     "DataProtection",
			Keywords: []string{"data protection", "encryption", "data loss prevention", "dlp", "masking", "tokenization", "pii", "phi", "classification"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(encryption\s+(at\s+rest|in\s+transit|in\s+use)|field.level\s+encryption)\b`),
				regexp.MustCompile(`(?i)\b(dlp|data\s+loss\s+prevention)\b`),
				regexp.MustCompile(`(?i)\b(tokenization|masking|anonymization|pseudonymization)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default": RiskHigh,
				"phi":     RiskCritical,
				"pii":     RiskCritical,
				"at_rest": RiskMedium,
			},
			VerificationRules: []string{
				"Verify data classification schema is enforced at rest and in transit",
				"Confirm encryption covers all sensitive fields and backup copies",
				"Validate data masking and tokenization for non-production environments",
			},
			ExplainabilityTemplate: "Architecture contains sensitive data but does not specify encryption scope, classification enforcement, or data masking.",
		},
		{
			Name:     "DataRetention",
			Keywords: []string{"retention", "data retention", "archive", "deletion", "purge", "lifecycle", "data expiry", "ttl"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(data\s+retention|retention\s+policy|archive\s+policy)\b`),
				regexp.MustCompile(`(?i)\b(auto.delete|purge|ttl|expiry|data\s+lifecycle)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":   RiskMedium,
				"regulated": RiskHigh,
				"no_policy": RiskCritical,
			},
			VerificationRules: []string{
				"Verify retention periods align with regulatory and contractual requirements",
				"Confirm automated deletion and purge mechanisms are tested",
				"Validate audit trail for data deletion events",
			},
			ExplainabilityTemplate: "Architecture references data storage but does not specify retention periods, automated deletion, or purge validation.",
		},
		{
			Name:     "Privacy",
			Keywords: []string{"privacy", "gdpr", "ccpa", "consent", "right to be forgotten", "data subject", "privacy policy", "anonymization", "pseudonymization"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(gdpr|ccpa|lgpd|privacy\s+(law|regulation|impact))\b`),
				regexp.MustCompile(`(?i)\b(data\s+subject|consent|opt.out|right\s+to)\b`),
				regexp.MustCompile(`(?i)\bprivacy\s+by\s+(design|default)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default": RiskHigh,
				"gdpr":    RiskCritical,
				"consent": RiskMedium,
			},
			VerificationRules: []string{
				"Verify consent management and audit trail for data subject requests",
				"Confirm privacy impact assessments are completed for new features",
				"Validate data minimization and purpose limitation controls",
			},
			ExplainabilityTemplate: "Architecture contains personal data but does not specify consent management, privacy impact assessments, or data minimization.",
		},
		{
			Name:     "Logging",
			Keywords: []string{"logging", "log", "audit log", "application log", "event log", "syslog", "log aggregation"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(audit\s+log|security\s+log|application\s+log)\b`),
				regexp.MustCompile(`(?i)\b(log\s+aggregation|centralized\s+logging|siem)\b`),
				regexp.MustCompile(`(?i)\b(elk|splunk|datadog|cloudwatch\s+logs)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default": RiskMedium,
				"audit":   RiskHigh,
				"missing": RiskCritical,
			},
			VerificationRules: []string{
				"Verify all security-relevant events are logged with sufficient context",
				"Confirm log integrity and tamper detection mechanisms",
				"Validate log retention and accessibility for incident response",
			},
			ExplainabilityTemplate: "Architecture contains application components but does not specify logging coverage, integrity controls, or retention requirements.",
		},
		{
			Name:     "Monitoring",
			Keywords: []string{"monitoring", "metrics", "health check", "observability", "telemetry", "uptime", "sla"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(health\s+check|readiness|liveness|probe)\b`),
				regexp.MustCompile(`(?i)\b(prometheus|grafana|datadog|newrelic|dynatrace)\b`),
				regexp.MustCompile(`(?i)\b(slo|sla|error\s+budget)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default": RiskMedium,
				"sla":     RiskHigh,
				"missing": RiskCritical,
			},
			VerificationRules: []string{
				"Verify monitoring coverage for all critical components and paths",
				"Confirm alerting thresholds are tuned and tested",
				"Validate observability data does not contain sensitive fields",
			},
			ExplainabilityTemplate: "Architecture contains service components but does not specify monitoring coverage, alerting thresholds, or observability hygiene.",
		},
		{
			Name:     "Alerting",
			Keywords: []string{"alerting", "alert", "notification", "pagerduty", "on-call", "escalation", "runbook"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(alert\s+(rule|policy|threshold)|notification\s+channel)\b`),
				regexp.MustCompile(`(?i)\b(pagerduty|opsgenie|victorops|on.call)\b`),
				regexp.MustCompile(`(?i)\bescalation\s+(policy|matrix|procedure)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default": RiskHigh,
				"runbook": RiskMedium,
				"missing": RiskCritical,
			},
			VerificationRules: []string{
				"Verify alert fatigue is managed through deduplication and severity tuning",
				"Confirm on-call escalation paths are documented and tested",
				"Validate alert response runbooks cover security-critical scenarios",
			},
			ExplainabilityTemplate: "Architecture contains monitoring but does not specify alerting rules, escalation paths, or runbook coverage.",
		},
		{
			Name:     "Auditability",
			Keywords: []string{"auditability", "audit", "audit trail", "audit log", "non-repudiation", "accountability", "forensic"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(audit\s+trail|audit\s+log|audit\s+record)\b`),
				regexp.MustCompile(`(?i)\b(non.repudiation|accountability|forensic\s+readiness)\b`),
				regexp.MustCompile(`(?i)\bimmutable\s+log\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":   RiskHigh,
				"immutable": RiskMedium,
				"missing":   RiskCritical,
			},
			VerificationRules: []string{
				"Verify all privileged and sensitive actions produce immutable audit records",
				"Confirm audit logs are protected from tampering and deletion",
				"Validate audit log accessibility for compliance and incident response",
			},
			ExplainabilityTemplate: "Architecture contains sensitive operations but does not specify immutable audit logging, tamper protection, or forensic readiness.",
		},
		{
			Name:     "Backups",
			Keywords: []string{"backup", "backups", "snapshot", "restore", "replication", "copy", "dump"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(backup\s+(policy|schedule|retention)|snapshot\s+policy)\b`),
				regexp.MustCompile(`(?i)\b(disaster\s+recovery|dr|business\s+continuity)\b`),
				regexp.MustCompile(`(?i)\b(cross.region|cross.az|geo.redundant)\s+backup\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":     RiskMedium,
				"encrypted":   RiskMedium,
				"unencrypted": RiskCritical,
				"untested":    RiskHigh,
			},
			VerificationRules: []string{
				"Verify backup encryption and key management for all copies",
				"Confirm restore procedures are tested on a scheduled basis",
				"Validate backup access is restricted to least-privilege roles",
			},
			ExplainabilityTemplate: "Architecture references data persistence but does not specify backup encryption, restore testing, or geographic distribution.",
		},
		{
			Name:     "DisasterRecovery",
			Keywords: []string{"disaster recovery", "dr", "failover", "site recovery", "recovery point objective", "rpo", "recovery time objective", "rto"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(rpo|rto|recovery\s+(point|time)\s+objective)\b`),
				regexp.MustCompile(`(?i)\b(failover|failback|site\s+recovery|warm\s+standby|hot\s+standby)\b`),
				regexp.MustCompile(`(?i)\bdr\s+(plan|test|runbook)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":  RiskHigh,
				"rpo_rto":  RiskMedium,
				"untested": RiskCritical,
			},
			VerificationRules: []string{
				"Verify RPO and RTO targets are defined and measured",
				"Confirm DR failover and failback procedures are tested annually",
				"Validate data consistency and integrity after recovery events",
			},
			ExplainabilityTemplate: "Architecture references availability but does not specify DR targets, failover procedures, or recovery testing.",
		},
		{
			Name:     "Availability",
			Keywords: []string{"availability", "uptime", "high availability", "ha", "redundancy", "load balancer", "replica"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(high\s+availability|ha|multi.az|multi.region)\b`),
				regexp.MustCompile(`(?i)\b(load\s+balancer|lb|reverse\s+proxy|cdn)\b`),
				regexp.MustCompile(`(?i)\b(active.active|active.passive|cluster|replica)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":      RiskMedium,
				"sla":          RiskHigh,
				"single_point": RiskCritical,
			},
			VerificationRules: []string{
				"Verify availability targets are defined and monitored with SLOs",
				"Confirm single points of failure are identified and mitigated",
				"Validate failover and recovery mechanisms are tested regularly",
			},
			ExplainabilityTemplate: "Architecture contains service components but does not specify availability targets, redundancy design, or SLO monitoring.",
		},
		{
			Name:     "Resilience",
			Keywords: []string{"resilience", "fault tolerance", "graceful degradation", "circuit breaker", "retry", "throttle", "bulkhead", "timeout"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(circuit\s+breaker|retry\s+policy|backoff|timeout)\b`),
				regexp.MustCompile(`(?i)\b(graceful\s+degradation|bulkhead|throttle|rate\s+limit)\b`),
				regexp.MustCompile(`(?i)\bfault\s+tolerance\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default": RiskMedium,
				"chaos":   RiskMedium,
				"missing": RiskHigh,
			},
			VerificationRules: []string{
				"Verify circuit breakers and retry policies are configured for dependencies",
				"Confirm graceful degradation paths are tested under load",
				"Validate timeout and bulkhead settings protect cascading failures",
			},
			ExplainabilityTemplate: "Architecture contains service dependencies but does not specify resilience patterns, circuit breakers, or degradation paths.",
		},
		{
			Name:     "ThirdPartyRisk",
			Keywords: []string{"third party", "third-party", "external service", "saas", "api dependency", "external integration"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(third.party|thirdparty|external\s+(service|integration|api))\b`),
				regexp.MustCompile(`(?i)\b(saas|hosted\s+service|managed\s+service)\b`),
				regexp.MustCompile(`(?i)\b(data\s+processor|subprocessor)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":    RiskHigh,
				"assessed":   RiskMedium,
				"unassessed": RiskCritical,
			},
			VerificationRules: []string{
				"Verify third-party security assessments are completed and current",
				"Confirm contractual security requirements are defined and monitored",
				"Validate data minimization and egress controls for third-party integrations",
			},
			ExplainabilityTemplate: "Architecture contains third-party integrations but does not specify security assessments, contractual controls, or data minimization.",
		},
		{
			Name:     "VendorRisk",
			Keywords: []string{"vendor", "supplier", "contractor", "procurement", "vendor management", "supplier risk"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(vendor\s+(management|risk|assessment)|supplier\s+(audit|review))\b`),
				regexp.MustCompile(`(?i)\bprocurement\s+(security|policy|review)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":    RiskMedium,
				"assessed":   RiskLow,
				"unassessed": RiskCritical,
			},
			VerificationRules: []string{
				"Verify vendor risk assessments are performed before onboarding",
				"Confirm vendor access is reviewed and revoked upon contract termination",
				"Validate vendor incident notification requirements are in contracts",
			},
			ExplainabilityTemplate: "Architecture references vendor dependencies but does not specify vendor risk assessments, access reviews, or incident notification.",
		},
		{
			Name:     "SupplyChain",
			Keywords: []string{"supply chain", "dependency", "software supply chain", "sbom", "artifact", "build pipeline", "ci/cd"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(supply\s+chain|software\s+supply\s+chain|sbom)\b`),
				regexp.MustCompile(`(?i)\b(ci/cd|build\s+pipeline|artifact\s+repository)\b`),
				regexp.MustCompile(`(?i)\b(dependency\s+(scan|check|update)|dependabot|snyk)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":  RiskHigh,
				"sbom":     RiskMedium,
				"unsigned": RiskCritical,
			},
			VerificationRules: []string{
				"Verify SBOM generation and dependency vulnerability scanning",
				"Confirm build artifacts are signed and integrity-checked",
				"Validate CI/CD pipeline security controls and least-privilege access",
			},
			ExplainabilityTemplate: "Architecture contains build dependencies but does not specify SBOM generation, artifact signing, or pipeline security.",
		},
		{
			Name:     "InfrastructureSecurity",
			Keywords: []string{"infrastructure", "server", "host", "os hardening", "patch management", "baseline", "cve", "vulnerability"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(os\s+hardening|server\s+hardening|host\s+security)\b`),
				regexp.MustCompile(`(?i)\b(patch\s+management|vulnerability\s+(scan|management))\b`),
				regexp.MustCompile(`(?i)\b(cis\s+baseline|security\s+baseline|config\s+audit)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":   RiskMedium,
				"hardened":  RiskLow,
				"unpatched": RiskCritical,
			},
			VerificationRules: []string{
				"Verify infrastructure hardening baselines are applied and audited",
				"Confirm patch management covers OS, middleware, and firmware",
				"Validate vulnerability scanning and remediation SLAs",
			},
			ExplainabilityTemplate: "Architecture contains infrastructure but does not specify hardening baselines, patch management, or vulnerability scanning.",
		},
		{
			Name:     "NetworkSegmentation",
			Keywords: []string{"network segmentation", "subnet", "vpc", "firewall", "microsegmentation", "zero trust", "ztna", "dmz"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(vpc|subnet|network\s+acl|security\s+group)\b`),
				regexp.MustCompile(`(?i)\b(microsegmentation|zero\s+trust|ztna|dmz)\b`),
				regexp.MustCompile(`(?i)\b(east.west|north.south|traffic\s+segmentation)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":    RiskHigh,
				"zero_trust": RiskMedium,
				"flat":       RiskCritical,
			},
			VerificationRules: []string{
				"Verify network segmentation rules enforce least-privilege traffic paths",
				"Confirm microsegmentation policies are defined and monitored",
				"Validate lateral movement is restricted between segments and tiers",
			},
			ExplainabilityTemplate: "Architecture contains network references but does not specify segmentation rules, microsegmentation, or lateral movement controls.",
		},
		{
			Name:     "CloudSecurity",
			Keywords: []string{"cloud", "aws", "azure", "gcp", "cloud security", "shared responsibility", "csp", "cloud provider"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(aws|azure|gcp|cloud\s+(provider|platform|service))\b`),
				regexp.MustCompile(`(?i)\b(shared\s+responsibility|csp|cloud\s+security\s+posture)\b`),
				regexp.MustCompile(`(?i)\b(cloud\s+config|cloud\s+trail|cloud\s+watch|guardduty)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":   RiskMedium,
				"misconfig": RiskCritical,
				"shared":    RiskMedium,
			},
			VerificationRules: []string{
				"Verify cloud security posture management (CSPM) is enabled and reviewed",
				"Confirm shared responsibility boundaries are documented and enforced",
				"Validate cloud IAM policies follow least privilege and are reviewed",
			},
			ExplainabilityTemplate: "Architecture contains cloud resources but does not specify CSPM coverage, shared responsibility boundaries, or IAM review cadence.",
		},
		{
			Name:     "ContainerSecurity",
			Keywords: []string{"container", "docker", "containerd", "image", "runtime", "container security", "image scanning"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(docker|containerd|podman|oci)\b`),
				regexp.MustCompile(`(?i)\b(container\s+image\s+scan|image\s+security|runtime\s+security)\b`),
				regexp.MustCompile(`(?i)\b(seccomp|apparmor|selinux|cap.drop)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":   RiskHigh,
				"scanned":   RiskMedium,
				"unscanned": RiskCritical,
			},
			VerificationRules: []string{
				"Verify container images are scanned for vulnerabilities before deployment",
				"Confirm runtime security profiles (seccomp, AppArmor) are enforced",
				"Validate container registries enforce image signing and immutability",
			},
			ExplainabilityTemplate: "Architecture contains container workloads but does not specify image scanning, runtime security, or registry controls.",
		},
		{
			Name:     "KubernetesSecurity",
			Keywords: []string{"kubernetes", "k8s", "kubectl", "cluster", "pod", "namespace", "service mesh", "helm", "cni"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(kubernetes|k8s|eks|aks|gke|openshift)\b`),
				regexp.MustCompile(`(?i)\b(pod\s+security|network\s+policy|rbac|psp|opa)\b`),
				regexp.MustCompile(`(?i)\b(service\s+mesh|istio|linkerd|cilium)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":      RiskHigh,
				"rbac":         RiskMedium,
				"unrestricted": RiskCritical,
			},
			VerificationRules: []string{
				"Verify Kubernetes RBAC and network policies are defined and enforced",
				"Confirm pod security standards and admission controllers are active",
				"Validate cluster audit logging and control plane hardening",
			},
			ExplainabilityTemplate: "Architecture contains Kubernetes but does not specify RBAC, network policies, pod security, or control plane hardening.",
		},
		{
			Name:     "OperationalSecurity",
			Keywords: []string{"operational security", "opsec", "runbook", "playbook", "sop", "operational procedure", "shift handover"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(runbook|playbook|sop|standard\s+operating\s+procedure)\b`),
				regexp.MustCompile(`(?i)\b(shift\s+handover|change\s+approval|emergency\s+access)\b`),
				regexp.MustCompile(`(?i)\b(security\s+operations|soc|mssp|managed\s+security)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default": RiskMedium,
				"runbook": RiskMedium,
				"missing": RiskHigh,
			},
			VerificationRules: []string{
				"Verify operational security procedures are documented and accessible",
				"Confirm emergency access and break-glass procedures are tested",
				"Validate SOC coverage and escalation paths for security incidents",
			},
			ExplainabilityTemplate: "Architecture contains operational components but does not specify runbooks, emergency access procedures, or SOC coverage.",
		},
		{
			Name:     "Compliance",
			Keywords: []string{"compliance", "regulatory", "framework", "sox", "hipaa", "pci", "gdpr", "iso27001", "nist", "fedramp", "soc2"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(hipaa|pci\s+dss|gdpr|sox|iso\s*27001|nist|fedramp|soc\s*2)\b`),
				regexp.MustCompile(`(?i)\b(compliance\s+(framework|control|audit|gap))\b`),
				regexp.MustCompile(`(?i)\b(regulatory\s+(requirement|obligation|reporting))\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":      RiskHigh,
				"audited":      RiskMedium,
				"noncompliant": RiskCritical,
			},
			VerificationRules: []string{
				"Verify compliance controls are mapped to framework requirements",
				"Confirm compliance audits and gap assessments are scheduled",
				"Validate evidence collection and retention for regulatory examinations",
			},
			ExplainabilityTemplate: "Architecture references regulated data but does not specify compliance mapping, audit scheduling, or evidence retention.",
		},
		{
			Name:     "Governance",
			Keywords: []string{"governance", "policy", "security policy", "risk management", "governance framework", "board", "ciso"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(governance\s+(framework|model|committee)|risk\s+(appetite|register|committee))\b`),
				regexp.MustCompile(`(?i)\b(security\s+council|ciso|board\s+risk|steering\s+committee)\b`),
				regexp.MustCompile(`(?i)\b(policy\s+(review|exception|violation))\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":  RiskMedium,
				"reviewed": RiskLow,
				"stale":    RiskHigh,
			},
			VerificationRules: []string{
				"Verify security policies are reviewed and approved on a defined cadence",
				"Confirm governance committees review risk register and exceptions",
				"Validate policy exceptions are time-bound and risk-assessed",
			},
			ExplainabilityTemplate: "Architecture contains governance references but does not specify policy review cadence, risk committee scope, or exception handling.",
		},
		{
			Name:     "ChangeManagement",
			Keywords: []string{"change management", "change control", "release", "deployment", "rollback", "blue green", "canary"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(change\s+(control|management|approval|advisory|board))\b`),
				regexp.MustCompile(`(?i)\b(blue.green|canary|rollback|feature\s+flag|deployment\s+gate)\b`),
				regexp.MustCompile(`(?i)\b(release\s+(pipeline|automation|approval))\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":   RiskMedium,
				"automated": RiskLow,
				"manual":    RiskHigh,
			},
			VerificationRules: []string{
				"Verify all production changes are tracked and approved before deployment",
				"Confirm rollback and recovery procedures are tested for each change",
				"Validate change windows and emergency change procedures are documented",
			},
			ExplainabilityTemplate: "Architecture contains deployment flows but does not specify change approval, rollback procedures, or deployment gates.",
		},
		{
			Name:     "IncidentResponse",
			Keywords: []string{"incident response", "ir", "breach", "security incident", "forensic", "containment", "eradication", "recovery"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(incident\s+(response|plan|playbook|retainer))\b`),
				regexp.MustCompile(`(?i)\b(breach\s+(notification|response|recovery)|forensic\s+(readiness|investigation))\b`),
				regexp.MustCompile(`(?i)\b(containment|eradication|recovery|lessons\s+learned)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":  RiskHigh,
				"tested":   RiskMedium,
				"untested": RiskCritical,
			},
			VerificationRules: []string{
				"Verify incident response plan is documented and tested annually",
				"Confirm incident communication channels and escalation paths are defined",
				"Validate forensic evidence collection and chain-of-custody procedures",
			},
			ExplainabilityTemplate: "Architecture contains sensitive operations but does not specify incident response plans, testing cadence, or forensic readiness.",
		},
		{
			Name:     "DetectionEngineering",
			Keywords: []string{"detection", "detection engineering", "siem", "soar", "edr", "xdr", "ids", "ips", "threat hunting", "ioc"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(siem|soar|edr|xdr|ids|ips|ndr|mdr)\b`),
				regexp.MustCompile(`(?i)\b(threat\s+hunting|detection\s+(rule|logic|use\s+case))\b`),
				regexp.MustCompile(`(?i)\b(ioc|indicator|ttp|att&ck|mitre)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":  RiskHigh,
				"coverage": RiskMedium,
				"blind":    RiskCritical,
			},
			VerificationRules: []string{
				"Verify detection rules cover all critical attack paths and data flows",
				"Confirm detection tuning and false positive management are performed",
				"Validate detection coverage against MITRE ATT&CK framework",
			},
			ExplainabilityTemplate: "Architecture contains security tooling but does not specify detection coverage, rule tuning, or ATT&CK mapping.",
		},
		{
			Name:     "TrustBoundaries",
			Keywords: []string{"trust boundary", "trust zone", "security boundary", "perimeter", "segmentation", "zero trust"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(trust\s+(boundary|zone)|security\s+(boundary|perimeter))\b`),
				regexp.MustCompile(`(?i)\b(zero\s+trust|implicit\s+trust|trust\s+model)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":   RiskHigh,
				"verified":  RiskMedium,
				"undefined": RiskCritical,
			},
			VerificationRules: []string{
				"Verify all trust boundaries are documented and validated",
				"Confirm no implicit trust exists between segments or tiers",
				"Validate authentication and authorization enforcement at each boundary",
			},
			ExplainabilityTemplate: "Architecture contains multiple tiers but does not specify trust boundaries, zero trust enforcement, or boundary validation.",
		},
		{
			Name:     "HumanProcess",
			Keywords: []string{"human process", "manual process", "social engineering", "phishing", "training", "awareness", "procedure"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(manual\s+(process|step|approval)|human\s+in\s+the\s+loop)\b`),
				regexp.MustCompile(`(?i)\b(phishing|social\s+engineering|security\s+awareness|training)\b`),
				regexp.MustCompile(`(?i)\b(background\s+check|personnel\s+security|insider\s+risk)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":   RiskMedium,
				"trained":   RiskLow,
				"untrained": RiskHigh,
			},
			VerificationRules: []string{
				"Verify security awareness training is mandatory and tracked",
				"Confirm phishing simulation results are reviewed and remediated",
				"Validate background checks and personnel security for sensitive roles",
			},
			ExplainabilityTemplate: "Architecture contains human process dependencies but does not specify training, phishing defenses, or personnel security.",
		},
		{
			Name:     "InsiderThreat",
			Keywords: []string{"insider threat", "insider risk", "malicious insider", "data exfiltration", "user behavior", "ueba", "dlp"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(insider\s+(threat|risk)|malicious\s+insider|data\s+exfiltration)\b`),
				regexp.MustCompile(`(?i)\b(user\s+behavior|ueba|dlp|insider\s+deterrence)\b`),
				regexp.MustCompile(`(?i)\b(privilege\s+abuse|unauthorized\s+access|data\s+theft)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":     RiskHigh,
				"monitored":   RiskMedium,
				"unmonitored": RiskCritical,
			},
			VerificationRules: []string{
				"Verify user behavior analytics and insider threat detection are deployed",
				"Confirm least-privilege access limits data exfiltration blast radius",
				"Validate data loss prevention (DLP) covers endpoints and egress channels",
			},
			ExplainabilityTemplate: "Architecture contains sensitive data access but does not specify insider threat detection, UEBA, or DLP coverage.",
		},
		{
			Name:     "SessionSecurity",
			Keywords: []string{"session", "session security", "token", "jwt", "cookie", "csrf", "session fixation", "session hijack"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(session\s+(management|security|timeout|rotation|fixation|hijack))\b`),
				regexp.MustCompile(`(?i)\b(jwt|oauth\s+token|refresh\s+token|access\s+token)\b`),
				regexp.MustCompile(`(?i)\b(csrf|xsrf|same.site|secure\s+cookie|http.only)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":  RiskHigh,
				"rotation": RiskMedium,
				"missing":  RiskCritical,
			},
			VerificationRules: []string{
				"Verify session tokens are rotated on privilege level change",
				"Confirm session timeout and concurrent session limits are enforced",
				"Validate CSRF protection and secure cookie attributes are present",
			},
			ExplainabilityTemplate: "Architecture contains user sessions but does not specify token rotation, timeout enforcement, or CSRF protection.",
		},
		{
			Name:     "APISecurity",
			Keywords: []string{"api security", "api gateway", "rest", "graphql", "grpc", "openapi", "swagger", "api versioning", "rate limit"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(api\s+(gateway|security|versioning|rate\s+limit)|openapi|swagger)\b`),
				regexp.MustCompile(`(?i)\b(rest|graphql|grpc|webhook|api\s+key)\b`),
				regexp.MustCompile(`(?i)\b(bola|idor|api\s+auth|api\s+validation)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default": RiskHigh,
				"gateway": RiskMedium,
				"exposed": RiskCritical,
			},
			VerificationRules: []string{
				"Verify API authentication and authorization are enforced at the gateway",
				"Confirm rate limiting and input validation are configured for all endpoints",
				"Validate API versioning and deprecation policies are defined and communicated",
			},
			ExplainabilityTemplate: "Architecture contains API endpoints but does not specify gateway security, rate limiting, or input validation.",
		},
		{
			Name:     "ObjectLevelAuthorization",
			Keywords: []string{"object level authorization", "bola", "authorization", "resource access", "data ownership", "row level security", "rls"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(bola|broken\s+object\s+level\s+authorization)\b`),
				regexp.MustCompile(`(?i)\b(object\s+level|resource\s+level|row\s+level|field\s+level)\s+(auth|access|security)\b`),
				regexp.MustCompile(`(?i)\b(data\s+ownership|tenant\s+isolation|multi.tenant\s+auth)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":  RiskCritical,
				"enforced": RiskMedium,
				"missing":  RiskCritical,
			},
			VerificationRules: []string{
				"Verify object-level authorization checks are performed on every request",
				"Confirm resource access is scoped to data ownership and tenant context",
				"Validate BOLA and IDOR testing is included in security assessments",
			},
			ExplainabilityTemplate: "Architecture contains data resources but does not specify object-level authorization, BOLA defenses, or ownership checks.",
		},
		{
			Name:     "BusinessContinuity",
			Keywords: []string{"business continuity", "bcp", "continuity", "continuity of operations", "cob", "critical function", "essential service"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(business\s+continuity|bcp|continuity\s+of\s+operations|cob)\b`),
				regexp.MustCompile(`(?i)\b(critical\s+function|essential\s+service|minimum\s+viable\s+product|mvp)\b`),
				regexp.MustCompile(`(?i)\b(continuity\s+test|exercise|tabletop|walkthrough)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":  RiskHigh,
				"tested":   RiskMedium,
				"untested": RiskCritical,
			},
			VerificationRules: []string{
				"Verify business continuity plan covers all critical functions and dependencies",
				"Confirm continuity tests and exercises are scheduled and documented",
				"Validate RTO and RPO alignment with business impact analysis",
			},
			ExplainabilityTemplate: "Architecture references critical functions but does not specify business continuity planning, testing, or impact analysis.",
		},
		{
			Name:     "DataGovernance",
			Keywords: []string{"data governance", "data quality", "data lineage", "data catalog", "data steward", "master data", "metadata"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(data\s+governance|data\s+quality|data\s+lineage|data\s+catalog)\b`),
				regexp.MustCompile(`(?i)\b(data\s+steward|data\s+owner|master\s+data|metadata\s+management)\b`),
				regexp.MustCompile(`(?i)\b(data\s+dictionary|schema\s+registry|data\s+profiling)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":  RiskMedium,
				"enforced": RiskLow,
				"missing":  RiskHigh,
			},
			VerificationRules: []string{
				"Verify data ownership and stewardship roles are assigned and active",
				"Confirm data lineage and impact analysis are available for critical data",
				"Validate data quality rules and monitoring are enforced at ingestion",
			},
			ExplainabilityTemplate: "Architecture references data flows but does not specify data governance, lineage, quality rules, or stewardship.",
		},
		{
			Name:     "EncryptionAtRest",
			Keywords: []string{"encryption at rest", "data at rest", "disk encryption", "volume encryption", "database encryption", "storage encryption"},
			Patterns: []*regexp.Regexp{
				regexp.MustCompile(`(?i)\b(encryption\s+at\s+rest|data\s+at\s+rest)\b`),
				regexp.MustCompile(`(?i)\b(disk\s+encryption|volume\s+encryption|storage\s+encryption|db\s+encryption)\b`),
				regexp.MustCompile(`(?i)\b(tde|transparent\s+data\s+encryption|bitlocker|luks)\b`),
			},
			RiskMappings: map[string]RiskLevel{
				"default":  RiskHigh,
				"enforced": RiskMedium,
				"missing":  RiskCritical,
			},
			VerificationRules: []string{
				"Verify all persistent storage volumes use approved encryption at rest",
				"Confirm encryption keys for at-rest data are managed separately from data",
				"Validate encryption settings are audited and reported for compliance",
			},
			ExplainabilityTemplate: "Architecture contains persistent storage but does not specify encryption at rest, key separation, or compliance audit coverage.",
		},
	}

	for _, cat := range categories {
		te.Categories[cat.Name] = cat
	}
}

// Count returns the number of registered categories.
func (te *TaxonomyEngine) Count() int {
	return len(te.Categories)
}

// AllCategories returns a slice of all category names.
func (te *TaxonomyEngine) AllCategories() []string {
	var names []string
	for name := range te.Categories {
		names = append(names, name)
	}
	return names
}
