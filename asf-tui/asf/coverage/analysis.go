package coverage

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"time"
)

type AssumptionInput struct {
	ID          string
	Description string
	Component   string
	Category    string
	Keywords    []string
	Risk        string
}

type CoverageEngine struct {
	domain      string
	components  []string
	assumptions []AssumptionInput
}

func NewCoverageEngine(domain string, components []string, assumptions []AssumptionInput) *CoverageEngine {
	return &CoverageEngine{
		domain:      domain,
		components:  components,
		assumptions: assumptions,
	}
}

func (e *CoverageEngine) RunAll() *CoverageOutput {
	output := &CoverageOutput{
		Domain:      e.domain,
		GeneratedAt: time.Now().UTC().Format(time.RFC3339),
	}

	assessment := e.assessCoverage()
	output.Assessment = assessment

	output.BlindSpots = e.detectBlindSpots(assessment)
	output.DomainBlindSpots = GetDomainBlindSpots(e.domain)

	output.AttentionScore = e.computeAttentionScore(assessment)
	output.CISOView = e.buildCISOView(assessment, output.BlindSpots, output.DomainBlindSpots)

	return output
}

func (e *CoverageEngine) assessCoverage() *CoverageAssessment {
	catCounts := make(map[CoverageCategory]int)
	riskByCat := make(map[CoverageCategory]string)

	for _, a := range e.assumptions {
		cat := CoverageCategory(strings.ToLower(a.Category))
		catCounts[cat]++
		if a.Risk == "Critical" || a.Risk == "High" {
			if riskByCat[cat] == "" || a.Risk == "Critical" {
				riskByCat[cat] = a.Risk
			}
		}
	}

	var compResults []ComponentExpectations
	compSeen := make(map[string]bool)

	for _, comp := range e.components {
		cl := comp
		expectations := GetExpectations(cl)
		if len(expectations) == 0 {
			continue
		}
		if compSeen[cl] {
			continue
		}
		compSeen[cl] = true

		matched := filterMatchingExpectations(expectations, e.assumptions, cl)
		compResults = append(compResults, ComponentExpectations{
			Component:    cl,
			Expectations: matched,
		})
	}

	allCats := make(map[CoverageCategory]*CoverageMetric)
	for _, cat := range AllCategories {
		allCats[cat] = &CoverageMetric{
			Category: cat,
			Risk:     "Medium",
		}
	}

	for _, cr := range compResults {
		for _, exp := range cr.Expectations {
			m := allCats[exp.Category]
			m.ExpectedCount++
			if m.Risk == "Medium" && exp.Risk == "High" {
				m.Risk = "High"
			}
			if exp.Risk == "Critical" {
				m.Risk = "Critical"
			}
		}
	}

	for cat, count := range catCounts {
		if m, ok := allCats[cat]; ok {
			m.ObservedCount = count
		}
	}
	if rc := riskByCat[CatIdentity]; rc != "" {
		allCats[CatIdentity].Risk = rc
	}
	if rc := riskByCat[CatAuthorization]; rc != "" {
		allCats[CatAuthorization].Risk = rc
	}
	if rc := riskByCat[CatCryptography]; rc != "" {
		allCats[CatCryptography].Risk = rc
	}
	if rc := riskByCat[CatMonitoring]; rc != "" {
		allCats[CatMonitoring].Risk = rc
	}
	if rc := riskByCat[CatResilience]; rc != "" {
		allCats[CatResilience].Risk = rc
	}

	var metrics []CoverageMetric
	var gaps []CoverageGap

	for _, cat := range AllCategories {
		m := allCats[cat]
		if m.ExpectedCount == 0 {
			continue
		}
		m.CoveragePct = coveragePct(m.ObservedCount, m.ExpectedCount)

		expRisk := riskLevel(m.CoveragePct)
		if riskWeight(expRisk) > riskWeight(m.Risk) {
			m.Risk = expRisk
		}

		delta := attentionDelta(m.CoveragePct)
		m.AttentionDelta = delta
		m.Reason = attentionReason(cat, m.CoveragePct, m.ObservedCount, m.ExpectedCount)

		metrics = append(metrics, *m)

		if m.CoveragePct < 80.0 {
			gaps = append(gaps, CoverageGap{
				Category:       cat,
				ExpectedCount:  m.ExpectedCount,
				ObservedCount:  m.ObservedCount,
				MissingCount:   m.ExpectedCount - m.ObservedCount,
				CoveragePct:    m.CoveragePct,
				Risk:           m.Risk,
				Recommendation: gapRecommendation(cat),
			})
		}
	}

	sort.Slice(metrics, func(i, j int) bool {
		return metrics[i].CoveragePct < metrics[j].CoveragePct
	})

	return &CoverageAssessment{
		Categories:       metrics,
		ComponentResults: compResults,
		Gaps:             gaps,
	}
}

