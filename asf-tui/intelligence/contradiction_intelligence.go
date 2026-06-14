package intelligence

import (
	"fmt"
	"regexp"
	"strings"
)

// ─────────────────────────────────────────────────────────────
// PHASE 1 — CONTRADICTION DATA MODEL
// ─────────────────────────────────────────────────────────────

// ContradictionType categorizes the nature of the contradiction.
type ContradictionType string

const (
	ContradictionTypeAUTHENTICATION      ContradictionType = "AUTHENTICATION"
	ContradictionTypeAUTHORIZATION       ContradictionType = "AUTHORIZATION"
	ContradictionTypeENCRYPTION          ContradictionType = "ENCRYPTION"
	ContradictionTypeSECRETS             ContradictionType = "SECRETS"
	ContradictionTypeKEY_MANAGEMENT      ContradictionType = "KEY_MANAGEMENT"
	ContradictionTypeNETWORK             ContradictionType = "NETWORK"
	ContradictionTypeMONITORING          ContradictionType = "MONITORING"
	ContradictionTypeLOGGING             ContradictionType = "LOGGING"
	ContradictionTypeBACKUP              ContradictionType = "BACKUP"
	ContradictionTypeAVAILABILITY        ContradictionType = "AVAILABILITY"
	ContradictionTypeTHIRD_PARTY         ContradictionType = "THIRD_PARTY"
	ContradictionTypeCOMPLIANCE          ContradictionType = "COMPLIANCE"
	ContradictionTypeDATA_CLASSIFICATION ContradictionType = "DATA_CLASSIFICATION"
	ContradictionTypeIDENTITY            ContradictionType = "IDENTITY"
	ContradictionTypeGENERAL             ContradictionType = "GENERAL"
	ContradictionTypeTRUST                ContradictionType = "TRUST"
	ContradictionTypeTRUST_BOUNDARY      ContradictionType = "TRUST_BOUNDARY"
	ContradictionTypeCONTROL             ContradictionType = "CONTROL"
)

// CIEContradiction (Contradiction Intelligence Engine) is the rich contradiction model.
type CIEContradiction struct {
	ID                      string            `json:"id"`
	Type                    ContradictionType `json:"type"`
	Severity                RiskLevel         `json:"severity"`
	Confidence              float64           `json:"confidence"`
	Summary                 string            `json:"summary"`
	Description             string            `json:"description"`
	StatementA              Statement         `json:"statement_a"`
	StatementB              Statement         `json:"statement_b"`
	AffectedAssets          []string          `json:"affected_assets,omitempty"`
	AffectedComponents      []string          `json:"affected_components,omitempty"`
	AffectedControls        []string          `json:"affected_controls,omitempty"`
	AffectedTrustBoundaries []string          `json:"affected_trust_boundaries,omitempty"`
	Reasoning               string            `json:"reasoning"`
	Evidence                []string          `json:"evidence,omitempty"`
	Recommendations         []string          `json:"recommendations,omitempty"`
}

// Statement represents a normalized security claim.
type Statement struct {
	ID             string  `json:"id"`
	Source         string  `json:"source"`
	OriginalText   string  `json:"original_text"`
	NormalizedText string  `json:"normalized_text"`
	Category       string  `json:"category"`
	Subject        string  `json:"subject"`
	Predicate      string  `json:"predicate"`
	Object         string  `json:"object"`
	Confidence     float64 `json:"confidence"`
}

// ClaimExtractor extracts normalized claims from all architecture sources.
type ClaimExtractor struct{}

// NewClaimExtractor creates a claim extractor.
func NewClaimExtractor() *ClaimExtractor {
	return &ClaimExtractor{}
}

// ExtractClaims extracts all security claims from an architecture.
func (ce *ClaimExtractor) ExtractClaims(arch *ArchDescription, assumptions []Assumption) []Statement {
	var claims []Statement

	// Extract from assumptions
	for _, a := range assumptions {
		claims = append(claims, ce.normalizeAssumption(a))
	}

	// Extract from explicit assumptions
	for i, text := range arch.ExplicitAssumptions {
		claims = append(claims, Statement{
			ID:           fmt.Sprintf("EXP-%03d", i),
			Source:       "explicit_assumptions",
			OriginalText: text,
			Category:     "assumption",
			Confidence:   0.9,
		})
	}

	// Extract from security controls
	for category, controls := range arch.SecurityControls {
		for i, control := range controls {
			claims = append(claims, Statement{
				ID:           fmt.Sprintf("CTRL-%s-%03d", category, i),
				Source:       "security_controls",
				OriginalText: fmt.Sprintf("%s: %s", category, control),
				Category:     "control",
				Confidence:   0.85,
			})
		}
	}

	// Extract from policies
	for i, policy := range arch.Policies {
		claims = append(claims, Statement{
			ID:           fmt.Sprintf("POL-%03d", i),
			Source:       "policies",
			OriginalText: policy,
			Category:     "policy",
			Confidence:   0.8,
		})
	}

	// Extract from compliance
	for i, comp := range arch.Compliance {
		claims = append(claims, Statement{
			ID:           fmt.Sprintf("COMP-%03d", i),
			Source:       "compliance",
			OriginalText: comp,
			Category:     "compliance",
			Confidence:   0.9,
		})
	}

	// Extract from notes
	for i, note := range arch.Notes {
		claims = append(claims, Statement{
			ID:           fmt.Sprintf("NOTE-%03d", i),
			Source:       "notes",
			OriginalText: note,
			Category:     "note",
			Confidence:   0.7,
		})
	}

	// Extract from raw text
	claims = append(claims, ce.extractFromRawText(arch.RawText)...)

	return claims
}

func (ce *ClaimExtractor) normalizeAssumption(a Assumption) Statement {
	return Statement{
		ID:           a.ID,
		Source:       "assumption",
		OriginalText: a.Description,
		Category:     a.Category,
		Confidence:   a.Confidence,
	}
}

