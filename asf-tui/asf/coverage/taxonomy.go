package coverage

import "strings"

type componentRule struct {
	components   []string
	expectations []ExpectedAssumption
}

var taxonomyRules = []componentRule{
	{
		components: []string{"auth0", "okta", "keycloak", "idp", "identity", "sso"},
		expectations: []ExpectedAssumption{
			{CatIdentity, "MFA Enforcement", "Multi-factor authentication is enforced for all users", "Critical", "identity"},
			{CatIdentity, "SSO Configuration", "Single sign-on is configured across all applications", "High", "identity"},
			{CatIdentity, "Federation", "Identity federation is configured with trusted providers", "High", "identity"},
			{CatIdentity, "Admin Access Restriction", "Administrative access is restricted to authorized users", "Critical", "identity"},
			{CatAuthorization, "RBAC Configuration", "Role-based access control is configured", "Critical", "authorization"},
			{CatAuthorization, "Least Privilege", "Least privilege principle is enforced", "High", "authorization"},
			{CatMonitoring, "Access Audit Logging", "Access events are logged and monitored", "High", "monitoring"},
		},
	},
	{
		components: []string{"database", "db", "rds", "dynamodb", "cosmosdb", "bigquery"},
		expectations: []ExpectedAssumption{
			{CatAuthorization, "Database Access Control", "Database access is restricted to authorized users", "Critical", "authorization"},
			{CatCryptography, "Encryption at Rest", "Data is encrypted at rest", "Critical", "cryptography"},
			{CatCryptography, "Encryption in Transit", "Data is encrypted in transit", "Critical", "cryptography"},
			{CatMonitoring, "Database Audit Logging", "Database access and modifications are audited", "High", "monitoring"},
			{CatResilience, "Database Backups", "Database backups are configured", "High", "resilience"},
			{CatResilience, "Restore Testing", "Backup restore procedures are tested", "High", "resilience"},
		},
	},
	{
		components: []string{"kms", "key", "vault", "hsm", "keyvault", "secrets manager"},
		expectations: []ExpectedAssumption{
			{CatCryptography, "Key Rotation", "Encryption keys are rotated automatically", "Critical", "cryptography"},
			{CatCryptography, "Key Access Control", "Key access is restricted to authorized services", "Critical", "cryptography"},
			{CatCryptography, "Key Backup", "Key backup and recovery procedures exist", "High", "cryptography"},
			{CatAuthorization, "KMS Access Policy", "KMS access is controlled via IAM policies", "Critical", "authorization"},
			{CatMonitoring, "Key Usage Auditing", "Key usage is logged and monitored", "High", "monitoring"},
		},
	},
	{
		components: []string{"siem", "splunk", "elastic", "log", "audit", "cloudtrail", "cloudwatch"},
		expectations: []ExpectedAssumption{
			{CatMonitoring, "Centralized Logging", "Logs are collected in a central SIEM system", "Critical", "monitoring"},
			{CatMonitoring, "Alerting", "Security alerts are configured for critical events", "Critical", "monitoring"},
			{CatMonitoring, "Log Retention", "Logs are retained per compliance requirements", "High", "monitoring"},
			{CatMonitoring, "Log Integrity", "Logs are immutable and tamper-proof", "High", "monitoring"},
		},
	},
	{
		components: []string{"apigateway", "api", "gateway", "kong", "apigee", "ambassador"},
		expectations: []ExpectedAssumption{
			{CatAuthorization, "API Authentication", "API requests are authenticated", "Critical", "authorization"},
			{CatAuthorization, "API Authorization", "API requests are authorized per policy", "Critical", "authorization"},
			{CatMonitoring, "API Audit Logging", "API access is logged", "High", "monitoring"},
			{CatOperational, "Rate Limiting", "API rate limiting is configured", "High", "operational"},
		},
	},
	{
		components: []string{"s3", "bucket", "storage", "blob", "object storage"},
		expectations: []ExpectedAssumption{
			{CatAuthorization, "Bucket Access Control", "Storage bucket access is restricted", "Critical", "authorization"},
			{CatCryptography, "Encryption at Rest", "Data in storage is encrypted at rest", "Critical", "cryptography"},
			{CatMonitoring, "Access Logging", "Storage access is logged", "High", "monitoring"},
			{CatResilience, "Data Backup", "Storage data is backed up", "High", "resilience"},
		},
	},
	{
		components: []string{"backup", "backupservice", "snapshot", "restore"},
		expectations: []ExpectedAssumption{
			{CatResilience, "Backup Schedule", "Backups are performed on a regular schedule", "High", "resilience"},
			{CatResilience, "Restore Testing", "Restore procedures are tested periodically", "Critical", "resilience"},
			{CatResilience, "Backup Encryption", "Backups are encrypted at rest and in transit", "High", "resilience"},
			{CatResilience, "Disaster Recovery", "Disaster recovery plan exists and is tested", "Critical", "resilience"},
		},
	},
	{
		components: []string{"jenkins", "ci/cd", "github actions", "gitlab ci", "pipeline", "circleci"},
		expectations: []ExpectedAssumption{
			{CatOperational, "Secrets Management", "CI/CD secrets are managed securely", "Critical", "operational"},
			{CatOperational, "Pipeline Security", "CI/CD pipeline is hardened against tampering", "Critical", "operational"},
			{CatOperational, "Artifact Signing", "Build artifacts are signed and verified", "High", "operational"},
		},
	},
	{
		components: []string{"webapp", "app", "application", "service"},
		expectations: []ExpectedAssumption{
			{CatAuthorization, "Authentication", "Application authentication is enforced", "Critical", "authorization"},
			{CatAuthorization, "Session Management", "User sessions are managed securely", "High", "authorization"},
			{CatCryptography, "TLS Configuration", "TLS is configured for all endpoints", "Critical", "cryptography"},
		},
	},
}

