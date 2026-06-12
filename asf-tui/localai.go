package main

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InstalledModelInfo represents a model discovered from Ollama.
type InstalledModelInfo struct {
	Name       string
	Size       string
	ModifiedAt string
}

// CatalogEntry is a model in the ASF catalog.
type CatalogEntry struct {
	Info      ModelInfo
	Installed bool
	Active    bool
}

// OtherModel is a model found in Ollama but not in the ASF catalog.
type OtherModel struct {
	Info   InstalledModelInfo
	Active bool
}

type refreshResult struct {
	catalog       []CatalogEntry
	otherModels   []OtherModel
	ollamaRunning bool
	statusMsg     string
}

type aiRefreshMsg struct {
	result refreshResult
}

type localaiModel struct {
	catalog        []CatalogEntry
	otherModels    []OtherModel
	selected       int
	section        int // 0 = catalog, 1 = other, 2 = actions
	actionSelected int
	showActions    bool
	downloading    bool
	downloadModel  string
	downloadProg   float64
	downloadStage  string
	downloadErr    string
	statusMsg      string
	modelMgr       *ModelManager
	dlProgress     *syncProgress
	config         *Config
	ollamaRunning  bool
	ollamaChecked  bool
	activeModel    string
}

func newLocalAIModel(cfg *Config) localaiModel {
	catalog := make([]CatalogEntry, len(SupportedModels))
	installedSet := make(map[string]bool, len(cfg.AI.InstalledModels))
	for _, name := range cfg.AI.InstalledModels {
		installedSet[name] = true
	}
	for i, mi := range SupportedModels {
		catalog[i] = CatalogEntry{
			Info:      mi,
			Installed: installedSet[mi.Name],
			Active:    mi.Name == cfg.AI.ActiveModel,
		}
	}
	return localaiModel{
		modelMgr:    NewModelManager(),
		config:      cfg,
		activeModel: cfg.AI.ActiveModel,
		catalog:     catalog,
	}
}

// refreshFromOllama queries Ollama and merges with the ASF catalog.
func (m localaiModel) refreshFromOllama() refreshResult {
	var result refreshResult
	if !m.modelMgr.CheckAvailable() || !m.modelMgr.CheckRunning() {
		result.catalog = make([]CatalogEntry, len(SupportedModels))
		for i, mi := range SupportedModels {
			result.catalog[i] = CatalogEntry{
				Info:   mi,
				Active: mi.Name == m.activeModel,
			}
		}
		return result
	}

	result.ollamaRunning = true

	ollamaModels, err := m.modelMgr.ListInstalledAPI()
	if err != nil {
		result.statusMsg = fmt.Sprintf("Ollama API error: %v", err)
	}

	installed := make(map[string]InstalledModelInfo)
	for _, om := range ollamaModels {
		installed[om.Name] = InstalledModelInfo{
			Name:       om.Name,
			Size:       formatSize(om.Size),
			ModifiedAt: formatTime(om.ModifiedAt),
		}
	}

	result.catalog = make([]CatalogEntry, len(SupportedModels))
	for i, mi := range SupportedModels {
		_, inst := installed[mi.Name]
		delete(installed, mi.Name)
		result.catalog[i] = CatalogEntry{
			Info:      mi,
			Installed: inst,
			Active:    mi.Name == m.activeModel && inst,
		}
	}

	for _, om := range ollamaModels {
		if _, found := installed[om.Name]; found {
			result.otherModels = append(result.otherModels, OtherModel{
				Info: InstalledModelInfo{
					Name:       om.Name,
					Size:       formatSize(om.Size),
					ModifiedAt: formatTime(om.ModifiedAt),
				},
				Active: om.Name == m.activeModel,
			})
		}
	}
	return result
}

func (m *localaiModel) applyRefresh(r refreshResult) {
	if r.ollamaRunning {
		m.catalog = r.catalog
		m.otherModels = r.otherModels
	}
	m.ollamaRunning = r.ollamaRunning
	if r.statusMsg != "" {
		m.statusMsg = r.statusMsg
	}
}

