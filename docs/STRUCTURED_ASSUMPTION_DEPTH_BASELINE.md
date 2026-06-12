# Structured Assumption Depth Baseline Report

## Benchmark File
- Path: `testdata/asftest.yaml`
- Content: Healthcare PHI Architecture
- Explicit Assumptions: 30
- Security Control Categories: 8 (authentication, authorization, encryption, logging, backup, network, monitoring, third_party)
- Compliance Targets: HIPAA, SOC2, ISO27001
- Expected Results: minimum_assumptions 25, minimum_critical 3, minimum_high 8, all 6 STRIDE categories

## Baseline Metrics (Before Fixes)
- **Total Assumptions:** 63 (33 analyzer + 30 explicit)
- **Risk Distribution:**
  - Critical: 3
  - High: 4
  - Medium: 10
  - Low: 46
- **Average Confidence:** Low (approx 43% for inferred, 75% for explicit)
- **STRIDE Distribution:** All categories present, but mapping was shallow.
- **Controls:** Generic templates only, enriched by YAML but not used for verification.
- **Compliance:** Shallow listing of framework names.
- **Validation Summary:** Failed (incorrectly reported STRIDE categories as missing due to timing bug).

## Current Metrics (After All Fixes)
- **Total Assumptions:** 48 (18 native + 30 explicit, deduplicated)
- **Risk Distribution:**
  - Critical: 3
  - High: 4
  - Medium: 9
  - Low: 32
- **Average Confidence:** 62% across 48 assumptions (21 high-confidence ≥ 70%)
- **STRIDE Distribution (deep mapped):**
  - Spoofing: 14
  - Tampering: 29
  - Repudiation: 10
  - Information Disclosure: 32
  - Denial of Service: 13
  - Elevation of Privilege: 21
- **Controls:** 15 architecture-specific controls (8 generic + 7 component-specific)
  - Generic: IDENTITY, AUTHENTICATION, AUTHORIZATION, ACCESS, NETWORK, ENCRYPTION, CONFIGURATION, DEPENDENCY, PROCESS, DATABASE, LOGGING, BACKUP, SESSION, THIRD_PARTY, DOCUMENTATION, GOVERNANCE
  - Component-specific: [Auth0] identity-provider controls, [PHIDatabase] database controls, [APIGateway] API gateway controls, [BackupService] storage controls, [KMS] encryption controls, [AuditLog] logging controls, [ThirdPartyAnalytics] third-party controls
- **Compliance:** Deep framework mapping with specific areas:
  - HIPAA: 6 specific safeguard areas (PHI access, encryption, audit, integrity, emergency access, automatic logoff)
  - SOC2: 5 trust services criteria areas
  - ISO27001: 5 Annex A control areas
- **Validation Summary:** Active validation with expected_results
  - 1 violation: expected ≥8 high findings, got 4
  - All criteria met: minimum assumptions (48 ≥ 25), minimum critical (3 ≥ 3), all 6 STRIDE categories present

## Key Improvements
1. **Deduplication:** Fixed by stripping `- ` prefixes in `normalizeText`, eliminating 15 duplicate assumptions.
2. **Risk Scoring:** PHI-related assumptions now properly scored with healthcare keyword boost (+4).
3. **STRIDE Timing:** Fixed by populating StrideDistribution before validation check in `RunAnalysis`.
4. **Compliance Depth:** Added framework-specific area mapping with detailed control references.
5. **Verification Wiring:** Security controls now mark matching assumptions as `PARTIALLY_VERIFIED` with confidence boost to 0.80+.
6. **Dynamic Confidence:** Explicit assumptions get 0.75 base confidence; boosted to 0.85 when supported by declared security controls.
7. **Architecture-Specific Controls:** Added 7 component-specific controls based on actual architecture components (PHIDatabase, Auth0, APIGateway, etc.).
