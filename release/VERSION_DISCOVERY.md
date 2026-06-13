# ASF Version Discovery

## Current Version

| Source | Version |
|--------|---------|
| `license.go:18` — `var ASFVersion` | **4.0.0** |
| `release/VERSION` | 2.4.0 (outdated) |
| `scripts/build-release.sh` default | 2.2.0 (outdated) |
| `install.sh` `ASF_VERSION=` example | 2.0.1 (outdated) |
| `install.ps1` `ASF_VERSION=` example | 2.0.1 (outdated) |
| Installer GitHub API fallback | v3.0.0 (outdated) |

## Verification

```
$ asf --version
ASF v4.0.0

$ asf doctor | grep Version
Version:      4.0.0
```

## Build Version Strategy

- **Source of truth:** `license.go` — `var ASFVersion = "4.0.0"`
- **Release version:** `4.0.0`
- **Binary naming:** `ASF-v4.0.0-{os}-{arch}` (matches `scripts/build-release.sh` convention, matches `install.sh` asset URL pattern)
- **Tag:** `v4.0.0` (reserved for repository owner)

## Files Not Modified

- `release/VERSION` — outdated (2.4.0); not updated to avoid conflicting with owner's versioning
- `scripts/build-release.sh` default — left at 2.2.0
- `install.sh`/`install.ps1` examples — left as-is

The build uses `-ldflags="-s -w"` without overriding `ASFVersion`; the source constant is used directly.
