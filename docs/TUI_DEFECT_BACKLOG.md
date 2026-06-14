# TUI Defect Backlog

## Methodology

Defects identified through:
1. TUI launch observations (via `script` PTY)
2. CLI mode testing (analyze, doctor, edge cases)
3. Full test suite execution
4. `go vet` / `go fmt` / `go build` static analysis
5. Code structure analysis
6. Logging infrastructure verification

---

## Critical

*No critical defects found.*

---

## High

### H-01: Log directory never created

**Location:** `main.go` — `initLogger()`  
**Reproduction:**
1. Launch ASF TUI or CLI
2. Check `~/.asf/logs/`
3. Directory does not exist

**Expected:** Log directory `~/.asf/logs/` is created on first launch if missing.  
**Actual:** Directory is never created. `asfLog` silently discards log output when the directory is absent (or panics if `*os.File` is nil).  
**Root Cause:** No `os.MkdirAll` call in the logger initialization path.  
**Recommended Fix:** Add `os.MkdirAll(filepath.Dir(logPath), 0755)` before opening the log file.

---

### H-02: Terminal size 0×0 at startup triggers minimum size warning but TUI continues rendering

**Location:** `app.go:554-556`  
**Reproduction:**
1. Launch TUI in an environment where terminal size cannot be determined (e.g., `script` PTY, CI, headless)
2. TUI prints "Terminal too small. Minimum: 60x12 Current: 0x0"
3. TUI continues to render full Dashboard with sidebar below the error

**Expected:** TUI should either block rendering (repeatedly display the size error until terminal is resized) or gracefully handle 0×0 by assuming a default size.  
**Actual:** The size check only runs once in `View()` — after displaying the error, the TUI renders the full UI underneath. The error message scrolls out of view.  
**Root Cause:** The size check is in `View()` (pure rendering) rather than in `Update()` with a blocking state. Minimum size enforcement is not preventing further rendering.  
**Recommended Fix:** Add a blocking state (`m.sizeOk bool`) that prevents content rendering until terminal meets minimum requirements. Set it in `tea.WindowSizeMsg` handler.

---

### H-03: viewport scroll position bypassed in analysisCompleteMsg and fileSelectedMsg

**Location:** `app.go:250-252`, `app.go:261-263`  
**Reproduction:**
1. Scroll down in results or analyze view
2. Run a new analysis or select a new file
3. Scroll position resets to 0

**Expected:** Scroll position should be saved/restored via `m.scrollY` map  
**Actual:** `m.vp.YOffset = 0` is hardcoded — bypasses scrollY map  
**Root Cause:** `analysisCompleteMsg` and `fileSelectedMsg` set YOffset directly instead of using the scrollY save/restore pattern  
**Recommended Fix:** Remove `m.vp.YOffset = 0` and `m.scrollY[...] = 0` lines — the `navigateTo` call already restores scroll position from the scrollY map

---

## Medium

### M-01: focusManager struct is dead code

**Location:** `app.go:47-50`, `app.go:123`, `app.go:173`  
**Reproduction:**
1. Search for all assignments to `m.focusMgr.activeView` — none found
2. Search for all reads of `m.focusMgr.subFocus` — none found

**Expected:** Focus manager should be used to track active sub-focus (e.g., which sidebar item has focus, which field in settings is being edited)  
**Actual:** Struct is instantiated but never read or written. `activeView` duplicates `Router.currentView`  
**Root Cause:** Focus tracking was split between Router and focusManager but focusManager was never wired up  
**Recommended Fix:** Either remove the unused struct or wire it into the focus flow (sidebar selection tracking, sub-focus for settings fields, etc.)

---

### M-02: Local AI Models view (localaiView) has no sidebar entry

**Location:** `app.go:131-148`  
**Reproduction:**
1. Launch TUI
2. Check sidebar for "AI Models" or "Local AI"
3. Not found
4. Only reachable via Dashboard → select "Local AI Models" or press `l`

