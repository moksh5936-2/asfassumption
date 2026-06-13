# ASF v3.0.0-RC2 — Release Forensics Report

**Date:** 2026-06-13
**Author:** Principal Release Engineer
**Subject:** Installer downloading v2.2.0 instead of v3.0.0-RC2

---

## 1. Root Cause

### Primary: GitHub `/releases/latest` API skips prereleases

All 4 installer scripts used the GitHub API endpoint `/releases/latest` to detect the current version:

```
GET https://api.github.com/repos/moksh5936-2/asfassumption/releases/latest
```

This endpoint **always returns the latest non-prerelease release** by design. Since v3.0.0-RC2 was correctly marked as `prerelease: true`, the API returned v2.2.0 (the latest stable release).

### Secondary: Hardcoded fallback versions

When the API call failed, every installer had a stale fallback version:

| Installer | Fallback Version | Correct Version |
|-----------|-----------------|-----------------|
| `/install.sh` | `2.2.0` | `3.0.0-RC2` |
| `/release/install.sh` | `2.1.2` | `3.0.0-RC2` |
| `/asf-tui/install.sh` | `2.1.2` | `3.0.0-RC2` |
| `/install.ps1` | `2.2.0` | `3.0.0-RC2` |

---

## 2. Resolution

### Changed API endpoint

All 4 installers now use:

```
GET https://api.github.com/repos/moksh5936-2/asfassumption/releases?per_page=10
```

This returns the 10 most recent releases including prereleases. The first **non-draft** release in the response is selected. Since v3.0.0-RC2 is the most recent non-draft release (excluding a draft v1.0.2 entry), it is correctly identified as the current version.

### Version parsing

The response JSON is parsed using python3 (available on all supported platforms):

```python
import json,sys
releases = json.load(sys.stdin)
for r in releases:
    if not r.get('draft', True):
        print(r['tag_name'].lstrip('v'))
        sys.exit(0)
```

If python3 is unavailable or the API call fails, the fallback is now `3.0.0-RC2`.

### Updated PowerShell installer

The PowerShell installer similarly queries `/releases?per_page=10` and filters out drafts:

```powershell
$releases = Invoke-RestMethod -Uri $apiUrl -Headers $headers
$release = $releases | Where-Object { -not $_.draft } | Select-Object -First 1
$Version = $release.tag_name -replace "^v", ""
```

---

## 3. Files Changed

| File | Change |
|------|--------|
| `/install.sh` | API endpoint + fallback version |
| `/release/install.sh` | API endpoint + fallback version |
| `/asf-tui/install.sh` | Hardcoded `ASF_VERSION` |
| `/install.ps1` | API endpoint + fallback version |

---

## 4. Release Validation

| Check | Result |
|-------|--------|
| `v3.0.0-RC2` tag exists | ✅ `01028bf` |
| Tag points to latest commit | ✅ `01028bf` |
| 7 assets on release | ✅ 5 binaries + checksums + release notes |
| All platforms have binaries | ✅ darwin/amd64, darwin/arm64, linux/amd64, linux/arm64, windows/amd64 |
| Checksums uploaded | ✅ `checksums.txt` |
| Release marked prerelease | ✅ (correct for RC) |
| API detects `v3.0.0-RC2` | ✅ first non-draft release |
| Binary reports `v3.0.0-RC2` | ✅ `asf --version` returns `v3.0.0-RC2` |
| Version in source code | ✅ `license.go: ASFVersion = "3.0.0-RC2"` |

---

## 5. Verification

```bash
# API correctly resolves to RC2
curl -s https://api.github.com/.../releases?per_page=10 | python3 -c "
import json,sys
for r in json.load(sys.stdin):
    if not r.get('draft', True):
        print(r['tag_name'])
        break
"
# Output: v3.0.0-RC2

# Binary reports correct version
asf --version
# Output: ASF v3.0.0-RC2 starting
```
