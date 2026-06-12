package extraction

import (
	"regexp"
	"strings"

	"asf-tui/asf/models"
)

var declarativePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)(?:only|just)\s+.+\s+(?:can|may|has|have|should|must|will)`),
	regexp.MustCompile(`(?i)(?:all|every|each)\s+.+\s+(?:are|is|shall|must|will|should)`),
	regexp.MustCompile(`(?i)(?:is|are)\s+(?:not|never)\s+.+`),
	regexp.MustCompile(`(?i)(?:is|are)\s+\w+\s+(?:encrypted|backed\s+up|logged|audited|monitored|reviewed|tested|approved|validated|isolated|segmented|restricted|limited)`),
	regexp.MustCompile(`(?i)(?:cannot|can\s+not|must\s+not|shall\s+not|should\s+not)\s+.+`),
	regexp.MustCompile(`(?i)(?:requir(?:e|es|ed|ing)|ensures?|guarantees?|provides?|protects?|prevents?|blocks?|restricts?|limits?)\s+.+`),
	regexp.MustCompile(`(?i)(?:ensure[s]?|guarantee[s]?)\s+that\s+.+`),
	regexp.MustCompile(`(?i)(?:accessed?|accessible|available)\s+(?:only|exclusively|solely)\s+.+`),
	regexp.MustCompile(`(?i)(?:encrypt(?:ed|s|ion)|back(?:ed|s)?\s+ups?|backups?|log(?:ged|s|ging)?|audit(?:ed|s|ing)?|monitor(?:ed|s|ing)?)`),
	regexp.MustCompile(`(?i)(?:configured?|set\s+up|deployed?|implemented?)\s+(?:to|with|as)\s+.+`),
	regexp.MustCompile(`(?i)(?:review|test|approv|validat|certif)\w+\s+(?:are|is|shall|must|should|will|performed|conducted)`),
	regexp.MustCompile(`(?i)(?:separated?|isolated?|segmented?|partitioned?)\s+.+`),
	regexp.MustCompile(`(?i)(?:is|are)\s+restricted\s+to\s+.+`),
	regexp.MustCompile(`(?i)manage[sd]?\s+.+\s+(?:access|permissions)`),
	regexp.MustCompile(`(?i)(?:security\s+)?(?:review|audit|assessment)s?\s+.+\s+(?:are|is|conducted|performed|scheduled)`),
}

var strongIndicatorPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)\b(?:only|never|always|all|every|must|shall)\b`),
	regexp.MustCompile(`(?i)\b(?:encrypt|backup|audit|require|ensure|guarantee|manage)\b`),
	regexp.MustCompile(`(?i)\b(?:isolated|segmented|restricted|limited|conducted)\b`),
}

var keywordTagMap = map[string]string{
	"access":         "access",
	"permission":     "access",
	"identity":       "identity",
	"mfa":            "identity",
	"authentication": "identity",
	"network":        "network",
	"firewall":       "network",
	"internet":       "network",
	"encrypt":        "configuration",
	"backup":         "configuration",
	"log":            "configuration",
	"audit":          "configuration",
	"review":         "governance",
	"approve":        "governance",
	"restrict":       "access",
	"manage":         "access",
	"process":        "process",
	"procedure":      "process",
	"test":           "process",
	"document":       "documentation",
	"policy":         "documentation",
	"depend":         "dependency",
	"integration":    "dependency",
}

type ClaimExtractor struct{}

func NewClaimExtractor() *ClaimExtractor {
	return &ClaimExtractor{}
}

func (ce *ClaimExtractor) Extract(text, sourceDocument, sourceLocation string) []models.Claim {
	var claims []models.Claim
	seen := make(map[string]bool)

	sentences := splitSentences(text)
	for _, sentence := range sentences {
		cleaned := strings.TrimSpace(sentence)
		if cleaned == "" || len(cleaned) < 15 {
			continue
		}
		if isDeclarative(cleaned) {
			normalized := strings.ToLower(strings.TrimSpace(cleaned))
			if seen[normalized] {
				continue
			}
			seen[normalized] = true
			confidence := computeConfidence(cleaned)
			tags := extractTags(cleaned)
			claim := models.NewClaim(sourceDocument, sourceLocation, cleaned, confidence, tags)
			claims = append(claims, claim)
		}
	}

	return claims
}

func splitSentences(text string) []string {
	re := regexp.MustCompile(`[.!?](\s+|$)`)
	matches := re.FindAllStringIndex(text, -1)
	if len(matches) == 0 {
		text = strings.TrimSpace(text)
		if text != "" {
			return []string{text}
		}
		return nil
	}

	var result []string
	start := 0
	for _, m := range matches {
		end := m[0] + 1
		s := strings.TrimSpace(text[start:end])
		if s != "" {
			result = append(result, s)
		}
		start = m[1]
	}
	tail := strings.TrimSpace(text[start:])
	if tail != "" {
		result = append(result, tail)
	}
	return result
}

func isDeclarative(text string) bool {
	for _, pat := range declarativePatterns {
		if pat.MatchString(text) {
			return true
		}
	}
	return false
}

func computeConfidence(text string) float64 {
	score := 0.5
	for _, pat := range strongIndicatorPatterns {
		if pat.MatchString(text) {
			score += 0.1
		}
	}
	if score > 0.95 {
		return 0.95
	}
	return score
}

func extractTags(text string) []string {
	var tags []string
	seen := make(map[string]bool)
	lowered := strings.ToLower(text)
	for keyword, tag := range keywordTagMap {
		if strings.Contains(lowered, keyword) {
			if !seen[tag] {
				seen[tag] = true
				tags = append(tags, tag)
			}
		}
	}
	return tags
}
