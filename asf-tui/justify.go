package main

import (
	"fmt"
	"strings"
)

// ──────────────────────────────────────────────
// Evidence Engine
// ──────────────────────────────────────────────

// EvidenceEngine traces assumptions back to architecture artifacts.
type EvidenceEngine struct {
	arch *ArchDescription
	sourcePath string
	sourceType string
}

// NewEvidenceEngine creates an evidence engine from a parsed architecture.
func NewEvidenceEngine(arch *ArchDescription, sourcePath string) *EvidenceEngine {
	sourceType := extractSourceType(sourcePath)
	return &EvidenceEngine{
		arch:       arch,
		sourcePath: sourcePath,
		sourceType: sourceType,
	}
}

func extractSourceType(path string) string {
	if idx := strings.LastIndex(path, "."); idx >= 0 {
		return path[idx:]
	}
	return ".txt"
}

// EvidenceResult holds matched evidence for a single assumption.
type EvidenceResult struct {
	MatchedComponents      []string
	MatchedRelationships   []string
	MatchedTrustBoundaries []string
	MatchedSecurityConcepts []string
	PrimarySourceNode      string
	SourceLine             int
	EvidenceCount          int
}

// TraceEvidence matches an assumption against the architecture description.
func (ee *EvidenceEngine) TraceEvidence(category string, keywords []string, text string) *EvidenceResult {
	result := &EvidenceResult{}

	if ee.arch == nil {
		return result
	}

	searchText := strings.ToLower(category + " " + text + " " + strings.Join(keywords, " "))

	// Match components — also identify the primary source node
	var primaryNode string
	for _, comp := range ee.arch.Components {
		label := strings.ToLower(comp.Label)
		if label == "" {
			continue
		}
		if strings.Contains(searchText, label) {
			result.MatchedComponents = append(result.MatchedComponents, comp.Label)
			if primaryNode == "" {
				primaryNode = comp.Label
			}
		}
	}

	// Match relationships
	for _, rel := range ee.arch.Relationships {
		srcLower := strings.ToLower(rel.Source)
		tgtLower := strings.ToLower(rel.Target)
		relLabel := strings.ToLower(rel.Label)

		if strings.Contains(searchText, srcLower) || strings.Contains(searchText, tgtLower) {
			result.MatchedRelationships = append(result.MatchedRelationships,
				fmt.Sprintf("%s → %s", rel.Source, rel.Target))
		}

		_ = relLabel
	}

	// Set source node from best match
	if primaryNode == "" {
		// Fall back to first keyword as implied source
		for _, kw := range keywords {
			if kw != "" {
				primaryNode = kw
				break
			}
		}
	}
	result.PrimarySourceNode = primaryNode

	// Match trust boundaries (inferred from relationships between different zones)
	trustZones := map[string]bool{
		"internet": true, "external": true, "public": true,
		"vpn": true, "gateway": true, "dmz": true,
	}
	for _, rel := range ee.arch.Relationships {
		srcLower := strings.ToLower(rel.Source)
		tgtLower := strings.ToLower(rel.Target)
		for _, comp := range ee.arch.Components {
			cl := strings.ToLower(comp.Label)
			if trustZones[cl] && (strings.Contains(searchText, cl) ||
				strings.Contains(srcLower, cl) || strings.Contains(tgtLower, cl)) {
				result.MatchedTrustBoundaries = append(result.MatchedTrustBoundaries,
					fmt.Sprintf("trust boundary at %s (%s to %s)", comp.Label, rel.Source, rel.Target))
			}
		}
	}

	// Match security concepts
	securityConcepts := map[string][]string{
		"authentication":    {"auth", "login", "password", "mfa", "sso", "identity", "credential"},
		"authorization":     {"access", "role", "permission", "rbac", "acl", "policy"},
		"encryption":        {"tls", "ssl", "https", "encrypt", "cipher", "signing"},
		"network_security":  {"firewall", "vpn", "segment", "gateway", "proxy"},
		"data_protection":   {"database", "backup", "storage", "encryption at rest"},
		"audit_logging":     {"log", "audit", "monitor", "alert"},
		"dependency":        {"third", "vendor", "supply chain", "external"},
		"session_management": {"session", "token", "cookie", "jwt"},
	}
	for concept, patterns := range securityConcepts {
		for _, p := range patterns {
			if strings.Contains(searchText, p) {
				result.MatchedSecurityConcepts = append(result.MatchedSecurityConcepts, concept)
				break
			}
		}
	}

	result.EvidenceCount = len(result.MatchedComponents) + len(result.MatchedRelationships) +
		len(result.MatchedTrustBoundaries) + len(result.MatchedSecurityConcepts)

	return result
}

