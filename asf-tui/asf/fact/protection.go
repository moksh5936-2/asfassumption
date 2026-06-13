package fact

import (
	"fmt"
	"strings"
)

// ProtectionLayer protects facts from being contradicted by assumptions.
// It is the core architectural fidelity mechanism.
type ProtectionLayer struct {
	Inventory *Inventory
}

// NewProtectionLayer creates a new protection layer.
func NewProtectionLayer(inventory *Inventory) *ProtectionLayer {
	return &ProtectionLayer{Inventory: inventory}
}

// ProtectionResult holds the result of a protection check.
type ProtectionResult struct {
	Allowed           bool   `json:"allowed"`
	Reason            string `json:"reason"`
	ContradictingFact *Fact  `json:"contradicting_fact,omitempty"`
	IsSuppression     bool   `json:"is_suppression"`
}

// CheckAssumption checks if an assumption would contradict a known fact.
// Returns ProtectionResult indicating whether the assumption is allowed.
func (pl *ProtectionLayer) CheckAssumption(assumptionText string) ProtectionResult {
	// Check against all facts
	for _, fact := range pl.Inventory.Facts {
		if pl.assumptionContradictsFact(fact, assumptionText) {
			return ProtectionResult{
				Allowed:           false,
				Reason:            fmt.Sprintf("Assumption contradicts fact: %s", fact.Text),
				ContradictingFact: &fact,
				IsSuppression:     true,
			}
		}
	}

	// Check if the assumption is a restatement of a fact
	if pl.assumptionRestatesFact(assumptionText) {
		return ProtectionResult{
			Allowed:       false,
			Reason:        "Assumption restates a known fact (not hidden)",
			IsSuppression: true,
		}
	}

	return ProtectionResult{
		Allowed:       true,
		Reason:        "Assumption does not contradict known facts",
		IsSuppression: false,
	}
}

// assumptionContradictsFact checks if an assumption contradicts a specific fact.
func (pl *ProtectionLayer) assumptionContradictsFact(fact Fact, assumptionText string) bool {
	// If the fact is negative, any positive assumption about the same control is a contradiction
	if fact.IsNegative {
		if pl.assumptionImpliesPositive(fact, assumptionText) {
			return true
		}
	}

	// If the fact is positive, any negative assumption is a contradiction
	if !fact.IsNegative {
		if pl.assumptionImpliesNegative(fact, assumptionText) {
			return true
		}
	}

	// Direct negation check
	if pl.isDirectNegation(fact.Text, assumptionText) {
		return true
	}

	// Semantic contradiction check
	if pl.isSemanticContradiction(fact, assumptionText) {
		return true
	}

	return false
}

// assumptionImpliesPositive checks if an assumption implies a positive control is present.
func (pl *ProtectionLayer) assumptionImpliesPositive(fact Fact, assumptionText string) bool {
	// Extract the control from the fact
	controlName := identifyControlName(fact.Text)
	if controlName == "" {
		return false
	}

	// Check if the assumption implies this control is present/positive
	positivePhrases := []string{
		"is required", "must be", "needs to be", "should be",
		"is enforced", "is enabled", "is configured", "is implemented",
		"is used", "is present", "is applied", "is mandatory",
		"requires", "mandates", "enforces", "ensures",
	}

	assumptionLower := strings.ToLower(assumptionText)
	for _, phrase := range positivePhrases {
		if strings.Contains(assumptionLower, phrase) && strings.Contains(assumptionLower, controlName) {
			return true
		}
	}

	return false
}

// assumptionImpliesNegative checks if an assumption implies a negative control is present.
func (pl *ProtectionLayer) assumptionImpliesNegative(fact Fact, assumptionText string) bool {
	controlName := identifyControlName(fact.Text)
	if controlName == "" {
		return false
	}

	negativePhrases := []string{
		"is disabled", "is not", "is absent", "is missing",
		"is not configured", "is not implemented", "is not enabled",
		"is optional", "is exempt", "is bypassed",
		"does not", "lacks", "missing", "absent",
	}

	assumptionLower := strings.ToLower(assumptionText)
	for _, phrase := range negativePhrases {
		if strings.Contains(assumptionLower, phrase) && strings.Contains(assumptionLower, controlName) {
			return true
		}
	}

	return false
}

