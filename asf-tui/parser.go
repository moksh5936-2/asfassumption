package main

import (
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/ledongthuc/pdf"
	"gopkg.in/yaml.v3"
)

type Component struct {
	ID    string
	Label string
}

type Relation struct {
	Source string
	Target string
	Label  string
}

type ArchDescription struct {
	Name          string
	Components    []Component
	Relationships []Relation
	Policies      []string
	RawText       string

	// Structured fields (populated from YAML/JSON definitions)
	ExplicitAssumptions []string
	SecurityControls    map[string][]string
	Compliance          []string
	ExpectedResults     map[string]interface{}
	ValidationCriteria  []string
	Notes               []string
}

func ParseArchitecture(path string) (*ArchDescription, error) {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".drawio":
		return parseDrawio(path)
	case ".mmd":
		return parseMermaid(path)
	case ".yaml", ".yml":
		return parseYAMLArch(path)
	case ".json":
		return parseJSONArch(path)
	case ".svg":
		return parseSVG(path)
	case ".png", ".jpg", ".jpeg":
		return parseImageOCR(path)
	case ".pdf":
		return parsePDF(path)
	case ".docx":
		return parseDOCX(path)
	case ".txt":
		return parseTextFile(path)
	case ".md":
		return parseMarkdown(path)
	default:
		return parseTextFile(path)
	}
}

func parseDrawio(path string) (*ArchDescription, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read drawio: %w", err)
	}

	if isGzipped(data) {
		reader, err := gzip.NewReader(bytes.NewReader(data))
		if err != nil {
			return nil, fmt.Errorf("decompress drawio (corrupt gzip): %w", err)
		}
		var buf bytes.Buffer
		_, err = buf.ReadFrom(reader)
		reader.Close()
		if err != nil {
			return nil, fmt.Errorf("decompress drawio: %w", err)
		}
		data = buf.Bytes()
	}

	var mxfile mxFile
	if err := xml.Unmarshal(data, &mxfile); err != nil {
		return nil, fmt.Errorf("parse drawio xml: %w", err)
	}

	desc := &ArchDescription{
		Name: fileBase(path),
	}

	cellMap := make(map[string]string)
	var edges []struct {
		src string
		tgt string
		lbl string
	}

	for _, diagram := range mxfile.Diagrams {
		for _, cell := range diagram.Graph.Root.Cells {
			id := cell.ID
			label := strings.TrimSpace(cell.Value)
			if label == "" {
				label = id
			}

			if cell.Source != "" || cell.Target != "" {
				srcLabel := cellMap[cell.Source]
				tgtLabel := cellMap[cell.Target]
				if srcLabel == "" {
					srcLabel = cell.Source
				}
				if tgtLabel == "" {
					tgtLabel = cell.Target
				}
				edges = append(edges, struct {
					src string
					tgt string
					lbl string
				}{srcLabel, tgtLabel, label})
			} else if !isStyleNoLabel(cell.Style) {
				desc.Components = append(desc.Components, Component{ID: id, Label: label})
				cellMap[id] = label
			}
		}
	}

	for _, e := range edges {
		desc.Relationships = append(desc.Relationships, Relation{
			Source: e.src,
			Target: e.tgt,
			Label:  e.lbl,
		})
	}

	desc.RawText = buildTextFromDiagram(desc.Name, desc.Components, desc.Relationships)
	return desc, nil
}

type mxFile struct {
	XMLName  xml.Name    `xml:"mxfile"`
	Diagrams []mxDiagram `xml:"diagram"`
}

type mxDiagram struct {
	Graph mxGraphModel `xml:"mxGraphModel"`
}

type mxGraphModel struct {
	Root mxRoot `xml:"root"`
}

type mxRoot struct {
	Cells []mxCell `xml:"mxCell"`
}

type mxCell struct {
	ID     string `xml:"id,attr"`
	Value  string `xml:"value,attr"`
	Style  string `xml:"style,attr"`
	Source string `xml:"source,attr"`
	Target string `xml:"target,attr"`
	Vertex int    `xml:"vertex,attr"`
	Edge   int    `xml:"edge,attr"`
}

