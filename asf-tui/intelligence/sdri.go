package intelligence

import (
	"fmt"
	"sort"
	"strings"
)

// ─────────────────────────────────────────────────────────────
// PHASE 1-2 — CONTROL FRAMEWORK & TAXONOMY
// ─────────────────────────────────────────────────────────────

type SDRIStrength string
type SDRIControlType string

const (
	SDRIStrengthWeak    SDRIStrength = "Weak"
	SDRIStrengthPartial SDRIStrength = "Partial"
	SDRIStrengthStrong  SDRIStrength = "Strong"

	SDRIControlPreventive SDRIControlType = "Preventive"
	SDRIControlDetective  SDRIControlType = "Detective"
	SDRIControlCorrective SDRIControlType = "Corrective"
)

type SDRIControl struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Category    string          `json:"category"`
	Description string          `json:"description"`
	ControlType SDRIControlType `json:"control_type"`
	Preventive  bool            `json:"preventive"`
	Detective   bool            `json:"detective"`
	Corrective  bool            `json:"corrective"`
	Strength    SDRIStrength    `json:"strength"`
	Evidence    []string        `json:"evidence,omitempty"`
	Coverage    string          `json:"coverage"`
	Status      string          `json:"status"`
}

type SDRIControlLibraryEntry struct {
	Name        string          `json:"name"`
	Category    string          `json:"category"`
	Description string          `json:"description"`
	ControlType SDRIControlType `json:"control_type"`
	Preventive  bool            `json:"preventive"`
	Detective   bool            `json:"detective"`
	Corrective  bool            `json:"corrective"`
}

func buildControlLibrary() map[string]SDRIControlLibraryEntry {
	return map[string]SDRIControlLibraryEntry{
		"MFA":                        {"MFA", "Authentication", "Multi-factor authentication for all user access", SDRIControlPreventive, true, false, false},
		"PasswordPolicy":             {"Password Policy", "Authentication", "Strong password complexity and rotation requirements", SDRIControlPreventive, true, false, false},
		"SessionManagement":          {"Session Management", "Authentication", "Session timeout, renewal, and invalidation controls", SDRIControlPreventive, true, false, false},
		"ConditionalAccess":          {"Conditional Access", "Authorization", "Context-aware access policies based on risk signals", SDRIControlPreventive, true, false, false},
		"RBAC":                       {"Role-Based Access Control", "Authorization", "Role-based permissions for least privilege access", SDRIControlPreventive, true, false, false},
		"ABAC":                       {"Attribute-Based Access Control", "Authorization", "Attribute-based fine-grained access policies", SDRIControlPreventive, true, false, false},
		"JustInTimeAccess":           {"Just-In-Time Access", "Authorization", "Temporary privileged access on demand", SDRIControlPreventive, true, false, false},
		"PrivilegedAccessManagement": {"Privileged Access Management", "Authorization", "Controls for managing and monitoring privileged accounts", SDRIControlPreventive, true, false, false},
		"IdentityGovernance":         {"Identity Governance", "Identity Governance", "Identity lifecycle management and certification", SDRIControlPreventive, true, false, false},
		"AccessReviews":              {"Access Reviews", "Identity Governance", "Periodic review and certification of access rights", SDRIControlDetective, false, true, false},
		"SecretsRotation":            {"Secrets Rotation", "Secrets Management", "Automatic rotation of secrets and credentials", SDRIControlPreventive, true, false, false},
		"SecretsVault":               {"Secrets Vault", "Secrets Management", "Centralized secure storage for secrets", SDRIControlPreventive, true, false, false},
		"SecretsScanning":            {"Secrets Scanning", "Secrets Management", "Scan code and configs for hardcoded secrets", SDRIControlDetective, false, true, false},
		"KeyRotation":                {"Key Rotation", "Key Management", "Automatic cryptographic key rotation", SDRIControlPreventive, true, false, false},
		"KeyAccessLogging":           {"Key Access Logging", "Key Management", "Logging of all key access and usage", SDRIControlDetective, false, true, false},
		"SeparationOfDuties":         {"Separation of Duties", "Key Management", "Split key management responsibilities across roles", SDRIControlPreventive, true, false, false},
		"NetworkSegmentation":        {"Network Segmentation", "Network Security", "Division of network into isolated segments", SDRIControlPreventive, true, false, false},
		"FirewallRules":              {"Firewall Rules", "Network Security", "Network firewall rules to restrict traffic", SDRIControlPreventive, true, false, false},
		"TLSEncryption":              {"TLS Encryption", "Network Security", "Encrypted communication using TLS", SDRIControlPreventive, true, false, false},
		"VPNAccess":                  {"VPN Access", "Network Security", "Encrypted remote access via VPN", SDRIControlPreventive, true, false, false},
		"IntrusionDetection":         {"Intrusion Detection", "Network Security", "Network-based intrusion detection system", SDRIControlDetective, false, true, false},
		"AuditLogging":               {"Audit Logging", "Logging", "Comprehensive audit logging of security events", SDRIControlDetective, false, true, false},
		"SIEMIntegration":            {"SIEM Integration", "Monitoring", "Centralized security event monitoring with SIEM", SDRIControlDetective, false, true, false},
		"RealTimeAlerting":           {"Real-Time Alerting", "Alerting", "Real-time security alerting and notification", SDRIControlDetective, false, true, false},
		"IncidentResponsePlan":       {"Incident Response Plan", "Incident Response", "Documented incident response procedures", SDRIControlCorrective, false, false, true},
		"IncidentResponseTeam":       {"Incident Response Team", "Incident Response", "Dedicated incident response team", SDRIControlCorrective, false, false, true},
		"AutomatedBackup":            {"Automated Backup", "Backup", "Automated regular data backup", SDRIControlCorrective, false, false, true},
		"DisasterRecoveryPlan":       {"Disaster Recovery Plan", "Disaster Recovery", "Documented disaster recovery procedures", SDRIControlCorrective, false, false, true},
		"DataEncryptionAtRest":       {"Data Encryption at Rest", "Data Protection", "Encryption of stored data", SDRIControlPreventive, true, false, false},
		"DataEncryptionInTransit":    {"Data Encryption in Transit", "Data Protection", "Encryption of data during transmission", SDRIControlPreventive, true, false, false},
		"DataClassification":         {"Data Classification", "Data Protection", "Classification of data by sensitivity level", SDRIControlPreventive, true, false, false},
		"DLPControls":                {"DLP Controls", "Data Protection", "Data loss prevention monitoring and controls", SDRIControlDetective, false, true, false},
		"DataRetentionPolicy":        {"Data Retention Policy", "Data Protection", "Policy for data retention and disposal", SDRIControlPreventive, true, false, false},
		"PrivacyImpactAssessment":    {"Privacy Impact Assessment", "Privacy", "Assessment of privacy impacts for new systems", SDRIControlPreventive, true, false, false},
		"ConsentManagement":          {"Consent Management", "Privacy", "Manage user consent for data processing", SDRIControlPreventive, true, false, false},
		"DataSubjectAccess":          {"Data Subject Access", "Privacy", "Process for handling data subject access requests", SDRIControlCorrective, false, false, true},
		"ThirdPartyDueDiligence":     {"Third Party Due Diligence", "Third Party", "Security assessment of third-party vendors", SDRIControlPreventive, true, false, false},
		"VendorRiskManagement":       {"Vendor Risk Management", "Third Party", "Ongoing vendor risk monitoring program", SDRIControlDetective, false, true, false},
		"CloudSecurityPosture":       {"Cloud Security Posture Management", "Cloud Security", "Continuous cloud security configuration assessment", SDRIControlDetective, false, true, false},
		"CloudAccessBroker":          {"Cloud Access Security Broker", "Cloud Security", "Security policy enforcement for cloud services", SDRIControlPreventive, true, false, false},
		"ContainerImageScanning":     {"Container Image Scanning", "Container Security", "Vulnerability scanning of container images", SDRIControlDetective, false, true, false},
		"ContainerRuntimeSecurity":   {"Container Runtime Security", "Container Security", "Runtime security monitoring for containers", SDRIControlDetective, false, true, false},
		"K8sRBAC":                    {"Kubernetes RBAC", "Kubernetes", "Kubernetes role-based access control", SDRIControlPreventive, true, false, false},
		"K8sAdmissionControllers":    {"Kubernetes Admission Controllers", "Kubernetes", "Admission control policies for Kubernetes", SDRIControlPreventive, true, false, false},
		"K8sNetworkPolicies":         {"Kubernetes Network Policies", "Kubernetes", "Network policy enforcement in Kubernetes", SDRIControlPreventive, true, false, false},
		"K8sPodSecurity":             {"Kubernetes Pod Security", "Kubernetes", "Pod security standards and policies", SDRIControlPreventive, true, false, false},
		"EndpointDetectionResponse":  {"Endpoint Detection & Response", "Endpoint Security", "EDR agent for endpoint threat detection", SDRIControlDetective, false, true, false},
		"Antivirus":                  {"Antivirus", "Endpoint Security", "Antivirus and anti-malware protection", SDRIControlPreventive, true, false, false},
		"DeviceCompliance":           {"Device Compliance", "Endpoint Security", "Device compliance and health checks", SDRIControlPreventive, true, false, false},
		"ChangeManagement":           {"Change Management", "Change Management", "Formal change management and approval process", SDRIControlPreventive, true, false, false},
		"VulnerabilityScanning":      {"Vulnerability Scanning", "Vulnerability Management", "Regular vulnerability scanning", SDRIControlDetective, false, true, false},
		"PatchManagement":            {"Patch Management", "Vulnerability Management", "Systematic patch management program", SDRIControlCorrective, false, false, true},
		"PenetrationTesting":         {"Penetration Testing", "Vulnerability Management", "Regular penetration testing", SDRIControlDetective, false, true, false},
	}
}

