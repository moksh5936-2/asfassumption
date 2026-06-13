package fidelity

import (
	"asf-tui/asf/fact"
	"fmt"
	"strings"
)

// RealContradiction represents a real contradiction.
// Real contradictions are:
// - Fact A vs Fact B
// - Fact vs Assumption
// NOT:
// - Self-comparison
// - Duplicate comparison
// - Same-source comparison
type RealContradiction struct {
	ID          string            `json:"id"`
	Severity    string            `json:"severity"`
	Type        string            `json:"type"` // "fact-fact", "fact-assumption"
	FactA       fact.Fact         `json:"fact_a"`
	FactB       *fact.Fact        `json:"fact_b,omitempty"`
	Assumption  *HiddenAssumption `json:"assumption,omitempty"`
	Description string            `json:"description"`
	Explanation string            `json:"explanation"`
	Resolution  string            `json:"resolution"`
}

// RealContradictionEngine detects real contradictions.
type RealContradictionEngine struct {
	factInventory      *fact.Inventory
	processedPairs     map[string]bool
	contradictionRules []ContradictionRule
}

// ContradictionRule defines a rule for detecting contradictions.
type ContradictionRule struct {
	Name        string
	Description string
	Detector    func(factA, factB fact.Fact) (bool, string)
}

// NewRealContradictionEngine creates a new real contradiction engine.
func NewRealContradictionEngine(inventory *fact.Inventory) *RealContradictionEngine {
	return &RealContradictionEngine{
		factInventory:  inventory,
		processedPairs: make(map[string]bool),
		contradictionRules: []ContradictionRule{
			{
				Name:        "mfa-required-vs-disabled",
				Description: "MFA is required but disabled",
				Detector:    detectMFAContradiction,
			},
			{
				Name:        "encryption-required-vs-disabled",
				Description: "Encryption is required but disabled",
				Detector:    detectEncryptionContradiction,
			},
			{
				Name:        "least-privilege-vs-shared-admin",
				Description: "Least privilege is enforced but shared admin exists",
				Detector:    detectLeastPrivilegeContradiction,
			},
			{
				Name:        "private-vs-public",
				Description: "Private network but public access allowed",
				Detector:    detectPrivatePublicContradiction,
			},
			{
				Name:        "immutable-log-vs-deletion",
				Description: "Immutable logs but deletion is allowed",
				Detector:    detectImmutableLogContradiction,
			},
			{
				Name:        "tls-required-vs-http",
				Description: "TLS required but HTTP traffic allowed",
				Detector:    detectTLSContradiction,
			},
			{
				Name:        "encryption-without-key-management",
				Description: "Encryption enabled but no key management",
				Detector:    detectEncryptionKeyManagementContradiction,
			},
			{
				Name:        "session-required-vs-no-rotation",
				Description: "Session security required but no rotation",
				Detector:    detectSessionRotationContradiction,
			},
			{
				Name:        "backup-required-vs-none",
				Description: "Backups required but none configured",
				Detector:    detectBackupContradiction,
			},
			{
				Name:        "audit-required-vs-disabled",
				Description: "Audit logging required but disabled",
				Detector:    detectAuditContradiction,
			},
		},
	}
}

// Detect finds real contradictions between facts.
func (e *RealContradictionEngine) Detect() []RealContradiction {
	var contradictions []RealContradiction
	facts := e.factInventory.Facts

	// Check all fact pairs
	for i := 0; i < len(facts); i++ {
		for j := i + 1; j < len(facts); j++ {
			pairKey := fmt.Sprintf("%s-%s", facts[i].ID, facts[j].ID)
			if e.processedPairs[pairKey] {
				continue
			}
			e.processedPairs[pairKey] = true

			// Check all rules
			for _, rule := range e.contradictionRules {
				if matched, explanation := rule.Detector(facts[i], facts[j]); matched {
					contradictions = append(contradictions, RealContradiction{
						ID:          fmt.Sprintf("contradiction-%s-%d", rule.Name, len(contradictions)),
						Severity:    "critical",
						Type:        "fact-fact",
						FactA:       facts[i],
						FactB:       &facts[j],
						Description: rule.Description,
						Explanation: explanation,
						Resolution:  fmt.Sprintf("Resolve %s vs %s", facts[i].Text, facts[j].Text),
					})
				}
			}
		}
	}

	return contradictions
}

// DetectFactAssumption finds contradictions between facts and assumptions.
func (e *RealContradictionEngine) DetectFactAssumption(assumptions []HiddenAssumption) []RealContradiction {
	var contradictions []RealContradiction
	facts := e.factInventory.Facts

	for _, a := range assumptions {
		for _, f := range facts {
			pairKey := fmt.Sprintf("%s-%s", f.ID, a.ID)
			if e.processedPairs[pairKey] {
				continue
			}
			e.processedPairs[pairKey] = true

			// Check if the assumption contradicts the fact
			if e.assumptionContradictsFact(f, a) {
				contradictions = append(contradictions, RealContradiction{
					ID:          fmt.Sprintf("contradiction-fact-asm-%d", len(contradictions)),
					Severity:    "critical",
					Type:        "fact-assumption",
					FactA:       f,
					Assumption:  &a,
					Description: fmt.Sprintf("Assumption contradicts fact: %s", f.Text),
					Explanation: fmt.Sprintf("Fact says '%s' but assumption '%s' contradicts it", f.Text, a.Description),
					Resolution:  "Either the fact or the assumption is wrong. Verify the architecture.",
				})
			}
		}
	}

	return contradictions
}

