package trust

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

// ChainEngine performs trust chain analysis.
type ChainEngine struct {
	graph *DependencyGraph
}

// NewChainEngine creates a chain analysis engine.
func NewChainEngine(graph *DependencyGraph) *ChainEngine {
	return &ChainEngine{graph: graph}
}

// FindTrustChains discovers all trust chains in the graph.
func (e *ChainEngine) FindTrustChains() []TrustChain {
	var chains []TrustChain
	seen := make(map[string]bool)

	// Find all root nodes (nodes with no dependents, i.e., no incoming edges)
	roots := e.findRootNodes()

	for _, root := range roots {
		if len(chains) >= maxTrustChains {
			break
		}
		// For each root, find all paths to leaf nodes
		paths := e.findAllPaths(root.ID)
		for _, path := range paths {
			if len(chains) >= maxTrustChains {
				break
			}
			key := strings.Join(path, "->")
			if seen[key] {
				continue
			}
			seen[key] = true

			chain := e.buildChain(path)
			chains = append(chains, chain)
		}
	}

	// Sort by length descending
	sort.Slice(chains, func(i, j int) bool {
		return chains[i].Length > chains[j].Length
	})

	return chains
}

// findRootNodes finds nodes with no incoming edges (no other assumptions depend on them).
func (e *ChainEngine) findRootNodes() []*AssumptionNode {
	var roots []*AssumptionNode
	for _, node := range e.graph.Nodes {
		if len(e.graph.GetDependents(node.ID)) == 0 {
			roots = append(roots, node)
		}
	}
	return roots
}

const (
	maxTrustChains  = 100
	maxChainDepth   = 20
	maxPathsPerRoot = 50
)

// findAllPaths finds all paths from a node to leaf nodes.
func (e *ChainEngine) findAllPaths(startID string) [][]string {
	var paths [][]string

	var dfs func(string, []string)
	dfs = func(current string, path []string) {
		if len(path) > maxChainDepth {
			return
		}
		if len(paths) >= maxPathsPerRoot {
			return
		}
		path = append(path, current)
		outgoing := e.graph.GetDependencies(current)
		if len(outgoing) == 0 {
			// Leaf node
			newPath := make([]string, len(path))
			copy(newPath, path)
			paths = append(paths, newPath)
			return
		}
		for _, edge := range outgoing {
			// Avoid cycles
			found := false
			for _, p := range path {
				if p == edge.TargetAssumption {
					found = true
					break
				}
			}
			if found {
				continue
			}
			dfs(edge.TargetAssumption, path)
		}
	}

	dfs(startID, []string{})
	return paths
}

// buildChain builds a TrustChain from a path.
func (e *ChainEngine) buildChain(path []string) TrustChain {
	chain := TrustChain{
		ID:       fmt.Sprintf("chain-%s-%s", path[0], path[len(path)-1]),
		Nodes:    path,
		Length:   len(path),
		RootNode: path[0],
		LeafNode: path[len(path)-1],
	}

	// Calculate chain confidence as product of node confidences
	confidence := 1.0
	for _, nodeID := range path {
		if node, ok := e.graph.Nodes[nodeID]; ok {
			confidence *= node.Confidence
		}
	}
	chain.Confidence = confidence

	// Chain risk = max risk in chain
	maxRisk := "Low"
	riskOrder := map[string]int{"Critical": 4, "High": 3, "Medium": 2, "Low": 1}
	for _, nodeID := range path {
		if node, ok := e.graph.Nodes[nodeID]; ok {
			if riskOrder[node.Risk] > riskOrder[maxRisk] {
				maxRisk = node.Risk
			}
		}
	}
	chain.Risk = maxRisk
	chain.DependencyCount = len(path) - 1

	return chain
}

// CascadeEngine performs cascade failure analysis.
type CascadeEngine struct {
	graph *DependencyGraph
}

// NewCascadeEngine creates a cascade failure analysis engine.
func NewCascadeEngine(graph *DependencyGraph) *CascadeEngine {
	return &CascadeEngine{graph: graph}
}