// isDirectNegation checks if the assumption directly negates the fact.
func (pl *ProtectionLayer) isDirectNegation(factText, assumptionText string) bool {
	// Normalize both
	factNorm := NormalizeFact(factText)
	assumptionNorm := NormalizeFact(assumptionText)

	// Check if the assumption is the negation of the fact
	negationPatterns := map[string]string{
		"enabled":         "disabled",
		"disabled":        "enabled",
		"required":        "not required",
		"not required":    "required",
		"enforced":        "not enforced",
		"not enforced":    "enforced",
		"configured":      "not configured",
		"not configured":  "configured",
		"implemented":     "not implemented",
		"not implemented": "implemented",
		"present":         "absent",
		"absent":          "present",
		"used":            "not used",
		"not used":        "used",
	}

	// Check if the fact contains a positive and the assumption contains the negative
	for pos, neg := range negationPatterns {
		if strings.Contains(factNorm, pos) && strings.Contains(assumptionNorm, neg) {
			// Also check if they reference the same control
			factControl := identifyControlName(factText)
			assumptionControl := identifyControlName(assumptionText)
			if factControl != "" && assumptionControl != "" && factControl == assumptionControl {
				return true
			}
			// If same keywords overall
			if factControl == "" && assumptionControl == "" {
				return true
			}
		}
	}

	return false
}

// isSemanticContradiction checks for semantic contradictions.
func (pl *ProtectionLayer) isSemanticContradiction(fact Fact, assumptionText string) bool {
	// Semantic rules
	factControl := identifyControlName(fact.Text)
	assumptionControl := identifyControlName(assumptionText)

	if factControl == "" || assumptionControl == "" {
		return false
	}

	// If different controls, no contradiction
	if factControl != assumptionControl {
		return false
	}

	// Check polarity
	factPositivity := pl.getPolarity(fact.Text)
	assumptionPositivity := pl.getPolarity(assumptionText)

	// If they have opposite polarity about the same control, it's a contradiction
	if factPositivity != 0 && assumptionPositivity != 0 && factPositivity != assumptionPositivity {
		return true
	}

	return false
}

// getPolarity returns 1 for positive, -1 for negative, 0 for neutral.
func (pl *ProtectionLayer) getPolarity(text string) int {
	lower := strings.ToLower(text)

	positiveWords := []string{"enabled", "required", "enforced", "configured", "implemented", "used", "present", "applied", "mandatory", "active", "on"}
	negativeWords := []string{"disabled", "not", "absent", "missing", "optional", "exempt", "bypassed", "lacks", "inactive", "off"}

	posCount := 0
	negCount := 0

	for _, w := range positiveWords {
		if strings.Contains(lower, w) {
			posCount++
		}
	}
	for _, w := range negativeWords {
		if strings.Contains(lower, w) {
			negCount++
		}
	}

	if negCount > 0 && posCount == 0 {
		return -1
	}
	if posCount > 0 && negCount == 0 {
		return 1
	}
	return 0
}

// assumptionRestatesFact checks if the assumption just restates a known fact.
func (pl *ProtectionLayer) assumptionRestatesFact(assumptionText string) bool {
	for _, fact := range pl.Inventory.Facts {
		// If the assumption is essentially the same text as the fact
		factNorm := NormalizeFact(fact.Text)
		assumptionNorm := NormalizeFact(assumptionText)

		// If they are very similar
		if factNorm == assumptionNorm {
			return true
		}

		// If the assumption is a slightly reworded version
		if factNorm != "" && assumptionNorm != "" && pl.textSimilarity(factNorm, assumptionNorm) > 0.8 {
			return true
		}

		// Check if the assumption contains the fact text
		if strings.Contains(assumptionNorm, factNorm) || strings.Contains(factNorm, assumptionNorm) {
			return true
		}
	}
	return false
}

// textSimilarity computes a simple similarity score between two texts.
func (pl *ProtectionLayer) textSimilarity(a, b string) float64 {
	aWords := strings.Fields(a)
	bWords := strings.Fields(b)

	if len(aWords) == 0 || len(bWords) == 0 {
		return 0.0
	}

	// Count shared words
	shared := 0
	for _, aw := range aWords {
		for _, bw := range bWords {
			if aw == bw {
				shared++
				break
			}
		}
	}

	// Jaccard similarity
	union := len(aWords) + len(bWords) - shared
	if union == 0 {
		return 0.0
	}
	return float64(shared) / float64(union)
}

// identifyControlName extracts the control name from text.
func identifyControlName(text string) string {
	lower := strings.ToLower(text)

	controls := []string{
		"mfa", "multi factor", "multi-factor",
		"encryption", "tls", "ssl",
		"auth0", "auth", "authentication", "authorization",
		"oauth", "oidc", "saml", "sso",
		"backup", "recovery", "restore",
		"hipaa", "soc2", "pci dss", "iso27001", "gdpr", "nist",
		"least privilege", "rbac", "role based",
		"audit", "logging", "monitoring",
		"waf", "firewall", "ids", "ips",
		"vpn", "network", "segment", "isolation",
		"cdn", "api gateway",
		"key management", "kms", "hsm", "certificate",
		"secret", "vault", "secrets manager",
		"data retention", "disaster recovery", "business continuity",
		"penetration testing", "vulnerability", "patch", "update",
		"zero trust", "dlp", "devops", "devsecops",
		"sbom", "sast", "dast", "sca", "code signing",
	}

	for _, control := range controls {
		if strings.Contains(lower, control) {
			return control
		}
	}

	return ""
}
