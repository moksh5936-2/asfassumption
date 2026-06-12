package intelligence

import (
	"fmt"
	"strings"
)

// ContradictionEngine detects contradictions within an assumption set.
type ContradictionEngine struct{}

// NewContradictionEngine creates a contradiction engine.
func NewContradictionEngine() *ContradictionEngine {
	return &ContradictionEngine{}
}

// DetectContradictions evaluates all assumptions and returns contradictions.
func (ce *ContradictionEngine) DetectContradictions(assumptions []Assumption) []Contradiction {
	var results []Contradiction
	results = append(results, ce.detectMFAExemption(assumptions)...)
	results = append(results, ce.detectPlaintextBackup(assumptions)...)
	results = append(results, ce.detectSharedAdmin(assumptions)...)
	results = append(results, ce.detectInternetAccessiblePrivate(assumptions)...)
	results = append(results, ce.detectMutableAudit(assumptions)...)
	results = append(results, ce.detectHTTPAllowed(assumptions)...)
	results = append(results, ce.detectEncryptionWithoutKeyManagement(assumptions)...)
	results = append(results, ce.detectSessionWithoutRotation(assumptions)...)
	return results
}

// detectMFAExemption: "MFA enforced" + "service accounts exempt" → contradiction.
func (ce *ContradictionEngine) detectMFAExemption(assumptions []Assumption) []Contradiction {
	var mfaEnforced bool
	var exempted bool
	var mfaIDs, exemptIDs []string
	for _, a := range assumptions {
		desc := strings.ToLower(a.Description)
		if strings.Contains(desc, "mfa") && (strings.Contains(desc, "enforce") || strings.Contains(desc, "required") || strings.Contains(desc, "mandatory")) {
			mfaEnforced = true
			mfaIDs = append(mfaIDs, a.ID)
		}
		if strings.Contains(desc, "service account") && strings.Contains(desc, "exempt") {
			exempted = true
			exemptIDs = append(exemptIDs, a.ID)
		}
	}
	if mfaEnforced && exempted {
		return []Contradiction{{
			Severity:            RiskCritical,
			Evidence:            append(mfaIDs, exemptIDs...),
			Explanation:         "MFA is enforced for users but service accounts are exempted, creating a privileged bypass path.",
			AffectedAssumptions: append(mfaIDs, exemptIDs...),
			RuleName:            "MFA_ENFORCED_WITH_EXEMPTION",
		}}
	}
	return nil
}

// detectPlaintextBackup: "encrypted" + "plaintext backup" → contradiction.
func (ce *ContradictionEngine) detectPlaintextBackup(assumptions []Assumption) []Contradiction {
	var encrypted bool
	var plaintextBackup bool
	var encIDs, pbIDs []string
	for _, a := range assumptions {
		desc := strings.ToLower(a.Description)
		if strings.Contains(desc, "encrypt") && !strings.Contains(desc, "unencrypted") {
			encrypted = true
			encIDs = append(encIDs, a.ID)
		}
		if strings.Contains(desc, "backup") && (strings.Contains(desc, "plaintext") || strings.Contains(desc, "unencrypted")) {
			plaintextBackup = true
			pbIDs = append(pbIDs, a.ID)
		}
	}
	if encrypted && plaintextBackup {
		return []Contradiction{{
			Severity:            RiskCritical,
			Evidence:            append(encIDs, pbIDs...),
			Explanation:         "Data is encrypted at rest but backups are stored in plaintext, defeating the protection.",
			AffectedAssumptions: append(encIDs, pbIDs...),
			RuleName:            "ENCRYPTED_WITH_PLAINTEXT_BACKUP",
		}}
	}
	return nil
}

