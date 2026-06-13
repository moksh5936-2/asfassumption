package confidencex

import (
	"math"
	"sort"
	"strings"
	"time"
)

type ExplainabilityEngine struct {
	inputs []ConfidenceInput
	domain string
}

func NewExplainabilityEngine(domain string, inputs []ConfidenceInput) *ExplainabilityEngine {
	return &ExplainabilityEngine{
		domain: domain,
		inputs: inputs,
	}
}

func (e *ExplainabilityEngine) RunAll() *ConfidenceOutput {
	output := &ConfidenceOutput{
		Domain:      e.domain,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
	}
	if len(e.inputs) == 0 {
		return output
	}

	breakdowns := make([]ConfidenceBreakdown, len(e.inputs))
	for i, in := range e.inputs {
		breakdowns[i] = e.analyzeAssumption(in)
	}

	output.Breakdowns = breakdowns
	output.CISOTrustView = e.buildCISOView(breakdowns)
	output.ArchitectReviewView = e.buildArchitectView(breakdowns)

	return output
}

func (e *ExplainabilityEngine) analyzeAssumption(in ConfidenceInput) ConfidenceBreakdown {
	posFactors := e.collectPositiveFactors(in)
	negFactors := e.collectNegativeFactors(in)
	factContribs := e.computeFactContributions(in)
	evContribs := e.computeEvidenceContributions(in)
	domainContrib := e.computeDomainContribution(in)
	trustContrib := e.computeTrustContribution(in)

	totalBiasExplanation := 0.0
	for _, f := range posFactors {
		totalBiasExplanation += f.Impact
	}
	for _, f := range negFactors {
		totalBiasExplanation += f.Impact
	}

	adjustedConf := in.Confidence + totalBiasExplanation
	if adjustedConf < 0 {
		adjustedConf = 0
	}
	if adjustedConf > 100 {
		adjustedConf = 100
	}

	stability := e.classifyStability(in.Confidence, posFactors, negFactors)
	stabilityReason := e.stabilityReason(stability, in.Confidence, in)

	return ConfidenceBreakdown{
		AssumptionID:            in.AssumptionID,
		AssumptionText:          in.AssumptionText,
		FinalConfidence:         math.Round(in.Confidence*100) / 100,
		AdjustedConfidence:      math.Round(adjustedConf*100) / 100,
		StabilityClass:          stability,
		StabilityReason:         stabilityReason,
		PositiveFactors:         posFactors,
		NegativeFactors:         negFactors,
		SupportingFacts:         factContribs,
		EvidenceContributions:   evContribs,
		DomainContribution:      domainContrib,
		TrustContribution:       trustContrib,
		WhyExists:               e.generateWhyExists(in, posFactors),
		WhyUncertain:            e.generateWhyUncertain(in, negFactors),
		WhatIncreasesConfidence: e.generateWhatIncreases(in, negFactors),
		WhatDecreasesConfidence: e.generateWhatDecreases(in, posFactors, domainContrib),
	}
}

