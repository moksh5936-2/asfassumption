# Supported Versions

## Current Release

| Version | Release Date | Support Level |
|---------|-------------|---------------|
| 1.0.0   | June 2026   | Active Development |

## Version Support Policy

| Version | Status |
|---------|--------|
| Latest minor (1.x) | Security fixes, critical bug fixes |
| Previous minor (0.x) | No longer supported |

## Platform Support

| Platform | Architecture | Status |
|----------|-------------|--------|
| macOS 12+ | arm64 | ✅ Supported |
| macOS 12+ | amd64 | ✅ Supported |
| Linux (glibc 2.28+) | amd64 | ✅ Supported |
| Linux (glibc 2.28+) | arm64 | ✅ Supported |
| Windows 10+ | amd64 | ⚠️ Experimental |

## Dependency Support

| Dependency | Version | Support |
|-----------|---------|---------|
| Go | 1.24+ | ✅ Build |
| Python | 3.8+ | ✅ Runtime (required) |
| Ollama | Latest | ⚠️ Optional (AI) |
| Tesseract | 5.x | ⚠️ Optional (OCR) |

## Release Cadence

- **Patch releases**: As needed for bug fixes and security issues
- **Minor releases**: Quarterly for feature additions
- **Major releases**: When breaking changes occur