// assumptionContradictsFact checks if an assumption contradicts a fact.
// Only flags if the assumption is about the SAME control as the fact.
func (e *RealContradictionEngine) assumptionContradictsFact(f fact.Fact, a HiddenAssumption) bool {
	// Check if the assumption is related to the fact's control
	if !fact.IsFactRelated(f.Text, a.Description) {
		return false
	}

	// If the fact is negative and the assumption implies positive
	if f.IsNegative {
		positivePhrases := []string{"is required", "must be", "needs to be", "should be", "is enforced", "is enabled", "is configured", "is implemented", "is used", "is present", "is applied", "is mandatory", "requires", "mandates", "enforces", "ensures"}
		for _, phrase := range positivePhrases {
			if strings.Contains(strings.ToLower(a.Description), phrase) {
				return true
			}
		}
	}

	// If the fact is positive and the assumption implies negative
	if !f.IsNegative {
		negativePhrases := []string{"is disabled", "is not", "is absent", "is missing", "is not configured", "is not implemented", "is not enabled", "is optional", "is exempt", "is bypassed", "does not", "lacks", "missing", "absent"}
		for _, phrase := range negativePhrases {
			if strings.Contains(strings.ToLower(a.Description), phrase) {
				return true
			}
		}
	}

	// Direct negation
	if e.isDirectNegation(f.Text, a.Description) {
		return true
	}

	return false
}

// isDirectNegation checks if the assumption directly negates the fact.
func (e *RealContradictionEngine) isDirectNegation(factText, assumptionText string) bool {
	factNorm := fact.NormalizeFact(factText)
	assumptionNorm := fact.NormalizeFact(assumptionText)

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

	for pos, neg := range negationPatterns {
		if strings.Contains(factNorm, pos) && strings.Contains(assumptionNorm, neg) {
			return true
		}
	}

	return false
}

// Contradiction detectors

func detectMFAContradiction(f1, f2 fact.Fact) (bool, string) {
	f1Lower := strings.ToLower(f1.Text)
	f2Lower := strings.ToLower(f2.Text)

	f1HasMFA := strings.Contains(f1Lower, "mfa") || strings.Contains(f1Lower, "multi factor")
	f2HasMFA := strings.Contains(f2Lower, "mfa") || strings.Contains(f2Lower, "multi factor")

	if !f1HasMFA || !f2HasMFA {
		return false, ""
	}

	f1Positive := !f1.IsNegative
	f2Negative := f2.IsNegative

	if f1Positive && f2Negative {
		return true, fmt.Sprintf("MFA is required/enabled (%s) but also disabled/exempt (%s)", f1.Text, f2.Text)
	}
	if f2Negative && f1Positive {
		return true, fmt.Sprintf("MFA is disabled (%s) but also required (%s)", f2.Text, f1.Text)
	}

	return false, ""
}

func detectEncryptionContradiction(f1, f2 fact.Fact) (bool, string) {
	f1Lower := strings.ToLower(f1.Text)
	f2Lower := strings.ToLower(f2.Text)

	f1HasEnc := strings.Contains(f1Lower, "encryption") || strings.Contains(f1Lower, "tls") || strings.Contains(f1Lower, "ssl")
	f2HasEnc := strings.Contains(f2Lower, "encryption") || strings.Contains(f2Lower, "tls") || strings.Contains(f2Lower, "ssl")

	if !f1HasEnc || !f2HasEnc {
		return false, ""
	}

	f1Positive := !f1.IsNegative
	f2Negative := f2.IsNegative

	if f1Positive && f2Negative {
		return true, fmt.Sprintf("Encryption is required (%s) but also disabled/none (%s)", f1.Text, f2.Text)
	}
	if f2Negative && f1Positive {
		return true, fmt.Sprintf("Encryption is disabled (%s) but also required (%s)", f2.Text, f1.Text)
	}

	return false, ""
}

func detectLeastPrivilegeContradiction(f1, f2 fact.Fact) (bool, string) {
	f1Lower := strings.ToLower(f1.Text)
	f2Lower := strings.ToLower(f2.Text)

	f1HasLP := strings.Contains(f1Lower, "least privilege") || strings.Contains(f1Lower, "rbac")
	f2HasSA := strings.Contains(f2Lower, "shared") || strings.Contains(f2Lower, "admin")

	if f1HasLP && f2HasSA {
		return true, fmt.Sprintf("Least privilege is enforced (%s) but shared admin exists (%s)", f1.Text, f2.Text)
	}

	return false, ""
}

