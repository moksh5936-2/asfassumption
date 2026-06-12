package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"
)

const (
	ollamaDefaultURL = "http://localhost:11434"
	ollamaTimeout    = 120
	apiTimeout       = 30
)

type ModelManager struct {
	ollamaCmd string
	baseURL   string
	client    *http.Client
}

var SupportedModels = []ModelInfo{
	{Name: "llama3.2:3b", Display: "Llama 3.2 (3B)", Size: "2.0 GB"},
	{Name: "qwen3:4b", Display: "Qwen 3 (4B)", Size: "2.5 GB"},
	{Name: "mistral:7b", Display: "Mistral (7B)", Size: "4.1 GB"},
	{Name: "phi4:14b", Display: "Phi-4 (14B)", Size: "9.1 GB"},
}

type ModelInfo struct {
	Name    string
	Display string
	Size    string
}

type OllamaModel struct {
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	ModifiedAt string `json:"modified_at"`
	Digest     string `json:"digest"`
}

type OllamaTagsResponse struct {
	Models []OllamaModel `json:"models"`
}

func NewModelManager() *ModelManager {
	cmd := findOllama()
	return &ModelManager{
		ollamaCmd: cmd,
		baseURL:   ollamaDefaultURL,
		client: &http.Client{
			Timeout: apiTimeout * time.Second,
		},
	}
}

func findOllama() string {
	for _, p := range []string{"/usr/local/bin/ollama", "/opt/homebrew/bin/ollama", "ollama"} {
		if _, err := exec.LookPath(p); err == nil {
			return p
		}
	}
	return "ollama"
}

// CheckAvailable checks if the ollama binary exists.
func (mm *ModelManager) CheckAvailable() bool {
	return exec.Command(mm.ollamaCmd, "--version").Run() == nil
}

// CheckRunning checks if the Ollama server is responding via the API.
func (mm *ModelManager) CheckRunning() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", mm.baseURL+"/api/tags", nil)
	if err != nil {
		return false
	}
	resp, err := mm.client.Do(req)
	if err != nil {
		return false
	}
	resp.Body.Close()
	return resp.StatusCode == http.StatusOK
}

// GetVersion returns the Ollama server version from /api/version.
func (mm *ModelManager) GetVersion() string {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", mm.baseURL+"/api/version", nil)
	if err != nil {
		return ""
	}
	resp, err := mm.client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	var v struct {
		Version string `json:"version"`
	}
	if json.NewDecoder(resp.Body).Decode(&v) == nil {
		return v.Version
	}
	return ""
}

// ListInstalledAPI queries /api/tags for installed models.
func (mm *ModelManager) ListInstalledAPI() ([]OllamaModel, error) {
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", mm.baseURL+"/api/tags", nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	resp, err := mm.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ollama api /api/tags: %w", err)
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	var result OllamaTagsResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}
	// Normalize names: strip :latest suffix
	for i := range result.Models {
		result.Models[i].Name = strings.TrimSuffix(result.Models[i].Name, ":latest")
	}
	return result.Models, nil
}

// IsModelInstalled checks if a specific model exists in Ollama via /api/tags.
func (mm *ModelManager) IsModelInstalled(modelName string) bool {
	models, err := mm.ListInstalledAPI()
	if err != nil {
		return false
	}
	name := strings.TrimSuffix(modelName, ":latest")
	for _, m := range models {
		if m.Name == name {
			return true
		}
	}
	return false
}

// ListInstalledNames returns just the names of installed models.
func (mm *ModelManager) ListInstalledNames() ([]string, error) {
	models, err := mm.ListInstalledAPI()
	if err != nil {
		return nil, err
	}
	names := make([]string, len(models))
	for i, m := range models {
		names[i] = m.Name
	}
	return names, nil
}

