package intelligence

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// ═══════════════════════════════════════════════════════
// SAMPI — Security Architecture Memory & Portfolio Intelligence Engine (ASF V10)
// Phases 1-16
// ═══════════════════════════════════════════════════════

// ── PHASE 1 — ARCHITECTURE MEMORY LAYER ──

type ArchitectureRecord struct {
	ArchitectureID string            `json:"architecture_id"`
	Name           string            `json:"name"`
	Domain         string            `json:"domain"`
	AnalysisDate   time.Time         `json:"analysis_date"`
	Version        string            `json:"version"`
	Assumptions    []Assumption      `json:"assumptions"`
	Threats        []Threat          `json:"threats"`
	AttackPaths    []AttackPath      `json:"attack_paths"`
	Findings       []SDRIFinding     `json:"findings"`
	Controls       []SDRIControl     `json:"controls"`
	Compliance     []string          `json:"compliance"`
	RiskScore      float64           `json:"risk_score"`
	Metadata       map[string]string `json:"metadata,omitempty"`
}

func SaveArchitectureRecord(dir string, rec *ArchitectureRecord) error {
	if rec == nil {
		return fmt.Errorf("nil record")
	}
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	path := filepath.Join(dir, rec.ArchitectureID+".json")
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func LoadArchitectureRecord(path string) (*ArchitectureRecord, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var rec ArchitectureRecord
	if err := json.Unmarshal(data, &rec); err != nil {
		return nil, err
	}
	return &rec, nil
}

// ── PHASE 2 — PORTFOLIO DATABASE ──

type Portfolio struct {
	Architectures []ArchitectureRecord `json:"architectures"`
	CreatedAt     time.Time            `json:"created_at"`
	UpdatedAt     time.Time            `json:"updated_at"`
}

func NewPortfolio() *Portfolio {
	return &Portfolio{
		Architectures: make([]ArchitectureRecord, 0),
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

func (p *Portfolio) AddArchitecture(rec ArchitectureRecord) {
	for i, a := range p.Architectures {
		if a.ArchitectureID == rec.ArchitectureID {
			p.Architectures[i] = rec
			p.UpdatedAt = time.Now()
			return
		}
	}
	p.Architectures = append(p.Architectures, rec)
	p.UpdatedAt = time.Now()
}

func (p *Portfolio) RemoveArchitecture(id string) {
	filtered := make([]ArchitectureRecord, 0, len(p.Architectures))
	for _, a := range p.Architectures {
		if a.ArchitectureID != id {
			filtered = append(filtered, a)
		}
	}
	p.Architectures = filtered
	p.UpdatedAt = time.Now()
}

func (p *Portfolio) GetArchitecture(id string) *ArchitectureRecord {
	for i := range p.Architectures {
		if p.Architectures[i].ArchitectureID == id {
			return &p.Architectures[i]
		}
	}
	return nil
}

func (p *Portfolio) Save(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	p.UpdatedAt = time.Now()
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func (p *Portfolio) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, p)
}

// LoadPortfolioDir loads all JSON files from a directory as architecture records
// and returns a Portfolio. If the directory has a portfolio.json, that's loaded instead.
func LoadPortfolioDir(dir string) (*Portfolio, error) {
	portfolioPath := filepath.Join(dir, "portfolio.json")
	if _, err := os.Stat(portfolioPath); err == nil {
		p := NewPortfolio()
		if err := p.Load(portfolioPath); err != nil {
			return nil, err
		}
		return p, nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	p := NewPortfolio()
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}
		rec, err := LoadArchitectureRecord(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}
		p.Architectures = append(p.Architectures, *rec)
	}
	if len(p.Architectures) == 0 {
		return nil, fmt.Errorf("no architecture records found in %s", dir)
	}
	return p, nil
}

// ── SAMPI INPUT ──

type SAMPIInput struct {
	Portfolio *Portfolio
}

// ── SAMPI RESULT ──

type SAMPIResult struct {
	RepeatedWeaknesses   []RepeatedWeakness        `json:"repeated_weaknesses"`
	EnterpriseRiskThemes []EnterpriseRiskTheme     `json:"enterprise_risk_themes"`
	ControlCoverage      []ControlCoverageItem     `json:"control_coverage"`
	Comparisons          []ArchitectureComparison  `json:"comparisons"`
	RiskTrends           []RiskTrendRecord         `json:"risk_trends"`
	SecurityDebt         SecurityDebt              `json:"security_debt"`
	AttackSurface        PortfolioAttackSurface    `json:"attack_surface"`
	CrownJewelInventory  []CrownJewelInventoryItem `json:"crown_jewel_inventory"`
	SharedDependencies   []SharedDependency        `json:"shared_dependencies"`
	BlastRadii           []BlastRadius             `json:"blast_radii"`
	ComplianceView       EnterpriseComplianceView  `json:"compliance_view"`
	Dashboard            PortfolioDashboard        `json:"dashboard"`
	Heatmaps             []ExecutiveHeatmap        `json:"heatmaps"`
	ProgramInsights      []SecurityProgramInsight  `json:"program_insights"`
}

// ── SAMPI ENGINE ──

type SAMPIEngine struct{}

func NewSAMPIEngine() *SAMPIEngine {
	return &SAMPIEngine{}
}

func (e *SAMPIEngine) Run(input SAMPIInput) *SAMPIResult {
	result := &SAMPIResult{}
	portfolio := input.Portfolio
	if portfolio == nil || len(portfolio.Architectures) == 0 {
		return result
	}

	// Phase 3: Cross-architecture risk discovery
	result.RepeatedWeaknesses = findRepeatedWeaknesses(portfolio)

	// Phase 4: Enterprise risk themes
	result.EnterpriseRiskThemes = aggregateEnterpriseRiskThemes(portfolio)

	// Phase 5: Control reuse analysis
	result.ControlCoverage = analyzeControlReuse(portfolio)

	// Phase 6: Architecture comparison
	result.Comparisons = generateAllComparisons(portfolio)

	// Phase 7: Risk trending
	result.RiskTrends = computeRiskTrends(portfolio)

	// Phase 8: Security debt engine
	result.SecurityDebt = computeSecurityDebt(portfolio)

	// Phase 9: Portfolio attack surface
	result.AttackSurface = computePortfolioAttackSurface(portfolio)

	// Phase 10: Crown jewel inventory
	result.CrownJewelInventory = buildCrownJewelInventory(portfolio)

	// Phase 11: Shared dependency analysis
	result.SharedDependencies = findSharedDependencies(portfolio)

	// Phase 12: Blast radius modeling
	result.BlastRadii = computeBlastRadii(portfolio)

	// Phase 13: Enterprise compliance view
	result.ComplianceView = computeEnterpriseCompliance(portfolio)

	// Phase 14: Portfolio dashboard
	result.Dashboard = buildPortfolioDashboard(portfolio, result)

	// Phase 15: Executive heatmaps
	result.Heatmaps = generateHeatmaps(portfolio, result)

	// Phase 16: Security program insights
	result.ProgramInsights = generateProgramInsights(result)

	return result
}

// ── PHASE 3 — CROSS-ARCHITECTURE RISK DISCOVERY ──

type RepeatedWeakness struct {
	FindingTitle          string   `json:"finding_title"`
	Category              string   `json:"category"`
	Severity              string   `json:"severity"`
	AffectedArchitectures []string `json:"affected_architectures"`
	OccurrenceCount       int      `json:"occurrence_count"`
	Systemic              bool     `json:"systemic"`
}

func findRepeatedWeaknesses(p *Portfolio) []RepeatedWeakness {
	index := make(map[string]*RepeatedWeakness)
	for _, arch := range p.Architectures {
		seen := make(map[string]bool)
		for _, f := range arch.Findings {
			key := normalizeFindingKey(f.Title, f.Category)
			if seen[key] {
				continue
			}
			seen[key] = true
			if rw, ok := index[key]; ok {
				rw.OccurrenceCount++
				rw.AffectedArchitectures = append(rw.AffectedArchitectures, arch.Name)
			} else {
				index[key] = &RepeatedWeakness{
					FindingTitle:          f.Title,
					Category:              f.Category,
					Severity:              f.Severity,
					AffectedArchitectures: []string{arch.Name},
					OccurrenceCount:       1,
				}
			}
		}
	}
	result := make([]RepeatedWeakness, 0, len(index))
	for _, rw := range index {
		rw.Systemic = rw.OccurrenceCount >= 3
		result = append(result, *rw)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].OccurrenceCount > result[j].OccurrenceCount
	})
	return result
}

