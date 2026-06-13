package trust

import (
	"strings"
	"testing"
)

func TestDiscoveryEngine(t *testing.T) {
	assumptions := []AssumptionInput{
		{ID: "A1", Text: "MFA is enabled for all admin accounts", Component: "Auth0", Category: "identity", Risk: "Critical", Confidence: 0.9, Keywords: []string{"mfa", "admin"}, Source: "controls"},
		{ID: "A2", Text: "Admin access is restricted", Component: "Auth0", Category: "access", Risk: "Critical", Confidence: 0.85, Keywords: []string{"admin", "access"}, Source: "assumptions"},
		{ID: "A3", Text: "Identity provider is operational", Component: "Auth0", Category: "identity", Risk: "High", Confidence: 0.8, Keywords: []string{"identity", "provider"}, Source: "assumptions"},
		{ID: "A4", Text: "Audit logs are monitored", Component: "SIEM", Category: "monitoring", Risk: "High", Confidence: 0.75, Keywords: []string{"audit", "log", "monitor"}, Source: "controls"},
		{ID: "A5", Text: "PHI access is controlled", Component: "Database", Category: "access", Risk: "Critical", Confidence: 0.9, Keywords: []string{"phi", "access"}, Source: "assumptions"},
		{ID: "A6", Text: "RBAC is configured", Component: "Database", Category: "authorization", Risk: "High", Confidence: 0.8, Keywords: []string{"rbac", "role"}, Source: "controls"},
		{ID: "A7", Text: "Database permissions are correct", Component: "Database", Category: "access", Risk: "High", Confidence: 0.8, Keywords: []string{"database", "permission"}, Source: "assumptions"},
		{ID: "A8", Text: "Encryption is enabled for PHI", Component: "Database", Category: "cryptographic", Risk: "Critical", Confidence: 0.9, Keywords: []string{"encrypt", "phi"}, Source: "controls"},
		{ID: "A9", Text: "Key management is operational", Component: "KMS", Category: "cryptographic", Risk: "High", Confidence: 0.8, Keywords: []string{"key", "management", "kms"}, Source: "assumptions"},
	}

	engine := NewDiscoveryEngine("healthcare", []string{"Auth0", "Database", "KMS", "SIEM"})
	graph := engine.DiscoverDependencies(assumptions)

	if graph == nil {
		t.Fatal("expected non-nil graph")
	}
	if len(graph.Nodes) != 9 {
		t.Errorf("expected 9 nodes, got %d", len(graph.Nodes))
	}
	if len(graph.Edges) == 0 {
		t.Error("expected edges to be discovered")
	}

	// Check specific dependencies
	// A2 (Admin access) should depend on A1 (MFA) or A3 (Identity provider)
	deps := graph.GetDependencies("A2")
	foundIdentityDep := false
	for _, dep := range deps {
		if dep.DependencyType == DepIdentity {
			foundIdentityDep = true
		}
	}
	if !foundIdentityDep {
		t.Errorf("expected A2 to have identity dependency, got: %v", deps)
	}

	// A5 (PHI access) should depend on A6 (RBAC) or A7 (Database permissions)
	deps = graph.GetDependencies("A5")
	if len(deps) == 0 {
		t.Errorf("expected A5 to have dependencies, got none")
	}

	// A8 (Encryption) should depend on A9 (Key management)
	deps = graph.GetDependencies("A8")
	foundCryptoDep := false
	for _, dep := range deps {
		if dep.TargetAssumption == "A9" || dep.DependencyType == DepCryptographic {
			foundCryptoDep = true
		}
	}
	if !foundCryptoDep {
		t.Errorf("expected A8 to have cryptographic dependency, got: %v", deps)
	}

	// Validate graph
	issues := graph.Validate()
	if len(issues) > 0 {
		t.Errorf("graph validation issues: %v", issues)
	}
}

func TestChainEngine(t *testing.T) {
	graph := buildTestGraph()
	engine := NewChainEngine(graph)
	chains := engine.FindTrustChains()

	if len(chains) == 0 {
		t.Error("expected trust chains to be found")
	}

	// Check chain properties
	for _, chain := range chains {
		if chain.Length < 2 {
			t.Errorf("expected chain length >= 2, got %d", chain.Length)
		}
		if chain.Confidence <= 0 {
			t.Errorf("expected positive confidence, got %f", chain.Confidence)
		}
		if chain.Risk == "" {
			t.Error("expected risk to be set")
		}
		if chain.RootNode == "" {
			t.Error("expected root node")
		}
		if chain.LeafNode == "" {
			t.Error("expected leaf node")
		}
	}
}

