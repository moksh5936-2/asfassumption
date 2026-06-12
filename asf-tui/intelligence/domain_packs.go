package intelligence

import (
	"sort"
	"strings"
)

// DomainPack holds domain-specific assumptions, controls, risk amplifiers, and compliance mappings.
type DomainPack struct {
	Name               string
	Assumptions        []Assumption
	Controls           []ControlDetail
	RiskAmplifiers     map[string]RiskLevel
	ComplianceMappings []string
}

// DomainEngine manages domain packs and detection.
type DomainEngine struct {
	packs map[string]*DomainPack
}

// NewDomainEngine creates a domain engine with all built-in packs.
func NewDomainEngine() *DomainEngine {
	de := &DomainEngine{
		packs: make(map[string]*DomainPack),
	}
	de.packs["Healthcare"] = HealthcarePack()
	de.packs["Fintech"] = FintechPack()
	de.packs["SaaS"] = SaaSPack()
	de.packs["Enterprise"] = EnterprisePack()
	de.packs["Kubernetes"] = KubernetesPack()
	de.packs["CloudNative"] = CloudNativePack()
	de.packs["VPN"] = VPNPack()
	de.packs["IdentityPlatform"] = IdentityPlatformPack()
	de.packs["DataPlatform"] = DataPlatformPack()
	return de
}

// GetPack returns a domain pack by name.
func (de *DomainEngine) GetPack(name string) *DomainPack {
	return de.packs[name]
}

// DetectDomain auto-detects the domain from architecture components and raw text.
func (de *DomainEngine) DetectDomain(arch *ArchDescription) string {
	if arch == nil {
		return ""
	}
	scores := make(map[string]int)
	raw := strings.ToLower(arch.RawText)

	for _, comp := range arch.Components {
		label := strings.ToLower(comp.Label)
		if strings.Contains(label, "phi") || strings.Contains(label, "patient") || strings.Contains(label, "health") || strings.Contains(label, "ehr") || strings.Contains(label, "clinic") || strings.Contains(label, "hospital") || strings.Contains(label, "medical") {
			scores["Healthcare"] += 3
		}
		if strings.Contains(label, "payment") || strings.Contains(label, "card") || strings.Contains(label, "bank") || strings.Contains(label, "transaction") || strings.Contains(label, "fraud") || strings.Contains(label, "pci") {
			scores["Fintech"] += 3
		}
		if strings.Contains(label, "tenant") || strings.Contains(label, "saas") || strings.Contains(label, "multi-tenant") || strings.Contains(label, "subscription") {
			scores["SaaS"] += 3
		}
		if strings.Contains(label, "ad") || strings.Contains(label, "ldap") || strings.Contains(label, "domain") || strings.Contains(label, "enterprise") || strings.Contains(label, "corporate") {
			scores["Enterprise"] += 2
		}
		if strings.Contains(label, "kubernetes") || strings.Contains(label, "k8s") || strings.Contains(label, "pod") || strings.Contains(label, "cluster") || strings.Contains(label, "helm") {
			scores["Kubernetes"] += 3
		}
		if strings.Contains(label, "cloud") || strings.Contains(label, "aws") || strings.Contains(label, "azure") || strings.Contains(label, "gcp") || strings.Contains(label, "lambda") || strings.Contains(label, "function") {
			scores["CloudNative"] += 2
		}
		if strings.Contains(label, "vpn") || strings.Contains(label, "tunnel") || strings.Contains(label, "remote access") {
			scores["VPN"] += 3
		}
		if strings.Contains(label, "sso") || strings.Contains(label, "federation") || strings.Contains(label, "identity platform") || strings.Contains(label, "auth0") || strings.Contains(label, "okta") {
			scores["IdentityPlatform"] += 3
		}
		if strings.Contains(label, "data lake") || strings.Contains(label, "warehouse") || strings.Contains(label, "etl") || strings.Contains(label, "pipeline") || strings.Contains(label, "analytics") {
			scores["DataPlatform"] += 3
		}
	}

	// Text-based scoring
	healthKeywords := []string{"phi", "hipaa", "patient", "health", "medical", "ehr", "clinical", "hl7", "fda"}
	fintechKeywords := []string{"pci", "payment", "card", "sox", "fraud", "transaction", "bank", "financial", "trading", "settlement"}
	saasKeywords := []string{"multi-tenant", "tenant", "saas", "subscription", "onboarding", "tenant isolation"}
	enterpriseKeywords := []string{"enterprise", "identity lifecycle", "access review", "governance", "compliance", "ad"}
	k8sKeywords := []string{"kubernetes", "k8s", "container", "pod", "namespace", "rbac", "network policy", "cni"}
	cloudKeywords := []string{"cloud native", "serverless", "faas", "lambda", "cloud function", "iam", "cloud security"}
	vpnKeywords := []string{"vpn", "tunnel", "ipsec", "ssl vpn", "remote access", "endpoint validation"}
	identityKeywords := []string{"sso", "federation", "saml", "identity platform", "session management", "token rotation"}
	dataKeywords := []string{"data governance", "data lineage", "data retention", "etl", "data warehouse", "catalog"}

	for _, kw := range healthKeywords {
		if strings.Contains(raw, kw) {
			scores["Healthcare"] += 2
		}
	}
	for _, kw := range fintechKeywords {
		if strings.Contains(raw, kw) {
			scores["Fintech"] += 2
		}
	}
	for _, kw := range saasKeywords {
		if strings.Contains(raw, kw) {
			scores["SaaS"] += 2
		}
	}
	for _, kw := range enterpriseKeywords {
		if strings.Contains(raw, kw) {
			scores["Enterprise"] += 2
		}
	}
	for _, kw := range k8sKeywords {
		if strings.Contains(raw, kw) {
			scores["Kubernetes"] += 2
		}
	}
	for _, kw := range cloudKeywords {
		if strings.Contains(raw, kw) {
			scores["CloudNative"] += 2
		}
	}
	for _, kw := range vpnKeywords {
		if strings.Contains(raw, kw) {
			scores["VPN"] += 2
		}
	}
	for _, kw := range identityKeywords {
		if strings.Contains(raw, kw) {
			scores["IdentityPlatform"] += 2
		}
	}
	for _, kw := range dataKeywords {
		if strings.Contains(raw, kw) {
			scores["DataPlatform"] += 2
		}
	}

	bestDomain := ""
	bestScore := 0
	// Sort domains alphabetically for deterministic tie-breaking
	var domains []string
	for domain := range scores {
		domains = append(domains, domain)
	}
	sort.Strings(domains)
	for _, domain := range domains {
		score := scores[domain]
		if score > bestScore {
			bestScore = score
			bestDomain = domain
		}
	}
	if bestScore < 3 {
		return ""
	}
	return bestDomain
}