func normalizeFindingKey(title, category string) string {
	key := strings.ToLower(strings.TrimSpace(title))
	key = strings.NewReplacer(" ", "_", "-", "_", "/", "_").Replace(key)
	if category != "" {
		key = key + "::" + strings.ToLower(strings.TrimSpace(category))
	}
	return key
}

// ── PHASE 4 — ENTERPRISE RISK THEMES ──

type EnterpriseRiskTheme struct {
	Name                  string `json:"name"`
	Description           string `json:"description"`
	RiskCount             int    `json:"risk_count"`
	AffectedArchitectures int    `json:"affected_architectures"`
	Severity              string `json:"severity"`
}

func aggregateEnterpriseRiskThemes(p *Portfolio) []EnterpriseRiskTheme {
	themeFindings := map[string][]string{
		"Identity Risk":        {"access_control", "authentication", "authorization", "mfa", "identity", "credential", "sso", "iam"},
		"Third Party Risk":     {"third_party", "vendor", "supply_chain", "outsource", "partner", "external"},
		"Data Protection Risk": {"encryption", "data_protection", "data_loss", "dlp", "classification", "masking", "tokenization"},
		"Cloud Risk":           {"cloud", "container", "kubernetes", "s3", "aws", "azure", "gcp", "serverless"},
		"Compliance Risk":      {"compliance", "audit", "regulatory", "hipaa", "pci", "sox", "gdpr", "reporting"},
		"Operational Risk":     {"availability", "backup", "disaster", "bcdr", "incident", "monitoring", "logging"},
	}
	type themeMatch struct {
		count int
		archs map[string]bool
	}
	themes := make(map[string]*themeMatch)
	for name := range themeFindings {
		themes[name] = &themeMatch{archs: make(map[string]bool)}
	}
	for _, arch := range p.Architectures {
		for _, f := range arch.Findings {
			lower := strings.ToLower(f.Title + " " + f.Description + " " + f.Category)
			for themeName, keywords := range themeFindings {
				for _, kw := range keywords {
					if strings.Contains(lower, kw) {
						themes[themeName].count++
						themes[themeName].archs[arch.Name] = true
						break
					}
				}
			}
		}
	}
	result := make([]EnterpriseRiskTheme, 0, len(themes))
	for name, tm := range themes {
		if tm.count == 0 {
			continue
		}
		sev := "Low"
		if tm.count >= 10 {
			sev = "Critical"
		} else if tm.count >= 5 {
			sev = "High"
		} else if tm.count >= 2 {
			sev = "Medium"
		}
		result = append(result, EnterpriseRiskTheme{
			Name:                  name,
			Description:           fmt.Sprintf("Aggregated risk across %d architectures", len(tm.archs)),
			RiskCount:             tm.count,
			AffectedArchitectures: len(tm.archs),
			Severity:              sev,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].RiskCount > result[j].RiskCount
	})
	return result
}

