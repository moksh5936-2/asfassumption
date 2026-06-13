package narrative

import (
	"fmt"
	"strings"
)

// NarrativeEngine generates architect-style narratives for ASF findings.
type NarrativeEngine struct {
	// Domain context for generating domain-specific narratives
	domain string

	// Architecture components for dependency mapping
	components []string

	// Relationships for dependency mapping
	relationships []string
}

// NewNarrativeEngine creates a new narrative engine.
func NewNarrativeEngine(domain string, components, relationships []string) *NarrativeEngine {
	return &NarrativeEngine{
		domain:        domain,
		components:    components,
		relationships: relationships,
	}
}

// GenerateNarrative creates a full narrative output for an analysis result.
func (e *NarrativeEngine) GenerateNarrative(
	archName string,
	assumptions []Assumption,
	controls []ControlDetail,
	trustBoundaries []TrustBoundary,
	contradictions []Contradiction,
	domain string,
	strideDist map[string]int,
	riskDist map[string]int,
) *NarrativeOutput {

	output := &NarrativeOutput{
		AssumptionNarratives: make([]AssumptionNarrative, 0, len(assumptions)),
	}

	// Build dependency map
	depMap := e.buildDependencyMap(assumptions)

	// Generate per-assumption narratives
	for _, a := range assumptions {
		narrative := e.generateAssumptionNarrative(a, depMap)
		output.AssumptionNarratives = append(output.AssumptionNarratives, narrative)
	}

	// Generate architecture overview
	output.ArchitectureOverview = e.generateArchitectureOverview(
		archName, assumptions, controls, trustBoundaries, domain,
	)

	// Generate full architect narrative text
	output.ArchitectNarrative = e.generateFullArchitectNarrative(
		output.ArchitectureOverview, output.AssumptionNarratives,
	)

	return output
}

// Assumption is the interface for the narrative engine.
// It mirrors the engine.go Assumption struct but uses a local copy for package independence.
type Assumption struct {
	ID                  string
	Description         string
	Component           string
	Category            string
	Risk                string
	STRIDECategories    []string
	Likelihood          int
	Impact              int
	Confidence          float64
	Keywords            []string
	SourceComponents    []string
	SourceRelationships []string
	Rationale           string
	EvidenceSources     []string
}

// ControlDetail mirrors the engine.go type.
type ControlDetail struct {
	Name                 string
	Category             string
	Description          string
	Rationale            string
	MitigatedAssumptions []string
	STRIDECategories     []string
}

// TrustBoundary mirrors the engine.go type.
type TrustBoundary struct {
	Type        string
	Components  []string
	RiskLevel   string
	Description string
}

// Contradiction mirrors the engine.go type.
type Contradiction struct {
	ID                  string
	Severity            string
	Description         string
	Explanation         string
	AffectedAssumptions []string
}

// generateAssumptionNarrative creates the architect-style narrative for one assumption.
func (e *NarrativeEngine) generateAssumptionNarrative(a Assumption, depMap map[string][]string) AssumptionNarrative {
	n := AssumptionNarrative{
		AssumptionID:     a.ID,
		AssumptionText:   a.Description,
		RiskLevel:        a.Risk,
		STRIDECategories: a.STRIDECategories,
		Confidence:       a.Confidence,
		DependsOn:        e.extractDependsOn(a),
		DownstreamImpact: depMap[a.ID],
	}

	// Context: What is the architectural context for this assumption?
	n.Context = e.generateContext(a)

	// Why ASF Identified It: What in the architecture triggered this finding?
	n.WhyASFIdentifiedIt = e.generateWhyIdentified(a)

	// Architectural Importance: Why does this matter to the architecture?
	n.ArchitecturalImportance = e.generateArchitecturalImportance(a, depMap[a.ID])

	// Failure Consequence: What happens if this assumption fails?
	n.FailureConsequence = e.generateFailureConsequence(a, depMap[a.ID])

	// Security Recommendation: What should be done?
	n.SecurityRecommendation = e.generateRecommendation(a)

	// Apply style enforcement
	n.Context = e.enforceStyle(n.Context)
	n.WhyASFIdentifiedIt = e.enforceStyle(n.WhyASFIdentifiedIt)
	n.ArchitecturalImportance = e.enforceStyle(n.ArchitecturalImportance)
	n.FailureConsequence = e.enforceStyle(n.FailureConsequence)
	n.SecurityRecommendation = e.enforceStyle(n.SecurityRecommendation)

	return n
}

