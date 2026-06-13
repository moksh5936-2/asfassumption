package fact

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// Extractor pulls explicit facts from architecture descriptions.
type Extractor struct {
	// Patterns for positive facts (X is enabled/required/used)
	positivePatterns []string
	// Patterns for negative facts (X is disabled/not required)
	negativePatterns []string
	// Control keywords
	controlKeywords []string
	// Requirement keywords
	requirementKeywords []string
}

// NewExtractor creates a new fact extractor.
func NewExtractor() *Extractor {
	return &Extractor{
		positivePatterns: []string{
			"(?i)\\b(mfa|multi\\s*factor)\\s*(is|are)\\s*(enabled|required|enforced|used|implemented|configured)",
			"(?i)\\b(encryption|tls|ssl)\\s*(is|are)\\s*(enabled|required|enforced|used|implemented|configured|applied)",
			"(?i)\\b(encryption)\\s*(is|are)\\s*(disabled|not\\s*(enabled|required|enforced))",
			"(?i)\\b(auth0|authn|authz|oauth|oidc|saml)\\s*(is|are)\\s*(used|configured|enabled|implemented)",
			"(?i)\\b(backups?)\\s*(is|are)\\s*(enabled|configured|used|implemented|performed|automated|daily|weekly|hourly)",
			"(?i)\\b(hipaa|soc2|pci\\s*dss|iso27001|gdpr|nist)\\s*(is|are)\\s*(required|compliant|enforced|mandated)",
			"(?i)\\b(least\\s*privilege|role\\s*based|rbac)\\s*(is|are)\\s*(enforced|enabled|configured|used|implemented)",
			"(?i)\\b(audit|logging)\\s*(is|are)\\s*(enabled|configured|implemented|required|enforced|immutable)",
			"(?i)\\b(waf|firewall|ids|ips)\\s*(is|are)\\s*(enabled|configured|used|implemented|deployed)",
			"(?i)\\b(network)\\s*(is|are)\\s*(segmented|isolated|private|restricted|internal)",
			"(?i)\\b(cdn)\\s*(is|are)\\s*(used|enabled|configured|implemented)",
			"(?i)\\b(monitoring|alerting|observability)\\s*(is|are)\\s*(enabled|configured|used|implemented)",
			"(?i)\\b(vpn)\\s*(is|are)\\s*(used|enabled|required|configured)",
			"(?i)\\b(automated)\\s*(backups?|recovery|restore|failover)\\s*(is|are|enabled)",
			"(?i)\\b(data)\\s*(is|are)\\s*(encrypted|at\\s*rest|in\\s*transit|masked|tokenized|anonymized)",
			"(?i)\\b(kms|hsm|key\\s*management)\\s*(is|are)\\s*(used|enabled|configured|implemented)",
			"(?i)\\b(secrets?)\\s*(is|are)\\s*(managed|encrypted|rotated|stored|in\\s*vault|in\\s*secrets\\s*manager)",
			"(?i)\\b(certificate)\\s*(is|are)\\s*(rotated|managed|auto\\s*renewed|pinned|validated)",
		},
		negativePatterns: []string{
			"(?i)\\b(mfa|multi\\s*factor)\\s*(is|are)\\s*(disabled|not\\s*(enabled|required|enforced)|optional|exempt)",
			"(?i)\\b(encryption|tls|ssl)\\s*(is|are)\\s*(disabled|not\\s*(enabled|required|used)|none|optional)",
			"(?i)\\b(auth0|authn|authz)\\s*(is|are)\\s*(not\\s*(used|configured|enabled|implemented)|disabled)",
			"(?i)\\b(hipaa|soc2|pci\\s*dss)\\s*(is|are)\\s*(not\\s*(required|compliant)|exempt)",
			"(?i)\\b(least\\s*privilege|rbac)\\s*(is|are)\\s*(not\\s*(enforced|enabled|configured)|disabled)",
			"(?i)\\b(audit|logging)\\s*(is|are)\\s*(not\\s*(enabled|configured|implemented)|disabled|optional)",
			"(?i)\\b(waf|firewall)\\s*(is|are)\\s*(not\\s*(enabled|configured)|disabled|absent)",
			"(?i)\\b(backups?)\\s*(is|are)\\s*(not\\s*(enabled|configured|performed)|disabled|manual)",
			"(?i)\\b(network)\\s*(is|are)\\s*(not\\s*(segmented|isolated)|public|open|exposed)",
			"(?i)\\b(vpn)\\s*(is|are)\\s*(not\\s*(used|enabled|required)|disabled)",
			"(?i)\\b(kms|hsm)\\s*(is|are)\\s*(not\\s*(used|enabled)|disabled)",
			"(?i)\\b(secrets?)\\s*(is|are)\\s*(not\\s*(managed|encrypted)|in\\s*(code|config|plaintext|source\\s*code))",
		},
		controlKeywords: []string{
			"mfa", "multi-factor", "encryption", "tls", "ssl", "auth0", "authn", "authz",
			"oauth", "oidc", "saml", "backup", "recovery", "restore", "hipaa", "soc2",
			"pci dss", "iso27001", "gdpr", "nist", "least privilege", "rbac", "audit",
			"logging", "waf", "firewall", "ids", "ips", "vpn", "network", "segment",
			"isolation", "cdn", "monitoring", "alerting", "observability", "key management",
			"kms", "hsm", "certificate", "rotation", "secret", "vault", "secrets manager",
			"api gateway", "rate limiting", "ddos", "csp", "cors", "csrf", "xss",
			"input validation", "sandbox", "container", "immutable", "data retention",
			"disaster recovery", "business continuity", "penetration testing", "vulnerability",
			"scanning", "patch", "update", "zero trust", "dlp", "iac", "gitops",
			"devops", "devsecops", "shift left", "supply chain", "sbom", "sast",
			"dast", "sca", "dependency", "iast", "code signing", "notarization",
			"attestation", "tamper", "evidence", "compliance", "regulatory", "governance",
		},
		requirementKeywords: []string{
			"required", "mandatory", "must", "shall", "needed", "necessary",
			"compliance", "regulatory", "audit", "policy", "standard", "framework",
		},
	}
}

