package intelligence

import (
	"fmt"
	"strings"
)

type SemanticPair struct {
	Category          ContradictionType
	ConceptA          []string
	ConceptB          []string
	Severity          RiskLevel
	BaseConfidence    float64
	Reasoning         string
	Impact            string
	Recommendations   []string
}

type SemanticEngine struct {
	pairs      []SemanticPair
	synonyms   map[string][]string
	normalizer map[string]string
}

func NewSemanticEngine() *SemanticEngine {
	se := &SemanticEngine{
		pairs:    defaultPairs(),
		synonyms: defaultSynonyms(),
	}
	return se
}

func defaultSynonyms() map[string][]string {
	return map[string][]string{
		"encrypted":            {"encryption", "encrypt", "enciphered", "cipher", "scrambled", "aes", "cryptographic"},
		"plaintext":            {"unencrypted", "not encrypted", "cleartext", "in the clear", "no encryption", "without encryption", "decrypted"},
		"mfa":                  {"multi-factor", "multi factor", "two-factor", "two factor", "2fa", "strong authentication"},
		"single factor":        {"password only", "password based only", "only password", "single factor", "no mfa", "without mfa", "no two-factor", "no multi-factor", "password without mfa"},
		"anonymous":            {"unauthenticated", "no login", "no authentication", "public access", "guest access"},
		"admin only":           {"administrators only", "admin restricted", "restricted to admin", "admin access only"},
		"everyone access":      {"all users access", "open to all", "public access", "unrestricted access", "anyone can access"},
		"least privilege":      {"least-privilege", "minimum necessary", "need to know", "minimum access", "principle of least privilege"},
		"full access":          {"unrestricted access", "admin access for all", "all users admin", "everyone admin", "full admin"},
		"restricted access":    {"limited access", "access restricted", "controlled access", "authorized only"},
		"open access":          {"no access control", "anyone can access", "without authorization", "no authorization required", "unrestricted access"},
		"staff only":           {"employees only", "internal staff", "staff access only"},
		"vendor access":        {"third party access", "external access", "third-party access", "contractor access"},
		"all access logged":    {"audit logging enabled", "logging required", "audit all access", "log all access", "audit enabled", "access logging"},
		"access not logged":    {"logging disabled", "no audit", "audit disabled", "not logged", "no logging"},
		"audit trail present":  {"audit trail required", "audit trail enabled", "immutable audit", "tamper-proof audit"},
		"audit disabled":       {"audit not enabled", "no audit trail", "audit trail absent", "audit turned off"},
		"monitoring enabled":   {"monitoring active", "active monitoring", "monitoring in place", "detection enabled"},
		"monitoring absent":    {"no monitoring", "monitoring disabled", "no detection", "unmonitored"},
		"private network":      {"internal network", "no public access", "private subnet", "isolated network", "vpc private"},
		"public exposure":      {"publicly accessible", "internet accessible", "exposed to internet", "internet facing", "public endpoint"},
		"isolated segment":     {"network isolation", "air gapped", "isolated network", "sandboxed", "vlan isolated"},
		"direct access":        {"directly accessible", "no network isolation", "direct connection", "unrestricted network"},
		"segmented":            {"network segmentation", "vlan", "dmz", "segmented network", "micro-segmented"},
		"flat network":         {"no segmentation", "single network", "unsegmented", "flat topology"},
		"internal only":        {"internal access only", "no external access", "intranet only", "internal facing"},
		"internet reachable":   {"internet facing", "public endpoint", "exposed to internet", "publicly reachable"},
		"backups encrypted":    {"encrypted backup", "backup encryption", "backup encrypted"},
		"backups plaintext":    {"plaintext backup", "unencrypted backup", "backup not encrypted", "backup in clear"},
		"redundant":            {"redundancy", "multiple copies", "multi-replica", "replicated"},
		"single copy":          {"no redundancy", "not replicated", "single replica", "no backup copy"},
		"replicated":           {"geo-replicated", "multi-region", "cross-region", "multi-az", "multiple regions"},
		"single region":        {"not replicated", "single location", "single region", "no replication"},
		"high availability":    {"ha", "highly available", "active-active", "active/passive", "failover", "clustered"},
		"single instance":      {"single point of failure", "no redundancy", "single node", "non-redundant", "spof"},
		"redundant path":       {"multiple paths", "multi-path", "redundant connection"},
		"single path":          {"single connection", "no redundant path", "single link"},
		"zero trust":           {"zero-trust", "never trust always verify"},
		"implicit trust":       {"trusted by default", "trust but verify", "trust all", "implicitly trusted"},
		"verified identity":    {"identity verification", "identity proofing"},
		"trusted by default":   {"implicitly trusted", "default trust", "auto-trusted"},
		"hipaa compliant":      {"hipaa compliance", "hipaa", "hipaa required"},
		"unencrypted phi":      {"plaintext phi", "phi not encrypted", "no phi encryption", "phi exposed", "without encryption", "phi without encryption"},
		"pci compliant":        {"pci dss", "pci compliance", "pci required", "pci-dss"},
		"card data exposed":    {"unencrypted card", "cardholder data in plaintext", "card data plaintext"},
		"soc2 compliant":       {"soc2", "soc 2", "soc2 compliance"},
		"logging absent":       {"no logging", "audit disabled", "not logged", "no audit trail"},
	}
}

