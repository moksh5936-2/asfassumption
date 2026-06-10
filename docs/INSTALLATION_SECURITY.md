# ASF Installation Security

> Version: 1.0.0 | June 2026

## Overview

This document reviews the security properties of the ASF distribution and installation system. It covers the `install.sh`, `install.ps1`, and GitHub Actions release workflow.

---

## 1. Threat Model

### Trust Boundaries

```
[User's Machine]                        [GitHub]
       │                                      │
       │ 1. curl pipe bash install.sh         │
       │─────────────────────────────────────>│
       │                                      │
       │ 2. Download binary + checksums.txt   │
       │<─────────────────────────────────────│
       │                                      │
       │ 3. Verify SHA-256 checksum          │
       │ (local, no network)                 │
       │                                      │
       │ 4. Install to /usr/local/bin/asf    │
```

### Assets

| Asset | Protection |
|-------|-----------|
| ASF binary | SHA-256 checksum from GitHub release |
| Install script | Served via HTTPS from GitHub RAW |
| User PATH | Script validates install directory |
| Existing installation | Backed up during upgrade |

### Attack Vectors

| Vector | Mitigation |
|--------|-----------|
| **Man-in-the-middle on install script** | Served over HTTPS from raw.githubusercontent.com |
| **Compromised GitHub release asset** | SHA-256 checksum from same release; user should verify against published hash |
| **Tampered checksums.txt** | Checksums.txt is served over HTTPS from the same release |
| **curl-pipe-bash injection** | Script is read-only from a trusted domain; HTTPS ensures integrity |
| **Temporary file leak** | All temp files are cleaned up via `trap` or try/finally |
| **Permission escalation** | Script prefers `/usr/local/bin` but falls back to `~/.local/bin` |
| **Upgrade failure** | Script downloads to temp dir first; existing binary is not replaced until verification passes |

---

## 2. Checksum Verification Flow

```
1. Download ASF binary
2. Download checksums.txt from same release
3. Look up expected hash for binary name in checksums.txt
4. Compute SHA-256 of downloaded binary
5. Compare: expected == computed?
   ├── YES → continue installation
   └── NO  → delete temp files, exit with error
```

**Limitation:** Both the binary and checksums.txt are served from the same GitHub release over the same HTTPS channel. A compromise of the GitHub release would affect both. For higher security, users should:

- Verify the checksum against a trusted source
- Use signed git tags and verify with `git tag -v`
- Pin to a specific version via `ASF_VERSION`

---

## 3. Script Security Properties

### install.sh

| Property | Implementation |
|----------|---------------|
| **HTTPS only** | All downloads use `https://` URLs |
| **Temp directory** | `mktemp -d` creates a secure temp directory |
| **Cleanup** | `trap 'rm -rf "${TMP_DIR}"' EXIT` ensures cleanup on any exit |
| **No eval** | No dynamic code execution via `eval` or `source` from untrusted input |
| **File size check** | Verifies `[ -s "${TMP_DIR}/asf" ]` after download |
| **Permission check** | Checks `[ -w "${INSTALL_DIR}" ]` before symlink |
| **PATH safety** | Warns if install directory not in PATH |
| **Upgrade safety** | Downloads to temp, verifies, then copies — never replaces in-place |

### install.ps1

| Property | Implementation |
|----------|---------------|
| **HTTPS only** | `Invoke-WebRequest` uses HTTPS |
| **Temp directory** | Uses `$env:TEMP\asf-install` |
| **Cleanup** | `Remove-Item -Recurse -Force` in try/finally |
| **No Invoke-Expression** | Uses `Invoke-WebRequest` (not `iex` for the download) |
| **File size check** | Checks `Length -eq 0` after download |
| **PATH safety** | Adds to user PATH, not system PATH |

---

## 4. GitHub Actions Security

### Workflow Permissions

```yaml
permissions:
  contents: write    # Required for creating releases
```

The workflow runs only on tag pushes matching `v*`. Artifacts are uploaded between jobs via `actions/upload-artifact`.

### Supply Chain

| Practice | Status |
|----------|--------|
| **GitHub Actions pinned versions** | ✅ Uses `@v4` major version pins |
| **Dependency cache** | ✅ `setup-go` with `cache: true` |
| **No third-party actions** | ✅ Only standard GitHub actions |
| **CGO_ENABLED=0** | ✅ Static binaries, no C library dependencies |
| **No secrets in build** | ✅ Only `github.token` for release creation |

---

## 5. Recommendation for Production

### Enhancements

1. **Code signing** — Sign macOS binaries with an Apple Developer ID for notarization
2. **Binary signatures** — Add GPG or Minisign signatures alongside checksums
3. **SLSA provenance** — Generate SLSA provenance attestations in CI
4. **Dependency scanning** — Add Dependabot for Go module vulnerability scanning
5. **SBOM generation** — Generate a Software Bill of Materials for each release
6. **Pin install script** — Encourage users to download and inspect before piping to bash

### Verification Commands

```bash
# Verify checksum manually
curl -sfLO https://github.com/moksh5936-2/asfassumption/releases/download/v1.0.0/checksums.txt
curl -sfLO https://github.com/moksh5936-2/asfassumption/releases/download/v1.0.0/ASF-v1.0.0-darwin-arm64
shasum -a 256 -c checksums.txt 2>/dev/null | grep OK

# Verify git tag signature
git tag -v v1.0.0

# Inspect the install script before running
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | less
```

---

## 6. Incident Response

If a compromised release is detected:

1. **Delete the GitHub release** — `gh release delete vX.Y.Z`
2. **Delete the git tag** — `git push --delete origin vX.Y.Z`
3. **Revoke any uploaded assets** — GitHub automatically removes them
4. **Publish security advisory** — Via GitHub Security Advisories
5. **Rotate any exposed keys** — If applicable

---

*Generated by `docs/INSTALLATION_SECURITY.md` — ASF Release Engineering Phase 9*