// ExtractFromText extracts facts from raw text.
func (e *Extractor) ExtractFromText(text string) *Inventory {
	inv := NewInventory()

	// Extract positive facts
	for _, pattern := range e.positivePatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllString(text, -1)
		for _, match := range matches {
			fact := e.classifyMatch(match, false)
			if fact.Text != "" && !inv.hasFact(fact.Text) {
				inv.Add(fact)
			}
		}
	}

	// Extract negative facts
	for _, pattern := range e.negativePatterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllString(text, -1)
		for _, match := range matches {
			fact := e.classifyMatch(match, true)
			if fact.Text != "" && !inv.hasFact(fact.Text) {
				inv.Add(fact)
			}
		}
	}

	// Extract requirement statements
	inv = e.extractRequirements(text, inv)

	// Extract control declarations from security_controls section
	inv = e.extractControlDeclarations(text, inv)

	return inv
}

// ExtractFromYAML extracts facts from structured YAML sections.
func (e *Extractor) ExtractFromYAML(securityControls map[string][]string, compliance []string, requirements []string, constraints []string) *Inventory {
	inv := NewInventory()

	// Extract security controls
	for controlName, values := range securityControls {
		for _, value := range values {
			fact := e.classifyControl(controlName, value)
			if fact.Text != "" && !inv.hasFact(fact.Text) {
				inv.Add(fact)
			}
		}
	}

	// Extract compliance facts
	for _, comp := range compliance {
		if comp != "" && !inv.hasFact(comp) {
			inv.Add(Fact{
				ID:         fmt.Sprintf("fact-compliance-%d", inv.Count()),
				Text:       comp,
				Source:     "yaml",
				Confidence: 1.0,
				Category:   "compliance",
				FactType:   "compliance",
			})
		}
	}

	// Extract requirements
	for _, req := range requirements {
		if req != "" && !inv.hasFact(req) {
			inv.Add(Fact{
				ID:         fmt.Sprintf("fact-req-%d", inv.Count()),
				Text:       req,
				Source:     "yaml",
				Confidence: 1.0,
				Category:   "compliance",
				FactType:   "requirement",
			})
		}
	}

	// Extract constraints
	for _, cons := range constraints {
		if cons != "" && !inv.hasFact(cons) {
			inv.Add(Fact{
				ID:         fmt.Sprintf("fact-const-%d", inv.Count()),
				Text:       cons,
				Source:     "yaml",
				Confidence: 1.0,
				Category:   "operational",
				FactType:   "constraint",
			})
		}
	}

	return inv
}

