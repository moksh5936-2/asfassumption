package main

import (
	"strings"
)

type StrideEngine struct {
	categoryRules map[string][]StrideCategory
	keywordRules  []keywordRule
}

type keywordRule struct {
	keywords []string
	stride   []StrideCategory
}

func (se *StrideEngine) GetKeywordRules() []keywordRule {
	return se.keywordRules
}

func (se *StrideEngine) GetCategoryRules() map[string][]StrideCategory {
	return se.categoryRules
}

func NewStrideEngine() *StrideEngine {
	return &StrideEngine{
		categoryRules: buildCategoryRules(),
		keywordRules:  buildKeywordRules(),
	}
}

func (se *StrideEngine) MapAssumption(category string, text string, keywords []string) []StrideCategory {
	seen := make(map[StrideCategory]bool)
	var result []StrideCategory

	if cats, ok := se.categoryRules[category]; ok {
		for _, s := range cats {
			if !seen[s] {
				seen[s] = true
				result = append(result, s)
			}
		}
	}

	searchText := strings.ToLower(text)
	for _, kw := range keywords {
		searchText += " " + strings.ToLower(kw)
	}
	for _, rule := range se.keywordRules {
		for _, kw := range rule.keywords {
			if strings.Contains(searchText, kw) {
				for _, s := range rule.stride {
					if !seen[s] {
						seen[s] = true
						result = append(result, s)
					}
				}
				break
			}
		}
	}

	return result
}

func buildCategoryRules() map[string][]StrideCategory {
	return map[string][]StrideCategory{
		"IDENTITY":       {StrideSpoofing, StrideElevationPriv},
		"AUTHENTICATION": {StrideSpoofing, StrideElevationPriv},
		"AUTHORIZATION":  {StrideElevationPriv, StrideInfoDisclosure},
		"ACCESS":         {StrideElevationPriv, StrideInfoDisclosure},
		"NETWORK":        {StrideInfoDisclosure, StrideDenialOfService, StrideTampering},
		"ENCRYPTION":     {StrideInfoDisclosure},
		"CONFIGURATION":  {StrideTampering},
		"DEPENDENCY":     {StrideDenialOfService, StrideTampering},
		"PROCESS":        {StrideRepudiation, StrideTampering},
		"DATABASE":       {StrideTampering, StrideInfoDisclosure},
		"LOGGING":        {StrideRepudiation, StrideTampering},
		"BACKUP":         {StrideInfoDisclosure, StrideDenialOfService},
		"SESSION":        {StrideSpoofing, StrideElevationPriv},
		"THIRD_PARTY":    {StrideTampering, StrideInfoDisclosure},
		"DOCUMENTATION":  {StrideRepudiation},
		"GOVERNANCE":     {StrideRepudiation, StrideTampering},
		"GENERAL":        {},
	}
}

func buildKeywordRules() []keywordRule {
	return []keywordRule{
		{keywords: []string{"idor", "insecure direct object"}, stride: []StrideCategory{StrideInfoDisclosure, StrideElevationPriv}},
		{keywords: []string{"bola", "broken object level"}, stride: []StrideCategory{StrideInfoDisclosure, StrideElevationPriv}},
		{keywords: []string{"session hijack", "session fixat", "session predict"}, stride: []StrideCategory{StrideSpoofing, StrideElevationPriv}},
		{keywords: []string{"audit log", "audit trail", "log immutable", "log tamper"}, stride: []StrideCategory{StrideRepudiation, StrideTampering}},
		{keywords: []string{"backup", "data loss", "data recover"}, stride: []StrideCategory{StrideInfoDisclosure, StrideDenialOfService}},
		{keywords: []string{"mfa", "multi factor", "two factor", "2fa"}, stride: []StrideCategory{StrideSpoofing}},
		{keywords: []string{"sql injection", "sqli", "nosql injec"}, stride: []StrideCategory{StrideTampering, StrideInfoDisclosure}},
		{keywords: []string{"key management", "key rotat", "key stor"}, stride: []StrideCategory{StrideInfoDisclosure}},
		{keywords: []string{"buffer overflow", "memory corrupt"}, stride: []StrideCategory{StrideTampering, StrideDenialOfService}},
		{keywords: []string{"cross site script", "xss"}, stride: []StrideCategory{StrideTampering, StrideInfoDisclosure}},
		{keywords: []string{"csrf", "cross site request"}, stride: []StrideCategory{StrideTampering, StrideElevationPriv}},
		{keywords: []string{"ssrf", "server side request"}, stride: []StrideCategory{StrideInfoDisclosure, StrideElevationPriv}},
		{keywords: []string{"privilege escal"}, stride: []StrideCategory{StrideElevationPriv}},
		{keywords: []string{"denial of serv", "dos", "ddos"}, stride: []StrideCategory{StrideDenialOfService}},
		{keywords: []string{"man in the middl", "mitm"}, stride: []StrideCategory{StrideSpoofing, StrideTampering, StrideInfoDisclosure}},
		{keywords: []string{"replay attack"}, stride: []StrideCategory{StrideSpoofing, StrideElevationPriv}},
		{keywords: []string{"tls", "ssl", "https"}, stride: []StrideCategory{StrideInfoDisclosure, StrideTampering}},
		{keywords: []string{"auth bypass", "authn bypass", "authentication bypass"}, stride: []StrideCategory{StrideSpoofing, StrideElevationPriv}},
		{keywords: []string{"rate limit"}, stride: []StrideCategory{StrideDenialOfService}},
		{keywords: []string{"supply chain"}, stride: []StrideCategory{StrideTampering, StrideDenialOfService}},
		{keywords: []string{"secret", "credential", "password", "token"}, stride: []StrideCategory{StrideSpoofing, StrideInfoDisclosure}},
		{keywords: []string{"firewall", "acl", "network segment"}, stride: []StrideCategory{StrideDenialOfService, StrideInfoDisclosure}},
		{keywords: []string{"encrypt", "decrypt", "cipher"}, stride: []StrideCategory{StrideInfoDisclosure, StrideTampering}},
		{keywords: []string{"signing", "signature"}, stride: []StrideCategory{StrideSpoofing, StrideTampering}},
		{keywords: []string{"certificate", "cert"}, stride: []StrideCategory{StrideSpoofing, StrideInfoDisclosure}},
		{keywords: []string{"oauth", "saml", "oidc"}, stride: []StrideCategory{StrideSpoofing, StrideElevationPriv}},
		{keywords: []string{"rbac", "abac", "access control"}, stride: []StrideCategory{StrideElevationPriv, StrideInfoDisclosure}},
		{keywords: []string{"monitoring", "alert", "detect"}, stride: []StrideCategory{StrideRepudiation}},
		{keywords: []string{"patch", "update"}, stride: []StrideCategory{StrideTampering, StrideDenialOfService}},
		{keywords: []string{"side channel"}, stride: []StrideCategory{StrideInfoDisclosure}},
		{keywords: []string{"race condition", "timing attack"}, stride: []StrideCategory{StrideElevationPriv, StrideDenialOfService}},
		{keywords: []string{"api gateway", "api key"}, stride: []StrideCategory{StrideSpoofing, StrideDenialOfService}},
		{keywords: []string{"container escape"}, stride: []StrideCategory{StrideElevationPriv}},
		{keywords: []string{"memory safe"}, stride: []StrideCategory{StrideTampering}},
	}
}
