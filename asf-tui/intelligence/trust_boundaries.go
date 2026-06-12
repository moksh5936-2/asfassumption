package intelligence

import (
	"fmt"
	"strings"
)

// TrustBoundaryEngine discovers trust boundaries from architecture topology.
type TrustBoundaryEngine struct{}

// NewTrustBoundaryEngine creates a trust boundary engine.
func NewTrustBoundaryEngine() *TrustBoundaryEngine {
	return &TrustBoundaryEngine{}
}

// DiscoverBoundaries analyzes the architecture and returns discovered trust boundaries.
func (tbe *TrustBoundaryEngine) DiscoverBoundaries(arch *ArchDescription) []TrustBoundary {
	if arch == nil {
		return nil
	}
	var boundaries []TrustBoundary
	boundaries = append(boundaries, tbe.discoverInternetBoundary(arch)...)
	boundaries = append(boundaries, tbe.discoverIdentityBoundary(arch)...)
	boundaries = append(boundaries, tbe.discoverTenantBoundary(arch)...)
	boundaries = append(boundaries, tbe.discoverVendorBoundary(arch)...)
	boundaries = append(boundaries, tbe.discoverNetworkBoundary(arch)...)
	boundaries = append(boundaries, tbe.discoverAdminBoundary(arch)...)
	boundaries = append(boundaries, tbe.discoverDataBoundary(arch)...)
	boundaries = append(boundaries, tbe.discoverCloudBoundary(arch)...)
	return boundaries
}

// discoverInternetBoundary detects boundaries between internet and internal components.
func (tbe *TrustBoundaryEngine) discoverInternetBoundary(arch *ArchDescription) []TrustBoundary {
	var internetComps []string
	var internalComps []string
	for _, comp := range arch.Components {
		label := strings.ToLower(comp.Label)
		if strings.Contains(label, "internet") || strings.Contains(label, "browser") || strings.Contains(label, "user") || strings.Contains(label, "client") || strings.Contains(label, "external") || strings.Contains(label, "public") {
			internetComps = append(internetComps, comp.Label)
		} else if !strings.Contains(label, "third") && !strings.Contains(label, "vendor") {
			internalComps = append(internalComps, comp.Label)
		}
	}
	if len(internetComps) > 0 && len(internalComps) > 0 {
		return []TrustBoundary{{
			Type:        "Internet",
			Components:  append(internetComps, internalComps...),
			RiskLevel:   RiskCritical,
			Description: fmt.Sprintf("Trust boundary between internet-facing components (%s) and internal components (%s)", strings.Join(internetComps, ", "), strings.Join(internalComps, ", ")),
		}}
	}
	return nil
}

// discoverIdentityBoundary detects boundaries around identity/authentication components.
func (tbe *TrustBoundaryEngine) discoverIdentityBoundary(arch *ArchDescription) []TrustBoundary {
	var idComps []string
	for _, comp := range arch.Components {
		label := strings.ToLower(comp.Label)
		if strings.Contains(label, "auth") || strings.Contains(label, "identity") || strings.Contains(label, "idp") || strings.Contains(label, "sso") || strings.Contains(label, "mfa") || strings.Contains(label, "login") {
			idComps = append(idComps, comp.Label)
		}
	}
	if len(idComps) > 0 {
		return []TrustBoundary{{
			Type:        "Identity",
			Components:  idComps,
			RiskLevel:   RiskCritical,
			Description: fmt.Sprintf("Identity trust boundary around %s", strings.Join(idComps, ", ")),
		}}
	}
	return nil
}

// discoverTenantBoundary detects boundaries between tenants in multi-tenant architectures.
func (tbe *TrustBoundaryEngine) discoverTenantBoundary(arch *ArchDescription) []TrustBoundary {
	var tenantComps []string
	raw := strings.ToLower(arch.RawText)
	for _, comp := range arch.Components {
		label := strings.ToLower(comp.Label)
		if strings.Contains(label, "tenant") || strings.Contains(label, "customer") || strings.Contains(label, "org") || strings.Contains(label, "workspace") {
			tenantComps = append(tenantComps, comp.Label)
		}
	}
	if strings.Contains(raw, "tenant") || strings.Contains(raw, "multi-tenant") || strings.Contains(raw, "multi tenant") {
		if len(tenantComps) == 0 {
			return []TrustBoundary{{
				Type:        "Tenant",
				Components:  []string{"Multi-Tenant System"},
				RiskLevel:   RiskCritical,
				Description: "Tenant trust boundary inferred from multi-tenant architecture references",
			}}
		}
		return []TrustBoundary{{
			Type:        "Tenant",
			Components:  tenantComps,
			RiskLevel:   RiskCritical,
			Description: fmt.Sprintf("Tenant isolation boundary around %s", strings.Join(tenantComps, ", ")),
		}}
	}
	return nil
}