// hasFact checks if an inventory already contains a similar fact.
func (inv *Inventory) hasFact(text string) bool {
	for _, f := range inv.Facts {
		if strings.EqualFold(f.Text, text) {
			return true
		}
	}
	return false
}

// classifyMatch turns a regex match into a structured fact.
func (e *Extractor) classifyMatch(match string, isNegative bool) Fact {
	// Determine the control being described
	controlName := e.identifyControl(match)
	factType := "control"
	category := "security"

	// Determine severity
	severity := "medium"
	if isNegative {
		severity = "high"
	}
	if strings.Contains(strings.ToLower(match), "mfa") || strings.Contains(strings.ToLower(match), "encryption") {
		severity = "critical"
	}

	return Fact{
		ID:         fmt.Sprintf("fact-%s-%d", strings.ToLower(controlName), len(match)),
		Text:       match,
		Source:     "text",
		Confidence: 0.95,
		Category:   category,
		FactType:   factType,
		IsNegative: isNegative,
		Severity:   severity,
	}
}

// classifyControl turns a structured control declaration into a fact.
func (e *Extractor) classifyControl(controlName, value string) Fact {
	factType := "control"
	category := "security"

	// Determine if negative
	isNegative := false
	negWords := []string{"disabled", "not enabled", "not required", "not configured", "not implemented", "none", "optional", "absent"}
	for _, nw := range negWords {
		if strings.Contains(strings.ToLower(value), nw) {
			isNegative = true
			break
		}
	}

	// Determine severity
	severity := "medium"
	if isNegative {
		severity = "high"
	}
	if strings.Contains(strings.ToLower(controlName), "mfa") || strings.Contains(strings.ToLower(controlName), "encryption") {
		severity = "critical"
	}

	return Fact{
		ID:         fmt.Sprintf("fact-ctrl-%s", strings.ToLower(controlName)),
		Text:       fmt.Sprintf("%s: %s", controlName, value),
		Source:     "yaml",
		Confidence: 1.0,
		Category:   category,
		FactType:   factType,
		IsNegative: isNegative,
		Severity:   severity,
	}
}

// identifyControl extracts the control name from a matched text.
func (e *Extractor) identifyControl(text string) string {
	// Try to find a control keyword
	lower := strings.ToLower(text)
	for _, kw := range e.controlKeywords {
		if strings.Contains(lower, kw) {
			return kw
		}
	}
	// Fallback: extract the first meaningful word
	words := strings.FieldsFunc(text, func(c rune) bool {
		return !unicode.IsLetter(c) && !unicode.IsDigit(c)
	})
	if len(words) > 0 {
		return words[0]
	}
	return "unknown"
}

// extractRequirements pulls requirement statements from text.
func (e *Extractor) extractRequirements(text string, inv *Inventory) *Inventory {
	// Pattern: "X must be Y" or "X is required" or "X shall be Y"
	requirementRe := regexp.MustCompile(`(?i)([A-Za-z][A-Za-z\s]*(?:must|shall|required|needs|necessary))\s*(?:be|have|use|implement|configure|enable|deploy|maintain)?\s*([^.,;\n]*)`)

	matches := requirementRe.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 2 && match[1] != "" {
			req := strings.TrimSpace(match[1])
			if req != "" && !inv.hasFact(req) {
				inv.Add(Fact{
					ID:         fmt.Sprintf("fact-req-%d", inv.Count()),
					Text:       req,
					Source:     "text",
					Confidence: 0.9,
					Category:   "compliance",
					FactType:   "requirement",
				})
			}
		}
	}

	return inv
}

