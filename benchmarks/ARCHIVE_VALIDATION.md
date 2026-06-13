# Benchmark Archive Validation

## Validation Date

2026-06-13

---

## Completeness

### Required Structure

```
benchmarks/
├── benchmark_index.md                               ✅
├── benchmark_history.md                             ✅
├── benchmark_timeline.md                            ✅
├── evidence_matrix.md                               ✅
├── ARCHIVE_VALIDATION.md                            ✅
├── benchmark_v1_3.7/
│   ├── summary.md                                   ✅
│   ├── report.md                                    ✅
│   ├── fixtures/README.md                           ✅
│   └── outputs/                                     ✅ (empty — no outputs preserved)
├── benchmark_v2_5.6/
│   ├── summary.md                                   ✅
│   ├── report.md                                    ✅
│   ├── fixtures/README.md                           ✅
│   └── outputs/                                     ✅ (empty — no outputs preserved)
└── benchmark_v3_8.3/
    ├── summary.md                                   ✅
    ├── report.md                                    ✅
    ├── fixtures/README.md                           ✅
    └── outputs/                                     ✅ (empty — no outputs preserved)
```

### File Count

| Directory | Files |
|-----------|-------|
| Root (`benchmarks/`) | 5 |
| `benchmark_v1_3.7/` | 3 (+1 empty dir) |
| `benchmark_v2_5.6/` | 3 (+1 empty dir) |
| `benchmark_v3_8.3/` | 3 (+1 empty dir) |
| **Total** | **14 files** |

---

## Data Accuracy Verification

| Data Point | Source | Archive Value | Match |
|------------|--------|---------------|-------|
| v2.2 score | `V29_TARGETED_REBENCHMARK.md` §6 | 3.7/10 | ✅ |
| v2.8 score | `V29_TARGETED_REBENCHMARK.md` §6 | 5.6/10 | ✅ |
| v3.0.0-RC2 score | `V300_RC2_CISO_READINESS.md` §5 | 8.3/10 | ✅ |
| v2.2 verdict | `V29_TARGETED_REBENCHMARK.md` §6 | NOT_READY | ✅ |
| v2.8 verdict | `V29_TARGETED_REBENCHMARK.md` §8 | IMPROVED_BUT_NOT_READY | ✅ |
| v3.0.0-RC2 verdict | `V300_RC2_CISO_READINESS.md` §8 | PILOT_READY | ✅ |
| v2.8 VERIFIED count | `V29_TARGETED_REBENCHMARK.md` Fixture E | 0 | ✅ |
| v3.0.0-RC2 VERIFIED count | `V300_RC2_CISO_READINESS.md` WS A | 11 | ✅ |
| v2.8 contradiction count | `V29_TARGETED_REBENCHMARK.md` Fixture C | 58 | ✅ |
| v3.0.0-RC2 contradiction count | `V300_RC2_CISO_READINESS.md` WS B | 6 | ✅ |
| v2.8 trust chains | `V29_TARGETED_REBENCHMARK.md` Fixture D | 100 | ✅ |
| v3.0.0-RC2 trust chains | `V300_RC2_CISO_READINESS.md` Fixture D | 100 | ✅ |
| v3.0.0-RC2 SDRI controls | `V300_RC2_CISO_READINESS.md` WS D | 10 | ✅ |

---

## Missing Artifacts

| Artifact | Reason |
|----------|--------|
| v2.2 raw CLI JSON output | Not preserved at time of benchmark |
| v2.2 fixture files | Created later for v2.8 re-benchmark |
| v2.8 raw CLI JSON output | Not preserved at time of benchmark |
| v3.0.0-RC2 raw CLI JSON output from benchmark | Not preserved (CI build output available in `dist/`) |
| Benchmark test logs from v2.2 | Not preserved |
| Benchmark test logs from v2.8 | Not preserved |

**Completeness:** All documentation artifacts are created. Raw test outputs from earlier benchmarks are not available. This is expected — the archive was created retrospectively.

---

## Validation Status

| Criterion | Status |
|-----------|--------|
| All required files exist | ✅ PASS |
| All scores match source reports | ✅ PASS |
| All verdicts match source reports | ✅ PASS |
| No fabricated evidence | ✅ PASS (all data sourced from existing documents) |
| Missing artifacts documented | ✅ PASS |
| Directory structure complete | ✅ PASS |
| Links use relative paths | ✅ PASS |

**Archive Completeness:** **100%** (all planned artifacts created; missing raw outputs documented)

---

## Archive Readiness Verdict

**BENCHMARK_ARCHIVE_COMPLETE**
