package review

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

type ReviewEngine struct {
	domain string
	inputs []ReviewInput
}

func NewReviewEngine(domain string, inputs []ReviewInput) *ReviewEngine {
	return &ReviewEngine{
		domain: domain,
		inputs: inputs,
	}
}

func (e *ReviewEngine) RunAll() *ReviewOutput {
	output := &ReviewOutput{
		Domain:      e.domain,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
	}

	if len(e.inputs) == 0 {
		return output
	}

	queue := e.buildQueue()
	output.Queue = queue

	output.Matrix = e.buildMatrix(queue.Items)
	output.Campaigns = e.buildCampaigns(queue.Items)
	output.CISODashboard = e.buildCISODashboard(queue.Items)
	output.DomainView = e.buildDomainView(queue.Items)

	return output
}

func (e *ReviewEngine) buildQueue() *ReviewQueue {
	items := make([]ReviewPriority, len(e.inputs))
	for i, in := range e.inputs {
		items[i] = e.scoreItem(in)
	}

	sort.Slice(items, func(i, j int) bool {
		return items[i].PriorityScore > items[j].PriorityScore
	})

	for i := range items {
		items[i].Rank = i + 1
	}

	critCount := 0
	highCount := 0
	medCount := 0
	lowCount := 0
	for _, it := range items {
		switch {
		case it.PriorityScore >= 80:
			critCount++
		case it.PriorityScore >= 60:
			highCount++
		case it.PriorityScore >= 30:
			medCount++
		default:
			lowCount++
		}
	}

	return &ReviewQueue{
		Items:         items,
		TotalItems:    len(items),
		CriticalCount: critCount,
		HighCount:     highCount,
		MediumCount:   medCount,
		LowCount:      lowCount,
	}
}

func (e *ReviewEngine) scoreItem(in ReviewInput) ReviewPriority {
	riskScore := riskWeight(in.Risk)
	trustScore := trustWeight(in.Centrality, in.FailureRadius, in.SupportCount)
	coverageScore := coverageWeight(in.CoverageGap)
	blindSpotScore := blindSpotWeight(in.BlindSpotScore)
	verificationScore := verificationWeight(in.VerificationPriority, in.VerificationConfidence)

	domainBoost := domainWeight(e.domain, in.Category)

	priorityScore := riskScore*0.25 + trustScore*0.20 + coverageScore*0.15 +
		blindSpotScore*0.15 + verificationScore*0.15 + domainBoost*0.10

	priorityScore = math.Round(priorityScore*10) / 10
	if priorityScore > 100 {
		priorityScore = 100
	}
	if priorityScore < 0 {
		priorityScore = 0
	}

	reviewValue := computeReviewValue(riskScore, trustScore, verificationScore, domainBoost)
	reviewEffort := computeReviewEffort(in.Category, in.Component, in.VerificationPriority)
	quadrant := computeQuadrant(reviewValue, reviewEffort)

	whyReview := buildWhyReview(in, priorityScore)
	whatToReview := buildWhatToReview(in.Category, in.Component)
	expectedEvidence := buildExpectedEvidence(in.Category)
	expectedOutcome := buildExpectedOutcome(in, reviewValue)
	expectedRiskReduction := buildRiskReduction(in, priorityScore)
	estimatedTime := buildEstimatedTime(in.Category, in.Component, reviewEffort)

	return ReviewPriority{
		AssumptionID:             in.AssumptionID,
		AssumptionText:           in.AssumptionText,
		Risk:                     in.Risk,
		Category:                 in.Category,
		Component:                in.Component,
		PriorityScore:            priorityScore,
		RiskContribution:         riskScore,
		TrustContribution:        trustScore,
		CoverageContribution:     coverageScore,
		VerificationContribution: verificationScore,
		BlindSpotContribution:    blindSpotScore,
		ReviewValue:              reviewValue,
		ReviewEffort:             reviewEffort,
		Quadrant:                 quadrant,
		WhyReview:                whyReview,
		WhatToReview:             whatToReview,
		ExpectedEvidence:         expectedEvidence,
		ExpectedOutcome:          expectedOutcome,
		ExpectedRiskReduction:    expectedRiskReduction,
		EstimatedTime:            estimatedTime,
	}
}

