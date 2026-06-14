# File Picker UX Refactor

## Old Workflow

```
Sidebar → File Explorer (as a primary screen/route)
↓
User manually navigates to a file
↓
Select file → returns to Analyze
↓
Type Analysis Mode
↓
Run
```

Problems:
- File Explorer treated as a primary application screen (sidebar entry, route, history item)
- Two navigation systems competing: sidebar and file explorer
- Dashboard had "Open File Explorer" quick action
- Help screen had "File Explorer" section
- `f` key mapped to global file explorer shortcut
- Tab/Shift+Tab pass-through for file explorer
- Bottom bar showed "F=Files" globally

## New Workflow

```
Dashboard or Sidebar
↓
Analyze Screen
↓
Select "Architecture File" → Enter
↓
File Picker Modal opens (overlays Analyze)
↓
Navigate directories, select file
↓
Enter → file selected, modal closes
↓
Back on Analyze with file shown
↓
Optionally: Select "Evidence Files" → Enter
↓
File Picker Modal opens (evidence mode)
↓
Select multiple files (each Enter adds + stays in picker)
↓
Esc closes picker
↓
Back on Analyze with evidence files listed
↓
Select Analysis Mode
↓
Run
```

## Key Design Decisions

1. **File Picker is a modal overlay, not a screen/route**
   - No sidebar entry
   - No navigation history
   - No separate route
   - No global key binding (`f` key removed)

2. **Two picker modes**
   - Architecture mode: single-select, closes on selection
   - Evidence mode: multi-select, user stays in picker after each selection

3. **Flag-based communication** (not cmd/msg indirection)
   - `analyzeModel.requestPicker` flag set on Enter
   - `updateAnalyze` checks flag and opens picker directly
   - Avoids message routing complexity

## State Diagram

```
[Analyze Screen]
     │
     ├── Enter on "Architecture File"
     │   └──→ [File Picker: Architecture Mode]
     │           ├── Enter on file → returns path → [Analyze Screen] (path set)
     │           └── Esc → cancels → [Analyze Screen] (no change)
     │
     ├── Enter on "Evidence Files"
     │   └──→ [File Picker: Evidence Mode]
     │           ├── Enter on file → adds to list → stays in [File Picker]
     │           └── Esc → cancels → [Analyze Screen] (list updated)
     │
     ├── Enter on Mode → toggles ASF Only / ASF + AI
     │
     └── Enter on "▶ Start Analysis"
         └──→ runs analysis → [Results Screen]
```

## Focus Diagram

```
Focus Ownership:
    Analyze Screen
        ↓ (Enter on path/evidence)
    File Picker Modal (traps all input)
        ↓ (Esc or file selected)
    Analyze Screen (restored)
```

While modal is open:
- Sidebar: disabled (no Tab cycling)
- Global keys: disabled (no q, ?, r, etc.)
- Search: disabled
- Only file picker keyboard handling active

## File Picker Keyboard Map

| Key | Action |
|-----|--------|
| ↑ / k | Move selection up |
| ↓ / j | Move selection down |
| Enter | Open directory / Select file |
| Backspace | Go to parent directory |
| . | Toggle hidden files |
| Tab | Toggle preview panel |
| / | Search files |
| Esc | Cancel / Close picker |

## Files Changed

| File | Change |
|------|--------|
| `filepicker.go` | NEW — FilePickerState with modal rendering |
| `filebrowser.go` | DELETED — replaced by filepicker.go |
| `analyze.go` | Redesigned with Architecture File + Evidence Files fields, flag-based picker trigger |
| `app.go` | Removed fileBrowserView from sidebar, routing, update dispatch, renderContent, Tab passthrough, f key, bottom bar hints; added pickerActive/filePicker fields, picker handling in updateAnalyze, picker overlay in View |
| `router.go` | Removed fileBrowserView from NavigateBack fallback |
| `help.go` | Removed File Explorer section, added File Picker section, removed f key |
| `tui_test.go` | Sidebar count 18→17, new file picker tests |
| `regression_test.go` | Removed fileBrowserView test cases, updated sidebar indices |

## Migration Notes

- No database migration needed
- No config migration needed
- All file browsing capability preserved as modal
- User workflow changed: Analyze → Select File → Enter (was: f → browse → select → back to analyze)
