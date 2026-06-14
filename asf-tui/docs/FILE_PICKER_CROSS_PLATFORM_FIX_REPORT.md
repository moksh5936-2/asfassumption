# File Picker Cross-Platform Fix Report

## Problem
The file picker modal in `asf-tui` exhibited two issues:
1. **`/reports` lock-in**: On first open, the picker defaulted to `./reports` regardless of the user's working directory, trapping the user in a fixed location.
2. **macOS (and cross-platform) path friction**: Path navigation used manual string manipulation and platform-agnostic assumptions (e.g., `/` separators, no `$TMPDIR` awareness).

## Root Causes

### 1. `pickerStartPath` unconditionally returning `/reports`
```go
// Old code (app.go)
func (m *mainModel) pickerStartPath() string {
    return "./reports"  // hardcoded, overridden only by lastPickerPaths
}
```
This meant every fresh analysis session started in `./reports`, even if the user was working in `/Users/me/projects/audit`.

### 2. Manual path manipulation
```go
// Old code (filepicker.go)
func (fp *filePickerState) handleKey(msg tea.KeyMsg) (tea.Cmd, bool) {
    // ... backspace did: path[:strings.LastIndex(path, "/")]
    // This breaks on macOS where paths have no trailing slash on root,
    // and completely fails on Windows with backslashes.
}
```

### 3. Missing macOS navigation keys
Key shortcuts for common macOS directories (`~` for home, `D` for `$TMPDIR`) were absent, requiring users to manually type paths.

## Fixes Applied

### File: `app.go` — `pickerStartPath` rewritten
Priority order:
1. `lastPickerPaths[mode]` — resume where you left off
2. `filepath.Dir(m.analyze.docPath())` — start near the architecture file
3. `os.Getwd()` — current working directory
4. `os.UserHomeDir()` — absolute fallback

All paths normalized via `filepath.Clean`.

### File: `filepicker.go` — path navigation
- `backspace` uses `filepath.Dir(path)` instead of manual string slicing
- `~` key jumps to `os.UserHomeDir()` via `filepath.Clean`
- `g` key jumps to `/` (Unix) or `filepath.VolumeName(path)+`\\` ` (Windows)
- `D` key jumps to `os.TempDir()` (respects `$TMPDIR` on macOS)
- `d` key redirects to trash (user-configured, safe fallback)
- `r` key refreshes current directory
- `Enter` on directory uses `filepath.Join` and `filepath.Clean`

### File: `app.go` — `openFilePickerMsg` handler
- `lastPickerPaths` map tracks last selected directory per mode (architecture vs evidence)
- `pickerStartPath` used for both the update function and the initial modal setup

### File: `app.go` — added `os` import
- Added `"os"` to imports for `os.Getwd`, `os.UserHomeDir`, `os.TempDir`

## File Changes Summary

| File | Changes |
|------|---------|
| `filepicker.go` | Rewritten path navigation (backspace, ~, g, d, D, r keys); all path ops use `filepath.*` |
| `app.go` | Rewritten `pickerStartPath` with priority-based fallback; added `lastPickerPaths` map; added `os` import |
| `picker_test.go` | 19 new tests covering all path scenarios, cross-platform safety, state preservation |
| `docs/FILE_PICKER_ACCEPTANCE.md` | Acceptance checklist |
| `docs/FILE_PICKER_CROSS_PLATFORM_FIX_REPORT.md` | This report |

## Regression Coverage
All existing tests in `tui_test.go` pass unchanged. The 19 new picker-specific tests cover the cross-platform concerns.
