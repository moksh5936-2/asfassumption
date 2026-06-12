package intelligence

import (
	"fmt"
	"strings"
)

// ─────────────────────────────────────────────────────────────
// PHASE 1 — TRUST ZONE MODEL
// ─────────────────────────────────────────────────────────────

// TrustZoneType represents the type of a trust zone.
type TrustZoneType string

const (
	ZoneINTERNET       TrustZoneType = "INTERNET"
	ZoneTHIRD_PARTY    TrustZoneType = "THIRD_PARTY"
	ZoneIDENTITY       TrustZoneType = "IDENTITY"
	ZoneCLIENT         TrustZoneType = "CLIENT"
	ZoneAPPLICATION    TrustZoneType = "APPLICATION"
	ZoneMANAGEMENT     TrustZoneType = "MANAGEMENT"
	ZoneDATA           TrustZoneType = "DATA"
	ZoneDATABASE       TrustZoneType = "DATABASE"
	ZoneBACKUP         TrustZoneType = "BACKUP"
	ZoneLOGGING        TrustZoneType = "LOGGING"
	ZoneMONITORING     TrustZoneType = "MONITORING"
	ZoneADMINISTRATIVE TrustZoneType = "ADMINISTRATIVE"
	ZonePRODUCTION     TrustZoneType = "PRODUCTION"
	ZoneSTAGING        TrustZoneType = "STAGING"
	ZoneDEVELOPMENT    TrustZoneType = "DEVELOPMENT"
	ZoneCI_CD          TrustZoneType = "CI_CD"
	ZoneSECRETS        TrustZoneType = "SECRETS"
	ZoneDMZ            TrustZoneType = "DMZ"
	ZoneVPN            TrustZoneType = "VPN"
	ZoneJUMP_HOST      TrustZoneType = "JUMP_HOST"
	ZoneNETWORK        TrustZoneType = "NETWORK"
	ZoneUNKNOWN        TrustZoneType = "UNKNOWN"
)

// TrustZone represents a security trust zone in the architecture.
type TrustZone struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Type        TrustZoneType `json:"type"`
	Sensitivity string        `json:"sensitivity"` // public, internal, confidential, restricted, critical
	Components  []string      `json:"components"`
	Description string        `json:"description"`
}

// ─────────────────────────────────────────────────────────────
// PHASE 2 — TRUST BOUNDARY MODEL
// ─────────────────────────────────────────────────────────────

// CrossingType represents the type of trust boundary crossing.
type CrossingType string

const (
	CrossingPUBLIC_TO_INTERNAL         CrossingType = "PUBLIC_TO_INTERNAL"
	CrossingTHIRD_PARTY_TO_INTERNAL    CrossingType = "THIRD_PARTY_TO_INTERNAL"
	CrossingIDENTITY_TO_APPLICATION    CrossingType = "IDENTITY_TO_APPLICATION"
	CrossingAPPLICATION_TO_DATA        CrossingType = "APPLICATION_TO_DATA"
	CrossingADMIN_TO_PRODUCTION        CrossingType = "ADMIN_TO_PRODUCTION"
	CrossingDEVELOPMENT_TO_PRODUCTION  CrossingType = "DEVELOPMENT_TO_PRODUCTION"
	CrossingBACKUP_ACCESS              CrossingType = "BACKUP_ACCESS"
	CrossingSECRETS_ACCESS             CrossingType = "SECRETS_ACCESS"
	CrossingMONITORING_ACCESS          CrossingType = "MONITORING_ACCESS"
	CrossingCLIENT_TO_APPLICATION      CrossingType = "CLIENT_TO_APPLICATION"
	CrossingAPPLICATION_TO_THIRD_PARTY CrossingType = "APPLICATION_TO_THIRD_PARTY"
	CrossingUNKNOWN                    CrossingType = "UNKNOWN"
)

// TBITrustBoundary represents a rich trust boundary from the Trust Boundary Intelligence Engine.
type TBITrustBoundary struct {
	ID                  string        `json:"id"`
	SourceZone          string        `json:"source_zone"`
	DestinationZone     string        `json:"destination_zone"`
	SourceZoneType      TrustZoneType `json:"source_zone_type"`
	DestinationZoneType TrustZoneType `json:"destination_zone_type"`
	CrossingType        CrossingType  `json:"crossing_type"`
	Risk                RiskLevel     `json:"risk"`
	Confidence          float64       `json:"confidence"`
	RequiredControls    []string      `json:"required_controls"`
	RequiredAssumptions []string      `json:"required_assumptions"`
	Threats             []string      `json:"threats"`
	MissingControls     []string      `json:"missing_controls,omitempty"`
	MissingAssumptions  []string      `json:"missing_assumptions,omitempty"`
	Reasoning           string        `json:"reasoning"`
	Recommendations     []string      `json:"recommendations,omitempty"`
	ComplianceMappings  []string      `json:"compliance_mappings,omitempty"`
}

