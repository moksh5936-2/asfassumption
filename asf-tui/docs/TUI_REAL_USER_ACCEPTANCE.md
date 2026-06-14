# TUI Real User Acceptance Test

## Test Environment

```
Terminal: real TTY (not PTY)
Size: 120×40 (minimum)
Shell: zsh
OS: macOS
Build: CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o asf-tui .
Binary: ./asf-tui
```

## Preparation

1. Open a real terminal at 120×40 (`resize` or Terminal.app > Preferences > Window > Columns 120, Rows 40)
2. `cd /Users/moksh/Project/cybersec/asf-tui`
3. `go build -o asf-tui .`
4. Verify binary exists: `ls -la asf-tui`
5. Launch: `./asf-tui`

---

## Test Sequence

### 1. Startup

| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 1.1 | Launch TUI | Startup screen renders immediately | ☐ |
| 1.2 | Observe top bar | Version, App Name visible, no truncation | ☐ |
| 1.3 | Observe sidebar | 18 items visible, scrollable if needed | ☐ |
| 1.4 | Observe bottom bar | Hints shown matching startupView | ☐ |
| 1.5 | Press Enter on "Dashboard" | Navigates to dashboardView | ☐ |

### 2. Dashboard

| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 2.1 | Dashboard loads | Quick actions visible, no startup overlay | ☐ |
| 2.2 | Arrow keys move highlight | Selection moves through quick actions | ☐ |
| 2.3 | Press Enter on "Analyze" | Navigates to analyzeView | ☐ |
| 2.4 | Navigate back (Esc or q) | Returns to dashboard, scroll preserved | ☐ |
| 2.5 | Press Enter on "Settings" | Navigates to settingsView | ☐ |
| 2.6 | Navigate back | Returns to dashboard | ☐ |
| 2.7 | Press Enter on "About" | Navigates to aboutView | ☐ |
| 2.8 | Navigate back | Returns to dashboard | ☐ |

### 3. Sidebar Navigation

| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 3.1 | Tab (from non-results view) | Moves to next sidebar item | ☐ |
| 3.2 | Shift+Tab (from non-results view) | Moves to previous sidebar item | ☐ |
| 3.3 | Down arrow / j | Moves down one sidebar item | ☐ |
| 3.4 | Up arrow / k | Moves up one sidebar item | ☐ |
| 3.5 | Enter (from sidebar) | Activates selected sidebar item | ☐ |
| 3.6 | Cycle to File Explorer (f) | Navigates to fileBrowserView | ☐ |
| 3.7 | Cycle to Analyze (r) | Navigates to analyzeView (or reviewView) | ☐ |
| 3.8 | Cycle to Help (?) | Navigates to helpView | ☐ |

### 4. File Explorer

| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 4.1 | Press f | File explorer opens | ☐ |
| 4.2 | Browse directories | Arrow keys move, Enter opens directory | ☐ |
| 4.3 | Select a file | File content shows in preview pane | ☐ |
| 4.4 | Tab/Shift+Tab | Cycles tree/preview panes | ☐ |
| 4.5 | Press Esc | Returns to previous view | ☐ |

### 5. Analyze

| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 5.1 | Press r (from dashboard) | Analyze view opens | ☐ |
| 5.2 | Path field is editable | Type a path, edits appear | ☐ |
| 5.3 | Tab through fields | Focus moves through path, mode, buttons | ☐ |
| 5.4 | Select mode (full/brief) | Selection changes | ☐ |
| 5.5 | Press Enter on "Run" | Analysis starts, progress bar shows | ☐ |
| 5.6 | Analysis completes | Results view opens automatically | ☐ |

### 6. Results Display

| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 6.1 | Results load | Summary tab (tab 0) shows by default | ☐ |
| 6.2 | Tab cycles result tabs | Each tab renders unique content | ☐ |
| 6.3 | Tab 4 shows Trust Chains | Only trust chains, no SPOFs mixed in | ☐ |
| 6.4 | Tab 11 shows SPOFs | Only SPOFs, no trust chains mixed in | ☐ |
| 6.5 | Scroll (arrow keys / j/k / pgup/pgdn) | Content scrolls, scroll indicator updates | ☐ |
| 6.6 | Press / | Search bar appears at top: "Search: █" | ☐ |
| 6.7 | Type text, Enter | Search closes, scroll to match | ☐ |
| 6.8 | Bottom bar shows /=Search | Hint visible on resultsView only | ☐ |

### 7. All Sidebar Entries (18 items)

| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 7.0 | Tab to Dashboard | Navigates to dashboardView | ☐ |
| 7.1 | Tab to File Explorer | Navigates to fileBrowserView | ☐ |
| 7.2 | Tab to Analyze | Navigates to analyzeView | ☐ |
| 7.3 | Tab to Summary | Navigates to resultsView, tab 0 | ☐ |
| 7.4 | Tab to Assumptions | Navigates to resultsView, tab 1 | ☐ |
| 7.5 | Tab to Verification | Navigates to resultsView, tab 2 | ☐ |
| 7.6 | Tab to Contradictions | Navigates to resultsView, tab 3 | ☐ |
| 7.7 | Tab to Trust Chains | Navigates to resultsView, tab 4 | ☐ |
| 7.8 | Tab to Single Points of Trust | Navigates to resultsView, tab 11 | ☐ |
| 7.9 | Tab to Assumption Impact Analysis | Navigates to resultsView, tab 5 | ☐ |
| 7.10 | Tab to Blind Spots | Navigates to resultsView, tab 6 | ☐ |
| 7.11 | Tab to SDRI | Navigates to resultsView, tab 9 | ☐ |
| 7.12 | Tab to Recommended Controls | Navigates to resultsView, tab 7 | ☐ |
| 7.13 | Tab to Security Design Review | Navigates to resultsView, tab 10 | ☐ |
| 7.14 | Tab to Reports/Exports | Navigates to resultsView, tab 8 | ☐ |
| 7.15 | Tab to Settings | Navigates to settingsView | ☐ |
| 7.16 | Tab to Help | Navigates to helpView | ☐ |
| 7.17 | Tab to About | Navigates to aboutView | ☐ |

### 8. Search

| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 8.1 | Press / on dashboard | No search bar (or global search, non-results) | ☐ |
| 8.2 | Navigate to results, press / | Search bar appears: "Search: █" | ☐ |
| 8.3 | Type characters | Characters appear in search bar | ☐ |
| 8.4 | Press Enter | Search closes | ☐ |
| 8.5 | Press / again, then Esc | Search closes | ☐ |
| 8.6 | Press /, type, press Enter again | Search re-opens, closes | ☐ |

### 9. Help

| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 9.1 | Press ? | Help view opens | ☐ |
| 9.2 | Scroll help content | All sections visible | ☐ |
| 9.3 | Sidebar section lists 18 items | Names match sidebar | ☐ |
| 9.4 | Press q or Esc | Returns to previous view | ☐ |

### 10. Settings

| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 10.1 | Navigate to Settings | Settings view renders | ☐ |
| 10.2 | Tab through fields | Focus moves between fields | ☐ |
| 10.3 | Press Enter to edit | Field becomes editable | ☐ |
| 10.4 | Press Esc | Returns to previous view | ☐ |

### 11. About

| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 11.1 | Navigate to About | About view renders | ☐ |
| 11.2 | Version displayed | Shows current version | ☐ |
| 11.3 | Press Esc/q | Returns to previous view | ☐ |

### 12. Quit

| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 12.1 | Press Q (capital Q) | TUI exits, returns to shell | ☐ |
| 12.2 | Launch TUI again | Starts fresh, no crash | ☐ |
| 12.3 | Ctrl+C | TUI exits cleanly | ☐ |

---

## Results Summary

| Category | Total Steps | Passed | Failed | N/A |
|----------|-------------|--------|--------|-----|
| 1. Startup | 5 | ☐ | ☐ | ☐ |
| 2. Dashboard | 8 | ☐ | ☐ | ☐ |
| 3. Sidebar | 8 | ☐ | ☐ | ☐ |
| 4. File Explorer | 5 | ☐ | ☐ | ☐ |
| 5. Analyze | 6 | ☐ | ☐ | ☐ |
| 6. Results | 8 | ☐ | ☐ | ☐ |
| 7. All Sidebar | 18 | ☐ | ☐ | ☐ |
| 8. Search | 6 | ☐ | ☐ | ☐ |
| 9. Help | 4 | ☐ | ☐ | ☐ |
| 10. Settings | 4 | ☐ | ☐ | ☐ |
| 11. About | 3 | ☐ | ☐ | ☐ |
| 12. Quit | 3 | ☐ | ☐ | ☐ |
| **Total** | **78** | **0** | **0** | **0** |

## Verdict

☐ **TUI_UX_ACCEPTED** — All 78 steps pass in real TTY at 120×40
☐ **TUI_UX_REJECTED** — One or more steps fail (see Failed column)

## Signed

```
Tester: _______________________
Date:   _______________________
Terminal: _____________________
Binary:  ______________________
```
