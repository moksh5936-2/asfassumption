# Structured Assumption Depth Root Cause Analysis

## 1. Explicit Assumption Processing
**Observation:** 30 explicit assumptions in YAML produce a mix of low-risk findings and duplicates.
**Root Cause:**
- The `analyzer` (extractor) identifies some explicit assumptions as claims based on declarative patterns but preserves the `- ` bullet point.
- `processExplicitAssumptions` adds the raw explicit assumptions.
- `normalizeText` did not strip the `- ` prefix, so the dedup logic treated `- MFA is enforced...` and `MFA is enforced...` as different assumptions.
**Fix:** Enhanced `normalizeText` regex to strip leading bullet markers (`- `, `* `, `• `). Added `mergeSourceMetadata` to merge source metadata when duplicates are detected.

## 2. Risk Scoring Sensitivity
**Observation:** PHI-related assumptions were not flagged as High/Critical.
**Root Cause:** `assessExplicitRisk` had insufficient weights for healthcare keywords (+3) and a threshold for "High" that was too high relative to the boosts. PHI-specific context was not strongly linked to risk levels.
**Fix:** Added dedicated PHI keyword boost (+4 score) for terms like "phi", "health", "hipaa", "patient", "medical", "ehr". This ensures PHI assumptions score at least Medium and often High/Critical.

## 3. STRIDE Validation Failure
**Observation:** `buildValidationSummary` reports missing STRIDE categories even when present in the final output.
**Root Cause:** `result.StrideDistribution` is populated *after* `buildResult()` returns, but `buildValidationSummary` was called *inside* `buildResult()` before StrideDistribution was populated.
**Fix:** Moved STRIDE distribution mapping and validation summary generation to `RunAnalysis()` after `buildResult()` completes, ensuring distribution is populated before validation.

## 4. Shallow Compliance Output
**Observation:** Compliance section only lists the framework names.
**Root Cause:** `buildComplianceOutput` lacks a mapping between framework names (e.g., "HIPAA") and the specific security domains/controls that ASF should evaluate.
**Fix:** Added `complianceFrameworkDetails` map with specific areas and control references for HIPAA, SOC2, ISO27001, PCI DSS, GDPR, and FedRAMP.

## 5. Verification Context Gap
**Observation:** `security_controls` in YAML are displayed as controls but don't help verify assumptions.
**Root Cause:** The verification engine in `asf/verification/engine.go` does not receive the `security_controls` map; it only looks at the provided evidence files.
**Fix:** Added `applySecurityControlVerification` in `engine.go` that checks each assumption against declared security controls. Matching assumptions are marked `PARTIALLY_VERIFIED` with confidence boosted to 0.80+.

## 6. Control Genericness
**Observation:** Recommended controls are mostly generic templates.
**Root Cause:** `generateControls` relies on a static template map. While `enhanceControlsWithSecurityControls` appends YAML controls, it doesn't generate new, architecture-specific recommendations based on the gaps.
**Fix:** Added `generateArchitectureSpecificControls` that inspects actual architecture components and generates component-specific controls (e.g., `[PHIDatabase] database-specific encryption and access logging`, `[Auth0] identity-provider MFA enforcement`).

## 7. Confidence Stagnation for Explicit Assumptions
**Observation:** All explicit assumptions receive a flat 0.75 confidence.
**Root Cause:** `processExplicitAssumptions` hard-codes confidence to 0.75 without considering whether the architecture definition includes supporting security controls.
**Fix:** Added `computeExplicitConfidence` that maintains 0.75 base but scans security controls for matching keywords. If a declared control supports the assumption, confidence is boosted to 0.85.

## 8. Missing Validation Summary in `buildResult`
**Observation:** `TestBuildResultWithExplicitAssumptions` failed because validation summary was not present in `result.Summary`.
**Root Cause:** `buildValidationSummary` was only called in `RunAnalysis`, not in `buildResult`. Direct calls to `buildResult` (e.g., in tests) did not generate validation summaries.
**Fix:** Added validation summary generation at the end of `buildResult` when `ExpectedResults` is defined.
