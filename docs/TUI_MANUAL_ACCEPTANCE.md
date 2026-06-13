# TUI Manual Acceptance Checklist

## Instructions

Run each step interactively in the TUI. Check PASS or FAIL for each.

**Prerequisites:** Ensure a valid `.yaml` architecture document exists at a known path (e.g., `~/test.yaml`).

---

### Step 1: Launch

| # | Action | Expected | Result |
|---|--------|----------|--------|
| 1.1 | Run `asf` from terminal | TUI loads, shows initial view (startup or dashboard) | ☐ PASS ☐ FAIL |
| 1.2 | Observe layout | Sidebar visible (left), content area (right), top bar, bottom bar | ☐ PASS ☐ FAIL |
| 1.3 | Observe sidebar | 16 items: Dashboard, File Explorer, Analyze, Assumptions, Verification, Contradictions, Trust Chains, Single Points of Trust, Assumption Impact Analysis, Blind Spots, SDRI, Recommended Controls, Security Design Review, Reports/Exports, Settings, Help | ☐ PASS ☐ FAIL |
| 1.4 | Press `?` | Help screen displays with all keyboard references including Sidebar Navigation section | ☐ PASS ☐ FAIL |
| 1.5 | Press `q` | Returns to previous view (dashboard or startup) | ☐ PASS ☐ FAIL |

### Step 2: File Explorer & Analysis

| # | Action | Expected | Result |
|---|--------|----------|--------|
| 2.1 | Press `f` | File explorer opens | ☐ PASS ☐ FAIL |
| 2.2 | Navigate to a `.yaml` file using ↑↓, Enter on folders, Backspace for parent | File selected | ☐ PASS ☐ FAIL |
| 2.3 | Press `.` | Hidden files toggle | ☐ PASS ☐ FAIL |
| 2.4 | Press Enter on a file | File selected, auto-navigates to Analyze view | ☐ PASS ☐ FAIL |
| 2.5 | In Analyze view, verify file path is shown | Path matches selected file | ☐ PASS ☐ FAIL |
| 2.6 | Press Enter to start analysis | Analysis runs, progress indicator visible | ☐ PASS ☐ FAIL |
| 2.7 | Wait for analysis to complete | Auto-navigates to Results Summary tab | ☐ PASS ☐ FAIL |

### Step 3: Results Browsing

| # | Action | Expected | Result |
|---|--------|----------|--------|
| 3.1 | Observe Results Summary tab | Shows aggregate counts (assumptions, verification, contradictions, etc.) | ☐ PASS ☐ FAIL |
| 3.2 | Press Tab (on results) | Switches to next result tab | ☐ PASS ☐ FAIL |
| 3.3 | Press Shift+Tab | Switches to previous result tab | ☐ PASS ☐ FAIL |
| 3.4 | Navigate to Assumptions tab | Lists all extracted assumptions with risk colors | ☐ PASS ☐ FAIL |
| 3.5 | Navigate to Verification tab | Shows verification status counts | ☐ PASS ☐ FAIL |
| 3.6 | Navigate to Contradictions tab | Shows detected contradictions | ☐ PASS ☐ FAIL |
| 3.7 | Navigate to Trust tab | Shows trust chains and SPOFs | ☐ PASS ☐ FAIL |
| 3.8 | Navigate to Impact tab | Shows priority queue and CISO view | ☐ PASS ☐ FAIL |
| 3.9 | Navigate to Blind Spots tab | Shows coverage gaps and blind spots | ☐ PASS ☐ FAIL |
| 3.10 | Navigate to SDRI tab | Shows executive summary, control inventory, coverage | ☐ PASS ☐ FAIL |
| 3.11 | Navigate to Controls tab | Shows recommended controls | ☐ PASS ☐ FAIL |
| 3.12 | Navigate to Security Design Review tab | Shows design findings, weaknesses, remediations | ☐ PASS ☐ FAIL |
| 3.13 | Navigate to Reports tab | Shows narrative, campaigns, export options | ☐ PASS ☐ FAIL |

### Step 4: Search