func isStyleNoLabel(style string) bool {
	return strings.Contains(style, "ellipse") ||
		strings.Contains(style, "rhombus") ||
		style == ""
}

func isGzipped(data []byte) bool {
	return len(data) > 2 && data[0] == 0x1f && data[1] == 0x8b
}

func parseMermaid(path string) (*ArchDescription, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read mermaid: %w", err)
	}
	content := string(data)
	desc := &ArchDescription{Name: fileBase(path)}

	nodeRe := regexp.MustCompile(`([A-Za-z0-9_]+)\[(.*?)\]`)
	edgeRe := regexp.MustCompile(`([A-Za-z0-9_]+)\s*-->`)

	nodeNames := make(map[string]string)
	for _, m := range nodeRe.FindAllStringSubmatch(content, -1) {
		id := m[1]
		label := strings.TrimSpace(m[2])
		label = strings.Trim(label, "(){}")
		if label == "" {
			label = id
		}
		nodeNames[id] = label
		desc.Components = append(desc.Components, Component{ID: id, Label: label})
	}

	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "-->") {
			parts := edgeRe.FindStringSubmatch(line)
			if len(parts) < 2 {
				continue
			}
			srcID := parts[1]

			remainder := line[strings.Index(line, "-->")+3:]
			tgtID := extractMermaidNodeID(remainder)

			src := nodeNames[srcID]
			if src == "" {
				src = srcID
			}
			tgt := nodeNames[tgtID]
			if tgt == "" {
				tgt = tgtID
			}

			label := extractMermaidEdgeLabel(line)
			desc.Relationships = append(desc.Relationships, Relation{
				Source: src, Target: tgt, Label: label,
			})
		}
	}

	desc.RawText = buildTextFromDiagram(desc.Name, desc.Components, desc.Relationships)
	return desc, nil
}

func extractMermaidNodeID(s string) string {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "|") {
		if idx := strings.Index(s[1:], "|"); idx >= 0 {
			s = strings.TrimSpace(s[idx+2:])
		}
	}
	if idx := strings.IndexAny(s, " \t["); idx >= 0 {
		s = s[:idx]
	}
	return s
}

func extractMermaidEdgeLabel(line string) string {
	pipeRe := regexp.MustCompile(`\|(.*?)\|`)
	matches := pipeRe.FindStringSubmatch(line)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

func parsePDF(path string) (*ArchDescription, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open pdf: %w", err)
	}
	defer f.Close()

	reader, err := r.GetPlainText()
	if err != nil {
		return nil, fmt.Errorf("extract pdf text: %w", err)
	}
	text, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read pdf text: %w", err)
	}

	return &ArchDescription{
		Name:    fileBase(path),
		RawText: string(text),
	}, nil
}

func parseDOCX(path string) (*ArchDescription, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, fmt.Errorf("open docx: %w", err)
	}
	defer r.Close()

	var documentXML []byte
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return nil, fmt.Errorf("read word/document.xml: %w", err)
			}
			documentXML, err = io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return nil, fmt.Errorf("read word/document.xml content: %w", err)
			}
			break
		}
	}
	if documentXML == nil {
		return nil, fmt.Errorf("word/document.xml not found in docx archive")
	}

	type wDoc struct {
		Body struct {
			Paragraphs []struct {
				Runs []struct {
					Text string `xml:"t"`
				} `xml:"r"`
			} `xml:"p"`
		} `xml:"body"`
	}

	var doc wDoc
	if err := xml.Unmarshal(documentXML, &doc); err != nil {
		return nil, fmt.Errorf("parse word/document.xml: %w", err)
	}

	var lines []string
	for _, p := range doc.Body.Paragraphs {
		var line []string
		for _, run := range p.Runs {
			line = append(line, run.Text)
		}
		lines = append(lines, strings.Join(line, ""))
	}

	return &ArchDescription{
		Name:    fileBase(path),
		RawText: strings.Join(lines, "\n"),
	}, nil
}