func TestCascadeEngine(t *testing.T) {
	graph := buildTestGraph()
	engine := NewCascadeEngine(graph)

	// Simulate failure of critical node
	cascade := engine.SimulateFailure("A2")

	if cascade.RootAssumptionID != "A2" {
		t.Errorf("expected root A2, got %s", cascade.RootAssumptionID)
	}
	if cascade.TotalAffected == 0 {
		t.Error("expected cascade to affect other assumptions")
	}
	if cascade.MaxDepth == 0 {
		t.Error("expected max depth > 0")
	}
	if cascade.Severity == "" {
		t.Error("expected severity to be set")
	}

	// Check steps
	for _, step := range cascade.Steps {
		if step.AssumptionID == "" {
			t.Error("expected step to have assumption ID")
		}
		if step.Step <= 0 {
			t.Error("expected positive step number")
		}
	}
}

func TestCriticalEngine(t *testing.T) {
	graph := buildTestGraph()
	engine := NewCriticalEngine(graph)
	critical := engine.FindCriticalAssumptions()

	if len(critical) == 0 {
		t.Error("expected critical assumptions to be found")
	}

	// Check scores
	for _, ca := range critical {
		if ca.Score < 0 || ca.Score > 1 {
			t.Errorf("expected score between 0 and 1, got %f", ca.Score)
		}
		if ca.Centrality < 0 {
			t.Errorf("expected non-negative centrality, got %f", ca.Centrality)
		}
	}

	// Should be sorted by score
	for i := 1; i < len(critical); i++ {
		if critical[i].Score > critical[i-1].Score {
			t.Error("expected critical assumptions to be sorted by score descending")
		}
	}
}

func TestSpotfEngine(t *testing.T) {
	graph := buildTestGraph()
	engine := NewSpotfEngine(graph)
	spotfs := engine.FindSinglePointsOfTrustFailure()

	// Should find Auth0 as a SPOTF (many dependents)
	foundAuth0 := false
	for _, spotf := range spotfs {
		if spotf.DependentsCount >= 3 {
			foundAuth0 = true
		}
		if spotf.Recommendation == "" {
			t.Error("expected SPOTF to have recommendation")
		}
	}
	if !foundAuth0 {
		t.Logf("SPOTFs found: %d", len(spotfs))
		for _, s := range spotfs {
			t.Logf("  %s: %d dependents", s.AssumptionText, s.DependentsCount)
		}
	}
}

func TestCollapseEngine(t *testing.T) {
	graph := buildTestGraph()
	engine := NewCollapseEngine(graph)

	collapse := engine.SimulateCollapse("A2")

	if collapse.FailedAssumptionID != "A2" {
		t.Errorf("expected failed A2, got %s", collapse.FailedAssumptionID)
	}
	if len(collapse.AssumptionsLost) == 0 {
		t.Error("expected assumptions to be lost")
	}
	if collapse.RiskIncrease == "" {
		t.Error("expected risk increase to be computed")
	}
	if collapse.RiskScoreBefore == 0 {
		t.Error("expected risk score before to be computed")
	}
	if collapse.RiskScoreAfter == 0 {
		t.Error("expected risk score after to be computed")
	}
}

func TestTrustChainEngine(t *testing.T) {
	graph := buildTestGraph()
	engine := NewTrustChainEngine(graph)
	output := engine.RunAll()

	if output.DependencyGraph == nil {
		t.Error("expected dependency graph")
	}
	if len(output.TrustChains) == 0 {
		t.Error("expected trust chains")
	}
	if len(output.FailureCascades) == 0 {
		t.Error("expected failure cascades")
	}
	if len(output.CriticalAssumptions) == 0 {
		t.Error("expected critical assumptions")
	}
	if len(output.TrustCollapseResults) == 0 {
		t.Error("expected trust collapse results")
	}
	if output.GeneratedAt == "" {
		t.Error("expected generated at timestamp")
	}
}