// BoundaryWeakness represents a detected weakness at a trust boundary.
type BoundaryWeakness struct {
	ID              string    `json:"id"`
	BoundaryID      string    `json:"boundary_id"`
	Type            string    `json:"type"`
	Severity        RiskLevel `json:"severity"`
	Description     string    `json:"description"`
	Reasoning       string    `json:"reasoning"`
	Recommendations []string  `json:"recommendations,omitempty"`
}

// ─────────────────────────────────────────────────────────────
// TRUST BOUNDARY INTELLIGENCE ENGINE
// ─────────────────────────────────────────────────────────────

// TBIEngine (Trust Boundary Intelligence Engine) discovers trust zones, boundaries, and generates boundary-specific findings.
type TBIEngine struct {
	zonePatterns      map[TrustZoneType][]string
	controlLibrary    map[CrossingType][]string
	assumptionLibrary map[CrossingType][]string
	threatLibrary     map[CrossingType][]string
}

// NewTBIEngine creates a new trust boundary intelligence engine.
func NewTBIEngine() *TBIEngine {
	return &TBIEngine{
		zonePatterns: map[TrustZoneType][]string{
			ZoneINTERNET:       {"internet", "browser", "mobile", "client", "user", "public", "external", "web", "frontend", "cdn", "edge", "dmz"},
			ZoneTHIRD_PARTY:    {"third", "vendor", "partner", "supplier", "external", "outsourced", "saas", "api", "webhook", "callback", "integration"},
			ZoneIDENTITY:       {"auth", "identity", "idp", "sso", "oauth", "oidc", "saml", "login", "mfa", "directory", "ldap", "ad", "cognito", "auth0", "okta", "keycloak"},
			ZoneCLIENT:         {"client", "browser", "app", "mobile", "desktop", "user", "workstation", "device", "endpoint"},
			ZoneAPPLICATION:    {"app", "api", "gateway", "service", "microservice", "server", "backend", "webapp", "frontend", "lambda", "function", "container", "pod"},
			ZoneMANAGEMENT:     {"management", "orchestration", "control", "scheduler", "controller", "master", "manager", "admin"},
			ZoneDATA:           {"data", "storage", "cache", "object", "blob", "file", "datastore", "warehouse", "lake", "analytics"},
			ZoneDATABASE:       {"database", "db", "postgres", "mysql", "rds", "mongodb", "dynamodb", "redis", "sql", "nosql", "cassandra", "neo4j", "sqlite", "oracle", "mssql"},
			ZoneBACKUP:         {"backup", "restore", "archive", "snapshot", "replication", "dr", "disaster", "recovery"},
			ZoneLOGGING:        {"log", "audit", "trace", "event", "syslog", "elk", "splunk", "graylog", "fluentd", "journal"},
			ZoneMONITORING:     {"monitor", "alert", "metric", "observability", "prometheus", "grafana", "datadog", "newrelic", "cloudwatch", "nagios", "zabbix"},
			ZoneADMINISTRATIVE: {"admin", "console", "dashboard", "portal", "management", "operator", "root", "privileged", "break.glass", "firecall"},
			ZonePRODUCTION:     {"production", "prod", "live", "active", "primary", "main", "real"},
			ZoneSTAGING:        {"staging", "stage", "preprod", "uat", "qa", "test", "preview"},
			ZoneDEVELOPMENT:    {"development", "dev", "local", "sandbox", "experimental", "debug"},
			ZoneCI_CD:          {"ci", "cd", "pipeline", "jenkins", "gitlab", "github", "travis", "circleci", "build", "deploy", "artifact", "registry"},
			ZoneSECRETS:        {"vault", "secret", "key", "credential", "password", "token", "certificate", "cert", "kms", "hsm", "vault", "keeper", "lastpass"},
			ZoneDMZ:            {"dmz", "demilitarized", "perimeter", "edge", "public", "external"},
			ZoneVPN:            {"vpn", "virtual private", "tunnel", "remote access", "site-to-site"},
			ZoneJUMP_HOST:      {"jump", "bastion", "jump host", "jump server", "bastion host", "proxy", "gateway"},
			ZoneNETWORK:        {"network", "firewall", "waf", "ids", "ips", "load balancer", "lb", "ingress", "egress", "router", "switch", "gateway", "proxy", "cdn", "dns"},
		},
		controlLibrary: map[CrossingType][]string{
			CrossingPUBLIC_TO_INTERNAL:         {"TLS", "WAF", "DDoS Protection", "Rate Limiting", "Authentication", "MFA", "IP Restriction", "Bot Detection", "Request Validation"},
			CrossingTHIRD_PARTY_TO_INTERNAL:    {"TLS", "API Authentication", "API Key", "OAuth", "Rate Limiting", "Input Validation", "Webhook Verification", "Contract Testing"},
			CrossingIDENTITY_TO_APPLICATION:    {"Token Validation", "Session Management", "Token Rotation", "JWT Verification", "Replay Protection", "Token Binding", "Clock Skew Tolerance"},
			CrossingAPPLICATION_TO_DATA:        {"Authorization", "Query Parameterization", "Encryption at Rest", "Encryption in Transit", "Audit Logging", "Data Classification", "Least Privilege"},
			CrossingADMIN_TO_PRODUCTION:        {"MFA", "PAM", "Just-In-Time Access", "Session Recording", "Change Control", "Approval Workflow", "Audit Logging", "Time-Based Access"},
			CrossingDEVELOPMENT_TO_PRODUCTION:  {"Code Review", "CI/CD Approval", "Artifact Signing", "Deployment Gates", "Rollback Plan", "Immutable Infrastructure", "Secret Rotation"},
			CrossingBACKUP_ACCESS:              {"Backup Encryption", "Access Logging", "Restore Testing", "Immutable Backups", "Geographic Distribution", "Retention Policy"},
			CrossingSECRETS_ACCESS:             {"Secret Encryption", "Access Audit", "Rotation Policy", "Least Privilege", "HSM", "Dual Control", "Emergency Access"},
			CrossingMONITORING_ACCESS:          {"Log Integrity", "Tamper Detection", "Access Logging", "Alert Escalation", "SIEM Integration", "Log Retention"},
			CrossingCLIENT_TO_APPLICATION:      {"TLS", "Certificate Pinning", "Input Validation", "CSRF Protection", "CORS", "Content Security Policy", "Session Management"},
			CrossingAPPLICATION_TO_THIRD_PARTY: {"TLS", "Mutual TLS", "API Authentication", "Timeout", "Circuit Breaker", "Retry Policy", "Fallback"},
		},
		assumptionLibrary: map[CrossingType][]string{
			CrossingPUBLIC_TO_INTERNAL:         {"All public traffic is encrypted", "Authentication is enforced for all access", "Rate limiting prevents abuse", "WAF protects against injection attacks", "DDoS protection is active"},
			CrossingTHIRD_PARTY_TO_INTERNAL:    {"Third-party access is authenticated", "API keys are rotated regularly", "Webhook signatures are verified", "Third-party data is validated", "Contracts enforce security requirements"},
			CrossingIDENTITY_TO_APPLICATION:    {"Identity tokens are validated correctly", "Token signing keys are rotated", "Session expiration is enforced", "Token replay is prevented", "Identity provider availability is guaranteed"},
			CrossingAPPLICATION_TO_DATA:        {"Database credentials are protected", "Access is restricted to authorized queries", "All data access is logged", "Data is encrypted at rest", "Object-level authorization is enforced"},
			CrossingADMIN_TO_PRODUCTION:        {"Administrative access requires MFA", "Admin sessions are recorded", "Privileged access is reviewed", "Admin actions are audited", "Break-glass procedures exist"},
			CrossingDEVELOPMENT_TO_PRODUCTION:  {"Code is reviewed before deployment", "Deployments require approval", "Production secrets are not in code", "Rollbacks are tested", "Development cannot access production data"},
			CrossingBACKUP_ACCESS:              {"Backup data is encrypted", "Backup integrity is verified", "Restore procedures are tested", "Backup access is logged", "Backups are geographically distributed"},
			CrossingSECRETS_ACCESS:             {"Secrets are encrypted at rest", "Secret access is audited", "Secrets are rotated regularly", "Secret retrieval is authenticated", "Emergency access is controlled"},
			CrossingMONITORING_ACCESS:          {"Logs are tamper-evident", "Log access is audited", "Monitoring covers all critical paths", "Alerts are reviewed and acted upon", "Log retention meets compliance"},
			CrossingCLIENT_TO_APPLICATION:      {"Client communication is encrypted", "Client input is validated", "Sessions are securely managed", "CSRF protection is active", "Client certificates are validated"},
			CrossingAPPLICATION_TO_THIRD_PARTY: {"Outbound communication is encrypted", "API timeouts are configured", "Circuit breakers prevent cascade failure", "Third-party errors are handled", "Fallback mechanisms exist"},
		},
		threatLibrary: map[CrossingType][]string{
			CrossingPUBLIC_TO_INTERNAL:         {"Man-in-the-Middle", "DDoS", "Injection", "Authentication Bypass", "Credential Stuffing", "Bot Attacks"},
			CrossingTHIRD_PARTY_TO_INTERNAL:    {"Supply Chain Attack", "API Key Compromise", "Webhook Spoofing", "Data Injection", "Third-Party Breach"},
			CrossingIDENTITY_TO_APPLICATION:    {"Token Theft", "Session Hijacking", "Replay Attack", "Token Forgery", "Identity Provider Compromise"},
			CrossingAPPLICATION_TO_DATA:        {"SQL Injection", "Data Exfiltration", "Privilege Escalation", "Unauthorized Access", "Data Tampering"},
			CrossingADMIN_TO_PRODUCTION:        {"Privilege Abuse", "Admin Account Compromise", "Insider Threat", "Unauthorized Change", "Session Hijacking"},
			CrossingDEVELOPMENT_TO_PRODUCTION:  {"Deployment of Malicious Code", "Secret Leakage", "Production Data Exposure", "Unauthorized Change", "Rollback Failure"},
			CrossingBACKUP_ACCESS:              {"Backup Theft", "Ransomware", "Backup Tampering", "Unauthorized Restore", "Data Loss"},
			CrossingSECRETS_ACCESS:             {"Secret Extraction", "Secret Rotation Failure", "Unauthorized Access", "Key Compromise", "Insider Threat"},
			CrossingMONITORING_ACCESS:          {"Log Tampering", "Alert Suppression", "Monitoring Bypass", "Blind Spot Exploitation", "Audit Evasion"},
			CrossingCLIENT_TO_APPLICATION:      {"Session Hijacking", "CSRF", "XSS", "Client Compromise", "MITM"},
			CrossingAPPLICATION_TO_THIRD_PARTY: {"Data Leakage", "Third-Party Compromise", "Dependency Vulnerability", "API Abuse", "Timeout Exploitation"},
		},
	}
}

