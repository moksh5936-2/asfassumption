package intelligence

import (
	"strings"
	"testing"
)

func TestDKPIDetectHealthcare(t *testing.T) {
	d := NewDKPIDetector()
	arch := &ArchDescription{
		Name: "health",
		Components: []Component{
			{ID: "phi1", Label: "PHIDatabase"},
			{ID: "ehr1", Label: "EHRSystem"},
		},
		RawText: "HIPAA compliant system storing PHI.",
	}
	res := d.DetectDomain(arch)
	if res.PrimaryDomain != "healthcare" {
		t.Errorf("expected healthcare domain, got %s", res.PrimaryDomain)
	}
	if res.Confidence < 10 {
		t.Errorf("confidence too low: %.1f%%", res.Confidence)
	}
	if len(res.Rationale) == 0 {
		t.Error("expected rationale strings")
	}
}

func TestDKPIDetectFintech(t *testing.T) {
	d := NewDKPIDetector()
	arch := &ArchDescription{
		Name: "fintech",
		Components: []Component{
			{ID: "pay1", Label: "PaymentGateway"},
			{ID: "txn1", Label: "TransactionDB"},
		},
		RawText: "PCI DSS compliant payment processing with fraud detection.",
	}
	res := d.DetectDomain(arch)
	if res.PrimaryDomain != "fintech" {
		t.Errorf("expected fintech domain, got %s", res.PrimaryDomain)
	}
}

func TestDKPIDetectKubernetes(t *testing.T) {
	d := NewDKPIDetector()
	arch := &ArchDescription{
		Name: "k8s",
		Components: []Component{
			{ID: "pod1", Label: "Pod"},
			{ID: "ns1", Label: "Namespace"},
		},
		RawText: "Kubernetes cluster with RBAC and network policies.",
	}
	res := d.DetectDomain(arch)
	if res.PrimaryDomain != "kubernetes" {
		t.Errorf("expected kubernetes domain, got %s", res.PrimaryDomain)
	}
}

func TestDKPIDetectSaaS(t *testing.T) {
	d := NewDKPIDetector()
	arch := &ArchDescription{
		Name: "saas",
		Components: []Component{
			{ID: "t1", Label: "TenantManager"},
			{ID: "api1", Label: "APIGateway"},
		},
		RawText: "Multi-tenant SaaS platform with SSO integration.",
	}
	res := d.DetectDomain(arch)
	if res.PrimaryDomain != "saas" {
		t.Errorf("expected saas domain, got %s", res.PrimaryDomain)
	}
}

func TestDKPIDetectGovernment(t *testing.T) {
	d := NewDKPIDetector()
	arch := &ArchDescription{
		Name: "gov",
		Components: []Component{
			{ID: "cit1", Label: "CitizenDB"},
		},
		RawText: "FedRAMP compliant federal government system handling classified data.",
	}
	res := d.DetectDomain(arch)
	if res.PrimaryDomain != "government" {
		t.Errorf("expected government domain, got %s", res.PrimaryDomain)
	}
}

func TestDKPIDetectCriticalInfrastructure(t *testing.T) {
	d := NewDKPIDetector()
	arch := &ArchDescription{
		Name: "industrial",
		Components: []Component{
			{ID: "scada1", Label: "SCADAController"},
			{ID: "hmi1", Label: "HMI"},
		},
		RawText: "Industrial control system with SCADA and PLC.",
	}
	res := d.DetectDomain(arch)
	if res.PrimaryDomain != "critical_infrastructure" {
		t.Errorf("expected critical_infrastructure domain, got %s", res.PrimaryDomain)
	}
}

func TestDKPINoDetection(t *testing.T) {
	d := NewDKPIDetector()
	arch := &ArchDescription{
		Name:    "generic",
		RawText: "A generic application with standard components.",
	}
	res := d.DetectDomain(arch)
	if res.PrimaryDomain != "" {
		t.Errorf("expected no domain, got %s", res.PrimaryDomain)
	}
}

func TestDKPIEngineRunHealthcare(t *testing.T) {
	e := NewDKPIEngine()
	arch := &ArchDescription{
		Name: "health",
		Components: []Component{
			{ID: "phi1", Label: "PHIDatabase"},
		},
		RawText: "HIPAA compliant healthcare system.",
	}
	input := DKPIInput{
		Architecture: arch,
	}
	res := e.Run(input)
	if res == nil {
		t.Fatal("expected non-nil result")
	}
	if res.DetectedDomain.PrimaryDomain != "healthcare" {
		t.Errorf("expected healthcare domain, got %s", res.DetectedDomain.PrimaryDomain)
	}
	if res.ActivePack == nil {
		t.Fatal("expected active pack")
	}
	if len(res.InjectedThreats) == 0 {
		t.Error("expected injected threats")
	}
	if len(res.DomainControls) == 0 {
		t.Error("expected domain controls")
	}
	if len(res.Recommendations) == 0 {
		t.Error("expected recommendations")
	}
	if len(res.DomainCompliance) == 0 {
		t.Error("expected domain compliance")
	}
	if len(res.DomainEvidence) == 0 {
		t.Error("expected domain evidence")
	}
}

