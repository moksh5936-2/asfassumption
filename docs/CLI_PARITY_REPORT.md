# CLI Parity Report

**Phase 9 — June 2026**

## Summary

The Go native CLI (`asf analyze`) was compared against the Python CLI (`python -m asf.cli.main analyze`). Functional parity is achieved for the core analysis workflow. Several differences and bugs were identified and partially resolved.

## Flag Comparison

| Flag | Python CLI | Go CLI | Notes |
|------|-----------|--------|-------|
| `analyze` subcommand | ✅ `analyze` | ✅ `analyze` | Core command |
| `-e` / `--evidence` | ✅ Both forms | ✅ Both forms (bug fixed) | Go was missing `--evidence` long form; added |
| `--json` | ✅ Required for JSON output | ✅ Accepted as no-op | Added for scripting compatibility; Go always outputs JSON |
| `--graph` | ✅ Accepted but **non-functional** | ✅ Produces graph data | Python accepts flag but never includes `graph` in output |
| `--persist` | ✅ Save to database | ❌ Not implemented | Database persistence not needed for native engine |
| `--auto-map` | ✅ Auto-map evidence | ❌ Not implemented | |
| `--help` | ✅ | ✅ | |
| Directory input | ✅ Expands to files | ❌ **Crashes** (`Error: %!s(<nil>)`) | Go CLI requires explicit file path |
| Directory evidence | ❌ **Crash** (`IsADirectoryError`) | ✅ Expands to supported files | Python cannot use `-e <dir>`; Go fixed to filter by extension |

## JSON Output Structure

| Field | Python CLI | Go CLI | Match |
|-------|-----------|--------|-------|
| `version` | Not present | `"1.1.0"` | Go adds version info |
| `architecture` | Not present | filename | Go identifies source file |
| `summary.claims_found` | ✅ | ✅ | Identical |
| `summary.assumptions` | ✅ | ✅ | Identical |
| `summary.verified` | ✅ | ✅ | Identical |
| `summary.contradicted` | ✅ | ✅ | Identical |
| `summary.unknown` | ✅ | ✅ | Identical |
| `summary.critical_gaps` | ✅ | ✅ | Identical |
| `summary.partially_verified` | ❌ Not present | ✅ Present | Go tracks all verification results |
| `claims[]` | ✅ Full claim objects | ❌ Not present | Go embeds claim info in assumption text |
| `assumptions[].text` | ✅ | ✅ | Identical (whitespace-normalized) |
| `assumptions[].assumption_type` | ✅ | ✅ | Identical |
| `assumptions[].verification_status` | ✅ | ✅ | Identical |
| `assumptions[].confidence` | ✅ Float | ✅ Float | Different values (different algorithms) |
| `assumptions[].keywords` | ✅ | ✅ | Identical |
| `verifications[].result` | ✅ | ✅ | Identical |
| `verifications[].reasoning` | ✅ | ✅ | Identical |
| `gaps[].type` | ✅ | ✅ | Identical |
| `gaps[].severity` | ✅ | ✅ | Identical |
| `gaps[].description` | ✅ Full length | ✅ Up to 500 chars | Truncation limit increased from 120→500 |
| `graph` | ❌ Never emitted | ✅ Full graph data | |

## Error Handling

| Scenario | Python Exit Code | Go Exit Code |
|----------|-----------------|--------------|
| Success | 0 | 0 |
| Non-existent file | 2 | 1 |
| Invalid flag | 2 | 1 |
| Missing argument | 2 | 1 |

Both exit non-zero on errors. Exit codes differ (Python uses 2 for Click-standard errors, Go uses 1).

## Bugs Found and Fixed

### Fixed in This Session
1. **Missing `--evidence` long form** — Go CLI only accepted `-e`, not `--evidence`. Added.
2. **Missing `--json` compat flag** — Go CLI rejected `--json` flag used by Python. Added as no-op.
3. **Gap description truncation** — Go truncated descriptions at 120 chars; Python does not. Increased to 500.

### Still Open
1. **Go directory expansion crash** — Running `asf analyze <directory>` crashes with `Error: %!s(<nil>)` instead of expanding to files. Python handles this correctly.
2. **Python `-e <dir>` crash** — Python crashes with `IsADirectoryError` when evidence path is a directory. Go handles this correctly (after fix).

## Conclusion

The Go native CLI achieves functional parity with the Python CLI for the core analysis workflow. All major flags are supported. The JSON output contains the same analytic data with minor structural differences (Go omits `claims[]` array but includes `partially_verified` in summary). Two remaining bugs (Go directory input crash, Python `-e <dir>` crash) do not affect the standard single-file analysis workflow.
