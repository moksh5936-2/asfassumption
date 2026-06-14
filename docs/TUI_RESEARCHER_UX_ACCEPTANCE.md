# TUI Researcher UX — Acceptance Test Plan

## Scope
Validate 10 UX stabilization fixes for the ASF0 TUI result tabs.
All tests run from the **case workspace** with an active analysis result (press Enter then run an analysis, or open a `.yaml` file).

---

## FIX 1 — Content Clipping Fix
| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 1.1  | Run analysis on any `.yaml` architecture file | Result loads, case workspace opens | |
| 1.2  | Switch to Assumptions tab | All assumption rows are visible (scroll to see all) | |
| 1.3  | Switch to Trust tab | All trust chains, SPOFs, queue items are visible | |
| 1.4  | Resize terminal to minimum (60x10) | Content still usable, no invisible content | |
| 1.5  | Scroll with ↑↓/PgUp/PgDn/Home/End | Content scrolls completely, last line of each section reachable | |

**Pass condition:** Every item in every tab is reachable via scrolling. No content is permanently clipped below viewport bottom.

---

## FIX 2 — Selectable List Navigation (All Tabs)
| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 2.1  | Open Assumptions tab | First item has `▸` cursor | |
| 2.2  | Press `↓` / `j` | Cursor moves to next assumption | |
| 2.3  | Press `↑` / `k` | Cursor moves to previous assumption | |
| 2.4  | Loop around: press `↑` at top | Cursor stays at first item (no wrap) | |
| 2.5  | Loop around: press `↓` at bottom | Cursor stays at last item (no wrap) | |
| 2.6  | Press `→` to Verification tab | Cursor appears on first verification item | |
| 2.7  | Navigate Verification items with `↑`/`↓` | Selection moves correctly | |
| 2.8  | Repeat for Contradictions, Trust, Controls, SDRI tabs | Cursor visible and navigable in each | |

**Pass condition:** All 6 list-based tabs support `↑↓/jk` selection. Cursor never wraps. Selection stays within bounds.

---

## FIX 3 — Detail Expansion (Enter)
| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 3.1  | Select an assumption, press `Enter` | Expanded detail shows ID, Risk, Confidence, Component, Category, STRIDE, Evidence | |
| 3.2  | Press `Enter` again | Detail collapses, item list shows | |
| 3.3  | Select a contradiction, press `Enter` | Shows reason and affected assumptions | |
| 3.4  | Select a trust chain, press `Enter` | Shows chain nodes and diagram | |
| 3.5  | Select a control, press `Enter` | Shows rationale, category, mitigated assumptions | |
| 3.6  | Press `Esc` while detail is open | Detail closes | |

**Pass condition:** `Enter` toggles per-item detail expansion. `Esc` closes detail. Detail shows relevant fields for the selected item type.

---

## FIX 4 — Breadcrumbs
| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 4.1  | Open Assumptions tab | Breadcrumb shows: `ASF0 / Case: <name> / Assumptions / #1` | |
| 4.2  | Navigate to second item with `↓` | Breadcrumb updates to `... / #2` | |
| 4.3  | Press `Enter` to open detail | Breadcrumb appends ` / detail` | |
| 4.4  | Press `→` to Verification tab | Breadcrumb shows `ASF0 / Case: <name> / Verification / #1` | |
| 4.5  | Press `←` back to Assumptions | Selection and breadcrumb restored to previous state | |
| 4.6  | Switch to Overview tab | Breadcrumb shows `... / Overview` (no item number) | |

**Pass condition:** Breadcrumb reflects current position: case name → tab → item index → detail state.

---

## FIX 5 — Contextual Help
| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 5.1  | Open Assumptions tab | Hints bar shows `↑↓=Select Enter=Detail /=Search` | |
| 5.2  | Press `/` | Hints bar shows `filter:` prefix | |
| 5.3  | Switch to Overview tab | Hints bar shows `↑↓=Scroll` instead | |
| 5.4  | Switch to any list tab, press `Enter` to open detail | Hints continue to show `↑↓=Select Enter=Detail` | |