// FindSourceLine searches the raw architecture text for the line containing
// the best matching evidence for this assumption.
func (ee *EvidenceEngine) FindSourceLine(searchText string, evidence *EvidenceResult) int {
	if ee.arch == nil || ee.arch.RawText == "" {
		return 0
	}
	// Try to find the primary matched component or a keyword in the raw text
	searchFor := strings.ToLower(strings.TrimSpace(searchText))
	lines := strings.Split(ee.arch.RawText, "\n")
	for i, line := range lines {
		ll := strings.ToLower(strings.TrimSpace(line))
		if ll == "" {
			continue
		}
		if strings.Contains(searchFor, ll) || strings.Contains(ll, searchFor[:min(len(searchFor), 40)]) {
			return i + 1
		}
	}
	// Fallback: find first line containing any matched component
	if evidence != nil {
		for _, comp := range evidence.MatchedComponents {
			cl := strings.ToLower(comp)
			for i, line := range lines {
				if strings.Contains(strings.ToLower(line), cl) {
					return i + 1
				}
			}
		}
	}
	return 0
}

// BuildEvidenceSources builds the evidence source strings for an assumption.
func (ee *EvidenceEngine) BuildEvidenceSources(evidence *EvidenceResult) []string {
	var sources []string
	sources = append(sources, fmt.Sprintf("source: %s (%s)", ee.sourcePath, ee.sourceType))
	if evidence.PrimarySourceNode != "" {
		sources = append(sources, fmt.Sprintf("source node: %s", evidence.PrimarySourceNode))
	}
	if evidence.SourceLine > 0 {
		sources = append(sources, fmt.Sprintf("source line: %d", evidence.SourceLine))
	}
	for _, c := range evidence.MatchedComponents {
		sources = append(sources, fmt.Sprintf("component: %s", c))
	}
	for _, r := range evidence.MatchedRelationships {
		sources = append(sources, fmt.Sprintf("relationship: %s", r))
	}
	for _, tb := range evidence.MatchedTrustBoundaries {
		sources = append(sources, fmt.Sprintf("trust boundary: %s", tb))
	}
	for _, sc := range evidence.MatchedSecurityConcepts {
		sources = append(sources, fmt.Sprintf("security concept: %s", sc))
	}
	return sources
}

// BuildEvidenceSummary creates a top-level summary from all assumptions.
func (ee *EvidenceEngine) BuildEvidenceSummary(assumptions []Assumption) EvidenceSummary {
	es := EvidenceSummary{}
	seenFiles := make(map[string]bool)
	seenComps := make(map[string]bool)
	seenRels := make(map[string]bool)
	for _, a := range assumptions {
		for _, s := range a.EvidenceSources {
			es.TotalSources++
			if strings.HasPrefix(s, "source:") {
				f := strings.TrimPrefix(s, "source: ")
				if !seenFiles[f] {
					seenFiles[f] = true
					es.SourceFiles = append(es.SourceFiles, f)
				}
			}
		}
		for _, c := range a.SourceComponents {
			if !seenComps[c] {
				seenComps[c] = true
				es.TotalComponents++
			}
		}
		for _, r := range a.SourceRelationships {
			if !seenRels[r] {
				seenRels[r] = true
				es.TotalRelationships++
			}
		}
	}
	return es
}