// extractFromRawText scans raw architecture text for security claims.
func (ce *ClaimExtractor) extractFromRawText(rawText string) []Statement {
	var claims []Statement
	if rawText == "" {
		return claims
	}

	// Define patterns that look like security claims
	patterns := []struct {
		re       *regexp.Regexp
		category string
	}{
		{regexp.MustCompile(`(?i)(all|every|each)\s+\w+\s+(is|are|must|should|shall)\s+(encrypted|secured|protected|authenticated|authorized|monitored|logged|backed\s+up|tested|reviewed)`), "configuration"},
		{regexp.MustCompile(`(?i)(no|none|never)\s+\w+\s+(can|may|has|have|is|are)\s+(access|permission|exposed|public|external)`), "access"},
		{regexp.MustCompile(`(?i)(only|just|exclusively)\s+(authorized|authenticated|admin|privileged)\s+\w+\s+(can|may|has|have)`), "access"},
		{regexp.MustCompile(`(?i)(mfa|multi.?factor|two.?factor|2fa)\s+(is|are|must|should|shall|required|enforced|mandatory)`), "identity"},
		{regexp.MustCompile(`(?i)(tls|ssl|https)\s+(is|are|must|should|shall|required|enforced|mandatory)`), "network"},
		{regexp.MustCompile(`(?i)(private|internal|isolated|segmented)\s+(subnet|network|vpc|zone)`), "network"},
		{regexp.MustCompile(`(?i)(public|external|internet|exposed)\s+(accessible|access|facing|endpoint)`), "network"},
		{regexp.MustCompile(`(?i)(audit|log|monitor|alert)\s+(is|are|must|should|shall|required|enforced|immutable|tamper.?proof)`), "logging"},
		{regexp.MustCompile(`(?i)(backup|restore|recovery)\s+(is|are|must|should|shall|tested|encrypted|verified|regular)`), "backup"},
		{regexp.MustCompile(`(?i)(least\s+privilege|minimum\s+necessary|need.?to.?know)\s+(is|are|must|should|shall|enforced|applied)`), "authorization"},
		{regexp.MustCompile(`(?i)(rbac|abac|role.?based|access.?control)\s+(is|are|must|should|shall|enforced|implemented)`), "authorization"},
		{regexp.MustCompile(`(?i)(shared|generic|common)\s+(admin|account|credential|password|token)`), "identity"},
		{regexp.MustCompile(`(?i)(service\s+account|system\s+account|bot\s+account)\s+(exempt|bypass|skip|ignore|not\s+required)`), "identity"},
	}

	lines := strings.Split(rawText, "\n")
	for i, line := range lines {
		for _, pat := range patterns {
			if pat.re.MatchString(line) {
				claims = append(claims, Statement{
					ID:           fmt.Sprintf("RAW-%04d", i),
					Source:       "raw_text",
					OriginalText: strings.TrimSpace(line),
					Category:     pat.category,
					Confidence:   0.75,
				})
				break // only one claim per line
			}
		}
	}
	return claims
}

// ─────────────────────────────────────────────────────────────
// PHASE 3 — CONTRADICTION DETECTION RULES
// ─────────────────────────────────────────────────────────────

// CIEEngine (Contradiction Intelligence Engine) is the main contradiction engine.
type CIEEngine struct {
	claimExtractor *ClaimExtractor
	semanticEngine *SemanticEngine
}

// NewCIEEngine creates a new contradiction intelligence engine.
func NewCIEEngine() *CIEEngine {
	return &CIEEngine{
		claimExtractor: NewClaimExtractor(),
		semanticEngine: NewSemanticEngine(),
	}
}

// DetectAllContradictions runs all contradiction detection phases.
func (cie *CIEEngine) DetectAllContradictions(arch *ArchDescription, assumptions []Assumption, controls []ControlDetail, boundaries []TrustBoundary) []CIEContradiction {
	claims := cie.claimExtractor.ExtractClaims(arch, assumptions)

	var contradictions []CIEContradiction

	// Phase 3 — Explicit contradictions
	contradictions = append(contradictions, cie.detectAuthenticationContradictions(claims, assumptions)...)
	contradictions = append(contradictions, cie.detectAuthorizationContradictions(claims, assumptions)...)
	contradictions = append(contradictions, cie.detectEncryptionContradictions(claims, assumptions)...)
	contradictions = append(contradictions, cie.detectSecretsContradictions(claims, assumptions)...)
	contradictions = append(contradictions, cie.detectKeyManagementContradictions(claims, assumptions)...)
	contradictions = append(contradictions, cie.detectBackupContradictions(claims, assumptions)...)
	contradictions = append(contradictions, cie.detectMonitoringContradictions(claims, assumptions)...)
	contradictions = append(contradictions, cie.detectComplianceContradictions(claims, assumptions, arch)...)
	contradictions = append(contradictions, cie.detectNetworkContradictions(claims, assumptions)...)

	// Phase 4 — Implied contradictions
	contradictions = append(contradictions, cie.detectImpliedContradictions(assumptions)...)

	// Phase 5 — Trust boundary contradictions
	contradictions = append(contradictions, cie.detectTrustBoundaryContradictions(assumptions, boundaries)...)

	// Phase 6 — Control contradictions
	contradictions = append(contradictions, cie.detectControlContradictions(assumptions, controls)...)

	// Phase 7 — Compliance contradictions
	contradictions = append(contradictions, cie.detectComplianceFrameworkContradictions(assumptions, arch)...)

	// Phase 8 — Score all contradictions
	for i := range contradictions {
		contradictions[i] = cie.scoreContradiction(contradictions[i], assumptions, arch)
	}

	// Phase 9 — Semantic contradiction detection
	semantic := cie.semanticEngine.DetectSemanticContradictions(claims)
	for i := range semantic {
		semantic[i] = cie.scoreContradiction(semantic[i], assumptions, arch)
	}
	contradictions = append(contradictions, semantic...)

	return contradictions
}

// ─────────────────────────────────────────────────────────────
// AUTHENTICATION CONTRADICTIONS
// ─────────────────────────────────────────────────────────────

