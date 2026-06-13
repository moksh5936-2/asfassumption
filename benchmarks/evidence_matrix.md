# ASF Benchmark Evidence Matrix

## Finding-to-Fix Mapping

| Finding | First Observed | Fix Version | Evidence Location |
|---------|---------------|-------------|-------------------|
| Parser pollution | v2.2 (3.7) | v2.8 | `docs/BENCHMARK_IMPROVEMENT_REPORT.md`, `asf-tui/docs/V29_TARGETED_REBENCHMARK_REPORT.md` Fixture A |
| Fact protection (invented TLS/MFA) | v2.2 (3.7) | v2.8 | `docs/BENCHMARK_IMPROVEMENT_REPORT.md`, Fixture B |
| Low assumption recall (~15%) | v2.2 (3.7) | v2.8 | `docs/BENCHMARK_IMPROVEMENT_REPORT.md` |
| Only 8 assumption categories | v2.2 (3.7) | v2.8 | `docs/BENCHMARK_IMPROVEMENT_REPORT.md` |
| No trust chain generation | v2.2 (3.7) | v2.8 | `docs/BENCHMARK_IMPROVEMENT_REPORT.md` |
| No contradiction detection | v2.2 (3.7) | v2.8 | `docs/BENCHMARK_IMPROVEMENT_REPORT.md` |
| Verification engine 0 VERIFIED | v2.2 (3.7) | **v3.0.0-RC2** | `asf-tui/docs/V300_RC2_CISO_READINESS_REPORT.md` WS A |
| Contradiction noise (58→4, 83% FP) | v2.8 (5.6) | **v3.0.0-RC2** | `asf-tui/docs/V300_RC2_CISO_READINESS_REPORT.md` WS B |
| Contradiction self-comparison bugs | v2.8 (5.6) | **v3.0.0-RC2** | `intelligence/contradiction_intelligence.go` |
| Storage/backup context misclassification | v2.8 (5.6) | **v3.0.0-RC2** | `intelligence/contradiction_intelligence.go` — `storageExclude` |
| Trust chain data hidden from CLI | v2.8 (5.6) | **v3.0.0-RC2** | `analyze_cli.go` — trust chain fields |
| SDRI ignoring YAML controls | v2.8 (5.6) | **v3.0.0-RC2** | `engine.go` — `convertControlsToIntel` |
| CIARE variant name mismatch | v2.8 (5.6) | **v3.0.0-RC2** | `intelligence/ciare.go` — `ciareControlAliases` |
| categoryMap key casing (UNKNOWN→0) | v3.0.0-RC2 (8.3) | **v3.0.0-RC2** | `engine.go` — `categoryMap` uppercase fix |
| `asf-tui/asf/coverage/` gitignored | v3.0.0-RC2 (8.3) | **v3.0.0-RC2** | `.gitignore` — negation rule |

## Fix Severity

| Finding | Severity | Customer Impact |
|---------|----------|-----------------|
| Verification engine 0 VERIFIED | **Critical** | Architect cannot trust positive verification |
| Contradiction noise (83% FP) | **Critical** | Manual review of all contradictions required |
| Storage/backup misclassification | **High** | False encryption contradictions |
| Trust chain hidden from CLI | **High** | Cannot consume trust data programmatically |
| SDRI ignoring controls | **High** | False gap findings erode credibility |
| CIARE variant name mismatch | **Medium** | Some compliance mappings missed |
| categoryMap key casing | **Medium** | 17/32 assumptions stayed UNKNOWN |
| Parser pollution | **Medium** | Assumptions contain irrelevant text |
| coverage/ gitignored | **Low** | CI pipeline failure only |

## Workstream Clusters

| Workstream | Findings | Fixes | Evidence |
|------------|----------|-------|----------|
| **A — Verification** | 0 VERIFIED, categoryMap casing, concept matching | `applySecurityControlVerification`, `normalizedControlName`, `controlCategoryConcept`, `categoryMap` | Fixture E: 0→11→28 VERIFIED |
| **B — Contradiction** | 58→4, self-comparisons, storage context | Self-comparison guards, `storageExclude`, two-phase dedup, `detectBackupContradictions` | Fixture C: 58→6, 0 self-comparisons |
| **C — Trust Chain** | Hidden from CLI | 5 CLI fields, `convertAnalysisResultToCLI` | Fixture D: 100 chains in JSON |
| **D — SDRI** | Controls ignored, CIARE mismatches | `convertControlsToIntel`, `ciareControlAliases` | Fixture B: 10 YAML controls consumed |
| **Release** | coverage/ gitignored | `.gitignore` negation | CI `go vet` passes |
