package main

import (
	"fmt"
	"strings"
)

type AIEnhancer struct {
	model     *ModelManager
	modelName string
}

func NewAIEnhancer() *AIEnhancer {
	return &AIEnhancer{model: NewModelManager()}
}

type AIEnhancedResult struct {
	AdditionalAssumptions []AIAssumption
	RefinedRisks          []AIRiskRefinement
	MissingThreats        []string
	Recommendations       []string
	RawResponse           string
}

type AIAssumption struct {
	Description string
	Category    string
	Risk        string
	Reasoning   string
}

type AIRiskRefinement struct {
	AssumptionID  string
	OriginalRisk  RiskLevel
	SuggestedRisk RiskLevel
	Reasoning     string
}

func (ae *AIEnhancer) Enhance(result *AnalysisResult, modelName string) (*AIEnhancedResult, error) {
	ae.modelName = modelName

	if !ae.model.CheckAvailable() {
		return nil, fmt.Errorf("Ollama binary not found")
	}

	if !ae.model.CheckRunning() {
		return nil, fmt.Errorf("Ollama is not running. Start it with: ollama serve")
	}

	if !ae.model.IsModelInstalled(modelName) {
		return nil, fmt.Errorf("Model %q is not installed. Download it in AI Settings or choose another installed model", modelName)
	}

	prompt := ae.buildPrompt(result)
	response, err := ae.model.Generate(prompt, modelName)
	if err != nil {
		return nil, fmt.Errorf("AI generation failed: %w", err)
	}

	return ae.parseResponse(response, result), nil
}

func (ae *AIEnhancer) buildPrompt(result *AnalysisResult) string {
	var b strings.Builder
	b.WriteString("You are a senior security architect reviewing an ASF (Architecture Security Framework) analysis.\n\n")
	b.WriteString(fmt.Sprintf("Architecture: %s\n", result.ArchitectureName))
	b.WriteString(fmt.Sprintf("Analysis Mode: %s\n", result.AnalysisMode))
	b.WriteString(fmt.Sprintf("Total Assumptions Found: %d\n\n", result.TotalAssumptions))

	b.WriteString("Existing Assumptions:\n")
	for i, a := range result.Assumptions {
		b.WriteString(fmt.Sprintf("%d. [%s][%s] %s\n", i+1, a.Risk, a.Category, a.Description))
	}
	b.WriteString("\n")

	b.WriteString("Existing Controls:\n")
	for _, c := range result.Controls {
		b.WriteString(fmt.Sprintf("- %s: %s\n", c.ID, c.Description))
	}
	b.WriteString("\n")

	b.WriteString("STRIDE Distribution:\n")
	for cat, count := range result.StrideDistribution {
		b.WriteString(fmt.Sprintf("- %s: %d\n", cat, count))
	}
	b.WriteString("\n")

	b.WriteString(`Review this analysis and provide:

1. ADDITIONAL_SECURITY_ASSUMPTIONS: Security assumptions that ASF may have missed. List each with category and risk level.

2. RISK_REFINEMENTS: For any assumptions whose risk level seems incorrect, suggest a correction. Format each as:
   Assumption <ID> current=<risk> suggested=<risk> reason=<brief explanation>

3. MISSING_THREAT_SCENARIOS: Attack scenarios not covered by the existing assumptions.

4. RECOMMENDATIONS: Additional security controls or architecture improvements.

Be specific. Do not repeat assumptions already listed. Focus on what ASF might have missed.`)

	return b.String()
}

func (ae *AIEnhancer) parseResponse(response string, result *AnalysisResult) *AIEnhancedResult {
	enh := &AIEnhancedResult{
		RawResponse: response,
	}

	sections := strings.Split(response, "\n\n")
	currentSection := ""
	for _, section := range sections {
		lower := strings.ToLower(section)

		switch {
		case strings.Contains(lower, "additional_security_assumptions") || strings.Contains(lower, "additional security assumption"):
			currentSection = "assumptions"
		case strings.Contains(lower, "risk_refinement") || strings.Contains(lower, "risk refinement"):
			currentSection = "risks"
		case strings.Contains(lower, "missing_threat") || strings.Contains(lower, "missing threat"):
			currentSection = "threats"
		case strings.Contains(lower, "recommendation"):
			currentSection = "recommendations"
		default:
			switch currentSection {
			case "assumptions":
				if lines := parseBulletList(section); len(lines) > 0 {
					for _, line := range lines {
						enh.AdditionalAssumptions = append(enh.AdditionalAssumptions, AIAssumption{
							Description: line,
							Category:    inferCategory(line),
							Risk:        inferRisk(line),
						})
					}
				}
			case "risks":
				enh.RefinedRisks = append(enh.RefinedRisks, parseRiskRefinements(section)...)
			case "threats":
				enh.MissingThreats = append(enh.MissingThreats, parseBulletList(section)...)
			case "recommendations":
				enh.Recommendations = append(enh.Recommendations, parseBulletList(section)...)
			}
		}
	}

	return enh
}

func parseBulletList(text string) []string {
	var items []string
	for _, line := range strings.Split(text, "\n") {
		line = strings.TrimSpace(line)
		line = strings.TrimLeft(line, "-*•0123456789. ")
		line = strings.TrimSpace(line)
		if len(line) > 10 {
			items = append(items, line)
		}
	}
	return items
}

