package review

type ReviewEffort string

const (
	EffortLow    ReviewEffort = "Low"
	EffortMedium ReviewEffort = "Medium"
	EffortHigh   ReviewEffort = "High"
)

type ReviewValue string

const (
	ValueVeryHigh ReviewValue = "Very High"
	ValueHigh     ReviewValue = "High"
	ValueMedium   ReviewValue = "Medium"
	ValueLow      ReviewValue = "Low"
)

type PriorityQuadrant string

const (
	QuadHighValueLowEffort  PriorityQuadrant = "High Value / Low Effort"
	QuadHighValueHighEffort PriorityQuadrant = "High Value / High Effort"
	QuadLowValueLowEffort   PriorityQuadrant = "Low Value / Low Effort"
	QuadLowValueHighEffort  PriorityQuadrant = "Low Value / High Effort"
)

type ReviewInput struct {
	AssumptionID           string
	AssumptionText         string
	Risk                   string
	Category               string
	Component              string
	Centrality             float64
	Criticality            float64
	FailureRadius          int
	SupportCount           int
	DependencyCount        int
	VerificationPriority   string
	VerificationConfidence float64
	VerificationStatus     string
	CoverageGap            bool
	BlindSpotScore         float64
	Domain                 string
}

type ReviewPriority struct {
	AssumptionID             string           `json:"assumption_id"`
	AssumptionText           string           `json:"assumption_text"`
	Risk                     string           `json:"risk"`
	Category                 string           `json:"category"`
	Component                string           `json:"component,omitempty"`
	PriorityScore            float64          `json:"priority_score"`
	RiskContribution         float64          `json:"risk_contribution"`
	TrustContribution        float64          `json:"trust_contribution"`
	CoverageContribution     float64          `json:"coverage_contribution"`
	VerificationContribution float64          `json:"verification_contribution"`
	BlindSpotContribution    float64          `json:"blind_spot_contribution"`
	ReviewValue              ReviewValue      `json:"review_value"`
	ReviewEffort             ReviewEffort     `json:"review_effort"`
	Quadrant                 PriorityQuadrant `json:"quadrant"`
	Rank                     int              `json:"rank"`
	WhyReview                string           `json:"why_review"`
	WhatToReview             string           `json:"what_to_review"`
	ExpectedEvidence         string           `json:"expected_evidence"`
	ExpectedOutcome          string           `json:"expected_outcome"`
	ExpectedRiskReduction    string           `json:"expected_risk_reduction"`
	EstimatedTime            string           `json:"estimated_time"`
}

type ReviewQueue struct {
	Items         []ReviewPriority `json:"items"`
	TotalItems    int              `json:"total_items"`
	CriticalCount int              `json:"critical_count"`
	HighCount     int              `json:"high_count"`
	MediumCount   int              `json:"medium_count"`
	LowCount      int              `json:"low_count"`
}

type ReviewMatrix struct {
	HighValueLowEffort  []ReviewPriority `json:"high_value_low_effort"`
	HighValueHighEffort []ReviewPriority `json:"high_value_high_effort"`
	LowValueLowEffort   []ReviewPriority `json:"low_value_low_effort"`
	LowValueHighEffort  []ReviewPriority `json:"low_value_high_effort"`
}

type ReviewCampaign struct {
	Name        string           `json:"name"`
	Duration    string           `json:"duration"`
	Items       []ReviewPriority `json:"items"`
	TotalItems  int              `json:"total_items"`
	TotalEffort string           `json:"total_effort"`
}

type CISOReviewDashboard struct {
	HighestRiskAssumptions          []ReviewPriority `json:"highest_risk_assumptions"`
	HighestRiskUnknowns             []ReviewPriority `json:"highest_risk_unknowns"`
	MostValuableVerificationActions []ReviewPriority `json:"most_valuable_verification_actions"`
	GreatestRiskReduction           []ReviewPriority `json:"greatest_risk_reduction"`
	TotalAssumptions                int              `json:"total_assumptions"`
	CriticalAssumptions             int              `json:"critical_assumptions"`
	HighAssumptions                 int              `json:"high_assumptions"`
}

type DomainPrioritization struct {
	Domain        string           `json:"domain"`
	TopPriorities []ReviewPriority `json:"top_priorities"`
	FocusAreas    []string         `json:"focus_areas"`
}

type ReviewOutput struct {
	Queue         *ReviewQueue          `json:"queue,omitempty"`
	Matrix        *ReviewMatrix         `json:"matrix,omitempty"`
	Campaigns     []ReviewCampaign      `json:"campaigns,omitempty"`
	CISODashboard *CISOReviewDashboard  `json:"ciso_dashboard,omitempty"`
	DomainView    *DomainPrioritization `json:"domain_view,omitempty"`
	Domain        string                `json:"domain"`
	GeneratedAt   string                `json:"generated_at"`
}