// ──────────────────────────────────────────────
// Assumption Justification Engine
// ──────────────────────────────────────────────

// JustifyAssumption generates human-readable rationale for why an assumption exists.
func JustifyAssumption(category string, evidence *EvidenceResult) string {
	var parts []string

	if len(evidence.MatchedComponents) > 0 {
		comps := evidence.MatchedComponents
		if len(comps) > 3 {
			comps = comps[:3]
		}
		parts = append(parts, fmt.Sprintf("detected %d relevant component(s): %s",
			len(evidence.MatchedComponents), strings.Join(comps, ", ")))
	}

	if len(evidence.MatchedRelationships) > 0 {
		parts = append(parts, fmt.Sprintf("identified %d communication path(s) between components",
			len(evidence.MatchedRelationships)))
	}

	if len(evidence.MatchedTrustBoundaries) > 0 {
		parts = append(parts, fmt.Sprintf("crosses %d trust boundary/ies requiring security verification",
			len(evidence.MatchedTrustBoundaries)))
	}

	if len(evidence.MatchedSecurityConcepts) > 0 {
		concepts := evidence.MatchedSecurityConcepts
		if len(concepts) > 3 {
			concepts = concepts[:3]
		}
		parts = append(parts, fmt.Sprintf("relates to security concept(s): %s",
			strings.Join(concepts, ", ")))
	}

	if len(parts) == 0 {
		parts = append(parts, fmt.Sprintf("generated from category %s with %d keyword match(es)",
			category, 0))
	}

	return "ASF identified this assumption because " + strings.Join(parts, "; ") + "."
}

// ──────────────────────────────────────────────
// STRIDE Justification Engine
// ──────────────────────────────────────────────

// StrideJustifyEngine extends StrideEngine with justification output.
type StrideJustifyEngine struct {
	inner *StrideEngine
}

// NewStrideJustifyEngine creates a STRIDE justification engine.
func NewStrideJustifyEngine(inner *StrideEngine) *StrideJustifyEngine {
	return &StrideJustifyEngine{inner: inner}
}

// StrideResult holds both categories and justifications.
type StrideResult struct {
	Categories     []StrideCategory
	Justifications []StrideJustification
}

// Justify returns STRIDE categories with full justifications.
func (sje *StrideJustifyEngine) Justify(category string, text string, keywords []string, matchedComponents []string) *StrideResult {
	result := &StrideResult{}

	searchText := strings.ToLower(category + " " + text + " " + strings.Join(keywords, " "))

	// Get categories from the inner engine
	result.Categories = sje.inner.MapAssumption(category, text, keywords)

	// Build justification for each matched category
	for _, cat := range result.Categories {
		just := StrideJustification{
			Category:          cat,
			MatchedComponents: matchedComponents,
		}

		// Find which rules matched
		for idx, rule := range sje.inner.GetKeywordRules() {
			for _, kw := range rule.keywords {
				if strings.Contains(searchText, kw) {
					for _, scat := range rule.stride {
						if scat == cat {
							just.MatchedRuleIndexes = append(just.MatchedRuleIndexes, idx)
							just.MatchedKeywords = append(just.MatchedKeywords, kw)
						}
					}
					break
				}
			}
		}

		// Also check category rules
		for catRule, strideCats := range sje.inner.GetCategoryRules() {
			if strings.EqualFold(category, catRule) {
				for _, scat := range strideCats {
					if scat == cat {
						just.MatchedKeywords = append(just.MatchedKeywords, fmt.Sprintf("category:%s", catRule))
					}
				}
			}
		}

		just.Reason = buildStrideReason(cat, category, just.MatchedKeywords, matchedComponents)
		just.Confidence, just.ConfidenceReason = calculateStrideConfidence(cat, len(just.MatchedKeywords), len(matchedComponents))

		result.Justifications = append(result.Justifications, just)
	}

	return result
}

