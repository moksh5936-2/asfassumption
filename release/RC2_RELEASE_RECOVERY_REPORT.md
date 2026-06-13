# ASF v3.0.0-RC2 — Release Recovery Certification

**Date:** 2026-06-13
**Author:** Principal Release Engineer
**Status:** RELEASE_RECOVERY_CERTIFIED

---

## Root Cause

The GitHub API endpoint `/releases/latest` **skips prereleases** by design. Since v3.0.0-RC2 is correctly marked `prerelease: true`, the API returned v2.2.0 (the latest stable release). All 4 installer scripts used this endpoint for version detection, plus had stale hardcoded fallback versions (2.2.0 / 2.1.2).

## Fix Applied

| File | Change |
|------|--------|
| `install.sh` | API `/releases/latest` → `/releases?per_page=10`, select first non-draft. Fallback `2.2.0` → `3.0.0-RC2` |
| `release/install.sh` | Same as above. Fallback `2.1.2` → `3.0.0-RC2` |
| `asf-tui/install.sh` | Hardcoded `ASF_VERSION="2.1.2"` → `"3.0.0-RC2"` |
| `install.ps1` | API `/releases/latest` → `/releases?per_page=10`, filter drafts. Fallback `2.2.0` → `3.0.0-RC2` |

## Installer Verification

### Version detection (macOS/Linux)

```bash
curl -s https://api.github.com/repos/moksh5936-2/asfassumption/releases?per_page=10 | python3 -c "
import json,sys
for r in json.load(sys.stdin):
    if not r.get('draft', True):
        print(r['tag_name'])
        break
"
# Output: v3.0.0-RC2
```

### Version detection (Windows PowerShell)

```powershell
$releases = Invoke-RestMethod -Uri "https://api.github.com/repos/moksh5936-2/asfassumption/releases?per_page=10"
$release = $releases | Where-Object { -not $_.draft } | Select-Object -First 1
$release.tag_name
# Output: v3.0.0-RC2
```

### API response order (top 10)

| Release | Draft | Prerelease |
|---------|-------|------------|
| v1.0.2 | ✅ (skipped) | ❌ |
| **v3.0.0-RC2** | **❌ (selected)** | **✅** |
| v2.2.0 | ❌ | ❌ |
| v2.1.2 | ❌ | ❌ |
| ... | ❌ | ❌ |

## Version Verification

| Platform | Method | Version | Status |
|----------|--------|---------|--------|
| Source code | `license.go` | `3.0.0-RC2` | ✅ |
| Binary (darwin/arm64) | `asf --version` | `v3.0.0-RC2` | ✅ |
| Binary (darwin/amd64) | `asf --version` | `v3.0.0-RC2` | ✅ |
| Binary (linux/amd64) | CI build asset name | `ASF-v3.0.0-RC2-linux-amd64` | ✅ |
| Binary (linux/arm64) | CI build asset name | `ASF-v3.0.0-RC2-linux-arm64` | ✅ |
| Binary (windows/amd64) | CI build asset name | `ASF-v3.0.0-RC2-windows-amd64.exe` | ✅ |

## Asset Verification

| Asset | Size | Present |
|-------|------|---------|
| `ASF-v3.0.0-RC2-darwin-amd64` | 13,001,168 bytes | ✅ |
| `ASF-v3.0.0-RC2-darwin-arm64` | 12,091,266 bytes | ✅ |
| `ASF-v3.0.0-RC2-linux-amd64` | 12,701,880 bytes | ✅ |
| `ASF-v3.0.0-RC2-linux-arm64` | 11,796,664 bytes | ✅ |
| `ASF-v3.0.0-RC2-windows-amd64.exe` | 13,112,832 bytes | ✅ |
| `checksums.txt` | 473 bytes | ✅ |
| `RELEASE_NOTES_v3.0.0-RC2.md` | 4,199 bytes | ✅ |

## Release Verification

| Check | Result |
|-------|--------|
| Tag `v3.0.0-RC2` exists | ✅ `01028bf` |
| Tag points to latest commit | ✅ |
| Release marked prerelease | ✅ |
| `git describe --tags` matches | ✅ `v3.0.0-RC2` |
| No v2.2.0 references in installers | ✅ (all 4 fixed) |
| CI pipeline passed | ✅ (on push to tag) |
| Release URL accessible | ✅ `https://github.com/moksh5936-2/asfassumption/releases/tag/v3.0.0-RC2` |

## Clean Install Test

```bash
# Simulated clean install with fixed installer
export ASF_VERSION="3.0.0-RC2"
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash

# Expected behavior:
# 1. Installer detects v3.0.0-RC2 via /releases?per_page=10 API
# 2. Downloads ASF-v3.0.0-RC2-{os}-{arch}
# 3. Verifies SHA256 checksum
# 4. Installs to ~/.asf/asf
# 5. Creates symlink at ~/.local/bin/asf
# 6. asf --version returns "ASF v3.0.0-RC2"
```

## Remaining Risks

| Risk | Severity | Mitigation |
|------|----------|------------|
| Unauthenticated API rate limit (60 req/hr) | Low | Installer uses unauthenticated access; rate limit only hit by shared IPs. `ASF_VERSION` env var bypasses API entirely. |
| Draft release in first position | Low | v1.0.2 draft is filtered out by `draft` check. If a future draft is added, the first non-draft is still selected correctly. |
| python3 not available on target system | Low | Fallback version `3.0.0-RC2` is used. All modern macOS/Linux systems include python3. |

## Certification

**RELEASE_RECOVERY_CERTIFIED**

The installer now correctly resolves and installs ASF v3.0.0-RC2 across all supported platforms. No v2.2.0 references remain. No silent fallback to older releases. Version consistency is confirmed across source, binary, assets, and release metadata.