func parseRiskRefinements(text string) []AIRiskRefinement {
	var refinements []AIRiskRefinement
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		// Match: Assumption <ID> current=<risk> suggested=<risk> reason=<...>
		var id, current, suggested, reason string
		parts := strings.Fields(line)
		for i, p := range parts {
			switch {
			case p == "current=" && i+1 < len(parts):
				current = parts[i+1]
			case strings.HasPrefix(p, "current="):
				current = strings.TrimPrefix(p, "current=")
			case p == "suggested=" && i+1 < len(parts):
				suggested = parts[i+1]
			case strings.HasPrefix(p, "suggested="):
				suggested = strings.TrimPrefix(p, "suggested=")
			case p == "reason=" && i+1 < len(parts):
				// Collect rest of line as reason
				rest := strings.Join(parts[i+1:], " ")
				reason = strings.TrimRight(rest, ".")
			case strings.HasPrefix(p, "reason="):
				reason = strings.TrimPrefix(p, "reason=")
				reason = strings.TrimRight(reason, ".")
			case i == 0 && p == "Assumption" && i+1 < len(parts):
				id = parts[i+1]
			case i > 0 && p != "Assumption" && id == "" && !strings.Contains(p, "="):
				id = p
			}
		}
		if id != "" && current != "" && suggested != "" {
			refinements = append(refinements, AIRiskRefinement{
				AssumptionID:  id,
				OriginalRisk:  RiskLevel(current),
				SuggestedRisk: RiskLevel(suggested),
				Reasoning:     reason,
			})
		}
	}
	return refinements
}

func inferCategory(text string) string {
	lower := strings.ToLower(text)
	switch {
	case strings.Contains(lower, "authenticate") || strings.Contains(lower, "login") || strings.Contains(lower, "credential") || strings.Contains(lower, "session"):
		return "AUTHENTICATION"
	case strings.Contains(lower, "authorize") || strings.Contains(lower, "permission") || strings.Contains(lower, "access control") || strings.Contains(lower, "role"):
		return "AUTHORIZATION"
	case strings.Contains(lower, "encrypt") || strings.Contains(lower, "tls") || strings.Contains(lower, "ssl"):
		return "ENCRYPTION"
	case strings.Contains(lower, "network") || strings.Contains(lower, "firewall") || strings.Contains(lower, "segment"):
		return "NETWORK"
	case strings.Contains(lower, "log") || strings.Contains(lower, "audit"):
		return "LOGGING"
	case strings.Contains(lower, "backup") || strings.Contains(lower, "recover"):
		return "BACKUP"
	case strings.Contains(lower, "database") || strings.Contains(lower, "sql"):
		return "DATABASE"
	case strings.Contains(lower, "third") || strings.Contains(lower, "vendor") || strings.Contains(lower, "supply"):
		return "THIRD_PARTY"
	case strings.Contains(lower, "identity") || strings.Contains(lower, "mfa") || strings.Contains(lower, "sso"):
		return "IDENTITY"
	default:
		return "GENERAL"
	}
}

func inferRisk(text string) string {
	lower := strings.ToLower(text)
	switch {
	case strings.Contains(lower, "critical") || strings.Contains(lower, "severe"):
		return "Critical"
	case strings.Contains(lower, "high") || strings.Contains(lower, "significant"):
		return "High"
	case strings.Contains(lower, "low") || strings.Contains(lower, "minor"):
		return "Low"
	default:
		return "Medium"
	}
}

func mergeAIResults(original *AnalysisResult, ai *AIEnhancedResult) *AnalysisResult {
	if ai == nil {
		return original
	}

	for _, aa := range ai.AdditionalAssumptions {
		risk := RiskLevel(aa.Risk)
		stride := NewStrideEngine().MapAssumption(aa.Category, aa.Description, nil)
		original.Assumptions = append(original.Assumptions, Assumption{
			ID:          fmt.Sprintf("AI-%d", len(original.Assumptions)+1),
			Description: aa.Description,
			Category:    aa.Category,
			Risk:        risk,
			Stride:      stride,
			Confidence:  0.5,
		})
		original.TotalAssumptions++
		switch risk {
		case RiskCritical:
			original.CriticalCount++
		case RiskHigh:
			original.HighCount++
		case RiskMedium:
			original.MediumCount++
		case RiskLow:
			original.LowCount++
		}
	}

	for _, rr := range ai.RefinedRisks {
		for i := range original.Assumptions {
			if original.Assumptions[i].ID == rr.AssumptionID && original.Assumptions[i].Risk != rr.SuggestedRisk {
				original.Assumptions[i].Risk = rr.SuggestedRisk
				original.Assumptions[i].Rationale = fmt.Sprintf("Risk refined by AI: was %s, now %s. Reason: %s",
					rr.OriginalRisk, rr.SuggestedRisk, rr.Reasoning)
			}
		}
	}

	// Recompute risk counts after refinements
	original.CriticalCount = 0
	original.HighCount = 0
	original.MediumCount = 0
	original.LowCount = 0
	for _, a := range original.Assumptions {
		switch a.Risk {
		case RiskCritical:
			original.CriticalCount++
		case RiskHigh:
			original.HighCount++
		case RiskMedium:
			original.MediumCount++
		case RiskLow:
			original.LowCount++
		}
	}

	for i, r := range ai.Recommendations {
		original.Controls = append(original.Controls, ControlDetail{
			ID:          fmt.Sprintf("AI-CTRL-%03d", i+1),
			Description: r,
			Rationale:   "AI-generated control recommendation",
			Priority:    3,
		})
	}

	return original
}
