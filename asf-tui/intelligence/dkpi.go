package intelligence

import (
	"fmt"
	"sort"
	"strings"
)

// ─────────────────────────────────────────────────────────────
// DKPI — Domain Knowledge Pack Intelligence Engine (ASF V8)
// Phases 1-15
// ─────────────────────────────────────────────────────────────

// ── PHASE 1 — KNOWLEDGE PACK FRAMEWORK ──

// KnowledgePack is a domain-specific bundle of security knowledge.
type KnowledgePack struct {
	ID                   string                    `json:"id"`
	Name                 string                    `json:"name"`
	Industry             string                    `json:"industry"`
	Version              string                    `json:"version"`
	Description          string                    `json:"description"`
	DetectionKeywords    []string                  `json:"detection_keywords"`
	DetectionComponents  []string                  `json:"detection_components"`
	DetectionCompliance  []string                  `json:"detection_compliance"`
	CrownJewels          []string                  `json:"crown_jewels"`
	ExpectedControls     []KnowledgePackControl    `json:"expected_controls"`
	ExpectedEvidence     []KnowledgePackEvidence   `json:"expected_evidence"`
	ThreatPatterns       []KnowledgePackThreat     `json:"threat_patterns"`
	AttackPathTemplates  []KnowledgePackAttackPath `json:"attack_path_templates"`
	AssumptionPatterns   []KnowledgePackAssumption `json:"assumption_patterns"`
	ComplianceFrameworks []string                  `json:"compliance_frameworks"`
}

type KnowledgePackControl struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Priority    string `json:"priority"`
}

type KnowledgePackEvidence struct {
	Control  string   `json:"control"`
	Evidence []string `json:"evidence"`
}

type KnowledgePackThreat struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Severity    string `json:"severity"`
	Category    string `json:"category"`
	Target      string `json:"target"`
}

type KnowledgePackAttackPath struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Steps       []string `json:"steps"`
	Target      string   `json:"target"`
}

type KnowledgePackAssumption struct {
	Description string   `json:"description"`
	Category    string   `json:"category"`
	Risk        string   `json:"risk"`
	Component   string   `json:"component"`
	Keywords    []string `json:"keywords"`
}

// ── PHASE 2 — DOMAIN DETECTION ──

type DomainDetectionResult struct {
	PrimaryDomain string            `json:"primary_domain"`
	Confidence    float64           `json:"confidence"`
	Rationale     []string          `json:"rationale"`
	Matches       []DomainPackMatch `json:"matches"`
}

type DomainPackMatch struct {
	PackID     string   `json:"pack_id"`
	PackName   string   `json:"pack_name"`
	Score      int      `json:"score"`
	MaxScore   int      `json:"max_score"`
	Confidence float64  `json:"confidence"`
	Reasons    []string `json:"reasons"`
}

// DKPIResult is the output of the DKPI engine.
type DKPIResult struct {
	DetectedDomain  DomainDetectionResult   `json:"detected_domain"`
	ActivePack      *KnowledgePack          `json:"active_pack"`
	DomainFindings  []SDRIFinding           `json:"domain_findings"`
	Recommendations []string                `json:"recommendations"`
	EvidenceReqs    []KnowledgePackEvidence `json:"evidence_requirements"`
	InjectedThreats []Threat                `json:"injected_threats"`
}

type DKPIDetector struct{}

func NewDKPIDetector() *DKPIDetector {
	return &DKPIDetector{}
}

// DetectDomain performs auto-domain detection from architecture data.
func (d *DKPIDetector) DetectDomain(arch *ArchDescription) DomainDetectionResult {
	if arch == nil {
		return DomainDetectionResult{PrimaryDomain: "", Confidence: 0}
	}

	packs := buildKnowledgePacks()
	raw := strings.ToLower(arch.RawText)
	matches := make([]DomainPackMatch, 0)

	for _, pack := range packs {
		score := 0
		maxScore := 0
		reasons := make([]string, 0)

		// Component-based scoring
		for _, comp := range arch.Components {
			label := strings.ToLower(comp.Label)
			for _, kw := range pack.DetectionComponents {
				maxScore += 3
				if strings.Contains(label, strings.ToLower(kw)) {
					score += 3
					reasons = append(reasons, fmt.Sprintf("Component %q matched %q", comp.Label, kw))
				}
			}
		}

		// Keyword-based scoring (from raw text)
		for _, kw := range pack.DetectionKeywords {
			maxScore += 2
			if strings.Contains(raw, strings.ToLower(kw)) {
				score += 2
				reasons = append(reasons, fmt.Sprintf("Text matched keyword %q", kw))
			}
		}

		// Compliance-based scoring
		for _, c := range arch.Compliance {
			for _, pc := range pack.DetectionCompliance {
				maxScore += 5
				if strings.EqualFold(c, pc) {
					score += 5
					reasons = append(reasons, fmt.Sprintf("Compliance %q matched pack", c))
				}
			}
		}

		// Policy-based scoring
		for _, policy := range arch.Policies {
			pl := strings.ToLower(policy)
			for _, kw := range pack.DetectionKeywords {
				maxScore += 1
				if strings.Contains(pl, kw) {
					score += 1
					break
				}
			}
		}

		confidence := 0.0
		if maxScore > 0 {
			confidence = float64(score) / float64(maxScore) * 100
		}

		if score > 0 {
			matches = append(matches, DomainPackMatch{
				PackID:     pack.ID,
				PackName:   pack.Name,
				Score:      score,
				MaxScore:   maxScore,
				Confidence: confidence,
				Reasons:    reasons,
			})
		}
	}

	// Sort by score descending (alphabetically for ties)
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].Score != matches[j].Score {
			return matches[i].Score > matches[j].Score
		}
		return matches[i].PackID < matches[j].PackID
	})

	primaryDomain := ""
	confidence := 0.0
	rationale := make([]string, 0)

	if len(matches) > 0 {
		primaryDomain = matches[0].PackID
		confidence = matches[0].Confidence
		rationale = append(rationale, fmt.Sprintf("Top match: %s (score %d/%d, %.0f%% confidence)", matches[0].PackName, matches[0].Score, matches[0].MaxScore, confidence))
		if len(matches[0].Reasons) > 0 {
			rationale = append(rationale, matches[0].Reasons[:minInt(3, len(matches[0].Reasons))]...)
		}
	}

	if confidence < 10 {
		primaryDomain = ""
		rationale = append(rationale, "Confidence too low for reliable domain detection")
	}

	return DomainDetectionResult{
		PrimaryDomain: primaryDomain,
		Confidence:    confidence,
		Rationale:     rationale,
		Matches:       matches,
	}
}