func buildStrideReason(cat StrideCategory, assumptionCategory string, keywords []string, comps []string) string {
	switch cat {
	case StrideSpoofing:
		return fmt.Sprintf("identity verification required — %s %s",
			assumptionCategory, joinKeywords(keywords))
	case StrideTampering:
		return fmt.Sprintf("data integrity risk — %s %s",
			assumptionCategory, joinKeywords(keywords))
	case StrideRepudiation:
		return fmt.Sprintf("non-repudiation concern — %s %s",
			assumptionCategory, joinKeywords(keywords))
	case StrideInfoDisclosure:
		return fmt.Sprintf("information disclosure risk — %s %s",
			assumptionCategory, joinKeywords(keywords))
	case StrideDenialOfService:
		return fmt.Sprintf("availability risk — %s %s",
			assumptionCategory, joinKeywords(keywords))
	case StrideElevationPriv:
		return fmt.Sprintf("privilege escalation risk — %s %s",
			assumptionCategory, joinKeywords(keywords))
	}
	return fmt.Sprintf("%s mapping from %s", cat, assumptionCategory)
}

func joinKeywords(kw []string) string {
	if len(kw) == 0 {
		return ""
	}
	unique := make(map[string]bool)
	var parts []string
	for _, k := range kw {
		if !unique[k] {
			unique[k] = true
			parts = append(parts, k)
		}
	}
	if len(parts) > 5 {
		parts = parts[:5]
	}
	return "(matched: " + strings.Join(parts, ", ") + ")"
}

func calculateStrideConfidence(cat StrideCategory, keywordMatches, componentMatches int) (float64, string) {
	score := 0.3 // base confidence
	factors := []string{"base confidence 0.3"}

	if keywordMatches > 0 {
		kwScore := float64(keywordMatches) * 0.1
		if kwScore > 0.4 {
			kwScore = 0.4
		}
		score += kwScore
		factors = append(factors, fmt.Sprintf("+%.2f from %d keyword match(es)", kwScore, keywordMatches))
	}
	if componentMatches > 0 {
		compScore := float64(componentMatches) * 0.08
		if compScore > 0.3 {
			compScore = 0.3
		}
		score += compScore
		factors = append(factors, fmt.Sprintf("+%.2f from %d component match(es)", compScore, componentMatches))
	}
	if score > 0.95 {
		score = 0.95
	}

	return score, strings.Join(factors, "; ")
}

// ──────────────────────────────────────────────
// Risk Justification Engine
// ──────────────────────────────────────────────

// LikelihoodAnalyzer evaluates how likely an assumption is to be exploited.
type LikelihoodAnalyzer struct{}

