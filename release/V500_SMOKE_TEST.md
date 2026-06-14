# Smoke Test — ASF0 v5.0.0

## Results

| Test | Result |
|------|--------|
| `--version` | ✅ `ASF0 v5.0.0` |
| `--help` | ✅ All usage sections display |
| `doctor` | ✅ System diagnostic runs clean |
| TUI launch (3s) | ✅ No crash, renders normally |

## Version Output
```
ASF0 v5.0.0
```

The "A newer version (v4.0.2) is available" message is expected — it compares against the latest published GitHub release (v4.0.2). This will resolve after v5.0.0 is published.

## Doctor Output
Confirmed: OS (darwin/arm64), version (5.0.0), paths (config, cache, data, engine all present).
