# ASF0 v5.0.5 — TUI Rendering & Layout Fixes

## Overview
ASF0 v5.0.5 is a maintenance release that fixes five TUI rendering, scrolling, and layout bugs identified in v5.0.4. No engine behavior was changed.

## Fixed Since v5.0.4

1. **Viewport scrolling** — Resolved a bug where `selectedLine = 0` was being set at the end of `viewResults()`, wiping the correct line offset and preventing scroll navigation. Scroll logic now uses scroll-into-view instead of snap-to-top.

2. **Sidebar CASES section clipping** — The section rule width was hardcoded to `sidebarInnerWidth()-3` (21 chars), causing " CASES " (7 chars) to overflow the 24-char sidebar. Now calculated dynamically using `lipgloss.Width(labelStr)`.

3. **CIE contradiction display** — Contradictions detected by the semantic contradiction engine (CIE) were stored in `result.CIEContradictions` but never rendered in the TUI. A `mergeCIEContradictions()` function now converts them to the legacy format and appends them to `result.Contradictions` with dedup.

4. **Hints bar overflow gap** — When hints exceeded terminal width, `lipgloss.JoinVertical` padded all rows to the widest width, creating a visible gap. Hints are now padded to `m.width - 2` to stay within terminal bounds.

5. **Status bar 2-char gap** — The status bar fill calculation used `-4` but the bar has `Padding(0,1)`. Fixed to `-2`, making the bar exactly `m.width` wide.

## Verification
- Build: `go fmt`, `go vet`, `go build` all clean
- Tests: All 21 test packages pass (`go test -count=1 ./...`)
- Binary: Native binary verified — `ASF0 v5.0.5`
- Installer: Updated to download v5.0.5 assets

## Upgrade Notes
Users running v5.0.4 (or earlier) should upgrade to v5.0.5 to fix TUI rendering gaps, sidebar clipping, scroll navigation, and missing CIE contradiction display.

```bash
curl -fsSL https://raw.githubusercontent.com/moksh5936-2/asfassumption/main/install.sh | bash -s -- --upgrade
```

## Known Limitations
- TUI requires a terminal emulator (not testable in CI non-TTY environments)
- Version check warns about "newer version" until v5.0.5 is published
- Semantic contradiction engine detects 26 pairs across 9 categories (unchanged from v5.0.4)