// AnalyzeLikelihood determines likelihood score with justification.
func (la *LikelihoodAnalyzer) AnalyzeLikelihood(assumption *Assumption, evidence *EvidenceResult) (int, string, []LikelihoodFactor) {
	var factors []LikelihoodFactor
	score := 1 // base

	// Factor 1: Exposure level
	exposure := 1
	exposureReason := "internal component with limited exposure"
	for _, sc := range evidence.MatchedSecurityConcepts {
		switch sc {
		case "network_security":
			if containsAny(evidence.MatchedComponents, "internet", "gateway", "public") {
				exposure = 5
				exposureReason = "internet-exposed component"
			} else {
				exposure = 3
				exposureReason = "network-accessible component"
			}
		}
	}
	if strings.Contains(strings.ToLower(assumption.Category), "network") ||
		strings.Contains(strings.ToLower(assumption.Category), "internet") {
		exposure = 4
		exposureReason = "network/internet category"
	}
	factors = append(factors, LikelihoodFactor{
		Factor: "Exposure Level", Value: exposure, Reason: exposureReason,
	})
	score += exposure - 1

	// Factor 2: Authentication dependency
	authScore := 2
	authReason := "standard security boundary"
	for _, sc := range evidence.MatchedSecurityConcepts {
		if sc == "authentication" || sc == "authorization" {
			authScore = 4
			authReason = "dependent on authentication/authorization controls"
		}
	}
	if containsAny(evidence.MatchedComponents, "auth", "login", "sso", "mfa", "identity") {
		authScore = 4
		authReason = "authentication component identified"
	}
	factors = append(factors, LikelihoodFactor{
		Factor: "Authentication Dependency", Value: authScore, Reason: authReason,
	})
	score += authScore - 2

	// Factor 3: Attack complexity
	complexity := 2
	complexityReason := "moderate attack complexity"
	if len(evidence.MatchedComponents) > 3 {
		complexity = 3
		complexityReason = "multiple components increase attack surface"
	}
	if len(evidence.MatchedRelationships) > 5 {
		complexity = 4
		complexityReason = "complex relationship graph increases attack paths"
	}
	factors = append(factors, LikelihoodFactor{
		Factor: "Attack Surface Complexity", Value: complexity, Reason: complexityReason,
	})
	score += complexity - 2

	// Clamp to 1-5
	if score < 1 {
		score = 1
	}
	if score > 5 {
		score = 5
	}

	reason := fmt.Sprintf("likelihood %d/5 based on exposure(%d), auth dependency(%d), attack complexity(%d)",
		score, exposure, authScore, complexity)

	return score, reason, factors
}

// ImpactAnalyzer evaluates the potential impact of an assumption being violated.
type ImpactAnalyzer struct{}

// AnalyzeImpact determines impact score with justification.
func (ia *ImpactAnalyzer) AnalyzeImpact(assumption *Assumption, evidence *EvidenceResult) (int, string, []ImpactFactor) {
	var factors []ImpactFactor
	score := 1 // base

	// Factor 1: Data classification
	dataScore := 2
	dataReason := "standard business data"
	for _, sc := range evidence.MatchedSecurityConcepts {
		if sc == "data_protection" {
			dataScore = 4
			dataReason = "sensitive data handling"
		}
	}
	if containsAny(evidence.MatchedComponents, "database", "db", "storage", "backup") {
		dataScore = 4
		dataReason = "database/storage component with sensitive data"
	}
	if containsAny(evidence.MatchedComponents, "pii", "phi", "financial", "health", "payment", "card") {
		dataScore = 5
		dataReason = "regulated data (PII/PHI/financial)"
	}
	factors = append(factors, ImpactFactor{
		Factor: "Data Classification", Value: dataScore, Reason: dataReason,
	})
	score += dataScore - 2

	// Factor 2: Regulatory exposure
	regScore := 1
	regReason := "no direct regulatory exposure detected"
	for _, comp := range evidence.MatchedComponents {
		cl := strings.ToLower(comp)
		if strings.Contains(cl, "health") || strings.Contains(cl, "hipaa") {
			regScore = 5
			regReason = "healthcare data (HIPAA regulated)"
		} else if strings.Contains(cl, "payment") || strings.Contains(cl, "card") || strings.Contains(cl, "pci") {
			regScore = 5
			regReason = "payment data (PCI DSS regulated)"
		} else if strings.Contains(cl, "financial") || strings.Contains(cl, "sox") || strings.Contains(cl, "audit") {
			regScore = 4
			regReason = "financial data (SOX regulated)"
		} else if strings.Contains(cl, "pii") || strings.Contains(cl, "gdpr") || strings.Contains(cl, "privacy") {
			regScore = 4
			regReason = "personal data (GDPR/CCPA regulated)"
		}
	}
	factors = append(factors, ImpactFactor{
		Factor: "Regulatory Exposure", Value: regScore, Reason: regReason,
	})
	score += regScore - 1

	// Factor 3: Business criticality
	criticalScore := 2
	criticalReason := "standard business process"
	for _, comp := range evidence.MatchedComponents {
		cl := strings.ToLower(comp)
		if strings.Contains(cl, "core") || strings.Contains(cl, "main") || strings.Contains(cl, "primary") {
			criticalScore = 4
			criticalReason = "core business component"
		}
	}
	if len(evidence.MatchedRelationships) > 8 {
		criticalScore = 4
		criticalReason = "highly interconnected — cascading failure risk"
	}
	factors = append(factors, ImpactFactor{
		Factor: "Business Criticality", Value: criticalScore, Reason: criticalReason,
	})
	score += criticalScore - 2

	// Clamp to 1-5
	if score < 1 {
		score = 1
	}
	if score > 5 {
		score = 5
	}

	reason := fmt.Sprintf("impact %d/5 based on data classification(%d), regulatory(%d), criticality(%d)",
		score, dataScore, regScore, criticalScore)
	return score, reason, factors
}

