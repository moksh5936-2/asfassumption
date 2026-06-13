# ASF v3.0.0-RC2 — Release Notes

**Release Candidate 2 for hostile independent benchmark evaluation.**

---

## Overview

ASF (Architecture Security Framework) v3.0.0-RC2 is a release candidate that fixes all 4 blocker workstreams identified in the v2.8.0 hostile benchmark evaluation. The weighted score improved from **5.6/10 (IMPROVED_BUT_NOT_READY)** to **8.3/10 (PILOT_READY)**.

This is a release-engineering-only candidate. No new features, no architecture changes, no roadmap expansion.

---

## Major Improvements

### Verification Engine (Blocker Score: 2/10 → 8/10)
- Rewritten `applySecurityControlVerification` to produce VERIFIED, PARTIALLY_VERIFIED, CONTRADICTED, and UNKNOWN statuses
- `normalizedControlName` map for text-to-control matching (e.g., "MFA" → "Admin_MFA")
- `controlCategoryConcept` map for concept-based matching (e.g., "encrypted" → any encryption control)
- Expanded `categoryMap` covering 30+ assumption categories
- Result: 11 VERIFIED + 4 PARTIALLY_VERIFIED on the positive verification fixture (was 0 + 1)

### Contradiction Precision (Blocker Score: 4/10 → 9/10)
- Self-comparison guards in all cross-product loops
- Storage/backup context isolation preventing false TLS/transport contradictions
- Two-phase deduplication normalizing text-based pairs then deduplicating by type+summary
- Result: 6 total contradictions on fixture C (was 58), 100% precision (was ~17%)

### Trust Chain CLI Exposure (Blocker Score: 7/10 → 10/10)
- 5 trust-chain fields added to CLI JSON output: `trust_chains`, `failure_cascades`, `critical_assumptions`, `single_points_of_trust_failure`, `trust_collapse_results`
- New types: `cliTrustChain`, `cliFailureCascade`, `cliCascadeResult`, `cliCriticalAssumption`, `cliSinglePointOfTrustFailure`, `cliTrustCollapseResult`
- Result: 100 chains, 25 failure cascades, 19 SPOFs in CLI JSON output (was 0)

### SDRI Control Awareness (New: 10/10)
- SDRI pipeline now receives YAML-enriched controls via `convertControlsToIntel`
- CIARE alias map normalizes variant control names
- Result: 10 controls from YAML flow into SDRI engine

---

## Benchmark Progress

| Metric | v2.8.0 | v3.0.0-RC2 | Delta |
|--------|--------|-------------|-------|
| Weighted Score | 5.6/10 | **8.3/10** | +2.7 |
| Verdict | IMPROVED_BUT_NOT_READY | **PILOT_READY** | |
| Contradiction Precision | 58 (17%) | **6 (100%)** | −52 |
| VERIFIED Count | 0 | **11** | +11 |
| PARTIALLY_VERIFIED | 1 | **4** | +3 |
| Trust Chains in CLI JSON | 0 fields | **5 fields** | +5 |
| SDRI Controls from YAML | 0 | **10** | +10 |

### Fixture Results

| Fixture | Before | After | Status |
|---------|--------|-------|--------|
| A — Parser Pollution | PASS | PASS | ✅ |
| B — Explicit Insecure | PASS | PASS | ✅ |
| C — True Contradictions | FAIL (58) | PASS (6) | ✅ |
| D — Trust Chain | FAIL (no CLI) | PASS (100 chains) | ✅ |
| E — Positive Verification | FAIL (0 VERIFIED) | PASS (11 VERIFIED) | ✅ |
| F — Blind Spot Review | PASS | PASS | ✅ |

---

## Known Limitations

1. **17 UNKNOWN assumptions** remain on the positive verification fixture — generated assumptions with categories like `Compliance` and `Privacy` that don't match existing control maps
2. **CIARE compliance coverage** remains 0% across all fixtures — the framework-mapping model is an architectural limitation outside this RC scope
3. **Version comparison glitch** — `--version-check` may show an incorrect "newer version available" message due to non-semver comparison
4. **RC build** intended for hostile third-party evaluation, not production deployment

---

## Build Artifacts

| Platform | Binary | Size |
|----------|--------|------|
| macOS ARM64 | `asf-darwin-arm64` | 16.3 MB |
| macOS AMD64 | `asf-darwin-amd64` | 17.4 MB |
| Linux AMD64 | `asf-linux-amd64` | 17.3 MB |
| Linux ARM64 | `asf-linux-arm64` | 16.2 MB |
| Windows AMD64 | `asf-windows-amd64.exe` | 17.8 MB |

All binaries built with CGO_ENABLED=0, -trimpath, -buildvcs=true.

---

## Upgrade Notes

- No database migrations required
- No configuration changes required
- Drop-in replacement for v2.8.0
- JSON output schema extended with trust-chain fields (backward compatible)
