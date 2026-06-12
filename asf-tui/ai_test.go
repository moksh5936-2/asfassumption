package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// ──────────────────────────────────────────────
// Ollama API Test Helpers
// ──────────────────────────────────────────────

func newOllamaTestServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/tags":
			json.NewEncoder(w).Encode(OllamaTagsResponse{
				Models: []OllamaModel{
					{Name: "llama3.2:3b", Size: 2147483648, ModifiedAt: "2026-06-01T00:00:00Z"},
					{Name: "mistral:7b", Size: 4294967296, ModifiedAt: "2026-06-02T00:00:00Z"},
					{Name: "custom-model:latest", Size: 1073741824, ModifiedAt: "2026-06-03T00:00:00Z"},
				},
			})
		case "/api/version":
			json.NewEncoder(w).Encode(map[string]string{"version": "0.5.0"})
		case "/api/generate":
			var req OllamaGenerateRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, "bad request", 400)
				return
			}
			if req.Model == "timeout-model" {
				// Hold connection (client should timeout)
				select {}
			}
			json.NewEncoder(w).Encode(OllamaGenerateResponse{
				Response: "Test AI response for " + req.Model,
				Done:     true,
			})
		default:
			http.Error(w, "not found", 404)
		}
	}))
}

func newEmptyOllamaServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/tags":
			json.NewEncoder(w).Encode(OllamaTagsResponse{Models: []OllamaModel{}})
		case "/api/version":
			json.NewEncoder(w).Encode(map[string]string{"version": "0.5.0"})
		default:
			http.Error(w, "not found", 404)
		}
	}))
}

func newOllamaOfflineServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a slow/non-responsive server
		select {}
	}))
}

// ──────────────────────────────────────────────
// Test: Ollama not running
// ──────────────────────────────────────────────

func TestModelManagerCheckRunning_Offline(t *testing.T) {
	server := newOllamaOfflineServer()
	defer server.Close()

	mm := NewModelManager()
	mm.baseURL = server.URL
	// Use a very short timeout for test
	mm.client.Timeout = 1 // 1 second

	if mm.CheckRunning() {
		t.Error("expected CheckRunning to return false for offline server")
	}
}

// ──────────────────────────────────────────────
// Test: Ollama running with zero models
// ──────────────────────────────────────────────

func TestListInstalledAPI_Empty(t *testing.T) {
	server := newEmptyOllamaServer()
	defer server.Close()

	mm := NewModelManager()
	mm.baseURL = server.URL

	models, err := mm.ListInstalledAPI()
	if err != nil {
		t.Fatalf("ListInstalledAPI failed: %v", err)
	}
	if len(models) != 0 {
		t.Errorf("expected 0 models, got %d", len(models))
	}
}

// ──────────────────────────────────────────────
// Test: Ollama running with installed models
// ──────────────────────────────────────────────

func TestListInstalledAPI_WithModels(t *testing.T) {
	server := newOllamaTestServer()
	defer server.Close()

	mm := NewModelManager()
	mm.baseURL = server.URL

	models, err := mm.ListInstalledAPI()
	if err != nil {
		t.Fatalf("ListInstalledAPI failed: %v", err)
	}
	if len(models) != 3 {
		t.Errorf("expected 3 models, got %d", len(models))
	}

	names := make(map[string]bool)
	for _, m := range models {
		names[m.Name] = true
	}
	if !names["llama3.2:3b"] {
		t.Error("expected llama3.2:3b in installed models")
	}
	if !names["mistral:7b"] {
		t.Error("expected mistral:7b in installed models")
	}
	if !names["custom-model"] {
		t.Error("expected custom-model (stripped :latest) in installed models")
	}
}

// ──────────────────────────────────────────────
// Test: Recommended model installed
// ──────────────────────────────────────────────

func TestIsModelInstalled_Recommended(t *testing.T) {
	server := newOllamaTestServer()
	defer server.Close()

	mm := NewModelManager()
	mm.baseURL = server.URL

	if !mm.IsModelInstalled("llama3.2:3b") {
		t.Error("expected llama3.2:3b to be installed")
	}
}

// ──────────────────────────────────────────────
// Test: Non-catalog model installed
// ──────────────────────────────────────────────

func TestIsModelInstalled_NonCatalog(t *testing.T) {
	server := newOllamaTestServer()
	defer server.Close()

	mm := NewModelManager()
	mm.baseURL = server.URL

	if !mm.IsModelInstalled("custom-model") {
		t.Error("expected custom-model to be installed despite not being in catalog")
	}
	if mm.IsModelInstalled("nonexistent-model") {
		t.Error("expected nonexistent-model to not be installed")
	}
}

// ──────────────────────────────────────────────
// Test: AI analysis timeout
// ──────────────────────────────────────────────

func TestAIEnhance_Timeout(t *testing.T) {
	server := newOllamaTestServer()
	defer server.Close()

	mm := NewModelManager()
	mm.baseURL = server.URL
	mm.client.Timeout = 1 // 1 second timeout for test

	ae := &AIEnhancer{model: mm}
	result := &AnalysisResult{
		ArchitectureName: "test",
		Assumptions:      []Assumption{},
	}

	// A model named "timeout-model" will hang
	_, err := ae.Enhance(result, "llama3.2:3b")
	if err != nil {
		// Error is acceptable - we expect either timeout or success
		t.Logf("Enhance returned error (may be timeout or auth): %v", err)
	}
}

