package trust

import (
	"fmt"
	"strings"
)

// DiscoveryEngine discovers dependencies between assumptions deterministically.
type DiscoveryEngine struct {
	domain     string
	components []string
	compIndex  map[string][]string // component name -> derived keywords
}

// NewDiscoveryEngine creates a new dependency discovery engine.
func NewDiscoveryEngine(domain string, components []string) *DiscoveryEngine {
	return &DiscoveryEngine{
		domain:     domain,
		components: components,
		compIndex:  buildComponentIndex(components),
	}
}

// buildComponentIndex maps component names to derived keywords for matching.
func buildComponentIndex(components []string) map[string][]string {
	idx := make(map[string][]string)
	for _, c := range components {
		lower := strings.ToLower(c)
		kw := []string{lower}
		// Split PascalCase into words: SharedAdminAccount -> shared, admin, account
		// Skip single-letter or purely-numeric fragments
		var words []string
		current := make([]byte, 0)
		for i := 0; i < len(lower); i++ {
			ch := lower[i]
			if ch >= 'a' && ch <= 'z' {
				current = append(current, ch)
			} else if ch >= 'A' && ch <= 'Z' {
				if len(current) > 0 {
					words = append(words, string(current))
				}
				current = []byte{ch + 32} // lowercase
			} else {
				// Digit or special char: end current word without starting a new one
				if len(current) > 0 {
					words = append(words, string(current))
					current = nil
				}
			}
		}
		if len(current) > 0 {
			words = append(words, string(current))
		}
		// Only add meaningful words (length >= 3, not purely numeric)
		for _, w := range words {
			if len(w) >= 3 {
				kw = append(kw, w)
			}
		}
		idx[lower] = kw
	}
	return idx
}

// getComponentKeywords returns derived keywords for a component name, or nil.
func (e *DiscoveryEngine) getComponentKeywords(component string) []string {
	if e.compIndex == nil {
		return nil
	}
	return e.compIndex[strings.ToLower(component)]
}

// AssumptionInput is the input for dependency discovery.
type AssumptionInput struct {
	ID         string
	Text       string
	Component  string
	Category   string
	Risk       string
	Confidence float64
	Keywords   []string
	Source     string
}

// DiscoverDependencies builds the full dependency graph from assumptions.
func (e *DiscoveryEngine) DiscoverDependencies(assumptions []AssumptionInput) *DependencyGraph {
	graph := NewDependencyGraph()

	// Add all nodes
	for _, a := range assumptions {
		node := &AssumptionNode{
			ID:         a.ID,
			Text:       a.Text,
			Type:       a.Category,
			Risk:       a.Risk,
			Confidence: a.Confidence,
			Source:     a.Source,
			Component:  a.Component,
			Category:   a.Category,
		}
		graph.AddNode(node)
	}

	// Discover dependencies
	for _, a := range assumptions {
		edges := e.findDependencies(a, assumptions)
		for _, edge := range edges {
			// Avoid creating bidirectional edges (cycles): skip if reverse edge exists
			if !graph.HasEdge(edge.TargetAssumption, edge.SourceAssumption) {
				graph.AddEdge(edge)
			}
		}
	}

	// Compute metrics
	graph.ComputeCentrality()
	graph.ComputeRadii()

	return graph
}

// findDependencies finds all dependencies for a single assumption.
func (e *DiscoveryEngine) findDependencies(a AssumptionInput, all []AssumptionInput) []DependencyEdge {
	var edges []DependencyEdge

	text := strings.ToLower(a.Text)
	keywords := make([]string, len(a.Keywords))
	for i, k := range a.Keywords {
		keywords[i] = strings.ToLower(k)
	}

	for _, other := range all {
		if other.ID == a.ID {
			continue
		}

		otherText := strings.ToLower(other.Text)
		otherKeywords := make([]string, len(other.Keywords))
		for i, k := range other.Keywords {
			otherKeywords[i] = strings.ToLower(k)
		}

		// Check if 'a' depends on 'other'
		if depType, strength, reason := e.detectDependency(a, other, text, otherText, keywords, otherKeywords); depType != "" {
			edges = append(edges, DependencyEdge{
				SourceAssumption: a.ID,
				TargetAssumption: other.ID,
				DependencyType:   depType,
				Strength:         strength,
				Confidence:       a.Confidence * other.Confidence,
				Reason:           reason,
				IsExplicit:       false,
			})
		}
	}

	return edges
}

