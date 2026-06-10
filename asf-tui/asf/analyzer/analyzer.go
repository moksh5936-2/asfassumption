package analyzer

import (
	"asf-tui/asf/assumption"
	"asf-tui/asf/confidence"
	"asf-tui/asf/evidence"
	"asf-tui/asf/extraction"
	"asf-tui/asf/gaps"
	"asf-tui/asf/graph"
	"asf-tui/asf/ingestion"
	"asf-tui/asf/models"
	"asf-tui/asf/verification"
)

type Analyzer struct {
	Pipeline           *ingestion.Pipeline
	ClaimExtractor     *extraction.ClaimExtractor
	AssumptionEngine   *assumption.Engine
	EvidenceLoader     *evidence.Loader
	EvidenceMapper     *evidence.Mapper
	VerificationEngine *verification.Engine
	ConfidenceEngine   *confidence.Engine
	GapEngine          *gaps.Engine
	GraphModel         *graph.Model
}

func New() *Analyzer {
	return &Analyzer{
		Pipeline:           ingestion.NewPipeline(),
		ClaimExtractor:     extraction.NewClaimExtractor(),
		AssumptionEngine:   assumption.NewEngine(),
		EvidenceLoader:     evidence.NewLoader(),
		EvidenceMapper:     evidence.NewMapper(),
		VerificationEngine: verification.NewEngine(),
		ConfidenceEngine:   confidence.NewEngine(),
		GapEngine:          gaps.NewEngine(),
		GraphModel:         graph.NewModel(),
	}
}

type AnalyzeResult struct {
	Result models.AnalysisResult
	Graph  graph.GraphData
}

func (a *Analyzer) Analyze(documentPaths []string, evidencePaths []string) (*AnalyzeResult, error) {
	var result models.AnalysisResult

	claims := a.processDocuments(documentPaths)
	assumptions := a.AssumptionEngine.ConvertMany(claims)

	var evidenceRecords []models.Evidence
	if len(evidencePaths) > 0 {
		evidenceRecords = a.loadEvidence(evidencePaths)
	}

	var verifications []models.Verification
	for _, assumption := range assumptions {
		matchingEvidence := a.findMatchingEvidence(assumption, evidenceRecords)
		verification := a.VerificationEngine.Verify(assumption, matchingEvidence)
		verification.Confidence = a.ConfidenceEngine.ComputeVerificationConfidence(verification, matchingEvidence)
		verifications = append(verifications, verification)

		assumption.Confidence = a.ConfidenceEngine.ComputeAssumptionConfidence([]models.Verification{verification})
		switch verification.Result {
		case models.VerificationResultVERIFIED:
			assumption.VerificationStatus = models.VerificationStatusVERIFIED
		case models.VerificationResultCONTRADICTED:
			assumption.VerificationStatus = models.VerificationStatusCONTRADICTED
		case models.VerificationResultPARTIALLY_VERIFIED:
			assumption.VerificationStatus = models.VerificationStatusIN_REVIEW
		}

		for i := range assumptions {
			if assumptions[i].ID == assumption.ID {
				assumptions[i] = assumption
				break
			}
		}
	}

	gapsResult := a.GapEngine.GenerateGaps(assumptions, verifications)

	result = models.AnalysisResult{
		Claims:        claims,
		Assumptions:   assumptions,
		Evidence:      evidenceRecords,
		Verifications: verifications,
		Gaps:          gapsResult,
	}

	a.GraphModel.Build(result)
	graphData := a.GraphModel.ExportData()

	return &AnalyzeResult{
		Result: result,
		Graph:  graphData,
	}, nil
}

func (a *Analyzer) processDocuments(paths []string) []models.Claim {
	var allClaims []models.Claim
	for _, path := range paths {
		text, err := a.Pipeline.ParseText(path)
		if err != nil {
			continue
		}
		meta := a.Pipeline.DetectType(path)
		_ = meta
		filename := path
		for i := len(path) - 1; i >= 0; i-- {
			if path[i] == '/' || path[i] == '\\' {
				filename = path[i+1:]
				break
			}
		}
		claims := a.ClaimExtractor.Extract(text, filename, path)
		allClaims = append(allClaims, claims...)
	}
	return allClaims
}

func (a *Analyzer) loadEvidence(paths []string) []models.Evidence {
	var evidenceList []models.Evidence
	for _, path := range paths {
		ev, err := a.EvidenceLoader.Load(path)
		if err != nil {
			continue
		}
		evidenceList = append(evidenceList, *ev)
	}
	return evidenceList
}

func (a *Analyzer) findMatchingEvidence(assumption models.Assumption, evidenceList []models.Evidence) []models.Evidence {
	compatibleTypes := a.EvidenceMapper.GetCompatibleSourceTypes(assumption.AssumptionType)
	if len(compatibleTypes) == 0 {
		return evidenceList
	}
	var matched []models.Evidence
	for _, ev := range evidenceList {
		for _, ct := range compatibleTypes {
			if ev.SourceType == ct {
				matched = append(matched, ev)
				break
			}
		}
	}
	if len(matched) == 0 {
		return evidenceList
	}
	return matched
}