// SimulateFailure simulates the failure of a single assumption.
func (e *CascadeEngine) SimulateFailure(nodeID string) FailureCascade {
	node, ok := e.graph.Nodes[nodeID]
	if !ok {
		return FailureCascade{}
	}

	cascade := FailureCascade{
		RootAssumptionID:   nodeID,
		RootAssumptionText: node.Text,
		Severity:           node.Risk,
	}

	// BFS to find all downstream nodes that will fail
	visited := make(map[string]bool)
	queue := []string{nodeID}
	step := 0

	for len(queue) > 0 {
		levelSize := len(queue)
		for i := 0; i < levelSize; i++ {
			current := queue[0]
			queue = queue[1:]

			if current != nodeID {
				if n, ok := e.graph.Nodes[current]; ok && !visited[current] {
					visited[current] = true
					step++
					cascade.Steps = append(cascade.Steps, CascadeResult{
						Step:           step,
						AssumptionID:   current,
						AssumptionText: n.Text,
						Severity:       n.Risk,
						Reason:         fmt.Sprintf("Depends on %s which has failed", nodeID),
					})
				}
			}

			// Find all nodes that depend on current
			for _, edge := range e.graph.Outgoing[current] {
				if !visited[edge.TargetAssumption] {
					queue = append(queue, edge.TargetAssumption)
				}
			}
		}
	}

	cascade.TotalAffected = len(cascade.Steps)
	cascade.MaxDepth = e.computeMaxDepth(nodeID)

	return cascade
}

// computeMaxDepth computes the maximum depth of the cascade.
func (e *CascadeEngine) computeMaxDepth(nodeID string) int {
	maxDepth := 0
	visited := make(map[string]bool)

	var dfs func(string, int)
	dfs = func(current string, depth int) {
		if visited[current] {
			return
		}
		visited[current] = true

		if depth > maxDepth {
			maxDepth = depth
		}
		for _, edge := range e.graph.Outgoing[current] {
			dfs(edge.TargetAssumption, depth+1)
		}
	}
	dfs(nodeID, 0)
	return maxDepth
}

// SimulateAllFailures simulates failure for all critical assumptions.
func (e *CascadeEngine) SimulateAllFailures() []FailureCascade {
	var cascades []FailureCascade

	for _, node := range e.graph.Nodes {
		if node.Risk == "Critical" || node.Risk == "High" || node.FailureRadius > 3 {
			cascade := e.SimulateFailure(node.ID)
			if cascade.TotalAffected > 0 {
				cascades = append(cascades, cascade)
			}
		}
	}

	// Sort by total affected descending
	sort.Slice(cascades, func(i, j int) bool {
		return cascades[i].TotalAffected > cascades[j].TotalAffected
	})

	return cascades
}

// CriticalEngine detects critical assumptions.
type CriticalEngine struct {
	graph *DependencyGraph
}

// NewCriticalEngine creates a critical assumption detection engine.
func NewCriticalEngine(graph *DependencyGraph) *CriticalEngine {
	return &CriticalEngine{graph: graph}
}

// FindCriticalAssumptions finds the most critical assumptions in the graph.
func (e *CriticalEngine) FindCriticalAssumptions() []CriticalAssumptionResult {
	var results []CriticalAssumptionResult

	for _, node := range e.graph.Nodes {
		score := e.computeCriticalityScore(node)
		if score >= 0.5 {
			// Collect dependency types
			depTypes := make(map[DependencyType]bool)
			for _, edge := range e.graph.GetDependencies(node.ID) {
				depTypes[edge.DependencyType] = true
			}
			for _, edge := range e.graph.GetDependents(node.ID) {
				depTypes[edge.DependencyType] = true
			}
			var dtList []DependencyType
			for dt := range depTypes {
				dtList = append(dtList, dt)
			}

			results = append(results, CriticalAssumptionResult{
				AssumptionID:    node.ID,
				AssumptionText:  node.Text,
				Centrality:      node.Centrality,
				SupportCount:    node.SupportCount,
				FailureRadius:   node.FailureRadius,
				TrustRadius:     node.TrustRadius,
				Risk:            node.Risk,
				Score:           score,
				DependencyTypes: dtList,
			})
		}
	}

	// Sort by score descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results
}

// computeCriticalityScore computes a criticality score for a node.
func (e *CriticalEngine) computeCriticalityScore(node *AssumptionNode) float64 {
	// Factors:
	// - Centrality (0.3)
	// - Support count (0.25)
	// - Failure radius (0.25)
	// - Risk level (0.2)

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

	centralityScore := node.Centrality
	supportScore := float64(node.SupportCount) / float64(len(e.graph.Nodes))
	if supportScore > 1.0 {
		supportScore = 1.0
	}
	failureRadiusScore := float64(node.FailureRadius) / float64(len(e.graph.Nodes))
	if failureRadiusScore > 1.0 {
		failureRadiusScore = 1.0
	}

	score := (centralityScore * 0.3) +
		(supportScore * 0.25) +
		(failureRadiusScore * 0.25) +
		(riskScore * 0.2)

	return score
}