// GetPack retrieves a knowledge pack by ID.
func (d *DKPIDetector) GetPack(id string) *KnowledgePack {
	for _, p := range buildKnowledgePacks() {
		if p.ID == id {
			return &p
		}
	}
	return nil
}

// ── PHASE 2-8 — KNOWLEDGE PACK DEFINITIONS ──

func buildKnowledgePacks() []KnowledgePack {
	return []KnowledgePack{
		healthcarePack(),
		fintechPack(),
		kubernetesPack(),
		saasPack(),
		governmentPack(),
		criticalInfrastructurePack(),
	}
}

func healthcarePack() KnowledgePack {
	return KnowledgePack{
		ID:                   "healthcare",
		Name:                 "Healthcare",
		Industry:             "Healthcare",
		Version:              "1.0",
		Description:          "Healthcare domain knowledge pack for architectures handling PHI, patient records, and clinical systems.",
		DetectionKeywords:    []string{"phi", "hipaa", "patient", "health", "medical", "ehr", "clinical", "hl7", "fda", "clinic", "hospital", "emr", "protected health", "healthcare", "clinician"},
		DetectionComponents:  []string{"phi", "patient", "health", "medical", "ehr", "clinic", "hospital", "emr"},
		DetectionCompliance:  []string{"HIPAA", "HITRUST"},
		CrownJewels:          []string{"PHI Database", "EHR System", "Patient Records", "Clinical Data Store"},
		ComplianceFrameworks: []string{"HIPAA", "HITRUST", "FDA 21 CFR Part 11"},
		ExpectedControls: []KnowledgePackControl{
			{Name: "BreakGlassAccess", Description: "Emergency break-glass access for PHI with immutable audit logging", Category: "PrivilegedAccess", Priority: "Critical"},
			{Name: "PHIAuditLogging", Description: "Immutable audit logging for all PHI access and modification events", Category: "AuditLogging", Priority: "Critical"},
			{Name: "DataEncryptionAtRest", Description: "Encryption of PHI at rest using FIPS 140-2 validated algorithms", Category: "DataProtection", Priority: "Critical"},
			{Name: "DataEncryptionInTransit", Description: "Encryption of PHI in transit over all communication channels", Category: "DataProtection", Priority: "Critical"},
			{Name: "AccessReviews", Description: "Periodic reviews of PHI access privileges", Category: "IdentityGovernance", Priority: "High"},
			{Name: "MinimumNecessaryAccess", Description: "Minimum necessary access controls for all PHI interactions", Category: "Authorization", Priority: "Critical"},
			{Name: "PatientDataExport", Description: "Controls for patient data export and portability", Category: "DataProtection", Priority: "High"},
			{Name: "MedicalDeviceTrust", Description: "Trust validation for connected medical devices", Category: "DeviceSecurity", Priority: "High"},
		},
		ExpectedEvidence: []KnowledgePackEvidence{
			{Control: "BreakGlassAccess", Evidence: []string{"Break glass policy document", "Emergency access logs", "Break glass usage reports"}},
			{Control: "PHIAuditLogging", Evidence: []string{"PHI audit log configuration", "Audit log review procedures", "Log retention policy"}},
			{Control: "DataEncryptionAtRest", Evidence: []string{"Encryption policy", "Key management documentation", "Encryption configuration audits"}},
			{Control: "DataEncryptionInTransit", Evidence: []string{"TLS configuration", "Certificate inventory", "Network encryption scan reports"}},
			{Control: "AccessReviews", Evidence: []string{"Access review schedule", "Completed review reports", "Remediation tracking reports"}},
			{Control: "MinimumNecessaryAccess", Evidence: []string{"Access control policy", "Role definitions", "User permission audits"}},
		},
		ThreatPatterns: []KnowledgePackThreat{
			{Name: "PHI Theft via Compromised Identity", Description: "Attacker compromises a clinical identity to exfiltrate PHI from EHR systems", Severity: "Critical", Category: "Information Disclosure", Target: "PHI Database"},
			{Name: "Ransomware on Healthcare Systems", Description: "Ransomware encrypts clinical data stores and disrupts patient care operations", Severity: "Critical", Category: "Denial of Service", Target: "EHR System"},
			{Name: "Insider PHI Access Abuse", Description: "Insider with legitimate access exfiltrates PHI for personal gain", Severity: "High", Category: "Information Disclosure", Target: "Patient Records"},
			{Name: "Medical Device Compromise", Description: "Attacker compromises connected medical device to pivot to clinical network", Severity: "Critical", Category: "Elevation of Privilege", Target: "Medical Device"},
		},
		AttackPathTemplates: []KnowledgePackAttackPath{
			{Name: "Internet to PHI Exfiltration", Description: "External attacker compromises identity provider, accesses EHR, exfiltrates PHI", Steps: []string{"Internet", "Identity Provider", "API Gateway", "EHR System", "PHI Database"}, Target: "PHI Database"},
			{Name: "Insider to Patient Record Access", Description: "Insider with clinical privileges accesses patient records outside need-to-know", Steps: []string{"Insider Workstation", "Clinical App", "Patient Records"}, Target: "Patient Records"},
		},
		AssumptionPatterns: []KnowledgePackAssumption{
			{Description: "Healthcare system must ensure PHI access is audited and meets HIPAA audit controls (164.312(b))", Category: "Auditability", Risk: "Critical", Component: "PHI", Keywords: []string{"phi", "hipaa", "audit", "healthcare"}},
			{Description: "Healthcare system must implement break-glass access for emergency PHI retrieval with immutable audit logging", Category: "PrivilegedAccess", Risk: "Critical", Component: "PHI", Keywords: []string{"phi", "break-glass", "emergency", "audit"}},
			{Description: "Healthcare system must encrypt PHI at rest and in transit with approved algorithms", Category: "DataProtection", Risk: "Critical", Component: "PHI", Keywords: []string{"phi", "encryption", "hipaa"}},
			{Description: "Healthcare system must enforce minimum necessary access for all PHI interactions", Category: "Authorization", Risk: "Critical", Component: "PHI", Keywords: []string{"phi", "minimum necessary", "access"}},
		},
	}
}