// ─────────────────────────────────────────────────────────────
// PHASE 3 — EXPECTED CONTROL ENGINE
// ─────────────────────────────────────────────────────────────

func buildExpectedControlMap() map[string][]string {
	return map[string][]string{
		"database":     {"Access Control", "Data Encryption at Rest", "Audit Logging", "Automated Backup", "Vulnerability Scanning"},
		"db":           {"Access Control", "Data Encryption at Rest", "Audit Logging", "Automated Backup", "Vulnerability Scanning"},
		"auth":         {"MFA", "Session Management", "Conditional Access", "Password Policy", "Audit Logging"},
		"identity":     {"MFA", "Session Management", "Conditional Access", "Password Policy", "Audit Logging"},
		"api":          {"TLS Encryption", "Authentication", "Rate Limiting", "Audit Logging", "Input Validation"},
		"gateway":      {"TLS Encryption", "Authentication", "Rate Limiting", "Audit Logging", "Firewall Rules"},
		"web":          {"TLS Encryption", "Authentication", "Session Management", "Audit Logging", "Vulnerability Scanning"},
		"app":          {"TLS Encryption", "Authentication", "Session Management", "Audit Logging", "Vulnerability Scanning"},
		"kms":          {"Key Rotation", "Key Access Logging", "Separation of Duties", "Audit Logging", "Access Control"},
		"secrets":      {"Secrets Rotation", "Secrets Vault", "Secrets Scanning", "Access Control", "Audit Logging"},
		"vault":        {"Secrets Rotation", "Secrets Vault", "Secrets Scanning", "Access Control", "Audit Logging"},
		"cache":        {"Access Control", "TLS Encryption", "Audit Logging", "Data Encryption at Rest", "Network Segmentation"},
		"queue":        {"Access Control", "TLS Encryption", "Audit Logging", "Data Encryption at Rest", "Network Segmentation"},
		"storage":      {"Access Control", "Data Encryption at Rest", "Data Encryption in Transit", "Audit Logging", "Backup"},
		"bucket":       {"Access Control", "Data Encryption at Rest", "Data Encryption in Transit", "Audit Logging", "Backup"},
		"loadbalancer": {"TLS Encryption", "Firewall Rules", "DDoS Protection", "Audit Logging", "Network Segmentation"},
		"lb":           {"TLS Encryption", "Firewall Rules", "DDoS Protection", "Audit Logging", "Network Segmentation"},
		"dns":          {"DNSSEC", "Access Control", "Audit Logging", "Rate Limiting", "Redundancy"},
		"kubernetes":   {"K8s RBAC", "K8s Admission Controllers", "K8s Network Policies", "K8s Pod Security", "Audit Logging"},
		"k8s":          {"K8s RBAC", "K8s Admission Controllers", "K8s Network Policies", "K8s Pod Security", "Audit Logging"},
		"cluster":      {"K8s RBAC", "K8s Admission Controllers", "K8s Network Policies", "K8s Pod Security", "Audit Logging"},
		"container":    {"Container Image Scanning", "Container Runtime Security", "Vulnerability Scanning", "Audit Logging", "Access Control"},
		"vpn":          {"MFA", "VPN Access", "Session Management", "Audit Logging", "Device Compliance"},
		"firewall":     {"Firewall Rules", "Network Segmentation", "Intrusion Detection", "Audit Logging", "Change Management"},
		"logging":      {"Audit Logging", "SIEM Integration", "Data Encryption at Rest", "Access Control", "Backup"},
		"monitoring":   {"Real-Time Alerting", "SIEM Integration", "Incident Response Plan", "Audit Logging", "Access Control"},
		"certificate":  {"Key Rotation", "Key Access Logging", "Separation of Duties", "Certificate Transparency", "Access Control"},
		"certs":        {"Key Rotation", "Key Access Logging", "Separation of Duties", "Certificate Transparency", "Access Control"},
		"cdn":          {"TLS Encryption", "DDoS Protection", "Access Control", "Audit Logging", "Network Segmentation"},
		"email":        {"TLS Encryption", "DKIM", "SPF", "DMARC", "Audit Logging"},
		"smtp":         {"TLS Encryption", "Authentication", "Rate Limiting", "Audit Logging", "Access Control"},
		"thirdparty":   {"Third Party Due Diligence", "Vendor Risk Management", "Access Control", "Audit Logging", "Data Encryption in Transit"},
		"endpoint":     {"Endpoint Detection & Response", "Antivirus", "Device Compliance", "Patch Management", "Vulnerability Scanning"},
		"device":       {"Endpoint Detection & Response", "Antivirus", "Device Compliance", "Patch Management", "Vulnerability Scanning"},
		"backup":       {"Automated Backup", "Data Encryption at Rest", "Access Control", "Disaster Recovery Plan", "Audit Logging"},
		"dr":           {"Disaster Recovery Plan", "Automated Backup", "Incident Response Plan", "Data Replication", "Testing Schedule"},
	}
}

// ─────────────────────────────────────────────────────────────
// PHASE 4 — OBSERVED CONTROL ENGINE
// ─────────────────────────────────────────────────────────────

