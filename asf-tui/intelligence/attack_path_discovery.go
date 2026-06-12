package intelligence

import (
	"fmt"
	"sort"
	"strings"
)

// ─────────────────────────────────────────────────────────────
// PHASE 1 — ATTACK PATH DATA MODEL
// ─────────────────────────────────────────────────────────────

type AttackPath struct {
	ID                  string       `json:"id"`
	Name                string       `json:"name"`
	Description         string       `json:"description"`
	EntryPoint          string       `json:"entry_point"`
	TargetAsset         string       `json:"target_asset"`
	AttackSteps         []AttackStep `json:"attack_steps,omitempty"`
	RequiredAssumptions []string     `json:"required_assumptions,omitempty"`
	RequiredConditions  []string     `json:"required_conditions,omitempty"`
	ExploitedThreats    []string     `json:"exploited_threats,omitempty"`
	AffectedComponents  []string     `json:"affected_components,omitempty"`
	AffectedBoundaries  []string     `json:"affected_boundaries,omitempty"`
	Likelihood          float64      `json:"likelihood"`
	Impact              float64      `json:"impact"`
	RiskScore           float64      `json:"risk_score"`
	Confidence          float64      `json:"confidence"`
	DetectionDifficulty string       `json:"detection_difficulty"`
	BusinessImpact      string       `json:"business_impact"`
	Recommendations     []string     `json:"recommendations,omitempty"`
	KillChainPhases     []string     `json:"kill_chain_phases,omitempty"`
	MITREATTACK         []string     `json:"mitre_attack,omitempty"`
	STRIDECategories    []string     `json:"stride_categories,omitempty"`
}

// ─────────────────────────────────────────────────────────────
// PHASE 2 — ATTACK STEP MODEL
// ─────────────────────────────────────────────────────────────

type AttackStep struct {
	SequenceNumber     int    `json:"sequence_number"`
	SourceComponent    string `json:"source_component"`
	TargetComponent    string `json:"target_component"`
	Action             string `json:"action"`
	Threat             string `json:"threat"`
	RequiredAssumption string `json:"required_assumption"`
	ControlBypassed    string `json:"control_bypassed"`
	Reasoning          string `json:"reasoning"`
	STRIDECategory     string `json:"stride_category"`
}

// ─────────────────────────────────────────────────────────────
// PHASE 3 — ATTACK GRAPH
// ─────────────────────────────────────────────────────────────

type AttackGraphNode struct {
	ID   string  `json:"id"`
	Name string  `json:"name"`
	Type string  `json:"type"` // "component", "zone", "entry", "target"
	Risk float64 `json:"risk"`
}

type AttackGraphEdge struct {
	Source  string   `json:"source"`
	Target  string   `json:"target"`
	Type    string   `json:"type"` // "relationship", "boundary", "threat"
	Threats []string `json:"threats,omitempty"`
	Risk    float64  `json:"risk"`
}

type AttackGraph struct {
	Nodes []AttackGraphNode `json:"nodes,omitempty"`
	Edges []AttackGraphEdge `json:"edges,omitempty"`
}

// ─────────────────────────────────────────────────────────────
// PHASE 4 — ENTRY POINT
// ─────────────────────────────────────────────────────────────

type EntryPoint struct {
	Component string  `json:"component"`
	ZoneType  string  `json:"zone_type"`
	Exposure  float64 `json:"exposure"` // 0.0-1.0
	Reasoning string  `json:"reasoning"`
}

// ─────────────────────────────────────────────────────────────
// PHASE 5 — TARGET ASSET
// ─────────────────────────────────────────────────────────────

type TargetAsset struct {
	Component   string  `json:"component"`
	Sensitivity string  `json:"sensitivity"` // "low", "medium", "high", "critical"
	Value       float64 `json:"value"`       // 0.0-1.0
	Reasoning   string  `json:"reasoning"`
}

// ─────────────────────────────────────────────────────────────
// PHASE 6 — THREAT CHAIN
// ─────────────────────────────────────────────────────────────

type ThreatChain struct {
	ID        string   `json:"id"`
	Threats   []string `json:"threats,omitempty"`
	Path      []string `json:"path,omitempty"` // component names
	RiskScore float64  `json:"risk_score"`
	Reasoning string   `json:"reasoning"`
}

// ─────────────────────────────────────────────────────────────
// PHASE 7 — APD ENGINE
// ─────────────────────────────────────────────────────────────

type EntryPointRule struct {
	Name     string
	Keywords []string
	Exposure float64
	ZoneType string
}

type TargetAssetRule struct {
	Name        string
	Keywords    []string
	Sensitivity string
	Value       float64
}

type AttackPathRule struct {
	Name       string
	EntryType  string
	TargetType string
	Action     string
	Threats    []string
}

type APDEngine struct {
	entryPointRules  map[string][]EntryPointRule
	targetAssetRules map[string][]TargetAssetRule
	attackPathRules  map[string][]AttackPathRule
	mitreMapping     map[string][]string
	killChainMapping map[string]string
}

type APDRunResult struct {
	AttackPaths       []AttackPath        `json:"attack_paths,omitempty"`
	ThreatChains      []ThreatChain       `json:"threat_chains,omitempty"`
	AttackGraph       AttackGraph         `json:"attack_graph"`
	EntryPoints       []EntryPoint        `json:"entry_points,omitempty"`
	TargetAssets      []TargetAsset       `json:"target_assets,omitempty"`
	TopPaths          []AttackPath        `json:"top_paths,omitempty"`
	Summary           string              `json:"summary"`
	KillChainCoverage map[string]int      `json:"kill_chain_coverage,omitempty"`
	MITREMapping      map[string][]string `json:"mitre_mapping,omitempty"`
	BusinessImpacts   []string            `json:"business_impacts,omitempty"`
}

// ─────────────────────────────────────────────────────────────
// PHASE 8 — CORE ENGINE METHODS
// ─────────────────────────────────────────────────────────────

