package assumption

import (
	"regexp"
	"strings"

	"asf-tui/asf/models"
)

type typePattern struct {
	atype    models.AssumptionType
	patterns []*regexp.Regexp
}

var typePatterns = []typePattern{
	{
		atype: models.AssumptionTypeIDENTITY,
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\b(?:mfa|multi.?factor|identity|authentication|password|credential|role|group)\b`),
			regexp.MustCompile(`(?i)\b(?:only|just)\s+.+\s+(?:can|may|has)\s+.+\b(?:access|login|authenticate)\b`),
		},
	},
	{
		atype: models.AssumptionTypeACCESS,
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\b(?:access|permission|acl|allow|deny|grant|read|write|execute|admin)\b`),
			regexp.MustCompile(`(?i)\b(?:only|just)\s+.+\s+(?:can|may|has)\s+(?:access|permission)\b`),
			regexp.MustCompile(`(?i)\b(?:restricted?|limited?|blocked?)\s+(?:to|access)\b`),
		},
	},
	{
		atype: models.AssumptionTypeNETWORK,
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\b(?:network|firewall|internet|subnet|vpc|vlan|segment|isolate|expose)\b`),
			regexp.MustCompile(`(?i)\b(?:not\s+)?(?:accessible|reachable)\s+(?:from|via|over)\b`),
			regexp.MustCompile(`(?i)\b(?:no\s+)?public\s+(?:access|exposure|facing)\b`),
		},
	},
	{
		atype: models.AssumptionTypeCONFIGURATION,
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\b(?:encrypt|backup|log|audit|monitor|config|setting|parameter)\b`),
			regexp.MustCompile(`(?i)\b(?:encrypt(?:ed|ion)|backup|log(?:ging|ged)|audit(?:ing|ed)?|monitor(?:ing|ed)?)\b`),
			regexp.MustCompile(`(?i)\b(?:enabled?|disabled?|configured?)\s+(?:by|with|to|as)\b`),
		},
	},
	{
		atype: models.AssumptionTypePROCESS,
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\b(?:process|procedure|workflow|review|approve|test|sign.?off|approval)\b`),
			regexp.MustCompile(`(?i)\b(?:must|shall|should)\s+(?:be\s+)?(?:reviewed|tested|approved|validated)\b`),
		},
	},
	{
		atype: models.AssumptionTypeDOCUMENTATION,
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\b(?:document|policy|runbook|procedure|guide|manual|readme|wiki)\b`),
			regexp.MustCompile(`(?i)\b(?:as\s+(?:per|described|documented|stated))\b`),
		},
	},
	{
		atype: models.AssumptionTypeDEPENDENCY,
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\b(?:depend|integration|connect|communicate|rel(y|ies)|upstream|downstream)\b`),
			regexp.MustCompile(`(?i)\b(?:requires?|depends?\s+on|relies?\s+on)\b`),
		},
	},
	{
		atype: models.AssumptionTypeGOVERNANCE,
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`(?i)\b(?:review|audit|compliance|regulat|policy|standard|framework|govern)\b`),
			regexp.MustCompile(`(?i)\b(?:annually|quarterly|monthly|regularly|periodically)\b`),
			regexp.MustCompile(`(?i)\b(?:reviewed|audited|approved|certified)\s+(?:annually|quarterly|monthly|regularly|periodically)\b`),
		},
	},
}

var stopwords = map[string]bool{
	"the": true, "and": true, "for": true, "are": true, "but": true,
	"not": true, "you": true, "all": true, "can": true, "has": true,
	"have": true, "may": true, "must": true, "shall": true, "should": true,
	"will": true, "with": true, "from": true, "that": true, "this": true,
	"each": true, "every": true, "than": true, "then": true, "just": true,
	"been": true, "were": true, "was": true, "its": true, "also": true,
	"per": true, "via": true,
}

type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

func (ae *Engine) Convert(claim models.Claim) *models.Assumption {
	atype := classify(claim.Text)
	if atype == nil {
		return nil
	}

	text := buildAssumptionText(claim.Text, *atype)
	keywords := extractKeywords(claim.Text)

	assumption := models.NewAssumption(claim.ID, text, *atype, keywords)
	return &assumption
}

func (ae *Engine) ConvertMany(claims []models.Claim) []models.Assumption {
	var result []models.Assumption
	for _, c := range claims {
		a := ae.Convert(c)
		if a != nil {
			result = append(result, *a)
		}
	}
	return result
}

func classify(text string) *models.AssumptionType {
	scores := make(map[models.AssumptionType]int)
	for _, tp := range typePatterns {
		score := 0
		for _, pat := range tp.patterns {
			matches := pat.FindAllString(text, -1)
			score += len(matches)
		}
		if score > 0 {
			scores[tp.atype] = score
		}
	}

	if len(scores) == 0 {
		return nil
	}

	var best models.AssumptionType
	bestScore := 0
	for _, tp := range typePatterns {
		sc := scores[tp.atype]
		if sc > bestScore {
			best = tp.atype
			bestScore = sc
		}
	}
	return &best
}

func buildAssumptionText(text string, atype models.AssumptionType) string {
	prefixes := map[models.AssumptionType]string{
		models.AssumptionTypeIDENTITY:      "System assumes identity posture: ",
		models.AssumptionTypeACCESS:        "System assumes access control: ",
		models.AssumptionTypeNETWORK:       "System assumes network posture: ",
		models.AssumptionTypeCONFIGURATION: "System assumes configuration state: ",
		models.AssumptionTypePROCESS:       "System assumes process compliance: ",
		models.AssumptionTypeDOCUMENTATION: "System assumes documentation accuracy: ",
		models.AssumptionTypeDEPENDENCY:    "System assumes dependency relationship: ",
		models.AssumptionTypeGOVERNANCE:    "System assumes governance compliance: ",
	}
	prefix := prefixes[atype]
	if prefix == "" {
		prefix = "System assumes: "
	}
	return prefix + text
}

func extractKeywords(text string) []string {
	re := regexp.MustCompile(`\b[a-zA-Z]{3,}\b`)
	words := re.FindAllString(strings.ToLower(text), -1)
	var result []string
	for _, w := range words {
		if !stopwords[w] {
			result = append(result, w)
		}
	}
	return result
}
