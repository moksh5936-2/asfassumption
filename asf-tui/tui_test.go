package main

import (
	"testing"

	"asf-tui/asf/trust"
	"asf-tui/asf/verify"
)

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		bytes int64
		want  string
	}{
		{0, "0 B"},
		{1, "1 B"},
		{1000, "1000 B"},
		{1024, "1024 B"},
		{2048, "2.0 KB"},
		{1048576, "1024.0 KB"},
		{1073741824, "1024.0 MB"},
		{1 << 30, "1024.0 MB"},
		{(1 << 30) + 1, "1.0 GB"},
	}
	for _, tt := range tests {
		got := formatFileSize(tt.bytes)
		if got != tt.want {
			t.Errorf("formatFileSize(%d) = %q, want %q", tt.bytes, got, tt.want)
		}
	}
}

func TestPadRight(t *testing.T) {
	tests := []struct {
		s    string
		n    int
		want string
	}{
		{"hello", 8, "hello   "},
		{"hi", 2, "hi"},
		{"test", 3, "test"},
	}
	for _, tt := range tests {
		got := padRight(tt.s, tt.n)
		if got != tt.want {
			t.Errorf("padRight(%q, %d) = %q, want %q", tt.s, tt.n, got, tt.want)
		}
	}
}

func TestCountRisk(t *testing.T) {
	assumptions := []Assumption{
		{ID: "A1", Risk: "Critical"},
		{ID: "A2", Risk: "High"},
		{ID: "A3", Risk: "Medium"},
		{ID: "A4", Risk: "Critical"},
		{ID: "A5", Risk: "Low"},
	}
	if got := countRisk(assumptions, "Critical"); got != 2 {
		t.Errorf("countRisk(Critical) = %d, want 2", got)
	}
	if got := countRisk(assumptions, "High"); got != 1 {
		t.Errorf("countRisk(High) = %d, want 1", got)
	}
	if got := countRisk(assumptions, "Low"); got != 1 {
		t.Errorf("countRisk(Low) = %d, want 1", got)
	}
	if got := countRisk(assumptions, "None"); got != 0 {
		t.Errorf("countRisk(None) = %d, want 0", got)
	}
}

func TestEmptyResultRendersEmptyStates(t *testing.T) {
	s := NewStyles(Themes["ASF0"])
	r := &AnalysisResult{}

	cases := []struct {
		name string
		fn   func() string
	}{
		{"Assumptions", func() string { return renderResultAssumptions(s, r, "", 80) }},
		{"Verification", func() string { return renderResultVerification(s, r, 80) }},
		{"Contradictions", func() string { return renderResultContradictions(s, r, "", 80) }},
		{"Controls", func() string { return renderResultControls(s, r, "", 80) }},
		{"BlindSpots", func() string { return renderResultBlindSpots(s, r, 80) }},
		{"Impact", func() string { return renderResultImpact(s, r, 80) }},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.fn()
			if got == "" {
				t.Errorf("empty state for %q should not be empty", c.name)
			}
		})
	}
}

