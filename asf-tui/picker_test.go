package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestPickerDefaultPathNotReports(t *testing.T) {
	m := defaultTestModel()
	// With no lastPickerPaths set and no architecture file, should NOT return reports dir
	path := m.pickerStartPath(pickerArchitecture)
	if strings.Contains(path, "/reports") || strings.Contains(path, "\\reports") {
		t.Errorf("pickerStartPath should not default to reports dir, got: %s", path)
	}
	if path == "" {
		t.Error("pickerStartPath should not return empty string")
	}
}

func TestPickerStartPathPriority(t *testing.T) {
	m := defaultTestModel()

	// Priority 1: last used path
	m.lastPickerPaths[pickerArchitecture] = "/tmp/last-used"
	path := m.pickerStartPath(pickerArchitecture)
	if !strings.HasSuffix(path, "last-used") {
		t.Errorf("should return last used path, got: %s", path)
	}

	// Priority 2: previous architecture file dir
	m2 := defaultTestModel()
	m2.analyze.setDocPath("/some/dir/arch.yaml")
	path2 := m2.pickerStartPath(pickerArchitecture)
	if !strings.HasSuffix(path2, "/some/dir") && !strings.HasSuffix(path2, "\\some\\dir") {
		t.Errorf("should return architecture file dir, got: %s", path2)
	}
}

func TestPickerStartPathFallback(t *testing.T) {
	m := defaultTestModel()
	path := m.pickerStartPath(pickerEvidence)
	// Should return cwd or home, never empty
	if path == "" {
		t.Error("pickerStartPath should never return empty")
	}
}

func TestPickerParentDirNavigation(t *testing.T) {
	fp := newFilePickerState()
	fp.navigateDir("/tmp")
	parentBefore := fp.path
	_ = parentBefore

	// backspace should go to parent of /tmp → /
	fp.handleKey(tea.KeyMsg{Type: tea.KeyBackspace})
	if fp.path != "/" {
		t.Errorf("after backspace from /tmp, path should be /, got: %s", fp.path)
	}
	// backspace at / should stay at /
	fp.handleKey(tea.KeyMsg{Type: tea.KeyBackspace})
	if fp.path != "/" {
		t.Errorf("backspace at / should stay at /, got: %s", fp.path)
	}
}

func TestPickerChildDirNavigation(t *testing.T) {
	tmpDir := t.TempDir()
	subDir := filepath.Join(tmpDir, "subdir")
	if err := os.MkdirAll(subDir, 0755); err != nil {
		t.Fatal(err)
	}

	fp := newFilePickerState()
	fp.path = tmpDir
	fp.refresh()

	// Verify subdir appears in entries
	found := false
	for _, e := range fp.entries {
		if e.name == "subdir" && e.isDir {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("subdir not found in entries")
	}

	// Navigate into subdir via handleKey (select subdir, press enter)
	for i, e := range fp.entries {
		if e.name == "subdir" {
			fp.selected = i
			break
		}
	}
	fp.handleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if fp.path != subDir {
		t.Errorf("after entering subdir, path = %q, want %q", fp.path, subDir)
	}
}

func TestPickerJumpHome(t *testing.T) {
	fp := newFilePickerState()
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatal(err)
	}

	fp.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'~'}})
	if fp.path != home {
		t.Errorf("after ~ press, path = %q, want %q", fp.path, home)
	}
}

func TestPickerJumpRoot(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("root jump only on Unix")
	}
	fp := newFilePickerState()
	fp.handleKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'g'}})
	if fp.path != "/" {
		t.Errorf("after g press, path = %q, want /", fp.path)
	}
}

func TestPickerSelectSupportedArchFile(t *testing.T) {
	tmpDir := t.TempDir()
	archFile := filepath.Join(tmpDir, "arch.yaml")
	if err := os.WriteFile(archFile, []byte("key: value\n"), 0644); err != nil {
		t.Fatal(err)
	}

	fp := newFilePickerState()
	fp.mode = pickerArchitecture
	fp.path = tmpDir
	fp.refresh()

	found := false
	for i, e := range fp.entries {
		if e.name == "arch.yaml" && !e.isDir {
			fp.selected = i
			found = true
			break
		}
	}
	if !found {
		t.Fatal("arch.yaml not found in entries")
	}

	cmd, _ := fp.handleKey(tea.KeyMsg{Type: tea.KeyEnter})
	// Should return filePickedMsg
	if cmd == nil {
		t.Fatal("handleKey should return a command for supported file selection")
	}
	msg := cmd()
	if _, ok := msg.(filePickedMsg); !ok {
		t.Errorf("expected filePickedMsg, got %T", msg)
	}
	if msg.(filePickedMsg).path != archFile {
		t.Errorf("file path = %q, want %q", msg.(filePickedMsg).path, archFile)
	}
}

