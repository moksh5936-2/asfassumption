# TUI Information Architecture — ASF v3 Rebuild

## Screen Hierarchy

```
┌─ Startup [one-shot welcome] ─────────────────┐
│  menu: Analyze | Results | AI Models |       │
│        Settings | About | Exit               │
│  → Dashboard (Esc or Enter)                  │
└──────────────────────────────────────────────┘

┌─ Dashboard [primary landing] ────────────────┐
│  System Status (version, mode, AI, theme)    │
│  Quick Actions → Analyze | AI Models |       │
│                  Settings | About            │
│  Recent Files list (future)                  │
│  → Any action screen via Enter/keys          │
│  → File Explorer (f)                        │
│  → Help (?)                                 │
└──────────────────────────────────────────────┘

┌─ Analyze [analysis setup] ───────────────────┐
│  Document Path (text input)                  │
│  Evidence Path (text input)                  │
│  Mode: ASF-Only / ASF+AI                    │
│  ► Start Analysis                           │
│  → Results (on analysis complete)           │
│  → File Explorer (f)                        │
└──────────────────────────────────────────────┘

┌─ Results [9-tab output] ─────────────────────┐
│  Tab: Summary | Assumptions | Verification  │
│        Contradictions | Trust | Impact       │
│        Blind Spots | Controls | Reports     │
│  Each tab has scrollable content             │
│  → Export (e) → Export sub-screen           │
│  → Review (r) → Review sub-screen           │
│  → Validation (v) → Validation sub-screen   │
└──────────────────────────────────────────────┘

┌─ File Explorer [MC-style file picker] ───────┐
│  Path bar | file list (dirs + files)         │
│  Supported file filtering                    │
│  Hidden file toggle                          │
│  Search mode (/ key)                         │
│  → Analyze (Enter on supported file)        │
│  → Dashboard (Esc)                           │
└──────────────────────────────────────────────┘

┌─ AI Models [Ollama manager] ─────────────────┐
│  Ollama status bar                           │
│  Recommended ASF Models catalog              │
│  Other installed models                      │
│  Per-model actions: Download | Activate |    │
│                    Delete                     │
│  → Dashboard (Esc)                           │
└──────────────────────────────────────────────┘

┌─ Settings [config editor] ───────────────────┐
│  Theme | Fox Style | Analysis Depth |        │
│  Risk Threshold | STRIDE | Controls |        │
│  Export Format | Export Directory |          │
│  AI Enhancement | Active Model               │
│  → Dashboard (Esc)                           │
└──────────────────────────────────────────────┘

┌─ Help [keyboard reference] ──────────────────┐
│  Global | Navigation | Analyze | Results |   │
│  File Explorer | Review | Export | Settings  │
│  Search | File Types                         │
│  → Previous screen (Esc)                     │
└──────────────────────────────────────────────┘

┌─ About [version/license] ────────────────────┐
│  Version | License | Description | Tech      │
│  Keyboard summary                            │
│  → Previous screen (Esc)                     │
└──────────────────────────────────────────────┘

── Sub-screens (contextual, non-sidebar) ──────

┌─ Export [format/output selection] ───────────┐
│  Format list (json, md, html, csv, pdf,      │
│             narrative, trust, coverage, etc)  │
│  Output path                                 │
│  Confirm (y)                                 │
│  → Results (Esc)                             │
└──────────────────────────────────────────────┘

┌─ Review [assumption review] ─────────────────┐
│  Mode: Browse (list) / Detail (single)       │
│  Status tagging: Accept | Reject | Modified  │
│  Note editing on each assumption             │
│  → Results (Esc)                             │
└──────────────────────────────────────────────┘

┌─ Validation [dev traceability] ──────────────┐
│  Mode: Browse (list) / Detail (full trace)   │
│  Shows evidence, STRIDE, risk, confidence    │
│  → Results or Review (Esc)                  │
└──────────────────────────────────────────────┘
```

## Navigation Paths

| From | To | Trigger |
|------|----|---------|
| Startup | Dashboard | `Enter` on any menu item, or `Esc` |
| Startup | Any screen | `Enter` on menu item directly |
| Dashboard | Analyze | Select "Analyze Architecture" or press `a` |
| Dashboard | AI Models | Select "Local AI Models" or press `l` |
| Dashboard | Settings | Select "Settings" or press `s` |
| Dashboard | About | Select "About" or press `i` |
| Dashboard | File Explorer | Press `f` |
| Dashboard | Help | Press `?` |
| Analyze | Results | Analysis complete (auto-navigate) |
| Analyze | File Explorer | Press `f` for path selection |
| Results | Export | Press `e` |
| Results | Review | Press `r` |
| Results | Validation | Press `v` |
| File Explorer | Analyze | `Enter` on supported file |
| Any screen | Help | Press `?` |
| Any screen | Previous | `Esc` (uses view history stack) |
| Any non-startup | Quit | `Ctrl+C` or `Q` |

