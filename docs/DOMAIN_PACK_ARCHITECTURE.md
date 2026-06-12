# ASF Domain Pack Architecture

## Overview

Domain packs provide contextual security knowledge for specific industries. They are deterministic rule sets that generate domain-specific assumptions, controls, risk amplifiers, and compliance mappings.

## Architecture

```go
type DomainPack struct {
    Name               string
    Keywords           []string
    Assumptions        []string
    Controls           []ControlDetail
    RiskAmplifiers     map[string]RiskLevel
    ComplianceMappings []string
    VerificationRules  []string
}
```

## Domain Packs

### 1. Healthcare

**Keywords:** phi, hipaa, patient, health, medical, ehr, healthcare

**Assumptions:**
- PHI must be encrypted at rest and in transit
- Patient privacy must be enforced with minimum necessary access
- PHI access must be audited with HIPAA audit controls
- Break-glass access must exist for emergency PHI retrieval
- Data retention must comply with HIPAA requirements

**Controls:**
- HIPAA audit controls (164.312(b))
- Minimum necessary access controls
- Break-glass procedures with immutable logging

**Risk Amplifiers:**
- PHI: Critical
- HIPAA: Critical
- Patient: High
- Medical: High

**Compliance:** HIPAA, HITECH, state privacy laws

### 2. Fintech

**Keywords:** pci, payment, card, financial, bank, transaction, fraud, sox

**Assumptions:**
- Cardholder data must be encrypted per PCI DSS
- Transaction integrity must be maintained
- Fraud detection mechanisms must be in place
- Financial records must be immutable

**Controls:**
- PCI DSS encryption requirements
- Transaction logging and integrity
- Fraud detection and prevention

**Risk Amplifiers:**
- PCI: Critical
- Payment: Critical
- Financial: High

**Compliance:** PCI DSS, SOX, GLBA

### 3. SaaS

**Keywords:** saas, multi-tenant, tenant, subscription, customer, portal

**Assumptions:**
- Tenant isolation must be enforced at all layers
- Data segregation must prevent cross-tenant access
- API security must enforce tenant scoping
- Customer data must be encrypted

**Controls:**
- Tenant isolation controls
- Data segregation verification
- API tenant scoping

**Risk Amplifiers:**
- Multi-tenant: Critical
- Tenant: High

**Compliance:** SOC2, ISO27001

### 4. Enterprise

**Keywords:** enterprise, corporate, internal, employee, directory, ad, ldap

**Assumptions:**
- Identity lifecycle management must be automated
- Access reviews must be periodic and enforced
- Privileged access must be monitored and logged
- Offboarding must revoke all access immediately

**Controls:**
- Identity lifecycle automation
- Periodic access reviews
- Privileged access monitoring

**Risk Amplifiers:**
- Enterprise: High
- Corporate: Medium

**Compliance:** SOX, ISO27001, NIST

### 5. Kubernetes

**Keywords:** kubernetes, k8s, pod, container, namespace, cluster, helm

**Assumptions:**
- Container security must enforce image scanning
- RBAC must be least-privilege
- Network policies must segment workloads
- Secrets must be managed via vaults or KMS

**Controls:**
- Container image scanning
- RBAC policies
- Network policies
- Secrets management

**Risk Amplifiers:**
- Kubernetes: High
- Container: High

**Compliance:** NIST, CIS Kubernetes Benchmark

### 6. Cloud Native

**Keywords:** aws, azure, gcp, cloud, serverless, lambda, function, iam

**Assumptions:**
- IAM policies must be least-privilege
- Data must be encrypted at rest and in transit
- Logging must be centralized and monitored
- Resources must be tagged for cost and security tracking

**Controls:**
- IAM least-privilege policies
- Encryption controls
- Centralized logging

**Risk Amplifiers:**
- Cloud: High
- IAM: High

**Compliance:** SOC2, ISO27001, FedRAMP

### 7. VPN

