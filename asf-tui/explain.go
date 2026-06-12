package main

import (
	"time"
)

// EvidenceSource tracks where evidence for an assumption came from.
type EvidenceSource struct {
	FilePath                string   `json:"file_path"`
	FileType                string   `json:"file_type"`
	MatchedComponents       []string `json:"matched_components"`
	MatchedRelationships    []string `json:"matched_relationships"`
	MatchedTrustBoundaries  []string `json:"matched_trust_boundaries"`
	MatchedSecurityConcepts []string `json:"matched_security_concepts"`
}

// EvidenceSummary is the top-level evidence container on AnalysisResult.
type EvidenceSummary struct {
	TotalSources       int      `json:"total_sources"`
	TotalComponents    int      `json:"total_components"`
	TotalRelationships int      `json:"total_relationships"`
	SourceFiles        []string `json:"source_files"`
}

// StrideJustification explains why a STRIDE category was assigned.
type StrideJustification struct {
	Category           StrideCategory `json:"category"`
	Reason             string         `json:"reason"`
	MatchedRuleIndexes []int          `json:"matched_rule_indexes"`
	MatchedKeywords    []string       `json:"matched_keywords"`
	MatchedComponents  []string       `json:"matched_components"`
	Confidence         float64        `json:"confidence"`
	ConfidenceReason   string         `json:"confidence_reason"`
}

// LikelihoodFactor describes a factor that contributes to likelihood.
type LikelihoodFactor struct {
	Factor string `json:"factor"`
	Value  int    `json:"value"`
	Reason string `json:"reason"`
}

// ImpactFactor describes a factor that contributes to impact.
type ImpactFactor struct {
	Factor string `json:"factor"`
	Value  int    `json:"value"`
	Reason string `json:"reason"`
}

// RiskJustification explains how risk was calculated.
type RiskJustification struct {
	Likelihood        int                `json:"likelihood"`
	LikelihoodReason  string             `json:"likelihood_reason"`
	LikelihoodFactors []LikelihoodFactor `json:"likelihood_factors"`
	Impact            int                `json:"impact"`
	ImpactReason      string             `json:"impact_reason"`
	ImpactFactors     []ImpactFactor     `json:"impact_factors"`
	RiskScore         int                `json:"risk_score"`
	RiskLevel         RiskLevel          `json:"risk_level"`
	RiskReason        string             `json:"risk_reason"`
	Confidence        float64            `json:"confidence"`
	ConfidenceReason  string             `json:"confidence_reason"`
}

// ReviewRecord tracks the human review status of an assumption.
type ReviewRecord struct {
	Status    string    `json:"status"` // Proposed, Accepted, Rejected, Modified
	Notes     string    `json:"notes"`
	Timestamp time.Time `json:"timestamp"`
	Reviewer  string    `json:"reviewer"`
}

// ValidationRecord stores data needed for future precision/recall studies.
type ValidationRecord struct {
	AssumptionID      string           `json:"assumption_id"`
	Description       string           `json:"description"`
	GeneratedEvidence []string         `json:"generated_evidence"`
	AssignedRisk      RiskLevel        `json:"assigned_risk"`
	RiskScore         int              `json:"risk_score"`
	Confidence        float64          `json:"confidence"`
	STRIDECategories  []StrideCategory `json:"stride_categories"`
	ArchReviewResult  string           `json:"arch_review_result"` // Accepted, Rejected, Modified
	ArchNotes         string           `json:"arch_notes"`
	ReviewTimestamp   time.Time        `json:"review_timestamp"`
}

// The following fields are added to the existing Assumption struct via a helper
// that returns the explainability info. The struct itself lives in engine.go.
// We use composition: Assumption + ExplainabilityExtension.

// ExplainabilityExtension holds all explainability data for an assumption.
// This is set alongside the Assumption in the analysis pipeline.
type ExplainabilityExtension struct {
	EvidenceSources      []string              `json:"evidence_sources"`
	SourceComponents     []string              `json:"source_components"`
	SourceRelationships  []string              `json:"source_relationships"`
	Rationale            string                `json:"rationale"`
	StrideJustifications []StrideJustification `json:"stride_justifications"`
	RiskJustification    *RiskJustification    `json:"risk_justification"`
	Review               ReviewRecord          `json:"review"`
}

// AttachExplainability attaches an ExplainabilityExtension to an Assumption
// by setting the extra fields that exist on the struct.
// ControlDetail links a control to the specific assumptions and threats it mitigates.
type ControlDetail struct {
	ID                     string           `json:"id"`
	Description            string           `json:"description"`
	Rationale              string           `json:"rationale"`
	Category               string           `json:"category"`
	MitigatedAssumptionIDs []string         `json:"mitigated_assumption_ids"`
	MitigatedSTRIDE        []StrideCategory `json:"mitigated_stride"`
	Priority               int              `json:"priority"` // 1=highest
}

func attachExplainability(a *Assumption, ext *ExplainabilityExtension) {
	if ext == nil {
		return
	}
	a.EvidenceSources = ext.EvidenceSources
	a.SourceComponents = ext.SourceComponents
	a.SourceRelationships = ext.SourceRelationships
	a.Rationale = ext.Rationale
	a.StrideJustifications = ext.StrideJustifications
	a.RiskJustification = ext.RiskJustification
	a.ReviewStatus = ext.Review.Status
	a.ReviewNotes = ext.Review.Notes
	a.ReviewTimestamp = ext.Review.Timestamp
}