// discoverVendorBoundary detects boundaries around third-party/vendor components.
func (tbe *TrustBoundaryEngine) discoverVendorBoundary(arch *ArchDescription) []TrustBoundary {
	var vendorComps []string
	for _, comp := range arch.Components {
		label := strings.ToLower(comp.Label)
		if strings.Contains(label, "third party") || strings.Contains(label, "third-party") || strings.Contains(label, "thirdparty") || strings.Contains(label, "vendor") || strings.Contains(label, "external") || strings.Contains(label, "saas") || strings.Contains(label, "analytics") || strings.Contains(label, "stripe") || strings.Contains(label, "sendgrid") {
			vendorComps = append(vendorComps, comp.Label)
		}
	}
	if len(vendorComps) > 0 {
		return []TrustBoundary{{
			Type:        "Vendor",
			Components:  vendorComps,
			RiskLevel:   RiskHigh,
			Description: fmt.Sprintf("Vendor trust boundary around %s", strings.Join(vendorComps, ", ")),
		}}
	}
	return nil
}

// discoverNetworkBoundary detects boundaries between network segments.
func (tbe *TrustBoundaryEngine) discoverNetworkBoundary(arch *ArchDescription) []TrustBoundary {
	var dmzComps []string
	var internalComps []string
	for _, comp := range arch.Components {
		label := strings.ToLower(comp.Label)
		if strings.Contains(label, "dmz") || strings.Contains(label, "gateway") || strings.Contains(label, "proxy") || strings.Contains(label, "lb") || strings.Contains(label, "load balancer") || strings.Contains(label, "waf") {
			dmzComps = append(dmzComps, comp.Label)
		} else if strings.Contains(label, "app") || strings.Contains(label, "service") || strings.Contains(label, "backend") || strings.Contains(label, "internal") {
			internalComps = append(internalComps, comp.Label)
		}
	}
	if len(dmzComps) > 0 && len(internalComps) > 0 {
		return []TrustBoundary{{
			Type:        "Network",
			Components:  append(dmzComps, internalComps...),
			RiskLevel:   RiskHigh,
			Description: fmt.Sprintf("Network segmentation boundary between DMZ/edge (%s) and internal (%s)", strings.Join(dmzComps, ", "), strings.Join(internalComps, ", ")),
		}}
	}
	return nil
}

// discoverAdminBoundary detects boundaries around administrative access.
func (tbe *TrustBoundaryEngine) discoverAdminBoundary(arch *ArchDescription) []TrustBoundary {
	var adminComps []string
	for _, comp := range arch.Components {
		label := strings.ToLower(comp.Label)
		if strings.Contains(label, "admin") || strings.Contains(label, "management") || strings.Contains(label, "console") || strings.Contains(label, "dashboard") || strings.Contains(label, "portal") || strings.Contains(label, "ops") {
			adminComps = append(adminComps, comp.Label)
		}
	}
	if len(adminComps) > 0 {
		return []TrustBoundary{{
			Type:        "Admin",
			Components:  adminComps,
			RiskLevel:   RiskCritical,
			Description: fmt.Sprintf("Administrative trust boundary around %s", strings.Join(adminComps, ", ")),
		}}
	}
	return nil
}

// discoverDataBoundary detects boundaries around sensitive data stores.
func (tbe *TrustBoundaryEngine) discoverDataBoundary(arch *ArchDescription) []TrustBoundary {
	var dataComps []string
	for _, comp := range arch.Components {
		label := strings.ToLower(comp.Label)
		if strings.Contains(label, "database") || strings.Contains(label, "db") || strings.Contains(label, "storage") || strings.Contains(label, "cache") || strings.Contains(label, "blob") || strings.Contains(label, "s3") || strings.Contains(label, "data lake") || strings.Contains(label, "warehouse") {
			dataComps = append(dataComps, comp.Label)
		}
	}
	if len(dataComps) > 0 {
		return []TrustBoundary{{
			Type:        "Data",
			Components:  dataComps,
			RiskLevel:   RiskCritical,
			Description: fmt.Sprintf("Data trust boundary around %s", strings.Join(dataComps, ", ")),
		}}
	}
	return nil
}

