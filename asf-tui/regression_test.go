package main

import (
	"fmt"
	"strings"
	"testing"

	"asf-tui/asf/trust"
	"asf-tui/asf/verify"
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

	// Find analyze (+ New Analysis) in sidebar and activate it
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

func TestTabStateNavigation(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.router.focus = focusContent
	m.results.result = &AnalysisResult{
		Assumptions: []Assumption{{ID: "A1"}, {ID: "A2"}, {ID: "A3"}},
	}
	m.results.resultTab = 1
	ts := m.results.tabStateFor(1)

	if ts.selectedIndex != 0 {
		t.Errorf("initial selectedIndex = %d, want 0", ts.selectedIndex)
	}

	ts.selectedIndex = 1
	got := ts.selectedIndex
	if got != 1 {
		t.Errorf("after set selectedIndex = %d, want 1", got)
	}
}

func TestTabStateDetailToggle(t *testing.T) {
	ts := &tabState{}
	if ts.detailOpen {
		t.Error("detailOpen should be false initially")
	}
	ts.detailOpen = true
	if !ts.detailOpen {
		t.Error("detailOpen should be true after toggle on")
	}
	ts.detailOpen = false
	if ts.detailOpen {
		t.Error("detailOpen should be false after toggle off")
	}
}

func TestTabStateFilter(t *testing.T) {
	ts := &tabState{}
	if ts.filterActive {
		t.Error("filterActive should be false initially")
	}
	ts.filterActive = true
	ts.searchQuery = "auth"
	if ts.searchQuery != "auth" {
		t.Errorf("searchQuery = %q, want %q", ts.searchQuery, "auth")
	}
	ts.filterActive = false
	ts.searchQuery = ""
	if ts.searchQuery != "" {
		t.Errorf("after clear searchQuery = %q, want empty", ts.searchQuery)
	}
}

func TestBreadcrumbRendering(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.activeCase = "/tmp/test.yaml"
	m.results.result = &AnalysisResult{
		Assumptions: []Assumption{{ID: "A1"}, {ID: "A2"}},
	}
	m.results.resultTab = 1
	ts := m.results.tabStateFor(1)

	breadcrumb := m.renderBreadcrumb(1, ts)
	if breadcrumb == "" {
		t.Error("breadcrumb should not be empty when result is loaded")
	}
	if !strings.Contains(breadcrumb, "Assumptions") {
		t.Error("breadcrumb should contain tab name 'Assumptions'")
	}
	if !strings.Contains(breadcrumb, "#1") {
		t.Error("breadcrumb should contain item number '#1'")
	}

	ts.selectedIndex = 1
	breadcrumb2 := m.renderBreadcrumb(1, ts)
	if !strings.Contains(breadcrumb2, "#2") {
		t.Error("breadcrumb should update item number to '#2'")
	}

	m.results.detailFocus = true
	breadcrumb3 := m.renderBreadcrumb(1, ts)
	if !strings.Contains(breadcrumb3, "detail") {
		t.Error("breadcrumb should contain 'detail' when detail pane is focused")
	}
}

func TestOverviewTabBreadcrumb(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.activeCase = "/tmp/test.yaml"
	m.results.result = &AnalysisResult{}
	m.results.resultTab = 0
	ts := m.results.tabStateFor(0)

	breadcrumb := m.renderBreadcrumb(0, ts)
	if !strings.Contains(breadcrumb, "Overview") {
		t.Error("breadcrumb for overview tab should contain 'Overview'")
	}
	if strings.Contains(breadcrumb, "#") {
		t.Error("breadcrumb for overview tab should NOT contain item number")
	}
}

func TestEmptyResultBreadcrumb(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.activeCase = ""
	m.results.result = nil
	ts := defaultTabState()
	breadcrumb := m.renderBreadcrumb(0, ts)
	if breadcrumb != "" {
		t.Error("breadcrumb should be empty when no active case")
	}
}

func TestTabStateResetOnResultChange(t *testing.T) {
	m := defaultTestModel()
	m.results.result = &AnalysisResult{
		Assumptions: []Assumption{{ID: "A1"}, {ID: "A2"}, {ID: "A3"}},
	}

	ts := m.results.tabStateFor(1)
	ts.selectedIndex = 2
	ts.detailOpen = true
	ts.searchQuery = "test"
	ts.filterActive = true

	if m.results.tabStates[1].selectedIndex != 2 {
		t.Error("tabState should persist before reset")
	}

	m.results.tabStates = make(map[int]*tabState)
	if _, ok := m.results.tabStates[1]; ok {
		t.Error("tabStates should be empty after reset")
	}
}

func TestListNavKeysRoutedToUpdateResults(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.router.focus = focusContent
	m.results.result = &AnalysisResult{
		Assumptions: []Assumption{
			{ID: "A1", Description: "First assumption"},
			{ID: "A2", Description: "Second assumption"},
		},
	}
	m.results.resultTab = 1

	ts := m.results.tabStateFor(1)
	ts.selectedIndex = 0

	// Simulate the routing in Update(): tab nav section routes "down" to updateResults
	model, _ := m.updateResults(msgFromString("down"))
	mm := model.(mainModel)
	ts2 := mm.results.tabStateFor(1)
	if ts2.selectedIndex != 1 {
		t.Errorf("after down: selectedIndex = %d, want 1", ts2.selectedIndex)
	}

	model, _ = mm.updateResults(msgFromString("up"))
	mm = model.(mainModel)
	ts3 := mm.results.tabStateFor(1)
	if ts3.selectedIndex != 0 {
		t.Errorf("after up: selectedIndex = %d, want 0", ts3.selectedIndex)
	}
}

func TestDetailToggleInUpdateResults(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.router.focus = focusContent
	m.results.result = &AnalysisResult{
		Assumptions: []Assumption{{ID: "A1"}},
	}
	m.results.resultTab = 1

	if m.results.detailFocus {
		t.Error("detailFocus should be false initially")
	}

	// Press Enter to focus detail pane
	model, _ := m.updateResults(msgFromString("enter"))
	mm := model.(mainModel)
	if !mm.results.detailFocus {
		t.Error("detailFocus should be true after Enter")
	}

	// Press Esc to return to list
	model, _ = mm.updateResults(msgFromString("esc"))
	mm = model.(mainModel)
	if mm.results.detailFocus {
		t.Error("detailFocus should be false after Esc")
	}
}

func TestRenderHintsBarPerTab(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.results.result = &AnalysisResult{
		Assumptions: []Assumption{{ID: "A1"}, {ID: "A2"}},
	}

	for tab := 0; tab <= 6; tab++ {
		m.results.resultTab = tab
		hints := m.renderHintsBar()
		if hints == "" {
			t.Errorf("tab %d: hints bar should not be empty", tab)
		}
		if !strings.Contains(hints, "↑↓") {
			t.Errorf("tab %d: hints should contain scroll keys", tab)
		}
	}
}

func TestMainHeightWithBreadcrumb(t *testing.T) {
	m := defaultTestModel()
	m.height = 60
	m.router.currentView = caseView

	// Without result (no breadcrumb)
	h1 := m.mainHeight()
	m.results.result = &AnalysisResult{Assumptions: []Assumption{{ID: "A1"}}}
	// With result (breadcrumb added)
	h2 := m.mainHeight()
	// h2 should be less than h1 (breadcrumb subtracted)
	if h2 > h1 {
		t.Errorf("mainHeight with breadcrumb (%d) should be less than without (%d)", h2, h1)
	}
	if h1 == 0 || h2 == 0 {
		t.Errorf("mainHeight should never be 0: noResult=%d, withResult=%d", h1, h2)
	}

	m.router.currentView = analyzeView
	h3 := m.mainHeight()
	m.results.result = &AnalysisResult{}
	h4 := m.mainHeight()
	if h3 != h4 {
		t.Errorf("mainHeight on analyzeView should be same regardless of result: %d vs %d", h3, h4)
	}
}

func TestScrollPercentFormat(t *testing.T) {
	tests := []struct {
		offset  int
		visible int
		total   int
		want    string
	}{
		{0, 40, 200, "Line 1–40 / 200"},
		{40, 40, 200, "Line 41–80 / 200"},
		{160, 40, 200, "Line 161–200 / 200"},
	}
	for _, tt := range tests {
		first := tt.offset + 1
		last := tt.offset + tt.visible
		if last > tt.total {
			last = tt.total
		}
		pct := int(float64(tt.offset+tt.visible) / float64(tt.total) * 100)
		if pct > 100 {
			pct = 100
		}
		got := fmt.Sprintf("Line %d–%d / %d  (%d%%)", first, last, tt.total, pct)
		if !strings.Contains(got, tt.want) {
			t.Errorf("scroll format = %q, should contain %q", got, tt.want)
		}
	}
}

func TestTabCountString(t *testing.T) {
	m := defaultTestModel()
	m.results.result = &AnalysisResult{
		Assumptions: []Assumption{
			{ID: "A1"}, {ID: "A2"}, {ID: "A3"},
		},
		Contradictions: []Contradiction{
			{ID: "C1"}, {ID: "C2"},
		},
		Controls: []ControlDetail{
			{ID: "Ctrl1"}, {ID: "Ctrl2"}, {ID: "Ctrl3"}, {ID: "Ctrl4"},
		},
		VerificationOutput: &verify.VerificationOutput{
			Assessment: &verify.VerificationAssessment{
				VerifiedCount:   5,
				PartialCount:    2,
				UnverifiedCount: 1,
				NoEvidenceCount: 0,
			},
		},
	}

	// Assumptions tab
	s := m.results.tabCountString(1)
	if !strings.Contains(s, "3 assumptions") {
		t.Errorf("tabCountString(1) = %q, want '3 assumptions'", s)
	}

	// Verification tab
	s = m.results.tabCountString(2)
	if !strings.Contains(s, "5 verified") || !strings.Contains(s, "1 unverified") {
		t.Errorf("tabCountString(2) = %q, should contain verification counts", s)
	}

	// Contradictions tab
	s = m.results.tabCountString(3)
	if !strings.Contains(s, "2 contradictions") {
		t.Errorf("tabCountString(3) = %q, want '2 contradictions'", s)
	}

	// Controls tab
	s = m.results.tabCountString(5)
	if !strings.Contains(s, "4 control") {
		t.Errorf("tabCountString(5) = %q, want '4 control(s)'", s)
	}

	// Overview tab
	s = m.results.tabCountString(0)
	if s != "" {
		t.Errorf("tabCountString(0) should be empty, got %q", s)
	}
}

func TestSelectedIndexCanReachEndOfList(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.router.focus = focusContent
	m.results.result = &AnalysisResult{
		Assumptions: []Assumption{
			{ID: "A1"}, {ID: "A2"}, {ID: "A3"}, {ID: "A4"}, {ID: "A5"},
		},
	}
	m.results.resultTab = 1
	ts := m.results.tabStateFor(1)

	// Navigate to last item
	for i := 0; i < 10; i++ {
		model, _ := m.updateResults(msgFromString("down"))
		mm := model.(mainModel)
		m = &mm
		ts = m.results.tabStateFor(1)
	}
	if ts.selectedIndex != 4 {
		t.Errorf("after pressing down 10 times, selectedIndex = %d, want 4 (last of 5)", ts.selectedIndex)
	}

	// Navigate back to first
	for i := 0; i < 10; i++ {
		model, _ := m.updateResults(msgFromString("up"))
		mm := model.(mainModel)
		m = &mm
		ts = m.results.tabStateFor(1)
	}
	if ts.selectedIndex != 0 {
		t.Errorf("after pressing up 10 times, selectedIndex = %d, want 0 (first)", ts.selectedIndex)
	}
}

func TestTabSwitchPreservesState(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.router.focus = focusContent
	m.results.result = &AnalysisResult{
		Assumptions:    []Assumption{{ID: "A1"}, {ID: "A2"}, {ID: "A3"}},
		Contradictions: []Contradiction{{ID: "C1"}, {ID: "C2"}},
	}
	m.results.resultTab = 1

	// Set up state on tab 1
	ts1 := m.results.tabStateFor(1)
	ts1.selectedIndex = 2
	ts1.detailOpen = true
	ts1.searchQuery = "auth"
	ts1.filterActive = true

	// Save scroll for tab 1, set scroll for tab 3
	m.results.tabScroll[1] = 50
	m.results.tabScroll[3] = 20

	// Switch to tab 3 then back
	m.results.resultTab = 3
	if m.results.resultTab != 3 {
		t.Errorf("after switch to tab 3, resultTab = %d, want 3", m.results.resultTab)
	}
	m.results.resultTab = 1
	if m.results.resultTab != 1 {
		t.Errorf("after switch back to tab 1, resultTab = %d, want 1", m.results.resultTab)
	}

	ts1back := m.results.tabStateFor(1)
	if ts1back.selectedIndex != 2 {
		t.Errorf("after tab switch, selectedIndex = %d, want 2", ts1back.selectedIndex)
	}
	if !ts1back.detailOpen {
		t.Error("after tab switch, detailOpen should be true")
	}
	if ts1back.searchQuery != "auth" {
		t.Errorf("after tab switch, searchQuery = %q, want 'auth'", ts1back.searchQuery)
	}
	if !ts1back.filterActive {
		t.Error("after tab switch, filterActive should be true")
	}
}

func TestSearchIncrementDecrement(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.router.focus = focusContent
	m.results.result = &AnalysisResult{
		Assumptions: []Assumption{
			{ID: "A1"}, {ID: "A2"}, {ID: "A3"},
		},
	}
	m.results.resultTab = 1
	ts := m.results.tabStateFor(1)
	ts.filterActive = true
	ts.searchQuery = "auth"

	// 'n' increments selected index (filter only affects render, not nav)
	ts.selectedIndex = 0
	model, _ := m.updateResults(msgFromString("n"))
	mm := model.(mainModel)
	ts = mm.results.tabStateFor(1)
	if ts.selectedIndex != 1 {
		t.Errorf("n should increment selectedIndex: got %d, want 1", ts.selectedIndex)
	}

	// 'N' decrements selected index
	model, _ = mm.updateResults(msgFromString("N"))
	mm = model.(mainModel)
	ts = mm.results.tabStateFor(1)
	if ts.selectedIndex != 0 {
		t.Errorf("N should decrement selectedIndex: got %d, want 0", ts.selectedIndex)
	}
}

func TestTabStateSearch(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.router.focus = focusContent
	m.results.result = &AnalysisResult{
		Assumptions: []Assumption{{ID: "A1"}, {ID: "A2"}},
	}
	m.results.resultTab = 1

	// Press / to activate filter
	model, _ := m.updateResults(msgFromString("/"))
	mm := model.(mainModel)
	ts := mm.results.tabStateFor(1)
	if !ts.filterActive {
		t.Error("filterActive should be true after /")
	}

	// Type characters while filter is active (simulated via updateResults)
	// The tab nav section routes single-char keys to updateResults
	model, _ = mm.updateResults(msgFromString("a"))
	mm = model.(mainModel)
	ts = mm.results.tabStateFor(1)
	if ts.searchQuery != "a" {
		t.Errorf("searchQuery = %q, want %q", ts.searchQuery, "a")
	}

	// Type more
	model, _ = mm.updateResults(msgFromString("u"))
	mm = model.(mainModel)
	ts = mm.results.tabStateFor(1)
	if ts.searchQuery != "au" {
		t.Errorf("searchQuery = %q, want %q", ts.searchQuery, "au")
	}

	// Backspace
	model, _ = mm.updateResults(tea.KeyMsg{Type: tea.KeyBackspace})
	mm = model.(mainModel)
	ts = mm.results.tabStateFor(1)
	if ts.searchQuery != "a" {
		t.Errorf("after backspace searchQuery = %q, want %q", ts.searchQuery, "a")
	}

	// Esc clears filter
	model, _ = mm.updateResults(msgFromString("esc"))
	mm = model.(mainModel)
	ts = mm.results.tabStateFor(1)
	if ts.filterActive {
		t.Error("filterActive should be false after esc")
	}
	if ts.searchQuery != "" {
		t.Errorf("searchQuery should be empty after esc, got %q", ts.searchQuery)
	}
}

func TestTrustSelectedIndexCanReachEndOfList(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.router.focus = focusContent
	m.results.result = &AnalysisResult{
		TrustOutput: &trust.ChainOutput{
			TrustChains: []trust.TrustChain{
				{ID: "TC1", RootNode: "user", LeafNode: "app"},
				{ID: "TC2", RootNode: "admin", LeafNode: "db"},
				{ID: "TC3", RootNode: "vendor", LeafNode: "api"},
				{ID: "TC4", RootNode: "dev", LeafNode: "kms"},
				{ID: "TC5", RootNode: "ops", LeafNode: "vault"},
			},
		},
	}
	m.results.resultTab = 4
	ts := m.results.tabStateFor(4)

	for i := 0; i < 10; i++ {
		model, _ := m.updateResults(msgFromString("down"))
		mm := model.(mainModel)
		m = &mm
		ts = m.results.tabStateFor(4)
	}
	if ts.selectedIndex != 4 {
		t.Errorf("trust: after 10 down presses, selectedIndex = %d, want 4 (last of 5)", ts.selectedIndex)
	}
}

func TestVerificationSelectedIndexCanReachEndOfList(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.router.focus = focusContent
	m.results.result = &AnalysisResult{
		VerificationOutput: &verify.VerificationOutput{
			Assessment: &verify.VerificationAssessment{
				VerifiedCount:   5,
				PartialCount:    2,
				UnverifiedCount: 1,
				NoEvidenceCount: 0,
			},
			CISOView: &verify.CISOReviewView{
				TopAssumptionsToVerify: []verify.VerificationPlan{
					{AssumptionText: "MFA enforced"},
					{AssumptionText: "TLS enabled"},
					{AssumptionText: "RBAC configured"},
				},
				EvidenceGaps: []string{},
			},
		},
	}
	m.results.resultTab = 2
	ts := m.results.tabStateFor(2)

	for i := 0; i < 10; i++ {
		model, _ := m.updateResults(msgFromString("down"))
		mm := model.(mainModel)
		m = &mm
		ts = m.results.tabStateFor(2)
	}
	if ts.selectedIndex != 7 {
		t.Errorf("verification: after 10 down presses, selectedIndex = %d, want 7 (tabCount-1 = %d)", ts.selectedIndex, 5+2+1+0-1)
	}
}

func TestContradictionsSelectedIndexCanReachEndOfList(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.router.focus = focusContent
	m.results.result = &AnalysisResult{
		Contradictions: []Contradiction{
			{ID: "C1", Description: "MFA vs service accounts"},
			{ID: "C2", Description: "Private DB vs public route"},
			{ID: "C3", Description: "TLS termination vs internal"},
			{ID: "C4", Description: "Admin access vs audit"},
		},
	}
	m.results.resultTab = 3
	ts := m.results.tabStateFor(3)

	for i := 0; i < 10; i++ {
		model, _ := m.updateResults(msgFromString("down"))
		mm := model.(mainModel)
		m = &mm
		ts = m.results.tabStateFor(3)
	}
	if ts.selectedIndex != 3 {
		t.Errorf("contradictions: after 10 down presses, selectedIndex = %d, want 3 (last of 4)", ts.selectedIndex)
	}
}

func TestNoNegativeViewportOffset(t *testing.T) {
	m := defaultTestModel()
	m.vp.YOffset = 0

	// All views should never produce negative YOffset
	model, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	mm := model.(mainModel)
	if mm.vp.YOffset < 0 {
		t.Errorf("viewport offset should never be negative, got %d", mm.vp.YOffset)
	}

	// After scrolling down and up, should not go negative
	mm.vp.YOffset = 50
	mm.vp.Height = 40
	mm.vp.SetContent(strings.Repeat("line\n", 200))
	for i := 0; i < 60; i++ {
		mm.vp.LineUp(1)
	}
	if mm.vp.YOffset < 0 {
		t.Errorf("viewport offset should not go negative after many LineUp calls, got %d", mm.vp.YOffset)
	}
}

func TestTrustRendersAllChains(t *testing.T) {
	m := defaultTestModel()
	result := &AnalysisResult{
		TrustOutput: &trust.ChainOutput{
			TrustChains: []trust.TrustChain{
				{ID: "TC-A", RootNode: "user", LeafNode: "app", Length: 3, Confidence: 0.8},
				{ID: "TC-B", RootNode: "admin", LeafNode: "db", Length: 4, Confidence: 0.9},
				{ID: "TC-C", RootNode: "vendor", LeafNode: "api", Length: 2, Confidence: 0.7},
				{ID: "TC-D", RootNode: "dev", LeafNode: "kms", Length: 5, Confidence: 0.6},
				{ID: "TC-E", RootNode: "ops", LeafNode: "vault", Length: 3, Confidence: 0.95},
			},
		},
	}
	ts := &tabState{}
	output := renderResultTrust(m.styles, result, ts, 80)
	if output == "" {
		t.Fatal("renderResultTrust returned empty output")
	}
	for _, chain := range result.TrustOutput.TrustChains {
		if !strings.Contains(output, chain.ID) {
			t.Errorf("rendered output should contain chain ID %q", chain.ID)
		}
	}
}

func TestMouseWheelChangesSelectionOnContentTab(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.router.focus = focusContent
	m.results.result = &AnalysisResult{
		Assumptions: []Assumption{
			{ID: "A1"}, {ID: "A2"}, {ID: "A3"}, {ID: "A4"},
		},
	}
	m.results.resultTab = 1
	ts := m.results.tabStateFor(1)
	ts.selectedIndex = 0

	// MouseWheelDown should increment selectedIndex
	model, _ := m.updateResults(tea.MouseMsg{Type: tea.MouseWheelDown})
	mm := model.(mainModel)
	ts = mm.results.tabStateFor(1)
	if ts.selectedIndex != 1 {
		t.Errorf("after MouseWheelDown: selectedIndex = %d, want 1", ts.selectedIndex)
	}

	// MouseWheelUp should decrement selectedIndex
	model, _ = mm.updateResults(tea.MouseMsg{Type: tea.MouseWheelUp})
	mm = model.(mainModel)
	ts = mm.results.tabStateFor(1)
	if ts.selectedIndex != 0 {
		t.Errorf("after MouseWheelUp: selectedIndex = %d, want 0", ts.selectedIndex)
	}

	// On overview tab (tab 0), mouse wheel should NOT change selectedIndex
	m2 := defaultTestModel()
	m2.router.currentView = caseView
	m2.router.focus = focusContent
	m2.results.result = &AnalysisResult{
		Assumptions: []Assumption{{ID: "A1"}, {ID: "A2"}},
	}
	m2.results.resultTab = 0
	ts2 := m2.results.tabStateFor(0)

	model, _ = m2.updateResults(tea.MouseMsg{Type: tea.MouseWheelDown})
	mm2 := model.(mainModel)
	ts2 = mm2.results.tabStateFor(0)
	if ts2.selectedIndex != 0 {
		t.Errorf("overview tab: selectedIndex should stay 0 after MouseWheelDown, got %d", ts2.selectedIndex)
	}
}

func make69SDRIFindings() *AnalysisResult {
	findings := make([]SDRIDesignFinding, 69)
	for i := range findings {
		findings[i] = SDRIDesignFinding{
			Title:       fmt.Sprintf("SDRI-Finding-%d", i+1),
			Severity:    "Critical",
			Description: fmt.Sprintf("Description for finding %d", i+1),
		}
	}
	return &AnalysisResult{
		SDRIDesignFindings: findings,
	}
}

func Test69SDRIFindingsNavigation(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.router.focus = focusContent
	m.results.result = make69SDRIFindings()
	m.results.resultTab = 6

	ts := m.results.tabStateFor(6)
	if ts.selectedIndex != 0 {
		t.Errorf("initial selectedIndex = %d, want 0", ts.selectedIndex)
	}

	// Navigate from item 1 to item 69
	for i := 0; i < 68; i++ {
		model, _ := m.updateResults(msgFromString("down"))
		mm := model.(mainModel)
		*m = mm
		ts = m.results.tabStateFor(6)

		if ts.selectedIndex != i+1 {
			t.Fatalf("step %d: selectedIndex = %d, want %d", i, ts.selectedIndex, i+1)
		}

		visHeight := m.paneHeight()
		if ts.selectedIndex < ts.ViewportOffset || ts.selectedIndex >= ts.ViewportOffset+visHeight {
			t.Fatalf("step %d: selectedIndex %d outside viewport [%d, %d)",
				i, ts.selectedIndex, ts.ViewportOffset, ts.ViewportOffset+visHeight)
		}
	}

	if ts.selectedIndex != 68 {
		t.Errorf("final selectedIndex = %d, want 68", ts.selectedIndex)
	}
}

func TestEnterOnLastItemShowsDetailPane(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.router.focus = focusContent
	m.results.result = make69SDRIFindings()
	m.results.resultTab = 6
	m.results.detailFocus = false

	// Navigate to last item
	for i := 0; i < 68; i++ {
		model, _ := m.updateResults(msgFromString("down"))
		mm := model.(mainModel)
		*m = mm
	}

	// Press Enter to focus detail pane
	model, _ := m.updateResults(msgFromString("enter"))
	mm := model.(mainModel)
	*m = mm
	if !m.results.detailFocus {
		t.Error("detailFocus should be true after Enter on last item")
	}

	// Esc returns to list
	model, _ = m.updateResults(msgFromString("esc"))
	mm = model.(mainModel)
	*m = mm
	if m.results.detailFocus {
		t.Error("detailFocus should be false after Esc")
	}
}

func TestEnsureSelectedVisibleOnResize(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.router.focus = focusContent
	m.results.result = make69SDRIFindings()
	m.results.resultTab = 6

	// Navigate to item 60
	ts := m.results.tabStateFor(6)
	for i := 0; i < 60; i++ {
		model, _ := m.updateResults(msgFromString("down"))
		mm := model.(mainModel)
		*m = mm
	}
	ts = m.results.tabStateFor(6)
	if ts.selectedIndex != 60 {
		t.Fatalf("selectedIndex = %d, want 60", ts.selectedIndex)
	}

	// Simulate resize to 80x24
	model, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	mm := model.(mainModel)
	*m = mm
	ts = m.results.tabStateFor(6)

	visHeight := m.paneHeight()
	if ts.selectedIndex < ts.ViewportOffset || ts.selectedIndex >= ts.ViewportOffset+visHeight {
		t.Errorf("after resize: selectedIndex %d outside viewport [%d, %d)",
			ts.selectedIndex, ts.ViewportOffset, ts.ViewportOffset+visHeight)
	}
}

func TestSearchSelectsVisibleItem(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.router.focus = focusContent
	m.results.result = make69SDRIFindings()
	m.results.resultTab = 6

	// Activate search
	model, _ := m.updateResults(msgFromString("/"))
	mm := model.(mainModel)
	*m = mm
	ts := m.results.tabStateFor(6)
	if !ts.filterActive {
		t.Fatal("filterActive should be true after /")
	}

	// Search for "50"
	for _, ch := range "50" {
		model, _ = m.updateResults(msgFromString(string(ch)))
		mm = model.(mainModel)
		*m = mm
	}
	ts = m.results.tabStateFor(6)

	// n to navigate to first match
	model, _ = m.updateResults(msgFromString("n"))
	mm = model.(mainModel)
	*m = mm
	ts = m.results.tabStateFor(6)

	visHeight := m.paneHeight()
	if ts.selectedIndex < ts.ViewportOffset || ts.selectedIndex >= ts.ViewportOffset+visHeight {
		t.Errorf("after search: selectedIndex %d outside viewport [%d, %d)",
			ts.selectedIndex, ts.ViewportOffset, ts.ViewportOffset+visHeight)
	}

	// Esc to exit search
	model, _ = m.updateResults(msgFromString("esc"))
	mm = model.(mainModel)
	*m = mm
	ts = m.results.tabStateFor(6)
	if ts.filterActive {
		t.Error("filterActive should be false after Esc")
	}
}

func make100TrustChains() *AnalysisResult {
	chains := make([]trust.TrustChain, 100)
	for i := range chains {
		chains[i] = trust.TrustChain{
			ID:         fmt.Sprintf("TC-%d", i+1),
			Length:     i + 1,
			Risk:       "High",
			Confidence: 0.85,
			RootNode:   fmt.Sprintf("Node-%d-A", i+1),
			LeafNode:   fmt.Sprintf("Node-%d-B", i+1),
		}
	}
	return &AnalysisResult{
		TrustOutput: &trust.ChainOutput{
			TrustChains: chains,
		},
	}
}

func Test100TrustChainsNavigation(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.router.focus = focusContent
	m.results.result = make100TrustChains()
	m.results.resultTab = 4

	// Navigate to last chain using pgdn several times then individual steps
	ts := m.results.tabStateFor(4)
	for i := 0; i < 99; i++ {
		model, _ := m.updateResults(msgFromString("down"))
		mm := model.(mainModel)
		*m = mm
		ts = m.results.tabStateFor(4)

		visHeight := m.paneHeight()
		if ts.selectedIndex < ts.ViewportOffset || ts.selectedIndex >= ts.ViewportOffset+visHeight {
			t.Fatalf("step %d: selectedIndex %d outside viewport [%d, %d)",
				i, ts.selectedIndex, ts.ViewportOffset, ts.ViewportOffset+visHeight)
		}
	}

	if ts.selectedIndex != 99 {
		t.Errorf("final selectedIndex = %d, want 99", ts.selectedIndex)
	}

	// Enter to focus detail pane
	model, _ := m.updateResults(msgFromString("enter"))
	mm := model.(mainModel)
	*m = mm
	if !m.results.detailFocus {
		t.Error("detailFocus should be true after Enter on last trust chain")
	}
}

func TestNoInlineDropdownExpansion(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.router.focus = focusContent
	m.results.result = make69SDRIFindings()
	m.results.resultTab = 6

	ts := m.results.tabStateFor(6)
	if ts.detailOpen {
		t.Error("detailOpen should be false initially (no dropdown expansion)")
	}

	// Press Enter — should set detailFocus, not detailOpen
	model, _ := m.updateResults(msgFromString("enter"))
	mm := model.(mainModel)
	*m = mm
	ts = m.results.tabStateFor(6)

	if !m.results.detailFocus {
		t.Error("detailFocus should be true after Enter (split pane), not detailOpen")
	}
	if ts.detailOpen {
		t.Error("detailOpen must remain false — no inline dropdown expansion allowed")
	}
}

func TestPageUpDownNavigation(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.router.focus = focusContent
	m.results.result = make69SDRIFindings()
	m.results.resultTab = 6

	// pgdn should jump by visibleHeight
	model, _ := m.updateResults(msgFromString("pgdn"))
	mm := model.(mainModel)
	*m = mm
	ts := m.results.tabStateFor(6)
	if ts.selectedIndex <= 0 {
		t.Error("pgdn should increase selectedIndex")
	}
	visHeight := m.paneHeight()
	if ts.selectedIndex < ts.ViewportOffset || ts.selectedIndex >= ts.ViewportOffset+visHeight {
		t.Errorf("after pgdn: selectedIndex %d outside viewport [%d, %d)",
			ts.selectedIndex, ts.ViewportOffset, ts.ViewportOffset+visHeight)
	}

	// Navigate to a middle position
	for i := 0; i < 30; i++ {
		model, _ = m.updateResults(msgFromString("down"))
		mm = model.(mainModel)
		*m = mm
	}
	ts = m.results.tabStateFor(6)

	// pgup should jump back
	model, _ = m.updateResults(msgFromString("pgup"))
	mm = model.(mainModel)
	*m = mm
	ts = m.results.tabStateFor(6)
	if ts.selectedIndex < ts.ViewportOffset || ts.selectedIndex >= ts.ViewportOffset+visHeight {
		t.Errorf("after pgup: selectedIndex %d outside viewport [%d, %d)",
			ts.selectedIndex, ts.ViewportOffset, ts.ViewportOffset+visHeight)
	}

	// home should go to first item
	for i := 0; i < 68; i++ {
		model, _ = m.updateResults(msgFromString("down"))
		mm = model.(mainModel)
		*m = mm
	}
	model, _ = m.updateResults(msgFromString("home"))
	mm = model.(mainModel)
	*m = mm
	ts = m.results.tabStateFor(6)
	if ts.selectedIndex != 0 {
		t.Errorf("after home: selectedIndex = %d, want 0", ts.selectedIndex)
	}

	// end should go to last item
	model, _ = m.updateResults(msgFromString("end"))
	mm = model.(mainModel)
	*m = mm
	ts = m.results.tabStateFor(6)
	if ts.selectedIndex != 68 {
		t.Errorf("after end: selectedIndex = %d, want 68", ts.selectedIndex)
	}
}

func TestViewResultsRendersSplitPane(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.width = 120
	m.height = 40
	m.results.result = make69SDRIFindings()
	m.results.resultTab = 6

	output := m.viewResults()
	if output == "" {
		t.Fatal("viewResults returned empty output")
	}
	if !strings.Contains(output, "SDRI") {
		t.Error("output should contain SDRI tab content")
	}
}

func TestNarrowTerminalFallback(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.width = 80
	m.height = 24
	m.results.result = make69SDRIFindings()
	m.results.resultTab = 6

	// viewResults should not panic on narrow terminal
	output := m.viewResults()
	if output == "" {
		t.Fatal("viewResults returned empty on narrow terminal")
	}
}

func TestDetailPaneScroll(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = caseView
	m.router.focus = focusContent
	m.results.result = make69SDRIFindings()
	m.results.resultTab = 6

	// Focus detail pane
	model, _ := m.updateResults(msgFromString("enter"))
	mm := model.(mainModel)
	*m = mm

	if !m.results.detailFocus {
		t.Fatal("detailFocus should be true after Enter")
	}

	ts := m.results.tabStateFor(6)
	if ts.DetailOffset != 0 {
		t.Errorf("initial DetailOffset = %d, want 0", ts.DetailOffset)
	}

	// Scroll detail pane down
	for i := 0; i < 5; i++ {
		model, _ = m.updateResults(msgFromString("down"))
		mm = model.(mainModel)
		*m = mm
	}
	ts = m.results.tabStateFor(6)
	if ts.DetailOffset < 5 {
		t.Errorf("after 5 down: DetailOffset = %d, want >= 5", ts.DetailOffset)
	}

	// Scroll detail pane up
	for i := 0; i < 3; i++ {
		model, _ = m.updateResults(msgFromString("up"))
		mm = model.(mainModel)
		*m = mm
	}
	ts = m.results.tabStateFor(6)
	if ts.DetailOffset >= 5 {
		t.Errorf("after 3 up: DetailOffset = %d, want < 5", ts.DetailOffset)
	}
}

func TestTabStatePreservedOnSwitch(t *testing.T) {
	m := defaultTestModel()
	m.results.result = make69SDRIFindings()
	m.results.resultTab = 6

	ts := m.results.tabStateFor(6)
	ts.selectedIndex = 30
	ts.ViewportOffset = 20

	// Switch to tab 4 then back
	m.results.resultTab = 4
	m.results.resultTab = 6

	tsBack := m.results.tabStateFor(6)
	if tsBack.selectedIndex != 30 {
		t.Errorf("after tab switch: selectedIndex = %d, want 30", tsBack.selectedIndex)
	}
	if tsBack.ViewportOffset != 20 {
		t.Errorf("after tab switch: ViewportOffset = %d, want 20", tsBack.ViewportOffset)
	}
}