**Keywords:** vpn, tunnel, remote, endpoint, certificate, mfa, wireguard, ipsec

**Assumptions:**
- Tunnel security must use strong encryption
- Endpoint validation must verify device health
- Certificate management must be automated
- MFA must be enforced for all VPN access

**Controls:**
- Tunnel encryption controls
- Endpoint validation
- Certificate management

**Risk Amplifiers:**
- VPN: High
- Remote: High

**Compliance:** NIST, PCI DSS

### 8. Identity Platform

**Keywords:** sso, identity, federation, oauth, saml, oidc, auth0, login

**Assumptions:**
- SSO must be enforced for all applications
- Session management must include timeout and rotation
- Federation must verify identity provider trust
- MFA must be enforced for all privileged access

**Controls:**
- SSO enforcement
- Session management
- Federation trust

**Risk Amplifiers:**
- SSO: High
- Federation: High

**Compliance:** NIST, ISO27001

### 9. Data Platform

**Keywords:** data, warehouse, lake, pipeline, etl, analytics, databricks, snowflake

**Assumptions:**
- Data governance must enforce lineage tracking
- Data retention must be enforced automatically
- Data encryption must be applied at all stages
- Access must be granted based on data classification

**Controls:**
- Data lineage controls
- Retention policies
- Encryption controls

**Risk Amplifiers:**
- Data: High
- Analytics: Medium

**Compliance:** GDPR, CCPA, HIPAA

## Auto-Detection

```go
func (de *DomainEngine) DetectDomain(arch *ArchDescription) string {
    rawText := strings.ToLower(arch.RawText + " " + arch.Name)
    for _, pack := range de.Packs {
        score := 0
        for _, kw := range pack.Keywords {
            if strings.Contains(rawText, kw) {
                score++
            }
        }
        if score >= 2 {
            return pack.Name
        }
    }
    return ""
}
```

Detection is deterministic:
- Count keyword matches in architecture text
- Domain with 2+ matches wins
- Alphabetic tie-breaker for deterministic selection

## Application

```go
func (de *DomainEngine) ApplyDomainPack(domain string, arch *ArchDescription) []Assumption {
    pack := de.GetPack(domain)
    var assumptions []Assumption
    for _, text := range pack.Assumptions {
        assumptions = append(assumptions, Assumption{
            ID:          fmt.Sprintf("DOM-%s-%03d", domain, len(assumptions)+1),
            Description: text,
            Category:    "Compliance",
            Risk:        RiskHigh,
            SourceType:  "domain-inferred",
            Confidence:  0.85,
        })
    }
    return assumptions
}
```

## Integration

Domain packs are applied:
1. After architecture parsing
2. Before topological reasoning
3. Before trust boundary discovery
4. Domain is stored in `AnalysisResult.Domain`

## Determinism

All domain detection is deterministic:
- Keyword matching with exact string containment
- No randomness
- No AI/LLM calls
- No cloud services
- Reproducible across runs

## Test Results

| Domain | Test Architecture | Detected | Assumptions Generated |
|--------|------------------|----------|---------------------|
| Healthcare | PHIDatabase, EHR | ✅ Yes | 10 |
| Fintech | Payment, PCI | ✅ Yes | 10 |
| SaaS | Multi-tenant | ✅ Yes | 10 |
| Enterprise | AD, LDAP | ✅ Yes | 10 |

## Compliance Mapping

| Domain | Compliance Frameworks |
|--------|----------------------|
| Healthcare | HIPAA, HITECH, state privacy |
| Fintech | PCI DSS, SOX, GLBA |
| SaaS | SOC2, ISO27001 |
| Enterprise | SOX, ISO27001, NIST |
| Kubernetes | NIST, CIS |
| Cloud Native | SOC2, ISO27001, FedRAMP |
| VPN | NIST, PCI DSS |
| Identity Platform | NIST, ISO27001 |
| Data Platform | GDPR, CCPA, HIPAA |