var domainBlindSpotRules = map[string][]DomainBlindSpot{
	"healthcare": {
		{"healthcare", "Break Glass Access", "Emergency break-glass access procedure for PHI access", "Critical", "Define and test break-glass access procedures"},
		{"healthcare", "PHI Audit Review", "Regular review of PHI access audit logs", "High", "Establish PHI audit review cadence"},
		{"healthcare", "Patient Data Retention", "Patient data retention and disposal policy", "High", "Define data retention policy per HIPAA requirements"},
	},
	"hipaa": {
		{"hipaa", "Break Glass Access", "Emergency break-glass access procedure for PHI access", "Critical", "Define and test break-glass access procedures"},
		{"hipaa", "PHI Audit Review", "Regular review of PHI access audit logs", "High", "Establish PHI audit review cadence"},
	},
	"fintech": {
		{"fintech", "Fraud Monitoring", "Real-time fraud detection and monitoring", "Critical", "Implement fraud detection system"},
		{"fintech", "Key Custody", "Cryptographic key custody and governance", "Critical", "Establish key custody procedures"},
		{"fintech", "Settlement Controls", "Transaction settlement verification controls", "High", "Implement settlement reconciliation"},
	},
	"cloud": {
		{"cloud", "IAM Review", "Regular IAM policy review and cleanup", "Critical", "Schedule periodic IAM access reviews"},
		{"cloud", "Federation Controls", "Cloud federation trust relationships are monitored", "High", "Review federation trust configurations"},
		{"cloud", "Secrets Rotation", "Cloud secrets and credentials are rotated", "Critical", "Implement automated secrets rotation"},
	},
	"kubernetes": {
		{"kubernetes", "Admission Controls", "Kubernetes admission controllers are configured", "Critical", "Configure OPA/Gatekeeper admission controls"},
		{"kubernetes", "Service Account Governance", "Kubernetes service accounts follow least privilege", "High", "Review and restrict service account permissions"},
		{"kubernetes", "Secret Management", "Kubernetes secrets are encrypted at rest", "Critical", "Enable KMS encryption for secrets"},
	},
	"vpn": {
		{"vpn", "Certificate Management", "VPN certificates are managed and rotated", "Critical", "Implement certificate lifecycle management"},
		{"vpn", "Split Tunnel Control", "Split tunneling is restricted per policy", "High", "Review split tunnel configuration"},
	},
}

func GetExpectations(component string) []ExpectedAssumption {
	cl := strings.ToLower(component)
	for _, rule := range taxonomyRules {
		for _, c := range rule.components {
			if strings.Contains(cl, c) {
				return rule.expectations
			}
		}
	}
	return nil
}

func GetComponentCategories(component string) []CoverageCategory {
	expectations := GetExpectations(component)
	seen := make(map[CoverageCategory]bool)
	var cats []CoverageCategory
	for _, e := range expectations {
		if !seen[e.Category] {
			seen[e.Category] = true
			cats = append(cats, e.Category)
		}
	}
	return cats
}

func GetDomainBlindSpots(domain string) []DomainBlindSpot {
	dl := strings.ToLower(domain)
	rules, ok := domainBlindSpotRules[dl]
	if !ok {
		for key, val := range domainBlindSpotRules {
			if strings.Contains(dl, key) {
				return val
			}
		}
	}
	return rules
}
