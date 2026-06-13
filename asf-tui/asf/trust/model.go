package trust

import (
	"fmt"
	"strings"
)

// DependencyType categorizes the nature of a dependency between assumptions.
type DependencyType string

const (
	DepIdentity       DependencyType = "IDENTITY"
	DepAuthorization  DependencyType = "AUTHORIZATION"
	DepCryptographic  DependencyType = "CRYPTOGRAPHIC"
	DepMonitoring     DependencyType = "MONITORING"
	DepOperational    DependencyType = "OPERATIONAL"
	DepThirdParty     DependencyType = "THIRD_PARTY"
	DepInfrastructure DependencyType = "INFRASTRUCTURE"
)

// AllDependencyTypes is the full list for iteration.
var AllDependencyTypes = []DependencyType{
	DepIdentity, DepAuthorization, DepCryptographic,
	DepMonitoring, DepOperational, DepThirdParty, DepInfrastructure,
}

// AssumptionNode represents a single assumption in the dependency graph.
type AssumptionNode struct {
	ID          string   `json:"id"`
	Text        string   `json:"text"`
	Type        string   `json:"type"`
	Risk        string   `json:"risk"`
	Confidence  float64  `json:"confidence"`
	Source      string   `json:"source"`
	Criticality float64  `json:"criticality"`
	Component   string   `json:"component,omitempty"`
	Category    string   `json:"category,omitempty"`
	STRIDE      []string `json:"stride,omitempty"`

	// Computed fields
	DependencyCount int     `json:"dependency_count"`
	SupportCount    int     `json:"support_count"`
	FailureRadius   int     `json:"failure_radius"`
	TrustRadius     int     `json:"trust_radius"`
	Centrality      float64 `json:"centrality"`
}

// DependencyEdge represents a directed edge from one assumption to another.
// Source depends on Target. If Target fails, Source may fail.
type DependencyEdge struct {
	SourceAssumption string         `json:"source_assumption"`
	TargetAssumption string         `json:"target_assumption"`
	DependencyType   DependencyType `json:"dependency_type"`
	Strength         float64        `json:"strength"`
	Confidence       float64        `json:"confidence"`
	Reason           string         `json:"reason"`
	IsExplicit       bool           `json:"is_explicit"`
}

// DependencyGraph holds the full graph of assumptions and their dependencies.
type DependencyGraph struct {
	Nodes map[string]*AssumptionNode `json:"nodes"`
	Edges []DependencyEdge           `json:"edges"`

	// Outgoing edges: node ID -> list of edges where this node is the source
	Outgoing map[string][]DependencyEdge `json:"-"`

	// Incoming edges: node ID -> list of edges where this node is the target
	Incoming map[string][]DependencyEdge `json:"-"`
}

// NewDependencyGraph creates an empty dependency graph.
func NewDependencyGraph() *DependencyGraph {
	return &DependencyGraph{
		Nodes:    make(map[string]*AssumptionNode),
		Edges:    make([]DependencyEdge, 0),
		Outgoing: make(map[string][]DependencyEdge),
		Incoming: make(map[string][]DependencyEdge),
	}
}

// AddNode adds an assumption node to the graph.
func (g *DependencyGraph) AddNode(node *AssumptionNode) {
	g.Nodes[node.ID] = node
}

// AddEdge adds a dependency edge to the graph.
func (g *DependencyGraph) AddEdge(edge DependencyEdge) {
	g.Edges = append(g.Edges, edge)
	g.Outgoing[edge.SourceAssumption] = append(g.Outgoing[edge.SourceAssumption], edge)
	g.Incoming[edge.TargetAssumption] = append(g.Incoming[edge.TargetAssumption], edge)

	// Update counts on nodes
	if src, ok := g.Nodes[edge.SourceAssumption]; ok {
		src.DependencyCount++
	}
	if tgt, ok := g.Nodes[edge.TargetAssumption]; ok {
		tgt.SupportCount++
	}
}

// HasEdge checks if a directed edge exists from source to target.
func (g *DependencyGraph) HasEdge(source, target string) bool {
	for _, edge := range g.Outgoing[source] {
		if edge.TargetAssumption == target {
			return true
		}
	}
	return false
}

// GetDependencies returns all edges where the given node is the source (depends on others).
func (g *DependencyGraph) GetDependencies(nodeID string) []DependencyEdge {
	return g.Outgoing[nodeID]
}

// GetDependents returns all edges where the given node is the target (others depend on it).
func (g *DependencyGraph) GetDependents(nodeID string) []DependencyEdge {
	return g.Incoming[nodeID]
}