// generateContext creates the architectural context section.
func (e *NarrativeEngine) generateContext(a Assumption) string {
	var parts []string

	// Identify the component
	if a.Component != "" {
		parts = append(parts, fmt.Sprintf("The %s component is", a.Component))
	} else {
		parts = append(parts, "The architecture")
	}

	// What role does it play?
	role := e.inferRole(a)
	if role != "" {
		parts = append(parts, role)
	}

	// What controls are mentioned?
	if len(a.Keywords) > 0 {
		kw := e.filterSecurityKeywords(a.Keywords)
		if len(kw) > 0 {
			parts = append(parts, fmt.Sprintf("with %s", strings.Join(kw, ", ")))
		}
	}

	// Build context sentence
	if len(parts) >= 2 {
		return strings.Join(parts, " ") + "."
	}

	// Fallback
	return fmt.Sprintf("The architecture describes %s.", a.Description)
}

// generateWhyIdentified explains what triggered this finding.
func (e *NarrativeEngine) generateWhyIdentified(a Assumption) string {
	// Use the rationale if available
	if a.Rationale != "" {
		// Clean up the rationale
		r := strings.TrimSpace(a.Rationale)
		// Remove trailing punctuation if we'll add our own
		r = strings.TrimRight(r, ".")
		return r + "."
	}

	// Use evidence sources
	if len(a.EvidenceSources) > 0 {
		return fmt.Sprintf("ASF detected references to %s in the architecture documentation.",
			strings.Join(a.EvidenceSources, ", "))
	}

	// Use keywords
	if len(a.Keywords) > 0 {
		return fmt.Sprintf("The architecture references %s, which implies this assumption.",
			strings.Join(a.Keywords, ", "))
	}

	// Fallback
	return fmt.Sprintf("ASF identified this assumption based on the presence of %s in the architecture.", a.Category)
}

// generateArchitecturalImportance explains why this matters.
func (e *NarrativeEngine) generateArchitecturalImportance(a Assumption, downstream []string) string {
	var parts []string

	// Risk level framing
	switch a.Risk {
	case "Critical":
		parts = append(parts, "This is a critical architectural dependency.")
	case "High":
		parts = append(parts, "This is a significant architectural dependency.")
	case "Medium":
		parts = append(parts, "This is an architectural dependency with moderate impact.")
	default:
		parts = append(parts, "This is a standard architectural dependency.")
	}

	// Downstream impact
	if len(downstream) > 0 {
		if len(downstream) == 1 {
			parts = append(parts, fmt.Sprintf("If this assumption fails, %s is affected.", downstream[0]))
		} else if len(downstream) <= 3 {
			parts = append(parts, fmt.Sprintf("If this assumption fails, %s are affected.",
				strings.Join(downstream, ", ")))
		} else {
			parts = append(parts, fmt.Sprintf("If this assumption fails, %d downstream components are affected, including %s.",
				len(downstream), strings.Join(downstream[:3], ", ")))
		}
	}

	// STRIDE framing
	if len(a.STRIDECategories) > 0 {
		parts = append(parts, fmt.Sprintf("The STRIDE categories %s indicate this assumption protects against multiple threat types.",
			strings.Join(a.STRIDECategories, ", ")))
	}

	return strings.Join(parts, " ")
}