func fintechPack() KnowledgePack {
	return KnowledgePack{
		ID:                   "fintech",
		Name:                 "Financial Technology",
		Industry:             "Financial Services",
		Version:              "1.0",
		Description:          "Fintech domain knowledge pack for architectures handling payments, transactions, and financial data.",
		DetectionKeywords:    []string{"pci", "payment", "card", "settlement", "fraud", "banking", "transaction", "aml", "kyc", "fintech", "swift", "ach", "sox", "financial", "trading"},
		DetectionComponents:  []string{"payment", "card", "bank", "transaction", "fraud", "settlement"},
		DetectionCompliance:  []string{"PCI-DSS", "PCI DSS", "SOX"},
		CrownJewels:          []string{"Payment Processor", "Settlement System", "Cardholder Data", "Transaction Database"},
		ComplianceFrameworks: []string{"PCI DSS", "SOX", "PSD2", "FFIEC", "GDPR"},
		ExpectedControls: []KnowledgePackControl{
			{Name: "FraudMonitoring", Description: "Real-time fraud monitoring for all payment transactions", Category: "DetectionEngineering", Priority: "Critical"},
			{Name: "DualAuthorization", Description: "Dual authorization for high-value transactions and settlement operations", Category: "AccessControl", Priority: "Critical"},
			{Name: "TransactionIntegrity", Description: "Transaction integrity verification with reconciliation", Category: "DataIntegrity", Priority: "Critical"},
			{Name: "PCIKeyManagement", Description: "HSM-backed key management for payment encryption keys", Category: "KeyManagement", Priority: "Critical"},
			{Name: "SettlementValidation", Description: "Automated settlement validation and reconciliation", Category: "DataIntegrity", Priority: "Critical"},
			{Name: "TransactionLogging", Description: "Immutable transaction logging with tamper-evident controls", Category: "AuditLogging", Priority: "Critical"},
			{Name: "PCIScopeControls", Description: "PCI DSS scope reduction with CDE segmentation", Category: "NetworkSegmentation", Priority: "Critical"},
			{Name: "CardholderEncryption", Description: "Encryption of cardholder data at rest and in transit", Category: "DataProtection", Priority: "Critical"},
		},
		ExpectedEvidence: []KnowledgePackEvidence{
			{Control: "FraudMonitoring", Evidence: []string{"Fraud detection rules configuration", "Fraud incident reports", "Monitoring dashboards"}},
			{Control: "DualAuthorization", Evidence: []string{"Dual control policy", "Transaction approval logs", "Dual authorization audit records"}},
			{Control: "TransactionIntegrity", Evidence: []string{"Reconciliation reports", "Exception handling logs", "Integrity verification procedures"}},
			{Control: "PCIKeyManagement", Evidence: []string{"HSM configuration", "Key rotation records", "Key access audit logs"}},
			{Control: "SettlementValidation", Evidence: []string{"Settlement reports", "Reconciliation logs", "Discrepancy resolution records"}},
			{Control: "TransactionLogging", Evidence: []string{"Transaction log configuration", "Log retention policy", "Log monitoring setup"}},
		},
		ThreatPatterns: []KnowledgePackThreat{
			{Name: "Payment Transaction Fraud", Description: "Attacker manipulates payment transactions to divert funds", Severity: "Critical", Category: "Tampering", Target: "Payment Processor"},
			{Name: "Settlement Manipulation", Description: "Attacker alters settlement records to conceal fraudulent transactions", Severity: "Critical", Category: "Tampering", Target: "Settlement System"},
			{Name: "Cardholder Data Theft", Description: "Attacker exfiltrates cardholder data from payment systems", Severity: "Critical", Category: "Information Disclosure", Target: "Cardholder Data"},
			{Name: "PCI Scope Escape", Description: "Attacker pivots from non-CDE systems into the cardholder data environment", Severity: "Critical", Category: "Elevation of Privilege", Target: "CDE"},
		},
		AttackPathTemplates: []KnowledgePackAttackPath{
			{Name: "Internet to Payment Fraud", Description: "External attacker exploits web application to manipulate payment transactions", Steps: []string{"Internet", "Web App", "API", "Payment Processor", "Settlement"}, Target: "Settlement System"},
			{Name: "Insider to Cardholder Data", Description: "Insider with non-CDE access pivots into cardholder data environment", Steps: []string{"Internal Network", "Non-CDE App", "CDE Boundary", "Cardholder Database"}, Target: "Cardholder Data"},
		},
		AssumptionPatterns: []KnowledgePackAssumption{
			{Description: "Fintech system must enforce PCI DSS requirements for cardholder data environment", Category: "Compliance", Risk: "Critical", Component: "Payment", Keywords: []string{"pci", "payment", "cardholder"}},
			{Description: "Fintech system must implement fraud detection and transaction integrity monitoring", Category: "Detection", Risk: "Critical", Component: "Payment", Keywords: []string{"fraud", "transaction integrity"}},
			{Description: "Fintech system must implement dual authorization for settlement operations", Category: "AccessControl", Risk: "Critical", Component: "Settlement", Keywords: []string{"dual authorization", "settlement"}},
		},
	}
}