func (cie *CIEEngine) detectAuthenticationContradictions(claims []Statement, assumptions []Assumption) []CIEContradiction {
	var results []CIEContradiction

	// MFA required vs MFA exempt
	mfaRequired := cie.findClaims(claims, []string{"mfa", "multi-factor", "two-factor", "2fa"}, []string{"required", "enforced", "mandatory", "all"})
	mfaExempt := cie.findClaims(claims, []string{"mfa", "multi-factor", "two-factor", "2fa"}, []string{"exempt", "bypass", "skip", "not required", "optional"})

	if len(mfaRequired) > 0 && len(mfaExempt) > 0 {
		for _, req := range mfaRequired {
			for _, ex := range mfaExempt {
				if req.ID == ex.ID {
					continue
				}
				results = append(results, CIEContradiction{
					ID:          fmt.Sprintf("CON-AUTH-%03d", len(results)),
					Type:        ContradictionTypeAUTHENTICATION,
					Severity:    RiskHigh,
					Confidence:  0.92,
					Summary:     "MFA required for all users but some accounts are exempt",
					Description: fmt.Sprintf("Claim '%s' requires MFA, but '%s' exempts accounts from MFA", req.OriginalText, ex.OriginalText),
					StatementA:  req,
					StatementB:  ex,
					Reasoning:   "MFA is a strong authentication control. Exempting any account class creates a bypass path that undermines the entire authentication posture.",
					Recommendations: []string{
						"Review all MFA exemptions and justify each with risk assessment",
						"Implement compensating controls for exempted accounts (e.g., IP restriction, certificate-based auth)",
						"Require MFA for service accounts or use workload identity",
					},
				})
			}
		}
	}

	// Admin MFA required vs break-glass exempt
	adminMFA := cie.findClaims(claims, []string{"admin", "administrator", "root"}, []string{"mfa", "multi-factor", "two-factor", "2fa"})
	breakglass := cie.findClaims(claims, []string{"break.glass", "emergency", "firecall"}, []string{"exempt", "bypass", "skip", "no mfa", "without mfa"})

	if len(adminMFA) > 0 && len(breakglass) > 0 {
		for _, am := range adminMFA {
			for _, bg := range breakglass {
				if am.ID == bg.ID {
					continue
				}
				results = append(results, CIEContradiction{
					ID:          fmt.Sprintf("CON-AUTH-%03d", len(results)),
					Type:        ContradictionTypeAUTHENTICATION,
					Severity:    RiskHigh,
					Confidence:  0.88,
					Summary:     "Administrative MFA required but break-glass accounts bypass it",
					Description: fmt.Sprintf("'%s' requires admin MFA, but '%s' allows bypass", am.OriginalText, bg.OriginalText),
					StatementA:  am,
					StatementB:  bg,
					Reasoning:   "Break-glass accounts are highly privileged and often targeted. Exempting them from MFA creates a critical attack path.",
					Recommendations: []string{
						"Require MFA for all break-glass access with hardware tokens",
						"Store break-glass credentials in a physical safe with dual-control",
						"Monitor break-glass usage with real-time alerting",
					},
				})
			}
		}
	}

	// All users authenticated vs anonymous access allowed
	allAuth := cie.findClaims(claims, []string{"all", "every", "each"}, []string{"authenticated", "login", "auth", "identity"})
	anonAccess := cie.findClaims(claims, []string{"anonymous", "guest", "public", "unauthenticated", "no login"}, []string{"access", "allowed", "permitted", "can"})

	if len(allAuth) > 0 && len(anonAccess) > 0 {
		for _, aa := range allAuth {
			for _, an := range anonAccess {
				if aa.ID == an.ID {
					continue
				}
				results = append(results, CIEContradiction{
					ID:          fmt.Sprintf("CON-AUTH-%03d", len(results)),
					Type:        ContradictionTypeAUTHENTICATION,
					Severity:    RiskHigh,
					Confidence:  0.85,
					Summary:     "All users must be authenticated but anonymous access is allowed",
					Description: fmt.Sprintf("'%s' requires authentication, but '%s' allows anonymous access", aa.OriginalText, an.OriginalText),
					StatementA:  aa,
					StatementB:  an,
					Reasoning:   "Anonymous access contradicts any claim that all users are authenticated. This creates an unauthenticated entry point.",
					Recommendations: []string{
						"Require authentication for all endpoints",
						"If public access is needed, segregate it to a separate service with no sensitive data",
					},
				})
			}
		}
	}

	return results
}

// ─────────────────────────────────────────────────────────────
// AUTHORIZATION CONTRADICTIONS
// ─────────────────────────────────────────────────────────────

func (cie *CIEEngine) detectAuthorizationContradictions(claims []Statement, assumptions []Assumption) []CIEContradiction {
	var results []CIEContradiction

	// Least privilege vs shared admin
	leastPriv := cie.findClaims(claims, []string{"least privilege", "least-privilege", "minimum necessary", "need to know"}, nil)
	sharedAdmin := cie.findClaims(claims, []string{"shared admin", "shared account", "generic admin", "common admin", "shared credential", "share a single", "share one", "share the same"}, nil)

	if len(leastPriv) > 0 && len(sharedAdmin) > 0 {
		for _, lp := range leastPriv {
			for _, sa := range sharedAdmin {
				if lp.ID == sa.ID {
					continue
				}
				results = append(results, CIEContradiction{
					ID:          fmt.Sprintf("CON-AUTHZ-%03d", len(results)),
					Type:        ContradictionTypeAUTHORIZATION,
					Severity:    RiskCritical,
					Confidence:  0.95,
					Summary:     "Least privilege claimed but shared admin accounts exist",
					Description: fmt.Sprintf("'%s' claims least privilege, but '%s' uses shared admin", lp.OriginalText, sa.OriginalText),
					StatementA:  lp,
					StatementB:  sa,
					Reasoning:   "Shared admin accounts violate least privilege because accountability is lost and any compromise affects all users of that account.",
					Recommendations: []string{
						"Eliminate shared admin accounts; use individual accounts with named credentials",
						"Implement privileged access management (PAM) for shared-role scenarios",
						"Enforce just-in-time (JIT) access for administrative tasks",
					},
				})
			}
		}
	}

	// RBAC enforced vs everyone has admin
	rbac := cie.findClaims(claims, []string{"rbac", "role-based", "abac", "attribute-based"}, []string{"enforced", "implemented", "required", "active"})
	allAdmin := cie.findClaims(claims, []string{"everyone", "all", "all users", "all staff", "all employees"}, []string{"admin", "administrator", "root", "privileged", "full access"})

	if len(rbac) > 0 && len(allAdmin) > 0 {
		for _, r := range rbac {
			for _, aa := range allAdmin {
				if r.ID == aa.ID {
					continue
				}
				results = append(results, CIEContradiction{
					ID:          fmt.Sprintf("CON-AUTHZ-%03d", len(results)),
					Type:        ContradictionTypeAUTHORIZATION,
					Severity:    RiskHigh,
					Confidence:  0.90,
					Summary:     "RBAC enforced but all users have admin access",
					Description: fmt.Sprintf("'%s' enforces RBAC, but '%s' grants admin to everyone", r.OriginalText, aa.OriginalText),
					StatementA:  r,
					StatementB:  aa,
					Reasoning:   "RBAC is meaningless if all users hold the highest-privilege role. The access control matrix collapses to a single column.",
					Recommendations: []string{
						"Implement role hierarchy with standard, elevated, and admin roles",
						"Enforce regular access recertification",
						"Apply segregation of duties (SoD) rules",
					},
				})
			}
		}
	}

	// Object-level authorization vs no authorization check
	objAuth := cie.findClaims(claims, []string{"object-level", "object level", "object authorization", "resource-level"}, []string{"enforced", "implemented", "required", "active"})
	noAuthCheck := cie.findClaims(claims, []string{"no authorization", "no access check", "no permission check", "bypass authorization"}, nil)

	if len(objAuth) > 0 && len(noAuthCheck) > 0 {
		for _, oa := range objAuth {
			for _, nac := range noAuthCheck {
				if oa.ID == nac.ID {
					continue
				}
				results = append(results, CIEContradiction{
					ID:          fmt.Sprintf("CON-AUTHZ-%03d", len(results)),
					Type:        ContradictionTypeAUTHORIZATION,
					Severity:    RiskHigh,
					Confidence:  0.87,
					Summary:     "Object-level authorization claimed but some paths bypass it",
					Description: fmt.Sprintf("'%s' enforces object-level auth, but '%s' bypasses it", oa.OriginalText, nac.OriginalText),
					StatementA:  oa,
					StatementB:  nac,
					Reasoning:   "Object-level authorization is only effective if applied to every access path. A single bypass defeats the control.",
					Recommendations: []string{
						"Enforce authorization at every API layer (gateway, service, database)",
						"Use a centralized authorization policy engine (e.g., OPA, Cedar)",
					},
				})
			}
		}
	}

	return results
}