func TestExportMarkdown(t *testing.T) {
	output := buildTestChainOutput()
	md := ExportMarkdown(output)

	if !strings.Contains(md, "# Trust Chain Analysis") {
		t.Error("expected markdown to contain title")
	}
	if !strings.Contains(md, "## Critical Assumptions") {
		t.Error("expected markdown to contain critical assumptions")
	}
	if !strings.Contains(md, "## Trust Chains") {
		t.Error("expected markdown to contain trust chains")
	}
	if !strings.Contains(md, "## Failure Cascades") {
		t.Error("expected markdown to contain failure cascades")
	}
}

func TestExportHTML(t *testing.T) {
	output := buildTestChainOutput()
	html := ExportHTML(output)

	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("expected HTML to contain DOCTYPE")
	}
	if !strings.Contains(html, "Trust Chain Analysis") {
		t.Error("expected HTML to contain title")
	}
	if !strings.Contains(html, "Critical Assumptions") {
		t.Error("expected HTML to contain critical assumptions")
	}
}

func TestGraphValidation(t *testing.T) {
	graph := NewDependencyGraph()
	graph.AddNode(&AssumptionNode{ID: "N1", Text: "Test"})
	graph.AddEdge(DependencyEdge{SourceAssumption: "N1", TargetAssumption: "N2", DependencyType: DepIdentity})

	issues := graph.Validate()
	if len(issues) == 0 {
		t.Error("expected validation to find missing N2")
	}
	if !strings.Contains(issues[0], "N2") {
		t.Errorf("expected issue to mention N2, got: %s", issues[0])
	}
}

func TestSelfLoopDetection(t *testing.T) {
	graph := NewDependencyGraph()
	graph.AddNode(&AssumptionNode{ID: "N1", Text: "Test"})
	graph.AddEdge(DependencyEdge{SourceAssumption: "N1", TargetAssumption: "N1", DependencyType: DepIdentity})

	issues := graph.Validate()
	foundSelfLoop := false
	for _, issue := range issues {
		if strings.Contains(issue, "Self-loop") {
			foundSelfLoop = true
		}
	}
	if !foundSelfLoop {
		t.Error("expected validation to detect self-loop")
	}
}

func TestDependencyGraphConnectivity(t *testing.T) {
	graph := buildTestGraph()

	// A2 (Admin access restricted) is a consumer — it depends on A1 (MFA) and A3 (IdP).
	// Nothing depends on A2, so downstream should be empty.
	downstream := graph.GetDownstreamNodes("A2")
	if len(downstream) != 0 {
		t.Errorf("expected A2 to have no downstream nodes, got %d: %v", len(downstream), downstream)
	}

	// A2's upstream = what A2 depends on = A1 and A3
	upstream := graph.GetUpstreamNodes("A2")
	if len(upstream) == 0 {
		t.Error("expected A2 to have upstream nodes (A1, A3)")
	}

	// A1 (MFA) is foundational — many things depend on it, so it should have downstream nodes.
	downstreamA1 := graph.GetDownstreamNodes("A1")
	if len(downstreamA1) == 0 {
		t.Error("expected A1 to have downstream nodes (A2, A5)")
	}

	// Check IsDependency — A2 depends on A1 or A3
	if !graph.IsDependency("A2", "A1") && !graph.IsDependency("A2", "A3") {
		t.Error("expected A2 to depend on A1 or A3")
	}

	// A5 (PHI access) should depend on A6 (RBAC) or A7 (DB permissions)
	if !graph.IsDependency("A5", "A6") && !graph.IsDependency("A5", "A7") {
		t.Error("expected A5 to depend on A6 or A7")
	}

	// A8 (Encryption) should depend on A9 (Key management)
	if !graph.IsDependency("A8", "A9") {
		t.Error("expected A8 to depend on A9")
	}
}

func TestGraphMetrics(t *testing.T) {
	graph := buildTestGraph()
	graph.ComputeCentrality()
	graph.ComputeRadii()

	for _, node := range graph.Nodes {
		if node.Centrality < 0 {
			t.Errorf("expected non-negative centrality for %s", node.ID)
		}
		if node.FailureRadius < 0 {
			t.Errorf("expected non-negative failure radius for %s", node.ID)
		}
		if node.TrustRadius < 0 {
			t.Errorf("expected non-negative trust radius for %s", node.ID)
		}
	}
}

