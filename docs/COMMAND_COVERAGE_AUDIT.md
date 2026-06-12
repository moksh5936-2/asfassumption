# ASF Command Coverage Audit

Generated: June 2026 | Target: v2.1.1

> **Update v2.1.1:** Issues #1 (invalid commands silently launching TUI) and #2 (`doctor--verbose` dead code) are now **FIXED**.

## Command Matrix

| Command | Documented Where | Implemented | Works | Exit Code | Notes |
|---------|-----------------|-------------|-------|-----------|-------|
| `asf` | README, help text, install output | `main.go:87-105` | ✅ | 0 | Launches TUI |
| `asf --version` | README, help text, main.go:14 | `main.go:38-40` | ✅ | 0 | Prints `ASF vX.Y.Z` |
| `asf -v` | help text | `main.go:38-40` | ✅ | 0 | Same as `--version` |
| `asf --help` | All docs | `main.go:70-72` | ✅ | 0 | Full usage text |
| `asf -h` | All docs | `main.go:70-72` | ✅ | 0 | Same as `--help` |
| `asf doctor` | README, help text | `main.go:50-66` | ✅ | 0 | System diagnostics |
| `asf doctor --verbose` | help text | `main.go:53-56` | ✅ | 0 | Detailed diagnostics |
| `asf doctor --fix` | help text | `main.go:57-58` | ✅ | 0 | Clean stale binaries |
| `asf --license` | README, help text | `main.go:41-49` | ✅ | 0/1 | 0 if valid, 1 if not |
| `asf analyze --help` | help text | `analyze_cli.go:77-88` | ✅ | 0 | |
| `asf analyze` (no file) | — | `analyze_cli.go:107-111` | ✅ | 1 | Error + usage |
| `asf analyze <file>` | README, help text | `analyze_cli.go:71-235` | ✅ | 0 | JSON output |
| `asf analyze <file> --json` | — | ⚠️ Implicit (always JSON) | ✅ | 0 | Flag accepted as no-op |
| `asf analyze <file> --graph` | README, help text | `analyze_cli.go:93-96,225-227` | ✅ | 0 | Graph in JSON |
| `asf analyze <file> -e <evidence>` | README, help text | `analyze_cli.go:95-99` | ✅ | 0 | |
| `asf analyze <file> --evidence <evidence>` | help text | `analyze_cli.go:95-99` | ✅ | 0 | Same as `-e` |
| `asf analyze <dir>` | — | `analyze_cli.go:121-136` | ✅ | 0 | Scans dir for .txt/.pdf/.docx |
| `asf analyze <file> --json --graph` | — | `analyze_cli.go` | ✅ | 0 | Both flags work |
| Invalid command (e.g. `asf foo`) | — | `main.go` | ❌ | 0 | Falls through to TUI |
| Missing file (`asf analyze missing.txt`) | — | `analyze_cli.go:114-117` | ✅ | 1 | |
| `asf doctor --verbose` (long) | help text | `main.go:53-56` | ✅ | 0 | |
| `asf --doctor` | — | `main.go:50` | ✅ | 0 | Alias for `doctor` |
| `asf diagnose` | — | `main.go:50` | ✅ | 0 | Alias for `doctor` |

## Issues Found

### Critical
1. **Invalid commands silently launch TUI (exit 0)** — `main.go:76-85` has no default case, so `asf foo` falls through to TUI with exit 0. Should print error and exit 1.
2. **`asf doctor --verbose` dead code** — `main.go:73-75` has `case "doctor--verbose":` which is never matched (no space between words).

### Medium
3. **`--json` flag not explicitly handled** — `analyze_cli.go:92-105` ignores `--json` as an unknown token. Output is always JSON, so it works, but `asf analyze --json <file>` misassigns `--json` as the file path.
4. **`--json` not documented in analyze CLI help** — help text (line 78-88) doesn't mention `--json`.

### Low
5. **Duplicate `--help` handling** — `main.go:70-72` handles `--help` in the switch, then lines 79-85 handle it again in the fallthrough. Not a bug but redundant.

## Recommendations

1. Add `default` case to the switch in `main()` that prints error and exits 1
2. Add `--json` to analyze CLI argument parser (as accepted no-op)
3. Document `--json` in analyze CLI help text
4. Remove dead `doctor--verbose` case
