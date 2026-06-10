package gaps

import (
	"asf-tui/asf/models"
)

type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

func (ge *Engine) GenerateGaps(assumptions []models.Assumption, verifications []models.Verification) []models.Gap {
	verMap := make(map[string]models.Verification)
	for _, v := range verifications {
		verMap[v.AssumptionID] = v
	}

	var gaps []models.Gap
	for _, a := range assumptions {
		v, hasVerification := verMap[a.ID]

		if !hasVerification {
			gaps = append(gaps, models.NewGap(
				a.ID,
				models.GapSeverityMEDIUM,
				models.GapTypeVERIFICATION,
				truncate("Assumption '"+a.Text+"...' has not been verified", 500),
				"No verification performed",
			))
			continue
		}

		switch v.Result {
		case models.VerificationResultCONTRADICTED:
			gapType := assumptionTypeToGapType(a.AssumptionType)
			severity := determineSeverity(a.AssumptionType, v)
			gaps = append(gaps, models.NewGap(
				a.ID,
				severity,
				gapType,
				truncate("Assumption contradicted: "+a.Text, 500),
				v.Reasoning,
			))

		case models.VerificationResultPARTIALLY_VERIFIED:
			gapType := assumptionTypeToGapType(a.AssumptionType)
			gaps = append(gaps, models.NewGap(
				a.ID,
				models.GapSeverityMEDIUM,
				gapType,
				truncate("Assumption only partially verified: "+a.Text, 500),
				v.Reasoning,
			))

		case models.VerificationResultUNKNOWN:
			gaps = append(gaps, models.NewGap(
				a.ID,
				models.GapSeverityLOW,
				models.GapTypeEVIDENCE,
				truncate("Insufficient evidence to verify: "+a.Text, 500),
				v.Reasoning,
			))
		}
	}

	return gaps
}

func determineSeverity(atype models.AssumptionType, v models.Verification) models.GapSeverity {
	if v.Confidence >= 0.8 {
		switch atype {
		case models.AssumptionTypeACCESS, models.AssumptionTypeIDENTITY, models.AssumptionTypeNETWORK:
			return models.GapSeverityCRITICAL
		case models.AssumptionTypeCONFIGURATION, models.AssumptionTypeGOVERNANCE:
			return models.GapSeverityHIGH
		default:
			return models.GapSeverityHIGH
		}
	}

	if v.Confidence >= 0.5 {
		return models.GapSeverityHIGH
	}

	return models.GapSeverityMEDIUM
}

func assumptionTypeToGapType(atype models.AssumptionType) models.GapType {
	mapping := map[models.AssumptionType]models.GapType{
		models.AssumptionTypeACCESS:       models.GapTypeACCESS,
		models.AssumptionTypeIDENTITY:     models.GapTypeIDENTITY,
		models.AssumptionTypeNETWORK:      models.GapTypeNETWORK,
		models.AssumptionTypeCONFIGURATION: models.GapTypeCONFIGURATION,
		models.AssumptionTypePROCESS:      models.GapTypePROCESS,
		models.AssumptionTypeDOCUMENTATION: models.GapTypeDOCUMENTATION,
		models.AssumptionTypeDEPENDENCY:   models.GapTypeDEPENDENCY,
		models.AssumptionTypeGOVERNANCE:   models.GapTypeGOVERNANCE,
	}
	if gt, ok := mapping[atype]; ok {
		return gt
	}
	return models.GapTypeVERIFICATION
}

func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen])
}