func parseTextFile(path string) (*ArchDescription, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}
	rawText := string(data)
	return &ArchDescription{
		Name:    fileBase(path),
		RawText: rawText,
	}, nil
}

// parseMarkdown extracts structured sections from Markdown architecture documents.
func parseMarkdown(path string) (*ArchDescription, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read markdown: %w", err)
	}
	content := string(data)
	desc := &ArchDescription{
		Name:    fileBase(path),
		RawText: content,
	}

	lines := strings.Split(content, "\n")
	section := ""
	subHeading := ""
	inCodeBlock := false
	var currentList []string

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Track code blocks
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			continue
		}

		// Detect sub-headings (###) within a section
		if strings.HasPrefix(trimmed, "###") {
			flushSection(desc, section, subHeading, currentList)
			subHeading = strings.ToLower(strings.TrimSpace(strings.Trim(trimmed, "#")))
			currentList = nil
			continue
		}

		// Detect section headers
		if strings.HasPrefix(trimmed, "#") {
			flushSection(desc, section, subHeading, currentList)
			section = strings.ToLower(strings.TrimSpace(strings.Trim(trimmed, "#")))
			subHeading = ""
			currentList = nil
			continue
		}

		// Detect list items
		isListItem := strings.HasPrefix(trimmed, "- ") || strings.HasPrefix(trimmed, "* ") ||
			regexp.MustCompile(`^\d+\.\s+`).MatchString(trimmed)

		if isListItem && trimmed != "" {
			item := strings.TrimSpace(regexp.MustCompile(`^[\s\-*\d.]+\s+`).ReplaceAllString(trimmed, ""))
			if item != "" {
				currentList = append(currentList, item)
			}
		} else if trimmed != "" && trimmed != "|" && !strings.HasPrefix(trimmed, "|") {
			// Non-list item ends current list context
			flushSection(desc, section, subHeading, currentList)
			currentList = nil
		}
	}

	// Flush any remaining items
	flushSection(desc, section, subHeading, currentList)

	// Extract compliance from metadata section
	if desc.Compliance == nil {
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if strings.Contains(strings.ToLower(trimmed), "compliance") && strings.Contains(trimmed, ":") {
				parts := strings.SplitN(trimmed, ":", 2)
				if len(parts) == 2 {
					items := strings.Split(parts[1], ",")
					for _, item := range items {
						item = strings.TrimSpace(item)
						item = strings.TrimPrefix(item, "**")
						item = strings.TrimSuffix(item, "**")
						item = strings.TrimSpace(item)
						if item != "" {
							desc.Compliance = append(desc.Compliance, item)
						}
					}
				}
			}
		}
	}

	return desc, nil
}

// flushSection populates ArchDescription fields based on the current section.
func flushSection(desc *ArchDescription, section, subHeading string, items []string) {
	if len(items) == 0 {
		return
	}
	switch section {
	case "assumptions":
		desc.ExplicitAssumptions = append(desc.ExplicitAssumptions, items...)
	case "security controls":
		if desc.SecurityControls == nil {
			desc.SecurityControls = make(map[string][]string)
		}
		// Use sub-heading as category if available, otherwise "general"
		cat := subHeading
		if cat == "" {
			cat = "general"
		}
		desc.SecurityControls[cat] = append(desc.SecurityControls[cat], items...)
	case "compliance":
		for _, item := range items {
			item = strings.TrimSpace(item)
			item = strings.TrimPrefix(item, "-")
			item = strings.TrimPrefix(item, "**")
			item = strings.TrimPrefix(item, "*")
			item = strings.TrimSuffix(item, "**")
			item = strings.TrimSpace(item)
			if item != "" {
				desc.Compliance = append(desc.Compliance, item)
			}
		}
	case "notes":
		for _, item := range items {
			item = strings.TrimSpace(item)
			item = strings.TrimPrefix(item, "-")
			item = strings.TrimPrefix(item, "**")
			item = strings.TrimPrefix(item, "*")
			item = strings.TrimSuffix(item, "**")
			item = strings.TrimSpace(item)
			if item != "" {
				desc.Notes = append(desc.Notes, item)
			}
		}
	}
}