// discoverCloudBoundary detects boundaries between cloud provider and on-premise.
func (tbe *TrustBoundaryEngine) discoverCloudBoundary(arch *ArchDescription) []TrustBoundary {
	var cloudComps []string
	var onPremComps []string
	for _, comp := range arch.Components {
		label := strings.ToLower(comp.Label)
		if strings.Contains(label, "aws") || strings.Contains(label, "azure") || strings.Contains(label, "gcp") || strings.Contains(label, "cloud") || strings.Contains(label, "lambda") || strings.Contains(label, "function") || strings.Contains(label, "serverless") {
			cloudComps = append(cloudComps, comp.Label)
		} else if strings.Contains(label, "on-prem") || strings.Contains(label, "onprem") || strings.Contains(label, "datacenter") || strings.Contains(label, "dc") || strings.Contains(label, "legacy") {
			onPremComps = append(onPremComps, comp.Label)
		}
	}
	if len(cloudComps) > 0 && len(onPremComps) > 0 {
		return []TrustBoundary{{
			Type:        "Cloud",
			Components:  append(cloudComps, onPremComps...),
			RiskLevel:   RiskHigh,
			Description: fmt.Sprintf("Cloud trust boundary between cloud (%s) and on-premise (%s)", strings.Join(cloudComps, ", "), strings.Join(onPremComps, ", ")),
		}}
	}
	return nil
}

