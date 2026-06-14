# Case-First Navigation Refactor

## Rationale

ASF0 is organized around **CASES**, not **FEATURES**. Previously, the sidebar had a WORK section (Review Queue, Validation Queue, Reports) that treated features as top-level navigation destinations. This forced users to context-switch away from the current case to access review, validation, or reports.

The case-first model makes the **Case Workspace** the primary interface. All engine outputs are accessible as **tabs within the case workspace**. Review, Validation, and Reports become **key-driven actions** on the current case (r, v, e) rather than sidebar destinations.

## Changes

### 1. Sidebar Structure (`router.go:sidebarTreeBase`)

**Before (12 visible nodes):**
```
CASES
  + New Analysis
  📁 <case files>
WORK
  📋 Review Queue
  ✓ Validation Queue
  📦 Reports
AI
  🧠 Local AI
SYSTEM
  ⚙ Settings
  ❓ Help
  ℹ About
```

**After (8 visible nodes):**
```
CASES
  + New Analysis
  📁 <case files>
AI
  🧠 Local AI
SYSTEM
  ⚙ Settings
  ❓ Help
  ℹ About
```

### 2. Case Workspace Tabs (`results.go`)

**Before (12 tabs):** Summary, Assumptions, Verification, Contradictions, Trust Chains, Impact, Blind Spots, Controls, Reports, SDRI, Security Design Review, SPOFs

**After (7 tabs):** Overview, Assumptions, Verification, Contradictions, Trust, Controls, SDRI

| Tab | Name | Content |
|-----|------|---------|
| 0 | Overview | Case metadata + risk distribution + quick status cards |
| 1 | Assumptions | All extracted assumptions (filterable by / search) |
| 2 | Verification | Verification assessment with status cards + confidence |
| 3 | Contradictions | Detected contradictions with severity |
| 4 | Trust | Trust chains + SPOFs + Priority Queue + CISO View |
| 5 | Controls | Recommended controls (filterable) |
| 6 | SDRI | SDRI executive summary + controls + coverage + design findings + architectural weaknesses + remediations + compliance |

### 3. Key-Driven Actions (`app.go:handleGlobalKey`)

| Key | Action | Scope |
|-----|--------|-------|
| `r` | Open Review mode for current case | caseView only |
| `v` | Open Validation for current case | caseView/reviewView |
| `e` | Open Reports/Export for current case | caseView only |
| `c` | Clear current case | caseView only |
| `←`/`h` | Previous tab | caseView only |
| `→`/`l` | Next tab | caseView only |
| `/` | Search within visible tab content | caseView only |
| `Esc` | Close review/validation/reports → back to case | review/validation/reports |

Review, Validation, and Reports remain as standalone views but are only accessible via keys from the case workspace. They are NOT in the sidebar.

### 4. Startup Screen (`app.go:viewStartup`)

Unchanged — full-screen branded overlay with fox art + Enter/?/q keys.

## File Map

| File | Changes |
|------|---------|
| `router.go` | Removed WORK section (3 items) from `sidebarTreeBase` |
| `results.go` | Reduced tabs from 12→7; merged Trust/Impact/SPOFs into Trust tab; merged SDRI/Security Design Review into SDRI tab; added `renderResultTabs()` to `viewResults()` |
| `app.go` | Added ←→/hl tab navigation for caseView; updated hints bar |
| `help.go` | Removed WORK from Sidebar Tree section; updated Case Workspace keys |
| `tui_test.go` | Updated `TestSidebarTree` (12→8 nodes, review not found), `TestNewResultsModel` (12→7 tabs), `TestResultTabCount` (tab 4 counts chains+SPOFs), `TestLocalAICasesWorkNavigation` (removed WORK views) |

## Verification

```
go build ./...     # PASS
go vet ./...       # PASS
go test ./...      # PASS (all 21 packages)
```