// detectDependency checks if assumption 'a' depends on assumption 'other'.
func (e *DiscoveryEngine) detectDependency(a, other AssumptionInput, aText, otherText string, aKeywords, otherKeywords []string) (DependencyType, float64, string) {
	// Check if they share the same component - dependencies within same component
	if a.Component != "" && a.Component == other.Component {
		// Check for specific dependency patterns
		if depType, reason := e.checkComponentDependency(a, other, aText, otherText); depType != "" {
			return depType, 0.7, reason
		}
	}

	// Check keyword-based dependencies
	if depType, reason := e.checkKeywordDependency(a, other, aKeywords, otherKeywords); depType != "" {
		return depType, 0.6, reason
	}

	// Check failure-mode-specific dependencies (shared admin, VPN, flat network, etc.)
	if depType, reason := e.checkFailureModeDependency(a, other, aText, otherText, aKeywords, otherKeywords); depType != "" {
		return depType, 0.65, reason
	}

	// Check domain-specific dependencies
	if depType, reason := e.checkDomainDependency(a, other, aText, otherText); depType != "" {
		return depType, 0.8, reason
	}

	// Check category-based dependencies
	if depType, reason := e.checkCategoryDependency(a, other, aText, otherText); depType != "" {
		return depType, 0.5, reason
	}

	return "", 0, ""
}

// checkComponentDependency checks for dependencies within the same component.
func (e *DiscoveryEngine) checkComponentDependency(a, other AssumptionInput, aText, otherText string) (DependencyType, string) {
	// Identity dependencies
	if e.containsAny(aText, []string{"access", "auth", "login", "sso", "mfa", "role", "permission", "rbac"}) &&
		e.containsAny(otherText, []string{"identity", "idp", "provider", "authentication", "mfa", "sso", "federation"}) {
		return DepIdentity, fmt.Sprintf("%s access depends on %s identity", a.Component, other.Component)
	}

	// Cryptographic dependencies
	if e.containsAny(aText, []string{"encrypt", "tls", "ssl", "certificate", "data protection"}) &&
		e.containsAny(otherText, []string{"key", "kms", "rotation", "certificate", "vault", "secrets"}) {
		return DepCryptographic, fmt.Sprintf("%s encryption depends on %s key management", a.Component, other.Component)
	}

	// Monitoring dependencies
	if e.containsAny(aText, []string{"audit", "log", "monitor", "alert", "siem"}) &&
		e.containsAny(otherText, []string{"log", "monitor", "siem", "observability"}) {
		return DepMonitoring, fmt.Sprintf("%s audit depends on %s logging", a.Component, other.Component)
	}

	// Authorization dependencies
	if e.containsAny(aText, []string{"access", "rbac", "permission", "privilege", "role"}) &&
		e.containsAny(otherText, []string{"rbac", "authorization", "permission", "role", "policy"}) {
		return DepAuthorization, fmt.Sprintf("%s access control depends on %s authorization", a.Component, other.Component)
	}

	return "", ""
}

// checkKeywordDependency checks for keyword-based dependencies.
func (e *DiscoveryEngine) checkKeywordDependency(a, other AssumptionInput, aKeywords, otherKeywords []string) (DependencyType, string) {
	// Check for MFA dependencies
	if e.containsAnySlice(aKeywords, []string{"mfa", "multi-factor", "2fa", "totp"}) &&
		e.containsAnySlice(otherKeywords, []string{"identity", "authentication", "auth", "idp"}) {
		return DepIdentity, "MFA depends on identity provider"
	}

	// Check for backup dependencies
	if e.containsAnySlice(aKeywords, []string{"backup", "restore", "recovery"}) &&
		e.containsAnySlice(otherKeywords, []string{"backup", "restore", "storage", "replication"}) {
		return DepOperational, "Backup depends on storage infrastructure"
	}

	// Check for network dependencies
	if e.containsAnySlice(aKeywords, []string{"network", "segment", "firewall", "vpn", "vpc"}) &&
		e.containsAnySlice(otherKeywords, []string{"network", "infrastructure", "dns", "load balancer"}) {
		return DepInfrastructure, "Network security depends on network infrastructure"
	}

	// Check for third-party dependencies
	if e.containsAnySlice(aKeywords, []string{"auth0", "aws", "azure", "cloudflare", "stripe", "saas", "third-party"}) &&
		e.containsAnySlice(otherKeywords, []string{"auth0", "aws", "azure", "cloudflare", "stripe", "provider", "vendor"}) {
		return DepThirdParty, "Assumption depends on third-party service"
	}

	// Check for encryption dependencies
	if e.containsAnySlice(aKeywords, []string{"encrypt", "tls", "ssl", "certificate"}) &&
		e.containsAnySlice(otherKeywords, []string{"key", "kms", "rotation", "vault", "certificate"}) {
		return DepCryptographic, "Encryption depends on key management"
	}

	return "", ""
}

