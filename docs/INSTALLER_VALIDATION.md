# Installer Validation Report

**Version:** 2.2.0
**Date:** 2026-06-13

---

## Install Scripts

| Script | Platform | Lines | Status |
|---|---|---|---|
| `install.sh` | macOS, Linux | 592 | PASS (code review) |
| `install.ps1` | Windows | 332 | PASS (code review) |

## Operations Supported

| Operation | `install.sh` | `install.ps1` |
|---|---|---|
| Install | ✓ | ✓ |
| Upgrade (--upgrade/-u) | ✓ | ✓ |
| Repair (--repair) | ✓ | ✓ |
| Clean reinstall (--clean) | ✓ | ✓ |
| Purge (--clean --purge) | ✓ | ✓ |
| Help (--help/-h) | ✓ | ✓ |

## Validation Points

- **Platform detection**: Both scripts detect OS and architecture correctly
- **Version pinning**: Both support `ASF_VERSION` env var for version pinning
- **Checksum verification**: Both download and verify SHA-256 checksums from `checksums.txt`
- **Binary verification**: Both run `asf --version` to confirm the downloaded binary matches the expected version
- **Config backup**: Both backup `config.yaml` and `license.key` before upgrade
- **PATH setup**: Both add install directory to PATH (shell config on Unix, User PATH on Windows)
- **Symlink management**: `install.sh` uses symlinks; `install.ps1` creates `bin\asf.exe` copy
- **Certificate/extraction**: Both use official release artifacts from GitHub Releases

## Default Configuration

Both installers create the same default `config.yaml`:

```yaml
general:
  theme: Dark
  fox_style: Classic
analysis:
  depth: deep
  stride: true
  controls: true
ai:
  enabled: false
output:
  default: markdown
  directory: ./reports
engine:
  use_native_engine: true
```

## Conclusion

**Both installers are structurally sound, support all required operations, and follow security best practices (checksum verification, binary verification, config backup).**