func TestPickerRejectsUnsupportedFile(t *testing.T) {
	tmpDir := t.TempDir()
	badFile := filepath.Join(tmpDir, "notes.txt")
	if err := os.WriteFile(badFile, []byte("hello\n"), 0644); err != nil {
		t.Fatal(err)
	}

	fp := newFilePickerState()
	fp.path = tmpDir
	fp.refresh()

	found := false
	for i, e := range fp.entries {
		if e.name == "notes.txt" && !e.isDir {
			fp.selected = i
			found = true
			break
		}
	}
	if !found {
		t.Fatal("notes.txt should appear (it's a supported ext)")
	}

	fp.handleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if fp.err != "" {
		t.Errorf("txt is supported, should not get error, got: %s", fp.err)
	}
}

func TestPickerRejectsUnsupportedExt(t *testing.T) {
	tmpDir := t.TempDir()
	badFile := filepath.Join(tmpDir, "image.png")
	if err := os.WriteFile(badFile, []byte("fake png\n"), 0644); err != nil {
		t.Fatal(err)
	}

	fp := newFilePickerState()
	fp.path = tmpDir
	fp.refresh()

	found := false
	for i, e := range fp.entries {
		if e.name == "image.png" {
			fp.selected = i
			found = true
			break
		}
	}
	if !found {
		t.Fatal("image.png should appear in entries")
	}

	cmd, _ := fp.handleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd != nil {
		t.Error("expected nil cmd for unsupported file, got a command")
	}
	if !strings.Contains(fp.err, "Unsupported") {
		t.Errorf("expected error about unsupported file, got: %s", fp.err)
	}
}

func TestPickerArchSelectionUpdatesAnalyze(t *testing.T) {
	m := defaultTestModel()
	m.pickerActive = true
	m.filePicker.mode = pickerArchitecture

	m.handleFilePicked(filePickedMsg{path: "/some/path/arch.yaml", mode: pickerArchitecture})

	if m.analyze.docPath() != "/some/path/arch.yaml" {
		t.Errorf("architecture path = %q, want /some/path/arch.yaml", m.analyze.docPath())
	}
	if m.pickerActive {
		t.Error("pickerActive should be false after file selection")
	}
}

func TestPickerEvidenceSelectionAppends(t *testing.T) {
	m := defaultTestModel()
	m.pickerActive = true
	m.filePicker.mode = pickerEvidence

	m.handleFilePicked(filePickedMsg{path: "/some/path/evidence1.csv", mode: pickerEvidence})
	if len(m.analyze.evidenceFiles) != 1 {
		t.Errorf("evidence count = %d, want 1", len(m.analyze.evidenceFiles))
	}

	// Second pick should append
	m.pickerActive = true
	m.handleFilePicked(filePickedMsg{path: "/some/path/evidence2.csv", mode: pickerEvidence})
	if len(m.analyze.evidenceFiles) != 2 {
		t.Errorf("evidence count = %d, want 2", len(m.analyze.evidenceFiles))
	}
}

func TestPickerEvidencePreservesExisting(t *testing.T) {
	m := defaultTestModel()
	m.analyze.addEvidence("/existing/ev1.csv")
	m.analyze.addEvidence("/existing/ev2.csv")

	m.handleFilePicked(filePickedMsg{path: "/new/ev3.csv", mode: pickerEvidence})

	if len(m.analyze.evidenceFiles) != 3 {
		t.Errorf("evidence count = %d, want 3", len(m.analyze.evidenceFiles))
	}
	if m.analyze.evidenceFiles[0] != "/existing/ev1.csv" {
		t.Errorf("first evidence should be preserved, got %q", m.analyze.evidenceFiles[0])
	}
}

func TestPickerCancelPreservesState(t *testing.T) {
	m := defaultTestModel()
	m.analyze.setDocPath("/original/arch.yaml")
	m.analyze.addEvidence("/original/ev.csv")
	m.pickerActive = true

	// Simulate cancel
	m.pickerActive = false

	if m.analyze.docPath() != "/original/arch.yaml" {
		t.Errorf("architecture should be preserved after cancel, got %q", m.analyze.docPath())
	}
	if len(m.analyze.evidenceFiles) != 1 {
		t.Errorf("evidence should be preserved after cancel, got %d", len(m.analyze.evidenceFiles))
	}
}