// ApplyDomainPack adds domain-specific assumptions to the result set.
func (de *DomainEngine) ApplyDomainPack(domain string, arch *ArchDescription) []Assumption {
	pack := de.GetPack(domain)
	if pack == nil || arch == nil {
		return nil
	}
	var results []Assumption
	for _, a := range pack.Assumptions {
		// Copy and enrich
		newA := a
		newA.SourceType = "domain-inferred"
		newA.SourceFile = arch.Name
		results = append(results, newA)
	}
	return results
}

// HealthcarePack returns domain-specific assumptions for healthcare.
func HealthcarePack() *DomainPack {
	return &DomainPack{
		Name: "Healthcare",
		Assumptions: []Assumption{
			{
				ID:          "DOM-HLT-001",
				Description: "Healthcare system must ensure PHI access is audited and meets HIPAA audit controls (164.312(b)).",
				Component:   "PHI",
				Category:    "Auditability",
				Risk:        RiskCritical,
				Likelihood:  5,
				Impact:      5,
				Confidence:  0.95,
				Keywords:    []string{"phi", "hipaa", "audit", "healthcare"},
				Rationale:   "HIPAA requires audit controls for PHI access; failure to audit creates regulatory violation and patient privacy risk.",
			},
			{
				ID:          "DOM-HLT-002",
				Description: "Healthcare system must enforce patient privacy and minimum necessary access for all PHI interactions.",
				Component:   "PHI",
				Category:    "Privacy",
				Risk:        RiskCritical,
				Likelihood:  4,
				Impact:      5,
				Confidence:  0.95,
				Keywords:    []string{"phi", "patient privacy", "minimum necessary", "healthcare"},
				Rationale:   "HIPAA Privacy Rule requires minimum necessary access; overexposure violates patient trust and regulation.",
			},
			{
				ID:          "DOM-HLT-003",
				Description: "Healthcare system must specify data retention periods aligned with medical record retention laws and patient deletion rights.",
				Component:   "PHI",
				Category:    "DataRetention",
				Risk:        RiskHigh,
				Likelihood:  3,
				Impact:      4,
				Confidence:  0.90,
				Keywords:    []string{"phi", "retention", "deletion", "medical records"},
				Rationale:   "Medical records have state-specific retention requirements; patient deletion requests must also be honored.",
			},
			{
				ID:          "DOM-HLT-004",
				Description: "Healthcare system must implement break-glass access for emergency PHI retrieval with immutable audit logging.",
				Component:   "PHI",
				Category:    "PrivilegeManagement",
				Risk:        RiskCritical,
				Likelihood:  3,
				Impact:      5,
				Confidence:  0.90,
				Keywords:    []string{"phi", "break-glass", "emergency", "audit"},
				Rationale:   "Emergency access to PHI without break-glass controls can delay care; without audit logging it violates HIPAA.",
			},
			{
				ID:          "DOM-HLT-005",
				Description: "Healthcare system must encrypt PHI at rest and in transit with approved algorithms and key management.",
				Component:   "PHI",
				Category:    "DataProtection",
				Risk:        RiskCritical,
				Likelihood:  4,
				Impact:      5,
				Confidence:  0.95,
				Keywords:    []string{"phi", "encryption", "key management", "healthcare"},
				Rationale:   "HIPAA Security Rule requires encryption of PHI at rest and in transit; weak encryption is a reportable breach.",
			},
			{
				ID:          "DOM-HLT-006",
				Description: "Healthcare system must support export controls for data leaving the jurisdiction and maintain data residency compliance.",
				Component:   "PHI",
				Category:    "Compliance",
				Risk:        RiskHigh,
				Likelihood:  3,
				Impact:      4,
				Confidence:  0.85,
				Keywords:    []string{"phi", "export control", "data residency", "jurisdiction"},
				Rationale:   "Cross-border PHI transfer may violate data residency laws; export controls and jurisdiction mapping are required.",
			},
		},
		Controls: []ControlDetail{
			{ID: "CTRL-HLT-001", Description: "Implement HIPAA audit controls for all PHI access", Category: "Auditability", Priority: "High"},
			{ID: "CTRL-HLT-002", Description: "Implement minimum necessary access controls for PHI", Category: "Privacy", Priority: "High"},
			{ID: "CTRL-HLT-003", Description: "Implement break-glass procedures with immutable logging", Category: "PrivilegeManagement", Priority: "High"},
		},
		RiskAmplifiers: map[string]RiskLevel{
			"PHI":     RiskCritical,
			"HIPAA":   RiskCritical,
			"Patient": RiskHigh,
			"Medical": RiskHigh,
		},
		ComplianceMappings: []string{"HIPAA", "HITECH", "FDA 21 CFR Part 11", "State Medical Privacy Laws"},
	}
}