**Pass condition:** Hints bar changes per-tab to show relevant navigation keys.

---

## FIX 6 — Per-Tab Footer Hints
| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 6.1  | Open Assumptions tab | Footer shows: `↑↓=Select │ Enter=Detail │ ←→=Tabs │ r=Review │ v=Validate │ e=Reports │ c=Clear │ /=Search │ Tab=Sidebar │ ?=Help │ q=Back │ Q=Quit` | |
| 6.2  | Press `/` | `filter:` appears alongside the hints | |
| 6.3  | Open Overview tab | `↑↓=Scroll` replaces `↑↓=Select` and `Enter=Detail` and `/=Search` | |

**Pass condition:** Footer hints correctly show the current tab's available actions.

---

## FIX 7 — Search Inside Tabs (/ + n/N + Esc)
| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 7.1  | Open Assumptions tab, press `/` | Filter prompt appears, characters typed build search query | |
| 7.2  | Type `authentication` | List filters to items matching "authentication" | |
| 7.3  | Press `↓` / `↓` / `n` | Moves between filtered items (next/prev match) | |
| 7.4  | Press `Esc` | Filter clears, full list restored, selection reset to top | |
| 7.5  | Open Contradictions tab, press `/` type `no encryption` | Contradictions filter by search term | |
| 7.6  | Open Trust tab, press `/` type `admin` | Trust chains/SPOFs/queue filter by search term | |
| 7.7  | Open Controls tab, press `/` type `auth` | Controls filter by search term | |
| 7.8  | Press `backspace` while filtering | Characters delete from search query | |

**Pass condition:** All 6 list tabs support `/` search. Filter is per-tab (switching tabs preserves filter state per tab). `Esc` clears filter. `n`/`N` navigate matches. `Backspace` edits query.

---

## FIX 8 — Per-Tab Scroll Memory
| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 8.1  | Open Assumptions tab, scroll 30 lines down | Content shows items starting from line 30 | |
| 8.2  | Press `→` to Contradictions tab | Contradictions tab shows from its own scroll position | |
| 8.3  | Press `←` back to Assumptions tab | Assumptions tab scroll restored to line 30 | |
| 8.4  | Scroll Assumptions to line 50, switch to Verification, switch back | Returns to line 50 (most recent scroll for that tab) | |

**Pass condition:** Each tab remembers its scroll position independently. Switching away and back restores the previous scroll position.

---

## FIX 9 — Viewport Layout Math
| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 9.1  | Maximize terminal to full screen | Header, sidebar, content, hints bar, status bar all fit exactly | |
| 9.2  | Verify no double scrollbars | Viewport scrolls content only, not the chrome | |
| 9.3  | Set terminal to 80x24 | Layout reflows correctly | |
| 9.4  | Set terminal to 60x10 | Minimum size warning or layout adapts | |
| 9.5  | Open Overview tab (summary cards) | All cards render fully, no truncation | |

**Pass condition:** `mainHeight()` correctly subtracts header (1), breadcrumb (1), hints bar (1), status bar (1). Content fills available space without overflow clipping.

---

## FIX 10 — Acceptance Test Doc
| Step | Action | Expected Result | Pass/Fail |
|------|--------|-----------------|-----------|
| 10.1 | Verify this document exists at `docs/TUI_RESEARCHER_UX_ACCEPTANCE.md` | File exists with all 10 fix sections | |
| 10.2 | Run through each test case in sections 1–9 | All cases pass | |

**Pass condition:** This document is reviewed and all 9 preceding fix sections are marked Pass.

---

## Notes
- Tests assume a valid `.yaml` architecture file is loaded (e.g., `asf-tui/arch/sample-advanced.yaml`).
- If no analysis is loaded, the case workspace shows "No results available. Run an analysis first."
- Keyboard navigation requires content focus (press `Tab` to leave sidebar if cursor appears in sidebar).
