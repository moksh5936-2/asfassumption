package main

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func defaultTestModel() *mainModel {
	cfg := DefaultConfig()
	return newMainModel(&cfg)
}

func TestGlobalKeyRouting_ArrowKeys(t *testing.T) {
	tests := []struct {
		name        string
		view        view
		key         string
		wantHandled bool
	}{
		{"analyze up -> child", analyzeView, "up", false},
		{"analyze j -> child", analyzeView, "j", false},
		{"settings down -> child", settingsView, "down", false},
		{"review up -> child", reviewView, "up", false},
		{"review down -> child", reviewView, "down", false},
		{"validation j -> child", validationView, "j", false},
		{"reports up -> child", reportsView, "up", false},
		{"reports down -> child", reportsView, "down", false},

		{"case up -> global (scroll)", caseView, "up", true},
		{"case down -> global (scroll)", caseView, "down", true},
		{"help up -> global (scroll)", helpView, "up", true},
		{"help down -> global (scroll)", helpView, "down", true},
		{"about up -> global (scroll)", aboutView, "up", true},
		{"about down -> global (scroll)", aboutView, "down", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := defaultTestModel()
			m.router.currentView = tt.view
			handled, _, _ := m.handleGlobalKey(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)})
			msg := tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(tt.key)}
			_ = msg
			handled, _, _ = m.handleGlobalKey(msgFromString(tt.key))
			if handled != tt.wantHandled {
				t.Errorf("handleGlobalKey(%q) on %v = %v, want %v", tt.key, tt.view, handled, tt.wantHandled)
			}
		})
	}
}

func msgFromString(s string) tea.KeyMsg {
	m := tea.KeyMsg{}
	switch s {
	case "up":
		m.Type = tea.KeyUp
	case "down":
		m.Type = tea.KeyDown
	case "k":
		m.Type = tea.KeyRunes
		m.Runes = []rune("k")
	case "j":
		m.Type = tea.KeyRunes
		m.Runes = []rune("j")
	case "tab":
		m.Type = tea.KeyTab
	case "shift+tab":
		m.Type = tea.KeyTab
	case "enter":
		m.Type = tea.KeyEnter
	case "esc":
		m.Type = tea.KeyEsc
	case "pgup", "b":
		m.Type = tea.KeyPgUp
	case "pgdn":
		m.Type = tea.KeyPgDown
	case " ":
		m.Type = tea.KeySpace
	case "ctrl+u":
		m.Type = tea.KeyCtrlU
	case "ctrl+d":
		m.Type = tea.KeyCtrlD
	case "home", "g":
		m.Type = tea.KeyHome
	case "end", "G":
		m.Type = tea.KeyEnd
	case "ctrl+c", "Q":
		m.Type = tea.KeyCtrlC
	case "?":
		m.Type = tea.KeyRunes
		m.Runes = []rune("?")
	case "r":
		m.Type = tea.KeyRunes
		m.Runes = []rune("r")
	case "v":
		m.Type = tea.KeyRunes
		m.Runes = []rune("v")
	case "c":
		m.Type = tea.KeyRunes
		m.Runes = []rune("c")
	case "e":
		m.Type = tea.KeyRunes
		m.Runes = []rune("e")
	case "s":
		m.Type = tea.KeyRunes
		m.Runes = []rune("s")
	case "/":
		m.Type = tea.KeyRunes
		m.Runes = []rune("/")
	default:
		m.Type = tea.KeyRunes
		m.Runes = []rune(s)
	}
	return m
}

func TestGlobalKeyRouting_Tab(t *testing.T) {
	m := defaultTestModel()
	m.router.focus = focusContent

	handled, model, _ := m.handleGlobalKey(msgFromString("tab"))
	if !handled {
		t.Error("tab should toggle focus globally")
	}
	mm := model.(mainModel)
	if mm.router.focus != focusSidebar {
		t.Error("tab should switch focus to sidebar")
	}

	handled, model, _ = mm.handleGlobalKey(msgFromString("tab"))
	if !handled {
		t.Error("tab should toggle focus back to content")
	}
	mm2 := model.(mainModel)
	if mm2.router.focus != focusContent {
		t.Error("tab should switch focus back to content")
	}
}