// ─────────────────────────────────────────────────────────────
// ENCRYPTION CONTRADICTIONS
// ─────────────────────────────────────────────────────────────

func (cie *CIEEngine) detectEncryptionContradictions(claims []Statement, assumptions []Assumption) []CIEContradiction {
	var results []CIEContradiction

	// All traffic encrypted vs HTTP/legacy allowed
	// Exclude storage/backup-context statements (e.g., "backups in plaintext") from transport-layer checks
	storageExclude := []string{"backup", "storage", "at rest", "disk", "store", "database"}
	allEncrypted := cie.findClaimsFiltered(claims, []string{"all", "every", "all traffic", "all communication", "all data"}, []string{"encrypted", "tls", "ssl", "https", "secure"}, storageExclude)
	httpAllowed := cie.findClaimsFiltered(claims, []string{"http", "unencrypted", "plaintext", "clear text", "http allowed", "port 80"}, nil, storageExclude)

	if len(allEncrypted) > 0 && len(httpAllowed) > 0 {
		for _, ae := range allEncrypted {
			for _, ha := range httpAllowed {
				if ae.ID == ha.ID {
					continue
				}
				results = append(results, CIEContradiction{
					ID:          fmt.Sprintf("CON-ENC-%03d", len(results)),
					Type:        ContradictionTypeENCRYPTION,
					Severity:    RiskHigh,
					Confidence:  0.90,
					Summary:     "All traffic encrypted but HTTP/unencrypted allowed",
					Description: fmt.Sprintf("'%s' requires encryption, but '%s' allows unencrypted traffic", ae.OriginalText, ha.OriginalText),
					StatementA:  ae,
					StatementB:  ha,
					Reasoning:   "Allowing HTTP alongside TLS creates a downgrade attack path. Attackers can force connections to HTTP to intercept data.",
					Recommendations: []string{
						"Enforce HTTPS-only with HSTS headers",
						"Block HTTP at the load balancer/WAF level",
						"Use TLS 1.3 with certificate pinning",
					},
				})
			}
		}
	}

	// TLS required vs TLS optional
	tlsReq := cie.findClaimsFiltered(claims, []string{"tls", "ssl", "https"}, []string{"required", "mandatory", "enforced", "must", "shall"}, storageExclude)
	tlsOpt := cie.findClaimsFiltered(claims, []string{"tls", "ssl", "https"}, []string{"optional", "not required", "not enforced", "not mandatory", "if available", "when possible"}, storageExclude)

	if len(tlsReq) > 0 && len(tlsOpt) > 0 {
		for _, tr := range tlsReq {
			for _, to := range tlsOpt {
				if tr.ID == to.ID {
					continue
				}
				results = append(results, CIEContradiction{
					ID:          fmt.Sprintf("CON-ENC-%03d", len(results)),
					Type:        ContradictionTypeENCRYPTION,
					Severity:    RiskHigh,
					Confidence:  0.88,
					Summary:     "TLS required but also optional depending on client",
					Description: fmt.Sprintf("'%s' requires TLS, but '%s' makes it optional", tr.OriginalText, to.OriginalText),
					StatementA:  tr,
					StatementB:  to,
					Reasoning:   "Optional TLS is equivalent to no TLS. An attacker can simply refuse TLS negotiation.",
					Recommendations: []string{
						"Enforce TLS for all connections (reject plaintext)",
						"Use TLS 1.3 with strict cipher suites",
					},
				})
			}
		}
	}

	return results
}

// ─────────────────────────────────────────────────────────────
// SECRETS CONTRADICTIONS
// ─────────────────────────────────────────────────────────────

func (cie *CIEEngine) detectSecretsContradictions(claims []Statement, assumptions []Assumption) []CIEContradiction {
	var results []CIEContradiction

	// Secrets in vault vs secrets in source code
	vaultSecrets := cie.findClaims(claims, []string{"vault", "secrets manager", "secrets management", "hashicorp", "aws secrets", "azure key vault"}, []string{"secrets", "credentials", "keys", "tokens", "passwords"})
	sourceSecrets := cie.findClaims(claims, []string{"source code", "repository", "git", "config file", "environment variable", "env var", "hardcoded", "embedded"}, []string{"secret", "credential", "password", "token", "api key", "key"})

	if len(vaultSecrets) > 0 && len(sourceSecrets) > 0 {
		for _, vs := range vaultSecrets {
			for _, ss := range sourceSecrets {
				if vs.ID == ss.ID {
					continue
				}
				results = append(results, CIEContradiction{
					ID:          fmt.Sprintf("CON-SEC-%03d", len(results)),
					Type:        ContradictionTypeSECRETS,
					Severity:    RiskCritical,
					Confidence:  0.93,
					Summary:     "Secrets managed in vault but also stored in source code",
					Description: fmt.Sprintf("'%s' uses a vault, but '%s' stores secrets in source", vs.OriginalText, ss.OriginalText),
					StatementA:  vs,
					StatementB:  ss,
					Reasoning:   "Storing secrets in source code defeats the purpose of a vault. Secrets in code are visible to all developers and leak through version control.",
					Recommendations: []string{
						"Remove all secrets from source code and configuration files",
						"Use secret injection at runtime (e.g., Kubernetes secrets, AWS Parameter Store)",
						"Scan repositories for secrets with tools like GitLeaks, TruffleHog",
					},
				})
			}
		}
	}

	return results
}

// ─────────────────────────────────────────────────────────────
// KEY MANAGEMENT CONTRADICTIONS
// ─────────────────────────────────────────────────────────────

func (cie *CIEEngine) detectKeyManagementContradictions(claims []Statement, assumptions []Assumption) []CIEContradiction {
	var results []CIEContradiction

	// Keys rotated vs keys never rotated
	keyRotated := cie.findClaims(claims, []string{"key", "keys", "encryption key", "signing key"}, []string{"rotated", "rotation", "periodically", "regularly", "scheduled", "automatic"})
	keyNeverRotated := cie.findClaims(claims, []string{"key", "keys", "encryption key", "signing key"}, []string{"never rotated", "no rotation", "static", "fixed", "hardcoded key", "never changed"})

	if len(keyRotated) > 0 && len(keyNeverRotated) > 0 {
		for _, kr := range keyRotated {
			for _, knr := range keyNeverRotated {
				if kr.ID == knr.ID {
					continue
				}
				results = append(results, CIEContradiction{
					ID:          fmt.Sprintf("CON-KM-%03d", len(results)),
					Type:        ContradictionTypeKEY_MANAGEMENT,
					Severity:    RiskHigh,
					Confidence:  0.91,
					Summary:     "Key rotation required but some keys are never rotated",
					Description: fmt.Sprintf("'%s' requires rotation, but '%s' has static keys", kr.OriginalText, knr.OriginalText),
					StatementA:  kr,
					StatementB:  knr,
					Reasoning:   "Static keys accumulate exposure over time. If a key is ever compromised, the window of exposure is indefinite.",
					Recommendations: []string{
						"Implement automatic key rotation (e.g., AWS KMS automatic rotation)",
						"Set maximum key age policy (e.g., 90 days for signing keys)",
						"Track key age and alert on expiration",
					},
				})
			}
		}
	}

	return results
}

