package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

func TestReportsLibrary_OpensAsLibrary(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = reportsView

	content := m.viewReports()
	if strings.Contains(content, "Select Export Format") {
		t.Error("Reports screen shows export format prompt (should show exported reports library)")
	}
	if strings.Contains(content, "Choose export format") {
		t.Error("Reports screen shows 'choose export format' (should show exported reports library)")
	}
}

func TestReportsLibrary_NoExportPromptByDefault(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = reportsView
	m.reportsV.entries = nil

	content := m.viewReports()
	if strings.Contains(content, "No exported reports yet.") {
		return
	}
	if strings.Contains(content, "No reports generated.") {
		t.Error("Reports screen shows old 'No reports generated.' message (should say 'No exported reports yet.')")
	}
}

func TestReportsLibrary_EmptyStateExplainsExport(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = reportsView
	m.reportsV.entries = nil

	content := m.viewReports()
	if !strings.Contains(content, "No exported reports yet.") {
		t.Error("Empty reports screen should show 'No exported reports yet.'")
	}
	if !strings.Contains(content, "Export a report from any case workspace") {
		t.Error("Empty reports screen should guide user to export from case workspace")
	}
	if strings.Contains(content, "Export") && strings.Contains(content, "Select Export Format") {
		t.Error("Empty reports screen should not show export format selection")
	}
}

func TestReportsLibrary_ListsExistingReports(t *testing.T) {
	dir := t.TempDir()
	orig := asfReportsDir
	asfReportsDir = func() string { return dir }
	defer func() { asfReportsDir = orig }()

	createFakeReport(t, dir, "test-report.pdf", 1024, time.Now().Add(-1*time.Hour))
	createFakeReport(t, dir, "another-report.json", 2048, time.Now())

	m := defaultTestModel()
	m.reportsV.entries = scanReportsDirs()

	content := m.viewReports()
	if !strings.Contains(content, "test-report.pdf") {
		t.Error("Reports screen should list test-report.pdf")
	}
	if !strings.Contains(content, "another-report.json") {
		t.Error("Reports screen should list another-report.json")
	}
}

func TestReportsLibrary_ReportMetadataDisplayed(t *testing.T) {
	dir := t.TempDir()
	orig := asfReportsDir
	asfReportsDir = func() string { return dir }
	defer func() { asfReportsDir = orig }()

	createFakeReport(t, dir, "test-report.pdf", 420*1024, time.Date(2026, 6, 15, 18, 20, 0, 0, time.UTC))

	m := defaultTestModel()
	m.reportsV.entries = scanReportsDirs()
	m.reportsV.selected = 0
	m.reportsV.mode = "detail"

	content := m.viewReports()
	if !strings.Contains(content, "test-report.pdf") {
		t.Error("Report detail should show file name")
	}
	if !strings.Contains(content, "PDF") {
		t.Error("Report detail should show format (PDF)")
	}
	if !strings.Contains(content, "Size") {
		t.Error("Report detail should show file size label")
	}
}

func TestReportsLibrary_RefreshReloadsReports(t *testing.T) {
	dir := t.TempDir()
	orig := asfReportsDir
	asfReportsDir = func() string { return dir }
	defer func() { asfReportsDir = orig }()

	m := defaultTestModel()
	m.reportsV.entries = scanReportsDirs()

	if len(m.reportsV.entries) != 0 {
		t.Error("Should start with empty reports")
	}

	createFakeReport(t, dir, "new-report.md", 512, time.Now())
	m.reportsV, _ = m.reportsV.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")})

	if len(m.reportsV.entries) != 1 {
		t.Errorf("After refresh should have 1 report, got %d", len(m.reportsV.entries))
	}
}

func TestReportsLibrary_SearchFiltersReports(t *testing.T) {
	m := defaultTestModel()
	m.reportsV.entries = []reportEntry{
		{name: "security-report.pdf", path: "/tmp/security-report.pdf", format: "PDF", modTime: time.Now(), size: 100, caseName: "security"},
		{name: "compliance-report.json", path: "/tmp/compliance-report.json", format: "JSON", modTime: time.Now(), size: 200, caseName: "compliance"},
	}

	m.reportsV.searchQuery = "security"

	content := m.viewReports()
	if !strings.Contains(content, "security-report.pdf") {
		t.Error("Filtered view should show matching report")
	}
	if strings.Contains(content, "compliance-report.json") {
		t.Error("Filtered view should not show non-matching report")
	}
}

func TestReportsLibrary_NavigateUpDown(t *testing.T) {
	m := defaultTestModel()
	m.reportsV.entries = []reportEntry{
		{name: "a.pdf", path: "/tmp/a.pdf", format: "PDF", modTime: time.Now(), size: 100},
		{name: "b.pdf", path: "/tmp/b.pdf", format: "PDF", modTime: time.Now(), size: 100},
		{name: "c.pdf", path: "/tmp/c.pdf", format: "PDF", modTime: time.Now(), size: 100},
	}

	m.reportsV.selected = 1
	m.reportsV, _ = m.reportsV.Update(tea.KeyMsg{Type: tea.KeyUp})
	if m.reportsV.selected != 0 {
		t.Error("up should decrement selected")
	}

	m.reportsV, _ = m.reportsV.Update(tea.KeyMsg{Type: tea.KeyDown})
	if m.reportsV.selected != 1 {
		t.Error("down should increment selected")
	}
}

