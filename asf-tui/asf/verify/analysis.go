package verify

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

type VerificationInput struct {
	ID          string
	Description string
	Component   string
	Category    string
	Risk        string
	Keywords    []string
}

type VerificationEngine struct {
	domain      string
	components  []string
	assumptions []VerificationInput
}

func NewVerificationEngine(domain string, components []string, assumptions []VerificationInput) *VerificationEngine {
	return &VerificationEngine{
		domain:      domain,
		components:  components,
		assumptions: assumptions,
	}
}

func (e *VerificationEngine) RunAll() *VerificationOutput {
	output := &VerificationOutput{
		Domain:      e.domain,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
	}

	assessment := e.buildAssessment()
	output.Assessment = assessment
	output.Roadmaps = e.buildRoadmaps(assessment.Plans)
	output.CISOView = e.buildCISOView(assessment)

	return output
}

func (e *VerificationEngine) buildAssessment() *VerificationAssessment {
	plans := e.createPlans()

	verifiedCount := 0
	partialCount := 0
	unverifiedCount := 0
	noEvidenceCount := 0

	for _, p := range plans {
		switch p.Status {
		case VsVerified:
			verifiedCount++
		case VsPartiallyVerified:
			partialCount++
		case VsUnverified:
			unverifiedCount++
		case VsNoEvidence:
			noEvidenceCount++
		}
	}

	total := len(plans)
	overallConfidence := 0.0
	if total > 0 {
		totalConf := 0.0
		for _, p := range plans {
			totalConf += p.Confidence
		}
		overallConfidence = math.Round(totalConf/float64(total)*10) / 10
	}

	return &VerificationAssessment{
		Plans:             plans,
		VerifiedCount:     verifiedCount,
		PartialCount:      partialCount,
		UnverifiedCount:   unverifiedCount,
		NoEvidenceCount:   noEvidenceCount,
		TotalAssumptions:  total,
		OverallConfidence: overallConfidence,
	}
}

func (e *VerificationEngine) createPlans() []VerificationPlan {
	var plans []VerificationPlan

	for _, a := range e.assumptions {
		cat := EvidenceCategory(toLower(a.Category))
		keywords := a.Keywords
		if len(keywords) == 0 {
			keywords = extractKeywords(a.Description)
		}

		evidenceReq := lookupEvidence(keywords, cat)
		actions := lookupActions(keywords, cat)
		stakeholders := lookupStakeholders(keywords, cat)
		whyVerify := lookupWhyVerify(keywords, cat)
		whatToReview := lookupWhatToReview(keywords, cat)
		whatEvidence := lookupWhatEvidence(keywords, cat)
		howToValidate := lookupHowToValidate(keywords, cat)
		expectedTime := lookupExpectedTime(keywords, cat)

		domainEvidence := lookupDomainEvidence(e.domain, cat, keywords)
		if len(domainEvidence) > 0 {
			evidenceReq = append(evidenceReq, domainEvidence...)
		}
		domainActions := lookupDomainActions(e.domain, cat, keywords)
		if len(domainActions) > 0 {
			actions = append(actions, domainActions...)
		}
		domainStakeholders := lookupDomainStakeholders(e.domain, cat, keywords)
		if len(domainStakeholders) > 0 {
			stakeholders = mergeStakeholders(stakeholders, domainStakeholders)
		}
		domainWhyVerify := lookupDomainWhyVerify(e.domain, cat, keywords)
		if domainWhyVerify != "" {
			whyVerify = domainWhyVerify
		}

		if len(evidenceReq) == 0 {
			evidenceReq = defaultEvidence(cat)
		}
		if len(actions) == 0 {
			actions = defaultActions(cat)
		}
		if len(stakeholders) == 0 {
			stakeholders = defaultStakeholders(cat)
		}

		status, confidence := computeConfidence(evidenceReq, nil, a.Risk)
		priority := computePriority(a.Risk, status, cat, e.domain)

		presCount := 0
		for _, ev := range evidenceReq {
			if !ev.Optional {
				presCount++
			}
		}
		var present []EvidenceSource
		var missing []EvidenceSource
		for _, ev := range evidenceReq {
			if !ev.Optional {
				present = append(present, ev)
			} else {
				missing = append(missing, ev)
			}
		}
		if len(present) == 0 {
			present = evidenceReq
			missing = nil
		}

		effort := computeEffort(len(actions), len(evidenceReq), priority)

		plan := VerificationPlan{
			AssumptionID:          a.ID,
			AssumptionText:        a.Description,
			Category:              cat,
			Risk:                  a.Risk,
			EvidenceRequired:      evidenceReq,
			EvidencePresent:       present,
			EvidenceMissing:       missing,
			Confidence:            confidence,
			Priority:              priority,
			Status:                status,
			Effort:                effort,
			Actions:               actions,
			Stakeholders:          stakeholders,
			WhyVerify:             whyVerify,
			WhatToReview:          whatToReview,
			WhatEvidenceToCollect: whatEvidence,
			HowToValidate:         howToValidate,
			ExpectedTime:          expectedTime,
		}
		plans = append(plans, plan)
	}

	return plans
}

