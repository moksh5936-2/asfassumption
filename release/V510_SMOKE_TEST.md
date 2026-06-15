# V510_SMOKE_TEST — ASF0 v5.1.0 Native Binary

Binary: `dist/ASF-v5.1.0-darwin-arm64`

## CLI Smoke Tests

| Test | Result | Output |
|------|--------|--------|
| `--version` shows ASF0 v5.1.0 | ✅ | `ASF0 v5.1.0` |
| `--help` shows usage | ✅ | All flags and commands displayed |
| `doctor` reports version 5.1.0 | ✅ | `Version: 5.1.0` |
| `doctor` reports all paths writable | ✅ | Config/cache/data dirs writable |
| `doctor` detects Ollama | ✅ | v0.30.8, running, 1 model |

## Checksum Verification

```
$ shasum -a 256 -c checksums.txt
ASF-v5.1.0-darwin-amd64: OK
ASF-v5.1.0-darwin-arm64: OK
ASF-v5.1.0-linux-amd64: OK
ASF-v5.1.0-linux-arm64: OK
ASF-v5.1.0-windows-amd64.exe: OK
```

All binaries verify. Checksums match.