func (e *ExplainabilityEngine) collectPositiveFactors(in ConfidenceInput) []ConfidenceFactor {
	var factors []ConfidenceFactor

	if in.Confidence >= 80 {
		factors = append(factors, ConfidenceFactor{
			Name: "High Baseline Confidence", Type: "positive", Impact: 5.0,
			Description: "Baseline confidence is already strong (>=80%)",
		})
	} else if in.Confidence >= 60 {
		factors = append(factors, ConfidenceFactor{
			Name: "Moderate Baseline Confidence", Type: "positive", Impact: 2.0,
			Description: "Baseline confidence is moderate (60-80%)",
		})
	}

	if len(in.SourceComponents) > 0 {
		factors = append(factors, ConfidenceFactor{
			Name: "Component Traceability", Type: "positive", Impact: 3.0,
			Description: "Assumption is traceable to specific architecture components",
		})
	}

	if len(in.Keywords) >= 3 {
		factors = append(factors, ConfidenceFactor{
			Name: "Rich Keyword Coverage", Type: "positive", Impact: 2.0,
			Description: "Assumption has strong keyword coverage from architecture text",
		})
	}

	if in.HasTrustChain {
		factors = append(factors, ConfidenceFactor{
			Name: "Trust Chain Support", Type: "positive", Impact: 5.0,
			Description: "Assumption is part of a verified trust chain",
		})
	}

	for _, src := range in.EvidenceSources {
		if strings.Contains(strings.ToLower(src), "explicit") || strings.Contains(strings.ToLower(src), "policy") || strings.Contains(strings.ToLower(src), "control") {
			factors = append(factors, ConfidenceFactor{
				Name: "Explicit Evidence Source", Type: "positive", Impact: 4.0,
				Description: "Evidence from explicit source: " + src,
			})
			break
		}
	}

	if in.Domain != "" && in.Domain == e.domain {
		factors = append(factors, ConfidenceFactor{
			Name: "Domain Alignment", Type: "positive", Impact: 3.0,
			Description: "Assumption aligns with detected domain: " + in.Domain,
		})
	}

	if in.Rationale != "" && len(in.Rationale) > 50 {
		factors = append(factors, ConfidenceFactor{
			Name: "Detailed Rationale", Type: "positive", Impact: 2.0,
			Description: "Assumption has a detailed rationale explaining its origin",
		})
	}

	if len(in.SupportingFactTexts) > 0 {
		pct := float64(len(in.SupportingFactTexts)) * 2.0
		if pct > 10 {
			pct = 10
		}
		factors = append(factors, ConfidenceFactor{
			Name: "Supporting Facts Present", Type: "positive", Impact: pct,
			Description: "Architecture provides supporting facts for this assumption",
		})
	}

	if in.DependencyCentrality > 0.5 {
		factors = append(factors, ConfidenceFactor{
			Name: "High Dependency Centrality", Type: "positive", Impact: 4.0,
			Description: "Assumption involves a centrally important dependency",
		})
	}

	return factors
}

func (e *ExplainabilityEngine) collectNegativeFactors(in ConfidenceInput) []ConfidenceFactor {
	var factors []ConfidenceFactor

	if in.Confidence < 30 {
		factors = append(factors, ConfidenceFactor{
			Name: "Low Baseline Confidence", Type: "negative", Impact: -5.0,
			Description: "Baseline confidence is very low (<30%)",
		})
	} else if in.Confidence < 50 {
		factors = append(factors, ConfidenceFactor{
			Name: "Below Average Baseline", Type: "negative", Impact: -2.0,
			Description: "Baseline confidence is below average (<50%)",
		})
	}

	if len(in.EvidenceSources) == 0 {
		factors = append(factors, ConfidenceFactor{
			Name: "No Evidence Sources", Type: "negative", Impact: -8.0,
			Description: "No evidence sources are linked to this assumption",
		})
	}

	if in.VerificationStatus == "" || in.VerificationStatus == "UNKNOWN" || in.VerificationStatus == "Unverified" {
		factors = append(factors, ConfidenceFactor{
			Name: "Unverified Status", Type: "negative", Impact: -5.0,
			Description: "Assumption has not been verified",
		})
	}

	if in.HasCoverageGap {
		factors = append(factors, ConfidenceFactor{
			Name: "Coverage Gap Detected", Type: "negative", Impact: -6.0,
			Description: "Assumption falls in a known coverage gap area",
		})
	}

	if in.BlindSpotScore > 50 {
		factors = append(factors, ConfidenceFactor{
			Name: "High Blind Spot Risk", Type: "negative", Impact: -5.0,
			Description: "Assumption is in a high-scoring blind spot area",
		})
	} else if in.BlindSpotScore > 20 {
		factors = append(factors, ConfidenceFactor{
			Name: "Moderate Blind Spot Risk", Type: "negative", Impact: -2.0,
			Description: "Assumption is in a moderate blind spot area",
		})
	}

	if len(in.SourceComponents) == 0 && len(in.Keywords) == 0 {
		factors = append(factors, ConfidenceFactor{
			Name: "No Architecture Traceability", Type: "negative", Impact: -7.0,
			Description: "Assumption cannot be traced to any architecture component",
		})
	}

	if in.Rationale == "" {
		factors = append(factors, ConfidenceFactor{
			Name: "Missing Rationale", Type: "negative", Impact: -3.0,
			Description: "No rationale explaining why this assumption was generated",
		})
	}

	if len(in.SupportingFactTexts) == 0 {
		factors = append(factors, ConfidenceFactor{
			Name: "No Supporting Facts", Type: "negative", Impact: -4.0,
			Description: "Architecture contains no facts that directly support this assumption",
		})
	}

	if !in.HasTrustChain {
		factors = append(factors, ConfidenceFactor{
			Name: "No Trust Chain", Type: "negative", Impact: -3.0,
			Description: "Assumption is not part of any trust chain",
		})
	}

	return factors
}