// RiskMatrix implements a 5x5 risk matrix.
type RiskMatrix struct{}

// Calculate computes risk score and level from likelihood and impact.
func (rm *RiskMatrix) Calculate(likelihood, impact int) (int, RiskLevel) {
	score := likelihood * impact

	var level RiskLevel
	switch {
	case score >= 20:
		level = RiskCritical
	case score >= 12:
		level = RiskHigh
	case score >= 5:
		level = RiskMedium
	default:
		level = RiskLow
	}

	return score, level
}

// RiskReason generates human-readable risk justification.
func (rm *RiskMatrix) RiskReason(likelihood, impact, score int, level RiskLevel) string {
	reason := fmt.Sprintf("risk score %d/25 (likelihood %d × impact %d) = %s", score, likelihood, impact, level)
	return reason
}

// ──────────────────────────────────────────────
// Confidence Engine
// ──────────────────────────────────────────────

// ConfidenceEngine calculates confidence in an assumption.
type ConfidenceEngine struct{}

// CalculateConfidence computes a confidence score with justification.
func (ce *ConfidenceEngine) CalculateConfidence(evidenceCount int, strideMatchCount int, componentCount int, relationshipCount int) (float64, string) {
	score := 0.1 // base
	var factors []string

	if evidenceCount > 0 {
		evScore := float64(evidenceCount) * 0.05
		if evScore > 0.3 {
			evScore = 0.3
		}
		score += evScore
		factors = append(factors, fmt.Sprintf("+%.2f from %d evidence point(s)", evScore, evidenceCount))
	}
	if strideMatchCount > 0 {
		stScore := float64(strideMatchCount) * 0.08
		if stScore > 0.25 {
			stScore = 0.25
		}
		score += stScore
		factors = append(factors, fmt.Sprintf("+%.2f from %d STRIDE rule match(es)", stScore, strideMatchCount))
	}
	if componentCount > 0 {
		compScore := float64(componentCount) * 0.06
		if compScore > 0.2 {
			compScore = 0.2
		}
		score += compScore
		factors = append(factors, fmt.Sprintf("+%.2f from %d component match(es)", compScore, componentCount))
	}
	if relationshipCount > 0 {
		relScore := float64(relationshipCount) * 0.04
		if relScore > 0.15 {
			relScore = 0.15
		}
		score += relScore
		factors = append(factors, fmt.Sprintf("+%.2f from %d relationship match(es)", relScore, relationshipCount))
	}
	if score > 0.95 {
		score = 0.95
	}

	reason := strings.Join(factors, "; ")
	if reason == "" {
		reason = "base confidence 0.10 (no evidence available)"
	}

	return score, reason
}

// ──────────────────────────────────────────────
// Pipeline Orchestrator
// ──────────────────────────────────────────────

// ExplainabilityPipeline orchestrates all explainability engines.
type ExplainabilityPipeline struct {
	evidenceEngine     *EvidenceEngine
	strideJustify      *StrideJustifyEngine
	likelihoodAnalyzer *LikelihoodAnalyzer
	impactAnalyzer     *ImpactAnalyzer
	riskMatrix         *RiskMatrix
	confidenceEngine   *ConfidenceEngine
}