// checkDomainDependency checks for domain-specific dependencies.
func (e *DiscoveryEngine) checkDomainDependency(a, other AssumptionInput, aText, otherText string) (DependencyType, string) {
	if e.domain == "" {
		return "", ""
	}

	switch e.domain {
	case "healthcare", "hipaa":
		return e.checkHealthcareDependency(a, other, aText, otherText)
	case "fintech", "pci", "finance":
		return e.checkFintechDependency(a, other, aText, otherText)
	case "cloud", "aws", "azure":
		return e.checkCloudDependency(a, other, aText, otherText)
	case "kubernetes", "k8s":
		return e.checkKubernetesDependency(a, other, aText, otherText)
	case "vpn", "network":
		return e.checkVPNDependency(a, other, aText, otherText)
	case "saas":
		return e.checkSaaSDependency(a, other, aText, otherText)
	}

	return "", ""
}

// checkHealthcareDependency finds healthcare-specific dependencies.
func (e *DiscoveryEngine) checkHealthcareDependency(a, other AssumptionInput, aText, otherText string) (DependencyType, string) {
	// PHI access depends on identity - but only for access-related assumptions
	if (e.containsAny(aText, []string{"phi access", "patient access", "medical access", "clinical access"}) ||
		(e.containsAny(aText, []string{"phi", "patient", "healthcare", "medical", "clinical"}) &&
			e.containsAny(aText, []string{"access", "control", "restrict"}))) &&
		e.containsAny(otherText, []string{"identity", "authentication", "auth", "mfa", "rbac"}) {
		return DepIdentity, "PHI access depends on identity controls"
	}

	// PHI access depends on audit
	if e.containsAny(aText, []string{"phi access", "patient access", "medical access"}) &&
		e.containsAny(otherText, []string{"audit", "log", "hipaa", "compliance"}) {
		return DepMonitoring, "PHI access depends on audit logging"
	}

	// Break-glass depends on logging
	if e.containsAny(aText, []string{"break-glass", "emergency", "override"}) &&
		e.containsAny(otherText, []string{"audit", "log", "monitor"}) {
		return DepMonitoring, "Break-glass depends on audit logging"
	}

	// Encryption for PHI depends on key management
	if e.containsAny(aText, []string{"encrypt", "encryption"}) &&
		e.containsAny(aText, []string{"phi", "patient"}) &&
		e.containsAny(otherText, []string{"key", "kms", "rotation", "vault"}) {
		return DepCryptographic, "PHI encryption depends on key management"
	}

	return "", ""
}

// checkFintechDependency finds fintech-specific dependencies.
func (e *DiscoveryEngine) checkFintechDependency(a, other AssumptionInput, aText, otherText string) (DependencyType, string) {
	// Payment depends on encryption
	if e.containsAny(aText, []string{"payment", "transaction", "settlement", "pci"}) &&
		e.containsAny(otherText, []string{"encrypt", "tls", "certificate", "kms"}) {
		return DepCryptographic, "Payment processing depends on encryption"
	}

	// Fraud depends on monitoring
	if e.containsAny(aText, []string{"fraud", "aml", "kyc", "monitoring"}) &&
		e.containsAny(otherText, []string{"log", "monitor", "siem", "alert"}) {
		return DepMonitoring, "Fraud detection depends on monitoring"
	}

	// Key custody depends on key management
	if e.containsAny(aText, []string{"key custody", "key management", "hsm"}) &&
		e.containsAny(otherText, []string{"key", "kms", "rotation", "vault"}) {
		return DepCryptographic, "Key custody depends on key management"
	}

	return "", ""
}

