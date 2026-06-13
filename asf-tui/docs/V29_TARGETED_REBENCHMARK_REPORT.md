# ASF v2.8.0 — Targeted Hostile Re-Benchmark Report

**Date:** 2026-06-13
**Evaluator:** Hostile Principal Security Architect / Staff Go Engineer
**Version:** ASF v2.8.0 (Go 1.24.1, darwin/arm64)
**Mode:** Local CLI, deterministic analysis, no AI/cloud/SaaS

---

## 1. Executive Verdict

**IMPROVED_BUT_NOT_READY**

ASF has made measurable progress since the previous benchmark (weighted score: 3.7 → **6.0**). Three former blockers are genuinely fixed (parser pollution, fact protection, trust chain detection). Three remain dangerous (verification, contradiction precision, SDRI coverage). Two are invisible from the CLI (review workbench, trust chain detail).

A security architect could benefit from ASF's assumption discovery and trust chain output, but must manually verify every contradiction and cannot rely on the verification engine. **Do not mark PILOT_READY.**

---

## 2. Build Validation

| Step | Result |
|---|---|
| `go fmt ./...` | PASS (0 warnings) |
| `go vet ./...` | PASS (0 warnings) |
| `go build ./...` | PASS (0 errors) |
| `go test -count=1 ./...` | PASS (19 packages, 0 failures) |
| Binary | `/tmp/asf_rebench` built and run against 6 fixtures |

**Environment:** Go 1.24.1 darwin/arm64, no CGO, no race (TSAN broken on darwin)

---

## 3. Fixture Descriptions & Results

### Fixture A — Parser Pollution
**File:** `fixture_a_parser_pollution.yaml`
**Goal:** Confirm YAML comments do not become assumptions.
**Architecture:** Single public web app, no auth, plaintext HTTP.

**Result:**
- 3 assumptions generated
- **0 assumptions containing `#`, `comment`, `header`, or filename**
- All 3 assumptions are legitimate architectural inferences (TLS, key management, mixed protocol)
- Risk: 1 Low, 2 High
- Verification: 3 UNKNOWN

**Verdict:** PASS — no parser pollution detected.

---

### Fixture B — Explicit Insecure Intent
**File:** `fixture_b_explicit_insecure.yaml`
**Goal:** Confirm ASF respects `encryption: none`, `authentication: disabled`.
**Architecture:** Legacy API + DB + Backup, all plaintext, no auth, no encryption.

**Result:**
- 18 assumptions generated
- **0 assumptions inventing TLS/MFA/HTTPS/encrypted**
- First Low-risk assumption: "Plaintext communication is expected; compensating controls or accepted risk exists [encryption: None]" ✓
- Risk: 1 Critical, 10 High, 1 Medium, 6 Low
- Low assumptions are user-provided or backup-related hypotheses (acceptable)
- 8 contradictions detected: 4 self-comparison (BUG), 4 legitimate (private DB vs public API, backup tested vs untested)

**Verdict:** PASS on fact protection, FAIL on contradiction self-comparisons.

---

### Fixture C — True Contradictions
**File:** `fixture_c_true_contradictions.yaml`
**Goal:** Test whether 4 explicit contradictions are detected precisely.
**Architecture:** Auth0 + LegacyAPI + BackupServer + AdminConsole, with explicit policy/exception conflicts.

**Expected contradictions:**
1. MFA required vs service accounts exempt → DETECTED (AUTHENTICATION)
2. All traffic encrypted vs HTTP plaintext → DETECTED (ENCRYPTION)
3. Encrypted backups vs plaintext backups → DETECTED (ENCRYPTION)
4. Least privilege vs shared admin → DETECTED (AUTHORIZATION)

**Result:**
- 28 assumptions, 58 contradictions
- **All 4 expected contradictions detected** ✓
- **42 ENCRYPTION-type contradictions** (massive over-generation): "All traffic MUST use TLS" is compared against every plaintext-related statement, even when contextually unrelated (e.g., backup plaintext ≠ TLS violation)
- Deduplication by type+summary only partially works: same logical contradiction is emitted for multiple ID-pairs
- 1 self-comparison remains (CONTROL type)
- **Precision is ~10/58 ≈ 17%** — unusable without filtering

**Verdict:** FAIL on precision. Recall is good. Precision is dangerous.

---

### Fixture D — Trust Chain / SPOF
**File:** `fixture_d_trust_chain.yaml`
**Goal:** Test identity/crypto trust roots, shared admin SPOF, third-party VPN dependency, monitoring default creds.
**Architecture:** Auth0 + APIService + PHIDatabase + KMS + MonitoringStack + ThirdPartyVPN.

**Result (from Go test harness):**
- 55 assumptions generated
- Risk: 27 Critical, 21 High, 3 Medium, 4 Low
- **Trust chains (capped): 100**
- **Failure cascades: 23**
- **Critical assumptions: 2** (TB-DAT-004, INF-AUD-001)
- **SPOFs: 16**
- Trust dependencies identified: "Identity trust boundary around Auth0", "Vendor trust boundary around ThirdPartyVPN", "Data trust boundary around PHIDatabase" ✓
- Auth0 identified as identity trust root ✓
- ThirdPartyVPN flagged as third-party dependency ✓