// FintechPack returns domain-specific assumptions for fintech.
func FintechPack() *DomainPack {
	return &DomainPack{
		Name: "Fintech",
		Assumptions: []Assumption{
			{
				ID:          "DOM-FIN-001",
				Description: "Fintech system must enforce PCI DSS requirements for cardholder data environment, including segmentation and access controls.",
				Component:   "Payment",
				Category:    "Compliance",
				Risk:        RiskCritical,
				Likelihood:  5,
				Impact:      5,
				Confidence:  0.95,
				Keywords:    []string{"pci", "payment", "cardholder", "fintech"},
				Rationale:   "PCI DSS non-compliance results in fines, brand damage, and potential loss of card processing privileges.",
			},
			{
				ID:          "DOM-FIN-002",
				Description: "Fintech system must implement SOX controls for financial reporting integrity, change management, and access logging.",
				Component:   "Financial",
				Category:    "Compliance",
				Risk:        RiskCritical,
				Likelihood:  4,
				Impact:      5,
				Confidence:  0.95,
				Keywords:    []string{"sox", "financial reporting", "change management", "fintech"},
				Rationale:   "SOX requires financial reporting controls; absence creates audit failure and regulatory penalties.",
			},
			{
				ID:          "DOM-FIN-003",
				Description: "Fintech system must implement fraud detection and transaction integrity monitoring for all payment flows.",
				Component:   "Payment",
				Category:    "DetectionEngineering",
				Risk:        RiskCritical,
				Likelihood:  4,
				Impact:      5,
				Confidence:  0.95,
				Keywords:    []string{"fraud", "transaction integrity", "payment", "fintech"},
				Rationale:   "Fraud detection gaps lead to financial loss, customer harm, and regulatory enforcement.",
			},
			{
				ID:          "DOM-FIN-004",
				Description: "Fintech system must implement dedicated key management for payment encryption and HSM-backed key storage.",
				Component:   "Payment",
				Category:    "KeyManagement",
				Risk:        RiskCritical,
				Likelihood:  4,
				Impact:      5,
				Confidence:  0.95,
				Keywords:    []string{"key management", "hsm", "payment encryption", "fintech"},
				Rationale:   "Payment keys require HSM-backed storage and strict rotation; compromise leads to mass card reissue.",
			},
		},
		Controls: []ControlDetail{
			{ID: "CTRL-FIN-001", Description: "Implement PCI DSS segmentation and access controls", Category: "Compliance", Priority: "High"},
			{ID: "CTRL-FIN-002", Description: "Implement SOX change management and financial reporting logging", Category: "Compliance", Priority: "High"},
			{ID: "CTRL-FIN-003", Description: "Implement real-time fraud detection and transaction monitoring", Category: "DetectionEngineering", Priority: "High"},
		},
		RiskAmplifiers: map[string]RiskLevel{
			"Payment":     RiskCritical,
			"PCI":         RiskCritical,
			"SOX":         RiskCritical,
			"Financial":   RiskHigh,
			"Transaction": RiskHigh,
		},
		ComplianceMappings: []string{"PCI DSS", "SOX", "GDPR", "PSD2", "FFIEC"},
	}
}

