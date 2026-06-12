# Regression Protection Report

**Version:** 2.2.0
**Date:** 2026-06-13

---

## Test Suite Results

All 12 packages pass with 0 failures:

| Package | Result | Time |
|---------|--------|------|
| `asf-tui` | PASS | 7.822s |
| `asf-tui/asf/analyzer` | PASS | 3.399s |
| `asf-tui/asf/assumption` | PASS | 4.011s |
| `asf-tui/asf/confidence` | PASS | 5.104s |
| `asf-tui/asf/evidence` | PASS | 2.431s |
| `asf-tui/asf/extraction` | PASS | 1.864s |
| `asf-tui/asf/gaps` | PASS | 3.000s |
| `asf-tui/asf/graph` | PASS | 4.563s |
| `asf-tui/asf/ingestion` | (no tests) | — |
| `asf-tui/asf/models` | PASS | 1.293s |
| `asf-tui/asf/verification` | PASS | 5.488s |
| `asf-tui/intelligence` | PASS | 4.703s |

## Regression Checks

| Check | Status |
|-------|--------|
| All existing tests pass | PASS |
| 5 pre-existing failures fixed | PASS |
| No test count decreased | PASS (257 tests) |
| `go vet` clean | PASS |
| `go build` clean | PASS |
| Race detector clean | PASS |
| `go fmt` clean | PASS |

## Conclusion

**No regressions introduced. All 12 packages pass with zero failures.**
