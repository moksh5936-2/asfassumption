# ASF0 v5.0.3 — Build Validation

## Status: CI-Deferred

Go is not installed on the local release engineering system. Build validation is deferred to the GitHub Actions CI pipeline on tag push.

## CI Pipeline (`.github/workflows/ci.yml`)

| Job | Trigger | Status |
|-----|---------|--------|
| lint | tag push `v*` | Runs `golangci-lint` |
| build | tag push `v*` | Runs `go build -v -o asf .` + `file/asf --version` |
| test | tag push `v*` | Runs `go test ./... -v -count=1 -timeout 120s` |
| release | tag push `v*` | Builds matrix binaries, checksums, completions, creates release |

## Binary Build Flags (CI release job)
- `CGO_ENABLED=0`
- `-trimpath`
- `-ldflags="-s -w -X main.ASFVersion=${GITHUB_REF_NAME#v}"`

## Binary Matrix (CI)
| OS | Arch | Output |
|----|------|--------|
| linux | amd64 | `ASF-v5.0.3-linux-amd64` |
| linux | arm64 | `ASF-v5.0.3-linux-arm64` |
| darwin | amd64 | `ASF-v5.0.3-darwin-amd64` |
| darwin | arm64 | `ASF-v5.0.3-darwin-arm64` |
| windows | amd64 | `ASF-v5.0.3-windows-amd64.exe` |

## Verification Required Post-CI
- [ ] `go fmt ./...` passes
- [ ] `go vet ./...` passes
- [ ] All tests pass: `go test -count=1 ./...`
- [ ] Semantic tests pass: `go test -count=1 -run Semantic ./...`
- [ ] Build succeeds: `go build ./...`
