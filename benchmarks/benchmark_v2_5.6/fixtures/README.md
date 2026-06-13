# Benchmark 2 — Fixture Inventory (v2.8)

## Available Fixtures

Six dedicated benchmark fixtures were created for the v2.8 targeted re-benchmark:

| Fixture | File | Purpose |
|---------|------|---------|
| A — Parser Pollution | `testdata/v29_rebenchmark/fixture_a_parser_pollution.yaml` | Verify YAML comments don't become assumptions |
| B — Explicit Insecure | `testdata/v29_rebenchmark/fixture_b_explicit_insecure.yaml` | Verify `encryption: none` is respected |
| C — True Contradictions | `testdata/v29_rebenchmark/fixture_c_true_contradictions.yaml` | Test contradiction precision |
| D — Trust Chain | `testdata/v29_rebenchmark/fixture_d_trust_chain.yaml` | Test trust chain / SPOF detection |
| E — Positive Verification | `testdata/v29_rebenchmark/fixture_e_positive_verification.yaml` | Test positive verification |
| F — Blind Spot Review | `testdata/v29_rebenchmark/fixture_f_blind_spot_review.yaml` | Test review workbench |

## Location

All fixtures are in `asf-tui/testdata/v29_rebenchmark/`.

## Test Harness

Benchmark tests in `asf-tui/benchmark_test.go` execute these fixtures and produce structured results.

## Missing Artifacts

| Artifact | Status | Note |
|----------|--------|------|
| Raw CLI JSON output | NOT_AVAILABLE | Not preserved for v2.8 |