// NewExplainabilityPipeline creates a new pipeline.
func NewExplainabilityPipeline(arch *ArchDescription, sourcePath string, strideEngine *StrideEngine) *ExplainabilityPipeline {
	return &ExplainabilityPipeline{
		evidenceEngine:     NewEvidenceEngine(arch, sourcePath),
		strideJustify:      NewStrideJustifyEngine(strideEngine),
		likelihoodAnalyzer: &LikelihoodAnalyzer{},
		impactAnalyzer:     &ImpactAnalyzer{},
		riskMatrix:         &RiskMatrix{},
		confidenceEngine:   &ConfidenceEngine{},
	}
}

// Explain processes a single assumption through the full pipeline.
func (ep *ExplainabilityPipeline) Explain(a *Assumption) {
	if a == nil {
		return
	}

	// Phase 1: Trace evidence
	evidence := ep.evidenceEngine.TraceEvidence(a.Category, a.Keywords, a.Description)

	// Phase 1b: Populate source traceability
	searchText := a.Category + " " + a.Description + " " + strings.Join(a.Keywords, " ")
	evidence.SourceLine = ep.evidenceEngine.FindSourceLine(searchText, evidence)
	a.SourceNode = evidence.PrimarySourceNode
	a.SourceLine = evidence.SourceLine

	// Phase 2: Build evidence sources
	a.EvidenceSources = ep.evidenceEngine.BuildEvidenceSources(evidence)
	a.SourceComponents = evidence.MatchedComponents
	a.SourceRelationships = evidence.MatchedRelationships

	// Phase 3: Generate rationale
	a.Rationale = JustifyAssumption(a.Category, evidence)

	// Phase 4: STRIDE justification
	strideResult := ep.strideJustify.Justify(a.Category, a.Description, a.Keywords, evidence.MatchedComponents)
	a.Stride = strideResult.Categories
	a.StrideJustifications = strideResult.Justifications

	// Phase 5: Risk justification
	lh, lhReason, lhFactors := ep.likelihoodAnalyzer.AnalyzeLikelihood(a, evidence)
	im, imReason, imFactors := ep.impactAnalyzer.AnalyzeImpact(a, evidence)
	a.Likelihood = lh
	a.Impact = im
	riskScore, riskLevel := ep.riskMatrix.Calculate(lh, im)
	a.Risk = riskLevel
	a.RiskJustification = &RiskJustification{
		Likelihood:        lh,
		LikelihoodReason:  lhReason,
		LikelihoodFactors: lhFactors,
		Impact:            im,
		ImpactReason:      imReason,
		ImpactFactors:     imFactors,
		RiskScore:         riskScore,
		RiskLevel:         riskLevel,
		RiskReason:        ep.riskMatrix.RiskReason(lh, im, riskScore, riskLevel),
	}

	// Phase 6: Confidence scoring
	strideMatchCount := 0
	for _, j := range strideResult.Justifications {
		strideMatchCount += len(j.MatchedKeywords)
	}
	conf, confReason := ep.confidenceEngine.CalculateConfidence(
		evidence.EvidenceCount,
		strideMatchCount,
		len(evidence.MatchedComponents),
		len(evidence.MatchedRelationships),
	)
	a.Confidence = conf
	a.RiskJustification.Confidence = conf
	a.RiskJustification.ConfidenceReason = confReason
}

// BuildEvidenceSummary creates the top-level evidence summary.
func (ep *ExplainabilityPipeline) BuildEvidenceSummary(assumptions []Assumption) EvidenceSummary {
	return ep.evidenceEngine.BuildEvidenceSummary(assumptions)
}

// ──────────────────────────────────────────────
// Helpers
// ──────────────────────────────────────────────

func containsAny(items []string, targets ...string) bool {
	for _, item := range items {
		il := strings.ToLower(item)
		for _, t := range targets {
			if strings.Contains(il, strings.ToLower(t)) {
				return true
			}
		}
	}
	return false
}