func kubernetesPack() KnowledgePack {
	return KnowledgePack{
		ID:                   "kubernetes",
		Name:                 "Kubernetes / Cloud Native",
		Industry:             "Cloud Native Infrastructure",
		Version:              "1.0",
		Description:          "Kubernetes domain knowledge pack for container orchestration platforms.",
		DetectionKeywords:    []string{"kubernetes", "k8s", "container", "pod", "namespace", "helm", "operator", "cri", "cni", "kube", "docker", "registry", "service mesh"},
		DetectionComponents:  []string{"kubernetes", "k8s", "pod", "cluster", "helm", "container", "namespace"},
		DetectionCompliance:  []string{"CIS"},
		CrownJewels:          []string{"Cluster Admin", "Secrets Store", "Container Registry", "etcd"},
		ComplianceFrameworks: []string{"CIS Kubernetes Benchmark", "NIST SP 800-190", "SOC 2", "PCI DSS"},
		ExpectedControls: []KnowledgePackControl{
			{Name: "K8sRBAC", Description: "Kubernetes RBAC with least-privilege service accounts", Category: "KubernetesSecurity", Priority: "Critical"},
			{Name: "K8sAdmissionControllers", Description: "Admission controller policies for pod security", Category: "KubernetesSecurity", Priority: "Critical"},
			{Name: "K8sNetworkPolicies", Description: "Network policies to restrict pod-to-pod traffic", Category: "NetworkSegmentation", Priority: "Critical"},
			{Name: "K8sSecretsProtection", Description: "External secrets management with vault integration", Category: "SecretsManagement", Priority: "Critical"},
			{Name: "ContainerImageScanning", Description: "Vulnerability scanning of container images in CI/CD pipeline", Category: "ContainerSecurity", Priority: "High"},
			{Name: "ContainerRuntimeSecurity", Description: "Runtime security monitoring for containers", Category: "ContainerSecurity", Priority: "High"},
			{Name: "K8sPodSecurity", Description: "Pod security standards and security contexts", Category: "KubernetesSecurity", Priority: "Critical"},
			{Name: "K8sAuditLogging", Description: "Kubernetes API server audit logging", Category: "AuditLogging", Priority: "High"},
		},
		ExpectedEvidence: []KnowledgePackEvidence{
			{Control: "K8sRBAC", Evidence: []string{"RBAC configuration files", "Service account inventory", "RBAC audit reports"}},
			{Control: "K8sAdmissionControllers", Evidence: []string{"Admission controller configuration", "Policy definitions", "Admission audit logs"}},
			{Control: "K8sNetworkPolicies", Evidence: []string{"Network policy manifests", "Traffic flow documentation", "Network policy audit"}},
			{Control: "K8sSecretsProtection", Evidence: []string{"Vault configuration", "Secrets rotation policy", "Access policies"}},
			{Control: "ContainerImageScanning", Evidence: []string{"Scan reports", "Image approval process", "Vulnerability remediation tracking"}},
			{Control: "ContainerRuntimeSecurity", Evidence: []string{"Runtime security configuration", "Security event logs", "Alert configuration"}},
		},
		ThreatPatterns: []KnowledgePackThreat{
			{Name: "Container Escape", Description: "Attacker escapes container to access host system and other containers", Severity: "Critical", Category: "Elevation of Privilege", Target: "Container Host"},
			{Name: "Cluster Admin Compromise", Description: "Attacker gains cluster admin via compromised service account or RBAC misconfiguration", Severity: "Critical", Category: "Elevation of Privilege", Target: "Cluster Admin"},
			{Name: "Secrets Theft from etcd", Description: "Attacker accesses etcd to steal secrets stored in plaintext", Severity: "Critical", Category: "Information Disclosure", Target: "etcd"},
			{Name: "Supply Chain Attack via Compromised Image", Description: "Attacker compromises container image with malware in registry", Severity: "Critical", Category: "Tampering", Target: "Container Registry"},
		},
		AttackPathTemplates: []KnowledgePackAttackPath{
			{Name: "Internet to Cluster Admin", Description: "External attacker exploits web app pod, escapes to node, escalates to cluster admin", Steps: []string{"Internet", "Ingress", "Web App Pod", "Node Compromise", "Cluster Admin"}, Target: "Cluster Admin"},
			{Name: "Compromised Image to Data Exfiltration", Description: "Attacker injects malicious image pulled from registry, exfiltrates data from pods", Steps: []string{"Compromised Image", "Container Registry", "Pod Deployment", "Data Store", "Exfiltration"}, Target: "Data Store"},
		},
		AssumptionPatterns: []KnowledgePackAssumption{
			{Description: "Kubernetes architecture must enforce pod security standards and admission controller policies", Category: "ContainerSecurity", Risk: "Critical", Component: "Kubernetes", Keywords: []string{"kubernetes", "pod security", "admission controller"}},
			{Description: "Kubernetes architecture must enforce RBAC with least-privilege service accounts", Category: "KubernetesSecurity", Risk: "Critical", Component: "Kubernetes", Keywords: []string{"kubernetes", "rbac", "service account"}},
			{Description: "Kubernetes architecture must implement network policies to restrict pod-to-pod traffic", Category: "NetworkSegmentation", Risk: "High", Component: "Kubernetes", Keywords: []string{"kubernetes", "network policy"}},
			{Description: "Kubernetes architecture must implement external secrets management with no plaintext secrets in manifests", Category: "SecretsManagement", Risk: "Critical", Component: "Kubernetes", Keywords: []string{"kubernetes", "secrets", "vault"}},
		},
	}
}

