package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

type TelemetryEvent struct {
	Event     string            `json:"event"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version"`
	OS        string            `json:"os"`
	Arch      string            `json:"arch"`
	Props     map[string]string `json:"props,omitempty"`
}

type TelemetryStore struct {
	mu       sync.Mutex
	Enabled  bool             `json:"enabled"`
	Events   []TelemetryEvent `json:"events"`
	filePath string
}

var telemetry *TelemetryStore
var telemetryOnce sync.Once

func getTelemetryStore() *TelemetryStore {
	telemetryOnce.Do(func() {
		store := &TelemetryStore{
			Enabled:  false,
			Events:   make([]TelemetryEvent, 0),
			filePath: asfTelemetryPath(),
		}
		if data, err := os.ReadFile(store.filePath); err == nil {
			json.Unmarshal(data, store)
		}
		telemetry = store
	})
	return telemetry
}

func initTelemetry(cfg *Config) {
	store := getTelemetryStore()
	store.mu.Lock()
	defer store.mu.Unlock()

	store.Enabled = cfg != nil && cfg.Telemetry.OptIn

	if store.Enabled {
		recordTelemetrySync(store, TelemetryEvent{
			Event:     "session_start",
			Timestamp: time.Now(),
			Version:   ASFVersion,
			OS:        runtime.GOOS,
			Arch:      runtime.GOARCH,
			Props: map[string]string{
				"mode": "tui",
			},
		})
	}
}

func recordTelemetrySync(store *TelemetryStore, event TelemetryEvent) {
	store.Events = append(store.Events, event)
	maxEvents := 1000
	if len(store.Events) > maxEvents {
		store.Events = store.Events[len(store.Events)-maxEvents:]
	}
	saveTelemetrySync(store)
}

func saveTelemetrySync(store *TelemetryStore) {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		debugLog.Printf("telemetry marshal: %v", err)
		return
	}
	if err := os.MkdirAll(filepath.Dir(store.filePath), 0755); err != nil {
		debugLog.Printf("telemetry mkdir: %v", err)
		return
	}
	if err := os.WriteFile(store.filePath, data, 0644); err != nil {
		debugLog.Printf("telemetry write: %v", err)
	}
}

func recordTelemetry(event string, props map[string]string) {
	store := getTelemetryStore()
	store.mu.Lock()
	defer store.mu.Unlock()

	if !store.Enabled {
		return
	}

	recordTelemetrySync(store, TelemetryEvent{
		Event:     event,
		Timestamp: time.Now(),
		Version:   ASFVersion,
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		Props:     props,
	})
}

func TelemetrySummary() string {
	store := getTelemetryStore()
	store.mu.Lock()
	defer store.mu.Unlock()

	if !store.Enabled {
		return "Telemetry: disabled"
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("Telemetry: enabled (%d events)\n", len(store.Events)))

	eventCounts := make(map[string]int)
	for _, e := range store.Events {
		eventCounts[e.Event]++
	}

	for event, count := range eventCounts {
		b.WriteString(fmt.Sprintf("  %s: %d\n", event, count))
	}

	return b.String()
}

func init() {
	store := getTelemetryStore()
	if store.Enabled {
		recordTelemetry("app_start", map[string]string{
			"version": ASFVersion,
		})
	}
}