func extractObservedControls(arch *ArchDescription, assumptions []Assumption, controls []ControlDetail) map[string][]string {
	observed := make(map[string][]string)

	for _, comp := range arch.Components {
		key := strings.ToLower(comp.Label)
		var found []string

		// Check architecture security controls metadata
		if arch.SecurityControls != nil {
			if compControls, ok := arch.SecurityControls[comp.ID]; ok {
				found = append(found, compControls...)
			}
		}

		// Check control details from domain packs
		for _, ctrl := range controls {
			if strings.EqualFold(ctrl.Component, comp.ID) || strings.EqualFold(ctrl.Component, comp.Label) {
				if !containsControl(found, ctrl.Name) {
					found = append(found, ctrl.Name)
				}
			}
		}

		// Check for evidence of controls in assumption evidence
		for _, a := range assumptions {
			if !strings.Contains(strings.ToLower(a.Component), key) &&
				!strings.Contains(strings.ToLower(a.Description), key) {
				continue
			}
			for _, ev := range a.EvidenceSources {
				evLower := strings.ToLower(ev)
				ctrlName := inferControlFromEvidence(evLower)
				if ctrlName != "" && !containsControl(found, ctrlName) {
					found = append(found, ctrlName)
				}
			}
		}

		if len(found) > 0 {
			observed[comp.Label] = found
		}
	}

	return observed
}

func containsControl(list []string, name string) bool {
	for _, c := range list {
		if strings.EqualFold(c, name) {
			return true
		}
	}
	return false
}

func inferControlFromEvidence(ev string) string {
	if strings.Contains(ev, "mfa") || strings.Contains(ev, "multi-factor") || strings.Contains(ev, "2fa") {
		return "MFA"
	}
	if strings.Contains(ev, "encrypt") && (strings.Contains(ev, "rest") || strings.Contains(ev, "at-rest") || strings.Contains(ev, "aes")) {
		return "Data Encryption at Rest"
	}
	if strings.Contains(ev, "encrypt") && (strings.Contains(ev, "transit") || strings.Contains(ev, "tls") || strings.Contains(ev, "https")) {
		return "Data Encryption in Transit"
	}
	if strings.Contains(ev, "audit") || strings.Contains(ev, "log") {
		return "Audit Logging"
	}
	if strings.Contains(ev, "backup") {
		return "Automated Backup"
	}
	if strings.Contains(ev, "firewall") || strings.Contains(ev, "acl") || strings.Contains(ev, "whitelist") {
		return "Firewall Rules"
	}
	if strings.Contains(ev, "rbac") || strings.Contains(ev, "role") {
		return "RBAC"
	}
	if strings.Contains(ev, "vpn") {
		return "VPN Access"
	}
	if strings.Contains(ev, "scan") || strings.Contains(ev, "vulnerability") {
		return "Vulnerability Scanning"
	}
	if strings.Contains(ev, "siem") || strings.Contains(ev, "splunk") || strings.Contains(ev, "elk") {
		return "SIEM Integration"
	}
	if strings.Contains(ev, "patch") {
		return "Patch Management"
	}
	if strings.Contains(ev, "secret") || strings.Contains(ev, "vault") {
		return "Secrets Vault"
	}
	if strings.Contains(ev, "rotation") {
		return "Secrets Rotation"
	}
	if strings.Contains(ev, "iam") || strings.Contains(ev, "identity") {
		return "Identity Governance"
	}
	if strings.Contains(ev, "incident") || strings.Contains(ev, "ir") {
		return "Incident Response Plan"
	}
	if strings.Contains(ev, "third") || strings.Contains(ev, "vendor") {
		return "Third Party Due Diligence"
	}
	if strings.Contains(ev, "container") || strings.Contains(ev, "docker") {
		return "Container Image Scanning"
	}
	if strings.Contains(ev, "k8s") || strings.Contains(ev, "kubernetes") {
		return "K8s RBAC"
	}
	if strings.Contains(ev, "edr") || strings.Contains(ev, "endpoint") {
		return "Endpoint Detection & Response"
	}
	if strings.Contains(ev, "antivirus") || strings.Contains(ev, "malware") {
		return "Antivirus"
	}
	return ""
}

// ─────────────────────────────────────────────────────────────
// PHASE 5 — CONTROL GAP DETECTION
// ─────────────────────────────────────────────────────────────

func detectControlGaps(expected, observed map[string][]string) map[string][]string {
	library := buildControlLibrary()
	missing := make(map[string][]string)

	allComponents := make(map[string]bool)
	for comp := range expected {
		allComponents[comp] = true
	}
	for comp := range observed {
		allComponents[comp] = true
	}

	for comp := range allComponents {
		exp := normalizeControlNames(expected[comp], library)
		obs := normalizeControlNames(observed[comp], library)
		obsSet := make(map[string]bool)
		for _, c := range obs {
			obsSet[c] = true
		}
		var compMissing []string
		for _, expCtrl := range exp {
			if !obsSet[expCtrl] {
				compMissing = append(compMissing, expCtrl)
			}
		}
		if len(compMissing) > 0 {
			missing[comp] = compMissing
		}
	}

	return missing
}

func normalizeControlNames(controls []string, library map[string]SDRIControlLibraryEntry) []string {
	var result []string
	for _, c := range controls {
		found := false
		lower := strings.ToLower(c)
		for key, entry := range library {
			if strings.ToLower(key) == lower || strings.ToLower(entry.Name) == lower {
				result = append(result, entry.Name)
				found = true
				break
			}
		}
		if !found {
			result = append(result, c)
		}
	}
	return result
}

// ─────────────────────────────────────────────────────────────
// PHASE 6 — COVERAGE SCORING
// ─────────────────────────────────────────────────────────────

type SDRICoverage struct {
	Category string  `json:"category"`
	Expected int     `json:"expected"`
	Observed int     `json:"observed"`
	Coverage float64 `json:"coverage"`
	Level    string  `json:"level"`
}

func computeCoverageScores(expected, observed map[string][]string) []SDRICoverage {
	library := buildControlLibrary()
	categoryExpected := make(map[string]int)
	categoryObserved := make(map[string]int)

	allComponentControls := make(map[string]map[string]bool)
	for comp, ctrls := range expected {
		if allComponentControls[comp] == nil {
			allComponentControls[comp] = make(map[string]bool)
		}
		for _, c := range normalizeControlNames(ctrls, library) {
			allComponentControls[comp][c] = true
			cat := lookupCategory(c, library)
			categoryExpected[cat]++
		}
	}
	for comp, ctrls := range observed {
		if allComponentControls[comp] == nil {
			allComponentControls[comp] = make(map[string]bool)
		}
		for _, c := range normalizeControlNames(ctrls, library) {
			cat := lookupCategory(c, library)
			if allComponentControls[comp][c] {
				categoryObserved[cat]++
			} else {
				allComponentControls[comp][c] = true
				categoryObserved[cat]++
			}
		}
	}

	var results []SDRICoverage
	categories := make([]string, 0, len(categoryExpected))
	for cat := range categoryExpected {
		categories = append(categories, cat)
	}
	sort.Strings(categories)

	for _, cat := range categories {
		exp := categoryExpected[cat]
		obs := categoryObserved[cat]
		cov := 0.0
		if exp > 0 {
			cov = float64(obs) / float64(exp) * 100
		}
		results = append(results, SDRICoverage{
			Category: cat,
			Expected: exp,
			Observed: obs,
			Coverage: cov,
			Level:    coverageLevel(cov),
		})
	}

	return results
}

func coverageLevel(cov float64) string {
	switch {
	case cov >= 90:
		return "Excellent"
	case cov >= 75:
		return "Strong"
	case cov >= 50:
		return "Good"
	case cov >= 25:
		return "Fair"
	default:
		return "Poor"
	}
}

func lookupCategory(controlName string, library map[string]SDRIControlLibraryEntry) string {
	for key, entry := range library {
		if strings.EqualFold(key, controlName) || strings.EqualFold(entry.Name, controlName) {
			return entry.Category
		}
	}
	return "General"
}

// ─────────────────────────────────────────────────────────────
// PHASE 7 — CONTROL EFFECTIVENESS
// ─────────────────────────────────────────────────────────────