// ─────────────────────────────────────────────────────────────
// BACKUP CONTRADICTIONS
// ─────────────────────────────────────────────────────────────

func (cie *CIEEngine) detectBackupContradictions(claims []Statement, assumptions []Assumption) []CIEContradiction {
	var results []CIEContradiction

	// Backups tested vs restore process unknown
	backupTested := cie.findClaims(claims, []string{"backup", "backups"}, []string{"tested", "verified", "validated", "regularly tested", "periodically tested"})
	restoreUnknown := cie.findClaims(claims, []string{"restore", "recovery", "restoration"}, []string{"unknown", "not tested", "untested", "not documented", "not defined", "no procedure"})

	if len(backupTested) > 0 && len(restoreUnknown) > 0 {
		for _, bt := range backupTested {
			for _, ru := range restoreUnknown {
				if bt.ID == ru.ID {
					continue
				}
				results = append(results, CIEContradiction{
					ID:          fmt.Sprintf("CON-BAK-%03d", len(results)),
					Type:        ContradictionTypeBACKUP,
					Severity:    RiskMedium,
					Confidence:  0.82,
					Summary:     "Backups tested but restore process is unknown",
					Description: fmt.Sprintf("'%s' tests backups, but '%s' has no restore process", bt.OriginalText, ru.OriginalText),
					StatementA:  bt,
					StatementB:  ru,
					Reasoning:   "A backup that cannot be restored is worthless. The restore process is the actual security control.",
					Recommendations: []string{
						"Document and test the full restore procedure quarterly",
						"Perform disaster recovery drills with realistic scenarios",
						"Measure and track RTO/RPO targets",
					},
				})
			}
		}
	}

	// Encrypted data vs plaintext backups
	encryptedData := cie.findClaims(claims, []string{"all", "every", "data", "everything", "all data"}, []string{"encrypted", "encryption", "encrypt"})
	plaintextBackup := cie.findClaims(claims, []string{"backup", "backups"}, []string{"plaintext", "unencrypted", "not encrypted", "no encryption", "clear text"})

	if len(encryptedData) > 0 && len(plaintextBackup) > 0 {
		for _, ed := range encryptedData {
			for _, pb := range plaintextBackup {
				if ed.ID == pb.ID {
					continue
				}
				results = append(results, CIEContradiction{
					ID:          fmt.Sprintf("CON-BAK-%03d", len(results)),
					Type:        ContradictionTypeBACKUP,
					Severity:    RiskCritical,
					Confidence:  0.92,
					Summary:     "Data is encrypted but backups are stored in plaintext",
					Description: fmt.Sprintf("'%s' requires encryption, but '%s' stores backups in plaintext", ed.OriginalText, pb.OriginalText),
					StatementA:  ed,
					StatementB:  pb,
					Reasoning:   "Plaintext backups bypass encryption controls. If the primary data is encrypted but backups are not, the protection is incomplete and the backup becomes an easy target.",
					Recommendations: []string{
						"Encrypt all backups using the same or stronger encryption than primary storage",
						"Implement key management for backup encryption keys",
						"Verify backup encryption with periodic restore and validation tests",
					},
				})
			}
		}
	}

	return results
}

// ─────────────────────────────────────────────────────────────
// MONITORING CONTRADICTIONS
// ─────────────────────────────────────────────────────────────

func (cie *CIEEngine) detectMonitoringContradictions(claims []Statement, assumptions []Assumption) []CIEContradiction {
	var results []CIEContradiction

	// All logs monitored vs alerts not reviewed
	logsMonitored := cie.findClaims(claims, []string{"log", "logs", "monitoring", "alert", "detection"}, []string{"monitored", "reviewed", "analyzed", "checked", "all", "every"})
	alertsNotReviewed := cie.findClaims(claims, []string{"alert", "alerting", "notification", "alarm"}, []string{"not reviewed", "not checked", "ignored", "not acted", "backlog", "not monitored", "not responded"})

	if len(logsMonitored) > 0 && len(alertsNotReviewed) > 0 {
		for _, lm := range logsMonitored {
			for _, anr := range alertsNotReviewed {
				if lm.ID == anr.ID {
					continue
				}
				results = append(results, CIEContradiction{
					ID:          fmt.Sprintf("CON-MON-%03d", len(results)),
					Type:        ContradictionTypeMONITORING,
					Severity:    RiskMedium,
					Confidence:  0.80,
					Summary:     "All logs monitored but alerts are not reviewed",
					Description: fmt.Sprintf("'%s' monitors logs, but '%s' ignores alerts", lm.OriginalText, anr.OriginalText),
					StatementA:  lm,
					StatementB:  anr,
					Reasoning:   "Monitoring without response is theater. If alerts are not reviewed, the monitoring system provides no actual security value.",
					Recommendations: []string{
						"Establish SOC procedures with defined SLA for alert response",
						"Implement automated response for high-confidence alerts",
						"Track alert resolution metrics and MTTR",
					},
				})
			}
		}
	}

	return results
}

// ─────────────────────────────────────────────────────────────
// COMPLIANCE CONTRADICTIONS
// ─────────────────────────────────────────────────────────────

func (cie *CIEEngine) detectComplianceContradictions(claims []Statement, assumptions []Assumption, arch *ArchDescription) []CIEContradiction {
	var results []CIEContradiction

	// HIPAA without audit logging
	if cie.hasCompliance(arch, "HIPAA") {
		if !cie.hasClaim(claims, []string{"audit", "audit log", "audit logging", "audit trail"}, []string{"required", "enforced", "implemented", "active", "present"}) {
			results = append(results, CIEContradiction{
				ID:          fmt.Sprintf("CON-COMP-%03d", len(results)),
				Type:        ContradictionTypeCOMPLIANCE,
				Severity:    RiskHigh,
				Confidence:  0.90,
				Summary:     "HIPAA compliance claimed but audit logging is not specified",
				Description: "HIPAA requires audit logging (§164.312(b)), but the architecture does not specify audit logging controls.",
				Reasoning:   "HIPAA Security Rule requires audit controls (§164.312(b)) to record and examine access and activity. Without audit logging, HIPAA compliance cannot be demonstrated.",
				Recommendations: []string{
					"Implement comprehensive audit logging for all PHI access",
					"Log user identity, timestamp, action, and data accessed",
					"Ensure logs are immutable and tamper-evident",
				},
			})
		}
	}

	// SOC2 without access management
	if cie.hasCompliance(arch, "SOC2") {
		if !cie.hasClaim(claims, []string{"access", "access control", "access management", "authorization"}, []string{"enforced", "implemented", "managed", "reviewed", "controlled"}) {
			results = append(results, CIEContradiction{
				ID:          fmt.Sprintf("CON-COMP-%03d", len(results)),
				Type:        ContradictionTypeCOMPLIANCE,
				Severity:    RiskHigh,
				Confidence:  0.85,
				Summary:     "SOC2 compliance claimed but access management is not specified",
				Description: "SOC2 requires access management (CC6.1-CC6.3), but the architecture does not specify access controls.",
				Reasoning:   "SOC2 Trust Services Criteria require logical access controls, access removal, and access review. These must be documented and implemented.",
				Recommendations: []string{
					"Implement RBAC with least privilege",
					"Enforce access reviews and recertification",
					"Automate access removal upon termination",
				},
			})
		}
	}

	// PCI DSS without encryption
	if cie.hasCompliance(arch, "PCI", "PCI DSS") {
		if !cie.hasClaim(claims, []string{"encrypt", "encryption", "tls", "ssl", "cryptographic"}, []string{"required", "enforced", "implemented", "active", "present"}) {
			results = append(results, CIEContradiction{
				ID:          fmt.Sprintf("CON-COMP-%03d", len(results)),
				Type:        ContradictionTypeCOMPLIANCE,
				Severity:    RiskCritical,
				Confidence:  0.95,
				Summary:     "PCI DSS compliance claimed but encryption is not specified",
				Description: "PCI DSS requires encryption (Req 3.4, 4.1) for cardholder data, but the architecture does not specify encryption.",
				Reasoning:   "PCI DSS Requirement 3.4 mandates rendering PAN unreadable, and Req 4.1 requires encrypting transmission. These are mandatory, not optional.",
				Recommendations: []string{
					"Encrypt cardholder data at rest with AES-256",
					"Use TLS 1.3 for all transmission of cardholder data",
					"Implement key management per PCI DSS Req 3.6",
				},
			})
		}
	}

	return results
}