func buildTextFromDiagram(name string, components []Component, relations []Relation) string {
	var b strings.Builder

	b.WriteString("## Topology\n\n")
	if len(relations) > 0 {
		parts := make([]string, 0, len(relations))
		for _, r := range relations {
			if r.Label != "" && r.Label != r.Source && r.Label != r.Target {
				parts = append(parts, fmt.Sprintf("[%s] --%s--> [%s]", r.Source, r.Label, r.Target))
			} else {
				parts = append(parts, fmt.Sprintf("[%s] --> [%s]", r.Source, r.Target))
			}
		}
		b.WriteString(strings.Join(parts, "\n"))
		b.WriteString("\n\n")
	} else {
		b.WriteString("Components identified but no explicit relationships mapped.\n\n")
	}

	b.WriteString("## Components\n\n")
	for _, c := range components {
		b.WriteString(fmt.Sprintf("- %s\n", c.Label))
	}
	b.WriteString("\n")

	b.WriteString("## Documented Policy\n\n")
	for _, r := range relations {
		protocol := r.Label
		if protocol == "" || protocol == r.Source || protocol == r.Target {
			protocol = "a secure protocol"
		}
		b.WriteString(fmt.Sprintf("%s connects to %s using %s.\n", r.Source, r.Target, protocol))
		encProtocol := protocol
		if strings.EqualFold(encProtocol, "SQL") || strings.EqualFold(encProtocol, "HTTP") {
			encProtocol = "TLS"
		}
		b.WriteString(fmt.Sprintf("All communication between %s and %s MUST use %s encryption.\n", r.Source, r.Target, encProtocol))
	}
	if len(relations) == 0 {
		b.WriteString("Standard enterprise security policies apply to all components.\n")
	}
	b.WriteString("\n")

	b.WriteString("## Access Control\n\n")
	for _, c := range components {
		label := strings.ToLower(c.Label)
		if strings.Contains(label, "database") || strings.Contains(label, "db") {
			b.WriteString(fmt.Sprintf("Only authorized applications may access %s.\n", c.Label))
			b.WriteString(fmt.Sprintf("Database access to %s is restricted to database administrators only.\n", c.Label))
		} else if strings.Contains(label, "gateway") || strings.Contains(label, "proxy") || strings.Contains(label, "lb") || strings.Contains(label, "load") {
			b.WriteString(fmt.Sprintf("All traffic MUST pass through %s for authentication.\n", c.Label))
			b.WriteString(fmt.Sprintf("%s enforces access control policies.\n", c.Label))
		} else if strings.Contains(label, "auth") || strings.Contains(label, "identity") || strings.Contains(label, "sso") || strings.Contains(label, "mfa") || strings.Contains(label, "idp") {
			b.WriteString(fmt.Sprintf("All administrative access requires multi-factor authentication through %s.\n", c.Label))
		} else if strings.Contains(label, "user") || strings.Contains(label, "browser") || strings.Contains(label, "client") {
			b.WriteString(fmt.Sprintf("VPN required for remote access from %s.\n", c.Label))
			b.WriteString(fmt.Sprintf("%s authenticates with corporate credentials.\n", c.Label))
		} else if strings.Contains(label, "app") || strings.Contains(label, "server") || strings.Contains(label, "service") {
			b.WriteString(fmt.Sprintf("%s authenticates with Active Directory credentials.\n", c.Label))
		}
	}
	b.WriteString("\n")

	b.WriteString("## Network Security\n\n")
	for _, r := range relations {
		src := strings.ToLower(r.Source)
		tgt := strings.ToLower(r.Target)
		if strings.Contains(src, "internet") || strings.Contains(tgt, "internet") ||
			strings.Contains(src, "browser") || strings.Contains(tgt, "browser") ||
			strings.Contains(src, "user") || strings.Contains(tgt, "user") {
			b.WriteString(fmt.Sprintf("TLS termination required at the entry point for %s to %s connection.\n", r.Source, r.Target))
		}
		if strings.Contains(tgt, "database") || strings.Contains(tgt, "db") {
			b.WriteString(fmt.Sprintf("Database is in private subnet, not internet accessible.\n"))
		}
	}
	b.WriteString("\n")

	b.WriteString("## Trust Boundaries\n\n")
	for i, r := range relations {
		boundaryTypes := []string{"network boundary", "authentication boundary", "data boundary", "access boundary"}
		bt := boundaryTypes[i%len(boundaryTypes)]
		b.WriteString(fmt.Sprintf("- Between %s and %s (%s)\n", r.Source, r.Target, bt))
	}
	if len(relations) == 0 && len(components) > 0 {
		b.WriteString("- Between external and internal components\n")
		b.WriteString("- Between application and data layer\n")
	}
	b.WriteString("\n")

	b.WriteString("## Assumptions to Consider\n\n")
	seenType := make(map[string]bool)
	for _, c := range components {
		label := strings.ToLower(c.Label)
		if strings.Contains(label, "vpn") || strings.Contains(label, "gateway") {
			if !seenType["vpn"] {
				b.WriteString("- What if VPN is unavailable?\n")
				b.WriteString("- What if MFA is not enforced on VPN?\n")
				b.WriteString("- What if VPN Gateway credentials are compromised?\n")
				seenType["vpn"] = true
			}
		}
		if strings.Contains(label, "database") || strings.Contains(label, "db") {
			if !seenType["db"] {
				b.WriteString("- What if database credentials are leaked?\n")
				b.WriteString("- What if backup restore is untested?\n")
				b.WriteString("- What if database has a public route?\n")
				seenType["db"] = true
			}
		}
		if strings.Contains(label, "auth") || strings.Contains(label, "mfa") || strings.Contains(label, "identity") || strings.Contains(label, "idp") {
			if !seenType["auth"] {
				b.WriteString("- What if MFA provider is unavailable?\n")
				b.WriteString("- What if authentication service is down?\n")
				seenType["auth"] = true
			}
		}
		if strings.Contains(label, "app") || strings.Contains(label, "server") || strings.Contains(label, "service") {
			if !seenType["app"] {
				b.WriteString("- What if application server is unavailable?\n")
				b.WriteString("- What if application credentials are leaked?\n")
				seenType["app"] = true
			}
		}
	}
	if len(seenType) == 0 {
		b.WriteString("- What if any component is unavailable?\n")
		b.WriteString("- What if encryption between components is compromised?\n")
		b.WriteString("- What if access control policies are misconfigured?\n")
	}
	b.WriteString("\n")

	return b.String()
}