// detectSharedAdmin: "least privilege" + "shared admin account" → contradiction.
func (ce *ContradictionEngine) detectSharedAdmin(assumptions []Assumption) []Contradiction {
	var leastPrivilege bool
	var sharedAdmin bool
	var lpIDs, saIDs []string
	for _, a := range assumptions {
		desc := strings.ToLower(a.Description)
		if strings.Contains(desc, "least privilege") || strings.Contains(desc, "least-privilege") {
			leastPrivilege = true
			lpIDs = append(lpIDs, a.ID)
		}
		if strings.Contains(desc, "shared admin") || strings.Contains(desc, "shared account") || strings.Contains(desc, "generic admin") {
			sharedAdmin = true
			saIDs = append(saIDs, a.ID)
		}
	}
	if leastPrivilege && sharedAdmin {
		return []Contradiction{{
			Severity:            RiskCritical,
			Evidence:            append(lpIDs, saIDs...),
			Explanation:         "Least privilege is claimed but shared admin accounts exist, violating accountability and access control.",
			AffectedAssumptions: append(lpIDs, saIDs...),
			RuleName:            "LEAST_PRIVILEGE_WITH_SHARED_ADMIN",
		}}
	}
	return nil
}

// detectInternetAccessiblePrivate: "private subnet" + "internet accessible" → contradiction.
func (ce *ContradictionEngine) detectInternetAccessiblePrivate(assumptions []Assumption) []Contradiction {
	var privateSubnet bool
	var internetAccessible bool
	var psIDs, iaIDs []string
	for _, a := range assumptions {
		desc := strings.ToLower(a.Description)
		if strings.Contains(desc, "private subnet") || strings.Contains(desc, "private network") || strings.Contains(desc, "internal network") {
			privateSubnet = true
			psIDs = append(psIDs, a.ID)
		}
		if strings.Contains(desc, "internet accessible") || strings.Contains(desc, "publicly accessible") || strings.Contains(desc, "exposed to internet") {
			internetAccessible = true
			iaIDs = append(iaIDs, a.ID)
		}
	}
	if privateSubnet && internetAccessible {
		return []Contradiction{{
			Severity:            RiskCritical,
			Evidence:            append(psIDs, iaIDs...),
			Explanation:         "Network is claimed to be private but also described as internet accessible, creating a direct exposure.",
			AffectedAssumptions: append(psIDs, iaIDs...),
			RuleName:            "PRIVATE_SUBNET_INTERNET_ACCESSIBLE",
		}}
	}
	return nil
}

// detectMutableAudit: "immutable audit" + "log deletion allowed" → contradiction.
func (ce *ContradictionEngine) detectMutableAudit(assumptions []Assumption) []Contradiction {
	var immutableAudit bool
	var logDeletion bool
	var iaIDs, ldIDs []string
	for _, a := range assumptions {
		desc := strings.ToLower(a.Description)
		if strings.Contains(desc, "immutable audit") || strings.Contains(desc, "immutable log") || strings.Contains(desc, "tamper-proof") {
			immutableAudit = true
			iaIDs = append(iaIDs, a.ID)
		}
		if strings.Contains(desc, "log deletion") || strings.Contains(desc, "delete audit") || strings.Contains(desc, "audit purge") {
			logDeletion = true
			ldIDs = append(ldIDs, a.ID)
		}
	}
	if immutableAudit && logDeletion {
		return []Contradiction{{
			Severity:            RiskCritical,
			Evidence:            append(iaIDs, ldIDs...),
			Explanation:         "Audit logs are claimed to be immutable but log deletion is allowed, violating non-repudiation.",
			AffectedAssumptions: append(iaIDs, ldIDs...),
			RuleName:            "IMMUTABLE_AUDIT_WITH_DELETION",
		}}
	}
	return nil
}