// ─────────────────────────────────────────────────────────────
// PHASE 3 — AUTOMATIC ZONE DISCOVERY
// ─────────────────────────────────────────────────────────────

// DiscoverZones analyzes the architecture and returns discovered trust zones.
func (tbi *TBIEngine) DiscoverZones(arch *ArchDescription) []TrustZone {
	if arch == nil || len(arch.Components) == 0 {
		return nil
	}

	var zones []TrustZone
	zoneMap := make(map[TrustZoneType][]string)

	// Classify each component into a zone
	for _, comp := range arch.Components {
		zoneType := tbi.classifyComponent(comp.Label)
		zoneMap[zoneType] = append(zoneMap[zoneType], comp.Label)
	}

	// Create zones from classified components
	zoneID := 1
	for zoneType, components := range zoneMap {
		if len(components) == 0 {
			continue
		}
		sensitivity := tbi.inferSensitivity(zoneType, arch)
		zones = append(zones, TrustZone{
			ID:          fmt.Sprintf("ZONE-%03d", zoneID),
			Name:        string(zoneType),
			Type:        zoneType,
			Sensitivity: sensitivity,
			Components:  components,
			Description: fmt.Sprintf("Trust zone %s containing %d component(s): %s", zoneType, len(components), strings.Join(components, ", ")),
		})
		zoneID++
	}

	return zones
}