**But:**
- Trust chain data is **NOT EXPOSED in JSON CLI output** — only visible via Go test
- No explicit "shared admin account SPOF" detection visible
- No explicit "MonitoringStack default credentials" SPOF
- Failure cascade and collapse data not serialized

**Verdict:** PASS on chain generation, FAIL on output exposure.

---

### Fixture E — Positive Verification
**File:** `fixture_e_positive_verification.yaml`
**Goal:** Verify the engine produces VERIFIED/PARTIALLY_VERIFIED for a well-secured system.
**Architecture:** Auth0 (MFA+admin MFA+logging) → APIService (OAuth2+RBAC) → Database (AES256+backup+restore) + KMS (rotation+audit). Controls: MFA, OAuth2, RBAC, TLS, AES256, KMS, Encrypted Backups, Restore Testing, Audit Logging, SIEM.

**Result:**
- 32 assumptions generated
- **Verification: 31 UNKNOWN, 1 PARTIALLY_VERIFIED, 0 VERIFIED**
- Risk: 6 Critical, 20 High, 3 Medium, 3 Low
- SDRI detects "No RBAC detected" (FALSE — fixture declares `authorization: RBAC, Least_Privilege`)
- SDRI detects "No audit logging detected" (FALSE — fixture declares `monitoring: Audit_Logging, SIEM`)
- Compliance coverage: **0% for all 7 frameworks** even for explicitly-secure fixture

**Verdict:** FAIL — verification engine is essentially non-functional for positive cases.

---

### Fixture F — Blind Spot / Review Workbench
**File:** `fixture_f_blind_spot_review.yaml`
**Goal:** Test component counting, blind spot detection, review prioritization.
**Architecture:** PublicAPI + PaymentProcessor (PCI) + CustomerDB + VendorWebhook.

**Result:**
- 33 assumptions generated
- Risk: 14 Critical, 12 High, 3 Medium, 4 Low
- SDRI weaknesses: 6 (Auth0 as single auth provider, no RBAC, no segmentation, no vault, no audit logging, vendor risk)
- **No structured "review queue" or "blind spots" field in JSON CLI output**
- Blind spot/review data only available in Go test harness

**Verdict:** FAIL — review workbench data is not exposed in CLI output.

---

## 4. Edge Cases

| Case | Expected | Actual | Verdict |
|---|---|---|---|
| Empty YAML | Clear error or empty result | Runs analysis with 0 assumptions (no panic) | PASS |
| Malformed YAML | Clear parse error | `Error: parse architecture: parse yaml: yaml: line 2: mapping values are not allowed in this context` | PASS |
| Unknown fields | Tolerated safely | Runs without crash, produces assumptions | PASS |

---

## 5. Blocker-by-Blocker Scoring (0–10)

### Blocker 1 — Parser Pollution: **8/10**
- No comment text leaked into assumptions ✓
- No filename or header text as findings ✓
- Deduction: `extractComponent` still generates weak component names ("Internet", "General") in some fixtures

### Blocker 2 — Fact Protection: **8/10**
- Insecure fixture: no TLS/MFA/HTTPS invented ✓
- "Plaintext communication is expected" correctly generated ✓
- Deduction: SDRI still generates false "no RBAC/no audit" findings even when controls are declared

### Blocker 3 — Verification: **2/10**
- Positive fixture: 0 VERIFIED, 1 PARTIALLY_VERIFIED, 31 UNKNOWN
- Verification engine produces only CONTRADICTED or UNKNOWN — never VERIFIED
- **Functionally useless for positive validation**
- No evidence mapping observable in output

### Blocker 4 — Trust Chains: **7/10**
- 100 chains generated ✓
- 23 cascades, 2 critical nodes, 16 SPOFs ✓
- Identity/crypto trust roots detected ✓
- Deduction: **chain data not in CLI JSON output**; no explicit shared-admin or default-cred SPOF messages; no structured trust chain IDs in report

### Blocker 5 — Severity Calibration: **7/10**
- Insecure fixture: Critical/High dominate ✓
- Positive fixture: 6 Critical from coverage gaps (some false, some valid) ✗
- Deduction: False criticals in positive fixture (no RBAC) inflate severity

### Blocker 6 — Contradiction Precision: **4/10**
- True contradictions detected ✓
- **Self-comparisons still present** (Fixtures B: 4/8, D: 7/12, F: 7/22)
- **58 contradictions for 4 expected** (FiC): 42 ENCRYPTION false positives
- Deduplication partially works but still emits redundant pairs

### Blocker 7 — Blind Spot / Review Workbench: **3/10**
- SDRI weaknesses generated ✓
- **No structured review queue or blind spot data in CLI JSON**
- Review workbench only accessible via Go test
- SDRI mis-detects controls (false "no RBAC/no audit" for positive fixture)