**Expected:** All primary views should be reachable from the sidebar  
**Actual:** `localaiView` is not in `sidebarEntries`. It is the only non-startup, non-error view without a sidebar entry.  
**Root Cause:** Sidebar was expanded to 16 items but localAI was omitted  
**Recommended Fix:** Add "AI Models" entry to sidebarEntries, or document it as intentionally hidden

---

### M-03: Review and Validation views have no sidebar entries

**Location:** `app.go:131-148`  
**Reproduction:**
1. Launch TUI, press `r` from results to enter review mode
2. Press `q` to go back
3. No way to return to review or validation from sidebar

**Expected:** Review and validation can be accessed from sidebar  
**Actual:** Not in sidebar — only reachable from results view contextually  
**Root Cause:** These are context-dependent sub-views, not primary navigation destinations  
**Recommended Fix:** Either add them to sidebar (greyed out when no results) or document them as contextual screens reachable only from results

---

### M-04: About view has no sidebar entry

**Location:** `app.go:131-148`  
**Reproduction:** Navigate Help → no "About" entry. No sidebar entry for About.  
**Expected:** About should be accessible from sidebar or Help screen  
**Actual:** Only reachable via Dashboard → "About" or `i` key  
**Recommended Fix:** Add "About" as a Help sub-item or add to sidebar

---

### M-05: Results tabs 4 (Trust) and 7 (Single Points of Trust) share the same tab index

**Location:** `app.go:138-139`  
**Reproduction:** Navigate to "Trust Chains" then "Single Points of Trust" — both set `m.results.resultTab = 4`  
**Expected:** Single Points of Trust should be on a different tab from Trust Chains  
**Actual:** Both share tab 4. The `renderTrustOutput` function renders both trust chains and SPOFs in a single view, so functionally this is consistent but the sidebar implies they're distinct.  
**Root Cause:** `ActivateSidebarTab()` returns tab 4 for both entries  
**Recommended Fix:** Either split SPOF into its own results tab (tab 11) or rename sidebar entry to "Trust / SPOFs"

---

## Low

### L-01: Sidebar width is hardcoded

**Location:** `app.go:64`  
**Detail:** `sidebarWidth: 23` in `newLayoutManager()` — not configurable  
**Recommended Fix:** Make configurable via settings or dynamic calculation based on longest item name

### L-02: View history truncated silently at 50 entries

**Location:** `app.go:78-82`  
**Detail:** `maxHistory = 50` — history drops oldest entries when exceeded. Reasonable but undocumented.  
**Status:** Fixed in v4.0.1 but the constant could be documented

### L-03: No mouse wheel support on dashboard, analyze, settings, review, validation views

**Location:** `app.go:201-209`  
**Detail:** Mouse wheel events are handled globally but `handleGlobalKey` only scrolls the viewport for results/help/about views. On other views, mouse wheel is passed through but those child models don't handle MouseMsg.  
**Recommended Fix:** Add MouseMsg handling to viewport for content-heavy screens

### L-04: TUI binary (11MB unstripped, 11MB stripped) could be smaller

**Detail:** Build uses `-ldflags="-s -w"` which strips DWARF/symbol table. The binary is 11MB. Likely acceptable for a Go binary with embedded engine.  
**Status:** Informational

### L-05: "Security Design Review" sidebar entry shares tab 10 with SDRI

**Location:** `app.go:142-143`  
**Detail:** "SDRI" = tab 9, "Security Design Review" = tab 10. These are different tabs but the naming could confuse users about what content each shows.  
**Recommended Fix:** Consider renaming "SDRI" to "SDRI Overview" and "Security Design Review" to "SDR Details"

---

## Defect Count Summary

| Severity | Count |
|----------|-------|
| Critical | 0 |
| High | 3 |
| Medium | 5 |
| Low | 5 |
| **Total** | **13** |
