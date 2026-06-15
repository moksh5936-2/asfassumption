# V512_VERSION_AUDIT — ASF0 v5.1.2

## Version Bump: 5.1.1 → 5.1.2

| File | Line | Change | Status |
|---|---|---|---|
| `asf-tui/license.go` | 18 | `ASFVersion = "5.1.1"` → `"5.1.2"` | ✅ Done |
| `release/VERSION` | 1 | `5.1.0` → `5.1.2` | ✅ Done |
| `install.sh` | 372 | `print('v5.1.1')` → `print('v5.1.2')` | ✅ Done |
| `install.sh` | 373 | `echo "v5.1.1"` → `echo "v5.1.2"` | ✅ Done |
| `install.sh` | 376 | `LATEST_VERSION="v5.1.1"` → `LATEST_VERSION="v5.1.2"` | ✅ Done |
| `INSTALL.md` | 1 | `v5.1.1` → `v5.1.2` (all) | ✅ Done |
| `README.md` | 40+ | `v5.1.0` → `v5.1.2` (all) | ✅ Done |

## Verification
```bash
go run . --version
# Expected: ASF0 v5.1.2
```