// SaaSPack returns domain-specific assumptions for SaaS.
func SaaSPack() *DomainPack {
	return &DomainPack{
		Name: "SaaS",
		Assumptions: []Assumption{
			{
				ID:          "DOM-SAS-001",
				Description: "SaaS architecture must enforce multi-tenancy isolation at network, data, and compute layers.",
				Component:   "Tenant",
				Category:    "TrustBoundaries",
				Risk:        RiskCritical,
				Likelihood:  4,
				Impact:      5,
				Confidence:  0.95,
				Keywords:    []string{"multi-tenancy", "tenant isolation", "saas"},
				Rationale:   "Tenant isolation failures lead to cross-tenant data leakage, reputation loss, and regulatory action.",
			},
			{
				ID:          "DOM-SAS-002",
				Description: "SaaS architecture must enforce data segregation and prevent tenant data commingling in shared storage.",
				Component:   "Tenant",
				Category:    "DataProtection",
				Risk:        RiskCritical,
				Likelihood:  4,
				Impact:      5,
				Confidence:  0.95,
				Keywords:    []string{"tenant", "data segregation", "commingling", "saas"},
				Rationale:   "Shared storage without segregation allows one tenant to access another tenant's data.",
			},
			{
				ID:          "DOM-SAS-003",
				Description: "SaaS architecture must implement API security controls per tenant, including rate limiting and object-level authorization.",
				Component:   "API",
				Category:    "APISecurity",
				Risk:        RiskHigh,
				Likelihood:  4,
				Impact:      4,
				Confidence:  0.90,
				Keywords:    []string{"api", "tenant", "rate limiting", "bola", "saas"},
				Rationale:   "Tenant-scoped API security prevents one tenant from exhausting resources or accessing another tenant's objects.",
			},
		},
		Controls: []ControlDetail{
			{ID: "CTRL-SAS-001", Description: "Implement tenant isolation at network and compute layers", Category: "TrustBoundaries", Priority: "High"},
			{ID: "CTRL-SAS-002", Description: "Implement data segregation and row-level security per tenant", Category: "DataProtection", Priority: "High"},
			{ID: "CTRL-SAS-003", Description: "Implement tenant-scoped API rate limiting and object-level authorization", Category: "APISecurity", Priority: "High"},
		},
		RiskAmplifiers: map[string]RiskLevel{
			"Tenant":      RiskCritical,
			"MultiTenant": RiskCritical,
			"SaaS":        RiskHigh,
		},
		ComplianceMappings: []string{"SOC 2", "ISO 27001", "GDPR", "CCPA"},
	}
}

