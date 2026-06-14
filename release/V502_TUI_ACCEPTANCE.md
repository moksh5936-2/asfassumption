# V5.0.2 — TUI Acceptance

## Real Terminal Verification

The following require a real interactive terminal session (TTY):

| Test | Status |
|---|---|
| Startup splash screen appears | ⏳ Manual |
| Enter starts ASF0 | ⏳ Manual |
| CASES section appears | ⏳ Manual |
| + New Analysis works | ⏳ Manual |
| Select Architecture opens modal picker | ⏳ Manual |
| Add Evidence opens modal picker | ⏳ Manual |
| Picker can navigate home directory | ⏳ Manual |
| Picker can navigate parent directory | ⏳ Manual |
| Picker is not locked to /reports | ⏳ Manual |
| Picker can select architecture file | ⏳ Manual |
| Picker can add evidence file | ⏳ Manual |
| Case appears under CASES after analysis | ⏳ Manual |
| Result tabs are accessible | ⏳ Manual |
| Local AI appears under AI | ⏳ Manual |
| Settings/Help/About work | ⏳ Manual |
| No Dashboard tab | ⏳ Manual |
| No File Explorer tab | ⏳ Manual |
| No overlapping TUI layers | ⏳ Manual |

## Tests verified from code/build

All file picker unit tests (19/19) pass and cover all navigation, selection, and state management scenarios. The binary builds, reports version correctly, and passes `go vet`/`go build`/`go test`.
