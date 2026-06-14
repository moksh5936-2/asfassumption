# Release Notes — ASF0 v5.0.0

## Overview

ASF0 v5.0.0 is a major release focused on TUI architecture stabilization, codebase cleanup, and the restoration of Local AI as a first-class feature. This release consolidates four iterations of TUI development (v4.0.0–v4.0.2) into a polished, production-ready experience.

## Major Changes

### Local AI Restored as Core Tab
- Local AI is now a first-class sidebar tab (previously only accessible via filebrowser workflow)
- Dedicated `localai.go` module with its own view, update, and event handling
- Full Ollama integration with model selection, chat interface, and streaming responses
- AI section added to sidebar navigation alongside CASES, WORK, and SYSTEM

### TUI Architecture Improvements
- Simplified app state management with consolidated model types
- Router system centralized for keyboard-driven navigation
- Sidebar redesigned with clear visual hierarchy and icon indicators
- Bottom status bar shows active context and version info
- Dashboard removed (functionality absorbed into main workspace views)
- Startup screen removed (direct entry to workspace)
- Filebrowser consolidated into modal overlay pattern

### Codebase Cleanup
- 1,989 lines removed, 1,694 added (net −295 lines)
- Regression tests refactored for maintainability
- Removed dead code: dashboard.go, startup.go, standalone filebrowser.go
- Styles extracted and centralized in styles.go
- Validation logic simplified and deduplicated
- Export, review, analyze, results, and help modules all streamlined

### Build & Release
- Version bumped from v4.0.2 to v5.0.0
- All binaries built with `CGO_ENABLED=0 -trimpath -ldflags="-s -w"`
- 5 platform binaries (darwin-arm64, darwin-amd64, linux-amd64, linux-arm64, windows-amd64)
- SHA256 checksums provided for all binaries

## Fixed Issues
- Local AI tab now accessible directly from sidebar
- Nil pointer panics guarded across view rendering
- Filebrowser modal properly scoped to context
- Keyboard navigation routing conflicts resolved
- Test suite hardened and all tests passing

## Breaking UX Changes
| Before | After |
|--------|-------|
| Dashboard on launch | Direct workspace entry |
| Standalone filebrowser | Modal file selection |
| Startup screen | Removed (faster launch) |
| Local AI in filebrowser | Dedicated sidebar tab |

## Installation

### macOS (Apple Silicon / Intel)
```bash
curl -sfL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash
```

### Linux (AMD64 / ARM64)
```bash
curl -sfL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash
```

### Windows
```powershell
powershell -c "irm https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.ps1 | iex"
```

Or download the binary directly from the [release assets](https://github.com/moksh5936-2/asfassumption/releases/tag/v5.0.0).

## Upgrade Notes
- Run `asf doctor` after upgrading to verify your installation
- Configuration and cache are preserved across versions
- No breaking changes to analysis results or export formats

## Known Limitations
- Local AI requires Ollama running locally (no bundled model)
- TUI requires a terminal with 24-bit color support and ≥80×24 dimensions
- Windows users may need Windows Terminal or ConEmu for full TUI support
