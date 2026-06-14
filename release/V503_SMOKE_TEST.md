# V504 — Native Binary Smoke Test

| Check | Result |
|---|---|
| `--version` shows `ASF0 v5.0.4` | ✅ |
| `--help` shows usage | ✅ |
| `doctor` shows version 5.0.4 | ✅ |
| `doctor` detects existing installations | ✅ |
| `doctor` OS detection (darwin/arm64) | ✅ |
| `doctor` paths/config writable | ✅ |
| TUI launch | ⚠️ Requires real terminal — automated smoke test deferred |

## Details

- Binary: `ASF-v5.0.4-darwin-arm64`
- Version output: `ASF0 v5.0.4`
- Doctor: All systems operational, all paths writable, no crashes
- Binary size: 12 MB (darwin-arm64)

**Status: SMOKE_OK**
