# TUI Real TTY UX Audit

## Environment

| Attribute | Value |
|-----------|-------|
| Platform  | macOS (Darwin) arm64 |
| Terminal  | PTY (Perl IO::Pty) at 120×40 |
| Build     | ASF v4.0.2 (CGO_ENABLED=0, -trimpath, -ldflags="-s -w") |
| Shell     | zsh |
| Date      | 2026-06-14 |

## Audit Method

Real PTY-based manual interaction using Perl IO::Pty module. Each key press is sent via `syswrite` to the PTY master, output is read from the PTY master after ANSI escape sequence stripping.

## Dashboard

**Status: PASS**

| Check | Result |
|-------|--------|
| Dashboard renders on startup | PASS |
| "Welcome to ASF" displayed | PASS |
| Version "ASF v4.0.2" in top bar | PASS |
| Sidebar rendered (18 items) | PASS |
| Bottom bar with context hints | PASS |
| Fox art displayed | PASS |
| Quick actions menu rendered | PASS |

The Dashboard shows System Status (Version, Mode, AI, Theme), Quick Actions (Analyze Architecture, Local AI Models, Settings, About), and a vertical sidebar with all 18 navigation items.

## Sidebar Navigation

**Status: PASS**

| Check | Result |
|-------|--------|
| All 18 sidebar items rendered | PASS |
| Tab cycles forward | PASS |
| Shift+Tab cycles backward | PASS |
| Down arrow selects next item | PASS |
| Up arrow selects previous item | PASS |
| vim j/k navigation | PASS |
| Selection highlight via color | PASS |

The sidebar includes: Dashboard, File Explorer, Analyze, Summary, Assumptions, Verification, Contradictions, Trust Chains, Single Points of Trust, Assumption Impact Analysis, Blind Spots, SDRI, Recommended Controls, Security Design Review, Reports / Exports, Settings, Help, About.

Tab immediately activates the selected view (auto-navigate). Arrow keys move selection within the current view's local menu.

## Help Screen

**Status: PASS**

| Check | Result |
|-------|--------|
| `?` opens Help | PASS |
| Keyboard shortcuts displayed | PASS |
| All sections rendered | PASS |
| `q` returns to previous view | PASS |
| `Esc` returns to previous view | PASS |

Help screen includes sections: Global, Navigation, Dashboard, Analyze, Results, File Explorer, AI Models, Review Mode, Export, Settings, Search, Sidebar Navigation, and Supported File Types.

## File Explorer

**Status: PASS**

| Check | Result |
|-------|--------|
| `f` opens File Explorer | PASS |
| Sidebar item opens File Explorer | PASS |
| Directory listing shown | PASS |
| Esc returns to previous view | PASS |
| Context hints in bottom bar | PASS |

The File Explorer opens in `./reports` by default (config output directory). Navigates with arrow keys, supports Tab for preview, `.` for hidden files toggle.

## Search

**Status: PASS (fixed)**

| Check | Result |
|-------|--------|
| `/` activates search bar | PASS |
| Search bar rendered in content | PASS |
| Query text visible while typing | PASS |
| Enter/Esc closes search | PASS |
| n/N scrolls viewport | PASS |
| Bottom bar shows /=Search on Results | PASS |

Search bar was previously invisible - fixed in v4.0.2. Now renders a `/search>_` prompt at the top of the content area when active.

## Analyze View

**Status: PASS**

| Check | Result |
|-------|--------|
| `r` opens Analyze | PASS |
| Document Path field | PASS |
| Evidence Path field | PASS |
| Mode selection (ASF Only, ASF+AI) | PASS |
| Start Analysis action | PASS |
| Progress bar during analysis | PASS |
| Esc cancels running analysis | PASS |

## Results View

**Status: PASS**

| Check | Result |
|-------|--------|
| Analysis complete switches to Results | PASS |
| Tab bar with all 12 result tabs | PASS |
| Summary tab renders | PASS |
| Assumptions tab renders | PASS |
| Verification tab renders | PASS |
| Contradictions tab renders | PASS |
| Trust Chains tab renders | PASS |
| Single Points of Trust tab renders | PASS |
| Impact Analysis tab renders | PASS |
| Blind Spots tab renders | PASS |
| SDRI tab renders | PASS |
| Controls tab renders | PASS |
| Security Design Review tab renders | PASS |
| Reports/Exports tab renders | PASS |
| Tab/Shift+Tab cycles through tabs | PASS |
| Sidebar entries map to correct tabs | PASS |

## Scrolling

**Status: PASS**

| Check | Result |
|-------|--------|
| j/k scroll long content | PASS |
| Up/Down scroll long content | PASS |
| PageUp/PageDown (b/Space) | PASS |
| Home/End (g/G) | PASS |
| Mouse wheel | PASS |
| Scroll percentage indicator | PASS |
| Scroll position preserved per view | PASS |
| No random scroll reset | PASS |

## Key Bindings

**Status: PASS**

| Key | Action | Result |
|-----|--------|--------|
| Ctrl+C / Q | Force quit | PASS |
| q | Back / navigate previous view | PASS |
| Esc | Back / cancel | PASS |
| ? | Toggle help | PASS |
| Tab | Cycle sidebar (activate) | PASS |
| Shift+Tab | Cycle sidebar reverse | PASS |
| f | Open file explorer | PASS |
| r | Open Analyze / Review mode | PASS |
| / | Start search (Results only hint) | PASS |
| ↑/k | Move up / Scroll up | PASS |
| ↓/j | Move down / Scroll down | PASS |
| PgUp/b | Page up | PASS |
| PgDn/Space | Page down | PASS |
| Home/g | Go to top | PASS |
| End/G | Go to bottom | PASS |
| e | Export (results view) | PASS |
| c | Clear results | PASS |
| s | Save settings | PASS |
| a/l/s/i | Dashboard quick actions | PASS |
| Number keys | Open recent file | PASS |

## Terminal Resize

**Status: PASS**

The TUI handles WindowSizeMsg correctly, recalculating layout on resize.

## Empty States

**Status: PASS**

All empty states render cleanly:
- No analysis run
- No assumptions
- No contradictions
- No trust chains
- No file selected

## Issues Found

### Resolved in v4.0.2

1. **Search bar invisible** → Fixed: Search bar now renders in content area
2. **Dead focusManager code** → Removed
3. **Scroll reset on analysis complete** → Fixed: uses navigateTo pattern
4. **Trust Chains/SPOFs combined** → Split into separate tabs
5. **Sidebar missing Summary, About** → Added
6. **Browse hints showed /=Search globally** → Results-only now
7. **Log directory not created** → Fixed: creates ~/.asf/logs/

### Remaining (minor)

1. **File browser default path** `./reports` directory doesn't exist until first export
2. **n/N search** does viewport scrolling, not semantic search (intentional, search is text highlight navigation)
3. **Results tab cycling via Tab** limited to results view; sidebar Tab returns to results view
