# ASF Assumption Engine V3

## Architecture

The V3 engine introduces an **Intelligence Layer** that sits between the native analyzer and the output formatter. This layer is responsible for:

1. **Domain Detection** - Identifying the domain from architecture keywords
2. **Topological Reasoning** - Inferring hidden assumptions from component relationships
3. **Trust Boundary Discovery** - Identifying security boundaries
4. **Contradiction Detection** - Finding conflicting statements
5. **Quality Scoring** - Ranking assumptions by value
6. **Explainability** - Generating human-readable WHY explanations

## Engine Flow

```
Architecture Input
    ↓
ParseArchitecture() → ArchDescription
    ↓
Native Analyzer (existing pipeline)
    ↓
buildResult() → AnalysisResult with basic assumptions
    ↓
Intelligence Engine
    ├─ Detect Domain
    ├─ Apply Domain Pack
    ├─ Topological Reasoning
    ├─ Trust Boundary Discovery
    ├─ Explainability
    ├─ Quality Scoring
    └─ Contradiction Detection
    ↓
Merge Results (deduplication)
    ↓
Final AnalysisResult with contradictions, boundaries, domain, quality scores
    ↓
Export (JSON, Markdown, CSV, HTML, PDF)
```

## Intelligence Engine Components

### Taxonomy Engine
- 40+ categories with keyword patterns
- Each category has: keywords, risk mappings, verification rules, explainability templates
- Classification is deterministic regex-based

### Reasoning Engine
- 8 topological rule sets:
  1. Database + PHI → key management, audit logging, object-level auth
  2. Identity Provider (Auth0) → MFA, session security, token validation
  3. API Gateway → rate limiting, auth validation, logging
  4. KMS → key rotation, access restriction, backup
  5. Backup Service → encryption, testing, geographic distribution
  6. Third Party → vendor risk, equivalent controls, data minimization
  7. Admin Console → MFA, break-glass, audit
  8. Audit Log → immutability, retention, tamper detection

### Domain Engine
- 10 domain packs with auto-detection
- Healthcare: PHI, HIPAA, auditability, patient privacy, retention
- Fintech: PCI DSS, SOX, fraud detection, transaction integrity
- SaaS: Multi-tenancy, tenant isolation, data segregation
- Enterprise: Identity lifecycle, access reviews, compliance
- Kubernetes: Container security, RBAC, network policies
- Cloud Native: IAM, encryption, logging, monitoring, resilience
- VPN: Tunnel security, endpoint validation, certificate management
- Identity Platform: SSO, MFA, session management, federation
- Data Platform: Data governance, lineage, retention, encryption

### Contradiction Engine
- 8 contradiction rules:
  1. MFA_ENFORCED_WITH_EXEMPTION
  2. ENCRYPTED_WITH_PLAINTEXT_BACKUP
  3. LEAST_PRIVILEGE_WITH_SHARED_ADMIN
  4. PRIVATE_SUBNET_WITH_INTERNET_ACCESS
  5. IMMUTABLE_AUDIT_WITH_LOG_DELETION
  6. TLS_REQUIRED_HTTP_ALLOWED
  7. ENCRYPTION_WITHOUT_KEY_MANAGEMENT
  8. SESSION_WITHOUT_ROTATION

### Trust Boundary Engine
- 8 boundary types:
  1. Internet
  2. Identity
  3. Tenant
  4. Vendor
  5. Network
  6. Administrative
  7. Data
  8. Cloud

### Quality Engine
- 6 dimensions:
  - Hiddenness (0-1): How non-obvious the assumption is
  - Impact (0-1): Severity if violated
  - Novelty (0-1): How unique the assumption is
  - Architectural Relevance (0-1): How well it maps to components
  - Risk (0-1): Risk level of the assumption
  - Confidence (0-1): Confidence in the assumption
- Composite: 0.20*Hiddenness + 0.20*Impact + 0.20*Novelty + 0.15*Relevance + 0.15*Risk + 0.10*Confidence
- Link-encryption assumptions penalized to 0.2
- Domain-specific inferred assumptions boosted to 0.9

### Explainability Engine
- Every assumption gets a WHY explanation
- Evidence from architecture components
- Missing controls identified
- Architecture context provided
- Confidence and category justification

## Integration Points

### Native Engine → Intelligence Engine
```go
// Convert native assumptions to intelligence format
intelAssumptions := convertAssumptionsToIntel(result.Assumptions)

// Run intelligence engine
ie := intelligence.NewIntelligenceEngine()
intelResult := ie.RunWithExistingAssumptions(intelArch, intelAssumptions)

// Convert back and merge
result.Assumptions = mergeAssumptions(result.Assumptions, convertIntelAssumptions(intelResult.Assumptions))
result.Contradictions = convertIntelContradictions(intelResult.Contradictions)
result.TrustBoundaries = convertIntelTrustBoundaries(intelResult.TrustBoundaries)
result.Domain = intelResult.Domain
```

### Backward Compatibility
- All existing fields preserved
- New fields are optional (omitempty)
- Existing CLI commands work unchanged
- Existing TUI unaffected
- Existing exports unaffected

## Determinism

All reasoning is rule-based:
- Regex patterns for classification
- Keyword matching for domain detection
- Topology rules for inference
- String comparison for contradiction detection
- No randomness, no AI, no cloud services

## Performance

- Intelligence engine adds ~10-20ms to analysis
- No external dependencies
- All processing in-memory
- No network calls

## Test Coverage

- Benchmark harness: recall, precision, F1 score
- Contradiction detection: 2+ contradictions detected
- Trust boundary discovery: 3+ boundaries detected
- Domain detection: Healthcare detected
- Quality scoring: assumptions ranked correctly
- No regressions in existing tests
