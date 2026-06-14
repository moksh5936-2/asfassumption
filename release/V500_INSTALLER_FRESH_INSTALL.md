# Fresh Install Test — ASF0 v5.0.0

## Procedure
Executed the production installer from raw.githubusercontent.com:
```bash
curl -sfL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash
```

## Results
| Step | Result |
|------|--------|
| Detected existing v4.0.1 | ✅ Correct |
| Prompted for `--upgrade` | ✅ Correct |
| Upgrade with `--upgrade` | ✅ Success |
| v5.0.0 downloaded | ✅ 12M from GitHub |
| SHA-256 checksum verified | ✅ Match |
| Binary verified | ✅ `ASF0 v5.0.0` |
| PATH configured | ✅ `~/.local/bin/asf` symlink |
| `asf --version` | ✅ `ASF0 v5.0.0` |
| `asf doctor` | ✅ Native engine active |
| TUI launch (2s) | ✅ No crash |

## Conclusion
✅ Fresh install/upgrade pipeline confirmed working end-to-end for v5.0.0.