// EnterprisePack returns domain-specific assumptions for enterprise.
func EnterprisePack() *DomainPack {
	return &DomainPack{
		Name: "Enterprise",
		Assumptions: []Assumption{
			{
				ID:          "DOM-ENT-001",
				Description: "Enterprise architecture must implement identity lifecycle management with automated provisioning and deprovisioning.",
				Component:   "Identity",
				Category:    "Identity",
				Risk:        RiskHigh,
				Likelihood:  4,
				Impact:      4,
				Confidence:  0.90,
				Keywords:    []string{"identity lifecycle", "provisioning", "deprovisioning", "enterprise"},
				Rationale:   "Stale accounts and orphan access create persistent attack surface; lifecycle automation is required.",
			},
			{
				ID:          "DOM-ENT-002",
				Description: "Enterprise architecture must implement periodic access reviews and recertification for all privileged roles.",
				Component:   "Access",
				Category:    "Authorization",
				Risk:        RiskHigh,
				Likelihood:  4,
				Impact:      4,
				Confidence:  0.90,
				Keywords:    []string{"access review", "recertification", "privileged", "enterprise"},
				Rationale:   "Access reviews prevent privilege creep and ensure continued need-to-know for sensitive roles.",
			},
			{
				ID:          "DOM-ENT-003",
				Description: "Enterprise architecture must implement governance controls with security committee oversight and policy enforcement.",
				Component:   "Governance",
				Category:    "Governance",
				Risk:        RiskMedium,
				Likelihood:  3,
				Impact:      3,
				Confidence:  0.85,
				Keywords:    []string{"governance", "security committee", "policy enforcement", "enterprise"},
				Rationale:   "Enterprise governance ensures security policy consistency across business units and geographies.",
			},
		},
		Controls: []ControlDetail{
			{ID: "CTRL-ENT-001", Description: "Implement identity lifecycle automation", Category: "Identity", Priority: "High"},
			{ID: "CTRL-ENT-002", Description: "Implement quarterly access reviews and recertification", Category: "Authorization", Priority: "High"},
			{ID: "CTRL-ENT-003", Description: "Implement security governance committee and policy enforcement", Category: "Governance", Priority: "Medium"},
		},
		RiskAmplifiers: map[string]RiskLevel{
			"Enterprise": RiskHigh,
			"AD":         RiskHigh,
			"Governance": RiskMedium,
		},
		ComplianceMappings: []string{"SOX", "ISO 27001", "NIST CSF", "COBIT"},
	}
}

// KubernetesPack returns domain-specific assumptions for Kubernetes.
func KubernetesPack() *DomainPack {
	return &DomainPack{
		Name: "Kubernetes",
		Assumptions: []Assumption{
			{
				ID:          "DOM-K8S-001",
				Description: "Kubernetes architecture must enforce pod security standards and admission controller policies.",
				Component:   "Kubernetes",
				Category:    "ContainerSecurity",
				Risk:        RiskCritical,
				Likelihood:  4,
				Impact:      5,
				Confidence:  0.95,
				Keywords:    []string{"kubernetes", "pod security", "admission controller", "psp"},
				Rationale:   "Unrestricted pod configurations allow privilege escalation, host namespace access, and container escape.",
			},
			{
				ID:          "DOM-K8S-002",
				Description: "Kubernetes architecture must enforce RBAC with least-privilege service accounts and namespace isolation.",
				Component:   "Kubernetes",
				Category:    "KubernetesSecurity",
				Risk:        RiskCritical,
				Likelihood:  4,
				Impact:      5,
				Confidence:  0.95,
				Keywords:    []string{"kubernetes", "rbac", "service account", "namespace isolation"},
				Rationale:   "Overly permissive RBAC and default service accounts create cluster-wide lateral movement paths.",
			},
			{
				ID:          "DOM-K8S-003",
				Description: "Kubernetes architecture must implement network policies to restrict pod-to-pod traffic and segment workloads.",
				Component:   "Kubernetes",
				Category:    "NetworkSegmentation",
				Risk:        RiskHigh,
				Likelihood:  4,
				Impact:      4,
				Confidence:  0.90,
				Keywords:    []string{"kubernetes", "network policy", "cni", "segmentation"},
				Rationale:   "Default allow-all pod networking enables lateral movement; network policies are required for zero trust.",
			},
			{
				ID:          "DOM-K8S-004",
				Description: "Kubernetes architecture must implement secrets management with external vault integration and no plaintext secrets in manifests.",
				Component:   "Kubernetes",
				Category:    "SecretsManagement",
				Risk:        RiskCritical,
				Likelihood:  4,
				Impact:      5,
				Confidence:  0.95,
				Keywords:    []string{"kubernetes", "secrets", "vault", "plaintext", "manifest"},
				Rationale:   "Plaintext secrets in manifests or etcd expose credentials; external vault integration is required.",
			},
		},
		Controls: []ControlDetail{
			{ID: "CTRL-K8S-001", Description: "Enforce pod security standards and admission controllers", Category: "ContainerSecurity", Priority: "High"},
			{ID: "CTRL-K8S-002", Description: "Implement least-privilege RBAC and namespace isolation", Category: "KubernetesSecurity", Priority: "High"},
			{ID: "CTRL-K8S-003", Description: "Implement network policies for all namespaces", Category: "NetworkSegmentation", Priority: "High"},
			{ID: "CTRL-K8S-004", Description: "Integrate external vault for secrets management", Category: "SecretsManagement", Priority: "High"},
		},
		RiskAmplifiers: map[string]RiskLevel{
			"Kubernetes": RiskCritical,
			"Container":  RiskHigh,
			"Pod":        RiskHigh,
		},
		ComplianceMappings: []string{"NIST SP 800-190", "CIS Kubernetes Benchmark", "PCI DSS"},
	}
}

