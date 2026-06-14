# File Picker Route Map

## Old Route Map (Before)

```
Sidebar (18 items)
├── Dashboard (dashboardView)
├── File Explorer (fileBrowserView)     ← REMOVED
├── Analyze (analyzeView)
├── Summary (resultsView, tab 0)
├── ...
├── About (aboutView)

Global keys:
    f → fileBrowserView                 ← REMOVED

Bottom bar:
    F=Files  R=Analyze  ?=Help  Q=Quit  ← REMOVED "F=Files"

Tab handling:
    resultsView, fileBrowserView → child ← REMOVED fileBrowserView

NavigateBack fallback:
    fileBrowserView → dashboardView     ← REMOVED

Help:
    "f  Open file explorer"             ← REMOVED
    File Explorer section               ← REMOVED
```

## New Route Map (After)

```
Sidebar (17 items)
├── Dashboard (dashboardView)
├── Analyze (analyzeView)               ← File picker is INSIDE Analyze, not a route
├── Summary (resultsView, tab 0)
├── ...
├── About (aboutView)

Global keys:
    (no f key)
    r → analyzeView (still works)

Bottom bar:
    R=Analyze  ?=Help  Q=Quit           ← No "F=Files"

Tab handling:
    resultsView → child only
    (no fileBrowserView passthrough)

NavigateBack fallback:
    (no fileBrowserView)

Help:
    "r  Run analysis (open Analyze view)"
    File Picker section (documented under Analyze workflow)
```

## Routing Table

| View | Sidebar Index | Global Key | Tab Passthrough | NavigateBack Fallback |
|------|:---:|:---:|:---:|:---:|
| startupView | — | — | — | — |
| dashboardView | 0 | — | — | → dashboardView |
| analyzeView | 1 | r | — | → dashboardView |
| resultsView | 2 | — | Tab/Shift+Tab | → dashboardView |
| localaiView | — | — | — | → dashboardView |
| settingsView | 14 | — | — | → dashboardView |
| aboutView | 16 | — | — | → dashboardView |
| exportView | — | e | — | → dashboardView |
| reviewView | — | r (results) | — | → dashboardView |
| validationView | — | v | — | → dashboardView |
| helpView | 15 | ? | — | → dashboardView |

## File Picker (Not a Route)

The file picker is NOT in the route table. It is a modal overlay within the Analyze screen:

```
State: mainModel.pickerActive = true
Enter: mainModel.updateAnalyze() → sets pickerActive → View() renders overlay
Exit:  filePickedMsg or filePickerCancelledMsg → pickerActive = false → Analyze screen restored
```

Route/Router interaction:
- Router.currentView remains `analyzeView` during picker operation
- No history entry created
- No sidebarSel change
- No view change

## Message Flow

```
User presses Enter on "Architecture File"
    ↓
analyzeModel.handleEnter() → sets requestPicker = true
    ↓
updateAnalyze() → checks flag → creates filePicker → sets pickerActive = true
    ↓
View() → detects pickerActive → renders picker overlay
    ↓
User navigates and presses Enter
    ↓
filePicker.handleKey() → returns filePickedMsg
    ↓
Update() → receives filePickedMsg → handleFilePicked() → updates analyzeModel
    ↓
pickerActive = false → View() renders Analyze screen with file shown
```