## Sidebar Design

**Items** (in order, always shown when sidebar is open):

| # | Item | View | Key |
|---|------|------|-----|
| 1 | Dashboard | dashboardView | `Tab` cycle |
| 2 | Analyze | analyzeView | `Tab` cycle |
| 3 | Results | resultsView | `Tab` cycle |
| 4 | File Explorer | fileBrowserView | `Tab` cycle |
| 5 | AI Models | localaiView | `Tab` cycle |
| 6 | Settings | settingsView | `Tab` cycle |
| 7 | About | aboutView | `Tab` cycle |
| 8 | Help | helpView | `Tab` cycle |

- Quit removed from sidebar (Ctrl+C/Q always works)
- Sidebar default open, 23 chars wide
- Active view highlighted with accent color
- `Tab` / `Shift+Tab` cycles sidebar selection AND navigates to that view
- `Esc` returns to previous view via history stack

## Layout Regions

```
┌── Top Bar (1 line) ──────────────────────────┐
│ ASF v3.4  │  /path/to/file.yaml  │  ready    │
├──┬───────────────────────────────────────────┤
│  │                                           │
│S │          Main Panel                       │
│i │      (per-view viewport)                  │
│d │                                           │
│e │                                           │
│b │                                           │
│a │                                           │
│r │                                           │
├──┴───────────────────────────────────────────┤
│  Tab:Nav • Esc:Back • ?:Help • Scroll 50%   │
└──────────────────────────────────────────────┘
```

| Region | Height | Content |
|--------|--------|---------|
| Top Bar | 1 line | `version │ current_file │ status_msg` |
| Sidebar | full height - 3 | 8 items with active highlight |
| Main Panel | full height - 3 | Per-view viewport |
| Bottom Bar | 1 line | Context-sensitive key hints + scroll% |

## Viewport Strategy

Each major view gets its **own** `viewport.Model` instance:

| View | Viewport Required | Reason |
|------|------------------|--------|
| Dashboard | Yes | Status + quick actions can exceed height |
| Analyze | Yes | Path inputs + modes + start button |
| Results | Yes + per-tab scroll tracking | 9 tabs, each can be long |
| File Explorer | Yes | Directory listing can be long |
| AI Models | Yes | Catalog + other models + actions |
| Settings | Yes | 10 settings items |
| Help | Yes | Keyboard reference sections |
| About | Maybe | Short content, but consistent |
| Export | Yes | Format list + output path + confirmation |
| Review | Yes | Browse list + detail content |
| Validation | Yes | Browse list + full trace content |

**Implementation**: Instead of shared `m.vp`, each sub-model holds its own `viewport.Model`. Main layout only sizes the main panel area. Per-view render functions return the view's viewport view.

## Unused Screens to Remove

These CLI-only workflows are NOT mapped to TUI and will remain CLI-only:
- `asf doctor` — CLI diagnostics
- `asf --version-check` — CLI version check
- `asf --license` — CLI license display
- CSV/PDF export via CLI (TUI uses Export screen instead)

## Keyboard Shortcut Master Plan

| Key | Scope | Action |
|-----|-------|--------|
| `Ctrl+C`, `Q` | Global | Force quit |
| `?` | Global | Toggle help overlay |
| `Esc` | Global | Go back (history stack) |
| `Tab` | Global | Next sidebar item |
| `Shift+Tab` | Global | Previous sidebar item |
| `↑`/`k` | Global | Scroll up / Navigate up |
| `↓`/`j` | Global | Scroll down / Navigate down |
| `PgUp`/`b` | Global | Page up |
| `PgDn`/`Space` | Global (not on results) | Page down |
| `Home`/`g` | Global | Go to top |
| `End`/`G` | Global | Go to bottom |
| `Ctrl+U` | Global | Half page up |
| `Ctrl+D` | Global | Half page down |
| `f` | Global | Open file explorer |
| `/` | Global | Start search in current panel |
| `n`/`N` | Global (when search active) | Next/previous match |
| `Enter` | Screen-specific | Select / Confirm / Edit |
| `a` | Dashboard | Analyze |
| `l` | Dashboard | AI Models |
| `s` | Dashboard/Settings | Settings / Save |
| `i` | Dashboard | About |
| `Tab` | Results | Next results tab |
| `Shift+Tab` | Results | Previous results tab |
| `e` | Results | Export |
| `r` | Results | Review |
| `v` | Results/Review | Validation |
| `h` | File Explorer | Parent directory |
| `.` | File Explorer | Toggle hidden files |
| `s` | Review (detail) | Accept status |
| `r` | Review (detail) | Reject status |
| `m` | Review (detail) | Modified status |
| `n` | Review (detail) | Edit note |
| `←`/`→` | Settings (editing) | Change value |
| `y` | Export | Confirm export |
