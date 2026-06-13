# ASF Benchmark History

## Score Progression

| Version | Score | Verdict |
|---------|-------|---------|
| v2.2 | 3.7 / 10 | NOT_READY |
| v2.8 | 5.6 / 10 | IMPROVED_BUT_NOT_READY |
| v3.0.0-RC2 | 8.3 / 10 | PILOT_READY |

---

## Major Improvements

### v2.2 → v2.8

| Area | Before | After | Delta |
|------|--------|-------|-------|
| Architectural Fidelity | 4/10 | 7/10 | +3 |
| Assumption Discovery | 5/10 | 7/10 | +2 |
| Contradiction Accuracy | 2/10 | 4/10 | +2 |
| Consequence Quality | 3/10 | 6/10 | +3 |
| Trust Chains | 2/10 | 7/10 | +5 |
| Blind Spot Detection | 2/10 | 3/10 | +1 |
| Enterprise Utility | 4/10 | 6/10 | +2 |

**Key fixes:**
- Parser pollution eliminated — YAML comments no longer leak into assumptions
- Fact protection added — insecure architecture (encryption: none) no longer invents TLS/MFA
- Trust chain generation added — 100 chains, 23 cascades, 16 SPOFs
- Risk calibration improved — critical/high risks dominate in insecure architectures
- Assumption categories expanded from 8 to 24
- Recall improved from ~15% to 67.5%
- Precision improved from ~50% to 84.9%

### v2.8 → v3.0.0-RC2

| Area | Before | After | Delta |
|------|--------|-------|-------|
| Parser Pollution | 8/10 | 9/10 | +1 |
| Fact Protection | 8/10 | 9/10 | +1 |
| Verification | 2/10 | 8/10 | +6 |
| Trust Chains | 7/10 | 10/10 | +3 |
| Severity Calibration | 7/10 | 8/10 | +1 |
| Contradiction Precision | 4/10 | 9/10 | +5 |
| Blind Spot / Review | 3/10 | 5/10 | +2 |

**Key fixes (4 workstreams):**

1. **Verification engine rewritten** — `applySecurityControlVerification` produces VERIFIED/PARTIALLY_VERIFIED/CONTRADICTED/UNKNOWN. Maps control names, categories, and concepts to assumptions.
2. **Contradiction precision fixed** — self-comparison guards, storage/backup context isolation, two-phase dedup. 58 contradictions reduced to 6.
3. **Trust chain CLI JSON exposure** — 5 trust-chain fields (chains, cascades, SPOFs, collapse, critical) serialized in CLI JSON output.
4. **SDRI control awareness fixed** — YAML `security_controls` flow into SDRI engine instead of being ignored. CIARE alias matching normalizes variant names.

---

## Remaining Risks

### Non-Blocking (v3.0.0-RC2)

1. **UNKNOWN assumptions in Fixture E** — 17 assumptions remained UNKNOWN due to non-matching assumption categories (Compliance, Privacy) or gap-oriented phrasing. These are engine-generated discovery hypotheses, not false negatives. **Note:** This was later reduced to 0 UNKNOWN after the categoryMap key casing fix. See `benchmarks/benchmark_v3_8.3/summary.md`.

2. **CIARE 0% coverage** — Compliance framework coverage remains 0% across all fixtures. This is an architectural limitation of the framework-mapping model, not a verification blocker.

3. **3 legacy contradictions** — 3 boolean-flag based contradictions in Fixture C are genuine (MFA exemption, plaintext backup, shared admin). Not false positives.

### Previously Blocking (resolved)

| Risk | Resolved In | Evidence |
|------|-------------|----------|
| Parser pollution | v2.8 | Fixture A: 0 comment/text leakage |
| No positive verification | v3.0.0-RC2 | Fixture E: 11 VERIFIED + 4 PARTIAL |
| Contradiction noise (58→4) | v3.0.0-RC2 | Fixture C: 6 total contradictions |
| Self-comparison bugs | v3.0.0-RC2 | 0 self-comparisons across all fixtures |
| Trust chain hidden from CLI | v3.0.0-RC2 | All 5 fields in JSON output |
| SDRI ignoring controls | v3.0.0-RC2 | 10 YAML controls detected |
| categoryMap mixed-case keys | v3.0.0-RC2 | 0 UNKNOWN after fix (was 17) |