// GetDownstreamNodes returns all nodes that transitively depend on the given node.
// These are nodes that will be affected if the given node fails.
func (g *DependencyGraph) GetDownstreamNodes(nodeID string) []string {
	visited := make(map[string]bool)
	var result []string

	var dfs func(string)
	dfs = func(current string) {
		for _, edge := range g.Incoming[current] {
			if !visited[edge.SourceAssumption] {
				visited[edge.SourceAssumption] = true
				result = append(result, edge.SourceAssumption)
				dfs(edge.SourceAssumption)
			}
		}
	}

	// Find all nodes that depend on nodeID
	for _, edge := range g.Incoming[nodeID] {
		if !visited[edge.SourceAssumption] {
			visited[edge.SourceAssumption] = true
			result = append(result, edge.SourceAssumption)
			dfs(edge.SourceAssumption)
		}
	}

	return result
}

// GetUpstreamNodes returns all nodes that the given node transitively depends on.
// These are the nodes that this node requires to function.
func (g *DependencyGraph) GetUpstreamNodes(nodeID string) []string {
	visited := make(map[string]bool)
	var result []string

	var dfs func(string)
	dfs = func(current string) {
		for _, edge := range g.Outgoing[current] {
			if !visited[edge.TargetAssumption] {
				visited[edge.TargetAssumption] = true
				result = append(result, edge.TargetAssumption)
				dfs(edge.TargetAssumption)
			}
		}
	}

	for _, edge := range g.Outgoing[nodeID] {
		if !visited[edge.TargetAssumption] {
			visited[edge.TargetAssumption] = true
			result = append(result, edge.TargetAssumption)
			dfs(edge.TargetAssumption)
		}
	}

	return result
}

// IsDependency checks if source depends on target (direct or transitive).
func (g *DependencyGraph) IsDependency(sourceID, targetID string) bool {
	visited := make(map[string]bool)

	var dfs func(string) bool
	dfs = func(current string) bool {
		for _, edge := range g.Outgoing[current] {
			if edge.TargetAssumption == targetID {
				return true
			}
			if !visited[edge.TargetAssumption] {
				visited[edge.TargetAssumption] = true
				if dfs(edge.TargetAssumption) {
					return true
				}
			}
		}
		return false
	}

	visited[sourceID] = true
	return dfs(sourceID)
}

// ComputeCentrality computes betweenness-like centrality for each node.
// Nodes that appear on many dependency paths are more central.
func (g *DependencyGraph) ComputeCentrality() {
	for _, node := range g.Nodes {
		// Count how many nodes depend on this node
		dependents := g.GetDependents(node.ID)
		// Count how many nodes this node depends on
		dependencies := g.GetDependencies(node.ID)
		// Centrality = (dependents + dependencies) / total nodes
		node.Centrality = float64(len(dependents)+len(dependencies)) / float64(len(g.Nodes))
	}
}

// ComputeRadii computes trust and failure radius for each node.
// Trust radius = upstream nodes (what this node depends on)
// Failure radius = downstream nodes (what depends on this node)
func (g *DependencyGraph) ComputeRadii() {
	for _, node := range g.Nodes {
		node.TrustRadius = len(g.GetUpstreamNodes(node.ID))
		node.FailureRadius = len(g.GetDownstreamNodes(node.ID))
	}
}

// Validate checks the graph for consistency.
func (g *DependencyGraph) Validate() []string {
	var issues []string

	for _, edge := range g.Edges {
		if _, ok := g.Nodes[edge.SourceAssumption]; !ok {
			issues = append(issues, fmt.Sprintf("Edge references unknown source: %s", edge.SourceAssumption))
		}
		if _, ok := g.Nodes[edge.TargetAssumption]; !ok {
			issues = append(issues, fmt.Sprintf("Edge references unknown target: %s", edge.TargetAssumption))
		}
	}

	// Check for self-loops
	for _, edge := range g.Edges {
		if edge.SourceAssumption == edge.TargetAssumption {
			issues = append(issues, fmt.Sprintf("Self-loop detected: %s", edge.SourceAssumption))
		}
	}

	return issues
}