func detectPrivatePublicContradiction(f1, f2 fact.Fact) (bool, string) {
	f1Lower := strings.ToLower(f1.Text)
	f2Lower := strings.ToLower(f2.Text)

	f1HasPrivate := strings.Contains(f1Lower, "private") || strings.Contains(f1Lower, "segmented") || strings.Contains(f1Lower, "isolated")
	f2HasPublic := strings.Contains(f2Lower, "public") || strings.Contains(f2Lower, "internet") || strings.Contains(f2Lower, "accessible")

	if f1HasPrivate && f2HasPublic {
		return true, fmt.Sprintf("Private network claimed (%s) but public access allowed (%s)", f1.Text, f2.Text)
	}

	return false, ""
}

func detectImmutableLogContradiction(f1, f2 fact.Fact) (bool, string) {
	f1Lower := strings.ToLower(f1.Text)
	f2Lower := strings.ToLower(f2.Text)

	f1HasImmutable := strings.Contains(f1Lower, "immutable") || strings.Contains(f1Lower, "write-once")
	f2HasDelete := strings.Contains(f2Lower, "delete") || strings.Contains(f2Lower, "remove") || strings.Contains(f2Lower, "modify")

	if f1HasImmutable && f2HasDelete {
		return true, fmt.Sprintf("Immutable logs claimed (%s) but deletion/modification allowed (%s)", f1.Text, f2.Text)
	}

	return false, ""
}

func detectTLSContradiction(f1, f2 fact.Fact) (bool, string) {
	f1Lower := strings.ToLower(f1.Text)
	f2Lower := strings.ToLower(f2.Text)

	f1HasTLS := strings.Contains(f1Lower, "tls") || strings.Contains(f1Lower, "ssl")
	f2HasHTTP := strings.Contains(f2Lower, "http") || strings.Contains(f2Lower, "plaintext") || strings.Contains(f2Lower, "unencrypted")

	if f1HasTLS && f2HasHTTP {
		return true, fmt.Sprintf("TLS required (%s) but HTTP/plaintext allowed (%s)", f1.Text, f2.Text)
	}

	return false, ""
}

func detectEncryptionKeyManagementContradiction(f1, f2 fact.Fact) (bool, string) {
	f1Lower := strings.ToLower(f1.Text)
	f2Lower := strings.ToLower(f2.Text)

	f1HasEnc := strings.Contains(f1Lower, "encryption")
	f2HasKey := strings.Contains(f2Lower, "key") || strings.Contains(f2Lower, "kms") || strings.Contains(f2Lower, "hsm")

	if f1HasEnc && !f2HasKey {
		return false, ""
	}

	f2Negative := strings.Contains(f2Lower, "no") || strings.Contains(f2Lower, "not") || strings.Contains(f2Lower, "none")

	if f1HasEnc && f2Negative && f2HasKey {
		return true, fmt.Sprintf("Encryption enabled (%s) but no key management (%s)", f1.Text, f2.Text)
	}

	return false, ""
}

func detectSessionRotationContradiction(f1, f2 fact.Fact) (bool, string) {
	f1Lower := strings.ToLower(f1.Text)
	f2Lower := strings.ToLower(f2.Text)

	f1HasSession := strings.Contains(f1Lower, "session")
	f2HasRotate := strings.Contains(f2Lower, "rotate") || strings.Contains(f2Lower, "rotation") || strings.Contains(f2Lower, "refresh")

	if f1HasSession && !f2HasRotate {
		return false, ""
	}

	f2Negative := strings.Contains(f2Lower, "no") || strings.Contains(f2Lower, "not") || strings.Contains(f2Lower, "none")

	if f1HasSession && f2Negative && f2HasRotate {
		return true, fmt.Sprintf("Session security required (%s) but no rotation (%s)", f1.Text, f2.Text)
	}

	return false, ""
}

func detectBackupContradiction(f1, f2 fact.Fact) (bool, string) {
	f1Lower := strings.ToLower(f1.Text)
	f2Lower := strings.ToLower(f2.Text)

	f1HasBackup := strings.Contains(f1Lower, "backup")
	f2HasNone := strings.Contains(f2Lower, "none") || strings.Contains(f2Lower, "no backup") || strings.Contains(f2Lower, "not configured")

	if f1HasBackup && !f1.IsNegative && f2HasNone && f2.IsNegative {
		return true, fmt.Sprintf("Backups required (%s) but none configured (%s)", f1.Text, f2.Text)
	}

	return false, ""
}

func detectAuditContradiction(f1, f2 fact.Fact) (bool, string) {
	f1Lower := strings.ToLower(f1.Text)
	f2Lower := strings.ToLower(f2.Text)

	f1HasAudit := strings.Contains(f1Lower, "audit") || strings.Contains(f1Lower, "logging")
	f2HasDisabled := strings.Contains(f2Lower, "disabled") || strings.Contains(f2Lower, "not enabled") || strings.Contains(f2Lower, "not configured")

	if f1HasAudit && !f1.IsNegative && f2HasDisabled && f2.IsNegative {
		return true, fmt.Sprintf("Audit logging required (%s) but disabled (%s)", f1.Text, f2.Text)
	}

	return false, ""
}