| # | Action | Expected | Result |
|---|--------|----------|--------|
| 4.1 | Navigate to Assumptions tab, press `/` | Search bar appears | ☐ PASS ☐ FAIL |
| 4.2 | Type search query (e.g., "auth") | Filters assumptions in real-time | ☐ PASS ☐ FAIL |
| 4.3 | Press `n` | Scroll to next match | ☐ PASS ☐ FAIL |
| 4.4 | Press `N` | Scroll to previous match | ☐ PASS ☐ FAIL |
| 4.5 | Press Enter or Esc | Exit search mode | ☐ PASS ☐ FAIL |

### Step 5: Scrolling

| # | Action | Expected | Result |
|---|--------|----------|--------|
| 5.1 | Press ↓ or j on results tab | Scroll down one line | ☐ PASS ☐ FAIL |
| 5.2 | Press ↑ or k | Scroll up one line | ☐ PASS ☐ FAIL |
| 5.3 | Press PgDn or Space | Page down | ☐ PASS ☐ FAIL |
| 5.4 | Press PgUp or b | Page up | ☐ PASS ☐ FAIL |
| 5.5 | Press Home or g | Scroll to top | ☐ PASS ☐ FAIL |
| 5.6 | Press End or G | Scroll to bottom | ☐ PASS ☐ FAIL |
| 5.7 | Observe bottom bar scroll indicator | Shows "X-Y/Z (N%)" | ☐ PASS ☐ FAIL |

### Step 6: Sidebar Navigation

| # | Action | Expected | Result |
|---|--------|----------|--------|
| 6.1 | Press Tab (when not on results/filebrowser) | Sidebar selection advances to next item | ☐ PASS ☐ FAIL |
| 6.2 | Press Shift+Tab | Sidebar selection moves backward | ☐ PASS ☐ FAIL |
| 6.3 | Cycle through each of the 16 items | Each navigates to correct view or results tab | ☐ PASS ☐ FAIL |
| 6.4 | Navigate to Dashboard | Quick actions render | ☐ PASS ☐ FAIL |
| 6.5 | Navigate to Settings | Settings render, Enter to edit, s to save | ☐ PASS ☐ FAIL |
| 6.6 | Navigate to Help | Full keyboard reference displays | ☐ PASS ☐ FAIL |

### Step 7: Export

| # | Action | Expected | Result |
|---|--------|----------|--------|
| 7.1 | From any results tab, press `e` | Export dialog opens | ☐ PASS ☐ FAIL |
| 7.2 | ↑↓ to select format | Selection highlights | ☐ PASS ☐ FAIL |
| 7.3 | Enter to confirm | Export path confirmation | ☐ PASS ☐ FAIL |
| 7.4 | Esc to cancel | Returns to results | ☐ PASS ☐ FAIL |

### Step 8: Review & Validation

| # | Action | Expected | Result |
|---|--------|----------|--------|
| 8.1 | From results, press `r` | Review mode opens | ☐ PASS ☐ FAIL |
| 8.2 | Press `s` to accept, `r` to reject, `m` modified, `n` note | Each action updates assumption status | ☐ PASS ☐ FAIL |
| 8.3 | Press `v` from review | Validation mode opens | ☐ PASS ☐ FAIL |
| 8.4 | Press `q` to navigate back | Returns to previous view | ☐ PASS ☐ FAIL |

### Step 9: Keyboard Conflicts

| # | Action | Expected | Result |
|---|--------|----------|--------|
| 9.1 | On results view, press `r` | Opens review mode (not analyze) | ☐ PASS ☐ FAIL |
| 9.2 | On review view, press `r` | Rejects current assumption | ☐ PASS ☐ FAIL |
| 9.3 | On results view, press `/` | Opens search (not scroll) | ☐ PASS ☐ FAIL |
| 9.4 | On analyze view, press `r` | Opens analyze (no-op if already there) | ☐ PASS ☐ FAIL |

### Step 10: Quit

| # | Action | Expected | Result |
|---|--------|----------|--------|
| 10.1 | Press Q | Force quit (no prompt) | ☐ PASS ☐ FAIL |
| 10.2 | Or press Ctrl+C | Force quit | ☐ PASS ☐ FAIL |

---

## Session Log

| Date | Tester | Steps Completed | Failures | Notes |
|------|--------|-----------------|----------|-------|
|      |        |                 |          |       |
