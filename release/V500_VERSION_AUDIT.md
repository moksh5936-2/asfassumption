# Version Audit — ASF0 v5.0.0

## Source of Truth
`asf-tui/license.go:18` → `var ASFVersion = "5.0.0"`

## Files Updated

| File | Old Value | New Value |
|------|-----------|-----------|
| `asf-tui/license.go` | `4.0.2` | `5.0.0` |
| `release/VERSION` | `4.0.0` | `5.0.0` |
| `scripts/build-release.sh` | `4.0.0` | `5.0.0` |
| `scripts/build-release.ps1` | `1.0.0` | `5.0.0` |
| `install.sh` (fallback) | `v4.0.0` | `v5.0.0` |
| `install.ps1` (fallback) | `v4.0.0` | `v5.0.0` |
| `release/install.sh` (fallback) | `v3.0.0` | `v5.0.0` |
| `asf-tui/install.sh` | `3.0.0` | `5.0.0` |
| `README.md` (header) | `v2.2.0` | `v5.0.0` |
| `README.md` (download table) | `v2.2.0` | `v5.0.0` |
| `README.md` (download URLs) | `v2.2.0` | `v5.0.0` |
| `README.md` (verify) | `v2.2.0` | `v5.0.0` |
| `README.md` (intel engines) | `v2.2.0` | `v5.0.0` |
| `README.md` (binary ref) | `v2.1.1+` | `v5.0.0+` |

## Unchanged
Go source files using `ASFVersion` variable (dynamic, no hardcoded version):
- `main.go`, `app.go`, `about.go`, `help.go`, `doctor.go`, `analyze_cli.go`, `version_check.go`

CI/CD (`release.yml`) extracts version from git tag at build time — no change needed.

## Verification
```
go run . --version  → ASF0 v5.0.0
```