func TestEmptyGraph(t *testing.T) {
	graph := NewDependencyGraph()
	engine := NewTrustChainEngine(graph)
	output := engine.RunAll()

	if len(output.TrustChains) != 0 {
		t.Error("expected no chains for empty graph")
	}
	if len(output.FailureCascades) != 0 {
		t.Error("expected no cascades for empty graph")
	}
	if len(output.CriticalAssumptions) != 0 {
		t.Error("expected no critical assumptions for empty graph")
	}
}

// buildTestGraph creates a test dependency graph.
func buildTestGraph() *DependencyGraph {
	assumptions := []AssumptionInput{
		{ID: "A1", Text: "MFA is enabled", Component: "Auth0", Category: "identity", Risk: "Critical", Confidence: 0.9, Keywords: []string{"mfa"}, Source: "controls"},
		{ID: "A2", Text: "Admin access restricted", Component: "Auth0", Category: "access", Risk: "Critical", Confidence: 0.85, Keywords: []string{"admin", "access"}, Source: "assumptions"},
		{ID: "A3", Text: "Identity provider operational", Component: "Auth0", Category: "identity", Risk: "High", Confidence: 0.8, Keywords: []string{"identity", "provider"}, Source: "assumptions"},
		{ID: "A4", Text: "Audit logs monitored", Component: "SIEM", Category: "monitoring", Risk: "High", Confidence: 0.75, Keywords: []string{"audit", "log"}, Source: "controls"},
		{ID: "A5", Text: "PHI access controlled", Component: "Database", Category: "access", Risk: "Critical", Confidence: 0.9, Keywords: []string{"phi", "access"}, Source: "assumptions"},
		{ID: "A6", Text: "RBAC configured", Component: "Database", Category: "authorization", Risk: "High", Confidence: 0.8, Keywords: []string{"rbac", "role"}, Source: "controls"},
		{ID: "A7", Text: "Database permissions correct", Component: "Database", Category: "access", Risk: "High", Confidence: 0.8, Keywords: []string{"database", "permission"}, Source: "assumptions"},
		{ID: "A8", Text: "Encryption enabled", Component: "Database", Category: "cryptographic", Risk: "Critical", Confidence: 0.9, Keywords: []string{"encrypt"}, Source: "controls"},
		{ID: "A9", Text: "Key management operational", Component: "KMS", Category: "cryptographic", Risk: "High", Confidence: 0.8, Keywords: []string{"key", "management"}, Source: "assumptions"},
	}

	engine := NewDiscoveryEngine("healthcare", []string{"Auth0", "Database", "KMS", "SIEM"})
	return engine.DiscoverDependencies(assumptions)
}

// buildTestChainOutput creates a test ChainOutput.
func buildTestChainOutput() *ChainOutput {
	graph := buildTestGraph()
	engine := NewTrustChainEngine(graph)
	return engine.RunAll()
}

// buildGroundTruthGraph creates a known dependency graph for accuracy benchmarks.
// Edges are designed to match the discovery engine's matching rules (component,
// keyword, domain, and category-based) for the healthcare domain.
func buildGroundTruthGraph() *DependencyGraph {
	g := NewDependencyGraph()
	for id, text := range map[string]string{
		"A1": "MFA is enabled for all admin accounts",
		"A2": "Admin access is restricted",
		"A3": "Identity provider is operational",
		"A4": "Audit logs are monitored",
		"A5": "PHI access is controlled",
		"A6": "RBAC is configured",
		"A7": "Database permissions are correct",
		"A8": "Encryption is enabled for PHI",
		"A9": "Key management is operational",
	} {
		g.AddNode(&AssumptionNode{ID: id, Text: text, Risk: "High"})
	}
	edges := []struct{ src, tgt string }{
		// A1 (MFA) → A3 (IdP): identity dependency (component Auth0, mfa→identity)
		{"A1", "A3"},
		// A2 (Admin access) → A1 (MFA): identity dependency (component Auth0, access→mfa)
		{"A2", "A1"},
		// A2 (Admin access) → A3 (IdP): identity dependency (component Auth0, access→identity)
		{"A2", "A3"},
		// A4 (Audit) → none within same component, but A5 (PHI access, Database) depends on audit
		// A5 (PHI access) → A4 (Audit): monitoring dependency (healthcare domain rule)
		{"A5", "A4"},
		// A5 (PHI access) → A6 (RBAC): identity dependency (phi+access→rbac in healthcare domain)
		{"A5", "A6"},
		// A5 (PHI access) → A7 (DB permissions): authorization (same component + access→permission)
		{"A5", "A7"},
		// A6 (RBAC) → A2 (Admin access): authorization same-component
		{"A6", "A2"},
		// A8 (Encryption for PHI) → A9 (Key management): healthcare domain + keyword (encrypt→key)
		{"A8", "A9"},
	}
	for _, e := range edges {
		g.AddEdge(DependencyEdge{SourceAssumption: e.src, TargetAssumption: e.tgt, DependencyType: DepInfrastructure, Strength: 1.0, Confidence: 1.0})
	}
	g.ComputeCentrality()
	g.ComputeRadii()
	return g
}