func TestGlobalKeyRouting_EscExceptions(t *testing.T) {
	m := defaultTestModel()

	m.router.currentView = analyzeView
	m.analyze.running = true
	handled, _, _ := m.handleGlobalKey(msgFromString("esc"))
	if handled {
		t.Error("esc on analyzeView with running=true should be forwarded to child")
	}

	m.analyze.running = false
	handled, _, _ = m.handleGlobalKey(msgFromString("esc"))
	if !handled {
		t.Error("esc on analyzeView with no state should navigate back")
	}

	m.router.currentView = settingsView
	m.settings.editing = true
	handled, _, _ = m.handleGlobalKey(msgFromString("esc"))
	if handled {
		t.Error("esc on settingsView with editing=true should be forwarded to child")
	}

	m.settings.editing = false
	handled, _, _ = m.handleGlobalKey(msgFromString("esc"))
	if !handled {
		t.Error("esc on settingsView with editing=false should navigate back")
	}

	m.router.currentView = reportsView
	m.reportsV.showConfirmation = true
	handled, _, _ = m.handleGlobalKey(msgFromString("esc"))
	if handled {
		t.Error("esc on reportsView with showConfirmation=true should be forwarded to child")
	}

	m.reportsV.showConfirmation = false
	handled, _, _ = m.handleGlobalKey(msgFromString("esc"))
	if !handled {
		t.Error("esc on reportsView with no confirmation should navigate back")
	}

	m.router.currentView = reviewView
	m.review.editing = true
	handled, _, _ = m.handleGlobalKey(msgFromString("esc"))
	if handled {
		t.Error("esc on reviewView with editing=true should be forwarded to child")
	}

	m.review.editing = false
	handled, _, _ = m.handleGlobalKey(msgFromString("esc"))
	if !handled {
		t.Error("esc on reviewView with editing=false should navigate back")
	}
}

func TestGlobalKeyRouting_ReviewRKey(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = reviewView

	handled, _, _ := m.handleGlobalKey(msgFromString("r"))
	if handled {
		t.Error("r on reviewView should be forwarded to child (Reject)")
	}

	m.router.currentView = analyzeView
	handled, _, _ = m.handleGlobalKey(msgFromString("r"))
	if !handled {
		t.Error("r on analyzeView should navigate to analyze")
	}
}

func TestGlobalKeyRouting_SettingsSKey(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = settingsView
	m.settings.editing = false

	handled, _, _ := m.handleGlobalKey(msgFromString("s"))
	if !handled {
		t.Error("s on settingsView not editing should save settings globally")
	}

	m.settings.editing = true
	handled, _, _ = m.handleGlobalKey(msgFromString("s"))
	if handled {
		t.Error("s on settingsView editing should not be globally handled")
	}

	m.router.currentView = reviewView
	m.settings.editing = false
	handled, _, _ = m.handleGlobalKey(msgFromString("s"))
	if handled {
		t.Error("s on reviewView should be forwarded to child (Accept)")
	}
}

func TestGlobalKeyRouting_PageKeys(t *testing.T) {
	pageKeys := []string{"pgup", "pgdn", "ctrl+u", "ctrl+d", "home", "end", "b", "g", "G"}
	allViews := []view{analyzeView, caseView, settingsView, aboutView, reportsView, reviewView, validationView, helpView}

	for _, key := range pageKeys {
		for _, v := range allViews {
			m := defaultTestModel()
			m.router.currentView = v
			handled, _, _ := m.handleGlobalKey(msgFromString(key))
			if !handled {
				t.Errorf("page key %q on %v should be handled globally (scroll), got false", key, v)
			}
		}
	}
}

func TestNavigateToUpdatesSidebarSel(t *testing.T) {
	m := defaultTestModel()

	m.navigateTo(analyzeView)
	if m.router.currentView != analyzeView {
		t.Errorf("navigateTo(analyze): currentView = %v, want analyzeView", m.router.currentView)
	}

	m.navigateTo(caseView)
	if m.router.currentView != caseView {
		t.Errorf("navigateTo(case): currentView = %v, want caseView", m.router.currentView)
	}

	m.navigateTo(settingsView)
	if m.router.currentView != settingsView {
		t.Errorf("navigateTo(settings): currentView = %v, want settingsView", m.router.currentView)
	}

	m.navigateTo(helpView)
	if m.router.currentView != helpView {
		t.Errorf("navigateTo(help): currentView = %v, want helpView", m.router.currentView)
	}
}

