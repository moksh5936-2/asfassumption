package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type ModelManager struct {
	ollamaCmd string
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

func NewModelManager() *ModelManager {
	cmd := findOllama()
	return &ModelManager{ollamaCmd: cmd}
}

func findOllama() string {
	for _, p := range []string{"/usr/local/bin/ollama", "/opt/homebrew/bin/ollama", "ollama"} {
		if _, err := exec.LookPath(p); err == nil {
			return p
		}
	}
	return "ollama"
}

func (mm *ModelManager) CheckAvailable() bool {
	return exec.Command(mm.ollamaCmd, "--version").Run() == nil
}

func (mm *ModelManager) ListInstalled() ([]string, error) {
	var stdout bytes.Buffer
	cmd := exec.Command(mm.ollamaCmd, "list")
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

func (mm *ModelManager) StartDownload(modelName string, sp *syncProgress) {
	if !mm.CheckAvailable() {
		sp.Fail("Ollama not found. Install from https://ollama.ai")
		return
	}

	cmd := exec.Command(mm.ollamaCmd, "pull", modelName)
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
		sp.Fail(fmt.Sprintf("ollama pull: %v", err))
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

func (mm *ModelManager) DeleteModel(modelName string) error {
	return exec.Command(mm.ollamaCmd, "rm", modelName).Run()
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

func (mm *ModelManager) Generate(prompt string, model string) (string, error) {
	body := OllamaGenerateRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}
	data, _ := json.Marshal(body)

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("ollama api: %w", err)
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("read response: %w", err)
	}

	var result OllamaGenerateResponse
	if err := json.Unmarshal(raw, &result); err != nil {
		return "", fmt.Errorf("parse response: %w", err)
	}
	return result.Response, nil
}
