# Benchmark 2 — Targeted Re-Benchmark (v2.8)

## Executive Summary

ASF v2.8 showed measurable progress from the v2.2 baseline. Three former blockers were genuinely fixed: parser pollution, fact protection, and trust chain detection. However, three critical blockers remained: the verification engine was still non-functional (0 VERIFIED), contradiction precision was dangerously low (58 contradictions for 4 expected), and trust chain/review data was hidden from CLI JSON output.

A security architect could benefit from ASF's assumption discovery and trust chain output, but would need to manually verify every contradiction and could not rely on the verification engine.

---

## Score

**5.6 / 10**

## Verdict

**IMPROVED_BUT_NOT_READY**

---

## Major Findings

### Top Strengths

1. **Parser pollution fixed** — Fixture A: 0 assumptions containing comment/header/filename text
2. **Fact protection works** — Fixture B: `encryption: none` respected, no TLS/MFA invented
3. **Trust chain generation** — 100 chains, 23 cascades, 16 SPOFs detected
4. **Edge case handling** — Empty YAML, malformed YAML, unknown fields all handled gracefully

### Top Weaknesses

1. **Verification engine non-functional** — Fixture E: 0 VERIFIED, 1 PARTIALLY_VERIFIED, 31 UNKNOWN
2. **Contradiction precision at 17%** — Fixture C: 58 contradictions for 4 expected
3. **Self-comparison bugs** — Fixtures B (4/8), D (7/12), F (7/22) contained self-edges
4. **Trust chain data hidden from CLI** — All trust chain fields absent from `asf analyze --json`
5. **SDRI ignoring declared controls** — Fixture E with MFA/RBAC/TLS gets "No RBAC detected"
6. **Review workbench hidden from CLI** — No blind spot or review queue in JSON output

### Blocking Issues

| Issue | Score | Impact |
|-------|-------|--------|
| Verification non-functional | 2/10 | Architect cannot verify secure designs |
| Contradiction precision | 4/10 | 58 contradictions for 4 expected; 83% false positives |
| Trust chain CLI exposure | 7/10 | Trust data generated but inaccessible via CLI |
| SDRI control ignorance | 3/10 | Generates false gap findings |

### Actions Taken (Leading to v3.0.0-RC2)

1. **Verification engine rewrite** — `applySecurityControlVerification` rewritten with `normalizedControlName`, `controlCategoryConcept`, `categoryMap`
2. **Contradiction self-comparison guard** — Edge filtering in `findClaimsFiltered`
3. **Storage/backup context isolation** — `storageExclude` prevents backup claims from being misclassified
4. **Two-phase dedup** — CIE contradictions deduplicated in `deduplicateCIEContradictions`
5. **Trust chain JSON exposure** — 5 fields added to CLI output schema
6. **SDRI control pipeline** — `convertControlsToIntel` feeds YAML controls into SDRI
7. **CIARE alias mapping** — `ciareControlAliases` normalizes variant control names

### Outcome

Verdict: **Do not ship to customers.** Verification engine returning 0 VERIFIED is a showstopper for any security architect evaluation.

---

## Source Documents

- `asf-tui/docs/V29_TARGETED_REBENCHMARK_REPORT.md` — Full re-benchmark report
- `asf-tui/docs/V29_CORE_HARDENING_REPORT.md` — Hardening analysis
- `asf-tui/testdata/v29_rebenchmark/` — Fixture files (Fixtures A-F)