func (e *ReviewEngine) buildMatrix(items []ReviewPriority) *ReviewMatrix {
	matrix := &ReviewMatrix{}

	for _, it := range items {
		switch it.Quadrant {
		case QuadHighValueLowEffort:
			matrix.HighValueLowEffort = append(matrix.HighValueLowEffort, it)
		case QuadHighValueHighEffort:
			matrix.HighValueHighEffort = append(matrix.HighValueHighEffort, it)
		case QuadLowValueLowEffort:
			matrix.LowValueLowEffort = append(matrix.LowValueLowEffort, it)
		case QuadLowValueHighEffort:
			matrix.LowValueHighEffort = append(matrix.LowValueHighEffort, it)
		}
	}

	sortMatrix := func(s []ReviewPriority) {
		sort.Slice(s, func(i, j int) bool {
			return s[i].PriorityScore > s[j].PriorityScore
		})
	}
	sortMatrix(matrix.HighValueLowEffort)
	sortMatrix(matrix.HighValueHighEffort)
	sortMatrix(matrix.LowValueLowEffort)
	sortMatrix(matrix.LowValueHighEffort)

	return matrix
}

func (e *ReviewEngine) buildCampaigns(items []ReviewPriority) []ReviewCampaign {
	if len(items) == 0 {
		return nil
	}

	campaigns := []ReviewCampaign{
		{Name: "30 Minute Review Plan", Duration: "30 minutes"},
		{Name: "2 Hour Review Plan", Duration: "2 hours"},
		{Name: "1 Day Review Plan", Duration: "1 day"},
		{Name: "1 Week Review Plan", Duration: "1 week"},
	}

	effortOrder := map[ReviewEffort]int{EffortLow: 0, EffortMedium: 1, EffortHigh: 2}

	sorted := make([]ReviewPriority, len(items))
	copy(sorted, items)
	sort.Slice(sorted, func(i, j int) bool {
		pi := effortOrder[sorted[i].ReviewEffort]
		pj := effortOrder[sorted[j].ReviewEffort]
		if pi != pj {
			return pi < pj
		}
		return sorted[i].PriorityScore > sorted[j].PriorityScore
	})

	for ci := range campaigns {
		var selected []ReviewPriority
		accumulated := 0

		for _, it := range sorted {
			already := false
			for _, s := range selected {
				if s.AssumptionID == it.AssumptionID {
					already = true
					break
				}
			}
			if already {
				continue
			}

			effortMin := effortMinutes(it.ReviewEffort)
			if accumulated+effortMin > campaignMinutes(ci) {
				continue
			}
			selected = append(selected, it)
			accumulated += effortMin

			if ci == 0 && len(selected) >= 3 {
				break
			}
		}

		if len(selected) > 0 {
			campaigns[ci].Items = selected
			campaigns[ci].TotalItems = len(selected)
			campaigns[ci].TotalEffort = fmt.Sprintf("%d min", accumulated)
		}
	}

	return campaigns
}

func (e *ReviewEngine) buildCISODashboard(items []ReviewPriority) *CISOReviewDashboard {
	if len(items) == 0 {
		return nil
	}

	topN := 10
	if len(items) < topN {
		topN = len(items)
	}

	highestRisk := make([]ReviewPriority, topN)
	copy(highestRisk, items[:topN])

	var unknowns []ReviewPriority
	for _, it := range items {
		if it.VerificationContribution >= 40 && len(unknowns) < 10 {
			unknowns = append(unknowns, it)
		}
	}

	var mostValuable []ReviewPriority
	for _, it := range items {
		if it.ReviewValue == ValueVeryHigh || it.ReviewValue == ValueHigh {
			if len(mostValuable) < 10 {
				mostValuable = append(mostValuable, it)
			}
		}
	}

	var riskReduction []ReviewPriority
	sorted := make([]ReviewPriority, len(items))
	copy(sorted, items)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].PriorityScore > sorted[j].PriorityScore
	})
	riskReductionN := 10
	if len(sorted) < riskReductionN {
		riskReductionN = len(sorted)
	}
	riskReduction = make([]ReviewPriority, riskReductionN)
	copy(riskReduction, sorted[:riskReductionN])

	critCount := 0
	highCount := 0
	for _, it := range items {
		if it.Risk == "Critical" {
			critCount++
		} else if it.Risk == "High" {
			highCount++
		}
	}

	return &CISOReviewDashboard{
		HighestRiskAssumptions:          highestRisk,
		HighestRiskUnknowns:             unknowns,
		MostValuableVerificationActions: mostValuable,
		GreatestRiskReduction:           riskReduction,
		TotalAssumptions:                len(items),
		CriticalAssumptions:             critCount,
		HighAssumptions:                 highCount,
	}
}