// ── PHASE 5 — CONTROL REUSE ANALYSIS ──

type ControlCoverageItem struct {
	ControlName        string  `json:"control_name"`
	Category           string  `json:"category"`
	ArchitecturesWith  int     `json:"architectures_with"`
	ArchitecturesTotal int     `json:"architectures_total"`
	CoveragePercent    float64 `json:"coverage_percent"`
}

func analyzeControlReuse(p *Portfolio) []ControlCoverageItem {
	controlIndex := make(map[string]map[string]bool)
	total := len(p.Architectures)
	if total == 0 {
		return nil
	}
	for _, arch := range p.Architectures {
		for _, c := range arch.Controls {
			key := normalizeControlKey(c.Name)
			if controlIndex[key] == nil {
				controlIndex[key] = make(map[string]bool)
			}
			controlIndex[key][arch.Name] = true
		}
	}
	result := make([]ControlCoverageItem, 0, len(controlIndex))
	for key, archs := range controlIndex {
		coverage := float64(len(archs)) / float64(total) * 100
		result = append(result, ControlCoverageItem{
			ControlName:        key,
			ArchitecturesWith:  len(archs),
			ArchitecturesTotal: total,
			CoveragePercent:    coverage,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].CoveragePercent < result[j].CoveragePercent
	})
	return result
}

func normalizeControlKey(name string) string {
	key := strings.ToLower(strings.TrimSpace(name))
	return strings.NewReplacer(" ", "_", "-", "_", "/", "_").Replace(key)
}

// ── PHASE 6 — ARCHITECTURE COMPARISON ENGINE ──

type ArchitectureComparison struct {
	ArchitectureA     string  `json:"architecture_a"`
	ArchitectureB     string  `json:"architecture_b"`
	SharedAssumptions int     `json:"shared_assumptions"`
	SharedThreats     int     `json:"shared_threats"`
	SharedControls    int     `json:"shared_controls"`
	RiskScoreA        float64 `json:"risk_score_a"`
	RiskScoreB        float64 `json:"risk_score_b"`
	RiskDelta         float64 `json:"risk_delta"`
	SimilarityScore   float64 `json:"similarity_score"`
}