func assessControlStrength(obsControls map[string][]string, assumptions []Assumption) map[string]SDRIStrength {
	strength := make(map[string]SDRIStrength)

	for comp, ctrls := range obsControls {
		for _, ctrl := range ctrls {
			key := comp + ":" + ctrl
			s := assessSingleControl(ctrl, comp, assumptions)
			if existing, ok := strength[key]; !ok || s == SDRIStrengthWeak || (s == SDRIStrengthPartial && existing == SDRIStrengthStrong) {
				strength[key] = s
			}
			strength[key] = s
		}
	}

	return strength
}

func assessSingleControl(controlName, component string, assumptions []Assumption) SDRIStrength {
	weakEvidence := 0
	partialEvidence := 0
	strongEvidence := 0

	for _, a := range assumptions {
		if !strings.Contains(strings.ToLower(a.Component), strings.ToLower(component)) &&
			!strings.Contains(strings.ToLower(a.Description), strings.ToLower(component)) {
			continue
		}

		ctrlLower := strings.ToLower(controlName)
		for _, ev := range a.EvidenceSources {
			evLower := strings.ToLower(ev)

			if !strings.Contains(evLower, ctrlLower) {
				continue
			}

			if strings.Contains(evLower, "strong") || strings.Contains(evLower, "comprehensive") ||
				strings.Contains(evLower, "automated") || strings.Contains(evLower, "enforced") {
				strongEvidence++
			} else if strings.Contains(evLower, "partial") || strings.Contains(evLower, "manual") ||
				strings.Contains(evLower, "limited") || strings.Contains(evLower, "basic") {
				partialEvidence++
			} else {
				weakEvidence++
			}
		}

		for _, ev := range a.EvidenceSources {
			evLower := strings.ToLower(ev)
			if strings.Contains(evLower, "no "+ctrlLower) || strings.Contains(evLower, "missing "+ctrlLower) ||
				strings.Contains(evLower, "absent") || strings.Contains(evLower, "not implemented") {
				return SDRIStrengthWeak
			}
		}
	}

	if strongEvidence > 0 {
		return SDRIStrengthStrong
	}
	if partialEvidence > 0 || weakEvidence > 1 {
		return SDRIStrengthPartial
	}
	if weakEvidence > 0 {
		return SDRIStrengthWeak
	}

	return SDRIStrengthPartial
}

// ─────────────────────────────────────────────────────────────
// PHASE 8 — SECURITY DESIGN REVIEW FINDINGS
// ─────────────────────────────────────────────────────────────

type SDRIFinding struct {
	ID                 string   `json:"id"`
	Title              string   `json:"title"`
	Description        string   `json:"description"`
	Severity           string   `json:"severity"`
	Category           string   `json:"category"`
	AffectedComponents []string `json:"affected_components,omitempty"`
	AffectedControls   []string `json:"affected_controls,omitempty"`
	BusinessImpact     string   `json:"business_impact"`
	Recommendation     string   `json:"recommendation"`
	Reasoning          string   `json:"reasoning"`
}

func generateDesignFindings(missing map[string][]string, weaknesses []SDRIArchitecturalWeakness,
	attackPaths []AttackPath, coverage []SDRICoverage) []SDRIFinding {

	var findings []SDRIFinding
	findingID := 1

	// Findings from missing controls
	for comp, missingCtrls := range missing {
		for _, ctrl := range missingCtrls {
			severity := severityForMissingControl(ctrl)
			impact := impactForMissingControl(ctrl)
			findings = append(findings, SDRIFinding{
				ID:                 fmt.Sprintf("SDR-F-%03d", findingID),
				Title:              fmt.Sprintf("Missing %s on %s", ctrl, comp),
				Description:        fmt.Sprintf("The component %s is missing %s, which is an expected security control.", comp, ctrl),
				Severity:           severity,
				Category:           categorizeControl(ctrl),
				AffectedComponents: []string{comp},
				AffectedControls:   []string{ctrl},
				BusinessImpact:     impact,
				Recommendation:     fmt.Sprintf("Implement %s on %s to address this control gap.", ctrl, comp),
				Reasoning:          fmt.Sprintf("%s is required for %s based on security best practices and expected control patterns.", ctrl, comp),
			})
			findingID++
		}
	}

	// Findings from architectural weaknesses
	for _, w := range weaknesses {
		findings = append(findings, SDRIFinding{
			ID:                 fmt.Sprintf("SDR-F-%03d", findingID),
			Title:              w.Pattern + " - " + strings.Join(w.Components, ", "),
			Description:        w.Description,
			Severity:           w.Severity,
			Category:           "Architectural Weakness",
			AffectedComponents: w.Components,
			BusinessImpact:     w.Impact,
			Recommendation:     w.Recommendation,
			Reasoning:          w.Description,
		})
		findingID++
	}

	// Findings from poor coverage
	for _, c := range coverage {
		if c.Level == "Poor" || c.Level == "Fair" {
			findings = append(findings, SDRIFinding{
				ID:             fmt.Sprintf("SDR-F-%03d", findingID),
				Title:          fmt.Sprintf("Weak Control Coverage in %s", c.Category),
				Description:    fmt.Sprintf("Control coverage in %s is %s (%.1f%%). Only %d of %d expected controls are observed.", c.Category, c.Level, c.Coverage, c.Observed, c.Expected),
				Severity:       severityForCoverage(c.Level),
				Category:       c.Category,
				BusinessImpact: fmt.Sprintf("Weak %s controls increase the attack surface and reduce the organization's ability to detect and respond to threats.", c.Category),
				Recommendation: fmt.Sprintf("Improve control coverage in %s from %.1f%% to at least 75%%.", c.Category, c.Coverage),
				Reasoning:      fmt.Sprintf("Only %d of %d expected controls are present in %s, leaving significant gaps.", c.Observed, c.Expected, c.Category),
			})
			findingID++
		}
	}

	if len(findings) == 0 {
		return findings
	}

	sort.SliceStable(findings, func(i, j int) bool {
		return severityRank(findings[i].Severity) < severityRank(findings[j].Severity)
	})

	return findings
}

func severityRank(s string) int {
	switch s {
	case "Critical":
		return 0
	case "High":
		return 1
	case "Medium":
		return 2
	case "Low":
		return 3
	}
	return 4
}

func severityForMissingControl(ctrl string) string {
	high := map[string]bool{
		"MFA": true, "Data Encryption at Rest": true, "Data Encryption in Transit": true,
		"RBAC": true, "Secrets Vault": true, "Secrets Rotation": true, "Key Rotation": true,
		"Firewall Rules": true, "K8s RBAC": true, "K8s Admission Controllers": true,
		"Audit Logging": true,
	}
	critical := map[string]bool{
		"Access Control": true, "Authentication": true,
	}
	if critical[ctrl] {
		return "Critical"
	}
	if high[ctrl] {
		return "High"
	}
	return "Medium"
}

func impactForMissingControl(ctrl string) string {
	impacts := map[string]string{
		"MFA":                        "Increased risk of account compromise and unauthorized access",
		"Data Encryption at Rest":    "Potential data exposure if storage is compromised",
		"Data Encryption in Transit": "Data vulnerable to interception during transmission",
		"RBAC":                       "Users may have excessive permissions leading to privilege escalation",
		"Secrets Vault":              "Secrets stored insecurely, increasing breach impact",
		"Secrets Rotation":           "Compromised secrets remain valid indefinitely",
		"Key Rotation":               "Cryptographic keys may be compromised without rotation",
		"Audit Logging":              "Inability to detect and investigate security incidents",
		"Access Control":             "Complete lack of access restrictions, critical vulnerability",
	}
	if impact, ok := impacts[ctrl]; ok {
		return impact
	}
	return "Control gap increases security risk for affected components"
}

