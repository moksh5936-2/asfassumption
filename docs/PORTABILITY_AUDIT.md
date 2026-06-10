# Portability Audit

## Complete Findings: Hardcoded Development Paths

### CRITICAL — Hardcoded Developer Machine Paths

| # | File | Line | Code | Issue | Severity |
|---|------|------|------|-------|----------|
| 1 | `asf-tui/engine.go` | 106 | `pythonPath: "/Users/moksh/Project/cybersec/.venv/bin/python"` | Dev machine Python venv | **FIXED** |
| 2 | `asf-tui/engine.go` | 107 | `projectDir: "/Users/moksh/Project/cybersec"` | Dev machine project root | **FIXED** |
| 3 | `asf-tui/engine.go` | 182 | `cmd.Dir = e.projectDir` | Propagated hardcoded path | **FIXED** |
| 4 | `benchmark/experiments/build_master_ground_truth.py` | 11-12 | `/Users/moksh/Project/cybersec/...` | Benchmark script paths | Not in binary |
| 5 | `benchmark/experiments/run_asf_tests.py` | 430, 520 | `/Users/moksh/Project/cybersec/...` | Benchmark script paths | Not in binary |
| 6 | `benchmark/experiments/generate_campaign_001.py` | 4 | `/Users/moksh/Project/cybersec/...` | Benchmark script paths | Not in binary |
| 7 | `benchmark/experiments/asf_diagnostic_tests.md` | ~80 refs | `/Users/moksh/Project/cybersec/...` | Generated docs | Not in binary |

### HIGH — CWD-Dependent Defaults (Acceptable with config override)

| # | File | Line | Code | Issue |
|---|------|------|------|-------|
| 8 | `asf-tui/config.go` | 48 | `"./reports"` | Default export relative to CWD — user-configurable |
| 9 | `asf-tui/results.go` | 51 | `"./reports"` | Export path — uses config default |
| 10 | `asf/config.py` | 7 | `Path("asf.config.yaml")` | Python engine default — CWD-relative |

### MEDIUM — Platform-Specific Binary Paths (Acceptable fallbacks)

| # | File | Line | Code | Issue |
|---|------|------|------|-------|
| 11 | `asf-tui/parser.go` | 621 | `/usr/local/bin/tesseract`, `/opt/homebrew/bin/tesseract` | Binary search paths |
| 12 | `asf-tui/model.go` | 39 | `/usr/local/bin/ollama`, `/opt/homebrew/bin/ollama` | Binary search paths |

### FIXED — Go Source Changes

| File | Before | After |
|------|--------|-------|
| `asf-tui/engine.go:106-107` | Hardcoded `/Users/moksh/...` paths | `discoverPythonPath()` runtime search |
| `asf-tui/engine.go:182` | `cmd.Dir = e.projectDir` | `cmd.Dir = asfCacheDir()` |
| `asf-tui/config.go:48-71` | Hardcoded `~/.asf/config.yaml` | XDG/AppData paths via `paths.go` |
| `asf-tui/license.go:26-31` | Hardcoded `~/.asf/license.key` | XDG/AppData paths via `paths.go` |
| `asf-tui/main.go:24` | `~/.asf/license.key` in message | Dynamic `asfLicensePath()` |
| `asf-tui/main.go` | No `--help`/`doctor` | Added `asf doctor`, `--help` |
| `asf-tui/paths.go` | Did not exist | Created with XDG/AppData strategy |

### FIXED — Installer Changes

| File | Before | After |
|------|--------|-------|
| `install.sh` | Config at `~/.asf/config.yaml` | Config at XDG path per platform |
| `install.ps1` | Config at `%LOCALAPPDATA%\ASF` | Config at `%APPDATA%\ASF` |

### NOT FIXED (Benchmark/Test/CI Scripts — Not in Binary Distribution)

These files contain hardcoded paths but are NOT part of the distributed ASF binary:

- `benchmark/experiments/build_master_ground_truth.py`
- `benchmark/experiments/run_asf_tests.py`
- `benchmark/experiments/generate_campaign_001.py`
- `benchmark/experiments/asf_diagnostic_tests.md`
- `tests/validation_harness.py` (uses `__file__`-based paths — acceptable for tests)
- `tests/test_ingestion.py` (uses `__file__`-based paths)
- `tests/test_cli.py` (uses `__file__`-based paths)
- `tests/test_analyzer.py` (uses `__file__`-based paths)

These are test/benchmark files that are not compiled into the Go binary. They are excluded from release builds.
