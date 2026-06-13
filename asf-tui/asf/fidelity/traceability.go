package fidelity

import (
	"asf-tui/asf/fact"
	"fmt"
	"strings"
)

// TraceabilityRecord tracks where an assumption came from.
type TraceabilityRecord struct {
	AssumptionID       string   `json:"assumption_id"`
	SourceType         string   `json:"source_type"` // "fact-derived", "relationship-derived", "component-derived", "domain-derived"
	SourceFactID       string   `json:"source_fact_id,omitempty"`
	SourceFactText     string   `json:"source_fact_text,omitempty"`
	SourceComponent    string   `json:"source_component,omitempty"`
	SourceRelationship string   `json:"source_relationship,omitempty"`
	Reason             string   `json:"reason"`
	Evidence           []string `json:"evidence,omitempty"`
}

// TraceabilityEngine adds traceability to assumptions.
type TraceabilityEngine struct{}

// NewTraceabilityEngine creates a new traceability engine.
func NewTraceabilityEngine() *TraceabilityEngine {
	return &TraceabilityEngine{}
}

// BuildTraceability builds traceability records for assumptions.
func (e *TraceabilityEngine) BuildTraceability(assumptions []HiddenAssumption, facts *fact.Inventory) []TraceabilityRecord {
	var records []TraceabilityRecord

	for _, a := range assumptions {
		record := TraceabilityRecord{
			AssumptionID: a.ID,
			SourceType:   a.SourceType,
			Reason:       a.Reason,
		}

		// Build evidence
		if a.SourceFactID != "" {
			record.SourceFactID = a.SourceFactID
			record.SourceFactText = a.SourceFactText
			record.Evidence = append(record.Evidence, fmt.Sprintf("Source fact: %s", a.SourceFactText))
		}

		if a.ComponentID != "" {
			record.SourceComponent = a.ComponentID
			record.Evidence = append(record.Evidence, fmt.Sprintf("Source component: %s (%s)", a.ComponentID, a.ComponentLabel))
		}

		records = append(records, record)
	}

	return records
}

// FormatTraceability formats a traceability record for display.
func (e *TraceabilityEngine) FormatTraceability(record TraceabilityRecord) string {
	var parts []string
	parts = append(parts, fmt.Sprintf("Assumption: %s", record.AssumptionID))
	parts = append(parts, fmt.Sprintf("Source: %s", record.SourceType))
	if record.SourceFactText != "" {
		parts = append(parts, fmt.Sprintf("Source Fact: %s", record.SourceFactText))
	}
	if record.SourceComponent != "" {
		parts = append(parts, fmt.Sprintf("Source Component: %s", record.SourceComponent))
	}
	parts = append(parts, fmt.Sprintf("Reason: %s", record.Reason))
	if len(record.Evidence) > 0 {
		parts = append(parts, fmt.Sprintf("Evidence: %s", strings.Join(record.Evidence, "; ")))
	}

	return strings.Join(parts, "\n")
}

// ValidateTraceability checks if traceability is complete for all assumptions.
func (e *TraceabilityEngine) ValidateTraceability(assumptions []HiddenAssumption) (valid int, invalid int, missing []string) {
	for _, a := range assumptions {
		if a.Reason == "" {
			invalid++
			missing = append(missing, a.ID)
			continue
		}

		// Check if source type is valid
		validSources := map[string]bool{
			"fact-derived":           true,
			"relationship-derived":   true,
			"component-derived":      true,
			"domain-derived":         true,
			"trust-boundary-derived": true,
		}

		if !validSources[a.SourceType] {
			invalid++
			missing = append(missing, a.ID)
			continue
		}

		valid++
	}

	return valid, invalid, missing
}