func generateAllComparisons(p *Portfolio) []ArchitectureComparison {
	archs := p.Architectures
	if len(archs) < 2 {
		return nil
	}
	comparisons := make([]ArchitectureComparison, 0)
	for i := 0; i < len(archs); i++ {
		for j := i + 1; j < len(archs); j++ {
			comparisons = append(comparisons, compareTwoArchitectures(&archs[i], &archs[j]))
		}
	}
	return comparisons
}

func compareTwoArchitectures(a, b *ArchitectureRecord) ArchitectureComparison {
	sharedS := countSharedStrings(assumptionIDs(a.Assumptions), assumptionIDs(b.Assumptions))
	sharedT := countSharedStrings(threatIDs(a.Threats), threatIDs(b.Threats))
	sharedC := countSharedStrings(controlIDs(a.Controls), controlIDs(b.Controls))
	maxFields := max(len(a.Assumptions), len(b.Assumptions)) +
		max(len(a.Threats), len(b.Threats)) +
		max(len(a.Controls), len(b.Controls))
	similarity := 0.0
	if maxFields > 0 {
		similarity = float64(sharedS+sharedT+sharedC) / float64(maxFields) * 100
	}
	return ArchitectureComparison{
		ArchitectureA:     a.Name,
		ArchitectureB:     b.Name,
		SharedAssumptions: sharedS,
		SharedThreats:     sharedT,
		SharedControls:    sharedC,
		RiskScoreA:        a.RiskScore,
		RiskScoreB:        b.RiskScore,
		RiskDelta:         b.RiskScore - a.RiskScore,
		SimilarityScore:   similarity,
	}
}

func countSharedStrings(a, b []string) int {
	set := make(map[string]bool)
	for _, s := range a {
		set[s] = true
	}
	count := 0
	for _, s := range b {
		if set[s] {
			count++
		}
	}
	return count
}

func assumptionIDs(assumptions []Assumption) []string {
	ids := make([]string, len(assumptions))
	for i, a := range assumptions {
		ids[i] = a.ID
	}
	return ids
}

func threatIDs(threats []Threat) []string {
	ids := make([]string, len(threats))
	for i, t := range threats {
		ids[i] = t.ID
	}
	return ids
}

func controlIDs(controls []SDRIControl) []string {
	ids := make([]string, len(controls))
	for i, c := range controls {
		ids[i] = c.ID
	}
	return ids
}

// ── PHASE 7 — RISK TRENDING ──

type RiskTrendRecord struct {
	ArchitectureID string  `json:"architecture_id"`
	Name           string  `json:"name"`
	PreviousScore  float64 `json:"previous_score"`
	CurrentScore   float64 `json:"current_score"`
	Direction      string  `json:"direction"`
}

func computeRiskTrends(p *Portfolio) []RiskTrendRecord {
	if len(p.Architectures) < 2 {
		return nil
	}
	byArch := make(map[string][]ArchitectureRecord)
	for _, arch := range p.Architectures {
		byArch[arch.ArchitectureID] = append(byArch[arch.ArchitectureID], arch)
	}
	trends := make([]RiskTrendRecord, 0)
	for id, versions := range byArch {
		if len(versions) < 2 {
			continue
		}
		sort.Slice(versions, func(i, j int) bool {
			return versions[i].AnalysisDate.Before(versions[j].AnalysisDate)
		})
		prev := versions[len(versions)-2].RiskScore
		curr := versions[len(versions)-1].RiskScore
		dir := "Stable"
		if curr < prev-0.5 {
			dir = "Improving"
		} else if curr > prev+0.5 {
			dir = "Worsening"
		}
		trends = append(trends, RiskTrendRecord{
			ArchitectureID: id,
			Name:           versions[len(versions)-1].Name,
			PreviousScore:  prev,
			CurrentScore:   curr,
			Direction:      dir,
		})
	}
	return trends
}

// ── PHASE 8 — SECURITY DEBT ENGINE ──

type SecurityDebt struct {
	Score             float64    `json:"score"`
	LongstandingCount int        `json:"longstanding_count"`
	RepeatedCount     int        `json:"repeated_count"`
	IgnoredCount      int        `json:"ignored_count"`
	TopDebts          []DebtItem `json:"top_debts"`
}

type DebtItem struct {
	Description  string `json:"description"`
	Architecture string `json:"architecture"`
	Category     string `json:"category"`
	Severity     string `json:"severity"`
	Age          string `json:"age"`
}