// YAML/JSON architecture definition format
type archDefinition struct {
	Name        string `yaml:"name" json:"name"`
	Description string `yaml:"description" json:"description"`
	Components  []struct {
		Name        string `yaml:"name" json:"name"`
		Type        string `yaml:"type" json:"type"`
		Description string `yaml:"description" json:"description"`
	} `yaml:"components" json:"components"`
	Relationships []struct {
		Source      string `yaml:"source" json:"source"`
		Target      string `yaml:"target" json:"target"`
		Protocol    string `yaml:"protocol" json:"protocol"`
		Description string `yaml:"description" json:"description"`
	} `yaml:"relationships" json:"relationships"`
	Policies []string `yaml:"policies" json:"policies"`

	// Structured YAML/JSON fields (populated from architecture definitions)
	Metadata *struct {
		Name       string   `yaml:"name" json:"name"`
		Version    string   `yaml:"version" json:"version"`
		Purpose    string   `yaml:"purpose" json:"purpose"`
		Compliance []string `yaml:"compliance" json:"compliance"`
	} `yaml:"metadata" json:"metadata"`
	System *struct {
		Name        string `yaml:"name" json:"name"`
		Description string `yaml:"description" json:"description"`
	} `yaml:"system" json:"system"`
	Assumptions        []string               `yaml:"assumptions" json:"assumptions"`
	SecurityControls   map[string][]string    `yaml:"security_controls" json:"security_controls"`
	ExpectedResults    map[string]interface{} `yaml:"expected_results" json:"expected_results"`
	ValidationCriteria []string               `yaml:"validation_criteria" json:"validation_criteria"`
	Notes              []string               `yaml:"notes" json:"notes"`
}