func (e *VerificationEngine) buildRoadmaps(plans []VerificationPlan) []VerificationRoadmap {
	priorityOrder := map[VerificationPriority]int{
		VpCritical: 0,
		VpHigh:     1,
		VpMedium:   2,
		VpLow:      3,
	}

	sorted := make([]VerificationPlan, len(plans))
	copy(sorted, plans)
	sort.Slice(sorted, func(i, j int) bool {
		pi := priorityOrder[sorted[i].Priority]
		pj := priorityOrder[sorted[j].Priority]
		if pi != pj {
			return pi < pj
		}
		return sorted[i].Confidence < sorted[j].Confidence
	})

	var roadmaps []VerificationRoadmap
	for _, p := range sorted {
		if len(p.Actions) == 0 {
			continue
		}
		steps := make([]VerificationAction, len(p.Actions))
		for i, a := range p.Actions {
			steps[i] = VerificationAction{
				Step:        i + 1,
				Action:      a.Action,
				Description: a.Description,
				Stakeholder: a.Stakeholder,
			}
		}
		roadmaps = append(roadmaps, VerificationRoadmap{
			AssumptionID:   p.AssumptionID,
			AssumptionText: p.AssumptionText,
			Priority:       p.Priority,
			Steps:          steps,
			Effort:         p.Effort,
			Stakeholders:   p.Stakeholders,
		})
	}

	return roadmaps
}

func (e *VerificationEngine) buildCISOView(assessment *VerificationAssessment) *CISOReviewView {
	if assessment == nil || len(assessment.Plans) == 0 {
		return nil
	}

	priorityOrder := map[VerificationPriority]int{
		VpCritical: 0,
		VpHigh:     1,
		VpMedium:   2,
		VpLow:      3,
	}

	sorted := make([]VerificationPlan, len(assessment.Plans))
	copy(sorted, assessment.Plans)
	sort.Slice(sorted, func(i, j int) bool {
		pi := priorityOrder[sorted[i].Priority]
		pj := priorityOrder[sorted[j].Priority]
		if pi != pj {
			return pi < pj
		}
		return sorted[i].Confidence < sorted[j].Confidence
	})

	topN := 10
	if len(sorted) < topN {
		topN = len(sorted)
	}
	topToVerify := make([]VerificationPlan, topN)
	copy(topToVerify, sorted[:topN])

	var highRiskUnverified []VerificationPlan
	for _, p := range sorted {
		if (p.Risk == "Critical" || p.Risk == "High") &&
			(p.Status == VsUnverified || p.Status == VsNoEvidence) &&
			len(highRiskUnverified) < 10 {
			highRiskUnverified = append(highRiskUnverified, p)
		}
	}

	var evidenceGaps []string
	gapSeen := make(map[string]bool)
	for _, p := range assessment.Plans {
		for _, miss := range p.EvidenceMissing {
			key := fmt.Sprintf("%s: %s", string(p.Category), miss.Name)
			if !gapSeen[key] {
				gapSeen[key] = true
				evidenceGaps = append(evidenceGaps, fmt.Sprintf("%s (%s)", key, p.Priority))
			}
		}
	}

	var backlog []VerificationPlan
	for _, p := range sorted {
		if p.Status == VsUnverified || p.Status == VsNoEvidence {
			backlog = append(backlog, p)
		}
	}
	if len(backlog) > 20 {
		backlog = backlog[:20]
	}

	criticalCount := 0
	highCount := 0
	mediumCount := 0
	lowCount := 0
	for _, p := range assessment.Plans {
		switch p.Priority {
		case VpCritical:
			criticalCount++
		case VpHigh:
			highCount++
		case VpMedium:
			mediumCount++
		case VpLow:
			lowCount++
		}
	}

	return &CISOReviewView{
		TopAssumptionsToVerify: topToVerify,
		HighestRiskUnverified:  highRiskUnverified,
		EvidenceGaps:           evidenceGaps,
		VerificationBacklog:    backlog,
		CriticalCount:          criticalCount,
		HighCount:              highCount,
		MediumCount:            mediumCount,
		LowCount:               lowCount,
	}
}

