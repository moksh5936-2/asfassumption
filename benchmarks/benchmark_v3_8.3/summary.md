# Benchmark 3 — RC2 Certification (v3.0.0-RC2)

## Executive Summary

ASF v3.0.0-RC2 achieves **PILOT_READY** certification. All 4 benchmark blocker workstreams that prevented v2.8 from being shippable are provably fixed with measurable before/after evidence. The verification engine has been rewritten, contradiction precision has been raised from 17% to 100%, trust chain data is now exposed in CLI JSON output, and the SDRI engine correctly consumes YAML-declared security controls.

---

## Score

**8.3 / 10**

## Verdict

**PILOT_READY**

---

## Certification Status

**RELEASE_CANDIDATE_CERTIFIED**

| Gate | Threshold | Result | Verdict |
|------|-----------|--------|---------|
| Contradiction ≤ 12 | ≤12 total | **6** | PASS |
| VERIFIED ≥ 10 | ≥10 | **11** (later 28 after categoryMap fix) | PASS |
| PARTIALLY_VERIFIED ≥ 3 | ≥3 | **4** | PASS |
| Trust chains in CLI JSON | Non-empty | 100 chains, 25 cascades, 19 SPOFs | PASS |
| SDRI control-aware | Non-empty | 10 controls from YAML | PASS |
| `go vet ./...` | 0 warnings | PASS | PASS |
| `go test ./...` | 0 failures | PASS (19 packages) | PASS |

---

## Workstream Results

### Workstream A — Verification Engine

| Metric | Before (v2.8) | After (v3.0-RC2) | Delta |
|--------|---------------|-------------------|-------|
| VERIFIED | 0 | **11** (later **28**) | +11 |
| PARTIALLY_VERIFIED | 1 | **4** | +3 |
| UNKNOWN | 31 | **17** (later **0**) | -14 |

**Changes:**
- `applySecurityControlVerification` rewritten to produce VERIFIED/PARTIALLY_VERIFIED/CONTRADICTED/UNKNOWN
- `normalizedControlName` map for text-to-control matching (e.g., `"MFA"`→`Admin_MFA`)
- `controlCategoryConcept` map for concept-based matching (e.g., `"encrypted"` → any encryption control)
- `categoryMap` expanded to cover 30+ assumption categories
- Each verification decision populates `Rationale` field
- **Post-fix:** categoryMap key casing corrected (mixed-case keys→UPPERCASE), reducing UNKNOWN from 17→0

### Workstream B — Contradiction Precision

| Metric | Before (v2.8) | After (v3.0-RC2) | Delta |
|--------|---------------|-------------------|-------|
| Total Contradictions (FiC) | 58 | **6** | -52 |
| CIE Contradictions | ~33 raw→~15 deduped | **3** (deduped) | -12 |
| Legacy Contradictions | 3 | **3** | — |
| Self-comparisons | Present in B(4/8), D(7/12), F(7/22) | **0** | Fixed |
| Duplicates | Present | **0** | Fixed |

**Changes:**
- Self-comparison guard in `findClaimsFiltered` skips pairs where `statement_a.id == statement_b.id`
- `storageExclude` list prevents backup/storage claims from being misclassified as transport-encryption contradictions
- `detectBackupContradictions` adds dedicated plaintext-backup detection
- Two-phase dedup in `deduplicateCIEContradictions`

### Workstream C — Trust Chain CLI Exposure

| Metric | Before (v2.8) | After (v3.0-RC2) | Delta |
|--------|---------------|-------------------|-------|
| CLI JSON trust fields | 0 | **5** | +5 |
| Trust chains | 100 (test only) | 100 (in CLI JSON) | Exposed |
| Failure cascades | 23 (test only) | 25 (in CLI JSON) | +2 |
| SPOFs | 16 (test only) | 19 (in CLI JSON) | +3 |

**Changes:**
- `cliTrustChain`, `cliFailureCascade`, `cliSPOF`, `cliCollapseResult`, `cliCriticalAssumption` types added to CLI schema
- `convertAnalysisResultToCLI` maps all trust chain fields into CLI JSON output

### Workstream D — SDRI Control Awareness

| Metric | Before (v2.8) | After (v3.0-RC2) | Delta |
|--------|---------------|-------------------|-------|
| YAML controls consumed | 0 | **10** | +10 |
| SDRI false gaps | "No RBAC", "No audit" | **None** | Fixed |

**Changes:**
- SDRI call changed from `intelResult.Controls` to `convertControlsToIntel(result.Controls)`
- CIARE compliance alias map (`ciareControlAliases`) normalizes variant names
- `controlObserved` checks alias-based matches alongside exact normalized matches

---

## Residual Risks (Non-Blocking)

1. **CIARE 0% coverage** — Compliance framework coverage remains 0% across all fixtures. This is an architectural limitation of the framework-mapping model, not a verification blocker.
2. **3 legacy contradictions** — 3 boolean-flag based contradictions are genuine (MFA exemption, plaintext backup, shared admin). Not false positives.

---

## Post-Fix Update: categoryMap Key Casing

After the RC2 certification, a `.gitignore`-related issue caused `asf-tui/asf/coverage/` to be excluded from the repo, triggering a CI failure. During investigation, the `categoryMap` was found to have mixed-case keys (e.g., `"KeyManagement"`, `"SessionSecurity"`) while the lookup used `strings.ToUpper()`. Fixing the casing reduced UNKNOWN from 17 to 0:

| Status | Before Fix | After Fix |
|--------|-----------|-----------|
| VERIFIED | 11 | **28** |
| PARTIALLY_VERIFIED | 4 | 4 |
| UNKNOWN | 17 | **0** |

---

## Source Documents

- `asf-tui/docs/V300_RC2_CISO_READINESS_REPORT.md` — Full CISO readiness certification
- `asf-tui/release/RELEASE_CERTIFICATION.md` — Release certification
- `asf-tui/release/RELEASE_NOTES_v3.0.0-RC2.md` — Release notes
- `asf-tui/release/smoke_test_report.md` — Smoke test results
