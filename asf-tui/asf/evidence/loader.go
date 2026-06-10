package evidence

import (
	"path/filepath"
	"strings"

	"asf-tui/asf/ingestion"
	"asf-tui/asf/models"
)

type Loader struct {
	pipeline *ingestion.Pipeline
}

func NewLoader() *Loader {
	return &Loader{
		pipeline: ingestion.NewPipeline(),
	}
}

func (l *Loader) Load(path string) (*models.Evidence, error) {
	sourceType := l.pipeline.DetectType(path)
	records, err := l.pipeline.ParseToRecords(path)
	if err != nil {
		return nil, err
	}

	ev := models.NewEvidence(filepath.Base(path), sourceType, records)
	return &ev, nil
}

type Mapper struct{}

func NewMapper() *Mapper {
	return &Mapper{}
}

func (m *Mapper) GetCompatibleSourceTypes(atype models.AssumptionType) []models.SourceType {
	switch atype {
	case models.AssumptionTypeACCESS:
		return []models.SourceType{models.SourceTypeCSV, models.SourceTypeJSON, models.SourceTypeACLList}
	case models.AssumptionTypeIDENTITY:
		return []models.SourceType{models.SourceTypeCSV, models.SourceTypeJSON, models.SourceTypeIAMExport}
	case models.AssumptionTypeNETWORK:
		return []models.SourceType{models.SourceTypeCSV, models.SourceTypeJSON, models.SourceTypeFirewallRules, models.SourceTypeSecurityGroups, models.SourceTypeRouteTables}
	case models.AssumptionTypeCONFIGURATION:
		return []models.SourceType{models.SourceTypeCSV, models.SourceTypeJSON, models.SourceTypeConfigExport}
	case models.AssumptionTypePROCESS:
		return []models.SourceType{models.SourceTypeCSV, models.SourceTypeJSON, models.SourceTypeAuditLog}
	case models.AssumptionTypeGOVERNANCE:
		return []models.SourceType{models.SourceTypeCSV, models.SourceTypeJSON, models.SourceTypeAuditLog}
	case models.AssumptionTypeDOCUMENTATION:
		return []models.SourceType{models.SourceTypeCSV, models.SourceTypeJSON}
	case models.AssumptionTypeDEPENDENCY:
		return []models.SourceType{models.SourceTypeCSV, models.SourceTypeJSON}
	default:
		return nil
	}
}

func FindField(record map[string]interface{}, candidates []string) string {
	recLower := make(map[string]string)
	for k, v := range record {
		recLower[strings.ToLower(k)] = strings.ToLower(toString(v))
	}

	for _, candidate := range candidates {
		if val, ok := recLower[strings.ToLower(candidate)]; ok {
			return val
		}
	}
	return ""
}

func toString(v interface{}) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}
