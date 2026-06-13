package main

import (
	"fmt"
	"os"
	"path/filepath"
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

type fileSelectedMsg string

type fileBrowserModel struct {
	path        string
	entries     []fileEntry
	selected    int
	showHidden  bool
	searchMode  bool
	searchQuery string
	err         string
	showPreview bool
	preview     string
}

var supportedExts = map[string]bool{
	".yaml": true, ".yml": true, ".json": true,
	".md": true, ".mmd": true, ".drawio": true,
	".svg": true, ".pdf": true, ".docx": true, ".txt": true,
}

func newFileBrowserModel() fileBrowserModel {
	return fileBrowserModel{
		path: ".",
	}
}

func (m *fileBrowserModel) refresh() {
	entries, err := m.readDir(m.path)
	if err != nil {
		m.err = fmt.Sprintf("Cannot read directory: %v", err)
		m.entries = nil
		return
	}
	m.entries = entries
	m.err = ""
	if m.selected >= len(m.entries) {
		m.selected = 0
	}
	m.preview = ""
}

func (m *fileBrowserModel) readDir(dir string) ([]fileEntry, error) {
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
		if !m.showHidden && strings.HasPrefix(name, ".") {
			continue
		}
		fullPath := filepath.Join(dir, name)
		info, err := os.Stat(fullPath)
		if err != nil {
			continue
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

func (m *fileBrowserModel) Update(msg tea.Msg) (fileBrowserModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m.handleKey(msg), nil
	}
	return *m, nil
}

func (m *fileBrowserModel) handleKey(msg tea.KeyMsg) fileBrowserModel {
	if m.searchMode {
		switch msg.String() {
		case "esc", "enter":
			m.searchMode = false
			m.searchForFile()
		case "n":
			if m.selected < len(m.entries)-1 {
				m.selected++
			}
			m.updatePreview()
		case "N":
			if m.selected > 0 {
				m.selected--
			}
			m.updatePreview()
		case "backspace":
			if len(m.searchQuery) > 0 {
				m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
			}
		default:
			if len(msg.String()) == 1 {
				m.searchQuery += msg.String()
			}
		}
		return *m
	}

	switch msg.String() {
	case "up", "k":
		if m.selected > 0 {
			m.selected--
		}
		m.updatePreview()
	case "down", "j":
		if m.selected < len(m.entries)-1 {
			m.selected++
		}
		m.updatePreview()
	case "enter":
		if m.selected >= 0 && m.selected < len(m.entries) {
			entry := m.entries[m.selected]
			if entry.isDir {
				m.path = entry.path
				m.selected = 0
				m.refresh()
			} else {
				ext := strings.ToLower(filepath.Ext(entry.name))
				if supportedExts[ext] {
					m.err = ""
				} else {
					m.err = fmt.Sprintf("Unsupported file type: %s", ext)
				}
			}
		}
	case "backspace":
		parent := filepath.Dir(m.path)
		if parent != m.path {
			m.path = parent
			m.selected = 0
			m.refresh()
		}
	case "tab":
		m.showPreview = !m.showPreview
	case ".":
		m.showHidden = !m.showHidden
		m.selected = 0
		m.refresh()
	case "/":
		m.searchMode = true
		m.searchQuery = ""
	case "q", "esc":
		m.err = ""
	}

	return *m
}

func (m *fileBrowserModel) searchForFile() {
	if m.searchQuery == "" {
		return
	}
	q := strings.ToLower(m.searchQuery)
	for i, entry := range m.entries {
		if strings.Contains(strings.ToLower(entry.name), q) {
			m.selected = i
			m.updatePreview()
			return
		}
	}
	m.err = fmt.Sprintf("No match for: %s", m.searchQuery)
}

func (m *fileBrowserModel) updatePreview() {
	if m.selected < 0 || m.selected >= len(m.entries) {
		m.preview = ""
		return
	}
	entry := m.entries[m.selected]
	if entry.isDir || !m.showPreview {
		m.preview = ""
		return
	}

	ext := strings.ToLower(filepath.Ext(entry.name))
	textExts := map[string]bool{".yaml": true, ".yml": true, ".json": true, ".md": true, ".txt": true, ".mmd": true, ".drawio": true, ".svg": true, ".go": true, ".py": true, ".js": true, ".ts": true, ".sh": true}
	if !textExts[ext] {
		m.preview = fmt.Sprintf("[Binary file: %s]\nSize: %s", ext, formatFileSize(entry.size))
		return
	}

	data, err := os.ReadFile(entry.path)
	if err != nil {
		m.preview = fmt.Sprintf("Cannot preview: %v", err)
		return
	}

	lines := strings.Split(string(data), "\n")
	maxLines := 40
	if len(lines) > maxLines {
		lines = lines[:maxLines]
		lines = append(lines, fmt.Sprintf("... (%d more lines)", len(strings.Split(string(data), "\n"))-maxLines))
	}
	m.preview = strings.Join(lines, "\n")
}

func (m mainModel) viewFileBrowser() string {
	if m.fileBrowse.path == "" {
		m.fileBrowse.path = "."
	}
	if m.fileBrowse.entries == nil {
		m.fileBrowse.refresh()
	}
	return m.renderFileBrowserContent()
}

func (m mainModel) renderFileBrowserContent() string {
	s := m.styles
	fb := &m.fileBrowse

	header := s.Title.Render("File Explorer")

	breadcrumb := s.DimText.Render(fmt.Sprintf("  %s", fb.path))

	// Error or search prompt
	var statusLine string
	if fb.err != "" {
		statusLine = s.StatusBad.Render("  " + fb.err)
	} else if fb.searchMode {
		statusLine = s.StatusWarn.Render(fmt.Sprintf("  Search: %s█  [n/N: next/prev match]", fb.searchQuery))
	}

	// Column headers
	colName := s.SectionItem.Render("  Name")
	colSize := s.SectionItem.Render(fmt.Sprintf("%*s", 10, "Size"))
	colModified := s.SectionItem.Render(fmt.Sprintf("%*s", 17, "Modified"))

	var rows []string
	for i, entry := range fb.entries {
		style := s.SectionItem
		prefix := "  "
		if i == fb.selected {
			style = s.MenuSelected
			prefix = "▸ "
		}

		icon := "📄"
		if entry.isDir {
			icon = "📁"
		}

		sizeStr := ""
		if !entry.isDir {
			sizeStr = fmt.Sprintf("%10s", formatFileSize(entry.size))
		} else {
			sizeStr = fmt.Sprintf("%10s", "<DIR>")
		}

		modStr := entry.modTime.Format("2006-01-02 15:04")

		ext := strings.ToLower(filepath.Ext(entry.name))
		isSupported := supportedExts[ext] || entry.isDir

		marker := ""
		if !entry.isDir && !isSupported {
			marker = s.StatusWarn.Render(" [unsupported]")
		}

		line := fmt.Sprintf("%s%s %-40s %s %s%s", prefix, icon, entry.name, sizeStr, modStr, marker)

		if i == fb.selected && fb.searchMode {
			q := strings.ToLower(fb.searchQuery)
			if q != "" && strings.Contains(strings.ToLower(entry.name), q) {
				line = s.StatusGood.Render("  > ") + line[3:]
			}
		}

		rows = append(rows, style.Render(line))
	}

	fileList := lipgloss.JoinVertical(lipgloss.Left, rows...)

	// Preview panel
	mainContent := fileList
	if fb.showPreview && fb.preview != "" && m.mainWidth() > 80 {
		leftWidth := m.mainWidth() * 3 / 5
		if leftWidth > 60 {
			leftWidth = 60
		}
		rightWidth := m.mainWidth() - leftWidth - 2

		previewTitle := s.Section.Render("Preview")
		previewBody := lipgloss.NewStyle().
			Width(rightWidth).
			MaxWidth(rightWidth).
			Render(fb.preview)
		previewContent := lipgloss.JoinVertical(lipgloss.Left, previewTitle, previewBody)

		// Truncate preview to viewport height
		rightPanel := s.BorderBox.
			Width(rightWidth + 4).
			Render(previewContent)

		mainContent = lipgloss.JoinHorizontal(lipgloss.Top,
			lipgloss.NewStyle().Width(leftWidth).Render(fileList),
			rightPanel,
		)
	}

	m.fileBrowse.err = ""

	return lipgloss.JoinVertical(lipgloss.Left,
		header, breadcrumb,
		statusLine,
		colName+"  "+colSize+"  "+colModified,
		mainContent,
	)
}

func (m mainModel) updateFileBrowser(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if !m.fileBrowse.searchMode && msg.String() == "h" {
			return m, func() tea.Msg { return navigateMsg{to: helpView} }
		}
		if !m.fileBrowse.searchMode && msg.String() == "enter" {
			if m.fileBrowse.selected >= 0 && m.fileBrowse.selected < len(m.fileBrowse.entries) {
				entry := m.fileBrowse.entries[m.fileBrowse.selected]
				if !entry.isDir {
					ext := strings.ToLower(filepath.Ext(entry.name))
					if supportedExts[ext] {
						m.fileBrowse.err = ""
						return m, func() tea.Msg { return fileSelectedMsg(entry.path) }
					}
				}
			}
		}
	}

	var cmd tea.Cmd
	m.fileBrowse, cmd = m.fileBrowse.Update(msg)
	return m, cmd
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