// detectHTTPAllowed: "TLS required" + "HTTP allowed" → contradiction.
func (ce *ContradictionEngine) detectHTTPAllowed(assumptions []Assumption) []Contradiction {
	var tlsRequired bool
	var httpAllowed bool
	var tlsIDs, httpIDs []string
	for _, a := range assumptions {
		desc := strings.ToLower(a.Description)
		if strings.Contains(desc, "tls required") || strings.Contains(desc, "tls is required") || strings.Contains(desc, "tls mandatory") || strings.Contains(desc, "https only") || strings.Contains(desc, "tls enforcement") {
			tlsRequired = true
			tlsIDs = append(tlsIDs, a.ID)
		}
		if strings.Contains(desc, "http allowed") || strings.Contains(desc, "http is allowed") || strings.Contains(desc, "http permitted") || strings.Contains(desc, "plaintext http") {
			httpAllowed = true
			httpIDs = append(httpIDs, a.ID)
		}
	}
	if tlsRequired && httpAllowed {
		return []Contradiction{{
			Severity:            RiskCritical,
			Evidence:            append(tlsIDs, httpIDs...),
			Explanation:         "TLS is required but HTTP is still allowed, enabling downgrade and man-in-the-middle attacks.",
			AffectedAssumptions: append(tlsIDs, httpIDs...),
			RuleName:            "TLS_REQUIRED_HTTP_ALLOWED",
		}}
	}
	return nil
}

// detectEncryptionWithoutKeyManagement: "encryption" without "key management" → implicit contradiction.
func (ce *ContradictionEngine) detectEncryptionWithoutKeyManagement(assumptions []Assumption) []Contradiction {
	var encrypted bool
	var keyManagement bool
	var encIDs, kmIDs []string
	for _, a := range assumptions {
		desc := strings.ToLower(a.Description)
		if strings.Contains(desc, "encryption") || strings.Contains(desc, "encrypted") {
			encrypted = true
			encIDs = append(encIDs, a.ID)
		}
		if strings.Contains(desc, "key management") || strings.Contains(desc, "kms") || strings.Contains(desc, "key rotation") {
			keyManagement = true
			kmIDs = append(kmIDs, a.ID)
		}
	}
	if encrypted && !keyManagement {
		return []Contradiction{{
			Severity:            RiskHigh,
			Evidence:            encIDs,
			Explanation:         "Encryption is claimed but no key management controls are specified, making encryption incomplete and potentially ineffective.",
			AffectedAssumptions: encIDs,
			RuleName:            "ENCRYPTION_WITHOUT_KEY_MANAGEMENT",
		}}
	}
	return nil
}

// detectSessionWithoutRotation: "session" without "rotation" → implicit contradiction.
func (ce *ContradictionEngine) detectSessionWithoutRotation(assumptions []Assumption) []Contradiction {
	var session bool
	var rotation bool
	var sessIDs, rotIDs []string
	for _, a := range assumptions {
		desc := strings.ToLower(a.Description)
		if strings.Contains(desc, "session") || strings.Contains(desc, "token") || strings.Contains(desc, "jwt") {
			session = true
			sessIDs = append(sessIDs, a.ID)
		}
		if strings.Contains(desc, "rotation") || strings.Contains(desc, "refresh") || strings.Contains(desc, "renew") {
			rotation = true
			rotIDs = append(rotIDs, a.ID)
		}
	}
	if session && !rotation {
		return []Contradiction{{
			Severity:            RiskHigh,
			Evidence:            sessIDs,
			Explanation:         "Session or token management is present but rotation or renewal is not specified, increasing session hijacking risk.",
			AffectedAssumptions: sessIDs,
			RuleName:            "SESSION_WITHOUT_ROTATION",
		}}
	}
	return nil
}

// CountBySeverity returns a map of contradiction counts per severity.
func (ce *ContradictionEngine) CountBySeverity(contradictions []Contradiction) map[RiskLevel]int {
	counts := make(map[RiskLevel]int)
	for _, c := range contradictions {
		counts[c.Severity]++
	}
	return counts
}

// GetAffectedAssumptionIDs returns all unique assumption IDs affected by contradictions.
func GetAffectedAssumptionIDs(contradictions []Contradiction) []string {
	seen := make(map[string]bool)
	var result []string
	for _, c := range contradictions {
		for _, id := range c.AffectedAssumptions {
			if !seen[id] {
				seen[id] = true
				result = append(result, id)
			}
		}
	}
	return result
}

// FormatContradiction returns a human-readable string for a contradiction.
func FormatContradiction(c Contradiction) string {
	return fmt.Sprintf("[%s] %s (rule: %s) affects %d assumption(s)", c.Severity, c.Explanation, c.RuleName, len(c.AffectedAssumptions))
}