func categorizeControl(ctrl string) string {
	library := buildControlLibrary()
	for _, entry := range library {
		if entry.Name == ctrl {
			return entry.Category
		}
	}
	return "General"
}

func severityForCoverage(level string) string {
	switch level {
	case "Poor":
		return "High"
	case "Fair":
		return "Medium"
	}
	return "Low"
}

// ─────────────────────────────────────────────────────────────
// PHASE 9 — ARCHITECTURAL WEAKNESS DETECTION
// ─────────────────────────────────────────────────────────────

type SDRIArchitecturalWeakness struct {
	ID             string   `json:"id"`
	Pattern        string   `json:"pattern"`
	Description    string   `json:"description"`
	Severity       string   `json:"severity"`
	Components     []string `json:"components,omitempty"`
	Impact         string   `json:"impact"`
	Recommendation string   `json:"recommendation"`
}

func detectArchitecturalWeaknesses(arch *ArchDescription, assumptions []Assumption) []SDRIArchitecturalWeakness {
	var weaknesses []SDRIArchitecturalWeakness
	weakID := 1

	compLabels := make([]string, len(arch.Components))
	compMap := make(map[string]string)
	for i, c := range arch.Components {
		compLabels[i] = c.Label
		compMap[c.ID] = c.Label
	}

	// Pattern: Single Point of Failure - one auth provider, no backup identity
	authCount := 0
	var authComponents []string
	for _, c := range arch.Components {
		lower := strings.ToLower(c.Label)
		if strings.Contains(lower, "auth") || strings.Contains(lower, "identity") || strings.Contains(lower, "sso") || strings.Contains(lower, "idp") {
			authCount++
			authComponents = append(authComponents, c.Label)
		}
	}
	if authCount == 1 {
		weaknesses = append(weaknesses, SDRIArchitecturalWeakness{
			ID:             fmt.Sprintf("SDR-W-%03d", weakID),
			Pattern:        "Single Point of Failure",
			Description:    fmt.Sprintf("%s is the only authentication provider with no backup identity source. If it fails, the entire system is unavailable.", authComponents[0]),
			Severity:       "Critical",
			Components:     authComponents,
			Impact:         "Complete authentication failure and system downtime if the single provider is compromised or unavailable",
			Recommendation: "Implement a backup identity provider or failover mechanism for authentication",
		})
		weakID++
	}

	// Pattern: Excessive Trust / Flat Access - no RBAC or segmentation
	hasRBAC := false
	hasSegmentation := false
	for _, a := range assumptions {
		for _, ev := range a.EvidenceSources {
			evLower := strings.ToLower(ev)
			if strings.Contains(evLower, "rbac") || strings.Contains(evLower, "role") {
				hasRBAC = true
			}
			if strings.Contains(evLower, "segment") || strings.Contains(evLower, "network policy") {
				hasSegmentation = true
			}
		}
	}
	if !hasRBAC && len(arch.Components) > 2 {
		weaknesses = append(weaknesses, SDRIArchitecturalWeakness{
			ID:             fmt.Sprintf("SDR-W-%03d", weakID),
			Pattern:        "Excessive Trust / Flat Access",
			Description:    "No RBAC or role-based access controls detected across the architecture. Components may have excessive trust relationships.",
			Severity:       "High",
			Components:     compLabels,
			Impact:         "Any compromised component can access all other components without restriction",
			Recommendation: "Implement RBAC and network segmentation to enforce least privilege access",
		})
		weakID++
	}
	if !hasSegmentation && len(arch.Components) > 3 {
		weaknesses = append(weaknesses, SDRIArchitecturalWeakness{
			ID:             fmt.Sprintf("SDR-W-%03d", weakID),
			Pattern:        "Flat Network Topology",
			Description:    "No network segmentation or network policies detected. All components appear to be on a flat network.",
			Severity:       "High",
			Components:     compLabels,
			Impact:         "Lateral movement is unhindered once any component is compromised",
			Recommendation: "Implement network segmentation and firewall rules to isolate components by trust zone",
		})
		weakID++
	}

	// Pattern: Secrets Exposure Risk
	hasSecretsVault := false
	for _, a := range assumptions {
		for _, ev := range a.EvidenceSources {
			evLower := strings.ToLower(ev)
			if strings.Contains(evLower, "vault") || strings.Contains(evLower, "secret") {
				hasSecretsVault = true
			}
		}
	}
	if !hasSecretsVault && len(arch.Components) > 0 {
		weaknesses = append(weaknesses, SDRIArchitecturalWeakness{
			ID:             fmt.Sprintf("SDR-W-%03d", weakID),
			Pattern:        "Secrets Exposure Risk",
			Description:    "No centralized secrets vault detected. Secrets may be stored in configuration files or code.",
			Severity:       "High",
			Components:     compLabels,
			Impact:         "Secrets can be exposed through code leaks, config errors, or compromised repositories",
			Recommendation: "Implement a secrets vault (e.g., HashiCorp Vault, AWS Secrets Manager) for all secrets",
		})
		weakID++
	}

	// Pattern: Weak Monitoring
	hasLogging := false
	hasSIEM := false
	hasAlerting := false
	for _, a := range assumptions {
		for _, ev := range a.EvidenceSources {
			evLower := strings.ToLower(ev)
			if strings.Contains(evLower, "audit") || strings.Contains(evLower, "log") {
				hasLogging = true
			}
			if strings.Contains(evLower, "siem") || strings.Contains(evLower, "splunk") || strings.Contains(evLower, "elk") {
				hasSIEM = true
			}
			if strings.Contains(evLower, "alert") || strings.Contains(evLower, "notif") {
				hasAlerting = true
			}
		}
	}
	if !hasLogging {
		weaknesses = append(weaknesses, SDRIArchitecturalWeakness{
			ID:             fmt.Sprintf("SDR-W-%03d", weakID),
			Pattern:        "Weak Monitoring",
			Description:    "No audit logging detected across the architecture. Security events cannot be traced.",
			Severity:       "Critical",
			Components:     compLabels,
			Impact:         "Complete inability to detect, investigate, or audit security incidents",
			Recommendation: "Implement comprehensive audit logging across all components",
		})
		weakID++
	} else if !hasSIEM && !hasAlerting {
		weaknesses = append(weaknesses, SDRIArchitecturalWeakness{
			ID:             fmt.Sprintf("SDR-W-%03d", weakID),
			Pattern:        "Weak Monitoring",
			Description:    "Logging is present but no SIEM or alerting detected. Logs exist but are not actively monitored.",
			Severity:       "High",
			Components:     compLabels,
			Impact:         "Security incidents may go undetected despite logs being collected",
			Recommendation: "Implement SIEM integration and real-time alerting on security events",
		})
		weakID++
	}

	// Pattern: Weak Key Management
	hasKeyRotation := false
	for _, a := range assumptions {
		for _, ev := range a.EvidenceSources {
			evLower := strings.ToLower(ev)
			if strings.Contains(evLower, "rotation") || strings.Contains(evLower, "key rotation") {
				hasKeyRotation = true
			}
		}
	}
	for _, c := range arch.Components {
		lower := strings.ToLower(c.Label)
		if (strings.Contains(lower, "kms") || strings.Contains(lower, "key") || strings.Contains(lower, "certificate") || strings.Contains(lower, "crypto")) && !hasKeyRotation {
			weaknesses = append(weaknesses, SDRIArchitecturalWeakness{
				ID:             fmt.Sprintf("SDR-W-%03d", weakID),
				Pattern:        "Weak Key Management",
				Description:    fmt.Sprintf("%s is present but no key rotation controls detected. Cryptographic keys may not be rotated.", c.Label),
				Severity:       "High",
				Components:     []string{c.Label},
				Impact:         "Compromised cryptographic keys remain valid indefinitely",
				Recommendation: "Implement automated key rotation and key access logging for " + c.Label,
			})
			weakID++
		}
	}

	// Pattern: Weak Third Party Governance
	if len(arch.Components) > 0 {
		hasThirdPartyControl := false
		for _, a := range assumptions {
			for _, ev := range a.EvidenceSources {
				evLower := strings.ToLower(ev)
				if strings.Contains(evLower, "third") || strings.Contains(evLower, "vendor") || strings.Contains(evLower, "supply") {
					hasThirdPartyControl = true
				}
			}
		}
		for _, c := range arch.Components {
			if strings.Contains(strings.ToLower(c.Label), "third") || strings.Contains(strings.ToLower(c.Label), "vendor") || strings.Contains(strings.ToLower(c.Label), "external") {
				if !hasThirdPartyControl {
					weaknesses = append(weaknesses, SDRIArchitecturalWeakness{
						ID:             fmt.Sprintf("SDR-W-%03d", weakID),
						Pattern:        "Weak Third Party Governance",
						Description:    fmt.Sprintf("%s is a third-party component but no vendor risk management controls detected.", c.Label),
						Severity:       "Medium",
						Components:     []string{c.Label},
						Impact:         "Third-party security posture is unknown, creating supply chain risk",
						Recommendation: "Implement third-party due diligence and vendor risk management for " + c.Label,
					})
					weakID++
				}
			}
		}
	}

	return weaknesses
}

