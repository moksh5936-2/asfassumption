# ASF Build System

> Version: 1.0.0 | Updated: June 2026

## Overview

ASF is a Go application with no dynamic linking. The build process produces a single statically-linked binary for each target platform. The Python ASF engine is a separate package installed via `pip`.

## Prerequisites

### Required

| Tool | Version | Purpose |
|------|---------|---------|
| Go | 1.24+ | Compilation |
| Python | 3.8+ | ASF engine (runtime, not build-time) |

### Optional

| Tool | Purpose |
|------|---------|
| `goreleaser` | Automated multi-platform releases |
| `docker` | Containerized cross-compilation |
| `upx` | Binary compression (not recommended for security tooling) |

## Build Commands

### Standard Build

```bash
cd asf-tui
go build -o asf-tui .
```

### Version Build (with embedded version)

```bash
go build -ldflags="-X 'main.version=1.0.0'" -o asf-tui .
```

### Cross-Compilation

| Platform | Command |
|----------|---------|
| Linux AMD64 | `GOOS=linux GOARCH=amd64 go build -o asf-linux-amd64 .` |
| Linux ARM64 | `GOOS=linux GOARCH=arm64 go build -o asf-linux-arm64 .` |
| macOS Intel | `GOOS=darwin GOARCH=amd64 go build -o asf-darwin-amd64 .` |
| macOS ARM | `GOOS=darwin GOARCH=arm64 go build -o asf-darwin-arm64 .` |
| Windows AMD64 | `GOOS=windows GOARCH=amd64 go build -o asf-windows-amd64.exe .` |

### Testing

```bash
# Run all tests
go test ./...

# With race detection
go test -race ./...

# With coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

### Code Quality

```bash
# Vet
go vet ./...

# Format
go fmt ./...

# Lint (requires golangci-lint)
golangci-lint run ./...
```

## Build Output

```
asf-tui/
  asf-tui              # darwin/arm64 binary (11.9MB)
release/
  asf-darwin-arm64     # Release binary (copied from build)
  install.sh           # Installer script
  VERSION              # Version manifest
  checksums.txt        # SHA-256 checksums
```

## Automation Scripts

Two release automation scripts are provided:

- `scripts/build-release.sh` — Unix/macOS release build
- `scripts/build-release.ps1` — Windows PowerShell release build

## Release Workflow

```
1. Update version in asf-tui/license.go
2. Run build-release.sh or build-release.ps1
3. Verify checksums
4. Test binary: ./asf-darwin-arm64 --version
5. Create GitHub release
6. Upload artifacts
7. Update release/README.md
```

## Docker Build (Future)

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY asf-tui/ .
RUN go build -o asf-tui .

FROM alpine:3.19
COPY --from=builder /app/asf-tui /usr/local/bin/asf
CMD ["asf"]
```

## Notes

- The binary size (~11.9MB) comes from Bubble Tea and the Go runtime. UPX compression is not recommended for security applications as it interferes with code signing and anti-virus scanning.
- The `go.sum` file must be committed to the repository for reproducible builds.
- Build tags are not currently used. All features compile into every binary.