// ─────────────────────────────────────────────────────────────
// NETWORK CONTRADICTIONS
// ─────────────────────────────────────────────────────────────

func (cie *CIEEngine) detectNetworkContradictions(claims []Statement, assumptions []Assumption) []CIEContradiction {
	var results []CIEContradiction

	// Private network vs public access
	privateNet := cie.findClaims(claims, []string{"private", "internal", "isolated", "segmented", "no public", "not exposed"}, []string{"network", "subnet", "vpc", "access", "only", "internal"})
	publicAccess := cie.findClaims(claims, []string{"public", "external", "internet", "exposed", "publicly accessible"}, []string{"access", "reachable", "available", "endpoint", "api", "port"})

	if len(privateNet) > 0 && len(publicAccess) > 0 {
		for _, pn := range privateNet {
			for _, pa := range publicAccess {
				if pn.ID == pa.ID {
					continue
				}
				results = append(results, CIEContradiction{
					ID:          fmt.Sprintf("CON-NET-%03d", len(results)),
					Type:        ContradictionTypeNETWORK,
					Severity:    RiskHigh,
					Confidence:  0.88,
					Summary:     "Private network claimed but public access is allowed",
					Description: fmt.Sprintf("'%s' claims private network, but '%s' allows public access", pn.OriginalText, pa.OriginalText),
					StatementA:  pn,
					StatementB:  pa,
					Reasoning:   "A system cannot be both private and publicly accessible. Public exposure contradicts the network isolation claim.",
					Recommendations: []string{
						"Remove public access or move to a DMZ",
						"Implement network segmentation with explicit firewall rules",
						"Use VPN or private connectivity for external access",
					},
				})
			}
		}
	}

	return results
}

// ─────────────────────────────────────────────────────────────
// PHASE 4 — IMPLIED CONTRADICTIONS
// ─────────────────────────────────────────────────────────────

func (cie *CIEEngine) detectImpliedContradictions(assumptions []Assumption) []CIEContradiction {
	var results []CIEContradiction

	// PHI access restricted + Database publicly accessible
	hasPHI := cie.hasAssumption(assumptions, []string{"phi", "patient", "health", "medical", "ehr", "hipaa"}, nil)
	dbPublic := cie.hasAssumption(assumptions, []string{"database", "db", "storage"}, []string{"public", "internet", "external", "accessible", "exposed"})

	if hasPHI && dbPublic {
		results = append(results, CIEContradiction{
			ID:          fmt.Sprintf("CON-IMP-%03d", len(results)),
			Type:        ContradictionTypeDATA_CLASSIFICATION,
			Severity:    RiskCritical,
			Confidence:  0.85,
			Summary:     "PHI data present but database is publicly accessible",
			Description: "Architecture contains PHI data but the database is accessible from public networks. This creates an implied contradiction between data sensitivity and network exposure.",
			Reasoning:   "PHI data requires strict access controls (HIPAA §164.312). Public database access violates the minimum necessary standard and creates a direct breach path.",
			Recommendations: []string{
				"Move PHI database to private subnet with no public access",
				"Implement API gateway with authentication for all data access",
				"Use VPC endpoints and private links for all database connections",
			},
		})
	}

	// Least privilege + Admin access granted to all developers
	hasLeastPrivilege := cie.hasAssumption(assumptions, []string{"least privilege", "least-privilege", "minimum necessary"}, nil)
	allDevsAdmin := cie.hasAssumption(assumptions, []string{"developer", "dev", "engineer", "team"}, []string{"admin", "administrator", "root", "privileged", "full access"})

	if hasLeastPrivilege && allDevsAdmin {
		results = append(results, CIEContradiction{
			ID:          fmt.Sprintf("CON-IMP-%03d", len(results)),
			Type:        ContradictionTypeAUTHORIZATION,
			Severity:    RiskHigh,
			Confidence:  0.82,
			Summary:     "Least privilege claimed but all developers have admin access",
			Description: "Architecture claims least privilege but grants admin access to all developers. This is an implied contradiction.",
			Reasoning:   "Least privilege requires that users have only the minimum access necessary. Granting admin to all developers violates this principle and creates broad attack surface.",
			Recommendations: []string{
				"Implement role-based access with developer, senior, and admin roles",
				"Use just-in-time (JIT) elevation for admin tasks",
				"Enforce quarterly access reviews",
			},
		})
	}

	// Encryption at rest + No key management
	hasEncryption := cie.hasAssumption(assumptions, []string{"encrypted", "encryption", "encrypt", "aes", "cipher"}, []string{"rest", "at rest", "storage", "database"})
	noKeyMgmt := !cie.hasAssumption(assumptions, []string{"key", "kms", "hsm", "vault", "key management", "key rotation"}, []string{"management", "rotation", "storage", "policy", "escrow"})

	if hasEncryption && noKeyMgmt {
		results = append(results, CIEContradiction{
			ID:          fmt.Sprintf("CON-IMP-%03d", len(results)),
			Type:        ContradictionTypeKEY_MANAGEMENT,
			Severity:    RiskHigh,
			Confidence:  0.80,
			Summary:     "Encryption at rest present but key management is not specified",
			Description: "Architecture specifies encryption at rest but does not document key management. This is an implied contradiction because encryption without key management is incomplete.",
			Reasoning:   "Encryption is only as strong as its key management. Without key generation, storage, rotation, and access controls, encrypted data is at risk.",
			Recommendations: []string{
				"Document key management lifecycle (generation, distribution, rotation, destruction)",
				"Use a centralized key management service (e.g., AWS KMS, HashiCorp Vault)",
				"Implement key access audit logging",
			},
		})
	}

	// Session management + No token rotation
	hasSession := cie.hasAssumption(assumptions, []string{"session", "token", "jwt", "cookie", "sso"}, []string{"management", "required", "enforced", "active"})
	noRotation := !cie.hasAssumption(assumptions, []string{"rotation", "refresh", "reissue", "expire", "timeout", "revoke"}, []string{"token", "session", "key"})

	if hasSession && noRotation {
		results = append(results, CIEContradiction{
			ID:          fmt.Sprintf("CON-IMP-%03d", len(results)),
			Type:        ContradictionTypeIDENTITY,
			Severity:    RiskMedium,
			Confidence:  0.75,
			Summary:     "Session management present but token rotation not specified",
			Description: "Architecture has session management but does not specify token rotation or expiration. This is an implied contradiction.",
			Reasoning:   "Session tokens that never rotate or expire create a persistent attack surface. A stolen token can be used indefinitely.",
			Recommendations: []string{
				"Implement token rotation with refresh tokens",
				"Set token expiration (e.g., 15 minutes for access tokens, 7 days for refresh tokens)",
				"Enable token revocation on logout or security event",
			},
		})
	}

	return results
}

