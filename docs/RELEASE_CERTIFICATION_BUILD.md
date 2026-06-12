# Release Certification — Build

## Commands Run

```bash
go fmt ./...        # 2 files formatted (minor whitespace)
go vet ./...        # clean, no warnings
go test -count=1 ./... # all pass
go build -o /tmp/asf-cert . # success
```

## Results

| Metric | Value |
|--------|-------|
| Packages | 10 |
| Test Functions | 181 |
| Pass | 181 |
| Fail | 0 |
| go vet warnings | 0 |
| go fmt changes | 2 (minor, non-functional) |
| Build time | 0.967s |
| Binary size | 12M |
| Binary type | Mach-O 64-bit arm64 |
| Build errors | 0 |

## Packages

| Package | Tests | Time | Status |
|---------|-------|------|--------|
| asf-tui | multiple | 4.510s | PASS |
| asf-tui/asf/analyzer | multiple | 3.161s | PASS |
| asf-tui/asf/assumption | multiple | 0.728s | PASS |
| asf-tui/asf/confidence | multiple | 2.596s | PASS |
| asf-tui/asf/evidence | multiple | 1.372s | PASS |
| asf-tui/asf/extraction | multiple | 1.905s | PASS |
| asf-tui/asf/gaps | multiple | 3.691s | PASS |
| asf-tui/asf/graph | multiple | 5.055s | PASS |
| asf-tui/asf/ingestion | 0 | — | no test files |
| asf-tui/asf/models | multiple | 5.588s | PASS |
| asf-tui/asf/verification | multiple | 6.276s | PASS |

## Verdict

✅ **PASS** — Build is clean, tests pass, binary builds successfully.