// ListInstalled uses CLI ollama list (fallback).
func (mm *ModelManager) ListInstalled() ([]string, error) {
	var stdout bytes.Buffer
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, mm.ollamaCmd, "list")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("ollama list: %w", err)
	}
	var models []string
	scanner := bufio.NewScanner(&stdout)
	first := true
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if first {
			first = false
			continue
		}
		parts := strings.Fields(line)
		if len(parts) > 0 {
			name := parts[0]
			name = strings.TrimSuffix(name, ":latest")
			models = append(models, name)
		}
	}
	return models, nil
}

// syncProgress tracks download progress in a thread-safe manner.
type syncProgress struct {
	mu      sync.Mutex
	Percent float64
	Stage   string
	Done    bool
	Err     string
}

func (sp *syncProgress) Set(percent float64, stage string) {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.Percent = percent
	sp.Stage = stage
}

func (sp *syncProgress) Finish() {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.Percent = 100
	sp.Stage = "Complete"
	sp.Done = true
}

func (sp *syncProgress) Fail(err string) {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	sp.Err = err
	sp.Done = true
}

func (sp *syncProgress) Snapshot() (float64, string, bool, string) {
	sp.mu.Lock()
	defer sp.mu.Unlock()
	return sp.Percent, sp.Stage, sp.Done, sp.Err
}

// StartDownload pulls a model via `ollama pull <model>`.
func (mm *ModelManager) StartDownload(modelName string, sp *syncProgress) {
	if !mm.CheckAvailable() {
		sp.Fail("Ollama not found. Install from https://ollama.ai")
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), ollamaTimeout*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, mm.ollamaCmd, "pull", modelName)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		sp.Fail(fmt.Sprintf("pipe: %v", err))
		return
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		sp.Fail(fmt.Sprintf("start: %v", err))
		return
	}

	scanner := bufio.NewScanner(stdout)
	total := 0
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "pulling") {
			total++
			prog := float64(total) * 10
			if prog > 90 {
				prog = 90
			}
			sp.Set(prog, fmt.Sprintf("Pulling %s...", extractDigest(line)))
		} else if strings.Contains(line, "verifying") {
			sp.Set(92, "Verifying...")
		} else if strings.Contains(line, "writing") {
			sp.Set(95, "Writing model...")
		}
		time.Sleep(100 * time.Millisecond)
	}

	if err := cmd.Wait(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			sp.Fail(fmt.Sprintf("Download timed out after %ds", ollamaTimeout))
		} else {
			sp.Fail(fmt.Sprintf("ollama pull: %v", err))
		}
		return
	}
	sp.Finish()
}

func extractDigest(line string) string {
	parts := strings.Fields(line)
	if len(parts) >= 2 {
		d := parts[len(parts)-1]
		if len(d) > 12 {
			d = d[:12] + "..."
		}
		return d
	}
	return ""
}

// DeleteModel removes a model via `ollama rm`.
func (mm *ModelManager) DeleteModel(modelName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, mm.ollamaCmd, "rm", modelName).Run()
}

func ModelShortName(name string) string {
	parts := strings.SplitN(name, ":", 2)
	return parts[0]
}

type OllamaGenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaGenerateResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// Generate calls /api/generate with a context timeout.
func (mm *ModelManager) Generate(prompt string, model string) (string, error) {
	return mm.GenerateWithTimeout(prompt, model, ollamaTimeout)
}

// GenerateWithTimeout calls /api/generate with a custom timeout in seconds.
func (mm *ModelManager) GenerateWithTimeout(prompt string, model string, timeoutSec int) (string, error) {
	body := OllamaGenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}
	data, _ := json.Marshal(body)

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "POST", mm.baseURL+"/api/generate", bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := mm.client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("AI generation timed out after %ds", timeoutSec)
		}
		return "", fmt.Errorf("ollama api: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama returned status %d: %s", resp.StatusCode, strings.TrimSpace(string(raw)))
	}

	var result OllamaGenerateResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}
	return result.Response, nil
}