func (e *ExplainabilityEngine) computeFactContributions(in ConfidenceInput) []FactContribution {
	var contribs []FactContribution
	for i, ft := range in.SupportingFactTexts {
		if i >= len(in.SupportingFactIDs) {
			break
		}
		contrib := 5.0 + float64(len(ft))/100.0*5.0
		if contrib > 15 {
			contrib = 15
		}
		category := ""
		if i < len(in.FactCategories) {
			category = in.FactCategories[i]
			switch category {
			case "security":
				contrib *= 1.2
			case "compliance":
				contrib *= 1.1
			}
		}
		if contrib > 15 {
			contrib = 15
		}
		pos := !strings.Contains(strings.ToLower(ft), "no ") && !strings.Contains(strings.ToLower(ft), "not ") && !strings.Contains(strings.ToLower(ft), "disabled") && !strings.Contains(strings.ToLower(ft), "absent")

		contribs = append(contribs, FactContribution{
			FactID:       in.SupportingFactIDs[i],
			FactText:     ft,
			Contribution: math.Round(contrib*10) / 10,
			IsPositive:   pos,
		})
	}
	return contribs
}

func (e *ExplainabilityEngine) computeEvidenceContributions(in ConfidenceInput) []EvidenceContribution {
	var contribs []EvidenceContribution
	if len(in.EvidenceSources) > 0 {
		for _, src := range in.EvidenceSources {
			impact := 5.0
			label := "Weak"
			srcLower := strings.ToLower(src)
			if strings.Contains(srcLower, "policy") || strings.Contains(srcLower, "control") || strings.Contains(srcLower, "explicit") {
				impact = 12.0
				label = "Strong"
			} else if strings.Contains(srcLower, "config") || strings.Contains(srcLower, "document") || strings.Contains(srcLower, "design") {
				impact = 8.0
				label = "Moderate"
			}
			contribs = append(contribs, EvidenceContribution{
				EvidenceID: src,
				Present:    true,
				Impact:     impact,
				Label:      label,
			})
		}
	} else {
		categories := []string{"architecture diagram", "component specification", "data flow description", "security control mapping", "compliance requirement"}
		for _, cat := range categories {
			contribs = append(contribs, EvidenceContribution{
				EvidenceID: cat,
				Present:    false,
				Impact:     -3.0,
				Label:      "Missing",
			})
		}
	}
	return contribs
}

