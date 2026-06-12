# ASF Progress Pipeline Fix Report

## Problem Summary

Three interacting bugs caused goroutine leaks and progress bar oscillation during ASF analysis:

### B2: Progress Channel Never Closed — Goroutine Leak

`engine.go:RunAnalysis` sends up to 6 `AnalysisProgress` values through the `progress` channel then returns, but **never calls `close(progress)`**. The consumer goroutine in `analyze.go:runAnalysisCmd` runs `for range progress` which blocks **indefinitely** because the channel is never closed. This leaks one goroutine per analysis invocation.

If `RunAnalysis` returns early (parse error, temp file failure, ASF engine error), the consumer goroutine leaks immediately. Even on success, after the consumer drains the final `Percent:100/Complete:true` message, it blocks forever waiting for more.

### B9: Progress Bar Oscillates 90→100→90→100

The `progressTickMsg` handler in `analyze.go` has:

```go
if m.progress < 100 {
    m.progress += 10
}
if m.progress >= 100 {
    m.progress = 90  // reset! causes oscillation
}
```

When real progress reaches 100 (from the engine's `AnalysisProgress{Percent:100}`), the next **tick handler** fires, resets progress to 90, and on the following tick it jumps back to 100. This creates infinite oscillation because:
1. The `analysisCompleteMsg` handler sets `m.analyze.progress = 100`
2. The tick handler sees `progress >= 100`, resets to 90, schedules another tick
3. Next tick sets 100 again, repeat forever

### B3: Fake Progress Fights Real Completion

The fake progress ticks (incrementing by 10 every 500ms) collide with the real completion signal from the engine. The tick system resets progress after real completion, and the consumer goroutine can block the producer when the channel buffer (size 10) fills during repeated analyses.

---

## Changes in `engine.go`

### 1. Close progress channel via defer

**Line 122:** Added `defer close(progress)` immediately after the first progress send.

```go
progress <- AnalysisProgress{Percent: 5, Stage: "Parsing Architecture..."}
defer close(progress)   // <-- new
```

This guarantees `close(progress)` runs when `RunAnalysis` returns — whether by:
- Normal completion (line ~185)
- Early error return (parse failure, temp file error, ASF engine error)
- **Panic** (Go's `defer` runs during panic unwinding)

The consumer goroutine's `for range progress` loop will exit when the channel closes, eliminating the goroutine leak.

---

## Changes in `analyze.go`

### 2. Fix progress tick oscillation (lines 76–84)

**Before (oscillating):**
```go
case progressTickMsg:
    if !m.running {
        return m, nil
    }
    if m.progress < 100 {
        m.progress += 10
        m.stage = analyzeStage(int(m.progress))
    }
    if m.progress >= 100 {
        m.progress = 90                      // BUG: resets to 90
        m.stage = "Finalizing Results..."
    }
    return m, m.progressCmd()
```

**After (capped at 99):**
```go
case progressTickMsg:
    if !m.running {
        return m, nil
    }
    if m.progress < 99 {
        m.progress += 10
        m.stage = analyzeStage(int(m.progress))
    }
    return m, m.progressCmd()
```

Key changes:
- **Cap fake progress at 99** (`m.progress < 99` instead of `< 100`), reserving 100 for the real completion signal
- **Remove the `>= 100` → reset to 90 branch entirely** — no oscillation possible
- When `m.progress` reaches 90 (9 ticks), subsequent ticks are no-ops (progress stays at 99)
- The real `analysisCompleteMsg` handler sets `progress = 100` once and finalizes

### 3. Clear stale result on new analysis (line 178)

Added `m.result = nil` in `startAnalysis` to prevent stale results from a previous analysis being visible:

```go
m.running = true
m.progress = 0
m.stage = "Initializing..."
m.statusMsg = ""
m.result = nil                        // <-- new: clear stale result
return m, tea.Batch(m.progressCmd(), m.runAnalysisCmd())
```

### 4. Consumer goroutine exits cleanly

In `runAnalysisCmd` (lines 189–202), the consumer goroutine:

```go
go func() {
    for range progress {
    }
}()
```

now properly exits when `RunAnalysis` returns because `defer close(progress)` guarantees the channel is closed. No code change needed here — the fix in `engine.go` enables this.

### 5. Completion handler (lines 307–313)

The `analysisCompleteMsg` handler already sets `m.analyze.running = false` **before** returning. Since `progressTickMsg` checks `!m.running` first and returns `nil` (no tick rescheduled), the tick chain is permanently broken after completion. No change needed here.

---

## How the Goroutine Leak is Fixed

| Before | After |
|---|---|
| `RunAnalysis` returns, never closes `progress` | `defer close(progress)` runs when `RunAnalysis` returns |
| Consumer `for range progress` blocks forever | Consumer exits because closed channel terminates the range loop |
| 1 leaked goroutine per analysis | Zero leaked goroutines |
| If `RunAnalysis` errors early, same leak | `defer` fires even on error return / panic |

## How the Oscillation is Fixed

The tick handler now **never resets progress backward**. Once fake progress reaches 99 (after 10 ticks at +10 each), it stays there. The only path to `progress = 100` is the `analysisCompleteMsg` handler, which fires once and also sets `running = false`, stopping all future ticks.

## How Channel Close is Guaranteed

`defer close(progress)` is placed after the first send so:
- It fires on every return path (success, error, panic)
- The progress channel is always closed exactly once
- The consumer goroutine always terminates

## Expected Behavior After Fix

1. **No goroutine leaks**: Run analysis repeatedly without resource exhaustion
2. **Progress bar climbs 0→10→20→...→90→99** via fake ticks
3. **Progress bar jumps 99→100** when real completion arrives
4. **No oscillation**: Progress never goes backward
5. **Analysis completes normally**: Final percent is 100, stage is "Complete"
6. **No deadlocks**: Channel close unblocks the consumer goroutine
7. **Clean slate**: Starting a new analysis clears any stale previous result