func computeConfidence(required []EvidenceSource, present []EvidenceSource, risk string) (VerificationStatus, float64) {
	if len(required) == 0 {
		return VsUnverified, 30.0
	}

	requiredCount := 0
	for _, ev := range required {
		if !ev.Optional {
			requiredCount++
		}
	}
	if requiredCount == 0 {
		requiredCount = len(required)
	}

	if requiredCount == 0 {
		return VsUnverified, 0.0
	}

	presentCount := 0
	for _, ev := range required {
		if !ev.Optional {
			presentCount++
		}
	}

	ratio := float64(presentCount) / float64(requiredCount)

	baseConfidence := ratio * 100.0

	riskBonus := 0.0
	switch risk {
	case "Critical":
		riskBonus = -5.0
	case "High":
		riskBonus = -3.0
	case "Medium":
		riskBonus = 0.0
	default:
		riskBonus = 2.0
	}

	confidence := baseConfidence + riskBonus
	if confidence > 100 {
		confidence = 100
	}
	if confidence < 0 {
		confidence = 0
	}

	var status VerificationStatus
	switch {
	case confidence >= 90:
		status = VsVerified
	case confidence >= 70:
		status = VsPartiallyVerified
	case confidence >= 30:
		status = VsUnverified
	default:
		status = VsNoEvidence
	}

	confidence = math.Round(confidence)
	return status, confidence
}

func computePriority(risk string, status VerificationStatus, category EvidenceCategory, domain string) VerificationPriority {
	riskScore := 0
	switch risk {
	case "Critical":
		riskScore = 5
	case "High":
		riskScore = 4
	case "Medium":
		riskScore = 2
	default:
		riskScore = 1
	}

	statusPenalty := 0
	switch status {
	case VsNoEvidence:
		statusPenalty = 3
	case VsUnverified:
		statusPenalty = 2
	case VsPartiallyVerified:
		statusPenalty = 1
	case VsVerified:
		statusPenalty = 0
	}

	total := riskScore + statusPenalty

	if category == EvCatCryptography || category == EvCatIdentity {
		total++
	}

	if domain != "" && domain != "general" {
		total++
	}

	switch {
	case total >= 8:
		return VpCritical
	case total >= 6:
		return VpHigh
	case total >= 4:
		return VpMedium
	default:
		return VpLow
	}
}

func computeEffort(actionCount, evidenceCount int, priority VerificationPriority) VerificationEffort {
	base := actionCount * 2
	base += evidenceCount * 3

	switch {
	case base >= 15:
		return EffortHigh
	case base >= 8:
		return EffortMedium
	default:
		return EffortLow
	}
}

func defaultEvidence(cat EvidenceCategory) []EvidenceSource {
	common := []EvidenceSource{
		{Type: SourcePolicyDocument, Name: "Security Policy", Description: "Security policy documentation", Optional: false},
		{Type: SourceAuditLog, Name: "Audit Records", Description: "Audit log records", Optional: true},
	}
	return common
}

func defaultActions(cat EvidenceCategory) []VerificationAction {
	return []VerificationAction{
		{Step: 1, Action: fmt.Sprintf("Review %s security controls", cat), Description: "Verify security controls are implemented", Stakeholder: "Security Architect"},
		{Step: 2, Action: fmt.Sprintf("Review %s configuration", cat), Description: "Verify configuration follows best practices", Stakeholder: "Security Engineer"},
		{Step: 3, Action: fmt.Sprintf("Review %s audit logs", cat), Description: "Verify audit logs show evidence of control operation", Stakeholder: "Security Engineer"},
	}
}

func defaultStakeholders(cat EvidenceCategory) []string {
	return []string{"Security Architect", "Security Engineer"}
}

func mergeStakeholders(a, b []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range a {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	for _, s := range b {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}

func extractKeywords(description string) []string {
	var keywords []string
	seen := make(map[string]bool)

	words := strings.Fields(toLower(description))
	for _, w := range words {
		w = strings.Trim(w, ".,;:!?\"'()[]{}")
		if len(w) < 3 {
			continue
		}
		stopWords := map[string]bool{
			"the": true, "and": true, "for": true, "are": true, "but": true,
			"not": true, "you": true, "all": true, "any": true, "can": true,
			"had": true, "her": true, "was": true, "one": true, "our": true,
			"out": true, "has": true, "have": true, "been": true, "some": true,
			"same": true, "each": true, "than": true, "that": true, "this": true,
			"very": true, "just": true, "with": true, "from": true, "they": true,
			"also": true, "into": true, "over": true, "such": true,
			"will": true, "would": true, "should": true, "could": true, "must": true,
			"more": true, "most": true, "much": true, "many": true, "well": true,
		}
		if stopWords[w] {
			continue
		}
		if !seen[w] {
			seen[w] = true
			keywords = append(keywords, w)
		}
	}
	return keywords
}
