package fidelity

import (
	"asf-tui/asf/fact"
	"asf-tui/asf/fidelity"
	"fmt"
	"os"
	"testing"
)

// BenchmarkResult holds the results for a single benchmark.
type BenchmarkResult struct {
	Name                  string
	TotalFacts            int
	RespectedFacts        int
	ContradictedFacts     int
	TotalAssumptions      int
	ValidAssumptions      int
	InvalidAssumptions    int
	TotalContradictions   int
	RealContradictions    int
	FalseContradictions   int
	FidelityScore         float64
	AssumptionQuality     float64
	ContradictionAccuracy float64
	NoveltyScore          float64
	Overall               string
	Passed                bool
}

// RunBenchmark runs a single benchmark.
func RunBenchmark(name string, facts []fact.Fact, components []fidelity.Component, relationships []fidelity.Relationship, domain string, expectedAssumptions int, expectedContradictions int) BenchmarkResult {
	// Create inventory
	inv := fact.NewInventory()
	for _, f := range facts {
		inv.Add(f)
	}

	// Generate hidden assumptions
	engine := fidelity.NewHiddenAssumptionEngine(inv, domain)
	assumptions := engine.Generate(inv, components, relationships)

	// Detect contradictions
	contradictionEngine := fidelity.NewRealContradictionEngine(inv)
	contradictions := contradictionEngine.Detect()
	contradictions = append(contradictions, contradictionEngine.DetectFactAssumption(assumptions)...)

	// Build traceability
	traceabilityEngine := fidelity.NewTraceabilityEngine()
	traceability := traceabilityEngine.BuildTraceability(assumptions, inv)

	// Score
	scorer := fidelity.NewFidelityScorer(inv)
	score := scorer.Compute(assumptions, contradictions, traceability)

	// Validate assumptions
	valid, invalid, _ := traceabilityEngine.ValidateTraceability(assumptions)

	// Validate contradictions
	realCount := 0
	falseCount := 0
	for _, c := range contradictions {
		if c.Type == "fact-fact" || c.Type == "fact-assumption" {
			realCount++
		} else {
			falseCount++
		}
	}

	passed := score.Score >= 0.9 && score.AssumptionQuality >= 0.7 && score.ContradictionAccuracy >= 0.9

	return BenchmarkResult{
		Name:                  name,
		TotalFacts:            score.TotalFacts,
		RespectedFacts:        score.RespectedFacts,
		ContradictedFacts:     score.ContradictedFacts,
		TotalAssumptions:      len(assumptions),
		ValidAssumptions:      valid,
		InvalidAssumptions:    invalid,
		TotalContradictions:   len(contradictions),
		RealContradictions:    realCount,
		FalseContradictions:   falseCount,
		FidelityScore:         score.Score,
		AssumptionQuality:     score.AssumptionQuality,
		ContradictionAccuracy: score.ContradictionAccuracy,
		NoveltyScore:          score.NoveltyScore,
		Overall:               score.Overall,
		Passed:                passed,
	}
}