// classifyComponent determines the trust zone type for a component label.
func (tbi *TBIEngine) classifyComponent(label string) TrustZoneType {
	lower := strings.ToLower(label)
	scores := make(map[TrustZoneType]int)

	for zoneType, patterns := range tbi.zonePatterns {
		for _, pattern := range patterns {
			if strings.Contains(lower, pattern) {
				scores[zoneType]++
			}
		}
	}

	// Priority list for deterministic tie-breaking
	priority := []TrustZoneType{
		ZoneINTERNET, ZoneIDENTITY, ZoneCLIENT, ZoneAPPLICATION,
		ZoneDATA, ZoneDATABASE, ZoneSECRETS, ZoneADMINISTRATIVE,
		ZoneTHIRD_PARTY, ZoneLOGGING, ZoneMONITORING, ZoneBACKUP,
		ZoneMANAGEMENT, ZonePRODUCTION, ZoneSTAGING, ZoneDEVELOPMENT,
		ZoneCI_CD, ZoneDMZ, ZoneVPN, ZoneJUMP_HOST, ZoneNETWORK,
	}

	bestType := ZoneUNKNOWN
	bestScore := 0
	for _, zoneType := range priority {
		score, ok := scores[zoneType]
		if ok && score > bestScore {
			bestScore = score
			bestType = zoneType
		}
	}

	return bestType
}

// inferSensitivity determines the sensitivity level of a zone based on its type and architecture context.
func (tbi *TBIEngine) inferSensitivity(zoneType TrustZoneType, arch *ArchDescription) string {
	// Check for PHI/PCI/sensitive data in architecture
	rawText := strings.ToLower(arch.RawText + " " + arch.Name)
	hasPHI := strings.Contains(rawText, "phi") || strings.Contains(rawText, "hipaa") || strings.Contains(rawText, "patient") || strings.Contains(rawText, "medical") || strings.Contains(rawText, "health")
	hasPCI := strings.Contains(rawText, "pci") || strings.Contains(rawText, "payment") || strings.Contains(rawText, "card") || strings.Contains(rawText, "cardholder")
	hasSensitive := strings.Contains(rawText, "sensitive") || strings.Contains(rawText, "confidential") || strings.Contains(rawText, "secret") || strings.Contains(rawText, "private")

	switch zoneType {
	case ZoneINTERNET:
		return "public"
	case ZoneDATABASE, ZoneDATA:
		if hasPHI || hasPCI {
			return "critical"
		}
		if hasSensitive {
			return "restricted"
		}
		return "confidential"
	case ZoneSECRETS:
		return "critical"
	case ZoneADMINISTRATIVE:
		return "restricted"
	case ZonePRODUCTION:
		return "restricted"
	case ZoneTHIRD_PARTY:
		return "external"
	case ZoneIDENTITY:
		return "restricted"
	case ZoneCLIENT:
		return "public"
	default:
		return "internal"
	}
}

