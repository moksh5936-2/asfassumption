# Result Refresh Fix Report (B3)

## Problem

After the first analysis completed successfully, re-running analysis with a new architecture document never updated the results screen. The user would see stale data — old architecture name, timestamp, and assumption counts — from the first run.

**Root cause** (in `asf-tui/results.go:423-432`):

```go
if m.analyze.result != nil && m.results.result == nil {
    m.results.result = m.analyze.result
    ...
}
m.analyze.result = nil
```

The nil gate `m.results.result == nil` meant the transfer from `analyze.result` → `results.result` happened only once. On subsequent analyses, `m.results.result` was already non-nil, so the gate blocked the transfer. Worse, `m.analyze.result = nil` ran unconditionally at line 430, destroying the new result before it could ever be displayed.

## What Was Changed

### 1. `results.go` — `updateResults` (line 423)

| Before | After |
|---|---|
| `if m.analyze.result != nil && m.results.result == nil` | `if m.analyze.result != nil` |
| `m.analyze.result = nil` after `Update()` call | `m.analyze.result = nil` inside the `if` block, after transfer |
| — | `m.results.exportComplete = false` |
| — | `m.results.exportPath = ""` |
| — | `m.results.expanded = map[int]bool{}` |

### 2. `analyze.go` — `startAnalysis` (line 182)

Added `m.result = nil` so the old result pointer is cleared when a new analysis begins. This ensures:
- No accidental stale pointer propagation
- A clean separation between the old completed analysis and the new one in progress

## Why the Nil Gate Was There and Why Removing It Is Safe

The nil gate was defensive: it prevented overwriting the results view while the user was looking at it mid-analysis. However, `updateResults` is only called when `m.currentView == resultsView`. The flow after starting a new analysis is:

1. User starts analysis → stays on `analyzeView` (progress bar shown)
2. Analysis completes → `updateAnalyze` sets `m.analyze.result`, then sends `navigateMsg{to: resultsView}`
3. Next frame: `navigateMsg` switches to `resultsView`
4. `updateResults` is called for the first time after the new result exists

Since `updateResults` only runs when already on `resultsView`, and navigation to `resultsView` happens immediately after the result is available, there is no window where a stale result could be shown. The `m.results.result` pointer is always replaced atomically with the new result by the time the user sees the view.

The additional resets (`exportComplete`, `exportPath`, `expanded`) are safe because:
- **export state**: A fresh analysis has not been exported yet; carrying over a previous export confirmation would be misleading
- **expanded map**: Section expand/collapse state is purely visual; resetting prevents rendering old content from a different architecture

## How Results Always Reflect the Latest Analysis

The end-to-end flow for each analysis is:

```
startAnalysis()
  │
  ├─ m.result = nil                    // Clear stale result pointer
  │
  ▼
analysisCompleteMsg arrives at updateAnalyze()
  │
  ├─ m.analyze.result = msg.result     // Store new result
  ├─ navigate to resultsView
  │
  ▼
updateResults() called (on resultsView)
  │
  ├─ m.analyze.result != nil → YES     // Gate removed — always enters
  ├─ m.results.result = m.analyze.result  // Replace results
  ├─ Reset export state, expanded map
  ├─ m.analyze.result = nil            // Handoff complete
  │
  ▼
viewResults() renders current m.results.result
  │
  ├─ Shows correct ArchitectureName
  ├─ Shows correct AnalysisDate
  ├─ Shows correct Assumption counts
  └─ Export uses latest result
```

No stale data survives between analyses. Each new analysis produces a fresh `AnalysisResult` that fully replaces the previous one in the results model.
