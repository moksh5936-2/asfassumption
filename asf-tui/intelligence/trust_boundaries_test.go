package intelligence

import (
	"strings"
	"testing"
)

func TestTrustBoundaryEngineDiscoverInternet(t *testing.T) {
	tbe := NewTrustBoundaryEngine()
	arch := &ArchDescription{
		Name: "web",
		Components: []Component{
			{ID: "u1", Label: "User"},
			{ID: "gw1", Label: "Gateway"},
			{ID: "app1", Label: "App"},
		},
	}
	boundaries := tbe.DiscoverBoundaries(arch)
	found := false
	for _, b := range boundaries {
		if b.Type == "Internet" {
			found = true
			if b.RiskLevel != RiskCritical {
				t.Errorf("expected Critical risk for internet boundary, got %s", b.RiskLevel)
			}
		}
	}
	if !found {
		t.Error("expected Internet boundary")
	}
}

func TestTrustBoundaryEngineDiscoverIdentity(t *testing.T) {
	tbe := NewTrustBoundaryEngine()
	arch := &ArchDescription{
		Name: "auth",
		Components: []Component{
			{ID: "idp1", Label: "IdentityProvider"},
			{ID: "sso1", Label: "SSO"},
		},
	}
	boundaries := tbe.DiscoverBoundaries(arch)
	found := false
	for _, b := range boundaries {
		if b.Type == "Identity" {
			found = true
			if b.RiskLevel != RiskCritical {
				t.Errorf("expected Critical risk for identity boundary, got %s", b.RiskLevel)
			}
		}
	}
	if !found {
		t.Error("expected Identity boundary")
	}
}

func TestTrustBoundaryEngineDiscoverTenant(t *testing.T) {
	tbe := NewTrustBoundaryEngine()
	arch := &ArchDescription{
		Name: "saas",
		Components: []Component{
			{ID: "t1", Label: "TenantA"},
			{ID: "t2", Label: "TenantB"},
		},
		RawText: "multi-tenant architecture",
	}
	boundaries := tbe.DiscoverBoundaries(arch)
	found := false
	for _, b := range boundaries {
		if b.Type == "Tenant" {
			found = true
			if b.RiskLevel != RiskCritical {
				t.Errorf("expected Critical risk for tenant boundary, got %s", b.RiskLevel)
			}
		}
	}
	if !found {
		t.Error("expected Tenant boundary")
	}
}

func TestTrustBoundaryEngineDiscoverVendor(t *testing.T) {
	tbe := NewTrustBoundaryEngine()
	arch := &ArchDescription{
		Name: "vendor",
		Components: []Component{
			{ID: "v1", Label: "ThirdPartyService"},
			{ID: "app1", Label: "App"},
		},
	}
	boundaries := tbe.DiscoverBoundaries(arch)
	found := false
	for _, b := range boundaries {
		if b.Type == "Vendor" {
			found = true
			if b.RiskLevel != RiskHigh {
				t.Errorf("expected High risk for vendor boundary, got %s", b.RiskLevel)
			}
		}
	}
	if !found {
		t.Error("expected Vendor boundary")
	}
}

func TestTrustBoundaryEngineDiscoverNetwork(t *testing.T) {
	tbe := NewTrustBoundaryEngine()
	arch := &ArchDescription{
		Name: "net",
		Components: []Component{
			{ID: "dmz1", Label: "DMZ"},
			{ID: "app1", Label: "App"},
			{ID: "db1", Label: "Backend"},
		},
	}
	boundaries := tbe.DiscoverBoundaries(arch)
	found := false
	for _, b := range boundaries {
		if b.Type == "Network" {
			found = true
			if b.RiskLevel != RiskHigh {
				t.Errorf("expected High risk for network boundary, got %s", b.RiskLevel)
			}
		}
	}
	if !found {
		t.Error("expected Network boundary")
	}
}