func saasPack() KnowledgePack {
	return KnowledgePack{
		ID:                   "saas",
		Name:                 "SaaS / Cloud Service",
		Industry:             "Software as a Service",
		Version:              "1.0",
		Description:          "SaaS domain knowledge pack for multi-tenant cloud service architectures.",
		DetectionKeywords:    []string{"multi-tenant", "tenant", "saas", "subscription", "onboarding", "sso", "auth0", "okta", "stripe", "s3", "api gateway", "cloud storage"},
		DetectionComponents:  []string{"tenant", "saas", "api gateway", "auth", "storage"},
		DetectionCompliance:  []string{"SOC2", "SOC 2", "ISO27001"},
		CrownJewels:          []string{"Tenant Data", "User Database", "API Gateway", "Identity Provider"},
		ComplianceFrameworks: []string{"SOC 2", "ISO 27001", "GDPR", "CCPA"},
		ExpectedControls: []KnowledgePackControl{
			{Name: "TenantIsolation", Description: "Tenant isolation at network, data, and compute layers", Category: "TrustBoundaries", Priority: "Critical"},
			{Name: "DataSegregation", Description: "Row-level security and data segregation per tenant", Category: "DataProtection", Priority: "Critical"},
			{Name: "APIRateLimiting", Description: "Per-tenant API rate limiting and throttling", Category: "APISecurity", Priority: "High"},
			{Name: "BOLAAuthorization", Description: "Object-level authorization to prevent broken object level access", Category: "APISecurity", Priority: "Critical"},
			{Name: "TenantAuditLogging", Description: "Per-tenant audit logging for compliance and troubleshooting", Category: "AuditLogging", Priority: "High"},
			{Name: "IdentityGovernance", Description: "Identity lifecycle management across tenants", Category: "Identity", Priority: "High"},
			{Name: "VendorAccessControls", Description: "Vendor and third-party access controls for SaaS operations", Category: "ThirdParty", Priority: "High"},
			{Name: "AvailabilityMonitoring", Description: "Multi-tenant availability monitoring and SLA tracking", Category: "Monitoring", Priority: "High"},
		},
		ExpectedEvidence: []KnowledgePackEvidence{
			{Control: "TenantIsolation", Evidence: []string{"Tenant isolation architecture diagram", "Isolation test results", "Network segmentation config"}},
			{Control: "DataSegregation", Evidence: []string{"Row-level security configuration", "Data segregation audit", "Tenant boundary tests"}},
			{Control: "APIRateLimiting", Evidence: []string{"Rate limiting configuration", "Throttling policies", "Rate limit monitoring"}},
			{Control: "BOLAAuthorization", Evidence: []string{"Authorization policy", "Object-level access logs", "BOLA vulnerability scan reports"}},
		},
		ThreatPatterns: []KnowledgePackThreat{
			{Name: "Cross-Tenant Data Access", Description: "Attacker exploits tenant isolation gap to access another tenant's data", Severity: "Critical", Category: "Information Disclosure", Target: "Tenant Data"},
			{Name: "API Abuse via Compromised Tenant", Description: "Compromised tenant API key used to overwhelm system resources", Severity: "High", Category: "Denial of Service", Target: "API Gateway"},
			{Name: "Broken Object Level Authorization", Description: "Attacker manipulates object IDs to access unauthorized tenant resources", Severity: "Critical", Category: "Information Disclosure", Target: "User Database"},
		},
		AttackPathTemplates: []KnowledgePackAttackPath{
			{Name: "Compromised Tenant to Cross-Tenant Access", Description: "Attacker compromises one tenant's credentials then exploits isolation gaps", Steps: []string{"Internet", "Compromised Tenant", "API Gateway", "Shared Data Store", "Another Tenant's Data"}, Target: "Another Tenant's Data"},
		},
		AssumptionPatterns: []KnowledgePackAssumption{
			{Description: "SaaS architecture must enforce multi-tenancy isolation at network, data, and compute layers", Category: "TrustBoundaries", Risk: "Critical", Component: "Tenant", Keywords: []string{"multi-tenancy", "tenant isolation"}},
			{Description: "SaaS architecture must enforce data segregation and prevent tenant data commingling", Category: "DataProtection", Risk: "Critical", Component: "Tenant", Keywords: []string{"data segregation", "commingling"}},
			{Description: "SaaS architecture must implement API security controls per tenant including rate limiting and object-level authorization", Category: "APISecurity", Risk: "High", Component: "API", Keywords: []string{"api", "rate limiting", "bola"}},
		},
	}
}

