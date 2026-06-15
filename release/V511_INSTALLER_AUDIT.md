# V511 — Installer Audit

## Installer (`install.sh`)

| Check | Status |
|---|---|
| REPO correct (`moksh5936-2/asfassumption`) | ✓ (line 20) |
| Default version v5.1.1 | ✓ (lines 372-376) |
| Asset URL pattern `ASF-${LATEST_VERSION}-${OS}-${ARCH}` | ✓ (line 386) |
| Download URL `releases/download/${LATEST_VERSION}/${BINARY_NAME}` | ✓ (line 387) |
| Checksums URL `releases/download/${LATEST_VERSION}/checksums.txt` | ✓ (line 388) |
| No stale v5.1.0 fallback | ✓ (updated to v5.1.1) |
| No stale v5.0.x references | ✓ (none in installer) |
| macOS detection | ✓ (darwin) |
| Linux detection | ✓ (linux) |
| Windows detection | ✓ (install.ps1) |
| Whitespace check on URLs | ✓ (line 391) |
| Binary name prefix check `ASF-` | ✓ (line 394) |

## Install Guide (`INSTALL.md`)

| Check | Status |
|---|---|
| Root `INSTALL.md` updated to v5.1.1 | ✓ |
| `release/INSTALL_v5.1.1.md` created | ✓ |
| No stale v5.1.0 references in install guides | ✓ |
| All 5 platform download URLs correct | ✓ |
| Quick install curl command correct | ✓ |
| Upgrade command correct | ✓ |
| Clean install command correct | ✓ |

**INSTALLER_AUDIT_PASSED**