func TestTrustBoundaryEngineDiscoverAdmin(t *testing.T) {
	tbe := NewTrustBoundaryEngine()
	arch := &ArchDescription{
		Name: "admin",
		Components: []Component{
			{ID: "adm1", Label: "AdminConsole"},
			{ID: "db1", Label: "Database"},
		},
	}
	boundaries := tbe.DiscoverBoundaries(arch)
	found := false
	for _, b := range boundaries {
		if b.Type == "Admin" {
			found = true
			if b.RiskLevel != RiskCritical {
				t.Errorf("expected Critical risk for admin boundary, got %s", b.RiskLevel)
			}
		}
	}
	if !found {
		t.Error("expected Admin boundary")
	}
}

func TestTrustBoundaryEngineDiscoverData(t *testing.T) {
	tbe := NewTrustBoundaryEngine()
	arch := &ArchDescription{
		Name: "data",
		Components: []Component{
			{ID: "db1", Label: "Database"},
			{ID: "s3", Label: "S3"},
		},
	}
	boundaries := tbe.DiscoverBoundaries(arch)
	found := false
	for _, b := range boundaries {
		if b.Type == "Data" {
			found = true
			if b.RiskLevel != RiskCritical {
				t.Errorf("expected Critical risk for data boundary, got %s", b.RiskLevel)
			}
		}
	}
	if !found {
		t.Error("expected Data boundary")
	}
}

func TestTrustBoundaryEngineDiscoverCloud(t *testing.T) {
	tbe := NewTrustBoundaryEngine()
	arch := &ArchDescription{
		Name: "hybrid",
		Components: []Component{
			{ID: "aws1", Label: "AWS"},
			{ID: "dc1", Label: "DataCenter"},
		},
	}
	boundaries := tbe.DiscoverBoundaries(arch)
	found := false
	for _, b := range boundaries {
		if b.Type == "Cloud" {
			found = true
			if b.RiskLevel != RiskHigh {
				t.Errorf("expected High risk for cloud boundary, got %s", b.RiskLevel)
			}
		}
	}
	if !found {
		t.Error("expected Cloud boundary")
	}
}

func TestTrustBoundaryEngineGenerateAssumptions(t *testing.T) {
	tbe := NewTrustBoundaryEngine()
	boundaries := []TrustBoundary{
		{Type: "Internet", Components: []string{"User", "Gateway"}, RiskLevel: RiskCritical},
		{Type: "Data", Components: []string{"Database"}, RiskLevel: RiskCritical},
	}
	assumptions := tbe.GenerateAssumptions(boundaries)
	if len(assumptions) != 2 {
		t.Errorf("expected 2 assumptions, got %d", len(assumptions))
	}
	for _, a := range assumptions {
		if a.Category != "TrustBoundaries" {
			t.Errorf("expected TrustBoundaries category, got %s", a.Category)
		}
		if a.Rationale == "" {
			t.Error("expected rationale for boundary assumption")
		}
	}
}

func TestTrustBoundaryEngineNilArch(t *testing.T) {
	tbe := NewTrustBoundaryEngine()
	boundaries := tbe.DiscoverBoundaries(nil)
	if len(boundaries) != 0 {
		t.Errorf("expected 0 boundaries for nil arch, got %d", len(boundaries))
	}
}

func TestSummarizeBoundaries(t *testing.T) {
	boundaries := []TrustBoundary{
		{Type: "Internet", Components: []string{"User"}, RiskLevel: RiskCritical},
		{Type: "Data", Components: []string{"DB"}, RiskLevel: RiskCritical},
	}
	summary := SummarizeBoundaries(boundaries)
	if summary == "No trust boundaries discovered" {
		t.Error("expected non-empty summary")
	}
	if !strings.Contains(summary, "Internet") {
		t.Error("expected Internet in summary")
	}
	if !strings.Contains(summary, "Data") {
		t.Error("expected Data in summary")
	}
}

func TestSummarizeBoundariesEmpty(t *testing.T) {
	summary := SummarizeBoundaries(nil)
	if summary != "No trust boundaries discovered" {
		t.Errorf("expected empty summary, got %s", summary)
	}
}