// ─────────────────────────────────────────────────────────────
// PHASE 5 — TRUST BOUNDARY CONTRADICTIONS
// ─────────────────────────────────────────────────────────────

func (cie *CIEEngine) detectTrustBoundaryContradictions(assumptions []Assumption, boundaries []TrustBoundary) []CIEContradiction {
	var results []CIEContradiction

	for _, boundary := range boundaries {
		switch boundary.Type {
		case "Internet":
			// Internet boundary + PHI access
			if cie.hasAssumption(assumptions, []string{"phi", "patient", "health", "medical", "ehr", "hipaa"}, nil) {
				results = append(results, CIEContradiction{
					ID:                      fmt.Sprintf("CON-TB-%03d", len(results)),
					Type:                    ContradictionTypeTRUST_BOUNDARY,
					Severity:                RiskCritical,
					Confidence:              0.90,
					Summary:                 "Internet trust boundary crosses PHI data path",
					Description:             fmt.Sprintf("Internet trust boundary at %s exposes PHI data to external networks. This contradicts HIPAA requirements for controlled access.", strings.Join(boundary.Components, ", ")),
					Reasoning:               "HIPAA requires that PHI access is controlled and logged. Internet exposure creates a direct breach path that violates the minimum necessary standard.",
					AffectedTrustBoundaries: []string{boundary.Type},
					Recommendations: []string{
						"Move PHI components behind a VPN or private endpoint",
						"Implement API gateway with strong authentication and rate limiting",
						"Encrypt all PHI in transit with TLS 1.3",
					},
				})
			}

		case "Vendor":
			// Vendor boundary + internal-only claim
			if cie.hasAssumption(assumptions, []string{"internal only", "internal-only", "no external", "no third-party", "no vendor"}, nil) {
				results = append(results, CIEContradiction{
					ID:                      fmt.Sprintf("CON-TB-%03d", len(results)),
					Type:                    ContradictionTypeTRUST_BOUNDARY,
					Severity:                RiskHigh,
					Confidence:              0.88,
					Summary:                 "Vendor trust boundary contradicts internal-only claim",
					Description:             fmt.Sprintf("Vendor trust boundary at %s contradicts claims that the system is internal-only.", strings.Join(boundary.Components, ", ")),
					Reasoning:               "A system cannot be both internal-only and have vendor components. Vendor access introduces third-party risk that must be explicitly managed.",
					AffectedTrustBoundaries: []string{boundary.Type},
					Recommendations: []string{
						"Document vendor access with data processing agreements",
						"Implement vendor monitoring and audit access",
						"Require vendor security assessments",
					},
				})
			}
		}
	}

	return results
}

// ─────────────────────────────────────────────────────────────
// PHASE 6 — CONTROL CONTRADICTIONS
// ─────────────────────────────────────────────────────────────

func (cie *CIEEngine) detectControlContradictions(assumptions []Assumption, controls []ControlDetail) []CIEContradiction {
	var results []CIEContradiction

	// Check each control against assumptions
	for _, control := range controls {
		controlText := strings.ToLower(control.Name + " " + control.Description)

		// MFA control + service account bypass assumption
		if strings.Contains(controlText, "mfa") || strings.Contains(controlText, "multi-factor") {
			if cie.hasAssumption(assumptions, []string{"service account", "system account", "bot", "automation"}, []string{"bypass", "exempt", "skip", "no mfa", "without mfa"}) {
				results = append(results, CIEContradiction{
					ID:               fmt.Sprintf("CON-CTRL-%03d", len(results)),
					Type:             ContradictionTypeCONTROL,
					Severity:         RiskHigh,
					Confidence:       0.90,
					Summary:          "MFA control exists but service accounts bypass it",
					Description:      fmt.Sprintf("Control '%s' enforces MFA, but assumptions allow service accounts to bypass MFA.", control.Name),
					Reasoning:        "A control is only effective if it covers all access paths. Service accounts bypassing MFA create a gap that defeats the control.",
					AffectedControls: []string{control.ID},
					Recommendations: []string{
						"Extend MFA control to service accounts using workload identity",
						"Implement certificate-based authentication for service accounts",
						"Monitor service account access with anomaly detection",
					},
				})
			}
		}

		// Encryption control + plaintext backup assumption
		if strings.Contains(controlText, "encrypt") || strings.Contains(controlText, "encryption") {
			if cie.hasAssumption(assumptions, []string{"backup", "restore"}, []string{"plaintext", "unencrypted", "no encryption", "not encrypted"}) {
				results = append(results, CIEContradiction{
					ID:               fmt.Sprintf("CON-CTRL-%03d", len(results)),
					Type:             ContradictionTypeCONTROL,
					Severity:         RiskCritical,
					Confidence:       0.92,
					Summary:          "Encryption control exists but backups are plaintext",
					Description:      fmt.Sprintf("Control '%s' enforces encryption, but assumptions allow plaintext backups.", control.Name),
					Reasoning:        "Plaintext backups bypass encryption controls. If the primary data is encrypted but backups are not, the protection is incomplete.",
					AffectedControls: []string{control.ID},
					Recommendations: []string{
						"Extend encryption control to backup systems",
						"Encrypt backups at rest and in transit",
						"Test backup encryption with periodic restore validation",
					},
				})
			}
		}
	}

	return results
}

// ─────────────────────────────────────────────────────────────
// PHASE 7 — COMPLIANCE FRAMEWORK CONTRADICTIONS
// ─────────────────────────────────────────────────────────────

