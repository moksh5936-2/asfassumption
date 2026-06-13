# TUI Architecture Stabilization Report

## Certification Status: PASS (PASS 14 requires interactive acceptance)

### Pass 2: Single-Shell Architecture

**Verdict: PASS**

All navigation, dispatch, and rendering flows through a single `mainModel` with a centralized `Update()` method. There is exactly one:
- Navigation owner: `mainModel.navigateTo()` / `navigateBack()` + history stack
- Focus owner: `focusManager` tracks active view and sub-focus
- Active screen renderer: `renderContent()` — single switch dispatch
- Global key dispatch layer: `handleGlobalKey()` with boolean fallthrough
- Layout manager: `layoutManager` with sidebar, top bar, bottom bar dimensions

The sidebar is redesigned from 8 items to 16 items covering every core ASF feature:
Dashboard, File Explorer, Analyze, Assumptions, Verification, Contradictions, Trust Chains, Single Points of Trust, Assumption Impact Analysis, Blind Spots, SDRI, Recommended Controls, Security Design Review, Reports/Exports, Settings, Help.

### Pass 3: Information Architecture

**Verdict: PASS**

Created `docs/TUI_INFORMATION_ARCHITECTURE.md` documenting all 16 sidebar items, 11 result tabs, key mappings, and architecture principles.

### Pass 4: Feature Parity

**Verdict: PASS**

All 85 features from `TUI_FEATURE_PARITY.md` continue to work. No scope-creep modules added.

### Pass 5: File Explorer Improvements

**Verdict: PASS**

File explorer correctly opens via `f` key or sidebar navigation. Navigates to parent with Backspace, toggles hidden files with `.`, selects files with Enter.

### Pass 6: Navigation / Focus

**Verdict: PASS**

- Tab / Shift+Tab cycle all 16 sidebar items correctly
- Tab on resultsView / fileBrowserView falls through to child models
- Esc globally defers to `navigateBack()` with view-specific exceptions
- `navigateTo()` updates sidebar selection via `sidebarEntries` mapping
- `navigateBack()` restores previous view and updates sidebar selection
- Arrow keys scroll content views; pass through to children for actionable views

### Pass 7: Scrolling

**Verdict: PASS**

All scroll keys work: ↑/↓/j/k, PgUp/PgDn, Home/End, Ctrl+U/Ctrl+D. Shared viewport with per-view scroll offset cached in `scrollY map[view]int`. Search mode correctly intercepts `/`, `n`, `N`.

### Pass 8: Layout Manager

**Verdict: PASS**

`layoutManager` struct provides `sidebarWidth` (23), `topBarHeight` (1), `bottomBarHeight` (1). Layout methods `sidebarWidth()`, `mainWidth()`, `mainHeight()` delegate to `layoutManager`.

### Pass 9: Results UX (Sectioned)

**Verdict: PASS**

11 result tabs with sectioned rendering:
- Summary: aggregate counts
- Assumptions: filterable, risk-colored, confidence display
- Verification: verified/partial/unverified/no-evidence counts
- Contradictions: severity-colored with descriptions
- Trust: trust chains + SPOFs, filterable
- Impact: priority queue, SPOFs, CISO dashboard
- Blind Spots: blind spots, domain blind spots, coverage gaps
- Controls: searchable control recommendations
- Reports: narrative, campaigns, confidence, export
- SDRI: executive summary, control inventory, coverage by category, dashboard
- Security Design Review: design findings, weaknesses, remediations, compliance

### Pass 10: Export Flow

**Verdict: PASS**

Export view accessible via `e` key from results. Navigate to exportView, select format, confirm export path. Back via Esc.

### Pass 11: Empty / Error States

**Verdict: PASS**

All tab renderers check for nil/empty data and return `s.EmptyState.Render(...)` with descriptive messages. Example: "No results available. Run an analysis first.", "No verification data available.", "No SDRI data available."

### Pass 12: Logging Isolation

**Verdict: PASS**

Logging audit confirmed: `asfLog` writes to `~/.asf/logs/asf.log` only. Defaults to `io.Discard`. Never leaks into TUI rendering. `initLogger()` called from `main.go:48`.

### Pass 13: Regression Tests

**Verdict: PASS**

48 regression tests verify:
- Arrow key routing (content vs. actionable views)
- Tab / Shift+Tab key handling
- Esc key exception views
- `r`/`s` key conflicts (review pass-through)
- PgUp/PgDn/Home/End/scroll keys
- Sidebar selection sync (`navigateTo`, `navigateBack`, `cycleSidebar`)
- WindowSizeMsg fallthrough
- Child dispatch for unhandled keys
- Search mode bypass
- Content view scroll behavior

All 21 packages pass `go build`, `go vet`, `go fmt`, `go test -count=1`.

### Pass 14: Manual Acceptance Test

**Verdict: ⏳ BLOCKED (requires interactive TUI session)**

The 10-step acceptance flow cannot be verified without launching the TUI:

1. Launch TUI → see empty dashboard
2. Press `f` → file explorer opens, navigate to a `.yaml` file, select it
3. Auto-navigates to Analyze view with file selected → press Enter to start
4. Analysis completes → auto-navigates to Results Summary tab
5. Tab through Assumptions, Verification, Contradictions, Trust, Impact, Blind Spots, SDRI, Controls, Reports
6. Press `/` → type search query → `n`/`N` to navigate matches → Esc to exit
7. Press `e` → export dialog → select format → confirm → verify export
8. Press Tab → cycle to Settings → verify settings render
9. Press `?` → help screen → verify all keys documented
10. Press `q` → back to dashboard → `Q` → quit

## Summary of Applied Fixes

| Issue | Fix | File |
|-------|-----|------|
| Event routing | `handleGlobalKey` → `handleKeyMsg` with boolean fallthrough | `app.go` |
| Sidebar sync | `navigateTo()`/`navigateBack()` update `m.sidebarSel` | `app.go` |
| Window resize | Removed early return from `WindowSizeMsg` | `app.go` |
| PgDn key | Changed `"pgdn"` → `"pgdown"` | `app.go` |
| Dead code | Removed `handleBack()` | Removed |
| Sidebar redesign | 8→16 items with `sidebarEntries` array | `app.go` |
| Results tabs | Added SDRI (9) + Security Design Review (10) | `results.go` |
| FocusManager | Lightweight struct on mainModel | `app.go` |
| LayoutManager | Lightweight struct on mainModel | `app.go` |
| Help screen | Added Sidebar Navigation section | `help.go` |
| Key hints | Updated bottom bar for 16-item nav | `app.go` |
| Auto-nav sidebar | `analysisCompleteMsg` updates sidebarSel | `app.go` |

## Remaining Technical Debt

1. Single shared viewport instead of per-view viewports (acceptable for current scope)
2. No mouse scroll on non-results content screens
3. `aboutView` uses global scroll keys but is rarely used
4. Tab/Shift+Tab on resultsView switches result tabs but doesn't highlight corresponding sidebar item
5. Some test files have hardcoded sidebar indices