// ─────────────────────────────────────────────────────────────
// PHASE 10 — REMEDIATION PRIORITIZATION
// ─────────────────────────────────────────────────────────────

type SDRIRemediation struct {
	ID                 string   `json:"id"`
	Priority           int      `json:"priority"`
	Description        string   `json:"description"`
	RiskScore          float64  `json:"risk_score"`
	BusinessImpact     string   `json:"business_impact"`
	Effort             string   `json:"effort"`
	Category           string   `json:"category"`
	Recommendation     string   `json:"recommendation"`
	AffectedComponents []string `json:"affected_components,omitempty"`
}

func prioritizeRemediations(findings []SDRIFinding, weaknesses []SDRIArchitecturalWeakness,
	attackPaths []AttackPath, coverage []SDRICoverage) []SDRIRemediation {

	var remediations []SDRIRemediation

	// Collect all attack path risks per component
	compRisk := make(map[string]float64)
	compThreatCount := make(map[string]int)
	for _, ap := range attackPaths {
		for _, comp := range ap.AffectedComponents {
			if ap.RiskScore > compRisk[comp] {
				compRisk[comp] = ap.RiskScore
			}
			compThreatCount[comp]++
		}
	}

	for _, f := range findings {
		maxRisk := 0.0
		for _, comp := range f.AffectedComponents {
			if compRisk[comp] > maxRisk {
				maxRisk = compRisk[comp]
			}
		}
		if maxRisk == 0 && len(attackPaths) > 0 {
			maxRisk = 0.3
		}

		score := severityScore(f.Severity) + maxRisk
		effort := estimateEffort(f.Category, f.Severity)

		remediations = append(remediations, SDRIRemediation{
			ID:                 f.ID,
			Description:        f.Title,
			RiskScore:          score,
			BusinessImpact:     f.BusinessImpact,
			Effort:             effort,
			Category:           f.Category,
			Recommendation:     f.Recommendation,
			AffectedComponents: f.AffectedComponents,
		})
	}

	sort.SliceStable(remediations, func(i, j int) bool {
		return remediations[i].RiskScore > remediations[j].RiskScore
	})

	for i := range remediations {
		remediations[i].Priority = i + 1
	}

	if len(remediations) > 20 {
		remediations = remediations[:20]
	}

	return remediations
}

func severityScore(severity string) float64 {
	switch severity {
	case "Critical":
		return 0.9
	case "High":
		return 0.7
	case "Medium":
		return 0.5
	case "Low":
		return 0.3
	}
	return 0.1
}

func estimateEffort(category, severity string) string {
	if severity == "Critical" || severity == "High" {
		return "High"
	}
	if category == "Architectural Weakness" {
		return "High"
	}
	return "Medium"
}

// ─────────────────────────────────────────────────────────────
// PHASE 11 — EXECUTIVE SUMMARY ENGINE
// ─────────────────────────────────────────────────────────────

func buildExecutiveSummary(findings []SDRIFinding, missing map[string][]string,
	weaknesses []SDRIArchitecturalWeakness, coverage []SDRICoverage, remediations []SDRIRemediation) string {

	criticalCount := 0
	highCount := 0
	mediumCount := 0
	lowCount := 0
	catCoverage := make(map[string]float64)
	var worstCategory string
	worstCov := 100.0

	for _, f := range findings {
		switch f.Severity {
		case "Critical":
			criticalCount++
		case "High":
			highCount++
		case "Medium":
			mediumCount++
		case "Low":
			lowCount++
		}
	}

	for _, c := range coverage {
		catCoverage[c.Category] = c.Coverage
		if c.Coverage < worstCov {
			worstCov = c.Coverage
			worstCategory = c.Category
		}
	}

	totalMissing := 0
	for _, mc := range missing {
		totalMissing += len(mc)
	}

	totalCoverage := 0.0
	if len(coverage) > 0 {
		for _, c := range coverage {
			totalCoverage += c.Coverage
		}
		totalCoverage = totalCoverage / float64(len(coverage))
	}

	summary := fmt.Sprintf("Security Design Review Summary — %d findings (%d critical, %d high, %d medium, %d low), %d missing controls across %d components, %d architectural weaknesses, %d prioritized remediations. Overall control coverage: %.1f%%. Highest risk area: %s (%.1f%% coverage).",
		len(findings), criticalCount, highCount, mediumCount, lowCount,
		totalMissing, len(missing), len(weaknesses), len(remediations),
		totalCoverage, worstCategory, worstCov)

	return summary
}

// ─────────────────────────────────────────────────────────────
// PHASE 12 — CONTROL COVERAGE DASHBOARD
// ─────────────────────────────────────────────────────────────

func buildCoverageDashboard(coverage []SDRICoverage) map[string]float64 {
	dashboard := make(map[string]float64)
	for _, c := range coverage {
		dashboard[c.Category] = c.Coverage
	}
	return dashboard
}

// ─────────────────────────────────────────────────────────────
// PHASE 13 — DOMAIN-SPECIFIC CONTROL EXPECTATIONS
// ─────────────────────────────────────────────────────────────