func defaultPairs() []SemanticPair {
	return []SemanticPair{
		// ─── ENCRYPTION ─────────────────────────────────────
		{
			Category:       ContradictionTypeENCRYPTION,
			ConceptA:       []string{"encrypted"},
			ConceptB:       []string{"plaintext"},
			Severity:       RiskCritical,
			BaseConfidence: 0.95,
			Reasoning:      "Data cannot be both encrypted and stored in plaintext. This contradicts the fundamental principle of data protection.",
			Impact:         "Data protection controls are inconsistent and may leave sensitive data exposed.",
			Recommendations: []string{
				"Encrypt all data consistently across primary and secondary storage",
				"Verify encryption policies apply to all data pipelines",
				"Conduct encryption audit across all storage tiers",
			},
		},
		{
			Category:       ContradictionTypeENCRYPTION,
			ConceptA:       []string{"encrypted"},
			ConceptB:       []string{"plaintext"},
			Severity:       RiskCritical,
			BaseConfidence: 0.92,
			Reasoning:      "Data encrypted at rest must also be encrypted in backups. Plaintext backups defeat the purpose of encryption.",
			Impact:         "Backup storage becomes an easy target for data exfiltration despite encryption controls on primary storage.",
			Recommendations: []string{
				"Are backups encrypted?",
				"Are encryption policies enforced on backup systems?",
				"Implement key management for backup encryption keys",
			},
		},
		// ─── AUTHENTICATION ────────────────────────────────
		{
			Category:       ContradictionTypeAUTHENTICATION,
			ConceptA:       []string{"mfa"},
			ConceptB:       []string{"single factor"},
			Severity:       RiskHigh,
			BaseConfidence: 0.90,
			Reasoning:      "MFA and single-factor authentication are mutually exclusive. If some access paths allow single factor, MFA enforcement is incomplete.",
			Impact:         "Authentication controls are inconsistent, creating a weaker link that attackers can exploit.",
			Recommendations: []string{
				"Enforce MFA across all access paths without exception",
				"Implement compensating controls where MFA is not feasible",
				"Audit all authentication paths for MFA enforcement",
			},
		},
		{
			Category:       ContradictionTypeAUTHENTICATION,
			ConceptA:       []string{"mfa"},
			ConceptB:       []string{"anonymous"},
			Severity:       RiskHigh,
			BaseConfidence: 0.88,
			Reasoning:      "MFA and anonymous access are direct contradictions. Any anonymous access path bypasses all authentication controls.",
			Impact:         "Unauthenticated access paths exist despite strong authentication controls.",
			Recommendations: []string{
				"Remove anonymous access or isolate it to non-sensitive services",
				"Require authentication for all endpoints",
			},
		},
		{
			Category:       ContradictionTypeAUTHENTICATION,
			ConceptA:       []string{"verified identity"},
			ConceptB:       []string{"anonymous"},
			Severity:       RiskHigh,
			BaseConfidence: 0.87,
			Reasoning:      "If identity is verified for some users but anonymous access is allowed, the identity verification is bypassable.",
			Impact:         "Identity verification controls are undermined by anonymous access paths.",
			Recommendations: []string{
				"Extend identity verification to all access paths",
				"Remove anonymous access or segregate to isolated environment",
			},
		},
		// ─── AUTHORIZATION ─────────────────────────────────
		{
			Category:       ContradictionTypeAUTHORIZATION,
			ConceptA:       []string{"restricted access"},
			ConceptB:       []string{"open access"},
			Severity:       RiskCritical,
			BaseConfidence: 0.93,
			Reasoning:      "A system cannot simultaneously restrict and open access. This creates direct conflict in access control enforcement.",
			Impact:         "Access control is inconsistent. Sensitive resources may be exposed despite stated restrictions.",
			Recommendations: []string{
				"Apply access controls consistently across all resources",
				"Authorize sensitive resources only",
			},
		},
		{
			Category:       ContradictionTypeAUTHORIZATION,
			ConceptA:       []string{"least privilege"},
			ConceptB:       []string{"full access"},
			Severity:       RiskCritical,
			BaseConfidence: 0.95,
			Reasoning:      "Least privilege and full access are mutually exclusive. Granting full access to anyone violates the principle of least privilege.",
			Impact:         "Every user with full access represents an unnecessary privilege that expands the attack surface.",
			Recommendations: []string{
				"Implement role-based access with tiered permissions",
				"Remove full-access roles from all non-essential users",
				"Audit and recertify all privileged access",
			},
		},
		{
			Category:       ContradictionTypeAUTHORIZATION,
			ConceptA:       []string{"staff only"},
			ConceptB:       []string{"vendor access"},
			Severity:       RiskHigh,
			BaseConfidence: 0.85,
			Reasoning:      "If access is restricted to staff, vendor access creates a direct contradiction that bypasses the staff-only policy.",
			Impact:         "Third-party access may expose systems to vendor risk despite staff-only access controls.",
			Recommendations: []string{
				"Document and approve all vendor access with data processing agreements",
				"Implement vendor access monitoring with anomaly detection",
				"Time-box vendor access with automatic revocation",
			},
		},
		{
			Category:       ContradictionTypeAUTHORIZATION,
			ConceptA:       []string{"admin only"},
			ConceptB:       []string{"everyone access"},
			Severity:       RiskCritical,
			BaseConfidence: 0.92,
			Reasoning:      "Admin-only access and everyone access are direct contradictions. If everyone has access, the admin-only restriction is meaningless.",
			Impact:         "Sensitive admin functionality may be exposed to all users.",
			Recommendations: []string{
				"Implement proper role separation",
				"Restrict admin functions to authorized administrators",
			},
		},
		// ─── LOGGING ───────────────────────────────────────
		{
			Category:       ContradictionTypeLOGGING,
			ConceptA:       []string{"all access logged"},
			ConceptB:       []string{"access not logged"},
			Severity:       RiskHigh,
			BaseConfidence: 0.90,
			Reasoning:      "It is impossible to log all access if some access paths have logging disabled.",
			Impact:         "Security monitoring gaps exist, allowing undetected unauthorized access.",
			Recommendations: []string{
				"Enable logging on all access paths without exception",
				"Verify logging coverage with periodic audits",
			},
		},
		{
			Category:       ContradictionTypeLOGGING,
			ConceptA:       []string{"audit trail present"},
			ConceptB:       []string{"audit disabled"},
			Severity:       RiskHigh,
			BaseConfidence: 0.88,
			Reasoning:      "An audit trail cannot be present if audit is disabled. This creates a gap in audit coverage.",
			Impact:         "Non-repudiation and forensic capabilities are compromised.",
			Recommendations: []string{
				"Enable audit trails across all systems",
				"Implement immutable audit storage",
			},
		},
		{
			Category:       ContradictionTypeLOGGING,
			ConceptA:       []string{"monitoring enabled"},
			ConceptB:       []string{"monitoring absent"},
			Severity:       RiskMedium,
			BaseConfidence: 0.82,
			Reasoning:      "Monitoring cannot be both enabled and absent. This indicates inconsistent security observability.",
			Impact:         "Parts of the system may be unmonitored, allowing attacks to go undetected.",
			Recommendations: []string{
				"Implement comprehensive monitoring coverage",
				"Identify and remediate monitoring gaps",
			},
		},
		// ─── NETWORK ───────────────────────────────────────
		{
			Category:       ContradictionTypeNETWORK,
			ConceptA:       []string{"private network"},
			ConceptB:       []string{"public exposure"},
			Severity:       RiskHigh,
			BaseConfidence: 0.90,
			Reasoning:      "A network cannot be both private and publicly exposed. This creates a direct network security contradiction.",
			Impact:         "Network isolation claims are invalidated by public exposure.",
			Recommendations: []string{
				"Remove public endpoints or move to DMZ",
				"Implement VPN for external access",
			},
		},
		{
			Category:       ContradictionTypeNETWORK,
			ConceptA:       []string{"isolated segment"},
			ConceptB:       []string{"direct access"},
			Severity:       RiskHigh,
			BaseConfidence: 0.88,
			Reasoning:      "An isolated segment with direct access paths is effectively not isolated.",
			Impact:         "Network isolation controls are bypassed by direct access paths.",
			Recommendations: []string{
				"Enforce network access controls at all connection points",
				"Remove direct access paths to isolated segments",
			},
		},
		{
			Category:       ContradictionTypeNETWORK,
			ConceptA:       []string{"segmented"},
			ConceptB:       []string{"flat network"},
			Severity:       RiskMedium,
			BaseConfidence: 0.85,
			Reasoning:      "Network segmentation and flat network topology are directly contradictory.",
			Impact:         "Network is not meaningfully segmented despite segmentation claims.",
			Recommendations: []string{
				"Implement proper network segmentation with firewall rules",
				"Verify segmentation effectiveness with penetration testing",
			},
		},
		{
			Category:       ContradictionTypeNETWORK,
			ConceptA:       []string{"internal only"},
			ConceptB:       []string{"internet reachable"},
			Severity:       RiskHigh,
			BaseConfidence: 0.90,
			Reasoning:      "An internal-only system cannot be reachable from the internet. These claims are mutually exclusive.",
			Impact:         "Internal systems are exposed to internet-borne attacks.",
			Recommendations: []string{
				"Remove internet-facing endpoints from internal systems",
				"Implement proper network segmentation with internet-facing DMZ",
			},
		},
		// ─── BACKUP ────────────────────────────────────────
		{
			Category:       ContradictionTypeBACKUP,
			ConceptA:       []string{"backups encrypted"},
			ConceptB:       []string{"backups plaintext"},
			Severity:       RiskCritical,
			BaseConfidence: 0.92,
			Reasoning:      "Backup encryption policies are inconsistent. Some backups are encrypted while others are not.",
			Impact:         "Unencrypted backups expose data despite encryption policies for primary backup targets.",
			Recommendations: []string{
				"Enforce encryption on all backup targets",
				"Audit backup encryption status regularly",
			},
		},
		{
			Category:       ContradictionTypeBACKUP,
			ConceptA:       []string{"redundant"},
			ConceptB:       []string{"single copy"},
			Severity:       RiskMedium,
			BaseConfidence: 0.80,
			Reasoning:      "Redundancy and single copy are directly contradictory. Data resilience is compromised.",
			Impact:         "Single copy represents a single point of failure for data durability.",
			Recommendations: []string{
				"Implement data replication to at least two locations",
				"Test failover and restore from each replica",
			},
		},
		{
			Category:       ContradictionTypeBACKUP,
			ConceptA:       []string{"replicated"},
			ConceptB:       []string{"single region"},
			Severity:       RiskMedium,
			BaseConfidence: 0.80,
			Reasoning:      "Geo-replication and single-region storage are directly contradictory.",
			Impact:         "Single-region deployment is vulnerable to region-wide outages.",
			Recommendations: []string{
				"Deploy to multiple regions for geo-redundancy",
				"Test cross-region failover and data replication",
			},
		},
		// ─── AVAILABILITY ─────────────────────────────────
		{
			Category:       ContradictionTypeAVAILABILITY,
			ConceptA:       []string{"high availability"},
			ConceptB:       []string{"single instance"},
			Severity:       RiskHigh,
			BaseConfidence: 0.88,
			Reasoning:      "A single instance cannot provide high availability. This contradicts the fundamental definition of HA.",
			Impact:         "The system has a single point of failure despite high availability claims.",
			Recommendations: []string{
				"Deploy at least two instances behind a load balancer",
				"Implement health checks and automated failover",
			},
		},
		{
			Category:       ContradictionTypeAVAILABILITY,
			ConceptA:       []string{"redundant path"},
			ConceptB:       []string{"single path"},
			Severity:       RiskMedium,
			BaseConfidence: 0.82,
			Reasoning:      "Redundant paths and a single path are directly contradictory. Network resilience is compromised.",
			Impact:         "A single path failure causes complete connectivity loss despite redundancy claims.",
			Recommendations: []string{
				"Implement redundant network paths with automatic failover",
				"Test path failover scenarios regularly",
			},
		},
		// ─── TRUST ─────────────────────────────────────────
		{
			Category:       ContradictionTypeTRUST,
			ConceptA:       []string{"zero trust"},
			ConceptB:       []string{"implicit trust"},
			Severity:       RiskHigh,
			BaseConfidence: 0.90,
			Reasoning:      "Zero trust and implicit trust are fundamentally incompatible security models.",
			Impact:         "Trust model is inconsistent, undermining the security architecture.",
			Recommendations: []string{
				"Adopt a consistent trust model across the architecture",
				"Implement verify-explicitly for all access decisions",
			},
		},
		{
			Category:       ContradictionTypeTRUST,
			ConceptA:       []string{"verified identity"},
			ConceptB:       []string{"trusted by default"},
			Severity:       RiskHigh,
			BaseConfidence: 0.85,
			Reasoning:      "If some entities are trusted by default, the identity verification claim is not universally applied.",
			Impact:         "Untrusted entities may gain access through default-trust paths.",
			Recommendations: []string{
				"Remove default-trust relationships",
				"Verify identity for all access requests",
			},
		},
		// ─── COMPLIANCE ───────────────────────────────────
		{
			Category:       ContradictionTypeCOMPLIANCE,
			ConceptA:       []string{"hipaa compliant"},
			ConceptB:       []string{"unencrypted phi"},
			Severity:       RiskCritical,
			BaseConfidence: 0.95,
			Reasoning:      "HIPAA requires encryption of PHI at rest and in transit. Unencrypted PHI directly violates HIPAA Security Rule.",
			Impact:         "HIPAA compliance is not achievable with unencrypted PHI.",
			Recommendations: []string{
				"Encrypt all PHI at rest and in transit",
				"Implement key management for encryption keys",
				"Conduct HIPAA risk assessment",
			},
		},
		{
			Category:       ContradictionTypeCOMPLIANCE,
			ConceptA:       []string{"pci compliant"},
			ConceptB:       []string{"card data exposed"},
			Severity:       RiskCritical,
			BaseConfidence: 0.95,
			Reasoning:      "PCI DSS prohibits storing sensitive cardholder data unencrypted. Exposed card data violates PCI DSS requirements.",
			Impact:         "PCI DSS compliance is not achievable with exposed cardholder data.",
			Recommendations: []string{
				"Encrypt cardholder data at rest with strong cryptography",
				"Tokenize or truncate PAN where possible",
				"Conduct PCI DSS assessment",
			},
		},
		{
			Category:       ContradictionTypeCOMPLIANCE,
			ConceptA:       []string{"soc2 compliant"},
			ConceptB:       []string{"logging absent"},
			Severity:       RiskHigh,
			BaseConfidence: 0.88,
			Reasoning:      "SOC2 requires monitoring and logging (CC7.2). Without logging, SOC2 compliance cannot be achieved.",
			Impact:         "SOC2 audit will fail due to insufficient logging and monitoring.",
			Recommendations: []string{
				"Implement comprehensive audit logging",
				"Enable monitoring for all security-relevant events",
			},
		},
	}
}

