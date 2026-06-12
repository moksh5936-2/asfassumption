package verification

import (
	"regexp"
	"strconv"
	"strings"

	"asf-tui/asf/evidence"
	"asf-tui/asf/models"
)

type Engine struct{}

func NewEngine() *Engine {
	return &Engine{}
}

func (ve *Engine) Verify(assumption models.Assumption, evidenceRecords []models.Evidence) models.Verification {
	var matchedEvidenceIDs []string
	result := models.VerificationResultUNKNOWN
	confidence := 0.0
	var reasoningParts []string
	details := make(map[string]interface{})

	for _, ev := range evidenceRecords {
		if len(ev.Records) == 0 {
			reasoningParts = append(reasoningParts, "No structured records in "+ev.Source)
			continue
		}

		matchedEvidenceIDs = append(matchedEvidenceIDs, ev.ID)
		checkResult := checkAssumptionAgainstEvidence(assumption, ev)

		if checkResult != nil {
			reasoningParts = append(reasoningParts, checkResult.reasoning)
			details[ev.Source] = checkResult.details

			switch checkResult.result {
			case models.VerificationResultVERIFIED:
				if result != models.VerificationResultCONTRADICTED {
					result = models.VerificationResultVERIFIED
				}
				if checkResult.confidence > confidence {
					confidence = checkResult.confidence
				}
			case models.VerificationResultCONTRADICTED:
				result = models.VerificationResultCONTRADICTED
				if checkResult.confidence > confidence {
					confidence = checkResult.confidence
				}
			case models.VerificationResultPARTIALLY_VERIFIED:
				if result != models.VerificationResultCONTRADICTED && result != models.VerificationResultVERIFIED {
					result = models.VerificationResultPARTIALLY_VERIFIED
				}
				if checkResult.confidence > confidence {
					confidence = checkResult.confidence
				}
			}
		}
	}

	if len(matchedEvidenceIDs) == 0 {
		reasoningParts = append(reasoningParts, "No matching evidence available for verification")
	}

	reasoning := strings.Join(reasoningParts, "; ")
	if reasoning == "" {
		reasoning = "No evidence processed"
	}

	return models.NewVerification(assumption.ID, matchedEvidenceIDs, result, confidence, reasoning, details)
}

type checkResult struct {
	result     models.VerificationResult
	confidence float64
	reasoning  string
	details    map[string]interface{}
}

func checkAssumptionAgainstEvidence(assumption models.Assumption, ev models.Evidence) *checkResult {
	switch assumption.AssumptionType {
	case models.AssumptionTypeACCESS:
		return checkAccess(assumption, ev)
	case models.AssumptionTypeIDENTITY:
		return checkIdentity(assumption, ev)
	case models.AssumptionTypeNETWORK:
		return checkNetwork(assumption, ev)
	case models.AssumptionTypeCONFIGURATION:
		return checkConfiguration(assumption, ev)
	case models.AssumptionTypeGOVERNANCE:
		return checkGovernance(assumption, ev)
	case models.AssumptionTypePROCESS:
		return checkGovernance(assumption, ev)
	}
	return nil
}

func extractResourceKeywords(text string) []string {
	re := regexp.MustCompile(`\b[A-Z][a-z]+\b`)
	matches := re.FindAllString(text, -1)
	if len(matches) > 5 {
		matches = matches[:5]
	}
	return matches
}