func NewAPDEngine() *APDEngine {
	return &APDEngine{
		entryPointRules: map[string][]EntryPointRule{
			"internet": {
				{Name: "Internet Facing", Keywords: []string{"internet", "browser", "public", "web", "api", "gateway", "cdn", "edge", "dmz", "external", "mobile", "client"}, Exposure: 1.0, ZoneType: "INTERNET"},
			},
			"third_party": {
				{Name: "Third Party", Keywords: []string{"third", "vendor", "partner", "saas", "external", "integration"}, Exposure: 0.8, ZoneType: "THIRD_PARTY"},
			},
			"api_gateway": {
				{Name: "API Gateway", Keywords: []string{"gateway", "api", "ingress", "edge"}, Exposure: 0.7, ZoneType: "APPLICATION"},
			},
			"vpn": {
				{Name: "VPN", Keywords: []string{"vpn", "tunnel", "remote"}, Exposure: 0.5, ZoneType: "VPN"},
			},
			"admin": {
				{Name: "Admin Portal", Keywords: []string{"admin", "portal", "dashboard", "console", "management"}, Exposure: 0.6, ZoneType: "ADMINISTRATIVE"},
			},
		},
		targetAssetRules: map[string][]TargetAssetRule{
			"critical": {
				{Name: "Critical PHI", Keywords: []string{"phi", "hipaa", "patient", "medical", "health"}, Sensitivity: "critical", Value: 1.0},
				{Name: "Critical Payment", Keywords: []string{"payment", "pci", "card", "financial", "transaction"}, Sensitivity: "critical", Value: 1.0},
				{Name: "Critical Database", Keywords: []string{"database", "db", "pii", "secrets", "kms", "vault", "production", "cluster", "registry", "artifact"}, Sensitivity: "critical", Value: 0.95},
			},
			"high": {
				{Name: "High API", Keywords: []string{"api", "service", "application", "cache", "queue", "broker"}, Sensitivity: "high", Value: 0.8},
			},
			"medium": {
				{Name: "Medium Logging", Keywords: []string{"logging", "monitoring", "backup", "config"}, Sensitivity: "medium", Value: 0.5},
			},
			"low": {
				{Name: "Low Dev", Keywords: []string{"dev", "staging", "test", "build", "sandbox"}, Sensitivity: "low", Value: 0.2},
			},
		},
		attackPathRules: map[string][]AttackPathRule{
			"web_to_database": {
				{Name: "Web to Database", EntryType: "internet", TargetType: "database", Action: "SQL Injection", Threats: []string{"SQL Injection", "Data Exfiltration"}},
			},
			"api_to_secrets": {
				{Name: "API to Secrets", EntryType: "api_gateway", TargetType: "secrets", Action: "Secrets Theft", Threats: []string{"Secrets Theft", "Key Compromise"}},
			},
			"client_to_identity": {
				{Name: "Client to Identity", EntryType: "client", TargetType: "identity", Action: "Credential Stuffing", Threats: []string{"Credential Stuffing", "MFA Bypass"}},
			},
			"vpn_to_admin": {
				{Name: "VPN to Admin", EntryType: "vpn", TargetType: "admin", Action: "Privilege Escalation", Threats: []string{"Privilege Escalation", "Admin Account Compromise"}},
			},
			"third_party_to_data": {
				{Name: "Third Party to Data", EntryType: "third_party", TargetType: "database", Action: "Data Exfiltration", Threats: []string{"Data Exfiltration", "Supply Chain Attack"}},
			},
		},
		mitreMapping: map[string][]string{
			"Credential Stuffing":            {"T1110 - Brute Force"},
			"Valid Accounts":                 {"T1078 - Valid Accounts"},
			"SQL Injection":                  {"T1190 - Exploit Public-Facing Application"},
			"Privilege Escalation":           {"T1078 - Valid Accounts", "T1098 - Account Manipulation"},
			"Data Exfiltration":              {"T1048 - Exfiltration Over Alternative Protocol"},
			"Identity Compromise":            {"T1078 - Valid Accounts"},
			"Network Sniffing":               {"T1040 - Network Sniffing"},
			"DoS":                            {"T1498 - Network Denial of Service"},
			"Injection":                      {"T1190 - Exploit Public-Facing Application"},
			"Authentication Bypass":          {"T1550 - Use Alternate Authentication Material"},
			"Access Control Bypass":          {"T1098 - Account Manipulation"},
			"Secrets Theft":                  {"T1552 - Unsecured Credentials"},
			"Cross-Zone Lateral Movement":    {"T1021 - Remote Services"},
			"Man-in-the-Middle":              {"T1557 - Man-in-the-Middle"},
			"TLS Downgrade":                  {"T1557 - Man-in-the-Middle"},
			"Session Hijacking":              {"T1563 - Remote Service Session Hijacking"},
			"Token Replay":                   {"T1563 - Remote Service Session Hijacking"},
			"DDoS":                           {"T1498 - Network Denial of Service"},
			"Bot Attack":                     {"T1498 - Network Denial of Service"},
			"WAF Bypass":                     {"T1190 - Exploit Public-Facing Application"},
			"Cache Poisoning":                {"T1495 - Defacement"},
			"Message Queue Poisoning":        {"T1557 - Man-in-the-Middle"},
			"XML External Entity":            {"T1190 - Exploit Public-Facing Application"},
			"Broken Authentication":          {"T1078 - Valid Accounts"},
			"Sensitive Data Exposure":        {"T1040 - Network Sniffing"},
			"Broken Access Control":          {"T1098 - Account Manipulation"},
			"Security Misconfiguration":      {"T1190 - Exploit Public-Facing Application"},
			"Cross-Site Scripting":           {"T1189 - Drive-by Compromise"},
			"Insecure Deserialization":       {"T1190 - Exploit Public-Facing Application"},
			"Insufficient Logging":           {"T1562 - Impair Defenses"},
			"API Abuse":                      {"T1498 - Network Denial of Service"},
			"Request Smuggling":              {"T1190 - Exploit Public-Facing Application"},
			"Client Compromise":              {"T1200 - Exploitation for Client Execution"},
			"Log Tampering":                  {"T1565 - Data Manipulation"},
			"Missing Detection":              {"T1562 - Impair Defenses"},
			"Alert Fatigue":                  {"T1562 - Impair Defenses"},
			"Backup Exposure":                {"T1048 - Exfiltration Over Alternative Protocol"},
			"Restore Failure":                {"T1491 - Defacement"},
			"Backup Corruption":              {"T1491 - Defacement"},
			"Vendor Compromise":              {"T1195 - Supply Chain Compromise"},
			"SaaS Breach":                    {"T1195 - Supply Chain Compromise"},
			"Identity Provider Failure":      {"T1498 - Network Denial of Service"},
			"Privilege Abuse":                {"T1078 - Valid Accounts"},
			"Admin Account Compromise":       {"T1078 - Valid Accounts"},
			"VPN Tunnel Compromise":          {"T1557 - Man-in-the-Middle"},
			"Jump Host Compromise":           {"T1021 - Remote Services"},
			"Firewall Misconfiguration":      {"T1190 - Exploit Public-Facing Application"},
			"Key Exposure":                   {"T1552 - Unsecured Credentials"},
			"Rotation Failure":               {"T1552 - Unsecured Credentials"},
			"Key Misuse":                     {"T1552 - Unsecured Credentials"},
			"Unauthorized Data Access":       {"T1078 - Valid Accounts"},
			"Data Tampering":                 {"T1565 - Data Manipulation"},
			"Token Validation Failure":       {"T1550 - Use Alternate Authentication Material"},
			"Recovery Flow Abuse":            {"T1078 - Valid Accounts"},
			"Identity Provider Compromise":   {"T1078 - Valid Accounts"},
			"Certificate Validation Failure": {"T1557 - Man-in-the-Middle"},
			"Credential Exposure":            {"T1040 - Network Sniffing"},
			"Certificate Compromise":         {"T1552 - Unsecured Credentials"},
			"VPN Tunnel Hijacking":           {"T1557 - Man-in-the-Middle"},
			"MFA Bypass":                     {"T1078 - Valid Accounts"},
			"Data Exposure":                  {"T1040 - Network Sniffing"},
			"Backup Data Exposure":           {"T1048 - Exfiltration Over Alternative Protocol"},
			"Undetected Activity":            {"T1562 - Impair Defenses"},
			"DDoS Attack":                    {"T1498 - Network Denial of Service"},
			"Insider Threat":                 {"T1078 - Valid Accounts"},
			"Data Leakage":                   {"T1048 - Exfiltration Over Alternative Protocol"},
			"Secret Extraction":              {"T1552 - Unsecured Credentials"},
			"Lateral Movement":               {"T1021 - Remote Services"},
			"Default":                        {"T1190 - Exploit Public-Facing Application"},
		},
		killChainMapping: map[string]string{
			"Credential Stuffing":            "Initial Access",
			"SQL Injection":                  "Initial Access",
			"Privilege Escalation":           "Privilege Escalation",
			"Data Exfiltration":              "Exfiltration",
			"Identity Compromise":            "Credential Access",
			"Network Sniffing":               "Collection",
			"DoS":                            "Impact",
			"Injection":                      "Initial Access",
			"Authentication Bypass":          "Credential Access",
			"Access Control Bypass":          "Privilege Escalation",
			"Secrets Theft":                  "Credential Access",
			"Cross-Zone Lateral Movement":    "Lateral Movement",
			"Man-in-the-Middle":              "Collection",
			"TLS Downgrade":                  "Collection",
			"Session Hijacking":              "Credential Access",
			"Token Replay":                   "Credential Access",
			"DDoS":                           "Impact",
			"Bot Attack":                     "Reconnaissance",
			"WAF Bypass":                     "Initial Access",
			"Cache Poisoning":                "Impact",
			"Message Queue Poisoning":        "Execution",
			"XML External Entity":            "Initial Access",
			"Broken Authentication":          "Credential Access",
			"Sensitive Data Exposure":        "Collection",
			"Broken Access Control":          "Privilege Escalation",
			"Security Misconfiguration":      "Initial Access",
			"Cross-Site Scripting":           "Initial Access",
			"Insecure Deserialization":       "Initial Access",
			"Insufficient Logging":           "Defense Evasion",
			"API Abuse":                      "Impact",
			"Request Smuggling":              "Initial Access",
			"Client Compromise":              "Initial Access",
			"Log Tampering":                  "Defense Evasion",
			"Missing Detection":              "Defense Evasion",
			"Alert Fatigue":                  "Defense Evasion",
			"Backup Exposure":                "Exfiltration",
			"Restore Failure":                "Impact",
			"Backup Corruption":              "Impact",
			"Vendor Compromise":              "Initial Access",
			"SaaS Breach":                    "Exfiltration",
			"Identity Provider Failure":      "Impact",
			"Privilege Abuse":                "Privilege Escalation",
			"Admin Account Compromise":       "Credential Access",
			"VPN Tunnel Compromise":          "Lateral Movement",
			"Jump Host Compromise":           "Lateral Movement",
			"Firewall Misconfiguration":      "Initial Access",
			"Key Exposure":                   "Credential Access",
			"Rotation Failure":               "Credential Access",
			"Key Misuse":                     "Privilege Escalation",
			"Unauthorized Data Access":       "Collection",
			"Data Tampering":                 "Impact",
			"Token Validation Failure":       "Credential Access",
			"Recovery Flow Abuse":            "Credential Access",
			"Identity Provider Compromise":   "Credential Access",
			"Certificate Validation Failure": "Collection",
			"Credential Exposure":            "Collection",
			"Certificate Compromise":         "Credential Access",
			"VPN Tunnel Hijacking":           "Lateral Movement",
			"MFA Bypass":                     "Credential Access",
			"Data Exposure":                  "Exfiltration",
			"Backup Data Exposure":           "Exfiltration",
			"Undetected Activity":            "Defense Evasion",
			"DDoS Attack":                    "Impact",
			"Insider Threat":                 "Privilege Escalation",
			"Data Leakage":                   "Exfiltration",
			"Secret Extraction":              "Credential Access",
			"Lateral Movement":               "Lateral Movement",
			"Default":                        "Initial Access",
		},
	}
}