func (e *CoverageEngine) detectBlindSpots(assessment *CoverageAssessment) []BlindSpot {
	var spots []BlindSpot
	seen := make(map[string]bool)

	type blindSpotRule struct {
		trigger     func() bool
		category    CoverageCategory
		title       string
		description string
		risk        string
		score       float64
		component   string
	}

	var rules []blindSpotRule

	compLower := make(map[string]bool)
	for _, c := range e.components {
		cl := strings.ToLower(c)
		compLower[cl] = true
	}

	hasKeyword := func(c string, keywords []string) bool {
		cl := strings.ToLower(c)
		for _, kw := range keywords {
			if strings.Contains(cl, kw) {
				return true
			}
		}
		return false
	}

	hasCatAssumption := func(cat CoverageCategory) bool {
		for _, a := range e.assumptions {
			if strings.ToLower(a.Category) == string(cat) {
				return true
			}
		}
		return false
	}

	hasComponentAssumption := func(component string, keywords []string) bool {
		for _, a := range e.assumptions {
			clComp := strings.ToLower(a.Component)
			clDesc := strings.ToLower(a.Description)
			for _, kw := range keywords {
				if strings.Contains(clComp, kw) || strings.Contains(clDesc, kw) {
					return true
				}
			}
		}
		return false
	}

	keyManagementKeywords := []string{"key", "kms", "rotation", "vault"}
	for cl := range compLower {
		if hasKeyword(cl, []string{"kms", "key", "vault", "hsm"}) {
			rules = append(rules, blindSpotRule{
				func() bool { return !hasComponentAssumption(cl, keyManagementKeywords) },
				CatCryptography, "Key Rotation Missing",
				fmt.Sprintf("Component %s is present but no key rotation assumption found", cl),
				"Critical", 95.0, cl,
			})
			rules = append(rules, blindSpotRule{
				func() bool { return !hasCatAssumption(CatCryptography) },
				CatCryptography, "Encryption Missing",
				fmt.Sprintf("Component %s is present but no encryption assumption found", cl),
				"Critical", 90.0, cl,
			})
		}
		if hasKeyword(cl, []string{"auth0", "okta", "idp", "identity"}) {
			rules = append(rules, blindSpotRule{
				func() bool { return !hasComponentAssumption(cl, []string{"mfa"}) },
				CatIdentity, "MFA Missing",
				fmt.Sprintf("Identity provider %s is present but no MFA assumption found", cl),
				"Critical", 95.0, cl,
			})
			rules = append(rules, blindSpotRule{
				func() bool { return !hasComponentAssumption(cl, []string{"admin", "restrict"}) },
				CatIdentity, "Admin Access Restrictions Missing",
				fmt.Sprintf("Identity provider %s is present but no admin restriction assumption found", cl),
				"Critical", 90.0, cl,
			})
		}
		if hasKeyword(cl, []string{"backup", "snapshot"}) {
			rules = append(rules, blindSpotRule{
				func() bool { return !hasComponentAssumption(cl, []string{"restore"}) },
				CatResilience, "Restore Testing Missing",
				fmt.Sprintf("Backup component %s is present but no restore testing assumption found", cl),
				"High", 88.0, cl,
			})
		}
		if hasKeyword(cl, []string{"apigateway", "api", "gateway"}) {
			rules = append(rules, blindSpotRule{
				func() bool { return !hasComponentAssumption(cl, []string{"rate", "limit", "throttle"}) },
				CatOperational, "Rate Limiting Missing",
				fmt.Sprintf("API gateway %s is present but no rate limiting assumption found", cl),
				"High", 82.0, cl,
			})
		}
		if hasKeyword(cl, []string{"siem", "splunk", "elastic", "cloudtrail"}) {
			rules = append(rules, blindSpotRule{
				func() bool { return !hasComponentAssumption(cl, []string{"alert"}) },
				CatMonitoring, "Alerting Missing",
				fmt.Sprintf("Monitoring component %s is present but no alerting assumption found", cl),
				"High", 85.0, cl,
			})
		}
		if hasKeyword(cl, []string{"jenkins", "ci", "pipeline", "gitlab"}) {
			rules = append(rules, blindSpotRule{
				func() bool { return !hasComponentAssumption(cl, []string{"secret", "credential"}) },
				CatOperational, "Secrets Management Missing",
				fmt.Sprintf("CI/CD component %s is present but no secrets management assumption found", cl),
				"Critical", 92.0, cl,
			})
		}
		if hasKeyword(cl, []string{"database", "db", "rds"}) {
			rules = append(rules, blindSpotRule{
				func() bool { return !hasComponentAssumption(cl, []string{"encrypt"}) },
				CatCryptography, "Database Encryption Missing",
				fmt.Sprintf("Database %s is present but no encryption assumption found", cl),
				"Critical", 93.0, cl,
			})
		}
	}

	blindID := 0
	for _, rule := range rules {
		if rule.trigger() {
			blindID++
			id := fmt.Sprintf("BS-%03d", blindID)
			if seen[id] {
				continue
			}
			seen[id] = true
			spots = append(spots, BlindSpot{
				Category:            rule.category,
				BlindSpotID:         id,
				Title:               rule.title,
				Description:         rule.description,
				Risk:                rule.risk,
				Score:               rule.score,
				Component:           rule.component,
				Domain:              e.domain,
				TrustChainImpact:    trustChainImpact(rule.category),
				ConsequenceSeverity: consequenceSeverity(rule.risk),
				ComplianceRelevance: complianceRelevance(rule.category),
				Recommendation:      blindSpotRecommendation(rule.title, rule.component),
			})
		}
	}

	sort.Slice(spots, func(i, j int) bool {
		return spots[i].Score > spots[j].Score
	})

	return spots
}

