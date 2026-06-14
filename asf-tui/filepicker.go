package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type fileEntry struct {
	name    string
	path    string
	isDir   bool
	size    int64
	modTime time.Time
	mode    os.FileMode
}

type pickerMode int

const (
	pickerArchitecture pickerMode = iota
	pickerEvidence
)

type openFilePickerMsg struct {
	mode pickerMode
}

type filePickedMsg struct {
	path string
	mode pickerMode
}

type filePickerCancelledMsg struct{}

var archExts = map[string]bool{
	".yaml": true, ".yml": true, ".json": true,
	".md": true, ".mmd": true, ".drawio": true,
	".svg": true, ".pdf": true, ".docx": true, ".txt": true,
}

var evidenceExts = map[string]bool{
	".csv": true, ".json": true, ".yaml": true,
	".yml": true, ".txt": true, ".md": true,
	".pdf": true, ".docx": true,
}

type filePickerState struct {
	path        string
	entries     []fileEntry
	selected    int
	showHidden  bool
	searchMode  bool
	searchQuery string
	err         string
	showPreview bool
	preview     string
	mode        pickerMode
}

func (fp *filePickerState) supportedExts() map[string]bool {
	switch fp.mode {
	case pickerArchitecture:
		return archExts
	case pickerEvidence:
		return evidenceExts
	default:
		return archExts
	}
}

func newFilePickerState() filePickerState {
	return filePickerState{
		path: ".",
	}
}

func (fp *filePickerState) navigateDir(dir string) {
	fp.path = dir
	fp.selected = 0
	fp.searchMode = false
	fp.searchQuery = ""
	fp.refresh()
}

func (fp *filePickerState) refresh() {
	entries, err := fp.readDir(fp.path)
	if err != nil {
		fp.err = fmt.Sprintf("Cannot read directory: %v", err)
		fp.entries = nil
		return
	}
	fp.entries = entries
	fp.err = ""
	if fp.selected >= len(fp.entries) {
		fp.selected = 0
	}
	fp.preview = ""
}

func (fp *filePickerState) readDir(dir string) ([]fileEntry, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}
	dir = absDir

	if runtime.GOOS == "windows" {
		vol := filepath.VolumeName(dir)
		cleaned := filepath.Clean(dir)
		if cleaned == vol || cleaned == vol+string(filepath.Separator) {
			return fp.listWindowsDrives()
		}
	}

	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	names, err := f.Readdirnames(-1)
	if err != nil {
		return nil, err
	}

	var entries []fileEntry
	for _, name := range names {
		if !fp.showHidden && strings.HasPrefix(name, ".") {
			continue
		}
		fullPath := filepath.Join(dir, name)
		info, err := os.Stat(fullPath)
		if err != nil {
			if os.IsNotExist(err) || os.IsPermission(err) {
				continue
			}
			if fi, e2 := os.Lstat(fullPath); e2 == nil {
				info = fi
			} else {
				continue
			}
		}
		entries = append(entries, fileEntry{
			name:    name,
			path:    fullPath,
			isDir:   info.IsDir(),
			size:    info.Size(),
			modTime: info.ModTime(),
			mode:    info.Mode(),
		})
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].isDir != entries[j].isDir {
			return entries[i].isDir
		}
		return strings.ToLower(entries[i].name) < strings.ToLower(entries[j].name)
	})

	return entries, nil
}

func (fp *filePickerState) listWindowsDrives() ([]fileEntry, error) {
	var entries []fileEntry
	for _, d := range "ABCDEFGHIJKLMNOPQRSTUVWXYZ" {
		drive := string(d) + ":\\"
		info, err := os.Stat(drive)
		if err == nil {
			entries = append(entries, fileEntry{
				name:  string(d) + ":",
				path:  drive,
				isDir: true,
				mode:  info.Mode(),
			})
		}
	}
	return entries, nil
}

