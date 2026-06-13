# TUI Audit — ASF v3 Pre-Rebuild

## Current Screens (10 views)

| View | File | Lines | Description |
|------|------|-------|-------------|
| startupView | startup.go | 154 | Welcome screen with fox art, 6-item menu |
| dashboardView | dashboard.go | 106 | System status + Quick Actions menu |
| analyzeView | analyze.go | 320 | Document/evidence path entry, mode selection, start button |
| resultsView | results.go | 1673 | 33-section scrollable list of analysis results |
| localaiView | localai.go | 580 | Ollama model management (catalog, install, delete, activate) |
| settingsView | settings.go | 244 | 10-setting editor with cycle values |
| aboutView | about.go | 72 | Version, license, description, keyboard reference |
| exportView | export.go | ~100 | Export format selection + confirmation |
| reviewView | review.go | 271 | Browse/detail review of assumptions with status tagging |
| validationView | validation.go | 213 | Developer validation view with full traceability |

## Broken Scrolling Areas

1. **Shared viewport model**: A single `viewport.Model` (`m.vp`) is shared across ALL views. Content is set via `m.vp.SetContent(content)` in the main `View()` method (app.go:302). Keyboard scroll keys manually adjust `m.vp.YOffset` (app.go:155-173) instead of using the viewport's built-in `LineDown`/`LineUp`/`ViewDown`/`ViewUp` methods.

2. **Viewport height off-by-one**: `availLines = m.height - helpLines - 2` (app.go:295) may clip one line of content or leave a gap.

3. **Results section list** (results.go:1673): All 33 sections are rendered as a single lipgloss-joined string. There is NO viewport per section. When a section is "expanded", its content is concatenated into one giant string. The only scrolling mechanism is the shared main viewport. This causes:
   - Content clipped by borders
   - Long sections force the entire page to scroll
   - No per-section scroll state
   - No scroll indicator (e.g., "Line 40-80 / 320")

4. **Trust chains, SPOFs, verification, contradictions** are all rendered as joined strings inside the results content — no individual viewports.

5. **Review model** (review.go): No viewport in detail mode; long assumption descriptions will overflow.

6. **Validation detail** (validation.go:93-190): Full traceability rendered as `strings.Builder` then boxed — no viewport. Overflow content is clipped.

## Features Missing from TUI

- **File explorer**: Must type paths manually (analyze.go:111-143)
- **Search/filter**: None in any screen
- **Help screen**: Only `renderHelp()` bottom bar with per-view keys
- **Recent files**: Config tracks them but no TUI menu for quick re-analysis
- **Run analysis from file selection**: Must type path, then go to analyze screen
- **Export directory picker**: Must type path
- **Per-section scroll indicators**: None
- **Empty states**: Some exist (results.go:93-98, validation.go:32-37) but missing from many sections
- **Error states**: Most parse errors surface as raw error strings or status messages
- **"Terminal too small" warning**: None
- **Content clipping detection**: None

## CLI-Only Workflows

1. `asf analyze <file>` — CLI analysis bypasses TUI entirely
2. `asf doctor` — Diagnostics not in TUI
3. `asf --version-check` — Not in TUI
4. `asf --license` — License display not in TUI
5. CSV/PDF export: Only available through CLI export path
6. Narrative report generation: Only through CLI

## Raw Log Leakage Points

1. **main.go:52**: `asfLog.Printf("ASF v%s starting", ASFVersion)` — writes to log file, OK
2. **main.go:115**: `asfLog.Printf("config path: %s", asfConfigPath())` — OK
3. **main.go:118**: `debugLog.Printf("runtime dirs: %v", err)` — OK
4. **localai.go:251,377,406**: `debugLog.Printf(...)` — OK, goes to file
5. **config.go**: Load/Save use `asfLog/errorLog` — OK
6. **engine.go**: Uses `asfLog/errorLog` — OK
7. **No stdout/stderr output in TUI mode**: Clean
8. **Risks**: If any engine code calls `fmt.Print` or `log.Print`, it would corrupt the TUI display

## Resize Issues

1. **WindowSizeMsg** sets `m.width` and `m.height` (app.go:108-109) but does NOT update any sub-model viewport dimensions
2. Viewport height is recalculated in `View()` method (app.go:295) — OK but happens every frame
3. Width used in layouts but no sub-model re-layouts on resize
4. `m.styles = NewStyles(...)` re-created on resize — OK
5. **Potential crash**: If height < minimum (~6 lines for help + padding), `availLines` could become 0 or negative (clamped at 1)

## Navigation Issues

1. `q` only works from `startupView` — doesn't quit from dashboard, settings, etc.
2. No direct keyboard shortcuts for many sections (e.g., no key for Results from Dashboard)
3. `esc` handler calls `m.vp.YOffset = 0` which resets scroll — resets scroll when going back
4. No breadcrumb or path indicator
5. No tab-based navigation between result sections
6. Review/validation mode string matching: `reviewMode`, `mode` must exactly match strings

## Content Clipping Issues

1. **Border clipping**: `BorderBox` style (styles.go:163-168) has `Padding(1, 2)` which reduces available content area
2. **Results expanded content**: When a section is expanded, its entire content is rendered inline — if content is taller than viewport, bottom is clipped with no scroll indicator
3. **Assumption descriptions**: Truncated at 50-60 chars in list views
4. **Trust chains**: Long dependency paths clipped by terminal width
5. **No scroll percentage indicator**: User can't tell if there's more content

## Key Findings Summary

| Issue | Severity | Impact |
|-------|----------|--------|
| Single shared viewport | Critical | Scrolling broken across all screens |
| No file explorer | Critical | Users must type paths manually |
| Results 1673-line monolith | Critical | Unmaintainable, content clipped |
| No per-section viewports | High | Content overflow clipped |
| No help screen | Medium | Users must guess shortcuts |
| No search/filter | Medium | Cannot find specific results |
| No scroll indicators | Medium | Users don't know content length |
| Resize does not update viewports | Medium | Layout corruption possible |
| No terminal-too-small check | Low | Crash on tiny terminals |
| Logs: engine may leak to stdout | Low | Untested |