// generateFailureConsequence describes what happens if the assumption fails.
func (e *NarrativeEngine) generateFailureConsequence(a Assumption, downstream []string) string {
	var parts []string

	// Immediate consequence
	consequence := e.inferConsequence(a)
	parts = append(parts, consequence)

	// Cascading impact
	if len(downstream) > 0 {
		if len(downstream) == 1 {
			parts = append(parts, fmt.Sprintf("This cascades to %s.", downstream[0]))
		} else {
			parts = append(parts, fmt.Sprintf("This cascades to %d components: %s.",
				len(downstream), strings.Join(downstream, ", ")))
		}
	}

	// Business impact framing
	if a.Risk == "Critical" || a.Risk == "High" {
		parts = append(parts, "This represents a material risk to the architecture.")
	}

	return strings.Join(parts, " ")
}

// generateRecommendation creates a specific security recommendation.
func (e *NarrativeEngine) generateRecommendation(a Assumption) string {
	// Try to infer specific controls from the assumption
	controls := e.inferControls(a)
	if len(controls) > 0 {
		return fmt.Sprintf("Implement %s to satisfy this assumption.", strings.Join(controls, ", "))
	}

	// Fallback based on category
	switch a.Category {
	case "identity":
		return "Verify identity controls are documented, tested, and monitored."
	case "access":
		return "Verify access controls enforce least privilege and are regularly audited."
	case "network":
		return "Verify network segmentation and traffic controls are enforced."
	case "configuration":
		return "Verify configuration is hardened, version-controlled, and change-managed."
	case "process":
		return "Verify processes are documented, reviewed, and followed."
	default:
		return fmt.Sprintf("Verify %s controls are in place and tested.", a.Category)
	}
}

// buildDependencyMap builds a map of assumption ID -> downstream components.
func (e *NarrativeEngine) buildDependencyMap(assumptions []Assumption) map[string][]string {
	depMap := make(map[string][]string)

	for _, a := range assumptions {
		var downstream []string
		for _, other := range assumptions {
			if other.ID == a.ID {
				continue
			}
			// If other assumption references the same component, it depends on this one
			if a.Component != "" && other.Component == a.Component {
				continue // Same component, not downstream
			}
			// Check if other assumption shares keywords or relationships
			if e.sharesRelationship(a, other) {
				downstream = append(downstream, other.Component)
			}
		}
		// Deduplicate
		seen := make(map[string]bool)
		var unique []string
		for _, d := range downstream {
			if d != "" && !seen[d] {
				seen[d] = true
				unique = append(unique, d)
			}
		}
		depMap[a.ID] = unique
	}

	return depMap
}

// sharesRelationship checks if two assumptions share a relationship.
func (e *NarrativeEngine) sharesRelationship(a, b Assumption) bool {
	for _, r1 := range a.SourceRelationships {
		for _, r2 := range b.SourceRelationships {
			if r1 == r2 && r1 != "" {
				return true
			}
		}
	}
	return false
}

// extractDependsOn extracts what an assumption depends on.
func (e *NarrativeEngine) extractDependsOn(a Assumption) []string {
	var deps []string
	if a.Component != "" {
		deps = append(deps, a.Component)
	}
	for _, sc := range a.SourceComponents {
		if sc != "" && sc != a.Component {
			deps = append(deps, sc)
		}
	}
	return deps
}

// inferRole infers the architectural role from the assumption.
func (e *NarrativeEngine) inferRole(a Assumption) string {
	desc := strings.ToLower(a.Description)
	component := strings.ToLower(a.Component)

	if strings.Contains(desc, "authentication") || strings.Contains(component, "auth") {
		return "responsible for authentication"
	}
	if strings.Contains(desc, "authorization") || strings.Contains(desc, "access") {
		return "responsible for access control"
	}
	if strings.Contains(desc, "encryption") || strings.Contains(desc, "encrypt") {
		return "responsible for data protection"
	}
	if strings.Contains(desc, "logging") || strings.Contains(desc, "audit") {
		return "responsible for observability"
	}
	if strings.Contains(desc, "network") || strings.Contains(desc, "firewall") {
		return "responsible for network security"
	}
	if strings.Contains(desc, "database") || strings.Contains(desc, "storage") {
		return "responsible for data storage"
	}
	if strings.Contains(desc, "api") || strings.Contains(desc, "gateway") {
		return "responsible for API management"
	}
	return ""
}

