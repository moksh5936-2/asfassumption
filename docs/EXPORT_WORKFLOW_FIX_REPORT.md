# Export Workflow Fix Report (B4)

## Problem Summary

The ASF TUI export workflow had a fake-success bug in two places:

### 1. Orphan export model (export.go:770-774)

The `exportModel.Update` method handled the `y` confirmation key by setting `m.done = true` but **never called `ExportResult()`**. The `exportPath` field remained an empty string (Go zero value). The UI rendered "✓ Export Complete" with a blank path, suggesting a file was written when none was.

### 2. Bypassed export model (results.go:49-57)

The results view's `e` (export) handler called `ExportResult()` directly with a **hardcoded directory** (`./reports`) and used `m.exportFormat` — set once from config. This bypassed the export model's UI-based format selection entirely. The user was never shown the format picker, and the output directory always defaulted to `./reports` regardless of the config's `output.directory` setting.

## Root Cause

Two disconnected code paths existed for export:
- A dead-end UI (`exportView`) that showed a confirmation screen but never performed I/O.
- A silent direct export (`results.go` `e` key) that ignored the UI entirely.

Neither path gave the user a meaningful choice of format or respected the configured output directory.

## What Was Changed

### `asf-tui/export.go`

1. **Added fields to `exportModel` struct** (`export.go:744-746`):
   - `result *AnalysisResult` — the analysis data to export
   - `outputDir string` — the directory to write into (from config or `./reports` fallback)
   - `err error` — captures export failures for display

2. **`y` key now performs the export** (`export.go:775-786`):
   ```go
   case "y":
       if m.showConfirmation && !m.done {
           if m.result != nil {
               path, err := ExportResult(m.result, m.format, m.outputDir)
               if err != nil {
                   m.err = err
               } else {
                   m.done = true
                   m.exportPath = path
               }
           }
       }
   ```

3. **Error view added** (`export.go:796-807`):
   When `ex.err != nil`, the UI displays "✗ Export Failed" with the error message in `StatusBad` style. Previously, failures were silently swallowed.

4. **`Esc` handler resets error state** (`export.go:774`).

### `asf-tui/results.go`

5. **`e` key navigates to export view** (`results.go:49-52`):
   ```go
   case "e":
       if m.result != nil {
           return m, func() tea.Msg { return navigateMsg{to: exportView} }
       }
   ```
   Instead of performing a hardcoded export, it now sends a `navigateMsg` that causes `mainModel` to set up the export model and switch to the format-picker UI.

### `asf-tui/app.go`

6. **Export model initialization on navigation** (`app.go:153-166`):
   When a `navigateMsg{to: exportView}` is handled, the export model is populated with:
   - `result` ← the current analysis result from `m.results.result`
   - `outputDir` ← `m.config.Output.Directory` (falls back to `"./reports"` if empty)
   - `format` ← `exportFormatFromConfig(m.config)` as the default (user can change it)
   - All state fields (`selected`, `done`, `exportPath`, `showConfirmation`, `err`) are reset to initial values.

7. **Help text updated** (`app.go`):
   Added `"y: Confirm export"` to the export view help bar.

## How Export Now Works

1. **User is on Results view** — presses `e`.
2. **App navigates to Export view** — shows format picker (JSON, Markdown, HTML, CSV, PDF). Default matches config `output.default`.
3. **User selects format** with `↑/↓` and presses `Enter`.
4. **Confirmation screen** — shows "Export as <format>?" with `Y` to confirm / `Esc` to cancel.
5. **User presses `Y`** — `ExportResult()` is called with the **current result**, the **chosen format**, and the **configured output directory**.
6. **On success** — "✓ Export Complete" is shown with the **exact absolute/relative path** of the written file.
7. **On failure** — "✗ Export Failed" is shown with the error message.

## Error Handling

| Scenario | Behavior |
|---|---|
| `outputDir` does not exist / cannot be created | `MkdirAll` fails → error shown |
| File write permission denied | `os.WriteFile` fails → error shown |
| `result` is nil | `y` key is a no-op (guard clause), view unchanged |
| Format is unrecognized | `ExportResult` switch hits default (empty path) → returns empty string with nil err; no file is written. This is a pre-existing edge case in `ExportResult` (no default case returns error). |

## Files Modified

| File | Lines Changed | Purpose |
|---|---|---|
| `asf-tui/export.go` | 738-786, 796-807 | Added result/outputDir/err fields; actual export call on `y`; error view |
| `asf-tui/results.go` | 49-52 | Changed `e` key from direct export to navigation |
| `asf-tui/app.go` | 153-166, 269 | Export model setup on navigation; help text update |
