# File Picker Cross-Platform Acceptance

## Test Results

| # | Test | Status |
|---|------|--------|
| 1 | Picker does not default to /reports unless last used | ✅ |
| 2 | Picker starts from current working directory or home | ✅ |
| 3 | Picker can navigate parent directory | ✅ |
| 4 | Picker can navigate into child directory | ✅ |
| 5 | Picker can jump home | ✅ |
| 6 | Picker can jump to root | ✅ |
| 7 | Picker can select supported architecture file | ✅ |
| 8 | Picker rejects unsupported file extension | ✅ |
| 9 | Architecture selection updates New Analysis state | ✅ |
| 10 | Evidence selection appends evidence file | ✅ |
| 11 | Evidence selection preserves existing evidence | ✅ |
| 12 | Cancel preserves previous state | ✅ |
| 13 | Search filters directory entries | ✅ |
| 14 | Paths with spaces work | ✅ |
| 15 | filepath.Clean/Abs used correctly | ✅ |
| 16 | Windows path handling does not use Unix-only separators | ✅ |
| 17 | Last picker path saved per mode on selection | ✅ |
| 18 | Start path priority: last used > doc dir > cwd/home | ✅ |
| 19 | Start path never empty | ✅ |

**All 19 tests pass on macOS (CI: Linux/Windows TBD).**

## Root Causes Fixed

### `/reports` lock-in
- **Cause**: `pickerStartPath` unconditionally returned `./reports` for architecture mode when no path was previously used.
- **Fix**: Priority-based start path: `lastPickerPaths[mode]` > `filepath.Dir(docPath())` for architecture > `os.Getwd()` / `os.UserHomeDir()` fallback. All paths use `filepath.Clean` and `filepath.Abs`.

### Cross-platform path navigation
- **Cause**: Manual `/` concatenation for "go to parent" and navigation commands assumed Unix separators.
- **Fix**: All path manipulation uses `filepath.Dir`, `filepath.Join`, `filepath.Clean` from the stdlib, which handle `/` vs `\` automatically.

### Missing navigation keys
- **Cause**: Only `Esc`, `Enter`, arrows, and `.` were bound. Keys `~` (home), `g` (root), `d` (trash), `D` (temp) were missing.
- **Fix**: Added all missing keybindings. `D` (Shift-d) on macOS now navigates to `$TMPDIR` (which points to `/var/folders/...`), while on Linux/Windows it uses `os.TempDir()`.

## Build
- `go build ./...` — clean
- `go vet ./...` — clean
- `go test ./...` — all packages pass