// inferConsequence infers the failure consequence from the assumption.
func (e *NarrativeEngine) inferConsequence(a Assumption) string {
	desc := strings.ToLower(a.Description)
	category := strings.ToLower(a.Category)

	if strings.Contains(desc, "mfa") || strings.Contains(desc, "multi-factor") {
		return "Without MFA, administrative accounts are susceptible to credential compromise."
	}
	if strings.Contains(desc, "encrypt") {
		return "Without encryption, data is exposed in transit or at rest."
	}
	if strings.Contains(desc, "access") || strings.Contains(desc, "privilege") {
		return "Without access restrictions, unauthorized parties can interact with protected resources."
	}
	if strings.Contains(desc, "logging") || strings.Contains(desc, "audit") {
		return "Without logging, security events go undetected."
	}
	if strings.Contains(desc, "backup") || strings.Contains(desc, "recovery") {
		return "Without backups, data loss is unrecoverable."
	}
	if strings.Contains(desc, "network") || strings.Contains(desc, "segmentation") {
		return "Without network controls, lateral movement is possible."
	}
	if category == "identity" {
		return "Without identity controls, authentication can be bypassed."
	}
	if category == "access" {
		return "Without access controls, unauthorized access is possible."
	}
	if category == "network" {
		return "Without network controls, traffic is uncontrolled."
	}
	return "If this assumption fails, the security control it represents is ineffective."
}

// inferControls infers specific security controls from the assumption.
func (e *NarrativeEngine) inferControls(a Assumption) []string {
	desc := strings.ToLower(a.Description)
	var controls []string

	if strings.Contains(desc, "mfa") || strings.Contains(desc, "multi-factor") {
		controls = append(controls, "multi-factor authentication")
	}
	if strings.Contains(desc, "encrypt") {
		controls = append(controls, "encryption at rest and in transit")
	}
	if strings.Contains(desc, "access") {
		controls = append(controls, "role-based access control")
	}
	if strings.Contains(desc, "logging") {
		controls = append(controls, "centralized logging and monitoring")
	}
	if strings.Contains(desc, "backup") {
		controls = append(controls, "automated backup and recovery")
	}
	if strings.Contains(desc, "firewall") {
		controls = append(controls, "network firewall rules")
	}
	if strings.Contains(desc, "patch") || strings.Contains(desc, "update") {
		controls = append(controls, "vulnerability management")
	}
	if strings.Contains(desc, "audit") {
		controls = append(controls, "regular access audits")
	}
	if strings.Contains(desc, "segmentation") {
		controls = append(controls, "network segmentation")
	}
	if strings.Contains(desc, "certificate") || strings.Contains(desc, "tls") {
		controls = append(controls, "TLS certificate management")
	}
	if strings.Contains(desc, "waf") {
		controls = append(controls, "web application firewall")
	}
	if strings.Contains(desc, "ddos") {
		controls = append(controls, "DDoS protection")
	}
	if strings.Contains(desc, "rate limit") {
		controls = append(controls, "rate limiting")
	}
	if strings.Contains(desc, "token") {
		controls = append(controls, "secure token management")
	}
	if strings.Contains(desc, "secret") {
		controls = append(controls, "secrets management")
	}
	if strings.Contains(desc, "api key") {
		controls = append(controls, "API key rotation and storage")
	}
	if strings.Contains(desc, "tenant") {
		controls = append(controls, "tenant isolation controls")
	}
	if strings.Contains(desc, "data retention") {
		controls = append(controls, "data retention policies")
	}
	if strings.Contains(desc, "pentest") || strings.Contains(desc, "penetration") {
		controls = append(controls, "regular penetration testing")
	}
	if strings.Contains(desc, "dlp") {
		controls = append(controls, "data loss prevention")
	}
	if strings.Contains(desc, "break-glass") {
		controls = append(controls, "break-glass procedures with audit logging")
	}
	if strings.Contains(desc, "hipaa") {
		controls = append(controls, "HIPAA administrative, physical, and technical safeguards")
	}
	if strings.Contains(desc, "pci") {
		controls = append(controls, "PCI DSS controls")
	}
	if strings.Contains(desc, "gdpr") {
		controls = append(controls, "GDPR data protection controls")
	}
	if strings.Contains(desc, "soc2") {
		controls = append(controls, "SOC 2 trust service criteria")
	}
	if strings.Contains(desc, "iso") {
		controls = append(controls, "ISO 27001 controls")
	}
	if strings.Contains(desc, "nist") {
		controls = append(controls, "NIST CSF controls")
	}

	return controls
}