// ─────────────────────────────────────────────────────────────
// PHASE 4 — BOUNDARY DISCOVERY
// ─────────────────────────────────────────────────────────────

// DiscoverBoundaries discovers trust boundaries between zones.
func (tbi *TBIEngine) DiscoverBoundaries(zones []TrustZone, arch *ArchDescription) []TBITrustBoundary {
	if len(zones) == 0 || arch == nil {
		return nil
	}

	var boundaries []TBITrustBoundary
	boundaryID := 1

	// Find relationships that cross zone boundaries
	for _, rel := range arch.Relationships {
		sourceZone := tbi.findZoneForComponent(zones, rel.Source)
		targetZone := tbi.findZoneForComponent(zones, rel.Target)

		if sourceZone == nil || targetZone == nil {
			continue
		}

		// Skip if same zone
		if sourceZone.ID == targetZone.ID {
			continue
		}

		crossingType := tbi.inferCrossingType(sourceZone.Type, targetZone.Type)
		boundary := TBITrustBoundary{
			ID:                  fmt.Sprintf("BOUND-%03d", boundaryID),
			SourceZone:          sourceZone.Name,
			DestinationZone:     targetZone.Name,
			SourceZoneType:      sourceZone.Type,
			DestinationZoneType: targetZone.Type,
			CrossingType:        crossingType,
			Risk:                tbi.scoreBoundaryRisk(sourceZone, targetZone, crossingType, arch),
			Confidence:          0.85,
			RequiredControls:    tbi.controlLibrary[crossingType],
			RequiredAssumptions: tbi.assumptionLibrary[crossingType],
			Threats:             tbi.threatLibrary[crossingType],
			Reasoning:           fmt.Sprintf("Boundary crossing from %s (%s) to %s (%s) via %s", sourceZone.Name, sourceZone.Type, targetZone.Name, targetZone.Type, rel.Label),
		}
		boundaries = append(boundaries, boundary)
		boundaryID++
	}

	return boundaries
}

// findZoneForComponent finds the zone that contains a component.
func (tbi *TBIEngine) findZoneForComponent(zones []TrustZone, componentLabel string) *TrustZone {
	for i := range zones {
		for _, comp := range zones[i].Components {
			if comp == componentLabel {
				return &zones[i]
			}
		}
	}
	return nil
}

// inferCrossingType determines the crossing type based on source and destination zone types.
func (tbi *TBIEngine) inferCrossingType(source, dest TrustZoneType) CrossingType {
	// Public/internet to internal
	if source == ZoneINTERNET || source == ZoneCLIENT {
		if dest == ZoneAPPLICATION || dest == ZoneIDENTITY {
			return CrossingPUBLIC_TO_INTERNAL
		}
		if dest == ZoneTHIRD_PARTY {
			return CrossingAPPLICATION_TO_THIRD_PARTY
		}
	}

	// Third party to internal
	if source == ZoneTHIRD_PARTY {
		if dest == ZoneAPPLICATION || dest == ZoneDATABASE || dest == ZoneDATA {
			return CrossingTHIRD_PARTY_TO_INTERNAL
		}
	}

	// Identity to application
	if source == ZoneIDENTITY {
		if dest == ZoneAPPLICATION || dest == ZoneADMINISTRATIVE {
			return CrossingIDENTITY_TO_APPLICATION
		}
	}

	// Application to data
	if source == ZoneAPPLICATION || source == ZoneMANAGEMENT {
		if dest == ZoneDATABASE || dest == ZoneDATA || dest == ZoneBACKUP {
			return CrossingAPPLICATION_TO_DATA
		}
		if dest == ZoneSECRETS {
			return CrossingSECRETS_ACCESS
		}
		if dest == ZoneMONITORING || dest == ZoneLOGGING {
			return CrossingMONITORING_ACCESS
		}
		if dest == ZoneTHIRD_PARTY {
			return CrossingAPPLICATION_TO_THIRD_PARTY
		}
	}

	// Admin to production
	if source == ZoneADMINISTRATIVE || source == ZoneMANAGEMENT {
		if dest == ZonePRODUCTION || dest == ZoneAPPLICATION || dest == ZoneDATABASE {
			return CrossingADMIN_TO_PRODUCTION
		}
	}

	// Development to production
	if source == ZoneDEVELOPMENT || source == ZoneSTAGING || source == ZoneCI_CD {
		if dest == ZonePRODUCTION || dest == ZoneAPPLICATION || dest == ZoneDATABASE {
			return CrossingDEVELOPMENT_TO_PRODUCTION
		}
	}

	// Backup access
	if source == ZoneBACKUP || dest == ZoneBACKUP {
		return CrossingBACKUP_ACCESS
	}

	// Client to application
	if source == ZoneCLIENT {
		if dest == ZoneAPPLICATION || dest == ZoneIDENTITY {
			return CrossingCLIENT_TO_APPLICATION
		}
	}

	return CrossingUNKNOWN
}

