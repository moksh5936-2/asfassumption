package graph

import (
	"asf-tui/asf/models"
)

type Node struct {
	ID       string                 `json:"id"`
	NodeType string                 `json:"node_type"`
	Attrs    map[string]interface{} `json:"-"`
}

type Edge struct {
	Source string `json:"source"`
	Target string `json:"target"`
	Key    int    `json:"key"`
	Attrs  map[string]interface{} `json:"-"`
}

type GraphData struct {
	Nodes     []map[string]interface{} `json:"nodes"`
	Edges     []map[string]interface{} `json:"edges"`
	NodeCount int                      `json:"node_count"`
	EdgeCount int                      `json:"edge_count"`
}

type Model struct {
	edges []struct {
		src    string
		dst    string
		key    int
		attrs  map[string]interface{}
	}
	nodes map[string]map[string]interface{}
}

func NewModel() *Model {
	return &Model{
		nodes: make(map[string]map[string]interface{}),
	}
}

func (m *Model) Build(result models.AnalysisResult) {
	m.nodes = make(map[string]map[string]interface{})
	m.edges = nil

	for _, claim := range result.Claims {
		label := claim.Text
		if len(label) > 60 {
			label = label[:60] + "..."
		}
		m.nodes[claim.ID] = map[string]interface{}{
			"type":                "Claim",
			"label":               label,
			"full_text":           claim.Text,
			"source_document":     claim.SourceDocument,
			"extraction_confidence": claim.ExtractionConfidence,
		}
	}

	for _, assumption := range result.Assumptions {
		label := assumption.Text
		if len(label) > 60 {
			label = label[:60] + "..."
		}
		m.nodes[assumption.ID] = map[string]interface{}{
			"type":                "Assumption",
			"label":               label,
			"full_text":           assumption.Text,
			"assumption_type":     string(assumption.AssumptionType),
			"verification_status": string(assumption.VerificationStatus),
			"confidence":          assumption.Confidence,
		}
		m.edges = append(m.edges, struct {
			src   string
			dst   string
			key   int
			attrs map[string]interface{}
		}{
			src: assumption.ClaimID,
			dst: assumption.ID,
			key: len(m.edges),
			attrs: map[string]interface{}{"relationship": "GENERATES"},
		})
	}

	for _, ev := range result.Evidence {
		label := ev.Source
		if idx := stringsLastIndex(ev.Source, "/"); idx >= 0 {
			label = ev.Source[idx+1:]
		}
		m.nodes[ev.ID] = map[string]interface{}{
			"type":         "Evidence",
			"label":        label,
			"source_type":  string(ev.SourceType),
			"record_count": len(ev.Records),
			"confidence":   ev.Confidence,
		}
	}

	for _, verification := range result.Verifications {
		m.nodes[verification.ID] = map[string]interface{}{
			"type":   "Verification",
			"label":  "Verification: " + string(verification.Result),
			"result": string(verification.Result),
			"confidence": verification.Confidence,
		}
		m.edges = append(m.edges, struct {
			src   string
			dst   string
			key   int
			attrs map[string]interface{}
		}{
			src: verification.AssumptionID,
			dst: verification.ID,
			key: len(m.edges),
			attrs: map[string]interface{}{"relationship": "VERIFIES"},
		})
		for _, evID := range verification.EvidenceUsed {
			m.edges = append(m.edges, struct {
				src   string
				dst   string
				key   int
				attrs map[string]interface{}
			}{
				src: evID,
				dst: verification.ID,
				key: len(m.edges),
				attrs: map[string]interface{}{"relationship": "SUPPORTS"},
			})
		}

		var assumption *models.Assumption
		for _, a := range result.Assumptions {
			if a.ID == verification.AssumptionID {
				assumptionCopy := a
				assumption = &assumptionCopy
				break
			}
		}
		if assumption != nil {
			var rel string
			switch verification.Result {
			case models.VerificationResultVERIFIED:
				rel = "SUPPORTS"
			case models.VerificationResultCONTRADICTED:
				rel = "CONTRADICTS"
			default:
				rel = "RELATES_TO"
			}
			m.edges = append(m.edges, struct {
				src   string
				dst   string
				key   int
				attrs map[string]interface{}
			}{
				src: verification.ID,
				dst: verification.AssumptionID,
				key: len(m.edges),
				attrs: map[string]interface{}{"relationship": rel},
			})
		}
	}

	for _, gap := range result.Gaps {
		m.nodes[gap.ID] = map[string]interface{}{
			"type":        "Gap",
			"label":       "Gap: " + string(gap.Type) + " (" + string(gap.Severity) + ")",
			"gap_type":    string(gap.Type),
			"severity":    string(gap.Severity),
			"description": gap.Description,
		}
		m.edges = append(m.edges, struct {
			src   string
			dst   string
			key   int
			attrs map[string]interface{}
		}{
			src: gap.AssumptionID,
			dst: gap.ID,
			key: len(m.edges),
			attrs: map[string]interface{}{"relationship": "IDENTIFIES"},
		})
	}
}

func (m *Model) ExportData() GraphData {
	var nodes []map[string]interface{}
	for id, attrs := range m.nodes {
		nodeDict := make(map[string]interface{})
		for k, v := range attrs {
			nodeDict[k] = v
		}
		nodeDict["id"] = id

		if nodeType, ok := nodeDict["type"]; ok {
			nodeDict["node_type"] = nodeType
			delete(nodeDict, "type")
		}

		nodes = append(nodes, nodeDict)
	}

	var edges []map[string]interface{}
	for _, e := range m.edges {
		edgeDict := map[string]interface{}{
			"source": e.src,
			"target": e.dst,
			"key":    e.key,
		}
		for k, v := range e.attrs {
			edgeDict[k] = v
		}
		edges = append(edges, edgeDict)
	}

	return GraphData{
		Nodes:     nodes,
		Edges:     edges,
		NodeCount: len(nodes),
		EdgeCount: len(edges),
	}
}

func stringsLastIndex(s, substr string) int {
	for i := len(s) - len(substr); i >= 0; i-- {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