func (e *APDEngine) Run(arch *ArchDescription, threats []Threat, boundaries []TBITrustBoundary, zones []TrustZone, assumptions []Assumption) *APDRunResult {
	if arch == nil {
		return &APDRunResult{Summary: "No architecture provided"}
	}

	// Phase 1-2: Discover entry points and target assets
	entryPoints := e.discoverEntryPoints(arch)
	targetAssets := e.discoverTargetAssets(arch)

	// Phase 3: Build attack graph
	attackGraph := e.buildAttackGraph(arch, threats, boundaries, zones, entryPoints, targetAssets)

	// Phase 4-6: Construct attack paths and threat chains
	attackPaths := e.constructAttackPaths(arch, entryPoints, targetAssets, boundaries, threats, assumptions)
	threatChains := e.buildThreatChains(attackPaths, threats)

	// Phase 7: Risk scoring
	for i := range attackPaths {
		e.scoreAttackPath(&attackPaths[i], entryPoints, targetAssets, boundaries, threats, assumptions)
	}

	// Phase 8: Prioritization
	topPaths := e.prioritizePaths(attackPaths)

	// Phase 9: Business impact
	businessImpacts := e.generateBusinessImpacts(targetAssets)

	// Phase 10: Detection difficulty
	for i := range attackPaths {
		attackPaths[i].DetectionDifficulty = e.assessDetectionDifficulty(attackPaths[i], assumptions)
	}

	// Phase 11: Recommendations
	for i := range attackPaths {
		attackPaths[i].Recommendations = e.generateRecommendations(attackPaths[i], entryPoints, boundaries)
	}

	// Phase 12: Kill chain mapping
	killChainCoverage := make(map[string]int)
	for i := range attackPaths {
		attackPaths[i].KillChainPhases = e.mapKillChainPhases(attackPaths[i].AttackSteps)
		for _, phase := range attackPaths[i].KillChainPhases {
			killChainCoverage[phase]++
		}
	}

	// Phase 13: MITRE ATT&CK mapping
	mitreMapping := make(map[string][]string)
	for i := range attackPaths {
		attackPaths[i].MITREATTACK = e.mapMITREATTACK(attackPaths[i].AttackSteps)
		for _, mitre := range attackPaths[i].MITREATTACK {
			mitreMapping[mitre] = appendUnique(mitreMapping[mitre], attackPaths[i].ID)
		}
	}
	// Sort MITRE mapping keys deterministically
	var mitreKeys []string
	for k := range mitreMapping {
		mitreKeys = append(mitreKeys, k)
	}
	sort.Strings(mitreKeys)

	// Phase 14: STRIDE mapping
	for i := range attackPaths {
		attackPaths[i].STRIDECategories = e.mapSTRIDE(attackPaths[i].AttackSteps)
	}

	// Phase 15: Summary generation
	summary := e.generateSummary(attackPaths, threatChains, killChainCoverage, mitreKeys, businessImpacts)

	return &APDRunResult{
		AttackPaths:       attackPaths,
		ThreatChains:      threatChains,
		AttackGraph:       attackGraph,
		EntryPoints:       entryPoints,
		TargetAssets:      targetAssets,
		TopPaths:          topPaths,
		Summary:           summary,
		KillChainCoverage: killChainCoverage,
		MITREMapping:      mitreMapping,
		BusinessImpacts:   businessImpacts,
	}
}

