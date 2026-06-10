package evidence

import (
	"testing"
)

func TestFindField(t *testing.T) {
	rec := map[string]interface{}{
		"User":   "alice",
		"GROUp":  "Finance",
		"ACCESS": "read",
	}
	if v := FindField(rec, []string{"user", "username"}); v != "alice" {
		t.Errorf("Expected 'alice', got '%s'", v)
	}
	if v := FindField(rec, []string{"group", "department"}); v != "finance" {
		t.Errorf("Expected 'finance', got '%s'", v)
	}
	if v := FindField(rec, []string{"permission"}); v != "" {
		t.Errorf("Expected '', got '%s'", v)
	}
}

func TestFindFieldMissing(t *testing.T) {
	rec := map[string]interface{}{"name": "bob"}
	if v := FindField(rec, []string{"user", "username", "identity"}); v != "" {
		t.Errorf("Expected empty string, got '%s'", v)
	}
}

func TestGetCompatibleSourceTypes(t *testing.T) {
	m := NewMapper()
	types := m.GetCompatibleSourceTypes("ACCESS")
	if len(types) == 0 {
		t.Error("Expected compatible types for ACCESS")
	}
}

func TestLoader(t *testing.T) {
	l := NewLoader()
	if l == nil {
		t.Fatal("Expected non-nil loader")
	}
}
