package ingestion

import (
	"archive/zip"
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"asf-tui/asf/models"
	"github.com/ledongthuc/pdf"
)

type Pipeline struct{}

func NewPipeline() *Pipeline {
	return &Pipeline{}
}

func (p *Pipeline) DetectType(filepathStr string) models.SourceType {
	ext := strings.ToLower(filepath.Ext(filepathStr))
	switch ext {
	case ".pdf":
		return models.SourceTypePDF
	case ".docx":
		return models.SourceTypeDOCX
	case ".txt":
		return models.SourceTypeTXT
	case ".csv":
		return models.SourceTypeCSV
	case ".json":
		return models.SourceTypeJSON
	default:
		return models.SourceTypeUNKNOWN
	}
}

func (p *Pipeline) ParseText(filepathStr string) (string, error) {
	ftype := p.DetectType(filepathStr)
	switch ftype {
	case models.SourceTypePDF:
		return parsePDF(filepathStr)
	case models.SourceTypeDOCX:
		return parseDOCX(filepathStr)
	case models.SourceTypeTXT:
		return parseTXT(filepathStr)
	case models.SourceTypeCSV:
		return parseCSVToText(filepathStr)
	case models.SourceTypeJSON:
		return parseJSONToText(filepathStr)
	default:
		return parseTXT(filepathStr)
	}
}

func (p *Pipeline) ParseToRecords(filepathStr string) ([]map[string]interface{}, error) {
	ftype := p.DetectType(filepathStr)
	switch ftype {
	case models.SourceTypeCSV:
		return parseCSVToRecords(filepathStr)
	case models.SourceTypeJSON:
		return parseJSONToRecords(filepathStr)
	default:
		return nil, nil
	}
}

func parseTXT(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func parseCSVToText(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	r := csv.NewReader(f)
	records, err := r.ReadAll()
	if err != nil {
		return "", err
	}

	var sb strings.Builder
	for i, row := range records {
		if i == 0 {
			sb.WriteString("headers: " + strings.Join(row, ", "))
		} else {
			sb.WriteString(strings.Join(row, ", "))
		}
		sb.WriteString("\n")
	}
	return sb.String(), nil
}

func parseCSVToRecords(path string) ([]map[string]interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	r := csv.NewReader(f)
	rows, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(rows) < 2 {
		return nil, nil
	}

	headers := rows[0]
	var result []map[string]interface{}
	for _, row := range rows[1:] {
		rec := make(map[string]interface{})
		for i, h := range headers {
			if i < len(row) {
				rec[h] = row[i]
			}
		}
		result = append(result, rec)
	}
	return result, nil
}

func parseJSONToText(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func parseJSONToRecords(path string) ([]map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(data, &result); err == nil {
		return result, nil
	}

	var single map[string]interface{}
	if err := json.Unmarshal(data, &single); err == nil {
		return []map[string]interface{}{single}, nil
	}

	return nil, fmt.Errorf("cannot parse JSON: not an array or object")
}

func parsePDF(path string) (string, error) {
	f, r, err := pdf.Open(path)
	if err != nil {
		return "", fmt.Errorf("pdf open: %w", err)
	}
	defer f.Close()

	plainTextReader, err := r.GetPlainText()
	if err != nil {
		return "", fmt.Errorf("pdf get text: %w", err)
	}

	data, err := io.ReadAll(plainTextReader)
	if err != nil {
		return "", fmt.Errorf("pdf read: %w", err)
	}

	return string(data), nil
}

type wDoc struct {
	Body wBody `xml:"http://schemas.openxmlformats.org/wordprocessingml/2006/main body"`
}

type wBody struct {
	Paragraphs []wParagraph `xml:"http://schemas.openxmlformats.org/wordprocessingml/2006/main p"`
}

type wParagraph struct {
	Runs []wRun `xml:"http://schemas.openxmlformats.org/wordprocessingml/2006/main r"`
}

type wRun struct {
	Text string `xml:"http://schemas.openxmlformats.org/wordprocessingml/2006/main t"`
}

func parseDOCX(path string) (string, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return "", fmt.Errorf("docx open: %w", err)
	}
	defer r.Close()

	var xmlData []byte
	for _, f := range r.File {
		if f.Name == "word/document.xml" {
			rc, err := f.Open()
			if err != nil {
				return "", fmt.Errorf("docx read document.xml: %w", err)
			}
			defer rc.Close()

			xmlData, err = io.ReadAll(rc)
			if err != nil {
				return "", fmt.Errorf("docx read: %w", err)
			}
			break
		}
	}

	if xmlData == nil {
		return "", fmt.Errorf("docx: word/document.xml not found")
	}

	var doc wDoc
	if err := xml.Unmarshal(xmlData, &doc); err != nil {
		return "", fmt.Errorf("docx xml parse: %w", err)
	}

	var text strings.Builder
	for _, p := range doc.Body.Paragraphs {
		for _, r := range p.Runs {
			text.WriteString(r.Text)
		}
		text.WriteString("\n")
	}
	return text.String(), nil
}