func (e *ReviewEngine) buildDomainView(items []ReviewPriority) *DomainPrioritization {
	if e.domain == "" || e.domain == "general" || len(items) == 0 {
		return nil
	}

	dl := strings.ToLower(e.domain)

	focusAreas := domainFocusAreas(dl)

	var filtered []ReviewPriority
	for _, it := range items {
		cl := strings.ToLower(it.Category)
		for _, fa := range focusAreas {
			if strings.Contains(cl, strings.ToLower(fa)) {
				filtered = append(filtered, it)
				break
			}
		}
	}

	if len(filtered) == 0 {
		filtered = items[:min(5, len(items))]
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].PriorityScore > filtered[j].PriorityScore
	})

	topN := min(10, len(filtered))
	return &DomainPrioritization{
		Domain:        e.domain,
		TopPriorities: filtered[:topN],
		FocusAreas:    focusAreas,
	}
}

func riskWeight(risk string) float64 {
	switch risk {
	case "Critical":
		return 95.0
	case "High":
		return 75.0
	case "Medium":
		return 50.0
	default:
		return 25.0
	}
}

func trustWeight(centrality float64, failureRadius int, supportCount int) float64 {
	score := 0.0
	score += math.Min(centrality*10, 40)
	score += math.Min(float64(failureRadius)*5, 30)
	score += math.Min(float64(supportCount)*3, 30)
	return math.Min(score, 95)
}

func coverageWeight(hasGap bool) float64 {
	if hasGap {
		return 70.0
	}
	return 10.0
}

func blindSpotWeight(score float64) float64 {
	if score <= 0 {
		return 10.0
	}
	return math.Min(score, 95)
}

func verificationWeight(priority string, confidence float64) float64 {
	base := 0.0
	switch priority {
	case "Critical":
		base = 80.0
	case "High":
		base = 60.0
	case "Medium":
		base = 40.0
	default:
		base = 20.0
	}

	if confidence >= 90 {
		base -= 30
	} else if confidence >= 70 {
		base -= 10
	} else if confidence <= 30 {
		base += 20
	}

	return math.Max(base, 0)
}

func domainWeight(domain string, category string) float64 {
	dl := strings.ToLower(domain)
	cl := strings.ToLower(category)

	switch dl {
	case "healthcare", "hipaa":
		if cl == "identity" || cl == "monitoring" {
			return 90.0
		}
		return 70.0
	case "fintech", "financial":
		if cl == "cryptography" || cl == "operational" {
			return 90.0
		}
		return 70.0
	case "cloud", "aws", "azure", "gcp":
		if cl == "identity" || cl == "authorization" {
			return 85.0
		}
		return 65.0
	case "kubernetes", "k8s":
		if cl == "authorization" || cl == "operational" {
			return 90.0
		}
		return 65.0
	default:
		return 50.0
	}
}

func computeReviewValue(risk, trust, verification, domain float64) ReviewValue {
	avg := (risk + trust + verification + domain) / 4.0
	switch {
	case avg >= 80:
		return ValueVeryHigh
	case avg >= 60:
		return ValueHigh
	case avg >= 35:
		return ValueMedium
	default:
		return ValueLow
	}
}

func computeReviewEffort(category, component, verifyPriority string) ReviewEffort {
	cl := strings.ToLower(category)

	if cl == "resilience" || cl == "third_party" {
		return EffortHigh
	}
	if cl == "monitoring" || cl == "operational" {
		if component != "" {
			return EffortMedium
		}
		return EffortLow
	}
	if component != "" {
		return EffortMedium
	}
	if verifyPriority == "Critical" {
		return EffortMedium
	}
	return EffortLow
}

func computeQuadrant(value ReviewValue, effort ReviewEffort) PriorityQuadrant {
	isHighValue := value == ValueVeryHigh || value == ValueHigh
	isLowEffort := effort == EffortLow

	switch {
	case isHighValue && isLowEffort:
		return QuadHighValueLowEffort
	case isHighValue && !isLowEffort:
		return QuadHighValueHighEffort
	case !isHighValue && isLowEffort:
		return QuadLowValueLowEffort
	default:
		return QuadLowValueHighEffort
	}
}

func buildWhyReview(in ReviewInput, score float64) string {
	var reasons []string
	if in.Risk == "Critical" || in.Risk == "High" {
		reasons = append(reasons, fmt.Sprintf("Risk level is %s", in.Risk))
	}
	if in.Centrality >= 0.5 {
		reasons = append(reasons, fmt.Sprintf("Trust centrality is %.1f (supports %d downstream assumptions)", in.Centrality, in.SupportCount))
	}
	if in.FailureRadius >= 3 {
		reasons = append(reasons, fmt.Sprintf("Failure radius affects %d assumptions", in.FailureRadius))
	}
	if in.VerificationPriority == "Critical" || in.VerificationPriority == "High" {
		reasons = append(reasons, fmt.Sprintf("Verification priority is %s", in.VerificationPriority))
	}
	if in.VerificationConfidence < 50 {
		reasons = append(reasons, fmt.Sprintf("Verification confidence is low (%.0f%%)", in.VerificationConfidence))
	}
	if in.CoverageGap {
		reasons = append(reasons, "Sits in a coverage gap area")
	}
	if in.BlindSpotScore > 50 {
		reasons = append(reasons, fmt.Sprintf("Has blind spot severity score of %.0f", in.BlindSpotScore))
	}

	if len(reasons) == 0 {
		reasons = append(reasons, fmt.Sprintf("Priority score is %.0f/100", score))
	}
	return strings.Join(reasons, "; ")
}

