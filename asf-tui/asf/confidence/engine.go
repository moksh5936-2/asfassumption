package confidence

import (
	"strings"
	"time"

	"asf-tui/asf/models"
)

type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

func (ce *Engine) ComputeVerificationConfidence(verification models.Verification, evidenceRecords []models.Evidence) float64 {
	if verification.Result == models.VerificationResultUNKNOWN {
		return verification.Confidence
	}

	freshness := computeFreshness(evidenceRecords)
	coverage := computeCoverage(verification, evidenceRecords)
	completeness := computeCompleteness(verification)

	weightedScore := freshness*0.3 + coverage*0.4 + completeness*0.3

	baseConfidence := verification.Confidence
	if weightedScore > 0 {
		return baseConfidence * weightedScore
	}
	return baseConfidence * 0.5
}

func (ce *Engine) ComputeAssumptionConfidence(verifications []models.Verification) float64 {
	if len(verifications) == 0 {
		return 0.0
	}

	var sum float64
	for _, v := range verifications {
		sum += v.Confidence
	}
	return sum / float64(len(verifications))
}

func computeFreshness(evidenceRecords []models.Evidence) float64 {
	if len(evidenceRecords) == 0 {
		return 0.3
	}

	now := time.Now()
	var totalScore float64
	for _, ev := range evidenceRecords {
		age := now.Sub(ev.Timestamp)
		switch {
		case age < 24*time.Hour:
			totalScore += 1.0
		case age < 7*24*time.Hour:
			totalScore += 0.9
		case age < 30*24*time.Hour:
			totalScore += 0.7
		case age < 90*24*time.Hour:
			totalScore += 0.5
		case age < 365*24*time.Hour:
			totalScore += 0.3
		default:
			totalScore += 0.1
		}
	}
	return totalScore / float64(len(evidenceRecords))
}

func computeCoverage(verification models.Verification, evidenceRecords []models.Evidence) float64 {
	if len(evidenceRecords) == 0 {
		return 0.2
	}

	usedCount := len(verification.EvidenceUsed)
	if usedCount == 0 {
		return 0.2
	}

	ratio := float64(usedCount) / float64(len(evidenceRecords))
	if ratio > 1.0 {
		ratio = 1.0
	}
	return 0.2 + ratio*0.8
}

func computeCompleteness(verification models.Verification) float64 {
	reasoning := verification.Reasoning
	if reasoning == "" || reasoning == "No evidence processed" {
		return 0.1
	}

	indicatorCount := 0
	indicators := []string{";", ":", "found", "user", "asset", "resource", "MFA", "access"}
	for _, ind := range indicators {
		if strings.Contains(strings.ToLower(reasoning), strings.ToLower(ind)) {
			indicatorCount++
		}
	}

	return 0.3 + (float64(indicatorCount)/float64(len(indicators)))*0.7
}
