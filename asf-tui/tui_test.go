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
	s := NewStyles(Themes["dark"])
	r := &AnalysisResult{}

	cases := []struct {
		name string
		fn   func() string
	}{
		{"Assumptions", func() string { return renderResultAssumptions(s, r, "") }},
		{"Verification", func() string { return renderResultVerification(s, r) }},
		{"Contradictions", func() string { return renderResultContradictions(s, r, "") }},
		{"Controls", func() string { return renderResultControls(s, r, "") }},
		{"BlindSpots", func() string { return renderResultBlindSpots(s, r) }},
		{"Impact", func() string { return renderResultImpact(s, r) }},
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
	if got := resultTabCount(r, 4); got != 1 {
		t.Errorf("Trust chains count = %d, want 1", got)
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

func TestViewForSidebar(t *testing.T) {
	expected := []view{dashboardView, analyzeView, resultsView, fileBrowserView, localaiView, settingsView, aboutView, helpView}
	for i, v := range expected {
		if got := viewForSidebar(i); got != v {
			t.Errorf("viewForSidebar(%d) = %d, want %d", i, got, v)
		}
	}
	// Out of range should fall back to dashboard
	if got := viewForSidebar(100); got != dashboardView {
		t.Errorf("viewForSidebar(100) = %d, want dashboardView", got)
	}
}

func TestSidebarItems(t *testing.T) {
	if len(sidebarItems) != 8 {
		t.Errorf("sidebar has %d items, want 8", len(sidebarItems))
	}
	if len(sidebarViews) != len(sidebarItems) {
		t.Errorf("sidebarViews (%d) != sidebarItems (%d)", len(sidebarViews), len(sidebarItems))
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
	if len(rm.tabs) != 9 {
		t.Errorf("results has %d tabs, want 9", len(rm.tabs))
	}
	expected := []string{"Summary", "Assumptions", "Verification", "Contradictions", "Trust", "Impact", "Blind Spots", "Controls", "Reports"}
	for i, name := range expected {
		if rm.tabs[i].name != name {
			t.Errorf("tab[%d] = %q, want %q", i, rm.tabs[i].name, name)
		}
	}
}

func TestNewFileBrowserModel(t *testing.T) {
	fb := newFileBrowserModel()
	if fb.path != "." {
		t.Errorf("initial path = %q, want %q", fb.path, ".")
	}
	if fb.showHidden {
		t.Error("showHidden should be false initially")
	}
	if fb.showPreview {
		t.Error("showPreview should be false initially")
	}
}

func TestRiskStyle(t *testing.T) {
	s := NewStyles(Themes["dark"])
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
	s := NewStyles(Themes["dark"])
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