func checkAccess(assumption models.Assumption, ev models.Evidence) *checkResult {
	textLower := strings.ToLower(assumption.Text)
	records := ev.Records

	onlyRe := regexp.MustCompile(`(?i)only\s+(.+?)\s+(?:can|may|has|should)`)
	onlyMatch := onlyRe.FindStringSubmatch(textLower)
	var restrictedGroup string
	if len(onlyMatch) >= 2 {
		restrictedGroup = strings.TrimSpace(onlyMatch[1])
	}

	var usersOutside, usersInside []string
	resourcesFound := make(map[string]bool)

	for _, rec := range records {
		user := evidence.FindField(rec, []string{"user", "username", "identity", "principal", "name", "email"})
		resource := evidence.FindField(rec, []string{"resource", "application", "system", "target", "service", "scope"})
		permission := evidence.FindField(rec, []string{"permission", "access", "role", "right", "privilege", "action"})
		group := evidence.FindField(rec, []string{"group", "department", "team", "unit", "division", "org"})

		if resource != "" {
			resourcesFound[resource] = true
		}
		if user != "" && permission != "" {
			if restrictedGroup != "" {
				userInGroup := group != "" && (strings.Contains(group, restrictedGroup) || strings.Contains(restrictedGroup, group))
				if userInGroup {
					usersInside = append(usersInside, user)
				} else {
					usersOutside = append(usersOutside, user)
				}
			}
		}
	}

	uniqueOutside := unique(usersOutside)
	uniqueInside := unique(usersInside)

	resourcesList := make([]string, 0, len(resourcesFound))
	for r := range resourcesFound {
		resourcesList = append(resourcesList, r)
	}

	details := map[string]interface{}{
		"expected_group":      restrictedGroup,
		"users_outside_group": uniqueOutside,
		"users_inside_group":  uniqueInside,
		"resources_found":     resourcesList,
		"total_records":       len(records),
	}

	if len(uniqueOutside) > 0 && restrictedGroup != "" {
		limit := 5
		if len(uniqueOutside) < limit {
			limit = len(uniqueOutside)
		}
		return &checkResult{
			result:     models.VerificationResultCONTRADICTED,
			confidence: 0.92,
			reasoning:  "Found " + strconv.Itoa(len(uniqueOutside)) + " user(s) outside '" + restrictedGroup + "' with access: " + strings.Join(uniqueOutside[:limit], ", "),
			details:    details,
		}
	}

	if restrictedGroup != "" && len(uniqueInside) > 0 {
		return &checkResult{
			result:     models.VerificationResultVERIFIED,
			confidence: 0.78,
			reasoning:  "Only users in '" + restrictedGroup + "' found with access (" + strconv.Itoa(len(uniqueInside)) + " users)",
			details:    details,
		}
	}

	return &checkResult{
		result:     models.VerificationResultUNKNOWN,
		confidence: 0.3,
		reasoning:  "Could not determine access patterns from evidence",
		details:    details,
	}
}

func checkIdentity(assumption models.Assumption, ev models.Evidence) *checkResult {
	textLower := strings.ToLower(assumption.Text)
	records := ev.Records
	details := make(map[string]interface{})

	hasMFA := strings.Contains(textLower, "mfa") || strings.Contains(textLower, "multi-factor") || strings.Contains(textLower, "multifactor")

	if hasMFA {
		var mfaUsers, noMFAUsers []string
		for _, rec := range records {
			user := evidence.FindField(rec, []string{"user", "username", "identity", "name", "email"})
			mfaStatus := evidence.FindField(rec, []string{"mfa", "mfa_enabled", "multi_factor", "2fa", "totp"})

			if user != "" {
				if mfaStatus != "" && isTrue(mfaStatus) {
					mfaUsers = append(mfaUsers, user)
				} else {
					noMFAUsers = append(noMFAUsers, user)
				}
			}
		}

		details["mfa_enabled_users"] = unique(mfaUsers)
		details["mfa_disabled_users"] = unique(noMFAUsers)

		if len(noMFAUsers) > 0 && len(mfaUsers) == 0 {
			return &checkResult{
				result:     models.VerificationResultCONTRADICTED,
				confidence: 0.95,
				reasoning:  "MFA not enabled for " + strconv.Itoa(len(noMFAUsers)) + " user(s)",
				details:    details,
			}
		}
		if len(noMFAUsers) > 0 && len(mfaUsers) > 0 {
			return &checkResult{
				result:     models.VerificationResultPARTIALLY_VERIFIED,
				confidence: 0.6,
				reasoning:  "MFA enabled for " + strconv.Itoa(len(mfaUsers)) + " user(s) but missing for " + strconv.Itoa(len(noMFAUsers)),
				details:    details,
			}
		}
		if len(mfaUsers) > 0 && len(noMFAUsers) == 0 {
			return &checkResult{
				result:     models.VerificationResultVERIFIED,
				confidence: 0.85,
				reasoning:  "MFA enabled for all " + strconv.Itoa(len(mfaUsers)) + " user(s)",
				details:    details,
			}
		}
	}

	return &checkResult{
		result:     models.VerificationResultUNKNOWN,
		confidence: 0.3,
		reasoning:  "No identity evidence matched assumption",
		details:    details,
	}
}

