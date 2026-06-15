# V510_VERSION_AUDIT — ASF0 v5.1.0 Version Bump

## Files Updated (5.0.5 → 5.1.0)

| File | Change |
|------|--------|
| `asf-tui/license.go:18` | `"5.0.5"` → `"5.1.0"` |
| `asf-tui/Makefile:4` | `5.0.5` → `5.1.0` |
| `asf-tui/install.sh:15` | `"5.0.5"` → `"5.1.0"` |
| `scripts/build-release.sh:9` | `5.0.5` → `5.1.0` |
| `install.sh:372` | `print('v5.0.5')` → `print('v5.1.0')` |
| `install.sh:373` | `echo "v5.0.5"` → `echo "v5.1.0"` |
| `install.sh:376` | `LATEST_VERSION="v5.0.5"` → `LATEST_VERSION="v5.1.0"` |
| `release/install.sh:372` | `print('v5.0.5')` → `print('v5.1.0')` |
| `release/install.sh:373` | `echo "v5.0.5"` → `echo "v5.1.0"` |
| `release/install.sh:376` | `LATEST_VERSION="v5.0.5"` → `LATEST_VERSION="v5.1.0"` |
| `release/VERSION` | `5.0.5` → `5.1.0` |
| `release/INSTALL.md` | All v5.0.5 → v5.1.0 |
| `README.md` | All v5.0.0 → v5.1.0 (multiple lines) |

## Verification

```
$ go run . --version
ASF0 v5.1.0
```

All version references updated.