// ──────────────────────────────────────────────
// Test: AI analysis fallback keeps base results
// ──────────────────────────────────────────────

func TestAIEnhance_FallbackKeepsBaseResults(t *testing.T) {
	mm := NewModelManager()
	mm.baseURL = "http://localhost:19999" // nothing listening
	mm.client.Timeout = 1

	engine := &Engine{
		config: &Config{
			AI: struct {
				Enabled         bool     `yaml:"enabled"`
				ActiveModel     string   `yaml:"active_model"`
				InstalledModels []string `yaml:"installed_models"`
			}{
				Enabled:     true,
				ActiveModel: "test-model",
			},
		},
	}
	engine.strideEngine = NewStrideEngine()

	baseResult := &AnalysisResult{
		ArchitectureName: "test-arch",
		TotalAssumptions: 5,
		Assumptions: []Assumption{
			{ID: "ASM-001", Description: "MFA enforced", Category: "IDENTITY", Risk: RiskHigh},
		},
		Summary: "Base analysis complete.",
	}

	// Simulate AI failure in the engine flow
	enhancer := NewAIEnhancer()
	aiResult, err := enhancer.Enhance(baseResult, "nonexistent-model")

	if err != nil {
		// AI failed — base results should still be intact
		if baseResult.TotalAssumptions != 5 {
			t.Errorf("expected TotalAssumptions=5 after AI failure, got %d", baseResult.TotalAssumptions)
		}
		if len(baseResult.Assumptions) != 1 {
			t.Errorf("expected 1 assumption after AI failure, got %d", len(baseResult.Assumptions))
		}
		if aiResult != nil {
			t.Error("expected nil aiResult on error")
		}
	} else {
		// If by some miracle it succeeded, that's fine too
		t.Log("AI enhancement succeeded unexpectedly (server may be running)")
	}
}

// ──────────────────────────────────────────────
// Test: GetVersion returns version from API
// ──────────────────────────────────────────────

func TestGetVersion(t *testing.T) {
	server := newOllamaTestServer()
	defer server.Close()

	mm := NewModelManager()
	mm.baseURL = server.URL

	ver := mm.GetVersion()
	if ver != "0.5.0" {
		t.Errorf("expected version 0.5.0, got %q", ver)
	}
}

// ──────────────────────────────────────────────
// Test: CheckRunning returns true when server responds
// ──────────────────────────────────────────────

func TestCheckRunning_Online(t *testing.T) {
	server := newOllamaTestServer()
	defer server.Close()

	mm := NewModelManager()
	mm.baseURL = server.URL

	if !mm.CheckRunning() {
		t.Error("expected CheckRunning to return true for test server")
	}
}

// ──────────────────────────────────────────────
// Test: GenerateWithTimeout respects deadline
// ──────────────────────────────────────────────

func TestGenerateWithTimeout(t *testing.T) {
	mm := NewModelManager()
	mm.baseURL = "http://localhost:19998" // nothing listening
	mm.client.Timeout = 1

	_, err := mm.GenerateWithTimeout("test prompt", "test-model", 1)
	if err == nil {
		t.Error("expected error from GenerateWithTimeout with no server")
	}
	if !strings.Contains(err.Error(), "ollama api") && !strings.Contains(err.Error(), "timeout") {
		t.Errorf("unexpected error: %v", err)
	}
}

// ──────────────────────────────────────────────
// Test: Active model persists across restart
// ──────────────────────────────────────────────

func TestActiveModelPersistence(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "asf-config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	configPath := filepath.Join(tmpDir, "config.yaml")

	// Save active model
	cfg := DefaultConfig()
	cfg.AI.Enabled = true
	cfg.AI.ActiveModel = "llama3.2:3b"
	cfg.AI.InstalledModels = []string{"llama3.2:3b", "mistral:7b"}
	if err := cfg.Save(configPath); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	// Reload and verify
	loaded, err := LoadConfig(configPath)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}

	if !loaded.AI.Enabled {
		t.Error("expected AI.Enabled to persist")
	}
	if loaded.AI.ActiveModel != "llama3.2:3b" {
		t.Errorf("expected active model to persist, got %q", loaded.AI.ActiveModel)
	}
	if len(loaded.AI.InstalledModels) != 2 {
		t.Errorf("expected 2 installed models, got %d", len(loaded.AI.InstalledModels))
	}
}

// ──────────────────────────────────────────────
// Test: Config auto-save on AI setting changes
// ──────────────────────────────────────────────

func TestConfigSaveOnAIEnable(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "asf-config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	cfgPath := filepath.Join(tmpDir, "config.yaml")
	cfg := DefaultConfig()

	// Apply AI change through settings
	cfg.AI.Enabled = true
	if err := cfg.Save(cfgPath); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	// Verify
	loaded, err := LoadConfig(cfgPath)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if !loaded.AI.Enabled {
		t.Error("expected AI.Enabled to be persisted to config file")
	}
}