// Summary returns a string summary of the graph.
func (g *DependencyGraph) Summary() string {
	var parts []string
	parts = append(parts, fmt.Sprintf("Dependency Graph: %d nodes, %d edges", len(g.Nodes), len(g.Edges)))

	// Count by dependency type
	typeCounts := make(map[DependencyType]int)
	for _, edge := range g.Edges {
		typeCounts[edge.DependencyType]++
	}
	for _, dt := range AllDependencyTypes {
		if count := typeCounts[dt]; count > 0 {
			parts = append(parts, fmt.Sprintf("  %s: %d", dt, count))
		}
	}

	// Critical nodes
	var criticalNodes []string
	for _, node := range g.Nodes {
		if node.Criticality >= 0.8 || node.Centrality >= 0.5 {
			criticalNodes = append(criticalNodes, node.ID)
		}
	}
	if len(criticalNodes) > 0 {
		parts = append(parts, fmt.Sprintf("Critical nodes: %d", len(criticalNodes)))
	}

	return strings.Join(parts, "\n")
}

// TrustChain represents a complete chain of dependencies from a root to a leaf.
type TrustChain struct {
	ID              string   `json:"id"`
	Nodes           []string `json:"nodes"`
	Length          int      `json:"length"`
	Confidence      float64  `json:"confidence"`
	Risk            string   `json:"risk"`
	DependencyCount int      `json:"dependency_count"`
	RootNode        string   `json:"root_node"`
	LeafNode        string   `json:"leaf_node"`
}

// CascadeResult represents a single step in a failure cascade.
type CascadeResult struct {
	Step             int      `json:"step"`
	AssumptionID     string   `json:"assumption_id"`
	AssumptionText   string   `json:"assumption_text"`
	Severity         string   `json:"severity"`
	AffectedAssets   []string `json:"affected_assets,omitempty"`
	AffectedControls []string `json:"affected_controls,omitempty"`
	Reason           string   `json:"reason"`
}

// FailureCascade represents the full cascade from a failed assumption.
type FailureCascade struct {
	RootAssumptionID   string          `json:"root_assumption_id"`
	RootAssumptionText string          `json:"root_assumption_text"`
	Steps              []CascadeResult `json:"steps"`
	TotalAffected      int             `json:"total_affected"`
	Severity           string          `json:"severity"`
	MaxDepth           int             `json:"max_depth"`
}

// CriticalAssumptionResult represents a critical assumption analysis.
type CriticalAssumptionResult struct {
	AssumptionID    string           `json:"assumption_id"`
	AssumptionText  string           `json:"assumption_text"`
	Centrality      float64          `json:"centrality"`
	SupportCount    int              `json:"support_count"`
	FailureRadius   int              `json:"failure_radius"`
	TrustRadius     int              `json:"trust_radius"`
	Risk            string           `json:"risk"`
	Score           float64          `json:"score"`
	DependencyTypes []DependencyType `json:"dependency_types"`
}

// SinglePointOfTrustFailure represents a detected single point of trust failure.
type SinglePointOfTrustFailure struct {
	NodeID          string           `json:"node_id"`
	AssumptionText  string           `json:"assumption_text"`
	DependentsCount int              `json:"dependents_count"`
	DependentNodes  []string         `json:"dependent_nodes"`
	DependencyTypes []DependencyType `json:"dependency_types"`
	Recommendation  string           `json:"recommendation"`
}

// TrustCollapseResult represents the result of a trust collapse simulation.
type TrustCollapseResult struct {
	FailedAssumptionID   string   `json:"failed_assumption_id"`
	FailedAssumptionText string   `json:"failed_assumption_text"`
	AssumptionsLost      []string `json:"assumptions_lost"`
	ControlsLost         []string `json:"controls_lost,omitempty"`
	AssetsExposed        []string `json:"assets_exposed,omitempty"`
	RiskIncrease         string   `json:"risk_increase"`
	RiskScoreBefore      float64  `json:"risk_score_before"`
	RiskScoreAfter       float64  `json:"risk_score_after"`
	AffectedComponents   []string `json:"affected_components,omitempty"`
}

// ChainOutput is the top-level container for all chain analysis results.
type ChainOutput struct {
	DependencyGraph      *DependencyGraph            `json:"dependency_graph"`
	TrustChains          []TrustChain                `json:"trust_chains"`
	FailureCascades      []FailureCascade            `json:"failure_cascades"`
	CriticalAssumptions  []CriticalAssumptionResult  `json:"critical_assumptions"`
	SinglePointsOfTrust  []SinglePointOfTrustFailure `json:"single_points_of_trust"`
	TrustCollapseResults []TrustCollapseResult       `json:"trust_collapse_results"`
	Domain               string                      `json:"domain"`
	GeneratedAt          string                      `json:"generated_at"`
}