func computeSecurityDebt(p *Portfolio) SecurityDebt {
	findingsByKey := make(map[string][]struct {
		arch string
		sev  string
		cat  string
	})
	for _, arch := range p.Architectures {
		for _, f := range arch.Findings {
			key := normalizeFindingKey(f.Title, f.Category)
			findingsByKey[key] = append(findingsByKey[key], struct {
				arch string
				sev  string
				cat  string
			}{arch.Name, f.Severity, f.Category})
		}
	}
	longstanding := 0
	repeated := 0
	debts := make([]DebtItem, 0)
	for key, entries := range findingsByKey {
		if len(entries) >= 3 {
			repeated++
			debts = append(debts, DebtItem{
				Description:  key,
				Architecture: entries[0].arch,
				Category:     entries[0].cat,
				Severity:     entries[0].sev,
				Age:          "repeated_across_architectures",
			})
		}
		if len(entries) >= 1 {
			longstanding++
		}
	}
	sort.Slice(debts, func(i, j int) bool {
		return severityWeight(debts[i].Severity) > severityWeight(debts[j].Severity)
	})
	if len(debts) > 10 {
		debts = debts[:10]
	}
	score := 0.0
	if len(findingsByKey) > 0 {
		score = (float64(repeated)/float64(len(findingsByKey)))*100 +
			(float64(longstanding)/float64(len(findingsByKey)))*50
		if score > 100 {
			score = 100
		}
	}
	return SecurityDebt{
		Score:             score,
		LongstandingCount: longstanding,
		RepeatedCount:     repeated,
		TopDebts:          debts,
	}
}

func severityWeight(sev string) int {
	switch strings.ToLower(sev) {
	case "critical":
		return 5
	case "high":
		return 4
	case "medium":
		return 3
	case "low":
		return 2
	default:
		return 1
	}
}

// ── PHASE 9 — PORTFOLIO ATTACK SURFACE ──

type PortfolioAttackSurface struct {
	InternetExposure   int `json:"internet_exposure"`
	ThirdPartyExposure int `json:"third_party_exposure"`
	IdentityExposure   int `json:"identity_exposure"`
	CloudExposure      int `json:"cloud_exposure"`
	TotalExposure      int `json:"total_exposure"`
}

func computePortfolioAttackSurface(p *Portfolio) PortfolioAttackSurface {
	var s PortfolioAttackSurface
	for _, arch := range p.Architectures {
		for _, threat := range arch.Threats {
			lower := strings.ToLower(threat.Name + " " + threat.Description)
			switch {
			case containsAny(lower, []string{"internet", "external", "public", "exposed", "dmz"}):
				s.InternetExposure++
			case containsAny(lower, []string{"third_party", "vendor", "partner", "supply", "outsource"}):
				s.ThirdPartyExposure++
			case containsAny(lower, []string{"identity", "credential", "auth", "sso", "iam", "login"}):
				s.IdentityExposure++
			case containsAny(lower, []string{"cloud", "aws", "azure", "gcp", "container", "kubernetes"}):
				s.CloudExposure++
			}
		}
	}
	s.TotalExposure = s.InternetExposure + s.ThirdPartyExposure + s.IdentityExposure + s.CloudExposure
	return s
}

// ── PHASE 10 — CROWN JEWEL INVENTORY ──

type CrownJewelInventoryItem struct {
	Name             string `json:"name"`
	ArchitectureName string `json:"architecture_name"`
	Category         string `json:"category"`
	ThreatCount      int    `json:"threat_count"`
	RiskLevel        string `json:"risk_level"`
}

func buildCrownJewelInventory(p *Portfolio) []CrownJewelInventoryItem {
	assetIndex := make(map[string]*CrownJewelInventoryItem)
	for _, arch := range p.Architectures {
		for _, t := range arch.Threats {
			for _, asset := range t.AffectedAssets {
				key := strings.ToLower(asset) + "::" + arch.Name
				if item, ok := assetIndex[key]; ok {
					item.ThreatCount++
				} else {
					assetIndex[key] = &CrownJewelInventoryItem{
						Name:             asset,
						ArchitectureName: arch.Name,
						Category:         classifyAssetCategory(asset),
						ThreatCount:      1,
						RiskLevel:        "Medium",
					}
				}
			}
		}
	}
	result := make([]CrownJewelInventoryItem, 0, len(assetIndex))
	for _, item := range assetIndex {
		item.RiskLevel = riskLevelForThreatCount(item.ThreatCount)
		result = append(result, *item)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].ThreatCount > result[j].ThreatCount
	})
	if len(result) > 50 {
		result = result[:50]
	}
	return result
}