// GenerateAssumptions creates assumptions for each discovered boundary.
func (tbe *TrustBoundaryEngine) GenerateAssumptions(boundaries []TrustBoundary) []Assumption {
	var assumptions []Assumption
	for _, b := range boundaries {
		switch b.Type {
		case "Internet":
			assumptions = append(assumptions, Assumption{
				ID:          fmt.Sprintf("TB-INT-%03d", len(assumptions)+1),
				Description: fmt.Sprintf("Internet-facing trust boundary requires TLS termination, WAF, and DDoS protection at %s", strings.Join(b.Components[:minInt(3, len(b.Components))], ", ")),
				Component:   strings.Join(b.Components, ", "),
				Category:    "TrustBoundaries",
				Risk:        RiskCritical,
				Likelihood:  5,
				Impact:      5,
				Confidence:  0.90,
				Keywords:    []string{"internet", "trust boundary", "tls", "waf", "ddos"},
				Rationale:   "Internet-facing components are exposed to global threat actors; edge protection is mandatory.",
			})
		case "Identity":
			assumptions = append(assumptions, Assumption{
				ID:          fmt.Sprintf("TB-IDT-%03d", len(assumptions)+1),
				Description: fmt.Sprintf("Identity trust boundary requires MFA, session hardening, and token validation at %s", strings.Join(b.Components[:minInt(3, len(b.Components))], ", ")),
				Component:   strings.Join(b.Components, ", "),
				Category:    "TrustBoundaries",
				Risk:        RiskCritical,
				Likelihood:  4,
				Impact:      5,
				Confidence:  0.90,
				Keywords:    []string{"identity", "trust boundary", "mfa", "session", "token"},
				Rationale:   "Identity boundary breaches lead to account takeover and lateral movement; strong controls are required.",
			})
		case "Tenant":
			assumptions = append(assumptions, Assumption{
				ID:          fmt.Sprintf("TB-TEN-%03d", len(assumptions)+1),
				Description: fmt.Sprintf("Tenant trust boundary requires isolation, object-level authorization, and data segregation at %s", strings.Join(b.Components[:minInt(3, len(b.Components))], ", ")),
				Component:   strings.Join(b.Components, ", "),
				Category:    "TrustBoundaries",
				Risk:        RiskCritical,
				Likelihood:  4,
				Impact:      5,
				Confidence:  0.90,
				Keywords:    []string{"tenant", "trust boundary", "isolation", "bola", "segregation"},
				Rationale:   "Tenant isolation failures lead to cross-tenant data leakage and regulatory action.",
			})
		case "Vendor":
			assumptions = append(assumptions, Assumption{
				ID:          fmt.Sprintf("TB-VND-%03d", len(assumptions)+1),
				Description: fmt.Sprintf("Vendor trust boundary requires vendor risk assessment, data minimization, and egress monitoring at %s", strings.Join(b.Components[:minInt(3, len(b.Components))], ", ")),
				Component:   strings.Join(b.Components, ", "),
				Category:    "TrustBoundaries",
				Risk:        RiskHigh,
				Likelihood:  4,
				Impact:      4,
				Confidence:  0.85,
				Keywords:    []string{"vendor", "trust boundary", "risk assessment", "egress", "third-party"},
				Rationale:   "Vendor integrations introduce supply chain risk; boundary controls and assessments are required.",
			})
		case "Network":
			assumptions = append(assumptions, Assumption{
				ID:          fmt.Sprintf("TB-NET-%03d", len(assumptions)+1),
				Description: fmt.Sprintf("Network trust boundary requires segmentation, ACL enforcement, and lateral movement prevention at %s", strings.Join(b.Components[:minInt(3, len(b.Components))], ", ")),
				Component:   strings.Join(b.Components, ", "),
				Category:    "TrustBoundaries",
				Risk:        RiskHigh,
				Likelihood:  4,
				Impact:      4,
				Confidence:  0.85,
				Keywords:    []string{"network", "trust boundary", "segmentation", "acl", "lateral movement"},
				Rationale:   "Network segmentation gaps allow lateral movement after initial compromise.",
			})
		case "Admin":
			assumptions = append(assumptions, Assumption{
				ID:          fmt.Sprintf("TB-ADM-%03d", len(assumptions)+1),
				Description: fmt.Sprintf("Administrative trust boundary requires MFA, break-glass, and command logging at %s", strings.Join(b.Components[:minInt(3, len(b.Components))], ", ")),
				Component:   strings.Join(b.Components, ", "),
				Category:    "TrustBoundaries",
				Risk:        RiskCritical,
				Likelihood:  4,
				Impact:      5,
				Confidence:  0.90,
				Keywords:    []string{"admin", "trust boundary", "mfa", "break-glass", "command logging"},
				Rationale:   "Admin boundaries are high-value targets; compromise leads to full system control.",
			})
		case "Data":
			assumptions = append(assumptions, Assumption{
				ID:          fmt.Sprintf("TB-DAT-%03d", len(assumptions)+1),
				Description: fmt.Sprintf("Data trust boundary requires encryption, access controls, and audit logging at %s", strings.Join(b.Components[:minInt(3, len(b.Components))], ", ")),
				Component:   strings.Join(b.Components, ", "),
				Category:    "TrustBoundaries",
				Risk:        RiskCritical,
				Likelihood:  4,
				Impact:      5,
				Confidence:  0.90,
				Keywords:    []string{"data", "trust boundary", "encryption", "access control", "audit"},
				Rationale:   "Data boundaries protect sensitive information; encryption and access controls are mandatory.",
			})
		case "Cloud":
			assumptions = append(assumptions, Assumption{
				ID:          fmt.Sprintf("TB-CLD-%03d", len(assumptions)+1),
				Description: fmt.Sprintf("Cloud trust boundary requires IAM alignment, encrypted transit, and data residency controls at %s", strings.Join(b.Components[:minInt(3, len(b.Components))], ", ")),
				Component:   strings.Join(b.Components, ", "),
				Category:    "TrustBoundaries",
				Risk:        RiskHigh,
				Likelihood:  3,
				Impact:      4,
				Confidence:  0.85,
				Keywords:    []string{"cloud", "trust boundary", "iam", "data residency", "transit encryption"},
				Rationale:   "Cloud boundaries require IAM alignment and data residency to prevent cross-environment leakage.",
			})
		}
	}
	return assumptions
}

// SummarizeBoundaries returns a summary string of all boundaries.
func SummarizeBoundaries(boundaries []TrustBoundary) string {
	if len(boundaries) == 0 {
		return "No trust boundaries discovered"
	}
	var parts []string
	for _, b := range boundaries {
		parts = append(parts, fmt.Sprintf("%s boundary (%s) with %d component(s)", b.Type, b.RiskLevel, len(b.Components)))
	}
	return strings.Join(parts, "; ")
}