func (se *SemanticEngine) matchConcepts(text string, concepts []string) bool {
	lower := strings.ToLower(text)
	for _, concept := range concepts {
		if se.textMatches(lower, concept) {
			return true
		}
		if syns, ok := se.synonyms[concept]; ok {
			for _, syn := range syns {
				if se.textMatches(lower, syn) {
					return true
				}
			}
		}
	}
	return false
}

func (se *SemanticEngine) textMatches(text, phrase string) bool {
	if phrase == "" {
		return false
	}
	if strings.Contains(text, phrase) {
		if se.isNegated(text, phrase) {
			return false
		}
		words := strings.Fields(phrase)
		checkWord := phrase
		if len(words) > 1 {
			checkWord = words[0]
		}
		if se.isUnPrefixed(text, checkWord) {
			return false
		}
		return true
	}
	words := strings.Fields(phrase)
	if len(words) <= 1 {
		return false
	}
	for _, w := range words {
		if !strings.Contains(text, w) || se.isNegated(text, w) {
			return false
		}
		if se.isUnPrefixed(text, w) {
			return false
		}
	}
	return true
}

func (se *SemanticEngine) isNegated(text, phrase string) bool {
	lower := strings.ToLower(text)
	idx := strings.Index(lower, phrase)
	if idx < 0 {
		return false
	}
	if idx > 0 && lower[idx-1] == ' ' {
		prefix := lower[:idx-1]
		lastSpace := strings.LastIndex(prefix, " ")
		prevWord := prefix
		if lastSpace >= 0 {
			prevWord = prefix[lastSpace+1:]
		}
		if prevWord == "no" || prevWord == "not" || prevWord == "without" {
			return true
		}
	}
	return false
}