// SpotfEngine detects single points of trust failure.
type SpotfEngine struct {
	graph *DependencyGraph
}

// NewSpotfEngine creates a single point of trust failure detection engine.
func NewSpotfEngine(graph *DependencyGraph) *SpotfEngine {
	return &SpotfEngine{graph: graph}
}

// FindSinglePointsOfTrustFailure finds nodes that are single points of trust failure.
func (e *SpotfEngine) FindSinglePointsOfTrustFailure() []SinglePointOfTrustFailure {
	var results []SinglePointOfTrustFailure

	for _, node := range e.graph.Nodes {
		dependents := e.graph.GetDependents(node.ID)
		if len(dependents) >= 3 {
			// Collect dependent nodes
			var depNodes []string
			depTypes := make(map[DependencyType]bool)
			for _, edge := range dependents {
				depNodes = append(depNodes, edge.SourceAssumption)
				depTypes[edge.DependencyType] = true
			}

			// Check if this is a single point of trust failure
			// Criteria: many dependents, high centrality, critical or high risk
			isSPOTF := false
			if node.Risk == "Critical" || node.Risk == "High" {
				isSPOTF = true
			} else if node.Centrality >= 0.3 {
				isSPOTF = true
			} else if len(dependents) >= 5 {
				isSPOTF = true
			}

			if isSPOTF {
				var dtList []DependencyType
				for dt := range depTypes {
					dtList = append(dtList, dt)
				}

				recommendation := e.generateSPOTFRecommendation(node, dtList)

				results = append(results, SinglePointOfTrustFailure{
					NodeID:          node.ID,
					AssumptionText:  node.Text,
					DependentsCount: len(dependents),
					DependentNodes:  depNodes,
					DependencyTypes: dtList,
					Recommendation:  recommendation,
				})
			}
		}
	}

	// Sort by dependents count descending
	sort.Slice(results, func(i, j int) bool {
		return results[i].DependentsCount > results[j].DependentsCount
	})

	return results
}

// generateSPOTFRecommendation generates a recommendation for a SPOTF.
func (e *SpotfEngine) generateSPOTFRecommendation(node *AssumptionNode, depTypes []DependencyType) string {
	// Check if this is a third-party dependency
	for _, dt := range depTypes {
		if dt == DepThirdParty {
			return fmt.Sprintf("Implement redundant %s provider or backup service", node.Text)
		}
	}

	// Check if this is an identity dependency
	for _, dt := range depTypes {
		if dt == DepIdentity {
			return fmt.Sprintf("Implement backup identity provider or federated authentication")
		}
	}

	// Check if this is a cryptographic dependency
	for _, dt := range depTypes {
		if dt == DepCryptographic {
			return fmt.Sprintf("Implement key escrow or backup key management")
		}
	}

	// Check if this is an infrastructure dependency
	for _, dt := range depTypes {
		if dt == DepInfrastructure {
			return fmt.Sprintf("Implement redundant infrastructure or failover")
		}
	}

	return fmt.Sprintf("Implement redundancy or backup for %s", node.Text)
}

// CollapseEngine performs trust collapse simulation.
type CollapseEngine struct {
	graph *DependencyGraph
}

// NewCollapseEngine creates a trust collapse simulation engine.
func NewCollapseEngine(graph *DependencyGraph) *CollapseEngine {
	return &CollapseEngine{graph: graph}
}

// SimulateCollapse simulates the collapse of a critical assumption.
func (e *CollapseEngine) SimulateCollapse(nodeID string) TrustCollapseResult {
	node, ok := e.graph.Nodes[nodeID]
	if !ok {
		return TrustCollapseResult{}
	}

	result := TrustCollapseResult{
		FailedAssumptionID:   nodeID,
		FailedAssumptionText: node.Text,
		RiskScoreBefore:      e.computeRiskScoreBefore(),
	}

	// Find all assumptions that will be lost
	cascadeEngine := NewCascadeEngine(e.graph)
	cascade := cascadeEngine.SimulateFailure(nodeID)
	for _, step := range cascade.Steps {
		result.AssumptionsLost = append(result.AssumptionsLost, step.AssumptionID)
	}

	// Find affected components
	affectedComponents := make(map[string]bool)
	for _, step := range cascade.Steps {
		if n, ok := e.graph.Nodes[step.AssumptionID]; ok && n.Component != "" {
			affectedComponents[n.Component] = true
		}
	}
	for comp := range affectedComponents {
		result.AffectedComponents = append(result.AffectedComponents, comp)
	}

	// Risk increase
	result.RiskScoreAfter = e.computeRiskScoreAfter(nodeID, result.AssumptionsLost)
	if result.RiskScoreAfter > result.RiskScoreBefore {
		result.RiskIncrease = fmt.Sprintf("Risk increased from %.2f to %.2f", result.RiskScoreBefore, result.RiskScoreAfter)
	} else {
		result.RiskIncrease = "Risk unchanged"
	}

	return result
}