// ─────────────────────────────────────────────────────────────
// PHASE 5 & 6 — REQUIRED CONTROL AND ASSUMPTION LIBRARY
// (Already defined in engine initialization)
// ─────────────────────────────────────────────────────────────

// ─────────────────────────────────────────────────────────────
// PHASE 7 — BOUNDARY RISK ENGINE
// ─────────────────────────────────────────────────────────────

// scoreBoundaryRisk scores the risk of a boundary crossing.
func (tbi *TBIEngine) scoreBoundaryRisk(source, target *TrustZone, crossing CrossingType, arch *ArchDescription) RiskLevel {
	score := 0

	// Base risk by crossing type
	switch crossing {
	case CrossingPUBLIC_TO_INTERNAL:
		score += 4
	case CrossingTHIRD_PARTY_TO_INTERNAL:
		score += 3
	case CrossingIDENTITY_TO_APPLICATION:
		score += 3
	case CrossingAPPLICATION_TO_DATA:
		score += 3
	case CrossingADMIN_TO_PRODUCTION:
		score += 4
	case CrossingDEVELOPMENT_TO_PRODUCTION:
		score += 3
	case CrossingSECRETS_ACCESS:
		score += 4
	case CrossingBACKUP_ACCESS:
		score += 2
	case CrossingMONITORING_ACCESS:
		score += 2
	case CrossingCLIENT_TO_APPLICATION:
		score += 2
	case CrossingAPPLICATION_TO_THIRD_PARTY:
		score += 2
	default:
		score += 1
	}

	// Sensitivity boost
	if source.Sensitivity == "critical" || target.Sensitivity == "critical" {
		score += 2
	}
	if source.Sensitivity == "restricted" || target.Sensitivity == "restricted" {
		score += 1
	}

	// PHI/PCI context
	rawText := strings.ToLower(arch.RawText + " " + arch.Name)
	if strings.Contains(rawText, "phi") || strings.Contains(rawText, "hipaa") || strings.Contains(rawText, "patient") || strings.Contains(rawText, "medical") {
		score += 2
	}
	if strings.Contains(rawText, "pci") || strings.Contains(rawText, "payment") || strings.Contains(rawText, "card") || strings.Contains(rawText, "cardholder") {
		score += 2
	}

	// Identity systems
	if source.Type == ZoneIDENTITY || target.Type == ZoneIDENTITY {
		score += 1
	}

	// Map score to risk level
	switch {
	case score >= 6:
		return RiskCritical
	case score >= 4:
		return RiskHigh
	case score >= 2:
		return RiskMedium
	default:
		return RiskLow
	}
}

// ─────────────────────────────────────────────────────────────
// PHASE 8 — BOUNDARY WEAKNESS DETECTION
// ─────────────────────────────────────────────────────────────

// DetectWeaknesses detects missing controls and assumptions at boundaries.
func (tbi *TBIEngine) DetectWeaknesses(boundaries []TBITrustBoundary, arch *ArchDescription, existingAssumptions []Assumption) []BoundaryWeakness {
	var weaknesses []BoundaryWeakness

	// Build map of existing assumptions for quick lookup
	existingMap := make(map[string]bool)
	for _, a := range existingAssumptions {
		lower := strings.ToLower(a.Description)
		existingMap[lower] = true
		// Also index by keywords
		for _, kw := range a.Keywords {
			existingMap[strings.ToLower(kw)] = true
		}
	}

	// Also check raw text and explicit assumptions
	for _, text := range arch.ExplicitAssumptions {
		existingMap[strings.ToLower(text)] = true
	}

	weaknessID := 1
	for _, boundary := range boundaries {
		// Check missing controls
		for _, control := range boundary.RequiredControls {
			if !tbi.isControlPresent(control, existingMap, arch) {
				weaknesses = append(weaknesses, BoundaryWeakness{
					ID:              fmt.Sprintf("WEAK-%03d", weaknessID),
					BoundaryID:      boundary.ID,
					Type:            "Missing Control",
					Severity:        boundary.Risk,
					Description:     fmt.Sprintf("Boundary %s is missing control: %s", boundary.ID, control),
					Reasoning:       fmt.Sprintf("Crossing type %s requires %s, but it is not documented in the architecture", boundary.CrossingType, control),
					Recommendations: []string{fmt.Sprintf("Document and implement %s at the %s boundary", control, boundary.CrossingType)},
				})
				weaknessID++
			}
		}

		// Check missing assumptions
		for _, assumption := range boundary.RequiredAssumptions {
			if !tbi.isAssumptionPresent(assumption, existingMap) {
				weaknesses = append(weaknesses, BoundaryWeakness{
					ID:              fmt.Sprintf("WEAK-%03d", weaknessID),
					BoundaryID:      boundary.ID,
					Type:            "Missing Assumption",
					Severity:        boundary.Risk,
					Description:     fmt.Sprintf("Boundary %s is missing assumption: %s", boundary.ID, assumption),
					Reasoning:       fmt.Sprintf("Crossing type %s requires the assumption that %s, but it is not documented", boundary.CrossingType, assumption),
					Recommendations: []string{fmt.Sprintf("Document and verify that %s", assumption)},
				})
				weaknessID++
			}
		}
	}

	return weaknesses
}