func (e *ExplainabilityEngine) computeDomainContribution(in ConfidenceInput) *DomainContribution {
	if in.Domain == "" || e.domain == "" {
		return nil
	}
	domain := in.Domain
	if domain == "" {
		domain = e.domain
	}
	influence := 0.0
	reason := ""
	strength := StrengthWeak

	switch strings.ToLower(domain) {
	case "healthcare":
		if strings.Contains(strings.ToLower(in.AssumptionText), "phi") ||
			strings.Contains(strings.ToLower(in.AssumptionText), "hipaa") ||
			strings.Contains(strings.ToLower(in.AssumptionText), "patient") ||
			strings.Contains(strings.ToLower(in.AssumptionText), "encrypt") ||
			strings.Contains(strings.ToLower(in.AssumptionText), "audit") ||
			strings.Contains(strings.ToLower(in.AssumptionText), "access") ||
			strings.Contains(strings.ToLower(in.AssumptionText), "key") ||
			strings.Contains(strings.ToLower(in.AssumptionText), "mfa") {
			influence = 12.0
			reason = "PHI requires strict access controls, encryption, and audit logging per HIPAA"
			strength = StrengthStrong
		} else {
			influence = 5.0
			reason = "Healthcare domain context provides partial support"
			strength = StrengthModerate
		}
	case "fintech":
		if strings.Contains(strings.ToLower(in.AssumptionText), "encrypt") ||
			strings.Contains(strings.ToLower(in.AssumptionText), "pci") ||
			strings.Contains(strings.ToLower(in.AssumptionText), "fraud") ||
			strings.Contains(strings.ToLower(in.AssumptionText), "transaction") ||
			strings.Contains(strings.ToLower(in.AssumptionText), "key") {
			influence = 12.0
			reason = "Financial data requires strong encryption and PCI-DSS compliance"
			strength = StrengthStrong
		} else {
			influence = 5.0
			reason = "Fintech domain context provides partial support"
			strength = StrengthModerate
		}
	case "kubernetes", "cloud native":
		if strings.Contains(strings.ToLower(in.AssumptionText), "pod") ||
			strings.Contains(strings.ToLower(in.AssumptionText), "secret") ||
			strings.Contains(strings.ToLower(in.AssumptionText), "rbac") ||
			strings.Contains(strings.ToLower(in.AssumptionText), "network") ||
			strings.Contains(strings.ToLower(in.AssumptionText), "policy") {
			influence = 10.0
			reason = "Kubernetes domain expects pod security, RBAC, and network policies"
			strength = StrengthStrong
		} else {
			influence = 4.0
			reason = "Kubernetes domain context provides partial support"
			strength = StrengthModerate
		}
	default:
		if strings.Contains(strings.ToLower(in.AssumptionText), strings.ToLower(domain)) {
			influence = 8.0
			reason = "Assumption directly references detected domain: " + domain
			strength = StrengthModerate
		} else {
			influence = 2.0
			reason = "General domain context provides baseline support"
			strength = StrengthWeak
		}
	}

	return &DomainContribution{
		Domain:    in.Domain,
		Influence: influence,
		Reason:    reason,
		Strength:  strength,
	}
}

func (e *ExplainabilityEngine) computeTrustContribution(in ConfidenceInput) *TrustContribution {
	return &TrustContribution{
		HasTrustChain:          in.HasTrustChain,
		ChainInfluence:         map[bool]float64{true: 5.0, false: -3.0}[in.HasTrustChain],
		DependencyCentrality:   in.DependencyCentrality,
		FailureRadiusInfluence: float64(in.FailureRadius) * 0.5,
	}
}

func (e *ExplainabilityEngine) classifyStability(confidence float64, posFactors, negFactors []ConfidenceFactor) ConfidenceStability {
	netBias := 0.0
	for _, f := range posFactors {
		netBias += math.Abs(f.Impact)
	}
	for _, f := range negFactors {
		netBias += math.Abs(f.Impact)
	}

	switch {
	case confidence >= 85 && netBias <= 15:
		return StabilityVeryStable
	case confidence >= 70 && netBias <= 25:
		return StabilityStable
	case confidence >= 50:
		return StabilityModerate
	case confidence >= 30:
		return StabilityWeak
	default:
		return StabilityHighlySpeculative
	}
}

func (e *ExplainabilityEngine) stabilityReason(stability ConfidenceStability, confidence float64, in ConfidenceInput) string {
	switch stability {
	case StabilityVeryStable:
		return "High confidence with strong supporting factors and minimal uncertainty"
	case StabilityStable:
		return "Good confidence supported by multiple positive factors"
	case StabilityModerate:
		return "Moderate confidence — some supporting evidence exists but gaps remain"
	case StabilityWeak:
		return "Low confidence — limited supporting evidence, significant gaps present"
	case StabilityHighlySpeculative:
		return "Very low confidence — highly speculative, minimal evidence available"
	default:
		return "Insufficient data to determine stability"
	}
}

