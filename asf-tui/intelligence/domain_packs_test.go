package intelligence

import (
	"strings"
	"testing"
)

func TestDomainEngineDetectHealthcare(t *testing.T) {
	de := NewDomainEngine()
	arch := &ArchDescription{
		Name: "health",
		Components: []Component{
			{ID: "phi1", Label: "PHIDatabase"},
			{ID: "app1", Label: "PatientPortal"},
		},
		RawText: "HIPAA compliant system storing PHI.",
	}
	domain := de.DetectDomain(arch)
	if domain != "Healthcare" {
		t.Errorf("expected Healthcare domain, got %s", domain)
	}
}

func TestDomainEngineDetectFintech(t *testing.T) {
	de := NewDomainEngine()
	arch := &ArchDescription{
		Name: "fintech",
		Components: []Component{
			{ID: "pay1", Label: "PaymentGateway"},
			{ID: "fraud1", Label: "FraudDetection"},
		},
		RawText: "PCI DSS compliant payment processing.",
	}
	domain := de.DetectDomain(arch)
	if domain != "Fintech" {
		t.Errorf("expected Fintech domain, got %s", domain)
	}
}

func TestDomainEngineDetectSaaS(t *testing.T) {
	de := NewDomainEngine()
	arch := &ArchDescription{
		Name: "saas",
		Components: []Component{
			{ID: "t1", Label: "TenantManager"},
			{ID: "api1", Label: "API"},
		},
		RawText: "Multi-tenant SaaS platform with tenant isolation.",
	}
	domain := de.DetectDomain(arch)
	if domain != "SaaS" {
		t.Errorf("expected SaaS domain, got %s", domain)
	}
}

func TestDomainEngineDetectKubernetes(t *testing.T) {
	de := NewDomainEngine()
	arch := &ArchDescription{
		Name: "k8s",
		Components: []Component{
			{ID: "pod1", Label: "Pod"},
			{ID: "svc1", Label: "Service"},
		},
		RawText: "Kubernetes cluster with network policies.",
	}
	domain := de.DetectDomain(arch)
	if domain != "Kubernetes" {
		t.Errorf("expected Kubernetes domain, got %s", domain)
	}
}

func TestDomainEngineDetectCloudNative(t *testing.T) {
	de := NewDomainEngine()
	arch := &ArchDescription{
		Name: "cloud",
		Components: []Component{
			{ID: "fn1", Label: "Lambda"},
			{ID: "db1", Label: "DynamoDB"},
		},
		RawText: "Cloud-native serverless architecture.",
	}
	domain := de.DetectDomain(arch)
	if domain != "CloudNative" {
		t.Errorf("expected CloudNative domain, got %s", domain)
	}
}

func TestDomainEngineDetectVPN(t *testing.T) {
	de := NewDomainEngine()
	arch := &ArchDescription{
		Name: "vpn",
		Components: []Component{
			{ID: "vpn1", Label: "VPNGateway"},
			{ID: "client1", Label: "RemoteClient"},
		},
		RawText: "IPSec VPN tunnel for remote access.",
	}
	domain := de.DetectDomain(arch)
	if domain != "VPN" {
		t.Errorf("expected VPN domain, got %s", domain)
	}
}

func TestDomainEngineDetectIdentityPlatform(t *testing.T) {
	de := NewDomainEngine()
	arch := &ArchDescription{
		Name: "idp",
		Components: []Component{
			{ID: "sso1", Label: "SSO"},
			{ID: "fed1", Label: "Federation"},
		},
		RawText: "Identity platform with SSO and SAML federation.",
	}
	domain := de.DetectDomain(arch)
	if domain != "IdentityPlatform" {
		t.Errorf("expected IdentityPlatform domain, got %s", domain)
	}
}

func TestDomainEngineDetectDataPlatform(t *testing.T) {
	de := NewDomainEngine()
	arch := &ArchDescription{
		Name: "data",
		Components: []Component{
			{ID: "etl1", Label: "ETLPipeline"},
			{ID: "dw1", Label: "DataWarehouse"},
		},
		RawText: "Data platform with lineage and governance.",
	}
	domain := de.DetectDomain(arch)
	if domain != "DataPlatform" {
		t.Errorf("expected DataPlatform domain, got %s", domain)
	}
}

func TestDomainEngineDetectUnknown(t *testing.T) {
	de := NewDomainEngine()
	arch := &ArchDescription{
		Name: "generic",
		Components: []Component{
			{ID: "app1", Label: "App"},
		},
		RawText: "Generic application.",
	}
	domain := de.DetectDomain(arch)
	if domain != "" {
		t.Errorf("expected empty domain for generic app, got %s", domain)
	}
}

func TestDomainEngineApplyDomainPack(t *testing.T) {
	de := NewDomainEngine()
	arch := &ArchDescription{
		Name: "health",
		Components: []Component{
			{ID: "phi1", Label: "PHI"},
		},
	}
	assumptions := de.ApplyDomainPack("Healthcare", arch)
	if len(assumptions) == 0 {
		t.Fatal("expected healthcare domain assumptions")
	}
	foundAudit := false
	foundPrivacy := false
	for _, a := range assumptions {
		if a.Category == "Auditability" && strings.Contains(a.Description, "HIPAA") {
			foundAudit = true
		}
		if a.Category == "Privacy" && strings.Contains(a.Description, "patient privacy") {
			foundPrivacy = true
		}
	}
	if !foundAudit {
		t.Error("expected HIPAA audit assumption")
	}
	if !foundPrivacy {
		t.Error("expected patient privacy assumption")
	}
}

func TestDomainEngineApplyDomainPackNil(t *testing.T) {
	de := NewDomainEngine()
	assumptions := de.ApplyDomainPack("NonExistent", nil)
	if len(assumptions) != 0 {
		t.Errorf("expected 0 assumptions for nil arch, got %d", len(assumptions))
	}
}

func TestDomainPackControls(t *testing.T) {
	de := NewDomainEngine()
	pack := de.GetPack("Fintech")
	if pack == nil {
		t.Fatal("expected Fintech pack")
	}
	if len(pack.Controls) == 0 {
		t.Error("expected controls in Fintech pack")
	}
	if len(pack.ComplianceMappings) == 0 {
		t.Error("expected compliance mappings in Fintech pack")
	}
	if len(pack.RiskAmplifiers) == 0 {
		t.Error("expected risk amplifiers in Fintech pack")
	}
}

func TestDomainPackAllPacksExist(t *testing.T) {
	de := NewDomainEngine()
	required := []string{"Healthcare", "Fintech", "SaaS", "Enterprise", "Kubernetes", "CloudNative", "VPN", "IdentityPlatform", "DataPlatform"}
	for _, name := range required {
		if de.GetPack(name) == nil {
			t.Errorf("missing domain pack: %s", name)
		}
	}
}

func TestDomainEngineDetectEnterprise(t *testing.T) {
	de := NewDomainEngine()
	arch := &ArchDescription{
		Name: "enterprise",
		Components: []Component{
			{ID: "ad1", Label: "ActiveDirectory"},
			{ID: "app1", Label: "CorporateApp"},
		},
		RawText: "Enterprise identity lifecycle and access reviews.",
	}
	domain := de.DetectDomain(arch)
	if domain != "Enterprise" {
		t.Errorf("expected Enterprise domain, got %s", domain)
	}
}