// SimulateAllCollapses simulates collapse for all critical assumptions.
func (e *CollapseEngine) SimulateAllCollapses() []TrustCollapseResult {
	var results []TrustCollapseResult

	for _, node := range e.graph.Nodes {
		if node.Risk == "Critical" || node.Risk == "High" || node.Centrality >= 0.3 {
			result := e.SimulateCollapse(node.ID)
			if len(result.AssumptionsLost) > 0 {
				results = append(results, result)
			}
		}
	}

	// Sort by assumptions lost descending
	sort.Slice(results, func(i, j int) bool {
		return len(results[i].AssumptionsLost) > len(results[j].AssumptionsLost)
	})

	return results
}

// computeRiskScoreBefore computes the overall risk score before collapse.
func (e *CollapseEngine) computeRiskScoreBefore() float64 {
	score := 0.0
	for _, node := range e.graph.Nodes {
		switch node.Risk {
		case "Critical":
			score += 4.0
		case "High":
			score += 3.0
		case "Medium":
			score += 2.0
		case "Low":
			score += 1.0
		}
	}
	return score
}

// computeRiskScoreAfter computes the risk score after a collapse.
func (e *CollapseEngine) computeRiskScoreAfter(failedNodeID string, lostAssumptions []string) float64 {
	score := 0.0
	lostSet := make(map[string]bool)
	for _, id := range lostAssumptions {
		lostSet[id] = true
	}
	lostSet[failedNodeID] = true

	for _, node := range e.graph.Nodes {
		if lostSet[node.ID] {
			// Lost assumptions increase risk (they no longer protect)
			switch node.Risk {
			case "Critical":
				score += 8.0
			case "High":
				score += 6.0
			case "Medium":
				score += 4.0
			case "Low":
				score += 2.0
			}
		} else {
			// Remaining assumptions still contribute
			switch node.Risk {
			case "Critical":
				score += 4.0
			case "High":
				score += 3.0
			case "Medium":
				score += 2.0
			case "Low":
				score += 1.0
			}
		}
	}
	return score
}

// TrustChainEngine is the main engine that orchestrates all chain analysis.
type TrustChainEngine struct {
	graph *DependencyGraph
}

// NewTrustChainEngine creates the main trust chain engine.
func NewTrustChainEngine(graph *DependencyGraph) *TrustChainEngine {
	return &TrustChainEngine{graph: graph}
}

// RunAll runs all chain analysis engines.
func (e *TrustChainEngine) RunAll() *ChainOutput {
	output := &ChainOutput{
		DependencyGraph:      e.graph,
		GeneratedAt:          time.Now().Format(time.RFC3339),
		TrustChains:          make([]TrustChain, 0),
		FailureCascades:      make([]FailureCascade, 0),
		CriticalAssumptions:  make([]CriticalAssumptionResult, 0),
		SinglePointsOfTrust:  make([]SinglePointOfTrustFailure, 0),
		TrustCollapseResults: make([]TrustCollapseResult, 0),
	}

	// Trust chains
	chainEngine := NewChainEngine(e.graph)
	output.TrustChains = chainEngine.FindTrustChains()

	// Failure cascades
	cascadeEngine := NewCascadeEngine(e.graph)
	output.FailureCascades = cascadeEngine.SimulateAllFailures()

	// Critical assumptions
	criticalEngine := NewCriticalEngine(e.graph)
	output.CriticalAssumptions = criticalEngine.FindCriticalAssumptions()

	// Single points of trust failure
	spotfEngine := NewSpotfEngine(e.graph)
	output.SinglePointsOfTrust = spotfEngine.FindSinglePointsOfTrustFailure()

	// Trust collapse
	collapseEngine := NewCollapseEngine(e.graph)
	output.TrustCollapseResults = collapseEngine.SimulateAllCollapses()

	return output
}