// CloudNativePack returns domain-specific assumptions for cloud-native.
func CloudNativePack() *DomainPack {
	return &DomainPack{
		Name: "CloudNative",
		Assumptions: []Assumption{
			{
				ID:          "DOM-CLD-001",
				Description: "Cloud-native architecture must implement IAM with least-privilege policies, MFA, and periodic access reviews.",
				Component:   "Cloud",
				Category:    "Identity",
				Risk:        RiskCritical,
				Likelihood:  4,
				Impact:      5,
				Confidence:  0.95,
				Keywords:    []string{"cloud", "iam", "least privilege", "mfa", "access review"},
				Rationale:   "Cloud IAM misconfiguration is a leading cause of breaches; least privilege and MFA are essential.",
			},
			{
				ID:          "DOM-CLD-002",
				Description: "Cloud-native architecture must enforce encryption at rest and in transit for all data stores and communication channels.",
				Component:   "Cloud",
				Category:    "DataProtection",
				Risk:        RiskCritical,
				Likelihood:  4,
				Impact:      5,
				Confidence:  0.95,
				Keywords:    []string{"cloud", "encryption", "at rest", "in transit", "data protection"},
				Rationale:   "Cloud data stores without encryption are exposed to insider access and provider compromise.",
			},
			{
				ID:          "DOM-CLD-003",
				Description: "Cloud-native architecture must implement centralized logging, monitoring, and alerting for all workloads and control plane events.",
				Component:   "Cloud",
				Category:    "Monitoring",
				Risk:        RiskHigh,
				Likelihood:  4,
				Impact:      4,
				Confidence:  0.90,
				Keywords:    []string{"cloud", "logging", "monitoring", "alerting", "control plane"},
				Rationale:   "Cloud control plane and workload events must be centralized for detection and compliance.",
			},
			{
				ID:          "DOM-CLD-004",
				Description: "Cloud-native architecture must implement resilience patterns including auto-scaling, circuit breakers, and multi-region failover.",
				Component:   "Cloud",
				Category:    "Resilience",
				Risk:        RiskHigh,
				Likelihood:  3,
				Impact:      4,
				Confidence:  0.90,
				Keywords:    []string{"cloud", "resilience", "auto-scaling", "circuit breaker", "multi-region"},
				Rationale:   "Cloud-native workloads require resilience patterns to handle provider failures and traffic spikes.",
			},
		},
		Controls: []ControlDetail{
			{ID: "CTRL-CLD-001", Description: "Implement least-privilege IAM with MFA", Category: "Identity", Priority: "High"},
			{ID: "CTRL-CLD-002", Description: "Implement encryption at rest and in transit", Category: "DataProtection", Priority: "High"},
			{ID: "CTRL-CLD-003", Description: "Implement centralized logging and monitoring", Category: "Monitoring", Priority: "High"},
			{ID: "CTRL-CLD-004", Description: "Implement resilience patterns and multi-region failover", Category: "Resilience", Priority: "High"},
		},
		RiskAmplifiers: map[string]RiskLevel{
			"Cloud":      RiskHigh,
			"AWS":        RiskHigh,
			"Azure":      RiskHigh,
			"GCP":        RiskHigh,
			"Serverless": RiskHigh,
		},
		ComplianceMappings: []string{"SOC 2", "ISO 27017", "FedRAMP", "CSA CCM"},
	}
}