func checkNetwork(assumption models.Assumption, ev models.Evidence) *checkResult {
	textLower := strings.ToLower(assumption.Text)
	records := ev.Records
	details := make(map[string]interface{})

	isExposed := strings.Contains(textLower, "internet") || strings.Contains(textLower, "public") || strings.Contains(textLower, "exposed") || strings.Contains(textLower, "external")
	isIsolated := strings.Contains(textLower, "isolat") || strings.Contains(textLower, "segment") || strings.Contains(textLower, "private")

	var exposures, isolations []string
	for _, rec := range records {
		asset := evidence.FindField(rec, []string{"asset", "resource", "host", "server", "system", "name", "service"})
		publicVal := evidence.FindField(rec, []string{"public", "exposed", "internet_facing", "is_public", "exposure"})

		if asset != "" {
			if publicVal != "" && isTrue(publicVal) {
				exposures = append(exposures, asset)
			} else {
				isolations = append(isolations, asset)
			}
		}
	}

	details["exposed_assets"] = unique(exposures)
	details["isolated_assets"] = unique(isolations)

	if isIsolated && len(exposures) > 0 {
		limit := 5
		if len(exposures) < limit {
			limit = len(exposures)
		}
		return &checkResult{
			result:     models.VerificationResultCONTRADICTED,
			confidence: 0.9,
			reasoning:  "Claimed isolated but found " + strconv.Itoa(len(exposures)) + " exposed asset(s): " + strings.Join(unique(exposures)[:limit], ", "),
			details:    details,
		}
	}

	if isIsolated && len(exposures) == 0 {
		return &checkResult{
			result:     models.VerificationResultVERIFIED,
			confidence: 0.8,
			reasoning:  "All " + strconv.Itoa(len(isolations)) + " asset(s) appear isolated",
			details:    details,
		}
	}

	if isExposed && len(exposures) == 0 && len(isolations) > 0 {
		negation := strings.Contains(textLower, "no") || strings.Contains(textLower, "not") || strings.Contains(textLower, "never")
		if negation {
			return &checkResult{
				result:     models.VerificationResultVERIFIED,
				confidence: 0.85,
				reasoning:  "Confirmed: no exposure found across " + strconv.Itoa(len(isolations)) + " asset(s)",
				details:    details,
			}
		}
	}

	if isExposed && len(exposures) > 0 {
		return &checkResult{
			result:     models.VerificationResultVERIFIED,
			confidence: 0.85,
			reasoning:  "Found " + strconv.Itoa(len(exposures)) + " exposed asset(s) as expected",
			details:    details,
		}
	}

	return &checkResult{
		result:     models.VerificationResultUNKNOWN,
		confidence: 0.3,
		reasoning:  "Could not verify network posture",
		details:    details,
	}
}