func TestReportsLibrary_EnterTogglesDetail(t *testing.T) {
	m := defaultTestModel()
	m.reportsV.entries = []reportEntry{
		{name: "a.pdf", path: "/tmp/a.pdf", format: "PDF", modTime: time.Now(), size: 100},
	}
	m.reportsV.mode = "browse"

	m.reportsV, _ = m.reportsV.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if m.reportsV.mode != "detail" {
		t.Error("Enter should switch to detail mode")
	}

	m.reportsV, _ = m.reportsV.Update(tea.KeyMsg{Type: tea.KeyEnter})
	if m.reportsV.mode != "browse" {
		t.Error("Enter in detail mode should switch back to browse mode")
	}
}

func TestReportsLibrary_DeleteRemovesEntry(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "to-delete.pdf")
	os.WriteFile(path, []byte("test"), 0644)

	m := defaultTestModel()
	m.reportsV.entries = []reportEntry{
		{name: "to-delete.pdf", path: path, format: "PDF", modTime: time.Now(), size: 100},
	}
	m.reportsV.selected = 0
	m.reportsV.mode = "detail"

	// Press d to initiate delete confirmation
	m.reportsV, _ = m.reportsV.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")})
	if !m.reportsV.confirmDelete {
		t.Fatal("d should set confirmDelete")
	}

	// Press y to confirm
	m.reportsV, _ = m.reportsV.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("y")})

	if len(m.reportsV.entries) != 0 {
		t.Error("Delete should remove entry from list")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("Delete should remove file from disk")
	}
}

func TestExportAction_StillWorksFromCaseWorkspace(t *testing.T) {
	m := defaultTestModel()
	m.results.result = sampleAnalysisResult()
	m.router.currentView = caseView
	m.activeCase = "test.yaml"

	handled, model, _ := m.handleGlobalKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("e")})
	if !handled {
		t.Error("e in caseView should be handled")
	}
	m2 := model.(mainModel)
	if !m2.exportActive {
		t.Error("e in caseView should set exportActive")
	}
}

func TestExportAction_ExportedReportAppearsInReports(t *testing.T) {
	dir := t.TempDir()
	orig := asfReportsDir
	asfReportsDir = func() string { return dir }
	defer func() { asfReportsDir = orig }()
	origProject := projectReportsDir
	projectReportsDir = func() string { return dir }
	defer func() { projectReportsDir = origProject }()

	m := defaultTestModel()
	m.config.Output.Directory = dir
	m.results.result = sampleAnalysisResult()

	path, err := ExportResult(m.results.result, ExportJSON, dir)
	if err != nil {
		t.Fatalf("Export failed: %v", err)
	}
	if path == "" {
		t.Fatal("Export returned empty path")
	}

	m.reportsV.entries = scanReportsDirs()
	found := false
	for _, e := range m.reportsV.entries {
		if e.path == path {
			found = true
			break
		}
	}
	if !found {
		t.Error("Exported report should appear in Reports library")
	}
}

func TestReportsLibrary_DeleteCancelsOnN(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "cancel-delete.pdf")
	os.WriteFile(path, []byte("test"), 0644)

	m := defaultTestModel()
	m.reportsV.entries = []reportEntry{
		{name: "cancel-delete.pdf", path: path, format: "PDF", modTime: time.Now(), size: 100},
	}
	m.reportsV.selected = 0
	m.reportsV.mode = "detail"

	m.reportsV, _ = m.reportsV.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("d")})
	if !m.reportsV.confirmDelete {
		t.Fatal("d should set confirmDelete")
	}

	m.reportsV, _ = m.reportsV.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	if m.reportsV.confirmDelete {
		t.Error("n should cancel delete confirmation")
	}
	if len(m.reportsV.entries) != 1 {
		t.Error("Entry should remain after cancelled delete")
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("File should remain after cancelled delete")
	}
}

func TestReportsLibrary_ExportKeyInReportsShowsMessage(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = reportsView

	handled, model, _ := m.handleGlobalKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("e")})
	if !handled {
		t.Error("e in reportsView should be handled")
	}
	m2 := model.(mainModel)
	if m2.statusMsg == "" {
		t.Error("e in reportsView should set statusMsg")
	}
	if !strings.Contains(m2.statusMsg, "Export is available from an active case") {
		t.Error("statusMsg should guide user to export from case workspace")
	}
}

func TestFilePicker_Unrestricted(t *testing.T) {
	m := defaultTestModel()
	m.filePicker = newFilePickerState()
	m.filePicker.path = "/tmp"
	m.filePicker.mode = pickerArchitecture
	m.filePicker.refresh()

	if m.filePicker.path != "/tmp" {
		t.Error("File picker should allow navigation outside reports directory")
	}
}

func sampleAnalysisResult() *AnalysisResult {
	return &AnalysisResult{
		ArchitectureName: "test-arch",
		AnalysisDate:     time.Now(),
		Assumptions: []Assumption{
			{ID: "ASM-001", Description: "Test assumption", Risk: "Medium", Confidence: 0.85},
		},
		Summary: "Test summary",
	}
}

func createFakeReport(t *testing.T, dir, name string, size int64, modTime time.Time) {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, make([]byte, size), 0644); err != nil {
		t.Fatalf("Failed to create fake report: %v", err)
	}
	if err := os.Chtimes(path, modTime, modTime); err != nil {
		t.Fatalf("Failed to set mod time: %v", err)
	}
}
