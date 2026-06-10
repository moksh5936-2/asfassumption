package models

import (
	"crypto/rand"
	"fmt"
	"time"
)

func generateID(prefix string) string {
	b := make([]byte, 6)
	rand.Read(b)
	return fmt.Sprintf("%s_%x", prefix, b)
}

type Claim struct {
	ID                  string    `json:"id"`
	SourceDocument      string    `json:"source_document"`
	SourceLocation      string    `json:"source_location,omitempty"`
	Text                string    `json:"text"`
	ExtractionConfidence float64  `json:"extraction_confidence"`
	CreatedAt           time.Time `json:"created_at"`
	Tags                []string  `json:"tags"`
}

func NewClaim(sourceDocument, sourceLocation, text string, confidence float64, tags []string) Claim {
	return Claim{
		ID:                  generateID("clm"),
		SourceDocument:      sourceDocument,
		SourceLocation:      sourceLocation,
		Text:                text,
		ExtractionConfidence: confidence,
		CreatedAt:           time.Now().UTC(),
		Tags:                tags,
	}
}

type Assumption struct {
	ID                 string            `json:"id"`
	ClaimID            string            `json:"claim_id"`
	Text               string            `json:"text"`
	AssumptionType     AssumptionType    `json:"assumption_type"`
	VerificationStatus VerificationStatus `json:"verification_status"`
	Confidence         float64           `json:"confidence"`
	CreatedAt          time.Time         `json:"created_at"`
	Keywords           []string          `json:"keywords"`
}

func NewAssumption(claimID, text string, atype AssumptionType, keywords []string) Assumption {
	return Assumption{
		ID:                 generateID("asm"),
		ClaimID:            claimID,
		Text:               text,
		AssumptionType:     atype,
		VerificationStatus: VerificationStatusPENDING,
		Confidence:         0.0,
		CreatedAt:          time.Now().UTC(),
		Keywords:           keywords,
	}
}

type Evidence struct {
	ID         string            `json:"id"`
	Source     string            `json:"source"`
	SourceType SourceType        `json:"source_type"`
	Timestamp  time.Time         `json:"timestamp,omitempty"`
	Content    interface{}       `json:"content,omitempty"`
	Confidence float64           `json:"confidence"`
	Metadata   map[string]interface{} `json:"metadata,omitempty"`
	Records    []map[string]interface{} `json:"records"`
}

func NewEvidence(source string, sourceType SourceType, records []map[string]interface{}) Evidence {
	return Evidence{
		ID:         generateID("evd"),
		Source:     source,
		SourceType: sourceType,
		Timestamp:  time.Now().UTC(),
		Confidence: 0.8,
		Metadata:   make(map[string]interface{}),
		Records:    records,
	}
}

type Verification struct {
	ID            string            `json:"id"`
	AssumptionID  string            `json:"assumption_id"`
	EvidenceUsed  []string          `json:"evidence_used"`
	Result        VerificationResult `json:"result"`
	Confidence    float64           `json:"confidence"`
	Reasoning     string            `json:"reasoning"`
	CreatedAt     time.Time         `json:"created_at"`
	Details       map[string]interface{} `json:"details,omitempty"`
}

func NewVerification(assumptionID string, evidenceUsed []string, result VerificationResult, confidence float64, reasoning string, details map[string]interface{}) Verification {
	return Verification{
		ID:           generateID("vrf"),
		AssumptionID: assumptionID,
		EvidenceUsed: evidenceUsed,
		Result:       result,
		Confidence:   confidence,
		Reasoning:    reasoning,
		CreatedAt:    time.Now().UTC(),
		Details:      details,
	}
}

type Gap struct {
	ID             string     `json:"id"`
	AssumptionID   string     `json:"assumption_id"`
	Severity       GapSeverity `json:"severity"`
	Type           GapType     `json:"type"`
	Description    string     `json:"description"`
	EvidenceDetail string     `json:"evidence_detail"`
	CreatedAt      time.Time  `json:"created_at"`
}

func NewGap(assumptionID string, severity GapSeverity, gtype GapType, description, evidenceDetail string) Gap {
	return Gap{
		ID:             generateID("gap"),
		AssumptionID:   assumptionID,
		Severity:       severity,
		Type:           gtype,
		Description:    description,
		EvidenceDetail: evidenceDetail,
		CreatedAt:      time.Now().UTC(),
	}
}

type AnalysisResult struct {
	Claims        []Claim        `json:"claims"`
	Assumptions   []Assumption   `json:"assumptions"`
	Evidence      []Evidence     `json:"evidence,omitempty"`
	Verifications []Verification `json:"verifications"`
	Gaps          []Gap          `json:"gaps"`
}

func (r *AnalysisResult) ClaimsFound() int {
	return len(r.Claims)
}

func (r *AnalysisResult) AssumptionsFound() int {
	return len(r.Assumptions)
}

func (r *AnalysisResult) VerifiedCount() int {
	count := 0
	for _, v := range r.Verifications {
		if v.Result == VerificationResultVERIFIED {
			count++
		}
	}
	return count
}

func (r *AnalysisResult) ContradictedCount() int {
	count := 0
	for _, v := range r.Verifications {
		if v.Result == VerificationResultCONTRADICTED {
			count++
		}
	}
	return count
}

func (r *AnalysisResult) UnknownCount() int {
	count := 0
	for _, v := range r.Verifications {
		if v.Result == VerificationResultUNKNOWN {
			count++
		}
	}
	return count
}

func (r *AnalysisResult) PartiallyVerifiedCount() int {
	count := 0
	for _, v := range r.Verifications {
		if v.Result == VerificationResultPARTIALLY_VERIFIED {
			count++
		}
	}
	return count
}

func (r *AnalysisResult) CriticalGaps() int {
	count := 0
	for _, g := range r.Gaps {
		if g.Severity == GapSeverityCRITICAL {
			count++
		}
	}
	return count
}

type Summary struct {
	ClaimsFound       int `json:"claims_found"`
	Assumptions       int `json:"assumptions"`
	Verified          int `json:"verified"`
	PartiallyVerified int `json:"partially_verified"`
	Contradicted      int `json:"contradicted"`
	Unknown           int `json:"unknown"`
	CriticalGaps      int `json:"critical_gaps"`
}

func (r *AnalysisResult) BuildSummary() Summary {
	return Summary{
		ClaimsFound:       r.ClaimsFound(),
		Assumptions:       r.AssumptionsFound(),
		Verified:          r.VerifiedCount(),
		PartiallyVerified: r.PartiallyVerifiedCount(),
		Contradicted:      r.ContradictedCount(),
		Unknown:           r.UnknownCount(),
		CriticalGaps:      r.CriticalGaps(),
	}
}
