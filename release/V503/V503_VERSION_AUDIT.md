# ASF0 v5.0.3 — Version Audit

## Files Updated

| File | Old Value | New Value |
|------|-----------|-----------|
| `asf-tui/license.go:18` | `ASFVersion = "5.0.2"` | `ASFVersion = "5.0.3"` |
| `asf-tui/Makefile:4` | `VERSION ?= 5.0.2` | `VERSION ?= 5.0.3` |
| `release/VERSION` | `5.0.0` | `5.0.3` |
| `release/install.sh:372` | `v5.0.0` fallback | `v5.0.3` |
| `release/install.sh:373` | `v5.0.0` fallback | `v5.0.3` |
| `release/install.sh:376` | `v5.0.0` fallback | `v5.0.3` |

## Files Verified (no change needed)

| File | Reason |
|------|--------|
| `asf-tui/.github/workflows/ci.yml` | Uses `GITHUB_REF_NAME#v` (dynamic) |
| `asf-tui/version_check.go` | Strips "v" prefix dynamically |
| `release/install.sh` | Download URL uses `LATEST_VERSION` variable |

## Cross-Reference: ASFVersion usages (13 files, all via var)
All usages refer to `ASFVersion` variable — no hardcoded strings:
- about.go, analyze_cli.go, app.go, cli.go, doctor.go, help.go, license.go
- main.go, telemetry.go, version_check.go, visuals.go

**Result:** All version references updated to v5.0.3.
