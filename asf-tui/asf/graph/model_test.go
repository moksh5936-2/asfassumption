package graph

import (
	"testing"

	"asf-tui/asf/models"
)

func TestBuildGraph(t *testing.T) {
	m := NewModel()
	result := makeTestResult()
	m.Build(result)
	data := m.ExportData()
	if data.NodeCount == 0 {
		t.Error("Expected non-zero nodes")
	}
	if data.EdgeCount == 0 {
		t.Error("Expected non-zero edges")
	}
}

func TestExportJson(t *testing.T) {
	m := NewModel()
	result := makeTestResult()
	m.Build(result)
	data := m.ExportData()

	if len(data.Nodes) != data.NodeCount {
		t.Errorf("Node count mismatch: %d vs %d", len(data.Nodes), data.NodeCount)
	}
	if len(data.Edges) != data.EdgeCount {
		t.Errorf("Edge count mismatch: %d vs %d", len(data.Edges), data.EdgeCount)
	}

	nodeTypes := make(map[string]int)
	for _, n := range data.Nodes {
		nt, ok := n["node_type"].(string)
		if !ok {
			t.Errorf("Node missing node_type: %v", n)
			continue
		}
		nodeTypes[nt]++
	}

	if nodeTypes["Claim"] != 1 {
		t.Errorf("Expected 1 Claim node, got %d", nodeTypes["Claim"])
	}
	if nodeTypes["Assumption"] != 1 {
		t.Errorf("Expected 1 Assumption node, got %d", nodeTypes["Assumption"])
	}
}

func TestBuildFullResult(t *testing.T) {
	m := NewModel()
	r := buildFullResult()
	m.Build(r)
	data := m.ExportData()

	claimNodes := 0
	assumptionNodes := 0
	evidenceNodes := 0
	verificationNodes := 0
	gapNodes := 0

	for _, n := range data.Nodes {
		switch n["node_type"].(string) {
		case "Claim":
			claimNodes++
		case "Assumption":
			assumptionNodes++
		case "Evidence":
			evidenceNodes++
		case "Verification":
			verificationNodes++
		case "Gap":
			gapNodes++
		}
	}

	if claimNodes != 2 {
		t.Errorf("Expected 2 Claim nodes, got %d", claimNodes)
	}
	if assumptionNodes != 2 {
		t.Errorf("Expected 2 Assumption nodes, got %d", assumptionNodes)
	}
	if evidenceNodes != 1 {
		t.Errorf("Expected 1 Evidence node, got %d", evidenceNodes)
	}
	if verificationNodes != 2 {
		t.Errorf("Expected 2 Verification nodes, got %d", verificationNodes)
	}
	if gapNodes != 1 {
		t.Errorf("Expected 1 Gap node, got %d", gapNodes)
	}

	if data.EdgeCount < 5 {
		t.Errorf("Expected at least 5 edges, got %d", data.EdgeCount)
	}
}

func TestGraphRelationships(t *testing.T) {
	m := NewModel()
	r := buildFullResult()
	m.Build(r)
	data := m.ExportData()

	rels := make(map[string]int)
	for _, e := range data.Edges {
		rel, ok := e["relationship"].(string)
		if !ok {
			continue
		}
		rels[rel]++
	}

	if rels["GENERATES"] < 2 {
		t.Errorf("Expected at least 2 GENERATES edges, got %d", rels["GENERATES"])
	}
	if rels["VERIFIES"] < 2 {
		t.Errorf("Expected at least 2 VERIFIES edges, got %d", rels["VERIFIES"])
	}
	if rels["IDENTIFIES"] < 1 {
		t.Errorf("Expected at least 1 IDENTIFIES edge, got %d", rels["IDENTIFIES"])
	}
}

func TestNodeAttributes(t *testing.T) {
	m := NewModel()
	r := buildFullResult()
	m.Build(r)
	data := m.ExportData()

	for _, n := range data.Nodes {
		if _, ok := n["id"]; !ok {
			t.Error("Node missing id")
		}
		if _, ok := n["node_type"]; !ok {
			t.Error("Node missing node_type")
		}
	}
}

func makeTestResult() models.AnalysisResult {
	claim := models.NewClaim("test.txt", "/path/test.txt", "Only admins can access.", 0.6, []string{"access"})
	assumption := models.NewAssumption(claim.ID, "System assumes access control: Only admins can access.", models.AssumptionTypeACCESS, []string{"access"})
	evidence := models.NewEvidence("acl.csv", models.SourceTypeCSV, []map[string]interface{}{
		{"user": "alice", "role": "admin"},
	})
	verification := models.NewVerification(assumption.ID, []string{evidence.ID}, models.VerificationResultVERIFIED, 0.85, "All access OK", nil)
	gap := models.NewGap(assumption.ID, models.GapSeverityLOW, models.GapTypeEVIDENCE, "Minor gap", "details")

	return models.AnalysisResult{
		Claims:        []models.Claim{claim},
		Assumptions:   []models.Assumption{assumption},
		Evidence:      []models.Evidence{evidence},
		Verifications: []models.Verification{verification},
		Gaps:          []models.Gap{gap},
	}
}

func buildFullResult() models.AnalysisResult {
	c1 := models.NewClaim("doc.txt", "doc.txt", "Only admins can access.", 0.7, []string{"access"})
	c2 := models.NewClaim("doc.txt", "doc.txt", "All data is encrypted.", 0.6, []string{"configuration"})

	a1 := models.NewAssumption(c1.ID, "System assumes access control: Only admins can access.", models.AssumptionTypeACCESS, []string{"access"})
	a2 := models.NewAssumption(c2.ID, "System assumes config: All data is encrypted.", models.AssumptionTypeCONFIGURATION, []string{"encrypt"})

	ev := models.NewEvidence("audit.csv", models.SourceTypeCSV, []map[string]interface{}{
		{"user": "alice", "status": "active"},
		{"user": "bob", "status": "active"},
	})

	v1 := models.NewVerification(a1.ID, []string{ev.ID}, models.VerificationResultVERIFIED, 0.85, "Only admins have access", nil)
	v2 := models.NewVerification(a2.ID, []string{ev.ID}, models.VerificationResultCONTRADICTED, 0.72, "Encryption not enabled for all", nil)

	g1 := models.NewGap(a2.ID, models.GapSeverityCRITICAL, models.GapTypeCONFIGURATION, "Critical gap", "Contradicted by evidence")

	return models.AnalysisResult{
		Claims:        []models.Claim{c1, c2},
		Assumptions:   []models.Assumption{a1, a2},
		Evidence:      []models.Evidence{ev},
		Verifications: []models.Verification{v1, v2},
		Gaps:          []models.Gap{g1},
	}
}