// ─────────────────────────────────────────────────────────────
// ENTRY POINT DISCOVERY
// ─────────────────────────────────────────────────────────────

func (e *APDEngine) discoverEntryPoints(arch *ArchDescription) []EntryPoint {
	var entryPoints []EntryPoint
	seen := make(map[string]bool)

	// Sort keys for deterministic iteration
	var keys []string
	for k := range e.entryPointRules {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, comp := range arch.Components {
		lowerLabel := strings.ToLower(comp.Label)
		for _, key := range keys {
			rules := e.entryPointRules[key]
			for _, rule := range rules {
				for _, kw := range rule.Keywords {
					if strings.Contains(lowerLabel, kw) && !seen[comp.Label] {
						entryPoints = append(entryPoints, EntryPoint{
							Component: comp.Label,
							ZoneType:  rule.ZoneType,
							Exposure:  rule.Exposure,
							Reasoning: fmt.Sprintf("Component '%s' matches keyword '%s' for entry point rule '%s'", comp.Label, kw, rule.Name),
						})
						seen[comp.Label] = true
						break
					}
				}
			}
		}
	}

	// Sort by component name for determinism
	sort.Slice(entryPoints, func(i, j int) bool {
		return entryPoints[i].Component < entryPoints[j].Component
	})

	return entryPoints
}

// ─────────────────────────────────────────────────────────────
// TARGET ASSET IDENTIFICATION
// ─────────────────────────────────────────────────────────────

func (e *APDEngine) discoverTargetAssets(arch *ArchDescription) []TargetAsset {
	var targetAssets []TargetAsset
	seen := make(map[string]bool)

	// Order of sensitivity: critical > high > medium > low
	sensitivityOrder := []string{"critical", "high", "medium", "low"}
	for _, sens := range sensitivityOrder {
		rules, ok := e.targetAssetRules[sens]
		if !ok {
			continue
		}
		for _, comp := range arch.Components {
			if seen[comp.Label] {
				continue
			}
			lowerLabel := strings.ToLower(comp.Label)
			for _, rule := range rules {
				for _, kw := range rule.Keywords {
					if strings.Contains(lowerLabel, kw) {
						targetAssets = append(targetAssets, TargetAsset{
							Component:   comp.Label,
							Sensitivity: rule.Sensitivity,
							Value:       rule.Value,
							Reasoning:   fmt.Sprintf("Component '%s' matches keyword '%s' for target asset rule '%s'", comp.Label, kw, rule.Name),
						})
						seen[comp.Label] = true
						break
					}
				}
			}
		}
	}

	// Sort by component name for determinism
	sort.Slice(targetAssets, func(i, j int) bool {
		return targetAssets[i].Component < targetAssets[j].Component
	})

	return targetAssets
}

// ─────────────────────────────────────────────────────────────
// ATTACK GRAPH GENERATION
// ─────────────────────────────────────────────────────────────

func (e *APDEngine) buildAttackGraph(arch *ArchDescription, threats []Threat, boundaries []TBITrustBoundary, zones []TrustZone, entryPoints []EntryPoint, targetAssets []TargetAsset) AttackGraph {
	var graph AttackGraph
	nodeMap := make(map[string]bool)

	// Add component nodes
	for _, comp := range arch.Components {
		graph.Nodes = append(graph.Nodes, AttackGraphNode{
			ID:   comp.Label,
			Name: comp.Label,
			Type: "component",
			Risk: 0.0,
		})
		nodeMap[comp.Label] = true
	}

	// Add zone nodes
	for _, zone := range zones {
		zoneNodeID := "zone-" + zone.ID
		if !nodeMap[zoneNodeID] {
			graph.Nodes = append(graph.Nodes, AttackGraphNode{
				ID:   zoneNodeID,
				Name: zone.Name,
				Type: "zone",
				Risk: 0.0,
			})
			nodeMap[zoneNodeID] = true
		}
	}

	// Add entry point nodes
	for _, ep := range entryPoints {
		entryNodeID := "entry-" + ep.Component
		if !nodeMap[entryNodeID] {
			graph.Nodes = append(graph.Nodes, AttackGraphNode{
				ID:   entryNodeID,
				Name: ep.Component,
				Type: "entry",
				Risk: ep.Exposure,
			})
			nodeMap[entryNodeID] = true
		}
	}

	// Add target asset nodes
	for _, ta := range targetAssets {
		targetNodeID := "target-" + ta.Component
		if !nodeMap[targetNodeID] {
			graph.Nodes = append(graph.Nodes, AttackGraphNode{
				ID:   targetNodeID,
				Name: ta.Component,
				Type: "target",
				Risk: ta.Value,
			})
			nodeMap[targetNodeID] = true
		}
	}

	// Sort nodes by ID for determinism
	sort.Slice(graph.Nodes, func(i, j int) bool {
		return graph.Nodes[i].ID < graph.Nodes[j].ID
	})

	// Add relationship edges
	for _, rel := range arch.Relationships {
		graph.Edges = append(graph.Edges, AttackGraphEdge{
			Source:  rel.Source,
			Target:  rel.Target,
			Type:    "relationship",
			Threats: []string{},
			Risk:    0.0,
		})
	}

	// Add boundary edges
	for _, boundary := range boundaries {
		riskVal := riskLevelToFloat(boundary.Risk)
		graph.Edges = append(graph.Edges, AttackGraphEdge{
			Source:  boundary.SourceZone,
			Target:  boundary.DestinationZone,
			Type:    "boundary",
			Threats: boundary.Threats,
			Risk:    riskVal,
		})
	}

	// Add threat edges
	for _, threat := range threats {
		affected := make([]string, len(threat.AffectedComponents))
		copy(affected, threat.AffectedComponents)
		sort.Strings(affected)
		for i := 0; i < len(affected); i++ {
			for j := i + 1; j < len(affected); j++ {
				graph.Edges = append(graph.Edges, AttackGraphEdge{
					Source:  affected[i],
					Target:  affected[j],
					Type:    "threat",
					Threats: []string{threat.Name},
					Risk:    threat.RiskScore,
				})
			}
		}
	}

	// Add entry-to-component edges
	for _, ep := range entryPoints {
		graph.Edges = append(graph.Edges, AttackGraphEdge{
			Source:  "entry-" + ep.Component,
			Target:  ep.Component,
			Type:    "relationship",
			Threats: []string{},
			Risk:    0.0,
		})
	}

	// Add target-to-component edges
	for _, ta := range targetAssets {
		graph.Edges = append(graph.Edges, AttackGraphEdge{
			Source:  "target-" + ta.Component,
			Target:  ta.Component,
			Type:    "relationship",
			Threats: []string{},
			Risk:    0.0,
		})
	}

	// Sort edges by source, then target for determinism
	sort.Slice(graph.Edges, func(i, j int) bool {
		if graph.Edges[i].Source == graph.Edges[j].Source {
			return graph.Edges[i].Target < graph.Edges[j].Target
		}
		return graph.Edges[i].Source < graph.Edges[j].Source
	})

	return graph
}

// ─────────────────────────────────────────────────────────────
// PATH CONSTRUCTION
// ─────────────────────────────────────────────────────────────

func (e *APDEngine) constructAttackPaths(arch *ArchDescription, entryPoints []EntryPoint, targetAssets []TargetAsset, boundaries []TBITrustBoundary, threats []Threat, assumptions []Assumption) []AttackPath {
	adj := e.buildAdjacencyList(arch)
	var attackPaths []AttackPath
	pathID := 1

	// Sort entry points and target assets for determinism
	for _, ep := range entryPoints {
		for _, ta := range targetAssets {
			if ep.Component == ta.Component {
				continue
			}
			paths := e.findPathsDFS(adj, ep.Component, ta.Component, 10)
			for _, path := range paths {
				if len(path) < 2 {
					continue
				}
				ap := e.buildAttackPath(path, ep, ta, boundaries, threats, assumptions, pathID)
				attackPaths = append(attackPaths, ap)
				pathID++
			}
		}
	}

	return attackPaths
}

func (e *APDEngine) buildAdjacencyList(arch *ArchDescription) map[string][]string {
	adj := make(map[string][]string)
	for _, comp := range arch.Components {
		adj[comp.Label] = []string{}
	}
	for _, rel := range arch.Relationships {
		adj[rel.Source] = append(adj[rel.Source], rel.Target)
		adj[rel.Target] = append(adj[rel.Target], rel.Source)
	}
	for k, v := range adj {
		adj[k] = uniqueSortedStrings(v)
	}
	return adj
}

func uniqueSortedStrings(slice []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range slice {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	sort.Strings(result)
	return result
}

func (e *APDEngine) findPathsDFS(adj map[string][]string, start, end string, maxDepth int) [][]string {
	var paths [][]string
	var dfs func(current string, path []string, visited map[string]bool)
	dfs = func(current string, path []string, visited map[string]bool) {
		if len(path) > maxDepth {
			return
		}
		if current == end {
			newPath := make([]string, len(path))
			copy(newPath, path)
			paths = append(paths, newPath)
			return
		}
		for _, next := range adj[current] {
			if visited[next] {
				continue
			}
			visited[next] = true
			dfs(next, append(path, next), visited)
			visited[next] = false
		}
	}
	visited := make(map[string]bool)
	visited[start] = true
	dfs(start, []string{start}, visited)
	return paths
}

func (e *APDEngine) buildAttackPath(path []string, ep EntryPoint, ta TargetAsset, boundaries []TBITrustBoundary, threats []Threat, assumptions []Assumption, pathID int) AttackPath {
	ap := AttackPath{
		ID:             fmt.Sprintf("APD-%03d", pathID),
		Name:           fmt.Sprintf("%s via %s", ta.Component, ep.Component),
		Description:    fmt.Sprintf("Attacker enters through %s and reaches %s", ep.Component, ta.Component),
		EntryPoint:     ep.Component,
		TargetAsset:    ta.Component,
		BusinessImpact: e.computeBusinessImpact(ta),
	}

	// Build attack steps
	for i := 0; i < len(path)-1; i++ {
		step := e.buildAttackStep(path, i, boundaries, threats, assumptions)
		ap.AttackSteps = append(ap.AttackSteps, step)
	}

	// Collect affected components
	ap.AffectedComponents = make([]string, len(path))
	copy(ap.AffectedComponents, path)

	// Collect affected boundaries
	ap.AffectedBoundaries = e.findCrossedBoundaries(path, boundaries)

	// Collect required assumptions
	ap.RequiredAssumptions = e.findAssumptionsForPath(path, assumptions)

	// Collect exploited threats
	ap.ExploitedThreats = e.findThreatsForPath(path, threats)

	// Required conditions
	ap.RequiredConditions = e.buildRequiredConditions(path, boundaries)

	return ap
}

func (e *APDEngine) buildAttackStep(path []string, index int, boundaries []TBITrustBoundary, threats []Threat, assumptions []Assumption) AttackStep {
	source := path[index]
	target := path[index+1]
	action := e.inferAction(source, target, boundaries)
	threat := e.findThreatForStep(source, target, threats)
	requiredAssumption := e.findAssumptionForComponent(source, assumptions)
	controlBypassed := e.findControlBypassed(source, target, boundaries)
	reasoning := fmt.Sprintf("Step %d: Attacker moves from %s to %s via %s. Threat: %s. Control bypassed: %s.",
		index+1, source, target, action, threat, controlBypassed)
	stride := e.inferSTRIDE(action)

	return AttackStep{
		SequenceNumber:     index + 1,
		SourceComponent:    source,
		TargetComponent:    target,
		Action:             action,
		Threat:             threat,
		RequiredAssumption: requiredAssumption,
		ControlBypassed:    controlBypassed,
		Reasoning:          reasoning,
		STRIDECategory:     stride,
	}
}

func (e *APDEngine) inferAction(source, target string, boundaries []TBITrustBoundary) string {
	// Check if this step crosses a boundary
	for _, boundary := range boundaries {
		if (strings.Contains(source, boundary.SourceZone) && strings.Contains(target, boundary.DestinationZone)) ||
			(strings.Contains(source, boundary.DestinationZone) && strings.Contains(target, boundary.SourceZone)) {
			return "Cross-Zone Lateral Movement"
		}
	}

	// Infer from component names
	sourceLower := strings.ToLower(source)
	targetLower := strings.ToLower(target)

	if strings.Contains(targetLower, "database") || strings.Contains(targetLower, "db") {
		return "SQL Injection"
	}
	if strings.Contains(targetLower, "secret") || strings.Contains(targetLower, "vault") || strings.Contains(targetLower, "kms") {
		return "Secrets Theft"
	}
	if strings.Contains(targetLower, "identity") || strings.Contains(targetLower, "auth") || strings.Contains(targetLower, "idp") {
		return "Identity Compromise"
	}
	if strings.Contains(sourceLower, "api") || strings.Contains(sourceLower, "gateway") {
		return "API Abuse"
	}
	if strings.Contains(sourceLower, "vpn") || strings.Contains(sourceLower, "tunnel") {
		return "VPN Tunnel Compromise"
	}
	if strings.Contains(sourceLower, "admin") || strings.Contains(sourceLower, "management") {
		return "Privilege Escalation"
	}
	if strings.Contains(sourceLower, "client") || strings.Contains(sourceLower, "browser") {
		return "Credential Stuffing"
	}
	if strings.Contains(sourceLower, "internet") || strings.Contains(sourceLower, "web") {
		return "Injection"
	}
	if strings.Contains(targetLower, "cache") {
		return "Cache Poisoning"
	}
	if strings.Contains(targetLower, "queue") || strings.Contains(targetLower, "broker") {
		return "Message Queue Poisoning"
	}
	if strings.Contains(sourceLower, "log") || strings.Contains(sourceLower, "monitor") {
		return "Network Sniffing"
	}
	if strings.Contains(sourceLower, "firewall") || strings.Contains(sourceLower, "waf") {
		return "WAF Bypass"
	}

	return "Lateral Movement"
}

func (e *APDEngine) findThreatForStep(source, target string, threats []Threat) string {
	for _, threat := range threats {
		for _, comp := range threat.AffectedComponents {
			if comp == source || comp == target {
				return threat.Name
			}
		}
	}
	return "Unknown"
}

func (e *APDEngine) findAssumptionForComponent(component string, assumptions []Assumption) string {
	for _, assumption := range assumptions {
		for _, comp := range assumption.SourceComponents {
			if comp == component {
				return assumption.Description
			}
		}
		if assumption.Component == component {
			return assumption.Description
		}
	}
	return "No documented assumption"
}

func (e *APDEngine) findControlBypassed(source, target string, boundaries []TBITrustBoundary) string {
	for _, boundary := range boundaries {
		if (strings.Contains(source, boundary.SourceZone) && strings.Contains(target, boundary.DestinationZone)) ||
			(strings.Contains(source, boundary.DestinationZone) && strings.Contains(target, boundary.SourceZone)) {
			if len(boundary.MissingControls) > 0 {
				return strings.Join(boundary.MissingControls, ", ")
			}
			return "Boundary control not enforced"
		}
	}
	return "None"
}

func (e *APDEngine) inferSTRIDE(action string) string {
	actionLower := strings.ToLower(action)
	strideMap := map[string]string{
		"credential stuffing":         "Spoofing",
		"sql injection":               "Tampering, Information Disclosure",
		"privilege escalation":        "Elevation of Privilege",
		"data exfiltration":           "Information Disclosure",
		"dos":                         "Denial of Service",
		"injection":                   "Tampering, Information Disclosure",
		"identity compromise":         "Spoofing",
		"network sniffing":            "Information Disclosure",
		"secrets theft":               "Information Disclosure",
		"access control bypass":       "Elevation of Privilege",
		"cross-zone lateral movement": "Information Disclosure",
		"api abuse":                   "Denial of Service",
		"waf bypass":                  "Tampering",
		"cache poisoning":             "Tampering",
		"message queue poisoning":     "Tampering",
		"vpn tunnel compromise":       "Information Disclosure",
		"session hijacking":           "Spoofing",
		"token replay":                "Spoofing",
		"ddos":                        "Denial of Service",
		"bot attack":                  "Denial of Service",
		"log tampering":               "Tampering, Repudiation",
		"missing detection":           "Repudiation",
		"alert fatigue":               "Repudiation",
		"backup exposure":             "Information Disclosure",
		"restore failure":             "Denial of Service",
		"backup corruption":           "Tampering, Denial of Service",
		"vendor compromise":           "Spoofing, Information Disclosure, Elevation of Privilege",
		"saas breach":                 "Information Disclosure",
		"insider threat":              "Elevation of Privilege, Tampering",
		"data leakage":                "Information Disclosure",
		"secret extraction":           "Information Disclosure",
		"lateral movement":            "Information Disclosure",
		"man-in-the-middle":           "Information Disclosure, Tampering",
		"tls downgrade":               "Information Disclosure, Tampering",
		"credential exposure":         "Information Disclosure",
		"certificate compromise":      "Information Disclosure, Spoofing",
		"default":                     "Information Disclosure",
	}

	// Sort keys for deterministic iteration
	var keys []string
	for k := range strideMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		val := strideMap[key]
		if strings.Contains(actionLower, key) {
			return val
		}
	}
	return "Information Disclosure"
}

func (e *APDEngine) findCrossedBoundaries(path []string, boundaries []TBITrustBoundary) []string {
	var crossed []string
	seen := make(map[string]bool)
	for i := 0; i < len(path)-1; i++ {
		for _, boundary := range boundaries {
			if (strings.Contains(path[i], boundary.SourceZone) && strings.Contains(path[i+1], boundary.DestinationZone)) ||
				(strings.Contains(path[i], boundary.DestinationZone) && strings.Contains(path[i+1], boundary.SourceZone)) {
				if !seen[boundary.ID] {
					crossed = append(crossed, boundary.ID)
					seen[boundary.ID] = true
				}
			}
		}
	}
	sort.Strings(crossed)
	return crossed
}

func (e *APDEngine) findAssumptionsForPath(path []string, assumptions []Assumption) []string {
	var result []string
	seen := make(map[string]bool)
	for _, comp := range path {
		for _, assumption := range assumptions {
			for _, srcComp := range assumption.SourceComponents {
				if srcComp == comp && !seen[assumption.Description] {
					result = append(result, assumption.Description)
					seen[assumption.Description] = true
				}
			}
			if assumption.Component == comp && !seen[assumption.Description] {
				result = append(result, assumption.Description)
				seen[assumption.Description] = true
			}
		}
	}
	sort.Strings(result)
	return result
}

func (e *APDEngine) findThreatsForPath(path []string, threats []Threat) []string {
	var result []string
	seen := make(map[string]bool)
	for _, comp := range path {
		for _, threat := range threats {
			for _, affectedComp := range threat.AffectedComponents {
				if affectedComp == comp && !seen[threat.Name] {
					result = append(result, threat.Name)
					seen[threat.Name] = true
				}
			}
		}
	}
	sort.Strings(result)
	return result
}

func (e *APDEngine) buildRequiredConditions(path []string, boundaries []TBITrustBoundary) []string {
	var conditions []string
	seen := make(map[string]bool)
	for i := 0; i < len(path)-1; i++ {
		for _, boundary := range boundaries {
			if (strings.Contains(path[i], boundary.SourceZone) && strings.Contains(path[i+1], boundary.DestinationZone)) ||
				(strings.Contains(path[i], boundary.DestinationZone) && strings.Contains(path[i+1], boundary.SourceZone)) {
				for _, assumption := range boundary.RequiredAssumptions {
					if !seen[assumption] {
						conditions = append(conditions, assumption)
						seen[assumption] = true
					}
				}
			}
		}
	}
	sort.Strings(conditions)
	return conditions
}

// ─────────────────────────────────────────────────────────────
// THREAT CHAINING
// ─────────────────────────────────────────────────────────────

func (e *APDEngine) buildThreatChains(attackPaths []AttackPath, threats []Threat) []ThreatChain {
	var chains []ThreatChain
	chainID := 1

	for _, ap := range attackPaths {
		chainThreats := e.findThreatsForPath(ap.AffectedComponents, threats)
		if len(chainThreats) == 0 {
			continue
		}
		riskScore := 0.0
		for _, tName := range chainThreats {
			for _, threat := range threats {
				if threat.Name == tName {
					riskScore += threat.RiskScore
					break
				}
			}
		}
		chains = append(chains, ThreatChain{
			ID:        fmt.Sprintf("TC-%03d", chainID),
			Threats:   chainThreats,
			Path:      ap.AffectedComponents,
			RiskScore: clamp(riskScore/float64(len(chainThreats)), 0.0, 1.0),
			Reasoning: fmt.Sprintf("Threat chain along path %s → %s: %d threats mapped", ap.EntryPoint, ap.TargetAsset, len(chainThreats)),
		})
		chainID++
	}

	// Sort by risk score descending
	sort.Slice(chains, func(i, j int) bool {
		return chains[i].RiskScore > chains[j].RiskScore
	})

	return chains
}

// ─────────────────────────────────────────────────────────────
// RISK SCORING
// ─────────────────────────────────────────────────────────────

func (e *APDEngine) scoreAttackPath(ap *AttackPath, entryPoints []EntryPoint, targetAssets []TargetAsset, boundaries []TBITrustBoundary, threats []Threat, assumptions []Assumption) {
	// Find entry point exposure
	var exposure float64
	for _, ep := range entryPoints {
		if ep.Component == ap.EntryPoint {
			exposure = ep.Exposure
			break
		}
	}

	// Find target asset value
	var impact float64
	for _, ta := range targetAssets {
		if ta.Component == ap.TargetAsset {
			impact = ta.Value
			break
		}
	}

	// Compute control coverage
	controlCoverage := e.computeControlCoverage(ap, boundaries, assumptions)

	// Likelihood
	likelihood := exposure * (1.0 - controlCoverage)

	// Risk score
	riskScore := likelihood * impact

	// Adjust by boundary count
	boundaryCount := len(ap.AffectedBoundaries)
	riskScore += float64(boundaryCount) * 0.1

	// Adjust by threat count
	threatCount := len(ap.ExploitedThreats)
	riskScore += float64(threatCount) * 0.05

	// Cap
	riskScore = clamp(riskScore, 0.0, 1.0)
	likelihood = clamp(likelihood, 0.0, 1.0)
	impact = clamp(impact, 0.0, 1.0)

	// Confidence
	confidence := 0.5 + float64(len(ap.RequiredAssumptions))*0.05 + float64(len(ap.AttackSteps))*0.02
	confidence = clamp(confidence, 0.0, 1.0)

	ap.Likelihood = likelihood
	ap.Impact = impact
	ap.RiskScore = riskScore
	ap.Confidence = confidence
}

func (e *APDEngine) computeControlCoverage(ap *AttackPath, boundaries []TBITrustBoundary, assumptions []Assumption) float64 {
	coverage := 0.0

	// Count boundary controls
	for _, boundaryID := range ap.AffectedBoundaries {
		for _, boundary := range boundaries {
			if boundary.ID == boundaryID {
				coverage += float64(len(boundary.RequiredControls)) * 0.03
			}
		}
	}

	// Check assumptions for controls
	for _, assumption := range assumptions {
		lowerDesc := strings.ToLower(assumption.Description)
		if strings.Contains(lowerDesc, "control") || strings.Contains(lowerDesc, "firewall") || strings.Contains(lowerDesc, "waf") || strings.Contains(lowerDesc, "mfa") || strings.Contains(lowerDesc, "encryption") {
			coverage += 0.05
		}
	}

	return clamp(coverage, 0.0, 0.9)
}

// ─────────────────────────────────────────────────────────────
// PRIORITIZATION
// ─────────────────────────────────────────────────────────────

func (e *APDEngine) prioritizePaths(attackPaths []AttackPath) []AttackPath {
	// Sort by risk score descending, then by name for stability
	sorted := make([]AttackPath, len(attackPaths))
	copy(sorted, attackPaths)
	sort.SliceStable(sorted, func(i, j int) bool {
		if sorted[i].RiskScore == sorted[j].RiskScore {
			return sorted[i].Name < sorted[j].Name
		}
		return sorted[i].RiskScore > sorted[j].RiskScore
	})

	// Take top 10
	if len(sorted) > 10 {
		sorted = sorted[:10]
	}
	return sorted
}

// ─────────────────────────────────────────────────────────────
// BUSINESS IMPACT
// ─────────────────────────────────────────────────────────────

func (e *APDEngine) computeBusinessImpact(ta TargetAsset) string {
	switch ta.Sensitivity {
	case "critical":
		lowerLabel := strings.ToLower(ta.Component)
		if strings.Contains(lowerLabel, "phi") || strings.Contains(lowerLabel, "hipaa") || strings.Contains(lowerLabel, "patient") || strings.Contains(lowerLabel, "medical") {
			return "HIPAA breach, regulatory fines, patient notification"
		}
		if strings.Contains(lowerLabel, "payment") || strings.Contains(lowerLabel, "pci") || strings.Contains(lowerLabel, "card") {
			return "PCI DSS violation, financial fraud, chargebacks"
		}
		if strings.Contains(lowerLabel, "identity") || strings.Contains(lowerLabel, "auth") || strings.Contains(lowerLabel, "idp") || strings.Contains(lowerLabel, "sso") {
			return "Enterprise-wide account compromise, SSO breach"
		}
		return "Data breach, customer notification, regulatory inquiry"
	case "high":
		lowerLabel := strings.ToLower(ta.Component)
		if strings.Contains(lowerLabel, "database") || strings.Contains(lowerLabel, "db") {
			return "Data breach, customer notification, regulatory inquiry"
		}
		return "Service disruption, revenue loss, SLA breach"
	case "medium":
		return "Operational impact, recovery effort"
	case "low":
		return "Development impact, minimal business risk"
	default:
		return "Unknown business impact"
	}
}

func (e *APDEngine) generateBusinessImpacts(targetAssets []TargetAsset) []string {
	var impacts []string
	seen := make(map[string]bool)

	for _, ta := range targetAssets {
		impact := e.computeBusinessImpact(ta)
		if !seen[impact] {
			impacts = append(impacts, impact)
			seen[impact] = true
		}
	}

	sort.Strings(impacts)
	return impacts
}

// ─────────────────────────────────────────────────────────────
// DETECTION DIFFICULTY
// ─────────────────────────────────────────────────────────────

func (e *APDEngine) assessDetectionDifficulty(ap AttackPath, assumptions []Assumption) string {
	// Check for negative indicators first
	for _, assumption := range assumptions {
		lowerDesc := strings.ToLower(assumption.Description)
		if strings.Contains(lowerDesc, "no logging") || strings.Contains(lowerDesc, "blind spot") || strings.Contains(lowerDesc, "no monitoring") {
			return "Very Hard"
		}
	}

	// Check for positive indicators
	hasSIEM := false
	hasMonitoring := false
	hasLogging := false
	for _, assumption := range assumptions {
		lowerDesc := strings.ToLower(assumption.Description)
		if strings.Contains(lowerDesc, "siem") {
			hasSIEM = true
		}
		if strings.Contains(lowerDesc, "monitoring") {
			hasMonitoring = true
		}
		if strings.Contains(lowerDesc, "logging") {
			hasLogging = true
		}
	}

	if hasSIEM || hasMonitoring {
		return "Easy"
	}
	if hasLogging {
		return "Moderate"
	}
	return "Hard"
}

// ─────────────────────────────────────────────────────────────
// RECOMMENDATIONS
// ─────────────────────────────────────────────────────────────

func (e *APDEngine) generateRecommendations(ap AttackPath, entryPoints []EntryPoint, boundaries []TBITrustBoundary) []string {
	var recommendations []string

	// Prevent recommendations based on entry point
	for _, ep := range entryPoints {
		if ep.Component == ap.EntryPoint {
			lowerZone := strings.ToLower(ep.ZoneType)
			if strings.Contains(lowerZone, "identity") || strings.Contains(lowerZone, "admin") {
				recommendations = append(recommendations, "Prevent: Implement MFA and PAM for administrative access")
			} else if strings.Contains(lowerZone, "internet") || strings.Contains(lowerZone, "web") {
				recommendations = append(recommendations, "Prevent: Deploy WAF and DDoS protection at internet-facing edge")
			} else if strings.Contains(lowerZone, "network") || strings.Contains(lowerZone, "vpn") {
				recommendations = append(recommendations, "Prevent: Enforce network segmentation and VPN hardening")
			} else if strings.Contains(lowerZone, "application") || strings.Contains(lowerZone, "api") {
				recommendations = append(recommendations, "Prevent: Implement API rate limiting and input validation")
			} else {
				recommendations = append(recommendations, "Prevent: Implement defense in depth for entry point")
			}
			break
		}
	}

	// Detect recommendations
	recommendations = append(recommendations, "Detect: Deploy SIEM with alerting for anomalous access patterns")
	recommendations = append(recommendations, "Detect: Enable comprehensive logging and audit trail coverage")

	// Respond recommendations
	recommendations = append(recommendations, "Respond: Establish credential rotation and revocation procedures")
	recommendations = append(recommendations, "Respond: Implement network isolation and incident response playbook")

	// Boundary-specific recommendations
	for _, boundaryID := range ap.AffectedBoundaries {
		for _, boundary := range boundaries {
			if boundary.ID == boundaryID {
				for _, control := range boundary.RequiredControls {
					rec := fmt.Sprintf("Prevent: Implement required control '%s' at boundary %s", control, boundary.ID)
					recommendations = append(recommendations, rec)
				}
			}
		}
	}

	return uniqueSortedStrings(recommendations)
}

// ─────────────────────────────────────────────────────────────
// KILL CHAIN MAPPING
// ─────────────────────────────────────────────────────────────

func (e *APDEngine) mapKillChainPhases(steps []AttackStep) []string {
	var phases []string
	seen := make(map[string]bool)
	for _, step := range steps {
		phase := e.killChainMapping[step.Action]
		if phase == "" {
			phase = "Initial Access"
		}
		if !seen[phase] {
			phases = append(phases, phase)
			seen[phase] = true
		}
	}
	sort.Strings(phases)
	return phases
}

// ─────────────────────────────────────────────────────────────
// MITRE ATT&CK MAPPING
// ─────────────────────────────────────────────────────────────

func (e *APDEngine) mapMITREATTACK(steps []AttackStep) []string {
	var mitre []string
	seen := make(map[string]bool)
	for _, step := range steps {
		mappings, ok := e.mitreMapping[step.Action]
		if !ok || len(mappings) == 0 {
			mappings = e.mitreMapping["Default"]
		}
		for _, m := range mappings {
			if !seen[m] {
				mitre = append(mitre, m)
				seen[m] = true
			}
		}
	}
	sort.Strings(mitre)
	return mitre
}

// ─────────────────────────────────────────────────────────────
// STRIDE MAPPING
// ─────────────────────────────────────────────────────────────

func (e *APDEngine) mapSTRIDE(steps []AttackStep) []string {
	var stride []string
	seen := make(map[string]bool)
	for _, step := range steps {
		categories := strings.Split(step.STRIDECategory, ", ")
		for _, cat := range categories {
			cat = strings.TrimSpace(cat)
			if cat != "" && !seen[cat] {
				stride = append(stride, cat)
				seen[cat] = true
			}
		}
	}
	sort.Strings(stride)
	return stride
}

// ─────────────────────────────────────────────────────────────
// SUMMARY GENERATION
// ─────────────────────────────────────────────────────────────

func (e *APDEngine) generateSummary(attackPaths []AttackPath, threatChains []ThreatChain, killChainCoverage map[string]int, mitreKeys []string, businessImpacts []string) string {
	// Top 3 by risk
	var top3 []string
	for i, ap := range attackPaths {
		if i >= 3 {
			break
		}
		top3 = append(top3, fmt.Sprintf("%s (Risk: %.2f)", ap.Name, ap.RiskScore))
	}

	// Kill chain summary
	var killChainSummary []string
	for phase, count := range killChainCoverage {
		killChainSummary = append(killChainSummary, fmt.Sprintf("%s: %d", phase, count))
	}
	sort.Strings(killChainSummary)

	// MITRE summary
	mitreSummary := fmt.Sprintf("MITRE techniques mapped: %d", len(mitreKeys))

	return fmt.Sprintf(
		"Attack Path Discovery: %d attack paths identified, %d threat chains. Top 3 by risk: %s. Kill chain coverage: %s. %s. Business impacts: %d categories.",
		len(attackPaths),
		len(threatChains),
		strings.Join(top3, "; "),
		strings.Join(killChainSummary, ", "),
		mitreSummary,
		len(businessImpacts),
	)
}

// ─────────────────────────────────────────────────────────────
// UTILITY
// ─────────────────────────────────────────────────────────────

func riskLevelToFloat(r RiskLevel) float64 {
	switch r {
	case RiskCritical:
		return 1.0
	case RiskHigh:
		return 0.8
	case RiskMedium:
		return 0.5
	case RiskLow:
		return 0.2
	default:
		return 0.0
	}
}
