package main

import (
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeTempFile(t *testing.T, dir, pattern, content string) string {
	t.Helper()
	path := filepath.Join(dir, pattern)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}
	return path
}

func TestParseDrawio_Valid(t *testing.T) {
	xml := `<?xml version="1.0" encoding="UTF-8"?>
<mxfile><diagram><mxGraphModel><root>
<mxCell id="0" />
<mxCell id="1" parent="0" />
<mxCell id="c1" value="WebServer" style="rounded=1;" vertex="1" parent="1">
  <mxGeometry x="10" y="10" width="100" height="50" as="geometry" />
</mxCell>
<mxCell id="c2" value="Database" style="rounded=1;" vertex="1" parent="1">
  <mxGeometry x="200" y="10" width="100" height="50" as="geometry" />
</mxCell>
<mxCell id="e1" value="TLS" source="c1" target="c2" edge="1" parent="1">
  <mxGeometry relative="1" as="geometry" />
</mxCell>
</root></mxGraphModel></diagram></mxfile>`
	dir := t.TempDir()
	path := writeTempFile(t, dir, "test.drawio", xml)

	desc, err := parseDrawio(path)
	if err != nil {
		t.Fatalf("parseDrawio: %v", err)
	}
	if desc.Name != "test.drawio" {
		t.Errorf("Name = %q, want %q", desc.Name, "test.drawio")
	}
	if len(desc.Components) != 2 {
		t.Fatalf("got %d components, want 2", len(desc.Components))
	}
	if desc.Components[0].Label != "WebServer" {
		t.Errorf("Component[0] = %q, want %q", desc.Components[0].Label, "WebServer")
	}
	if desc.Components[1].Label != "Database" {
		t.Errorf("Component[1] = %q, want %q", desc.Components[1].Label, "Database")
	}
	if len(desc.Relationships) != 1 {
		t.Fatalf("got %d relationships, want 1", len(desc.Relationships))
	}
	if desc.Relationships[0].Label != "TLS" {
		t.Errorf("Relationship label = %q, want %q", desc.Relationships[0].Label, "TLS")
	}
	if desc.RawText == "" {
		t.Error("RawText is empty")
	}
}

func TestParseDrawio_Gzipped(t *testing.T) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	gw.Write([]byte(`<mxfile><diagram><mxGraphModel><root><mxCell id="c1" value="API" style="rounded=1" vertex="1" parent="1"><mxGeometry x="10" y="10" width="80" height="40" as="geometry"/></mxCell></root></mxGraphModel></diagram></mxfile>`))
	gw.Close()

	dir := t.TempDir()
	path := writeTempFile(t, dir, "test.drawio", buf.String())

	desc, err := parseDrawio(path)
	if err != nil {
		t.Fatalf("parseDrawio (gzipped): %v", err)
	}
	if len(desc.Components) == 0 {
		t.Error("expected at least one component from gzipped drawio")
	}
}

func TestParseDrawio_Malformed(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "bad.drawio", "not xml at all {{{")

	_, err := parseDrawio(path)
	if err == nil {
		t.Fatal("expected error for malformed drawio")
	}
}

func TestParseDrawio_Empty(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "empty.drawio", "")

	_, err := parseDrawio(path)
	if err == nil {
		t.Fatal("expected error for empty drawio file")
	}
}

func TestParseMermaid_Valid(t *testing.T) {
	mmd := "graph TD\nA[WebServer]\nB[Database]\nA-->|HTTP|B"
	dir := t.TempDir()
	path := writeTempFile(t, dir, "test.mmd", mmd)

	desc, err := parseMermaid(path)
	if err != nil {
		t.Fatalf("parseMermaid: %v", err)
	}
	if len(desc.Components) == 0 {
		t.Error("expected at least one component")
	}
	if len(desc.Relationships) == 0 {
		t.Error("expected at least one relationship")
	}
	if desc.RawText == "" {
		t.Error("RawText is empty")
	}
}

func TestParseMermaid_Malformed(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "bad.mmd", "this is not valid mermaid syntax @@@")

	desc, err := parseMermaid(path)
	if err != nil {
		t.Fatalf("parseMermaid should not error on malformed input: %v", err)
	}
	// Mermaid parser is lenient — should return empty but not error
	if desc == nil {
		t.Fatal("got nil desc")
	}
}

func TestParseMermaid_Empty(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "empty.mmd", "")

	desc, err := parseMermaid(path)
	if err != nil {
		t.Fatalf("parseMermaid (empty): %v", err)
	}
	if desc == nil {
		t.Fatal("got nil desc")
	}
}

func TestParseSVG_Valid(t *testing.T) {
	svg := `<?xml version="1.0"?>
<svg xmlns="http://www.w3.org/2000/svg">
  <text x="10" y="20">WebServer</text>
  <rect x="10" y="10" width="80" height="40"/>
  <text x="200" y="20">Database</text>
  <rect x="200" y="10" width="80" height="40"/>
  <line x1="90" y1="30" x2="200" y2="30"/>
</svg>`
	dir := t.TempDir()
	path := writeTempFile(t, dir, "test.svg", svg)

	desc, err := parseSVG(path)
	if err != nil {
		t.Fatalf("parseSVG: %v", err)
	}
	foundWeb := false
	foundDB := false
	for _, c := range desc.Components {
		if strings.Contains(c.Label, "WebServer") {
			foundWeb = true
		}
		if strings.Contains(c.Label, "Database") {
			foundDB = true
		}
	}
	if !foundWeb {
		t.Error("expected component containing 'WebServer'")
	}
	if !foundDB {
		t.Error("expected component containing 'Database'")
	}
}