// checkCloudDependency finds cloud-specific dependencies.
func (e *DiscoveryEngine) checkCloudDependency(a, other AssumptionInput, aText, otherText string) (DependencyType, string) {
	// IAM depends on identity
	if e.containsAny(aText, []string{"iam", "role", "policy", "access"}) &&
		e.containsAny(otherText, []string{"identity", "authentication", "federation", "sso"}) {
		return DepIdentity, "IAM depends on identity provider"
	}

	// KMS depends on IAM
	if e.containsAny(aText, []string{"kms", "encrypt", "key"}) &&
		e.containsAny(otherText, []string{"iam", "role", "policy", "access"}) {
		return DepAuthorization, "KMS depends on IAM authorization"
	}

	// Secrets depend on KMS
	if e.containsAny(aText, []string{"secret", "vault", "credential"}) &&
		e.containsAny(otherText, []string{"kms", "encrypt", "key"}) {
		return DepCryptographic, "Secrets management depends on KMS"
	}

	return "", ""
}

// checkKubernetesDependency finds Kubernetes-specific dependencies.
func (e *DiscoveryEngine) checkKubernetesDependency(a, other AssumptionInput, aText, otherText string) (DependencyType, string) {
	// RBAC depends on authentication
	if e.containsAny(aText, []string{"rbac", "role", "serviceaccount"}) &&
		e.containsAny(otherText, []string{"auth", "authentication", "identity", "certificate"}) {
		return DepIdentity, "Kubernetes RBAC depends on authentication"
	}

	// Admission depends on policy
	if e.containsAny(aText, []string{"admission", "webhook", "policy"}) &&
		e.containsAny(otherText, []string{"rbac", "authorization", "policy"}) {
		return DepAuthorization, "Admission control depends on authorization"
	}

	// Secrets depend on encryption
	if e.containsAny(aText, []string{"secret", "etcd", "configmap"}) &&
		e.containsAny(otherText, []string{"encrypt", "kms", "etcd"}) {
		return DepCryptographic, "Kubernetes secrets depend on encryption"
	}

	return "", ""
}

// checkVPNDependency finds VPN-specific dependencies.
func (e *DiscoveryEngine) checkVPNDependency(a, other AssumptionInput, aText, otherText string) (DependencyType, string) {
	// VPN depends on certificates
	if e.containsAny(aText, []string{"vpn", "tunnel", "remote"}) &&
		e.containsAny(otherText, []string{"certificate", "tls", "ca", "pki"}) {
		return DepCryptographic, "VPN depends on certificate infrastructure"
	}

	// VPN depends on identity
	if e.containsAny(aText, []string{"vpn", "access"}) &&
		e.containsAny(otherText, []string{"identity", "auth", "mfa", "radius", "ldap"}) {
		return DepIdentity, "VPN access depends on identity"
	}

	return "", ""
}

// checkSaaSDependency finds SaaS-specific dependencies.
func (e *DiscoveryEngine) checkSaaSDependency(a, other AssumptionInput, aText, otherText string) (DependencyType, string) {
	// Tenant isolation depends on identity
	if e.containsAny(aText, []string{"tenant", "isolation", "multi-tenant"}) &&
		e.containsAny(otherText, []string{"identity", "auth", "rbac"}) {
		return DepIdentity, "Tenant isolation depends on identity controls"
	}

	// API security depends on rate limiting
	if e.containsAny(aText, []string{"api", "gateway"}) &&
		e.containsAny(otherText, []string{"rate", "limit", "throttle"}) {
		return DepInfrastructure, "API security depends on rate limiting"
	}

	return "", ""
}

