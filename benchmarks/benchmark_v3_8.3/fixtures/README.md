# Benchmark 3 — Fixture Inventory (v3.0.0-RC2)

## Available Fixtures

Same 6 fixtures as v2.8, located at `asf-tui/testdata/v29_rebenchmark/`:

| Fixture | File | Purpose |
|---------|------|---------|
| A — Parser Pollution | `fixture_a_parser_pollution.yaml` | Verify YAML comments don't become assumptions |
| B — Explicit Insecure | `fixture_b_explicit_insecure.yaml` | Verify `encryption: none` is respected |
| C — True Contradictions | `fixture_c_true_contradictions.yaml` | Test contradiction precision |
| D — Trust Chain | `fixture_d_trust_chain.yaml` | Test trust chain / SPOF detection |
| E — Positive Verification | `fixture_e_positive_verification.yaml` | Test positive verification |
| F — Blind Spot Review | `fixture_f_blind_spot_review.yaml` | Test review workbench |

## Benchmark Tests

Benchmark tests in `asf-tui/benchmark_test.go` execute all 6 fixtures:

| Test | Description |
|------|-------------|
| `TestBenchmarkAsftestYAML` | Baseline asftest.yaml execution |
| `TestBenchmarkContradictionDetection` | Basic contradiction detection |
| `TestBenchmarkTrustBoundaryDiscovery` | Trust boundary discovery |
| `TestBenchmarkDomainDetection` | Domain detection accuracy |
| `TestBenchmarkReportExport` | Report export functionality |
| `TestBenchmarkContradictionPrecision` | Fixture C — contradiction count ≤ 12 |
| `TestBenchmarkPositiveVerification` | Fixture E — VERIFIED ≥ 10, PARTIAL ≥ 3 |
| `TestBenchmarkTrustChainExposure` | Fixture D — trust chain data in CLI JSON |
| `TestBenchmarkSDRIControlAwareness` | Fixture B — SDRI controls consumed |
| `TestBenchmarkAllFixtures` | All 6 fixtures run end-to-end |

## Outputs

- `asf-tui/release/smoke_test_report.md` — Smoke test results on native binary
- `asf-tui/release/build_validation.txt` — Build validation output

## Test Harness

The `benchmark_test.go` test harness at `asf-tui/benchmark_test.go` provides structured before/after assertions for all 4 workstreams.

## CI Workflow

The `.github/workflows/release.yml` pipeline:
1. Builds and tests on every tag push
2. Runs `go vet`, `go test`, `go build`
3. Produces 5 platform binaries with SHA256 checksums
4. Creates GitHub release with all assets