func getDomainSpecificControls(domain string) map[string][]string {
	domainPacks := map[string]map[string][]string{
		"healthcare": {
			"database": {"PHI Logging", "Break Glass Access", "Audit Trails", "Data Encryption at Rest", "Access Control"},
			"app":      {"Audit Trails", "Consent Management", "Data Classification", "Access Control", "Privacy Controls"},
			"api":      {"Audit Logging", "Data Encryption in Transit", "Break Glass Access", "Rate Limiting", "Authentication"},
			"identity": {"Break Glass Access", "Access Reviews", "MFA", "Session Management", "Audit Logging"},
			"storage":  {"Data Encryption at Rest", "PHI Logging", "Access Control", "Audit Trails", "Data Retention Policy"},
		},
		"fintech": {
			"database": {"Data Encryption at Rest", "Access Control", "Audit Logging", "Transaction Logging", "Dual Control"},
			"app":      {"Fraud Monitoring", "Settlement Integrity", "Dual Control", "Transaction Logging", "Audit Logging"},
			"api":      {"Fraud Monitoring", "Transaction Logging", "Rate Limiting", "Data Encryption in Transit", "Authentication"},
			"payment":  {"PCI DSS Compliance", "Fraud Monitoring", "Dual Control", "Transaction Logging", "Tokenization"},
			"ledger":   {"Settlement Integrity", "Dual Control", "Transaction Logging", "Audit Logging", "Access Control"},
		},
		"kubernetes": {
			"k8s":       {"K8s RBAC", "K8s Admission Controllers", "K8s Network Policies", "K8s Pod Security", "Audit Logging"},
			"cluster":   {"K8s RBAC", "K8s Admission Controllers", "K8s Network Policies", "K8s Pod Security", "Container Image Scanning"},
			"container": {"Container Image Scanning", "Container Runtime Security", "Vulnerability Scanning", "Access Control", "Audit Logging"},
			"api":       {"K8s RBAC", "Authentication", "Audit Logging", "Rate Limiting", "TLS Encryption"},
			"secret":    {"Secrets Vault", "Secrets Rotation", "Secrets Scanning", "Access Control", "Audit Logging"},
		},
		"saas": {
			"app":      {"MFA", "Conditional Access", "Session Management", "Audit Logging", "Data Encryption at Rest"},
			"identity": {"MFA", "Conditional Access", "Access Reviews", "Privileged Access Management", "Audit Logging"},
			"api":      {"TLS Encryption", "Authentication", "Rate Limiting", "Audit Logging", "Data Encryption in Transit"},
			"storage":  {"Data Encryption at Rest", "Access Control", "Audit Logging", "Data Retention Policy", "DLP Controls"},
		},
		"zero_trust": {
			"gateway":  {"MFA", "Conditional Access", "Continuous Verification", "TLS Encryption", "Session Management"},
			"identity": {"MFA", "Conditional Access", "Just-In-Time Access", "Privileged Access Management", "Continuous Verification"},
			"app":      {"MFA", "RBAC", "Just-In-Time Access", "Audit Logging", "Data Encryption in Transit"},
			"endpoint": {"Device Compliance", "Endpoint Detection & Response", "Continuous Verification", "Antivirus", "Patch Management"},
			"network":  {"Network Segmentation", "Micro-segmentation", "TLS Encryption", "Firewall Rules", "Intrusion Detection"},
		},
	}

	domain = strings.ToLower(domain)
	if pack, ok := domainPacks[domain]; ok {
		return pack
	}
	return nil
}

// ─────────────────────────────────────────────────────────────
// PHASE 14 — COMPLIANCE CONTROL ALIGNMENT
// ─────────────────────────────────────────────────────────────

type SDRIComplianceMapping struct {
	Framework string   `json:"framework"`
	Coverage  float64  `json:"coverage"`
	Controls  []string `json:"controls,omitempty"`
	Status    string   `json:"status"`
}

func computeComplianceAlignment(coverage []SDRICoverage, observed map[string][]string, frameworks []string) []SDRIComplianceMapping {
	frameworkControls := map[string]map[string]bool{
		"HIPAA": {
			"Data Encryption at Rest": true, "Data Encryption in Transit": true,
			"Audit Logging": true, "Access Control": true, "Automated Backup": true,
			"Incident Response Plan": true, "Data Classification": true, "Privacy Impact Assessment": true,
			"Data Retention Policy": true, "Data Subject Access": true,
		},
		"SOC2": {
			"Audit Logging": true, "Access Control": true, "Data Encryption at Rest": true,
			"Data Encryption in Transit": true, "Incident Response Plan": true,
			"Vulnerability Scanning": true, "Change Management": true, "Real-Time Alerting": true,
			"SIEM Integration": true, "Incident Response Team": true,
		},
		"ISO27001": {
			"Access Control": true, "Audit Logging": true, "Incident Response Plan": true,
			"Data Encryption at Rest": true, "Data Encryption in Transit": true,
			"Vulnerability Scanning": true, "Patch Management": true, "Change Management": true,
			"Third Party Due Diligence": true, "Business Continuity": true,
		},
		"PCI-DSS": {
			"Firewall Rules": true, "Data Encryption at Rest": true, "Data Encryption in Transit": true,
			"Vulnerability Scanning": true, "Access Control": true, "MFA": true,
			"Audit Logging": true, "Penetration Testing": true, "Incident Response Plan": true,
			"Key Rotation": true,
		},
	}

	// Collect all observed control names
	allObserved := make(map[string]bool)
	for _, ctrls := range observed {
		for _, c := range ctrls {
			allObserved[c] = true
		}
	}

	var mappings []SDRIComplianceMapping
	if len(frameworks) == 0 {
		frameworks = []string{"HIPAA", "SOC2", "ISO27001", "PCI-DSS"}
	}

	for _, fw := range frameworks {
		controls, ok := frameworkControls[fw]
		if !ok {
			continue
		}
		total := len(controls)
		met := 0
		var metControls []string
		for ctrl := range controls {
			if allObserved[ctrl] {
				met++
				metControls = append(metControls, ctrl)
			}
		}
		cov := 0.0
		if total > 0 {
			cov = float64(met) / float64(total) * 100
		}
		sort.Strings(metControls)
		mappings = append(mappings, SDRIComplianceMapping{
			Framework: fw,
			Coverage:  cov,
			Controls:  metControls,
			Status:    coverageLevel(cov),
		})
	}

	return mappings
}

// ─────────────────────────────────────────────────────────────
// PHASE 15 — SECURITY ARCHITECT REASONING
// ─────────────────────────────────────────────────────────────

func generateArchitectReasoning(f SDRIFinding, attackPaths []AttackPath, assumptions []Assumption) string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("**Why this matters:** %s\n\n", f.Description))
	b.WriteString(fmt.Sprintf("**Business impact:** %s\n\n", f.BusinessImpact))

	// Find related attack paths
	var relatedPaths []string
	for _, ap := range attackPaths {
		for _, comp := range f.AffectedComponents {
			for _, apComp := range ap.AffectedComponents {
				if strings.EqualFold(comp, apComp) {
					relatedPaths = append(relatedPaths, ap.Name)
					break
				}
			}
		}
	}
	if len(relatedPaths) > 0 {
		b.WriteString(fmt.Sprintf("**Exploitable via attack paths:** %s\n\n", strings.Join(relatedPaths, ", ")))
	} else {
		b.WriteString("**Exploitation scenario:** This control gap can be exploited by an attacker targeting the affected components. Without this control, the attack surface is enlarged.\n\n")
	}

	b.WriteString(fmt.Sprintf("**Recommended control:** %s\n", f.Recommendation))

	// Add evidence-based reasoning
	var supportingEvidence []string
	for _, comp := range f.AffectedComponents {
		for _, a := range assumptions {
			if strings.Contains(strings.ToLower(a.Component), strings.ToLower(comp)) ||
				strings.Contains(strings.ToLower(a.Description), strings.ToLower(comp)) {
				for _, ev := range a.EvidenceSources {
					supportingEvidence = append(supportingEvidence, ev)
				}
			}
		}
	}
	if len(supportingEvidence) > 0 {
		b.WriteString("\n**Supporting evidence:**\n")
		seen := make(map[string]bool)
		for _, ev := range supportingEvidence {
			if !seen[ev] {
				b.WriteString(fmt.Sprintf("- %s\n", ev))
				seen[ev] = true
			}
		}
	}

	return b.String()
}

// ─────────────────────────────────────────────────────────────
// SDRI ENGINE
// ─────────────────────────────────────────────────────────────

type SDRIResult struct {
	Controls                []SDRIControl               `json:"controls,omitempty"`
	ExpectedControls        map[string][]string         `json:"expected_controls,omitempty"`
	ObservedControls        map[string][]string         `json:"observed_controls,omitempty"`
	MissingControls         map[string][]string         `json:"missing_controls,omitempty"`
	CoverageByCategory      []SDRICoverage              `json:"coverage_by_category,omitempty"`
	DesignFindings          []SDRIFinding               `json:"design_findings,omitempty"`
	ArchitecturalWeaknesses []SDRIArchitecturalWeakness `json:"architectural_weaknesses,omitempty"`
	Remediations            []SDRIRemediation           `json:"remediations,omitempty"`
	ExecutiveSummary        string                      `json:"executive_summary"`
	CoverageDashboard       map[string]float64          `json:"coverage_dashboard,omitempty"`
	ComplianceAlignments    []SDRIComplianceMapping     `json:"compliance_alignments,omitempty"`
}