// VPNPack returns domain-specific assumptions for VPN.
func VPNPack() *DomainPack {
	return &DomainPack{
		Name: "VPN",
		Assumptions: []Assumption{
			{
				ID:          "DOM-VPN-001",
				Description: "VPN architecture must enforce tunnel encryption with strong cipher suites and perfect forward secrecy.",
				Component:   "VPN",
				Category:    "NetworkSegmentation",
				Risk:        RiskCritical,
				Likelihood:  4,
				Impact:      5,
				Confidence:  0.95,
				Keywords:    []string{"vpn", "tunnel encryption", "pfs", "cipher"},
				Rationale:   "Weak VPN tunnel encryption allows interception and decryption; strong ciphers and PFS are mandatory.",
			},
			{
				ID:          "DOM-VPN-002",
				Description: "VPN architecture must implement endpoint validation including device health, certificate checks, and posture assessment.",
				Component:   "VPN",
				Category:    "InfrastructureSecurity",
				Risk:        RiskHigh,
				Likelihood:  4,
				Impact:      4,
				Confidence:  0.90,
				Keywords:    []string{"vpn", "endpoint validation", "device health", "certificate", "posture"},
				Rationale:   "VPN without endpoint validation allows compromised or unmanaged devices onto the internal network.",
			},
			{
				ID:          "DOM-VPN-003",
				Description: "VPN architecture must implement certificate management with automated issuance, rotation, and revocation.",
				Component:   "VPN",
				Category:    "CertificateManagement",
				Risk:        RiskHigh,
				Likelihood:  3,
				Impact:      4,
				Confidence:  0.90,
				Keywords:    []string{"vpn", "certificate", "rotation", "revocation", "pki"},
				Rationale:   "VPN certificate expiry or compromise without automated rotation and revocation breaks remote access security.",
			},
		},
		Controls: []ControlDetail{
			{ID: "CTRL-VPN-001", Description: "Enforce strong tunnel encryption with PFS", Category: "NetworkSegmentation", Priority: "High"},
			{ID: "CTRL-VPN-002", Description: "Implement endpoint validation and posture assessment", Category: "InfrastructureSecurity", Priority: "High"},
			{ID: "CTRL-VPN-003", Description: "Implement automated certificate lifecycle management", Category: "CertificateManagement", Priority: "High"},
		},
		RiskAmplifiers: map[string]RiskLevel{
			"VPN":    RiskHigh,
			"Tunnel": RiskHigh,
			"Remote": RiskHigh,
		},
		ComplianceMappings: []string{"NIST SP 800-114", "ISO 27001", "CIS Controls"},
	}
}

// IdentityPlatformPack returns domain-specific assumptions for identity platforms.
func IdentityPlatformPack() *DomainPack {
	return &DomainPack{
		Name: "IdentityPlatform",
		Assumptions: []Assumption{
			{
				ID:          "DOM-IDP-001",
				Description: "Identity platform must implement SSO with secure token exchange and session binding.",
				Component:   "IdentityPlatform",
				Category:    "Authentication",
				Risk:        RiskCritical,
				Likelihood:  4,
				Impact:      5,
				Confidence:  0.95,
				Keywords:    []string{"sso", "token exchange", "session binding", "identity platform"},
				Rationale:   "SSO without secure token binding is vulnerable to token theft, replay, and session hijacking.",
			},
			{
				ID:          "DOM-IDP-002",
				Description: "Identity platform must enforce MFA with phishing-resistant methods (FIDO2, WebAuthn) for all users.",
				Component:   "IdentityPlatform",
				Category:    "Authentication",
				Risk:        RiskCritical,
				Likelihood:  4,
				Impact:      5,
				Confidence:  0.95,
				Keywords:    []string{"mfa", "fido2", "webauthn", "phishing resistant", "identity platform"},
				Rationale:   "Legacy MFA (SMS, TOTP) is vulnerable to phishing; FIDO2/WebAuthn is required for high-assurance identity.",
			},
			{
				ID:          "DOM-IDP-003",
				Description: "Identity platform must implement session management with rotation, concurrent limits, and global sign-out.",
				Component:   "IdentityPlatform",
				Category:    "SessionSecurity",
				Risk:        RiskHigh,
				Likelihood:  4,
				Impact:      4,
				Confidence:  0.90,
				Keywords:    []string{"session", "rotation", "concurrent limits", "global sign-out", "identity platform"},
				Rationale:   "Session management gaps allow session hijacking, credential stuffing, and persistent unauthorized access.",
			},
			{
				ID:          "DOM-IDP-004",
				Description: "Identity platform must implement federation security with metadata validation, certificate pinning, and assertion encryption.",
				Component:   "IdentityPlatform",
				Category:    "Authentication",
				Risk:        RiskHigh,
				Likelihood:  3,
				Impact:      4,
				Confidence:  0.90,
				Keywords:    []string{"federation", "saml", "metadata", "assertion encryption", "identity platform"},
				Rationale:   "Federation without metadata validation and assertion encryption is vulnerable to man-in-the-middle and assertion injection.",
			},
		},
		Controls: []ControlDetail{
			{ID: "CTRL-IDP-001", Description: "Implement SSO with secure token binding", Category: "Authentication", Priority: "High"},
			{ID: "CTRL-IDP-002", Description: "Enforce phishing-resistant MFA", Category: "Authentication", Priority: "High"},
			{ID: "CTRL-IDP-003", Description: "Implement session management with global sign-out", Category: "SessionSecurity", Priority: "High"},
			{ID: "CTRL-IDP-004", Description: "Implement federation metadata validation and assertion encryption", Category: "Authentication", Priority: "High"},
		},
		RiskAmplifiers: map[string]RiskLevel{
			"SSO":        RiskCritical,
			"Federation": RiskHigh,
			"MFA":        RiskHigh,
			"Session":    RiskHigh,
		},
		ComplianceMappings: []string{"NIST 800-63", "FIDO Alliance", "ISO 27001", "GDPR"},
	}
}