func checkConfiguration(assumption models.Assumption, ev models.Evidence) *checkResult {
	textLower := strings.ToLower(assumption.Text)
	records := ev.Records
	details := make(map[string]interface{})

	_ = strings.Contains(textLower, "encrypt")
	_ = strings.Contains(textLower, "backup")

	compliant := 0
	nonCompliant := 0
	var examplesCompliant, examplesNonCompliant []string

	for _, rec := range records {
		resource := evidence.FindField(rec, []string{"resource", "system", "asset", "service", "name", "component"})
		enabled := evidence.FindField(rec, []string{"enabled", "status", "state", "active", "value", "configuration"})

		if resource != "" {
			if enabled != "" && isTrue(enabled) {
				compliant++
				if len(examplesCompliant) < 3 {
					examplesCompliant = append(examplesCompliant, resource)
				}
			} else if enabled != "" && isFalse(enabled) {
				nonCompliant++
				if len(examplesNonCompliant) < 3 {
					examplesNonCompliant = append(examplesNonCompliant, resource)
				}
			}
		}
	}

	details["compliant"] = compliant
	details["non_compliant"] = nonCompliant
	details["examples_compliant"] = examplesCompliant
	details["examples_non_compliant"] = examplesNonCompliant

	if nonCompliant > 0 && compliant == 0 {
		return &checkResult{
			result:     models.VerificationResultCONTRADICTED,
			confidence: 0.9,
			reasoning:  "Configuration not applied: " + strconv.Itoa(nonCompliant) + " non-compliant resource(s)",
			details:    details,
		}
	}

	if nonCompliant > 0 && compliant > 0 {
		return &checkResult{
			result:     models.VerificationResultPARTIALLY_VERIFIED,
			confidence: 0.5,
			reasoning:  "Partially compliant: " + strconv.Itoa(compliant) + " OK, " + strconv.Itoa(nonCompliant) + " non-compliant",
			details:    details,
		}
	}

	if compliant > 0 && nonCompliant == 0 {
		return &checkResult{
			result:     models.VerificationResultVERIFIED,
			confidence: 0.85,
			reasoning:  "All " + strconv.Itoa(compliant) + " resource(s) compliant with configuration",
			details:    details,
		}
	}

	return &checkResult{
		result:     models.VerificationResultUNKNOWN,
		confidence: 0.3,
		reasoning:  "Could not verify configuration from evidence",
		details:    details,
	}
}

func checkGovernance(assumption models.Assumption, ev models.Evidence) *checkResult {
	records := ev.Records
	details := make(map[string]interface{})

	reviewed := 0
	notReviewed := 0
	for _, rec := range records {
		status := evidence.FindField(rec, []string{"status", "reviewed", "approved", "state", "completed"})
		if status != "" && isTrue(status) {
			reviewed++
		} else {
			notReviewed++
		}
	}

	details["reviews_completed"] = reviewed
	details["reviews_pending"] = notReviewed

	if notReviewed > 0 && reviewed == 0 {
		return &checkResult{
			result:     models.VerificationResultCONTRADICTED,
			confidence: 0.88,
			reasoning:  "No governance reviews completed (" + strconv.Itoa(notReviewed) + " pending)",
			details:    details,
		}
	}

	if notReviewed > 0 {
		return &checkResult{
			result:     models.VerificationResultPARTIALLY_VERIFIED,
			confidence: 0.55,
			reasoning:  strconv.Itoa(reviewed) + " reviews done, " + strconv.Itoa(notReviewed) + " pending",
			details:    details,
		}
	}

	if reviewed > 0 {
		return &checkResult{
			result:     models.VerificationResultVERIFIED,
			confidence: 0.85,
			reasoning:  "All " + strconv.Itoa(reviewed) + " governance reviews completed",
			details:    details,
		}
	}

	return &checkResult{
		result:     models.VerificationResultUNKNOWN,
		confidence: 0.3,
		reasoning:  "No governance evidence",
		details:    details,
	}
}

func isTrue(val string) bool {
	val = strings.ToLower(val)
	return val == "true" || val == "yes" || val == "enabled" || val == "1" || val == "active" || val == "on"
}

func isFalse(val string) bool {
	val = strings.ToLower(val)
	return val == "false" || val == "no" || val == "disabled" || val == "0" || val == "inactive" || val == "off"
}

func unique(strs []string) []string {
	seen := make(map[string]bool)
	var result []string
	for _, s := range strs {
		if !seen[s] {
			seen[s] = true
			result = append(result, s)
		}
	}
	return result
}