// BenchmarkHealthcare runs the healthcare benchmark.
func BenchmarkHealthcare(b *testing.B) {
	facts := []fact.Fact{
		{ID: "f1", Text: "MFA is enabled for all users", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f2", Text: "Encryption is enabled for data at rest and in transit", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f3", Text: "HIPAA compliance is required", Source: "yaml", FactType: "compliance", Category: "compliance", IsNegative: false},
		{ID: "f4", Text: "Backups are automated daily", Source: "yaml", FactType: "control", Category: "availability", IsNegative: false},
		{ID: "f5", Text: "WAF is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f6", Text: "Audit logging is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f7", Text: "VPN is used for admin access", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f8", Text: "Role-based access control is enforced", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
	}

	components := []fidelity.Component{
		{ID: "comp1", Label: "Patient Database"},
		{ID: "comp2", Label: "API Gateway"},
		{ID: "comp3", Label: "Auth0 Service"},
		{ID: "comp4", Label: "Audit Log"},
		{ID: "comp5", Label: "Load Balancer"},
		{ID: "comp6", Label: "CDN"},
		{ID: "comp7", Label: "Message Queue"},
	}

	relationships := []fidelity.Relationship{
		{Source: "API Gateway", Target: "Patient Database", Label: "queries"},
		{Source: "Auth0 Service", Target: "API Gateway", Label: "authenticates"},
		{Source: "Load Balancer", Target: "API Gateway", Label: "routes"},
		{Source: "CDN", Target: "API Gateway", Label: "caches"},
		{Source: "Admin", Target: "Patient Database", Label: "manages"},
	}

	result := RunBenchmark("Healthcare", facts, components, relationships, "healthcare", 15, 0)
	reportBenchmarkResult(b, result)
}

// BenchmarkFintech runs the fintech benchmark.
func BenchmarkFintech(b *testing.B) {
	facts := []fact.Fact{
		{ID: "f1", Text: "MFA is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f2", Text: "Encryption is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f3", Text: "PCI DSS compliance is required", Source: "yaml", FactType: "compliance", Category: "compliance", IsNegative: false},
		{ID: "f4", Text: "Backups are enabled", Source: "yaml", FactType: "control", Category: "availability", IsNegative: false},
		{ID: "f5", Text: "WAF is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f6", Text: "Audit logging is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f7", Text: "Fraud detection is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f8", Text: "Tokenization is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
	}

	components := []fidelity.Component{
		{ID: "comp1", Label: "Payment Processor"},
		{ID: "comp2", Label: "API Gateway"},
		{ID: "comp3", Label: "Fraud Detection Service"},
		{ID: "comp4", Label: "Token Vault"},
		{ID: "comp5", Label: "Database"},
		{ID: "comp6", Label: "Load Balancer"},
	}

	relationships := []fidelity.Relationship{
		{Source: "API Gateway", Target: "Payment Processor", Label: "routes"},
		{Source: "Fraud Detection Service", Target: "Payment Processor", Label: "monitors"},
		{Source: "Payment Processor", Target: "Token Vault", Label: "stores tokens"},
		{Source: "Payment Processor", Target: "Database", Label: "stores transactions"},
		{Source: "Load Balancer", Target: "API Gateway", Label: "routes"},
	}

	result := RunBenchmark("Fintech", facts, components, relationships, "fintech", 10, 0)
	reportBenchmarkResult(b, result)
}

// BenchmarkCloud runs the cloud benchmark.
func BenchmarkCloud(b *testing.B) {
	facts := []fact.Fact{
		{ID: "f1", Text: "MFA is enabled for all IAM users", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f2", Text: "Encryption is enabled (KMS)", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f3", Text: "AWS GuardDuty is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f4", Text: "CloudTrail is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f5", Text: "VPC is configured with private subnets", Source: "yaml", FactType: "control", Category: "network", IsNegative: false},
		{ID: "f6", Text: "Security groups are restricted", Source: "yaml", FactType: "control", Category: "network", IsNegative: false},
		{ID: "f7", Text: "IAM policies use least privilege", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f8", Text: "AWS WAF is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f9", Text: "S3 buckets are encrypted and versioned", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f10", Text: "RDS has automated backups", Source: "yaml", FactType: "control", Category: "availability", IsNegative: false},
		{ID: "f11", Text: "CloudWatch monitoring is enabled", Source: "yaml", FactType: "control", Category: "monitoring", IsNegative: false},
		{ID: "f12", Text: "AWS Config is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
	}

	components := []fidelity.Component{
		{ID: "comp1", Label: "API Gateway"},
		{ID: "comp2", Label: "Lambda Functions"},
		{ID: "comp3", Label: "RDS Database"},
		{ID: "comp4", Label: "S3 Storage"},
		{ID: "comp5", Label: "CloudFront CDN"},
		{ID: "comp6", Label: "EC2 Instances"},
		{ID: "comp7", Label: "EKS Cluster"},
		{ID: "comp8", Label: "DynamoDB"},
	}

	relationships := []fidelity.Relationship{
		{Source: "API Gateway", Target: "Lambda Functions", Label: "invokes"},
		{Source: "Lambda Functions", Target: "RDS Database", Label: "queries"},
		{Source: "Lambda Functions", Target: "S3 Storage", Label: "reads/writes"},
		{Source: "Lambda Functions", Target: "DynamoDB", Label: "queries"},
		{Source: "CloudFront CDN", Target: "API Gateway", Label: "routes"},
		{Source: "EKS Cluster", Target: "RDS Database", Label: "queries"},
		{Source: "EC2 Instances", Target: "S3 Storage", Label: "reads"},
	}

	result := RunBenchmark("Cloud", facts, components, relationships, "cloud", 15, 0)
	reportBenchmarkResult(b, result)
}

// BenchmarkKubernetes runs the kubernetes benchmark.
func BenchmarkKubernetes(b *testing.B) {
	facts := []fact.Fact{
		{ID: "f1", Text: "RBAC is enforced", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f2", Text: "Network policies are configured", Source: "yaml", FactType: "control", Category: "network", IsNegative: false},
		{ID: "f3", Text: "Pod security policies are enforced", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f4", Text: "Secrets are encrypted at rest", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f5", Text: "Admission controllers are enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f6", Text: "Container images are scanned", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f7", Text: "Resource quotas are set", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f8", Text: "Node auto-scaling is enabled", Source: "yaml", FactType: "control", Category: "availability", IsNegative: false},
		{ID: "f9", Text: "Cluster logging is enabled", Source: "yaml", FactType: "control", Category: "monitoring", IsNegative: false},
		{ID: "f10", Text: "Cluster monitoring is enabled", Source: "yaml", FactType: "control", Category: "monitoring", IsNegative: false},
	}

	components := []fidelity.Component{
		{ID: "comp1", Label: "API Server"},
		{ID: "comp2", Label: "etcd"},
		{ID: "comp3", Label: "Ingress Controller"},
		{ID: "comp4", Label: "Service Mesh"},
		{ID: "comp5", Label: "Application Pods"},
		{ID: "comp6", Label: "Monitoring Stack"},
		{ID: "comp7", Label: "Logging Stack"},
		{ID: "comp8", Label: "CI/CD Pipeline"},
	}

	relationships := []fidelity.Relationship{
		{Source: "Ingress Controller", Target: "API Server", Label: "routes"},
		{Source: "API Server", Target: "etcd", Label: "stores data"},
		{Source: "Service Mesh", Target: "Application Pods", Label: "routes traffic"},
		{Source: "CI/CD Pipeline", Target: "Application Pods", Label: "deploys"},
		{Source: "Monitoring Stack", Target: "Application Pods", Label: "monitors"},
		{Source: "Logging Stack", Target: "Application Pods", Label: "collects logs"},
	}

	result := RunBenchmark("Kubernetes", facts, components, relationships, "kubernetes", 15, 0)
	reportBenchmarkResult(b, result)
}

// BenchmarkVPN runs the VPN benchmark.
func BenchmarkVPN(b *testing.B) {
	facts := []fact.Fact{
		{ID: "f1", Text: "VPN is enabled", Source: "yaml", FactType: "control", Category: "network", IsNegative: false},
		{ID: "f2", Text: "MFA is enabled for VPN access", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f3", Text: "VPN logs are enabled", Source: "yaml", FactType: "control", Category: "monitoring", IsNegative: false},
		{ID: "f4", Text: "VPN certificates are rotated", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f5", Text: "Split tunneling is disabled", Source: "yaml", FactType: "control", Category: "network", IsNegative: true},
		{ID: "f6", Text: "VPN client is managed", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f7", Text: "Network segmentation is enforced", Source: "yaml", FactType: "control", Category: "network", IsNegative: false},
		{ID: "f8", Text: "Firewall rules are strict", Source: "yaml", FactType: "control", Category: "network", IsNegative: false},
	}

	components := []fidelity.Component{
		{ID: "comp1", Label: "VPN Gateway"},
		{ID: "comp2", Label: "Authentication Server"},
		{ID: "comp3", Label: "Internal Network"},
		{ID: "comp4", Label: "Firewall"},
		{ID: "comp5", Label: "Logging Server"},
		{ID: "comp6", Label: "Certificate Authority"},
		{ID: "comp7", Label: "DNS Server"},
	}

	relationships := []fidelity.Relationship{
		{Source: "VPN Gateway", Target: "Authentication Server", Label: "authenticates"},
		{Source: "VPN Gateway", Target: "Internal Network", Label: "routes"},
		{Source: "Firewall", Target: "Internal Network", Label: "protects"},
		{Source: "Logging Server", Target: "VPN Gateway", Label: "collects logs"},
		{Source: "Certificate Authority", Target: "VPN Gateway", Label: "issues certs"},
		{Source: "DNS Server", Target: "VPN Gateway", Label: "resolves"},
	}

	result := RunBenchmark("VPN", facts, components, relationships, "vpn", 10, 0)
	reportBenchmarkResult(b, result)
}

// BenchmarkSaaS runs the SaaS benchmark.
func BenchmarkSaaS(b *testing.B) {
	facts := []fact.Fact{
		{ID: "f1", Text: "MFA is enabled for tenant admins", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f2", Text: "Encryption is enabled for tenant data", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f3", Text: "Tenant isolation is enforced", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f4", Text: "API rate limiting is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f5", Text: "Audit logging is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f6", Text: "DLP is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		{ID: "f7", Text: "Data retention policies are enforced", Source: "yaml", FactType: "control", Category: "compliance", IsNegative: false},
		{ID: "f8", Text: "Backup is enabled", Source: "yaml", FactType: "control", Category: "availability", IsNegative: false},
		{ID: "f9", Text: "Monitoring is enabled", Source: "yaml", FactType: "control", Category: "monitoring", IsNegative: false},
		{ID: "f10", Text: "Penetration testing is performed", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
	}

	components := []fidelity.Component{
		{ID: "comp1", Label: "Tenant Portal"},
		{ID: "comp2", Label: "API Gateway"},
		{ID: "comp3", Label: "Application Server"},
		{ID: "comp4", Label: "Database"},
		{ID: "comp5", Label: "Object Storage"},
		{ID: "comp6", Label: "Cache"},
		{ID: "comp7", Label: "Message Queue"},
		{ID: "comp8", Label: "CDN"},
		{ID: "comp9", Label: "Monitoring Stack"},
		{ID: "comp10", Label: "Backup Service"},
	}

	relationships := []fidelity.Relationship{
		{Source: "Tenant Portal", Target: "API Gateway", Label: "routes"},
		{Source: "API Gateway", Target: "Application Server", Label: "routes"},
		{Source: "Application Server", Target: "Database", Label: "queries"},
		{Source: "Application Server", Target: "Object Storage", Label: "reads/writes"},
		{Source: "Application Server", Target: "Cache", Label: "reads/writes"},
		{Source: "Application Server", Target: "Message Queue", Label: "publishes"},
		{Source: "CDN", Target: "Tenant Portal", Label: "caches"},
		{Source: "Monitoring Stack", Target: "Application Server", Label: "monitors"},
		{Source: "Backup Service", Target: "Database", Label: "backs up"},
	}

	result := RunBenchmark("SaaS", facts, components, relationships, "saas", 15, 0)
	reportBenchmarkResult(b, result)
}

func reportBenchmarkResult(b *testing.B, result BenchmarkResult) {
	b.Logf("\n=== %s Benchmark ===", result.Name)
	b.Logf("Fidelity Score: %.1f%%", result.FidelityScore*100)
	b.Logf("Assumption Quality: %.1f%%", result.AssumptionQuality*100)
	b.Logf("Contradiction Accuracy: %.1f%%", result.ContradictionAccuracy*100)
	b.Logf("Novelty Score: %.1f%%", result.NoveltyScore*100)
	b.Logf("Overall: %s", result.Overall)
	b.Logf("Facts: %d total, %d respected, %d contradicted", result.TotalFacts, result.RespectedFacts, result.ContradictedFacts)
	b.Logf("Assumptions: %d total, %d valid, %d invalid", result.TotalAssumptions, result.ValidAssumptions, result.InvalidAssumptions)
	b.Logf("Contradictions: %d total, %d real, %d false", result.TotalContradictions, result.RealContradictions, result.FalseContradictions)
	b.Logf("Passed: %v", result.Passed)
}

// TestAllBenchmarks runs all benchmarks in sequence.
func TestAllBenchmarks(t *testing.T) {
	var results []BenchmarkResult

	// Healthcare
	{
		facts := []fact.Fact{
			{ID: "f1", Text: "MFA is enabled for all users", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
			{ID: "f2", Text: "Encryption is enabled for data at rest and in transit", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
			{ID: "f3", Text: "HIPAA compliance is required", Source: "yaml", FactType: "compliance", Category: "compliance", IsNegative: false},
			{ID: "f4", Text: "Backups are automated daily", Source: "yaml", FactType: "control", Category: "availability", IsNegative: false},
			{ID: "f5", Text: "WAF is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
			{ID: "f6", Text: "Audit logging is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
			{ID: "f7", Text: "VPN is used for admin access", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		}
		components := []fidelity.Component{
			{ID: "comp1", Label: "Patient Database"},
			{ID: "comp2", Label: "API Gateway"},
			{ID: "comp3", Label: "Auth0 Service"},
			{ID: "comp4", Label: "Audit Log"},
			{ID: "comp5", Label: "Load Balancer"},
		}
		relationships := []fidelity.Relationship{
			{Source: "API Gateway", Target: "Patient Database", Label: "queries"},
			{Source: "Auth0 Service", Target: "API Gateway", Label: "authenticates"},
		}
		results = append(results, RunBenchmark("Healthcare", facts, components, relationships, "healthcare", 15, 0))
	}

	// Fintech
	{
		facts := []fact.Fact{
			{ID: "f1", Text: "MFA is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
			{ID: "f2", Text: "Encryption is enabled", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
			{ID: "f3", Text: "PCI DSS compliance is required", Source: "yaml", FactType: "compliance", Category: "compliance", IsNegative: false},
		}
		components := []fidelity.Component{
			{ID: "comp1", Label: "Payment Processor"},
			{ID: "comp2", Label: "API Gateway"},
		}
		results = append(results, RunBenchmark("Fintech", facts, components, nil, "fintech", 10, 0))
	}

	// Cloud
	{
		facts := []fact.Fact{
			{ID: "f1", Text: "MFA is enabled for all IAM users", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
			{ID: "f2", Text: "Encryption is enabled (KMS)", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		}
		components := []fidelity.Component{
			{ID: "comp1", Label: "API Gateway"},
			{ID: "comp2", Label: "Lambda Functions"},
		}
		results = append(results, RunBenchmark("Cloud", facts, components, nil, "cloud", 15, 0))
	}

	// Kubernetes
	{
		facts := []fact.Fact{
			{ID: "f1", Text: "RBAC is enforced", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
			{ID: "f2", Text: "Network policies are configured", Source: "yaml", FactType: "control", Category: "network", IsNegative: false},
		}
		components := []fidelity.Component{
			{ID: "comp1", Label: "API Server"},
			{ID: "comp2", Label: "etcd"},
		}
		results = append(results, RunBenchmark("Kubernetes", facts, components, nil, "kubernetes", 15, 0))
	}

	// VPN
	{
		facts := []fact.Fact{
			{ID: "f1", Text: "VPN is enabled", Source: "yaml", FactType: "control", Category: "network", IsNegative: false},
			{ID: "f2", Text: "MFA is enabled for VPN access", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		}
		components := []fidelity.Component{
			{ID: "comp1", Label: "VPN Gateway"},
			{ID: "comp2", Label: "Authentication Server"},
		}
		results = append(results, RunBenchmark("VPN", facts, components, nil, "vpn", 10, 0))
	}

	// SaaS
	{
		facts := []fact.Fact{
			{ID: "f1", Text: "MFA is enabled for tenant admins", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
			{ID: "f2", Text: "Encryption is enabled for tenant data", Source: "yaml", FactType: "control", Category: "security", IsNegative: false},
		}
		components := []fidelity.Component{
			{ID: "comp1", Label: "Tenant Portal"},
			{ID: "comp2", Label: "API Gateway"},
		}
		results = append(results, RunBenchmark("SaaS", facts, components, nil, "saas", 15, 0))
	}

	// Report summary
	t.Log("\n=== Benchmark Summary ===")
	allPassed := true
	for _, r := range results {
		if !r.Passed {
			allPassed = false
		}
		t.Logf("%s: Fidelity=%.1f%%, Quality=%.1f%%, Accuracy=%.1f%%, Novelty=%.1f%%, %s, Passed=%v",
			r.Name, r.FidelityScore*100, r.AssumptionQuality*100, r.ContradictionAccuracy*100, r.NoveltyScore*100, r.Overall, r.Passed)
	}

	if !allPassed {
		t.Log("Some benchmarks did not pass. Review fidelity scores.")
	}

	// Write report
	report := "# Architectural Fidelity Benchmark Report\n\n"
	for _, r := range results {
		report += fmt.Sprintf("## %s\n\n", r.Name)
		report += fmt.Sprintf("- Fidelity Score: %.1f%%\n", r.FidelityScore*100)
		report += fmt.Sprintf("- Assumption Quality: %.1f%%\n", r.AssumptionQuality*100)
		report += fmt.Sprintf("- Contradiction Accuracy: %.1f%%\n", r.ContradictionAccuracy*100)
		report += fmt.Sprintf("- Novelty Score: %.1f%%\n", r.NoveltyScore*100)
		report += fmt.Sprintf("- Overall: %s\n", r.Overall)
		report += fmt.Sprintf("- Passed: %v\n\n", r.Passed)
	}

	os.WriteFile("/Users/moksh/Project/cybersec/asf-tui/benchmark/fidelity/BENCHMARK_REPORT.md", []byte(report), 0644)
}