func formatSize(bytes int64) string {
	if bytes == 0 {
		return ""
	}
	switch {
	case bytes > 1<<30:
		return fmt.Sprintf("%.1f GB", float64(bytes)/float64(1<<30))
	case bytes > 1<<20:
		return fmt.Sprintf("%.1f MB", float64(bytes)/float64(1<<20))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func formatTime(t string) string {
	if t == "" {
		return ""
	}
	parsed, err := time.Parse(time.RFC3339, t)
	if err != nil {
		if len(t) >= 10 {
			return t[:10]
		}
		return t
	}
	return parsed.Format("Jan 2, 2006")
}

func (m localaiModel) Update(msg tea.Msg) (localaiModel, tea.Cmd) {
	if !m.ollamaChecked {
		m.ollamaChecked = true
		return m, m.startRefreshCmd()
	}

	switch msg := msg.(type) {
	case aiRefreshMsg:
		m.applyRefresh(msg.result)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.downloading {
				break
			}
			m.moveUp()
		case "down", "j":
			if m.downloading {
				break
			}
			m.moveDown()
		case "enter":
			if m.downloading {
				break
			}
			if m.showActions {
				return m.executeAction()
			}
			if m.section == 0 && m.selected < len(m.catalog) {
				m.showActions = true
				m.actionSelected = 0
			} else if m.section == 1 {
				m.showActions = true
				m.actionSelected = 0
			}
		case "esc":
			if m.showActions {
				m.showActions = false
			}
		}

	case aiDownloadTickMsg:
		if !m.downloading || m.dlProgress == nil {
			return m, nil
		}
		pct, stage, done, err := m.dlProgress.Snapshot()
		m.downloadProg = pct
		m.downloadStage = stage
		if err != "" {
			m.downloading = false
			m.downloadErr = err
			m.statusMsg = fmt.Sprintf("Download failed: %s", err)
			return m, nil
		}
		if done {
			m.downloading = false
			m.downloadProg = 100
			m.applyRefresh(m.refreshFromOllama())
			m.statusMsg = fmt.Sprintf("✓ %s installed", m.downloadModel)
			m.showActions = false
			m.config.AI.InstalledModels = nil
			for _, c := range m.catalog {
				if c.Installed {
					m.config.AI.InstalledModels = append(m.config.AI.InstalledModels, c.Info.Name)
				}
			}
			for _, o := range m.otherModels {
				m.config.AI.InstalledModels = append(m.config.AI.InstalledModels, o.Info.Name)
			}
			if err := m.config.Save(ConfigPath()); err != nil {
				debugLog.Printf("save config after download: %v", err)
			}
			return m, nil
		}
		return m, m.pollDownloadCmd()
	}
	return m, nil
}

func (m *localaiModel) updateSection() {
	if m.selected < len(m.catalog) {
		m.section = 0
	} else if len(m.otherModels) > 0 && m.selected > len(m.catalog) {
		m.section = 1
	}
}

func (m *localaiModel) moveUp() {
	if m.showActions {
		if m.actionSelected > 0 {
			m.actionSelected--
		}
		return
	}
	if m.selected > 0 {
		m.selected--
		m.updateSection()
		return
	}
	// Wrap to end
	total := m.totalItems()
	if total > 0 {
		m.selected = total - 1
		m.updateSection()
	}
}

func (m *localaiModel) moveDown() {
	if m.showActions {
		actions := m.actionsForModel()
		if m.actionSelected < len(actions)-1 {
			m.actionSelected++
		}
		return
	}
	total := m.totalItems()
	if m.selected < total-1 {
		m.selected++
		m.updateSection()
		return
	}
	m.selected = 0
	m.updateSection()
}

func (m *localaiModel) totalItems() int {
	n := len(m.catalog)
	if len(m.otherModels) > 0 {
		n++ // section gap
		n += len(m.otherModels)
	}
	return n
}

func (m localaiModel) actionsForModel() []string {
	actions := []string{"Download", "Set as Active"}
	if m.section == 0 && m.selected < len(m.catalog) {
		if m.catalog[m.selected].Installed {
			actions = append(actions, "Delete")
		}
	}
	if m.section == 1 {
		actions = []string{"Set as Active", "Delete"}
	}
	return actions
}

func (m localaiModel) executeAction() (localaiModel, tea.Cmd) {
	var modelName string

	if m.section == 0 && m.selected < len(m.catalog) {
		modelName = m.catalog[m.selected].Info.Name
	} else if m.section == 1 {
		idx := m.selected - len(m.catalog) - 1
		if idx >= 0 && idx < len(m.otherModels) {
			modelName = m.otherModels[idx].Info.Name
		}
	}

	if modelName == "" {
		return m, nil
	}

	actions := m.actionsForModel()
	if m.actionSelected >= len(actions) {
		return m, nil
	}

	switch actions[m.actionSelected] {
	case "Download":
		if !m.modelMgr.CheckAvailable() {
			m.statusMsg = "Ollama not found. Install from https://ollama.ai"
			return m, nil
		}
		if !m.modelMgr.CheckRunning() {
			m.statusMsg = "Ollama is not running. Start it with: ollama serve"
			return m, nil
		}
		m.downloadModel = modelName
		m.downloading = true
		m.downloadProg = 0
		m.downloadStage = "Starting download..."
		m.downloadErr = ""

		sp := &syncProgress{}
		m.dlProgress = sp
		go m.modelMgr.StartDownload(modelName, sp)

		return m, m.pollDownloadCmd()

	case "Set as Active":
		m.config.AI.ActiveModel = modelName
		m.config.AI.Enabled = true
		m.activeModel = modelName
		m.applyRefresh(m.refreshFromOllama())
		if err := m.config.Save(ConfigPath()); err != nil {
			m.statusMsg = fmt.Sprintf("Failed to save config: %v", err)
		} else {
			m.statusMsg = fmt.Sprintf("✓ Active model set to %s", modelName)
		}
		m.showActions = false

	case "Delete":
		delErr := m.modelMgr.DeleteModel(modelName)
		m.applyRefresh(m.refreshFromOllama())
		if delErr != nil {
			m.statusMsg = fmt.Sprintf("Delete failed: %v", delErr)
		} else {
			m.statusMsg = fmt.Sprintf("✓ %s deleted", modelName)
		}
		m.config.AI.InstalledModels = nil
		for _, c := range m.catalog {
			if c.Installed {
				m.config.AI.InstalledModels = append(m.config.AI.InstalledModels, c.Info.Name)
			}
		}
		for _, o := range m.otherModels {
			m.config.AI.InstalledModels = append(m.config.AI.InstalledModels, o.Info.Name)
		}
		if m.config.AI.ActiveModel == modelName {
			m.config.AI.ActiveModel = ""
			m.config.AI.Enabled = false
			m.activeModel = ""
		}
		if err := m.config.Save(ConfigPath()); err != nil {
			debugLog.Printf("save config after delete: %v", err)
		}
		m.showActions = false
	}
	return m, nil
}

type aiDownloadTickMsg struct{}

func (m localaiModel) pollDownloadCmd() tea.Cmd {
	return tea.Tick(400*time.Millisecond, func(t time.Time) tea.Msg {
		return aiDownloadTickMsg{}
	})
}

func (m localaiModel) startRefreshCmd() tea.Cmd {
	return func() tea.Msg {
		return aiRefreshMsg{result: m.refreshFromOllama()}
	}
}

func (m mainModel) viewLocalAI() string {
	s := m.styles
	lm := m.localai

	if lm.downloading {
		barWidth := 40
		filled := int(float64(barWidth) * lm.downloadProg / 100.0)
		bar := ""
		for i := 0; i < barWidth; i++ {
			if i < filled {
				bar += "█"
			} else {
				bar += "░"
			}
		}
		return lipgloss.JoinVertical(lipgloss.Left,
			s.Title.Render("Local AI Models"),
			s.BorderBox.Render(
				lipgloss.JoinVertical(lipgloss.Center,
					s.SectionItem.Render(fmt.Sprintf("Downloading: %s", lm.downloadModel)),
					s.ProgressBar.Render(bar),
					s.Value.Render(fmt.Sprintf("%.0f%%", lm.downloadProg)),
					s.SectionItem.Render(lm.downloadStage),
				),
			),
		)
	}

	// Ollama status header
	ollamaStatus := s.StatusWarn.Render("Ollama is not running")
	if lm.ollamaRunning {
		ollamaStatus = s.StatusGood.Render("Ollama is running")
	}
	if !lm.ollamaChecked {
		ollamaStatus = s.SectionItem.Render("Checking Ollama...")
	}

	var rows []string
	rows = append(rows,
		s.Title.Render("Local AI Models"),
		s.Subtitle.Render("Download and manage local AI models for enhanced analysis"),
		s.BorderBox.Render(
			lipgloss.JoinVertical(lipgloss.Left,
				ollamaStatus,
			),
		),
	)

	if !lm.ollamaRunning && lm.ollamaChecked {
		rows = append(rows, "",
			s.StatusWarn.Render("  Ollama is not running. Install from https://ollama.ai and start with: ollama serve"),
		)
	}

	// Recommended ASF Models section
	var catalogItems []string
	for i, entry := range lm.catalog {
		style := s.SectionItem
		prefix := "  "
		if !lm.showActions && lm.section == 0 && i == lm.selected {
			style = s.MenuSelected
			prefix = "▸ "
		}
		status := ""
		if entry.Installed {
			status = s.StatusGood.Render(" ✓ Installed")
		}
		if entry.Active {
			status += s.StatusGood.Render(" [ACTIVE]")
		}
		catalogItems = append(catalogItems, style.Render(fmt.Sprintf("%s%s  (%s)%s",
			prefix, entry.Info.Display, entry.Info.Size, status)))
	}

	if len(catalogItems) > 0 {
		rows = append(rows, "",
			s.Section.Render("Recommended ASF Models"),
			lipgloss.JoinVertical(lipgloss.Left, catalogItems...),
		)
	}

	// Other installed models section
	if len(lm.otherModels) > 0 {
		var otherItems []string
		for i, om := range lm.otherModels {
			idx := len(lm.catalog) + 1 + i
			style := s.SectionItem
			prefix := "  "
			if !lm.showActions && lm.section == 1 && idx == lm.selected {
				style = s.MenuSelected
				prefix = "▸ "
			}
			status := ""
			if om.Active {
				status = s.StatusGood.Render(" [ACTIVE]")
			}
			line := fmt.Sprintf("%s%s", prefix, om.Info.Name)
			if om.Info.Size != "" {
				line += fmt.Sprintf("  (%s)", om.Info.Size)
			}
			if om.Info.ModifiedAt != "" {
				line += fmt.Sprintf("  %s", om.Info.ModifiedAt)
			}
			otherItems = append(otherItems, style.Render(line+status))
		}
		rows = append(rows, "",
			s.Section.Render("Other Installed Ollama Models"),
			lipgloss.JoinVertical(lipgloss.Left, otherItems...),
		)
	}

	// Action menu
	if lm.showActions && !lm.downloading {
		var actionItems []string
		actions := lm.actionsForModel()
		for i, action := range actions {
			style := s.SectionItem
			prefix := "  "
			if i == lm.actionSelected {
				style = s.MenuSelected
				prefix = "▸ "
			}
			actionItems = append(actionItems, style.Render(prefix+action))
		}

		var modelDisplay string
		if lm.section == 0 && lm.selected < len(lm.catalog) {
			modelDisplay = lm.catalog[lm.selected].Info.Display
		} else {
			idx := lm.selected - len(lm.catalog) - 1
			if idx >= 0 && idx < len(lm.otherModels) {
				modelDisplay = lm.otherModels[idx].Info.Name
			}
		}
		if modelDisplay != "" {
			rows = append(rows, "",
				s.Section.Render(fmt.Sprintf("Actions for %s:", modelDisplay)),
				lipgloss.JoinVertical(lipgloss.Left, actionItems...),
			)
		}
	}

	if lm.statusMsg != "" {
		rows = append(rows, "", s.StatusGood.Render("  "+lm.statusMsg))
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (m mainModel) updateLocalAI(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.localai, cmd = m.localai.Update(msg)
	return m, cmd
}
