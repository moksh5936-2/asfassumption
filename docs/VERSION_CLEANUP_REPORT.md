# Version Cleanup Report

## Summary

All hardcoded stale version strings referencing old versions (v0.1.0, v2.0.0, v2.0.1) have been updated to use the `ASFVersion` constant (v2.1.1) or directly replaced with v2.1.1.

## Fixes Applied

### dashboard.go (MUST_UPDATE)
- **File:** `asf-tui/dashboard.go:57`
- **Change:** `version := "v0.1.0"` → `version := "v" + ASFVersion`
- The constant `ASFVersion = "2.1.1"` is defined in `asf-tui/license.go:14` (same `package main`), so no new import was needed.

### README.md (MUST_UPDATE — 2 fixes)
| Line | Before | After |
|------|--------|-------|
| 212 | `# Expected: ASF v2.0.0` | `# Expected: ASF v2.1.1` |
| 469 | `ASF v2.0.0+` | `ASF v2.1.1+` |

### docs/INSTALLER_AND_COMMAND_RELIABILITY_REPORT.md (MUST_UPDATE)
- All `v2.0.1` → `v2.1.1` references updated (header, body, version table — 12 instances).

### docs/COMMAND_COVERAGE_AUDIT.md (MUST_UPDATE)
- `Target: v2.0.1` → `Target: v2.1.1` (line 3).

## HISTORICAL_OK References (intentionally preserved)

These are historical records of past versions and must not be changed:

| Pattern | File | Count |
|---------|------|-------|
| v0.1.0 | `docs/PRODUCTION_READINESS_AUDIT.md` | 2 |
| 2.0.0 / v2.0.0 | `CHANGELOG.md` | 2 |
| 2.0.0 / v2.0.0 | `docs/PRODUCTION_READINESS_AUDIT.md` | 4 |
| 2.0.0 / v2.0.0 | `docs/GO_NATIVE_SINGLE_BINARY_RELEASE_REPORT.md` | ~20 |
| 2.0.1 / v2.0.1 | `CHANGELOG.md` | 1 |
| 2.0.1 / v2.0.1 | `docs/PRODUCTION_READINESS_AUDIT.md` | 3 |

## DEPENDENCY_OK References (externally pinned, not ASF versions)

| Pattern | File | Reason |
|---------|------|--------|
| `pydantic>=2.0.0` | `pyproject.toml` | Python dependency pin |
| `v2.0.1` | `asf-tui/go.mod`, `asf-tui/go.sum` | Go dependency (`go-osc52/v2`) |
| `v2.0.1` | `docs/TECHNICAL_REFERENCE.md` | Go dependency docs |
| `ASF_VERSION=2.0.1` | `install.sh`, `release/install.sh` | Comment example for version pinning |
| `$env:ASF_VERSION="2.0.1"` | `install.ps1` | Comment example for version pinning |
| `ASF_VERSION="2.0.1"` | `scripts/test-installer.sh` | Test fixture pinning old version |

## Result

No stale version strings remain in active code or current documentation. The dashboard now correctly displays `ASFVersion` (v2.1.1).