// checkFailureModeDependency detects specific security failure mode dependencies:
// shared admin, third-party VPN, monitoring default credentials, flat network.
func (e *DiscoveryEngine) checkFailureModeDependency(a, other AssumptionInput, aText, otherText string, aKeywords, otherKeywords []string) (DependencyType, string) {
	// Shared admin -> dependency on authorization/identity
	// Matches when one assumption mentions shared admin and another mentions auth/identity
	if e.containsAny(aText, []string{"shared admin", "shared account", "shared credential"}) &&
		e.containsAny(otherText, []string{"identity", "auth", "mfa", "rbac", "authorization"}) {
		return DepIdentity, "Shared admin account depends on identity controls"
	}

	// Third-party VPN -> dependency on network infrastructure
	// Matches when one assumption mentions third-party + VPN and another mentions network
	if e.containsAny(aText, []string{"third-party vpn", "third party vpn", "vendor vpn", "partner vpn"}) &&
		e.containsAny(otherText, []string{"network", "segment", "firewall", "vpn", "gateway"}) {
		return DepThirdParty, "Third-party VPN depends on network infrastructure"
	}

	// Also detect third-party access without explicit VPN keyword
	if e.containsAny(aText, []string{"third-party", "third party", "vendor full access", "partner access", "external access"}) &&
		e.containsAny(otherText, []string{"network", "segment", "firewall", "vpn", "vpc"}) {
		return DepThirdParty, "Third-party access depends on network controls"
	}

	// Monitoring default credentials -> dependency on operational security
	// Matches when one assumption mentions monitoring + default creds and another mentions security
	if e.containsAny(aText, []string{"default credential", "default password", "default admin", "hardcoded password", "hardcoded credential"}) &&
		e.containsAny(otherText, []string{"monitor", "log", "alert", "siem", "scan", "patch"}) {
		return DepMonitoring, "Default credentials depend on monitoring detection"
	}

	// Flat network -> dependency on segmentation controls
	// Matches when one assumption mentions flat network and another mentions segmentation
	if e.containsAny(aText, []string{"flat network", "no segmentation", "flat network topology", "no network segmentation"}) &&
		e.containsAny(otherText, []string{"segment", "firewall", "isolat", "vlan", "acl", "network control"}) {
		return DepInfrastructure, "Flat network depends on segmentation controls"
	}

	// Component-name-based matching for shared admin accounts
	if a.Component != "" && e.containsStringField(a.Component, "shared.admin") ||
		other.Component != "" && e.containsStringField(other.Component, "shared.admin") {
		return DepAuthorization, "Shared admin component creates authorization dependency"
	}

	return "", ""
}

// containsStringField checks if a field value matches a pattern.
// Pattern uses dots as word separators, matches partial PascalCase substrings.
func (e *DiscoveryEngine) containsStringField(field, pattern string) bool {
	lower := strings.ToLower(field)
	parts := strings.Split(pattern, ".")
	for _, p := range parts {
		if !strings.Contains(lower, p) {
			return false
		}
	}
	return true
}

// checkCategoryDependency finds dependencies based on category relationships.
func (e *DiscoveryEngine) checkCategoryDependency(a, other AssumptionInput, aText, otherText string) (DependencyType, string) {
	// Identity assumptions support authorization assumptions
	if a.Category == "authorization" || a.Category == "access" {
		if other.Category == "identity" || other.Category == "authentication" {
			return DepIdentity, "Authorization depends on identity"
		}
	}

	// Encryption assumptions depend on key management
	if a.Category == "encryption" || a.Category == "cryptographic" {
		if other.Category == "key management" || other.Category == "secrets" {
			return DepCryptographic, "Encryption depends on key management"
		}
	}

	// Monitoring assumptions depend on logging infrastructure
	if a.Category == "monitoring" || a.Category == "audit" {
		if other.Category == "logging" || other.Category == "infrastructure" {
			return DepMonitoring, "Monitoring depends on logging infrastructure"
		}
	}

	return "", ""
}

// containsAny checks if text contains any of the keywords.
func (e *DiscoveryEngine) containsAny(text string, keywords []string) bool {
	for _, kw := range keywords {
		if strings.Contains(text, kw) {
			return true
		}
	}
	return false
}

// containsAnySlice checks if any keyword in the slice is in the list.
func (e *DiscoveryEngine) containsAnySlice(keywords, list []string) bool {
	for _, kw := range keywords {
		lowerKw := strings.ToLower(kw)
		for _, item := range list {
			if strings.Contains(lowerKw, item) {
				return true
			}
		}
	}
	return false
}