func (e *ExplainabilityEngine) generateWhyExists(in ConfidenceInput, factors []ConfidenceFactor) string {
	var parts []string

	if len(in.SourceComponents) > 0 {
		comps := in.SourceComponents
		if len(comps) > 3 {
			comps = comps[:3]
		}
		parts = append(parts, "detected in architecture components: "+strings.Join(comps, ", "))
	}
	if len(in.Keywords) > 0 {
		kw := in.Keywords
		if len(kw) > 5 {
			kw = kw[:5]
		}
		parts = append(parts, "keyword patterns matched: "+strings.Join(kw, ", "))
	}
	if in.Rationale != "" {
		parts = append(parts, in.Rationale)
	}
	if in.Component != "" {
		parts = append(parts, "associated with component: "+in.Component)
	}

	hasStrong := false
	for _, f := range factors {
		if strings.Contains(f.Name, "Explicit Evidence") || strings.Contains(f.Name, "Trust Chain") {
			hasStrong = true
		}
	}
	if hasStrong {
		parts = append(parts, "supported by explicit evidence and trust chain analysis")
	}

	if len(parts) == 0 {
		return "ASF inferred this assumption from architecture patterns and keyword analysis"
	}
	return "ASF identified this assumption because it was " + strings.Join(parts, "; ")
}

func (e *ExplainabilityEngine) generateWhyUncertain(in ConfidenceInput, factors []ConfidenceFactor) string {
	var reasons []string

	for _, f := range factors {
		if f.Type == "negative" {
			reasons = append(reasons, strings.ToLower(f.Description))
		}
	}
	if len(in.EvidenceSources) == 0 {
		reasons = append(reasons, "no direct evidence has been mapped to this assumption")
	}
	if in.VerificationStatus == "" || in.VerificationStatus == "Unverified" {
		reasons = append(reasons, "this assumption has not been verified against actual system state")
	}
	if len(in.SourceComponents) == 0 {
		reasons = append(reasons, "this assumption has no direct traceability to architecture components")
	}

	if len(reasons) == 0 {
		return "ASF is confident in this finding — all available evidence is consistent"
	}

	unique := uniqueStrings(reasons)
	if len(unique) > 3 {
		unique = unique[:3]
	}
	return "ASF has limited confidence because " + strings.Join(unique, ", ")
}

func (e *ExplainabilityEngine) generateWhatIncreases(in ConfidenceInput, factors []ConfidenceFactor) string {
	var suggestions []string

	if len(in.EvidenceSources) == 0 {
		suggestions = append(suggestions, "Map explicit evidence sources to this assumption")
	}
	if in.VerificationStatus == "" || in.VerificationStatus == "Unverified" {
		suggestions = append(suggestions, "Perform verification against actual system configuration")
	}
	if !in.HasTrustChain {
		suggestions = append(suggestions, "Establish a trust chain for this assumption")
	}
	if in.BlindSpotScore > 20 {
		suggestions = append(suggestions, "Address the underlying blind spot to reduce uncertainty")
	}
	if in.Rationale == "" {
		suggestions = append(suggestions, "Document the rationale for why this assumption was generated")
	}
	if len(in.SupportingFactTexts) == 0 {
		suggestions = append(suggestions, "Identify architecture facts that support or refute this assumption")
	}

	if len(suggestions) == 0 {
		return "All confidence-increasing measures have been applied"
	}
	return strings.Join(suggestions, "; ")
}

func (e *ExplainabilityEngine) generateWhatDecreases(in ConfidenceInput, factors []ConfidenceFactor, domainContrib *DomainContribution) string {
	var risks []string

	if domainContrib != nil && domainContrib.Influence > 5 {
		risks = append(risks, "Domain pack influence may over-weight this assumption")
	}
	if len(in.SupportingFactTexts) > 0 {
		risks = append(risks, "Review fact accuracy — incorrect facts would reduce confidence")
	}
	if in.HasTrustChain {
		risks = append(risks, "Trust chain dependencies may introduce hidden assumptions")
	}
	if in.Confidence < 50 {
		risks = append(risks, "Low baseline confidence suggests fundamental uncertainty")
	}
	for _, f := range factors {
		if strings.Contains(f.Name, "Blind Spot") {
			risks = append(risks, "Blind spot indicators suggest possible incomplete analysis")
			break
		}
	}

	if len(risks) == 0 {
		return "No significant confidence-reducing factors identified"
	}
	return strings.Join(risks, "; ")
}