func (e *CoverageEngine) computeAttentionScore(assessment *CoverageAssessment) float64 {
	if len(assessment.Categories) == 0 {
		return 100.0
	}

	totalWeight := 0.0
	weightedScore := 0.0

	for _, cat := range assessment.Categories {
		w := riskWeight(cat.Risk)
		if w == 0 {
			w = 1.0
		}
		totalWeight += w
		weightedScore += w * cat.CoveragePct
	}

	rawScore := 0.0
	if totalWeight > 0 {
		rawScore = weightedScore / totalWeight
	}

	penalty := 0.0
	for _, gap := range assessment.Gaps {
		if gap.Risk == "Critical" {
			penalty += 5.0
		} else if gap.Risk == "High" {
			penalty += 3.0
		}
	}
	if penalty > 30.0 {
		penalty = 30.0
	}

	final := math.Round(rawScore - penalty)
	if final < 0 {
		final = 0
	}
	if final > 100 {
		final = 100
	}
	return final
}

func (e *CoverageEngine) buildCISOView(assessment *CoverageAssessment, blindSpots []BlindSpot, domainSpots []DomainBlindSpot) *CISOView {
	if assessment == nil {
		return nil
	}

	topN := 10
	if len(blindSpots) < topN {
		topN = len(blindSpots)
	}
	var topBlindSpots []BlindSpot
	if topN > 0 {
		topBlindSpots = blindSpots[:topN]
	}

	var dangerous []BlindSpot
	for _, bs := range blindSpots {
		if bs.Risk == "Critical" && len(dangerous) < 5 {
			dangerous = append(dangerous, bs)
		}
	}

	var areas []string
	for _, gap := range assessment.Gaps {
		if gap.Risk == "Critical" || gap.Risk == "High" {
			areas = append(areas, fmt.Sprintf("%s (%.0f%% coverage)", gap.Category, gap.CoveragePct))
		}
	}

	for _, ds := range domainSpots {
		if ds.Risk == "Critical" || ds.Risk == "High" {
			areas = append(areas, fmt.Sprintf("%s: %s", ds.Domain, ds.MissingArea))
		}
	}

	var highestRisk []BlindSpot
	for _, bs := range blindSpots {
		if bs.Score >= 90 && len(highestRisk) < 5 {
			highestRisk = append(highestRisk, bs)
		}
	}
	if len(highestRisk) == 0 && len(blindSpots) > 0 {
		highestRisk = append(highestRisk, blindSpots[0])
	}

	return &CISOView{
		TopBlindSpots:               topBlindSpots,
		DangerousMissingAssumptions: dangerous,
		AreasRequiringReview:        areas,
		HighestRiskUnknowns:         highestRisk,
	}
}