func (se *SemanticEngine) isUnPrefixed(text, word string) bool {
	if len(word) < 2 {
		return false
	}
	lower := strings.ToLower(text)
	idx := strings.Index(lower, word)
	if idx < 0 {
		return false
	}
	if idx >= 2 && lower[idx-2:idx] == "un" {
		before := idx - 2
		if before == 0 || lower[before-1] == ' ' || lower[before-1] == '-' || lower[before-1] == '_' {
			return true
		}
	}
	return false
}

func (se *SemanticEngine) matchesOpposite(text string, pair SemanticPair) bool {
	for _, ca := range pair.ConceptA {
		if se.textMatches(strings.ToLower(text), ca) {
			return true
		}
		if syns, ok := se.synonyms[ca]; ok {
			for _, syn := range syns {
				if se.textMatches(strings.ToLower(text), syn) {
					return true
				}
			}
		}
	}
	return false
}

func (se *SemanticEngine) DetectSemanticContradictions(claims []Statement) []CIEContradiction {
	var results []CIEContradiction

	var filtered []Statement
	for _, c := range claims {
		// Skip structured labels (security_controls) — they are identifiers, not semantic statements
		if c.Source == "security_controls" {
			continue
		}
		filtered = append(filtered, c)
	}
	claims = filtered

	for i := 0; i < len(claims); i++ {
		for j := i + 1; j < len(claims); j++ {
			a, b := claims[i], claims[j]
			if a.ID == b.ID {
				continue
			}

			for _, pair := range se.pairs {
				aA := se.matchConcepts(a.OriginalText, pair.ConceptA)
				aB := se.matchConcepts(a.OriginalText, pair.ConceptB)
				bA := se.matchConcepts(b.OriginalText, pair.ConceptA)
				bB := se.matchConcepts(b.OriginalText, pair.ConceptB)

				if aA && aB {
					continue
				}
				if bA && bB {
					continue
				}

				if (aA && bB) || (aB && bA) {
					if !se.isContextExcluded(a, b, pair) {
						results = append(results, se.buildContradiction(a, b, pair))
					}
				}
			}
		}
	}

	return se.deduplicate(results)
}