func classifyAssetCategory(name string) string {
	lower := strings.ToLower(name)
	switch {
	case containsAny(lower, []string{"database", "data", "store", "warehouse", "bucket"}):
		return "Data Store"
	case containsAny(lower, []string{"api", "gateway", "service", "endpoint"}):
		return "Service"
	case containsAny(lower, []string{"auth", "identity", "iam", "login", "sso"}):
		return "Identity"
	case containsAny(lower, []string{"network", "firewall", "vpn", "load_balancer", "proxy"}):
		return "Network"
	case containsAny(lower, []string{"pod", "container", "cluster", "node", "vm", "server"}):
		return "Compute"
	case containsAny(lower, []string{"app", "ui", "web", "portal", "dashboard"}):
		return "Application"
	default:
		return "Infrastructure"
	}
}

func riskLevelForThreatCount(count int) string {
	switch {
	case count >= 5:
		return "Critical"
	case count >= 3:
		return "High"
	case count >= 2:
		return "Medium"
	default:
		return "Low"
	}
}

// ── PHASE 11 — SHARED DEPENDENCY ANALYSIS ──

type SharedDependency struct {
	DependencyName      string   `json:"dependency_name"`
	Category            string   `json:"category"`
	UsedByArchitectures []string `json:"used_by_architectures"`
	UsageCount          int      `json:"usage_count"`
	RiskLevel           string   `json:"risk_level"`
}

func findSharedDependencies(p *Portfolio) []SharedDependency {
	depIndex := make(map[string]map[string]bool)
	for _, arch := range p.Architectures {
		for _, c := range arch.Controls {
			depName := normalizeDependencyName(c.Name)
			if depIndex[depName] == nil {
				depIndex[depName] = make(map[string]bool)
			}
			depIndex[depName][arch.Name] = true
		}
	}
	result := make([]SharedDependency, 0, len(depIndex))
	for dep, archs := range depIndex {
		rl := "Low"
		uc := len(archs)
		switch {
		case uc >= 10:
			rl = "Critical"
		case uc >= 5:
			rl = "High"
		case uc >= 3:
			rl = "Medium"
		}
		archNames := make([]string, 0, len(archs))
		for n := range archs {
			archNames = append(archNames, n)
		}
		sort.Strings(archNames)
		cat := classifyDependencyCategory(dep)
		result = append(result, SharedDependency{
			DependencyName:      dep,
			Category:            cat,
			UsedByArchitectures: archNames,
			UsageCount:          uc,
			RiskLevel:           rl,
		})
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].UsageCount > result[j].UsageCount
	})
	return result
}

func normalizeDependencyName(name string) string {
	return normalizeControlKey(name)
}

func classifyDependencyCategory(name string) string {
	lower := strings.ToLower(name)
	switch {
	case containsAny(lower, []string{"mfa", "auth", "identity", "sso", "iam", "ldap", "oauth"}):
		return "Identity"
	case containsAny(lower, []string{"kms", "key", "certificate", "crypto", "hsm", "vault"}):
		return "Cryptography"
	case containsAny(lower, []string{"firewall", "vpn", "network", "proxy", "waf", "gateway"}):
		return "Network Security"
	case containsAny(lower, []string{"log", "monitor", "siem", "splunk", "elk", "datadog"}):
		return "Monitoring"
	case containsAny(lower, []string{"encrypt", "dlp", "tokenize", "mask"}):
		return "Data Protection"
	default:
		return "Security Control"
	}
}

// ── PHASE 12 — BLAST RADIUS MODELING ──

type BlastRadius struct {
	ComponentName    string   `json:"component_name"`
	ArchitectureName string   `json:"architecture_name"`
	FailureImpact    string   `json:"failure_impact"`
	AffectedSystems  []string `json:"affected_systems"`
	Severity         string   `json:"severity"`
}

func computeBlastRadii(p *Portfolio) []BlastRadius {
	radii := make([]BlastRadius, 0)
	for _, arch := range p.Architectures {
		seen := make(map[string]bool)
		for _, t := range arch.Threats {
			for _, asset := range t.AffectedAssets {
				key := asset + "::" + arch.Name
				if seen[key] {
					continue
				}
				seen[key] = true
				impact := describeFailureImpact(asset, t)
				if t.Severity == "" {
					t.Severity = "Medium"
				}
				radii = append(radii, BlastRadius{
					ComponentName:    asset,
					ArchitectureName: arch.Name,
					FailureImpact:    impact,
					AffectedSystems:  findAffectedSystems(&arch, asset),
					Severity:         string(t.Severity),
				})
			}
		}
	}
	sort.Slice(radii, func(i, j int) bool {
		return severityWeight(radii[i].Severity) > severityWeight(radii[j].Severity)
	})
	if len(radii) > 20 {
		radii = radii[:20]
	}
	return radii
}