func parseYAMLArch(path string) (*ArchDescription, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read yaml: %w", err)
	}
	var def archDefinition
	if err := yaml.Unmarshal(data, &def); err != nil {
		return nil, fmt.Errorf("parse yaml: %w", err)
	}
	return buildFromDefinition(&def, path)
}

func parseJSONArch(path string) (*ArchDescription, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read json: %w", err)
	}
	var def archDefinition
	if err := json.Unmarshal(data, &def); err != nil {
		return nil, fmt.Errorf("parse json: %w", err)
	}
	return buildFromDefinition(&def, path)
}

func buildFromDefinition(def *archDefinition, path string) (*ArchDescription, error) {
	desc := &ArchDescription{
		Name:       filepath.Base(path),
		Components: make([]Component, 0),
		Policies:   def.Policies,
	}

	text := def.Description + "\n\n"
	if len(def.Components) > 0 {
		text += "## Components\n\n"
		for _, c := range def.Components {
			desc.Components = append(desc.Components, Component{
				ID:    c.Name,
				Label: c.Name,
			})
			text += fmt.Sprintf("- %s (%s): %s\n", c.Name, c.Type, c.Description)
		}
		text += "\n"
	}
	if len(def.Relationships) > 0 {
		text += "## Relationships\n\n"
		for _, r := range def.Relationships {
			proto := r.Protocol
			if proto == "" {
				proto = "a secure protocol"
			}
			desc.Relationships = append(desc.Relationships, Relation{
				Source: r.Source,
				Target: r.Target,
				Label:  proto,
			})
			text += fmt.Sprintf("%s -> %s [%s]\n", r.Source, r.Target, proto)
		}
		text += "\n"
	}
	if len(def.Policies) > 0 {
		text += "## Policies\n\n"
		for _, p := range def.Policies {
			text += fmt.Sprintf("- %s\n", p)
		}
		text += "\n"
	}

	// Extract structured fields
	if def.Metadata != nil {
		desc.Compliance = def.Metadata.Compliance
	}
	if len(def.Assumptions) > 0 {
		desc.ExplicitAssumptions = def.Assumptions
		text += "## Explicit Assumptions\n\n"
		for _, a := range def.Assumptions {
			text += fmt.Sprintf("- %s.\n", a)
		}
		text += "\n"
	}
	if len(def.SecurityControls) > 0 {
		desc.SecurityControls = def.SecurityControls
		text += "## Security Controls\n\n"
		for category, controls := range def.SecurityControls {
			text += fmt.Sprintf("- %s: %s\n", category, strings.Join(controls, ", "))
		}
		text += "\n"
	}
	if def.ExpectedResults != nil {
		desc.ExpectedResults = def.ExpectedResults
	}
	if len(def.ValidationCriteria) > 0 {
		desc.ValidationCriteria = def.ValidationCriteria
	}
	if len(def.Notes) > 0 {
		desc.Notes = def.Notes
	}

	desc.RawText = text
	if len(desc.Components) > 0 {
		desc.RawText = buildTextFromDiagram(desc.Name, desc.Components, desc.Relationships)
		// Append structured content to the diagram text for claim extraction
		if len(def.Assumptions) > 0 {
			desc.RawText += "\n## Explicit Assumptions\n\n"
			for _, a := range def.Assumptions {
				desc.RawText += fmt.Sprintf("- %s.\n", a)
			}
			desc.RawText += "\n"
		}
		if len(def.SecurityControls) > 0 {
			desc.RawText += "## Security Controls\n\n"
			for category, controls := range def.SecurityControls {
				desc.RawText += fmt.Sprintf("- %s: %s\n", category, strings.Join(controls, ", "))
			}
			desc.RawText += "\n"
		}
	}
	return desc, nil
}

// SVG parser
type svgRoot struct {
	XMLName xml.Name  `xml:"svg"`
	Texts   []svgText `xml:"text"`
	Rects   []svgRect `xml:"rect"`
	Circles []svgCirc `xml:"circle"`
	Lines   []svgLine `xml:"line"`
	Paths   []svgPath `xml:"path"`
	Groups  []svgG    `xml:"g"`
}