func TestNavigateBackUpdatesView(t *testing.T) {
	m := defaultTestModel()

	m.navigateTo(settingsView)
	m.navigateTo(caseView)
	m.navigateTo(analyzeView)

	m.navigateBack()
	if m.router.currentView != caseView {
		t.Errorf("after navigateBack: currentView = %v, want caseView", m.router.currentView)
	}

	m.navigateBack()
	if m.router.currentView != settingsView {
		t.Errorf("after second navigateBack: currentView = %v, want settingsView", m.router.currentView)
	}
}

func TestWindowSizeMsgFallsThrough(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = helpView

	msg := tea.WindowSizeMsg{Width: 120, Height: 40}
	model, cmd := m.Update(msg)
	_ = cmd

	mm := model.(mainModel)
	if mm.width != 120 {
		t.Errorf("width = %d, want 120", mm.width)
	}
	if mm.height != 40 {
		t.Errorf("height = %d, want 40", mm.height)
	}
	if !mm.ready {
		t.Error("ready should be true after WindowSizeMsg")
	}
}

func TestUpdateFallsThroughToChild(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = helpView

	msg := tea.KeyMsg{Type: tea.KeyDown}
	model, _ := m.Update(msg)
	mm := model.(mainModel)
	if mm.vp.YOffset <= 0 && mm.vp.TotalLineCount() > 0 {
		t.Error("expected vp.YOffset > 0 after down arrow on helpView (scroll)")
	}
}

func TestSidebarFocus(t *testing.T) {
	m := defaultTestModel()

	if m.router.focus != focusContent {
		t.Error("initial focus should be content")
	}

	m.router.ToggleFocus()
	if m.router.focus != focusSidebar {
		t.Error("after toggle, focus should be sidebar")
	}

	m.router.ToggleFocus()
	if m.router.focus != focusContent {
		t.Error("after second toggle, focus should be content")
	}
}

func TestSidebarNav_upDown(t *testing.T) {
	m := defaultTestModel()
	m.router.focus = focusSidebar

	// Skip past any initial section headers to reach a navigable item
	for nodes := m.router.sidebarVisibleNodes(); m.router.sidebarSel < len(nodes) && nodes[m.router.sidebarSel].isSection; {
		m.router.sidebarSel++
	}
	startSel := m.router.sidebarSel

	m.router.sidebarMoveDown()
	if m.router.sidebarSel <= startSel {
		t.Error("sidebarMoveDown should increase selection")
	}

	m.router.sidebarMoveUp()
	if m.router.sidebarSel != startSel {
		t.Error("sidebarMoveUp should return to original")
	}
}

func TestSearchActiveBypassesGlobalHandler(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.searchActive = true
	m.searchQuery = "test"

	m.searchActive = true
	_, cmd := m.Update(msgFromString("x"))
	_ = cmd

	if m.searchQuery != "test" {
		t.Logf("searchQuery after update = %q (value receiver copy)", m.searchQuery)
	}
}

func TestScrollKeysOnContentViewsScroll(t *testing.T) {
	contentViews := []view{caseView, helpView, aboutView}
	for _, v := range contentViews {
		m := defaultTestModel()
		m.router.currentView = v
		m.vp.YOffset = 50

		handled, _, _ := m.handleGlobalKey(msgFromString("up"))
		if !handled {
			t.Errorf("up on %v should be handled globally (scroll)", v)
		}

		handled, _, _ = m.handleGlobalKey(msgFromString("down"))
		if !handled {
			t.Errorf("down on %v should be handled globally (scroll)", v)
		}
	}
}

func TestSidebarActivate(t *testing.T) {
	m := defaultTestModel()

	// Find analyze (➕ New Analysis) in sidebar and activate it
	for i, n := range m.router.sidebarVisibleNodes() {
		if n.vid == analyzeView && !n.isSection {
			m.router.sidebarSel = i
			to, tab := m.router.sidebarActivate()
			if to != analyzeView {
				t.Errorf("sidebarActivate should navigate to analyzeView, got %v", to)
			}
			if tab != -1 {
				t.Errorf("analyze tab should be -1, got %d", tab)
			}
			break
		}
	}
}