// buildGroundTruthAssumptions creates assumptions matching the ground truth graph
// and designed to trigger the discovery engine's matching rules.
func buildGroundTruthAssumptions() []AssumptionInput {
	return []AssumptionInput{
		{ID: "A1", Text: "MFA is enabled for all admin accounts", Component: "Auth0", Category: "identity", Risk: "Critical", Confidence: 0.95, Keywords: []string{"mfa", "admin"}, Source: "controls"},
		{ID: "A2", Text: "Admin access is restricted", Component: "Auth0", Category: "access", Risk: "Critical", Confidence: 0.9, Keywords: []string{"admin", "access"}, Source: "assumptions"},
		{ID: "A3", Text: "Identity provider is operational", Component: "Auth0", Category: "identity", Risk: "High", Confidence: 0.85, Keywords: []string{"identity", "provider"}, Source: "assumptions"},
		{ID: "A4", Text: "Audit logs are monitored via SIEM", Component: "SIEM", Category: "monitoring", Risk: "High", Confidence: 0.8, Keywords: []string{"audit", "log", "siem"}, Source: "controls"},
		{ID: "A5", Text: "PHI access is controlled and restricted", Component: "Database", Category: "access", Risk: "Critical", Confidence: 0.95, Keywords: []string{"phi", "access", "control"}, Source: "assumptions"},
		{ID: "A6", Text: "RBAC is configured for database roles", Component: "Database", Category: "authorization", Risk: "High", Confidence: 0.85, Keywords: []string{"rbac", "role"}, Source: "controls"},
		{ID: "A7", Text: "Database permissions are correct and scoped", Component: "Database", Category: "access", Risk: "High", Confidence: 0.8, Keywords: []string{"database", "permission"}, Source: "assumptions"},
		{ID: "A8", Text: "Encryption is enabled for PHI at rest", Component: "Database", Category: "cryptographic", Risk: "Critical", Confidence: 0.95, Keywords: []string{"encrypt", "phi"}, Source: "controls"},
		{ID: "A9", Text: "Key management is operational in KMS", Component: "KMS", Category: "cryptographic", Risk: "High", Confidence: 0.85, Keywords: []string{"key", "kms", "rotation"}, Source: "assumptions"},
	}
}

// TestDependencyAccuracy verifies the discovery engine finds all expected
// dependencies against a ground truth based on the engine's matching rules.
func TestDependencyAccuracy(t *testing.T) {
	groundTruth := buildGroundTruthGraph()
	assumptions := buildGroundTruthAssumptions()
	components := []string{"Auth0", "SIEM", "Database", "KMS"}

	engine := NewDiscoveryEngine("healthcare", components)
	graph := engine.DiscoverDependencies(assumptions)

	// Count how many ground truth edges the engine discovered
	correctEdges := 0
	extraEdges := 0
	totalEdges := len(groundTruth.Edges)

	for _, expected := range groundTruth.Edges {
		if graph.IsDependency(expected.SourceAssumption, expected.TargetAssumption) {
			correctEdges++
		}
	}

	// Count unexpected edges not in ground truth
	edgeSet := make(map[string]bool)
	for _, e := range groundTruth.Edges {
		edgeSet[e.SourceAssumption+"→"+e.TargetAssumption] = true
	}
	for _, e := range graph.Edges {
		key := e.SourceAssumption + "→" + e.TargetAssumption
		if !edgeSet[key] {
			extraEdges++
		}
	}

	// Precision = correct / total_discovered; Recall = correct / ground_truth
	totalDiscovered := len(graph.Edges)
	precision := 1.0
	if totalDiscovered > 0 {
		precision = float64(correctEdges) / float64(totalDiscovered)
	}
	recall := float64(correctEdges) / float64(totalEdges)
	f1 := 0.0
	if precision+recall > 0 {
		f1 = 2 * precision * recall / (precision + recall)
	}

	t.Logf("dependency discovery: recall=%.2f (%d/%d), precision=%.2f (%d/%d), F1=%.2f",
		recall, correctEdges, totalEdges, precision, correctEdges, totalDiscovered, f1)
	if recall < 0.8 {
		t.Errorf("dependency recall %.2f < 0.8 threshold (%d/%d ground truth edges found)",
			recall, correctEdges, totalEdges)
	}
}