func buildWhatToReview(category, component string) string {
	cl := strings.ToLower(category)
	switch cl {
	case "identity":
		return "Identity provider configuration, MFA enforcement, SSO federation, admin access controls"
	case "authorization":
		return "RBAC configuration, role definitions, permission assignments, access review records"
	case "cryptography":
		return "KMS configuration, key rotation policy, encryption standards, certificate management"
	case "monitoring":
		return "SIEM integration, log sources, alert rules, monitoring coverage, retention policies"
	case "resilience":
		return "Backup configuration, restore test results, DR plan, RTO/RPO definitions"
	case "third_party":
		return "Vendor security posture, SOC reports, contract SLAs, integration security"
	case "operational":
		return "Secrets management, CI/CD pipeline security, rate limiting, patch management"
	default:
		return fmt.Sprintf("Security controls and configuration for %s", category)
	}
}

func buildExpectedEvidence(category string) string {
	cl := strings.ToLower(category)
	switch cl {
	case "identity":
		return "Identity policy documents, IdP configuration exports, access log samples"
	case "authorization":
		return "RBAC role matrices, access review reports, permission audit records"
	case "cryptography":
		return "KMS configuration, key rotation logs, certificate inventory, encryption policies"
	case "monitoring":
		return "SIEM configuration, log samples, alert rule definitions, monitoring runbooks"
	case "resilience":
		return "Backup success reports, restore test records, DR exercise documentation"
	case "third_party":
		return "Vendor security assessments, SOC 2 reports, integration architecture diagrams"
	case "operational":
		return "Secrets management configuration, CI/CD pipeline definitions, rate limit configs"
	default:
		return "Relevant security documentation and configuration artifacts"
	}
}

func buildExpectedOutcome(in ReviewInput, value ReviewValue) string {
	base := "Clarify verification status of this assumption"
	if in.VerificationConfidence < 50 {
		base += " and increase confidence from current " + fmt.Sprintf("%.0f%%", in.VerificationConfidence)
	}
	if value == ValueVeryHigh || value == ValueHigh {
		base += "; significant uncertainty removed for critical architecture component"
	}
	return base
}

func buildRiskReduction(in ReviewInput, score float64) string {
	if in.Risk == "Critical" || in.Risk == "High" {
		if in.VerificationConfidence < 50 {
			return "High risk reduction — verifying a critical unverified assumption removes significant architectural uncertainty"
		}
		return "Moderate risk reduction — confirming an already partially verified critical assumption"
	}
	if score >= 60 {
		return "Moderate risk reduction — addressing this priority item improves overall security posture"
	}
	return "Low risk reduction — verifying this assumption provides incremental improvement"
}

func buildEstimatedTime(category, component string, effort ReviewEffort) string {
	if component != "" {
		switch effort {
		case EffortLow:
			return "1-2 hours"
		case EffortMedium:
			return "3-4 hours"
		case EffortHigh:
			return "4-8 hours"
		}
	}

	switch effort {
	case EffortLow:
		return "30 min - 1 hour"
	case EffortMedium:
		return "2-3 hours"
	case EffortHigh:
		return "4-6 hours"
	}
	return "TBD"
}

func domainFocusAreas(domain string) []string {
	switch domain {
	case "healthcare", "hipaa":
		return []string{"identity", "monitoring", "cryptography", "phi", "clinical"}
	case "fintech", "financial":
		return []string{"cryptography", "operational", "monitoring", "settlement", "fraud", "custody"}
	case "cloud", "aws", "azure", "gcp":
		return []string{"identity", "authorization", "cryptography", "federation", "iam"}
	case "kubernetes", "k8s":
		return []string{"authorization", "operational", "identity", "rbac", "admission", "secrets"}
	default:
		return []string{"identity", "authorization", "cryptography", "monitoring"}
	}
}

func effortMinutes(effort ReviewEffort) int {
	switch effort {
	case EffortLow:
		return 30
	case EffortMedium:
		return 120
	case EffortHigh:
		return 240
	default:
		return 60
	}
}

func campaignMinutes(campaignIdx int) int {
	switch campaignIdx {
	case 0:
		return 30
	case 1:
		return 120
	case 2:
		return 480
	case 3:
		return 2400
	default:
		return 480
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