func (cie *CIEEngine) detectComplianceFrameworkContradictions(assumptions []Assumption, arch *ArchDescription) []CIEContradiction {
	var results []CIEContradiction

	complianceMap := map[string][]string{
		"HIPAA":    {"audit", "access", "encryption", "integrity", "backup", "retention", "phi"},
		"SOC2":     {"access", "monitoring", "change", "encryption", "backup", "incident"},
		"ISO27001": {"access", "risk", "supplier", "asset", "incident", "business continuity"},
		"PCI DSS":  {"encryption", "access", "monitoring", "testing", "network", "physical"},
		"GDPR":     {"consent", "access", "deletion", "encryption", "breach", "processor"},
		"FedRAMP":  {"access", "encryption", "monitoring", "incident", "backup", "contingency"},
	}

	for _, comp := range arch.Compliance {
		compLower := strings.ToLower(comp)
		requirements, exists := complianceMap[compLower]
		if !exists {
			continue
		}

		for _, req := range requirements {
			if !cie.hasAssumption(assumptions, []string{req}, nil) {
				results = append(results, CIEContradiction{
					ID:          fmt.Sprintf("CON-COMP-FW-%03d", len(results)),
					Type:        ContradictionTypeCOMPLIANCE,
					Severity:    RiskHigh,
					Confidence:  0.85,
					Summary:     fmt.Sprintf("%s compliance requires %s but not documented", comp, req),
					Description: fmt.Sprintf("%s compliance framework requires %s controls, but the architecture does not document them.", comp, req),
					Reasoning:   fmt.Sprintf("%s has specific requirements for %s. Without documented controls, compliance cannot be demonstrated or audited.", comp, req),
					Recommendations: []string{
						fmt.Sprintf("Document %s controls for %s compliance", req, comp),
						fmt.Sprintf("Implement %s and verify with control testing", req),
						"Add compliance controls to the architecture definition",
					},
				})
			}
		}
	}

	return results
}

// ─────────────────────────────────────────────────────────────
// PHASE 8 — SCORING
// ─────────────────────────────────────────────────────────────

func (cie *CIEEngine) scoreContradiction(con CIEContradiction, assumptions []Assumption, arch *ArchDescription) CIEContradiction {
	score := 0.0

	// Base: severity
	switch con.Severity {
	case RiskCritical:
		score += 1.0
	case RiskHigh:
		score += 0.7
	case RiskMedium:
		score += 0.4
	case RiskLow:
		score += 0.2
	}

	// PHI presence boosts score
	if cie.hasAssumption(assumptions, []string{"phi", "patient", "health", "medical", "ehr", "hipaa"}, nil) {
		score += 0.3
		if con.Severity == RiskHigh {
			con.Severity = RiskCritical
		}
	}

	// PCI presence boosts score
	if cie.hasAssumption(assumptions, []string{"pci", "payment", "card", "cardholder", "financial"}, nil) {
		score += 0.3
		if con.Severity == RiskHigh {
			con.Severity = RiskCritical
		}
	}

	// Identity systems boost score
	if con.Type == ContradictionTypeAUTHENTICATION || con.Type == ContradictionTypeIDENTITY {
		score += 0.2
	}

	// Trust boundary involvement
	if len(con.AffectedTrustBoundaries) > 0 {
		score += 0.15
	}

	// Multiple affected components
	if len(con.AffectedComponents) > 2 {
		score += 0.1
	}

	// Direct control conflict
	if con.Type == ContradictionTypeCONTROL {
		score += 0.2
	}

	// Clamp confidence
	if score > 1.0 {
		score = 1.0
	}
	con.Confidence = score

	return con
}

// ─────────────────────────────────────────────────────────────
// HELPER METHODS
// ─────────────────────────────────────────────────────────────

func (cie *CIEEngine) findClaims(claims []Statement, subjects, predicates []string) []Statement {
	return cie.findClaimsFiltered(claims, subjects, predicates, nil)
}

func (cie *CIEEngine) findClaimsFiltered(claims []Statement, subjects, predicates, exclude []string) []Statement {
	var matches []Statement
	for _, c := range claims {
		text := strings.ToLower(c.OriginalText)
		subjMatch := len(subjects) == 0
		predMatch := len(predicates) == 0
		for _, s := range subjects {
			if strings.Contains(text, strings.ToLower(s)) {
				subjMatch = true
				break
			}
		}
		for _, p := range predicates {
			if strings.Contains(text, strings.ToLower(p)) {
				predMatch = true
				break
			}
		}
		if !subjMatch || !predMatch {
			continue
		}
		excluded := false
		for _, e := range exclude {
			if strings.Contains(text, strings.ToLower(e)) {
				excluded = true
				break
			}
		}
		if excluded {
			continue
		}
		matches = append(matches, c)
	}
	return matches
}

func (cie *CIEEngine) hasAssumption(assumptions []Assumption, subjects, predicates []string) bool {
	for _, a := range assumptions {
		text := strings.ToLower(a.Description)
		subjMatch := len(subjects) == 0
		predMatch := len(predicates) == 0
		for _, s := range subjects {
			if strings.Contains(text, strings.ToLower(s)) {
				subjMatch = true
				break
			}
		}
		for _, p := range predicates {
			if strings.Contains(text, strings.ToLower(p)) {
				predMatch = true
				break
			}
		}
		if subjMatch && predMatch {
			return true
		}
	}
	return false
}

func (cie *CIEEngine) hasClaim(claims []Statement, subjects, predicates []string) bool {
	for _, c := range claims {
		text := strings.ToLower(c.OriginalText)
		subjMatch := len(subjects) == 0
		predMatch := len(predicates) == 0
		for _, s := range subjects {
			if strings.Contains(text, strings.ToLower(s)) {
				subjMatch = true
				break
			}
		}
		for _, p := range predicates {
			if strings.Contains(text, strings.ToLower(p)) {
				predMatch = true
				break
			}
		}
		if subjMatch && predMatch {
			return true
		}
	}
	return false
}

func (cie *CIEEngine) hasCompliance(arch *ArchDescription, keywords ...string) bool {
	for _, c := range arch.Compliance {
		cl := strings.ToLower(c)
		for _, kw := range keywords {
			if strings.Contains(cl, strings.ToLower(kw)) {
				return true
			}
		}
	}
	return false
}

// CountContradictionsBySeverity returns counts per severity.
func CountContradictionsBySeverity(contradictions []CIEContradiction) map[RiskLevel]int {
	counts := make(map[RiskLevel]int)
	for _, c := range contradictions {
		counts[c.Severity]++
	}
	return counts
}

// BuildContradictionSummary creates a human-readable summary.
func BuildContradictionSummary(contradictions []CIEContradiction) string {
	if len(contradictions) == 0 {
		return "No contradictions detected."
	}
	counts := CountContradictionsBySeverity(contradictions)
	return fmt.Sprintf(
		"Contradictions Found: %d (Critical: %d, High: %d, Medium: %d, Low: %d)",
		len(contradictions),
		counts[RiskCritical],
		counts[RiskHigh],
		counts[RiskMedium],
		counts[RiskLow],
	)
}
