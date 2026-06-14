# ASF0 v5.0.2

## Overview

ASF0 v5.0.2 is a focused bug-fix release that resolves critical issues with the modal file picker across macOS, Linux, and Windows. Users can now freely navigate the filesystem to select architecture and evidence files without being locked to the `/reports` directory. This release also includes TUI stability improvements from v5.0.1.

## Fixed Issues

- **Fixed modal file picker cross-platform navigation** — Rewrote path handling using Go `filepath` utilities instead of manual string concatenation; works on macOS, Linux, and Windows.
- **Fixed macOS file selector behavior** — File picker can now navigate `~`, `~/Desktop`, `~/Downloads`, `/`, `/Volumes`, and paths with spaces. macOS `$TMPDIR` is respected.
- **Fixed picker being locked to `/reports`** — Start directory now follows a priority chain: last-used → document directory → current working directory → home directory. No hardcoded paths.
- **Added free filesystem traversal** — Navigation keys: `~` (home), `g` (root), `d` (Downloads), `D` (Desktop), `r` (refresh), `Backspace` (parent), `/` (search), `.` (hidden files).
- **Improved architecture file selection** — Mode-specific extension filtering (`.yaml`, `.yml`, `.json`, `.md`, `.mmd`, `.drawio`, `.svg`, `.pdf`, `.docx`, `.txt`).
- **Improved evidence file selection** — Mode-specific extension filtering (`.csv`, `.json`, `.yaml`, `.yml`, `.txt`, `.md`, `.pdf`, `.docx`); duplicate prevention.
- **Preserved case-first ASF0 workflow** — + New Analysis → select architecture → add evidence → analyze → case-first result view.
- **Preserved Local AI sidebar support** — AI section with model management remains intact.

## Installation

### macOS / Linux (curl)
```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash
```

### macOS (Homebrew, future)
```bash
brew install asf
```

### Windows
Download `ASF-v5.0.2-windows-amd64.exe` from the GitHub release page.

### Upgrade existing
```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash -s -- --upgrade
```

## Upgrade Notes

- No breaking changes.
- Existing configuration at `~/.asf/config.yaml` is preserved.
- Existing evidence files and case data are preserved.
- If you were affected by the `/reports` lock, upgrade removes the restriction immediately.

## Known Limitations

- Windows drive switching is supported but requires manual path entry for drives beyond the current volume.
- Manual TUI acceptance (visual verification in a real terminal) is recommended after upgrade.