func governmentPack() KnowledgePack {
	return KnowledgePack{
		ID:                   "government",
		Name:                 "Government / Public Sector",
		Industry:             "Government & Public Sector",
		Version:              "1.0",
		Description:          "Government domain knowledge pack for architectures handling citizen data, classified systems, and federal networks.",
		DetectionKeywords:    []string{"government", "fedramp", "federal", "classified", "citizen", "public sector", "nist", "agency", "civic", "state", "municipal", "sovereign"},
		DetectionComponents:  []string{"government", "federal", "classified", "citizen", "agency"},
		DetectionCompliance:  []string{"FedRAMP", "NIST800-53", "NIST SP 800-53"},
		CrownJewels:          []string{"Citizen Data", "Classified System", "Federal Network", "Government Database"},
		ComplianceFrameworks: []string{"FedRAMP", "NIST SP 800-53", "FISMA", "EO 14028"},
		ExpectedControls: []KnowledgePackControl{
			{Name: "PrivilegedAccessReviews", Description: "Quarterly privileged access reviews with recertification", Category: "IdentityGovernance", Priority: "Critical"},
			{Name: "DataSovereignty", Description: "Data residency and sovereignty controls for citizen data", Category: "DataProtection", Priority: "Critical"},
			{Name: "NetworkSegmentation", Description: "Network segmentation with classified/unclassified boundaries", Category: "NetworkSecurity", Priority: "Critical"},
			{Name: "AuditRequirements", Description: "Comprehensive audit logging meeting FISMA requirements", Category: "AuditLogging", Priority: "Critical"},
			{Name: "PersonnelVetting", Description: "Personnel security clearance and background check verification", Category: "PersonnelSecurity", Priority: "Critical"},
			{Name: "PhysicalSecurity", Description: "Physical security controls for government data centers", Category: "PhysicalSecurity", Priority: "High"},
			{Name: "SupplyChainSecurity", Description: "Supply chain risk management for government systems", Category: "SupplyChain", Priority: "Critical"},
			{Name: "ContinuousMonitoring", Description: "Continuous monitoring and authorization (NIST SP 800-137)", Category: "Monitoring", Priority: "High"},
		},
		ExpectedEvidence: []KnowledgePackEvidence{
			{Control: "PrivilegedAccessReviews", Evidence: []string{"Access review schedule", "Completed review certifications", "Remediation tracking"}},
			{Control: "DataSovereignty", Evidence: []string{"Data residency policy", "Data classification schema", "Jurisdiction mapping"}},
			{Control: "NetworkSegmentation", Evidence: []string{"Network architecture diagram", "Cross-domain solution documentation", "Boundary protection procedures"}},
			{Control: "AuditRequirements", Evidence: []string{"Audit log configuration", "Log retention policy", "Audit trail integrity controls"}},
			{Control: "SupplyChainSecurity", Evidence: []string{"SBOM records", "Vendor assessments", "Supply chain risk management plan"}},
		},
		ThreatPatterns: []KnowledgePackThreat{
			{Name: "Citizen Data Breach", Description: "Attacker exfiltrates citizen PII from government databases", Severity: "Critical", Category: "Information Disclosure", Target: "Citizen Data"},
			{Name: "Privileged Access Abuse", Description: "Privileged user abuses elevated access to modify sensitive government records", Severity: "Critical", Category: "Tampering", Target: "Government Database"},
			{Name: "Cross-Domain Pivot", Description: "Attacker compromises unclassified system and pivots to classified network", Severity: "Critical", Category: "Elevation of Privilege", Target: "Classified System"},
			{Name: "Supply Chain Compromise", Description: "Attacker introduces backdoor via compromised government vendor/supplier", Severity: "Critical", Category: "Tampering", Target: "Federal Network"},
		},
		AttackPathTemplates: []KnowledgePackAttackPath{
			{Name: "Internet to Citizen Data Exfiltration", Description: "External attacker exploits web portal vulnerability to access citizen database", Steps: []string{"Internet", "Web Portal", "API", "Government Database", "Citizen Data"}, Target: "Citizen Data"},
			{Name: "Unclassified to Classified Pivot", Description: "Attacker compromises unclassified network and pivots across classified boundary", Steps: []string{"Unclassified Network", "Cross-Domain Solution", "Classified Network", "Classified Database"}, Target: "Classified System"},
		},
		AssumptionPatterns: []KnowledgePackAssumption{
			{Description: "Government system must implement privileged access reviews with recertification for all high-privilege roles", Category: "IdentityGovernance", Risk: "Critical", Component: "Government", Keywords: []string{"government", "privileged access", "recertification"}},
			{Description: "Government system must enforce data sovereignty and residency for citizen data", Category: "DataProtection", Risk: "Critical", Component: "Citizen", Keywords: []string{"data sovereignty", "citizen data", "residency"}},
			{Description: "Government system must implement network segmentation between classified and unclassified environments", Category: "NetworkSecurity", Risk: "Critical", Component: "Government", Keywords: []string{"government", "classified", "segmentation"}},
			{Description: "Government system must meet FISMA audit requirements with continuous monitoring", Category: "AuditLogging", Risk: "Critical", Component: "Government", Keywords: []string{"fisma", "audit", "continuous monitoring"}},
		},
	}
}

func criticalInfrastructurePack() KnowledgePack {
	return KnowledgePack{
		ID:                   "critical_infrastructure",
		Name:                 "Critical Infrastructure",
		Industry:             "Critical Infrastructure",
		Version:              "1.0",
		Description:          "Critical infrastructure knowledge pack for SCADA, ICS, industrial control systems, and utility networks.",
		DetectionKeywords:    []string{"scada", "ics", "industrial", "power", "utility", "plant", "plc", "rtu", "hmi", "modbus", "dnp3", "opc", "historian", "grid", "generator", "substation", "water", "energy"},
		DetectionComponents:  []string{"scada", "ics", "plc", "rtu", "hmi", "controller", "sensor", "actuator"},
		DetectionCompliance:  []string{"NERC", "NIST800-82", "IEC 62443"},
		CrownJewels:          []string{"SCADA Controller", "Historian Database", "Safety System", "Control Network"},
		ComplianceFrameworks: []string{"NIST SP 800-82", "IEC 62443", "NERC CIP"},
		ExpectedControls: []KnowledgePackControl{
			{Name: "OperationalContinuity", Description: "Operational continuity controls for industrial processes", Category: "Resilience", Priority: "Critical"},
			{Name: "PhysicalSecurity", Description: "Physical security for industrial control facilities", Category: "PhysicalSecurity", Priority: "Critical"},
			{Name: "RemoteAccessControls", Description: "Secure remote access for ICS maintenance and support", Category: "AccessControl", Priority: "Critical"},
			{Name: "SafetySystemIntegrity", Description: "Safety system integrity and independence from control network", Category: "Safety", Priority: "Critical"},
			{Name: "NetworkSegmentationICS", Description: "Purdue model network segmentation for ICS/SCADA", Category: "NetworkSecurity", Priority: "Critical"},
			{Name: "PatchManagementICS", Description: "ICS-specific patch management with change advisory board", Category: "ChangeManagement", Priority: "High"},
			{Name: "IncidentResponseICS", Description: "Incident response procedures for industrial control incidents", Category: "IncidentResponse", Priority: "High"},
			{Name: "AirGapMonitoring", Description: "Monitoring of air-gap boundaries and data diodes", Category: "Monitoring", Priority: "High"},
		},
		ExpectedEvidence: []KnowledgePackEvidence{
			{Control: "OperationalContinuity", Evidence: []string{"Failover configuration", "Disaster recovery plan", "Resilience test results"}},
			{Control: "PhysicalSecurity", Evidence: []string{"Facility access logs", "Physical security policy", "Security assessment reports"}},
			{Control: "RemoteAccessControls", Evidence: []string{"Remote access policy", "VPN configuration", "Session recording logs"}},
			{Control: "SafetySystemIntegrity", Evidence: []string{"Safety system architecture", "Independence verification", "Safety certification"}},
			{Control: "NetworkSegmentationICS", Evidence: []string{"Purdue model diagram", "Firewall rule sets", "Segment access logs"}},
			{Control: "PatchManagementICS", Evidence: []string{"ICS patch policy", "Change advisory board minutes", "Patch deployment records"}},
		},
		ThreatPatterns: []KnowledgePackThreat{
			{Name: "SCADA Controller Compromise", Description: "Attacker compromises SCADA controller to disrupt industrial processes", Severity: "Critical", Category: "Tampering", Target: "SCADA Controller"},
			{Name: "Safety System Bypass", Description: "Attacker bypasses safety system to cause physical damage", Severity: "Critical", Category: "Tampering", Target: "Safety System"},
			{Name: "ICS Ransomware", Description: "Ransomware encrypts historian and control systems halting production", Severity: "Critical", Category: "Denial of Service", Target: "Control Network"},
			{Name: "Remote Access Pivot", Description: "Attacker compromises remote access VPN and pivots to ICS network", Severity: "Critical", Category: "Elevation of Privilege", Target: "Remote Access"},
		},
		AttackPathTemplates: []KnowledgePackAttackPath{
			{Name: "Internet to SCADA Disruption", Description: "External attacker exploits remote access to manipulate SCADA controllers", Steps: []string{"Internet", "Remote Access VPN", "Corporate Network", "ICS DMZ", "SCADA Controller"}, Target: "SCADA Controller"},
			{Name: "Insider to Safety System Manipulation", Description: "Insider with ICS access manipulates safety system parameters", Steps: []string{"Insider Workstation", "HMI", "Safety System"}, Target: "Safety System"},
		},
		AssumptionPatterns: []KnowledgePackAssumption{
			{Description: "Critical infrastructure must implement Purdue model network segmentation for ICS/SCADA environments", Category: "NetworkSecurity", Risk: "Critical", Component: "SCADA", Keywords: []string{"scada", "ics", "purdue", "segmentation"}},
			{Description: "Critical infrastructure must enforce physical security controls for industrial control facilities", Category: "PhysicalSecurity", Risk: "Critical", Component: "ICS", Keywords: []string{"physical security", "ics", "facility"}},
			{Description: "Critical infrastructure must implement safety system independence from control network", Category: "Safety", Risk: "Critical", Component: "Safety", Keywords: []string{"safety system", "independence", "ics"}},
			{Description: "Critical infrastructure must enforce secure remote access for ICS maintenance with session recording", Category: "AccessControl", Risk: "Critical", Component: "ICS", Keywords: []string{"remote access", "ics", "maintenance"}},
		},
	}
}