// filterSecurityKeywords filters keywords to only security-relevant ones.
func (e *NarrativeEngine) filterSecurityKeywords(keywords []string) []string {
	var result []string
	for _, k := range keywords {
		lower := strings.ToLower(k)
		if len(k) > 3 && !isCommonWord(lower) {
			result = append(result, k)
		}
	}
	return result
}

// isCommonWord checks if a word is too common to be useful.
func isCommonWord(w string) bool {
	common := []string{"the", "and", "for", "are", "but", "not", "you", "all", "can", "had", "her", "was", "one", "our", "out", "day", "get", "has", "him", "his", "how", "its", "may", "new", "now", "old", "see", "two", "way", "who", "boy", "did", "she", "use", "her", "way", "many", "oil", "sit", "set", "run", "eat", "far", "sea", "eye", "ago", "off", "too", "any", "say", "man", "try", "ask", "end", "why", "let", "put", "say", "she", "try", "way", "own", "say", "too", "old", "tell", "very", "when", "much", "would", "there", "their", "what", "said", "each", "which", "will", "about", "could", "other", "after", "first", "never", "these", "think", "where", "being", "every", "great", "might", "shall", "still", "those", "while", "this", "that", "with", "have", "from", "they", "know", "want", "been", "good", "much", "some", "time", "very", "when", "come", "here", "just", "like", "long", "make", "many", "over", "such", "take", "than", "them", "well", "were"}
	for _, c := range common {
		if w == c {
			return true
		}
	}
	return false
}

// enforceStyle applies style rules to remove banned phrases.
func (e *NarrativeEngine) enforceStyle(text string) string {
	result := text
	for _, rule := range BannedPhrases {
		// Case-insensitive replacement using a simple approach
		result = replaceCaseInsensitive(result, rule.Pattern, rule.Replacement)

		// Also handle title case
		titlePattern := strings.Title(rule.Pattern)
		titleReplacement := strings.Title(rule.Replacement)
		result = replaceCaseInsensitive(result, titlePattern, titleReplacement)
	}
	// Clean up double spaces
	for strings.Contains(result, "  ") {
		result = strings.ReplaceAll(result, "  ", " ")
	}
	return strings.TrimSpace(result)
}

// replaceCaseInsensitive replaces all occurrences of pattern in text with replacement,
// matching case-insensitively.
func replaceCaseInsensitive(text, pattern, replacement string) string {
	lowerText := strings.ToLower(text)
	lowerPattern := strings.ToLower(pattern)

	var result strings.Builder
	start := 0
	for {
		idx := strings.Index(lowerText[start:], lowerPattern)
		if idx == -1 {
			result.WriteString(text[start:])
			break
		}
		idx += start
		result.WriteString(text[start:idx])
		result.WriteString(replacement)
		start = idx + len(pattern)
	}
	return result.String()
}
