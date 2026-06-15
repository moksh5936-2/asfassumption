# V511 — Release Certification

## Certification Checklist

| # | Criterion | Status | Evidence |
|---|---|---|---|
| 1 | Correct source verified | ✓ | `V511_SOURCE_VERIFICATION.md` — all features present |
| 2 | Code committed | ✓ | commit `9d75c65` — "release: ASF0 v5.1.1" |
| 3 | Tag v5.1.1 pushed | ✓ | `git push origin v5.1.1` |
| 4 | GitHub release created | ✓ | `gh release create v5.1.1` — published |
| 5 | All binaries uploaded | ✓ | 5 platform binaries: darwin-arm64, darwin-amd64, linux-amd64, linux-arm64, windows-amd64 |
| 6 | Checksums verified | ✓ | `shasum -a 256` matches `checksums.txt` on release |
| 7 | Installer installs v5.1.1 | ✓ | `install.sh` defaults to `v5.1.1`, `BINARY_NAME` correct |
| 8 | Upgrade installs v5.1.1 | ✓ | `--upgrade` flag targets `v5.1.1` |
| 9 | Binary shows intended TUI | ✓ | `ASF0 v5.1.1` confirmed via `--version` |
| 10 | Split-pane results work | ✓ | acceptance tests pass (11/11) |
| 11 | Selected item stays visible | ✓ | `ensureSelectedVisible` invariant enforced + viewport offset fix |
| 12 | No dropdown details remain | ✓ | `TestNoInlineDropdownExpansion` passes |
| 13 | Tests pass | ✓ | `go test -count=1 .` — all 19 packages pass |
| 14 | Race tests pass | ✓ | `go test -race -count=1 -short .` — PASS |
| 15 | Build passes | ✓ | `go fmt`, `go vet`, `go build` — all clean |

## Build Validation

| Command | Status |
|---|---|
| `go fmt ./...` | PASS |
| `go vet ./...` | PASS |
| `go build ./...` | PASS |
| `go test -count=1 .` | PASS (main) |
| `go test -race -count=1 -short .` | PASS |
| `go test -count=1 ./...` | PASS (all 19 packages) |

## Assets

| Asset | Size | SHA256 |
|---|---|---|
| `ASF-v5.1.1-darwin-arm64` | 12M | `2c5d9f7ed2e384b71c5933251614c1519e35c37de27ee38ee5870a35e36c7143` |
| `ASF-v5.1.1-darwin-amd64` | 13M | `701874865f4bed475929054f6b70ce9f532853e7ae641f8fef286a67c0f6e520` |
| `ASF-v5.1.1-linux-amd64` | 13M | `5e0833886a98b0d65ed63f2cf67a05771e67bdcf69fefe3e5510ef75a722c6fb` |
| `ASF-v5.1.1-linux-arm64` | 12M | `8288505be9c72870fb63ffbe505ce1b906067309d3bdaa63aac9afd91f7b03c4` |
| `ASF-v5.1.1-windows-amd64.exe` | 13M | `9a041ea07dae2fd0c7da8fd51cb11f8fb6a6706be17e29552274466d7ac453fb` |

## Remaining Risks

- Controls tab lacks `[Gap]`/`[Control]`/`[Coverage]` categorization prefixes (data model limitation, not a blocker)
- Verification tab group labels cause minor header line estimation imprecision (acceptable for small lists)
- `SearchActive` field on `tabState` is defined but currently unused directly

## Verdict

```
V511_RELEASE_CERTIFIED
```