func (fp *filePickerState) handleKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	if fp.searchMode {
		switch msg.String() {
		case "esc", "enter":
			fp.searchMode = false
			fp.searchForFile()
		case "n":
			if fp.selected < len(fp.entries)-1 {
				fp.selected++
			}
			fp.updatePreview()
		case "N":
			if fp.selected > 0 {
				fp.selected--
			}
			fp.updatePreview()
		case "backspace":
			if len(fp.searchQuery) > 0 {
				fp.searchQuery = fp.searchQuery[:len(fp.searchQuery)-1]
			}
		default:
			if len(msg.String()) == 1 {
				fp.searchQuery += msg.String()
			}
		}
		return nil, true
	}

	switch msg.String() {
	case "up", "k":
		if fp.selected > 0 {
			fp.selected--
		}
		fp.updatePreview()
	case "down", "j":
		if fp.selected < len(fp.entries)-1 {
			fp.selected++
		}
		fp.updatePreview()
	case "enter":
		if fp.selected >= 0 && fp.selected < len(fp.entries) {
			entry := fp.entries[fp.selected]
			if entry.isDir {
				fp.navigateDir(entry.path)
			} else {
				ext := strings.ToLower(filepath.Ext(entry.name))
				if fp.supportedExts()[ext] {
					fp.err = ""
					return func() tea.Msg { return filePickedMsg{path: entry.path, mode: fp.mode} }, true
				} else {
					fp.err = fmt.Sprintf("Unsupported file type: %s", ext)
				}
			}
		}
	case "backspace":
		parent := filepath.Dir(fp.path)
		if parent != fp.path {
			fp.navigateDir(parent)
		}
	case "~":
		home, err := os.UserHomeDir()
		if err == nil {
			fp.navigateDir(home)
		} else {
			fp.err = fmt.Sprintf("Cannot find home directory: %v", err)
		}
	case "g":
		if runtime.GOOS != "windows" {
			fp.navigateDir("/")
		}
	case "d":
		home, err := os.UserHomeDir()
		if err == nil {
			fp.navigateDir(filepath.Join(home, "Downloads"))
		} else {
			fp.err = fmt.Sprintf("Cannot find home directory: %v", err)
		}
	case "D":
		home, err := os.UserHomeDir()
		if err == nil {
			fp.navigateDir(filepath.Join(home, "Desktop"))
		} else {
			fp.err = fmt.Sprintf("Cannot find home directory: %v", err)
		}
	case "r":
		fp.refresh()
	case "tab":
		fp.showPreview = !fp.showPreview
	case ".":
		fp.showHidden = !fp.showHidden
		fp.selected = 0
		fp.refresh()
	case "/":
		fp.searchMode = true
		fp.searchQuery = ""
	case "esc":
		return func() tea.Msg { return filePickerCancelledMsg{} }, true
	}

	return nil, true
}

func (fp *filePickerState) searchForFile() {
	if fp.searchQuery == "" {
		return
	}
	q := strings.ToLower(fp.searchQuery)
	for i, entry := range fp.entries {
		if strings.Contains(strings.ToLower(entry.name), q) {
			fp.selected = i
			fp.updatePreview()
			return
		}
	}
	fp.err = fmt.Sprintf("No match for: %s", fp.searchQuery)
}

func (fp *filePickerState) updatePreview() {
	if fp.selected < 0 || fp.selected >= len(fp.entries) {
		fp.preview = ""
		return
	}
	entry := fp.entries[fp.selected]
	if entry.isDir || !fp.showPreview {
		fp.preview = ""
		return
	}

	ext := strings.ToLower(filepath.Ext(entry.name))
	textExts := map[string]bool{".yaml": true, ".yml": true, ".json": true, ".md": true, ".txt": true, ".mmd": true, ".drawio": true, ".svg": true, ".go": true, ".py": true, ".js": true, ".ts": true, ".sh": true}
	if !textExts[ext] {
		fp.preview = fmt.Sprintf("[Binary file: %s]\nSize: %s", ext, formatFileSize(entry.size))
		return
	}

	data, err := os.ReadFile(entry.path)
	if err != nil {
		fp.preview = fmt.Sprintf("Cannot preview: %v", err)
		return
	}

	lines := strings.Split(string(data), "\n")
	maxLines := 40
	if len(lines) > maxLines {
		lines = lines[:maxLines]
		lines = append(lines, fmt.Sprintf("... (%d more lines)", len(strings.Split(string(data), "\n"))-maxLines))
	}
	fp.preview = strings.Join(lines, "\n")
}

