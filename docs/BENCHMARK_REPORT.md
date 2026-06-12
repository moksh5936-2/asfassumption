# Benchmark Validation Report

**Version:** 2.1.2
**Date:** 2026-06-13
**Platform:** darwin/arm64

---

## Benchmark Datasets

All 5 benchmark datasets execute successfully and produce valid results:

| Dataset | File | Results |
|---|---|---|
| Healthcare | `testdata/attack_paths/healthcare.yaml` | Full pipeline |
| Fintech | `testdata/attack_paths/fintech.yaml` | Full pipeline |
| Cloud/Native | `testdata/attack_paths/cloud_native.yaml` | Full pipeline |
| Kubernetes | `testdata/attack_paths/kubernetes.yaml` | Full pipeline |
| VPN | `testdata/attack_paths/vpn.yaml` | Full pipeline |

## Intelligence Engine Benchmarks

All engine benchmarks pass. Key measurements:

| Benchmark | Iterations | ns/op | B/op | allocs/op |
|---|---|---|---|---|
| BenchmarkTMIEngine/auth0_saas | 285 | 4,184,561 | 1,972,186 | 33,135 |
| BenchmarkTMIEngine/minimal | 5,482 | 218,289 | 62,908 | 1,319 |
| BenchmarkTMIEngine/empty | 5,349 | 222,692 | 61,715 | 1,295 |
| BenchmarkTMIEngine/no_relationships | 4,666 | 256,677 | 69,696 | 1,396 |
| BenchmarkTMIEngine/http | 2,042 | 584,182 | 177,626 | 3,422 |

## Pipeline Engine Metrics (from benchmark_intelligence_test.go)

All 5 benchmark datasets produce consistent results with full pipeline execution:

| Metric | Healthcare | Fintech | Cloud | K8s | VPN |
|---|---|---|---|---|---|
| Assumptions | 73 | 74 | 57 | 40 | 20-31 |
| Critical | 19 | 34 | 12 | 11 | 3 |
| High | 31 | 16 | 29 | 14 | 7-15 |
| Threats | 29 | 27 | 18 | 13 | 7-10 |
| Attack Paths | 15 | 42 | 5 | 4 | 1-3 |
| Controls | 20 | 20 | 20 | 20 | 20 |
| CIE Contradictions | 43 | 85 | 23 | 23 | 8-17 |
| TBI Zones | 6 | 6 | 6 | 6 | 2-4 |
| TBI Weaknesses | 55 | 57 | 38 | 10 | 0-11 |
| ERN Themes | 4 | 5 | 4 | 4 | 4 |
| ERN Exposure | Severe | Severe | Severe | Severe | Severe |
| SDI Recommendations | 7 | 8 | 7 | 7 | 4 |
| SDT Change Impacts | 4 | 4 | 4 | 2 | 1-2 |
| SDT Control Drifts | 10 | 10 | 10 | 10 | 10 |

## Conclusion

- All benchmark datasets produce valid, consistent results
- All intelligence engines execute without errors
- Pipeline completes for every dataset configuration
- No crashes, panics, or unexpected behavior observed

**Status: ALL BENCHMARKS PASS**