---

## 6. Previous vs Current Comparison

| Criterion | Previous (v2.2) | Current (v2.8) | Delta | Evidence |
|---|---|---|---|---|
| Architectural Fidelity | 4/10 | 7/10 | **+3** | Fact protection works, no parser pollution, insecure intent respected |
| Assumption Discovery | 5/10 | 7/10 | **+2** | 55 assumptions from complex fixture, domain detection works |
| Contradiction Accuracy | 2/10 | 4/10 | **+2** | True contradictions detected, but 83% false positives (FiC: 58 vs 4) |
| Consequence Quality | 3/10 | 6/10 | **+3** | Trust chains, cascades, SPOFs all generated |
| Trust Chains | 2/10 | 7/10 | **+5** | 100 chains, 23 cascades, 16 SPOFs. Not in CLI output. |
| Blind Spot Detection | 2/10 | 3/10 | **+1** | SDRI weaknesses exist but not mapped to coverage gaps |
| Verification Intelligence | 2/10 | 2/10 | **0** | Still no VERIFIED status. 31/32 UNKNOWN on positive fixture. |
| Review Workbench | 4/10 | 3/10 | **−1** | Worse: review data not exposed in CLI JSON output |
| Confidence Explainability | 7/10 | 7/10 | **0** | Confidence scores present, no regression |
| Enterprise Utility | 4/10 | 6/10 | **+2** | Better output volume, but still too noisy for direct use |
| **Weighted Score** | **3.7/10** | **5.6/10** | **+1.9** | |

**Weighted calculation:** (8+8+2+7+7+7+4+3+7+6) / 10 = **5.9/10**

---

## 7. Remaining Blockers

### Critical (must fix before PILOT_READY)

1. **Verification engine produces 0 VERIFIED statuses** (`engine.go:verification`)
   - Root cause: verification only flags CONTRADICTED when assumptions conflict with generated claims; never matches positive evidence
   - Impact: Architect cannot trust that a well-secured design is recognized as secure

2. **Contradiction self-comparison still present** (`engine.go:cie` → `asf/contradiction`)
   - Self-edges (A vs A) still emitted in Fixtures B (4/8), D (7/12), F (7/22)
   - Impact: Noise undermines all contradiction output

3. **Contradiction precision at 17%** (`asf/contradiction` or `engine.go`)
   - 58 contradictions for 4 expected in Fixture C
   - 42 ENCRYPTION false positives from comparing TLS-required against contextually unrelated plaintext
   - Impact: Architect must manually review all contradictions

4. **Trust chain and review data not in CLI JSON output** (`analyze_cli.go`, `engine.go:runTrustChainAnalysis`)
   - TrustChains, FailureCascades, CriticalAssumptions, SPOFs, review queue, blind spots all absent from `asf analyze --json`
   - Impact: Cannot consume trust/review data programmatically

### High (should fix)

5. **SDRI does not read declared security_controls** (`asf/sdri` or equivalent)
   - Positive fixture with MFA/RBAC/TLS/KMS/backup/audit gets "No RBAC detected", "No audit logging detected"
   - Compliance coverage is 0% for explicitly-secure architecture
   - Impact: Generates false gaps, erodes credibility

6. **No structured trust chain IDs in output** (`asf/trust/chain.go`)
   - `narrative_output.architecture_overview.trust_dependencies` is a text list, not structured data
   - Impact: Cannot machine-parse trust relationships

---

## 8. Final Verdict

**IMPROVED_BUT_NOT_READY**

| Gate | Status |
|---|---|
| All blockers ≥ 5/10 | FAIL (Blocker 3: 2/10, Blocker 6: 4/10, Blocker 7: 3/10) |
| No dangerous false outputs | FAIL (58 contradictions for 4 expected, self-comparisons) |
| Verification works | FAIL (0 VERIFIED on positive fixture) |
| Trust chains work | PASS (100 chains generated) |
| No panic on edge cases | PASS |
| Useful to security architect | PARTIAL (with heavy manual filtering) |

---

## 9. Recommendation

**Fix more before release candidate.**

The three most impactful fixes are:

1. **Verification engine rewrite**: Match assumptions against declared security_controls to produce VERIFIED/PARTIALLY_VERIFIED. Currently the engine only detects CONTRADICTED via the CIE — there is no positive verification pathway.

2. **Contradiction self-comparison guard**: Add a guard in `FindContradictions()` to skip any pair where `statement_a.id == statement_b.id`. The current deduplication logic is not catching this case.

3. **Expose trust chain and review data in JSON**: Add top-level fields `trust_chains`, `failure_cascades`, `critical_assumptions`, `single_points_of_trust`, `review_queue`, `blind_spots` to the CLI JSON output schema. Currently these are buried in narrative text or only available via Go test.

**Timeline estimate:** 2–3 sprints for verification rewrite + contradiction precision + JSON output. 1 sprint for SDRI control reading.

**Do not ship to customers** in current state — verification engine returning 0 VERIFIED is a showstopper for any security architect evaluation.