// ── PHASE 9 — DOMAIN THREAT INJECTION ──

func injectDomainThreats(pack *KnowledgePack, existingThreats []Threat) []Threat {
	if pack == nil {
		return nil
	}
	injected := make([]Threat, 0)
	for i, t := range pack.ThreatPatterns {
		// Skip if a similar threat already exists
		duplicate := false
		for _, et := range existingThreats {
			if strings.Contains(strings.ToLower(et.Description), strings.ToLower(t.Name[:minInt(20, len(t.Name))])) {
				duplicate = true
				break
			}
		}
		if duplicate {
			continue
		}
		injected = append(injected, Threat{
			ID:          fmt.Sprintf("DKPI-T-%s-%03d", pack.ID, i+1),
			Name:        t.Name,
			Description: t.Description,
			Category:    ThreatCategory(t.Category),
			Severity:    RiskLevel(t.Severity),
			Likelihood:  4,
			Impact:      5,
			RiskScore:   20,
			Confidence:  0.75,
			Reasoning:   "Domain-specific threat from " + pack.Name + " knowledge pack",
		})
	}
	return injected
}

// ── PHASE 10 — DOMAIN ATTACK PATHS ──

func generateDomainAttackPaths(pack *KnowledgePack) []AttackPath {
	if pack == nil {
		return nil
	}
	paths := make([]AttackPath, 0)
	for i, ap := range pack.AttackPathTemplates {
		steps := make([]AttackStep, len(ap.Steps))
		for j, step := range ap.Steps {
			steps[j] = AttackStep{
				SequenceNumber:  j + 1,
				SourceComponent: step,
				Action:          "exploit",
				Reasoning:       "Domain-specific attack path step from " + pack.Name + " pack",
			}
			if j > 0 {
				steps[j].SourceComponent = ap.Steps[j-1]
				steps[j].TargetComponent = step
			}
		}
		paths = append(paths, AttackPath{
			ID:              fmt.Sprintf("DKPI-AP-%s-%03d", pack.ID, i+1),
			Name:            ap.Name,
			Description:     ap.Description,
			EntryPoint:      ap.Steps[0],
			TargetAsset:     ap.Target,
			AttackSteps:     steps,
			RiskScore:       0.7,
			Confidence:      0.75,
			BusinessImpact:  "Domain-specific attack from " + pack.Name + " knowledge pack",
			Recommendations: []string{"Review architecture against " + pack.Name + " attack path template"},
		})
	}
	return paths
}

// ── PHASE 11 — DOMAIN CONTROL EXPECTATIONS ──

func domainControlExpectations(pack *KnowledgePack) []SDRIControl {
	if pack == nil {
		return nil
	}
	controls := make([]SDRIControl, 0)
	for i, c := range pack.ExpectedControls {
		controls = append(controls, SDRIControl{
			ID:          fmt.Sprintf("DKPI-CTRL-%s-%03d", pack.ID, i+1),
			Name:        c.Name,
			Category:    c.Category,
			Description: c.Description,
			ControlType: SDRIControlPreventive,
			Coverage:    "Expected",
			Status:      "Domain Expected",
			Strength:    SDRIStrengthStrong,
		})
	}
	return controls
}

// ── PHASE 12 — DOMAIN COMPLIANCE PACKS ──

func domainCompliancePack(pack *KnowledgePack) []string {
	if pack == nil {
		return nil
	}
	return pack.ComplianceFrameworks
}

// ── PHASE 13 — DOMAIN EVIDENCE REQUIREMENTS ──

func domainEvidenceRequirements(pack *KnowledgePack) []KnowledgePackEvidence {
	if pack == nil {
		return nil
	}
	return pack.ExpectedEvidence
}

// ── PHASE 14 — DOMAIN-SPECIFIC RECOMMENDATIONS ──