func describeFailureImpact(asset string, threat Threat) string {
	return fmt.Sprintf("If %s fails, %s could impact %s", asset, threat.Name, strings.Join(threat.AffectedComponents, ", "))
}

func findAffectedSystems(arch *ArchitectureRecord, asset string) []string {
	sys := make([]string, 0)
	for _, t := range arch.Threats {
		for _, a := range t.AffectedComponents {
			if strings.EqualFold(a, asset) {
				sys = append(sys, t.Name)
			}
		}
	}
	return sys
}

// ── PHASE 13 — ENTERPRISE COMPLIANCE VIEW ──

type EnterpriseComplianceView struct {
	Frameworks         []ComplianceFrameworkCoverage `json:"frameworks"`
	TotalArchitectures int                           `json:"total_architectures"`
}

type ComplianceFrameworkCoverage struct {
	Framework          string  `json:"framework"`
	ArchitecturesWith  int     `json:"architectures_with"`
	TotalArchitectures int     `json:"total_architectures"`
	Coverage           float64 `json:"coverage"`
}

func computeEnterpriseCompliance(p *Portfolio) EnterpriseComplianceView {
	fwIndex := make(map[string]int)
	total := len(p.Architectures)
	for _, arch := range p.Architectures {
		for _, fw := range arch.Compliance {
			key := strings.ToUpper(strings.TrimSpace(fw))
			fwIndex[key]++
		}
	}
	fws := make([]ComplianceFrameworkCoverage, 0, len(fwIndex))
	for fw, count := range fwIndex {
		cov := 0.0
		if total > 0 {
			cov = float64(count) / float64(total) * 100
		}
		fws = append(fws, ComplianceFrameworkCoverage{
			Framework:          fw,
			ArchitecturesWith:  count,
			TotalArchitectures: total,
			Coverage:           cov,
		})
	}
	sort.Slice(fws, func(i, j int) bool {
		return fws[i].Coverage > fws[j].Coverage
	})
	return EnterpriseComplianceView{
		Frameworks:         fws,
		TotalArchitectures: total,
	}
}

// ── PHASE 14 — PORTFOLIO DASHBOARD ──

type PortfolioDashboard struct {
	TotalArchitectures int            `json:"total_architectures"`
	TotalFindings      int            `json:"total_findings"`
	TotalThreats       int            `json:"total_threats"`
	TotalAttackPaths   int            `json:"total_attack_paths"`
	TotalControls      int            `json:"total_controls"`
	AverageCoverage    float64        `json:"average_coverage"`
	AverageRiskScore   float64        `json:"average_risk_score"`
	ComplianceCount    int            `json:"compliance_count"`
	RiskDistribution   map[string]int `json:"risk_distribution"`
}

func buildPortfolioDashboard(p *Portfolio, result *SAMPIResult) PortfolioDashboard {
	dash := PortfolioDashboard{
		TotalArchitectures: len(p.Architectures),
		RiskDistribution:   make(map[string]int),
	}
	totalCoverage := 0.0
	totalRisk := 0.0
	for _, arch := range p.Architectures {
		dash.TotalFindings += len(arch.Findings)
		dash.TotalThreats += len(arch.Threats)
		dash.TotalAttackPaths += len(arch.AttackPaths)
		dash.TotalControls += len(arch.Controls)
		totalRisk += arch.RiskScore
		if len(arch.Compliance) > 0 {
			dash.ComplianceCount++
		}
	}
	if len(p.Architectures) > 0 {
		dash.AverageRiskScore = totalRisk / float64(len(p.Architectures))
	}
	if dash.TotalControls > 0 && len(p.Architectures) > 0 {
		totalCoverage = float64(dash.TotalControls) / float64(len(p.Architectures))
	}
	dash.AverageCoverage = totalCoverage
	for _, rw := range result.RepeatedWeaknesses {
		sev := strings.ToLower(rw.Severity)
		if sev == "" {
			sev = "unknown"
		}
		dash.RiskDistribution[sev]++
	}
	return dash
}

// ── PHASE 15 — EXECUTIVE HEATMAPS ──

type ExecutiveHeatmap struct {
	ArchitectureName string  `json:"architecture_name"`
	RiskScore        float64 `json:"risk_score"`
	FindingCount     int     `json:"finding_count"`
	ControlCount     int     `json:"control_count"`
	ComplianceCount  int     `json:"compliance_count"`
	RiskBand         string  `json:"risk_band"`
}