func filterMatchingExpectations(expectations []ExpectedAssumption, assumptions []AssumptionInput, component string) []ExpectedAssumption {
	var matched []ExpectedAssumption
	cl := strings.ToLower(component)

	for _, exp := range expectations {
		keywords := keywordsForTitle(exp.Title)
		found := false
		for _, a := range assumptions {
			al := strings.ToLower(a.Component)
			if !strings.Contains(al, cl) && !strings.Contains(cl, al) {
				continue
			}
			desc := strings.ToLower(a.Description)
			for _, kw := range keywords {
				if strings.Contains(desc, kw) {
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if !found {
			matched = append(matched, exp)
		}
	}

	return matched
}

func coveragePct(observed, expected int) float64 {
	if expected == 0 {
		return 100.0
	}
	pct := float64(observed) / float64(expected) * 100.0
	if pct > 100 {
		pct = 100
	}
	return math.Round(pct*10) / 10
}

func riskLevel(coverage float64) string {
	switch {
	case coverage >= 90:
		return "Low"
	case coverage >= 70:
		return "Medium"
	case coverage >= 50:
		return "High"
	default:
		return "Critical"
	}
}

func riskWeight(risk string) float64 {
	switch risk {
	case "Critical":
		return 5.0
	case "High":
		return 3.0
	case "Medium":
		return 2.0
	default:
		return 1.0
	}
}

func attentionDelta(pct float64) float64 {
	switch {
	case pct >= 90:
		return 5.0
	case pct >= 70:
		return 3.0
	case pct >= 50:
		return 1.0
	default:
		return 0.5
	}
}

func attentionReason(cat CoverageCategory, pct float64, observed, expected int) string {
	switch {
	case pct >= 90:
		return fmt.Sprintf("Good coverage in %s (%d/%d)", cat, observed, expected)
	case pct >= 70:
		return fmt.Sprintf("Moderate coverage in %s (%d/%d), investigate gaps", cat, observed, expected)
	case pct >= 50:
		return fmt.Sprintf("Low coverage in %s (%d/%d), needs attention", cat, observed, expected)
	default:
		return fmt.Sprintf("Critical gap in %s (%d/%d), immediate review needed", cat, observed, expected)
	}
}

func gapRecommendation(cat CoverageCategory) string {
	switch cat {
	case CatIdentity:
		return "Review identity controls: MFA, SSO, federation, and admin access restrictions"
	case CatAuthorization:
		return "Review authorization: RBAC, least privilege, API authentication, and access policies"
	case CatCryptography:
		return "Review cryptographic controls: encryption at rest/transit, key rotation, and KMS access"
	case CatMonitoring:
		return "Review monitoring: SIEM integration, audit logging, alerting, and log retention"
	case CatResilience:
		return "Review resilience: backup schedule, restore testing, DR plan, and backup encryption"
	case CatThirdParty:
		return "Review third-party dependencies: vendor risk, supply chain, and provider trust"
	case CatOperational:
		return "Review operational controls: secrets management, patch management, and configuration"
	default:
		return "Review missing assumptions in this category"
	}
}

func trustChainImpact(cat CoverageCategory) string {
	switch cat {
	case CatIdentity:
		return "Identity failures cascade to all downstream authentication decisions"
	case CatAuthorization:
		return "Authorization gaps cascade to privilege escalation risks"
	case CatCryptography:
		return "Cryptographic failures cascade to data confidentiality and integrity"
	case CatMonitoring:
		return "Monitoring gaps cascade to delayed incident detection"
	case CatResilience:
		return "Resilience gaps cascade to extended recovery times"
	default:
		return "Gaps in this category affect overall security posture"
	}
}

func consequenceSeverity(risk string) string {
	switch risk {
	case "Critical":
		return "Severe"
	case "High":
		return "Major"
	case "Medium":
		return "Moderate"
	default:
		return "Minor"
	}
}

func complianceRelevance(cat CoverageCategory) string {
	switch cat {
	case CatIdentity:
		return "SOX, HIPAA, PCI-DSS, SOC2"
	case CatAuthorization:
		return "SOX, HIPAA, PCI-DSS, SOC2, FedRAMP"
	case CatCryptography:
		return "PCI-DSS, HIPAA, FedRAMP, GDPR"
	case CatMonitoring:
		return "HIPAA, PCI-DSS, SOC2, FedRAMP"
	case CatResilience:
		return "SOC2, PCI-DSS, HIPAA"
	default:
		return "General security compliance"
	}
}

func blindSpotRecommendation(title, component string) string {
	switch title {
	case "Key Rotation Missing":
		return "Enable automatic key rotation in " + component
	case "MFA Missing":
		return "Enable MFA for all users in " + component
	case "Admin Access Restrictions Missing":
		return "Restrict administrative access in " + component + " to authorized users only"
	case "Restore Testing Missing":
		return "Implement and test restore procedures for " + component
	case "Rate Limiting Missing":
		return "Configure rate limiting for API gateway " + component
	case "Alerting Missing":
		return "Configure security alerts for " + component
	case "Secrets Management Missing":
		return "Integrate secrets management solution for " + component + " pipelines"
	case "Database Encryption Missing":
		return "Enable encryption at rest and in transit for " + component
	case "Encryption Missing":
		return "Implement encryption controls for " + component
	default:
		return "Review and address " + title + " for " + component
	}
}

func keywordsForTitle(title string) []string {
	kwMap := map[string][]string{
		"MFA Enforcement":          {"mfa", "multi-factor", "2fa", "two-factor"},
		"SSO Configuration":        {"sso", "single sign-on", "federat"},
		"Federation":               {"federat", "saml", "oidc", "identity provider"},
		"Admin Access Restriction": {"admin", "administrative", "privileged"},
		"RBAC Configuration":       {"rbac", "role", "role-based"},
		"Least Privilege":          {"least privilege", "principle of least"},
		"Access Audit Logging":     {"audit", "log", "access log"},
		"Database Access Control":  {"database access", "db access", "restrict"},
		"Encryption at Rest":       {"encrypt", "at rest", "encryption"},
		"Encryption in Transit":    {"tls", "ssl", "in transit", "https"},
		"Database Audit Logging":   {"audit", "database log", "db audit"},
		"Database Backups":         {"backup", "snapshot"},
		"Restore Testing":          {"restore", "recovery", "disaster"},
		"Key Rotation":             {"rotation", "key rotation", "rotate"},
		"Key Access Control":       {"key access", "kms access", "key restrict"},
		"Key Backup":               {"key backup", "key recovery"},
		"KMS Access Policy":        {"kms policy", "key policy", "kms access"},
		"Key Usage Auditing":       {"key audit", "key usage", "kms audit"},
		"Centralized Logging":      {"centralized", "log aggregation", "log collection"},
		"Alerting":                 {"alert", "notif", "incident"},
		"Log Retention":            {"retention", "log retention", "retain"},
		"Log Integrity":            {"immutable", "tamper", "log integrity", "write-once"},
		"API Authentication":       {"api auth", "api key", "token", "jwt"},
		"API Authorization":        {"api authz", "api policy", "api permission"},
		"API Audit Logging":        {"api log", "api audit"},
		"Rate Limiting":            {"rate limit", "throttle", "quota"},
		"Bucket Access Control":    {"bucket policy", "s3 policy", "access control", "block public"},
		"Data Backup":              {"backup", "snapshot", "replication"},
		"Backup Schedule":          {"backup schedule", "backup frequency", "regular backup"},
		"Backup Encryption":        {"backup encrypt", "backup crypto"},
		"Disaster Recovery":        {"dr", "disaster recovery", "failover", "recovery plan"},
		"Secrets Management":       {"secret", "credential", "vault", "token"},
		"Pipeline Security":        {"pipeline", "ci/cd", "build"},
		"Artifact Signing":         {"sign", "artifact", "verify"},
		"Authentication":           {"auth", "login", "authenticat"},
		"Session Management":       {"session", "cookie", "token"},
		"TLS Configuration":        {"tls", "ssl", "https", "certificate"},
	}
	if kw, ok := kwMap[title]; ok {
		return kw
	}
	return []string{strings.ToLower(title)}
}