func (se *SemanticEngine) isContextExcluded(a, b Statement, pair SemanticPair) bool {
	if pair.Category == ContradictionTypeENCRYPTION {
		aText := strings.ToLower(a.OriginalText)
		bText := strings.ToLower(b.OriginalText)
		storageTerms := []string{"backup", "at rest", "storage", "disk", "database", "s3", "blob", "file", "archive"}
		transportTerms := []string{"traffic", "transit", "network", "http", "https", "tls", "connection", "communication"}

		aIsStorage := false
		bIsStorage := false
		aIsTransport := false
		bIsTransport := false

		for _, t := range storageTerms {
			if strings.Contains(aText, t) { aIsStorage = true }
			if strings.Contains(bText, t) { bIsStorage = true }
		}
		for _, t := range transportTerms {
			if strings.Contains(aText, t) { aIsTransport = true }
			if strings.Contains(bText, t) { bIsTransport = true }
		}

		if (aIsTransport && bIsStorage) || (aIsStorage && bIsTransport) {
			return true
		}
	}
	return false
}

func (se *SemanticEngine) buildContradiction(a, b Statement, pair SemanticPair) CIEContradiction {
	cat := pair.Category
	severity := pair.Severity
	confidence := pair.BaseConfidence

	desc := fmt.Sprintf("Statement A: %s\nStatement B: %s\n\n%s",
		a.OriginalText, b.OriginalText, pair.Reasoning)

	return CIEContradiction{
		ID:          fmt.Sprintf("CON-SEM-%03d", 0),
		Type:        cat,
		Severity:    severity,
		Confidence:  confidence,
		Summary:     fmt.Sprintf("%s: %s contradicts %s", cat, a.OriginalText, b.OriginalText),
		Description: desc,
		StatementA:  a,
		StatementB:  b,
		Reasoning:   pair.Reasoning + " " + pair.Impact,
		Recommendations: pair.Recommendations,
	}
}

func (se *SemanticEngine) deduplicate(contradictions []CIEContradiction) []CIEContradiction {
	seen := make(map[string]bool)
	var result []CIEContradiction

	for _, c := range contradictions {
		aText := strings.ToLower(strings.TrimSpace(c.StatementA.OriginalText))
		bText := strings.ToLower(strings.TrimSpace(c.StatementB.OriginalText))
		if aText == bText {
			continue
		}
		if aText > bText {
			aText, bText = bText, aText
		}
		key := string(c.Type) + "|" + aText + "|" + bText
		if seen[key] {
			continue
		}
		seen[key] = true
		result = append(result, c)
	}

	for i := range result {
		result[i].ID = fmt.Sprintf("CON-SEM-%03d", i+1)
	}
	return result
}