// isControlPresent checks if a control is present in the architecture.
func (tbi *TBIEngine) isControlPresent(control string, existingMap map[string]bool, arch *ArchDescription) bool {
	lower := strings.ToLower(control)

	// Check existing assumptions
	if existingMap[lower] {
		return true
	}

	// Check security controls
	for category, controls := range arch.SecurityControls {
		for _, c := range controls {
			if strings.Contains(strings.ToLower(c), lower) || strings.Contains(strings.ToLower(category), lower) {
				return true
			}
		}
	}

	// Check raw text
	if strings.Contains(strings.ToLower(arch.RawText), lower) {
		return true
	}

	// Check policies
	for _, policy := range arch.Policies {
		if strings.Contains(strings.ToLower(policy), lower) {
			return true
		}
	}

	return false
}

// isAssumptionPresent checks if an assumption is present in the architecture.
func (tbi *TBIEngine) isAssumptionPresent(assumption string, existingMap map[string]bool) bool {
	lower := strings.ToLower(assumption)

	// Check exact match
	if existingMap[lower] {
		return true
	}

	// Check keyword overlap
	keywords := strings.Fields(lower)
	for _, kw := range keywords {
		if len(kw) > 3 && existingMap[kw] {
			return true
		}
	}

	return false
}

// ─────────────────────────────────────────────────────────────
// PHASE 9 — TRUST ASSUMPTION GENERATION
// ─────────────────────────────────────────────────────────────

// GenerateAssumptions generates assumptions for each boundary.
func (tbi *TBIEngine) GenerateAssumptions(boundaries []TBITrustBoundary) []Assumption {
	var assumptions []Assumption
	assumptionID := 1

	for _, boundary := range boundaries {
		for _, reqAssumption := range boundary.RequiredAssumptions {
			assumptions = append(assumptions, Assumption{
				ID:          fmt.Sprintf("TBI-%03d", assumptionID),
				Description: fmt.Sprintf("Boundary %s (%s → %s): %s", boundary.ID, boundary.SourceZone, boundary.DestinationZone, reqAssumption),
				Component:   fmt.Sprintf("%s → %s", boundary.SourceZone, boundary.DestinationZone),
				Category:    string(boundary.CrossingType),
				Risk:        boundary.Risk,
				Confidence:  boundary.Confidence,
				Keywords:    []string{string(boundary.CrossingType), "trust boundary", "boundary"},
				SourceType:  "tbi-generated",
				SourceFile:  "trust_boundary_intelligence",
			})
			assumptionID++
		}
	}

	return assumptions
}

// ─────────────────────────────────────────────────────────────
// PHASE 10 — COMPLIANCE ENRICHMENT
// ─────────────────────────────────────────────────────────────

// EnrichCompliance maps boundaries to compliance requirements.
func (tbi *TBIEngine) EnrichCompliance(boundaries []TBITrustBoundary, arch *ArchDescription) []TBITrustBoundary {
	complianceMap := map[string][]string{
		"hipaa":    {"Access Controls", "Audit Logging", "Encryption", "Integrity", "Transmission Security"},
		"soc2":     {"Access Management", "Monitoring", "Change Management", "Encryption", "Backup"},
		"iso27001": {"Access Control", "Risk Management", "Supplier Management", "Asset Management", "Incident Management"},
		"pci":      {"Encryption", "Access Control", "Monitoring", "Testing", "Network Security"},
		"gdpr":     {"Consent", "Access", "Deletion", "Encryption", "Breach Notification"},
		"nist":     {"Access Control", "Encryption", "Monitoring", "Incident Response", "Backup"},
	}

	for i := range boundaries {
		var mappings []string
		for _, comp := range arch.Compliance {
			compLower := strings.ToLower(comp)
			// Find matching compliance framework by checking if compLower contains map key
			var matchedKey string
			var requirements []string
			for key, reqs := range complianceMap {
				if strings.Contains(compLower, key) {
					matchedKey = key
					requirements = reqs
					break
				}
			}
			if matchedKey == "" {
				continue
			}

			// Check if this boundary type is relevant to this compliance
			if tbi.isBoundaryRelevantForCompliance(boundaries[i].CrossingType, matchedKey) {
				for _, req := range requirements {
					mappings = append(mappings, fmt.Sprintf("%s: %s", comp, req))
				}
			}
		}
		boundaries[i].ComplianceMappings = mappings
	}

	return boundaries
}