func generateHeatmaps(p *Portfolio, result *SAMPIResult) []ExecutiveHeatmap {
	heatmaps := make([]ExecutiveHeatmap, 0, len(p.Architectures))
	for _, arch := range p.Architectures {
		band := riskBand(arch.RiskScore)
		heatmaps = append(heatmaps, ExecutiveHeatmap{
			ArchitectureName: arch.Name,
			RiskScore:        arch.RiskScore,
			FindingCount:     len(arch.Findings),
			ControlCount:     len(arch.Controls),
			ComplianceCount:  len(arch.Compliance),
			RiskBand:         band,
		})
	}
	sort.Slice(heatmaps, func(i, j int) bool {
		return heatmaps[i].RiskScore > heatmaps[j].RiskScore
	})
	return heatmaps
}

func riskBand(score float64) string {
	switch {
	case score >= 8:
		return "Critical"
	case score >= 5:
		return "High"
	case score >= 3:
		return "Medium"
	default:
		return "Low"
	}
}

// ── PHASE 16 — SECURITY PROGRAM INSIGHTS ──

type SecurityProgramInsight struct {
	Area      string `json:"area"`
	Insight   string `json:"insight"`
	Priority  string `json:"priority"`
	Rationale string `json:"rationale"`
}

func generateProgramInsights(result *SAMPIResult) []SecurityProgramInsight {
	insights := make([]SecurityProgramInsight, 0)

	if len(result.RepeatedWeaknesses) > 0 {
		insights = append(insights, SecurityProgramInsight{
			Area:      "Systemic Weakness Remediation",
			Insight:   fmt.Sprintf("Found %d repeated weaknesses across the portfolio", len(result.RepeatedWeaknesses)),
			Priority:  "High",
			Rationale: "Repeated weaknesses indicate systemic issues requiring enterprise-wide remediation programs",
		})
	}

	for _, theme := range result.EnterpriseRiskThemes {
		if theme.RiskCount >= 5 {
			insights = append(insights, SecurityProgramInsight{
				Area:      theme.Name,
				Insight:   fmt.Sprintf("%d findings across %d architectures", theme.RiskCount, theme.AffectedArchitectures),
				Priority:  theme.Severity,
				Rationale: fmt.Sprintf("Enterprise-wide %s requires coordinated investment", theme.Name),
			})
		}
	}

	if result.SecurityDebt.Score > 0 {
		insights = append(insights, SecurityProgramInsight{
			Area:      "Security Debt Reduction",
			Insight:   fmt.Sprintf("Security debt score: %.1f — %d longstanding, %d repeated findings", result.SecurityDebt.Score, result.SecurityDebt.LongstandingCount, result.SecurityDebt.RepeatedCount),
			Priority:  "High",
			Rationale: "Security debt erodes trust and increases breach probability",
		})
	}

	if len(result.ControlCoverage) > 0 {
		lowCoverage := 0
		total := len(result.ControlCoverage)
		for _, cc := range result.ControlCoverage {
			if cc.CoveragePercent <= 50 {
				lowCoverage++
			}
		}
		if lowCoverage > 0 {
			insights = append(insights, SecurityProgramInsight{
				Area:      "Control Coverage Gaps",
				Insight:   fmt.Sprintf("%d of %d controls have below 50%% coverage", lowCoverage, total),
				Priority:  "High",
				Rationale: "Low control coverage represents blind spots in the security program",
			})
		}
	}

	if result.AttackSurface.TotalExposure > 0 {
		insights = append(insights, SecurityProgramInsight{
			Area:      "Attack Surface Reduction",
			Insight:   fmt.Sprintf("Portfolio attack surface: %d exposures (internet: %d, third-party: %d, identity: %d, cloud: %d)", result.AttackSurface.TotalExposure, result.AttackSurface.InternetExposure, result.AttackSurface.ThirdPartyExposure, result.AttackSurface.IdentityExposure, result.AttackSurface.CloudExposure),
			Priority:  "Medium",
			Rationale: "Reducing attack surface lowers the probability of successful breaches",
		})
	}

	if len(result.SharedDependencies) > 0 {
		criticalDeps := 0
		for _, sd := range result.SharedDependencies {
			if sd.RiskLevel == "Critical" || sd.RiskLevel == "High" {
				criticalDeps++
			}
		}
		if criticalDeps > 0 {
			insights = append(insights, SecurityProgramInsight{
				Area:      "Shared Dependency Risk",
				Insight:   fmt.Sprintf("%d shared dependencies present concentration risk across %d architectures", criticalDeps, result.Dashboard.TotalArchitectures),
				Priority:  "Medium",
				Rationale: "Shared dependencies create single points of failure for the enterprise",
			})
		}
	}

	return insights
}

// ── CONVENIENCE ──

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