type SDRIEngine struct{}

func NewSDRIEngine() *SDRIEngine {
	return &SDRIEngine{}
}

func (e *SDRIEngine) Run(arch *ArchDescription, assumptions []Assumption, controls []ControlDetail,
	attackPaths []AttackPath, threats []Threat, domain string) *SDRIResult {

	result := &SDRIResult{}

	// Phase 3: Expected controls
	result.ExpectedControls = generateExpectedControls(arch)

	// Phase 4: Observed controls
	result.ObservedControls = extractObservedControls(arch, assumptions, controls)

	// Phase 5: Gap detection
	result.MissingControls = detectControlGaps(result.ExpectedControls, result.ObservedControls)

	// Phases 1-2: Build control list
	result.Controls = buildControlList(result.ExpectedControls, result.ObservedControls, result.MissingControls, assumptions, arch)

	// Phase 6: Coverage scoring
	result.CoverageByCategory = computeCoverageScores(result.ExpectedControls, result.ObservedControls)

	// Phase 12: Coverage dashboard
	result.CoverageDashboard = buildCoverageDashboard(result.CoverageByCategory)

	// Phase 9: Architectural weaknesses
	result.ArchitecturalWeaknesses = detectArchitecturalWeaknesses(arch, assumptions)

	// Phase 8: Design findings
	result.DesignFindings = generateDesignFindings(result.MissingControls, result.ArchitecturalWeaknesses, attackPaths, result.CoverageByCategory)

	// Phase 15: Architect reasoning
	for i := range result.DesignFindings {
		if result.DesignFindings[i].Reasoning == "" {
			result.DesignFindings[i].Reasoning = generateArchitectReasoning(result.DesignFindings[i], attackPaths, assumptions)
		}
	}

	// Phase 10: Remediation prioritization
	result.Remediations = prioritizeRemediations(result.DesignFindings, result.ArchitecturalWeaknesses, attackPaths, result.CoverageByCategory)

	// Phase 13: Domain-specific controls
	domainSpecific := getDomainSpecificControls(domain)
	if domainSpecific != nil {
		for comp, ctrls := range domainSpecific {
			result.ExpectedControls[comp] = append(result.ExpectedControls[comp], ctrls...)
		}
		// Re-run gap detection with domain-specific controls
		result.MissingControls = detectControlGaps(result.ExpectedControls, result.ObservedControls)
		result.CoverageByCategory = computeCoverageScores(result.ExpectedControls, result.ObservedControls)
		result.CoverageDashboard = buildCoverageDashboard(result.CoverageByCategory)
	}

	// Phase 14: Compliance alignment
	result.ComplianceAlignments = computeComplianceAlignment(result.CoverageByCategory, result.ObservedControls, arch.Compliance)

	// Phase 11: Executive summary
	result.ExecutiveSummary = buildExecutiveSummary(result.DesignFindings, result.MissingControls,
		result.ArchitecturalWeaknesses, result.CoverageByCategory, result.Remediations)

	return result
}

func generateExpectedControls(arch *ArchDescription) map[string][]string {
	expected := make(map[string][]string)
	ctrlMap := buildExpectedControlMap()
	for _, comp := range arch.Components {
		lower := strings.ToLower(comp.Label)
		bestMatch := ""
		bestLen := 0
		for key := range ctrlMap {
			if strings.Contains(lower, key) && len(key) > bestLen {
				bestMatch = key
				bestLen = len(key)
			}
		}
		if bestMatch != "" {
			expected[comp.Label] = ctrlMap[bestMatch]
		}
	}
	return expected
}

func buildControlList(expected, observed, missing map[string][]string,
	assumptions []Assumption, arch *ArchDescription) []SDRIControl {

	library := buildControlLibrary()
	strength := assessControlStrength(observed, assumptions)
	allNames := make(map[string]bool)
	var controls []SDRIControl
	ctrlID := 1

	// Merge all unique control names and their assignments
	controlAssignments := make(map[string]map[string]string) // controlName -> comp -> status
	for comp, ctrls := range expected {
		if controlAssignments[comp] == nil {
			controlAssignments[comp] = make(map[string]string)
		}
		for _, ctrl := range ctrls {
			status := "Missing"
			isObserved := false
			for _, oCtrl := range observed[comp] {
				if strings.EqualFold(oCtrl, ctrl) {
					isObserved = true
					break
				}
			}
			if isObserved {
				status = "Present"
			}
			controlAssignments[comp][ctrl] = status
			allNames[ctrl] = true
		}
	}

	for comp, ctrls := range observed {
		if controlAssignments[comp] == nil {
			controlAssignments[comp] = make(map[string]string)
		}
		for _, ctrl := range ctrls {
			if _, exists := controlAssignments[comp][ctrl]; !exists {
				controlAssignments[comp][ctrl] = "Present"
				allNames[ctrl] = true
			}
		}
	}

	ctrlNames := make([]string, 0, len(allNames))
	for name := range allNames {
		ctrlNames = append(ctrlNames, name)
	}
	sort.Strings(ctrlNames)

	for _, name := range ctrlNames {
		entry, inLib := library[name]
		if !inLib {
			for _, e := range library {
				if e.Name == name {
					entry = e
					inLib = true
					break
				}
			}
		}
		if !inLib {
			entry = SDRIControlLibraryEntry{
				Name:        name,
				Category:    "General",
				Description: name,
				ControlType: SDRIControlPreventive,
			}
		}

		cat := entry.Category
		desc := entry.Description
		ctrlType := entry.ControlType
		prev := entry.Preventive
		det := entry.Detective
		corr := entry.Corrective

		foundWeak := false
		foundPartial := false
		foundStrong := false
		coverage := "Missing"
		status := "Missing"

		for comp := range controlAssignments {
			s := controlAssignments[comp][name]
			if s == "Present" {
				coverage = "Present"
				status = "Present"
				key := comp + ":" + name
				if st, ok := strength[key]; ok {
					switch st {
					case SDRIStrengthStrong:
						foundStrong = true
					case SDRIStrengthPartial:
						foundPartial = true
					case SDRIStrengthWeak:
						foundWeak = true
					}
				}
			}
		}

		var effStrength SDRIStrength
		if foundStrong {
			effStrength = SDRIStrengthStrong
		} else if foundPartial {
			effStrength = SDRIStrengthPartial
		} else if foundWeak {
			effStrength = SDRIStrengthWeak
		} else if coverage == "Present" {
			effStrength = SDRIStrengthPartial
		} else {
			effStrength = SDRIStrengthWeak
		}

		if coverage == "Missing" {
			status = "Missing"
			effStrength = SDRIStrengthWeak
		}

		var evidence []string
		for _, a := range assumptions {
			for _, ev := range a.EvidenceSources {
				if strings.Contains(strings.ToLower(ev), strings.ToLower(name)) {
					evidence = append(evidence, ev)
				}
			}
		}

		controls = append(controls, SDRIControl{
			ID:          fmt.Sprintf("CTRL-%03d", ctrlID),
			Name:        name,
			Category:    cat,
			Description: desc,
			ControlType: ctrlType,
			Preventive:  prev,
			Detective:   det,
			Corrective:  corr,
			Strength:    effStrength,
			Evidence:    evidence,
			Coverage:    coverage,
			Status:      status,
		})
		ctrlID++
	}

	return controls
}