// DataPlatformPack returns domain-specific assumptions for data platforms.
func DataPlatformPack() *DomainPack {
	return &DomainPack{
		Name: "DataPlatform",
		Assumptions: []Assumption{
			{
				ID:          "DOM-DAT-001",
				Description: "Data platform must implement data governance with ownership, stewardship, and quality rules.",
				Component:   "DataPlatform",
				Category:    "DataGovernance",
				Risk:        RiskHigh,
				Likelihood:  3,
				Impact:      4,
				Confidence:  0.90,
				Keywords:    []string{"data governance", "ownership", "stewardship", "quality", "data platform"},
				Rationale:   "Data without governance lacks accountability, quality assurance, and regulatory traceability.",
			},
			{
				ID:          "DOM-DAT-002",
				Description: "Data platform must implement data lineage tracking for all ingestion, transformation, and egress paths.",
				Component:   "DataPlatform",
				Category:    "DataGovernance",
				Risk:        RiskHigh,
				Likelihood:  3,
				Impact:      4,
				Confidence:  0.90,
				Keywords:    []string{"data lineage", "ingestion", "transformation", "egress", "data platform"},
				Rationale:   "Data lineage is required for impact analysis, breach investigation, and regulatory reporting.",
			},
			{
				ID:          "DOM-DAT-003",
				Description: "Data platform must enforce retention policies with automated deletion, archival, and legal hold capabilities.",
				Component:   "DataPlatform",
				Category:    "DataRetention",
				Risk:        RiskHigh,
				Likelihood:  3,
				Impact:      4,
				Confidence:  0.90,
				Keywords:    []string{"retention", "deletion", "archival", "legal hold", "data platform"},
				Rationale:   "Retention policy gaps lead to regulatory violations, excessive storage costs, and discovery risks.",
			},
			{
				ID:          "DOM-DAT-004",
				Description: "Data platform must encrypt data at rest and in transit with key management and rotation for all data tiers.",
				Component:   "DataPlatform",
				Category:    "DataProtection",
				Risk:        RiskCritical,
				Likelihood:  4,
				Impact:      5,
				Confidence:  0.95,
				Keywords:    []string{"data platform", "encryption", "key management", "rotation", "data tier"},
				Rationale:   "Data platform encryption gaps expose data across ingestion, storage, and egress to unauthorized access.",
			},
		},
		Controls: []ControlDetail{
			{ID: "CTRL-DAT-001", Description: "Implement data governance framework", Category: "DataGovernance", Priority: "High"},
			{ID: "CTRL-DAT-002", Description: "Implement data lineage tracking", Category: "DataGovernance", Priority: "High"},
			{ID: "CTRL-DAT-003", Description: "Implement retention and legal hold automation", Category: "DataRetention", Priority: "High"},
			{ID: "CTRL-DAT-004", Description: "Implement encryption and key management for all data tiers", Category: "DataProtection", Priority: "High"},
		},
		RiskAmplifiers: map[string]RiskLevel{
			"DataLake":  RiskHigh,
			"Warehouse": RiskHigh,
			"ETL":       RiskHigh,
			"Analytics": RiskHigh,
		},
		ComplianceMappings: []string{"GDPR", "CCPA", "SOX", "HIPAA", "PCI DSS"},
	}
}