type svgText struct {
	Content string `xml:",chardata"`
	X       string `xml:"x,attr"`
	Y       string `xml:"y,attr"`
}

type svgRect struct {
	Width  string `xml:"width,attr"`
	Height string `xml:"height,attr"`
}

type svgCirc struct {
	R string `xml:"r,attr"`
}

type svgLine struct {
	X1 string `xml:"x1,attr"`
	Y1 string `xml:"y1,attr"`
	X2 string `xml:"x2,attr"`
	Y2 string `xml:"y2,attr"`
}

type svgPath struct {
	D string `xml:"d,attr"`
}

type svgG struct {
	Texts []svgText `xml:"text"`
	Rects []svgRect `xml:"rect"`
}

func parseSVG(path string) (*ArchDescription, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read svg: %w", err)
	}
	var root svgRoot
	if err := xml.Unmarshal(data, &root); err != nil {
		return nil, fmt.Errorf("parse svg: %w", err)
	}

	components := make(map[string]bool)
	var comps []Component
	var rels []Relation

	for _, t := range root.Texts {
		label := strings.TrimSpace(t.Content)
		if label != "" && !components[label] {
			components[label] = true
			comps = append(comps, Component{ID: label, Label: label})
		}
	}
	for _, g := range root.Groups {
		for _, t := range g.Texts {
			label := strings.TrimSpace(t.Content)
			if label != "" && !components[label] {
				components[label] = true
				comps = append(comps, Component{ID: label, Label: label})
			}
		}
	}

	pathLabels := extractLabelsFromPaths(data)

	for i := 0; i < len(comps)-1; i++ {
		rels = append(rels, Relation{
			Source: comps[i].Label,
			Target: comps[i+1].Label,
			Label:  pathLabels,
		})
	}

	desc := &ArchDescription{
		Name:       filepath.Base(path),
		Components: comps,
	}
	desc.RawText = buildTextFromDiagram(desc.Name, desc.Components, rels)
	return desc, nil
}

func extractLabelsFromPaths(data []byte) string {
	re := regexp.MustCompile(`>([^<]+)<`)
	matches := re.FindAllStringSubmatch(string(data), -1)
	var labels []string
	for _, m := range matches {
		s := strings.TrimSpace(m[1])
		if len(s) > 1 && len(s) < 40 && !strings.Contains(s, "<") && !strings.Contains(s, ">") {
			labels = append(labels, s)
		}
	}
	if len(labels) > 0 {
		return strings.Join(labels, " ")
	}
	return "a secure protocol"
}

func parseImageOCR(path string) (*ArchDescription, error) {
	tesseractPath := findTesseract()
	if tesseractPath == "" {
		return nil, fmt.Errorf("tesseract not found. Install with: brew install tesseract (macOS), apt install tesseract-ocr (Linux), or download from https://github.com/tesseract-ocr/tesseract")
	}

	cmd := exec.Command(tesseractPath, path, "stdout", "-l", "eng")
	var stdout, stderr strings.Builder
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("tesseract failed: %w (stderr: %s)", err, stderr.String())
	}

	text := strings.TrimSpace(stdout.String())
	if text == "" {
		return nil, fmt.Errorf("no text extracted from image")
	}

	desc := &ArchDescription{
		Name:    filepath.Base(path),
		RawText: fmt.Sprintf("# Architecture: %s\n\n## Extracted Text\n\n%s\n\n## Notes\n\nThis is OCR output from an architecture diagram image. Quality depends on image clarity and text legibility. For best results, use Draw.io (.drawio), Mermaid (.mmd), or SVG format.", filepath.Base(path), text),
	}
	return desc, nil
}

func findTesseract() string {
	for _, p := range []string{"/usr/local/bin/tesseract", "/opt/homebrew/bin/tesseract", "/usr/bin/tesseract", "tesseract"} {
		if _, err := exec.LookPath(p); err == nil {
			return p
		}
	}
	return ""
}