func domainRecommendations(pack *KnowledgePack) []string {
	if pack == nil {
		return nil
	}
	recs := make([]string, 0)
	for _, c := range pack.ExpectedControls {
		recs = append(recs, fmt.Sprintf("[%s] %s: %s", c.Priority, c.Name, c.Description))
	}
	if len(pack.CrownJewels) > 0 {
		recs = append(recs, fmt.Sprintf("Identify and protect crown jewels: %s", strings.Join(pack.CrownJewels, ", ")))
	}
	return recs
}

// ── PHASE 15 — DOMAIN CONFIDENCE BOOSTING ──

func boostDomainConfidence(assumptions []Assumption, domain string, confidence float64) []Assumption {
	if domain == "" || confidence <= 0 {
		return assumptions
	}
	boosted := make([]Assumption, len(assumptions))
	copy(boosted, assumptions)
	for i := range boosted {
		if boosted[i].Confidence < 0.95 {
			boosted[i].Confidence += 0.25 * (confidence / 100.0)
			if boosted[i].Confidence > 0.95 {
				boosted[i].Confidence = 0.95
			}
		}
		boosted[i].Keywords = append(boosted[i].Keywords, "domain:"+domain)
	}
	return boosted
}

// ── PHASE 9-12 — ENRICH CIE / TBI / TMI / APD / SDRI / CIARE RESULTS ──

func enrichControlStrength(controls []SDRIControl, pack *KnowledgePack) []SDRIControl {
	if pack == nil {
		return controls
	}
	enriched := make([]SDRIControl, len(controls))
	copy(enriched, controls)
	for i, c := range enriched {
		for _, ec := range pack.ExpectedControls {
			if normalizeControlName(c.Name) == normalizeControlName(ec.Name) {
				enriched[i].Category = ec.Category
				if enriched[i].Coverage == "Missing" || enriched[i].Coverage == "" {
					enriched[i].Coverage = "Expected by " + pack.Name + " pack"
				} else if enriched[i].Coverage == "Partial" {
					enriched[i].Coverage = "Enhanced"
				}
			}
		}
	}
	return enriched
}

// ── MAIN DKPI ENGINE ──

type DKPIEngine struct{}

func NewDKPIEngine() *DKPIEngine {
	return &DKPIEngine{}
}

type DKPIInput struct {
	Architecture        *ArchDescription
	ExistingAssumptions []Assumption
	ExistingThreats     []Threat
	ExistingControls    []SDRIControl
	ExistingFindings    []SDRIFinding
	Domain              string
	Compliance          []string
}

type DKPIFindings struct {
	ID                 string   `json:"id"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	Severity           string   `json:"severity"`
	Category           string   `json:"category"`
	AffectedComponents []string `json:"affected_components,omitempty"`
	BusinessImpact     string   `json:"business_impact"`
	Recommendation     string   `json:"recommendation"`
}

func (e *DKPIEngine) Run(input DKPIInput) *DKPIEngineResult {
	detector := NewDKPIDetector()

	// Phase 2: Domain detection
	detected := detector.DetectDomain(input.Architecture)

	// Fall back to declared domain if detection yielded nothing
	if detected.PrimaryDomain == "" && input.Domain != "" {
		id := strings.ToLower(strings.ReplaceAll(input.Domain, " ", "_"))
		detected.PrimaryDomain = id
		detected.Confidence = 50.0
		detected.Rationale = []string{"Using declared domain: " + input.Domain}
	}

	// Get active knowledge pack
	var activePack *KnowledgePack
	if detected.PrimaryDomain != "" {
		p := detector.GetPack(detected.PrimaryDomain)
		if p == nil {
			// Try matching by name
			for _, pack := range buildKnowledgePacks() {
				if strings.EqualFold(pack.Name, detected.PrimaryDomain) || strings.EqualFold(pack.ID, detected.PrimaryDomain) {
					p = &pack
					break
				}
			}
		}
		activePack = p
	}

	result := &DKPIEngineResult{
		DetectedDomain: detected,
		ActivePack:     activePack,
	}

	// Phase 9: Inject domain threats
	if activePack != nil {
		result.InjectedThreats = injectDomainThreats(activePack, input.ExistingThreats)
	}

	// Phase 10: Generate domain attack paths
	if activePack != nil {
		result.GeneratedAttackPaths = generateDomainAttackPaths(activePack)
	}

	// Phase 11: Domain control expectations
	if activePack != nil {
		result.DomainControls = domainControlExpectations(activePack)
	}

	// Phase 12: Domain compliance packs
	if activePack != nil {
		result.DomainCompliance = domainCompliancePack(activePack)
	}

	// Phase 13: Domain evidence requirements
	if activePack != nil {
		result.DomainEvidence = domainEvidenceRequirements(activePack)
	}

	// Phase 14: Domain-specific recommendations
	if activePack != nil {
		result.Recommendations = domainRecommendations(activePack)
	}

	// Phase 15: Domain confidence boosting (returns boosted assumptions)
	if activePack != nil && len(input.ExistingAssumptions) > 0 {
		result.BoostedAssumptions = boostDomainConfidence(input.ExistingAssumptions, detected.PrimaryDomain, detected.Confidence)
	}

	// Phase 11: Enrich existing controls with domain context
	if activePack != nil && len(input.ExistingControls) > 0 {
		result.EnrichedControls = enrichControlStrength(input.ExistingControls, activePack)
	}

	return result
}

type DKPIEngineResult struct {
	DetectedDomain       DomainDetectionResult   `json:"detected_domain"`
	ActivePack           *KnowledgePack          `json:"active_pack"`
	InjectedThreats      []Threat                `json:"injected_threats"`
	GeneratedAttackPaths []AttackPath            `json:"generated_attack_paths"`
	DomainControls       []SDRIControl           `json:"domain_controls"`
	DomainCompliance     []string                `json:"domain_compliance"`
	DomainEvidence       []KnowledgePackEvidence `json:"domain_evidence"`
	Recommendations      []string                `json:"recommendations"`
	BoostedAssumptions   []Assumption            `json:"boosted_assumptions"`
	EnrichedControls     []SDRIControl           `json:"enriched_controls"`
}