func TestParseSVG_Empty(t *testing.T) {
	svg := `<?xml version="1.0"?><svg xmlns="http://www.w3.org/2000/svg"></svg>`
	dir := t.TempDir()
	path := writeTempFile(t, dir, "empty.svg", svg)

	desc, err := parseSVG(path)
	if err != nil {
		t.Fatalf("parseSVG (empty): %v", err)
	}
	if desc == nil {
		t.Fatal("got nil desc")
	}
}

func TestParseTextFile_Empty(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "empty.txt", "")

	desc, err := parseTextFile(path)
	if err != nil {
		t.Fatalf("parseTextFile: %v", err)
	}
	if desc.Name == "" {
		t.Error("Name should not be empty")
	}
}

func TestParseTextFile_Content(t *testing.T) {
	content := "The application uses TLS 1.3 for all communications.\nDatabase is encrypted at rest."
	dir := t.TempDir()
	path := writeTempFile(t, dir, "test.txt", content)

	desc, err := parseTextFile(path)
	if err != nil {
		t.Fatalf("parseTextFile: %v", err)
	}
	if desc.RawText != content {
		t.Errorf("RawText = %q, want %q", desc.RawText, content)
	}
}

func TestParseYAMLArch_Malformed(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "bad.yaml", "name: broken\n  indentation: bad\n  : :")

	_, err := parseYAMLArch(path)
	if err == nil {
		t.Fatal("expected error for malformed yaml")
	}
}

func TestParseJSONArch_Malformed(t *testing.T) {
	dir := t.TempDir()
	path := writeTempFile(t, dir, "bad.json", "{this is not json}")

	_, err := parseJSONArch(path)
	if err == nil {
		t.Fatal("expected error for malformed json")
	}
}

func TestParseYAMLArch_Valid(t *testing.T) {
	yaml := `name: TestArch
description: A test architecture
components:
  - name: WebServer
    type: server
    description: Main web server
  - name: Database
    type: database
    description: Primary database
relationships:
  - source: WebServer
    target: Database
    protocol: TLS`
	dir := t.TempDir()
	path := writeTempFile(t, dir, "test.yaml", yaml)

	desc, err := parseYAMLArch(path)
	if err != nil {
		t.Fatalf("parseYAMLArch: %v", err)
	}
	if len(desc.Components) != 2 {
		t.Errorf("got %d components, want 2", len(desc.Components))
	}
	if len(desc.Relationships) != 1 {
		t.Errorf("got %d relationships, want 1", len(desc.Relationships))
	}
}

func TestParseJSONArch_Valid(t *testing.T) {
	j := `{
		"name": "TestArch",
		"components": [
			{"name": "Web", "type": "server", "description": "web"},
			{"name": "DB", "type": "database", "description": "db"}
		],
		"relationships": [
			{"source": "Web", "target": "DB", "protocol": "TLS"}
		]
	}`
	dir := t.TempDir()
	path := writeTempFile(t, dir, "test.json", j)

	desc, err := parseJSONArch(path)
	if err != nil {
		t.Fatalf("parseJSONArch: %v", err)
	}
	if len(desc.Components) != 2 {
		t.Errorf("got %d components, want 2", len(desc.Components))
	}
	if len(desc.Relationships) != 1 {
		t.Errorf("got %d relationships, want 1", len(desc.Relationships))
	}
}

func TestArchDefinition_StructuredFields(t *testing.T) {
	yaml := `name: Structured
metadata:
  name: Structured
  version: "1.0"
  purpose: testing
  compliance:
    - SOC2
    - ISO27001
assumptions:
  - TLS is enforced for all external communications
  - MFA is required for admin access
security_controls:
  authentication:
    - MFA
    - OAuth2
  encryption:
    - AES-256
    - TLS 1.3
expected_results:
  minimum_assumptions: 5
  minimum_critical: 1
validation_criteria:
  - All external connections are encrypted
  - Admin access requires MFA
notes:
  - Architecture reviewed on 2024-01-01`
	dir := t.TempDir()
	path := writeTempFile(t, dir, "structured.yaml", yaml)

	desc, err := parseYAMLArch(path)
	if err != nil {
		t.Fatalf("parseYAMLArch: %v", err)
	}
	if len(desc.Compliance) == 0 {
		t.Error("expected compliance items")
	}
	if len(desc.ExplicitAssumptions) == 0 {
		t.Error("expected explicit assumptions")
	}
	if len(desc.SecurityControls) == 0 {
		t.Error("expected security controls")
	}
	if desc.ExpectedResults == nil {
		t.Error("expected expected results")
	}
	if len(desc.ValidationCriteria) == 0 {
		t.Error("expected validation criteria")
	}
	if len(desc.Notes) == 0 {
		t.Error("expected notes")
	}
}