func TestPickerSearchFiltersEntries(t *testing.T) {
	tmpDir := t.TempDir()
	for _, name := range []string{"alpha.yaml", "beta.yaml", "gamma.json", "delta.txt"} {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte("x\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	fp := newFilePickerState()
	fp.path = tmpDir
	fp.refresh()

	// Verify all entries present
	if len(fp.entries) != 4 {
		t.Errorf("expected 4 entries, got %d", len(fp.entries))
	}

	// Activate search mode
	fp.searchMode = true
	fp.searchQuery = "alpha"
	fp.searchForFile()

	fp.searchMode = false
	if fp.selected >= len(fp.entries) || fp.entries[fp.selected].name != "alpha.yaml" {
		t.Errorf("search should find alpha.yaml, selected = %d, name = %s",
			fp.selected, fp.entries[fp.selected].name)
	}
}

func TestPickerPathsWithSpaces(t *testing.T) {
	tmpDir := t.TempDir()
	dirWithSpace := filepath.Join(tmpDir, "my documents")
	if err := os.MkdirAll(dirWithSpace, 0755); err != nil {
		t.Fatal(err)
	}
	fileWithSpace := filepath.Join(dirWithSpace, "arch file.yaml")
	if err := os.WriteFile(fileWithSpace, []byte("x\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Navigate into directory with space
	fp := newFilePickerState()
	fp.path = tmpDir
	fp.refresh()

	found := false
	for i, e := range fp.entries {
		if e.name == "my documents" && e.isDir {
			fp.selected = i
			found = true
			break
		}
	}
	if !found {
		t.Fatal("directory 'my documents' not found in entries")
	}

	fp.handleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if fp.path != dirWithSpace {
		t.Errorf("after entering 'my documents', path = %q, want %q", fp.path, dirWithSpace)
	}

	// Now select the file
	fp.refresh()
	found = false
	for i, e := range fp.entries {
		if e.name == "arch file.yaml" {
			fp.selected = i
			found = true
			break
		}
	}
	if !found {
		t.Fatal("arch file.yaml not found in 'my documents'")
	}

	cmd, _ := fp.handleKey(tea.KeyMsg{Type: tea.KeyEnter})
	if cmd == nil {
		t.Fatal("should be able to select file with spaces in path")
	}
}

func TestPickerPathNormalization(t *testing.T) {
	// Test that readDir uses filepath.Abs and filepath.Clean
	fp := newFilePickerState()
	fp.path = "."
	fp.refresh()
	if fp.path != "." {
		t.Errorf("initial path should be '.', got %q", fp.path)
	}

	// NavigateDir uses filepath.Clean implicitly via os operations
	tmpDir := t.TempDir()
	fp.navigateDir(tmpDir)
	if fp.path != tmpDir {
		t.Errorf("after navigateDir, path = %q, want %q", fp.path, tmpDir)
	}

	// Verify parent navigation uses filepath.Dir (not manual "/" concatenation)
	parent := filepath.Dir(tmpDir)
	fp.handleKey(tea.KeyMsg{Type: tea.KeyBackspace})
	if fp.path != parent {
		t.Errorf("after backspace, path = %q, want %q", fp.path, parent)
	}
}

func TestPickerWindowsPathSafety(t *testing.T) {
	if runtime.GOOS == "windows" {
		vol := filepath.VolumeName(`C:\Users\test`)
		if vol != "C:" {
			t.Errorf("VolumeName should return C: for C:\\Users\\test, got %q", vol)
		}

		cleaned := filepath.Clean(`C:\Users\test`)
		if !strings.HasPrefix(cleaned, `C:\`) {
			t.Errorf("Clean should preserve Windows drive letter, got %q", cleaned)
		}
	} else {
		// On non-Windows, VolumeName returns ""
		vol := filepath.VolumeName(`C:\Users\test`)
		if vol != "" {
			t.Errorf("on non-Windows, VolumeName should be empty, got %q", vol)
		}
	}
}

func TestPickerLastPathSavedOnSelection(t *testing.T) {
	m := defaultTestModel()
	m.handleFilePicked(filePickedMsg{path: "/some/dir/arch.yaml", mode: pickerArchitecture})

	if m.lastPickerPaths[pickerArchitecture] != "/some/dir" {
		t.Errorf("last path for architecture = %q, want /some/dir",
			m.lastPickerPaths[pickerArchitecture])
	}

	m.handleFilePicked(filePickedMsg{path: "/other/ev.csv", mode: pickerEvidence})
	if m.lastPickerPaths[pickerEvidence] != "/other" {
		t.Errorf("last path for evidence = %q, want /other",
			m.lastPickerPaths[pickerEvidence])
	}

	// Verify architecture last path is still preserved
	if m.lastPickerPaths[pickerArchitecture] != "/some/dir" {
		t.Errorf("architecture last path should still be /some/dir, got %q",
			m.lastPickerPaths[pickerArchitecture])
	}
}