// isBoundaryRelevantForCompliance checks if a boundary type is relevant to a compliance framework.
func (tbi *TBIEngine) isBoundaryRelevantForCompliance(crossing CrossingType, compliance string) bool {
	relevance := map[string][]CrossingType{
		"hipaa":    {CrossingPUBLIC_TO_INTERNAL, CrossingIDENTITY_TO_APPLICATION, CrossingAPPLICATION_TO_DATA, CrossingADMIN_TO_PRODUCTION},
		"soc2":     {CrossingPUBLIC_TO_INTERNAL, CrossingIDENTITY_TO_APPLICATION, CrossingAPPLICATION_TO_DATA, CrossingMONITORING_ACCESS, CrossingADMIN_TO_PRODUCTION},
		"iso27001": {CrossingPUBLIC_TO_INTERNAL, CrossingTHIRD_PARTY_TO_INTERNAL, CrossingIDENTITY_TO_APPLICATION, CrossingAPPLICATION_TO_DATA, CrossingADMIN_TO_PRODUCTION},
		"pci":      {CrossingPUBLIC_TO_INTERNAL, CrossingIDENTITY_TO_APPLICATION, CrossingAPPLICATION_TO_DATA, CrossingSECRETS_ACCESS},
		"gdpr":     {CrossingPUBLIC_TO_INTERNAL, CrossingIDENTITY_TO_APPLICATION, CrossingAPPLICATION_TO_DATA, CrossingAPPLICATION_TO_THIRD_PARTY},
		"nist":     {CrossingPUBLIC_TO_INTERNAL, CrossingIDENTITY_TO_APPLICATION, CrossingAPPLICATION_TO_DATA, CrossingADMIN_TO_PRODUCTION, CrossingMONITORING_ACCESS},
	}

	for _, relevant := range relevance[compliance] {
		if relevant == crossing {
			return true
		}
	}
	return false
}

// ─────────────────────────────────────────────────────────────
// PHASE 11 — SUMMARY GENERATION
// ─────────────────────────────────────────────────────────────

// BuildSummary creates a human-readable summary of trust boundary analysis.
func (tbi *TBIEngine) BuildSummary(zones []TrustZone, boundaries []TBITrustBoundary, weaknesses []BoundaryWeakness) string {
	if len(boundaries) == 0 {
		return "No trust boundaries detected."
	}

	// Count by risk
	riskCounts := make(map[RiskLevel]int)
	for _, b := range boundaries {
		riskCounts[b.Risk]++
	}

	// Count weaknesses
	weaknessCounts := make(map[string]int)
	for _, w := range weaknesses {
		weaknessCounts[w.Type]++
	}

	return fmt.Sprintf(
		"Trust Boundary Analysis: %d zones, %d boundaries (Critical: %d, High: %d, Medium: %d, Low: %d), %d weaknesses detected (%d missing controls, %d missing assumptions)",
		len(zones),
		len(boundaries),
		riskCounts[RiskCritical],
		riskCounts[RiskHigh],
		riskCounts[RiskMedium],
		riskCounts[RiskLow],
		len(weaknesses),
		weaknessCounts["Missing Control"],
		weaknessCounts["Missing Assumption"],
	)
}

// ─────────────────────────────────────────────────────────────
// MAIN ORCHESTRATOR
// ─────────────────────────────────────────────────────────────

// Run executes the full trust boundary intelligence pipeline.
func (tbi *TBIEngine) Run(arch *ArchDescription, existingAssumptions []Assumption) (*TBIRunResult, error) {
	if arch == nil {
		return nil, fmt.Errorf("architecture is nil")
	}

	// Phase 3: Discover zones
	zones := tbi.DiscoverZones(arch)

	// Phase 4: Discover boundaries
	boundaries := tbi.DiscoverBoundaries(zones, arch)

	// Phase 10: Enrich compliance
	boundaries = tbi.EnrichCompliance(boundaries, arch)

	// Phase 8: Detect weaknesses
	weaknesses := tbi.DetectWeaknesses(boundaries, arch, existingAssumptions)

	// Phase 9: Generate assumptions
	assumptions := tbi.GenerateAssumptions(boundaries)

	// Build summary
	summary := tbi.BuildSummary(zones, boundaries, weaknesses)

	return &TBIRunResult{
		Zones:       zones,
		Boundaries:  boundaries,
		Weaknesses:  weaknesses,
		Assumptions: assumptions,
		Summary:     summary,
	}, nil
}

// TBIRunResult holds the complete output of the TBI engine.
type TBIRunResult struct {
	Zones       []TrustZone        `json:"zones"`
	Boundaries  []TBITrustBoundary `json:"boundaries"`
	Weaknesses  []BoundaryWeakness `json:"weaknesses"`
	Assumptions []Assumption       `json:"assumptions"`
	Summary     string             `json:"summary"`
}