func (e *ExplainabilityEngine) buildCISOView(breakdowns []ConfidenceBreakdown) *CISOTrustView {
	if len(breakdowns) == 0 {
		return nil
	}

	sorted := make([]ConfidenceBreakdown, len(breakdowns))
	copy(sorted, breakdowns)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].FinalConfidence > sorted[j].FinalConfidence
	})

	view := &CISOTrustView{}

	maxTrusted := 5
	if len(sorted) < maxTrusted {
		maxTrusted = len(sorted)
	}
	view.MostTrustedFindings = sorted[:maxTrusted]

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].FinalConfidence < sorted[j].FinalConfidence
	})
	maxLeast := 5
	if len(sorted) < maxLeast {
		maxLeast = len(sorted)
	}
	view.LeastTrustedFindings = sorted[:maxLeast]

	var criticalLow []ConfidenceBreakdown
	for _, b := range breakdowns {
		if b.FinalConfidence < 50 && strings.Contains(strings.ToLower(b.AssumptionText), "critical") {
			criticalLow = append(criticalLow, b)
		}
	}
	if criticalLow == nil {
		criticalLow = []ConfidenceBreakdown{}
	}
	sort.Slice(criticalLow, func(i, j int) bool {
		return criticalLow[i].FinalConfidence < criticalLow[j].FinalConfidence
	})
	maxCritLow := 5
	if len(criticalLow) < maxCritLow {
		maxCritLow = len(criticalLow)
	}
	view.MostCriticalLowConfidence = criticalLow[:maxCritLow]

	var unknowns []ConfidenceBreakdown
	for _, b := range breakdowns {
		if b.FinalConfidence < 30 {
			unknowns = append(unknowns, b)
		}
	}
	if unknowns == nil {
		unknowns = []ConfidenceBreakdown{}
	}
	sort.Slice(unknowns, func(i, j int) bool {
		return unknowns[i].FinalConfidence < unknowns[j].FinalConfidence
	})
	maxUnknowns := 5
	if len(unknowns) < maxUnknowns {
		maxUnknowns = len(unknowns)
	}
	view.HighestRiskUnknowns = unknowns[:maxUnknowns]

	return view
}

func (e *ExplainabilityEngine) buildArchitectView(breakdowns []ConfidenceBreakdown) *ArchitectReviewView {
	if len(breakdowns) == 0 {
		return nil
	}

	view := &ArchitectReviewView{}

	var requiring []ConfidenceBreakdown
	var weak []ConfidenceBreakdown
	var strong []ConfidenceBreakdown

	for _, b := range breakdowns {
		if b.FinalConfidence < 40 {
			requiring = append(requiring, b)
		} else if b.FinalConfidence < 60 {
			weak = append(weak, b)
		} else {
			strong = append(strong, b)
		}
	}

	sort.Slice(requiring, func(i, j int) bool {
		return requiring[i].FinalConfidence < requiring[j].FinalConfidence
	})
	sort.Slice(weak, func(i, j int) bool {
		return weak[i].FinalConfidence < weak[j].FinalConfidence
	})
	sort.Slice(strong, func(i, j int) bool {
		return strong[i].FinalConfidence > strong[j].FinalConfidence
	})

	maxItems := 10
	if len(requiring) > maxItems {
		requiring = requiring[:maxItems]
	}
	if len(weak) > maxItems {
		weak = weak[:maxItems]
	}
	if len(strong) > maxItems {
		strong = strong[:maxItems]
	}
	if requiring == nil {
		requiring = []ConfidenceBreakdown{}
	}
	if weak == nil {
		weak = []ConfidenceBreakdown{}
	}
	if strong == nil {
		strong = []ConfidenceBreakdown{}
	}

	view.RequiringValidation = requiring
	view.WeakSupport = weak
	view.StrongSupport = strong

	return view
}

func uniqueStrings(s []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, v := range s {
		if !seen[v] {
			seen[v] = true
			result = append(result, v)
		}
	}
	return result
}