// extractControlDeclarations extracts structured control declarations.
func (e *Extractor) extractControlDeclarations(text string, inv *Inventory) *Inventory {
	// Pattern: "security_controls: X enabled" or "controls: X"
	controlRe := regexp.MustCompile(`(?i)(security_controls?|controls?|policies?|configurations?)\s*[:=]\s*(.+)`)

	matches := controlRe.FindAllStringSubmatch(text, -1)
	for _, match := range matches {
		if len(match) >= 2 && match[2] != "" {
			ctrl := strings.TrimSpace(match[2])
			if ctrl != "" && !inv.hasFact(ctrl) {
				inv.Add(Fact{
					ID:         fmt.Sprintf("fact-ctrl-dec-%d", inv.Count()),
					Text:       ctrl,
					Source:     "text",
					Confidence: 0.9,
					Category:   "security",
					FactType:   "control",
				})
			}
		}
	}

	return inv
}

// ComponentFactExtractor extracts facts from component labels.
type ComponentFactExtractor struct{}

// NewComponentFactExtractor creates a new component fact extractor.
func NewComponentFactExtractor() *ComponentFactExtractor {
	return &ComponentFactExtractor{}
}

// ExtractFromComponent extracts facts from a component label.
func (e *ComponentFactExtractor) ExtractFromComponent(id, label string) *Inventory {
	inv := NewInventory()
	lower := strings.ToLower(label)

	// Check for explicit control keywords in component name
	if strings.Contains(lower, "mfa") || strings.Contains(lower, "auth") {
		inv.Add(Fact{
			ID:          fmt.Sprintf("fact-comp-%s-auth", id),
			Text:        fmt.Sprintf("Component %s (%s) is authentication-related", id, label),
			Source:      "component",
			Confidence:  0.9,
			Category:    "security",
			FactType:    "control",
			ComponentID: id,
		})
	}

	if strings.Contains(lower, "kms") || strings.Contains(lower, "hsm") || strings.Contains(lower, "vault") {
		inv.Add(Fact{
			ID:          fmt.Sprintf("fact-comp-%s-kms", id),
			Text:        fmt.Sprintf("Component %s (%s) is key-management-related", id, label),
			Source:      "component",
			Confidence:  0.9,
			Category:    "security",
			FactType:    "control",
			ComponentID: id,
		})
	}

	if strings.Contains(lower, "waf") || strings.Contains(lower, "firewall") {
		inv.Add(Fact{
			ID:          fmt.Sprintf("fact-comp-%s-waf", id),
			Text:        fmt.Sprintf("Component %s (%s) is firewall/WAF-related", id, label),
			Source:      "component",
			Confidence:  0.9,
			Category:    "security",
			FactType:    "control",
			ComponentID: id,
		})
	}

	if strings.Contains(lower, "log") || strings.Contains(lower, "audit") {
		inv.Add(Fact{
			ID:          fmt.Sprintf("fact-comp-%s-log", id),
			Text:        fmt.Sprintf("Component %s (%s) is logging/audit-related", id, label),
			Source:      "component",
			Confidence:  0.9,
			Category:    "security",
			FactType:    "control",
			ComponentID: id,
		})
	}

	return inv
}

// NormalizeFact normalizes a fact text for comparison.
func NormalizeFact(text string) string {
	// Lowercase, remove extra spaces, remove punctuation
	text = strings.ToLower(text)
	text = strings.ReplaceAll(text, ",", " ")
	text = strings.ReplaceAll(text, ".", " ")
	text = strings.ReplaceAll(text, ";", " ")
	text = strings.ReplaceAll(text, ":", " ")
	text = strings.Join(strings.Fields(text), " ")
	return text
}

// IsFactRelated checks if a given text is related to a fact.
func IsFactRelated(fact, text string) bool {
	fNorm := NormalizeFact(fact)
	tNorm := NormalizeFact(text)

	// If the text contains the core fact keywords
	factWords := strings.Fields(fNorm)
	matchCount := 0
	for _, fw := range factWords {
		if len(fw) > 3 && strings.Contains(tNorm, fw) {
			matchCount++
		}
	}

	// If more than 50% of fact words match, it's related
	if len(factWords) > 0 && float64(matchCount)/float64(len(factWords)) >= 0.5 {
		return true
	}

	return false
}