func (m mainModel) renderFilePicker(width, height int) string {
	s := m.styles
	fp := &m.filePicker

	if fp.entries == nil {
		fp.refresh()
	}

	modeLabel := "Select Architecture File"
	if fp.mode == pickerEvidence {
		modeLabel = "Add Evidence File"
	}

	header := s.Title.Render(modeLabel)

	breadcrumb := s.DimText.Render(fmt.Sprintf("  ASF0 / New Analysis / Select Architecture / File Picker: %s", fp.path))

	var statusLine string
	if fp.err != "" {
		statusLine = s.StatusBad.Render("  " + fp.err)
	} else if fp.searchMode {
		statusLine = s.StatusWarn.Render(fmt.Sprintf("  Search: %s█  [n/N: next/prev match]", fp.searchQuery))
	}

	colName := s.SectionItem.Render("  Name")
	colSize := s.SectionItem.Render(fmt.Sprintf("%*s", 10, "Size"))
	colModified := s.SectionItem.Render(fmt.Sprintf("%*s", 17, "Modified"))

	var rows []string
	for i, entry := range fp.entries {
		style := s.SectionItem
		prefix := "  "
		if i == fp.selected {
			style = s.MenuSelected
			prefix = "▸ "
		}

		sizeStr := ""
		if !entry.isDir {
			sizeStr = fmt.Sprintf("%10s", formatFileSize(entry.size))
		} else {
			sizeStr = fmt.Sprintf("%10s", "<DIR>")
		}

		modStr := entry.modTime.Format("2006-01-02 15:04")

		ext := strings.ToLower(filepath.Ext(entry.name))
		isSupported := fp.supportedExts()[ext] || entry.isDir

		marker := ""
		if !entry.isDir && !isSupported {
			marker = s.StatusWarn.Render(" [unsupported]")
		}

		line := fmt.Sprintf("%s%-40s %s %s%s", prefix, entry.name, sizeStr, modStr, marker)

		if i == fp.selected && fp.searchMode {
			q := strings.ToLower(fp.searchQuery)
			if q != "" && strings.Contains(strings.ToLower(entry.name), q) {
				line = s.StatusGood.Render("  > ") + line[3:]
			}
		}

		rows = append(rows, style.Render(line))
	}

	fileList := lipgloss.JoinVertical(lipgloss.Left, rows...)

	mainContent := fileList
	if fp.showPreview && fp.preview != "" && width > 80 {
		leftWidth := width * 3 / 5
		if leftWidth > 60 {
			leftWidth = 60
		}
		rightWidth := width - leftWidth - 2

		previewTitle := s.Section.Render("Preview")
		previewBody := lipgloss.NewStyle().
			Width(rightWidth).
			MaxWidth(rightWidth).
			Render(fp.preview)
		previewContent := lipgloss.JoinVertical(lipgloss.Left, previewTitle, previewBody)

		rightPanel := s.BorderBox.
			Width(rightWidth + 4).
			Render(previewContent)

		mainContent = lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Width(leftWidth).Render(fileList),
			rightPanel,
		)
	}

	fp.err = ""

	hintLine := s.DimText.Render("↑↓ Select | Enter Open | Backspace Parent | . Hidden | / Search | Tab Preview | ~ Home | g Root | d Downloads | D Desktop | r Refresh | Esc Cancel")

	return lipgloss.JoinVertical(lipgloss.Left,
		header, breadcrumb,
		statusLine,
		colName+"  "+colSize+"  "+colModified,
		mainContent,
		"",
		hintLine,
	)
}

func formatFileSize(bytes int64) string {
	switch {
	case bytes > 1<<30:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(1<<30))
	case bytes > 1<<20:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(1<<20))
	case bytes > 1<<10:
		return fmt.Sprintf("%.1f KB", float64(bytes)/float64(1<<10))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}