func TestResultTabCount(t *testing.T) {
	r := &AnalysisResult{
		Assumptions: []Assumption{{ID: "A1"}, {ID: "A2"}},
		Contradictions: []Contradiction{
			{ID: "C1", Description: "test"},
		},
		TrustOutput: &trust.ChainOutput{
			TrustChains: []trust.TrustChain{
				{ID: "TC1", Confidence: 0.8},
			},
			SinglePointsOfTrust: []trust.SinglePointOfTrustFailure{
				{AssumptionText: "SPOF1"},
			},
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

	if got := resultTabCount(r, 1); got != 2 {
		t.Errorf("Assumptions count = %d, want 2", got)
	}
	if got := resultTabCount(r, 3); got != 1 {
		t.Errorf("Contradictions count = %d, want 1", got)
	}
	if got := resultTabCount(r, 4); got != 2 {
		t.Errorf("Trust count (chains + SPOFs) = %d, want 2", got)
	}
}

func TestSupportedExts(t *testing.T) {
	exts := []string{".yaml", ".yml", ".json", ".md", ".mmd", ".drawio", ".svg", ".pdf", ".docx", ".txt"}
	for _, ext := range exts {
		if !supportedExts[ext] {
			t.Errorf("expected %q to be supported", ext)
		}
	}
	if supportedExts[".exe"] {
		t.Error("expected .exe to be unsupported")
	}
	if supportedExts[".go"] {
		t.Error("expected .go to be unsupported")
	}
}

func TestAddRecentFile(t *testing.T) {
	m := &mainModel{}

	m.addRecentFile("a.yaml")
	if len(m.recentFiles) != 1 || m.recentFiles[0] != "a.yaml" {
		t.Errorf("after first add: %v", m.recentFiles)
	}

	m.addRecentFile("b.yaml")
	if len(m.recentFiles) != 2 || m.recentFiles[0] != "b.yaml" {
		t.Errorf("after second add: %v", m.recentFiles)
	}

	m.addRecentFile("a.yaml")
	if len(m.recentFiles) != 2 || m.recentFiles[0] != "a.yaml" {
		t.Errorf("after dedup add: %v", m.recentFiles)
	}

	m.addRecentFile("")
	if len(m.recentFiles) != 2 {
		t.Errorf("empty path should not be added: %v", m.recentFiles)
	}
}

func TestSidebarTree(t *testing.T) {
	r := newRouter()
	nodes := r.sidebarVisibleNodes()
	// 3 sections + 5 items = 8 (no case entries initially)
	expected := 8
	if len(nodes) != expected {
		t.Errorf("sidebar has %d visible nodes, want %d", len(nodes), expected)
	}
	analyzeFound := false
	reviewFound := false
	settingsFound := false
	for _, n := range r.sidebarVisibleNodes() {
		if n.vid == analyzeView && !n.isSection {
			analyzeFound = true
		}
		if n.vid == reviewView && !n.isSection {
			reviewFound = true
		}
		if n.vid == settingsView && !n.isSection {
			settingsFound = true
		}
	}
	if !analyzeFound {
		t.Error("sidebar should contain analyzeView (+ New Analysis)")
	}
	if reviewFound {
		t.Error("sidebar should NOT contain reviewView (removed from sidebar)")
	}
	if !settingsFound {
		t.Error("sidebar should contain settingsView")
	}
}

type testModel struct {
	vpYOffset int
	vpHeight  int
	vpTotal   int
}

func (m *testModel) scrollPercent() string {
	total := m.vpTotal
	visible := m.vpHeight
	offset := m.vpYOffset
	if total <= visible || total == 0 {
		return "All"
	}
	pct := int(float64(offset+visible) / float64(total) * 100)
	if pct > 100 {
		pct = 100
	}
	to := offset + visible
	if to > total {
		to = total
	}
	return ""
}

func TestScrollPercentLogic(t *testing.T) {
	// Test directly using the same logic as viewportScrollPercent
	tests := []struct {
		offset  int
		visible int
		total   int
		all     bool
	}{
		{0, 40, 40, true},
		{0, 40, 30, true},
		{0, 40, 0, true},
		{0, 40, 200, false},
		{160, 40, 200, false},
	}
	for _, tt := range tests {
		all := tt.total <= tt.visible || tt.total == 0
		if all != tt.all {
			t.Errorf("scroll all check (offset=%d, visible=%d, total=%d) = %v, want %v",
				tt.offset, tt.visible, tt.total, all, tt.all)
		}
	}
}

func TestNewResultsModel(t *testing.T) {
	rm := newResultsModel()
	if len(rm.tabs) != 7 {
		t.Errorf("results has %d tabs, want 7", len(rm.tabs))
	}
	expected := []string{"Overview", "Assumptions", "Verification", "Contradictions", "Trust", "Controls", "SDRI"}
	for i, name := range expected {
		if rm.tabs[i].name != name {
			t.Errorf("tab[%d] = %q, want %q", i, rm.tabs[i].name, name)
		}
	}
}

func TestNewFilePickerState(t *testing.T) {
	fp := newFilePickerState()
	if fp.path != "." {
		t.Errorf("initial path = %q, want %q", fp.path, ".")
	}
	if fp.showHidden {
		t.Error("showHidden should be false initially")
	}
	if fp.showPreview {
		t.Error("showPreview should be false initially")
	}
	if fp.mode != pickerArchitecture {
		t.Errorf("initial mode = %d, want pickerArchitecture (%d)", fp.mode, pickerArchitecture)
	}
}

func TestFilePickerMode(t *testing.T) {
	fp := newFilePickerState()
	fp.mode = pickerEvidence
	if fp.mode != pickerEvidence {
		t.Errorf("mode = %d, want pickerEvidence (%d)", fp.mode, pickerEvidence)
	}
}

func TestRiskStyle(t *testing.T) {
	s := NewStyles(Themes["ASF0"])
	// Verify styles render without panic and produce non-empty output
	cases := []struct {
		risk RiskLevel
	}{
		{"Critical"},
		{"High"},
		{"Medium"},
		{"Low"},
		{"Unknown"},
	}
	for _, c := range cases {
		st := riskStyle(s, c.risk)
		rendered := st.Render("test")
		if rendered == "" {
			t.Errorf("riskStyle(%q) rendered empty", c.risk)
		}
	}
}

func TestConfidenceStyle(t *testing.T) {
	s := NewStyles(Themes["ASF0"])
	cases := []struct {
		pct int
	}{
		{90},
		{60},
		{30},
	}
	for _, c := range cases {
		st := confidenceStyle(s, c.pct)
		rendered := st.Render("test")
		if rendered == "" {
			t.Errorf("confidenceStyle(%d) rendered empty", c.pct)
		}
	}
}

func TestNewLocalAIModel(t *testing.T) {
	cfg := DefaultConfig()
	cfg.AI.Enabled = true
	cfg.AI.ActiveModel = "llama3.2:3b"
	cfg.AI.InstalledModels = []string{"llama3.2:3b", "mistral:7b"}
	lm := newLocalAIModel(&cfg)

	if lm.activeModel != "llama3.2:3b" {
		t.Errorf("activeModel = %q, want %q", lm.activeModel, "llama3.2:3b")
	}
	if len(lm.catalog) != len(SupportedModels) {
		t.Errorf("catalog has %d entries, want %d", len(lm.catalog), len(SupportedModels))
	}

	found := false
	for _, c := range lm.catalog {
		if c.Info.Name == "llama3.2:3b" {
			found = true
			if !c.Installed {
				t.Error("llama3.2:3b should be marked installed")
			}
			if !c.Active {
				t.Error("llama3.2:3b should be marked active")
			}
		}
	}
	if !found {
		t.Error("llama3.2:3b not found in catalog")
	}
}

func TestLocalAIViewRender(t *testing.T) {
	m := defaultTestModel()
	m.router.currentView = localAIView
	content := m.renderContent()
	if content == "" {
		t.Error("Local AI view should render non-empty content")
	}
}

func TestLocalAIViewSwitch(t *testing.T) {
	m := defaultTestModel()
	m.router.SetView(localAIView)
	if m.router.currentView != localAIView {
		t.Errorf("currentView = %v, want localAIView (%v)", m.router.currentView, localAIView)
	}

	// Verify scroll state is tracked for localAIView
	m.scrollY[localAIView] = 42
	m.restoreScroll()
	if m.vp.YOffset != 42 {
		t.Errorf("scroll y offset = %d, want 42", m.vp.YOffset)
	}
}

func TestLocalAISidebarEntry(t *testing.T) {
	m := defaultTestModel()
	found := false
	for _, n := range m.router.sidebarVisibleNodes() {
		if n.vid == localAIView && !n.isSection {
			found = true
			break
		}
	}
	if !found {
		t.Error("Local AI sidebar entry not found in sidebar nodes")
	}
}

func TestLocalAIAnalysisMode(t *testing.T) {
	a := newAnalyzeModel(nil)
	modeFound := false
	for _, item := range a.items {
		if item.value == ModeASFAndAI {
			modeFound = true
			break
		}
	}
	if !modeFound {
		t.Error("ASF Engine + Local AI mode not found in analyze menu")
	}
}

func TestLocalAISidebarNavigation(t *testing.T) {
	m := defaultTestModel()
	m.router.focus = focusSidebar

	// Find Local AI entry and navigate to it
	var localAIIndex int
	found := false
	for i, n := range m.router.sidebarVisibleNodes() {
		if n.vid == localAIView && !n.isSection {
			localAIIndex = i
			found = true
			break
		}
	}
	if !found {
		t.Skip("Local AI sidebar entry not found")
	}

	m.router.sidebarSel = localAIIndex
	to, tab := m.router.sidebarActivate()
	if to != localAIView {
		t.Errorf("sidebarActivate should navigate to localAIView, got %v", to)
	}
	if tab != -1 {
		t.Errorf("Local AI tab should be -1, got %d", tab)
	}
}

func TestLocalAICasesWorkNavigation(t *testing.T) {
	// Test that CASES / AI / SYSTEM navigation still works with Local AI added
	m := defaultTestModel()
	m.router.focus = focusSidebar

	viewsToCheck := []view{analyzeView, localAIView, settingsView, helpView, aboutView}
	for _, target := range viewsToCheck {
		for i, n := range m.router.sidebarVisibleNodes() {
			if n.vid == target && !n.isSection {
				m.router.sidebarSel = i
				to, _ := m.router.sidebarActivate()
				if to != target {
					t.Errorf("sidebarActivate should navigate to %v, got %v", target, to)
				}
				break
			}
		}
	}
}

func TestLocalAIRouteDoesNotConflict(t *testing.T) {
	// Verify localAIView is a distinct value
	if localAIView == analyzeView || localAIView == caseView || localAIView == settingsView {
		t.Error("localAIView conflicts with an existing view")
	}
	if localAIView == helpView || localAIView == aboutView || localAIView == reviewView {
		t.Error("localAIView conflicts with an existing view")
	}
	if localAIView == validationView || localAIView == reportsView {
		t.Error("localAIView conflicts with an existing view")
	}
}

func TestAnalyzeStage(t *testing.T) {
	if s := analyzeStage(0); s == "" {
		t.Error("analyzeStage(0) should not be empty")
	}
	if s := analyzeStage(50); s == "" {
		t.Error("analyzeStage(50) should not be empty")
	}
	if s := analyzeStage(100); s != "Generating Gap Analysis..." {
		t.Errorf("analyzeStage(100) = %q, want 'Generating Gap Analysis...'", s)
	}
}
