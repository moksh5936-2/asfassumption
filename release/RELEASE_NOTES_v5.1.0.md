# ASF0 v5.1.0 — UX Stabilization and Discoverability

## Overview

ASF0 v5.1.0 is a UX stabilization release that fixes five TUI layout and scrolling bugs identified in v5.0.5, introduces a complete onboarding experience, redesigns the Overview tab to surface critical findings, and aligns all documentation to v5.1.0. No engine behavior was changed.

## Major Changes

### Bounded Layout Engine
- Sidebar rendering now caps visible rows to available terminal height
- Content body clips to remaining space after header, breadcrumb, hints, and status bar
- Prevents overflow rendering that previously caused layout drift

### Sidebar Viewport
- Sidebar uses a three-value offset pattern (selectedIndex + viewportOffset + visibleHeight)
- MaintainSidebarOffset() recenters selection when it scrolls out of view
- Selection always visible regardless of sidebar tree depth
- Resize recalculates viewport offset to preserve selection visibility

### Selection-Follow-Scroll
- Content viewport scroll position tracked per-view
- Scroll restored when navigating back to a previously visited view
- Sidebar selection syncs with content view when navigating

### Onboarding Experience (New)
- Lazygit-style startup menu with 4 selectable items + hotkeys
- Quick Tour (7 slides, accessible via `?` on startup)
- Empty state guidance for first-time analyze view
- Enhanced empty state for results workspace
- Navigation: ↑↓ enter for menu, ← → for tour slides, q/esc to return

### Overview Redesign
- Top 3 critical/high contradictions displayed inline on Overview tab
- Top 5 Single Points of Trust Failure displayed inline
- Top 3 SDRI critical/high findings displayed inline
- Users can now see the most important findings without navigating to tabs 3-6

### Documentation Alignment
- Root `README.md` updated to v5.1.0
- `release/INSTALL.md` updated from v4.0.0 to v5.1.0
- All installer scripts reference v5.1.0 assets
- Build scripts default to v5.1.0

## Fixed Issues

- **Drifting sidebar**: Sidebar content could overflow below status bar when terminal was short
- **Disappearing selection**: Selected sidebar item scrolled out of view when tree was deep
- **Duplicate breadcrumbs**: Inner breadcrumb in viewResults() duplicated the outer breadcrumb bar
- **Inaccessible findings**: Contradictions, SPOFs, and SDRI required 3-6 tab keystrokes to discover
- **Clipped content**: Content view clipped at incorrect offset after breadcrumb removal

## Files Changed
- `asf-tui/app.go` — startup screen, tourView wiring, tour key handling, renderContent
- `asf-tui/analyze.go` — onboarding empty state for first-time users
- `asf-tui/results.go` — Overview tab surfaces contradictions/SPOFs/SDRI, enhanced empty state
- `asf-tui/router.go` — sidebar viewport with offset tracking, maintainSidebarOffset
- `asf-tui/tour.go` — new: 7-slide Quick Tour model and rendering
- `asf-tui/license.go` — version bumped to 5.1.0
- `asf-tui/Makefile` — version bumped to 5.1.0
- `asf-tui/install.sh` — updated to v5.1.0
- `README.md` — updated to v5.1.0
- `release/INSTALL.md` — updated from v4.0.0 to v5.1.0
- `install.sh` — fallback version updated to v5.1.0
- Various installer/build scripts — version references updated

## Installation

### macOS (arm64)
```bash
curl -sfLO https://github.com/moksh5936-2/asfassumption/releases/download/v5.1.0/ASF-v5.1.0-darwin-arm64
chmod +x ASF-v5.1.0-darwin-arm64
mkdir -p ~/.local/bin
cp ASF-v5.1.0-darwin-arm64 ~/.local/bin/asf
export PATH="$PATH:$HOME/.local/bin"
asf --version
# Expected: ASF0 v5.1.0
```

### Quick Install
```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash
```

## Upgrade from v5.0.x
```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash -s -- --upgrade
```

## Known Limitations

- Semantic contradiction engine detects 26 pairs across 9 categories (unchanged from v5.0.5)
- Version check warns about "newer version" until v5.1.0 is published as latest
- Quick Tour requires interactive terminal (keyboard navigation)