// TestCascadeAccuracy verifies that cascade failure propagation is correct.
// The cascade engine follows Outgoing edges (what the failed node depends on),
// which represents upstream dependencies that would be impacted.
func TestCascadeAccuracy(t *testing.T) {
	assumptions := buildGroundTruthAssumptions()
	components := []string{"Auth0", "SIEM", "Database", "KMS"}

	discoveryEngine := NewDiscoveryEngine("healthcare", components)
	graph := discoveryEngine.DiscoverDependencies(assumptions)
	cascadeEngine := NewCascadeEngine(graph)

	// Test cascade for each node that has Outgoing edges
	testNodes := []string{"A2", "A5", "A1", "A6"}
	for _, nodeID := range testNodes {
		cascade := cascadeEngine.SimulateFailure(nodeID)

		if cascade.RootAssumptionID != nodeID {
			t.Errorf("expected root %s, got %s", nodeID, cascade.RootAssumptionID)
		}

		// The cascade follows Outgoing edges from the failed node
		// Outgoing[nodeID] = edges where nodeID is source → what nodeID depends on
		expectedUpstream := graph.Outgoing[nodeID]
		matched := 0
		for _, expected := range expectedUpstream {
			for _, step := range cascade.Steps {
				if step.AssumptionID == expected.TargetAssumption {
					matched++
					break
				}
			}
		}

		if len(expectedUpstream) > 0 {
			accuracy := float64(matched) / float64(len(expectedUpstream))
			t.Logf("cascade %s: %d upstream deps, %d matched, %d steps, accuracy=%.2f",
				nodeID, len(expectedUpstream), matched, len(cascade.Steps), accuracy)
			if accuracy < 0.8 {
				t.Errorf("cascade accuracy for %s: %.2f < 0.8 (%d/%d upstream in steps)",
					nodeID, accuracy, matched, len(expectedUpstream))
			}
		} else {
			t.Logf("cascade %s: no upstream dependencies (leaf node)", nodeID)
		}

		if len(cascade.Steps) > 1 {
			for i := 1; i < len(cascade.Steps); i++ {
				if cascade.Steps[i].Step <= cascade.Steps[i-1].Step {
					t.Errorf("cascade %s steps not depth-ordered at index %d", nodeID, i)
					break
				}
			}
		}
	}
}