func TestDKPIEngineRunFintech(t *testing.T) {
	e := NewDKPIEngine()
	arch := &ArchDescription{
		Name: "fintech",
		Components: []Component{
			{ID: "pay1", Label: "PaymentGateway"},
		},
		RawText: "PCI DSS compliant payment system.",
	}
	input := DKPIInput{Architecture: arch}
	res := e.Run(input)
	if res == nil || res.ActivePack == nil {
		t.Fatal("expected active pack for fintech")
	}
	if len(res.InjectedThreats) == 0 {
		t.Error("expected fintech threats")
	}
	hasPCI := false
	for _, c := range res.DomainCompliance {
		if strings.Contains(strings.ToUpper(c), "PCI") {
			hasPCI = true
			break
		}
	}
	if !hasPCI {
		t.Error("expected PCI compliance in fintech pack")
	}
}

func TestDKPIEngineRunFallbackDomain(t *testing.T) {
	e := NewDKPIEngine()
	arch := &ArchDescription{
		Name:    "custom",
		RawText: "Custom architecture.",
	}
	input := DKPIInput{
		Architecture: arch,
		Domain:       "Healthcare",
	}
	res := e.Run(input)
	if res == nil {
		t.Fatal("expected non-nil result")
	}
	if res.DetectedDomain.PrimaryDomain == "" {
		t.Error("expected fallback domain detection")
	}
}

func TestDKPIThreatDeduplication(t *testing.T) {
	e := NewDKPIEngine()
	arch := &ArchDescription{
		Name: "health",
		Components: []Component{
			{ID: "phi1", Label: "PHIDatabase"},
		},
		RawText: "HIPAA healthcare system.",
	}
	input := DKPIInput{
		Architecture:    arch,
		ExistingThreats: []Threat{}, // empty threats - should still inject
	}
	res := e.Run(input)
	if res == nil || res.ActivePack == nil {
		t.Fatal("expected active pack")
	}
	if len(res.InjectedThreats) < 1 {
		t.Error("expected at least 1 injected threat")
	}
}

func TestDKPIBoostConfidence(t *testing.T) {
	assumptions := []Assumption{
		{ID: "ASM-001", Description: "PHI is encrypted at rest", Confidence: 0.70},
		{ID: "ASM-002", Description: "Access is logged", Confidence: 0.50},
	}
	boosted := boostDomainConfidence(assumptions, "healthcare", 85.0)
	if len(boosted) != len(assumptions) {
		t.Errorf("expected %d assumptions, got %d", len(assumptions), len(boosted))
	}
	for _, a := range boosted {
		if a.Confidence < 0.70 {
			t.Errorf("assumption %s confidence unexpectedly low: %.2f", a.ID, a.Confidence)
		}
		if !containsString(a.Keywords, "domain:healthcare") {
			t.Error("expected domain keyword on boosted assumption")
		}
	}
}

func TestDKPIEnrichControlStrength(t *testing.T) {
	controls := []SDRIControl{
		{ID: "CTRL-001", Name: "MFA", Category: "Access Control", Coverage: "Partial"},
	}
	pack := &KnowledgePack{
		ID:   "healthcare",
		Name: "Healthcare",
		ExpectedControls: []KnowledgePackControl{
			{Name: "MFA", Category: "Identity", Priority: "High"},
		},
	}
	enriched := enrichControlStrength(controls, pack)
	if len(enriched) != 1 {
		t.Fatal("expected 1 enriched control")
	}
	if enriched[0].Category != "Identity" {
		t.Errorf("expected category 'Identity', got '%s'", enriched[0].Category)
	}
	if enriched[0].Coverage != "Enhanced" {
		t.Errorf("expected coverage 'Enhanced', got '%s'", enriched[0].Coverage)
	}
}

func TestDKPIAllPacksHaveExpectedStructure(t *testing.T) {
	packs := buildKnowledgePacks()
	for _, p := range packs {
		t.Run(p.ID, func(t *testing.T) {
			if p.Name == "" {
				t.Error("pack missing name")
			}
			if p.Industry == "" {
				t.Error("pack missing industry")
			}
			if len(p.DetectionKeywords) == 0 {
				t.Error("pack missing detection keywords")
			}
			if len(p.CrownJewels) == 0 {
				t.Error("pack missing crown jewels")
			}
			if len(p.ExpectedControls) == 0 {
				t.Error("pack missing expected controls")
			}
			if len(p.ThreatPatterns) == 0 {
				t.Error("pack missing threat patterns")
			}
			if len(p.AttackPathTemplates) == 0 {
				t.Error("pack missing attack path templates")
			}
			if len(p.AssumptionPatterns) == 0 {
				t.Error("pack missing assumption patterns")
			}
		})
	}
}

func TestDKPIEngineRunAllPacks(t *testing.T) {
	packs := buildKnowledgePacks()
	e := NewDKPIEngine()
	for _, p := range packs {
		t.Run(p.ID, func(t *testing.T) {
			arch := &ArchDescription{
				Name: p.ID,
				Components: []Component{
					{ID: "comp1", Label: p.CrownJewels[0]},
				},
				RawText: strings.Join(p.DetectionKeywords[:minInt(3, len(p.DetectionKeywords))], " "),
			}
			input := DKPIInput{Architecture: arch}
			res := e.Run(input)
			if res == nil {
				t.Fatal("expected result")
			}
			if res.ActivePack == nil {
				t.Fatalf("expected active pack for %s", p.ID)
			}
			if len(res.InjectedThreats) == 0 {
				t.Errorf("expected threats for %s", p.ID)
			}
			if len(res.DomainControls) == 0 {
				t.Errorf("expected controls for %s", p.ID)
			}
			if len(res.Recommendations) == 0 {
				t.Errorf("expected recommendations for %s", p.ID)
			}
			if len(res.DomainCompliance) == 0 {
				t.Errorf("expected compliance for %s", p.ID)
			}
		})
	}
}

func containsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}
