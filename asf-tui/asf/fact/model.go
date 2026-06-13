package fact

// Fact represents an explicit statement from the architecture.
// Facts are ground truth. They must be preserved.
// Assumptions must NOT contradict facts.
type Fact struct {
	ID          string  `json:"id"`
	Text        string  `json:"text"`
	Source      string  `json:"source"` // "yaml", "json", "mermaid", "text", "markdown"
	Confidence  float64 `json:"confidence"`
	Category    string  `json:"category"`               // "security", "compliance", "operational", "infrastructure"
	FactType    string  `json:"fact_type"`              // "control", "requirement", "constraint", "policy", "compliance", "configuration"
	ComponentID string  `json:"component_id,omitempty"` // which component this fact belongs to
	IsNegative  bool    `json:"is_negative"`            // "MFA disabled" -> true
	Severity    string  `json:"severity,omitempty"`     // "critical", "high", "medium", "low"
}

// Inventory holds all facts extracted from an architecture.
type Inventory struct {
	Facts       []Fact            `json:"facts"`
	ByType      map[string][]Fact `json:"by_type,omitempty"`
	ByCategory  map[string][]Fact `json:"by_category,omitempty"`
	ByComponent map[string][]Fact `json:"by_component,omitempty"`
}

// NewInventory creates a new empty inventory.
func NewInventory() *Inventory {
	return &Inventory{
		Facts:       make([]Fact, 0),
		ByType:      make(map[string][]Fact),
		ByCategory:  make(map[string][]Fact),
		ByComponent: make(map[string][]Fact),
	}
}

// Add adds a fact to the inventory.
func (inv *Inventory) Add(f Fact) {
	inv.Facts = append(inv.Facts, f)
	inv.ByType[f.FactType] = append(inv.ByType[f.FactType], f)
	inv.ByCategory[f.Category] = append(inv.ByCategory[f.Category], f)
	if f.ComponentID != "" {
		inv.ByComponent[f.ComponentID] = append(inv.ByComponent[f.ComponentID], f)
	}
}

// Count returns the number of facts.
func (inv *Inventory) Count() int {
	return len(inv.Facts)
}

// HasType checks if any fact of the given type exists.
func (inv *Inventory) HasType(factType string) bool {
	fs, ok := inv.ByType[factType]
	return ok && len(fs) > 0
}

// FindByType returns facts of a given type.
func (inv *Inventory) FindByType(factType string) []Fact {
	return inv.ByType[factType]
}

// FindByComponent returns facts for a given component.
func (inv *Inventory) FindByComponent(componentID string) []Fact {
	return inv.ByComponent[componentID]
}

// FidelityScore computes the ratio of facts that are respected (not contradicted).
// Returns 0.0-1.0.
func (inv *Inventory) FidelityScore(respected int, contradicted int) float64 {
	total := len(inv.Facts)
	if total == 0 {
		return 1.0
	}
	return float64(respected) / float64(total)
}