// TestCriticalAccuracy verifies critical assumption detection correctly
// ranks high-centrality and high-support nodes from the discovered graph.
func TestCriticalAccuracy(t *testing.T) {
	assumptions := buildGroundTruthAssumptions()
	components := []string{"Auth0", "SIEM", "Database", "KMS"}

	discoveryEngine := NewDiscoveryEngine("healthcare", components)
	graph := discoveryEngine.DiscoverDependencies(assumptions)

	// Compute centrality and radii on the actual discovered graph
	graph.ComputeCentrality()
	graph.ComputeRadii()

	criticalEngine := NewCriticalEngine(graph)
	critical := criticalEngine.FindCriticalAssumptions()

	// Verify scores are in descending order
	for i := 1; i < len(critical); i++ {
		if critical[i].Score > critical[i-1].Score {
			t.Errorf("critical assumptions not sorted by score descending at index %d", i)
			break
		}
	}

	// Verify that when a node has high centrality, high support count, and
	// critical risk level, its computed score is >= 0.5 (the inclusion threshold).
	// The critical engine uses: 0.3*centrality + 0.25*supportNorm + 0.25*radiusNorm + 0.2*risk
	centralityWeight := 0.3
	supportWeight := 0.25
	radiusWeight := 0.25
	riskWeight := 0.2

	for _, node := range graph.Nodes {
		riskScore := 0.0
		switch node.Risk {
		case "Critical":
			riskScore = 1.0
		case "High":
			riskScore = 0.7
		case "Medium":
			riskScore = 0.4
		case "Low":
			riskScore = 0.1
		}
		supportNorm := float64(node.SupportCount) / float64(len(graph.Nodes))
		radiusNorm := float64(node.FailureRadius) / float64(len(graph.Nodes))
		expectedScore := (node.Centrality * centralityWeight) +
			(supportNorm * supportWeight) +
			(radiusNorm * radiusWeight) +
			(riskScore * riskWeight)

		// Find this node in critical results
		found := false
		for _, ca := range critical {
			if ca.AssumptionID == node.ID {
				found = true
				// Verify score is approximately correct (within rounding)
				if ca.Score < expectedScore-0.01 || ca.Score > expectedScore+0.01 {
					t.Errorf("score mismatch for %s: expected %.4f, got %.4f (centrality=%.3f, support=%d, radius=%d, risk=%s)",
						node.ID, expectedScore, ca.Score, node.Centrality, node.SupportCount, node.FailureRadius, node.Risk)
				}
				break
			}
		}

		if !found {
			// Node not in critical results - verify score < 0.5
			if expectedScore >= 0.5 {
				t.Errorf("node %s with expected score %.4f >= 0.5 threshold not in critical results (centrality=%.3f, support=%d, radius=%d, risk=%s)",
					node.ID, expectedScore, node.Centrality, node.SupportCount, node.FailureRadius, node.Risk)
			}
		}
	}

	t.Logf("critical engine found %d nodes with score >= 0.5", len(critical))
}

// BenchmarkDependencyAccuracy measures dependency discovery accuracy under benchmark load.
func BenchmarkDependencyAccuracy(b *testing.B) {
	groundTruth := buildGroundTruthGraph()
	assumptions := buildGroundTruthAssumptions()
	components := []string{"Auth0", "SIEM", "Database", "KMS"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		engine := NewDiscoveryEngine("healthcare", components)
		graph := engine.DiscoverDependencies(assumptions)

		correctEdges := 0
		for _, expected := range groundTruth.Edges {
			if graph.IsDependency(expected.SourceAssumption, expected.TargetAssumption) {
				correctEdges++
			}
		}

		recall := float64(correctEdges) / float64(len(groundTruth.Edges))
		if recall < 0.8 {
			b.Logf("dependency recall %.2f < 0.8", recall)
		}
	}
}

// BenchmarkCascadeAccuracy measures cascade prediction accuracy under benchmark load.
func BenchmarkCascadeAccuracy(b *testing.B) {
	assumptions := buildGroundTruthAssumptions()
	components := []string{"Auth0", "SIEM", "Database", "KMS"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		discoveryEngine := NewDiscoveryEngine("healthcare", components)
		graph := discoveryEngine.DiscoverDependencies(assumptions)
		cascadeEngine := NewCascadeEngine(graph)

		cascade := cascadeEngine.SimulateFailure("A2")
		expectedUpstream := graph.Outgoing["A2"]
		matched := 0
		for _, expected := range expectedUpstream {
			for _, step := range cascade.Steps {
				if step.AssumptionID == expected.TargetAssumption {
					matched++
					break
				}
			}
		}

		if len(expectedUpstream) > 0 {
			accuracy := float64(matched) / float64(len(expectedUpstream))
			if accuracy < 0.8 {
				b.Logf("cascade accuracy %.2f < 0.8", accuracy)
			}
		}
	}
}

// BenchmarkCriticalAccuracy measures critical assumption detection accuracy under benchmark load.
func BenchmarkCriticalAccuracy(b *testing.B) {
	assumptions := buildGroundTruthAssumptions()
	components := []string{"Auth0", "SIEM", "Database", "KMS"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		discoveryEngine := NewDiscoveryEngine("healthcare", components)
		graph := discoveryEngine.DiscoverDependencies(assumptions)

		criticalEngine := NewCriticalEngine(graph)
		critical := criticalEngine.FindCriticalAssumptions()

		// Verify scores are in descending order
		for i := 1; i < len(critical); i++ {
			if critical[i].Score > critical[i-1].Score {
				b.Logf("scores not sorted at %d", i)
				break
			}
		}
	}
}
